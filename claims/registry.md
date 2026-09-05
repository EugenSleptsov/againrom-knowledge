# Claim registry

This file is the public index of claim ledgers. The detailed research chronology,
confidence-review notes, experiment write-ups, disassembly and raw evidence live in the
private research repository named by [`SOURCE.md`](../SOURCE.md); they are not part of
the public functional edition.

A claim ID is permanent. Read the claim in its owning ledger and check
[`retracted.md`](retracted.md) for later correction or narrowing. The command
`go run ./tools/claim <ID>` performs both lookups without loading a whole ledger.

## ROM1 ledgers

| Area | Ledger |
|---|---|
| AI and targeting | [`ai.md`](ai.md) |
| ALM scenario/map data | [`alm.md`](alm.md) |
| Animation | [`anim.md`](anim.md) |
| Data.bin | [`databin.md`](databin.md) |
| Dialogue | [`dialogue.md`](dialogue.md) |
| Hall of fame | [`fame.md`](fame.md) |
| Heroes and character generation | [`hero.md`](hero.md) |
| Inventory | [`inv.md`](inv.md) |
| Items and equipment | [`item.md`](item.md) |
| Magic | [`magic.md`](magic.md) |
| Menus | [`menu.md`](menu.md) |
| Mission flow | [`mission.md`](mission.md) |
| Movement | [`move.md`](move.md) |
| Palettes | [`pal.md`](pal.md) |
| Party | [`party.md`](party.md) |
| Inline registry | [`reg.md`](reg.md) |
| RES/LM container | [`res.md`](res.md) |
| Saves | [`sav.md`](sav.md) |
| Session/campaign state | [`session.md`](session.md) |
| Shops | [`shop.md`](shop.md) |
| 16/16A sprites | [`spr16a.md`](spr16a.md) |
| 256 sprites | [`spr256.md`](spr256.md) |
| Tavern | [`tavern.md`](tavern.md) |
| Terrain and map presentation | [`terrain.md`](terrain.md) |
| Text and encoding | [`text.md`](text.md) |
| Town UI and flow | [`town.md`](town.md) |
| Triggers | [`trigger.md`](trigger.md) |
| Units and combat | [`unit.md`](unit.md) |
| Audio/video | [`video.md`](video.md) |

## ROM2 ledgers

| Area | Ledger |
|---|---|
| Engine-level observations | [`rom2-engine.md`](rom2-engine.md) |
| Session/network observations | [`rom2-session.md`](rom2-session.md) |

The source snapshot also contains `claims/rom2-asset.md`, but that ledger was not
included in public snapshot `k1`. Public ROM2 format pages that depend on `R2-ASSET-*`
therefore have incomplete public provenance until a separately reviewed functional
edition of that ledger is published. Do not treat the missing ledger as silently
approved or reconstruct it from other material.

## Status and confidence

- `✔ promoted` — the claim is also used by a public format/specification page.
- `● active` — current claim not promoted into a format page.
- `✖ retracted` — former wording has been corrected, narrowed or refuted; see
  [`retracted.md`](retracted.md).

Confidence records the strength and scope of the research evidence. It is not a legal
clearance, a compatibility guarantee or permission to redistribute original game
material.

## Publication boundary

The public ledgers should state independently expressed functional conclusions,
confidence, scope, status and the identity of the private evidence that earned them.
Instruction listings, decompiler output, reconstructable shipped-content tables and
internal review narrative are publication-review material under
[`PUBLICATION.md`](../PUBLICATION.md), not required parts of this index.
