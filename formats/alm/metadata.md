# Metadata and cell planes

[Reference](format.md)

## type-0 metadata payload (632 bytes) — `ALM-META-008…ALM-META-010`

The type 0 payload is 632 bytes: 48 bytes of scalar fields, a 64-byte map
name, two u32 values, a 64-byte description and seven further 64-byte text
slots. — ALM-META-008

The offsets below are relative to the payload. The fixed reads total 632 bytes. — ALM-FRAME-031 (amended; framing retained)

| Off | Type | Field | Notes | Claim |
|-----|------|-------|-------|-------|
| +0x00 | u32 | **W** | map width — read into the map object | ALM-HDR-001, ALM-META-024 |
| +0x04 | u32 | **H** | map height (`W≠H` occurs, e.g. 112×144) — read into the map object | ALM-HDR-001, ALM-META-024 |
| +0x08 | f32 | **angle** | Radians; loaded into map `M+0x18` and terrain `P+0x20`. The terrain store fills four bytes of a double. Forced relight replaces it from the sun globals before the known read. Installed angles are whole degrees: `±45, ±44, ±18, 36, 28, 25, 22`; the common π/4 encoding is `0x3f490fda`. | ALM-META-027 (amended; payload-angle clause retained), ALM-META-091, ALM-META-092, ALM-CORP-093, TERR-LIGHT-149, TERR-LOAD-152 |
| +0x0c | u32 | scalar (stored) | Loaded into map `M+0x1c` and terrain `P+0x2c`; no identified reader in the named direct-consumer paths. Installed values: `360,480,645,720,1080`. Meaning Unknown; numeric agreement with the sun clock does not identify a consumer. | ALM-META-026, ALM-META-091, ALM-META-092, ALM-CORP-093 |
| +0x10 | u32 | scalar (stored) | Loaded into map `M+0x20` and terrain ambient byte `P+0x1c`. Forced relight overwrites the terrain value before its known read. Installed values span `0..33`; this is not a validation bound. | ALM-META-026, ALM-META-091, ALM-META-092, ALM-CORP-093 |
| +0x14 | u32 | scalar (stored) | Loaded into map `M+0x24` and terrain range byte `P+0x1d`; the four-byte terrain store also writes `P+0x1e/0x1f/0x20`. Relight replaces the range before its known read. Installed values span `27..64`; this is not a validation bound. | ALM-META-026, ALM-META-091, ALM-META-092, ALM-CORP-093 |
| +0x18 | u32 | bitmask | **Terrain tile-group mask.** The map loader discards its local copy; terrain stores it at `P+0x28`. Bit i selects group `(i>>2)+1`, variants `(i&3)*4 .. +3`. Installed values use bits `0..12`: groups 1–3 and group 4 variants 0–3. | ALM-META-026, TERR-LOAD-152 |
| +0x1c | u32 | **#players** | player-record count = `type5_size / 76` ; installed values `3..9` (editor caps at 16) | ALM-META-025 |
| +0x20 | u32 | **#objects** | object-record count (base records; extensions are additional bytes, not records) | ALM-META-025 |
| +0x24 | u32 | **#units** | unit-record count = `type6_size / 70` in the version-990 form | ALM-META-025 |
| +0x28 | u32 | (unused) | Read and discarded; meaning Unknown. Type7 uses its own count words | ALM-META-025 |
| +0x2c | u32 | **#type 8 records** | count used as the case-8 loop bound; `0 ⟺ type8 empty` | ALM-META-025 |
| +0x30 | char[64] | **name** | NUL-terminated ASCII. Installed names use at most 21 bytes and can be empty; 21 is not a field-width limit. | ALM-META-010 |
| +0x70 | u32 | scalar (stored) → `map+0xd4` | **Multiplayer-mode source:** `[0x005cd758]+0xc = (map+0xd4 > 1)`, read during player construction. A value 1 selects single-player loading even in a loose map such as RU `Horror.alm`; it does not identify `scenario.res` membership. A player-slot/MP-capacity interpretation beyond this Boolean consumer remains unestablished. | ALM-META-026, ALM-MODE-070 |
| +0x74 | u32 | scalar (stored) | standalone `1..5`, campaign `1` | ALM-META-026 |
| +0x78 | char[64] | **description** | NUL-terminated code-page text (ASCII or Windows-1251). Installed text uses at most 36 bytes; 36 is not a field-width limit. | ALM-META-010 |
| +0xb8 | 7×64 B | **text slots** | fixed array of 7 slots `[3×u32 prefix][char[52] text @+12]`, default `"<None>"`. Loaded (part of a 512-B block read from +0x78). Campaign maps fill slots 4/6 with trigger/quest text (`"mission complit"`, `"Start1"`, …) | ALM-META-028 |

The negative consumer statements above are limited to the named map/terrain
paths. Indirect calls, aliases and the skipped-relight/message-path cases
remain Unknown; absence of a reader is not a semantic default.
— TERR-LIGHT-149, TERR-LIGHT-150, TERR-LIGHT-151

The record's `typeId = 0` and opaque four-byte word are [record-header
fields](container.md) at `+0x0c/+0x10`; they are not payload fields.

The whole type-0 record is read by `rom.exe`'s `.alm` loader (`FUN_00512369`, case 0) in
16 reads totalling exactly **632** bytes: twelve `u32` (`W/H`, the light `angle`, three
stored scalars, the discarded bitmask, the five content-record counts), `0x40` (`name`),
two `u32`, and one `0x200` block covering `description` + the 7 text slots. The `π/4`
constant is `+0x08` (`0x3F490FDA`).

## Grid layers type 1 / type 2 / type 3 — `ALM-GRID-012`, `ALM-GRID-013`,
`ALM-GRID-014`, `ALM-GRID-032`

Each grid record stores **`W·H` cells starting at payload+0** — the payload being the
bytes after the record's 20-byte header — with `payloadSize` exactly `2·W·H` (type 1) /
`W·H` (type 2, type 3). The payload is **pure grid**: nothing is overlaid on it, no cell is
lost, and the record's `typeId`/opaque word live in the header (`ALM-GRID-032`). Cell layout is
row-major, `W` cells per row:

```
cell (col, row)  ->  element index  row*W + col      (0 <= col < W, 0 <= row < H)
type1: u16 at payload + 2*(row*W + col)
type2: u8  at payload +    row*W + col
type3: u8  at payload +    row*W + col
```

The world ingest addresses `tiles[row*W+col]`, `heights[row*W+col]` and
`objects[row*W+col]` from the payload origin. — ALM-GRID-032

| Layer | Name | Cell | Encoding | Claim |
|-------|------|------|------------------------------------------------|-------|
| type 1 | **Tiles** | u16 LE | **tile-index word**: **bits 0–9 = tile index**, **bit 13 (`0x2000`) = impassable flag**, bits 10–12/14–15 unused in installed maps. Terrain class is *derived* from the index by the loader (below), not a raw high byte | ALM-GRID-012 |
| type 2 | **Altitudes** | u8 | Height/altitude, copied to the simulation height buffer. Installed heights are below `0x80`; this is not a general u8 limit. The earlier range derived from the displaced grid origin is retracted. | ALM-GRID-013, TERR-LIGHT-016, TERR-LIGHT-028 |
| type 3 | **Objects** | u8 | Static-object code: 0 is empty; nonzero `c` selects `objects.reg` section `c-1`. Missing type 3 is zero-filled. World ingest additionally derives runtime block value 5 for nonzero cells. | ALM-GRID-014, ALM-CLS-035 |

### type 1 tile word — terrain resolution (`rom.exe`)

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

The displaced Cost/Pass vectors associated with superseded `REG-FMT-017`
are invalid. Use the per-class values in the table above.

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
specified in [TERRAIN](../terrain/format.md) (`TERR-IDX-003`, `TERR-SEM-004`): `g = (w & 0x1fff) >> 6` selects a
`terrain.3d/tileG-VV.bmp` strip (`G=(g>>2)+1`, `V=(g&3)*4+((w>>4)&3)`) and `w & 0xf`
selects the 32×32 sub-cell. `tile1/2` = land strips, `tile3` = animated water, `tile4` =
road; bit-13 tiles composite over `dirt.bmp`.

**Runtime passability** — corrected and completed by `TERR-PASS-049…TERR-PASS-051`.
The map load derives **three** 256×256 byte planes at fixed
stride 256, addressed `(row<<8)|col` regardless of `W`,`H`:

| plane | at | initial fill | built from |
|---|---|---|---|
| movement cost | `sim+0x00000` | `0x01`, overwritten for every in-bounds cell | the blended terrain cost |
| block bits | `sim+0x10000`, copied to `sim+0x20000` | `0` | type 1 + type 3 + the border |
| height | `sim+0x9451c` | `0` | type 2 (Altitudes) |

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
how a bridge crosses water. Specified in [TERRAIN](../terrain/format.md) → "Structures on the block plane"
(`TERR-STRUCT-068`…`072`, `TERR-PASS-073`); the `kind`→class resolution is `ALM-CLS-036` above.

The record-header opaque word is separate from metadata payload+0x08.
No numeric interpretation follows from its bit-pattern census. — ALM-HEADER-098

For payload+0x0c, the selected editor reader stores four bytes at object+0x1c
and its writer submits four bytes from receiver+0x1c. The writer caller uses
document+0x50. Publication between constructor and document, intervening
control changes, and native round-trip preservation remain Unknown.
— ALM-EDSCALAR-102

The game and landscape read prefixes transfer this scalar unchanged only
when their lower read completes. Resolving that stream edge does not close
later map/landscape aliases or establish meaning. — TERR-STREAM-157
