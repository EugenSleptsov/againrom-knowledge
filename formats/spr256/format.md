# SPR256 (`.256` sprite) — specification (partial)

Level 3. Promoted claims only. **In progress** — structure (`SPR256-STRUCT-001`), RLE
pixel decoding (`SPR256-RLE-007`, `SPR256-RLE-022`), palette colour/transparency
(`SPR256-PAL-011`, `SPR256-KEY-044`), the `spritesb` overlay class
(`SPR256-OVL-014`), the trailer/frame-count and residual reads
(`SPR256-TRLR-016`, `SPR256-TRLR-021`), and the loader in `rom.exe`
(`SPR256-RLE-020`, `SPR256-EXC-020`) established. The loader claims close the `0xC0`
opcode question (alias of `0x80`) and the Bucket-B inner block (runtime-inert). Still
open: the overlay blend rule and the 2-file attack/pickup anomaly (both still need
runtime/`rom.exe`).

Seen as: `*.256` inside archives. A paletted, multi-frame sprite.

## Minimap cursor sheets (`SPR256-CURSOR-046`)

The five mode cursors are `cursors/{smove,sattack,sdefend,spatrol,scast}.256`. Each is one
palette-bearing 16x16 frame, parses whole to trailer count 1 and is byte-identical across EN/RU.
The executable registers each at hotspot `(0,0)` with the large static-cursor period argument.
Construction, overview drawing and the physical minimap handler establish their use; the `s`
filename prefix is not the identification evidence (SPR256-CURSOR-046). Pixels, geometry and
count are data; hotspot and registration order are executable constants.

## Structure (standard form — `SPR256-STRUCT-001`, 1370/1384 non-empty)

```
[ 1024 B palette ]              256 × [B,G,R,0] (SPR256-PAL-003/011)
repeat frames:
   u32  width
   u32  height                                                  SPR256-STRUCT-001
   u32  dataSize
   u8   data[dataSize]          RLE pixels — see below          SPR256-RLE-007
[ 4 B trailer ]  = [31-bit frameCount][bit31 = has-palette]     SPR256-TRLR-016/021

frameCount = number of frame records before the trailer         SPR256-COUNT-002
             AND trailer & 0x7FFFFFFF (the two agree exactly)   SPR256-TRLR-016/021
```

- `u32@0` of the payload is palette entry 0 (≈always 0), **not** a frame count.
- The **trailer is `[31-bit frameCount][bit31 = has-palette flag]`** (LE). The loader
  (`rom.exe`, `SPR256-TRLR-021`) computes `frameCount = trailer & 0x7FFFFFFF` and reads the
  1024-byte palette **iff `trailer & 0x80000000`** (SPR256-TRLR-021). For palette sprites
  bit 31 is set (`0x8000_00NN`), which also makes the trailer's first `u32` an impossible
  width — the **frame-list terminator sentinel** a walker stops at (SPR256-TRLR-016). So
  the frame count *is* stored (COUNT-002 refined) and the high half is a 1-bit palette
  flag, not a `0x8000` word.
- Corpus maxima: width ≤ 640, height ≤ 480, ≤ 256 frames/sprite (SPR256-CORPUS-006); the
  trailer's low `u16` spans the same `[1..256]`.

## Frame pixel data — RLE (`SPR256-RLE-007`, 23 791/23 791 frames)

`data[dataSize]` is a run-length stream of single-byte controls, `[2-bit opcode | 6-bit
count]`. A cursor moves left→right and wraps to the next row at `width`; the background
is transparent.

```
c & 0xC0 == 0x00   (0x00–0x3F)  literal      : emit next (c & 0x3F) bytes as indices
c & 0xC0 == 0x40   (0x40–0x7F)  blank rows   : emit (c & 0x3F) fully-transparent rows
c & 0xC0 == 0x80   (0x80–0xBF)  transparent  : emit (c & 0x3F) transparent pixels
c & 0xC0 == 0xC0   (0xC0–0xFF)  == 0x80 (loader aliases it); unused in data  RLE-009/020
```

The loader (`rom.exe`, `SPR256-RLE-020`) dispatches on `c & 0xC0` with a three-way
`if(0x00)/else if(0x40)/else`; the `else` handles **both `0x80` and `0xC0`** — so `0xC0`
is a decode-time **alias of `0x80`** (transparent skip), not a fourth opcode
(SPR256-RLE-020). Data never emits it (0/23 895 frames, RLE-009). Verified in two
structurally-opposite blit variants (`FUN_0x0044dd40` stencil, `FUN_0x0044dde0` colour);
the `& 0xc03f` decode idiom is shared by 14 blitters. The loader also confirms this exact
opcode assignment (SPR256-RLE-022).

- **Exact / lossless** (SPR256-RLE-008): every row's tokens sum to exactly `width`,
  the stream yields exactly `height` rows and is consumed with no leftover byte, and
  `Σ literal + Σ transparent + Σ blankRows·width == Σ width·height` corpus-wide.
- Blank-row opcodes occur only at a row boundary (column 0).
- The grammar is **corpus-universal** (SPR256-RLE-010): all 14 non-standard payloads'
  frames decode identically; a Bucket-B trailing section follows a cleanly-decoded
  frame (its own structure is still open).
- "Transparent" here = *no pixel emitted*.

## Palette (`SPR256-PAL-011`)

The 1024-byte leading region is 256 entries of **`[B, G, R, reserved]`** (Win32
`RGBQUAD`; byte order **BGR**, 4th byte 0). Byte order is pinned against the game's own
copy of the canonical IBM VGA 16-colour palette in `cursors/default.256`, which matches
the standard index→colour assignments (1=blue, 4=red, 6=brown, …) **only under BGR**
(SPR256-PAL-011). The 4th byte is 0 in 99.96 % of entries; the non-zero remainder
(138/350 720) is **entirely** in `cursors/pickup.256` (95) + `cursors/attack.256` (43) —
the OVL-015 anomaly — where the leading 1024 B is not a standard palette (SPR256-PAL-018).

- **Index 0 is reserved, and it is NOT a colour key** — never emitted as a literal pixel
  across 9.54 M literals (SPR256-PAL-012), and the blitters do not test it: neither
  `FUN_0044db00` (.256) nor `FUN_00451ae0` (.16a) contains a compare or a branch on the
  source value, so every literal is written unconditionally (SPR256-KEY-044). All
  transparency is the RLE skip/blank opcodes. Treating index 0 as transparent changes no
  shipped pixel and diverges on the first sheet a customisation adds.
- **Zeroed-low-palette sprites (181)** (SPR256-PAL-013): two sub-classes below.

## Overlay layers — `spritesb.256` (`SPR256-OVL-014`)

A `spritesb.256` "b"-variant is an **overlay decal**, not a standalone sprite
(SPR256-OVL-014): **one frame per base `sprites.256` frame** (180/180 pairs match),
~12× sparser (median literal density 0.022 vs 0.266), and its **own low palette is
zero on purpose** — its literal pixels index slots that are black in its own palette
but coloured in the **base sibling's** palette (179/180 ≥99 %, median 100 %). So it is
composited over the base frame **through the base palette**. Rendering it standalone
with its own palette yields sparse black dots — that is expected, not a decode fault.
The engine's exact blend (glow / team-tint / replace) is not yet derived (needs
`rom.exe`/runtime).

- **Open anomaly** (SPR256-OVL-015): `cursors/attack.256` + `cursors/pickup.256` are the
  only zero-low-palette sprites with **no base sibling**; they carry their own partial
  palette (low 16 zero) and dense content, and decode to non-image bands. Cause open.

## The equipment tree — portraits, not world sprites (`SPR256-EQUIP-042`)

`graphics\equipment` is addressed by `rom.exe`'s own literals and is **not** part of the world
sprite path. Four figure directories — `ffighter`, `fmage`, `mfighter`, `mmage` — two layer
directories under each — `primary`, `secondary` — and 987 file nodes, identical between the EN and
RU roots.

```
equipment/<figure>/primary/<7 digits>.256     784 sheets   1 frame each, 160x240
equipment/<figure>/secondary/<7 digits>.256   144 sheets   1 frame each, 160x240
equipment/<figure>/<1..31>.256                 59 sheets   the head, indexed by the face byte
```

Four of the 784 `primary` sheets yield 0 frames. Every seven-digit leaf name is produced by exactly
one `u16` through `formats/item/format.md`'s formatter — 928 of 928 inverted, none unproduced, none
ambiguous, over the entire 65 536-value domain.

The one-frame, 160×240 shape is what separates this tree from a world sheet: the sixteen hero body
sheets under `units/heroes*` carry 129 to 216 frames of 24×40 to 40×48, one per direction and phase.
An equipment sheet cannot be composited into a world frame, and the routine that loads them is a
virtual method of the drawable whose identified callers are the info-window routines
(`formats/hero/format.md` §6b).

## Unit sheet layout (`SPR256-UNIT-024`)

A `.256` is a flat frame array; nothing inside the file says what any frame is. For the 34
`units.reg` classes the sheet is **block-addressed**, and the block bases are arithmetic on the
class's own phase counts, computed fresh at every draw. With `MB, MV, AT, DY, BN, ID` the
class's `MoveBeginPhases, MovePhases, AttackPhases, DyingPhases, BonePhases, IdlePhases`:

```
Flip == 0  ->  S = 16, D = 8            Flip != 0  ->  S = 9, D = 5

  0                       stand    S frames        index = the facing itself (16-way)
  S                       move     D*(MB + MV)     dir*(MB+MV) + MB + moveTL[phase mod len]
  S + D*(MB+MV)           attack   D*AT            dir*AT      + attackTL[phase]
  S + D*(MB+MV+AT)        dying    D*DY            dir*DY      + phase/2
  S + D*(MB+MV+AT+DY)     bone     D*BN            dir*BN      + (stage - 2)
  S + D*(MB+MV+AT+DY)     idle     D*ID            dir*ID      + idleTL[phase]
```

`dir` is 16-way for the standing block (`(facing - 8) & 0xf`) and 8-way everywhere else
(`((facing - 8) >> 1) & 7`). On a `Flip` sheet only directions `0..8` / `0..4` are stored and
the rest are drawn mirrored. **Bone and idle share one base** — no shipped class carries both.
The dying and bone blocks are read out of `classes[Dying]`'s sheet with *its* counts and *its*
`Flip` (`REG-UNITS-050`).

Sheet length is therefore `S + D*(MB+MV+AT+DY+max(BN,ID))`, which reproduces the frame count of
**33 of the 34 shipped unit sheets exactly** and overshoots none. (The 34th, `Unit2`
*Unarmed Fighter with Shield*, spells a `Flip = 1` sheet but inherits `Flip = 0`; it is a defect
in the registry, and four shipped map placements reach it.)

**Every one of the 34 unit sheets mixes frame sizes**, so the draw anchor must be built from the
frame actually being drawn (`SPR256-FRAME-023`, `TERR-SPR-043`).

## Structure sheet layout (`SPR256-STR-040`, `SPR256-STR-041`)

A `structures.reg` class opens **two** sheets, `graphics\structures\<File>.256` and
`<File>b.256` — the `spritesb` overlay class of `SPR256-OVL-014`, drawn with the same frame
index through `vt+0x34` and gated on `[0x005eb520]`.

Neither sheet is block-addressed like a unit's and neither is `Index`-addressed like an
object's. It is a **row-major grid of tile-sized frames**, `TileWidth` columns wide and
**`FullHeight`** rows tall — `FullHeight`, not `TileHeight`, so the grid is taller than the
footprint by the structure's overhang. **All 130 structure sheets in the install carry
uniformly `32×32` frames**, which is why no anchor pixel is needed: a frame *is* a map cell.

```
cellIndex(k, c) = k * TileWidth + c          k = image row 0..FullHeight-1
                                             c = image column 0..TileWidth-1

  [ 0                          , TW*FH )         BASE    frame = cellIndex
  [ TW*FH                       , TW*FH + P*L )  ANIM    frame = TW*FH + (p-1)*L + rank(c,k)
  [ frames - TW*FH              , frames )       RUIN    frame = frames - TW*FH + cellIndex

  TW = TileWidth   FH = FullHeight   L = AnimMask's non-'-' count   P = max timeline value
  p    = the class's current animation value (0 = not animating this frame)
  rank = how many non-'-' AnimMask cells precede this one
```

The ruin block is addressed **from the end of the file**; nothing stores its base. Which
block a cell takes is decided by `TERR-STRUCT-102`. Shipped identity:
`frames = k·TW·FH + P·L`, `k ∈ {1,2}`, on **64 of 66** classes — `k = 2` on 39 (a ruin grid
ships), `k = 1` on 25 (none does, and `frames − TW·FH` is then 0, so a destroyed structure
re-draws its intact art). The two exceptions are `bridge1v` and `bridge2`, whose drawable is
a C++ subclass with a 9- and a 14-frame nine-patch selector (`TERR-STRUCT-105`).

## Known variants / exceptions

- **No-palette variant** (6 projectile arrows): the 1024-byte palette is absent;
  frames begin at offset 0. This is the **same format with the trailer's has-palette bit
  clear** (SPR256-TRLR-021) — the loader simply skips the palette read; bit31=0 so the
  trailer is a plain `u32` frame count (16). The external palette source is still TBD
  (SPR256-VAR-004).
- **Frame + trailing section** (8 UI / global-map / goblin-arrow sprites): a normal
  palette + one frame + its count-trailer (`0x80000001`), then an **appended secondary
  section bracketed `0x80000001 … 0x80000001`** (8/8). The inner block is opaque — not a
  frame, not a same-width RLE image; 4–16 B (small UI/projectile) or 303–2521 B dense
  (world-map tiles + shop frame, using `0xC0-0xFF` bytes). **The shipped loader ignores
  it**: the `CSprite256` ctor indexes exactly `frameCount = 1` frames (the final
  `0x80000001`), so the inner is loaded but never indexed/decoded/blitted — runtime-inert
  (SPR256-EXC-020). Authoring origin still open but non-runtime. (SPR256-EXC-005/017)

## Open (own experiments)

- **Overlay blend rule** — `spritesb` is a base-palette overlay (`SPR256-OVL-014`); the exact
  runtime composite (additive glow / team-tint / replace) needs `rom.exe`/runtime.
- **attack/pickup banded anomaly** (`SPR256-OVL-015`) — 2 files, cause open; the
  reserved-byte census establishes that these are the *only* `.256` with a non-zero
  reserved palette byte
  (SPR256-PAL-018), so their 1024 B region is not a palette. Cause still needs `rom.exe`.
- ~~**Frame-block / animation roles**~~ — *closed for units by `SPR256-UNIT-024`*
  and **for `structures.reg` by `SPR256-STR-040`, `SPR256-STR-041`** (above: a `TileWidth ×
  FullHeight` grid of 32×32 frames in three blocks, identity 64/66 with both exceptions
  predicted). `objects.reg` needs no block map: an object sheet is a flat list its class
  indexes by `Index`, `Index + timeline[phase]` or `0` (`TERR-SPR-042`, `REG-OBJ-046`).
- **No-palette variant's** external palette source remains open (`SPR256-VAR-004`).

## Closed

- **The loader in `rom.exe`** (`SPR256-TRLR-021`, `SPR256-RLE-020`,
  `SPR256-EXC-020`) — `CSprite256` ctor `FUN_0x00428c50` (vtable
  `0x00597418`) parses the file; draw dispatch (`+0x40` = `FUN_0x00429190`) feeds a family
  of ~42 DirectDraw blitters (`0x0044c660`‥`0x00452xxx`) that decode the RLE. This is the
  authority behind the three items below.
- **RLE `0xC0` quadrant** — the loader aliases `0xC0` to `0x80` (transparent skip); no 4th
  opcode. Data never emits it (0/23 895). Resolved (`SPR256-RLE-019`, `SPR256-RLE-020`).
- **Bucket-B inner block** — the shipped loader loads it but never indexes/decodes/blits
  it (frameCount = 1 from the final trailer); runtime-inert. Resolved
  (`SPR256-EXC-017`, `SPR256-EXC-020`); authoring origin non-runtime, out of scope.
- **Frame count + trailer** — the count is walk-derived (SPR256-COUNT-002) **and** stored
  as `trailer & 0x7FFFFFFF`, the two agreeing exactly (1370/1370); the trailer =
  `[31-bit frameCount][bit31 = has-palette]`, and bit31 both gates the palette read and
  acts as the end-of-frames sentinel (`SPR256-TRLR-016`, `SPR256-TRLR-021`). The no-palette
  arrows are the same format with bit31 = 0.
- **RLE decoding** — grammar established and proven lossless (`SPR256-RLE-007`,
  `SPR256-RLE-008`);
  the per-row = `width` invariant holds for all 23 791 frames.
- **Palette colour + transparency** — entry `[B,G,R,0]`, byte order BGR pinned to the
  VGA palette, index 0 reserved but **not** a colour key (`SPR256-PAL-011`, `SPR256-PAL-012`,
  `SPR256-KEY-044`); the reserved 4th byte is non-zero only in the 2 anomaly cursors
  (`SPR256-PAL-018`).
- **`spritesb` overlay class** — a base-palette overlay, 1 frame per base frame
  (`SPR256-OVL-014`); only 2 zero-palette files remain anomalous.
