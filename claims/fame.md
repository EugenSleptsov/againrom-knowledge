# Claim registry — FAME (famehall.dat)

Level 2 ledger. Index: [registry.md](registry.md) · spec:
[`formats/fame/format.md`](../formats/fame/format.md).
IDs are permanent. `✔ promoted` means the claim is reflected in the spec;
`● active (amended)` and `● partially retracted` require reading the correction
scope in [retracted.md](retracted.md).

The shipped EN/RU copies contain one distinct table. FAME-HDR-001 through
FAME-DEFAULT-008 retain that sample and the prior writer witness.
FAME-READER-009 through FAME-BOUNDARY-016 supply the current bounded lifecycle:
original static instructions, checked with synthetic instruction vectors under
explicit environment hooks. The older validator's six rejected mutants are
facts about that validator, not original-reader rejection evidence. Unknowns
inside the earlier writer-only row are superseded only to the extent of the
new field-specific rows; no global tail semantics or native score session is
claimed.

FAME-021 through FAME-025 add static upstream producers, conditional source
identity and a bounded startup precision request. Their new measurements decode
files only; no original instruction or process executes. Earlier upstream
Unknowns are superseded only where these rows supply an explicit receiver path.

## Hall of fame record and score lifecycle

| ID | Claim | Confidence | Status | Evidence |
|---|---|---|---|---|
| FAME-HDR-001 | A four-byte little-endian record count opens the file, followed immediately by records; no magic or separate version field. | High | ✔ promoted | [EXP-0012](../experiments/EXP-0012-famehall/), [EXP-0086](../experiments/EXP-0086-saves/), [EXP-0331 reader/writer](../experiments/EXP-0331-fame-record-lifecycle/RESULT.md) |
| FAME-REC-002 | Record framing is `[u32 nameSpan][nameSpan bytes][three 32-bit words]`, size 16+nameSpan; in memory it is a 16-byte record with CString pointer+0 and words+4/+8/+c. | High | ● active (amended) | [EXP-0012](../experiments/EXP-0012-famehall/), [EXP-0331](../experiments/EXP-0331-fame-record-lifecycle/RESULT.md), FAME-READER-009 |
| FAME-NAME-003 | All ten shipped names are NUL-terminated ASCII with prefix=strlen+1. | High | ● active (partially retracted) | [EXP-0012](../experiments/EXP-0012-famehall/), [EXP-0086](../experiments/EXP-0086-saves/), [EXP-0331](../experiments/EXP-0331-fame-record-lifecycle/RESULT.md), FAME-STRING-010 |
| FAME-SCORE-004 | The first 32-bit word is the record score used for insertion and decimal display. | High | ● active (partially retracted) | [EXP-0012](../experiments/EXP-0012-famehall/), [EXP-0331](../experiments/EXP-0331-fame-record-lifecycle/RESULT.md), FAME-INSERT-011, FAME-DISPLAY-012 |
| FAME-UNK-005 | The two trailing words contribute 80 zero bytes in the shipped ten-record table. | High / Unknown | ● active (amended) | [EXP-0012](../experiments/EXP-0012-famehall/), [EXP-0331](../experiments/EXP-0331-fame-record-lifecycle/RESULT.md), FAME-TAILS-013 |
| FAME-DEFAULT-006 | The sample is consistent with the shipped **default** table (10 seeded names, clean descending ladder ending at a `0` entry, all-zero trailing fields) — i.e. | Low | ● active | [EXP-0012](../experiments/EXP-0012-famehall/) |
| FAME-WRITE-007 | **`famehall.dat` is written by a raw `CFile`, on `WM_DESTROY`, and shares nothing with the save path.** `FUN_00489de0(this, CFile*)` loads `[[file]+0x40]` — the `CFile::Write` slot — into `EBP` once and calls it: first `obj+0x138` as 4 ... | High / Unknown | ● active | [EXP-0086](../experiments/EXP-0086-saves/) |
| FAME-DEFAULT-008 | EN/RU famehall.dat are byte-identical 228-byte files with count 10 and nameSpan=strlen+1 on all ten records. | High | ● active (amended) | [EXP-0143](../experiments/EXP-0143-actor-names/), [EXP-0331](../experiments/EXP-0331-fame-record-lifecycle/RESULT.md), FAME-SEED-015, FAME-DISPLAY-012 |
| FAME-READER-009 | Complete 00489ea0 reads count, resizes the object+134/+138 array, then reads each prefixed name and all three words through CFile slot+3c. | High / Medium / Unknown | ✔ promoted | [EXP-0331](../experiments/EXP-0331-fame-record-lifecycle/RESULT.md), evidence/reader.txt, writer.txt, array-size.txt, instruction-vectors.json |
| FAME-STRING-010 | Reader 0048a047..0048a05e passes nameSpan unchanged to Read into a 1024-byte stack buffer, then calls 572a2a; that helper calls imported lstrlenA before CString construction. | High / Medium / Unknown | ✔ promoted | [EXP-0331](../experiments/EXP-0331-fame-record-lifecycle/RESULT.md), evidence/reader.txt, string-from-char.txt, writer.txt, local-facts.json, instruction-vectors.json |
| FAME-INSERT-011 | Complete 00488a50 inserts before the first stored score<=new score using signed JGE at 00488ab3; otherwise appends. | High / Medium | ✔ promoted | [EXP-0331](../experiments/EXP-0331-fame-record-lifecycle/RESULT.md), evidence/insert.txt, record-copy.txt, campaign-constructor.txt, array-size.txt, instruction-vectors.json |
| FAME-DISPLAY-012 | Complete initializer 00455490 obtains campaign via frame+548 and copies campaign count+138 to display+74. | High / Unknown | ✔ promoted | [EXP-0331](../experiments/EXP-0331-fame-record-lifecycle/RESULT.md), evidence/display-init.txt, display.txt, references.json, local-facts.json |
| FAME-TAILS-013 | The two independently streamed words at record+8/+c are opaque carried values in the traced ordinary lifecycle. | High / Medium / Unknown | ✔ promoted | [EXP-0331](../experiments/EXP-0331-fame-record-lifecycle/RESULT.md), evidence/reader.txt, writer.txt, record-copy.txt, score-producer.txt, default-producer.txt, display.txt, instruction-vectors.json |
| FAME-PRODUCER-014 | Complete 00488c00 gets source from frame+d0 then server+3f54; if nonnull, copies source+e4 char buffer into the record name. | High / Medium / Unknown | ✔ promoted | [EXP-0331](../experiments/EXP-0331-fame-record-lifecycle/RESULT.md), evidence/score-producer.txt, score-caller-slice.txt, float-to-integer.txt, local-facts.json, instruction-vectors.json |
| FAME-SEED-015 | Startup slice 00471282..004712cc takes fallback 0048a740 when file-open returns 0 or file length is 0; nonempty input reaches 00489ea0. | High / Medium | ✔ promoted | [EXP-0331](../experiments/EXP-0331-fame-record-lifecycle/RESULT.md), evidence/startup-caller-slice.txt, file-length.txt, default-producer.txt, random.txt, instruction-vectors.json |
| FAME-BOUNDARY-016 | The measured population is one executable identity from EN/RU, one distinct shipped FAME table, eight complete lifecycle bodies, three caller slices and 12 helper entries. | High / Unknown | ✔ promoted | [EXP-0331](../experiments/EXP-0331-fame-record-lifecycle/RESULT.md), evidence/inputs.json, references.json, local-facts.json, probe.py, instruction-vectors.json |
| FAME-021 | Campaign+124 (A) is zeroed by constructor00487190 and reset00487640. | High / Unknown | ✔ promoted | [EXP-0359](../experiments/EXP-0359-fame-score-sources/EXP-0359.md), selected operations/time_accumulation and campaign_save/load_scalars; SESS-TICK-004, SAV-890, SAV-892 |
| FAME-022 | Campaign+128 (B) is a received hostile corpse-transition counter on the selected client arm. | High / Unknown | ✔ promoted | [EXP-0359](../experiments/EXP-0359-fame-score-sources/EXP-0359.md), selected operations/hostile_transition_increment, tables.json, stage-model.tsv; ANIM-DEATH-007, HERO-REVIVE-068, UNIT-VPLAYER-022, SAV-598, SAV-892 |
| FAME-023 | On actor-state projection004e7de3, effective mask4 copies simulation actor+130 as an unconverted dword through append004f0e74. | High / Unknown | ✔ promoted | [EXP-0359](../experiments/EXP-0359-fame-score-sources/EXP-0359.md), selected operations/projection_xp_stage, append_dword, packet_payload_pack, client_xp_and_stage; HERO-XP-077, SAV-HEROXP-063, ANIM-MSG-005 |
| FAME-024 | FAME source is the cached client drawable at [[frame+d0]+3f54]. | High / Medium / Unknown | ✔ promoted | [EXP-0359](../experiments/EXP-0359-fame-score-sources/EXP-0359.md), selected operations/source_cache, preview_cache, preview_file_xp, preview_xp_zero, resume_projection, score_inputs; FAME-PRODUCER-014, SAV-890, SAV-948, SAV-HUMPROJECT-461 |
| FAME-025 | The PE entry00556a00 reaches call00558360 in its selected prefix; that dispatcher conditionally calls the stored pointer005c98f8, whose image value is00554310. | High / Unknown | ✔ promoted | [EXP-0359](../experiments/EXP-0359-fame-score-sources/EXP-0359.md), selected operations/startup_initializer_call, precision_request, control_merge, control_precision_mapping, integer_conversion; FAME-PRODUCER-014 |

### FAME-HDR-001

A four-byte little-endian record count opens the file, followed immediately by records; no magic or separate version field. The shipped table has count 10. Writer 00489de0 writes object+138 and loops over that variable; reader 00489ea0 consumes the same header. Ten is not a constant file-count requirement.

**Confidence.** High for framing and the variable count; sample count is a whole-file measurement

### FAME-REC-002

Record framing is `[u32 nameSpan][nameSpan bytes][three 32-bit words]`, size 16+nameSpan; in memory it is a 16-byte record with CString pointer+0 and words+4/+8/+c. Reader and writer transfer these fields in this order. The shipped ten records tile 228 bytes and the writer emits only its record stream. Exact EOF is not a reader acceptance check: trailing input is ignored after count records. The original six rejected mutants tested the repository validator only.

**Confidence.** High for reader/writer framing and finite sample arithmetic; no malformed-I/O guarantee

**Amended.** The complete reader consumes the declared records and returns without comparing its cursor to file length. A synthetic trailing suffix remains unread. The writer's framing, exact shipped-file tiling and the old validator's finite mutant results stand.

### FAME-NAME-003

All ten shipped names are NUL-terminated ASCII with prefix=strlen+1. The writer loads CString nDataLength from [data-8], adds 1 and writes that many bytes; that is its arithmetic, not a reader validation rule. The reader consumes the prefixed span and then scans to the first NUL. ASCII-only and prefix=strlen+1 are withdrawn as universal input constraints.

**Confidence.** High for the sample, writer length arithmetic and corrected reader distinction

**Amended.** The writer uses CString stored length plus one. The reader consumes the prefixed span, then separately calls a NUL-scanning constructor; it does not check ASCII, a final NUL or equality with the scanned length. Early NUL and non-ASCII synthetic controls distinguish those rules. Shipped-name observations stand.

### FAME-SCORE-004

The first 32-bit word is the record score used for insertion and decimal display. The shipped scores descend strictly from 70006 to 0, with two above 65535. Strict descending order is a sample fact, not a format requirement: load preserves order and insertion allows ties, comparing signed 32-bit values. The width remains four bytes; unsigned ranking is not supported.

**Confidence.** High for width and the bounded insertion/display role; sample order is measured, not generalized

**Amended.** Load preserves stored order. Insertion uses signed comparison, puts a new equal score before an old equal score, and does not sort the existing table. Display uses the first word with `%d`. Four-byte width and the shipped table's strict order stand.

### FAME-UNK-005

The two trailing words contribute 80 zero bytes in the shipped ten-record table. Their broader meaning cannot be decoded from that sample. The earlier reserved/unknown label supplied no reserved-zero rule; the current ordinary transfer and producer roles are FAME-TAILS-013.

**Confidence.** High for the zero-byte measurement; Unknown for a semantic label beyond the traced lifecycle

**Amended.** Reader, copy and writer carry both words independently. Both located producers zero them; the complete display body reads neither directly. The sample measurement stands and no broader semantic label or universal zero rule is established.

### FAME-DEFAULT-006

The sample is consistent with the shipped **default** table (10 seeded names, clean descending ladder ending at a `0` entry, all-zero trailing fields) — i.e. not yet modified by play. Context, not a structural claim

**Confidence.** Low

### FAME-WRITE-007

**`famehall.dat` is written by a raw `CFile`, on `WM_DESTROY`, and shares nothing with the save path.** `FUN_00489de0(this, CFile*)` loads `[[file]+0x40]` — the `CFile::Write` slot — into `EBP` once and calls it: first `obj+0x138` as 4 bytes (the count), then per record `strlen+1` as 4 bytes, the string that many bytes, and `+0x04`/`+0x08`/`+0x0c` as 4 bytes each, striding `0x10` through the array at `obj+0x134`. **No `Asg&` header, no `CArchive`, no class record, no tag, no compression** — every one of which the save writer `FUN_004cee8f` uses. Its **only** caller is `FUN_00471ee0` (`EnumRefs callto:489de0`: 1 hit, 1 owner, 0 orphan), which opens the file `CFile::Open("famehall.dat", 0x1001, NULL)` = `modeCreate\|modeWrite` and is itself reached only through `0x00599bec` (`EnumRefs callto:471ee0`: 1 hit, 0 owners, 0 orphan) — a `pfn` slot in the six-dword message-map entry array at `0x00599bd8` whose `nMessage` dword is **2**, `WM_DESTROY`. So “written in the same minute as the last save” is the application exiting, not a shared writer. **And the round trip is byte-exact, twice**: the game rewrote the file at 00:24 and again at 00:40 with fresh mtimes, and both results hash `1845f772…1cbe9d`, **identical to the pristine EN root and to the pristine RU root** (the 00:40 copy was compared with `cmp` — 0 differing offsets over 228 bytes — and then discarded, per the corpus manifest's own addendum). That addendum asks what the game did at 00:40 that rewrites the hall of fame without changing it; this row answers it — **it exited** — a read‑modify‑write reproducing its own input, which is the test a single shipped sample could never give

**Confidence.** **High** for the writer, the mode, the trigger and the record arithmetic (one routine read whole; the `CFile::Write` slot is loaded once and used for every field; both caller sets are `EnumRefs callto:` with 0 orphan; the message-map record is read out of `.rdata`) / **High** that the shipped bytes survive a round trip (three sha256 over three copies, one of them written by the game) / **Unknown** what the three per-record dwords mean beyond `FAME-SCORE-004`'s width — the writer names nothing, and nothing in this session changed the table

### FAME-DEFAULT-008

EN/RU famehall.dat are byte-identical 228-byte files with count 10 and nameSpan=strlen+1 on all ten records. Their names equal UI string-table indices 263..272, unchanged between those localized tables. These are two copies of one distinct content sample. The former choice between file seeding and independent UI-table drawing is now resolved on the located paths: fallback 0048a740 populates records from those indices, while display 00455af0 reads record names. This does not establish the exact RNG state or history of the shipped table.

**Confidence.** High for whole-file identity, the retained source-table correspondence and the bounded source/display resolution

**Amended.** The fallback body constructs records from UI indices 263..272; the complete display body takes names from those records. Missing/zero-byte files select fallback, while a nonempty zero-count file loads an empty table. No exact shipped RNG state or native reseed witness is claimed.

### FAME-READER-009

Complete 00489ea0 reads count, resizes the object+134/+138 array, then reads each prefixed name and all three words through CFile slot+3c. Count 0 clears the array. No local sorting, ten-record clamp, score/tail validation or EOF check occurs. Signed loop/allocation arithmetic and unchecked reads prevent interpreting this as acceptance of every 32-bit count. Thirteen records, unsorted signed extremes/ties and nonzero tails survive the bounded instruction reader/writer vectors unchanged; extra suffix bytes remain unread.

**Confidence.** High for complete local control/data flow; Medium for conditional instruction execution as evidence beyond that static body; exceptional inputs Unknown

### FAME-STRING-010

Reader 0048a047..0048a05e passes nameSpan unchanged to Read into a 1024-byte stack buffer, then calls 572a2a; that helper calls imported lstrlenA before CString construction. No local positive-length, maximum, ASCII or terminator check exists. Early NUL normalizes the next write; non-ASCII bytes are copied. A 1024-byte NUL-terminated span fits. Synthetic unterminated/zero-length second records reuse prior buffer bytes, establishing missing validation under the probe environment rather than universal malformed-file behavior. Writer length is CString nDataLength+1.

**Confidence.** High for prefix-versus-scan and absent local checks; Medium for buffer-reuse instruction examples; native encoding, upstream name limits and failures Unknown

### FAME-INSERT-011

Complete 00488a50 inserts before the first stored score<=new score using signed JGE at 00488ab3; otherwise appends. Ties put the new record first. Copies/shifts preserve the complete 16-byte record, and the final comparison trims to object+12c (constructor default 10). There is no name deduplication or full sorting of old rows. The pointer/record-copy and comparison mechanics are read; vector outcomes assume explicit allocator/string/copy hooks and normal return from untraced helper 0048c0f0 on the newly appended empty slot before memmove.

**Confidence.** High for comparison, direct copy stores and limit branches; Medium for complete successful insertion across the untraced helper boundary

### FAME-DISPLAY-012

Complete initializer 00455490 obtains campaign via frame+548 and copies campaign count+138 to display+74. Complete body 00455af0 traverses 16-byte records in stored order, formats one-based rank with `%d.`, passes record name+0 to drawing, and formats record+4 with `%d` before helper 00468f60. It reads neither record+8 nor+c directly. Stored body pointers occur at 005990f8/005990a4. The final number helper, draw callbacks and native pixels were not traced or executed.

**Confidence.** High for initializer/body field uses and literal formats; downstream text transformation and full UI lifecycle Unknown

### FAME-TAILS-013

The two independently streamed words at record+8/+c are opaque carried values in the traced ordinary lifecycle. Reader 00489ea0, writer 00489de0 and record-copy 0048c230 preserve each; both found producers 00488c00/0048a740 initialize each to 0 and do not change either before insertion. Complete display 00455af0 reads neither directly. Nonzero synthetic values survive read/write and surviving insertion records under the documented environment hooks. These are not padding, a proved reserved-zero constraint, or proved globally unused storage. No semantic level/mission/difficulty/time interpretation or other nonzero producer is established.

**Confidence.** High for the named bodies and independent transfer/zero stores; Medium for hooked insertion survival; Unknown for other consumers, native nonzero lifecycle and broader purpose

### FAME-PRODUCER-014

Complete 00488c00 gets source from frame+d0 then server+3f54; if nonnull, copies source+e4 char buffer into the record name. With signed A=campaign+124, B=campaign+128, C=source+108, its x87 operations compute C/(A*stored 10.0)*B if A!=0, else C*stored_binary64(2e-6)*B. Helper 0055458c selects truncation for signed 64-bit FISTP; low EAX becomes the first word. Both remaining words retain 0. The sole located direct call is 00473ec5, inside the retained caller slice; broader scheduling and upstream meanings are not resolved. Five instruction examples declare FPCW 037f; the A 0/B3/C1000000 example gives 5, so exact-rational rounding or native initial FPU state must not be inferred.

**Confidence.** High for immediate source offsets, local arithmetic, conversion and zero stores; Medium for conditional numeric examples; upstream meaning, overflow/exception and native FPU state Unknown

### FAME-SEED-015

Startup slice 00471282..004712cc takes fallback 0048a740 when file-open returns 0 or file length is 0; nonempty input reaches 00489ea0. Complete seed body uses UI indices 263..271 with base 70000 decreasing by 7000, plus a random-derived signed remainder modulo 5000; index 272 gets score 0. All ten names are copied into records, both trailing words remain 0, and all records use 00488a50. The random helper returns 0..32767; exact shipped RNG state is not established. Synthetic startup vectors select ten seeds for missing/zero-byte files and zero records for a four-byte zero header.

**Confidence.** High for branch predicates, source indices and producer stores; Medium for startup with file/allocator environment hooks; no native deletion/reseed witness

### FAME-BOUNDARY-016

The measured population is one executable identity from EN/RU, one distinct shipped FAME table, eight complete lifecycle bodies, three caller slices and 12 helper entries. Whole-text E8/E9 byte-form and all-section literal-DWORD censuses find three insert call sites in two producer bodies, but do not exclude computed calls or aliases. The 26 instruction vectors execute original local bytes with memory-only hooks, including untraced0048c0f0, storage/lifetime, file, allocator and memory-copy boundaries. They establish conditional discriminators, not a native process, OS I/O, UI, exception or universal malformed-input contract. No owner artifact or original process was required.

**Confidence.** High for the measured population and method boundary; unresolved behavior explicitly Unknown

### FAME-021

Campaign+124 (A) is zeroed by constructor00487190 and reset00487640. Frontend message425 calls that reset with frame+548. In the admitted message41d completion arm (mode other than 0/1/3, frame+3bc nonzero), 00473e4a..00473e69 adds signed trunc(N/16), N=[[005cd758]+4], to the current A with ordinary 32-bit addition. SESS-TICK-004 identifies N as simulation sub-ticks; this is accumulation in groups of 16 sub-ticks, not proved wall-clock seconds or exact equality to the separate full-tick counter. The addition precedes SAV-890's final-score predicate/call. Campaign SAVE/LOAD transfer A unchanged as four bytes; no local clamp or positive-value test constrains the loaded word or this addition.

**Confidence.** High for named stores, signed division and admitted ordering; global writer closure, mission-entry clock baseline and native elapsed-time/session totals Unknown

### FAME-022

Campaign+128 (B) is a received hostile corpse-transition counter on the selected client arm. Dispatcher004104e8 saves drawable+15a before mask8 applies the incoming stage. The table00418873 sends new stages2/3/4 to00412209; only old unsigned stage<2 and relation word mask1 admit add1/store at00412264, through frame+548. The relation is local CPlayer+38 indexed by drawable owner CPlayer+4 (UNIT-VPLAYER-022). Stage1 is the death/fall arm; stage5 has another arm (ANIM-DEATH-007). This increment checks no attacker, killer or prior-ID ledger. Constructor/reset zero B; LOAD restores its raw word. CUnit constructor0045ae30 sets stage0, so a newly admitted CUnit receiving2/3/4 can satisfy the arm; repeated2-to3 does not. A later reset below 2 can permit another increment. No local increment clamp exists.

**Confidence.** High for the conditional local transition and counter stores; Unknown for native reconstruction ordering, initial-corpse double counting, global totals and a once-per-death interpretation

### FAME-023

On actor-state projection004e7de3, effective mask4 copies simulation actor+130 as an unconverted dword through append004f0e74. The state-packet constructor binds vtable0059c320 slot8 to packer004f0c9e, which selects6e/6f/70/6c by mask width and passes the current payload bytes to transport packing. Client004104e8's common state arm decodes that flag4 dword to005cd6b0 and stores it at drawable+108 (0041181c); without flag4 the copy is skipped. This raw transfer has no established Human-class prerequisite. On the Human/Humanoid source path established by HERO-XP-077 and SAV-HEROXP-063, the field carries the independently stored six-slot experience aggregate; the Human LOAD restores its raw signed aggregate without local recomputation from saved skill slots. Its six-slot meaning for other actor classes or an unproved final cached source remains Unknown. Transport success, object continuity and final delivery order are conditional.

**Confidence.** High for the original producer/packet/client field copies and the qualified Human/Humanoid experience meaning; other-class six-slot interpretation, final cached class and native latest-value delivery Unknown

### FAME-024

FAME source is the cached client drawable at [[frame+d0]+3f54]. State reception stores the newly constructed drawable there when its map count+9c4 equals 1 (0041171b..0041178c); the assignment has no ownership/selection predicate. Preview paths also allocate CUnit and assign this cache; character-file reading transfers16 bytes at source+fc, including+108, while temporary derive0047bf40 explicitly zeros source+108. These are distinct source producers, not a universal fresh-game C=0 guarantee. Resume setup requests all-state projection of the then-current Player+34 (SAV-948), but initial packet ordering/cache identity and the last delivered experience/corpse update before final00488c00 are unproved. The final call uses current cached fields after the completion-time addition; an intervening source virtual call in004201c9 remains a side-effect boundary.

**Confidence.** High for selected cache assignments and local final use; Medium for the conditional cached-actor lifecycle; native main-hero identity, preview replacement and latest-value timing Unknown

### FAME-025

The PE entry00556a00 reaches call00558360 in its selected prefix; that dispatcher conditionally calls the stored pointer005c98f8, whose image value is00554310. After returning helpers, that initializer calls0055c180, which requests value10000 under mask30000 through005648d0/00564890. The conversion pair maps this request to x87 precision bits0200 (53-bit) while preserving incoming rounding bits. This is a conditional local startup precision request, not a measured full control word or proof it survives to FAME. Separate FNINIT005648f8 and other control setters have no established path to the score call here. The score still follows FAME-PRODUCER-014's signed A/B/C loads and x87 order;0055458c temporarily forces truncation for signed64 FISTP and restores prior control, with low EAX stored as score. No native rounding, exception/overflow outcome or safe gameplay maximum is established.

**Confidence.** High for selected startup/control mapping and conversion instructions; native initial/final FPCW, exceptional conversions and complete earned-score session Unknown


## Earned score source and campaign ending

| ID | Claim | Confidence | Status | Evidence |
|---|---|---|---|---|
| FAME-026 | The FAME source is one cached drawable: the one registered when the client map was empty, not a chosen hero. Its experience word arrives only through mask-4 state packets that pass the sender's owner filter. | High / Medium / Unknown | ✔ promoted | [EXP-0448](../experiments/EXP-0448-score-ending/EXP-0448.md) |
| FAME-027 | No game-code instruction writes the x87 control word, but game code calls C runtime entries that do. The word at the score call is Unknown; on a 33 x 33 x 1199 grid 53-bit and 64-bit precision differ on 1471 points. | High / Unknown | ✔ promoted | [EXP-0448](../experiments/EXP-0448-score-ending/EXP-0448.md) |
| FAME-028 | At five of 11 world-step call sites the client pump runs adjacent to the step. Resume projects Player+34, the unfiltered global list and world-list actors with actor+13c below 5, with no stage baseline. | High / Medium / Unknown | ✔ promoted | [EXP-0448](../experiments/EXP-0448-score-ending/EXP-0448.md) |
| FAME-029 | Credits end on a key press, a mouse press (by slot shape) or the scroll end; the hall of fame ends on its button. Message 442 is the only SAVE caller not excluded by SAV-972; its reachability is Unknown. | High / Medium / Unknown | ✔ promoted | [EXP-0448](../experiments/EXP-0448-score-ending/EXP-0448.md); SAV-972 |

### FAME-026

Instrument: capstone decode of the matched EN/RU `rom.exe` (SHA-256 in `evidence/inputs.tsv`), listings in `evidence/listings.txt`, address assertions in `evidence/checks.tsv`.

Cache assignment. The cached drawable is the one registered when the map held one element, which is the first registration on an empty map. Receiver `004104e8` stores the drawable it just built or found into `[frame+d0]+3f54` at `0041178c`. The store runs only in the new-object arm, after the drawable is registered, when the map element count at `frame+9c4` equals 1 (`0041176b`, `00411777`). The count is the element count of the id-keyed map at `frame+9b8`. The store has no class, owner or hero test. The score producer `00488c00` reads only this pointer (`00488c57`) and loads `source+108` (`00488c7d`); it reads neither `Player+34` nor any other drawable. The score therefore uses one actor's experience aggregate, not a party sum.

Writers of the cache word, from the census of all decoded `.text` (`evidence/cache-field-writers.tsv`, 9 stores): zero at `00402c3e`, `00410401`, `00420e68`; the receiver store at `0041178c`; zero at `00475b71` after the map is cleared; zero at `00476a87` in the terminal reset route (FAME-029); and three preview-object stores at `0047b048`, `0047b364`, `0047b7bb` (FAME-024).

Experience delivery. `source+108` is written only at `0041181c`, from the flag-4 dword of a state packet. The sender `004e7de3` builds that dword from `actor+130` (`004e7f7c`). Its owner filter removes bit 4 by `and 0x507b` (`004e7e9e`) when the actor's owner `actor+14` differs from the recipient and the actor's virtual slot `+30` returns 0. When that slot returns nonzero, it removes bit 4 by `and 0x50fb` (`004e7f11`) unless the kind word `actor+0e` lies in `0x21..0x3f` (`004e7efb`, `004e7f0c`). An actor owned by the recipient is never filtered. The gain routine `004f72d7` returns without effect unless `actor+0e` lies in `0x21..0x3f` (`004f72eb`, `004f72fc`) and, after updating `+130`, sends mask 4 (`004f7725`) or `0x1f001304` on a level change (`004f7705`) with recipient 0, which `004e7de3` expands to every eligible player.

Which actor is first. Entry routine `004d303e` (mission entry and resume, SAV-948) projects `Player+34` with mask -1 at `004d3470`, then calls `004e9b4e` (`004d347e`), which projects `Player+34` again before every other actor (`004e9b70`, FAME-028). `Player+34` is the participant's own named hero (SAV-HERO-059). The actor-add announcer `004e9312` also sends mask -1 to human players (`004e93a7`, six callers in `evidence/calls.tsv`); whether an announcement through `004e9312` reaches the empty client map before `004d3470` was not traced, and is the live alternative to the hero being the cached drawable.

**Confidence.** High for the cache assignment and its predicate, the full writer census, the single reader and field, the sender's filter masks and the gain gate, each read at instruction level. Medium for the hero being the cached drawable: it follows from the order inside two routines, with no run, and the announcer route `004e9312` is untraced. Unknown: whether the local transport delivers every packet in send order, whether an earlier announcement can precede the hero projection, and the meaning of the virtual slot `+30`.

### FAME-027

Instrument: capstone linear sweep of `.text` (no control-flow recovery); control-word instruction census `evidence/control-census.tsv` and `evidence/control-by-routine.tsv`; raw-byte cross-check `evidence/control-rawbytes.tsv`; caller tables `evidence/calls.tsv` and `evidence/dword-refs.tsv`; imports `evidence/imports.tsv`; arithmetic model `evidence/arith-probe.tsv`.

Inventory. The sweep decodes 71 x87 control or environment instructions (FLDCW, FNSTCW, FNINIT, FNCLEX, FRSTOR, FLDENV, FNSAVE, FNSTENV, FXRSTOR, LDMXCSR, XRSTOR). The sweep attributes each to the preceding `push ebp; mov ebp,esp` (column `prologue_attributed`); frameless library routines have no such prologue, so the grouping in `control-by-routine.tsv` is attribution and not routine boundaries. One, at `0043ae0a`, is a jump-table dword decoded as an instruction (`evidence/data-words.tsv`). The other 70 lie at or above `00553f70`, the C runtime range. The raw-byte scan of `.text` for memory FLDCW, FLDENV, FRSTOR, FXRSTOR, LDMXCSR, XRSTOR and FNINIT byte patterns finds 97 matches, 45 below `00553f70`: 44 lie inside other instructions and one is the jump-table data above. No game-code instruction below `00553f70` writes the control word.

Game code calls runtime entries that do. `00554654` (called from `0040a521` and `0040a88d`) reaches `fldcw [005a0a68]` at `0055466c`. `00557374` (called from `004e9067`) and `005574a4` (12 sites including `004f7261` and `004f7276`) reach the reloads of the same constant at `00557390` and `005574c0`. The constant at `005a0a68` is `027f`: precision 53-bit, nearest rounding, all exceptions masked. Four of the 20 call sites of the merge helper `005624e0` request `0x133f` (64-bit precision, nearest rounding): `00556621`, `00567e1e`, `00567f41` and `00569f4e`, each preceded by `push 0x133f` (`push_133f` rows). The other 16 push a saved word; whether they restore it was not audited.

Persistent setter. `00564890` merges a request into the hardware word as `(value & mask) | (old & ~mask)` (`005648ab`..`005648b3`). Its only direct caller is wrapper `005648d0`, whose only caller is `0055c180` (`push 0x30000; push 0x10000`, `0055c180`, `0055c185`). `0055c180` is called from the startup initializer `00554310` (`0055431f`) and from the FNINIT entry `005648f0` (`005648fa`). Entry `005648f0` has no direct call and no dword reference in any section. Value 0x10000 under mask 0x30000 requests precision 53-bit (FAME-025). The mask leaves rounding and exception bits as inherited. Conversion helper `0055458c` saves the word (`00554593`), sets rounding bits 11 by `or ah, 0xc` (`0055459b`), runs FISTP to a 64-bit slot and reloads the saved word (`005545a8`). The producer stores the low dword.

Probe. The model rounds each x87 result to a 53-bit or 64-bit significand in the order of `00488c71`..`00488cb6` (FILD B, FILD C, A times 10.0, divide, multiply by B; or C times binary64 2e-6 times B when A is 0) and truncates. The grid is A 0..32, B 0..32 and the 1199-value C list `range(0,3001,3)` plus `range(3000,200001,997)` (3000 appears twice): 1,305,711 points. On this grid 53-bit and 64-bit precision give different scores on 1471 points, 53-bit differs from the exact rational on 648 and 64-bit on 833. Example A=1, B=15, C=42: exact 63, 53-bit 63, 64-bit 62. The real domain of A, B and C in play is not established, so 1471 is a count on this grid and not a frequency.

**Confidence.** High for the instruction-level absence of control-word writes in game code below `00553f70`, the game-to-runtime calls listed, the setter's merge formula and call graph, and the conversion helper. Unknown for the control word at the score call: the runtime entries above set `027f` or `0x133f`, the 16 saved-word callers were not audited for restoring it, and the model assumes ideal x87 rounding without exponent effects. Unknown: the control word at process start (the loader and the 14 imported DLLs listed in `evidence/imports.tsv`, ADVAPI32, COMCTL32, DDRAW, DSOUND, GDI32, KERNEL32, SHELL32, USER32, WINMM, WINSPOOL, WSOCK32, comdlg32, ole32 and smackw32, run code outside this image), rounding and exception bits at the call, and any entry into `005648f0`.

### FAME-028

Instrument: capstone decode; listings of the pump, step and resume bodies; `evidence/calls.tsv`.

Order. The client pump `004104e8` polls the queue through `004e7670` in a loop (`0041056a`, exit at `0041868e`/`004186a6`). A packet whose opcode byte equals the argument ends the call with 1 (`004186a2`); with an empty queue the call returns 0 when the argument is 0 and bit 0 of `frame+3dc` is clear (`004186b2`..`004186cc`), and otherwise polls again (`004186ce`). The world-step routines are called directly from 11 sites (`004d88f1` at `00410433`, `0041d47c`, `004716ef`, `00474010`, `00478bd2`, `00478c0a`, `0047bdcc`; `004d2551` at `00475293`, `00475415`, `004755ab`, `00477e5b`). The pump is called right after the step at five of them (step, then pump): `0041d47c`/`0041d489`, `004716ef`/`004716fc`, `00478bd2`/`00478bdf` (world step `004d88f1`) and `00475293`/`004752a0`, `00475415`/`00475422` (world step `004d2551`). The step body `004d891a` increments the world tick, runs actor subticks and ends with the queue flush `004e9e7c`; actor state packets, including stage mask 8, are produced inside it. The client reads the old stage at `00411828` before applying a mask-8 stage at `0041184f`, so a stage packet produced by step N is counted by the pump that follows step N. The other pump sites `004104b2`, `0047569b`, `00477f01`, `0047982a` and `0047be5e` follow a timer wait or poll loop, not an adjacent step.

Loaded corpses. Resume `004e9b4e` projects `Player+34` first (`004e9b70`). It then walks the global actor list `[00609558]` (MOVE-TICK-013) and projects each actor with mask -1 and no stage test (`004e9bb0`, `004e9bbf`). It then walks the world list at `[[5cd758]+14]+c` and projects each actor whose `actor+13c` is below 5 (`004e9c1b`, `004e9c2f`). The mask carries the stage (bit 8 survives both filters). A drawable built by `0045ae30` starts at stage 0 (FAME-022), so a hostile corpse at stage 2, 3 or 4 meets the increment arm; the arm sends 0 to 2/3/4 to `00412209` with old stage below 2. No instruction of the resume bodies stores to the campaign counter (`004e9b4e`..`004e9c89` listing), and the campaign LOAD restores the saved raw word (FAME-022). So loaded corpses are neither skipped nor baselined by these bodies.

**Confidence.** High for the call order at the five sites, the pump's loop and exit, the resume projection loops and their guards, and the absence of a counter store in them. Medium for the counting consequence: it needs delivery to an empty client map, the hostile relation word and a saved counter that already includes the corpse. Unknown: whether the local transport returns a packet to the pump of the step that sent it, how many stage 5 actors stay on the actor list, native behaviour.

### FAME-029

Instrument: capstone decode; listings `q4_*`; `evidence/vtables.tsv`; `evidence/calls.tsv` (sites that push a message id).

Credits (SAV-970 screen). Activation `00438180` sets the scroll offset at `+94` to 0x1e0 (`00438183`). The draw routine `00438280` runs a step only when `timeGetTime` has advanced by 23 ms or more since its stored value (`004382f7`); a step decrements the offset, and the first visible line index at `+98` becomes the absolute offset divided by the line height once the offset is negative (`00438306`..`00438317`). When `+98` reaches or passes the clamped last line index it sends 445 to itself through `00438850` (`0043863a`, `00438640`). Slots 27 (`00438250`, one argument, then the base handler `004c5369`) and 21 (`00438270`, three arguments) of the credits vtable also call `00438850` (`evidence/vtables.tsv`); by argument shape they are the key-down and mouse-press slots, the same slots 27 and 21 in which the hall-of-fame class (slots `00455890`, `00455990`) carries its input handling. Message 445 reaches `004c52f3`, which with `modal+5c` nonzero posts 44c (SAV-971). The run time is a function of the line count of the credits text and the line height; neither was read.

Hall of fame (FAME-DISPLAY-012 screen). Routine `00455250` (slot 30 of the hall vtable at base `00599078`, `evidence/vtables.tsv`) builds a button with command 0x445 and stores the rectangle `0x230,0x1a0,0x25c,0x1c0` at `+a8`. Mouse-down inside the rectangle records a pressed state (`00455990`); mouse-up inside it calls `00455cb0`, which sends 445 (`00455a60`, `00455ad6`, `00455cb6`). Key-down slot `00455890` calls only the base handler `004c5369`, which maps Tab, Enter and Escape to objects at `[this+38]+50/+54` (`004c53c7`..). The screen's own methods (`00455250`..`00455cc0`) contain no `timeGetTime` call.

Terminal route. After the close chain SAV-971 describes, `004769c0` calls campaign reset `00487640` (`004769e4`), clears the drawable map array (`00476a29`..`00476a67`) and zeroes the score cache word (`00476a87`).

Producers of 41d (`evidence/calls.tsv`, `push_41d`): `0043c700`, `0043c75e`, `004476c9` (dialog or panel constructors binding command 41d to a control), `00446d76` (after a 446 modal close), `0047407d` (frame code when the selected mission number is a multiple of 10), `00479561`, `004795dd`.

SAVE. Writer `00478c40` has two callers: `00475d50` (SAV-972) and `0047500c` in the arm for message 442. That arm copies text string 0x37 and the literal at `005be4e4` into `frame+134` and `frame+234` first (`00474fb1`..`00475008`), so message 442 saves from copied fixed text, with no test of `frame+3dc`; the content of that text was not read. Message 442 has two producers: command id 442 of a frame child constructed at `00472400`, and a PostMessage at `00478a19` after a mode-2 load handler.

**Confidence.** High for the 23 ms scroll gate, the end test, the 445 sends, the vtable slot contents, the terminal clear and the two SAVE callers. Medium for naming slots 21 and 27 as the mouse-press and key-down inputs (argument shape and the sibling class, not a dispatcher trace) and for the order credits then hall of fame then reset (SAV-971 scope). Unknown: which of the seven 41d producers is the Victory acknowledgment, the credits duration, whether Tab, Enter or Escape dismiss the hall of fame, whether the frame child with command 442 is visible and clickable while either screen is up, and base-class timers.
