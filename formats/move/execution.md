# Coordination, ticks and cell transitions

[Reference](format.md)

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
  term — see [group speed gates](rate.md).


## Tick ordering (`MOVE-TICK-013…017`, `MOVE-ID-016`)

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


## Movement step

Position is `actor[4]` = `*(actor+0x10)`: `+0x00` x cell, `+0x01` y cell, `+0x02` packed cell,
`+0x04` x fraction, `+0x05` y fraction, `0x80` = centred. `FUN_00548c60` treats `(cell<<8)|fraction`
as one 16-bit value per axis and adds the signed per-axis step `mover+0xb0` / `mover+0xb1`; when the
cell byte changes it hands `FUN_0054abb0` the **old** packed cell, then writes the new position, then
calls `FUN_00544d00`; when `mover+0xac >= mover+0xaa` it snaps both fractions to `0x80`. A transit
takes `mover+0xaa = ceil(256 / |step|)` ticks. Speed never enters a label — two units of different
speed pick the same route.

### What the cell-boundary calls do (`TERR-CELLREC-146`, `TERR-FOOTPRINT-147`)

`FUN_00544d00(map, actor)` and `FUN_00545230(map, actor, x, y)` are the enter and the leave of the
map's **cell record**, the per-cell structure [TERRAIN](../terrain/format.md) specifies. Both read the
actor's movement domain through `vt+0x20` and pick a slot from it: domain 1 or 2 uses the record's
`payload+0x04`, domain 3 uses `payload+0x08`. The per-cell entry refuses other
domain values; the detach statement here is scoped to domains 1/2/3. `FUN_00544d00` reads the
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
and progress 3 completion can occur despite a refused destination entry. Local
selected serializers subsequently emit that state; no rollback or next-SAVE
reconciliation is established. — SAV-CELLFAIL-583, SAV-CROSSNEXT-585

Entry writes+82/+83 from cached low-byte cell-X/Y getter results
(`00544a10/00544a20`), then+84/+85 from the low bytes of full-X/Y getters
(`005449e0/005449f0`, hence Position+04/+05). Cell-X/Y were sampled before
`0054a620`; full-X/Y are sampled after it. Word+72 is that callback's AX,
not a proved numeric default; the synthetic `0x1357` cut does not establish a state default. Those accessors are
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


## Restored actor registration and the reached turn

LOAD rebuilds Player+20 from group membership, then the global actor list from
those Player lists, skipping actor+4c mask0x08. The restored registration helper
appends the exact actor pointer; the creator helper's ID assignment is a separate
path. The later manager+4 callback invokes actor+24, not actor+50. Unit uses the base hook;
Humanoid and Human delegate to it. Stage BYTE+13c=0 admits
mover reference repair at+7c. These local paths do not by themselves establish
preservation through every callback. — SAV-LOADREG-878, SAV-LOADHOOK-879

The three measured actor classes share the+18 tick. It processes attached effects
before signed HP and order admission; HP>0 and actor+3c!=0 are necessary for the
selected order call. Within an admitted order prefix, progress+9=0 without the
status hold lets pending byte+8=10 choose the explicit turn arm. Unequal current
and desired facing call the turn with the same actor; equality clears the pending
byte. The inactive short-arc turn arm snaps without reading mover+a, whereas the
other arm reads that actor's allocated byte. — MOVE-EVENT-060

A selected continuous Effect callback can call the same actor's derive first:
its original+38 -> +40 -> generic+48 route passes the unchanged actor to+50.
Unit's derive differs from Humanoid/Human's, whose reached store replaces mover+a
with low8(actor+8c). Empty/ineligible effects and special identities can skip that
route. This conditional order does not identify the first restored effect or the
first post-LOAD mover access. — MOVE-EVENT-061

The selected frontend resume can bootstrap a server call, but earlier world,
frontend, phase-dependent full-tick, command and world-object callbacks remain
between LOAD and the selected actor consumer. Absolute first producer/read order,
the byte at that read and native resume timing remain Unknown. — SAV-FIRSTMOVE-880
