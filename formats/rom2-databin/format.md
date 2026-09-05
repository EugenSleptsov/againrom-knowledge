# ROM2 Data.bin database — identity survey

Level 3. Promoted, evidence-backed claims only. Basis: `R2-ASSET-004`
(`world.res`/`world_srv.res` partial match, now superseded on reach by
`R2-ASSET-029`/`030`/`031`, not retracted), `R2-ASSET-005` (`templates.bin`
distinctness, reconfirmed by `R2-ASSET-032`), `R2-ASSET-025` (`templates.bin`
string absence in `allods2.exe`), `R2-ASSET-026` (the `Data.bin` loading
subsystem in `allods2.exe`), `R2-ASSET-029` through `033` (full corpus-wide
grammar, delta location, title schema, `templates.bin` reader search,
`a2server.exe` classification). Cross-reference only, not evidence: ROM1's own
[`formats/databin/format.md`](../databin/format.md) (`DAT-GRAM-003`,
`DAT-LOC-001`, `DAT-SCHEMA-007`).

Seen as: `world.res:data/data.bin` (149 971 B), `world_srv.res:data/data.bin`
(152 527 B), and the ROM2 root's `templates.bin` (61 802 B) — three separate
files, none assumed in advance to be "the" ROM2 analog of ROM1's single
`world.res:data/data.bin` (88 327 B).

## Positive control

DAT-GRAM-003's wire grammar, replayed unmodified against ROM1's own
`world.res:data/data.bin` on both the EN and RU preserved roots, reaches exactly 0
residue at 88 327 bytes on both — the same result DAT-GRAM-003 already reports, and
the basis for trusting the two ROM2 results below rather than a transcription
error.

## `world.res` / `world_srv.res` result: the grammar, corpus-wide

DAT-GRAM-003's grammar, corrected in exactly one place — group C's raw byte block
widened from ROM1's 10 bytes to 14 — reaches **exact 0 residue** on both files: all
declared entry counts for all 8 groups (A–H) are fully consumed, on both the client
and the server file independently. No other group's shape changes. `R2-ASSET-029`.

A negative control rules out the correction being a general parsing fix rather than
a ROM2-specific fact: the same 14-byte grammar, replayed against ROM1's own EN/RU
`data.bin`, does not walk cleanly — it diverges at offset 88 325/88 327. Under the
unchanged 10-byte grammar, ROM1's own file walks to exact 0 residue at 88 327 (the
Positive control above); the 14-byte divergence is a new number the widened
grammar alone produces, not a reproduction of `R2-ASSET-004`'s own two divergence
offsets, which are ROM2's own file sizes and report nothing about ROM1's file.
Swept independently over every group-C raw width from 0 to 80 bytes, 14 is the
**only** width that reaches 0 residue on either ROM2 file, and 10 is the **only**
width that reaches 0 residue on either ROM1 file — the width is unique on both
sides, not merely a value that happens to close.

The widened group-C block is independently confirmed at instruction level in
`allods2.exe`: the class whose real constructor runs a 7×`u16` (14-byte) zero-init
loop also stamps the vtable whose own Serialize entry contains a literal `PUSH 0xe`
(14) around its raw-block read/write call — three signals, tied to one compiled
class by a directly-read vtable pointer, agreeing on 14. `R2-ASSET-029`.

### Per-group Ghidra identity

All eight groups' own compiled class is individually pinned by one call-graph
chain, read directly from `allods2.exe`: the Data.bin Serialize dispatcher's own
program order names each `SetSize` address's letter; each `SetSize` calls exactly
one ctor-loop; each ctor-loop calls exactly one real constructor and passes it the
class's own element byte size; each real constructor stamps exactly one vtable
literal; each vtable's own slot+8 entry is that class's Serialize body. No step
assumes letters follow address order.

| Group | SetSize | real ctor | element size | vtable | Serialize | ROM1 kind |
|---|---|---|---:|---|---|---|
| A | `00541150` | `00543550` | 0x68 (104) | `005dac78` | `00506fff` | `kMatShape` |
| B | `00541a80` | `00543710` | 0x1c (28) | `005dacb8` | `00506f55` | `kBase` |
| C | `00541460` | `005431e0` | 0x40 (64) | `005dac20` | `00506fa3` | `kArmor` (widened) |
| D | `00541770` | `005435c0` | 0x44 (68) | `005dac98` | `00507107` | `kMagicItem` |
| E | `00541d90` | `00543780` | 0x30 (48) | `005dacd8` | `00507486` | `kNStr2` |
| F | `005420a0` | `00543870` | 0x30 (48) | `005dacf8` | `00507d03` | `kNStr10` |
| G | `005423b0` | `00543960` | 0x1c (28) | `005dad18` | `00507e8e` | `kBase` (forwards) |
| H | `005426c0` | `005439d0` | 0x20 (32) | `005dad38` | `00508293` | `kSpell` |

B and G share the identical `kBase` behaviour (a name plus one member-delegated
param array, no group-specific fields) and the identical element size, so they are
not distinguished by shape — only by which construction chain reaches which
vtable. `FindVtables` slot+8 over `5dac00:5dae00` returns 18 candidate vtables, 9
of them 7-slot; one of the 9, `005dac40`, shares B's own Serialize address
(`00506f55`) but is reached by none of the eight construction chains above and is
excluded from the family on that ground. `R2-ASSET-029`.

### Per-group result

| Group | Collections | Declared entries (client / server) | Titles | Walk |
|---|---|---|---|---|
| A | Shapes, Materials | 7 / 7, 16 / 16 | 11 | clean |
| B | Magic | 50 / 50 | 30 | clean |
| C | Armors, Shields, Weapons | 31/10/28 (both files) | 18 | clean (14-byte raw block) |
| D | MagicItems | 97 / 97 | 4 | clean |
| E | Units | 242 / 242 | 64 | clean (span differs, see delta) |
| F | Humans | 290 / 299 | 28 | clean (count AND span differ) |
| G | Buildings | 180 / 180 | 9 | clean |
| H | Spells | 35 / 35 | 24 | clean |

"Declared entries" is the collection's own stored `u32` count (index 0 is a
reserved null on groups C–H per DAT-SCHEMA-004, matching ROM1); "Titles" counts the
group's own name-title plus its column titles.

### The 2556-byte client/server delta, located

`R2-ASSET-012` measured the whole-file size difference (152 527 − 149 971 = 2556
B). Per-group span accounting over the clean corpus-wide walk locates it exactly:
`Units` (+10 bytes, same declared count 242 on both files — longer per-entry
content) and `Humans` (+2546 bytes, 299 server rows against 290 client rows).
10 + 2546 = 2556, matching `R2-ASSET-012`'s total with no remainder. Every other
group's span is byte-identical between the two files. `R2-ASSET-030`.

## Column-title schema vs `DAT-SCHEMA-007`

Read via the clean corpus-wide walk and diffed column-for-column against ROM1's own
published list. Title text only; no destination/width/streamer-site data (that
cross-reference is `DAT-SCHEMA-007`'s own, not reproduced here). Counts below
exclude each group's own leading name-title (DAT-SCHEMA-004: slot `i` is title
`i+1`), so they read one lower than the "Titles" column in the per-group result
table above, which does count it. `R2-ASSET-031`.

| Group | Collection | ROM1 titles | ROM2 titles | Result |
|---|---|---|---|---|
| A | Shapes+Materials | 10 | 10 | identical |
| B | Magic | 29 | 29 | identical |
| C | Armors | 17 | 17 | identical |
| D | MagicItems | 3 | 3 | identical |
| E | Units | 56 | 63 | 44-title shared prefix, 1-title shared suffix; ROM1's 11-title `treasure.3`/`Power`/`Spell`/`Probability` block replaced by ROM2's 18, adding `serverID`, `knownSpells`, `skillFire`, `skillWater`, `skillAir`, `skillEarth`, `skillAstral` |
| F | Humans | 27 | 27 | identical |
| G | Buildings | 8 | 8 | identical |
| H | Spells | 23 | 23 | 1 title differs by two literal quote characters only (`Radius, Length/2` vs `"Radius, Length/2"`) |

The full ROM2 title list (180 rows across all 8 groups; the name-title is excluded
from slot numbering exactly as `databin-slots.csv` excludes it, so the same
column carries the same slot number on both sides) is the evidence behind this
table; per-title text is not reproduced further here beyond the two differences
named above.

A per-collection param-array width census (the numeric `dwordArray` most entry
kinds carry, distinct from title text) finds every collection's width unchanged
between the two games except Units: Magic 28, Armors/Shields/Weapons 17,
MagicItems 2, Humans 26, Buildings 6, Spells 22 identical on both sides; Units
widens from 55 to 62, tracking the 56→63 title-count growth above net of the one
trailing string-type column that is not itself a numeric param. This is an
independent wire-level confirmation that group E's title growth is a real column
change, not a title-array-only edit. `R2-ASSET-031`.

## `templates.bin` result

A distinct file from `world.res`'s own data.bin (different size, different
content), but whether it shares DAT-GRAM-003's own wire grammar at all is
**Unknown**. The clean corpus-wide grammar (above) reaches the identical
divergence point `R2-ASSET-005` already found under the unmodified grammar —
group A, `Shapes` entry 112, offset 35 004 of 61 802 — because group A's own shape
is untouched by the group-C-only correction; this reconfirms but does not advance
`R2-ASSET-005`'s own reach. `R2-ASSET-032`.

`R2-ASSET-005`'s open question — a `u16`-count/`u8`-length-prefix reading vs a
content-free `u32`-count/`u32`-length reading of the same head bytes — is not
resolved by this survey.

A reader-identity search adds a bounded, separate fact: an exhaustive ASCII string
search (`templates.bin`, `Templates.bin`, `TEMPLATES.BIN`, `templates`, case-
sensitive) over `allods2.exe`'s and `a2server.exe`'s own loaded memory images —
a raw byte scan of every loaded memory block, not a defined-data-only search —
finds **zero** occurrences of any of the four spellings in either binary — against
a same-search positive control that DOES find `Data.bin` and `World.res` in both,
proving the search finds a filename that is present. Widened independently to all
12 `.exe`/`.dll` files under the preserved ROM2 root, under exact-case ASCII,
case-insensitive ASCII and UTF-16LE: the two game binaries stay at zero under
every encoding, and the positive controls reproduce. The literal string
`templates.bin` DOES occur once elsewhere in the same root, in `ROM2 Map
Editor.exe`, as a standalone filename among resource/configuration strings; no
cross-reference check was run on that binary, so this locates where the name is
referenced textually, not that the editor's own code reads the file by it. This
rules out a literal-name reference in the two game binaries under every spelling
and encoding tried; it does not rule out a dynamically composed path, and is
scoped to exactly the population searched. `R2-ASSET-032`.

## `a2server.exe` classification (Q5)

`allods2.exe`'s Data.bin-reading population — the Serialize dispatcher, the shared
`IsStoring()` helper, the 8 group-level Serialize call sites, and the 8 SetSize
call sites — classified against `a2server.exe`'s own function census under
`tools/enginematch`'s strict tier (`mnemHash`+`instrCount`+`byteLen`): the
dispatcher and 7 of 8 group Serialize call sites show 0 strict matches; the shared
`IsStoring()` helper shows exactly 1 (specific); one short, generic 11-instruction
wrapper and the SetSize family (a generic MFC `CArray::SetSize` template shape)
show 3–81 matches each, reported as same-shape hash-collision noise, not identity.

The noise reading is measured, not only argued, two ways. A self-collision census
(the identical strict triple, joined against `allods2.exe`'s own function census)
finds every discriminating address unique inside its own binary (matches itself
only), while the wrapper matches 102 other functions and the SetSize shapes match
11 and 5. A relocation- and `rel32`-tolerant raw byte search of `a2server.exe`'s
whole file image — an independent method sharing no code with the census —
corroborates the 0 matches for the dispatcher and the 7 substantive Serialize
bodies, and additionally finds 0 matches for all 8 SetSize addresses individually,
ruling out a relocation artifact as the reason the strict tier shows 0 there.
`R2-ASSET-033`.

## `allods2.exe`'s own loading code (Q5: which file does the map-reader's
## type-id resolution consult)

No literal reference to `templates.bin` (or two variant spellings) exists
anywhere in `allods2.exe`'s image, under any of 4 encodings/casings (exact
ASCII, case-insensitive ASCII, exact UTF-16LE, case-insensitive UTF-16LE) — the
case-sensitivity control confirms the instrument itself is working: `Data.bin`,
capital D, IS found, at 4 addresses, and lowercase `data.bin` is absent.
**Correction, finding 8:** `templates.bin` is not absent from the ROM2 root as
a whole — a root-wide search of all twelve `.exe`/`.dll` files in
`gameversions/rom2-ru` finds the literal string at file offset `0x7afdc` in
`ROM2 Map Editor.exe`, as a standalone NUL-terminated string among a run of
resource/configuration filenames, in ASCII and ASCII-case-insensitive matching;
it occurs in no other binary in the root. The negative is scoped to
`allods2.exe` itself (the map-reading binary Q5 asks about), not to the
preserved root: `templates.bin` is referenced by the map editor and by nothing
the game itself runs, which narrows toward H5-databin more than a single-binary
absence does on its own. `R2-ASSET-025`.
`allods2.exe` does contain its own copy of the Data.bin loading shape
`DAT-LOC-001` already publishes for `rom.exe`: a once-only guard flag, a loose
file at a path ending `Data.bin` falling back to the same relative path inside
`World.res`, and on failure of both, the same 11 named source tables
(`.txt` here, `.csv` on ROM1) then a `Data.bin` rewrite. **Correction, finding
9:** the eleven table names and the two path literals are now confirmed read
directly by a decompile of the function that reads them, `FUN_00501c9e` —
called from the guarded entry point between its "Parsing .txt files" and
"Writing new .bin file" messages, at the position ROM1's own `DAT-LOC-001`
places `FUN_004da351` — rather than inferred from their `.rdata` adjacency to
the two `Data.bin` literals alone, which is how the original pass read them.
Its only two callers are consistent with a one-time process/session-startup
call, not a per-map site — a structural signal, not a full trace. This
experiment did NOT locate ROM2's own counterpart of ROM1's
type-id-to-definition resolver (`ALM-CLS-036`), so whether an `.alm` record's
raw type id actually resolves against this Data.bin-built table remains
untraced; within `allods2.exe`, Data.bin is the structurally-present candidate
and templates.bin has no textual footprint at all, without the consumer link
itself being confirmed. `R2-ASSET-026`.

## Not yet surveyed

Entry CONTENT for any group beyond declared counts and byte spans (no field is
read); which of `templates.bin`'s two candidate head readings, if either, is its
real wire grammar, or what it holds past `Shapes` entry 111; the actual
`.alm`-record-to-Data.bin-table consumer function, the link that would settle
the map-reader's own type-id resolution outright; whether `templates.bin` is
read by the client/server through a dynamically composed path this survey's
string search cannot see, by `ROM2 Map Editor.exe`'s own code under the literal
name located there (no cross-reference check was run on that binary), or under
a needle beyond the `templates.bin`/`Data.bin` families searched here; a
non-Data.bin negative-control population for the Q5 strict tier, to
independently confirm the collision-noise reading of the 11-instruction
wrapper and the SetSize-family addresses beyond the self-collision and
relocation-tolerant checks already run; which specific other MFC
`CArray::SetSize` template instantiations, if any, the SetSize shapes collide
with beyond the 18 addresses this survey names.
