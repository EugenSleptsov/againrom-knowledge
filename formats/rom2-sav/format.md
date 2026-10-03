<a id="rom2-single-player-save-bsg"></a>

# ROM2 single-player save (`Bsg&`)

A ROM2 `game*.sav` has the ROM1 [`Asg&` envelope](../sav/format.md) with a
different magic word. The simulation document shares its head, Player-list
position, world-present byte and trailer with ROM1; the world half, the
Player body and the physical tail differ. This page lists the differences
only; where it says "as ROM1", the [ROM1 SAV references](../sav/format.md)
apply. Field meanings and a decoded owner save remain Unknown. — R2-SESSION-017,
R2-SESSION-018

## Envelope

| File offset | Width | Field | Rule |
|---:|---:|---|---|
| `0x00` | 4 | magic | `42 73 67 26` = `Bsg&` (0x26677342); ROM1 is `Asg&` |
| `0x04` | 4 | `blobEnd` | Writer patches the first offset after the blob; the reader discards it, the slot list seeks to it |
| `0x08` | 4 | version | Writer emits `0x0BAD0002`; reader rejects values of `0x0BAD0001` and below |
| `0x0c` | 4 | `blobBytes` | Length of the blob |

The blob is `u32 outWords` and word-codec packets, as ROM1
[Encoding](../sav/encoding.md#word-codec). A 256-byte label follows the blob.
The slot list reads the label at `blobEnd`. — R2-SESSION-017

## Decoded document

```text
head                  u32, u32, CString, 11 x u32, mission u32, difficulty u32
Players               u32 metadata, then list of Player (below)
dead actors           list
worldPresent          u8
if worldPresent != 0:
    list              first world list
    list              second world list
    terrain           sparse key list, sub-object, u32 identity key
    session           10,374 raw bytes
    list              third world list
    list              ROM2-only object list
u32 0xBADFACE1, u32 0xBADFACE1
trailer state         400 bytes
```

The three lists before and after terrain and session sit where ROM1 has
Buildings, SpellEffects and Sacks; their classes were not read. The ROM2
reader reads both trailer words unconditionally and compares neither. ROM1 reads
the second only when the first equals 0xBADFACE1. The ROM2 world-half load builds the map from the named scenario before
reading terrain. — R2-SESSION-018

### Terrain and session

Terrain stores one `u32` per cell whose byte in plane two is above 0x0f, over
cell indices `0x807` up to `0xedee` (exclusive). The key is
`(index << 16) | (plane-two byte << 8) | plane-one byte`. The key list is
followed by one sub-object written through its own virtual slot and a `u32`
identity key. The ROM1 cell rows do not appear. List count form and the
sub-object layout are Unknown.

Session is 4,000, 1,000, 48, 400 and 4,908 raw bytes, then `u8`, `u8`, `u32` and
three `u32`: 10,374 bytes against ROM1's 4,374. Only the first and fifth
widths differ from ROM1. — R2-SESSION-019

## Player body

The scalar prefix has the same 51 bytes and type order as the ROM1 Player,
with different runtime offsets and the same XOR key on the two obfuscated
fields. After it, ROM2 writes 2,560 raw bytes, then the group list, then 36 raw
bytes. ROM1 has the group list, 32 raw bytes and a Diary; ROM2 has no Diary.
— R2-SESSION-020

## Physical tail

After the label ROM2 writes the [`&YA1` store](../rom2-reg/format.md) with the
ROM1 root names (`Character`, `GameOptions`, `SpellBook`, `Objects`,
`Inventory`, `Projectiles`, `Fog`). It then writes in order:

| Block | Layout |
|---|---|
| Pair list | `u32 count`, `count` pairs of `u32`, then two `u32` |
| Scenario record | 4,096 raw, 320 raw (4 by 4 cells of 20 bytes), `u32 count`, per node `u32`, `u32`, raw 16, then `u32` current index |
| Nine blocks | each `u16`, `u16`, `u32 length`, `length` raw bytes |

No ROM1 campaign record is found in this position; no ROM1 campaign routine was
compared byte for byte, so that the three blocks replace it is Medium. Field
meanings of the three blocks are Unknown. — R2-SESSION-021

## Related

The character file is a separate format: [ROM2 character file](../rom2-a2c/format.md).
