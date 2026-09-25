# TAVERN — mercenary hire

Claims about the inn's hiring hall. The inn is a town building, not a file
format. Its other half, the mission-giver (`InnMission` / `InnNPC`, EXP-0060),
is covered by `REG-SCN-064`; the two halves share the view and the record and
nothing else. Format of this file: [registry.md](registry.md).

## Terms

The terms follow the engine's own indexing.

- A type is the 1-based subscript that `byte [unit + 0x15b]` (client) and
  `byte [unit + 0x14c]` (server) carry. It is also the `[npc<t>]` section number
  in `scenario.res::npc.reg` and the index into `[General] MercenaryCount`.
  There are fifteen types.
- A pool is how many units of one type exist.
- A shelf is what a given mission offers.

## Types, shelf, hire, price and level

| ID | Claim | Confidence | Status | Evidence |
|---|---|---|---|---|
| MERC-TYPE-001 | A mercenary is a `CUnit` carrying a 1..15 type id, and the type is the index everything else is keyed by. | High | ● active (amended) | [EXP-0062](../experiments/EXP-0062-mercenaries/) |
| MERC-SHELF-002 | The tavern offers `[Mission<n>] Mercenaries` ∩ a permanent unlock list ∩ a non-empty pool: three sets under three separate keys. | High / Medium / Unknown | ● active (partially retracted) | [EXP-0062](../experiments/EXP-0062-mercenaries/) |
| MERC-HIRE-003 | A hire is one flag per type, and it takes the whole squad of that type, not one man. | High | ● active | [EXP-0062](../experiments/EXP-0062-mercenaries/) |
| MERC-PRICE-004 | `cost(type, mission) = (PriceA + n · PriceB) × unitPrice(mission)`, and `PriceA`/`PriceB` are the `npc.reg` fields whose use `SHOP-NPC-012` left Unknown. | High / Medium | ● active | [EXP-0062](../experiments/EXP-0062-mercenaries/) |
| MERC-LEVEL-005 | A mercenary's level is a step function of the mission number alone, and it selects a different `Data.bin` template rather than scaling one. | High | ● active | [EXP-0062](../experiments/EXP-0062-mercenaries/) |
| MERC-CMD-007 | The tavern is a simulation class the image names, `Tavern`, driven by command `0x37` (rebuild the roster) and command `0x38` (spawn and pay for the hire). | High / Unknown | ● active | [EXP-0062](../experiments/EXP-0062-mercenaries/) |

### MERC-TYPE-001

- `FUN_00466940` builds the inn's list from the campaign document's
  `CMapPtrToPtr` (`[campaignScreen+0xd0] + 0x9bc`, size `+0x9c0`, count
  `+0x9c4`). It admits a value iff `CObject::IsKindOf` returns true for the
  `CRuntimeClass` at `0x599208`, whose fields decode as name `0x5bc98c` =
  `"CUnit"`, object size `0x1b0`, schema `0xffff`, no `CreateObject`.
- The type id is `byte [unit + 0x15b]`. Three routines use it as a subscript:
  `word [record+0x60 + (t−1)·2]` at `0047f456`, the same at `00466a74`, and
  `dword [record+0x88 + (t−1)·4]` at `0048a12e`.
- On the server the same id is `byte [obj + 0x14c]`, written at `005059a5` from
  the spawn loop's own counter.
- Fifteen is the loop bound of the roster builder `FUN_005053ce`: `t = 1..2`
  builds `Unit` objects from the literals `"Catapult"` and `"Ballista"`
  (`0x5c331c`, `0x5c3320`), and `t = 3..15` builds `Human` objects from
  `"NPC%02d_%d"`. It matches the 15 elements `[General] MercenaryCount` ships.
- The inn draws each type with `graphics\interface\inn\Unit%d\sprites.16a` keyed
  by the same byte (`0047f889`).

**Confidence.** High. The class is read off its own `CRuntimeClass` record, and
the three subscript computations are named instructions on the repaired table.
`refto:599208` returns 13 hits in 12 owners, 0 in orphan or undisassembled code.

**Amended.** `TAVERN-TALKPIC-016` (EXP-0400) corrects this claim's attribution
of `0047f94b`: that address is the talk loop's non-Hero branch.

### MERC-SHELF-002

- `FUN_00480560`, the inn view's `vt+0x80`, walks `Mercenaries`
  (`campaignScreen+0x5e4`/`+0x5e8` = `record+0x9c`/`+0xa0`). It is the only
  reader of `Mercenaries` on an image-wide `EnumRefs disp:`. It keeps an element
  only if `FUN_0048a190` finds it in the array at `record+0xac`
  (`m_pData +0xb0`, `m_nSize +0xb4`). `FUN_00466940` then drops any type whose
  pool is 0.
- Only `FUN_0048a1e0`, an append, fills the `+0xac` list. Its two callers are
  vtable slot 0 of the two mission-record classes: `FUN_004878b0` (main record,
  vptr `0x59a470`) and `FUN_00487100` (sub-mission record, vptr `0x59a460`).
  Each drains its own record's `[Mission<n>] EnableMercenary` (`record+0x30`,
  `REG-SCN-059`'s destination) into the main record's list, then
  `SetSize(0,−1)`s it. An unlock is consumed once and never cleared.
- The chooser `FUN_004878f0` has exactly one caller, `00473f08`, inside the
  mission-end arm of `FUN_00473110`. It picks the record whose mission number
  equals `record+0x118`, the accepted mission. `EnableMercenary` is therefore a
  completion reward.
- Corpus, all three roots (`evidence/offers.csv`): 13 values over 5 sections.
  `[Mission10]` unlocks nine types; side missions 41, 71, 111 and 121 unlock one
  each.
- Four types are named by a `Mercenaries` key before they are unlocked: type 10
  named at `Mission40` and unlocked by side mission 41, type 8 at 70 / 71, type
  1 at 10 (and 20, and 110) / 111, type 5 at 120 / 121.
- Being named is not being shown: `FUN_00480560` applies the filter before it
  builds the list, so the tavern draws nothing for a type not yet unlocked. On
  the main-mission ladder in ascending order, the first shelf that can show each
  of the four is 50, 80, 120, 130 respectively (`evidence/checks.txt` C8, all
  three roots). A consumer keeps "named by this mission" and "hireable in this
  mission's town" apart.

**Confidence.** High for the three sets: every step is a named instruction.
`callto:48a1e0` returns 2 hits and `callto:4878f0` 1, both on the repaired table
with `.rdata` slots included. Both writers are reachable only through a vtable,
so an `imm:`-shaped sweep would report zero callers for each. Medium that no
fourth set narrows the shelf: the reader enumeration is `disp:` over the array
object and its data and size fields, image-wide, with 0 orphan hits, but it
cannot see a pointer laundered out of the record. Unknown when within a chapter
a side mission's unlock starts to bite.

**Unknown.** The first-shelf figures walk the main-mission ladder in ascending
order, a reading of the data and not of a routine. A side mission completed
during its own chapter's town visit would apply its `EnableMercenary` at
`00473f08` with the record still on that chapter. Whether the record advances at
that point turns on `FUN_00488970`, which `REG-SCN-065` describes as deleting
the finished sub-mission and `SHOP-TOWN-022` glosses as advancing the main
mission. One measurement narrows this without settling it. The other mission-end
arm, `00473f20`'s `FUN_00488970` + `FUN_0048a360`, is dead: its guard reads
`[0x005be138]`, a global with 2 references image-wide, both reads, both in
`FUN_00473110`, and no writer. Its static value is 10 against a
`CMP EAX,0xa ; JLE`, which therefore always jumps. Every mission end takes
`FUN_004884d0` at `00473f35`. This leaves `SHOP-TOWN-022`'s reading of the
reachable site at `00473fb5` untouched.

**Amended.** The clause "Four types appear on a shelf one main mission before
the side mission that unlocks them, so the intersection is observable in play"
is retracted ([`retracted.md`](retracted.md), EXP-0062). The first missions at
which the four can appear are 50 / 80 / 120 / 130, not 40 / 70 / 10 / 120.

### MERC-HIRE-003

- `FUN_0048a110` (hire) and `FUN_0048a150` (dismiss) are mirror routines over
  `dword [record+0x88 + (t−1)·4]`. Each is gated on `FUN_0048a190` and has
  exactly one caller (`FUN_00480e80`, `FUN_00480f60`).
- The flag's only other writers are the record reset `FUN_00487640`
  (zero-fill 15) and `FUN_004885a0` (zero-fill 15, `MERC-DEATH-006`).
- Nothing is committed inside the inn. Leaving it runs `FUN_00480ad0`, which
  calls `FUN_0041d4ad`. That routine culls every non-hero `CUnit` from the
  document map and sends command `0x38` carrying `FUN_00488500`'s vector
  `out[i] = working[i] × hired[i]` (`IMUL DX, word [hired + i·4]` at
  `0048852b`): the entire surviving pool of each hired type.
- The server's `FUN_005056f1` then runs `for (i = 1; i <= n; i++)`, creating one
  object per head.
- The client's affordability gate is `FUN_00480d70`. `record+0x138` is the local
  player's money (`[[view+0x68]+0x9b4]+0x0c`) minus the price of every type
  already hired, and `FUN_00480e80` refuses a hire whose price exceeds it. The
  check is local only. The money moves once, on the server, by
  `FUN_004faff7(−Tavern+0x9c, 0)` at `005059ea`.

**Confidence.** High. Both flag routines, the vector's `IMUL`, the spawn loop's
bound and the single debit are read at instruction level. `callto:` was run on
all five routines, `.rdata` slots included.

### MERC-PRICE-004

- Server (`FUN_005056f1`, `005057d7`…`00505831`): `PriceA` =
  `FUN_005218b0(FUN_004c86a0(0x5f0758, t))` = `[npc + 0x18]` and `PriceB` =
  `FUN_005218d0(…)` = `[npc + 0x1c]`, the two fields `FUN_0048c510` writes
  (`SHOP-NPC-012`). `n` is the hired count for that type, and the third factor
  is `FUN_005050ee`.
- Client preview (`FUN_0047f410`, `0047f494`…`0047f4b0`):
  `(byte[unit+0x155] · avail + byte[unit+0x156]) × (i16)[unit+0x158]`.
- The mapping between the two is witnessed on both sides of the wire, not
  inferred from the arithmetic. `FUN_004e7de3` emits `{type, PriceA, PriceB}` in
  that order under mask `0x40000000` (`004e851c`, `004e8539`, `004e8564`), and
  `FUN_004104e8` receives them into `{+0x15b, +0x156, +0x155}` (`00411dc7`,
  `00411dfd`, `00411e33`). So `+0x156` is `PriceA` and `+0x155` is `PriceB`.
- `+0x158` is not sent: the client runs the ladder itself at `00411e80`.
- `unitPrice` is a jump table at `0x5051b1` on `mission/10 − 3`, arms
  `10 15 20 40 60 80 100 600 800 1000 6000 8000 10000`, everything outside
  `0..12` → 0. Its first argument, the type id, is pushed and read by no
  instruction, so the factor is per chapter and not per mercenary.
- Corpus, all three roots identical (`evidence/types.csv`,
  `evidence/costs.csv`): `PriceA` 28..80 and `PriceB` 0..15, present on exactly
  the 13 types whose pool is nonzero. The five pool-1 types all carry
  `PriceB = 0`, a flat price where the per-head term cannot bite.

**Confidence.** High for both formulas, read at instruction level, and for the
wire mapping, witnessed by the sender's and the receiver's own field order.
`callto:5218b0` and `callto:5218d0` return 2 hits each on the repaired table,
and the `refto:599208`-style orphan check is clean. Medium that this is the only
price path: `FUN_005050ee` short-circuits to a flat 450 when
`[0x005eb5a4] == 2`.

**Unknown.** Which mode writes 2 to `[0x005eb5a4]`. The global has 5 writers,
all in the `0x470000` app/menu family, and none was read.

### MERC-LEVEL-005

- `FUN_005051e5(mission)` computes `k = mission/10 − 1`, rejects `k > 14`,
  indexes the 15-byte arm table at `0x505257`
  (`00 00 00 00 00 01 01 01 01 02 02 02 03 03 03`) and jumps through `0x505243`
  to four arms returning 1, 2, 3, 4.
- Missions 10–50 give level 1, 60–90 level 2, 100–120 level 3 and 130–150
  level 4.
- Both consumers build a class name from it: `FUN_005056f1` at
  `00505919`/`00505929` and `FUN_0041ce98` at `0041d445`. Each formats
  `"NPC%02d_%d"` (`0x5c81b4` / `0x5c81a8`) and hands it to `FUN_004f8e78`
  (`new Human(0x1e8)`).
- Types 1 and 2 take no level: they are `Unit(0x198)` objects named by the
  literals `"Catapult"` / `"Ballista"`.
- A level is therefore a different character record, not a levelled-up one.
- Census over `world.res:data/data.bin`, all three roots: all 33 reachable
  `NPC<t>_<lvl>` names exist, and so do the 19 unreachable ones plus `Catapult`
  and `Ballista`, 60/60. The reading requires this, because `FUN_005053ce`
  builds every one of the fifteen types unconditionally at whatever level the
  town is at.

**Confidence.** High. The ladder is transcribed from its own arm table and jump
table, and both consumers' `sprintf` sites are named. The `Data.bin` census is a
whole-string count over the length-prefixed names on 3/3 roots, and the reading
predicted it before the file was opened.

### MERC-CMD-007

- `Tavern` is `0x59c7d0`-vtabled, `0xa0` bytes, `: Building : Token`, in the
  `.data` `CRuntimeClass` table (`EXP-0057/evidence/runtime-classes.csv`). The
  single-player instance is the static `0x609a70`, the counterpart of
  `SHOP-ENTRY-003`'s `0x60a120`.
- Command `0x37` (`FUN_0041ce98`, opcode written at `0041d40b`, payload
  `dword cmd+0x0a` = `FUN_005051e5(mission)`) reaches
  `FUN_005053ce(player, level)`. That routine destroys the previous roster,
  refunds `Tavern+0x9c` in full, and rebuilds one object for every type 1..15 at
  that level.
- Command `0x38` (`FUN_0041d4ad`, `0041d907`; `word cmd+0x0e` = the mission
  number, `dword cmd+0x0a` = element count, payload from `0041d98e`'s `memcpy`)
  reaches `FUN_005056f1(player, vector)`. That routine reads element 0 as the
  mission, `RemoveAt(0,1)`, spawns `n` of each remaining type and debits the
  accumulated cost once.
- Both arms sit in `FUN_004d5dd8` at `004d7020` / `004d706e`, and both call the
  static instance.

**Confidence.** High. Both senders' opcode stores and payload writes, and both
dispatcher arms, are read at instruction level, and the payload layouts match
field for field between sender and handler. Unknown for a map-placed `Tavern`.

**Unknown.** Whether a map-placed `Tavern` is reachable in play. Only the static
instance was followed, and no analogue of `SHOP-ENTRY-016`'s `server+0x14c` test
was looked for on this class.

## The pool across missions

| ID | Claim | Confidence | Status | Evidence |
|---|---|---|---|---|
| MERC-DEATH-006 | At every mission end a hired type's pool becomes its live tally, an idle type regains one head up to `MercenaryCount`, and every hire flag is cleared. | High | ● active (amended) | [EXP-0062](../experiments/EXP-0062-mercenaries/) |
| MERC-POOL-011 | The shelf filter and the inn's price preview read one location, the working-pool `CWordArray` at campaign-record `+0x5c`; the server's own `n` is a different storage. | High / Unknown | ● active | [EXP-0395](../experiments/EXP-0395-mercenary-pool-persistence/) |
| MERC-POOL-012 | Nothing inside the tavern writes the pool; its four element writers and three storage routines are reached only from record construction, new-campaign reset, mission end, LOAD and record destruction. | High / Medium | ● active | [EXP-0395](../experiments/EXP-0395-mercenary-pool-persistence/) |

### MERC-DEATH-006

- Death and recovery are an explicit restore, not an absent write, and they are
  two different mechanisms.
- At the end of every mission `FUN_00473110` calls `FUN_004201c9` (`00473e75`),
  whose first act is `record->FUN_004885c0()` (`00420236`; `callto:4885c0` = 1
  hit).
- `FUN_004885c0` collects the local player's live units through `FUN_004667f0`
  (owner `[unit+0x14] == doc+0x9b4`, corpse stage `byte[unit+0x15a] == 0`) and
  tallies them per type into a 15-element scratch array (`00488679`…`0048869e`).
- It then merges per type (`004886a9`…`004886e7`):
  - a hired type gets `working[t] = the live tally`, so every man that died is
    gone from the pool for good;
  - a type not hired gets `working[t] += 1`, and only while
    `working[t] < MercenaryCount[t]` (`CMP DX, word[record+0x74 + i]` / `JNC`).
- Then `FUN_004885a0` zeroes all fifteen hire flags (`004886eb`, its only
  caller), so no hire survives a mission.
- The tally runs before `FUN_004201c9`'s own cull. The cull keeps only `CUnit`s
  with `+0x18c & 1` (player character; `PARTY-FLAG-003`, `PARTY-CULL-004`),
  `(i16)[+0xfc] > −10` and the player's owner id, and drops every mercenary from
  the map.
- A consumer must implement the consequence: a squad wiped on a mission takes as
  many further missions to rebuild as it lost men, and only while it is left at
  home.

**Confidence.** High. The merge is six instructions on one branch. Both in-play
writers of the pool's elements are enumerated, this routine and the record
reset, and the cap is a read of the pristine `MercenaryCount` array. `callto:`
on both routines returns 1 hit each, `.rdata` slots included.

**Amended.** EXP-0077 corrected the cull's gloss of `+0x18c & 1` from hero to
player character (`PARTY-FLAG-003`, `PARTY-CULL-004`;
[`retracted.md`](retracted.md)). `MERC-POOL-012` (EXP-0395) narrows the writer
enumeration: it is complete for the in-play paths and not for the document path.
The campaign reader `FUN_00489580` writes the same elements directly at
`004897a8` and resizes the array at `00489783`, so "both writers" is not the
whole writer population of the storage ([`retracted.md`](retracted.md)). The
merge, its arms, its cap and the flag clearing are unaffected.

### MERC-POOL-011

- `MERC-SHELF-002`'s collector `FUN_00466940` and `MERC-PRICE-004`'s client
  preview `FUN_0047f410` each reach the campaign object through
  `CALL 0x00573196` and its `vt+0x7c`. Each then subscripts the same `m_pData`
  at `campaignObject+0x5a8` = `record+0x60` by `type − 1`:
  - `00466a6e MOV EAX,dword ptr [EDX + 0x5a8]`, then
    `00466a74 CMP word ptr [EAX + ECX*0x2 + -0x2],0x0` with `JBE`, which drops a
    type whose unsigned word is zero;
  - `0047f450 MOV EDX,dword ptr [EAX + 0x5a8]`, then
    `0047f456 MOV BP,word ptr [EDX + EBX*0x2 + -0x2]` and
    `0047f45d MOV dword ptr [EDI],EBP`, storing exactly the value
    `0047f4a0 IMUL EDX,dword ptr [EDI]` multiplies by `PriceB`.
- `MERC-PRICE-004`'s client `n` is therefore the working pool. `FUN_0047f410`
  also reads the pristine array's `m_pData` at `+0x5bc` (`0047f462`, `0047f479`)
  into a second out-parameter: the inn holds both the current and the full
  headcount and prices on the first.
- Image-wide census over the campaign-object-relative form of all three record
  collections: `disp:5a4` 0 hits; `disp:5a8` exactly two, both reads, the two
  above; `disp:5ac` 0; `disp:5b8` 0; `disp:5bc` 1 read; `disp:5c0` 0. The
  bare-immediate form of all six is 0 instructions. The 739 `imm:` hits over the
  same six are `0x5a4xxx`-shaped constants and `0x5b8xxx`-shaped string
  addresses that the prefix-matching mode returns, inspected as such.
- The server's price operand is a different storage. `FUN_005056f1` runs
  `0050581e CALL 0x004adcf0` (`ElementAt` on its argument array),
  `00505825 MOV DX,word ptr [EAX]`, `00505828 IMUL EDX,dword ptr [EBP + -0x24]`:
  a stack `CWordArray` built from command `0x38`'s payload, whose elements are
  `working[i] × hired[i]` (`MERC-HIRE-003`). For a hired type the two agree
  numerically and are two locations.
- Width: `PARTY-MERC-007` carries the pool across a mission boundary as "fifteen
  integers" and `SAV-629` quotes it as a "fifteen-integer mercenary pool-count".
  It is fifteen u16 words at `+0x60`; the fifteen dwords at `+0x88` beside it
  are the hire flags. Both rows cite `MERC-DEATH-006`'s merge, so both mean this
  array; only the word width is corrected, and neither row's finding changes.
- Evidence files: `evidence/enum-abs-release.txt`,
  `evidence/d-pool-release.txt`, `evidence/corpus-searches.txt`. Also cites
  `REG-SCN-064`.

**Confidence.** High. A closed census, not plausibility, excludes two separate
storages for the two readers: the two `disp:5a8` hits are the image's only two,
and neither offset has any immediate form at all. It also excludes the server
reading the client's location: `0050581e` is `ElementAt` on a stack array, not a
campaign-object displacement.

**Unknown.** `MERC-PRICE-004`'s flat-450 mode, which this claim does not touch.

### MERC-POOL-012

- `FUN_0048a110` (hire) and `FUN_0048a150` (dismiss) were read whole. Each is
  under thirty instructions, is gated on `FUN_0048a190`, and touches exactly one
  location: `0048a13f MOV dword ptr [EAX],0x1` and
  `0048a17f MOV dword ptr [EAX],0x0`, with `EAX = [EDI + 0x88] + (t−1)·4`.
  Neither references `+0x5c`, `+0x60` or `+0x64`, so a hire sets a flag and
  leaves the headcount alone until mission end.
- The four element writers:
  - `004877cb MOV word ptr [EAX + EDI*0x1],DX` in the record reset
    `FUN_00487640`, copying `working[i] = pristine[i]` over fifteen iterations
    (`004877be`…`004877df`) after `004877b5` has read `[General] MercenaryCount`
    into the pristine array;
  - `004886c5 MOV word ptr [ECX + EAX*0x1],DX` and
    `004886db MOV word ptr [ECX],DX`, the hired and idle arms of the mission-end
    merge `FUN_004885c0` (`MERC-DEATH-006`);
  - `004897a8 CALL ESI` in the campaign reader `FUN_00489580`, a raw read of
    `2·count` bytes straight into `m_pData` (`004897a4 PUSH EAX` supplies the
    byte count).
- The three storage routines:
  - the empty-array constructor `0x005707ee` at `00487203`;
  - `SetSize` (`0x00570858`, which zero-fills a grown tail and frees and NULLs
    `m_pData` at size 0) at `00487688` with 15 and at `00489783` with the
    document's own count;
  - the teardown `0x00570821`, which the record destructor `FUN_004873a0` calls
    on each pool object at `004875a9` and `004875b6`. That routine is a
    destructor by its own SEH frame, unwind-state stores, element-free loops and
    base-vptr stores. `0x00570821`'s own body was not read, so its effect on the
    storage is a role name from its call position.
- The events reaching those seven:
  - the constructor `FUN_00487190` has one call site, `00471870` (`SAV-600`);
  - the reset has three, `00471102`, `00473ab1` (`SESS-START-034`'s new-campaign
    arm) and `004769e4`, whose own only caller is `FUN_00473110` at `00473a84`;
    all three are new-campaign initialisation;
  - the merge is reached only from `00420236`, at mission end;
  - the reader is reached from `004776fe` and `00477fe1`, the two campaign LOAD
    drivers (`SAV-CAMPTAIL-070`).
- `FUN_00488500`, which builds the wire vector on the way out of the inn, reads
  `[ESI + 0x60]` and `[ESI + 0x88]` and writes only its argument.
- Record-relative `+0x64` has exactly one object reference in the whole campaign
  family, the SAVE count read; the other eleven hits at that displacement inside
  `00487000..0048c000` are `[ESP + 0x64]` stack locals. So neither tavern reader
  bounds-checks the subscript against `m_nSize`.
- Evidence files: `evidence/enum-rel-release.txt`, `evidence/range-release.txt`,
  `evidence/d-pool-release.txt`, `evidence/d-drivers-release.txt`,
  `evidence/d-reset-callers-release.txt`, `evidence/callers-release.txt`,
  `evidence/callers2-release.txt`. Also cites `MERC-HIRE-003` and `REG-SCN-067`.

**Confidence.** High for the four element writes, the three storage operations
and for hiring not debiting the pool. Both flag routines are short and were read
whole, so "no write" is a reading of the routine rather than an absence
argument; it excludes the live alternative that a hire debits immediately and
the mission-end merge only restores. Medium that there is no eighth writer. The
campaign-object-relative census is closed. The record-relative census was
crossed against all 35 record-base owners, whose 47 `+0x5c`/`+0x60`/`+0x64`
instructions split into 20 `[ESP + ...]` stack locals and 27 object references
read in place. Every one resolves to `[EDI+0xd0]` or `[EDI+0xec]`, a field
compared against the constant `0x445`, a halved dimension in a centring
computation, or an indirect call through a slot.

**Unknown.** A routine that receives the record pointer as a plain argument
while appearing in neither census; 1 093 bytes of `00487000..0048c000` are
undisassembled.

## Roster cells, pictures, statistics and voice

| ID | Claim | Confidence | Status | Evidence |
|---|---|---|---|---|
| TAVERN-MERCVOICE-008 | The fourteen mercenary bios ship on all three roots, but 16 of 21 mercenary voice parts are EN/RU-only; where the pre-release root carries voice, its byte size matches EN's, not RU's. | High / Medium | ● active | [EXP-0393](../experiments/EXP-0393-inn-text-closure/) |
| TAVERN-ORDER-015 | The mercenary cells follow the walk of the campaign document's actor map, not any persisted list. | High / Medium | ● active | [EXP-0400](../experiments/EXP-0400-tavern-roster-presentation/) |
| TAVERN-TALKPIC-016 | A talk-only cell draws the sheet of the object `FUN_00422b29(InnNPC[j])` returns: `HeroMage`/`HeroFighter` when `+0x18c` bit 0 is set, else `Unit<+0x15b>`; the selected cell animates it. | High / Medium | ● active | [EXP-0400](../experiments/EXP-0400-tavern-roster-presentation/) |
| TAVERN-TALKSTATS-017 | Selecting a talk-only cell shows that object's statistics in the upper left panel only when its `+0x18c` bit `0x40` is clear, and a synthesised object is created with it set. | High / Medium | ● active | [EXP-0400](../experiments/EXP-0400-tavern-roster-presentation/) |

### TAVERN-MERCVOICE-008

- The fourteen bios, `text/inn/mercenary/npc*.txt`, ship on all three roots: 0
  asymmetric.
- Their voice does not: 16 of 21 `inn/mercenary/*.wav` parts are EN/RU-only.
- `npc09p3.wav` and `npc13p3.wav` ship on EN with no RU counterpart: a
  text/voice split inside the same bio, EN richer.
- A fourth, owner-supplied pre-release root carries none of the 21 mercenary
  voice parts except `npc04p1.wav`, `npc10p1.wav` and `npc35p1.wav`. For those
  three its byte size equals the EN file's exactly (190188/179332/630792 bytes
  respectively), while RU's own recording of the same three differs
  (174496/371180/634548).
- This is the hiring hall's own text/voice census, distinct from the
  mission-giver half this ledger excludes (`DLG-INNVOICE-030`).

**Confidence.** High for the presence and byte-size figures: direct per-root
reads over a named, closed set of 21+14 paths
(`evidence/dialogue-presence.csv`). Medium for reading the three size matches as
unlocalized or reused audio: consistent with a byte-identical payload, but only
size was compared, not a payload hash, so exact identity is not established.

### TAVERN-ORDER-015

- `FUN_00466940` visits the map's buckets from index 0 and each bucket's chain
  from its head (`004669b7`..`00466a0f`). It keeps a `CUnit` whose type is in
  the unlock-filtered `Mercenaries` list and whose working pool is non-zero
  (`00466a45`..`00466a7a`), and appends in walk order. The inn view copies the
  list unchanged (`00480716`..`00480748`).
- The bucket is `((key & 0xffff) >> 4) % size` (`0042705c`..`00427069`), and the
  client inserts at the bucket head (`004116a4`..`004116c4`). Units with
  consecutive ids inside one 16-id block therefore come out newest first.
- `Mercenaries`, the unlock list and the pools decide membership only.
- Owner screenshot of the original EN tavern for a mission-130 SAV
  (`Mercenaries` `14,6,10,13,4,8,7,3,1,9,2,12,5`): the cells read
  `10 9 8 7 6 5 4 3 2 1 14 13 12` from position 0. That is descending type in
  two runs, the pattern of stock units created in type order with the id block
  boundary between types 10 and 12.
- The owner's statistics panels tie the three 420,000 cells (types 3, 4, 5) to
  `NPC04_4`, `NPC03_4`, `NPC05_4` field for field. Under the seat's
  left-to-right reading this puts types 5, 4, 3 at positions 5, 6, 7, as the
  rule predicts.
- Evidence files: `evidence/listing.txt`, `evidence/model-scores.txt`,
  `evidence/roster-grid.csv`, `evidence/mage-templates.csv`. Cites
  `MERC-SHELF-002`, `MERC-LEVEL-005`, `SAV-928` (the same walk, for membership),
  `MERC-CMD-007`, `SAV-918`, `SAV-926` (command `0x37` at tavern entry rebuilds
  types 1..15 in fixed type order, so the ordering ids are assigned after any
  load).

**Confidence.** High that the order is the actor-map walk. Rivals excluded on
the same 14 cells: `Mercenaries` key order matches 2, ascending type 2,
descending type 1 (`evidence/model-scores.txt`). Medium for the id-block split:
the stock units' ids were not derived from code, so where a split falls is
fitted to one screenshot, and the mage tie rests on the seat's reading of the
owner's click order.

### TAVERN-TALKPIC-016

- Loader `0047f905`..`0047f95d`. When `+0x18c` bit 0 is set, bit 1 picks
  `HeroMage` over `HeroFighter`.
- The object is a live `CUnit` passing the npc section's `Flags` terms
  (`DLG-SPEAKER-023`) or, when none passes, the synthesised object with
  `+0x18c = 0x48` plus the token bits and `+0x15b` = the npc id (`00421fe5`,
  `00422016`).
- Paint advances a cell's frame `(frame+1) % frameCount` only while its position
  equals the selection and more than 125 ms have passed
  (`0047f232`..`0047f26f`); other cells hold their frame.
- A sheet that does not ship aborts: the `.16a` constructor formats
  `"FATAL ERROR: can't load "` + path, shows it and calls `abort`
  (`00428cb6`..`00428d11`, `0056f643`).
- Over the 22 `InnNPC` elements per root (EN = RU):
  - npc22 draws `HeroMage`, npc23 and npc25 `HeroFighter` on either arm;
  - on the synthesised arm npc2, 30, 32, 41, 52, 62 and 64 draw their shipped
    `Unit<id>`;
  - npc90 (Mission40) and npc59 (Mission70) have no `Unit<id>` sheet. They are
    the only Human candidates a stock mercenary template of the stage's level
    passes on sex, class and face (`NPC10_1`, `NPC08_2`), paired with side
    missions 41 and 71, which unlock types 10 and 8.
- The mission-130 SAV's candidate npc25 draws
  `graphics\interface\inn\HeroFighter\sprites.16a`, a shared sheet, not a
  per-character one.
- `0047f94b` is this loader's non-Hero talk branch; this corrects
  `MERC-TYPE-001`'s attribution of that address.
- Evidence files: `evidence/listing.txt`, `evidence/talk-candidates.csv`,
  `evidence/inn-sheets.csv`. Also cites `REG-NPC-088`, `DLG-SPEAKER-022`,
  `MERC-CMD-007`, `SAV-926` and `MERC-LEVEL-005` (the level-1 type-10 stock unit
  `NPC10_1` at Mission40).

**Confidence.** High for the rule and for the Hero candidates' sheets. Rival
excluded: a sheet keyed by the `InnNPC` value alone, which would give npc25 the
unshipped `Unit25` and abort. Medium that npc90 and npc59 resolve to a live
stock mercenary: `Platoon`, `+0x15a` and the walk order were not evaluated, and
no running original was observed at those stages.

### TAVERN-TALKSTATS-017

- `FUN_0047d4a0` passes `view+0xec[sel − view+0xc8]` to `FUN_0047d4f0`. That
  routine blits `LeftStats.bmp` and `LeftPicture.bmp`, and calls the shared
  character/unit panel `FUN_00460480` with rect `(12,0)-(172,238)` only when
  `+0x18c & 0x40` is clear (`0047d55b`..`0047d5ba`).
- Below the panel it draws the object's composed figure through a temporary
  `allods%d.$$$` file when `+0x18c & 0x11` is set (`0047d5c6`..`0047d6d0`), else
  its `graphics\infowindow\<InfoPicture>.bmp` (`0047d70b`..`0047d791`).
- The synthesiser stores `0x48` (`00421fe5`) and afterwards only ORs in token
  bits: `1` for `Hero`, `0x10` for `Human`, `4`, `2` and the hero-relative sex
  and class bits. None clears a bit.
- The client class setter keeps only bit `0x80` of the old value and sets `9`
  for heroes or `0x18` for human classes (`0045f886`..`0045f8df`).
- A synthesised candidate therefore shows no statistics. Below the panel, one
  with `Hero` or `Human` shows its composed face figure. One with neither
  (`0x48 & 0x11 = 0`) takes the `infowindow` arm with the class name the object
  holds, and the synthesiser sets `+0x20` only under a `Picture` token.
- The one such shipped record is npc2 (Mission110, `Flags` `Platoon`), identical
  on both roots.
- For the mission-130 SAV the owner saw statistics for npc25's cell, so that
  cell holds a live actor passing `Hero,Face,!Female,!Mage`. When several pass,
  the first in the actor-map walk is taken. The owner identifies Brian.
- Evidence file: `evidence/listing.txt`. Cites `TEXT-UI-034`, `DLG-SPEAKER-022`,
  `HERO-APPEAR-053` (names the same two literals in `FUN_0047d4f0`).

**Confidence.** High for the gate and the value at creation. Rival excluded:
statistics drawn from the npc record itself; no `npc.reg` field reaches
`FUN_00460480`, which takes the object. Medium that the synthesised object keeps
bit `0x40` for its lifetime. The nine immediate dword stores to `+0x18c` in the
linear listing are `0x48` (`00421fe5`), `0x29`/`0x2b` (`0047b30b`/`0047b31c`),
`0` (`0045b25b`, `0047c6f9`) and four `0` stores (`004eb0bc`, `004edffc`,
`004ee3c3`, `004ef0a9`) set aside as another class's field: their routines zero
`+0x190` beside it, call an import on `this+0x3b0`, or pass the old value to an
import before zeroing it. Register stores such as the `and al,0xf7` merge at
`0045f7d1` (clears bit 3 only) were not all traced, nor other virtual slots or
pointers to the object. Medium also that no other path sets bit `0x40` on a live
actor.

**Unknown.** npc2's picture, and whether a live actor answers its `Platoon`
term.
