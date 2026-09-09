# Sprite encoding and pixels

[Reference](format.md)

## Sprite container

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

- **`frameCount = trailer & 0x7FFFFFFF`, and `trailer & 0x80000000` gates the 1024-byte
  palette read** — both quoted from `FUN_00428c50`, the constructor `.16a` shares with `.256`
  (SPR16A-TRLR-012, SPR256-TRLR-021). Read the trailer; do not only walk. The `16+16`
  alternative is refuted: the loader reads the trailer with a single 4-byte read and touches
  it only as a dword.
- **For a `.16` the same four bytes are a plain `u32` count with no flag** — `FUN_004284e0`
  applies no mask and never reads a palette (SPR16A-TRLR-012). Masking a `.16` trailer is
  harmless on everything that ships, but it is not what the engine does.
- **The loader validates nothing** (SPR16A-RDR-017). It indexes exactly `frameCount` records,
  `cursor += 12 + dataSize`, with no size threshold, no `w`/`h` check and no bound against the
  end of the buffer. There is no rejection path and no "null frame": a `dataSize == 0` record
  is a bare 12-byte header. A bounded decoder must check record extents before reading them.
- Every frame is compressed: `dataSize != w*h` and `!= w*h*2` (SPR16A-STRUCT-001).
- Installed `.16a` sheets have uniform frame dimensions. The frame grammar
  still stores dimensions on every record; this observation is not a general
  format limit. — SPR16A-BOUND-016

## Word RLE

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

- Each row's literal and skip counts sum to width. The stream describes
  height rows and has no remaining byte after the last row. — SPR16A-RLE-003
- Blank-row opcodes start at column 0.
- The `0b11` quadrant takes the blank-row arm, because bit 14 is tested
  before bit 15; installed streams do not use this alias. — SPR16A-RLE-002 The count field is 14 bits (maximum 16383); smaller
  installed maxima are not encoding limits.
- A transparent span emits no pixel.

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
ships, so a `.16a` pixel is between 2/16 and 16/16 of the art. Three installed projectile sheets use a single
palette index and encode their variation through alpha (`SPR16A-PROJ-026`). Straight-alpha RGBA reproduces
this exactly up to the framebuffer's 5-6-5 quantization; a binary mask cannot.

The `/18` source variant and a destination table quantized to 8 192 entries a row are **not
a display mode**: `FUN_0044ba10` sets the flag that selects them from `GlobalMemoryStatus`,
under 24 MB of physical RAM (`PAL-MODE4-010`).

Installed literal words leave bit 0 and bits 13–15 clear. Bit0 being zero
makes the LUT offset even. This installed pattern is not an arbitrary-input
validation rule. The optional 1024-byte palette has 256 entries of `[B,G,R,0]`.
— SPR16A-PIX-011, SPR16A-PAL-008

The direct-RGB models (`SPR16A-PIX-004`, `SPR16A-PIX-005`), unused-palette
interpretation (`SPR16A-PAL-006`), high-byte/gamma model (`SPR16A-PIX-009`)
and 5+7 or 8×32 bank models (`SPR16A-PIX-010`) are retracted. The established
address split is 4 level bits and 8 palette-index bits.

## Structural bounds

| Quantity | Encoding |
|---|---|
| Frame width, height, dataSize | u32 each |
| `.16a` frame count | Low 31 bits of final trailer |
| `.16` frame count | Complete final u32 |
| `.16a` RLE count | 14 bits, maximum 16383 |
| `.16` RLE count | 6 bits, maximum 63 |

An indexed record occupies `12+dataSize` bytes. Check those extents against
the containing input and the decoded cursor against the frame bounds.
`.16a` data is a u16 stream and therefore has even byte length.
The original loader does not impose these defensive extent checks.
Installed maxima, uniform sheet dimensions and unused literal bit patterns
are not arbitrary-input admission limits. — SPR16A-RDR-017, SPR16A-BOUND-016

## Unknowns

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
- `font5.16a` has no identified loader; its runtime use remains Unknown.

## Read and write sequence

1. Select `.16a` or `.16` from the resource contract; do not infer a shared
   codec from the frame-header shape.
2. Read the final trailer. For `.16a`, extract the low-31-bit count and optional
   palette flag. For `.16`, use the complete count and no palette.
3. Walk exactly that many `12+dataSize` records from the selected origin.
4. Decode `.16a` with word controls and its shade-table blend, or `.16` with
   byte controls, low-nibble-first literals and the caller's text ramp.
5. For text, load the paired `.dat` advance table and use its space rule.

An ordinary `.16a` emitter writes the optional palette, frame headers/RLE data
and count/flag trailer. A `.16` emitter writes its glyph records and plain
count trailer; an odd literal run uses the final high nibble as zero padding.
Widths, encoded lengths and counts must describe the emitted records. No
writer for arbitrary retained overwrite tails is established.
— SPR16A-TRLR-012, SPR16A-RLE-002, SPR16A-FONT-013,
SPR16A-FONT-014, SPR16A-FONT-018
