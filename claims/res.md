# RES — the RES / LM container (`&YA1`)

Claims about the `&YA1` container that `.res` and `.lm` files use: its header
and node records, the `rom.exe` code that opens, reads and resolves it, the EN
and RU shipped corpora, and the rule that decides which archive answers a path.
EXP-0017, EXP-0032 and EXP-0034 read the reader, the lookup, the writer and the
node-insert. Between them they cover every field this ledger names, because the
header is a node record and the `.reg` attestations carry the `.res` layout too.
Spec: [`formats/res/format.md`](../formats/res/format.md). Format of this file:
[registry.md](registry.md).

## Header, node record and tree

| ID | Claim | Confidence | Status | Evidence |
|---|---|---|---|---|
| RES-MAGIC-001 | The magic is `26 59 41 31` ("&YA1"); as a little-endian u32 it is `0x31415926`, π's first 8 digits, which is arithmetic only and claims no intent. | High | ✔ promoted | [EXP-0001](../experiments/EXP-0001-res-container/) |
| RES-HDR-002 | The header is 24 bytes, and payload data begins at offset 24. | High | ✔ promoted | [EXP-0001](../experiments/EXP-0001-res-container/) |
| RES-HDR-003 | `u32@16` is the registry offset; a registry `[@16,EOF)` length ≡ 0 (mod 32) is an EN-corpus packer fact, not a format rule. | High | ✔ promoted (amended) | [EXP-0001](../experiments/EXP-0001-res-container/), [EXP-0051](../experiments/EXP-0051-res-ru/) |
| RES-HDR-004 | `u32@20` is the node count, the engine's only size input (allocation and bulk-read length); `count = regLen/32` holds on the EN corpus only. | High | ✔ promoted (amended) | [EXP-0001](../experiments/EXP-0001-res-container/), [EXP-0051](../experiments/EXP-0051-res-ru/) |
| RES-HDR-005 | `u32@8` is the root count: the number of top-level, unowned nodes. | High | ✔ promoted | [EXP-0001](../experiments/EXP-0001-res-container/) |
| RES-HDR-006 | Header `u32@4` and `u32@12` are the root directory node's `off` and `type` fields (`RES-HDR-017`, `RES-HDR-018`). | High | ● active | [EXP-0001](../experiments/EXP-0001-res-container/), [EXP-0014](../experiments/EXP-0014-res-opaque/), [EXP-0017](../experiments/EXP-0017-res-open/) |
| RES-NODE-007 | A node is 32 bytes, `{u32, u32 off, u32 size, u32 type, char[16] name}`, and shipped names are 0xCD-padded. | High | ✔ promoted | [EXP-0001](../experiments/EXP-0001-res-container/) |
| RES-NODE-008 | Node `type` 0 is a file, a byte range in `[24,regOff)`; `type` 1 is a directory, a first-child index plus a count. | High | ✔ promoted | [EXP-0001](../experiments/EXP-0001-res-container/) |
| RES-TREE-009 | Corpus-wide, the nodes form a tree: all reachable, no cycles, and file ranges tile `[24,regOff)` exactly (12 standalone containers / 4592 nodes, 0 violations). | High | ✔ promoted (amended) | [EXP-0001](../experiments/EXP-0001-res-container/) |
| RES-SCOPE-010 | `.LM` files use the same container, and an empty archive is valid; `Allods/*.RES` use an identical container. | High | ✔ promoted | [EXP-0001](../experiments/EXP-0001-res-container/) |
| RES-NODE-011 | Node `u32@0` is the node record's reserved word: 0 in standalone archives, and none of the four scanned lookup, path-walk and finalize-sort functions dereferences it. | High | ● active (amended) | [EXP-0001](../experiments/EXP-0001-res-container/), [EXP-0014](../experiments/EXP-0014-res-opaque/), [EXP-0017](../experiments/EXP-0017-res-open/) |
| RES-HDR-012 | On the 60-blob EN widened corpus, header `@0x04` is 0 or the non-root-node count; it is the root node's `off`, and the two-value law is EN-only. | High | ● active (amended, superseded) | [EXP-0014](../experiments/EXP-0014-res-opaque/), [EXP-0017](../experiments/EXP-0017-res-open/) |
| RES-HDR-013 | Header `@0x0C` is a per-file constant in `{1, 17}` with no third value in 60 blobs; it is the root node's type/flags word, and bit 4 marks sorted children. | High | ● active (amended, superseded) | [EXP-0014](../experiments/EXP-0014-res-opaque/), [EXP-0017](../experiments/EXP-0017-res-open/) |
| RES-NODE-014 | Tail-registry node `@0x00` is 0 for every node of every reachable tail-registry `&YA1` (4592 nodes, 11 archives): a reserved word, not a hash or id. | High | ● active (amended) | [EXP-0014](../experiments/EXP-0014-res-opaque/), [EXP-0017](../experiments/EXP-0017-res-open/) |
| RES-SCOPE-015 | The `&YA1` magic labels two distinct formats, the tail-registry archive and the inline REG-style record store; in this install's 60 blobs every nested `&YA1` is inline. | High / Medium | ● active | [EXP-0014](../experiments/EXP-0014-res-opaque/) |
| RES-NODE-016 | In `rom.exe`, the 24-byte `&YA1` header is a directory-node record: the tree's virtual root node, sharing its first 16 bytes with a node. | High | ● active | [EXP-0017](../experiments/EXP-0017-res-open/) |
| RES-HDR-017 | In `rom.exe`, header `@0x04` is the root directory node's `off`, the array index of the first top-level node, and the name lookup dereferences it. | High | ● active | [EXP-0017](../experiments/EXP-0017-res-open/) |
| RES-HDR-018 | In `rom.exe`, header `@0x0C` is the root directory node's type/flags word: `1` is a directory, and bit 4 (`0x10`) is the children-sorted flag the lookup branches on. | High | ● active (superseded) | [EXP-0017](../experiments/EXP-0017-res-open/) |
| RES-NODE-019 | In `rom.exe`, node `@0x00` is the node record's reserved word, and none of the four scanned lookup, path-walk and finalize-sort functions dereferences it. | High | ● active (partially retracted) | [EXP-0017](../experiments/EXP-0017-res-open/) |

### RES-MAGIC-001

- The bytes are `26 59 41 31` ("&YA1"); the LE u32 is `0x31415926`, π's first
  8 digits. The π reading is arithmetic and claims nothing.
- `RES-CODE-020` finds the immediate `0x31415926` in exactly four `rom.exe`
  functions: two `CMP` readers and two `MOV` writers.

**Confidence.** High. This is the case the scale's carve-out is written for: a
magic constant has no second model, and here it is not even corpus-only. The one
alternative worth naming, "these four bytes are data rather than a signature",
is what the two `CMP`s exclude.

### RES-HDR-002

**Confidence.** High, and not from a fit. `RES-NODE-016` shows the header is a
node record, read field by field into the object's node slots by
`FUN_004c90a0`, and the tail open takes the payload base from it.

### RES-HDR-003

- `u32@16` is the registry offset: the tail open seeks by this word.
- The registry `[@16,EOF)` length ≡ 0 (mod 32) holds on the EN corpus only. RU
  `MAIN.RES` ships `(EOF−@16) = 16151 ≡ 23`, and the engine never computes that
  length at all (`RES-OPEN-026`).
- The node array is still 32-stride (`RES-GEOM-028`).

**Confidence.** High for the offset semantics: the tail open seeks by this
word, and the 32-byte stride is pinned by the `SHL 0x5` in the shared lookup
(`REG-REC-032`). The region-length ≡ 0 clause is EN corpus only.

**Amended.** EXP-0051 narrowed the second clause: the mod-32 identity was
worded as a format rule and is an EN-corpus packer fact
([`retracted.md`](retracted.md)). The offset semantics stand.

### RES-HDR-004

- The tail open reads `u32@20` as the allocation and bulk-read length; it is
  the engine's only size input.
- `= regLen/32` holds on EN 12/12 and fails on RU `MAIN.RES`
  (`504 ≠ 16151/32`). A reader deriving the count from region length rejects a
  file the engine opens (`RES-OPEN-026`, `RES-ACCEPT-031`).

**Confidence.** High for the count semantics, read as the allocation and read
length by the tail open. The `regLen/32` identity is EN corpus only.

**Amended.** EXP-0051 narrowed `= regLen/32`: it was worded as a format
identity and is an EN-corpus fact ([`retracted.md`](retracted.md)). The count
semantics stand.

### RES-HDR-005

**Confidence.** High. It is the root node's own `size` field, the child count a
directory node carries, which `FUN_004ce8e0` iterates as the loop bound
(`RES-NODE-016`/`REG-REC-032`).

### RES-HDR-006

- EXP-0017 resolved the meaning of both words: they are the root directory
  node's `off` and `type` fields (`RES-HDR-017` / `RES-HDR-018`).
- Their value laws are pinned on the widened corpus in `RES-HDR-012` (`@4`) and
  `RES-HDR-013` (`@12`).

**Confidence.** High, the grade of `RES-HDR-017` and `RES-HDR-018`, whose
instruction-level readings resolve the meaning. The value laws carry their own
grades in `RES-HDR-012` and `RES-HDR-013`.

### RES-NODE-007

- Layout `{u32, u32 off, u32 size, u32 type, char[16] name}`; the name is
  0xCD-padded.
- The `0xCD` padding is a measurement. It is the CRT's uninitialised-fill byte:
  a fact about the packer, not a format constant.

**Confidence.** High. `REG-REC-032` attests the same record four ways in
`rom.exe`: the loader `FUN_004cae80`, the lookup `FUN_004ce8e0`, the accessors
`FUN_004ccc20`/`FUN_004cc670` and the node-insert writer `FUN_004ce6d0`. The
stride is the `SHL 0x5` in the loader, the lookup and the node-insert, and
`name @+0x10` is fixed by the writer's 15-character clamp.

### RES-NODE-008

**Confidence.** High. `RES-HDR-018` reads the flag word at instruction level
(`1` = directory, bit 4 = sorted), and `REG-REC-032` shows the directory case
computes `&records[node->off]` and iterates `node->size`, which is what
"first-child index + count" asserts.

### RES-TREE-009

- The measurement is exhaustive over the 12 standalone containers (11 tail
  archives and the empty `KIDS.LM`): 4592 nodes, 0 violations.
- Exact tiling is the invariant EXP-0030 showed can hold under two framings at
  once, so it corroborates the layout rather than pinning it. The reader pins it
  (`RES-NODE-016`, `RES-HDR-017/018`). This claim covers only the measurement.

**Confidence.** High as a measurement.

**Amended.** The population was worded "60 blobs / 4592 nodes". EXP-0001's
`evidence/archives.csv` measures 12 containers whose node counts sum to 4592,
with every violation column 0. The 60 blobs are EXP-0014's widened reach, which
adds 48 inline stores that neither EXP-0001 nor EXP-0014 tree-checked
([`retracted.md`](retracted.md)).

### RES-SCOPE-010

**Confidence.** High. A measurement over every such file in this install. The
extension plays no part in the binary either: one reader validates the magic
(`RES-CODE-020`).

### RES-NODE-011

- Observed 0 in standalone archives.
- `RES-NODE-014` widened the observation and refuted the hash/id reading.
- EXP-0017 resolved the meaning: the node record's reserved word. No
  instruction in the four scanned functions, the lookup `FUN_004ce8e0`, the
  path-walk `FUN_004ce800`/`FUN_004ce9e0` and the finalize-sort `FUN_004ce660`,
  dereferences it (`RES-NODE-019`).

**Confidence.** High, the grade of `RES-NODE-014` (always zero) and of
`RES-NODE-019`, which resolve the meaning. The absence of a `[node+0]`
dereference is High over the four scanned lookup, path-walk and finalize-sort
functions only, the scope `RES-NODE-019` keeps.

**Unknown.** Whether a node-touching function outside the four scanned reads
`[node+0]`. `RES-NODE-019` names such functions (`FUN_004c9320`, `FUN_004ce8c0`,
`FUN_004c9ad0`, `FUN_004ccc20`/`FUN_004cc670`, and the node-insert
`FUN_004ce6d0`, which stores `0` there).

**Amended.** The former headline and Confidence said the word is never
dereferenced by lookup, descent or sort, and took `RES-NODE-019`'s High for that
without its scope. `RES-NODE-019`'s "the four functions that touch a node"
clause is withdrawn ([`retracted.md`](retracted.md)), so this claim now carries
the same scope and the same Unknown.

### RES-HDR-012

- Header `@0x04` ∈ `{0, non-root-node count}` over the 60-blob widened corpus:
  12 standalone + 44 nested `.reg` + 4 save state stores.
- `@0x04 = 0` in 57 blobs: all 48 inline REG stores, the empty archive and 8/11
  tail archives.
- `@0x04 = nNodes−roots` (Σ dir child-counts) in exactly
  graphics(3120)/main(491)/sfx(274).
- Not recomputed: VIDEO4/VIDEO8/movies/speech/world have 30–359 non-root nodes
  yet store 0.
- Not a checksum: `blob_sum32` is unrelated.
- The populate-vs-zero choice has no structural predictor: depth, dirs, nesting
  and contains-REG are all falsified.

**Confidence.** High for the value law on the EN corpus. The meaning carries
`RES-HDR-017`'s grade (EXP-0017).

**Amended.** `RES-HDR-017` (EXP-0017) resolves the meaning: `@0x04` is the root
directory node's `off`, the index of the first top-level node, and
`{0, ownedNodes}` is the packer's roots-first vs roots-last layout order, not a
recomputed count. It supersedes the pointer "Which archives populate it →
`rom.exe` writer routine": the tail archives come from an external packer, and
`rom.exe` contains only the inline writer (EXP-0017, `RES-GEOM-028`,
[`retracted.md`](retracted.md)). `RES-HDR-029` (EXP-0051) narrows the scope: the
`{0, nNodes−roots}` law is EN-corpus only, because RU `SFX.RES` stores a
mid-table `186`.

### RES-HDR-013

- `17` in every inline REG-style store (48/48) + `graphics.res` + `main.res`;
  `1` in the other 9 tail archives + the empty archive.
- `17 ⟺ contains-REG` is falsified: main `@17` has 0 nested, VIDEO4 `@1` has
  18.
- `17 ⟺ inline-flavor` is falsified: graphics and main are tail archives.
- `17 ⟺ has-dirs` is falsified: VIDEO4 has 15 dirs `@1`.
- Bit 4 (`0x10`) is the discriminator: the children-sorted flag
  (`RES-HDR-018`).

**Confidence.** High for the domain and the inline⟹17 invariant. The meaning
carries `RES-HDR-018`'s grade (EXP-0017).

**Amended.** `RES-HDR-018` (EXP-0017) resolves the meaning: `@0x0C` is the root
directory node's `type/flags` (`1` = directory), and bit 4 is the
children-sorted flag the reader tests to pick binary vs linear search
(`17` = dir+sorted), set by the writer's finalize-sort. It supersedes the
EXP-0014 reading of bit 4 as "a format-variant / writer-generation flag"
([`retracted.md`](retracted.md)). `RES-HDR-030` (EXP-0051) narrows the
association: RU ships `1` on all 11 archives including graphics and main, so
"17 on graphics+main" was EN packing. The domain `{1,17}` holds over all 23.

### RES-NODE-014

- Exactly 0 in 4592 nodes over 11 archives. The per-node hash/id/checksum
  candidate is refuted: a hash would vary.
- Reserved and always zero.
- The inline REG record `@0x00` is the separate `A` field (`REG-FMT-017`), not
  this word.
- Narrows `RES-NODE-011`.

**Confidence.** High for always-zero. The meaning carries `RES-NODE-019`'s
grade (EXP-0017).

**Amended.** `RES-NODE-019` (EXP-0017) resolves the meaning: the node record's
reserved word, aligned with the header's magic. None of the four scanned lookup,
path-walk and finalize-sort functions dereferences it; whether another
node-touching function reads it is Unknown.

### RES-SCOPE-015

- Tail-registry archive: `u32@0x10` = registry byte offset (EXP-0001).
- Inline REG-style record store: `u32@0x10` = record count `R`
  (EXP-0006/0011).
- Across the 60 reachable blobs, every nested `&YA1` (inside an archive or a
  save) is the inline flavor; the tail-registry flavor occurs only as a
  standalone top-level `.res`/`.LM`.

**Confidence.** High that there are two flavours: `rom.exe` has two distinct
validators for the same magic, the tail open `FUN_004c90a0` and the inline read
`FUN_004cae80` (`RES-CODE-020`), so the distinction is in the code, not only in
the bytes. Medium for the distribution: it is corpus-only over this install's 60
blobs, and nothing read in the binary forbids a nested tail-registry archive. It
is a fact about the shipped data and must not be relied on as a structural rule.

### RES-NODE-016

- Header and node share their first 16 bytes:
  `[+0 reserved][+4 off][+8 size][+0xc type/flags]`.
- As the virtual root node, the header's `off` is `@0x04`, its `size` is
  `@0x08` rootCount and its `type` is `@0x0C`. `@0x10`/`@0x14`
  (regOffset/nodeCount) occupy the slot where a real node keeps its 16-byte
  name.
- The tail open `FUN_004c90a0` reads the header words into the object's node
  slots and copies the base filename into `this+0x10`.
- The serializer `FUN_004cafe0 @0x004cafe7` passes the object as both `this` and
  the root node to the tree-sort.

**Confidence.** High. The identification rests on one instruction, the
serializer handing the same pointer as object and as root node, which turns
"the header looks like a node" into "the header is the node".

### RES-HDR-017

- `FUN_004ce8e0 @0x004ce8fb`: `MOV ESI,[node+4]; SHL 5; ADD [this+0x28]` →
  `&firstChild`.
- So `@0x04` is read and used; it is not recomputed and not a checksum.
- EXP-0014's `{0, nNodes−roots}` law is the packer's serialization order:
  roots-first (`0`) vs roots-last (`ownedNodes`, e.g. graphics/main/sfx).
- Resolves the `rom.exe`/Unknown flag on `RES-HDR-012`, whose value law is
  intact.

**Confidence.** High. The dereference is quoted, so the field is read and used.
The corpus law it replaces had two candidate explanations and no way to choose
between them.

### RES-HDR-018

- The reader branches on bit 4: `FUN_004ce8e0 @0x004ce909 TEST AL,0x10; JZ`.
  Set gives a binary search over sorted children (`bsearch`, cmp `0x4ceaf0`);
  clear gives a linear scan through the CRT `_strnicmp`, ASCII case-insensitive
  (`RES-LOOKUP-023`, `RES-041`). So `{1,17} = {dir, dir+sorted}`.
- The writer sets it in the finalize-sort
  `FUN_004ce660 @0x004ce6c3 OR AL,0x10; MOV [node+0xc]`.
- Bit 4 is EXP-0014's discriminator.
- Bit 31, tested at `FUN_004c90a0 @0x004c92a8` to gate a seek, is clear in all
  shipped archives.
- Resolves the meaning flag on `RES-HDR-013`; its domain and the inline⟹17
  invariant are intact.

**Confidence.** High. The reader branches on the bit and the writer sets it,
both quoted: read and write sides, the strongest form available for a flag.

**Amended.** The linear branch was worded "a linear `strncmp` scan".
`RES-LOOKUP-023` (EXP-0034) disassembled its comparator `FUN_00557030` as the
CRT `_strnicmp`, which folds ASCII `A`-`Z` on both operands, and supersedes that
label ([`retracted.md`](retracted.md)). The bit-4 branch stands.

### RES-NODE-019

- Scanned: `FUN_004ce8e0` (reads `+4/+8/+0xc/+0x10`),
  `FUN_004ce800`/`FUN_004ce9e0` (`name+0x10`, `type+0xc`) and `FUN_004ce660`
  (`+4/+8/+0xc`). No instruction in them dereferences `[node+0]`.
- Reserved and ignored, confirmed at instruction level, not merely "0 in the
  corpus".
- Resolves the meaning flag on `RES-NODE-014`.

**Confidence.** High. An absence, established by enumerating every dereference
in the four scanned lookup, path-walk and finalize-sort functions; that
distinction from a corpus zero is the one the confidence scale turns on.

**Unknown.** Whether a node-touching function outside the four scanned reads
`[node+0]`.

**Amended.** The clause that the four scanned functions are "the four functions
that touch a node" is withdrawn ([`retracted.md`](retracted.md)). Other
functions read or write node words: `FUN_004c9320` (`RES-ORDER-033`,
`RES-MASK-035`), the endpoint `FUN_004ce8c0` (`RES-040`), `FUN_004c9ad0`
(`RES-MASK-035`), the accessors `FUN_004ccc20`/`FUN_004cc670` and the
node-insert `FUN_004ce6d0`, which stores `0` into `node+0x00` (`REG-REC-032`).
The absence over the four scanned functions stands. The headline's unscoped
reading, that no lookup, descent or finalize-sort instruction dereferences the
word, is narrowed to the four scanned functions
([`retracted.md`](retracted.md)).

## Reader code, entry names and the tail open

| ID | Claim | Confidence | Status | Evidence |
|---|---|---|---|---|
| RES-CODE-020 | An address map of the `&YA1` reader, lookup, writer, value accessors, node insert, stream class and archive manager in `rom.exe` (sha256 `942e9b72…7d367d03`). | High | ● active (partially retracted) | [EXP-0017](../experiments/EXP-0017-res-open/), [EXP-0032](../experiments/EXP-0032-reg-text-and-kinds/), [EXP-0034](../experiments/EXP-0034-res-name-codec/), [EXP-0040](../experiments/EXP-0040-enumeration-audit/), [EXP-0051](../experiments/EXP-0051-res-ru/), [EXP-0052](../experiments/EXP-0052-res-resolution/), [EXP-0348](../experiments/EXP-0348-container-field-operations/) |
| RES-TEXT-021 | In `rom.exe`, the `.res` reader applies no byte-to-character conversion to archive entry names: one verbatim bulk read loads every 16-byte name field. | High | ● active | [EXP-0034](../experiments/EXP-0034-res-name-codec/), [EXP-0040](../experiments/EXP-0040-enumeration-audit/) |
| RES-TEXT-022 | Over all 12 standalone `&YA1` containers of this install, 0 of 45 865 entry-name bytes are `>=0x80`; the 55 distinct values span `0x27..0x7A`. | High | ● active | [EXP-0034](../experiments/EXP-0034-res-name-codec/) |
| RES-LOOKUP-023 | In `rom.exe`, `.res` child-name lookup reuses `.reg`'s shared `FUN_004ce8e0`; its linear branch is the CRT `_strnicmp`, case-insensitive with a 15-character bound. | High / Medium | ● active (partially retracted) | [EXP-0034](../experiments/EXP-0034-res-name-codec/), [EXP-0040](../experiments/EXP-0040-enumeration-audit/), [EXP-0348](../experiments/EXP-0348-container-field-operations/) |
| RES-IDENT-024 | In `rom.exe`, `FUN_004ce800` matches a path's leading segment against the resolving object's own stored name by a direct, case-sensitive byte compare. | High / Medium | ● active (amended) | [EXP-0034](../experiments/EXP-0034-res-name-codec/), [EXP-0052](../experiments/EXP-0052-res-resolution/) |
| RES-OPEN-026 | In `rom.exe`, the tail open accepts a file when the stream opens, `u32@0 == 0x31415926` and `malloc(@0x14×32)` succeeds; no instruction checks the registry geometry. | High | ● active | [EXP-0051](../experiments/EXP-0051-res-ru/) |
| RES-PATH-025 | In `rom.exe`, path splitting treats `\` and `/` as interchangeable terminators and rewrites neither; the clause denying any lower-casing pass is retracted. | High | ● active (partially retracted) | [EXP-0034](../experiments/EXP-0034-res-name-codec/), [EXP-0052](../experiments/EXP-0052-res-resolution/) |
| RES-OPEN-027 | In `rom.exe`, the registry is located by `@0x10` alone and sized by `@0x14×32` alone; bytes past it are unreachable, and the module cannot observe EOF. | High | ● active | [EXP-0051](../experiments/EXP-0051-res-ru/) |

### RES-CODE-020

- Core: tail open+validate `FUN_004c90a0` (magic `CMP` `@0x004c91f8`); inline
  read+validate `FUN_004cae80` (`@0x004caeaf`); shared name lookup
  `FUN_004ce8e0`; path walk `FUN_004ce800`/`FUN_004ce9e0`; file-resolve
  `FUN_004c9320`; finalize-sort `FUN_004ce660` (writer-only, sole caller
  `FUN_004cafe0`); serializer `FUN_004cafe0`; empty ctor `FUN_004cb0a0`;
  INI→REG compiler `FUN_004cb340`; record comparator `0x4ceaf0`. The magic
  immediate `0x31415926` occurs in exactly four functions: the two `CMP`
  readers `FUN_004c90a0` and `FUN_004cae80`, and the two `MOV` writers
  `FUN_004cb0a0` (`@0x004cb13b`) and `FUN_004cb340` (`@0x004cb624`).
- EXP-0032, value accessors and node insert: `FUN_004cc670` GetString,
  `FUN_004ccc20` GetInt, `FUN_004ccda0` GetDouble; int-array getters
  `FUN_004cd240`/`004cd130`; setters `FUN_004ccad0`, `004cccb0`, `004cd910`,
  `004cda60`, `004cdbb0`, `004cdd10`, `004ce020`; `FUN_004cd370`, labelled
  create-subkey; `FUN_004ce6d0` node insert, which fixes `value`/`size` as
  block-start/count and the 15-char name clamp; `FUN_004cb160` Open.
- EXP-0034: path-resolve wrapper `FUN_004ce8c0`; the linear branch's comparator
  `FUN_00557030` (`_strnicmp`); the sorted branch's search driver `FUN_00557140`
  (`bsearch`, not a comparator: EXP-0040 read it, and the comparison arrives as
  a function-pointer argument); `FUN_004c9b10`, sole caller of `FUN_004c9320`.
- EXP-0051, the tail open's stream class: ctor+open `FUN_00573b09` (vtable
  `0x0059d4a4`); read `FUN_00573dba` (raw `ReadFile`, vt+0x3c); seek
  `FUN_00573e3f` (raw `SetFilePointer`, vt+0x30); throw helpers
  `FUN_00579973`/`FUN_0057988e`; allocator `FUN_00554390 → FUN_005543b0` (CRT
  malloc-retry); the length family `FUN_00579c7d/cff/d81 → FUN_00579db0`, the
  image's sole `GetFileSize` call site.
- EXP-0052, the archive manager and its search: singleton `0x005f2140`
  (archive `CObArray` `+0x00`/vt `0x59bb00`, directory `CStringArray`
  `+0x14`/vt `0x59bae8`, loose-file `CObArray` `+0x28`/vt `0x59bad0`); ctor
  `FUN_004c93b0` and static-init/destructor thunks `FUN_004c9380`,
  `ORPHAN[004c9370..004c937a]`, `ORPHAN[004c93a0..004c93aa]`; free registrars
  `FUN_004c9910` AddArchive, `FUN_004c9940` AddDirectory, `FUN_004c9a20`
  update.lst-loader; methods `FUN_004c95e0`, `FUN_004c9970`, `FUN_004c9ad0`
  (node mask); path fold `FUN_004c9580` + CRT `tolower` `FUN_005562d0` + gate
  `0x005c1d00`; resolver `FUN_004c9b10`; public open `FUN_004c9f10` (manager
  vt slot `0x0059bb40`); directory-existence check `FUN_004ca300`; startup
  `FUN_004709e0` with catch funclets `FUN_00470e0d`, `FUN_00470f35`,
  `FUN_004710c7`; the on-demand `World.res` site `FUN_004db20d`.
- Archive-name string pool `0x005be34c..0x005be434`, plus `World.res`
  `0x005c6934` and the default extension `".res"` `0x005c1d04`.

**Confidence.** High. An address map of a named binary, pinned to its sha256,
so checkable rather than believable; "exactly four functions" is a whole-image
immediate scan, not a sample.

**Amended.** The 004cd370 create-subkey label is partially retracted by
`REG-106` (EXP-0348, [`retracted.md`](retracted.md)): it is a typed array
getter, the dword-array getter with inline output resize logic; `004cc540`
remains the separately identified section creator, and the other entries keep
their scope. EXP-0040 corrected the label of `FUN_00557140`, first given as one
of the two lookup comparators ([`retracted.md`](retracted.md)).

### RES-TEXT-021

- The tail open `FUN_004c90a0` reads the entire node table, every 16-byte name
  field included, through one verbatim bulk stream `Read()` call, with no
  per-name loop or transform between the allocation and the read.
- No code-page API (`MultiByteToWideChar`, `WideCharToMultiByte`,
  `CharToOemA`, `OemToCharA`, `GetACP`, `GetOEMCP`, `LCMapStringA`) has any
  caller in `0x004c9000-0x004cf000` (whole-image reachability sweep).
  The callers lie in `0x0055xxxx–0x0057xxxx` (EXP-0040).
- The one name-shaped operation in the function, deriving the archive's own base
  filename from its OS open-path, is also a raw byte copy with no fold.
- `REG-TEXT-036` reached the same conclusion for `.reg`; this claim re-derives it
  independently from `.res`'s own distinct reader and does not assume it carries
  over.
- No code page is backed by the binary for entry names: CP866, CP1251 or any
  other code page applied to display a name is a rendering choice, not a
  transcription, as `REG-TEXT-036` finds for `.reg` text.

**Confidence.** High. An absence over a whole-image reachability sweep, plus the
single bulk `Read` that leaves no room for a per-name transform, re-derived from
`.res`'s own reader rather than inherited from `REG-TEXT-036`. EXP-0040
re-checked it on a repaired function table (the vtable blind spot closed, 607
functions created): 0 code-page references inside `0x004c9000–0x004cf000` and 0
anywhere below `0x550000`, so none from the game's own code at all, which is
stronger than this claim states; 0 in orphan code.

### RES-TEXT-022

- Instrument: `tools/resnametext`, sweeping all 12 standalone `&YA1`
  containers. A node's name is the bytes at record `+0x10` truncated at the
  first `0x00`/`0xCD` (`RES-NODE-007`'s own convention).
- 0 of 45 865 examined name bytes are `>=0x80` across all 4 592 corpus nodes;
  55 distinct values, range `0x27..0x7A`.
- This reproduces `EXP-0001`'s `nameBad=0` invariant as an explicit histogram.
- The corpus cannot discriminate any candidate code page that agrees on ASCII:
  the same limitation `REG-TEXT-037` recorded for `.reg`.

**Confidence.** High for this install and this sweep.

### RES-LOOKUP-023

- In-tree child-name lookup during `.res` path resolution calls the exact
  shared function `.reg` uses (`FUN_004ce8e0`, `RES-HDR-017`/`REG-REC-032`).
- Unsorted/linear branch: the exact CRT `_strnicmp`, addressed as
  `FUN_00557030`, disassembled and byte-identical in shape to `REG-TEXT-036`'s
  description. ASCII `A`-`Z` fold only: `CMP AH,0x41`/`CMP AH,0x5a`/`ADD AH,0x20`
  on both operands, no locale active. Case-insensitive, 15-character bound.
- Sorted branch: the `bsearch` driver `FUN_00557140`, a generic bisection loop,
  is handed the same comparator address (`0x4ceaf0`) that `.reg`'s sorted
  lookups cite: the identical shared code, not an independent instance.
  EXP-0034 did not read the comparator's instructions; its byte range had no
  function in the Ghidra auto-analysis. `RES-041` and `REG-100` read them:
  `0x4ceaf0` compares unsigned bytes, case-sensitively.
- EXP-0040 read `FUN_00557140`: it is the CRT's generic `bsearch` driver, not a
  comparator. The comparison function arrives as a pointer argument and is
  called indirectly at `00557181 CALL dword ptr [ESP+0x2c]`, so reading this
  function could never settle the sorted branch's semantics, which live in the
  callee.

**Confidence.** High for the linear branch, disassembled. Medium for the sorted
branch's comparator on this claim's own evidence: the address is confirmed
shared, and EXP-0034 did not read its instructions. EXP-0040's reading of
`FUN_00557140` supports the Medium for that better reason. The comparator's
semantics carry `RES-041`'s High. EXP-0040's audit used this claim as its
positive control: flagged on shape, it came back correct.

**Amended.** The clause that the sorted branch's actual comparator is still
unidentified is partially retracted by `RES-041` and `REG-100` (EXP-0348,
[`retracted.md`](retracted.md)): `004ceaf0` is unsigned bytewise and
case-sensitive. Manager query preprocessing is separate and precedes it, so the
inner comparison alone does not describe original query-case behaviour. The
linear comparison stands. EXP-0040 corrected the label of `FUN_00557140`, first
given as a lookup comparator ([`retracted.md`](retracted.md)).

### RES-IDENT-024

- The leading path segment handed to a resolving object (`FUN_004ce800`) is
  matched against that object's own stored name (`this+0x10`) by a direct,
  case-sensitive byte compare (`CMP AL,CL`): no fold, no code-page involvement.
  It is a distinct, narrower comparison than `RES-LOOKUP-023`'s in-tree fold.
- Corroboration: `FUN_004c9b10` (the sole caller of `FUN_004c9320`) retries the
  same path string against what its loop bound indicates are several resolving
  objects in turn. This is consistent with the check deciding which loaded
  archive a path belongs to before the remainder resolves through the
  case-insensitive in-tree lookup.
- The role, from EXP-0052: the several resolving objects are the manager
  singleton's archive `CObArray`. `FUN_004c9b10` iterates it ascending and takes
  the first non-zero result, and this comparison makes at most one archive
  answer: the leading segment must equal the archive's identity name, or
  `FUN_004ce800` returns 0 (`RES-IDENT-034`, `RES-ORDER-033`).

**Confidence.** High for the comparison itself: `CMP AL,CL` against
`this+0x10`, case-sensitive, read at instruction level. Medium for the role as
this claim infers it from a retry loop's bound; `RES-IDENT-034` (EXP-0052) reads
the loop and the check rather than inferring from the bound, and holds the role
at High.

**Amended.** EXP-0052 confirmed the inferred role (`RES-IDENT-034`) and
corrected the mechanism this claim assumed: the case-sensitive compare is
workable only because every lookup path is lower-cased in full beforehand
(`RES-CASE-036`).

### RES-OPEN-026

- `FUN_004c90a0`, read from a fresh listing, has a complete acceptance
  predicate of three conditions:
  - the stream open succeeds (`FUN_00573b09`, throw on fail);
  - `u32@0 == 0x31415926` (`CMP @0x004c91f8`), the only value check in the
    routine;
  - `malloc(@0x14×32)` succeeds (`@0x004c929f`).
- Seek and read reject only on OS errors: `FUN_00573e3f` is a raw
  `SetFilePointer` (past-EOF legal), `FUN_00573dba` a raw `ReadFile`.
- A short registry read is silent: the byte count is discarded (`@0x004c92cf`),
  and the `malloc`'d table is never zeroed (`FUN_005543b0`, CRT malloc-retry
  shape).
- No instruction computes `EOF − @0x10`, tests 32-byte alignment or validates
  `@0x14`. `RES-HDR-003/004`'s mod-32 and `count = regLen/32` identities are
  corpus facts about packers, not engine checks.
- The seek is gated on `@0x0C` bit 31 clear (`TEST @0x004c92a8`). Bit 31 is
  clear in all 23 shipped archives, both installs.

**Confidence.** High. Every clause is an instruction reading with its anchor
quoted, and the falsification hunt for a geometry check came back empty three
independent ways (`RES-OPEN-027`).

### RES-PATH-025

- `FUN_004ce800`/`FUN_004ce9e0` treat `\` (`0x5c`) and `/` (`0x2f`) as
  interchangeable terminators, tested symmetrically at every position. Neither
  function rewrites one into the other, so no separator-normalization pass
  exists, or is needed. A consumer's `\`→`/` normalization remains unbacked.
- A whole-image reachability sweep, extended beyond `EXP-0032`'s code-page
  needles to check for a lower-casing routine, finds
  `CharLowerA`/`CharLowerBuffA`/`_strlwr`/`_strupr` absent from the binary
  entirely: no import symbol at all.

**Confidence.** High for the separator half, read at every test position.

**Unknown.** The 95-literal figure below is `RES-CASE-036`'s. It conflicts
with `RES-IDENT-034`'s census; see the Unknown on both cards.

**Amended.** The lower-casing half, that no separate case-normalization pass
exists anywhere in `rom.exe`, is retracted and corrected by `RES-CASE-036`
(EXP-0052, [`retracted.md`](retracted.md)). It carried High, but the absence
was established over the import table, which cannot see an inlined fold: the
bound was real and did not cover the claim made from it. A whole-path fold
exists: `FUN_004c9580`, applied to every name and every lookup path entering
the archive manager, gated on `[0x005c1d00]`, which is `1` in the image and
written by no instruction. It is required, because `RES-IDENT-024`'s
archive-identity compare is case-sensitive and 95 shipped literals capitalise
that segment. The fold is a hand-written loop over the CRT `tolower`
(`FUN_005562d0`), which a search for those four names could never reach. A consumer's lower-casing
matches what the engine does.

### RES-OPEN-027

- `@0x10` alone locates the registry (the seek target), and `@0x14×32` alone
  sizes it (the allocation and the single bulk read). Bytes in
  `[@0x10 + @0x14×32, EOF)` are unreachable.
- The walk (`FUN_004ce8e0`) bounds itself by `node.off`/`node.size` only.
  `[this+0x20]` (`nodeCount`) is never read there, so even
  `dir.off+dir.size ≤ nodeCount` is corpus-only.
- Payload resolve (`FUN_004c9320`) returns `{stream, node.off, node.size}` with
  no comparison against `@0x10` or EOF.
- The module cannot observe EOF at all. `GetFileSize` has exactly one call site
  in the image (`@0x00579e0a` in `FUN_00579db0`), reachable only via the stream
  vmethods `+0x18/+0x1c/+0x20` (`FUN_00579c7d/cff/d81`, vtable-dispatch only),
  and 0 of the image's 628 `CALL [reg+0x18/0x1c/0x20]` dispatch sites lie in
  `0x4c9000..0x4cf000`.
- Instrument: `EnumRefs` `sym:`/`callto:`/`re:` on the repaired function table.
  It sees orphan code (1 orphan hit, `@0x0044c902`, outside the module); its
  blind spot is bytes never disassembled.

**Confidence.** High. An enumeration with its instrument and blind spot stated
per AGENTS.md; the `sym:` import route is the immune instrument, and it agrees.

## RU corpus

| ID | Claim | Confidence | Status | Evidence |
|---|---|---|---|---|
| RES-GEOM-028 | RU `MAIN.RES` is a fully valid 504-node tail archive plus 23 stale trailing bytes that no engine read touches; the residue is unique in both corpora. | Medium | ● active | [EXP-0051](../experiments/EXP-0051-res-ru/) |
| RES-HDR-029 | The RU corpus refutes `@0x04 ∈ {0, nNodes−roots}` as a format law: RU `SFX.RES` stores 186, the index of its first top-level node. | Medium | ● active | [EXP-0051](../experiments/EXP-0051-res-ru/) |
| RES-HDR-030 | RU ships `@0x0C = 1` on all 11 archives, including `GRAPHICS.RES` and `MAIN.RES`, so EN's `17` on graphics and main is packing, not format. | Medium | ● active | [EXP-0051](../experiments/EXP-0051-res-ru/) |

### RES-GEOM-028

- File: RU `MAIN.RES`, sha256 `be919916…aa95c40`, `gameversions\ru`. Header
  `(regOff 4906389, count 504)`.
- Engine-style parse: 0 violations; 462 files + 42 dirs, 3 roots, reach
  504/504, payloads tile `[24, regOff)` exactly. Region `16151 = 504×32 + 23`.
- The residue `[4922517, EOF)` aligns as a node-record byte-suffix duplicating
  record 499 (`onmap19.256`, `size 0x707`, `type 0`; the full 16-byte name field
  ends exactly at EOF) on a +23-shifted grid: an in-place overwrite layer of an
  earlier registry write, EXP-0050's font-tail shape.
- Per `RES-OPEN-026/027`, no engine read ever touches it.
- Unique in both corpora: residue 0 on EN 12/12 and RU 10/11.
- Leading and short rival framings score dead: 502/504 bad types, reach 3/504;
  `16128 ≤ 16151`.

**Confidence.** Medium. The corpus measurement is exact. The overwrite-layer
interpretation is Medium: the record-suffix fit discriminates against noise, but
the external packer was not read, and no tail writer exists in `rom.exe`.

### RES-HDR-029

- RU `SFX.RES` stores `186`, with its 48 top-level nodes at indices 186..233:
  186 owned nodes before them and 88 after, verified by an unowned-set
  computation independent of `@0x04`.
- `@0x04` is exactly `RES-HDR-017`'s "array index of the first top-level node",
  confirmed on a mid-table value the EN corpus never exhibited. The
  `{0, n−roots}` dichotomy was EN-packer serialization order.
- RU values: `MAIN.RES` 501 (=504−3, roots-last), `WORLD.RES` 3 (=4−1,
  roots-last vs EN roots-first), 8 others 0.
- Amends `RES-HDR-012`.

**Confidence.** Medium (corpus). It raises the reach of `RES-HDR-017`'s High:
the dereference now has a discriminating witness.

### RES-HDR-030

- EN ships `17` on graphics+main; the "17 on graphics/main" association was EN
  packing, not format.
- The flag is honest in both. RU `MAIN.RES` root children are unsorted
  (`text, graphics, id`); EN's are sorted (`graphics, id, text`). This
  corroborates `RES-HDR-018`: a reader honouring bit 4 linear-scans the RU file
  correctly, and a reader assuming graphics/main are always sorted would bsearch
  an unsorted array.
- The domain over all 23 archives remains `{1, 17}`.
- Amends `RES-HDR-013`.

**Confidence.** Medium (corpus).

## Archive set, resolution and the acceptance contract

| ID | Claim | Confidence | Status | Evidence |
|---|---|---|---|---|
| RES-SET-032 | In `rom.exe`, the archives the engine opens come only from its 11 `AddArchive` call sites: ten literals at startup and `World.res` on demand. | High | ● active | [EXP-0052](../experiments/EXP-0052-res-resolution/) |
| RES-ORDER-033 | In `rom.exe`, a path resolves through the archives, then the directories, each in ascending registration order with the first hit winning; there is no priority field and no sort. | High | ● active | [EXP-0052](../experiments/EXP-0052-res-resolution/) |
| RES-IDENT-034 | In `rom.exe`, an archive's identity name is a mandatory, case-sensitive leading path segment, so the archive list is a dispatch rather than a priority. | High | ● active | [EXP-0052](../experiments/EXP-0052-res-resolution/) |
| RES-MASK-035 | In `rom.exe`, `update.lst` is the archive layer's only config-driven input and only override; it masks entries to loose files and adds no archive. | High | ● active | [EXP-0052](../experiments/EXP-0052-res-resolution/) |
| RES-CASE-036 | In `rom.exe`, `FUN_004c9580` lower-cases every path entering the archive manager; the fold is always on and required. | High | ● active | [EXP-0052](../experiments/EXP-0052-res-resolution/) |
| RES-DIR-037 | In `rom.exe`, `dir[0]` is the working directory captured before `main`, and `AddArchive` uses `dir[0]` alone, so an archive not beside the CWD never opens. | High / Medium | ● active | [EXP-0052](../experiments/EXP-0052-res-resolution/) |
| RES-COLL-038 | Over the EN and RU corpora, the 8 live archives show 0 identity collisions, 0 cross-archive path collisions and 0 within-archive case-insensitive duplicates. | Medium | ● active | [EXP-0052](../experiments/EXP-0052-res-resolution/) |
| RES-ACCEPT-031 | A consumer accepts exactly what the engine accepts: assert the magic, read exactly `@0x14` records at `@0x10` when `@0x0C` bit 31 is clear, and never derive the count from region length. | High / Medium / Unknown | ● active (amended) | [EXP-0051](../experiments/EXP-0051-res-ru/) |

### RES-SET-032

- The tail open `FUN_004c90a0` has exactly one call site in the image, inside
  `AddArchive` `FUN_004c95e0`. The set is therefore exactly the `AddArchive`
  (`FUN_004c9910`) call sites: 11 in 2 functions.
- Ten are literal filenames pushed unconditionally by the startup routine
  `FUN_004709e0`, in this order: `graphics.res` `@00470e47`, `main.res`
  `@00470e54`, `patch.res` `@00470e61`, `music.res` `@00470e75`, `video4.res`
  `@00470f1e`, `video8.res` `@00470f2b`, `sfx.res` `@0047108f`, `movies.res`
  `@0047109c`, `scenario.res` `@004710a9`, `speech.res` `@004710bd`.
- The eleventh, `World.res` `@004db32a` (`FUN_004db20d`), is added on demand at
  world load and only after a direct open of `Data.bin` has failed, so the list
  grows after startup and that archive lands last.
- Each is opened as `dir[0] + "\" + tolower(name)`, with `.res` appended only
  when the given name has no `.` (`@004c90ff`; all ten already do).
- No directory scan and no config file contributes an archive: the image's 8
  `FindFirstFileA` sites include none in `0x4c9000..0x4cf000`, and `update.lst`
  (the sole config input) only masks (`RES-MASK-035`).
- A failed open throws and is swallowed by the startup's catch funclets
  `FUN_00470e0d`/`FUN_00470f35`/`FUN_004710c7`, which set the "no CD" flag
  `0x005eb5a4` and resume, so a missing archive is simply absent from the list.
- The `-4x`/`-8x` command-line test at `00470e84..00470f14` selects only the
  string copied to `0x005f0050`. Both video archives register on either branch;
  the paths converge at `00470f14`.

**Confidence.** High. An enumeration with instrument and blind spot stated per
AGENTS.md: `EnumRefs callto:`/`refto:` on the repaired function table, which
sees orphan code. It finds 0 orphan hits on the open and the registrar, and the
2 orphan hits on the singleton are named in `RES-DIR-037`; it is blind only to
bytes never disassembled. The bound rests on the tail open's single call site,
which is what makes this a set rather than a sample.

### RES-ORDER-033

- `FUN_004c9b10` is `FUN_004c9320`'s sole caller. Its own sole caller is
  `FUN_004c9f10`, the manager's virtual Open at vtable slot `0x0059bb40`, with
  32 call sites in 28 functions.
- It searches two tiers of the singleton at `0x005f2140`, both by ascending
  index, first hit winning:
  1. The archive `CObArray` at `+0x00` (`m_pData @+0x04`, `m_nSize @+0x08`):
     `MOV ECX,[EAX+ESI*0x4]` `@004c9b5f`, `CALL FUN_004c9320` `@004c9b62`.
     `JNZ @004c9b69` leaves the loop on the first non-zero result and
     `JNZ @004c9b78` returns it, so no later archive is ever consulted.
  2. On exhaustion, the directory `CStringArray` at `+0x14`
     (`m_pData @+0x18`, `m_nSize @+0x1c`): `MOV EDI,[EAX+EBX*0x4]` `@004c9bd7`
     builds `dir[j] + "\" + path` and opens it as a plain file on disk;
     `JNZ @004c9cdd` returns on the first success.
- Both tiers yield the same 12-byte record `{stream, offset, size}`: an archive
  entry gives `{archive stream, node.off, node.size}` (`@004c9353..004c9365`), a
  loose file `{file object, 0, GetLength()}` (`@004c9d86..004c9d92`).
- Registration order fixes the order, because the container appends:
  `FUN_004c95e0` reads `m_nSize` into the index (`@004c9714`, `@004c9723`),
  grows by one and writes `array[oldCount] = newArchive` `@004c98b5`.
- There is no priority field and no sort. The resolver dereferences only
  `m_pData` and `m_nSize`, and the only archive-object fields read anywhere in
  the chain are `+0x10` (identity), `+0x34` (stream) and the node words.

**Confidence.** High. Each clause is an instruction reading with its anchor
quoted, from a raw listing rather than the decompiler. The live rivals, prepend,
last-match-wins, a priority field, a sorted key and a second resolver, were each
killed separately, and the "one archive at a time" framing dies on the loop
itself.

### RES-IDENT-034

- `FUN_004ce800` compares the lookup path's leading bytes against the resolving
  object's stored name (`this+0x10`) one byte at a time, case-sensitively
  (`CMP AL,CL @004ce81f`), and returns 0 on any mismatch
  (`JNZ @004ce821` → `XOR EAX,EAX`). It descends the tree only once the identity
  is consumed, or the path reaches a separator, or the identity is empty.
- The check is skipped entirely when the path begins with `\` or `/`
  (`@004ce80a`, `@004ce811`). Only then is each archive searched from its own
  root, and only then can `RES-ORDER-033`'s order decide an outcome.
- The tail open `FUN_004c90a0` derives the identity from the path it opened:
  the basename (walk back to `\` or `/`, `@004c918a..@004c919c`), truncated at
  the first `.` (`@004c9187`), clamped to 15 characters, with bit 28
  (`0x10000000`) OR'd into the root node's type word when it was longer
  (`@004c91bf`).
- So `graphics.res` answers to `graphics\…` and to nothing else, and the
  shipped archives are disjoint namespaces by construction.
- Corroboration from the shipped call sites: every path-shaped literal in
  `.data` begins with one of these identities: 335 `graphics\`, 95
  `SFX\`+`sfx\`, 44 `main\`, 21 `music\`, 10 `movies\`, 7 `World\`, 7
  `scenario\`+`Scenario\`, 4 `video4\`, 4 `speech\`, 1 `patch\`. This includes
  `main\graphics\chrgen\leftup.bmp`, where `graphics` is a subdirectory of
  `main`.
- Resolves `RES-IDENT-024`'s Medium role clause.

**Confidence.** High. The comparison, its case-sensitivity, its skip conditions
and the identity's derivation are each read at instruction level with anchors.
This supplies exactly the dispatch rule `RES-IDENT-024` said "a reader that
needs [it] must have read"; the literal census is independent corroboration,
not the basis.

**Unknown.** This census counts 95 `SFX\`+`sfx\` literals and 7 `World\`;
`RES-CASE-036` counts 95 literals that capitalise the identity segment, `World\`
and `Scenario\` among them. Both hold only if the lower-case `sfx\` literals
number exactly the capitalised literals of every other identity. The EXP-0052
record states the 95 capitalised total with one `SFX\` example and no
per-identity case split, so it does not settle which count is right.

### RES-MASK-035

- `FUN_004c9a20` (sole call site `@004710e3`, argument
  `0x005be340 = "update.lst"`) opens the file with `fopen`, reads lines of at
  most `0xff` bytes, skips any line whose first character is not alphanumeric,
  trims `\n`/`\r` and passes each to `FUN_004c9ad0`.
- `FUN_004c9ad0` resolves the line against the archive array ascending, first
  hit winning (`@004c9ae6`, `@004c9af0`), and sets bit 29 on the found node's
  in-memory type word: `OR dword ptr [EAX+0xc],0x20000000 @004c9b00`.
- `FUN_004c9320` then returns the scalar `1` for that node instead of a record
  (`TEST [ESI+0xc],0x20000000 @004c9339` → `MOV EAX,1 @004c9342`).
  `RES-ORDER-033`'s phase 1 treats `1` as "stop searching archives":
  `CMP EAX,0x1 @004c9b75` falls through to the directory tier.
- The net effect: a masked entry is served from a loose file on disk instead of
  from its archive.
- Neither shipped install contains an `update.lst`, so the mechanism ships
  unused.
- `Data_updated.lst`, referenced by `FUN_004dc65d`, is not routed here:
  `FUN_004c9a20` has one call site.

**Confidence.** High. Every clause is an instruction reading with its anchor.
The bit is written in one place and read in one place, and both are quoted: the
read/write pair that is the strongest available form for a flag.

### RES-CASE-036

- Every path entering the archive manager passes `FUN_004c9580`: the name in
  `AddArchive` (`@004c9920`, again at `@004c9615`), the directory in
  `AddDirectory` (`@004c9950`, `@004c999c`), and the lookup path in both the
  public open `FUN_004c9f10` (`@004c9f42`) and the resolver `FUN_004c9b10`
  (`@004c9b41`).
- When `[0x005c1d00] != 0`, `FUN_004c9580` copies the string byte by byte
  through the CRT `tolower` `FUN_005562d0` (`CMP 0x41`/`CMP 0x5a`/`ADD 0x20` on
  its no-locale path).
- `[0x005c1d00]` is `01 00 00 00` in the file image and has exactly one
  reference in the whole image, the read itself (`EnumRefs refto:5c1d00`), so
  nothing can turn the fold off: it is unconditionally on.
- It is load-bearing, not cosmetic. 95 of the shipped lookup literals capitalise
  the identity segment (`SFX\ChrGen\Skill\MAstral.wav`, `Scenario\…`,
  `World\…`), and `RES-IDENT-034`'s identity compare is case-sensitive, so
  without the fold each of those 95 would fail to find its archive.
- `CharLowerA`/`CharLowerBuffA`/`_strlwr`/`_strupr` are absent from the import
  table; the fold is a hand-written loop over the CRT `tolower`, which no search
  for those four names could reach.
- Corrects a clause of `RES-PATH-025` (see `retracted.md`), whose evidence was
  sound and whose scope was not. Its separator half is untouched and stands: no
  `\`↔`/` rewriting exists anywhere.

**Confidence.** High. The gate's value is read from the image, its reference
count enumerated with the instrument stated, the transform read at instruction
level, and the 95-literal census shows the fold is required rather than merely
present.

**Unknown.** The 95 capitalised literals conflict with `RES-IDENT-034`'s
census, which counts 95 `SFX\`+`sfx\` literals and 7 `World\`. Both hold only
if the lower-case `sfx\` literals number exactly the capitalised literals of
every other identity. The EXP-0052 record gives no per-identity case split.
The fold's existence and gate do not rest on the count.

### RES-DIR-037

- `FUN_004c93b0` constructs the manager singleton `0x005f2140` before `main`. It
  is reached from the static-init thunk pair `ORPHAN[004c9370..004c937a]`
  (`CALL 004c9380` = `MOV ECX,0x5f2140; JMP 004c93b0`, then `JMP 004c9390` =
  `PUSH 0x4c93a0; CALL atexit`), with the registered destructor thunk
  `ORPHAN[004c93a0..004c93aa]`.
- It initialises three arrays (`+0x00` archives, vtable `0x59bb00`; `+0x14`
  directories, `0x59bae8`; `+0x28` open loose files, `0x59bad0`), then calls
  `GetCurrentDirectoryA` and `AddDirectory(cwd)` (`@004c9425`, `@004c9432`). So
  `dir[0]` is the process working directory captured before `main`.
- Startup then adds, in order:
  1. `INSTALLDIR` from `HKLM\SOFTWARE\1C\Allods` (`@00470c50`, conditional on
     `RegQueryValueExA`, and `SetCurrentDirectoryA`'d first `@00470c43`);
  2. `GetTempPathA` with its trailing separator stripped (`@00470cef`,
     unconditional);
  3. the `CD` registry value + `Allods` (`@00470dfc`, conditional).
- Two further `AddDirectory` sites exist outside startup:
  `FUN_004dc65d @004dc6d7` and `FUN_0053a200 @0053a30e`.
- `AddArchive` uses `dir[0]` alone: `XOR EAX,EAX`, then an inlined `GetAt`
  bounds assert `@004c9622` and `MOV EDI,[ECX+EAX*0x4]` `@004c963b`, with no
  loop. An archive not sitting beside the CWD therefore never opens, whatever
  other directories are registered.
- At `@00470d83..@00470db0` the CD volume label is compared against
  `"ALLODS"`, and the result is followed by an unconditional `JMP`, so it is not
  branched on.

**Confidence.** High for the code: each site and the `dir[0]` indexing are read
at instruction level, and both orphan hits that `EnumRefs refto:5f2140` returned
are named here rather than dropped. Medium for which directories a given launch
actually holds: the CWD and the two registry values are properties of the
environment, not of the image, and were not observed at runtime.

### RES-COLL-038

- Instrument: `tools/resorder` replays the `RES-SET-032` registration walk over
  each install root, with an engine-style parse per `RES-ACCEPT-031`.
- 8 live archives in both releases, at indices 0 graphics · 1 main · 2 patch ·
  3 sfx · 4 movies · 5 scenario · 6 speech · 7 world.
- Over those: 0 identity collisions; 0 cross-archive path collisions (EN 3980
  distinct internal paths, RU 3976); 0 within-archive case-insensitive duplicate
  paths. No entry name in either release is reachable through two archives, and
  none is ambiguous inside one.
- Three of the ten registered names are never live: `music.res`, `video4.res`
  and `video8.res` ship under `Allods\` only, and `RES-DIR-037`'s
  `dir[0]`-only rule means their opens throw and are swallowed.
- `patch.res` is not an override archive. It holds exactly one entry,
  `patch.txt` (EN 1418 B, RU 1480 B), addressed as `patch\patch.txt`, and is
  registered after `graphics` and `main` in any case.
- RU containers are upper-cased on disk (`MAIN.RES`), while the engine opens
  the lower-cased name; Win32 resolves that case-insensitively.
- The corpus geometry reproduces EXP-0051 exactly: RU `MAIN.RES` residue 23, RU
  `SFX.RES` `@0x04` = 186, `@0x0C` = 17 on EN graphics/main and 1 on all RU.

**Confidence.** Medium, capped there deliberately. A census over 23 containers
can show that no collision occurs, never that one is impossible; what forbids
collisions is `RES-IDENT-034`, a separate claim graded High. The
"3 of 10 never live" clause is a property of the shipped layout plus a normal
launch, not of the image: a launch whose CWD were `Allods\` would open those
three and none of the others.

### RES-ACCEPT-031

- Assert, engine-checked:
  - `u32@0 == "&YA1"`;
  - with `@0x0C` bit 31 clear, locate the registry at `@0x10`, take exactly
    `@0x14` 32-byte records, and ignore `[@0x10 + @0x14×32, EOF)`;
  - with bit 31 set, read the node table at the stream position after the six
    header dwords (`RES-039`);
  - do not require `(EOF−@0x10) % 32 == 0`;
  - do not derive the count from region length; that is what rejected RU
    `MAIN.RES`.
- Corpus clauses (EN+RU, 23 archives):
  - `@0x10 + @0x14×32 ≤ EOF` (residue 0 on 22/23, 23 on RU `MAIN.RES`);
  - `24 ≤ @0x10 ≤ EOF`;
  - types ∈ {0,1};
  - dir ranges in-table;
  - payloads tile `[24, @0x10)` exactly;
  - `@0x0C ∈ {1,17}` with bit 4 honest.
- A reader may assert the corpus clauses as hardening, knowing it is then
  stricter than the engine, which tolerates even a silently-truncated registry
  read and unchecked dir bounds.

**Confidence.** High for the assert clauses, instruction-level per
`RES-OPEN-026/027`; the bit 31 set clause carries `RES-039`'s High. Medium for
the corpus clauses, the cap for a corpus measurement.

**Unknown.** Engine runtime behaviour outside the corpus envelope: garbage node
words and short reads. No defensive path exists; what the engine checks is
pinned, and what it would do is not.

**Amended.** `RES-039` (EXP-0348) closes the `@0x0C` bit 31 item, formerly in
the Unknown: bit 31 set suppresses the seek to `@0x10`, and the tail open
reads the table at the stream position after the six header dwords. The
`@0x10` locate clause is narrowed to bit 31 clear, which holds on all 23
shipped archives (`RES-OPEN-026`; [`retracted.md`](retracted.md)).

## Relocated template row

| ID | Claim | Confidence | Status | Evidence |
|---|---|---|---|---|
| RES-HDR-001 | _statement_ | High | ● active (amended) | [EXP-0001](../experiments/EXP-0001-res-container/) |

### RES-HDR-001

- This row is not a finding about the RES/LM container. `claims/registry.md`
  carried it, with its table header, as generic contributor scaffolding for any
  not-yet-populated area: "copy this row into a format section when its first
  claim lands".
- The knowledge repository's `k1` edition dropped it from the public index as
  internal review narrative, under that page's own Publication-boundary rule.
- It sits here because the id it names belongs to this ledger's own
  `RES-HDR-*` series, which otherwise runs `002`–`030`, and was never promoted
  into a published row here.
- The move touches and resolves nothing about the id's confidence, status or
  evidence; the id remains exactly as unpromoted as it was.

**Amended.** The Evidence link read `../experiments/EXP-0001-res-header/`, a
directory that does not exist; the only EXP-0001 directory is
`EXP-0001-res-container`, and the link now names it
([`retracted.md`](retracted.md)). The statement stays `_statement_`.

## Shared container field operations

| ID | Claim | Confidence | Status | Evidence |
|---|---|---|---|---|
| RES-039 | The complete tail open 004c90a0 treats header +0x0c as flags, not a two-value enum: bit 31 alone changes where the node table is read. | High | ✔ promoted | [EXP-0348](../experiments/EXP-0348-container-field-operations/) |
| RES-040 | RES lookup, descent and endpoint are distinct: 004ce800 requires `(kind & 0x40000001) == 1` before another component; endpoint 004ce8c0 rejects bits 4/30 but not bit 0 alone. | High / Unknown | ✔ promoted | [EXP-0348](../experiments/EXP-0348-container-field-operations/) |
| RES-041 | RES sorted lookup uses `REG-100`'s unsigned-byte, case-sensitive comparator 004ceaf0; under the named no-locale state, 004c9580 first folds the query's ASCII. | High / Unknown | ✔ promoted | [EXP-0348](../experiments/EXP-0348-container-field-operations/) |

### RES-039

- In 32 single-bit controls on a bounded header, bit 31 alone changes node-table
  positioning: clear seeks to header +0x10; set reads at the stream position
  after six header dwords.
- Distinct tables at offsets 24/64 discriminate the routes.
- The table is bulk-read without per-node kind validation.
- This is RES tail framing, not REG inline pool framing.
- Evidence: V166..V198.

**Confidence.** High for the complete open and the synthetic stream/allocator
controls. Native stream failure remains outside scope.

### RES-040

- Shared 004ce800 requires `(kind & 0x40000001) == 1` before another component.
- Endpoint 004ce8c0 rejects bits 4/30 but not bit 0 alone. Name matching need
  not reject the same node.
- Wrapper 004c9320 then applies the established bit-29 loose-file sentinel.
- Deletion/copy follows `REG-104`; no whole-dword two-value restriction follows.
- Evidence: V001..V066.

**Confidence.** High for the complete local bodies and one-bit controls.

**Unknown.** Whole-manager/native file access and malformed paths.

### RES-041

- The comparator is `REG-100`'s unsigned-byte, case-sensitive 004ceaf0; it is no
  longer unidentified.
- Under the named no-locale state, 004c9580 first folds ASCII through 005562d0;
  caller 004c9b41 binds preprocessing to resolver input.
- Controls turn both R/KEY and r/key into r/key: sorted lookup matches stored
  key and misses stored Key; unsorted lookup matches either.
- Query folding does not modify stored names. The high-byte control preserves
  its high byte while folding its ASCII suffix.
- Storage, lookup and display remain separate; no universal encoding follows.
- Evidence: V199..V208.

**Confidence.** High for the retained comparator evidence, fresh
bodies/arguments and conditional composition controls.

**Unknown.** Nonzero locale and complete native-manager/display behavior.

## Open questions

- How many path-shaped `rom.exe` `.data` literals capitalise their identity
  segment. `RES-IDENT-034` counts 95 `SFX\`+`sfx\` and 7 `World\`;
  `RES-CASE-036` counts 95 capitalised literals, `World\` and `Scenario\` among
  them, and `RES-PATH-025` repeats that figure. The discriminating measurement is a case-split census
  of those literals: per identity, how many spell the leading segment other
  than the lower-case archive identity, with the `SFX\`/`sfx\` and
  `Scenario\`/`scenario\` splits.
