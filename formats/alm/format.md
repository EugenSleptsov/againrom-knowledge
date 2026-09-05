# ALM map container (`M7R␀`) — specification

Level 3. Promoted, evidence-backed claims only. Sources: `ALM-HDR-001`,
`ALM-SEC-002…ALM-SEC-004`, `ALM-LOC-007`, `ALM-FRAME-031` (container);
`ALM-META-008…ALM-META-010`, `ALM-GRID-012…ALM-GRID-014`, `ALM-GRID-032`
(type-0 metadata and grid-layer layout); `ALM-META-024…ALM-META-028` (the `rom.exe`
map-load routine and type-0 field meanings); `ALM-CNT-017`, `ALM-UNIT-018`,
`ALM-OBJ-019`, `ALM-GRP-020`,
`ALM-TRIG-021`, `ALM-TRIG-022`, `ALM-CODE-023` (the variable content sections type4–9);
`ALM-GRID-012…ALM-GRID-014`, `ALM-TERR-015`, `ALM-TERR-016` (grid-layer **semantics**
from `rom.exe` and `world.res:data/map.reg`); **`ALM-FRAME-031`, `ALM-GRID-032`,
`ALM-PLACE-033` (the corrected record framing, grid base and placement→cell conversion)**;
**`ALM-OBJ-034`, `ALM-CLS-035…ALM-CLS-038`, `ALM-OWN-039`, `ALM-UNIT-040`,
`ALM-GRP-041`, `ALM-CLS-042` (the placed-content → class binding)**; and
**`ALM-REQ-055`, `ALM-REQ-056`, `ALM-ORD-057`, `ALM-META-058`, `ALM-RDR-059`,
`ALM-CORP-060` (the acceptance contract, absent-record behaviour and four readers)**.
Corpus: the EN root's 10 loose `.alm` maps **+** the 28 embedded
in `scenario.res` = 38 maps, extended by `ALM-CORP-060` to both preserved installs — the RU
root's 6 loose files and its own 28 embedded, 72 walked files
over 44 distinct names; 0 tiling violations; grid cell-count exact 38/38;
content-section counts + coordinate-in-bounds tests exact 38/38; falsification passed.
**Decoded:** the container framing, the type-0 metadata layout + field meanings, the
type1 (**Tiles**) / type2 (**Altitudes** = height) / type3 (object occupancy overlay)
grid layers including the grid base, the tile→terrain resolution and the derived runtime
passability grid, and the type4–9
content-section record layouts (units, objects, groups, the trigger effect/Drop list, and
the marker trees), **and — `ALM-OBJ-034`, `ALM-CLS-035…ALM-CLS-038`, `ALM-OWN-039`,
`ALM-UNIT-040`, `ALM-GRP-041`, `ALM-CLS-042` — the class binding: the type-3 code, the type-4
`kind` and the type-6 `+0x08`/`+0x0a`/`+0x10` all resolve to a graphics-registry class,
plus the owner field both record types share and the full type-5/type-6 field maps.**
**Decoded by `ALM-TRIG-044…ALM-TRIG-047`, `ALM-UNIT-048`, `ALM-TRIG-049`,
`ALM-TRIG-050`** (all Medium — corpus closure, no consuming instruction read):
the type-7/8/9 record grammars below; the six comparison codes and the trigger's
`+0xb4` word are since resolved at instruction level (`TRIG-CMP-006`, `TRIG-FIRE-007`;
see `formats/trigger`). **Still Unknown:** the trigger's 64-byte junk field, the
`Target_Item` domain, and the type-9 `+0x00` tag word (the other type-9 header fields
are carried by `ALM-EFFREC-071`…`ALM-EFFPOP-073` and `UNIT-M10CELL-054` below).

Seen as: `*.alm` / `*.ALM` at the install root, and as `N.alm` file-node payloads
inside `scenario.res` (an `&YA1` archive of 28 campaign maps).

**Selected EN/RU member equality (`SAV-SUFF-301`).** Read-only payload hashing finds all six
ALMs selected by the current SAV population byte-identical across the preserved roots:
`scenario.res/10.alm`, `20.alm`, `30.alm`, `31.alm`, `40.alm` and `41.alm`. This is
exact equality for that selected subset only. It does not establish equality of the
other ALM members, runtime-derived state, or compatibility when SAV and external roots
are crossed. — SAV-SUFF-301

> **Changed by `ALM-FRAME-031` — if you implemented the superseded framing, re-read this box.**
> The earlier record split was **8 bytes too early**. The file header is
> **20** bytes (not 12), each record header is `[tag][hdrLen][payloadSize][typeId][f32]`
> (the `typeId`/`f32` belong to the *header*, not the payload), there are no per-section
> `f0`/`f1` fields, and there is **no 8-byte trailer**. Consequences: **every payload offset
> in the old spec reads −8**, and the three grid layers begin 8 bytes later — a decoder on
> the old base renders **4 cells off in X on type1 (u16)** and **8 cells off on type2/type3
> (u8)**. `ALM-FRAME-031`, `ALM-GRID-032`.

## At a glance

A 20-byte file header and a chain of **`recordCount` length-prefixed, typed records** — 10 in every
authored map, but the count is a header field the loader gates only `>= 3`, and one shipped file
carries 4 (`ALM-CORP-060`). The records tile the file with no gaps, no overlaps and no trailer.

```
+----------+-------------------------------------------------+
|  header  |          recordCount records (chain)            |
|  20 B    |   each: 20-B header + payloadSize B payload      |
+----------+-------------------------------------------------+
0        0x14                                               EOF
```

**Required: types 1 and 2 only.** See *Acceptance contract* before implementing a reader.

## File header (20 bytes)

| Off | Type | Field | Notes | Claim |
|-----|------|-------|-------|-------|
| 0x00 | char[4] | magic | `4D 37 52 00` = `"M7R␀"` — the **only** field the loader validates | ALM-HDR-001 |
| 0x04 | u32 | **hdrLen** | `= 20` — this header's **own length**, and it is *used*: the header helper peeks 8 bytes, seeks back, then reads `dword[cursor+4]` bytes. At offset 0 that reads 20 | ALM-HDR-001, ALM-FRAME-031 |
| 0x08 | u32 | dataSize | `= 4·W·H + 72` on 38/38, but the loader never reads it and the corrected framing does **not** explain the 72 (see notes) | ALM-HDR-001 |
| 0x0C | u32 | **recordCount** | the loader's loop bound, gated `≥ 3` and nothing more; `= 10` in every authored map, `= 4` in `ru/Horror.alm` | ALM-META-024, ALM-FRAME-031, ALM-CORP-060 |
| 0x10 | u32 | **formatVersion** | `= 990`; gated `≤ 1001`; selects version-conditional record fields at `0x3b6/0x3d8/0x3da/0x3dd/0x3e9`, and `== 1000` would skip record headers entirely (unexercised by any shipped map) | ALM-META-024, ALM-FRAME-031 |

`W`,`H` themselves live in the type-0 record payload (see below), not the header.
`rom.exe` `FUN_00512353` reads this header as one `Read(dest, 0x14)`.

## Record header (20 bytes) — repeated `recordCount`×, from 0x14

| Off | Type | Field | Notes | Claim |
|-----|------|-------|-------|-------|
| +0x00 | u32 | tag | constant `7` on all 380 shipped record headers; stored to the map object (`+0x08`) and never read again | ALM-SEC-002, ALM-FRAME-031 |
| +0x04 | u32 | hdrLen | constant `20` (= this header's own length); read, not validated | ALM-SEC-002 |
| +0x08 | u32 | payloadSize | payload byte length; the payload follows immediately | ALM-SEC-002 |
| +0x0C | u32 | **typeId** | the record's type, `0..9`; **this is the word the loader's `switch` dispatches on** (10-entry jump table) | ALM-SEC-003, ALM-FRAME-031 |
| +0x10 | f32 | **perMapConst** (`selectorA`) | a per-map `f32` (e.g. `0xBFC02B6D`), byte-identical in every record header of a map; 3 discrete values across the corpus | ALM-SEC-003, ALM-META-027 |

The payload is `payloadSize` bytes of **pure data** — nothing is overlaid on it and no
identifier precedes it. `rom.exe` `FUN_00512332` reads this header as one
`Read(dest, 0x14)` (skipped iff `formatVersion == 1000`), then each `case` reads the
payload directly. The clinching corpus fact: **twelve shipped records have
`payloadSize < 8`** (an empty `type8` on eight maps, an empty `type4` on `scn:31`, a
4-byte `type9` on three), which is only well-formed if `typeId`/`f32` are header fields —
an 8-byte payload-head identity cannot fit in a 0- or 4-byte payload.

## Record roster — what the writer emits, and what the loader requires

**These are two different things (`ALM-SEC-003`, `ALM-REQ-055`, `ALM-REQ-056`,
`ALM-ORD-057`, `ALM-META-058`, `ALM-RDR-059`).** Authored maps emit 10 records in
typeId order **`0, 1, 2, 3, 5, 4, 9, 8, 6, 7`** (`ALM-SEC-003`), and 71 of the 72 `.alm` files in
the two preserved installs do. The loader requires **none of that**: it iterates the header's
`recordCount` (gated only `>= 3`), dispatches per record on `typeId`, and afterwards demands only
**type1 and type2** — see *Acceptance contract* below. A consumer that requires ten records rejects
`gameversions\ru\Horror.alm`, which the engine loads (`ALM-REQ-055`, `ALM-CORP-060`).

Sizes by type (`ALM-SEC-004`):

| typeId | payloadSize | Role |
|--------|-------------|------|
| 0 | 632 (constant) | **metadata** — decoded below (`ALM-META-008…ALM-META-010`) |
| 1 | `2·W·H` | **grid: tile word** (u16/cell) — decoded below (`ALM-GRID-012`, `ALM-GRID-032`) |
| 2 | `W·H` | **grid: relief/height** (u8/cell) — decoded below (`ALM-GRID-013`, `ALM-GRID-032`) |
| 3 | `W·H` | **grid: static-object placement layer** (u8/cell) — cell code `c ≠ 0` is `objects.reg` section index `c − 1` (`ALM-GRID-014`, `ALM-CLS-035`) |
| 4 | `20·#4` (+8·ext) | **placed structures/buildings** — 20-B records, `#4 = meta+0x20` (the loader's own case-4 loop bound); every record is a `structures.reg` class (`ALM-OBJ-019`, `ALM-PLACE-033`, `ALM-CLS-036`) |
| 5 | `76·#5` | **player/group roster** — 76-B named records, `#5 = meta+0x1c` (`ALM-GRP-020`, `ALM-PLACE-033`) |
| 6 | `70·#6` | **placed units** — 70-B records, `#6 = meta+0x24` (`ALM-UNIT-018`, `ALM-PLACE-033`) |
| 7 | variable | **trigger effect/instant list + Drop table** — `[u32 entryCount][named blocks]` (`ALM-TRIG-021`) |
| 8 | variable, sometimes **0** | **counted geometry records**, `meta+0x2c` of them, no count word — *not* a tree and *not* the box/circle geometry (the `ALM-TRIG-022` gloss is refuted; see below and `retracted.md`) |
| 9 | variable (small) | **tile-marker list** — `[u32 count@+0][records]` (`ALM-CODE-023`, `ALM-TRIG-049`) |

The three fixed-record sections carry their **count in the type-0 metadata**
(`+0x1c`=type5, `+0x20`=type4, `+0x24`=type6, `ALM-CNT-017`); the variable trigger
sections carry their own count word. This is a second, independent proof of the
76/20/70 strides.

Also in the metadata, and load-bearing for a *writer*: `+0x28` and `+0x2c` are the loader's case-7
and case-8 loop bounds, so all five of `+0x1c`..`+0x2c` must agree with the records actually present
— a zero-length record needs a zero count or the sequential parse desynchronises. Shipped data obeys
this exactly: 15 zero-length records over 72 files, 15 zero counts, 0 disagreements (`ALM-REQ-056`).

## Acceptance contract — what a reader must and must not require

Placed-structure footprint registration must not be reused as the physical-attack token
size. The base Building attack vtable supplies size 1 even though `+0x60/+0x61` carry its
registered rectangle (`UNIT-STRUCTREACH-063`). Its destructor detaches the footprint, but
the physical strike's HP write does not itself invoke that destructor. The HP-to-destruction
scheduler and derived-class overrides remain Unknown in this bounded route
(`UNIT-STRUCTSTOP-066`); no new type-4 byte layout is implied.

Read from `FUN_00512369` and its consumers (`ALM-REQ-055`, `ALM-REQ-056`, `ALM-ORD-057`,
`ALM-RDR-059`). The engine's own names for the two mandatory records are in its error strings.

| record | absent → | why |
|---|---|---|
| **type1** *(“Tiles”)* | **reject** — loader status 5, `"Tiles block not found"` | the world builder `FUN_00548550` reads `map+0x0c` as u16/cell with no null test |
| **type2** *(“Altitudes”)* | **reject** — loader status 6, `"Altitudes block not found"` | same, `map+0x14` as u8/cell |
| type3 | **default**: a `W·H` plane of zeroes | same consumer reads `map+0x10` as i8/cell; nonzero = “an object blocks this cell”, so all-zero is a valid empty plane |
| type5 | **default**: one group record — `new(0x48)`, `FUN_0048d9d0(1,1)`, then `CPlayer+0x0c = 5000` (`ALM-REQ-056`), which is the field a record's own `+0x08` fills (`ALM-SCALAR-087`). Every shipped map authors `0` in that position on record 0 (`ALM-SCALAR-089`) | the group list must be non-empty |
| type0 | **default**: `W = H = 16`, empty name and text, `+0x70 = +0x74 = 1` | but see the ordering rules — this default is only safe if types 1 and 4..9 are absent too |
| type4, 6, 7, 8, 9 | **skip** | their loops are bounded by type-0 counts (type9's by its own first word), which are 0 when the records are absent |
| `typeId >= 10` | **skip** by `Seek(payloadSize)`, no error | the switch's `default` arm |

Header gates, in order: magic `M7R␀`; `recordCount >= 3`; `formatVersion <= 1001`. A duplicate
`typeId` is **not** rejected — the case runs again and overwrites the pointer.

**Ordering.** The permutation is not enforced, but three precedence relations are, because two
lengths and five loop bounds are carried in state that only one earlier case writes:

1. **type0 before type1** — the tile grid is allocated `W·H·2` from the *metadata's* `W`/`H`
   (default `16,16`), not from the record's own `payloadSize`.
2. **type2 before type3** — type3's allocation and read length is the **type-2 record's**
   `payloadSize`.
3. **type0 before types 4..9** — those cases' loop bounds are metadata fields.

**Four readers ship, with three different acceptance tests** (`ALM-RDR-059`). Beside the loader,
`rom.exe` has the `*.alm` browser scan's accept test — magic, `recordCount >= 2`, a type-0 record,
and type-0 payload **`+0x70 >= 2`**, so a map below that threshold is silently unlisted
(`ALM-META-058`) — and a light parser that handles typeIds 0..3 only, walks to EOF rather than by
count, and treats its two header failures as advisory message boxes without stopping. `Map
Editor.exe` carries that light parser's two literals and none of the loader's six; its own reader is
unread.

## type-0 metadata payload (632 bytes) — `ALM-META-008…ALM-META-010`

A fixed record: a scalar header, a 64-byte map name, two u32, a 64-byte description,
then a 440-byte trailing region of default-empty text slots. Offsets are constant
across all 38 maps; only `typeId` is a corpus-constant *value* (`ALM-META-008`).

Offsets are **`ALM-FRAME-031`-corrected** (each is 8 lower than the superseded framing; the
byte positions in the file are the same). They are re-derived independently by the
loader's own read sequence, whose lengths sum to exactly 632.

| Off | Type | Field | Notes | Claim |
|-----|------|-------|-------|-------|
| +0x00 | u32 | **W** | map width — read into the map object | ALM-HDR-001, ALM-META-024 |
| +0x04 | u32 | **H** | map height (`W≠H` occurs, e.g. 112×144) — read into the map object | ALM-HDR-001, ALM-META-024 |
| +0x08 | f32 | **angle** | radians. Read into the map object at `M+0x18` and into the landscape object at `P+0x20`, the **light/sun angle** slot (editor "Show Time / Light"). Every shipped value is a whole number of degrees: `±45 ±44 ±18 36 28 25 22`, modal `0x3f490fda`, one ULP below `float32(π/4)`. **Neither copy is read**: `M+0x18` is never loaded, and the landscape copy is overwritten from the sun globals before its first read | ALM-META-027, ALM-META-091, ALM-META-092, ALM-CORP-093 |
| +0x0c | u32 | scalar (stored) | five values only — `360`×62, `480`×2, `645`×2, `720`×4, `1080`×2 over 72 maps. Read into the map object at `M+0x1c` and into the landscape object at `P+0x2c`. **Neither copy is read anywhere**; role undetermined. `360` is also where the sun model's minute clock starts, which is an agreement between value sets, not a consumer | ALM-META-026, ALM-META-091, ALM-META-092, ALM-CORP-093 |
| +0x10 | u32 | scalar (stored) | `0..33`, mode `21`. Read into the map object at `M+0x20` and into the landscape ambient byte `P+0x1c`. **Neither copy is read**: `M+0x20` is never loaded, and the landscape byte is overwritten from a sun global before its only read | ALM-META-026, ALM-META-091, ALM-META-092, ALM-CORP-093 |
| +0x14 | u32 | scalar (stored) | `27..64`, `64`×37 and `63`×22. Read into the map object at `M+0x24` and into the landscape range byte `P+0x1d`. **Neither copy is read**, and the landscape read is four bytes wide into a one-byte slot, so its upper three bytes land on `P+0x1e/0x1f/0x20` | ALM-META-026, ALM-META-091, ALM-META-092, ALM-CORP-093 |
| +0x18 | u32 | bitmask | **terrain tile-group mask.** Read into a local and discarded by the map loader, but the landscape builder stores it at `P+0x28` and the tile loader consumes it: bit `i` selects tile group `(i>>2)+1`, variants `(i&3)*4 .. +3`. Eight values over 72 maps, all `≤0x1fff` — groups 1–3 entire plus group 4 variants 0–3, which is the shipped file set | ALM-META-026, TERR-LOAD-152 |
| +0x1c | u32 | **#players** | player-record count = `type5_size / 76` (38/38); `∈{3..9}` (editor caps at 16) | ALM-META-025 |
| +0x20 | u32 | **#objects** | object-record count (base records; extensions are additional bytes, not records) | ALM-META-025 |
| +0x24 | u32 | **#units** | unit-record count = `type6_size / 70` (38/38) | ALM-META-025 |
| +0x28 | u32 | (unused) | An earlier reading treated this as #type7, but case 7 reads its **own** count from the stream and ignores this slot (differs on 32/38 maps); read-and-discarded, meaning Unknown | ALM-META-025 |
| +0x2c | u32 | **#type8 records** | count used as the case-8 loop bound; `0 ⟺ type8 empty` | ALM-META-025 |
| +0x30 | char[64] | **name** | ASCII, NUL-terminated; ≤21 B used; empty on most campaign maps | ALM-META-010 |
| +0x70 | u32 | scalar (stored) -> `map+0xd4` | **the multiplayer-mode source**: `[0x005cd758]+0xc = (map+0xd4 > 1)` (`004e1bc1`/`004e1bc8`/`004e1bd1`), read by the player build at `004e238f`. All 28 `scenario.res` maps carry `1` on all three roots. The converse fails: the RU root's loose `Horror.alm` carries `1` while the EN file of that name carries `16`, and `ALM-CORP-060` shows the RU file is the EN file's first 262 876 bytes with exactly two bytes changed — `recordCount` 10 → 4 and offset `0x98`, which is this field, 16 → 1. The value identifies a single-player load, not a map stored in `scenario.res`. EN/install histogram `1:28 4:2 8:2 12:1 16:5` over 38 maps; RU `1:29 4:1 8:2 16:2` over 34. A player-slot / MP-capacity reading remains consistent with the values | ALM-META-026, ALM-MODE-070 |
| +0x74 | u32 | scalar (stored) | standalone `1..5`, campaign `1` | ALM-META-026 |
| +0x78 | char[64] | **description** | code-page text (ASCII or Windows-1251), NUL-terminated; ≤36 B used | ALM-META-010 |
| +0xb8 | 7×64 B | **text slots** | fixed array of 7 slots `[3×u32 prefix][char[52] text @+12]`, default `"<None>"`. Loaded (part of a 512-B block read from +0x78). Campaign maps fill slots 4/6 with trigger/quest text (`"mission complit"`, `"Start1"`, …) | ALM-META-028 |

The record's own `typeId = 0` and the per-map `selectorA` float are **header** fields
(`+0x0C`/`+0x10` of the record header), not payload — see the box at the top.

The whole type-0 record is read by `rom.exe`'s `.alm` loader (`FUN_00512369`, case 0) in
16 reads totalling exactly **632** bytes: twelve `u32` (`W/H`, the light `angle`, three
stored scalars, the discarded bitmask, the five content-record counts), `0x40` (`name`),
two `u32`, and one `0x200` block covering `description` + the 7 text slots. The `π/4`
constant is `+0x08` (`0x3F490FDA`).

## Grid layers type1 / type2 / type3 — `ALM-GRID-012`, `ALM-GRID-013`,
`ALM-GRID-014`, `ALM-GRID-032`

Each grid record stores **`W·H` cells starting at payload+0** — the payload being the
bytes after the record's 20-byte header — with `payloadSize` exactly `2·W·H` (type1) /
`W·H` (type2, type3). The payload is **pure grid**: nothing is overlaid on it, no cell is
lost, and the record's `typeId`/`f32` live in the header (`ALM-GRID-032`). Cell layout is
row-major, `W` cells per row:

```
cell (col, row)  ->  element index  row*W + col      (0 <= col < W, 0 <= row < H)
type1: u16 at payload + 2*(row*W + col)
type2: u8  at payload +    row*W + col
type3: u8  at payload +    row*W + col
```

This is exactly how `rom.exe` addresses them: the loader does `malloc(W·H·2)` +
`Read(ptr, W·H·2)` with an unadjusted destination, and the sim ingest `FUN_00548550`
reads `tiles[row·W+col]`, `heights[row·W+col]`, `type3[row·W+col]`. The size accounting
`dataSize = 20 + 4·20 (the four grid+meta headers…) …` is unchanged in value:
`4·W·H + 72` holds 38/38.

> **If you implemented the superseded `ALM-GRID-011` spec:** it placed the grid 8 bytes earlier, so
> your type1 render is **4 cells left** of the engine's and your type2/type3 layers **8
> cells left** — and therefore also 4 cells out of register with each other.

| Layer | Name | Cell | Encoding (`ALM-GRID-012…ALM-GRID-014`: `rom.exe` tables + corpus) | Claim |
|-------|------|------|------------------------------------------------|-------|
| type1 | **Tiles** | u16 LE | **tile-index word**: **bits 0–9 = tile index**, **bit 13 (`0x2000`) = impassable flag**, bits 10–12/14–15 unused (0.000% over 880 552 cells). Terrain class is *derived* from the index by the loader (below), not a raw high byte | ALM-GRID-012 |
| type2 | **Altitudes** | u8 | **height / altitude**. Named "Altitudes" in `rom.exe`; copied to the sim height buffer; height correlates with terrain class (Mountain highest, Water lowest) → a height field, **not** illumination. **⚠ 2026-07-25 sweep: the range "0..~246" is refuted and the per-class means (55.1 / 35.2) are stale** — both were measured on the superseded `ALM-GRID-011` base; every map's old-base maximum was a record-header byte (`TERR-LIGHT-016`), and `TERR-LIGHT-028`(d) finds no shipped height byte ≥ `0x80` | ALM-GRID-013 |
| type3 | **Objects** | u8 | **static-object placement layer** (0 = empty): code `c ≠ 0` resolves to `objects.reg` **section index `c − 1`** — see below. Optional block, zero-filled if absent; the sim ingest additionally marks every nonzero cell object-occupied (runtime block value **5**), which is a *derived* effect, not what the layer is. Census re-run on the corrected base, 38 maps: **86.19–97.94 % zero, 71 099 nonzero cells, 98 distinct codes 1..246** — the old-base figures reproduce exactly, so the earlier "stale base / the `246` is an injected header byte" warning is withdrawn | ALM-GRID-014, ALM-CLS-035 |

### type1 tile word — terrain resolution (`rom.exe`)

`terrainType, passability = f(tileIndex = word & 0x3ff)` (`FUN_00548720`):

- **strip group = bits 6–9** indexes a hardcoded (primary, secondary) terrain-type pair
  table at `world+0x54156`; **blend variant = bits 0–5** selects primary vs secondary
  and a **5-level** blend of their passability (auto-tiling terrain transitions).
- tile-index range **`[512,768)`** (bits 8–9 = `0b10`) = **Water** (special case).
- **bit 13** forces the cell impassable regardless of tile.

**Terrain enum** (1-based in the lookup; = `world.res:data/map.reg` `Terrain` record
order): `1 Land · 2 Grass · 3 Flowers · 4 Sand · 5 Cracked · 6 Stones · 7 Savanna ·
8 Mountain · 9 Water · 10 Road` (`ALM-TERR-015`). `map.reg` supplies a per-terrain `Cost`
and `Pass` scalar; `rom.exe` reads **only the ten `Cost*` keys**, into the byte table
`world+0x54176+class` (slot 0 = `0xff`), with its own hardcoded defaults
`[8,8,9,14,6,12,11,16,8,6]`.

**Restored by `ALM-TERR-043`** on the corrected `.reg` framing:

| class | 1 Land | 2 Grass | 3 Flowers | 4 Sand | 5 Cracked | 6 Stones | 7 Savanna | 8 Mountain | 9 Water | 10 Road |
|---|---|---|---|---|---|---|---|---|---|---|
| `Cost` | 8 | 8 | 8 | 14 | 6 | 12 | 8 | 16 | 8 | 6 |
| `Pass` | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 1 | 1 | 1 |

The previously published `Cost = [0,0,0,0,0,0,0,1,1,1]` / `Pass = [8,8,14,6,12,8,16,8,6,7]`
were read on the superseded `REG-FMT-017` record framing and are **regenerated exactly** by pairing
name(record i) with value(record i+1) — including the trailing `7`, which is `ScanShift`.

**`Pass` is inert.** No `Pass*` key string exists anywhere in `rom.exe` or in
`Map Editor.exe`; impassability is hardcoded (terrain class 8, plus the raw water bit
test). `Cost` is a **movement cost**, consumed only by the pathfinder's per-step add and
by the move-duration divide (`TERR-COST-052`). `map.reg` also carries a `Path Finding`
section — `SpeedMultiplier 8`, `StaticScanAhead 5`, `DynamicScanAhead 3`,
`StaticRefreshRate 16`, `DynamicRefreshRate 32`, `DynamicByStaticLookup 3`,
`StaticIsntNeeded 5` — and `Scanning`/`ScanShift 7`. The `.reg` name field holds 15
characters and `rom.exe` asks for keys up to 21; both lookup paths truncate the request
to 15 first, so every key still resolves.

**Terrain graphic** (which picture a tile-word draws) is a separate render mapping,
specified in `formats/terrain` (`TERR-IDX-003`, `TERR-SEM-004`): `g = (w & 0x1fff) >> 6` selects a
`terrain.3d/tileG-VV.bmp` strip (`G=(g>>2)+1`, `V=(g&3)*4+((w>>4)&3)`) and `w & 0xf`
selects the 32×32 sub-cell. `tile1/2` = land strips, `tile3` = animated water, `tile4` =
road; bit-13 tiles composite over `dirt.bmp`.

**Runtime passability** — corrected and completed by `TERR-PASS-049…TERR-PASS-051`.
The map load derives **three** 256×256 byte planes at fixed
stride 256, addressed `(row<<8)|col` regardless of `W`,`H`:

| plane | at | initial fill | built from |
|---|---|---|---|
| movement cost | `sim+0x00000` | `0x01`, overwritten for every in-bounds cell | the blended terrain cost |
| block bits | `sim+0x10000`, copied to `sim+0x20000` | `0` | type1 + type3 + the border |
| height | `sim+0x9451c` | `0` | type2 (Altitudes) |

The block byte is a **bitmask**, not an enum: `1 = bit0`, `5 = bit0|bit2`,
`0x1f = bits 0..4`; bit 5 marks a cell carrying a runtime record, bits 6/7 a ground/air
occupant on the dynamic plane. A cell blocks a mover iff `block[cell] & mover.mask != 0`
over the mover's `n×n` footprint, mask `0x41` ground / `0x44` / `0x82` air. Bit 1 is set
by nothing but the border, so **only the border stops an air mover**. The block arms, in
the order the ingest applies them:

```
if (w & 0x2000)                block = 1     tile-word bit 13
if (classify(w & 0x3ff) == 8)  block = 1     Mountain
if ((w & 0x300) == 0x200)      block = 1     water range, tested on the raw word
if (type3[cell] != 0)          block = 5     static object  (assignment: this one wins)
8-cell border                  block = 0x1f
```

Both block planes are **save state** (`TERR-PASS-053`); the cost and height planes are
not. So is the cell-record map the next paragraph names.

**That is the whole of what the ingest writes — it is not the whole of what blocks.** A placed
**type-4 structure** never goes through it: it attaches, in its own constructor and strictly after
the ingest, to a per-cell record, and one routine then recomputes that cell's cost byte and *both*
block bytes from the record. Its footprint can **clear** bits 0 and 2 as well as set them, which is
how a bridge crosses water. Specified in `formats/terrain` → "Structures on the block plane"
(`TERR-STRUCT-068`…`072`, `TERR-PASS-073`); the `kind`→class resolution is `ALM-CLS-036` above.

## Content sections type4–9 — `ALM-CNT-017`, `ALM-UNIT-018`, `ALM-OBJ-019`,
`ALM-GRP-020`, `ALM-TRIG-021`, `ALM-TRIG-022`, `ALM-CODE-023`; offsets corrected by
`ALM-FRAME-031`

Like the grids, the content payloads are pure: the first record starts at payload+0.
`X`/`Y` fields are `u32` fixed-point `/256` (integer part = tile, low byte usually
`0x80` = tile centre).

**Placement anchor → terrain cell is a bare `>>8`** (`ALM-PLACE-033`): no origin is
subtracted, no border inset added, no rounding applied. The engine shifts the object
anchors in the loader itself (`SAR ,8`) and the unit anchors when it walks the loaded
unit list. Combined with the grid base above, a placement at `(X,Y)` stands on tile
`(X>>8, Y>>8)` of every grid layer.

### type4 — placed structures/buildings (`ALM-OBJ-019`, `ALM-OBJ-034`, `ALM-CLS-036/037`)

`count = meta+0x20` — and that is the loader's own case-4 loop bound (`CMP ECX,[EBP-0x10]`
at `005128fa`), not merely a corpus identity. Base record = **20 bytes**; a record whose
`kind` is `0x21` appends **8 more bytes**. With that rule the walk consumes every payload
exactly, **38/38**.

| Off | Type | Field | Notes |
|-----|------|-------|-------|
| +0x00 | u32 | **X** | `/256`, `0x80`-centred. **This is the axis the tile grid strides by 1** and the block plane's **low** byte — the same axis as `meta+0x00` and as the Buildings column `sizeX` (`ALM-OBJ-061`) |
| +0x04 | u32 | **Y** | `/256`. The other axis: the grid's row index and the plane key's high byte |
| +0x08 | u32 | **kind** | **the class**: a `structures.reg` `ID` (= section index + 1). Corpus domain `1..66`, 65 distinct, never 0. `0x21` ⇒ an 8-byte extension follows; `0x22`/`0x23` ⇒ the `Shop` branch (below) |
| +0x0c | u16 | **stock value cap** | the `Shop` branch stores `value × 1000` at the live object's `+0x70` and forwards it to the stock generator's `template+0x90` (`SHOP-CAP-004`). **Not a durability** — that label is retracted; 23 distinct values over the 50 shipped records |
| +0x0e | u32 | **owner** | low u16 used: a **1-based physical slot** in the type-5 array, 3141/3141 |
| +0x12 | u16 | **trigger target id** | sign-extended into the live object's `+0x10`; it is the id a type-7 `Target_Structure` parameter names, 16/16 (`ALM-TRIG-046`) |
| (+0x14 | u32 | *ext, low axis* | present iff `kind == 0x21`. **Low byte only** → `obj+0x60`, the extent along the `+0x00` axis (`ALM-OBJ-062`). Bytes `+0x16`/`+0x17` reach no instruction) |
| (+0x18 | u32 | *ext, high axis* | **low byte only** → `obj+0x61`, the extent along the `+0x04` axis; overrides the definition's own footprint — `structures.reg`'s `VariableSize`. Bytes `+0x1a`/`+0x1b` reach no instruction) |

The override is selected by `(w & 0xff) + (h & 0xff) > 0` at `00504509`…`0050450f`, **not**
by the kind (`TERR-STRUCT-090`): an extension record with two zero bytes would take the
table arm. All 8 shipped extension records are class 33 and carry two non-zero bytes.

**Which C++ class a record becomes** (`ALM-CLS-037`, read from `FUN_004e2462` at
`004e24f0`…`004e2503`; the names are the MFC runtime-class table's own, `SHOP-CLS-001`):

```
kind == 0x22 or 0x23  ->  Shop,     0x74 bytes (FUN_00505cbd, vtable PTR_FUN_0059c828)
anything else         ->  Building, 0x6c bytes (FUN_005042b6, vtable PTR_FUN_0059c738)
```

`0x22`/`0x23` are `structures.reg` `ID` 34/35 = **`Shop 1` / `Shop 2`**, the two `Usable`
shops; 50 of 3141 shipped records take that branch. Note that `kind == 0x21` — the
*extension* discriminator — is an **object**: the two discriminators are different
questions with adjacent immediates. **`Shop` derives from `Building` and its constructor
opens by calling `Building`'s at `00505cea`**, so a shop resolves and attaches a footprint
like any other placement — with the two extent bytes pushed as literal `0`, i.e. always
from the table (`ALM-CLS-063`). Over the shipped maps that is 450 attached cells, 400 of
them blocking and 50 doorways.

**How `kind` reaches a class.** The engine does not subscript `structures.reg` with it
directly: `FUN_0050445d` uses the byte as a **1-based** index into a placeable-definition
table (`0x609be0`, guarded `kind != 0 && kind <= count − 1`), and that definition supplies
the footprint (`sizeX × sizeY` in tiles + the two per-cell masks), the HP pair and one
further byte. That table is the **Buildings collection of `world.res:data/data.bin`**
(`DAT-LOC-001`, `DAT-BLD-005`); the 1-based law is the file's own entry-0 skip.

### type5 — player/group roster (`ALM-GRP-020`, `ALM-GRP-041`)

`count = meta+0x1c`. **76-byte** fixed records, 38/38 exact. Every field below is a read
boundary in the loader's case-5 sequence (`4+4+4+32+16×2 = 76`).

| Off | Type | Field | Notes |
|-----|------|-------|-------|
| +0x00 | u32 | **colour slot** | → `player+0x08`, carried `+1` to `Player+0x44`; the unit body draw indexes the 17 shade objects with it (`ALM-PLAYER-069`, `PAL-SHADE-012`, `PAL-SHADE-013`). Sparse and unordered, e.g. `Cross.ALM` `8 2 3 5 1 6 4 13`. The owner lookup `FUN_004fb534` compares the record's own 1-based ordinal (below), **not** this word |
| +0x04 | u32 | **human-participant flag** | `0`/`1` → `Player+0x28`, the word `UNIT-OWNER-009` reads as *a human participant owns this* (`ALM-GRP-041` as amended, `ALM-PLAYER-069`). Every campaign map authors `0` on record 0; every loose-map record and every record above ordinal 0 authors `1` |
| +0x08 | u32 | **scalar; effect Unknown** | `0` or `5000` and nothing else, 357/357 records over both roots; all 71 record-0 entries author `0`, and 26 of the 286 later records do too (`ALM-SCALAR-089`). Loaded into `CPlayer+0x0c`. The seven-function family has one direct read, in a copy constructor called twice by a direct-reference-free owner (`ALM-SCALAR-087`). A selected population of **17 map-side plus 7 family routines** has no map-side direct-displacement read attributed to this field, but computed pointers and whole-object mover/serializer/`memcpy` paths remain unsearched (`ALM-SCALAR-088`). Do not treat it as inert. Not redundant with `+0x04`: 41 records pair `+0x04 = 1` with `+0x08 = 0` |
| +0x0c | char[32] | **name** | NUL-terminated ASCII: `Self, Monsters, Villagers, Neutral, Enemy, Beasts, Guards, Peasants, Orcs, Trolls, …` |
| +0x2c | u16[16] | **relations** | one per editor player slot (the editor caps at 16) — the diplomacy row |

The object each record is loaded into is the interface `CPlayer`, RTTI object size `0x48`, vtable `0x0059a560` — **six dwords: five function entries and a null at `0x59a574`**; `0x59a578` begins the next class (`ALM-CPLAYER-090`). Its fields take the record's reads in order: `+0x08` the colour word, `+0x30` the human-participant flag, `+0x0c` the scalar above, `+0x10` the 32-byte name, `+0x34` the diplomacy `CWordArray`, and `+0x04` the loader-written ordinal. Both non-copy constructors leave `+0x0c` at `0` (`UNIT-VPLAYER-022`); the map loader writes the authored scalar or the absent-section default.

A record's **physical slot + 1** is what type-4 `+0x0e` and type-6 `+0x14` store
(`ALM-OWN-039`); the loader itself writes `slot+1` to `player+0x04`. Do **not** use the
`+0x00` colour word for that: it fails on 22.7 % / 27.7 % of shipped records.

### type6 — placed units (`ALM-UNIT-018`)

`count = meta+0x24`. **70-byte** fixed records, 38/38 exact; **all 8094 records'
`(X,Y)` land inside the map** (a wrong stride would scatter them OOB). The loader reads
the 70 bytes in **28** stream reads, of which **18** are displaced and three pass
through a stack local, so a file offset is usually not the runtime offset. From file
`+0x2c` to `+0x40` the displacement is a uniform `−4`: the skill run at file `+0x35`
is `rec+0x31` and the resistance run at file `+0x3b` is `rec+0x37` (`ALM-TAILMAP-079`).
Offsets in the table below are file offsets. The record's only sentinel is `0xFFFF`,
and it belongs to the four `u16` at `+0x20`…`+0x27`
(`ALM-TAILU16-082`); the earlier reading that put `0xFF` inventory sentinels in the two
byte runs at `+0x35` and `+0x3b` is refuted — no byte of `+0x2c`…`+0x3f` is `0xFF` on any
of the 12 085 records of the two preserved roots, and those runs are stat overrides
(`ALM-TAILRUN-083`).

| Off | Type | Field | Notes |
|-----|------|-------|-------|
| +0x00 | u32 | X | `/256`, `0x80`-centred |
| +0x04 | u32 | Y | `/256` |
| +0x08 | i16 | **class** | a `units.reg` **`ID`** (domain `1..80`; inside the 34-element `ID` set on 8094/8094, inside the 34 section indices on 17.9 %) |
| +0x0a | u16 | **class (2nd key)** | only consulted when `+0x08 ∈ {0x1a,0x1b}` or `≥ 0x40`; also the `npc.reg` index on the NPC path |
| +0x0c | u32 | flags | bit 0 ⇒ resolve as an NPC; bit 2 sets the live unit's `+0x4b` high bit |
| +0x10 | u32 | **definition id** | present only when `formatVersion > 0x3da`; when nonzero and `!= 0xcdcdcdcd` it **overrides** `+0x08`/`+0x0a` (matched on the definition's parameter `0x18`). **But it is read only when `+0x0c` bit 0 is clear** — `FUN_004e26bb` tests the NPC flag *outside* the definition-id test, so a record carrying both takes the npc arm and this field is never read (`MISSION-ARM-006`) |
| +0x14 | u32 | **owner** | low u16 used: 1-based physical slot in the type-5 array, 8094/8094 |
| +0x18 | u32 | **type-8 link** | 1-based, bounded by `meta+0x2c`; `0` = none. The loader writes this record's `+0x40` word into entry `[value−1]` of `mapObj+0x2dc`. Guard satisfied 8094/8094; 43 records use it |
| +0x1c | u32 | — | domain exactly `0..14` over 12 085 records, 6967 even against 1127 odd on EN and 3857 against 134 on RU; **no reader in `rom.exe`** — role Unknown, and the surviving explanation is the editor's own loader (`ALM-TAILDIR-081`) |
| +0x20 | u16 | **current health** | absent value `0xFFFF`; applied to `actor+0x94` (`UNIT-PLACE-034`). Authored on 30 records per root, over seven `scenario.res` maps (`ALM-TAILU16-082`) |
| +0x22 | u16 | — | absent value `0xFFFF`; authored on exactly the records `+0x20` is and with a different value set, and **no reader** (`ALM-TAILU16-082`) |
| +0x24 | u16 | **current mana** | absent value `0xFFFF`; applied to `actor+0x9a` (`UNIT-PLACE-034`). `0xFFFF` on 12 085/12 085 records — the consumer exists and no shipped map exercises it (`UNIT-PLACEIDLE-088`) |
| +0x26 | u16 | — | `0xFFFF` on 12 085/12 085 records, and **no reader** (`ALM-TAILU16-082`) |
| +0x28 | u16 ×2 | — | present only when `formatVersion > 0x3b6`; they land at `rec+0x4c`/`rec+0x4e` and **nothing reads either** — `disp:4e` is empty image-wide. `+0x28` is `0` on every record; `+0x2a` carries four values on 348 records per root (`ALM-TAILVER-084`) |
| +0x2c | u8 ×4 | **stat overrides** | Body, Mind, Spirit, Reaction in that file order → `actor+0x84`/`+0x88`/`+0x8a`/`+0x86`, zero = absent (`UNIT-PLACE-034`, file order fixed by `ALM-TAILMAP-079`). Authored on 3 records per root |
| +0x30 | u8 | — | lands at `rec+0x2c`; **no reader**, and `0` on every record of both roots (`ALM-TAILRUN-083`) |
| +0x31 | u8 ×2 | — | → `actor+0xbe` and `actor+0xc0`; `+0x32` present only when `formatVersion > 0x3d8`. `+0x31` is `0` on every record and `+0x32` nonzero on one (`UNIT-PLACE-034`, `UNIT-PLACEIDLE-088`) |
| +0x33 | u8 ×2 | — | land at `rec+0x2f`/`rec+0x30`; **no reader**, and `0` on every record of both roots. `disp:30`'s two hits inside the loader are a vtable call and a dword store, both outside case 6 (`ALM-TAILRUN-083`) |
| +0x35 | u8 ×6 | **skill overrides** | → `actor+0xa8 + 2i`, zero = absent — but the loop's index starts at **1**, so only `+0x36`…`+0x3a` are applied and `+0x35`, the `Skill.General` slot, is stored and never read (`UNIT-PLACESKILL-086`) |
| +0x3b | u8 ×5 | **elemental resistances** | all five applied to `actor+0xc4 + 2i`, zero = absent, after the spawner's re-derivation call (`UNIT-PLACERESIST-087`) |
| +0x40 | u16 | **unit id** | what a type-7 `Target_Unit` parameter names (`ALM-TRIG-046`); 39 distinct over `scn:110`'s 39 records. Also copied into the linked type-8 entry's `+0x3c` |
| +0x42 | u32 | **group id** | what a type-7 `Target_Group` parameter names, 169/169 (`ALM-TRIG-046`); **not** unique per record — 19 distinct over those same 39. The map keeps its running maximum at `mapObj+0x90`, i.e. the next free **group** id (`ALM-UNIT-048`; `ALM-UNIT-040` had the two labels the wrong way round) |

The record is **70 bytes only at `formatVersion == 990`**: three of the reads above are
version-gated (`> 0x3b6`, `> 0x3d8`, `> 0x3da`), so a different version yields a different
stride (`ALM-UNIT-040`).

**Who ever reads a placed-unit record.** The loader stores the record's *pointer* into a
4-byte-element `CObArray` at `mapObj+0x318`, and exactly three routines can hold one: the
loader itself, the placement spawner `FUN_004e26bb`, and the map object's destructor,
which frees the record without reading a field. Between them they read `+0x00`, `+0x04`,
`+0x08`, `+0x0c`, `+0x10`, `+0x14`, `+0x20`, `+0x24`, `+0x28`…`+0x2b`, `+0x2d`, `+0x2e`,
`+0x31+i` (`i=1..5`), `+0x37+i` (`i=0..4`), `+0x3c`, `+0x40`, `+0x44` and `+0x48` of the
runtime record — nothing else (`ALM-TAILHOLD-080`). That closure is what makes the "no
reader" entries above decidable, because a displacement sweep at a small offset finds
every structure in the image with a field there.

**A consumer can skip the whole tail and still be right about 7919 of the EN root's 8094
placements, and 3839 of the RU root's 3991** (`UNIT-PLACEIDLE-088`). Two of the tail's
consumers — current mana and `actor+0xbe` — are never exercised by any shipped map, and
the busiest field in it, the sixth skill byte at `+0x3a`, moves 119 records.

**Which shipped file a placement reads** (`MISSION-ARM-006`, `FUN_004e26bb`). The outer
discriminator is the class key, not a flag, and the four arms are ordered:

```
+0x08 >= 0x1a                          -> data.bin Units,  matched on params 0x1d/0x1e
+0x08 <  0x1a and +0x0c bit 0          -> npc.reg [npc<+0x0a>] -> its DataBinID
                                          -> data.bin Humans on param 0x18 (serverID)
+0x08 <  0x1a, bit0 clear, +0x10 != 0  -> data.bin Humans on param 0x18, +0x10 being the id
+0x08 <  0x1a, bit0 clear, +0x10 == 0  -> data.bin Humans on param 0x10 (typeID)
```

Corpus, 8094 records / 38 maps: **6672 / 15 / 1405 / 2** in that order. The definition lookup
`FUN_004de63e` searches **backwards**, never tests index 0, and **returns 0 on a miss** — a
valid index, not an error (`MISSION-DEF-007`).

Those four arms also select the Humans constructor mode. Both direct Humans arms pass zero. The
npc arm queries the literal flag `Hero` on its `npc.reg` section and passes that Boolean. The
constructor always streams the Humans row's slot-16 `typeID` to `actor+0x0e`, then replaces it with
the player-character `gender+0x21/+0x23` value only when this mode is non-zero. Consequently a
Humans placement is not by itself a persistent party actor: mission 20's npc without `Hero` and its
three definition-id Humans retain table types outside the mission-end keep band
(`PARTY-M20-030`, `PARTY-M20-031`).

### type7 — the trigger script: three counted arrays (`ALM-TRIG-044`…`047`)

The payload is **not** one list of nodes. It is three arrays, back to back, each preceded
by its own `u32` count, and the walk consumes every payload exactly on **38/38** maps:

```
[u32 nAct ][ nAct  x 796-byte node    ]   the "THEN" vocabulary, Description Instants.ini
[u32 nCond][ nCond x 796-byte node    ]   the "IF"  vocabulary, Description Checks.ini
[u32 nTrg ][ nTrg  x 184-byte trigger ]   binds conditions to actions
```

`ALM-TRIG-021`'s `entryCount` is `nAct`; the previously reported "8-byte trailer" is
`[u32 nCond = 0][u32 nTrg = 0]`, which is what a skirmish map with no script looks like.
Corpus totals: 797 actions, 680 conditions, 421 triggers.

**Node — 796 bytes, identical for an action and a condition:**

| Off | Type | Field | Notes |
|-----|------|-------|-------|
| +0x000 | `char[64]` | editor label | the map author's own name; NUL-terminated inside a fixed buffer whose tail is uninitialised editor memory |
| +0x040 | `u32` | **opcode** | the `ID` declared in the `.ini` (`EDITOR-023`) |
| +0x044 | `u32` | **id** | unique within its own list, 38/38 lists of each kind; what a trigger references |
| +0x048 | `u32` | — | `0` on all 1477 shipped nodes |
| +0x04c | `u32[10]` | `value[Par0..Par9]` | **indexed by slot, not packed** |
| +0x074 | `u32[10]` | `type[Par0..Par9]` | `0` = slot unused |
| +0x09c | `char[64][10]` | `pname[Par0..Par9]` | the `.ini`'s `Par<i>_NAME`; `"<None>"` when unused |

Parameter **type codes** (`ALM-TRIG-046`), and what a value of that type names:

| Code | `.ini` type | The value is |
|---|---|---|
| 1 | `int`, `Enum` | a literal |
| 2 | `Target_Group` | the type-6 record's `+0x42` group id |
| 3 | `Target_Player` | 1..8 |
| 4 | `Target_Unit` | the type-6 record's `+0x40` unit id |
| 5 / 6 | `X` / `Y` | a plain tile coordinate — 505/505 inside `[0,W)×[0,H)` |
| 7 | `Const` | the literal the `.ini` declares for this variant; it is what selects among the eleven signatures sharing opcode 6 |
| 8 | `Target_Item` | domain Unknown (values 2..36) |
| 9 | `Target_Structure` | the type-4 record's `+0x12` word |

`Target_Building` is declared by the editor and placed by no shipped map.

**Trigger — 184 bytes:**

| Off | Type | Field |
|-----|------|-------|
| +0x00 | `char[64]` | name |
| +0x40 | 64 B | never a printable string (0 of 421); holds editor heap addresses |
| +0x80 | `u32[3][2]` | three `(left, right)` **condition ids** — pairs, never half-empty (1263/1263) |
| +0x98 | `u32[4]` | up to four **action ids** |
| +0xa8 | `u32[3]` | one **comparison code** per pair — `0 ==`, `1 !=`, `2 >`, `3 <`, `4 >=`, `5 <=`, dispatched through a 6-entry table; the three pairs are **ANDed with short-circuit** and a code above 5 is permanently false (`TRIG-CMP-006`) |
| +0xb4 | `u32` | the **once flag**: `1` = fire at most once per session (the byte latch at `session+0xbec4+index` gates re-entry), `0` = re-evaluate and re-fire every full tick; corpus 387 one-shot / 34 repeating, the 34 including every standalone map's single trigger (`TRIG-FIRE-007`) |

References are node **ids**, not array indices: 735/735 action and 1092/1092 condition
references resolve by id, against 85 % / 92 % by index.

### type8 — authored loot, and type9 — caster payload (`ALM-SACK-065`, `ALM-TRIG-049`)

Neither is a named-node tree, and **type8 has no count word at all**. Both walks consume
every payload exactly, 38/38.

> **What they are NOT (`retracted.md` → `ALM-TRIG-022`).** These are **not** the
> geometry behind the trigger machinery's box and circle checks — those checks read their own
> parameter slots and never touch a type-8 record. The section table above carried
> "condition/marker regions (box/circle X/Y) — serialized tree" until
> 2026-08-01, contradicting this heading four hundred lines further down its own file.
>
> **And type-8 is not caster payload either (`ALM-SACK-065`, `ALM-SACK-066`,
> `ALM-LIM-067`, `ALM-ORD-068`).** `FUN_004e4f3e` is the only consumer of
> either section, but it has **two** loops over two different lists. The first walks `mapObj+0x2f0`,
> which case 9 fills, and that is where the `spellbook` / `building - caster` strings are. The second
> (`004e5483`..`004e59cf`) walks `mapObj+0x2dc`, which case 8 fills, and it builds an item container
> per record and calls the sack maker `FUN_0050f5aa` at `004e59be`. **type8 is the map's authored
> loot**, and on a shipped single-player campaign map it is the only sack source that runs at load.

```
type9: [u32 count] then count records of 26 + 6*n bytes
       +0x00 u32 tag
       +0x04 u32 X · +0x08 u32 Y
       +0x0c u16 A · +0x0e u16 B · +0x10 u16 C
       +0x12 u32 spellRaw · +0x16 u32 n
       then n elements of 6 B = (u16 kind, u16 low, u16 high)

type8: no count word - the record count is meta+0x2c (the loader's case-8 bound).
       Records of 20 + 10*n bytes; the head is 20 B only at formatVersion >= 0x3dd,
       16 B below it (005133f8), and every shipped map is 0x3de.

       +0x00 u32 n      element count. A LOCAL: the arm never stores it (00513448)
       +0x04 u32 owner  -> obj+0x3c. 0 = a sack on the ground; non-zero = the id of an
                          actor that already exists, looked up in a map FUN_004e4f3e
                          builds from the registry 0x00609558 keyed by actor+0x8.
                          OVERWRITTEN at load by case 6 for any record a type-6
                          placement's +0x18 links to (00512f9e) -- 43 of 43 agree
       +0x08 u32 X      -> obj+0x40. 004e5949 SAR EAX,0x8 -> position byte 0
       +0x0c u32 Y      -> obj+0x44. 004e5969 SAR EAX,0x8 -> position byte 1
       +0x10 u32 gold   -> obj+0x48. FUN_0050f5aa arg 3; 0050f6ac adds it to sack+0x3c
                          (version-gated: 0 below 0x3dd)
       then n elements of 10 B:
       +0x00 u32 item   low u16 is a PACKED code -- FUN_004dcf92 cuts bits 8..11 (class),
                        12..15, 5..7, 0..4 and FUN_004dd02a allocates 0x84 Weapon (1),
                        0x68 (2), 0x68 (3..13), 0x50 Item (14), null otherwise. Class 14
                        takes the whole low BYTE as its index instead of bits 0..4
       +0x04 u16        read ONLY on the owner!=0 arm (004e58df): 0 -> append to the
                        actor's carrier at actor+0x7c, non-zero -> actor vt+0x3c
       +0x06 u32 link   1-based index into the type-9 list, 0 = none (004e5605)

       A record with owner==0 and a non-empty container becomes a Sack (new 0x44) at the
       cell; if a sack is already there FUN_0050f5aa merges the container and adds the gold,
       so two records on one cell become one sack.
```

Corpus, both roots: 30 of 38 / 28 of 34 maps carry a type8 record; 137 / 133 ground sacks,
43 / 43 owner records, 181 / 177 elements, 636 074 / 125 074 gold; 0 cells out of the map,
0 links out of range, 0 item codes resolving to null. 23 of 28 check-opcode-14 nodes name a
type8 ground cell in the direct axis order and 0 in the swapped one.

`ALM-EFFLINK-072` and `ALM-EFFPOP-073` read the link operation and type9 consumer end to end. The loader and consumer each
subtract one from a positive file link before indexing the type9 list; zero means none and file
value one reaches zero-based record zero. Case 8 writes the referring element ordinal and owning
type8 pointer to the target type9 object's runtime `+0x00/+0x04`; these are runtime fields, not
wire offsets.
The consumer accepts a linked recipe only when its file X/Y are both zero and the backlink exists.
It then appends Effects in this order:

1. non-zero A becomes kind `A + 43`, with B/C as its operands;
2. non-zero low word of `spellRaw` becomes kind 41, or kind 42 for a Book;
3. every tail element becomes `(kind, low, high)` in file order, with input kind 41 remapped to
   runtime kind 49.

This construction finishes before the owner-zero branch. A ground sack and actor stock therefore
receive the same enchanted item operation; the link is not a ground-only or owner-only field.
Across all 28 campaign maps per root there are 176 type8 records, 177 elements, 199 type9 records
and 74 positive links. Every link is in range, every linked type9 has X/Y zero and exactly one
reference, and the EN/RU linked rows produce the same Effect sequences (`ALM-EFFREC-071` through
`ALM-EFFPOP-073`).

### `rom.exe` cross-reference (`ALM-CODE-023`)

The separate type-9 cell-entry arm is `UNIT-M10CELL-054`. On version-990
records, X or Y nonzero and unsigned A below 4 select it. The target key is
`u8(X) | (u8(Y)<<8)`. Its six-byte payload is:

| Byte | File source | Role |
|---|---|---|
| 0 | low byte of spellRaw | spell/operation |
| 1 | byte 2 of spellRaw | power |
| 2,3 | low bytes of tail entry 0 words 0,1 | caster source x,y |
| 4,5 | low bytes of tail entry 1 words 0,1 | operation-26 relocation x,y |

This arm reads two tail entries without a length guard. Their third words do
not feed this payload. These coordinate roles do not replace the linked-item
Effect interpretation: X=Y=0 is that other arm's eligibility condition. A bit 2
selects a separate building-caster branch with a building-key lookup; the cell
arm performs no such lookup. Mission 10 has exactly two cell selections,
(22,64) and (21,63), both spell 13 / power 1, and no building-caster selection.
The remaining nine records have zero X/Y. Runtime admission and lifetime are
`UNIT-M10ENTRY-055` through `UNIT-M10LIFE-057`.

`%d.alm` is built + loaded in `FUN_00477c00`→`FUN_00572a2a` (`CMap` family).
`FUN_004d403c` reads groups (type5, `"…no groups"`) + the drop location (type7, `"…no
drop location in .alm…"`); `FUN_004f12d7` reads Outpost/repopper effects (type7). All
via a generic named-node accessor family (`FUN_004e17d2` find-by-name, `FUN_0051ab50`
get-param, `FUN_0051ac40`/`IsEmpty` iterate). This historical cross-reference
does not classify types 8/9 as named-node trees: their distinct record consumers
are `ALM-SACK-065`, `ALM-EFFREC-071` and `UNIT-M10CELL-054`.

## Trailer — there isn't one

The superseded framing saw an 8-byte trailer because it split every record 8 bytes early; those two
`u32` are the **last 8 bytes of the type-7 payload** (`ALM-TRL-005`/`ALM-TRL-030`, both
retracted by `ALM-FRAME-031`). The chain ends exactly at EOF.

## Reading algorithm

```
assert bytes[0:4] == "M7R\0"
version   = u32(bytes, 0x04)             # == 20
dataSize  = u32(bytes, 0x08)             # == 4*W*H + 72
count     = u32(bytes, 0x0C)             # the loader's loop bound; gated >= 3, NOT == 10
fmtVer    = u32(bytes, 0x10)             # == 990
off = 0x14
for k in 0..count-1:
    tag, hdrLen = u32(off), u32(off+4)    # == 7, 20
    size        = u32(off+8)
    typeId      = u32(off+12)             # <- the loader's switch; >= 10 is skipped by seek
    perMapConst = f32(off+16)
    payload     = bytes[off+20 : off+20+size]
    off        += 20 + size
assert off == len(bytes)                  # no trailer
# acceptance: require types 1 and 2 only. Default type3 to a W*H zero plane, type5 to one
# group record, type0's fields to W=H=16 / empty text / +0x70=+0x74=1. Skip the rest.
# Do NOT require 10 records, the full set {0..9}, or the physical order.

# grids (payload = pure data):
#   type1: u16 cell (col,row) at payload + 2*(row*W + col)
#   type2: u8  cell (col,row) at payload +   (row*W + col)
#   type3: u8  cell (col,row) at payload +   (row*W + col)
#
# content records (meta = type-0 payload):
#   type5: 76-B records, count = u32(meta,0x1c)            -> id @ rec+0, name @ rec+0x0c
#   type6: 70-B records, count = u32(meta,0x24)            -> X @ rec+0, Y @ rec+4 (/256),
#          class = units.reg ID @ rec+8, owner slot @ rec+0x14
#   type4: 20-B records, count = u32(meta,0x20)            -> X @ rec+0, Y @ rec+4,
#          kind = structures.reg ID @ rec+8, owner slot @ rec+0x0e;
#          kind == 0x21 -> 8 extra bytes follow (footprint override);
#          kind in {0x22,0x23} -> CStructure, else CObject
#
# class binding:
#   type3 cell code c != 0  ->  objects.reg    section index c - 1
#   type4 kind              ->  structures.reg ID (= section index + 1)
#   type6 rec+8             ->  units.reg      ID (sparse 1..80; NOT the section index)
#   owner (type4 +0x0e, type6 +0x14) -> 1-based physical slot in the type-5 array
#   type7: [u32 nAct][nAct x 796][u32 nCond][nCond x 796][u32 nTrg][nTrg x 184]
#          node    : name[64], opcode, id, 0, value[10], type[10], pname[10][64]
#          trigger : name[64], 64 B junk, 3x(condId,condId), 4x actionId, 3x opcode, flag
#          a trigger slot holds a node's id (node+0x44), NOT its array index
#   type9: [u32 count][count x (26 + 6n)], n = u32 at rec+0x16; X,Y = tile indices
#   type8: NO count word; meta+0x2c records of (20 + 10n), n = u32 at rec+0; the
#          head is 20 B at formatVersion >= 0x3dd and 16 B below it;
#          X,Y at rec+8/+0xc are 0x80-centred /256 anchors
#
# placement -> tile: cell = (X >> 8, Y >> 8). Nothing else: no origin, no inset.
```

## Invariants (hold across all 38 EN maps)

These are **corpus** invariants of authored maps — what a *writer* may be held to — not the reader's
acceptance test, which is *Acceptance contract* above. Two of them are broken by a shipped file
outside this corpus: `ALM-CORP-060` walked both preserved installs (72 files, 44 distinct names) and
`gameversions\ru\Horror.alm` carries **four** records, typeIds `{0,1,2,3}`, while satisfying every
other line below (`ALM-CORP-060`).

- Exactly **10 records** (**EN corpus only**; `recordCount` is a header field the loader gates only
  `>= 3`); `tag == 7` and `hdrLen == 20` in every record header (380/380 EN, 720/720 over both roots).
- Every payload is in-bounds: `off + 20 + size <= len`.
- **Tiling is exact:** `20 + Σ(20 + payloadSize) == fileSize` — no gaps/overlaps/trailer.
- `dataSize == 4·W·H + 72`, with `W,H` read from the type-0 payload.
- typeIds are `{0..9}` once each, in physical order `0,1,2,3,5,4,9,8,6,7` (**EN corpus only**; the
  loader enforces neither the set nor the permutation — only three precedence relations,
  `ALM-ORD-057`). “Once each” holds 72/72 over both roots; the set and the order hold 71/72.
- **Content counts (38/38):** `#type5 == meta+0x1c`, `#type4 == meta+0x20`,
  `#type6 == meta+0x24`; `76·#type5 == size5`, `70·#type6 == size6`; the type-4 walk with
  the `kind == 0x21` rule consumes each payload exactly (38/38).
- **Placement coordinates (38/38):** all **11 235** type-4 + type-6 anchors satisfy
  `0 ≤ X/256 < W`, `0 ≤ Y/256 < H`.
- **Class keys resolve (38/38):** every type-6 `+0x08` is in `units.reg`'s `ID` set
  (8094/8094); every type-4 `kind` is in `1..66` = `structures.reg`'s `ID` domain
  (3141/3141, never 0); every owner is in `[1, #type5]` (11 235/11 235); every type-6
  `+0x18` satisfies the loader's `meta+0x2c` guard (8094/8094); **71 035 of 71 099**
  nonzero type-3 codes resolve inside the 82-entry object array under `c − 1`.
- **Placements are legal on the corrected grids:** the derived-impassable rate
  (water | bit-13 | type-3 nonzero) over those 11 235 anchors is **4.9 %**, versus 20.7 %
  on the superseded `ALM-GRID-011` base.
- **The trigger script closes (38/38, `ALM-TRIG-044`/`049`/`050`):**
  the **physical record order is `0,1,2,3,5,4,9,8,6,7`** on 38/38, and the two departures
  from the typeId sequence are exactly the two dependencies the loader carries: case 6
  indexes the type-8 list and case 8 indexes the type-9 list, so 9 precedes 8 precedes 6
  (`ALM-ORD-068`);
  `4 + 796·nAct + 4 + 796·nCond + 4 + 184·nTrg == size7`; the type-9 walk
  (`26 + 6n` per record) and the type-8 walk (`meta+0x2c` records of `20 + 10n`) each
  consume their payload exactly. Each of the three is the **only** model in its parameter
  family that does.
- **The script's references resolve (`ALM-TRIG-045`/`046`/`047`):** 1476 of 1477 nodes
  match a `Description *.ini` group on all ten `(name, type)` slots; 735/735 action and
  1092/1092 condition references name an existing node id; 1263/1263 condition pairs are
  both-set or both-clear; 505/505 `X`/`Y` parameters are inside the grid; 169/169
  `Target_Group` values are a type-6 `+0x42`.

## Notes & open questions

- No ASCII chunk tags: records are typed numerically (`typeId` in the record header),
  unlike a RIFF-style container.
- `hdrLen = 20` is a **stored** field equal to the header length. It does **not** pin
  *where* the header starts — `[f0][f1][tag][hdrLen][size]` and
  `[tag][hdrLen][size][typeId][f32]` are both 20 bytes, which is how the 8-byte
  mis-split survived until `ALM-FRAME-031` corrected it.
- **No grid identifier overlay** (`ALM-FRAME-031`, `ALM-GRID-032`): the `[typeId][f32]` are record-header
  fields; the type1/2/3 payloads are pure `W·H` grids from payload+0. Anything reading
  the grid at the old base is 4 cells (type1) / 8 cells (type2/3) off in X.
- **Grid semantics resolved (`ALM-GRID-012…ALM-GRID-014`):** type1 = **Tiles** (tile-index word + bit-13
  impassable flag; terrain class derived via `rom.exe`'s strip→terrain-pair + blend
  tables), type2 = **Altitudes** (height, decided from the loader's own error name +
  Mountain-high/Water-low corpus stat), type3 = the **static-object placement
  layer** (`ALM-CLS-035`: code `c ≠ 0` → `objects.reg` section index `c − 1`;
  the sim ingest's block value 5 is derived from it). Terrain vocabulary =
  `world.res:data/map.reg` `Terrain` (10 classes; only `Cost` is read — `ALM-TERR-043`).
- **Content sections type4–9 decoded (`ALM-CNT-017`, `ALM-UNIT-018`, `ALM-OBJ-019`,
  `ALM-GRP-020`, `ALM-TRIG-021`, `ALM-TRIG-022`, `ALM-CODE-023`):** type4 objects (20 B), type5 groups
  (76 B), type6 units (70 B) are fixed-record tables whose counts are mirrored in the
  type-0 metadata; type7 is the trigger script — three counted arrays, actions then
  conditions then triggers (`ALM-TRIG-044`); type8 is the map's **authored loot**
  (`ALM-SACK-065`) and type9 the caster payload (`ALM-TRIG-049`). `type8` is empty on
  skirmish maps, which have no script.
- **Loader-confirmed (`ALM-META-024…ALM-META-028`, `ALM-FRAME-031`):** the `.alm` reader
  `FUN_00512369` reads a 20-byte **file** header (magic, 20, dataSize, `count`,
  `version`), then per record a 20-byte header and dispatches on its `typeId`; case 0
  reads W/H, the light angle, the stored scalars, and the five content-record **counts**
  (`+0x1c` players, `+0x20` objects, `+0x24` units, `+0x28`, `+0x2c`) — 632 bytes in
  total; the `+0x18` bitmask is read and discarded. `type8` is empty exactly when
  `+0x2c = 0`.
- **`dataSize = 4·W·H + 72` is no longer explained.** The superseded `ALM-HDR-001` reading treated the 72 as
  `12 (file hdr) + 3×20 (grid section hdrs)`; corrected that is `20 + 60 = 80`. The identity
  still holds 38/38, but the loader never reads the field, so it is writer-side only and its
  decomposition is Unknown (a legacy layout is plausible — the loader still honours
  `formatVersion == 1000`, which reads no record headers).
- **Class binding decoded (`ALM-OBJ-034`, `ALM-CLS-035…ALM-CLS-038`, `ALM-OWN-039`,
  `ALM-UNIT-040`, `ALM-GRP-041`, `ALM-CLS-042`):** which field names a class, and against which
  registry key, for all three placement kinds — see the type-3 note and the type4/5/6
  sections above (`ALM-CLS-035…038`, `ALM-OWN-039`, `ALM-UNIT-040`, `ALM-GRP-041`).
  **The one hop it could not close:** type-4 `kind` and type-6 `+0x08`/`+0x10` are
  resolved through a **placeable-definition database** (`0x609b18`, sub-collections
  `0x609bb8` for units and `0x609be0` for buildings) whose *populating file was not
  found*. Until it is, a consumer can name the class but not the definition's own
  parameters (footprint, HP, and whatever else the collection carries). Also open: the
  instruction that turns a definition into the live object's `+0x20`, the subscript the
  renderer and hit-test use on the class arrays.
- **Still undecoded:** the trigger's 64-byte junk field, the `Target_Item` domain, and
  the type-9 `+0x00` tag word (`ALM-TRIG-047`, `ALM-TRIG-049`, `ALM-TRIG-050`);
  the type-6 fields that **no routine able to hold a record pointer reads** — file `+0x1c` (domain
  `0..14`), `+0x22`, `+0x26`, `+0x28`, `+0x2a`, `+0x30`, `+0x33`, `+0x34` and the
  `Skill.General` byte at `+0x35` (`ALM-TAILDIR-081`, `ALM-TAILU16-082`,
  `ALM-TAILVER-084`, `UNIT-PLACESKILL-086`) — which are Unknown in role but no longer
  Unknown in whether they matter to this image;
  the exact gameplay meaning of the stored type-0
  scalars `+0x0c/+0x10/+0x14/+0x74` — for the first three the question is now bounded rather
  than open: each has a named destination on both the map object and the landscape object and
  a reader on neither (`ALM-META-091`, `ALM-META-092`, `ALM-CORP-093`), so the meaning is
  recoverable from a writer, not from this image; which trigger consumes which text slot.
  **Answered since:** the type-7/8/9 leaf grammar (`ALM-TRIG-044`…`047`, `049`, `050`);
  the six comparison codes and the trigger `+0xb4` once flag (`TRIG-CMP-006`,
  `TRIG-FIRE-007`); the type5 `+0x00` colour slot and `+0x04` human-participant flag
  (`ALM-PLAYER-069`, with `ALM-GRP-041`); the type5 `+0x08` scalar as far as static evidence
  reaches — destination `CPlayer+0x0c`, one family copy read with two direct call sites,
  the 17-map-side + 7-family direct-displacement population, its `{0,5000}` corpus law,
  and the six-dword CPlayer vtable (`ALM-SCALAR-087`, `ALM-SCALAR-088`,
  `ALM-SCALAR-089`, `ALM-CPLAYER-090`). The preregistered whole-image mover search is
  incomplete, so consumer effect remains Unknown;
  the last 8 bytes of the type-7 payload, which are the second and third arrays' count
  words and are zero exactly when a map has no conditions and no triggers; type-4 `+0x12`,
  which is what a `Target_Structure` parameter names; and the **64 type-3 cells**
  (`scn:131`, `scn:150`, rows 0..2 only) whose `c − 1` indexes past the 82-entry object
  array — head-of-buffer residue, not placements (`ALM-CLS-051`). A consumer must still
  bounds-check, and must also not trust an *in-range* code in those two maps' rows 0..2.
- `scenario.res` embeds 28 campaign maps (`N.alm`) that all conform to this spec —
  independent, same-header-different-body confirmation of the container.

The preserved top-level files are a separate population: the measured EN root has 38 ALMs
including archive members, RU has 34. RU top-level Horror declares 415 type-4 records but has no
type-4 section. Its decoded placement count is zero; runtime acceptance is Unknown. Do not infer
placements from metadata when their payload is absent. Ordered registration replay also separates
rectangle, presence-mask and accepted-cell counts (`UNIT-AREAPOP-075`, `UNIT-STRUCTCELL-070`).
