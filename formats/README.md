# Format specifications (level 3)

One folder per ROM1-native format, each holding a `format.md` that contains **only
promoted, evidence-backed claims**. Anything speculative stays in an experiment
until it earns its way here.

## Target formats

Derived purely from the bounded install census promoted as `INV-CORPUS-001`,
`INV-SIG-002`…`004` and `INV-SCOPE-005`.
Status reflects how much of the format is specified, not whether a file merely
opens.

| Format | Seen as | Nature | Status |
|--------|---------|--------|--------|
| **RES / LM** | `*.res`, `KIDS.LM` | Container archive, magic `26 59 41 31` | ☑ [specified](res/format.md) (core; `@4`/`@12`/node`@0` open) |
| **REG** | `*.reg` (inside `.res`) | Nested archive / sprite manifest | ☐ not started |
| **SPR256** | `*.256` (inside `.res`) | 8-bit paletted sprite, multi-frame | ◐ [structure](spr256/format.md) done (`SPR256-STRUCT-001`…`SPR256-CORPUS-006`); RLE/palette/roles open |
| **SPR16A** | `*.16a`, `*.16` (inside `.res`) | 16-bit-word sprite | ☑ [specified](spr16a/format.md) (`SPR16A-STRUCT-001`, `SPR16A-RLE-002`, `SPR16A-RLE-003`, `SPR16A-PIX-011`); framebuffer packing + `.16` fonts open |
| **PAL** | `*.pal` (inside `.res`) | Palette — the per-tier and per-owner recolour | ☑ [specified](pal/format.md) (both shapes, the shade-table build, the `Palette` selector and `face`; the owner index’s writer open) |
| **ALM** | `*.alm`/`*.ALM`, in root + `scenario.res`, magic `M7R\0` | Map | ☐ not started |
| **SAV** | `game*.sav` | Save game, magic `Asg&2a` | ☐ not started |
| **FAME** | `famehall.dat` | Hall-of-fame records | ☐ not started |
| **MENU** | `main.res:graphics/mainmenu/*` | Main-menu asset contract (mask + overlays) | ☑ [specified](menu/format.md) (base/mask/overlays + hit-test + placement; disable-source/cmd-ids open) |
| **VIDEO/MUSIC** | `Allods/VIDEO*.RES`, `Allods/MUSIC.RES` | A/V payloads | ☐ not started |
| **TERRAIN** | `graphics.res:terrain.3d` | Tile graphics + cell geometry | ◐ [specified](terrain/format.md) (tile word → pixel, geometry, passability) |
| **DAT** | `world.res:data/data.bin` | Placeable-definition database | ◐ [specified](databin/format.md) (grammar + the 3 placement tables) |

## Not a file format — simulation areas with their own spec

The same three levels apply to engine *behaviour* we have read out of `rom.exe`: an experiment, a
claim ledger, a promoted spec. These folders hold no file layout.

| Area | Spec | Status |
|------|------|--------|
| **MOVE** — unit movement & path selection | [`move/format.md`](move/format.md) | ◐ search, costs, termination, route, reservation, step, refresh policy specified; tick order + the blocked-cell verdict table open |
| **SHOP** — stock, price & trade | [`shop/format.md`](shop/format.md) | ◐ object graph, generator, candidate pool, enchantment stage, both price formulas and the refusal specified; the unit price's own derivation, the RNG seed and whether stock survives a save open |
| **TAVERN** — mercenary hire | [`tavern/format.md`](tavern/format.md) | ◐ the type space, the shelf gate, the hire vector, the price and level ladders and the death/recovery rule specified; what a mercenary *is* once spawned, the inn's art and the multiplayer path open |
| **HERO** — chargen, the derived-stat graph, the combat loop | [`hero/format.md`](hero/format.md) | ◐ see `claims/hero.md` |
| **MAGIC** — the spellbook, the cast, the resistance | [`magic/format.md`](magic/format.md) | ◐ see `claims/magic.md` |
| **ANIM** — the simulation ↔ presentation boundary | [`anim/format.md`](anim/format.md) | ◐ see `claims/anim.md` |
| **SESSION** — the game's own lifecycle: object, phase, clock, map load, mission end | [`session/format.md`](session/format.md) | ◐ the two counters, the rate ladder, the map-load order, the mission-end arm and the hero-creation chain specified; where a **win** is decided, the trigger machinery and the network client's clock open |
| **DIALOGUE** — the window a mission's script speaks through, and its lifecycle | [`dialogue/format.md`](dialogue/format.md) | ◐ the announcement transport, the client dispatcher, the window class and its six entries, both negatives, the pager, the tag scan and the wrap/clamp specified; the conditional tag arms' effects, the composed speech name and the lose chain past `0x41e` open |
| **TRIGGER** — the runtime that evaluates a map's authored mission script | [`trigger/format.md`](trigger/format.md) | ◐ the one-second pass, both dispatch tables read whole, the comparison arms, the fire-once latch, the save contract, the binder and the win/lose path specified; the helper routines behind a handful of arms open; what turns the outcome announcement into `0x41d` is specified for the win side in [`dialogue/format.md`](dialogue/format.md) |
| **UNIT** — the non-hero actor: template → instance → combat inputs | [`unit/format.md`](unit/format.md) | ◐ creation specified end to end (ctor defaults, 38 streamed slots, equipment, the one-time fold) and there is **no derive** on this arm; `actor+0x14`, the equipment-name parse and six constructor call sites open |
| **ITEM** — the item, its container and the sack | [`item/format.md`](item/format.md) | ◐ the four item classes, the definition binding, the container and its load, the fourteen equipment slots, both move commands, the pick-up, death and the sack's whole lifecycle specified; the `Armors`/`Shields`/`Magic Items` columns, `item+0x44`'s value space and the per-class `Equip` bodies open |

## Explicitly out of scope (not ROM1-native)

Standard or third-party formats we only *identify*, never re-derive:

- `Map Editor.opt` — Microsoft OLE2 compound document (`D0 CF 11 E0 …`).
- `*.bmp`, `*.wav` inside archives — standard Windows BMP / RIFF WAVE.
- GOG / system wrappers: `goggame-*`, `unins000*`, `ddraw.dll`, `smackw32.dll`,
  `aqrit.cfg`, `webcache.zip`, `Help/*.htm`, `Hints/*.gif`, `*.ico`, `*.lnk`.
- `rom.exe`, `Map Editor.exe` — PE executables (the engine itself, not a data format).

Status key: ☐ not started · ◐ in progress · ☑ specified (core) · ★ complete.
