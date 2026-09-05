# ROM2 sprite / palette containers — identity survey

Level 3. Promoted, evidence-backed claims only. Basis: `R2-ASSET-007` (`.16a`/
`.256`), `R2-ASSET-008` (`.pal`), `R2-ASSET-009` (`.16`). Cross-reference only, not
evidence: ROM1's own [`formats/spr16a/format.md`](../spr16a/format.md),
[`formats/spr256/format.md`](../spr256/format.md),
[`formats/pal/format.md`](../pal/format.md), and the pre-existing ROM1 claims
`SPR256-EXC-017`/`SPR256-EXC-020`, `SPR16A-FONT-014`/`SPR16A-FONT-021`.

Seen as: `.16a`, `.16`, `.256` (inside `.res` containers, mainly `graphics.res`)
and `.pal` (same), found by the extension census over all 11 containers — every
result below is the FULL population found that way, not a capped sample (the
committed `spr-sample.tsv`/`pal-sample.tsv` files themselves cap at 8 rows per
container/extension for size; the counts here do not).

## `.16a` / `.256` result

Identical container framing. `.16a`: **623/623** (100%) walk with exact tiling
from the trailer-derived frame count through every frame's `[w,h,dataSize,data]`
header to the trailer boundary, 0 exceptions (ROM1 EN and RU: 542/542 each, same
code, positive control). `.256`: **1916/1929** (99.3%) walk exactly clean; of the
remainder, 8 are 0-byte stub entries (not evaluable — e.g. `cursors/cast.256`,
`cursors/defend.256`) and **5 show residue** before the trailer despite every
frame validating internally. All 5 are **sha256-identical** to the
identically-named node in ROM1 RU — the already-published `SPR256-EXC-017`
Bucket-B pattern (a bracketed secondary section appended after the frame data,
resolved by `SPR256-EXC-020` as loader-inert) reproduced byte for byte, not a
ROM2-specific extension. `R2-ASSET-007`.

## `.pal` result

Identical to one of ROM1's two known shapes (BMP colour table at a fixed seek, or
a flat 16×1024-byte block) on **159/159** (100%) of the full population — 0
matching neither shape. Of the 156 matching the BMP shape, all **156/156** also
pass the strict field-value test `bfOffBits == 0x436 && biBitCount == 8` — 0
partial matches. Both ROM1 EN and RU show the identical proportional breakdown on
their own smaller population, same code: **82/82** total, **79/79**
shape-A-and-strict, 3/3 shape B, 0 matching neither, on both roots — the
positive control this test previously lacked. `R2-ASSET-008`.

## `.16` result

A complete population of 3 files. All 3 have no palette (matching ROM1's own rule
that a bare `.16` never carries one) and every individual frame in every file
validates (no frame's own header/data ever exceeds file bounds). Only 1 of 3
(`font3.16`, 64 frames) tiles exactly to the trailer; the other 2 (`font1.16` and
`font2.16`, both 224 frames) leave unexplained bytes between the last frame and
the trailer — 17716 and 32 bytes respectively. Both residue files are
**sha256-identical** to the identically-named node in ROM1 RU: this is the
already-published `SPR16A-FONT-014`/`SPR16A-FONT-021` finding (older
in-place-overwrite build layers of the same font, each boundary run the
byte-suffix of a record cut by the next build's own end) reproduced byte for byte
in ROM2, not a ROM2-specific residue pattern. `R2-ASSET-009`.

## Not yet surveyed

RLE/pixel decoding, palette colour values, whether any `.16a`/`.256`/`.pal`/`.16`
file OUTSIDE this survey's full-population walk (there is none — the walk covers
every entry the extension census found) differs. What the `.16`/`.256` residue
bytes hold at the field level is not re-decoded here — `SPR16A-FONT-014`/`-021`
and `SPR256-EXC-017`/`-020` are its own, cross-referenced, not repeated.
