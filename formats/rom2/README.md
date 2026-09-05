# ROM2 format specifications (level 3)

Rage of Mages II asset and session formats, read from the preserved `gameversions/rom2-ru`
root (`pipeline/ROM2-MULTIPLAYER.md`). Kept apart from [`formats/`](../README.md), which is
ROM1-only: a page here cites only `R2-ASSET-*`/`R2-SESSION-*` claims, never a ROM1 claim
ID as evidence, though it may name a ROM1 page as cross-reference where the two formats
agree. `formats/README.md`'s own index and target-format table are unmodified by
this family.

Every page below is an **identity survey** result, not an independent re-derivation:
it reports whether a ROM1 grammar, applied unmodified to ROM2 bytes, parses cleanly,
and where it first diverges when it does not. A new field a divergence exposes is
named as unread, never guessed at.

| Surface | Seen as | Result | Page |
|---------|---------|--------|------|
| **RES container** | `*.res`/`*.RES`, 11 files at the install root | Identical to ROM1's `&YA1` container, 11/11, 0 tiling exceptions; extension vocabulary is NOT a subset of ROM1's — `.pkt`/`.smk` appear in neither ROM1 root | [`rom2-res/format.md`](../rom2-res/format.md) |
| **ALM map** | `*.alm`/`*.ALM` at the root + nested in `scenario.res`, magic `M7R\0` | File-header field layout identical, 83/83; record 0's declared size is not ROM1's and needs a documented +16 correction; once applied, the full record chain past record 0 is also identical, 83/83, 0 residue — `recordCount`/`formatVersion` remain genuine content differences | [`rom2-alm/format.md`](../rom2-alm/format.md) |
| **Data.bin database** | `world.res:data/data.bin`, `world_srv.res:data/data.bin`, root `templates.bin` | `world.res`'s own data.bin matches ROM1's wire grammar through group A's title array (102 shared bytes) and 2 more exact-matching collections, then diverges inside group C; `templates.bin` is a distinct, earlier-diverging file whose own wire grammar is Unknown (a rival head reading fits the same bytes) | [`rom2-databin/format.md`](../rom2-databin/format.md) |
| **Inline registry** | nested `&YA1` payloads inside `.res` containers, found by magic | Identical to ROM1's REG-FMT-031 grammar, 19/19 | [`rom2-reg/format.md`](../rom2-reg/format.md) |
| **Sprite / palette** | `*.16a`, `*.16`, `*.256`, `*.pal` inside `.res` containers | `.16a`/`.256`/`.pal` container framing identical over the full population (`.pal` also 156/156 on a strict field-value test); the `.256`/`.16` files with trailer residue are sha256-identical to ROM1 RU's own — a reproduced ROM1 phenomenon, not a ROM2 extension | [`rom2-spr/format.md`](../rom2-spr/format.md) |
| **Text** | `main.res:text/*.txt`, `patch.res:patch.txt` | Container convention preserved and extended (52 more files); byte-range encoding differs (66.4% vs ROM1 RU's 100.0000%; 0 of ROM2's high bytes fall in one of the two source blocks at all) | [`rom2-text/format.md`](../rom2-text/format.md) |
| **Session wire frame** | the 8-byte record header read by `CBufferManager::ReceiveData` in `allods2.exe`/`a2server.exe` | One 8-byte header, four fixed-offset fields (payload length, codec selector, codec input size, passthrough byte), read identically in both binaries; a second transport's 150-byte bound is the same 142-byte payload cap restated with the header included, not an independent constant | [`rom2-net/format.md`](../rom2-net/format.md) |

## Not yet surveyed

Sprite pixel/RLE decoding, palette colour values, the `.alm` record body past
record 0, `data.bin`'s groups D-H, the specific new fields the divergences above
expose, the session wire frame's own record payload grammar past its 8-byte
header, and any semantic decoding of ROM2-specific content are all out of
frame for this container/record-layout identity survey and remain open for
later work.
