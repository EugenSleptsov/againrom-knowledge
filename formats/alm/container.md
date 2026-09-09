# Container and load rules

[Reference](format.md)

## Structure

The header gives `recordCount`. Each record supplies `payloadSize` and
`typeId`. The standard record chain ends at EOF; record count and type order
are data, subject to the loader dependencies below. — ALM-FRAME-031,
ALM-CORP-060

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
| 0x08 | u32 | dataSize | Installed version-990 value `4·W·H + 72`; the loader ignores it. The writer-side meaning of 72 is Unknown | ALM-HDR-001 |
| 0x0C | u32 | **recordCount** | the loader's loop bound, gated `≥ 3` and nothing more; `= 10` in every authored map, `= 4` in `ru/Horror.alm` | ALM-META-024, ALM-FRAME-031, ALM-CORP-060 |
| 0x10 | u32 | **formatVersion** | `= 990`; gated `≤ 1001`; selects version-conditional record fields at `0x3b6/0x3d8/0x3da/0x3dd/0x3e9`, and `== 1000` would skip record headers entirely (unexercised by any shipped map) | ALM-META-024, ALM-FRAME-031 |

`W`,`H` themselves live in the [type 0 record payload](metadata.md), not the header.
`rom.exe` `FUN_00512353` reads this header as one `Read(dest, 0x14)`.

## Record header (20 bytes) — repeated `recordCount`×, from 0x14

| Off | Type | Field | Notes | Claim |
|-----|------|-------|-------|-------|
| +0x00 | u32 | tag | installed value `7`; stored to the map object (`+0x08`) and never read again | ALM-SEC-002, ALM-FRAME-031 |
| +0x04 | u32 | hdrLen | constant `20` (= this header's own length); read, not validated | ALM-SEC-002 |
| +0x08 | u32 | payloadSize | payload byte length; the payload follows immediately | ALM-SEC-002 |
| +0x0C | u32 | **typeId** | the record's type, `0..9`; **this is the word the loader's `switch` dispatches on** (10-entry jump table) | ALM-SEC-003, ALM-FRAME-031 |
| +0x10 | f32 | **perMapConst** (`selectorA`) | a per-map `f32` (e.g. `0xBFC02B6D`), byte-identical in every record header of a map | ALM-SEC-003, ALM-META-027 |

Each record payload starts immediately after its 20-byte header and has no
additional type identifier. The loader reads the header with `Read(dest,0x14)`
except when `formatVersion == 1000`, then dispatches the payload grammar by
the header's typeId. A payload may be empty or shorter than eight bytes.

## Record roster — what the writer emits, and what the loader requires

The complete authored order is `0,1,2,3,5,4,9,8,6,7`. The primary loader
dispatches the declared recordCount, requires only types 1/2 and permits
other record sets. The four-record RU Horror.alm satisfies these structural
checks. — ALM-SEC-003, ALM-REQ-055, ALM-REQ-056, ALM-ORD-057,
ALM-META-058, ALM-RDR-059, ALM-CORP-060

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
| 8 | variable, sometimes **0** | **authored loot**, `meta+0x2c` records with no count word. The box/circle reading in `ALM-TRIG-022` is retracted; use `ALM-SACK-065` |
| 9 | variable (small) | **caster list** — `[u32 count@+0][records]` (`ALM-CODE-023`, `ALM-TRIG-049`) |

The three fixed-record sections carry their **count in the type-0 metadata**
(`+0x1c`=type 5, `+0x20`=type 4, `+0x24`=type 6, `ALM-CNT-017`); the variable trigger
sections carry their own count word.

Also in the metadata, and load-bearing for a *writer*: `+0x28` and `+0x2c` are the loader's case-7
and case-8 loop bounds, so all five of `+0x1c`..`+0x2c` must agree with the records actually present
— a zero-length record needs a zero count or the sequential parse desynchronises. — ALM-REQ-056

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
| **type 1** *(“Tiles”)* | **reject** — loader status 5, `"Tiles block not found"` | the world builder `FUN_00548550` reads `map+0x0c` as u16/cell with no null test |
| **type 2** *(“Altitudes”)* | **reject** — loader status 6, `"Altitudes block not found"` | same, `map+0x14` as u8/cell |
| type 3 | **default**: a `W·H` plane of zeroes | same consumer reads `map+0x10` as i8/cell; nonzero = “an object blocks this cell”, so all-zero is a valid empty plane |
| type 5 | **default**: one group record — `new(0x48)`, `FUN_0048d9d0(1,1)`, then `CPlayer+0x0c = 5000` (`ALM-REQ-056`), which is the field a record's own `+0x08` fills (`ALM-SCALAR-087`). Every shipped map authors `0` in that position on record 0 (`ALM-SCALAR-089`) | the group list must be non-empty |
| type 0 | **default**: `W = H = 16`, empty name and text, `+0x70 = +0x74 = 1` | but see the ordering rules — this default is only safe if types 1 and 4..9 are absent too |
| type 4, 6, 7, 8, 9 | **skip** | their loops are bounded by type-0 counts (type 9's by its own first word), which are 0 when the records are absent |
| `typeId >= 10` | **skip** by `Seek(payloadSize)`, no error | the switch's `default` arm |

Header gates, in order: magic `M7R␀`; `recordCount >= 3`; `formatVersion <= 1001`. A duplicate
`typeId` is **not** rejected — the case runs again and overwrites the pointer.

**Ordering.** The permutation is not enforced, but three precedence relations are, because two
lengths and five loop bounds are carried in state that only one earlier case writes:

1. **type 0 before type 1** — the tile grid is allocated `W·H·2` from the *metadata's* `W`/`H`
   (default `16,16`), not from the record's own `payloadSize`.
2. **type 2 before type 3** — type 3's allocation and read length is the **type-2 record's**
   `payloadSize`.
3. **type 0 before types 4..9** — those cases' loop bounds are metadata fields.

**Four readers ship, with three different acceptance tests** (`ALM-RDR-059`). Beside the loader,
`rom.exe` has the `*.alm` browser scan's accept test — magic, `recordCount >= 2`, a type-0 record,
and type-0 payload **`+0x70 >= 2`**, so a map below that threshold is silently unlisted
(`ALM-META-058`) — and a light parser that handles typeIds 0..3 only, walks to EOF rather than by
count, and treats its two header failures as advisory message boxes without stopping. `Map
Editor.exe` carries that light parser's two literals and none of the loader's six; its own reader is
unread.

## End of file

The standard chain ends after its final payload. There is no separate
eight-byte trailer; ALM-TRL-005 and ALM-TRL-030 are retracted. The final two
words in a type 7 payload with no conditions or triggers are its zero count
words. — ALM-FRAME-031, ALM-TRIG-044

## Reading algorithm

```
assert bytes[0:4] == "M7R\0"
hdrLen    = u32(bytes, 0x04)             # == 20
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
#          kind in {0x22,0x23} -> Shop, else Building
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

## Write sequence

1. Select the format version. This reference's fixed-size placement layouts
   describe version 990; apply each documented version branch for other values.
2. Build metadata and payloads. Derive counts from the actual records and use
   `W*H` cells in each present grid. A metadata count alone creates no payload.
3. Emit the 20-byte file header: magic, header length, dataSize, recordCount,
   formatVersion. Version-990 installed maps use `dataSize=4*W*H+72`; its
   decomposition remains Unknown.
4. Emit each record header followed by its payload. The installed full-map
   order is `0,1,2,3,5,4,9,8,6,7`. Preserve metadata before dependent records,
   type 2 before type 3, and type 9 before type 8 before type 6.
5. Emit exact type-specific counts, optional fields and extensions. Do not add
   a trailer or duplicate counts where the grammar uses metadata.
6. Validate record endpoints, grids, class/owner relations and script IDs.
   A type 7 trigger refers to a node's ID, not its array position.

The structural size relation is `20 + sum(20+payloadSize) = fileSize` for
the standard complete form. Types 1 and 2 are required by the primary loader;
the complete installed type set is not an acceptance requirement.
— ALM-FRAME-031, ALM-ORD-057, ALM-ORD-068, ALM-TRIG-044,
ALM-TRIG-045, ALM-TRIG-046, ALM-TRIG-047, ALM-CORP-060

## Unknowns and compatibility

| Field or path | Limit |
|---|---|
| Header dataSize | Installed `4*W*H+72`; ignored by the loader; writer accounting of 72 Unknown |
| Type0 scalars +0x0c/+0x10/+0x14/+0x74 | Exact widths/destinations retained; general authoring meanings Unknown |
| Type5 +0x08 | Loaded to `CPlayer+0x0c`; copy use is known, complete gameplay consumer effect Unknown |
| Type6 +0x1c/+0x22/+0x26/+0x28/+0x2a/+0x30/+0x33/+0x34/+0x35 | Known stored fields; no named read in the described original field-consumer set; authoring roles Unknown |
| Trigger +0x40 | 64 raw bytes; no meaningful text grammar established |
| Trigger Target_Item | Authored values exist; complete value domain/meaning Unknown |
| Type9 +0x00 | Exact tag width known; meaning Unknown |
| Type0 text slots | Fixed layout; complete trigger-to-slot relation Unknown |
| Editor reader | Separate reader; complete acceptance contract Unknown |

— ALM-HDR-001, ALM-META-091, ALM-META-092, ALM-CORP-093,
ALM-SCALAR-087, ALM-SCALAR-088, ALM-SCALAR-089, ALM-CPLAYER-090,
ALM-TAILDIR-081, ALM-TAILU16-082, ALM-TAILVER-084,
UNIT-PLACESKILL-086, ALM-TRIG-047, ALM-TRIG-049, ALM-TRIG-050

The placeable-definition source is [Data.bin](../databin/format.md), including
footprints and actor fields. The earlier unknown source-file clause is
resolved by DAT-LOC-001 and DAT-ACT-006.

The first three rows of campaign maps 131 and 150 contain type 3 residue,
including 64 codes whose `c-1` exceeds the 82-entry object array. Bounds-check
all object codes; an in-range code in those rows does not establish a valid authored placement. — ALM-CLS-051

The preserved RU loose `Horror.alm` declares four records and has no type 4 or
type 6 payload despite nonzero metadata counts. Decode no placements from
those absent records; native runtime continuation remains Unknown. A Building's
rectangle, selected presence-mask cells and successfully attached cells are
distinct quantities. — ALM-CORP-060, UNIT-AREAPOP-075, UNIT-STRUCTCELL-070
