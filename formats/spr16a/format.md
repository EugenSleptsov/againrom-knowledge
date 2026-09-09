<a id="spr16a-16a--16-sprite--specification"></a>

# SPR16A sprites and `.16` fonts

`.16a` uses the SPR256 frame container with u16-control RLE. `.16` uses a
separate loader, byte-control glyph data and a plain u32 frame count. The
trailer word therefore has a different interpretation in each format. — SPR16A-STRUCT-001,
SPR16A-TRLR-012, SPR16A-FONT-013

A `.16a` literal selects a source shade-table entry and a destination blend
level. A `.16` literal contains two four-bit text intensities. Neither is a
literal framebuffer word. — SPR16A-PIX-011, SPR16A-FONT-013

The `.16a` source offset is the complete raw u16; the normal and low-memory
destination-table origins differ. Reversed clipping also depends on run
partition. The selected decoder contract is in [encoding](encoding.md).
— SPR16A-080, SPR16A-081

## Wire layout

| Part | .16a | .16 |
|---|---|---|
| Leading palette | 1024 bytes only when trailer bit 31 is set | None |
| Frame header | Three u32: width, height, encoded byte count | Same widths |
| Frame data | u16 controls; 14-bit counts; literal u16 lookup offsets | Byte controls; 6-bit counts; paired 4-bit intensities |
| Final u32 | Low 31-bit frame count and bit 31 palette flag | Complete frame count |

The standard size is `paletteBytes+sum(12+dataSize)+4`. All multibyte
fields are little-endian. The original loaders index the declared frames
without defensive bounds checks. — SPR16A-STRUCT-001, SPR16A-TRLR-012,
SPR16A-RDR-017, SPR16A-FONT-013

## Read and write order

Select the codec from the resource contract, read the trailer, then walk
the declared frames from the palette-dependent origin. Decode .16a word
RLE through its palette/blend tables; decode .16 byte RLE through the
caller's text ramp and paired .dat advance table. To write the standard
form, emit the optional palette, frame headers and encoded data, then the
count/flag trailer. Retained overwrite tails have no general writer.
— SPR16A-RLE-002, SPR16A-FONT-014, SPR16A-FONT-018

## Reference map

| Reference | Contents |
|---|---|
| <a id="structure-16a-standard-form--spr16a-struct-001-542542"></a><a id="frame-pixel-data--16-bit-word-rle-spr16a-rle-002-spr16a-rle-003-27272727-frames"></a><a id="defensible-limits--a-decision-not-a-measurement"></a><a id="open-own-experiments"></a><a id="closed"></a><a id="sprite-container"></a><a id="word-rle"></a><a id="literal-word-colour--the-display-model-spr16a-pix-011"></a><a id="structural-bounds"></a><a id="unknowns"></a><a id="read-and-write-sequence"></a> [Sprite encoding and pixels](encoding.md) | Container, word RLE, blend, bounds and emission |
| <a id="font-atlases--five-of-them-spr16a-font-007-spr16a-font-013spr16a-font-015-spr16a-font-020spr16a-font-022"></a><a id="the-high-half--cp437-letters-plus-cyrillic-in-a-hybrid-arrangement-spr16a-font-020"></a><a id="a-font-is-two-nodes--fontndat-is-the-advance-table-spr16a-font-018"></a><a id="16-glyph-pixel-grammar-spr16a-font-013"></a><a id="the-16-section-chain-spr16a-font-014-spr16a-font-021"></a> [Font atlases and glyphs](fonts.md) | Byte glyph RLE, text slots and advance tables |
| <a id="cursor-sheets-spr16a-cursor-046"></a><a id="the-item-icon-subtree-spr16a-icon-027"></a><a id="projectile-sheets"></a><a id="spell-art-which-sheets-one-cast-can-reach"></a><a id="which-consumer-draws-each-spell-reachable-sheet-spr16a-part-030"></a><a id="frame-demand-of-the-mark-path"></a> [Cursor, item and spell resources](resources.md) | Cursor, inventory, projectile and spell sheet consumers |

Runtime member offsets identify original in-memory fields. Wire offsets and
byte order are stated separately in the relevant layouts.
