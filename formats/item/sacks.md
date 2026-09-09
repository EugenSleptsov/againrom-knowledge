# Sacks, death and loot

[Reference](format.md)

## Sacks

| Off | Type | Meaning |
|-----|------|---------|
| +0x04 | u32 | runtime id, from `MOVE-ID-016`'s bitmap |
| +0x10 | ptr | position |
| +0x1c | i32 | total value = gold + `Σ item+0x1c`; recomputed by `FUN_0050f4b3` |
| +0x3c | i32 | gold |
| +0x40 | ptr | a container of the same class an actor holds |

**One sack per cell.** `FUN_0050f5aa(pos, container, gold)` looks the cell up first; when a
sack is already there it pours the incoming container into it and adds the gold. A new sack
that fails to register at its cell is deleted and the call returns 0.

The local registration key is `Sack+10 -> Position+02`, not Position's other
cell word. Dynamic bit 0 rejects before lookup. An existing nonzero Sack slot
rejects even when it already holds the supplied Sack. An existing empty slot
is written without plane recomputation; creation instead zeroes 52 bytes,
captures current Cost/Static, sets the record-present flag, stores the Sack
and recomputes. — SAV-SACKENTRY-590

Local removal refuses only a missing record before its write path. It clears
payload+10 without an occupied-slot or pointer-equality test, recomputes, and
tests the four occupant slots, layer count+02 and operation+2c for deletion.
Other residue and the layer pointers themselves do not retain the record.
— SAV-SACKREMOVE-591

Deletion restores Cost and Static baselines, preserving Static bit 4, but does
not replace Dynamic with restored Static or clear Dynamic bit 5. Dynamic keeps
the preceding recompute result, with the conditional bit 4 OR. This is a local
write-set, not proof that allocator or higher-level callbacks preserve every
other field. — SAV-SACKPLANES-592

The lookup used by merge/create first requires Static bit 5 and then a matching
node. Registration success dispatches the collection append; refusal dispatches
the deleting-destructor slot. Two known removal-caller slices ignore the local
return before subsequent collection/transfer calls. Complete merge, destructor,
transfer and exceptional effects are separate from these established cell
transitions. — SAV-SACKCALLER-593

**A sack does not tick and does not expire.** Its `vt+0x14` and `vt+0x18` are empty stubs and
it is never inserted into the actor tick list — it is created, merged into, and destroyed by a
pick-up. Nothing ages it.

A pick-up takes **everything**: the gold is credited to the looter's `Player`, the sack's whole
container is poured into `actor+0x7c` one unit at a time, the sack is unregistered and deleted.
There is no per-item take from a sack in the protocol, and no capacity, distance or ownership
test at execution.

## Death

`FUN_004f4f5d`, in order:

```
1. state := 16; leave the world
2. unequip actor+0x78 into the container
3. unequip actor+0x74 into the container ONLY IF the weapon's Data.bin parameter 15
   ("sutableFor") != 0 -- i.e. unless the weapon suits no class at all
4. actor->vt+0x44()  -- the STRIP. Empty body on the base actor; on both humanoid
   classes FUN_004f70f8: for i = 1..12 inclusive, unequip actor+0x198+4i into the
   container. A null slot costs nothing: vt+0x40 returns 0 and the append returns
   before storing. This is step 4 and everything below it sees the result.
5. suppress := templateName contains "NPC"   OR   (multiplayer AND Player+0x5c != 0)
6. if suppress:  DELETE the container outright -- with the armour already in it --
                 and install an empty one
7. gold := 0; if typeID > 0x40 and rand()%100 < param[0x26]:
               gold := param[0x27] + rand()%param[0x28]
8. if container is non-empty OR gold != 0:
       no existing Sack: the new Sack ADOPTS the container object itself
       existing Sack: DRAIN into its container, deleting the source container
9. the corpse is given a fresh empty container
```

The container identity is retained only by the new-Sack adoption branch.
Strip/reinsertion and an existing-Sack drain apply the whole/split/merge rules;
Weapon unequip also changes its owned Spell. `ITEM-DEATH-012`, `ITEM-GROUNDMOVE-130`
and `ITEM-SPELLMOVE-132` distinguish these outcomes. A mercenary (`NPC%02d_%d`)
leaves nothing at all. Because the strip precedes the
emptiness test, a body that wore anything leaves a sack even if it carried nothing.

**Death is not overridden per class.** `vt+0x18`, the tick that reaches `FUN_004f4f5d`, is
`FUN_004f37be` on all three actor vtables. Only `vt+0x44` differs, and its empty base body is
structural rather than stylistic: `Unit` is `0x198` bytes, so the loop's first read on a base
actor would be four bytes past the end of the object.

**Disposal is a different path from dropping.** `Humanoid::~Humanoid` (`FUN_004f6f09`) also
walks `1..12`, but calls `worn[i]->vt+0x04(1)` -- the deleting destructor -- and never touches
the container. On a normal death it finds twelve nulls, because the strip ran first.

## Sack producers

Five owners, complete on the repaired function table:

```
FUN_0050f6d2   the wrapper used by death and by both drop arms
FUN_004d5dd8   command 0x23, at the requested cell and at the hero's cell
FUN_004f2176   the mission .ini's "Items" list, at map load  (no .ini ships)
FUN_004f1fc6   random treasure at map load, and a multiplayer top-up to W*H/400
FUN_004e4f3e   .alm type-8 authored loot, at map load
```

The repaired `EnumRefs` reference-manager and `.rdata` populations contain three direct
packed-item-factory owners: `.alm` type-8 load, script instant 12 and a transient UI formatter that
immediately destroys its object. The sack maker has six direct calls in five owners; the death
inventory routine has one caller. Computed targets remain the bounded blind spot of those image
reference populations (`ITEM-PRODUCER-091`).

## Authored staffs, Dragons and death

Every Human equipment cell containing `castSpell` is a slot-0 Weapon cell: 46 per root, belonging
to 36 ordinary template names and ten with the exact uppercase substring `NPC`. Every underlying
weapon is `Staff` or `Shaman Staff`, both `sutableFor = 2`. Across the 28 campaign maps, 76 Human
placements resolve to those definitions per root. Death moves the same worn Weapon into the actor's
container before producing a sack for 49 non-`NPC` placements; it destroys that container for the
27 `NPC` placements. This split is for the shipped single-player campaign maps. In multiplayer,
`Player+0x5c != 0` is a second suppression input; loose-map non-`NPC` rows therefore leave death
eligibility unevaluated (`ITEM-AUTHCAST-086`, `ITEM-AUTHDROP-087`, `ITEM-DEATH-012`,
`ALM-MODE-070`).

Dragon rows 112–115 author `Flame Thrower` in slot 0, leave slot 1 empty and
have no castSpell. `Flame Thrower` has `sutableFor=0`, so the death gate does
not move it into the container. No Dragon staff or enchanted-item death
producer is established in type 8 stock, same-cell ground elements,
item/container script targets, the shared death body, or its enumerated
direct/reference callers. This negative is limited to those families.
Flame Thrower is innate-weapon state. — ITEM-DRAGON-088, ITEM-DRAGDROP-089
