# RES / LM container (`&YA1`) — public functional specification

Level 3 functional edition. Promoted claim families are recorded in
[`claims/res.md`](../../claims/res.md). This page states the binary grammar and lookup
behaviour needed by an independent reader/writer; instruction addresses and
reverse-engineering listings remain private evidence.

## Overview

A tail-registry RES/LM archive consists of:

```text
+---------+---------------------+-----------------------------+-------+
| header  | payload data        | node registry               | slack |
| 24 B    | [0x18, regOffset)   | nodeCount * 32-byte nodes   | 0..n  |
+---------+---------------------+-----------------------------+-------+
```

The registry location and size come from header fields, not from EOF. Shipped files can
contain trailing bytes after the registry, so deriving `nodeCount` from file length is
not compatible with all observed releases. Bytes past the node array are unreachable to
the original reader. — `RES-HDR-003`, `RES-HDR-004`, `RES-GEOM-028`, `RES-OPEN-026`,
`RES-OPEN-027`

The same container is seen under more than one extension, including a valid empty
archive with no nodes. The extension plays no part in the binary: one reader validates
the magic. — `RES-SCOPE-010`

## Header

All integers are little-endian.

| Offset | Type | Meaning |
|---|---|---|
| `0x00` | `u32` | magic `0x31415926` (`&YA1` as bytes) |
| `0x04` | `u32` | root node `off`: first top-level child index |
| `0x08` | `u32` | root node `size`: top-level child count |
| `0x0C` | `u32` | root node type/flags; bit `0x10` means children are sorted |
| `0x10` | `u32` | byte offset of the node registry |
| `0x14` | `u32` | number of 32-byte nodes |

The header is 24 bytes and payload data begins at offset 24. The magic value is the only
field the original validates. — `RES-HDR-002`, `RES-MAGIC-001`

The header's first 16 bytes have the same logical shape as a node's first 16 bytes, so
the header functions as the tree's virtual root node. — `RES-NODE-016`,
`RES-HDR-017`, `RES-HDR-018`

`0x08` is that root node's own child count. The value domains of `0x04` and `0x0C` are
per-file constants rather than quantities a reader may derive: `0x04` is a legal array
index anywhere in the table, not a count or a checksum, and `0x0C` carries the directory
type with the optional sorted flag and no third value in the examined corpus. —
`RES-HDR-005`, `RES-HDR-012`, `RES-HDR-013`

## Node record

Each registry node is 32 bytes:

| Offset | Type | File node | Directory node |
|---|---|---|---|
| `0x00` | `u32` | reserved | reserved |
| `0x04` | `u32` | payload byte offset | first child node index |
| `0x08` | `u32` | payload length | child count |
| `0x0C` | `u32` | type `0` | type `1`, optionally OR `0x10` for sorted children |
| `0x10` | `char[16]` | NUL-terminated name | NUL-terminated name |

Observed writers may leave non-semantic padding values after a short name. Readers
should use the terminator/name bound rather than treating padding bytes as content. —
`RES-NODE-007`, `RES-NODE-008`, `RES-NODE-019`

`0x00` is reserved. It is zero on every node of every examined tail-registry archive and
is never read by lookup, descent or sort, which refutes reading it as a per-node hash,
id or checksum. — `RES-NODE-011`, `RES-NODE-014`

## Minimal read algorithm

```text
require u32(file, 0x00) == 0x31415926

rootOff     = u32(file, 0x04)
rootCount   = u32(file, 0x08)
rootFlags   = u32(file, 0x0C)
regOffset   = u32(file, 0x10)
nodeCount   = u32(file, 0x14)

for i in [0, nodeCount):
    node[i] = parse32(file[regOffset + i*32 : regOffset + (i+1)*32])

root.children = node[rootOff : rootOff + rootCount]
```

A file node addresses `file[node.off : node.off + node.size]`. A directory node addresses
`node[node.off : node.off + node.size]`.

Defensive implementations should bounds-check these ranges even where the original is
permissive; such hardening is a safety choice, not an additional format invariant.

## Directory lookup

A directory's `0x10` flag indicates that its children are sorted. The original uses a
sorted lookup in that case and a linear lookup otherwise. Both perform ASCII-oriented
case-insensitive child-name comparison. — `RES-HDR-018`, `RES-LOOKUP-023`

Path processing has several important compatibility properties:

- both `/` and `\\` are accepted as component separators;
- archive-manager input paths are lowercased before resolution in the reached path;
- child lookup is case-insensitive over the original byte-oriented ASCII domain;
- no Unicode normalization or general code-page conversion is part of the archive
  grammar. — `RES-TEXT-021`, `RES-PATH-025`, `RES-CASE-036`

No byte-to-character conversion happens anywhere on the read path, and the sibling REG
text store reaches the same conclusion through its own distinct reader. No examined
archive name byte reaches the high half of the byte range, so the corpus cannot
discriminate between candidate renderings that agree on ASCII: a code page for an entry
name is a display convention of the consumer, not a property of the format. —
`RES-TEXT-022`, `REG-TEXT-036`

## Archive identities and resolver

An opened archive has an identity derived from its archive filename. Resource paths are
normally of the form:

```text
<archive-identity>\<path-inside-archive>
```

The resolver searches registered archives in registration order and returns the first
match, then falls back to configured loose-file directories. Because normal paths carry
an archive identity, the shipped archives mostly form disjoint namespaces rather than a
single overlay-by-priority system. — `RES-SET-032`, `RES-ORDER-033`, `RES-IDENT-034`

The leading path segment is matched against the resolving archive's own stored name by a
direct case-sensitive byte compare — a narrower comparison than the case-insensitive
child lookup used for every later component — which is why the whole path is lowercased
first. The loose-file tier is anchored on the process working directory captured before
entry, so which archives exist at all is an environment property rather than an archive
property. Over the live archives of each preserved release the measured identity and
cross-archive path collision counts are both zero, which is what makes registration
order a search order and not a priority. — `RES-IDENT-024`, `RES-DIR-037`, `RES-COLL-038`

The original also supports an update-list mechanism that can mark an archive node so
resolution falls through to a loose file. A compatible implementation may model this as
an explicit override from archive member to filesystem resource. — `RES-MASK-035`

The complete shipped archive-name inventory and path corpus are evidence about one game
installation, not part of the container grammar, and are intentionally omitted here.

## Acceptance versus corpus invariants

The format grammar and the shipped corpus should not be conflated.

Observed corpus properties include well-formed trees, in-range payloads and exact
payload tiling on the examined archives. However, some releases carry trailing registry
slack, and some header-value patterns differ between EN and RU data. Therefore a reader
must use the explicit registry offset/count fields rather than promoting one release's
packing pattern into a universal requirement. — `RES-ACCEPT-031`, `RES-GEOM-028`,
`RES-HDR-029`, `RES-HDR-030`

The nodes form a tree with every node reachable and no cycles, and file payload ranges
tile the data region exactly, over every blob measured. That is a measurement of the
shipped corpus, not a rule the original reader enforces. — `RES-TREE-009`

Useful defensive checks for independent software include:

```text
regOffset >= 24
regOffset + nodeCount*32 <= fileSize
file node: off + size <= regOffset
folder node: off + size <= nodeCount
```

These checks protect the replacement implementation; they are not all asserted by the
original reader.

## Two `&YA1` families

The same magic is used by two structurally distinct storage families:

1. this **tail-registry RES/LM archive**, where `0x10` is a registry byte offset;
2. the **inline REG-style store**, where the same position belongs to a different record
   framing.

Do not select the parser from the magic alone. Use the surrounding file/context and the
appropriate specification. — `RES-SCOPE-015`

## Container digest versus selected member

A differing whole-container digest does not imply that the members a consumer selects
differ. Across the two preserved roots several containers have different whole-file
digests while every map and store member selected by the current save corpus is
byte-identical between them; one measured data payload differs at equal length. Compare
the member actually used, not the archive that carries it, and do not read a differing
container digest as evidence about an unmeasured member. The member inventory and the
individual digests are corpus evidence and stay private. — `SAV-SUFF-301`

## Writer guidance

A conservative writer for the tail-registry form should:

1. reserve/write the 24-byte root header;
2. append payload bytes and record each file node's byte range;
3. emit a flat 32-byte node array with directory child ranges;
4. set the root child range in the header;
5. write `regOffset` and `nodeCount` explicitly;
6. set the sorted-child flag only when the relevant child range is actually sorted under
   the lookup comparison used by the consumer.

There is no interoperability need to reproduce incidental padding/slack values from the
original packer.

## Publication boundary

The private research contains exact executable addresses, reader/writer instruction
sequences, complete shipped archive inventories and corpus hashes used to establish the
rules above, including the address map of the original reader, lookup and writer. They
are not required to implement the public binary grammar and are not reproduced in this
edition. — `RES-CODE-020`

Unknown or malformed-input behaviour should remain explicit rather than being inferred
from a third-party implementation.
