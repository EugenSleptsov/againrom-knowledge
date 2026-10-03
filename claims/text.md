# TEXT — a string byte to a glyph on screen

Claims about the path a text byte travels: what the engine does to it between the shipped file
and the glyph blit, which atlas record it lands on, and what governs the whole path. Not a file
format. The pieces it sits on belong elsewhere and are not restated here:
[`claims/spr16a.md`](spr16a.md) owns the atlases as *files* (their pixel grammar, their record
count, the `.dat` sidecar, what the high half holds), and [`claims/res.md`](res.md) owns the
container the strings and atlases are read from. Spec:
[`formats/text/format.md`](../formats/text/format.md). Format of this file:
[registry.md](registry.md). IDs are permanent.

**The one thing to carry away.** The engine has a **language selector** and a **code-page
converter**, and both had been read as absent. `SPR16A-TXT-023` published "no code-page pass" and
therefore "how the RU game displays Russian text is open"; the pass exists, is called on every
byte of every string the engine draws, and is the reason the atlases' Cyrillic sits where it
does. The refuted clauses are in [`retracted.md`](retracted.md).

## Code page, selector and atlas index

| ID | Claim | Confidence | Status | Evidence |
|---|---|---|---|---|
| TEXT-CONV-001 | Every string byte the engine draws passes through a code-page converter first — `FUN_004562f0` — and it is the missing half of the index rule. | High | ● active (partially retracted) | [EXP-0097](../experiments/EXP-0097-ru-text/), amended [EXP-0103](../experiments/EXP-0103-ru-text-domain/) |
| TEXT-LANG-002 | The engine is language-conditional, the switch is one dword, and the game names its own language in its own data. | High | ● active | [EXP-0097](../experiments/EXP-0097-ru-text/) |
| TEXT-INDEX-003 | The atlas subscript is `record = (byte)(conv(b) − 0x20)`, and it is unbounded. | High / Medium | ● active | [EXP-0097](../experiments/EXP-0097-ru-text/) |
| TEXT-FIT-004 | The shipped RU text and the shipped atlases fit the converter exactly, in both directions, and neither fits without it. | High | ● active (superseded) | [EXP-0097](../experiments/EXP-0097-ru-text/), remeasured [EXP-0103](../experiments/EXP-0103-ru-text-domain/) |
| TEXT-IN-005 | A second converter runs the other way, and it is the input path: `FUN_00468a90` maps CP1251's Cyrillic onto the engine's own CP866. | High / Medium | ● active | [EXP-0097](../experiments/EXP-0097-ru-text/) |
| TEXT-LOWER-006 | The case fold is language-conditional too: `FUN_00468ac0` is a CP866-aware `tolower`. | High / Unknown | ● active | [EXP-0097](../experiments/EXP-0097-ru-text/) |
| TEXT-DOM-010 | The drawn-byte converter is NOT one-to-one, and on the Russian selector the failure is exactly 64 collision pairs. | High | ● active | [EXP-0103](../experiments/EXP-0103-ru-text-domain/) |
| TEXT-ALIAS-011 | A colliding byte is not an error: it silently draws the other byte's glyph, in bounds, with no clamp and no substitute. | High / Medium | ● active | [EXP-0103](../experiments/EXP-0103-ru-text-domain/) |
| TEXT-SEL0-012 | The non-injectivity belongs to selector 1 alone; on every other selector the pass is the identity on all 256 values. | High | ● active | [EXP-0103](../experiments/EXP-0103-ru-text-domain/) |
| TEXT-FIT2-013 | The shipped Russian text lies entirely inside the source set, measured on a corpus 46× larger than the one previously reported — and the previous population was selected by a filter that leans toward the answer. | High / Medium | ● active | [EXP-0103](../experiments/EXP-0103-ru-text-domain/) |
| TEXT-CAP-018 | The limit the text path implies, stated for the customisation seam: 160 distinguishable glyphs on the Russian selector, 224 on any other, and the ceiling is the byte, not the atlas. | High | ● active | [EXP-0103](../experiments/EXP-0103-ru-text-domain/) |

### TEXT-CONV-001

The routine is nine tests and two adds: `004562f0 MOV EAX,[0x005eb57c]` / `CMP EAX,0x1` /
`MOV AL,byte ptr [ESP + 0x4]` / `JNZ 00456317` — so on any selector but 1 the argument is returned
untouched — then `004562fe CMP AL,0x80` / `JC` / `CMP AL,0xAF` / `JA` / `ADD EAX,0x30`
(**`0x80..0xAF` → `0xB0..0xDF`**), and `0045630c CMP AL,0xE0` / `JC` / `CMP AL,0xEF` / `JA` /
`ADD EAX,0x10` (**`0xE0..0xEF` → `0xF0..0xFF`**). Every other byte is returned unchanged, and
neither add can carry out of `AL` (`0xAF+0x30 = 0xDF`, `0xEF+0x10 = 0xFF`), so the caller's `AL` is
always clean. `SPR16A-FONT-018` read the atlas subscript as `SUB AL,0x20` on the string byte; the
instruction two before it is `0045788e CALL 0x004562f0`, and the `SUB` is applied to that call's
**result** (`00457893`…`00457899`: `CALL` / `MOV ECX,[ESP+0x28]` / `MOV BL,AL` / `SUB BL,0x20`).
~~The map is **injective** on `0x80..0xFF` — the two moved blocks land on `0xB0..0xDF` and
`0xF0..0xFF`, which no unmoved byte occupies (`evidence/converter-properties.txt`) — so no two
source bytes can collide on one glyph~~ — **REFUTED by `TEXT-DOM-010`.** `0xB0..0xDF` and
`0xF0..0xFF` *are* unmoved bytes and *do* occupy those cells: nothing removes them from the
fall-through path, so the map has 64 collision pairs and image 192 of 256. Two source bytes collide
on one glyph 64 times over, and what is drawn is `TEXT-ALIAS-011`. **Everything else in this row
stands** — the routine, the two adds, the block boundaries, the `CALL` two instructions before the
`SUB`, and the byte-cleanliness of `AL` were all re-read from both installs by EXP-0103 and are
unchanged

**Confidence.** **High** for the routine and the map, which EXP-0103 re-derived byte-for-byte / the
injectivity clause carried **High** and was wrong; see [`retracted.md`](retracted.md). The clause's
own stated reason named its counterexample, which is why it survived review

**Amended.** The injectivity clause struck through above is refuted by `TEXT-DOM-010` (EXP-0103;
[`retracted.md`](retracted.md), REFUTED). The routine, the two adds, the block boundaries, the
`CALL` before the `SUB` and the byte-cleanliness of `AL` stand.

### TEXT-LANG-002

`[0x005eb57c]` has exactly **one writer** image-wide and three readers (`EnumRefs refto:005eb57c` on
the repaired table: 4 hits, 4 owners, **0 in orphan or undisassembled code**; the companion
`imm:5eb57c` returns 0, so no site reaches it as an immediate). The writer is `FUN_00468810`, and it
computes the value from a **resource**: it opens `main\id` (the literal at `0x5bd094`, `004688f4`),
reads it into a stack buffer, takes that buffer's **last character** — `00468931 SCASB.REPNE` /
`NOT ECX` / `DEC ECX` gives the length, then `00468936 MOVSX EAX,byte ptr [EBP + ECX*0x1 + -0xd1]`,
which is `buf[strlen-1]` — and stores `0046893e SUB EAX,0x30` / `00468941 MOV [0x005eb57c],EAX`,
i.e. **the trailing ASCII digit as a number**. Measured, per root: EN `main\id` is 9 bytes
`656e676c6973682030` = `"english 0"` → selector **0**; RU is 9 bytes `7275737369616e2031` =
`"russian 1"` → selector **1**. `FUN_00468810` has one caller (`FUN_004709e0` at `004710f5`).
**`rom.exe` is byte-identical between the two roots** — sha256 `942e9b72…d367d03` on both — so the
entire language difference in the executable's behaviour is this one resource; a consumer that ships
one binary and two data sets is doing what the engine does

**Confidence.** **High.** The writer enumeration names its instrument and that instrument's blind
spot is not live here: `refto:` sees reads and writes through the reference manager, `imm:` covers
the address-as-immediate route and returns 0, and a whole-image `range:` census of the two text
modules shows **0 orphan bytes**. The two selector values are file measurements with the bytes
quoted, and the executable identity is our own hash of both roots

### TEXT-INDEX-003

The subtraction is byte-wide (`SUB BL,0x20`, `00457899`; `SUB AL,0x20` at `00457bde` and
`00456363`/`00456397`), and the result is widened with `AND …,0xff` before it indexes (`004578e3`,
`0045793b`, `00457c10`, `0045636d`). The three accessors it reaches all subscript the frame-pointer
table with **no compare against the record count**: `FUN_00428be0` (`vt+0x34`, the glyph blit) is
`00428be0 MOV EAX,[ECX + 0xc]` / `MOV ECX,[ESP + 0xc]` / `00428beb MOV EAX,[EAX + ECX*0x4]`, and
`FUN_00428b00` (`vt+0x20`, width → frame `+0x0`) and `FUN_00428b10` (`vt+0x24`, height → frame
`+0x4`) have the identical shape. So **a byte below `0x20` reads out of bounds**: `b − 0x20` wraps
to `0xE0..0xFF` = records **224..255**, past the last of the 224 (`SPR16A-FONT-015`), and the table
holds exactly `count` pointers (`SPR16A-RDR-017`) — the loaded pointer is whatever follows it, and
it is dereferenced. There is **no substitute glyph, no clamp and no skip**; the only byte the draw
loop special-cases is record **0** (space), which draws nothing and adds `GetHeight(0) >> 1`
(`00457910 TEST BL,BL` / `JZ 00457931` / `CALL [EAX+0x24]` / `SAR EAX,1`) before falling into the
shared advance. Since the loop is `strlen`-bounded a NUL never arrives, but **a newline or tab does
if a caller passes one** — so line splitting is the caller's obligation, not the text routine's

**Confidence.** **High** for the index rule and for the absence of a bound (six quoted instructions
across three accessors, all three read whole and each 13 bytes long, plus the two `SUB`/`AND` pairs)
/ **Medium** for the *consequence* of an out-of-range record: the load and dereference are read, but
no shipped string was found that reaches the draw routines carrying a control byte, so the fault is
predicted from the instructions rather than witnessed. Whether any caller splits lines before
drawing is EXP-0098's area and is not answered here

### TEXT-FIT-004

Domain side: over the RU root's text surfaces — 30 nodes, **1 952 high bytes** — every high byte
lies in `0x80..0xAF` or `0xE0..0xEF`, the two blocks `TEXT-CONV-001` moves, and **0 lie anywhere
else** (`0xB0..0xDF`: 0; `0xF0..0xFF`: 0 on `MAIN.RES` and `patch.res`, whose 1 895 bytes are the
localised text). Image side: the converter's image is the 64 cells `0xB0..0xDF ∪ 0xF0..0xFF`, and
those cells carry ink on **64/64** in `font1`, `font2`, `font4` and `font5`; in the RU root's
`font2` — the one node the RU release replaced (`SPR16A-FONT-022`) — the high-half ink is
**exactly** that image, `0` cells of ink outside it, while the EN `font2` has **35** (the CP437
accent records the RU build blanked). Closure: under the converter **1 952 / 1 952** RU high bytes
land on a record that has ink, on all four 224-record atlases — **100.00 %**. Under the identity
rule the same bytes score **41.39 %** on font1/font4/font5 and **2.77 %** on font2, whose 54 on-ink
hits are all residue from `world.res`, so the localised text alone scores **0 of 1 895** there: the
RU release's own smallest font would draw **nothing at all** for every Russian byte. Population and
its blind spot are stated in the probe: text nodes are `.txt`/`.ini`/`.lst`/extensionless payloads
plus the embedded string tables of `.reg` and `data.bin`; **a display string embedded in any other
node type is not counted**. ~~30 nodes, 1 952 high bytes~~ — **the figure and the population are
SUPERSEDED by `TEXT-FIT2-013`**, which measures 419 nodes and 87 293 high bytes on the same root and
reproduces the 100 % fit exactly. Both of this row's filters lean toward the answer: the node gate
requires ASCII majority, which a majority-Cyrillic file fails by construction and which admitted 39
of 419 RU text nodes, and the run filter counts a byte as a letter only when it already lies in the
two blocks. **The conclusion is unchanged and now rests on 46× the bytes**

**Confidence.** **High.** This is corpus evidence, which caps at Medium on its own — what lifts it
is that the two sides were measured independently and the fit is exact and two-sided: the code's
*image* (read from instructions, before the census ran) equals the atlas's *ink set*, and the
corpus's *alphabet* equals the code's *domain*, with zero exceptions on either. The rival — that
strings and atlases agreed all along and no converter is needed — is refuted by the same numbers at
0/1 895 on font2. The figures carried **High** when believed and were measured through a biased
population; see [`retracted.md`](retracted.md)

**Amended.** The figure and the population struck through above are superseded by `TEXT-FIT2-013`
(EXP-0103; [`retracted.md`](retracted.md), SUPERSEDED). The conclusion stands.

### TEXT-IN-005

`00468a94 CMP AL,0x80` / `JC` returns any ASCII byte **untouched** — this arm calls no CRT routine
at all — then `00468a98 MOV ECX,[0x005eb57c]` / `DEC ECX` / `JNZ` returns unchanged on any other
selector. On selector 1: `00468aa1 CMP AL,0xC0` / `JC` / `CMP AL,0xF0` / `JNC` /
`00468aa9 ADD AL,0xC0`, and `00468aac CMP AL,0xF0` / `JC` / `00468ab0 ADD AL,0xF0`. `ADD AL,0xC0` is
`− 0x40` and `ADD AL,0xF0` is `− 0x10` in byte arithmetic, so the map is `0xC0..0xEF → 0x80..0xAF`
and `0xF0..0xFF → 0xE0..0xEF` — exactly the inverse of the block layout the display side assumes,
from the encoding a Windows edit control hands over. **6 call sites, 6 distinct owners, 0 orphan**
(`EnumRefs callto:00468a90`). This is the reason the RU root's `README.TXT` reads as CP1251 while
its data strings read as CP866 (`SPR16A-TXT-023` measured both and could not join them): they are
two different sides of the same boundary

**Confidence.** **High** for the map (every branch and both adds are quoted from an 18-instruction
routine read whole) and for the call-site count (instrument named, 0 orphan) / **Medium** for
calling it *the input path*: the six owners were enumerated but not read, so which surface each
serves — typed hero name, chat, a file read — is not established here

### TEXT-LOWER-006

`00468ac0 MOV EAX,[0x005eb57c]` / `DEC EAX` / `JZ 00468adb` selects the arm; **every other selector
falls through to the CRT `tolower`** at `0x005562d0` (`00468ad2`), so the English build's folding is
the C library's and nothing else. On selector 1: `00468adf CMP AL,0x80` / `JC` / `CMP AL,0x90` /
`JNC` / `00468ae7 ADD AL,0x20` folds `0x80..0x8F → 0xA0..0xAF`, and `00468aea CMP AL,0x90` / `JC` /
`CMP AL,0xA0` / `JNC` / `00468af2 ADD AL,0x50` folds `0x90..0x9F → 0xE0..0xEF`; anything else
reaches the CRT call at `00468af5`. Those are exactly CP866's two uppercase blocks onto its two
lowercase blocks — the same block geometry `TEXT-CONV-001` moves, which is what makes the three
routines one system rather than three coincidences. **3 call sites, 3 distinct owners, 0 orphan**
(`EnumRefs callto:00468ac0`); one owner, `FUN_004c3118`, calls **both** this and `TEXT-IN-005`'s
routine

**Confidence.** **High** for the fold and the CRT fallthrough (both arms quoted; the fallthrough is
a `CALL` to a named address on two of the three exits) and for the call-site count (instrument
named, 0 orphan) / **Unknown** what the three callers do with it — no case-insensitive comparison
was traced to a surface. Note this fold is **not** the archive path fold: `RES-CASE-036`'s
`FUN_00456580` is a different routine reaching the CRT `tolower` directly, and it is not
selector-gated

### TEXT-DOM-010

`FUN_004562f0` is 42 bytes read whole out of both roots' `rom.exe` by virtual address
(`evidence/converter-domain.txt`; the two dumps are identical and the probe checks its own
transcription against them). The first arm's image is `0xB0..0xDF`, and **nothing removes
`0xB0..0xDF` from the fall-through path**: `00456304 JA` sends every byte above `0xAF` to
`0045630c`, and `0045630e JC` sends everything below `0xE0` straight to the plain `RET 0x4` at
`00456317`. The same shape returns `0xF0..0xFF` unchanged via `00456312 JA`, onto the second arm's
image. So under selector 1 the map has **image 192 of 256**, **64 output values with two preimages
each** — `0xB0` from `0x80` and `0xB0` … `0xDF` from `0xAF` and `0xDF` (48 pairs), `0xF0` from
`0xE0` and `0xF0` … `0xFF` from `0xEF` and `0xFF` (16 pairs) — and it is **not injective on
`0x80..0xFF`**. The routine has no third test and no default arm; the collision is structural, not a
data accident. There are 2^64 maximal one-to-one subsets of size 192, so **the code alone does not
name a domain**: two are natural, the source set `0x00..0xAF` u `0xE0..0xEF` and the map's own image
`0x00..0x7F` u `0xB0..0xDF` u `0xF0..0xFF`, and only the shipped text picks between them
(`TEXT-FIT2-013`)

**Confidence.** **High.** Every branch is a quoted instruction from one 42-byte routine read end to
end and verified byte-for-byte against both installs, and the map was then enumerated over the
**whole** 0..255 range rather than sampled, for both selector values. This **refutes**
`TEXT-CONV-001`'s clause "the map is injective on `0x80..0xFF`", which carried High; that clause's
own stated reason — "the two moved blocks land on `0xB0..0xDF` and `0xF0..0xFF`, which no unmoved
byte occupies" — names its counterexample, since those two ranges *are* the unmoved bytes occupying
them. See [`retracted.md`](retracted.md)

### TEXT-ALIAS-011

Under selector 1, byte `b` in `0xB0..0xDF` selects record `b − 0x20` = 144..191 — **the same
record** as byte `b − 0x30` — and byte `b` in `0xF0..0xFF` selects record 208..223, the same as
`b − 0x10` (`evidence/conv-table.csv` gives the aliasing partner for all 256 values). Every one of
those records is **inside** a 224-record atlas, so unlike `TEXT-INDEX-003`'s sub-`0x20` case the
subscript does not even leave the frame table: the blit is well-formed and a real glyph appears.
`FUN_004562f0` has two `RET`s and no default branch, and the three accessors subscript with no
compare — `00428b07`, `00428b17` and `00428beb` are all `MOV reg,[base + idx*4]`, each routine read
whole. Per `SPR16A-FONT-020` the records concerned hold Cyrillic, so the visible consequence is that
the RU build cannot show the 64 characters those source bytes name in their own page: each draws a
Russian letter instead

**Confidence.** **High** for the aliasing and for the absence of any error arm — the record
arithmetic is the byte-wide `SUB`/`AND` pair already published in `TEXT-INDEX-003`, re-read here,
and all three accessors were read end to end rather than decompiled / **Medium** for *which letter*
appears, which rests on `SPR16A-FONT-020`'s reading of the atlas rather than on anything measured
this round

### TEXT-SEL0-012

`004562f5 CMP EAX,0x1` / `004562fc JNZ 00456317` reaches a bare `RET 0x4` with `AL` untouched, so
the returned byte is the argument. Enumerated over the whole range at selector 0: **0 bytes moved,
image 256, 0 collisions** (`evidence/converter-domain.txt`). Two consequences a consumer must carry.
First, the two selectors do not merely differ in *which* glyph a high byte draws — they differ in
whether the mapping is **information-preserving at all**, so a reimplementation cannot model the
pass as one table with a language parameter unless that table is allowed to be non-injective in one
column. Second, the reachable record set differs: at selector 0 all 224 records of a 224-record
atlas are selectable, at selector 1 only 160 (`TEXT-CAP-018`). The upper 24 bits of the returned
`EAX` are the selector's on this arm, and are clean only because the shipped selectors are 0 and 1;
every caller reads `AL`, so it does not matter, but a consumer returning a full word would inherit a
bug the original does not have

**Confidence.** **High.** Two quoted instructions and a bare `RET`, plus an exhaustive 256-value
enumeration at that selector. The alternative this rules out is the one a reader of `TEXT-CONV-001`
would form — that the converter is a code-page table applied always — and it is ruled out by the
`JNZ` target being the function's own exit

### TEXT-FIT2-013

Population, stated so it can be attacked: every node of every archive on the root whose
**extension** is `txt`/`ini`/`lst`/none and which passes a **code-page-neutral** gate (under 2 %
control bytes; it counts control bytes, never letters). RU: **419 nodes, 146 709 bytes, 87 293 bytes
at or above `0x80`, of which 0 lie in `0xB0..0xDF` and 0 in `0xF0..0xFF`** — 100.0000 % inside
`0x80..0xAF` u `0xE0..0xEF`. EN control, same probe: 413 nodes, 132 727 bytes, **one** high byte in
the entire corpus (`0xE5`). Widening to every NUL-terminated run in every non-sampled node adds only
`.alm` and `.bin` residue, and that residue is **root-invariant** — 1 002 bytes outside the set on
RU against 1 006 on EN, 111 against 111 — so it is binary the run heuristic mistakes for text, not
localised text. **Why `TEXT-FIT-004`'s 1 952 over 30 nodes does not reproduce:** its node gate
requires ASCII majority, which a single-byte Cyrillic file fails by construction, admitting **39 of
419** RU text nodes; and its run filter credits a byte as a letter only when it already lies in the
two blocks. Run over the same nodes in one pass, that filter reports 99.78 % fit where the neutral
one reports 98.75 %, and on EN 52.49 % against 26.99 % — it roughly doubles the apparent fit by
declining to look outside the blocks (`evidence/corpus-en.txt`, `evidence/corpus-ru.txt`)

**Confidence.** **High** for the count itself — it is exhaustive over a mechanically defined
population, the search for a counterexample was neutral over the whole high range by construction,
and zero were found / **Medium for the consequence** that the RU release never relies on a colliding
byte, and the cap is not negotiable: this is corpus evidence, and the *population* is a filter
choice that no reading of the engine forces. What would lift it is a decode of the nodes that carry
display strings, which this round did not do

### TEXT-CAP-018

The subscript is byte-wide (`SUB BL,0x20`, then `AND …,0xff`), so a font can address at most **256**
records; the shipped atlases provide **224**, and `font3` 64. Enumerated over all 256 inputs: at
selector 0 every one of the 224 records is selectable and 32 bytes (`0x00..0x1F`) read past the end;
at selector 1 only **160** are selectable — 96 from `0x20..0x7F`, 48 from the first arm's image
`0xB0..0xDF`, 16 from the second's `0xF0..0xFF` — and **64 records are unreachable dead weight**,
exactly the character codes `0x80..0xAF` and `0xE0..0xEF`. So a Russian string can name 160 distinct
glyphs, not 224, and the missing 64 are not missing art but unaddressable slots. Three separate
things would each have to change to lift it, and a consumer should know which: adding records past
224 needs only the atlas and its `.dat`, because the readers validate no frame header at all
(`SPR16A-RDR-017`); reaching records 224..255 needs the byte-wide subscript widened, which changes
code and not files; reaching the 64 dead records at selector 1 needs a third test inside the
converter, which changes the meaning of every shipped Russian byte and so **cannot** be done without
rewriting the text. That last one is the seam: the code-page pass is where language support is
bolted on, and it is not extensible in place

**Confidence.** **High.** The arithmetic is an exhaustive enumeration over all 256 values against
record counts taken from each shipped file's own trailer, with no free parameter, and the three
"what would have to change" clauses each name the instruction or the file that carries the limit.
What this does **not** establish is whether anything downstream — a save file, a network path, a
`.dat` consumer — assumes 224 rather than 160; no such sweep was run

## Text API, atlases and tilde markup

| ID | Claim | Confidence | Status | Evidence |
|---|---|---|---|---|
| TEXT-API-007 | The text API is four routines on one class, and all three that touch a string byte convert. | High / Medium | ● active | [EXP-0097](../experiments/EXP-0097-ru-text/) |
| TEXT-FONT3-008 | One shipped atlas cannot render Russian at all, and it is a structural fact rather than a gap in its art. | High / Medium | ● active | [EXP-0097](../experiments/EXP-0097-ru-text/) |
| TEXT-TILDE-009 | `~` is markup, not a character, and a consumer that measures text must know it. | High / Unknown | ● active (amended, partially retracted, superseded) | [EXP-0097](../experiments/EXP-0097-ru-text/), amended [EXP-0103](../experiments/EXP-0103-ru-text-domain/) |
| TEXT-FONT-015 | Every value the converter can produce reaches a record that has ink, on every 224-record atlas, on both roots — so "converts correctly, then draws nothing" does not happen. | High / Medium | ● active | [EXP-0103](../experiments/EXP-0103-ru-text-domain/) |
| TEXT-FONT2-016 | `font2` is the only atlas the roots do not share, and the RU build blanked its entire Latin high half — which is free under selector 1 and destroys text under selector 0. | High / Medium | ● active | [EXP-0103](../experiments/EXP-0103-ru-text-domain/) |
| TEXT-TILDE2-017 | The two byte loops do not test `~` at the same point: the measurer tests the raw byte, the draw tests the converted one. The conclusion survives; the published reason does not. | High | ● active | [EXP-0103](../experiments/EXP-0103-ru-text-domain/) |
| TEXT-079 | A lone `~` in the font1-3 draw underlines the next glyph and takes no pen advance: a line from x, that glyph's advance-table width long, one row below frame 0, in ramp entry 15. `~~` is one literal glyph; the font4 draw has no tilde arm. | High / Unknown | ● active | [EXP-0410](../experiments/EXP-0410-dialogue-layout/) |

### TEXT-API-007

The font object is 0x10 bytes — `+0x0` vtable `0x599138`, `+0x4` the glyph sprite, `+0x8` the `.dat`
advance table, `+0xc` the letter spacing (`SPR16A-FONT-018`) — and its method table has 12 slots of
which `+0x14` is `FUN_004577f0` (draw) and `+0x18` is a two-instruction getter. Beside the vtable
sit `FUN_00457b50` (a second draw, taking a length rather than a flag word), `FUN_00456320` (measure
a string's width) and `FUN_004563e0` (a wrap/layout pass that owns no byte loop of its own — it
calls the measurer five times and never indexes an atlas). The converter's callers are exactly the
first three: **5 call sites, 3 distinct owners, 0 orphan** (`EnumRefs callto:004562f0` — `0045635b`,
`0045638f` in the measurer; `0045788e`, `004578a7` in the draw; `00457bd6` in the second draw). Four
font objects are constructed, all in `FUN_00458780`, from whole string literals
`graphics\font1\font1` … `graphics\font4\font4` (`0x5bc95c`, `0x5bc944`, `0x5bc92c`, `0x5bc914`)
with spacing `2`, into four globals: **font1 `[0x005e88b8]` 179 refs / 48 owners, font2
`[0x005e92f8]` 135 / 22, font3 `[0x005e9c28]` 4 / 4, font4 `[0x005e9bc0]` 89 / 20**, all with 0
orphan hits

**Confidence.** **High** for the converter's call-site set *within the text module*: a `range:`
census of `00456260..00456400` and `00457600..00458900` reports **0 orphan bytes and 0
undisassembled bytes that are code** (17 and 68 bytes of padding), so no hidden routine sits between
these and the enumeration cannot have dropped one there / **Medium** for "every text path in the
image converts", and the blind spot is named and is real: the glyph blit is reached as `vt+0x34`, so
a caller that computed a subscript itself and dispatched virtually would appear in **no** `callto:`
sweep. None was found; none was ruled out / **Medium** for the per-font reference counts as a
measure of *surfaces* — they count references, and the owners were not read

### TEXT-FONT3-008

`font3.16` has **64** records, not 224 (`SPR16A-FONT-015`), so under `record = char − 0x20` it
covers characters `0x20..0x5F` only — no lowercase ASCII and **none** of the converter's image:
`0 of 64` of the cells `0xB0..0xDF ∪ 0xF0..0xFF` exist in it, on both roots. Its ink is narrower
still: **11 records** carry any, at characters `0x30..0x39` and `0x46` — the ten digits and one
letter — with `53 of 64` advances zero. It is constructed like the others and read from three sites
(`FUN_00459c50`, `FUN_004588a0`, `FUN_004b08b0`). So the area's question has no single answer: any
surface drawn with `font3` is numeric whatever the selector says, and a consumer must not treat the
four atlases as interchangeable

**Confidence.** **High** for the record count, the covered character range and the ink census (the
count is the file's own trailer under an already-closed walk, and the range follows from the index
rule with no free parameter; the ink figures are a decode of every record on both roots) /
**Medium** for which *surface* `font3` serves — its three readers were enumerated, not read

### TEXT-TILDE-009

~~Both byte loops test it, at the same point and in the same way, on the value **after**
conversion~~ — **SUPERSEDED by `TEXT-TILDE2-017`**: the draw tests the **converted** byte
(`004578ae CMP BL,0x5e`, after the `CALL` at `0045788e`), the measurer tests the **raw** one
(`0045634c CMP AL,0x7e`, nine bytes *before* its `CALL` at `0045635b`). The two agree on all 256
inputs anyway — nothing but `0x7e` converts to `0x7e` — so the rule below is unchanged, but a
consumer must not convert before testing in the measurer. The two addresses this row already gave
are correct; the sentence joining them was not: `004578ae CMP BL,0x5e` in the draw and
`0045634c CMP AL,0x7e` in the measurer (`0x5e` is `0x7e − 0x20`, the same char). **Doubled**, it is
the literal glyph: the draw takes `004578b7 CMP AL,BL` / `JZ 00457909` into the ordinary path and
then skips its partner (`00457950 CMP byte ptr [ESP+0x20],BL` / `JNZ` /
`00457956 INC dword ptr [ESP+0x2c]`), and the measurer does the same at `00456383`. **Single**, it
draws a **rule** and occupies no width: the draw calls `0044fcd0` — a Bresenham line routine — with
the colour taken from `word ptr [arg+0x1e]` (`004578c8`), from `(x, y + GetHeight(0))` to
`(x + dat[j], y + GetHeight(0))` (`004578d3`…`004578ff`), `j` the glyph index of the byte after
the `~` (`TEXT-079`), and then **jumps past the advance
block** (`00457907 JMP 0045794b`), so `x` does not move; the measurer's single-`~` arm likewise
reaches `004563c3` without adding anything to its running total. The two routines therefore agree,
which is what makes this a rule rather than a quirk of one of them

**Confidence.** **High** (both arms of both routines are quoted, and the two agree — a consumer's
width calculation and its draw would desynchronise if either were read wrong) / **Unknown** what the
rule is *for*, and which shipped strings use it: the five dialogue families and string 77 hold no
`~` (`TEXT-079`), and no other text was searched. The "after conversion"
clause carried **High** and was wrong for the measurer; see [`retracted.md`](retracted.md)

**Amended.** The clause that both byte loops test `~` at the same point and after conversion is
superseded by `TEXT-TILDE2-017` (EXP-0103; [`retracted.md`](retracted.md), SUPERSEDED): the measurer
tests the raw byte. The markup rule stands.

A second correction (EXP-0410): the single-`~` sentence read `to (x + dat[0x5e], y + GetHeight(0))`,
and the Unknown read "the argument struct whose `+0x1e` supplies the colour was not identified" and
"no census of `~` in the shipped text was run". The line ends at `x + dat[j]`, `j` the recoded byte
after the `~`, and the colour word is entry 15 of the fifth argument, the ink ramp (`TEXT-079`). The
dialogue families and string 77 were counted: no `~`. The rule, the doubled case, the no-advance
jump and the measurer's arm stand.

### TEXT-FONT-015

Record counts taken from each file's own trailer under a walk that is exact: `font1` 224, `font2`
224, `font3` 64, `font4` 224, `font5` 224, identically on both roots, with 0 bytes of slack between
the last record and the trailer on `font3`, `font4` and `font5`. Decoding **every** record of every
sheet on both roots (`evidence/font-records-en.csv`, `evidence/font-records-ru.csv`, 960 rows per
root): of the 64 character codes `0xB0..0xDF` u `0xF0..0xFF` that are the converter's image, **64 of
64 carry ink** in `font1`, `font2`, `font4` and `font5` on **both** installs. Consequently, under
selector 1 the only byte in `0x20..0xFF` that selects an inkless record on those four atlases is
`0x20` itself — which the draw loop special-cases before it ever blits (`00457910 TEST BL,BL` /
`JZ 00457931`). The 32 bytes `0x00..0x1F` still read past the last record, which is `TEXT-INDEX-003`
measured per font rather than predicted: records 224..255 against a count of 224. `font3` is the
exception and is already published (`TEXT-FONT3-008`): 64 records, 11 with ink, and 0 of the image's
64 cells exist in it at all

**Confidence.** **High** for the record counts and the ink census — the counts are each file's own
trailer under the reader's own predicate (`SPR16A-RDR-017`), the walk is exact, and the census is a
decode of every record on both roots rather than a sample / **Medium** for the decoder: the `.16`
and `.16a` grammars are this repo's published claims (`SPR16A-FONT-013`, `SPR16A-RLE-002`) re-run
here, not re-derived, so an error in them would move these figures. "Ink" is counted as any pixel
the decoder set, which is a decode property, not a rendering

### TEXT-FONT2-016

`font2.16` is 13 892 bytes on EN and 9 436 on RU (sha256 `c262119b…cb0c071b` against
`c5108777…ed66cdef`); `font1`, `font3`, `font4` and `font5` are byte-identical across the installs.
Both `font2` files declare 224 records and walk exactly. EN carries ink on **194** records, RU on
**159**; the 35-record difference lies **entirely** in character codes `0x80..0xAF`, and those are
precisely the records selector 1 can never select, because no byte converts into `0x80..0xAF`
(`TEXT-DOM-010`). So the RU release deleted exactly the art its own selector cannot reach, and lost
nothing. The consequence runs the other way and is a real trap for a consumer: under **selector 0**
the RU root's `font2` selects an inkless record for **65** byte values — `0x20` plus the whole of
`0x80..0xAF` and `0xE0..0xEF` — against 30 for the EN file. A build that ships RU data with the EN
selector, or that omits the selector entirely, loses text on that font and only on that font,
silently and with no fault

**Confidence.** **High** for the counts and for the localisation of the difference: two files, both
walked exactly, every record decoded on both roots, and the blanked set intersected against the
converter's reachable set by enumeration rather than by inspection / **Medium** for the selector-0
consequence, which is predicted from the index rule and the ink census and was not witnessed
running. This overturns nothing in `SPR16A-FONT-022`; it adds *which* records and *why* it cost the
RU build nothing

### TEXT-TILDE2-017

In the measurer `FUN_00456320` the test is `0045634c CMP AL,0x7e`, where `AL` was loaded three bytes
earlier by `00456349 MOV AL,[EDI+EBX]` — the string byte — and the `CALL 0x004562f0` does not run
until `0045635b`, **nine bytes later**, after the branch has already been taken. In the draw
`FUN_004577f0` the test is `004578ae CMP BL,0x5e` on `BL`, which is `00457897 MOV BL,AL` /
`00457899 SUB BL,0x20` applied to the result of `0045788e CALL 0x004562f0`. `TEXT-TILDE-009` says
both test the converted value "at the same point and in the same way"; they do not. The two
nonetheless agree on all 256 inputs, and the reason is `TEXT-DOM-010`: `conv` is the identity below
`0x80` and sends everything at or above it to `0xB0` or higher, so `conv(b)` equals `0x7e` if and
only if `b` equals `0x7e`. The rule `TEXT-TILDE-009` published is therefore correct and stays; what
changes is where a consumer must put the test. A reimplementation following that row's description
literally will convert before testing in the measurer — harmless today, and it stops being harmless
the moment the converter is touched, which is what `TEXT-CAP-018` is about

**Confidence.** **High.** Four instruction addresses, all inside routines this experiment dumped and
decoded from both installs, and the agreement is proved by the enumeration rather than assumed. The
rival — that the addresses are one routine read twice, or that a second `CMP` exists after the
`CALL` — is killed by the measurer's own listing, in which the only `CMP` against `0x7e` before
`00456383` is the one at `0045634c`

### TEXT-079

- Each loop step of `FUN_004577f0` recodes the byte at `i` and the byte at `i + 1`
  (`FUN_004562f0`) and subtracts 0x20 from both, so `~` becomes 0x5e. The tilde arm
  (`004578ae`..`00457907`) runs when byte `i` is 0x5e and byte `i + 1` is not (`004578b7 CMP
  AL,BL` / `JZ`). When both are 0x5e the ordinary glyph path draws one `~` glyph and the loop
  skips the partner (`00457950`..`00457956`): the doubled case of `TEXT-TILDE-009`.
- The arm calls the line routine `FUN_0044fcd0(x1, y1, x2, y2, colour)`, which plots both
  endpoints, with
  - x1 = x, the pen position;
  - x2 = x + `dat[j]`, where `j` is the recoded index of byte `i + 1` (the stack slot written
    at `004578b1` and read at `004578da`) and `dat` is the font's advance table, without the
    letter spacing;
  - y1 = y2 = y + `vt+0x24(0)`, the row below frame 0's height, which is 15 in font1
    (`TEXT-078`);
  - colour = the 16-bit word at `ramp + 0x1e` (`004578c8`): entry 15 of the fifth argument,
    the top level of the ramp that inks the text.
- The arm then jumps past the advance block (`00457907 JMP 0045794b`), so the pen stays at x
  and the next step draws the next glyph there. The line is `dat[j] + 1` pixels long and lies
  under that glyph.
- Through the wrapper `FUN_00456b50` both calls draw the line, the first `d` pixels right and
  down in the shadow ramp's entry 15 (`TEXT-078`).
- A tilde in the last byte of a string reads its partner from the terminating NUL, index 0xe0
  after the recode, which is one dword past the end of font1's 224-entry table
  (`evidence/metrics.tsv`). The line length then depends on what follows the table in memory.
- The font4 draw `FUN_00457b50` (`TEXT-067`) has no tilde arm. Its loop recodes each byte,
  subtracts 0x20 and either draws the glyph or advances for a space, so it draws `~` as glyph
  0x5e and advances by that entry, while the measurer skips a lone `~` (`TEXT-TILDE2-017`). A
  font4 string that holds a tilde is measured narrower than it is drawn. No font4 string was
  searched for one.
- Census: the five dialogue families of both roots contain no `~` byte (`DLG-LINE-038`), and
  the button label, string 77, contains none (`DLG-BUTTON-039`). No other text was searched.

**Confidence.** High for the arm's operands and the doubled case: `FUN_004577f0` and
`FUN_0044fcd0` are read whole and each operand is traced to its stack slot, which settles the
extent as the next glyph's advance and not the tilde glyph's own entry, 0x5e. High for the
missing arm in `FUN_00457b50`, read whole, which holds no compare against 0x5e.

**Unknown.** The read past the end of the advance table for a final tilde, what lies after the
table, and whether any shipped string outside the searched population reaches the arm.

## Draw routines

| ID | Claim | Confidence | Status | Evidence |
|---|---|---|---|---|
| TEXT-065 | `DrawText`'s font4 (`.16a`-class) callee unconditionally null-dereferences the argument `DrawText` passes as a literal zero, so `DrawText` cannot draw font4 text; only font1-3 (`.16`-class) pass through it. | High | ● active | [EXP-0399](../experiments/EXP-0399-text-presentation/EXP-0399.md) |
| TEXT-066 | Font4's constructor builds its `.16a` sprite and its own shading table once, at font-load time — not per `DrawText` call. | High | ● active | [EXP-0399](../experiments/EXP-0399-text-presentation/EXP-0399.md) |
| TEXT-067 | The second draw routine's vtable slot is a bare no-op for font1-3 (`.16`-class) and the published alpha compositor for font4 (`.16a`-class), which makes it the confirmed font4 draw path. | High | ● active | [EXP-0399](../experiments/EXP-0399-text-presentation/EXP-0399.md) |
| TEXT-068 | Neither text draw routine issues an outline/shadow/second-offset pass around a glyph; the tilde markup case is a substitute pixel-producing call, not an additional one. | High | ● active | [EXP-0399](../experiments/EXP-0399-text-presentation/EXP-0399.md) |
| TEXT-069 | No surviving claim in the corpus asserts anti-aliasing, supersampling or sub-pixel filtering anywhere in the sprite/text pixel path, and neither text draw routine's own body contains filtering logic. | Medium | ● active | [EXP-0399](../experiments/EXP-0399-text-presentation/EXP-0399.md) |
| TEXT-070 | Colour reaches a font4 pixel through a table-pointer selector, not a level shift: the `SHL EDX,0x9` level addressing belongs to `FUN_00428fc0`, not to font4's real receiver `FUN_0042b970`. | High / Medium / Unknown | ● active | [EXP-0399](../experiments/EXP-0399-text-presentation/EXP-0399.md) |
| TEXT-071 | `FUN_00428fc0` unconditionally dereferences the argument `DrawText` supplies as a literal zero, on both branches, before any pixel work: a null-pointer-plus-8 fault, not a benign argument miscount. | High / Unknown | ● active | [EXP-0399](../experiments/EXP-0399-text-presentation/EXP-0399.md) |
| TEXT-078 | The font's two draw routines read flag values 1, 2, 4 and 8 as anchors: 1 and 2 move x left by the string's width and by half of it, 4 and 8 move y up by frame 0's width, not its height, and by half of it. | High / Unknown | ● active | [EXP-0410](../experiments/EXP-0410-dialogue-layout/) |

### TEXT-065

`DrawText`'s glyph-draw call is a polymorphic vtable call, but its font4 (`.16a`-class) callee
unconditionally null-dereferences the exact argument `DrawText` passes as a literal zero, before any
pixel work — `DrawText` cannot be the routine that draws font4 text; only font1-3 (`.16`-class) can
pass through it without crashing.

`FUN_004577f0` at `0x0045792c` issues `CALL dword ptr [EAX+0x34]` after pushing exactly five stack
dwords, the first one pushed (which becomes the callee's own arg5) a literal `PUSH 0x0` at
`0x00457914`. For font1-3 (`.16`-class, vtable `0x005973e0`) this resolves to
`FUN_00428be0`→`FUN_0044d890`, the already-published opaque byte blitter (`SPR16A-FONT-013`,
`SPR16A-072`); `FUN_00428be0`'s full body never reads its own arg5 at all (only arg1/this, arg3 and
arg4 feed the six values it forwards to `FUN_0044d890`) — the constant `DrawText` pushes there is
simply inert for this class. Font4 alone is a `.16a`-class object (vtable `0x005974f8`, confirmed
directly: its constructor `FUN_00457990` calls `FUN_0042b920` at `0x00457a45`, which sets the
sprite's own vptr to `0x005974f8`), so the same instruction would instead resolve to `FUN_00428fc0`
— but `FUN_00428fc0` dereferences arg5 unconditionally on both of its branches (`MOV ECX,[ESP+0x18]`
then `MOV ESI,[ECX+0x8]` at `0x00428fd7`/`0x00428fde` and `0x00429005`/`0x0042900c`) before any
pixel-producing work, reading it as a shadeObj pointer (`TERR-LIGHT-059`'s own name for this
argument, in its ordinary, working caller). Because `DrawText` supplies arg5 as the constant `0`,
this dereference reads `[0x00000008]` on every call — an access violation under ordinary Win32
memory protection (the low 64KB is never mapped), unconditionally, on both branches, before a single
pixel is touched. `DrawText` therefore cannot be the routine that draws font4 glyphs in working
shipped code: the two draw routines split by sprite class, not merely by convention — `DrawText`
draws `.16`-class fonts (font1-3) only, and font4 is drawn, where it is drawn at all, through the
second draw routine instead (`FUN_00457b50`, `TEXT-067`)

**Confidence.** **High.** Both halves of the crash argument are direct instruction reads (the
constant push, and the unconditional dereference on both of `FUN_00428fc0`'s branches,
`evidence/d-text.txt`/`evidence/d-428fc0.txt`), and the consequence follows from ordinary Win32
memory protection, not a probabilistic model. This experiment did not observe a crash at runtime —
no game state was run — but a deterministic access violation on every call needs only the two reads
above to state, not a runtime witness. `TEXT-071` documents the argument-count asymmetry this same
reading also shows; this row states its practical consequence

### TEXT-066

`FUN_00457990`, the function `SPR16A-FONT-018` already names as font4's sole construction site,
calls `FUN_0042b920` (the `.16a` constructor) at `0x00457a45`, then at `0x00457b15` calls
`FUN_00428ad0` with `(count=0x10, mode=4, tint=0)` — `BuildShade(16,4,0)`, the exact signature
`PAL-MODE4-010`/`SPR16A-PIX-011` already read as the `.16a` shading-LUT construction arm (no
day/night tint reference on this arm). Font1-3's construction site (`FUN_00457640`) reaches neither
call; those fonts' sprite objects carry a NULL palette pointer (`SPR16A-FONT-013`)

**Confidence.** **High** for the instruction read: both calls and their literal pushed arguments are
quoted in `evidence/d-457b15-owner.txt`. The interpretation of mode 4/tint 0 as "the `.16a` arm" is
inherited from `PAL-MODE4-010`

### TEXT-067

The second draw routine's class split is now settled at both ends: font1-3 (`.16`-class) reach a
bare no-op through this routine's vtable slot — it draws literally nothing for those fonts — while
font4 (`.16a`-class) reaches the published alpha compositor, argument-count-consistent, and (per
`TEXT-065`) is this experiment's confirmed font4 draw path, not merely its best candidate.

`FUN_00457b50`'s glyph write is `CALL dword ptr [EDX+0x18]` (vt+0x18), not `+0x34`, after pushing
exactly five stack dwords (the same shape as `DrawText`'s own call). For the `.16` class this slot
is `FUN_0042baa0`, whose entire disassembled body (read for this correction,
`evidence/d-042baa0.txt`) is one instruction: `0042baa0 RET 0x14` — a bare return popping its five
stack args and doing nothing else; font1-3 draw no pixels at all through this routine. For the
`.16a` class the same slot is `FUN_0042b970`, which `SPR16A-078` (promoted) already reads selecting
`FUN_00451ae0` (forward) or `FUN_00451e50` (reversed) — the exact compositor `SPR16A-ALPHA-025`
names for the `.16a` alpha blend; its own body (already disassembled and committed by `EXP-0356`,
`experiments/EXP-0356-sprite-decoder-traces/evidence/receivers.txt`, read here, not regenerated)
ends `RET 0x14` on both branches, popping exactly the five dwords `FUN_00457b50` pushes. Given
`TEXT-065`'s finding that `DrawText`'s font4 dispatch faults unconditionally, this routine is no
longer merely argument-count-consistent among two candidates — it is the only one of the two draw
routines that can produce font4 pixels in working code

**Confidence.** **High** for both halves. The `.16` half is now a direct, complete disassembly of
the entire callee body (one instruction), not an absence/Unknown. The `.16a` half inherits
`SPR16A-078`'s own instruction read, this experiment's own read of `FUN_00457b50`'s dispatch
instruction, and the `RET 0x14` cross-check from `EXP-0356`'s already-committed evidence

### TEXT-068

Full bodies of `FUN_004577f0` and `FUN_00457b50` (`evidence/d-text.txt`) show exactly one
pixel-producing call per ordinary per-character loop iteration, one destination coordinate pair per
call. The tilde branch is a different call in the same slot, not a second call alongside it: for a
single `~`, the draw routine's per-character glyph call is replaced by one call to `FUN_0044fcd0` (a
Bresenham line routine, already published — `TEXT-TILDE-009`), drawing a horizontal rule instead of
a glyph, and the advance step is skipped for that character — still exactly one pixel-producing call
per character, never two. A caller wanting a drop-shadow effect needs two separate `DrawText` calls
at different colours/offsets; neither draw routine provides a second pass on top of its own single
call, ordinary or tilde

**Confidence.** **High** — a full read of both routines' complete bodies settles a structural
absence rather than bounding a search, and `TEXT-TILDE-009` already reads the tilde branch as
instead-of, not in-addition-to (it jumps past the advance block rather than falling through to it)

### TEXT-069

Searched: `supersampl` (0 rows), `anti-alias` (one claim ID, listed twice — once in its ledger, once
in `retracted.md`, as every retracted row is), `sub-pixel` (the same claim). The one directly
on-point claim, `SPR16A-PIX-009` ("the low byte is only anti-alias/sub-pixel coverage"), is
retracted; its own overturn argues for the destination-blend model `TEXT-065` also relies on, not a
coverage/AA model. What can look like smoothing is (a) the pre-graded intensity levels already baked
into the `.16` atlas at author time (`SPR16A-FONT-013`), and (b) font4's destination blend
(`TEXT-065`), a general sprite-compositing mechanism reused for text, not a text-specific AA
algorithm

**Confidence.** **Medium** — the absence is bounded to the three search terms named and the two draw
routines' own bodies read in full, not a whole-image search for every possible filtering
implementation

### TEXT-070

Colour/shading reaches a font4 pixel through a table-pointer selector, not a level shift — the
`SHL EDX,0x9` level addressing this experiment previously attributed to font4's draw-time argument
belongs to `FUN_00428fc0`, which `TEXT-065` now shows is not font4's actual draw path; font4's real
receiver (`FUN_0042b970`) uses a different mechanism entirely.

Font1-3: unchanged — a caller-selected 16-entry ramp from a small fixed set (thirteen ramps,
`SPR16A-FONT-013`), built by `FUN_00457c40` — whose only static caller image-wide is `FUN_0044c920`
(`evidence/callers.txt`), not either draw routine or the font constructor, so ramp construction is a
subsystem separate from drawing. Font4, drawn via `FUN_00457b50`→`FUN_0042b970` (`TEXT-067`): its
own fourth argument is a table-pointer override — `TEST EDX,EDX` / `JZ` / `MOV ECX,EDX` takes the
caller-supplied pointer if it is non-zero, otherwise (the `JZ` target) `MOV ECX,[ECX+0x1c]` falls
back to a pointer stored on the sprite object itself, `this+0x1c`. There is no `SHL` or comparable
shift anywhere in `FUN_0042b970`'s body — the level-shaped, bit-shifted table-row addressing this
experiment previously described belongs to `FUN_00428fc0`, a function `TEXT-065` now shows
`DrawText` cannot safely reach for a font4 object, and no known caller in this experiment's evidence
reaches it with a font object either (its established caller is `TERR-LIGHT-059`'s ordinary
lit-sprite path, unrelated to text) — so that level-shift description does not describe font4's
actual draw-time colour mechanism. What `this+0x1c` holds is not established here: it is a distinct
field from the `this+0x14` shading table `FUN_00428ad0`/`TEXT-066` builds (`LEA EDI,[ESI+0x14]` in
`FUN_00428ad0`, not `+0x1c`) — whether font4's default draw-time table is the `BuildShade` output or
a different field is Unknown. No claim in the corpus describes a distinct low-resolution text/UI
plane later stretched, and neither draw routine nor `FUN_00428be0`/`FUN_0042b970` computes a second
scale factor for x/y before handing coordinates onward — text composes directly into the
framebuffer's own channel-mask format, dispatch-layer stride only

**Confidence.** **High** for the ramp-vs-table-pointer split, the sole-caller census, and
`FUN_0042b970`'s own argument-4 mechanism (each a direct, complete disassembly read). Explicit
**Unknown**, newly scoped: what `this+0x1c` holds and whether it derives from `BuildShade`'s
`this+0x14` output. **Medium** retained for "no separate scaling stage" (the pixel writers
themselves were only partly read)

### TEXT-071

The argument-count mismatch this experiment first found between `DrawText`'s font4 dispatch and its
resolved callee is not a benign miscount: `FUN_00428fc0` unconditionally dereferences the exact
argument `DrawText` supplies as a literal zero, on both of its branches, before any pixel work — an
unconditional null-pointer-plus-8 fault, not merely a stack-argument-count anomaly.

`DrawText` (`FUN_004577f0`) pushes exactly five stack dwords before `CALL dword ptr [EAX+0x34]` at
`0x0045792c`, the first pushed (its own arg5) a literal `0` (`0x00457914`). For font1-3 the resolved
callee `FUN_00428be0` never reads arg5 at all — harmless. For font4 the resolved callee
`FUN_00428fc0` reads arg5 as a shadeObj pointer and dereferences `arg5+8` unconditionally on both
branches (`0x00428fd7`/`0x00428fde`, `0x00429005`/`0x0042900c`) before any pixel-producing work —
with arg5 forced to `0` by `DrawText`, this reads `[0x00000008]`, an access violation under ordinary
Win32 memory protection, on every call, both branches. The second draw routine `FUN_00457b50` also
pushes exactly five stack dwords before `CALL dword ptr [EDX+0x18]` at `0x00457c01`; for font4 the
resolved callee `FUN_0042b970` (already disassembled and committed by `EXP-0356`) ends `RET 0x14` on
both branches, popping exactly the five pushed, and its own fourth argument is tested for zero
before use (`TEST EDX,EDX`/`JZ`) rather than dereferenced unconditionally like `FUN_00428fc0`'s arg5
— safe against the same kind of caller-supplied zero. So `DrawText`'s font4 dispatch is not merely
argument-count-inconsistent, it is a deterministic fault, while the second draw routine's font4
dispatch is both argument-count-consistent and unconditional-dereference-free. This experiment now
states positively that `DrawText` draws `.16`-class fonts only and `FUN_00457b50` draws the
`.16a`-class font (font4); the open question narrows from "which of the two routines draws font4" to
"which callers/screens invoke `FUN_00457b50` with font4 active" — no caller-side census of
`FUN_00457b50` was run

**Confidence.** **High** for the argument-5 dereference and its consequence, and for the resolved
routine split (`DrawText` draws `.16` only, `FUN_00457b50` draws `.16a` only) — each rests on direct
instruction reads cross-checked against Win32 memory-protection semantics, not a probabilistic
count. Explicit **Unknown**, narrowed: which callers/screens invoke `FUN_00457b50` with font4 active
— that needs a caller-side census this experiment did not run

### TEXT-078

- `FUN_004577f0`, the draw slot `vt+0x14` (`TEXT-API-007`), takes x, y, the string, a flag word
  and a ramp (`RET 0x14`). It tests the flag byte at `004577f7`, `0045781b`, `0045782e` and
  `00457843`, one value each, and reads no other bit. The second draw `FUN_00457b50`
  (`TEXT-067`) repeats the four tests at `00457b57`, `00457b79`, `00457b8a` and `00457b9d` with
  the same operations.
- The tests are independent, so the effects add.

| Flag value | Effect on the origin |
|---|---|
| 1 | x = x - width of the string (`FUN_00456320`) |
| 2 | x = x - (width >> 1) |
| 4 | y = y - `vt+0x20(0)` of the glyph sprite |
| 8 | y = y - (`vt+0x20(0)` >> 1) |

- Slot `+0x20` of both sprite class tables, `0x005973e0` (font1-3) and `0x005974f8` (font4),
  read as raw dwords, is `FUN_00428b00`. It returns the first dword of frame 0's record, its
  width. Slot `+0x24` is `FUN_00428b10`, which returns the second dword, its height.
  `FUN_004577f0` uses the height getter for the advance of a space (half of it) and for the
  underline row (`TEXT-079`), so the two vertical anchors take the width where the height would
  be expected.
- Font1 frame 0 measures 16 wide and 15 high on both roots (`evidence/metrics.tsv`). Value 4
  subtracts 16 and value 8 subtracts 8; the height would give 15 and 7.
- Value 10, the anchor of the dialogue button label (`DLG-BUTTON-039`), centres the string on x
  by its measured width and puts the top of a font1 string's cell at y - 8, so the cell spans
  rows y - 8 to y + 6.
- `FUN_00456b50(x, y, string, anchor, ramp, d)` calls the font's draw slot twice with the same
  anchor: at (x + d, y + d) with the shadow ramp that the font's slot `+0x18` returns
  (`FUN_00457980`, `0x005e9be8`), then at (x, y) with `ramp`. Each call applies the anchors
  itself, so the shadow keeps the ink's alignment. `TEXT-068` stands: the offset pass is the
  caller's second call, not a pass inside the routine.

**Confidence.** High. Both routines are read whole, each test and both getters are named
instructions, and the class-table slots are raw dwords. The font1 frame 0 size is a decode of
`font1.16` on both roots.

**Unknown.** Which callers pass which flag values was not enumerated; the dialogue's own values,
0 for the text lines (`DLG-LINE-038`) and 10 for the button label, are the only ones established
here.

## Root corpora

| ID | Claim | Confidence | Status | Evidence |
|---|---|---|---|---|
| TEXT-ROOT-014 | On the text path the two installs are not variants of one corpus — they are two corpora, and the executable is not the difference. | High / Medium | ● active | [EXP-0103](../experiments/EXP-0103-ru-text-domain/) |
| TEXT-ITEMNAME-019 | Item names are a text corpus like any other, and `TEXT-ROOT-014`'s finding holds on them exactly: the two roots ship two item-name corpora and one key file. | High | ● active | **[EXP-0142](../experiments/EXP-0142-item-names/)** |
| TEXT-PATTERN-020 | The Russian root ships two nodes that look exactly like a name-composition table, and `rom.exe` names neither, so the composition model they support cannot be the live one. | High / Medium | ● active | **[EXP-0142](../experiments/EXP-0142-item-names/)** |
| TEXT-BATTLEROOT-063 | An owner-supplied pre-release root's battle-event corpus is confined to three missions; where it shares a scene, speakers differ only by same-length substitution, while EN and RU differ in all three ways. | High / Medium | ● active | [EXP-0393](../experiments/EXP-0393-inn-text-closure/EXP-0393.md); extends TEXT-ROOT-014 |

### TEXT-ROOT-014

`rom.exe` is byte-identical, sha256
`942e9b72610eeba2f3d74930ee47cbc4476b85a942c14348f874ec7b7d367d03` on both, verified here by hashing
both files rather than taken from `TEXT-LANG-002`. Over the union of nodes with a text extension —
**no heuristic in this measurement at all**, membership is the extension and the verdict is a
payload SHA-256 — there are **437** nodes: **422 differ in bytes, 12 exist only on RU, 1 only on EN,
and 2 are identical.** The two identical ones are `main.res/text/cutpaths.txt` and
`main.res/text/heropicture.txt`, both configuration lists rather than prose; the EN-only one is
`graphics.res/version.txt`. The 12 RU-only nodes are three `text/battle/mNNN/eventNN.txt`, seven
`text/inn/npc/npc*.txt`, `text/material.txt` and `text/pattern.txt`. Over **all** nodes: 3 151
identical, 898 differing, 18 EN-only, 14 RU-only. This is the text-path counterpart of `AI-ROOT-049`
and it is a stronger result in the same direction: where the map path differed in a few files, the
text path shares essentially nothing. A consumer cannot treat one root's text as a translation of
the other's, and cannot key a mission's text off a node list taken from either root alone

**Confidence.** **High** for every count: the node sets come from the container walk already closed
in `claims/res.md`, membership is by extension, and every verdict is a SHA-256 of the payload on
each root (`evidence/diff-nodes.csv` carries both hashes for all 932 rows) / **Medium** for "reaches
something a player sees": the RU-only nodes sit under `text/battle/` and `text/inn/npc/` beside
nodes that plainly do, but no reader was traced to any of them this round

### TEXT-ITEMNAME-019

`main.res:text/itemname.txt` is 7 703 bytes and 416 lines on the EN root and 9 582 bytes and 416
lines on the RU root, and **0 of the 416 lines are identical** while `text/itemname.bin`, the 416
`u16` keys that pair them, is **byte-identical** (sha256
`adb09ccc7a210ead12812901f847e299747138030cdd085addb1a1895a7c1e24`). The RU corpus carries 8 035
high bytes over its 416 lines and the EN corpus **0**; the longest RU line is 44 bytes against 29
EN. Every one of those high bytes therefore reaches the screen through `TEXT-CONV-001`'s converter
under selector 1, which `TEXT-LANG-002` derives from `main\id` and not from the executable, and
`TEXT-CAP-018`'s 160-of-224 ceiling applies to item names unchanged. The store the display actually
reads is `ITEM-DISPNAME-036`'s map; `world.res:data/data.bin`'s own item row names are **identical**
on the two roots and are not the localisation seam (`DAT-ITEMNAME-010`)

**Confidence.** High: two node payloads compared byte for byte on both roots and a line-by-line
comparison of a complete 416-line corpus; the onward clause about the converter is inherited from
the TEXT ledger rather than re-derived here

### TEXT-PATTERN-020

`MAIN.RES` carries `text/pattern.txt` (14 144 B) and `text/material.txt` (143 B) which the EN
`main.res` does not. `pattern.txt` is an `English=English` table whose left column is the
composed-name form (`Bronze Plate Bracers=Bronze Plate Bracers`) and `material.txt` is the 15-line
English material list with `Wood` written `Wooden`, which is the very substitution the Armor and
Shield constructors' `FUN_004dba5f` performs. Together they are a complete composition worksheet.
The refutation is a **raw byte scan of the whole image**, not an xref sweep: the byte strings
`pattern.txt` and `material.txt` occur **nowhere** in `rom.exe`, which is byte-identical on the two
roots, so no build of the shipped executable can open either. They are localisation tooling left in
the archive

**Confidence.** High for the absence, the instrument being a raw scan of the whole file for the
literal rather than any sweep that can be blind to orphan code, a vtable or a static initialiser /
Medium for `pattern.txt` being a worksheet rather than data for some other consumer: nothing was
found that reads it, which is not the same as establishing what it was for

### TEXT-BATTLEROOT-063

A fourth, owner-supplied pre-release data root's battle-event corpus is not a subset of EN's or RU's
shape — it is confined to three missions end to end — and where a scene is shared with that root its
speaker assignment always differs by same-length substitution, never by reorder or length change; EN
and RU disagree with each other in all three ways.

Reproducing `TEXT-ROOT-014` from an independent registry+resource walk: EN 225, RU 228, RU-only
exactly `m100/event09`, `m130/event07`, `m150/event10` (matching that row's published count, now
named). The pre-release root's 23 files sit entirely inside missions 41, 51 and 91; every other
mission's battle events (all 225/228 of them) are absent from it, and three of its own 23
(`m41/event10`, `m41/event11`, `m41/event12`) exist on neither preserved install. Of 231 total rows
in the three-root union, 18 show a different `<npc=..>` tag sequence between roots that both ship
the file. Split by which roots are being compared: the 5 rows where the pre-release root disagrees
with EN or RU are all same-length substitutions, never a reorder or a length change (e.g.
`m91/event05`: EN/RU `25,23,21`, pre-release `21,22,21`). The remaining 13 rows are EN-vs-RU
disagreements, where the pre-release root ships nothing to compare, and there all three kinds occur:
3 are pure reorders of the same speaker multiset (e.g. `m10/event02`: EN `51,51,21`, RU `21,51,51`),
5 are same-length substitutions (e.g. `m131/event10`: EN `62,62,62,21,62` vs. RU `74,74,74,21,74`),
and 5 are outright length changes — RU's own sequence has a different element count from EN's, not
merely a different order or value (e.g. `m101/event06`: EN 1 element `21`, RU 2 elements `44,21`;
`m90/event06`: EN 3 elements, RU 4; `m70/event02`: EN 12 elements, RU 14; `m100/event05`: EN 10, RU
9 — the sole case in either direction where EN is longer)

**Confidence.** High for the file-set counts and the RU-only three-file reproduction of
`TEXT-ROOT-014` (both are closed-population node censuses) / High for the
same-length-substitution-only pattern on the pre-release comparisons and for the
reorder/substitution/length-change split on the EN-vs-RU comparisons — this is the `<npc=..>` markup
itself (`DLG-MARKUP-007`'s own vocabulary), directly parsed, not inferred, and the classification
(reorder = same multiset; substitution = same length, different multiset; length change = different
element count) is a mechanical comparison, not a reading / **Medium** for what a substituted,
reordered or lengthened speaker sequence means in play: no scene was traced to a consumer, only its
markup was compared

## String tables

| ID | Claim | Confidence | Status | Evidence |
|---|---|---|---|---|
| TEXT-STRTAB-023 | Sixteen text files are loaded into one shared line-pointer array, in a fixed order, and that order *is* the index space. | High | ● active (partially retracted) | [EXP-0143](../experiments/EXP-0143-actor-names/), amended [EXP-0228](../experiments/EXP-0228-ui-text-indices/) |
| TEXT-NAMETAB-026 | Which of the sixteen tables localise, measured line by line on both roots — and the answer is *almost all of them*, with the exceptions naming themselves. | High / Medium | ● active (amended) | [EXP-0143](../experiments/EXP-0143-actor-names/), amended [EXP-0228](../experiments/EXP-0228-ui-text-indices/) |
| TEXT-STRTAB-030 | `FUN_00423220(this,id)` is a two-field array-index accessor that returns an address, not a Win32 `STRINGTABLE` lookup, and no STRINGTABLE resource is read on the tip popup's caption path. | High | ● active | [EXP-0198](../experiments/EXP-0198-tip-popup-presentation/) |
| TEXT-STRTAB-031 | A caller census of the index accessor corroborates `TEXT-STRTAB-023`'s one-shared-array reading; a 2-site gap between two scan instruments resolves to two further thiscall thunks on the same object. | Medium | ● active | [EXP-0198](../experiments/EXP-0198-tip-popup-presentation/) |
| TEXT-UI-032 | The sixteen table paths and bases are shared, but the complete line-array length is root-specific: EN 1,568, RU 1,527. | High | ● active | [EXP-0228](../experiments/EXP-0228-ui-text-indices/) |
| TEXT-UI-033 | The selected global and local string accessors have no unclassified static consumer or stored callback in the owned image. | Medium | ● active | [EXP-0228](../experiments/EXP-0228-ui-text-indices/) |

### TEXT-STRTAB-023

`FUN_00468550(this, node)` reads the node's whole payload into one heap buffer, then walks it: it
appends the current line's start pointer to a growing pointer array at **`[0x005eb3d4]`** (count
`[0x5eb3d8]`, capacity `[0x5eb3dc]`, grow-by `[0x5eb3e0]`, doubling through
`FUN_00572824`/`FUN_00572860`), scans forward to the next `CR`, overwrites that byte with a NUL **in
place**, and advances the cursor by **2**; the loop ends when the cursor passes the buffer end, and
the table object keeps how many lines it added at `+0x08` and the index it started at at `+0x0c`.
Nothing validates a line, nothing counts them in advance, and no code-page pass runs at load — the
bytes are used exactly as shipped. `EnumRefs callto:00468550` returns **16 hits / 2 owners / 0
orphan**: fifteen consecutive calls in `FUN_004709e0` at `004711a6`…`00471278` and one in
`FUN_00438670` at `004386fe`, each with its own `this` and its own path literal, read off the `PUSH`
of every site — in order, `main\text\` `main.txt`, `heropicture.txt`, `stats.txt`, `spells.txt`,
`spell.txt`, `dialogs.txt`, `unitname.txt`, `building.txt`, `itemname.txt`, `sites.txt`,
`npcnames.txt`, `cutscene.txt`, `cutpaths.txt`, `tunes.txt`, then `patch.res::patch.txt`, then
`main\text\credits.txt`. Two indexing styles coexist and a consumer must not confuse them:
`FUN_004687f0(table, i)` adds the table's own `+0x0c` base, so a table's lines are numbered **from 0
within that table**; a direct `MOV ECX,[0x005eb3d4]` / `MOV EDX,[ECX + 4·k]` uses the **global**
subscript, and `main.txt` is loaded first so its own line numbers are the global ones. Measured
separately on both roots, the bases run `main.txt` 0, `heropicture` 274, `stats` 300, `spells` 350,
`spell` 374, `dialogs` 402, `unitname` 568, `building` 649, `itemname` 715, `sites` 1131, `npcnames`
1150, `cutscene` 1245, `cutpaths` 1259, `tunes` 1273, `patch` 1294, `credits` 1361. ~~The sixteen
files hold 1,568 lines on both roots.~~ **EXP-0228 refutes that total alone:** `credits.txt` has 207
EN lines and 166 RU lines, so the arrays hold **1,568 EN / 1,527 RU**; every earlier table count and
base stands (`TEXT-UI-032`). Three limits a consumer inherits: a file **must** end every line with
`CR` and the byte after it is skipped whatever it is (the scan has no end check, so an unterminated
last line runs off the buffer); the index space is **positional**, so inserting a line renumbers
everything after it *in every later file*; and no reader bounds-checks a subscript

**Confidence.** **High.** The loader is one routine read whole, the call-site count names its
instrument with 0 orphan, every path literal is read off its own `PUSH`, and the load order fixes
each root's bases arithmetically from its own line counts. The old equal-total clause was High and
wrong; it is recorded in [`retracted.md`](retracted.md). `main.txt` at 0 remains a read of the call
order

**Amended.** EXP-0228 refutes the both-root total of 1,568 lines alone (`TEXT-UI-032`;
[`retracted.md`](retracted.md), REFUTED). The loader, the file order, the bases, the array structure
and the two index forms stand.

### TEXT-NAMETAB-026

Per file, lines that are byte-identical `en` to `ru` (`evidence/string-tables.txt`, corrected by
EXP-0228): `spells` 0/24, `spell` 0/28, `building` 0/66, `itemname` 0/416, `sites` 0/19, `tunes`
0/21 — **completely** translated; `stats` 4/50, `npcnames` 16/95, `dialogs` 19/166, `main` 17/274,
`unitname` 42/81, `patch` 46/67 — translated except for structural, empty and format-only lines; and
two that are **not localised at all**, `heropicture` 26/26 and `cutpaths` 14/14, both of which hold
resource paths rather than prose (`HERO-APPEAR-052` already reads the first as a path list).
`credits` has unequal lengths: **15 equal positions of 166 paired, plus 41 EN-only positions**
(`TEXT-UI-032`), replacing the ambiguous old `15/207` ratio. The discriminator is not the ratio but
*what* stays identical: in `unitname` the 42 identical lines are exactly the 42 empty ones and all
39 that carry a name differ (`UNIT-NAMETAB-041`), and in `main.txt` the 17 identical lines are seven
structural ones plus the ten hall-of-fame names (`FAME-DEFAULT-008`). So a table that mixes
translated and untranslated lines is doing so deliberately, per line, not by omission. G1
consequence, stated plainly: **every word the program chooses about a person or a thing is authored
in these files and moves with the root** — the one naming surface that does not is the executable's
own default-hero literals (`SESS-DEFNAME-032`)

**Confidence.** **High** for the census — sixteen files, both roots, complete, compared positionally
on the loader's own split with unequal tails reported separately / **Medium** for calling
`heropicture` and `cutpaths` "paths rather than prose", which is read off their contents and one
prior claim, not off a consumer

**Amended.** EXP-0228 corrected the `credits` count: 15 equal positions of 166 paired plus 41
EN-only positions (`TEXT-UI-032`) replace the old `15/207` ratio.

### TEXT-STRTAB-030

Read whole (`evidence/disasm-strtab-loader-423220.txt`, `RET 0x4`): `MOV ECX,[EAX+0x4]` loads the
object's own `+4` field (`this=[EBP-4]`), `MOV EDX,[EBP+8]` loads the sole stack argument `id`,
`LEA EAX,[ECX+EDX*4]` computes `+4field + id*4` and returns that address; the routine never
dereferences it, and every caller dereferences the returned address itself (`MOV EDX,[EAX]`). For
the fixed table object `this=0x005eb3d0` that the popup's two captions use, `+4` reads
`[0x005eb3d4]`, the address `TEXT-STRTAB-023` already publishes as the base of the one shared
line-pointer array that `FUN_00468550` fills from sixteen text files. `TOWN-185` (`EXP-0197`)
recorded that resolving the two captions would need the Win32 `STRINGTABLE` format decoded, and that
this experiment did not build it; no STRINGTABLE resource exists on this call path. The id is a
plain subscript into the array `TEXT-STRTAB-023` already documents, and `TOWN-206` (`EXP-0198`)
resolves both captions through it, on both roots

**Confidence.** **High.** The routine is four instructions, read in full, with an unambiguous
`LEA`-only (no-dereference) return shape and a `RET 0x4` matching its one stack argument. The alias
of `this+4` to `TEXT-STRTAB-023`'s own array base at `[0x005eb3d4]` is a direct address match on the
fixed global `this=0x005eb3d0`, not an inference

### TEXT-STRTAB-031

A caller census of the index accessor independently corroborates `TEXT-STRTAB-023`'s
one-shared-array reading, and a 2-site gap between two scan instruments resolves to two further
thiscall thunks on the same object, not a defect in either scan.

Ghidra's own reference index (`ScanField sym:FUN_00423220`,
`evidence/scan-423220-callers/SCAN_SUMMARY.md`) finds **21** call sites to `FUN_00423220` through
`this=0x005eb3d0`. A text-match scan for the immediate operand `0x5eb3d0` (`ScanField imm:5eb3d0`,
`evidence/scan-strtab-imm-5eb3d0/SCAN_SUMMARY.md`) finds **23** sites setting up that same `this`.
The 2-site difference is `00468090` and `004680b0`, each a five-byte `MOV ECX,0x5eb3d0` immediately
followed by an unconditional `JMP` (`evidence/disasm-strtab-candidate-loaders.txt`,
`evidence/disasm-strtab-thunk-4680b0.txt`), to `0x004691b0` and `0x004691d0` respectively, neither
of which is `FUN_00423220` (`0x00423220`). `00468090` is called directly as `FUN_00468090`;
`004680b0` is registered through `atexit`-shaped code at `004680a0` (`PUSH 0x004680b0` then
`CALL 0x00553f10`, a routine that converts its own argument call's return into `0` on success or
`-1` on failure). Both extra sites are thiscall setups for the same object reaching a method other
than the index accessor through a short forwarding thunk. This corroborates, from an instrument
independent of `TEXT-STRTAB-023`'s own `EnumRefs callto:00468550` load-site enumeration, that the
object at `0x005eb3d0` is a small, fixed-shape struct reached through a handful of named accessor
and lifecycle routines rather than reconstructed ad hoc at each call site

**Confidence.** **Medium.** Corpus agreement between two scans is capped at Medium by this ledger's
own confidence rule. The 2-site discrepancy was traced to two disassembled, byte-confirmed thunks
rather than left unexplained, which is evidence for that resolution; the underlying shared-array
reading is `TEXT-STRTAB-023`'s own High-confidence claim and is not re-graded here

### TEXT-UI-032

The first fifteen tables have equal counts and bases in both roots. `credits.txt`, loaded last at
base 1361, has 207 EN lines and 166 RU lines. Of the 166 paired positions, 15 are byte-identical and
151 differ; EN has 41 additional trailing positions. Every EXP-0228 requested consumer addresses an
earlier table, so the correction changes no selected index

**Confidence.** **High.** Both complete archives were parsed independently with the loader's CR
split; hashes and every row are preserved in `input-manifest.json`, `table-summary.tsv`, and
`strings.tsv`

### TEXT-UI-033

The repaired 8,005-function project finds global accessor `0x00423220`: 21 direct calls / 8 owners /
0 orphan; local accessor `0x004687f0`: 302 / 57 / 0; line-array address `0x005eb3d4`: 230 references
/ 47 owners / 0 orphan. The table object's 23 immediate references resolve to the 21 accessor calls
plus the lifecycle thunks at `0x00468090` and `0x004680b0`. A whole-file dword scan finds zero
stored pointers to either accessor, and neither occurs in an `.rdata` slot

**Confidence.** **Medium.** The named call, reference, orphan, pointer, callback, initializer, and
vtable representations are complete. A runtime-computed target or bulk copy that never stores either
accessor address remains outside the negative

## UI label consumers

| ID | Claim | Confidence | Status | Evidence |
|---|---|---|---|---|
| TEXT-UI-034 | The shared character/unit panel selects its displayed actor name as `unitname.txt[typeID]`. | High | ● active | [EXP-0228](../experiments/EXP-0228-ui-text-indices/) |
| TEXT-UI-035 | The same panel's persistent statistic captions come from fixed global `main.txt` slots, not `stats.txt`. | High | ● active | [EXP-0228](../experiments/EXP-0228-ui-text-indices/) |
| TEXT-UI-036 | The character/unit panel supplies its own numeric grammar and gates caption groups by visibility level. | High / Medium | ● active | [EXP-0228](../experiments/EXP-0228-ui-text-indices/) |
| TEXT-UI-037 | `stats.txt` has exactly three accessor calls, all in the item-description formatter, and none in the character/unit panel. | High | ● active | [EXP-0228](../experiments/EXP-0228-ui-text-indices/) |
| TEXT-UI-038 | The item-description formatter uses ten executable-authored format literals around its selected `stats.txt` labels and values. | High | ● active | [EXP-0228](../experiments/EXP-0228-ui-text-indices/) |
| TEXT-UI-039 | The shop hover consumer maps four stock rectangles to global slots 62..65 and the shopkeeper rectangle to slot 61. | High | ● active | [EXP-0228](../experiments/EXP-0228-ui-text-indices/) |
| TEXT-UI-040 | The shop identify modal draws slots 79 and 80 as two raw lines, so the `%d` in slot 80 is literal on this consumer. | High | ● active | [EXP-0228](../experiments/EXP-0228-ui-text-indices/) |
| TEXT-UI-041 | The shop control consumer chooses slots 72, 70, 71, and 73 for Undo, Buy, Sell, and Exit. | High | ● active | [EXP-0228](../experiments/EXP-0228-ui-text-indices/) |
| TEXT-UI-042 | The client opcode `0xbe` save-acknowledgement arm selects global slot 203 for a town notice. | High / Medium | ● active | [EXP-0228](../experiments/EXP-0228-ui-text-indices/) |
| TEXT-UI-043 | Mission outcome panels select global slots 140 and 141: success is `Mission Completed` / `Миссия выполнена`, failure is `Mission Failed` / `Миссия провалена`. | High | ● active | [EXP-0228](../experiments/EXP-0228-ui-text-indices/) |
| TEXT-UI-044 | The mission Esc menu initializes its visible controls from local `dialogs.txt` slots 34, 35 or 76, 36, 37, 38, 39, and 40. | High | ● active | [EXP-0228](../experiments/EXP-0228-ui-text-indices/) |
| TEXT-UI-045 | The town Esc menu uses local `dialogs.txt` slots 34, 35, 37, 77, and 40. | High | ● active | [EXP-0228](../experiments/EXP-0228-ui-text-indices/) |
| TEXT-UI-046 | The world map uses three localized index forms. | High / Medium | ● active | [EXP-0228](../experiments/EXP-0228-ui-text-indices/) |
| TEXT-UI-047 | The selection status uses two fixed lines for root-specific grammar. | High | ● active | [EXP-0228](../experiments/EXP-0228-ui-text-indices/) |

### TEXT-UI-034

In `FUN_00460480`, `[actor+0x20]` supplies the local index, object `0x005eb4c0` selects
`unitname.txt`, and `0x004605c3` calls the local accessor before the draw. All 81 positions are
stable across roots; every non-empty name is localized

**Confidence.** **High.** The index load, object identity, accessor call, and draw are one retained
instruction path, and both roots were extracted positionally

### TEXT-UI-035

Slots 15..18 are the primary statistics; 19..20 Health/Mana; 21..22 Sight/Speed; 23..26 damage,
absorption, attack, defence; 27..28 group headings; 30..34 weapon skills; 36..40 magic schools;
41..45 resistances; 35 Weight; 46 XP; 190..191 the conditional spellcaster and armour-piercing
flags. The Russian root supplies the same indices with Russian bytes

**Confidence.** **High.** `FUN_00460480` was retained through its terminator and every selected slot
is independently reproduced from both archives

### TEXT-UI-036

Primary statistics require level above 4; Sight/Speed above 1; attack/damage above 2;
defence/absorption above 3; skills/resistances above 6. Positive pools gate Health/Mana, and actor
kind selects the skill family. The consumer uses `%d`, `%d/%d`, `%d.%d`, and `#%s: %d-%d`; none is a
line in `stats.txt`

**Confidence.** **High** for the comparisons and format operands, read from one complete retained
consumer and byte-checked in `format-literals.tsv` / **Medium** for calling the controlling value
“visibility level”; its provenance was not decoded in this experiment

### TEXT-UI-037

Object `0x005ea230` reaches `FUN_004687f0` at `0x00484239`, `0x00484828`, and `0x00484869` inside
`FUN_00484160`; the local index comes from the item's encoded tag. The other local-accessor calls in
that function use different table objects

**Confidence.** **High.** The repaired-project immediate/object population and each call site's
receiver distinguish the three retained calls from same-function calls to other tables

### TEXT-UI-038

The forms are `#%s %+d`, `#%s %d`, `-%d`, `#%s: %d`, `#%s: %5.1f`, `#%s: +%d`, `#%s: +%d%%`,
`#%s %s%s`, ` %s %s%s`, and `#%s: %d-%d`. Thus the table supplies statistic names while code
supplies sign, punctuation, numeric precision, range, and suffix composition

**Confidence.** **High.** Every literal's address and raw bytes are in `format-literals.tsv`, each
push is inside the retained `FUN_00484160` range, and all 50 local labels are reproduced for EN and
RU

### TEXT-UI-039

The pairs are Armor/Броня, Weapons/Оружие, Magic items/Магические предметы, Scrolls, books &
potions/Свитки, книги и пузырьки, and Shopkeeper/Продавец. `FUN_004ab26b` computes `62+i` for the
four groups and uses fixed 61 for the shopkeeper

**Confidence.** **High.** Both call sites and their immediate/register arguments are in the complete
global-accessor population; the five slots are reproduced from both roots

### TEXT-UI-040

`FUN_004ab37e` looks up `Do you want to identify this item` and `for %d gold coins?` and draws each
centered; no format call or price operand occurs between lookup and draw. Both roots ship the same
English bytes for both positions

**Confidence.** **High** for this consumer's lookup-to-draw path and both-root bytes. No claim is
made that another, unlocated consumer cannot format the same line

### TEXT-UI-041

`FUN_004acb50` uses each slot twice while constructing the four controls; EN/RU pairs are
Undo/Отмена, Buy/Купить, Sell/Продать, Exit/Выход. Existing presentation evidence determines where a
separate or composed number is drawn; this row fixes the table selection

**Confidence.** **High.** All eight calls belong to one enumerated global-accessor owner and use
four fixed arguments reproduced from both roots

### TEXT-UI-042

At `0x00418510` it pushes 203, calls the global accessor at `0x0041851a`, and passes
`Your character is saved` / `Ваш персонаж сохранен` to the notice routine with numeric argument
`0x7530`. The argument's unit is not established here

**Confidence.** **High** for the opcode arm, index, bytes, and numeric operand, all in one retained
path / **Medium** for the human label “save acknowledgement”, based on the arm's surrounding
save-data handling rather than a named symbol

### TEXT-UI-043

The message-handler arms construct the two panels from those fixed positions; neither phrase is
mission-script dialogue

**Confidence.** **High.** Both construction arms and fixed array reads are retained, and both roots
reproduce the selected bytes

### TEXT-UI-044

The first branch offers Save Game and Load Game; the alternate mode offers Diplomacy in their place.
The remaining entries are Game Options, Sound Options, Quest Objectives, End Quest, and Return to
Game. Russian uses the same local positions, and the embedded tilde remains hotkey markup

**Confidence.** **High.** `FUN_0043c23b` was retained through all eight constructors and each
receiver is the dialogs table object; both roots preserve the index sequence

### TEXT-UI-045

They are Save Game, Load Game, Sound Options, Abort Game, and Return to Game, with localized Russian
bytes at the same positions. This is a distinct five-control constructor, not a condition on the
mission menu

**Confidence.** **High.** `FUN_0043c90b` supplies five fixed indices to the dialogs accessor in one
retained constructor; the extracted roots agree on positions

### TEXT-UI-046

A non-zero mission payment selects global slot 262 and formats `%s: %d c`; the return entry
concatenates local `sites.txt[0]` then global slot 261; hit-tested marker `i` returns `sites.txt[i]`
only after its availability test. The fixed pairs are Payment/Оплата, Plagat/Плагат, and Return to
town/Вернуться в город; the sites table has 19 localized positions in both roots

**Confidence.** **High** for the indices, format operands, order, gate, and corpus / **Medium** for
the visible spacing of the return composition, because this experiment did not decode the string
helper's separator behaviour

### TEXT-UI-047

Zero selected draws slots 47 and 48: EN `No units` + `selected`, RU `Персонаж` + `не выбран`. Two or
more draws slots 49 and 50 plus the count: EN `Units` + `selected:`, RU `Выбрано` + `персонажей:`.
Exactly one selected actor takes the full unit-panel path instead

**Confidence.** **High.** The count branches and four direct global-array reads are inside one
retained routine, and the line pairs are reproduced byte for byte from both roots

## Hover help

| ID | Claim | Confidence | Status | Evidence |
|---|---|---|---|---|
| TEXT-HOVER-048 | Hover help crosses a shared 500 ms cursor-idle threshold, rather than being requested by each screen on every paint. | High | ● active | [EXP-0361](../experiments/EXP-0361-hover-tooltips/) |
| TEXT-HOVERSET-049 | The inspected base-widget family has 96 static tables and 28 distinct hover getters. | High / Medium | ● active | [EXP-0361](../experiments/EXP-0361-hover-tooltips/) |
| TEXT-HOVERCHAR-050 | Character labels and generator icons return shipped explanatory prose. | High | ● active | [EXP-0361](../experiments/EXP-0361-hover-tooltips/) |
| TEXT-HOVERROOM-051 | The same delayed route serves room hotspots, navigation and inventory areas. | High | ● active | [EXP-0361](../experiments/EXP-0361-hover-tooltips/) |
| TEXT-HOVERTEXT-052 | Hover content is a mixture of installed prose and formatted live information, not a single ready-made string per icon. | High | ● active | [EXP-0361](../experiments/EXP-0361-hover-tooltips/) |
| TEXT-HOVERPAINT-053 | Normal hover help preserves source-authored line breaks and has its own layout, separate from introductory tips. | High | ● active | [EXP-0361](../experiments/EXP-0361-hover-tooltips/) |
| TEXT-080 | Spellbook hover getter `004b0350` joins a name and mana line with damage, range and duration lines and at most one caption of `main[182..187,217]`; the last caption with a value replaces earlier ones. | High | ● active | [EXP-0435](../experiments/EXP-0435-tooltip-sources/) |
| TEXT-081 | Spellbook values come from fold `004190f0` over the selected actors: per cell, sums or minima and maxima of spell-record fields; captions follow level formulas for ten spell ids; in the spellbook, ids 17, 27 and 28 never fill a caption. | High / Medium | ● active | [EXP-0435](../experiments/EXP-0435-tooltip-sources/) |
| TEXT-082 | Generator attribute hover `0042ca00` tests five rectangles per attribute row: `main[155+i]`, `main[273]`, `%s = %d`, then `%+d` of the next-point cost and of the refund, the last two comma-grouped. | High / Medium | ● active | [EXP-0435](../experiments/EXP-0435-tooltip-sources/) |
| TEXT-083 | Map-list hover `00444cd8` chooses by cursor x alone: the row's description left of 300 px, then `dialogs.txt[134]`, `[135]`, `[136]` per column; no map field selects 135 or 136. | High | ● active | [EXP-0435](../experiments/EXP-0435-tooltip-sources/) |
| TEXT-084 | The inherited hover hint at `+0x3c` has writers found by the traced pattern in three base constructors and setter `004bd615`; 252 vslot-6 call sites exist, of which 8 pass a `dialogs.txt` lookup. | High / Medium | ● active | [EXP-0435](../experiments/EXP-0435-tooltip-sources/) |
| TEXT-085 | The 172 hint-forwarding constructor call sites pass null (63), a static address (54) or a text-table entry (55, in `dialogs.txt` and `patch.txt`); 17 final classes inherit the getter. | High / Medium | ● active | [EXP-0435](../experiments/EXP-0435-tooltip-sources/) |

### TEXT-HOVER-048

`FUN_0042a810`, called on controller `005cd790` by `0047555b`, samples imported `timeGetTime` at
`0042a825`; `+58` accumulates elapsed time using `+5c`. `0042a95d..0042a975` requires signed new
elapsed >=500 and <25500, and unsigned old elapsed <500. With `+9c ==0`, `+a0 !=0` and a root at
frame `+cc`, `0042a9a1` calls `004bd564`, then `0042a9b1` calls the result's vslot `+14`, copies the
text and raises nonempty help. The normal pointer-change route `0042a5b0` clears elapsed/visibility
and refreshes the clock; the polling branch excludes a held right button. Elapsed >=25500 hides an
existing box. `0042b900` hides without clearing elapsed; frame mouse-button cases201..206 call it,
while `0042b0e0` enters a marquee and hides help. A failed crossing is not retried merely because
the pointer stays still

**Confidence.** **High for the identified static arithmetic and branches.** The
complete181-instruction poll has no unresolved jump and is checked by two decoders. Native polling
cadence, all focus paths, the outer rendering gate and arbitrary intervening field writers are not
claimed

### TEXT-HOVERSET-049

Aligned `.rdata`/`.data` tables matching vslots08/0c/10=`00401950/00401970/004bdcc6` and carrying30
code pointers yield21 specialized text-return bodies, six local zero-return bodies and inherited
getter `004bd631` in 69 tables. The inherited getter returns the owned CString at `+3c` unless
flag20 is set; vslot18 (`004bd615`) copies the supplied hint. `004bd564` searches children
recursively in reverse stored order, then the current widget's screen rectangle. `0042a810` asks
exactly that widget and has no parent-text retry after an empty answer

**Confidence.** **High for the measured signature population and local bodies; Medium as a UI
inventory.** Tables with different base methods, computed replacement tables and all
setter/constructor bindings for inherited hints were not exhaustively recovered. A null getter does
not prove screen-wide absence. Source-level names are not original symbols

### TEXT-HOVERCHAR-050

Shared card helper `004610b0`, reached from generator, tavern and mission character getters, selects
`main.txt[155..158]` for primary attributes,159/160 for health/mana,161/162 for
damage/attack,163/164 for armour/defence,165..170 for
weight/sight/speed/skills/resistance/experience,171..175 or 176..180 for weapon/magic skills and 188
for monster weapon resistance. Its disclosure level, actor kind, ownership and rectangle gates
remain part of the result. Generator getters `0042ca00`, `0042f120`, `004363b0` and `00432db0` also
bind free points273, skill pictures171..180, difficulty/template/continue/back247..255 and
name-entry256. Attribute +/- and value rectangles additionally format current values instead of
returning only prose

**Confidence.** **High for the reached getters, source indices and bounded control branches.** Both
original roots contain the same indexed populations. Native reachability in every transient state
and every inherited button hint remain unproved

### TEXT-HOVERROOM-051

`004b72c0` binds town mask1/4/2/8/10hex to main233/234/236/237/235. `004bba00` binds school
weapon/magic icons to 171..180 through helper `004bb0e0`'s visual-slot permutation1,2,4,3,5.
`004ab26b` binds shopkeeper61 and stock groups62..65. Shop shelf/backpack/table getters
`004a41b0/004a5280/004a5de0` and inventory bar `00483970` return scrolling/background labels54..60,
or call the item formatter for an occupied cell; inventory gold uses 74. Character getter `00490c60`
selects book/backpack/doll/menu labels8..14 and hero/portrait arrows52..53 or 121..122, then
card/item descriptions. Tavern getter `0047d880` likewise splits candidate card and worn-item mask
regions. Each getter retains its active, held-item, state-mask, selection and geometric conditions

**Confidence.** **High for local dispatch, text sources and source geometry.** The table does not
turn every visible room button into a hover target, and does not establish every full-screen
occlusion/focus state

### TEXT-HOVERTEXT-052

Spellbook getter `004b0350` reads names from `spells.txt` via object005eb3c0, formats labels
main117/118/123/124/182..187/217 with present live fields, and requires its availability bit. Card
helper `004610b0` builds a monster spell list from main192 plus `spell.txt` object005eb4b0.
Shop/equipment getters call the existing item formatter `00484160`; world-map getter `004654f0`
reads `sites.txt` only for a discovered or offered marker. Map-list getter `00444cd8` returns
row-owned text or `dialogs.txt[134..136]`. Paired RES measurements reproduce eight relevant tables
with equal EN/RU counts: main274, stats50, spells24, spell28, dialogs166, unitname81, itemname416,
sites19. The two independent readers agree on all1742 selected line lengths/hashes

**Confidence.** **High for the specific getters, table bindings and measured inputs.** Exact
complete item-formatter semantics remain with their existing claims; arbitrary runtime values, all
inherited hints and a claim that every line is localized are excluded

### TEXT-HOVERPAINT-053

`0042a700` splits the owned text on byte23hex (`#`), counts lines and measures their maximum width.
`0042b430` places a width=maxWidth+11, height=14*lines+5 rectangle above the stored cursor, shifting
left for right overflow and down at the top bound. `0042b510` draws split lines through `00456b50`
using font global005e92f8 (`font2`, TEXT-API-007), at insets5/4 and pitch14. No word-wrap pass
occurs in these three complete routines. Introductory tip constructor `004c7422` and its TipsMode
gate are a different route; the hover controller reads its own `+a0` field rather than that
preference

**Confidence.** **High for the complete normal layout routines and separate gates.** Native pixels,
all redraw/focus paths and oversized/corrupt authored lines remain unproved

### TEXT-080

Getter `004b0350` converts the cursor (`005cd7a0`, `005cd7a4`) with `004b0bf0` to a cell index, which is -1
when x < 0, y < 0, y >= 75 or x > 0x1c7. It reads the object at the singleton view's `+0xd0` and returns null
unless the cell is not negative and bit `cell` of that object's dword `+0x148` is set. Five CStrings are built; each `0056f3d4` call is
MFC-style `FormatV` (`0056f12e`: `GetBuffer`, `vsprintf`, `ReleaseBuffer(-1)`) and replaces its destination.

| Slot | Content | Template |
|---|---|---|
| A | `spells.txt[cell]` (object `005eb3c0`), `main[117]`, dword `+0x2cc+4*cell` | `%s#%s: %d` |
| B | `main[118]`, pair `+0x14c`/`+0x1ac` | `#%s: %d` or `#%s: %d-%d` |
| C | `main[123]`, pair `+0x20c`/`+0x26c` | same |
| D | `main[124]`, pair `+0x38c`/`+0x3ec`, each value times 0.0625 | `#%s: %5.1f` or `#%s: %5.1f-%5.1f` |
| E | one caption block, below | below |

Slots B..D are skipped when the pair's max is 0; an equal pair formats one value, otherwise the range. The
seven caption blocks run in code order 182, 183, 184, 185, 187, 186, 217 and all write slot E, so the last
block with a value wins. A block runs when its max exceeds -65535 (182, 184) or is nonzero (the rest). Each
caption's label is `main[N]` for its own number. The pairs and templates, equal value then range, are:

| Caption | Pair (min/max) | Templates |
|---|---|---|
| 182 | `+0x44c`/`+0x4ac` | `#%s: %d`, `#%s: %d...%d` |
| 183 | `+0x50c`/`+0x56c` | `#%s: +%d`, `#%s: +%d...+%d` |
| 184 | `+0x5cc`/`+0x62c` | `#%s: %d`, `#%s: %d...%d` |
| 185 | `+0x68c`/`+0x6ec` | `#%s: +%d%%`, `#%s: +%d...+%d%%` |
| 187 | `+0x74c`/`+0x7ac` | as 185 |
| 186 | `+0x80c`/`+0x86c` | `#%s: %d`, `#%s: %d-%d` |
| 217 | `+0x8cc`/`+0x92c` | as 186 |

The result is `%s%s%s%s%s` of A, B, C, D, E into static buffer `005f1cb8`, so the lines are `#`-separated
(TEXT-HOVERPAINT-053). The two roots run the same executable bytes (SHA-256 `942e9b72...7d03`).

**Confidence.** **High.** The routine `004b0350..004b08a4` is read whole: every Format call, its destination
frame slot, template address and operand displacement is in the evidence, and the instruction facts listed
there are checked against the decoded listing. The single-caption rule follows from all seven blocks writing
one slot before one join.

**Unknown.** Which object sits at the singleton's `+0xd0` at run time is bounded by TEXT-081; the cell
geometry behind `004b0bf0` beyond the bounds above was not decoded.

### TEXT-081

Fold `004190f0` (30 call sites) clears, on its `this`, the unit count (`+0x140`), flags (`+0x144`) and spell
mask (`+0x148`), then initialises per-cell arrays of 24 dwords: sums `+0x14c`, `+0x1ac` to 0; minima
`+0x20c`, `+0x2cc`, `+0x38c`, `+0x44c`, `+0x50c`, `+0x5cc`, `+0x68c`, `+0x74c`, `+0x80c`, `+0x8cc` to 0xffff;
maxima `+0x26c`, `+0x32c`, `+0x3ec`, `+0x56c`, `+0x6ec`, `+0x7ac`, `+0x86c`, `+0x92c` to 0; maxima `+0x4ac`
and `+0x62c` to 0xffff0001. It walks the selection map at `+0x9b8` when `+0x9c4` is nonzero.

An actor with a zero `+0x7c` is skipped silently. Otherwise the unit count `+0x140` is incremented and the
mask `+0x148` ORs the actor's `+0x18` (`004193e8..0041940f`). Then, if the singleton view's `+0x3dc` is not 1
and lacks bit 2, flags become 8 and only the per-cell fold is skipped. For such an actor the getter's mask test
can pass while the mana dword stays at 0xffff, so no value lines show. The actor's class string must equal
`CUnit` for the spell loop; its `+0x20` must be 0x17 or 0x18 and its `+0x18` nonzero. For each set bit `b` of
`+0x18` (0..23) the spell id is the byte at `005c232c+4*b`, and `004fdd96(id)` builds a record. The skill byte
is the actor's byte at `+0x14a` plus the index `[[record+4]+0xc]+8`, the stat byte is at `+0x139`,
`004fe1cf(skill, stat)` fills the record, and level `L` is skill plus stat minus 30, clamped to 0..100.

Per cell: record byte `+0xe` is summed into `+0x14c`, bytes `+0xe` and `+0xf` into `+0x1ac`; the word at
record `+0xc` is stored, not folded, into `+0x2cc` and `+0x32c`, so the last actor visited sets the mana value
shown; min and max of record byte `+0x9` and of word `+0x10` fill the 123 and 124 pairs. The seven accessors
take `L` and return a value only for their spell ids, else 0:

| Caption | Accessor | Spell ids | Value for level L |
|---|---|---|---|
| 182 | `004fe670` | 24; 7, 28 | L/15+1; -(L/15+1) |
| 183 | `004fe5a9` | 5, 16, 10, 22 | L/2 |
| 184 | `004fe62c` | 12, 17 | -(L/30+1) |
| 185 | `004fe541` | 23 | 4L/5+20 |
| 187 | `004fe575` | 27 | 4L/5+20 |
| 186 | `004fe4ed` | 14 | min(L/20+2, 7) |
| 217 | `004fe5fb` | 18 | L/10+3 |

The id table holds ids 1..10, 12..16 and 18..26 in 24 cells; ids 17, 27 and 28 are not in it. Caption 187
therefore cannot receive a value from this fold, and 182 and 184 receive none for ids 28 and 17. Cells
carrying captions: 182 cells 6 and 13; 183 cells 4, 7, 16, 19; 184 cell 11; 185 cell 5; 186 cell 9;
217 cell 23.

**Confidence.** **High** for the arrays, accessors, formulas, id table and cell mapping, read at instruction
level. **Medium** that the fold's `this` is the object the getter reads at the singleton's `+0xd0`: both use
the same field displacements, but the 30 callers were not all traced. Which spell each id names was not
re-derived.

**Unknown.** What `+0x20` values 0x17 and 0x18 and the `+0x3dc` bits name; the callers' firing times; whether
any state yields an actor counted and masked that fails the `+0x3dc` test; the meaning of record bytes `+0x9`,
`+0xe`, `+0xf` and word `+0x10` beyond their use here. A second consumer of the accessors exists: twelve
calls at `00484548..00484682`, in the item formatter family (TEXT-HOVERTEXT-052), call six of the seven
accessors, including `004fe575`, and format `main[182..187]`; what it shows for ids outside the book was not
read. Stores to `+0x74c` and `+0x7ac` occur only at `004192af`, `004192c0`, `00419aac` and `00419b0b`, all in
the fold (store-pattern scan of the code map).

### TEXT-082

Getter `0042ca00` returns null unless `[[view+0x5c]+0x104]` is nonzero. For attribute row `i` = 0..3 it
point-tests five rectangles, advancing a rectangle-table pointer by 0x30 and a second pointer by 0x10 per
row, and returns at the first hit; after row 3 it returns null.

| Rectangle | Returned text |
|---|---|
| 1 | `main[155+i]` |
| 2 | `main[273]` |
| 3 | `%s = %d`: the label CString `[view+0x1c0][i]` (set from `main[15..18]`), then the value dword `view+0x1d0+4*i` |
| 4 | `%+d` of -(T(v+1)-T(v)), v the same value dword, from `0042d050` |
| 5 | `%+d` of T(v)-T(v-1), from `0042d080` |

`T` is `004d3718` (HERO-COST-002). Both signed texts go to static `005e4068` and then through `00468f60`,
which keeps the string when its length is 3 or less, or 4 when the first character is `-` or `+`, and
otherwise loops, taking the right three characters and joining them to the rest with `,`. Rectangles 1, 2
and 3 return without grouping. The EN `main[155..158]` lines carry 2, 5, 3 and 4 `#` separators;
`main[15..18]` and `main[273]` carry none.

**Confidence.** **High** for the rectangle order, indices, templates, values read and the cost functions,
read in the complete getter and callees. **Medium** for the grouping output shape: the loop and constants
are read, but no output string was produced.

**Unknown.** The rectangles' pixel positions were not measured; the owner of `[view+0x5c]+0x104` was not
named.

### TEXT-083

Getter `00444cd8` reads the control's screen left `L` and cursor x `005cd7a0`. If x < L+0x12c it takes the row
under the cursor, `[control+0x84]` plus the row from `004c179d(y)`, and returns its record's CString at `+8`
(the description, below) when the vslot `+0x78` validity call accepts the row, else null. If x < L+0x186 it
returns `dialogs.txt[134]`; if x < L+0x1a4, `dialogs.txt[135]`; otherwise `[136]`. No record field enters
the test.

The row text `%s#%dx%d#%d#%d` (at `00445937`) takes the record's name (`+4`), `+0x14` minus 16, `+0x18` minus
16, `+0xc` and `+0x10`. The draw routine's column advances (0x1e plus 0x10e, 0x5a, 0x1e) equal the getter's
thresholds, so 134, 135 and 136 caption the size, `+0xc` and `+0x10` columns. The loader `00447211` fills
the record from a `.alm`: it compares the tag with `M7R`, then width and height go to `+0x14` and `+0x18`,
a 40-byte seek follows, then the 64-byte name (payload `+0x30`) to `+4`, two dwords (payload `+0x70`, `+0x74`)
to `+0xc` and `+0x10`, and a 512-byte block from payload `+0x78` whose line feeds become `#` before its leading
C string is stored at `+8`. The reads run only when the `M7R` compare matches and a header dword is at least
2 (`004472c5..004472cd`). The loader has one return-1 path, which tests `+0xc` > 1, and three return-0 paths;
`00444742` drops a map whose loader returns 0.

The EN `dialogs.txt[134..136]` lines are 11, 26 and 20 bytes, RU 12, 32 and 28, none with a `#`.

**Confidence.** **High.** The getter and loader are read whole, the thresholds equal the drawn column advances,
and both roots share the executable. The column-to-field pairing follows draw order and arguments, not a
label.

**Unknown.** What values
`+0x74` takes beyond ALM-META-026 is outside this read.

### TEXT-084

Three base constructors, `004bc770`, `004bc880` (three stack parameters) and `004bc98f` (six), construct the
CString at `+0x3c` and set flags `+0x18` to 1. The latter two call `CString::operator=(LPCSTR)` when their
text parameter (the third, respectively sixth) is nonzero; `004bc770` assigns the static `005f2120`.
Setter `004bd615` (vslot 6 in all 96 widget tables) assigns its argument to `+0x3c`; getter `004bd631`
returns it unless flag 0x20 is set (TEXT-HOVERSET-049). The instrument's pattern finds ten CString
construct, assign or copy calls on a `+0x3c` receiver: six in the three base constructors, one in the setter
and three in other functions (`00429f70`, `004df385`, `005168e0`) whose receiver class was not resolved.

The call `[reg+0x18]` occurs at 252 sites image-wide with unresolved receivers. Eight sit within one routine,
`00438e87`: four pairs, each passing `dialogs.txt[23]` and then `[75]` to vslot 6. No other site has a
text lookup within the 14 preceding instructions.

**Confidence.** **High** for the constructors, setter, getter and the listed sites, read from the code map
(11 788 entry points, 454 310 instructions). **Medium** for the absence of other writers: the 252 receivers
and three other functions were not classified.

**Unknown.** What the 244 other vslot-6 receivers are and what text they pass; whether a lookup stored in a
local before the call is missed by the 14-instruction window.

### TEXT-085

From the two base constructors with a text parameter, the parameter was traced backwards through 24
forwarding constructors to 172 call sites in 72 functions that call 23 distinct constructors, by the
argument's source:

| Source | Sites |
|---|---|
| null | 63 |
| address of an uninitialised data cell | 51 |
| address of a literal in initialised data | 3 |
| `dialogs.txt` lookup | 52 |
| `patch.txt` lookup | 3 |

The 51 data cells are single addresses at four-byte strides, 50 in `005e4090..005e4168` and one at `005f2138`;
no instruction other than the argument push references them, so they read as empty strings. The three
literals (`005c1bb4`, `005c1bcc`, `005c1be8`; 18, 21, 20 bytes, not read from a text resource) are passed
by `004c2353`. The 52 `dialogs.txt` sites use 51 distinct indices; the three `patch.txt` sites (object
`005ea668`, from `patch.res`, 67 lines) use 52, 53 and 54, all in `0043d847`. One site (`0044457e`) resolves
only by hand, to `dialogs.txt[117]`, because overlapping decodes hide its push. One site, `004450cc`, builds
the map list with `dialogs.txt[118]`, and its class `005988c8` has its own getter `00444cd8` (TEXT-083).
The other 171 sites construct 17 final classes, all carrying getter `004bd631`; 29 of the 171 sites construct
the base class itself with null. These are 17 of the 69 inherited-getter tables; the other 52 are not reached
by this trace.

**Confidence.** **High** for the counts under the stated trace: the code map, table reader and operand facts
are reproduced from both roots, the two roots share the executable, and the table lines exist in both.
**Medium** that the uninitialised cells hold empty strings, since a pointer-based writer is not excluded.

**Unknown.** The population of the 52 tables not reached; which constructed controls are ever shown; whether
the literal-text controls are reachable in play.

## Character-generation labels

| ID | Claim | Confidence | Status | Evidence |
|---|---|---|---|---|
| TEXT-CHARGEN-027 | The final character generator has a mixed resource-to-control contract: three text-overlay navigation buttons, one localised caption bitmap, and two five-entry pictorial skill banks. | High | ● active | [EXP-0157](../experiments/EXP-0157-chargen-labels/) |
| TEXT-CHARGEN-028 | The earlier precreation screen binds the name prompt directly, but sex and class are one four-way pictorial choice with no separate text labels. | High | ● active | [EXP-0157](../experiments/EXP-0157-chargen-labels/) |
| TEXT-CHARGEN-029 | Character-generation localisation selects content by root, not label index, and none of the requested static labels is executable-authored. | High | ● active | [EXP-0157](../experiments/EXP-0157-chargen-labels/) |

### TEXT-CHARGEN-027

`FUN_0042f680` constructs the statistic, card, navigation and skill children and initialises them at
`0042f91b..0042f943`. Navigation construction reads global string slots **238, 239 and 260** at
`0042d591..0042d5d0`; the measured EN/RU pairs are `Accept` / `Принять`, `Reset` / `Сбросить`, and
`Back` / `Назад`. Click indices 0/1/2 call continue, reset and back at `0042d96e`, `0042d95a`,
`0042d946`. The statistic panel instead loads `main.res::graphics\chrgen\leftup.bmp` at `0042bec9`,
blits it, then draws four formatted statistic values and the remaining-points number. The 160×238
bitmap itself carries EN `Body`, `Agility`, `Mind`, `Spirit`, `Free pool` or RU `Сила`, `Ловкость`,
`Разум`, `Дух`, `Свободные очки`; its hashes differ by root. Slots 15..18 are metadata and slots
155..158 plus 273 are hover prose, not the persistent caption draw. `FUN_0042e250` constructs
exactly five fighter images (`sword`, `axe`, `Mace`, `Pike`, `Bow`) or five mage images (`fire`,
`water`, `air`, `earth`, `astral`); `FUN_0042de70` draws images with no font call. Hover index `i`
alone returns slots **171..175** or **176..180** at `0042f161..0042f17e`

**Confidence.** **High for the static draw contract.** The executable binding and event paths
discriminate an overlaid-text model from a bitmap/picture model, while the complete both-root
resource comparison and decoded bitmap provide an independent content witness. No runtime frame is
claimed; transient tooltip reachability remains unwitnessed

### TEXT-CHARGEN-028

`00433c0f PUSH 0x464` / `00433c16 CALL FUN_004329b0` constructs the name-entry child and `00433c2a`
stores it at owner `+0x1bc`, closing `TEXT-NAMEIN-024`'s Medium surface identification. The draw
reads `[0x005eb3d4]+0x1f4`, global slot **125**, at `0043526d..00435288`: EN `Character name:`, RU
`Имя персонажа:`. The same child's tooltip getter returns slot **256**, EN
`Type character name here`, RU `Введите здесь имя персонажа`. `FUN_004340b0` loads four
`PreCreate\Heroes` image families (`mf`, `mm`, `ff`, `fm`) and `FUN_00434d70` draws their states
without a class- or sex-string read. Hit testing produces one value 0..3; acceptance passes the
selected field to `FUN_0047c590`, which masks `0xc0` and dispatches four character templates. The
earlier confirmation is a byte-identical `PreCreate\ButtonOk.bmp` with baked `OK`; a non-empty name
lets it emit continue message `0x445`

**Confidence.** **High for construction, draw and four-way dispatch.** A model with two labelled
selectors predicts two label sources or text draws and is refuted by the complete control
constructor and draw. The semantic gloss of what the `OK` confirms is not encoded and remains
interpretation

### TEXT-CHARGEN-029

Both roots use the same global string slots and resource paths in the same byte-identical
executable. Their active `main.res` supplies different bytes at every cited label/help slot and
different pixels in `graphics\chrgen\leftup.bmp`; the pictorial class/sex and skill resources and
`ButtonOk.bmp` are byte-identical. `FUN_00468810` opens `main\id`, reads its last character and
stores digit-minus-`0` at `[0x005eb57c]`: `english 0` → 0, `russian 1` → 1. That value gates byte
conversion (`TEXT-LANG-002`, with `TEXT-CONV-001`'s injectivity clause explicitly retracted by
`TEXT-DOM-010`); it does not branch any cited consumer to another index. Static wording is therefore
table bytes or bitmap pixels. Typed names and numbers are generated values, not labels. The seams
differ: wording is positional data; persistent statistic captions require editing a fixed 160×238
image; replacing one of five skill images is data-only, while adding a sixth choice, a third class
bank or separate sex/class controls changes the fixed consumer loops and executable behaviour

**Confidence.** **High for the source and selector mechanism.** Same-index/different-root bytes,
same-path/different-root pixels, identical executable and direct consumers discriminate the
bilingual-table and code-literal rivals. The customisation classification follows the named fixed
loops and files; no save or command format implication was examined

## Character-name entry

| ID | Claim | Confidence | Status | Evidence |
|---|---|---|---|---|
| TEXT-NAMEIN-024 | The typed-name field is a text-entry class with a hard ten-character cap applied *before* conversion — and `TEXT-IN-005`'s six owners are now all read. | High / Medium | ● active (partially retracted) | [EXP-0143](../experiments/EXP-0143-actor-names/), amended [EXP-0406](../experiments/EXP-0406-precreate-name/EXP-0406.md) |
| TEXT-COLL-025 | A typed name *can* reach a colliding byte at runtime: of 192 reachable glyph records, 48 collide, and 16 of those three ways. | High / Unknown | ● active | [EXP-0143](../experiments/EXP-0143-actor-names/) |
| TEXT-073 | Both traced openings of the pre-create screen load the name field from `main+0x430` after replacing `Unnamed` or `npcnames.txt` entries 20..23 with entry 20, so a new campaign without `-name` opens at EN `Danath`, RU `Данас`. | High / Unknown | ● active | [EXP-0406](../experiments/EXP-0406-precreate-name/EXP-0406.md) |
| TEXT-074 | A left-button press on hero slot 0, 1, 2 or 3 writes `npcnames.txt` entry 20, 21, 23 or 22 into the name field only when the slot differs from the last pressed one and the text is `Unnamed` or one of entries 20..23. | High / Medium / Unknown | ● active | [EXP-0406](../experiments/EXP-0406-precreate-name/EXP-0406.md) |
| TEXT-075 | The name field only appends: a converted byte at or above `0x20` goes to the end while the text is under 10 bytes, key 8 removes the last byte, and nothing selects or replaces, so the first accepted keystroke extends the seeded name. | High / Unknown | ● active | [EXP-0406](../experiments/EXP-0406-precreate-name/EXP-0406.md) |
| TEXT-076 | The caret is `\|` (`0x7c`) appended to the drawn copy of the name, in the name's font and colour; it flips at the first paint more than 500 ms after the last flip, character message under the cap or construction, with no focus test. | High / Unknown | ● active | [EXP-0406](../experiments/EXP-0406-precreate-name/EXP-0406.md) |
| TEXT-077 | The prompt `main.txt[125]` and the name are left-aligned font4 draws at origin + (224,305) and + (224,321); full-level inks are RGB(65,47,20) and RGB(101,39,61) on the normal ramp, (57,41,17) and (89,34,54) on the low-memory `/18` ramp. | High / Unknown | ● active | [EXP-0406](../experiments/EXP-0406-precreate-name/EXP-0406.md) |

### TEXT-NAMEIN-024

The class's vtable is ~~**`0x005979c0`**~~ **`0x005979a8`** (30 slots, `+0x00..+0x74`);
~~slot 23 (`+0x5c`)~~ slot 29 (`+0x74`) is
`FUN_00432c80`, eighteen instructions read whole: `00432c83 MOV EAX,[ESI + 0x60]` (the MFC
`CString`'s `m_pchData`), `00432c86 CMP dword ptr [EAX + -0x8],0xa` (its `nDataLength`), `JGE` →
**return 0, the keystroke is discarded**; otherwise `00432c91 CALL 0x00468a90` — the input converter
— and `00432c9c CALL 0x00432ac0`, which sets a caret-blink pair (`+0x70 = 1`,
`+0x74 =` ~~`GetTickCount`~~ `timeGetTime`), tests `CMP BL,0x20 / JC` so **every byte below `0x20` is dropped**, and
appends through `CString::operator+=` at `0x00572dcf`. `FUN_00432b00` is backspace — `Left(len-1)` —
dispatched from `FUN_00432c60` on key **8** and no other. `FUN_00432db0` on the same class returns
`stringTable[256]` when its owner's `+0x1e0` is non-zero, and index 256 of `TEXT-STRTAB-023`'s table
is the character-name hint. So the cap is on the **stored byte length**, and because the converter
is one byte in, one byte out, ten characters and ten bytes are the same limit in either language.
**The six owners of `FUN_00468a90`** (`EnumRefs callto:00468a90`: 6 hits, 6 owners, 0 orphan), each
read: `FUN_004bf1d3` and `FUN_004c3118` convert a keystroke only to compare it with `this+0x74`
(`FUN_004c3118` also accepts `0x0d`) — hotkey matchers that store nothing; `FUN_0044057c` converts
whole strings out of an array of **0x104**-byte records into a list control — a file browser, so one
owner really is a file read; `FUN_004bf8ed` appends into a width-measured, wrapping multi-line
buffer; `FUN_00437690` appends any byte `>= 0x20` with **no** length cap; and `FUN_00432c80` is the
capped field. That closes the Medium `TEXT-IN-005` published about itself

**Confidence.** **High** for the cap, the drop-below-`0x20` rule, the backspace key and the six
owners' shapes (every one is a cited instruction from a routine read end to end, and the enumeration
names its instrument with 0 orphan) / **Medium** that this class is the *character-name* field: the
tie is `FUN_00432db0`'s tooltip index plus the corpus fact that all ten shipped hall-of-fame names
are nine characters or fewer, which is consistent with a cap of ten but does not discriminate it
from any larger cap. A screen that instantiates the class was not traced

**Amended.** EXP-0406 (`TEXT-075`, `TEXT-076`; [`retracted.md`](retracted.md), PARTIALLY
RETRACTED) corrects three words of the mechanism. The 30-slot table starts at `0x005979a8`
(`EnumRefs vt:5979a8:30`); `0x005979c0` is its slot 6, so the entry `0x00597a1c` that holds
`FUN_00432c80` is slot 29 (`+0x74`), the character slot of `TOWN-211`'s message map. The caret
stamp calls `[0x00632fa4]`, which the image's import table binds to WINMM `timeGetTime`;
`GetTickCount` is `[0x00632cbc]`, which neither `FUN_00432a80` nor `FUN_00432ac0` calls. The cap,
the drop below `0x20`, the backspace key, the tooltip index and the six owners stand. The Medium
surface identification was closed by `TEXT-CHARGEN-028`'s `00433c16` construction; `TEXT-075`
reads that construction against this vtable.

### TEXT-COLL-025

Composing the input map (`TEXT-IN-005`) with the display map (`TEXT-DOM-010`) over what a keystroke
can store under selector 1: `0x00..0x7F` passes both untouched; `0x80..0xBF` is stored unchanged;
`0xC0..0xEF` is stored as `0x80..0xAF`; `0xF0..0xFF` is stored as `0xE0..0xEF`. Stored `0xC0..0xDF`
and `0xF0..0xFF` are therefore **unreachable** through the keyboard. At the draw, stored
`0x80..0xAF` becomes records 144..191, stored `0xB0..0xBF` becomes records 144..159, stored
`0xE0..0xEF` becomes records 208..223. Enumerated over all 256 typed values: **192** distinct
records are reachable and **48** of them collide. Records **160..191** take two typed preimages
each, `0x90..0xAF` and `0xD0..0xEF`; records **144..159** take **three**, `0x80..0x8F`, `0xB0..0xBF`
and `0xC0..0xCF` — the first and third colliding at the *input* converter and the first and second
at the *display* converter, and the third is the CP1251 block carrying the first sixteen Cyrillic
capitals. So 64 typed values are lost, the same 64 `TEXT-ALIAS-011` counts, but spread over 48
records rather than 64 pairs, because the input converter folds a second pair on top of the display
converter's. Two different names can be drawn identically, with no clamp and no substitute. The two
cross-language cases fall out of the same table: **a Latin name on the Russian install is safe**
(every byte below `0x80`, both maps identity), and **a Cyrillic name on the English install** takes
the non-1 selector where both maps are the identity, so a CP1251 byte `0xC0..0xFF` is stored raw and
drawn as record `0xA0..0xDF` — inside the 224-record atlas, so a glyph appears and nothing faults.
This closes the open item `docs/status/text.md` carried, positively

**Confidence.** **High** for the composition and the reachable set (both maps are already published
at High from routines read end to end, and the composition is an exhaustive enumeration over 256
inputs with no free parameter, taken through the one storage path `TEXT-NAMEIN-024` reads) /
**Unknown** which *glyph* each colliding record draws, which is `SPR16A-FONT-020`'s question, and
therefore whether a player would notice. The prediction that would settle it, and what refutes it,
is in the experiment

### TEXT-073

The pre-create screen is `main+0x360`, constructed at `0047269c` by `FUN_00432f90`, which stores
vtable `0x00597a20` (`00433124`) and runs the create routine `FUN_00433390` (`0043312a`). Its opener
`FUN_00479f50` has two callers (`EnumRefs callto:479f50`: 2 hits, 2 owners, 0 orphan): `00473b08`
in the `0x425` new-campaign arm of `FUN_00473110`, and `00475bce` in the result router
`FUN_004757b0`, reached from its detailed-generation arm (`00475999`) and by fall-through from its
`+0x35c` arm when that result is not `0x446` and `main+0x490` bit 2 is clear
(`00475ba2..00475bca`). The opener sets `pre+0x1c4` to a copy of the string at
`main+0x430` (`00479fa1..00479fca`), `pre+0x1d0` to `main+0x65c` (`00479fec`) and `pre+0x1cc` to
`main+0x494 >> 6` (`0047a001`), then calls `vt+0x80` (`0047a00f`), the enter routine
`FUN_00433c90`.

The enter clears both button groups, sets `pre+0x1cc = 0` (`00433d2d`), lights hero element 0
(`00433d3b`) and lights difficulty element `pre+0x1d0` (`00433d4d`). It compares `pre+0x1c4` byte
for byte with `npcnames.txt` entries 20, 21, 22 and 23 and with `Unnamed` (`0x5b9adc`); any match
replaces `pre+0x1c4` with entry 20 (`00433e80..00433e93`). The entries come from `FUN_004687f0` on
the table object `0x5ea658`, which `00471232..0047123c` loads from `main\text\npcnames.txt`; the
index counts CR-terminated lines from 0 (`TEXT-STRTAB-023`). The enter then sets
`pre+0x1d8 = pre+0x1cc` (`00433ea4`), copies `pre+0x1c4` into the name field through `FUN_00432b70`
(`00433ebd`) and sets `pre+0x1e0 = 1` (`00433fca`).

On a new campaign `main+0x430` is written by `FUN_00477450` (`00473aa4`) and again by
`FUN_0047b490` on `main+0x420` with three zero arguments (`00473aef..00473b01`; `EDI` is zeroed at
`FUN_00473110`'s entry, `0047314f`, in EXP-0208's committed export, and the path to the switch
does not write it). Both copy the bytes after `-name` in the string `SESS-HERO-013`
reads as the command line, at most 31 bytes, cut at the first space (`00477503..0047752a`,
`0047b5a6..0047b5c7`), or `Unnamed` when `-name` is absent (`00477530..00477556`,
`0047b5cc..0047b5ec`). Without `-name` the first opening therefore shows entry 20: EN `Danath`
(`44 61 6e 61 74 68`), RU `Данас` (`84 a0 ad a0 e1`). With `-nameX` it shows `X`; `-name X` stores
an empty string, which the enter keeps.

The detailed screen's Back reopens through `00475bce` with `main+0x430` holding the text the
pre-create forwarded (`00475a45..00475a6d`). The enter keeps a typed name, turns any of the five
default texts into entry 20, keeps the difficulty and resets the hero to slot 0, whatever
`pre+0x1cc` the opener copied.

`Master Oberic` (`0x5b9ab0`) has one reference, `00433b97` in `FUN_00433390`, which writes it into
`pre+0x1c4`. Both openers overwrite `pre+0x1c4` before the enter, so no traced path displays it.

**Confidence.** **High** for the two traced openings, the replacement rule, the seeding and the
per-root bytes: each step is a cited instruction in a routine read end to end, the caller census of
`FUN_00479f50` is complete for direct calls, and the entry bytes come from both roots'
`main.res::text/npcnames.txt`.

**Unknown.** What object `FUN_005847a4` returns: this experiment reads only its `+4`, `+0x70` use
and takes the command-line identity from `SESS-HERO-013`. Whether a path outside `FUN_00479f50`
dispatches the pre-create's `vt+0x80`: no census of `+0x80` virtual calls was run. The `0x425` arm
posts `0x42f` instead of opening when `main+0x6bc` is 3 (`00473ad0`); `FUN_00477450` stores 2 there
(`00477460`), and the three calls between (`00473ab1`, `00473abc`, `00473acb`) were not read. What
`main+0x430` holds when the `+0x35c` route opens the screen was not traced.

### TEXT-074

The pre-create's left-button slot `vt+0x54` is `FUN_00435750` (`TOWN-211` maps message `0x201` to
`+0x54`). It calls `FUN_004349b0(x, y, 1)` at `00435798`, which reads the region code under the
point from `Mask.bmp` (`FUN_004348e0`) and, for the four hero codes, lights that element and stores
its slot in `pre+0x1cc`: `0x50` → 0 (`00434a7b`), `0x64` → 1 (`00434abf`), `0x78` → 2 (`00434add`),
`0x8c` → 3 (`00434a99`). The click routine then switches on the same code (`004357c5`). Each hero
arm tests the last pressed slot `pre+0x1d8`, reads the field text (`FUN_00436e80`) and writes an
`npcnames.txt` entry through `FUN_00432b70` only when the text equals entry 20, 21, 22 or 23 or
`Unnamed`:

- `0x50`, slot 0, male fighter: skipped when `pre+0x1d8` is 0 (`004358e7`); writes entry 20
  (`00435a4e..00435a7b`); stores 0 (`00435a93`).
- `0x64`, slot 1, female fighter: skipped when it is 1 (`00435d1b`); writes entry 21
  (`00435e86..00435eb3`); stores 1 (`00435ec5`, `00435ed8`).
- `0x78`, slot 2, female mage: skipped when it is 2 (`00435ee7`); writes entry 23
  (`00436053..00436080`); stores 2 (`00436098`).
- `0x8c`, slot 3, male mage: skipped when it is 3 (`00435b0c`); writes entry 22
  (`00435c7a..00435ca3`); stores 3 (`00435cbd`).

The store runs whether or not the text was written. `EnumRefs disp:1d8` finds 33 hits in 20
owners; the hits in pre-create routines are these arms and the enter, which stores the entered
slot 0 there (`00433ea4`), so `pre+0x1d8` holds the last pressed slot. The other 18 owners were not
traced to a pre-create pointer. The field text has five writers (`EnumRefs callto:432b70`: 5 hits,
2 owners, 0 orphan): the enter and these four arms. A typed name survives
every press. The slot meanings are `SESS-HERO-013`'s. From the first opening's text the four
presses leave: slot 0 EN `Danath`, RU `Данас` (no write); slot 1 `Naira`, `Найра`; slot 2
`Reniesta`, `Рениеста`; slot 3 `Fergard`, `Фергард`. The probe reads each entry's bytes on both
roots.

The pointer-move slot `vt+0x4c`, `FUN_00435730`, calls `FUN_004349b0(x, y, flags & 1)` at
`00435742`. A move with bit 0 set over a portrait therefore stores its slot in `pre+0x1cc` without
entering a name arm and without updating `pre+0x1d8`.

**Confidence.** **High** for the four arms, their tests, the written entries and the writer census:
every branch is a cited instruction and the switch table is decoded from the image. **Medium** for
the pointer-move clause: it rests on `TOWN-211`'s `0x200` → `+0x4c` slot and on reading bit 0 of
the move's flags as the Win32 left-button bit, and no drag was observed.

**Unknown.** Whether any of the 18 `disp:1d8` owners outside the pre-create routines writes
`pre+0x1d8` through a pre-create pointer.

### TEXT-075

The name field is the class whose 30-slot vtable is `0x005979a8` (`EnumRefs vt:5979a8:30`: no slot
without a function). Three instructions store that address (`imm:5979a8`, `refto:5979a8`: 3 hits,
3 owners, 0 orphan): the constructors `FUN_00432930` (`00432969`) and `FUN_004329b0` (`00432a0a`)
and the destructor `FUN_00432a30` (`00432a48`). `FUN_00432930` has no caller. `FUN_004329b0` has
one, `00433c16` in the pre-create's create routine, which builds the field with id `0x464` and rect
(224,310)-(362,337) (`00433bfa..00433c0f`) and stores it at `pre+0x1bc` (`00433c2a`). The text is
the MFC `CString` at `+0x60`.

The class's own slots, each read whole:

- `+0x74`, `FUN_00432c80` (`TOWN-211`: message `0x102`): returns 0 when the length at `[+0x60]-8`
  is 10 or more (`00432c86`); otherwise converts the byte with `FUN_00468a90` (`TEXT-IN-005`) and
  calls `FUN_00432ac0`, which drops a byte below `0x20` (`00432ad8`) and appends any other with
  `CString::operator+=` (`00432ae9`).
- `+0x6c`, `FUN_00432c60` (message `0x100`): on key 8 calls `FUN_00432b00`, which replaces the text
  with its `Left(length - 1)` (`00432b1e..00432b3c`); every other key returns 0.
- `+0x54`, `FUN_00432c30`: when the point is inside the field, calls `FUN_004bd232` on the owner
  with the field (`00432c51`), which gives the field the owner's `+0x38` focus; the text is
  untouched.
- `+0x4c`, `FUN_00432b80`: sets the draw table `+0x64` to `+0x68` when the point is inside the
  field and to `+0x6c` when it is not; `FUN_00432a80` stores the same table in both.
- `+0x2c` paints (`TEXT-076`); `+0x14` returns tooltip slot 256 (`TEXT-NAMEIN-024`).

The inherited slots `+0x30`, `+0x50`, `+0x58..+0x68` and `+0x70` return without a store
(`FUN_00423850`, `FUN_00423860`, `FUN_00436e00..FUN_00436e30`, `FUN_00423c80`, `FUN_00423c90`).
The only other writer found is the setter `FUN_00432b70` (`TEXT-074`). No routine selects
text, moves an insertion point or replaces the text on a keystroke, so the first accepted keystroke
extends the seeded name. From each seed, typing stops at 10 bytes: EN `Danath` takes 4 more bytes,
RU `Данас` 5; `Naira` and `Найра` 5; `Fergard` and `Фергард` 3; `Reniesta` and `Рениеста` 2. The
cap limits typing only: the setter and the `-name` seed (`TEXT-073`) are not capped.

**Confidence.** **High** for the class identity, the construction census and the edit model: every
class slot is read end to end and every reference census names its instrument with 0 orphan.

**Unknown.** Whether a keystroke reaches the field before the field is pressed. By `TOWN-211` and
`TOWN-139`, a keyboard message goes to the owner's focused child `+0x38` and, unconsumed, to every
child's `vt+0x48`, which for this class ends in the slots above; this experiment did not re-read
that walk, what sets the owner's `+0x38` at opening, or how the frame delivers keyboard messages to
the pre-create. No `disp:` sweep of `+0x60` ran, so a write through a field pointer outside the
class is not excluded.

### TEXT-076

`FUN_00432cb0` is the field's paint slot `+0x2c` (`callto:432cb0`: the vtable entry `0x005979d4`
only). It reads `timeGetTime` (`00432cdb`, import slot `0x00632fa4`) and, when the caret flag
`+0x70` is set, draws a copy of the text with `|` (`0x7c`) appended by MFC's
`operator+(CString, char)` (`00432cfc..00432d16`); otherwise it draws the text alone
(`00432d2a..00432d32`). The draw is one font4 `vt+0x14` call with the field's table `+0x64`
(`00432d37..00432d68`, `TEXT-077`), so the caret is a glyph of the drawn string: it starts at the
text's advance sum and takes the text's colour. Font4 record 92 (`0x7c - 0x20`) inks columns 0..2
of rows 0..14, 45 pixels with 13 at level 15, and advances 4. The advance sums of the seeds are EN
`Danath` 54, `Naira` 42, `Fergard` 59, `Reniesta` 63 and RU `Данас` 47, `Найра` 50, `Фергард` 66,
`Рениеста` 72 pixels (`font4.dat` plus spacing 2, `SPR16A-FONT-018`).

After the draw, when `now - [+0x74]` exceeds 500 (`00432d6b..00432d78`, `CMP ECX,0x1f4` then
`JBE`), the paint stores `now` in `+0x74` and flips `+0x70` (`00432d7a..00432d83`). The pair is set
to (1, now) at construction by `FUN_00432a80` (`00432a97`, `00432aa4`; `callto:432a80`: the two
constructors) and by `FUN_00432ac0` before its `0x20` test (`00432ac4`, `00432acb`;
`callto:432ac0`: `00432c9c` only). Every character message under the cap, a dropped control byte
included, therefore shows the caret and restarts its phase; a character message at the cap, the key
8 handler and the setter do not. The paint reads no focus state: its fields are the owner's
`+0x08` and `+0x0c` and the field's `+0x08`, `+0x14`, `+0x60`, `+0x64`, `+0x70` and `+0x74`, so
the caret blinks whether or not the field holds focus.

**Confidence.** **High** for the caret glyph, its position and colour, the 500 ms rule and the
reset sites: the paint and both reset routines are read whole, and the census of each reset routine
is complete.

**Unknown.** How often the pre-create repaints. The flip happens only inside a paint, so each
visible phase lasts 500 ms plus the wait for the next paint.

### TEXT-077

The pre-create paint `FUN_00434d70` (`vt+0x2c`) draws the prompt at `0043524d..00435288`: string
`[[0x5eb3d4]+0x1f4]`, global slot 125 (`TEXT-CHARGEN-028`), at x = field `+0x08` + screen `+0x08`
and y = field `+0x0c` + screen `+0x0c` - 5, flags 0, table `[[0x5e9c78]+8]`, through
`[0x5e9bc0]`'s `vt+0x14`. `[0x5e9bc0]` is font4 (`TEXT-API-007`), and its vtable `0x00599158`
holds `FUN_00457b50` at `+0x14`. The field paint (`TEXT-076`) draws the name through the same
routine at x = field `+0x08` + screen `+0x08` and y = field `+0x14` + screen `+0x0c` - 16, where 16
is what the glyph sheet's `vt+0x24(0)` returns for font4's 16x16 record 0, with flags 0 and table
`+0x64` = `[[0x5e9c30]+8]` (`FUN_00432a80`). With the field rect (224,310)-(362,337) (`TEXT-075`)
the prompt's cells start at screen origin + (224,305) and the name's at + (224,321).
`FUN_00457b50` moves x only for flag values 1 (right) and 2 (centre) and y only for 4 and 8, so
flags 0 is left and top alignment and no text width enters either position. The screen origin is
`[0x5ea210]`, `[0x5ea214]` (`00472679..0047269c`, the pre-create's construction), which is
`((W-640)/2, (H-480)/2)` (`SHOP-VIEW-044`): (0,0) at 640x480 and (192,144) at 1024x768, the two
sizes `00471942..0047195f` selects.

Every inked literal word of `font4.16a` indexes palette entry 255 except four words at 254 in
record 6 (11,261 words and 4; 193 of 224 records carry ink), and the file is byte-identical on both
roots. The table pointer passed as the fourth argument selects the source table (`SPR16A-078`).
`FUN_00457c40` fills two 256-entry palettes and builds each ramp with
`FUN_00427df0(palette, 16, 4, 0)`: `0x5e8ef8` bytes 0, 1, 2 = (20, 47, 65)·i/255 (`0045852a`,
`0045850d`, `004584f0`) into `[0x5e9c78]` (`00458708`), and `0x5e8ac0` bytes 0, 1, 2 =
(61, 39, 101)·i/255 (`0045859f`, `00458578`, `00458542`) into `[0x5e9c30]` (`0045874e`). Mode 4
(`00428197`) scales every entry by level/16 for levels 1..16, or by level/18 when `[0x5eb570]` is
non-zero (`PAL-MODE4-010`), and packs the three bytes in the order `SPR16A-PAL-008` reads as
`[B,G,R,0]`. A level-15 word reads the table row built with level 16 (`SPR16A-ALPHA-025`). With the
`/16` table (`[0x5eb570]` = 0) that row is `pal · 16/16` and the normal-memory blit adds no
destination term, so the prompt's opaque ink is RGB(65,47,20) and the name's RGB(101,39,61) before
packing; lower levels blend these with the background. `[0x5eb570]`'s only writer sets it when
`GlobalMemoryStatus` reports less than 24,000,000 bytes of physical memory (`PAL-MODE4-010`). A
level-15 word then reads RGB(57,41,17) and RGB(89,34,54) from the `/18` table before packing, and
that branch's blit reads destination row `L` rather than `1+L` (`SPR16A-080`), so it also adds row
15, a sixteenth of the quantized background (`SPR16A-ALPHA-025`). The strings differ by root, EN
`Character name:` with 127 pixels of advance and RU `Имя персонажа:` with 132; placement, font and
colours do not, because the executable and both font4 files are byte-identical across the roots.

**Confidence.** **High** for placement, alignment, font, the two tables, the palette values and
both ramps' level-15 entries: each is a cited instruction or a decoded byte, and the channel order
follows `SPR16A-PAL-008` through the mode-4 arm that also serves file palettes.

**Unknown.** The pixel a display shows: packing uses the framebuffer's channel widths, which are
runtime state. Whether a native original runs with `[0x5eb570]` set, which `PAL-MODE4-010` leaves
Unknown.

## Save-label chooser

| ID | Claim | Confidence | Status | Evidence |
|---|---|---|---|---|
| TEXT-SAVELABEL-054 | The established byte-indexed code-page converter `FUN_004562f0` is not called by any of the nine SAVE/LOAD dialog functions, directly or through its only wrapper. | Medium | ● active | [EXP-0372](../experiments/EXP-0372-save-label-encoding/EXP-0372.md) |
| TEXT-SAVELABEL-055 | Withdrawn: the capped name class was said to be constructed nowhere by a literal vtable store; its vtable is `0x005979a8`, three instructions store it, and `TEXT-075` finds its one construction at `00433c16`. | — | ✖ retracted | [EXP-0372](../experiments/EXP-0372-save-label-encoding/EXP-0372.md), retracted by [EXP-0406](../experiments/EXP-0406-precreate-name/EXP-0406.md) |
| TEXT-SAVELABEL-057 | The SAVE label's copy primitive `FUN_00554700`/`FUN_00554771` (`SAV-LABELTAIL-236`'s NUL-bounded copy into application storage) filters no byte value other than `0x00` within the disassembled range read. | High | ● active | [EXP-0372](../experiments/EXP-0372-save-label-encoding/EXP-0372.md) |
| TEXT-SAVELABEL-058 | Neither catalogued text-rendering mechanism is reached by any of the nine SAVE/LOAD dialog functions, extending `TEXT-SAVELABEL-054` to the draw functions, the glyph blit and the GDI text-out imports. | High / Medium | ● active | [EXP-0376](../experiments/EXP-0376-save-label-draw/EXP-0376.md) |
| TEXT-SAVELABEL-059 | The four MFC `CDC` GDI-wrapper stubs sit at matching offsets in the `CDC`, `CClientDC`, `CWindowDC` and `CPaintDC` tables; `imm:`/`disp:`/`refto:` find no reference to `0x59d200`/`0x59d280`/`0x59d300`/`0x59d380`, slot 11 of each. | High / Unknown | ● active (partially retracted) | [EXP-0376](../experiments/EXP-0376-save-label-draw/EXP-0376.md), partially retracted by [EXP-0406](../experiments/EXP-0406-precreate-name/EXP-0406.md) |
| TEXT-SAVELABEL-060 | `FUN_004bd28f`, the child-vector-by-id accessor `TOWN-354` reads, sits immediately before the copy `SAV-SAVELABEL-1017` traced into `dialog+0x68`, and it is a system-wide utility, not SAVE/LOAD-specific. | High / Medium | ● active | [EXP-0376](../experiments/EXP-0376-save-label-draw/EXP-0376.md) |
| TEXT-SAVELABEL-061 | By elimination among the two catalogued text-rendering mechanisms and their sole GDI-wrapper path, the SAVE/LOAD chooser's label pixels are consistent with native list-control painting outside `rom.exe`'s code. | Medium / Unknown | ● active (amended) | [EXP-0376](../experiments/EXP-0376-save-label-draw/EXP-0376.md), amended [EXP-0406](../experiments/EXP-0406-precreate-name/EXP-0406.md) |

### TEXT-SAVELABEL-054

`callto:004562f0` finds 3 direct callers; through its only wrapper `FUN_00456320`
(`callto:00456320`), 14 more callers, 15 distinct owning functions total (`FUN_004577f0` and
`FUN_00457b50` call both the converter directly and the wrapper, so they are counted once, not
twice). None of the nine SAVE/LOAD dialog functions this experiment and `SAV-LABELTAIL-236` traced —
`0043a230`, `0043a2b7`, `0043aa85`, `0043ae76`, `0043af4e`, `0043ba59`, `00473110`, `00478af0`,
`00478c40` — is among the 15. `callto:` also scans every `.rdata` dword for a literal match to
either function's address; none is found, so vtable-dispatched entry to either is excluded

**Confidence.** **Medium.** A bounded, reproducible, two-levels-deep reverse-direct-call search
(direct callers of the converter, plus callers of its one known wrapper) over a named function
population, with vtable-dispatched entry ruled out by the same scan. It does not establish what
mechanism, if any, draws the label instead — no paint/`WM_DRAWITEM` call chain was traced, and a
caller reached only through a third level, or an indirect call through a pointer not stored in
`.rdata`, would not appear

### TEXT-SAVELABEL-055

- What stands is the search. No disassembled instruction references `0x005979c0` in the
  1,977,344-byte EN `rom.exe` (`imm`, `disp`, `refto` and `text` modes; EXP-0406 repeats `imm` and
  `refto` with 0 hits), and the same modes find the SAVE dialog's vtable `0x00597df8` at
  `0043aecb`.
- What fell is the address. `0x005979c0` is slot 6 of the class's 30-slot vtable, which starts at
  `0x005979a8`. That address is stored by the literal-vtable-store idiom in both constructors
  (`00432969`, `00432a0a`) and in the destructor (`00432a48`), and the constructor `FUN_004329b0`
  is called once, at `00433c16`, where the pre-create screen builds its name field (`TEXT-075`).
- The class is not the save-label producer for that reason instead: its one direct construction
  belongs to the pre-create screen (`callto:4329b0` 1 hit, `callto:432930` 0 hits).

**Confidence.** The High was earned for the absence of `0x005979c0` and stated about the class; it
stands only for that address. See [`retracted.md`](retracted.md).

**Amended.** Retracted as a whole by EXP-0406 (`TEXT-075`, [`retracted.md`](retracted.md)). The
withdrawn headline read: "`TEXT-NAMEIN-024`'s capped text-entry class is not constructed anywhere
in the image by the literal-vtable-store idiom every other read constructor in this codebase uses."
Its Medium clause, that a class built by another means would escape the search, was not the gap;
the searched address was.

### TEXT-SAVELABEL-057

The SAVE label's own copy primitive, `FUN_00554700`/`FUN_00554771` (the function `SAV-LABELTAIL-236`
already named as the NUL-bounded copy into application storage), filters no byte value other than
`0x00` within the disassembled range read.

The byte-wise fallback loop has exactly one value comparison, `TEST DL,DL` / `JZ` at
`0055477d..0055478d`, testing only for the terminator; every byte `0x01..0xFF` it copies is stored
unchanged. Its aligned-dword path uses a separate mechanism, the CRT's standard
`0x7efefeff`/`0x81010100` zero-byte detector plus four byte selectors
(`TEST AL,AL`/`TEST AH,AH`/`TEST EAX,0xff0000`/`TEST EAX,0xff000000`) — five value-bearing tests,
not one, each locating the terminating `0x00` rather than filtering byte content; this detector's
full shape is visible within the committed evidence in the adjacent routine at `0055472c..0055475b`,
syntactically identical CRT machinery not itself on this call's executed path. The committed
disassembly range (`00554700..005547a4`) stops before the primitive's own executed instance of this
detector (beginning `00554796`) completes; the function's own tail through its `RET` was not read
here. This is the opposite of `TEXT-NAMEIN-024`'s capped field, which drops every byte below `0x20`
and caps stored length at ten — the SAVE label path has neither rule in the range read. A live
corpus witness is `SAV-SAVELABEL-1018`

**Confidence.** **High** for the instruction read within the captured range: every test located is a
terminator/zero-byte-position test, none a value filter, and the byte-wise loop's single
`TEST DL,DL` is the only comparison gating what gets stored. The primitive's own tail beyond
`005547a4` was not read

### TEXT-SAVELABEL-058

Neither of this codebase's two catalogued text-rendering mechanisms is reached by any of the nine
SAVE/LOAD dialog functions, extending `TEXT-SAVELABEL-054`'s converter-only negative to the draw
functions, the glyph blit, and the whole GDI text-out import surface.

`callto:004577f0`/`callto:00457b50`/`callto:00428be0` (the font-atlas draw functions and glyph blit
already established for other text surfaces) find, in the whole 1,977,344-byte image, exactly one
reference each — the function's own vtable slot on its owning font object's table
(`0059914c`/`0059916c`/`00597414`), never a literal call site; the same `.rdata` scan that catches
vtable-slot stores finds no second reference. `refto:` (validated against a known-good control
before use: `callto:00632f50` finds 0 hits for `GetDC`, whose caller `FUN_00428870` is already named
in `EXP-0355`'s own committed evidence, showing `callto:` alone would silently miss every import
call — `refto:` correctly finds it) for the six GDI-text imports this image actually links
(`ExtTextOutA`, `TextOutA`, `TabbedTextOutA`, `DrawTextA`, `GetTextExtentPoint32A`,
`GetTextExtentPointA`) finds callers in seven functions system-wide — four MFC `CDC` wrapper stubs
(`00550cf2`/`00550d0e`/`00550d33`/`00550d6b`), `CDC::FillSolidRect` (`0058123f`, itself implemented
with `ExtTextOutA`, an MFC idiom not this codebase's own), and two further unidentified functions —
`FUN_00581269` (reached only from `FUN_005812e0`, four literal sites) and `FUN_0057b3ac` (the sole
caller of `GetTextExtentPointA`; its own callers were not searched by this experiment) — none of the
nine SAVE/LOAD dialog functions among them, and neither `FUN_005812e0`'s own callers nor
`FUN_0057b3ac`'s were searched, so an unsearched path from the nine reaching either through a
further level is unfound, not excluded. Within the nine functions, the only references to any of the
four known font-object globals (`font1`..`font4`, `0x5e88b8`/`0x5e92f8`/`0x5e9c28`/`0x5e9bc0`) are
eighteen occurrences (six in `FUN_0043a2b7`, seven in `FUN_0043af4e`, five in `FUN_00473110`) of one
repeated idiom — load the font1 pointer, then a compile-time-constant caption selector (a pushed
index `0x97`/`0x1d`/`0x90` into the caption-lookup accessor `FUN_004687f0`, or, for `FUN_00473110`,
a fixed `+0x34c` displacement off a second string-table global) into a control-construction helper —
a fixed-caption/warning control, not a per-item draw call; none reads a loop variable, list index,
or the label buffer. The remaining two of the nine dialog functions, `FUN_0043a230` and
`FUN_0043ae76`, carry no reference to any of the four font-object globals at all
(`evidence/disasm-ctors.txt`)

**Confidence.** **High** for the search itself (six modes, reproducible, whole-image scope).
**Medium** for "this rules out the two catalogued mechanisms": `callto:`/`refto:` cannot see a
virtual call dispatched through a font-object pointer without a literal reference to the callee's
own address (`CALL dword ptr [reg+0x14]`-shaped instructions were not enumerated system-wide), so an
uncatalogued caller reached only that way would not appear; the font-global check is complete for
the four *known* font objects only

### TEXT-SAVELABEL-059

The four MFC `CDC` GDI-wrapper stubs this image links (`00550cf2`/`00550d0e`/`00550d33`/`00550d6b`,
wrapping `TextOutA`/`ExtTextOutA`/`TabbedTextOutA`/`DrawTextA`) sit at matching offsets across four
near-identical vtables (`0x59d200`/`0x59d280`/`0x59d300`/`0x59d380`), identical over the 20 slots
read (`vt:` dumped 20 slots per base, spacing `0x80` between bases; 20 is what was read, not a
measured total table length)~~, and none of those four vtables is constructed anywhere in the
image — extending `TEXT-SAVELABEL-055`'s vtable-never-constructed pattern from the character-name
class to this GDI-wrapper family~~.

`callto:` for three of the four wrappers (`00550d0e`/`00550d33`/`00550d6b`) finds exactly four
`RDATA-SLOT` entries each (one per vtable, at matching offsets) and no literal caller; the fourth,
`00550cf2` (the `TextOutA` wrapper), was not reverse-searched by `callto:` in this experiment, but
the same four `vt:` dumps place it at slot `+0x38` in all four vtables. `imm:`/`disp:`/`refto:` for
all four ~~vtable base~~ addresses find 0 hits each across the whole image, the same three-mode
idiom `TEXT-SAVELABEL-055` used for its own control case. `refto:0058123f` (`CDC::FillSolidRect`, the one
function in this family with a real body) also finds 0 callers anywhere

**Confidence.** **High** for the search (three modes times four addresses, one reproducible
process, zero hits; the `vt:` slot pattern shown for all four wrappers, and `callto:` literal-caller
absence measured for three of the four; `00550cf2`'s own literal-caller count was not separately
measured). ~~**Medium** that this makes the family unreachable: the same caveat
`TEXT-SAVELABEL-055` already named applies unchanged — an immediate in undisassembled bytes, or a
class constructed some other way, would not appear~~

**Unknown.** Whether a constructed `CDC`, `CWindowDC` or `CPaintDC` object reaches one of the four
text wrappers on a save-label path.

**Amended.** EXP-0406 ([`retracted.md`](retracted.md), PARTIALLY RETRACTED) withdraws the
never-constructed clause and its Medium. The four searched addresses are slot 11 (`+0x2c`) of
tables that start `0x2c` lower, at `0x0059d1d4`, `0x0059d254`, `0x0059d2d4` and `0x0059d354`
(`vt:` 31 slots each). The dword before each start is an RTTI locator whose type descriptor names
`.?AVCDC@@`, `.?AVCClientDC@@`, `.?AVCWindowDC@@` and `.?AVCPaintDC@@`, and the four wrappers are
slots 25..28 (`+0x64..+0x70`) of every table. Each start is stored by its constructor (`0057bdc6`,
`0057cef0`, `0057cfa4`, `0057d056`) and its destructor (`0057bf1b`, `0057cf54`, `0057d008`,
`0057d0bc`). `callto:` finds six calls to the `CDC` constructor `0057bdc2`, three of them from the
other three constructors; three to the `CWindowDC` constructor `0057cf85`; two to the `CPaintDC`
constructor `0057d039`; and none to the `CClientDC` constructor `0057ced1`. The `CPaintDC` call at
`00472009` is in `FUN_00471fe0`. `callto:471fe0` finds only the `.rdata` dword `0x00599c1c`, the
last of six dwords at `0x00599c08` that read `0x0f`, 0, 0, 0, `0x0c`, `0x00471fe0`: the layout of
an MFC message-map entry binding message `0x0f` (`WM_PAINT`) to that routine.
`TEXT-SAVELABEL-055`, the pattern this card extended, is retracted. The zero-hit search stands for
the four interior addresses only.

### TEXT-SAVELABEL-060

`FUN_004bd28f` — the child-vector-by-id accessor `TOWN-354` already reads (walks the child vector at
`this+0x1c`, returns the child whose `+0x4` equals the passed id) — sits immediately before the copy
`SAV-SAVELABEL-1017` already traced into `dialog+0x68` at the SAVE dispatcher's rename-commit
handler. This experiment adds that the same accessor is a system-wide utility, not
SAVE/LOAD-dialog-specific.

`callto:004bd28f` finds 189 call sites across 54 distinct owning functions image-wide,
`FUN_0043ba59` among them (three sites: `0043bacc`/`0043baff`/`0043bb4f`, using tags `4`/`1`/`3`);
at the `0043baff` site the returned child's own vtable slot `+0x3c` is then called indirectly
(`CALL dword ptr [EDX+0x3c]`) before the already-published copy into `dialog+0x68` at `0043bb2c`
(that copy's fixed-capacity source is `SAV-LABELTAIL-236`'s 256-byte stack buffer, not re-measured
here). The returned child's concrete class/vtable identity — narrower than `TOWN-354`'s own scope,
which already places it in the dialog's own child vector — and the `+0x3c` method's own body were
not located in this experiment

**Confidence.** **High** for the accessor's body, established by `TOWN-354`. **Medium** for the
genericness count measured here (a named, reproducible, whole-image `callto:` search) and for "this
accessor feeds the already-published copy": the call ordering and tag values are a direct
instruction read, but the returned child's concrete class and the `+0x3c` method's own body were not
traced, so whether it reads back already-rendered widget state or something else remains open

### TEXT-SAVELABEL-061

By elimination among this codebase's two catalogued text-rendering mechanisms (`TEXT-SAVELABEL-058`)
and their sole GDI-wrapper access path (`TEXT-SAVELABEL-059`), the SAVE/LOAD chooser's per-item
label pixels are consistent with native Win32/MFC list-control default painting occurring entirely
outside `rom.exe`'s own code — an inference by elimination among named candidates, not a positive
trace of an external paint call, and it does not exclude an uncatalogued in-game draw routine
reached only by virtual dispatch this search cannot enumerate.

Because no draw primitive inside this image is shown reached, no byte-to-glyph table is shown
indexed for this field either, so `TEXT-SAVELABEL-054`'s open item narrows rather than closes: EN
and RU cannot differ in a mechanism this image does not exhibit for this field, because the two
lawful executables are the same file (independently re-verified in this experiment, not only cited
from `EXP-0372`) — any difference would have to come from outside `rom.exe`, unreachable by static
analysis. What the draw path does with a byte outside 7-bit printable ASCII, and what bounds a
label's drawn (as opposed to retrieved) length in a chooser row, are Unreachable by this experiment:
both require observing the native control's own paint step, which needs a running original

**Confidence.** **Medium** for the elimination and for "EN/RU cannot differ here" (both rest on
`TEXT-SAVELABEL-058`/`-059`'s own Medium bounds). **Unreachable**, not merely Unknown, for the
non-ASCII draw handling and the visible length bound: both need a running original, excluded by the
static-analysis-only boundary

**Amended.** EXP-0406 partially retracts `TEXT-SAVELABEL-059`: the `CDC`, `CWindowDC` and
`CPaintDC` classes whose tables hold the GDI text wrappers are constructed in the image, so the
wrappers are not unreachable by construction. The GDI-wrapper leg of this elimination rests on
`TEXT-SAVELABEL-058`'s traced-path negative alone: no traced SAVE/LOAD chooser path reaches the
wrappers or the text-out imports.

## Help text

| ID | Claim | Confidence | Status | Evidence |
|---|---|---|---|---|
| TEXT-086 | `main\text\help.txt` is read once at startup as one NUL-terminated string into the global string object at `0x005eb4d0`, outside the sixteen-file line table, with no code-page pass at load. | High | ● active | [EXP-0437](../experiments/EXP-0437-f1-help/) |
| TEXT-087 | Both roots' `help.txt` are CRLF-terminated paragraph lines with no other control byte and no NUL; EN is 965 bytes of 7-bit text in 35 pieces, RU 1313 bytes with 837 high bytes of 34 values in 39 pieces. | High | ● active | [EXP-0437](../experiments/EXP-0437-f1-help/) |
| TEXT-088 | Help is wrapped at 408 px by the dialogue splitter and wrapper, overflows the 12-line body on both roots, and is rewrapped at 382 px beside a scroll bar; it is not paged. | High / Medium / Unknown | ● active | [EXP-0437](../experiments/EXP-0437-f1-help/) |
| TEXT-089 | Each root's `help.txt` holds exactly one doubled tilde and no lone tilde, so help draws one literal `~` glyph and never reaches the underline arm; the OK label holds no tilde. | High | ● active | [EXP-0437](../experiments/EXP-0437-f1-help/) |

### TEXT-086

- The startup routine `FUN_004709e0` pushes `0x005eb4d0` and the literal at `0x005be2fc` (`main\text\help.txt`, between the `main\text\heropicture.txt` and `main\text\main.txt` literals) at `0047118a`/`0047118f` and calls `FUN_004682d0` (`00471194`). The sixteen-file table loads that follow (`MOV ECX,table` / `PUSH path` / `CALL 00468550`, `TEXT-STRTAB-023`) do not include it.
- `FUN_004682d0` opens the node (`FUN_004c9f10`), takes its size (`FUN_004ca130`), allocates size plus 1 (`FUN_00554390`), reads the whole payload (`FUN_004ca140`), stores a NUL at `buffer[size]` and copies the buffer into the string object (`FUN_00572b52`, `FUN_00572860`). Nothing splits the text into lines, and no `FUN_004562f0` call lies in the routine.
- `0x005eb4d0` has four references image-wide (`EnumRefs refto:5eb4d0`, 4 hits, 4 owners, 1 in orphan code): the load push, two `MOV ECX,0x5eb4d0` stubs at `00468110` and `00468130` (the object's destructor call), and the single read at `00473a52`, the help case of `MENU-051`. Help is the only consumer.
- The converter acts at draw time (`TEXT-CONV-001`), so the stored bytes are the file's bytes.
- Each root's `main.res` holds exactly one entry named `text/help.txt` (`tools/helptxt`).

**Confidence.** High: the load routine is read whole, the reference enumeration names its instrument and returns one reader, and the file population is a listing of both roots' archives.

### TEXT-087

- Instrument: `tools/helptxt` over each root's `main.res`, which prints counts, never text. Both files end with `CRLF`; bare CR 0, bare LF 0, NUL 0, other control bytes 0.
- EN: 965 bytes, 34 `CRLF`, 35 pieces on splitting at `CRLF`, 131 space bytes, no byte at or above `0x80`. The two empty pieces are the last two (a blank line before the final `CRLF`). Longest piece 60 bytes.
- RU: 1313 bytes, 38 `CRLF`, 39 pieces, 157 space bytes, 837 bytes at or above `0x80` taking 34 distinct values, and the first byte is above `0x80`. The two empty pieces are piece 1 (a blank line after the first line) and the last (the file ends after one `CRLF`). Longest piece 61 bytes.
- The files therefore share terminator, container and markup, and differ in language bytes, blank-line placement (EN trailing, RU after the first line) and line counts. All 837 RU high bytes lie in the converter's source set for the Russian selector, `0x80..0xAF` and `0xE0..0xEF` (`TEXT-DOM-010`, `TEXT-FIT2-013`): 0 occurrences fall outside it.
- `help.txt` carries no section, page or escape marker other than the doubled tilde of `TEXT-089`.

**Confidence.** High for the census: each figure is a direct count over the archive entry of each root, with the output committed. The meaning of a byte (which glyph a high byte draws on the Russian selector) is `TEXT-CONV-001`'s, not re-derived.

### TEXT-088

- The body control of `MENU-052` gives the text to `FUN_00456ab0`: the splitter `FUN_00456900` cuts at `CRLF` and trims left, the wrapper `FUN_004563e0` fits words to the width with the measure `FUN_00456320` over font 1's `.dat` advances plus spacing 2 (`DLG-LINE-038`, `TEXT-079`).
- Instrument: a transcription of that splitter, wrapper and measure in `tools/helptxt`, over each root's `font1.dat` and `help.txt`, no emulation. Wrap width 408 (body rectangle width). EN: 33 non-empty pieces wrap to 35 lines; RU: 37 pieces wrap to 47 lines. No wrapper stall on either.
- `FUN_004be44e` compares the line count times the line height with the body height 211 (`MENU-052`): 35 lines against 211 EN, 47 against 211 RU, a scroll bar on both. The width then drops by `0x1a` to 382 and the text is rewrapped: EN 35 lines (widest 376 px), RU 48 lines (widest 365 px). The comparison does not change with the line height read (15 or 17): 35 lines times 15 is 525.
- Arithmetic only: with 12 visible lines, lines minus visible would be 23 (EN) and 36 (RU). No routine read sets the scroll bar's range, so that formula is an assumption.
- Each paragraph's first line is indented 10 px and every line except a paragraph's last and the final line is justified (`DLG-LINE-038`), with a 1 px shadow. No key or marker pages the text.

**Confidence.** High for the pipeline and the scroll decision (routines read whole, the decision holds for any line height read). Medium for the exact line counts and the 12-line window: the wrapper is transcribed from its listing and checked against the dialogue rules, not emulated, and the body height is derived arithmetic (`MENU-052`).

**Unknown.** The scroll bar's range and so the scroll extent (the 23 and 36 above are unread-formula arithmetic), the initial scroll position, and native floating-point justification (`DLG-LINE-038`'s Unknown applies).

### TEXT-089

- Instrument: `tools/helptxt` byte census. EN `help.txt` holds two `~` bytes at offsets 103 and 104; RU holds two at offsets 95 and 96. Each pair is adjacent, so the file holds one doubled tilde and no lone tilde.
- By `TEXT-079` a doubled tilde draws one literal `~` glyph and the loop skips the partner, and the underline arm runs only for a lone tilde. Help never reaches the arm. The measurer counts a doubled tilde once, so the wrap widths of `TEXT-088` treat the pair as one glyph.
- The first line of `dialogs.txt`, the OK label, holds no tilde in either root (length 2 EN, 7 RU).
- This extends the tilde census of `TEXT-079`, which covered the five dialogue families and one button label: it adds `help.txt` and the OK label for both roots.

**Confidence.** High: a direct byte count over the file of each root. The glyph drawn follows `TEXT-079`'s High clause.

## Spellbook item card

| ID | Claim | Confidence | Status | Evidence |
|---|---|---|---|---|
| TEXT-090 | For an item whose stream carries kind 42, the formatter appends a spell fragment to the name line, ` of <spell>` EN and ` с заклинанием <spell>` RU, and adds no prefix to the buffer it composes. | High / Medium | ● active | [EXP-0438](../experiments/EXP-0438-resolver-regen-card/) |
| TEXT-091 | The spell name in that fragment is line `spell id - 1` of `main\text\spell.txt` (28 lines on both roots), inserted unchanged; EN and RU differ only in text-file content, since `rom.exe` is identical and holds the one format literal. | High / Medium | ● active | [EXP-0438](../experiments/EXP-0438-resolver-regen-card/) |
| TEXT-092 | The same formatter gives a non-book item a new `#` line `casts <spell>` for effect kind 41 and a `#Magic:` header from marker `0x33`; a book (class `0xe00`) takes the ` of <spell>` form for kind 41 and no header. | High / Medium | ● active | [EXP-0438](../experiments/EXP-0438-resolver-regen-card/) |

### TEXT-090

- `FUN_00484160` sets the buffer at `0x005f0358` to the empty string
  (`0048417e`). It looks the item's code word `item+0x06` up in the name map
  `0x005eb410` (`00484189`..`00484198`, `ITEM-DISPNAME-036`), copies the found
  name into the buffer (`004841bc`..`004841e9`), then loops over the attribute
  stream once per entry, the count being the byte at `item+0x09`
  (`004841eb`..`004841f3`, `004848f5`..`00484909`). Each pass reads a tag byte,
  forms `tag - 1`, indexes the byte table `0x00484954` (limit `0x32`,
  `00484210`) and jumps through the dword table `0x00484930` (`0048421e`).
- Tag 42 (`teachSpell`, `ITEM-EFFKEY-147`) is table entry 4, `004842e4`. The
  arm reads one stream byte, the spell id; loads `main[90]` and `main[91]`
  from the shared string array `[0x005eb3d4]` at `+0x168` and `+0x16c`; takes
  the spell name as line `id - 1` of the table object `0x005eb4b0` through
  `FUN_004687f0` (`004842fc`..`00484304`); and calls the formatter
  `FUN_005545d0` with the literal `0x005bef54`, ` %s %s%s`
  (`0048430d`, `00484385`). The arguments are `main[90]`, the spell name,
  `main[91]`. The result is appended to the end of the buffer
  (`004848cd`..`004848f3`). The literal has a leading space and no `#`, and the
  name line carries no text before the name from this arm, so the fragment
  continues the name line.
- Strings (`tools/spellcard`, both roots, 274 `main.txt` lines): EN `main[90]`
  is `of` and RU `main[90]` is `с заклинанием`; `main[91]` is empty on both. The
  five book codes `0xe13`..`0xe17` name `Book` in EN and `Книга` in RU. An
  item named Book that carries kind 42 therefore reads `Book of <spell>` in EN
  and `Книга с заклинанием <spell>` in RU.
- The formatter has four call sites: `0047d9c0`, `00483a5e`, `00491043` and
  `004a2cbd` (`EnumRefs callto:484160`, 0 orphan hits). The hover painter
  splits the buffer at `#` (`TEXT-HOVERPAINT-053`). The other three painters
  were not read.

**Confidence.** High for the composition: the arm, its pushes, the literal and
the append are cited instructions read whole, and the name copy precedes the
loop. Medium for what a player sees, because only the hover painter's
treatment of `#` is established and that painter was not traced for this card.

**Unknown.** Whether a non-hover painter inserts a title or a break around the
buffer; the three other callers were not read. Which shipped item rows carry
kind 41 or kind 42: no evidence file lists the stream of any row, so which kind
a shipped book uses, and so whether a shipped book card shows this fragment, is
open.

### TEXT-091

- The spell id byte is one-based: the arm decrements it (`004842fc DEC EAX`)
  before the table call. The table object `0x005eb4b0` is
  `main\text\spell.txt` (`0x005be28c`; `0x005be2a0` is `spells.txt`,
  `0x005be2b8` `stats.txt`, `0x005be2e8` `main.txt`; `TEXT-STRTAB-023`). The
  accessor `FUN_004687f0(table, i)` returns that table's line `i`
  (`UNIT-NAME-039`).
- Both roots carry 28 `spell.txt` lines; row lengths are 4..21 characters EN
  and 4..23 RU, one line per spell and no per-case column. The formatter has
  no second index, so the RU row is shown in the form it is stored in, after
  `с заклинанием`.
- `rom.exe` has one SHA-256 on both roots (`inputs.tsv`), so the format
  literals `0x005bef10`..`0x005bef64` are identical; the only EN/RU difference
  in the card is the content of the text files.

**Confidence.** High for the row selection and for the absence of a case-form
choice. Medium for the grammatical form of any one RU row: the rows were
counted and hashed, not read for inflection, and the card shows whatever form
a row holds.

### TEXT-092

- Tag 41 (`castSpell`) is table entry 3, `00484314`. It reads the spell id and
  tests the item class `item+0x06 & 0xf00` against `0xe00`
  (`00484328`..`00484338`). For class `0xe00` it takes the ` of` form of
  `TEXT-090` (`main[90]`, `main[91]`, literal `0x005bef54`;
  `0048433a`..`0048435c`). For any other class it loads `main[92]` and
  `main[93]` (`+0x170`, `+0x174`) and the literal `0x005bef48`, `#%s %s%s`
  (`0048435e`..`0048437b`), which starts a new line with the verb word. EN
  `main[92]` is `casts` and RU `main[92]` is `с заклинанием`; `main[93]` is
  empty on both.
- The marker tag `0x33` (51) is table entry 7, `004847c2`. For class `0xe00`
  it skips to the loop tail (`004847d7 JZ 0x004848f5`). For any other class it
  appends the literal `0x005bef18`, `#`, then `main[189]` (`+0x2f4`,
  `0048480f`): EN `Magic:`, RU `магия:`. The writer emits the marker first in
  every stream built by `FUN_00509363` (`ITEM-152`), so a non-book card
  carries `#Magic:` before its spell lines.
- The other literals the arms use are `0x005bef10` `#%s %+d`, `0x005bef40`
  `#%s: %d` and `0x005bef64` `#%s %d`; they are not part of this question.

**Confidence.** High for the arms and literals, read whole. Medium for the
on-screen result, for the reason in `TEXT-090`.

**Unknown.** Which other effects a Book row carries, and so whether a shipped
book card has lines beyond the fragment.

## Notice and pause lines

| ID | Claim | Confidence | Status | Evidence |
|---|---|---|---|---|
| TEXT-094 | `main.txt[119]` is the Pause panel's body text: one load of slot 119 exists among the 30 hits for displacement or immediate `0x1dc`, in the Pause arm; it is a modal panel, not a message-line line. | High / Medium | ● active | [EXP-0446](../experiments/EXP-0446-toggle-notice/) |
| TEXT-095 | Slots 94..107 and 218..220 are read by the seven setting arms and slots 108..116 by the speed arms; two other routines have scaled-index loads at those displacements, not traced to their base. | Medium | ● active | [EXP-0446](../experiments/EXP-0446-toggle-notice/) |

### TEXT-094

- Population: `EnumRefs disp:1dc` (25 hits, 20 owners) and `imm:1dc` (5 hits, 5 owners) over the EN `rom.exe`, 0 hits in orphan or undisassembled code, and `refto:5eb3d4` (230 hits, 47 owners, `TEXT-STRTAB-023`). A function can read slot 119 from the shared array only if it is one of those 47 owners; three of the 30 hit owners are (`FUN_0042bca0`, `FUN_00472b80`, `FUN_0047eb20`).
- `FUN_0042bca0` stores the constant `0x22` into `[ESI+0x1dc]` (`0042bd15`), a field write. `FUN_0047eb20` takes `LEA EBX,[ESI+0x1dc]` (`0047ec80`), the address of a field of its own object. `FUN_00472b80` loads `[[0x005eb3d4]+0x1dc]` at `00472c09`, inside the Pause arm: it allocates `0x78` bytes, calls `FUN_0043fbbf` with rectangle (0x20, 0x30, 0x260, 0x1b0) and shows the panel through `FUN_00476810`, the panel `MENU-052` describes. The arm requires `campaign+0x6bc == 2` and `campaign+0x3dc == 1` (`00472bd0`, `00472bdd`).
- No toggle or speed arm loads slot 119, and none draws on the Pause panel: the two use unrelated surfaces (`MENU-059`).

**Confidence.** High for the Pause load and for the surface (named instructions). Medium for "the only load": the owner rule excludes the other 27 hits but misses a pointer copied out of the array before use and a load through an index register with another displacement.

**Unknown.** Loads of slot 119 through a copy of the array pointer or a computed index.

### TEXT-095

- Indexed loads at the notice displacements, from `EnumRefs disp:` at 0x178..0x1d0 step 4 and 0x368..0x370, in owners of `refto:5eb3d4`: `FUN_0040f38e` at `0040f9d4`, `0040fa58`, `0040fabd`, `0040fb25`, `0040fb8a`, `0040fbd5` and `0040fc27` (the seven setting arms), and `FUN_00472b80` at `00472dbe` and `00472e26` (the speed arms).
- Four further `[base + index*4 + disp]` loads exist at those displacements: `004b044b` in `FUN_004b0350` (`0x1ac`) and `0042cd90`, `0042cdcc`, `0042cde5` in `FUN_0042ca00` (`0x1d0`). Those routines are the spellbook hover getter and the generator attribute getter, whose array loads `TEXT-080` and `TEXT-082` read at other displacements (`main[182..187,217]`, `main[155+i]`). The register that holds the base of these four loads was not traced.
- The remaining owner hits at these displacements are stores or fixed-field accesses (for example in `FUN_00434d70`, `FUN_0047eb20`, `FUN_0042bca0`); they were classified by operand shape, not each traced.

**Confidence.** Medium: the nine arm loads are named instructions; the exclusion of the other hits rests on the owner rule, operand shape and two earlier claims.

**Unknown.** The base register of the four indexed loads in `FUN_004b0350` and `FUN_0042ca00`, and any array read with a computed index outside the swept displacements.
