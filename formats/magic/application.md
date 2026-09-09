# Application and ordinary Effect lifetime

[Reference](format.md)

## Effect ownership and credit

An attached ordinary Effect has source state at+44 separate from the target's
credited actor at+40. Default Effect construction clears its source, copy
preserves it, and continuous same-id refresh retains the old source. Template
construction's null-parser arm is different: it does not write+44, and its
ordinary reach/allocation value remains Unknown. — SAV-906

The new-Effect SAV route uses the default factory and omits source+44 from its
44-byte record. An archive back-reference preserves the loaded Effect alias.
The selected post-load paths do not reconstruct its source. In contrast, the
target's raw saved credit key is explicitly remapped on world LOAD, or cleared
on a lookup miss; the no-world document arm clears it. — SAV-907, SAV-908

For nonzero Token8 periodic damage, victim HP changes before source admission.
A null source skips the source callback, negative source HP clears the source,
and zero/positive HP admit virtual+64. A duration-only timer does not make that
periodic call. Later death credit reads target+40 and can name a different actor;
an intervening source+48 callback precedes the source+60 award. Native lifetime,
later repairs and first-consumer chronology remain Unknown. — SAV-909, SAV-910

## Application

```
Spell::Apply(caster, target, x, y)                                  FUN_004feadb
  power, then FUN_004fe25d reloads and raises spell+0x09 and fills
  spell+0x0e / +0x0f / +0x10 (the spell-power fields); an item cast mutates its weapon-owned Spell too

  if (base + spread) != 0 and the spell is neither 6 (Heal) nor 11 (Drain Life):
        an Effect_DirectDamage (0x60 B) is built; its combat block at +0x48 gets
            block+0x13 = base, block+0x14 = spread, block+0x15 = the Sphere
        so SCHOOL i IS DAMAGE KIND i, and the block carries no physical damage

  ...and then RETURNS TO THE ATTACHMENT TAIL. The per-spell arm below is reached ONLY when the
  scratch pair is 0, or the spell is 6 (Heal) or 11 (Drain Life) -- so the seven spells with
  damage columns never execute an arm at all, and the arm they share (0x005001fc) is dead code.

  the per-spell arm, dispatched through the 29-dword table at 0x0050066f (004feec9), which
  collapses the 28 shipped ids onto 17 addresses:

        0x005001fc  fire_arrow fire_ball wall_of_fire acid_stream lightning
                    prismatic_spray meteor_storm            (never entered -- see above)
        0x004ffa2e  the four Protection spells              0x004ff2bb  fire_sacrifice
        0x004ffc2c  haste, slow                             0x004ff05b  heal
        0x004fff81  bless, curse                            0x004ff18e  drain_life
        0x004feef6  freezing_cloud                          0x004fefac  poison_cloud
        0x004ff8e8  light                                   0x004ff98a  darkness
        0x004ffb2a  shield                                  0x004ff88d  wall_of_earth
        0x004ffdb9  stone_curse                             0x005000e4  invisibility
        0x004ff53b  control_spirit                          0x004ff283  teleport

  Eleven arms build their Effect from the row's own `Effects` column (FUN_00501156 ->
  FUN_00502eee), which supplies the KIND and, where the arm writes no +0x40, the magnitude;
  the parser reads ONE record, so stone_curse's second (`defence=-20`) is never built. Five
  use the plain ctor FUN_00501105 instead -- bless, curse, invisibility, wall_of_earth and
  the dead default -- so those spells' `Effects` columns are never parsed at all.

  Selected shapes:
        the four Protection spells share one arm: magnitude = power/2, duration as above
        Bless / Curse                            magnitude = (power*4)/5 + 20, negated for Curse
        Haste / Slow                             magnitude = power/15 + 1,     negated for Slow
        Heal                                     refused across the diplomacy table
        Fire Sacrifice     base = min(health + mana, 255), spread = min(min(health + mana +
                                                 power, 512) - base, 255); leaves 1 health and
                                                 0 mana; stamped FIRE by a hard-coded setter
                                                 rather than by the Sphere table. It also
                                                 builds a protectionFire-100 / 32-tick effect
                                                 and NEVER ATTACHES IT -- no AddEffect, no
                                                 list insert, and the pointer is not read again
        Stone Curse        duration = T(10), then x (100 - target.protectionEarth)/100 with a
                                                 floor of 1 tick -- the only place in the image
                                                 where a resistance shortens an effect
        Bless / Curse      the magnitude is a PROBABILITY: FUN_004fbc92 takes the maximum (23)
                                                 or the minimum (27) of the damage roll when
                                                 magnitude > U[0,100], i.e. 20/101 at power 0
                                                 and 100/101 at power 100
        Prismatic Spray    the cast's id==14 arm calls FUN_004fe92e; its ordered selected-victim
                                                 list receives one complete apply per entry.
                                                 The common item wrapper refuses id 14 after a
                                                 caster-item route has already run this selector
        Drain Life                               moves health from the target to the caster
        Control Spirit                           kills a nearby actor whose +0x13c is 2 (health
                                                 := -10001, stage := 5) and raises a "Ghost"
                                                 from it -- a 0x198-byte actor built from the
                                                 Data.bin template named by the literal at
                                                 0x005c80ac, placed at the corpse's cell and
                                                 owned by the CASTER's Player: Reaction
                                                 halved+1, health and healthMax halved, combat
                                                 block and defence copied, and MIND AND SPIRIT
                                                 COPIED VERBATIM - the two stats are the only
                                                 fields not reduced. The power enters no
                                                 expression in this arm
        Teleport                                 moves the caster to the target point

  every lasting effect gets +0x3c = the kind, +0x3d |= 1 (duration), +0x40 = the magnitude,
  +0x42 = the ticks, +0x0c = the spell id, +0x0e = spellId*2 + 8 (an art index)

  attachment:  DistributionSystem == 1 -> a PointEffect, which REQUIRES a target unit; with
                                          none it prints "Spell, oops - can't cast point
                                          effect of x,y" and the effect is DROPPED.
                                          PointEffect+0x41 = (Defensive != 1), the inverse of
                                          cached Spell+0x0a = (Defensive == 1).
               otherwise               -> an AreaEffect sized by the Radius column (unscaled),
                                          with effect+0x4c = (AreaEffectDuaration << 4), and
                                          += (power << 4)/10 and +0x08 = 1 when that is
                                          nonzero; DistributionSystem 5 instead sets +0x08 = 2
                                          and +0x4c = 0
               DeliverySystem == 2     -> the whole thing is wrapped in a SpellTransport that
                                          flies at SpellEffectSpeed and delivers on arrival;
                                          Lightning and Prismatic Spray overwrite its counter
                                          with the literal 10
               anything else           -> nothing is attached at all
```

**Which spell is aimed at what.** `Spell Target` (title 5, `getParam(row, 4)` at `FUN_004fe0f9`)
routes the *cast*: 1 carries the target unit, anything else carries the point. `Distribution
system` (title 9) shapes the *apply*. They agree on 27 of 28 rows — `Shield` ships `Spell Target 2`
with `Distribution 1`. Nothing in the area collector `FUN_0053ddd0` consults the diplomacy table:
**`Heal` is the only arm in the whole apply that does**, and an area spell burns its own party.

**Attachment to an actor** (`FUN_005014ae`) uses Effect Token+0x0c as identity, including
non-spell Potion id 0, and sets `actor+0x144 |= 1 << id` — the bitmask the
damage resolver reads for Bless and Curse and the cast reads for Invisibility. **Ordinary timed
attachments share one record per id**: an effect of the same id already present is un-applied,
given the new magnitude and duration, and re-applied (an incoming `continuous` one has only its
counter refreshed; `MAGIC-ATTACH-016`, `MAGIC-POISONREFRESH-158`). **Bless and Curse
annihilate**: casting either on an actor carrying the other removes that one and applies nothing.
`FUN_0050134f` ticks `+0x42` down, clears the id bit on expiry and ORs in a fifth `+0x3d` bit
(`[0x0059bc18] = 128`) that is not one of the five duration words; a `+0x42` above 9600 never counts
down at all, which no shipped spell can reach.

See [overlap rules](overlap.md) for Fire Wall/Poison Cloud and
[training and weapon casting](training.md) for award producers.

### Non-spell Potion lifetime

`MAGIC-CONSUME-142`, `MAGIC-CONSUME-143` and `MAGIC-CONSUME-144` qualify this attachment mechanism. Only mode bits 1
and 2 are timed; charges bit 4 alone is not. Item use turns mode 0 into singleuse 8.
Current-health/mana potions add 30 or 100, capped at maxima. The four +1 attribute
potions change the live attribute without changing its modifier byte and leave no
timer. The common effect-dispatch tail immediately invokes Human derive, which caps it at 50 plus signed modifier.

The five timed rows are absorption+50 for 480 ticks, health/mana regeneration+100
for 960, and the two regeneration bonuses+250 for 1920. All retain Token id 0.
On repeat, attachment finds the old id 0 even if the new characteristic differs.
Non-continuous replacement un-applies the old effect, overwrites only magnitude
and duration, then re-applies the old characteristic. It neither appends an
independent timer nor changes old kind/mode. Incoming continuous refreshes duration
only. A class-inapplicable mana effect can still consume its Potion.

The actor ticks attachments before the health branch, except terminal act `0x10`.
Expiry un-applies non-continuous effects, clears the id bit and marks removal;
continuous expiry does not unapply. Saves carry the attachment's remaining
counter, kind/mode/id, live attributes and the already-applied modifier block.
These readers do not replay a fresh Potion or reset its timer. Native drink-save-load behavior and town-to-mission preservation remain Unknown.

### Prismatic Spray's selected victims (`MAGIC-SPRAY-134`…`MAGIC-SPRAY-137`)

`FUN_004fe92e` constructs a local output list and calls
`FUN_0053ddd0(session,caster,primary,out,cap)`. The cap is
`(u8)min(power/20+2,7)`: book power comes from Sphere 2, while an item/staff uses the signed
kind-`0x29` effect's `+0x42`. It is a final victim count, not a distance.

The selector first alarms the primary and applies the directional hostility flip. It asks
`FUN_005365e0` for the caster group's shared candidate lists: sight is cleared once and stamped by
every member; the global actor list is walked head to tail; the first member supplies diplomacy and
the See-invisible exception; health-below-one candidates move from A to B. It then copies A followed
by B into a 100-entry pointer region at `[ESP+0x60]`, with no bound check. Parallel scores are written
to `[session+0xd74+4*i]`; this scratch region's capacity is Unknown. A score is
`((edgeDistance<<8)+turnCost)&0xffff`; B score is
`((((edgeDistance<<8)+turnCost)&0xffff)<<8)`. If A was empty, the builder has already moved all of B
into A, so an only-corpse population uses the A formula. With living A, B remains in its own list;
the bytes do not prove that its transformed score can never pass the selector thresholds because
the reachable edge-distance domain was not established.

On the nonempty-A path, up to `cap` winners are selected by repeated strict-minimum scans. Equal
scores preserve source-list order. Ten winner indices and ten score thresholds begin at 65530. A
winner index below 60000 gates only stamping its score scratch to 65500; append separately requires
both the winner index and saved threshold below 65000. The primary is appended first, selected
secondaries follow in rank order, a secondary equal to the primary is skipped, and the tail is
removed until `count <= cap`. Thus the primary bypasses secondary visibility, diplomacy and corpse admission but
still consumes one place. With nonempty A, cap zero removes it. With empty A, however, an earlier
arm appends the primary and returns before reading the cap. Ordinary book/item power produces 2..7,
so this exception does not exceed the authored cap. Neither pointer-copy loop nor a custom cap above
ten is checked. More than 100 candidates leaves the pointer stack region; a cap above ten leaves the
winner/threshold regions. Whether a candidate score leaves its session-owned scratch region is
Unknown.

The same ordered list is applied, then sent. Opcode `0x8a` carries an ordinary caster; `0x8c`
carries a packed cell for a temporary caster. Both client arms refill source `+0xac/+0xb0` with the
victim ids in order, and picture 36 draws one link for each id. The list is therefore both the
simulation apply order and the visible branch order.

`Spell+0x09` is separate. Max Range plus the power bonus is copied to order `+0x14` and controls
cast approach/admission. It neither limits group sight nor enters `FUN_0053ddd0`. Script instant 24
structurally creates the unit-targeted state that can converge here, but neither preserved campaign
root authors spell 14 in instant 24. Instant 21 makes the point-targeted state with a null target;
a custom spell-14 record is predicted to fault at the selector's immediate primary dereference.
That prediction was not run, and temporary-caster group ownership remains Unknown.

## Effect duration

`effect+0x3d` carries the game's own five words, parsed by `FUN_00503890`:

```
permanent   0     no counter
duration    1     +0x42 counts remaining ticks
continuous  2     +0x42 counts, and the effect re-applies
charges     4     +0x42 counts uses
singleuse   8     applied once; does NOT raise a stat's cap byte (HERO derived-stat step 0)
```

Any of `1 / 2 / 4` also makes the effect dispatch read the magnitude as a **signed 16-bit** value at
`+0x40`; without them it is a signed 32-bit value.
