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
| DLG-WIN-001 | One `0x84`-byte panel class with one constructor draws every line of authored dialogue, and its six calling routines are the complete set of entries into it. | High | ● active (amended) | [EXP-0098](../experiments/EXP-0098-mission-text/) |
| DLG-PATH-002 | Script and outcome announcements share one transport; the failure close is selected by its own stored panel pointer. | High | ● active (partially retracted, amended) | [EXP-0098](../experiments/EXP-0098-mission-text/), [EXP-0101](../experiments/EXP-0101-mission-end/), [EXP-0274](../experiments/EXP-0274-defeat-modes/) |
| DLG-ABSENT-003 | A named event text that does not ship produces nothing: no window, no fallback, no fault, no blocked state; this answers `MISSION-TEXT-005`'s Unknown. | High | ● active | [EXP-0098](../experiments/EXP-0098-mission-text/) |
| DLG-EMPTY-004 | The window's `"Nothing to say"` fallback is reached only by an existing file that yields no part 1, and no shipped file on either root reaches it. | High | ● active | [EXP-0098](../experiments/EXP-0098-mission-text/) |
| DLG-LIFE-005 | Only input ends the display, nothing queues, and an announcement that arrives while any dialog is open is discarded. | High / Medium | ● active (amended, partially retracted) | [EXP-0098](../experiments/EXP-0098-mission-text/), [EXP-0108](../experiments/EXP-0108-panel-modality/) |
| DLG-READ-006 | The event text is read at fire time, not at map load, and it is read twice per announcement. | High | ● active | [EXP-0098](../experiments/EXP-0098-mission-text/) |

### DLG-WIN-001

- `FUN_004217be(name)` allocates `0x84` bytes (`0042180a PUSH 0x84`) and
  constructs the panel with `FUN_004c5682(id=9, 30, 120, 610, 360)`. Rects are
  `{left,top,right,bottom}`, fixed by `FUN_00447d20` being
  `MOV EAX,[EAX+0xc] / SUB EAX,[ECX+0x4]`, `bottom − top`.
- It then adds:
  - a portrait (`FUN_004c4995`, id 12, 30,54–118,168), only when `panel+0x7c`;
  - a text control (`FUN_004be10e`, id 10, at 128,36–428,172 with the portrait
    and 48,36–428,172 without, its base constructor shrinking the height to
    135, `DLG-RECT-037`), whose initial content is the literal
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
- Each of the six formats a name under a `main.res` text prefix. The two
  mercenary-hall owners share `text/inn/mercenary/`, so the six owners name
  five families, and all five ship (`evidence/family-en.txt`: 318/22/14/4/33
  nodes).
- The child rects are panel-relative, and the drawn layout settles it. Read as
  screen coordinates, the portrait's top edge (54) would sit 66 px above the
  panel's own (120). Read as panel-relative, all four children fall inside the
  580×240 constructor rectangle and inside the 488×232 rectangle drawn
  (`DLG-PANEL-035`) with nothing overhanging: 118 < 488, 428 < 488, 172 < 232,
  198 < 232. Four independent literals would have to coincide for that.

**Confidence.** High for the geometry, the two layouts, the button command and
the string-table index: each is a named `PUSH` in one routine read whole, and an
instruction, not a fit, settles the rect convention. High for the entry
enumeration with its instrument stated: `callto:` on the repaired table, which
also scans `.rdata` dwords holding the address, 0 orphan on both sweeps. Its
blind spot is a construction reached only through a vtable slot; the class
census `EnumRefs range:4c52b0:4c58a0` (13 functions, 1468 function bytes,
0 orphan bytes, 0 orphan runs) closes that hole for this class.

**Amended.** The family bullet read "Each of the six formats a name in its own
family, and all six families ship", beside five node counts.
`evidence/family-en.txt` gives `FUN_0047e410` (`inn\mercenary\npc%02d`) and
`FUN_004b4c70` (`inn\mercenary\npc35`) the same prefix,
`text/inn/mercenary/` (14 nodes): six owners, five families. The entry
enumeration and the node counts stand.

A second correction (EXP-0410): the overhang bullet read "all four children
fall inside the 580×240 panel with nothing overhanging: 118 < 580, 428 < 580,
172 < 240, 198 < 240", and the text control bullet gave no height for the live
rectangle. 580×240 is the constructor rectangle, whose size `FUN_004c553f`
replaces with 488×232 when the panel is drawn (`DLG-PANEL-035`); the argument
holds on both sizes. The text control's base constructor shrinks its height to
135 (`DLG-RECT-037`). The child literals and the entry enumeration stand.

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

- Three inputs reach the same command `0x46f`: a click on the button, Enter
  (`0x0d`) and Escape (`0x1b`). The button's own key handler `FUN_004bf182`
  answers Enter by posting `0x46f`; the panel's key slot `vt+0x6c`
  `FUN_004c574d` answers Escape by self-sending it, and has the same arm for
  Enter, reached only when no child takes the key (`DLG-KEYS-040`).
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
closes it": the container dispatcher `FUN_004bd9dc` and the base key handler
`FUN_004c5369` were read and add none, but the panel's key-release and
system-key slots and the senders of `0x445` outside the class were not
(`DLG-KEYS-040`). Medium for which handler answers Enter, which rests on the
dispatcher's child order.

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

A second correction (EXP-0410): bullet 1 read "through `vt+0x6c`
`FUN_004c574d`, the keys `0x0d` (RETURN) and `0x1b` (ESCAPE), both of which
self-send `0x46f`". The container dispatcher offers a key to the panel's
children before `FUN_004c574d`, and the button's key handler takes Enter first,
so `FUN_004c574d` is the route for Escape and only the fallback for Enter. The
command is the same, so the three inputs, the pager loop and the drop stand.
The Confidence paragraph's Medium clause read "`FUN_004c52f3`'s default arm
`FUN_004bd9dc` and the base key handler's siblings were not read line by
line"; both were read.

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
| DLG-MARKUP-007 | The content model is a tag scan, not a grammar; its vocabulary is fifteen lowercase literals, four of which occur in no event file on either root, and `iamfighter` and `fighter` in no EN one. | High / Medium | ● active (amended) | [EXP-0098](../experiments/EXP-0098-mission-text/) |
| DLG-FACE-008 | The whole file decides once whether the window has a portrait, each part decides which portrait it shows, and the two tests differ. | High / Unknown | ● active | [EXP-0098](../experiments/EXP-0098-mission-text/) |
| DLG-WRAP-009 | The text is wrapped and measured against its own font, and the window shows at most 7 lines of 17 px: the control scrolls by key only, and no shipped block has more than 7 lines. | High / Medium | ● active (amended, partially retracted) | [EXP-0098](../experiments/EXP-0098-mission-text/), amended [EXP-0410](../experiments/EXP-0410-dialogue-layout/) |
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
- Never used in either root's event files: `npcalive=`, `npcdead=`,
  `sound=`, `tune=`. `iamfighter` and `fighter` each count once in RU's and
  never in EN's; `fighter` is matched only through `iamfighter`
  (`evidence/markup-*.txt`, `DLG-TAGARM-027`).

**Confidence.** High for the vocabulary and the scan: one routine's own string
operands, read with `StrDump fnstr:`, and the lowercasing is a named call.
Medium for what each conditional literal does: only `part=`, `npc=` and
`tips=` were followed to their consumers, and the four `iam*` arms and the four
npc-flag arms were read as tests, not as effects.

**Amended.** The headline read "five of them are dead in the shipped corpus".
Its own list is four literals unused on both roots plus `iamfighter` and
`fighter`, unused on EN only (`evidence/markup-en.txt` 0, `markup-ru.txt` 1
each), over the 225 EN / 228 RU event files. `DLG-NPCTAG-019` bounds the
`sound=` zero to that corpus: on EN `sound=` occurs 7 times, all in `text/inn`.
It also finds one `iamfigter` misspelling per root, which no arm matches.
`DLG-TAGARM-027` reads the four `iam*` arms and the four speaker (npc-flag) arms
to their effects, and `DLG-SOUND-028` follows `sound=` to the pager; neither
follows `npcalive=`, `npcdead=` or `tune=`, and for those three the Medium
clause stands.

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
  the panel class answers `0x46d`. The control's own key handler `FUN_004c224e`
  moves the top line on PageUp, PageDown, Up and Down while it holds focus
  (`DLG-KEYS-040`); the base key handler `FUN_004c5369` binds Tab and the four
  arrows to focus movement.
- The constructor rectangle is 300×136 px with the portrait and 380×136
  without, and the control's base constructor shrinks its height to 135
  (`DLG-RECT-037`), so with the pitch of 17 the window shows at most 7 lines.
  The longest shipped part body is 230 characters (EN) / 211 (RU), and no
  shipped block wraps to more than 7 lines (`DLG-LINE-038`).
- The per-glyph advance that turns 300 px into a character count is
  EXP-0097's, and no figure here depends on it.

**Confidence.** High for the wrap, the clamp formula, the default pitch and
the control's key scroll: named instructions plus two one-line accessors read
whole. Medium for "no shipped block has more than 7 lines" (`DLG-LINE-038`, a
transcription of the wrapper, not an emulation) and for whether the text
control holds focus when the panel opens, which decides whether its key scroll
is live (`DLG-KEYS-040`).

**Amended.** The headline read "this window clamps the line count rather than
scrolling". The fifth bullet read "The rect is 300×136 px with the portrait and
380×136 without", with no live height, and the fourth ended "The base key
handler `FUN_004c5369` binds Tab and the four arrows to focus movement, not
scrolling". The Medium clause read "text past the clamp is unreachable in this
window": `FUN_004bd9dc`, the panel's default command arm, was not read, and
whether `panel+0x38`'s navigation links are ever filled was not established.
EXP-0410 corrects three points. The live rectangle is 135 px high, so the cap
is `floor(135 / 17)` = 7 lines and not 8. The text control has its own key
handler, so text past the clamp is reachable by key in principle. `FUN_004bd9dc`
was read: it offers a key to the focused child, then to every child in order,
then to the panel's own slot. The wrap, the clamp formula, the default pitch
and the absence of a scrollbar child stand.

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
| DLG-DRAW-015 | Nothing else is drawn differently while a panel is up: seven instructions in six routines read the panel bit, and none of them is in a draw path. | High / Medium | ● active (amended) | [EXP-0108](../experiments/EXP-0108-panel-modality/) |
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
- That leaves seven readers of bit 3, in six routines:
  - `0047167f` (`TEST EAX,0x4008`, the halt gate; `imm:4008` = 1 hit
    image-wide);
  - `004738df` (the announcement drop, `DLG-LIFE-005`);
  - `004761f2` (the close arm, `DLG-CLOCK-014`);
  - four early-outs in three routines:
    `004902b1 TEST byte ptr [EAX + 0x3dc],0xa` in `FUN_00490280`, `00490c9d`
    and `00490fc5` in `FUN_00490c60`, and `004b9c04` in `FUN_004b9bc0`. Each
    returns or jumps past its body when the bit is set.
- None of the seven lies in the terrain/sprite/shroud blit family
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
the four early-outs are input paths rather than draw paths: each was read at
its guard and the two instructions after it, not end to end. `FUN_00490280`
and `FUN_00490c60` read the cursor globals `[0x005cd7a0]`/`[0x005cd7a4]` that
the pacer's own edge-scroll block reads at `0047547d`…`004754fc`, and
`FUN_004b9bc0`'s not-taken path calls the panel command router `FUN_004c52f3`.

**Amended.** The headline read "six instructions read the panel bit" and the
card "six readers of bit 3" and "three early-outs", against seven listed
addresses. `evidence/enum-en.txt` and section 11 of
`evidence/rom-panel-excerpt.md` show seven bit-3 tests in six routines;
`FUN_00490c60` holds two (`00490c9d`, `00490fc5`), so the early-outs are four
instructions in three routines. The draw-path negative holds for all seven
addresses.

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
| DLG-NPCTAG-019 | On EN the `<npc=` tag occurs in four `main.res` text families, which corrects two shipped-usage points of `DLG-MARKUP-007`. | High / Medium | ● active (amended, partially retracted) | [EXP-0141](../experiments/EXP-0141-npc-fields-and-tag/) |
| DLG-FIGURE-020 | The speaker's figure is the world figure compositor's output, called on the speaker's own drawable with the stencil surface null. | High | ● active (amended) | [EXP-0166](../experiments/EXP-0166-dialogue-dress/) |
| DLG-FIGURE-021 | No layer is excluded at the dialogue site, and the head slot is equipment slot 6, drawn in both halves of the compositor. | High | ● active | [EXP-0166](../experiments/EXP-0166-dialogue-dress/) |
| DLG-SPEAKER-022 | The speaker is a live actor when one matches the section, and otherwise a synthesised drawable with twelve empty equipment slots. | High / Unknown | ● active | [EXP-0166](../experiments/EXP-0166-dialogue-dress/) |
| DLG-SPEAKER-023 | The live-actor predicate is seventeen `Flags` tokens plus two record comparisons and a state gate, and most shipped speakers are on the composed-figure arm. | High / Medium | ● active | [EXP-0166](../experiments/EXP-0166-dialogue-dress/) |
| DLG-DRESS-024 | A live named speaker is drawn in its spawn outfit; mission 40's `npc25` joins its section, its `Data.bin` template and its map label in three independent directions. | High / Medium / Unknown | ● active (amended, partially retracted) | [EXP-0166](../experiments/EXP-0166-dialogue-dress/), [EXP-0192](../experiments/EXP-0192-mission20-party-boundary/) |
| DLG-SPEAKER-041 | Among the party's heroes, the stock mercenaries and the placements of `scn:140.alm`, no actor answers `npc62` (Human, Mage, Female, Face 4), so the chapter-140 tavern's dialogue draws the synthesised, unequipped face sheet `fmage\4.256`. | High / Medium / Unknown | ● active | [EXP-0413](../experiments/EXP-0413-speaker-figures/) |
| DLG-SYNTH-042 | When no live actor answers `npc21` to `npc24`, the synthesiser builds a Hero drawable whose face is the `Face` key of the `npc.reg` archetype chosen by its sex and class bits; the tokens `Me` and `!Mage` have no arm. | High / Unknown | ● active | [EXP-0413](../experiments/EXP-0413-speaker-figures/) |
| DLG-FACEBYTE-043 | For an npc-arm Humans actor in zero mode, `FUN_004f9065`'s tail ORs the gender local into bit 7 of `actor+0x4b` after the streamer's face store; a stored face of 0 names sheet 0, which neither root ships, and the load aborts. | High / Medium / Unknown | ● active | [EXP-0413](../experiments/EXP-0413-speaker-figures/) |

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
- RU figures come from a byte census of the raw file, not a node walk. RU
  carries 762 tags with the identical 59 distinct values.
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
the two shipped misspellings stand. `retracted.md` also withdraws the reason
given for the RU raw-file census: "because the RU `MAIN.RES` cannot be opened
by this repository's reader (`DLG-READER-011`)". `DLG-READER-011` records
EXP-0101's repair of that reader, which precedes EXP-0141 and opens RU
`MAIN.RES` as 504 nodes. The RU figures stand as a raw-file census, and so
does the Medium grade for RU carrying no fifth family.

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
- The panel then blits a 72 x 96 window of the 160 x 240 canvas, or a 72 x 92
  window on the default arm, to `(8,7)` of the 88 x 108 surface
  (`REG-NPC-089`, `DLG-PORTRAIT-036`).
- Every equipment-tree `.256` on both roots is a single 160 x 240 frame: 59 face
  sheets, 780 of 784 `primary/` and 144 `secondary/` parsed, and 4 `primary/`
  nodes too short to carry one. A layer therefore covers the canvas rather than
  being placed on it.

**Confidence.** High. Each load-bearing run was byte-verified against the raw
image through the PE section table on three roots, 0 mismatches, and both
vtable slots were read as raw dwords at fixed addresses, which consults no call
graph.

**Amended.** The blit bullet read "to `(8,7)` of the 88 x 108 canvas". In this
card the canvas is the 160 x 240 scratch surface; the 88 x 108 target is the
surface `FUN_00429520(0x58, 0x6c)` allocates and the panel receives
(`evidence/listing-dialogue-portrait.txt`). The blit geometry stands.

A second correction (EXP-0410): the blit bullet read "a fixed 72 x 96 window".
The window is 72 x 96 when the second word of the record's portrait quadruple
is not -1 and 72 x 92 when it is -1 (`DLG-PORTRAIT-036`). The canvas, the
target surface and the destination `(8,7)` stand.

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
- `npc25` keeps that outfit across the boundary by constructor mode
  (`PARTY-M20-031`): its exact `Hero` flag makes the npc arm request the
  player-character typeID overwrite, after which the culls leave the worn
  slots untouched.
- This is not universal to Humans: mission 20 transfers four dressed Humans
  that retain out-of-band table typeIDs and are removed (`PARTY-M20-031`).

G2: dress remains template data plus later equipment changes; persistence is a
separate classification decision.

**Confidence.** High for the mission-40 join and for the constructor-mode
reason `npc25` keeps its outfit, which rests on `PARTY-M20-031`. Medium for
the 166-row census.

**Unknown.** How the other 28 composed-arm tags resolve at run time.

**Amended.** EXP-0192 (`PARTY-M20-031`) narrows the claim, and `retracted.md`
withdraws the universal cross-boundary sentence: "A speaker who joins the party
keeps that outfit across a mission boundary", and "`PARTY-BAND-027` puts every
Humans-arm actor inside the surviving typeID band". Outfit persistence is
conditional on actor survival, not a property of the Humans arm: `npc25`
survives because it is flagged `Hero`, and mission 20's four dressed Humans
join and are then removed. The mission-40 outfit and all equipment counts
stand. The Confidence paragraph graded "the corrected constructor reason", a
label the card did not carry. The graded clause is the `npc25` bullet's
`Hero`-flag reason, which `PARTY-M20-031` names constructor mode; the bullet
and the Confidence paragraph now name it.

### DLG-SPEAKER-041

- `npc62` reads the same on both roots: `Flags` `Human,Mage,Female,Face`, `Face` 4, no `Picture`, no portrait quadruple, no `DataBinID`.
- `FUN_00422b29` returns the first client actor that passes every term of `FUN_00422399` (`DLG-SPEAKER-023`) and otherwise builds the synthesised drawable (`DLG-SPEAKER-022`). For `npc62` the terms read actor bits `0x10`, `2` and `4` of `+0x18c` and the face at `+0x24`. The list is the client's actor map at `this+0x9b8`. The inn view's walk `FUN_00466940` reads the same client fields `+0x9c4`, `+0x9c0` and `+0x9bc` (`00466984`, `004669be`, `004669ca`, in the listing `TAVERN-ORDER-015` cites).
- The searched population is the party's hero drawables and the stock units the tavern builds for the stage's level (`TAVERN-ORDER-015`), and the placements of `scn:140.alm`, added in case the client keeps the mission map's actors. No actor map of the tavern's client was enumerated, and that the client holds the placements of `scn:140.alm` is assumed from the chapter number, not established.
  - `FUN_0045f850` classes a hero drawable by its type id in `[0x20,0x40)` (`0045f877`..`0045f8c5`): bits `(old & 0x80) | 9`, plus Female and Mage from `typeID - 0x21`. Bit `0x10` is never set, so the `Human` term rejects every party hero.
  - The stock units are the 52 shelf rows of the `Data.bin` Humans table (rows 47 to 98: 13 types at four levels). Twelve are mages, `NPC03`, `NPC04` and `NPC05` at each level, with faces 2, 2 and 1 at every level; `NPC04` is the female one. The `Face` 4 term rejects all twelve.
  - `scn:140.alm` has 171 type-6 records on both roots: 170 Unit-arm records, one npc-arm record and no type-key or definition-ID record. No Unit-arm record has a type word in the hero range `[0x20,0x40)` (0 of 6672 EN, 0 of 3366 RU), so its drawable keeps the flags word at 0 (`0045f91f`, `0045b25b`; `UNIT-PICT-035`) and the `Human` term rejects it. The npc-arm record is `npc26` in hero mode, row 43 `PC_Elf` (face 4, male): the `Human` term rejects a hero drawable.
- The tavern grid's cell and left panel for the same record are separate presentations (`TAVERN-TALKPIC-016`, `TAVERN-TALKSTATS-017`) and are not read here.
- Over the 215 Humans rows of each root, the terms human band, mage type id 23 or 24, female by the gender column or its default, and face column 4 pass one row on each root: 208 `M141_MageTooth`, server id 517. Eight rows differ between the roots, none a mage row: rows 54, 67, 80, 93, 159, 162, 165 and 168 have type id 3 on EN and 4 on RU, and the last four also face 1 on EN and 29 on RU. It is placed only by the definition-ID arm of the spawner (id `0x205`, secondary word zero): `scn:141.alm` record 0 on both roots, and `Beast.ALM` (5 records) and `Cross.ALM` (1) beside the EN executable.
- Placement census: 38 EN maps with 8094 type-6 records (unit 6672, npc 15, type-key 2, definition-ID 1405) and 34 RU maps with 3991 (3366, 15, 1, 609); each definition id resolves to exactly one row. Each record was read as the actor its arm builds: the row's class bits (the type-key arm's row is the first Humans row with that type id, `ALM-CLS-052`), then the arm's own store of face (the secondary word's low byte) and of the female bit (flags bit 2), which the type-key arm and the definition-ID arm with a nonzero secondary word make (`004e29d9`..`004e2acc`).
  - The actors that pass the four terms are 57 on EN and 1 on RU: row 208 (7 EN, 1 RU) and, on EN only, 50 records of `Horror.alm` whose mage rows take face 4 and the female flag from the placement. The same counts result when the record's own type word stands for the actor's type id. The type-key records (2 EN, 1 RU) and the zero-mode npc records (11 per root) build none. None of them lies in `scn:140.alm`.
- The synthesised drawable starts `+0x18c` at `0x48` (`00421fe5`); `Human` ORs `0x10` (`0042208e`), `Female` `4` (`004220c5`) and `Mage` `2` (`004220ff`), giving `0x5e`. The record has no `Start`, so the face is `record+0xc`, 4 (`00422358`..`00422367`).
- `FUN_00421b46` takes the composed-figure arm because `0x5e & 0x11` is nonzero. `FUN_0045ed10` selects the directory by `0x5e & 6` = 6, `graphics\equipment\fmage\`, and names the face sheet `graphics\equipment\fmage\4.256` (`%s%d.256`). The archive holds `fmage` sheets 1 to 5 on both roots.
- Layers: no equipment (`DLG-SPEAKER-022`); no hero back layer, since `0x5e & 3` is 2 and the gate wants 1 (`HERO-FIGURE-144`); a horse layer only when `[+0x20]` lies in `0x11..0x15` (`HERO-DOLL-078`). With no portrait quadruple the picture window is the 72x92 default over the `t_back` and `t_border` surface (`DLG-PORTRAIT-036`).

**Confidence.** High for the resolver, the filter terms, the synthesiser's bits and face, the class setter's bits and the sheet named: named instructions and installed data. The executables are one program, so RU adds no confirmation of a code fact; the data facts were measured per root. Medium that no actor answers at chapter 140: the population is the party's hero drawables, the 52 shelf rows and the 171 placements of `scn:140.alm`, each read as the actor its arm builds; the Unit-arm result rests on `UNIT-PICT-035`. It excludes an actor entered by a runtime spawn, a trigger or a path this experiment did not read, and no save's actor map was enumerated. The passing placements of `scn:141.alm` and of the EN root maps reach the tavern's client only if that client holds another map's actors; that was not examined.

**Unknown.** The `+0x20` word of the synthesised drawable. `FUN_00421f46` writes it only under a `Picture` token, the constructor does not write it, and the object comes from `FUN_00572824` and the runtime allocator (`FUN_00554390`, then `FUN_005543b0`, not read). The horse layer draws only if that word lies in `0x11..0x15`. Which rows of the canvas the window frames, and the drawn pixels (`DLG-PORTRAIT-036`).

### DLG-SYNTH-042

- `npc21` to `npc24` read the same on both roots: `Flags` `Hero,Me,Start`, `Hero,Mage,!MySex,Start`, `Hero,!Mage,!MySex,Start` and `Hero,!MyClass,MySex,Start`. None carries `Face`, `Picture` or a portrait quadruple; all four carry `DataBinID` 26.
- A live actor answers first when one passes the terms (`DLG-SPEAKER-023`). `Me` is bit `0x20`, which the client sets on the actor it makes primary (`00411765`..`004117b3`). The synthesiser runs only when no actor passes.
- `FUN_00421f46` starts `+0x18c` at `0x48` and sets its bits from these tokens only (`0042205a`..`0042224a`): `Hero` ORs 1, `Human` `0x10`, `Female` 4, `Mage` 2; `MySex` and `MyClass` OR the primary actor's own 4 and 2 from `[this+0x3f54]+0x18c`; `!MySex` and `!MyClass` OR the inverted values. `Me`, `!Me`, `!Mage`, `!Female`, `!Hero`, `!Human`, `Platoon` and `Face` have no arm. `Picture` sets `+0x20` (`00422275`).
- `Start` reads `scenario\npc.reg`, takes the section chosen by `bits & 6` (`MaleFighter` 0, `MaleMage` 2, `FemaleFighter` 4, `FemaleMage` 6) and stores that section's `Face` key, default 1, as the face (`004222bd`..`00422347`). A record without `Start` takes its own `Face` (`00422358`). The four keys read 5, 3, 1 and 1 in that order on both roots.
- Bits and figure by the primary's sex and class, identical on both roots:
  - `npc21`: `0x49` for every primary; `MaleFighter`, face 5, `mfighter\5.256`.
  - `npc22`: male primary `0x4f`, `FemaleMage`, face 1, `fmage\1.256`; female primary `0x4b`, `MaleMage`, face 3, `mmage\3.256`.
  - `npc23`: male primary `0x4d`, `FemaleFighter`, face 1, `ffighter\1.256`; female primary `0x49`, `MaleFighter`, face 5, `mfighter\5.256`.
  - `npc24`: male fighter `0x4b` `MaleMage` 3; male mage `0x49` `MaleFighter` 5; female fighter `0x4f` `FemaleMage` 1; female mage `0x4d` `FemaleFighter` 1.
  - The 16 results name four sheets, `mfighter\5.256`, `mmage\3.256`, `ffighter\1.256` and `fmage\1.256`; all four ship on both roots.
- Every result has bit 0 set, so `FUN_00421b46` takes the composed-figure arm, with twelve empty equipment slots (`DLG-SPEAKER-022`). The compositor adds the hero back layer, `backm` or `backf` by bit 2, to every result whose `bits & 3` is 1: the results `0x49` and `0x4d`. The results `0x4b` and `0x4f` draw the face sheet alone (`HERO-FIGURE-144`). The `[+0x20]` word that gates the horse layer is unwritten here as in `DLG-SPEAKER-041`.
- The synthesised figure is an archetype selected by the record's tokens and the primary's sex and class. The synthesiser reads no party list and no key names a fixed picture per record. Only `npc21` is independent of the primary. A live actor that passes the terms is drawn instead (`DLG-SPEAKER-023`); the hero arm sets bit 0 and the client sets `Me` on the primary, which are the two terms of `npc21`.
- The filter's and the synthesiser's `MySex` and `MyClass` terms read the primary through `[this+0x3f54]` with no null test (`00422562`, `00422139`).

**Confidence.** High for the token arms, the archetype mapping, the `Face` keys and the table: one routine read whole, the registry of both roots read whole, each sheet looked up in the installed archive.

**Unknown.** Which arm a dialogue takes for a given party, that is, whether a live hero passes the terms. Whether a dialogue can run with a null primary. Which rows of the canvas the window frames (`DLG-PORTRAIT-036`).

### DLG-FACEBYTE-043

- Case: an npc-arm placement (`FUN_004e26bb`, flags bit 0) whose record lacks `Hero`. `FUN_0048ceb0` returns whether the record holds the `Hero` token and the arm passes it to `FUN_004f8e78` as the constructor mode; mode 0 is zero mode. The arm joins the record's `DataBinID` to the Humans row whose server id, slot 24, equals it (`FUN_004de63e`, a scan from the last row down that never tests element 0) and passes the row's name. The arm stores nothing to `+0x4b`. The type-key arm (`004e29d9`..`004e2a04`) and the definition-ID arm (`004e2a9f`..`004e2acc`) do.
- `FUN_004f9065` resolves the row again by name: it scans the table from element 1 and takes the first row whose name equals the given one (`004f9281`..`004f92c8`; `FUN_00447fa0` reaches the runtime's exact string compare `FUN_00553d80`, byte-wise or lead-byte aware by the flag at `0x630854`). The 210 non-empty names of each root are unique, also ignoring case, so both routes reach one row. The 5 rows without a type id have empty names.
- Stores to `actor+0x4b` on the constructor path, in order:
  1. `004f30ff MOV byte ptr [ECX+0x4b],1` in `FUN_004f30a2`: the default.
  2. `FUN_004f974d` passes `&actor+0x4b` (`004f98cc ADD EAX,0x4b`) to the byte slot reader `FUN_00523410` for Humans slot 17, `face`. The reader stores only when the getter `FUN_00523500` does not return -1 (`00523428`, `00523441`), so a column of -1 leaves the default.
  3. `004f95bd MOV byte ptr [EAX+0x4b],CL` in `FUN_004f9065`: the suffix face, when it is above 0. A name with a `.` sets the gender local to 1 when the character after the dot is `f` and to 0 otherwise (`004f9123`..`004f912e`), and the digits after that character give the suffix face (`004f9105`).
  4. `004f9637`..`004f9648`, the zero-mode tail: `MOV EAX,[gender local]`, `SHL EAX,7`, `OR DL,AL` into the byte, `MOV byte ptr [EAX+0x4b],DL`. Only the low bit of the local reaches bit 7.
- The gender local defaults to 1 (`004f90aa`), takes the name suffix when there is one, and then takes Humans slot 18, `gender ( is female? )`, unless that value is -1 (`004f9302`..`004f932c`). A non-zero mode, a `Hero` record, sets the type id to the local plus `0x21` or `0x23` and skips the tail (`004f95c6`..`004f95e4`).
- Input: the gender local. Of 215 Humans rows per root, 184 carry a gender value and 31 default to 1, 26 of them in the human band; the counts are equal on both roots, and no row name holds a `.`.
- Client: the first state message carries the byte (`PAL-FACE-005` identifies the byte and traces its carriage). `FUN_004e935d` calls `FUN_004e7de3` with mask -1 (`004e939a`, `004e93a7`) and `FUN_004e7de3` appends the byte under mask bit `0x4000` (`004e8249`..`004e826a`). The creation path of `FUN_004104e8` passes it to `FUN_0045f850`: for a type id below `0x1a` bit 7 becomes the Female bit `4`, the face is `byte & 0x7f`, and type ids 23 and 24 add Mage (`0045f8c8`..`0045f91a`). The update path stores the raw byte to `+0x24` (`00411eff`).
- Placements: 15 npc-arm placements per root, 4 in hero mode and 11 in zero mode. Each zero-mode placement resolves to one row. The sex term of the record agrees with the modelled bit 7 in 11 of 11 (one `Female`, ten `!Female`), and 10 of 11 pass their own record's filter against the actor built. `npc55` does not: row `F_BrigandLeader3` has face 25 and the record `Face` 21.
- Face 0: `FUN_0045ed10` formats `%s%d.256` with `[drawable+0x24]` (`0045f0da`..`0045f0ec`) and no instruction remaps 0. No figure directory holds sheet 0: `mfighter` 1 to 31, `mmage` 1 to 13, `ffighter` 1 to 10, `fmage` 1 to 5, both roots. The sprite constructor `FUN_00428c50` calls `FUN_004c9f10`, which returns 0 for a missing node, then formats `FATAL ERROR: can't load ` plus the path and calls `FUN_0056f643` (`00428cb6`..`00428d11`), the abort path of `TAVERN-TALKPIC-016`. A figure drawn for a face of 0 aborts the program.
- No shipped source stores face 0: the 210 human-band rows per root have effective faces 1 to 31, the 5 rows without a type id have none, and the type-key and definition-ID placements store no low byte 0 (2 and 1 type-key, 433 and 17 definition-ID records with a nonzero secondary word, EN and RU).

**Confidence.** High for the writer chain, its input and the face-0 code path: named instructions in routines read whole, and the data results per root. Medium that no other store follows the streamer's on this path. Instrument: `EnumRefs` `disp:4b`, `disp:4a`, `disp:49`, `disp:48` and `imm:4b` on the repaired project, reduced to memory-destination writes that start at `0x4b` (byte), `0x4a` (word) or `0x48` (dword, qword) and to `ADD reg,0x4b`, stack bases set aside: 79 hits in 62 owners. Four owners are read whole (`FUN_004e26bb`, `FUN_004f30a2`, `FUN_004f9065`, `FUN_004f974d`); 58 are classed only by constructor vtable constants, virtual slots and callers: 21 constructors of other classes, 8 virtual methods, 29 free functions. Blind spots: bulk copies (`REP MOVS`, `memcpy`), pointer tables, free functions handed the actor, ten non-stack `LEA` forms at `+0x48`, and `FUN_004e3591`, which the placement driver calls after the spawner (`004e1e50`) and which stores 1 at `+0x48` through `[Z+0x3c]` in eight places: `SAV-GRPAI-563` types two of them (`004e419e`, `004e4206`) as `AI+0x48` of a group's AI block, and the other six were not typed.

**Unknown.** Which caller of `FUN_004e935d` carries a placement's first message. Whether any update message carries mask bit `0x4000` for a person: the update path would store bit 7 into the face. The face of a person stored by a save, a trigger or a mod.

## Message number, mission number, tag arms and sound

| ID | Claim | Confidence | Status | Evidence |
|---|---|---|---|---|
| DLG-MSGNUM-025 | The message number is a dword from the script parameter to the `%02d` field, and the value 255 is byte-identical to the mission-lost sentinel. | High | ● active | [EXP-0171](../experiments/EXP-0171-script-message/) |
| DLG-MISSION-026 | `campaign+0x660` is the mission number and the map file number at once, and its writers are inside the record embedded at `campaign+0x548`, not at that displacement. | High / Medium | ● active | [EXP-0171](../experiments/EXP-0171-script-message/) |
| DLG-TAGARM-027 | The eight conditional markup arms test the player's hero or the speaker, and substring nesting makes four of them fire twice. | High / Medium | ● active (amended) | [EXP-0171](../experiments/EXP-0171-script-message/) |
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
  `female`, `iammale` satisfies `male`, `iammage` satisfies `mage`, and
  `iamfighter` satisfies `fighter`. `iamfemale` also contains `male`, but the
  `male` arm re-tests `004c662f Find("female") == -1` and
  `004c663f JNZ 0x004c6671` skips its speaker test for any body containing
  `female`. That is the one guard.
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

**Amended.** The nesting bullet read "`iamfemale` also satisfies `female` and
`male`", beside the `male` arm's own `Find("female") == -1` guard.
`evidence/markup-arms.txt` shows that guard's `004c663f JNZ 0x004c6671`
jumping past the `male` arm's speaker test whenever the body contains
`female`, so an `iamfemale` body reaches the `female` arm's test only. The
four double firings, the two gates and the census stand.

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
| DLG-INNVOICE-030 | Text and voice presence for the inn's mission-giver NPCs do not track each other in either direction, and differ per root independently of registry addressing. | High | ● active (amended) | [EXP-0393](../experiments/EXP-0393-inn-text-closure/) |

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
alternative "presence is 1:1" is refuted by seven named voice counterexamples
running in both directions plus the text-level asymmetry count, not by an
aggregate mismatch count alone.

**Amended.** The Confidence paragraph read "five named voice
counterexamples" against the seven `.wav` names the card lists.
`evidence/dialogue-presence.csv` carries all seven: two RU-only, five
EN-only.

## Panel geometry, line placement and input

| ID | Claim | Confidence | Status | Evidence |
|---|---|---|---|---|
| DLG-PANEL-035 | At 640x480 the panel is drawn 488x232 at (76,124): the 580x240 constructor rectangle is snapped and centred, then painted as a nine-piece `lm.256` frame with an 8 px shadow band. | High / Unknown | ● active | [EXP-0410](../experiments/EXP-0410-dialogue-layout/) |
| DLG-PORTRAIT-036 | The portrait pane is an 88x108 picture surface at panel offset (30,54): a 72x94 black fill, a 72x96 or 72x92 picture window at (8,7), and a border frame whose opening is 72x92. | High / Unknown | ● active | [EXP-0410](../experiments/EXP-0410-dialogue-layout/) |
| DLG-RECT-037 | All six opening routines (seven call sites) build the identical panel, because the constructor takes only a name; every shipped node takes the portrait layout, so the text rectangle is 300x135 at (204,160). | High | ● active | [EXP-0410](../experiments/EXP-0410-dialogue-layout/) |
| DLG-LINE-038 | Each wrapped line is drawn from the rectangle's left, 17 px below the previous one, justified by widening the word gaps except on a paragraph's last line and one-word lines; a paragraph's first line is indented 10 px. | High / Medium | ● active | [EXP-0410](../experiments/EXP-0410-dialogue-layout/) |
| DLG-BUTTON-039 | The button is drawn from lines and a font-1 label with no art and no fill: a two-colour bevel, the label centred at (316,308) in gold ink (brown on hover), and a shadow offset of 2 px (4 while pressed). | High / Unknown | ● active | [EXP-0410](../experiments/EXP-0410-dialogue-layout/) |
| DLG-KEYS-040 | Enter, Escape and a click on the button each raise command `0x46f`, which turns the page or, on the last page, closes the panel; Space has no dialogue action, and the panel has no accept or decline state. | High / Medium / Unknown | ● active (amended) | [EXP-0410](../experiments/EXP-0410-dialogue-layout/) |

### DLG-PANEL-035

- `FUN_004217be` passes the literals 30, 120, 610, 360 to `FUN_004c5682`
  (`DLG-WIN-001`). Only the width and height of that rectangle, 580x240, are
  used: the base constructor `FUN_004c54d1` calls `FUN_004c553f`, which
  replaces the size and centres it on the screen.
- `W' = ((W - 8) / 96) * 96 + 8` and
  `H' = ((((H - 104) + a) >> 6) << 6) + 104`, with `a` = 63 when `H - 104` is
  negative and 0 otherwise. Then `left = (screenW - W') >> 1` and
  `top = (screenH - H') >> 1`. 580x240 becomes 488x232.
- `FUN_00471790` sets the screen words `0x005ea208` and `0x005ea20c` to
  640x480, 800x600 or 1024x768 (`0047199a`, `004719b6`) by the option strings
  `-640`, `-800` and `-1024` it names. The panel is drawn at (76,124),
  (156,184) and (268,268), 488x232 on all three.
- The background bitmap argument is 0 (`FUN_004c4c93` stores it at
  `panel+0x64`), so the panel has no bitmap background. The painter
  `FUN_004c4da2` (`vt+0x30`) draws a frame from `interface/lm.256`
  (global `0x005ef980`) instead.
- It clamps the rectangle to the screen and shrinks its right and bottom by
  8, giving a 480x224 body. The `lm.256` frames measure 0: 96x64, 1: 48x48,
  2: 96x48, 3: 48x48, 4: 48x64, 5: 48x64, 6: 48x48, 7: 96x48, 8: 48x48, the
  same on both roots.
  - Corners: frame 1 at (L,T), frame 3 at (R-48,T), frame 6 at (L,B-48),
    frame 8 at (R-48,B-48).
  - Edges: `(bw - 96) / 96` = 4 tiles of frame 2 at `x = L+48+96i`, `y = T`,
    and of frame 7 at `y = B-48`; `(bh - 96) / 64` = 2 tiles of frame 4 at
    `x = L`, `y = T+48+64j`, and of frame 5 at `x = R-48`.
  - Interior: frame 0 in 4 by 2 tiles from (L+48,T+48).
  - The tiles cover 96 + 96*4 = 480 by 96 + 64*2 = 224, so no body pixel is
    left to a background.
- Frames 3, 5, 6, 7 and 8 are drawn a second time 8 px right and down through
  `vt+0x1c` (blit mode 6), which fills the 8 px band right of and below the
  body.
- `FUN_00476810` darkens the whole screen behind the panel once
  (`DLG-DIM-013`). The root container `campaign+0xcc` spans the screen from
  (0,0) (`FUN_00472070`), so panel coordinates are screen coordinates.

**Confidence.** High for the size, the position and the tiling: `FUN_004c553f`,
`FUN_004c4da2` and `FUN_00471790` are read whole, and every figure is
arithmetic on the nine frame sizes measured from both roots' `lm.256`. The
rival, a panel drawn at its constructor rectangle (30,120)-(610,360), is
refuted because `FUN_004c553f` replaces all four coordinates.

**Unknown.** Whether a frame piece is drawn at an offset stored in its sprite
record (origins were not measured), and what the shadow pieces put in the 8 px
band (blit mode 6 was not decoded).

### DLG-PORTRAIT-036

- `FUN_004217be` adds the portrait child `FUN_004c4995` (id 12, bitmap 0,
  mode 2) at panel-relative (30,54)-(118,168), 88x114, only when `panel+0x7c`
  is set (`DLG-FACE-008`). At 640x480 it is (106,178)-(194,292).
- Its paint `FUN_004c4a9d` first fills (L+8,T+7)-(L+80,T+101), 72x94, with
  black (`FUN_0044f990`), then blits the picture surface at (L,T) keyed on
  colour 0 (`vt+0x38`). The surface is 88x108, six rows shorter than the
  child rectangle.
- `FUN_00421b46` builds the surface (`UNIT-PICT-036`, `DLG-FIGURE-020`), in
  this order: cleared; `interface/t_back.bmp` (160x240, 24 bpp) copied opaque
  through the source window to (8,7); the figure canvas (160x240) copied keyed
  through the same source window to (8,7); `interface/t_border.256` frame 0
  (88x108) drawn at (0,0). The canvas is freed afterwards.
- The source window depends on the second word of the record's portrait
  quadruple at `+0x20` (`PortraitY1`, `REG-NPC-089`), tested at `00421dce`.
  Not -1: (x0, 0x90 - y1)-(x0 + 72, 0xf0 - y1), 72x96, with x0 the first
  word. -1: (0x24,0x8c)-(0x6c,0xe8), 72x92.
- The border frame, measured by walking its control stream under
  `SPR256-RLE-007`, `SPR256-RLE-020`, `SPR256-STRUCT-001` and
  `SPR256-TRLR-021` (identical on both roots): 2231 control bytes, 1626
  written and 7878 unwritten pixels. The unwritten region that contains
  (44,55) has 6506 pixels and the bounding box (8,7)-(80,99), 72x92, and does
  not fill that box. The frame writes 118 pixels inside the 72x92 window, 404
  inside the 72x96 window and 262 inside the black fill.
- The border is drawn last, so the picture the player sees is the part of the
  window that the frame leaves unwritten. At 640x480 the window is
  (114,185)-(186,277) for the -1 case and (114,185)-(186,281) otherwise; the
  four extra rows of the second lie under 286 border pixels.

**Confidence.** High for the child rectangle, the fill, the surface size, the
layer order and the two windows: each is an operand of a named instruction in
`FUN_004c4a9d` or `FUN_00421b46`, both read whole. High for the border
measurement under the stated instrument (a walk of the control stream, not a
render). Which arm a shipped record takes was not classified
(`DLG-SPEAKER-022`).

**Unknown.** Which rows of the 240-row canvas the window frames, the open
question on record (`REG-NPC-091` grades the bottom-up reading Medium for the
flat arm only), and the pixel blend of the keyed copy.

### DLG-RECT-037

- `FUN_004217be(name)` is the only constructor (`DLG-WIN-001`) and takes the
  name alone, and it shows the panel itself (`FUN_00476810`), so no caller has
  a step between construction and display. Every rectangle below is a
  constructor literal offset by the panel origin (`DLG-PANEL-035`).
- The routines, their call sites, name formats and shipped nodes. A node is a
  `main.res` text node whose base name starts with the stem, and every one
  contains `npc`, the whole-file test of `DLG-FACE-008`.

| Routine | Call site | Name format | EN nodes | RU nodes |
|---|---|---|---|---|
| `FUN_00473110`, mission event | `004739ac` | `battle\m%d\event%02d` | 225 | 228 |
| `FUN_00480fd0`, inn NPC | `00481060`, `004810a5` | `inn\NPC\npc%02dm%d` | 22 | 29 |
| `FUN_0047e410`, mercenary hall | `0047e51a` | `inn\mercenary\npc%02d` | 14 | 14 |
| `FUN_004b4c70`, mercenary hall | `004b4d2a` | `inn\mercenary\npc35` | 14 | 14 |
| `FUN_004a8bc3`, shop keeper | `004a90fe` | `shop\npc31m%d` | 4 | 4 |
| `FUN_004b7f20`, training hall | `004b8383` | `training\npc34m%d` | 3 | 3 |

- The two mercenary routines share one folder, so the 14 nodes count once:
  268 EN and 278 RU nodes in all, none without `npc`
  (`evidence/strings.tsv`). Every shipped dialogue therefore takes the
  portrait layout; the no-portrait layout below is not reached by shipped
  data. `DLG-WIN-001`'s 318/22/14/4/33 count every node under each folder
  prefix, a different population.
- The rectangles at 640x480, as left, top, right, bottom:

| Item | Rectangle | Size | Source |
|---|---|---|---|
| panel constructor | 30,120,610,360 | 580x240 | literal; size only |
| panel drawn | 76,124,564,356 | 488x232 | `DLG-PANEL-035` |
| frame body | 76,124,556,348 | 480x224 | `DLG-PANEL-035` |
| portrait child | 106,178,194,292 | 88x114 | literal 30,54,118,168 |
| portrait surface blit | 106,178,194,286 | 88x108 | `DLG-PORTRAIT-036` |
| black fill | 114,185,186,279 | 72x94 | `DLG-PORTRAIT-036` |
| picture window, -1 / other | 114,185,186,277 / 114,185,186,281 | 72x92 / 72x96 | `DLG-PORTRAIT-036` |
| text control, constructor | 204,160,504,296 | 300x136 | literal 128,36,428,172 |
| text control, live | 204,160,504,295 | 300x135 | base constructor |
| text control, no portrait | 124,160,504,295 | 380x135 | literal 48,36,428,172 |
| button | 276,296,356,322 | 80x26 | literal 200,172,280,198 |
| button label cell, EN / RU | 305,300,327,315 / 281,300,351,315 | 22x15 / 70x15 | `DLG-BUTTON-039` |

- At 800x600 and 1024x768 every rectangle moves by (80,60) and (192,144),
  the shift of the panel origin (`evidence/layout.tsv`).
- The text control's base constructor `FUN_004c156b` recomputes the height
  from the font: `n = floor(136 / (h + 4))` rows of `h + 4`, plus 2. With
  font 1 (`h` = 15) that is 7 rows and 7 * 19 + 2 = 135. `FUN_004be10e` then
  sets the pitch to `h + 2` = 17, and the control shows
  `min(floor(135 / 17), line count)` = min(7, line count) lines.

**Confidence.** High. The literals are `PUSH` operands read in one routine;
the derived values are arithmetic on the snap, the frame sizes and the font-1
metrics measured on both roots (`evidence/metrics.tsv`). The stem counts are
directory counts of both roots' `main.res` over the five folders named, so
"every node has `npc`" is a census of those nodes and not a statement about a
node that does not ship.

### DLG-LINE-038

- `FUN_004be3a7` paints the text control through
  `FUN_00457400(font, rect, first, last, lines, ink, pitch)`: `first` is the
  control's top line (0 for every shipped block), `last` is `first` plus the
  visible count (`+0x8c`, at most 7), and `pitch` is 17. The wrapped lines
  come from `FUN_004be229` through `FUN_00456ab0` at the rectangle's full
  width, 300 with the portrait (`DLG-WRAP-009`). Each line keeps a trailing
  space, and the last line of each paragraph piece ends in CR.
- For line `i`, counted from 0 over the whole array of `n` lines:
  - the line is a paragraph's first when `i = 0` or line `i - 1` ends in CR,
    and then `p` = 10, the font-1 `.dat` dword 32; otherwise `p` = 0;
  - it is justified when `i != n - 1` and it does not end in CR;
  - the drawn copy has its CR removed;
  - `x0 = rect.left + p` and `y = rect.top + 17 * (i - first)`.
- A justified line goes to `FUN_00456fe0(x0, y, W, line, ink)` with
  `W = rect.width - p`. It trims the right end, then repeatedly trims the left
  end and cuts a word at the first space, so runs of blanks collapse, and it
  sums the words' widths `m` (`FUN_00456320`). With two or more words
  `gap = (W - sum) / (words - 1)` in floating point.
  - Word `k` is drawn at `trunc(xacc)`, `xacc` starting at `x0` and becoming
    `(xacc + m_k) + gap` after each word, so the first word starts at `x0` and
    the last ends at about `x0 + W`. The natural space width is not used.
  - A one-word line is drawn at `x0`.
- Every other line is drawn whole at `(x0, y)`, left aligned: a paragraph's
  last line (which includes the line before an authored break) and the last
  array line. No line is centred.
- Every draw goes through `FUN_00456b50(x, y, text, 0, ink, 1)`, which calls
  the font's draw slot `vt+0x14`, `FUN_004577f0`, twice: first at
  `(x + 1, y + 1)` with the shadow ramp `0x005e9be8` (flat 8,8,8), then at
  `(x, y)` with the ink ramp `0x005e8878`, whose level 15 is (255,255,255)
  (`MISSION-MSGLINE-056`). The first line's top edge is the rectangle top:
  no ascent offset is added.
- At 640x480 with the portrait, line tops are 160, 177, 194, 211, 228, 245
  and 262; the left is 204, or 214 on a paragraph's first line; the justify
  width is 300, or 290 on that line (`evidence/text.tsv`).
- Shipped population, both roots' five families, from a transcription of the
  splitter, the wrapper and the measure (`tools/dlglayout`, no emulation):
  EN 688 blocks and RU 732, one per tag containing `part=`. The 535 EN and
  543 RU mission event blocks equal the tag-body census of `DLG-TAGARM-027`.
  - Lines per block, 1 to 7: EN 48, 154, 172, 104, 76, 64, 70; RU 74, 171,
    145, 115, 87, 76, 64. No block exceeds 7 lines.
  - Authored breaks (a CRLF inside a block): EN 0 blocks; RU 12 (7 mission
    event, 4 inn, 1 shop).
  - No block has a one-word line that is not a paragraph's last, an unclosed
    line wider than 300, an empty text, a glyph outside `font1.dat`, or a
    wrapper stall, and no block contains a `~` byte, so the tilde markup
    (`TEXT-079`) never reaches a shipped dialogue.

**Confidence.** High for the placement rule: `FUN_00457400`, `FUN_00456fe0`
and `FUN_00456b50` are read whole, and each operand is traced to its source.
The alternatives of one left x for every line and of centring are excluded by
the justify path and by the absence of any centring flag in the call. Medium
for the population figures: the wrapper is transcribed from its listing and
checked by hand traces, not emulated or observed.

**Unknown.** Whether the x87 extended-precision intermediates of the running
sum move any shipped word by a pixel against a double-precision transcription;
it could matter only where a running sum lands within about 10^-13 of an
integer.

### DLG-BUTTON-039

- The button is `FUN_004be6e3`: id 11, panel-relative (200,172)-(280,198),
  80x26, font 1 (`[0x005e88b8]`), command `0x46f`, and the label is
  string-table entry 77 (`DLG-LANG-010`): "Ok" in EN, 22 px wide, and
  "Принять" in RU, 70 px wide. It has no `~` accelerator in either root.
- `FUN_004be98d` paints it from primitives. There is no art and no interior
  fill: the parent repaints its frame under the rectangle first, then the
  label is drawn, then the bevel.
- The two colours are RGB (41,69,63), light, and RGB (7,12,9), dark, each
  channel reduced to the screen format's width and shift. With `R' = R - 1`
  and `B' = B - 1`, lines run through `FUN_0044fcd0` (both ends drawn) and
  pixels through `FUN_00450390`:
  - colour 1: vertical `x = R'`, `y = T+2..B'-2`; vertical `x = R'-1`,
    `y = T+1..B'-1`; horizontal `y = B'`, `x = L+2..R'-2`; horizontal
    `y = B'-1`, `x = L+1..R'-1`; pixel (R'-2,B'-2);
  - colour 2: horizontal `y = T`, `x = L+2..R'-2`; vertical `x = L`,
    `y = T+2..B'-2`; pixel (L+1,T+1).
- Idle, colour 1 is dark and colour 2 light, which is a raised bevel with a
  one-pixel light top and left and a two-pixel dark bottom and right. Pressed,
  which needs the pressed flag (`+0x6c`, set by the mouse press) and the
  cursor inside the rectangle, colour 1 is light and colour 2 dark.
- The label is drawn through `FUN_00456b50` at
  `(L + (R' - L) / 2 + 1, T + (B' - T) / 2)` = (316,308) at 640x480, with the
  anchor flags 0xa (`TEXT-078`), which put the cell at (305,300)-(327,315) in
  EN and (281,300)-(351,315) in RU. The shadow offset is 2 idle and 4 pressed.
- The ink ramp is `0x005e9b88` idle, level 15 = (185,159,73), and
  `0x005e9b48` while the hover flag (`+0x68`, set by the mouse move
  `FUN_004bef8e`) is set, level 15 = (150,90,0). The shadow ramp is
  `0x005e9be8`, flat (8,8,8). The ramp entries are from CPU emulation of the
  executable's own ramp builder `FUN_00457c40` over synthetic memory
  (`evidence/ramps.tsv`); the original was not run.
- The paint ends by testing flag 1 of the control word `+0x18` (`vt+0x20`,
  `FUN_004bd409`, an AND). With it clear the routine remaps the button's
  screen rectangle through `FUN_0044fad0` at level 3 (`004beea2`), the shade of
  `DLG-DIM-013`. The word starts as 3: the base constructor stores 1
  (`FUN_004bc98f`, `004bca21`) and the button constructor ORs 2 (`004be81b`).
  `FUN_004217be` appends the button and shows the panel without another write,
  so the dialogue button is drawn undimmed. No routine read here clears flag 1.

**Confidence.** High for the geometry, the colour assignment, the label
anchor and the offsets: every operand is read in `FUN_004be98d`, listed whole.
High for the ramp values under the stated instrument.

**Unknown.** How the display's colour depth quantises the ramp entries and the
two bevel colours; the reduction to the screen format was read, not applied to
a format.

### DLG-KEYS-040

- A key reaches the panel through the campaign window's key handler
  `FUN_00472b80` (message-map record for `0x100`). Enter and Space take its
  forwarding arm `0x00472e4b`. Escape takes arm `0x00472c4c`, which posts
  `0x416` only when `campaign+0x3dc` is 1 and `0x41f` only when it is 0 with
  `campaign+0x3b4` empty; a panel on screen sets bit 3, so Escape forwards as
  well (`MENU-ESC-010`, `MENU-INPUT-016`). `FUN_00475250` forwards `WM_CHAR`
  the same way. The arm calls `vt+0x48` of the root container `campaign+0xcc`.
- The class's show slot `vt+0x80` (`FUN_004c584a`) calls the base show
  `FUN_004c5275` at `004c585c`. That routine calls `vt+0x24(1)`, which is
  `FUN_004bd456` in the base and dialogue class tables (`0x0059b710`,
  `0x0059b798`): it saves the root's capture object at `root+0x40` and stores
  the panel in `root+0x34`. `FUN_004bd9dc` gives a mouse message
  (`0x200`..`0x206`, `0x400`) to that object before any child
  (`004bda2d`..`004bda50`), so the panel receives the mouse while it is up.
- The root offers the key to its focus object, `root+0x38`, which is the
  panel after `FUN_004c5275` (`vt+0x28(1)`, `FUN_004bd4c5`). The panel's
  `vt+0x48` is `FUN_004c5797`; a message other than `0x46f` goes through
  `FUN_004c52f3` to the container dispatcher `FUN_004bd9dc`, which offers a
  key to the focused child `panel+0x38`, then to every child in order
  (portrait, text control, button), then to the panel's own key slot
  `FUN_004c574d`.
- Enter (`0x0d`): the button's key handler `FUN_004bf182` answers in the
  second step. It posts `0x46f` to the main window and returns 1, so the
  panel's own arm never sees the key.
- Escape (`0x1b`): no child takes it, and `FUN_004c574d` calls
  `vt+0x48(0x46f)` on the panel directly and returns 1. That routine has the
  same arm for `0x0d`; every other key goes to the base handler
  `FUN_004c5369`.
- A click on the button: the press (`FUN_004bf0a0`) sets the pressed flag,
  the release (`FUN_004bf0e7`) clears it and, with the point inside the
  rectangle, posts `0x46f`. The main window's default arm (`00475053`)
  forwards a posted `0x46f` to the root, and `FUN_004c5797` answers it.
- Space (`0x20`) has no dialogue action. The text control takes only
  `0x21`, `0x22`, `0x26` and `0x28`; the button takes only `0x0d`;
  `FUN_004c574d` takes `0x0d` and `0x1b`; the base handler's table covers
  `0x09` and `0x25`..`0x28`. The button's `WM_CHAR` handler `FUN_004bf1d3`
  fires only on its accelerator `+0x74`, which is 0 for label 77.
- `0x46f` in `FUN_004c5797` calls the pager `FUN_004c5886`, which adds 1 to
  `panel+0x68` and looks up that part (`FUN_004c607a`, `DLG-MARKUP-007`).
  - Found: the text of child 10 is replaced (`FUN_004be229`), children 10
    and 12 repaint, and the panel stays open. A dialogue of N parts takes N
    presses.
  - Not found, the last page or a suppressed part (`DLG-TAGARM-027`): when
    `panel+0x80` is nonzero the panel posts `0x45b` to the main window with
    that value, then calls `FUN_004c52f3(0x445)`, which hides the panel
    (`vt+0x84`) and posts `0x44c`, whose arm is `0047333c CALL 0x004757b0`
    (`DLG-LIFE-005`).
- The `tips=` arm of the pager (`004c68d7`..`004c694f`) reads the five
  characters after `tips=` in the part's tag with `%d` into `panel+0x80`.
  The `0x45b` arm of the main window (`00474cef`) does nothing while the
  global `0x005eb52c` is 0.
  - Census over the five families' shipped blocks (`evidence/wrap.tsv`):
    `tips=` is in 12 EN and 12 RU mission event blocks, with the numbers 1 to
    7, and in no inn, mercenary, shop or training block on either root. No
    shipped town dialogue sets `panel+0x80`, so the `0x45b` post exists only
    at the close of some mission event dialogues.
- A key on the focused text control: `FUN_004c224e` acts on `0x21`
  (`FUN_004c1939`), `0x22` (`FUN_004c1993`), `0x26` (`FUN_004c18e1`) and
  `0x28` (`FUN_004c1913`) when focus flag 4 is set. Each moves the top line
  through `FUN_004be2b0`, which clamps it to `[0, lines - visible]`, then
  sends `0x46d` to the panel and returns 1; `0x23`, `0x24`, `0x25` and `0x27`
  return 0. With 7 or fewer lines the range is 0, and no shipped block has
  more (`DLG-LINE-038`).
- On the last page of a quest offer Enter, Escape and the click do the same
  close, so the alternative of Escape declining is refuted. The panel has
  three children and no accept state: its fields are the page `+0x68`, the
  name `+0x6c`, an object pointer `+0x70` that the pager deletes and replaces,
  the payload `+0x74` and its size `+0x78`, the portrait flag `+0x7c` and the
  tips value `+0x80`, written by the constructor, the loader `FUN_004c5f06`,
  the pager and the part lookup `FUN_004c607a`, and by neither key handler.
  Registration does not wait for the close:
  the shop calls `FUN_0048a360` at `004a9110` and the training hall at
  `004b8393`, right after the constructor returns, and the inn appends the
  mission to its queue at `00481060` for `FUN_00480ad0` to commit when the
  player leaves (`REG-SCN-062`, `REG-SCN-064`, `DLG-ZEROARM-029`).

**Confidence.** High for what each handler does with each key and for the
command flow: every arm is a named instruction in a routine read whole, and
the answer to the last-page question follows from one command. Medium that
the button rather than the panel answers Enter (it rests on the dispatcher's
child order, derived by reading and not observed, and both arms send the same
command), that the text control holds focus at open (its base constructor sets
flag 2, the portrait lacks it, and the show routine focuses the first child
with flags 1 and 2, `FUN_004bd0a4`), and that the posted `0x46f` reaches the
panel before another child of the root.

**Unknown.** Whether siblings of the panel under the root react to Space; the
handlers for key release (`0x101`), system keys and the panel's own `0x101` and
`0x102` slots, which were not read; what a click outside the button does once
the panel holds the capture (the panel's mouse slots were not read); senders of
`0x445` and `0x446` outside this class, which were not enumerated; what the
`0x45b` arm shows and where `0x005eb52c` is set.

**Amended.** DIALOGUE-044 and DIALOGUE-045 establish the named captured mouse
slots and outside-button release path. Native event ordering remains Unknown.


## Open questions

- Which 96 rows of the 240-row composition the portrait pane frames
  (`DLG-FIGURE-020`). `FUN_00421b46` blits a 72 x 96 window of the
  160 x 240 canvas, or 72 x 92 on the default arm, to `(8,7)` of the 88 x 108
  surface the panel shows, with source top `0x90 - PortraitY1` or the default
  `(0x24, 0x8c)-(0x6c, 0xe8)` (`REG-NPC-089`, `DLG-PORTRAIT-036`). Reading
  those canvas rows as picture rows needs the canvas's row order. `REG-NPC-091` grades the bottom-up reading Medium on the
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
## Captured mouse input

| ID | Claim | Confidence | Status | Evidence |
|---|---|---|---|---|
| DIALOGUE-044 | The captured dialogue panel routes mouse input to its children; an outside left press returns 0 without a pager or close command, and its root does not retry sibling hit tests. | High / Medium | ✔ promoted | [EXP-0415](../experiments/EXP-0415-dialogue-queue/) |
| DIALOGUE-045 | The dialogue button posts its command only for a release inside its rectangle after a press; a release outside clears the press and posts no command. | High | ✔ promoted | [EXP-0415](../experiments/EXP-0415-dialogue-queue/) |

### DIALOGUE-044

The constructor installs table `0059b798`. Show `004c5275` calls slot `+0x24`,
`004bd456`, which stores the panel in its parent capture pointer `+0x34` and
saves the old pointer at `+0x40`. The campaign root is constructed by
`004bc98f` with table `0059b188`.

`004bd9dc` gives messages `0x200..0x206` and `0x400` to capture first. After
that call it jumps to `004bdad2`; a zero return runs the root's own mouse slot,
not `004bd87e`'s sibling loop. The root's down/up/double/right slots return 0.
This is return-value propagation to the parent dispatcher, not a panel call
into the parent or a second hit test behind the panel.

The dialogue's own slots `+0x4c/+0x50/+0x58/+0x5c/+0x60/+0x64/+0x68` return 0.
Its left-down slot `+0x54` is `004c560d`. When panel `+0x40 == +0x34`, it clears
`+0x40`, recursively hit-tests children through `004bd564`, and calls a found
child's `+0x54`. No hit or a hit on the panel itself reaches `00436df0`, which
returns 0. A point outside all child rectangles is also excluded from the
preceding `004bd87e` child-message pass by `004ade50`'s `PtInRect` call.

Text child down returns 0. Its up sends `0x472` and double click sends `0x444`
to the panel; neither is `0x46f`, and neither enters the panel's `0x445/0x446`
close arm. Its drag route scrolls only with mouse flag 1 and child `+0x90`
nonzero. Scroll helpers are reused from `DLG-KEYS-040`, not reread here.

**Confidence.** High for the named slots, branches, returned values and absence
of pager/close commands in these bodies. The instrument reads 50 selected
functions, raw PE vtable dwords and independent branch destinations; this is
not a whole-image absence claim. Medium for the composed campaign input route:
the capture/root dispatch is static and assumes the shown root/panel identity.

**Unknown.** Native input ordering, actual OS events, capture changes by outside
callers, window teardown during dispatch and panel types with another vtable.
The text scroll helper bodies and drawing/sound callees are outside this read.

### DIALOGUE-045

Button table `0059b308` maps `+0x54` to `004bf0a0` and `+0x58` to `004bf0e7`.
Down requires a parent and enabled flag 1. With pressed field `+0x6c` zero it
calls `004beeb4(1)`. That routine sets `+0x6c = 1` and, when flag 8 is clear,
requests capture in the button's parent, the dialogue panel.

Up has the same parent/enabled prerequisite. With `+0x6c != 0` it calls
`004beeb4(0)`, clearing the flag and restoring capture when flag 8 was set.
Only then does `004bf153` test the release point against the screen rectangle
through `004ade50`; its false branch skips `004bf15c..004bf16e`, the command
post. No pressed flag means no post even for a release inside. The dialogue
constructor binds command `0x46f`; that command's pager/close behavior is
`DLG-KEYS-040`.

**Confidence.** High for the two handler bodies, flag stores, capture calls and
release hit-test branch. These are checked original instructions, not a native
click experiment. The rectangle's right and bottom edges use `PtInRect`.

**Unknown.** Native delivery after a press, arbitrary capture replacement,
disabling or deleting the button between events, and sound/paint side effects.
