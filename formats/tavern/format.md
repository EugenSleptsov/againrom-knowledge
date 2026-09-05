# TAVERN — mercenary hire — specification (partial)

Level 3. Promoted, evidence-backed claims `MERC-TYPE-001`…`MERC-DEATH-006`. Ledger:
`claims/tavern.md`.

**Status: partial (◐).** The type space, the shelf gate, the hire, the price, the level, the
death/recovery rule and both commands are read at instruction level. **Not specified:** what a
mercenary *is* once spawned (that is `formats/databin/format.md`'s Humans table); the inn's own
art and layout; the multiplayer path, including whether a map-placed `Tavern` is reachable and
what mode makes `FUN_005050ee` return its flat 450.

This is not a file format. It is the second half of the town building whose mission half is
specified by `REG-SCN-064` / `REG-SCN-065`; it sits on `scenario.res::scenario.reg`,
`scenario.res::npc.reg` and `world.res:data/data.bin`.

## At a glance

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

## The type space

There are fifteen mercenary types and they are one index space shared by four files:
`[General] MercenaryCount`'s fifteen elements, `npc.reg`'s `[npc<t>]` sections, the runtime
object's own byte (`CUnit +0x15b` on the client, the server object's `+0x14c`), and the
`Data.bin` template name. Two of the fifteen ship empty — `MercenaryCount` is 0, no `PriceA` /
`PriceB`, no mission offers them, none unlocks them.

A consumer implementing this needs no separate "mercenary table": the type id **is** the npc
id, and `MercenaryCount[t-1]` **is** the headcount.

## What the tavern shows

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

An owner-preserved 55-file save corpus corroborates this boundary and its counterpart at the type-10
row above (`Humans[54] NPC10_1`) at population scale rather than one fixture. Every classKey-58 save
in that corpus sits at main mission 20 — where, as above, type 14 has no shelf at all — and matches
the mission-20 transfer's own signature exactly (three actors, typeWord `0x0a`, a bundled
classKey-201/`M10_Merchant` identity); a hire is excluded at every one of those states for the same
reason the shelf is empty. Every classKey-54 save in the same corpus instead sits on side-mission map
`41.alm` — its own selected-mission field reads 41, not the campaign record's main-mission value of
50 — which places exactly three `Humans[54]`/`NPC10_1` actors directly under the Player's own group;
no corpus row's own permanent-unlock array, including these three, ever admits type 10, so a hire is
excluded here too, on the same shelf-filter mechanism this page's own type-10 row states above. The
map's own authored placement, not a hire or a transfer, accounts for this class. — SAV-622, SAV-623,
SAV-624, SAV-625, SAV-626, SAV-627, SAV-628, SAV-629

## Hiring

The hire button toggles `hired[t]`, one flag per **type**. Leaving the inn sends the vector
`pool[t] * hired[t]` for `t = 1..15` together with the mission number, and the simulation
spawns that many of each type. **There is no way to hire one man of a type with four.**

Money is checked locally against `player.money − Σ cost(t) over already-hired t`, and is
actually debited once, on the simulation side, as the spawn runs. A hire never survives past
the mission it was made for.

## The price

```
cost(t, mission) = (PriceA[t] + n * PriceB[t]) * unitPrice(mission)
```

`n` is the number actually taken, which for a hire is the whole pool. `PriceA` is the constant
term, `PriceB` the per-head term; the five pool-1 types all ship `PriceB = 0`. `unitPrice`
depends on the **mission number only** — the type id is passed to the routine and never read.

A witnessed hire's debit confirms this formula end to end against one type's own already-published
price fields, `n` equal to that type's whole hired pool (SAV-614).

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
