# SAV object programmes

[Format reference](format.md) · [Encoding](encoding.md)

Each row is an ordered body programme after its archive tag. `Base` means a
direct call to that base serializer; it adds no tag. Widths exclude the tag
and any newly introduced class descriptor. Offsets name the runtime source.
— SAV-MEMBER-036, SAV-CLASSSER-172, SAV-CONTSER-188

## Token-derived bodies

| Class | Ordered body after its base |
|---|---|
| `Token` | [37-byte Token head](token.md#token-head) |
| `VirtualCaster` | Token; u8 `+3c`; raw 6 from buffer pointed to by `+40` |
| `Unit` | [Unit programme](actors.md#unit-programme), beginning with Token |
| `Humanoid` | Unit; raw 24 from `+1cc`; `objref` at `+198+4*i`, `i=1..12`; `objref` at `+1e4` |
| `Human` | Humanoid; own store arm adds no bytes |
| `Effect` | Token; u8 `+3c`; u8 `+3d`; u32 `+40`; u8 `+0c` |
| `Effect_DirectDamage` | Effect; raw 24 from `+48` |
| `SpellEffect` | Token; u8 `+40`; u8 `+41` |
| `PointEffect` | SpellEffect; `objref<Effect>` at `+48`; raw u32 `+44` |
| `AreaEffect` | SpellEffect; u8 `+48,+49,+4a,+4b`; u16 `+4c`; `objref<Effect>` at `+44` |
| `SpellTransport` | SpellEffect; `objref<SpellEffect>` at `+44`; `objref<AreaEffect>` at `+48`; u16 `+4c` |
| `Item` | Token; `list32<Effect>` at `+20`; u16 `+40,+42`; u8 `+44,+45,+46`; u16 `+48,+4a`; u8 `+47` |
| `Shield` | Item; raw 22 from `+50` (its own store arm is empty; this transfer is inherited/shared) |
| `Armor` | Item; raw 22 from `+52`; u8 `+50` |
| `Weapon` | Item; raw 24 from `+52`; raw 22 from `+6a`; u8 `+50`; `objref` at `+80` |
| `Building` | Token; raw 22 from `+52`; u8 `+40`; u16 `+42,+44,+46`; u8 `+48,+60,+61`; u32 `+64,+68` |
| `Outpost` | Building; u32 `+84,+88,+80,+8c`; `Count(n)`; raw `8*n` from embedded `+6c` |
| `Tavern` | Building; u32 `+9c` |
| `Shop` | Building; u32 `+70`, the value cap |
| `Sack` | Token; u32 `+3c`; direct [container](items.md#container) at `*(+40)` |

These sequences are the schema-1 `.data` Token lineage. PointEffect LOAD
immediately remaps its trailing raw `+44` key: hit installs the pointer;
miss writes null. No extra bytes are consumed. Outpost's array has only the
`CObject` runtime descriptor; its element meanings remain Unknown.
— SAV-MEMBER-036, SAV-HUMAN-043, SAV-CLASSSER-172, SAV-CLASSSER-173,
SAV-CLASSSER-174, SAV-CLASSSER-175, SAV-CLASSSER-176, SHOP-SAVE-015

## Fixed extents and variable extents

| Class | Body bytes |
|---|---:|
| Spell | 9 |
| Token | 37 |
| SpellEffect | 39 |
| VirtualCaster, Effect | 44 |
| Effect_DirectDamage | 68 |
| Building | 77 |
| Tavern, Shop | 81 |
| PointEffect | At least 45 |
| AreaEffect | At least 47 |
| SpellTransport | At least 45 |
| Outpost | `95+8*n`, or `99+8*n` with a wide Count |

The minima for reference-bearing records use null references. A reference can
instead introduce a class and nested body. Item is 53 bytes **including** its
u32 Effect-list count, plus the encoded Effect reference entries:
`37 Token + 4 count + 12 scalar bytes`. Unit has a 603-byte fixed part and is
at least `609+L` with its name of length `L` and three null references;
Human adds at least 24+26 bytes.
Counts, strings, presence flags and nested references move the endpoint.
— SAV-MEMBER-036, SAV-UNITLEN-045, SAV-CLASSSER-173, SAV-CLASSSER-176

Building's Token map-unit ID joins the authored type-4 record. Its 77-byte body
is fixed. — SAV-BLDG-037

## Non-Token schema-1 bodies

| Class | Ordered body |
|---|---|
| `Player` | [Player programme](player.md#wire-fields) |
| `Diary` | `Count(n)`, raw `4*n`; `Count(m)`, raw `2*m`; u32 owner key |
| `Spell` | u8 `+08,+09,+0a`; u16 `+0c`; u32 `this` identity key |
| `Spellbook` | u32 `+18`; u32 `n` from array `+04`; `n-1` references for indices `1..n-1` |
| `TableLine` | CString at `+4`; direct CDWordArray at `+8` |
| `CMultiShopShelf` | Empty |
| `CMultiShopInstance` | Empty |
| `CMultiShopTemplate` | Empty |

Spellbook index zero is omitted in both directions. Each slot has its own
reference operation; new Spell objects and aliases follow the intended graph.
— SAV-DIARY-042, SAV-SPELL-044, SAV-SPELLBK-041, SAV-CONTSER-188

TableLine's grammar is for the descriptor-created base object. Other vtables
advertise the same descriptor and some use different serializers; those
unregistered derived programmes are not covered by this row. The three
CMultiShop slot-2 bodies are `RET 4`. Shop's own programme excludes stock.
— SAV-CONTSER-189, SAV-CONTSER-193, SHOP-SAVE-015

## Schema-0 bodies

| Class | Ordered body |
|---|---|
| `CStringArray` | `Count(n)`, `n` CStrings |
| `CDWordArray` | `Count(n)`, raw `4*n` |
| `CWordArray` | `Count(n)`, raw `2*n` |
| `CByteArray` | `Count(n)`, raw `n` |
| `CMapStringToString` | `Count(n)`, `n` pairs of CString key and CString value |
| `CMapStringToOb` | `Count(n)`, `n` pairs of CString key and `objref<CObject>` value |
| `CDib` | BMP file header[14], info/palette[`bfOffBits-14`], pixels[`imageBytes`] |

Both maps walk current bucket chains without sorting. CMapStringToOb values
participate in the shared archive index. Diary directly dispatches its two
arrays and therefore emits no class tag for them at those sites.
— SAV-CONTSER-189, SAV-CONTSER-190, SAV-CONTSER-191, SAV-CONTSER-192,
SAV-DIARY-042

CDib flushes the archive and writes through the underlying file. Store writes
`BM`, `bfOffBits=0x36+4*paletteCount`, `bfSize=bfOffBits+imageBytes`, a 40-byte
info header, palette and pixels. LOAD requires `BM`, obtains the middle size
from `bfOffBits`, and uses nonzero `biSizeImage` or the padded-row formula for
pixels. It does not validate `bfSize`. — SAV-CONTSER-192

## Embedded collections

| Source | Body | Meaning |
|---|---|---|
| Unit `+15c` | `Count(n)`, `n*u16` | Static packed-cell route |
| Unit `+178` | Same | Dynamic packed-cell route |
| `*(*(Unit+158)+90)` | Same | Actor order's own patrol path |
| Group `+20` | Same | Meaning Unknown |
| `*(*(Group+3c)+4c)` | Same, after raw AI80 | Group AI path |
| `*(Unit+140)` | Spellbook, presence-gated | Known spells |
| Diary `+04/+18` | CDWordArray / CWordArray | [Diary](diary.md) |

The unnamed word-list class uses vtable `0059c430`, constructor `00523040`
and serializer `00523100`. Each saved element is node `+8`; nodes themselves
are not saved. Runtime append grows the pool by ten nodes and teardown frees
it. LOAD uses that same append path. The Unit route elements are packed cells
`(y<<8)|x`; actor-order and Group paths have separate identities.
— SAV-EMBED-039, SAV-WLIST-040, MOVE-ROUTE-004, SAV-630, SAV-631, SAV-632, SAV-633, SAV-634,
SAV-GRPPATROL-570

## Programme scope

The tables cover the known creator-backed schema-1 game and schema-0 runtime
programmes. Schema `0xffff` descriptors with null creators are not
stream-creatable. Runtime-built, copied and aliased descriptors, plus an
unregistered derived object advertising TableLine's descriptor, remain outside
these programmes. — SAV-CLASS-033, SAV-SERPOP-047, SAV-READPOP-255,
SAV-CONTSER-189

A base programme can execute within a derived body without its own class tag.
An embedded programme can likewise execute without a descriptor. Absence of a
class name from a save does not make that programme unreachable.
— SAV-PRODDIRECT-207
