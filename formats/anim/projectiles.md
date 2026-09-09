# Projectiles and effect presentation

[Reference](format.md)

## Projectile state

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


## Projectile clock

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


## Picture 34/36 polylines

`ANIM-BOLTDRAW-034`. The two picture ids that draw a path do not draw the projectile's
own sprite at all. Their draw arms iterate a list of 8-byte records held on the object (the list
itself, and the geometry that fills it, are [MAGIC projectile geometry](../magic/projectiles.md)):

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


## Spell-effect presentation (`ANIM-044`…`ANIM-047`, late-pass scope amended)

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

The projectile registration slot is a no-op. The relevant map composition uses
the actual selector/storage join; the earlier universal unit-body wording of
`ANIM-047` is retracted (`ANIM-CATEGORY-084`, `ANIM-AIRPASS-086`):

```
earlier cell phases: alternate CUnit and CBackPack, then ordinary CUnit/static objects
within the main cell phase: area selector 1, wall_of_fire/wall_of_earth
complete registration-selector-3 (CAirUnit) shadow sweep
separate collection walk; projectile type interpretation remains Medium
area selector 0, freezing_cloud/poison_cloud
complete registration-selector-3 (CAirUnit) body sweep
later: marker/bar calls, then shroud
```

The complete projectile insertion/type population and native pixel overlap in
that collection remain Unknown. A dispatch order does not promise a particular
Meteor-over-unit pixel result (`ANIM-047`, amended).
