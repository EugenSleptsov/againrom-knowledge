# Encoding and lookup

[Reference](format.md)

## Structure

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

The header contains six u32 fields. Its first four fields form a virtual
root node without a stored name. Opening the registry supplies that name from
the file path. — REG-FMT-031

| Off | Field | Notes |
|-----|-------|-------|
| 0x00 | magic | `&YA1` = `0x31415926` LE; mismatch → the game's `"bad signature"` |
| 0x04 | root `value` | index of the root's first child record |
| 0x08 | root `size` | number of top-level children |
| 0x0C | root `kind` | **17** = `0x11` in installed files = subkey (bit 0) \| sorted (bit 4). The raw loader does not require this value (`REG-099`, `REG-101`). |
| 0x10 | `R` | total record count |
| 0x14 | — | read into the object but used by no accessor found; **Unknown** |

Then `R × 32 B` of records at `0x18`, a `u32` **pool byte length**, and the pool:

```
poolStart = 0x18 + R*32 + 4
registryEnd = 0x18 + R*32 + 4 + poolLen
```

## Record (32 bytes)

The record layout in REG-REC-032 is retained; its lookup clause is partially
retracted, as narrowed by REG-100.

| Off | Type | Field | Notes |
|-----|------|-------|-------|
| 0x00 | u32 | — | explicitly zeroed by the writer; read by no accessor found; **Unknown** |
| 0x04 | u32 | `value` | int32 · pool byte offset · child-block **start** index · **low dword of a double** |
| 0x08 | u32 | `size` | byte length · child count · **high dword of a double** |
| 0x0C | u32 | `kind` | bitfield, see below |
| 0x10 | char[16] | `name` | 15 significant characters + NUL |

Record `i` starts at `0x18 + 32*i`; its value is at `0x1c + 32*i`. These
REG-REC-032 offsets are retained; its lookup clause is narrowed by REG-100.

### Kind — a bitfield, not an enum (`REG-KIND-033`)

```
type       = kind & 0x0E          tested as AND 0xe / CMP by every accessor
bit 0      = node is a subkey
bit 4      = children are sorted -> the lookup bsearch()es instead of scanning
bit 28     = the name was longer than 15 chars and got truncated (no shipped record sets it)
```

### Kinds and their storage (`REG-KIND-034`)

| Kind | Meaning | Storage |
|---:|---|---|
| 0 | string | `poolBase + value`, `size` bytes **including** the NUL |
| 1 | subkey | children = records `[value, value + size)` |
| 2 | int32 (signed) | `value` |
| **4** | **double** | **`value` = low dword, `size` = high dword — 8 bytes in the record, no pool access** |
| 6 | int32[] | `poolBase + value`, `size/4` LE int32s |
| 10 | double[] | `poolBase + value`, `size/8` LE doubles |

A key name does not determine its kind. `Mercenaries`, `InnNPC`, `InnMission`,
`EnableMercenary` and `AddTextDocument` can hold a bare int32 or an int32 array.
The content model treats a scalar as a one-element list. The behavior of
`FUN_004cd240` itself on kind 2 remains Unknown. Empty-string/array alternatives
also occur. — REG-KIND-056, REG-KEY-045

For `(kind>>1)&7 == 4` (kind 8/9), the string converter skips the first
four pool bytes, copies `size−4` bytes and terminates with two NULs. Its
producer and semantic meaning remain Unknown. Other conversion classes
return `"UNKNOWN TYPE. CANT CONVERT"`.

The game's own names for the types, from its error strings: *int*, *double*, *int array*,
*double array*, *string array or single string*.

### Kind 4 — the double (`REG-DBL-035`)

```
reader  FUN_004ccda0 :  AND EDX,0xe ; CMP DL,0x4 ; FLD qword ptr [node + 0x4]
writer  FUN_004cb340 :  *(double *)(node + 4) = atof(text) ; node->kind = 4
```

Kind 4 stores the binary64 value directly at record `+4`; it has no pool allocation. — REG-DBL-035

### Pool encoding (`REG-VAL-025`, as amended)

Kind 0/6/10 records address `poolBase + value`; `size` is a byte length.
Values do not depend on a cursor or the preceding record. — REG-VAL-025

### Tree walk (`REG-VAL-028`, as amended)

A subkey addresses the contiguous child range `[value,value+size)` in the
record array. The root has the same relation. — REG-VAL-028

## Text (`REG-TEXT-036`, `REG-TEXT-037`)

The loader transfers the record array and pool verbatim. The string accessor
copies from `poolBase + value`; it performs no code-page conversion.
— REG-TEXT-036, REG-TEXT-037

The unsorted name lookup's `_strnicmp` fast path folds `A`–`Z` only and leaves
bytes `>= 0x80` untouched. The sorted lookup instead compares raw bytes with no
case folding (`REG-100`). The CRT locale-dependent arm remains outside that
conditional fast-path statement.

Registry text has no established non-ASCII display mapping. Keep its bytes
distinct from any display encoding chosen by a consuming application.
— REG-TEXT-036, REG-TEXT-037

## Name matching, ordering and the 15-character clamp

The parent kind's bit 4 selects two different comparisons. With bit 4 clear,
`004ce8e0` linearly scans children using `_strnicmp(query, child+0x10, 15)`.
With bit 4 set and a nonempty list, it copies at most 15 query bytes into a
NUL-terminated temporary record and calls `bsearch` with `004ceaf0`. That
comparator compares unsigned bytes to NUL, **case-sensitively**. A root or leaf
case variant can therefore match in an unsorted container and miss in a sorted
one. Nonzero CRT locale state is not covered by the unsorted ASCII-fold result.
The lookup clause of REG-REC-032 and the comparator clause of REG-KEY-054
are narrowed by this result. — REG-100

The node-insert helper clamps a new name to 15 bytes plus NUL and marks an
overlong input with kind bit 28; typed setters may later replace the kind.
The raw loader does not clamp names. A 15-byte stored name
matches a longer query sharing that prefix on both lookup routes; a raw name
occupying all 16 bytes without NUL matches that query on the 15-byte linear
route but misses the sorted route's truncated temporary key. Arbitrary
unterminated names and comparison reads outside the record remain outside the
valid producer-name domain. — REG-099, REG-100

The installed registry roots set bit 4; their ordinary section children do
not. Section names therefore use sorted lookup, and per-key names ordinarily
use linear lookup. REG-KEY-054's unconditional comparator/writer clauses are
partially retracted; use REG-100's comparator rules.

- Sections are stored in **lexicographic**, not numeric, order — `MapObject1, MapObject10,
  MapObject11, …, MapObject2`. The order in which a loader's `"<Prefix>%d"` counter visits
  sections is **not** the record order.
- A container advertising bit 4 must order its children for the case-sensitive
  comparator. A container with bit 4 clear takes the linear route. — REG-100

Neither lookup route checks a child's kind before matching its name. The linear
route selects the first equal child; sorted duplicate blocks of 2, 3 and 4
equal names select indices 0, 1 and 1 in the original bsearch body. Lookup does
not enforce the parent subkey bit or skip a child carrying bit 30. The selected
integer getter checks only `kind & 0x0e == 2` after two name lookups; missing
names return its caller's default, while a mismatched value type reaches the
exception-throw boundary. These are accessor facts, not structural or
application acceptance rules. — REG-101

Some stored keys are 15-byte prefixes of longer authoring names:
`MinimalGuardRan` → `MinimalGuardRange`, `AddPictureDocum` →
`AddPictureDocument`, `ScenarioMission` → `ScenarioMissionCount`.
The stored prefix is the lookup name. — REG-NAME-055

## Raw parsing and writing of unrecognized entries

The raw loader reads six header dwords, the full `R*32` record block, a pool
length and the pool. It checks the signature and allocation results but does
not enumerate names, validate child kinds, enforce ordering or deduplicate
records. Unrecognized entries remain in the raw registry. Known lookups remain unchanged
when entries do not collide under the selected comparator and the advertised
order remains true. This does not establish application LOAD acceptance.
— REG-099, REG-100, REG-101

The raw writer sorts through `004ce660`, then writes the supplied registry's
header, complete record array and pool. Unrecognized entries can survive this direct round trip. Sorting can alter lookup behavior: a lowercase
root matched the unsorted case-folding route before writing and missed the
case-sensitive route after writing set the root's sorted bit. The sort helper's
recursive candidate is indexed by `parent.value + parent.size`, rather than
its loop counter; it does not establish a conventional recursive sort of every
descendant. Internal insertion/copy helpers and an application's fresh-registry
producer are separate paths. — REG-102

## Read and write sequence

1. Read six little-endian header dwords; check magic `0x31415926`.
2. Read `R` records of 32 bytes from `0x18`.
3. Read u32 `poolLen`, then exactly that many pool bytes.
4. Interpret each value by its own kind. A subkey refers to records; string
   and array kinds refer to pool bytes; an integer or double is inline.
5. Resolve names with the parent kind's bit-4 comparator. Check record and
   pool ranges before dereferencing them in a bounded decoder.

To emit a registry, assign contiguous child ranges, encode each typed value,
build the pool and write the header, records, u32 pool length and pool.
Set each sorted flag only when its child block has the required unsigned-byte,
case-sensitive ordering. The original raw writer's sorting behavior is stated
above; it is not a general recursive normalization guarantee. Preserve unknown
values and the header's unnamed final dword when rewriting an existing store.
— REG-FMT-031, REG-KIND-034, REG-099, REG-100, REG-102
