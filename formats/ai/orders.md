# Orders, formation and execution

[Reference](format.md)

## Player orders

Nineteen opcodes, `0x14..0x26`, dispatched by `FUN_004d5dd8`'s order space (19 direct dwords at
`0x4d86ae`; [SESSION](../session/format.md) owns the command itself). **Four are empty** — `0x15`, `0x20`, `0x22`,
`0x23` — and every live arm first tests the AI manager `[0x005f21c4]` and does nothing when it is
null.

**Every order builds a new group.** Before the switch, the prologue resolves the commanded actors,
allocates a fresh `0x48`-byte group onto the player's list and fills it with exactly them. The
group's AI block is constructed with **`grpAI+0x20 = 0`**, so a commanded group always *starts* at
group order 0 — which is the gate the per-actor state machine runs behind. Nothing about the load-time
1-or-3 stance survives a player order.

The orders then split by what they leave that byte at:

| opcode | dispatches to | leaves `grpAI+0x20` | per-actor `actor+0x50` |
|---|---|---:|---|
| `0x14` | `FUN_00534e20` | 0 | `0x16` |
| `0x16`, `0x1c` | `FUN_005340a0(col,row)` | **4** — `Move` | from the argument |
| `0x17` | `FUN_00534ac0` — **guard** | **1** | `0xb` |
| `0x18` | `FUN_00534cb0` — **aggressive** | **3** | `0xc` |
| `0x19` | `FUN_005348c0(target)` | 0 | `3`, then `0xc` |
| `0x1a` | `FUN_00534390(col,row)` | **5** — `Swarm 2` | from the argument |
| `0x1b` | `FUN_00534f80(target)` | 0 | `0xc`, `8` |
| `0x1d` | `FUN_00535300(col,row)` — **patrol** | 0 | `0xa` |
| `0x1e`, `0x25` | `FUN_00533d00(target, spell)` | 0 | `0xd`, then `0xc` |
| `0x1f`, `0x26` | `FUN_00533ed0(col,row, spell)` | 0 | `0xe`, then `0xc` |
| `0x21` | `FUN_005308a0(col,row)` — **pick up a sack** | *unwritten* → 0 | `2` |
| `0x24` | `FUN_0052f8b0(building)` | *unwritten* → 0 | `0xf` |

**Four orders act through a group arm** (guard, aggressive, move, `0x1a`); **ten act by writing a
per-actor state** and leaving the group at 0, which is the only value under which
`FUN_00533ae0` evaluates `actor+0x50` at all. That is why patrol and the pick-up are per-actor
behaviours and guard and aggressive are group behaviours — the difference is one byte written at the
end of the order routine.

**Nine of the fourteen live arms have the dispatcher as their only caller in the image** —
`0x14`, `0x19`, `0x1a`, `0x1b`, `0x1d`, `0x1e`, `0x1f`, `0x21`, `0x24` — while the map's script
reaches `FUN_005340a0`, `FUN_00534390` and `FUN_00534ac0` directly and **never enters the
dispatcher**.

Patrol shares `FUN_005301f0` between player and script routes; the
player helper `FUN_00535300` tail-calls it (`AI-PATROL-017`). Pickup's
immediate-form actor+0x50=2 writers that reach an actor both belong to
`FUN_004d5dd8`. Broader computed-value producers are outside that writer set.

The received object-target command (`cmd+4 == 3`) resolves signed-word id `cmd+0x0e`
against the unit collection first, then the structure collection. Opcode `0x19` forwards
that pointer to the shared attack order. The scorer must not return `0xffffff`, and the
target must not be the member itself, before the member retains it at `ord+0x0c`.
The physical branch proceeds through order 5, shared approach and action 3, copying the
pointer to `actor+0x5c`; no separate structure action or per-strike footprint lookup is
present in these bodies. This is a conditional received-command route, not proof that
every Building instance is safely admitted or that ordinary hover produces the command:
the scorer/AI still carry actor-oriented field reads on the supplied object.
Claim: `UNIT-STRUCTORDER-062`.

**Customisation.** The four empty order slots are absent implementations and free. The order
opcode's ceiling `0x26` is a compiled constant. Neither `actor+0x50` nor `actor+0x54` is serialized,
so extending either state space costs nothing in a shipped save — but `ord+0x08` and `ord+0x09`
*are*, inside a raw `0x94`-byte block, so a new field before them breaks every save the original
wrote while a new *value* in them does not.

## Formation

Every `Player` owns a `0x20`-byte settings block, allocated and constructed by the `Player`
constructor (`004fac7d PUSH 0x20` → `FUN_0052ccb0` → `004facb8 MOV [Player+0x30]`). Its **last
byte**, `+0x1f`, is the **formation mode**; the constructor writes the only nonzero default it
has, `2`. The block is serialized raw, `0x20` bytes, by `FUN_005392d0`, so the mode is in every
save.

| value | behaviour |
|---:|---|
| `0` | never in formation; the group centroid is not even computed |
| `2` | in formation **iff** every member is within `AImanager+0xa824` cells (Chebyshev) of the group centroid — the threshold is a compile-time `2` |
| anything else | in formation **unconditionally**; no spread test runs |

One routine writes it, `FUN_00537e90(player, mode)`, and it has two callers, both authored
surfaces:

- **player command `0x46`, sub-code 2** — the player from `cmd+0x05` through `FUN_004d4c6e`, the
  value from `cmd+0x0e` **remapped** `0→0`, `1→2`, `2→1`, anything else `→2`, so the wire can only
  produce `{0,1,2}`;
- **trigger instant id 7**, which the shipped `Description Instants.ini` names **`Set formation`**
  and which writes the parameter **raw** — so a map may author a value the UI cannot.

What the mode decides is both halves of a group Move / Swarm-2 order: whether each member is sent
to `target + (memberCell − centroidCell)` or every member to the same cell, **and** whether the
group's rate override is written at all. [MOVE](../move/format.md) owns that; the gate is the same
local flag in both cases.

`scn:110.alm` authors `Set formation = 2`, equal to the constructor default.

## Order execution

Between the state machine and the actor tick sits `FUN_005310e0`, called once per actor per tick at
`004f398c` — four instructions before the tick dispatches on `actor+0x54`. It preserves
`actor+0x54` only when it is `2` or `0xf` and clears it otherwise, then dispatches on the order
object's **progress** byte `ord+0x09` (255 in range, four live) and, when that is 0, on its
**pending order** byte `ord+0x08` (15 slots, three empty). "When that is 0" means, in practice, the
tick the actor is standing exactly on a cell centre: the state machine re-decides the order at group
rate, the order machine executes it at cell rate.

| `ord+0x08` | what it does | inputs |
|---|---|---|
| `1` | walk to `ord+0x0a` | the cell |
| `2` | attack `ord+0x0c` **with no distance test** | the target |
| `4` | close on the actor `ord+0x18` | stop distance `ord+0x14` |
| `5` | pursue and attack `ord+0x0c` | facing, reach `actor+0x12c`, stop distance `ord+0x14` |
| `6` | as `5`, plus the auto-engage suppression, and turn instead of path | same |
| `7` | the sack pick-up bridge | — |
| `8` | cast `ord+0x30` at the actor `ord+0x28` | facing, `ord+0x14` |
| `9` | cast at the cell `ord+0x3c` | facing, `ord+0x14` |
| `0xa` | turn in place until `mover+0x01 == mover+0x00` | the desired facing |
| `0xb` | idle: re-face to `facing + 0x21 + (190·rand()/0x8000)`, a near-uniform byte, and only when `ord+0x54` is set or on `rand() < 0xcd` | — |
| `0xc` | walk to `ord+0x0a`, then the `actor+0x54 = 2` action | the cell |
| `0xf` | reach `ord+0x0a` within `ord+0x14`, then `actor+0x54 = 0xf` | the cell, `ord+0x14` |
| `3`, `0xd`, `0xe` | empty — the switch's own default | — |

The in-position test that `5`, `6` and `8` share is **not** a centre-to-centre distance: it is the
actor's current facing equal to the 8-way direction to the target, **and** an edge-to-edge distance
that subtracts both token sizes, so two touching actors measure 1. Reach is re-read from
`actor+0x12c` at the test; `ord+0x14` is used only as the mover's stop distance, inside which the
actor turns to face and stands still.

Worked example — the sack pick-up, which is the whole chain in one order:

1. order `0x21` → `FUN_005308a0`: `actor+0x50 = 2`, `ord+0x0a` = the sack's cell, `ord+0x08 = 0`.
2. `FUN_0052ce50` arm 2: not at `ord+0x0a` → `ord+0x08 = 1`, walk. **At it → `ord+0x08 = 7`.**
3. `FUN_005310e0` pending arm 7 → **`actor+0x54 = 2`**.
4. the actor tick's `actor+0x54 == 2` arm loots the sack under the actor's feet.
5. next tick, arm 7 sees `+0x54` already 2 and completes: `actor+0x50 = 0xc`, `ord+0x08 = 0`.

So a unit told to pick something up ends the order **hunting**, not idle. The same relay carries
every other order that has a "walk there, then do a thing" shape.
