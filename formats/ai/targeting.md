# Target acquisition and sight

[Reference](format.md)

## Target-selection sequence

1. **Population — sight.** Clear `world+0x82ef0` (`0x10000` bytes, one per cell). Stamp the
   deciding actor's line of sight into it (below). If the actor is AI-owned and remembers being
   attacked, mark that cell too — for 20 ticks. Then walk the global on-map actor list and keep
   every actor standing on a marked cell.
2. **Relation.** Drop every candidate `c` for which `matrix[me.Player+0x04][c.Player+0x04] & 1`
   is 0. Drop an invisible candidate (`c+0x144 & 0x8000`) unless some member of the decider's
   group is within that member's own `order+0x71` of it.
3. **The dead.** Move survivors with `health < 1` to a second list. If nothing living is left
   and that list is not empty, move it back: a unit with only corpses in view targets a corpse.
4. **Selection.** Minimise the footprint-aware edge distance `d` (1 = touching), with `d + 1`
   for a flying candidate when the decider does not fly; break ties on the 16-way turn cost from
   the mover's current facing. Seed the best distance at `reach + 1`, and **discard the winner
   if its distance exceeds `actor+0x12c`**.
5. **Outcome.** A winner sets `order+8 = 6`, `order+0xc = target`, `order+0x14 = reach`. No
   winner: a human participant's unit tries to heal; an AI-owned unit goes to the guard state,
   which walks it back to its post whenever it is farther from the post than `actor+0xa5`.

**Sight decides the population and reach decides the pick.** An engine that acquires on sight
alone starts fights the game does not; one that scans only within reach acquires through walls,
and acquires things it cannot see. Neither radius is a substitute for the other.

### Prismatic Spray is a distinct consumer of the group population (`AI-SPRAY-266`, `AI-SPRAY-267`)

`FUN_0053ddd0` is the fifth direct caller of the group builder. It consumes shared list A head to
tail and then list B head to tail; it does not run a radius search. The ordinary secondary gates are
therefore the builder's group-wide visibility, first-member diplomacy and See-invisible rule, then
its living/corpse split. The primary spell target is separate: the selector alarms it and performs
the hostility flip before the builder call, appends it first, and skips an equal secondary later.
It can therefore survive conditions that would remove it as a secondary.

Prismatic rank differs from ordinary acquisition. List A uses
`((edgeDistance<<8)+turnCost)&0xffff`; list B uses
`((((edgeDistance<<8)+turnCost)&0xffff)<<8)`. Repeated strict-minimum scans preserve source order on
equal stored scores. Selection first requires a score below 65530 and append later requires the
saved threshold below 65000. A low enough B score can therefore compete beside living A; whether a
valid layout can supply the required edge distance is Unknown. If there are no living entries, the
builder has already moved B back into A and the only-corpses population uses the A formula. The
spell's capped byte limits the final output list after ranking; it is neither actor reach nor order
`+0x14`, and the primary consumes one position on the ranked path. When the builder returns empty A, the selector instead
appends the primary and returns before reading the cap. Its stack pointer region holds 100 entries
and its winner/threshold regions ten; neither population overflow nor a custom cap above ten is
checked. Candidate scores are session scratch at `+0xd74`, whose capacity is Unknown.

### The sight stamp, complete (`AI-GROUPSEE-068`, `AI-LOS-087`…`AI-LOS-091`, `AI-SIGHT-092`…`AI-SIGHT-094`)

The sight object is embedded at `world+0x58ee8`. Its init runs in each of the four world
constructors, **before the map is loaded**, and builds three tables that depend on
`k = [Scanning] ScanShift` of `World\Data\map.reg` (ships **7**) and on nothing else — no terrain,
no actor, no tick. Two of them serve the stamp: `step`, `+0x22000`, and `cost`, `+0x28000`, both
`u16`/`i16` over a 41×41 window inside a 64×64 grid, addressed as `base + (a<<7) + 2b` with `a`
the column offset and `b` the row offset, both `0..40` around a centre of `(20,20)`.

```
step[a][b]  = the cell one Bresenham step toward the centre, as (i8 dx, i8 dy) in {-1,0,+1}
              zone by slope:  j <  i>>1 -> column step | j > 2i -> row step | else diagonal
              the builder's four quadrant mirrors collide on the j==0 axis, so it repairs
              cells (+1,0) and (-1,0) with four literal stores at the end
cost[a][b]  = ftol( (1<<k) * sqrt(i*i + j*j) / max(i,j) )   -- the mean length of one step
              along that ray, in 1/(1<<k) cell.  k=7: 128 on the axes, 181 on the diagonals

stamp(actor at cell C, scanRange = actor+0xa5):
  acc[*][*] = 0 ; acc[20][20] = (1 << (k-1)) + (scanRange << k)      -- a budget in 1/128 cell
  vis[C] = 1 ; alt = (i8) height[C]                                  -- height = world+0x9451c
  for r = 1 .. 19:                            -- Chebyshev rings; ring 20 is unreachable
     allBlocked = true
     for each of the ring's four edges, 2r+1 cells each, in edge order:
        skip the cell unless 8 <= absCol <= W-9 and 8 <= absRow <= H-9   (world+0x58ee0..3)
        p = step[a][b]
        v = acc[a+p.dx][b+p.dy] - cost[a][b] - (i8) height[cell] + alt
        acc[a][b] = v                         -- stored EVEN when it blocks
        if v <= 0: continue                   -- blocked: NOT pruned, the ray is not cut
        vis[cell] = 1 ; allBlocked = false
     if allBlocked: stop
```

The stamp uses a height-sensitive budget:
`acc(n) = seed − Σcost − Σh(cell) + n·h(observer)`.
The observer's height is added at every step; each traversed cell's height is
charged once. On flat ground the region is a disc of radius `scanRange`
(145 cells at 6). A single high cell is charged once and does not cast a
continuing shadow. A negative-budget cell is still written and the march
continues. The minimum step cost is `1<<k = 128`; installed altitude bytes
`0..127` cannot restore a negative budget, while an authored byte `>=0x80`
can. All three stamp sites center the window on the actor's own cell.

Drawn fog uses the same predecessor/cost tables and an algebraically equal
seed for whole-cell sight. The shared algorithm has two parameters:
— AI-SIGHT-093

- **the seed's width.** The fog reads the whole `u16` at `actor+0xa4` — sight in 1/256 cell — and
  shifts it right by `8 − k`; the sight stamp reads only its **high byte** `actor+0xa5`, the
  whole-cell radius. A monster's sight is whole cells and the two coincide; a hero's is
  `ftol(((mind + reaction)/25 + 4) × 256)` and is not, so **the AI grants a hero less sight than
  his own fog shows him** — 105 cells against 145 at `actor+0xa4 = 1535` (`AI-SIGHT-094`).
- **the playable rectangle.** `8 .. W−9` here against the view's `7 .. W−8`, and byte-wide compares
  here against 32-bit ones there (`TERR-FOG-118`).

`actor+0xa4`'s six writers, and the two forms of sweep that see them, are `AI-SIGHT-092`; a
`disp:a5` sweep sees two of the six.

**A unit only closes distance when it is given a pursuit order** — `order+8 = 5`,
`order+0xc = target`, `order+0x14 = reach` as the stop distance. Acquisition never produces one;
the guard and engage arms do.

**Two clauses suppress engagement entirely**, both requiring all three of `actor+0x4c & 4`,
`actor+0x12c < 2` and a human-participant owner (`Player+0x28 == 0`).
