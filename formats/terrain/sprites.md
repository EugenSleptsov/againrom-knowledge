# Sprite placement and composition

[Reference](format.md)

## Sprite placement — where a unit or object stands (`TERR-SPR-038…TERR-SPR-041`)

Sprite placement applies one vertical terrain lift to the frame; sprite
blitters do not use the sloped-terrain step table. The draw call carries
scalar placement rather than four corner heights. TERR-SPR-038's fifth-
argument label is amended: this argument is the shadow's per-row X slope,
not a brightness level. — TERR-SPR-038, TERR-SPR-066

```
alt(col,row) = ( h(wc,wr) + h(wc+1,wr) + h(wc,wr+1) + h(wc+1,wr+1) ) / 4     trunc toward zero
               four MOVSX reads of the raw type2 grid; wc = col+scrollX, wr = row+scrollY
               built by FUN_00407117 pass 3 into CMapView+0xc0        (TERR-SPR-039)

anchorX = (CenterX - Width /2) + frameWidth /2      class +0x20, +0x18   (TERR-SPR-040)
anchorY = (CenterY - Height/2) + frameHeight/2      class +0x24, +0x1c
    frameWidth/frameHeight are the BODY pass's own drawn frame; the SHADOW pass measures
    frame 0 instead, and frames of one .256 need not share a size  (TERR-SPR-043)

dstX = col*32 + 16 - anchorX          ( - T on the SHADOW pass only, see below )
dstY = row*32 + 16 - anchorY - alt

dst is the frame's TOP-LEFT. With frameW=Width and frameH=Height this reduces to:
    frame pixel (CenterX, CenterY)  lands on  (col*32 + 16, row*32 + 16 - alt)
```

**The shadow pass's extra `dstX` term is a pixel count, not the slope.** Two sun-derived values
travel together and they are not interchangeable (`TERR-SHDW-130`, and `TERR-SPR-067`'s naming of
one as the other is amended in `retracted.md`):

```
s = ftol( tan(theta) * 65536.0 )                         the 16.16 per-row X slope; blit arg 5
T = ftol( tan(theta) * (floor(frameH/2) + floor(Height/2) - CenterY) )      a PIXEL count
      the multiplicand is frameH - anchorY for even frameH, one less for odd:
      the two /2 are separate truncations
```

Only `T` enters `dstX`. `s` goes to the blitter, which shears about the **bottom edge of the blit
rectangle** all by itself; subtracting `T` moves that pivot up to the **anchor row**, so that the
composed placement of image row `r` is

```
X(r) = (col*32 + 16 - anchorX) + tan(theta) * (anchorY - r)        residue <= 2 px
```

A structure writes the same rule upward instead of downward and therefore **adds** where a sprite
subtracts — `dstX = col*32 + ftol(tan(theta)*((FullHeight-k)*32 - ShadowY))` — because
`(FullHeight-k)*32` is a height above the image bottom while `anchorY` is a depth below the frame
top. The two signs are both correct and neither may be copied onto the other path
(`TERR-SHDW-131`).

**Two altitude models per frame.** The terrain raster uses `+0xb4` — one *corner*, `r*32 − h`
(`TERR-GEOM-031`). Everything standing on it uses `+0xc0` — the *mean of four corners*. A port must
implement both; using the terrain mesh to lift a unit puts it at a corner's height instead of the
quad centre's.

The four-corner mean equals a corner height on flat cells and can also
equal it on sloped cells. Flatness therefore cannot be tested from that
equality. The lift uses the four named samples and truncation toward zero.
Installed contact points can lie up to two pixels below their own terrain
span at the cell center, within pixels drawn by neighboring cells.
— TERR-SPR-039

The occlusion test `FUN_0040d027(classId, col, row, alt)` recomputes the same `dstY`, converts the
sprite's top and bottom back to terrain rows with the picker `FUN_0041a9f9`, and skips the draw iff
every covered cell's shroud level is `0x10` (fully dark).

### Which frame a static object draws (`TERR-SPR-042`, `TERR-SPR-043`, `TERR-TILE-044`)

The `frame` argument of the object draw call is not simply the class's `Index`. `FUN_00407b1a`'s
type-3 pass writes it on three arms and then overrides it once (`TERR-SPR-042`):

```
c    = type3[row*W + col]                     0 -> nothing here
k    = classes[c - 1]                         objects.reg section index   (ALM-CLS-035)
imp  = tile[row*W + col] & 0x2000             bit 13 -- see below
anim = ( (tile[i] | tile[i+1] | tile[i+W] | tile[i+W+1]) & 0xc000 ) == 0xc000

if k.timelineLen != 0 and anim and not (imp and k.DeadObject != -1):
      phase = (animCtr + col*(row+1)) % k.timelineLen        # SIGNED remainder
      frame = k.Index + k.timeline[phase]                    # timeline = class+0x2c
elif imp and k.DeadObject != -1:
      k     = classes[k.DeadObject]           # the whole sprite is swapped
      frame = 0
else: frame = k.Index

if animations_disabled:  frame = 0            # DAT_005bcef0 == 0, -noanimation / -detail0
sprite = Files[k.File]                        # graphics.res!objects/<path>.256
```

`k.timeline` is not a registry key: the loader run-length expands `AnimationFrame[i]`
`AnimationTime[i]` times into `class+0x2c` and stores the length at `class+0x3c`
(`REG-OBJ-046`). `animCtr` is `CMapView+0xa70`, the same counter the water phase uses,
unshifted here.

**The two flag bits of the tile word** (`TERR-TILE-044`). Bit 13 is the one `TERR-DIRT-017`
already uses for the dirt composite; the object pass reads the same bit for the `DeadObject`
swap, and `FUN_0041f3d4` reads it to decide whether a cell still holds a standing destructible.
It is **rewritten at runtime** by an RLE decoder (`FUN_0041f2b7`) over the interior from `(8,8)`,
so the `.alm` field supplies only its initial state.

See [fog and visibility](fog.md) for tile bits 14/15, and
[sprite composition](composition.md) for frame passes and unit shadows.

### Sprite lighting — how bright a unit or object is drawn (`TERR-LIGHT-059…TERR-LIGHT-064`)

A sprite **is** lit, per sprite rather than per pixel, through the same builder and the same kind of
level-indexed LUT the terrain uses. The `.256` class dispatches on vtable `0x00597418` (stored by its
own payload constructor `FUN_00428c50`) and exposes four blit entry points, which split in two:

| slot | routine | pixel loops | reads | role |
|---|---|---|---|---|
| `vt+0x14` | `FUN_00428f40` | `0044db00`, `0044ea40` | the **source** index | body, plain |
| `vt+0x34` | `FUN_00428fc0` | `0044dde0`, `0044eca0` | source + destination + `[0x005eb578]` | body, `spritesb` overlay |
| `vt+0x1c` | `FUN_00429040` | `0044e010/e240/eef0/f140` | the **destination** only | silhouette (shadow) |
| `vt+0x3c` | `FUN_00429100` | `0044e460/e750/f380/f690` | the **destination** only | silhouette (shadow), sheared |

The lit pair takes `(dstX, dstY, frame, level, shadeObj, mirror)` and forms

```
row = shadeObj[+8] + (level << 9)          # 512-B row stride, 256 u16 entries
dst = row[srcIndex]                        # 0044dbed MOV AL,[ESI] ; 0044dbf8 MOV AX,[EBX+EAX*2]
```

The silhouette pair takes no shading object: it advances the source without reading it, reads the
destination pixel, `SHR 3`, and looks *that* up in the shroud table `[0x005e8420]` at row
`arg4 × [0x005e42f0]` — so those two passes recolour what is already on screen under the sprite's
outline. `vt+0x3c`'s fifth argument is the shadow's **16.16 per-row X slope** from the sun angle, not
a brightness (`TERR-SPR-066`, correcting `TERR-SPR-038`); `vt+0x1c`'s caller instead offsets `dstX`
by `shear/2000`.

**Which pass is which.** An object cell issues `vt+0x3c` twice (shadow) then `vt+0x14` + `vt+0x34`
(body). A unit is dispatched twice over the drawable grid `CMapView+0x90`, both times at
`idx = (row+3)*(visCols+6) + (col+3)` under `drawable+0x78 == 0`: `vt+0x2c` = `FUN_0045bf00` gets the
`+0xc0` altitude and draws the **shadow**, `vt+0x28` = `FUN_0045b3f0` gets the `+0xb0` light level and
draws the **body** (`TERR-SPR-065`, correcting `TERR-SPR-048`, partially retracted).

**The level.** `CMapView+0xb0` is a per-frame byte grid of `(visCols+6) × (visRows+10)`, filled by
`FUN_004050da`:

```
memset(grid, DAT_005eb494 >> 2, size)                 # ambient 0x0e by day -> 3, everywhere
for each entry of the lit-object list this+0x9f0:
    flags & 0x1000  -> grid[cell] = 0                 # brightest
    flags & 0x20000 -> grid[cell] = 0xc               # x0.5
    flags & 0x8     -> flicker 0 / 0xc, + a radius-1 stamp into +0xa8
for every cell:                                       # +0xa8 = light-source stamps, 0xff = none
    q[k] = max(0, a8corner[k] - 0x20)  (0xdf -> ambient)
    if all four corners were 0xff: keep the value above
    else grid[cell] = min( (q0+q1+q2+q3) >> 4 , ambient >> 2 )
```

Lower level = brighter, so the sweep can only brighten; the only darkening mechanism is a
`0x20000`-flagged object. Objects and units read the same grid, so nothing can light one differently
from the other.

**Mode 2 — the table.** `FUN_00427df0`'s arm at `0x00428054` (the dispatch table sits at
`0x00428484`), built as `(nLevels = 0x10, mode = 2, useTint = 1)`:

```
out_chan = clamp( ((palette_chan + skyTint_chan) * (nLevels - level) * 2) / nLevels, 0, 255 )
```

so **row `L` has gain `2(nLevels − L)/nLevels`**: row 0 → ×2.0, row `nLevels/2` → ×1.0 (identical to
no table at all), row 15 → ×0.125.

**The two ladders are one ladder.** `2(16 − S)/16 = (96 − (4S + 32))/32` exactly, and on the same
palette the rows are bit-identical (0 of 256 entries differ at `L = 32, 36 … 64`). Shipped daytime:
ground `L = 46` → ×1.5625, sprite `S = 3` → ×1.6250 — half a step apart because `ambient/4 = 3.5`
truncates to 3. Both indices descend from the same byte `DAT_005eb494`, whose only two readers in the
image are `FUN_004050da` and `FUN_00484f40`.

**Which table each drawable gets.** Terrain: one, from the tile palette. Each object sheet: its own
`(0x10, 2, 1)` table from its own BMP palette (`FUN_0046b7e0`), passed as `spriteObj + 0x14`. A unit:
selected on the `units.reg` `Palette` key at `class+0x98`: `0` selects one of
16 shared owner tables from `human.pal`; `1` selects `class+0x9c[0]`; `> 1`
selects `class+0x9c[face-1]`, where face is the actor's Data.bin tier. The last
arm is not owner-indexed (`PAL-KEY-002`, `PAL-OWN-007`, correcting that clause
of `TERR-LIGHT-064`). The owner-index producer and first-message lifetime are
specified by `PAL-RULE-021` and `PAL-FIRST-022`; see [palette selection](../pal/format.md).
`MAGIC-STONEDRAW-084` identifies both greyscale overrides as the stone-curse
draw arm; its stated confidence and the `drawable+0x15a` Unknown remain.
The terrain.3d tiles share one palette; sprite ramps are built from each
sprite's own palette.

**Reimplementation guidance.** Light every sprite through its own palette's 16-row mode-2 table at the
cell's `+0xb0` level. Drawing sprites unlit is not a small error: the shipped swordsman frame means
104.9 luma at the engine's row 3 and 65.2 with no table, against 106.1–118.3 for the terrain in the
same crop (`tools/terrgfx -mode sprlightart`).

### Terrain output range (`TERR-LIGHT-063`)

For the shared terrain.3d palette and installed per-vertex levels 30..70 at
θ=0.78539815, mode 3 gives the following output ranges:

| level | gain | entries with a clipped channel | entries packing to `0xffff` | pixels clipped | pixels pure white | mean luma |
|---:|---:|---:|---:|---:|---:|---:|
| 30 (brightest shipped) | ×2.0625 | 89 | 33 | 16.99 % | 5.09 % | 153.25 |
| 46 (flat, daytime) | ×1.5625 | 47 | 11 | 10.59 % | 0.10 % | 119.49 |
| 64 (unattenuated) | ×1.0000 | 8 | 1 | 0.30 % | 0.00 % | 79.48 |
| 70 (darkest shipped) | ×0.8125 | 0 | 0 | 0 % | 0 % | 64.14 |

**The mapping saturates**: bright shipped terrain reaching pure white is what the engine does, on
about a twentieth of all tile pixels at the brightest level a shipped map produces. Pixel weights are
over 651 266 `terrain.3d` pixels; the level range carries its θ, framing base and window because the
ceiling moves with all three (`TERR-LIGHT-029`).
