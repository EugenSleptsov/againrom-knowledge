# HERO — character generation and the derived-stat graph

Level 3 spec. Promoted claim basis: `HERO-STAT-001`, `HERO-COST-002`,
`HERO-BUDGET-004`, `HERO-EQUIP-017`, `HERO-MOD-016`, `MAGIC-SPELL-001`,
`MAGIC-CAST-003`, `HERO-CADENCE-112`…`HERO-CADENCE-115`.
Not a file format — the simulation area that turns four player-chosen numbers into a hero. Everything
here was read out of `rom.exe` at instruction level; the addresses and evidence bounds are in the
promoted claim rows cited by each section. The spell side is
[`formats/magic`](../magic/format.md).

## Integer conventions — read this before any formula below

City sale reaches full derive only if `004f36f7` finds a changed truncating
signed16 load/capacity quotient. It uses stored prior load, updates the load
word with wrapping arithmetic, and has no zero-divisor guard (`SAV-CITYSALE-513`).
A successful school purchase instead ends with unconditional Human derive.
Derive rebuilds live fields and clamps current health/mana but does not rebuild
the modifier block from equipment: it consumes retained modifiers, with the
conditional negative-speed reset of `+d8`. Its direct stores leave base and XP
blocks intact; school updates selected skill/base/XP before the derive call.
This bounded producer result does not establish runtime callback chronology or
safe initial values (`SAV-CITYDERIVE-515`).

- `ftol` is `__ftol` (`0055458c`). It forces x87 rounding-control to **11 = truncate toward zero**
  before `FISTP qword`, then returns the 64-bit integer in `EDX:EAX`. **It truncates. It does not
  round.** A caller that stores only EAX narrows that result to the low 32 bits.
- `/` between integers is `IDIV` or the `CDQ`/`SUB`/`SAR 1` idiom: also **truncation toward zero**.
- The **only** rounding-to-nearest in this whole area is the `+ 0.5` inside `T(n)` below.
- Every intermediate of the health and mana chains is stored back into a **16-bit** field, so the
  truncations compound; a consumer that carries doubles through will drift.
- `pow(1.1, s)` is `FUN_004f7233` and `log₁.₁(x)` is `FUN_004f7253`. The `1.1` is pushed as `pow`'s
  **first** argument, so the stat is the **exponent**.

## 1 — the four primary stats

| name on screen (RU) | engine name | actor | record | chargen panel row |
|---|---|---|---|---|
| Сила | `Body` | `+0x84` u16 | `+0x138` u8 | 0 |
| Ловкость | `Reaction` | `+0x86` u16 | `+0x13b` u8 | 1 |
| Разум | `Mind` | `+0x88` u16 | `+0x139` u8 | 2 |
| Дух | `Spirit` | `+0x8a` u16 | `+0x13a` u8 | 3 |

The engine names are `scenario\npc.reg` key literals, read by the pre-create loader in that order.
The record order is **not** the panel order; the panel reads `0x138, 0x13b, 0x139, 0x13a`.

The same crossing appears in three further places, each an independent witness of the binding: the
effect dispatch's kind table (`3 mind → +0x88`, `4 reaction → +0x86`), the object serialiser, which
sends the four **as single bytes** and reads them back from stream offsets `+0x28 → Body`,
`+0x2b → Reaction`, `+0x29 → Mind`, `+0x2a → Spirit` (skipping any that arrives 0), and the
record writer, which copies `+0x84→+0x138`, `+0x88→+0x139`, `+0x86→+0x13b`, `+0x8a→+0x13a`. **The
wire is byte-wide**, which is a ceiling on a stat independent of the cap below.

Every reader of the two magic-relevant stats is enumerated in `formats/magic` section 8. Outside
magic: Mind buys sight and multiplies experience gain; Spirit buys the mana maximum and the five
elemental protections. Neither is read by movement, by the hit resolver or by regeneration.

## 2 — character generation

```
T(n) = ftol( 0.349 * pow(1.15, n - 1) + 0.5 )        the CUMULATIVE cost of one stat at n

start        all four stats = 25, pool = 100
"+" on v     refuse if pool < T(v+1) - T(v)  or  v >= 45
             else  v += 1 ;  pool -= T(v_old+1) - T(v_old)
"-" on v     refuse if v <= 15
             else  v -= 1 ;  pool += T(v_old) - T(v_old-1)
budget       accepted iff  T(body) + T(reaction) + T(mind) + T(spirit) <= 140
             otherwise all four are forced back to 25
identity     pool == 140 - sum T(stat)   (because 4*T(25) + 100 == 140)
```

The refund is **exactly symmetric**: `refund(v) == cost(v-1)` for every `v` in `[16..45]`.
The counter drawn on the creation panel is the **remaining** pool.

Reachability, given the budget: one stat alone reaches **42** with the others left at 25 and **43**
with the others floored at 15; all four together reach **34**. **The click bound 45 cannot be reached
in character generation** (`T(45) = 164`; `HERO-BUDGET-004`).

Nothing in the point-buy depends on class, race, or the other three stats, and no registry key or
class table participates. Character generation also sets exactly **one** skill slot — to `20` when the
global `[0x005eb5a4] == 2`, otherwise `10` — zeroing slots 1..5 first.

### What character generation puts in his hand

`FUN_004f992c(slot, value)` zeroes skill slots 1..5, writes the one chosen slot, recomputes, then
`CMP ECX,0xa` (`004f9a00`) selects one of two sets of five weapons, built from **string literals**
through `FUN_0050d670`. `value` is 20 only when `[0x005eb5a4] == 2` **and** `campaign+0x6bc == 2`
(`0047c191`…`0047c1bf`); otherwise 10 (`HERO-START-039`).

```
slot  skill      value <= 10 (ordinary)          value > 10
1     Blade      Iron Short Sword                Uncommon Steel Two Handed Sword
2     Axe        Uncommon Bronze Axe             Uncommon Steel Axe
3     Bludgen    Uncommon Bronze Mace            Uncommon Steel Mace
4     Pike       Bronze Pike                     Uncommon Steel Pike
5     Shooting   Uncommon Wood Short Bow         Uncommon Magic Wood Short Bow
mage  (staff)    Wood Staff {castSpell=Fire_Arrow:10}   ... :20
```

No instruction anywhere in the mapped image names the `Data.bin` Weapons row `BareHands`, so an
unarmed hero has an all-zero modifier and `active = 0` — both skill terms are skipped and the roll
is `[d, 2d]` with `d = ftol(1.1^body/20)`, which is 0 below body 32 (`HERO-BARE-037`).

## 3 — the six skill slots

`actor+0xa8 + 2i`, `i = 0..5`. Slot 0 is `Skill.General`; slots 1..5 are one shared set that the class
renames, `Data.bin`'s own column titles giving both names at once:

| slot | fighter | mage |
|---|---|---|
| 1 | Blade | Fire |
| 2 | Axe | Water |
| 3 | Bludgeon | Air |
| 4 | Pike | Earth |
| 5 | Shooting | Astral |

Base copies live at `+0x116 + 2i` (`i = 1..5`) and bonuses at `+0xe8 + 2i`. There is no
sixth..tenth slot: `+0xb4` immediately follows and is written as the damage byte.

The word at `+0x116` is inside the serialized block and both school price arms index it when their
slot is 0, but it is **not** a General base copy: all twelve restore, snapshot and bonus-fold loops
start at 1 and end at 5. Five are in the play-award routine, two each in purchase and derive, and
one each in load, creation and death loss. A General purchase changes live `+0xa8` and `xp[0]`
without refreshing `+0x116`.

Slot 0 is the ranged-weapon accuracy skill. A ranged weapon copies `skill[0]` to the to-hit
modifier and sets active skill to 0. Melee weapons select one of slots 1..5 instead. General
experience also enters the six-slot total used by the health and mana maximum formulas. It is not
included in the five-slot victim XP-value sum.

The spell Sphere is a raw index into this vector on three readers and one writer. The writer's
accessor `FUN_004fe11e` takes the Spell row at `+0x04`, pushes parameter index 2, calls the row
accessor and returns byte 0; the paired `Spells` title is `Sphere`. Its caller masks the result to a
byte before indexing the vector. Both shipped roots have 0 Sphere-0 rows among 28 spells; custom
Sphere 0 therefore activates an existing General seam without changing the spell schema.

The bonus slots are written by **effects only**, and by two class-gated sets of six arms over the
same six fields: `fighterSkill0`/`skillBlade…skillShooting` behind the fighter predicate,
`mageSkill0`/`skillFire…skillAstral` behind the mage one. An effect naming a fighter skill does
nothing to a mage and vice versa. Equipment reaches them only through an item's effect list.
For a recomputing human, the two slot-0 effect kinds are inert: they write bonus `+0xe8`, while the
derive restores and folds bonuses only for indices 1..5.

## 4 — experience, and how a skill level moves

Two fields per slot, and they are **not** two views of one number: the level `+0xa8 + 2i` and the
experience `+0x1cc + 4i`. `actor+0x130` is their running sum.

```
S(n)   = ftol( (pow(1.1, n) - 1) * 1000 )       what a slot at level n accounts for
S^-1(x)= ftol( log_1.1( x/1000 + 1 ) )          used by the loss ONLY
invariant  actor+0x130 == sum over i = 0..5 of xp[i]      preserved by all four writers
at creation  xp[i] = S(skill[i])                          FUN_004f7c28, both build paths
```

### the raise — `FUN_004f72d7`, `vt+0x5c` of both human classes

```
gate    typeID (+0x0e) in [0x21, 0x3f]             player-character mode, not every Human
        source, if given and not already dead (health >= 0):
          same owner (+0x14)                        -> nothing
          diploMatrix[srcPlayer][myPlayer] & 2       -> nothing
amount  ftol( amount * (mind/30 + 0.25) )           signed; zero is not an early refusal
slot    +0x4c & 4 ? the caller's slot, which must be > 0
                  : the caller's slot must be 0, and the slot used is +0xb6 (active), > 0
        skill[slot] must be < 100, tested BEFORE the award
cap     amount = min(amount, S(level+1) - S(level))         negative amounts remain negative
                                                          [0x005cd758]+0x0c ? that / 5 : that
pay     xp[slot] += amount ; actor+0x130 += amount          can reduce XP; no level-decrement arm
raise   if xp[slot] > S(skill[slot]):  skill[slot] += 1     ONE level, never two
after   a raise re-snapshots base[1..5] and calls vt+0x50 (the whole derive);
        on the +0x4c & 4 arm it also re-powers every spell in the book,
        power = clamp(skill[Sphere] + Mind - 30, 0, 100)
no raise -> skill[i] = min(base[i] + bonus[i], 100) for i = 1..5
```

Three feeds, all `vt` slots of the same two vtables and all **stubs** on the non-human one:
`vt+0x60` a kill, `ftol(victim+0x1c * 0.5)`; `vt+0x64` a landed hit,
`ftol(victim+0x1c * 0.5 * dealt / victim.healthMax + 1)`; `vt+0x68` a cast,
`round(manaCost/2)`. `victim+0x1c` is the `Units` `XPvalue` column. Each resolves the slot as
*"the spell's `Sphere` if a spell id was supplied and `+0x4c & 4` is set, otherwise 0"*. The
clear-bit arm accepts that zero and substitutes the current weapon skill at `+0xb6`; a fighter's
weapon-borne spell damage therefore trains its weapon, not a school, only after the typeID gate
above. A low-type map Human receives no item hit, damage or kill award. Item context suppresses the
`vt+0x68` cast feed, but later direct damage, Drain Life, Slow, Stone Curse, Curse and Poison Cloud
can each reach `vt+0x64`. On the fighter route those spell-side events run before the physical XP
test. That later test rereads post-spell health and refuses at `-10` or below, so the physical blow
does not necessarily add its own award. Their cadence follows affected targets and ticks.

The damage feed first requires a victim owner with `+0x28 != 0`; kill has no `+0x28` test. Both
require a victim owner and refuse when owner `+0x5c` and multiplayer are both nonzero. Poison ticks
reject only a zero signed amount and retain a recorded caster at health zero, clearing it only below
zero. A custom negative item power can therefore turn Poison into a healing tick that still reduces
the credited Humanoid's slot and aggregate XP. The shipped Catapult and Ballista carry castSpell
weapons, but use the non-human vtable's three stubs: they can release the rider and never train.

The kill feed's second argument is the final victim attribution byte. The resolver retains prior
state for null source/definition, clears the actor for an ownerless defined source, and otherwise
writes source plus damage kind for bit 4 or zero for clear bit. A non-Defensive PointEffect can
replace it with actual id after payload; an AreaEffect has a distinct writer gated by a nonzero,
zero-extended target movement-domain byte.
Lethal payload damage does not change either condition for ordinary shipped direct-damage shapes.
Drain Life constructs neither envelope and writes no attribution. Poison Cloud can seed id 8 on
area application but does not refresh attribution on later ticks, leaving both eventual kill paths
dependent on the victim's runtime history.

Consequences a consumer must carry:

- **No play award raises `Skill.General`.** Both arms exclude slot 0, and every restore /
  re-snapshot loop runs `i = 1..5`. Opcode 61 can raise it when supplied slot 0, although the
  shipped school cannot supply that selection.
- **A `+0x4c & 4` carrier earns nothing from melee or missile combat** — those arrive with slot 0
  and its arm refuses that.
- **The first award into a slot always raises it**, because creation leaves `xp[i] == S(skill[i])`
  and the test is a strict `>`.
- **A slot at 100 accrues nothing at all** — the ceiling test precedes the payment.
- **Item-spell progress uses these same six slots and fields.** There is no item proficiency or
  separate save state.

### the loss — `FUN_004f7a53`, charged by conditional repair

For `i = 0..5` (**General included**): `xp[i] = xp[i] * 9 / 10`, then `skill[i] = S^-1(xp[i])`; the
total is rebuilt and `base[1..5]` re-snapshotted. Its one caller `FUN_004d3755` runs it only for
an existing primary with **zero-extended byte** `actor+0x13c != 0`, then restores HP and MP.
The loss tail clears the stage. This is not a campaign death or reload charge.

There are three caller sites of that accepter (`HERO-DEFEAT-136`):

| Entry | Admission to repair |
|---|---|
| Reporter | Actual `server+0x0c != 0`, participant active, latch at least 2, entry-active byte nonzero, signed primary HP below -53; campaign paths are excluded |
| Session opcode `0x48` | Resolved player; no mode/latch restriction in this handler, so an existing staged primary can be charged if the command is delivered |
| Session opcode `0xbe` | The call occurs only with a null primary; it creates rather than charging an existing dead primary |

Observed `0x48` production is character acceptance, not a failure-panel command. Teardown,
direct healing and imported-character restoration also clear stage without this penalty.
Ordinary UI delivery of a staged-primary `0x48` remains Unknown. The loss arithmetic's prior
level-band result remains in `HERO-SKILLLOSS-075` (amended); it does not establish mode reachability.

### the purchase — `FUN_004f7cd7`, one command

Server command opcode **61** and no other. Requires `server+0x0c == 0`; slot from the packet,
bounded `0..5`, so it is the only shipped server operation that reaches `Skill.General`; price
`ftol(1.1^word[actor+0x116+2*slot] * 200)` charged against `player+0x38`. For slots 1..5 the word
is the maintained base; for General it is the decoupled serialized word above. Then
`skill[slot] += 1` **with no ceiling**, `xp[slot] = S(new level) + 1`, base slots 1..5
re-snapshotted, derive called.

The first crossing of each fixed-width stage is below. Order is by the numeric level/base input;
live General and its serialized price shadow are independent inputs, so this is not one execution
timeline.

| order | input transition | wide calculation | stored/consumed value | G2 class |
|---:|---|---|---|---|
| 1 | new level 152 → 153 | `S(152)=1,957,437,581`; `S(153)=2,153,181,439` | purchase stores low-dword `S(153)+1 = 0x8056f100`, signed `-2,141,785,856` | signed-`i32` policy; treating the dword as unsigned needs a consumer-wide engine change, not a schema change |
| 2 | new level 160 → 161 | stored `S(n)+1`: `4,195,942,440` → `4,615,536,784` | low dword wraps to `0x131b8090`, signed `320,569,488` | `u32` storage; lifting needs wider per-slot experience, total, and save fields, not a wider `Humans` level column |
| 3 | maintained price base 169 → 170 | `P(169)=1,978,763,028`; `P(170)=2,176,639,331` | low dword `0x81bce163`, signed `-2,118,327,965`, is compared to the purse with signed `JLE` | signed-`i32` policy inside the existing price dword |
| 4 | maintained price base 177 → 178 | `P(177)=4,241,654,286`; `P(178)=4,665,819,714` | low dword wraps to `0x161ac242`, signed `370,852,418` | `u32` storage; lifting needs wider price, purse, and message state |
| 5 | new level 385 → 386 | `S(n)`: `8.6334382265275638e18` → `9.4967820491803197e18` | the `FISTP qword` input first exceeds signed-qword range | conversion-policy boundary; the unexecuted original result is Unknown |
| 6 | maintained price base 402 → 403 | `P(n)`: `8.7274913946611507e18` → `9.6002405341272678e18` | the `FISTP qword` input first exceeds signed-qword range | conversion-policy boundary; the unexecuted original result is Unknown |
| 7 | new level `0x7fff` → `0x8000` | the stored word is 32768 | `MOVSX` passes `-32768` to `S` | signed-`i16` policy inside the existing word; no schema or save-width change |
| 8 | old level `0xffff` + 1 | mathematical result 65536 | `ADD AX,1` stores `0x0000` | `u16` storage/format; lifting needs wider actor, `Humans`, and save fields |

The raw helper body, caller and low-dword stores are asserted by `tools/general`; its boundary CSV
asserts the exact values above on both roots. The total at `actor+0x130` is another signed dword:
because it accumulates all six slots, its sign and wrap thresholds depend on the other five values
and can precede the per-slot rows above. Between the first dword wrap and the qword-input boundary,
later sign and wrap transitions repeat modulo `2^32`; the table records the first crossing of each
width stage, not every repetition. Runtime behaviour downstream of every unexecuted
transition, including the exact CRT result after either qword-range crossing, remains Unknown.

**Store both fields.** Experience is not a display of the skill vector and the skill vector is not a
display of the experience: after any raise `xp[slot]` lies in `(S(n-1), S(n)]` for level `n`, and
only the loss ever brings them back into agreement.

## 5 — the derive, `FUN_004f7dfc` (`vt+0x50` of both human classes)

Reached **only** through the `.rdata` slots `0x0059c498` and `0x0059c520`. The order is normative:

```
 0  stat caps    stat_i = min(stat_i, 50 + (int8)actor[0xd4 + i])            i = 0..3
 1  healthMax    h = body * (fighter ? 2 : 1)             -> +0x96 UNCONDITIONALLY (004f7f53),
                 [skip the two lines below if that product is 0 -- 004f7f66 tests the
                  value just stored, NOT the Data.bin HealthMax column, which is gone]
                 h = ftol( h + log1.1(XP/5000 + 1) * (fighter ? 2 : 1) )
                 h = ftol( h * (pow(1.1, body)/100 + 1) )                    -> +0x96
 2  manaMax      [gate FIRST, on the streamed Data.bin ManaMax column at 004f8014:
                  if it is 0, current mana = 0 and the whole of 2 is skipped]
                 m = spirit * 2                                              -- overwrite second
                 m = ftol( m + log1.1(XP/5000 + 1) * (fighter ? 1 : 2) )
                 m = ftol( m * (pow(1.1, spirit)/100 + 1) )                  -> +0x9c
 3  sight        ftol( ((mind + reaction)/25 + 4) * 256 )                    -> +0xa4 u16, 1/256 cell
 4  capacity     body * 10 + 1                                               -> +0x92
 5  speed        reaction < 12 ? reaction : reaction/5 + 12
                 + 10 if typeID == 0x13 or 0x15                              -> +0x8c
 6  load         weight(+0x8e) + (money < 64000 ? money/2 : set load = 32000) -> +0x90
 7  overload     if load >= capacity: speed -= load/capacity ; speed = max(speed, 6)
 8  damage       dmgSpread = dmgBase = ftol( pow(1.1, body) / 20 )           -> +0xb5, +0xb4
                 the pair is (base, SPREAD): the roll is base + U[0, spread] -- section 8
                 NOTE the Body term is added to BOTH bytes, so it raises the roll's
                 minimum once and its maximum twice; step 11's skill term hits ONLY
                 the base, and nothing on either path is multiplicative
 9  toHit        ftol( (pow(1.1, body) + pow(1.1, reaction)) / 5 )           -> +0xa6
10  skills       skill[i] = min(100, base[i] + bonus[i])                     i = 1..5
11  active skill toHit += 3 * skill[active] ; dmgBase += skill[active]/5     active = +0xb6, 0 = none
12  zero         memset(actor+0xbe, 0, 0x16)   -- defence, absorption, both resistance arrays
13  defence      reaction / 3                                                -> +0xbe
14  protections  prot[i] = spirit / 2                        i = 1..5, +0xc4 .. +0xcc (u16)
15  modifiers    FUN_004f54f8(actor+0xd4, actor)             -- section 5a
16  clamps       current health/mana clamped to their maxima; defence, absorption and load
                 floored at 0; prot[i] = clamp(min(spirit/2 + 70, prot[i]), 0, 100);
                 skill[i] = clamp(skill[i], 0, 100)
17  mover        [actor+0x154] + 0xa = (byte)speed            -- this is how speed reaches movement
18  effects      if actor+0x140: walk the spell list
```

Order is load-bearing at four points: the caps run first, so everything downstream sees the capped
value; health and mana are computed **before** the skills are restored, so a skill change reaches them
only on the next recompute; the `memset` runs between defence's inputs and defence's write, so all of
`+0xbe..+0xd3` restart from zero every time; and the equipment pass runs **between** the bare values
and the clamps, so the clamps bound equipment rather than stats.

**Current health `+0x94` and current mana `+0x9a` are live state, not derived.** Only the maxima are.

**`БРОНЯ` (absorption, `+0xc0`) is derived from no stat** — zeroed at step 12 and never rewritten by
the derive. It and `ЗАЩИТА` (defence, `+0xbe`) are two separate fields with two different sources,
and both are supplied at step 15.

### Step 0 in full — what bounds a stat *in principle*

```
stat_i = min(stat_i, 50 + (int8)actor[0xd4 + i])     i = 0..3, and the smaller is WRITTEN BACK
```

The `50` is a hardcoded immediate, four times over. It is **not** class-dependent, **not** per level
and **not** in any file. The modifier byte is signed and is written by **effects only**: an effect
that raises a stat raises the modifier by the same amount (unless its own flag bit 3 is set, in which
case the next recompute takes the raise back), and clamps **the stat alone to 100**.

```
a stat no effect has ever touched      capped at 50   (chargen can only reach 43)
the ceiling in principle               100            imposed by the effect arm, not by the cap
fighters vs mages                      identical      the class flag reaches neither operand
```

The `100` is the same `CMP EAX,0x64` in all four arms, and each arm owns one modifier byte:
`+0xd4` Body, `+0xd5` Reaction, `+0xd6` Mind, `+0xd7` Spirit. Each of the four was enumerated
image-wide and each returns the same shape — **4 hits, 2 owners, 0 in orphan code**: the derive
reads it twice, the effect dispatch reads and writes it. Because the cap is step 0 and every effect
arm ends in a full recompute, **whatever reads a stat sees the capped value**, and an untouched
stat is therefore 50, not its nominal maximum.

Consumer note: the modifier is a **byte** and is not clamped, so a long enough chain of raises wraps
it negative and the cap with it.

### When the derive runs on a map-placed person, and what its inputs are (`HERO-HP-071`, `HERO-HP-072`)

The same routine is the `vt+0x50` of both human vtables, and on the `.alm` placement path it runs
**twice** before the actor is playable:

```
FUN_004f9065  004f964e  CALL FUN_004f7c28       ; actor+0x130 <- the six skill columns
              004f965b  CALL [vptr + 0x50]      ; the derive
              004f966b  actor+0x94 := actor+0x96
FUN_004e26bb  004e2ca8..004e2d19                ; the placement record's stat overrides
              004e2d28  CALL [vptr + 0x50]      ; the derive again, on the overridden Body
              004e2d41  actor+0x94 := rec+0x20  ; only when rec+0x20 != -1
```

So for such an actor the **`Data.bin` `HealthMax` column is dead**: the streamer writes it into
`+0x96` (and `+0x94`), step 1 overwrites `+0x96` at `004f7f53`, and the surviving
`+0x94 = min(+0x94, +0x96)` clamp at `004f84c4` is itself overwritten by `004f966b`. The inputs
that decide a placed person's health maximum are Body (capped at 50 by step 0, the block being
zeroed at spawn), the six skill columns through step 1's experience term, and whether the `ManaMax`
column is positive — which sets the class bit at `004f94d0` and therefore the `fighter ? 2 : 1`
multiplier. Over the 215 shipped `Humans` rows the derived maximum equals the column on **1** and
reaches **20.2x** the column on one row. The `-1` defaults for such an actor are **not** the base
constructor's: `FUN_004f6ded` runs after it and raises Mind and Spirit to 30, healthMax and health
to 50, and speed to 16.

## 5a — step 15, the modifier fold `FUN_004f54f8(actor+0xd4, actor)`

The actor carries **three parallel copies of one layout**, built together by every constructor:

| live | modifier | base | contents |
|---|---|---|---|
| `+0xa6` | `+0xe6` | `+0x114` | `[toHit u16][6 skill u16][dmgBase u8][dmgSpread u8][active u8][2 u8][elemBase u8][elemSpread u8][elemKind u8]` |
| `+0xbe` | `+0xfe` | — | `[defence u16][absorption u16][6 protection u16][6 damage-kind u8]` — 0x16 B |

The modifier block is `actor+0xd4 … +0x113`, exactly 0x40 bytes: four stat-cap bytes, then
`+0xd8` speed, `+0xda` capacity, `+0xdc` healthMax, `+0xde` healthRegen, `+0xe0` manaMax,
`+0xe2` manaRegen, `+0xe4` sight, then the two sub-blocks at `+0x12` and `+0x2a`. The fold is:

```
speed += +0xd8 ; capacity += +0xda ; healthMax += +0xdc ; manaMax += +0xe0 ; sight += +0xe4
if the resulting speed < 0:  the MODIFIER +0xd8 is zeroed, not the speed
+0xbe block += +0xfe block   (all 0x16 bytes)
+0xa6      += +0xe6 ; dmgBase += +0xf4 ; dmgSpread += +0xf5 ; four more bytes added, the last ASSIGNED
```

The six skill words are **not** in this fold — the derive adds `+0xe8 + 2i` itself at step 10.
The complete attack fold is 53 instructions with seven `ADD`s (one word and
six bytes) and one byte assignment. It also skips active-index mirror `+0xf6`
and final bytes `+0xfc/+0xfd`; the active live index is written separately.
— HERO-FOLD-035, SAV-HUMFOLD-446

The fold is not a universal loaded-state boundary. Before its later clears,
derive reads `word[a8+2*b6]` for any nonzero byte `b6`; the selected indexing
step has no upper bound. Malformed synthetic indices10,13,32,39,42 make
`bc/bd`, slot-zero `c2/c3`, General shadow `e8/e9`, mirror/secondary pair
`f6/f7`, and unnamed `fc/fd` affect to-hit/damage. Ordinary producer reach of
those indices is Unknown. The damage resolver also has an unchecked
`target+ce+attacker.b6` byte index. — SAV-HUMINDEX-464

The actor-state sender computes wrapped sums of live physical/secondary/
elemental damage bytes with an effective damage mask, without a local derive
call. This is presentation, not the combat formula. Phase-12 list order puts
that actor's projection before its regeneration read of `de/e2`; ordinary
subticks walk Effects before orders/actions. Generic Effect dispatch reads
`d8` before its own later derive, while special Effect identity17 on a Human
changes `e4/a4` by signed magnitude times256 without generic derive.
The admitted strike independently consumes live absorption, resistance,
secondary protection and elemental selector under distinct component gates.
— SAV-HUMPROJECT-461, SAV-HUMTICK-462, SAV-HUMSTRIKE-463

None of these conditional routes establishes the absolute first read after
loading: earlier commands, callbacks and other receivers may already have
changed the actor. Nonzero capacity/secondary-modifier producers and ordinary
residual-field semantics remain Unknown; source-byte passthrough and the
temporary chargen Human's derive do not close those authoring inputs.
— SAV-HUMFIRST-465

Fresh hero and hired Human creation share one Human constructor chain. Its
embedded block calls leave live-tail `bc/bd` unchanged and clear modifier-tail
`fc/fd`. A 488-byte Human allocation rounds to496, above the image's480-byte
small-block threshold, and reaches imported `HeapAlloc` with flags0; returned
payload contents are not established. Both pairs reach raw SAV writes, but
the constructor clear does not close later Effect/indexed writes or first-save
values. — SAV-HUMALLOC-504, SAV-HUMNEW-505, SAV-HUMNEWSAVE-507

The156 shipped nonempty Human weapon-definition cells per root all join to
attackType1..5. Weapon equip copies the low byte only for signed type below10;
every type at least10 and removal set active selector0. Thus positive type10/42
does not yield the known tail aliases. This finite producer bound excludes
neither negative/custom values nor archive, starting-skill and later lifecycle
aliases. It does not establish a default vector. — SAV-HUMSEL-506,
SAV-HUMNEWSAVE-507

The original-process checkpoint verifies only process-local startup-directory
redirection in a disposable copy. It did not load the ordinary save present
there or arm Human observations. First selected-field reads/writes after
normal load therefore remain Unknown; the startup result does not decide
whether an earlier derive replaced any saved Human field.
— SAV-HUMRUNTIME-476

The proposed restricted-token write guard failed on synthetic outside DELETE
and NULL-DACL write-open controls before any original child was launched.
This adds no Human identity, loaded field value or chronological read. The
ordinary-load prerequisite remains an independently verified write/child
boundary; failure of this one profile is not proof that all safe execution
is unavailable. — SAV-INITGUARD-488

## 5b — who fills the modifier block

**Equipment** has class-specific stores. Armor/Shield add into the modifier and
live defensive copies. Weapon writes modifiers and calls derive:

```
Armor   slot = item+0x50, stored at actor+0x198 + 4*slot ; block item+0x52 -> +0xfe and +0xbe
Shield  single slot actor+0x78                           ; block item+0x50 -> +0xfe and +0xbe
        equipping one drops a two-handed weapon
Weapon  single slot actor+0x74, melee (attackType < 10):
          dmgBase mod += w+0x60 ; dmgSpread mod += w+0x61 ; defence mod += w+0x6a ;
          toHit mod += w+0x52  ; active skill = low byte of cached definition parameter5
        ranged (attackType 0xb / 0xc):
          the damage goes to +0xf9/+0xfa, +0xfb = 1 or 2, toHit mod = skill[0], active skill = 0
unequip subtracts additive terms, but ranged toHit assignment is not inverted:
          removal subtracts w+0x52 from the General value assigned to +0xe6
```

The universal inverse rule does not hold: prior to-hit modifier 5, General 17
and item to-hit 11 give melee 16 then 5, but ranged 17 then 6 in the bounded
original instruction slices. Weapon directly writes modifiers and invokes
derive; Armor/Shield explicitly add into both defensive copies. Full runtime
cycles and transitive callback interleaving remain Unknown. The bounded local
event sequence is specified below.
— HERO-EQUIP-017, SAV-HUMEQUIP-447

### Local equipment sequence

Armor removal refreshes negative weight, subtracts defence, clears its slot and
sets flags, then removes Effects. Shield removal refreshes and subtracts first
too, but removes Effects before flag update and slot clear. Weapon removes
Effects first, then changes attack modifiers/active skill and derives;
range/timing, weight, flags, owned-Spell deletion and slot clear follow. On
attach, Weapon prepares its Spell before eviction and derives before
timing/range and weight. Armor/Shield store their slot and add defence before
weight and Effects; Armor sets flags before Effects, Shield after. Same-slot
eviction calls the item; opposite-hand two-handed eviction calls the actor
wrapper and reinserts before installing the new item (`SAV-EQUIPORDER-552`,
`HERO-EQUIP-017`, `ITEM-ARMFOLD-033`).

The active selector comes from definition parameter5 cached after displacement
and the new-weapon slot store; melee copies its low byte, ranged11/12 and
removal assign zero (`SAV-HUMEQUIP-447`). Range uses a different source:
after derive, attach/removal reread current Weapon byte `+50` and add/subtract
its minus-one delta in actor byte `+12c`. Removal assigns timing8/4. This is
not a guaranteed inverse cycle across callbacks (`SAV-EQUIPORDER-552`).

State0 Effects run in forward list order, on removal as well as attach. Each
normal-return general Effect dispatch reaches actor `+50` before the next item
Effect; an empty list does not. The iterator prefetches one next-node pointer,
not all remaining Effect values. Actual callback list mutations remain Unknown
(`SAV-EQUIPEFFECT-553`). Command22 adds container removal/reinsertion and a final
zero-weight refresh; that refresh derives only on a changed load quotient.
Humanoid/Human wrappers add no trailing derive; Unit's trailing `+54` is Token
value, not `+50` derive (`SAV-EQUIPCALL-554`, `ITEM-EQUIP-006`). Derive and Effect
reads therefore occur between local stores; no atomic final-state or pure
callback contract follows (`SAV-EQUIPOBS-555`).

A weapon's numbers are `round(column × shapeFactor × materialFactor)` (`+0.5` then `ftol`); runtime
column *i* is `Data.bin` title *i+1*, so they are the shipped titles `@.physicalMin`,
`@.physicalMax`, `@.toHit`, `#.deIrnce`, `weight`. **The two factors are now decoded**
(`ITEM-SCALE-017`, partially retracted; the ladder below is the `ITEM-LADDER-019` correction): each is one f64 out of the Shapes record `[item+0x46]` and the Materials record
`[item+0x45]`, at `record + 0x20 + 8×(title−1)` — so the damage factor is the column the game itself
names `@.damage`, the to-hit factor `@.toHit`, the defence factor `#.defence`. **Both ends of the
damage pair take the same `@.damage` factor**, so an item scales symmetrically; the actor does not. **The damage pair is not stored as it is
written**: `w+0x61 = round(physicalMax × s × m) − w+0x60`, i.e. the spread. On the defensive block
the first word is `#.deIrnce` (defence) and the second `#.absorbtion` (`БРОНЯ`).

**Effects**, one per typed increment. `FUN_00501a22` takes `kind = effect+0x3c` (`0..0x31`) and
`v = param × effect+0x40`, jumps through a 50-entry table, and ends every arm with a full recompute.
The kinds are `Data.bin`'s Magic rows by index and by name: `price`; the four stats (each with its
cap byte); `health`/`healthMax`/`healthRegeneration`; `mana`/`manaMax`/`manaRegeneration` (the last
two mage-only); `toHit`; `damageMin`/`damageMax`; `defence`; `absorbtion`; `speed`; `rotationSpeed`;
`scanRange`; six protections; twelve class-gated skill slots; three lores; `castSpell`/`teachSpell`;
six damage elements; `damageBonus`. Fifteen arms branch on `target->vt+0x30()` — modifier copy when
the target recomputes, live field when it does not. `itemLore`, `magicLore`, `creatureLore` and
`castSpell` are **dead arms**.

## 6 — the class discriminator

`actor+0x4c` bit 2 is the **fighter/mage** flag, derived when the streamed mana maximum is positive
and also set by spellbook construction. `FUN_0051cb70` tests it. Sex is a **different** test:
`FUN_0051cb90` returns `typeID ∈ {0x22, 0x24}`. In non-zero player-character constructor mode the
overwrite is `typeID = gender + (mage ? 0x23 : 0x21)` — gender in the addend, class in the base.
Zero-mode map Humans retain their Data.bin typeID, so this sex predicate is not universal to the
Humans class (`PARTY-M20-031`). The derive calls only the class test; its two conditional edges are the
health and mana multipliers at steps 1 and 2. Outside it there are fifteen more, all in the effect
dispatch: twelve skill arms in two gated sets of six over the same six fields, plus `manaMax` and
`manaRegeneration`. The fighter predicate is the exact negation of the mage one; the polarity is
fixed by the health multiplier, which doubles when the bit is **clear**. On the UI side the same axis is `record+0x18c` bit 1 (mage) and bit
2 (female), which pick the fighter/mag skill column, the `chrgen1m`/`chrgen1f` tips and the
`FacesMM`/`FacesMF`/`FacesFM`/`FacesFF` list.

## 6a — how the character is DRAWN (`HERO-APPEAR-040`…`HERO-APPEAR-046`)

The class discriminator above decides *stats*. What decides the **drawn unit class** is a separate
chain, it lives on the client, and it does not read the stat discriminator at all.

What the server sends is not a drawn class. `FUN_004e7de3` pushes, under field-mask bit `0x4000`,
`FUN_00521890(actor)` — nine instructions returning `word[actor+0xe]`, the `typeID` — and then the
face byte `actor+0x4b`. The client parks them in `[0x005cd5c8]` / `[0x005cd5cc]` and assigns the
first to `drawable+0x20`, which is the subscript the frame selector `FUN_0045bf00` uses on the
`units.reg` class array `0x005eb674`.

`FUN_0045f850` then splits on that id:

| id | meaning | what happens to `+0x20` |
|---|---|---|
| `< 0x1a`, `>= 0x40` | a class a map places | left alone — it **is** the `units.reg` `ID` |
| `[0x20, 0x40)` | a player's character | discarded; see below |

The shipped roster occupies `1..27` and `64..80`, so `[0x20,0x40)` holds no class and is free for the
four archetypes `0x21..0x24`. For them the routine computes `typeID - 0x21`, files bit 0 into
`drawable+0x18c` bit 2 and bit 1 into bit 1 — the **same** assignment the chargen record uses, so
bit 2 is the sex axis and bit 1 the mage axis — sets bits 0 and 3, and stores the literal `1`.

`FUN_0045fb00`, the next instruction in the caller, does the real work. It returns at once unless
`+0x18c & 1`. Otherwise it builds a **body name** from a twelve-slot array of visible equipment at
`drawable+0x15c`, which the equipment messages fill from twelve `u16` at `msg+0xc`:

```
name = registryString[(byte[slot0 + 6] & 0x1f) - 1]        slot 0 absent -> index 0
if slot1 != 0            name += "_"
if mage and name == "unarmed"   name = (action == 6 ? "mage_st" : "mage")
```

and then maps the name to the drawn class with a seventeen-arm chain:

| name | id | name | id | name | id |
|---|---|---|---|---|---|
| `unarmed` | 1 | `axeman_` | 8 | `pikeman_` | 13 |
| `unarmed_` | 2 | `axeman2h` | 9 | `archer` | 14 |
| `swordsman` | 3 | `clubman` | 10 | `bowman` | 14 |
| `swordsman_` | 4 | `clubman_` | 11 | `xbowman` | 15 |
| `swordsman2h` | 5 | `pikeman` | 12 | `mage_st` | 24 |
| `axeman` | 7 | | | `mage` | 23 |

Seventeen names, sixteen distinct ids -- `archer` and `bowman` are the only pair that share one. `bowman` ships no sheet on either root and is dead.

**The pixels do not come from `units.reg`'s `File`.** The same routine composes

```
graphics\units\ <dir> \ <name> \sprites.256      -> drawable+0x194
graphics\units\ <dir> \ <name> \spritesb.256     -> drawable+0x198
```

where `<dir>` is `heroes` for a mage, `heroes_l` for a fighter with nothing in slot 7, and otherwise
`material.reg`'s `Material[word[slot7 + 6] >> 12].Path` — sixteen blocks carrying `heroes` for
0–7 and 14 and `heroes_l` for 8–13 and 15. So the class record supplies only the **geometry**:
phases, `Width`/`Height`, `CenterX`/`CenterY`, `Flip` and the three timelines.

`+0x19c` caches the last name and `+0x1ac` the last material index; nothing is reloaded unless one of
them changed.

**It is recomputed while the mission runs**, on the ordinary field-masked state message (opcodes
`0x6c`/`0x6e`/`0x6f`/`0x70`) as well as on the two equipment messages (`0x9c`, `0x76`) and from the
character screen. The server picks between the two equipment opcodes on the same `[0x21,0x40)` test.

**A consumer must therefore not store a drawn class on a player's character.** It must keep the
twelve visible-equipment slots, recompute the name whenever they change, load the sheet by path, and
take the geometry from the resulting `units.reg` record. For a unit a map places it must do none of
this: that actor's drawn class is its `typeID` and nothing derives it.

Claims: `HERO-APPEAR-040`…`HERO-APPEAR-046`, `UNIT-APPEAR-030`. The name list's order and the
second sheet's consumer were open here and are settled in §6b.

## 6b — the twelve slots, the `u16` in each, and where the sheets come from (`HERO-APPEAR-047`…`HERO-APPEAR-055`, `HERO-FIGURE-062`)

**Who sends them.** `FUN_004e873b` sends nothing at all unless `actor->vt+0x30()` — the humanoid
predicate — is true, so a non-humanoid actor's twelve slots stay null for its whole life. With no
named recipient it walks the player list and calls itself once per player.

**Which field becomes which slot.** Both opcodes visit `actor+0x74`, `actor+0x78`, then
`actor+0x198+4i` for `i = 3..12` — `ITEM-EQUIP-006`'s own 1..12 equipment numbering. Wire slot `k`
is equipment slot `k+1`, and an empty slot is sent as 0 rather than omitted.

**The two opcodes do not carry the same bytes.**

```
0x9c  twelve bare u16 at msg+0xc+2k, each = word[item+0x40]
0x76  twelve serialised item records from msg+0x13, each 7 + len bytes:
        rec+0  u16   word[item+0x40]        -> slot+0x06
        rec+2  u16   item+0x42, the count   -> slot+0x10
        rec+4  u8    bit7 = item+0x14 != 0  -> slot+0x08
                     bit6 = item+0x8  != 0
        rec+5  u8                           -> slot+0x09
        rec+6  u8    a length               -> slot+0x0a
        rec+7  ...   len bytes              -> malloc'd at slot+0x0c
      rec+0 == 0 means the slot is empty.  An empty slot is sent as a
      throwaway Item constructed, serialised and destroyed.
```

**The `u16`.** Every bit is consumed, by `FUN_00483c80`, which turns it into a seven-digit name:

```
  bits 15..12  A = v >> 12          the material.reg index  (§6a's slot-7 axis)
  bits 11..8   B = (v >> 8) & 0xf   the equipment slot, 1..12
  bits  7..5   C = (v >> 5) & 7     unnamed
  bits  4..0   D = v & 0x1f         the item's Data.bin definition row index

  name = sprintf("%02d%02d%1d%02d", A, B, C, D)
  or     sprintf("%02d%02d%03d",    A, B, v & 0xff)   when B == 14
```

`FUN_00525d60` assembles it from `item+0x46` (A), the item's kind (B), `item+0x45` (C) and
`item+0x0c` (D), each masked to a byte and `OR`ed — so a definition row index above 31 would corrupt
C.

**What the name addresses.** Two trees, and they are different pictures:

```
graphics\equipment\ <figure> \ primary   \ <name>.256    the info window's figure, per slot
graphics\equipment\ <figure> \ secondary \ <name>.256    slots 3, 8, 9, and 7 for a mage
graphics\equipment\ <figure> \ <face>.256                 the head, from drawable+0x24
graphics\inventory\ <name>.16a                            the icon; missing -> "Invalid item weared "

<figure> is one of mfighter, mmage, ffighter, fmage, chosen by (drawable+0x18c & 6) —
the mage bit and the SEX bit together.
```

`FUN_0045ed10`, virtual slot `+0x80` of the drawable class, loops all twelve slots and builds these.
**It is not the world renderer.** Every one of the 928 shipped equipment sheets holds exactly one
frame of 160×240, against a hero body sheet's 129–216 frames of 24×40 to 40×48, and the routine's
identified callers are the info-window routines. In the world nothing is composited over the body:
a held weapon is visible because §6a's body name changes.

**The name list is a shipped file.** `main\text\heropicture.txt`, loaded into the object at
`0x5ea220` at startup, 25 non-blank lines, identical on both roots, indexed by `D - 1`:

```
 0 unarmed      5 swordsman2h  10 axeman      15 pikeman   20 archer
 1 swordsman    6 clubman      11 axeman2h    16 pikeman   21 xbowman
 2 swordsman    7 clubman      12 mage_st     17 axeman    22 Sonic Beam
 3 swordsman    8 clubman      13 mage_st     18 axeman2h  23 Flame Thrower
 4 swordsman    9 clubman      14 pikeman     19 archer    24 swordsman
```

`bowman` never occurs, which is why §6a's seventeenth arm is dead. Two lines name no shipped body
directory. The list is 25 long while D is five bits wide.

**The second sheet is drawn.** `FUN_0045bf00` blits `drawable+0x194` and then `drawable+0x198` at
the same frame index behind one `+0x18c & 1` gate, the second centred from its own width and height
and the class record's `+0x30`/`+0x38`. Its six arguments are `dstX`, `dstY`, frame,
`[0x005eb4a0]`, sun shear, mirror — `TERR-SPR-038`'s push list, so the fourth is
`TERR-SPR-066`'s **shroud level** and the second sheet belongs to the silhouette family.
**`0x26` is not one of them** (§6c): `HERO-FIGURE-062` corrects the partially retracted
`HERO-APPEAR-056` by identifying it as the preceding gate's key, not a blit argument.

**A consumer must therefore** keep twelve slots per humanoid actor, fill them from the equipment
slots offset by one, decode each `u16` into the four fields above, use D of slot 0 against
`heropicture.txt` for the body name, A of slot 7 against `material.reg` for the directory, and treat
`graphics\equipment` as a portrait tree that never touches the world sprite.

Claims: `HERO-APPEAR-047`…`HERO-APPEAR-055`, `HERO-FIGURE-062`, `ITEM-APPEAR-023`, `ITEM-APPEAR-024`,
`UNIT-APPEAR-031`, `SPR256-EQUIP-042`.
Not established here: field C (§6c narrows it); which garment each of slots 2..6 and 8..11 is.

## 6c — composing the figure: two surfaces, and one of them is a hit map (`HERO-FIGURE-057`…`HERO-FIGURE-064`)

`FUN_0045ed10` builds up to twelve layers per figure and then draws each of them **twice**, into two
caller-supplied surfaces bound in turn with `vt+0x28`:

```
  pass 1   vt+0x18 = FUN_00428e60   ->  FUN_0044ea40 / FUN_0044db00
           (x, y, frame, palRow, mode)      palette = layer+0x1c + (palRow << 9)
           the picture

  pass 2   vt+0x40 = FUN_00429190   ->  FUN_0044dd40
           (x, y, frame, tag)               tag is the SIXTH argument of the blitter
           a byte-per-pixel stencil: for an opaque run FUN_0044dd40 copies no source
           pixel at all, it stores the tag  (0044dda2 / 0044dda8)
```

**The tag is the equipment slot number**, `drawable index + 1`, over all 28 draw sites; the head
layer takes 0. So the second surface is a map from screen position to equipment slot, which is what
an info window with clickable equipment needs, and it is a second, independent derivation of the
`k+1` slot map §6b gives from the wire.

**Draw order is program order and depends on the mage bit** (`0045f1fd`):

```
  mage       p7 head p11 p9 s9 p3 s3 p6 p4 p8 p0 p5 s7
  non-mage   p11 p10 p6 p3 p4 p8 p9 p0 p7 s3 p5 s8 s9   then  p1 or p0
```

Three asymmetries a consumer must reproduce rather than smooth over: index **2** (equipment slot 3)
is drawn by neither half and ships no sheet; index **10** (slot 11) is drawn only by the non-mage
half and ships no sheet; index **8** (slot 9) is **tagged but not colour-blitted for a mage**.

**Which held layer is on top.** `FUN_00483e70` takes slot 0's field D, indexes `heropicture.txt`,
and returns 1 for `bowman`, `archer`, `xbowman`, `axeman2h`, `swordsman2h`, `mage_st` — every
two-handed and ranged body. Predicate 1 → index **0** last; predicate 0 → index **1** last. Field D
is the only field of the `u16` this routine reads.

**The list accessor has no bound.** `FUN_004687f0` is `[[0x005eb3d4] + ([list+0xc] + i)*4]`, six
instructions, no test — and this caller runs only when slot 0 is **occupied**, so `D = 0` indexes at
**−1**. A consumer must clamp at both ends; `D ≥ 26` runs off the other.

**And the same drawable field feeds sound.** `FUN_0045e890` picks one of eight 48-byte voice banks
on the mage bit, the composed-figure bit, **whether slot 0 is occupied**, and the sex bit:

```
  mage bit set                          -> 0x5f1b20   m_mage    / f_mage
  composed-figure bit set, mage clear    -> 0x5f1bf8   mf_hero   / ff_hero
  neither, slot 0 occupied               -> 0x5f1c58   mf_merc   / ff_merc
  neither, slot 0 empty                  -> 0x5f1b98   m_peasant / f_peasant
```

`FUN_004af540` fills them at start-up from `sfx\<bank>\` + `select1/2`, `command1/2`, `retreat`,
`defend`, `idle`, `easy`, `hard`, `die` — 80 of 80 present on both roots. So an armed humanoid and
an unarmed one do not merely look different, they sound different.

Claims: `HERO-FIGURE-057`…`HERO-FIGURE-064`, `UNIT-FIGURE-032`, `ITEM-APPEAR-025`.
Not established: what the effect-list keys `0x26`/`0x30` name; what `[0x005ef990]`/`[0x005ef994]`
are; which garment each of indices 2..6 and 8..11 is.

## 6d — persistent companion joins

There are five shipped persistent-companion producer origins. Mission 30 town activation consumes
`AddHero=22` and constructs Fergard or Reniesta. Map instants transfer an existing `Hero` actor in
missions 40, 70, 100, and 140. A dynamic npc with `DataBinID=26` selects Humans server id

```
26 + sex/class selector + 4*floor(mission/40)
```

The selector bits are sex at bit 0 and class at bit 1. `MySex`, `MyClass`, and their inverted forms
derive them from the primary hero. The ten conditional source variants are recorded by
`HERO-JOIN-120` and compared across roots by `HERO-JOIN-128`.

Brian is mission-40 unit 6, `npc25`, fixed Humans row 42 `PC_Paladin`. Mission-70 Naira exists for a
male primary and is unit 151, `npc23`, Humans row 31 `PC_Naira_2`. Their implementation payload is:

| field | Brian | Naira |
|---|---|---|
| class / sex / actor type / definition type | fighter / male / `0x21` / 5 | fighter / female / `0x22` / 14 |
| face / figure / weapon body / geometry | 1 / `mfighter` / `swordsman2h` / 5 | 1 / `ffighter` / `archer` / 14 |
| Body / Reaction / Mind / Spirit | 41 / 39 / 25 / 21 | 39 / 41 / 22 / 26 |
| six skill levels | 0 / 25 / 3 / 0 / 0 / 1 | 0 / 24 / 5 / 0 / 0 / 38 |
| six skill XP values | 0 / 9834 / 331 / 0 / 0 / 100 | 0 / 8849 / 610 / 0 / 0 / 36404 |
| aggregate XP | 10265 | 45863 |
| max HP / MP | 157 / 0 | 177 / 0 |
| speed / sight raw / capacity / equipped load | 18 / 1679 / 411 / 457 | 20 / 1669 / 391 / 179 |
| to-hit / defence / absorption | 97 / 76 / 5 | 177 / 88 / 1 |
| physical base / spread | 20 / 15 | 17 / 11 |
| physical / fire / water / air / earth / astral protection | 0 / 10 / 10 / 10 / 10 / 10 | 0 / 32 / 13 / 29 / 13 / 13 |
| known spells / carried item instances | none / 0 | none / 0 |
| occupied worn slots | 1, 6, 7, 8, 9, 10, 12 | 1, 4, 5, 6, 7, 8, 9, 10, 12 |

Skill order is General, Blade/Fire, Axe/Water, Bludgeon/Air, Pike/Earth, Shooting/Astral. Sight is
the raw 1/256-cell word. Item codes, effects, prices, weights, damage and defence are in
`HERO-JOIN-125`.

The map ownership routine transfers the same actor object. It does not reconstruct the Human or
write `Player+0x34`. Current health and other mutable state must therefore be copied from the live
source actor rather than reset to the source-initial value. The transferred type remains in
`0x21..0x24`, so the actor keeps the ordinary Human inventory, worn slots, skill, spellbook and save
capabilities. The join does not make the companion primary.

Claims: `HERO-JOIN-120` through `HERO-JOIN-128`.

## 7 — resolving one hit (`FUN_004fbc92`, the actor's `vt+0x4c`)

For a unit target, the decoded blow and direct-damage effect use this resolver (`HERO-DAMAGE-022`).
Building targets use a different resolver (`UNIT-STRUCTDAMAGE-064`). Its second argument is a `0x16`-byte
combat block of the `+0xa6` layout — the attacker's own live block for a melee strike, an
`Effect_DirectDamage`'s copy at `+0x48` for a spell (`formats/magic`).

```
rand(n)     = rand() * (n+1) >> 15          uniform on [0, n] INCLUSIVE ; 0 when n == 0

dmg   = A.dmgBase + rand(A.dmgSpread)
        bless (spell 23) active on the attacker and rand(100) < its magnitude
                                            -> dmg = A.dmgBase + A.dmgSpread, no roll
        curse (spell 27) likewise           -> dmg = A.dmgBase, no roll
hit   iff  A.toHit + (rand(200) - 100) > target.defence
      or   attacker+0x4c & 0x10                 set by attackKind == 3: Bat_Sonic, Bee
      or   the same roll >= 90              (the roll < -90 arm is inert: it re-assigns 0)
miss  -> the physical component is 0, the elemental one still runs if A has no physical damage

dmg  -= target.absorption                   FLAT, and physical only ; floored at 0
dmg   = ftol( dmg * (100 - target[0xce + A.activeSkill]) / 100 + 0.75 )

second component, run whenever that pair is nonzero -- NOT gated on the hit,
and no absorption is subtracted from it:
        v = A[+0x11] + rand(A[+0x12])
        v = ftol( v * (100 - target.protection[2]) / 100 + 0.75 ) ; max(v, 0)
        A[+0x11]/A[+0x12] = actor+0xb7/+0xb8, written only by the units streamer's
        attackKind == 1 arm -- which NO shipped class takes, so this pair is 0 in play

elemental component, applied iff the swing landed OR A has no physical damage at all:
        v = A.elemBase + rand(A.elemSpread)
        p = target.protection[A.elemKind]   kind 1..5 = Fire/Water/Air/Earth/Astral
        v = ftol( v * (100 - p) / 100 + 0.75 ) ; max(v, 0)

result = max(sum, 0) ; the melee strike subtracts it from target+0x94
side effect: target+0x40 = the attacker, target+0x48 = A.elemKind when the attacker is a mage
```

The six bytes at `+0xce … +0xd3` are therefore **resistance by weapon class**, indexed by the
attacker's own active skill slot (0 = General, 1..5 = blade/axe/bludgeon/pike/shooting), and they
are never re-derived for a human — an equipped item is their only source. A protection is a
percentage, clamped to `[0,100]` by step 16, so 100 is immunity.

## 8 — regeneration (`FUN_004f4204`, the actor's `vt+0x14`)

The regeneration route first requires dword `actor+0x54 != 16` and signed
health>0. Health0 is untouched; negative health belongs to the separate decay
route. The helper computes signed `i32(server.sub - actor.due)` from server
`+04` and actor `+138`; a result>80 selects local rate3, otherwise1. This is
the arithmetic on an invocation, not a new scheduler or post-load timing rule.
— HERO-REGEN-021 (amended), SAV-REGENORDER-531

| arm | current / maximum / period | modifier / remainder | extra local gates |
|---|---|---|---|
| health | `+94 / +96 / +98` | `+de / +a2` | current<maximum, period!=0, signed server full-counter remainder modulo4==0 |
| mana | `+9a / +9c / +9e` | `+e2 / +a3` | current<maximum; no period or full-counter filter |

All eight pool/max/period/modifier loads are signed16; the remainder loads
are unsigned8. Let `i32`/`i16` mean signed interpretation of the low32/16 bits,
`u8` mean low8 bits, and `T(a,b)` mean signed division truncated toward zero.
For a reached arm, with factor2 for health and1 for mana:

```
base = i32(current * 100 + remainder)
n = i32(maximum * factor)
n = i32(n * (modifier + 100))
n = i32(n * rate)
acc = i32(base + T(n, period))
q = T(acc, 100)
remainder = u8(acc - q * 100)       # first persistent store
current = i16(q)                    # second persistent store
current = min(i16(current), maximum) # signed read-back, third store
```

The `+100` sum is32-bit, not narrowed to a word. The routine actually performs
two separate `IDIV100` operations, storing the first remainder before the
second quotient. There is no lower clamp. The final upper-bound store does not
clear the remainder. Negative remainders become bytes157..255 and reload
unsigned; initial stored bytes100..156 are also consumed without validation.
Health32766/max32767/period1/modifier0/remainder0/rate1 yields32764 after
quotient98300 narrows, not32767. — SAV-REGENWIDTH-528, SAV-REGENSTORE-529

Health period0 skips its arm. Reached mana period0 causes a divide fault before
either mana store. Products and accumulators can wrap without fault. Under
stable gate operands, valid memory and rate1/3, the admitted product cannot
equal INT_MIN modulo2^32, so signed division overflow is excluded; injected
arm entries or corrupt locals are outside that proof. No operating-system
fault/recovery outcome has been witnessed. — SAV-REGENFAULT-530

The health stores at `004f42f1/004f4305/004f4349` precede callback
`004f435d -> 004e7de3`. Mana then rereads its operands, without another
health test or rate calculation, and stores at `004f43dc/004f43f0/004f4434`
before callback `004f4448`. The local body is not transactional. An unchanged,
normally returning first callback permits a later mana divide fault after
the health stores; actual callback mutations and loaded chronology remain
Unknown. — SAV-REGENORDER-531

The six pool/max/period words and two remainder bytes are serialized at their
existing widths; the modifier words are inside the earlier raw64 block at
`+d4`. Serialization does not normalize these values or establish a safe
initial vector. — SAV-REGENWIRE-532

The modifiers scale the base rate rather than adding a pool amount. No stat
is directly read: Spirit contributes through `manaMax`, Mind not at all.
— HERO-REGEN-021 (amended)

## 10 — the combat loop around that hit

`FUN_004f37be` is the actor's `vt+0x18`, run once per tick for every actor on the tick list
(`claims/move.md`). The attack lives in one arm of it, as a three-phase cycle:

```
actor+0x58   sub-phase        actor+0x6c   countdown byte
0            wind-up          = attackChargeTime + extra,  then sub-phase 5
                              extra = (dist*256 + 128)/200 when dist > 1, else 0
5            counting down    the blow lands on the tick the byte reaches EXACTLY 0
                              -> FUN_004fb942 -> unit: FUN_004fba0e -> target->vt+0x4c
                                                Building: FUN_004fbbc3 -> FUN_005046dc
                              then sub-phase 7
7            recovery         = attackRelaxTime + rand(3) + humanoidPenalty
                              at 0 -> sub-phase 0 and completion byte +0x136 = 1

humanoidPenalty = clamp(IDIV(runtimeWeaponWeight + 5*(30-Reaction), 12), 0, 12)
                  only with an equipped weapon and the Humanoid predicate; otherwise 0

one blow every   attackChargeTime + attackRelaxTime + rand(3) + extra
                 + humanoidPenalty + 2 actor ticks
                 the final 2 cross the completion-latch/order-machine boundary
                 62 ms per tick at the shipped speed index (16 tps)
```

Movement speed does not enter; Reaction does through `humanoidPenalty`. Neither the number of actors,
the tick order, nor adjacency changes the stored tick count: an actor
whose countdown is not zero cannot strike this tick, and one whose countdown is zero strikes
whatever `actor+0x5c` holds if it is in reach.

```
actor+0x5c    the combat target      a SECOND field; the destination is in the mover / order block
actor+0x54    3 = attacking, 1 = moving, 0xd/0xe = casting, 0x10 = torn down
actor+0x12c   reach, in cells        ctor default 1, raised by the weapon (weapon[0x50] - 1)
actor+0xa5    scanRange              read by the AI module, NOT by the attack path
actor+0x49    tokenSize (footprint)  vt+0x1c
actor+0x4a    movementType           > 1 leaves no corpse (Ghost, Bee, Bat_Sonic, Dragon)

in reach iff   attacker+0x12c >= FUN_004fb702(attacker, target), where
               d = max(|dx|, |dy|)                                1/256-cell units
               d -= ((sizeA + sizeB) << 7) - 256
               result = (d <= 0x180) ? 1 : (d + 0x40) >> 8
an order also requires the actor to be FACING the target before state 3 is entered
out of reach -> a move-to-actor order and state 1; the strike sub-phases run only in state 3
```

The start routine itself does **not** reject out-of-reach targets: it computes distance for
extra delay. Application rechecks reach (`HERO-CADENCE-115`, `UNIT-STRUCTREACH-063`). The
Building branch reads only combat bytes `+0x13/+0x14`, with a positive-spread gate and flat
subtraction 5; ordinary physical `+0x0e/+0x0f` is ignored. A physical Building hit writes the
word HP without clamping, unlike the direct-effect Building consumer, which clamps at zero.
Reach below 2 versus at least 2 selects physical message `0x71` versus `0x72`, not another
server strike callback (`UNIT-STRUCTDAMAGE-064`, `UNIT-STRUCTDELIVERY-065`).

> **⚠ REFUTED, and this paragraph is what it looked like while believed.** It read:
> *"Acquisition (the one route decoded, `HERO-TARGET-024`) runs from the order machine when the
> mover raises an event, not every tick: candidates within `actor+0x12c` whole cells, filtered to
> the hostile by the player diplomacy matrix, nearest wins."* `HERO-TARGET-024` is explicitly
> amended and **refuted** on that Medium clause by `AI-ACQUIRE-002`; this spec was never updated
> with it. Three of its four assertions are wrong: the rule an **unordered** unit runs is
> `FUN_0052e110` (13 callers) and not `FUN_005327d0` (one); the candidate set is everything the
> unit can **see** — a line-of-sight footprint of radius `actor+0xa5` — and `actor+0x12c` bounds
> only the **selection**, so *sight decides the population and reach decides the pick*; and
> **neither route takes the nearest** — the tie is broken on turn cost. Both routes fall back to
> **corpses** when no living enemy survives the filter. **Read `formats/ai` for acquisition; do
> not read it here.**

**A hit makes the two players mutually hostile** (bit 0, both directions) — that, not a per-actor
assignment, is retaliation (`HERO-AGGRO-028`, bounded by `AI-DIPLO-004`'s bit 1).

```
death, from the killing blow to the freed id  (HERO-DEATH-026, SESS-TICK-004,
                                               SESS-TICK-006, HERO-DWELL-065,
                                               HERO-FINISH-066, HERO-DYETICK-067,
                                               HERO-REVIVE-068, HERO-DECAY-069,
                                               HERO-ZERO-070)
  TWO clocks act on a dying body when stepping is admitted. Campaign outcome panels
  pause frontend stepping before teardown (SESS-DEFEAT-065); primary-fall reporting
  can also replace owned actors' HP with -50 (MISSION-DEFEAT-045).

  A. vt+0x18, FUN_004f37be, once per SUB-TICK
  tick 0        health <= 0 : defence halved, cast interrupted,
                              countdown actor+0x6c = dyingTime - 1
                              the corpse STILL OCCUPIES ITS CELLS
  while > 0     countdown -= 1                       decremented ONLY while positive
  at 0          the byte PINS there -- nothing writes it again while dying -- and
                from this sub-tick on the tail re-runs on EVERY invocation:
                  movementType > 1 -> health = -1000    (Ghost, Bee, Bat_Sonic, Dragon)
                  if health <= -10 -> torn down          <-- a STANDING condition,
                    every footprint cell released, the reserved next cell released,
                    the actor unlinked from the tick list onto the dead list,
                    gold dropped: treasureMin.1 + rand(treasureMax.1) when
                                  treasure.1 Gold > rand(100)

  B. vt+0x14, FUN_004f4204, once per FULL tick -- the SAME routine that regenerates
     the living, reached over the LIVE list, which a corpse is still in
  health > 0    regeneration
  health < 0    health -= 1 every FOURTH full tick        <-- what makes A's standing
                                                              condition come true
  health == 0   entered (JLE) and left untouched (JGE)    <-- NOTHING HAPPENS, EVER

  C. after the teardown, FUN_004f52ee over the DEAD list, once per FULL tick
  then          one health point per two FULL ticks; stages at -10, -20, -40
  at -600       the runtime id's bitmap bit is freed and actor+0x4 = 0

  so a ground unit felled to -h, 1 <= h <= 9, and left alone:
    >= dyingTime-1 sub-ticks   ~437 ms   the countdown, then it pins
    (10-h) x 4 full ticks      up to ~36 s   B carries health to -10
                                             A tears it down: cells free
    590 x 2 full ticks         ~19.5 min     C: the id returns to the pool
  and a ground unit felled to EXACTLY 0 does none of it, for the whole mission.

experience, per LANDED blow and only to a human class (the monster slot is a stub)
  xp = ftol( target.XPvalue * 0.5 * dmg / target.healthMax + 1 )
  then the Mind scaling of section 4; refused when the two share a player, or when the
  diplomacy byte has bit 1 -- and, past that, capped and credited to ONE slot by
  the rules of section 4, which is where it becomes a skill level
```

**`dyingTime` is a floor on the dwell, not the whole of it, and a corpse is not out of reach.**
The countdown is in **sub-ticks** — 62.5 ms each at the shipped default speed, so the human
default `dyingTime = 8` floors the dwell at ≈ 437 ms and not the ≈ 7 s a full-tick reading gives
(`HERO-DYETICK-067`). What ends the dwell is the standing `health <= -10`, and health keeps
moving after the killing blow: the strike tests only the **attacker**'s health and subtracts
unconditionally, so a body can be struck again (`HERO-FINISH-066`). A consumer must therefore
implement three things this spec used to imply away. **What ends the dwell without anyone doing
anything is the regeneration tick** — `vt+0x14` is dispatched over the *live* list, which a body
that has not been torn down is still in, and its dying arm removes one health point every fourth
**full** tick while health is strictly negative (`HERO-DECAY-069`). So a body felled at −1…−9
reaches −10 by itself in at most ≈ 36 s and is torn down then; nothing about it changes visibly in
the meantime, because the corpse stage byte does not move until −10. **The exception is one
value.** The arm is entered on `health <= 0` and left on `health >= 0`, so a ground unit reduced to
**exactly 0** — a blow whose damage equals the remaining health, which nothing prevents, since the
subtraction has no floor — is a corpse no clock can move: cells, reserved next cell and runtime id
held for the whole mission at corpse stage 1 (`HERO-ZERO-070`). Flying classes cannot reach it,
their arm forcing −1000 regardless. And the first tick's
effects are **reversible**: a heal that carries health back above 0 clears the corpse stage,
doubles the halved defence and floors health at 1, so a corpse inside the window is a live actor
with non-positive health (`HERO-REVIVE-068`). −10 is written at three sites, all in
`FUN_004fba0e` — the damage numeral, the experience/retaliation call, and the death arm's own
guard — and it is one line, not three rules.

Ranged is the same machine: the only melee/ranged test in the simulation is
`attacker+0x12c > 1`, and it picks between two **network opcodes** (0x71 melee, 0x72 shot). No
projectile actor carries the damage — the shot resolves through the same `FUN_004fbc92` after the
distance-proportional `extra` ticks.

## 9 — what a consumer still cannot reproduce

~~An item's own numbers cannot be predicted from its `Data.bin` template row~~ — **closed by
`ITEM-SCALE-017` and `ITEM-DMGCOL-018`**: the two doubles are the Shapes and Materials records' own `@.damage`/`@.toHit`/
`#.defence` columns at `+0x20 + 8×(title−1)`, and `formats/item` carries the fill
(`ITEM-SCALE-017`, `ITEM-DMGCOL-018`, both **partially retracted** by `ITEM-LADDER-019`). An effect's `kind`/`scalar` are read but not sourced — what builds an effect from a Magic row is
unread, and so is the `target->vt+0x30()` predicate the dispatch and every equip path branch on. Its
`+0x3d` **is** sourced: `permanent 0 / duration 1 / continuous 2 / charges 4 / singleuse 8`, the
game's own words (`formats/magic`). The two regeneration periods `+0x98` and
`+0x9e` have constructor/template producers (`HERO-REGEN-021`); serialization
can restore other values (`SAV-REGENWIRE-532`). Still open: ~~the school's skill training (price and experience granted); the experience
*distribution*'s two class arms in `FUN_004f72d7`~~ — **closed by `HERO-SKILLUP-073` and
`HERO-SKILLBUY-076`**, section 4: both arms
read, the price of a bought class-specific level is `ftol(1.1^base × 200)` and the granted amounts are the three
feeds listed there; what fills `actor+0x4c` bit 2 (`= & 4`, and note that `FUN_00523f50`'s `& 2` is
a **different** bit of the same byte), the flag those arms and three other combat decisions branch
on; and the membership of the actor list the acquisition scan walks. ~~Whether any shipped route
produces the packet carrying command opcode 61, so whether a level can be bought at all in a
shipped campaign~~ — **closed by `HERO-SKILLBUY-076` and `HERO-GENERAL-089`**: the school Train widget is the sole direct producer,
and its mapping reaches slots 1..5 but not General. The first crossing of every purchase width is
located but not executed: per-slot XP changes signed interpretation at new level 153 and wraps its
dword at 161; price does the same at maintained bases 170 and 178; the two `FISTP qword` inputs
leave signed-qword range at 386 and 403; level interpretation changes at `0x7fff -> 0x8000`; and
level storage wraps at `0xffff -> 0x0000`. The six-slot total can cross its own dword boundaries
earlier, depending on the other five slots; the derive's clamp reaches none of these transitions.
~~the character sheet's display formatting, including whether it prints `[base, base+spread]`~~ —
**closed by `HERO-SHEET-038`**: `FUN_00460480` zero-extends both bytes, **adds** them and prints
`"%d-%d"` (`0x005bcbd8`) at `00460aa1`…`00460ac2`, so the sheet shows exactly the roll's bounds
(`HERO-SHEET-038`). The values still leave the server raw, in `FUN_004e7de3`/`FUN_0047bf40`.

## The name (`HERO-NAME-079`, `HERO-TYPED-080`)

`actor+0x80` is an MFC `CString` and it is the actor's name. It is constructed and destroyed by
every actor constructor and by the destructor-shaped sixth, so a monster carries one too; it is
per **instance**, never per class. Four writers, and no shipped surface draws it for a non-hero:

| Writer | What it assigns |
|---|---|
| `FUN_004f9065` `004f9230` | the literal `"Unknown"` — only inside a branch no shipped asset reaches |
| `FUN_004f9065` `004f9255`/`004f9276` | a default-name array element, same dead branch |
| `FUN_004d3755` `004d3d8a` | the `CString` at the accepting request's `+0x18` — the **typed** name |
| `FUN_004e26bb` `004e292d` | the `.alm` spawner's own local, on the `Humans` sub-arm only |
| `FUN_004e0f26` `004e10cf` | the `Name` key of a `Humans.Hero` block that ships in no root |

Limits for a consumer: **10 bytes**, refused at the keystroke by the entry control
(`TEXT-NAMEIN-024`); charset as `TEXT-COLL-025` bounds it; storage is a `CString`, so there is no
fixed buffer to overflow; and the field is **not inert** — `FUN_004e32c2` skips any actor whose
name is empty before evaluating its sex/class predicate, so an empty name changes behaviour.

The two default-name arrays — 10 and 7 slots at `0x00609580` and `0x00609c18`, holding **9** male
and **6** female names with slot 0 empty — are `.rdata` literals built by two static initialisers,
so they do not localise; and the gate that reads them requires a template name beginning `Hero`,
which nothing in either install provides.

## Character-generation controls (`HERO-CHARGEN-082`…`HERO-CHARGEN-085`)

Character generation is two screens and two boundaries. Pre-create owns name, class/sex and
difficulty. Its Forward copies those fields into the main-frame draft, constructs a temporary
archetype record and opens detailed generation. Detailed Play joins or resolves the participant and
only then sends command `0x48`, which creates and installs the live actor.

Reset on the detailed screen writes pool 100 and Body/Reaction/Mind/Spirit 25. It preserves the
selected skill index, name, class/sex, difficulty and appearance. Preservation of the skill selector
does not mean preservation of the preview record: the reset reconstructs a fresh archetype actor,
zeros the five class skills, restores the chosen starting skill, recomputes all derived mirrors and
releases/rebuilds the twelve equipment-object slots (`HERO-CHARGEN-082`).

Detailed Back returns to pre-create and restores name, difficulty and class/sex. Stats, skill and
appearance are not written at that moment. The next pre-create Forward nevertheless reloads the
archetype defaults unconditionally, so stat and skill edits do not survive the round trip. Pre-create
Back is the destructive cancel: it releases the temporary record and returns to the participant/lobby
screen (`HERO-CHARGEN-083`).

The name field's UI limit remains ten stored bytes. Final Play adds server gates in this order:
non-empty; not the whole string `Self` or `Computer` under an ASCII-case-insensitive comparison; not
a byte-exact name already held by another participant; and no second character on a returning
participant. The first three failures are localized string-table indices 193, 194 and 195
(`HERO-CHARGEN-084`). The ten-byte limit is UI policy, not a server-storage limit: the participant
and actor both use `CString`. Lifting it changes no shipped asset bytes, but a consumer must also
replace the fixed preview-name copy rather than widening the control alone.

A failed Play is not transactional. It closes the detailed screen, sets the ready bit, resets the
connection and releases the temporary preview before the server rejects the name. It creates neither
a participant nor a live actor. Successful `0x48` is the commit point and copies `Player+0x18` to
`actor+0x80` unchanged (`HERO-CHARGEN-085`). Back/Play close is latched and idempotent while the screen
is inactive; repeated Reset reconstructs the same draft.
## Starting templates and non-item state

The four class/sex combinations select `PC_Danath`, `PC_Naira`, `PC_Fergard`, or `PC_Reniesta`.
Their ten `Humans` equipment cells can construct only Weapon, Shield, and Armor objects because the
cell position selects the class. Character generation adds one further direct weapon. It has no
MagicItems construction arm. No `Humans` equipment cell in either preserved root contains `Quest`.

After those paths, live-hero construction calls `FUN_004d3f6b`. When `server+0x0c == 0` and
`[0x005eb5a4] != 2`, the helper constructs `Quest Documents` by row name and appends it to the hero's
carried container before runtime-id assignment. The campaign construction path establishes the
first gate. The second gate was not measured in a fresh campaign, so actual item presence remains
Unknown (`ITEM-DOC-069`).

The character-generation command carries no purse value. The selected template, stats, appearance,
name, and weapon do not change the owning `Player` purse (`HERO-START-081`).
