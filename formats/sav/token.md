# SAV Token and Position

[Format reference](format.md) · [Object programmes](objects.md)

## Token head

Every Token-derived body begins with these 37 bytes, serialized by
`00510e5c`. Player and Diary do not have this head.
— SAV-TOKEN-034, SAV-MEMBER-036

| Body offset | Type | Source | Meaning |
|---:|---|---|---|
| `0` | raw 12 | `*(this+10)` | Position object below |
| `12` | u32 | `+04` | Runtime creation-order ID |
| `16` | u8 | `+0c` | Registry/definition selector; class-specific use |
| `17` | u16 | `+0e` | Type word; Human LOAD compares it with `0x21` |
| `19` | u32 | `+08` | Low word is the map-unit ID on the mapped producer path |
| `23` | u16 | `+18` | Recipient publication mask in located send paths |
| `25` | u32 | `+1c` | Item value/price; other classes' meanings Unknown |
| `29` | u32 | `this` | Object's saved-address identity key |
| `33` | u32 | `+14` | Saved-address reference, actor owner on the actor path |

Both archive arms transfer the stated widths. Coincident runtime offsets in
another object are not the same field. Token `+18` is also Unit `+18` because
Unit calls Token without adjusting `this`. The mask is not a class discriminator. — SAV-OBJ-014, SAV-ID-015, SAV-TOKEN-034,
SAV-636, SAV-653, SAV-654, SAV-678, ITEM-VALUE-115,
SHOP-CONSUME-073

## Position

| Position offset | Type | Meaning |
|---:|---|---|
| `+00` | u8 | Cell X |
| `+01` | u8 | Cell Y |
| `+02` | u16 | Explicit packed cell `(cellY<<8)\|cellX` |
| `+04` | u8 | Sub-cell X; named constructors initialize `0x80` |
| `+05` | u8 | Sub-cell Y; named constructors initialize `0x80` |
| `+06` | u16 | Bytes left unmanaged by the named constructors/copies |
| `+08` | u32 | Caller-supplied terrain pointer/key |

Accessors return `fullX=256*cellX+subX`, `fullY=256*cellY+subY`.
The u16 spanning `+00/+01` and explicit `+02` encode the same cell in the
normal producer relation. The archive copies all 12 bytes, including bytes
that constructors/copies omit. LOAD overwrites the fresh current-terrain
pointer before terrain reconstruction. — SAV-TOKENPOS-074, SAV-TOKENPTR-075

No named constructor/copy/lifecycle path assigns a meaning to `+06/+07`.
Zero is a Medium-confidence authoring placeholder for those two bytes only;
a preserving reader retains them. Aliased/bulk copies and unwitnessed runtime
paths remain outside that conclusion. — SAV-TOKENPOS-074, SAV-TOKENLOAD-095

## Identity map

Token LOAD binds the saved key to the object just constructed in the map at
`[0x005cd758]+0x88`. The trailing reference uses `00527d30`: missing key
becomes null. Player's own key and Diary's trailing reference use the same
map at their named sites. A writer assigns a unique nonzero key per defined
object and references the target's key. Archive indices and map-unit IDs are
separate. The null-on-miss rule is narrowed to these named resolver sites;
Position, cell and order repair use the different policies below.
— SAV-TOKEN-034, SAV-PTRMAP-035, SAV-ARCHREL-253

Reference-repair policies differ by field:

| Field or family | Hit | Miss / zero |
|---|---|---|
| Token trailing `+14`, named Player/Diary/Group references | Replace with live pointer | Null at their named resolver |
| Position terrain `+08` | Replace with live terrain | Retain saved word |
| Ten cell-payload object fields | Replace with live pointer | Retain saved word |
| Stage-zero order keys and mover `+7c` | Replace with live pointer | Retain saved word |
| World-lifecycle Unit `+40/+64/+68` | Replace with live pointer | Null |

Use the exact call-site rule, not one generic fixup for all u32 fields.
Position's nonzero key must agree with the terrain record for the named
world-rebind route. Zero or a dummy pointer is not a general substitute.
An unresolved nonzero Position key remains unresolved; coordinates or record
order do not supply a replacement.
— SAV-IDCENSUS-254, SAV-TOKENLOAD-092, SAV-CELLLOAD-112, SAV-HUMRESUME-460,
SAV-GRPLOAD-560, SAV-847, SAV-908

## Reach of Position repair

World LOAD registers terrain, then dispatches live/dead actors and top-level
SpellEffects through their lifecycle hooks. Sacks load after the earlier
live-list construction and do not receive that pass. The complete loader has
no Building/loaded-Sack manager callback and establishes no nested Item/Effect
walk. City/no-world LOAD constructs no terrain and calls no such manager.
— SAV-TOKENLOAD-093, SAV-TOKENLOAD-094, SAV-CELLLOAD-112, SAV-DEADLOAD-125

Building, Outpost, Tavern, Shop and Sack serializers contain no direct Position
rebind/clamp call after their base body. Building/Sack tick bodies are empty.
Building placement/removal reads cell X/Y; Sack reads the packed cell, with
terrain supplied separately. Nested callbacks, later resave and actual first
Position use remain Unknown. — SAV-POSLOAD-140

## Session entry and later sends

The later session join walks Building and Sack managers. Both senders read
cell and sub-cell through full-coordinate getters. Sack sends full words;
Building shifts right eight and emits cell bytes. Neither sender uses packed
`+02`, unmanaged `+06` or terrain `+08`. Earlier Position use can precede
these mandatory-prefix consumers. — SAV-POSTLOAD-221

Before entry, nonzero campaign `+6b8` can call the ordinary wrapper
`004d2551`; nonzero server `+2c` admits the same sub-tick used by normal pacing.
It increments the counter, drains commands, dispatches world entries and the
global tick list. The zero campaign arm executes no tick. Queued/computed
callbacks can precede the later sender even though Building/Sack `vt+18`
bodies are no-ops. — SAV-FIRSTTICK-348, SAV-FIRSTTICK-349,
SAV-FIRSTTICK-350

Actor entry mask `-1` can admit equipment/container traversal, subject to
class, recipient ownership and type gates. An Item's Effects require nonzero
appearance `+40`; plain Item `+1c==-1` bypasses the common Effect walk.
Automatic Sack entry does not visit its container. Pickup drains quantity-one
Items, can deep-copy Effects on a split, deletes the source container, nulls
Sack `+40`, destroys the Sack and requests another gated actor send.
— SAV-POSTLOAD-222, SAV-POSTLOAD-223

## Recipient publication mask

| Operation | Token `+18` rule |
|---|---|
| `004f2885` | Return intersection with Player `+2c` |
| `004f28a8` | Return 1 if that intersection is zero |
| `004f28d0` | OR recipient bits into the field |
| `004f2639` reset | Zero `+18`, `+1c` and `+4` |
| Located revoke sites | AND with complement of recipient bits |

One revoke site immediately republishes. Other send paths use a dedup test,
a positive intersection test or inline OR; a null recipient can publish to
every bit. Zero therefore records current revoked/unpublished state, not
proof that publication never happened. Terminal stage does not determine the
mask; revoke and republish can change it during the actor lifetime.
Universal fog-of-war meaning and the complete class lifetime remain Unknown.
— SAV-678, SAV-653, SAV-664
