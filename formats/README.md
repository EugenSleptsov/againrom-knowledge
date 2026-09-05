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
| **REG** | `*.reg` (inside `.res`) | Nested archive / sprite manifest | ☑ [specified](reg/format.md) (`REG-FMT-031`, `REG-REC-032`, `REG-KIND-033`/`REG-KIND-034`, `REG-VAL-024`…`REG-VAL-029`, `REG-KEY-044`, `REG-ROSTER-052`, `REG-CUT-053`; 64 claim rows); `data/map.reg`'s own per-key table open |
| **SPR256** | `*.256` (inside `.res`) | 8-bit paletted sprite, multi-frame | ◐ [structure](spr256/format.md) done (`SPR256-STRUCT-001`…`SPR256-CORPUS-006`); RLE/palette/roles open |
| **SPR16A** | `*.16a`, `*.16` (inside `.res`) | 16-bit-word sprite | ☑ [specified](spr16a/format.md) (`SPR16A-STRUCT-001`, `SPR16A-RLE-002`, `SPR16A-RLE-003`, `SPR16A-PIX-011`); framebuffer packing + `.16` fonts open |
| **PAL** | `*.pal` (inside `.res`) | Palette — the per-tier and per-owner recolour | ☑ [specified](pal/format.md) (both shapes, the shade-table build, the `Palette` selector and `face`; the owner index’s writer open) |
| **ALM** | `*.alm`/`*.ALM`, in root + `scenario.res`, magic `M7R\0` | Map | ☑ [specified](alm/format.md) (`ALM-HDR-001`, `ALM-FRAME-031`, `ALM-META-008`…`ALM-META-010`, `ALM-GRID-012`…`ALM-GRID-014`, `ALM-CLS-035`…`ALM-CLS-038`, `ALM-CORP-060`; 87 claim rows); the placeable-definition database's populating file, the trigger's 64-byte junk field and `Target_Item` open |
| **SAV** | `game*.sav` | Save game, magic `Asg&2a` | ◐ [specified (partial)](sav/format.md) (`SAV-HDR-001`, `SAV-FRAME-021`, `SAV-STREAM-010`/`SAV-STREAM-013`, `SAV-OBJ-014`/`SAV-OBJ-016`, `SAV-CELLLOAD-108`…`SAV-CELLLOAD-113`, `SAV-FULLREAD-252`; 300 claim rows); container/transport/object framing specified, most of eleven observed classes' fields named only by offset |
| **FAME** | `famehall.dat` | Hall-of-fame records | ◐ [specified (single-sample)](fame/format.md) (`FAME-HDR-001`…`FAME-WRITE-007`, `FAME-DEFAULT-008`; 8 claim rows); writer confirmed byte-exact on replay, the two trailing per-record fields' meaning open |
| **MENU** | `main.res:graphics/mainmenu/*` | Main-menu asset contract (mask + overlays) | ☑ [specified](menu/format.md) (base/mask/overlays + hit-test + placement; disable-source/cmd-ids open) |
| **VIDEO/MUSIC** | `Allods/VIDEO*.RES`, `Allods/MUSIC.RES` | A/V payloads | ☑ [specified](video/format.md) (functional edition; music in `MUSIC.RES` — `VIDEO-MUSIC-001`…`VIDEO-MUSIC-012`; SFX — `VIDEO-SFX-013`…`VIDEO-SFX-021`; cutscenes — `VIDEO-029`…`VIDEO-036`, `VIDEO-045`…`VIDEO-051`; 36 claim rows); Smacker decoder internals and hardware/driver timing intentionally out of scope |
| **TERRAIN** | `graphics.res:terrain.3d` | Tile graphics + cell geometry | ◐ [specified](terrain/format.md) (tile word → pixel, geometry, passability) |
| **DAT** | `world.res:data/data.bin` | Placeable-definition database | ◐ [specified](databin/format.md) (grammar + the 3 placement tables) |
| **TEXT** | `main.res:text/*.txt`, `patch.res:patch.txt` | One-byte string tables + display/input encoding | ☑ [specified](text/format.md) (functional edition; `TEXT-STRTAB-023`, `TEXT-NAMEIN-024`, `TEXT-COLL-025`, `TEXT-CHARGEN-027`…`TEXT-CHARGEN-029`, `TEXT-UI-032`…`TEXT-UI-047`; 45 claim rows); a Unicode mapping and malformed-byte behaviour are intentionally out of scope |

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
| **AI** — target acquisition, the diplomacy matrix, and idle/guard/patrol behaviour | [`ai/format.md`](ai/format.md) | ◐ the consulted relation, candidate population, selection, the three radii, who runs the AI and how often, the group/per-actor state machines, guard and patrol specified (`AI-FILTER-001`…`AI-GUARD-007`, `AI-TICK-008`…`AI-DIFF-016`; 192 claim rows); order execution (`formats/move`), the four player-issued group orders, the candidate scorer and the line-of-sight predicate open |
| **MISSION** — starting a campaign mission and ending it in a win | [`mission/format.md`](mission/format.md) | ◐ the player-placement routine, the drop-cell RNG, the seat search, the four placement arms, the definition lookup and the entry/two-container writer boundary specified (`MISSION-DROP-002`, `PARTY-ORIGIN-010`…`PARTY-GATE-013`, `MISSION-VICTORY-035`; 38 claim rows in `claims/mission.md`); the `DataBinID == 26` sentinel arithmetic and the inter-mission campaign-state arm open |
| **TOWN** — town-exterior and tavern-interior reactions and animation | [`town/format.md`](town/format.md) | ◐ shop/tavern/school/gate entrance reactions, the bird/horse/baba/dervish episodes, tavern-interior draw/clocks and school training presentation specified (`TOWN-158`, `TOWN-211`, `TOWN-399`…`TOWN-447`); the wider town/world-map/rooms domain in `claims/town.md` (224 rows) is mostly not yet reified in a format page |

## ROM2

ROM2 (Rage of Mages II) work in this repository is format research only — an identity
survey of ROM2 bytes against ROM1's own grammar, not a second game this project builds
(owner direction). The family index, with its own per-surface result column, is
[`rom2/README.md`](rom2/README.md); compiled-code relationship to ROM1
(`claims/rom2-engine.md`, `R2-ENGINE-*`) is a structural binary comparison, not a stored
byte format, and has no page here.

| Format | Seen as | Nature | Status |
|--------|---------|--------|--------|
| **RES** | 11 files at the ROM2 install root | Container archive, magic `&YA1` (ROM1's own) | ☑ [identity survey](rom2-res/format.md) (`R2-ASSET-001`, `R2-ASSET-012`; 0 violations, 11/11 clean); `world.res`/`world_srv.res` payload-level equality on the 4 shared same-size entries open |
| **ALM** | `*.alm`/`*.ALM` at the root + nested in `scenario.res`, magic `M7R\0` | Map | ◐ [identity survey](rom2-alm/format.md) (`R2-ASSET-002`, `R2-ASSET-003`, `R2-ASSET-017`…`R2-ASSET-027`, `R2-SESSION-010`; header + record framing/dispatch 83/83 clean once a +16 correction is applied); record content past the 20-byte header, for any of the 83 maps, undecoded |
| **Data.bin** | `world.res`/`world_srv.res:data/data.bin`, root `templates.bin` | Placeable-definition database | ◐ [identity survey](rom2-databin/format.md) (`R2-ASSET-004`, `R2-ASSET-005`, `R2-ASSET-025`, `R2-ASSET-026`, `R2-ASSET-029`…`R2-ASSET-033`; `world.res`/`world_srv.res` grammar closed, 0 residue); `templates.bin`'s own grammar and every group's entry content undecoded |
| **Inline registry** | nested `&YA1` payloads inside `.res` containers, found by magic | Nested archive / key-value tree | ◐ [identity survey](rom2-reg/format.md) (`R2-ASSET-006`; 19/19 clean tiling); record value content and the `kind` bitfield unsurveyed on ROM2 data |
| **Sprite / palette** | `*.16a`, `*.16`, `*.256`, `*.pal` inside `.res` containers | Sprite + palette containers | ◐ [identity survey](rom2-spr/format.md) (`R2-ASSET-007`…`R2-ASSET-009`; container framing 623/623, 1916/1929, 159/159 and 3/3 across the four kinds, every divergence cross-referenced to a pre-existing ROM1 finding); RLE/pixel decoding and palette colour values unsurveyed |
| **Text** | `main.res:text/*.txt`, `patch.res:patch.txt` | String tables | ◐ [identity survey](rom2-text/format.md) (`R2-ASSET-010`, `R2-ASSET-011`; container preserved and extended); byte-range encoding measured at 66.36% overlap with ROM1's own two source blocks, not 100%, and the specific 8-bit encoding is not named |
| **Session wire frame** | `CBufferManager::ReceiveData`, read by `allods2.exe`/`a2server.exe` | 8-byte socket record header | ◐ [identity survey](rom2-net/format.md) (`R2-SESSION-003`; the header's four fixed-offset fields specified); the record's own payload grammar past the header is unsurveyed |

## Explicitly out of scope (not ROM1-native)

Standard or third-party formats we only *identify*, never re-derive:

- `Map Editor.opt` — Microsoft OLE2 compound document (`D0 CF 11 E0 …`).
- `*.bmp`, `*.wav` inside archives — standard Windows BMP / RIFF WAVE.
- GOG / system wrappers: `goggame-*`, `unins000*`, `ddraw.dll`, `smackw32.dll`,
  `aqrit.cfg`, `webcache.zip`, `Help/*.htm`, `Hints/*.gif`, `*.ico`, `*.lnk`.
- `rom.exe`, `Map Editor.exe` — PE executables (the engine itself, not a data format).

Status key: ☐ not started · ◐ in progress · ☑ specified (core) · ★ complete.
