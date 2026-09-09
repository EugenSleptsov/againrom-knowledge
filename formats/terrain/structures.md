# Structure placement

[Reference](format.md)

## Structure placement — where a building stands (`TERR-STRUCT-100…TERR-STRUCT-107`)

A `structures.reg` placement is **not** the unit/object shape above. It has no canvas and no
anchor pixel; its anchor is a **cell** and every frame of its sheet is exactly one tile
(`TERR-STRUCT-100`, `SPR256-STR-040`).

```
anchor cell (ac, ar) = (fine.x >> 8, fine.y >> 8)      installed fine = cellByte*256 + 128
                                                        i.e. the footprint's min col / min row

drawn once per footprint cell (c, r), c in [ac, ac+TileWidth), r in [ar, ar+TileHeight):
    COL0   = c - ac                                     local column
    ROW0   = r - ar                                     local row
    rowTop = ROW0 - TileHeight + FullHeight
    limit  = (ROW0 != 0) ? rowTop : 0                   the BACK row also draws the overhang

    for k = rowTop down to limit:
        frame = block(k*TileWidth + COL0)               see SPR256-STR-041
        dstX  = (c - scrollX) * 32
        dstY  = (r - scrollY) * 32 - lift - (k_first - k)*32       ( - obj+0x10, always 0 )
```

`lift` is the drawable's own `+0x68`: a **bilinear** sample of the cell's four corner heights
at the object's sub-cell position (`TERR-STRUCT-106`) — which at a cell centre, where every
structure sits, equals the four-corner mean the unit path uses. `dst` is the frame's
top-left, as above.

So a `TileWidth × TileHeight` building draws `TileWidth × FullHeight` tiles: the bottom
`TileHeight` rows sit on the footprint, the top `FullHeight − TileHeight` rows hang above the
back row.

**Shadow** (`TERR-STRUCT-103`, `TERR-SHDW-136`): the same loop and the same `dstY`, through the
sprite's `vt+0x3c` with `[0x005eb49c]`, and `dstX` displaced **per strip** by
`ftol( tan(sunAngle) · ((FullHeight − k)·32 − ShadowY) )`. That is not the only shear: the call
also passes the **live 16.16 slope** into blit argument 5 (`0045a3ff`/`0045a407`), so each strip's
rows are sheared again inside the blit. The two compose into one pivot — `X = col*32` at the screen
row `dstY_k + h` of the strip with `(FullHeight − k)·32 = ShadowY`, `h` being the strip frame's
height. Wooden bridges override the shadow with a `RET`: they cast none.

**`ShadowY` is also a suppression sentinel** (`TERR-SHDW-136`). It is read by one instruction,
`0045a3cb FILD dword ptr [ESI + 0x30]`, and no compare in the routine or the loader touches it.
Over the 66 shipped classes it is 28…55 on 51 of them and then `10000` on four and `20000` on
eleven, against a maximum `FullHeight·32` of 192 px. Those fifteen push the strip 660–11 400 px
sideways, past any viewport, so the blitter rejects it: **they cast no shadow, and the mechanism is
displacement rather than a branch.** They are the flat and the hollow — four bridges, four graves,
a cave, three wells, a teleport, a campfire, `magic`. A consumer must not clamp the field.

**Order** (`TERR-STRUCT-104`, amended): the drawable registers its prepared
footprint into `CMapView+0x94`. The early body phase takes `Flat!=0`; the main
cell phase takes `Flat==0` before the selector-2 CUnit calls in that cell.
Both body gates require signed `+0x78<2`. The former main-body `==0` clause is
retracted: that later test guards work after the body call. Structure shadows
have an earlier separate phase (`ANIM-DRAWGATE-087`, `ANIM-CELL-085`).

**Ownership never reaches the sprite** — only the minimap blip `vt+0x34` (`TERR-STRUCT-107`).
