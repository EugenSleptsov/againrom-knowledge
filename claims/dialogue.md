# DIALOGUE — the window a mission's script speaks through, and its lifecycle

Claims about the presentation side of `claims/mission.md` and
`claims/trigger.md`. Those say what the script decides and which resource its
decision names; these say what object draws it, everything that can enter that
path, what happens when the resource is absent, what ends the display and what
the window measures. The glyph rule, how a text byte becomes a pixel, and the
font atlases are EXP-0097's area; this ledger cites it rather than restating it.
Spec: [`formats/dialogue/format.md`](../formats/dialogue/format.md). Format of
this file: [registry.md](registry.md).

## Terms

- A root is the install a measurement was read from; every measurement names
  its root. `rom.exe` is byte-identical in the EN and RU roots
  (`sha256 942e9b72…d367d03`, 1 977 344 B), so every address is a fact about
  both, and only resources can differ by language.

## Window class, announcement path and lifecycle

| ID | Claim | Confidence | Status | Evidence |
|---|---|---|---|---|
| DLG-WIN-001 | One `0x84`-byte panel class with one constructor draws every line of authored dialogue, and its six calling routines are the complete set of entries into it. | High | ● active | [EXP-0098](../experiments/EXP-0098-mission-text/) |
| DLG-PATH-002 | Script and outcome announcements share one transport; the failure close is selected by its own stored panel pointer. | High | ● active (partially retracted, amended) | [EXP-0098](../experiments/EXP-0098-mission-text/), [EXP-0101](../experiments/EXP-0101-mission-end/), [EXP-0274](../experiments/EXP-0274-defeat-modes/) |
| DLG-ABSENT-003 | A named event text that does not ship produces nothing: no window, no fallback, no fault, no blocked state; this answers `MISSION-TEXT-005`'s Unknown. | High | ● active | [EXP-0098](../experiments/EXP-0098-mission-text/) |
| DLG-EMPTY-004 | The window's `"Nothing to say"` fallback is reached only by an existing file that yields no part 1, and no shipped file on either root reaches it. | High | ● active | [EXP-0098](../experiments/EXP-0098-mission-text/) |
| DLG-LIFE-005 | Only input ends the display, nothing queues, and an announcement that arrives while any dialog is open is discarded. | High / Medium | ● active (amended) | [EXP-0098](../experiments/EXP-0098-mission-text/), [EXP-0108](../experiments/EXP-0108-panel-modality/) |
| DLG-READ-006 | The event text is read at fire time, not at map load, and it is read twice per announcement. | High | ● active | [EXP-0098](../experiments/EXP-0098-mission-text/) |

### DLG-WIN-001

- `FUN_004217be(name)` allocates `0x84` bytes (`0042180a PUSH 0x84`) and
  constructs the panel with `FUN_004c5682(id=9, 30, 120, 610, 360)`. Rects are
  `{left,top,right,bottom}`, fixed by `FUN_00447d20` being
  `MOV EAX,[EAX+0xc] / SUB EAX,[ECX+0x4]`, `bottom − top`.
- It then adds:
  - a portrait (`FUN_004c4995`, id 12, 30,54–118,168), only when `panel+0x7c`;
  - a text control (`FUN_004be10e`, id 10, at 128,36–428,172 with the portrait
    and 48,36–428,172 without), whose initial content is the literal
    `"Nothing to say"` at `0x005b89c0`/`0x005b89d0`;
  - always a button (`FUN_004be6e3`, id 11, 200,172–280,198) whose label is
    string-table entry 77 and whose command is `0x46f`.
- Then `FUN_00476810` shows it.
- `EnumRefs callto:4c5682` returns 1 hit, 1 owner, 0 in orphan code:
  `FUN_004217be` is the class's only constructor.
- `EnumRefs callto:4217be` returns 7 hits, 6 owners, 0 in orphan: the mission
  event text (`FUN_00473110`), the inn's NPCs (`FUN_00480fd0`, ×2), the
  mercenary hall (`FUN_0047e410`, `FUN_004b4c70`), the shop keeper
  (`FUN_004a8bc3`) and the training hall (`FUN_004b7f20`).
- Each of the six formats a name in its own family, and all six families ship
  (`evidence/family-en.txt`: 318/22/14/4/33 nodes).
- The child rects are panel-relative, and the drawn layout settles it. Read as
  screen coordinates, the portrait's top edge (54) would sit 66 px above the
  panel's own (120). Read as panel-relative, all four children fall inside the
  580×240 panel with nothing overhanging: 118 < 580, 428 < 580, 172 < 240,
  198 < 240. Four independent literals would have to coincide for that.

**Confidence.** High for the geometry, the two layouts, the button command and
the string-table index: each is a named `PUSH` in one routine read whole, and an
instruction, not a fit, settles the rect convention. High for the entry
enumeration with its instrument stated: `callto:` on the repaired table, which
also scans `.rdata` dwords holding the address, 0 orphan on both sweeps. Its
blind spot is a construction reached only through a vtable slot; the class
census `EnumRefs range:4c52b0:4c58a0` (13 functions, 1468 function bytes,
0 orphan bytes, 0 orphan runs) closes that hole for this class.

### DLG-PATH-002

- `004ea1b3` stamps static packet `00609c38` with `0xb6`; `004ea2f6` stamps the
  supplied outcome opcode; both use `004e74fe`.
- Client `004104e8` selects `opcode-3` through index `0041879b[188]` and table
  `004186e3[46]`: `0xb6 -> 0x433` with its number, `0xb4 -> 0x433/0xff`,
  `0xb5 -> 0x430`.
- The failure sentinel routes to `0x431`: constructor `00446de2`, panel stored
  at `frontend+0x110`, title 141. The vtable slot `00598bc0` contains
  `00447063`, not `00446d35`.
- The base close stores the result and posts `0x44c`. `004757b0` matches
  `+0x110` and maps `0x445 -> 0x41e` (teardown/menu), otherwise `0x418` (load
  selection).
- The separate win programme may post `0x41d` under its fast-transition
  predicate or build its own panel.

**Confidence.** High for the named instruction paths and the discriminating
branch and vtable counterexamples. No original runtime timing or GUI
observation is claimed.

**Amended.** `MISSION-DEFEAT-046` (EXP-0274) refutes the old dialogue-sentinel
close attribution and the later added neighbouring-class `0x41d` hop.
`retracted.md` records the withdrawn text: the failure close tested
`frontend+0x414 == 0xff`, and its class handler `00446d35` posted `0x41d`
upstream. The actual failure close matches stored panel pointer `+0x110` and
result `panel+0x60`, and its handler is `00447063`. The shared transport
stands.

### DLG-ABSENT-003

- In `FUN_00473110`'s `0x433` arm the path is built, the reader is constructed
  on the stack, and `0047396c CALL 0x004c9f10` is made with both optional
  arguments zero (`0047395d PUSH EDI / PUSH EDI`, `EDI = 0`): flags 0 and no
  error-out pointer, so the routine's own error code 2 is discarded.
- `FUN_004c9f10` returns 1 only when `FUN_004c9b10(name)` resolves and 0
  otherwise. `00473994 CMP ESI,EDI` / `00473996 JZ 0x004739b1` then skips the
  call to the window builder, and the arm returns 1 as if it had worked.
- The only durable effect is `campaign+0x414`, written before the test at
  `004738c1`.
- Corpus, EN root: of 242 numbers the 28 campaign scripts raise, 223 name a
  shipped file and 19 are silent. RU raises the identical 242 (0 differing
  maps) and silences 16 (`evidence/join.txt`).
- This reproduces `MISSION-TEXT-005`'s 223/242 from two independent
  instruments, the raised set from `scenario.res` and the shipped set from
  `main.res`, and makes the figure a behavioural statement rather than a
  discrepancy.

**Confidence.** High. The guard, its two zero arguments and the skip are named
instructions in one routine read whole. The alternative readings, a fallback
resource, an empty window or a logged failure, each predict an instruction that
is not there, and the routine has no other exit between the open and the
builder.

### DLG-EMPTY-004

- The window's "on show" slot `vt+0x80` is `FUN_004c584a`, which calls the
  pager `FUN_004c5886` once and discards its return.
- The pager returns 0 when `FUN_004c607a` cannot find a tag containing
  `part=<n>`. On that path nothing sets the text control, so the window opens
  on the constructor's own literal `"Nothing to say"` and its button then
  closes it.
- Corpus: 0 of 225 EN and 0 of 228 RU event files lack a `part=1` tag, and 0
  files on either root have a gap in their part sequence before their own
  maximum (`evidence/markup-en.txt`, `markup-ru.txt`).
- The string is unreachable from shipped mission data, so a consumer that omits
  it diverges on nothing the original shows. It is reachable from an authored
  file, which is why it is published.

**Confidence.** High for the mechanism: the ignored return and the two literals
are named instructions, and the pager's zero return is its own `if`. High for
the corpus half: exhaustive over every event file of both roots, and the
measurement could have failed, since one file with a part-1 typo would show it.

### DLG-LIFE-005

- Three inputs reach the same command: the button (`0x46f`), and, through
  `vt+0x6c` `FUN_004c574d`, the keys `0x0d` (RETURN) and `0x1b` (ESCAPE), both
  of which self-send `0x46f`.
- `FUN_004c5797` answers it by calling the pager. Non-zero: redraw children 10
  and 12 and stay open. Zero: post `0x445`, which `FUN_004c52f3` turns into
  `vt+0x84` (free the text buffer) plus `0x44c`, whose window arm is
  `0047333c CALL 0x004757b0`.
- There is no timer: the class census over `4c52b0..4c58a0` finds 13
  functions, 0 orphan bytes, and none references `0x113`.
- One instruction settles overlap. `FUN_00476810` sets bit 3 of
  `campaign+0x3dc` when it shows a panel, and the `0x433` arm tests
  `004738df TEST byte ptr [EBP+0x3dc],0x8` and returns without building
  anything when it is set.
- Announcements are therefore dropped, not deferred: a consumer that queues
  them shows text the original never showed. The drop at `004738df` does not
  depend on where the clears of the bit live.

**Confidence.** High for the three inputs, the pager loop, the absence of a
timer and the drop: each is a named instruction. Medium for "no other route
closes it": `FUN_004c52f3`'s default arm `FUN_004bd9dc` and the base key
handler's siblings were not read line by line.

**Amended.** EXP-0108 refutes the clear enumeration in both halves, and
`retracted.md` withdraws it. The withdrawn sentence: "The bit is cleared by
nine instructions, all nine inside `FUN_004757b0`" (`EnumRefs imm:fffffff7`:
19 hits / 10 owners / 0 in orphan), with the grade "the gate's set and its nine
clears are enumerated with the instrument and its orphan count stated".
`imm:fffffff7` is a dword-immediate sweep and is blind by construction to the
byte-width form of the same operation. It missed `00475c9c AND AL,0xf7` and
`004761f6 AND AL,0xf7` inside `FUN_004757b0`, so the count is at least twelve,
and `00474a7d AND AL,0xf7` inside `FUN_00473110`, between
`00474a77 MOV EAX,[EBP+0x3dc]` and `00474a7f MOV [EBP+0x3dc],EAX`, so not all
clears are in `FUN_004757b0`. That arm also sets the bit itself
(`00474a59 OR EDX,0x8`) without going through `FUN_00476810`: a shutdown
bracket that wipes the screen through `FUN_0044f990` at level 0, calls
`FUN_0047a600` and posts `WM_CLOSE`. The miss is what `AGENTS.md`'s
enumeration rule 4 guards against: reduce by the field's own width and read
every byte- and word-wide hit. The three inputs, the absent timer, the drop at
`004738df` and "dropped, not deferred" stand.

### DLG-READ-006

- `FUN_00473110` opens the resource once only to decide whether to proceed
  (`DLG-ABSENT-003`) and discards the payload.
- `FUN_004217be` then passes the relative name (`004739a1 MOV EAX,[ESP+0x14]`,
  the string before the `"main\text\"` and `".txt"` concatenations) to the
  constructor, which stores it at `panel+0x6c`.
- `FUN_004c5f06` rebuilds the same path from it
  (`004c5f43 PUSH 0x5c1c44 "main\text\"`, `004c5f37 PUSH 0x5c1c3c ".txt"`),
  opens it again, and copies the whole payload into `panel+0x74`
  (`FUN_004ca130` size → `FUN_00554390` alloc → `FUN_004ca140` read → NUL).
- Nothing loads a mission's event files at map load: the only openers of this
  name family are those two, and both are inside the fire-time path.
- `vt+0x84` frees `panel+0x74` on close, so the text is resident only while the
  window is.

**Confidence.** High. Both opens are named instructions in two routines read
whole, and the relative-name hand-off is the instruction that makes the second
one possible. The "at load" rival predicts an opener in the map-load path, and
`FUN_00473110`'s arm is reached only by a posted `0x433`.

## Content model, layout and language roots

| ID | Claim | Confidence | Status | Evidence |
|---|---|---|---|---|
| DLG-MARKUP-007 | The content model is a tag scan, not a grammar; its vocabulary is fifteen lowercase literals, and five of them are dead in the shipped corpus. | High / Medium | ● active | [EXP-0098](../experiments/EXP-0098-mission-text/) |
| DLG-FACE-008 | The whole file decides once whether the window has a portrait, each part decides which portrait it shows, and the two tests differ. | High / Unknown | ● active | [EXP-0098](../experiments/EXP-0098-mission-text/) |
| DLG-WRAP-009 | The text is wrapped and measured against its own font, and this window clamps the line count rather than scrolling. | High / Medium | ● active | [EXP-0098](../experiments/EXP-0098-mission-text/) |
| DLG-LANG-010 | Only resources depend on the language root, and the roots differ in three mission events: EN silences three announcements that RU speaks. | High / Medium | ● active | [EXP-0098](../experiments/EXP-0098-mission-text/) |
| DLG-READER-011 | Before its EXP-0101 repair, this repository's shared `.res` reader rejected the RU `MAIN.RES`, so every tool using it was blind to the RU half of `main.res`. | High | ● active (amended) | [EXP-0098](../experiments/EXP-0098-mission-text/), [EXP-0101](../experiments/EXP-0101-mission-end/) |

### DLG-MARKUP-007

- `FUN_004c607a(part, …)` walks the payload for `<` … `>`, lowercases each tag
  body (`FUN_00572f3e` → `FUN_00557bf0`) and substring-searches it.
- The literals, in the order the routine tests them: `part=%d`, `npc=`,
  `iamfemale`, `iammale`, `iammage`, `iamfighter`, `npcalive=`, `npcdead=`,
  `female`, `male`, `mage`, `fighter`, `sound=`, `tune=`, `tips=`.
- The part match is `CString::Find("part=%d")`, a substring test: a tag reading
  `part=10` also satisfies a search for part 1, and the first matching tag in
  file order wins.
- Corpus, both roots: 225 EN / 228 RU event files, 513 / 518 parts, max part
  13, and 0 parts on either root whose number is a prefix of an earlier part
  number in the same file. The hazard exists and the shipped data never trips
  it.
- Never used on either root: `npcalive=`, `npcdead=`, `sound=`, `tune=`.
  `iamfighter` is used once on RU and never on EN (`evidence/markup-*.txt`).

**Confidence.** High for the vocabulary and the scan: one routine's own string
operands, read with `StrDump fnstr:`, and the lowercasing is a named call.
Medium for what each conditional literal does: only `part=`, `npc=` and
`tips=` were followed to their consumers, and the four `iam*` arms and the four
npc-flag arms were read as tests, not as effects.

### DLG-FACE-008

- The constructor's `FUN_004c5f06` lowercases the entire payload and sets
  `panel+0x7c = (Find("npc") >= 0)` (`004c6014 PUSH 0x5c1c50 "npc"`).
  `FUN_004217be` reads that flag to choose between the two layouts.
- Inside `FUN_004c607a` the same field is rewritten from the current part's own
  tag, `panel+0x7c = (Find("npc=") >= 0)`, and `FUN_004c5886` uses it to decide
  whether to refresh child 12.
- A file in which only one part names a speaker therefore gets the portrait
  pane for all of its parts, and the parts without `npc=` leave the previous
  face standing.
- The first test has no `=`: prose containing the letters `npc` would open the
  portrait layout on a file with no speaker at all.
- Corpus: on both roots the two tests agree on every file, 225/225 EN and
  228/228 RU, so the shipped corpus cannot discriminate them, and the
  separation is established from the instructions alone.

**Confidence.** High for the two tests and their different consumers: three
named instructions in two routines.

**Unknown.** Whether the disagreement is ever reachable in shipped data. By the
census above it is not, and no shipped file exercises the no-portrait layout
through this entry point.

### DLG-WRAP-009

- `FUN_004be229(ctrl, text)` calls `FUN_00456ab0(ctrl+8, text)`, which splits
  the string and calls `FUN_004563e0(rect, line)` per piece into a single
  static accumulator at `DAT_005e9ba8` (reset by `FUN_0056fc7a(0,-1)`), then
  copies it into the control's own line array at `ctrl+0x64`.
- It then computes
  `ctrl+0x8c = min( Height(ctrl+8) / ctrl+0x60 , GetSize(ctrl+0x64) )`
  (`004be25b`…`004be2a3`), with `FUN_00447d20` = `bottom − top` and
  `FUN_004478f0` = `[this+8]`, the array size.
- The line pitch `ctrl+0x60` defaults to font height + 2 (`FUN_004be10e`:
  `param_10 == 0 → param_1[0x17] + 2`). `FUN_004217be` passes 0, so the
  dialogue window always takes the default.
- The control can scroll: `vt+0x80` `FUN_004be2b0` clamps a top line into
  `[0, size − visible]`, redraws, drives a scrollbar child and notifies with
  command `0x46d`. `FUN_004217be` gives it no scrollbar child, and no arm of
  the panel class answers `0x46d`. The base key handler `FUN_004c5369` binds
  Tab and the four arrows to focus movement, not scrolling.
- The rect is 300×136 px with the portrait and 380×136 without. The longest
  shipped part body is 230 characters (EN) / 211 (RU).
- The per-glyph advance that turns 300 px into a character count is
  EXP-0097's, and no figure here depends on it.

**Confidence.** High for the wrap, the clamp formula and the default pitch:
four named instructions plus two one-line accessors read whole. Medium for
"text past the clamp is unreachable in this window": `FUN_004bd9dc`, the
panel's default command arm, was not read, and whether `panel+0x38`'s
navigation links are ever filled was not established.

### DLG-LANG-010

- `rom.exe` is byte-identical across both preserved roots and the GOG install
  (`sha256 942e9b72…`, 1 977 344 B), so every geometry, index and format string
  in this ledger is shared.
- The two language-carrying inputs are
  `main.res::text/battle/m<n>/event<NN>.txt` and the indexed table
  `main.res::text/main.txt`.
- The compiled indices survive the change of root: 274 lines on both, with
  entry 77 (the button), 140 (the win panel) and 141 (the lose panel) present
  on both. The EN lines measure 2, 17 and 14 bytes and are ASCII; the RU lines
  7, 16 and 16 bytes and are non-ASCII (`evidence/strtab-*.txt`).
- The event corpora differ: EN 225 files, RU 228. The three RU-only files are
  m100/event09, m130/event07 and m150/event10. All three are raised by scripts
  that are identical on both roots (28/28 maps, 0 differing raised sets), so
  those three announcements are silent in English and spoken in Russian.
- `speech.res` carries no `battle` subtree on either root (0 of 355 EN / 349 RU
  nodes), so the `speech\…\….wav` the pager composes never resolves for a
  mission event.

**Confidence.** High for the identical image and the three-file delta: a byte
hash, and two exhaustive node censuses whose per-mission diff is three lines.
High for the string table's index stability: the indices are compiled
constants and both roots have the same line count, which a compiled index
requires and a differing count would have refuted. Medium for the speech
clause: the `%s` of `"npc%02de%sp%d"` was not pinned, so the absence rests on
there being no `battle` node at all rather than on a composed name missing.

### DLG-READER-011

- As it stood before the EXP-0101 commit, `internal/rom.OpenArchive` required
  `(len(file) − registryOffset) % 32 == 0` (`internal/rom/res.go:38`).
- The RU root's `MAIN.RES` is 4 922 540 bytes with registry offset 4 906 389.
  The remainder is 504 records of 32 bytes plus 23 trailing bytes, so the
  reader rejected it outright with `bad registry offset 4906389`.
- The registry itself is well formed and parses cleanly: 504 nodes against
  EN's 494, the count `RES-*`'s own EN↔RU compare already published.
- This is a fact about the instrument, not about the game (`METHODOLOGY.md` →
  *facts about the instrument*). It is recorded because it silently removes a
  whole root from any measurement: `tools/campaign -mode msgs` returned
  `open: …\main.res: bad registry offset` on RU
  (`evidence/raised-vs-shipped-ru.txt`), and any lane that measures "the RU
  text" through `OpenArchive` gets an error it may read as absence.
- EXP-0098 reproduced the parse locally (`tools/misstext` `openRes`) rather
  than patch a file its sibling lanes share.

**Confidence.** High. The constant is one line of committed source, the
arithmetic is exact, and the tolerant re-parse yields a node count that agrees
with an independently published EN↔RU figure. `internal/rom/res_test.go`
verifies the repair; its synthetic cases (the repository's own bytes) pin the
count-from-`@0x14` law and the residue tolerance without the install.

**Amended.** EXP-0101 repaired the reader on 2026-08-02.
`internal/rom.OpenArchive` now implements `RES-ACCEPT-031` directly: magic,
registry at `@0x10`, exactly `@0x14` records of 32 bytes, trailing bytes
ignored (`Archive.Residue` reports them). It opens 9/9 EN and 8/8 RU archives,
RU `MAIN.RES` as 504 nodes / 462 files, which is `RES-GEOM-028`'s own figure.
The measurement above describes the instrument before that commit; it is no
longer true of the tree, and `tools/campaign -mode msgs` on RU now returns
output rather than `open: … bad registry offset`.

## Modality while a panel is up

| ID | Claim | Confidence | Status | Evidence |
|---|---|---|---|---|
| DLG-STOP-012 | The world stops while a dialogue panel is displayed, and one instruction in the whole image decides it. | High | ● active | [EXP-0108](../experiments/EXP-0108-panel-modality/) |
| DLG-DIM-013 | When a panel is shown the whole screen behind it is darkened once, destructively, to 13/16 brightness, through the shroud table `TERR-FOG-084` pinned. | High | ● active | [EXP-0108](../experiments/EXP-0108-panel-modality/) |
| DLG-CLOCK-014 | The time a dialogue panel stops is discarded, not caught up, and the close's guard names the only state the engine expects the panel to interrupt. | High | ● active | [EXP-0108](../experiments/EXP-0108-panel-modality/) |
| DLG-DRAW-015 | Nothing else is drawn differently while a panel is up: six instructions read the panel bit, and none of them is in a draw path. | High / Medium | ● active | [EXP-0108](../experiments/EXP-0108-panel-modality/) |
| DLG-ENTRY-016 | All six entry points and both mission-outcome panels take the identical mechanism; they differ in the state each runs from, plus one close arm. | High / Medium | ● active | [EXP-0108](../experiments/EXP-0108-panel-modality/) |
| DLG-MODAL-017 | A second, stronger modality, a nested message loop, exists in the image; no dialogue entry uses it, and it makes harmless the one panel the halt gate cannot see. | High | ● active | [EXP-0108](../experiments/EXP-0108-panel-modality/) |

### DLG-STOP-012

- `FUN_00476810` is `DLG-WIN-001`'s "then `FUN_00476810` shows it" and the only
  routine that puts a panel on the session window (`callto:476810` = 19 hits /
  6 owners / 0 in orphan). It sets bit 3 of `campaign+0x3dc` for every panel
  it is handed (`0047682e OR AL,0x8`, `00476837 MOV [EBX+0x3dc],EAX`).
- The MFC idle handler loads that word at `0047165f` and, once the phase is the
  campaign's, tests it: `0047167a CMP EDX,0x2` / `0047167f TEST EAX,0x4008` /
  `00471684 JNZ 0x004716d8`. `004716d8` sets the keep-idling flag and jumps to
  the handler's common tail, running no pacer arm at all.
- `EnumRefs imm:4008` is 1 hit / 1 owner / 0 in orphan: that test exists once
  in the image.
- Against every other branch of the same routine, the halted branch omits:
  - the sub-tick `FUN_004d2551` (`callto:4d2551` = 4 hits / 4 owners / 0 in
    orphan: `FUN_00475280`, `FUN_004753c0`, `FUN_00475550`, `FUN_00477c00`;
    the branch calls none of them);
  - `ANIM-CLOCK-001`'s `0x401` presentation tick;
  - the `0x402` post;
  - `FUN_004104e8` on the map view;
  - `FUN_004d88f1`'s command drain (`SESS-PAUSE-020`), which makes the halt
    stricter than the pause.
- All five branches then converge on `00471729` and do identical work there,
  so nothing in the tail can be what a panel changes.
- The rival that would make this cosmetic is a timer, and the import table
  refutes it: `SetTimer`, `KillTimer`, `timeSetEvent` and `timeKillEvent` are
  absent from rom.exe's entire import directory (`SESS-TIMER-022`). No
  callback exists that could advance `server+0x04` behind the panel; its only
  incrementer is inside `FUN_004d891a` (`SESS-TICK-004`; `callto:4d891a` =
  2 hits, both in the server family).
- This evidence also corrects `DLG-LIFE-005` (filed in `retracted.md`). That
  claim's "the bit is cleared by nine instructions, all nine inside
  `FUN_004757b0`" is wrong in both halves, because its instrument
  `imm:fffffff7` is blind to the byte-width form: `004761f6 AND AL,0xf7`, the
  clear the dialogue panel itself takes, is one it missed, and `00474a7d` is in
  `FUN_00473110`.
- Bit 3 is not exclusively `FUN_00476810`'s to set: `00474a59 OR EDX,0x8` sets
  it directly in a shutdown bracket that darkens through `FUN_0044f990` at
  level 0 and posts `WM_CLOSE`. Neither correction touches this claim, whose
  subject is the gate rather than the writer set.

**Confidence.** High. Every step is a named instruction in one of three
routines read whole. The gate is an image-wide enumeration returning 1 hit with
0 orphan, and the tick's issuer set another returning 4 with 0. An import-table
absence, the instrument this repository treats as immune, closes the timer
rival.

### DLG-DIM-013

- Inside `FUN_00476810`, between a DirectDraw Lock (`FUN_0044c3b0` →
  `vt+0x64`, with `0044c3cd MOV dword ptr [0x005e4398],0x6c`, the
  `DDSURFACEDESC` size) and Unlock (`FUN_0044c410` → `vt+0x80`), the routine
  calls the remap with the whole surface: `00476861 PUSH 0x3` /
  `PUSH [0x005ea20c]` / `PUSH [0x005ea208]` / `PUSH 0x0` / `PUSH 0x0` /
  `00476869 CALL 0x0044fad0`.
- `0x005ea200` is a RECT (`FUN_0044cf00` hands it to `CopyRect`, IAT
  `0x00632f8c`), so those two globals are its right and bottom.
- `FUN_0044fad0(x0,y0,x1,y1,L)` clips x to `[0x005e4408]`/`[0x005e4410]` and y
  to `[0x005e440c]`/`[0x005e4414]`, and per pixel does
  `dst = LUT[L·[0x005e42f0] + (dst>>3 or dst)]` through `[0x005e8420]`
  (`0044fad7`, `0044fadd`, `0044fae6`). The
  `0044fb4d CMP dword ptr [0x005eb570],0x0` split chooses the 13- or 16-bit
  index. This is the table, stride global and pair of index forms
  `TERR-FOG-037` describes for the shroud blitters.
- `TERR-FOG-084` fixes its law as `out = (in × (16 − L)) >> 4`, so L = 3 is a
  per-channel gain of 13/16 = 0.8125.
- Re-executed over the whole 16-bit pixel space
  (`tools/panelmodal -mode shade`): 65 535 of 65 536 pixels get strictly
  darker, 1 (black) is unchanged, and none is brightened. Mean channel level
  0.5000 → 0.3937 in RGB565 and → 0.3911 in RGB555.
- It is not a render state and not per-frame: one call, straight into the
  locked framebuffer, before the panel's own `vt+0x34` draw. It survives only
  because `DLG-STOP-012`'s gate stops everything that would repaint.
- `callto:44fad0` = 11 hits / 11 owners / 0 in orphan. The other ten are in
  the `0x43xxxx`, `0x459xxx`, `0x4abxxx` and `0x4bexxx`–`0x4c3xxx` control
  families. Five also push level 3 (`004beea2`, `004c0a48`, `004c0ffb`,
  `004c12cd`, `004c3e25`), one pushes 10 (`004c1e72`), one 8 (`00459a52`),
  and three pass the level in a register (`004aba3f`, `004bf73d`, `00431779`).
- `FUN_00476810` is the only one of the eleven that passes the screen rect:
  `refto:5ea208` (65 hits / 47 owners / 0 in orphan) contains `0047685b` and
  none of the other ten owners. None of `FUN_00476810`'s five sibling show
  routines calls it at all.

**Confidence.** High for the call, its five arguments, the Lock/Unlock pair and
the once-ness: each is a named instruction in two routines read whole, and the
whole-screen extent is a property of the argument list, not an observation that
nothing smaller was seen. The numeric result is no stronger than
`TERR-FOG-084`, which is High: the probe re-derives nothing and applies that
law. It first reproduces that experiment's own two anchors (L = 16 all-zero,
L = 8 `== (px>>1) & mask`, 65536/65536 in both layouts), so the law applied is
visibly the one pinned.

### DLG-CLOCK-014

- The dialogue panel is stored in no `campaign+…` slot: `FUN_004217be`
  allocates it, hands it to `00421a3f CALL 0x00476810`, and stores it nowhere.
  Its close therefore runs `FUN_004757b0`'s default arm, after all 25
  stored-panel comparisons fail.
- The default arm: `004761ec MOV EAX,[EBP+0x3dc]` / `004761f2 TEST AL,0x8` /
  `004761f6 AND AL,0xf7` / `004761f8 CMP EAX,EBX` (EBX = 1, from `004757ce`) /
  `004761fa MOV [EBP+0x3dc],EAX` / `00476200 JNZ` /
  `00476202 CMP dword ptr [EBP+0x6bc],0x2` / `00476209 JNZ` /
  `0047620b CALL dword ptr [0x00632fa4]` (`timeGetTime`) →
  `00476211 MOV [EBP+0x3fc],EAX`, `00476219 MOV [EBP+0x3e4],0`,
  `0047621f MOV [EBP+0x400],0`.
- `campaign+0x3e4` is the pacer phase. `FUN_004753c0` rebases its deadline base
  `campaign+0x3ec` to `timeGetTime()` only at phase 0 (`004753f3 TEST EDI,EDI` /
  `004753f7 CALL EBX` / `004753f9 MOV [ESI+0x3ec],EAX`).
- Zeroing the phase makes the first iteration after the close rebase and issue
  one sub-tick. An untouched phase would run the catch-up `00475461 JBE` until
  the `AND EDX,0xf` wrap brought it back to 0: up to fifteen sub-ticks fired
  back to back (`SESS-CLOCK-021`).
- The `CMP EAX,EBX` makes this discriminating rather than incidental: the
  restart happens only when clearing bit 3 leaves the word at exactly 1,
  `SESS-SCREEN-003`'s map-session-only state, and only in phase 2. A panel
  closed over the town (`campaign+0x3dc == 0`, `SHOP-TOWN-023`) takes the same
  arm, clears the same bit and restarts nothing.
- The close tail then redraws the frame window
  (`0047627d CALL dword ptr [EDX+0x34]` on `campaign+0xcc`), which repaints
  `DLG-DIM-013`'s darkened pixels. It is skipped only when the word is 1 and
  `[mapview+0x80] == 0`; then the resumed pacer's next frame repaints them.

**Confidence.** High. One routine read whole, every step a named instruction.
The raw bytes at `004761ec`, `8b 85 dc 03 00 00 a8 08 74 3b 24 f7 3b c3 …`,
carry the `JZ rel8 = 0x3b` that fixes the no-op exit at `00476231`
independently of any disassembler.

### DLG-DRAW-015

- No lighting change, no palette change and no per-frame dimming.
- `EnumRefs disp:3dc`, whole image: 242 hits / 80 owners / 0 in orphan or
  undisassembled code. Every hit that could involve bit 3 was followed to the
  instruction that tests it, a step a `disp:` sweep cannot take.
- The one register-carried mask, `00473042 TEST byte ptr [ESI + 0x3dc],DL`, is
  bit 0: `00472fba MOV EDX,0x1` at the head of `FUN_00472fb0`, never reassigned
  before it. The four byte-wide loads `00490dbf`, `0049148c`, `004b08eb`,
  `004b0c47` test `0x3`, `0x1`, `0x2`, `0x2` respectively.
- That leaves six readers of bit 3:
  - `0047167f` (`TEST EAX,0x4008`, the halt gate; `imm:4008` = 1 hit
    image-wide);
  - `004738df` (the announcement drop, `DLG-LIFE-005`);
  - `004761f2` (the close arm, `DLG-CLOCK-014`);
  - three early-outs: `004902b1 TEST byte ptr [EAX + 0x3dc],0xa` in
    `FUN_00490280`, `00490c9d` and `00490fc5` in `FUN_00490c60`, and
    `004b9c04` in `FUN_004b9bc0`. Each returns or jumps past its body when the
    bit is set.
- None of the six lies in the terrain/sprite/shroud blit family
  (`0x0044dxxx`–`0x00452xxx`) or the map-view draw family
  (`0x00404xxx`–`0x0040cxxx`) that `TERR-*` establishes.
- The two composite masks that read the word inside the `0x0049xxxx` family do
  not contain bit 3: `004915cd TEST …,0x226` is bits 1, 2, 5, 9 and
  `0049195b TEST …,0x627` is bits 0, 1, 2, 5, 9, 10.
- A consumer needs no second render path for "a panel is up", only
  `DLG-DIM-013`'s one-shot remap.

**Confidence.** High for the reader enumeration with its instrument and orphan
count stated, for the register-carried mask being bit 0, and for the two mask
arithmetics. Its blind spot: a wholesale `REP MOVSD`/`memcpy` of the campaign
object carries no displacement, and no `disp:` sweep can see it. Medium that
the three early-outs are input paths rather than draw paths: each was read at
its guard and the two instructions after it, not end to end. `FUN_00490280`
and `FUN_00490c60` read the cursor globals `[0x005cd7a0]`/`[0x005cd7a4]` that
the pacer's own edge-scroll block reads at `0047547d`…`004754fc`, and
`FUN_004b9bc0`'s not-taken path calls the panel command router `FUN_004c52f3`.

### DLG-ENTRY-016

- `callto:4217be` = 7 hits / 6 owners / 0 in orphan, reproducing
  `DLG-WIN-001`'s six. All seven reach the same `00421a3f CALL 0x00476810`
  inside the constructor, so the bit, the darkening and the gate are the same
  at every entry, and whether the entry points differ reduces to whether the
  state differs.
- The gate needs three conditions at once:
  - `campaign+0x3dc & 1`, a map session (`00471665 TEST AL,0x1` diverts first);
  - `campaign+0x6bc == 2`, the campaign;
  - `campaign+0x6b8 != 0`, this process owns the simulation (`00471678`
    diverts a pure network client to `FUN_00475610` before the mask is ever
    tested).
- The mission-event entry, `FUN_00473110`'s `0x433` arm, is the one that runs
  with a map session on screen.
- The inn, the two mercenary-hall entries, the shop keeper and the training
  hall run from the town, which `SHOP-TOWN-023` fixes at
  `campaign+0x3dc == 0`. There `FUN_00476810` still darkens the screen and sets
  bit 3, and the halt and the clock restart are no-ops: a mechanism that
  executes and changes nothing observable.
- `DLG-PATH-002`'s two outcome panels use the same routine:
  `0047411a CALL 0x00476810` for the win panel (string-table entry 140, stored
  at `campaign+0x114`) and `004741a7` for the lose panel (entry 141,
  `campaign+0x110`). Both darken and both halt.
- They diverge only on close. `FUN_004757b0` dispatches on 25 stored panel
  pointers. `+0x110` has its own arm (`0047615b`: clears bit 3, posts `0x41e`
  or `0x418`, and does not restart the clock). `+0x114` appears in none of the
  25 and falls through to the default arm with the dialogue panel.

**Confidence.** High for the shared path, the two outcome arms and the 25-slot
list. `FUN_004757b0` was read whole and every comparison form in it collected,
not only the common `CMP ESI,dword ptr [EBP+disp]`: two of the twenty-five are
a plain `CMP ESI,EAX` after a separate load, and reading only the common form
is how a slot list loses a member. Medium for the five town entries: this claim
fixes the state the gate requires and takes the town's state from
`SHOP-TOWN-023`. The five callers' own reachability was not read, so "no town
dialogue can ever run with bit 0 set" is not claimed.

### DLG-MODAL-017

- `FUN_0047c8f0(panel)` calls `FUN_00476810` and then runs a nested
  `GetMessageA` loop: `0047c903 MOV ESI,dword ptr [0x00632f34]`
  (`GetMessageA`), `0047c914 CALL ESI`,
  `0047c926 CMP dword ptr [ESP + 0x14],0x44c`, `TranslateMessage`
  `[0x00632f64]` and `DispatchMessageA` `[0x00632f68]`. It loops until the
  close message `0x44c` arrives, re-posts it and returns.
- The loop never returns to `CWinThread`'s idle processing, so `FUN_00471600`
  does not run at all: nothing is paced, gated or repainted, whatever
  `campaign+0x3dc` and `+0x6bc` hold.
- `callto:47c8f0` = 18 hits / 11 owners / 0 in orphan, and none of
  `DLG-WIN-001`'s six entry points is among the owners. Four of the eighteen
  are inside `FUN_00473110` itself (`00473cda`, `00473dbb`, `004746e9`,
  `00474777`), in arms other than the `0x433` one. That routine therefore
  appears in both call lists with a different meaning in each, and the two
  lists must be separated by call site rather than by owner.
- Two of the eleven owners, `FUN_0043aa85` (`0043aca1`) and `FUN_0043ba59`
  (`0043bd73`), store their panel at `campaign+0x3a8` immediately before the
  call. That is the single pointer `FUN_00476810` special-cases:
  `00476818 MOV EAX,dword ptr [EBX + 0x3a8]` / `0047681f CMP ESI,EAX` /
  `00476829 OR AH,0x80`, bit 15, which the `0x4008` mask does not contain.
- The image therefore contains a panel whose show sets a state bit the idle
  gate never tests, a halt that cannot fire. It is not a defect only because
  that class's modality comes from the message loop instead: the
  `DLG-EMPTY-004` shape, twice over.
- `disp:3a8` = 10 hits / 6 owners / 0 in orphan: those two stores, the
  constructor's initialiser `00471a46`, one reader in `FUN_004b72c0`, and one
  clearing arm, `004761b7 AND CH,0x7f`.

**Confidence.** High. The loop, its terminating message, the enumeration with
its orphan count, the `+0x3a8` special case and its clearing arm are each named
instructions.

## The npc tag, the speaker and its figure

| ID | Claim | Confidence | Status | Evidence |
|---|---|---|---|---|
| DLG-NPCTAG-018 | The number in `<npc=N>` is the `npc<n>` section id of `scenario.res::npc.reg`, established by a subscript path, not by two shipped files agreeing. | High / Medium | ● active | [EXP-0141](../experiments/EXP-0141-npc-fields-and-tag/) |
| DLG-NPCTAG-019 | On EN the `<npc=` tag occurs in four `main.res` text families, which corrects two shipped-usage points of `DLG-MARKUP-007`. | High / Medium | ● active (amended) | [EXP-0141](../experiments/EXP-0141-npc-fields-and-tag/) |
| DLG-FIGURE-020 | The speaker's figure is the world figure compositor's output, called on the speaker's own drawable with the stencil surface null. | High | ● active | [EXP-0166](../experiments/EXP-0166-dialogue-dress/) |
| DLG-FIGURE-021 | No layer is excluded at the dialogue site, and the head slot is equipment slot 6, drawn in both halves of the compositor. | High | ● active | [EXP-0166](../experiments/EXP-0166-dialogue-dress/) |
| DLG-SPEAKER-022 | The speaker is a live actor when one matches the section, and otherwise a synthesised drawable with twelve empty equipment slots. | High / Unknown | ● active | [EXP-0166](../experiments/EXP-0166-dialogue-dress/) |
| DLG-SPEAKER-023 | The live-actor predicate is seventeen `Flags` tokens plus two record comparisons and a state gate, and most shipped speakers are on the composed-figure arm. | High / Medium | ● active | [EXP-0166](../experiments/EXP-0166-dialogue-dress/) |
| DLG-DRESS-024 | A live named speaker is drawn in its spawn outfit; mission 40's `npc25` joins its section, its `Data.bin` template and its map label in three independent directions. | High / Medium / Unknown | ● active (amended) | [EXP-0166](../experiments/EXP-0166-dialogue-dress/), [EXP-0192](../experiments/EXP-0192-mission20-party-boundary/) |

### DLG-NPCTAG-018

- `FUN_004c607a` calls the registry loader `FUN_0048c510` on entry
  (`004c60a3`), scans `<`…`>`, lower-cases the tag body (`FUN_00572f3e`,
  `DLG-MARKUP-007`), and requires
  `Find(sprintf("part=%d", caller's part)) >= 0`. Then `Find("npc=")` →
  `004c61de CString::Mid(idx+4, 5)` → `004c621a sscanf(mid, "%d", *[EBP+0xc])`.
- In the caller `FUN_004c5886` that out-parameter is `[EBP-0x24]`, and it
  reaches the npc array unmodified: `004c5ea5 MOV EAX,dword ptr [EBP + -0x24]` /
  `PUSH` / `004c5eb2 CALL 0x00421b46`. The argument passes through
  `FUN_00421b46` (`[EBP+0x8]`) to `FUN_00422b29` to `FUN_00422399` /
  `FUN_00421f46` with no arithmetic in any frame, ending at
  `00422026 MOV ECX,dword ptr [0x005f075c]` /
  `0042202f MOV EAX,dword ptr [ECX + EDX*0x4]` and at `00421d98`…`00421da1`.
  That is the array whose element *i* the loader fills from section `npc<i>`
  (`REG-NPC-088`).
- The same local also feeds `sprintf("npc%02de%sp%d")` at `004c5aa5`, the
  speech leaf, so the tag's number names a file as well as a section. A second,
  gated call at `004c5b3e` is reached only for `21 <= N <= 24` (`REG-NPC-090`).
- `FUN_00422399` uses the record to decide whether an already-live actor is
  this npc: it compares `record+0xc` against `actor+0x24` under the `Flags`
  token `Face` (`00422a7a`/`00422a7d`) and `record+0x10` against `actor+0x20`
  under the token `Picture` (`00422acf`/`00422ad2`).
- Corpus, both roots: 59 distinct tag values, the same 59 on each root, and
  59/59 name an `npc<n>` section that has a record; 46 records no tag names.
  That census is corpus agreement and is not the reason for the grade.

G2: the tag's domain is the array bound, 0..255, and the shipped ids reach
132. An author may add a section and tag it with no code change, but a tag
naming a slot with no section dereferences a null array element in
`FUN_00421f46`.

**Confidence.** High for the identity: an unbroken register-to-subscript path,
every hop a named instruction, and no instruction between the `sscanf` and the
`MOV EAX,[ECX + EDX*0x4]` alters the value. Medium for the 59/59 census.

### DLG-NPCTAG-019

- Walking `main.res` node by node on EN, `<npc=` occurs 718 times over 268
  `.txt` nodes in four families: `text/battle` 225 nodes / 535 tags,
  `text/inn` 36 / 135, `text/shop` 4 / 8, `text/training` 33 / 40
  (`evidence/tag-families-en.csv`).
- RU figures come from the raw file, because the RU `MAIN.RES` cannot be opened
  by this repository's reader (`DLG-READER-011`). RU carries 762 tags with the
  identical 59 distinct values.
- Correction 1: `DLG-MARKUP-007` records `sound=` as "never used on either
  root". That census was `text/battle`, its own 225 EN event files. `sound=`
  occurs 7 times on EN, all in `text/inn` (as `sound="npcNmNpa"`,
  `sound="imfNmNpN"` and siblings). The clause is true of its corpus and false
  of the file: three node families and 183 tags lay outside it.
- `npcalive=`, `npcdead=` and `tune=` remain 0 across all four families, so
  those three are dead in shipped EN data by a wider census.
- Correction 2: EN and RU each ship one `iamfigter`, a misspelling of the
  parser's literal `iamfighter`; `iamfigter` occurs in no `rom.exe` string.
  `CString::Find("iamfighter")` cannot match it, so that arm is unreachable for
  that part on either root. `iamfighter` itself is used once on RU and never
  on EN, as `DLG-MARKUP-007` reports.

G2: the vocabulary is fixed in the image. An author may use any of the fifteen
literals in any of the four families, and the three dead literals are usable
today with no code change.

**Confidence.** High for the per-family counts and the two shipped misuses: an
exhaustive node walk of one root and a byte census of both. Medium that no
fifth family carries the tag on RU: the RU count matches the raw-file grep
exactly, but RU was not walked node by node.

**Amended.** EXP-0171 corrects the literal-run enumeration, and `retracted.md`
withdraws it. The withdrawn sentence: "The parser's `.data` literal run is
twenty-four strings at `0x5c1c44`…`0x5c1d04`, of which fifteen are the tag
vocabulary and the rest are `main\text\`, `%d` ×4, `.txt`, `.res` and the
format `npc%02de%sp%d`." Read through the PE section table,
`[0x5c1c44, 0x5c1d0a)` holds twenty-five non-empty strings plus a lone `0x01`
byte at `0x5c1d00`. `%d` occurs five times, not four. `npc`, `Start` and a
second `female` at `0x5c1cc8` are in the run and were not listed. `.txt` and
`npc%02de%sp%d` are not in it: the format string is at `0x5c1c0c`, below the
run, and `.txt` is at `0x5be578`. `Start` is the omission that matters: it is
the `npc.reg` key read at `004c6579` that gates all four speaker-conditional
arms (`DLG-TAGARM-027`). The per-family tag counts, the `sound=` correction and
the two shipped misspellings stand.

### DLG-FIGURE-020

- `FUN_00421b46(npcId)` (`UNIT-PICT-036`) has one caller, `FUN_004c5886`,
  which fetches child 12 at `004c5e4d` and hands it the return value at
  `004c5eb2`/`004c5ebb`.
- It allocates two surfaces through `FUN_00429520(w,h)`: `(0x58, 0x6c)` =
  88 x 108, the one the panel receives, and `(0xa0, 0xf0)` = 160 x 240, a
  scratch canvas.
- It branches on `00421ccc MOV EAX,[EBP-0x14]` / `00421ccf MOV ECX,[EAX+0x18c]` /
  `00421cd5 AND ECX,0x11` / `00421cda JNZ 0x00421d62`. Clear takes
  `UNIT-PICT-036`'s flat-BMP path. Set takes `00421d62 PUSH 0x0` /
  `PUSH [EBP-0x34]` / `PUSH 0x0` / `00421d72 CALL dword ptr [EDX+0x80]` on the
  speaker.
- Slot `+0x80` of both drawable vtables holds `FUN_0045ed10`, read as raw
  dwords `0x00599468 -> 0045ed10` and `0x005994f0 -> 0045ed10`: the world
  figure compositor of `HERO-APPEAR-051` and `UNIT-FIGURE-032`.
- That routine is `RET 0xc`, and its three parameters are read at their
  consumers:
  - arg2 `[ESP+0x8a0]` is the picture surface, cleared at entry;
  - arg3 `[ESP+0x8a4]` is `HERO-FIGURE-057`'s stencil surface, skipped whole
    when zero (`0045ed75 CMP EBP,EBX` / `0045ed79 JZ`);
  - arg1 `[ESP+0x89c]` is a third surface reached only through
    `0045f71c CMP EBP,EBX` / `JZ`.
- The dialogue passes `(0, canvas, 0)`, so the colour pass runs and the
  click-map pass does not.
- The panel then blits a fixed 72 x 96 window of the 160 x 240 canvas to
  `(8,7)` of the 88 x 108 canvas (`REG-NPC-089`).
- Every equipment-tree `.256` on both roots is a single 160 x 240 frame: 59 face
  sheets, 780 of 784 `primary/` and 144 `secondary/` parsed, and 4 `primary/`
  nodes too short to carry one. A layer therefore covers the canvas rather than
  being placed on it.

**Confidence.** High. Each load-bearing run was byte-verified against the raw
image through the PE section table on three roots, 0 mismatches, and both
vtable slots were read as raw dwords at fixed addresses, which consults no call
graph.

### DLG-FIGURE-021

- `FUN_0045ed10` has no parameter that selects layers: `RET 0xc`, three
  parameters, all three surfaces (`DLG-FIGURE-020`).
- The layer set comes only from `drawable+0x15c + 4i` for `i = 0..11`
  (`0045ee58`, `0045ef69`; an empty slot skipped at `0045ee61`) and from
  `drawable+0x18c`. The dialogue's layer set is therefore every other caller's,
  `HERO-FIGURE-059`'s two orders included.
- Equipment slot 6 is the head slot, fixed by two instruments that read none of
  each other's evidence:
  - the `Armors` collection's `Slot` column (`ITEM-ARMSLOT-031`), whose value 6
    selects exactly `Hat`, `Cap`, `Low Hat`, `Helm`, `Soft Helm`,
    `Chain Helm`, `Full Helm` and `Plate Helm` and none of the other 22 rows on
    either root;
  - the compositor, which stamps `primary[5]` with the byte `6` at `0045f6b0`
    and `0045f421`. The stencil byte is the equipment slot number
    (`HERO-FIGURE-058`), and drawable index `k` is slot `k+1`
    (`HERO-APPEAR-050`).
- `primary[]` is based at `ESP+0x28` (`0045ee50`), so `[ESP+0x3c]` is
  `primary[5]`. It is colour-blitted in both halves:
  `0045f53a`…`0045f550 CALL dword ptr [EAX+0x18]` non-mage and
  `0045f2f1`…`0045f307` mage.

G2: an author may put any item in slot 6 through that item's own `Armors` row.
The head layer's presence in the dialogue figure is not separately
controllable, and suppressing it is a code change.

**Confidence.** High. The parameter count is the routine's own `RET`
immediate, and the pushes at the call site account for all three. The four
head-slot runs are byte-verified on three roots. The slot identification joins
a data column to an image immediate through the index rather than through a
name.

### DLG-SPEAKER-022

- `FUN_00422b29(npcId)` calls `00422b39 CALL 0x00422399` and, only when that
  returns 0, `00422b53 CALL 0x00421f46`.
- `FUN_00422399` returns a live client actor (`DLG-SPEAKER-023`).
- `FUN_00421f46` allocates `0x1b0` bytes and constructs at `FUN_0045ae30`
  (`REG-NPC-088`), whose first act is `0045ae3b MOV ECX,0xc` /
  `0045ae40 XOR EAX,EAX` / `0045ae42 LEA EDI,[ESI+0x15c]` /
  `0045ae62 REP STOSD`, clearing the twelve visible-equipment slots.
- Nothing on the synthesiser's path writes them. Two searches:
  - a byte scan of the whole extent `[00421f46..00422399)`, 1107 bytes on each
    of three roots, finds 0 occurrences of the dword `0x0000015c`, the encoding
    any dword displacement of that field must carry;
  - `HERO-APPEAR-054` independently enumerates the array's writers as the four
    message stores in `FUN_004104e8` plus the constructor and the destructor,
    none of which is on this path.
- A synthesised speaker contributes no equipment layer, and its figure is the
  face sheet `graphics\equipment\<figure>\<face>.256` alone (`HERO-DOLL-078`).
- A live speaker's figure carries whatever its twelve slots hold at that
  moment. Those slots change only by a re-send (`HERO-APPEAR-054`), so
  equipping an actor changes what its dialogue figure wears.

**Confidence.** High for the two arms and for the empty array. They are named
instructions, and the negative rests on a published writer enumeration and not
on the extent scan alone, which cannot see a wholesale copy of the containing
object.

**Unknown.** Which arm a given shipped dialogue takes at run time, which is
per-mission state.

### DLG-SPEAKER-023

- `FUN_00422399(npcId)` walks the client actor list at `this+0x9b8` and filters
  by runtime class against `0x599208` (`00422504 CALL 0x0057272f` /
  `0042250b JZ`).
- It combines one term per token present in the section's `Flags`, in this
  order, with each `PUSH`'s address: `Me` `004225c5`, `Mage` `00422612`,
  `Female` `0042265f`, `Hero` `004226ac`, `Human` `004226f9`, `MySex`
  `00422746`, `MyClass` `00422793`, `!Me` `004227e0`, `!Mage` `0042282d`,
  `!Female` `0042287a`, `!Hero` `004228c7`, `!Human` `00422914`, `!MySex`
  `00422961`, `!MyClass` `004229ae`, `Platoon` `004229fb`, `Face` `00422a52`,
  `Picture` `00422aa7`.
- The bits of `+0x18c`: bit 0 `Hero`, bit 1 `Mage`, bit 2 `Female`, bit 4
  `Human` and bit 5 `Me`. Bit 5 is read at `0042253e AND EAX,0x20` and tested
  at `004225e7 CMP dword ptr [EBP-0x38],0x0`. `MySex` and `MyClass` compare
  bits 2 and 1 against the player's own drawable at `[this+0x3f54]+0x18c`.
- `Face` compares `record+0xc` against `actor+0x24` (`00422a7a`/`00422a7d`),
  and `Picture` compares `record+0x10` against `actor+0x20`
  (`00422acf`/`00422ad2`).
- `00422afb MOV AL,byte ptr [EDX+0x15a]` / `00422b01 CMP EAX,0x1` / `JLE`
  rejects any candidate above 1, and the first survivor is returned at
  `00422b13`.
- Corpus, all three roots identically: of 105 `npc<n>` sections carrying a
  `Flags` list, 57 carry `Hero` or `Human` and take the composed-figure arm and
  48 take the flat arm. Of the 36 distinct `<npc=N>` values in
  `main.res::text/battle/m*`, the same 36 on EN and RU, 29 are composed-figure
  and 7 flat.

**Confidence.** High for the token list and the bits: a complete read of one
routine's string operands, each with its own `PUSH` address, and each bit read
at its own `AND`. Medium for the two censuses, which are corpus agreement.

### DLG-DRESS-024

- Mission 40's sole npc-arm placement is `npc25`, the same section its
  dialogue tags. Its `Flags` are `Hero,Face,!Female,!Mage`, and `DataBinID=42`
  resolves `Humans[42] PC_Paladin`, matching the map's `Paladin` label.
- Seven cells equip slots 1, 6, 7, 8, 9, 10 and 12, including a plate helm.
- Across both roots, 166 Humans rows name equipment: 156 fill slot 1, 46 slot 2,
  36 slot 4, 39 slot 5, 110 slot 6, 148 slot 7, 81 slot 8, 64 slot 9, 105
  slot 10 and 123 slot 12.
- A live named speaker is drawn in its spawn outfit.
- `npc25` keeps that outfit across the boundary because its exact `Hero` flag
  makes the npc arm request the player-character typeID overwrite, after which
  the culls leave the worn slots untouched.
- This is not universal to Humans: mission 20 transfers four dressed Humans
  that retain out-of-band table typeIDs and are removed (`PARTY-M20-031`).

G2: dress remains template data plus later equipment changes; persistence is a
separate classification decision.

**Confidence.** High for the mission-40 join and the corrected constructor
reason. Medium for the 166-row census.

**Unknown.** How the other 28 composed-arm tags resolve at run time.

**Amended.** EXP-0192 (`PARTY-M20-031`) narrows the claim, and `retracted.md`
withdraws the universal cross-boundary sentence: "A speaker who joins the party
keeps that outfit across a mission boundary", and "`PARTY-BAND-027` puts every
Humans-arm actor inside the surviving typeID band". Outfit persistence is
conditional on actor survival, not a property of the Humans arm: `npc25`
survives because it is flagged `Hero`, and mission 20's four dressed Humans
join and are then removed. The mission-40 outfit and all equipment counts
stand.

## Message number, mission number, tag arms and sound

| ID | Claim | Confidence | Status | Evidence |
|---|---|---|---|---|
| DLG-MSGNUM-025 | The message number is a dword from the script parameter to the `%02d` field, and the value 255 is byte-identical to the mission-lost sentinel. | High | ● active | [EXP-0171](../experiments/EXP-0171-script-message/) |
| DLG-MISSION-026 | `campaign+0x660` is the mission number and the map file number at once, and its writers are inside the record embedded at `campaign+0x548`, not at that displacement. | High / Medium | ● active | [EXP-0171](../experiments/EXP-0171-script-message/) |
| DLG-TAGARM-027 | The eight conditional markup arms test the player's hero or the speaker, and substring nesting makes four of them fire twice. | High / Medium | ● active | [EXP-0171](../experiments/EXP-0171-script-message/) |
| DLG-SOUND-028 | `sound=` names a wave and plays nothing; the pager plays it, and the name it composes when the tag is absent depends on the resource's own leaf name. | High / Medium | ● active | [EXP-0171](../experiments/EXP-0171-script-message/) |

### DLG-MSGNUM-025

- The sender writes it whole: `004ea1f4 MOV dword ptr [EDX + 0xa],EAX`.
- The client dispatcher's `0xb6` arm reads it back whole at
  `004105d9 MOV ECX,dword ptr [EAX + 0xa]` and posts at `004105eb PUSH 0x433`
  through `[0x00632f5c]`, with wParam the number and lParam 0.
- The handler `FUN_00473110` takes wParam into ESI
  (`00473148 MOV ESI,dword ptr [ESP + 0xb4]`), compares `004738bb CMP ESI,0xff`,
  and stores `004738c1 MOV dword ptr [EBP + 0x414],ESI` before the branch and
  on both paths. `004738c7 JNZ 0x004738df` sends the equal case to
  `004738ce PUSH 0x431`, the mission-lost panel.
- The dispatcher's own lose arm posts that value: `0041062b PUSH 0xff` then
  `00410630 PUSH 0x433`.
- A script raising message number 255 is therefore indistinguishable from a
  lost mission at the first instruction that inspects it. `campaign+0x414`, the
  field `00475dea CMP dword ptr [EBP + 0x414],0xff` reads on the lose chain
  (`MISSION-PATH-015`), is the same field an ordinary announcement's number
  lands in.
- Past the sentinel the arm gates once on an open dialog
  (`004738df TEST byte ptr [EBP + 0x3dc],0x8`, `DLG-LIFE-005`), formats the path
  from `004738f5` and `004738fb`, opens at `0047396c`, and skips the panel on a
  null result (`00473994 CMP ESI,EDI` then `00473996 JZ 0x004739b1`,
  `DLG-ABSENT-003`).
- No `campaign+0x6bc` test occurs between `004738bb` and `004739d8`, so the
  event-text path is not gated on game mode.

G2: message numbers 0..254 and 256 upward are free; 255 is reserved by the
image. Lifting that costs a code change but no shipped bytes, since the corpus
maximum is 25.

**Confidence.** High. Three routines were read at instruction level, both
accesses to `packet+0x0a` are dword on their own operands, and the sentinel is
the same immediate at both ends. The mode clause is scoped to the two address
ranges named and says nothing about the dispatcher.

### DLG-MISSION-026

- Four format strings take it, each referenced once (`EnumRefs imm:`, 1 hit /
  1 owner / 0 orphan each):
  - `battle\m%d\event%02d` at `0x5be58c`, read at `004738fb`;
  - `%d.alm` at `0x5be89c`, read at `00477d32`;
  - `main\text\battle\m%d\briefing.txt` at `0x5be7a0`, read at `00478997`;
  - `m%d` at `0x5be928`, read at `0047a321`.
- The map-file build is gated `00477d29 CMP dword ptr [EBX + 0x6bc],0x2`. The
  other arm takes the map-name CString at `[EBX + 0x6b4]`, which the `-map`
  command-line arm writes at `00473b3b`.
- A campaign mission is loaded by number and any other game by name, so the
  event-text directory number is the map file number by construction rather
  than by convention.
- Corpus, both roots: the 28 `text/battle/m<N>` directories are the 28 campaign
  map numbers, in both directions.
- The field has no writer at its own displacement. `EnumRefs disp:660` returns
  7 hits / 3 owners / 0 in orphan and every hit is a read; `imm:660` returns 0;
  `disp:658`, `disp:668`, `disp:66c` and `disp:670` through `disp:67c` are
  empty or foreign.
- It is the record's `+0x118`, and `0x548 + 0x118 = 0x660`. The writers are
  `00488488` and `00488498` in `FUN_00488460`, `00488448` in `FUN_00488420`,
  `004877f5` in the record constructor `FUN_00487640` with ESI = 0, and
  `00487322` in `FUN_00487190`.
- `FUN_00488460` refuses a number above the record's `+0x4` and reads `-1` as
  the value at `+0x4`. `FUN_00488420` accepts only a number that some entry of
  the `0x4c`-byte list at `[ECX+0x4c]` with count `[ECX+0x50]` carries at its
  own `+0x4`.
- `EnumRefs callto:488460`: 6 hits, 6 owners, 0 in orphan. The campaign's own
  start passes ten, at `00473ac1 PUSH 0xa` with ECX = `campaign+0x548`.

G2: the mission number set is data, not a compiled table; no run of 10, 20, 30
exists anywhere in the image as dwords or as words. A new mission costs a list
entry and an `<n>.alm`, and its text directory has to carry the same number.

**Confidence.** High for the field identity, its four consumers and the mode
gate: each is a named instruction, both literals' reference sets are
enumerated with the instrument and its orphan count stated, and the
embedded-record arithmetic is exact. Medium for the writer set being complete:
a whole-record `memcpy` carries no displacement and would appear in neither
sweep (`docs/INSTRUMENT.md` rules 4 and 9).

### DLG-TAGARM-027

- This claim closes the Medium clause of `DLG-MARKUP-007`, which read these
  arms as tests and not as effects.
- The four `iam*` arms test the player's own hero, subject
  `[EBP-0x14] + 0xd0` dereferenced through `+0x3f54`. Each abandons the tag and
  resumes the scan when tag and bit disagree:
  - `004c622d` `iamfemale` with `004c6253 AND EDX,0x4` and `JNZ`, rejecting
    when the bit is clear;
  - `004c627a` `iammale` with `004c62a0 AND EAX,0x4` and `JZ`;
  - `004c62c7` `iammage` with `004c62ed AND ECX,0x2` and `JNZ`;
  - `004c6314` `iamfighter` with `004c633a AND EDX,0x2` and `JZ`.
- Bit `0x4` is female and bit `0x2` is spellcaster on the actor's `+0x18c`,
  read off these four arms' own polarities.
- The four sex and class arms test the speaker instead, behind two gates: the
  npc section's `Start` key (`004c6579 PUSH 0x5c1cb0`, then
  `004c65b4 TEST EAX,EAX` and `JZ 0x004c66f3`) and a resolved speaker
  (`004c65cb CALL 0x00422399`, then `004c65d7 JZ 0x004c66f3`). Either failing
  skips all four. They are `004c65dd` `female` with `004c65f7 AND EDX,0x4` and
  `JNZ`, `004c661e` `male`, `004c6671` `mage` and `004c66b2` `fighter`, on the
  speaker's own `+0x18c`.
- The scan is one forward pass over the tags. An arm that disagrees jumps to
  `004c6a1d`, which resumes at `004c60b1` while the cursor is not at the
  terminator and returns 0 when it is; the pager turns a 0 into `0x445`, which
  closes the window (`DLG-LIFE-005`). A suppressed tag is therefore not a
  suppressed part: another tag declaring the same part can still match, which
  is how the shipped `iamfemale`/`iammale` pairs are authored.
- An arm that does not disagree falls through to the next arm on the same tag
  body, and `CString::Find` is a substring test. So `iamfemale` also satisfies
  `female` and `male`, `iammale` satisfies `male`, `iammage` satisfies `mage`,
  and `iamfighter` satisfies `fighter`. One guard exists: the `male` arm
  re-tests `004c662f Find("female") == -1` first.
- A part tagged for a female player therefore additionally requires a female
  speaker, and an `iamfemale`/`iammale` pair read to a female player by a male
  speaker satisfies neither tag, so that part is absent and the window closes.
- Corpus, `text/battle`, both roots (`evidence/tag-bodies.txt`): 535 EN tag
  bodies in 225 files and 543 RU in 228. `iamfemale` 9 and `iammale` 9 on each
  root; `iammage` 2 EN and 3 RU; `iamfighter` 0 EN and 1 RU; `female` 20 EN and
  22 RU; `male` 40 EN and 44 RU, with every `iam*` body counted in those.
- `mage` 2 EN / 3 RU and `fighter` 0 EN / 1 RU are matched only through
  `iammage` and `iamfighter`, so neither literal is used standalone on either
  root. `npcalive=`, `npcdead=`, `sound=` and `tune=` are 0.
- The misspelling `iamfigter` (`DLG-NPCTAG-019`) appears once on each root, at
  `npc=52` part 1, and matches no literal.

**Confidence.** High for the eight arms and the two gates: one routine read at
instruction level, every test a named instruction. High for the census: an
exhaustive tag-body walk of both roots' `text/battle` nodes, whose 535 EN bodies
reproduce `DLG-NPCTAG-019`'s independently published figure. Medium that a
nesting-suppressed part is absent in play: the fall-through and the loop tail
are read from branch targets, not observed.

### DLG-SOUND-028

- The parser's arm at `004c66f3` finds `sound=`. When the tag is absent it
  writes the empty literal `0x5f212c` to the caller's out-parameter
  (`004c6709`). Otherwise it skips six characters and takes the value to the
  next quote when the first character is one (`004c6733 CMP EDX,0x22`), or to
  the next semicolon when it is not (`004c680e PUSH 0x3b`), storing it with
  `CString::operator=`. It reaches no sound routine.
- Its only caller is the pager `FUN_004c5886` (`EnumRefs callto:4c607a`: 1 hit,
  1 owner, 0 in orphan).
- The pager takes the resource name at `panel+0x6c`, cuts it at the last
  backslash (`004c59e7 PUSH 0x5c`), lowercases the leaf and tests `Left(5)`
  against `event` (`004c5a19 PUSH 0x5c1c04`, `004c5a1e PUSH 5`, compared at
  `004c5a49`).
- A leaf beginning `event` composes `npc%02de%sp%d` at `0x5c1c0c`; any other
  leaf composes `%sp%d` at `0x5c1c24`. Both are discarded when the `sound=`
  value is non-empty: `004c5c5f` re-tests it and `004c5c69 JZ` skips the
  composition otherwise.
- The name that is used is built at `004c5cc7`, `004c5cfc` and `004c5d34` from
  `speech\` at `0x5c1c34`, the resource path's directory prefix, the value, and
  `.wav` at `0x5c1c2c`.
- The mission-event family is therefore the one family selected by a file-name
  prefix rather than by its directory, and a `sound=` tag on it overrides the
  composed speech name rather than adding to it.
- Corpus, both roots: `sound=` and `tune=` occur 0 times in `text/battle`, so
  the override is unexercised by shipped data. `DLG-NPCTAG-019` found `sound=`
  7 times on EN, all in `text/inn`.

G2: authored audio on a mission announcement needs a `sound=` tag and a
matching `speech.res` node, with no code change and no change to shipped bytes.

**Confidence.** High for the arm, the caller enumeration and the family test:
one routine's own instructions plus a `callto:` sweep with its instrument and
orphan count stated. High for the corpus zero: an exhaustive tag-body walk of
both roots. Medium for the concatenation order, which is read from the operand
order of three CString helpers whose arity was inferred from their other call
sites rather than pinned.

## Inn mission-giver text and voice

| ID | Claim | Confidence | Status | Evidence |
|---|---|---|---|---|
| DLG-ZEROARM-029 | The inn's "nothing to offer" zero arm never advances `InnMission` whatever its shipped text says, and that text is not uniformly backward-only: it splits by data root. | High / Medium | ● active | [EXP-0393](../experiments/EXP-0393-inn-text-closure/) |
| DLG-INNVOICE-030 | Text and voice presence for the inn's mission-giver NPCs do not track each other in either direction, and differ per root independently of registry addressing. | High | ● active | [EXP-0393](../experiments/EXP-0393-inn-text-closure/) |

### DLG-ZEROARM-029

- This claim extends `REG-SCN-064`, which established the mechanism
  (`InnMission[i]==0` speaks a line keyed by the live main mission,
  `00481038`/`0048107a`) without classifying its content.
- It reads every file the shipped corpus lets the zero arm reach: 9 of 13
  addressed combos ship.
- All 9 recap something the player has already done: "...we have obtained the
  Cloak..." (EN/RU Mission110), "My vow is fulfilled. The Envoy is avenged..."
  (EN/RU Mission130), "you had taught them a good lesson... even before we
  reached the city" (demo Mission50), "did you deal with the turtle?" (demo
  Mission90).
- EN and RU Mission50, a shared thirteen-part, mostly past-tense backstory,
  close with an explicit reciprocal proposal, "Help me to find his assassins
  and I will help you find the Cloak and Scroll", which reads as an offer in
  ordinary language.
- EN/RU Mission110 and Mission130 close with a forward-pointing suggestion
  short of a proposal ("Let's go check out the shop", "We should ask the
  shopkeeper").
- Only the pre-release root's three shipped instances (Mission50,
  Mission90 ×2) carry no forward-pointing content at all.
- Four more addressed zero-arm combos ship on no root read here (demo
  Mission60 ×2, Mission110, Mission130) and are not classified: no text exists
  to read.
- Only openings and closings are quoted, bounded to 140 runes per file; each
  file's own byte and tagged-part count is reported alongside, never its full
  text (`evidence/zeroarm-text.txt`).

**Confidence.** High for the mechanism reproduction: an independent walk of the
same registry reaches the same 22/22/24 addressed combos
`REG-SCN-064`/`REG-INNCLOSE-116` already report. Medium for the content
classification: it is a reading of 9 files against a recap/forward-pointing
distinction, not a traced consumer of any in-game flag, and none was found or
claimed to exist. It rules out a single uniform answer across the corpus: a
"never any forward content" reading is refuted by EN/RU Mission50, and a "reads
like an ordinary quest offer everywhere" reading is refuted by the 3 demo
instances and by the registry state itself, which never changes at any of the
9.

### DLG-INNVOICE-030

- Population: the union of all three roots' `text/inn/npc/*.txt` (32 distinct
  names) and the matching voice parts (149 distinct `inn/npc/*.wav` names), in
  `text/inn/npc/` and `speech.res::inn/npc/`.
- 31 of the 32 text paths are asymmetric across roots (`REG-INNCLOSE-116`'s
  closure counts). The one exception, `npc32m51.txt`, ships identically present
  on EN, RU and the pre-release root alike.
- Voice compounds the asymmetry. `npc25m150p-.wav` and `npc90m41p2.wav` ship on
  RU with no EN part of that number.
  `npc25m150pa.wav`/`npc25m150pb.wav`/`npc25m150pc.wav` and `npc30m151p6.wav`,
  `npc25m140p5.wav` ship on EN with no RU part of that number.
- Both directions occur inside files whose EN and RU text both ship, so text
  presence does not predict voice-part completeness.
- The pre-release root ships text and voice for the same four stems
  (`npc32m51`, `npc32m90`, `npc90m50`, `npc93m90`) and neither form for any of
  the other 28 `npc/` files. Two of those four (`npc32m51`, `npc32m90`) also
  carry demo-only voice parts under a distinct `mf_<stem>p<n>.wav` name absent
  from EN/RU entirely, alongside the `npc<stem>p<n>.wav` parts of the same stem.
- `DLG-WIN-001` already names this class, "the inn's NPCs (FUN_00480fd0, ×2)",
  as one of its six entries.

**Confidence.** High. Every figure is a direct per-root file-presence count
over a named, closed union of paths (`evidence/dialogue-presence.csv`). The
alternative "presence is 1:1" is refuted by five named voice counterexamples
running in both directions plus the text-level asymmetry count, not by an
aggregate mismatch count alone.

## Open questions

- Which 96 rows of the 240-row composition the portrait pane frames
  (`DLG-FIGURE-020`). `FUN_00421b46` blits a fixed 72 x 96 window of the
  160 x 240 canvas to `(8,7)` of the 88 x 108 surface the panel shows, with
  source top `0x90 - PortraitY1` or the default `(0x24, 0x8c)-(0x6c, 0xe8)`
  (`REG-NPC-089`). Reading those canvas rows as picture rows needs the
  canvas's row order. `REG-NPC-091` grades the bottom-up reading Medium on the
  flat arm, from 33/33 renders judged by eye, and does not address the composed
  arm, whose filler is the sprite blitter rather than a DIB loader. Bottom-up
  frames the top of the figure; top-down frames the lower legs; the two are up
  to 122 rows apart. The measurement: whether the loader that fills the canvas
  from a BMP walks rows forwards or backwards, and whether `FUN_0044dde0`'s
  colour blit walks them the same way.
- Which arm a given shipped dialogue takes at run time (`DLG-SPEAKER-022`,
  `DLG-DRESS-024`). `DLG-SPEAKER-022` establishes both arms, and
  `DLG-DRESS-024` traces one section to a live placement. The other 28
  composed-figure tags depend on whether an ordinary map-placed human satisfies
  the section's predicate in that mission, which is per-map state. The
  measurement: for each tagging mission, resolve its `.alm` humans' `+0x18c`
  bits and face bytes against the tagged sections' `Flags` and `Face`.
