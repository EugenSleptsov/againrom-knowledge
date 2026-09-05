# UNIT — the non-hero actor: template → instance → combat inputs (partial)

Level 3. Promoted, evidence-backed claims only. Ledger:
[`claims/unit.md`](../../claims/unit.md).

**Status: partial (◐).** Creation is specified end to end for the arm that reaches the
`Data.bin` **Units** collection: the constructor's defaults, the 38 streamed slots, the
equipment step, the spellbook step, the one-time modifier fold, and every combat input
`claims/hero.md` established. Not specified: what `actor+0x14` is, what
`FUN_0050d670` does with a tier/material name, and the six unread constructor call sites.

## What this covers, and what owns the rest

### Physical orders against a Building

The received object-target attack route can resolve a Building id and retain its pointer
through the shared order and approach machine; there is no separate structure action in
the read route. The scorer, caster alternatives and normal approach predicates still
apply. Successful eligibility of every structure class is not established
(`UNIT-STRUCTORDER-062`).

The base Building vtable returns token size 1, independently of its registered rectangle.
Approach checks facing and edge distance, while application uses the separately rounded
strike distance. Start computes range-dependent delay without rejecting range; application
does reject it (`UNIT-STRUCTREACH-063`).

For live combat block `A=attacker+0xa6`, Building damage is zero when the block is null,
maximum HP is zero, or `u8(A+0x14)==0`. Otherwise it is
`max(u8(A+0x13) + U[0,u8(A+0x14)] - 5, 0)`. Ordinary physical damage, to-hit, defence,
absorption and elemental-kind selection are not read by this resolver. Consequently an
accepted physical attack need not damage a structure (`UNIT-STRUCTDAMAGE-064`).

Both short- and long-reach physical presentation remain in the actor's countdown, with
distance extra at start. Physical application subtracts from word HP without a zero clamp;
the separate effect consumer clamps to zero. Retained unredirected orders inherit the
shared two boundary ticks as well as charge/recovery and their modifiers
(`UNIT-STRUCTDELIVERY-065`).

Target HP zero does not cancel the named strike body. The resolver's same-cell cleanup
before a lethal result is not the Building destructor. Actual destruction scheduling,
retained-pointer lifetime, obstruction release and mid-attack save/load remain Unknown
(`UNIT-STRUCTSTOP-066`).

The reached notification builds a Building state message (`0x82`) containing HP,
not a death-specific message. Its admitted client arm updates drawable HP; neither
local arm branches on HP sign. Packet virtuals, send and flush descendants remain
an unresolved execution boundary (`UNIT-STRUCTCONT-078`).

A separate `+0x14` callback increments HP for selector byte 15 or 16 when the
session counter is divisible by 60 and `0 <= HP < maximum`. Zero is admitted;
negative HP is not. Its live scheduling and membership are not established
(`UNIT-STRUCTZERO-080`). A retained cell reference and mask keep the same lookup
and movement interpretation regardless of HP (`UNIT-STRUCTNEXT-079`). The bounded
caller search does not establish whether those inputs survive an actual lethal
application; synchronous/deferred removal and permanent registration remain
alternatives (`UNIT-STRUCTBOUND-081`).

| | owner |
|---|---|
| the `Data.bin` file, its grammar and its schema law | [`formats/databin`](../databin/format.md) |
| the **human** arm — chargen, the derive `FUN_004f7dfc`, the damage resolver | [`claims/hero.md`](../../claims/hero.md) |
| the client drawable's frame, phases and corpse stage | [`formats/anim`](../anim/format.md) |
| movement, pathing, tick order | [`formats/move`](../move/format.md) |
| **this file** — everything between a Units row and a monster that can hit you | — |

## Which file carries which number

`units.reg` is the **client drawable's** class record (`REG-UNITS-049`) and the simulation
module never reads that array (`REG-UNITS-061`). `Data.bin` **Units** is the **simulation
actor's**. The join is `units.reg ID == Data.bin typeID`; the tier `face` (1..4) exists only
on the `Data.bin` side, so one `units.reg` class covers a whole tier family.

Three pairs look like duplicates and are not:

| looks the same | drawable (`units.reg`) | actor (`Data.bin`) |
|---|---|---|
| footprint | `TileSize` → `CUnit vt+0x20` | `tokenSize` → `actor+0x49` (`TERR-MOVE-054`) |
| attack timing | `AttackDelay` / `ShootDelay` | `attackChargeTime` / `attackRelaxTime` → `+0x134`/`+0x135` |
| death | `Dying` / `DyingPhases` — a corpse **sheet** | `dyingTime` — how long the corpse holds its cells |

## Creation, in order

1. **Allocate `0x198` bytes and run a family constructor.** All five install vtable
   `0x59c3c0`; `FUN_004f2cb8` also builds the four blocks: live `+0xa6` (0x16 B), live
   `+0xbe` (0x16 B), modifier `+0xd4` (0x40 B), base `+0x114`.
2. **`FUN_004f30a2`** sets the container fields — `+0x3c = 0`, `+0x0e = 0`, `+0x49 = 1`,
   `+0x4a = 1`, `+0x4b = 1`, `+0x4c = 0` — allocates the mover (`+0x154`), the order block
   (`+0x158`) and the inventory (`+0x7c`), then calls `FUN_004f3317`:

   ```
   body 30  reaction 30  mind 20  spirit 20     speed 10   rotation 8 (mover+0x0a)
   scanRange 5  sight 0   capacity = body*10    health 30 / healthMax 30
   healthRegenPeriod 100   manaMax 0 / mana 0 / +0xa0 = manaMax   manaRegenPeriod 50
   attackCharge 8   attackRelax 4   reach 1
   ```

   **Every `−1` cell in the template leaves exactly these values.**
3. **Find the class by name** (`FUN_004f59de`): the Units collection is walked from index
   26, jumping 28 → 63, which is exactly the shipped non-empty band. On a match,
   `actor+0x0c` = the index and `actor+0x3c` = the entry.
4. **Stream 38 slots** (`FUN_004f5604`) — the map below.
5. **Equip the `EquipItem` strings** (up to two): a name containing `"Shield"` becomes a
   Shield, otherwise a Weapon; each goes through `actor->vt+0x3c`. The name is
   `[<Shapes name> ][<Materials name> ]<Weapons name>`.
6. **Build the spellbook** if `Spell 1 > 0`: `actor+0x4c |= 2`, a `Spellbook` at
   `actor+0x140`, and per pair `[actor+0x158]+0x78+4i` = spell id,
   `[actor+0x158]+0x84+4i` = probability × `0x147`.
7. **Special arms:** `face == 4` sets the five live skill words to 30; `typeID ∈ {0x47,
   0x48}` (Dragon, Daemon) ORs `actor+0x4c |= 6` and forces a spellbook.
8. **Fold the modifier block once** — `FUN_004f54f8(actor+0xd4, actor)` — onto the streamed
   values, with **no zeroing first**. This is the one structural difference from the hero,
   whose fold runs inside a derive that has just `memset` the `+0xbe` block.
9. `XPvalue` re-read into `actor+0x1c`; the mover is registered.
10. **The scenario-setting adjustment**, applied by the `.alm` spawner `FUN_004e26bb` to the actor
    it has just built — see §*The setting that scales a unit*. **Non-heroes only**: the humans arm
    of the same routine jumps over the block (`004e2acf`).

## The slot map (Units)

Slot `i` is column title `i+1`. Store width is the helper's: `u16` `FUN_005233c0`/`005234b0`,
`u8` `FUN_00523410`, `u32` `FUN_00523460`. **Every helper skips its store when the value is
`−1` and advances the cursor anyway.**

```
 0 body            -> +0x84 u16      19..23 prot Fire..Astral -> +0xc4 +0xc6 +0xc8 +0xca +0xcc u16
 1 reaction        -> +0x86 u16      24..28 res.Blade..Shooting-> +0xcf +0xd0 +0xd1 +0xd2 +0xd3 u8
 2 mind            -> +0x88 u16      29 typeID              -> +0x0e  u16
 3 spirit          -> +0x8a u16      30 face                -> +0x4b  u8
 4 healthMax       -> +0x96 (+0x94)  31 tokenSize           -> +0x49  u8
 5 HP regen period -> +0x98 u16      32 movementType        -> +0x4a  u8
 6 manaMax         -> +0x9c (+0x9a, +0xa0)
 7 MP regen period -> +0x9e u16      33 dyingTime           -> dropped; re-read by vt+0x6c
 8 speed           -> +0x8c u16      34 Withdraw            -> [+0x158]+0x40 u32
 9 rotationSpeed   -> [+0x154]+0x0a  35 Wimpy               -> [+0x158]+0x44 u32
10 scanRange       -> +0xa5 u8       36 See invisible       -> [+0x158]+0x71 u8
11 physicalMin     -> local          37 XPvalue             -> +0x1c  u32
12 physicalMax     -> local; spread = max - min
13 attackKind      -> local PRE-SET TO 0; routes the pair:
                        <= 0 (and -1) -> +0xb4 base / +0xb5 spread
                          1           -> +0xb7 / +0xb8
                          2           -> +0xb9 / +0xba, both negated
                          3           -> +0xb4 / +0xb5 and actor+0x4c |= 0x10 (auto-hit)
14 toHit           -> +0xa6 (then +0xa8 <- +0xa6)
15 defence         -> +0xbe u16      slots 38..54 (treasure, Power, Spell/Probability,
16 absorbtion      -> +0xc0 u16      Spell Power) are NOT streamed; other routines read them
17 attackChargeTime-> +0x134 u8
18 attackRelaxTime -> +0x135 u8
```

The **Humans** streamer `FUN_004f974d` consumes 23 slots and differs as `UNIT-DIFF-002`
records — six skills instead of a damage pair, `toHit = Skill.General`, and no regeneration,
absorption, protection or resistance column at all.

## There is no derive

`vt+0x50` on `0x59c3c0` is `FUN_004f5946`, whole:

```
if (actor+0x14) actor+0xa0 = (u16)actor+0x9c * [actor+0x14]+0x58 / 100
else            actor+0xa0 = actor+0x9c
```

`actor+0xa0` is the mana floor the heal AI `FUN_0052eb60` will not cast below. Nothing else
is recomputed, ever. It is invoked whenever `load/capacity` changes bucket
(`FUN_004f36f7` @`004f37b5`), which every equip path reaches.

## The combat inputs

| input | non-hero source | moved by an item? |
|---|---|---|
| damage `(base, spread)` `+0xb4`/`+0xb5` | `physicalMin`, `physicalMax − physicalMin`, via `attackKind` | yes, through the fold |
| absorption `+0xc0` | the `absorbtion` column | yes (armour/shield) |
| `target+0xc6` | the `prot Water` column | yes |
| five protections `+0xc4…+0xcc` | the `prot *` columns | yes |
| damage-kind resistance `+0xce + k` | `res.*` fill `+0xcf…+0xd3`; **`+0xce` (k = 0) is filled by nothing** | no |
| the attacker's `k` = `actor+0xb6` | 0 unless a **melee** weapon sets it to `@.attackType` | yes |
| reach `+0x12c` | **1**, then `+= @.range − 1` per weapon — `@.range = −1` stores 1, so a `Pike` or `Long Sword` stays at 1, and reach is **independent** of the melee/projectile arm | yes |
| attack period `+0x134`/`+0x135` | the two template columns, **assigned over** by `@.charge`/`@.relax` | yes, on 16/56 shipped classes |
| speed `+0x8c` | the `speed` column (ctor 10) — no `reaction` formula on this arm | via the fold's `+0xd8` |
| healthMax `+0x96` | the `healthMax` column (ctor 30) — no XP term on this arm | via the fold's `+0xdc` |

## Corpus, both shipped roots (EN and RU `world.res`, different files)

```
56 parameterised Units rows, param arrays all 55 long
attackKind:  -1 x48, 3 x8 (Bat_Sonic x4, Bee x4) -- never 1 or 2
EquipItem:   30 classes none, 26 exactly one, 0 two, 0 containing "Shield"
             26/26 parse as [tier ][material ]weapon, residue 0
reach:       histogram {1:38 4:8 5:4 8:4 20:2} -- 38 at reach 1, of which 30 are unarmed and 8
             carry a weapon whose @.range cell is -1 (Pike, Long Sword)
             of the 18 above 1, 14 take the melee arm (@.attackType < 10), 4 the projectile arm
cadence:     weapon moves the template pair on 16 classes
spells:      12 classes carry a nonzero Spell 1
prototypes:  53 of 56 inside FUN_00477b10's gate, 0 (typeID, face) collisions
Units differ between the two roots on 0 of 56x55 parameters
```

## The owner, `actor+0x14`

Set by the ALM spawner from the placement's type-5 group id through the player manager
`[0x00609544]` (`004e274f`…`004e2f16`). It is a **`Player`** (`0x70` B, ctor `FUN_004faa81`).
Three of its fields are read outside the shop:

| field | meaning | who reads it |
|---|---|---|
| `+0x28` | **0 = a human participant owns this unit**; 1 = a scenario-authored owner (also the constructor default); 2 = a group whose `HumanFriend` key is `"Yes"` | 25 sites, incl. the route budget (`MOVE-TERM-003`) and the experience payout (`HERO-KILL-027`) |
| `+0x38` | money (`SHOP-BUY-009`); defaulted to 100 on the join path | the shop |
| `+0x58` | `95`, the percentage `vt+0x50` scales `manaMax` by into `actor+0xa0` | `FUN_004f5946`, `FUN_004f7dfc` |

Two consequences a consumer must carry: a unit owned by a human participant gets the **flat
1000-generation** static route budget instead of `max(scalar, D>>2) + D`, and **killing** such a
unit pays the killer **no experience** (`004f7910`…`004f7920`).

## The setting that scales a unit

`[0x005cd758]+0x84` is a dword on the server singleton with **exactly three values**. A consumer
that instantiates a map must carry it; a consumer that resolves a hit need not, because the state
is read at spawn and thereafter lives in the actor's own stat fields.

| `+0x84` | what the spawner does to the actor |
|---|---|
| 1 | `actor+0x96 := ftol(healthMax × 0.66)`, then `actor+0x94 := actor+0x96` |
| 2 | nothing |
| 3 | `actor+0xa6 += 50`, `actor+0xbe += 50`, `actor+0x96 := ftol(healthMax × 1.5)`, then `actor+0x94 := actor+0x96` |

`0.66` and `1.5` are the doubles at `0x0059bc68` / `0x0059bc70` (the display uses a second pair with
the same values). The spawner's second guard is `server+0x0c == 0`, the single-player condition.

**Where the value comes from.** The constructor defaults it to `2`; three sites in the campaign
module set it to `campaign+0x65c + 1`; the save's load arm restores it but **only when
`1 <= v <= 3`**, which is where the value set is stated by the image rather than inferred.
`campaign+0x65c` is the selected index of a **three-button control on the character pre-create
screen**, drawn from `graphics\interface\chrgen\PreCreate\Levels\level0..2`. The control carries
no text, so what it is *called* is not established (`UNIT-GATE-012`…`014`).

## What the placement record overrides

The `.alm` spawner `FUN_004e26bb` applies one more block to the actor it has just built,
after the difficulty adjustment above and after its own re-derivation call
(`004e2d28 CALL dword ptr [EAX + 0x50]`). Every value in it comes from the map's type-6
placement record, and every one is guarded, so an unauthored record changes nothing
(`UNIT-PLACE-034`).

| runtime record | actor | absent value |
|---|---|---|
| `+0x28`, `+0x2b`, `+0x29`, `+0x2a` | `+0x84` Body, `+0x86` Reaction, `+0x88` Mind, `+0x8a` Spirit | `0` |
| `+0x20` u16 | `+0x94` current health | `-1` |
| `+0x24` u16 | `+0x9a` current mana | `-1` |
| `+0x2d`, `+0x2e` | `+0xbe`, `+0xc0` | `0` |
| `+0x37 + i`, `i = 0..4` | `+0xc4 + 2i`, the five elemental protections | `0` |
| `+0x31 + i`, `i = 1..5` | `+0xa8 + 2i`, five of the six skill words | `0` |

The two byte-run loops differ in one place: the protection loop initialises its index to
**0** and the skill loop to **1**. So all five protection bytes the record carries are
applied and only five of the six skill bytes are — `Skill.General`, slot 0 at `+0xa8`, is
the one a placement record cannot set, although the record carries a byte for it and three
shipped placements author it (`UNIT-PLACESKILL-086`, `UNIT-PLACERESIST-087`). Because the
block runs after the re-derivation, it is the last word inside the spawner: on a non-hero
the `spirit / 2` protection fill and its `+70` clamp have already run.

**In the shipped campaign the block is almost inert** (`UNIT-PLACEIDLE-088`). Over the
12 085 placements of both preserved roots the mana word is absent on every one, `+0xbe`'s
source byte is zero on every one, `+0xc0`'s is nonzero on one, and the busiest field in
the whole tail — the sixth skill byte — moves 119. Counted per record over the slots the
spawner actually reads, the block changes nothing on 7919 of the EN root's 8094
placements and 3839 of the RU root's 3991.

## What the information display is given

`FUN_0045f850` (`CUnit`, one caller) copies **25 stores** off the prototype actor in one fixed
order (`UNIT-PANEL-010`). In order: the four stats, health and
its maximum, mana and its maximum, toHit, defence, absorption, the damage pair, speed, sight, the
five damage-kind resistances, the five elemental protections, and a derived byte
(`2` when the typeID is `0x49`, else `0`).

**It reads no `units.reg` field**, so `InfoPicture` and `DescText` do not arrive here.

**Four values are computed on the way out.** When the app's `+0x6bc` is 2, `[0x005cd758]` is
non-null and that object's `+0x84` is not 2:

| `+0x84` | what the display is given |
|---|---|
| 1 | `healthMax := ftol(healthMax × 0.66)`, then `health := healthMax` |
| 3 | `toHit += 50`, `defence += 50`, `healthMax := ftol(healthMax × 1.5)`, then `health := healthMax` |
| 2 | untouched |

`0.66` and `1.5` are the doubles at `0x00599260` / `0x00599268`. `+0x84` is the setting above, and
this table is the display **mirroring** the spawner's rule for a prototype that never went through
it — not a second, display-only rule.

**Not established:** which value appears where. From `+0x14a` on, the block is read by *computed
index* (`FUN_004190f0` `0041953a`), so no displacement sweep names a consumer — the layout needs
the interface layer (`UNIT-PANEL-011`). `FUN_0045f850` has no multi-selection arm.

## What class it is drawn as (`UNIT-APPEAR-030`)

Its own `typeID`, unchanged, and nothing derives it.

The server sends `word[actor+0xe]` — `FUN_00521890`, nine instructions, no arithmetic — under field
mask bit `0x4000`, together with a face byte `actor+0x4b`. The client assigns it to `drawable+0x20`,
which is what the frame selector `FUN_0045bf00` subscripts the `units.reg` class array with. For a
class a map places, `FUN_0045f850` leaves that field alone: its two outer arms are `>= 0x40` and
`< 0x1a`, and the shipped roster occupies exactly `1..27` and `64..80`, so both arms are total on the
shipped population. The appearance routine `FUN_0045fb00` gates on `drawable+0x18c` bit 0, which no
arm sets for these actors, and returns immediately.

So the drawn class of a non-hero actor is a **stored** property and needs no recomputation, and its
sheet is the class record's own `File` in the ordinary way (`REG-UNITS-049`).

This is the opposite of the hero arm, where the id is discarded on arrival and the drawn class is
recomputed from equipment on every state message — `formats/hero/format.md` §6a. **One code path
cannot serve both.** A reimplementation that derives a class for every actor is wrong for every unit
a map places; one that stores a class for every actor is wrong for every player character.

**And it never carries visible equipment.** The equipment sender `FUN_004e873b` calls
`actor->vt+0x30()` before either of its opcode arms and returns without sending anything when it is
false — the same humanoid predicate whose false arm prints *"Trying to takeoff armor from non
humanoid"* (`FUN_005232d0` returns 0 on the base actor class, `FUN_00523530` returns 1 on both human
classes). So a non-humanoid actor's twelve client-side visible-equipment slots stay null for its
whole life, and there is nothing for the hero arm to derive from even if it were reached
(`UNIT-APPEAR-031`).

Claims: `UNIT-APPEAR-030`, `UNIT-APPEAR-031`, `HERO-APPEAR-040`…`045`.

## Open

`FUN_0050d670`'s tier/material parse and the `{castSpell=…}` suffix; the six unread constructor
call sites, i.e. whether any path creates a unit without the fold; `FUN_00532e60`, the only
routine outside the equip pair that writes reach; `FUN_0045f850`, the client-side consumer of
the prototype cache; `Player+0x5c`, the second gate on the experience payout; and whether
anything changes `Player+0x28` after setup.

## The name a unit is shown by (`UNIT-NAME-039`…`UNIT-NAMETAB-041`)

A class's displayed name is **line `ID` of `main\text\unitname.txt`**, not `units.reg`'s `DescText`.
The information display `FUN_00460480` reads it at `004605b1 MOV EDX,[EBP + 0x20] / PUSH EDX /
MOV ECX,0x5eb4c0 / CALL 0x004687f0`, where `[EBP+0x20]` is `drawable+0x20` (the typeID) and
`0x5eb4c0` is the table object. The file has 81 lines against the class array's 81 slots
(`maxID + 1`); 33 of the 34 shipped classes have a non-empty line and every one of those 33 differs
between the roots, while all 42 empty lines are identical. Six non-empty lines (17, 18, 20, 67, 77,
78) belong to no shipped class. The one class with an empty line is ID 2, `Unarmed Fighter with
Shield`.

`DescText` is a second naming system with no reader in `rom.exe`, byte-identical on both roots, and
it disagrees with the shown name in content — `Death Star`/`Daemon`, `Ghost`/`Spirit`,
`Goblin`/`Goblin Pikeman`. A consumer that renders `DescText` renders the wrong string, in the
wrong language.

An **instance** name exists as well — `actor+0x80`, an MFC `CString` on every actor — but no shipped
surface draws it for a non-hero, so two units of one class are not distinguishable by name.

## Multi-cell Building area targets

The simulation object is Building, not the interface CStructure (`UNIT-CLASS-023`). Its presence
mask registers the same pointer in multiple cell records, in y-then-x order. A collision can leave
an attached prefix; the geometry alone is not the accepted-cell set (`UNIT-STRUCTCELL-070`).

Ring and blast visit cell slots, not unique objects. The ordinary cloud pulse omits the Building
slot. For a direct-damage inner effect, each successful lookup can therefore reach a fresh HP
application (`UNIT-AREAVISIT-071`, `UNIT-AREADIRECT-072`). Building's token-size getter returns 1,
so fireball's size-squared divide does not normalize its rectangle. Direct effect application
clamps negative post-subtraction HP to zero and does not reject current HP zero
(`UNIT-AREAHP-073`).

Destructor cleanup walks the mask again, not the accepted prefix, and clears a found `+0xc`
without an identity check. HP-to-destructor timing, live partial-registration teardown and alias
lifetime across real damage calls remain Unknown (`UNIT-STRUCTDETACH-074`). Installed EN/RU
payload and replay populations are explicitly separate (`UNIT-AREAPOP-075`).

## Mission-10 cell-entry caster


`UNIT-M10CELL-054` identifies the two original authored cells, (22,64) and
(21,63), with spell 13 and power 1. They contain source coordinates, not tower
references. Their type-9 source binding is in the ALM format page.

`UNIT-M10ENTRY-055` establishes the admission event: the square-footprint
attachment loop calls the cell helper for each cell. Movement domains 1/2 can
cast from an existing record with spell byte neither zero nor 26; domain 3
does not. The cast request precedes the occupied-ground-slot rejection. The
indexed spell predicate returns literal 1 in this image, so the target is the
entering actor, not the fallback current-cell target. This local arm checks no
allegiance, range, visibility or tower state.

The helper creates a runtime-id-zero actor at the source coordinate, with
Mind=30, skill[Sphere]=power, charge/recovery=1, order 13 and the entering actor
as target. It follows the ordinary spell-13 damage and direct-client effect
routes (`UNIT-M10CAST-056`). The visual flight value 5 is not a repeat timer.

The cell tail is not consumed or time-gated by admission. Another qualifying
attachment attempt may create another caster; standing still is not itself
this event. Zero spell disables this arm and 26 selects another operation.
The persisted cell tail and transient casting actor have different save
lifetimes. No tower-health dependency is present in this local path, but
post-destruction cleanup elsewhere and mission-10 runtime/save-load outcomes
remain Unknown (`UNIT-M10LIFE-057`).
