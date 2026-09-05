# ROM2 ALM map (`M7R\0`) — identity survey

Level 3. Promoted, evidence-backed claims only. Basis: `R2-ASSET-002` (file
header, record-0 metadata), `R2-ASSET-003` (record chain past record 0,
`recordCount`/`formatVersion`), `R2-SESSION-010` (code-side dispatcher arity,
cross-referencing `R2-ASSET-003`'s record-type count from a running binary
rather than a stored file), `R2-ASSET-017`, `R2-ASSET-018`, `R2-ASSET-019`,
`R2-ASSET-020`, `R2-ASSET-021`, `R2-ASSET-022`, `R2-ASSET-023`, `R2-ASSET-024`,
`R2-ASSET-027` (cross-game population/order, per-type payload layout, the
`allods2.exe` record dispatcher, and map extents). Cross-reference only, not
evidence: ROM1's own
[`formats/alm/format.md`](../alm/format.md) (`ALM-HDR-001`, `ALM-FRAME-031`,
`ALM-SEC-004`, `ALM-UNIT-018`, `ALM-META-024`, `ALM-META-025`, `ALM-CLS-036`,
`ALM-OBJ-034`, `ALM-TRIG-021`, `ALM-SACK-065`, `ALM-GRP-041`, cross-reference only).

Seen as: `*.alm`/`*.ALM` at the ROM2 install root (37 files), plus `M7R\0`-magic
payloads nested inside `scenario.res` under a `<n>.alm` name (46 payloads) — 83
maps total, not the ~30 this survey's own brief assumed.

## Result

**File header: identical field layout, not identical values.** All 83 maps carry
magic `4D 37 52 00` ("M7R\0") and file-header field `hdrLen=20` (0 exceptions on
either). Every map's first record reads as type 0 with a `>=0x28`-byte payload,
from which `W`/`H` (map dimensions) are read at payload offsets `0x00`/`0x04`.
ALM-HDR-001's own corpus arithmetic identity `dataSize = 4*W*H + 72` — spanning
the file header's `dataSize` field and record 0's own `W`/`H` fields, three
independently-positioned reads, unaffected by the correction below since `W`/`H`
are read before it applies — holds on **83/83**, 0 exceptions. `R2-ASSET-002`.

**Record 0's own declared size is not ROM1's, and undercounts its own true
span.** Record 0's `payloadSize` field is a corpus-wide constant **644** (ROM1's
own constant: **632**, `ALM-SEC-004`) — 12 bytes more. Separately, record 0's
TRUE span (its header's end to record 1's header's start) is **660** bytes, 16
more than even the declared 644: a `type0PayloadOverhang` this survey's own
tool must add to its cursor arithmetic to reach record 1 correctly. Replayed
against ROM1's own two preserved roots as a falsification check, that same +16
correction does NOT also produce a clean walk (ROM1 EN 0/38 clean, ROM1 RU 0/34
clean, both desyncing at record 1 uniformly) — ruling out "any additive fudge
factor would clean up any corpus." `R2-ASSET-002`.

**Past record 0: identical record framing, once corrected.** Applying the
16-byte correction, all 83 maps reach **exact EOF with 0 residue**, walking the
full 13-record chain with ROM1's unmodified ALM-FRAME-031 20-byte record header
(`tag`/`hdrLen`/`payloadSize`/`typeId`/`f32`) and using `payloadSize` unmodified
to advance. The originally reported divergence within the first four records —
and the far-out `typeId` values read there (e.g. `343807103`, `4292870144`) —
were this survey's own probe carrying record 0's undercounted declared size
forward: a cursor desync, not a ROM2 grammar change. `R2-ASSET-003`.

**Content differences, distinct from framing.** `recordCount` (file header
`+0x0C`) is a corpus-wide constant **13** on all 83 maps (ROM1: constant 10 —
3 more record types recorded, not a framing change). `formatVersion` (`+0x10`)
is **1300** on 45 of the 46 `scenario.res`-nested maps and **1600** on the 37
root maps plus that one exception (ROM1: constant 990). All 13 `typeId` values
`0..12` occur exactly once per map, corpus-wide, in one canonical order —
`[0,1,2,3,5,11,4,9,8,6,7,10,12]`, 83/83 uniform. `R2-ASSET-003`.

**The 13-versus-10 count also holds on the code side, from a shared parsing
primitive.** The per-record header parser is one function, byte-identical
across `rom.exe`, `allods2.exe` and `a2server.exe`, each called exactly once
by its own binary's per-record dispatch loop. Counting each dispatch loop's
own jump-table targets directly: `rom.exe`'s has exactly 10 arms, both ROM2
binaries' have exactly 13 — the same count this page already reports from
record content, now read independently from running code rather than a
stored file. `R2-SESSION-010`.

**Per-type payload-size formulas.** Replaying ROM1's own formulas (`ALM-SEC-004`,
`ALM-GRID-032`, `ALM-GRP-020`, `ALM-UNIT-018`, transcribed unmodified) against
each record's declared `payloadSize`: types 1 (`2*W*H`), 2/3 (`W*H`), 5
(`76*nPlayers`) match exactly on **83/83**, 0 exceptions. Type 6 (ROM1:
`70*nUnits`) matches on **0/83** — but on all 83, the declared value instead
equals exactly `48*nUnits` for the same `nUnits` read from record 0's own
metadata, a clean corpus-uniform alternate constant, not scatter around 70. Why
48 bytes/unit rather than ROM1's 70 is fully accounted for: summing case 6's own
field widths, read from `allods2.exe`'s dispatcher, against each file's own
`formatVersion` reproduces the declared `payloadSize` on 154 of 154 type-6
records in both games (48 at ROM2's two `formatVersion` values, 70 at ROM1's
990). Types 10, 11 and 12 now have their own published per-type formulas as
well (below); types 4, 7, 8, 9 have no published ROM1 per-type formula to
replay. `R2-ASSET-003`, `R2-ASSET-024`, `R2-ASSET-023`.

**ROM2's own order and typeId set strictly contain ROM1's.** ROM1's own two-root
corpus (38 EN maps, 33 of 34 RU maps) carries typeId set `{0..9}` in one
canonical order, `[0,1,2,3,5,4,9,8,6,7]`, uniform. Removing ROM2's three new
types `{10,11,12}` from ROM2's own order leaves ROM1's order byte for byte —
type 11 is inserted between 5 and 4, types 10 and 12 are appended after 7 — not
an independently renumbered sequence. The one RU exception, `Horror.alm`,
declares `recordCount=4` (types `{0,1,2,3}` only) in its own file header — but
this is not an independently authored 4-record map: the RU file is
byte-for-byte identical to the EN root's own full-10-record `Horror.alm` for
its first 262,876 bytes except two bytes (the rewritten `recordCount` and one
byte inside record 0's own payload), and the EN file continues past that point
with a 127,050-byte type-6 record the RU file does not carry. RU's `Horror.alm`
is the EN file truncated after record 3 with `recordCount` rewritten to match,
not a second, differently-authored map. `R2-ASSET-017`.

**Per-typeId, tested against each type's own correct structural rule, the
shared record types either match or diverge uniformly — never split by
`formatVersion`.** Types 1, 2, 3, 5, 7, 8, 9 parse ROM1's own formula to 0
residue on 100% of both `formatVersion` groups (1300, 1600); types 0 and 6 fail
against ROM1's own formula on 100% of both groups (matching their own alternate
constants above); type 4's apparent `formatVersion`-correlated failure under
ROM1's own already-published exact-match rule (`kind==0x21`, `ALM-CLS-036`) —
45/45 clean at 1300, only 19/38 clean at 1600 — disappears (83/83 pass, both
groups) once the dispatcher's own exact kind-extension condition is replayed
instead. That one divergence, scored against ROM1's own rule, is not merely
present but perfectly partitioned: all 19 files failing ROM1's rule are exactly
the 19 files whose `kind` sets bit 24 (`0x1000021`), and all 19 are
`formatVersion` 1600 — 0 of the 45 `formatVersion`-1300 files carry that kind.
The dispatcher's own bit-24 test is not itself version-gated (it is evaluated
unconditionally on every file), but the corpus content that exercises it is a
1600-era authoring feature: no shared type shows a genuine `formatVersion`-gated
SHAPE change once each is tested against the rule that actually governs it, but
type 4's own extended-kind CONTENT is real and cleanly `formatVersion`-bound.
`R2-ASSET-018`, `R2-ASSET-019`.

**The record dispatcher, read from `allods2.exe`.** `FUN_0053ea3a` is ROM2's
counterpart of ROM1's `FUN_00512369` (`ALM-META-024`), located by cross-reference
from its own already-identical-classified header/record-read callees. Its
acceptance gate keeps ROM1's own magic and `recordCount>=3` checks but
widens the `formatVersion` ceiling from ROM1's `<=1001` to `<=1600` — the
engine's own compiled boundary, not merely what ships. Its `switch` has exactly
13 cases (`0..12`) plus a generic, non-rejecting `default`. Types 10, 11 and 12
read their own element counts from type-0's own version-gated fields, never
from their own record payload, through the same generic allocate-then-read
mechanism every case (shared or new) uses — ROM2 extended the existing
dispatcher rather than adding a second one. Sourcing a type's element count
from record 0 rather than its own payload is not unique to the three new
types: cases 1/2/3 (from `W`/`H`) and cases 4, 5, 6, 8 (from four more type-0
fields) use the identical mechanism; only cases 7 and 9 read a count word from
their own payload. What is distinctive about 10, 11 and 12 is narrower —
their own type-0 count fields are themselves `formatVersion`-gated (zero below
the threshold, so the case runs but iterates zero times), where every other
type-0-sourced count is populated unconditionally on every accepted file. Case
`0xb` (type 11)'s body runs three separate flat-array loops bounded by three
counts read in record 0 — a per-element shape resembling, but not the same
count provenance as, ROM1's own type-7 (whose three counts are read from its
OWN payload, not from record 0). Each of the three new types now has a
published per-type size formula, replayed 83/83: type 10 is `16*meta[+0x30]`
(flat 16-byte elements at `map+0x310`); type 11 is
`12*meta[+0x34] + 84*meta[+0x38] + 12*meta[+0x3c]` (three flat arrays at
`map+0x338`/`map+0x324`/`map+0x34c`); type 12 is `28 + 28*meta[+0x40]` (one
fixed 28-byte head at `map+0x374` plus a flat array at `map+0x360`). Field
*meaning* — what each byte of an element represents — is not decoded. Type 4's
28-byte-extension trigger is `kind==0x21 || (kind&0x1000000)!=0` — ROM1's own
`kind==0x21` test plus an independent bit-24 flag; the corpus cannot
discriminate this from a low-byte mask `(kind&0xff)==0x21` (both reach 83/83),
since the only `kind` value with bit 24 set, `0x1000021`, also has low byte
`0x21`. `R2-ASSET-019`, `R2-ASSET-020`, `R2-ASSET-021`, `R2-ASSET-022`,
`R2-ASSET-023`.

**Versioned field growth, and a full account of type 6's stride and record 0's
span.** Case 0 gates four field groups on four `formatVersion` thresholds (all
below both this corpus's own values, so all four groups are corpus-uniformly
present, summing to exactly 28 bytes — record 0's own true 632→660 span growth,
`R2-ASSET-002`, exactly). The complete field-width sum for case 0 (not only the
four version-gated groups) is `48+4+12+4+64+4+4+8+512=660` at ROM2's own
`formatVersion` values and `632` at ROM1's 990 — the two corpus constants
reproduced from the code, on 155 of 155 record-0 payloads in both games; which
specific bytes the file's own declared `payloadSize` (644) counts, versus which
land in the 16-byte undeclared overhang, stays Unknown — that split is a
property of the save path, and this survey's Ghidra pass reads the load path
only, which reads all 660 bytes in fixed order without ever consulting the
declared 644. Case 6 gates six fields on six thresholds; one of them, `0x44c`
(1100), runs the opposite direction from every other threshold in the
dispatcher and is corpus-uniformly ABSENT (every file's `formatVersion`, 1300
or 1600, exceeds 1100) — the largest single term but not the whole account. The
complete field-width sum for case 6 is `36+4+4+0+4=48` at ROM2's own
`formatVersion` values and `36+0+4+28+2=70` at ROM1's 990 (unconditional
fields 36; a `0x47e` gate +4; a `0x3db` gate +4; the `0x44c` block ±28; a
`0x456` gate 2 or 4) — reproduced on 154 of 154 type-6 records in both games,
closing what this survey previously reported as a partial, code-verified
contributor rather than a full reconciliation. `R2-ASSET-024`.

**Extents: ROM2 stays within ROM1's own observed corpus maxima on W, H, object
and unit counts, and within an already-published, directly-read structural
capacity on player count.** Measured against both games' own preserved-root
corpora: `W`/`H` both cap at 256 in both corpora (not exceeded); `#objects`
ROM2 max 324 vs ROM1 max 478; `#units` ROM2 max 785 vs ROM1 max 1815 — ROM2
under ROM1's own maximum on both (ROM1's own maximum is carried by the EN
root's `Horror.alm` alone: the RU root's own copy of that file declares the
identical metadata value with no backing type-6/type-4 record, being a
truncation of the EN file, `R2-ASSET-017`). `#players` ROM2 max 15 vs ROM1's
own observed corpus maximum of 9 — but ROM1's own type-5 record already
documents a fixed 16-slot diplomacy-row field (`ALM-GRP-041`, `AI-DIPLO-005`)
ROM1's own shipped corpus never filled past 9; ROM2's own type-5 record CONTENT
confirms the same capacity directly, not only by size match: case 5 reads four
leading fields (widths 4, 4, 4, 32) then a fixed 16-iteration loop (the loop
bound is the literal immediate `0x10`, not a per-file count) of 2-byte fields,
`4+4+4+32+16×2=76` — ROM1's own `ALM-GRP-041` arithmetic reproduced field for
field, not only matched by total size (`76*nPlayers`, `R2-ASSET-003`).
Separately: 0 of 83 ROM2 files pass ROM1's own documented header gate
(`formatVersion<=1001`) — every file fails on `formatVersion` alone, a direct
value comparison, not a claim about what ROM1's own loader would do with an
out-of-range `#players` value specifically. `R2-ASSET-027`.

## Not yet surveyed

Field-level MEANING of the 3 additional record types' own per-element bytes
(13 vs ROM1's 10) — `R2-ASSET-023` now places their exact widths, count-field
offsets and destination collections, not only dispatch mechanism, but what
each byte represents is undecoded; any record's own field content past its
outer 20-byte header, for any of the 83 maps (types 4, 7, 8 and 9 have no
published ROM1 per-type size formula to replay; every other type now does);
which specific bytes of record 0's 28-byte true growth (`R2-ASSET-024`) the
file's own declared `payloadSize` (644, `R2-ASSET-002`) counts, versus which
land in the 16-byte undeclared overhang — the total is accounted for
field-by-field (`R2-ASSET-024`), the split between declared count and overhang
is not, and is a save-path question this survey's load-path-only Ghidra pass
cannot answer; whether ROM1's loader enforces any hard reject on `#players`
above 9 specifically, as opposed to the type-5 record's own 16-slot field
capacity; the ROM2 root binaries other than `allods2.exe` and
`ROM2 Map Editor.exe`, and needles beyond the `templates.bin`/`Data.bin`
families, for either string-search result on `formats/rom2-databin/format.md`.
