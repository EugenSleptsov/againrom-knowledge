# SAV campaign record

[Format reference](format.md) · [Application state](application.md)

The application campaign at `+0x548` follows the `&YA1` store on the same
file. Writer `00489270` and reader `00489580` use the same counted sequence.
The scalar suffix is not the entire campaign; arrays and markers make its
size variable. — SAV-CAMPTAIL-070, SAV-CAMPPROG-071

## Wire programme

All counts in this record are **plain u32**, not MFC `Count`. Source offsets
below are relative to the campaign record, or the child/base record where
specified. Pointer/count pairs describe in-memory storage, not wire pointers.
— SAV-CAMPPROG-071

The base record is written once at the head and once per child:

| Order | Wire | Source | Meaning |
|---:|---|---|---|
| 1 | Six u32 | `+04,+08,+0c,+10,+14,+18` | Mission, MapObject, Payment, shop lower/upper bounds, announce latch |
| 2 | u32 count, count × u16 | Count `+24`, array `+20` | AddHero collection at `+1c` |
| 3 | u32 count, count × u16 | Count `+38`, array `+34` | EnableMercenary collection at `+30` |

After the main base record:

| Order | Wire | Count / payload source | Meaning |
|---:|---|---|---|
| 1 | u32 count; each child = base + u32 age | `+50 / +4c`, child age `+48` | Side-mission children |
| 2 | u32 count; count × u16; count × u16 | `+64 / +60` then `+74` | Working / pristine per-type mercenary counts |
| 3 | u32 count; count × u32 | `+8c / +88` | Per-type hire flags |
| 4 | u32 count; count × u16 | `+a0 / +9c` | Current main mission's Mercenaries shelf |
| 5 | u32 count; count × u16 | `+b4 / +b0` | Permanent mercenary unlocks |
| 6 | u32 count; count × u16 | `+c8 / +c4` | InnNPC |
| 7 | u32 count; count × u16 | `+dc / +d8` | InnMission |
| 8 | u32 count; count × u16 | `+104 / +100` | TCMission |
| 9 | u32 count; count × u16 | `+f0 / +ec` | ShopMission |
| 10 | u32 count; pairs of u32 value, u32 kind | `+14c / +148` | Documents: kind 1 text, kind 0 picture |
| 11 | Seven u32 | Table below, in listed order | Scalar suffix |
| 12 | u32 count; counted markers | `+160 / +15c` | Selected-mission marker cache |

The two mercenary count arrays share one serialized count. InnNPC/InnMission
are semantically paired but have separate wire counts. TCMission precedes
ShopMission on the wire despite its higher memory address. With every count
zero, the complete record is 104 bytes. — SAV-CAMPPROG-071,
SAV-CAMPAIGN-076, SAV-CAMPAIGN-077, SAV-CAMPAIGN-078, SAV-CAMPAIGN-079,
SAV-CAMPAIGN-080, SAV-CAMPAIGN-081, SAV-CAMPAIGN-082, SAV-CAMPAIGN-083,
SAV-CAMPAIGN-084, SAV-CAMPAIGN-085, SAV-CAMPAIGN-086

## Scalar suffix

| Order | Runtime source | Meaning / restoration |
|---:|---|---|
| 1 | `+118` | Selected mission |
| 2 | `+114` | Raw value; both application load paths increment it into server difficulty `+84` |
| 3 | `+110` | AutoGetMission |
| 4 | `+11c` | LastMission |
| 5 | `+120` | First-MapPoint flag, computed on SAVE |
| 6 | `+124` | Mission time |
| 7 | `+128` | Raw value with an unnamed floating-point reader; nonzero producer Unknown |

`+120=1` restores MapPoint zero. Otherwise selected mission `+118` resolves
through the external mission-to-map-object registry. This stores a relation
to campaign map data, not an X/Y coordinate. `View/X` and `View/Y` in YA1 are
separate application fields. — SAV-CAMPPOS-072, SAV-892

`+114/+128` are raw passthrough, unlike computed `+120`. Constructors/reset
zero them; the local campaign reader tail does not consume them again.
The `+114` singleton chain reaches the three-way ALM placement law. `+128` has
`00488c00` as a floating-point reader, but the producer of its nonzero saved
values remains Unknown. Neither field is spare storage. — SAV-598, SAV-599, SAV-600,
SAV-601, SAV-602, UNIT-GATE-012, UNIT-GATE-013

## Markers

```text
u32 value
u32 n             writer uses strlen(string)+1
n string bytes    includes terminating NUL
u32
u32
```

A marker occupies `17+strlen(string)` bytes. The seven scalars plus marker
count occupy 32 bytes only when marker count is zero. Both serializer arms
include the nonzero-marker grammar. The LOAD contract is narrowed to restoring
these bytes: picture-pointer reconstruction runs separately on world-map entry,
not synchronously in either SAV loader. — SAV-CAMPMARK-073, SAV-609

## Progress and collections

| Event | Located campaign effect |
|---|---|
| Ordinary main progression | Retains active main record; rejects a lower main request |
| Side-mission completion / age 2 | Removes its child record |
| Accept town candidate | Shortens its building array; keeps InnNPC/InnMission paired |
| Announce candidate | Sets the record's own announce latch |
| Complete mission | Drains EnableMercenary into permanent unlocks; clears hire flags |
| Activate town | Drains AddHero |
| Acquire document | Accumulates `(value,kind)` pair |

There is no completed-mission collection in this grammar. The marker cache is
world-map presentation state; surviving heroes are objects in the Player
roster, not campaign records. — SAV-CAMPAIGN-076, SAV-CAMPAIGN-077,
SAV-CAMPAIGN-078, SAV-CAMPAIGN-079, SAV-CAMPAIGN-080, SAV-CAMPAIGN-081,
SAV-CAMPAIGN-082, SAV-CAMPAIGN-083, SAV-CAMPAIGN-084, SAV-CAMPAIGN-085,
SAV-CAMPAIGN-086

## Mercenary state

Mercenaries `+9c` is the current main mission's shelf, not a cumulative union.
The tavern reader consumes the stored array directly; it does not recompute
membership from the mission number. Mission changes can replace rather than accumulate this shelf. — SAV-606

Hire sets its type's hire-flag dword. Mission entry does not clear that flag;
mission end does. Zero flags mean no outstanding hire, not that no hire ever
happened: named producers are hire, dismiss, reset and mission-end zeroing.
— SAV-615, SAV-616, MERC-HIRE-003

Actor class/type/equipment signatures do not encode hire origin. The same
Human shape can come from hiring, authored map placement or scripted transfer.
Creation-order IDs are session-local rather than a provenance tag.
— SAV-617, SAV-622, SAV-629

## Tavern eligibility and selection

The stored main-mission shelf is narrowed to permanently unlocked types before
the tavern builds its offers. The client then collects live CUnit objects of
those types with a nonzero unsigned working-pool word at `type-1`. These are
separate conditions; shelf membership
does not check the pool index's range. A type must have a corresponding
working-pool entry. — MERC-SHELF-002, SAV-928

| Input / consumer | Rule |
|---|---|
| Campaign shelf `+9c`, count `+a0` | Supplies candidate types |
| Permanent unlocks `+b0`, count `+b4` | Filters the shelf before client collection |
| Working pool `+60`, count `+64` | Collector tests the unsigned u16 at `type-1` |
| Client document map | Zero map count returns an empty result; nonzero count traverses buckets/next links even for an empty eligible set |
| Live map entries | Require a finite coherent map, acyclic stable links, valid objects and returning services |
| Mercenary selection on inn activation | Stores 0 for a nonempty result, -1 for an empty result |
| Caption / price refresh | Signed checks `selection <= count-1` / `selection < count` admit -1; neither is a lower-bound check |
| InnNPC | Separate append loop; it does not supply an empty-mercenary fallback selection |

The local empty-selection stores and indexed consumers are established.
Whether selection/storage survive intervening UI/resource calls unchanged
remains conditional. Do not infer a universal empty-save failure from the
negative-index condition. — SAV-928, SAV-929

Saved eligible types do not prove that live stock was published, selected and
kept valid through entry. Full original click order and the first failing
consumer remain Unknown. The [Player command prerequisite](player.md#tavern-command-prerequisite)
is a separate relation. — SAV-932

## Terminal campaign

Nonzero LastMission with selected mission divisible by ten reaches an earlier
mission-end branch to score production and message `0x428`. It locally
bypasses ordinary reward/advance and enters UI code with unresolved virtual
callbacks. This does not establish a final usable city or terminal save.
— SAV-890

The inner registry loader rejects a request above
`10*ScenarioMissionCount` before campaign stores/grants. Its outer higher-request
caller still writes current `+04` and selected `+118` and reports success.
Those numbers can therefore change with older registry/documents retained;
the terminal dispatcher may bypass the path entirely. An inferred increment
does not supply a successor campaign. — SAV-891

Campaign SAVE/LOAD has no separate terminal grammar or local LastMission
rejection. On the selected missing-map lookup, default object one becomes
zero-based index zero. Later callbacks and UI remain acceptance boundaries.
Native terminal-city LOAD remains Unknown. — SAV-892

## Extension boundary

Existing counts admit more entries of their existing shape. Added fields,
changed widths or reordered fields are outside the original fixed programme;
there is no version-selected campaign grammar arm. Preserve the programme or
version an extension outside it. Unnamed `+114/+128` cannot be repurposed as
padding. — SAV-CAMPPROG-071, SAV-599, SAV-600, SAV-602
