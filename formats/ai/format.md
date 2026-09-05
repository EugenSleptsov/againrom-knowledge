# AI — how a unit chooses whom to fight, and what it does when nobody is ordering it (partial)

Level 3. Promoted, evidence-backed claims only. Core acquisition is `AI-FILTER-001`…`AI-GUARD-007`;
the authored behaviour layer is `AI-TICK-008`…`AI-DIFF-016`. Ledger:
[`claims/ai.md`](../../claims/ai.md).

**Status: partial (◐).** Specified end to end: the relation a unit consults, where it comes
from, the candidate population, the selection, the three radii, **who runs the AI and how often,
the group and per-actor state machines, guard, patrol, and where a behaviour is chosen**. Not
specified: how the resulting order is *executed* (that is `formats/move`), the four player-issued
group orders, the candidate scorer, and the line-of-sight predicate itself. See the ledger's
*Open, written out*.

## Saved Group continuation

Group LOAD restores the80-byte AI block at `*(Group+3c)`, then replaces its serialized +4c
pointer with a fresh word list and restores that list. Constructor order0/enabled1 are not a
replacement for saved AI values. The normal resume linker receives argument0; its initial-stance
branch is skipped, but its Group+1c index insertion is not. Saved +1c remains a literal selector,
distinct from file-local object pointer keys. Dynamic first-save initialization is not closed.
— SAV-GRPIDENT-562, SAV-GRPAI-563

This is not whole-linker or first-tick preservation. A selected later link arm can write AI+48
through a resolved actor's Group. Group tick reads membership count before order+20 and writes
ff when empty. Arm0 reads actor ownership, dispatches the actor, then searches the live member
list before advancing. The rate helper follows actor+70 to Group+3c and reads AI byte+44.
Transitive effects and first reached chronology remain Unknown. — SAV-GRPAI-563

One Group dispatch selects order+20 once. Arms2/4/5 receive AI word+0a;
changing the byte in a returning callee does not restart selection. Every arm,
including an unrecognized byte, reaches a fresh current-membership withdraw
tail. The empty-group ff store is at entry only: becoming empty during an arm
does not itself rewrite the order. — SAV-GRPDISPATCH-568

Order0's inline successor search and the helper used by order3, ff and the
withdraw tail find the just-dispatched actor again in the current list. Missing
actor or final node ends that stage; a new successor can extend it. The tail
starts again at the current head. Actual callback mutations and deallocation
safety remain Unknown. — SAV-GRPMUTATE-569

The AI-owned path is not the actor's active patrol ring. The path-copy setter
walks AI+4c head-to-tail but prepends every point into actor order+90. Source
A,B,C becomes own,C,B,A; current order+02 is C, chosen before prepending own.
This is a setter rule, not a normal-LOAD reinitialization rule.
— SAV-GRPPATROL-570

Actor-order LOAD preserves raw cursor+02 and latch+04 while replacing list
pointer+90. Patrol consumes a nonzero latch before guard, then only guard
result0/b reaches continuation. That continuation sets the latch even before
arrival. Arrival searches for the first equal cell value, advances or wraps;
missing value reaches an unguarded null read. The cursor is not a node index.
— SAV-PATROLCURSOR-571

Group SAVE consumes its current embedded word list, AI block/list, members and
tail fields. No direct Group-dispatch instruction dereferences either Group
word-list payload; this is not a whole-call-tree absence. Group+20 semantics,
unnamed AI fields and original LOAD-to-dispatch-to-SAVE chronology remain
Unknown. — SAV-GRPSAVENEXT-572

## What this covers, and what owns the rest

| | owner |
|---|---|
| the blow itself — timer, roll, damage, experience | [`claims/hero.md`](../../claims/hero.md), [`claims/unit.md`](../../claims/unit.md) |
| walking there, pathing, tick order | [`formats/move`](../move/format.md) |
| the block/cost planes and the cell-record hash | [`formats/terrain`](../terrain/format.md) |
| the `.alm` type-5 player record's own layout | [`formats/alm`](../alm/format.md) |
| **this file** — the decision, from "nobody told me" to "that one" | — |

## Player input contract (`AI-INPUT-121`, `AI-SELECT-122`, `AI-PANEL-123`, `AI-MINIMAP-124`, `AI-KEY-125`, `AI-CURSOR-126`, `AI-INPUT-127`)

Input is edge-driven. The base dispatcher maps physical Windows messages to each surface's
vtable, gives a capture object first refusal, otherwise stops at the first child containing the
point, and gives a focused child first refusal on keyboard input. The map therefore never receives
a point inside the right column. There is no mission `WM_MOUSEWHEEL`, `WM_CANCELMODE` or
`WM_CAPTURECHANGED` handler (AI-INPUT-121/127, SESS-INPUT-037).

### Map mouse

| input | state / target | result |
|---|---|---|
| left down or double-click | marquee global zero | call OS `ClipCursor` and begin marquee; no engine capture or action yet |
| repeated left down or double-click | marquee global nonzero | consumed no-op; preserve the first origin and clip |
| left up | click-sized, no selected actor or select cursor | selection |
| left up | rectangle selection | select/toggle owned non-structures with strictly more than half overlap; equality fails |
| left up | attack cursor, `CUnit` target | order `0x19` |
| left up | attack cursor, other target or ground | move `0x16` |
| left up | move / swarm | `0x16` / `0x1a` |
| left up | defend, target / ground | `0x1b` / no order |
| left up | cast, unit / ground | spell `0x1e` / `0x1f`, item `0x25` / `0x26` |
| left up | pickup / town / patrol | `0x21` / `0x24` / `0x1d` |
| right down | map | set engine capture object and store pan origin |
| right move | non-zero signed `(delta pixels)/8` | pan by cells, mark drag |
| right up | drag / click | release only / send cancel `0x405` |
| right double-click | any | no-op |

The selection/order rows require an active marquee at `[0x005cd82c]`. Left up with that global
zero skips the `ClipCursor`, hit-update and selection/cursor-order prefix. It is a no-op when no
item is selected, but it still reaches the subsequent `session+0x3cc` item placement arm and can
emit `0x23` when `item+6 == 0xffff` or transfer record `0x22/0x32` otherwise. This edge ordering matters for repeated or unmatched button
messages (AI-INPUT-121).

The click/drag threshold is `screenW*10/640`, hence 10/12/16 pixels. The outer test sends an
axis strictly greater than the threshold to selection; the selection routine itself calls a
rectangle a click only when both axes are strictly smaller. This makes equality context-dependent
and must not be simplified to one `>=` test (AI-INPUT-121).

Selection: plain click replaces with the topmost intersecting non-structure regardless of owner,
or a structure whose class `+0x64` is zero. Shift-click and Shift-rectangle toggle an owned
qualifier only when the old selection summary at `view+0x144` has bit `0x04` clear. If a plain
click first selected a foreign actor, that bit is set and either Shift gesture over an owned actor
is a no-op; ground and foreign candidates also preserve the old selection. A plain rectangle
replaces only if at least one owned non-structure has strictly more than half overlap; equality
does not qualify. Alt after a single map selection expands to that object's group but does not
centre (AI-SELECT-122).

### Minimap mouse

Minimap actions occur on **left down**. Default sends camera `0x406`; move `0x16`; attack
`0x19` with any object id and `0x1a` without one; defend `0x1b` with an id and nothing without;
cast emits nothing; patrol `0x1d`. Left-drag repeats on every move. Right down and right-drag send
camera `0x406`. If both button bits are present on one move, the left test runs first and returns,
so only left semantics occur. No capture is taken. Left up only clears a special cursor; left
double-click, right up and right double-click are no-ops (AI-MINIMAP-124).

### Command cells and keys

| action | panel/key result | map click result |
|---|---|---|
| Attack / A | arm 1 | target `0x19`, otherwise move `0x16` |
| Move / M | arm 2 | `0x16` |
| Guard / G | immediate `0x17` | — |
| Defend / D | arm 4 | target `0x1b`, ground no-op |
| Cast / C | arm 5 | `0x1e/0x1f` or `0x25/0x26` |
| Swarm / S | arm 6 | `0x1a` |
| Stand Ground / T | immediate `0x18` | — |
| Retreat / R | immediate `0x14` | — |

Every consuming map/minimap action clears the armed mode. Guard, Stand Ground and Retreat repeat
on Windows key repeat; armed modes are idempotently reasserted. Although internal mode 8, the
patrol cursor and opcode `0x1d` exist, no shipped panel or key route leaves Patrol armed: R is the
immediate Retreat command. Exposing Patrol is an authored extension (AI-PANEL-123,
AI-CURSOR-126).

### Explicit Retreat execution

Ordinary-command Group construction has a narrower cleanup rule than a full sweep:
`0050fa33` deletes at most the first old Group rejected by `00449800`, then allocates.
The selected prologue does not directly mint Group+1c or populate the embedded
Group+20 word list. Constructor/list state and later behavioral callbacks are distinct
sources; Guard's GroupAI+20 order byte is not Group+20. — SAV-GRPCMD-578

Retreat is immediate command admission, not guaranteed immediate interruption.
The panel requires an active enabled cell, no selected inventory item and an eligible
inventory-bearing selection. Empty selection, a foreign first owner, or the structure
summary flag disables the panel. The producer includes entries with nonnull `actor+0x7c`
and caps the command list at 253. The dispatcher requires its AI manager and a resolvable
first actor; missing later actors are skipped. Its resolver rejects completed-death action
`0x10`, not every nonpositive-HP death phase (`AI-RETREAT-270`).

The command sets group order 0, actor state `0x16`, pending order 0 and current action 0.
It clears order fields `+0x38`, `+0x50`, `+0x60` and mover field `+0x7c`; it aligns desired
facing to current facing and releases a noncurrent reserved cell only at a cell centre.
It does not clear order-progress `+9` or all old target/spell slots. Conditional spell-object
cleanup is a separate dispatcher helper (`AI-RETREAT-271`).

Nonzero progress takes precedence over the next pending order. Progress 1 restores attack;
progress 2 restores the existing cast action. Both await counter `ord+0x15 > 2` and the
actor's animation flag `+0x136`. Progress 3 completes the current movement step to its cell
centre. Progress 4 remains held while its status bit is set. The setter's action reset is
therefore not proof that the old action was cancelled (`AI-RETREAT-272`).

While group order is 0, state `0x16` recomputes a flee order through the first withdrawal
helper below. Its radius is mover field `+8`, not a player-supplied destination. A zero
positive-HP count falls back to ordinary acquisition, which can choose a pursuit, idle or
autoheal order. The local state arm and helper have no timer, arrival termination or state
`0x16` reset; route failure also uses ordinary acquisition. Sustained whole-session
persistence and visible completion timing remain Medium/Unknown (`AI-RETREAT-273`).

Later admitted Move installs group order 4. Target attack replaces actor state with 3 or
fallback `0xc`; the two cast commands replace it with `0xd`/`0xe` or fallback `0xc`.
Nonpositive HP is handled before order execution in the unit tick and clears actor state
and action. Retreat does not bypass that death branch (`AI-RETREAT-275`).

The character panel is a second physical producer of Cast. Its `WM_LBUTTONUP` handler sends
selected item kinds 1, 2 and 5..8 to vslot `+0x7c`. The nested inventory grid's
`WM_LBUTTONDBLCLK` handler also forwards a resolved regular item to this same vslot. For an admitted owned single selection whose
item flags contain both `0x10` and `0x01`, `FUN_00492410` stores the inventory slot and effect
`0x29` context, then calls `FUN_0041b439(0x0a)`; the routine remaps that value to mode 5. The next
normal map down/up therefore uses the item branch and emits `0x25` for a unit or `0x26` for
ground. The same character-panel vmethod contains the sibling transfer producers `0x22/0x32`.
This path is not a ninth command-panel cell and must not be lost by enumerating only the eight-cell
object.

The nested inventory grid also produces a transfer on physical `WM_LBUTTONUP`. Live content at
`grid+0x84` is hit-tested through vslot `+0x88`. With no selected item at `session+0x3cc`, the
handler stops after that hit-test. With a selected item, the signed hit result is passed unchanged
to vslot `+0xa4 = FUN_00482820`. This grid's destination vslot returns code 2. Its insertion
routine merges an equal stack first, inserts at a hit in `0..count-1`, and otherwise inserts before
a trailing gold sentinel or appends. A margin or outside hit of `-1` therefore remains an item
transfer. `00482931` emits `0x22/0x32`, then clears the selected item and posts `0x46d`.

The grid's other physical mouse slots do not hide another order. Left-down scrolls the inventory
viewport by one through one of two edge strips. Left-button move over a live cell, with no item
already selected, selects a regular item or gold amount and builds its drag cursor; Shift changes
the requested amount from one regular item or 1000 gold to the source stack's quantity. That move
emits no order. Right down, right up and right double-click are consumed no-ops. These handlers
matter to the physical-input census even though only the later left-up reaches the transfer builder.

The same inventory-grid double-click has a separate gold-cell arm. A cell whose `item+6` word is
`0xffff` opens the prebuilt 296x168 Drop Gold modal at `(100,H-200)` only while the local purse is
positive. Enter or an inside left-down/up on `Ok` posts `0x445`; Esc or the cancel button posts
`0x446`. The action parses the edit, uses the amount returned by the owner resolver, subtracts it from the purse and emits
opcode `0x23` with that amount and the current map cell. Cancel closes without an order. The direct
caller set of `FUN_0041ce3e` is exactly the map selected-item arm and this dialog
(`AI-PANEL-123`, `SESS-INPUT-037`).

Other default mission keys:

| key | result |
|---|---|
| Esc | mission/town menu, unless focused child handles it |
| F1 / F2 | help / save |
| F3 | load in phase 2, Diplomacy otherwise |
| F4 / F9 | explicit no-op |
| F5..F8 | Ctrl assigns quick-spell slots 0..3; plain key selects the binding, with gated Cast-mode arming when the book is closed |
| F12 | toggle an unidentified global |
| Pause | modal pause text |
| numpad `+/-` | speed index, clamped 0..8; Ctrl+plus selects the unpaced/max-speed idle loop, Ctrl+minus restores paced mode and resets its epoch |
| digits/numpad digits | select group; Shift augments; Ctrl assigns; Alt selects and centres; Ctrl wins |
| arrows | pan one cell per repeat |
| E | select all owned exact-name `CUnit` objects |
| Space | toggle/close inventory and spellbook together |
| B or Q / I or backtick | spellbook / inventory |
| Enter | open and focus text entry; subsequent map shortcuts are suppressed |
| Tab with character panel focused | send panel custom message `0x412`; precise effect Unknown |
| Ctrl+F/H/L/N/O/U/W | retreat mode, health, smoothing, day/night, flying damage, autoheal, formation |
| Alt+S | screenshot |

Modifier latches set on down, clear on up and all clear on focus loss. Keydown does not inspect the
repeat flag. Backspace clears three unidentified global containers; Alt+B..Y except S emits an
unidentified outbound record `0x46`. Tab's `0x412` route exists only when focus dispatch reaches
the character panel. Unlisted keys have no mission action (AI-KEY-125).

### Quick-slot state and lifecycle

The four signed quick-slot values live on the shared spellbook controller at
`campaign+0xec`, offsets `+0x64/+0x68/+0x6c/+0x70`. Construction writes `-1`
to all four. Current spell `+0x60` and selected-item override `+0x74` are
separate fields. The reached keys do not index actor/group-owned bindings.
Message `0x411` and controller right-down/up clear current spell only;
selection-summary recomputation has no direct binding write. Book open/close
attaches/detaches the same controller. Selection carry beyond these bodies is
Medium; no-load mission entry and new-campaign reset remain Unknown.
— AI-QUICKOWNER-280

F5–F8 send `0x417`, slot 0–3 and the Ctrl latch to that controller. Ctrl copies
current spell when it is not `-1`; otherwise it uses a nonnegative spell-grid
mouse hit. Each assignment empties all other equal bindings to `-1`, so a
duplicate moves rather than swaps. The hover arm neither checks current spell
availability nor changes the current-spell field. Ctrl still reaches the
key handler's subsequent arming checks. — AI-QUICKASSIGN-278

A plain key copies a nonnegative binding into current spell. If the book is
closed, an available stored-index bit in `view+0x148` requests action 9,
which maps to mode 5 only when Cast capability is present. It does not open
the book or directly emit a cast order. An empty `-1` slot leaves local mode
and current spell unchanged; an unavailable nonempty slot changes current
spell but requests no new mode. Neither local branch cancels an existing
mode. An open book skips mode arming. Generic child dispatch and malformed
saved indices remain outside this bounded result. — AI-QUICKINVOKE-279

The original save programme writes the four indices in F5–F8 order to
`SpellBook/Shortcuts` and restores them to the same controller. It separately
saves/restores `SpellBook/Pressed`. The reached slots are persistent save
state, not merely a transient UI cache; a populated original save/reload was
not runtime-witnessed. — AI-QUICKSAVE-281

### Book identity and selection predicates

For valid book cell `i=0..23`, current spell and quick bindings hold `i`.
Map targeting passes `i+1`; opcodes1e/1f carry this ordinal as a byte at
command+10. The dispatcher rewrites it through the byte at
`005c2328+4*(i+1)` before actor-book lookup. That is the same storage as the
statistics table at `005c232c+4*i`. Cell5 means command6 and intrinsic spell23,
not spell6. Lookup uses actor+140 and returns the sparse ID-indexed element;
its pointer is parked at actor+44. Malformed indices and stale actor spell
pointers remain outside this valid-cell contract. — AI-SPELLIDENT-286

Selection recomputation ORs object+18 into view+148 for every selected object
with nonnull+7c, before class/type/session predicates. Count and the primary
object use this same accepted population. The two book producers independently
filter those selected objects by the chosen cell bit; the append helper caps
the command at253 members. Snapshot freshness and all upstream selection
restrictions remain Unknown. — AI-SPELLPOP-287

Cast capability is narrower than availability. In session states+3dc==1 or
with bit2 set, each accepted `CUnit` with+20 equal17h/18h and nonzero+18 can
OR200h into view+144. This is the current loop object, not only the primary
object; the primary-only clause of AI-PANEL-061 is corrected. Other session
states assign summary8. The owner comparison still uses the primary. The
mode helper returns0 for summary bit4, otherwiseEFh plus10h for200h; the
panel separately disables on zero count or summary24h. Item action0Ah bypasses
the mode-bit refusal. Static predicates do not prove every synthetic mixed
population is reachable. — AI-SPELLCAP-288

Mouse cell selection requires the cell's availability bit. A nonnegative
shortcut can become current before its later arming guard, so unavailable
does not mean impossible as current. The reached map Cast branch takes
nonnegative item+74 before current+60 and rejects a negative result. Its
book target flag is indexed by cell at005bcef8:14 of24 entries are nonzero.
A unit-like hit plus nonzero flag routes to actor production, other hits to
cell production; without a hit a nonzero flag emits no order.
— AI-SPELLGUARD-289

Item opcodes25h/26h instead carry the word from controller+78 at command+10.
The consumer retrieves that inventory item and constructs a Spell from its
first kind29h effect's byte+40. It skips the book ordinal translation and
the producer's learned-spell bit filter. Its target predicate uses005bcfb8.
Item snapshot synchronization and the complete use/refund lifecycle are
separate boundaries. — AI-SPELLITEM-290

## The objects

| | what |
|---|---|
| **session** | `[0x005f21c4]`, `0xc320` bytes. `+0xa50` world, `+0xa9b8` `MinimalGuardRange`, `+0xa9c4` diplomacy, `+0xba4` / `+0xbc4` the live / dead candidate lists — each an outer object whose list sub-object begins **four bytes in**, at `+0xba8` / `+0xbc8`, so their element counts are `+0xbb4` / `+0xbd4` (`list + 0xc`). The pair is **shared scratch**: per-actor acquisition rebuilds it too, so nothing in it survives another routine's call |
| **world** | `0xa4558` bytes. `+0x10000` block plane, `+0x5400c` the occupancy scratch list, `+0x540b8` the cell-record hash, `+0x82ef0` the visibility map, `+0x9451c` heights, `+0xa4554` the actor list |
| **actor** | `+0x10` position, `+0x14` `Player`, `+0x4c` flags (bit 3 = off map), `+0x5c` combat target, `+0x94`/`+0x96` health, `+0xa4` sight (`u16`, 1/256 cell; `+0xa5` is its whole-cell high byte), `+0x12c` reach, `+0x144` state bits (`0x8000` invisible), `+0x154` mover, `+0x158` order block |
| **order block** | `+0x00` guard post cell, `+0x02` current patrol waypoint, `+0x08` state, `+0x0c` target, `+0x14` stop distance, `+0x20` the group's assigned target, `+0x58`/`+0x5a` remembered attacker cell and its age, `+0x70` follow stop distance, `+0x71` see-invisible radius, `+0x90` the patrol waypoint list |
| **group** | `0x48` bytes, a `CObList` (`+0x04` head, `+0x0c` count) in `player+0x24`. `+0x1c` the authored group id, `+0x3c` the AI record, `+0x44` the owning `Player` |
| **group AI record** | `0x50` bytes. `+0x00` current guard post, `+0x20` the group order, `+0x24`/`+0x28` fine and cell centroid, `+0x2a` spread, `+0x2b` max member sight, `+0x2c` derived notice radius, `+0x2d` working radius, `+0x38` base radius, `+0x45` AI enable (default 1), `+0x4c` the patrol path |

## The decision, in the order the engine takes it

1. **Population — sight.** Clear `world+0x82ef0` (`0x10000` bytes, one per cell). Stamp the
   deciding actor's line of sight into it (below). If the actor is AI-owned and remembers being
   attacked, mark that cell too — for 20 ticks. Then walk the global on-map actor list and keep
   every actor standing on a marked cell.
2. **Relation.** Drop every candidate `c` for which `matrix[me.Player+0x04][c.Player+0x04] & 1`
   is 0. Drop an invisible candidate (`c+0x144 & 0x8000`) unless some member of the decider's
   group is within that member's own `order+0x71` of it.
3. **The dead.** Move survivors with `health < 1` to a second list. If nothing living is left
   and that list is not empty, move it back: a unit with only corpses in view targets a corpse.
4. **Selection.** Minimise the footprint-aware edge distance `d` (1 = touching), with `d + 1`
   for a flying candidate when the decider does not fly; break ties on the 16-way turn cost from
   the mover's current facing. Seed the best distance at `reach + 1`, and **discard the winner
   if its distance exceeds `actor+0x12c`**.
5. **Outcome.** A winner sets `order+8 = 6`, `order+0xc = target`, `order+0x14 = reach`. No
   winner: a human participant's unit tries to heal; an AI-owned unit goes to the guard state,
   which walks it back to its post whenever it is farther from the post than `actor+0xa5`.

**Sight decides the population and reach decides the pick.** An engine that acquires on sight
alone starts fights the game does not; one that scans only within reach acquires through walls,
and acquires things it cannot see. Neither radius is a substitute for the other.

### Prismatic Spray is a distinct consumer of the group population (`AI-SPRAY-266`, `AI-SPRAY-267`)

`FUN_0053ddd0` is the fifth direct caller of the group builder. It consumes shared list A head to
tail and then list B head to tail; it does not run a radius search. The ordinary secondary gates are
therefore the builder's group-wide visibility, first-member diplomacy and See-invisible rule, then
its living/corpse split. The primary spell target is separate: the selector alarms it and performs
the hostility flip before the builder call, appends it first, and skips an equal secondary later.
It can therefore survive conditions that would remove it as a secondary.

Prismatic rank differs from ordinary acquisition. List A uses
`((edgeDistance<<8)+turnCost)&0xffff`; list B uses
`((((edgeDistance<<8)+turnCost)&0xffff)<<8)`. Repeated strict-minimum scans preserve source order on
equal stored scores. Selection first requires a score below 65530 and append later requires the
saved threshold below 65000. A low enough B score can therefore compete beside living A; whether a
valid layout can supply the required edge distance is Unknown. If there are no living entries, the
builder has already moved B back into A and the only-corpses population uses the A formula. The
spell's capped byte limits the final output list after ranking; it is neither actor reach nor order
`+0x14`, and the primary consumes one position on the ranked path. When the builder returns empty A, the selector instead
appends the primary and returns before reading the cap. Its stack pointer region holds 100 entries
and its winner/threshold regions ten; neither population overflow nor a custom cap above ten is
checked. Candidate scores are session scratch at `+0xd74`, whose capacity is Unknown.

### The sight stamp, complete (`AI-GROUPSEE-068`, `AI-LOS-087`…`AI-LOS-091`, `AI-SIGHT-092`…`AI-SIGHT-094`)

The sight object is embedded at `world+0x58ee8`. Its init runs in each of the four world
constructors, **before the map is loaded**, and builds three tables that depend on
`k = [Scanning] ScanShift` of `World\Data\map.reg` (ships **7**) and on nothing else — no terrain,
no actor, no tick. Two of them serve the stamp: `step`, `+0x22000`, and `cost`, `+0x28000`, both
`u16`/`i16` over a 41×41 window inside a 64×64 grid, addressed as `base + (a<<7) + 2b` with `a`
the column offset and `b` the row offset, both `0..40` around a centre of `(20,20)`.

```
step[a][b]  = the cell one Bresenham step toward the centre, as (i8 dx, i8 dy) in {-1,0,+1}
              zone by slope:  j <  i>>1 -> column step | j > 2i -> row step | else diagonal
              the builder's four quadrant mirrors collide on the j==0 axis, so it repairs
              cells (+1,0) and (-1,0) with four literal stores at the end
cost[a][b]  = ftol( (1<<k) * sqrt(i*i + j*j) / max(i,j) )   -- the mean length of one step
              along that ray, in 1/(1<<k) cell.  k=7: 128 on the axes, 181 on the diagonals

stamp(actor at cell C, scanRange = actor+0xa5):
  acc[*][*] = 0 ; acc[20][20] = (1 << (k-1)) + (scanRange << k)      -- a budget in 1/128 cell
  vis[C] = 1 ; alt = (i8) height[C]                                  -- height = world+0x9451c
  for r = 1 .. 19:                            -- Chebyshev rings; ring 20 is unreachable
     allBlocked = true
     for each of the ring's four edges, 2r+1 cells each, in edge order:
        skip the cell unless 8 <= absCol <= W-9 and 8 <= absRow <= H-9   (world+0x58ee0..3)
        p = step[a][b]
        v = acc[a+p.dx][b+p.dy] - cost[a][b] - (i8) height[cell] + alt
        acc[a][b] = v                         -- stored EVEN when it blocks
        if v <= 0: continue                   -- blocked: NOT pruned, the ray is not cut
        vis[cell] = 1 ; allBlocked = false
     if allBlocked: stop
```

**Three things a consumer gets wrong by default.** (a) It is a **budget, not a radius**, and the
dominant term is the observer's own altitude: unrolled, `acc(n) = seed − Σcost − Σh(cell) +
n·h(observer)`, so the observer's height is added at *every* step while a cell's own height is
charged *once*. On flat ground the region is a disc of radius `scanRange` (145 cells at 6); over
the shipped corpus the same unit sees 29 to 1 248 cells depending only on where it stands, and a
single tall cell costs its height once and casts no shadow behind itself. (b) A blocking cell is
**not pruned** — the march writes the negative value and walks on. Nothing is ever re-lit behind a
blocker on shipped maps only because the cheapest step costs `1<<k = 128` while every shipped
altitude byte is `0..127`, so the rise can never reach 128; an authored byte `>= 0x80` breaks that
by one unit. (c) The window centre is the actor's **own** cell in all three stamp sites.

**This is the same algorithm the drawn fog runs** (`AI-SIGHT-093`), re-derived from both
listings and executed against itself: the two `pred`/`cost` tables differ in **0 of 1681** cells,
the seeds are algebraically equal for a whole-cell sight, and over 43 640 corpus observers more
than 27 cells from any map edge the two regions never differ. A consumer needs **one**
implementation, parameterised by two things:

- **the seed's width.** The fog reads the whole `u16` at `actor+0xa4` — sight in 1/256 cell — and
  shifts it right by `8 − k`; the sight stamp reads only its **high byte** `actor+0xa5`, the
  whole-cell radius. A monster's sight is whole cells and the two coincide; a hero's is
  `ftol(((mind + reaction)/25 + 4) × 256)` and is not, so **the AI grants a hero less sight than
  his own fog shows him** — 105 cells against 145 at `actor+0xa4 = 1535` (`AI-SIGHT-094`).
- **the playable rectangle.** `8 .. W−9` here against the view's `7 .. W−8`, and byte-wide compares
  here against 32-bit ones there (`TERR-FOG-118`).

`actor+0xa4`'s six writers, and the two forms of sweep that see them, are `AI-SIGHT-092`; a
`disp:a5` sweep sees two of the six.

**A unit only closes distance when it is given a pursuit order** — `order+8 = 5`,
`order+0xc = target`, `order+0x14 = reach` as the stop distance. Acquisition never produces one;
the guard and engage arms do.

**Two clauses suppress engagement entirely**, both requiring all three of `actor+0x4c & 4`,
`actor+0x12c < 2` and a human-participant owner (`Player+0x28 == 0`).

## The diplomacy matrix

50 × 50 bytes at `session+0xa9c4`, all zero at construction — with nothing written, no unit ever
finds an enemy. Cell `[i][j]` is indexed by `Player+0x04`, stride 50; the engine reads three
bits of it.

| bit | meaning |
|---|---|
| 0 | *i* treats *j* as hostile. This is the only bit acquisition tests. |
| 1 | locked: combat may not turn this pair hostile. Forced on every diagonal. |
| 2 | read (`AND …,0x7`) and unused by any shipped map. |

**Row 0 and column 0 are not padding.** Column 0 holds each player's owner kind and row 0 marks
the slot live; the join routine writes both and reads them back. A consumer that treats index 0 as
unused loses the join rule.

### The six writers

Ordered by when they run. Only the last two are symmetric, and only one of the six re-scans.

| # | when | routine | writes | symmetric? | respects bit 1? |
|---|---|---|---|---|---|
| 1 | map load | `FUN_004e1924` | the `.alm` type-5 record's sixteen `u16` at file `+0x2c` into columns 1…16 of that player's row (low byte only), then forces the diagonal to 2 | one way | n/a |
| 2 | a player joins | `FUN_0053d8a0` | column 0 = owner kind (kind 2 stored as 0), row 0 = 1, then for every live slot one of four template bytes chosen by whether the two column-0 flags agree; forces the diagonal to 2 | both | no |
| 3 | a player leaves | `FUN_0053d9a0` | row 0 = 0 for that slot | n/a | n/a |
| 4 | a mission join | `FUN_004d303e` / `FUN_004d8963` | clones a reference player's row *and* column, allies with it, forces neutrality or 2 between participants | both | no |
| 5 | script action **10** | `FUN_00539be0` | `matrix[p0][p1] = (v &~ 3) + p2` | **one way** | **no — it clears it** |
| 6 | session command **0x45** | `FUN_004d5dd8` | assigns the whole row `matrix[setter][*]` from a `u16` array in the command body, each element `& 7` | **one way** | **no** |
| 7 | a blow landing, or a spell cast | `FUN_0053d9b0` | sets bit 0 (`OR AL,0x1`) in both directions, each direction **separately** gated on `(cell & 3) == 0` | both | **yes** |

The template of writer 2 is four bytes at `session+0xa9bc`, set by the constructor to
`{1, 1, 0, 0}`: with those values a joining player is hostile both ways to exactly the live
players whose column-0 flag differs from his, and neutral to the rest.

### What a change reaches

Writers 1–6 are **bare**: they change the byte and nothing else. The one exception is the
spell-cast entry into writer 7, which rebuilds the caster's group candidate list in the same
instruction stream. Everything else propagates only because a consumer re-reads the cell — which
for the 95.6 % of hostile placements under a group order means the next AI tick, since the
candidate list is rebuilt and re-filtered on every evaluation. **An order already issued keeps
running**: making a faction friendly does not stop a unit that is already attacking, it stops that
unit being re-selected.

### What the script can do with it

Action opcode **10** writes `matrix[p0][p1]`, one direction, clearing the low two bits and then
**adding** `p2`. Nothing clamps `p2`, so a value ≥ 4 carries into the bits above and ≥ 256
truncates in the byte store; shipped maps use only `{0, 1, 2}`. Check opcode **10** reads
`matrix[A][B] & 3` into a script slot — **narrower than the action can write**.

Shipped content uses both: 11 action nodes over 7 maps and 5 check nodes over 5 maps. A mutual
change costs **two** nodes, and the shipped maps author them in pairs
(*"Alliance Self→Peasants"* beside *"Alliance Peasants→Self"*).

### Persistence and the wire

The matrix is serialised **verbatim** as part of the 2508-byte (`8 + 2500`) sub-object at
`session+0xa9bc`, in the world half of a save, and read back verbatim. Nothing re-derives it from
the map on a save load. Outbound, a changed row is broadcast as message `0xb9` carrying the whole
row widened to `u16`; inbound, session command `0x45` receives a row in the same width. The
relation is byte-wide in the simulation and **word-wide on the wire in both directions**.

**The relation is directional.** Over the 38 shipped maps, 102 of 866 ordered off-diagonal
pairs disagree with their mirror on bit 0, across 19 maps: guards attack monsters that will not
start a fight with guards. A symmetric store is the wrong shape.

`World\Data\ai.reg` contributes one value on this path: `[Scanning] MinimalGuardRange`, a floor
on a group's notice radius. The code default is 10; **the shipped file sets 8**. The file has four
records in total and its other scalar, `[Tasker] IntelligentCons` = 15, has no literal anywhere in
the image — nothing reads it.


## Who runs it, and how often

The AI is **one slot of the full tick**: `server+0x04 % 16 == 6`. It is driven **per group**, never
per actor — the driver walks the player list, then each player's group list, and runs the group
order machine once per group, skipping a group whose AI-enable byte is 0. At the shipped speed
index that is one AI pass per about 992 ms.

## Two levels of behaviour

The **group order** (`grpAI+0x20`) is evaluated first. Eight values are live. The names in
parentheses are the shipped editor catalogue's labels for the script command that sets each; three
of the four this table used to leave open do **not** describe what the arm does.

| value | what the group does |
|---|---|
| 0 | run **each member's own state** — the only value under which `actor+0x50` matters |
| 1 | **guard**: build the list, clip it to the notice radius around the group centroid, score it, engage; a targetless member walks back to **its own post `ord+0x00`** (not the centroid) and, once there and idle, takes the idle-turn order `ord+0x08 = 0xb`. This is the order the load walk gives **every other** player's groups |
| 2 (`Swarm`) | build the list, score it, then per member: engage the scored target; else walk to the commanded cell `grpAI+0x0a` if not already standing on it and idle; else cast (scenario owner) or run the heal AI (a human participant's unit). **No formation and no spread** |
| 3 (`Stand Ground`) | **stand ground**: build the list with **no radius clip**, score it with the *reach-bounded* scorer, engage what is already in reach and otherwise turn on the spot. There is no walk anywhere in the arm. This is the order the load walk gives the **player's own** groups |
| 4 (`Move`) | walk each member to **its own** `ord+0x0a` — the arm ignores the cell the dispatcher passes it; the destination arrives from the command's setter. On arrival: stop, clear the route list, set the stop distance to the member's reach, re-acquire |
| 5 (`Swarm 2`) | **arm 4 with a pre-emption.** The arm builds the candidate list and reads its element count, `AImanager+0xbb4`: **zero — nothing visible — and it tail-calls arm 4 with both arguments forwarded**; non-zero and it runs 2's body without the walk. The command's setter is arm 4's setter to the byte (one immediate apart), so the per-member state the fallback runs on is arm 4's own. **A vetoed group is not the same as a blind one:** candidates present but all vetoed by the preference matrix leaves every `ord+0x20` at 0 inside this arm, and both of its zero-target branches — idle turn plus optional cast, or the heal AI — **stand still** |
| 0x11 (`Roam`) | keep a Group-AI cell in `grpAI+0x0a`; re-roll when the farthest member's distance is **< 10** or byte counter`+0x15` is **> 50**, by stepping **20 cells** in one of eight random compass directions and rejecting anything outside the playable rectangle; reset the counter on acceptance, call arm 5 evaluation and increment the counter after normal return |
| 0xff | forced when the group is empty: run the standing acquisition for each member |

Values 6..0x10 and 0x12..0xfe do nothing. **After every arm the withdraw tail runs**, below.

Roam's local result is the Group-AI cell/counter program and the Swarm2 call,
not a proved copy into member destinations. Swarm2 forwards the cell only on
its zero-candidate Move fallback; Move ignores that argument and uses each
member's current `ord+0x0a`. Reached movement can depend on member order and
called helpers. These bodies establish neither guaranteed wandering motion
nor absence of all other wander paths. — AI-ROAM-025, AI-SWARM2GATE-107,
AI-MOVE-023

The **actor state** (`actor+0x50`) is a 27-arm switch. **Its value at construction is `0xb`, guard.**
The arms a consumer needs:

| value | what the actor does |
|---|---|
| 0 | idle |
| 1, 2 | re-issue or complete a stored move |
| 3, 4 | engage a stored target, or the first hostile in the block around a stored cell |
| 8, 0x11 | follow an actor; the stop distance is `order+0x70`, or the actor's sight when that is 0 |
| **0xa** | **patrol** |
| **0xb** | **guard** |
| **0xc** | acquire (the decision above) with no post and no leash |
| 0xd, 0xe, 0xf, 0x17 | further order forms; `0x17` also overwrites both of the actor's radii with 5 |
| 0x16 | **explicit Retreat** — installed by the player command through a register-form store; repeatedly calls the first automatic-withdrawal helper, without that tail's HP-threshold admission (`AI-RETREAT-271`, `AI-RETREAT-274`) |
| anything else | clear the sub-tick gate and acquire |

## The group engagement machine — how a candidate becomes a target

Everything here is what a group order arm runs before it touches a member; a consumer that gets the
per-actor states right and this wrong reproduces neither who is attacked nor when.

**1. One sight map for the whole group.** The builder clears the `0x10000`-byte visibility array
once, then stamps *every* member's sight into it, so a candidate seen by any member is a candidate
for all. For an AI-owned member the cell it was last struck from is forced into the map as well and
expires after **20** group ticks.

**2. The sweep is global.** Every actor on the map is tested against that map — head to tail of the
world actor list, no neighbourhood and no spatial index.

**3. One decider.** The diplomacy filter runs with the group's **first member** as `me`, so one
player row governs the whole group and one member's group list governs the See-invisible exception.

**4. Corpses are parked, not dropped.** Candidates below 1 health move to a second list; if the
first list ends up empty **the whole second list is moved back**. A group that can see only corpses
therefore has a candidate list and will attack one.

**5. Score, assign, and forget.** Each member is given the cheapest candidate in `ord+0x20`, or 0.
This happens on every evaluation: there is no memory of last tick's target anywhere.

**6. The list's element count is a decision input, and it is read at two widths.** Group order 5
branches on the full dword `+0xbb4`; the scorer tests only its low byte, as it does for the group's
own member count. A count that is a non-zero multiple of `0x100` therefore passes the gate and
scores nothing — a bounded divergence rather than an observed behaviour, and one that raising an
actor cap would make reachable.

The cost, lower being better:

```
d    = Chebyshev(member, candidate) in cells
k    = candidate movement domain, forced to 0 when the candidate's reach > 1
pref = M[member domain][k]        for a member of reach <= 1   (melee)
     = M[0][k]                    for a member of reach >  1   (ranged)
d   := 1        when d <= member reach          (ranged members only)
     = d + 1 - member reach   otherwise         (ranged members only)
cost = (d << 8) + turnCost(member facing -> direction of candidate)
```

then `pref == 0` -> **never pick this candidate**; `pref == 1` -> `cost * 1.5`; `pref == 4` ->
`cost * 0.75`, the last two only when the candidate is single-cell with Mind >= 15. A candidate
carrying active spell id 20 costs a flat **+127**.

**M**, indexed `[member domain][candidate domain]`, is a compile-time constant:

| member \ candidate | 0 | 1 | 2 | 3 (flier) |
|---|---|---|---|---|
| 0 | 2 | 1 | 2 | 4 |
| **1 (ground)** | 4 | 2 | 1 | **0** |
| **2** | 2 | 1 | 4 | **0** |
| 3 (flier) | 2 | 1 | 2 | 4 |

Two consequences a consumer will not guess. **A melee ground creature never auto-selects a flier**
— and every human actor is domain 1, because the Humans table ships no `MovementType` cell and the
actor constructor defaults the field to 1. **A ranged creature uses row 0 always**, which contains
no zero, so reach and not class decides whether anything will chase a flier. Fliers are placed on
30 of the 38 EN maps.

Group order **3** uses a second scorer that additionally refuses any candidate whose distance term
exceeds 1 — that is the whole reason its members never move.

**Two clocks.** The group decides once per full tick; the order it writes is executed by the actor
tick, which runs far more often. Nothing in the group arms is sticky because nothing needs to be.

**The notice radius is frozen.** The group's geometry — centroid, spread, max sight — is recomputed
on every guard tick, but the only reader of the derived radius is the routine that *issues* a guard
order. The circle follows the group; its size is whatever it was when guard was last issued.


## Guard — the whole of motion without an order

**Where the post comes from is decided one level up, and it is decided twice.** A group under
group-order 1 or 3 never runs the per-actor state machine at all, so for it steps 1–4 below are
not the live path; what anchors its members is the **stance setter** that put the group under that
order. Every stance setter, before it writes the group order byte, walks the members and writes
each one's post: from the actor's current cell if it is standing on a cell centre, and from the
cell it is **stepping into** if it is mid-move. Only the two guard forms have that second branch;
the two stand-ground setters take the current cell unconditionally. At map load the placements are
all centred — every placement entry point of the position record writes the centre offsets — and
the spawner runs seventeen bytes before the stance walk, so a load-time post *is* the spawn cell.
Steps 1–4 are the live path only for a group at order **0**: one built by a player command, or —
because the load-time stance walk's session gate skips both setters together — every group on a
multiplayer load.

1. If the guard post is unset, **it becomes the actor's current cell**. Home is emergent: nothing
   authors it, and it is not the spawn point recorded in the map — it is wherever the actor
   happened to be when guard first ran.
2. If the actor is farther from the post than its **sight** radius, it is ordered to walk back.
   *This step is behaviourally dead*: step 3 overwrites it whenever it fires and step 4 writes
   exactly the same order. It survives only when the engage refuses and the acquire finds nothing,
   which needs a human owner.
3. The occupancy block of **Chebyshev radius 5** around the **post** — not around the actor — is
   filtered for hostiles and any survivor is engaged with a pursuit order. The 5 is `mover+0x08`,
   a **hard-coded constant** written once by the mover constructor; no shipped byte carries it and
   no class varies it.
4. With nothing to fight: away from the post, walk back; at the post, run the standing acquisition.

**This is what ends a pursuit.** Nothing inside a pursuit order measures elapsed time, health, or
the distance the pursuer has covered; the decision is re-taken from scratch every group tick, and
its only geometric input is *post to target*. Back a target out of the 5-cell block and the actor
is ordered home on the next tick, however far it has already chased; keep it inside and the actor
re-engages, however far it has already chased. **The post moves only when something re-issues a
stance.** A creature nobody commands keeps its spawn cell for the whole mission; a unit given
Guard again somewhere else is re-anchored there, and so is the 5-cell block. Besides the stance
setters, `ord+0x00` is written by the three target-died teardowns, the arrived-branch of the
order-progress epilogue, and patrol's re-anchor. A consumer must not model the post as immutable,
and must not model it as *"set on the first tick of guard"* either — for a load-time group it was
set before the first tick ran.

The second terminator is the route search: when it returns an empty path the mover sets
`mover+0x98`, and the order-progress epilogue cancels the order and re-acquires within **reach**.
A consumer that omits this will leave units frozen against unreachable targets.

This actor guard path does not define the Group Roam program above. An unordered actor stands on its post
until something hostile comes into the block around that post. Guard is the whole of motion for a
group left at its **load-time** order — it is not the whole of motion without a *player* order,
because the map's script can put a group on patrol (below).

Two fields of the group AI record belong to this arm and to nothing else: `grpAI+0x30` is a
**has-members latch** and `grpAI+0x34` a **tick counter with no reader**. The latch is what gates
the radius roll below.

## Patrol

State `0xa`. The actor **runs guard first**, so a patrolling actor still fights; if guard leaves it
idle, the patrol advances. The waypoints are a **ring**: on arrival at the current one the list is
searched for that cell and the next node is taken, falling back to the head.

**Two live surfaces put an actor into this state, and the one shipped maps use is the map's own
script.** The `.alm` type-7 action opcode 6 is a group command; its first parameter `14` is the
shipped catalogue's `Group Command : Patrol`, taking `X`, `Y` and a target group. The command:

1. stops every member and sets the **group order to 0**, which is what lets the per-actor state
   machine run at all;
2. per member: `actor+0x50 = 0xa`, guard post := the member's current cell, and the waypoint list
   is emptied and rebuilt as **exactly two nodes — the member's own current cell, then `(Y<<8)|X`**;
3. the current waypoint is set to the commanded cell.

So a shipped patrol is one commanded point and wherever the creature was standing, walked back and
forth forever. `order+0x04` is a **re-anchor latch**: set on every advance and consumed at the next
entry to move the guard post to the actor's current cell — which is why the guard leash in step 2 of
the previous section never pulls a patrolling actor home.

The same per-member routine serves a **player-issued** patrol command, a sibling of the guard and
aggressive commands in the same dispatcher.

The other surface is the **mission description file** `World\Mission\<n>.ini`, a section of `x;y`
lines with both coordinates strictly inside 8..135; a line without a `;` at index 1 or above is
skipped with an operator message, and a group whose path came out empty is not put into patrol. It
copies the group's path into each member and adds that member's own cell to it. **No such file ships
in the GOG install** — 0 nodes in any archive, 0 files under any `Mission` directory in either
preserved root — so no shipped patrol comes from *there*.

**Corpus:** 14 patrol commands ship, on 8 campaign maps, naming groups of 34 placed creatures, none
of them the player's. Both preserved roots carry the same 14.

A consumer implementing the ring should note that the engine's own search dereferences null when the
current waypoint is not a member of the list. No live setter can produce that state.

## The third radius

Reach decides the pick and sight decides the population (above). The **notice radius is a third
thing and it belongs to the group**:

```
centroid   = mean of the members' fine positions
grpAI+0x2c = max over members of ( Chebyshev(centroid, member) + member.sight )
grpAI+0x2d = grpAI+0x38 = max( grpAI+0x2c, [Scanning] MinimalGuardRange )
on the tick grpAI+0x30 flips from 0:
           grpAI+0x2d = grpAI+0x38 + 4 + rand3
           rand3      = (3 * rand()) / 0x8000 - 1,  i.e. exactly {-1, 0, +1},
                        each on about a third of the RNG's range
           (the +4 is on that edge only; the empty-group edge rolls without it)
```

`0x8000` is not a literal in the roll: it is `AImgr+0x00`, written once by the AI manager's base
constructor and never again, and it is `RAND_MAX + 1`. **`n * rand() / AImgr[0]` is the module's
uniform-range idiom** — fifteen sites use it — so any such quotient is `floor(n·rand()/32768)`,
uniform on `0…n−1`. The idle-turn arm is the same idiom at `n = 190`.

The radius clips the group's **candidate list** around the centroid under group order 1. It is
never authored, and it never decides whether a blow lands.

## The player's order vocabulary

Nineteen opcodes, `0x14..0x26`, dispatched by `FUN_004d5dd8`'s order space (19 direct dwords at
`0x4d86ae`; `formats/session` owns the command itself). **Four are empty** — `0x15`, `0x20`, `0x22`,
`0x23` — and every live arm first tests the AI manager `[0x005f21c4]` and does nothing when it is
null.

**Every order builds a new group.** Before the switch, the prologue resolves the commanded actors,
allocates a fresh `0x48`-byte group onto the player's list and fills it with exactly them. The
group's AI block is constructed with **`grpAI+0x20 = 0`**, so a commanded group always *starts* at
group order 0 — which is the gate the per-actor state machine runs behind. Nothing about the load-time
1-or-3 stance survives a player order.

The orders then split by what they leave that byte at:

| opcode | dispatches to | leaves `grpAI+0x20` | per-actor `actor+0x50` |
|---|---|---:|---|
| `0x14` | `FUN_00534e20` | 0 | `0x16` |
| `0x16`, `0x1c` | `FUN_005340a0(col,row)` | **4** — `Move` | from the argument |
| `0x17` | `FUN_00534ac0` — **guard** | **1** | `0xb` |
| `0x18` | `FUN_00534cb0` — **aggressive** | **3** | `0xc` |
| `0x19` | `FUN_005348c0(target)` | 0 | `3`, then `0xc` |
| `0x1a` | `FUN_00534390(col,row)` | **5** — `Swarm 2` | from the argument |
| `0x1b` | `FUN_00534f80(target)` | 0 | `0xc`, `8` |
| `0x1d` | `FUN_00535300(col,row)` — **patrol** | 0 | `0xa` |
| `0x1e`, `0x25` | `FUN_00533d00(target, spell)` | 0 | `0xd`, then `0xc` |
| `0x1f`, `0x26` | `FUN_00533ed0(col,row, spell)` | 0 | `0xe`, then `0xc` |
| `0x21` | `FUN_005308a0(col,row)` — **pick up a sack** | *unwritten* → 0 | `2` |
| `0x24` | `FUN_0052f8b0(building)` | *unwritten* → 0 | `0xf` |

**Four orders act through a group arm** (guard, aggressive, move, `0x1a`); **ten act by writing a
per-actor state** and leaving the group at 0, which is the only value under which
`FUN_00533ae0` evaluates `actor+0x50` at all. That is why patrol and the pick-up are per-actor
behaviours and guard and aggressive are group behaviours — the difference is one byte written at the
end of the order routine.

**Nine of the fourteen live arms have the dispatcher as their only caller in the image** —
`0x14`, `0x19`, `0x1a`, `0x1b`, `0x1d`, `0x1e`, `0x1f`, `0x21`, `0x24` — while the map's script
reaches `FUN_005340a0`, `FUN_00534390` and `FUN_00534ac0` directly and **never enters the
dispatcher**.

**That is a fact about the routines, not about the behaviours.** A player-only routine can
still reach a shared behaviour one level down, and patrol does: `FUN_00535300` has one caller,
but it tail-calls `FUN_005301f0`, which the script's own group-command dispatcher also calls
(`AI-PATROL-017`). The **pick-up** is the one that survives the test as a behaviour — the
complete immediate-form writer set of `actor+0x50 = 2` is three sites, and the two that reach
an actor are both inside `FUN_004d5dd8`.

The received object-target command (`cmd+4 == 3`) resolves signed-word id `cmd+0x0e`
against the unit collection first, then the structure collection. Opcode `0x19` forwards
that pointer to the shared attack order. The scorer must not return `0xffffff`, and the
target must not be the member itself, before the member retains it at `ord+0x0c`.
The physical branch proceeds through order 5, shared approach and action 3, copying the
pointer to `actor+0x5c`; no separate structure action or per-strike footprint lookup is
present in these bodies. This is a conditional received-command route, not proof that
every Building instance is safely admitted or that ordinary hover produces the command:
the scorer/AI still carry actor-oriented field reads on the supplied object.
Claim: `UNIT-STRUCTORDER-062`.

**Customisation.** The four empty order slots are absent implementations and free. The order
opcode's ceiling `0x26` is a compiled constant. Neither `actor+0x50` nor `actor+0x54` is serialized,
so extending either state space costs nothing in a shipped save — but `ord+0x08` and `ord+0x09`
*are*, inside a raw `0x94`-byte block, so a new field before them breaks every save the original
wrote while a new *value* in them does not.

## Formation — a per-player mode, and the only thing that makes a group move as a group

Every `Player` owns a `0x20`-byte settings block, allocated and constructed by the `Player`
constructor (`004fac7d PUSH 0x20` → `FUN_0052ccb0` → `004facb8 MOV [Player+0x30]`). Its **last
byte**, `+0x1f`, is the **formation mode**; the constructor writes the only nonzero default it
has, `2`. The block is serialized raw, `0x20` bytes, by `FUN_005392d0`, so the mode is in every
save.

| value | behaviour |
|---:|---|
| `0` | never in formation; the group centroid is not even computed |
| `2` | in formation **iff** every member is within `AImanager+0xa824` cells (Chebyshev) of the group centroid — the threshold is a compile-time `2` |
| anything else | in formation **unconditionally**; no spread test runs |

One routine writes it, `FUN_00537e90(player, mode)`, and it has two callers, both authored
surfaces:

- **player command `0x46`, sub-code 2** — the player from `cmd+0x05` through `FUN_004d4c6e`, the
  value from `cmd+0x0e` **remapped** `0→0`, `1→2`, `2→1`, anything else `→2`, so the wire can only
  produce `{0,1,2}`;
- **trigger instant id 7**, which the shipped `Description Instants.ini` names **`Set formation`**
  and which writes the parameter **raw** — so a map may author a value the UI cannot.

What the mode decides is both halves of a group Move / Swarm-2 order: whether each member is sent
to `target + (memberCell − centroidCell)` or every member to the same cell, **and** whether the
group's rate override is written at all. `formats/move/format.md` owns that; the gate is the same
local flag in both cases.

Shipped content authors it **once** — one `Set formation` node in the whole 38-map corpus,
`scn:110.alm`, carrying `2`, the constructor's own default.

## Order progress, and how an order becomes an act

Between the state machine and the actor tick sits `FUN_005310e0`, called once per actor per tick at
`004f398c` — four instructions before the tick dispatches on `actor+0x54`. It preserves
`actor+0x54` only when it is `2` or `0xf` and clears it otherwise, then dispatches on the order
object's **progress** byte `ord+0x09` (255 in range, four live) and, when that is 0, on its
**pending order** byte `ord+0x08` (15 slots, three empty). "When that is 0" means, in practice, the
tick the actor is standing exactly on a cell centre: the state machine re-decides the order at group
rate, the order machine executes it at cell rate.

| `ord+0x08` | what it does | inputs |
|---|---|---|
| `1` | walk to `ord+0x0a` | the cell |
| `2` | attack `ord+0x0c` **with no distance test** | the target |
| `4` | close on the actor `ord+0x18` | stop distance `ord+0x14` |
| `5` | pursue and attack `ord+0x0c` | facing, reach `actor+0x12c`, stop distance `ord+0x14` |
| `6` | as `5`, plus the auto-engage suppression, and turn instead of path | same |
| `7` | the sack pick-up bridge | — |
| `8` | cast `ord+0x30` at the actor `ord+0x28` | facing, `ord+0x14` |
| `9` | cast at the cell `ord+0x3c` | facing, `ord+0x14` |
| `0xa` | turn in place until `mover+0x01 == mover+0x00` | the desired facing |
| `0xb` | idle: re-face to `facing + 0x21 + (190·rand()/0x8000)`, a near-uniform byte, and only when `ord+0x54` is set or on `rand() < 0xcd` | — |
| `0xc` | walk to `ord+0x0a`, then the `actor+0x54 = 2` action | the cell |
| `0xf` | reach `ord+0x0a` within `ord+0x14`, then `actor+0x54 = 0xf` | the cell, `ord+0x14` |
| `3`, `0xd`, `0xe` | empty — the switch's own default | — |

The in-position test that `5`, `6` and `8` share is **not** a centre-to-centre distance: it is the
actor's current facing equal to the 8-way direction to the target, **and** an edge-to-edge distance
that subtracts both token sizes, so two touching actors measure 1. Reach is re-read from
`actor+0x12c` at the test; `ord+0x14` is used only as the mover's stop distance, inside which the
actor turns to face and stands still.

Worked example — the sack pick-up, which is the whole chain in one order:

1. order `0x21` → `FUN_005308a0`: `actor+0x50 = 2`, `ord+0x0a` = the sack's cell, `ord+0x08 = 0`.
2. `FUN_0052ce50` arm 2: not at `ord+0x0a` → `ord+0x08 = 1`, walk. **At it → `ord+0x08 = 7`.**
3. `FUN_005310e0` pending arm 7 → **`actor+0x54 = 2`**.
4. the actor tick's `actor+0x54 == 2` arm loots the sack under the actor's feet.
5. next tick, arm 7 sees `+0x54` already 2 and completes: `actor+0x50 = 0xc`, `ord+0x08 = 0`.

So a unit told to pick something up ends the order **hunting**, not idle. The same relay carries
every other order that has a "walk there, then do a thing" shape.

## Where a behaviour is chosen

Nothing in the `.alm` authors a behaviour **per unit**; everything below names a **group**.

### At load

The map contributes two things: the **group partition** (the type-6 group id, which puts units of
one player into one group) and the **roster order**. Two setters exist:

| setter | members | post it anchors | group order |
|---|---|---|---|
| guard | `actor+0x50 = 0xb` | current cell, or the cell being stepped into | 1 |
| aggressive / stand ground | `actor+0x50 = 0xc` | current cell, always | 3 |

Both also clear each member's order byte, so a stance issued mid-order cancels it.

On the map-load path, **in single player only**, the groups of the player in **type-5 slot 0** are
set aggressive and every other player's groups are set to guard. On all 38 maps the **EN** root
ships that slot is named `Self` (32 `Self`, 6 `self`). The single-player test and the argument test
gate **both** branches together: fail either and no group is set, so a multiplayer load leaves
every group at the constructor's `0` and hands every member to its own state machine — where the
per-actor guard arm anchors the post instead, at the first tick.

So the constructor's `0` — the value that hands a group to its members' own states — **survives no
map**: over the whole corpus the load-time byte is 1 for **8085 of 8094** EN placements and 3 for
the other 9, and 0 for none (`AI-CENSUS-046`). Slot 0 owns placements on only four campaign maps.

The mission description file can override this per group: key `Patrol` names a waypoint section, and
key `StandGround = Yes` selects the aggressive setter. An absent key means guard. No such file ships.

### At run time — the map's own script

**This is what shipped maps use, and a consumer that implements only the load-time setters will play
a different game.** The type-7 script's action opcode 6 is a group command dispatching on its own
first parameter. Eleven signatures are authorable, ten implemented, and nine appear in shipped
content — **135 nodes over 20 campaign maps, and none on any of the ten loose multiplayer maps**:

| par0 | catalogue name | shipped nodes | group order it leaves |
|---|---|---:|---|
| 1 | Guard | 1 | 1, via the load-time guard setter |
| 2 | Swarm | 6 | 2, with the target cell in `grpAI+0x0a` |
| 3 | Stand Ground | 24 | 3 |
| 4 | Move | 29 | 4 |
| 5 | Swarm 2 | 40 | 5 |
| 10 | Attack | 5 | 0 — per member, engage the named unit |
| 11 | Defend | 6 | 0 — per member, the follow states |
| 14 | Patrol | 14 | 0 — per member, state `0xa` (above) |
| 15 | Follow | 10 | 0 |
| 17 | Roam | 0 | `0x11` |
| 18 | Dwell | 0 | **no case — authorable and inert** |

A subcommand-10 attack is not a group-wide copy of one order. The helper first clears each
member's route and reservation state, then forks on identity: the named member enters state `0xc`
with its post set to its own cell and no pending order; every other member enters state 3 with the
named member at `ord+0x0c`, its own `actor+0x12c` at `ord+0x14`, and no pending order. The next
member-state pass therefore acquires on the subject and engages on everyone else
(`AI-SCRIPTATTACK-120`).

Subcommand 11 is the defend form already specified by `AI-DEFEND-111` and
`AI-FOLLOWSET-116`: the named member acquires, every other member stores it at `ord+0x10` and
defends it at the authored range. This closes the two per-member helpers formerly left unread;
subcommand 15 remains the related follow form, which keeps the range but omits defending on the
subject's behalf.

A group's order is therefore **not** its load-time 1 or 3 for the whole mission: it is whatever the
last command left. The four commands that leave 0 are the ones that hand the group back to its
members' own states.

### Which break-off rule that puts in front of a player

A creature at order **0** breaks off over a 5-cell Chebyshev block around its own **post**; one at
order **1** over `grpAI+0x2d` about the group **centre**; one at order **3** not at all. The census
(`AI-CENSUS-047`) settles which a player meets. The 135 command nodes name **91 of 2264 groups** and
**361 of 8094 placements**, and **45 of them are named by no trigger's action slot**, so they cannot
be dispatched at all. Between those two bounds:

| group order | break-off rule | placements, every node fires | placements, uncarried nodes dropped |
|---|---|---:|---:|
| 0 | post, Chebyshev 5 about `ord+0x00` | 68 (0.8 %) | 34 (0.4 %) |
| 1 | **centre, `grpAI+0x2d` about `grpAI+0x28`** | **7739 (95.6 %)** | **7819 (96.6 %)** |
| 2 | swarm — no clip published | 37 | 37 |
| 3 | none | 106 | 92 |
| 4 | move — no clip published | 10 | 14 |
| 5 | swarm 2 | 88 | 86 |
| — | commanded to more than one order | 46 | 12 |

Restricted to creatures hostile to slot 0 — the only ones that can pursue the player — the centre
rule takes **7293 of 7571 (96.3 %)** and the post rule **52 (0.7 %)**. **A consumer that implements
only the per-actor post leash implements the rule 1 % of the corpus uses.** The post rule is what a
player's *own* units run, because every player-issued order allocates a fresh group at order 0.

Eight of the 28 campaign maps carry no group command at all (`20 30 31 41 51 61 101 121`); on
mission 10, the first, 32 of 35 placements stay under the centre rule (`AI-CENSUS-048`).

## Withdrawing — what a ranged monster does when you close to melee

Once per full tick, for **every member of every group and after whatever the order arm did**, the
dispatcher runs two tests in this order and takes the first that fires:

| test | radius | threshold |
|---|---|---|
| 1 | the actor's own `pth+0x08` | `health <= ord+0x44` — the `Data.bin` Units column **`Wimpy`** |
| 2 | the literal **2** cells | `health <= ord+0x40` — the `Data.bin` Units column **`Withdraw`** |

Both collect the occupants of the square of that radius, keep the ones the diplomacy matrix calls
hostile, and require the set to be non-empty. Both thresholds compare against **current** health as
a signed word, so a value of 0 means *never* and a value equal to `healthMax` means *always*.

The first reaction helper counts entries with `health > 0` in a byte. A zero count calls
ordinary acquisition. Its averaging loop then sums fine coordinates over the **whole filtered
list**, not a living-only subset, and divides by the low byte of the list count. The second
helper has a nonempty-list gate without a positive-HP scan. Their actual mixed/all-dead
runtime populations and dense-count wrap remain Unknown (`AI-RETREAT-274`).

Both call the same distance-3 away picker and install pending move 1. The picker clamps to
`[8, dimension-9]`; a zero coordinate difference becomes 1. It has no passability test.
Route search, not the picker, handles inaccessible cells. These geometry clauses survive
the partial correction to `AI-WITHDRAW-028`.

**Shipped threshold source, not explicit-command eligibility.** The **Humans** table has
neither column; that is not an immunity to explicit Retreat. Of the 56 parameterised Units rows in
both shipped roots, twelve carry a nonzero `Withdraw` and they are exactly the three ranged
families — `Goblin_Sling` (30), `Orc_Bow` (60), `Bat_Sonic` (15), four tiers each. Because the
threshold is absolute health and the four tiers of a family share the tier-1 value while `healthMax`
scales ×1.6 per tier, a **tier-1 ranged monster withdraws on every contact** and tiers 2/3/4 only
below 62 % / 39 % / 24 % of their own maximum. Twenty-four rows carry a nonzero `Wimpy`
(`Goblin_Pike` 10, `Goblin_Sling` 8, `Bee` 3, `Squirrel` 5, `Bat_Sonic` 4, `Dragon` 63).

**Two things a consumer must not implement.** The image contains a spawn-time classifier that would
replace both numbers with `healthMax`, `healthMax/2` or `healthMax/4` according to a shooter/caster
class byte. It **never executes those writes**: its only call site is inside the base actor
constructor, before the derived vtable is installed, so its `is-humanoid` gate always reads 0.
Implementing it "correctly" discards the shipped table. The same class byte selects the per-class
percentages of the runtime command that resets both thresholds (three modes: off, 10 %, 30 %); since
nothing ever sets the byte, only the first column of that table is reachable.

## Difficulty

The three-way pre-create control does **not** reach behaviour. No instruction in the AI module reads
the state it sets; `ai.reg` has no per-level dimension; group counts come from the map, the cadence
from the speed index, and the radii from geometry. The level's only consumer is the spawner's
health / to-hit / defence adjustment (`formats/unit`).

## Cost

The candidate sweep is **linear in the number of on-map actors**, twice: once to stamp
visibility, once to test each actor's cell. There is no spatial index anywhere on the path, and
the visibility map is cleared and rebuilt per decision.
