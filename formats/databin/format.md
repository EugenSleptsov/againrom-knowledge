# DAT `Data.bin` — the placeable-definition database — specification (partial)

Level 3. Promoted, evidence-backed claims only. The populator, grammar, streamed
slot maps and corpus closure are `DAT-LOC-001`…`DAT-ACT-006`. Ledger:
`claims/databin.md`.

**Status: partial (◐).** The file, the wire grammar (exact tiling, 0 residue), the
schema law and the three placement-facing tables (Buildings, Humans, Units) are
specified with their consumer paths; the other eight tables are framed (names,
counts, param arrays decoded) but their slots are named only by the shipped column
titles — no consumer of them has been read.

Seen as: `world.res:data/data.bin` (88 327 B in the shipped install), or a loose
`World\Data\Data.bin` which takes precedence. If neither exists, `rom.exe` parses
eleven `;`-separated `.csv` tables from the same directory and **rewrites**
`Data.bin` (`DAT-LOC-001`) — the `.bin` is a serialized image of the CSVs, column
titles included.

## At a glance

Eight class-groups, each `[column-title string array][1..3 collections]`; a
collection is `[u32 count][entries]`. Groups C–H write entries `1..count-1` — entry
0 is a reserved null, which makes every consumer index **1-based**.

```
group A  titles(11)  Shapes(5, all)      Materials(16, all)     entry = [name][9 doubles]
group B  titles(30)  Magic(50, all)                             entry = [name][params]
group C  titles(18)  Armors(30) Shields(9) Weapons(27)          entry = [name][params][10 raw][dwords]
                     the 10 raw = five u16 material masks, one per Shapes row  DAT-MATMASK-020
group D  titles(4)   MagicItems(49)                             entry = [name][params][1 raw][string]
group E  titles(57)  Units(118; 56 parameterised)               entry = [name][params][2 strings]
group F  titles(28)  Humans(215; 210 parameterised)             entry = [name][params][10 strings]
group G  titles(9)   Buildings(66)                              entry = [name][params]
group H  titles(24)  Spells(28)                                 entry = [name][params][1 string]
```

Wire primitives: `CString` = u8 length (0xFF → u16) + bytes; a title array = u16
count + CStrings; a param array = u16 count + raw u32 little-endian values; the
collection count is a plain u32. Counts above are the shipped file's stored entries
(`DAT-GRAM-003`).

## The schema law

The complete per-collection column list, re-expressed as slots and carrying each
streamed slot's actor destination, is promoted as `DAT-SCHEMA-007`.
Two things a raw title dump does not show: a group's title array is **one array
shared by its sibling collections** (Armors + Shields + Weapons share 18 titles;
Shapes + Materials share 11, and that element kind has **no param array at all** —
its record is 9 doubles), and the param arrays are **full** — 56/56 Units rows carry
55 values, 210/210 Humans rows carry 26.

Param slot `i` is CSV column `i+1` (column 0 is the entry name); the trailing
"equipment" column, when the class has string slots, is `,`-separated with `{...}`
groups and fills the entry's extra strings instead. **An empty cell is stored as −1,
and the engine's readers skip the store on −1**, so −1 always means "constructor
default" (`DAT-SCHEMA-004`, `DAT-ACT-006`).

## The placement-facing tables

- **Buildings** — index = `structures.reg` `ID` (1..66; `DAT-BLD-005`). Slots:
  `0 sizeX, 1 sizeY` (footprint in tiles — overridden by the ALM type-4 `kind==0x21`
  extension), `2 scanRange → obj+0x48`, `3 healthMax → obj+0x44`, `4 Passability`,
  `5 BuildingPresent` (title-named only). Resolved from an ALM type-4 `kind` by
  direct 1-based subscript (`ALM-CLS-036`/`ALM-CLS-053`), or by name
  (`FUN_005241d0`, backwards, `"Invalid building %s created"`).
- **Humans** — searched on slot `0x10 typeID` for type-6 keys `< 0x40 ∉ {26,27}`,
  and on slot `0x18 serverID` for the type-6 `+0x10` override id. Streamed at spawn
  (`FUN_004f974d`): stats/HP/mana/speed/skills, `0x15 TokenSize → actor+0x49`,
  `0x16 MovementType → actor+0x4a`; ten equipment strings resolved by name against
  Weapons/Shields/Armors; `0x19 knownSpells` a bitmask. Slot `0x10` streams to
  `actor+0x0e`. `FUN_004f9065` overwrites it with `gender+0x21/+0x23` only when the
  constructor-mode argument is non-zero. ALM definition-id and explicit typeID arms
  pass zero; the npc arm passes the exact `Hero` flag result (`PARTY-M20-031`). A map
  Human can therefore retain its authored table typeID.
- **Units** — searched on slots `0x1d typeID` + `0x1e face` for all other type-6
  keys. Streamed at spawn (`FUN_004f5604`): slots **0–37 only**, in order — stats,
  the two regeneration periods, the damage pair through the `attackKind` switch,
  protections, resists, `0x1f tokenSize → +0x49`, `0x20 movementType → +0x4a`
  (shipped values 2 = Ghost/Bee, 3 = Bat_Sonic/Dragon — the `0x44`/`0x82` block
  masks). The **treasure and spell slots (38–54) are not streamed**: the spell pairs
  are read by `FUN_004f59de` into the spellbook and the order block
  (`UNIT-SPELL-007`) and the treasure columns by the kill payout (`HERO-KILL-027`).
  Complete slot → actor-field map with widths and instruction addresses:
  [`formats/unit`](../unit/format.md), `UNIT-STREAM-001`.

Corpus closure: every shipped placement resolves — 3141/3141 type-4, 8079/8079
in-scope type-6 (`ALM-CLS-052`).

## Open

The Shapes/Materials 9-double record's slot meanings; the armor 10-byte block and
second dword array; Magic/Spells/MagicItems semantics; the `Start ID`/`Tiles`
building columns (no stored param); the server→client creation-message encoding; and
how `FUN_0050d670` parses an equipment cell's tier and material words — the name
grammar `[<Shapes> ][<Materials> ]<Weapons>` closes 26/26 over the shipped Units
`EquipItem` column but is a corpus fit, not a read (`UNIT-EQUIP-005`).

**A second shipped sample exists.** `gameversions/ru/WORLD.RES` differs from the EN
`world.res`, and so does its `Data.bin` — but the **Units** collection is identical
in every name, string and parameter, while **Humans** differ on 16 parameters over
6 slots. Any corpus figure taken from the Humans table must say which root it read.
## Buildings presence-mask consumer

The Buildings presence mask is consumed as `(row*width+col)&31`, with byte width/height and a
32-bit mask; it is not an arbitrary-size bitset. Registration visits only its set positions and
can stop at the first occupied Building slot. The declared rectangle therefore does not prove
that all selected cells attached (`UNIT-STRUCTCELL-070`). Both preserved Data.bin inputs yield
66 identical numeric Buildings rows, but the corresponding installed placement populations
differ (`UNIT-AREAPOP-075`).

## Campaign-start document item

`MagicItems[28]` is named `Quest Documents` on both roots. Its parameter array is `[-1,5]` under
the titles `Magic Items`, `Price`, `weight`, `Effects`; price is therefore -1. Its packed item code
is `0x0e1c`.

The `Humans` table cannot place this MagicItems row in starting equipment. Cell 0 constructs a
Weapon, cell 1 a Shield, and cells 2..9 Armor. Across all 215 stored rows in each preserved root,
the 909 non-empty cells contain no `Quest` string (`DAT-DOC-021`).
