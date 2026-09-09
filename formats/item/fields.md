# Item fields and definitions

[Reference](format.md)

## Structure

Carried source-code 2 transfer to the city table subtracts signed16 unit weight
times unsigned16 detached count from the inventory load, then inserts into the
tray container. A partial quantity invokes the Item split virtual; insertion's
`vt+50` is Item stackability, not actor derive. The carrying actor has no
direct stat recompute in this dispatcher arm. This protocol result does not
close the ordinary UI partial-pickup gesture or notification callback effects
(`SAV-CITYMOVE-512`).

```
Item      0x50 : Token          vtable 0x59c8e0 (27 slots)   defs <- Data.bin "Magic Items"
  Armor   0x68 : Item           vtable 0x59c950              defs <- Data.bin "Armors"
  Shield  0x68 : Item           vtable 0x59c9a8              defs <- Data.bin "Shields"
  Weapon  0x84 : Item           vtable 0x59ca00              defs <- Data.bin "Weapons"
Sack      0x44 : Token          vtable 0x59ca88 (19 slots)
container 0x24 : CObList        no vtable of its own; one per actor and one per sack
```


## Item fields

The bounded per-class Item definition lookup does not extend unchanged to
actors. Unit uses its saved Units row without a bounds check; exact Humanoid
leaves a null definition; Human uses its saved Humans row only below type 33,
otherwise row 5. Its saved row byte is not rewritten. — ITEM-DEF-002,
SAV-ACTORBIND-544

Only the fields this specification pins. `+0x00..+0x3b` is the `Token` base, whose head
`FUN_00510e5c` serializes as `SAV-TOKEN-034`'s 37 bytes.

| Off | Type | Meaning |
|-----|------|---------|
| +0x08 | u32 | flag word; ORed together when two stacks merge |
| +0x0c | u8 | index of this item's row in its class's `Data.bin` collection |
| +0x14 | ptr | owning `Player` — stamped by the shop tray's destination-4 arm (`SHOP-TRAY-025`) |
| +0x1c | i32 | price (`SHOP-PRICE-011`); the `Token` "value" slot |
| +0x20 | — | ordered `CObList` of 0x48-byte Effects (`ITEM-EFFOBJ-072`) |
| +0x3c | ptr | the resolved definition row; re-derived on every load, never stored |
| +0x40 | u16 | packed item code; shelf spell id belongs to Effect+0x40, not this field (`ITEM-VALUE-115`) |
| +0x42 | u16 | **stack count**, default 1 |
| +0x44 | u8 | item kind; MagicItems prefix classification writes `3` Potion, `4` Scroll, `5` Book, otherwise 0. Literal 3 forces stackable |
| +0x45 | u8 | **Shapes row index** (the quality tier) — `FUN_004db6b5(name)`, the **5**-row table; entry 0 is `Common`, whose `@.damage` is 0.2 |
| +0x46 | u8 | **Materials row index** — `FUN_004db801(name)`, the **16**-row table; entry 0 is `Iron`. The `Armor` constructor bounds this one, not `+0x45`: `0050c3fe MOV DL,byte ptr [ECX + 0x46]` / `0050c401 CMP EDX,0xf` |
| +0x47 | u8 | serialized, unread |
| +0x48 | u16 | `ftol(shape.MagCap × material.MagCap)` — the only fill value taking no column |
| +0x4a | i16 | **per-unit weight**, default 1; `ftol(col3 × shape.weight × material.weight + 0.5)` |
| +0x50 | u8 | slot (`Armor`/`Shield`) or runtime column 11, else the literal 1 (`0050dc44`) |
| +0x52 | u16 | **to-hit bonus** — `Weapon::Equip` adds it to the actor's to-hit modifier `+0xe6` |
| +0x60,+0x61 | u8 | damage `(base, spread)`; see [item formulas](effects.md) |
| +0x6a | u16 | defence bonus (`#.deIrnce`) — added to the actor's `+0xfe` |

`IsStackable` = `vt+0x50` = `FUN_00508635`, inherited unchanged by all three subclasses:

```
stackable(item) = (item.kind == 3) || (item.effects.count == 0)
```

Splitting `qty` off a stack of `n` (`FUN_0050ebf9`):

```
if n <= qty:  the whole element leaves the container
else:         n -= (qty - 1); detached = item->vt+0x40(); detached.count = qty
```

The four detach-one virtuals invoke their concrete copy constructor. Copy walks the source effect
list in order and allocates a fresh Effect for every node. Weapon copy also reconstructs a fresh
owned Spell from the source Spell id. One exception matters to persistence: base Item copy omits
serialized byte `+0x47` without first running the default Item constructor, so a clone/split leaves
that byte indeterminate (`ITEM-EFFSPLIT-074`).


## Magic-item value

`FUN_00508486` initially copies an ordinary `MagicItems` row's parameter 0 to `item+0x1c` as a
signed `i32`, without a sign branch at assignment. Scroll and Book are name-derived exceptions:
Scroll sums spell-derived kind-41 values and falls back to parameter 0 only when the sum is zero,
while Book reads only the first Effect's spell value, or zero when empty (`ITEM-VALUE-115`). Mission
40's `Quest Item31` therefore leaves construction at 10,000 from its row, not from an enchantment.
The later descriptor path does treat `item+0x1c == -1` specially; the assignment rule is not a
global denial of sentinel semantics (`ITEM-MAGVAL-090`, `ITEM-WEAR-058`, `ALM-M40-074`).
