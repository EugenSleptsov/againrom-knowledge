# SAV simulation document

[Format reference](format.md) · [Encoding](encoding.md)

## Runtime source names

Document `world` is the server singleton `[0x005cd758]`. Terrain
`[0x005f22c8]` is a separate object. The names in the source columns below
do not identify a single shared allocation. — SAV-657

## Head

Offsets in this table are relative to decoded byte 0. Let `p` be the first
byte after the encoded map-name CString. For short ANSI `10.alm`, `p=0x0f`.
The map name is reopened as `Scenario\` plus the name when the mission number
is nonzero. Mission start obtains that number with `atoi(mapName)`.
— SAV-MAP-005, SAV-HEAD-025

| Decoded offset | Wire | Original source | Meaning or rule |
|---:|---|---|---|
| `0x00` | u32 | `world+0x04` | Sub-tick counter; zeroed at mission start |
| `0x04` | u32 | `world+0x00` | Full-clock value; restore the separately stored value |
| `0x08` | CString | `world+0x28` | Map name |
| `p+0` | 11 × u32 | See ordered list below | Literal saved values |
| `p+44` | u32 | `world+0x80` | Mission number |
| `p+48` | u32 | `world+0x84` | Difficulty; LOAD retains a stored value only in 1..3 |
| `p+52` | u32 | `playerList+0x20` | Meaning Unknown |
| `p+56` | u32 | Player list count | Followed by that many `objref<Player>` operations |

The eleven-value source order is
`+11c, +124, +128, +12c, +130, +134, +138, +13c, +148, +144, +140`.
The final three are not in address order. They are not the difficulty field
and receive no Boolean normalization or difficulty clamp.
— SAV-STREAM-010, SAV-HEAD-025, UNIT-GATE-012, SAV-657, SAV-WHEADLOAD-521

## Ordered grammar

`FUN_004d0cb7` writes the following after the head:

| Order | Wire | Source or exact serializer |
|---:|---|---|
| 1 | Player bodies, as counted by the head | `00527d70`, typed Player archive references |
| 2 | `list32<Unit>` | Dead-actor manager `00510401 -> 00527ac0` |
| 3 | u8 world-present discriminator | Writer derives 0/1 from `world+0x2c`; loader tests zero/nonzero |
| 4, world only | `list32<Building>` | `005102db -> 00527a10` |
| 5, world only | `list32<SpellEffect>` | `005114b6 -> 00527f90` |
| 6, world only | Block array, cell table, u32 terrain key | `00544a60`; [World](world.md) |
| 7, world only | raw 4374 | Session `00539310` |
| 8, world only | `list32<Sack>` | `0051149d -> 00527e90` |
| 9 | u32 discriminator | Writer emits `0xBADFACE1` |
| 10, marker matches | u32 | Global `[0x00609b0c]` |
| 11 | raw 400 | Block pointed to by `world+0x118`, `0053e800` |
| 12, odd endpoint only | 1 byte | Outer word-compression alignment; not a document member |

The Building, SpellEffect, Sack and dead-actor managers are reached through
four fields of the separate object at `world+0x14`: `+0`, `+4`, `+8`, `+c`
respectively. The dead list is not `world+0xc`. Every top-level list uses a
plain u32 count and complete archive references; an empty list is four bytes.
— SAV-ROSTER-024, SAV-SHAPE-023, SAV-DOC-053, SAV-657

## City and world shapes

Both shapes contain the entire roster, dead-list framing and trailer.
A city/no-world save omits Buildings, SpellEffects, terrain, session and Sacks
from the conditional part. It can retain the previous map name with mission
number zero. Unit/Human objects in the roster remain legal irrespective of
whether the conditional world half is present. The located no-world path
restores and detaches actors without constructing terrain. The narrowed
placement-latch claim supplies no universal relation between `Player+3d` and
the city/world discriminator.
— SAV-CITY-030, SAV-ROSTER-024, SAV-SHAPE-023, SAV-POSTLOAD-220

The roster is nested as Player -> Groups -> actors -> owned/equipped objects.
Group records are direct inline programmes, with no class tag. The map producer
creates Players from type-5 slots, actors from type-6 records plus the hero,
and Buildings from type-4 records. These relations do not fix the file counts.
Sacks and their contents are created at runtime.
— SAV-OBJ-016, SAV-MEMBER-036, SAV-DIARY-042, SAV-HUMAN-043

## Trailer

The loader always consumes a discriminator and 400 bytes. It consumes the
intervening global dword only when the discriminator is `0xBADFACE1`. Thus
marker-present and marker-absent trailers occupy 408 and 404 bytes. A marker
inside an earlier raw member does not terminate that member.
— SAV-TRAIL-026, SAV-DOC-053, SAV-FULLREAD-252, SAV-790

The 400-byte block is a separately allocated, exclusively owned world object,
zero-filled by its constructor and freed at world teardown. It is not session
storage or padding. Both archive directions transfer all `0x190` bytes.

| Block offset | Width | Meaning |
|---:|---:|---|
| `0x00` | 4 | Turn-trace toggle |
| `0x04` | 4 | Script-trace toggle |
| `0x08..0x18f` | 392 | Literal stored state; consumers Unknown |

Only the first two dwords have identified consumers. The remaining fields'
meanings and pointer/alias/computed access remain Unknown. The separate global `[0x00609b0c]`
has initialization and archive read/write sites but no established nonzero
ordinary producer.
— SAV-646, SAV-647, SAV-648, SAV-649, SAV-790, SAV-791

Command 46 parameter 80 reaches the toggle helper only after Player resolution,
`byte[cmd+4]==0`, nonzero `[0x5f21c4]` and unsigned `Player+68 > 50`.
Subcommands 3/19 toggle the first/second dword: zero becomes 1, any nonzero
becomes 0. Ordinary UI emission and production of the gate byte remain Unknown.
The helper's debug help names subcommand 19 `Script tracing on/off`.
— SAV-694, SAV-647, SAV-791

An odd logical endpoint receives one transport byte; its value is not
constrained or consumed by the document routine. An even endpoint receives
none. Fresh SAVE computes parity again rather than preserving a source pad.
— SAV-TRAIL-026, SAV-DECPAD-238

## Eleven world-head values

The constructor explicitly zeroes all eleven dwords. Its separate memset
`[+a4,+118)` and indexed initializer ending at `+114` do not cover them.
Mode setup writes `+c/+8/+150` elsewhere. Constructor zeros are not an
all-path first-SAVE vector. — SAV-WHEADINIT-520, SAV-WHEADLIMIT-525

| Field | Identified conditional consumer |
|---|---|
| `+11c` | Nonzero restricts map-unit creation to `Player+28==0`; zero also admits the later `Mission/Players` loader with `Monsters` fallback |
| `+124` | Zero admits the map-construction `Outposts` reader |
| `+128` | Zero selects `00530b60(actor)`; nonzero selects `0052fdf0(actor,Player+34,0)`, Defend |
| `+12c/+134` | Either admits a helper on a new hero; `+134` also admits it on eligible reconstructed actors |
| `+130` | Nonzero writes new-actor word `+8c=0x40` after its earlier derive |
| `+138` | Nonzero constructs `PlasmaSword` and dispatches new-actor `vt+3c` |
| `+148/+13c` | `+148` enables the modulo-5 full-clock watchdog; `+13c` selects actor `+c0` clearing or missing-actor countdown/replacement |
| `+144/+140` | Constructor and archive transport known; other consumers Unknown |

The hero helper writes six words 100 at actor `+102..+10c`, six bytes 100 at
`+10e..+113`, then dispatches `vt+50`. Its flag tests do not run on the
existing-hero return arm. The watchdog decrements a nonzero countdown;
finding zero on entry reloads 2 and attempts replacement.
— SAV-WHEADMAP-522, SAV-WHEADHERO-523, SAV-WHEADWATCH-524

No general post-load reset or nonzero ordinary producer is established for
these eleven values. Undecoded locations, aliases, computed calls and
unexpanded callees remain open.
— SAV-WHEADWATCH-524, SAV-WHEADLIMIT-525
