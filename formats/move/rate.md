# Movement rate and formation gates

[Reference](format.md)

## Movement rate and clock (`MOVE-RATE-029`…`034`)

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

**Turning advances before stepping.** A next-cell step starts only when current
facing byte+0 already equals desired byte+1. Its mismatch arm calls the turn
routine and returns even if that call reaches the target facing. The full
active DWORD+a0 matters: only zero admits the short-arc (<=32) snap. An active
turn, or a larger fresh turn, advances current by RotationSpeed byte+a along
the shorter arc, clamps at desired and wraps modulo256. An exact 128 tie takes
addition. This leaf has no route/list access; the older route-destruction
interpretation is partially retracted in MOVE-TURN-031.

The turn caller writes BYTE+a4=ceil(pre-step shorter arc/RotationSpeed) after
the advancing call, rather than decrementing a countdown. A fresh short snap
writes 1. It resets BYTE+9d when inactive, increments it on each call, sets the
active DWORD and clears that dword if the updated facings match. The division
arm requires nonzero RotationSpeed. The selected move/order caller chain can repeat
this step; complete order scheduling and route/callback effects remain
Unknown. These are local call boundaries, not a whole-action tick count —
MOVE-TURN-044. Serialized local treatment and the post-LOAD frontier are
SAV-TURNLOAD-822.

**The turning byte has several producers.** Unit initialization stores the
180-byte mover allocation at actor+154. Its constructor writes byte+a=16,
then actor defaults write 8 to the same object. Unit table slot 9 and Human
table slot 7 address that byte; their byte helper preserves the incoming value
on -1 and otherwise stores low8. The former final-default interpretation in
MOVE-TURN-031 is partially retracted — MOVE-RATE-052.

Human derive later writes low8(actor+8c) to mover+a after its speed derivation
and modifier fold. It can therefore replace the independent table value.
The exact Unit derive slot lacks this assignment; the measured Humanoid and
Human slots share it — MOVE-RATE-053. Effect selector 18 first adds into the
target's mover+a modulo256, then calls that same target's derive. This local
effect write alone does not establish a surviving Human bonus — MOVE-RATE-054.

The complete turn leaf gets its actor from the stack argument and reloads
actor+154; incomingECX is not its mover receiver. It changes current facing
only. The selected next-cell arm with actor+184=0 calls it and returns without
a Position change. The turn caller also uses byte+a for its estimate, while
the selected positional-rate body reads actor+8c or the formation override.
All three direct leaf callers are selected; eight of eleven direct turn-caller
sites and all unresolved indirect/rebased accesses remain outside this local
proof. Full scheduling and elapsed time remain Unknown — MOVE-RATE-055.

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
(`Set formation`), whose installed use carries the constructor default
(`AI-SPREAD-038`, `AI-FORM-037`).

## Area effects (`MOVE-AREA-038`)

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

## Pre-search gates (`MOVE-GATE-039`, `MOVE-STEP-040`)

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

The attack tail has a separate gate. The machine's tail at `0x0053165f` is not
gated, and it can install an attack (`MAGIC-ACTGATE-079`); it runs only when `mover+0x98` is
non-zero, and that flag's three setters — `00549092`, `005494a4`, `00549b8d` — are all inside
`FUN_00548f70`, `FUN_005492a0` and `FUN_00549a90`, which are reached through the gated order arms.
A refused actor can therefore never raise the flag again, and the tail clears it at `00531673`.
