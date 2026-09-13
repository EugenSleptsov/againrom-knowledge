# Pointer-hover help

## Shared timing and dispatch

ROM1 has a common cursor controller which polls the current pointer, accumulates
elapsed `timeGetTime` milliseconds, and requests text from the widget beneath
the pointer. Its normal hover request occurs when elapsed time first crosses
500 ms. The test requires `new >= 500`, `new < 25500` and unsigned `old < 500`.
It is a threshold crossing, not a repeated request on every subsequent frame.
The actual display occurs at the next admitted controller poll, so 500 ms is
the threshold rather than a measured native frame timestamp. — TEXT-HOVER-048

The controller walks children in reverse stored order, recursively, then tests
the current widget's screen rectangle. It calls the selected widget's text
getter at virtual offset `+0x14`, copies the returned string and displays a
nonempty result. A null/empty child result does not make the dispatcher ask the
parent for another string. The normal path requires a root widget, no active
marquee and the controller's help-enable flag. Getter-specific state checks
still apply. — TEXT-HOVER-048, TEXT-HOVERSET-049

An admitted pointer-position change hides help, clears elapsed time and records
a new clock sample. The observed polling branch does not do this while the
right mouse button is held. At accumulated elapsed time 25500 ms the visible
box is cleared, giving about 25 seconds after the normal reveal. The separate
hide routine, also called after mouse button messages `0x201..0x206`, clears
visibility without restarting elapsed time; a hidden box
does not automatically reappear at an unchanged pointer position. Entering a
marquee also hides it. These are the identified controller routes, not a claim
about every focus or modal transition. — TEXT-HOVER-048

## Text targets and sources

Indices below are zero-based. `main.txt` is `main.res::text/main.txt`, whose
local indices equal global indices. Other files use their table-local index.
Source availability and a getter branch establish a static route; native
visibility under every state is not implied. — TEXT-HOVERSET-049,
TEXT-HOVERTEXT-052

| Surface | Hover target | Source |
|---|---|---|
| Command panel | Eight command icons | `main.txt[0..7]` |
| Character panel | Open/close book and backpack; doll/statistics; menu | `main.txt[8..14]`, chosen by current state |
| Character navigation | Previous/next hero or portrait | `main.txt[52..53]` or `[121..122]`, by mode |
| Statistics card | Four attributes, health, mana, damage, attack, armour, defence, weight, sight, speed, skills, resistances and experience | `main.txt[155..180]`; monster weapon resistance `[188]`; stat-specific visibility/ownership gates |
| Statistics card | Monster's assigned spell list | `main.txt[192]`, joined with names from `spell.txt` and runtime spell bits |
| Precreation | Difficulty, four hero templates, continue/back and name field | `main.txt[247..256]` |
| Final character generator | Four attributes, free points and ten class-specific skill pictures | `main.txt[155..158]`, `[273]`, `[171..180]`; numeric controls also compose their label and value |
| Town | Shop, school, tavern, mission exit and menu hotspots | `main.txt[233..237]`, in hotspot order rather than index order |
| School | Current five weapon or magic skill icons | `main.txt[171..180]`, with class and visual-slot permutation `1,2,4,3,5` |
| Tavern | Candidate statistics and equipped items | Shared statistics getter and item formatter |
| Shop | Shopkeeper and four stock groups | `main.txt[61..65]` |
| Shop and inventory | Scroll arrows; backpack, transaction table and shelf backgrounds | `main.txt[54..60]` |
| Shop, backpack and equipment | Item under the pointer, or gold | Item formatter; `main.txt[74]` for the inventory gold sentinel |
| Spellbook | Available spell cell | `spells.txt` name plus `main.txt` field labels and live costs/strength/range/duration/other present values |
| World map | Available site marker | `sites.txt[marker index]`, subject to discovery/mission availability |
| Map-selection list | Map row and three metadata columns | Row-owned text; `dialogs.txt[134..136]` |

Character and generator mappings are established by TEXT-HOVERCHAR-050.
Room, navigation and item-area mappings are established by TEXT-HOVERROOM-051.
The composed text and local-table branches are established by TEXT-HOVERTEXT-052.
Existing TEXT-UI-039 and MENU-COMBAT-019 still apply to their own narrower rows.

The inspected family comprises 96 aligned widget tables carrying the same
three base methods: 28 distinct `+0x14` getters, of which 21 contain specialized
text-return paths, six return zero, and one reads the widget's owned hint
string at `+0x3c` for 69 tables. That inherited getter suppresses its string
when widget flag `0x20` is set; the setter at `+0x18` copies its argument.
All assignments to these inherited hint fields, computed replacement tables
and control families with different base methods remain outside the inventory.
The six null getters do not prove that their entire screens lack help.
— TEXT-HOVERSET-049

## Presentation

The string's `#` bytes are line separators. The normal layout measures each
source line, keeps the maximum width and places the box above the pointer:
width is measured width plus 11 pixels, height is `14 * lineCount + 5`.
It shifts left at the right screen edge and down at the top screen bound.
Text uses `font2`, with a five-pixel left inset, four-pixel top inset and
14-pixel line pitch. The decoded path splits authored lines rather than
automatically wrapping them by words. Long descriptions already carry these
separators in both installed languages. — TEXT-HOVERPAINT-053

This lifecycle is separate from introductory `text/tips/*.txt` panels, whose
constructor, controls and persistent `TipsMode` gate are documented in the town
claims. Disabling an introductory panel is not evidence that the hover
controller's distinct enable field is disabled. — TEXT-HOVERPAINT-053,
TOWN-184, TOWN-185, TOWN-186

## Limits

The timings, branches and resources above come from static executable/data
inspection and paired resource measurements. Native frame timing, all focus
transitions, all inherited hint assignments and a universal inventory of every
possible control remain unproved. A title or tooltip-shaped string in a table
alone does not establish a control binding. — TEXT-HOVERSET-049
