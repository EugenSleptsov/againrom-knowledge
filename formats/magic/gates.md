# Effect action gates

[Reference](format.md)

## Effect action gates

`MAGIC-ACTGATE-079`, `MAGIC-ACTKEY-080`, `MAGIC-MASKREAD-081`, `MAGIC-AIBIT-082`.

### The gate

`actor+0x144` is a bitmask of the spell ids currently attached to the actor (`MAGIC-ATTACH-016`).
`FUN_005310e0`, the per-actor order machine, reads it before it reads anything about the order:

```
0053114e  MOV EAX,dword ptr [ESI + 0x144]
00531154  MOV DL,0x4
00531156  TEST EAX,EAX
00531158  JZ  0x00531171
0053115a  TEST EAX,0x100000
0053115f  JZ  0x00531171
00531161  MOV EAX,dword ptr [ESI + 0x158]
00531167  MOV CL,byte ptr [EAX + 0x9]
0053116a  TEST CL,CL
0053116c  JNZ 0x00531171
0053116e  MOV byte ptr [EAX + 0x9],DL
```

`ord+0x09` is the progress byte of `AI-PROGRESS-034`. The order switch of `AI-ORDER-039` — walk,
attack, cast at an actor, cast at a cell, and eleven others — is entered only while it is 0
(`0053117c JZ 0x0053125a`). Progress value 4 routes to the arm at `0x00531228`, which sets
`actor+0x54 = 0x1a`, re-tests the same bit, calls nothing, and clears the progress byte only when
the bit is gone.

So one bit refuses every kind of action, and it does so above the point where the kind is chosen.
Value 4 is written at `0053116e` and at no other instruction in the image.

### The one escape, and its bound

The refusing arm is not the end of the tick. Both its exits reach the machine's common tail at
`0x0053165f`, which is not gated on the mask: when `mover+0x98` is non-zero it clears that flag and,
for any command state but 1, `0xa` and `0x17`, runs target acquisition `FUN_005327d0`. If that finds
a target the reach test accepts, it calls `FUN_00531b10`, whose whole body is

```
ord+0x09 = 1 ; ord+0x15 = 0 ; actor+0x54 = 3 ; actor+0x5c = ord+0x0c
```

— order kind 2's own attack install, replacing the parked 4. The gate does not re-park while the
byte reads 1, so progress arm 1 runs, holding `actor+0x54 = 3` and counting `ord+0x15` until it
exceeds 2 with `actor+0x136` set, then returning the byte to 0; the gate parks the actor again on
the following tick.

The escape cannot repeat inside one refusal. `mover+0x98` has three setters — `00549092` in
`FUN_00548f70`, `005494a4` in `FUN_005492a0`, `00549b8d` in `FUN_00549a90` — and all three are inside
movement executors that are reachable only from an order arm (`MOVE-GATE-039`). The refusing arm
calls none of them, and the tail clears the flag itself. **So a consumer must allow at most one
attack to start at the beginning of a refusal, and none afterwards.**

### What reaches the gate, and what does not

The bit index is the spell id. `FUN_005014ae` sets `1 << effect+0x0c`, and the per-spell arm stamps
`effect+0x0c` from `spell+0x8`. The gate's immediate `0x00100000` is `1 << 20`, and entry 20 of the
image's spell-name array at `0x005c5b58` is `stone_curse`.

Nothing about the effect record reaches the gate: not the kind parsed from the `Effects` column
(`effect+0x3c`), not the mode (`+0x3d`), not the magnitude (`+0x40`), not the duration (`+0x42`).
Spell 20's own `Effects` column parses to `absorbtion=+5`, which is applied through the effect-kind
dispatch of § 5 and is unrelated to the refusal. **An implementation built from the column alone
produces damage absorption and no immobilisation.**

### The observable

| Quantity | While the bit is set |
|---|---|
| position | unchanged; every call site of the routine that displaces an actor lies inside one call of the order machine, and the refusing arm calls none of them (`MOVE-GATE-039`) |
| a step in flight | completes; the gate fires only from progress 0, which is the tick the actor stands on a cell centre (`MOVE-STEP-040`) |
| orders | none of the fifteen arms runs; the machine's tail still runs and may install an attack once, see above |
| act state `actor+0x54` | forced to `0x1a` every tick, which is outside the actor tick's own fifteen-arm switch, so no state arm runs either (`ANIM-PARK-039`) |
| the drawable | drawn grey, with the frame index replaced by the facing so it no longer follows the animation clock; keyed separately on the effect record's `+0x0e` (`MAGIC-STONEDRAW-084`, § 15) |
| end | when `FUN_0050134f`'s per-tick countdown of `effect+0x42` reaches zero and clears the bit |

Duration before resistance is `ftol(1.025^power × SpellDuration × 16)` ticks — 160 at power 0, 262
at 20, 549 at 50, 1890 at 100 on the shipped row — then multiplied by `(100 − target+0xca)/100`
with a floor of one tick (`MAGIC-SING-019` d; a different clause of this id, the item-cast
universal reach, was retracted).

Two further, spell-independent ways an actor stops acting, which a consumer must not confuse with
this one: `ord+0x09 = 0xff`, written when the queued command is `actor+0x50 == 0x17` and cleared
only by `FUN_00532e60`; and `actor+0x54 == 0x10`, which makes the actor tick return at `004f37e9`
before the order machine runs at all.

### The AI side

The same bit is read by the AI without gating anything. `FUN_00535d30` adds 127 to a melee target's
cost when it carries the bit. `FUN_0053e510` will not cast spell 20 on a target that already carries
it, and `FUN_0053e840` generalises that test to any spell with `1 << spell+0x8`.

### Customisation limits

The immediate `0x00100000`, the parked value 4, the act state `0x1a`, the six progress arms and the
255-byte index table are `.text` code and immediates. Changing which spell immobilises, or making a
second spell do it, changes `ROM.EXE`. The duration is data: the row's `Spell Duration` column and
the target's `protectionEarth`. The bitmask is one dword, so the spell id space is bounded at 32 by
the field width regardless of how many rows `Data.bin` carries.
