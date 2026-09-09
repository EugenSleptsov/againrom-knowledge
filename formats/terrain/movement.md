# Passability and cell records

[Reference](format.md)

## Tile word → movement (`TERR-PASS-049…TERR-PASS-051`, `TERR-COST-052`, `TERR-PASS-053`)

The same word drives a second, wholly separate reading: whether a **ground unit** may stand on the
cell. It is the sim side, not the render side, and it uses a different bit split — the whole
decision runs on `w & 0x3ff` plus bit 13, so **bits 10–12 and 14–15 never reach it** (0 of 880 704
shipped cells set any of them). Sim-side grid layout, the ingest and the persistence are specified
in [ALM](../alm/format.md) → "Runtime passability"; what belongs here is the tile word's own part.

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

A ground mover is blocked by tile bit 13 (`w&0x2000`), Mountain class 8,
raw water bits `(w&0x300)==0x200`, a nonzero type 3 object cell, or the
eight-cell border. Bit13 also controls the render-path dirt composite;
that graphical use does not replace the other movement blockers.

### Who is asking — the mover (`TERR-MOVE-054…TERR-MOVE-057`)

The `n×n` footprint and the mask are **not** properties of a `units.reg` class. Both come from the
simulation actor instance the predicate is called on (base vtable `0x59c3c0`, derived `0x59c448`,
`0x59c4d0` — the only family that allocates the `0xb4`-byte mover and stores it at `+0x154`):

```
vt+0x1c  FUN_00523210  MOV AL,byte ptr [ECX + 0x49] ; RET      the footprint side n
vt+0x20  FUN_00523230  MOV AL,byte ptr [ECX + 0x4a] ; RET      the mask selector, 1/2/3
base ctor FUN_004f30a2 004f30e0 MOV byte ptr [EAX+0x49],1
                       004f30e7 MOV byte ptr [ECX+0x4a],1      constructor defaults
```

The constructor's1 values are defaults. Data.bin Units/Humans streamers
overwrite tokenSize/movementType when the parameter is not−1. Installed
Ghost/Bee use domain 2 (mask0x44), Bat_Sonic/Dragon use 3 (mask0x82), and
footprints include sizes 2 and 3. The universal 1×1/0x41 consequence is
retracted by `TERR-MOVE-057`. The units.reg drawable class array is not a
simulation input: TileSize belongs to drawable CUnit, while Z selects the
drawable class. The simulation-domain clause of `REG-UNITS-061` is retained
as amended. — TERR-MOVE-054, TERR-MOVE-055

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
back. With installed maps, v0=16 and SpeedMultiplier=8, the resulting
speed spans 4..32 (mode16); these are data-derived values within the clamp.

## Cell records

[Cell records and Building attachment](cells.md) specify the runtime payload,
registration, recompute, masks, detach and saved state.
