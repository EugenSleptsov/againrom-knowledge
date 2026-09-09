<a id="fame-hall-of-fame-famehalldat--specification-core"></a>

# FAME hall of fame (`famehall.dat`)

The file stores a count followed by variable-length names and three words per
record. The ordinary reader preserves stored order and both opaque tail
words. Insertion and default seeding are separate operations. The tails are
not padding or required-zero validation fields. — FAME-READER-009,
FAME-INSERT-011, FAME-TAILS-013, FAME-SEED-015

## Wire layout

All integer fields occupy four little-endian bytes. The file has no magic,
version, archive wrapper or compression. — FAME-HDR-001, FAME-REC-002,
FAME-WRITE-007

```text
[count:32]
repeat count times:
    [nameSpan:32][name bytes:nameSpan][score:32][tail1:32][tail2:32]
```

The writer emits `4 + sum(16 + nameSpan)` bytes. Exact end-of-file after the
records is a writer and shipped-file property; the reader does not check for
trailing bytes. — FAME-REC-002, FAME-READER-009

| Field | Ordinary writer / producer | Reader / consumer | Claims |
|---|---|---|---|
| `count` | Current array count; not a constant ten | Resizes the record array and controls the record loop; zero clears it | FAME-HDR-001, FAME-READER-009 |
| `nameSpan` | CString stored byte length plus one | Number of bytes requested from the file; not a terminator-validation rule | FAME-STRING-010 |
| `name bytes` | CString bytes including the final NUL | Builds a CString by scanning from the local buffer to the first NUL; that string supplies display text | FAME-STRING-010, FAME-DISPLAY-012 |
| `score` | First record word; computed or seeded by the two located producers | Transferred unchanged; insertion compares signed 32-bit values and display formats with `%d` | FAME-INSERT-011, FAME-DISPLAY-012, FAME-PRODUCER-014, FAME-SEED-015 |
| `tail1` | Zero in both located producers | Read, copied and written unchanged; no direct read in the traced display body | FAME-TAILS-013 |
| `tail2` | Zero in both located producers | Same independent four-byte transfer as `tail1` | FAME-TAILS-013 |

In memory, each record is 16 bytes: a CString data pointer at `+0`, followed
by the three words at `+4/+8/+c`. The containing object has its array pointer
at `+134`, count at `+138`, and insertion limit at `+12c`. Its constructor sets
that limit to ten and initializes an empty array. — FAME-REC-002,
FAME-READER-009, FAME-INSERT-011

## Loading, strings and writing

The startup path opens `famehall.dat` for reading. An open failure or a file of
zero bytes takes the default-seeding path. A nonempty file goes to the reader:
a four-byte zero count produces an empty table, not seeded defaults. The
writer uses raw file writes; the previously traced exit handler opens the
file with create/write mode. — FAME-SEED-015, FAME-WRITE-007

The reader neither sorts the records nor clamps their count to ten. It reads
all three words without score-range or zero-tail tests. A loaded table can
therefore retain ties, nonzero tails and a different order until another
operation changes it. Signed loop/allocation arithmetic means that absence of
a guard is not a promise that every 32-bit count succeeds. — FAME-READER-009,
FAME-INSERT-011, FAME-BOUNDARY-016

The name buffer occupies 1,024 bytes. The reader passes the supplied prefix
directly to `Read`, then passes the buffer to a char-string constructor that
calls `lstrlenA`. There is no local ASCII, positive-length, length-bound or
final-NUL check. A well-formed ordinary name has one terminating NUL inside
that buffer. A 1,024-byte span with 1,023 non-NUL bytes followed by NUL fits the
observed buffer; it is not an enforced maximum-name validation rule.
— FAME-STRING-010

An early NUL shortens the in-memory string even though the full prefixed span
has been consumed; the writer then emits the shorter string. Non-ASCII byte
values are not rejected here. Zero-length or unterminated spans can reuse
previous buffer content; these instruction examples establish missing
validation, not a portable malformed-name contract. The effective native
encoding, upstream name-entry limits and rendering remain separate questions.
— FAME-STRING-010, FAME-BOUNDARY-016

## Insertion and display

Insertion walks the stored order and inserts before the first existing score
less than or equal to the new score, using signed comparisons. A new equal
score precedes the old equal score. If no such position exists, it appends.
The operation copies whole records and trims the array to its configured
limit after insertion. It does not deduplicate names or repair an already
unsorted table. The normal constructor's limit is ten; the load path has no
corresponding clamp. — FAME-INSERT-011

The display initializer copies the current record count. The display body
walks the array in its stored order, formats a one-based rank with `%d.`, passes
the record name to drawing, and formats the score with `%d` before a further
number-formatting helper. Neither trailing word is read directly by that
body. The helper's final text transformation and native rendered pixels are
outside this contract. — FAME-DISPLAY-012, FAME-TAILS-013

<a id="located-producers"></a>

## Score producers

The non-default producer takes its name from a char buffer at a live source
object's `+e4`. With signed words `A = campaign+124`, `B = campaign+128` and
`C = source+108`, it computes `C / (A * 10.0) * B` when A is nonzero, or
`C * stored_binary64(2e-6) * B` otherwise, using the observed x87 instruction
order. A helper truncates to a signed 64-bit integer and the producer stores
the low 32 bits as score. Both other words remain zero. These offsets identify
immediate inputs; their full upstream meaning and native floating-point
control state are Unknown. The expression is not a promise of exact rational
rounding at integer boundaries. — FAME-PRODUCER-014

The fallback producer obtains names from UI string-table indices 263 through
272. For the first nine it starts a score base at 70000, subtracts 7000 per
row and adds a random-derived remainder modulo 5000. It assigns zero score
to the last row. Both trailing words stay zero and all ten records pass
through the same insertion routine. — FAME-DEFAULT-008,
FAME-SEED-015

<a id="sample-facts-and-remaining-boundaries"></a>

## Unknowns

Names are not restricted to ASCII and loaded scores need not be strictly
descending. Those earlier clauses of FAME-NAME-003 and FAME-SCORE-004 are
partially retracted. Installed content does not establish prior-play history.
— FAME-NAME-003, FAME-SCORE-004, FAME-DEFAULT-006, FAME-DEFAULT-008

Allocation failure, short reads, overflowing lengths/counts, exceptions and
string lifetime are outside the ordinary successful-transfer contract.
The effective name encoding, upstream limits and final display transformation
remain Unknown. Neither tail word has an established level, mission,
difficulty or time meaning; a full native nonzero-tail lifecycle and complete
score-production session remain unverified. — FAME-TAILS-013, FAME-BOUNDARY-016

## Read and write sequence

1. Read u32 count and process that many records in stored order.
2. Read u32 nameSpan and that many name bytes. The ordinary CString value
   ends at the first NUL; the complete span is still consumed.
3. Read score, tail1 and tail2 as independent four-byte words.
4. Retain order, ties and tails. Loading does not sort, trim or deduplicate.

To write, emit the current count, then each CString's byte length plus one,
its bytes and final NUL, and all three stored words. The output size is
`4 + sum(16 + nameSpan)`. Use zero tails only for the two specified record
producers; preserve them in ordinary record transfers. — FAME-REC-002,
FAME-WRITE-007, FAME-STRING-010, FAME-TAILS-013
