# Stock and generation

[Reference](format.md)

## Structure

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

The later draw copies each admitted Book once, then draws 1..8 additional
candidates. These effect-bearing Books/Scrolls are non-stackable and each copy
has count 1; copy preserves stored price without recalculation. Finally the
caller adds six explicit Potions (health/mana regeneration, medium and big
health/mana restoration), each with a separate count 51..100. These six additions
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
