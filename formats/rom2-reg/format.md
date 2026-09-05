# ROM2 inline registry (`&YA1`, nested) — identity survey

Level 3. Promoted, evidence-backed claims only. Basis: `R2-ASSET-006`. Cross-reference
only, not evidence: ROM1's own [`formats/reg/format.md`](../reg/format.md)
(`REG-FMT-031`, `REG-REC-032`) and `RES-SCOPE-015` (found by magic, not by name).

Seen as: any file-typed node, in any of the 11 `.res` containers, whose payload
begins with magic `&YA1` — the same by-magic search RES-SCOPE-015 used on ROM1,
not a `.reg`-extension filter.

## Result

Identical to ROM1's REG-FMT-031 grammar. The by-magic search over all 11
containers finds **19** nested registries (`graphics.res` 5, `scenario.res` 1,
`sfx.res` 1, `video.res` 8, `world.res` 2, `world_srv.res` 2), and **19/19 parse
cleanly**: the `0x18`-byte header reads, `recordCount` 32-byte records read, and
the payload tiles `0x18 + recordCount*32 + 4 + poolLen` exactly to the payload's
own end — 0 exceptions. `world.res` and `world_srv.res` each nest the identical
pair `data/ai.reg` (156 B) and `data/map.reg` (1020 B), matching the independent
size comparison in [`rom2-res/format.md`](../rom2-res/format.md). `R2-ASSET-006`.

Sampled key names (first 12 distinct per registry, printed only because the
payload's own tiling was already verified exact) read as plausible schema
identifiers, e.g. `Files`, `Global`, `Object0..Object106` (`objects/objects.reg`),
`Scanning`, `Tasker`, `IntelligentCons`, `MinimalGuardRan` (`data/ai.reg`),
`Common`, `Fading1`, `Fading2`, `startx`, `starty` (`video.res`'s `*/01.reg`
family).

## Not yet surveyed

Record VALUE content and the `kind` bitfield's meaning on ROM2 data (REG-KIND-033/
034/035 are ROM1 findings, not replayed here); whether any registry's key SET
differs from an equivalent ROM1 registry (no ROM1 registry of the same name was
compared field-for-field, only the container grammar).
