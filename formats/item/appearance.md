# Names and pictures

[Reference](format.md)

## Display names (`ITEM-DISPNAME-036`…`ITEM-NAMELIMIT-042`)

Nothing in the image composes an item name for display. The name is a **stored line**, selected
by a `CMap<u16, const char*>` built once at start-up from two `main.res` nodes
(`ITEM-DISPNAME-036`):

```
main\text\itemname.bin   N u16 keys, little-endian, no header   (shipped: 832 B = 416 keys)
main\text\itemname.txt   N lines, one display name each          (shipped: 416 lines)
map[key[i]] = line[i]     FUN_00468380, count = filesize(bin)/2
```

The key is the **packed item code** — the same `u16` an authored `.alm` element carries
(`ITEM-CODE-029`) — and its fields are fixed by the diagnostic formatter `FUN_00483c80`
(`ITEM-NAMEKEY-037`):

```
bits 12..15  material   index into Materials      (item+0x46)
bits  8..11  class      1 Weapons, 2 Shields, 3..13 Armors, 14 MagicItems
                        for an Armors row the value is also the row's own Slot column
bits  5.. 7  shape      index into Shapes                       (item+0x45)
bits  0.. 4  row        1-based index into the class's collection (item+0x0c)

class 14 only: bits 0..7 are the whole MagicItems index and there is no shape field
```

Reading a name is `FUN_00483e50`: take the interface element's `u16` at `+0x6`, `Lookup`. On a
**miss** the element is not drawn under a fallback — it is dropped from the list and a
diagnostic naming the code is emitted, so the name table is the inventory's admission list
(`ITEM-NAMEMISS-039`). The numeric form the diagnostic carries is seven digits with no
separator: `%02d%02d%1d%02d`, or `%02d%02d%03d` for class 14.

Consequences a consumer must not miss:

- **The name is free text, not a composition.** `Bronze`/`Amulet`/`Common` reads `Beard`; the
  `Elven Adamantium Amulet` key reads `Adamantium Amulet` in English with the tier word absent
  and carries a tier word in Russian. No grammar produces both (`ITEM-NAMEPOP-038`).
- **Localisation lives here and only here.** `Data.bin`'s own item row names are byte-identical
  on the two roots (`DAT-ITEMNAME-010`); `itemname.txt` differs on all 416 lines and
  `itemname.bin` not at all (`TEXT-ITEMNAME-019`).
- **Enchantment is invisible to the name.** The key has four fields and none is an effect, so an
  enchanted piece shows the same line as a plain one (`ITEM-BRACE-041`).
- **416 of 6 064 expressible items have a name.** The rest cannot be shown at all.
- The `[tier ][material ]shape` grammar (`UNIT-EQUIP-005`) is the **parse** direction only: it is
  how an authored string in a `Units.EquipItem` cell, a `Humans` equipment cell or a mission
  `.ini` becomes an item, and it is finished before anything is drawn (`ITEM-NAMEPARSE-040`).

## Appearance field `item+0x40` (`ITEM-APPEAR-023`, `ITEM-APPEAR-024`)

Every item carries one `u16` at `+0x40` that is the whole of its appearance. It is not stored: five
routines write it and all five write the result of `FUN_00525d60`, which is

```
item+0x40 = ((item+0x46 & 0xff) << 12)      the Shapes/Materials index -> material.reg
          | ((kind      & 0xff) <<  8)      the equipment slot, 1..12
          | ((item+0x45 & 0xff) <<  5)      3 bits, unnamed
          |  (item+0x0c & 0xff)             the Data.bin definition row index
```

`kind` is the literal `1` in the `Weapon` builder (which also stores it at `item+0x50`), `2` in
`FUN_0050cfa2`, `1` in `FUN_0050d5d4`, `byte[item+0x50]` in `FUN_0050c53a`, and the caller's
argument in the base `Item` builder. The `OR` is not a shift into disjoint fields, so a `+0x0c`
above 31 would corrupt the `+0x45` field.

`FUN_00483c80` turns the word into a seven-digit name — `"%02d%02d%1d%02d"` of the four fields, or
`"%02d%02d%03d"` of the top two and the whole low byte when the slot field is 14 — and that name
addresses both `graphics\inventory\<name>.16a`, the icon, and
`graphics\equipment\<figure>\<layer>\<name>.256`, the figure layer. `FUN_00483d70` tries to open
the icon and returns 1 when it cannot; the `0x76` equipment message then shows the engine's own
`"Invalid item weared "` for five seconds. *Weared* is the engine's word for the equipment slots.

**`FUN_00483c80` is not a method on `Item`.** Its `this` is a display record whose `+0x06` is a
copy of `item+0x40`, made by the record's only two constructor call sites:

```
0047c368  MOV DX,word ptr [ECX + 0x40]     0047c458  MOV AX,word ptr [EDX + 0x40]
0047c36f  CALL FUN_00484ac0                0047c465  CALL FUN_00484ac0
00484ac8  MOV word ptr [ESI + 6],AX        00484ad3  MOV byte ptr [ESI + 0xa],AL
```

So a consumer needs two fields, not one: the simulation's `item+0x40` and the panel record's own
copy. The record also carries `+0x04` u16, `+0x08` u8, `+0x09` u8, `+0x10` and `+0x14` dwords
(`FUN_00484990`, its copy constructor); what those are is unread.

**The picture word is rebuilt at construction, not carried.** `FUN_004dd02a` passes only three
values into each constructor — shape, material, row — and the kind nibble comes from the class:
the literal 1 for a `Weapon`, 2 for a `Shield`, the caller's argument for the base `Item`, and,
for an `Armor`, `byte[item+0x50]`, which the fill has just taken from the row's own `Slot`
column. So an authored code whose kind nibble disagrees with its `Armors` row's `Slot` allocates
by the nibble and **draws by the `Slot`**. Installed codes agree with their row Slots
(`ITEM-PICT-051`).

**A second self-check message.** Beside `"Invalid item weared "` sits `"Invalid item in
inventory "` (`0x5b83ac`), shown for **10 000 ms** at `004134e7`, gated on `FUN_00483e50` — which
is not a file open but a keyed lookup of the same word in the object at `0x5eb410`, and returns
**non-zero on success**, the opposite polarity to `FUN_00483d70`'s. `graphics\inventory\` is
composed at five sites in all: `00481f27`, `004838c0`, `00483db1` (the check), `00491815`,
`004a3571`. Inside `graphics.res` the nodes are under `inventory/`, not `graphics/inventory/` —
the composed path's first segment names the container (`ITEM-PICT-049`).

### The registry — which items have a picture (`ITEM-PICT-046`…`ITEM-PICT-052`)

The picture-name encoding covers the full u16 namespace. Installed item
icons use 416 words: 367 equipment pictures and 49 class-14 MagicItems.
Allowed shape/material combinations for group-C definitions are stored in
the ten raw bytes at runtime `entry+0x1c`: five u16 masks, one per Shapes
row, with one bit per Materials row.

```
allowed(row, shape, material) = (row.mask[shape] & (1 << material)) != 0
```

Every installed allowed combination has a picture. `Bone` is absent from
every row's masks and has no art. Adding a permitted combination requires
its mask bit and its matching graphics.res node. MagicItems are outside
this mask scheme. — ITEM-PICT-050, DAT-MATMASK-020

**Reading a name.** Absent words do not default alike: `FUN_004db6b5` returns `0` (`Common`) for
a missing `Shapes` word (`004db7ef XOR AL,AL`) and `FUN_004db801` returns `15` (`None`) for a
missing `Materials` word (`004db931 MOV AL,0x0f`). Both scan their collection downwards from the
last row. Use these distinct defaults when constructing a picture word. — ITEM-PICT-048

What the word is used for on the drawing side is [hero appearance](../hero/appearance.md).

Claims: `ITEM-APPEAR-023`, `ITEM-APPEAR-024`, `ITEM-PICT-046`…`ITEM-PICT-051`.
