# MENU — menu surface composition contracts — specification

Level 3. Promoted, evidence-backed claims only. Inventory assets are `MENU-ASSET-001`…`002`,
the hit mask and state are `MENU-MASK-003`…`MENU-STATE-007`, and the two in-play Esc menus are
`MENU-ESC-010` (partially retracted)…`MENU-INPUT-016`. This is a **composition contract**, not a bitmap codec — the images are
standard Windows BMP; what is re-derived here is how the engine assembles them into
an interactive surface.

Two different surface kinds are specified here. The pre-game main menu is a
full-screen bitmap surface with a pixel hit mask. The two menus Esc raises during
play carry no art of their own and are built from text rows over a sprite
nine-patch; they are specified in their own section at the end.

Inputs pinned: `main.res` (the bitmaps) and `rom.exe`
(sha256 `942e9b72610eeba2f3d74930ee47cbc4476b85a942c14348f874ec7b7d367d03`; the loader,
hit-test, and coordinate tables). All from the owned install.

## Asset set (18 files under `graphics/mainmenu/`)

| Entry | Dims | Bpp | Role |
|-------|------|-----|------|
| `menu_.bmp` | 640×480 | 24 | base brooch background (drawn first) |
| `menumask.bmp` | 640×480 | 8 | hit mask — palette **indices**, not colour |
| `button1..8.bmp` | see below | 24 | per-button **hover/highlight** overlay |
| `button1..8p.bmp` | see below | 24 | per-button **pressed** overlay |

The engine (`rom.exe` `FUN_00485e00`) loads them as
`main\graphics\MainMenu\{menu_.bmp, MenuMask.bmp, button%d.bmp, button%dp.bmp}`, with the
button loop running `n = 1..8` — so there are exactly **8 buttons**, each with a hover and
a pressed variant (MENU-ASSET-001/002). The mask is loaded by a distinct 8-bit loader
(`FUN_00429ab0`) that keeps raw indices; the colour bitmaps by `FUN_00429330`.

## Hit mask (`menumask.bmp`)

- 640×480, 8bpp. The palette is the **identity grayscale ramp** `palette[i]=(i,i,i)`, so
  the palette is decorative — the **8-bit index value** is the semantic (MENU-MASK-003).
- The hit-test `FUN_00486260` reads the index at the cursor as
  `idx = maskBuf[(y − top)·640 + (x − left)]` (menu origin `(left,top) = (0,0)`), then:

  | index | button | index | button |
  |------:|:------:|------:|:------:|
  | 0x80 | 1 | 0xc0 | 5 |
  | 0x90 | 2 | 0xd0 | 6 |
  | 0xa0 | 3 | 0xe0 | 7 |
  | 0xb0 | 4 | 0xf0 | 8 |

  Asset number = `idx/16 − 7`. Index 0 (background) and **every** other value — including
  the anti-alias edge ramp `0x10,0x12,…,0x1e` (= hot ≫ 3) and stray pixels — select **no
  button** (switch default) (MENU-MASK-004). Of the 43 index values present in the shipped
  mask, only these 8 are "hot".

## Overlay placement (static `rom.exe` tables)

Two contiguous tables, 8 entries × 16 B = `{x, y, w, h}` (int32 LE):

- **normal/hover** — `DAT_0059a2d8` (file offset `0x198ed8`)
- **pressed** — `DAT_0059a358` (immediately after, `+0x80`)

`FUN_00485c20` places button *i* at screen rect `(x, y) … (x+w, y+h)`, offset by the menu
origin (0,0). Overlays draw **1:1** — each `(w,h)` equals the BMP's pixel dimensions,
verified 16/16; and each normal rect brackets that button's mask region, 8/8
(MENU-GEOM-005/006).

| btn | hover (x,y,w,h) | pressed (x,y,w,h) |
|----:|-----------------|-------------------|
| 1 | 112, 64, 212, 136 | 116, 64, 208, 138 |
| 2 | 84, 88, 236, 148 | 88, 88, 236, 152 |
| 3 | 84, 236, 236, 152 | 88, 236, 236, 152 |
| 4 | 112, 276, 212, 136 | 116, 272, 208, 140 |
| 5 | 320, 60, 212, 140 | 320, 64, 212, 140 |
| 6 | 324, 88, 236, 148 | 324, 88, 232, 152 |
| 7 | 324, 236, 236, 152 | 324, 236, 236, 152 |
| 8 | 324, 276, 208, 136 | 320, 272, 212, 140 |

Layout: two columns of 4 (buttons 1–4 left, 5–8 right), each column top→bottom.

## Interaction / draw state

- Object fields (menu `this`): `+0x6c` hover-overlay ptr array, `+0x80` pressed-overlay
  ptr array, `+0x94` hover placement RECTs, `+0xa8` pressed placement RECTs, `+0xb8`
  background, `+0xbc` mask, `+0xd8` hovered index, `+0xdc` shown index, `+0xe0`
  pressed-latch index, `+0xe4` per-button disable bitfield.
- Selection (`FUN_00486260`): while not pressed, the **hover** overlay of the hovered
  button is composited over the base at its normal rect; on mouse-down over the latched
  button, its **pressed** overlay is composited at the pressed rect. A set bit in the
  `+0xe4` disable field suppresses the overlay for that button (MENU-STATE-007).
- Click (`FUN_00486530`): posts a per-button command message; button 8 → WM_CLOSE (quit).
  The other message ids are not decoded (Medium).

## The in-mission command panel (`MENU-COMBAT-017`…`MENU-COMBAT-019`)

The mission root owns a map viewport and a fixed 160-pixel right column. The right column owns,
in draw and hit-test order, the minimap, command panel, character panel and a lower filler:

| surface | local rectangle | 640x480 screen rectangle | 800x600 | 1024x768 |
|---|---|---|---|---|
| minimap | `(0,0,160,158)` | `(480,0,640,158)` | `(640,0,800,158)` | `(864,0,1024,158)` |
| command | `(0,158,160,238)` | `(480,158,640,238)` | `(640,158,800,238)` | `(864,158,1024,238)` |
| character | `(0,238,160,480)` | `(480,238,640,480)` | `(640,238,800,480)` | `(864,238,1024,480)` |
| filler | `(0,480,160,H)` | empty | `(640,480,800,600)` | `(864,480,1024,768)` |

The command panel has no child controls. Its eight hit cells are 34x34, row-major, at
`left=8+34*(i&3)`, `top=7+34*(i>>2)`. Right and bottom are exclusive, as in `PtInRect`.
Inventory, spellbook, portrait and Esc controls are in the character-panel sibling and must not
be attached to this object. That sibling is not display-only: an admitted inventory-item release
can arm Cast mode 5, and its transfer branches can produce item orders `0x22/0x32`
(MENU-COMBAT-017, AI-PANEL-123). Its nested inventory grid also transfers the selected item on
physical left-up through vslot `+0xa4`. Destination code 2 accepts a valid hit as an insertion index
and uses merge/insert-before-gold/append fallback for `-1` or an out-of-range hit.
The grid also owns two left-down edge strips that scroll its viewport, and left-button move starts
an item or money drag before that release. Its three right-button edges are consumed no-ops. None
of these inventory controls is physically part of the command-panel object.

| i | rect | tooltip / key | action |
|---:|---|---|---|
| 0 | `(8,7)-(42,41)` | Attack / A | arm mode 1 |
| 1 | `(42,7)-(76,41)` | Move / M | arm mode 2 |
| 2 | `(76,7)-(110,41)` | Guard / G | immediate order `0x17` |
| 3 | `(110,7)-(144,41)` | Defend / D | arm mode 4 |
| 4 | `(8,41)-(42,75)` | Cast / C | arm mode 5; suppressed while cast UI is open |
| 5 | `(42,41)-(76,75)` | Swarm / S | arm mode 6 |
| 6 | `(76,41)-(110,75)` | Stand Ground / T | immediate order `0x18` |
| 7 | `(110,41)-(144,75)` | Retreat / R | immediate order `0x14` |

Tooltip source is `main.res::text/main.txt[i]`. Russian labels, in the same order, are
`Атаковать`, `Идти`, `Охранять`, `Защищать`, `Колдовать`,
`Идти в боевой готовности`, `Держать позицию`, `Отступить`; the accelerator letters remain
Latin. A tooltip appears only while the panel and cell are enabled, no special cursor or child
dialog is active, and screen-state bits `0x0a` are clear (MENU-COMBAT-019).

The draw is four 160x80, 24-bpp BMPs. Inactive: `HeadsR.bmp`. Active: draw `CommandBarR.bmp`,
copy the same 34x34 cell from `CommandEmpR.bmp` over every disabled cell, then copy the selected
cell from `CommandDnR.bmp`. Icons and borders are embedded in these full-panel images; there is
no separate hover image. Left down acts immediately, double-click aliases it, and left-drag
repeats it as the pointer crosses enabled cells. Right up sends map cancel `0x405`; the other
right edges are no-ops (MENU-COMBAT-018/019).

**Customisation boundary.** Geometry, cell order and action binding are executable constants.
Labels and every visible pixel are resource data. Enlarging the panel or adding a ninth cell
requires code and new art; replacing a label or embedded icon changes only resource bytes.

## The Drop Gold modal (`MENU-COMBAT-017`, `AI-PANEL-123`)

The mission root always constructs a 296x168 dialog at `(100,H-200)-(396,H-32)` and stores it at
`campaign+0x10c`. It is attached only when physical `WM_LBUTTONDBLCLK` on the nested inventory
grid hits the sentinel gold cell (`item+6 == 0xffff`) while the local purse is positive.

Its children, in construction order, are an amount edit id `0x989685` at local
`(30,65)-(266,85)`, action button id `0x989681` at `(68,118)-(138,138)` posting `0x445`, cancel
id `0x989682` at `(158,118)-(228,138)` posting `0x446`, and two text rows at
`(20,20)-(276,40)` and `(20,40)-(276,60)`. Button labels are
`dialogs.txt[0..1]`: EN `Ok` / `Cancel`, RU `Принять` / `Отменить`; auxiliary strings
`[46..49]` supply its auxiliary strings. The modal input contract is `AI-PANEL-123`: buttons arm
on left down and post on an inside left up.
Enter aliases action; Esc aliases cancel.

Action parses the edit, resolves the entered amount, subtracts it from the purse and emits opcode
`0x23` at the current map cell before closing. Cancel closes without an order. Its screen placement
and child geometry are executable constants; all six strings are localised `main.res` data. The
generic dialog-frame drawing assets were not re-enumerated as part of the command-panel asset
census.

## The in-play Esc menus (`MENU-ESC-010` (partially retracted)…`MENU-INPUT-016`)

Two surfaces, not one. `VK_ESCAPE` reaches the frame window's `WM_KEYDOWN` handler
`FUN_00472b80`, which reads the UI state word `campaign+0x3dc`:

| state | meaning | posts | constructor | screen rect | rows |
|---|---|---|---|---|---|
| `== 1` | map session on screen | `0x416` | `FUN_0043c23b` | `(100,60)-(440,400)`, 340×340 | 7 of 8 built |
| `== 0` | town | `0x41f` | `FUN_0043c90b` | `(100,100)-(440,340)`, 340×240 | 5 |
| other | — | nothing | — | key forwarded to `campaign+0xcc` | — |

The town arm additionally requires the `campaign+0x3b4` CString to be empty. The
mission surface is also raised by a 32×32 button on the character-panel sibling at panel-local
`(0x7e,0xce)-(0x9e,0xee)`, whose tooltip is global string index 14 (MENU-ESC-010, partially
retracted: this row is the owning-panel correction).

Entries, in screen order. Labels come from `main.res::text/dialogs.txt` through the
descriptor at `0x005ea678`; `FUN_004687f0(table,i)` resolves to
`[0x005eb3d4][table->+0xc + i]` (MENU-ITEM-011, MENU-ITEM-012):

| row | mission label (index) | msg | town label (index) | msg |
|----:|---|---|---|---|
| 1 | `~Save Game` (0x22) | `0x41a` | `~Save Game` (0x22) | `0x41a` |
| 2 | `~Load Game` (0x23) **or** `Diplomacy` (0x4c) | `0x418` / `0x43c` | `~Load Game` (0x23) | `0x418` |
| 3 | `Game ~Options` (0x24) | `0x41b` | `Sou~nd Options` (0x25) | `0x422` |
| 4 | `Sou~nd Options` (0x25) | `0x422` | `Abort Game` (0x4d) | `0x41c` |
| 5 | `~Quest Objectives` (0x26) | `0x420` | `~Return to Game` (0x28) | `0x446` |
| 6 | `~End Quest` (0x27) | `0x41c` | | |
| 7 | `~Return to Game` (0x28) | `0x446` | | |

Row 2 of the mission menu is `~Load Game` when `campaign+0x6bc == 2` and `Diplomacy`
otherwise; only one of the two is built. `0x446` is the family-wide close message,
not a distinct action. The first argument of the button constructor is a **control
id, not a row index**.

Disabled at construction — mission menu only; every town row is built enabled:

| row | disabled when |
|---|---|
| `~Save Game` | `campaign+0x6bc` is 0 or 1 |
| `~Load Game` | no file matches `game*.sav` in the game directory (`FUN_0043c0ae`) |
| `Sou~nd Options` | `FUN_00449d90()` returns 5, its `this+0x9c == 0` arm |
| `~Quest Objectives` | `campaign+0x6bc != 2` |

**Accelerator.** The constructor passes an immediate, then `FUN_004be6e3` overwrites
it with the character following a single `~` in the label, lowercased by
`FUN_00468ac0` (CP866-aware). `~~` is an escape. Labels without a `~` keep the
immediate, which on the English release is `Diplomacy` → `D` and `Abort Game` → `E`;
on the Russian release both labels carry a `~` and the accelerator is a different,
language-dependent letter. Hard-coding the accelerators reproduces one release only
(MENU-KEY-013).

**The shared message is not a shared action.** Both surfaces post `0x41c`, and
`FUN_00473110` branches on `campaign+0x3dc`: `== 1` raises the five-row confirmation
`FUN_0043c66f` (`Change Map` or `~Victory!` at `0x41d`, `~Exit to Main Menu` and
`Exit to ~Windows` both at `0x41e`, `~Return to Game`), `== 0` raises the three-row
`FUN_0043cb6b` (the same two exit rows and `~Return to Game`), both at
`(100,100)-(440,340)` (MENU-ITEM-012).

**Geometry.** Each row is a panel-local rect `(40, 40+30(n−1), width−48, 40+30n)`:
252 px wide, 30 px tall, the first starting 40 px below the panel top
(`FUN_0043c129` with step `0x1e`, `FUN_0044a360` forcing left `0x28` and right
`width − 0x30`).

**Art.** Neither surface loads a bitmap: the family's art-load slot is `RET`. The
frame is a nine-patch drawn by `FUN_004c4da2` out of `graphics.res::interface/lm.256`
(one writer of `[0x005ef980]`, one reader). Frames 0..8: centre 96×64, corners 48×48,
top and bottom edges 96×48, left and right edges 48×64. Inset 48, horizontal step 96,
vertical step 64. Two passes: a drop shadow of the right column and bottom row offset
`+8,+8`, then the full nine-patch at the panel rect shrunk by 8 on right and bottom.
Frames 9..17 are a second nine-patch at 32/48 px, unused by this routine
(MENU-ART-014).

**Compositing and modality** are not part of this contract and are specified by the
claims: the whole screen is darkened once, destructively, at shade level 3
(`MENU-STOP-015`, `DLG-DIM-013`); the world stops on the `0x4008` idle gate
(`MENU-STOP-015`, `DLG-STOP-012`); Esc also closes, and the panel does not capture
the mouse (`MENU-INPUT-016`).

## Open / not established

- The `+0xe4` disable bits' source (which buttons start disabled) and the non-WM_CLOSE
  command-message meanings — needs further `rom.exe`/runtime.
- The DIB loader's row order (top-down vs flip) is inferred from the 8/8 placement bracket
  rather than read directly; the screen y here is top-down.
- How the Esc menus' nine-patch covers a frame width that is not `96 + 96k`. The tile
  counts are `(w−96)/96` and `(h−96)/64` with truncating division, which for the
  332×332 frame these panels have leaves a 44-px band inside the right edge and one
  inside the bottom edge unaccounted for by any placed tile. `MENU-ART-014` retains this as a
  falsifiable prediction rather than a settled fact.
- What each entry's raised panel then does. `MENU-ITEM-011` and `MENU-ITEM-012` name the constructor and rect each
  message raises and stops there; the two exit rows of both `0x41c` confirmations post the same
  `0x41e` and differ only in control id, and how the exit target is distinguished is unread.
- Whether a click landing outside an Esc panel changes anything, as opposed to merely being
  delivered to the window under the cursor, which is all `MENU-INPUT-016` claims.
