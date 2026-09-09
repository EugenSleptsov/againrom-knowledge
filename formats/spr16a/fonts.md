# Font atlases and glyphs

[Reference](format.md)

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

Glyph record `k` is selected by display byte `32+k`. The four 224-record
atlases cover bytes 32..255; `font3` covers 32..95. Record 0 is a blank
space glyph. Use the separate advance table for spacing.
— SPR16A-FONT-007, SPR16A-FONT-018

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

The atlas arrangement is not CP866, CP1251 or KOI8-R. The RU `font2.16`
blanks its accent records at bytes `0x80..0x9a`, `0xa0..0xa5`, `0xe1` and
`0xef`; the other named atlases retain their EN forms. — SPR16A-FONT-020,
SPR16A-FONT-022

The no-conversion clause of SPR16A-TXT-023 is partially retracted. The
selector-dependent byte transform is specified by
[TEXT](../text/format.md#encoding-model); atlas indexing happens after that
transform. — SPR16A-TXT-023, TEXT-CONV-001 (partially retracted), TEXT-LANG-002

### A font is two nodes — `fontN.dat` is the advance table (SPR16A-FONT-018)

Each font object loads **`<base>.16`/`.16a` *and* `<base>.dat`**. The sidecar is one `u32`
per glyph (896 B = 224, or 256 B = 64 for `font3`) and text layout is

```
x += dat[glyph] + spacing          // spacing = 2, from the construction site
x += height(0)/2 + dat[0] + spacing   // for glyph 0, the space
```

The cell width is distinct from `dat[glyph]`. Use the advance value and
the space special case above for text layout. — SPR16A-FONT-018

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
                     skip branch
```

- A complete glyph stream consumes its data and produces exactly `width*height`
  pixels/skips. Literal runs do not cross a row end. A final zero high nibble
  is padding; a zero low nibble is a written intensity value. — SPR16A-FONT-013
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

The loader indexes exactly `frameCount` records from the front. Bytes after
that live frame list are not indexed. The named residual sections are older
in-place overwrite layers and cut record suffixes. They do not extend the
live count or define a new section grammar. — SPR16A-FONT-014,
SPR16A-FONT-021, SPR16A-FONT-022
