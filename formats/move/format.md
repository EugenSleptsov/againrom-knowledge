# MOVE — unit movement and path selection — specification (partial)

Level 3. Promoted, evidence-backed claims only. Search, costs, termination, extraction,
reservation, tick and refresh are `MOVE-SEARCH-001`…`MOVE-REFRESH-012`; tick ordering and
save/load regrouping are `MOVE-TICK-013`…`MOVE-TICK-017`; substitute goals are
`MOVE-ALT-018`…`MOVE-ORDER-023` (partially retracted); movement domains are `MOVE-DOM-024`…`MOVE-DOM-028`;
the rate and formation gates are `MOVE-RATE-029`…`MOVE-GROUP-037`. Ledger: `claims/move.md`.

**Status: partial (◐).** The search, its cost model, its termination, the substitute goal, the route
extraction, the reservation mechanism, the sub-cell step, the refresh policy and the tick order are
read at instruction level, as are the rate law, the turn and the tick that drives them. Not
specified: what `FUN_0054abb0` does on a cell boundary (`FUN_00544d00` and `FUN_00545230` are
specified below by `TERR-CELLREC-146` and `TERR-FOOTPRINT-147`); when a
group order writes the group speed term (`grpAI+0x44`) rather than leaving it 0; the verdict table
of the blocked-cell decision
`FUN_00533540`; the exact entry cell of the contact-ring picker (its *shape* is read, its
quadrant→edge mapping is not checked cell by cell); the order layer above movement; whether
buildings tick from another container.

The movement predicates do not consult a Building's HP. If lethal application
leaves its cell reference, mask and other cell inputs unchanged, the next lookup,
recompute and movement predicate retain the same interpretation
(`UNIT-STRUCTNEXT-079`). Whether those inputs actually survive to that next
consumer is not established by the bounded lifetime search
(`UNIT-STRUCTBOUND-081`). Ruin appearance alone does not establish passage release.

This is not a file format. It is the simulation area sitting on the three planes
(`formats/terrain/format.md` → passability) and the mover block (`TERR-MOVE-054…058`).

## At a glance

```
per unit, per move order (FUN_00548f70, once per tick):
  mid-transit?          -> advance the sub-cell position and return
  target changed, or  N dynamic searches since the last static one?
                        -> STATIC search  (plane world+0x10000, no occupancy) -> static route
  no dynamic route, or M ticks since the last dynamic search?
                        -> next static waypoint adjacent AND blocked?
                             -> wait, facing it. no step, no search.
                           else DYNAMIC search (plane world+0x20000, with occupancy)
                                aimed a few waypoints down the static route -> dynamic route
  any search that fails to label its goal
                        -> substitute the labelled cell with the SMALLEST LABEL near the
                           goal (or near the target actor) and route to that instead
  step: claim the next cell on the shared plane, turn toward it, or advance the sub-cell
        position by the frozen per-axis step -- exactly once per SUB-TICK
```

## The movement domain (`MOVE-DOM-024…028`)

One byte, `actor+0x4a`, read through slot `+0x20` of all three simulation-actor vtables
(`FUN_00523230`). It is **not** a passability property with side effects — six things branch on it:

| what | rule | where |
|---|---|---|
| block mask | `1 → 0x41`, `2 → 0x44`, `3 → 0x82`; anything else stores **nothing** and leaves the constructor's `0x41` | `FUN_0054b120` |
| occupancy bit | `< 3` sets/tests bit 6, `== 3` bit 7, on the dynamic plane | `FUN_0054abb0`/`ac70`/`ad20`/`af40` |
| step cost | `== 1` reads `cost[dst]`; every other value takes flat 2 straight / 3 diagonal | 9 sites, `MOVE-COST-002` |
| speed | `== 1` divides by `cost[cell]` and tilts by height; 2 and 3 take the raw class speed; 0 and `> 3` give **0** | `FUN_0054d210`, `FUN_0054a620` |
| AI targeting | a domain-3 candidate counts **one cell farther** to a decider that is not domain 3 | `0052e18f`, `AI-ACQUIRE-002` |
| corpse | `> 1` slams the decay counter to −1000 → immediate teardown, so no corpse | `004f3907` |

**The mask is a cache, written twice in the image and never refreshed**: `FUN_00545c00` (the mover
constructor, `0x41`) and `FUN_0054b120`'s three arms, whose only two call sites are inside object
construction. It survives a save because `Unit::Serialize` hands the whole `0xb4`-byte mover to
`CArchive::Read`/`Write` (`FUN_0054d4c0`) before its own storing/loading branch.

**What each domain may enter.** Terrain and objects are *different bits*, so the domains are not
three grades of one thing. Counts are the 38 shipped maps' static plane after the structure pass,
880 704 cells:

```
                                    dom 1  0x41   dom 2  0x44   dom 3  0x82
terrain      water                      77 845             0             0
             Mountain (class 8)        117 572             0             0
             tile-word bit 13            3 944             0             0
non-terrain  .alm type-3 object         69 676        69 676             0
             structure (Passability)    13 527        13 527             0
             8-cell border             158 976       158 976       158 976
             ---------------------------------------------------------------
             blocked                   441 540       242 179       158 976  (50.1 / 27.5 / 18.1 %)
             + runtime: bit 6 occupant  yes           yes            no
                        bit 7 occupant  no            no             yes
```

So: **domain 2 crosses water and mountain and is stopped by every object, building and ground
occupant; domain 3 is stopped only by the map border and by another domain-3 occupant.** They
disagree on 83 203 cells. 8-connected, the free set is in more than one piece on 23 of 38 maps for
domain 1 and on **0 of 38** for domain 3 (whose one piece is the whole interior).

**Nothing narrows the verdict afterwards.** The predicate is `block[cell] & mover[5]` over the
`n × n` footprint and nothing else: no height term anywhere on the search path (the height plane
has 6 references image-wide, none in the search), no corner rule in the extraction, no plane test
at step time, and no order-time check — the block planes have no reader outside
`0x00523be0…0x0054f680`, so no interface or campaign routine can run one.

**Who is what.** Every one of the 210 `Data.bin` Humans rows carries `movementType = −1` → ground.
Of 57 named Units rows, 16 are non-ground: Ghost and Bee (2), Bat_Sonic and Dragon (3), four
difficulty variants each. The two domain-3 classes are exactly the two with `units.reg` `Z != 0`,
i.e. the two `CAirUnit`s. The hireable roster is Catapult, Ballista and 13 Humans, so **no
player-ownable unit is a non-ground mover** — with one open exception, `Control Spirit`, which
spawns a caster-owned `Ghost` (domain 2).

## The search — `FUN_00541e80`

A **double-buffered label-correcting wave**. Not Dijkstra (no extract-min), not A\* (no `h`).

```
label:    u16 plane at world+0x30000, 65536 cells, 256-stride, 0xffff = unlabelled
frontier: two lists used alternately, both 4096 entries, counts u16 at
          world+0x54008 and world+0x5400a
            footprint 1  -> u16 packed cells at world+0x545b4 / world+0x565b4
            footprint >1 -> byte pairs   x world+0x50008 / y world+0x51008
                                     and x world+0x52008 / y world+0x53008
seed:     label[src] = 0, src pushed to list A
generation:  for every cell in the current list, for each of the 9 neighbours
             (centre included, scan order dx = -1,0,+1 outer, dy = -1,0,+1 inner):
               skip if (blockPlane[c] & mover.mask) != 0 over the mover's n×n footprint
               g = label[cur] + step(cur -> c)
               if g < label[c]:  label[c] = g;  append c to the other list
             then swap the lists
stop when:   label[goal] != 0xffff  (tested between generations)
          or the current list is empty
          or generations >= budget
```

**Step costs** (`MOVE-COST-002`), where `cost[]` is the byte plane at `world+0`:

| mover | straight | diagonal |
|---|---|---|
| `movementType == 1` (ordinary ground) | `cost[dst]` | `cost[dst] + (cost[dst] >> 1)` |
| any other movement type | `2` | `3` |

The diagonal is 3/2 of the straight step **truncated**, and the cost byte is the *destination*
cell's. There is no heuristic and no distance term in any label.

**Budget** (`MOVE-TERM-003`), with `D = max(|Δx|, |Δy|)`:

| search | footprint 1 | footprint > 1 |
|---|---|---|
| static | `max(StaticScanAhead, D>>2) + D`, or **1000** under the override below | `StaticScanAhead + D` |
| dynamic | `max(DynamicScanAhead, D>>2) + D` | `DynamicScanAhead + D` |

**The 1000-generation override is not an edge case — implement it or a human player's units will
refuse long detours.** It applies to the static footprint-1 search when both hold
(`005422fe`…`0054235d`):

1. `actor+0x14` — the **owning `Player`** — has `+0x28 == 0`, which is true exactly when a *human
   participant* owns the unit. Every scenario-authored owner carries 1, or 2 for a group whose
   `HumanFriend` key reads `"Yes"`; the `Player` constructor's own default is 1, and the only
   store of 0 in the image is the session-join path (`UNIT-OWNER-009`).
2. the **goal's own `tokenSize × tokenSize` footprint** is free of the mover's domain mask
   `mover+0x5` on the static plane `world+0x10000`. A `tokenSize <= 0` skips the scan and takes
   the override.

Outside the override, one generation advances the wave by one ring, so the slack over the
straight-line distance is only `max(scalar, D>>2)` rings. A detour needing more fails, and the
search then substitutes a goal near the requested one (below) and extracts to that instead.
**There is no node budget, no frontier capacity check and no closed set.**

## The substitute goal (`MOVE-ALT-018…022`)

When the loop ends with `label[goal] == 0xffff`, `FUN_00541e80` picks a substitute **itself** — the
two pickers are called from its own tail and from nowhere else in the image — and feeds it to the
ordinary route extraction as the destination. The substitute is **not** written back to the actor:
the ordered target (`mover+0x74`/`+0x76`) and the order block are untouched, so the next search
starts from the same request again.

```
label[goal] != 0xffff  -> ordinary extraction to the goal
staticFlag != 0        -> ring picker around the requested cell, limit = (D>>2) + 4   [altTarget ignored]
staticFlag == 0        -> altTarget != 0 ? contact-ring picker(mover, altTarget)
                                         : ring picker around the requested cell, limit = 8
picker returned 0      -> no route (static: the route list is freed)
```

`altTarget` — the search's eighth argument — is a pointer to a **target actor**, not a cell. It is
supplied only by the go-to-actor order (`FUN_005492a0` → `FUN_00549a90`) and is consulted only on
the dynamic branch.

**Both pickers choose by the same rule and read only one plane: the label plane `world+0x30000`.**
A candidate qualifies iff it carries a label, i.e. iff the wave that has just failed reached it —
which already folded in the plane it ran on, the mover's footprint and its mask. Among candidates
the **smallest label** wins, so the substitute is the cell *cheapest to reach from the mover* under
the step costs above, **not** the cell nearest the click. Neither picker contains or calls a
passability predicate, a cost read or an occupancy test.

*Ring picker* `FUN_0054bac0(world; mover, 0, cell, limit)` — rings `r = 1 … limit-1` around the
requested cell; per ring, `i = -r..r` and four probes `(x+i, y+r)`, `(x+i, y-r)`, `(x+r, y+i)`,
`(x-r, y+i)`; the **whole ring** is scanned and the strict minimum kept; the loop stops after the
first ring that yielded anything. The centre is never probed. The mover's footprint side is fetched
and branched on but the two branches are the same code. Returns the packed cell or 0.

*Contact-ring picker* `FUN_0054b420(world; mover, target)` — the box
`x ∈ [tx-nM-k, tx+nT+k]`, `y ∈ [ty-nM-k, ty+nT+k]` for `k = 0..7`, `nM`/`nT` the mover's and
target's footprint sides; at `k = 0` its perimeter is exactly the set of mover origins whose
footprint touches the target's without overlapping. The walk enters where the line between the two
actors' fine footprint centres crosses the box and proceeds in **both** directions to close the
ring; the strict minimum label wins. Returns the packed cell or 0.

On the **static** branch only, the substitute's Chebyshev distance from the request is compared with
1 — or 2 when the requested cell carries a cell record (`TERR-PASS-051` bit 5) with a non-zero
`+0x0c` — and exceeding it queues a UI message. It does not cancel the move.

## The route — `FUN_005436b0` (dynamic) / `FUN_005433a0` (static)

Walk downhill through the label field from the terminating cell to the seed, choosing at each step
the 3×3 neighbour minimising `label[nb] + step(nb -> current)` with the same four cost arms.
Constraints: the neighbour must be labelled, and `8 <= x <= W+8`, `8 <= y <= H+8`
(`world+0x50000` = W, `world+0x50004` = H). Passability is **not** re-tested and there is **no corner
rule** — a diagonal between two blocked cells is legal.

**Tie-break, and it is not symmetric:** the straight/centre arm accepts a candidate on `<=`, the
diagonal arm only on `<`. Among equal-cost neighbours the **last straight one in scan order** wins,
and a diagonal never displaces an equal straight.

Output: a doubly-linked list of 12-byte nodes on the actor — `+0x00` toward the goal, `+0x04` toward
the unit, `+0x08` the packed cell `(y<<8)|x` — built goal-first, the unit's own cell unlinked before
returning, consumed from the tail. A walk exceeding **1000** steps discards the whole route.

Two routes per unit, one per plane:

| | static (unit-blind) | dynamic (unit-aware) |
|---|---|---|
| plane read | `world+0x10000` | `world+0x20000` |
| tail / head / count | `actor+0x160` / `+0x164` / `+0x168` | `actor+0x17c` / `+0x180` / `+0x184` |
| free list / pool | `actor+0x16c` / `+0x170` | `actor+0x188` / `+0x18c` |

## Parameters — `data/map.reg` `[Path Finding]`

Read by `FUN_00547eb0`. The shipped file carries all seven equal to the code defaults, in both the
EN and RU releases.

| key (file clamps names to 15 chars) | world | default = shipped |
|---|---|--:|
| `SpeedMultiplier` | `+0x58db4` | 8 |
| `StaticScanAhead` | `+0x585b4` | 5 |
| `DynamicScanAhead` | `+0x585b8` | 3 |
| `StaticRefreshRate` | `+0x585bc` | 16 |
| `DynamicRefreshRate` | `+0x585c0` | 32 |
| `DynamicByStaticLookup` | `+0x585c4` | 3 |
| `StaticIsntNeeded` | `+0x585c8` | 5 |

## Coordination between units

The **dynamic block plane is the only channel**, and it carries two things: where units *are* and
where they are *going*.

- **Occupancy.** `FUN_0054abb0` ORs the mover's own bit (`0x40` ground for movementType < 3, `0x80`
  air for 3) over its n×n footprint; `FUN_0054ac70` clears it. The cell is recorded in `mover+0xa6`.
  A dynamic search brackets itself with clear/restore at the unit's own cell so it does not block
  itself.
- **Reservation.** Before stepping, `FUN_00549990` reads the next route cell into `mover+0x06` and,
  if it differs from `mover+0x80`, releases the old claim (`FUN_0054af40`) and marks the new one
  (`FUN_0054ad20`) — on a cell the unit has **not yet entered**. `mover+0x80` is the unit's intended
  cell. The pair is asymmetric: the claim ORs unconditionally, the release skips any cell inside the
  unit's current footprint, so releasing never un-occupies the ground it stands on.
- Every mask includes its own domain's occupancy bit — `0x41` = terrain bit 0 + ground bit 6,
  `0x44` = object bit 2 + bit 6, `0x82` = border bit 1 + air bit 7 — which is why the same predicate
  serves both planes, and why the static search is unit-blind: bits 6/7 are never set on
  `world+0x10000`.
- **Blocked next cell:** the unit turns to face it and waits. It never pushes, swaps or steps aside.
  The step routine itself tests no plane, so two units whose routes predate each other's claims can
  enter the same cell.
- **No priority.** The tick loop (`FUN_0050fca9` → the actor's `vt+0x18`) applies no sort and no
  priority key, and claims land in the shared plane immediately — so whichever unit the loop reaches
  first that tick takes the cell.
- **Formation, conditionally** (`MOVE-FORM-036`; `MOVE-ORDER-023`'s "no formation" is retracted).
  A group Move / Swarm-2 order runs one of two arms. **In formation**, each member is ordered to
  `target + (memberCell − groupCentroidCell)` — the offsets are kept at `ord+0x24`/`ord+0x26` — so
  the destinations differ by construction. **Out of formation**, and for `FUN_005370f0`'s group
  order 2, every member gets the same loop-invariant cell and spreads only because each one's own
  search fails at the crowded cell and substitutes independently against a plane that already
  carries the earlier movers' claims. Which arm runs is the same gate that decides the group rate
  term — see *When the group term is set* below.

## The tick order — what "first" means (`MOVE-TICK-013…017`, `MOVE-ID-016`)

The loop walks a pooled doubly-linked list embedded at `+4` of the manager at `[0x00609558]`
(12-byte nodes: next/prev/element; CPlex blocks of 10), head→tail. **The order is insertion
history and nothing else** — the family's only insert is AddTail; no sort, no AddHead, no
mid-list insert exists.

```
fresh map:    party first, then the map's type-6 records in record order
              (Humans, Units, Sacks tick here; no Building insert was found)
during play:  death       -> unlink (no hole), actor moves to the dead list *(world+0xc)
              spawn/summon-> AddTail (one sack creator is the actor tick itself)
              garrison    -> unlink;  return to map -> AddTail (loses its old position)
              owner change-> tick position unchanged (only player/group lists move)
save/load:    the tick list is NOT serialized. The stream carries player -> group -> actors
              (each list head->tail); the loader rebuilds the tick list as
                for each player (manager order):
                  for its list (groups in creation order, actors in group order):
                    AddTail, skipping off-map actors (actor+0x4c bit 3)
              => within-group relative order survives; the cross-player interleave does NOT.
                 A save/load cycle can change which of two contending units moves first.
```

The **runtime id** (`actor+0x04`, the SAV head id) is *not* the order: it is the lowest free bit
of the bitmap `0x62c7e0`, assigned at insert, freed (and the field zeroed) when a corpse reaches
decay stage 5, reused by the next spawn, and restored exactly across save/load (read back and
re-marked by the head serializer `FUN_00510e5c`). A consumer reproducing contention must keep the
**list**, not sort by id.

## The step

Position is `actor[4]` = `*(actor+0x10)`: `+0x00` x cell, `+0x01` y cell, `+0x02` packed cell,
`+0x04` x fraction, `+0x05` y fraction, `0x80` = centred. `FUN_00548c60` treats `(cell<<8)|fraction`
as one 16-bit value per axis and adds the signed per-axis step `mover+0xb0` / `mover+0xb1`; when the
cell byte changes it hands `FUN_0054abb0` the **old** packed cell, then writes the new position, then
calls `FUN_00544d00`; when `mover+0xac >= mover+0xaa` it snaps both fractions to `0x80`. A transit
takes `mover+0xaa = ceil(256 / |step|)` ticks. Speed never enters a label — two units of different
speed pick the same route.

### What the cell-boundary calls do (`TERR-CELLREC-146`, `TERR-FOOTPRINT-147`)

`FUN_00544d00(map, actor)` and `FUN_00545230(map, actor, x, y)` are the enter and the leave of the
map's **cell record**, the per-cell structure `formats/terrain/format.md` specifies. Both read the
actor's movement domain through `vt+0x20` and pick a slot from it: domain 1 or 2 uses the record's
`payload+0x04`, domain 3 uses `payload+0x08`. The per-cell entry refuses other
domain values; the detach statement here is scoped to domains1/2/3. `FUN_00544d00` reads the
footprint side once (`00544d12`) and writes the actor into **every** one of the `n x n` cells it
covers, through `FUN_00544ec0`, one call per cell; each call fails when that cell's slot is already
taken, and one failure stops further iteration without local rollback of earlier
cells or the mover+72/+82..85 caches. Only successful cell entries reach `FUN_005456d0`,
which recomputes that cell's cost byte and both block-plane bytes from the record — which is where
dynamic bits 6 and 7 come from, and where an area effect reaches the search at all
(`MOVE-AREA-038`).

The former all-or-nothing interpretation is withdrawn. A2x2 footprint can retain
its completed row-major prefix on refusal. The step caller stops its old-cell
detach loop on false but continues the position rewrite, and it does not test
entry's return before the center test. At center, the known dynamic-route cleanup
and progress3 completion can occur despite a refused destination entry. Local
selected serializers subsequently emit that state; no rollback or next-SAVE
reconciliation is established. — SAV-CELLFAIL-583, SAV-CROSSNEXT-585

Entry writes+82/+83 from cached low-byte cell-X/Y getter results
(`00544a10/00544a20`), then+84/+85 from the low bytes of full-X/Y getters
(`005449e0/005449f0`, hence Position+04/+05). Cell-X/Y were sampled before
`0054a620`; full-X/Y are sampled after it. Word+72 is that callback's AX,
not a proved numeric default; `0x1357` is only the probe cut. Those accessors are
cuts backed by `SAV-TOKENPOS-074`, not direct loads in the entry caller.
Successful detach instead directly loads Position+00/+01/+04/+05 into
+86/+87/+88/+89, in that order, after recompute/optional removal. It reloads
actor+10 for each source. Neither sequence establishes callback purity or an
atomic entry snapshot. — SAV-CELLFAIL-583, SAV-CELLLEAVE-584

Two consequences a consumer must reproduce. A unit of footprint side `n` is present in `n²` cell
records while it stands, so anything walking cells finds it `n²` times — that is what makes
`fire_ball`'s footprint-squared divide a normalisation (`MAGIC-FIREDIV-047`). And `FUN_005456d0`
assigns the dynamic byte from the static byte before rebuilding bits 6 and 7 from the record, so
occupancy written straight onto the dynamic plane for a cell that holds a record does not survive
the next recompute of that cell.

## The rate, and the clock it runs against (`MOVE-RATE-029`…`034`)

The rate is computed **once per cell transit**, by `FUN_0054d210`, from the actor's current facing
and cell, and is then frozen until the next cell.

```
dir  = ((facing + 0x10) >> 5) & 7          facing is a byte, 8 directions x 32 units
dx   = [ 0, +1, +1, +1,  0, -1, -1, -1]    world+0x58eb0, clockwise from north
dy   = [-1, -1,  0, +1, +1, +1,  0, -1]    world+0x58eb8
dst  = src + ((dy[dir] << 8) + dx[dir])    world+0x58ec0, built from the two above

speed = grpAI+0x44 (u8) if nonzero         the group's slowest member's Speed, set by the
        else actor+0x8c (i16)              last FORMATION group move; else the class Speed

domain == 1 (ground):
   d = clamp((i8)(height[src] - height[dst]), -32, +32)     downhill is d > 0
   v = SpeedMultiplier * speed                              map.reg [Path Finding], ships 8
   v = v + ((v * d) >> 6)                                   arithmetic shift; uphill reduces
   c = ((u8)(cost[src] + cost[dst])) >> 1 ; if c == 0 -> 8  byte-wide add, wraps at 256
   v = v / c                                                signed
domain != 1 (the other movement domains):
   v = speed                                                no multiplier, no slope, no cost
v = clamp(v, 1, 63)

stepX,stepY = (dx*dy == 0) ? (v*dx, v*dy)                   straight: an 8-bit IMUL
                           : (trunc(v*dx*K), trunc(v*dy*K)) diagonal: K = 0.707 at 0x59cd98
ticks       = ceil(256 / (stepX != 0 ? |stepX| : |stepY|))  -> mover+0xaa
```

Per tick, the position advances by `(stepX, stepY)` in 1/256ths of a cell per axis and `mover+0xac`
counts up; on reaching `mover+0xaa` the fractions snap to the centre and the surplus is discarded.

### When the group term is set, and when it is not (`MOVE-GATE-035`…`037`)

`grpAI+0x44` is written by exactly two routines — the group Move setter `FUN_005340a0` (group order
4) and the Swarm-2 setter `FUN_00534390` (order 5) — and by each of them **only when the order is
issued in formation**. One local flag decides it, and the same flag decides where each member is
sent:

```
flag = 1
mode = [[grp+0x44] + 0x30] + 0x1f          the owning Player's formation mode, default 2
if mode != 2:                              0 = never, anything else = always
    flag = mode
else:
    centre = mean over members of (fineX, fineY) >> 8       FUN_00533210, unsigned divide
    for each member:
        if max(|mx - cx|, |my - cy|) > AImanager+0xa824:    Chebyshev, whole cells
            flag = 0                                        threshold = 2, a code constant

for each member:
    if flag:  order to (X + (mx - cx), Y + (my - cy))   offsets kept at ord+0x24 / ord+0x26
              min = min(min, member Speed)               seeded 0xfa
    else:     order to (X, Y)                            every member the same cell
if flag:  grpAI+0x44 = min
grpAI+0x20 = 4 or 5 ; grpAI+0x0a = (Y << 8) | X
```

**Nothing clears it.** The only instruction in the image that writes `grpAI+0x44 = 0` is inside
`FUN_005355c0`, and that routine is unreachable — no call, no reference, no immediate, and no
occurrence of its address as a dword anywhere in the image (`AI-DEAD-036`). What resets the term in
practice is allocation: every **player** order builds a new group, whose record the constructor
zeroes. A **scenario-authored** group persists, so its rate outlives the order that set it and
survives every later group command, a member dying, a member being ordered away, and a save/load —
the record is (de)serialized raw, `0x50` bytes, by `FUN_005391d0`.

A consumer that wants one rule: **the group term is live for a group iff the last Move/Swarm-2
order that group received was issued in formation, and it then stays live until another such order
replaces it.**

**Turning is a separate rate and a separate cost.** A step only starts when `mover+0x00` already
equals the desired facing. A turn of less than 33 units snaps and costs 1 tick; a larger one costs
`ceil(min(arc, 256-arc) / mover+0x0a)` ticks and **destroys the dynamic route** first. `mover+0x0a`
is the `RotationSpeed` column (ctor default `0x10`).

**The clock.** One `FUN_00548c60` per actor per **sub-tick** — the counter `server+0x04`, paced by
`FUN_004753c0` against `timeGetTime` at `campaign+0x3f0 = 1000/R` ms, `R` from the nine-arm ladder
`{8,10,12,14,16,20,24,28,32}` defaulting to index 4 (`SESS-CLOCK-005`). Nothing in the step routine
reads elapsed time: the rate is per tick, so cells/tick is invariant and cells/second moves with the
game-speed setting. The same loop iteration issues the `0x401` presentation tick that drives
animation, immediately after the simulation tick; `campaign+0x3dc & 1` clear stops both
(`SESS-PACE-018`).

Worked example, the shipped defaults: `Speed` 16, `SpeedMultiplier` 8, both cells cost 8, level
ground → `v = 16`, step 16, `ticks = 16` — exactly one full tick, ≈ 992 ms, per cell.

**Customisation limits (G2).** `v` is clamped to `[1,63]`; the transit is a whole number of ticks, so
`v ∈ [1,63]` yields only 27 distinct transit times and the shipped 22 `Speed` values yield 15
distinct times on cost-8 terrain; the slope term is a `>>6`; a diagonal is 0.97..1.15× of `√2 ×` the
straight time rather than exactly `√2`. Only `Speed`, `RotationSpeed` (`data.bin`),
`SpeedMultiplier` and the `Cost*` alphabet (`map.reg`) are carried by a shipped file
(`MOVE-LIMIT-033`). The two gates above are **neither**: the spread threshold is a compile-time `2`
and the formation mode is per-player runtime state whose only authored surface is trigger instant 7
(`Set formation`), which ships once in the whole corpus, carrying the constructor's own default
(`AI-SPREAD-038`, `AI-FORM-037`).

## What an area effect can do to movement (`MOVE-AREA-038`)

Every plane read on the movement path is `TEST byte ptr [cell + plane],reg` with `mover+0x5`
reloaded per cell, in `FUN_0054bd20` and `FUN_00543060` on the static plane, `FUN_0054bf10`,
`FUN_0054c100`, `FUN_0054c2f0`, `FUN_0054c6f0` and `FUN_00541dd0` on the dynamic plane, plus the
copies inlined in the driver. No movement routine reads an area-effect layer slot, a cell-record
occupant slot, or a plane bit by immediate, and the only mask values that exist are `0x41`, `0x44`
and `0x82`. An area effect therefore reaches the search through exactly two writes, both made by
`FUN_005456d0` when it recomputes a cell from its record:

- the **Wall of Earth** layer slot `payload+0x20` sets bits 0 and 2 on both planes, so masks `0x41`
  and `0x44` are blocked and `0x82` is not — movement domains 1 and 2 stop, domain 3 crosses;
- **any** occupied layer slot multiplies the cell's cost byte by 4 (an 8-bit `SHL`), which the
  `movementType == 1` cost arm reads inline at the destination cell. Two layers on a cost-16 cell
  truncate the byte to 0.

The step-duration routine `FUN_0054d210` does not read the cost plane directly: it goes through
`FUN_0054e5e0`, which for a cell with any layer returns `cost >> 2` and **stores that value back**
into the plane.

`FUN_00541dd0` is a standalone `n x n` footprint query on the dynamic plane, instruction-for-
instruction `FUN_00543060` apart from the displacement. It is not part of the search: 5 call sites
in 2 owners, `FUN_0052e6d0` and `FUN_0054e220`.

## Refresh policy

| trigger | effect |
|---|---|
| a cell transit completes | the **whole dynamic route** is freed (`FUN_005495f0`) |
| `mover+0x78` (per tick) > `DynamicRefreshRate` | dynamic re-search; counter reset |
| ordered target != `mover+0x74` | static re-search; the dynamic route is freed |
| `mover+0x09` (dynamic searches since the last static one) > `StaticRefreshRate` | static re-search |

The dynamic search does not aim at the final goal: it takes the static route's tail waypoint when it
is more than `DynamicByStaticLookup` cells away, else the node `DynamicByStaticLookup + 1` further
along, and the final goal only when the static route holds `StaticIsntNeeded` nodes or fewer. A
waypoint is popped once the unit is within `DynamicByStaticLookup` of it. Every counter is per-unit
and reset on use — **nothing is staggered**, and neither terrain change nor target movement
invalidates a stored route by itself.

## What stops a mover before the search runs (`MOVE-GATE-039`, `MOVE-STEP-040`)

Movement is not gated inside the search or the step. It is gated by the per-actor order machine,
one switch above the walk order.

The routine that writes an actor's position is `FUN_00548c60`. Closing the call graph above it
(`EnumRefs callto:`, 0 orphan at every step):

```
FUN_00548c60  <- FUN_005495f0, FUN_00549990
FUN_005495f0  <- FUN_005310e0 @005311f6, FUN_00548f70, FUN_005492a0
FUN_00549990  <- FUN_00548f70, FUN_005492a0
FUN_00548f70  <- FUN_005310e0 @0053129c, FUN_00531970
FUN_005492a0  <- FUN_005310e0 @005313c0, FUN_00532100
FUN_00531970  <- FUN_005310e0 only        FUN_00532100 <- FUN_005310e0 only
FUN_005310e0  <- FUN_004f37be (actor vtable slot 6), FUN_00531070 (unreachable)
```

`FUN_00531070` is a loop that calls the machine over a list; `callto:` returns 0 hits for it, and a
raw scan of every section for the stored dword `0x00531070` also returns 0, so nothing reaches it.
So **an actor is displaced only from inside one call of `FUN_005310e0`**, and any state of that
machine which does not reach the walk arm or progress arm 3 leaves the position untouched. The
actor is not slowed and does not drift.

The states that do this:

| state | set by | effect |
|---|---|---|
| `ord+0x09 = 4` | `0053116e`, when `actor+0x144 & 0x100000` (spell 20) and the byte is 0 | no order arm runs until the bit clears (`MAGIC-ACTGATE-079`) |
| `ord+0x09 = 0xff` | `00531789`, when the queued command is `actor+0x50 == 0x17` | no order runs; cleared only by `FUN_00532e60` |
| `actor+0x54 == 0x10` | elsewhere | `FUN_004f37be` returns at `004f37e9` before the machine runs |

A refusal takes effect only from `ord+0x09 == 0`, which `AI-ORDER-039` fixes as the tick the actor
stands on a cell centre (`FUN_00545c30` testing `pos+0x4 == pos+0x5 == 0x80`). An actor in transit
is in progress state 3, whose arm still calls `FUN_005495f0` and clears the byte on arrival. **A
stopped unit is therefore always aligned to the grid, never caught between cells.**

The same closure bounds the one way out of a refusal. The machine's tail at `0x0053165f` is not
gated, and it can install an attack (`MAGIC-ACTGATE-079`); it runs only when `mover+0x98` is
non-zero, and that flag's three setters — `00549092`, `005494a4`, `00549b8d` — are all inside
`FUN_00548f70`, `FUN_005492a0` and `FUN_00549a90`, which the closure above puts behind an order arm.
A refused actor can therefore never raise the flag again, and the tail clears it at `00531673`.
