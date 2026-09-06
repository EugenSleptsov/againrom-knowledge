# SHOP — stock, price and trade — specification (partial)

Level 3. Promoted, evidence-backed claims only. Class, generator, pool, enchantment and commit
paths are `SHOP-CLS-001`…`SHOP-NPC-012` (amended); lifecycle, serialization and rounding are
`SHOP-LIFE-013`…`SHOP-ROUND-017`; campaign/town inputs are `SHOP-MISSION-018`…`SHOP-TOWN-023`;
the complete equipment/effect populations and live enchantment generator are
`SHOP-EFFPOOL-061`…`SHOP-EFFSEED-072`. Ledger: `claims/shop.md`.

**Status: partial (◐).** The object graph, the stock lifecycle, the generator's structure and
inputs, the value window, the tier/material/shape triple, the enchantment stage, the price
formula, both commit paths, the refusal, the seed, the save behaviour, **both** providers'
ceilings **and the counter itself** — the table's five places, the move command, the two commit
buttons and the duplicate rules — are read at instruction level. The complete Magic Items effect
population, payloads, pricing, capacity debit, retry paths, and draw order are also specified.
**Not specified:** what advances the session tick counter the restock period counts; whether the
interface can pick up part of a stack; and the exact clock value and intervening draw count at the
first shelf after a mission/save load. **Closed since:** *at which moment of a mission the campaign
ceiling is applied* (see *When the ceiling is applied*), and *whether a tray can hold two items at
once* (see *The counter* — it holds five).

This is not a file format. It is the simulation area sitting on `ALM-CLS-037` (which map
records become shops) and on `formats/databin/format.md` (the item collections it draws from).

## At a glance

### Actor state across trade boundaries

The carried-to-table arm changes containers and Item ownership without calling
the actor load helper (`SAV-CITYMOVE-512`). Sale completion calls that helper
on the requested actor after purse/tray updates. It compares truncating signed
quotients of stored old and refreshed load over stored capacity; only a changed
quotient dispatches Human derive. Zero capacity is not guarded in the helper
(`SAV-CITYSALE-513`). Individual return calls the same helper, while bulk
clear/leave uses a separate return loop without that call in its direct path.
Notification/destructor aliases remain outside the bounded no-call observation
(`SAV-CITYRETURN-514`).

```
placement (ALM type-4, kind 34 or 35)
  |
  |  FUN_004e2462:  new Shop(0x74)                       ; Shop : Building : Token
  |                 cap = (i16)rec+0x06 * 1000  -> Shop::SetCap
  v
Shop                             +0x6c -> CMultiShopTemplate (0x98)
                                 +0x70 =  the value cap (ctor default 10 000 000)
CMultiShopTemplate               +0x04 =  OPEN CUSTOMER count  (the generator's gate)
                                 +0x08 =  pending-restock counter
                                 +0x0c =  4 x CMultiShopShelf(0x1c), codes 1, 2, 4, 3
                                 +0x7c =  CObArray of open CMultiShopInstance (count at +0x84)
                                 +0x90 =  the value cap, copied from Shop+0x70
                                 +0x94 =  the owning object the per-instance tick passes on
CMultiShopShelf                  +0x00 vptr  +0x04 code  +0x08 CArray of items
CMultiShopInstance (0xa0)        one per customer: +0x04 the same 4 shelves,
                                 +0x74 customer, +0x78 the tray list; max 250 open
```

## Where a shop comes from

`Shop` is instantiated at exactly three sites in the image:

| site | provider |
|---|---|
| `FUN_004e2462`, the ALM type-4 walker | one `Shop` per placement whose `kind` is `0x22`/`0x23` |
| `FUN_004d4d42` | a placement-new into the **static** object at `0x0060a120` |
| `FUN_00505acf` | the MFC `CreateObject` thunk — no direct caller; reached through `CArchive` |

The shop-by-id resolver `FUN_004d4db1` returns the static object when `server+0x14c != 0` and
otherwise looks the id up in the map's objects behind an `IsKindOf(&Shop)` guard. **Both
providers yield the same class and the same state block**, so everything below is shared
(`SHOP-ENTRY-003`).

**`server+0x14c` is the single-player condition**, and it decides which of the two you get
(`SHOP-ENTRY-016`):

```
map type-0 +0x70  (participants; 1 on 56/56 campaign maps, 4/8/12/16 on loose ones)
   -> map+0xd4
   -> server+0x0c  = (participants > 1)          ; "multiplayer"
   -> server+0x14c = (server+0x0c == 0)          ; "single player"
   -> FUN_004d4db1: if (server+0x14c) return &static shop, WITHOUT reading the id
```

Consequence a consumer must reproduce: **on a campaign map, a shop placed on the map is
built by the ALM walker and can never be opened.** Exactly one shipped campaign map has such a
placement (`scn:120.alm`, id 35, cap 100 000). The map arm of the resolver is the multiplayer
path and nothing else.

## The one per-shop parameter

A shop takes exactly **one** number from outside — the **upper bound on the value of an item
the generator may stock**, `Shop+0x70` → `CMultiShopTemplate+0x90`. Each provider fills it its
own way, and the two are **not** the same arithmetic:

| provider | source | arithmetic |
|---|---|---|
| map placement (multiplayer) | type-4 word at file `+0x0c` (in-memory `rec+0x06`) | `(i16)rec+0x06 × 1000` (`SHOP-CAP-004`) |
| the static town shop (campaign) | `scenario.reg` `[Mission<n>] ShopMaxPrice` | **verbatim, no multiplier** (`SHOP-MISSION-018`) |

Shipped map range, EN root, 50 placements: 23 distinct words, caps **4 000 … 10 000 000**. The
constructor default (10 000 000) is also the largest shipped value; it is not the typical one.

### The campaign chain

```
scenario.res:scenario.reg  [Mission<n>] ShopMaxPrice     ; absent -> the getter's default 0
   -> FUN_004879f0   00487cc2  MOV [record+0x14],EAX     ; the mission record
   -> the record is embedded at campaignScreen+0x548
   -> 00473ff1 / 00478beb  MOV ECX,[obj+0x55c]           ; its only two readers; no writer
   -> FUN_0041da75   0041da88  cmd+0x09 = 0x3f           ; the staging command at 0x609c38
                     0041da9b  cmd+0x0e = ShopMaxPrice
                     0041da92  cmd+0x0a = ShopMinPrice   ; packed and never read
   -> FUN_004d5dd8 case 0x3f, gated server+0x0c == 0     ; SINGLE PLAYER ONLY
                     004d7189  MOV EDX,[cmd+0x0e]
                     004d7192  Shop::SetCap  on 0x60a120 ; the static town shop
                     004d719c  Shop::Generate
   -> 00505e1d  shop+0x70 = value   ->  00524d8d  template+0x90 = value   ->  the generator's max
```

The shipped ladder is strictly increasing and **two missions ship no key at all**, so their
ceiling is the loader's own default, 0. No sub-mission section carries either shop key. **Those
two missions are also the two during which the town cannot be entered, so ceiling 0 never
reaches a shop the player can open** — see *When the ceiling is applied*.

```
Mission10  Mission20      (no ShopMaxPrice)          ceiling 0
Mission30 …150            1000  3000  7000  10000  30000  70000  100000
                          300000  700000  1000000  3000000  7000000  10000000
ShopMinPrice, same order  0 0 0 0 0 300 400 500 500 600 800 1000 2000 3500 5000   ; never read
```

**`ShopMinPrice` has no consumer.** The window's lower bound is the literal `0` at every site
that sets it, so a consumer must not implement a rising floor: the campaign raises a ceiling,
it does not move a band (`SHOP-MISSION-019`). What the constant 0 does do is exclude
`Data.bin`'s negative-price sentinels — the deleted `rem` weapon row and 25 `"Quest …"` magic
items, all `-1`.

## When the ceiling is applied

Command `0x3f` is sent from **exactly two** places, and neither is a mission's start
(`SHOP-TOWN-022`).

```
(a) end of a mission -- FUN_00473110, message 0x41d, campaignScreen+0x6bc == 2
      FUN_004768c0                 switch to the global-map view
      00473f35 FUN_004884d0        record+0x110 (AutoGetMission) != -1 ?
                                     yes -> load THAT mission's record, travel to its map
                                            object, SEND NOTHING
                                     no  -> fall through
      00473fb5 FUN_00488970        the finished mission was main -> advance record+0x04 by 10
                                   and LOAD the next mission's record; else delete the
                                   completed sub-mission
      00473fba                     travel to global-map object 0 (the city)
      00473ff1 / 00473ff7          read record+0x14 / +0x10 -- of the record as it is NOW
      00474005 FUN_0041da75        build command 0x3f

(b) load game -- FUN_00478af0, one caller (message 0x419)
      FUN_00477560                 the save's [CurrentState] InBattle
        != 0 -> FUN_00477c00(1), resume the battle, SEND NOTHING
        == 0 -> 00478bff FUN_0041da75, then PostMessage 0x42e (enter town)
```

Three consequences a consumer must implement:

1. **The ceiling is the *upcoming* mission's.** Finishing main mission *N* caps the town shop at
   `[Mission(N+10)] ShopMaxPrice`, because `FUN_00488970` has already loaded that record.
   Finishing a *side* mission leaves the record alone and re-sends the same value.
2. **Every homecoming re-rolls the stock**, because the `0x3f` arm does `Shop::SetCap` **and**
   `Shop::Generate` unconditionally in single player (`SHOP-LIFE-013`).
3. **Before the first homecoming there is no stock.** Single-player `Shop::Generate` has no other
   caller, and a new campaign goes hero creation → `0x42f` → mission 10 without a town.

The campaign's own routing, from `scenario.reg` and `rom.exe` (`REG-SCN-063`):

```
new campaign      FUN_00487640 -> FUN_00488460(10); hero creation; PostMessage 0x42f (battle)
Mission10 ends    [Mission10] AutoGetMission = 20  -> straight to mission 20, no town
Mission20 ends    no AutoGetMission                -> town; record advances to Mission30
                                                      SetCap 1000 + Generate  <- the first ever
Mission30..150    13 town visits, ceilings 1000 ... 10 000 000
Mission150        LastMission = 1
```

`AutoGetMission` (default **-1**, `00487b1e`/`00487b30`) ships **once** in the whole campaign, and
`LastMission` (default 0) once. Values identical in the live, EN and RU roots.

## The stock lifecycle

The generator `FUN_005073e2` is the only one in the image, and it **clears before it fills**:
after its gate it calls `FUN_005073a9`, four passes of `FUN_00506670`, each of which destroys
every element of a shelf and then `RemoveAll`s the array. It is never an append.

Its gate is `template+0x04 < 1` — **"no customer currently has the shop open"**, not "the
shelves are empty". `+0x04` is incremented by `CMultiShopTemplate::CreateInstance`
(`FUN_00507162`) and decremented by the instance destructor (`FUN_005067c2`).

```
single player   city-screen command -> Shop::SetCap(cmd+0x0e) ; Shop::Generate   (unconditional)
                => the stock is rebuilt every time that command runs

multiplayer     Shop vt+0x14 -> FUN_00507db0, gated on server+0x0c != 0:
                    if (server+0x00 % 180 == 0) { template+0x08 += 1; TryRestock(); }
                TryRestock (FUN_00507e56):
                    if (template+0x08 > 0 && template+0x04 == 0) { Generate(); template+0x08 = 0; }
                and the customer-leaves path FUN_00507a1c calls TryRestock as its last act
                => a request raised while someone is inside is held, and fires when the shop empties

either way      Shop::Open -> if the shelves are empty, Generate() first, then create the instance
```

Reached from `Shop::Open` (`FUN_005079a7`), `Shop::Generate` (`FUN_00505e54`) and
`FUN_00507e56` — three callers, 0 in orphan code (`SHOP-GEN-005` (partially retracted), `SHOP-LIFE-013`,
`SHOP-LIFE-014`).

### Persistence

**Nothing of the stock is saved.** `Shop::Serialize` (`FUN_00505fcc`, MFC vtable slot `+0x08`,
whose only reference is that slot) writes `Building::Serialize` plus one `int` — `shop+0x70`,
the value cap — and never touches `shop+0x6c`. `CMultiShopTemplate` does not override
`Serialize` at all: its slot `+0x08` is the `CObject` no-op stub `FUN_00401950`. A loaded game
therefore starts with an empty shop (`SHOP-SAVE-015`).

| shelf | code | drawn | source | second stage |
|---|---|---|---|---|
| 0 | 1 | 100 | `Data.bin` Shields (`0x609b40`) + Armors (`0x609b54`) | — |
| 1 | 2 | 100 | `Data.bin` Weapons (`0x609b68`) | — |
| 2 | 4 | 20 | the union of the two above | **enchantment required** |
| 3 | 3 | `rand(1..8)` | `Data.bin` **Spells** (`0x609bf4`), via `FUN_00509d3f` — two subtypes `0x2a`/`0x29` priced from param slots `0x15`/`0x14` | — |

Then six literal potions are appended by name — `Potion Health Regeneration`,
`Potion Medium Healing`, `Potion Big Healing`, `Potion Mana Regeneration`,
`Potion Medium Mana`, `Potion Big Mana` — each at quantity `rand(1..50) + 50`.

### The candidate pool

For a collection and a tier, `FUN_005095ee` walks the entries (index from 1) and, per entry,
the 16 bits of the `u16` mask at `entry+0x1c + 2·tier`. Each set bit is a candidate
`(tier, materialBit, shapeIndex)`, admitted iff

```
v = ftol( param[2] * Materials[m].+0x30 * Shapes[t].+0x30 )      ; 005096af..005096bd
min <= v <= max          ; min = 0 at every site that sets it, max = template+0x90
```

`param[2]` is the entry's `Price` column; `Materials` is `0x609b18` (16 entries, indexed by the
mask bit) and `Shapes` is `0x609b2c` (5 entries, indexed by the tier); `+0x30` is double #2 of
the 0x48-byte block those entries carry at `element+0x20` (`FUN_004df328`). **So the window
bounds the item's own price** — the same product `item+0x1c` is built from, minus its `+0.5`
(`SHOP-POOL-021`; `SHOP-POOL-006`'s base-only reading is retracted). The five masks × 16 bits
are `DAT-OBJ-002`'s "+10 raw bytes" on an armor/shield/weapon entry. The item stores the triple
at `+0x45` (tier), `+0x46` (material), `+0x0c` (shape).

Every generator call asks for tier `5`, which `FUN_005098e2` reads as **all five**, `0..4`.

Measured over the EN `Data.bin` (`tools/placedb -mode shopwindow`): **367** mask-admitted
triples, values `[0 … 768 000]`; six entries are masked out at every tier (`BareHands`, `rem`,
`Sonic Beam`, `Flame Thrower`, `Boulder Thrower`, `Plasma Sword`). Admitted per campaign
ceiling: `0 → 4`, `1000 → 142`, `3000 → 172`, `7000 → 214`, `10 000 → 236`, `30 000 → 286`,
`70 000 → 318`, `100 000 → 330`, `300 000 → 354`, `700 000 → 365`, **`≥ 1 000 000 → 367`** —
i.e. the ceiling stops discriminating from `Mission120` on (`SHOP-MISSION-020`). **The `0 -> 4`
row is arithmetic about a state no player reaches** (`SHOP-TOWN-023`): the campaign never opens
the town while the ceiling is 0, and until the first homecoming the static town shop has been
neither capped nor generated, so it is empty rather than four-item.

`FUN_00509ac8`'s kinds **3, 4 and 5** — the Magic-Items pools filtered by the names `Potion`,
`Scroll`, `Book` — are reachable only from `FUN_005094fd`, which has **0 callers**. The four
wrappers the generator uses pass kinds 1 and 2 only.

### The draw

`FUN_0050b055(n, k, dest)` picks `n` candidates uniformly with replacement (retry budget
`10 n`), clones each, and sets a stack size: `1` if the item's `vt+0x50` is false, else
`rand(1..k)`. For the armour and weapon shelves `FUN_0050abb3` then overrides it **from the
tier**: `rand(1..8)`, `rand(1..4)`, `rand(1..2)`, `1`, `1` for tiers 0..4.

The **consumables** pool is different (`SHOP-CONSUME-073`, `SHOP-CONSUME-074`). `FUN_00509d3f`
constructs named Book templates with a kind-42 Effect and named Scroll templates
with a kind-41 Effect. The spell id is Effect+0x40; Item+0x40 remains its packed
item identity. Books use Spells parameter 21; Scrolls use parameter 20 as base
price. Both test that unscaled price against the inclusive window.

For Scrolls, `n = min(100, 10·floor(max/base))` and `power = rand(1..n)` go to
**Effect+0x42**, not the Item count. Stored Item price becomes
`base*(floor(power/10)+1)`, with no second ceiling test. Both loops consider
ids1..27; only ids2..5,7..10,13..16,19..26 construct Books. Slow28 is excluded.

The later draw copies each admitted Book once, then draws1..8 additional
candidates. These effect-bearing Books/Scrolls are non-stackable and each copy
has count1; copy preserves stored price without recalculation. Finally the
caller adds six explicit Potions (health/mana regeneration, medium and big
health/mana restoration), each with a separate count51..100. These six additions
have no price-window test. The Scroll-power law and Potion-count law are not
interchangeable.

### The enchantment stage

Only the Magic Items shelf demands this second stage. Its complete base pool is 367 mask-admitted
triples on each preserved root: 193 Armor, 36 Shield, and 138 Weapon. The two listings are
row-identical. Tier and material choose the base triple, base price, and MagCap; they do not choose
an effect tier.

The 50-row `Magic` collection holds 24 cumulative endpoint columns: fighter slots 1..12 followed
by mage slots 1..12. For fighter flag `C` and equipment slot `s`, the endpoint parameter is
`s+3+(1-C)*12`. The selector draws inclusive `1..finalEndpoint`; kind `k` wins when
`endpoint[k-1] < draw <= endpoint[k]`. Positive endpoint differences are weights. The complete
table has 331 positive class/slot rows; the shipped shelf pool reaches 273 of them and 36 distinct
kinds. Weapon is slot 1, Shield slot 2, and Armor uses the definition's own slot. `UsableBy` does
not filter this selector (`SHOP-EFFPOOL-061`, `SHOP-EFFWEIGHT-062`).

For a non-cast kind with cost `c`, declared range `a..z`, budget `B`, and remaining MagCap `M`:

```text
capMax     = trunc(M/c)
budgetMax  = trunc((70/c) * log base 1.5 (B/(50*M) - 1))
payloadMax = min(capMax, budgetMax, z)
```

Non-positive `B` or `M`, or `payloadMax<a`, returns null without selecting another kind. Otherwise
the payload is `max(rand(1..payloadMax),a)`, worth `payload*c` points. Damage kinds 44..48 first
consume and discard that common draw, cap the maximum at 255, then draw byte base `1..max` and
byte spread `1..floor(max/2)`; their point value is `(base+spread)*c` (`SHOP-EFFPAY-063`,
`SHOP-EFFRANGE-064`).

`castSpell` draws Stone Curse or Drain Life for a fighter, and Fire Arrow, Lightning, Prismatic
Spray, Stone Curse, or Drain Life for a mage. With spell scalar `S`, it requires
`r=log base 2(B/(10*S)) > 0`; the power maximum is
`min(100,trunc(30*(1.2^r-1)))`. Fighter cast uses `currentPrice*10` as `B`. All ten mage-only
weapon triples take the forced-cast arm with
`B=min(2*ceiling-currentPrice,currentPrice*100)` (`SHOP-EFFCAST-065`).

For summed non-cast points `N`, the price addition is
`trunc((pow(N/70,1.5)+1)*N*50)`. Each retained cast adds
`trunc(10*S*pow(log base 1.2(power/30+1),2))`. After every append the stored price is clamped to
9,999,999, and the next call recalculates `2*ceiling-currentStoredPrice`. The final ordinary path
debits MagCap by `N` (`SHOP-EFFPRICE-066`, `SHOP-EFFCAP-069`).

One effect is required. The second gate is `rand(0..100)<50`, exactly `50/101`; the third is an
independent `rand(0..100)<25`, exactly `25/101`, and is drawn even when the second gate is false
unless an earlier return has ended the routine. A required null rejects the item. An optional null
retains earlier effects and returns success; an optional cast is randomized then freed and also
returns success. Those early returns skip the final MagCap debit. Selector attempt 100 cannot
survive, so only attempts 1..99 succeed. The outer item loop tests `attempt<=10*n`, allowing 201
rejections for the Magic Items request of 20 (`SHOP-EFFORDER-067`, `SHOP-EFFRETRY-068`).

### Randomness

`FUN_00554a60` is the CRT `rand()`: `seed = seed·0x343fd + 0x269ec3`, result `(seed>>16)&0x7fff`.
`FUN_00504003(n)` = uniform on `0..n`; `FUN_0050402c(n)` = uniform on `1..n`.

**Every seed in the image is a clock.** The only writer of the seed slot is `srand`
(`FUN_00554a50`), with three call sites: `FUN_00453ec0` and `FUN_00503fd6` seed with
`timeGetTime()`, `FUN_0052b2b0` (the server constructor) with `time(0)`. Nothing derives a seed
from a mission id, a map name, or saved stock. The shop generator does not reseed. The mission/save
map loader reaches `FUN_00503fd6`; ordinary homecoming which does not cross that loader continues
the existing stream. Exact clock value and intervening draws at the first shelf remain Unknown
(`SHOP-RNG-008`, `SHOP-EFFSEED-072`).

## Trading

A customer opens the shop through `Shop::Open`, which creates a `CMultiShopInstance` (max 250
live). Items move between the shelves and the instance's tray at `+0x78`; the tray element's
`+0x14` is its current owner — `0` = the shop, otherwise the customer's `Player`.

```
buy   (FUN_005069d0): for each tray element with owner == 0:
                        total = (u16)elem+0x42 * elem+0x1c
                        if (Player+0x38 < total) STOP THE WHOLE COMMIT
                        Player+0x38 -= total;  owner = Player;  -> customer+0x7c

sell  (FUN_00506893): for each tray element with owner != 0 and elem+0x1c != 0:
                        Player+0x38 += ftol(0.5 * ((u16)elem+0x42 * elem+0x1c) + 0.5)
                        owner = 0;  -> back onto a shelf
```

So **the shop pays exactly half of what it charges**, rounded half-up, and an unaffordable
item aborts the rest of the purchase rather than being skipped — **but see *The counter* below:
the shipped client never sends a basket that can reach the abort** (`SHOP-BUY-009`,
`SHOP-TRAY-026`).

The `ceil` is applied to `quantity × price` as a whole, not per unit. Selling `n` odd-priced
items as one stack pays `ceil(n·p/2)`; selling them one at a time pays `n·(p+1)/2`, up to
`⌊n/2⌋` more. **A consumer that halves per unit will disagree with the engine by up to one coin
per pair, on odd prices only** (`SHOP-SELL-010`).

The unit price `item+0x1c` is a property of the item. Neither commit path reads anything of
the customer except the money at `Player+0x38`, and the generator never sees a customer — so
two customers at one shop pay the same (`SHOP-PRICE-011`).

### Where the price comes from, and every rounding

`__ftol` (`0x0055458c`) rounds **toward zero**. Four rounding sites exist and no more
(`SHOP-ROUND-017`):

```
Armor  FUN_0050c740 |
Shield FUN_0050d0ce  >  item+0x1c = ftol(base * material.+0x30 * tier.+0x30 + 0.5)   ; half up
Weapon FUN_0050dc4e |    base     = param slot 2 of the item's own Data.bin entry
                         material = Materials[item+0x46] (0x609b18), tier = [item+0x45] (0x609b2c)

enchantment addend  FUN_00502a10(n) = ftol((pow(n/70.0, 1.5) + 1.0) * n * 50.0)      ; TRUNCATES
                    item+0x1c += that, then clamp to 9 999 999

sell payout         FUN_00506893 = ftol(0.5 * (qty * price) + 0.5)                    ; half up
buy                 FUN_005069d0 — integer IMUL only, no floating point at all
```

### Which shelf an item returns to

`FUN_00506aea`, given the item's class code at `item+0x44`:

| condition | shelf | code |
|---|---|---|
| `item+0x44 ∈ {3, 4, 5}` | 3 | 3 |
| `item->vt+0x50()` is false | 2 | 4 |
| `item+0x44 == 2` | 1 | 2 |
| `item+0x44 == 1` | 0 | 1 |
| otherwise | 2 | 4 — with the log line `"Item of strange type is returned to shop - placed to Magic Items"` |

### The counter — the table, the two buttons, the move

The shop screen is one view (vtable `0x0059ac88`) with five children, each of a class built at
exactly one site image-wide. Three are item grids; each computes `cols`/`rows` from its own
rectangle at **80 px per cell**, allocates `cols·rows` CRects, and shares one hit test that walks
exactly those rectangles (`SHOP-TRAY-024`):

```
child   rect (l,t,r,b)   cols x rows   container code (vt+0xa8)   what it is
+0x68   0,0,164,303        2 x 3 = 6   parent+0x132 + 5 (5..8)    the shop's four shelves, scrolling
+0x70   0,303,480,390      5 x 1 = 5   4                          THE TABLE  (graphics\interface\ShopTable.bmp)
+0x6c   0,390,480,480      5 x 1 = 5   2                          the party backpack, scrolling
```

The other two children and the whole layout are `SHOP-SCREEN-030`…`SHOP-SCREEN-039`,
whose claim rows carry every number with its citing address:

```
child   rect (l,t,r,b)      id     background art          measured
+0x68     0,  0,164,303   1002   ShopInv.bmp             164x303   cells (1+80c, 31+80r) 80x80, 2x3
+0x74   164,  0,480,303   1005   ShopFrame.256           316x303   the merchant's room
+0x70     0,303,480,390   1003   ShopTable.bmp           472x87    cells (32+80c, 303)   80x80, 5x1
+0x6c     0,390,480,480   1001   (none)                            cells (32+80c, 395)   80x80, 5x1
+0x78   464,  0,640,238   1006   ShopMenu.bmp            176x238   the four command buttons
+0x7c   480,238,640,480      7   (composed at run time)            THE CHARACTER PANEL, borrowed
```

**Every rect above is view-relative, and the view's own rect is
`((screenW-640)/2, (screenH-480)/2, screenW-(screenW-640)/2, screenH-(screenH-480)/2)` — four
globals with one writer each in `FUN_00471790`. The shop screen is a fixed 640x480 panel centred
on the display, so at the shipped 640x480 default the origin is `(0,0)` and every rect above is
also a screen coordinate** (`SHOP-VIEW-044`). Activation clears the whole screen to colour 0 first.

`ShopTable.bmp` is 472 px wide and the table's blit asks for 480 (`SHOP-SCREEN-032`).

The shelf grid carries two 72x32 arrow rects, `(46,0,118,32)` and `(46,271,118,303)`, each
scrolling by one row of two cells (`SHOP-SCREEN-033`). The four shelves are chosen by clicking the
room picture. **The merchant panel holds two four-rect arrays and they are different**: the hit
loop reads `panel+0x60 + 0x10*i` and the draw loop reads `panel+0xa0 + 0x10*i`, and the art loader
names the folder as `4 - i` (`SHOP-SHELF-047`, correcting `SHOP-SCREEN-034`):

```
i   hit rect                draw rect               size     folder
0   354,110,459,295         353,108,433,220         80x112   shopanim\04
1   169,110,274,295         197,108,277,220         80x112   shopanim\03
2   314,  5,454,105         313, 20,445,108        132x88    shopanim\02
3   172,  5,314,105         201, 20,313,108        112x88    shopanim\01
```

Each draw rect is exactly the size of its own eleven loaded frame files. Opening a shelf resets the rack
to its first item and stores the shelf index in `view+0x132`, which is `0x64` while none is open
(`SHOP-SCREEN-034`).

The merchant himself is drawn at `(277,112)`, 76x176, from
`movies\shopanim\Pose2-3\1.bmp` — two immediate offsets
`(+113, +112)` from the merchant panel's own top-left, not a rect. His Yes and No animations
replace him at the same anchor (`SHOP-MERCHANT-046`).

#### Shop-interior progression

Field offsets and mode masks in this section are hexadecimal; counts and
milliseconds are decimal.

These are static code results on the identical preserved EN/RU executable, not an
observed original-runtime recording. Merchant-panel paint `004ab37e` is vtable
`0059ad10+2c`. When `view+148` is nonzero and unsigned
`timeGetTime()-005f1b14 >= 100`, it calls the sound hook, stores a new clock
sample, advances enabled racks, then advances at most one merchant mode. Drawing
follows that block even if the block was skipped. No elapsed-time catch-up loop
is present. Thus the local gate is100 ms, not guaranteed visible10 fps. The
initial clock value is `now-100`; the two paint timestamps are static globals,
not one clock per series (`SHOP-ANIMATION-081`).

| series | start and normal progression | termination or replacement |
|---|---|---|
| Rack `i=0..3`, folder `4-i` | Selection loads files1..11, sets bit`1<<i`, resets index0. Eligible paint increments; reaching9 rewrites the index to3. The steady loop uses indices3..8, files4..9. | Switching sets the old rack index9. A paint before the next eligible advance may draw file10. Advance9→10 disables/frees it before drawing; file11 is not a normal draw. |
| Base merchant | `Pose2-3/1.bmp`, `panel+128`; drawn when no merchant mode bit is set. | Replaced by the highest-priority enabled mode below. |
| Idle, bit10 | One bitmap at+12c is replaced per step; index+254 increments modulo30 and loads `Pose2-3/(index+1).bmp`. From0, files2..28 are drawn. | At index28, file29 is loaded, the bit clears and the index resets before the draw. |
| Yes, bit20 | Array+1e0 has base pose in slot0 and Yes files2..12 in slots1..11. | Advance to12 clears the bit, resets index0 and frees the array before the draw. |
| No, bit40 | Array+210 has base pose in slot0 and No files2..12 in slots1..11. | Same twelve-step boundary and cleanup as Yes. |

The first rack/base reaction frame can be skipped if the first paint already
advances it. Selecting the current shelf does nothing. Entry selects rack0,
folder04; clicking a different rack also arms Yes. The four rack folders each
contain11 files; Pose2-3 contains29, Yes13 and No12, with equal corresponding
file hashes across roots. Yes/1, Yes/13 and No/1 are not loaded by these reaction
loaders. Counts describe inventory, not frames guaranteed visible
(`SHOP-ANIMATION-082`, `SHOP-ANIMATION-083`).

Idle is armed only when none of bits10/20/40 is set and elapsed time since
`005f1b10` reaches a freshly computed `5000+1000*(rand()%5)` threshold. The
threshold is recomputed at every eligible pass, including passes whose mode
bits prevent arming. It is not a single uniformly sampled5–9 second delay.
Completion of each merchant mode refreshes the idle timestamp
(`SHOP-ANIMATION-081`).

Both advance and draw prioritize idle, then Yes, then No. Sell-with-items and
affordable buy load Yes and OR20; unaffordable buy loads No and OR40. These
triggers and the rack-click Yes arm do not reset the shared index or clear other
mode bits. Repeated input therefore is not a proved fresh restart. The two
parameterized constructors zero the index. Exit clears flags and frees static
and animated resources; entry reloads resources and resets rack selection.
Neither named transition resets the merchant index or static timestamps.
Object reuse, inherited transition effects and actual mid-animation re-entry
remain Unknown (`SHOP-ANIMATION-084`).

The sound hook runs before frame advancement. While idle it requests view+b4
at index1, view+94 at10/14/18/22, and view+b8 at24. With no merchant mode it can
request view+b0 after its separate redrawn30-second-base threshold,
`004ae650(view+70+20b0)==0`, and non-null/guard tests. The table-query helper's
meaning is unexpanded. Shelf selection requests+a0; sell+a8;
affordable buy+a4; refusal+90; entry+ac and conditionally+b0. The common request
helper requires a non-null sample and guard result0. Its priority128 argument
is not a sound id. Sample identity, audible overlap, presentation timing and
the unexpanded entry+bc helper remain Unknown
(`SHOP-ANIMATION-085`, `VIDEO-SFX-013`).

The button panel's four rects, top to bottom, are `(494,15,614,67)` clear the table with the purse
printed on it, `(483,67,623,113)` **buy** with the buy total, `(483,114,623,160)` **sell** with the
sell total, and `(494,160,614,212)` clear and leave with the sum (`SHOP-SCREEN-035`).

One cell draws a background chosen from five (`backinvs.bmp` for a shelf item that is both
affordable and usable, `backinv.bmp` for a usable one elsewhere, `backinvg.bmp` otherwise), then the
item icon, then the quantity as `"%d"` at `(cell.left+10, cell.bottom-15)` when it exceeds 1, then a
price plaque right-aligned to the cell's right edge. The plaque is `costm(d+1).bmp` on the player's
side and `costs(d+1).bmp` on the shop's, `d = floor(log10(price))` clamped to 6, and the number on
it is `ceil(price/2)` on the player's side and `price` on the shop's (`SHOP-SCREEN-036`,
`SHOP-SCREEN-037`). The grids draw no item name and no characteristics panel
(`SHOP-SCREEN-039`).

Of the 134 resources the screen names, 132 are byte-identical on the EN and RU installs; the two
that differ are `main\text\tips\shop1.txt` and `shop2.txt` (`SHOP-SCREEN-038`).
Both are shown in **one** widget of a different class (vtable `0x0059ba18`, ctor
`FUN_004c7422(id, l, t, r, b, text)`), id 1011, at `(0,162,312,298)` inside the **merchant panel** —
view-relative `(164,162,476,298)`. It exists only when `[0x005eb52c]` is non-zero, and `shop2.txt`
replaces `shop1.txt` at most once per visit, behind the latch `view+0x8c` (`SHOP-TIP-045`).

The money cell — an element whose `+0x6` is `0xffff` — draws `backinv.bmp` and then
`graphics\interface\money\money.16a`, a single 80x80 `.16a` frame that
composites over it (`SHOP-MONEY-048`).

### The character panel, `+0x7c`

The sixth region is **not built by the shop**. It is `campaign+0xe0`, the third of four widgets the
main-frame builder stacks in the mission screen's 160-pixel right strip, at `(0,238,160,480)`. On
activation the shop takes it from `campaign+0xd4`, offsets its rect by the literal `640 - width`,
re-sets it and adopts it; on deactivation it offsets back and hands it over. One object, two parents
(`SHOP-FIGURE-041`).

It draws the shown character **composed in that character's own equipment**: virtual slot `+0x80` of
the member's drawable — the world figure compositor, which has no layer selector, so the head is
included — into a 160x240 canvas blitted whole at `(panel.left, panel.top + 2)`. The composition is
cached on a string key at `panel+0x7c` and redone when the key changes or when `member+0x18c` bit 3
is set. When `member+0x18c & 0x11` is zero the same routine draws a flat picture instead
(`SHOP-FIGURE-042`).

The **party picker** is the panel's own, not the shop's. Two 32x32 rects, panel-relative
`(1,205,33,237)` and `(119,205,151,237)` — `(481,443,513,475)` and `(599,443,631,475)` at 640x480 —
post messages `0x414` and `0x415`. The shop view's arms step `view+0x130`, a `u16` index into the
roster `CArray` at `view+0x108`, with wrap-around in both directions, and then in the same call:

```
view+0x6c . vt+0x90(member+0xc8)        bind the item strip's list
[view+0x6c] + 0x90 = member+0xdc        a POINTER to that member's own scroll base
FUN_004a2cd0 ; FUN_004a2fa0             rebuild and repaint the strip
member+0x18c |= 8                       recompose the figure
view+0x14c   |= 9                       invalidate
```

So one step changes the figure and the strip together, and **the strip at `(0,390,480,480)` is the
shown member's own container rather than a party-wide backpack** (`SHOP-PICKER-043`).

**The table is one strip of five places, shared by both sides of the deal.** It does not scroll;
its display-list add refuses to append a sixth element; and the server-side container
(`instance+0x78`) has no capacity of its own. Which side an element is on is ownership, not
position — `elem+0x14` on the server, the container code on the screen — and it is drawn with
`myitem.256` + `costm1..7` (yours) versus `shopitem.256` + `costs1..7` (the shop's), seven price
plaques per side, one per digit of the `9 999 999` clamp.

Every move is command **`0x22`** — `FUN_0041c98f(srcCode, srcIdx, dstCode, dstPos, qty)`, packed at
`cmd+0x0c/0x0e/0x0d/0x10/0x12`. Container codes are shared with the inventory screen: `1`
equipment, `2` the party backpack, `4` the tray, `5..8` the four shelves; the shop arm requires both
codes in `4..8`. The shelf helper `FUN_00506552` retains its published split rule:
`stack -= (qty−1)`, detach one unit through the item's `vt+0x40`, give the detached object
`qty`. The list helper `FUN_0050ebf9` splits only when count exceeds requested quantity.
Otherwise `0050ec52` branches to whole removal at `0050ec86`, bypassing the split virtual
and returning the same Item pointer; count one/request one follows that whole-item branch.
This correction does not remeasure the shelf helper or city-sale lifecycle. The destination-4
arm stamps `item+0x14 = Player` (`SHOP-TRAY-025`).

The screen keeps three running numbers over the table (`FUN_004a95ec`) and gates three commands on
them (`SHOP-TRAY-026`):

```
purse      = campaignScreen+0x9b4 -> +0x0c
buyTotal   = -SUM over elements whose code != 2 of  price * qty        (full price)
sellTotal  = +SUM over elements whose code == 2 of  ceil(price/2) * qty  (half, PER UNIT)

0x33 buy    only if buyTotal  != 0 and purse + buyTotal >= 0   (else the refusal animation)
0x34 sell   only if sellTotal != 0
0x35 clear the table
```

**Two consequences a consumer must reproduce.** (1) Because the client refuses an unaffordable
*total* and the server debits monotonically, the buy loop's abort is unreachable through the
shipped UI. (2) The screen halves **per unit** and the shop pays **per stack**, so on an odd price
the displayed sell total exceeds the coins received by `⌊qty/2⌋` (`SHOP-TRAY-027`) — the same delta
as splitting the stack. A consumer that makes them agree has removed a visible engine behaviour.

### Duplicates

Two items are "the same thing" iff `FUN_00508671` says so: equal item code, then **both** stackable →
yes without reading effects; one stackable → no; **neither** stackable → only if their ordered
effect lists at `+0x20` agree. Differing enchantment therefore separates equal-code non-stackable
items. Equal-code Potions are stackable even with effects and compare equal without examining them.

The 100/100/20/`rand(1..8)` random draws use plain `CObArray::Add`, so equal draws remain separate.
The six literal Potions and items *returned* to a shelf instead use the merging insert
`FUN_005062e8`; the family's other merge routine, `FUN_00506434`, has 0 callers. A fresh shelf can
therefore hold duplicate random draws as separate objects — 200 draws from at most 367 admitted
triples — while literal-Potion additions and returned equal items merge (`SHOP-DUP-028`).

## The runtime class table

`.data` carries 28 MFC `CRuntimeClass` records (`{name, objectSize, schema, CreateObject, base,
next}`, 24 bytes). The shop-relevant hierarchy, with the image's own names:

```
CObject
├── Player 0x70                     ; +0x38 = money
├── CMultiShopShelf 0x1c   CMultiShopInstance 0xa0   CMultiShopTemplate 0x98
└── Token 0x3c
    ├── Building 0x6c  ├── Outpost 0xac  ├── Tavern 0xa0  └── Shop 0x74
    ├── Item 0x50      ├── Armor 0x68    ├── Shield 0x68  └── Weapon 0x84
    └── Sack 0x44
```

The full 28-row runtime-class population, including the unit/spell/effect side, is bounded by
`SHOP-CLS-001`.

## Consumer notes

- Do not read the shop record's `+0x0c` as durability, and do not assume it is constant.
- Do not model stock as a fixed list. **Do not model generation as once-only either**: it is a
  clear-and-refill, gated only on nobody having the shop open, and it is driven every
  city-screen entry in single player and every 180 ticks in multiplayer.
- **Do not persist a shop's stock across a save.** The engine does not.
- **Do not seed the generator from anything stable.** The engine seeds from a clock, so
  re-entering a town re-rolls the shop; a mission-seeded generator is a different game.
- On a **campaign** map, a placed shop is unreachable — route every shop lookup to the single
  town shop when the session has one participant.
- Buy and sell prices differ by exactly a factor of two; the sell `ceil` applies to the whole
  stack, not per unit — **and the screen's own total halves per unit, so it over-promises by
  `⌊qty/2⌋` on odd prices. Reproduce both, not one.**
- **The table holds five places, shared by both sides.** Not five per side, and not one at a time.
- A failed affordability check stops the whole commit — but the client refuses the basket first, so
  a faithful client never reaches it.
- **Do not dedup random-draw stock.** Use the merging insert for the six literal-Potion additions
  and returned items. Differing effects separate equal-code non-stackable items; equal-code Potions
  ignore those effect differences.
## Document access item exclusion

Shop generation cannot create the document access item at campaign start. No town or stock
generation occurs before mission 10, and mission 10 routes directly to mission 20. Later live stock
generation uses weapons, armor, shields, spells, and six literal potions. The MagicItems name-filter
routine has no live caller. `Quest Documents` also has price -1, below the live pool floor 0
(`SHOP-DOC-029`).
