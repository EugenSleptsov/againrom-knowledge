# ROM2 RES container (`&YA1`) — identity survey

Level 3. Promoted, evidence-backed claims only. Basis: `R2-ASSET-001` (container
identity), `R2-ASSET-012` (`world.res` vs `world_srv.res`). Cross-reference only,
not evidence: ROM1's own [`formats/res/format.md`](../res/format.md)
(`RES-MAGIC-001`, `RES-HDR-002`, `RES-ACCEPT-031`).

Seen as: 11 files at the ROM2 install root — `MUSIC.RES`, `graphics.res`,
`main.res`, `movies.res`, `patch.res`, `scenario.res`, `sfx.res`, `speech.res`,
`video.res`, `world.res`, `world_srv.res`.

## Result

ROM2's `.res` container is the same format as ROM1's, byte for byte, at the level
this survey tested. `internal/rom.OpenArchive` — ROM1's own reader, unmodified —
opens all 11 files; every one carries magic `26 59 41 31` ("&YA1"). Replaying
RES-ACCEPT-031's stricter corpus invariant (every file node's payload inside
`[24,regOff)`, exact contiguous tiling with no gap or overlap, printable-ASCII
names) finds **0 violations across all 11 containers** — the same invariant ROM1's
own corpus satisfies. `R2-ASSET-001`.

The entry-extension vocabulary found inside these 11 containers is **not** a
subset of ROM1's own. 12 of ROM2's 14 non-empty extension categories match
ROM1's exactly (`.16`, `.16a`, `.256`, `.alm`, `.bin`, `.bmp`, `.dat`, `.pal`,
`.reg`, `.txt`, `.wav`, plus the extensionless entry both roots have), but
`.pkt` (2 files — `data/itemname.pkt` in `world.res` and `world_srv.res`) and
`.smk` (8 files, all in `video.res` — Smacker video) appear in neither ROM1 EN's
nor ROM1 RU's own container population, confirmed by running the identical
extension census rooted at each ROM1 install as a positive control.

## Population note

This survey's own brief named 10 top-level `.res` files; the actual population is
**11** (`MUSIC.RES` was the omission). A count taken from a list rather than a
directory read undercounts here.

## Client/server split

`world.res` and `world_srv.res` hold the identical 5-entry `data/` directory
(`ai.reg`, `data.bin`, `itemname.bin`, `itemname.pkt`, `map.reg` — 0 client-only, 0
server-only names). 4 of 5 are byte-identical in size; `data/data.bin` differs
(client 149 971 B, server 152 527 B) — see
[`rom2-databin/format.md`](../rom2-databin/format.md) for what those two files' own content
comparison shows. `R2-ASSET-012`.

## Not yet surveyed

Whether every entry's PAYLOAD (not just its size) is identical between `world.res`
and `world_srv.res` for the 4 same-size entries; whether any other container pair
(there is only one client/server pair, `world`/`world_srv`) exists.
