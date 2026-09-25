# PARTY — ownership, the roster, and what crosses a map change

Not a file format. Claims about what a player owns and what of it is still the
player's on the next map. The ledger sits between four ledgers that each own
one piece and none of which owns the join: [`tavern.md`](tavern.md) owns the
mercenary pool, [`hero.md`](hero.md) the hero's stats, [`alm.md`](alm.md) the
map's type-5 group roster (`ALM-GRP-041`), and [`sav.md`](sav.md) the save
container. These claims add the carrier: the engine has no party object, and
one thing that crosses a mission boundary is a serialized object graph; one
thing, not the only one (`PARTY-PERSIST-014`). Format of this file:
[registry.md](registry.md).

## Terms

Two vocabularies meet here and are kept apart.

- The party is two containers on the `Player` plus a pointer
  (`PARTY-ORIGIN-010`, the claim to read first). The group list `+0x24` is what
  a save stores, the flat index `+0x20` is what the placement walk reads, and
  `+0x20` is rebuilt from `+0x24` on load. Which one a consumer implements
  decides which half of the game breaks.
- A `Player` (`0x70` bytes, `CRuntimeClass` at `0x5c24f8`) is the simulation's
  owner object, one per map type-5 slot.
- A group (`0x48` bytes) is a container of actors held by a `Player`.
  `ALM-GRP-041`'s "group" is the map record a `Player` is built from, not this
  object.
- The world is the server object at `[0x005cd758]`; `SESS-*` claims call it
  the server singleton.

## Ownership, roster and the client cull

| ID | Claim | Confidence | Status | Evidence |
|---|---|---|---|---|
| PARTY-OWN-001 | Ownership is one pointer, `actor+0x14`, to a `Player`; the `Player` carries two identities, its slot `+0x04` and the map's type-5 id word `+0x08`. | High | ● active | [EXP-0077](../experiments/EXP-0077-party-and-save/) |
| PARTY-ROSTER-002 | There is no party object and no member list: the roster is `Player` → group collection → group → the group's own actor list, and membership is containment. | High / Medium / Unknown | ● active (amended) | [EXP-0077](../experiments/EXP-0077-party-and-save/) |
| PARTY-FLAG-003 | On the client, side membership is the dword flag word `CUnit+0x18c`: the end-of-mission cull keeps bit 0, the tavern tally reads bit 4 (mercenary), and bit 1 splits that bucket. | High / Medium / Unknown | ● active | [EXP-0077](../experiments/EXP-0077-party-and-save/) |
| PARTY-CULL-004 | At the end of every mission the client document is cut down to one player and its own surviving player characters; everything else is destroyed. | High | ● active | [EXP-0077](../experiments/EXP-0077-party-and-save/) |

### PARTY-OWN-001

- The load arm of `Player::Serialize` (`FUN_00511089`) walks the player's own
  groups and writes `actor+0x14 = player` for every actor it finds
  (`00511404`, `MOV dword ptr [EAX + 0x14],ECX`).
- Every consumer that asks "is this mine" compares that field against the
  local player object: the client's end-of-mission cull at `00420360`
  (`CMP EDX,dword ptr [ECX + 0x9b4]`) and the tavern's live-unit collector at
  `004668b4`/`004668bd` (`SHOP`/`MERC-DEATH-006` already used the second).
- The `Player` carries two identities, and they are not the same number:
  - `+0x04` is `slot+1`: the index it is re-seated at in the document's player
    array (`00420aab`) and the value the wire uses to name a player
    (`004200b0`, `004ea568`);
  - `+0x08` is the map's own type-5 id word, the value the owner lookup
    `FUN_004fb534` compares its argument against (`004fb555`,
    `CMP dword ptr [ECX + 0x8],EAX`).
- `ALM-GRP-041` publishes that word from the loader's side; this claim is the
  consumer's side.

**Confidence.** High. The write is one instruction inside the serializer, and
the three reads are named instructions in three unrelated modules. The two
identities are told apart by their two distinct consumers, not by their values.

### PARTY-ROSTER-002

- The roster has three levels:
  `Player → group collection → group → the group's own actor list`.
- `Player+0x24` is a collection of groups, serialized by `FUN_005102f4`. On
  load it reads a count, then `new(0x48)` + `FUN_0050f81b` per element
  (`00510396`, `005103b3`).
- Each group is serialized by `FUN_00511938`. It writes its actor list through
  `FUN_00527ac0` (an object list: count via `FUN_0044a4b0`, then
  `ar << element` per node), then its own group id at `+0x1c` and two object
  references at `+0x40`/`+0x44` that are id-fixed-up on load (`FUN_00527e40`,
  `FUN_00527d30`).
- The players hang off one global list, `[0x00609544]`, whose serializer is
  `FUN_0051102f`.
- `SAV-OBJ-016` found inventories and units nesting under their owners with an
  untagged ~100-byte structure between sibling units: that structure is a group
  record.
- Count: the corpus's five `Player`s are the map's five type-5 slots
  (`SAV-OBJ-016`; the count is the fourth head u32, `SAV-HEAD-025`). Cap: a map
  may declare at most 16 (`ALM-META-025`, `ALM-GRP-041`'s sixteen diplomacy
  words).
- No cap on groups per player and none on actors per group was found. Both
  containers are unbounded MFC collections, and no compare against a limit was
  seen on either insert path.

**Confidence.** High for the three-level structure and the group record's
fields: every level is a serializer read end to end, and the shape it writes is
the shape `SAV-OBJ-016` measured in the file. Medium for the 5 and the 16: both
are facts about shipped maps carried from `alm.md`, not bounds read off an
instruction.

**Unknown.** A per-group or per-player member cap. The absence was found by
reading the two insert paths, not by an image-wide sweep.

**Amended.** `SAV-GRPFLD-060` (EXP-0149) supersedes the group-id clause at
`+0x1c` ([`retracted.md`](retracted.md)). `+0x1c` is an authored id only for
groups the map supplies. The group constructor `FUN_0050f81b` writes `+0x40`,
`+0x44`, `+0x3c` and the embedded list at `+0x20`, never `+0x1c`, so for a
group the engine creates at run time the serializer writes out uninitialised
memory. The three-level roster, the two serializers and the load-side
construction are untouched.

### PARTY-FLAG-003

- `FUN_004667f0` partitions the local player's live units into two buckets in
  one pass:
  - bucket A takes those with `TEST AL,0x1` (`004668c7`) and corpse stage
    `byte+0x15a == 0`;
  - bucket B takes those with `TEST AL,0x10` (`004668f1`) and the same stage
    test. Bucket B only is split into two counters by `TEST AL,0x2`
    (`00466920`): `+0x2c` if set, `+0x28` if clear.
- The tavern's pool tally `FUN_004885c0` reads bucket B: it takes the array at
  `0x5ea0d8` and the count at `0x5ea0dc`, the `+0x18` / `+0x1c` fields of the
  collector at `0x5ea0c0`, and subscripts it by the mercenary type
  `byte+0x15b` (`00488687`). So bit 4 is "mercenary" and bit 0 is not. The
  end-of-mission cull keeps bit 0 (`PARTY-CULL-004`).
- `MERC-DEATH-006`'s parenthetical gloss "`+0x18c & 1` (hero)" is right about
  the shipped effect and wrong about the field: the bit means player character.
- `FUN_0047aee0` sets the word as a whole-word literal `0x29`, or `0x2b` when
  `byte[EDI+0x74] & 0x40` (`0047b30b`, `0047b31c`). That routine builds a hero
  out of a character record it has just read with eighteen `CFile::Read`
  calls.
- Bits 2, 3 and 7 are `OR`ed in at other sites (`0047b332`, `00411f14`,
  `004185ac`). Bit 5 is set only by those two literals.
- Instrument: `EnumRefs disp:18c`, whole image, 232 hits over 81 owners, 0 in
  orphan or undisassembled code. Its blind spot applies here: the word is
  copied wholesale from one drawable to another at eight sites in
  `FUN_00421f46` and three in `FUN_0047c590`, and any `REP MOVSD` of a
  containing struct would carry no displacement at all and be invisible.

**Confidence.** High for what the three bits do: each of the five tests above
is a named instruction, and the bucket-to-consumer binding is fixed by the two
static offsets the tally reads, not by which bucket looks apt. Medium that bit 0
means "player character" in general: one writer sets it, and this install ships
one campaign, so a second kind of bit-0 unit would not have been seen.

**Unknown.** The remaining bits, and whether the server actor carries a
counterpart field at all: none was looked for.

**Amended.** `PARTY-ENDCULL-026` reads the server's end-of-mission cull. There
is no counterpart flag word: the server's membership predicate at a mission
boundary is a range test on the actor's typeID `word+0x0e`.

### PARTY-CULL-004

- `FUN_004201c9` has a single caller, `FUN_00473110` at `00473e75`
  (`MERC-DEATH-006`). It walks the document's object map `doc+0x9b8`, and for
  each `CUnit` (`IsKindOf` against the `CRuntimeClass` at `0x599208`,
  `0042031a`) applies three tests, all of which must pass:
  - `[+0x18c] & 1` (`0042033e`);
  - `(i16)[+0xfc] > −10` (`0042034f`, `CMP EDX,-0xa` / `JLE`);
  - `[+0x14] == [doc+0x9b4]` (`00420360`).
- A survivor is reset in place: `byte+0x15a = 0` (corpse stage),
  `byte+0x84 = 0`, `dword+0xa0 = 0`, then `vt+0x14(0)`
  (`0042036b`…`00420393`).
- Everything else is destroyed through `vt+0x04(1)` and unlinked: every
  `CUnit` that fails any test, every non-`CUnit` object, the second map
  `doc+0x9d4` in its entirety, `doc+0x3f3c` and `doc+0x80`.
- The player array `doc+0x9a0` is then emptied from index 1 upward of every
  player except `doc+0x9b4` (`00420911`, `CMP EAX,dword ptr [EDX + 0x9b4]` /
  `JZ`), and the survivor is written back at index `[player+0x04]`
  (`00420aab`).
- The kill list is built first and drained second: a `CWordArray` of map keys
  at `[EBP-0x28]`, appended by `SetAtGrow(m_nSize, key)` and drained with
  `RemoveAt(0,1)`.

**Confidence.** High. The three tests, the three resets and the two loops are
named instructions in one routine with one caller. The "everything else" clause
is not an inference: the `IsKindOf`-false arm at `004203bb` appends to the same
kill list as the failed-test arm at `00420398`.

## The carry across a map change

| ID | Claim | Confidence | Status | Evidence |
|---|---|---|---|---|
| PARTY-CARRY-005 | What crosses a map change is one compressed `CArchive` graph with one root actor plus sixteen bytes of `Player` scalars, the same bytes in the network packet and the `.chr` file. | High | ● active | [EXP-0077](../experiments/EXP-0077-party-and-save/) |
| PARTY-LOSS-006 | The carry importer re-homes the actor: a carried character arrives at full pools, off the map, with no position, order or route, and named after the account. | High / Medium / Unknown | ● active | [EXP-0077](../experiments/EXP-0077-party-and-save/) |
| PARTY-MERC-007 | Mercenaries do not cross a map change as objects; what crosses is a count per type in the tavern's mission record, so one party needs two persistences. | High | ● active | [EXP-0077](../experiments/EXP-0077-party-and-save/) |
| PARTY-SESSION-008 | Apart from three zeroing stores, `FUN_004cfdb0` is the only writer of the world, player-list and actor-registry globals; that the server is rebuilt at every mission boundary is contested. | High / Medium / Unknown | ● active (contested) | [EXP-0077](../experiments/EXP-0077-party-and-save/) |
| PARTY-GROUP-009 | A client's group is the one whose `Player` the client is assigned, named by slot and nothing else; hiring a mercenary changes no ownership field. | High / Medium | ● active | [EXP-0077](../experiments/EXP-0077-party-and-save/) |

### PARTY-CARRY-005

- The exporter `FUN_004cf510(actor)` opens a `CMemFile`, calls
  `FUN_00517c00(ar, actor)` (`WriteObject`, one root: a single `ar << actor`),
  pads the memory file to an even length and compresses it with
  `FUN_00517510`. That is the save writer's compressor, and the image has
  exactly 2 call sites for it: this one and `FUN_004cee8f`.
- It lays the result into the static packet at `0x60c7b0`:
  - `byte+0x09 = 0xbe`;
  - `word+0x07 = [actor+0x14]+0x04`, the owning player's slot;
  - `dword+0x0a` = the payload word count;
  - `memcpy(pkt+0x0e, &{player+0x38, +0x48, +0x4c, +0x54}, 16)`;
  - the compressed graph at `pkt+0x1e`.
- The same 16 + N bytes are written to `chr\<profile>\<u><u>.chr` (format
  string `0x5c5df8`) when the profile name at `player+0x64` is non-empty, and
  to `<u><u>.chr` when it is not.
- The client re-sends them: `FUN_0042008d` builds the identical `0xbe` packet
  out of `campaign+0x114` / `campaign+0x118` (`FUN_0047bd50` at `0047bda4`,
  guarded by `[ESI+0x118] != 0` at `0047bd93`). When the blob is absent it
  sends the chargen message `0x48` instead (`FUN_0041fd61`).
- The server's arm (`FUN_004d5dd8` at `004d7eb6` / `004d7f15`) can take the
  payload from the wire or read a `.chr` off disk into the same packet field
  before calling the importer.
- The importer `FUN_004cf8d8` has one decompressor call, and the image has
  exactly 2 of those: this one and the save loader.

**Confidence.** High. Both call-site enumerations are `EnumRefs callto:` on the
repaired table with `.rdata` slots included, and each returns 2 hits / 2
owners, 0 orphan. The single root is the single `WriteObject` at `004cf56f`.
The 16 bytes are one `memcpy` with a named source, and the importer assigns
those four dwords back to the same four `Player` offsets at
`004cfaa6`…`004cfac1`.

### PARTY-LOSS-006

After `ReadObject`, the importer `FUN_004cf8d8`:

- overwrites the actor's own name with the receiving `Player`'s CString
  (`004cfa9b`, `[actor+0x80] = player+0x18`);
- frees three sub-objects and installs fresh ones:
  - `actor+0x154` freed and replaced by `new(0xb4)` + `FUN_00545c00`
    (`004cfad9`…`004cfb38`);
  - `actor+0x158` released via `FUN_00517c40(1)` and replaced by `new(0x94)` +
    `FUN_0052cbb0`;
  - `actor+0x10` freed and replaced by `new(0xc)` + `FUN_00544530`;
- zeroes five fields `+0x5c`, `+0x64`, `+0x44`, `+0x68`, `+0x40`
  (`004cfc51`…`004cfc79`);
- calls `FUN_004f4b39`;
- copies `word+0x96 → word+0x94` and `word+0x9c → word+0x9a`, a pair of
  "current ← maximum" restores;
- sets `byte+0x13c = 0`, clears bit 3 of `byte+0x4c` (`AND DL,0xf7`), sets
  `word+0x18 = 0`, and calls `vt+0x50`.

The actor's name comes from the account rather than from the object. The four
`Player` scalars `+0x38`, `+0x48`, `+0x4c`, `+0x54` ride separately because
they are not in the actor. The complementary loss is `PARTY-MERC-007`.

**Confidence.** High for the field list and the replacements: twenty-odd
consecutive instructions in one routine, each naming its offset. Medium for the
reading of `+0x94`/`+0x9a` as current-from-maximum and of
`+0x40`/`+0x44`/`+0x5c` as map-bound state: the pattern is the same one the
world `Serialize`'s no-map arm applies to every actor it restores
(`SAV-ROSTER-024`), which is corroboration from a second site rather than a
reading of either field.

**Unknown.** What `+0x154`, `+0x158` and `+0x10` are; only that they are
per-map and are never carried.

### PARTY-MERC-007

- The reason is structural rather than a rule.
- Mercenaries are not in the carried graph: the graph has one root, and that
  root is a player character.
- They do not reach it as references: they hang off the `Player`'s groups
  (`PARTY-ROSTER-002`), and the `Player` is not serialized into the blob; only
  four of its dwords are.
- On the client the same units are dropped by `PARTY-CULL-004`, because they
  carry `+0x18c` bit 4 and not bit 0 (`PARTY-FLAG-003`).
- What is carried instead is the count per type, in the mission record the
  tavern keeps: `MERC-DEATH-006`'s merge, and the tavern's own rebuild of the
  whole roster on the next town visit (`MERC-CMD-007`, command `0x37`).
- A consumer implements two different persistences for one party: the hero as
  an object, the squads as fifteen integers.

**Confidence.** High. Each half is a published, instruction-level claim in
another ledger. This claim adds that the two are exhaustive: the carried
payload is one `WriteObject` and sixteen bytes, so there is no third channel.

**Amended.** `MERC-POOL-011` refines the width: it names that array's storage
and the width of its elements.

### PARTY-SESSION-008

- The world is the object `[0x005cd758]`; `SESS-*` claims call it the server
  singleton.
- `FUN_004cfdb0` is the only writer of the three globals the roster lives in:
  `[0x005cd758]` (the world), `[0x00609544]` (the player list) and
  `[0x00609558]` (the actor registry). The only other writes are
  `FUN_004cec1d`'s three `MOV …,0x0`.
- Published reading: nothing survives in memory, because the single-player
  server is destroyed and rebuilt at every mission boundary. `FUN_004762a0`
  allocates `0x174` bytes, calls `FUN_004cec1d` on the old world and
  `FUN_004cfdb0` on the new one. A carry is therefore a serialization or it
  does not happen.
- Contesting reading (EXP-0100, `PARTY-PERSIST-014`): the premise has no
  instruction behind it. `FUN_004762a0` makes both constructor calls on the
  object it has just allocated (`004762d6`/`004762ee`, both `MOV ECX,EAX`), so
  it builds one world and destroys none. It has five call sites, and the
  per-mission starter `FUN_00477c00` is not one of them; `FUN_004d00e9` runs on
  the existing object (`00477dc8`). What stands is the global-writer
  enumeration, not the cadence, and therefore not the conclusion "a carry is
  therefore a serialization or it does not happen".
- What replaces the rebuild is `FUN_004d00e9`, read whole by `SESS-LOAD-009`,
  which holds the mission-number parse and the save-versus-map branch. The
  routine finishes by zeroing `world+0x00` and `world+0x04`, the two counters
  `SESS-TICK-004` names and `SAV-HEAD-025` measures in the file, and by setting
  `world+0x2c = 1` (`004d0519`). That is the flag the save writer tests to
  decide whether a world half exists at all (`SAV-SHAPE-023`).
- The mission `.ini` it builds a path for, `World\Mission\<n>.ini`, does not
  ship: no `.ini` entry exists in any of the eight `.res` archives (0 of 4 080
  entries), and no `World` directory exists on disk.

**Confidence.** Medium for the lifecycle, demoted by EXP-0100: the
global-writer enumeration is `EnumRefs refto:` over three addresses with 0
orphan hits and stands, but "destroyed and rebuilt at every mission boundary"
was read off `FUN_004762a0`'s two constructor calls, and both are on the new
object. High for `world+0x2c`'s single write. High for the `.ini`'s absence: a
whole-corpus enumeration of archive entry names plus a directory listing, the
instrument class `AGENTS.md` calls immune.

**Unknown.** What `FUN_004e0b0c` does when the file is missing, and therefore
what a shipped campaign loses by its absence.

**Amended.** Contested by EXP-0100 (`PARTY-PERSIST-014`,
[`retracted.md`](retracted.md)). The lifecycle half, graded High when
published, is Medium. The retraction entry (verdict AMBIGUOUS) finds no
instruction behind the premise and leaves the conclusion's falsity
unestablished: whether the shipped campaign's mission-to-mission edge reaches
the teardown `FUN_00476340` was not read there.

### PARTY-GROUP-009

- The client stamps `word[pkt+0x05] = word[[doc+0x9b4]+0x04]` on both messages
  it sends about its character: the carry `0xbe` (`004200b0`) and the chargen
  `0x48` (`0041fdb2`). That value is `slot+1` (`PARTY-OWN-001`).
- The server's arm resolves the slot to a `Player` and passes it to the
  importer as the second argument. The imported actor's name, money and three
  other scalars are then written from and to that object (`PARTY-LOSS-006`).
- Hiring a mercenary changes no ownership field. `MERC-HIRE-003`'s spawn loop
  `FUN_005056f1` creates one object per head on the server, and the object's
  owner is the player it was created for. The hire flag
  `dword[record+0x88 + (t−1)·4]` lives in the campaign mission record, not on
  any unit.
- The three ways a unit becomes a player's are one mechanism and two origins:
  the field is always `actor+0x14`, and it is written by the map loader for
  what the map gives, by the spawner for what is bought, or by the importer for
  what is brought.

**Confidence.** High for the slot stamp and the hire's non-effect on ownership:
both are named instructions, the second in a published claim. Medium that the
resolution in the server's arm is by slot: the arm's player argument was read
as a live value at the call site, and the lookup that produced it was not
traced to its own compare. The discriminator is one listing of
`FUN_004d5dd8`'s `0xbe` prologue.

## Where a player's units come from

| ID | Claim | Confidence | Status | Evidence |
|---|---|---|---|---|
| PARTY-ORIGIN-010 | The party is two `Player` containers: the group list `+0x24` is saved, and the flat index `+0x20` is derived from it on load and read only by the placement walk. | High | ● active | [EXP-0100](../experiments/EXP-0100-party-origin/) |
| PARTY-WRITE-011 | Eleven routines write the player's two containers by call, and the group-list and group-member appends share one owner set, so the two containers are always written together. | High / Medium / Unknown | ● active | [EXP-0100](../experiments/EXP-0100-party-origin/) |
| PARTY-INSTALL-012 | Giving a player a unit is one straight-line sequence of five writes, duplicated rather than shared at every entry point; a mercenary spawn writes four, without the hero pointer. | High / Medium | ● active | [EXP-0100](../experiments/EXP-0100-party-origin/) |
| PARTY-GATE-013 | No hero, no mission: `FUN_004d303e` rejects a client whose `Player+0x34` is null, and that is the only membership test on the way in. | High | ● active | [EXP-0100](../experiments/EXP-0100-party-origin/) |
| PARTY-PERSIST-014 | The single-player carry has a memory arm: `FUN_004e2327` can reuse the surviving `Player` for the new map's slot 1, and the campaign start selects that arm. | Medium | ● active | [EXP-0100](../experiments/EXP-0100-party-origin/) |

### PARTY-ORIGIN-010

- The `Player` constructor `FUN_004faa81` builds both containers (`004fabf3`,
  `004fac31`):
  - `+0x20` is `new(0x20)` + `FUN_0050fb5e`, a four-byte vtable (`0x59caf0`)
    wrapping one `CObList` at its own `+0x04`: the player's flat actor index;
  - `+0x24` is `new(0x1c)` + `FUN_0050fa1d`, the group list `PARTY-ROSTER-002`
    names.
- Besides the two containers, the party is a pointer and two back-pointers.
  `+0x34` is the hero pointer, zeroed there (`004fab31`). `actor+0x14` is the
  owning player (`PARTY-OWN-001`) and `actor+0x70` the owning group.
- `+0x20` is derived, not stored. `Player::Serialize`'s load arm iterates
  `[player+0x24]`'s groups and, for every actor in every group, appends it to
  `[player+0x20]` and stamps `actor+0x14` (`005113bd`, `005113f6`, `005113f9`,
  `00511404`). That is why `SAV-PLAYER-028`'s seventeen fields contain `+0x24`
  and `+0x34` and not `+0x20`.
- `FUN_004d403c` reads only `[player+0x20]+4` (`004d4231`…`004d42f7`).
- A consumer that implements the groups alone starts every mission with nobody
  standing on the map, and one that implements the flat index alone loses the
  party at the first save.
- The two globals `[0x00609544]` and `[0x00609558]` are the same class
  (`FUN_004cfdb0` at `004cffe5`/`004d0040`), so "my units" and "all units" are
  the same shape at two scopes.

**Confidence.** High. Every field is a named instruction in a constructor or a
serializer read end to end. The derivation is not an inference but the load
arm's own direction of copy, groups → index. The rival "one collection seen
twice" is excluded by the two `new`s of different sizes with different
constructors.

### PARTY-WRITE-011

- Instrument: `EnumRefs "re:CALL 0x0050fbee"` (the flat index's append), 32
  hits, 20 owners, 0 in orphan or undisassembled code. It prints the
  instruction that loaded `ECX`, so a hit with no function name cannot be lost.
- Twelve hits load `ECX` from a `+0x20` displacement:
  - `FUN_004d3755` ×2 (chargen, `SESS-HERO-014`);
  - `FUN_004d5dd8` (the `0xbe` carry arm, `PARTY-CARRY-005`);
  - `FUN_004e26bb` (the map's own type-6 spawner, `MISSION-ARM-006`);
  - `FUN_005056f1` (the mercenary spawn, `MERC-HIRE-003`);
  - `FUN_00511089` (`Player::Serialize`, the load arm);
  - `FUN_004d403c` (the `.ini` "Humans" arm);
  - `FUN_004d1e14`, `FUN_004d8fcd`, `FUN_004f164c`, `FUN_004feadb`,
    `FUN_00504da1`, which EXP-0100 did not read.
- The group list's append (`EnumRefs "re:CALL 0x0051a050"`: 13 hits, 11
  owners, 0 orphan) and the group's own actor append (`callto:50f950`: 15 hits,
  11 owners, 0 orphan) have the same owner set. That is the evidence that the
  two containers are always written together.
- `callto:51a200`, the raw `CObList::AddTail` beneath the three append helpers,
  returns 4 hits, 4 owners, 0 orphan: the three helpers plus `FUN_004d1f42`. So
  nothing else reaches these lists by call.
- Blind spot: a wholesale `REP MOVSD` of a `Player` carries neither call nor
  displacement and is invisible to every sweep here.

**Confidence.** High for the enumeration being complete by call: three sweeps
on the repaired table, 0 orphan hits in each, and the raw-`AddTail` sweep closes
the route around the helpers. Medium that these eleven are the complete set of
origins: five were read at instruction level and six were classified by their
append site alone.

**Unknown.** Whether any path copies a whole `Player`.

### PARTY-INSTALL-012

- The `0xbe` carry arm (`004d8024`…`004d80de`) writes, in order:
  1. `actor+0x04 = FUN_004d9fed() & 0xffff`, a fresh runtime id;
  2. `actor+0x14 = player`;
  3. `player+0x34 = actor`;
  4. `[player+0x20].AddTail(actor)`;
  5. `new(0x48)` + `FUN_0050f81b` → `[player+0x24].AddTail(group)` →
     `FUN_0050f950(group, actor)`.
- The chargen arm `FUN_004d3755` writes the identical five
  (`004d3e8d`…`004d3f02`, `player+0x34` at `004d3f2c`).
- The mercenary spawn `FUN_005056f1` writes four of them (`005057c6`,
  `005059b1`, `005059be`, `005059ca`): everything but `+0x34`, which is what
  makes a mercenary not a hero.
- A carried character therefore arrives in a group of its own, never in the
  group it left.
- `FUN_004d3755`'s first arm is a repair rather than a creation. On
  `player+0x34 != 0` it tests `byte[actor+0x13c]`:
  - if that is 0, it returns the hero untouched (`004d3783`, `004d379b`,
    `004d37a3`), the state `PARTY-LOSS-006`'s importer leaves a carried actor
    in;
  - a nonzero value restores `word+0x94 ← +0x96` and `word+0x9a ← +0x9c`,
    re-indexes into `+0x20` and rebuilds a group if `actor+0x70` is null
    (`004d37be`…`004d37eb`).
- Its one non-dispatcher caller is `FUN_004d8963`'s revive arm, which calls it
  with all-zero stats and then re-runs the placement walk (`004d8a67`,
  `004d8a73`).

**Confidence.** High for the two sequences and the branch: three routines read
at instruction level, each write naming its own offset; the duplication is a
fact about two listings, not a reading. Medium for the reading of `+0x13c` as
the flag that separates "repair a dead hero" from "leave a carried one alone":
the two states are inferred from `PARTY-LOSS-006`'s zeroing and from the revive
caller, and the field itself was not read.

### PARTY-GATE-013

- `FUN_004d303e`'s single caller is the session dispatcher's `0x04` arm
  (`004d82cb`, inside the arm at `004d828a`). It tests `player+0x34` first
  (`004d3108`). On null it logs
  `"Client %s tries to enter mission without Hero. Rejected."`
  (`0x5c609c`/`0x5c60cc`), sends `"You can't enter mission without Hero"`
  (`0x5c60d4`) and returns without placing anything.
- An empty flat index is not tested and is not an error: it places nobody and
  the mission runs.
- Past the gate the routine clears `player+0x3c` (`SAV-FLAG-027`'s outcome
  latch, `004d31b2`) and calls the walk once, under a second latch:
  `byte player+0x3d`, set to 1 immediately after (`004d3342`, `004d3358`). That
  byte is the tenth field of `Player::Serialize`'s record (`SAV-PLAYER-028`),
  so a player placed and then saved is never re-placed on load.
- The walk begins with a hygiene pass on the group list that nothing else does:
  `GetHead`, and if there is none, log `"Oops - player has no groups %-["`
  (`0x5c6164`) and create one. It then deletes empty groups from the head until
  the head is non-empty or is the last one (`004d4063`…`004d4126`).

**Confidence.** High. One routine and one loop were read at instruction level,
and three of the four arms are named by the image's own strings. The `+0x3d`
latch is a store two instructions after the call and a field of a published
serializer record.

### PARTY-PERSIST-014

- `FUN_004e2327` builds one `Player` per type-5 record, except that for roster
  slot 1 it reuses the player already in the global list when two conditions
  hold: `[0x005cd758]+0x0c == 0`, and the list holds exactly one player
  (`004e238a`…`004e23b1`, three compares in a row, `GetSize` then `GetHead`).
- That mode flag is written once, in the server's own constructor, as
  `arg0 < 2` (`004cfe46` `CMP dword ptr [EBP+8],0x2` / `004cfe4a SETL AL` /
  `004cfe53`). The campaign-start routine passes 2 (`FUN_00477450`,
  `00477454 PUSH 0x2` → `0047746a CALL FUN_004762a0`), so on the campaign path
  the flag is 0 and the reuse arm is the live one.
- The server object is not rebuilt per mission. `EnumRefs refto:5cd758` gives
  172 hits, 93 owners, 0 orphan, with exactly 3 writes: `004762f0` the
  allocator, `004cfea9` its constructor, `0047635e` the teardown's `= 0`.
- `EnumRefs "re:CALL 0x004762a0"` gives 5 hits, 5 owners, 0 orphan, and the
  per-mission starter `FUN_00477c00` is not one of them: it calls
  `FUN_004d00e9` with `ECX = [0x005cd758]`, the existing object (`00477dc8`).
- This contests `PARTY-SESSION-008`'s "a carry is therefore a serialization or
  it does not happen", whose premise was that `FUN_004762a0` destroys the old
  world. Both of its constructor calls, `FUN_004cec1d` at `004762d8` and
  `FUN_004cfdb0` at `004762f5`, are made on the newly allocated object
  (`004762d6 MOV ECX,EAX`), so that routine builds one world and destroys none.

**Confidence.** Medium. The reuse arm, the mode flag and the write enumeration
are instruction-level, and the campaign's argument is a literal `PUSH 2`, but
whether the shipped campaign's mission→mission edge always avoids the teardown
was not established. Discriminator: `FUN_00473110` is the campaign state
machine and holds both `FUN_00477c00` (`00474cd9`) and four calls to the
teardown `FUN_00476340`. The `00474cd9` arm jumps straight to the common exit
`00473193` without reaching any of them, but which arms the campaign drives, in
what order, was not read.

**Amended.** `PARTY-PERSIST-028` (EXP-0165) reads the mission-end arm: it
reaches no `FUN_00476340`, so the memory carry is the live one on the campaign's
mission-to-mission edge. Whether the `0xbe` carry also fires on that edge, and
what builds the human `Player` before the campaign's first map load, stay open.

## The purse, and AddHero companions

| ID | Claim | Confidence | Status | Evidence |
|---|---|---|---|---|
| PARTY-MONEY-015 | `Player::Player` stores a zero purse at `004fab27`; the conclusion that nothing credits it before the first mission is withdrawn, because the participant factory stores 100. | High | ● active (partially retracted) | [EXP-0151](../experiments/EXP-0151-mission-documents/), [EXP-0156](../experiments/EXP-0156-start-money/) |
| PARTY-MONEY-016 | Each `Player` owns one 32-bit purse at `+0x38`, and the two persistence paths, the save record and the mission carry, carry that same field. | High / Medium | ✔ promoted | [EXP-0152](../experiments/EXP-0152-money-cycle/) |
| PARTY-ADDHERO-017 | `[Mission<n>] AddHero[]` creates a companion player character, neither the primary character nor a mercenary, when the town view activates. | High / Medium | ● active | [EXP-0153](../experiments/EXP-0153-addhero-consumer/) |
| PARTY-MONEY-018 | Withdrawn: a new campaign was said to reach the first playable mission with `Player+0x38 == 0`; the participant factory writes 100 first (`PARTY-MONEY-024`). | — | ✖ retracted | [EXP-0154](../experiments/EXP-0154-campaign-start/), [EXP-0156](../experiments/EXP-0156-start-money/) |
| PARTY-MONEY-024 | A fresh campaign starts with 100 in the human participant's `Player+0x38`; zero is only the constructor intermediate. | High / Medium | ✔ promoted | [EXP-0156](../experiments/EXP-0156-start-money/) |

### PARTY-MONEY-015

- The constructor fact stands: `Player::Player` stores 0 at `004fab27`.
- The registry, type-5, mission-10 type-8 and instant-23 censuses remain
  correct in their stated scopes.
- Participant factory `FUN_004d35e6` directly tests `Player+0x38` and stores
  `0x64` when it is zero (`004d36cb`…`004d36d4`). Its one direct caller is the
  join routine `FUN_004d27dd`, which then calls `FUN_004faff7` with delta 0 at
  `004d2fad`.
- The original call-site instrument could not see the direct store it was
  being used to exclude.

**Confidence.** High retained for the constructor and the four scoped data
negatives.

**Amended.** EXP-0156 refutes the complete no-credit conclusion
(`PARTY-MONEY-024`, [`retracted.md`](retracted.md)); it had been graded
Medium. The withdrawn headline read: "The purse starts at zero and nothing at
or before the campaign's first mission credits it." Two lawful
original-runtime saves decode to 100, including an automatic save at sub-tick
1.

### PARTY-MONEY-016

- `Player::Serialize` loads the dword at `0051111f`, passes it through the XOR
  involution `0x5c073f4d` and writes one u32 (`SAV-OBF-029`). This record is in
  the campaign half and therefore in both save shapes.
- Independently, the mission carry copies sixteen bytes from
  `{Player+0x38,+0x48,+0x4c,+0x54}` into the packet and character file, and its
  importer writes those four dwords back to the same offsets
  (`PARTY-CARRY-005`).
- The purse is therefore per participant, not a field on the hero or one
  campaign-global balance.
- Arithmetic readers decide signed interpretation: the shop affordability gate
  uses a signed `JGE`, while instant 23's bare `ADD` has no clamp and wraps
  modulo 2^32 (`TRIG-MONEY-028`).
- G2: widening the purse changes the `Player` memory layout, save member width,
  sixteen-byte carry block and packet ABI. It does not require changing a
  shipped map or registry file.

**Confidence.** High for object, offset, width, save and carry: two
independently read paths name the same field and width. Medium for signed
interpretation as a general law: individual operations are read, but no type
declaration survives in the image.

### PARTY-ADDHERO-017

- `FUN_004b4880` calls `FUN_0048a6e0` on the current mission record. That
  routine sends one command `0x49` per `u16` and clears the array.
- Command `0x49` calls `FUN_004d8fcd`, which assigns the existing player as
  owner, appends the actor to `Player+0x20`, creates and appends a group to
  `Player+0x24`, and inserts the actor into that group.
- It does not write `Player+0x34` and does not use the mercenary constructor or
  type field. The actor is therefore a companion player character, not the
  primary character and not a mercenary.
- The serialized group list preserves it in saves. Complete mission-to-mission
  memory persistence retains `PARTY-PERSIST-014`'s Medium limit.

**Confidence.** High for activation, insertion, classification and save
persistence: complete instruction paths, with the distinct primary and
mercenary alternatives excluded. Medium for uninterrupted mission-to-mission
memory persistence.

### PARTY-MONEY-018

- The retained facts: the constructor still writes a 32-bit zero, opcode `0x48`
  still carries no purse, and the scoped mission and data negatives stand.
- What failed is the jump from those facts to playable state: the participant
  factory writes 100 directly before hero creation, outside the enumerated
  delta-helper callers.
- The first-sub-tick automatic save and a later explicit save both decode to
  100.

**Confidence.** The retained subclaims keep their original grades. The start
value and the complete no-grant conclusion are retracted.

**Amended.** Retracted by EXP-0156: the claim's own prediction failed
(`PARTY-MONEY-024`, [`retracted.md`](retracted.md)). The withdrawn headline
read: "A new campaign reaches the first playable mission with
`Player+0x38 == 0`, independent of character-generation choices."

### PARTY-MONEY-024

- `Player::Player` stores 0 at `004fab27`.
- The participant factory `FUN_004d35e6`, whose only direct caller is join
  routine `FUN_004d27dd`, then executes `CMP [Player+0x38],0` / `JNZ` /
  `MOV [Player+0x38],0x64` at `004d36cb`…`004d36d4`. It therefore preserves an
  existing nonzero purse and gives 100 to a fresh one.
- Join next calls `FUN_004faff7(player,0,1)` at `004d2fad`; its store at
  `004fb00a` adds zero, leaving 100, and emits opcode `0x67`. The client copies
  packet `+0x0a` into its local participant balance at `004169c5`.
- Chargen opcode `0x48` contains stats and appearance/class/sex but no purse,
  and resolves this already existing `Player`. Thus no purse write occurs after
  chargen confirmation and before play; the playable state inherits the join
  value.
- Two lawful original-runtime mission-10 saves independently decode the human
  Fergard record through `stored XOR 0x5c073f4d`: 100 at sub-tick 1/full tick 0
  and 100 at sub-tick 129/full tick 8.
- Last-writer vocabulary: `004d36d4` is the last value-changing start write;
  `004fb00a` is the later physical store and UI notification with delta zero.

**Confidence.** High for value, owner, ordering, save representation and
notification: every hop is a named instruction, `callto:4d35e6` is 1 hit / 1
owner / 0 orphan, and two different-tick original-runtime saves discriminate the
constructor rival. Medium that the experiment's typed writer list is image-wide
complete: `disp:38` is 889 accesses / 367 owners / 2 orphan, and whole-object
copies remain a blind spot. The start chain does not depend on that universal.

## Joins and the mission-end boundary

| ID | Claim | Confidence | Status | Evidence |
|---|---|---|---|---|
| PARTY-JOIN-025 | A mid-mission join is one routine, `FUN_004d1e14`, reached from trigger instants 19 and 22: `PARTY-INSTALL-012`'s install sequence without the hero pointer or a fresh runtime id. | High | ● active | [EXP-0165](../experiments/EXP-0165-join-persistence/) |
| PARTY-ENDCULL-026 | The server's end-of-mission cull `FUN_004d05b2` decides which players and actors exist on the next map; it keeps an actor iff `0x21 <= word[actor+0x0e] < 0x40`. | High | ● active | [EXP-0165](../experiments/EXP-0165-join-persistence/) |
| PARTY-BAND-027 | Withdrawn: the server survival band `[0x21,0x40)` was said to be the Humans creation arm's typeID range; the band is real, but that arm's output is conditional (`PARTY-M20-031`). | High | ✖ retracted | [EXP-0192](../experiments/EXP-0192-mission20-party-boundary/) |
| PARTY-PERSIST-028 | The campaign's mission-to-mission edge preserves the surviving human `Player`, its name and every actor that first passes the client and server filters, not an unfiltered roster. | High | ● active (amended) | [EXP-0165](../experiments/EXP-0165-join-persistence/), [EXP-0192](../experiments/EXP-0192-mission20-party-boundary/) |
| PARTY-JOINCORPUS-029 | The shipped join corpus holds 28 runtime instant-19/22 nodes on 16 maps, 22 of which hand actors to player 1; the handed Humans are not all on the kept side. | High / Medium | ● active (amended) | [EXP-0165](../experiments/EXP-0165-join-persistence/), [EXP-0192](../experiments/EXP-0192-mission20-party-boundary/) |

### PARTY-JOIN-025

- Trigger instants 19 (`Change unit's owner`) and 22 (`Change group's owner`)
  both dispatch to `FUN_004d1e14`: arm `00539ed1` once, and arm `00539ee8` once
  per group member (`TRIG-ACT-004`'s table).
- In order, the routine:
  1. removes the actor from its current group if it has one (`004d1e4e`,
     `FUN_0050f996`);
  2. removes it from the old owner's flat index (`004d1e60`, `FUN_005194c0` on
     `[old+0x20]+4`);
  3. sets `actor+0x14 = newPlayer` (`004d1e6b`);
  4. appends to `[new+0x20]` (`004d1e78`, `FUN_0050fbee`);
  5. allocates `new(0x48)` + `FUN_0050f81b` (`004d1e7f`, `004d1e9a`);
  6. appends that group to `[new+0x24]` (`004d1ec8`, `FUN_0051a050`);
  7. inserts the actor into it (`004d1ed4`, `FUN_0050f950`).
- It does not write `player+0x34` and does not allocate a fresh runtime id,
  which is what separates it from the `0xbe` carry arm and from chargen.
- It then clears both the old and the new owner's `+0x2c` mask from
  `actor+0x18` (`004d1efd`, `004d1f1c`) and broadcasts through `FUN_004e7de3`.
- A joined actor arrives in a group of its own, as a carried character does.

**Confidence.** High. One routine read end to end; every field is named by its
own instruction, and the two absences are two of the five writes
`PARTY-INSTALL-012` lists for the other entry points.

### PARTY-ENDCULL-026

- Its call site is the campaign state machine's mission-end arm, `00473e94`,
  nineteen instructions after `PARTY-CULL-004`'s client cull.
- Three passes:
  1. per player with `+0x28 == 0` and a hero, drop the hero from the map's cell
     structure;
  2. empty the tick list `[0x00609558]+4` (`MOVE-TICK-009`) and expire every
     effect hanging off each of the player's actors;
  3. the roster cull.
- The roster cull destroys a `Player`, removes it from `[0x00609544]` and calls
  `vt+0x04(1)` through `FUN_00518100` when `+0x28 != 0`, or when the `CString`
  at `player+0x18` compares byte-equal to the literal `"Self"` at `0x005c5e88`
  (`004d07ec`...`004d0816`; the compare is `FUN_004c7ed0` to `FUN_00447fe0` to
  the CRT byte loop at `0x00553d80`).
- A surviving player has its placement latch cleared (`byte+0x3d = 0`,
  `004d0823`, the field `PARTY-GATE-013` names), and then its flat index `+0x20`
  is walked.
- An actor is kept iff `0x21 <= word[actor+0x0e] < 0x40` (`004d086f`,
  `004d087d`). Otherwise it is destroyed through `FUN_004fafaa`, which unlinks it
  from its group, drops that group from `player+0x24` if it is left empty and
  its owner is human, removes it from `+0x20` and calls `vt+0x04(1)`.
- A kept actor is reset in place: `word+0x94` from `+0x96`, `word+0x9a` from
  `+0x9c`, `byte+0x13c = 0`, `FUN_004f4b39`, then `+0x5c`, `+0x64`, `+0x44`,
  `+0x68`, `+0x40` zeroed (`004d0890`...`004d08f5`).
- Neither routine writes an inventory field, not the sack `+0x7c`
  (`TRIG-ADDITEM-027`) and not the twelve worn slots `+0x198`
  (`SAV-CARRY-050`), so a survivor keeps what it carried.
- This is the server counterpart `PARTY-FLAG-003` recorded as never looked for,
  and it is not a flag word: the server asks the actor's typeID.

**Confidence.** High. One routine was read end to end with its five helpers;
each test and each store is a named instruction, both removals are calls to the
object's own `vt+0x04(1)`, and the literal was read out of the PE through its
section table.

### PARTY-BAND-027

- `FUN_004f9065` first streams `Humans` slot 16 into `actor+0x0e`, and
  overwrites it with `gender+0x21` or `+0x23` only when its constructor-mode
  argument is non-zero.
- Mission 20 supplies zero on its definition-id arm and supplies the exact
  `npc.reg` `Hero` flag on its npc arm. Its four transferred Humans therefore
  retain `typeID` `0x17,0x0a,0x0a,0x0a` and all fail the band. See
  `PARTY-M20-031`.

**Confidence.** High for the refutation: the conditional constructor branch,
both spawner arguments, both culls and two original-save boundaries
discriminate it from the creation-arm model.

**Amended.** Retracted as a whole by EXP-0192 (`PARTY-M20-030` to
`PARTY-M20-032`, [`retracted.md`](retracted.md)). The withdrawn text read: "The
survival band `[0x21,0x40)` is the Human creation arm's own typeID range" and
"A script handover of a Human is permanent".

### PARTY-PERSIST-028

- `FUN_00473110` reaches the client cull and then `FUN_004d05b2` before town
  construction, and reaches no teardown.
- `FUN_004e2327` consequently reuses the lone surviving player and skips the
  map name assignment.
- Repetition of that reuse is unbounded, but membership at each boundary is
  conditional: mission 20's four transferred Humans are gone before reuse
  (`PARTY-M20-031`).

**Confidence.** High for the ordered boundary and reuse arms, read end to end.
Original saves independently witness the narrowed membership.

**Amended.** Narrowed by EXP-0192 (`PARTY-M20-031`);
[`retracted.md`](retracted.md) records the unqualified roster-persistence
reading as superseded. The withdrawn sentence read: "Persistence is therefore
unbounded across boundaries". Reuse is unbounded only for the Player and the
actors that survive both filters; it is not a bypass around them. The surviving
Player reuse arm and the name rule stand.

### PARTY-JOINCORPUS-029

- Mission 20 alone supplies four counterexamples to the previous partition of
  the handed Humans as the kept side: one npc-arm Human without the `Hero` flag
  and three definition-id-arm Humans, all with out-of-band table typeIDs.
- Mission 40's measured recruit still survives, for the narrower reason that
  `npc25` has the exact `Hero` flag, so its npc arm passes a non-zero
  constructor mode and performs the player-character overwrite.
- The published mission-40 saves and `AddHero` census otherwise stand.

**Confidence.** High for the corrected mission-20 and mission-40
classifications at instruction level. Medium for the corpus counts and the save
observation.

**Amended.** Narrowed by EXP-0192 (`PARTY-M20-030`, `PARTY-M20-031`);
[`retracted.md`](retracted.md) records the Humans-as-kept-side partition as
refuted: a creation arm does not decide the cull side. The node counts, the
mission-40 saves and the `AddHero` census stand.

## Mission 20's transferred actors

| ID | Claim | Confidence | Status | Evidence |
|---|---|---|---|---|
| PARTY-M20-030 | Mission 20 transfers exactly four live actors to player 1 when `T08` fires: group 16 units 136–139. | High / Medium | ● active | [EXP-0192](../experiments/EXP-0192-mission20-party-boundary/) |
| PARTY-M20-031 | A transferred Human survives a mission boundary by constructor mode and resulting typeID, not by being a Human. | High | ● active | [EXP-0192](../experiments/EXP-0192-mission20-party-boundary/) |
| PARTY-M20-032 | Mission 20's transferred actors have no town, tavern, later-mission or post-boundary save identity. | High / Medium | ● active | [EXP-0192](../experiments/EXP-0192-mission20-party-boundary/) |

### PARTY-M20-030

- Unit 136 is `npc52`, runtime name Sarindar, `Humans[201] M10_Merchant`,
  table `typeID=0x17`, with a Bronze Amulet and Rare None Robe.
- Units 137–139 are three unnamed instances of definition 114,
  `Humans[58] NPC14_1`, table `typeID=0x0a`, each with Iron Mace, Hard Leather
  Small Shield, Uncommon Hard Leather Helm, Bronze Cuirass, Uncommon Hard
  Leather Mail, Bronze Bracers, Leather Gauntlets and Leather Boots.
- EN and RU resolve the same rows, typeIDs and equipment.
- An original mission-20 save independently holds the same four under the human
  player with runtime ids 111–114, class keys 201/58, type words `0x17/0x0a`,
  2/6 worn slots and empty packs.

**Confidence.** High for identities, count and fields: the
trigger/group/placement/Data.bin join agrees on both roots, and an original
save observes a different representation. Medium for the translated display
name, observed in the EN save only.

### PARTY-M20-031

- `FUN_004f9065` streams Humans slot 16 to `actor+0x0e`; only a non-zero mode
  argument enters `004f95c0`'s overwrite to a player-character value.
- `FUN_004e26bb` passes zero from the definition-id and explicit typeID arms,
  and passes the exact `"Hero"` flag lookup from the npc arm.
- The server cull keeps only `[0x21,0x40)`. Independently, the client
  classifier gives mission 20's `0x17` and `0x0a` actors no keep bit, and the
  client cull requires that bit.
- Mission 20's sole reachable victory route calls the client cull, then the
  server cull, before town; its authored VIP condition is unreferenced. All four
  actors are therefore detached and destroyed at completion if still present.
- Brian in mission 40 survives because `npc25` is flagged `Hero`, not because
  it took the Humans arm.

**Confidence.** High. Two independent culls agree; the constructor mode and each
spawner arm are named instructions; mission-20 and mission-40 data discriminate
the competing creation-arm rule.

### PARTY-M20-032

- Closely paired mission save `game0009` contains all four under Danath. Town
  save `game0010` has matching Player and Danath allocator identities but none
  of the four; matching values do not prove that both files came from one
  process. Its second actor is fresh town `AddHero=22` companion Reniesta.
- Mid-mission save/load serializes the four through the Player group graph. A
  town save cannot recreate records already culled, and load remaps pointer
  identities between processes.
- Tavern stock is independent, even though the client classifier puts these
  actors in the bit-4 bucket the live tally reads: the simulation and client
  constructors clear the mercenary type, the broadcaster omits that field group
  for zero, and the tally matches only 1..15.
- Mission 20 lists locked type 1 and therefore has an empty shelf. Mission 30
  offers type 14, whose fresh level-1 roster object reuses the same `NPC14_1`
  template and equipment as the three removed guards. Matching dress is
  template reuse, not actor conversion.
- Later mission-31 and mission-40 saves contain only Danath and Reniesta under
  the human Player. Because they are from another session, they corroborate
  roster absence but establish no identity comparison.

**Confidence.** High for the static cull, campaign creation and shelf/template
paths. Medium for the paired-save boundary; for the zero-type stock dependency,
because no image-wide writer sweep is committed; for later-session absence; and
for the predicted fresh type-14 runtime identity, which was not viewed in the
tavern UI.

## Open questions

- The remaining bits of `+0x18c` (`PARTY-FLAG-003`), and the eleven sites that
  copy the whole word between drawables.
- What the three replaced sub-objects are (`PARTY-LOSS-006`: `+0x154`,
  `+0x158`, `+0x10`), and the two group references `+0x40`/`+0x44`
  (`PARTY-ROSTER-002`).
- A second campaign. Every figure about how many (five players, one carried
  root, fifteen mercenary types) is a fact about the one campaign this install
  ships. The structures are read from code and do not depend on it; the counts
  do.
- Whether the `0xbe` carry also fires on the campaign's mission-to-mission edge,
  and what builds the human `Player` before the campaign's first map load, the
  one boundary `FUN_004e2327`'s reuse arm cannot serve (`PARTY-PERSIST-014`,
  `PARTY-PERSIST-028`).
- The six writers of `player+0x20` classified only by their append site
  (`PARTY-WRITE-011`): `FUN_004d1e14`, `FUN_004d8fcd`, `FUN_004f164c`,
  `FUN_004feadb`, `FUN_00504da1`, and the second `FUN_004d3755` site. Each is
  one listing.
- `actor+0x13c`, the flag that decides whether `FUN_004d3755` repairs a hero or
  leaves it alone (`PARTY-INSTALL-012`), and `actor+0x70`'s writers. One writer
  of `+0x13c` is named: the mission-end cull sets it to 0 on every survivor
  (`PARTY-ENDCULL-026`), the state that makes the repair arm return a carried
  hero untouched. What sets it non-zero, other than corpse decay, is not read.
