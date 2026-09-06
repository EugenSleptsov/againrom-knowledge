# REG — registry (nested `&YA1` key-value tree) — specification

Level 3. Promoted, evidence-backed claims only. Claim basis: location and key catalogues
(`REG-LOC-016`, `REG-LOC-038`, `REG-KEY-044`); values
(`REG-VAL-024`…`REG-VAL-029`); framing,
record layout, kinds and text handling from `rom.exe`'s registry class
(`REG-FMT-031`, `REG-REC-032`, `REG-KIND-033`, `REG-KIND-034`, `REG-TEXT-036`,
`REG-TEXT-037`); corrected
per-key tables (`REG-VAL-043`); and the widened registry/roster tables
(`REG-CUT-053`, `REG-ROSTER-052`).

> **⚠ Correction (`REG-FMT-031`, `REG-REC-032`).** Everything published before this
> correction placed the record
> fields **4 bytes early**, so each key's name was paired with the **next** key's value.
> The numbers were real; the key each belonged to was wrong. Any decoder written against
> the old spec reads every scalar off by one key. See `claims/registry.md` → Standing
> corrections.

Seen as: file nodes ending `.reg` inside a RES container. The install holds **44** of
them in **12** containers (`REG-LOC-038`) — the five class registries in `graphics.res`,
three in `scenario.res`, one in `sfx.res`, two in `world.res`, and **33 cutscene
registries** in `Allods/VIDEO4.RES` (18) / `VIDEO8.RES` (15). **43 of the 44 now have a
published per-key value table**; the exception is `data/map.reg`.

**Selected EN/RU member equality (`SAV-SUFF-301`).** Read-only payload hashing finds five
located SAV-reconstruction inputs byte-identical across the preserved roots:
`scenario.res/globalmap.reg`, `scenario.res/npc.reg`, `scenario.res/scenario.reg`,
`world.res/data/map.reg` and `world.res/data/ai.reg`. This is exact byte equality for
the named payloads, not a claim of runtime-derived semantic equality, crossed-root
compatibility or equality of unmeasured REG members. — SAV-SUFF-301

## At a glance

A `.reg` payload is **itself an `&YA1`** (same magic as the outer RES container,
`REG-LOC-016`) but used as a **hierarchical key-value store**, not a file archive. After a
24-byte header it is a flat array of fixed 32-byte **records** forming a key/subkey tree,
then a `u32` pool length, then a **pool** holding string and array values (`REG-FMT-031`).

```
+---------+-----------------------------+-------+--------------------------+
| header  |   records (32 B each)       |poolLen|   string / array pool    |
| 0x18 B  |   key & subkey tree         | u32   |   DescText, File, arrays |
+---------+-----------------------------+-------+--------------------------+
0        0x18                    0x18+R*32   +4                          EOF
```

### Header (0x18 B) — `REG-FMT-031`

`rom.exe`'s loader (`FUN_004cae80`) performs exactly six 4-byte reads before the bulk
read, which is what fixes the header length. The first four land on the registry object's
own first 16 bytes — i.e. **the header is the root node**, minus its name (which the Open
call fills in from the file path, never from the file).

| Off | Field | Notes |
|-----|-------|-------|
| 0x00 | magic | `&YA1` = `0x31415926` LE; mismatch → the game's `"bad signature"` |
| 0x04 | root `value` | index of the root's first child record |
| 0x08 | root `size` | number of top-level children |
| 0x0C | root `kind` | always **17** = `0x11` = subkey (bit 0) \| sorted (bit 4) — not a magic constant |
| 0x10 | `R` | total record count |
| 0x14 | — | read into the object but used by no accessor found; **Unknown** |

Then `R × 32 B` of records at `0x18`, a `u32` **pool byte length**, and the pool:

```
poolStart = 0x18 + R*32 + 4
0x18 + R*32 + 4 + poolLen == payloadLen      exactly, on 44/44 registries
```

## Record (32 bytes) — `REG-REC-032`

| Off | Type | Field | Notes |
|-----|------|-------|-------|
| 0x00 | u32 | — | explicitly zeroed by the writer; read by no accessor found; **Unknown** |
| 0x04 | u32 | `value` | int32 · pool byte offset · child-block **start** index · **low dword of a double** |
| 0x08 | u32 | `size` | byte length · child count · **high dword of a double** |
| 0x0C | u32 | `kind` | bitfield, see below |
| 0x10 | char[16] | `name` | 15 significant characters + NUL |

Stride is 32 (`SHL 0x5` in the loader's allocation, the lookup's addressing and the
node-insert's placement).

> **Which word is the value.** In absolute file terms record `i` occupies
> `[0x18 + 32i, 0x18 + 32(i+1))` and its value word is **`0x1C + 32i`** — eight bytes
> *before* the `kind` word. The word *after* the name, `0x3C + 32i`, is record `i+1`'s
> value; reading it as record `i`'s is the superseded four-bytes-early error. The same word serves an
> integer key and a string key's pool offset — there is no per-kind displacement.

This is attested four independent ways, so no single misread can carry it:

| # | function | what it shows |
|---|---|---|
| 1 | loader `FUN_004cae80` | six 4-byte header reads, then `Read(records, R<<5)` — fixes the header at `0x18` and the array origin |
| 2 | lookup `FUN_004ce8e0` | `records + node[+0x4] * 32`, iterate `node[+0x8]`, compare `node[+0x10]` |
| 3 | accessors `FUN_004ccc20` / `FUN_004cc670` | `*(u32*)(node+4)` for an int; `poolBase + *(int*)(node+4)` with length `*(u32*)(node+8)` for a string |
| 4 | node-insert `FUN_004ce6d0` | places a new child at `records[parent[+0x4] + parent[+0x8]]` and does `parent[+0x8]++`; writes the name at `node+0x10` **clamped to 15 chars** (so the field is 16 B, not 20); stores `0` into `node+0x00` |

### Kind — a bitfield, not an enum (`REG-KIND-033`)

```
type       = kind & 0x0E          tested as AND 0xe / CMP by every accessor
bit 0      = node is a subkey
bit 4      = children are sorted -> the lookup bsearch()es instead of scanning
bit 28     = the name was longer than 15 chars and got truncated (no shipped record sets it)
```

### Kinds and their storage (`REG-KIND-034`)

| kind | meaning | where the value lives | corpus count |
|-----:|---------|-----------------------|-------------:|
| 0 | string | `poolBase + value`, `size` bytes **including** the NUL | 873 |
| 1 | subkey | children = records `[value, value + size)` | 512 |
| 2 | int32 (signed) | `value` | 2 792 |
| **4** | **double** | **`value` = low dword, `size` = high dword — 8 bytes in the record, no pool access** | 122 |
| 6 | int32[] | `poolBase + value`, `size/4` LE int32s | 322 |
| 10 | double[] | `poolBase + value`, `size/8` LE doubles | 0 |

> **`kind` belongs to the record, not to the key name** (`REG-KIND-056`). Nine key names in
> the install carry two kinds. Four are the empty-string-vs-array pair (`REG-KEY-045`);
> the other five are in `scenario.reg`, where `Mercenaries`, `InnNPC`, `InnMission`,
> `EnableMercenary` and `AddTextDocument` are spelled as a **bare int32** in some sections and
> an **int32 array** in others. A typed reader must accept a scalar where an array is expected
> and treat it as a one-element list; reading only kind 6 loses 18 of the 31
> `InnNPC`/`InnMission` records. How `FUN_004cd240` itself handles a kind-2 record is Unknown.

`rom.exe`'s string converter carries one further case (`(kind>>1)&7 == 4`, i.e. kind 8/9:
skip 4 bytes of the pool item, copy `size − 4`, terminate with two NULs). **No writer in
the binary produces it and no record in the corpus is one — its meaning is Unknown**, and
this spec deliberately does not guess. Anything else → `"UNKNOWN TYPE. CANT CONVERT"`.

The game's own names for the types, from its error strings: *int*, *double*, *int array*,
*double array*, *string array or single string*.

### Kind 4 — the double (`REG-DBL-035`)

```
reader  FUN_004ccda0 :  AND EDX,0xe ; CMP DL,0x4 ; FLD qword ptr [node + 0x4]
writer  FUN_004cb340 :  *(double *)(node + 4) = atof(text) ; node->kind = 4
```

The 122 instances in the install are `startfade` / `endfade` under `Fading<n>` in 31 of the
33 cutscene registries, and every one is exactly `0x0000000000000000` (0.0, x61) or
`0x3FF0000000000000` (1.0, x61) — a fade-in then a fade-out, matching each file's own
`nFadings`. Reading any other 8-byte window of the record yields denormal garbage, which is
how the layout was falsified rather than fitted.

### Pool encoding (`REG-VAL-025`, as amended)

Each kind-0/6/10 record stores its pool byte offset **directly** in its own `value` and
its byte length in `size`; `rom.exe` reads `poolBase + node->value`. There is no cursor
and no "snapshot slot" — the superseded `record[i−1].C` rule was the 4-byte shift of
this same fact. Items tile the pool exactly: no gap, no overlap, no tail, on all 44
registries.

### Tree walk (`REG-VAL-028`, as amended)

A subkey's `value` is the **start** index of its child block and `size` is the child
count; the lookup computes `&records[parent->value]` and iterates `parent->size` entries.
The root block plus every subkey block cover the record array **exactly once** — 0
unreachable records, 0 doubly-reached — on all 44 registries. No tiling carve-out is
needed.

## Text (`REG-TEXT-036`, `REG-TEXT-037`)

**`rom.exe` performs no byte-to-character conversion on registry text.** The loader
bulk-reads the record array and the pool verbatim; the string accessor is a `strncpy` from
`poolBase + value`. No code page is named on the path, and no Windows conversion API is
reachable from the registry class — `MultiByteToWideChar` / `WideCharToMultiByte` /
`LCMapStringA` / `GetACP` / `GetOEMCP` are called only from the statically-linked MFC/CRT,
the `CharToOemA` / `OemToCharA` thunks have zero callers, and `setlocale` appears nowhere.

The only byte transform anywhere on the path is the *name lookup*'s `_strnicmp`, whose
fast path folds `A`–`Z` only and leaves bytes `>= 0x80` untouched.

> **Therefore a code page for `.reg` is a rendering choice of ours, not a property of the
> format.** Names and string values are **opaque bytes**. Consumers that must display them
> should state the convention they apply (this repo's tools print ASCII and escape the
> rest); CP866 and CP1251 are both defensible renderings and neither is the game's.

Corpus check with its parameters: sweeping every file under the install, classifying
nested `&YA1` payloads by root `kind == 17`, taking names as the 16 bytes at record `+0x10`
truncated at the first NUL and string values as kind-0 pool items — **0 bytes `>= 0x80`**
occur in 4 621 name fields or 873 string values. Range `0x20…0x7A`. No shipped registry
datum exercises a code page at all.

## Name matching, ordering and the 15-character clamp

Keys are looked up by `_strnicmp` over **15** characters, case-insensitively in the ASCII
range. Two keys agreeing in their first 15 characters are indistinguishable to the game.
When a subkey's `kind` has bit 4 set, its children are assumed **sorted** and are found by
`bsearch`; otherwise by linear scan.

**Bit 4 is set on exactly one node per registry: the root** (`REG-KEY-054`). Over all 44
registries the root's `kind` is `17 = subkey|sorted` 44/44 and **0 of the 512 subkey records**
set it, so the `bsearch` path is reachable only for a top-level section name and every per-key
lookup inside a section is a linear scan. The data honours the obligation that creates: the
root's child list is sorted under the lookup's own comparator on **44/44** (only 81 of the 512
non-root lists happen to be). Two consequences:

- Sections are stored in **lexicographic**, not numeric, order — `MapObject1, MapObject10,
  MapObject11, …, MapObject2`. The order in which a loader's `"<Prefix>%d"` counter visits
  sections is **not** the record order.
- A writer must emit root children sorted, or the game will `bsearch` an unsorted list.

**The clamp, from the data side** (`REG-NAME-055`): over all 4 621 shipped records, name
lengths run 1..15, **0** are 16 or longer, **69** are exactly 15, and **0 sibling pairs
anywhere in the install are indistinguishable** under the 15-character compare — the clamp
never bites on shipped data. A 15-character key may still be a truncation of a longer authored
name, and `rom.exe` sometimes carries the original: `MinimalGuardRan` → `MinimalGuardRange`,
`AddPictureDocum` → `AddPictureDocument`, `ScenarioMission` → **`ScenarioMissionCount`**.

## Cross-checks that the framing has to satisfy

Measured on our own bytes under both framings (`tools/regtext -grounds`); each is a place
where an off-by-one key would show:

| observation | this spec | superseded four-bytes-early spec |
|---|---|---|
| `u32` at `0x18 + 32*R` vs sum of `size` over kind-0/6/10 nodes | equal, **44/44**, and the pool then ends at EOF **44/44** | that word is name padding, with no account of why it reads 5967 / 4161 / 3060 |
| `[Global] UnitCount` vs actual `Unit*` sections | `34` vs 34 | `33` vs 34 |
| `[Global] FileCount` vs actual `Files` children | `33` vs 33 | `0` vs 33 |
| `[Unit0] AttackPhases` vs `AttackAnimTime`/`Frame` lengths | `7` vs 7, 7 | `128` vs 7, 7 |
| `[Unit0] MovePhases` vs `MoveAnimTime`/`Frame` lengths | `8` vs 8, 8 | `2` vs 8, 8 |
| `[Unit0]` selection box | `X1 48 < X2 80`, `Y1 48 < Y2 90` | `X1 80 > X2 48`, `Y1 90 > Y2 16` — inverted |
| kind-4 values | exactly `0.0` and `1.0` | a kind-4 node's byte length reads `0x3FF00000` |
| Counter censuses (`REG-CUT-053`, `REG-SFX-057`, `REG-SCN-059`) on the registries this table never used: 37 declared counters (`SfxCount`, `ObjectCount`, `TotalMissions`, `ScenarioMissionCount`, 33× `nFadings`) against counts derived by walking section names | agrees **37/37** | agrees **1/37** — and the one hit is `TotalMissions`, whose successor record happens to hold the same 15 |

(`CenterX = 64` is the midpoint of `SelectionX1/X2` on `Unit0`; `CenterY = 78` is *not* the
midpoint of `SelectionY1/Y2` = 69. The X relation is noted, not generalised.)

## `units.reg` — unit classes (`REG-UNITS-018`, `REG-ROSTER-019` (partially retracted, corrected by `REG-ROSTER-052`))

34 classes `Unit0..Unit33` (`[Global] UnitCount` = **34** on the corrected framing —
`REG-UNITS-018`; the `33` this spec used to print was the
off-by-one, and `33` is `[Global] FileCount`) under a root with a shared `Files`
sprite-descriptor section. Per-class keys (all game-authored ASCII):

| Key | Meaning |
|-----|---------|
| `ID` | class id |
| `File` | **scalar index into `Files[]`** — *not* a path. `units/` + `Files[File]` + `.256` (`REG-VAL-029`; see *Resolved values* below) |
| `Index` | frame index / base |
| `Palette` | palette ref |
| `DescText` | human-readable name (`Human Two-Handed Swordsman`, `Ogre (Lord)`). **Nothing in `rom.exe` reads it** (`UNIT-NAME-039`), it is byte-identical on both roots on 34/34 classes and pure ASCII, and its longest shipped value is 31 bytes + NUL against the `0x20` inline field — so it is editor-facing metadata with a 31-character ceiling (`REG-DESC-096`) |
| `InfoPicture` | portrait leaf, **without directory or extension** — the engine builds `graphics\infowindow\<InfoPicture>.bmp`, or `<InfoPicture><tier>.bmp` with the digit dropped at tier 1, from `class+0xd8` (`UNIT-PICT-036`). Read only when the drawable's `+0x18c & 0x11` is clear, i.e. for `ID >= 0x1a`; on the thirteen human classes below that it is **dead data** (`UNIT-PICT-035`) |
| `Parent` | inherit-from class |
| `InMapEditor` | editor-visible flag |
| `MovePhases`, `MoveBeginPhases` | walk / walk-start phase counts |
| `AttackPhases`, `DyingPhases`, `BonePhases`, `IdlePhases` | per-state phase counts |
| `Flip` | **selects the sheet layout, not a blit option** — see *The loaded class record* below (`REG-UNITS-051`) |
| `Width`, `Height` | frame canvas |
| `CenterX`, `CenterY` | draw anchor |
| `SelectionX1/X2/Y1/Y2` | selection box |
| `Z` | draw layer |
| `TileSize` | footprint |
| `Dying`, `AttackDelay` | death behaviour / attack cooldown |
| `Projectile`, `ShootDelay`, `ShootOffset` | ranged-only (present on ~9–11 classes) |
| `AttackAnimTime/Frame`, `MoveAnimTime/Frame`, `IdleAnimTime/Frame` | per-phase kind-6 tracks |
| `Sound` | sound ref |

**Resolved values (`REG-VAL-029`).** The root's 36 children are `Files`, `Global` and the
34 classes. Each class's sprite is the `Files` leaf table (`File0..File32`) indexed by the
class's own **`File` scalar**: `units/` + `Files[File]` + `.256` — e.g. Human Archer has
`File = 8` -> `humans\archer\archer`, Human CrossBowMan `File = 9` -> `humans\xbowman\xbowman`.
(The superseded four-bytes-early spec said "indexed by `ID`"; that was the off-by-one — the shifted `ID`
slot held the true `File` value.) **The full `DescText` → sprite table is `REG-ROSTER-052`:**
all 34 classes carry their own `File` and
`DescText`, all 34 resolve to a `.256` node that exists, and `File` equals the section's own
`%d` on **33/34** — the exception is `Unit33` "Unarmed Fighter with Shield" (`File = 21`),
which shares `Unit21` "Unarmed Fighter"'s sheet, the table's only shared entry. The pairing
`REG-ROSTER-019` published, `Files[class.ID]`, **differs on 34/34 and is out of range on
14/34**; do not consume it. `Flip`'s only clean value is `1` (10 classes).
`Projectile` / `ShootDelay` / `ShootOffset` appear on the ranged/caster classes and carry
readable values (Human Archer: `Projectile = 1`, `ShootDelay = 21`, `AttackDelay = 19`).
The `Anim*Frame` / `Anim*Time` / `ShootOffset` / `Sound` kind-6 arrays decode as per-phase
int tracks.

> **⚠ "and their lengths agree with the matching `*Phases` scalar" is withdrawn**
> (`REG-UNITS-018`). Re-derived on this framing: `AttackPhases` agrees on 24/24
> classes carrying both, but `MovePhases` on 15/18 (`Goblin` 8 vs 10, `Ghost` 3 vs 4,
> `Goblin slinger` 8 vs 10) and `IdlePhases` on 4/5 (`Ghost` 3 vs 4);
> `MoveBeginPhases`/`DyingPhases`/`BonePhases` have **no** `*AnimTime`/`*AnimFrame` key,
> so no length exists to compare. `objects.reg` `Phases` vs `len(AnimationTime)` agrees on
> 7/16, `structures.reg` on 12/14. A kind-6 array carries its own length (`size/4`,
> `REG-KIND-034`) and the loader iterates **that**, never `Phases`. Read the array length;
> treat `*Phases` as an independent scalar.

> **`Width`/`Height` withdrawn and re-measured (`REG-VAL-030`); SPR256 comparison
> resolved (`REG-VAL-043`).** The published `Width`/`Height` = `128×64` (Dragon
> `160×80`) was measured on the shifted framing and is withdrawn for good — no table in
> this repo should quote it. Corrected: `Width`/`Height` = **128×128** on 17 classes
> (Dragon **160×160**), `CenterX` = 64 (Dragon 80), `CenterY` per-class 58…84. Whether
> that canvas matches the SPR256 frame bounds was left open by `REG-VAL-030`;
> `REG-VAL-043` re-measured it with the corrected canvas **and** the corrected `Files[class.File]`
> sprite index (not `Files[class.ID]`) and it resolves cleanly: **18/18 within-canvas,
> 1/18 exact** (`Ogre (Lord)`, 128×128 = 128×128) — the canvas is a per-class draw-bounds
> box that contains, not equals, the sprite's own pixel extent. The "placeholder 128
> before `Width`" observation was `Width`'s own value read one key early and is not
> re-examined.

### `Z` picks the C++ class — and only the draw layer (`REG-UNITS-061`)

`FUN_004104e8` is the only site in `rom.exe` that allocates a `0x1b0` unit object, and it branches
on `class+0x104` = **`Z`**:

```
0041152d  CMP dword ptr [ECX + 0x104],0x0     ; Z, on classes[DAT_005eb674][subscript]
00411534  JZ  -> PUSH 0x1b0 ; CALL 0045ae30   ; Z == 0 -> CUnit      (CRuntimeClass 0x599208)
00411536      PUSH 0x1b0 ; CALL 0x00461560    ; Z != 0 -> CAirUnit   (CRuntimeClass 0x599220)
```

`CAirUnit` derives from `CUnit` and overrides exactly one behavioural slot of 33 (`vt+0x38`): the
draw-layer registration, layer **3** unconditionally against `CUnit`'s **2** alive / **4** dead
(the corpse stage is `unit+0x15a`, `REG-UNITS-050`). It is a z-order flag on the *drawable*; it
does not reach the simulation, whose mover fields are per-instance bytes no registry key can touch
(`TERR-MOVE-055`). Corpus: `Z != 0` on **2 of 34** classes — `Sonic Bat` (ID 70) and `Dragon`
(ID 71), both `Z = 96`.

### The loaded `units.reg` class record (`REG-UNITS-049`, `REG-UNITS-050`, `REG-UNITS-051`)

`FUN_0046ba80` allocates **`0x12c`** bytes per class and fills, in this order:

```
+0x04 ID            +0x08 section counter  +0x0c File         +0x10 Index
+0x14 MovePhases    +0x18 MoveBeginPhases  +0x1c AttackPhases
+0x20 DyingPhases   +0x24 BonePhases       +0x28 IdlePhases
+0x2c Width         +0x30 Height           +0x34 CenterX      +0x38 CenterY
+0x3c CArray<int> Move   timeline   (m_pData +0x40, m_nSize +0x44)   len copy +0x50
+0x54 CArray<int> Attack timeline   (m_pData +0x58, m_nSize +0x5c)   len copy +0x68
+0x6c CArray<int> Idle   timeline   (m_pData +0x70, m_nSize +0x74)   len copy +0x80
+0x84 SelectionX1  +0x88 SelectionY1  +0x8c SelectionX2  +0x90 SelectionY2
+0x94 Dying        +0x98 Palette      +0x9c.. palette objects  +0xac.. 0x400-B buffers
+0xbc CArray Sound (m_nSize +0xc4)
+0xd0 TileSize     +0xd4 Projectile   +0xd8 InfoPicture[0x10]
+0xe8 CArray ShootOffset (m_nSize +0xf0)
+0xfc ShootDelay   +0x100 AttackDelay +0x104 Z  +0x108 Flip  +0x10c DescText[0x20]
```

Defaults when there is no `Parent`: **`-1`**, except `IdlePhases`, `Dying`, `Palette`,
`Projectile`, `ShootDelay`, `AttackDelay`, `Z` and `Flip` (**`0`**) and `TileSize` (**`1`**).
**`File` inherits here** — the opposite of `objects.reg`, where its default is the literal
`-1` (`REG-OBJ-046`).

The three timelines at `+0x3c` / `+0x54` / `+0x6c` are **built, not stored**: the same
run-length expansion `REG-OBJ-046` describes, run once per `<track>AnimTime`/`<track>AnimFrame`
pair. Only the **move** arm of the renderer takes the animation phase modulo the length; the
idle and attack arms index the array raw.

**`Dying` is an `ID`, and it names the class whose sheet holds the corpse.** The dying and
corpse draw arms substitute `classes[Dying]` wholesale — its `File`, its `Flip` and all six of
its phase counts. A live unit's own `DyingPhases`/`BonePhases` therefore describe frames drawn
for whoever names *it*; 20 of 34 classes name themselves, and the humans name a sibling.
The live object's corpse stage is a byte: `0` alive, `1` the last dying frame, `>= 2` bone
frame `stage - 2`. See `SPR256-UNIT-024` for what the counts index into.

**`Flip` selects between two sheet layouts.** `0` -> 16 standing frames, 8 stored directions,
never mirrored. Nonzero -> 9 and 5, with the other half drawn by reflection and a mirror flag
handed to the blit. It is never passed to a blitter as data.

## Which key a class table is indexed by (`REG-KEY-044`)

The three graphics class registries are loaded into flat arrays — `units.reg` →
`DAT_005eb674` (`FUN_0046ba80`), `objects.reg` → `DAT_005ef8c4` (`FUN_0046af50`),
`structures.reg` → `DAT_005eb62c` (`FUN_0046e8a0`) — by walking
`i = 0 … [Global]<Count> − 1` and formatting the section name `"<Prefix>%d"`. `ID` is
stored on the class (`+0x04`; structures `+0x0c`) and the counter `i` on the class too
(`+0x08`; structures have no such field). **But the array subscript the class is stored
at is per-registry, and only `objects.reg`'s is the section index** (amended
2026-07-27 — the earlier text said all three were):

| registry | array | store | subscript | `m_nSize` | slots vs classes |
|---|---|---|---|---:|---|
| `units.reg` | `DAT_005eb674` | `SetAtGrow(ID)` @ `0046c826` | **`ID`** | `maxID+1` = 81 | 34 populated, **47 NULL** (index 0 among them) |
| `objects.reg` | `DAT_005ef8c4` | `Add()` @ `0046b57e` | append order **= section index = `ID`** | 82 | 82, no holes |
| `structures.reg` | `DAT_005eb62c` | `SetAtGrow(ID)` @ `0046ee5f` | **`ID`** | `maxID+1` = 67 | 66 populated, index 0 NULL |

| registry | `[Global]` counter | sections | `ID` domain | relation to the section index |
|---|---|---:|---|---|
| `units.reg` | `UnitCount` = 34 | `Unit0..Unit33` | `1..80`, 34 distinct | **neither** — `ID == index` 0/34, `index+1` 0/34 (sparse) |
| `objects.reg` | `ObjectCount` = 82 | `Object0..Object81` | `0..81`, 82 distinct | `ID == index` on **82/82** |
| `structures.reg` | `Count` = 66 | `Structure0..Structure65` | `1..66`, 66 distinct | `ID == index + 1` on **66/66** |

So a `.alm` placement's `ID` (a type-4 record's `kind` is a `structures.reg` `ID`, a
type-6 record's `+0x08` a `units.reg` `ID` — `ALM-CLS-036`, `ALM-CLS-038`) is used
**directly as the array subscript; no translation.** A units-array consumer must
null-check instead, because 47 of the 81 slots are empty. **No lookup by `ID` over these
arrays exists anywhere in the image** — and for units and structures none is needed,
because the array *is* the by-`ID` index.

`Parent` is an **`ID`** as well (`array[Parent]`, `0046bc5f/65` and `0046b122/28`), read
by the units and objects loaders only — the structures loader never reads the key, and
`Parent` occurs on 0/66 structures. Inheritance is **eager, per key, at load**; nothing
resolves lazily at an accessor. **The guard is not the same for scalar and array keys**
(`REG-KEY-045`):

| keys | mechanism | inherits when | an explicit empty value |
|---|---|---|---|
| `File` + the 19 scalars `Index`…`AttackDelay` | the parent's resolved field is passed as the *default* argument of `FUN_004ccc20(section, key, default)` | the record is **absent** (`004ccc42`) | **overrides** the parent — `0` is returned as `0`; a non-int kind *throws* |
| `Move`/`Idle`/`AttackAnimTime`, the three `*AnimFrame`, `ShootOffset` | `FUN_004cd240` reads into a `CArray`; the loader then re-reads the **parent's section** by the name formatted from the parent class's `+0x08` | `dest.m_nSize == 0` **after** the child read (`0046c607/0b`) — presence is never tested | **does not** clear it: a present record with `kind == 0 && size < 2` takes a *success* path `SetSize(0,-1)` (`004cd2b6`), so the guard fires and the parent's array is inherited |

**There is therefore no way to spell "clear this inherited array".** A consumer must treat
an empty array key as *absent*, and must not extend that to scalars. Pinned on the corpus:
`units.reg`'s `AttackAnimTime`/`AttackAnimFrame` are kind 6 on 32 records and **kind 0,
size 1** on exactly one — `Unit33` (`ID` 2, `Parent = 3`), which therefore inherits
`Unit0`'s `AttackAnimTime = [2 2 1 1 1 2 2]` despite writing `= ""`.

`objects.reg` cannot witness the `ID`-vs-section-index question: its two keys are equal
on 82/82 *and* it is the one registry whose loader genuinely appends in order.

## Sibling registries

The 4-bytes-early tables carried historically by partially retracted `REG-VAL-030` are
superseded. The corrected framing is `REG-FMT-031`/`REG-REC-032`, and the corrected per-key
measurements are `REG-VAL-043` and the rows below.

- **`objects.reg`** — **82** objects (`REG-OBJ-039`; the previously published 56 was
  `[Global] FileCount`, the 56-entry sprite-file table size, transposed for
  `ObjectCount` — a slip independent of the framing bug that the framing re-check
  surfaced). Keys: `Index`, `Phases`, `DeadObject` (−1..67), `FireObject` (−2..−1),
  `IconID`, `AnimationTime/Frame`, `Width/Height` (32..128), `CenterX/Y`, `Parent`
  (`REG-OBJ-020`, `REG-OBJ-039`). The **in-memory class record** is `0x6c` bytes and
  `REG-OBJ-046` reads it off the loader `FUN_0046af50`:

  ```
  +0x04 ID        +0x08 section counter   +0x0c File      +0x10 Index
  +0x14 Phases    +0x18 Width             +0x1c Height    +0x20 CenterX
  +0x24 CenterY   +0x28 CArray vptr       +0x2c m_pData   +0x30 m_nSize
  +0x34 m_nMaxSize +0x38 m_nGrowBy        +0x3c = m_nSize +0x40 FireObject
  +0x44 DeadObject +0x48 InMapEditor(def 0) +0x4c DescText[0x1f]
  ```

  Two things a key list cannot show. **`File` does not inherit** — its `FUN_004ccc20`
  default is the literal `-1`, unlike every other scalar here, whose default is the
  parent's already-resolved field. And the array at `+0x2c` is **built, not stored**: the
  loader run-length expands the two kind-6 tracks, appending `AnimationFrame[i]` exactly
  `AnimationTime[i]` times, and `+0x3c` holds the result's length. That length is the
  modulus of the renderer's animation phase and the array is what it indexes
  (`TERR-SPR-042`). Shipped lengths: 0 ×52, 24 ×2, 28 ×21, 105 ×7 — and 39 of the 82
  classes get theirs only through the `Parent` array fallback.

  `DeadObject` and `FireObject` are **not the same kind of field** (`REG-OBJ-047`).
  `DeadObject` is a subscript into this same class array — −1 on 54 classes, otherwise one
  of 21 section indices, all in range, all naming a class whose `DescText` ends "(dead)".
  `FireObject` has one recovered reader in `rom.exe`, the ambient-sound loop
  `FUN_0041e28d`: its whole shipped value space is `{-2 ×21, -1 ×61}`, and the reader only
  compares it against `0` and `-2` to pick which sound slot a burnt cell feeds. The
  21 classes carrying −2 are exactly the 21 that are some class's `DeadObject`. Neither
  key occurs in any other registry in the install.
- **`structures.reg`** — 66 structures (held, `REG-STR-040`, partially retracted; see
  retracted.md). The **in-memory class record**
  is `0xa4` bytes and `REG-STR-080` reads it off the loader `FUN_0046e8a0`:

  ```
  +0x04 CSprite256 <File>.256    +0x08 CSprite256 <File>b.256   +0x0c ID
  +0x10 TileWidth   +0x14 TileHeight   +0x18 FullHeight   +0x1c Phases
  +0x20 SelectionX1 +0x24 SelectionY1  +0x28 SelectionX2  +0x2c SelectionY2
  +0x30 ShadowY     +0x34 AnimMask*    +0x38 live cells   +0x3c timeline length
  +0x40 CArray(timeline)             +0x54 Picture[0x10]  +0x64 Indestructible
  +0x68 DescText[0x20]  +0x88 VariableSize  +0x8c Usable  +0x90 Flat
  +0x94 LightRadius  +0x98 LightPulse  +0x9c CString path  +0xa0 sheets-loaded
  ```

  Three things a key list cannot show. **There is no `Parent` key in this loader**, so a
  structure class inherits nothing — unlike `objects.reg`/`units.reg` (`REG-KEY-045`).
  **`File` is a path string**, not a `Files[]` index: the class holds
  `"graphics\structures\" + File` and opens `<path>.256` / `<path>b.256` **lazily**, on the
  first draw, through `FUN_0046e720` (`REG-STR-081`). And **`AnimMask` is a picture of the
  sprite grid, not of the footprint** — its buffer is `FullHeight * TileWidth + 1` bytes,
  `'-'` means "this grid cell never animates", `+0x38` counts the rest, and the draw
  indexes it by `k*TileWidth + c` (`REG-STR-082`, `SPR256-STR-041`). `AnimTime`/
  `AnimFrame`/`AnimMask` are read **only when `Phases > 1`**; the timeline is the same
  run-length expansion `objects.reg` uses. Corpus: `AnimMask` length `= TileWidth ×
  FullHeight` on 14/14 classes that spell one, alphabet `{'-', '*'}`; ten further classes
  spell `Phases > 1` and no animation keys at all, so `Phases > 1` does not mean animated.

  `FullHeight` = **1..6**
  (not `0..20` — that range is `SelectionX1`'s own corrected value, the old `FullHeight`
  column having read the next key's true value); `TileHeight` = **1..6** (held);
  `TileWidth` = **1..11** (wider than the old bundled `1..6`, which was in fact
  `TileHeight`'s range). Also `AnimMask`, `Phases`, `Flat`, `Indestructible`, `Usable`,
  `VariableSize`, `LightPulse`, `LightRadius`, `ShadowY`, `SelectionX1..Y2`, `IconID`,
  `Picture` (`REG-STR-021`, `REG-STR-040`) — `class+0x54`, `0x10` bytes, the same `graphics\infowindow\<name>.bmp` leaf as `units.reg`'s `InfoPicture` and never carrying a tier (`UNIT-PICT-037`). Thirteen of the 66 sections name `magic`, `ruins` or `sphinx`, which the RU root does not ship (`REG-PICT-083`).
- **`projectiles.reg`** — 31 projectiles (held, `REG-PROJ-041`). `SFX`'s published
  `12..62` is **withdrawn, not re-attributable** — on the corrected base `SFX` is present
  on 23/31 and always reads `1`, a presence flag rather than a variable id; no single
  re-measured neighbour reproduces `12..62` closely enough to name. `Width`/`Height`
  newly measured at `12..128`. `File` is a kind-0 **string** sprite path (e.g.
  `archer\arrow`), not a `Files[]` index like `units`/`objects`. Also `Phases`,
  `RotationPhases`, `Homing`, `Flip`, `Palette` (`REG-PROJ-022`, `REG-PROJ-041`).
- **`material.reg`** — 16 `Material*` each with a `Path`, no `[Global]` counter — held
  unchanged (`REG-MAT-042`).

## The other registries (`REG-CUT-053`, `REG-SFX-057`, `REG-NPC-058`, `REG-SCN-059`, `REG-AI-060`)

The promoted per-key tables and their corpus bounds are recorded in the claims cited above.
`data/map.reg` is the one registry in the install still without a published table.

- **`sfx.res::sfx.reg`** — the **sound-slot table** (`REG-SFX-057`). `[Global] SfxCount = 564`
  is the **highest slot id, not the entry count**: `[Sfx]` holds **115** kind-0 strings named
  `Sfx<n>` over a sparse id space whose maximum is exactly 564. Each value is a path relative
  to `sfx.res`, no extension (`click00`, `units\sword`, `ambient\river`, `magic\firewall`).
  Every **nonzero** element of a `units.reg` class's `Sound[]` array is an existing slot
  (**153/153**; 17 further elements read `0`), and `REG-OBJ-047`'s ambient slots resolve —
  `0x3c..0x3e` = `ambient\bird1/2/3`, `0x46` = `ambient\crow`.
  The loader and consumers are pinned by `VIDEO-SFX-013`…`VIDEO-SFX-021`.
  `FUN_004af540` allocates `SfxCount+1` pointers,
  clears them and constructs present keys `Sfx1` through `Sfx<SfxCount>`; the data value is
  capacity, not a live-event count. Both roots' 115 rows join to the archive 112 times and miss
  three paths: slot 62 `ambient\bird3`, slot 80 `ambient\wind1` and slot 81 `ambient\wind2`.
  Every non-zero inherited `units.reg Sound[]` value still resolves. Runtime selectors are
  separate: fixed slots, inherited elements and `500+value` formulas read this pointer array,
  while literal-path samples and conditional voice banks bypass it (`VIDEO-SFX-013`..`021`).
- **`scenario.res::npc.reg`** (`REG-NPC-058`) — **105** `npc<n>` sections (ids 1..132, sparse),
  four archetype blocks (`MaleFighter`, `MaleMage`, `FemaleFighter`, `FemaleMage`, each
  `Face Body Reaction Mind Spirit Skill`), and `Multiplayer` with four kind-6 face lists
  (`FacesMM` 6, `FacesMF` 15, `FacesFM` 4, `FacesFF` 8). Per-npc, out of 105: `Flags` 105 (kind
  0), `Face` 86 (1..30), `Picture` 33 (64..80), `PortraitX1` 48 (7..69), `PortraitY1` 48 (8..72),
  `PortraitX2` 10 (48..54), `PortraitY2` 10 (15..32), `DataBinID` 23, `PriceA` 13 (28..80),
  `PriceB` 13 (0..15). **`Flags` is a comma-separated token list with `!`
  negation** over exactly 11 tokens — `Human, Mage, Female, Face, Picture, Platoon, Hero, Me,
  Start, MySex, MyClass` — in 15 distinct combinations, all 11 present as whole strings in
  `rom.exe`. **`DataBinID` is a `Templates.ini` id** (`EDITOR-023`): 23/23 inside the declared
  `1..1088`, and 13/13 of those at or above the table's lowest labelled slot name a slot
  (`509 → M10_Witch`, `1051 → F_Knight3`, …); its one located reader writes it into a message
  body at `+0xa` behind the byte tag `0x49` (`REG-NPC-090`).

  The in-memory record is **`0x30`** bytes, allocated only when `Flags` is a non-empty string,
  and held in an array indexed by the section number itself (`npc<i>` → element *i*, 256 slots):

  | Offset | Key | Default |
  |-------:|-----|--------:|
  | `+0x08` | `Flags` (CString) | `""` |
  | `+0x0c` | `Face` | 0 |
  | `+0x10` | `Picture` | 0 |
  | `+0x14` | `DataBinID` | 0 |
  | `+0x18` | `PriceA` | 0 |
  | `+0x1c` | `PriceB` | 0 |
  | `+0x20` | `RECT { PortraitX1, PortraitY1, PortraitX2, PortraitY2 }` | -1 each |

  **`Face` and `Picture` name no picture resource** (`REG-NPC-088`). They are copied into a
  dialogue actor's **typeID** (`+0x20`) and **face** (`+0x24`), each gated on the matching token
  of the same section's `Flags` list, and the portrait that follows is the ordinary
  `graphics\infowindow\<InfoPicture><face>.bmp` leaf (`UNIT-PICT-036`) with the digit dropped at
  face 1. With the `Start` token the face comes instead from one of the four archetype sections,
  chosen by the player's own sex/class bits. Where `Picture` is present `Face` is 1..4, a tier;
  where absent it is 1..30 and the section takes the composed-doll path (`UNIT-PICT-035`).

  **`Portrait*` is a Win32 `RECT`** — the loader hands the four consecutive values to `CopyRect`,
  which fixes the order as `left, top, right, bottom` — and only `left` and `top` are read
  (`REG-NPC-089`). They cut a fixed **72 x 96** source window, `(X1, 144-Y1)`…`(X1+72, 240-Y1)`,
  out of a 160 x 240 16-bit canvas; `PortraitY1` absent (the loader's `-1`) selects a fixed
  72 x 92 default at `(36, 140)`. `PortraitX2`/`PortraitY2` have no reader in the image.

  `npc21`..`npc24` are a four-slot block the image knows by literal number — two routines gate on
  `21 <= n <= 24` and a third subscripts `array + id*4 + 0x54` — and they are exactly the four
  sections carrying `Start` (`REG-NPC-090`).
- **`scenario.res::scenario.reg`** (`REG-SCN-059`) — `[General] TotalMissions = 15`,
  `ScenarioMissionCount = 15`, `MercenaryCount` a 15-element array; **24** `[Mission<n>]`
  sections (15 main, numbered by tens, plus 9 sub-missions). Keys, with presence out of 24:
  `Mercenaries` 15, `InnNPC` 13, `InnMission` 13, `ShopMinPrice` 13 (0..5000), `ShopMaxPrice`
  13 (1000..10 000 000), `Payment` 9 (700..700 000), `EnableMercenary` 5, `ShopMission` 4,
  `TCMission` 3, `AddTextDocument` 2, `AddHero`/`AddPictureDocument`/`AutoGetMission`/
  `LastMission` 1 each. Five of these are the per-record-kind keys above. **Every npc-shaped
  reference resolves to an existing `npc<n>` section: 143/143.** What `MercenaryCount` counts
  is **Unknown** (roster length fails 12/15); `TotalMissions` is named by no literal in either
  executable.
- **`scenario.res::globalmap.reg`** (`REG-SCN-059`) — `[General] ObjectCount = 35` matching
  **35** `[MapObject<n>]` sections (1-based), each `MapPoint` = 2 ints, `MapRect` = 4 ints, 5
  with a `Picture` string; plus `[MissionObjects]`, **28** `Mission<n>` keys naming a map-object
  number. **28/28 name an existing `MapObject<n>`**; 24/28 name a `scenario.reg` mission and
  24/24 of those missions appear here (extras: 81, 101, 141, 151). The two geometry keys have
  different consumers (`TOWN-116`, `TOWN-118`): `MapRect` becomes an indexed-query hit rectangle;
  `MapPoint` is a world-route endpoint and marker anchor. All 35 shipped MapPoints land on
  index-2 graph nodes in `graphics.res::global.map/pathmap.bmp`; that bitmap has 46 further
  nodes with no MapObject. A MapObject whose `Picture` is `"nothing"` receives a static return-enable
  bit, while a picture-bearing object's indexed return is gated by the runtime marker cache (`TOWN-123`).
  `MapPoint` coordinates and PathMap pixels are therefore one coupled data seam: moving one
  without placing an index-2 node at the same coordinate leaves route lookup without an endpoint.
- **`world.res::data/ai.reg`** (`REG-AI-060`) — 4 records, empty pool: `[Scanning]
  MinimalGuardRange = 8` and `[Tasker] IntelligentCons = 15`. Neither `Tasker` nor any run
  beginning `IntelligentCons` occurs in `rom.exe` or `Map Editor.exe`.
- **The 33 video-cutscene registries** (`REG-CUT-053`) — **one schema with two optional section
  families**, in three shipped shapes:

  ```
  Common { startx, starty, nFadings [, nPanaramings] }
  Fading<n>     { startframe, endframe, startfade, endfade }   (kind 4 doubles)
  Panaraming<n> { startframe, endframe, stepx, stepy }
  ```

  27 files carry `Common` + `Fading`, 4 add `Panaraming` (`M10/01`, `M50/01` in both
  containers), 2 carry `Common` alone with `nFadings = 0`. `nFadings` equals the `Fading`
  section count **33/33** and `nPanaramings` the `Panaraming` count **4/4**;
  `startx == starty == 0` on 33/33; `startframe <= endframe` on 65/65; `startfade`/`endfade`
  take only `0.0` and `1.0` (61 each — the whole kind-4 population, `REG-DBL-035`), the first
  `Fading` always fading in and the last always out. Each registry sits beside a `.smk` node of
  its own basename **33/33**, and the 15 node paths present in both `VIDEO4.RES` and
  `VIDEO8.RES` are identical in resolved content **15/15** while their videos differ in size.
  What the frame numbers index is **Unknown** — it needs the `.smk` frame count.

## Map Editor text catalogs (`EDITOR-023`)

The install ships parallel game-authored vocabulary as plain text at its root:
`Description Checks.ini` (22 trigger **condition** opcodes, id→name→params),
`Description Instants.ini` (44 trigger **action** opcodes), `Templates.ini`
(1088-slot placeable-template id→name table). Trigger types: `Target_Unit`,
`Target_Group`, `Target_Player`, `Target_Item`, `Target_Structure`,
`Target_Building`, `X`, `Y`, `int`, `Enum`, `Const`.

## `projectiles.reg` — and why its `ID` is an address

`REG-PROJ-086`, read from the loader `FUN_0046e0f0` rather than fitted.

```
[Global] Count = 31
[Projectile<n>]   n = 0 .. Count-1        (the section index is NOT the id)
  File            string, e.g. archer\arrow  -- a path, as in structures.reg, not a Files[] index
  A16             1 -> the sheet is .16a, 0 -> .256          default 0
  ID              the key everything else uses                default -1
  Phases          frames per facing                           default -1
  RotationPhases  facings in the sheet                        default 16
  Width, Height   the DRAW's centring offsets, not the art    default 64, 64
  Palette         0 -> use the shared projectiles.pal          default 0
  Homing          stored at +0x2c and read by nothing          default 0
  Flip            halves the sheet: 9 stored facings, 7 mirrored   default 0
  SFX             1 wherever present on the shipped file       default 0
```

The loader grows its array to `ID + 1` and stores the record at `array[ID]`, so **`ID` is the
address**: the shipped 1..62 domain over 31 rows leaves 32 null slots, and a null slot makes every
consumer draw nothing rather than fault. Two rows may name one `File` (the shipped file does, twice
for `catap2\sprites`). Two `projectiles/*` directories — `smoke0` and `smoke1` — are loaded by
literal path outside the loop and are in no section.

The id space is shared with `units.reg`, whose `Projectile` key selects a unit's shot from the same
array; 9 of the 34 shipped classes carry one, and two of those point at a spell's art.

## `AddHero` town activation

A mission record stores `AddHero` as a `u16` array. Town activation walks the current record's array, sends one command `0x49` per value, and then clears the array. The array is not a mission-completion action.

The shipped value `[Mission30] AddHero = 22` selects `npc22`. Its `Mage,!MySex` flags select Humans row 28 (`PC_Fergard`) when the existing primary character is female and row 29 (`PC_Reniesta`) when the existing primary character is male. The displayed name is the corresponding localized entry in `text/npcnames.txt`. The special `npc21..24` branch replaces the stored `DataBinID` with `26 + selector`.
