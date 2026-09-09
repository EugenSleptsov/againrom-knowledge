# Definitions and creation

[Reference](format.md)

## Definition sources

`units.reg` is the **client drawable's** class record (`REG-UNITS-049`) and the simulation
module never reads that array (`REG-UNITS-061`, amended; this simulation-domain clause is retained). `Data.bin` **Units** is the **simulation
actor's**. The join is `units.reg ID == Data.bin typeID`; the tier `face` (1..4) exists only
on the `Data.bin` side, so one `units.reg` class covers a whole tier family.

Drawable and simulation fields have different consumers:

| looks the same | drawable (`units.reg`) | actor (`Data.bin`) |
|---|---|---|
| footprint | `TileSize` → `CUnit vt+0x20` | `tokenSize` → `actor+0x49` (`TERR-MOVE-054`) |
| attack timing | `AttackDelay` / `ShootDelay` | `attackChargeTime` / `attackRelaxTime` → `+0x134`/`+0x135` |
| death | `Dying` / `DyingPhases` — a corpse **sheet** | `dyingTime` — how long the corpse holds its cells |

## Creation sequence

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
    it has just built — see *Unit scaling setting*. **Non-heroes only**: the humans arm
    of the same routine jumps over the block (`004e2acf`).

## Units slot map

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

## Derived-state boundary

`vt+0x50` on `0x59c3c0` is `FUN_004f5946`, whole:

```
if (actor+0x14) actor+0xa0 = (u16)actor+0x9c * [actor+0x14]+0x58 / 100
else            actor+0xa0 = actor+0x9c
```

`actor+0xa0` is the mana floor the heal AI `FUN_0052eb60` will not cast below. Nothing else
is recomputed, ever. It is invoked whenever `load/capacity` changes bucket
(`FUN_004f36f7` @`004f37b5`), which every equip path reaches.

## Combat inputs

| input | non-hero source | moved by an item? |
|---|---|---|
| damage `(base, spread)` `+0xb4`/`+0xb5` | `physicalMin`, `physicalMax − physicalMin`, via `attackKind` | yes, through the fold |
| absorption `+0xc0` | the `absorbtion` column | yes (armour/shield) |
| `target+0xc6` | the `prot Water` column | yes |
| five protections `+0xc4…+0xcc` | the `prot *` columns | yes |
| damage-kind resistance `+0xce + k` | `res.*` fill `+0xcf…+0xd3`; **`+0xce` (k = 0) is filled by nothing** | no |
| the attacker's `k` = `actor+0xb6` | 0 unless a **melee** weapon sets it to `@.attackType` | yes |
| reach `+0x12c` | **1**, then `+= @.range − 1` per weapon — `@.range = −1` stores 1, so a `Pike` or `Long Sword` stays at 1, and reach is **independent** of the melee/projectile arm | yes |
| attack period `+0x134`/`+0x135` | the two template columns, **assigned over** by `@.charge`/`@.relax` | yes |
| speed `+0x8c` | the `speed` column (ctor 10) — no `reaction` formula on this arm | via the fold's `+0xd8` |
| healthMax `+0x96` | the `healthMax` column (ctor 30) — no XP term on this arm | via the fold's `+0xdc` |

## Installed Units definitions

The installed Units table has 56 parameterized rows with 55 parameters each.
attackKind is−1 or 3; Bat_Sonic and Bee use 3. Equipment strings contain at
most one weapon and no Shield. Installed reach values are 1,4,5,8,20;
Pike and Long Sword retain reach 1 because their range column is−1. These
are installed data values, not arbitrary-template limits. The constructor
and slot map above determine behavior for each authored field.

## Owner field, `actor+0x14`

Set by the ALM spawner from the placement's type-5 group id through the player manager
`[0x00609544]` (`004e274f`…`004e2f16`). It is a **`Player`** (`0x70` B, ctor `FUN_004faa81`).
Three of its fields are read outside the shop:

| field | meaning | who reads it |
|---|---|---|
| `+0x28` | **0 = a human participant owns this unit**; 1 = a scenario-authored owner (also the constructor default); 2 = a group whose `HumanFriend` key is `"Yes"` | 25 sites, incl. the route budget (`MOVE-TERM-003`) and the experience payout (`HERO-KILL-027`) |
| `+0x38` | money (`SHOP-BUY-009`); defaulted to 100 on the join path | the shop |
| `+0x58` | `95`, the percentage `vt+0x50` scales `manaMax` by into `actor+0xa0` | `FUN_004f5946`, `FUN_004f7dfc` |

A unit owned by a human participant gets the **flat
1000-generation** static route budget instead of `max(scalar, D>>2) + D`, and **killing** such a
unit pays the killer **no experience** (`004f7910`…`004f7920`).

## Unit scaling setting

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

## Placement overrides

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

Installed placements leave current mana and the actor+0xbe override absent.
Other tail values can be authored, including actor+0xc0 and skill overrides;
absence in most records does not remove their consumer rules.
— UNIT-PLACEIDLE-088

## Unknowns

`FUN_0050d670`'s tier/material parse and the `{castSpell=…}` suffix; the six unread constructor
call sites, i.e. whether any path creates a unit without the fold; `FUN_00532e60`, the only
routine outside the equip pair that writes reach; `FUN_0045f850`, the client-side consumer of
the prototype cache; `Player+0x5c`, the second gate on the experience payout; and whether
anything changes `Player+0x28` after setup.
