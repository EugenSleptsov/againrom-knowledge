# SPR16A (`.16a` / `.16` sprite) — specification

Level 3. Promoted claims only. The `.16a` container, the 16-bit word-RLE grammar, the
**pixel display model** (`SPR16A-PIX-011`, read out of the owned `rom.exe` and confirmed visually),
the **frame-count trailer** (`SPR16A-TRLR-012`), the **`.16` glyph grammar**
(`SPR16A-FONT-013`) and the **corpus bounds** (`SPR16A-BOUND-016`) are resolved.
**`SPR16A-RDR-017`, `SPR16A-FONT-018` and `SPR16A-FONT-019` read the readers**: both loaders, both
blitters and the font layer, closing every deferred reader-side question.
**`SPR16A-FONT-020…SPR16A-FONT-022` read the font high half and the tails**: Cyrillic in a hybrid arrangement, in both
releases' bytes, and the tails are older overwrite layers. Open: the exact per-display-mode
framebuffer packing (runtime), the `spritesb` blend rule, and how the RU release's CP866
strings meet these atlases.

**Read this first if you are writing a decoder.** `.16a` and `.16` are *not* one format with
a flag. They are parsed by two different routines, and the trailing `u32` means different
things to each (below). `.16a` is parsed by the **same** constructor as `.256`, which is why
the two containers are identical.

Seen as: `*.16a` (542) and `*.16` (3) file nodes inside `graphics.res` — the whole corpus;
an install-wide sweep of every `&YA1` container and every loose file finds no others. A
multi-frame sprite carrying a per-file 256-entry BGR0 table. It is the 16-bit analogue of
SPR256 (`formats/spr256`): same container shape and three RLE opcodes, but control tokens
and literal words are 16-bit.

## Cursor sheets (`SPR16A-CURSOR-046`)

The executable registers 23 `.16a` cursor sheets. Paths, hotspots, opaque period arguments and
file geometry are tabulated by `SPR16A-CURSOR-046`. Command-relevant entries:

| cursor | frames | size | hotspot | period argument |
|---|---:|---|---|---:|
| default | 1 | 32x32 | 5,5 | 2000000000 |
| move | 5 | 32x32 | 15,15 | 100 |
| swarm | 5 | 44x44 | 21,21 | 100 |
| attack | 10 | 32x32 | 3,3 | 100 |
| defend | 8 | 32x32 | 15,13 | 100 |
| select | 1 | 32x32 | 3,4 | 100 |
| patrol | 8 | 32x32 | 8,25 | 100 |
| cast | 14 | 32x32 | 15,15 | 100 |
| pickup | 15 | 32x32 | 12,13 | 66 |

All eight edge arrows are one 32x32 frame; small-default is one 16x16 frame; cantput is one
64x64 frame; town and backpack are one 32x32 frame; dice and wait are 15 and 10 32x32 frames.
All payloads are identical across EN/RU. The period argument's downstream time unit is not part of
the format contract. Frame count, dimensions and pixels are data; hotspot and period are executable
constants.

## The item-icon subtree (`SPR16A-ICON-027`)

The 416 `.16a` under `inventory/` are the most regular population in the archive and are worth
knowing about before writing a cache: **416 of 416 parse whole, 0 zero-frame, 0 multi-frame, 0
with bit 31 clear, and frame 0 is 80x80 on every one.** Their per-pixel level runs 1..15 — level
0 occurs in none of them and 15 does, so under `SPR16A-ALPHA-025` the alpha runs 2/16..16/16 and
the opaque value is used. sha256 per node: 416 identical across the two releases, 0 different.
A consumer may therefore size the item-icon cache at one 80x80 frame per node. It may **not**
generalise that to the rest of the corpus: `SPR16A-BOUND-016` mixes geometries freely, and the
`.256` equipment sheets beside these have four that yield 0 frames (`SPR16A-ICON-027`).

## Structure (`.16a`, standard form — `SPR16A-STRUCT-001`, 542/542)

```
[ 1024 B leading block ]           per-file 256 × [B,G,R,0] palette table
                                                                   SPR16A-PAL-008
repeat frames:
   u32  width
   u32  height                                                    SPR16A-STRUCT-001
   u32  dataSize
   u8   data[dataSize]             16-bit word-RLE — see below     SPR16A-RLE-002
u32  trailer                       [31-bit frameCount][bit31 = has-palette?]
                                                                   SPR16A-TRLR-012
```

- 542 `.16a` walk clean from offset **1024**; the residual tail is exactly **4**
  (as SPR256's trailer).
- **`frameCount = trailer & 0x7FFFFFFF`, and `trailer & 0x80000000` gates the 1024-byte
  palette read** — both quoted from `FUN_00428c50`, the constructor `.16a` shares with `.256`
  (SPR16A-TRLR-012, SPR256-TRLR-021). Read the trailer; do not only walk. The `16+16`
  alternative is refuted: the loader reads the trailer with a single 4-byte read and touches
  it only as a dword.
- **For a `.16` the same four bytes are a plain `u32` count with no flag** — `FUN_004284e0`
  applies no mask and never reads a palette (SPR16A-TRLR-012). Masking a `.16` trailer is
  harmless on everything that ships, but it is not what the engine does.
- A walker that ignores the trailer still terminates: bit 31 makes the trailer's `u32` an
  impossible width for the 542 `.16a`, and the 3 `.16` run out of bytes.
- **The loader validates nothing** (SPR16A-RDR-017). It indexes exactly `frameCount` records,
  `cursor += 12 + dataSize`, with no size threshold, no `w`/`h` check and no bound against the
  end of the buffer. There is no rejection path and no "null frame": a `dataSize == 0` record
  is a bare 12-byte header. **Your decoder must impose the bounds the engine does not** —
  *Defensible limits*, below.
- Every frame is compressed: `dataSize != w*h` and `!= w*h*2` (SPR16A-STRUCT-001).
- Frames within one sheet are **always the same size** — 0 of 542 sheets mix, the opposite of
  `.256`, where 131 of 1384 do (SPR16A-BOUND-016 vs SPR256-FRAME-023). A `.16a` consumer may
  cache one frame's bounds for the sheet; a `.256` consumer may not.
- Corpus maxima and the caps we would choose above them: *Defensible limits*, below.

## Frame pixel data — 16-bit word-RLE (`SPR16A-RLE-002`, `SPR16A-RLE-003`, 2727/2727 frames)

`data[dataSize]` is a run-length stream of **u16 LE control words**,
`[2-bit op | 14-bit count]`. A cursor moves left→right and wraps to the next row at
`width`; the background is transparent.

```
cw >> 14 == 0b00   literal   : the next n u16 words are 16-bit colour pixels
cw >> 14 == 0b01   blank rows: emit n fully-transparent rows (at a row boundary)
cw >> 14 == 0b10   skip      : emit n transparent pixels in the current row
cw >> 14 == 0b11   BLANK ROWS again — unused in ROM1 data, but the dispatch is two
                   sequential bit tests and bit 14 is tested FIRST, so both-bits-set
                   takes the blank-rows arm                      SPR16A-RLE-002
            n = cw & 0x3FFF
```

> **Do not carry `.256`'s alias across.** `.256` and the `.16` byte blitter both fold their
> fourth quadrant into *skip*; the `.16a` u16 blitter folds it into *blank rows*, because the
> two routines test the same two control bits in opposite order. Nothing shipped exercises
> either alias.

- **Exact / lossless** (SPR16A-RLE-003): every row's tokens sum to exactly `width`,
  the stream yields exactly `height` rows and is consumed with no leftover byte, and
  `Σ literal + Σ skip + Σ blankRows·width == Σ width·height` corpus-wide
  (34 639 040 == 34 639 040).
- Blank-row opcodes occur only at a row boundary (column 0).
- The `0b11` opcode quadrant never appears; observed count maxima literal ≤137,
  skip ≤551, blank ≤128 (the 14-bit field allows 16 383).
- "Transparent" here = *no pixel emitted*.

## Literal-word colour — the display model (`SPR16A-PIX-011`)

A literal u16 word is **not a colour**; it is an **even byte offset** into a runtime
`[16][256]u16` LUT the executable builds from the file's palette. This forces:

```
paletteIndex = (word >> 1) & 0xFF     // bits 1..8   (0..255)
level        = (word >> 9) & 0x0F     // bits 9..12  (0..15)  -- shipped range 1..15
src          = srcLUT[level][paletteIndex]      // palette scaled by (level+1)/16
pixel        = src + destTable[1 + level][old_fb]   // u16 add, then framebuffer write
```

**And the level is the pixel's alpha** (`SPR16A-ALPHA-025`). The destination
table `FUN_0044ba10` builds has **17** rows and its row `k` is every representable pixel
scaled by `(16 - k)/16`; the blitter's base skips exactly one row, so the row it uses for
level `L` scales the old pixel by `(15 - L)/16` while the source row scales the art by
`(L + 1)/16`. The two weights sum to 16, so the write is an exact 16-step linear blend:

```
out = palette[index] * (level+1)/16  +  destination * (15-level)/16
```

Level 15 is the only opaque value (its destination row is all zeros) and level 0 never
ships, so a `.16a` pixel is between 2/16 and 16/16 of the art. A decoder that renders the
palette colour opaquely is not slightly off — three shipped projectile sheets use a single
palette index and are *entirely* alpha (`SPR16A-PROJ-026`). Straight-alpha RGBA reproduces
this exactly up to the framebuffer's 5-6-5 quantization; a binary mask cannot.

The `/18` source variant and a destination table quantized to 8 192 entries a row are **not
a display mode**: `FUN_0044ba10` sets the flag that selects them from `GlobalMemoryStatus`,
under 24 MB of physical RAM (`PAL-MODE4-010`).

- Consistent with the surviving static fact: `OR=0x1FFE` over all **2,535,421** literals
  (bits 13–15 and bit 0 always zero — the low bit is 0 because LUT offsets are even).
- The leading 1024 B is the genuine per-file **256-entry `[B,G,R,0]` table** indexed by
  `paletteIndex`, distinct in 497/542 files (`SPR16A-PAL-008`).
- Instruction data flow closes `0x00457b15 → 0x00428ad0 → 0x00427df0` (LUT build) and
  `0x0042b970 → 0x00451ae0` (blit); 2727/2727 frames, all 2,535,421 literals, and a
  4096-pair synthetic fixture pass; altered mask/stride/selector controls fail.
- Confirmed visually on the anchor sprites (fire warm, freeze cyan, portraits legible with
  skin, mage purple, coins gold).

**Retracted en route** (kept as the negative-result record): direct-RGB
(`SPR16A-PIX-004`, `SPR16A-PIX-005`), the "block unused" misread (`SPR16A-PAL-006`), flat high-byte + gamma
(`SPR16A-PIX-009`), and the `5+7` / `8×32` bank models (`SPR16A-PIX-010`) — all failed owner visual
comparison; the `5+7` split is directly contradicted by the executable's `4+8` addressing.

## Font atlases — five of them (`SPR16A-FONT-007`, `SPR16A-FONT-013…SPR16A-FONT-015`, `SPR16A-FONT-020…SPR16A-FONT-022`)

The install ships **five** font atlases in two container forms (SPR16A-FONT-015), and
`rom.exe` loads **four** of them (SPR16A-FONT-018):

| file | form | glyphs | cell | advance range | notes |
|---|---|---:|---|---|---|
| `font1/font1.16` | byte-control | 224 | 16×15 | 0..14 | + an appended section chain |
| `font2/font2.16` | byte-control | 224 | 8×10 | 0..7 | + appended sections; **the one node the RU release replaces** |
| `font3/font3.16` | byte-control | 64 | 8×6 | 0..3 | 53 of the 64 records are empty |
| `font4/font4.16a` | **ordinary `.16a`** | 224 | 16×16 | 0..15 | decoded by the `.16a` path above |
| `font5/font5.16a` | **ordinary `.16a`** | 224 | 24×24 | 0..20 | **loaded by nothing in the install** |

Record 0 is empty in all five (`dataSize` 1 for the `.16`, 2 for the `.16a` — a single
blank-rows control). **Record `k` is character `32+k`**: this is not an inference from
224 = 256−32; the review instrument renders `font1` records 32..62 as `@ A B C … ^` and
`font4` records 33..63 as `A B … _`, and the engine's own `DrawText` subtracts `0x20` from
every character before indexing. The four 224-record atlases therefore cover 32..255 and
`font3` covers 32..95 — of whose 64 records 53 are empty, leaving the digits and a few marks.

### The high half — CP437 letters plus Cyrillic in a hybrid arrangement (SPR16A-FONT-020)

Records 96..223 of every 224-record atlas hold, under `char = record + 32`:

```
0x80..0x9A   CP437's accented Latin (Ç ü é … Ö Ü); its ¢£¥₧ƒ dropped (blank)
0xA0..0xA5   á í ó ú ñ Ñ; CP437's ª..» dropped
0xB0..0xCF   А..Я        — the box-drawing region carries the Cyrillic uppercase
0xD0..0xDF   а..п
0xE0..0xEF   blank except ß at 0xE1
0xF0..0xFF   р..я
```

This arrangement is **none of CP866 / CP1251 / KOI8-R** (each needs letters where these
fonts are blank); 17 of 18 Latin/Cyrillic homoglyph pairs (A/А … x/х) decode pixel-identical
at exactly these positions in font1 and font2. It is the same in the **English** release's
bytes: after the 2026-07-29 GOG update swapped the install to the Russian release (and the
EN reinstall restored the pinned bytes), font1/font3/font4/font5 and all five `.dat` proved
**byte-identical across the two releases** (sha256 equal). Only `font2.16` differs: the RU
build blanks its 35 accent records (chars `0x80..0x9A`, `0xA0..0xA5`, `0xE1`, `0xEF`) and
carries a different dead residue (SPR16A-FONT-022). Note the measured mismatch a consumer
should know about: the RU release's own data strings are **CP866** and its README **CP1251**
(SPR16A-TXT-023, partially retracted) — neither matches these atlases'
arrangement, and `rom.exe` converts nothing; how the RU game shows Russian text is an open
question.

### A font is two nodes — `fontN.dat` is the advance table (SPR16A-FONT-018)

Each font object loads **`<base>.16`/`.16a` *and* `<base>.dat`**. The sidecar is one `u32`
per glyph (896 B = 224, or 256 B = 64 for `font3`) and text layout is

```
x += dat[glyph] + spacing          // spacing = 2, from the construction site
x += height(0)/2 + dat[0] + spacing   // for glyph 0, the space
```

**The glyph cell is never the advance** — `dat[g]` is strictly less than the record width on
all 960 glyphs of all five atlases. A consumer that advances by the cell renders monospaced
text at the wrong pitch; the review instrument shows both layouts side by side.

### `.16` glyph pixel grammar (SPR16A-FONT-013)

A `.16` glyph record has the same `[u32 w][u32 h][u32 dataSize][data]` shape, but `data` is a
**byte** control stream, not the `.16a` word stream:

```
control byte c:  op = c >> 6 ,  n = c & 0x3F
  op 00  literal   : the next n BYTES follow, each carrying TWO 4-bit pixels,
                     LOW nibble = the left-hand pixel. A ZERO HIGH nibble ENDS
                     the run and is PAD, not a pixel — the run is 2n-1 px long.
                     A zero LOW nibble is a WRITTEN pixel, value 0, not transparent.
  op 01  blank rows: emit n fully-transparent rows (at a row boundary)
  op 10  skip      : emit n transparent pixels in the current row
  op 11  SKIP too  : the dispatch tests 0x00 then 0x40, so 0x80 and 0xC0 share the
                     skip branch (0 of 14 845 shipped control bytes reach it)
```

- **Exact**: 687/687 glyph records across every section of all three `.16` satisfy the same
  invariant the `.16a` grammar is held to — every row lands exactly on `width`, exactly
  `height` rows result, the stream is consumed with no leftover byte, and **no run ever passes
  the row end**. `Σ tokens == Σ w·h` per file, with zero slack.
- The pad is not a convention we picked: over 24 689 literal data bytes a zero nibble occurs
  **only** as the high nibble of a run's final byte (2 896 times; 0 mid-run; 0 in any low
  nibble). That is what fixes both the odd-run length and the nibble order — the pad must be
  the run's *last* pixel, so the low nibble is the left-hand one.
- Observed pixel values: `{4,5,6,7,8,9,11,13,15}`. **The 4-bit value is an intensity level of
  a text colour the caller chooses** (SPR16A-FONT-013). It is *not* a level into the `.16a`
  path's 16-level LUT — a `.16` object has no palette and never builds one — and *not* a
  palette index. The blitter indexes a **16-entry `u16` table passed as its sixth argument**
  (`MOV AX,[EBX + EAX*2]`, so the whole table is 32 bytes), and the engine builds thirteen of
  them, each entry `base * k / 15` packed to the active framebuffer format. Rendering the
  value as coverage is therefore not a display convention — it is what the format means:

```
pixel = ramp[v]        // ramp = 16 u16 entries chosen by the text's caller
                       // ramp[k] ~ baseColour * k / 15, packed to the framebuffer
```

- The write is **opaque** — unlike the `.16a` u16 blitter there is no destination read and no
  additive blend. A `.16` glyph replaces the pixels it covers.

### The `.16` section chain (SPR16A-FONT-014, SPR16A-FONT-021)

`font1`/`font2` do not end at their count trailer. Each continues into further
`[records][4-byte trailer]` sections, every trailer an identical copy of the first:

```
EN font1.16  224 recs 16x15 | e0000000 | (26 B) | 11 recs | e0000000 | (12 B) | 115 recs | e0000000
EN font2.16  224 recs  8x10 | e0000000 | (57 B) | 49 recs | e0000000
RU font2.16  224 recs  8x10 | e0000000 | (28 B) | e0000000          (font1.16 identical EN/RU)
font3.16      64 recs  8x 6 | 00000040
```

**All of it is dead** (`SPR16A-FONT-014`): the loader indexes exactly `frameCount`
records from the front of the buffer and has no other bound. **And it is now explained**
(`SPR16A-FONT-021`): each section is the remnant of an **older build of the same
file**, overwritten in place without truncation — SPR256's Bucket-B shape (`SPR256-EXC-017`)
shown mechanically. Every boundary run is the byte-suffix of a record cut by the next build's
end (the 12-byte run equals the last 12 bytes of a solid-`0xd` 16×15 record byte-for-byte;
the 57-byte run the last 57 of a solid 8×10 record, cut mid-`h`-field), and counting backward
from each remnant's own trailer (all say 224) fixes the positions: font1's 11 art records are
**positions 213..223 = х..я** — an older, bolder draw of the same glyphs the live font keeps
there — and the 115/49 solid-`0xd` cells are placeholder records of a not-yet-drawn stage,
positions 109..223 / 175..223. **A decoder should stop at the count**; the boundary runs need
no framing because nothing frames them — they are cut record suffixes.

## Defensible limits — a decision, not a measurement

`SPR16A-BOUND-016` measures what shipped; the format itself bounds nothing below its field
types. A decoder needs caps anyway, so the two are kept apart: the middle column is evidence,
the right column is a choice.

| quantity | measured max (whole corpus) | field type | suggested cap (a decision) |
|---|---:|---|---:|
| frame width | 552 | `u32` | 2048 |
| frame height | 128 | `u32` | 2048 |
| frames per sheet | 224 | 31 bits | 4096 |
| `dataSize` per frame | 20 892 | `u32` | 2^24 |
| payload size | 206 328 | — | 2^26 |
| `.16a` literal run | 137 words | 14 bits | the field's own 16 383 |
| `.16a` skip run | 551 px | 14 bits | the field's own 16 383 |
| `.16a` blank-row run | 128 rows | 14 bits | the field's own 16 383 |
| `.16` literal run | 8 bytes (16 px) | 6 bits | the field's own 63 |
| `.16` skip run | 15 px | 6 bits | the field's own 63 |
| `.16` blank-row run | 15 rows | 6 bits | the field's own 63 |

Hard checks a decoder can assert from a measurement rather than from taste:

- `.16a` `dataSize` is **even** in 2727/2727 frames — it is a `u16` stream.
- `dataSize >= 2`; 145 shipped frames are wholly empty (one blank-rows control).
- `(literal word & 0xE001) == 0` — 0 violations in 2 535 421 words.
- The level field `(word >> 9) & 0xF` is **never 0** in shipped data (values 1..15 only);
  all 256 palette indices do occur.
- Frames within a sheet share `(w,h)` — 0 of 542 sheets violate it.

## Open (own experiments)

- **The framebuffer's own channel widths and shifts** — six runtime globals the table
  builders and the blend all read, so the *absolute* colour of a drawn pixel is derived only
  up to the assumption that they are 5-6-5. The same residual applies to the `.16` glyph
  ramps. What is **no longer** open: the arithmetic of `src + dest` itself, and the pair of
  variants that used to be recorded here as display modes — they are a low-memory fallback
  (`SPR16A-ALPHA-025`, `PAL-MODE4-010`).
- **Which of the thirteen ramps a given UI string gets** — traced to the blit argument, not
  through the UI code that chooses it.
- **The `^` escape** in the `.16` `DrawText` — a two-character sequence that diverts to a
  separate draw call; located, not decoded.
- **Animation frame-block roles** — as with SPR256, per-class block sums vs registry
  phase counts, once cross-referenced.
- **`font5.16a` ships and nothing loads it** — whether it is a cut asset or reached by a path
  we have not found is not a format question, but it is unexplained.
- **The RU text pipeline** — the RU release's data strings are CP866, its README CP1251, its
  fonts the hybrid arrangement, and `rom.exe` (byte-identical across releases) converts
  nothing (SPR16A-TXT-023, partially retracted). Only runtime observation of the
  RU game can show what its text looks like on screen.

## Closed

- **The font high half and the tails** — CP437 accents + Cyrillic in the hybrid arrangement,
  byte-identical EN↔RU except `font2.16`; the appended sections are older in-place-overwrite
  layers of the same font, boundary runs = cut-record suffixes (`SPR16A-FONT-020`,
  `SPR16A-FONT-021`, `SPR16A-FONT-022`).
- **The readers** — two loaders (`.16a` shares `.256`'s), the trailer's two readings, the
  bit-31 palette gate, the total absence of frame-header validation, both blitters' fourth
  quadrant, the `.16` glyph ramp and the `.dat` advance table (`SPR16A-RDR-017`,
  `SPR16A-FONT-018`, `SPR16A-FONT-019`, and the amendments to `SPR16A-TRLR-012`,
  `SPR16A-RLE-002`, `SPR16A-FONT-013`, `SPR16A-FONT-014`).
- **`.16a` container** — `[1024B block][frames][4B trailer]`, frame
  `[u32 w][u32 h][u32 dataSize][data]`, **count in the trailer's low 31 bits**
  (`SPR16A-STRUCT-001` for the layout; `SPR16A-TRLR-012` for the count and loader).
- **`.16` glyph grammar** — byte control; two 4-bit pixels per literal byte, low nibble first,
  with an odd-run pad; 687/687 exact (`SPR16A-FONT-013`).
- **The corpus census** — five font atlases, the measured bounds, and the invariants a decoder
  can assert (`SPR16A-FONT-015`, `SPR16A-BOUND-016`).
- **16-bit word-RLE grammar** — `[2-bit op][14-bit count]` words; literal / blank-rows
  / skip; proven exact/lossless on all 2727 frames (`SPR16A-RLE-002`, `SPR16A-RLE-003`).
- **Pixel display model** — literal word = even byte offset into a runtime `[16][256]u16`
  LUT: `index=(word>>1)&0xFF`, `level=(word>>9)&0xF`, source × `(level+1)` + level-selected
  destination term (`SPR16A-PIX-011`; instruction-derived, corpus/synthetic +
  visually confirmed). Retracted `SPR16A-PIX-004`, `SPR16A-PIX-005`, `SPR16A-PIX-009`
  and `SPR16A-PIX-010` colour models are the negative-result record.

## Projectile sheets

`SPR16A-PROJ-024`. The 31 `projectiles.reg` rows resolve to **24 `.16a`** sheets and
**7 `.256`**; which extension is read is the registry's `A16` key, not a probe of the bytes. Frame
geometry is uniform inside every sheet and the walked frame count equals the trailer count 31 of 31.

The consumer-facing law is the draw's: a sheet must hold `Phases * 9` frames when the registry's
`Flip` is set and `Phases * RotationPhases` when it is not. **29 of the 31 hold exactly that.** The
two that do not are `healing` (7 `Phases`, 8 frames — one frame unreachable) and `goblin\arrow`
(one frame plus an 8-byte residue before the count trailer, the `SPR256-EXC-017` shape).

`Width`/`Height` in the registry are the draw's centring offsets and need not equal the frame: 8 of
the 31 rows omit them and take the loader's `0x40` default against a real 12x12 or 64x96 frame.

**How one of these sheets is drawn** (`SPR16A-ALPHA-025`, `REG-PROJ-087`). `A16` picks the extension *and* the C++ class —
the `.256` and `.16a` sprite vtables differ in exactly two slots, the destructor and `+0x18`, so
the class is what chooses the blit. The registry's `Palette` is a **boolean**, "this sheet carries
its own colour table": non-zero makes the loader build the sprite a 16-level lookup (mode 4, no
tint, for `.16a`), zero makes it build nothing and the draw hands the blit the shared table from
`projectiles.pal` instead. It agrees with the sheet's own trailer bit 31 on **31 of 31** rows —
and the three rows where it is absent are the three `.256` sheets whose bit 31 is clear, so
**their frames begin at offset 0, not 1024** (`REG-PROJ-087`). Level 0 occurs 0 times in 247 883
lit projectile pixels; seven of the 24 sheets contain no opaque pixel at all
(`SPR16A-PROJ-026`).

## Spell art: which sheets one cast can reach

`SPR16A-CAST-028`. The 28 shipped spells joined to `projectiles.reg` by `2*id + 8` (the
cast picture) and `2*id + 9` (the burst picture), each sheet walked on both preserved roots with
identical results.

```
cast parity  (+8), 15 defined rows
  firebolt(10) fireball(12) p_fire(18) healing(20) poison_d(24) p_water(28) Drain(30)
  lightnin(34) chain(36) p_air(40) shield(44) p_earth(52) bless(54) teleport(60) Curse(62)

burst parity (+9), 8 defined rows
  fireexpl(13) firewall(15) smallxpl(17) freeze(23) poison(25) acid(27) wall(47) Meteor(51)

both parities defined: fire_ball and poison_cloud, and no other spell
neither parity defined: light, invisibility, darkness, stone_curse, haste,
                        control_spirit, slow
```

**Orientation.** Exactly two of the 23 spell-reachable rows carry `RotationPhases = 16` with
`Flip = 1`: `firebolt` and `fireball`, both `Phases = 4`, both walking 36 frames, which is
`Phases * 9` under the mirrored-facing law. Every other spell-reachable row is `RotationPhases = 1`,
so its frame index is the phase alone and it stores no per-direction row.

**Frame law.** Walked count equals `Phases * RotationPhases` on 22 of 23; the exception is
`healing`, 8 walked against `Phases = 7`, the frame the draw can never reach.

**Reachability.** Of the 31 shipped rows, `steam`(8) is named by no spell id at either parity and by
no `units.reg` `Projectile` key.

**Release parity.** The 23 spell-reachable sheets, the two `smoke%d` sheets the draw uses and
`projectiles.reg` itself are 26 of 26 byte-identical across the two releases.

## Which consumer draws each spell-reachable sheet (`SPR16A-PART-030`)

Promoted from `SPR16A-MARK-029` (partially retracted) and `SPR16A-PART-030`. The 23 spell-reachable rows above are drawn
by four different consumers.

**Cast parity, `2*id + 8`, 15 rows.** Ten are actor-bound marks drawn from the mark array on the
actor: `p_fire`(18) `healing`(20) `poison_d`(24) `p_water`(28) `Drain`(30) `p_air`(40) `shield`(44)
`p_earth`(52) `bless`(54) `Curse`(62). Five are moving cast objects: `firebolt`(10) `fireball`(12)
`lightnin`(34) `chain`(36) `teleport`(60). `healing` and `Drain` are in both groups, because the
cast-spawn switch also gives them a one-tick flight.

**Burst parity, `2*id + 9`, 8 rows.** All eight are burst art. Four of them are additionally painted
on ground cells under an area effect: `firewall`(15) `freeze`(23) `poison`(25) `wall`(47). The other
four are burst-only: `fireexpl`(13) `smallxpl`(17) `acid`(27) `Meteor`(51).

Every row in the mark group and the ground group has `RotationPhases = 1`. The two rotating sheets
are both in the cast-object group.

### Frame demand of the mark path

Corrected for Heal and Drain Life by `SPR16A-031`.

The frame index a mark carries comes from an engine immediate, not from the sheet. Each arm can
produce a bounded set:

| record index | sheet | frames the arm can index | declared `Phases` |
|---|---|---|---|
| 0x12 0x1c 0x28 0x34 | `p_fire` `p_water` `p_air` `p_earth` | 6 | 6 |
| 0x14 | `healing` | 8 | 7 |
| 0x18 | `poison_d` | 6 | 8 |
| 0x1e | `Drain` | 8 | 9 |
| 0x2c | `shield` | 5 | 5 |
| 0x36 | `bless` | 5 | 5 |
| 0x3e | `Curse` | 5 | 5 |

Seven of the ten agree exactly. `poison_d` carries two frames no actor mark can reach and `Drain`
carries one, frame 8. `healing` is the inverse registry exception: the sheet walks eight frames,
`Phases` is 7, and its mark reaches all eight, 0 through 7. The terminal phase is produced because
the builder tests the old phase below 7, then increments and copies it (`SPR16A-031`).

Nine of the ten rows are 12x12 in the registry and `healing` is 16x16; the mark draw subtracts
exactly `Width/2` and `Height/2`, so those are centring halves here as everywhere else.
