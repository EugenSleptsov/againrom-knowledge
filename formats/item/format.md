# ITEM — the item, its container and the sack — specification (partial)

Level 3. Promoted, evidence-backed claims only. Core container behaviour is
`ITEM-CLASS-001`…`ITEM-CARRY-015`; the corrected damage ladder is `ITEM-LADDER-019`…`ITEM-PANEL-022`;
equipment-cell effects are `ITEM-EFFGRAM-070`…`ITEM-EFFSAVE-077`. Ledger: `claims/item.md`.

**Corrected 2026-08-03.** The scale records' double ladder in *A weapon's numbers* stood one
slot high in retracted `ITEM-SCALE-017` and `ITEM-DMGCOL-018`; the corrected ladder is
`ITEM-LADDER-019`…`ITEM-WEAPCOL-021`. Every factor figure published before
that date used the neighbouring column. See `claims/retracted.md` → `ITEM-SCALE-017`,
`ITEM-DMGCOL-018`.

**Status: partial (◐).** The class hierarchy, the record, the definition binding, the
container, the fourteen equipment slots, both move commands, the pick-up, death, the sack's
whole lifecycle, equipment-cell effect grammar and every serializer are read at instruction
level. **Not specified:** the per-column meaning of the `Data.bin` `Armors` / `Shields` /
`Magic Items` rows beyond the named consumers below; `item+0x44` producers outside the bounded
item constructor family; the meaning of serialized `item+0x47`; the transition that puts an
actor into tick state 2; the screen.

This is not a file format. It is the simulation's carrying layer, sitting on
`world.res:data/data.bin` for definitions and inside `claims/sav.md`'s stream for persistence.
Generation and pricing of the same objects are `formats/shop/format.md`'s.

## At a glance

Carried source-code2 transfer to the city table subtracts signed16 unit weight
times unsigned16 detached count from the inventory load, then inserts into the
tray container. A partial quantity invokes the Item split virtual; insertion's
`vt+50` is Item stackability, not actor derive. The carrying actor has no
direct stat recompute in this dispatcher arm. This protocol result does not
close the ordinary UI partial-pickup gesture or notification callback effects
(`SAV-CITYMOVE-512`).

```
Item      0x50 : Token          vtable 0x59c8e0 (27 slots)   defs <- Data.bin "Magic Items"
  Armor   0x68 : Item           vtable 0x59c950              defs <- Data.bin "Armors"
  Shield  0x68 : Item           vtable 0x59c9a8              defs <- Data.bin "Shields"
  Weapon  0x84 : Item           vtable 0x59ca00              defs <- Data.bin "Weapons"
Sack      0x44 : Token          vtable 0x59ca88 (19 slots)
container 0x24 : CObList        no vtable of its own; one per actor and one per sack
```

## The item record

The bounded per-class Item definition lookup does not extend unchanged to
actors. Unit uses its saved Units row without a bounds check; exact Humanoid
leaves a null definition; Human uses its saved Humans row only below type33,
otherwise row5. Its saved row byte is not rewritten. — ITEM-DEF-002,
SAV-ACTORBIND-544

Only the fields this specification pins. `+0x00..+0x3b` is the `Token` base, whose head
`FUN_00510e5c` serializes as `SAV-TOKEN-034`'s 37 bytes.

| Off | Type | Meaning |
|-----|------|---------|
| +0x08 | u32 | flag word; ORed together when two stacks merge |
| +0x0c | u8 | index of this item's row in its class's `Data.bin` collection |
| +0x14 | ptr | owning `Player` — stamped by the shop tray's destination-4 arm (`SHOP-TRAY-025`) |
| +0x1c | i32 | price (`SHOP-PRICE-011`); the `Token` "value" slot |
| +0x20 | — | ordered `CObList` of 0x48-byte Effects (`ITEM-EFFOBJ-072`) |
| +0x3c | ptr | the resolved definition row; re-derived on every load, never stored |
| +0x40 | u16 | packed item code; shelf spell id belongs to Effect+0x40, not this field (`ITEM-VALUE-115`) |
| +0x42 | u16 | **stack count**, default 1 |
| +0x44 | u8 | item kind; MagicItems prefix classification writes `3` Potion, `4` Scroll, `5` Book, otherwise 0. Literal 3 forces stackable |
| +0x45 | u8 | **Shapes row index** (the quality tier) — `FUN_004db6b5(name)`, the **5**-row table; entry 0 is `Common`, whose `@.damage` is 0.2 |
| +0x46 | u8 | **Materials row index** — `FUN_004db801(name)`, the **16**-row table; entry 0 is `Iron`. The `Armor` constructor bounds this one, not `+0x45`: `0050c3fe MOV DL,byte ptr [ECX + 0x46]` / `0050c401 CMP EDX,0xf` |
| +0x47 | u8 | serialized, unread |
| +0x48 | u16 | `ftol(shape.MagCap × material.MagCap)` — the only fill value taking no column |
| +0x4a | i16 | **per-unit weight**, default 1; `ftol(col3 × shape.weight × material.weight + 0.5)` |
| +0x50 | u8 | slot (`Armor`/`Shield`) or runtime column 11, else the literal 1 (`0050dc44`) |
| +0x52 | u16 | **to-hit bonus** — `Weapon::Equip` adds it to the actor's to-hit modifier `+0xe6` |
| +0x60,+0x61 | u8 | **the damage pair, `(base, spread)` — NOT `(min, max)`** (see below) |
| +0x6a | u16 | defence bonus (`#.deIrnce`) — added to the actor's `+0xfe` |

`IsStackable` = `vt+0x50` = `FUN_00508635`, inherited unchanged by all three subclasses:

```
stackable(item) = (item.kind == 3) || (item.effects.count == 0)
```

Splitting `qty` off a stack of `n` (`FUN_0050ebf9`):

```
if n <= qty:  the whole element leaves the container
else:         n -= (qty - 1); detached = item->vt+0x40(); detached.count = qty
```

The four detach-one virtuals invoke their concrete copy constructor. Copy walks the source effect
list in order and allocates a fresh Effect for every node. Weapon copy also reconstructs a fresh
owned Spell from the source Spell id. One exception matters to persistence: base Item copy omits
serialized byte `+0x47` without first running the default Item constructor, so a clone/split leaves
that byte indeterminate (`ITEM-EFFSPLIT-074`).

## Equipment-cell effects (`ITEM-EFFGRAM-070`…`ITEM-EFFSAVE-077`)

The first `{` in a `Humans` or `Units` equipment cell ends the right-trimmed item head. The effect
tail begins after it and ends at the first `}`; absent close means end-of-string and text after the
first close is ignored. The tail is split on commas from left to right. Each accepted element is
appended to an ordered `CObList`; unknown elements are skipped independently, and valid duplicates
remain in order. The parser has no count comparison or fixed element budget.

An element is case-folded and has `key=operand[:suffix]`. The first `=` and first `:` are the
delimiters. Index 0 of the 50-entry key table is a rejected sentinel; kinds 1..49 are:

```text
 1 price              2 body               3 mind               4 reaction
 5 spirit             6 health             7 healthMax          8 healthRegeneration
 9 mana              10 manaMax           11 manaRegeneration  12 toHit
13 damageMin         14 damageMax          15 defence           16 absorbtion
17 speed             18 rotationSpeed      19 scanRange         20 protection0
21 protectionFire    22 protectionWater    23 protectionAir     24 protectionEarth
25 protectionAstral  26 fighterSkill0      27 skillBlade        28 skillAxe
29 skillBludgeon     30 skillPike          31 skillShooting     32 mageSkill0
33 skillFire         34 skillWater         35 skillAir          36 skillEarth
37 skillAstral       38 itemLore           39 magicLore         40 creatureLore
41 castSpell         42 teachSpell         43 damage            44 damageFire
45 damageWater       46 damageAir          47 damageEarth       48 damageAstral
49 damageBonus
```

| Kinds | Operand | `Effect+0x40..+0x43` |
|---:|---|---|
| 1..40, 49 | signed-decimal prefix | mode 0/8: i32; modes 1/2/4: observed i16 magnitude + u16 ticks |
| 41 | spell name, optional power | u16 spell id + observed i16 power |
| 42 | spell name, optional mode | u16 spell id + optional u16 ticks |
| 43..48 | first-hyphen `low-high`, optional mode | u8 low + u8 `(high-lowByte)` + optional u16 ticks |

Exact `permanent` and `singleuse` select modes 0 and 8. Otherwise searched substrings `charges`,
`duration` and `continuous` select 4, 1 and 2; their count is stored as `count<<4` in the high u16.
Unknown suffix text is mode 0. Kind 41 always treats its suffix as power. General decimal parsing
keeps a leading run from `+-1234567890`: empty is zero, `12junk` is 12, and a non-empty prefix for
which `%d` makes no assignment leaves an indeterminate value. Out-of-range input remains accepted,
but its concrete CRT result is Unknown. Cast power admits no minus, so `:-20` is zero. Range bytes
wrap and the second is spread, not maximum (`ITEM-EFFMODE-073`).

One Effect is 0x48 bytes: Token base, u8 kind `+0x3c`, u8 mode `+0x3d`, operand union `+0x40`,
transient dword `+0x44`. Effect copy preserves Token state `+0x0c` and transient `+0x44`; equality
compares only kind, the full operand dword and mode. Item equality first compares item code. Both
stackable means equal immediately, exactly one stackable means unequal, and only two non-stackable
items compare their ordered lists in lockstep. A Potion remains stackable with effects because its
item kind is 3 (`ITEM-STACK-003`, `ITEM-EFFOBJ-072`).

For freshly parsed equipment Effects, Token state `+0x0c` is zero. Equip walks the list and
dispatches each value with multiplier +1; unequip walks it with -1. Every kind arm then calls target
`vt+0x50` recompute. Kinds 1..37 target price, actor attributes, pools, combat, protection and skill
fields; 38..40 diagnose not implemented; 41 is a general no-op; 42 teaches a missing Spell; 43 and
49 share damage-base; 44..48 install elemental base/spread. `ITEM-EFFDISP-075` and the experiment's
`effect-keys.tsv` retain the exact per-kind target, clamp, class gate, jump arm and shipped count.
Token-state values 8, 12 and 17 select special lifecycle paths and are not grammar kinds.

Shipped per root: 2,386 cells, 935 non-empty, 260 braced, 266 accepted Effects and one rejected
token. `Humans` contributes 258 braced cells, `Units` two. Corresponding EN/RU raw cells are equal.
The complete population and the one nested-brace anomaly are `ITEM-EFFPOP-071`.

### A weapon's numbers — `FUN_0050dadc`

Every scaled value is `ftol(row.column × shape.f × material.f + 0.5)`, where `shape` is the
5-row quality table at `0x609b2c` subscripted by `item+0x45` and `material` the 16-row table at
`0x609b18` subscripted by `item+0x46` (`0050d923 CMP EAX,0xf` bounds `+0x46`, and
`FUN_004db6b5`, which writes `+0x45`, searches `this + 0x14`).

**The record's double ladder — corrected by `ITEM-LADDER-019`…`ITEM-WEAPCOL-021`; this file
carried it one slot high until 2026-08-03.** Both tables share one schema (`Data.bin` group A: a name plus nine f64) and the
runtime record is **`0x68` bytes** (`FUN_0051b970`, `0051b97a IMUL EAX,EAX,0x68`). The
serialized tail is `0x48`, so the head is `0x20` and double *j* sits at **`+0x20 + 8j`**. Group A
ships eleven titles for one name and nine doubles, so one title has no slot: the correspondence
is *j* ↔ shipped title *j+1*, leaving `Level` unserialized and putting the two identically-zero
doubles on the non-numeric `Abbreviation` and `Materials` columns.

```
disp    +0x20        +0x28      +0x30   +0x38   +0x40     +0x48    +0x50      +0x58         +0x60
title   Abbreviation Materials  Price   weight  @.damage  @.toHit  #.defence  #.absorption  MagCap
```

```
w+0x4a  weight  = col3 × shape[+0x38] × material[+0x38]    title: weight
w+0x60  dmgBase = col6 × shape[+0x40] × material[+0x40]    title: @.damage   <- @.physicalMin
w+0x61  dmgSprd = col7 × shape[+0x40] × material[+0x40] − w+0x60               <- @.physicalMax
w+0x52  toHit   = col8 × shape[+0x48] × material[+0x48]    title: @.toHit
w+0x6a  defence = col9 × shape[+0x50] × material[+0x50]    title: #.defence
w+0x48          =        shape[+0x60] × material[+0x60]    title: MagCap  (no column, no +0.5)
w+0x50  reach   = col0xb verbatim, or 1 when that column is −1 (no factor)
```

`+0x58` (`#.absorption`) is read by the armor and shield fills and by **nothing on the weapon
path** — which is one of the three things that fix the ladder's origin. The other two: nine
doubles starting at `+0x28` would put the last one past the end of a `0x68`-byte record, and the
`+0x30` price factor (`SHOP-PRICE-011`) would then read `0.0000` on all 21 shipped rows.

**`w+0x60` is re-read as an already-rounded unsigned byte before the subtraction**
(`0050db5d MOV DL,byte ptr [ECX + 0x60]`, `0050db66 FSUBP`), so the spread is
`round(max×f) − round(min×f)` and not `round((max−min)×f)`. The two encodings are
indistinguishable while one weapon is equipped and diverge as soon as a second source contributes,
because the actor's fold sums bases and spreads **separately** (`HERO-FOLD-035`).

**A shape factor is never 1.** `Common`'s `@.damage` is `0.2000`, so a name carrying no leading
shape word yields one fifth of the row's columns, and the `Weapons` columns are on a different
scale from every quantity they meet. Worked, identical on both roots: `Iron Short Sword` = shape
`Common` (0.2) × material `Iron` (1.0) over columns `23..40` gives **`(5, 3)`**, `@.attackType` 1,
`+0x52` 5, price 250. `Uncommon Steel Two Handed Sword` = `0.24 × 1.3351 = 0.3204` over `41..82`
gives **`(13, 13)`**.

Names are parsed by `FUN_0050d670`: a leading shape word, then a leading material word, then
the remainder looked up in the class's own collection (`FUN_0051c240(0x609b68, …)` for `Weapons`).
Weapons row 0 is `BareHands` and **no instruction in the mapped image names it**, so it reaches an
actor only through a Units `EquipItem` column.

### An armour's and a shield's numbers — `FUN_0050c53a`, `FUN_0050cfa2` (`ITEM-ARMFILL-032`)

The same routine twice, offset by two bytes because the armour spends `+0x50` on its `Slot`. The
three item fills share the `Data.bin` group C title list, so a column index means the same thing
in all three, and each factor's own group A title names the same attribute:

```
armor+0x50  slot       = col4 verbatim (no factor)              title: Slot
armor+0x52  defence    = col9  x shape[+0x50] x material[+0x50]  title: #.deIrnce  x #.defence
armor+0x54  absorption = col10 x shape[+0x58] x material[+0x58]  title: #.absorbtion x #.absorption
armor+0x48             =         shape[+0x60] x material[+0x60]  title: MagCap (no column)
armor+0x4a  weight     = col3  x shape[+0x38] x material[+0x38]  title: weight

shield+0x50 defence    = col9   ... same factors, same order ...
shield+0x52 absorption = col10  ...
shield+0x48            =        ...
shield+0x4a weight     = col3   ...
```

**The absorption term is the one scaled term in any of the three fills with no `+0.5`.** Every
other scaled value in this file is `ftol(x + 0.5)`; `armor+0x54` and `shield+0x52` are `ftol(x)`
(`0050c674` → `0050c677`, `0050d012` → `0050d015`, against `0050c6a1` and `0050d03f`). It is not
cosmetic: **188 of 3 120** shipped (row x shape x material) combinations differ under the two
rules. The shield fill reads **no** `Slot` column — it has no `PUSH 4` — although all nine
shipped `Shields` rows carry `Slot = 2`.

### What equipping an armour or a shield does (`ITEM-ARMFOLD-033`)

Both classes embed one **0x16-byte modifier block**, `memset` to zero by
`FUN_004fa49b` → `FUN_004fa4b1`, at `armor+0x52` and `shield+0x50`. The two scaled words above
are its members `+0x00` and `+0x02`; **nothing writes the other twelve.** `Armor::Equip` and
`Shield::Equip` pass the block's **address** to `FUN_004fa4eb` twice — once with `actor+0xfe`,
once with `actor+0xbe` — and that routine folds fourteen members (`HERO-FOLD-033`). So:

```
an armour or shield contributes  defence (block +0x00) and absorption (block +0x02)
and contributes exactly zero to  the six protections  actor+0xc2..+0xcc
                                 the six damage kinds actor+0xce..+0xd3
```

There is **no per-slot difference**: the slot byte indexes `actor+0x198 + 4i` and is read nowhere
else, and neither block add sits under a branch on it. The Unit equip wrapper's
`vt+0x54` call (`004f4dbb`) is **not** an equipment recompute. Humanoid/Human
equip/removal wrappers omit that call (`SAV-EQUIPCALL-554`, `ITEM-EQUIP-006`).
The humanoid `+54` target `FUN_004fc22a` sums
`actor+0x1cc + 4i` for `i = 1..5`, the skill-experience caches, and stores `ftol(sum x 0.01)` into
`actor+0x1c`; the base class's slot is a bare `return actor+0x1c`, and the caller discards the
result. Equipment is **stored on the actor, never recomputed from the worn array.**

The helper set does not specify a common event order. Armor removal refreshes
negative weight, subtracts both defensive blocks, clears its slot, sets flags,
then removes Effects. Shield removal removes Effects after those block
subtractions but before flags and slot clear. Shield attach sets flags after
Effects; Armor attach sets them before (`SAV-EQUIPORDER-552`, `ITEM-ARMFOLD-033`).
Weapon attach prepares its Spell before old-weapon eviction and calls derive
before timing/range and weight stores; Weapon removal starts with Effects and
clears the owned Spell and actor slot near the end (`HERO-EQUIP-017`).

Weapon's serialized byte `+50` is the current range operand, not the active
skill selector. After derive, attach adds `u8(Weapon+50)-1` to actor range
byte `+12c`; removal subtracts it, both with byte wrapping. Removal timing
is a literal8/4 assignment. The active melee selector instead uses the low
byte of cached definition parameter5; ranged11/12 and removal assign zero
(`SAV-EQUIPORDER-552`, `SAV-HUMEQUIP-447`). Callback mutation remains Unknown.

Both Effect walks are forward. Each state0 Effect's general normal-return
dispatch derives before the next Effect. One next-node pointer is prefetched;
later node contents remain live, so this is not a pre-collected batch
(`SAV-EQUIPEFFECT-553`). Command22 takes the source object, calls the actor
wrapper, reinserts the displaced/removed object and refreshes zero weight.
Derive from that final refresh is conditional, not an unconditional final pass
(`SAV-EQUIPCALL-554`). Named reads can occur before later local stores; a
direct-body no-store result does not establish transitive callback purity or
an atomic equipment update (`SAV-EQUIPOBS-555`).

### The `Weapons` row's columns

Group C ships eighteen titles for a name and seventeen slots, so `slot i` = shipped title `i+1`
with no spare. Four uses anchor it at instruction level:

```
slot  2  Price          FUN_0050dc4e, item+0x1c
slot  3  weight         the fill's w+0x4a
slot  5  @.attackType   FUN_0050def2 0050df76, branched CMP 0xa / 0xb / 0xc
slot  6  @.physicalMin  the fill's w+0x60
slot  7  @.physicalMax  the fill's w+0x61
slot  8  @.toHit        the fill's w+0x52
slot  9  #.deIrnce      the fill's w+0x6a
slot 0xb @.range        the fill's w+0x50; Weapon::Equip 0050e15c actor+0x12c += w+0x50 − 1
slot 0xc @.charge       Weapon::Equip 0050e126 -> actor+0x134
slot 0xd @.relax        Weapon::Equip 0050e156 -> actor+0x135
slot 0xe 2 handed       Weapon::Equip 0050df3b CMP dword ptr [EAX],0x2, the shield-removal gate
```

`Weapon::Equip` has **three** arms on slot 5. `< 0xa` is melee: it adds `w+0x60`/`w+0x61` into
the modifier pair `actor+0xf4`/`+0xf5`, adds `w+0x6a` into `+0xfe` and `w+0x52` into `+0xe6`,
**assigns** `actor+0xf9/+0xfa/+0xfb` from `w+0x65/+0x66/+0x67`, and sets `actor+0xb6` to the
attack type. `== 0xb` and `== 0xc` add `w+0x60`/`w+0x61` into `actor+0xf9`/`+0xfa` instead and
assign the kind byte `+0xfb` to 1 or 2, leaving `actor+0xb6` zero — so the shipped
`Flame Thrower` feeds the *second* damage component, not the first.

### A weapon's own damage line

The game draws it, and it reads the weapon's own two bytes. `FUN_0047bf40` takes the equipped
weapon at `actor+0x74` and appends two attributes through `FUN_00484af0`, which stores a
`(tag, value)` byte pair at `this+0x0c + 2 × this+0x09`:

```
0047c3b9  MOV CL,byte ptr [EAX + 0x61]   gate: a weapon with zero spread draws no damage line
0047c3c0  MOV AL,byte ptr [EAX + 0x60]   -> attribute tag 0x0d
0047c3d0  MOV DL,byte ptr [ECX + 0x61]   -> attribute tag 0x0e
```

`FUN_00484160` formats the list — `ECX = tag − 1`, bounded `CMP ECX,0x32`, indexed through the
byte table at `0x484954` into the jump table at `0x484930`. The `0x0d` arm at `0048422d` consumes
**both** pairs: it prints the base, steps over the next tag byte, does `00484298 ADD EBX,EAX` and
prints `base + spread`. So the item's line is `[base, base+spread]` — the same composition the
character sheet uses (`HERO-SHEET-038`) — and tag `0x0e`'s own table entry is the default arm,
which nothing reaches. The shipped `Iron Short Sword` draws `5-8`.

## The container

`0x24` bytes: an MFC `CObList` (`+0x00..+0x1b`) plus two dwords.

| Off | Type | Meaning |
|-----|------|---------|
| +0x1c | u32 | the index the **next** `Add` inserts at; constructed to 10 000, i.e. "the end" |
| +0x20 | u32 | the **load** = `Σ (i16 weight × u16 count)` over the elements |

Homes: `actor+0x7c` (`FUN_004f30a2` at `004f3282`), `sack+0x40`. There is **no slot count and
no capacity**: an insert appends whenever `index >= count`, and nothing in the class or in any
arm of the move command compares the count or the load against a maximum.

What the load does, and the only thing it does (`FUN_004f7dfc`):

```
load = actor.ownWeight(+0x8e)
if actor.container:
    if container.load < 64000:  load += container.load / 2      (truncated toward zero)
    else:                       load  = 32000
capacity = body*10 + 1                                          (actor+0x92)
if load >= capacity:  speed -= load / capacity;  speed = max(speed, 6)
```

### Whole-object transfer and serialized state

Whole extraction (`qty >= count`) returns the existing Item pointer; only
`qty < count` reaches split. With count1 and a non-merging destination, the
inspected take/add/list bodies do not store any serialized Item Token or
Position field. This is a bounded direct-body result, not a transitive callback
guarantee. Container draining takes one unit repeatedly, so a larger stack can
still be split. — ITEM-WHOLE-128

A merge requires equal codes and both objects stackable. It retains the
destination Item, adds the incoming count, ORs incoming flags into destination
Token `+0x08`, and deletes the incoming Item. It does not directly copy the
incoming Position or `+0x14` into the retained Item. — ITEM-MERGE-129

Pickup writes each incoming Item `+0x08 = 1` before the drain. A pickup merge
therefore produces `oldDestinationFlags OR 1`. GiveAll and XferItem direct
bodies do not perform that stamp. New-Sack adoption retains the supplied
container; an existing Sack instead drains/deletes it. Sack Position is its
own allocation, copied from the supplied drop Position, not a rewrite of the
contained Item Positions. — ITEM-GROUNDMOVE-130

## Equipment

Fourteen pointer fields, all outside the container.

```
slot 1  -> actor+0x74                       any actor
slot 2  -> actor+0x78                       any actor
slot 3..12 -> actor+0x198 + 4*slot          Humanoid only, gated on actor->vt+0x30()
                                            (false arm: "Error - Trying to takeoff armor
                                             from non humanoid")
```

The UI index in the command is `slot - 1`. **The array is thirteen dwords, `actor+0x198` through
`actor+0x1c8`, and element 0 is dead storage**: the constructor `FUN_004f6ded` clears `0..12`,
while `Humanoid::Serialize`, the death strip and the destructor all run `1..12`, and
`Armor::Equip` refuses a piece whose `Slot` is 0 outright. So no code path can fill index 0 and
none reads it.

**`Armor::Equip` indexes the array with the slot number and does not skip 1 or 2**, so
`actor+0x19c` and `actor+0x1a0` are ordinary slots — written by `Equip`, round-tripped by both
arms of the serializer, stripped on death and deleted by the destructor. The one thing that
cannot reach them is the **take-off command**, which spends those two numbers on `actor+0x74`
and `actor+0x78` instead. No shipped `Armors` row uses either.

**Which slot an item takes is a property of the item, never of the actor or of where the item
was named.**

```
Weapon::Equip  FUN_0050def2  -> actor+0x74            no refusal arm exists
Shield::Equip  FUN_0050d2de  -> actor+0x78            refuses when shield+0x0c == 0
Armor::Equip   FUN_0050c8d0  -> actor+0x198 + 4*part  refuses when armor+0x50 == 0, and when
                                                      actor->vt+0x30() is 0
armor+0x50 = param 4 of the piece's Armors row (FUN_0050c53a 0050c579); shipped title "Slot",
             values 1..12; FUN_0050c53a discards a row above 12 WITHOUT restoring +0x50 to 0
```

`Armors`, `Shields` and `Weapons` share ONE 18-title array, verified identical from all three
collections on both roots. Two of its slots have consumers read at instruction level:

```
slot  4  "Slot"        armor+0x50; which of the twelve fields the piece takes
slot 15  "sutableFor"  a TWO-BIT MASK over consumer classes: bit 0 = fighter, bit 1 = mage.
                       FUN_0050cb03 (Armor), FUN_0050d467 (Shield) and FUN_0050e449 (Weapon)
                       each clear the item descriptor byte element+0x08 and copy bit 0 into
                       its bit 1 and bit 1 into its bit 2; the death gate tests the whole
                       value against 0. Shipped: Weapons 0 on three monster attacks (Boulder
                       Thrower, Flame Thrower, Sonic Beam), 1 on nineteen martial weapons,
                       2 on the two staves, 3 on BareHands and Plasma Sword, and one row
                       (rem) with no parameter array; Shields all 1; Armors 1 on nineteen
                       metal rows, 2 on nine cloth rows, 3 on Amulet and Ring.
```

### Equipment transfer boundary

Armor, Shield and Weapon retain their Item allocation through the named
equip/unequip slot operations. The inspected item-class bodies write actor
slots/stats, not Item Token/Position fields. Effect apply/remove and actor
recompute callbacks are an unresolved transitive write frontier. Base Item
equip is instead consumption: positive actor HP permits changing zero Effect
mode to8, Effect attachment, then Item deletion. — ITEM-EQUIPMOVE-131

Weapon equip finds the first kind41 Effect and, if present, deletes any old
owned Spell and reconstructs one from the low byte of Effect `+0x40`. For a
valid nonzero id, the new Spell's serialized `+0x09`, `+0x0a` and `+0x0c` come
from Spells row parameter6 low byte, parameter18 == 1, and parameter1 low word.
Unequip deletes a nonnull owned Spell and clears Weapon `+0x80`, while leaving
the source Effect in place. No Item split is needed for this nested identity
change. An allocator may reuse the old numerical address. — ITEM-SPELLMOVE-132

### Weapon-borne spell state (`ITEM-CASTSTATE-056`)

A `castSpell` record remains an ordinary ordered Effect of kind 41. General effect dispatch does
nothing for that kind beyond target recompute. A separate list finder returns the first kind-41
record; Weapon construction uses its spell-id byte to create a distinct owned Spell at
`weapon+0x80`, and Weapon copy/split reconstruct it from id. The list record remains the source of
comparison and persistence while the derived object holds cast runtime state (`ITEM-CASTLINK-076`).

A weapon carrying `castSpell` keeps its constructed `Spell` at `weapon+0x80`. Both release routes
temporarily set `actor+0x64` to that `Spell` and `actor+0x68` to the weapon, then clear both after
the attempt. The non-null item field suppresses mana use and the immediate half-mana training
award. The item effect's `+0x40` is the byte spell id and `+0x42` is a raw signed-`i16` power. The
ordinary staff path decrements no charge and preserves the effect record, `weapon+0x80` pointer and
pointed-to object identity. It does rewrite that `Spell` on every apply: `+0x09` is reloaded from Max
Range and raised by `power/3` for Teleport or `power/30` otherwise, while `+0x0e`, `+0x0f` and
`+0x10` receive the damage base, spread and duration scratch. `+0x09` is serialized; the three
scratch fields are not.
The authored Data.bin equipment corpus carries 48 such weapons: 46 Human staffs plus Catapult and
Ballista Unit weapons. Their 18 distinct powers are
`{1,5,10,15,25,30,34,35,40,50,60,63,65,70,82,90,98,99}` on each root; 10 and 20 are not an authored
boundary. Runtime stock is larger: the shop generates castSpell weapons with ids
`{1,11,13,14,20}` and price-derived random power capped at 100. The two siege Units can apply their rider, but their plain Unit vtable uses empty cast,
damage and kill award slots, so they gain no training.
`Unit::Serialize` stores `actor+0x68` as an archive object reference and `actor+0x64` as a raw
`u32`; a save during the temporary interval carries both bit patterns, but only `+0x68` is remapped.
Whether the raw Spell pointer remains valid after any load is untested.

A separate post-cast arm destroys an item of kind `0x0e` and its `Spell`; the shipped staff does not
have that kind. Prismatic Spray id 14 is route-dependent: caster admission applies the fan before
the common wrapper refuses the id, while the fighter rider reaches only that refusal. Training is
not stored on the item: the common sink first requires recipient type id in `[0x21,0x3f]`, after
which later damage awards enter the actor's school or current weapon slot.

The selected-inventory Cast UI cannot admit a Weapon despite accepting kind 2 into its vmethod. It
requires descriptor bits `0x10|0x01`; a castSpell Weapon has bit 4 from its effect but its class
descriptor supplies only suitability bits 1 and 2. Bit 0 is a MagicItems-class bit. The server arm
would accept a crafted weapon item order, but no Weapon row can make the shipped UI emit one.

### The wear rule (`ITEM-WEAR-055`, `ITEM-WEAR-057`)

The column decides whether a character may equip or use the item, and the whole rule is one
13-instruction predicate in the CLIENT. The simulation applies none of it.

```
FUN_00460440(member, element) =
      (element+0x08 & 2) && !(member+0x18c & 2)      item allows fighter, member is no mage
   || (element+0x08 & 4) &&  (member+0x18c & 2)      item allows mage,    member is a mage

member+0x18c bit 1 = mage. FUN_0045f850 writes it from (typeID - 0x21) & 2 for a typeID in
   [0x20, 0x40). In non-zero player-character constructor mode, Human builds typeID =
   sex + 0x23 for a mage, sex + 0x21 otherwise (004f95d5 / 004f95e4), and bit 2 of the
   same dword is sex. Zero-mode map Humans retain their Data.bin typeID; the classifier's
   low-type arm sets the mage bit independently for 0x17 and 0x18 (`PARTY-M20-031`).

element+0x08 = the item descriptor byte the simulation builds and the client reads
   bit 0  0x01  the item is a MagicItems-backed Item, not an equipment piece
   bit 1  0x02  sutableFor bit 0
   bit 2  0x04  sutableFor bit 1
   bit 4  0x10  the item carries an effect of kind 0x29 (castSpell)
   bit 5  0x20  the item has at least one effect; cleared again when item+0x44 == 3 (Potion)

Two call sites, whole image (EnumRefs callto:460440, 2 hits / 2 owners / 0 orphan):
   FUN_004a2fe0 004a3198  shop cell background only: false -> backinvg.bmp, the empty-cell
                          background. The icon is drawn either way and the item is still sold.
   FUN_00492410 004924e3  the equipment doll's drop handler. False sets the refusal flag; the
                          routine returns 0 without reaching this->vt+0x80, the only call that
                          turns the drop into a move. No command 0x22 is produced.

Command 0x22's handler FUN_004d5dd8 and all seven routines its equip destination reaches read
no class field and not this column. Their Data.bin parameter reads are 5, 0xc, 0xd, 0xe.
```

A `MagicItems` item takes its restriction from its ROW NAME instead. `FUN_005092c9` writes
`0x07` into the descriptor byte and clears bit 1 when `item+0x44 == 5`; `item+0x44` is a prefix
test on the row name at `record+0x4` through `FUN_0056f10f` (`strstr`, `-1` for no match):
`Potion` = 3, `Book` = 5, `Scroll` = 4, otherwise 0. Shipped, 49 rows, both roots: 13 Potion,
5 Scroll, 5 Book, 26 neither. The five `Book` rows are the only mage-only `MagicItems`, and
they are refused a second time by the `teachSpell` gate (`ITEM-WEAR-058`).

A refusal returns `this` rather than the displaced item, and every caller treats a non-null
return as "put this in the container" — so a refused piece lands in the backpack. **No slot
restricts by type**: the one runtime-class test on the path, `Humanoid`'s `vt+0x38`
(`FUN_004f705b`, testing against `Armor`), branches between two byte-for-byte identical calls.

One cross-slot rule: `Shield::Equip` reads param `0x0e` of the already-worn weapon's row and,
when it is `2`, takes the weapon off through `actor->vt+0x40` and appends it to the container
itself. `Weapon::Equip` carries the mirror arm for an already-worn shield.

### Equipment named by a `Humans` row

The ten trailing strings of a `Humans` row are ten equipment cells, and the cell **position**
picks the C++ class with the string never consulted:

```
cell 0      -> Weapon  (FUN_004f9065 004f93a7, PUSH 0x84, ctor FUN_0050d670)
cell 1      -> Shield  (                004f9403, PUSH 0x68, ctor FUN_0050ccd4)
cell 2..9   -> Armor   (                          PUSH 0x68, ctor FUN_0050c2ea)
each        -> actor->vt+0x3c(item), which on a Human is FUN_004f7099: it calls
               actor->vt+0x38(item) and appends to actor+0x7c only the RETURN value
```

This is not how the `Units` arm does it — there the class is chosen by searching the string for
the literal `"Shield"` (`UNIT-EQUIP-005`). A name resolves through
`FUN_004db944` (cut from `{`) → `FUN_004db6b5` (a Shapes word) → `FUN_004db801` (a Materials
word) → `FUN_004dba5f` → `FUN_0051c240` (exact lookup, 0 = not found). `FUN_004dba5f` puts back
the shape word the **material** implies — `"Soft "` for a material containing `Leather`,
`"Wooden "` for one containing `Wood` — and the `Weapon` constructor is the one that does not
call it. The `Shield` constructor alone then cuts the literal `" Shield"` off the residue.

```
actor->vt+0x38  Equip(item)    -> item->vt+0x38(actor); actor->vt+0x54(); returns the displaced item
actor->vt+0x40  Unequip(item)  -> item->vt+0x3c(actor); actor->vt+0x54(); returns the item
Item::vt+0x38                  -> refuses if actor.health <= 0; else attaches every effect of
                                  the item to the actor and DELETES the item (a consumable)
Item::vt+0x3c                  -> prints "Unknown item takeoff"; only the three subclasses
                                  implement a real take-off
```

## What an item is called — the display name (`ITEM-DISPNAME-036`…`ITEM-NAMELIMIT-042`)

Nothing in the image composes an item name for display. The name is a **stored line**, selected
by a `CMap<u16, const char*>` built once at start-up from two `main.res` nodes
(`ITEM-DISPNAME-036`):

```
main\text\itemname.bin   N u16 keys, little-endian, no header   (shipped: 832 B = 416 keys)
main\text\itemname.txt   N lines, one display name each          (shipped: 416 lines)
map[key[i]] = line[i]     FUN_00468380, count = filesize(bin)/2
```

The key is the **packed item code** — the same `u16` an authored `.alm` element carries
(`ITEM-CODE-029`) — and its fields are fixed by the diagnostic formatter `FUN_00483c80`
(`ITEM-NAMEKEY-037`):

```
bits 12..15  material   index into Materials      (item+0x46)
bits  8..11  class      1 Weapons, 2 Shields, 3..13 Armors, 14 MagicItems
                        for an Armors row the value is also the row's own Slot column
bits  5.. 7  shape      index into Shapes                       (item+0x45)
bits  0.. 4  row        1-based index into the class's collection (item+0x0c)

class 14 only: bits 0..7 are the whole MagicItems index and there is no shape field
```

Reading a name is `FUN_00483e50`: take the interface element's `u16` at `+0x6`, `Lookup`. On a
**miss** the element is not drawn under a fallback — it is dropped from the list and a
diagnostic naming the code is emitted, so the name table is the inventory's admission list
(`ITEM-NAMEMISS-039`). The numeric form the diagnostic carries is seven digits with no
separator: `%02d%02d%1d%02d`, or `%02d%02d%03d` for class 14.

Consequences a consumer must not miss:

- **The name is free text, not a composition.** `Bronze`/`Amulet`/`Common` reads `Beard`; the
  `Elven Adamantium Amulet` key reads `Adamantium Amulet` in English with the tier word absent
  and carries a tier word in Russian. No grammar produces both (`ITEM-NAMEPOP-038`).
- **Localisation lives here and only here.** `Data.bin`'s own item row names are byte-identical
  on the two roots (`DAT-ITEMNAME-010`); `itemname.txt` differs on all 416 lines and
  `itemname.bin` not at all (`TEXT-ITEMNAME-019`).
- **Enchantment is invisible to the name.** The key has four fields and none is an effect, so an
  enchanted piece shows the same line as a plain one (`ITEM-BRACE-041`).
- **416 of 6 064 expressible items have a name.** The rest cannot be shown at all.
- The `[tier ][material ]shape` grammar (`UNIT-EQUIP-005`) is the **parse** direction only: it is
  how an authored string in a `Units.EquipItem` cell, a `Humans` equipment cell or a mission
  `.ini` becomes an item, and it is finished before anything is drawn (`ITEM-NAMEPARSE-040`).

## What an item looks like — `item+0x40` (`ITEM-APPEAR-023`, `ITEM-APPEAR-024`)

Every item carries one `u16` at `+0x40` that is the whole of its appearance. It is not stored: five
routines write it and all five write the result of `FUN_00525d60`, which is

```
item+0x40 = ((item+0x46 & 0xff) << 12)      the Shapes/Materials index -> material.reg
          | ((kind      & 0xff) <<  8)      the equipment slot, 1..12
          | ((item+0x45 & 0xff) <<  5)      3 bits, unnamed
          |  (item+0x0c & 0xff)             the Data.bin definition row index
```

`kind` is the literal `1` in the `Weapon` builder (which also stores it at `item+0x50`), `2` in
`FUN_0050cfa2`, `1` in `FUN_0050d5d4`, `byte[item+0x50]` in `FUN_0050c53a`, and the caller's
argument in the base `Item` builder. The `OR` is not a shift into disjoint fields, so a `+0x0c`
above 31 would corrupt the `+0x45` field.

`FUN_00483c80` turns the word into a seven-digit name — `"%02d%02d%1d%02d"` of the four fields, or
`"%02d%02d%03d"` of the top two and the whole low byte when the slot field is 14 — and that name
addresses both `graphics\inventory\<name>.16a`, the icon, and
`graphics\equipment\<figure>\<layer>\<name>.256`, the figure layer. `FUN_00483d70` tries to open
the icon and returns 1 when it cannot; the `0x76` equipment message then shows the engine's own
`"Invalid item weared "` for five seconds. *Weared* is the engine's word for the equipment slots.

**`FUN_00483c80` is not a method on `Item`.** Its `this` is a display record whose `+0x06` is a
copy of `item+0x40`, made by the record's only two constructor call sites:

```
0047c368  MOV DX,word ptr [ECX + 0x40]     0047c458  MOV AX,word ptr [EDX + 0x40]
0047c36f  CALL FUN_00484ac0                0047c465  CALL FUN_00484ac0
00484ac8  MOV word ptr [ESI + 6],AX        00484ad3  MOV byte ptr [ESI + 0xa],AL
```

So a consumer needs two fields, not one: the simulation's `item+0x40` and the panel record's own
copy. The record also carries `+0x04` u16, `+0x08` u8, `+0x09` u8, `+0x10` and `+0x14` dwords
(`FUN_00484990`, its copy constructor); what those are is unread.

**The picture word is rebuilt at construction, not carried.** `FUN_004dd02a` passes only three
values into each constructor — shape, material, row — and the kind nibble comes from the class:
the literal 1 for a `Weapon`, 2 for a `Shield`, the caller's argument for the base `Item`, and,
for an `Armor`, `byte[item+0x50]`, which the fill has just taken from the row's own `Slot`
column. So an authored code whose kind nibble disagrees with its `Armors` row's `Slot` allocates
by the nibble and **draws by the `Slot`**. On the shipped corpus they never disagree
(`ITEM-PICT-051`).

**A second self-check message.** Beside `"Invalid item weared "` sits `"Invalid item in
inventory "` (`0x5b83ac`), shown for **10 000 ms** at `004134e7`, gated on `FUN_00483e50` — which
is not a file open but a keyed lookup of the same word in the object at `0x5eb410`, and returns
**non-zero on success**, the opposite polarity to `FUN_00483d70`'s. `graphics\inventory\` is
composed at five sites in all: `00481f27`, `004838c0`, `00483db1` (the check), `00491815`,
`004a3571`. Inside `graphics.res` the nodes are under `inventory/`, not `graphics/inventory/` —
the composed path's first segment names the container (`ITEM-PICT-049`).

### The registry — which items have a picture (`ITEM-PICT-046`…`ITEM-PICT-052`)

Over the whole constructible population (66 shape/material rows x 16 x 5, plus 49 `MagicItems`):

```
addressable items 5329   icons that exist 416   missing 4913        (7.8 % drawn)
inventory nodes    416   claimed by an item 416  unclaimed 0  ambiguous 0
equipment leaves   367 distinct over 928 sheets  = 416 - 49, the class-14 rows
requests the shipped data makes:  181/181 + 177/177 .alm codes, 908/908 Humans cells, 0 misses
```

The name space is a bijection onto the 65 536 `u16` values, and 416 of them are used. The
5 329 is the *addressing space*, not a hole: which of its 80 combinations a row may actually
take is a **five-word bitmask on the row itself**, the ten raw bytes of a group-C `Data.bin`
entry at runtime `entry+0x1c`, one `u16` per `Shapes` row with one bit per `Materials` row:

```
0050966d  MOV DX,word ptr [EAX + ECX*0x2 + 0x1c]     row->mask[tier]
00509672  MOV EAX,0x1        0050967a  SHL EAX,CL     1 << material
0050967c  AND EDX,EAX        0050967e  TEST EDX,EDX   00509680  JZ <next material>
```

Bit set and picture present agree on **5 280 of 5 280** (row, tier, material) cells with **0**
disagreements, on both roots — 367 bits, 367 pictures, plus the 49 magic items outside the
scheme. `Bone` is set on no row at any tier, which is why it has no art. So the map is total
over what the data permits, and the mask is where a customisation widens it: ten bytes on one
row plus one node in `graphics.res`, changing no other shipped file (`ITEM-PICT-050`,
`DAT-MATMASK-020`).

**Reading a name.** Absent words do not default alike: `FUN_004db6b5` returns `0` (`Common`) for
a missing `Shapes` word (`004db7ef XOR AL,AL`) and `FUN_004db801` returns `15` (`None`) for a
missing `Materials` word (`004db931 MOV AL,0x0f`). Both scan their collection downwards from the
last row. Taking the material default as 0 makes 106 of the 908 shipped `Humans` cells
undrawable, all of them cloth (`ITEM-PICT-048`).

What the word is used for on the drawing side is `formats/hero/format.md` §6b.

Claims: `ITEM-APPEAR-023`, `ITEM-APPEAR-024`, `ITEM-PICT-046`…`ITEM-PICT-051`.

## Moving an item — command `0x22`

One command, five fields (`FUN_0041c98f`, `SHOP-TRAY-025`): `cmd+0x0c` source code,
`cmd+0x0e` source index, `cmd+0x0d` destination code, `cmd+0x10` destination position,
`cmd+0x12` quantity.

| Code | Source | Destination |
|------|--------|-------------|
| 1 | an equipment slot (`slot = cmd+0x0e + 1`) | equip; the displaced item goes to `cmd+0x0e` of the container |
| 2 | the actor's container, index `cmd+0x0e`, quantity `cmd+0x12` | the actor's container at index `cmd+0x10` |
| 3 | **the ground** — validates a sack on the actor's own cell, sets `container+0x1c = cmd+0x10` and `actor+0x50 = 2`; refuses with `"Invalid pickup order - no sack there."` | **the ground** — see below |
| 4 | a shop (`4..8` selects tray or shelf) | a shop; stamps `item+0x14 = Player` |

Dropping to the ground:

```
destX = cmd+0x10 & 0xff ; destY = (cmd+0x10 >> 8) & 0xff
if |actorCol - destX| <= 2 && |actorRow - destY| <= 2 :  drop at (destX,destY)
else:                                                    drop at the actor's own cell
```

Never refused — out of range it lands underfoot. Money is opcode `0x23` and uses the same
window; it refuses an amount `<= 0` or one exceeding `Player+0x38`, debits `Player+0x38`, and
adds the amount to the sack's `+0x3c`. CastSpell item use is orders `0x25`/`0x26`;
ordinary Potion use is session `0x22` destination1 (`ITEM-USE-113`).

### Carried Potion and Scroll activation

`ITEM-USE-112`, `ITEM-USE-113` and `ITEM-USE-114` distinguish input, reservation and commit. Mission backpack
double-click and dropping a carried item on the character panel reach the same
action. The shop backpack's double-click uses its borrowed character panel too.
The action requires one owned selected member and descriptor/class admission.
Ordinary Potion and Scroll descriptors admit both classes.

Without castSpell display bit 0x10, usable bit 1 emits a one-unit transfer to equip.
Base Item use returns the item unchanged at HP<=0. At HP>0 it applies every Effect
in order and destroys the detached item, even if healing was already capped or an
effect's class check made it a no-op. With both display bits it instead arms cast
mode and stores the carried slot. Arming does not remove the server item.

Map-up uses the item-specific target table, not simply the Spells target column.
It emits unit order `0x25` or point order `0x26`. The server removes one item and reserves
it at actor+68, building a Spell from its first kind-41 Effect. Point admission
requires Spells parameter 4 ==2 and otherwise restores at the original slot.
Missing/wrong first Effect has no local restore; no ordinary malformed-input
route is established.

Accepted cast-start sets actor+136=0; completion destroys the category-14 Item
after invoking apply, without an effect-success refund. Earlier cancellation can
restore it only under the act/+136/first-effect/id-match checks; cancellation
after accepted start and actor teardown can destroy it instead. Pathfinding and
stale-target failures are not all closed. Shop arming is established, but casting
while the town/shop covers the target map remains Unknown.

`ITEM-VALUE-115`: the general Book value reader uses only its first Effect's
spell id, or zero if empty. It does not select the last Effect. Scroll has a
separate kind41 summation. Shelf constructors override the stored price, and the
Item copy preserves that price and ordered Effect data without re-pricing.

## The sack

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
cell word. Dynamic bit0 rejects before lookup. An existing nonzero Sack slot
rejects even when it already holds the supplied Sack. An existing empty slot
is written without plane recomputation; creation instead zeroes52 bytes,
captures current Cost/Static, sets the record-present flag, stores the Sack
and recomputes. — SAV-SACKENTRY-590

Local removal refuses only a missing record before its write path. It clears
payload+10 without an occupied-slot or pointer-equality test, recomputes, and
tests the four occupant slots, layer count+02 and operation+2c for deletion.
Other residue and the layer pointers themselves do not retain the record.
— SAV-SACKREMOVE-591

Deletion restores Cost and Static baselines, preserving Static bit4, but does
not replace Dynamic with restored Static or clear Dynamic bit5. Dynamic keeps
the preceding recompute result, with the conditional bit4 OR. This is a local
write-set, not proof that allocator or higher-level callbacks preserve every
other field. — SAV-SACKPLANES-592

The lookup used by merge/create first requires Static bit5 and then a matching
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

## Where a sack comes from

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

All four Dragon rows, 112 through 115, author `Flame Thrower` in slot 0, leave slot 1 empty and have
no `castSpell`. `Flame Thrower` has `sutableFor = 0`; the death gate therefore does not move it into
the container. The 28 campaign maps contain 21 Dragon placements per root, split 2/6/4/9 over the
four rows. The complete producer join finds no Dragon staff or enchanted-item death producer: no
type-8 stock element, same-cell ground element, item/container script target or death-time factory
call in the shared death body or enumerated direct/reference population. Full-root placement counts
are 135 EN and 37 RU because the loose-map sets differ. The negative is bounded to those enumerated
families; Flame Thrower is innate-weapon state, not item evidence (`ITEM-DRAGON-088`,
`ITEM-DRAGDROP-089`).

## Magic-item value

`FUN_00508486` initially copies an ordinary `MagicItems` row's parameter 0 to `item+0x1c` as a
signed `i32`, without a sign branch at assignment. Scroll and Book are name-derived exceptions:
Scroll sums spell-derived kind-41 values and falls back to parameter 0 only when the sum is zero,
while Book reads only the first Effect's spell value, or zero when empty (`ITEM-VALUE-115`). Mission
40's `Quest Item31` therefore leaves construction at 10,000 from its row, not from an enchantment.
The later descriptor path does treat `item+0x1c == -1` specially; the assignment rule is not a
global denial of sentinel semantics (`ITEM-MAGVAL-090`, `ITEM-WEAR-058`, `ALM-M40-074`).

## Persistence

`Effect::Serialize` writes the 37-byte Token head, u8 kind, u8 mode, u32 operand and u8 Token
state `+0x0c`: 44 bytes of body. It omits transient Effect `+0x44`. The effect-list serializer
writes u32 count followed by archive object references in list order; load clears the list and
appends that count in archive order (`ITEM-EFFSAVE-077`).

```
Item::Serialize      head, effects, u16 +0x40, u16 +0x42, u8 +0x44, u8 +0x45, u8 +0x46,
                     u16 +0x48, u16 +0x4a, u8 +0x47; on load, resolve +0x3c from +0x0c
Armor/Shield/Weapon  Item::Serialize, then their own tail, then the same resolve against
                     their own collection
container            the CObList elements, then +0x1c, then +0x20
Sack::Serialize      head, +0x3c, container
Unit::Serialize      ... actor+0x74, actor+0x78, container(actor+0x7c) ...
Humanoid::Serialize  Unit::Serialize, raw 24 bytes at +0x1cc, the twelve slots
                     actor+0x198+4i (i = 1..12), one reference at +0x1e4
```

Item copy/split preserves the deep Effect list but leaves the serialized `Item+0x47` byte
indeterminate. Saving such a clone writes that byte; its meaning remains unknown.

## What crosses a mission boundary

`PARTY-CARRY-005`'s carried blob is a single `WriteObject` on one actor. That invokes the
chain above, so **a hero's equipment and backpack cross inside it**; `PARTY-LOSS-006`'s strip
list touches none of `+0x74`, `+0x78`, `+0x7c` or `+0x198…+0x1e4`. A mercenary is never a root
of that graph and is not reachable from one, so nothing of it crosses but `MERC-DEATH-006`'s
fifteen pool integers — and if it dies first its gear is destroyed by the `"NPC"` test above.
A consumer must implement two persistences for one party.
## Document access item

`MagicItems[28]` is internal row `Quest Documents`, packed code `0x0e1c`. The display name is
`Valuable Documents` in EN and `Официальные Документы` in RU. Its `Data.bin` price is -1.

An eligible inventory action checks `(code & 0x0f00) == 0x0e00` and `(code & 0x1f) == 28`, then
posts campaign message `0x463`. The campaign arm opens the document panel. The panel binds to the
campaign record at `campaign+0x548` and reads its document collection. The item is an access trigger;
it does not contain the document entries.

Player-hero construction calls `FUN_004d3f6b` after starting-skill and derived-stat application and
before runtime-id assignment. When `server+0x0c == 0` and `[0x005eb5a4] != 2`, the helper constructs
an Item from the exact row name `Quest Documents` and appends it to the hero's `+0x7c` container.
The campaign construction path establishes the first gate. The second gate was not measured in a
fresh campaign, so actual item presence at mission 10 start is Unknown (`ITEM-DOC-054`,
`ITEM-DOC-069`). The campaign record already contains text documents 1, 2, and 3; a non-empty
collection does not establish item possession or panel access.

The class-14 code uses the low byte as its row, but this UI test uses only five low bits. Rows
`28 + 32n` therefore alias at the gate. The existing packed representation permits at most 255
nonzero class-14 rows. Removing the aliases changes code. Adding a starting access item through a
runtime rule changes no shipped asset; adding one to `Data.bin` or the name resources changes those
files. The producer, destination and both producer gates are executable code; this experiment found
no data-only G2 seam for changing them.
