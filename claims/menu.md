# Claim registry — MENU (menu surfaces: the main menu, the documents panel, the in-play Esc menus)

Level 2 ledger. Index: [registry.md](registry.md) · spec: [`formats/menu/format.md`](../formats/menu/format.md). Legend: `✔ promoted` (also in the spec) · `● active` · `✖ retracted`. IDs are permanent.

This ledger covers the **composition contract** of a menu surface — which archive nodes form it,
what its hit regions mean, and where the engine places each element. It began with the main menu
(`MENU-ASSET-001`…`MENU-STATE-007`) and was extended on 2026-08-13 to a second full-screen bitmap
surface, the campaign documents panel (`MENU-DOC-009`), which is built the same way out of
`graphics.res` rather than `main.res`. On 2026-08-14 it was extended again, to the two in-play
menus Esc raises (`MENU-ESC-010`…`MENU-INPUT-016`). Those two are not bitmap surfaces: they carry
no art of their own and are drawn as a sprite nine-patch over a darkened frame, so the ledger's
subject is now the menu surface rather than the full-screen bitmap. The BMP *encoding* itself is
standard Windows BMP (out of scope to re-derive).

| ID | Claim | Confidence | Status | Evidence |
|----|-------|-----------|--------|----------|
| MENU-COMBAT-017 | The mission right column is one fixed 160-pixel container with four ordered children; the command panel is the second child and not the minimap or character panel. | High | ● active | [EXP-0190](../experiments/EXP-0190-combat-controls/) |
| MENU-COMBAT-018 | The eight command cells are 34x34 at panel-local `(8+34c,7+34r)`, row-major, and all visible states come from four full-panel BMPs. | High | ● active | [EXP-0190](../experiments/EXP-0190-combat-controls/) |
| MENU-COMBAT-019 | The eight cells and tooltip indices are one exact table: | High | ● active | [EXP-0190](../experiments/EXP-0190-combat-controls/) |
| MENU-ASSET-001 | The main menu is the 18-file `graphics/mainmenu/` subtree of `main.res`: a 640x480 24bpp base, a 640x480 8bpp hit mask, and two 24bpp overlay sets for buttons 1..8. | High | ✔ promoted | [EXP-0026](../experiments/EXP-0026-mainmenu-assets/) |
| MENU-ASSET-002 | There are **exactly 8 buttons**. | High | ✔ promoted | [EXP-0026](../experiments/EXP-0026-mainmenu-assets/) |
| MENU-MASK-003 | `menumask.bmp`'s 256-entry palette is the **identity grayscale ramp** `palette[i]=(i,i,i)` (all 256 entries) — it carries no colour meaning; the raw **8-bit index** is the semantic marker. | High | ✔ promoted | [EXP-0027](../experiments/EXP-0027-menu-hitmask/) |
| MENU-MASK-004 | The hit-test `FUN_00486260` reads the mask byte under the cursor and maps eight index values, 0x80 to 0xf0 in steps of 0x10, to buttons 1..8; every other index is no button. | High | ✔ promoted | [EXP-0027](../experiments/EXP-0027-menu-hitmask/) |
| MENU-GEOM-005 | Overlay placement is **two static `rom.exe` tables**: normal/hover at `DAT_0059a2d8` (VA; file `0x198ed8`) and pressed at `DAT_0059a358`, contiguous, 8 entries × 16 B = `{x,y,w,h}` int32 LE. | High / Medium | ✔ promoted | [EXP-0028](../experiments/EXP-0028-menu-overlay-geometry/) |
| MENU-GEOM-006 | Overlays draw **1:1** (no scale): each table entry's `(w,h)` equals the corresponding BMP's pixel dimensions — **16/16** exact (normal + pressed). | High | ✔ promoted | [EXP-0028](../experiments/EXP-0028-menu-overlay-geometry/) |
| MENU-STATE-007 | `button%d.bmp` is the **hover/highlight** overlay, `button%dp.bmp` the **pressed** overlay. | Medium | ● active | [EXP-0027](../experiments/EXP-0027-menu-hitmask/) |
| MENU-STRTAB-008 | Which surfaces the one global string index space names, from the sites read this round — the map a consumer needs before it can renumber anything. | Medium | ● active | [EXP-0143](../experiments/EXP-0143-actor-names/) |
| MENU-DOC-009 | The documents panel, whole: eleven bitmaps, three hit rectangles that equal their bitmaps, 21-line pages, and one entry point. | High / Medium | ● active | [EXP-0151](../experiments/EXP-0151-mission-documents/) |
| MENU-ESC-010 | Esc during play raises two different surfaces, one per UI state, and both are panels of one family rather than the main menu. | High | ✔ promoted (partially retracted) | [EXP-0162](../experiments/EXP-0162-ingame-menu/) |
| MENU-ITEM-011 | The mission menu constructs eight entries and shows seven; two of them are exclusive. | High / Medium | ✔ promoted | [EXP-0162](../experiments/EXP-0162-ingame-menu/) |
| MENU-ITEM-012 | The town menu is a different class with five entries, and one of them is a different action with a different word. | High / Unknown | ✔ promoted | [EXP-0162](../experiments/EXP-0162-ingame-menu/) |
| MENU-KEY-013 | The effective accelerator is the letter after `~` in the label, not the constructor's immediate, and the difference is visible in the shipped English text. | High | ✔ promoted | [EXP-0162](../experiments/EXP-0162-ingame-menu/) |
| MENU-ART-014 | Both Esc menus load no art of their own; the frame is a nine-patch out of one sprite bank, `graphics\interface\lm.256`, plus an 8-px drop shadow. | High / Unknown | ✔ promoted | [EXP-0162](../experiments/EXP-0162-ingame-menu/) |
| MENU-STOP-015 | Both Esc menus stop the world, and the mechanism is the idle-handler gate, not `MISSION-STOP-016`'s `server+0x2c`. | High | ● active | [EXP-0162](../experiments/EXP-0162-ingame-menu/) |
| MENU-INPUT-016 | Esc opens and Esc closes; the panel takes keyboard focus among its own rows and, as the root's capture object, receives every mouse message while it is up. | High / Unknown | ● active (partially retracted) | [EXP-0162](../experiments/EXP-0162-ingame-menu/) |
| MENU-CURSOR-046 | The main-menu transition is the one surface transition that ends on `select` rather than `default`. | High / Medium | ● active | [EXP-0214](../experiments/EXP-0214-cursor-surfaces/) |

### MENU-COMBAT-017

**The mission right column is one fixed 160-pixel container with four ordered children; the command panel is the second child and not the minimap or character panel.** `FUN_00472070` constructs root `campaign+0xcc`, map `+d0`, right column `+d4`, minimap `+d8`, command panel `+dc`, character panel `+e0` and filler `+e4`; `FUN_00472800` appends `d8,dc,e0,e4` under `d4`, then `d0,d4` under the root. Command-panel screen rectangles are `(480,158)-(640,238)`, `(640,158)-(800,238)` and `(864,158)-(1024,238)` at the three shipped resolutions. The command-panel class is `FUN_004901f0`, vtable `0x0059a6f0`, id 6, 160x80, with no children. Inventory, spellbook and Esc hit regions belong to the separate character-panel sibling. Its nested inventory-grid double-click can enter the same item Cast/transfer vmethod as character-panel left-up. Its physical left-up separately transfers the selected item through grid vslot `+0xa4` to destination code 2 and emits `0x22/0x32`, including when the hit-test returns `-1`. A gold-sentinel cell with a positive purse instead attaches the always-created Drop Gold modal at `campaign+0x10c`, rect `(100,H-200)-(396,H-32)`: edit id `0x989685`, action `0x989681`/`0x445`, cancel `0x989682`/`0x446`; Enter/Esc alias those buttons

**Confidence.** High (construction arguments, append order, inventory/dialog/button vtables, physical handlers and all child constructors are read end to end; the rectangles are exact evaluation of hard-coded arguments)

### MENU-COMBAT-018

**The eight command cells are 34x34 at panel-local `(8+34c,7+34r)`, row-major, and all visible states come from four full-panel BMPs.** Inactive draws `HeadsR.bmp`; active draws `CommandBarR.bmp`, then each disabled cell's same source rectangle from `CommandEmpR.bmp`, then the selected cell from `CommandDnR.bmp`. There is no hover bitmap, separate border or separate icon, and no child button object. All four assets are 160x80x24 uncompressed BMPs and byte-identical across EN/RU. Active, enable-mask, selected-index and dirty state are panel fields `+0x68,+0x64,+0x60,+0x6c`. A margin or disabled left-down clears only the selected overlay and sends index -1, leaving a previously armed map mode intact

**Confidence.** High (one draw routine and one left-down handler read whole, plus BMP-header and hash measurements on both roots)

### MENU-COMBAT-019

**The eight cells and tooltip indices are one exact table:** Attack/A, Move/M, Guard/G, Defend/D, Cast/C, Swarm/S, Stand Ground/T, Retreat/R. Cell 0..7 returns `main.txt[0..7]` only when active, enabled, no modal/special cursor is present and frame-state bits `0x0a` are clear. The RU file carries the corresponding CP866 strings `Атаковать`, `Идти`, `Охранять`, `Защищать`, `Колдовать`, `Идти в боевой готовности`, `Держать позицию`, `Отступить`, with the same Latin accelerators. Left down acts; double-click aliases it; left-drag re-enters it on every delivered move; right up cancels through map message `0x405`; the other right edges are no-ops. **G2:** cell geometry and index order are `.text`; changing them requires executable edits. Text and the four full-panel images are data, so changing a label or any embedded icon changes only its shipped resource bytes

**Confidence.** High for the event routes, predicates, indices and EN/RU bytes (handlers read whole and both resource roots measured)

### MENU-ASSET-001

The main menu is the **18-file** `graphics/mainmenu/` subtree of `main.res`: `menu_.bmp` (640×480, 24bpp base brooch background), `menumask.bmp` (640×480, **8bpp** hit mask), `button1..8.bmp` (24bpp hover overlays) and `button1..8p.bmp` (24bpp pressed overlays). All 18 sha256 distinct; overlay dims per `inventory.csv`. `rom.exe` requests them as `main\graphics\MainMenu\{menu_.bmp, MenuMask.bmp, button%d.bmp, button%dp.bmp}` (referenced only by the loader `FUN_00485e00`)

**Confidence.** High (the asset set is not inferred from the directory listing but from the **only** referrer of those path strings; "18 files" is then a listing)

### MENU-ASSET-002

There are **exactly 8 buttons**. The loader `FUN_00485e00` `sprintf`s `button%d.bmp`/`button%dp.bmp` in a loop `n = 1..8` (`do{…}while(n<8)`, start 1), storing hover overlays into the array at `this+0x6c` and pressed overlays at `this+0x80`, mask at `this+0xbc`, background at `this+0xb8`

**Confidence.** High (the count is the loop bound itself — `do{…}while(n<8)` from 1 — not a count of files that happen to exist)

### MENU-MASK-003

`menumask.bmp`'s 256-entry palette is the **identity grayscale ramp** `palette[i]=(i,i,i)` (all 256 entries) — it carries no colour meaning; the raw **8-bit index** is the semantic marker. The mask is loaded by `FUN_00429ab0` as raw indices (buffer at maskObj+0x10, `w·h` bytes, never rendered to colour)

**Confidence.** High (the identity ramp is a measurement over all 256 entries; "the index is the semantic marker" is carried by the loader, which keeps the buffer as raw indices — an absence of any colour conversion on that path)

### MENU-MASK-004

The hit-test `FUN_00486260` reads `idx = maskBuf[(y−top)·640 + (x−left)]` (640 = mask width) and maps eight index values to buttons: **0x80→button1, 0x90→button2, 0xa0→button3, 0xb0→button4, 0xc0→button5, 0xd0→button6, 0xe0→button7, 0xf0→button8** (asset number = idx/16 − 7). Index 0 (background) and every other value (incl. the edge ramp 0x10..0x1e = hot>>3 and stray AA pixels) fall to the switch default = **no button**. Verified: only these 8 of the 43 present indices are "hot", and each region is bracketed by its button's placement rect 8/8

**Confidence.** High (the eight mappings are `switch` cases read off the hit-test, and the *default* — every other index, including the edge ramp — is read too, which is what makes "no button" a fact rather than an assumption)

### MENU-GEOM-005

Overlay placement is **two static `rom.exe` tables**: normal/hover at `DAT_0059a2d8` (VA; file `0x198ed8`) and pressed at `DAT_0059a358`, contiguous, 8 entries × 16 B = `{x,y,w,h}` int32 LE. `FUN_00485c20` `SetRect`s each button at `(x,y)`→`(x+w,y+h)` offset by the menu origin. Origin = (0,0) (full-screen 640×480 mask; hit-test indexes it from (0,0)). Values listed in `placement-tables.csv` / `rom-layout-excerpt.md`

**Confidence.** High (static tables read at named addresses and the `SetRect` that consumes them) / Medium (the origin `(0,0)` — deduced from the mask being full-screen and indexed from `(0,0)`, not from an instruction that sets it)

### MENU-GEOM-006

Overlays draw **1:1** (no scale): each table entry's `(w,h)` equals the corresponding BMP's pixel dimensions — **16/16** exact (normal + pressed). Each normal rect **brackets** its button's mask region — **8/8**. These two discriminating checks jointly bind BMP dims ↔ hit-test index ↔ placement table into one consistent geometry

**Confidence.** High (16/16 exact and 8/8 bracketing across three independently-derived sources — dims from the files, indices from the mask, rects from the binary; a wrong pairing anywhere breaks one of the two. This row is the shape the confidence scale asks for and states its own discriminator without prompting)

### MENU-STATE-007

`button%d.bmp` is the **hover/highlight** overlay, `button%dp.bmp` the **pressed** overlay. In `FUN_00486260`: not-pressed & hovering button i → draw `this+0x6c[i]` at normal rect `this+0x94[i]`; mouse-down while still hovering the latched button i → draw `this+0x80[i]` at pressed rect `this+0xa8[i]`. `this+0xe4` is a per-button **disable bitfield** — bit i set suppresses the overlay. The click dispatcher `FUN_00486530` posts a per-button command message (button 8 → WM_CLOSE 0x10); message-id meanings beyond WM_CLOSE are undecoded

**Confidence.** Medium

### MENU-STRTAB-008

**Which surfaces the one global string index space names, from the sites read this round — the map a consumer needs before it can renumber anything.** Every entry is `MOV ECX,[0x005eb3d4]` then a subscript of `4 × index` (`TEXT-STRTAB-023`). Index **77** is the dialogue panel's button (`004219e3`, `FUN_004217be`). Indices **140/141** are the mission-outcome panels (`004740e2`, `00474172`). Indices **47..50** are the multi-selection panel, `FUN_00491930` reading `+0xbc/+0xc0/+0xc4/+0xc8` — a selected-unit **count**, not names — and the same routine reads `class+0xd8` (`InfoPicture`) off the `units.reg` class array. Index **256** is the text-entry field's tooltip, returned by `FUN_00432db0` when its owner's `+0x1e0` is non-zero, and that class is `TEXT-NAMEIN-024`'s ten-character name field. The in-mission message layer `FUN_004104e8` resolves **26** distinct indices — 85..89, 129, 142..149, 204..209, 221..226 — a contiguous-block pattern that matches item-pickup, alliance, join/rejoin and cheat messages, all of which concatenate a **player** name. `EnumRefs refto:5eb3d4` returns 230 hits over **47** owners, so this map covers a minority of the consumers: the remaining sites fold the index into a memory displacement and were not resolved to an index this round. The *other* fifteen tables are not reached this way at all — they are read through `FUN_004687f0(table, i)` with a table-local index, 302 hits over 57 owners — and the two surfaces this round followed there are the unit information panel reading `unitname.txt` (`UNIT-NAME-039`) and the character screen reading `npcnames.txt` at `FUN_0041fe3f` `00420003`/`0042002f` under an index built from the mage bit `+0x18c & 4`

**Confidence.** **Medium.** Each named index is a cited instruction and each agrees with the EN text at that line of `text/main.txt`, which is corroboration across two artefacts rather than a read of the drawing call; and the coverage is explicitly partial — 5 of 47 owners mapped. The three name-validation messages at 193/194/195 have **no** located consumer

### MENU-DOC-009

**The documents panel, whole: eleven bitmaps, three hit rectangles that equal their bitmaps, 21-line pages, and one entry point.** `FUN_0048de00` is the only referrer of any of the six document art path literals and loads them all: `+0x6c` = `graphics\interface\Docs\sheet.bmp`, a 3-element array at `+0x74` = `Arrows\{00_l,01_l,11_l}`, a 3-element array at `+0x88` = `Arrows\{00_r,01_r,11_r}`, a 4-element array at `+0x9c` = `OK\{Ok_off,Ok_on,Ok_l_off,Ok_l_on}`. Geometry is `FUN_0048dcf0`, in screen pixels: left arrow `(0,200,56,240)`, right arrow `(576,200,636,240)`, OK `(560,416,604,448)` — each **equal to its bitmap's pixel size**, 56×40 / 60×40 / 44×32, measured from the BMP headers on both roots, which checks the two readings against each other. `sheet.bmp` is 640×480×24 and `1.bmp` 464×344×24. `FUN_0048e700` draws sheet, left arrow, right arrow, OK, then the current document at the panel origin. A picture is blitted at the document RECT's top-left `(92,72)` at natural size, **8 px wider than the 456-px RECT**; text is drawn by `FUN_004571c0` with `ECX = [0x005e9bc0]`, `TEXT-API-007`'s font4, from line `+0x38` to `+0x38+0x15` with the pitch taken from the glyph sprite because the caller passes 0. `FUN_0048e820` selects array element 1 on hover and 2 on press (3 for OK) and 0 on leave, so the arrow file names read as `00` idle / `01` hover / `11` pressed; `Ok_on` is selected by no arm read here. Paging is `FUN_00486b80`/`FUN_00486ba0`, a step of **`0x15` = 21 lines** on `+0x38` bounded by the line count at `+0x18`, and `FUN_0048e5e0`/`FUN_0048e630` change **document** only when the page step returns 0, so one arrow pair walks pages first and documents second. OK posts `0x445`. The panel is 248 bytes, window id `0x4ba`, constructed at `00472742` into `app+0x3a4`; `FUN_00476710` is its only shower, called once, from the campaign state machine `FUN_00473110` at `004749c8` under `campaign+0x3dc == 1`, on the arm of **message `0x463`** — one of 112, read out of the PE's own index and jump tables on both roots

**Confidence.** High for the asset set (the only referrer of the literals), the geometry (immediates, cross-checked against BMP headers on both roots) and the page step / High for "one message raises it": the switch is read from the image's own tables rather than a listing, and both roots agree / **Medium** for "one entry point in the image": a control with the numerically equal id `0x463` exists in `FUN_004320d0`'s own id space and whether its notification can reach this handler was not ruled out

### MENU-ESC-010

**Esc during play raises two different surfaces, one per UI state, and both are panels of one family rather than the main menu.** `FUN_00472b80` is the frame window's `WM_KEYDOWN` handler — read off the MFC message map at `0x00599c50`, whose six dwords are `{0x100, 0, 0, 0, 0x10, 0x00472b80}` and which reproduces on both roots. Its `VK_ESCAPE` arm loads the UI state word (`00472c4c MOV ECX,[ESI + 0x3dc]`) and branches: `== 1` posts **`0x416`** (`00472c52 CMP ECX,0x1`, `00472c5c PUSH 0x416`), `== 0` **and** the `campaign+0x3b4` CString empty posts **`0x41f`** (`00472c6d CMP ECX,EDI`, `00472c7b CMP [ECX + -0x8],EDI`, `00472c89 PUSH 0x41f`), anything else falls to `00472e4b`, which forwards `WM_KEYDOWN` to the root container `campaign+0xcc`. `SESS-SCREEN-003` fixes `campaign+0x3dc == 1` as a map session on screen and `SHOP-TOWN-023` fixes `== 0` as the town. `FUN_00473110` case `0x416` re-checks `== 1` and builds `FUN_0043c23b` at `(100,60)-(440,400)` = **340×340** (`0047344d PUSH 0x190 / 0x1b8 / 0x3c / 0x64`); case `0x41f` is ungated and builds `FUN_0043c90b` at `(100,100)-(440,340)` = **340×240** (`00474652 PUSH 0x154 / 0x1b8 / 0x64 / 0x64`). Both then take the shared `0047466c … CALL 0x00476810`. The same `0x416` is posted by the ~~command panel's own~~ *character-panel sibling's* 32×32 button at panel-local `(0x7e,0xce)-(0x9e,0xee)` (`004916b6`, in `FUN_00491330`, gated `campaign+0x3dc & 1`), whose tooltip `FUN_00490c60` resolves to global string index 14 — `text/main.txt` line 15, `Main Menu <ESC>` on the EN root and the same line in Russian on the RU root. Neither path reaches the main-menu window: `callto:485960` is **1 hit / 1 owner / 0 orphan** (`00472407`), and the routine that shows that window, `FUN_004769c0`, first frees the session's object lists and plays `music\menu.wav` — it is the exit path, not a pause

**Confidence.** High for the handler identity (the message map is data, not a decompiler rendering, and both roots' bytes agree), the two branch conditions, both rects (immediates verified through the PE section table on both roots, `evidence/anchors.txt`) and the single main-menu construction (enumeration instrument named, 0 orphan) / the *meaning* of `campaign+0x3dc == 1` and `== 0` is taken from `SESS-SCREEN-003` and `SHOP-TOWN-023`, not re-derived here

**Amended.** The owning-panel noun is partially retracted; the button rectangle, tooltip and posted message stand. The correction is in `retracted.md` and `MENU-COMBAT-017`.

### MENU-ITEM-011

**The mission menu constructs eight entries and shows seven; two of them are exclusive.** `FUN_0043c23b` builds each row as `FUN_00449110(id, FUN_004687f0(index), font [0x005e88b8], 0, message, accelerator, colourPtr)` and stacks it with `FUN_0043c129(button, 0x1e)`. In screen order, with `text/dialogs.txt` line and posted message: `~Save Game` `0x41a`; then either `~Load Game` `0x418` when `campaign+0x6bc == 2` **or** `Diplomacy` `0x43c` otherwise — the two write the same local and only one `FUN_0043c129` follows; then `Game ~Options` `0x41b`, `Sou~nd Options` `0x422`, `~Quest Objectives` `0x420`, `~End Quest` `0x41c`, `~Return to Game` `0x446`. Label indices are `0x22 0x23 0x4c 0x24 0x25 0x26 0x27 0x28` on the descriptor `0x005ea678`, which `004711e7 PUSH 0x5be274` / `004711ec MOV ECX,0x5ea678` loads from `main\text\dialogs.txt`, and `FUN_004687f0(table,i)` resolves as `[0x005eb3d4][table->+0xc + i]`. Four rows carry a disable: `~Save Game` when `campaign+0x6bc` is 0 or 1; `~Load Game` when `FUN_0043c0ae()` is 0, and that routine is a `FindFirstFileA` on `<game dir>\game*.sav`; `Sou~nd Options` when `FUN_00449d90()` returns 5, its `this+0x9c == 0` arm; `~Quest Objectives` when `campaign+0x6bc != 2`. The **first argument** to `FUN_00449110` is a control id, **not** the row: it runs 1,2,3,4,5,6,8 with 7 taken by `Diplomacy`, while the row order is the `FUN_0043c129` call order above. `0x446` is the family-wide close, not a distinct action: `FUN_004c52f3` turns `0x445`/`0x446` into the standard `0x44c` teardown, and 39 owners in the image post it

**Confidence.** High for the entry list, its order, the messages, the label file binding and the exclusive pair (one routine read whole, the loader binding two adjacent instructions, and the labels re-read out of both roots' `main.res` where the RU strings are the same menu in different bytes) / **Medium** for the four disable rules: each is one cited call whose predicate was read, but whether any other site re-enables a row later was not swept

### MENU-ITEM-012

**The town menu is a different class with five entries, and one of them is a different action with a different word.** `FUN_0043c90b`, vtable `0x00598008` against the mission menu's `0x00597ef8`, same button and row helpers. In screen order: `~Save Game` `0x41a`, `~Load Game` `0x418`, `Sou~nd Options` `0x422`, **`Abort Game`** (label index `0x4d`) `0x41c`, `~Return to Game` `0x446`. No `Game Options`, no `Quest Objectives`, no `Diplomacy`, and Save precedes Load where the constructor's control ids are 2 and 1 — another instance of `MENU-ITEM-011`'s id-is-not-row point. A search of the constructor for the disable call `(*button)vt+0x1c(1,0)` returns **0**: every row is built enabled. `0x41c` is the same message the mission menu's `~End Quest` posts; `FUN_00473110`'s `0x41c` arm branches on `campaign+0x3dc`, and the two surfaces take **different** branches of it: `== 1` raises the five-row confirmation `FUN_0043c66f` (`Change Map` or `~Victory!` at `0x41d`, `~Exit to Main Menu` and `Exit to ~Windows` both at `0x41e`, `~Return to Game` at `0x446`), `== 0` raises the three-row `FUN_0043cb6b` (the same two exit rows and `~Return to Game`), both at `(100,100)-(440,340)`. So one message, two panels, and the mission row's own word for it is `~End Quest` while the town's is `Abort Game`

**Confidence.** High for the entry list, its order, the messages, the absence of disables and the two confirmation panels' own rows (three routines read whole, labels re-read on both roots) / **Unknown** for what the two exit rows do: both post `0x41e` and differ only in control id, and how the exit target is distinguished was not established

### MENU-KEY-013

**The effective accelerator is the letter after `~` in the label, not the constructor's immediate, and the difference is visible in the shipped English text.** `FUN_004be6e3` stores the passed accelerator, then walks the label from index 1: `~~` is an escape and skips two, otherwise the first character preceded by a single `~` replaces it via `FUN_00468ac0(ch) & 0xff`. `FUN_00468ac0` lowercases — `FUN_005562d0` normally, and when `[0x005eb57c] == 1` it lowercases CP866 Cyrillic instead (`+0x20` over `0x80..0x8f`, `+0x50` over `0x90..0x9f`). Consequences measured on both roots (`evidence/labels.txt`): on the EN root eleven of thirteen labels carry a `~` and take it, and the two that do not — `Diplomacy` and `Abort Game` — keep the constructor's `D` and `E`; on the RU root **all thirteen** carry a `~`, including the rows those two correspond to, so those two entries have a language-dependent accelerator while the code has not changed. The RU labels also place the `~` away from the first letter where the first letters would collide — mission rows 5, 6 and 7 mark bytes `0xa0`, `0xaa` and `0xad` at label offsets 1, 2 and 3 — which the immediate model cannot produce at all. One residue: mission row 5 passes `0x4d` = `M` and its EN label is `~Quest Objectives`, so the shipped accelerator is `Q` and the immediate is inert. A reimplementation that hard-codes the accelerators reproduces the English release and diverges on the Russian one

**Confidence.** High. The scan is one routine read whole, and the discriminator is not corpus agreement: the two roots carry disjoint label bytes and disagree about which entries carry a `~` at all, which no single-root reading could have separated from *the immediate is the accelerator*

### MENU-ART-014

**Both Esc menus load no art of their own; the frame is a nine-patch out of one sprite bank, `graphics\interface\lm.256`, plus an 8-px drop shadow.** The family's art-load slot `vt+0x78` is `FUN_004490e0`, whose whole body is `RET`. Painting is `FUN_004c4da2` through the global `[0x005ef980]`, which has exactly **one writer and one reader**: `refto:5ef980` = 29 hits / 2 owners / 0 orphan, the write `0046a903 MOV [0x005ef980],EAX` in the interface loader `FUN_00469f80` and 28 reads all inside `FUN_004c4da2`. The store pairs with the construct that precedes it, `0046a8e7 PUSH 0x5bd2b8` = `graphics\interface\lm.256` into the `.256` constructor `0x00428c50` (its two neighbours use the `.bmp` constructor `0x00429330`). Frames 0..8 are the nine-patch and their measured sizes match the placement arithmetic exactly: centre `96×64`, corners `48×48`, top and bottom edges `96×48`, left and right edges `48×64`, against an inset of `0x30` = 48, a horizontal step of `0x60` = 96 and a vertical step of `0x40` = 64. Frames 9..17 are a second nine-patch of the same shape at 32/48 px and are not referenced by `FUN_004c4da2`. The file is 39 943 bytes with the same sha256 on both roots. Two passes: `vt+0x1c` draws frames 3,6,8,7,5 — the right column and bottom row only — over a rect offset `+8,+8`, then `vt+0x18` draws all nine at the panel rect shrunk by 8 on right and bottom. Rows are laid out by `FUN_0043c129`: `panel+0x70 = panel+0x78`, `panel+0x78 = top + 0x1e`, then `button->SetRect(&panel+0x6c)` and `panel->AddChild(button)`, with `FUN_0044a360` forcing `panel+0x6c = 0x28` and `panel+0x74 = width − 0x30` and seeding top/bottom from `(0,0,0xf0,0x28)` — so each row is panel-local `(40, 40+30(n−1), width−48, 40+30n)`

**Confidence.** High for the bank's identity, its exclusivity, the frame sizes and the row layout (the enumeration names its instrument and returns one writer with 0 orphan; the frame sizes are read from the `.256` structure on both roots and agree independently with the code's immediates; the row helper is nine instructions read raw) / **Unknown** for how the tiling covers a frame width that is not `96 + 96k`: the counts are `(w−96)/96` and `(h−96)/64` with truncating division (`004c4f3a SUB EAX,0x60` / `CDQ` / `MOV ECX,0x60` / `IDIV`), which for this panel's 332×332 frame leaves a 44-px band inside the right edge and one inside the bottom edge that no placed tile accounts for. EXP-0162 states this as a prediction with its refutation rather than as a claim

### MENU-STOP-015

**Both Esc menus stop the world, and the mechanism is the idle-handler gate, not `MISSION-STOP-016`'s `server+0x2c`.** Both `FUN_00473110` arms end at `FUN_00476810`, which sets bit 3 of `campaign+0x3dc` (`0047682e OR AL,0x8`) — or bit 15 for the one panel stored at `campaign+0x3a8`, which is not either of these — and the MFC idle handler tests it once in the whole image: `0047167f TEST EAX,0x4008`, `EnumRefs imm:4008` = **1 hit / 1 owner / 0 orphan**, verified byte-for-byte on both roots. `DLG-STOP-012` measured that branch: it runs no pacer arm at all, which is stricter than the pause, and the timer rival is closed by an import-table absence (`SESS-TIMER-022`). Neither Esc path writes, reads or reaches `server+0x2c`; `MISSION-STOP-016`'s two writers are `004d051c` and `004d0957`, in the session-start and mission-teardown routines, and neither is on this path. Neither panel is stored in a `campaign+…` slot, so its close takes `FUN_004757b0`'s default arm, and `DLG-CLOCK-014` fixes what that does: clear bit 3, and if the word is then exactly 1 and the pacer phase is 2, zero the phase so the stopped time is **discarded** rather than caught up. In the town the word returns to 0 and the restart correctly does nothing. Compositing is `DLG-DIM-013` unchanged: between a DirectDraw Lock and Unlock, `00476861 PUSH 0x3` and `00476869 CALL 0x0044fad0` remap **the whole screen** in place through the shroud table `[0x005e8420]` at level 3, a per-channel gain of 13/16 by `TERR-FOG-084`'s law — one destructive pass, not a per-frame blend and not a translucent draw, and it survives only because this gate stops everything that would repaint

**Confidence.** High for the gate, its single-hit enumeration, the shade level and the absence of `server+0x2c` from this path (each a named instruction, the immediates re-read through the PE section table on both roots) / the *consequences* of the gate and of the close arm are `DLG-STOP-012` and `DLG-CLOCK-014`, both High, applied here rather than re-derived; this row establishes that the Esc menus take that path, not that the path behaves as those rows say

### MENU-INPUT-016

**Esc opens and Esc closes; the panel takes keyboard focus among its own rows and, as the root's capture object, receives every mouse message while it is up.** With the panel up `campaign+0x3dc` is 9, so neither arm of `FUN_00472b80`'s `VK_ESCAPE` case matches and the key is forwarded to the root container, reaching the panel's key slot `FUN_0043c1e6` — `0x26`/`0x28` move the focused row, everything else falls to `FUN_004c55d8`, which on `0x1b` posts `0x446`, the same message `~Return to Game` posts. The same word being 9 also makes the character panel's Main Menu button inert while the panel is up, because `FUN_00473110`'s `0x416` arm requires exactly 1. `FUN_00476810` adds the panel as the **last child of `campaign+0xcc`** (`00476830 MOV ECX,[EBX + 0xcc]`), the same parent the map view and the side panels have, and `FUN_004bccf4` only appends and sets the parent pointer. `FUN_004bd9dc`, the base dispatcher, sends a mouse message to the capture object at `this+0x34` when one is set and otherwise walks the children, delivering to the first whose rect contains the point and stopping there; ~~`FUN_004c5275`, the panel's show, sets `this+0x5c = 1` and calls `FUN_004bd0a4(1,0)`, which moves focus among the panel's **own** children and sets no capture. A click outside the panel is therefore delivered to the window under the cursor~~ *`FUN_004c5275`, the panel's show, sets `this+0x5c = 1`, calls `vt+0x24(1)` (`004c5290`), which is `FUN_004bd456` in both class tables and stores the panel in `root+0x34`, then `vt+0x28(1)` and `FUN_004bd0a4(1,0)`. `FUN_004bd9dc` gives a mouse message to `root+0x34` before any child, the root's own mouse slots are stubs, and the panel's children do not contain a point outside it, so a click outside the panel reaches no other child of the root.*

**Confidence.** High for the open/close path and for the routing rule (three routines read whole, and the forwarding target `0x00472e4b` is fixed by the `JNZ rel32` in the same listing) / **Unknown** for ~~what a click outside then *does*: the map view's and side panels' own handlers were not read this round, so this row claims delivery and not effect. EXP-0162 carries that as an open question~~ *what the panel's own handlers do beyond the slots read*

**Amended.** The mouse-capture clause of the headline and the last two body sentences are partially retracted; the open and close path, key forwarding, the `campaign+0x3dc` value 9, the append as last child of `campaign+0xcc` and the dispatcher rule stand. The correction is in `retracted.md`.

### MENU-CURSOR-046

**The main-menu transition is the one surface transition that ends on `select` rather than `default`.** `FUN_004769c0` is identified as the main-menu routine by its own `PUSH 0x5be5fc` at `0x00476b7f`, which is `music\menu.wav`; the push site was found by a whole-image dword scan of the string block rather than by a literal heuristic. It sets `campaign+0x3dc` bit 7 (`OR CL,0x80` at `0x00476b30`, `TOWN-373`), loads `wait` (slot 26) at `0x00476a8f` and calls the set-cursor adapter at `0x00476a93`, and at its tail loads **`select`** (slot 5, `0x005ef9c8`) at `0x00476b3f` and calls the adapter at `0x00476b43`. The other eleven routines of the same family end on `default` (slot 0) instead (`TOWN-372`). The same routine also pushes `World\Data\` at `0x004769ea`, `-serverid` at `0x00476c80`, `.srv` at `0x00476cd7` and `default.srv` at `0x00476d37`, so its range extends to at least `0x00476d37`

**Confidence.** High for the two cursor sets and the bit, all named instructions in the complete 69-caller population (`AI-CURSOR-175`), and for the slot identities, which come from `SPR16A-CURSOR-067`. **Medium** for the surface identification: it rests on the routine pushing the menu music path, and the alternative that `music\menu.wav` is loaded by some other screen is not excluded by the literal alone


`EXP-0214` was allocated `MENU-CURSOR-046`..`MENU-CURSOR-050` (5 ids) and spent
`MENU-CURSOR-046`, 1 id. **`MENU-CURSOR-047`..`MENU-CURSOR-050` (4 ids) are returned
unused**, none ever reissued. The next free `menu.md` id is therefore `MENU-CURSOR-051`.


## Help panel and Cast key

| ID | Claim | Confidence | Status | Evidence |
|---|---|---|---|---|
| MENU-051 | The campaign frame's key handler posts help message `0x434` on F1 in any state; the frame builds the help panel only when `campaign+0x3dc` equals exactly 1, so F1 is ignored in the town and over any panel that sets an overlay bit. | High | ● active | [EXP-0437](../experiments/EXP-0437-f1-help/) |
| MENU-052 | Help is the Pause panel's class and constructor call with the help text in place of the pause line: id 1, rectangle 576x384 snapped to 488x360 and centred, one OK button, no title, scrolling body in font 1. | High / Medium | ● active | [EXP-0437](../experiments/EXP-0437-f1-help/) |
| MENU-053 | Help stops the world through the Esc menus' idle gate and closes on the OK button or Esc; F1 over help is ignored; whether keys scroll the body depends on a focus rule not resolved here. | High / Medium / Unknown | ● active | [EXP-0437](../experiments/EXP-0437-f1-help/) |
| MENU-054 | The C key reaches the Cast arm only in a map session with no text entry open, a nonzero selection count and `view+0x144 & 0x24` clear; otherwise it falls to a second dispatch whose C entry returns 0. It reaches nothing in the town. | High | ● active (amended) | [EXP-0437](../experiments/EXP-0437-f1-help/) |
| MENU-055 | C arms Cast mode 5 only when a selected object set `view+0x144` bit `0x200`; for a nonempty selection with the bit clear the key is consumed with no message, sound call or state change. An armed Cast posts `0x408` and chooses no spell. | High / Medium / Unknown | ● active (amended) | [EXP-0437](../experiments/EXP-0437-f1-help/) |
| MENU-056 | F1 in the campaign frame reaches no accelerator, help file or WinHelp call; its only path is the `text/help.txt` panel. The library `ID_HELP` machinery is in the image and untraced. Bounded to the population named. | Medium | ● active | [EXP-0437](../experiments/EXP-0437-f1-help/) |

### MENU-051

- The F1 arm of the frame's `WM_KEYDOWN` handler `FUN_00472b80` is `00472c9a MOV EAX,[ESI+0x1c]` / `PUSH EDI` / `PUSH EDI` / `00472c9f PUSH 0x434` / `PUSH EAX` / `00472ca5 CALL [0x00632f5c]` (`PostMessageA` on the frame window), then `JMP 00472e5d`. No compare on `campaign+0x3dc` or `campaign+0x6bc` precedes it. Pause (key `0x13`) tests `+0x6bc == 2` and `+0x3dc == 1` before it acts, and F2 tests `+0x6bc == 2` and `+0x3dc` in {0, 1}, so the absence is a property of this arm.
- The frame procedure `FUN_00473110` case `0x434` is `00473a28 CMP [EBP+0x3dc],EBX` / `JNZ 00473193`, with `EBX = 1` set at `0047315d`. A different word leaves the case without building a panel or posting a message.
- Accepted: word 1, the map-session-only state of `SESS-SCREEN-003`. Ignored: word 0, the town (`SHOP-TOWN-023`), and every word with an overlay bit, which includes the Esc menu and Pause panel (exactly 9, `MENU-STOP-015`), any dialogue panel (`DLG-CLOCK-014`) and the help panel itself (`MENU-053`).
- The identity of the F1 key with this arm (the virtual-key compare and its dispatch) rests on `AI-KEY-125` and its `keyboard.tsv`; the committed listing of `FUN_00472b80` holds the arm.
- `rom.exe` is byte-identical on both roots (SHA-256 `942e9b72610eeba2f3d74930ee47cbc4476b85a942c14348f874ec7b7d367d03`), so EN and RU share the rule.

**Confidence.** High for both arms: each is a named instruction and the compare operand is fixed by the listing's raw bytes. The live alternatives named by the question (accepted in every state, accepted only in a map session, accepted in a map session and the town) are excluded by the missing compare in the key arm and the exact-equality compare in the frame case.

**Unknown.** Every value `campaign+0x3dc` can take: `SESS-SCREEN-003` leaves bits 1, 2, 9, 10 and 13 individually Unknown, so the ignored set is stated by the one accepted value. Frame states in which `FUN_00472b80` is not the key handler (the main menu surface) were not examined. Whether the Drop Gold modal (`MENU-COMBAT-017`), the chat text entry (`campaign+0xc0`) and the spellbook popup change the word: none is shown to, so F1 may open help over any of them.

### MENU-052

- Construction: `00473a52 MOV EDX,[0x005eb4d0]`, then pushes of `0`, `0`, the text pointer, `0x1b0`, `0x260`, `0x30`, `0x20`, `1` and `CALL 0043fbbf`. The Pause arm of `FUN_00472b80` makes the same call with `[[0x005eb3d4]+0x1dc]` as the text (line 119 of the shared string index space, `TEXT-STRTAB-023`). Both panels then go through `FUN_00476810` (`MENU-STOP-015`). The class vtable is `0x005982b8`, over the modal panel base `FUN_0043eb6a`, which stores the text pointer, a null second string and button case 0.
- Size and position: `FUN_004c553f` replaces the rectangle with `w = ((cw-8)/0x60)*0x60+8`, `h = ((ch-0x68)/0x40)*0x40+0x68` and centres it (`DLG-PANEL-035`). The constructor rectangle `(left, top, right, bottom)` is `(0x20, 0x30, 0x260, 0x1b0)`, 576x384, which snaps to 488x360. Origins are (76,60) at 640x480, (156,120) at 800x600 and (268,204) at 1024x768.
- Build routine `FUN_0043ebc5`: no title or footer string, one OK button of case 0 (control id 4, message `0x445`) whose label is the first line of `dialogs.txt` read through `FUN_004687f0(0)`, and a body created through vtable slot `+0x88`, which builds `FUN_004be03c` (a text control, vtable `0x0059b278`) over body rectangle `(40,56,448,272)` in font 1 `[0x005e88b8]` with colour pointer `[0x005ba21c]`.
- Text control: the base constructor `FUN_004c168a` shrinks the bottom to whole rows of font height plus 4 (11 rows, bottom 267, height 211); `FUN_004be03c` then sets the line pitch to font height plus 2 (17 for font 1's 15) and the visible count to height divided by pitch (12). `FUN_004be44e` adds a scroll bar (width 24, control id `0xdf23`, class `FUN_004c35cd`) and narrows the text width by `0x1a` when the wrapped line count times the line height exceeds the body height, then rewraps at the narrower width. `FUN_004be3a7` paints through `FUN_00457400` (`DLG-LINE-038`).
- `text/help.txt` gives a scroll bar on both roots (`TEXT-088`).

**Confidence.** High for the construction call, the shared class, the rectangle snap and the widget sequence (each routine read, operands traced to immediates). Medium for the body height 211 and the 12-line figure: they combine immediates read from `FUN_004c168a`, `FUN_004be03c` and `FUN_004be44e` by arithmetic, and no frame was observed.

**Unknown.** The colour `[0x005ba21c]` points to, and so the displayed colour of the help text (`MISSION-MSGLINE-056` documents the ink ramp `0x005e8878` for other callers; this pointer's identity was not traced). Native pixels.

### MENU-053

- Stop: the help panel is shown by `FUN_00476810`, which ORs `8` into `campaign+0x3dc` (`MENU-STOP-015`), darkens the screen once (`DLG-DIM-013`) and is stored in no `campaign+…` slot, so its close takes `FUN_004757b0`'s default arm (`DLG-CLOCK-014`): bit 3 is cleared and the stopped time is discarded. The word while help is open is 9, so F1 over help meets `MENU-051`'s compare and is ignored.
- Close: the panel's message handler `FUN_004c52f3` treats `0x445` and `0x446` alike, hiding the panel through vtable slot `+0x84` and posting `0x44c`. The panel's key handler `FUN_004c55d8` turns Esc (`0x1b`) into `0x446`. With help open the frame handler's Esc arm does not open the Esc menu: the word is neither 1 nor 0, so the arm falls through to the forwarder that delivers the key to the root and so to the panel (`MENU-INPUT-016`). The OK button's key handler `FUN_004bf182` posts its own message on Enter (`0xd`).
- Focus and navigation: `FUN_0043ebc5` calls vtable slot `+0x1c(0x10, 1)` on the OK button; `FUN_004c5369` moves focus by Tab and the four arrow keys over the neighbour links at `+0x48..+0x54`.
- Scroll: the scroll bar's key handler `FUN_004c224e` acts on Page Up, Page Down, Up and Down (`0x21`, `0x22`, `0x26`, `0x28`) by calling its set-position slot `+0x80` with a one-line or one-page step and then posts `0x46d`; it acts only when its own flag test `vtable+0x20(4)` succeeds and otherwise defers to the base handler. The frame handler forwards those four keys to the root.
- No key pages the text and no markup paginates it in the routines read: the body scrolls through the scroll bar.

**Confidence.** High for the stop mechanism, the close messages and the Esc routing (instruction reads, callee behaviour cited from the named claims). Medium for the scroll key set and its one-line or one-page steps (the handler's four arms are read whole, the step operands from `FUN_004c1939`, `FUN_004c1993`, `FUN_004c18e1` and `FUN_004c1913`).

**Unknown.** Whether the scroll bar's flag test succeeds while the OK button holds focus, so whether Up, Down, Page Up and Page Down scroll help from the keyboard at all. The scroll bar's mouse handling (clicks, thumb drag) and any mouse-wheel route were not read.

### MENU-054

- The map view's key handler `FUN_0040f38e` first requires `campaign+0x3dc == 1` (`0040f3c8`) and, when `campaign+0xc0` is nonzero (the chat text entry of `AI-KEY-125`), forwards the key to that object and returns. Digits take a separate arm.
- The letter group is dispatched through a byte table at `0x0040fd91` and a target table at `0x0040fd6d` (`JMP [EDX*4+0x0040fd6d]`). Its guard is `0040f681 CMP [EDX+0x3dc],1`, `0040f691 CMP [EAX+0x140],0` (the selection count; zero skips the group), `0040f6a1 MOV EDX,[ECX+0x144]` / `0040f6a7 AND EDX,0x24` / `JNZ 0040f789`. Key `C` (`0x43`, index 2 from base `0x41`) maps to table entry 1, the block at `0040f755`.
- A zero selection count (`JZ` at `0040f698`) or a set `0x24` bit (`JNZ` at `0040f6ac`) jumps to `0040f789`, which starts a second dispatch: key minus `0x20` indexes the byte table at `0x0040fddd` and the target table at `0x0040fda5`. For C the index is `0x23`, the byte is 13 and entry 13 is `0x0040fc47`, which jumps to the epilogue with `EAX = 0`. Only the block at `0040f755` ends with `EAX = 1`. When the text entry is open the handler forwards the key and returns 0 as well.
- That block is `CALL 0041ddea` (does a spellbook popup exist) / `TEST EAX,EAX` / `JNZ 0040f76b`, else `PUSH 5` / `CALL 0041b439`; it then sets `EAX = 1`. A second C press while the popup exists is consumed with no further call.
- Bit `0x4` of `view+0x144` is ownership by another player and bit `0x20` the structure flag (`AI-PANEL-061`). The town has no map view and no state word of 1, so the key reaches nothing there.
- The other Cast controls route to the same arm and are cited, not re-derived: the command-panel Cast cell (`MENU-COMBAT-018`, `MENU-COMBAT-019`), the quick-spell keys F5 to F8 (`AI-QUICKINVOKE-279`, `AI-KEY-125`) and the item Cast action through `FUN_0041b439(0x0a)` (`AI-PANEL-123`).

**Confidence.** High: `FUN_0040f38e`'s compares are read in the listing, the entry for `C` is read from the PE's table bytes (byte table `00 08 01 02 …`, `C` at index 2 yields 1, target table entry 1 is `0x0040f755`; the second dispatch's tables are committed beside it), and `FUN_0041ddea` and `FUN_0041b439` are read whole.

**Unknown.** What the caller does with the returned 0 (whether the key is passed on and who consumes it). How the selection count and ownership bit behave for selections that mix owners: `AI-SPELLCAP-288` bounds the capability side only.

**Amended.** `MENU-063` answers the Unknown about the caller of the returned 0: the frame handler never reads it and always calls the MFC default, so 0 only lets the root offer the key to its other children and its own key routine. The mixed-owner selection question stays Unknown.

### MENU-055

- `FUN_0041b439(5)` computes the enable mask with `FUN_00419e7b`: `view+0x144 & 4` returns 0, otherwise `0xef`, plus `0x10` when `view+0x144 & 0x200` (`AI-PANEL-060`). Cast is mode 5, tested as `mask & (1 << 4)`, so it passes only when `0x200` is set. `0x200` is set when any accepted selected object is a `CUnit` whose `+0x20` is `0x17` or `0x18` and whose `+0x18` is nonzero (`AI-SPELLCAP-288`).
- Refused: for a nonempty selection with `view+0x144 & 0x24` clear and the bit clear (no spell-capable object, that is no valid caster), `FUN_0041b439` does nothing after the mask test. Before the test it only calls `AfxGetThread` and `FUN_00419e7b` and, when the stored mode is already 5 and the request is not, closes the popup (not reached for a request of 5). It posts no message line, plays no sound and leaves `view+0x99c` unchanged. `FUN_0040f38e` still returns 1, so the key is consumed. An empty selection or a set `0x24` bit does not reach this routine (`MENU-054`).
- Armed: with the bit set, `view+0x99c` becomes 5 and the command panel receives `0x40d` with index 4. If no spellbook popup exists, `FUN_0041de5d` runs: it sizes the popup region (`FUN_004bcba8`), posts `0x408` to the right-column container `campaign+0xd4` and calls `FUN_0041e06d`, which records the popup's row count and sets `view+0x74 = 1`. No spell is chosen here: the order built on the next map click reads the stored current spell or item (`AI-SPELLGUARD-289`, `AI-PANEL-123`), and a negative result emits no cast.
- The no-spell case is not a separate refusal: the capability bit tests the selection, not the book.

**Confidence.** High for the mask test, the refusal path having no message or sound call and the armed-mode writes (`FUN_00419e7b` and `FUN_0041b439` read whole). Medium for the spellbook opening through `0x408`: the post is read but the receiver of `0x408` was not traced.

**Unknown.** Whether anything preselects a spell after C when none was chosen earlier in the session (the initial value of the controller's current spell was not traced), what `0x408` shows in the right column, and whether `FUN_00453b08`, called inside `FUN_0041de5d`, plays a sound.

**Amended.** `MENU-064` and `MENU-065` answer the Unknown: `0x408` is a refresh posted to the right-column container (three child panels write one field each) and shows no popup; the popup is `campaign+0xec`, and no code in the C path writes its current spell. Whether `FUN_00453b08` plays a sound stays Unknown.

### MENU-056

- Population searched: the PE resource directory of `rom.exe` (type ids 1, 3, 5, 12 and 14 only: no accelerator table, type 9); the image bytes for `.htm`, `.chm` and `winhlp` (no hit) and for `.hlp` and `winhelp` (one hit each, at file offsets `0x0019d52c` and `0x001cd86c`); and the references to every symbol or string containing `Help` (`EnumRefs sym:Help`).
- The path literal `main\text\help.txt` has one reference, the startup routine `FUN_004709e0` at `0047118f` (`TEXT-086`). The `WinHelpA` import has two call sites: `FUN_005822d6`, the frame-destroy routine, which calls `WinHelpA(hwnd, 0, 2, 0)` (null file, command 2) and is reached from `FUN_00471ee0` at `00471f92`; and `FUN_00577333`, a library virtual method held in 27 vtable slots with no direct caller.
- Other strings containing `help` belong to other surfaces: three `*_help_context` literals (`FUN_004c2353`) and `<Alt-h> This help` (`FUN_0053cdc0`). Neither is on the F1 path; their surfaces were not examined.
- Each root's `main.res` holds exactly one entry named `text/help.txt`.

**Confidence.** Medium: absence over the named population only. The library's standard `ID_HELP` command (`0xE145`) has dword hits at file offsets `0x0019c650` and `0x0019c654` and in code immediates; its producers and the virtual `WinHelp` method's callers were not traced, so a help route through an MFC modal-dialog hook is not excluded.

**Unknown.** Whether F1 inside a native modal dialog (the save and load chooser) reaches the library's `ID_HELP` handling, and what the `.hlp` literal at `0x0019d52c` is consumed by.

## Settings shortcut and speed notices

| ID | Claim | Confidence | Status | Evidence |
|---|---|---|---|---|
| MENU-057 | In a map session Ctrl+W, F, H, U, L, N and O each change one setting and post its new state's line, slot = base + state value, bases 94, 97, 100, 218, 102, 104, 106; W is retreat, F formation, L flying damage, O smoothing. | High | ● active | [EXP-0446](../experiments/EXP-0446-toggle-notice/) |
| MENU-058 | Numpad plus and minus without Ctrl, in a campaign session on the map screen, step the speed index by one, clamped to 0..8, and post slot 108 plus the clamped index through the duplicate-dropping post; Ctrl variants post nothing. | High | ● active | [EXP-0446](../experiments/EXP-0446-toggle-notice/) |
| MENU-059 | These notices are grey lines of 2000 ms appended below the existing lines of the map message line, the oldest removed past the capacity. Toggles are gated by the map screen word, the text entry and the Ctrl latch; speed also needs phase 2. | High / Medium | ● active | [EXP-0446](../experiments/EXP-0446-toggle-notice/) |

### MENU-057

- Dispatch: `FUN_0040f38e` is the map view's key-down handler (vtable slot at `0x005970f4`). A key outside its first jump table reaches `0040f67e`. With a nonempty selection and `view+0x144 & 0x24` clear, letters `A..T` go through the letter table at `0x0040fd6d`, where `F`, `H`, `L`, `N` and `O` reach `0040f789`; `W` lies outside that table and also reaches `0040f789`. A zero selection or a set `0x24` bit jumps straight to `0040f789`. The second dispatch there (key minus `0x20`, byte table `0x0040fddd`, target table `0x0040fda5`, all read from the PE by `tools/togglenotice`) sends vkeys `0x57`, `0x46`, `0x48`, `0x55`, `0x4c`, `0x4e` and `0x4f` to seven distinct arms. Selection state does not change which arm runs.
- Each arm opens with `CMP [0x005eb558],0` (the Ctrl latch, `AI-KEYMOD-059`) and skips everything when it is clear.

  | vkey | arm | state | change | slot | other call |
  |---|---|---|---|---|---|
  | `0x57` W | `0040f96b` | `view+0xa9c` | +1 mod 3 | 94 + state | `FUN_0041ccd2(state)` before the post |
  | `0x46` F | `0040f9ef` | `view+0xa98` | +1 mod 3 | 97 + state | `FUN_0041cd88(state)` before the post |
  | `0x48` H | `0040fa73` | `view+0xaa0` | 0 or 1 | 100 + state | none |
  | `0x55` U | `0040fad8` | `[0x005eb534]` | +1 mod 3 | 218 + state | `FUN_0041cd2d(state)` before the post |
  | `0x4c` L | `0040fb40` | `view+0xaa4` | 0 or 1 | 102 + state | none |
  | `0x4e` N | `0040fba5` | `[0x005eb528]` | new = (old == 0) | 104 + state | `FUN_00421a54(1)` after the post |
  | `0x4f` O | `0040fbf7` | `[0x005eb520]` | new = (old == 0) | 106 + state | none |

- Index: the post pushes `[[0x005eb3d4] + 4 * (base + state)]`, a line of the shared line array (`TEXT-STRTAB-023`), as `[array + state*4 + disp]` with `disp` `0x178`, `0x184`, `0x190`, `0x368`, `0x198`, `0x1a0` and `0x1a8`. The state read is the value just stored, so the line names the new state. Three-state settings run 0, 1, 2, 0. The arms hold no table of state names.
- `FUN_0041ccd2`, `FUN_0041cd88` and `FUN_0041cd2d` build record type `0x46` with subtype 1, 2 or 3 and the state at `+0xe`, with the word at `[view+0x9b4]+4` as player, and hand it to `FUN_004e74fe` on object `0x005f22d0`. The post is issued whatever that call does.
- The `keyboard.tsv` rows of `AI-KEY-125` for Ctrl+F, Ctrl+L, Ctrl+O and Ctrl+W name F retreat, W formation, L smoothing and O flying damage. The arms read the other way round. The arms agree with `ANIM-NUM-020`: `view+0xaa4` is flying damage and L toggles it. The H, N and U rows agree, and `MISSION-MSGPOST-058` and `ANIM-NUM-020` already carried the correct map; the cause of the `keyboard.tsv` error is not identified. `TERR-LIGHT-108` reads `[0x005eb528]` as the ShowTimeFlow option, which agrees with N as day/night.

**Confidence.** High: the vkey-to-arm mapping is read from the PE's index and target bytes (`key-arms.tsv`, with a test), and every arm, state field, displacement and call is a named instruction in the complete listing of `FUN_0040f38e` (624 instructions, 0 not disassembled). The setting names come from the text each slot holds, not from the state fields' other readers, which were not traced.

**Unknown.** What the receiver does with record `0x46` subtypes 1..3, and what `FUN_00421a54(1)` redraws.

### MENU-058

- Arms: the plus key (`0x6b`) at `00472d68` and the minus key (`0x6d`) at `00472dc8` of `FUN_00472b80`, the frame's key-down handler. With the Ctrl latch set the arm writes `campaign+0x40c` only (plus: 1; minus: 0 with a fresh stamp and cleared counters, `AI-KEY-125`) and posts nothing.
- Without Ctrl, `campaign+0x6bc == 2` (`00472d7f`, `00472df0`) and `campaign+0x3dc == 1` (`00472d8c`, `00472df9`) are both required, else the key ends there. Then `FUN_00477370(campaign+0x3f4 ± 1)` runs.
- `FUN_00477370` clamps its argument to 0..8 with a signed compare, stores it in `+0x3f4` and sets the rate from the ladder 8, 10, 12, 14, 16, 20, 24, 28, 32 (jump table `0x477424`).
- The post reads `+0x3f4` again, so at either end of the range the index is the clamped one. It pushes `[[0x005eb3d4] + 4 * (108 + index)]` (displacement `0x1b0`), ramp `0x5e9720` and lifetime `0x7d0`, and calls `FUN_00401fe0` on `[campaign+0xd0] + 0xa10`, the map view's list (`SESS-VIEW-028`).
- `FUN_00401fe0` drops a one-piece post equal to the newest line (`MISSION-MSGLINE-057`; compare re-read at `0040207a`..`00402089`). A second press at the end of the range, while the same line is still the newest, adds nothing; once the line has expired it is posted again.
- Slots 108..116 are the nine speeds in index order.

**Confidence.** High: both arms, the gates, the clamp and the index are named instructions in complete listings (`r-gates.txt`, `d-fns.txt`).

- Text entry: `FUN_00472b80` from `00472b80` to `00472bc9` and the arms `00472d68`..`00472e3a` hold no `campaign+0xc0` test. `AI-KEY-125` states that an open text entry consumes map shortcuts; no test for it was found in these bounds.

**Unknown.** Whether an earlier handler in the frame's message map consumes the key while a text entry is open.

### MENU-059

- Surface: the seven toggle arms call `FUN_00401e70` and the speed arms `FUN_00401fe0` on `view+0xa10`, the list `MISSION-MSGLINE-056` describes: in a campaign session drawn from (8, 8) in font 1 with a 1-pixel shadow, 17 pixels per line. The ramp `0x5e9720` is the grey one and the lifetime `0x7d0` is 2000 ms. The line is not a modal panel and has no rectangle of its own.
- Several notices: a post appends below the existing lines, a post past the capacity (14, 16 or 22 lines at 640, 800 or 1024 wide) removes the oldest, and lines expire oldest first, one per tick message, each lifetime counted from the removal of the line above it (`MISSION-MSGLINE-057`). Two toggles in a row therefore show two lines at once. `FUN_00401e70` never drops a duplicate; only the speed post does.
- Toggle gate: `campaign+0x3dc == 1` (`0040f3c8`, `0040f681`), the text entry `campaign+0xc0` closed (`0040f3d8`; an open one takes the key) and the Ctrl latch. On these arms `FUN_0040f38e` reads no `campaign+0x6bc`, player count, selection, ownership bit or option: the 624 instructions hold no `0x6bc` operand.
- Speed gate: phase `campaign+0x6bc == 2`, word `+0x3dc == 1` and no Ctrl (`MENU-058`).
- The frame arm `00472e41`..`00472e5a` forwards letter keys `A..Z` to `campaign+0xcc` with message `0x100`; where that route ends was not read, so the absence of other gates is bounded to `FUN_0040f38e` and the two speed arms.
- In phase 3 the same list is drawn at (0, 220) in font 2 (`MISSION-MSGLINE-056`); the toggle arms add no phase test.

**Confidence.** High for the surface, the lifetime, the ramp and the gates read from the listing. Medium for the dwell as wall-clock time and for the stacking timing: they inherit `MISSION-MSGLINE-057`'s Medium, since the interval of the `0x401` tick message was not measured.

**Unknown.** The `0x401` interval. Whether a phase-3 or networked session delivers these keys to `FUN_0040f38e` at run time: no original was run.

## F12, Backspace, Alt band and the Cast popup

| ID | Claim | Confidence | Status | Evidence |
|---|---|---|---|---|
| MENU-060 | In a map session F12 flips the static dword `0x005eb584`; while it is set `FUN_00407b1a` draws a box and a `%3.1f fps` line. No other instruction references the address. | High / Medium / Unknown | ● active | [EXP-0450](../experiments/EXP-0450-keyboard-remainder/) |
| MENU-061 | In a map session Backspace empties the map message line: `FUN_00401e40` on `view+0xa10` calls three routines with the arguments of `SetSize(0, -1)` on its text, colour and lifetime arrays and writes nothing else; the key returns 0. | High / Medium / Unknown | ● active | [EXP-0450](../experiments/EXP-0450-keyboard-remainder/) |
| MENU-062 | In a map session Alt plus a letter B..Y except S broadcasts a type `0x46` record with sub-selector `0x80` and index letter minus `A`; Alt+S is the screenshot. No local effect found in the handler. | High / Unknown | ● active | [EXP-0450](../experiments/EXP-0450-keyboard-remainder/) |
| MENU-063 | The map key handler's return 0 for C, Backspace and unmapped keys is dropped by the frame handler, which always calls the MFC default; 0 only lets the root offer the key to its other children and its own key routine. | High / Medium / Unknown | ● active | [EXP-0450](../experiments/EXP-0450-keyboard-remainder/) |
| MENU-064 | The Cast arm's `0x408` goes to the container `campaign+0xd4` and shows no popup; three child panels (Medium that they are all) write one field each. The spellbook popup is `campaign+0xec`, shown by `FUN_0041de5d`. | High / Medium / Unknown | ● active | [EXP-0450](../experiments/EXP-0450-keyboard-remainder/) |
| MENU-065 | In the spell-bar class the current slot `[+0x60]` is written by constructors, the `0x411`/`0x417` arms, the click routine and `-1` stores in `FUN_004b0fa0`/`FUN_004b1080`; `0x417` is pushed only by F5..F8; `FUN_0041b439(5)` is unread. | High / Medium / Unknown | ● active | [EXP-0450](../experiments/EXP-0450-keyboard-remainder/) |

### MENU-060

- Route: the frame key handler `FUN_00472b80` sends F12 (`0x7b`, byte `0x12` at `0x472ee8 + 0x73`, entry 18 of the 21-dword table at `0x472e94`) to `00472e4b`, which forwards message `0x100` to the root `campaign+0xcc` through `vt+0x48`. The map view's key handler `FUN_0040f38e` reaches its F12 arm only after `campaign+0x3dc == 1` (`0040f3c8`) and with the chat entry `campaign+0xc0` closed (`0040f3d8`; an open entry takes the key, `MENU-054`).
- Toggle: `0040f500` `XOR EDX,EDX` / `CMP [0x005eb584],0` / `SETZ DL` / `MOV [0x005eb584],EDX`, then `EAX = 1`. The cell lies past the raw `.data` bytes (image `.bss`), so it starts 0: the readout is off at load.
- Readers: `EnumRefs refto:5eb584` gives 4 hits, 2 owners, 0 orphan: the toggle (`0040f502`, `0040f50c`) and two reads in `FUN_00407b1a`, the map view's `0x402` frame routine of `MISSION-MSGLINE-056` (`0040cb45`, `0040cce6`). Among direct references the F12 arm is the only writer (a computed or indirect access is not excluded). The cell is in no save record found in this population.
- Effect (listing `rom-fps-draw.txt`): at `0040cb45` a set cell runs `FUN_0044f990` with five arguments (`[view+0x10]-0x78`, `0`, `[view+0x10]-0x1e`, `0x18`, a colour built from the pixel-format shifts at `0x5bcedc..0x5bcee4`), then `FUN_005545d0` (a `sprintf`) with the format `%3.1f fps` (`0x5b8274`) and the double at `view+0xd0`, then `FUN_00456b50` in font1 `[0x005e88b8]`. At `0040cce6` a set cell also calls `FUN_0044cda0` on a local rectangle. A clear cell skips both.
- The double: `view+0xd0` is stored by `FSTP` after `FILD` / `FDIVP` at `0040cb03`..`0040cb11`, `view+0xc8` is reduced by `0x3e8` and `view+0xc4` cleared (`0040cb1d`..`0040cb3b`): a recomputation once per 1000 view-tick units; the unit of `view+0xc8` and `view+0xc4` is unread, so wall-clock time is not claimed.
- Rejected: a gameplay switch (no direct reader outside the draw routine), a store with no reader, a grid or overlay (the only draw found is the number).

**Confidence.** High for the route, the toggle, the reader set and the format string (arm and tables read whole, string dumped from its own address). Medium for the box geometry (`[view+0x10]` is taken as the view width, not confirmed) and for the value being a frame rate: the label says so, the numerator and denominator of the divide were not read.

**Unknown.** Where the text lands inside the box, the unit of `view+0xc8` and `view+0xc4`, and the `FUN_0044cda0` rectangle's role. Whether a native run draws the box: no original was run.

### MENU-061

- Route: the frame handler's byte table sends Backspace (`0x08`, index 0) to `00472e4b`, the same forward as F12. In `FUN_0040f38e` the arm at `0040f66b` is `MOV ECX,[EBP-0x50]` / `ADD ECX,0xa10` / `CALL 0x00401e40` and then `JMP 0040fc47`, the exit of `MENU-054` with `EAX = 0`. Gate: `campaign+0x3dc == 1` and a closed chat entry.
- Routine: `FUN_00401e40` pushes `-1`, `0` and calls `0x0056fc7a` on `this+0x04`, `0x0057011d` on `this+0x18` and `0x00570490` on `this+0x2c`: the arguments of an MFC `CArray::SetSize(0, -1)` (count 0, default grow) on each array. It writes no other field: the redraw flag `+0x5c`, the tick times `+0x40`, `+0x44` and the capacity `+0x48` are untouched.
- Contents: the three arrays are the message line's text (`+0x04`), colour ramp (`+0x18`) and lifetime (`+0x2c`) of `MISSION-MSGLINE-056`, appended by `FUN_00401e70` and `FUN_00401fe0` (`MENU-059`). The next draw loops over count 0 and shows no line.
- Other callers (`EnumRefs callto:401e40`: 3 hits, 3 owners, 0 orphan): the message line's destructor `FUN_00401d10` (`00401d3b`) and `FUN_00472800` (`0047286d`, not read).
- Rejected: the chat text entry, the selection set, the order queue and the session record queue: the routine's only operand is `view+0xa10`.

**Confidence.** High for the route, the three operands and the absence of other writes (routine read whole). Medium for the three callees being `SetSize` of the three array classes: their bodies were not read, the argument pair and the layout of `MISSION-MSGLINE-056` carry the identification.

**Unknown.** What `FUN_00472800` does around `0047286d`. Whether the return of 0 lets a later sibling act on Backspace (`MENU-063`).

### MENU-062

- Route: `FUN_00472fb0` (the frame's `WM_SYSKEYDOWN` handler) tests the Alt context bit `lParam & 0x2000`, then in order: Alt+digit and Alt+numpad digit forward to the root (`AI-KEY-125`), `0x53` calls the screenshot routine `0x0044bdd0`, and a key in `0x42..0x59` with `campaign+0x3dc` bit 0 set calls `FUN_0041cde3(key - 0x41)` (`00473038`..`00473054`). Alt+A and Alt+Z fall out of the range test. After every path the handler calls the MFC default `0x0057638c`. `EnumRefs callto:41cde3`: 1 hit.
- Record: `FUN_0041cde3` fills the static record `0x609c38`: `+9 = 0x46`, `+5` = the word at `[[view+0x9b4]+4]` (the local player's id), `+7 = 0`, `+0xa = 0x80` (a dword), `+0xe` = the index; then `FUN_004e74fe` on `0x5f22d0` sends it. A zero word at `+7` is the broadcast addressee of `SESS-CMD-015`.
- Population: 23 keys (B..Y except S) post the same record shape; the index is the only difference. Receiver and effect: `AI-378`.
- Sender-side effect: none found in the handler, which does not touch the state word, the latches or the message line. `FUN_004e74fe` and the trailing call `FUN_0042b900` (`00473065`) are unread.

**Confidence.** High for the gate, the record fields and the single call site (handler and builder read whole, tables dumped).

**Unknown.** Whether a single-player session loops the record back to `FUN_004d5dd8`: the path from `FUN_004e74fe` to the drain loop was not re-read here (`SESS-CMD-015`).

### MENU-063

- Frame: `FUN_00472b80` forwards the key with `PUSH key` / `PUSH 0x100` / `CALL [EDX+0x48]` on `campaign+0xcc` (`00472e4b`..`00472e5a`) and never reads `EAX`: the next instruction at `00472e5d` is `MOV ECX,ESI` / `CALL 0x0057638c` (the MFC default). It returns `RET 0xc` with no value. Every arm of its table ends at `00472e5d`. A C key (`0x43`) takes byte `0x14` at `0x472ee8 + 0x3b`, entry 20, `00472e41`, which forwards letters `0x41..0x5a`: the letter route `MENU-059` left unread.
- Dispatcher: `FUN_004bd9dc` (root `vt+0x48`) for `0x100..0x102` calls the focus child `[this+0x38]` first when set, then the child walk `FUN_004bd87e`, then, when both returned 0, the object's own `vt+0x6c`, `+0x70` or `+0x74` (`004bdc6b`..`004bdc9f`); a nonzero result is returned at once. `FUN_004bd87e` stops on a nonzero child result and, for non-mouse messages, continues to the next child on 0 (`004bd9c0`..`004bd9ce`).
- Map view: `FUN_0040f38e` is the own `vt+0x6c` of the map view `campaign+0xd0`. Its return 0 reaches the root's walk as that child's result, so the root offers the key to its remaining children and then to its own routine before returning 0 to the frame, which discards it. For C this is the fall-through of `MENU-054`: no selection, a set `0x24` bit, a state word other than 1 or an open chat entry.
- Player-visible effect found: none. The key is neither consumed nor acted on by the frame beyond the default.

**Confidence.** High for the frame never using the return and always calling the default, and for the dispatcher's order (all three routines read whole). Medium that the other children do nothing with C or Backspace: their key-down routines were not read.

**Unknown.** The bodies of the root's own `vt+0x6c` and of the other children's key routines (the default stubs `FUN_00436e40`, `FUN_00436e50`), whether the focus child `[root+0x38]` is ever set in a map session, and the child order of the root's list at `+0x1c`.

### MENU-064

- Posting: `FUN_0041de5d` loads `campaign+0xec`, sizes it with `FUN_004bcba8` (a `FUN_004bd28f` child-3 test picks the height variant, and the overlay bit `0x2` adds a fixed `0x131,0x1e0,0x186` placement), attaches it with `FUN_004bccf4`, calls the sound routine `FUN_00453b08`, then `PUSH 0`, `PUSH 0`, `PUSH 0x408`, `CALL [vt+0x48]` on `campaign+0xd4` (`0041df7e`..`0041df97`) and `FUN_0041e06d`. The close routine `FUN_0041dfa6` posts the same `0x408` (`0041e045`).
- Receivers: `FUN_004bd9dc` on the container walks its children on 0. Three child routines act on `0x408` (jump entries of message minus `0x402`, read from the PE): `FUN_0048ec10` sets `[this+0x6c] = -100` (`0048ecc0`); `FUN_00491080` sets `[this+0x60] = 1` (`004911e1`); `FUN_00492c80` sets `[this+0x5c] = 1` (`00492d0e`). `FUN_00490380` has no `0x408` entry (`004904c2`, its default).
- Popup: `campaign+0xec` is built once by `FUN_004b02a0` (`FUN_00472070`, `004722fa`); it paints 24 cells from bits of `view+0x148` (`004b0a57`, `004b0ab4`), labels slots `F%d` (`0x5c0a44`) and highlights the cell equal to `[this+0x60]` (`FUN_004b08b0`); `FUN_004b0350` is its tooltip.
- `0x408` is a general refresh: 29 immediates in 20 routines (`EnumRefs imm:408`, `rom-enum-imm.txt`); the Cast open and close are two of them.
- Rejected: a `0x408` handler that creates or shows the popup (the three handlers write one field each).

**Confidence.** High for the posting chain, the three field writes and the popup's identity as the spell bar (listings read whole; tables read from the PE). Medium for the container's children being exactly these routines: the child slots were taken from the construction in `FUN_00472070` read earlier, and are not committed beside this card.

**Unknown.** What `+0x6c`, `+0x60` and `+0x5c` of the three panels mean (cursor, mode or redraw state). Whether `FUN_00453b08` plays a sound: it tests `[0x005e8430]`, walks an object list at `this+0x10`, and its arguments are `[0x005eb468]`, `0`, `0`, `0xdc`; the buffer chosen was not followed.

### MENU-065

- Writers of `[this+0x60]` in the spell-bar routines `FUN_004b0230`..`FUN_004b1080` (`EnumRefs disp:60`, `rom-enum-disp.txt`; a store through a pointer to `campaign+0xec` from another routine is outside this search): both constructors `FUN_004b0230` and `FUN_004b02a0` store `-1` (`004b024d`, `004b02da`); `FUN_004b0e40` stores `-1` for message `0x411` (`004b0f24`) and a slot value for `0x417` (`004b0f12`); `FUN_004b0f40`, the click routine, stores the hit cell when its bit in `view+0x148` is set (`004b0f84`); `FUN_004b0fa0` and `FUN_004b1080` store `-1` (`004b0fa0`, `004b1080`).
- Senders: `EnumRefs imm:411` finds one push, `FUN_00472070` at `00472309` (construction); `imm:417` finds four, the F5..F8 arms of `FUN_0040f38e` (`0040f524`, `0040f578`, `0040f5cc`, `0040f620`; census `rom-enum-imm.txt`). The C path's own routines `FUN_0040f38e` (arm `0040f755`), `FUN_0041de5d` and `FUN_0041e06d` contain neither. `FUN_0041b439(5)` is not re-read in this experiment (`MENU-055` read it).
- `FUN_0041e06d` writes `view+0x68`, clamps `view+0x60` (the view's field, not the popup's) and sets `view+0x74 = 1`.
- So C leaves the popup's current spell as it was: `-1` until a quick-slot key, a click or an earlier assignment sets it (`AI-SPELLGUARD-289`, `AI-PANEL-123`).

**Confidence.** High that the C path's routines read here select nothing (bounded to immediate pushes of `0x411` and `0x417` and to displacement stores to `+0x60` in the class's routines; `FUN_0041b439` is not re-read). Medium that no computed message id reaches `FUN_004b0e40`.

**Unknown.** Whether closing the popup resets the current spell: `FUN_0041dfa6` contains no `0x411` push. The message arms of `FUN_004b0fa0` and `FUN_004b1080`. Whether the item-Cast route of `AI-PANEL-123` assigns a slot before the popup shows.

## Hover help bindings and rectangles

| ID | Claim | Confidence | Status | Evidence |
|---|---|---|---|---|
| MENU-066 | Attribute-row hover rectangles are panel-relative: value x 82..102, cost 107..127, refund 132..152, y 54+32i..74+32i for rows i 0..3; one shared pool rectangle (46,181,123,203); the gate dword is 1 only while the owner screen runs. | High / Medium | ● active | [EXP-0454](../experiments/EXP-0454-hover-help/) |
| MENU-067 | In a map-list record `+0x14` and `+0x18` are the map file's first two header dwords, its width and height; the list row prints each minus 16 and the `.alm` loader seeks 0x28 bytes after reading them. | High / Medium | ● active | [EXP-0454](../experiments/EXP-0454-hover-help/) |
| MENU-068 | `dialogs.txt` 23 and 75 reach four Sound Options controls through the vslot 6 setter, `dialogs.txt` 117 reaches a text edit box through its constructor, and `patch.txt` 52..54 reach three Game Options controls as hint and caption. | High / Medium | ● active | [EXP-0454](../experiments/EXP-0454-hover-help/) |
| MENU-069 | Of the 52 inherited-getter widget tables not reached by the traced hint sites, 46 are built with a null hint text and 6 differ; of the 9 of 252 vslot-6 receivers that resolve, none is one of the 52. | High / Unknown | ● active | [EXP-0454](../experiments/EXP-0454-hover-help/) |

### MENU-066

The initialiser `0042bca0` writes the rectangles into the panel (vtable `0x00597550`) in panel coordinates; the getter `0042ca00` copies each, offsets it by the panel origin (`004bcc23`: the panel's `(+8, +0xc)` plus the same pair of every ancestor on the `+0x30` chain, added by `004c7c90`) and tests it with `PtInRect`, which is half-open: left <= x < right and top <= y < bottom. Tests run per row in the order label, pool, value, cost, refund (TEXT-082): five tests per row over 17 distinct rectangles (4 label, 12 value, cost and refund, 1 pool).

| Rectangle | Panel offset | left | top | right | bottom |
|---|---|---|---|---|---|
| value, row i | `+0xb4 + 0x30 i` | 82 | 54 + 32 i | 102 | 74 + 32 i |
| cost, row i | `+0xc4 + 0x30 i` | 107 | 54 + 32 i | 127 | 74 + 32 i |
| refund, row i | `+0xd4 + 0x30 i` | 132 | 54 + 32 i | 152 | 74 + 32 i |
| label, row i | `+0x64 + 0x10 i` | 16 | 57 + 33 i | 16 + W | 57 + 33 i + H |
| pool | `+0xa4` | 46 | 181 | 123 | 203 |

Rows i are 0..3. `W` is the glyph-advance sum of the row label on the font object `[0x5e9bc0]` (`00456320`) and `H` the vslot `+0x24` result of that object, both runtime values. The pool rectangle is one record tested on every row, so TEXT-082's second rectangle is the same on all rows; its y range lies below the value, cost and refund rectangles (last bottom 170).

The gate `[[view+0x5c]+0x104]` is a dword of the panel's owner object `A`, which the panel constructor `0042bbb0` receives as its last stack argument from `A`'s build routine `0042f330`. Within `A`'s methods three stores write it: 0 in the build routine (`0042f4c2`), 1 at the end of the start routine (`0042f9be`), 0 in the teardown routine (`0042fa5d`). Readers through `[this+0x5c]` are `0042ca00`, `0042f0b0` and `0042f120`. By these stores the hover answers only between the start and teardown of that screen.

**Confidence.** **High** for the rectangle immediates, strides, order and the three stores (instruction rows reproduced from both roots). **Medium** for the owner role and for "only between start and teardown": a whole-image scan for operands `0x104` (`scan104.txt`, 131 rows, displacement or immediate) was matched to `A`'s methods by owner routine, and stores in routines of other classes were taken to be other objects, not traced.

**Unknown.** `W` and `H` in pixels (font data not read), and so whether a label rectangle (bottom 156 + H for the last row) reaches the pool's top 181. The screen `A` implements, and whether anything besides start and teardown changes the dword. The panel origin's absolute value.

### MENU-067

The loader `00447211` reads 4 bytes into record `+0x14` (`00447377`, `00447383`) and 4 bytes into `+0x18` (`0044738b`, `00447397`), then seeks 0x28 (`0044739c`). These are the type-0 payload `+0` and `+4` dwords that ALM-META-091 stores as `M+0` and `M+4`. The row text routine pushes `+0x18 - 0x10` and `+0x14 - 0x10` (`0044591d`, `00445927`) for the `%dx%d` columns of TEXT-083, so the shown size is the raw size minus 16. The record constructor `004482f0` initialises only the three CStrings at `+0`, `+4` and `+8` (allocation 0x1c); the loader fills the numbers. The shown size equals the playable rectangle `(8, 8, W-9, H-9)` of TERR-SIGHT-116, which has width `W - 16`.

**Confidence.** **High** that the two fields are the raw dwords and that the display subtracts 16 (named instructions). **Medium** that the 16 is the 8-cell border on each side: the match with TERR-SIGHT-116 is arithmetic, and no instruction ties the two.

**Unknown.** Whether the first two header dwords are width and height in every shipped map (the loader does not name them).

### MENU-068

- `dialogs.txt` 23 and 75, builder `00438e87` (the Sound Options dialog, VIDEO-SFX-053): four setter pairs, each `PUSH 0x17` then `CALL [reg+0x18]` and `PUSH 0x4b` then `CALL [reg+0x18]` (`00439134`/`00439174`, `0043959e`/`004395de`, `004396b9`/`004396f9`, `00439978`/`004399b8`). The receivers, by the construction before each pair: a control of constructor `00448cb0` (`004390d4`, class `00598d18`), two of `004be6e3` (`00439538`, `00439653`; class `0059b308`) and one of `004c35cd` (`004398aa`; class `0059b598`). All three classes carry the inherited getter `004bd631` and their constructors take a hint at construction (`dialogs.txt` 10, 13, 15 and 19 respectively for the four controls, the first-pushed argument), so the setter replaces it later.
- `dialogs.txt` 117: pushed at `00444542` as the last argument of constructor `004bf28e` (`0044457e`, class `0059b380`, a text edit box), in the routine at `00444507`. It is a constructor hint.
- `patch.txt` 52, 53 and 54 (object `0x5ea668`), builder `0043d847` (the Game Options dialog, TOWN-OPTIONS-457): each is pushed as the last constructor argument of a `00448cb0` control (`0043dc18`/`0043dc6c`, `0043dd09`/`0043dd60`, `0043ddfd`/`0043de57`) and again to caption setter `00448c90` of the same control (`0043dc9f`..`0043dcaf` and the two following ones), which forwards to `this+0x64`. The same dialog also reads `dialogs.txt` 52, 53 and 54 for other controls.

**Confidence.** **High** for the push, constructor and setter addresses. **Medium** that the setter receivers are exactly those four controls: they come from the nearest constructor call before each pair, not a data-flow proof; **Medium** that the last constructor argument of `00448cb0`, `004be6e3` and `004bf28e` is the hint text, which follows the TEXT-085 trace.

**Unknown.** The other 243 vslot-6 receivers (a backward trace resolves 9 of 252; one resolved site is a constructor, `00518160`, with no widget table). What `00448c90`'s `this+0x64` object draws.

### MENU-069

The `00597ac0`..`0059ba18` population is the 69 widget tables whose `+0x14` getter is `004bd631` (TEXT-HOVERSET-049); TEXT-085 reaches 17 by the hint-forwarding sites and the other 52 are classified here from the constructor that stores each table's vtable address, followed through its base-constructor chain to the text parameter of `004bc98f`, `004bc880` or `004bc770`.

| Traced text of the 52 | Tables |
|---|---|
| null constant | 46 |
| no constructor in the code map | 2 (`00597c60`, `005983d8`) |
| forwards a parameter, base of a reached class | 2 (`0059b910`, `0059b9a0`) |
| static address plus a parameter | 1 (`0059b3f8`, base of two reached classes) |
| null and static in different constructors | 1 (`0059b688`) |

Null text means the base constructor assigns no hint and the hint is empty until a setter writes it; hover raises only a nonempty string (TEXT-HOVER-048). By direct callers of the storing constructors, 44 of the 52 have a builder caller, 5 only a derived constructor and 3 none in the code map. Of the 9 vslot-6 sites whose receiver resolves (of 252), 8 target reached classes (MENU-068) and 1 is a constructor with no widget table (`00518160`); none is one of the 52. The other 243 are unresolved.

**Confidence.** **High** for the table classification under the stated trace and for the 9 resolved receivers. **Unknown** whether any of the 52 shows a hint: a hint can arrive through any of the 243 unresolved vslot-6 receivers or a derived constructor passing text the trace read as null.

**Unknown.** Which of the 52 controls appear in play, and which of the 243 unresolved vslot-6 sites are widget setters. The three `+0x3c` writers outside the base constructors: `00429f80` constructs class `0x5974e0` (not one of the 96 tables) and `005168e0` class `0x59be78`; the receiver class of the one in `004df385` is unresolved.
