# Mission identifiers and entry/exit

[Reference](format.md)

## Identifier bands

`Target_Unit` is three id spaces, and a consumer that treats it as one will fail to resolve four
fifths of the shipped references:

```
value <  10001    a unit id from the map's own type-6 records
10001..11000      a hero ordinal, resolved against the live player list at run time
                  (and unconditionally unresolvable when the multiplayer flag is set)
value >  11000    an index into a static name table in the executable
```

## The drop table — where the player lands

The authored action `0x10002` is never dispatched. At build time its `value[0]` and `value[1]` are
truncated to bytes, packed `(y << 8) | x`, and appended to a `CWordArray` at scriptObj`+0x28` =
`mapObj+0x6c`. **Every installed map carries one** — campaign and skirmish alike, and
on 9 of the 10 loose maps it is the map's only script node.

`FUN_004d403c` consumes it, and three details decide whether a reimplementation lands the player
where the engine does (`MISSION-DROP-002`):

```
if (mapObj+0xc != 0 && (u16)player+0x60 != 0)   x,y = low,high byte of player+0x60
else if (array.GetSize() > 0)                   i   = FUN_00504003(GetSize()-1)   <- RANDOM index
                                                x,y = low,high byte of array[i]
if (x * y == 0)                                 x   = FUN_00504003(0x46) + 0x1e   <- 30..100
                                                y   = FUN_00504003(0x46) + 0x1e   <- independently
                                                log "no drop location in .alm - random used"
```

`FUN_00504003(n)` is `(rand() * (n+1)) / 32768` — **inclusive of `n`**, and not a modulo. The
array is **not** read at `[0]`.

## Starting the mission

The map places nobody for the player. Roster slot 1 owns no type-6 record on any campaign map,
and `FUN_004d403c` iterates the **player's own unit list** at `player+0x20`, placing the hero
(`player+0x34`) at the drop cell exactly and everything else within
`ftol(max(5.0, sqrt(nUnits) + K))` of it. A consumer that builds the party from the map file
starts every mission empty (`MISSION-START-001`).

The routine's second half — a by-name lookup of a node called **`"Humans"`** whose children are
`label#id` strings with their own `(x, y, radius)`, or `param[1] == -1` meaning "at the drop
cell, radius 8" — is the `.ini` overlay's entry point below, and **no shipped map has such a
node**, so on shipped data it never runs.

## Ending it

`instant 4` is the only thing that wins, and **every campaign map carries exactly one**
(`MISSION-WIN-003`). None of the ten loose maps carries any: a skirmish map cannot be won,
only lost. A win action may be referenced by more than one trigger — four on `81.alm` — so
bind by node id.

`check 18` is the "protect this unit" objective and it is authored as a trigger with **no
action at all**: the check writes no slot, so the pattern never passes, and its only effect is
the lose it raises when its unit dies. A check is armed by being **authored**, not by being
referenced — the builder gives every check node a slot and the pass evaluates every check
(`MISSION-VIP-004`).

## The register file is shared

The build assigns check node *i* the slot *i*, constants included, and an authored **variable**
addresses the same `session[0xbd34 + p0*4]` array with its own literal number. Any variable index
below a map's check-node count is overwritten every full tick. `60.alm` ships exactly that
collision — variables 32 and 33 against 64 check nodes (`MISSION-SLOT-008`).

## Mission text

`instant 2` carries a **number**, not a string. `FUN_00473110` opens
`main.res::text/battle/m<mission>/event<NN>.txt`; the same block holds `briefing.txt`,
`briefmap.txt`, `title.txt` and `tips<NN>.txt`. The files are markup with `<NPC=n,Part=k,…>`
tags whose `n` is an `npc.reg` section number (`MISSION-TEXT-005`). **A number whose file does
not ship is a silent no-op** — the window is never built and nothing is logged; the window
itself, its lifecycle and the tag vocabulary are [DIALOGUE](../dialogue/format.md)
(`DLG-WIN-001`…`DLG-MARKUP-007`).

## `World\Mission\<n>.ini`

Built and opened at every map load, absent from the install, and its absence is a **silent
no-op** — the open fails, the reader returns immediately, and the one call site discards the
result. It is a sectioned text overlay whose sections the engine looks for by name:
`Humans.Hero` (`Name`, `Bag.`, `Armor`, `Weapon`, `Shield`, `Item`), `Outposts`, `Patrol`,
`StandGround`, `Items`, `Players`, `Mission`, `Monsters`. **Nothing in the trigger machinery reads
it**; its consumers are hero import, patrol paths, outposts and item setup.
