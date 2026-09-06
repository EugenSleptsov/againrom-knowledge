# SAV save game (`Asg&`) — specification (partial)

## Equipment producer sequence

Stored equipment, modifier fields and Weapon-owned Spell state are not one
atomic update. Armor removal clears its actor slot before Effects; Shield
removal clears it after Effects; Weapon removes Effects before its own direct
fields/derive and deletes its Spell before the final slot clear. Attach has
class-specific order too: Weapon prepares its Spell before displacement and
derives before timing/range and weight. Armor/Shield store slots and add their
defensive blocks before positive-weight refresh, with flags before Effects for
Armor and after for Shield (`SAV-EQUIPORDER-552`, `HERO-EQUIP-017`,
`ITEM-ARMFOLD-033`).

Weapon byte `+50` survives serialization separately from definition binding.
The current byte, reread after derive, supplies the minus-one delta added to
or subtracted from actor range byte `+12c`; byte wrapping is retained, not
guaranteed inverse restoration (`SAV-EQUIPORDER-552`). Melee's active selector
instead takes the low byte of definition parameter5 cached after displacement
and the new-weapon slot store;
ranged11/12 and removal assign zero (`SAV-HUMEQUIP-447`, `HERO-EQUIP-017`).

State0 Effect apply/remove traverses forward and derives after each normal
general dispatch, before advancing. The iterator saves one next-node pointer
but reads that node's payload/link later; transitive mutation is not excluded
(`SAV-EQUIPEFFECT-553`). Command22 adds container operations and a final
zero-weight refresh. That refresh derives only on a changed quotient, and
Humanoid/Human wrappers add no trailing derive. The Unit wrapper's `+54` is
Token value, not derive (`SAV-EQUIPCALL-554`, `ITEM-EQUIP-006`). These named
intermediate consumers do not establish which state a real save observes,
callback purity, or a complete graph-wide write set (`SAV-EQUIPOBS-555`).

## City transaction producer boundary

Carried-to-table transfer updates container bookkeeping without an actor-load
refresh call in that dispatcher arm (`SAV-CITYMOVE-512`). Sale completion and
individual table-to-inventory return refresh actor load, but call Human derive
only when the signed16 load/capacity quotient changes; bulk cancellation uses
a different return loop (`SAV-CITYSALE-513`, `SAV-CITYRETURN-514`). A successful
school purchase updates skill/base/XP fields and unconditionally dispatches
derive (`SAV-CITYDERIVE-515`). These are different state-producing boundaries.

The selected writer emits existing Human spans `+a6/24`, `+be/22`, `+114/24`,
`+d4/64`, scalar words and XP, and stores the inventory's insertion index and
running load without recomputation. It also explicitly assigns Human
`u32 +148 = u8 +14c`. Thus raw-field preservation is not a claim that SAVE
never mutates any field. Embedded serializers and intervening callback
chronology remain Unknown; no universal next-SAVE or safe initial-state vector
follows (`SAV-CITYSTORE-516`).

Level 3. Promoted, evidence-backed claims only. Claim basis: container
(`SAV-HDR-001`, `SAV-VER-002`, `SAV-PTR-003`, `SAV-EMB-004`); body transport and extent
(`SAV-FRAME-021`, `SAV-CODEC-022`); decoded framing, objects and tables
(`SAV-STREAM-010`, `SAV-STREAM-013`, `SAV-OBJ-014`, `SAV-OBJ-016`, `SAV-ID-015`,
`SAV-CELLREC-017`,
`SAV-TAIL-018`); writer, reader and world `Serialize`
(`SAV-SHAPE-023`, `SAV-ROSTER-024`, `SAV-HEAD-025`, `SAV-TRAIL-026`); campaign
record (`SAV-CAMPTAIL-070`, `SAV-CAMPPROG-071`, `SAV-CAMPPOS-072`,
`SAV-CAMPMARK-073`, partially retracted);
terrain overlay, typed cell payload and post-load identity order
(`SAV-CELLLOAD-108`…`SAV-CELLLOAD-113`, `TRIG-CELLTAIL-035`); and the exact
byte-0-to-EOF reader and bounded witness census (`SAV-FULLREAD-252`, `SAV-READPOP-255`).
Writer completeness and its claim-insufficient boundary are recorded by `SAV-WRITERAUDIT-380`.
Original EN writer acceptance, the tested no-world minimum, resave canonicalization,
unsafe-value counterexample and mission-transition limit are recorded by
`SAV-ORIGWRITER-396`…`SAV-ORIGMISSION-400`.

**Status: partial (◐).** The container, the body's transport, the decoded stream's object framing
and its two cell-keyed tables are specified; every object instance is located and attributed. The
member layout of all eleven observed classes is read from its writer, though most fields are named
only by offset. `SAV-CELLLOAD-108`…`SAV-CELLLOAD-113` add the complete terrain/cell
load order and the 52-byte cell payload's typed lifecycle. Its static evidence covers both
byte-identical executable roots; its object join
covers 14 distinct preserved world-half saves. The first four corpus saves remain the stride and
static-bit-5 discriminator for the table itself.

The unified reader (`SAV-FULLREAD-252`) closes 55/55 current paths and 31/31 distinct
digests with zero physical or decoded interval gaps/overlaps. This is Medium-confidence
population closure: it establishes
framing, not meanings for opaque raw members or production/acceptance of unwitnessed classes.
`SAV-SUFF-300` reruns that reader and partitions the same population into **52 world-half paths / 29
digests** and **three no-world-half paths / two digests**. Those counts supersede the older 50/5
split; they remain a finite-corpus observation, not a rule for other producers. — SAV-SUFF-300

Seen as: `game####.sav` at the install root. Magic `41 73 67 26` = "Asg&".

## At a glance

A **16-byte header**, then one **compressed blob** that carries its own decompressed size, then an
uncompressed tail. The header's `0x0C` is the blob's byte length - what the reader allocates and
reads - and `0x04` is where the blob ends, patched in by the writer afterwards.

```
+----------+--------------------------------+------------------------------------+
| header   | compressed blob                | uncompressed tail                         |
| 16 B     | u32 outWords | run/literal code | 0x100 B label + &YA1 + campaign record    |
+----------+--------------------------------+------------------------------------+
0        0x10           0x14              @0x04              @0x04+0x100        EOF
                 (the blob is [0x10, @0x04) and @0x0C == @0x04 - 16, by construction)
```

## Header (16 bytes, little-endian)

| Off | Type | Field | Notes | Claim |
|-----|------|-------|-------|-------|
| 0x00 | char[4] | magic | `41 73 67 26` = "Asg&"; a mismatch reports *"Invalid save file."* | SAV-HDR-001, SAV-FRAME-021 |
| 0x04 | u32 | blobEnd | first byte **after** the blob; the uncompressed tail starts here. Written as 0, then patched with `Seek(4,0)` and the writer's `GetPosition()`. **The reader reads it and discards it** | SAV-FRAME-021 |
| 0x08 | u32 | version | `0x0BAD0002`. The reader requires `>=` this and otherwise reports *"Outdated save file."* | SAV-VER-002, SAV-FRAME-021 |
| 0x0C | u32 | blobBytes | the blob's byte length: the reader's allocation size and read length. `== blobEnd - 16` | SAV-FRAME-021 |

The blob begins at `0x10`. **There is no fifth header field**: the `u32` at `0x10` that earlier
rounds called `uncompressedWords` is the blob's *own* first dword, read by the decompressor - which
is why the same four bytes ride in a `.chr` character file and in a network packet, neither of
which has a container header at all.

## The two shapes

A save carries a **campaign half** always and a **world half** only sometimes. After the campaign
half the stream holds **one byte**: the writer emits `1` when the world snapshot follows and `0`
when it does not; the reader tests zero versus nonzero. The writer takes it from `world+0x2c`, set
when a map finishes loading. So a between-mission save and a mid-mission save are different shapes
and the file says which. The
**roster is in the campaign half** and is therefore in both. - SAV-SHAPE-023, SAV-ROSTER-024

```
campaign half : counters, map NAME, 11 zero dwords, mission no., difficulty,
                player list -> Player -> groups -> actor lists     (always)
shape byte    : writer 00/01; loader zero/nonzero                   (always)
world half    : world sub-objects, the map re-opened by name, terrain,
                block plane, cell records                           (only when 01)
trailer store : 0xBADFACE1, one dword, world+0x118's block          (always)
trailer load  : discriminator, conditional dword, world+0x118 block (always)
```

**Both shapes are now observed** (`SAV-SHAPE-023`, `SAV-CITY-030`; corpus 12 owner-produced saves): eleven carry `01` and
`game0010.sav` carries `00`. What the zero shape does *not* contain, each a world-half structure:
the block-plane array, the 54-byte cell-record table, the 4374-byte session block, and the class
records `Unit`, `Building`, `Sack`. What it *does* contain: the whole roster, and a campaign head
whose mission number is **0** while the map-name field still holds the previous mission's name —
so a consumer must take "is a mission in progress" from the shape byte, never from the map name.
— SAV-CITY-030

## The mission-outcome flag

`Player+0x3c` — the ninth field of the `Player` record, one byte, in the **campaign half**, and
therefore in **every** save including a between-mission one:

| Value | Meaning | Written at | Alongside |
|-------|---------|-----------|-----------|
| `0` | mission in progress | `004d8a7b`, and the `Player` ctor `004fab60`, and `004d31b2` | — |
| `1` | **mission complete** | `004d8cf4`, on `session+0xb3ac == 1` | log `"Logic - Mission Complete"`, announcement `0xb5` |
| `2` | mission failed | `004d8bcf` (no living hero) and `004d8c9d` (`session+0xb3b4 == 1`) | log `"Logic - Mission Failed"`, announcement `0xb4` |

The reporter `FUN_004d8963` runs once per full tick, tests **lose before win**, and acts only on
the player whose `+0x28` is `0` (the human participant) and whose `+0x3d` is non-zero.

**Do not use `session+0xb3ac` as the flag.** It is the win *counter* an authored action
increments, it lives in the world half, and it is gone the moment the mission ends — which is
exactly when a campaign needs to know the mission was won. — SAV-FLAG-027

## The `Player` record

Seventeen fields, one straight run, `FUN_00511089`:

```
CString  +0x18   name (length byte then chars, MFC short form)
u16      +0x04   1-based type-5 slot id
u32      +0x08   the same slot id again
8 bytes  +0x10   raw, undecoded
u8       +0x44   the ALM roster's own colour/palette-shade slot plus one, truncated to a byte
                 (004e2435 MOV EAX,[EDX+0x8]; 004e2438 ADD EAX,0x1; 004e243e MOV [ECX+0x44],AL
                 — the roster record's own +0x8, not copied verbatim), by the same ALM-to-session
                 copy that assigns +0x2c (ALM-GRP-041, PAL-SHADE-012, ALM-PLAYER-069). Correlates
                 1:1 with +0x2c's own nonzero
                 population on the preserved corpus: every record with +0x44==2 (the campaign
                 protagonist colour) has +0x2c==2, and no other value pair co-occurs — SAV-664
u32      +0x28   0 iff a human participant owns this player
u16      +0x2c   1 << (registration slot mod 16) for the participant's own slot at +0x04,
                 computed only when +0x28==0 (zero otherwise). The preserved corpus never
                 exercises a slot other than 1 (value 2) — SAV-664
u32      +0x38   XOR 0x5c073f4d  <- MONEY
u8       +0x3c   the mission-outcome latch
u8       +0x3d   retained-actor mission-entry placement latch
u32      +0x48   XOR 0x5c073f4d
u32      +0x50   0 in all 410 preserved records, matching the constructor default; no writer
                 located in a bounded search — SAV-665
u16      +0x54   dword in memory, min(v, 0x7fff) on the way out
u16      +0x4c   dword in memory, min(v, 0x7fff) on the way out
u32      +0x58   defaults to 95; overwritten to 0 or 50 for the same population +0x2c marks
                 nonzero, plus two zero-colour records landing on 0. The routine performing the
                 override was not located within one bounded address-range search — SAV-666
u32      +0x34   the participant's own STARTING CHARACTER: an identity key resolved through
                 the pointer map on load. Equal to exactly one actor's key in this Player's
                 own groups on 18/18 saves. Written at chargen and by the carry arm and NOT
                 by the mercenary spawn, which is what makes a mercenary not a hero
u32      this    the object's own pointer, written verbatim
```

The record does not end there: after the store/load merge `Player::Serialize` writes the
player's group set — `FUN_005102f4`: a `u32` group count then that many group records, each
`FUN_00511938` = the embedded `u16` list at `Group+0x20` (`WriteCount` then `2n`) +
`FUN_005391d0` on `*(Group+0x3c)` (80 raw bytes then that object's own `u16` list) + a `u32`
actor count with that many `ar << CObject*` + `u32 +0x1c` + `u32 +0x40` + `u32 +0x44` — and then 32 raw bytes
from `this` (`FUN_005392d0`). Of those 32 bytes, only `+0x1f` (`AI-FORM-037`'s formation mode) is
ever nonzero across all 410 preserved records; the other 31 are 0 with no exception, and the
archive call itself admits no per-field dispatch to find one — SAV-662, SAV-663.
**And then a `Diary`** — `00511394 CALL dword ptr [EDX + 0x8]` on
`*(Player+0x40)`, a 48-byte object built at `004fac36`/`004fac4f` whose constructor stores vtable
`0x59c4b8`. `Diary::Serialize` `00511427` writes a `CDWordArray` (`WriteCount` then `4n` raw), a
`CWordArray` (`WriteCount` then `2n` raw), and a `u32` resolved through the pointer map on load.
A freshly constructed `Diary` sizes both arrays once from a global catalog — `data.bin`'s Units
table (DAT-OBJ-002, ALM-CLS-038 as amended), so an array index is a unit type — through the same
`SetSize` entry point every subsequent load also reaches, resizing both arrays again from the
stream's own count; their element counts are always equal (119 on 540/542 preserved records, 0 on
the remaining 2, both the same project-written save and its unmodified resave). Each
pair satisfies `word[i] = 1024 - dword[i]` exactly, with the one located mutator touching only the
word side (the dword side's own writer is a bounded negative; only six of the 119 indices ever
carry a nonzero dword anywhere in the corpus). The trailing `u32` resolves, on every one of 408
nonzero-reference records among 410 `Player`-owned `Diary` records, to the Diary's own enclosing
`Player` and never any other object; a further 132 `Diary` records are actor-owned
(`Humanoid+0x1e4`) and all null — SAV-667, SAV-668, SAV-669.
Without that record a walk desynchronises inside the first `Player`. — SAV-MEMBER-036,
SAV-PLDIARY-054, SAV-HERO-059

Two traps a verbatim copy walks into. `+0x38` and `+0x48` pass through `FUN_00527e80`, which is
`return x ^ 0x5c073f4d` — an involution applied in both directions, so a reader that skips it sees
about 1.54 billion gold. `+0x54` and `+0x4c` are **silently saturated** at 32767 by the writer:
a larger value does not survive a save and nothing says so. — SAV-PLAYER-028, SAV-OBF-029

## The body's transport — a run/literal code over 16-bit words

Take the blob `[0x10, blobEnd)`. Its first `u32` is `outWords`; allocate `2 * outWords` bytes.
From `blob+4`, read one opcode byte `n`:

```
n <  0x80   ->  n literal 16-bit words follow  (2n bytes, copied out; consumes 2n+1)
n >= 0x80   ->  the next 16-bit word follows, emitted (n & 0x7f) times (consumes 3)
```

There is no end marker: the loop is bounded by the **source** span and never by the output count,
and it must then have emitted exactly `outWords` words. Both hold in 4/4 shipped saves. — SAV-PACK-007, SAV-EXT-009

```
observed operand domains, 4 saves       opcode values never seen, but all legal:
literal counts   1 .. 126               0x00  literal of 0 words, costs 1 byte
run counts       2 .. 127               0x7f  literal of 127 words
                                        0x80  run of 0 words, costs 3 bytes
                                        0x81  run of 1 word
```

Encoder side, if you write saves: emit a run whenever `word[i] == word[i+1]`, else a literal; write
`outWords` into the first dword; pad the uncompressed stream to an **even** byte length first. The
shipped encoder's output buffer is sized to its input, `inWords * 2`, with no headroom.
- SAV-CODEC-022

The repeated unit is a **word**, not a byte: 145–150 runs per save repeat a word whose two bytes
differ. Do not byte-swap it — that keeps the length and corrupts the content.

## The decoded body

| Off | Type | Field | Notes | Claim |
|-----|------|-------|-------|-------|
| +0x00 | u32 | `world+0x04` | zeroed at mission start; monotone with save order; `1` in the "Restart last mission" save | SAV-STREAM-010, SAV-HEAD-025 |
| +0x04 | u32 | `world+0x00` | `== (+0x00) >> 4` exactly, 4/4 | SAV-STREAM-010, SAV-HEAD-025 |
| +0x08 | u8 + ASCII | map name (`world+0x28`) | MFC `CString` short form; `06 "10.alm"`. On load it is re-opened from the install as `Scenario\` + name whenever the mission number is non-zero | SAV-MAP-005, SAV-HEAD-025 |
| +0x0f | u32 x11 | `world+0x11c` .. `+0x140` | Store order is `+0x11c +0x124 +0x128 +0x12c +0x130 +0x134 +0x138 +0x13c +0x148 +0x144 +0x140`; the last three are not ascending. Explicit constructor zeros and literal restore are established; all-path first-save values are not | SAV-HEAD-025, SAV-WHEADINIT-520, SAV-WHEADLOAD-521, SAV-WHEADLIMIT-525 |
| +0x3b | u32 | mission number (`world+0x80`) | `10`. Parsed with `atoi` from the map file name at mission start; also the `Scenario\` switch above | SAV-HEAD-025 |
| +0x3f | u32 | difficulty (`world+0x84`) | `2`. The load arm keeps the stored value only if it is 1..3 | SAV-HEAD-025, UNIT-GATE-012, SAV-657 |
| +0x43 | u32 | `playerList+0x20` | `6`. Meaning unknown | SAV-HEAD-025 |
| +0x47 | u32 | player count | `5` = this map's type-5 slots. Then that many `Player` objects, each with its groups and their actor lists | SAV-HEAD-025, SAV-ROSTER-024 |
| .. | — | object stream | to the stream's end: the object graph, the shape byte, two cell tables, a fixed region, the trailer | SAV-STREAM-013.. |

Offsets after the map name shift with its length; the table above is for `10.alm`.

### The document programme, end to end

The table above is the head of `FUN_004d0cb7`. The whole routine, store arm, is the complete
statement of what a save contains:

```
u32 +0x04 · u32 +0x00 · CString +0x28 · 13 x u32 (+0x11c..+0x140, +0x80, +0x84)
FUN_0051102f   u32 playerList+0x20, then FUN_00527d70: u32 count + count x `ar << CObject*`
CALL [vt+0]    on *(*(this+0x14)+0xc) at 004d0e5f — dead-actor list:
                 u32 count + count x `ar << CObject*`
u8             world half present (store 0/1; load zero/nonzero)
  FUN_005102db -> FUN_00527a10   u32 count + refs   Building
  FUN_005114b6 -> FUN_00527f90   u32 count + refs   SpellEffect
  FUN_00544a60                   block array + cell-record table + the terrain's identity key
  FUN_00539310                   4374 raw bytes
  FUN_0051149d -> FUN_00527e90   u32 count + refs   Sack
u32 trailer discriminator
  if == 0xbadface1: u32 [0x00609b0c]
FUN_0053e800                   raw 400 bytes from world+0x118
```

Both document arms join at `004d15dd` for the last call. `FUN_0053e800` passes exactly `0x190`
bytes to `CArchive::Read` or `Write`; the 400 observed zeros are loaded state, not free padding. If
the logical stream ends at an odd byte offset, the outer writer extends it by one byte before word
compression. That one byte is not consumed and is not value-constrained; no byte is added when the
logical endpoint is even. — SAV-TRAIL-026, SAV-DOC-053

`world+0x118` itself is a heap-allocated, exclusively-owned 400-byte block: world's own
constructor allocates it with `operator new(0x190)`, zero-fills it through a trivial constructor
(`REP STOSD`, no field-specific init) and stores the pointer at `this+0x118`; world's own
teardown frees and nulls it alongside four sibling owned-singleton globals, including the session
pointer — a separate allocation with a separate lifetime, not a session sub-object. Two of its 100
dwords are debug-trace verbosity toggles, not free padding or a journal: `+0x0` gates
"turn tracing," `+0x4` gates "script tracing," each with a paired on/off confirmation string and
four further consumer sites gating `Script:`-prefixed trace/error lines. Over 45 of the
corpus's 46 distinct preserved save streams — 99 files across 16 preserved-save directories
(byte-identical executables, above), minus the one stream this repository's own reader cannot
open — all 400 bytes are zero in every save, extending the finding above with the same result.
`[0x00609b0c]` is a plain global int with exactly three whole-image touch points — a
zero-initializer, the unconditional store-arm write shown above, and the marker-gated load-arm
re-read — zero in 44 of the 45 parsed saves and 520 in exactly one, an original-written stream and
not one of the six this project's own writer produced, returning to 0 in the very next save from
the same directory: not a monotonic counter, and the write that first produces a non-zero value is
not located.

That transfer is one call site, `004d15ed`, and both directions run it: the store arm tail-jumps
into the load arm's epilogue at `004d0efa`. `FUN_0053e800` has exactly one direct call site and no
stored pointer anywhere in the image, so it is not virtual, and neither is the block's constructor.
Direction is the only thing the archive's mode word selects: a set loading bit reaches `0057b908`,
which memcpies out of the archive's own buffer into the caller's, and a clear one reaches
`0057ba16`, which copies the other way; the count is `0x190` at every mode value tested. The block
is read unconditionally, even when the discriminator does not match — only the `[0x00609b0c]`
re-read is gated on it. Of the block's 100 dwords, only `+0x0` and `+0x4` have any located
consumer: the complete `[base+0x118]` displacement census over `.text` is 70 sites, and after
crossing it against the complete absolute-dword reference census for `[0x005cd758]` and
disassembling every survivor, the six that reach this block are exactly the trace-gate and toggle
sites already named, each touching dword 0 or dword 1 and no other. Dwords 2 through 99 have no
located consumer. That is a complete search of one instruction form, not a statement about the
image: a consumer handed the block pointer as an argument uses no `+0x118` displacement and would
not be found.
— SAV-646, SAV-647, SAV-648, SAV-649, SAV-790, SAV-791

This page's "world," in `world+0x118` and as the document routine's own `this` below, is one
specific global, `[0x005cd758]` — every disassembled `+0x118` dereference loads its base from that
address, including a same-routine reuse that ties the two phrases to the identical cached local.
`formats/session/format.md`'s own object table independently fixes the same address as the
"server singleton," `0x174` bytes; that page's own, separately-named "world"
(`[0x005f22c8]`, `0xa4558` bytes, the loaded map's terrain) is a different object, not this one. —
SAV-657

All 31 distinct current documents use world discriminator 0 or 1 and the marker-present trailer
arm. The exact load reader also implements noncanonical nonzero world values and the marker-absent
trailer arm from the same instructions; those capacities have synthetic tests but no
original-produced witness. Its CString primitive likewise implements 8/16/32-bit lengths and the
`ff | fffe` Unicode arm, while this population reaches only one-byte ANSI lengths. All 1,274
reached Unit containers take the present arm, so its absent arm is synthetic-only; embedded
Spellbook presence reaches both arms, 1,254 absent and 20 present. —
SAV-FULLREAD-252

`this` is the world (`[0x005cd758]`) itself. Its own `+0x14` field holds a pointer to a separate,
statically constructed object — not a second view of the world and not heap-allocated — whose four
dwords are the four list heads the programme reads, and whose `+0xc` is the dead-actor list
manager, not `world+0xc` directly (`SAV-657`). Its slot 0
is `FUN_00510401`, which calls `FUN_00527ac0` on the embedded `CObList` at manager `+4`. The record
is four bytes only when the list is empty. `SAV-DOC-053` identifies the same manager through its
construction, the death-path insertion, the decay walk and this serializer.

On load, `FUN_00527ac0` creates or resolves every actor through the archive's shared object index
and appends the resulting pointer. After world reconstruction, document load calls manager slot 1,
`FUN_0051041d`; it walks every dead actor and dispatches `actor->vt+0x24`. `Unit` implements that as
`FUN_00510dc0`; both humanoid tables use `FUN_00510e49`, a one-call wrapper around the same routine.
The actor rebind resolves position `+0x08` and actor references `+0x5c`, `+0x64`, `+0x44`, `+0x68`
and `+0x40` through the archive identity map. It does not join the actor to an authored `.alm`
record. — SAV-DEADLOAD-124, SAV-DEADLOAD-125

The dead-list manager does not derive actor state. `Unit::Serialize` restores timer `+0x6c` as an
`i8`, health `+0x94` as an `i16` and stage `+0x13c` as a `u8`, each through a separate stored
address. The post-load actor routine changes none of them. Its only stage gate is stage 0, which
repairs two embedded objects not used by dead actors. — SAV-DEADLOAD-126

The five top-level list serializers use one shape. Store: `IsStoring`, `GetHeadPosition`, `GetCount` written
as a plain `u32` (not `WriteCount`), then per node `GetNext` and `ar << CObject*`. Load: the count,
then per element `ar >> <class>` with the element's class descriptor pushed, then `AddTail`.

**There is no framing between top-level records and none is needed**: the count bounds the list and
each element is a full record written in place. `SAV-FULLREAD-252` joins the subsequently completed class
programmes, codec and physical tail into one call-site reader. It closes through the logical trailer
on **31 of 31 distinct current streams** (55/55 paths), then classifies transport alignment as
exactly zero bytes on 21 and one byte on ten before continuing through label, `&YA1`, campaign and
physical EOF. Every digest reports zero decoded/physical interval gaps or overlaps. A second
trailer-shaped dword inserted inside a bounded raw member does not move the endpoint, so marker scan
is not the grammar. — SAV-DOC-053, SAV-PLDIARY-054, SAV-UNITCORP-157, SAV-CLASSSER-175,
SAV-CLASSSER-177, SAV-FULLREAD-252

### The world-half load and post-load order

The load arm of the same `FUN_004d0cb7` has two stages. Its exact sequence is:

```
deserialize live actors/Players
deserialize dead actors
rebuild the global actor list
deserialize Buildings
deserialize SpellEffects
construct the named map and ALM terrain; publish the terrain global
deserialize saved block records, cell records and terrain identity
deserialize session
deserialize Sacks
rebuild triggers
post-load live actors
post-load dead actors
rebind all cell-record object keys
post-load session
post-load SpellEffects
```

The two orders are deliberate. Actor, Building and SpellEffect identity bindings exist before
terrain construction. The saved cell payload overlays the ALM-created cell hash before Sacks load;
the all-cell identity pass waits until after Sacks and triggers exist. SpellEffect object post-load
runs after that pass. — SAV-CELLLOAD-108, SAV-CELLLOAD-109

### Reconstruction provenance

A resumable world is not reconstructed from one source or by one generic fixup. The bounded route
has four non-interchangeable provenance families:

| provenance | examples | required order |
|---|---|---|
| SAV-direct | clocks, Player/actor/dead state, session latches, campaign collections | deserialize and retain until the named lifecycle consumes or replaces it |
| SAV identity | archive object graph, terrain key, cell occupancy/layer keys | bind file-local keys only when their owner population exists; key numbers are not cross-file identity |
| SAV overlay | terrain block records, cell payload and Fog | construct the selected external baseline first, then overwrite/insert/apply the saved delta |
| external-derived | ALM terrain/cell hash, compiled triggers, Data definition pointers, campaign MapPoint objects | select by saved map/index/relation or the class-specific constant branch, then join the resulting live objects to saved state |

City-shaped documents omit the world half and the located loader has separate fresh-construction
controls. That does **not** prove an independent route input: a live countermodel reads SAV
`world_half` and conditionally selects resume or fresh construction. The present corpus has no
crossing that fixes SAV bytes, external inputs and shape while varying route alone. The complete
per-field provenance is promoted in `SAV-RECON-268` and `SAV-RECON-269`. SAV-only and external-only
sufficiency are refuted on the located positive routes. Joint sufficiency of a hybrid remains
Unknown: constructor outputs and opaque direct fields are not closed, a computed pre-entry tick
and the first normal tick lack independent projections, and including those unresolved original
transitions makes a candidate circular. A full-world reconstruction story remains blocked; exact
framing and already closed subsystems are narrower decisions. — SAV-RECON-268, SAV-RECON-269,
SAV-RECON-270, SAV-SUFF-302, SAV-SUFF-303

### The group record and what it does not carry

`Group` is0x48 bytes. Its constructor and selected base constructors preserve incoming
Group+1c while explicitly constructing its embedded word list: vtable at+20,
zeros at+24/+28/+2c/+30/+34 and grow-by10 at+38. It supplies+3c and zeros+40/+44.
Nonzero byte patterns and valid source-pointer sentinels distinguish these writes
from inherited allocator zeros; the sentinels do not assign pointer semantics to+1c.
— SAV-GRPNEW-576

The selected allocator rounds the72-byte Group request to80. With the image's
initialized threshold480, a controlled existing-page small-block success changes
metadata without clearing the returned payload. A controlled threshold-zero route
passes flags0 to its imported heap call and does not clear the returned payload
afterward. Actual allocator mode, fresh-page/reuse history and OS-returned bytes
are not established. — SAV-GRPALLOC-577

The ordinary-command prologue's direct stores do not supply Group+1c or Group+20
contents. Its wrapper deletes at most the first rejected old Group before allocating;
it does not delete every rejected Group in that invocation. Guard and its selected
summary write GroupAI fields through Group+3c, not these Group fields. Actor/list
and behavioral callbacks remain separate mutation boundaries. — SAV-GRPCMD-578

- **The actor list has no order imposed on it.** `FUN_00527ac0` writes the `CObList` head to tail
  with no comparison anywhere in its store arm, and the load arm rebuilds it with `AddTail`. The
  group's own append `FUN_0050f950` detaches a unit from its previous group first, so a unit's
  position is destroyed and re-made whenever it changes group. Measured: the same five actors, by
  identity key, appear as one group of five, five groups of one, and three groups across one
  session. **Position is not identity.** — SAV-GRPORD-058
- `+0x1c` is an authored selector on the map path. Allocation retention is established only
  for the selected local routes and constructors above, not every first-SAVE producer.
  The older29-record population had10 zeros and15 distinct nonzero values; that corpus does
  not close runtime provenance. LOAD restores the dword literally, and normal resume indexes
  restored Groups by it before its initial-stance gate. It is not a remapped pointer key, nor
  is it proven harmless when dynamically produced. — SAV-GRPFLD-060, SAV-GRPIDENT-562,
  SAV-GRPNEW-576, SAV-GRPALLOC-577
- `+0x40` is0 in all29 measured records, but LOAD reads and immediately identity-remaps it;
  its gameplay meaning remains unclassified. `+0x44` has separate append-time and LOAD laws.
  Append copies the actor's current owner; LOAD subsequently overwrites that value from the
  saved Group reference. Missing keys become null. Neither field is established as the
  participant's own-character marker; that role belongs to `Player+0x34`. — SAV-GRPFLD-060,
  SAV-GRPLOAD-560, SAV-GRPOWNER-561, SAV-HERO-059

Normal LOAD uses the following local sequence. Calls are normal-return boundaries, not assertions
that transitive callees preserve other fields. Actor order is the existing head-to-tail rule.

| Boundary | Restored or constructed state |
|---|---|
| Player prefix | Register this Player's saved identity before reading its Groups |
| Group-list loader | Construct a fresh Group; serialize it; then append that Group to Player+24 |
| Group prefix | Restore embedded Group+20 word list and the owned AI block/list |
| Each Group member | Read Unit reference; detach old Group if present; append; stamp actor+70; copy actor+14 to Group+44 |
| Group tail | Restore literal +1c; read/remap +40; read/remap +44, replacing append's value |
| Player suffix | After the additional32-byte block and Diary, remap hero; append Group members to flat Player+20; stamp actor+14 with enclosing Player |

The Player suffix does not recopy Group+44 or rewrite actor+70. At that local boundary the
actor owner can be enclosing Player while Group owner is another resolved Player or null.
There is no proved fallback from a missing Group-owner key to containment.
— SAV-GRPLOAD-560, SAV-GRPOWNER-561

The AI block at `*(Group+3c)` restores80 raw bytes but not its old list pointer: `005391d0`
conditionally deletes the constructor list, reads the block, replaces AI+4c with a new word
list and restores its elements. Group+20's word-list meaning stays Unknown. Reaching group tick
with zero members sets AI order+20 toff before dispatch; nonempty reads the current order.
Arm0 re-finds each dispatched actor in the current Group before taking its successor. The rate
helper reads actor+70 -> Group+3c -> AI+44. A later conditional link arm writes AI+48, so
skipping resume's initial-stance branch does not prove whole-linker or first-tick preservation.
— SAV-GRPAI-563

The complete Group dispatcher chooses the restored/current order once and
passes AI word+0a to arms2/4/5. All byte values reach the shared withdraw tail,
which rereads current count/head. The ff overwrite occurs only at dispatch
entry; becoming empty later does not reset the byte during that invocation.
— SAV-GRPDISPATCH-568

Order0, order3, ff and the withdraw tail use a current-list identity search
for the just-dispatched actor before advancing. A missing actor ends that
stage; a newly appended successor can be visited. The withdraw stage starts
from the current head even after the selected order stage ended early.
These are mutation-sensitive caller laws, not proof that a loaded callback
actually mutates membership. — SAV-GRPMUTATE-569

The Group AI path and a member's patrol ring have separate pointer identities.
The path-copy setter reverses the AI list into actor order+90, chooses its head
as cursor+02, then prepends the actor's own cell. Actor-order LOAD instead
restores raw148 bytes and replaces pointer+90 with a new word list; it does
not substitute the Group path or constructor cursor for the saved values.
The ordinary walker searches the first equal cursor cell on arrival, advances
or wraps, and has an unguarded null edge when that value is missing. Its
re-anchor latch+04 is set whenever guard permits patrol continuation, not only
when a waypoint is reached. — SAV-GRPPATROL-570, SAV-PATROLCURSOR-571

The selected SAVE sequence consumes current Group+20 list, AI80/list, member
order and trailing+1c/+40/+44 values. Lists emit their current values, not node
identities. The AI and actor-order raw stores precede their respective list
stores and carry the newly allocated current pointers, which LOAD replaces
again. No atomic snapshot or original next-SAVE observation follows. The
direct Group dispatcher has no Group+20-element or AI+4c-element read;
transitive consumers and Group+20's meaning remain Unknown.
— SAV-GRPSAVENEXT-572

Authored-map miss copies loader+40 into Group+1c after construction. LOAD instead
restores the saved selector and ordered embedded-list values; it does not substitute
dynamic constructor defaults. The selected SAVE contains no direct normalization
of either queried field. Its list store precedes member serialization and its selector
read follows that serialization: a callback change can produce an earlier list count
and a later changed selector in one invocation. A preserving-cut prologue-to-store
trace is not a first genuine SAVE witness, atomic snapshot or safe authoring default.
— SAV-GRPFIRSTSAVE-579

### Object-stream framing

The stream is a Microsoft **`CArchive` object stream** and therefore names its own classes. A new
class is introduced by

```
FF FF | u16 schema | u16 nameLen | nameLen ASCII bytes
```

and the four shipped saves introduce **11**: `Player Human Weapon Item Armor Diary Effect Shield
Unit Building Sack`, all with `schema = 1`, in **first-use order** (the order differs between
saves; the set does not). — SAV-STREAM-010

Objects and classes share **one index counter, starting at 1, in stream order**: a class record
takes the next index, its first instance follows the record immediately and takes the next, and
every later instance of a known class is introduced by a u16 tag `0x8000 | classIndex` followed
by the instance's data. This walk reproduces all 393 tag words of the four saves, including the
values that differ between saves. The complete call-site reader widens this to 31 distinct
documents and 13,235 archive calls: 346 class+first-object, 3,950 prior-class new-object, 8,939
null and **zero genuine back-references**. A separate state replay agrees on every reached index
transition. Null is therefore observed; the plain-u16 back-reference arm is statically implemented
but still original-produced-unwitnessed. — SAV-STREAM-013, SAV-ARCHREL-253

### The placeable-object head — `Token::Serialize` `FUN_00510e5c`, 37 bytes

Every class deriving from `Token` — `Unit`, `Human`, `Item`, `Weapon`, `Armor`, `Shield`,
`Effect`, `Building`, `Sack` — begins with this, and nothing else does. `Player` begins with a
CString and its slot id; `Diary` begins with two embedded objects.

```
 file    width   source          what
  0..11  raw 12  *(this+0x10)    ONE BLOCK, copied out of a separate object the token points
                                 at (FUN_00544a30, PUSH 0xc): the position object below
 12..15  u32     this+0x04       runtime creation-order id (SAV-ID-015)
 16      u8      this+0x0c       registry key: Weapon/Armor/Shield resolve their class from it
 17..18  u16     this+0x0e       Human's load arm compares this against 0x21
 19..22  u32     this+0x08       its low u16 is SAV-ID-015's map unit id (head +0x13)
 23..24  u16     this+0x18       same byte as Unit's own +0x18 (Token is Unit::Serialize's
                                 first, un-adjusted component); recipient publication
                                 mask in located send paths — SAV-678; prior corpus SAV-655
 25..28  u32     this+0x1c       Item's own reader/writer is published (ITEM-VALUE-115,
                                 "Item+1c"); every other class's meaning is Unknown — SAV-636,
                                 SAV-654, SAV-656
 29..32  u32     this            IDENTITY KEY   — the object's own address
 33..36  u32     this+0x14       REFERENCE      — another object's identity key
```

Both `+0x18` and `+0x1c` serialize with the same helper pairs `+0x0e` (u16) and
`+0x4`/`+0x8`/`this`/`+0x14` (u32) already use, in both the store and load arm, and
`Token::Serialize` exposes no field boundary beyond width and offset — SAV-653. Over 3,170
position-verified corpus records, `+0x18` is 0 for every Armor, Effect, Item, Shield,
SpellTransport and Weapon record, 2 for every Building and Sack record, and 2 for 1,008 of 1,013
Human/Unit records (5 exceptions, unexplained) — SAV-655. `+0x1c` ranges 0..121319 across 90
distinct values with one sentinel, `0xffffffff`, used only by 34 of 70 Item records — consistent
with the shop-shelf price computation `ITEM-VALUE-115`/`SHOP-CONSUME-073` already publish for
`Item` specifically — SAV-656.

The raw position object at file offsets 0..11 is:

```
 object  width  meaning
 +0x00   u8     cell X
 +0x01   u8     cell Y
 +0x02   u16    packed cell = (cellY << 8) | cellX
 +0x04   u8     sub-cell X; constructors initialise it to 0x80
 +0x05   u8     sub-cell Y; constructors initialise it to 0x80
 +0x06   u16    not written by the five constructors or two copy paths bounded by SAV-TOKENPOS-074
 +0x08   u32    raw terrain pointer supplied by the caller
```

The accessors compute `fullX = cellX*256 + subX` and
`fullY = cellY*256 + subY`. The word at object `+0x00` is the same packed cell stored
explicitly at `+0x02`, which explains the corpus's two equal `u16` readings. The archive helper
writes or reads the whole object in one 12-byte operation. Its load arm therefore overwrites the
fresh constructor's current-terrain pointer before terrain exists.

World load later constructs terrain and reads one identity dword at the end of its record. It binds
that saved key to the newly constructed terrain in the archive map. Each dispatched `Token`
lifecycle hook then passes its Position to `FUN_0054d670`: zero returns unchanged, a non-zero key is
looked up, a hit replaces `+0x08`, and a miss remains unchanged. Therefore a live world Position
must carry the **same saved identity key as the terrain record**. The numeric key may be chosen by a
writer, but it is relational; neither zero nor a fixed dummy pointer is a general placeholder.

The Position family and resolved lifecycle paths never touch `+0x06..+0x07` except through the raw
archive operation, and copy/assignment omit them. Zero is the Medium-confidence authoring
placeholder. A byte-preserving reader may retain them. The remaining uncertainty is explicitly an
alias/bulk-copy and unwitnessed-runtime boundary, not a known field meaning. — SAV-OBJ-014,
SAV-TOKENPOS-074, SAV-TOKENPTR-075, SAV-TOKENLOAD-092, SAV-TOKENLOAD-095

Reach is class- and document-path-specific. World load dispatches live/dead actor lists and the
top-level SpellEffect list after terrain registration. Although the live-list insertion family includes
Human, Unit and Sack, Sacks deserialize after this pass and do not receive it here.
The complete loader does not call a Building or loaded-Sack lifecycle manager, does not establish a walk over
nested Items/Effects, and on the no-world/city arm constructs no terrain and invokes no lifecycle
manager. For those paths, a general `+0x08` placeholder remains Unknown pending the first later
dispatch, whole-Position replacement, dereference or destruction. — SAV-TOKENLOAD-093,
SAV-TOKENLOAD-094, SAV-CELLLOAD-112, SAV-DEADLOAD-125

The immediate static-object paths narrow that Unknown without closing it. After their base
serializer, Building, Outpost, Tavern, Shop and Sack make no direct call to the lifecycle hook,
Position rebind or either terrain-clamping Position helper. Their two distinct tick bodies are
empty. Building placement/removal directly read cell X/Y; Sack placement/removal directly read the
packed cell, with terrain supplied separately. These are bounded direct-body facts, not permission
to discard Position `+0x08`: nested calls, later events, resave behaviour and the first actual
post-load consumer remain Unknown. — SAV-POSLOAD-140

The later session-entry route supplies the first mandatory consumer without making it the absolute
first runtime event. Resume enters `FUN_00477c00(1)` and a city transition enters it with 0, which
admits fresh mission construction. Before session entry, nonzero campaign `+0x6b8` calls
`FUN_004d2551 -> FUN_004d891a`, an ordinary tick with computed actor/event dispatch. The later join
`FUN_004d303e` walks the Building and Sack managers. Both resolved sender arms call getters that
read Position `+0x00/+0x04` and `+0x01/+0x05`. Sack stores both full returned words. Building shifts
each word right eight and stores one byte, so only cell `+0x00/+0x01` influences its packet while
sub-cell `+0x04/+0x05` is read then discarded. Neither arm reads packed `+0x02/+0x03`, unmanaged
`+0x06/+0x07` or terrain `+0x08`. These are mandatory-prefix consumers unless the earlier tick
changed them; absolute-first order and terrain-key fate remain Unknown. A no-world-half SAV carries
zero loaded Buildings and Sacks; new-mission ALM objects are fresh. — SAV-POSTLOAD-220
(amended), SAV-POSTLOAD-221

The optional edge is not a special load-only repair tick. Its nonzero arm and the normal paced loop
both call the same wrapper, `FUN_004d2551`; when `server+0x2c` is nonzero that wrapper calls the same
ordinary sub-tick, which increments the sub-tick counter, drains queued commands, dispatches
`world+0x2c` entries and then the global ticking list. The zero campaign arm skips the wrapper
entirely. Therefore an executed P0-to-P1 edge is one ordinary sub-tick and the first normal paced
edge is the next wrapper call, while a skipped P0-to-P1 edge executes no tick instruction. —
SAV-FIRSTTICK-348

Building and Sack have complete no-op `vt+0x18` bodies, and the later Sack entry sender still omits
its container. Those direct facts do not close queued command effects or either computed virtual
receiver population. No live P0/P1/P2 value projection is established, so convergence,
persistence, replacement, consumption and destruction remain Unknown for Building, Sack and
Sack-nested Item/Effect state. — SAV-FIRSTTICK-349, SAV-FIRSTTICK-350

`Player+0x3d` controls retained-actor placement at that join. When the entry packet's `+0x0a` is
zero and the loaded latch is zero, `FUN_004d403c` places unseated actors: `FUN_005445d0` replaces
Position cell, packed cell, sub-cell and terrain pointer, and the placement search may clamp the
cell/sub-cell again. Already seated actors are skipped. The join then writes the latch to 1 on both
placement and no-placement arms; mission teardown clears it for the surviving human Player. —
SAV-POSTLOAD-220 (amended)

Actor entry mask `-1` admits equipment and carried-container walkers, but actor vtable eligibility,
recipient ownership and type-ID branches gate them before any Item is reached. Each reached Item
reference is consumed. Effect traversal then has two more guards: Item appearance `+0x40` must be
nonzero, and a plain Item with `+0x1c == -1` bypasses the common `Item+0x20` Effect walker. Automatic
Sack entry state does not traverse `Sack+0x40`. Pickup extracts quantity-one Items into the existing
actor container; a stack split can deep-copy Effects before destination insertion. The helper
deletes the source container, after which the caller nulls `Sack+0x40`, destroys the Sack and
requests another actor send under the same actor/recipient gates. — SAV-POSTLOAD-222,
SAV-POSTLOAD-223

The last two Token-head dwords define an identity key and one reference using that key map.
The load arm binds the key to the object it has just constructed, in a map at
`[0x005cd758] + 0x88` (`SetAt`, `00523f70`), and resolves the reference through that same map
(`FUN_00527d30`), substituting `0` when the key is absent. `Player`'s seventeenth field and
`Diary`'s `+0x2c` use the same mechanism at their named sites.

For these sites a writer must mint a unique nonzero identity key per object and write a
reference as its target's key. The value is arbitrary; the identity is not. Missing-key
nulling is this resolver's rule, not a universal property of the map. These keys are neither
Token creation-order/map-unit IDs nor MFC archive indices. — SAV-TOKEN-034, SAV-PTRMAP-035
(partially retracted), SAV-ARCHREL-253

Across 31 distinct complete documents the same call-site reader records 4,280 pointer-map
definitions and 9,382 references: 5,722 uniquely resolve, 2,963 are zero and 697 nonzero keys remain
unresolved, with no duplicate nonzero definition. All 697 unresolved rows are Position terrain
keys; they remain unresolved rather than being guessed from coordinates or record order. These are
body-level pointer-map relations, not CArchive back-references. — SAV-IDCENSUS-254

### Item movement is not a Token-head normalization rule

Whole extraction followed by nonmerge insertion retains the Item pointer.
The inspected container/list bodies make no Item Token/Position store; this
does not establish invariance through every actor callback or alias.
Quantity-one draining may split a larger stack. — ITEM-WHOLE-128

Merge retains the destination Item, adds counts and ORs its Token `+0x08`
with the incoming value, then deletes the incoming Item. Pickup first sets
incoming `+0x08 = 1`, so the retained merge word is `oldDestinationFlags OR 1`.
New Sack construction adopts the supplied container and separately constructs
Sack Position; an existing Sack drains and deletes the supplied container.
Death reaches both branches, not an unconditional container-identity transfer.
Neither branch is evidence for rewriting contained Item Position to Sack or
actor Position. — ITEM-MERGE-129, ITEM-GROUNDMOVE-130, ITEM-DEATH-012

The direct Armor/Shield/Weapon equip/unequip bodies write actor slots/stats;
their Effect and actor callbacks leave the complete transitive head write set
Unknown. Base Item equip can consume/delete the Item and mutate an Effect mode.
— ITEM-EQUIPMOVE-131

Weapon equip with a first kind41 Effect reconstructs its owned Spell from
the Effect's id byte; unequip deletes the owned Spell and clears the archive
object reference at Weapon `+0x80`. For a valid nonzero id, Spell `+0x08` comes
from the Effect, `+0x09` from Spells row parameter6 low byte, `+0x0a` from
parameter18 == 1, and `+0x0c` from parameter1 low word. Serialization writes
these fields and the new allocation's address. Same Item identity therefore
does not guarantee unchanged nested archive identity or values; address reuse
can hide the allocation change in numerical keys. — ITEM-SPELLMOVE-132

### The eleven class records

Read from each class's own `Serialize`, located from the constructor that stores the vtable
and taken from vtable slot 2. `list` = a `u32` count then that many `ar << CObject*`; `objref`
= one `ar << CObject*`; `raw N` = `CArchive::Write(ptr, N)`.

```
Token      the 37-byte head above
Effect     Token, u8 +0x3c, u8 +0x3d, u32 +0x40, u8 +0x0c                        = 44 B fixed
Item       Token, list(+0x20), u16 +0x40, u16 +0x42, u8 +0x44, u8 +0x45,
           u8 +0x46, u16 +0x48, u16 +0x4a, u8 +0x47                     = 53 B + the list
Shield     Item, raw 22 from +0x50                       (its own store arm is EMPTY)
Armor      Item, raw 22 from +0x52, u8 +0x50
Weapon     Item, raw 24 from +0x52, raw 22 from +0x6a, u8 +0x50, objref *(+0x80)
Building   Token, raw 22 from +0x52, u8 +0x40, u16 +0x42, u16 +0x44, u16 +0x46,
           u8 +0x48, u8 +0x60, u8 +0x61, u32 +0x64, u32 +0x68            = 77 B fixed
Sack       Token, u32 +0x3c, then on *(+0x40): list, u32 +0x1c, u32 +0x20
           (the container object at *(+0x40) is ITEM-CONT-004's own class, distinct from Item;
           its own +0x1c/+0x20 are not the numerically coincident Item+0x1c/Item+0x20 -- SAV-670)
Diary      CDWordArray(+0x04), CWordArray(+0x18), u32 +0x2c (a reference)
Spell      u8 +0x08, u8 +0x09, u8 +0x0a, u16 +0x0c, u32 this (identity KEY)     = 9 B fixed
Player     the 17 fields listed earlier, then u32 group count + that many group records
           (FUN_005102f4 -> FUN_00511938), then raw 32 from this (FUN_005392d0)
Unit       see below
Human      Humanoid::Serialize = Unit::Serialize + raw 24 XP[6] from +0x1cc
           + objref *(this + i*4 + 0x198) for i = 1..12   (index 0 SKIPPED)
           + objref *(this+0x1e4).  Human's OWN store arm writes nothing.
```

`Item`'s own `list(+0x20)` names "the Effect list" by an enforced base type, not by an exact one:
the underlying `objref` primitive checks every typed reference against the class the reader
requested — on a class name's first occurrence as well as on a repeat, and on an object
back-reference — and what it permits is any class derived from that request. This list therefore
accepts `Effect` and its one subclass `Effect_DirectDamage`, and raises the archive's own bad-class
error on any other name (`SAV-774`). Two writers were located for this list — the item-generation
`Effects=` grammar (`ITEM-EFFGRAM-070`) and, already published as the magic-shop effect generator's
own dispatch (`SHOP-EFFALT-071`, `SHOP-MAGIC-007`), a merge-or-append routine that reads the item's
own `+0x1c` price to size its interaction budget, not a use-time writer — and across a
98-file/5,134-item/234-member census spanning the full preserved corpus, every member either writer
produced is exactly `Effect`; the one derived class the reader would also accept was not observed
in this list, and no third writer was located (`SAV-775`). `Effect`'s own trailing `+0x0c` byte
(the fourth field in its row above) is 0 in all 234 census members, and both located writers
construct through the same default constructor that never sets it to anything else — from either
located writer an item's own `Effect` cannot carry the nonzero identity that
`MAGIC-ATTACH-016`/`MAGIC-DMG-005` document for a live, attached spell effect, because that stamp
is written by call sites neither item generation nor the shop generator's own merge routine
reaches (`SAV-776`).

The container's own tail, `u32 +0x1c` then `u32 +0x20` — `ITEM-CONT-004`'s shared container class
(a `CObList` with two dwords bolted on, embedded at `Sack+0x40` and, presence-gated, at
`Unit+0x7c`) — is `SAV-CITYSTORE-516`'s "insertion index" and "stored
load" (`ITEM-LOAD-005`'s corrected reading of `HERO-SIGHT-007`, maintained incrementally by
`ITEM-STACK-003`'s weight-times-count sites, not recomputed on read). `+0x1c` defaults at
construction to a sentinel (10000, `ITEM-CONT-004`) that exceeds any realistic item count, so the
generic insert primitive's `index<count?insert-before:AddTail` branch appends by default; 97.7% of
the preserved corpus (3667/3755 records) sits at that sentinel, and every one of the 88 exceptions
is on a `Human`-class owner, reached by `ITEM-CONT-004`'s own two named writers (a ground-pickup
order and an equip arm), both non-default and both exercised in the corpus. `+0x20` defaults to 0
and never exceeds 224 in the preserved corpus, three orders of magnitude under `ITEM-LOAD-005`'s
own `0xfa00` (64000) branch threshold. A transfer helper copies both fields verbatim between two
container instances; a reset pair zeroes both together — SAV-670, SAV-671, SAV-672.

`Unit::Serialize` `FUN_00510518`, before the store/load branch: `Token`; `list(+0x20)`;
`u16list(+0x15c)`; `u16list(+0x178)`; raw 24 from `+0xa6`; raw 22 from `+0xbe`; raw 24 from
`+0x114`; raw 64 from `+0xd4`; raw 180 from `*(+0x154)`; raw 148 from `*(+0x158)` followed by
`u16list(*(*(+0x158)+0x90))` — the last dispatch is inside `FUN_005390c0`, not in
`Unit::Serialize`. The two `u16list`s are the mover's static and dynamic route (element identity,
non-Serialize accessors and load consumption: `SAV-630`, `SAV-631`, `SAV-634`); the mask byte at
`*(+0x154)+5` inside the raw 180-byte mover block selects the plane `TERR-PASS-051` already names;
the `*(*(+0x158)+0x90)` list is the order's own copy of an AI patrol path (`SAV-GRPPATROL-570`,
`SAV-633`). Then its store arm:

```
u8 +0x49 +0x4a +0x4b +0x4c | raw4 +0x50 | raw4 +0x54 | raw4 +0x58 | u8 +0x60 +0x61 | u8 +0x6c
objref *(+0x74) | objref *(+0x78) | CString +0x80
u16 +0x84 +0x86 +0x88 +0x8a +0x8c +0x8e +0x90 +0x92 +0x94 +0x96 +0x98 +0x9a +0x9c +0x9e
u8 +0xa2 +0xa3 | u16 +0xa0 | u16 +0xa4 | u8 +0x12c | u32 +0x130 | u8 +0x134 +0x135 +0x136
u32 +0x138 | u8 +0x13c | u32 +0x148 (the byte at +0x14c is copied into it first) | u32 +0x144
objref *(+0x68)
u8 PRESENCE FLAG for +0x7c ; when 1, list(*(+0x7c)) + u32 +0x1c + u32 +0x20
u8 PRESENCE FLAG for +0x140; when 1, Spellbook::Serialize(*(+0x140))
u32 +0x5c | u32 +0x64 | u32 +0x44 | u32 +0x40 | u8 +0x48
```

These are **six** raw writes, not four. Their widths sum to 462, so `SAV-UNITLEN-045`'s fixed-part
arithmetic remains unchanged. The load arm repeats their order and widths. Each load-side presence
value is masked to one byte and tested against zero: every nonzero byte takes the arm, while the
store arm emits only 0 or 1. — SAV-UNITPROG-156

The `+0x7c` presence flag does not serialize the allocation's identity. The Unit constructor has
already created a fresh 0x24-byte container when the load arm runs; a nonzero flag restores the
count, references and two dword tails into that object. Death earlier transferred the actor's old
container object to a sack and installed a fresh empty container on the corpse. Thus the stage-3,
stage-4 and stage-5 record carries the corpse's replacement container contents, not the sack's
object identity. — SAV-DEADLOAD-130

### Actor definition and presentation binding

The restored row byte `+0c` and type word `+0e` are distinct selectors.
Let `U=[00609ba8]` and `H=[00609bbc]` be Units and Humans array bases.
After the selected complete serializer returns locally:

| exact class | definition pointer `+3c` |
|---|---|
| Unit | `U + 48*u8(+0c)` |
| Humanoid | 0 |
| Human | `H + 48*(u16(+0e)<33 ? u8(+0c) : 5)` |

The arithmetic wraps32-bit. Neither actor lookup checks collection size;
Item's size-checked lookup is different. Neither branch rewrites `+0c/+0e`.
The fixed5 arm covers every zero-extended word33..65535, not just a small
player-character set. — SAV-ACTORBIND-544

Unit copies the low byte of restored `+148` to `+14c`, then exact Unit
clears the byte. Exact Humanoid preserves it. Human keeps a nonzero byte
only when its selected pointer is nonnull and the definition-name search
for `NPC` does not return-1. This preserves the byte, not boolean1.
The backing dword is unchanged by these suffixes; Unit's later local SAVE
widens the effective byte back into `+148`. The conditional presentation
sender uses `+14c` to select a separate collection at `005f0758`.
The selected single-byte-character tests execute the original substring search rather than
injecting its return. NPC can occur inside the name; a match preserves2/255 and does not
enable an incoming zero. Other character modes and full runtime remain outside this check.
— SAV-ACTORDISPLAY-545

Archive creators call the class constructors before restore. Human's default
creator reaches `004f9065("Man_Unarmed",0,0)`; its table/equipment setup
and constructor-mode type rewrite are not rerun by the Human load suffix.
Fresh actor `+154/+158` pointer identities are construction-owned, while
their180/148 target bytes come from SAV. The order target's `+90` list is
reallocated by its load helper. These are not safe post-load defaults.
— SAV-ACTORCTOR-546

Footprint/domain `+49/+4a`, face/class bytes `+4b/+4c` and the mover
mask at target`+5` restore independently. Named getters/sender use the
retained selectors; the local post-read hooks repair references without
direct calls to the constructor mask builder or actor table lookups. This
slice establishes no domain/mask consistency normalization.
— SAV-ACTORINPUT-547

With actor stage`+13c == 0`, the hook calls order helper`0053b590` on the
fresh actor-owned `+158` object. It treats order dwords
`+0c,+10,+18,+20,+28,+30,+68` as keys in `[005cd758]+88`: zero is left
unchanged; a hit replaces the field with the mapped pointer; a miss preserves
the raw dword. Nonzero stage skips this repair. The helper supplies no class
or lifetime check, and this local boundary is not proof that every loaded
reference is safe before its first consumer. Order`+90` is different: LOAD
replaces it with a fresh list, not a mapped saved pointer.
— SAV-HUMRESUME-460, SAV-ACTORINPUT-547, SAV-ACTORCTOR-546

Human's conditional definition-name search is a local next use. Exact
Humanoid's null is not a harmlessness result: it shares Unit's `vt+58`
method, which passes `+3c+8` to a parameter lookup without a local null
guard. Actual exact-class acceptance, that method's loaded reach, later
callbacks, first frame/move/save order and original runtime remain Unknown.
— SAV-ACTORLIMIT-548

### Saved crossing and local cell failure

The retained mover/Position and route inputs are not normalized into a centered
actor by the selected post-read reference repair. With stage zero, mover+7c uses
the same hit-only identity replacement as the known hook; the mover's scalar
crossing bytes are not reset there. The known progress3 arm then continues its
selected step. These local statements do not establish the first actual loaded
consumer or global scheduling. — SAV-ACTORINPUT-547, SAV-HUMRESUME-460,
MOVE-STEP-040

Entry distinguishes a saved/existing cell payload from creation. Existing
domain1/2 entry can dispatch its trigger before refusing an occupied+04 slot.
Creation instead zeroes52 bytes, captures current cost/static baselines and takes
the post-creation store path. Reusing a record preserves its saved baseline/tail.
Domain3 has no corresponding trigger dispatch. — SAV-CELLENTRY-582

Footprint refusal stops iteration without undoing earlier cell writes or the
already written mover+72/+82..85 caches. Detach clears any nonzero domain slot,
not only a matching actor pointer; a missing node/zero slot refuses early. After
a successful clear, the record is removable only when its four occupant slots,
layer count and operation byte are zero. The remaining tail/residue does not
prevent deletion. — SAV-CELLFAIL-583, SAV-CELLLEAVE-584

The cache bytes have the following source order. Position is actor+10; all
offsets in the two cache columns are relative to actor+154's mover.

| Coordinate | Entry cache and accessor source | Successful detach cache and direct source |
|---|---|---|
| cell X | +82 = cached AL from `00544a10(Position)` | +86 = byte[Position+00] |
| cell Y | +83 = cached AL from `00544a20(Position)` | +87 = byte[Position+01] |
| fraction X | +84 = low8(`005449e0(Position)`) | +88 = byte[Position+04] |
| fraction Y | +85 = low8(`005449f0(Position)`) | +89 = byte[Position+05] |

The full-coordinate accessor laws make their low bytes the two fractions.
Entry caches cell-X/Y results before `0054a620`, writes its AX to word+72,
then stores+82/+83 and calls the full-coordinate accessors for+84/+85. The
accessors are explicit probe cuts backed by `SAV-TOKENPOS-074`; the `0x1357`
callback result is only a fixture, not an original value or calculation.
These four entry bytes are not proved to be one atomic snapshot. Detach instead
reloads actor+10 and directly loads each byte after recompute/optional removal;
its four stores are `0054538a/0054539c/005453ae/005453c0` in table order.
— SAV-CELLFAIL-583, SAV-CELLLEAVE-584

The selected step continues after a detach refusal and ignores the entry return.
Its centered cleanup can therefore empty actor+178's dynamic route and clear the
known mover progress/claim fields even when the destination slot refused entry.
The static route at actor+15c survives the named local instructions. Isolated
SAVE consumers then emit current Position12, static/dynamic list values, mover180
and current cell key/payload52; they do not reconcile this discrepancy. Unit's
prefix writes the two lists before the mover. This is conditional local state,
not an original resave observation or a safe authored default. — SAV-CROSSNEXT-585

### Sack cell transitions

The Sack cell key comes from its Position+02 word. Registration first rejects
Dynamic bit0. Existing nonzero payload+10 refuses, including the same pointer;
an empty existing slot is written without recompute. Creation captures current
Cost/Static in a new zero52-byte record before setting the present flag, then
stores the Sack and recomputes. Saved baseline bytes must not be replaced by
invented constructor values on reuse. — SAV-SACKENTRY-590

Removal returns0 for a missing record but otherwise clears+10 without a zero
or equality check. After recompute, the four occupant slots, count+02 and
operation+2c decide deletion. Other residue and individual layer pointers are
not independent retainers. — SAV-SACKREMOVE-591

Successful empty-node deletion restores Cost and Static baselines, preserving
current Static bit4. It does not restore Dynamic from those baselines or clear
Dynamic bit5; the preceding recomputation remains visible there. This direct
state transition is not whole-allocator or original-runtime proof.
— SAV-SACKPLANES-592

The merge/create accessor additionally gates on Static bit5. Its caller's
merge, append, transfer and destructor calls are distinct authorities. Two
selected removal-caller slices ignore the cell-removal return before their
remaining calls; local cell success or refusal therefore does not by itself
describe a complete original item transfer or next SAVE. — SAV-SACKCALLER-593

### Human live runs and authoring boundary

The raw-a6 tail is `+0xb4..+0xbd` (10 bytes), raw-be tail is
`+0xc0..+0xd3` (20), and raw-d4 is `+0xd4..+0x113` (64). These 94 bytes
are mixed live state and modifiers, not three scalar defaults. — SAV-HUMRUN-444

| live fields | modifier fields | role |
|---|---|---|
| `+0xb4/+0xb5` | `+0xf4/+0xf5` | physical base/spread bytes |
| `+0xb6` | `+0xf6` is not folded | active skill index, written directly by weapon paths |
| `+0xb7/+0xb8` | `+0xf7/+0xf8` | second damage base/spread bytes |
| `+0xb9/+0xba/+0xbb` | `+0xf9/+0xfa/+0xfb` | elemental base/spread and assigned kind |
| `+0xbc/+0xbd` | `+0xfc/+0xfd` | unnamed final bytes, not used by the named folds |
| `+0xbe/+0xc0` | `+0xfe/+0x100` | defence/absorption words |
| `+0xc2..+0xcc` | `+0x102..+0x10c` | six protection words; slot-zero gameplay meaning remains Unknown |
| `+0xce..+0xd3` | `+0x10e..+0x113` | six damage-kind resistance bytes |

The modifier prefix is four signed stat-cap bytes `+0xd4..+0xd7`, then
seven words: speed `+0xd8`, capacity `+0xda`, maximum health `+0xdc`, health
regeneration percentage `+0xde`, maximum mana `+0xe0`, mana regeneration
percentage `+0xe2`, and sight `+0xe4`. The attack modifier begins with to-hit
`+0xe6`, General shadow `+0xe8` and five class-skill modifiers `+0xea..+0xf2`.
The traced derive/award loops skip the General shadow.

The attack fold adds one word and six bytes, assigns elemental kind,
and skips skills, the active-index mirror and the final two bytes. The
defence fold adds eight words and six bytes. Additions wrap at their storage
widths before any later derive clamp. — SAV-HUMFOLD-446

Weapon kinds 11/12 assign the to-hit modifier from General; removal subtracts
the weapon's own to-hit word. This is not an inverse for arbitrary prior
modifier state. Armor/Shield explicitly add their blocks into both defensive
copies; Weapon invokes derive after direct modifier writes. Reconstructing
all modifiers as a sum of currently equipped objects is therefore not an
established authoring law. — SAV-HUMEQUIP-447

Current health/mana, regeneration remainders and earned XP are distinct from
these modifier words. Admitted health/mana arithmetic changes current
amounts and remainders while leaving the selected 94 bytes unchanged; a raised
active skill can change live damage through derive without changing its
modifier. These are bounded instruction results, not observed post-load game
events. — SAV-HUMMUT-448

Regeneration consumes current/max/period at `94/96/98` and `9a/9c/9e`, plus
modifiers `de/e2`, as signed16; the stored remainder bytes `a2/a3` reload
unsigned. Modifier plus100 is32-bit. Products and the accumulator wrap at32
bits, followed by signed division truncated toward zero. The exact admitted
formula and local gates are in the hero regeneration section.
— SAV-REGENWIDTH-528, HERO-REGEN-021 (amended)

Each reached arm stores the remainder's low byte, then the quotient's low
word, then the signed minimum of that stored word and the maximum. There is
no lower clamp or remainder reset at the upper bound. Thus the serialized
pair is not necessarily a normalized hundredths value: negative remainders
become unsigned bytes157..255. With unchanged returning callbacks, mana0,
max101, period1 and modifier-101 produces current-1/remainder255 and the
next admitted call produces current0/remainder54. — SAV-REGENSTORE-529

Health period0 skips its arithmetic. Reached mana period0 faults before any
mana store; preceding health stores lie before an intervening callback and
are not rolled back by this local body. Product/accumulator overflow wraps;
signed quotient overflow is excluded by the stable admitted word/rate domain.
Invalid-memory test cuts are not observed operating-system fault behavior.
Actual callback effects, first-load admission and saved post-fault state
remain Unknown. — SAV-REGENFAULT-530, SAV-REGENORDER-531

The Unit wire order is six two-byte members `94,96,98,9a,9c,9e`, then the two
one-byte members `a2,a3`; within this scalar run their offsets from member84
are16,18,20,22,24,26,28,29. The earlier raw64 block at `d4` contains modifier
words `de/e2` at block offsets10/14. These width-preserving transfers do not
rebuild or normalize the values. Wire order is not arithmetic store order,
and no safe invented period/remainder vector follows from it.
— SAV-REGENWIRE-532

The selected archive/load-hook bodies do not rebuild modifiers from
equipment references or directly invoke derive. A later derive consumes
the restored modifiers and active index. The exact first computed post-load
consumer is still Unknown, including city/world order and callbacks.
— SAV-HUMLOAD-445

The attack initializer clears 22 bytes, not its full 24-byte serialized extent;
the final two survive that helper. The modifier constructor separately clears
all 64 bytes. Safe authored values for the unnamed tails remain Unknown.
Capacity and second-component modifiers have fold consumers but no nonzero
producer in the selected paths; the General shadow and active-index mirror
are not interchangeable with their live counterparts. Later aliases, bulk
copies and non-enumerated producers remain open. — SAV-HUMGAPS-449

The expanded reference hook does not supply a derive-before-read guarantee.
`00510e49 -> 00510dc0 -> 00510fca` repairs separate Position, actor-reference,
mover and order keys under their field-specific rules above. With ordinary
disjoint receivers it touches none of the selected 94 bytes and invokes no
Human derive. World resume depends on
YA1 InBattle, frontend mode/authority and the server run gate. City resume
drains queued commands around client restoration. Earlier document calls,
queues and callbacks remain outside that hook's negative scope.
— SAV-HUMRESUME-460

There are concrete conditional consumers without a local derive call:

| admitted path | selected-field computation | boundary |
|---|---|---|
| actor-state sender, effective damage mask `200` | `(b4+b9+b7) mod 256`, `(b5+ba+b8) mod 256` | computed presentation, not combat; recipient/class mask applies |
| phase-12 Human full tick, positive health and deficit/period gates | signed `de/e2` regeneration | the same list iteration projects state before actor full tick |
| base Effect generic apply/remove | signed `d8` test/possible clear before eventual derive | attached Effect identity, mode and cadence select the path |
| special Effect identity17 on Human | `e4` and `a4` add/subtract signed magnitude times256 | modular word changes without generic derive |
| admitted strike | signed `c0` absorption, indexed resistance, `c6` secondary protection, `bb` elemental selector | physical/secondary/third-component gates differ |

The sender's virtual class test is not derive. The chargen routine that does
derive constructs a new temporary Human, so it cannot establish loaded-actor
ordering. Base Effect identity8 apply also computes from `c6` without generic
derive; other Effect subclasses and earlier command order remain open.
— SAV-HUMPROJECT-461, SAV-HUMTICK-462, SAV-HUMSTRIKE-463

Derive reads `word[a8+2*b6]` for every nonzero byte index before clearing its
secondary/elemental damage and defensive live fields. No upper bound occurs
in that indexing step. Synthetic indices10,13,32,39,42 select respectively
`bc/bd`, `c2/c3`, `e8/e9`, `f6/f7`, `fc/fd` and affect computed to-hit/damage.
The resolver independently indexes a resistance byte at `target+ce+b6`
without a six-slot local bound. These are malformed-index sensitivity
controls, not accepted-SAV states or evidence of ordinary producer reach.
— SAV-HUMINDEX-464

An admitted third damage component expects `bb=1..5`; zero reaches a
diagnostic arm rather than selecting protection slot zero. No ordinary
`c2` damage arm is established in this resolver. Nonzero capacity and
secondary-modifier producers, ordinary meanings of residual shadows/tails,
and the absolute first loaded computed read remain Unknown. Both an earlier
callback derive and a saved-live-field-first path still fit the static
frontier; neither zero safety nor an original-state author is established.
— SAV-HUMSTRIKE-463, SAV-HUMFIRST-465

New hero and hired Human creation share `004f8e78 -> 004f6d3f -> 004f2bb4`.
Their allocation request488 rounds to496, above the image's480-byte small-block
threshold, reaching imported `HeapAlloc` with flags0. OS-returned contents were
not measured. The four embedded block constructors preserve `bc/bd` while
clearing `fc/fd`; both pairs are included in the later raw archive writes.
That constructor-stage distinction is not a first-save default: attached
Effect dispatch, indexed writes and later events remain unclosed.
— SAV-HUMALLOC-504, SAV-HUMNEW-505, SAV-HUMNEWSAVE-507

For the156 shipped nonempty Human weapon-definition cells per root, equipment
produces only selectors1..5. Signed attackType values below10 are narrowed to
a byte; values at least10 and removal produce0. Positive types10/42 do not
produce the corresponding malformed tail aliases. Negative/custom inputs,
archive restore, starting-skill aliases and other lifecycle producers are not
bounded by this table join. No safe zero vector follows. — SAV-HUMSEL-506,
SAV-HUMNEWSAVE-507

The original-process instrument has reached only a startup path checkpoint:
on a verified disposable EN copy, it redirected the original INSTALLDIR
argument process-locally and observed the original CWD API return success.
It stopped before resource-directory registration and any witnessed ordinary
SAV load. No Human observation hooks were armed. This does not distinguish
derive-first, saved-live-field-first or earlier load-side replacement; actor
identity, effective later file lookup and chronological reads/writes remain
Unknown. The next requirement is a verified normal load route with those
observations, not another constructor or serializer inference.
— SAV-HUMRUNTIME-476

The proposed restricted-token extension failed before original execution:
in owned synthetic fixtures, its outside DELETE and NULL-DACL write-open
controls succeeded. No original copy, menu, SAV load or effective resource
path was observed by that test. Its cleanup flag also required correction
by actual DACL readback. A separately verified write/child boundary is still
needed; this one failed profile neither closes Human chronology nor proves
that safe original execution is impossible. — SAV-INITGUARD-488

### Character skill and experience state

Three parts of the programme carry one `Human`'s progress:

```
memory                         serialized as                         meaning
Unit+0xa8+2i, i = 0..5        bytes 2..13 of raw 24 from +0xa6     six u16 skill levels
Unit+0x130                    raw dword after Unit+0x12c            aggregate experience; simulation i32
Humanoid+0x1cc+4i, i = 0..5  first raw 24 after Unit::Serialize    six raw dwords; simulation i32
```

All three are read directly back into the object. The aggregate is redundant but stored:
`aggregate == sum(xp[i])` on 306 of 306 player-subtree `Human` records over 22 distinct save streams.
The file does not derive it from either levels or per-skill experience on load. The word at
`Unit+0xa6` precedes the skill array and is not skill slot 0; it differs from slot 0 on 306/306.

The slot association is by the same index. The initializer reads `skill[i]`, writes `xp[i]`, and
adds that value to the aggregate for `i = 0..5`; both experience-award arms use the same subscript
for the skill word and experience dword. Controlled saves move one experience slot and the
aggregate while leaving the level vector and other party members unchanged. — SAV-HEROXP-063,
SAV-HEROSKILL-064

Weapon-borne spell awards use these same fields. The save contains no item-cast proficiency or
award source tag: a caster's admitted spell-side event changes its school slot, and a fighter
rider's event changes the current weapon slot before the ordinary serializer copies the six-slot
state. — HERO-ITEMSKILL-096

The related runtime context is already in `Unit::Serialize`. `actor+0x68` is stored through the
archive object-reference primitive. Temporary `actor+0x64`, kill-credit `actor+0x40` and the signed
attribution byte `actor+0x48` are stored and restored as raw `u32`, raw `u32` and `u8`. Ordinary
item-cast completion clears `+0x64/+0x68`; a save during the live interval carries their bit
patterns, but only `+0x68` is reconstructed as an object reference. The raw `+0x64` Spell pointer
and credited actor at `+0x40` are not remapped. Their validity and consumption after any load,
including same-process save/load, remain untested. — MAGIC-ITEMKILL-117, ITEM-CASTSTATE-056

Progress belongs to the `Human` object whose body contains it. The enclosing group actor list holds
that object, `Token+0x14` resolves its owning `Player`, and `Player+0x34` selects the primary
character without changing the other members' record shape. Position is not part of the mapping:
five player-owned objects retain their own values while all five list positions move across a
no-world-half/world-half pair. The saved `Token` key is an old pointer used by the file's load map;
it must be unique inside one graph but is not a durable id across save lifetimes. — SAV-HEROID-065

**Customisation boundary.** The original programme fixes six `u16` level slots and six raw-dword
experience slots interpreted as `i32`. A seventh paired slot cannot fit the 24-byte experience
tail. Extending that tail
moves the equipment references that follow it; reusing other bytes destroys different serialized
state. Keep the six-slot programme for original saves and put a wider model behind a versioned
extension.

**A record's length is not a constant.** Counts, a `CString`, presence flags and archive object
references all move it. Known fixed complete bodies include `Spell` (9), `Token` (37),
`SpellEffect` (39), `VirtualCaster` (44), `Effect` (44), `Effect_DirectDamage` (68), `Building`
(77), `Tavern` (81) and `Shop` (81). A `Unit`'s fixed part sums to **603** bytes, so with its name
and three null references a `Unit` record is never shorter than **609 + L**; `Human` adds at least
24 + 26 on top. — SAV-MEMBER-036, SAV-UNITLEN-045, SAV-CLASSSER-173

### The eight newly closed schema-1 `.data` Token-subclass programmes

The schema-1 descriptors enumerated in `.data` contain `Token` plus nineteen descendants.
Descriptor → creator → constructor → literal vtable → slot 2 closes that enumerated population.
Runtime-built descriptors and descriptors outside `.data` remain outside the census and closure.
Eight programmes were new after the existing member rows and `SHOP-SAVE-015` were joined:

```
VirtualCaster       Token | u8 +0x3c | raw 6 at the buffer pointed to by +0x40
SpellEffect         Token | u8 +0x40 | u8 +0x41
PointEffect         SpellEffect | objref *(+0x48), load expects Effect | raw u32 +0x44
AreaEffect          SpellEffect | u8 +0x48 | u8 +0x49 | u8 +0x4a | u8 +0x4b |
                    u16 +0x4c | objref *(+0x44), load expects Effect
SpellTransport      SpellEffect | objref *(+0x44), load expects SpellEffect |
                    objref *(+0x48), load expects AreaEffect | u16 +0x4c
Effect_DirectDamage Effect | raw 24 from +0x48
Outpost             Building | u32 +0x84 | u32 +0x88 | u32 +0x80 | u32 +0x8c |
                    count n | raw 8*n from embedded +0x6c
Tavern              Building | u32 +0x9c
```

Store and load preserve each sequence. `PointEffect` then resolves the raw `+0x44` key through the
archive identity map: a hit installs the live pointer and a miss writes null. That repair consumes
no further bytes. `PointEffect`, `AreaEffect` and `SpellTransport` have minimum bodies of 45, 47
and 45 bytes, but an object reference can introduce a class record and nested body.

Outpost's embedded object has vtable `0059c7b8`, slot-2 serializer `005244b0`, and only the
`CObject` runtime descriptor. Its count is `u16`, or `u16 0xffff | u32` for a wide value; helper
`00524810` transfers `8n` raw bytes. The Outpost extent is `95 + 8n`, or `99 + 8n` with a wide
count. No element meaning or invented class name follows from that framing. — SAV-CLASSSER-172..176

One SHA-distinct preserved save carries a nested `SpellTransport` → `PointEffect` →
`Effect_DirectDamage` graph. The class-record offsets are 46,949, 47,008 and 47,064; the body
offsets are 46,969, 47,025 and 47,089. Programme replay consumes 196, 136 and 68 bytes including
nested bodies, with every nested start reproduced by the independent exact class-record scan. The
other five newly read classes have zero exact records and zero raw-name hits in 55 paths / 31
distinct digests. Those zeroes are corpus scope, not a format law. — SAV-CLASSSER-177

### The embedded objects a record serializes

Eight sites dispatch through `CALL dword ptr [EAX + 0x8]`, so the class is not in the caller's
listing. Each was reached from the instruction that builds the object — the enclosing
constructor, or the serializer's own load arm, which allocates a literal size — then from that
constructor's `MOV dword ptr [EAX],<vtable>`, then from **slot 0**, `GetRuntimeClass`.

```
site                            class                              Serialize
Unit+0x15c   embedded           (unnamed; GetRuntimeClass = CObject) 00523100   static route, SAV-630/631
Unit+0x178   embedded           (unnamed)                            00523100   dynamic route, SAV-630/631
*(*(Unit+0x158)+0x90)           (unnamed)                            00523100   order's own list, SAV-633
Group+0x4c                      (unnamed)                            00523100
Group+0x20   embedded           (unnamed)                            00523100
*(Unit+0x140)                   Spellbook   desc 005c30c0            00500bbe
Diary+0x04   embedded           CDWordArray desc 005c8ef0, schema 0  00570799
Diary+0x18   embedded           CWordArray  desc 005c8f28, schema 0  00570b53
```

**`u16list`** — the unnamed class, ctor `00523040`, 28 bytes, vtable `0059c430`. It has no
runtime descriptor of its own, so a consumer can only refer to it by that vtable. It writes a
**count** (`u16`, or `u16 0xffff` then `u32` when `n >= 0xffff`) and then `n` × `u16`. Every
element is the class's own node `+0x8` on both the save and load arm; for the two `Unit`-embedded
sites that node is `MOVE-ROUTE-004`'s own packed route cell `(y<<8)|x` (`SAV-630`). The class's own
runtime append (used identically outside `Serialize`) grows its node pool 10 at a time and a full
teardown releases that pool, both closed by `SAV-632`. Six functions read or write the two
`Unit`-embedded lists outside `Serialize` — two already-published route-extraction routines plus
four further seed/consume/clear drivers — `SAV-631`; the order list's own corpus shape is strictly
0-or-2 elements and matches `SAV-GRPPATROL-570`'s already-published setter, though its writer
population is not exhaustively confirmed — `SAV-633`. Loaded list state reaches the same runtime
append path as ordinary play with no separate recompute found in the module swept — `SAV-634`.

**`Spellbook`** — `u32 +0x18`, `u32 n` (the element count of the array at `+0x04`), then
**`n-1`** references, for `i = 1 .. n-1`: index 0 is skipped in both directions. The same
idiom appears in `Humanoid`'s twelve equipment references.

**`CDWordArray` / `CWordArray`** — the count encoding above, then `n × 4` and `n × 2` raw
bytes. So a `Diary` record is `count + 4n`, `count + 2m`, `u32`.

A group record (`FUN_00511938`) is, in file order: `u16list(+0x20)`; then on `*(+0x3c)`, raw 80
followed by `u16list(*(+0x4c))`; then a `list` of the group's actors; then `u32 +0x1c`, `u32
+0x40`, `u32 +0x44`. — SAV-EMBED-039, SAV-WLIST-040, SAV-SPELLBK-041, SAV-DIARY-042,
SAV-HUMAN-043, SAV-SPELL-044

### Why an `Effect` cannot be chained

`Effect` is 44 bytes and the number is right; the instances are simply not consecutive. An
`Effect` is an element of the **`+0x20` list** of the object it is attached to, so what follows
one is the rest of that object's own programme. `Building` chains at 77 because buildings *are*
written consecutively; nothing steps from one `Effect` to the next. — SAV-EFFCHAIN-046

### Walking the stream by programme

`tools/savdoc` walks the complete document by archive calls rather than candidate words. Across all
current preserved roots it closes **30 of 31 distinct streams** and **53 of 55 file instances**.
The one failing digest reaches a genuine `SpellTransport` after 49 successful `Unit`-derived
records; no `Unit` arm remains unframed. That stop is the tool's class-dispatch limit:
`SAV-CLASSSER-175` publishes the complete `SpellTransport` programme and nested graph, but this
document audit was not extended to consume it. All eighteen distinct streams in the
`SAV-UNITCORP-157` population close. —
SAV-UNITCORP-157, SAV-CLASSSER-175, SAV-CLASSSER-177

There is a rollback trap in using that partial walk as a document-shape oracle. `savdoc` records
early `Player` values while walking, then parses the later world tail into a tentative document.
When the genuine `SpellTransport` is unsupported, it restores that tentative document, including
the default `WorldHalf=false`, but the already recorded Player rows remain. Two path copies of the
same world document were consequently counted as no-world/latch-1. `SAV-FULLREAD-252`'s
exact reader handles
`SpellTransport` and is the shape authority: those paths are world-half, yielding the corrected
**52/3** path split above. — SAV-SUFF-300

`tools/savreplay` remains the older first-`Player` programme. Repeating its original fifteen-file
population gives exact agreement with the byte tag scan on **12 of 15**, not 14 of 15. On
`game0007.sav` the programme has 50 objects while the scan reports 44: four scan-only words occur
inside raw-180 payloads, advance an invented counter and misclassify later real records. Opposite
mutations discriminate the readings: a tag-shaped word inserted in raw Unit bytes changes only the
scan; neutralising a genuine class tag breaks only the exact programme. A byte scan is therefore
not independent framing evidence unless an archive boundary has already been proved. —
SAV-EFFCHAIN-046, SAV-TAGSCAN-158

`Building`'s 77 is the one length verified against files: stepping 77 bytes from the class
record chains **18/18** instances in six saves and **30/30** in five, stopping in every case on
the `0000` null-objref word, with the identity keys all distinct and the creation-order id
running `2..19` and `2..62`. Decoding the records field by field adds three internal checks:
`Building+0x40` reproduces the head's `u16` at +17 in 18/18 records; the dword at head **+19**
runs `1 2 3 19 20 26 28 30 31 32 33 43 …`, the map's own type-4 record ids; and the identity
keys ascend in `0x80`/`0x100` steps, an allocator's stride, while the reference field takes one
of a few values shared by groups of buildings. — SAV-BLDG-037, SAV-PTRMAP-035 (partially
retracted)

### Obfuscation and clamping

`XOR 0x5c073f4d` and `min(v, 0x7fff)` occur at four sites each, **all eight inside
`Player::Serialize`**. For every other class the bytes in the file are the bytes in memory.
The instrument is a listing sweep over the eleven bodies and everything they call, so it
cannot see a field already held transformed in memory. — SAV-OBFCEN-038

### The class set is 28 of the game's own, and 35 serializable in all

`.data` holds **28** MFC class descriptors with schema 1: `TableLine Token VirtualCaster Unit
Humanoid Diary Human Player SpellEffect PointEffect AreaEffect SpellTransport Spell Spellbook
Effect Effect_DirectDamage Building Outpost Tavern Shop CMultiShopShelf CMultiShopInstance
CMultiShopTemplate Item Armor Shield Weapon Sack`. The between-mission save `game0010.sav`
introduces a **`Spell`** class record, so eleven is not even the observed set. A reader built
on eleven names fails on a save carrying a spellbook, and by construction on a save from a
mission with a shop, a tavern or an outpost.

That 28 is the count of the **game's own** serializable classes. Relaxing the scan to any schema
and to `.rdata` finds **117** descriptor-shaped records: 82 at schema `0xffff` with a null
createObject (dynamic only, not creatable from a stream), the 28 at schema 1, and **7 at schema
0** — `CDib CStringArray CDWordArray CWordArray CByteArray CMapStringToString CMapStringToOb` —
which are serializable all the same and whose names live in `.rdata`. Two of them are `Diary`'s
members. That `Diary` route dispatches them directly and emits no class record; no schema-0 class
record occurs in the 55-path preserved corpus. This corpus fact is not a production ban.
Instruction-backed native generic routes can author exact records for zero-witness `Outpost`,
`Shop` and nested `AreaEffect`. Six zero-witness bodies have direct SAV reach without their own
descriptor, while six identities also have a static typed-loader-return to located generic-writer
chain conditional on that return. Ten classes have no SAV parent only within the explicit
106-function producer slice. Run-time-built, copied and aliased descriptors remain Unknown. —
SAV-CLASS-033, SAV-SERPOP-047, SAV-CONTSER-194, SAV-PRODPOP-205, SAV-PRODGEN-206,
SAV-PRODDIRECT-207, SAV-LOADRESAVE-208, SAV-PRODNOHIT-209

The unified reader joins all 35 creator-backed schema-0/1 descriptors to exact body programmes in
both byte-identical images. Identity-root coverage is narrower: 15 classes appear as archive
objects, `Spellbook` is embedded-only and 19 descriptor identities have neither witness. `Token`,
`Humanoid` and `SpellEffect` nevertheless execute as base programmes of derived objects, leaving
16 body programmes with no execution witness. Neither static closure nor base execution promotes
an absent exact identity to observed SAV production or acceptance. —
SAV-READPOP-255

### The fifteen non-`Token` descriptor bodies

The eight schema-1 classes outside `Token` are `TableLine`, `Diary`, `Player`, `Spell`,
`Spellbook` and the three `CMultiShop*` classes. The seven schema-0 classes are listed above. All
fifteen bodies now have exact framing. `Diary`, `Player`, `Spell`, `Spellbook`, `CDWordArray` and
`CWordArray` are described in the member and embedded-object sections above. The remaining bodies
are:

```text
TableLine          CString, Count(n), n * u32
CMultiShopShelf    empty
CMultiShopInstance empty
CMultiShopTemplate empty
CStringArray       Count(n), n * CString
CByteArray         Count(n), n raw bytes
CMapStringToString Count(n), n * (CString key, CString value)
CMapStringToOb     Count(n), n * (CString key, CObject* value)
CDib               BMP file header[14], info/palette[bfOffBits-14], pixels[imageBytes]
```

`Count(n)` is u16 when `n < 0xffff`; otherwise it is u16 `0xffff` followed by u32 `n`. It is an
element count. The dword/word arrays perform separate `4n`/`2n` transfers; the byte array passes
the same count as its byte length. Both maps walk current bucket chains, without sorting.
`CMapStringToOb` uses the archive object operation for each value, so nested values advance the
shared class/object counter and use its identity map.

`TableLine`'s descriptor creator selects vtable `0059be20`; its slot 2 writes the string at `+4`
and directly dispatches the array at `+8`. Nine other vtables return the same runtime-class
descriptor, but some have different serializers. The grammar above is therefore the load-created
base object, not a promise about an unregistered derived object that advertises the same class.

`CDib` flushes the archive and writes a BMP through its underlying file. Store authors `BM`,
`bfOffBits = 0x36 + 4*paletteCount` and `bfSize = bfOffBits + imageBytes`, then writes the 40-byte
info header, palette and pixels. Load requires `BM`, sizes the middle segment from `bfOffBits`, and
derives pixels from non-zero `biSizeImage` or the padded row formula. It never validates `bfSize`.

All three `CMultiShop` constructor-owned vtables use slot 2 `00401950`, exactly `RET 4`. The located
`Shop` writer still excludes stock. None of the nine newly checked classes has a class-record
witness in 55 paths / 31 distinct SAV digests. That is a corpus boundary, not a format ban. —
SAV-CONTSER-188..194

### The population

One `Player` per type-5 slot (5/5, names + slot ids); one `Human` **or** `Unit` per type-6 record
— the class chosen by the record's class key ({1,7,10,11,14,24} → Human, {69,73,74} → Unit on
this map) — plus the hero; one `Building` per type-4 record (18/18); sacks and their contents are
runtime objects. Inventory objects serialize inline after their owner, and the mechanism is the list:
`Item`, `Unit` and `Sack` each hold a counted collection whose elements are written with
`ar << CObject*`. Dead units keep their object with runtime id `0`.

These records are programmes, not fixed extents. `Item`, `Unit` and `Sack` carry counted or
presence-gated content; `Diary` carries two variable arrays; and a `Human` is a `Unit` plus 24
bytes and thirteen object references. `Building` is the bounded exception here: its 77-byte body
was derived independently from both writer and corpus. — SAV-MEMBER-036, SAV-DIARY-042,
SAV-HUMAN-043, SAV-UNITLEN-045, SAV-BLDG-037

### The block-plane record array

One structure inside the stream is decoded, and it is the one `rom.exe`'s `FUN_00544a60` was read
to emit (`TERR-PASS-053`):

```
u16 count
count x u32   (cell << 16) | (dyn << 8) | static
```

- `cell` is the 256-stride plane index `(row << 8) | col`, **strictly increasing**, and always
  inside the serializer's sweep window `[0x807, 0x807 + 0xe5e7)`.
- `static` is the byte at `world + 0x10000` for that cell; `dyn` is the byte at `world + 0x20000`.
  Bits 0..3 are the map-derived block bits (`TERR-PASS-049`), bit 4/5 runtime flags, bits 6/7 the
  occupancy bits — which appear in `dyn` only, never in `static`.
- Only cells with `dyn > 0x0f` are written. — SAV-BLOCK-011

**This array is a delta, not a plane.** It omits every cell whose only block bits are the map's own
`0x01`/`0x05` (2378 of 4071 nonzero in-window cells on `scn:10.alm`) and every border cell below
the window (647 of 2304). A consumer must derive the planes from the `.alm` first and then apply
these records over them; applying them to a zero plane yields a map with no terrain blocking.
— SAV-BLOCK-012

### The cell-record table (immediately after the block array)

```
CArchive::WriteCount(count)
count x { u16 packedCell; u8 payload[52]; }
```

The key set equals **exactly** the block-array cells whose **static** byte carries bit 5
(symmetric difference 0 in all four saves; counts 185/186/183/179) — block bit 5 means "this cell
has a record here". Stride 54 is discriminated against every neighbour. — SAV-CELLREC-017

Read from the routine rather than from the corpus: the table is the hash map embedded at
`terrain+0x540b4`, vtable `0x0059cd40`, `Serialize` `FUN_0054ff70`, which writes `WriteCount(count)`
and then, per node, **2 bytes at node+0x8 and 0x34 bytes at node+0xc**. The count is
`CArchive::WriteCount`, i.e. a `u16` below 0xffff and `u16 0xffff` + `u32` above it, not a plain
`u16`. There is no tail. Counts measured 185..197. — SAV-TERRKEY-056

The payload projection, at offsets after the packed-cell key, is:

| Payload | Width | Persisted role | Post-load/later state |
|---:|---:|---|---|
| `+0x00` | 1 | cost-plane baseline captured on record creation | remains live; plane recomputation reads it |
| `+0x01` | 1 | static-block baseline captured on creation | remains live; recomputation reads it |
| `+0x02` | 1 | occupied area-layer count | remains saved until layer attach/detach recounts six slots |
| `+0x03` | 1 | constructor zero; no bounded semantic consumer | persisted residue |
| `+0x04` | 4 | movement-domain-1/2 actor identity key | live pointer on successful archive lookup |
| `+0x08` | 4 | movement-domain-3 actor identity key | live pointer on successful archive lookup |
| `+0x0c` | 4 | Building identity key | live pointer on successful archive lookup |
| `+0x10` | 4 | Sack identity key | live pointer on successful archive lookup |
| `+0x14..+0x28` | 6 × 4 | SpellEffect keys indexed by map layer | each becomes a live pointer on successful lookup |
| `+0x2c` | 1 | trigger spell or operation | nonzero values other than 26 arm actor-entry casting; 26 arms relocation |
| `+0x2d` | 1 | trigger power | passed to either temporary-caster helper on actor entry |
| `+0x2e..+0x2f` | 2 | temporary-caster source x/y | passed as the first two coordinates to either caster helper |
| `+0x30..+0x31` | 2 | relocation x/y | read only under operation 26 |
| `+0x32..+0x33` | 2 | constructor zero; no bounded semantic consumer | persisted residue |

Terrain is constructed from the ALM before this table loads. For each saved record the reader
zeroes a 13-dword scratch, reads key and payload, looks up an existing node, allocates only on
absence, and copies all 13 dwords over the chosen node. It does not clear the constructed hash:
saved keys overlay, saved-only keys are inserted and construction-only keys survive.

After Sacks and triggers load, `FUN_0054d4f0` visits every node in that union. `FUN_0054d6a0`
passes the ten object dwords through the archive identity map. Success replaces the stored key with
a live pointer; zero and lookup failure preserve the saved dword. Every other payload byte is copied
back unchanged. Later actor attachment reads the unchanged trigger fields: `FUN_00544ec0` accepts a
nonzero `+0x2c` other than 26 and sends `+0x2c..+0x2f` to a temporary unit-target or cell-target
caster. The separate arrival reader `FUN_005495f0` accepts 26 and relocates through `+0x30/+0x31`.
Exact serialized joins over 14 saves resolve 2,373 of 2,373 populated actor, Building and Sack keys;
no saved layer key appears in that corpus. The no-consumer residue rows `+0x03` and `+0x32..+0x33`
are bounded to the 28 scratch owners plus the immediate-address owner and retain
alias/pointer-arithmetic blind spots. — SAV-CELLLOAD-109, SAV-CELLLOAD-110, SAV-CELLLOAD-111,
SAV-CELLLOAD-112, SAV-CELLLOAD-113

### The session block and the stream tail

After the cell-record table comes the **terrain object's own identity key** — a `u32` of `this`
written inline by `FUN_00544a60` at `00544bc4`, read back at `00544c56` and bound in the pointer map
at `00544c6c`; these are `SAV-CELLREC-032`'s four unattributed bytes (`SAV-TERRKEY-056`). Then the
session `Serialize`
`FUN_00539310`'s block — **4374** bytes, not the 4382 `SAV-TAIL-018` measured — then the `Sack`
objects, the marker and dword, then `world+0x118`'s 400-byte raw block and an optional one-byte
word-compression pad.

```
off      len   source            what it is
    0    400   session+0xbd34    100 signed int trigger result slots
  400   1000   session+0xbec4    1000 fire-once trigger latches
 1400     48   session+0x08
 1448    400   session+0xa828
 1848   2508   session+0xa9bc    the 50x50 diplomacy matrix starts at +8 of this
 4356      1   session+0xa48
 4357      1   session+0xa49
 4358      4   session+0xa4c
 4362      4   session+0xb3ac    the WIN counter
 4366      4   session+0xb3b0
 4370      4   session+0xb3b4    the LOSE counter
 4374                            end
```

The whole block is in the **world half** — the call sits after the shape byte, not before it — so
a between-mission save carries no trigger state at all. — SAV-SESS-031, SAV-CELLREC-032

Of the 100 slots at `session+0xbd34`, slot 93 (`session+0xbea8`) and slots 90–92 are the ones this
repository has found the engine writes by a literal, non-indexed displacement, independent of any
mission's own check/pattern/trigger content — a decode-anchored sweep of all 100 slots' own
displacements found a literal write at exactly these four and no others (`SAV-650`, `SAV-658`).
Slots 90–92 get this write only once, at session construction; slot 93 gets it at construction and
at two further runtime sites, each immediately followed by an unconditional call to the
check-evaluator then the pattern-evaluator (`formats/trigger/format.md`). Every other slot the
engine touches is reached only generically, at a script-supplied index, by the check-evaluator and
pattern-evaluator that execute a mission's own check/pattern/trigger definitions. On load, the
wholesale `Serialize` above restores every slot's raw file-stored value before the map trigger
programme is rebuilt. Its authored opcode `0x10002` check nodes overwrite their compiled slots
from node+0x48. The separate later call `0053b870` only repairs `session+0xa9c0` to the live session
address; adjacency to `0053b880` does not establish a register-writer call edge. —
SAV-650, SAV-651 (partially retracted), SAV-652 (partially retracted), SAV-658, SAV-710, SAV-711

## Uncompressed tail (`[bodyEnd, EOF)`)

The `SAV-EXT-009` region split is unchanged and now explained: this is the only part of
the file stored verbatim.

- **`[bodyEnd, bodyEnd + 0x100)` — label region.** Both save dialogs read this fixed 256-byte
  buffer, but the located selection-to-writer path passes only its NUL-terminated string onward.
  Bytes after the first NUL therefore do not retain source provenance and a fresh output may carry
  unrelated stale bytes there (`SAV-LABELTAIL-236`). The `E1 AC DF BA` bytes
  `SAV-PTR-003` locates 16 bytes *before* `bodyEnd` are **not** part of this region — they are the
  coded body's own last opcodes. — SAV-PTR-003
- **`[bodyEnd + 0x100, storeEnd)` — embedded `&YA1` state store.** One nested `&YA1`, the
  REG-style inline key-value variant (`SAV-EMB-004`, `REG-FMT-017`, `REG-REC-032`),
  with header `@8 = 9` = the 9
  top-level state keys `Character, CurrentState, Fog, GameOptions, Inventory, Objects,
  Projectiles, SpellBook, View`. Framing and value decode follow the REG format, so
  `storeEnd = bodyEnd + 0x100 + 0x18 + 32R + 4 + poolLen`. The leaf set is fixed: 28 records /
  19 leaves on a mid-mission save, 22 / 15 on a between-mission one, which lacks `Fog` and
  `Projectiles`. — SAV-EMB-004, SAV-TAILEXT-062
- **`[storeEnd, EOF)` — the variable campaign record at application `+0x548`.** The state-store
  writer is followed immediately by `FUN_00489270` on the same file; both load paths make the
  inverse call `FUN_00489580` after consuming the store. The complete programme below closes at EOF
  on all 26 preserved owner saves. Their observed record sizes are 268 (19 files), 308 (1), 310
  (4) and 368 (2). — SAV-TAILEXT-062, SAV-CAMPTAIL-070, SAV-CAMPPROG-071

### Bounded extension-survival check

Three apparent slack regions do not provide a provenance-preserving 16-byte carrier through the
located original load-to-fresh-save paths:

| candidate | bounded capacity / consumption | fresh-save provenance |
|---|---|---|
| bytes after the first NUL in the fixed 256-byte label | The two located save dialogs read all 256 bytes, but their selection path passes only a NUL-terminated string into application label storage. The 55-path corpus leaves 235..254 bytes after the first NUL; 38 paths contain nonzero residue there (469 bytes total). | The main writer emits the application buffer, not the source tail. Bytes after the new NUL can be unrelated old destination-buffer residue, so physical room and observed residue do not establish preservation. — SAV-LABELTAIL-236 |
| physical bytes after the campaign grammar endpoint | Both located main-load bodies return from the campaign reader without another file read or an EOF comparison. All 55 corpus paths end exactly at campaign EOF. | The located writer creates/truncates the destination, writes label, state store and campaign, then closes; it has no source-suffix copy route. A suffix is therefore outside the located read programmes and absent from their fresh output. — SAV-PHYSUFFIX-237 |
| decoded document transport gap | The exact gap after the logical trailer is zero bytes on 37 paths and one byte on 18; ten byte values occur. The reader does not compare its document cursor with decoded extent. | The writer independently adds one byte only when needed to make the logical stream even before word compression. A zero/one-byte parity gap cannot hold the 16-byte test frame, and the located routines do not copy a source pad into the output pad. — SAV-DECPAD-238 |

This is a bounded negative for those three regions and those located paths only. None supplies the
tested provenance-preserving custom carrier. It does **not** establish that the original accepts a
mutated label tail, suffix or decoded pad; it does not exclude another load/save route, a
semantically ignored field, internal raw padding, a runtime-built serializer, a multiplayer-only
shape or another carrier. — SAV-EXTSURV-239

### Campaign record after the state store

The base record is used once at the campaign head and again inside every child:

```text
6 × u32
u32 count(+0x24) · count × u16 from +0x20
u32 count(+0x38) · count × u16 from +0x34
```

The campaign writer then emits, and the reader consumes, this exact sequence:

```text
u32 childCount(+0x50) · childCount × (base + u32 +0x48) from +0x4c
u32 parallelCount(+0x64) · parallelCount × u16 from +0x60
                              · parallelCount × u16 from +0x74
u32 count(+0x8c) · count × u32 from +0x88
six independent count × u16 arrays:
  +0xa0/+0x9c, +0xb4/+0xb0, +0xc8/+0xc4,
  +0xdc/+0xd8, +0x104/+0x100, +0xf0/+0xec
u32 documentCount(+0x14c) · documentCount × (u32 value, u32 kind) from +0x148
u32 +0x118 · u32 +0x114 · u32 +0x110 · u32 +0x11c
u32 +0x120 · u32 +0x124 · u32 +0x128
u32 markerCount(+0x160) · markerCount × marker from +0x15c
```

With every count zero the record occupies 104 bytes. A marker is `u32 value`, `u32 n`, `n` bytes
of NUL-terminated string, then two `u32`; the writer sets `n = strlen + 1`, so its wire size is
`17 + strlen`. The seven scalars and the marker count are a 32-byte post-document suffix only when
`markerCount == 0`. Three corpora walked so far agree at zero markers with no exception:
`SAV-CAMPMARK-073`'s original 26-file corpus, `SAV-CAMPAIGN-086`'s separately-selected 55-file
corpus, and this page's own 55-file, 31-distinct-state corpus (`SAV-609`) — the last of which
recurs the first almost entirely (26 of its 31 distinct states are `SAV-CAMPMARK-073`'s own 26
saves, confirmed by SHA-256) rather than confirming it independently. This bounds the claim to the
population actually walked, not every ROM1 save: accessible saves outside these corpora remain
untested. The non-zero case is established by the paired programme, not by a save corpus. —
SAV-CAMPMARK-073 (partially retracted), SAV-609

Campaign `+0x118` is the selected mission. Field `+0x120` is 1 when the live world-map position
equals the first MapPoint. On load, 1 restores that first point; 0 looks up the MapPoint indexed by
`+0x118`. This record therefore stores a relation to campaign map data, not a raw coordinate pair.
The separate `&YA1` leaves `View/X` and `View/Y` remain a different record. — SAV-CAMPPOS-072

The object-side meanings are:

```text
base/child +0x04 mission             +0x08 MapObject
           +0x0c Payment             +0x10/+0x14 shop bounds
           +0x18 announce latch      +0x1c AddHero[]
           +0x30 EnableMercenary[]   child +0x48 age

+0x60 / +0x74  working / pristine per-type mercenary counts
+0x88           per-type hire flags
+0x9c           current mission Mercenaries[] shelf
+0xb0           permanent mercenary unlocks
+0xc4 / +0xd8   paired InnNPC[] / InnMission[]
+0xec / +0x100  ShopMission[] / TCMission[]
+0x148          documents (value, kind), kind 1 text and 0 picture
+0x110          AutoGetMission        +0x114 no name; record ctor/reset zero it, singleton chain reads it
+0x118          selected mission      +0x11c LastMission
+0x120          first-MapPoint flag   +0x124 mission time
+0x128          no name; record ctor/reset/FP-reader touch it; varies across saves, non-zero writer unlocated
+0x15c          selected-mission marker cache
```

`+0x114` and `+0x128` are raw passthrough in both the writer and the reader, unlike computed
`+0x120`, and neither is re-read inside the reader's own tail, unlike `+0x120` and `+0x118`.
`+0x114`'s app-relative form `campaign+0x65c` reaches the already-published server singleton
`[0x005cd758]+0x84`: both load paths increment the stored value into it, and its three readers
(named by `UNIT-GATE-012`, which grades its own reader enumeration Medium) include the `.alm`
spawner's three-way placement law. Every save walked for this measurement stores `1`, so the
increment always lands on that law's neutral arm. Both `+0x114` and `+0x128` are also zeroed by
the campaign record's own constructor and reset. `+0x128` additionally has a floating-point
reader, `FUN_00488c00`, that loads it once and never writes it; no writer of the non-zero values
`+0x128` takes in the walked corpus is located anywhere in the image, under either addressing
convention or a callee-reachability crosscheck from every function that computes the record's
own base pointer — the write mechanism for those values is an open Unknown, not merely an
unnamed field.
— SAV-598, SAV-599, SAV-600, SAV-601, SAV-602

The record carries no completed-mission collection. Main progress is the active main record and is
monotone: a lower requested main mission is rejected. A side mission exists while its child record
is retained and disappears on completion or age 2. Configured town candidates live in the three
building arrays; accepting one shortens that array (and the paired inn arrays), while announcing it
sets the mission record's separate latch. Completion drains `EnableMercenary[]` into permanent
unlocks. Town activation drains `AddHero[]`. Documents accumulate by pair. The marker cache is read
only by world-map presentation. — SAV-CAMPAIGN-076..086

`Mercenaries[]` (`+0x9c`) is a corpus-verified function of the record's own main mission, not a
cumulative union: across a 55-file, 31-distinct-state corpus spanning five distinct main missions,
its value matches the mission-number-keyed registry list for that save's current main mission with
zero exceptions, and a type present at mission 10/20 is absent at 30/40/50. Its own reader
(`FUN_00480560`) consumes the stored count and stored array directly; no instruction in its loop
reads a mission-number field, so membership is not recomputed live at tavern-render time either. A
separate, hire-bearing city-save pair outside that corpus leaves `Mercenaries[]` and every other
per-type array (working/pristine mercenary counts, permanent unlocks, `InnNPC[]`/`InnMission[]`/
`TCMission[]`/`ShopMission[]`) unchanged by the hire itself: in the one witnessed hire — one type,
one main mission — the only field the 268-byte campaign record changes is the hire-flag array's own
single dword element for that type. — SAV-606, SAV-614

The per-type hire-flag array (`+0x88`) is all-zero across the 55-file corpus, consistent with, but
not by itself a discriminating test of, its already-established unconditional mission-end zero. A
separate witnessed hire sets one element (the hired type) and the flag then survives mission entry
and stays set into mid-mission play, closing that corpus's own untested case directly: the array
reads zero throughout the 55-file corpus because none of its saves were taken between a hire and the
mission's actual end, not because a hire is never recorded. — SAV-607, SAV-615

A zero hire-flag array, including every file in the 55-file corpus, means no outstanding hire at
save time, not that a hire never happened: MERC-HIRE-003's writer census for that array is
exhaustive (hire, dismiss, the record reset's zero-fill, the mission-end zero-fill), so an all-zero
array is as consistent with a completed hire as with none. The corpus's own actor graph holds nine
fixed identities on a fixed per-file count and a narrow runID band, two of them unnamed classKeys 54
and 58 — Data.bin `Humans[]` rows 54 and 58, tavern types 10 and 14 (`NPC10_1`/`NPC14_1`, each level
1). A confirmed hire's own actors are classKey-58 `Human` actors inside a newly created group
under the Player, reproducing classKey 58's signature (typeWord, worn count, per-file count) field
for field — the differing absolute runID band is a session-local allocator artifact, not a
discriminator — so that signature is what a hire looks like, not evidence against one, for either
classKey: the same spawn loop produces it regardless of which type is hired.

classKey 58's own corpus population (14 files, 7 distinct states) is fully accounted for without a
hire: every state is main mission 20, where `Mercenaries[]` excludes type 14 and hire is therefore
structurally impossible (main-mission progress cannot run backward to reopen it), and the
already-published mission-20 scripted transfer (`PARTY-M20-030`; the trigger boundary and removal
are `formats/trigger/format.md`'s own subject) matches the population's count, typeWord and bundled
classKey-201 identity exactly, including a within-corpus negative control (1 of 15 mission-20 files)
consistent with that transfer's own trigger gate. classKey 54's own population (3 files, 2 distinct
states, all main mission 50) is a campaign map's own authored Player-slot placement, not a hire: a
hire is structurally excluded at every state this class occupies, because its own tavern type never
appears in the permanent-unlock array carried by any of the 55 corpus rows (including these three),
so the shelf filter that gates a hire cannot pass it regardless of `Mercenaries[]` itself already
naming it. Every classKey-54 file's own campaign record instead names side-mission map `41.alm` (the
selected-mission field, `+0x118`), which places exactly three of its own actors under the Player's
own slot, matching this class's count, typeWord and worn signature exactly, while the same files'
other four scenario-player actor counts — unrelated to this classKey — match that map's own roster
slot for slot. — SAV-608, MERC-HIRE-003, SAV-616, SAV-617, SAV-622, SAV-623, SAV-624, SAV-625,
SAV-626, SAV-627, SAV-628, SAV-629

Marker-picture rehydration is not LOAD-synchronous. The persisted marker record (value, string, two
`u32`) round-trips through SAVE/LOAD like any other campaign field, but the routine that rebuilds
the runtime picture pointer for a restored marker has exactly one caller in the whole image, the
world-map-enter routine, itself reached only through virtual dispatch — never from either SAV Load
routine or the campaign loader directly. A LOAD does not itself repaint a restored marker; the next
world-map entry does. — SAV-609

The fixed `game0020.sav` labelled `we have brian !!` is main 50, selected 41, with child records 41
and 51 and no 30 or 40 in any record or building array. Mission 40's Brian is transferred by a map
instant whose complete handover routine has no destination-roster comparison. Ordinary play blocks
a duplicate earlier, at the lower-main-mission load guard; the save's Brian object is in the Player
group graph, not this campaign record. A full 55-file actor census finds Brian in the Player group
graph seven times across the corpus (four distinct save states), never inside a campaign record,
corroborating this single-file finding at population scale. — SAV-CAMPAIGN-087, SAV-CAMPAIGN-088,
SAV-608

**G2 boundary.** Existing counts admit more entries of an existing shape. Adding a field, changing
an element width or reordering the programme is engine-class: the original reader has one fixed
sequence and no version-selected grammar arm. A compatible extension must retain the original
programme or introduce versioning outside it. `+0x114` has a named consumer and a statically
derived effect (SAV-599, via the already-published `UNIT-GATE-013`, not a runtime trace) and is
not free to repurpose. `+0x128` has three located consumers of its own (the
record's constructor, reset, and `FUN_00488c00`) but still no semantic name, and its non-zero
values still have no located producer (SAV-600, SAV-602): its offset-level layout is likewise
not permission to repurpose it, since an unidentified producer may still depend on it.

### The `Fog` section — the explored-terrain record

```
Fog/FirstState  int32     the tile-word state the first run carries; 0 on 14/14
Fog/Data        int32[]   run lengths, alternating state, sum == W*H
```

`Data` is a run-length encoding of **tile-plane bit 15** — the "has been seen at least once" bit
(`ANIM-TICK-011`, `TERR-TILE-079`) — over the whole `W × H` plane in its own linear order
`idx = col + row·W` (`TERR-EDGE-024`). The first run carries `FirstState`; the state flips at every
run boundary. Bit 14 is not persisted.

To restore: load the map's tile plane from the `.alm`, which carries bit 15 clear on every authored
cell (`TERR-FOG-087`), then walk the runs setting bit 15 on the cells whose run is in the set state.
That is what `FUN_00477c00` does at `00478870 OR word ptr [ESI],DI`. To write: emit run lengths over
`W·H` tile words, a run ending where `tile[i] & 0x8000` changes (`FUN_00478c40`). Between-mission
saves carry no `Fog` section at all. — SAV-FOG-061, TERR-FOG-145

Corpus: `sum(runs)` is 6400 on every save of `10.alm` and 20736 on every save of `20.alm`. Run
counts and set-cell counts ascend with play: 1 run / 0 cells on a restart save, 29/152 and 35/248 on
two saves minutes apart in one session, 273/3475 on the longest.

### Independently generated no-world writer boundary

`SAV-WRITER-284` independently generates two complete no-world SAV programmes without copying corpus
programme bytes. Both carry one `Player` / `Group` / `Human` graph. The larger also carries two
`Item` objects (including an existing-class tag), two plain-`u16` equipment back-references to those
items, and an embedded `Spellbook` / `Spell` graph. Every generated byte is covered by a declared
programme or file rule; an independent decoder and the pre-existing document walker agree on the
structure, and 18 targeted mutations distinguish the exercised framing, codec, archive, optional,
identity, nested-graph and tail rules. This establishes a bounded structurally valid writer for the
generated no-world subset, not compatibility with the original loader. — SAV-WRITER-284

Equality is checked only for the enumerated map/shape/label, selected Player/Human/graph and
tail-shape subset in `SAV-WRITER-284`. An altered Human skill or campaign scalar still closes both research
readers, so those values and every other unlisted emitted scalar are unverified. The semantic meanings
of 450 bytes inside each generated `Human`'s six raw blocks were not established by that writer
oracle. `SAV-HUMRUN-444` through `SAV-HUMGAPS-449` now distinguish selected live fields and
modifiers, but do not validate those generated values. The generated programme uses only the clear
world-half selector. Original acceptance, load-to-state fidelity, fresh-save
survival, present-world programmes and rare/boundary variants all remain Unknown.

### Independently generated world-present writer boundary

`SAV-WORLDWRITER-316` extends the independent generator with two complete `world-half = 1` programmes. The
1,510-byte minimal file exercises zero-count world collections, terrain overlays and Sacks; the
1,850-byte representative file carries one base Building, one Sack containing an Item containing an
Effect, one packed block record, one 54-byte cell overlay with separately resolved actor/Building/Sack
keys, and nine file-local identity keys. Its exact archive-operator boundary is 7/2/2: seven
new-class/first-object records for Player, Human, Item, Spell, Building, Sack and Effect; two
prior-class/new-object Item records; and two plain backreferences from the first two equipment slots
to the first two Human-container Items. Both files include the terrain identity key, exact 4,374-byte
session extent, a 25-record world state store with one 6,400-cell Fog run, and the 104-byte campaign
minimum with selected mission 10. Every emitted byte has a rule, two generations are identical, the
separate decoder consumes the full file, and the pre-existing `savdoc` walker independently closes
both documents and all seven emitted class families. Twenty-two exact mutations discriminate the
bounded grammar and enumerated checked subset. — SAV-WORLDWRITER-316

This is structural checked-subset A at Medium, not an original-compatible state vector. The equality
oracle exactly compares only the enumerated named edges, collections and checked values: the first
two equipment-to-container aliases, Player/Human/Building/Sack ownership and ancestry, terrain/cell
rows and keys, extents, Fog run and seven campaign wire values. A third equipment alias and a
reminted Spell identity both pass that oracle, so the remaining eleven equipment/reference slots and
Spell identity are unchecked.

The complete unchecked boundary is: eleven world-head dwords; Player/Group/Diary offset-only or raw
fields; Human runtime/class/type, skills, stats, XP, six raw Unit blocks and Position
`+0x06..+0x07`; the remaining eleven equipment/reference slots and Spell identity; Building's fixed
40-byte body; Sack scalars/tails; Item/Effect scalar bodies; all values in the 4,374-byte session;
the final raw 400 bytes; non-name application state; campaign meanings beyond the selected-mission
join; semantic safety of every declared zero/raw default; and every external ALM/Data/campaign
projection. An unchecked session-byte mutation still closes both research readers. Nonempty dead
actors and SpellEffects, a nonzero cell SpellEffect layer, multiple Players, AI-owned graphs and
rare/boundary populations are not generated. Original loader acceptance, observable loaded state and
fresh-save survival remain Unknown.

### Original-loader boundary of the generated fixtures

The selected representative no-world and world-present generated files both pass the save chooser,
but neither reaches a successful EN load transition. Each fresh process fails with access violation
`0xc0000005` at image offset `0x0007822d`. Thus chooser enumeration is not document acceptance, and
these two exact structurally closed programmes are not original-compatible. RU, post-load state and
distinct original resave remain Unknown. — SAV-ORIGLOAD-332

The first incompatibility is exact. Both files encode application-state record
`SpellBook/Shortcuts` as kind 6 with `Size = 0`. Original state load sizes the returned dword array as
`Size >> 2`, then unconditionally copies four dwords into campaign state; it faults on the first
copy because the zero-element record supplies no data pointer. A compatible candidate must therefore
supply a 16-byte `Shortcuts` record, but this is only a necessary boundary:
whether a corrected file passes the next consumer and every later incompatibility remain Unknown.
— SAV-ORIGFAULT-335

The four meanings are signed zero-based spell indices for F5, F6, F7 and F8,
in that order; `-1` is unbound. Writer `00478c40` reads the shared spellbook
controller at `campaign+0xec`, offsets `+0x64/+0x68/+0x6c/+0x70`, and emits
them as the kind-6 16-byte array. Restore paths `00477650` and `00477c00`
copy four values back without per-index validation. `SpellBook/Pressed`
separately carries current spell `+0x60`, while `IsOpen` carries book
visibility. All 27 preserved-root saves (EN 23, RU 4) have four `-1`
bindings and `Pressed=-1`; none witnesses populated-slot reload behaviour.
The static transport meaning is High; populated original reload and later
malformed-index behaviour remain unwitnessed. — AI-QUICKSAVE-281

The exact invalid-outer-magic neighbour does not appear in the chooser, while both valid-framing
generated positives do. This establishes only the bounded chooser discriminator; it does not define
the full chooser grammar or exercise the malformed file's deeper loader. — SAV-ORIGCHOOSER-333

### Claim-only writer completeness audit

A one-pass city/world matrix separates mandatory emission from optional population coverage. It
found three claim-determined writer corrections: `SpellBook/Shortcuts` must transport 16 bytes;
Spellbook authoring must accept `n-1` distinct records rather than hard-code one; and every emitted
cell-record key must equal the key set of block rows whose static byte has bit 5 set. Deterministic
city and world candidates now exercise a count-29 Spellbook, seven equipment aliases, the 16-byte
Shortcuts extent and the corrected block/cell relation.

The audited world candidate still has 25 state records and omits
`Projectiles/FreeIndex` and `Projectiles/IDs`. The missing wire details are now closed below, but
the semantic safety of opaque head, object, session, trailer, application-state and campaign
values remains explicit Unknown. The corrected candidates establish deterministic claim-backed
transport only, not original acceptance or state fidelity. — SAV-WRITERAUDIT-380

### World `Projectiles` application state

The world-present producer always emits a `Projectiles` root. `FreeIndex` is a kind-2 integer
copied from the manager's u16 allocator field. `IDs` is a kind-6 integer array: each live node's
u16 id is widened into one four-byte pool element. For every id the producer also emits a decimal
`Prj<id>` root with sixteen kind-2 leaves, in this order:

```text
x y z picture dir phase lastaction action actiondir actiontarget
actionx actiony actionz actionphase actionsegments actionspell
```

An empty manager still writes both leaves, so the canonical empty form has 28 records under nine
state roots. — SAV-PROJSTORE-428

The loader accepts a weaker form. Missing `FreeIndex` uses caller default 0; missing `IDs` leaves
the constructed word vector empty and skips the projectile loop. Kind 6 is read as `size/4`
elements, taking one low u16 from each four-byte element; kind 2 is accepted by a one-element
compatibility arm. Other kinds reach the typed getter error. Each listed id causes allocation,
sixteen defaulted leaf reads, world/terrain/id binding, manager insertion and a post-load helper.
Consequently the subtree is producer-complete but not loader-mandatory. — SAV-PROJLOAD-429

Across 55 lawful paths / 31 distinct documents, all 29 world-state documents use kind 2 and kind
6. Twenty-eight have empty IDs. The one nonempty document carries id 266, one complete `Prj266`,
and `FreeIndex=267`; this single case is insufficient to promote `max(id)+1` as a general relation.
— SAV-PROJCORP-430

A deterministic independent structural document closes the canonical empty tree at 28 records,
while an otherwise matched 25-record comparator omits it. This proves source-independent level-A
transport and refutes omission as a sufficient explanation for early rejection; it is not an
original load, action or resave witness. — SAV-WORLDSTATE-431

### Eleven world-head values

The original world constructor explicitly writes zero to all eleven words.
Its separate `memset` covers `[+a4,+118)`, and the initializer's indexed stores
end at `+114`; neither range supplies these zero values. Mode setup writes
`+c`, `+8` and `+150`, outside this block. These are construction facts, not a
guarantee that every first SAVE contains zeros. — SAV-WHEADINIT-520

The archive restores eleven raw u32 values in wire order before the world-half
flag. It neither normalizes them to Boolean values nor applies the later
difficulty clamp to them. Its leading reset concerns a separate global table.
Post-load dispatch and a possible pre-entry tick remain later consumers or
mutation opportunities, not an implicit reset to the constructor vector. —
SAV-WHEADLOAD-521

| world field | identified conditional consumer |
|---|---|
| `+11c` | Nonzero filters map-unit creation to Players with `+28 == 0`; zero also permits the later `Mission/Players` loader with `Monsters` fallback. |
| `+124` | Zero permits the `Outposts` registry reader during map construction. |
| `+128` | Chooses `00530b60(actor)` when zero or `0052fdf0(actor,Player+34,0)` when nonzero during actor reconstruction; the latter is Defend, not Follow. |
| `+12c` / `+134` | Either nonzero selects a helper on a newly constructed hero; `+134` also selects it independently on eligible reconstructed actors. The helper writes six words 100 at actor `+102..+10c`, six bytes 100 at `+10e..+113`, then invokes actor `vt+50`. |
| `+130` | Nonzero writes new-actor word `+8c = 40h` after its earlier derive. |
| `+138` | Nonzero constructs `PlasmaSword` and dispatches new-actor `vt+3c` with it. |
| `+148` / `+13c` | `+148` enables the modulo-5 full-clock watchdog; `+13c` selects its actor `+c0` clear or missing-actor countdown/replacement effects. |
| `+144` / `+140` | Constructor and wire behavior are established; non-serialization world consumers remain unidentified. |

The map branches are receiver-bound observations. The hero flag tests do not
run on the existing-hero return arm. The watchdog's nonzero countdown is
decremented; an invocation finding it already zero reloads 2 and attempts
replacement. These are instruction rules, not witnessed runtime effects. —
SAV-WHEADMAP-522, SAV-WHEADHERO-523, SAV-WHEADWATCH-524

A routine with matching `+148/+13c` clear operations has no found caller in
the reference-manager and aligned `.rdata` query. Its exact receiver and live
reachability remain Unknown; it is not a proved post-load reset. The finite
listing/raw-dword census also leaves undecoded locations, derived-pointer and
data-driven aliases, computed calls and unexpanded callees. No nonzero producer
other than archive restore was established. This does not prove unused fields
or an all-path first-save zero vector. — SAV-WHEADWATCH-524, SAV-WHEADLIMIT-525

Projectiles closure does not make arbitrary world values safe. Human values already have a later
observed shop hang/crash boundary, and unresolved world-head, Position/cell, session, campaign,
trailer and external-reconstruction values prevent a live-state no-copy/no-guess writer.
— SAV-WORLDFRONT-432

A source-faithful diagnostic demonstrates a narrower transport result. It rebuilds one
original-produced mission-30 vector as named object, terrain, session, state and campaign values.
The decoded 64,074-byte document remains equal, while an all-literal outer codec and a semantically
equal but physically canonicalized YA1 store ensure that source compressed packets and records are
not replayed. The complete reader closes the deterministic 65,826-byte output with no gaps,
overlaps or unresolved identities. This validates independent encoding of that one ROM1-derived
vector only; it supplies no producers for current implementation state and carries no original
chooser/load/action/resave result. — SAV-WORLDDIAG-433

### Original-readable generated no-world form

Original EN accepts an independently serialized no-world document whose writer
rebuilds the named object graph, CArchive registry, identity relations, state
store, campaign record and word codec. A complete 25-object semantic candidate
loaded with both heroes visible, supported the same shop route and produced a
distinct complete original save. This establishes the writer route for the
tested source-faithful values; it does not make arbitrary opaque values safe. —
SAV-ORIGWRITER-396

The smallest successful coherent reduction tested contains one Player, one
Group and one Human: exactly two CArchive objects. The Unit container is present
and empty; Item, Effect, Spellbook, equipment and the world half are absent. The
1,864-byte generated file and the original's distinct 1,860-byte output both
decode to complete 2,104-byte documents with 22 state records and no gaps or
overlaps. This is a topology minimum under the tested no-world grammar, not a
byte minimum over malformed inputs or alternative compression. —
SAV-ORIGMIN-397

The original output retains the two-object Player/Group/Human graph and exact
campaign bytes while reminting both identity keys. It also closes
`Inventory/IsOpen`; fourteen other checked state leaves remain semantically
equal. File identity numbers are therefore transient keys, while the resolved
relations are the durable part of the accepted graph. — SAV-ORIGCANON-398

Structural closure alone is insufficient. A different one-Human projected
candidate and the original's resave of it both load and close completely. The
generated candidate hangs on the shop route; after the original-produced
resave is loaded in a fresh process, that route crashes. Source-faithful
semantic candidates with 25, three and two objects pass that route. The exact
unsafe invented Human field remains Unknown. — SAV-ORIGVALUE-399

The one-Human city state can enter mission `30.alm` and produce complete
world-present saves. That transition reports a nonfatal failure to resolve hero
ordinal 10002 while running the second-hero `Remove Cure2` action. The mission
continues and saves, but original format acceptance does not guarantee that a
reduced party satisfies every authored mission reference. — SAV-ORIGMISSION-400

### Live-value authoring boundary

The accepted and failed two-object city programmes differ in more than their
Human/Unit bodies. Of 196 accounted semantic, relation, identity and physical
rows, 65 differ: 50 Human, 11 Player, one Group and three state rows. The
Player half includes 119-dword and 119-word Diary arrays plus a reference to
the enclosing Player, where the failed programme has empty arrays and a null
reference. File-local identities and state-pool offsets are not authoring
values. This census does not identify which logical row the original shop
consumer requires. — SAV-LIVEFIELD-412

Humans row 26 is named `PC_Danath`, and its 26 measured installed parameters
agree between EN and RU. Those parameters are initial definitions rather than
saved current state. Across all 65 differences, 26 have a named installed or
runtime producer chain, four are serializer/identity-derived and 35 still have
no producer within the bounded inventory. The latter include accepted Body,
Reaction, Mind, Spirit and HealthMax values that differ from their EN/RU
installed initializers without a named route producing the current values. A
writer must not obtain those 35
accepted values by copying or treating the source SAV as a default. —
SAV-LIVEPROD-413

Reciprocal diagnostic documents can transplant the Player/Group/Diary and
Human/Unit halves while retaining fresh base identities and rebuilding their
relations. Four such no-world documents are deterministic and close under two
independent complete readers at exactly the intended 11-row or 52-row
partition. That is structural evidence only: the transplanted SAV-derived
values are not safe authoring constants. One candidate has the separate
original-runtime boundary below; the other three remain unmeasured. —
SAV-LIVEREADER-414

For the exact 1,657-byte projected-base document carrying the complete accepted
Player/Group/Diary half, the EN original shows the chooser row and reaches the
city/shop route, then its process crashes while opening the shop. The requested
30-second city-stability and visible-Human subchecks were not separately
reported, and no resave was attempted. This one exact trial did not reproduce
the accepted shop result. It does not decide PGD-half sufficiency because
hidden runtime state, candidate identity, route variance and repeatability
remain open. It says neither that Player/Group/Diary is unnecessary nor that
Human/Unit is sufficient; joint dependence, individual causal fields, RU
behaviour and no-SAV producers remain Unknown. —
SAV-LIVEPGD-415

Save discovery uses the process directory after startup reads 32-bit HKLM
`SOFTWARE\\1C\\Allods` value `INSTALLDIR` and changes CWD to it. Running a copied executable with a
different caller-supplied working directory is therefore insufficient isolation while that value
exists; the registry-derived directory must itself be isolated and restored. — SAV-SAVEDIR-334

## Invariants (hold across the 4-save corpus)

- `bytes[0:4] == "Asg&"` and `u32@0x08 == 0x0BAD0002`.
- `20 <= u32@0x04 <= size` and `u32@0x0C == u32@0x04 − 16`.
- Decoding `[0x14, u32@0x04)` consumes the span exactly and emits exactly `u32@0x10` words.
- A `&YA1` magic is present at exactly `u32@0x04 + 0x100`; exactly one `Asg&` (offset 0) and one
  `&YA1`; zero `M7R`.
- The block-record array's cells are strictly increasing and lie in `[0x807, 0xEDEE)`.

### Rare and boundary fixture coverage

The current preserved population reaches a Spellbook count of 29, a nonempty nested SpellEffect
graph, and a dead-actor count of 22. It reaches at most seven populated references among the twelve
worn-equipment slots; its largest terrain collections contain 4,117 block records and 197 cell
records, with no 256-wide map fixture. These are finite population bounds, not representability
limits. — SAV-BOUNDARY-364

The 29-count Spellbook stores 28 consecutive `Spell` references for logical slots 1 through 28.
Each is a prior-class/new-object record tagged `0x8009`, so the sequence advances object identity
rather than aliasing one Spell repeatedly. — SAV-SPELLBK-365

One top-level SpellEffect-list record is a `SpellTransport` whose typed graph contains a
`PointEffect` and `Effect_DirectDamage`; its typed `AreaEffect` reference is null. Base
`SpellEffect` and non-null `AreaEffect` instances remain unwitnessed. — SAV-EFFECTGRAPH-366

Exact byte-0-to-EOF closure of these documents is established, but original load, loaded-state
projection and fresh original resave are not. — SAV-BOUNDARY-367

## Notes & open questions

- **Corpus limit:** the framing rows were derived on 4 saves of `10.alm`; `SAV-FLAG-027`,
  `SAV-CITY-030` and `SAV-SESS-031` add 12 saves of
  `10.alm` and `20.alm` including one between-mission save, one campaign, one difficulty, one
  release. Nothing here reaches a
  **256×256** map, whose rows 238..255 fall past the sweep window's high end — that is the
  measurement that would test the window from the other side, and it costs one play session. It
  would also test the runtime-id law (`SAV-ID-015`) on a second map.
- **A save records what the participant has seen, in the uncompressed tail** (`SAV-FOG-061`). The
  tile plane is not serialized anywhere in the document programme — that part of `SAV-SEEN-055`
  (retracted headline; this negative clause survives) was
  right — and the load arm rebuilds the terrain from the `.alm` before applying the block records
  (`SAV-LOAD-057`); the explored bit arrives afterwards, from the tail's `Fog` section. A consumer
  that discards explored terrain across a save/load diverges from the original. What the original
  *draws* from the restored bit is not measured: the render gate ORs four neighbouring corner words
  (`TERR-TILE-079`), so the drawn extent may exceed the set cells.
- **Nothing here reaches a body over 65 535 words**, so the `u16 count` before the block array and
  any wide form of it are untested at scale.
- **Four opcode values (`0x00`, `0x7f`, `0x80`, `0x81`) never occur** — a census: 0 of 12 516
  opcode positions. Their meaning is no longer Unknown: all four are legal and their behaviour is
  read off the codec's own arms (`SAV-CODEC-022`); the shipped encoder simply never emits them.
- **Members are decoded, and the actor's own are now partly named** (`SAV-UNITFLD-049`).
  The fourteen `u16` at `Unit+0x84`..`+0x9e` read `body, reaction, mind, spirit, speed, ?, ?,
  capacity, health, healthMax, healthRegen, mana, manaMax, manaRegen`; the alignment is fixed by
  `UNIT-CTOR-004`'s two constructor immediates `0x64` and `0x32` landing at `+0x98` and `+0x9e`
  on 82/82 `Human` records. **Health at `+0x94` is SIGNED** and a corpse carries it negative.
  Capacity at `+0x92` is 300 on 176/176 `Unit` records whatever their body, being computed in
  the constructor from the default body of 30 and never recomputed on that arm. `+0x8e` is the
  actor's own carried weight and `+0x90` its derived load — own weight plus half the inventory
  container's running weight sum (`ITEM-LOAD-005`, `HERO-SIGHT-007`); `+0xa4` is the `u16` sight
  radius in 1/256-cell units, high byte `+0xa5` (`AI-SIGHT-092`, `HERO-SIGHT-007`) — closed without
  new derivation by `SAV-635`. All three are read back out of the stream by the load arm and are
  not recomputed there: one original-written stream carries a load of 181 against an own weight of
  178 and an empty container — a value no recompute could produce — and the original loaded that
  stream, played, and wrote 181 back, twice (`SAV-794`). So a consumer must round-trip what the
  file supplies. `+0x8e` is maintained incrementally, not stored once: besides the derive, which
  has no direct caller and is reached only through a vtable slot, `FUN_004f36f7` has 13 direct
  callers, adds a 16-bit delta to `+0x8e` and then recomputes `+0x90` with the same arithmetic
  (`SAV-792`). That arithmetic is 16-bit and its 64 000 threshold is a signed compare, and the
  speed penalty reads the result with `MOVSX`, so a load that wraps past 0x8000 reads negative and
  cancels the penalty instead of deepening it; no preserved record reaches that range
  (`SAV-793`). Over 826 `Unit` records in the 31 preserved streams no writer of this project
  reaches, `+0xa4` takes six values, `1024 x {1, 1.25, 1.5, 1.75, 2, 2.25}`, with a zero low byte
  on all of them — in the 1/256-cell unit above, whole-cell sight radii of 4 to 9, which is the
  `scanRange` column the non-hero streamer registers; over 448 `Human` records in the same streams
  it spans 1331..2037, and the hero derive reproduces the stored value on 442 of them (`SAV-795`).
  `+0xa4 = 0` occurs in none of those 1 274 records. It occurs three times in the wider corpus, in
  one generated candidate of this project's, a copy of it, and the original's own resave of the
  first, all three agreeing in every scalar column: the original preserved that field rather than
  choosing it, and the constructor's own sight word is `0x0500` = 1280, not 0 (`SAV-796`,
  `claims/retracted.md`). `+0x18` is a different surface, `Token+0x18` (`SAV-636` below), and
  `+0x6c` is a countdown only while dying.
  `+0xa2/+0xa3` are the separately persisted regeneration remainder bytes
  (`SAV-REGENSTORE-529`, `SAV-REGENWIRE-532`), not unknown scalar padding.
  What a record carries is four distinct constructs: the inventory at `+0x7c` (presence-flagged,
  present on every measured actor; a new Sack adopts the supplied container on death, while
  an existing Sack drains and deletes it), the twelve worn
  armour references at `+0x198 + 4i` for `i = 1..12`, the `Spellbook` at `+0x140`
  (presence-flagged), and a thirteenth reference at `+0x1e4` (`SAV-CARRY-050`).
- ~~**A writer must know that the top level does not chain** (`SAV-TOPLVL-052`)~~ — **refuted by
  `SAV-PLDIARY-054`.** The apparent terminator was `Player::Serialize`'s missing
  `Diary` tail call; the top-level list writes its own `u32` count and then that many full records.
  What survives of the row: the human participant is the FIRST record - slot 1, `+0x28 == 0` -
  and its subtree is the whole of the player's own side, so a reader that wants
  the party needs no scan. Every actor's head reference at file 33..36 (`+0x14`) must equal the
  identity key of its owning `Player`: measured 258/258 with no exception (`SAV-OWNER-048`).
  The earlier apparent terminator is therefore not a boundary. Programme walks, not a pattern
  scan, reach later counted records, including the exact dead list below.
- **An object removed from play leaves the owner graph and stays in the exact top-level dead list**
  (`SAV-DEATH-051`, corrected by `SAV-DEADLOAD-124`). Its exact framing yields 101 observations:
  stages 3/4/5 = 5/71/25, and stage 5 coincides with runtime id 0 on 25/25. Do **not** test
  `health == -10001`: stage-5 health also reads -10007, -10011, -10014 and -10017. The complete
  decay consumer returns unchanged for every signed health below -10000; which external writer
  leaves each non--10001 residue remains Unknown (`SAV-DEADLOAD-127`, `SAV-DEADLOAD-129`).
  A later dead-actor load pass does rebind the raw terrain key at position `+0x08` through the
  archive identity map (`SAV-DEADLOAD-125`). Still unread: actual post-load cell occupancy; whether
  the optional pre-entry simulation tick first consumes, replaces or destroys a loaded
  Building/Sack Position before the class-specific entry sender (`SAV-POSLOAD-140`,
  `SAV-POSTLOAD-221`); live ordering for nested Items/Effects before entry or pickup;
  the meanings of the
  objects and arrays embedded at `Diary+0x04`, `Diary+0x18`, `Group+0x4c` and `*(Unit+0x140)`
  beyond their known framing — `Unit+0x15c`, `Unit+0x178` and `*(*(Unit+0x158)+0x90)` close below
  (`SAV-630`–`SAV-634`). No
  `Serialize` body in the static 35-descriptor population remains unread:
  `SAV-CLASSSER-172`…`SAV-CLASSSER-177` close the twenty schema-1 Token-lineage bodies, and
  `SAV-CONTSER-188`…`SAV-CONTSER-194` close the fifteen non-`Token` bodies.
  `SAV-PRODROOT-204`, `SAV-PRODPOP-205`, `SAV-PRODGEN-206`, `SAV-PRODDIRECT-207`,
  `SAV-LOADRESAVE-208` and `SAV-PRODNOHIT-209` classify the enumerated static producer
  routes in the class-set section above.
  Run-time-built, copied, aliased or separately implemented producers remain Unknown, and the
  ten-class no-hit is not a format-wide absence.
- **Load does re-run terrain ingest before applying saved terrain.** The complete arm constructs
  terrain from the named ALM, publishes it, then loads block records and the cell table. Cell nodes
  are overlaid or inserted without clearing the constructed hash (`SAV-LOAD-057`,
  `SAV-CELLLOAD-108`, `SAV-CELLLOAD-109`).
- **The shape byte clear arm is observed in preserved no-world saves.** Those paths retain the
  roster while omitting the world half, and the reader's post-load path restores and detaches the
  actors. Acceptance of a newly generated no-world programme by the original remains Unknown
  (`SAV-CITY-030`, `SAV-POSTLOAD-220` (amended), `SAV-WRITER-284`).

## Conditional command route to trailer trace toggles

Command 46 parameter 80 calls 0053cdc0 after resolving its Player — a decisive
subset of the route's conditions; `byte[cmd+4]==0` and global `[0x5f21c4]!=0`
also gate it. Unsigned Player+68 must exceed 50 ("unsigned" rests on the `jbe`
mnemonic; the twelve test vectors cover only gate bytes 50/51); subcommands
3/19 then toggle the first/second dword of world+118's 400-byte object. Zero
maps to 1 and any nonzero value to 0. The tests stop before logging. This
supplies a static incoming route, not an ordinary UI emitter or a proved
producer of the gate byte — though the callee's own subcommand 7 prints the
engine's debug-command help, naming subcommand 19's toggle "Script tracing
on/off" (`0x5c8990`), without establishing what reaches it. SAV-694 falsifies
SAV-647's earlier "0 direct E8 callers" reading for this same callee
(`claims/retracted.md`): the route above is one direct E8 call. The other 98
dwords have no located consumer at all, by the complete displacement census
above; what they hold, and whether an ordinary runtime route reaches the two
that do, remain unresolved. — SAV-694, SAV-647, SAV-791

## Primitive reader domain and minimal encodings

The generic count writer uses u16 below `0xffff` and `0xffff`+u32 otherwise. CString
uses a byte below `0xff`, `0xff`+u16 through `0xfffd`, and `0xff`+`0xffff`+u32 starting at `0xfffe`;
`0xff`+`0xfffe` is reserved by the CString reader as a Unicode-mode sentinel. The
primitive readers also accept nonminimal encodings of 1: generic wide, CString
word-extended and CString dword-extended. They consume 6/3/7 bytes respectively.
These are successful-primitive-IO paths, not proof that any large allocation,
full malformed archive or unobserved map/Player graph survives original LOAD.
— SAV-758, SAV-WLIST-040, SAV-FULLREAD-252

## Player percentage command versus archive transport

Player+0x58 has a non-constructor producer in command `0x46` parameter 3. Signed
input 0/1/2 maps to 100/50/0; 3..100 is copied; negative or above 100 leaves the
prior field unchanged. It is not a clamp or a universal 0..100 invariant. The
actor walk then dispatches each member's virtual `+0x50` derive unconditionally.
Unit `004f5946` (`UNIT-DERIVE-003`) consumes the same field through `actor+0x14`
as the percentage `HERO-MP-006` already names the AI's mana floor. SAVE/LOAD
transport does not apply this command normalization. Parameter 1, not 3, is the
withdraw arm. Player+0x50 and other raw settings remain unresolved. — SAV-726,
SAV-666, SESS-PARAM-017 (partially retracted)
## Session self-pointer and map-constant repair

The raw 2,508-byte span beginning at `session+0xa9bc` includes a self-pointer at `+0xa9c0`,
stream offsets 1852..1855. LOAD's `004d14f7` call to `0053b870` replaces it with ECX.
The complete thunk has no register-array access. Earlier `004d143f` calls `004e3591`
after session restore. Authored check opcode `0x10002` writes node+0x48 into the current
compiled slot and advances it; opcode `0x10003` does neither on the selected arm.
Successfully compiled preceding normal checks and constants advance that slot,
while resolution-rejected normal checks do not. Whole-LOAD execution, arbitrary
indices and all adjacent helper reachability are not established by local vectors.
— SAV-710, SAV-711

## Recipient publication mask

Token+0x18 is a u16 recipient mask in the located emission paths, both set and
cleared there. Helper 004f2885 returns its intersection with Player+0x2c;
004f28a8 returns 1 when that intersection is zero; 004f28d0 ORs the recipient's
bits into the field. Reset 004f2639 zeroes it along with +0x1c and +0x4, not
alone. Five sites in four functions clear one recipient's bit with the
complementary AND idiom, and one of them republishes to the same recipient
immediately after clearing it, so a zero bit is revoked state, not a proof a
recipient was never published. Sender 004e7de3 uses the negative test then the
setter (a dedup guard), 004e873b tests the intersection directly, and a third
function inlines the same OR-publish step twice, broadcasting to every bit
when its recipient argument is null — the discriminator against a class flag
or Boolean. This does not establish a universal fog-of-war meaning, a class
discriminator or the full lifetime of every class. The existing width and
literal transport remain. Over 1,420 decoded records in 45 preserved streams
the field is 2 on 1,404 and 0 on 16, with no other value. The terminal stage
does not decide which: of 26 records at stage 5 exactly half carry 0, and the
same three dead units appear at mask 2 in one save of a campaign and at mask 0
in another. The three stage-0 records carrying 0 are all in the one save a
preserved manifest records as the resave made after three mercenaries were
hired.
— SAV-678, SAV-653, SAV-664, SAV-797
