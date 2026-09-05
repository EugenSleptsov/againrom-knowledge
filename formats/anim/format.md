# ANIM — what drives a unit's animation — specification (partial)

Level 3. Promoted, evidence-backed claims only. The clock, action block, phase rules,
run lengths, message set, death chain and object boundary are `ANIM-CLOCK-001`…`ANIM-OBJ-008`;
the sheet/frame boundary is `SPR256-UNIT-024`, and simulation-side timing is `HERO-CADENCE-023`.
Ledger: `claims/anim.md`.

**Status: partial (◐).** The driver, its clock, its five-field state block, all eight action arms,
the eleven message opcodes that fill it and the whole death chain are read at instruction level.
Not specified: `CAirUnit`'s and `CProjectile`'s overrides; what distinguishes opcodes
`0x86`/`0x8a`/`0x8b`/`0x8c` from `0x6b`/`0x71`; `FUN_0040eaee`, the every-32-ticks call on the same
chain; and no runtime session has confirmed the clock by eye.

This is not a file format. It is the **boundary** between the simulation and what is drawn: the
sheet is `formats/spr256`, the frame switch and the blit are `formats/terrain`, the simulation-side
durations are `formats/hero` and `formats/move`.

## The one thing a consumer must not get wrong

**The animation frame is presentation state, and it is stepped by the simulation's clock.**

- **Not simulation state.** No simulation actor holds a frame, a phase or an animation state. It
  holds `actor+0x138` — the tick its current run is due to end — and nothing else. Every animated
  quantity is on the client drawable `CUnit`, which is not one of the 28 serializable classes, so
  none of it is saved and none of it needs to be reproduced bit for bit.
- **Not a free-running clock either.** The driver is called from the paced `0x401` game tick, at
  `1000/tps` ms with `tps ∈ {8,10,12,14,16,20,24,28,32}` selected by the game-speed index (default
  16 tps, 62 ms — `1000/tps` is a **truncating** divide, so the realised periods are
  `125 100 83 71 62 50 41 35 31` and only three of the nine are exact). Animation slows with game
  speed and stops when the tick stops.
- **And the walk is the one exception, deliberately.** Every timeline in this engine advances one
  step per tick — idle, attack, shoot, cast, death, and a placed object's — **except the walk**,
  which advances one step per **1/16 of a cell travelled**. So the walk runs at `16 / ticksPerCell`
  steps per tick against everyone else's 1, and the two coincide only when a cell takes exactly 16
  ticks: **78 of 262 shipped (class, definition) pairs**, the rest between 0.5× and 2.0×
  (`ANIM-WALK-013`, `ANIM-AMBIENT-016`). This is not a defect to be corrected in a reimplementation;
  it is what the engine does, and a consumer that ticks the walk like everything else will be right
  for 78 units and visibly wrong for 184.

## The state block on the drawable

```
+0x6c  facing, 16-way                       +0x94  the running action's clock
+0x70  phase, consumed by the frame switch  +0xa0  ticks the action still has to run
+0x74  draw state = a copy of +0x84         +0xbc  turn accumulator, 1/16 of a facing step
+0x84  action code, one byte                +0x15a corpse stage (the server's, shipped)
+0x85  the facing the action ends on        +0x88/+0x8c  remaining position delta, 1/256 cell
```

## The driver: once per game tick, per drawable

```
FUN_0040ec30 (map view message handler), arm msg == 0x401
  -> FUN_0040dcdd
       CMapView+0xa70 += 1                        the object / water animation counter
       for each drawable in the view's container:
            vt+0x3c()                             = FUN_0045cf00 for a CUnit
            returns false -> drop it from the walk

FUN_0045cf00:
  ticksRemaining (+0xa0) == 0  ->  state := 0
        IdlePhases != 0 : clock++ ; phase := clock % len(IdleAnimTime expansion)
        IdlePhases == 0 : phase := 0 ; clock := 0     (the stand / corpse fork draws)
  otherwise                     ->  switch (actionCode - 1), 8 arms
        ... the arm updates phase / position / facing ...
        state := actionCode ; ticksRemaining -= 1
```

## The eight action arms

| code | run | what advances the phase | run length `+0xa0` |
|---|---|---|---|
| 1 | move | `clock += \|step\|`, `phase = clock/16`; the draw takes it **mod** `len(MoveAnimTime expansion)` | the message's duration byte (the mover's per-step ticks) |
| 2 | — | nothing; the arm only assigns the state | — |
| 3 | attack | `phase = clock`, `clock++`; **no modulo** | `len(AttackAnimTime expansion)` |
| 4 | — | nothing | — |
| 5 | turn | facing interpolates by the shortest arc; the draw takes the **stand** arm | the message's duration byte |
| 6 | die | `phase = clock`, `clock++`; the draw halves it | `2 × classes[Dying].DyingPhases` |
| 7 | shoot | as attack, plus the projectile spawn on the tick `clock == ShootDelay` | `len(AttackAnimTime expansion)` |
| 8 | cast | as attack, over `ShootDelay` / `AttackDelay` | `len(AttackAnimTime expansion)` |

Codes **2 and 4 are set by no immediate anywhere in the image**, which is why the frame switch's
arms 2 and 4 fall to a default that draws the class id.

**A run's length is taken from the art, never from the duration the simulation ships.** The server
does put `attackChargeTime + attackRelaxTime` in the attack message and the client never reads it.
Over 262 shipped (class, definition) pairs the animation and the pause are equal on 119 and differ
by −36 … +26 ticks on the rest; the death run exceeds `dyingTime` on 148 of 262.

**A walk cycle is advanced by distance, not by time.** The clock `+0x94` accumulates the
**Euclidean** length of each tick's own displacement — `FSQRT` of `sx² + sy²`, truncated per tick —
so it is an odometer in 1/256 of a cell, and `phase = clock >> 4` is one timeline step per **1/16
of a cell**. No speed term appears anywhere in the animation code: a slow unit takes more ticks to
cover the same distance and holds each frame longer; it never skips one.

- **A straight cell is exactly 16 timeline steps, for every unit at every speed** — 262/262 pairs.
  A **diagonal** is `256√2 = 362` in that metric, so 20…22 steps after the per-tick truncation.
- **How many *cycles* those 16 steps are is data**: `class+0x50 = len(MoveAnimTime expansion)`, and
  `16/L` is the cycles per cell. The shipped convention is `MovePhases = 8` frames held two steps
  each, `L = 16`, **one cycle per cell** — 208 of 262 pairs. `Bee` and the two `Catapult`s run four
  cycles per cell, `Ghost` two, and 44 pairs have `16 mod L != 0` and are not cell-aligned at all.
- **The clock is not reset when a walk starts**, so consecutive cells continue one odometer. It
  *is* zeroed by any turn, attack, shoot or death message, and by standing still for a single tick
  when `IdlePhases == 0` — which is 29 of the 34 classes (`ANIM-WALK-013`…`015`).

## The idle run — state 0, and the fidget

State 0 is not "no animation": it is a **fork on the class's `IdlePhases`**, and a consumer that
draws a standing frame for every idle unit is wrong for five of the 34 shipped classes
(`ANIM-IDLE-009`, `ANIM-IDLE-010`).

```
state 0, IdlePhases != 0     clock++ once per tick; phase = clock % len(IdleAnimTime expansion)
                             frame = idleBase + dir*IdlePhases + idleTL[phase]   (NO modulo)
                             idleBase = S + D*(MoveBegin+Move+Attack+Dying)      = the BONE base
state 0, IdlePhases == 0     phase = 0, clock = 0
                             corpse stage >= 1 -> the dying/bone fork; else the standing frame
```

Three consequences, each of which a consumer gets wrong by default:

- **The corpse fork is unreachable for a class with idle art.** The draw reads `IdlePhases` before
  it reads the corpse stage, and the idle and bone blocks share one base — so a class either
  fidgets or leaves bones, never both, and the sheet layout is consistent because it has to be.
- **The idle clock is the action clock.** `+0x94` is not reset when an action ends, only when one
  starts, so a unit resumes its idle loop at `(whatever the last action left) % len` — two units of
  one class idle out of phase. Reproducing this needs no extra field; forgetting it makes a group of
  monsters flap in lockstep, which the engine never does.
- **The loop length is the expansion's length, not `IdlePhases`.** `Ghost` carries 3 frames and a
  4-step timeline `0 1 2 1` — a ping-pong. `IdlePhases` sizes the *block*; `class+0x80` sizes the
  *loop*.

Shipped data, all five classes, one step = one animation tick (62.5 ms at the default speed index):

| class | `IdlePhases` | timeline (expanded) | loop | cycle @16 tps | placements / 8094 |
|---|---:|---|---:|---:|---:|
| Bee (73) | 4 | `0 1 2 3` | 4 | 250 ms | 608 |
| Ghost (69) | 3 | `0 0 0 1 1 1 2 2 2 1 1 1` | 12 | 750 ms | 416 |
| Sonic Bat (70) | 6 | `0 0 1 1 2 2 3 3 4 4 5 5` | 12 | 750 ms | 580 |
| Dragon (71) | 7 | `0 0 1 1 … 6 6` | 14 | 875 ms | 135 |
| Death Star (72) | 7 | `0 0 1 1 … 6 6` | 14 | 875 ms | 1 |

The other 29 classes carry no `IdlePhases`, and their whole idle presentation is the standing
block: one frame per facing, changing only when the facing does.

## Directions

One 16-way facing field. The standing block indexes it whole, `(facing − 8) & 0xf`; every other
block halves it, `((facing − 8) >> 1) & 7`; `Flip` mirrors the upper half (`SPR256-UNIT-024`,
`REG-UNITS-051`). There is no separate turning run: while the turn action is active the unit is
drawn standing at a facing that advances by `(target − current)/ticksRemaining` sixteenths per tick,
wrapping by the shortest arc.

## The messages that fill the block

`FUN_004104e8` dispatches on `msg+0x9` biased by 3, through a byte arm-index table at `0x0041879b`
into a jump table at `0x004186e3`.

```
0x6b  action 1 move    facing, duration, and the per-direction delta from 0x005c23e8 / 0x005c2428
0x6d  action 5 turn    facing, duration
0x71  action 3 attack  facing            (refused while an action is running, or AttackPhases == 0)
0x72  action 7 shoot   facing, and the target's runtime id at +0x86
0x86  action 1 or 8 · 0x8a action 8 · 0x8b action 1 · 0x8c action 1
0x6c 0x6e 0x6f 0x70    the field-masked state sync — one bit per field:
      0 health  1 mana  2 +0x108  3 corpse stage  4 facing  5 position
```

## The death chain, end to end

```
server   actor+0x13c = 1 at death; FUN_004f52ee raises it to 2 / 3 / 4 as health passes
         -10 / -20 / -40, and to 5 below -600; every change is broadcast
client   the stage arrives in the state sync under mask bit 3; the client keeps the previous
         value and dispatches on the TRANSITION:
             0 -> 1   action 6, 2 x DyingPhases ticks, clock 0        the fall
             1 -> 1   action 6, 4 ticks, clock 2*DyingPhases - 4      the last two frames again
         when the run ends, no action runs and the frame switch's state-0 corpse fork draws
         on the stage alone:
             1        dying frame DyingPhases - 1, frozen
             2 3 4    bone frames 0 1 2
             5        the actor is gone
```

The thresholds never reach the client as numbers — it sees only the stage. The four
`movementType > 1` classes (`Ghost`, `Bee`, `Bat_Sonic`, `Dragon`) take `health = -1000` the moment
`dyingTime` expires, so they pass −10 without pausing, never occupy stages 2..4, and leave neither
corpse nor bones.

Corpus caution: `BonePhases ≥ 3` on 261 of 262 pairs. `Death Star` (ID 72, `Daemon`) carries the
loader's absent `-1`, which the corpse arm's `== 0` guard does not catch, so its bone index is
`BaseBone − dir + stage − 2` — inside the dying block, walking backwards with the facing.

## Objects

Same clock, different mechanism, no second driver. The routine that ticks every drawable also
increments `CMapView+0xa70`, and a placed object's frame is a pure function of that counter and its
cell — `Index + T[(animCtr + col·(row+1)) mod class+0x3c]`, recomputed at every draw, with no
per-instance state at all (`TERR-SPR-042`, `REG-OBJ-046`). The animated arm's gate is the four-corner `0xc000` test that 0 of
880 704 cells of a shipped **map file** satisfy (`TERR-TILE-044`) — but those two bits are
**runtime** state: `vt+0x48` ORs `0xc0` into the tile word's high byte at exactly those four corners
around every drawable, and the 32-tick sweep `FUN_0040eaee` clears bit 14 map-wide (`ANIM-TICK-011`).
Those two bits are the **fog of war**, and the gate holds on exactly the cells the local player can
currently see, so the arm **does** fire: a shipped map's fires and trees animate inside the field of
view and hold frame 0 outside it (`TERR-TILE-079`, `ANIM-OBJ-008`).

**Cadence, which is the thing the walk does not share.** `CMapView+0xa70` advances an object's
timeline **one step per tick**, unshifted. The ambient *repaint* is throttled by a second field:
`CMapView+0xa74` latches the counter, and a redraw is admitted only once it has advanced by more
than 3 — **one ambient frame per 4 ticks**, i.e. 248 ms at the default speed. Neither field is the
walk clock and neither is the mover's tick count: they are a third counter on the same pacer, and
no game-speed setting changes the ratio between them (`ANIM-AMBIENT-016`, `ANIM-PACE-017`).

## What a consumer needs

One tick source at the game-speed rate; per unit, the eight-field block above; the eight arms
**and the state-0 idle fork**; the frame switch of `formats/terrain`; the sheet layout of
`formats/spr256`. For objects, one
counter and a pure function of `(animCtr, cell)`. Nothing here belongs in a hash of simulation
state, and nothing here belongs in a save.

**And two cadences, not one.** Everything above runs at one timeline step per tick; the walk runs
at one step per 1/16 cell. A consumer that implements a single "animation frame counter" produces a
walk that is in step with the world for units of `Speed` 16–17 and out of step for every other, and
no amount of tuning the tick rate fixes it — the ratio is the same at all nine game speeds
(`ANIM-WALK-013`, `ANIM-AMBIENT-016`, `ANIM-PACE-017`).

## A projectile — a different animation model entirely

`ANIM-PROJ-025`, `ANIM-PROJ-026`. `CProjectile` overrides body draw, shadow and driver,
so nothing above this heading applies to it. The engine's own field names come from its savegame
section `[Prj%d]`:

```
+0x08 x        +0x6c dir          +0x84 action (b)       +0x88 actionx
+0x0c y        +0x70 phase        +0x85 actiondir (b)    +0x8c actiony
+0x10 z        +0x74 lastaction   +0x86 actiontarget (w) +0x90 actionz
+0x20 picture                                            +0x94 actionphase
                                                         +0xa0 actionsegments
                                                         +0xa4 actionspell
object size 0x14c; the save also carries [Projectiles] Count / FreeIndex / IDs
```

`picture` is a `projectiles.reg` `ID`. The driver runs once per tick:

```
if actionsegments == 0            -> the object is finished
if actiontarget != 0              -> actionx/y/z := the target's CURRENT position
step                              := (actionx - x) / actionsegments, per axis
actiondir                         := recomputed from the two coordinates
actionphase                       += 1
lastaction := action ; actionsegments -= 1
```

So **the life is a countdown the spawner sets, not a property of the art**, and **tracking is a
property of having a target, not of the `Homing` key** — which nothing in the driver or the draw
reads. The behaviour is a `switch` on `picture` (biased by 13, spanning 13..64): the default arm
travels with `phase = (actionphase / 2) % Phases`; eleven ids snap to the target and notify it; two
more do that into the target's own slot array; two take `phase` from a fixed 13-step ramp; one plays
once with `phase = actionphase`. Caster-attached, target-attached, travelling and area are therefore
**one mechanism**.

The draw is a second `switch` on the same field with a different bias (7) and span (7..60):

```
picture out of range, or its slot empty  -> nothing is drawn
x -= Width/2 ; y -= Height/2 - z - <view offset>
facing = (dir - 8) & 15 ; if Flip and facing > 8 { facing = 16 - facing ; mirror }
frame  = Phases * facing + phase          (frame = phase when RotationPhases == 1)
Palette == 0 -> the shared projectiles.pal, otherwise the sheet's own
the sheet is loaded on FIRST DRAW, not at start-up
picture 10 and 12 also blit a smoke sheet once per point of the object's trail array
```

## The projectile frame clock

`ANIM-PHASECLOCK-028`. A projectile's driver computes its sheet frame as

```
phase = (|actionphase| / 2) % Phases          Phases from projectiles.reg
phase = 0                                     when the picture has no registry row
```

so **effect art advances one sheet frame every two game ticks**, against one frame per tick for a
unit action. Four picture ids replace the rule:

| picture | rule |
|---|---|
| 34, 36 | a 13-entry constant ramp indexed by `actionphase - 1` (below) |
| 51 | `phase = actionphase`, raw |
| 60 | `phase = actionphase - 1`, no modulus |

`ANIM-BOLTRAMP-035` reads the 13-entry table out of the image. In order, for
`actionphase` 1 to 13, it yields

```
4, 3, 2, 1, 0, 1, 2, 1, 0, 1, 2, 3, 4
```

which is one value per tick and consumes exactly the 13-tick life those two pictures get. The range
0 to 4 is exactly `lightnin`'s 5 phases. Picture 36 adds `5 * (link index mod 7)` to the same ramp
(below), giving 0 to 34 against `chain`'s 35 phases.

Corroboration from shipped data: the only two burst lifetimes that are not the default 16 ticks are
18 and 22, against `acid`'s 9 phases and `fireexpl`'s 11 — `2 * Phases` in both cases, one full pass
under this clock and under no other divisor.

## The polyline draw of pictures 34 and 36

`ANIM-BOLTDRAW-034`. The two picture ids that draw a path do not draw the projectile's
own sprite at all. Their draw arms iterate a list of 8-byte records held on the object (the list
itself, and the geometry that fills it, are `formats/magic` section 14):

```
for i in 0 .. count-1:            count = object+0x118, records at object+0x114, stride 8
    x = i16 record[0] - 8
    y = i16 record[2] - 8
    frame = phase                             picture 34
    frame = phase + 5 * u8 record[6]          picture 36
    blit(sheet, frame, x, y)
```

Three properties a consumer must reproduce:

- **The sheet is a constant per arm**, `projectiles.reg` record 34 and record 36, not the
  projectile's own picture field. Re-pointing either spell's art at another row moves nothing.
- **No facing is folded.** Neither arm reads the direction, and both shipped records carry
  `RotationPhases = 1`, which already disables the `Phases * facing` term. A bolt's apparent
  bend is geometry, never rotation.
- **The 8-pixel centring is an immediate**, not the record's `Width`/`Height` halves the routine
  computed for the ordinary path. A replacement sheet of any size other than 16 by 16 draws
  off-centre.

Every point of one figure carries the same frame for picture 34, so the whole bolt flickers as one.
For picture 36 the per-record byte is the chain-link index modulo 7, so each branch is drawn from a
different fifth of the 35-frame sheet — the byte is a branch identifier, not an age.

## The act state a refused actor is parked in (`ANIM-PARK-039`)

`FUN_005310e0`, the per-actor order machine, runs immediately before the actor tick's own
act-state switch, and it can leave the actor in a state that switch has no arm for.

```
004f3994  MOV EAX,dword ptr [EDX + 0x54]
004f399d  SUB ECX,0x1
004f39a3  CMP dword ptr [EBP + -0x78],0xe
004f39a7  JA  0x004f41b3          the epilogue, not a fifteenth arm
```

The machine writes `actor+0x54 = 0x1a` from four sites: its early-out at `0053110e`, the refusing
progress arm at `0053122e`, the `ord+0x09 == 0xff` arm at `0053124e`, and the default progress arm.
`0x1a − 1 = 0x19` exceeds the switch's bound of `0xe`, so none of the fifteen arms runs and the tick
returns. The value means "the machine had nothing to run" and carries no information about why, so
it cannot be used to identify a spell, a death or a halt.

What the client draws for such an actor is decided on the drawable, not here. For `stone_curse` the
unit draw replaces the frame index with the unit's facing, so the frame no longer follows the
animation clock, and separately replaces the shade table with a greyscale one; both are keyed on the
effect record's `+0x0e` rather than on `actor+0x144` (`MAGIC-STONEDRAW-084`). The two are separate
mechanisms in the two hierarchies this ledger's header keeps apart: the draw substitution is
presentation and is not serialized, the action refusal is simulation and is.

## Reviewed spell-effect presentation (`ANIM-044`…`ANIM-047`)

The paced `0x401` tick drives both actor marks and projectiles. CUnit and CAirUnit invoke vtable
`+0x50` once before their action branch, and both bind it to the effect-list rebuild. A mark is
centred at

```
x = unit[+0x60] - Width/2 + dx
y = unit[+0x64] - unit[+0x68] - unit[+0x10] - Height/2 + dy - depth
```

Positive-depth marks draw before their actor body and non-positive marks after it, preserving array
order inside each pass (`ANIM-044`). Shield preserves the midpoint source order and appends every
Component A record before Component B; duplicate alpha stamps therefore remain visible
(`ANIM-045`).

Picture 51, Meteor, starts `actionphase=-1`. The driver increments before the draw, giving sixteen
phases:

```
phase 0..7   frame 8        y offset = 4*phase - 28
phase 8..15  frame phase-8  y offset = 0
```

It ends on frame 7 and never returns to frame 8. The world position stays at the accepted cell; only
the draw offset moves (`ANIM-046`).

Projectiles do not register in the cell grid. Their exact relevant map composition is:

```
selector-1 retained overlay: wall_of_fire, wall_of_earth
all unit shadows
separate projectile list
selector-0 retained overlay: freezing_cloud, poison_cloud
all unit bodies, including their actor marks
later: shroud
```

The exact traversal order among projectiles in that separate list remains Unknown, so overlapping
same-cell Meteors are specified only up to their list order (`ANIM-047`).
