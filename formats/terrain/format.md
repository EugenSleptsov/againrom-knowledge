# Terrain-tile graphics (`terrain.3d/*.bmp`) — specification

Level 3. Promoted, evidence-backed claims only. Sources: `TERR-LOC-001`,
`TERR-LOAD-002`, `TERR-IDX-003`, `TERR-SEM-004`, `TERR-VER-005`
(format, location and render mapping from `graphics.res` bytes and `rom.exe`'s
load/render routines), building on `ALM-GRID-013`, `ALM-TERR-015` and `TERR-SEM-004`
(tile-word semantics and terrain enum). Corpus: the 53 `terrain.3d/*` BMPs + the type1
grids of the 38 maps (880 704 cells — all of them; `TERR-VER-005` retires the
superseded "overlay corner" exclusion); 0 falsification violations.

This is the graphic counterpart to the ALM **type1 "Tiles"** grid: `formats/alm`
(`ALM-GRID-012`) decodes what a tile-word *means* for the sim; this doc decodes what it
*draws*.

> **Where the tile words come from (`TERR-GRID-027`).** Everything
> below is a statement about a tile *word* and about the `W·H` buffer the engine loads it
> into — none of it changed. What changed is **which file bytes fill that buffer**: the
> `.alm` type-1 payload starts 8 bytes later than `ALM-GRID-011` said, so cell `(col,row)`
> is word `row·W + col` of the record payload = word `row·W + col + 4` of the payload as
> the superseded `ALM-GRID-011` framing delimited it. **A terrain frame rendered on the old base is four cells off in X**
> (and out of register with a heights/overlay layer read the old way, which is off by
> eight). See `formats/alm` → "Grid layers".

## Location & format

The terrain images live inside **`graphics.res`** (the `RES-MAGIC-001` `&YA1` archive) under the
node prefix **`terrain.3d/`**. They are **standard 8-bpp Windows BMP** files (magic `BM`,
14-byte `BITMAPFILEHEADER` + 40-byte `BITMAPINFOHEADER`, 256-entry RGBQUAD palette,
`bfOffBits = 1078 = 14 + 40 + 1024`) — read with any standard BMP decoder; ROM adds no
wrapper. (`world.res` holds only terrain *config* — `data/map.reg` — not the images.)

| File(s) | Count | BMP W×H | 32×32 sub-cells | Role |
|---------|-------|---------|------|------|
| `dirt.bmp` | 1 | 32×128 | 4 | impassable-tile composite overlay |
| `tile1-00..15` | 16 | 32×448 | 14 | Land group (Grass / Cracked / Sand / Savanna) |
| `tile2-00..15` | 16 | 32×448 | 14 | Stones / Cracked-Stones / Flowers-Savanna / Mountain |
| `tile3-00..15` | 16 | 32×256 | 8 | **Water** (animated) |
| `tile4-00..03` | 4 | 32×448 | 14 | **Road** |

Each BMP is **32 px wide** and is a **vertical strip of 32×32 sub-cells** (top-to-bottom;
BMP rows are bottom-up per the standard, so sub-cell `k` occupies the *k*-th 32×32 block
counting from the strip origin the game reads — pixel byte offset `8 + k·0x400` into the
loaded image buffer). There is a parallel **legacy non-3d** set at `terrain\tile*.bmp`
that `rom.exe` can load instead; the flag `DAT_005e4418 & 2` selects the `terrain.3d` set,
which is the shipped path.

## Loader (`rom.exe FUN_00469620`)

The only referrer of the terrain path strings. It builds `tileG-VV.bmp` (`%d-%d` /
zero-padded `%d-0%d` for `V<10`) and loads each into a **128-slot pointer array `tiles[]`**
at `DAT_005ef6c0`:

```
tiles[(G-1)*16 + V]  <-  tileG-VV.bmp      G = 1..8, V = 0..15   (absent file -> null)
DAT_005ef8d4         <-  dirt.bmp
DAT_005ef8d8         =   (first non-null tile) + 0x14            (default handle)
```

On disk only `tile1/2/3` (16 each) + `tile4` (0..3) exist, so slots for groups 5–8 and
`tile4` V≥4 stay null. The loader then pre-composes the tiles into 256×256 DirectDraw work
surfaces (`DAT_005eb5a8[]`) and a shared palette (built from the tiles' own palette) for
8-bpp display.

**The loader takes one argument, and the map supplies it.** `FUN_00469620` is `cdecl` with a
single dword, read at `[esp+0x558]` and used as a 32-bit mask: the body runs exactly 32
iterations over `tiles[]` in steps of four slots, and enters iteration `i` only when
`mask & (1<<i)`. Iteration `i` fills `tiles[4i .. 4i+3]`, which is precisely the block the
render mapping reaches with strip group `g = i` and blend column `b = 0..3`, so one mask bit
is one value of `g`, i.e. one file group `G = (i>>2)+1` with variants `V = (i&3)*4 .. +3`.
The value comes from the `.alm` type-0 record: payload `+0x18` → landscape slot `P+0x28`,
pushed at the end of the map-build. Over 72 shipped maps it takes eight values, all
`≤ 0x1fff` — `0x1fff`×44, `0x0fff`×10, `0x00ff`×4, `0x0fbf`×4, `0x1fdf`×4, `0x0f7f`×2,
`0x1f1f`×2, `0x1fbf`×2 — never a bit above 12, which decodes to groups 1–3 entire plus
group 4 variants 0–3: the same set the disk holds. A clear bit leaves all four images that
tile-word group can select null (`TERR-LOAD-152`).

## Tile-word → graphic (render mapping)

Four render routines — `FUN_004058c7`, `FUN_00405e83`, `FUN_00406349`, `FUN_00406806`
(zoom/mode variants) — read the type1 tile-word `w` (u16) and select image + sub-cell with
**identical** arithmetic:

```
g    = (w & 0x1fff) >> 6          strip group      (bits 6..9; corpus range 0..12)
b    = (w >> 4) & 3               blend column     (bits 4..5)
sub  = w & 0xf                    sub-cell index    (bits 0..3)

image = tiles[ g*4 + b ]          == tileG-VV.bmp with:
          G = (g >> 2) + 1        file group  (1..4)
          V = (g & 3)*4 + b       file variant (00..15)
src   = image.pixels + 8 + sub*0x400        one 32x32 8-bpp cell (0x400 = 1024 B)
```

Mapping of strip group → file group → terrain (terrain names via `ALM-TERR-015`):

```
g  0  tile1-00..03  Grass  / Land        g  4  tile2-00..03  Stones          / Land
g  1  tile1-04..07  Cracked/ Land        g  5  tile2-04..07  Cracked         / Stones
g  2  tile1-08..11  Sand   / Land        g  6  tile2-08..11  Flowers         / Savanna
g  3  tile1-12..15  Savanna/ Land        g  7  tile2-12..15  Mountain        / Stones
g 8..11  tile3-*    Water (animated)     g 12  tile4-00..03  Road            / Land
```

- **Water (`g 8..11` = `tile3`)** is animated (fully decoded by `TERR-ANIM-006…TERR-ANIM-010`; see "Water
  animation" below). The renderer replaces the group with `8 + phase`,
  `phase = (g + (worldCol+1)*worldRow + (animCtr>>2)) & 3` when enabled (`DAT_005bcef0 ≠ 0`),
  else `phase = 0`. Water uses only 8 sub-cells (tile3 is 32×256).
- **Impassable tiles (bit 13 set, non-water)** are composited over `dirt.bmp` before blit
  (`dirt` sub-cell `= (col + row*5) & 3`), giving the "blocked ground" look. The composite is a
  **transparent-index (index 0) keyed overlay**, not a blend: copy the terrain sub-cell, then
  paint the dirt sub-cell's non-zero pixels over it (`FUN_0044c390`; `TERR-DIRT-017`). Drawn shaded.
- The four **corner altitudes** (type2 Altitudes grid, `ALM-GRID-013`) choose a flat vs
  sloped blit **and displace the pixels** — see "Cell geometry" below. A render that ignores
  them is wrong on **89.7 %** of the shipped corpus's cells.
- Every tile is drawn **relief-shaded** — see "Terrain lighting" below.

## Tile word → movement (`TERR-PASS-049…TERR-PASS-051`, `TERR-COST-052`, `TERR-PASS-053`)

The same word drives a second, wholly separate reading: whether a **ground unit** may stand on the
cell. It is the sim side, not the render side, and it uses a different bit split — the whole
decision runs on `w & 0x3ff` plus bit 13, so **bits 10–12 and 14–15 never reach it** (0 of 880 704
shipped cells set any of them). Sim-side grid layout, the ingest and the persistence are specified
in `formats/alm` → "Runtime passability"; what belongs here is the tile word's own part.

`rom.exe FUN_00548720` maps the masked word to a terrain class and a movement cost:

```
i = w & 0x3ff

if (i & 0x300) == 0x200:                     # the water range, taken before anything else
    if (i & 0xf) >= 8:      return 0xff, 8   # reject      (0 shipped cells)
    if (i & 0xf) == 4 and (i & 0x30) == 0x10:
                            return 1,    8   # Land        (0 shipped cells)
                            return 9,    8   # Water; the cost is the literal 8, never CostWater

s = i & 0xf ; b = (i >> 4) & 3 ; g = (i >> 6) & 0xf     # same split as the render mapping
if s >= 14:                 return 0xff, 8   # reject      (0 shipped cells)

pri, sec = pair[g]                           # world+0x54156 + 2g   (immediates, below)
sel      = level[b][s]                       # world+0x540d6 + 16b + s, in 1..5
cls      = sec if sel <= 2 else pri
cost     = [ cost[sec],
             (3*cost[sec] + cost[pri]) >> 2,
             (  cost[sec] +   cost[pri]) >> 1,
             (3*cost[pri] + cost[sec]) >> 2,
             cost[pri] ][sel - 1]            # cost[k] = world+0x54176+k, [0] = 0xff
return cls, cost
```

`pair[g]` is the same primary/secondary pairing the render mapping's terrain names come from,
which is what makes this an independent attestation of them from the movement side:

```
g  0 (Grass, Land)   g  4 (Stones, Land)     g  8..11  never written (water early-out)
g  1 (Cracked, Land) g  5 (Cracked, Stones)  g 12 (Road, Land)
g  2 (Sand, Land)    g  6 (Flowers, Savanna) g 13..15  never written — AND REACHABLE:
g  3 (Savanna, Land) g  7 (Mountain, Stones)           a map with g >= 13 reads uninitialised
                                                       heap. 0 shipped cells do it.

level[b][s], s = 0..13            cost[1..10] from world.res:data/map.reg  (Cost only)
b=0 [2,3,2,4,3,4,2,2,2,2,4,4,4,4]   Land 8 Grass 8 Flowers 8 Sand 14 Cracked 6
b=1 [3,5,3,3,1,3,2,4,2,2,4,2,4,4]   Stones 12 Savanna 8 Mountain 16 Water 8 Road 6
b=2 [2,3,2,4,3,4,2,4,2,2,4,2,4,4]
b=3 [5,5,5,5,5,5,2,2,2,2,4,4,4,4]
```

A cell is then blocked for a **ground** mover by any of three tile-word arms — `w & 0x2000`
(7 464 cells), `cls == 8` i.e. Mountain (136 622), and the raw water test `(w & 0x300) == 0x200`
(87 584), pairwise overlaps 860 / 370 / 0 — plus a nonzero type-3 cell and the 8-cell border.
Note the two roles bit 13 plays are the same bit but not the same rule: the render path composites
it over `dirt.bmp` (above), the sim path is one of three blockers, and it is the **smallest** of
them. Class census over the 38 maps: Land 233 611, Grass 104 421, Flowers 10 685, Sand 24 119,
Cracked 41 668, Stones 92 880, Savanna 118 009, Mountain 136 622, Water 87 584, Road 31 105, and
the reject value `0xff` **0** times.

### Who is asking — the mover (`TERR-MOVE-054…TERR-MOVE-057`)

The `n×n` footprint and the mask are **not** properties of a `units.reg` class. Both come from the
simulation actor instance the predicate is called on (base vtable `0x59c3c0`, derived `0x59c448`,
`0x59c4d0` — the only family that allocates the `0xb4`-byte mover and stores it at `+0x154`):

```
vt+0x1c  FUN_00523210  MOV AL,byte ptr [ECX + 0x49] ; RET      the footprint side n
vt+0x20  FUN_00523230  MOV AL,byte ptr [ECX + 0x4a] ; RET      the mask selector, 1/2/3
base ctor FUN_004f30a2 004f30e0 MOV byte ptr [EAX+0x49],1
                       004f30e7 MOV byte ptr [ECX+0x4a],1      <- the only stores in the image
```

The constructor's `1` is only the **default**. ~~So on a map loaded from a `.alm` every mover has a
1×1 footprint and the mask `0x41`, and the `0x44`/`0x82` arms of `FUN_0054b120` cannot be
reached~~ — **retracted by `TERR-MOVE-057`**: two of the four routines that hand these
bytes' *addresses* to a field-by-field (de)serializer (`FUN_004f5604`, `FUN_004f974d`,
`FUN_00500f46`, `FUN_00510518`) are the `Data.bin` param streamers, run inside the spawn path, so
`tokenSize`/`movementType` overwrite both whenever the table is not `−1`. Shipped: `movementType`
**2** on Ghost/Bee (mask `0x44`) and **3** on Bat_Sonic/Dragon (mask `0x82`, air), `tokenSize` 2 on
11 entries and 3 on 4 — all three mask arms and the `n×n` walk are live. The `units.reg` class array
`0x005eb674` is referenced 99 times by 20 functions, **all in the drawable/loader modules**, so no
registry key reaches the simulation at all. `TileSize` is *not* the footprint — it is what
**`CUnit`'s** `vt+0x20` returns, a different hierarchy at the same offset (`TERR-MOVE-054`,
`TERR-MOVE-055`; the one registry key that does select a C++ class is `Z`, and it selects the draw
layer — `REG-UNITS-061`).

### Movement speed (`TERR-MOVE-056`)

`FUN_0054d210(world, actor, u16 srcCell, u8 facing)`; the direction tables are written as
immediates by the world constructor.

```
dir = ((facing + 0x10) >> 5) & 0xff                                   0..7
dst = srcCell + (i16)world[0x58ec0 + 4*dir]
v0  = actor->[0x70]->[0x3c]->[0x44]  if nonzero  else (i16)actor->[0x8c]   (16 at construction)
if vt+0x20() != 1:  v = v0
else:
  d = clamp((i8)(height[src] - height[dst]), -32, +32)      height plane world+0x9451c
  v = SpeedMultiplier * v0                                  world+0x58db4, map.reg = 8
  v = v + ((v * d) >> 6)                                    SAR: uphill slows, downhill speeds up
  c = (u8)(cost[src] + cost[dst]) >> 1 ; if c == 0 then 8   byte-width add
  v = v / c                                                 signed idiv
v = clamp(v, 1, 63)
mover[0xa8] = v ; mover[0xae] = dir
dx = (i8)world[0x58eb0+dir] ; dy = (i8)world[0x58eb8+dir]
mover[0xb0] = dx*dy != 0 ? (i8)ftol(v*dx*0.707) : (i8)(v*dx)      double 0.707 @0x59cd98
mover[0xb1] = dx*dy != 0 ? (i8)ftol(v*dy*0.707) : (i8)(v*dy)
s = |mover[0xb0]| or |mover[0xb1]| if that is 0, or 1 ; mover[0xaa] = ceil(256 / s)
```

`cost(cell)` = `FUN_0054e5e0` is **not a pure read**: with `block[cell] & 0x20` and a nonzero byte
at `record+0xe` in the `world+0x540b8` table it shifts the stored cost right by 2 and writes it
back. Over every orthogonally adjacent cell pair of the 38 shipped maps at `v0 = 16` and
`SpeedMultiplier = 8`: `v ∈ [4..32]`, mode 16 — neither clamp bound is ever reached.

### Structures on the block plane (`TERR-STRUCT-068`…`072`, `TERR-PASS-073`)

The three tile-word arms, the type-3 cell and the border above are the whole of what the **ingest**
writes. A placed structure never goes through it. The sim class is the one the image names
`Building` (`CRuntimeClass 0x5c32d0`, `0x6c` bytes, vptr `0x59c738`; `Shop` derives from it), and it
attaches **in its own constructor**, after the ingest, to a per-cell record:

```
FUN_004e1924 (map load)  004e1c3f world ctor -> FUN_00547eb0 -> FUN_00548550 + FUN_00547d40
                         004e1c66 MOV [0x005f22c8], world        <- published, no null check later
                         004e1c76 FUN_004e2462   the .alm type-4 walker
   Building ctor FUN_005042b6 -> FUN_0050445d -> FUN_0054dbc0 -> FUN_0054d790 -> FUN_005456d0
   ~Building     FUN_005045e7 -> FUN_0054dc70 -> (payload+0x0c = 0, FUN_005456d0, free the record)
```

**The cell record**, `world+0x540b4`, keyed by the `u16` cell index, `0x34`-byte payload:

```
+0x00  u8    terrain-baseline COST byte        snapshotted at record creation (0054d8b7)
+0x01  u8    terrain-baseline STATIC block     snapshotted at record creation (0054d8c2)
+0x04  ptr   ground occupant  -> dynamic bit 6
+0x08  ptr   air occupant     -> dynamic bit 7
+0x0c  ptr   the Building
+0x10  ptr   set/cleared by FUN_005477b0 / FUN_005479b0; no plane arm reads it
+0x14..+0x28 six area-effect layer slots; each non-null one does cost <<= 2. +0x20 also blocks.
             SOURCED by TERR-STRUCT-074 and TERR-STRUCT-078: 005457e0 MOV ESI,0x6 ; 005457e5 CMP [ECX],0x0 ;
             005457ea SHL byte ptr [EAX],0x2 ; 005457ed ADD ECX,4 ; 005457f0 DEC ESI ; JNZ.
             WHICH SLOT IS WHICH (TERR-CELLREC-146): slot = +0x14 + 4*FUN_004fd8c3(spellId), and
             the two tables that function dispatches on are read out of the shipped image
             by tools/areamove:
               +0x14 spell  3 Wall of Fire     +0x24 spell 12 Light
               +0x18 spell  7 Freezing Cloud   +0x28 spell 17 Darkness
               +0x1c spell  8 Poison Cloud     +0x20 spell 19 WALL OF EARTH  <- the blocking one
             So a Wall of Earth is the only area effect that writes passability, and it does
             so through this arm rather than through anything in the area module. The
             registration is FUN_0054e730 (which then calls the recompute at 0054e829 or
             0054e97b) and the removal FUN_0054e9e0 (0054ea7c clears the slot, 0054eaff
             recomputes). MAGIC-WALLBLOCK-045, MAGIC-AREACOST-046.
+0x2c  u8
```

**Who takes `+0x04` and who takes `+0x08`** (`TERR-CELLREC-146`). Not *ground* and *air*
as such: the selector is the actor's movement-domain byte, read through `vt+0x20` in both directions
— `FUN_00544ec0` at `00544efc` (`JBE` rejects 0 and below, `<= 2` takes `+0x04`, `== 3` takes
`+0x08`) and `FUN_00545230` at `0054527f` (clears `+0x08` for 3, `+0x04` for 1 and 2). Domain 2 —
`Ghost` and `Bee`, `MOVE-DOM-028` — therefore shares the slot with ordinary ground movers. Both
slots hold at most one actor and a taken slot fails the entry (`00545001`, `005450d6`). `+0x10` is
the **sack** slot: `FUN_005477b0`'s one caller `FUN_0050f715` is the sack registration of
`ITEM-SACK-010`, and no plane arm reads the slot. Accessors by returned displacement: `+0x04`
`FUN_005463d0` / `FUN_00546520` / `FUN_0054a180`; `+0x08` `FUN_00546590`; `+0x0c` `FUN_0054df10`;
`+0x10` `FUN_00547be0` / `FUN_00547c60` / `FUN_00547cd0`.

**An actor occupies every cell of its footprint** (`TERR-FOOTPRINT-147`), the same shape
the building attach uses. `FUN_00544d00(map, actor)` reads the footprint side once
(`00544d12 CALL dword ptr [EDX + 0x1c]`), runs two nested loops both bounded by it (`00544e2f`,
`00544e49`) and calls `FUN_00544ec0` once per covered cell at `00544e6e`; a refusal from any covered
cell stops further iteration (`00544e75` to `00544e99 XOR EAX,EAX`), without local
rollback of earlier cells or mover+72/+82..85. Its five callers include the
sub-cell step `FUN_00548c60` and the arrival `FUN_005495f0`, so the `n x n` record entries are
rewritten on every cell transit, and `FUN_005456d0` derives dynamic bits 6 and 7 per cell from the
slot rather than from a separate footprint walk.

This is prefix-preserving failure, not atomic entry. A conflict at each ordinal
of a2x2 footprint leaves the earlier successful actor slots in place under the
selected normal-return paths. — SAV-CELLFAIL-583

Existing domain1/2 actor entry checks the trigger before testing occupied+04;
its caster call can precede an eventual refusal. Domain3 checks+08 without that
trigger arm. A missing record takes creation instead:52 zeroed bytes, current
cost/static baselines in+00/+01, then refetch and actor store. The creation path
does not revisit the existing-record trigger branch. Existing record reuse keeps
its baseline, tail and residue. — SAV-CELLENTRY-582

Sack registration is a separate programme. It reads Position+02 through
Sack+10, refuses Dynamic bit0, and rejects every nonzero existing+10 slot.
Writing an empty existing+10 returns without recompute. Missing-record creation
zeroes52 bytes, captures current Cost/Static before setting Static bit5, stores
the Sack, then recomputes. — SAV-SACKENTRY-590

Sack removal clears a present record's+10 without testing zero or identity,
then recomputes. Missing records return0. Its deletion predicate tests the four
occupant slots, byte+02 and byte+2c, not the six layer pointers or other residue.
— SAV-SACKREMOVE-591

The deletion arm restores Cost and Static from payload+00/+01 and preserves
current Static bit4. Dynamic is not copied from the restored baseline and its
record-present bit5 is not cleared: it retains the preceding recompute result,
plus the conditional bit4 OR. The recompute itself never reads Sack+10; its
other payload inputs remain active. Allocation-release effects are outside this
direct write-set. — SAV-SACKPLANES-592

Sack lookup separately requires Static bit5 before hash lookup. The creation
caller uses that lookup; registration's caller dispatches append or deletion,
and two selected removal callers continue without testing removal's result.
Those dispatches do not establish complete caller side effects.
— SAV-SACKCALLER-593

Detach tests the selected actor slot only for nonzero, not equality to the
actor argument. It clears and recomputes before testing whether four occupant
slots, layer-count+02 and operation+2c are all zero. Only that predicate enters
record deletion; other tail/residue bytes do not retain it. Successful detach
copies current Position cell/fractions into mover+86..89; missing-node/zero-slot
refusal returns before those stores. — SAV-CELLLEAVE-584

The exact cache mapping is entry+82/+83 = cached low bytes from cell-X/Y
accessors `00544a10/00544a20`, then+84/+85 = low bytes from full-X/Y accessors
`005449e0/005449f0` (Position+04/+05 under `SAV-TOKENPOS-074`). The first pair
is sampled before `0054a620`, the second after it; word+72 receives its AX.
The accessors and `0x1357` callback return are probe cuts, not proof of a runtime
+72 value or an atomic snapshot. After recompute/optional removal, detach
directly reads Position+00/+01/+04/+05 into+86/+87/+88/+89 in that order,
reloading actor+10 for each byte. — SAV-CELLFAIL-583, SAV-CELLLEAVE-584

**Bit 4** (`TERR-PASS-148`). Written by `FUN_0054e070`, four instructions that OR `0x10`
into both planes for one cell, called from the area module's per-cell add and blast on a fourth
argument meaning *the inner effect does damage*. `wall_of_earth` never reaches it. No mover mask
contains bit 4, and its one consuming reader is `FUN_0054e140`, which run-length encodes it over the
map interior into a `CArchive`. It is transmitted state, not passability.

The recompute reads the record through a **52-byte snapshot**: `005456ed` calls the map lookup and
`00545704 REP MOVSD` copies 13 dwords from `record+0x0c` into the scratch at `world+0x5402c`, and
every `payload+N` below is really `scratch+N`. A lookup miss returns at `005456f4` having written
no plane byte at all.

**The recompute** `FUN_005456d0(world, cellIndex)` — the only routine that turns a record into plane
bytes, and a no-op on a cell that has no record:

```
static  := payload+0x01 ; cost := payload+0x00        the terrain baseline, restored first
static  |= 0x20 ; dynamic := static                   bit 5 = this cell has a record
payload+0x04 ? dynamic |= 0x40                        ground occupant
payload+0x08 ? dynamic |= 0x80                        air occupant
if payload+0x0c:                                      a Building stands here
    pos = *(u8**)(obj + 0x10)                             ; a POINTER, see "the position object"
    bit = (cellRow - pos[1]) * obj[0x60] + (cellCol - pos[0])   ; SHL masks the count to 5 bits,
                                                          ; which is the only reason the engine's
                                                          ; own 32-bit intermediate is harmless
    if obj[0x64] & (1 << bit):  static |= 5 ; dynamic |= 5            ; BLOCKS
    else:                       static &= 0xfa ; dynamic &= 0xfa      ; OPENS
                                cost := costTable[5] (CostCracked, shipped 6)

  POLARITY, settled by branch displacement (`TERR-STRUCT-078`). It is decided by
  00545793 = 74 20: a JZ whose rel8 fixes the target at 00545795 + 0x20 = 005457b5, the
  block that begins b1 fa = MOV CL,0xfa. ZF is set when TEST finds the bit CLEAR, so
  CLEAR takes the AND-0xfa arm and SET falls through to OR 5. The fall-through measures
  exactly 0x20 bytes and the taken block exactly 0x28 (005457b3 = eb 28), so neither can
  be misaligned by a decode. The shipped Data.bin column title "Passability" is therefore
  the INVERSE of what the code does with it: a set bit is impassable.

  both arms load the operand once and then read-modify-write each plane separately:
    BLOCKS  0054579b MOV DL,0x5   ; 0054579d OR  BL,DL ; 0054579f MOV [EAX+0x10000],BL
                                  ; 005457ab OR  CL,DL ; 005457ad MOV [EAX+0x20000],CL
    OPENS   005457bb MOV CL,0xfa  ; 005457bd AND DL,CL ; 005457bf MOV [EAX+0x10000],DL
                                  ; 005457cb AND DL,CL ; 005457cd MOV [EAX+0x20000],DL
  There is no `AND r/m8,0xfa` in the routine: 0xfa reaches the byte only through CL.
for p in payload+0x14 .. +0x28: if p: cost <<= 2
if payload+0x20:  static |= 5 ; dynamic |= 5
if oldStatic & 0x10: static |= 0x10 ; dynamic |= 0x10  bit 4 is carried across every recompute
```

So a block byte is **derived state**, recomputed per cell from (terrain baseline, occupants,
building). A consumer must keep the baseline, not only the current byte, or it cannot demolish.

**The footprint** is a rectangle plus two 32-bit masks, all four from the class's `Data.bin`
Buildings entry (`FUN_0050445d`, params 0/1/4/5, the file's own column titles):

```
obj+0x60  u8    sizeX  the extent along the plane key's LOW byte  (param 0)
obj+0x61  u8    sizeY  the extent along its HIGH byte             (param 1)
obj+0x64  u32   "Passability"     the BLOCKING set    tested by FUN_005456d0 @00545790
obj+0x68  u32   "BuildingPresent" the ATTACH set      tested by FUN_0054dbc0 @0054dc06
```

`FUN_0054dbc0` walks `row = 0..sizeY-1`, `col = 0..sizeX-1`, bit index running continuously, and
attaches `cell = low16(((objRow+row) << 8) + objCol+col)` for each set bit of `BuildingPresent`.
This is addition, not bitwise OR: synthetic out-of-byte coordinates can carry into the other
coordinate. Neither registration body clips against the authored map dimensions
(`UNIT-STRUCTCELL-070`).
`sizeX` bounds the **inner** loop and is added to the position object's byte 0 — the same axis the
tile grid strides by 1 and the `.alm` type-4 record's `+0x00` carries (`ALM-OBJ-061`); the shipped
masks corroborate it, since `Horisontal Bridge` (6×4) and `Vertical Bridge` (4×6) are deck patterns
only under this assignment.
`FUN_0054d790` refuses a cell that already carries a building and the walk then aborts (2 refusals
over the 38 measured EN maps). A footprint over 32 cells aliases — the 11×4 `Castle` folds bits 32..43
onto 0..11. The `.alm` extension arm overrides `(w,h)` from the record and sets
`Passability = 0`, `BuildingPresent = 0xffffffff`; `kind == 0x21` is class id 33,
`Vertical Wooden Bridge`, so the arm exists to let a map author size a bridge, and it opens every
cell of the rectangle. **The arm is selected by `(w & 0xff) + (h & 0xff) > 0`** at
`00504509`…`0050450f`, not by the kind (`TERR-STRUCT-090`) — the caller-supplied bytes are file
`+0x14` → `obj+0x60` and `+0x18` → `obj+0x61` (`ALM-OBJ-062`), and the other two callers of
`FUN_0050445d` push literal zeros, so this arm has exactly one reachable caller. The
`Passability = 0` store happens on **both** sides of the arm's own `w·h > 32` test
(`00504559`, `00504568`, same immediate): the branch is dead (`TERR-STRUCT-077`).

Registration is not transactional. A collision aborts immediately and leaves earlier accepted
references in place; the constructor ignores the return value (`UNIT-STRUCTCELL-070`). Thus the
rectangle, mask-selected cells and successfully attached cells are different sets. The measured
EN/RU populations contain 17,062/9,741 selected cells but 17,057/9,741 accepted references under
ordered replay; neither EN collision exercises a partial prefix (`UNIT-AREAPOP-075`).

`~Building` independently calls `FUN_0054dc70`. This walks the current dimensions and mask again,
returns at a missing cell record, and clears an existing record's `+0xc` without checking pointer
identity. It does not replay a saved successful-attachment list. Therefore cleanup of a partially
registered object has a conditional ownership hazard; actual destructor reach and HP-to-destruction
ordering remain Unknown (`UNIT-STRUCTDETACH-074`). Ring and blast consumers read each current
cell reference anew, so a surviving alias can be visited repeatedly (`UNIT-AREAVISIT-071`).

The cell accessor does not read Building HP. Recompute reads Position, width and
the blocking mask, but no HP or class gate. Conditional on an unchanged reference,
mask and other cell inputs, HP 1, 0 and -1 therefore produce the same lookup and
blocking/opening contribution (`UNIT-STRUCTNEXT-079`). This conditional consumer
contract is not a lethal-hit lifetime rule: the bounded caller search leaves the
actual reference/mask survival boundary Unknown (`UNIT-STRUCTBOUND-081`).

**Two call sites, three callers** (`ALM-CLS-063`). `FUN_0050445d` is called only from
`FUN_005042b6` (`Building`'s ctor) and `FUN_0050433f`, but the `Shop` ctor `FUN_00505cbd` opens by
calling `FUN_005042b6` at `00505cea` — so a `kind ∈ {0x22,0x23}` placement attaches a footprint too,
always from the table. On the shipped maps that is 450 cells over 50 placements: a 3×3 square with
the bottom-left cell open as the doorway, 400 blocking cells and 373 cells that a ground mover could
otherwise walk through.

**The attach does not always recompute** (`TERR-STRUCT-076`). `FUN_0054d790` has two exits and only
the one that had to **create** the cell record ends `0054d994 CALL FUN_005456d0`. If a record was
already there with a null `payload+0x0c`, the building pointer is stored (`0054d809`) and the
routine returns 1 at `0054d85f` with no plane byte written — the cell keeps its old bytes until
something else recomputes it. At map load no record exists before the type-4 walk, so 0 of the
17 057 shipped attachments take that exit; a consumer that places a building at runtime must decide
what to do about it. The recompute's callers are eleven routines, not two:
`FUN_00544ec0` (×4), `FUN_00545230`, `FUN_005477b0`, `FUN_005479b0`, `FUN_0054d790`, `FUN_0054dc70`,
`FUN_0054e730` (×2), `FUN_0054e9e0`, `FUN_0054f680`, `FUN_00545de0`, `FUN_00545f50`.

**The position object** (`TERR-STRUCT-075`). `obj+0x10` is a pointer, allocated and stored by the
base actor constructor (`004f2523 PUSH 0xc`, `004f2525 CALL 0x00572824`, `004f2562 MOV [ECX+0x10],
EDX`) and written whole by `FUN_00544550`:

```
+0x00  u8   col           +0x01  u8   row
+0x02  u16  (row<<8)|col  the packed cell index the record map is keyed by
+0x04  u8   0x80          +0x05  u8   0x80    the half-cell sub-position
+0x08  ptr  the world
```

The chain from the file is complete **and unpermuted** (`ALM-OBJ-061`): the `.alm` type-4 record's
`+0x00` is read into the local at `[EBP-0x68]` (`00512903`) and `+0x04` into `[EBP-0x74]`
(`00512914`); those are `SAR ...,0x8`-ed at `005129ae`/`005129a7` and pushed **last** and
second-to-last, so they become the in-memory record's `+0x00` and `+0x02`; `FUN_004e2462` pushes
`byte[rec+0x00]` last into `FUN_00544550` (`004e25b8`, `004e25b1`), which makes it `[ESP+0x4]` and
therefore **byte 0**. So file `+0x00` is `col`. **The anchor is the object's own cell and the
footprint runs right and down** — nothing on that chain subtracts `(w-1)/2`, which is what excludes
the centred rival.

**A building can make terrain passable.** The `AND DL,CL` arm above — mask `0xfa`, applied once per
plane — erases the same two bits the terrain ingest set. Over the 38 maps / 3141 type-4 placements: **17 057 cells attached, 11 755 newly
blocked, 1113 opened** — 949 of the opened had been water, 145 a type-3 object, 19 Mountain, and 0
the bit-13 flag or the border. The masks read as pictures (`#` blocks, `o` opens):

```
38 Horisontal Bridge 6x4      52 Vertical Bridge 4x6      11 Church 4x3      27 Magic Symbol 3x3
   ######                        #oo#                        o##o               ooo
   oooooo                        #oo#                        ####               ooo
   oooooo                        #oo#                        o###               ooo
   ######                        #oo# (x6)
```

**And the opening matters: it is what connects the map** (`TERR-STRUCT-074`). Run 8-connected
reachability with the ground mask `0x41` — the neighbourhood `FUN_0054bd20` expands, full 3×3, **no
corner rule** — over both stages of the plane:

```
                                    stage A (ingest)   stage B (+ structures)
free cells to a ground mover               449 806            439 164
blocked                                    430 898            441 540
bridge-deck cells free                   391 of 1 299      1 299 of 1 299
bridge placements joining two components         0                    92   (of 104)
largest component, Islands.alm              22 189             35 135
                   Cross.ALM                15 549             39 575
                   scn:121.alm                 516              1 814
```

A consumer that implements the ingest and stops has a map whose halves do not meet. **Quote each
blocked-cell figure with the stage it was taken at**: the two numbers are two different maps.

**Save/load.** `Building::Serialize` (`FUN_005047c6`) round-trips `+0x40/+0x42/+0x44/+0x46/+0x48/
+0x60/+0x61/+0x64/+0x68` and does **not** re-attach; MFC's `CreateObject` ctor builds by the name
`"null"`, which misses the table and skips the attach. It does not need to: the world's Serialize
(`FUN_00544a60`) stores both plane sweeps (`TERR-PASS-053`) **and** the cell-record map itself
(`00544baa` storing, `00544c4c` loading) — the `u16 key + 0x34` records `SAV-CELLREC-017` measured,
whose key set equals the static-plane bit-5 set exactly in all four saves.

## Invariants (hold across all 38 maps, 880 704 cells, no exclusions)

- Strip group `g = (w&0x1fff)>>6 ∈ {0..12}` ⇒ file group `G ∈ {1,2,3,4}`.
- **Every referenced tile file exists**: 0 cells index a null/absent slot; `tile4` is
  referenced only at `g=12` → variants 00–03 (exactly the 4 files shipped).
- **Sub-cell in range for its BMP**: land (tile1/2/4, 14 cells) sub ≤ 13; **water (tile3,
  8 cells) sub ≤ 7** — 0 out-of-range on either. The shorter water strip being respected
  exactly is the discriminating check that pins the water grouping.

## Reading algorithm

```
# one-time: load terrain.3d/*.bmp into tiles[(G-1)*16+V], G=1..4, V=0..15; + dirt.bmp
def draw_cell(w, col, row):
    g   = (w & 0x1fff) >> 6
    b   = (w >> 4) & 3
    sub = w & 0xf
    if 8 <= g <= 11:                      # animated water
        g = 8 + (0 if anim_off else phase(g,col,row,anim_ctr))
    img = tiles[g*4 + b]                  # == tileG-VV.bmp, G=(g>>2)+1, V=(g&3)*4+b
    cell = img.subcell(sub)              # 32x32 block at pixel offset 8 + sub*0x400
    if (w & 0x2000) and not water:        # impassable, non-water
        cell = composite(cell, dirt.subcell((col + row*5) & 3))
    blit(cell, x=col*32, y=..., corner_heights=Altitudes[...])
```

## Cell geometry — where the pixels land (`TERR-GEOM-031…TERR-GEOM-036`)

The terrain raster is **not** flat. Each grid **vertex** is projected to a destination row and a
cell is drawn as the quad between its four projected corners.

```
vertex screen Y   V(c,r) = r*32 - (int8)Altitudes[worldCol + W*worldRow]
                            ^ screen row, not world row      ^ SIGNED (MOVSX), raw type2 byte
corners           yTL=V(c,r)  yTR=V(c+1,r)  yBL=V(c,r+1)  yBR=V(c+1,r+1)
destination X     col*32 + i, i = 0..31 — always 32 columns, NO horizontal term
```

**A larger altitude subtracts from the destination row index** (`FUN_00407117` @`004071c7`), i.e.
it moves the pixel toward the clip rectangle's `top` field — up the screen in the Win32 device
coordinates the clip rect is expressed in. The mesh lives at `CMapView+0xb4` and is rebuilt every
frame for a several-cell over-scan margin (`TERR-EDGE-026`); the *derived* `+0xc0` grid (mean of
four heights) belongs to a different consumer and must not be used here.

**Which blitter.** `flat` iff `yTL==yTR && yBL==yBR && yTL+32==yBL`, which is exactly *all four
corner altitudes equal* — no threshold, no epsilon. Corpus: **10.28 % flat, 89.72 % sloped**
(870 198 interior cells, 38 maps).

```
FLAT   (FUN_00450770): the axis-aligned rectangle [x,x+32) x [yTL,yTL+32);  src(i,j) -> (x+i, yTL+j)

SLOPED (FUN_004508f0): per destination column i = 0..31
    top(i)    = yTL walked toward yTR by the step table below
    bottom(i) = yBL walked toward yBR by the same table
    H         = bottom(i) - top(i)
    if H <= 0:  the column is SKIPPED — no pixel, no clamp, no fallback
    else rows top(i) .. bottom(i)-1          (top INCLUSIVE, bottom EXCLUSIVE)
         srcStep  = (32<<16) / H                       integer divide, 16.16
         srcRow(j)= (j * srcStep) >> 16                floor; j = 0..H-1, v starts at 0
         pixel    = LUT16[ level(j) + src[i + 32*srcRow(j)]*2 ]
         level(j) = (levelTop + 0x100 + j*((levelBot-levelTop)/H)) & 0xfffffe00
              levelTop = (c0<<9) + i*((c1-c0)<<4)   levelBot = (c2<<9) + i*((c3-c2)<<4)
```

So a column is **stretched or squeezed**, not shifted: 32 source rows are resampled into `H`
destination rows (corpus range `H` = 1…159). Flat is the degenerate case (`H = 32`,
`srcStep = 0x10000`). The shading ramp runs corner-to-corner over the *span*, not over 32 rows.

**The edge walk.** Both edges are rasterised from a startup-built step table at `0x005e4420`
(`.bss`; built by `FUN_0044ba10`, 128 rows × 128 bytes, row = `|Δy|`):

```
for d = 1..127:  s = (32<<16) / (d+1)                 # unsigned
                 acc = 0x8000
                 for k = 0..d:  acc += s;  T[d][k] = acc >> 16     # T[d][d] == 32 = sentinel

offset(i) = #{ k < d : T[d][k]    <= i }   # forward:  top edge if yTL<yTR, bottom if yBL>=yBR
offset(i) = #{ k < d : 32-T[d][k] <= i }   # mirrored: the other two cases
                                            # each edge latches once its index reaches d
closed form, exhaustively verified for d in [1,127]:
offset(i) = min(d, floor((2i+1)(d+1)/64))   # exact except d = 63 and d = 127
```

The convention is **inclusive at both ends** (`d+1` accumulator steps across 32 columns, sampled
at pixel centres) — *not* `d` steps. For `d >= 63` the drawn edge does not reach the far vertex by
column 31.

**Tiling, seams and degeneracies** (`TERR-GEOM-036`):

- A cell's bottom edge and the cell below's top edge are the same line but are walked with
  opposite quantisations. They agree for every slope **except `|Δy| = 63` and `= 127`**, where the
  upper cell paints one extra row that the lower cell repaints — an **overdraw, never a gap**
  (corpus: 9280 boundary columns, 9280 overdraw / 0 holes). Rows are drawn in ascending order, so
  the lower cell wins.
- Columns whose span collapses (`H <= 0`, i.e. the lower vertex is ≥ 32 above the upper one) are
  dropped: 3576 of 24 983 232 drawn columns (0.0143 %) in 353 cells. **They are not a hole in the
  frame.** `H <= 0` means the cell's own top and bottom edges have crossed, and the cells above and
  below share exactly those two edges, so the crossed range is painted by a neighbour. Whole-frame
  renders of Kids/Islands/Tomb (375 collapsed columns) on a surface pre-painted magenta leave **0**
  unwritten interior destination pixels. What is lost is the source pixels that column carried.
- The table has only 128 rows; `|Δh|` between adjacent signed heights could reach 255. The corpus
  maximum is exactly 127 — inside, but with nothing to spare. Behaviour beyond that is undefined.
- **Clipping is asymmetric:** a cell whose 32-px X range is not entirely inside the clip rect is
  **dropped whole** by both blitters, while vertical clipping is per-row.
- The engine's own hit-test (`FUN_0041a9f9`) bounds the cell by an *exact* lerp
  `y0 + ((y1-y0)*(x&31))/32` instead of the table walk, and the two disagree by up to **3 rows**.
  A renderer must reproduce the **table walk** — that is what is drawn.

## Water animation (`TERR-ANIM-006…TERR-ANIM-010`)

The water group is animated on a fixed-timestep logic clock (all from `rom.exe`):

- **Sequence.** A water cell keeps `b=(w>>4)&3` and sub-cell `w&0xf`; only the group rotates,
  so the drawn image walks `tile3` variant `V = phase*4 + b`, i.e. `tile3-{b, b+4, b+8, b+12}`,
  as `phase` steps `0→1→2→3→0`. The 16 `tile3` files = 4 phases × 4 blend columns. The
  `(worldCol+1)*worldRow` term (scroll-adjusted) offsets neighbours into a diagonal ripple.
- **Counter.** `animCtr = *(mapObj+0xa70)`; its only writes are init-0 (`FUN_004023c6`) and
  `+1` per logic tick (`FUN_0040dcdd`, the handler of window message `0x401`). `phase` uses
  `animCtr>>2` → one variant every **4 ticks**, full cycle every **16 ticks**.
- **Cadence.** The paced loop `FUN_004753c0` fires one `0x401` per `dtMs = 1000/tps` ms;
  `SetGameSpeed` (`FUN_00477370`) maps speed index `0..8` → `tps ∈ {8,10,12,14,16,20,24,28,32}`.
  Map-load defaults to **index 4 → 16 tps → 62 ms/tick → ~250 ms/frame → ~1 s/cycle**; the +/-
  keys and config key 2 change it. Per-frame ms = `4·(1000/tps)`, per-cycle = `16·(1000/tps)`.
- **Enable.** `DAT_005bcef0` (static default `1`); cleared by `-noanimation` / `-detail0` (then
  `phase=0`, water static). Verified: 38 maps / 87 580 water cells, all 4 phases resolve to a
  present `tile3` file (0 violations). See `claims/terrain.md` `TERR-ANIM-006…010`.

## Terrain lighting (`TERR-LIGHT-011…TERR-LIGHT-013`, `TERR-LIGHT-015`,
`TERR-LIGHT-018…TERR-LIGHT-023`, `TERR-LIGHT-028`, `TERR-LIGHT-029`,
`TERR-LIGHT-109`, `TERR-LIGHT-119…TERR-LIGHT-128`)

Terrain is **relief-shaded** (Lambertian slope shading), computed at runtime — the `.alm` stores
no lightmap:

- **Blit.** Each pixel is `LUT16[(level<<9 & 0xfffffe00) + srcIndex*2]`; `level` is bilinearly
  interpolated from four per-vertex brightness bytes (`FUN_00450770`/`FUN_004508f0`). The LUT is a
  palette brightness ramp rebuilt on light change.
- **Per-vertex brightness** (`FUN_00484f40`, grid at `P+0x18`, `P` = CMapView+0x80) — exact level:
  ```
  tanT   = tan|θ|                            H[] read MOVSX (signed bytes)
  stepH  = 32.0 / cos|θ|
  Δh₁    = H[x ∓ tanT, y+1] − H[x, y]        forward,  ONE row   (∓ = −tanT if θ≥0, +tanT if θ<0)
  Δh₂    = H[x, y] − H[x ∓ tanT, y−1]        backward, ONE row
  sᵢ     = atan2(Δhᵢ, stepH)
  axisᵢ  = clamp( L − R·sin(π/6 − sᵢ), 0, 95 )
  byte   = ftol( 0.5 · (axis₁ + axis₂) )     ← the shading level the blit LUT indexes
  ```
  `R = DAT_005eb498` (day 32), `L = (R>>1) + DAT_005eb494 + 32` (day 62). Flat day vertex → 46.

  **Both differences span exactly one cell, along the same (row) axis** (`TERR-LIGHT-028`). There is *no*
  two-cell central difference and *no* perpendicular `x` gradient anywhere in the routine — the two
  terms are the forward and backward differences of one axis, and the only `x±1` reads are the
  lateral operands of the `tanT` shear *within* rows `y+1`/`y−1`, i.e. the sample walks along the
  sun azimuth while `stepH` lengthens the baseline to match. The `0.5` is applied to the sum of the
  two **already-clamped levels**, never to a height difference — so the routine is not equal to any
  single-difference form once either arm clamps. **A port using a central difference over-contrasts**
  (it roughly doubles every gradient). Which lateral column axis 2 uses is selected by `y == 1`, not
  by the θ sign the code computes and discards (`DEC EDI` clobbers the compare's flags); the effect
  is confined to row 1 and is not measurable in the corpus range.

  Written region = the **strict interior `{1..W-2}×{1..H-2}`** only (`FUN_00484f40(1,1,0,0)`; the `±1`
  neighbour's sign follows the sun azimuth, so starting at 1 / ending at W-2 keeps both in bounds).
  The **outer ring of vertices is NOT computed** and has no one-sided fallback (`TERR-EDGE-025`, correcting an
  earlier "first row/col one-sided" reading); it is left at the non-zeroing allocator's value. The grid
  is allocated by `FUN_00484c70` (`malloc(W·H)`, unpadded) but never read from the file (computed).
  Constants read from the x87 disassembly: `32.0`, `π/6`, `95.0`, `0.5`.
- **The shading table** (`FUN_00427df0`, called via `FUN_00428ad0` from the relight driver
  `FUN_0046f1e0` as `(nLevels=0x60, mode=3, useTint=1)`). It is allocated `malloc(nLevels<<9)` —
  **96 rows × 256 entries × u16 = 49 152 B, row stride 512** — and filled per channel with

  ```
  out_chan = clamp( ((palette_chan + skyTint_chan) * (96 - level)) / 32, 0, 255 )   # trunc toward 0
  entry    = (R>>(8-rBits))<<rShift | (G>>(8-gBits))<<gShift | (B>>(8-bBits))<<bShift
  ```

  Integer, linear per channel; no gamma, no cross-channel term. The bits/shifts are the DirectDraw
  surface's (derived at runtime by `FUN_0044c920`; shipped = **RGB565**). The palette is the tile
  object's own BMP palette (`obj+0x20`); the table lands at `obj+0x1c` and the renderers pass
  `*(DAT_005ef8d8 + 8)`.

  **Row 64 is unattenuated** (`96−L = 32`); `L=0` is ×3.0 and `L=95` is ×1/32, so the level byte is
  an *attenuation* index. 96 rows is exactly the `[0,95]` clamp of the per-vertex byte above.
  **Flat daytime ground is level 46 → ×1.5625** — the shipped art is authored dark for that
  (pixel-weighted mean luminance 78.5 at ×1.0 vs 119.2 at ×1.5625), so rendering "unshaded" means
  level 46, not ×1.0. Only the **sky tint** enters the table; the intensity bytes shape the level.
  One table serves all terrain (the driver rebuilds the first tile slot and `break`s) — sound
  because all 53 `terrain.3d` palettes are byte-identical. For what this mapping actually puts on
  screen, including where it **clips to white**, see "Terrain output range" below; for the sprite
  table built by the same routine in mode 2, "Sprite lighting".
- **Sun** (`FUN_00468b10`): day/night cycle — hour `= (t/60)%24` sets the sky RGB tint + intensity
  bytes and sweeps the sun angle **θ ∈ [−0.78539815, +0.78539815]** — a truncated-π *quarter* each
  way, so the whole sweep is π/2. *(This paragraph previously read "[−π/2, +π/2]"; both bounds
  were taken from a disassembler's symbol names rather than the bytes — `TERR-LIGHT-109`,
  `retracted.md`.)* The clock `t` is **`fullTicks + 360`**, `fullTicks = campaign+0x3e0 >> 4` where
  `campaign+0x3e0` is the server sub-tick counter shipped as command `0x64`: **one in-game minute is
  one full tick**, a mission starts at 06:00, and an in-game hour is ≈ 59.5 s of real time at the
  shipped speed index (`SESS-TICK-026`, `SESS-TICK-027`). Five arms write θ:

  | band | hours | store | θ |
  |---|---|---|---|
  | cycle off | — | `00468b59` | literal `+0.78539815` |
  | day | 6…17 | `00468bff` | **computed** `−0.78539815 + m·0.0021816615`, `m = (t+360) % 720` |
  | dawn | 2…5 | `00468d6a` | literal `−0.78539815` |
  | dusk | 18…21 | `00468ee5` | literal `+0.78539815` |
  | night | 22,23,0,1 | `00468f54` | **computed** `+0.78539815 − m·0.0065449846`, `m = (t+120) % 240` |

  The **colour** sees six arms, not five: each twilight band splits in two at `00468c1e` /
  `00468d8e` and each half is a 120-minute ramp. With `m = t mod 120` and `q(k) = (k·m)/120`
  (one magic division, multiplier `0x88888889`, post-shift **6**, i.e. /120):

  | band | hours | R `0x5eb490` | G `0x5eb491` | B `0x5eb492` | ambient `0x5eb494` | amplitude `0x5eb498` | shroud `0x49c`/`0x4a0` |
  |---|---|---|---|---|---|---|---|
  | cycle off | — | 0 | 0 | 0 | 14 | 32 | 4 / 2 |
  | night | 22,23,0,1 | 0 | 12 | 48 | 32 | 8 | 8 / 4 |
  | dawn 1 | 2,3 | `q(24)` | 12 | `48 − q(48)` | `32 − q(4)` | `8 + q(12)` | 6 / 3 |
  | dawn 2 | 4,5 | `24 − q(24)` | `12 − q(12)` | 0 | `28 − q(14)` | `20 + q(12)` | 6 / 3 |
  | day | 6…17 | 0 | 0 | 0 | 14 | 32 | 4 / 2 |
  | dusk 1 | 18,19 | `q(24)` | 0 | `q(8)` | `14 + q(10)` | `32 − q(12)` | 6 / 3 |
  | dusk 2 | 20,21 | `24 − q(24)` | `q(12)` | `8 + q(40)` | `24 + q(8)` | `20 − q(12)` | 6 / 3 |

  Every value above is an **immediate in `.text`** or that quotient — the routine reads five
  addresses in total (the `ShowTimeFlow` flag and the four `.rdata` doubles the two
  computed-angle arms use) and **no table** (`TERR-LIGHT-119`, `TERR-LIGHT-120`). The three
  bytes are `R, G, B` in ascending address order (`TERR-LIGHT-124`), and the tint is taken
  **unsigned** and added to the palette channel *inside* the level multiply, so it can only
  brighten — all darkening is the level's (`TERR-LIGHT-127`). The six arms form one
  continuous programme: the largest step across any band join is **1**, and night's tuple is
  exactly dawn 1's opening tuple (`TERR-LIGHT-123`). Flat-ground level therefore runs
  **`[45,64]`** over a day against 46 by day, the map-wide sprite level 3 → 8, and terrain
  relief **flattens** as it darkens because `0x5eb498` falls 32 → 8 (`TERR-LIGHT-125`).

  **A consumer must not interpolate this per minute.** The relight cadence below samples each
  120-minute ramp **six** times, so only **24** distinct `(R,G,B,0x494,0x498)` tuples are
  reachable through the unforced path out of the 170 the arithmetic can produce; the forced
  path bypasses the modulo and can reach any of them (`TERR-LIGHT-128`).

  So 960 of 1440 full ticks in a day carry a moving sun (`TERR-LIGHT-110`). Relight fires from
  `FUN_00421a54` when `(clock & 0xf)==0 && ((clock>>4)+0x168) % 20 == 0` — once per **20 in-game
  minutes**, ≈ 19.8 s — or forced: sun → per-vertex grid → all tables → redraw (`TERR-LIGHT-114`).
  With the cycle off the model writes ambient `0x0e` / range `0x20`, i.e. the `L=62, R=32` daytime
  configuration. **The cycle is on by default**: `DAT_005eb528` is the `ShowTimeFlow` game option,
  set to 1 by a constructor the C++ static-initializer table runs before `main`, persisted to the
  registry and the savegame, and toggled at runtime by key `N` (`TERR-LIGHT-108`). Since both
  `stepH` and the lateral shear depend on θ, **no per-vertex level statistic is meaningful without
  stating its θ** (`TERR-LIGHT-029`).
- **θ propagates through one cache, and it is what leans every shadow.** `FUN_00484f40` is the only
  reader of `_DAT_005eb4a8` image-wide and its first act is to copy the double onto the terrain
  object at `+0x20` (`00484f57`/`00484f5a`); nothing reads the global again. `FUN_00484c10` reads
  that cache and returns **θ·2/3**, except inside a dead band where `−0.05 < θ < 0` returns `−0.05`
  and `0 < θ < +0.05` returns `+0.05`. All ten of its call sites are shadow passes — unit
  (`FUN_0045bf00`), structure (`FUN_0045a260`), object (inside `FUN_00407b1a`), `FUN_00459dc0` —
  and each takes `tan` of the result and multiplies by 65536.0 into the blitter's 16.16 per-row X
  slope. `tan(shear)` spans **−0.5774 … +0.5754** over a day and is fixed at **+0.57735** with the
  cycle off. The unit **body** `FUN_0045b3f0` calls the same routine and discards the result
  (`0045b44f FSTP ST0`) — the body never leans (`TERR-LIGHT-111`…`113`).
- **The map's stored light fields** land in the lighting slots — type-0 payload `+0x08` (f32 angle) →
  `P+0x20`, `+0x10`/`+0x14` → the intensity bytes `P+0x1c`/`P+0x1d` *(offsets on the `ALM-FRAME-031`
  record-payload base; this paragraph quoted the superseded `+0x10`/`+0x18`/`+0x1C` until the
  2026-07-25 sweep, although `TERR-LIGHT-023` had been re-anchored)* — but `FUN_00484f40`
  **overwrites all three from the sun globals before reading them**, and the θ slot is a `double`
  the loader fills only 4 bytes of. So they are *not* the map's initial lighting (amends
  `ALM-META-009` / `TERR-LIGHT-015`). **Nothing reads them**, and the search behind that is now
  bounded rather than open (`TERR-LIGHT-149`): the fourth stored scalar, payload `+0x0c` →
  `P+0x2c`, has zero accesses of any kind in the whole searched population (`TERR-LIGHT-150`);
  the overwrite is not merely first inside `FUN_00484f40` but forced by the load path, which
  reaches the relight driver on every path from the store that publishes the landscape object,
  behind one gate; and the values copied in cannot come from a file, because all 53 stores to
  the sun globals lie inside `FUN_00468b10` and no instruction anywhere loads one of their
  addresses into a register (`TERR-LIGHT-151`). The map-object copies of the same file bytes
  are never read either. The one scalar of the seven-read sequence that *is* consumed is
  payload `+0x18` → `P+0x28`, the tile-group mask above (`TERR-LOAD-152`).
- Corpus: 38/38 maps carry a height grid. *(The figures "non-flat, σ 7…41" are **withdrawn** by the
  2026-07-25 sweep — `TERR-LIGHT-016`: that census was measured on the superseded `ALM-GRID-011` grid base and
  every map's reported maximum is a record-header byte the old base injected. Whether any shipped
  map is flat is currently **unknown**; re-running the census on the corrected base settles it.)*
  Applying the exact formula above to all 38 grids at the engine's default θ = 0.78539815 over the
  interior window yields levels **30…70** (10 root maps alone: 30…65); the ceiling rises to 81 at
  the ends of the day's θ sweep. **0 / 859 768 vertices fall outside the 96 rows** — and cannot:
  the per-axis clamp confines the level to `[30, 89]` for *any* height field, and the floor 30 is a
  saturation (`sin(π/6 − s)` peaks at `s = −π/3`), not a property of the corpus.
  *(The superseded figure `30…74` was measured on the old `ALM-GRID-011` grid base by an
  approximate probe rule — see `TERR-LIGHT-029`.)*
  See `claims/terrain.md` `TERR-LIGHT-011…030`.

## Far-edge cells (`TERR-EDGE-024…TERR-EDGE-026`)

All render grids are **unpadded `W×H`** (`FUN_00484c70`: `malloc(W·H)` per byte grid, `malloc(W·H·2)`
for type1; corpus type1 `==2·W·H`, type2 `==W·H`, 38/38). A cell `(col,row)` samples four **vertex**
corners, but `W×H` vertices only fully bound `(W-1)×(H-1)` cells — so a cell in the last column/row
has no `+1` vertex in-grid. The engine does **not** clamp, duplicate, wrap, or use a side table:

- **Corner sourcing = raw flat row-major addressing.** The four renderers
  (`FUN_004058c7/00405e83/00406349/00406806`) read `grid[idx], grid[idx+1], grid[idx+W], grid[idx+W+1]`
  (`idx = col + row·W`) and pass the values to the blitters. So a **last-column** cell's `+1` corner
  reads `grid[(row+1)·W + 0]` — the **next row's column 0** (a row-major wrap, in-bounds except at the
  last row); a **last-row** cell's `+W` corner reads `grid[H·W + col]` — **past the `W×H` allocation**.
  The same flat scheme drives heights (projection `FUN_00407117`) and tile words (object passes).
- **The outer ring is drawn, not skipped.** No renderer special-cases it; the mesh is built (plus a
  several-cell over-scan margin beyond the grid) and blits are issued, gated only by screen culling.
  The one bounds check in the whole draw is a single `if (worldRow < H)` guard on the smoothed-height
  pass — its presence proves the draw reaches/exceeds the grid edge.
- **But the outer band's shading is degenerate.** The per-vertex brightness outer ring is never
  computed (interior `{1..W-2}×{1..H-2}` only) and the allocator does not zero it, so the far corners
  of the outermost cells are uninitialised / wrapped / out-of-bounds. The engine tolerates this
  because gameplay is bounded by the **derived 8-cell sim border** (`FUN_00547d40`, `ALM-TERR-016`) and
  the camera keeps that band at the extreme edge.
- **Reimplementation guidance:** the far edge has **no defined shading** to reproduce — a faithful
  port should **clamp-to-edge** (safe, defined) and treat the outer band as the non-gameplay border it
  is. See `claims/terrain.md` `TERR-EDGE-024…026`.

## Sprite placement — where a unit or object stands (`TERR-SPR-038…TERR-SPR-041`)

Sprites are **lifted** by the terrain, never **warped** onto it. Enumerated over the whole
`rom.exe` image, the geometry step table at `0x005e4420` has exactly four readers — the sloped
terrain blitter and the three shroud blitters — and **no sprite blitter indexes it**
(`TERR-SPR-038`). The sprite draw call carries two scalars and no corner Y. (`TERR-SPR-038`'s push
list named its fifth argument a `level`; it is the shadow's per-row X slope — `TERR-SPR-066`, and see
"Sprite lighting" below for the argument that really is a brightness.)

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

They agree on **every** flat cell (necessarily — four equal corners are their own mean) and, less
obviously, on **9.26 %** of sloped ones: over 136 291 interior cells of Kids/Islands/Tomb, 12 513
of the 135 129 sloped cells have `mean == corner`, 6227 of them with four *pairwise-distinct*
corners. "Sloped" is a statement about the corners being unequal, not about their mean; do not use
one as a test for the other. Where they differ, `mean − corner` spans **−73 … +55** destination
rows. Measured by `tools/terrgfx -mode sprart` (`TERR-SPR-039`),
which also puts the two models against each other in pixels: **134 666 / 136 291** ground-contact
points fall inside the destination-row span the cell's own terrain raster fills at screen column
`col*32+16`, **0** fall where nothing is drawn, none above their own span and at most **2 px**
below it. That measurement discriminates the divisor, the corner set, the index base and the
presence of the lift; it does **not** discriminate the rounding rule or the trunc-toward-zero
fixup — both rest on the listing alone (`TERR-SPR-039`).

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
so a `.alm` census of it is the initial state only.

**Bits 15..14 are the FOG OF WAR** (`TERR-TILE-079`, `TERR-FOG-080`, `TERR-FOG-081`).
They gate the animated object arm above and `FUN_00405e83`, the partial-repaint terrain renderer,
which skips every cell that fails the test. **No shipped cell sets either bit** (0 of 880 704) —
that is a fact about the **file**, and a customised `.alm` must keep it so, because the runtime
grid is the same words:

```
state  11  currently in sight     10  explored, not in sight     00  never seen
set    FUN_00462f90 = CUnit/CAirUnit vt+0x48, tail-jumping to the body at 004598c0:
       OR byte ptr [tile + idx*2 + {1,3,W*2+1,W*2-1}], 0xc0    0045997a 0045997f 0045998b 00459990
       -- ONE immediate, both bits, so a single stamped word satisfies the four-corner OR.
       SECOND writer of bit 15: 00478870 OR word ptr [ESI],DI in the save-load path
       FUN_00477c00, restoring the saved run-length record (SAV-FOG-061, TERR-FOG-145).
       It has a register source, so no imm: sweep can witness it.
clear  bit 14 only, map-wide, on a period: 0040eb44 AND word ptr [ESI],0xbfff  (FUN_0040eaee)
       bit 15 is cleared by nothing. imm:3fff = 25 hits/16 owners/0 orphan, imm:7fff =
       66/41/0; of the 91, the 8 with a memory destination are all whole-dword MOVs of a
       0x7fff/0x7fffffff sentinel -- not one AND of a 16-bit word. The two hits inside
       map/drawable code are read, not assumed: FUN_00451ae0's four are the sprite RLE
       opcode strip (TEST AX,0x4000 / 0x8000 above them) and 0045e9f8 is the CDQ/AND/SAR
       signed-divide-by-0x4000 idiom. Blind spot: a clear through a register carries no
       immediate at the store. This is why FUN_00459f50 tests == 0x8000 before == 0xc000.

guards, all three, before anything is stamped
  0  00462f90 CMP byte ptr [ECX+0x15a],0x3 ; JNC        decay stage >= 3 -> nothing
  1  004598e0 TEST byte ptr [[[view+0x9b4]+0x38] + [[this+0x14]+0x4]*2], 0x8
       bit 3 of a per-player u16 on the LOCAL participant's record, indexed by the
       drawable's owning Player+0x04. FUN_004757b0 writes flags[localPlayer] = 0x0a at
       session setup and its refresh loop preserves bit 3 (0047607e AND EDX,0x8); the only
       other writer is the message arm at 004162ec. In single player: the local player only.
  2  004598ea CMP word ptr [this+0x102],0x0             sight; the ctor leaves it 0
vt+0x44 = FUN_00462f80 repeats 0..2 per tick and adds
  3  (this+0x50 & 0x1f) == 0x10 and (this+0x54 & 0x1f) == 0x10     an eighth-of-a-cell grid
  4  this+0x8/+0xc != this+0xc0/+0xc4                              it moved since last stamp

the stamped set — CMapView+0x17cc, 41x41 int32, memset to 0 by FUN_00403c8c on EVERY stamp
  seed   mask[20][20] (= view+0x24ec) = (sight >> (8-k)) + (1 << (k-1)),  k = view+0x3f38 = 7
                                        so the unit of the field is 1/128 cell
  walk   Chebyshev rings r = 1..19, four edges each; stop at the first fully blocked ring
  cell   FUN_00403b77:  mask[dx][dy] = mask[pred[dx][dy]] - (cost[dx][dy] + h(cell) - h(obs))
         visible iff > 0.  h is map+0x10, signed bytes (the .alm type-2 Altitudes grid);
         h(obs) is sampled ONCE, so the term is per-cell altitude vs the observer's, NOT a
         slope along the ray.
  tables built once in the CMapView ctor by FUN_00403718:
         pred  view+0xaa8   41x41 int8 pairs, one step toward the observer; three zones,
                            j < i>>1 -> (-1,0), j > i<<1 -> (0,-1), else (-1,-1), mirrored
         cost  view+0x3210  41x41 int16 = ftol(128 * sqrt(i^2+j^2) / max(i,j))
                            axis 128, diagonal 181 -> the revealed region is a DISC
  margin ring cells outside [7, W-7) x [7, H-7) are never evaluated and stay 0
  clip   19 cells; max shipped scanRange is 12, so it never binds on shipped data

sight   drawable+0x102 <- actor+0xa4 by the state sync (0047c2a3 / 0047c2aa)
        hero:     FUN_004f7dfc writes ftol(((mind+reaction)/25 + 4) * 256)   -> 1/256 cell
        non-hero: the streamer's slot 11 targets actor+0xa5, the HIGH BYTE   -> whole cells
                  Data.bin Units title 11 "scanRange", Humans title 9 "ScanRange"
```

Consumer consequences: `Index` is in range on 82/82 shipped classes but a decoder should still
bounds-check it; `DeadObject` is a subscript into the same class array; `FireObject` is **not** a
class reference (`REG-OBJ-047`). **And the animated arm is not dead**: it fires on exactly the
cells the local player can currently see, so a shipped map's fires and trees animate inside the
field of view and hold frame 0 outside it. A consumer that implements only the plain arm
reproduces every shipped map at load time and never afterwards.

**Shroud / fog** (`TERR-FOG-037`): after the terrain, the same cell quad is darkened in place by
`FUN_00451200` (Gouraud level ramp, `dst = LUT[level][dst]`), or by `FUN_00450cf0` (level `0x10`,
degenerate quads filled with colour 0) or `FUN_00451710` (level `8`, degenerate quads halved as
`(px>>1) & mask`). These reuse the terrain quad and the same step-table edge walk.

**What the renderer branches on, and the levels it turns the pair into**
(`TERR-FOG-082`, `TERR-FOG-083`, `TERR-FOG-084`, `TERR-FOG-085`). The shroud pass does not read the
tile word directly: `FUN_00404135` projects the pair onto a **per-vertex** dword grid every frame.

```
grid    CMapView+0xa0, one dword per lattice vertex, (cols+7)*(rows+11) entries
        allocated with +0xa4 by FUN_00402c88 at +0x70 = (cols+7)*(rows+11)*4 bytes
        memset to 0 at 004043d7, filled, then memcpy'd to +0xa4 at 004046f2
        FUN_00404135 is its only content writer; FUN_00402f85 frees both

classify per vertex, off [CMapView+0x80]+0xc  (= map+0x0c, the SAME plane the sim uses)
        004044d0 MOV DX,word ptr [ECX+EAX*2] ; 004044d4 AND EDX,0xc000
        == 0xc000  -> level 0    00404644     11  in sight        NO shroud drawn at all
        == 0x8000  -> level 8    00404674     10  explored        half brightness
        otherwise  -> level 0x10 00404695     00  never seen      black
        01 is not a state: the stamp writes one immediate 0xc0, and nothing clears bit 15,
        so a renderer branches on the PAIR and needs no fourth arm.

dispatch FUN_00407b1a reads the four corner vertices (0040c1ce..0040c305) and branches
        (0040c421..0040c679) on all-equal-0 / all-equal-0x10 / all-equal-8 / otherwise.
        A cell whose corners DISAGREE is a gradient -- one fog value per cell cannot draw it.

table   [0x005e8420], built once by FUN_0044ba10: 17 rows (L = 0..16) of `stride` u16,
        stride = 65536 or 8192 by [0x005eb570]   (0044ba7f SHL ECX,4 ; ADD ECX,EAX ; SHL ECX,1)
        LUT[L][px] = each channel index scaled (v * (16-L)) >> 4, saturated at 0xff, repacked
        row 0 = identity   row 8 = (px>>1) & [0x005eb578]   row 16 = 0
        checked against the engine's own two table-free fast paths over all 65536 pixel
        values, RGB565 (mask 0x7bef) and RGB555 (0x3def): 65536/65536 on all three rows.
        [0x005eb578] = sum over channels of (0x7f >> (8-bits)) << shift   (0044bd29..0044bd85)

clock   FUN_0040eaee clears bit 14 over the WHOLE map (W*H words) then re-stamps every
        drawable in +0x9b8. Its one caller FUN_0040dcdd gates it:
          0040dd54 CMP dword ptr [0x005eb588],0x0 -- nonzero SKIPS the clear, freezing the layer
          0040dd8d AND EAX,0x1f ; JNZ             -- CMapView+0xa70, the ANIM-CLOCK-001 counter
        => once every 32 PRESENTATION ticks ~ 2 s at the default speed index, ~4 s at the
        slowest. Neither server+0x00 nor server+0x04 appears on this path: the fog is not
        simulation state, is not hashed, and stops when rendering stops.

gates   1. the shroud pixels above
        2. FUN_004597f0: OR the four corner words, & 0xc000 != 0xc000 -> drawable+0x78 = 1,
           the sprite pass's guard (00459814..00459839). A unit on a cell with no visible
           corner is NOT DRAWN. One visible corner suffices. +0x7c latches +0x10c on change.
        3. nothing in the simulation. imm:c000 = 85 hits/21 owners and every game-module
           owner is in 0x00404135..0x0048f784; that sweep is blind to byte-wide tests on the
           high half, so the load-bearing half is that the AI's own visibility input is a
           DIFFERENT array (below) and passability reads the block planes (MOVE-DOM-027).

persist BIT 15 IS SAVED; bit 14 is not. Not by the map: 0 of 880 704 EN cells and 0 of
        677 696 RU cells author either bit outside the [u32 typeId][f32] head overlay at
        cells 0..3 (one bit-15 hit per map, always at cell 3 -- the TERR-LIGHT-016 artefact),
        and the .alm supplies the plane on both arms of FUN_004d00e9. By the save: the
        uncompressed tail's &YA1 state store carries a section Fog with FirstState (int32)
        and Data (int32[]), a run-length encoding of bit 15 over W*H cells in the plane's own
        order idx = col + row*W. Store FUN_00478c40 00479343..004793dc; load FUN_00477c00
        0047883a..00478884, whose per-cell write is 00478870 OR word ptr [ESI],DI -- a
        register source, which is why an imm: sweep could not see it (SAV-FOG-061,
        TERR-FOG-145). sum(runs) == the map's cell count on 14/14 saves that carry the key.
        => a load restores the explored map. Persist bit 15; re-derive bit 14, which the
        32-tick clear at 0040eb44 and the stamp rebuild anyway.
        The 8192 B arithmetic that TERR-FOG-087 used against this measured the COMPRESSED
        body; the record is in the tail and costs 4..1092 bytes as runs.

edge    the map border is black because it is NEVER SEEN -- level 16, the same mechanism at
        its maximum, and TERR-FOG-080's 7-cell stamp margin is why it is never lit. There is
        no separate edge treatment. (Owner: "край карты чёрный, незатемнённый".)

reveal  permission to stamp is bit 3 of [[mapView+0x9b4]+0x38][player], and that array has
        THREE writers: FUN_004757b0's setup store, dispatcher opcode 33 (0x21, arm 004160c0)
        writing ONE entry at 004162ec, and dispatcher opcode 45 (0x2d, arm 00415c85) which
        resizes the array and memcpy's it WHOLESALE out of the message body --
        00416067 memcpy([[view+0x9b4]+0x38], msg+0xe, [msg+0xa]*2). Opcode = 3 + (slot -
        0x004186e3)/4 from the dispatcher's own SUB/JMP. A wholesale copy carries no
        displacement, which is why a store-form sweep cannot see the third one.
```

**The AI's vision is a second implementation of this same algorithm, not this one**
(`TERR-FOG-088`). It lives in an object embedded at `world+0x58ee8` — so `AI-SIGHT-006`'s byte
map at `+0x2a008` **is** `world+0x82ef0` — with the same recurrence and its own `pred`
(`+0x22000`), `cost` (`+0x28000`) and accumulator (`+0x24000`, whose centre cell `[20·64+20]` is
the address that row calls the field `fog+0x25450`). Its `k` is **`[Scanning] ScanShift` from
`World\Data\map.reg`** (ships 7), while the view's `k` is the compiled constant at `00402b61`;
they agree only because the shipped value equals the code default, so **editing `ScanShift` moves
the AI's sight and leaves the player's fog untouched**. Its byte map is cleared by
`FUN_005474f0`, whose two callers stamp **one actor** (`FUN_0052dea0`) or **a whole group**
(`FUN_005365e0`). Different storage, clock, consumer and lifetime — do not merge them.

Both of the server's tables are built by **`FUN_00546790`**, called once from the object's init
`FUN_00547510` in each of the four world constructors, before the map load; there is no other
writer and no other reader than `FUN_00546c70` (`TERR-SIGHT-115`, `AI-LOS-087`…`AI-LOS-091` —
the layout and the region are specified in `formats/ai/format.md`, which is where the stamp is
consumed). The same init builds a **fourth** grid at `+0x20000` holding `di² + dj²`, which serves
a disc test in `FUN_00546fa0`/`FUN_005470a0` and has nothing to do with line of sight.
**The view and server use one algorithm** (`TERR-FOG-117`, `AI-SIGHT-093`). Nineteen clauses
of the two match instruction for instruction, including the
**four trailing literal stores** — `00403b4b`…`00403b6c` on the view against `0054699e`…`005469b3`
on the server — that repair the two cells `(+1,0)` and `(-1,0)` whose slope classification leaves
them pointing diagonally. Those four are what `TERR-FOG-080` did not transcribe. Re-executed with
them, both sides give **9 / 21 / 45 / 69 / 105 / 145** cells on flat ground for `scanRange` 1..6 at
`k = 7` and 1253 at 19; without them the client alone gives **59 / 127 / 223 / 479 / 847 / 1183**,
which is every figure the old reading published. **145 is correct and 127 is withdrawn.** Two
things do differ and a consumer must carry both: the server seeds from `actor+0xa5`, the **high
byte** of the `u16` sight, so a **hero's sub-cell sight is truncated for the AI and not for the
fog** (up to 40 cells at sight 6 — `AI-SIGHT-094`), and its playable rectangle is one cell narrower
on every side (`TERR-FOG-118`).

**The playable rectangle** (`TERR-SIGHT-116`) is four bytes at `world+0x58ee0..0x58ee3` written by
the map load `FUN_00548550` as `(8, 8, W-9, H-9)`, with the same corners packed as words at
`+0x58ee4` = `0x808` and `+0x58ee6`. Every sight ring cell is tested against them; the compares
are **byte-wide** while the height and visibility indices are exact 32-bit, so on a 256-wide map a
true column of `-9..-11` wraps to `245..247`, passes, and reads the previous row.

### The order one frame is painted in, and everything that bounds a drawn silhouette

`TERR-SPR-137…TERR-SPR-141`, `TERR-FOG-142`, `TERR-LIGHT-143`, `TERR-SPR-144`.

`FUN_00407b1a` is **ten cell sweeps back to back**, not one composite pass. Its lowest back-edge
target is `0x004094ca`, so the 1570 instructions before that — rect/scroll arithmetic and the
terrain painting — contain no loop: **the terrain is finished before the first drawable is
touched**, and the shroud sweep is the tenth and last, after every sprite (`TERR-SPR-137`,
`TERR-FOG-142`).

```
sweeps 1..9  grid CMapView+0x8c / +0x90 / +0x94 / +0x98 / +0x9c,  all indexed
             idx = (row+3)*(visCols+6) + (col+3)
             outer = row,    -4  ..  visRows + 7      ascending
             inner = column, visCols + 3  ..  -4      DESCENDING
sweep 10     the shroud:  outer = column 0..visCols, inner = row 0..visRows+3, ascending
```

So within a sweep a cell can only be overpainted by a **greater row**, or the same row and a
**smaller column** — nothing is depth-tested, and the walk order is the visibility rule
(`TERR-SPR-138`). Two of the sweeps composite several planes inside one cell (sweep 4 issues seven
dispatches per cell); the `+0x90` plane's **shadow** (`vt+0x2c`, sweep 5) and **body** (`vt+0x28`,
sweep 8) are two whole sweeps apart, with byte-identical six-test guard chains, so on that plane
*every shadow is drawn before any body* — a body covers the shadow of a drawable above and to the
right of it (`TERR-SPR-139`).

**What bounds a silhouette.** One global device-space clip rectangle,
`0x005e4408` / `0x440c` / `0x4410` / `0x4414` = left/top/right/bottom, set from `*(CMapView+0xf4)`
at each of the four block boundaries by `FUN_0044ca60`, which **replaces** — there is no
`IntersectRect`, no clamp against the surface, and no per-cell destination window anywhere on the
draw path. `FUN_0044e460`, the sheared silhouette blit
`(dstX, dstY, w, h, src, level, shear)`, tests the box `[dstX - span, dstX + w + span] ×
[dstY, dstY + h]` with `span = |(shear*h) >> 16|`; entirely inside takes an unclipped loop,
otherwise it tries a whole-sprite reject on the same four edges and then clips **per row** (against
top/bottom) and **per pixel** (the sheared X against left/right). That is the opposite of the
terrain blitters, which drop a cell whole when its 32-px X range is not entirely inside
(`TERR-SPR-140`).

**Fog does not clip a shadow.** The only per-draw shroud test in the image is `FUN_0040d027`, and
it has **exactly one caller** — sweep 4's inlined object draw, where a true verdict skips the
object's shadow *and* body on one boolean. No unit shadow is subject to it. A silhouette crossing
fog is therefore drawn whole and then overpainted by sweep 10, so its visible edge is the shroud's
**per-vertex gradient**, never a cell boundary; on a fully visible cell the shroud draws nothing
and the shadow stands as drawn (`TERR-SPR-141`, `TERR-FOG-142`, `TERR-FOG-083`).

**Overlapping shadows compound.** The blit writes `table[level][dst >> 3]` back over the
destination with no guard against a second application and no per-pixel state that could hold one,
so where two silhouettes cross the ground is `((16-L)/16)²` rather than `(16-L)/16`
(`TERR-LIGHT-143`, on `TERR-FOG-084`'s law). A consumer that unions shadows into a coverage mask
and darkens once draws something the original never draws.

**Customisation (G2).** `[0x005bceec]` is the shadow switch: seven reads image-wide, all inside
`FUN_00407b1a`, each gating exactly one shadow draw and nothing else, so clearing it removes every
silhouette and changes no body, bar, terrain cell or shroud cell. It is a runtime global — no
shipped file's bytes move (`TERR-SPR-144`).

**Open:** the runtime value of `CMapView+0xf4`; the meaning of `drawable+0x78`, the field every
sweep filters on; what sweeps 2, 6 and 7 draw.

### A unit places itself — and its shadow places itself differently (`TERR-SPR-047`, `TERR-SPR-065…TERR-SPR-067`)

A **unit** does not use the cell-centre destination above. It is dispatched twice per frame over the
drawable grid `CMapView+0x90`, at `idx = (row+3)*(visCols+6) + (col+3)`, both times under
`drawable+0x78 == 0`, and both passes ignore the `(col, row, …)` they are handed:

```
vt+0x2c  FUN_0045bf00   the SHADOW   third argument = the CMapView+0xc0 four-corner altitude
vt+0x28  FUN_0045b3f0   the BODY     third argument = the CMapView+0xb0 light level

anchorX = frameW/2 + (CenterX - Width /2)      class +0x34, +0x2c   (units.reg record, REG-UNITS-049)
anchorY = frameH/2 + (CenterY - Height/2)      class +0x38, +0x30
    frameW/frameH via vt+0x20 / vt+0x24 at the DRAWN frame, in BOTH passes

BODY     dstX = unit+0x60 - anchorX                                    (0045bc58..0045bc74)
         dstY = unit+0x64 - anchorY - unit+0x10 - unit+0x68
SHADOW   dstX = unit+0x60 - anchorX - T                                (0045c63e..0045c65b)
         dstY = unit+0x64 - anchorY            - unit+0x68
             T = ftol(tan(theta) * (floor(frameH/2)+floor(Height/2)-CenterY)),  a PIXEL
             count loaded at 0045c645 -- NOT the 16.16 slope, which lives one slot away
             at frame base+0x18 and is pushed as blit argument 5 (TERR-SHDW-130)
```

**The two differ by exactly `(sunShear, −unit+0x10)`.** One rule will not serve both: applied to the
body, the shadow's rule puts every unit `unit+0x10` pixels too low; applied to the shadow, the body's
rule stands the shadow on the unit's feet instead of on the ground and drops its lean.
`unit+0x68` is subtracted by all four CUnit drawing routines (body, shadow, `vt+0x38`, the HP/mana
bars) and `unit+0x10` by every one except the shadow, so `+0x68` belongs to the plane the shadow lies
in and `+0x10` lifts what is drawn *at* the unit off it (`TERR-SPR-067`; the interpretation is
Medium, the reader split is not). Neither field has a writer inside `0045b000..0045f000`.

Frame selection is **duplicated**: each routine has its own nine-arm switch on `unit+0x74` with its
own jump table (body `0x0045bedc`, shadow `0x0045c918`), same arm shape, and the same arithmetic arm
for arm — so a shadow is a silhouette of the same pose. The one divergence is the states-2-and-4
default arm, which draws the class id in the shadow and the light level in the body: garbage in both,
which is the standing reason to think those two states never reach a draw.

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
draws the **body** (`TERR-SPR-065`, correcting `TERR-SPR-048`).

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
selected on `class+0x98` — `0` → one of 16 shared tables `[0x005eb65c][sprite[+8]]` over a 16 × `0x400`
palette blob read into `0x005eb6a8` at `units.reg` load; `1` → the class's own `class+0x9c`; `> 1` → a
per-owner `class+0x9c + (owner−1)*4`. Corpus: the 53 `terrain.3d` tiles share **one** palette and
**0 of 1373** palette-bearing `.256` sheets carries it, so **a sprite's ramp must be built from that
sprite's own palette**.

**Reimplementation guidance.** Light every sprite through its own palette's 16-row mode-2 table at the
cell's `+0xb0` level. Drawing sprites unlit is not a small error: the shipped swordsman frame means
104.9 luma at the engine's row 3 and 65.2 with no table, against 106.1–118.3 for the terrain in the
same crop (`tools/terrgfx -mode sprlightart`).

### Terrain output range (`TERR-LIGHT-063`)

The mode-3 mapping above, measured end to end against the one palette all 53 tiles share, with the
per-vertex level census re-run on the corrected base (`[30..70]` over 38 maps / 859 768 interior
vertices at θ = 0.78539815):

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

## Structure placement — where a building stands (`TERR-STRUCT-100…TERR-STRUCT-107`)

A `structures.reg` placement is **not** the unit/object shape above. It has no canvas and no
anchor pixel; its anchor is a **cell** and every frame of its sheet is exactly one tile
(`TERR-STRUCT-100`, `SPR256-STR-040`).

```
anchor cell (ac, ar) = (fine.x >> 8, fine.y >> 8)      fine = cellByte*256 + 128 on 3141/3141
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
back row. 1284 of 3141 shipped placements have such an overhang.

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

**Order** (`TERR-STRUCT-104`): the drawable registers into `CMapView+0x94` over its whole
footprint rectangle and the renderer's cell walk reads that plane twice — an early pass
taking `Flat != 0`, then the main pass taking `Flat == 0`, immediately before the unit plane
`CMapView+0x8c` at the same cell. Both passes are skipped by the fog state `+0x78`: the flat
pass needs `< 2` (seen at least once), the main pass needs `== 0` (in sight now).

**Ownership never reaches the sprite** — only the minimap blip `vt+0x34` (`TERR-STRUCT-107`).

## Open questions

- ~~The exact per-hour dawn/dusk **colour schedule** of the sun model.~~ **Closed by `TERR-LIGHT-119…TERR-LIGHT-128`**
  (`TERR-LIGHT-119`…`-128`): all six arms are in the table under *Sun* above, the magic division
  is /120 rather than /60, and nothing is read from a table. `TERR-LIGHT-014`'s two intensity
  figures are the cycle-off/day values only and are retracted as a description of the fields.
  *(The `Δh` neighbour per axis is closed by `TERR-LIGHT-028`.)* **Still open there:**
  why the object path's shroud level `0x5eb49c` is exactly twice the unit path's `0x5eb4a0`.
- ~~Which **θ a live session actually runs at** — whether shipped play has the day/night cycle
  enabled (`DAT_005eb528`) or takes the 0.78539815 default.~~ **Closed by `TERR-LIGHT-108`**
  (`TERR-LIGHT-108`): it is the `ShowTimeFlow` option, default 1, so the default session sweeps and
  the `0.78539815` row of the θ→level-range table is the *disabled* case, not the shipped one. What
  remains open is narrower: **what a given install's profile or savegame contains**, which is not a
  fact about the image.
- ~~The **CMapView scroll clamp** — how far the camera may scroll, and therefore whether the outermost
  1-cell ring ever becomes a *visible* pixel or always stays inside the culled over-scan margin.~~
  **Closed by `SESS-VIEW-030`**: the band is `8 <= origin <= dim − 8 − span` on each axis,
  enforced by four routines with identical arithmetic, so the outer 8 cells are never the viewport's
  origin and the outermost ring is reachable *only* through the draw's own 3–4 cell over-scan. The
  symbols this file writes as `scrollX`/`scrollY` are **`CMapView+0x5c` / `+0x60`, in cells**, and
  the span they are bounded against is `+0x64` / `+0x68` — `(rect.right − rect.left)/32` and
  `(rect.bottom − rect.top)/32`, i.e. **15×15 cells at 640×480, 20×18 at 800×600, 27×24 at 1024×768**,
  with `+0x68` shortened while an inventory or spell-book panel is open (`SESS-VIEW-028`). What a
  mission *opens* at is `clamp(heroCell − span/2)` (`MISSION-VIEW-019`).
- What the **uncomputed brightness ring actually contains at runtime**: the engine never writes it and
  the allocator never zeroes it, but a first-touch heap page may still arrive demand-zeroed from the
  OS. Not observed on the running game (`TERR-EDGE-025`).
- The **8-bpp** display path builds a different table (mode 3 is the 16-bpp one) — not covered.
  *(The three further blitters that index the same geometry step table are identified by `TERR-FOG-037` and
  are the shroud pass, not other pixel formats — see "Sprite placement" below.)*
- Why the drawn edge (step table) and the engine's own hit-test edge (exact lerp) disagree by up
  to 3 rows — reported, not explained (`TERR-GEOM-036`).
- Behaviour for `|Δh| >= 128` across a cell edge: the step table has no such row. Unreachable on
  the shipped corpus (max 127), uncharacterised.
- ~~Whether **anything** reads the map-stored type-0 light fields, given the lighting path
  overwrites them.~~ **Closed by `TERR-LIGHT-149`, `TERR-LIGHT-150` and `TERR-LIGHT-151`**
  — see "The map's stored light fields" above. Nothing reads them, bounded by a pointer
  reachability search whose blind spots those rows enumerate. What remains open is what the
  values *mean*, which no reader in this image can answer.
- What sets **bit 3 of the per-player fog flag** for a player other than the local one. The setter
  is one located message arm (`004162ec`) whose opcode is unread, so "shared vision in multiplayer"
  is the shape of the answer and not the answer (`TERR-FOG-080`).
- What each of the **six pointers** at cell-record `+0x14`…`+0x28` holds. The loop that quadruples
  the cost per non-null one, and the fourth one's extra block, are now sourced; the fields are not.
- ~~Whether the four-corner `0xc000` gate is ever satisfied in play.~~ **Closed by `TERR-TILE-079`, `TERR-FOG-080` and `TERR-FOG-081`** — see
  "Bits 15..14 are the FOG OF WAR" above.
- ~~Unit/sprite lighting uses the same sun globals and the same builder in the 16-level mode 2
  (neutral row `nLevels/2`), via a separate path (not covered here).~~ **Closed by `TERR-LIGHT-059…TERR-LIGHT-064`** —
  "Sprite lighting" above. What remains open in that area: what `class+0x98` counts, and so which of
  the 34 unit classes takes the shared, class-owned or per-owner table; what the `[0x005eb65c]+0x40`
  and mode-5 greyscale overrides at `0045b9a5`/`0045b9cb` are for; which node supplies the 16 KB
  palette blob `FUN_004ca140(0x5eb6a8, 0x4000)` reads; and what the second `vt+0x28` dispatch
  (`CMapView+0x9d4`, gated `+0x9e0`, arguments `(0,0,0)`) is for.
- ~~**What a unit's own `Draw` does.**~~ *(Closed for the frame and destination by `TERR-SPR-047…TERR-SPR-048`, and
  re-labelled by `TERR-SPR-065`: `vt+0x2c` = `FUN_0045bf00` draws the unit's **shadow**, `vt+0x28` =
  `FUN_0045b3f0` its **body** — `TERR-SPR-065`.)* Still open from that reading: the second placement
  model beside it, the unit's screen *rectangle* built by `OffsetRect` from its own pixel position
  minus two further terms.
- ~~The shroud's actual darkening curve: the 16→16-bpp remap table behind `[0x005e8420]` … was not
  read (`TERR-FOG-037`).~~ **Closed by `TERR-FOG-082…TERR-FOG-089`** — see "Fog of war" below. The table is **17 rows**
  and row `L` is `out_channel = (in_channel × (16 − L)) >> 4`; the same table is what the sprite
  silhouette passes index, so that half is closed with it.
- The `vt+0x34` overlay's blend rule — it reads the destination and `[0x005eb578]` alongside the
  source. `TERR-LIGHT-059` locates the consuming instructions (`0044dde0`, `0044eca0`) and does not read
  them; this is the standing `spritesb` residual (`SPR256-OVL-014`).
- ~~The altitude-driven **vertical mesh offset** … is runtime-render detail, out of scope here.~~
  **Closed by `TERR-GEOM-031…TERR-GEOM-036`** — it was not out of scope, it was the missing half of the render
  ("Cell geometry" above). Per-display-mode 8-bpp palette quantization remains uncovered.
- The legacy non-3d `terrain\tile*.bmp` fallback set is noted but not corpus-exercised.
- The type3 static-object placement layer's own graphics are a separate render path
  (`ALM-CLS-035`).
- ~~**Movement (`TERR-PASS-049…TERR-PASS-051`, `TERR-COST-052`, `TERR-PASS-053`):** which `units.reg` class is a ground mover … what fills `vt+0x1c`~~
  **Answered by `TERR-MOVE-054…TERR-MOVE-057`, and the answer is that there is no such binding**: both values are
  per-instance bytes of the simulation actor, written `1` by its constructor and by nothing else,
  and the `units.reg` class array is never referenced from the simulation module ("Who is asking"
  above). ~~What remains open in the same area: what fills those two bytes on the (de)serialization
  path, and therefore whether a mask other than `0x41` is ever produced at all~~ — **answered by
  `TERR-MOVE-057`**: the `Data.bin` param streamers, and yes (above). Also answered by
  `TERR-MOVE-058`: the `+0x70 → +0x3c → +0x44` speed override is a per-owner runtime byte, not the DB.
  **Block bit 2 is answered by `TERR-STRUCT-068…TERR-STRUCT-072`** ("Structures" below) — it is *an object
  occupies this cell*, written and erased only ever paired with bit 0. **Bits 3 and 4 remain open**
  beyond "part of the border value `0x1f`", with the narrower question that bit 4 is a live runtime
  flag every recompute deliberately preserves and no located instruction originates.
