<a id="tavern--mercenary-hire--specification-partial"></a>

# Tavern mercenary selection and hire

Campaign registry lists, persistent unlocks and per-type pool counts determine
the mercenary shelf. Hiring selects complete pools by type; mission end
updates losses/recovery and clears hire flags. Actor definitions come from
[Data.bin](../databin/format.md). MERC-SHELF-002's early-display inference
is retracted: being named in a roster does not display a locked type.
— MERC-TYPE-001, MERC-SHELF-002,
MERC-HIRE-003, MERC-PRICE-004, MERC-LEVEL-005, MERC-DEATH-006

The multiplayer tavern path and the mode selecting the flat price 450 remain
Unknown. Room drawing is specified by [TOWN](../town/format.md).

## Entry and saved Player prerequisites

The located campaign helper sends command 37 before computed tavern activation.
Its reached server arm resolves the command key against signed Player+04 and
exits without stock construction on a miss. A match reaches the stock handler's
refund: wrapped 32-bit addition to Player+38, followed by event67 even when the
refund is zero, before the old roster is inspected. This is a local static
order. Earlier UI/transport effects, actual client key and absolute first
runtime consumer remain Unknown. (SAV-918)

Saved Player slot and colour values alone do not identify a failing tavern
operation or establish the native construction prerequisites. — SAV-920

<a id="at-a-glance"></a>

## Structure

```
type t            1..15, the [npc<t>] section number and the MercenaryCount subscript
                  t = 1,2  -> Unit  "Catapult", "Ballista"
                  t = 3..15 -> Human "NPC%02d_%d" % (t, level(mission))
pool[t]           runtime, starts at [General] MercenaryCount[t-1]
hired[t]          0/1; persists in a save between a hire and mission end (SAV-614, SAV-615),
                  cleared at the end of every mission
enabled           a set of t, grows only, from [Mission<n>] EnableMercenary on completion
shelf(mission)    { t in [Mission<mission>] Mercenaries : t in enabled and pool[t] > 0 }
cost(t, mission)  (PriceA[t] + pool[t] * PriceB[t]) * unitPrice(mission)
level(mission)    m10..m50 -> 1   m60..m90 -> 2   m100..m120 -> 3   m130..m150 -> 4
unitPrice(m)      m30..m150 -> 10 15 20 40 60 80 100 600 800 1000 6000 8000 10000
                  anything else -> 0
```

<a id="the-type-space"></a>

## Type identity

There are fifteen mercenary types and they are one index space shared by four files:
`[General] MercenaryCount`'s fifteen elements, `npc.reg`'s `[npc<t>]` sections, the runtime
object's own byte (`CUnit +0x15b` on the client, the server object's `+0x14c`), and the
`Data.bin` template name. Two of the fifteen ship empty — `MercenaryCount` is 0, no `PriceA` /
`PriceB`, no mission offers them, none unlocks them.

A consumer implementing this needs no separate "mercenary table": the type id **is** the npc
id, and `MercenaryCount[t-1]` **is** the headcount.

<a id="what-the-tavern-shows"></a>

## Shelf filters

Three independent gates, in this order:

1. `[Mission<n>] Mercenaries` — the mission's shelf, a list of type ids. Reloaded whenever the
   main mission record loads, so the town shows the *upcoming* mission's shelf (the record has
   already advanced when the party comes home — `SHOP-TOWN-022`).
2. the **enabled** set — a permanent accumulation of every `EnableMercenary` element of every
   mission the player has *completed*. Never cleared short of a new campaign.
3. `pool[t] > 0`.

**A `Mercenaries` key may *name* a type that is not enabled yet, and then the tavern does not
show it at all** — the filter runs before the list is built. The shipped campaign does exactly
that for four types, whose unlocks are the reward for side missions 41, 71, 111 and 121:

```
type   named by Mercenaries from   unlocked by   first main-mission shelf that shows it
 10    mission 40                  mission 41    mission 50
  8    mission 70                  mission 71    mission 80
  1    missions 10, 20, 110        mission 111   mission 120
  5    mission 120                 mission 121   mission 130
```

A consumer must therefore keep the two apart: *named by this mission* and *hireable in this
mission's town* are different sets, and only the second is what the player sees. The right-hand
column above is computed by walking the main-mission ladder in ascending order; whether a side
mission's unlock can take effect inside its own chapter, before the next main mission, depends
on when the player does it and is **not** established here.

Mission 20 is a useful boundary case. Its record names type 1, but type 1 is not enabled until
side mission 111, so the mission-20 shelf is empty. The next main mission, 30, names type 14;
type 14 was enabled by mission 10 and has a positive pool, so it appears. Its level-1 constructor
uses `Humans[58] NPC14_1`, the same template and equipment as three map actors transferred and then
culled in mission 20. Those map actors do enter the client bit-4 bucket read by the live tally, but
their server and client mercenary-type bytes remain zero and the tally matches only types 1..15, so
they increment no pool slot. The agreement is template reuse, not conversion or identity
persistence (`PARTY-M20-032`).

A Human template key does not establish that an actor was hired. Mission 20
can transfer template 58 (`NPC14_1`) while its tavern has no eligible type-14
shelf. Side mission 41 directly places template 54 (`NPC10_1`) under the
Player's group; this does not enable type 10 or make those actors hires.
Use mission identity, the actual mercenary-type byte and the unlock/pool
relations. — PARTY-M20-032, SAV-623, SAV-624, SAV-625, SAV-626, SAV-627

## Hiring

The hire button toggles `hired[t]`, one flag per **type**. Leaving the inn sends the vector
`pool[t] * hired[t]` for `t = 1..15` together with the mission number, and the simulation
spawns that many of each type. **There is no way to hire one man of a type with four.**

Money is checked locally against `player.money − Σ cost(t) over already-hired t`, and is
actually debited once, on the simulation side, as the spawn runs. A hire never survives past
the mission it was made for.

<a id="the-price"></a>

## Price

```
cost(t, mission) = (PriceA[t] + n * PriceB[t]) * unitPrice(mission)
```

`n` is the number actually taken, which for a hire is the whole pool. `PriceA` is the constant
term, `PriceB` the per-head term; the five pool-1 types all ship `PriceB = 0`. `unitPrice`
depends on the **mission number only** — the type id is passed to the routine and never read.

## Levels

A mercenary's level is a step function of the mission number, and it selects a **different**
`Data.bin` template (`NPC03_1` … `NPC15_4`), not a scaling of one. All sixty names ship,
including the four for each of the two empty types, because the town roster builder constructs
every type unconditionally.

## Death and recovery

At the end of every mission, per type, in this order and before anything else is cleaned up:

```
if hired[t]:  pool[t] = (number of that type still alive)     # losses are permanent
else:         pool[t] = min(pool[t] + 1, MercenaryCount[t])   # one man back per mission
then:         hired[t] = 0 for all t
```

Then every mercenary is removed from the world; only heroes survive the cull. So a squad that
is wiped takes as many further missions to rebuild as it lost men, and only while it is left
at home — taking it out again freezes it at its current strength.

## Roster progress and selected-entry prerequisites

Command37 constructs types 1..15 independently of the eligible shelf. Its
fixed construction loops and empty-old-group conditional return do not prove
all constructors or callbacks complete. A missing owner group differs from a
valid empty group; cleanup needs a finite live chain. Identifier allocation
has an unbounded first dword scan, and null actor allocation can reach an
unchecked type store. Ordinary reach of those controls remains Unknown.
(SAV-926, SAV-927)

The client eligible collector can return empty. It still walks a nonempty
document map, so coherent buckets and acyclic stable chains matter even when
no type is eligible. Matching a supplied type does not validate its type-1
pool index. (SAV-928)

Activation copies the mercenary pointers, setting selection to 0 for nonzero
count and -1 for zero. The separate NPC list supplies no fallback selection in
that loop. Caption and price refresh admit -1 through signed upper-bound-only
checks; the local zero-storage controls read data[-1]. Preserving that state
through the complete intervening UI/resource work is Medium, and a later
selection write remains a live alternative. (SAV-929)

The separate party list must contain an entry marked with bit 0x20. Its helper
returns -1 if no marked entry exists; activation uses that index without an
absent-selection guard. Nonempty mercenary stock does not supply this party
prerequisite. (SAV-930)

A saved eligible set does not establish successful stock construction, list
activation or full click completion. The original first failing operation
across those paths remains Unknown. — SAV-931, SAV-932
