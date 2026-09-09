# Writing SAV files

[Format reference](format.md) · [Document](document.md) · [Application state](application.md)

A writer must satisfy the byte grammar and the relations between saved
objects. Opaque fields still need their actual state values; the grammar
does not supply a complete set of safe defaults. — SAV-WRITERAUDIT-380,
SAV-WORLDFRONT-432

## Emission sequence

1. Prepare Player/Group/actor and nested-object values. Preserve the intended
   alias graph and distinguish saved-address keys from gameplay IDs.
2. Emit the [document programme](document.md) in order: head, Player graph,
   dead list, shape byte, selected world half and trailer.
3. Include direct base/embedded serializers as well as tagged references.
   Each Player includes Groups, settings32 and Diary; each Unit includes all
   raw blocks and selected optional members.
4. For a world document, use the external ALM baseline and emit the saved
   terrain/cell overlays, terrain key, session4374 and counted world objects.
5. Emit `0xBADFACE1`, its global dword and trailer state400. Pad an odd logical
   endpoint by one byte, compress words and prefix `outWords`.
6. Write magic, placeholder `blobEnd`, version `0x0BAD0002`, `blobBytes` and
   blob. Patch `blobEnd` to the blob endpoint. Append label256, YA1 and the
   complete campaign record.
7. Validate physical/decoded intervals and the relations below. Loading the
   graph does not by itself validate the values consumed by the next action.

Archive indices are assigned by first use. Saved-address definitions need
unique nonzero keys; a reference uses its target's key and that field's
specific lookup timing/policy. The missing-key rule is narrowed by field:
Token's named resolver writes null, while stage-zero order repair retains the
raw key. — SAV-FRAME-021, SAV-CODEC-022, SAV-DOC-053, SAV-FULLREAD-252,
SAV-ARCHREL-253, SAV-PTRMAP-035, SAV-HUMRESUME-460

## Required relations

| Surface | Writer obligation |
|---|---|
| Envelope | `blobEnd=16+blobBytes`; blob includes u32 `outWords` |
| Word codec | Even decoded length, exact word count, complete opcode operands |
| CArchive | One shared class/object index; correct typed base; intended null/new/alias references |
| Player/Group | Intended owner/hero relations and head-to-tail actor order; early Group references need already-bound targets |
| Position | Matching terrain key for the named world rebind; coordinates/list position do not substitute for identity |
| Terrain | External baseline precedes overlays; cell keys correspond to block rows with static bit5 |
| Spellbook | n-1 reference operations for slots 1..n-1, each representing its intended object or null |
| Unit/Humanoid/Human | Six raw blocks, order list, presence-gated members, Humanoid XP and thirteen references |
| Shortcuts | Canonical kind-6 array holds four dwords /16 bytes |
| World application state | Producer emits Fog and Projectiles; empty Projectiles still includes FreeIndex and IDs |
| Campaign | Exact count/field order; shared count for the two mercenary-count arrays; MapPoint relation; NUL-inclusive marker strings |

A lookup can null a miss or retain its raw word, and some lookups occur only
after other objects load. Do not replace these rules with one generic fixup.
— SAV-GRPLOAD-560, SAV-GRPOWNER-561, SAV-TOKENLOAD-092, SAV-CELLLOAD-112,
SAV-UNITPROG-156, SAV-HUMAN-043, SAV-SPELLBK-041, SAV-ORIGFAULT-335,
SAV-PROJSTORE-428, SAV-CAMPPROG-071, SAV-CAMPPOS-072

## Values without a complete authoring contract

| Fields / relation | Known rule | Remaining Unknown |
|---|---|---|
| Human live values and modifiers | Restored literally; named later consumers may read before derive | Complete current-state producers, first consumer and safe residual values |
| Human constructor tails | Named constructors preserve `bc/bd` and clear `fc/fd` | Intervening writes and first-SAVE values |
| Player/Group/Diary raw fields | Defined widths, ownership and field-specific LOAD rules | Safe arbitrary defaults and complete city-consumer prerequisites |
| Tavern campaign lists | [Shelf, unlock and pool filters](campaign.md#tavern-eligibility-and-selection); empty collection selects -1 | Live stock publication, preserved selection/lifetime and full click order |
| Eleven world-head values | Constructor zero, literal restore, named conditional consumers | All-path first-SAVE values and unnamed consumers |
| Position/cell keys | Defined relations and class/path-specific repair | First use on paths without repair; safe unresolved-key placeholders |
| Session and trailer raw state | Exact transfer and named repairs | Unnamed field meanings and valid new-state values |
| Terminal campaign | Same counted grammar; separate terminal dispatch | A usable terminal city and its native consumer state |

Definition rows describe initial data; accepted current values are not
necessarily derivable from those rows alone. Constructor zero, a zero seen in
a save, or another save's accepted value is not a general default.
— SAV-HUMLOAD-445, SAV-HUMGAPS-449, SAV-HUMFIRST-465, SAV-HUMNEW-505,
SAV-900, SAV-LIVEPROD-413, SAV-WHEADLIMIT-525, SAV-TOKENLOAD-094,
SAV-CELLLOAD-113, SAV-790, SAV-791, SAV-892, SAV-928, SAV-929, SAV-932

## Compatibility and extensions

Structurally complete no-world files can load and later fail on the shop
route. The responsible value/consumer contract is not fully identified;
Player/Group/Diary and Human state cannot be declared independently sufficient.
Reducing the party can also leave an authored mission's hero reference
unresolved. Native next-action behavior remains a separate check from framing.
— SAV-ORIGVALUE-399, SAV-LIVEPGD-415, SAV-ORIGMISSION-400

The known descriptor programmes do not define runtime-built or unregistered
derived serializers. Nonminimal primitive encodings have known reader arms,
but do not establish complete malformed-file or allocation acceptance.
— SAV-READPOP-255, SAV-CONTSER-189, SAV-758

Keep original field widths/order and the six Human skill/XP slots. Added
fields require a versioned extension outside these fixed programmes.
Unrecognized YA1 entries can survive a raw registry write but have no general
fresh-application-SAVE preservation guarantee. Label residue, physical suffix
and decoded alignment are likewise not specified extension storage.
— SAV-HEROXP-063, SAV-CAMPPROG-071, SAV-915, REG-102, SAV-EXTSURV-239
