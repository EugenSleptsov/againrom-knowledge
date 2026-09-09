# Spell objects and casting

[Reference](format.md)

## Arithmetic

`ftol` truncates toward zero (`0055458c`, rounding-control 11); `IDIV` truncates toward zero.
`rand(n)` is `FUN_00504003`: `rand() * (n+1) >> 15`, uniform on `[0, n]` **inclusive**, 0 when
`n == 0`. Runtime `Data.bin` column *i* is title *i+1* throughout.

## Spell definitions

Direct-damage effect delivery to a Building differs from a physical actor action:
`FUN_00502c51` selects the Building resolver, subtracts from word `+0x42`, clamps a
negative result to zero and notifies only for positive damage. The physical Building
strike does not clamp and has a different notification predicate. The shared resolver
does not make the two delivery paths interchangeable (`UNIT-STRUCTDELIVERY-065`). This
does not establish every spell's Building admission or area-effect multiplicity.

`World\Data\Data.bin` group **H spells** ships **28** rows, indices 1..28. `rom.exe` carries the
same names as a 29-entry array at `0x005c5b58` whose entry 0 is `unused_spell_0`; the names agree
and the ID is the array index. Spell 0 is refused at every entrance.

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

## Spell objects and spellbooks

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

## Spell power

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
`skill/5` damage-floor term for a spell** — the weapon-skill terms of [HERO](../hero/format.md) step 11 have no
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
([HERO](../hero/format.md) step 0), so **without a cap-raising effect Mind contributes at most 20 points of
power**.

At the ceiling, in the integers the engine stores: `f` tops at **4.333**, so damage is 4.25×–4.33×
the `Data.bin` columns; a duration column multiplies by `1.025^100 = 11.81`; range gains at most
**3** cells and Teleport **33**; Prismatic Spray reaches **7 total victims**. Nothing overflows its store on the
shipped table — the largest damage base is 43 and spread 87 against byte fields, the largest
duration 6 312 against a `u16`.

## Cast sequence

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

## AI cast selection

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

## Mind and Spirit

Enumerated image-wide (`EnumRefs disp:88`: 461 hits / 235 owners / 4 orphan, all four stack frames;
`disp:8a`: 36 / 13 / 0), with the actor's `u16` width used as the filter and every byte- and
word-wide hit read:

```
Mind    -> the power term, one for one with the school skill    (the three spell-power sites above)
        -> Prismatic Spray's victim cap, through the one unclamped helper
        -> whether an AI mage casts at all                      (AI cast selection above)
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
