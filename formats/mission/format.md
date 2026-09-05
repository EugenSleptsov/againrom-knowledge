# MISSION — starting a campaign mission and ending it in a win — specification (partial)

Level 3. Promoted, evidence-backed claims only. Ledger: `claims/mission.md`.

**Status: partial (◐).** Read at instruction level: the player-placement routine, the drop-cell
decision and its RNG, the seat search and its failure, the four placement arms and their branch
order, the definition lookup, the npc→definition hop, the message-text path, and the entry gate
and two-container writer boundary (`PARTY-ORIGIN-010`…`PARTY-GATE-013`). Not specified: the arithmetic behind the
`DataBinID == 26` sentinel; whether the compiled check list is the whole node list or the validated
subset; which arm the campaign state machine drives between two missions.

**`MISSION-DROP-002` corrects a base name used throughout the older text: `FUN_004d403c`'s
`this` is the server singleton `[0x005cd758]`, not a map object.** Where this spec said `mapObj+0xc` read
`server+0x0c`, and `mapObj+0x6c` is `server+0x6c` — the same memory as `(server+0x44)+0x28`,
because `server+0x44` *is* the map sub-object. Numbers unchanged, object different.

This is not a file format. It is the seam between three that are: the map (`formats/alm`), the
script runtime (`formats/trigger`) and the class definitions (`formats/databin`, `formats/reg`).

## The one thing a consumer must not get wrong

**The map does not contain the player.** Not the hero, and — on 23 of the 28 campaign maps — not
one unit for him either; the five exceptions supply **19 units in total** and the walk moves even
those to the drop cell. What a map contributes is **one packed cell**, and the engine's job at load
is to take units that already exist on the player object and stand them on it. A consumer that
looks for the player's units among the map's placements starts every mission with an empty army.

**And "the player object" holds two containers, not one** (`PARTY-ORIGIN-010`):

```
Player +0x20  -> a flat index of every actor I own      <- the placement walk reads THIS
       +0x24  -> a list of groups, each with its own actor list   <- the SAVE stores this
       +0x34  -> the hero; null here means the mission is refused outright
```

`+0x20` is derived: `Player::Serialize`'s load arm rebuilds it by walking `+0x24`. Every routine
that gives the player a unit writes both, plus a fresh group, plus `actor+0x14` and `actor+0x70`.

## Starting

```
0. gate                    FUN_004d303e   ; session opcode 0x04. player+0x34 == 0 -> REJECT,
                                          ; "You can't enter mission without Hero", nothing placed.
                                          ; then a once-latch, byte player+0x3d, which is SAVED
1. compile the script      FUN_004e3591   ; also appends the 0x10002 cells to server+0x6c
2. group hygiene           FUN_004d403c   ; no groups -> make one ("Oops - player has no groups");
                                          ; then drop empty groups off the head, keeping the last
3. decide the drop cell    FUN_004d403c   ; player+0x60 override IF server+0x0c != 0 -- which the
                                          ; campaign path never is -- else a RANDOM array element,
                                          ; else rand 30..100 on each axis + a logged warning
4. position the party      FUN_004f4604   ; walk player+0x20 in list order; hero (player+0x34)
                                          ; exact at r=0; the rest within
                                          ; r = ftol(max(5.0, sqrt(nUnits) + [0x0059bc50]))
5. import the .ini's own units             ; node named "Humans"; no shipped map has one
```

The order matters: step 3 reads what step 1 wrote, through the same object
(`(server+0x44)+0x28` == `server+0x6c`).

**Seating can fail, per unit, and the mission starts anyway.** `FUN_004f4604(x,y,r)` makes
`r·r/2 + 2` random attempts inside `[x ± r/2] × [y ± r/2]`, then — only if `r > 0` — rasters that
whole square, x outer, first free cell wins. If nothing is free it returns 0, the engine logs
`"Error - can't place hero from previous mission on map."` **for any member, not just the hero**,
and the walk carries on: the surplus stays in the collections and off the map. Nothing anywhere
compares the party size against what the box can seat.

**What writes `player+0x60`** (the old open item): exactly one instruction, `004d6e1b`, in session
opcode `0x32`. It packs the *current cell* of one of the player's own units as `x | (y<<8)` and
stores it on that unit's owning `Player` — "start the next map where this unit is standing". On the
campaign path `server+0x0c` is 0, so the value is written and never read.

## Ending

```
instant 4  ->  win     exactly one authored node per campaign map, 28/28, and all 28
                       are reached by a trigger
                       zero on all ten loose maps -> a skirmish map cannot be won
instant 5  ->  lose    authored on 15 of 28, reached by a trigger on 12 of those
                       (was published as 14 of 28; MISSION-TYP-028 corrects it)
check   18 ->  lose    when its unit is dead; writes no slot, so the trigger that
                       carries it usually has NO action at all
                       authored on 7 EN maps (10 nodes), 6 RU maps (9 nodes)
```

A win node may be referenced by more than one trigger — four on `81.alm`, two on `130.alm` — so
resolve `THEN` slots by node **id**, never by position.

**The two lose arms are armed differently, and a consumer that treats them alike will get one
of them wrong.** An instant-5 node fires only when a trigger's action list reaches it. A
check-18 node is armed by being **authored**: the builder gives every check node a slot and the
evaluator runs every check once per full tick, so an unreferenced check-18 node still loses the
mission when its unit dies (`MISSION-VIP-004`). The test for "can this map be lost by its
script" is therefore *referenced* instant-5 nodes plus *authored* check-18 nodes.

Eleven of the 28 campaign maps have neither, on both roots — `30`, `41`, `51`, `61`, `80`, `90`,
`120`, `121`, `131`, `140`, `141` — so their scripts can only ever produce a win
(`MISSION-LOSE-025`). Independently of scripts, primary-character fall reaches the campaign
failure branch under the participant gates in `MISSION-DEFEAT-045`.

## `30.alm`, worked end to end

The smallest complete campaign mission read so far, and a useful conformance target: its whole
executable surface is five instant arms and four check arms.

```
build time   trig[0] has a zero left id in its first condition pair -> TRIG-BIND-010 drops
             the whole trigger. Its drop-location node still reaches the start, because
             the drop table is built from the node array, not through the trigger.

start        trig[7]  constant 0 <= constant 0        (always true, once=1)
                 do   instant 12  Unit=10001 Item=6   create item, add to hero ordinal 1

win          trig[6]  check 6 (unit 10001, unit 56) <= constant 3      (once=1)
                 do   instant  2  message 16
                 do   instant 13  Unit=10001 Item=6   resolve and destroy
                 do   instant 13  Unit=10002 Item=6   resolve and destroy
                 do   instant  4  WIN

lose         none reachable: one instant-5 node, referenced by nothing; no check-18 node

authored     trig[1] and trig[3] read the SAME check node (group 8 population == 0),
defect       so their two messages fire together. The group-5 check node the second
             trigger is named after is read by nothing. A consumer that "repairs"
             this diverges from the shipped map.
```

Check 6 is the Chebyshev metric of `TRIG-DIST-014` and answers `0xff` when either unit is dead.
Unit 10001 and 10002 are the first and second hero ordinals (`TRIG-REC-011`); the second exists
because `[Mission30] AddHero=22` inserts a companion before the mission starts. Unit 56 is a
placed NPC at cell (65,15).

**Nothing tests the created item.** The vocabulary has check 17, `Item in inventory`, and the
campaign authors 23 of them; this map authors none. The item is given, carried and destroyed,
and the win depends only on where the hero is standing (`MISSION-CURE-026`).

## Who is on the map, and out of which file

```
+0x08 >= 0x1a                          data.bin Units,  params 0x1d/0x1e     6672 / 8094
+0x08 <  0x1a and +0x0c bit 0          npc.reg [npc<+0x0a>].DataBinID
                                       -> data.bin Humans, param 0x18          15 / 8094
+0x08 <  0x1a, bit0 clear, +0x10 != 0  data.bin Humans, param 0x18 = +0x10   1405 / 8094
+0x08 <  0x1a, bit0 clear, +0x10 == 0  data.bin Humans, param 0x10 (typeID)     2 / 8094
```

The npc test is **outside** the definition-id test. A record with both takes the npc arm and its
own `+0x10` is dead. Two traps in the lookup:

- `FUN_004de63e` walks the definition collection **backwards from the last entry**, never tests
  index 0, and **returns 0 on a miss** — so an unresolved id silently becomes entry 0.
- `FUN_0048cbe0` returns `npc+0x14` (`DataBinID`) **unless it is 26**, which is a sentinel meaning
  "compose the template from this npc's own `Flags` tokens" — `Mage`, `Female`, `MySex`,
  `MyClass`. `[npc21]`, `Flags="Hero,Me,Start"`, is the player.

## Dialogue

`instant 2` carries a number. The engine opens
`main.res::text/battle/m<mission>/event<NN>.txt`, and the same family gives
`briefing.txt`, `briefmap.txt`, `title.txt`, `tips<NN>.txt`. The files are markup with
`<NPC=n,Part=k,…>` tags; `n` is an `npc.reg` section number, the same id space the npc placement
arm uses. Over the campaign, 223 of 242 raised numbers name a shipped file. **The 19 that do not are
silent**: `DLG-ABSENT-003` establishes that the failure path has no fallback, no message and no state
change (`formats/dialogue`). Three of the 19 are English-only — the Russian root ships
`m100/event09`, `m130/event07` and `m150/event10` from identical scripts.

## The register file is one array

Check node *i* gets slot *i* (constants included) and an authored **variable** indexes the same
`session[0xbd34 + p0*4]`. A variable number below the map's check-node count is clobbered every
full tick. `60.alm` ships that collision.

## `scenario.res`'s three registries

Beside the 28 maps: `scenario.reg` (the per-mission campaign parameters — shops, inn, mercenary
pool, payment, `AutoGetMission`, `LastMission`; `REG-SCN-059`), `npc.reg` (110 sections; the
placement arm above, the four chargen archetypes, and the `<NPC=n>` tag space; `REG-NPC-058`),
`globalmap.reg` (the world map and its mission→object mapping).

## Money across the boundary

Each `Player` owns one 32-bit purse at `+0x38`. It is not a hero field. A save writes it in the
campaign half as `u32(value ^ 0x5c073f4d)`, and the character carry transports it separately from
the actor graph as the first of four `Player` dwords. A mission transition must therefore preserve
the participant's purse even when it rebuilds or re-homes the hero (`PARTY-MONEY-016`).

`scenario.reg`'s per-mission `Payment` is passed on the completion path, not at mission entry
(`MISSION-MONEY-022`). The last callee on that path remains unread, so the specification does not
yet state its sign or recipient as a High-confidence purse operation. Both roots carry 24 mission
sections and nine non-zero values, totalling 1,289,700.

## Added companion at mission 30

Mission 20 makes mission 30 current before the town view activates. The first activation consumes mission 30's `AddHero = 22` and clears the array. The companion therefore exists before mission 30 starts.

Command `0x49` assigns the existing player as owner, inserts the actor into the player's flat actor index, creates a group, inserts that group into the player's group list, and inserts the actor into the group. It does not replace the primary character pointer and does not use the mercenary path. The companion is a persistent player character with separate actor state and inventory. Save persistence follows from the serialized group list. Complete mission-to-mission in-memory persistence remains Medium.
## New-campaign state at mission 10

The new-campaign path constructs the session, resets the campaign record, and loads mission 10
before the character-generation command creates the hero and before the first simulation tick.

The human participant's `Player+0x38` purse is a 32-bit zero from `Player::Player`. Character
generation, mission 10 entry, and the first tick do not change it on the measured path. Fighter or
mage, sex, face, stats, entered or default name, and starting weapon do not affect the purse. A
runtime observation remains required for the complete negative: any non-zero purse before a player
action refutes it (`PARTY-MONEY-018`).

Mission 10 appends three text-document entries to the campaign record before hero creation:
`(value=1,kind=1)`, `(2,1)`, and `(3,1)`. The collection is independent of the inventory item that
opens its panel. Hero construction calls a producer that appends `Quest Documents` to the hero only
when `server+0x0c == 0` and `[0x005eb5a4] != 2`, before runtime-id assignment. The first gate holds
on the campaign construction path; the second was not measured in a fresh campaign. The document
collection is present regardless, while actual access-item presence remains Unknown
(`MISSION-DOC-023`, `ITEM-DOC-054`, `ITEM-DOC-069`).

## Defeat: primary character, mode and frontend

The reporter requires an active human participant (`player+0x3d != 0`, `+0x28 == 0`).
Latch `+0x3c >= 2` takes a recovery branch and cannot fall through to script counters.
Actual mode `server+0x0c == 0` returns without resetting that latch. Nonzero mode additionally
requires `player+0x3f != 0` and primary signed HP below -53 before repair, placement and latch
zero. All shipped campaign map loads and campaign save entry use mode zero; the mode is not
equivalent to archive membership or transport ownership (`SESS-DEFEAT-064`).

Below latch 2, absent primary `player+0x34` or nonzero primary stage `+0x13c` wins over scripts.
The branch sets latch 2, sends `0xb4` only with joined byte `player+0x3e != 0` and mode zero,
and sets all actors in the flat ownership index to HP -50. It does not test whether any companion
is alive. Otherwise script loss `== 1` precedes success `== 1`, without a mode predicate.
Win-to-loss is locally admissible; loss-to-win merely because loss count exceeds 1 is not
(`MISSION-END-013`, `MISSION-DEFEAT-045`).

The actual failure panel uses vtable `00598b78`, handler `00447063`, and these two actions:

| Label (EN / RU) | Result and route |
|---|---|
| Exit to Main Menu / Выйти в главное меню | `0x445 -> 0x44c -> 0x41e`, session teardown, `0x421`, then main menu when screen word is zero |
| Load Game / Восстановить игру | `0x446 -> 0x44c -> 0x418`, save selection; disabled when no `game*.sav` is found |

Neither is Restart, Continue or a dead-primary repair command. The neighboring class's
`00446d35 -> 0x41d` is not this panel (`MISSION-PATH-015`, `MISSION-VICTORY-035`,
`MISSION-DEFEAT-046`). The shown panel sets screen bit 8, which pauses campaign idle stepping
even while the run dword stays 1 (`MISSION-STOP-016`, `SESS-DEFEAT-065`). Runtime delay, packet
backlog and ordinary-UI delivery of a repair command after fall remain Unknown.
