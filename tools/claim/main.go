// Command claim reads one claim out of the ledgers instead of the whole ledger.
//
// The ledgers under claims/ hold ~1200 rows over 2.0 MB. A row is ~1.6 KB and
// carries its own evidence, so the ledgers cannot be made smaller without losing
// what a claim is for. What can change is how they are read: opening
// claims/terrain.md to learn one fact costs 234 KB of an agent's context, and
// claims/ai.md another 203 KB, when the answer is one row.
//
// AGENTS.md routes every lookup through claims/registry.md -> claims/<fmt>.md,
// and PIPELINE-STATUS.md tells a brief-writer to check every id against the
// ledgers before sending it. Both are right; neither is a reason to load a
// ledger whole.
//
//	claim TERR-LOC-001 AI-DIPLO-004   one row each, with its retraction state
//	claim -k "sight range"            search every ledger, compact hits
//	claim -l terrain                  one ledger's ids as an index, not its prose
//	claim -stats                      what each ledger costs to open
//
// Exit 1 if any requested id is missing -- an unknown id is a defect in the
// brief that asked for it, not an empty result.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// A claim id is <FMT>-<TOPIC>-<NNN>, or <FMT>-<NNN> where the ledger needs no
// topic segment. The topic segment may carry digits (SPR256, SPR16A), which is
// what makes the naive [A-Z]+-[0-9]{3} pattern miss rows -- the same defect
// that has made claims/registry.md's per-area counts go stale three times.
// Anchor on the trailing -NNN instead.
//
// THE TOPIC SEGMENT IS OPTIONAL, and requiring it cost a whole ledger. On
// 2026-08-16 claims/town.md was published with 36 two-segment ids (TOWN-001 and
// up). This pattern required at least one segment between the prefix and the
// number, so every one of those rows answered "TOWN-036 is not in any ledger" --
// the tool reported a missing claim, which AGENTS.md defines as a defect in the
// brief that asked for it, for 36 rows that were present and correct. The same
// id shape had already slipped past check-claim-ids.sh and check-retraction-
// status.sh, both widened the same day; this tool was not, because nothing had
// asked it for a TOWN id yet. It is applied to a row's id cell alone (one call
// site), so relaxing the middle group cannot over-match prose.
var idRe = regexp.MustCompile(`\b([A-Z][A-Z0-9]*(?:-[A-Z0-9]+)*-[0-9]{3})\b`)

type row struct {
	id    string
	file  string // ledger base name, e.g. "terrain.md"
	line  int
	cells []string
	raw   string
}

func main() {
	var (
		keyword = flag.String("k", "", "search every ledger for a regexp and print compact hits")
		ledger  = flag.String("l", "", "list one ledger's ids as an index (no prose)")
		stats   = flag.Bool("stats", false, "print what each ledger costs to open")
		dir     = flag.String("dir", "", "claims directory (default: ./claims, then ../claims)")
		full    = flag.Bool("full", false, "with -k or -l, print whole rows instead of an index")
	)
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: claim [flags] ID [ID...]\n\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	root, err := claimsDir(*dir)
	if err != nil {
		fmt.Fprintln(os.Stderr, "claim:", err)
		os.Exit(2)
	}

	rows, err := load(root)
	if err != nil {
		fmt.Fprintln(os.Stderr, "claim:", err)
		os.Exit(2)
	}

	switch {
	case *stats:
		printStats(root, rows)
	case *ledger != "":
		os.Exit(printLedger(rows, *ledger, *full))
	case *keyword != "":
		os.Exit(printSearch(rows, *keyword, *full))
	default:
		if flag.NArg() == 0 {
			flag.Usage()
			os.Exit(2)
		}
		os.Exit(printIDs(rows, flag.Args()))
	}
}

func claimsDir(override string) (string, error) {
	if override != "" {
		return override, nil
	}
	for _, c := range []string{"claims", filepath.Join("..", "claims"), filepath.Join("..", "..", "claims")} {
		if fi, err := os.Stat(c); err == nil && fi.IsDir() {
			return c, nil
		}
	}
	return "", fmt.Errorf("no claims/ directory found (run from the repo root, or pass -dir)")
}

// load parses every ledger into rows. A ledger row is a table line whose first
// cell contains a claim id; header and separator lines have none, so they drop
// out without a special case.
func load(root string) ([]row, error) {
	files, err := filepath.Glob(filepath.Join(root, "*.md"))
	if err != nil {
		return nil, err
	}
	if len(files) == 0 {
		return nil, fmt.Errorf("no ledgers under %s", root)
	}
	sort.Strings(files)

	var out []row
	for _, f := range files {
		fh, err := os.Open(f)
		if err != nil {
			return nil, err
		}
		sc := bufio.NewScanner(fh)
		// Rows carry their evidence, so a few run past bufio's default 64 KB.
		sc.Buffer(make([]byte, 0, 1<<16), 4<<20)
		for n := 1; sc.Scan(); n++ {
			ln := sc.Text()
			if !strings.HasPrefix(strings.TrimSpace(ln), "|") {
				continue
			}
			cells := splitRow(ln)
			if len(cells) == 0 {
				continue
			}
			m := idRe.FindStringSubmatch(cells[0])
			if m == nil {
				continue
			}
			out = append(out, row{id: m[1], file: filepath.Base(f), line: n, cells: cells, raw: ln})
		}
		fh.Close()
		if err := sc.Err(); err != nil {
			return nil, fmt.Errorf("%s: %w", f, err)
		}
	}
	return out, nil
}

// splitRow splits a markdown table line on unescaped pipes that are not inside a
// backtick span. Claim text is full of `disp:a9c4 | 32 hits` and of escaped
// pipes, and a naive strings.Split puts the rest of the row in the wrong column
// -- which is invisible until a Status reads as evidence.
func splitRow(ln string) []string {
	ln = strings.TrimSpace(ln)
	if !strings.HasPrefix(ln, "|") {
		return nil
	}
	var (
		cells []string
		cur   strings.Builder
		tick  bool
	)
	rs := []rune(ln)
	for i := 0; i < len(rs); i++ {
		switch {
		// Only `\|` and `\\` are escapes in a table cell. A lone backslash is data
		// -- claim text carries Windows paths (`terrain\tile*.bmp`), and consuming
		// it here would silently rewrite the claim on its way to the reader.
		case rs[i] == '\\' && i+1 < len(rs) && (rs[i+1] == '|' || rs[i+1] == '\\'):
			cur.WriteRune(rs[i+1])
			i++
		case rs[i] == '`':
			tick = !tick
			cur.WriteRune(rs[i])
		case rs[i] == '|' && !tick:
			cells = append(cells, strings.TrimSpace(cur.String()))
			cur.Reset()
		default:
			cur.WriteRune(rs[i])
		}
	}
	cells = append(cells, strings.TrimSpace(cur.String()))
	// A leading and trailing pipe give an empty first and last cell.
	if len(cells) > 0 && cells[0] == "" {
		cells = cells[1:]
	}
	if len(cells) > 0 && cells[len(cells)-1] == "" {
		cells = cells[:len(cells)-1]
	}
	return cells
}

func printIDs(rows []row, ids []string) int {
	byID := map[string][]row{}
	for _, r := range rows {
		byID[r.id] = append(byID[r.id], r)
	}
	missing := 0
	for i, want := range ids {
		want = strings.ToUpper(strings.Trim(want, "`"))
		hits := byID[want]
		if len(hits) == 0 {
			fmt.Fprintf(os.Stderr, "claim: %s is not in any ledger\n", want)
			missing++
			continue
		}
		if i > 0 {
			fmt.Println()
		}
		for _, r := range hits {
			// The retraction state is the reason to cite a claim id rather than an
			// experiment, so it is never something the caller has to ask for.
			mark := ""
			if r.file == "retracted.md" {
				mark = "  ** RETRACTED -- this row is the overturn, not the claim **"
			}
			fmt.Printf("%s  (claims/%s:%d)%s\n", r.id, r.file, r.line, mark)
			for i, c := range r.cells {
				// The id cell repeats the heading unless it carries a qualifier --
				// retracted.md scopes a row to "the writer-set clause only", and that
				// scope is the difference between a retraction and a whole claim.
				if i == 0 && strings.Trim(c, "`* ") == r.id {
					continue
				}
				fmt.Printf("  %s\n", wrap(c, 92, "  "))
			}
		}
	}
	if missing > 0 {
		return 1
	}
	return 0
}

func printLedger(rows []row, name string, full bool) int {
	name = strings.TrimSuffix(filepath.Base(name), ".md") + ".md"
	var hits []row
	for _, r := range rows {
		if r.file == name {
			hits = append(hits, r)
		}
	}
	if len(hits) == 0 {
		fmt.Fprintf(os.Stderr, "claim: no ledger %s\n", name)
		return 1
	}
	fmt.Printf("claims/%s — %d rows\n\n", name, len(hits))
	for _, r := range hits {
		if full {
			fmt.Println(r.raw)
			continue
		}
		fmt.Printf("%-22s %-10s %s\n", r.id, status(r), summary(r, 96))
	}
	return 0
}

func printSearch(rows []row, pattern string, full bool) int {
	re, err := regexp.Compile("(?i)" + pattern)
	if err != nil {
		fmt.Fprintln(os.Stderr, "claim:", err)
		return 2
	}
	n := 0
	for _, r := range rows {
		if !re.MatchString(r.raw) {
			continue
		}
		n++
		if full {
			fmt.Printf("%s  (claims/%s:%d)\n%s\n\n", r.id, r.file, r.line, r.raw)
			continue
		}
		fmt.Printf("%-22s %-14s %-10s %s\n", r.id, strings.TrimSuffix(r.file, ".md"), status(r), summary(r, 78))
	}
	if n == 0 {
		fmt.Fprintf(os.Stderr, "claim: no row matches %q\n", pattern)
		return 1
	}
	fmt.Fprintf(os.Stderr, "\n%d rows. Re-run with -full for the evidence, or name the ids.\n", n)
	return 0
}

func printStats(root string, rows []row) {
	type agg struct {
		rows  int
		bytes int64
	}
	per := map[string]*agg{}
	for _, r := range rows {
		a := per[r.file]
		if a == nil {
			a = &agg{}
			per[r.file] = a
		}
		a.rows++
	}
	names := make([]string, 0, len(per))
	for n := range per {
		names = append(names, n)
		if fi, err := os.Stat(filepath.Join(root, n)); err == nil {
			per[n].bytes = fi.Size()
		}
	}
	sort.Slice(names, func(i, j int) bool { return per[names[i]].bytes > per[names[j]].bytes })

	var tb int64
	var tr int
	fmt.Printf("%-18s %10s %7s %10s\n", "ledger", "bytes", "rows", "bytes/row")
	for _, n := range names {
		a := per[n]
		fmt.Printf("%-18s %10d %7d %10d\n", n, a.bytes, a.rows, a.bytes/int64(max(a.rows, 1)))
		tb += a.bytes
		tr += a.rows
	}
	fmt.Printf("%-18s %10d %7d\n", "TOTAL", tb, tr)
	fmt.Fprintf(os.Stderr, "\nOpening one ledger costs its whole byte count. One row is the size above.\n")
}

func status(r row) string {
	for _, c := range r.cells {
		switch {
		case strings.Contains(c, "✖"):
			return "retracted"
		case strings.Contains(c, "✔"):
			return "promoted"
		case strings.Contains(c, "●"):
			return "active"
		}
	}
	if r.file == "retracted.md" {
		return "retracted"
	}
	return ""
}

// summary is the claim cell, flattened and cut. It exists so an index costs a
// line per claim rather than a row per claim.
func summary(r row, n int) string {
	s := ""
	if len(r.cells) > 1 {
		s = r.cells[1]
	}
	s = strings.NewReplacer("**", "", "`", "", "\t", " ").Replace(s)
	s = strings.Join(strings.Fields(s), " ")
	if len(s) > n {
		s = s[:n-1] + "…"
	}
	return s
}

func wrap(s string, width int, indent string) string {
	words := strings.Fields(s)
	if len(words) == 0 {
		return ""
	}
	var b strings.Builder
	col := 0
	for i, w := range words {
		if col > 0 && col+1+len(w) > width {
			b.WriteString("\n" + indent)
			col = 0
		} else if i > 0 {
			b.WriteString(" ")
			col++
		}
		b.WriteString(w)
		col += len(w)
	}
	return b.String()
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
