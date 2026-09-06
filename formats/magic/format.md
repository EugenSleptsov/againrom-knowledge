# MAGIC — the spellbook, the cast and the resistance

Level 3 spec. Promoted from [`claims/magic.md`](../../claims/magic.md), including
`MAGIC-SPELL-001`…`MAGIC-EFFMODE-009` and `MAGIC-CADENCE-126`…`MAGIC-CADENCE-127`.
Not a file format — the simulation area that turns
a learned spell into a wound or a buff. Everything here was read out of `rom.exe` at instruction
level; the addresses are in the cited claim rows. The actor's
own fields are [`formats/hero`](../hero/format.md); the parameter table is
[`formats/databin`](../databin/format.md).

`ftol` truncates toward zero (`0055458c`, rounding-control 11); `IDIV` truncates toward zero.
`rand(n)` is `FUN_00504003`: `rand() * (n+1) >> 15`, uniform on `[0, n]` **inclusive**, 0 when
`n == 0`. Runtime `Data.bin` column *i* is title *i+1* throughout.

## 1 — the spell table

Direct-damage effect delivery to a Building differs from a physical actor action:
`FUN_00502c51` selects the Building resolver, subtracts from word `+0x42`, clamps a
negative result to zero and notifies only for positive damage. The physical Building
strike does not clamp and has a different notification predicate. The shared resolver
does not make the two delivery paths interchangeable (`UNIT-STRUCTDELIVERY-065`). This
does not establish every spell's Building admission or area-effect multiplicity.

`World\Data\Data.bin` group **H spells** ships **28** rows, indices 1..28. `rom.exe` carries the
same names as a 29-entry array at `0x005c5b58` whose entry 0 is `unused_spell_0`; the two agree
28/28 by name, and the id space *is* the array index. Spell 0 is refused at every entrance.

```
 1 Fire Arrow            2 Fire Ball             3 Wall of Fire        4 Fire Sacrifice
 5 Protection from Fire  6 Heal                  7 Freezing Cloud      8 Poison Cloud
 9 Acid Stream          10 Protection from Water 11 Drain Life        12 Light
13 Lightning            14 Prismatic Spray      15 Invisibility      16 Protection from Air
17 Darkness             18 Shield               19 Wall of Earth     20 Stone Curse
21 Meteor Storm         22 Protection from Earth 23 Bless            24 Haste
25 Control Spirit       26 Teleport             27 Curse             28 Slow
```

The columns the cast path reads, by runtime index (title in quotes):

| idx | title | used for |
|---|---|---|
| 0 | `Complication Level` | added when cast state `0x0d`/`0x0e` still has a non-null `actor+0x64` Spell at recovery; weapon-diverted casts clear it first |
| 1 | `Mana Cost` | the whole cost, cached at `spell+0x0c` |
| 2 | `Sphere` | the school, 1..5 = Fire/Water/Air/Earth/Astral — also the damage kind |
| 5 | `Delivery System` | 1 = attach now, 2 = a flying missile |
| 6 | `Max Range` | cached at `spell+0x09`, then raised by the power |
| 7 | `Spell Effect Speed` | the missile's speed; the delay is `distance / speed` |
| 8 | `Distribution system` | 1 = a point effect on one unit, else an area effect |
| 9 | `Radius, Length/2` | the area effect's size |
| 11 | `Area Effect Duaration` | the fallback duration when `Spell Duration` is 0 |
| 14 | `Spell Duration` | the duration in game seconds, ×16 for ticks |
| 16, 17 | `damageMin`, `damageMax` | the damage pair, scaled by the power |
| 18 | `Defensive` | cached at `spell+0x0a` as `(value == 1)` |

## 2 — a `Spell` and a spellbook

```
Spell        0x14 bytes, vtable 0x0059c670
  +0x00  vtable
  +0x04  pointer to the Data.bin Spells row   (restored from the id on load)
  +0x08  u8   spell id 1..28
  +0x09  u8   Max Range   -- raised per cast
  +0x0a  u8   Defensive
  +0x0c  u16  Mana Cost
  +0x0e  u8   per-cast scratch: damage base
  +0x0f  u8   per-cast scratch: damage spread
  +0x10  u16  per-cast scratch: duration in ticks -- WRITTEN AND NEVER READ. Every arm that
              needs a duration recomputes it from the column; over the twelve routines that
              take a Spell* the only `word ptr [reg + 0x10]` is a stack argument.
Serialize stores +0x08, +0x09, +0x0a as bytes and +0x0c as a word. Nothing else.

Spellbook    0x1c bytes, at actor+0x140
  +0x00  vtable
  +0x04  CObArray (0x14 B) SUBSCRIPTED BY SPELL ID; a null element means "not known"
  +0x18  int, the id of the last successful lookup
```

Learning a spell is the effect kind `teachSpell` (42) and nothing else: with no book it does
nothing, with the id already present it does nothing, otherwise it constructs `Spell(id)` and
stores it at `book[id]`, then sets `actor+0x150 |= 0x400000`. Allocating a book sets `actor+0x4c`
bit 1, and bit 2 — the mage bit — when the actor has a nonzero `manaMax`.

## 3 — the power: what a school skill buys

```
power = clamp( skill[Sphere] + Mind - 30, 0, 100 )       skill[] is actor+0xa8 + 2i
      = the item's castSpell effect +0x42                when casting from an item -- UNCLAMPED
f     = power/30 + 1

damage pair   base   = ftol( damageMin * f )   ONLY when the column is > 0, else 0
              spread = ftol( damageMax * f ) - base    ditto -- a SPREAD, not a maximum
              so a -1 column scores 0, which is what lets a spell reach its own arm at all
duration      ftol( 1.025^power * SpellDuration * 16 ) ticks
              invisibility instead: min( ftol(1.05^power * 3 * 16), 65000 ) -- the literal 3,
              NOT its own SpellDuration column, which is read only as the > 0 gate
              when SpellDuration is 0: (AreaEffectDuaration << 4) + (power << 4)/10
cast range    MaxRange + power/30       (Teleport: MaxRange + power/3)
              ONLY when the MaxRange column is nonzero (004fe2b4): Fire Sacrifice and
              Shield ship 0 and gain nothing at any power
spray victims min( power/20 + 2, 7 ) final victims, primary included
              This lives in FUN_004fe92e, whose only caller is the cast's spell-14 arm. It is
              PRISMATIC SPRAY's output-list cap, not a radius. The ten AreaEffect spells take
              their radius from the Radius column, unscaled; Prismatic secondaries instead
              come from the caster group's visibility population.
```

Mind is worth exactly as much as the school skill. **There is no `+3 × skill` to-hit term and no
`skill/5` damage-floor term for a spell** — the weapon-skill terms of `formats/hero` step 11 have no
counterpart here, and a spell's damage never passes the to-hit roll at all.

**The domain, and the one place the clamp is missing.** Both terms are hard-bounded: a school skill
is clamped to `[0,100]` twice in the derive and a stat is clamped to 100 by every effect arm, so
`skill + Mind - 30` runs over `[-30, 170]` and **the `[0,100]` clamp binds from `skill + Mind = 130`
upward — it is reachable, not defensive**. Three routines compute the expression and they do not
agree: `Spell::Apply` and the Prismatic Spray fan clamp; `FUN_004f5998` returns `max(x, 0)` in `AL`
with **no upper clamp**, and its single caller uses it for the spray victim cap. That makes no observable
difference — the victim count saturates at 7 at power 100, exactly where the clamp begins, and over all
10 201 reachable `(skill, Mind)` pairs the two forms never disagree — but a consumer that clamps
everywhere is right by luck, not by construction. An untouched stat is capped at **50**
(`formats/hero` step 0), so **without a cap-raising effect Mind contributes at most 20 points of
power**.

At the ceiling, in the integers the engine stores: `f` tops at **4.333**, so damage is 4.25×–4.33×
the `Data.bin` columns; a duration column multiplies by `1.025^100 = 11.81`; range gains at most
**3** cells and Teleport **33**; Prismatic Spray reaches **7 total victims**. Nothing overflows its store on the
shipped table — the largest damage base is 43 and spread 87 against byte fields, the largest
duration 6 312 against a `u16`.

## 4 — casting

```
Spell::Cast(caster, target, x, y)                                   FUN_004fe6d3
  if the caster is a mage AND caster+0x68 == 0   (i.e. not an item cast):
        if spell.manaCost > caster.mana:  refuse, nothing happens
        caster.mana -= spell.manaCost
  a fighter is never charged; an item cast is never charged; the cost never scales
  casting at anyone but yourself sets any active invisibility's remaining duration to 1
  DeliverySystem == 2:  delay = distance(caster, target) / SpellEffectSpeed
                        Lightning and Prismatic Spray: delay = 5, flat
  else                  delay = 0
  the apply is queued on the world object with that delay
```

The cast itself runs through the actor's common action phases. Phase 0 loads charge
`actor+0x134`; phase 5 applies when that countdown reaches zero; phase 7 loads relax
`actor+0x135 + U[0,3]` and the equipped-Humanoid weight/Reaction penalty from
the [hero format](../hero/format.md). It also adds this Spell row's `Complication Level` when
state is `0x0d` or `0x0e` and `actor+0x64` is still non-null. For an uninterrupted retained
order cast, application-to-application cadence is

```
charge + relax + U[0,3] + humanoidPenalty + ComplicationLevel + 2 actor ticks
```

The caster weapon-diversion route enters state `0x0d`, but application clears `actor+0x68` and
`actor+0x64` before phase 7 tests the Spell pointer. Its cadence therefore omits
`ComplicationLevel`:

```
charge + relax + U[0,3] + humanoidPenalty + 2 actor ticks
```

The final two ticks are the completion latch crossing the order-before-action boundary, not another
stored timer. Insufficient mana returns from phase 0 before changing phase, countdown or the
completion latch, so it writes no failed-cast recovery. If a prior completion left that latch set,
the retained order attempts admission for three actor ticks, skips one while progress consumes the
stale latch, then repeats. With a zero latch, refusal continues to be attempted every actor tick.

## 5 — applying

```
Spell::Apply(caster, target, x, y)                                  FUN_004feadb
  power, then FUN_004fe25d reloads and raises spell+0x09 and fills
  spell+0x0e / +0x0f / +0x10 (section 3); an item cast mutates its weapon-owned Spell too

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
damage resolver reads for Bless and Curse and the cast reads for Invisibility. **Spells do not
stack**: an effect of the same id already present is un-applied, given the new magnitude and
duration, and re-applied (a `continuous` one has only its counter refreshed). **Bless and Curse
annihilate**: casting either on an actor carrying the other removes that one and applies nothing.
`FUN_0050134f` ticks `+0x42` down, clears the id bit on expiry and ORs in a fifth `+0x3d` bit
(`[0x0059bc18] = 128`) that is not one of the five duration words; a `+0x42` above 9600 never counts
down at all, which no shipped spell can reach.

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
These readers do not replay a fresh Potion or reset its timer. A new observed
drink-save-load sequence and town-to-mission preservation remain Unknown.

**Training is event-based, not one award per release.** Every apply from a book —
`caster+0x4c & 2` and `caster+0x68 == 0` — calls `caster->vt+0x68(target, spellId)`, an empty stub on
the monster vtable and `FUN_004f7980` on both hero ones. It awards `round(manaCost / 2)` into the
spell's `Sphere`. An item cast keeps `actor+0x68` non-null and gets **none** of that immediate award.

Later `vt+0x64` calls are independent: positive direct damage supplies the resolved amount; Drain
Life supplies the positive capped transfer; Slow id 28, Stone Curse and Curse id 27 each supply
`trunc(target.healthMax × 0.03)` once per accepted target before attachment; Poison Cloud id 8
supplies every nonzero signed continuous tick while its recorded caster has health `>= 0`. Negative
caster health clears the pointer; zero health does not. The award is
`trunc(target.XPvalue × 0.5 × amount / target.healthMax + 1)`. Before that calculation, the damage
caller requires a victim owner with `+0x28 != 0`, and refuses when owner `+0x5c` and multiplayer
`server+0x0c` are both nonzero. The sink first requires recipient
`typeID` in `[0x21,0x3f]`; low-type map Humans receive no award. After that gate, with
`actor+0x4c & 4` set it maps the actual spell id to `Sphere`; with the bit clear the sink substitutes
`actor+0xb6`, the current weapon skill. Thus an admitted caster item trains a school only through
admitted later events, while an admitted fighter rider trains the equipped weapon for its physical
hit and every spell-side event.

There is no early refusal when Mind scaling produces zero. The chosen arm adds zero and runs the
ordinary refold. Nor is a negative amount refused: a custom negative item power can invert Poison's
signed magnitude, so the tick heals through subtraction and still calls the award path. The positive
cap leaves a negative value unchanged; the chosen arm subtracts it from slot and aggregate XP but
has no level-decrement arm. Shipped positive-power Poison items do not exercise that G2 boundary.

`vt+0x60` is a delayed kill award. Its caller requires a victim owner and applies the same
owner-`+0x5c`/multiplayer refusal, but not the damage caller's owner-`+0x28` test. The direct-damage
resolver leaves prior attribution for a null source or null source definition, clears `victim+0x40`
for a definition-bearing source without an owner, and otherwise writes the source. It then writes
damage kind for a bit-4 source or zero for a clear-bit source. A
PointEffect enters its post-payload tail only with a recorded caster and `+0x41 != 0`, making it
non-Defensive only. A now-null caster definition actively clears prior `victim+0x40`; otherwise a
non-null owner admits the recorded-caster and actual-id writes, while a null owner leaves prior
state. An AreaEffect uses a separate tail: Wall of Earth returns before payload, Light after payload
but before attribution, and every other arm requires a nonzero target movement-domain byte plus a
recorded caster. A null definition or null owner clears `victim+0x40`; a surviving owner admits the
recorded-caster and actual-id writes. The caller zero-extends the byte, so ordinary domains 1..3 and
custom `0xfe` pass, while only zero fails. Movement domain is not health, so a lethal application on
an ordinary domain is attributed. Drain Life jumps to the apply epilogue
without either envelope and writes no
fresh attribution. Poison Cloud seeds id 8 during area application, but later ticks do not refresh
it; intervening combat can therefore change its eventual kill recipient. An **area** application
remains per target, and Poison Cloud is per qualifying tick; one release may pay zero, one or many
awards.

The authored Data.bin equipment corpus contains 48 castSpell weapon cells on each root: 46 Human
staffs and the Catapult/Ballista Unit weapons. Its 18 powers are
`{1,5,10,15,25,30,34,35,40,50,60,63,65,70,82,90,98,99}`; they all enter the same unclamped signed
`i16` item-power route. The live shop is an additional producer: unflagged generated weapons choose
ids `{1,11,13,14,20}` and a price-derived random power capped at 100. Plain Unit vtable
`0x0059c3c0` maps kill, damage and cast-award slots to empty
stubs, so those siege actors can release the rider but never train. Humanoid recipients then face
the type-id and later gates above.

An item apply passes the weapon's existing `Spell` to the same fill as a book cast. It reloads and
power-adjusts serialized `Spell+0x09` Max Range, then rewrites `+0x0e` damage base, `+0x0f` spread
and `+0x10` duration scratch. The item effect, `weapon+0x80` pointer and `Spell` identity survive;
the four filled fields do not remain unchanged.

`Unit::Serialize` preserves training as six levels, six per-skill experience dwords and one
aggregate. It also stores `actor+0x68` as an object reference, while `actor+0x64`, kill-credit
`actor+0x40` and attribution `actor+0x48` are raw `u32`, raw `u32` and `u8`. The credited actor is
therefore not remapped as an archive object reference on load.

**Prismatic Spray** (`FUN_004fe92e`) selects at most `min(power/20 + 2, 7)` total victims and
applies to each in selected order. Its secondaries come from group sight; this count is not a
radius and the ten AreaEffect spells use a different path.

The per-spell arms are **power consumers too** — twelve distinct expressions in all, of which the
four above are the shared ones. The magnitudes an arm writes to `effect+0x40`:

```
Freezing Cloud        -(power/15 + 1)          Shield              power/10 + 3
Light                   power/30 + 1           Haste / Slow      +-(power/15 + 1)
Darkness               -(power/30 + 1)         Bless / Curse     +-((power*4)/5 + 20)
the four Protections    power/2                Poison Cloud       ftol(template * f)
Fire Sacrifice       two totals, each capped at 0x200
```

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

## 5a — what decides whether an AI actor casts

Only a **mage** reaches the choice (`actor+0x4c & 4`, plus `[actor+0x14]+0x28 != 0`), and then:

```
if (Mind > 59 and rand()*100/0x8000 < 30):  hand back to the action dispatcher, cast nothing
otherwise, over ids 1..28 through Spellbook::Get:
        keep a spell iff  ManaCost <= current mana  AND  Defensive == 0
        pick one uniformly, order it
```

So **a monster never casts a defensive spell**, never chooses by school, skill or power, and Mind
above 59 makes it cast *less often*, not better. `0x8000` is the AI class's own RNG range, so both
rolls are exact.

## 6 — an effect's duration model

`effect+0x3d` carries the game's own five words, parsed by `FUN_00503890`:

```
permanent   0     no counter
duration    1     +0x42 counts remaining ticks
continuous  2     +0x42 counts, and the effect re-applies
charges     4     +0x42 counts uses
singleuse   8     applied once; does NOT raise a stat's cap byte (formats/hero, section 3 step 0)
```

Any of `1 / 2 / 4` also makes the effect dispatch read the magnitude as a **signed 16-bit** value at
`+0x40`; without them it is a signed 32-bit value.

## 7 — a weapon that casts

An item may carry an effect of kind `castSpell` (41). That effect is **never applied** — the
dispatch's arm for it is a bare jump to the epilogue. It is read as data instead:

- `FUN_005089a7(item)` finds it by kind, and its `+0x42` is the power, replacing section 3's
  expression;
- `effect+0x40` is the spell **id**, a byte; `FUN_0050e7a6` reads it and `FUN_0050e6f3` hangs a
  constructed `Spell` on `weapon+0x80`. Nothing anywhere resolves that spell by **name** — the
  name-taking twin `FUN_0050e603` is dead code, unreached by a `CALL`, by a pointer table in any
  section, or by any `rel32` in `.text`;
- **which class fires it decides which of two paths runs, and both exist.** The predicate is one:
  `actor+0x4c & 4` (set = caster, fixed at its write — it is OR'd in only when the streamed
  `ManaMax` column is positive, and the spellbook is allocated in the same arm).

```
actor tick FUN_004f37be, state 3 (attack), weapon != 0, weapon+0x80 != 0:

  caster    (actor+0x4c & 4)  -> state := 0x0d, +0x64 = weapon+0x80, +0x68 = weapon
                                 the strike is NOT called at all
                                 FUN_004fe6d3 validates; wind-up = actor+0x134 frames
                                 then Spell::Apply                       <- the cast
  non-caster                  -> FUN_004fb942 -> FUN_004fba0e            <- the strike
                                 ordinary rider iff dealt > 0 and post-hit health > 0
                                 Fire Ball rider when either test fails too
                                 Spell::Apply                            <- the rider
                                 then physical XP iff dealt > 0, saved pre-hit gate clear,
                                      and post-spell health > -10
```

`FUN_004fb942` has exactly one caller and is the tick's **else**, so a caster wielding such a
weapon reaches no to-hit roll, no damage roll and no absorption: it **cannot miss** and deals no
weapon damage. Both paths clear `+0x64`/`+0x68` afterwards, and an item whose kind is `0x0e` is
destroyed with its `Spell` after the cast — a staff is not.

The selected-inventory item Cast command is not a third shipped weapon path. Its UI arm requires
descriptor bits `0x10|0x01`; the effect walk supplies bit 4, but the Weapon descriptor clears the
byte and supplies only bits 1 and 2 from `sutableFor`. Bit 0 comes only from a base MagicItems
descriptor. The server's latent item-order arm would accept a crafted kind-2 command, so this is an
engine/client G2 boundary, not a Data.bin switch.

Fire Ball is the fallback rider whenever the positive-damage/positive-health conjunction fails. It
therefore triggers after a miss or zero damage and on living as well as dead targets. The item spell
runs before physical experience; if it moves health to `-10` or below, the later physical award is
suppressed.

Spell id 14 separates the routes. The caster validator runs the Prismatic Spray fan during
admission while the item context is live; the later common wrapper refuses id 14 and prevents a
second fan. The fighter rider never enters that validator and reaches only the refusing wrapper,
so its id-14 rider applies nothing.

The **generated** starting kit gives a caster a staff and only a staff (`Wood Staff
{castSpell=Fire_Arrow:20}` above tier 10, `:10` at or below, plain or `Uncommon` by gender) and a
non-caster one of five swords/axes/maces/pikes/bows. That is a distribution, not a rule: nothing on
the equip path refuses an item on account of class.

## 8 — what Mind and Spirit do, and what they do not

Enumerated image-wide (`EnumRefs disp:88`: 461 hits / 235 owners / 4 orphan, all four stack frames;
`disp:8a`: 36 / 13 / 0), with the actor's `u16` width used as the filter and every byte- and
word-wide hit read:

```
Mind    -> the power term, one for one with the school skill    (three sites, section 3)
        -> Prismatic Spray's victim cap, through the one unclamped helper
        -> whether an AI mage casts at all                      (section 5a)
        -> Control Spirit copies the TARGET's Mind into the Ghost
Spirit  -> the mana pool, through manaMax        \ each in ONE hop, inside the derive
        -> the five protections, through spirit/2 /
        -> Control Spirit copies the TARGET's Spirit into the Ghost
```

**Neither stat touches**: cast admission, the mana gate, the mana cost, the cast delay, the to-hit
question, the damage or resistance arithmetic itself, mana regeneration (whose base is `manaMax`),
effect stacking, or the spellbook — **the book's capacity is not bounded by a stat, what may be
learned is not chosen by one, and learning is not gated by one**. There is no cooldown, no
interruption and no failure chance on the cast; the only refusal is the mana gate.

The blind spot this enumeration cannot close: a wholesale `REP MOVSD` copy of an actor would carry
no displacement and no `disp:` sweep can see it.

## 9 — what a consumer still cannot reproduce

The `AreaEffect` per-tick consumer, including its staged geometry, is specified in sections 11–13.
**What makes an invisible unit unseeable**: bit 15 of `actor+0x144` has two located consumers
(`FUN_004fb853`, `FUN_004fe6d3`) and both only *cancel* it.
**What `caster+0x3c` is** (`FUN_00523360`): a nonzero value suppresses every `Delivery != 2` cast
outright and drops the caster from an effect's `+0x44`.
What builds an `Effect` from a `Data.bin` **Magic** row beyond its duration word — the kind lookup
and the one-record limit are read, the numeric value parse is not, which is why `Poison Cloud`'s
`-2` and `Stone Curse`'s `+5` are taken from the CSV rather than from the disassembly.
The value→word mapping of `Distribution system` (`Point` / `Round` / `Long` / `Phase` /
`Hang On Unit` are located at `0x005c6c10` and the shipped values do not fall in that order). What
writes the `block+0x11`/`+0x12` damage pair, which lands even on a missed swing. ~~The projectile's
per-tick step and its art~~ — closed by § 10 and `formats/anim/format.md`.
How an actor first acquires a spellbook, and who authors a weapon's spell — the two routines
that do both have no caller on any enumerated route. What `actor+0x13c` is, the byte `Control
Spirit` selects its victim by. What the AI does with the spell it picked (`FUN_0052eea0`).

## 10 — the two pictures a spell has

Both are settled by `MAGIC-ICON-024`, `MAGIC-ICON-025`, `MAGIC-PIC-026`, `MAGIC-PIC-027`.
They share a spell id and nothing else.

**The icon is a position, not an index.** `graphics\interface\SpellBook.bmp` is 480x85, 24bpp,
and the panel blits it **whole and once**. Nothing in the engine cuts it. What is per-spell is a
mask: 24 cells at

```
x = xBase + 6 + 38 * (slot % 12)
y = yBase + 6 + 38 * (slot / 12)      slot = 0 .. 23, cell 36 x 36
```

and every cell whose known-spell bit is **clear** has `SpellBack.bmp` (36x36) pasted over it. The
selected slot gets a 36x36 outline; four hot-key slots get a numeral. The
book's slot-to-spell mapping is a fixed 24-entry table:

```
slot   0  1  2  3  4  5  6  7  8  9 10 11 12 13 14 15 16 17 18 19 20 21 22 23
id     1  2  3  4  5 23 24 16 15 14 13 12  6  7  8  9 10 25 26 22 21 20 19 18
```

**Ids 11, 17, 27 and 28 — `drain_life`, `darkness`, `curse`, `slow` — are absent**, so four spells
have no cell in this book. `main.res::text/spell.txt` is 28 lines keyed by id;
`main.res::text/spells.txt` is 24 lines keyed by slot and is the tooltip text.

The current spell and quick binding are zero-based cells. A normal book
command carries `cell+1`, not the intrinsic ID. Dispatcher005c2328-based
lookup translates that byte; statistics use the same table through base+4.
Thus cell5 travels as command6 and resolves spell23. The mapped ID indexes
the actor's sparse book. — AI-SPELLIDENT-286

Book availability is the OR of+18 over selected objects with nonnull+7c;
actual book-command members additionally need the chosen cell bit.
The shared Cast bit uses a narrower per-object `CUnit`, type17h/18h and
nonzero-mask predicate under the session gate. It is not primary-only and
does not require all members to cast. Upstream snapshot freshness and full
mixed-selection runtime reachability remain Unknown.
— AI-SPELLPOP-287, AI-SPELLCAP-288

Mouse selection checks availability before storing the cell; a nonnegative
shortcut stores it before later arming checks. The book's target-production
flag uses the cell-indexed005bcef8 table, not intrinsic IDs. Item context
selects a different target table and opcodes25h/26h: command+10 is then a
word-sized inventory position, and the item effect supplies the constructed
Spell's ID. This is a separate identity route, not another book-table lookup.
— AI-SPELLGUARD-289, AI-SPELLITEM-290

**The map picture is arithmetic.** No `Spells` column names it. A cast's picture id is

```
picture = 2 * spellId + 8        the homing / attached form
picture = 2 * spellId + 9        the burst / area form
```

and that value indexes `projectiles.reg` by its **`ID`** key. The parity is load-bearing: an even
picture creates no map object at all — the caster is looked up and given the cast action — while an
odd one creates a projectile. Two further message opcodes create one unconditionally. A picture id
with no `projectiles.reg` row is dropped silently at both the spawn and the draw, which is why
`light`, `invisibility`, `darkness`, `stone_curse`, `haste`, `control_spirit` and `slow` put nothing
on the map, and `fire_ball` and `poison_cloud` own two pictures each.

**What a consumer must not do:** treat the icon as indexed art (it is a strip and a mask); treat
`Delivery System` or `Distribution system` as the picture selector (they are not); or give a spell
an id-independent sprite (the formula pins the two together, and re-pointing one spell's art means
editing `projectiles.reg`'s `ID`).

## 11 — what a cast draws, and on which tick it applies

Settled by `MAGIC-CASTANIM-029`, `MAGIC-CASTTICK-030`, `MAGIC-BURST-031`, `MAGIC-SENDER-032`.
Section 10's parity rule is unchanged; this section says what each parity is *for*.

**Every cast animates the caster.** The cast routine sends the picture message at the **start** of
the wind-up, carrying `2*spellId + 8`, which is even for every id. The client's response to an even
picture is not "draw nothing": it finds the **caster** and sets

```
unit.actionCode   = 8        the cast action
unit.actionClock  = 0
unit.ticksLeft    = len(expanded Attack timeline)     units.reg, the same run a melee attack uses
```

skipping the whole branch if the unit is already animating or if its class has `AttackPhases == 0`.
The action then advances one frame per game tick with **no modulus** — frame index equals the tick
index — and ends when `ticksLeft` reaches 0.

**The two clocks.** They are separate numbers and the original does not reconcile them.

```
simulation   windUp   = 8 ticks   (actor field; an equipped item may override it)
             recovery = 4 ticks   (same)
             the effect is applied on tick windUp, counting the starting tick as 0;
             a melee strike lands on that same tick, in the other branch of one if

client       the visible projectile is spawned on tick ShootDelay   (units.reg)
             a melee hit reaction fires on tick AttackDelay         (units.reg)
```

A consumer that wants the swing and the application to coincide must choose one. The simulation's
tick is the one that changes hashed state.

**Which spells put a picture on the map.** The apply builds one of two effect classes, chosen by
the `Distribution system` column (section 5):

| Column value | Class | Sends a map picture |
|---|---|---|
| 1 | `PointEffect` | no |
| 3, 4, 5 | `AreaEffect` | yes — one `2*spellId + 9` message per cell, per ring |

Shipped `Data.bin`: **18 spells are `PointEffect` and 10 are `AreaEffect`** — `fire_ball`,
`wall_of_fire`, `fire_sacrifice`, `freezing_cloud`, `poison_cloud`, `acid_stream`, `light`,
`darkness`, `wall_of_earth`, `meteor_storm`. Of those ten, `light` and `darkness` compute an id
`projectiles.reg` does not define, so eight spells draw a spreading picture. The `AreaEffect` walks
a cell pattern with its own arms for `fire_sacrifice`, `acid_stream` and `meteor_storm` and a
default for the rest.

**What a consumer must not do:** read `MAGIC-PIC-027`'s even-parity rule as licence to draw nothing
on a cast. Every cast of every spell animates its caster; only the map object is parity-gated.

## 12 — the map object a cast makes, and how long it lives

Settled by `MAGIC-CASTSPAWN-033`, `MAGIC-BURSTLIFE-034`, `MAGIC-DELIVER-035`. Section 11
describes the message; this section describes the object. The two are built by different code at
different times.

**The cast object is built by the caster, not by the message.** The even picture puts the caster
into the cast action (section 11). On the tick `actionphase == ShootDelay` the caster's own
`vt+0x5c` runs and allocates one `0x14c`-byte projectile:

```
picture         = 2*spellId + 8          taken from the caster's actionspell field
position        = the class's muzzle point for facing (dir - 8) & 0xe,
                  or the class bounding-box centre when the class has no muzzle table
                  or the picture is 60 (teleport)
actionx/y/z     = the target's current position, or the caster's own aim point
action          = 1
actionphase     = 0
actionsegments  = a hard-coded switch on the picture id (below)
```

`picture == 60` allocates a **second** projectile, copy-constructed, at the caster's own centre —
teleport draws two sprites.

**The flight length is an engine table, not data.** 51 index bytes over pictures 10..60 through an
8-entry jump table. Seven ids are non-zero:

| picture | sheet | spell | actionsegments |
|---|---|---|---:|
| 10 | `firebolt` | Fire Arrow | `distance / 200` |
| 12 | `fireball` | Fire Ball | `distance / 384` |
| 20 | `healing` | Heal | 1 |
| 30 | `Drain` | Drain Life | 1 |
| 34 | `lightnin` | Lightning | 13 |
| 36 | `chain` | Prismatic Spray | 13 |
| 60 | `teleport` | Teleport | 21 |
| any other | — | — | 0 |

`actionsegments = 0` makes the driver return finished on its first tick, before its own picture
switch. So 21 of the 28 spells put an object on the map that never executes a driver arm.

**The burst object does not move.** The odd-picture client arm builds it at a cell and sets
`actionx/actiony` to its own `x/y` and `actiontarget` to 0, so every per-tick step divides zero. Its
lifetime is `msg+0xf`, chosen by the sender:

```
AreaEffect staged walker 16 ticks, raised to 18 on the acid_stream arm
                         a stage runs immediately, then every 3 ticks
the fire_ball sender     22 ticks
```

Two of the eight burst sheets are given `2 * Phases`, which is exactly one pass under the projectile
frame clock (`formats/anim`): `acid` has 9 phases and gets 18, `fireexpl` has 11 and gets 22. The
other six run 16 ticks against the 22 or 30 a full pass would need.

**A burst also plays a sound**, id `500 + picture`, at volume `(10000 - screenDistance)/100`,
skipped for picture 51 (`Meteor`).

**Two flight-length rules exist and only one is normally used.** The simulation computes its own
value in the cast routine — 0, or `distance / Data.bin parameter 7` when `Delivery System` is 2, or
5 for spell ids 13 and 14 — and puts it in `msg+0xf`. On a normal cast that value is discarded,
because the even picture routes the message to the caster-animation branch and the projectile is
built later by the table above. The simulation's value is used only when the source has no
client-side runtime id, which rewrites the opcode to `0x8b`.

`Delivery System == 2` holds for exactly `fire_arrow`, `fire_ball`, `lightning` and
`prismatic_spray`, which are exactly the four spells the client's own table gives a travelling or
ramped arm, and exactly the four picture ids with a special draw arm. A consumer may treat that
column as "this spell throws something", but the number it produces is not the number the original
draws with.

## 13 — staged area effects

### Cell aliases and direct-damage dispatch

The ring and single-pass blast read cell occupants in `+4`, `+8`, `+0xc` order and immediately
apply to each pointer. No object set suppresses a pointer encountered in another cell. Ring
coordinates pass the interior gate; blast Building coordinates are truncated independently to
bytes without that gate. The persistent cloud pulse reads only `+4`, not the Building slot
(`UNIT-AREAVISIT-071`).

Inner-effect polymorphism matters: base Effect `0059c688+3c` is the ordinary timed-effect attach,
but DirectDamage `0059c6e0+3c` is `FUN_00502c51` itself. The read fire-sacrifice, acid-stream and
meteor-storm producers carry DirectDamage. They therefore do not gain duplicate suppression from
the base-effect non-stacking rule. Fireball alone takes the explicit copy/divide branch of the
local dispatcher; it is not the only direct-damage area route (`UNIT-AREADIRECT-072`). This narrows
the old base-class-only and universal-idempotence clauses in `MAGIC-AREAAPPLY-038` and
`MAGIC-FIREDIV-047`; the unit n² storage and spell-2 normalization remain valid for that unit path.

For Building `0059c738`, the size getter is 1 regardless of rectangle size. Every reached direct
call uses the separate Building damage core, subtracts a word at `+42`, clamps a negative signed
result to zero, and notifies for positive damage. With aliases held present, repeated calls also
continue at HP zero. Actual HP-to-destruction timing is Unknown; do not infer either immediate
detachment or an immortal target (`UNIT-AREAHP-073`, `UNIT-STRUCTDETACH-074`).

Settled by `MAGIC-AREAPULSE-037` and corrected/completed by `MAGIC-RING-048`.
`Distribution system == 5` builds an `AreaEffect` with mode 2, stage byte `+0x4b = 0`,
timer word `+0x4c = 0` and orientation byte

```
orientation = (directionByte & 0xff) >> 5       0..7
```

The shipped rows selecting this mode are ids 4 `fire_sacrifice`, 9 `acid_stream` and 21
`meteor_storm` on both roots.

### Stage clock and cell consumer

The first stage runs on the first tick. A stage reloads the timer with 2. Each following tick
decrements the timer and returns while its old value is positive, so later stages run three ticks
apart:

```
stage ticks = 0, 3, 6, ...
```

The prior two-tick reading treated the reload value as the interval. For every coordinate entry the
consumer computes target cell plus offset and truncates x and y to bytes. The position helper accepts
the cell exactly when `8 <= x <= mapWidth-9` and `8 <= y <= mapHeight-9`. It clamps its stored x and y
to that range even when it rejects the cell. A rejected cell is neither sent nor applied. An accepted
cell receives one client message for picture `2*spellId+9`, then one application attempt for occupant
slots `+0x4`, `+0x8` and `+0xc`, in that order. The picture and simulation therefore use the same
accepted cell list. The sack slot `+0x10` is not read.

After the full list, the consumer increments the byte stage counter. It sets completion byte
`+0x40 = 1` when the incremented stage is greater than or equal to the arm limit. The common effect
driver reaps it on that tick. Staged effects register no map layer and perform no separate cleanup.

### Fire Sacrifice

Fire Sacrifice ignores the orientation. Its two stages are fixed ordered lists relative to the
target cell:

```
stage 0: (-1, 1) (-1, 0) (-1,-1) ( 0, 1) ( 0,-1) ( 1, 1) ( 1, 0) ( 1,-1)
stage 1: (-2, 1) (-2, 0) (-2,-1) (-1, 2) ( 0, 2) ( 1, 2)
         (-1,-2) ( 0,-2) ( 1,-2) ( 2, 1) ( 2, 0) ( 2,-1)
```

The centre and the four corners of the radius-2 square are not visited.

### Acid Stream

Acid Stream has six stages. Define the two ordered source tables:

```
A_s = [(x,s) for x=-s..s]             s = 0..4
A_5 = []
B_s = [(s-i,i) for i=0..s]            s = 0..5
```

Even orientations use A and odd orientations use B:

| orientation | transformed offset |
|---:|---|
| 0 | `(+dx,-dy)` from A |
| 1 | `(+dx,-dy)` from B |
| 2 | `(+dy,+dx)` from A |
| 3 | `(+dx,+dy)` from B |
| 4 | `(-dx,+dy)` from A |
| 5 | `(-dx,+dy)` from B |
| 6 | `(-dy,+dx)` from A |
| 7 | `(-dx,-dy)` from B |

Stage 5 is empty for even orientations and contains six cells for odd orientations. The client
burst lifetime is 18 ticks rather than the default 16.

### Meteor Storm

Meteor Storm has 32 one-cell stages. Each stage calls `rand(5)` twice, in x-then-y order:

```
dx = first rand(5) - 2
dy = second rand(5) - 2
```

Each coordinate is in `[-2,3]`. The orientation and the `Radius` column are not used. Cells can
repeat across stages. A sampled cell is accepted only within the eight-cell interior:
`8 <= x <= mapWidth-9` and `8 <= y <= mapHeight-9`. A rejected edge sample is clamped internally
but sends, applies and draws nothing. An accepted stage applies to occupants and sends picture 51
at that cell with lifetime 16. The client consumes no further placement RNG.

Picture 51 is stationary at the accepted cell. Its sixteen paced client ticks draw frame 8 at
vertical offsets `-28,-24,-20,-16,-12,-8,-4,0`, then frames 0 through 7 at zero offset. It ends on
frame 7 and does not wrap (`MAGIC-093`, `ANIM-046`).

### Table layout and customisation limits

The two fixed arms use stage records in `ROM.EXE`:

```
record size 0xa4
+0x00  i32 count
+0x04  i32 dx[20]
+0x54  i32 dy[20]
```

Fire Sacrifice and Acid Stream take terminal counts 2 and 6 from image dwords. Meteor Storm takes
the immediate 32 and overwrites one count-1 record. The 20-cell capacity, stage bounds,
8-orientation transform, random domain, byte stage counter, byte map coordinates and word timer are
engine limits. Changing the shipped fixed geometry changes `ROM.EXE`, not `Data.bin`. An engine may
externalize replacement tables without changing any shipped file bytes. More than 20 cells, 255
stages, 8 directions or byte coordinates requires a wider engine representation.

## 14 — the drawn path of a travelling projectile

Settled by `MAGIC-BOLTGATE-069`, `MAGIC-BOLTSHAPE-070`, `MAGIC-BOLTLIST-071`,
`MAGIC-BOLTSTILL-072`, `MAGIC-TRAIL-073`, `MAGIC-BOLTEND-074`. Section 12 gives the
object and its lifetime; this section gives what is drawn between the two ends of a flight, which
is per-picture and is not the object's own sprite for two of the seven ids.

**Two picture ids draw a path and no other does.** The projectile calls its own `vt+0x50` once per
tick, on the path it takes when `actionsegments` is still non-zero and its `action` byte is 1, the
value the cast spawner writes. That routine's whole body is behind two comparisons on
the picture: 34 (`lightnin`, Lightning) and 36 (`chain`, Prismatic Spray). Every other picture
returns having touched nothing, and the geometry generator is reachable from nowhere else in the
image.

**A path is a list of 8-byte records held on the object.** The list is a `CArray` at object
`+0x110`: data pointer `+0x114`, count `+0x118`, capacity `+0x11c`.

```
+0x00  i16 x        screen coordinate
+0x02  i16 y        screen coordinate
+0x04  i16 zero     no reader located
+0x06  u8  tag      34 for picture 34; the chain-link index mod 7 for picture 36
+0x07  u8  zero
```

**The geometry is a bounded random walk.** Given the two endpoints, the generator takes their
distance, builds the shape on a horizontal segment of that length, then rotates every point onto
the real direction. The walk itself:

```
deflection step   (rand() % 7) * randomSign,  applied only when its magnitude is >= 2
abscissa step     (rand() % 50) * 0.01,       applied only when it is >= 0.15
deflection sum    clamped to [-3, +3]
step scale        multiplied by -1 each accepted step
walk ends         when the abscissa passes 0.7
smoothing         one pass over every second point, factor -0.5
acceptance        every point's distance from the straight line must be < 0.15 * length;
                  on failure the whole walk is generated again with the same arguments
```

So the drawn figure is ragged, different on every call, and confined to a lens of half-width
`0.15 * length` about the straight caster-to-target line. None of these numbers is data: the ten
constants are `.rdata` doubles and the two moduli are immediates.

**The list is rebuilt every tick and does not accumulate.**

```
picture 34   resolve actiontarget; generate one set with tag 34;
             SetSize(generatedCount) and overwrite  -> the list is REPLACED
             an unresolvable target leaves the previous list in place
picture 36   SetSize(0)  -> the list is CLEARED
             then, per entry of the projectile's own word array at +0xac (count +0xb0):
               resolve the id, generate one set with tag = (index mod 7), append
```

The generator always spans the whole source-to-target segment, and no argument carries a partial
length, so the whole figure exists from the first tick. The `+0xac` array is copied word by word
onto the projectile by the cast spawner from the caster's own array of the same offsets.

`MAGIC-SPRAY-134`…`MAGIC-SPRAY-137` close that array's producer. After applying the selected list, simulation opcode `0x8a`
or `0x8c` sends `count+1` words: caster identity or packed source cell first, then every selected
victim id in list order. The matching client arm clears the source's embedded word array, appends
those ids, writes its count and primary, and only then reaches the spawner copy above. The previously
proposed `004162d4` `SetSize` site belongs to a different list operation.

**These two projectiles do not move.** The driver computes a per-axis step before its switch, but
only the default arm applies it. Pictures 34 and 36 take an arm that writes the sheet frame and
nothing else, so `x`, `y` and `z` keep the spawner's values for the object's whole life and
`actionsegments = 13` is a countdown of ticks, not a division of a distance. A moving target is
followed by regenerating the figure onto its new position, not by travelling toward it.

**Nothing is created at the end.** On the tick `actionsegments` is already 0 the driver returns
before its switch and before the `vt+0x50` call; that path allocates nothing and sends nothing.
The list is a member of the object and ceases with it. Neither spell has a second picture:
`2*spellId + 9` gives 35 and 37, and no `projectiles.reg` row defines either on either root.

**The Fire Arrow and Fire Ball trail is the opposite mechanism.** Those two pictures move normally
and keep a `CObArray` at object `+0x138` (data `+0x13c`, count `+0x140`) of **past** positions,
each packed as `(y << 16) or x` and appended after the object has been moved for that tick. When
the count reaches **6** the oldest entry is removed first. So that trail accumulates, is bounded at
six, and expires one entry at a time — where a bolt's list is wholesale and vanishes with the
object.

**Customisation limits.** The two picture ids, both comparisons, the six draw arms, the ten shape
constants, the two moduli, the tag stride 5, the 8-pixel centring, the trail cap 6 and the 13-step
frame ramp are all `.text` or `.rdata`. A third spell with a drawn path, a straight bolt, a longer
trail or a differently sized bolt sheet all require an engine change; none is reachable from
`Data.bin` or `projectiles.reg`.
## 15 — what a lasting effect draws on its actor

Promoted claims: `MAGIC-MARK-059`,
`MAGIC-MARK-060`, `MAGIC-MARK-061`, `MAGIC-PROT-062`, `MAGIC-SHIELD-063` (superseded),
`MAGIC-BLESS-064`, `MAGIC-CLOUD-065` (partially retracted), `MAGIC-ACTOR-066`.

Corrected and completed by the promoted presentation claims
`MAGIC-089` through `MAGIC-093`, `ANIM-044` through `ANIM-047` and `SPR16A-031`.

### The two states an actor's marks are held in

The presentation side of an actor holds two arrays.

```
effect list     unit+0x124 array object, unit+0x128 data, unit+0x12c count
                one dword per element:  bits 31..16 kind, bits 15..0 countdown
mark array      unit+0x110 array object, unit+0x114 data, unit+0x118 count
                8 bytes per record:
                  +0x00 i16 dx      horizontal offset
                  +0x02 i16 dy      vertical offset
                  +0x04 i16 depth   vertical offset, and selects the draw pass
                  +0x06 u8  record  index into the projectile record array
                  +0x07 u8  phase   frame index for the blit
```

A `kind` is `2*spellId + 8`, the even half of section 10's picture pair.

### Lifecycle

The simulation opens and closes each element with a dedicated message. On attach, the effect's spell
id sets a bit in `actor+0x144` and, when that spell id is non-zero, an attach message is sent
carrying the target's runtime id and the effect's picture field. On expiry the bit is cleared and a
detach message is sent. A recast of the same spell evicts the old effect first, so it sends detach
then attach.

The client creates the element with countdown `0xffff` on attach, replacing any element of the same
kind, and removes it on detach. The countdown falls by one on every rebuild and the element is
removed when it reaches zero, but at `0xffff` that does not happen in play: the countdown is the
phase clock, and the element's lifetime is the effect's.

A third writer exists. The projectile driver dispatches an arriving projectile's picture over the
closed range 13..64, and two of its arms write the same list. Pictures 20 `healing` and 30 `Drain`
replace or append an element with countdown **32**, which does expire on its own. Pictures 18, 24,
28, 40, 44, 48, 52, 54, 56, 62 and 64 append one with countdown `0xffff` and no by-kind check, but
section 12's flight-length switch gives all eleven a length of 0, so on the cast path the driver
returns before reaching that arm.

`heal` and `drain_life` therefore run from a 32-rebuild element rather than from an attach message.
Their builders use countdown only as the unsigned fresh-spawn gate; their motion magnitude comes
from the separate constant argument `32T` (`MAGIC-089`, `MAGIC-090`).

### Rebuild

Every rebuild discards the mark array and re-derives it. For each element the kind selects an arm;
kinds outside `0x12..0x3e` and 35 of the 45 kinds inside it select an arm that appends nothing. Ten
kinds reach seven builders:

| kind | spell | records appended | phase input |
|---|---|---|---|
| 0x12 0x1c 0x28 0x34 | the four Protections | 1 each | countdown mod 6 |
| 0x14 | heal | carry own phase below 7; append 3..5 while countdown >= 8 | `32T`, countdown gate |
| 0x18 | poison_cloud | 1 | countdown mod 6 |
| 0x1e | drain_life | carry own phase below 7; append 3..5 while countdown >= 8 | `32T`, countdown gate |
| 0x2c | shield | two independently rasterised components | countdown mod 90 |
| 0x36 | bless | 20 | countdown mod 5 |
| 0x3e | curse | 20 | countdown mod 5 |

`haste`, `slow`, `stone_curse`, `invisibility`, `freezing_cloud`, `light`, `darkness`, `lightning`,
`prismatic_spray`, `acid_stream`, `wall_of_earth`, `meteor_storm`, `control_spirit`, `teleport`,
`fire_arrow`, `fire_ball`, `wall_of_fire` and `fire_sacrifice` append nothing. Each element is
dispatched independently, so an actor carrying several effects shows the union of their records.

Every builder scales its lengths by `T`, the actor class's `units.reg` `TileSize`, default 1.

### Drawing

The actor draw walks the mark array twice: once before the actor's own sprite, drawing only records
whose `depth` is positive, and once after, drawing only records whose `depth` is not positive. Both
compute

```
x = anchorX - Width/2  + dx
y = anchorY - Height/2 + dy - depth
```

where `Width` and `Height` are the projectile record's centring halves (section 10) and the anchor
is the actor's own draw position. Both blit the sheet at `record` with frame `phase`.

The `depth` field is therefore both a vertical offset and a z-order selector, which is how a mark
set can enclose the actor: records with a positive vertical term pass behind the sprite and records
with a non-positive one draw in front of it.

### The four Protections

One builder serves all four and appends exactly one record with `depth = 0`, `record = kind`, and

```
protection_from_fire   dx = -6  dy = -(T*32)
protection_from_water  dx = +6  dy = -(T*32)
protection_from_air    dx =  0  dy = -(T*32) - 6
protection_from_earth  dx =  0  dy = -(T*32) + 6
```

Four simultaneous Protections therefore stand at the four corners of a diamond of half-width 6
pixels, centred `T*32` pixels above the actor anchor. Nothing in the builder recognises that case;
it is the union of four independent placements.

### shield, bless and curse

`shield` concatenates two independently rasterised kind-44 components, A then B. The common midpoint
routine starts `x=0,y=r,err=2-2r`, emits sampled symmetry points in order
`(+x,+y),(-x,+y),(+x,-y),(-x,-y)`, including axis duplicates, and emits nothing at radius zero.

For `p=countdown%90`, Component A rasterises radius `16T` at interval 4 with `theta=4p` degrees:

```
q = trunc(y*28T/16T)
f = abs(4-trunc(abs(y)*5/16T))
A1 = (trunc(x*cos(theta)), -q, 11T+trunc(x*sin(theta)), phase f)
A2 = (-trunc(x*sin(theta)), -q, 11T+trunc(x*cos(theta)), phase f)
```

Component B computes

```
c45 = binary64(bits 0x3f96c16c16c16c17)
x1 = round53(p*c45)
x2 = round53(x1-1)
x3 = abs(x2)
x4 = round53(x3*float32(28T))
E = trunc(x4)
F = 4-floor(abs(p-45)/9)
R = trunc(sin(acos(E/(28T)))*16T)
```

and rasterises `R` at interval `2F`. Each `(u,v)` gives `(u,-(E+11T),v,phase F)` then
`(u,E-11T,v,phase F)`. At the wrapper's post-prologue stack base `P`, `[P+0x30]` is the `16T`
argument and `[P+0x34]` is overwritten with `E`; after pushing the interval, `FMUL [ESP+0x34]`
therefore reads `16T`, while `E` is at `[ESP+0x38]`. The envelope sequence is integer load, multiply
by the stored qword nearest `1/45`, subtract one, absolute value, multiply by the float32 `28T`
argument, then x87 truncation. CRT startup requests PC53, and a debugger observation after the
initializer returns reads control word `0x027f` on both lawful roots. The raw and recursive census
classifies all 127 FLDCW, FLDENV, FRSTOR, FSAVE/FNSAVE, FNINIT and FXRSTOR candidates. The only two
valid FSAVE instructions reset the environment temporarily, call one helper, then immediately
restore the saved environment with FRSTOR before any branch or return. No lasting PC64 writer is
present on the normal Shield path. PC53 round-to-nearest therefore applies after each arithmetic
instruction; `__ftol` changes the rounding-control bits temporarily and restores the saved word.
The stored qword is exactly `6405119470038039/2^58`, but retaining that rational without the PC=53
operation boundaries is a different and wrong model. Shipped `T=3` has `E=56` at `p=15`, `E=27` at
`p=30` and `E=28` at `p=60`. Exact-dyadic retention disagrees at `p=15`; division by 45 disagrees at
`p=30` and `p=60`. Round-down differs at phases 30, 60 and 75; round-toward-zero differs at phases
15, 30, 60 and 75. Round-up agrees with round-to-nearest on the 270 integer outputs, so the runtime
control word is the discriminator. At `p=0`, `F=-1` but `R=0`, so no invalid frame is emitted. That is the only
zero-output phase for every shipped `T={1,2,3}`; at `p=45`, `E=0` and `R=16T`. Every emitted B
record uses a frame in 0..4. The builder reads neither RNG nor previous marks (`MAGIC-091`,
`MAGIC-092`).

`bless` and `curse` each append 20 records: five trail steps of four records, on a circle of radius
20.0 centred `T*32` above the anchor. Angles are in degrees. `bless` steps 89, 71, 53, 35, 17 and
`curse` steps 0, 18, 36, 54, 72, so the two trails sweep opposite ways as the countdown falls. The
frame index of a step is `4 - step/18`, so five different frames are visible at once. A sine term
goes into `depth` rather than into `dy`.

### poison_cloud, heal and drain_life

`poison_cloud` appends one record at `dx = 0`, `dy = -(T*32)`, `depth = 0`, `record = 0x18`.

`heal` and `drain_life` first walk the previous mark array and carry their own records whose phase is
below 7. They increment the phase, then Heal subtracts and Drain adds
`s=trunc_toward_zero(1+32T/7)` to `dy`; shipped `T={1,2,3}` gives `s={5,10,14}`. The step is
constant throughout the effect. A phase-6 record is copied as visible phase 7, and phase 7 is then
dropped, so a cohort has eight outputs.

While unsigned countdown is at least 8, each builder appends `rand()%3+3` fresh particles. Every
particle consumes one further `rand()%360` angle `a`:

```
heal   dx=trunc(cos(a)*16T), dy=0,    depth=trunc(sin(a)*8T), record=0x14, phase=0
drain  dx=trunc(cos(a)*16T), dy=-32T, depth=trunc(sin(a)*8T), record=0x1e, phase=0
```

Countdown 32 through 8 produces 25 cohorts and 75–125 particles in total. Eight cohorts overlap in
steady state, 24–40 particles. One angle supplies both ellipse coordinates; there are not two
independent position samples (`MAGIC-089`, `MAGIC-090`).

### Two spells that change the actor instead

`stone_curse` and `invisibility` reach no mark arm. The actor draw looks their kinds up directly.

`0x30` gates two separate arms in `FUN_0045b3f0`, each additionally conditional on
`byte [drawable+0x15a] <= 2`. The first, at `0045b8a2`, discards the frame index the per-state arm
computed and substitutes `(drawable+0x6c - 8) & 0xf`, the facing at full 16-way resolution, mirrored
above 8. The frame therefore stops following the animation clock and changes only when the unit
turns. The second, at `0045b98b`, replaces the shade table the blit takes: element 16 of the array
at `[0x005eb65c]` for a class whose `Palette` is 0, otherwise a table built per draw from the
class's own palette with `FUN_00427df0(palette, 0x10, mode 5, 0)`, mode 5 being the greyscale ramp
`(r+g+b) * level * 2 / 3 / levels`. Both forms draw the unit grey. `FUN_0045bf00` carries the frame
arm a second time at `0045c29a` (`MAGIC-STONEDRAW-084`).

`0x26` gates the actor's whole sprite on a per-player bit. That bit is bit 3 of the 32-entry
`CWordArray` row at `CPlayer+0x34` on the local player, indexed by the drawn unit's owner index, and
it is set on exactly one entry — the local player's own (`UNIT-VISBIT-044`). The test is therefore
an ownership test: an invisible unit is drawn to its owner through a different blit slot and is not
drawn at all to any other client (`MAGIC-INVISOWN-085`). The map draw applies the same test before
drawing a unit at all.

The by-kind lookup `FUN_004599d0` has fourteen call sites over six routines and no consumer outside
the presentation layer; three of the six are consecutive vtable slots of the client unit class,
`FUN_0045b3f0`, `FUN_0045bf00` and `FUN_0045c940` (`MAGIC-EFFLOOK-083`).

### Customisation limits

The kind range, the dispatch table, the seven builders, the 6-pixel Protection radius, the four-way
Protection assignment, the `T*32`, `T*28`, `T*16`, `T*11` multipliers, the 20-mark trail, the
18-degree step, the 20.0 radius, the 90-step and 6-step and 5-step moduli, the eight-output particle
life and the 8-byte record are all engine code or `.text` immediates. Changing any of them changes
`ROM.EXE`. `TileSize`, the sheets behind each record index and their `Phases`, `Width` and `Height`
are data and can be replaced without touching the executable. A sheet with more phases than its arm
indexes gains nothing: `poison_d` declares 8 phases against 6 the engine can index, and `Drain`
declares 9 against 8.

## 16 — what a lasting effect does to the actor's actions

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
