# FAME hall-of-fame (`famehall.dat`) — specification (single-sample)

Level 3. Promoted, evidence-backed claims only. The single-file corpus and its
10/10 exact tiling are `FAME-HDR-001`, `FAME-REC-002`, `FAME-NAME-003`,
`FAME-SCORE-004`, and `FAME-UNK-005`. The 6/6 mutant boundary is our own
validator's bound (`FAME-HDR-001`, `FAME-REC-002`); the default-table
interpretation is `FAME-DEFAULT-006`.

**Status: the writer is read (`FAME-WRITE-007`), the content is still one table (◐).** `FUN_00489de0` writes
the count, then per record `strlen+1`, the string including its NUL, and three `u32` — striding
`0x10` through an in-memory array of `{CString*, u32, u32, u32}`. It is a **raw `CFile`**, not a
`CArchive`: no `Asg&`, no compression, no class record, nothing shared with the save path. Its one
caller opens the file `modeCreate\|modeWrite` and is the main frame's **`WM_DESTROY`** handler, so
the file is rewritten every time the game exits. A 2026-08-02 play session did exactly that and
produced **byte-identical** output to both pristine roots. — FAME-WRITE-007

**Status: single-sample content (◐).** The whole file *is* decoded — this is not a
framing-only spec — but the corpus is **one file**, almost certainly the shipped
default (untouched by play). So the record layout is well-supported (replicated 10×
inside the file + exact tiling + falsification), while the two trailing `u32` per
record are **all zero** here and their meaning is **Unknown**.

Seen as: `famehall.dat` at the install root, 228 bytes. No magic — the file opens on
a bare `u32` count. (The file changes if the game is played; hash-pin the sample.)

## At a glance

A high-score / hall-of-fame table: a 4-byte record count, then that many
variable-length, length-prefixed records that tile the file exactly. Each record
holds a name and a score; the table is sorted strictly descending by score.

```
+-------+-----------------------------------------------------------+
| count | record[0] | record[1] | … | record[count-1]               |
| u32   |                                                           |
+-------+-----------------------------------------------------------+
0     0x04                                                         EOF (228)

record = [ u32 nameLen ][ char[nameLen] name\0 ][ u32 score ][ u32 ][ u32 ]
                                                              \___ 8 reserved ___/
         size = 4 + nameLen + 12
```

## Header (4 bytes, little-endian)

| Off | Type | Field | Notes | Claim |
|-----|------|-------|-------|-------|
| 0x00 | u32 | count | number of records that follow (= 10) | FAME-HDR-001 |

Record stream begins at offset `0x04`. There is no ASCII magic and no separate
version field.

## Record (variable length, little-endian)

| Off (within record) | Type | Field | Notes | Claim |
|-----|------|-------|-------|-------|
| +0x00 | u32 | nameLen | byte length of the name field, **including** its trailing `\0` (`= strlen + 1`) | FAME-NAME-003 |
| +0x04 | char[nameLen] | name | NUL-terminated ASCII | FAME-NAME-003 |
| +0x04+nameLen | u32 | score | strictly descending across the table; spans > 16 bits | FAME-SCORE-004 |
| +0x08+nameLen | u32 | *reserved #1* | 0 in this sample — meaning unknown | FAME-UNK-005 |
| +0x0C+nameLen | u32 | *reserved #2* | 0 in this sample — meaning unknown | FAME-UNK-005 |

Record size `= 4 + nameLen + 12`.

## Sample measurements (not the original payload)

The measured file has ten records and is 228 bytes long. The observed scores are
strictly descending from 70006 to 0; two values exceed 65535, supporting the field-width
finding in `FAME-SCORE-004`. All 80 bytes in the two trailing dwords per record are
zero in this sample (`FAME-UNK-005`); this is not a proved universal zero default.

The original ten-name/ten-score content table is intentionally not reproduced in
this public edition. The unabridged research evidence is identified by
[SOURCE.md](../../SOURCE.md). `FAME-DEFAULT-008` retains the correspondence to the
installed UI table's lines 263..272 and the uncertainty about which copy is drawn.
These omissions change neither the file grammar nor the research confidence.

## Invariants (hold for the sample; enforced by the probe)

- `len(file) >= 4`; `count = u32@0`.
- Exactly `count` records walk: each `[u32 nameLen][nameLen bytes][12 bytes]` is
  in-bounds, `nameLen >= 1`, and `name[nameLen-1] == 0` (NUL-terminated).
- **Exact tiling:** after `count` records the cursor equals file size —
  `4 + Σ(4 + nameLen + 12) == 228`, 0 bytes unaccounted, no overrun.

The repository probe's own validator rejects six synthetic mutants breaking the
contract above (truncation, `count` too small/large, inflated length prefix,
removed NUL terminator) — 6/6 (`FAME-HDR-001`, `FAME-REC-002`). This is a bound
on our validator, not evidence from the original reader.

## Notes & open questions

- **Corpus limit:** one file, and it is almost certainly the untouched **default**
  table (10 seeded names, clean descending ladder ending at a `0` entry, all-zero
  trailing fields). We cannot see how `count` scales, whether it can differ from the
  record total, or the max name length. — FAME-DEFAULT-006
- **Reserved `[u32][u32]`** after each score are zero in every record, so they are
  observed but not decodable. A played-in `famehall.dat` (with non-zero values) would
  be needed to interpret them — candidates by analogy only: a level/rank, a
  mission/difficulty id, or a play-time.
- **Score width:** `u32` is forced (two entries exceed 65535); the descending order
  is read from the static table, not from watching the game write it.
