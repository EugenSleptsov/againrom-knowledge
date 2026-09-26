# ITEM — the item, its container and the sack

Claims about the simulation's carrying layer, which is not a file format: what
an item is as an object, the container a unit and a sack both hold one of, and
the commands that move an item between them. Definitions come from
`world.res:data/data.bin` ([`databin.md`](databin.md)) and the stream from
[`sav.md`](sav.md). [`shop.md`](shop.md) owns the generation and pricing of the
same objects and published one of these commands (`SHOP-TRAY-025`). This ledger
is not [`inv.md`](inv.md), the install corpus and signatures ledger
(`INV-CORPUS-001`, `INV-SIG-002`, …) about files on disk; no game-item claim
belongs there. Format of this file: [registry.md](registry.md).

## Terms

The vocabulary is the engine's own.

- An item is any instance of the four `Item`-rooted classes.
- The container is the unnamed `CObList` subclass a unit holds at `actor+0x7c`
  and a sack at `sack+0x40`. Its `+0x20` field is the load.
- A sack is the `Sack` class.
- Equipped always means one of the fourteen pointer fields of `ITEM-EQUIP-006`,
  never a member of the container.

## Consumable use

| ID | Claim | Confidence | Status | Evidence |
|---|---|---|---|---|
| ITEM-USE-112 | Carried consumables enter the shared character-panel action from two physical inventory paths, not inventory right-click. | High / Unknown | ✔ promoted | [EXP-0275](../experiments/EXP-0275-consumable-use/) |
| ITEM-USE-113 | Potion use commits through `Item::Equip`; cast consumables reserve one carried object before release. | High / Unknown | ✔ promoted | [EXP-0275](../experiments/EXP-0275-consumable-use/) |
| ITEM-USE-114 | A cast item's reservation, cast-start and destruction are different boundaries; effect success is not the refund predicate. | High / Unknown | ✔ promoted | [EXP-0275](../experiments/EXP-0275-consumable-use/) |
| ITEM-VALUE-115 | Book's general value reader takes its first Effect, not its last; `Item+40` is not the shelf spell id. | High | ✔ promoted | [EXP-0275](../experiments/EXP-0275-consumable-use/) |

### ITEM-USE-112

- Mission-grid vtable `59a1a0+5c → 4835b0`: a double-click selects one hit item
  and calls `session+e0` slot `+7c`. Panel left-up `491330` calls that slot for
  held source/place codes 1/2.
- Shop-backpack constructor `4a50a0` writes vtable `59ab08`. Its
  `+5c → 4a39c0` carried-source-2 branch calls the borrowed panel `view+7c`,
  again `+7c → 492410`. Activation `4a8bc3` takes the panel from `session+e0`.
- The action requires one selected member, a held item, matching player
  ownership and `460440` descriptor/class admission; ordinary Potion and Scroll
  admit both classes. The slot-0/1 member-type and teachSpell gates remain
  `ITEM-WEAR-057`.
- Without display bit `0x10`, usable bit 1 enters a one-unit `41c98f` transfer
  to destination 1.
- With both bits, `4b0d40` stores slot and castSpell id-minus-one, and
  `41b439(10)` forces mode 5 without the ordinary learned-spell command-mask
  gate. This also arms from the shop; it is not evidence of a town target-map
  release.

**Confidence.** High for the complete physical handlers, the constructor/vtable
connection and the gates.

**Unknown.** Release while a town/shop covers the mission map.

### ITEM-USE-113

- Session `0x22` resolves the actor, removes the requested source quantity,
  calls actor slot `+38` for destination 1, and returns any non-null result to
  its container.
- Human `4f705b` and base `4f4d98` forward to Item `508779`. For signed HP
  `>0` it walks every Effect in order, changes mode 0 to singleuse 8,
  attaches/applies each, recomputes actor Token value through vt+54 and
  destroys the detached Item. HP `<=0` returns it unchanged.
- No test requires a stat increase, so full-health healing or a
  class-inapplicable effect can still consume one. The ordinary UI requests one
  even from a stack.
- In order `0x25/26`, `4d6466` removes one via `50ebf9` and demands the first
  Effect kind 41. It sets actor+68 to the item, constructs a 0x14-byte Spell
  from the low byte of Effect+40 at actor+44, and builds initial values from
  the low byte of Effect+42.
- The cell arm restores at the original slot and deletes the temporary Spell if
  `4fe0d2` finds Spells target parameter 4 !=2.
- A missing or wrong first Effect exits after detachment without that restore;
  ordinary descriptor arming does not make that malformed command reachable.

**Confidence.** High for the complete dispatcher, wrappers, leaf and return
branches.

**Unknown.** Whether the malformed command is reachable in ordinary play.

### ITEM-USE-114

- Physical map-up `419ec1` uses the item target table through `4b0e20`: a
  unit-like hit and a nonzero flag emit `0x25` through `41c2d6`; otherwise an
  admitted point emits `0x26` through `41c0dc`.
- Producers require selected members with +7c!=0, but item context bypasses
  their learned-spell bit.
- The order dispatcher requires a resolved actor and AI context and cancels
  earlier item context in its prologue; then `533d00/533ed0` move the Spell
  into order+30.
- Progress `5310e0` waits for range/approach before act `0x0d/0x0e`. `4f37be`
  calls `4fe6d3`, whose mana test is bypassed when actor+68 is non-null;
  acceptance sets actor+136=0.
- At completion it invokes apply and destroys the item and Spell for packed
  Item category 14, independently of effect damage/attachment success.
- Cancellation `4f4c97` restores non-Weapon context only when its first Effect
  is kind 41, act is outside `0x0d/0x0e` or +136!=0, and the active Spell id
  matches the Effect signed id.
- `4f4b39` first tries that refund, then destroys unrefunded category 14
  context when +136==0. Destructor `4f3432` also destroys retained category 14
  context.
- UI mode reset `4b0d90` deletes a display copy, not the server Item.

**Confidence.** High for these branches and the direct caller census.

**Unknown.** Every pathfinding-failure and stale-target outcome, and wall-clock
timing.

### ITEM-VALUE-115

- `508486` Book arm `5085d9..508634` gets the list head once, reads signed
  Effect+40, assigns Spells parameter 21 to Item+1c and exits. An empty list
  means zero. It neither loops nor checks kind 42 there.
- Scroll's separate branch loops kind 41 effects and sums its spell/power
  formula, falling back to MagicItems parameter 0 only on a zero sum.
- The live shelf overrides the stored price after construction
  (`SHOP-CONSUME-073`).
- `5081f1 → 525d60` builds Item+40 from category 14, shape, material and
  definition index. The newly allocated Effect owns spell id +40 and power +42.
- `508340` copy preserves Item fields and each Effect without calling the value
  reader again.

**Confidence.** High. One head read and no Book back edge discriminate first
from last; connected constructors distinguish the two +40 owners.

## Item object, definition, stack, container and equipment

| ID | Claim | Confidence | Status | Evidence |
|---|---|---|---|---|
| ITEM-CLASS-001 | An item is a `Token`, in the same family as a unit and a building: four concrete classes, and the record is 0x50 bytes at its narrowest. | High | ● active | [EXP-0079](../experiments/EXP-0079-items-and-sacks/) |
| ITEM-DEF-002 | An item's definition is a row of `Data.bin`, and the C++ class picks the table. | High / Unknown | ● active (partially retracted) | [EXP-0079](../experiments/EXP-0079-items-and-sacks/), [EXP-0287](../experiments/EXP-0287-actor-definition-binding/), [EXP-0396](../experiments/EXP-0396-generated-mission-load-crash/) |
| ITEM-STACK-003 | A stack is one object with a count, a per-unit weight, and a virtual that splits it. | High / Unknown | ● active (amended, partially retracted) | [EXP-0079](../experiments/EXP-0079-items-and-sacks/), [EXP-0225](../experiments/EXP-0225-item-effect-grammar/) |
| ITEM-CONT-004 | The thing that holds items is one class, a `CObList` with two dwords bolted on: no slot count, no capacity, and a unit and a sack hold the same one. | High / Medium | ● active | [EXP-0079](../experiments/EXP-0079-items-and-sacks/) |
| ITEM-LOAD-005 | Carried weight does one thing, a speed penalty, and never refuses: the derive adds half the container's load to the actor's own weight. | High | ● active | [EXP-0079](../experiments/EXP-0079-items-and-sacks/) |
| ITEM-EQUIP-006 | Equipped is a different place from carried: fourteen pointer fields, not a marker in the list. | High / Medium / Unknown | ● active (amended, partially retracted) | [EXP-0079](../experiments/EXP-0079-items-and-sacks/) |

### ITEM-CLASS-001

- `EXP-0057`'s `CRuntimeClass` table names `Item 0x50 : Token`, with
  `Armor 0x68`, `Shield 0x68` and `Weapon 0x84` derived from `Item`. Each
  thunk's `PUSH <size>` is that size: `00507ea9` and the three sibling thunks
  `FUN_0050c1ca` / `FUN_0050cbb9` / `FUN_0050d51d`.
- `Item::Item` is `FUN_00507f42`: `Token::Token` (`004f241e`), vptr
  `0x59c8e0`, then the eight fields it writes:
  - `+0x3c = 0` (the definition pointer);
  - `+0x40 = 0` (u16);
  - `+0x42 = 1` (u16, the stack count);
  - `+0x44 = +0x45 = +0x46 = 0` (u8);
  - `+0x0c = 0` (u8, the definition row index);
  - `+0x4a = 1` (u16, the weight);
  - `+0x08 = 0` (dword, a flag word);
  - `+0x4c = 0`.
- Vtables and their extents, from `FindVtables 59c000:59d000 6`:
  `Item 0x59c8e0` has 27 slots. Three siblings follow in one run with no null
  between them: `Armor 0x59c950`, `Shield 0x59c9a8`, `Weapon 0x59ca00`. Each is
  identified by its own `GetRuntimeClass`
  (`MOV EAX,0x5c33a0 / 0x5c33b8 / 0x5c33d0`), the three records `refto:` finds
  4 / 3 / 3 references to, 0 orphan.
- The vtable slots this ledger uses: `+0x08` Serialize, `+0x38` Equip, `+0x3c`
  Unequip, `+0x40` detach-one, `+0x50` IsStackable.
- `Item`'s own `+0x3c` is `FUN_00508836`, which prints
  `"Unknown item takeoff"` and does nothing. The base class cannot be taken
  off; only the three subclasses override it.

**Confidence.** High. The sizes are each an allocator immediate beside a
`CRuntimeClass` record that names them; the four vtable bases are the
constructors' own `MOV [this],<addr>`; the extents are a run scan whose
boundaries are confirmed by the three `GetRuntimeClass` bodies landing at slot
0 of each.

### ITEM-DEF-002

- Each class's `Serialize` load arm resolves `+0x3c` from the object's own
  `+0x0c` against a different collection base: `Item → 0x609b7c` (`0051171f`),
  `Armor → 0x609b54` (`00511b2d`), `Shield → 0x609b40` (`00511019`),
  `Weapon → 0x609b68` (`00510502`).
- Against `DAT-OBJ-002`'s map
  `0x609b18 + {0x28 Shields, 0x3c Armors, 0x50 Weapons, 0x64 MagicItems}` those
  are Magic Items, Armors, Shields and Weapons, the four collections EXP-0049
  framed and left uninterpreted.
- Only class `Item` compares the row against the collection size: `GetSize`
  (`FUN_0051b870`), `CMP` against `+0x0c`, then `ElementAt(+0x0c) → +0x3c`, or
  `+0x3c = 0` when the index is out of range. `Armor`, `Shield` and `Weapon`
  call `FUN_0051b850`, the unchecked `ElementAt`, not a `GetSize`, and never
  take a `+0x3c = 0` branch (`SAV-1088`).
- `item+0x0c` is a u8 row index and `item+0x3c` the resolved entry pointer,
  re-derived at every load rather than stored.
- Human is not the identical bounded lookup: its load selects the saved row
  only for `u16 +0x0e < 33`, otherwise literal row5, through unchecked stride48
  arithmetic (`SAV-ACTORBIND-544`).
- The class is the table: there is no kind field selecting it, and an `Armor`
  can never name a Weapons row.

**Confidence.** High for the four table identities: each base is a literal in
its class's load routine, and the collection→offset map is `DAT-OBJ-002`'s,
carried independently by the CSV parser and the DB Serialize. The distinct
Human exception is bounded by `SAV-ACTORBIND-544`. The grade first rested on
four instances of one instruction sequence; three of the four lack the compare.

**Unknown.** Which column of those rows means what: `claims/databin.md`'s open
item, untouched here.

**Amended.** Two clauses are withdrawn ([`retracted.md`](retracted.md)).
EXP-0287 (`SAV-ACTORBIND-544`) refuted the Human-as-identical-bounded-lookup
clause; the Human bullet above replaces it. EXP-0396 (`SAV-1088`) partially
retracted the shared size-checked lookup: the claim had stated that all four
arms end with `GetSize` (`FUN_0051b870`/`FUN_0051b850`), `CMP` and
`ElementAt(+0x0c) → +0x3c` or `+0x3c = 0`; only class `Item` does.

### ITEM-STACK-003

- `item+0x42` (u16) is the count, default 1; `item+0x4a` (s16) is the per-unit
  weight, default 1.
- The container's running load is `Σ (s16)+0x4a × (u16)+0x42`, signed times
  unsigned (`MOVSX` / `MOVZX` / `IMUL`), at `0050e983`, `0050ea2f`,
  `0050ea8a`, `0050ec95` and `0050ebce`.
- Taking `qty` from a stack of `n > qty` is `count -= (qty−1)`, then `vt+0x40`
  detaches one unit, then the detached object's count is set to `qty`
  (`0050ec57`…`0050ec80`). `n <= qty` removes the whole element instead.
- Merging is the inverse and also ORs the two `+0x08` flag words together
  (`0050e9bf`).
- Stackable is a virtual, `vt+0x50`. `FUN_00508635`, inherited by all four item
  classes, is `item+0x44 == 3 || effects.empty`.
- `FUN_005081f1` writes literal 3 when a `MagicItems` row name begins `Potion`,
  4 for `Scroll` and 5 for `Book`, and otherwise leaves 0.
- `FUN_00508671` first compares item code. It returns equal immediately when
  both objects are stackable, unequal when exactly one is, and compares ordered
  effect lists only when neither is. A Potion therefore stays stackable with
  effects, and equal-code Potions compare equal without examining those
  effects.
- Enchantment separates equal-code items only when both are non-stackable.

**Confidence.** High: every arithmetic step is named; the four load-update sites
are complete within the container class; both stackability tests, the prefix
writer and the equality control flow are complete routines.

**Unknown.** Any `+0x44` producer beyond the bounded prefix family and the item
constructors.

**Amended.** EXP-0225 refuted the former final clause, "what separates two
otherwise identical items is their enchantment", which made enchantment a
universal separator; Potions are the explicit exception, and the last three
bullets replace it ([`retracted.md`](retracted.md)). Count, weight, split,
merge and the two `IsStackable` tests stand.

### ITEM-CONT-004

- `FUN_0050e7d3` is the whole constructor: `CObList::CObList` on `this`
  (`00526600`), `+0x1c = 0x2710` (10 000), `+0x20 = 0`.
- The object is `0x24` bytes: `PUSH 0x24` at all fifteen construction sites,
  and `callto:50e7d3` = 15 hits / 14 owners / 0 orphan.
- The base actor constructor parks one at `actor+0x7c` (`004f3282`), the
  `Sack` constructors at `sack+0x40` (`0050f2a3`, `0050f33b`, `0050f39e`).
- `+0x1c` is not a capacity: it is the index the next `Add` inserts at.
  `FUN_0050e8e2` pushes it straight into `FUN_0050e902`, which `AddTail`s
  whenever the index is `>= GetCount`. Both the pick-up order (`004d6791`) and
  the equip arm (`004d692c`) rewrite it before an add.
- `+0x20` is the load.
- No routine in the class compares either field against a limit, and neither
  insert path (`FUN_0050e902`, `FUN_0050e8e2`) nor any arm of the move command
  reads a maximum. An inventory is unbounded in both count and weight.

**Confidence.** High for the class shape and the two fields' arithmetic: the
constructor is 6 instructions, and the insert, the take and the move-all were
read whole. Medium for "nothing refuses". The instrument is `EnumRefs callto:`
on the two inserts, 5 hits / 3 owners and 31 hits / 25 owners, 0 orphan; the
item-module and command-dispatcher owners were read and the rest classified by
owner, not read line by line. Its blind spot is a caller that tests the load
before calling.

### ITEM-LOAD-005

- In the derive `FUN_004f7dfc`, `actor+0x90` (the load) is set to `actor+0x8e`
  (the actor's own weight) at `004f81b3`.
- If `actor+0x7c` is non-null, half the container's `+0x20` is then added
  (`CDQ / SUB / SAR 1` at `004f81e2`, truncating toward zero), unless
  `+0x20 >= 0xfa00` (64 000); then the load is assigned the flat `0x7d00`
  (32 000) at `004f8203`.
- The only consumer is `HERO-SPEED-008`'s penalty: `004f8220` compares the load
  against `actor+0x92` (capacity, `body×10+1`) and, only when
  `load >= capacity`, subtracts `load/capacity` from speed and floors the
  result at 6.
- `HERO-SIGHT-007` calls `[actor+0x7c]+0x20` "the carried money". It is the
  inventory container's running weight sum, maintained by the four
  `weight × count` sites of `ITEM-STACK-003`. Money is `Player+0x38` and is
  never an item (`ITEM-DROP-008`). That claim's arithmetic is untouched; its
  label is corrected here and in [`retracted.md`](retracted.md).

**Confidence.** High: the block is nine consecutive named instructions, and the
field it reads has its four writers enumerated inside one class.

### ITEM-EQUIP-006

- Two fields live on the base actor, `actor+0x74` and `actor+0x78`. Twelve more
  form an array `actor+0x198 + 4i` that only a `Humanoid` has (`Unit` is
  `0x198` bytes, `Humanoid`/`Human` `0x1e8`).
- The take-off arm of command `0x22` computes `slot = (s16)cmd+0x0e + 1` and
  dispatches: `slot == 1 → actor+0x74` (`004d6872`),
  `slot == 2 → actor+0x78` (`004d688d`),
  `3 <= slot <= 12 → [actor + slot*4 + 0x198]` (`004d68c0`).
- The last arm sits behind `actor->vt+0x30()`, whose false arm prints the
  engine's own `"Error - Trying to takeoff armor from non humanoid"`
  (`0x5c63cc`). The base actor's `vt+0x30` returns 0 (`FUN_005232d0`, three
  instructions), which makes those ten slots humanoid-only.
- Unit equip/removal wrappers `004f4d98/004f4e0a` forward to item
  `vt+0x38/+0x3c`, then call actor `vt+0x54`, whose Unit target only returns
  `actor+0x1c`. Humanoid/Human use `004f705b/004f70cf` and omit that trailing
  call. The separate derive slot is `+0x50`; `+0x54` is Token-value
  computation, not equipment derive (`ITEM-ARMFOLD-033`, `SAV-EQUIPCALL-554`).
- For accepted Armor/Shield/Weapon equip, the returned pointer is the same-slot
  displaced item or null. The command reinserts a nonnull return at the source
  index (`004d6959`), then calls load refresh with zero (`004d6963`).
- `Item`'s own `vt+0x38` (`FUN_00508779`) is not a slot at all. It refuses when
  `actor+0x94 <= 0` (dead); otherwise it attaches every effect of the item to
  the actor and deletes the item: a consumable being drunk.

**Confidence.** High for the slot map and the gate: each is a named
instruction, and the engine supplies the string that names the array as
armour. Medium for "no per-slot type restriction": none is applied anywhere on
the take-off path, but the per-class `vt+0x38` bodies of
`Armor`/`Shield`/`Weapon` were not all read, and a restriction could live
there.

**Unknown.** `actor+0x19c` and `actor+0x1a0`, array indices 1 and 2, which
`Humanoid::Serialize` writes and the take-off arm never addresses.

**Amended.** EXP-0288 (`SAV-EQUIPCALL-554`) refuted the universal actor
wrapper and recompute clauses, "Actor equip is `004f4d98` followed by recompute
through `+54`; removal similarly recomputes" ([`retracted.md`](retracted.md)).
The wrapper bullet above states the correction. `ITEM-HUMEQ-030` reads the
three per-class `Equip` bodies the Medium clause left unread, and
`ITEM-ARMSLOT-031` addresses the Unknown through the `Armors` `Slot` column.

## Moving items, drops, sacks and pick-ups

| ID | Claim | Confidence | Status | Evidence |
|---|---|---|---|---|
| ITEM-CMD-007 | Moving an item is one command, `0x22`, with a source code and a destination code, and the ground is one of them. | High / Unknown | ● active (amended, partially retracted) | [EXP-0079](../experiments/EXP-0079-items-and-sacks/), [EXP-0275](../experiments/EXP-0275-consumable-use/) |
| ITEM-DROP-008 | A drop lands where asked only within a Chebyshev distance of 2, otherwise at the dropper's own feet, and is never refused. | High | ● active | [EXP-0079](../experiments/EXP-0079-items-and-sacks/) |
| ITEM-PICK-009 | A pick-up is all-or-nothing: the whole sack, every item and all the gold, in one act, and nothing refuses it. | High / Unknown | ● active (amended, superseded) | [EXP-0079](../experiments/EXP-0079-items-and-sacks/) |
| ITEM-SACK-010 | A sack is a `Token`, not an actor, an item or a building, and there is one per cell. | High | ● active (contested) | [EXP-0079](../experiments/EXP-0079-items-and-sacks/) |
| ITEM-SACK-011 | A sack does not tick and does not expire; the engine has no such mechanism. | High / Medium | ● active (contested) | [EXP-0079](../experiments/EXP-0079-items-and-sacks/) |
| ITEM-DEATH-012 | Death passes the corpse's container to the sack merge/adopt wrapper, unless the template is an NPC, in which case the inventory is destroyed. | High / Unknown | ● active (amended, partially retracted) | [EXP-0079](../experiments/EXP-0079-items-and-sacks/), [EXP-0134](../experiments/EXP-0134-corpse-worn-armour/), [EXP-0285](../experiments/EXP-0285-whole-item-transfer/) |
| ITEM-SPAWN-013 | Map load calls two further sack makers, a mission `.ini` `"Items"` reader and a random treasure scatter; the claim that no sack comes from the `.alm` is retracted. | High / Medium | ● active (partially retracted) | [EXP-0079](../experiments/EXP-0079-items-and-sacks/) |
| ITEM-SAVE-014 | `Item::Serialize` writes the 37-byte Token head, the effect list and a fixed scalar tail; one container per actor and one per sack is serialized. | High | ● active (amended) | [EXP-0079](../experiments/EXP-0079-items-and-sacks/), [EXP-0225](../experiments/EXP-0225-item-effect-grammar/) |
| ITEM-CARRY-015 | A hero's equipment and backpack cross a mission boundary inside its carried object graph; nothing of a mercenary's crosses. | High | ● active | [EXP-0079](../experiments/EXP-0079-items-and-sacks/) |
| ITEM-PICK-016 | Picking up a sack is one execution path with two entry points: session source 3 for a sack underfoot and order opcode `0x21`, which walks to a named cell. | High | ● active | [EXP-0090](../experiments/EXP-0090-commands/) |

### ITEM-CMD-007

- The dispatcher `FUN_004d5dd8` has two switch spaces:
  - orders (`cmd+0x04 >= 1`): opcode−0x14 over 19 direct dwords at
    `0x4d86ae`;
  - session commands (`cmd+0x04 == 0`): opcode−2 through a byte index at
    `0x4d876e` into 29 dwords at `0x4d86fa`.
- `0x22` is in the second space, arm `004d666a`. The decode is confirmed by
  four opcodes published elsewhere landing on their known arms:
  `0x33`/`0x34`/`0x35` the shop's buy/sell/clear, `0x3f` `Shop::SetCap`,
  `0x38` the mercenary hire, `0xbe` the carried-character import.
- Codes, extending `SHOP-TRAY-025`'s list with the one it did not have:
  - source `1` = an equipment slot, `2` = the actor's container, `3` = the
    ground, `4` = a shop;
  - destination `1` = equip, `2` = the container, `3` = the ground, `4` = a
    shop;
  - `4..8` on both sides selects a shop shelf.
- Source 3 does not move anything. It checks that a sack sits on the actor's
  own cell (`FUN_00547c60`) and refuses with
  `"Invalid pickup order - no sack there."` if not. Otherwise it writes the
  destination index into `[actor+0x7c]+0x1c` and sets `actor+0x50 = 2`, a
  deferred order.
- Using a castSpell item takes order-space `0x25`/`0x26`; not every item use
  does. An ordinary potion reaches session `0x22` destination 1 and
  `Item::Equip`. The cast arm detaches one item, requires its first Effect kind
  `0x29`, and restores only the explicitly failed cell-target admission there;
  this is not a general failure refund. `ITEM-USE-113/114` separate
  reservation, release and cancellation.

**Confidence.** High: the two tables are dumped, the mapping is cross-checked
against four independently published opcodes, and every code test and every
store is a named instruction.

**Unknown.** How `actor+0x50 = 2` becomes the tick's state, on this claim's own
evidence (EXP-0079). EXP-0090 answers it (Amended).

**Amended.** EXP-0275 refuted the universal use and refund clauses, "Using an
item is a different opcode pair, 0x25/0x26" and "puts it back if the cast
fails"; the last bullet replaces them (`ITEM-USE-112/113/114`,
[`retracted.md`](retracted.md)). EXP-0090 (`ITEM-PICK-016`) answers the
Unknown: `FUN_0052ce50` arm 2 walks to `ord+0x0a` and writes `ord+0x08 = 7` on
arrival, `FUN_005310e0` arm 7 turns that into `actor+0x54 = 2`, and the tick's
arm 2 loots. Contradiction C-6 with `AI-STATE-011` (`retracted.md`, AMBIGUOUS)
is resolved in favour of both claims: they describe one execution path, entered
from opposite ends, and this claim's "deferred order" and that claim's "AI
state" are the same store of the same value into the same field. The entry
point this claim did not have is order-space opcode `0x21`, which names a cell
and walks to it.

### ITEM-DROP-008

- Destination 3 of command `0x22` allocates a fresh container (`004d6a4f`),
  puts the item in it, and unpacks `cmd+0x10` as `destX = low byte`,
  `destY = high byte`.
- Two independent tests follow: `abs(actorCol − destX)` and
  `abs(actorRow − destY)`, each `CMP EAX,0x2 / JG` (`004d6acf`, `004d6afb`),
  with no addition between them. The window is a 5×5 square, not a radius.
- Inside it the drop goes to a position built at `(destX,destY)`; outside it
  the same call is made with `actor+0x10`, the actor's own cell.
- Money is dropped by a different opcode, `0x23`, through the same geometry. It
  resolves a `Player` by id (`"Order error: no such Player "` on failure), takes
  that player's first actor as the dropper, refuses an amount `<= 0` or greater
  than `Player+0x38`, debits `Player+0x38`, and calls the same sack routine
  with a null container and the amount as gold.
- Gold on the ground is a field of the sack, not an item, and there is no coin
  object anywhere in the family.

**Confidence.** High: both windows are four named instructions in two arms of
one routine, and the money arm's three refusals and the single debit are named.

### ITEM-PICK-009

- `FUN_004f4e3e(sack)` credits `sack+0x3c` to the owner's `Player` money
  (`FUN_004faff7`) and stamps `+0x08 = 1` on every element. It pours the sack's
  entire container into `actor+0x7c` with `FUN_0050eaaf`, which moves one unit
  at a time and then destroys the source container. It then nulls `sack+0x40`,
  notifies the clients and deletes the sack.
- The routine has no capacity test, no distance test, no ownership test and no
  per-item selection. The protocol has no "take item *i* from the sack" at all,
  only source code 3, which names a destination slot and takes everything.
- Its two call sites (`callto:` on the repaired table, 2 hits / 2 owners /
  0 orphan) are a command handler `FUN_004d4e18` and the actor tick
  `FUN_004f37be`. The tick's arm at `004f3e76` first requires a sack at the
  actor's own cell, unregisters it and unlinks it before looting.
- That arm is selected by the tick's state dispatch on `actor+0x54 − 1`: state
  2 is the pick-up. The byte index at `0x4f41d9` has 15 entries
  (`004f39a3 CMP …,0xe`) over 6 dwords at `0x4f41c1`, so states `4..0xc` are on
  the default arm rather than out of range.

**Confidence.** High for the routine and the tick arm: both were read whole, and
the state→arm binding is the dumped table.

**Unknown.** The order writes `actor+0x50 = 2` but the tick dispatches on
`+0x54`; on this claim's own evidence, what carries `+0x50` into `+0x54` is not
established. `ITEM-PICK-016` closes it (Amended).

**Amended.** Superseded by EXP-0090 (`ITEM-PICK-016`, `AI-PROGRESS-034`;
[`retracted.md`](retracted.md)): the named discriminator, and the negative
behind it. The negative was a whole-image `EnumRefs re:` for
`dword ptr [reg + 0x54],<imm|reg>` that returned 198 hits over 99 owners with
9 in orphan or undisassembled code and no actor-module writer storing 2; it
enumerated one store form in one module and reported it as a fact about the
engine, the shape `PROCESS-LOG` names. The discriminator named was the write at
`0052cfd6`, with reading that arm as the whole remaining step. `0052cfd6` is
arm 0 of that table, which clears `+0x54`, and is adjacent to the `JMP` only
because it is the first arm.
The relay is two bytes wide and one routine away: `FUN_0052ce50` arm 2 writes
`ord+0x08 = 7` on arrival, and `FUN_005310e0` (`AI-PROGRESS-034`), called by
this tick four instructions earlier at `004f398c`, writes `actor+0x54 = 2` from
`EBX`. The 15-entry byte index in the last bullet corrects this claim's
narrower reading of the table.

### ITEM-SACK-010

- `Sack 0x44 : Token` (the `CRuntimeClass` at `0x5c33e8`), vtable `0x59ca88`,
  19 slots, the shortest in the family: it overrides its destructor and
  `Serialize` and nothing else.
- Its own fields:
  - `+0x04`, a runtime id drawn from `MOVE-ID-016`'s bitmap (`FUN_004d9fed` at
    `0050f3d4`);
  - `+0x3c`, gold;
  - `+0x40`, a container of the same class a unit holds;
  - `+0x1c`, the `Token` value slot that carries `XPvalue` on an actor and
    price on an item, recomputed by `FUN_0050f4b3` as `gold + Σ item+0x1c`.
- Three constructors: no-argument (the `CArchive` `CreateObject` thunk), by
  cell, and by cell and an existing container. The last takes the caller's
  container object by pointer (`0050f39e`) instead of allocating.
- All creation goes through `FUN_0050f5aa(pos, container, gold)`, which looks
  the cell up first (`FUN_00547c60`). When a sack is already there it pours the
  incoming container into it and adds the incoming gold rather than building a
  second. A new sack that fails to register at its cell (`FUN_005477b0`) is
  deleted and the call returns 0.
- A sack occupies a cell in the world's own sack registry `[0x005f22c8]`. It is
  not in the actor tick list `[0x00609558]`.

**Confidence.** High: the class record, the vtable extent, the three
constructors and the merge-or-create routine are all read at instruction level,
and the per-cell rule is the routine's own first act.

**Amended.** Contested by `MOVE-TICK-014` (contradiction C-2 in
`docs/REGISTRY-LOG.md`). That claim lists sack creation `FUN_004feadb` among the
callers of the tick-list insert `FUN_0050fc0a` and counts Sacks as members of
`[0x00609558]`. This claim reads sack registration as a direct id-allocator
call with no list insert (`ITEM-SACK-011`). Both readings stand.

### ITEM-SACK-011

- `Sack vt+0x14` and `vt+0x18`, the two slots `SESS-TICK-004` shows the loop
  invoking on every actor, are `FUN_00526b50` and `FUN_00526b60`: seven bytes
  each, prologue and `RET`.
- A sack never enters the loop. `FUN_0050f3c3` calls the id allocator
  `FUN_004d9fed` directly, without the `FUN_0050fc0a` AddTail every actor
  creator uses (`MOVE-TICK-014`), and `FUN_0050f715` appends it to the sack
  manager's own list instead.
- There is no timer field, no stage byte and no destroy-on-age path. A sack is
  created, merged into, and destroyed once, by a pick-up (`ITEM-PICK-009`).
- `SAV-OBJ-016`'s corpus progression "sacks 4/4/3/2" is therefore looting,
  which is how that claim's own Medium clause reads it; a gloss calling it
  decay is wrong.

**Confidence.** High that a sack does not tick: the two stubs are their whole
bodies, and the absence from the tick list is a positive reading of the sack's
own registration routine, not a failure to find something. Medium that no other
container ages it: the enumeration is `callto:` over the sack's destructor, its
registration and the create-or-merge routine (1, 1 and 6 hits, 0 orphan), which
cannot see a wholesale sweep of the manager's list that this experiment did not
identify.

**Amended.** Contested by `MOVE-TICK-014` (contradiction C-2 in
`docs/REGISTRY-LOG.md`), which names sack creation `FUN_004feadb` as a tick-list
inserter and counts Sacks as members, while this claim reads `FUN_0050f3c3`
calling the id allocator directly. `MOVE-TICK-014` itself also lists
`FUN_0050f3c3` among the direct allocator callers with no list insert. Both
readings stand.

### ITEM-DEATH-012

- `FUN_004f4f5d` has one caller, the actor tick `FUN_004f37be` at `004f392c`.
  It sets `actor+0x50 = actor+0x54 = 0x10` and leaves the world.
- It unequips `actor+0x78` into the container unconditionally, and unequips
  `actor+0x74` only if the weapon's own `Data.bin` parameter 15 is non-zero
  (`004f4fda` `PUSH 0xf` → `FUN_0051ab50`, `CMP dword ptr [EAX],0x0`).
- Four instructions, `004f5010`…`004f5018`, then dispatch `actor->vt+0x44`,
  which on both humanoid classes empties all twelve worn armour fields into the
  same container (`ITEM-CORPSE-034`). A person's corpse drops its armour too.
- It then computes a suppression flag: `CString::Find` of the literal `"NPC"`
  (`0x5c7d3c`) in the template's name at `[actor+0x3c]+4` (`004f5029`,
  `SETNZ`), OR'd with `Player+0x5c != 0` in multiplayer.
- When the flag is set, the whole container is deleted (`00523330` at
  `004f5078`) and replaced by an empty one before anything is dropped. A dead
  mercenary, whose templates are `MERC-LEVEL-005`'s `NPC%02d_%d`, leaves
  nothing.
- Gold is rolled only for `typeID > 0x40`, from parameters
  `0x26`/`0x27`/`0x28` (`HERO-KILL-027`'s three treasure columns).
- A sack is built by `FUN_004f4f30` → `FUN_0050f6d2(actor+0x10, container, gold)`
  iff the container is non-empty or gold is non-zero, and the corpse is then
  given a fresh empty container (`004f517b`).
- Only a newly created Sack receives the same container object; an existing
  Sack drains it through `0050eaaf` and deletes it (`ITEM-GROUNDMOVE-130`).
  Whole/split/merge Item identity depends on the helper predicates. Weapon
  unequip also changes its owned Spell (`ITEM-SPELLMOVE-132`).

**Confidence.** High for every clause stated here: each is a named instruction,
and the `"NPC"` literal is dumped from its own address. The "one routine read
whole" warrant that carried the grade is withdrawn: the routine held an unread
virtual call.

**Unknown.** The meaning of the weapon's parameter 15 was graded Unknown;
`ITEM-SUIT-035` closes it: its shipped title is `sutableFor`.

**Amended.** Two corrections ([`retracted.md`](retracted.md)). EXP-0134
(`ITEM-CORPSE-034`) refuted the completeness warrant "one routine read whole":
the published sequence stepped over `004f5010`…`004f5018`, the
`actor->vt+0x44` dispatch in the third bullet. EXP-0285 (`ITEM-GROUNDMOVE-130`)
refuted the unconditional container-identity sentence, "The sack receives the
same object the actor was carrying — identity, not a copy"; the last bullet
replaces it.

### ITEM-SPAWN-013

- `EnumRefs callto:50f5aa` on the repaired table: 6 hits, 5 distinct owners, 0
  in orphan or undisassembled code. This claim listed the owners as
  `FUN_0050f6d2` (the wrapper both drop arms and death use), the two arms of
  command `0x23`, `FUN_004f1fc6` and `FUN_004f2176`; `ITEM-SPAWN-026` holds the
  corrected membership.
- `FUN_004f1fc6` and `FUN_004f2176` are both called by `FUN_004d00e9`, the
  map-load routine `SESS-LOAD-009` read whole, and neither is in the `.alm`
  walker.
- `FUN_004f2176` reads a section the image names with the literal `"Items"`
  (`0x5c7c8c`) through `FUN_004e17d2`, builds each entry out of `Data.bin`
  (`0x609b18`, `FUN_004dcaf0`) and drops it as a sack at the coordinates the
  entry carries. What `FUN_004d00e9` is parsing at that point is
  `World\Mission\<n>.ini`, which `PARTY-SESSION-008` measured as not shipping
  in either root.
- `FUN_004f1fc6` scatters random treasure: `n` sacks (default
  `max(10, (W·H/400)/2)`) at cells `10 + rand()%(dim−20)`, each holding
  `rand()%3` items drawn by the shop's own generator at a 500 ceiling, and
  `100 + rand()%400` gold if it came out empty.
- The scatter is called a second time from `FUN_0050f511`, the sack manager's
  per-call pass, multiplayer only (`[0x005cd758]+0x0c != 0`), on a
  `rand()%20 < 2` roll, topping the map up to `W·H/400` live sacks.

**Confidence.** High for the enumeration figures and each owner's reading: the
enumeration is complete on the repaired table with `.rdata` slots included,
each owner was read, and the `"Items"` literal is the routine's only string
operand, from `StrDump fnstr:`. Medium that the `"Items"` entries are items
rather than something else the same section can hold: the per-entry parameter
slots were read as a shape, not matched against a shipped file, because none
ships.

**Amended.** EXP-0128 (`ITEM-SPAWN-026`, `ITEM-SPAWN-027`;
[`retracted.md`](retracted.md)) refuted the headline "Five things in the image
make a sack, and none of them is the `.alm`", the exclusivity clause, and the
conclusion "map-authored loot exists, in the mission `.ini` and not in any
`.alm` record set". The same `callto:50f5aa` enumeration, re-run, returns the
same 6 hits / 5 owners / 0 orphan with a different membership: `FUN_004e4f3e`,
reached only from the `.alm` loader `FUN_004e1924` and ungated there, calls the
sack maker at `004e59be` and is not in this claim's list; counting the two arms
of command `0x23` as two owners filled the fifth slot. Map-authored loot is in
the `.alm`, in the type-8 section (`ALM-SACK-065`). The reading of
`FUN_004f2176` and of the `.ini` stands; that it is the only authored source
does not. The map-load call site of `FUN_004f1fc6` is gated too, on
`server+0x120`, which makes the random scatter multiplayer-only at both sites
(`ITEM-SPAWN-027`).

### ITEM-SAVE-014

- `Item::Serialize` (`FUN_005115f4`) writes the shared 37-byte
  `Token::Serialize` head `FUN_00510e5c` (`SAV-TOKEN-034`), the effect list at
  `+0x20` (`FUN_00527b70`), then `u16 +0x40`, `u16 +0x42`, `u8 +0x44`,
  `u8 +0x45`, `u8 +0x46`, `u16 +0x48`, `u16 +0x4a`, `u8 +0x47`, in that order,
  storing and loading. The load arm then resolves the definition
  (`ITEM-DEF-002`).
- The effect list writes a u32 count followed by archive object references in
  list order. It loads by clearing, then appending the same count in archive
  order.
- One `Effect` body is the 37-byte Token head, then u8 kind, u8 mode, u32
  operand and u8 Token `+0x0c`: 44 bytes total. Transient Effect `+0x44` is not
  stored.
- `Armor`/`Shield`/`Weapon` each call Item first and add their own tail; Weapon
  also writes its owned Spell reference.
- The container's serializer is `FUN_00511a53`: the `CObList` elements
  (`FUN_005280b0`), then `+0x1c`, then `+0x20`. Insert index and load are
  stored, not recomputed.
- `Sack::Serialize` writes the head, `+0x3c`, then the container.
- `Unit::Serialize` (`FUN_00510518`) writes `actor+0x74` and `actor+0x78` as
  object references and, guarded on non-null, the container.
  `Humanoid::Serialize` (`FUN_00511760`) adds a raw 24-byte block at `+0x1cc`,
  the twelve slots `actor+0x198+4i` for `i = 1..12`, and one more reference at
  `+0x1e4`.
- `callto:511a53` = 3 hits / 2 owners / 0 orphan: the actor's two arms and the
  sack's. Exactly one container per actor and one per sack is serialized.

**Confidence.** High: the serializer bodies and field orders are complete, and
the Effect and list store/load arms agree.

**Amended.** EXP-0225 (`ITEM-EFFSAVE-077`) narrowed the shared Item-head
extent. This claim had retained `SAV-OBJ-014`'s retracted 16-byte head extent
after `SAV-TOKEN-034` had established 37 bytes ([`retracted.md`](retracted.md));
the first 16 bytes are a prefix, not the head's extent. Every Item tail field,
container member and enclosing actor/sack relation stands.

### ITEM-CARRY-015

- `PARTY-CARRY-005` established that what crosses is a single `WriteObject` on
  one actor plus sixteen bytes of `Player` scalars.
- `WriteObject` invokes that actor's own `Serialize`. For a hero that is
  `Human::Serialize → Humanoid::Serialize → Unit::Serialize` (`ITEM-SAVE-014`),
  which writes both hand slots, all twelve armour slots and the entire
  container as nested objects of the same graph. Equipment and backpack ride
  inside the carried blob.
- `PARTY-LOSS-006`'s strip list (three sub-objects freed, five fields zeroed,
  two pools refilled) touches none of them: `+0x74`, `+0x78`, `+0x7c` and
  `+0x198…+0x1e4` are not in it.
- A mercenary is in the same C++ family and its `Serialize` would write the
  same fields, but it is never a root and never reachable from one
  (`PARTY-MERC-007`: the graph has one root, and the `Player` is not in the
  blob). Nothing of a mercenary crosses except `MERC-DEATH-006`'s fifteen pool
  integers.
- A mercenary that dies first leaves nothing on the ground either, because its
  template name matches `ITEM-DEATH-012`'s `"NPC"` test.
- One party therefore has two persistences: the hero as an object graph with
  its gear inside it, the squads as counts. A consumer that treats them alike
  gets the mercenaries wrong.

**Confidence.** High: the chain is three serializers read whole plus one
published single-root claim, and the non-membership of the item fields in the
strip list is a comparison of two enumerated field lists.

### ITEM-PICK-016

- Order-space opcode `0x21`, arm `004d60ed`, looks a sack up at an arbitrary
  cell `(cmd+0x0a, cmd+0x0c)` in the world sack registry `[0x005f22c8]` via
  `FUN_00547be0`. It refuses with `"Sack not found at "` (`0x5c6368`) when there
  is none, writes the destination slot `cmd+0x0e` into `[actor+0x7c]+0x1c` as
  `ITEM-CMD-007`'s source-code-3 arm does, and calls
  `FUN_005308a0(actor, col, row)`.
- `FUN_005308a0`'s whole body is seven stores, among them `actor+0x54 = 0`
  (`0053092f`), `ord+0x50 = 0`, `actor+0x50 = 2` (`0053093b`),
  `ord+0x0a = (row<<8)|col` (`00530947`) and `ord+0x08 = 0` (`0053095e`).
- The relay then runs in four hops:
  1. `FUN_0052ce50` arm 2, the only consumer of `actor+0x50 == 2`, compares
     `ord+0x0a` against the actor's own cell (`FUN_00544a00(actor+0x10)`). Not
     there: `ord+0x08 = 1`, walk. There: `ord+0x08 = 7` (`0052d218`).
     `re:MOV byte ptr \[E.. \+ 0x8\],0x7$` is 1 hit / 1 owner / 0 orphan, so
     nothing else writes that value.
  2. The actor tick calls `FUN_005310e0` four instructions before it
     dispatches (`004f398c` vs `004f3994`).
  3. That routine's pending-order arm 7 (`005314df`) sets `actor+0x54 = 2`
     from `EBX`, live since `00531139`.
  4. The tick's own switch sends `actor+0x54 == 2` to `004f3e76`, the arm of
     `ITEM-PICK-009`: sack at the actor's cell, unregister, unlink,
     `FUN_004f4e3e`.
- On the next tick arm 7 finds `+0x54` already 2 and completes:
  `actor+0x50 = 0xc`, `ord+0x08 = 0`, `ord+0x50 = 1`.
- Session-space source 3 and order-space `0x21` are the same mechanism with and
  without a walk. The "precondition" `grpAI+0x20 == 0` is what the dispatcher
  establishes for every commanded set (`AI-CMD-033`), not a state the player
  must contrive.
- Contradiction C-6 dissolves: `ITEM-CMD-007` and `AI-STATE-011` read the same
  path from opposite ends.

**Confidence.** High. Each of the seven stores, each of the four hops and both
range tests is a named instruction. The alternative is enumerated away: a
second path would need a second `actor+0x50 == 2` consumer, and
`callto:52ce50` is 2 / 2 / 0 with one caller dead (`AI-DEAD-035`); or a second
`actor+0x54 = 2` writer, and the whole-image register-store sweep is 120 hits /
67 owners / 9 in orphan code, of which exactly one can carry a 2 because `EBX`
is 2 only between `00531139` and `00531502`.

## Weapon numbers and damage columns

| ID | Claim | Confidence | Status | Evidence |
|---|---|---|---|---|
| ITEM-SCALE-017 | A weapon's numbers are its `Data.bin` row's columns times one factor per attribute, each factor the product of a shape record and a material record. | High | ● active (partially retracted) | [EXP-0105](../experiments/EXP-0105-damage-fold/), [EXP-0106](../experiments/EXP-0106-weapon-columns/) |
| ITEM-DMGCOL-018 | An item's displayed damage and the pair it puts into the actor's combat block differ: the item holds `(base, spread)`, not the row's `(min, max)`. | High | ● active (partially retracted) | [EXP-0105](../experiments/EXP-0105-damage-fold/), [EXP-0106](../experiments/EXP-0106-weapon-columns/) |
| ITEM-LADDER-019 | The nine doubles of a shape or material record sit at `record + 0x20 + 8j`, fixed by the record stride `0x68`. | High | ● active | [EXP-0106](../experiments/EXP-0106-weapon-columns/) |
| ITEM-DMGFACT-020 | The weapon damage factor is the column the game names `@.damage`, at `record+0x40` = `f64[4]`, and a shape factor is never 1. | High | ● active | [EXP-0106](../experiments/EXP-0106-weapon-columns/) |
| ITEM-WEAPCOL-021 | The `Weapons` row's slot-to-title map is `slot i = title i+1`, anchored by four instruction-level uses; `Weapon::Equip` has three arms on slot 5. | High / Unknown | ● active | [EXP-0106](../experiments/EXP-0106-weapon-columns/) |
| ITEM-PANEL-022 | The item panel draws a weapon's damage line from the weapon's own two bytes as `[base, base+spread]`, composed the way the character sheet does. | High / Unknown | ● active | [EXP-0106](../experiments/EXP-0106-weapon-columns/) |

### ITEM-SCALE-017

- `FUN_0050dadc(this = weapon, arg = def row)` looks up the shape record
  `0x609b2c[item+0x45]` and the material record `0x609b18[item+0x46]` and
  writes six values:
  - `w+0x60` from runtime column 6 and `w+0x61` from column 7, both through the
    `+0x40` double;
  - `w+0x52` from column 8 through `+0x48`;
  - `w+0x6a` from column 9 through `+0x50`;
  - `w+0x4a` from column 3 through `+0x38`;
  - `w+0x48` from `+0x60 x +0x60` with no column and no rounding addend;
  - `w+0x50` from column 0xb verbatim, or 1 when that column is -1, through no
    factor at all.
- Names are parsed by `FUN_0050d670`: `FUN_004db6b5` writes `item+0x45` from
  the 5-row table, `FUN_004db801` writes `item+0x46` from the 16-row one, and
  `FUN_0051c240(0x609b68, remainder)` writes the row index `item+0x0c`. The
  grammar is a leading shape word, then a leading material word, then the rest
  as the Weapons row.
- The ten chargen literals resolve 10/10 on both roots.
- The runtime ladder is `+0x20 + 8j` with `f64[j] = title[j+1]`
  (`ITEM-LADDER-019`). `Common`'s `@.damage` is 0.2000, not 1.0, so a name
  carrying no leading shape word yields one fifth of the row's columns rather
  than the row's own numbers (`ITEM-DMGFACT-020`).

**Confidence.** High for the fill, the name parse and the 10/10 resolution: the
routine is read whole, and each factor is a named `FMUL` at a fixed
displacement. The ladder and every factor value derived from it are refuted.

**Amended.** EXP-0106 refuted two clauses, corrected in place
([`retracted.md`](retracted.md)). The runtime ladder was published as
`+0x28 + 8 x (title - 2)`; the record is `0x68` bytes
(`0051b97a IMUL EAX,EAX,0x68`) and nine doubles at `+0x28` would end past it
(`ITEM-LADDER-019`). The two collection labels were swapped: `0x609b18` is
Materials and `0x609b2c` the shape/tier table, which `DAT-OBJ-002` and
`SHOP-PRICE-011` had already published. The old rationale, pinned by five uses
landing on titles 4/5/6/7/9 in order, was the defect: ordering is invariant
under a uniform shift, so it fixes the ladder's spacing and never its origin.

### ITEM-DMGCOL-018

- The row ships `@.physicalMin` and `@.physicalMax` (runtime columns 6 and 7).
  After `FUN_0050dadc` the item carries `w+0x60 = round(min x f)` and
  `w+0x61 = round(max x f) - w+0x60`.
- The second byte subtracts the first: `0050db5d MOV DL,byte ptr [ECX + 0x60]`
  then `0050db66 FSUBP`, with `w+0x60` re-read as an already-rounded unsigned
  byte. The item holds `(base, spread)`, not `(min, max)`.
- For a single equipped weapon the two encodings carry the same information and
  a consumer cannot tell them apart. They diverge once a second contributor
  exists, because `FUN_004fa818` sums bases into `+0xb4` and spreads into
  `+0xb5` separately, so two sources' maxima never add (`HERO-FOLD-035`).
- Weapons row 0 is `BareHands` (`1..2`, attackType 3), and no instruction in the
  mapped image names it (`HERO-BARE-037`). It is a monster's `EquipItem` value
  and never a hero's default.
- Worked values: the factor is 0.2, so `Iron Short Sword` carries (5, 3), its
  `+0x52` is 5, and `Uncommon Steel Two Handed Sword` carries (13, 13)
  (`ITEM-DMGFACT-020`).

**Confidence.** High for the encoding: the subtraction is two named
instructions, the re-read's width is on the instruction itself, and
`ITEM-PANEL-022` adds a fourth independent kill of `(min, max)`. The worked
values first published are refuted. They were re-executed from both roots' own
`Data.bin` by a committed probe; the probe was faithful, the ladder it
re-executed was wrong, and reproducing an error on two roots looks like
corroboration.

**Amended.** EXP-0106 and owner measurement on the shipped game refuted the
worked values ([`retracted.md`](retracted.md), `ITEM-DMGFACT-020`): factor 1.0,
`Iron Short Sword` (23, 17), `+0x52` 1, and `Uncommon Steel Two Handed Sword`
(59, 58). The last bullet gives the corrected values. The encoding half stands.

### ITEM-LADDER-019

- `FUN_0051b810` → `FUN_0051b970` is the `GetAt` both scale collections use. It
  strides the record by `0x68` (`0051b97a IMUL EAX,EAX,0x68`, the same size
  `DAT-OBJ-002` publishes).
- The serialized record is a name plus `0x48` bytes = nine f64, so the head is
  `0x68 - 0x48 = 0x20` and `f64[j]` sits at `record + 0x20 + 8j`.
- The image reads seven of them, and the population is complete:
  - `+0x30` by the three price routines
    `FUN_0050c740`/`FUN_0050d0ce`/`FUN_0050dc4e` and by the shop window
    `FUN_005095ee`;
  - `+0x38`, `+0x50` and `+0x60` by all three item fills;
  - `+0x40` and `+0x48` by the weapon fill alone;
  - `+0x58` by the armor and shield fills alone.
- Group A ships 11 titles for one name and nine doubles, so one title has no
  slot. The six semantic uses fix the correspondence as `f64[j] = title[j+1]`.
  The two identically-zero doubles are then `Abbreviation` and `Materials`, two
  non-numeric CSV columns, and `Level` is left unserialized.
- Field-to-table pairing is a bound check, not a guess: `0050d923 CMP EAX,0xf`
  is applied to `item+0x46`, and `0x609b18` is the 16-row table, so `+0x46` is
  the material. `+0x45`, written by `FUN_004db6b5`, which searches
  `this + 0x14` = `0x609b2c`, is the shape.

**Confidence.** High. The origin is fixed by an immediate rather than by a fit,
and the rival is impossible rather than merely unsupported: `+0x28 + 8j` puts
`f64[8]` at `+0x68`, past the end of a `0x68`-byte record, while `+0x60` is
read by three routines. Two further discriminators agree. `+0x58` is read by
the two armour fills and never by the weapon, which only a correct ladder
explains. Under the rival the `Price` factor is `0.0000` on all 21 shipped
rows, so every item in the game would cost 0 against `SHOP-MISSION-020`'s
measured 4-of-367. The read population comes from `EnumRefs re:` over
`FMUL`/`FLD` with a record-shaped displacement on the repaired table: 101 hits,
0 in orphan or undisassembled code. It is blind only to a factor loaded through
a register-held address, which none of these is.

### ITEM-DMGFACT-020

- `Common` carries 0.2000 and `Iron` 1.0000. The shipped `Iron Short Sword`
  (shape `Common`, material `Iron`, row `Short Sword`, columns 6,7 = `23, 40`)
  resolves to factor 0.2 and carries:
  - `+0x60 = ftol(23 x 0.2 + 0.5)` = 5;
  - `+0x61 = ftol(40 x 0.2 - 5 + 0.5)` = 3;
  - `+0x52` = 5, `+0x6a` = 0, `+0x50` = 1, price 250.
- The values are identical on both roots, whose `Data.bin` files differ byte
  for byte.
- The other nine chargen literals: `Uncommon Bronze Axe` (5,6),
  `Uncommon Bronze Mace` (4,5), `Bronze Pike` (4,4),
  `Uncommon Wood Short Bow` (3,2), `Uncommon Steel Two Handed Sword` (13,13),
  `Uncommon Steel Axe` (10,11), `Uncommon Steel Mace` (8,9),
  `Uncommon Steel Pike` (9,9), `Uncommon Magic Wood Short Bow` (7,7).
- The `Weapons` columns are on a different scale from every quantity they
  meet: a row ships 23..40, the item carries 5..8, and a hero's own bare pair
  is 0..3 over the whole legal Body range. Only the scaled numbers are
  commensurable with a stat.
- `w+0x48 = ftol(MagCap_shape x MagCap_material)` is 0 for anything iron, which
  makes `SHOP-MAGIC-007`'s reject-and-redraw reachable. Under the superseded
  ladder every material read at least 1.0 and nothing was ever unenchantable.

**Confidence.** High. Each figure is a re-execution of `FUN_0050dadc` on
`ITEM-LADDER-019`'s ladder over both shipped roots, and the discriminating
evidence is `ITEM-LADDER-019`'s, not a fit. The clause this replaces was also
High and also re-executed by a committed probe; the difference is that the
ladder underneath this one is fixed by a record size rather than by an
ordering.

### ITEM-WEAPCOL-021

- Group C ships 18 titles for a name and 17 slots, with no spare, unlike group
  A, so `slot i = title i+1` exactly. Four independent uses land on it:
  - slot 5 `@.attackType`: `FUN_0050def2` reads it at `0050df76` and branches
    `CMP 0xa` / `CMP 0xb` / `CMP 0xc`; the shipped value set is
    `{1,2,3,4,5,11,-1}`;
  - slot 0xb `@.range`: the fill's `w+0x50`; `Weapon::Equip` does
    `actor+0x12c += w+0x50 - 1` at `0050e15c`, which is `UNIT-COMBAT-006`'s
    reach term;
  - slots 0xc/0xd `@.charge`/`@.relax`, assigned to `actor+0x134`/`+0x135` at
    `0050e126`/`0050e156`;
  - slot 0xe `2 handed`: `0050df3b CMP dword ptr [EAX],0x2`, the
    shield-removal gate.
- So slots 6 and 7 are `@.physicalMin`/`@.physicalMax`, and slots 2/3/8/9 are
  `Price`/`weight`/`@.toHit`/`#.deIrnce`.
- `Weapon::Equip` has three arms on slot 5, not one:
  - `< 0xa` is melee and adds `w+0x60`/`w+0x61` into `actor+0xf4`/`+0xf5`;
  - `== 0xb` and `== 0xc` add the same two bytes into `actor+0xf9`/`+0xfa` and
    assign the kind byte `+0xfb` to 1 or 2.
- The shipped `Flame Thrower` (`@.attackType` 11) therefore feeds the second
  damage component, not the first, and its type is never written to
  `actor+0xb6`.
- The melee arm also assigns `actor+0xf9/+0xfa/+0xfb` from
  `w+0x65/+0x66/+0x67`.

**Confidence.** High for the slot map: four uses, each a different instruction
on a different field, and the title list has no spare entry to absorb a shift,
unlike group A's ladder. High for the three arms, read from one listing.

**Unknown.** What writes `w+0x65/+0x66/+0x67`. `FUN_0050dadc` writes none of
them, so a third writer exists and was not enumerated.

### ITEM-PANEL-022

- `FUN_0047bf40` reads the equipped weapon at `actor+0x74` and appends two
  attributes through `FUN_00484af0`, a ten-instruction appender that stores a
  `(tag, value)` byte pair at `this+0x0c + 2 x this+0x09`:
  - `0047c3c0 MOV AL,byte ptr [EAX + 0x60]` as tag 0x0d;
  - `0047c3d0 MOV DL,byte ptr [ECX + 0x61]` as tag 0x0e.
- The whole block is gated on `0047c3b9`/`0047c3be`: a weapon with zero spread
  draws no damage line at all.
- `FUN_00484160` formats the list through a byte index table at `0x484954`
  into the jump table at `0x484930` (`ECX = tag - 1`, bounded `CMP ECX,0x32`).
- The `0x0d` arm at `0048422d` consumes both pairs. It prints the base with the
  format at `0x005bef64`, steps over the next tag byte, does
  `00484298 ADD EBX,EAX` and prints the one at `0x005bef60`. The item shows
  `[base, base+spread]`, and tag `0x0e`'s own table entry is the default arm,
  unreachable.
- The shipped `Iron Short Sword` therefore draws `5-8` while its `Data.bin` row
  ships `23..40`.
- This is a fourth independent kill of the `(min, max)` reading, alongside the
  resolver's roll, the fill's `FSUBP` and the hero sheet's `ADD`
  (`HERO-SHEET-038`), and the first that involves no inference from the hero
  path.

**Confidence.** High. Every hop is read at instruction level, and the two
format strings are dumped from `.rdata`. The tag tables are read out of the PE
by `tools/weaponcols -mode tags` rather than from a disassembler, because the
bytes have no containing function. The `0x0e`-is-unreachable clause is the one
inference: it follows from the `0x0d` arm consuming two pairs and from
`FUN_0047bf40` being the only site that emits `0x0d`, on an `EnumRefs disp:61`
of 25 hits over 15 owners with 0 in orphan code.

**Unknown.** The label text, which comes from `[0x005eb3d4]+0x1d8`. It is the
same string the character sheet's damage line uses, so the two surfaces are
labelled identically, but the string itself is a resource and was not
resolved.

## Appearance word

| ID | Claim | Confidence | Status | Evidence |
|---|---|---|---|---|
| ITEM-APPEAR-023 | `item+0x40` is the appearance word, assembled from three of the item's own bytes plus a kind rather than stored. | High / Unknown | ● active | [EXP-0115](../experiments/EXP-0115-held-item/) |
| ITEM-APPEAR-024 | A worn slot whose `graphics\inventory\` sheet is missing makes the engine show `"Invalid item weared "` on screen for five seconds. | High | ● active | [EXP-0115](../experiments/EXP-0115-held-item/) |
| ITEM-APPEAR-025 | The appearance `u16` has a second reader, which consumes field D alone, so field C reaches the drawn figure only through the sheet name. | High / Medium | ● active | [EXP-0119](../experiments/EXP-0119-figure-compositor/) |

### ITEM-APPEAR-023

- `ITEM-CLASS-001` records `Item::Item` zeroing `+0x40`. Five routines then
  write it, all with the result of the same builder, `FUN_00525d60`
  (`EnumRefs callto:525d60` = 5 hits / 5 owners / 0 orphan): `FUN_005081f1`
  (`005082ce`), `FUN_0050c53a` (`0050c72d`), `FUN_0050cfa2` (`0050d0c6`),
  `FUN_0050d5d4` (`0050d654`) and `FUN_0050d8e8` (`0050dac9`).
- Every site pushes the same four things in the same order: `item+0x0c`,
  `item+0x46`, `item+0x45`, then a kind.
- The builder, twenty-one instructions, returns
  `((item+0x46) << 12) | (kind << 8) | ((item+0x45) << 5) | (item+0x0c)`, each
  input masked to a byte first.
- `item+0x0c`, which `ITEM-CLASS-001` names as the definition row index, is the
  low five bits that `HERO-APPEAR-052`'s list is indexed by. `item+0x46` is
  `HERO-APPEAR-043`'s `material.reg` index. The kind is the equipment slot
  (`HERO-APPEAR-050`). `item+0x45` is three bits nothing here names.
- Because the builder masks to a byte and `OR`s rather than shifting into
  disjoint fields, a `+0x0c` above 31 would corrupt the `+0x45` field.

**Confidence.** High: all five call sites and the builder are read at
instruction level, and the argument slots are fixed by five sites pushing
identically rather than by one.

**Unknown.** What `item+0x45` means, and whether any shipped definition row
index exceeds 31.

### ITEM-APPEAR-024

- `FUN_00483d70` calls `FUN_00483c80` for the slot's seven-digit name
  (`00483da6`) and composes `graphics\inventory\` (`0x5beeb4`) + name + `.16a`
  (`0x5bc90c`). It tries to open it (`FUN_004c9f10` at `00483dfc`) and returns
  1 when it cannot (`00483e05`).
- The `0x76` message arm calls it per slot (`00413c0f`). On a true result it
  formats `0x5b83e4`, `"Invalid item weared "`, into the on-screen message
  routine `FUN_00401e70` with a 5000 ms lifetime (`00413c18`).
- Weared is the engine's own word for the twelve slots.
- `graphics.res` ships 416 `inventory/` nodes, all `.16a`, identical in count on
  both roots. That is a different tree and a different container from the 928
  `.256` figure sheets of `SPR256-EQUIP-042`, addressed by the same name.

**Confidence.** High: the composition, the open, the polarity of the return and
the message call are named instructions, and the node counts are exact archive
measurements on both roots.

### ITEM-APPEAR-025

- `ITEM-APPEAR-023` and `HERO-APPEAR-049` established the word's four fields
  and that `FUN_00483c80` formats all sixteen bits into a sheet name.
- The second reader is `FUN_00483e70` (`HERO-FIGURE-060`). Its first two
  instructions are `00483e70 MOV AL,byte ptr [ECX + 0x6]` and
  `00483e76 AND EAX,0x1f`: field D and nothing else. There is no shift, no
  second mask and no other read of `[ECX + 0x6]` in the routine, which is read
  end to end.
- It uses D to name a body and decides from that name which of the two held
  layers is painted last.
- On the whole figure path, field C (`item+0x45`) is consumed only by being
  spelled into the `%02d%02d%1d%02d` / `%02d%02d%03d` file name. A consumer that
  gets C wrong loses a sheet, not a behaviour.

**Confidence.** High for the second reader's field use: one routine read end to
end, both instructions verified against raw bytes on both roots. Medium, and
scoped, for "field C affects nothing else": the statement is about the figure
path (the compositor, the frame selector and this predicate), not the whole
image. `disp:45` is too common a displacement to carry an image-wide absence;
`docs/INSTRUMENT.md` rule 4 requires it be reduced by the field's width first,
which was not done.

## Authored sacks and item codes

| ID | Claim | Confidence | Status | Evidence |
|---|---|---|---|---|
| ITEM-SPAWN-026 | A sack has a sixth origin, the `.alm`, and on a single-player campaign map it is the only one that runs at load. | High | ● active | [EXP-0128](../experiments/EXP-0128-alm-items-section/) |
| ITEM-SPAWN-027 | The random treasure scatter runs on multiplayer maps only, at both of its call sites, so a campaign map has authored sacks and no random ones. | High | ● active | [EXP-0128](../experiments/EXP-0128-alm-items-section/) |
| ITEM-OWNED-028 | A `.alm` type-8 record whose `+0x04` is non-zero is not a sack: it stocks an actor that already exists. | High / Medium | ● active | [EXP-0128](../experiments/EXP-0128-alm-items-section/) |
| ITEM-CODE-029 | An authored item is a packed `u16`, and its bits 8..11, the nibble below the top one, pick the C++ class. | High / Unknown | ● active | [EXP-0128](../experiments/EXP-0128-alm-items-section/) |

### ITEM-SPAWN-026

- `EnumRefs callto:50f5aa` on the repaired table with `.rdata` slots included
  returns 6 hits, 5 distinct owners, 0 in orphan or undisassembled code. The
  owner set is `FUN_004f1fc6`, `FUN_004e4f3e`, `FUN_004f2176`, `FUN_0050f6d2`
  and `FUN_004d5dd8` (twice, the two arms of command `0x23`).
- `ITEM-SPAWN-013` published the same 6/5/0 figures with a different
  membership: it named the two arms of `0x23` as two owners and did not name
  `FUN_004e4f3e` at all.
- `FUN_004e4f3e` has one caller, `004e1e5f` inside the `.alm` loader
  `FUN_004e1924`, and there is no gate on it. It walks the type-8 record list
  at `map+0x2dc` and calls the sack maker at `004e59be` (`ALM-SACK-065` has the
  grammar).
- Map-authored loot is in the `.alm`, in the section `ALM-TRIG-050` measured
  the shape of and EXP-0080 read as spellbook payload.
- The `.ini` path `ITEM-SPAWN-013` describes is real and correctly read.
  `FUN_004e17d2` is a text parser (`004e1870 PUSH 0x3d` = `=`,
  `004e187d PUSH 0x7b` = `{`) filled by `FUN_004e0b0c`, which is
  `CStdioFile::Open(path, 0x4020)` splitting on `[` and `;`. Its file does not
  ship (`PARTY-SESSION-008`), so it contributes nothing to a shipped campaign.
- The other three origins are unchanged: the drop/death wrapper, the two arms
  of command `0x23`, and the random scatter, which `ITEM-SPAWN-027` shows never
  runs on a campaign map.

**Confidence.** High: one call enumeration, complete on the repaired table with
0 orphan hits, and every owner read. The missing owner is a `CALL` immediate in
a routine whose own single call site is a `CALL` immediate in the map loader,
so no enumeration instrument is needed to see either.

### ITEM-SPAWN-027

- `ITEM-SPAWN-013` reads the map-load call site as ungated and only the second,
  in `FUN_0050f511`, as multiplayer-gated. The map-load site is gated too:
  `004d04dc CMP dword ptr [EAX + 0x120],0x0` and `004d04e3 JNZ 0x004d04f2` skip
  `FUN_004f1fc6` whenever `server+0x120` is non-zero.
- `server+0x120` has two writers image-wide (`disp:120`, 34 hits / 25 owners /
  0 orphan, of which these two touch this object):
  - `004ced82 MOV dword ptr [EDX + 0x120],0x0` in the world reset
    `FUN_004cec1d`;
  - `004e1be8` in the map loader `FUN_004e1924`, which computes it from the
    map.
- The loader's computation: `004e1bc1 CMP dword ptr [EAX + 0xd4],0x1`,
  `004e1bc8 SETG CL` and `004e1bd1 MOV dword ptr [EDX + 0xc],ECX` set the
  server's multiplayer flag; then `004e1bdb CMP dword ptr [EAX + 0xc],0x0`,
  `004e1bdf SETZ CL` and `004e1be8` set `server+0x120` to its negation.
- `server+0x120 != 0` is exactly "this map is single-player", and the scatter
  is skipped exactly then.
- `ALM-REQ-056` records `map+0xd4` defaulting to 1 when the type-0 record is
  absent. `ALM-META-058` measures the related `+0x70` at 1 on 56/56 campaign
  maps against 4/8/12/16 on 15 of 16 loose ones: the same single-versus-multi
  split from the file side.
- A consumer that believed the old reading would give every campaign map
  `max(10, (W*H/400)/2)` sacks of random shop stock that the original never
  places: 10 on an 80x80 map, 25 on a 144x144 one, against the 4 or 5 the
  `.alm` authors.

**Confidence.** High: the gate is one `CMP`/`JNZ` pair at the call site, its
operand has a two-instruction derivation in the map loader, and the writer set
is an image-wide displacement enumeration with 0 orphan hits.

### ITEM-OWNED-028

- `FUN_004e4f3e` opens by walking the actor registry `[0x00609558]` and filling
  a map keyed by each actor's own `+0x08` (`004e4fe7 ADD EDX,0x8`,
  `004e4fee CALL 0x0051c8f0`, `004e4ff9 MOV dword ptr [EAX],ECX`).
- For each type-8 record it looks `+0x04` up in that map (`004e54ce`,
  `004e54ef`). The sack arm is guarded by
  `004e5921 CMP dword ptr [EAX + 0x3c],0x0` and `JNZ`, so a non-zero `+0x04`
  never reaches `FUN_0050f5aa`, and the record's coordinates and gold are never
  read.
- Each element goes to the resolved actor instead. When the element's `+0x04`
  is zero, `004e5900 MOV ECX,dword ptr [ECX+0x7c]` and
  `004e5903 CALL 0x0050e8e2` append it to the carrier `ITEM-CONT-004` places at
  `actor+0x7c`. When it is non-zero, `004e58f4 CALL dword ptr [EDX + 0x3c]`
  routes it through the actor's own vtable.
- Corpus, both roots: 43 such records on each, carrying 25 elements, of which 5
  have a non-zero element `+0x04`; the largest is 20 records on one map.
- The `.alm` authors starting inventories as well as ground loot, through one
  record type. A consumer that reads every type-8 record as a sack would place
  43 sacks that the original does not.

**Confidence.** High for the arm and the two destinations: each is a named
instruction on a path read end to end, and the `actor+0x7c` destination is the
field `ITEM-CONT-004` independently identified as the carrier. Medium for what
the ids identify: the key is `actor+0x8` and the map is built from the
registry, but nothing here reads what writes `actor+0x8` or what the shipped
values 1..N correspond to in the `.alm`'s own record sets.

### ITEM-CODE-029

- `FUN_004dcf92` takes the code and cuts four fields out of it:
  - `004dcfa3 SAR EAX,0x8` then `004dcfa6 AND EAX,0xf` (bits 8..11);
  - `004dcfb5 SAR ECX,0xc` then `AND ECX,0xf` (bits 12..15);
  - `004dcfc7 SAR EDX,0x5` then `AND EDX,0x7` (bits 5..7);
  - `004dcfd8 AND EAX,0x1f` (bits 0..4).
- It hands them to `FUN_004dd02a`, which allocates by the first:
  - 1 → `004dd05c PUSH 0x84`, ctor `FUN_0050d861`;
  - 2 → `004dd0be PUSH 0x68`, ctor `FUN_0050cf26`;
  - 3..13 → `004dd12b PUSH 0x68`, ctor `FUN_0050c4be`;
  - 14 → `004dd188 PUSH 0x50`, ctor `FUN_00507fb3`;
  - anything else → `004dd1d1 XOR EAX,EAX`, a null the map-load caller skips
    at `004e55b4`.
- Those four sizes are `ITEM-CLASS-001`'s `Weapon 0x84`, the two `0x68`
  siblings `Armor` and `Shield`, and `Item 0x50`.
- Class 14 differs twice: `004dcfea` routes it down a branch that passes 0, 0
  for the middle two fields, and `004dd00e AND EAX,0xff` takes the whole low
  byte, not five bits, as the index.
- Corpus, both roots, over the 181 / 177 authored elements: the classes used
  are {1, 2, 4, 5, 6, 7, 8, 9, 10, 12, 14}, 0 codes resolve to null, and 0
  elements have any bit set above the low `u16`.

**Confidence.** High: the four masks and the four allocation sizes are named
immediates on one path read end to end, and the sizes agree independently with
a class table read from `CRuntimeClass` records in another experiment.

**Unknown.** The meaning of the three ctor arguments: they are read here and
not followed into `Data.bin`.

## Humans spawn equipment and armour slots

| ID | Claim | Confidence | Status | Evidence |
|---|---|---|---|---|
| ITEM-HUMEQ-030 | The ten item names on a `Humans` row are equipment cells whose position picks the C++ class; the container receives only what equipping refuses or displaces. | High / Medium | ● active | [EXP-0131](../experiments/EXP-0131-human-spawn-equipment/) |
| ITEM-ARMSLOT-031 | The armour field a piece takes is the `Slot` column of its own `Armors` row, not the actor or the cell; it reaches `actor+0x19c`/`+0x1a0`, which take-off cannot. | High / Medium | ● active | [EXP-0131](../experiments/EXP-0131-human-spawn-equipment/) |

### ITEM-HUMEQ-030

- `FUN_004f9065`'s loop runs `i = 0..9` over `entry+0x1c`, skips an empty
  string (`004f9376 CALL 0x004c7c20`, `004f937d JLE`) and allocates by `i`
  alone:
  - `004f93a7 CMP dword ptr [EBP + -0x24],0x0` → `PUSH 0x84` / `FUN_0050d670`,
    a Weapon;
  - `004f9403 CMP … ,0x1` → `PUSH 0x68` / `FUN_0050ccd4`, a Shield;
  - otherwise `PUSH 0x68` / `FUN_0050c2ea`, an Armor.
- No arm reads the string before it allocates. The `Units` arm is the
  opposite: `UNIT-EQUIP-005` found the class chosen by searching for the
  literal `"Shield"`.
- All three end `CALL dword ptr [E?? + 0x3c]` on the actor (`004f93fb`,
  `004f9454`, `004f94a4`). `FUN_004f8e78` installs the `Human` vptr `0x59c4d0`
  at `004f8ea8` before calling the streamer at `004f8ecb`, so that slot is
  `FUN_004f7099` and not the base `FUN_004f4dc7`.
- `FUN_004f7099` calls `actor->vt+0x38(item)` (`004f70a2` pushes the item,
  `004f70a6` makes the actor the `this`) → `FUN_004f705b`. Its runtime-class
  test against `Armor`'s `0x5c33a0` selects between two byte-for-byte identical
  calls to `item->vt+0x38(actor)`, so no per-class restriction is applied. It
  appends to `actor+0x7c` only the return value.
- `Weapon::Equip` → `actor+0x74` (`0050df73`) and has no refusal arm at all.
- `Shield::Equip` → `actor+0x78` (`0050d38c`), refusing with `"Invalid shield"`
  when `shield+0x0c == 0` (`0050d2ec`).
- `Armor::Equip` refuses with `"Illegal armor"` when `armor+0x50 == 0`
  (`0050c8de`) or `actor->vt+0x30()` is 0 (`0050c915`). A refusal returns
  `this`, which is how an item reaches the backpack.
- One cross-cell arm exists. Cell 0 is equipped while `actor+0x78` is still
  empty, so `Weapon::Equip`'s two-hand arm cannot fire during a spawn.
  `Shield::Equip`'s mirror at `0050d350` reads param `0x0e` of the already-worn
  weapon's row and, on `2`, takes the weapon off and pushes it into
  `actor+0x7c` itself (`0050d381`).
- Corpus, both roots: 166 of 215 rows name an item; 909 cells → 908 worn, 1
  carried, 0 destination collisions. The one carried item is a shipped typo
  (`Gauntless`) that resolves to no row. 46 rows name both a weapon and a
  shield, and 0 of those weapons is two-handed, so the take-off arm never fires
  on shipped data.

**Confidence.** High for the dispatch and the three destinations. Every hop is
a named instruction in a routine read whole. The vtable slot is a raw pointer
read at a fixed address (`DecompAddrs vt:`, which consults no call graph and so
has none of `docs/INSTRUMENT.md` rule 7's blind spots), and the vptr install
precedes the call that uses it. Medium for the corpus figures, which are
agreement over two roots and rest on a name parse reproduced from four helpers
whose `CString` primitives are identified by use rather than by symbol.

### ITEM-ARMSLOT-031

- `FUN_0050c53a` zeroes `armor+0x50` (`0050c55c`) and resolves
  `armor+0x3c = ElementAt(armor+0x0c)` on `0x609b54`, the Armors collection of
  `ITEM-DEF-002` (`0050c56e`). Then `0050c579 PUSH 0x4` and
  `0050c58e MOV byte ptr [ECX + 0x50],DL` copy param 4 of that row into
  `armor+0x50`. The shipped column title for param 4 is `Slot`.
- `Armor::Equip` indexes with it directly:
  `0050c937 CMP dword ptr [ECX + EAX*0x4 + 0x198],0x0`; the occupant is
  unequipped through its own `vt+0x3c` (`0050c962`) and returned; then
  `0050c973 MOV dword ptr [EDX + ECX*0x4 + 0x198],EAX`.
- The array `ITEM-EQUIP-006` found at `actor+0x198 + 4i` is addressed by the
  same slot namespace as the take-off arm. Indices 1 and 2, `actor+0x19c` and
  `actor+0x1a0` (`ITEM-EQUIP-006`'s Unknown), are writable here and not
  addressable there, because take-off maps its slots 1 and 2 to `actor+0x74` /
  `actor+0x78` instead.
- `0050c599 CMP ECX,0xc` and `JLE` discard a `Slot` above 12 with
  `"Invalid armor part "`, without restoring `armor+0x50` to zero first, so such
  a piece would carry an index past the twelve.
- Corpus, both roots, all 30 shipped `Armors` rows: `Slot` ∈
  {4,5,6,7,8,9,10,12} with counts 1/1/8/5/4/3/4/4. 0 rows use 1, 2, 3, 11 or
  anything above 12, and 0 of the 909 shipped equipment cells produces two
  pieces in one field.

**Confidence.** High for the column, the index arithmetic and the
two-namespace mismatch: each is a named instruction, and `Armors` is fixed by
the collection base `ITEM-DEF-002` already carries. Medium for "parts 1, 2, 3
and 11 are never used", which is a census over both roots and nothing stronger.

## Armour and shield fill, worn strip and suitability

| ID | Claim | Confidence | Status | Evidence |
|---|---|---|---|---|
| ITEM-ARMFILL-032 | An `Armors` or `Shields` row fills four scaled `u16` words: absorption truncates, defence and weight round, and the uncolumned `+0x48` word is magic capacity. | High / Medium | ● active | [EXP-0132](../experiments/EXP-0132-armour-fold/) |
| ITEM-ARMFOLD-033 | Equipping an armour or a shield adds two of fourteen block fields, so a breastplate moves defence and absorption and no protection or damage-kind byte. | High / Unknown | ● active (amended, partially retracted) | [EXP-0132](../experiments/EXP-0132-armour-fold/) |
| ITEM-CORPSE-034 | At death every worn piece reaches the sack through a virtual strip call, and the worn array is a 1-based thirteen-dword array whose element 0 nothing can fill. | High / Medium | ● active | [EXP-0134](../experiments/EXP-0134-corpse-worn-armour/) |
| ITEM-SUIT-035 | Parameter 15 of a group-C row is titled `sutableFor` and is a two-bit mask over consumer classes; a weapon suitable for nobody is not dropped at death. | High / Medium | ● active (partially retracted, superseded) | [EXP-0134](../experiments/EXP-0134-corpse-worn-armour/) |

### ITEM-ARMFILL-032

- `FUN_0050c53a` (Armor) and `FUN_0050cfa2` (Shield) are the same routine
  twice, offset by two bytes because the armour spends `+0x50` on its `Slot`.
- Both resolve `mat = Materials[item+0x46]` (`0x609b18`) and
  `shp = Shapes[item+0x45]` (`0x609b2c`), with `f64[j]` at `record+0x20 + 8j`
  (`ITEM-LADDER-019`). They then write:
  - runtime column 10 as `ftol(col × mat.f64[7] × shp.f64[7])` into
    `armor+0x54` / `shield+0x52`, u16 (`0050c67f`, `0050d01d`);
  - runtime column 9 as `ftol(col × mat.f64[6] × shp.f64[6] + 0.5)` into
    `armor+0x52` / `shield+0x50`, u16 (`0050c6af`, `0050d04d`);
  - an uncolumned `ftol(mat.f64[8] × shp.f64[8])` into `+0x48`, u16
    (`0050c6c7`, `0050d065`);
  - runtime column 3 as `ftol(col × mat.f64[3] × shp.f64[3] + 0.5)` into
    `+0x4a`, u16, the weight (`0050c6f7`, `0050d095`).
- The armour alone also copies column 4 verbatim into `armor+0x50`
  (`ITEM-ARMSLOT-031`). The shield fill has no `PUSH 4` and reads no `Slot`.
- The column-10 term is the only scaled term in any of the three item fills
  with no `FADD [0x0059bc48]`. That constant is `0.5`, read from the image, so
  defence and weight round to nearest while absorption truncates. Over the
  shipped tables 188 of 3 120 (row × shape × material) combinations differ
  under the two rules.
- The identification is not fitted. The fill joins two title lists in two
  `Data.bin` groups by numeric index alone, and under the published ladders all
  four pairings hold at once, in all three item classes: col 3 `weight` ×
  `weight`, col 9 `#.deIrnce` × `#.defence`, col 10 `#.absorbtion` ×
  `#.absorption`, uncolumned × `MagCap`. A one-slot error in either ladder
  breaks four pairings rather than one. So `ITEM-SCALE-017`'s uncolumned
  `+0x48` word is magic capacity.

**Confidence.** High for the arithmetic and the ladder pairing. Both routines
were read end to end, and all 69 load-bearing instruction sites were re-checked
literally against both roots' `rom.exe` bytes through the PE section table with
no disassembler in the path, 69/69; `rom.exe` is byte-identical across roots.
Medium for the name `MagCap`, which rests on the factor table's own title and
on no reader of `+0x48`.

### ITEM-ARMFOLD-033

- Both classes embed one 0x16-byte block, built by `FUN_004fa49b` →
  `FUN_004fa4b1` = `memset(this, 0, 0x16)`, at `armor+0x52` (`0050c31b`) and
  `shield+0x50` (`0050cd05`). The two words `ITEM-ARMFILL-032` writes are its
  members `+0x00` and `+0x02`; nothing writes the other twelve.
- `Armor::Equip` stores the item at `actor[0x198 + 4·(armor+0x50)]`
  (`0050c973`), evicting any occupant through its `vt+0x3c` at `0050c962`.
  `Shield::Equip` stores it at `actor+0x78` (`0050d38c`).
- Both then pass the address of the whole block to `FUN_004fa4eb` twice,
  `this = actor+0xfe` (`0050c98a`, `0050d39f`) and `this = actor+0xbe`
  (`0050c99f`, `0050d3b4`), then pass `(s16)item+0x4a` to `FUN_004f36f7`.
- Armor sets `actor+0x150` bits `0x0c00` before the Effect walk; Shield sets
  them after it. Removal is class-specific: both refresh negative weight before
  subtracting the blocks, but Armor clears its slot and sets flags before
  removing Effects, while Shield removes Effects before flags and slot clear
  (`SAV-EQUIPORDER-552`).
- `FUN_004fa4eb` folds fourteen members (`HERO-FOLD-033`), so the six
  protections `actor+0xc2…+0xcc` and the six damage-kind bytes
  `actor+0xce…+0xd3` receive zero from any armour or shield.
- There is no per-slot difference. The slot byte indexes the array and is read
  nowhere else, and neither block add sits under a branch on it. The shipped
  variation is data: `Slot ∈ {4,5,6,7,8,9,10,12}` on `Armors` and `2` on all
  nine `Shields` rows, a column the shield path never reads.
- `vt+0x54` is not the recompute it looks like. The `vt+0x54`
  `ITEM-EQUIP-006` records is `FUN_004fc22a` on both humanoid vtables
  (`[0x59c448 + 0x54]` and `[0x59c4d0 + 0x54]`, raw dwords). It sums
  `actor+0x1cc + 4i` for `i = 1..5` (`004fc23a`, `004fc24c`, `004fc25b`), the
  skill-experience caches of `HERO-XP-010`, five dwords past the twelve slots,
  and stores `ftol(sum × 0.01)` into `actor+0x1c` (`004fc27e`).
- It reads neither the slot array nor either block. The base class's slot is a
  bare `return actor+0x1c` (`FUN_00523270`), and `FUN_004f4d98` discards the
  result, so the call is a side effect and equipment is stored, never
  recomputed.

**Confidence.** High. Every routine was read end to end and every site was
re-checked byte for byte on both roots, 69/69, with the vtable slots read as
raw `.rdata` dwords rather than through a call graph.

**Unknown.** Whether a later path writes a worn item's other twelve block
members. The effect and magic-item routines were not read, and this is two
routines read whole rather than a sweep.

**Amended.** The event-order clause is refuted ([`retracted.md`](retracted.md),
EXP-0288, `SAV-EQUIPORDER-552`). It stated a common flag/Effect order for both
classes and a removal sequence at "the same sites"; the per-class attach and
removal order above replaces it. The block layout, the fill population and the
`+54` Token-value meaning stand.

### ITEM-CORPSE-034

- `ITEM-DEATH-012`'s published sequence steps over a virtual call between the
  weapon arm at `004f4fda` and the `"NPC"` find at `004f5029`:
  `004f5010 MOV ECX,dword ptr [EBP + -0x3c]`, `004f5013 MOV EDX,dword ptr [ECX]`,
  `004f5015 MOV ECX,dword ptr [EBP + -0x3c]`,
  `004f5018 CALL dword ptr [EDX + 0x44]`. Two `rel8` in that listing pin
  `004f5010` (`004f4fd8 74 36`, `004f4ff0 74 1e`).
- Slot `+0x44` is `FUN_00523250` on the base actor `0x59c3c0`:
  `PUSH EBP`/`MOV EBP,ESP`/`PUSH ECX`/store/leave/`RET`, no body. On both
  `0x59c448` and `0x59c4d0` it is `FUN_004f70f8`.
- `FUN_004f70f8` loops `004f7101 MOV dword ptr [EBP + -0x4],0x1` …
  `004f7113 CMP dword ptr [EBP + -0x4],0xd` | `JGE`, i = 1..12 inclusive. Per
  slot it does `004f711f MOV EAX,dword ptr [EDX + ECX*0x4 + 0x198]`,
  `004f712f CALL dword ptr [EDX + 0x40]` and `004f7139 CALL 0x0050e8e2` with
  `ECX = actor+0x7c`.
- Empty slots are free: humanoid `vt+0x40` `FUN_004f70cf` returns 0 for a null
  argument (`004f70da`), and the append `FUN_0050e902` returns before storing
  when its item is 0 (`0050e90c`/`0050e912`). `Unit` is `0x198` bytes, so the
  empty base override keeps the loop off four bytes past the end of a monster.
- The strip runs before the `"NPC"` delete (`004f505e`), before the emptiness
  test (`004f513c`) and before the sack (`004f5159`). A body that wore anything
  leaves a sack even if it carried nothing, and a suppressed template's armour
  is destroyed with the container rather than left on it.
- The array is thirteen dwords. The ctor `FUN_004f6ded` clears `i = 0..12`,
  while the strip, the destructor `FUN_004f6f09` and both arms of the serializer
  `FUN_00511760` run 1..12, and `Armor::Equip` refuses `Slot == 0` outright
  (`0050c8de` | `0050c8e3`, `"Illegal armor"`). So `actor+0x198` is written by
  nothing and read by nothing. This closes `ITEM-EQUIP-006`'s
  `+0x19c`/`+0x1a0` Unknown: they are ordinary slots, stripped, serialized and
  destroyed like the other ten.
- Disposing is not dropping: `FUN_004f6f09` calls `worn[i]->vt+0x04(1)`, the
  deleting destructor, and names `actor+0x7c` nowhere.
- `FUN_004f7144` is the same loop plus `+0x74`/`+0x78` into a container passed
  in, and it is unreachable: `callto:` 0 hits. A raw dword scan of the whole
  image finds its address 0 times on both roots, against 2 for `FUN_004f70f8`
  (`0059c48c`, `0059c514`, each opening a 33-pointer `.rdata` run) and 1 for
  `FUN_00523250` (`0059c404` = `0x59c3c0 + 0x44`).
- Corpus, both roots: all 30 shipped `Armors` rows carry `Slot` in 1..12, so
  every shipped armour piece is stripped. Boots are `Slot` 12, the last index
  the loop reaches.

**Confidence.** High for the path, the range and the class split. The dispatch
is a named instruction; its two possible targets are fixed by a raw dword scan
over the whole image rather than by a call graph, and both were read whole; the
loop bound is an immediate; and two `rel8` in its own listing pin the branch
target the sequence turns on. Medium for "the strip is the only route by which
a worn piece reaches a container", which rests on a `disp:198` sweep whose
stated blind spot is a structure-wide `REP MOVSD` carrying no displacement.

### ITEM-SUIT-035

- `DAT-SCHEMA-007`(a)'s shared title array is verified by identity: the walk
  returns the same 18 strings for `Armors`, `Shields` and `Weapons`, byte for
  byte, on both roots. Slot 15 is `sutableFor`. Slot 4 is `Slot`, agreeing with
  `ITEM-ARMSLOT-031`'s independent reading of `0050c579 PUSH 0x4`.
- `FUN_0050cb03`, the `Armor` display arm, fetches the same column
  (`0050cb28 PUSH 0xf` → `FUN_0051ab50`, `0050cb38 MOV EDX,dword ptr [EAX]`)
  and tests it twice, independently: `0050cb40 AND EAX,0x1` →
  `0050cb4d OR DL,0x2` into the record's `+0x4`, and `0050cb59 AND ECX,0x2` →
  `0050cb66 OR AL,0x4` into the same byte. Two bits, two flags, no ordering.
- The death gate reads the value whole, `004f4fed CMP dword ptr [EAX],0x0`, so
  its meaning is neither bit set. This closes `ITEM-DEATH-012`'s Unknown and
  makes its weapon gate readable.
- Corpus, both roots, identical:
  - Weapons, 27 rows: `0` on exactly three, `Boulder Thrower` /
    `Flame Thrower` / `Sonic Beam`; `1` on nineteen martial weapons; `2` on
    `Shaman Staff` and `Staff`; `3` on `BareHands` and `Plasma Sword`; one row
    (`rem`) carries no parameter array at all.
  - Shields, 9 rows, all `1`.
  - Armors, 30 rows: `1` on nineteen metal pieces, `2` on the nine cloth pieces
    (`Cap`, `Cape`, `Cloak`, `Dress`, `Gloves`, `Hat`, `Low Hat`, `Robe`,
    `Shoes`), `3` on `Amulet` and `Ring`.
- The three zero rows are the innate attacks of large monsters, which is what
  the gate is for.

**Confidence.** High for the title and for the mask. The title array's identity
across the three collections is a measurement made on both roots, and the two
bits are tested separately in one routine read whole. Medium for reading bit 0
and bit 1 as the two player classes: the partition is compelling and is corpus
agreement, and an ordinal reading is excluded only because nothing orders
`BareHands`, `Amulet` and `Plasma Sword`.

**Amended.** The consumer clause of the Medium grade is superseded
([`retracted.md`](retracted.md), EXP-0168, `ITEM-WEAR-055`, `ITEM-WEAR-057`).
It read: no routine was found that refuses an item on account of this column;
`ITEM-HUMEQ-030` establishes the equip path applies no such restriction, so two
display flags are the only consumers read. The two bits are the inputs to a
refusal enforced in the client (`ITEM-WEAR-057`); the equip path is still
restriction-free, so `ITEM-HUMEQ-030` is untouched. `retracted.md` records the
fighter/mage reading as confirmed, and the title, the two-bit mask, the two
independent tests and the corpus census as standing.

## Item pictures

| ID | Claim | Confidence | Status | Evidence |
|---|---|---|---|---|
| ITEM-PICT-046 | The picture-name formatter `FUN_00483c80` is not an `Item` method: it reads a display record whose `+0x06` word is a copy of `item+0x40`. | High / Unknown | ● active | [EXP-0144](../experiments/EXP-0144-item-pictures/) |
| ITEM-PICT-047 | The picture map is partial over the 5 329 constructible items and total over what the data asks for; its inverse over the 416 `inventory/*.16a` nodes is exact. | High | ● active | [EXP-0144](../experiments/EXP-0144-item-pictures/) |
| ITEM-PICT-048 | The two name-word parsers return different not-found defaults, shape 0 `Common` and material 15 `None`, and that asymmetry makes every shipped item drawable. | High | ● active | [EXP-0144](../experiments/EXP-0144-item-pictures/) |
| ITEM-PICT-049 | The engine has two self-check messages about missing item art, and the second, `Invalid item in inventory`, gates on a keyed lookup and opens no file. | High / Unknown | ● active | [EXP-0144](../experiments/EXP-0144-item-pictures/) |
| ITEM-PICT-050 | A five-word material mask on each `Data.bin` row decides which (material, shape) pairs an item may take, and it predicts the picture set exactly: 416 legal items. | High | ● active | [EXP-0144](../experiments/EXP-0144-item-pictures/) |
| ITEM-PICT-051 | The authoring code and the picture word share a layout but are produced by different code; an armour draws by its row's `Slot`, not by the code's kind nibble. | High | ● active | [EXP-0144](../experiments/EXP-0144-item-pictures/) |
| ITEM-PICT-052 | The name table, picture tree and material masks, three files decoded three ways, each give exactly 416 items: 0 named-not-drawn, 0 drawn-not-named. | High | ● active | [EXP-0144](../experiments/EXP-0144-item-pictures/) |

### ITEM-PICT-046

- `ITEM-APPEAR-024`, `ITEM-APPEAR-025` and `HERO-APPEAR-049` all describe
  `FUN_00483c80` reading `word[slot + 0x6]` without naming the object. It is a
  display record, and `item+0x40` reaches it by an ordinary copy.
- `0047c368 MOV DX,word ptr [ECX + 0x40]` and
  `0047c458 MOV AX,word ptr [EDX + 0x40]` are each followed at `0047c36f` /
  `0047c465` by `CALL FUN_00484ac0`. Those are that record's only two
  constructor call sites.
- `FUN_00484ac0`'s body after the prologue is `00484ac6 MOV ESI,ECX`,
  `00484ac8 MOV word ptr [ESI + 6],AX`, `00484ad0 MOV EAX,dword ptr [ESP + 0xc]`,
  `00484ad3 MOV byte ptr [ESI + 0xa],AL`.
- Its copy constructor `FUN_00484990` copies `+0x04` u16, `+0x06` u16, `+0x08`
  u8, `+0x09` u8, `+0x10` dword, `+0x14` dword.
- The nine call sites of `FUN_00483c80` fetch their `this` from a `CObArray`
  (`00481f12 MOV ECX,dword ptr [EDX + EAX*0x4]`) or from a pointer array at
  `+0x15c + 4i` (`00413bb5`, `0045ef69`), never from an item minus `0x3a`.
- `item+0x40` and `record+0x06` are two fields, both correctly reported. A
  consumer that puts the appearance word on the item object alone has nothing
  to feed the panel with.

**Confidence.** High. Both copies are named instructions immediately before the
record's only two constructor calls, and the constructor and the copy
constructor were each read whole.

**Unknown.** Which class the record belongs to, and what its other five fields
are.

### ITEM-PICT-047

- `FUN_00525d60` and `FUN_00483c80` were re-executed over the whole
  constructible population: the 27 `Weapons`, 9 `Shields` and 30 `Armors` rows
  across 16 `Materials` and 5 `Shapes`, plus the 49 `MagicItems` rows, which
  take neither because `004dd014` / `004dd016` push 0, 0.
- That gives 5 329 addressable items. 416 have an `inventory/*.16a` node and
  4 913 do not: 7.8 %.
- The inverse is exact: all 416 nodes are claimed, 0 are claimed by no item, 0
  are claimed by more than one, and 0 carry a name no `u16` produces. Identical
  on both roots.
- The 4 913 are not a gap. `ITEM-PICT-050` finds the row's own five-word
  material mask, which admits exactly the 367 shape/material combinations that
  have art, so this figure measures the addressing space, not a hole in the
  game.
- The two `sprintf` arms cannot collide. Arm A runs only when the kind nibble is
  exactly 14 and prints it in fixed positions 2..3, while arm B's widths
  (2,2,1,2) consume all sixteen bits. The name space is a bijection onto 65 536
  `u16` values, and 416 of them are used.

**Confidence.** High. An exhaustive re-execution of two routines read at
instruction level over a closed population, inverted over the whole shipped
node set, reproduced on both roots.

### ITEM-PICT-048

- `FUN_004db6b5` ends `004db7ef XOR AL,AL`: no `Shapes` word means index 0,
  `Common`. `FUN_004db801` ends `004db931 MOV AL,0x0f`: no `Materials` word
  means index 15, the row literally named `None`.
- Both scan their collection downwards from the last row (`004db821` fetches the
  count, then `SUB EAX,1`, and the loop decrements). That is why `Hard Leather`
  is matched before `Leather` and `Magic Wood` before `Wood`.
- With those two returns the census closes. Over the shipped corpus, 181 of 181
  `.alm` type-8 authored item codes on the root with ten loose maps and 177 of
  177 on the other resolve to a node that exists, and 908 of 908 resolvable
  `Humans` equipment name cells do, on both roots: 1 089 / 1 085 requests,
  0 misses.
- The one cell that fails, `F_KnightLeader1` row 143 cell 5,
  `Leather Gauntless`, fails one layer earlier at the collection lookup and is
  `ITEM-HUMEQ-030`'s known typo.
- With a material default of 0 instead, exactly 106 cells miss, all of them the
  nine cloth rows whose only art is material 15. The two defaults are not
  interchangeable, and the corpus discriminates between them.

**Confidence.** High. The two returns are named immediates at the tail of two
routines, and the census is an exhaustive re-derivation over both roots that a
wrong default visibly breaks.

### ITEM-PICT-049

- `ITEM-APPEAR-024` publishes `"Invalid item weared "`. Beside it, `0x5b83ac`
  is `"Invalid item in inventory "`, composed at `004134b7` with the same
  seven-digit name `FUN_00483c80` produced at `00413493`.
- It is shown by the same routine `FUN_00401e70` with a 10 000 ms lifetime
  (`004134e7 PUSH 0x2710`), against the other's 5 000.
- Its gate is `004134de CALL FUN_00483e50` | `004134e3 TEST EAX,EAX` | `JNZ`,
  the opposite polarity to `FUN_00483d70`'s.
- `FUN_00483e50` is seven instructions that do no file I/O:
  `00483e51 MOV CX,word ptr [ECX + 0x6]` | `00483e5b MOV ECX,0x5eb410` |
  `00483e60 CALL 00571716`, a keyed lookup in one global. It has three call
  sites (`004134de`, `00413585`, `00413f0b`); `FUN_00483d70` has exactly one.
- The literal `graphics\inventory\` at `0x5beeb4` is composed at five sites:
  `00481f27`, `004838c0`, `00483db1`, `00491815`, `004a3571`. `00483db1` is the
  open-and-close check; the other four are loads on separate surfaces.
- The archive's own node paths begin at `inventory/`, not
  `graphics/inventory/`: the composed path's first segment names the container.

**Confidence.** High for the message, the lifetime, the polarity and the site
counts. Each is a named instruction, or an exhaustive scan of the image for the
string's address and for `rel32` targets.

**Unknown.** Which panel each of the four load sites draws, and what the object
at `0x5eb410` is keyed on beyond the word.

### ITEM-PICT-050

- The shop's candidate loop tests the mask in six instructions:
  `0050966d MOV DX,word ptr [EAX + ECX*0x2 + 0x1c]` (the row, subscripted by
  the tier), `00509672 MOV EAX,0x1` | `0050967a SHL EAX,CL` (`1 << material`),
  `0050967c AND EDX,EAX` | `0050967e TEST EDX,EDX` | `00509680 JZ`. A material
  whose bit is clear is skipped before the price window is computed.
- `row+0x1c` is the ten raw bytes the group-C `Data.bin` grammar has carried
  uninterpreted since `DAT-OBJ-002` (`DAT-MATMASK-020`): five `u16`, one per
  `Shapes` row.
- Against the art, on both roots, 5 280 of 5 280 (row, tier, material) cells
  agree: the mask bit is set exactly where a picture exists, 0 disagreements.
  367 bits are set and there are 367 pictures; the 49 magic items outside the
  scheme make 416.
- Per material the bits run Iron 24, Bronze 48, Steel 32, Silver 5, Gold 4,
  Mithrill 40, Adamantium 60, Meteoric 29, Wood 16, Magic Wood 19, Leather 14,
  Hard Leather 15, Dragon Leather 7, Bone 0, Crystal 13, None 41.
  `Bone` is legal on no row at any tier, which is why it has no art and why no
  generator can offer it.
- `ITEM-PICT-047`'s 4 913 pictureless items are therefore not items the data
  permits; they are the arithmetic product of two tables the mask never joins.
- G2: the mask is the customisation seam. Widening what an item may be made of
  is ten bytes on one `Data.bin` row plus one new node in `graphics.res`, and
  changes no other shipped file. The mask is a `u16` and `FUN_00525d60` `OR`s
  unmasked operands, so a customisation cannot exceed 16 materials, 8 shapes
  and 31 rows per class without changing the name scheme itself.

**Confidence.** High. The six instructions are one path read in sequence, and
the agreement is exhaustive over the whole cross product on both roots; a mask
off by one row, one tier or one bit would break it in dozens of cells at once.

### ITEM-PICT-051

- `ITEM-CODE-029`'s packed `u16` and `ITEM-APPEAR-023`'s `item+0x40` have the
  same four fields in the same bits, and on shipped data they are equal, but
  the picture word is rebuilt, not carried.
- `FUN_004dcf92` splits the code and `FUN_004dd02a` passes only three values to
  each constructor: `0050c4fc MOV DL,byte ptr [EBP + 0x8]` → `+0x45`,
  `0050c505 MOV CL,byte ptr [EBP + 0xc]` → `+0x46`,
  `0050c50e MOV AL,byte ptr [EBP + 0x10]` → `+0x0c`. The kind nibble is not
  passed.
- `FUN_0050d5d4` and `FUN_0050cfa2` then supply the literals 1 and 2. The
  armour fill `FUN_0050c53a` supplies `byte[item+0x50]`, which
  `ITEM-ARMSLOT-031` shows it has just written from the row's own `Slot` column.
- So an authored code whose kind nibble is 7 while its `Armors` row carries
  `Slot` 6 allocates by 7 and draws as 6.
- Corpus, both roots: on all 181 / 177 authored codes the nibble already equals
  the row's `Slot`, so the two encodings have never been observed to diverge. A
  consumer that stores the authored code as the appearance passes every shipped
  test and diverges on the first authored mismatch.

**Confidence.** High. The constructor's three stores and the three kind sources
are named instructions, and the corpus agreement is exhaustive on both roots.

### ITEM-PICT-052

- `ITEM-DISPNAME-036` reads `main.res:text/itemname.bin` as `filesize/2` `u16`
  keys into the name map and counts 416. `graphics.res:inventory/*.16a` read as
  leaf names counts 416. `ITEM-PICT-050`'s masks admit 367 shape/material
  forms, which with the 49 `MagicItems` rows is 416.
- Inverted through `FUN_00483c80` and compared as sets of `u16`: 416 in both,
  0 named but not drawn, 0 drawn but not named, on both roots.
- An item's name and its picture are addressed by the same value and are
  provisioned for exactly the same population. `ITEM-NAMEMISS-039`'s admission
  list (a name miss drops the element from the inventory rather than showing a
  fallback) never fires on a legal item, because no legal item lacks a name.
- This is the discriminating check for the whole addressing model: the two file
  sets are produced by different tools from different data, and a single wrong
  bit anywhere in the four fields would show as a non-empty symmetric
  difference.
- Consequence for `ITEM-NAMEPOP-038`: its "416 of 6 064 expressible items have
  a line; the other 5 648 cannot be shown at all" counts a cross product the
  data never joins. Under the row's own mask those 5 648 are not expressible,
  and nothing that can exist lacks either a name or a picture.

**Confidence.** High. An exact set identity between two independently decoded
shipped files, reproduced on both roots, corroborated a third time by a mask
read out of a third file.

## Display names and the packed item key

| ID | Claim | Confidence | Status | Evidence |
|---|---|---|---|---|
| ITEM-DISPNAME-036 | The displayed item name is a stored `itemname.txt` line selected through a `CMap<u16,const char*>` at `0x5eb410`, not composed and not the `Data.bin` row name. | High / Medium | ● active | [EXP-0142](../experiments/EXP-0142-item-names/) |
| ITEM-NAMEKEY-037 | The name key is the packed item code: material `k>>12`, class `(k>>8)&0xf`, shape `(k>>5)&7`, row `k&0x1f`; class 14 spends the whole low byte on the row. | High | ● active | [EXP-0142](../experiments/EXP-0142-item-names/) |
| ITEM-NAMEPOP-038 | 416 items have a name; 5 648 further keys the same encoding can address do not, and the shipped name set is not a product of its fields. | High / Medium | ● active (amended) | [EXP-0142](../experiments/EXP-0142-item-names/) |
| ITEM-NAMEMISS-039 | An inventory element whose code is not in the name map is dropped, not shown under a fallback name, so the name table is the inventory's admission list. | Medium | ● active | [EXP-0142](../experiments/EXP-0142-item-names/) |
| ITEM-NAMEPARSE-040 | `FUN_0050d670` is `Weapon::Weapon(const char*)`, and the tier/material/shape name grammar is an authoring parser that no display path reaches. | High / Medium | ● active | [EXP-0142](../experiments/EXP-0142-item-names/) |
| ITEM-BRACE-041 | No item-collection row name carries a brace: the `{...}` suffix belongs to the equipment cell, and none of it reaches the display key. | High / Medium | ● active (amended) | [EXP-0142](../experiments/EXP-0142-item-names/), [EXP-0225](../experiments/EXP-0225-item-effect-grammar/) |
| ITEM-NAMELIMIT-042 | The name decode implies customisation limits: a display name is a whole text line, the `u16` key's field widths are hard, and a `Data.bin` row name is a lookup key. | Medium | ● active | [EXP-0142](../experiments/EXP-0142-item-names/) |

### ITEM-DISPNAME-036

- `FUN_00468380` is the whole of the build, 206 bytes. It opens
  `main\text\itemname.bin` (`004683a8 PUSH 0x5bd04c`), takes the length
  (`004683c2 CALL 0x004ca130`), halves it (`004683c9 SHR EBP,0x1`) for the
  entry count, `malloc`s `2 x count` and reads the node whole.
- It then loops `i = 0..count-1` doing two things.
  `004683fa MOV ECX,0x5eb3b0` / `CALL 0x004687f0` fetches line `i` of the
  string array loaded from `main\text\itemname.txt`.
  `00468406 MOV AX,word ptr [EDI]` / `0046840a MOV ECX,0x5eb410` /
  `CALL 0x00571738` / `0041841a MOV dword ptr [EAX],EBX` stores that pointer in
  the map under the `u16` at `bin[2i]`.
- `FUN_00468190` is the map's `InitHashTable(10)`.
- The reader is `FUN_00483e50`, four instructions of body:
  `00483e51 MOV CX,word ptr [ECX + 0x6]` / `00483e5b MOV ECX,0x5eb410` /
  `00483e60 CALL 0x00571716`.
- The key blob is freed at `00468424` and the strings are not, so the map holds
  pointers into the string array. Nothing bounds `i` against that array's
  length, so a `.bin` longer than the `.txt` reads past it.
- The path strings are reached at their first byte, `0x005bd04c`, ten bytes
  before the `itemname` a string search lands on. An `EnumRefs` aimed at the
  search hit returns 0 hits for both nodes, which is what an absence would look
  like.

**Confidence.** High for the build and the read: two short routines read whole
with every immediate quoted, and the model re-executed over the whole shipped
corpus with 0 unresolved keys and 0 duplicates. Medium for the onward path: the
returned pointer is handed to interface code and is a `char*`, but no listing
was followed from there to `claims/text.md`'s converter.

### ITEM-NAMEKEY-037

- `FUN_00483c80` reads the same `+0x6` word and formats it with a `sprintf`.
  `if ((code & 0xf00) == 0xe00)` takes `%02d%02d%03d` of `code>>12`,
  `(code>>8)&0xf`, `code&0xff`; otherwise `%02d%02d%1d%02d` of `code>>12`,
  `(code>>8)&0xf`, `(code>>5)&7`, `code&0x1f`.
- Both literals are read from the raw section bytes (`0x005beefc`,
  `0x005beeec`) rather than from a decompiler symbol name. They carry no
  separator, so the diagnostic form is a bare seven-digit run.
- These are the same four cuts `FUN_004dcf92` makes to construct an item
  (`ITEM-CODE-029`). The encoding an authored `.alm` element carries and the
  encoding the interface looks a name up by are one encoding.
- The class field selects the collection by `FUN_004dd02a`'s allocation switch
  (1 `Weapons`, 2 `Shields`, 3..13 `Armors`, 14 `MagicItems`). For an armour it
  is also the equipment slot: over the 193 shipped armour keys the class field
  equals the row's own `Slot` column (`Armors` param 4, `ITEM-ARMSLOT-031`) 193
  times with 0 exceptions. Those two numbers live in two different files,
  `main.res:text/itemname.bin` and `world.res:data/data.bin`.

**Confidence.** High. The split is read off literals in one routine and is then
discriminated rather than merely agreed with: a wrong field boundary would have
to reproduce a 193/193 cross-file agreement by accident. All 416 keys resolve
to a real row with no duplicate.

### ITEM-NAMEPOP-038

- Both roots ship a 832-byte `text/itemname.bin`, byte-identical, sha256
  `adb09ccc7a210ead12812901f847e299747138030cdd085addb1a1895a7c1e24`, and a
  416-line `text/itemname.txt`.
- The 416 keys split `Weapons` 138, `Shields` 36, `MagicItems` 49 and `Armors`
  193. The armour keys use eight of the eleven armour slot values
  (4/5/6/7/8/9/10/12 to 11/14/52/25/26/20/22/23). The shape field takes all five
  shipped values (128/97/80/73/38 for
  `Common`/`Uncommon`/`Rare`/`Very Rare`/`Elven`).
- Re-executing the collections against the encoding gives 6 064 reachable
  (class, shape, material, row) keys, counting only armour rows whose own
  `Slot` matches the class, so 5 648 have no line.
- The set is not a product, and that refutes composition positively. Key
  `0x6582` is shape 4 `Elven`, material 6 `Adamantium`, `Armors` row 2
  `Amulet`: the EN line is `Adamantium Amulet` with the tier word absent, while
  the RU line for the same key carries a tier word. Key `0x1502` is Bronze /
  `Amulet` / `Common` and reads `Beard`, and `0x0502` (Iron) reads
  `Magic Beard`.

**Confidence.** High for the counts and the encoding: a complete re-execution
over both roots, 0 unresolved, 0 duplicate. Medium for `5 648` as a count of
reachable keys. It excludes armour rows whose `Slot` disagrees with the class
value, which is what the shipped data does, but nothing in the image forbids an
authored code from pairing any class value with any row.

**Amended.** EXP-0144: the `5 648` figure counts a cross product the data never
joins. Each `Data.bin` row carries a five-word material mask (`ITEM-PICT-050`)
that admits 367 shape/material forms, which with the 49 `MagicItems` rows is
the same 416 keys (`ITEM-PICT-052`). The Medium clause stands as written:
`5 648` counts the addressing space, not items the data permits.

### ITEM-NAMEMISS-039

- `FUN_004104e8` calls `FUN_00483c80` first (`00413493`), prefixes `0x5b83ac`,
  the literal `Invalid item in inventory ` (`004134b8`, concatenated at
  `004134c4`), and only then performs the lookup
  (`004134de CALL 0x00483e50`).
- `004134e3 TEST EAX,EAX` / `004134e5 JNZ 0x0041350b` takes a hit to the
  `CObArray` grow-and-store at `+0xc8`. The fall-through pushes `0x2710` and
  `0x5e9720`, calls `0x00401e70` through `[panel+0xa10]` and jumps past the
  append.
- A second site at `00413585` skips the append on a miss silently. A third at
  `0041362c` calls `Lookup` inline.
- `0041376c MOV word ptr [ECX + 0x6],0xffff` writes the empty-slot sentinel
  that `FUN_00483970` tests with `if (*(short *)(elem + 6) == -1)` before it
  asks for a name at all.

**Confidence.** Medium. The branch, the two pushes and the join are quoted from
one listing read at three sites, but `FUN_004104e8` is 7 327 instructions and
only these surroundings were read. That `0x00401e70` is a diagnostic sink
rather than a display call is inferred from its two literal arguments and was
not followed.

### ITEM-NAMEPARSE-040

- It calls `Item::Item` (`0050d692 CALL 0x00507f42`), installs the `Weapon`
  vtable (`0050d6b7 MOV dword ptr [EAX],0x59ca00`, `ITEM-CLASS-001`'s third
  sibling) and returns `RET 0x4`. The sibling constructors are `FUN_0050c2ea`
  (Armor), `FUN_0050ccd4` (Shield) and `FUN_00507fb3` (Item).
- Its body, in order:
  1. `FUN_004db944` cuts the `{...}` tail and returns it.
  2. `FUN_004db6b5` strips a `Shapes` word into `+0x45`.
  3. `FUN_004db801` strips a `Materials` word into `+0x46`.
  4. `FUN_0051c240` looks the residue up in `Weapons`
     (`0050d75f PUSH 0x609b68`) into `+0xc`; index 0 prints
     `Invalid weapon %s - no such ID` and abandons the fill.
  5. Only on success does `0050d811 CALL 0x0050d8e8` fill the item and
     `0050d829 CALL 0x005089f5` receive the tail.
- This confirms `MAGIC-SPELLHOP-023`'s reading independently, with the
  correction that the routine is one class's constructor and not a shared
  parse.
- The inverse also exists and is one routine: `FUN_004dcaf0(const char*)` runs
  the same strips and then tries `Armors`, `Weapons`, ` Shield` plus `Shields`,
  and a `MagicItems` walk before `Invalid item %s`. It is the only routine in
  the image that references all five item globals `0x609b18` / `0x609b40` /
  `0x609b54` / `0x609b68` / `0x609b7c`.

**Confidence.** High for the constructor: it is read whole with its vtable
store, its four helper calls and its two error arms, and the three sibling
constructors are identified by their own allocation sizes. Medium for the
clause that no display path reaches the parse. It rests on the display path
being the map lookup of `ITEM-DISPNAME-036` plus an owner sweep of the
collection accessor, whose 13 owners are all in the data, item and shop
modules. That accessor is one of several per-type `GetAt` thunks, so it bounds
one collection type rather than the question.

### ITEM-BRACE-041

- A complete walk of both databases returns 0 matching row names across
  `Weapons` (27), `Armors` (30), `Shields` (9) and `MagicItems` (49) per root.
- Equipment has 260 braced cells per root, `Humans` 258 and `Units` 2, hence
  520 combined. It has 261 opening-brace characters per root, 522 combined,
  because `Humans` runtime row 79 slot 4 contains a second opening brace before
  its first close.
- The first-open/first-close parser makes that cell one Gauntlets item with
  `defence=7` and rejects the swallowed `Rare Mithrill Chain Boots{defence=19`
  token as an unknown key. The separating space is optional.
- The display consequence: `ITEM-NAMEKEY-037` carries material, class, shape
  and row and no effect field, so an enchanted piece and a plain piece with the
  same four fields resolve to the same line.
- The `MagicItems` group-D extra string is distinct, unbraced input.

**Confidence.** High for framing and the display consequence, which follow from
complete parser control flow and the key's bit layout. Medium for every corpus
population, because the two roots' corresponding equipment text is identical.

**Amended.** EXP-0225: 520 is a combined count of matching cells, not brace
characters.

### ITEM-NAMELIMIT-042

- (a) A display name is a whole line of a text node. The shipped extremes are
  4..29 bytes on the EN root and 3..44 on the RU root; no length is compared
  anywhere on the path, and the only ceiling found is the drawing surface.
- (b) The key is a `u16` and its fields are hard. `Weapons` / `Armors` /
  `Shields` get 5 bits of row, so a 32nd row of any of those three cannot be
  named or authored at all (shipped 27/30/9). Material has 4 bits (shipped 16,
  full), shape 3 bits (shipped 5), and class 4 bits, of which 11 values are
  armour slots. `MagicItems` alone gets 8 bits of index (shipped 49).
- (c) A display name is not a lookup key: the key is numeric, so renaming a line
  renames the item and breaks nothing.
- (d) A `Data.bin` row name is a lookup key, in two directions, since the four
  constructors and `FUN_004dcaf0` resolve authored text against `Shapes`,
  `Materials` and the row collections. Renaming a row silently breaks every
  `Units.EquipItem` cell, every `Humans` equipment cell and every mission
  `.ini` item line that named it; `ITEM-HUMEQ-030`'s one shipped `Gauntless` is
  that failure already present in the shipped data.
- (e) Adding a name changes `main.res` alone, two nodes whose lengths must stay
  paired; adding an item changes `world.res:data/data.bin` as well.
- (f) `REG-NAME-055`'s 15-character registry clamp does not carry across: no
  clamp was found on either a row name or a display line.

**Confidence.** Medium. (b) and (c) follow from bit widths and from the key
being numeric, and (d) and (f) from complete sweeps of the shipped data and of
the parse routines. (a) rests on the absence of a compare along a path of which
only the lookup sites were read, and no experiment has authored a file and
observed the result.

## The document access item

| ID | Claim | Confidence | Status | Evidence |
|---|---|---|---|---|
| ITEM-DOC-053 | One item raises the documents panel and has a name and a picture; no shipped map, script instant 12 or image immediate places it, within the families searched. | High / Medium | ● active | [EXP-0151](../experiments/EXP-0151-mission-documents/) |
| ITEM-DOC-054 | The measured start-path censuses do not decide whether mission 10 starts with the document access item. | High / Unknown | ✔ promoted (amended) | [EXP-0154](../experiments/EXP-0154-campaign-start/), [EXP-0221](../experiments/EXP-0221-valuable-documents/) |
| ITEM-DOC-069 | Player-hero construction calls a gated `Quest Documents` producer before assigning the hero's runtime id. | High / Unknown | ✔ promoted | [EXP-0221](../experiments/EXP-0221-valuable-documents/) |

### ITEM-DOC-053

- `FUN_00492410` is the only site in `rom.exe` that posts message `0x463` to a
  window. It does so under a three-part gate on the object at `app+0x3cc`:
  - `004924a7 TEST byte ptr [ECX+0x8],0x6` clear;
  - `004924b3 AND ECX,0xf00` / `004924b9 CMP ECX,0xe00` (class field 14);
  - `004924c1 AND AL,0x1f` / `004924c3 CMP AL,0x1c` (low five bits 0x1c).
- `+0x06` is the packed item code of `ITEM-NAMEKEY-037`. The routine's other
  work is to store the same object into `[actor+0x15c + slot*4]`.
- Code `0x0e1c` resolves in `main.res::text/itemname.bin` to row 211:
  material 0 `Iron`, class 14 `MagicItems`, row 28, `Data.bin` row name
  `Quest Documents`, English display name `Valuable Documents`, key `0014028`.
- Class 14 spends the whole low byte on the row while the gate masks five bits,
  so every code with class 14 and row ≡ 28 (mod 32) passes: rows 28, 60, 92,
  124, 156, 188, 220, 252 for any material nibble.
- Corpus: 0 type-8 elements carry code `0x0e1c`, over every shipped map of
  either root, against 181 (EN) and 177 (RU) authored elements.
  `TRIG-ADDITEM-027`'s exhaustive instant-12 census gives 7 nodes with codes
  0xe1a, 0xe1d, 0xe1d, 0xe1e, 0xe23, 0xe24, 0xe25, none of them `0x0e1c`.
  `EnumRefs imm:e1c` finds no such immediate anywhere in the image.

**Confidence.** High for the gate: five named instructions on one path, and the
`PostMessageA` call is the message's only sender. High for the item's identity:
the key formatter and the name table are `ITEM-NAMEKEY-037`'s, and 211 is one
row of 416 with no duplicate. Medium for "nothing places it". Three producer
families were searched exhaustively: the `.alm` type-8 section, script
instant 12, and image immediates. Three were not: shop stock generation,
character generation, and any `Data.bin`-driven starting equipment. The absence
is scoped to what was searched, and the discriminator is a running original in
which the panel opens at all.

### ITEM-DOC-054

- The access item is `MagicItems[28]` `Quest Documents`, packed code `0x0e1c`,
  localized as `Valuable Documents` in EN and `Официальные Документы` in RU.
- `FUN_00492410` accepts an eligible class-14 item whose low five row bits equal
  28 and posts message `0x463`. The campaign arm opens the panel, and the panel
  binds separately to the campaign record at `campaign+0x548`.
- The searches found no `0x0e1c` in all 909 `Humans` equipment cells, the four
  selected character rows, direct weapons, mission 10's seven type-8 elements,
  all seven instant-12 grants, image immediates, campaign/session
  initialization, or live pre-town shop generation.
- `ITEM-DOC-069` identifies a conditional string producer called during hero
  construction.
- Runtime prediction: if `server+0x0c == 0` and `[0x005eb5a4] != 2` at that
  call, the new hero holds the item before its runtime id is assigned.
  Observing both predicates pass without the item refutes the producer reading.

**Confidence.** High for identity, localization, UI route, collection
separation, and each scoped producer search.

**Unknown.** Actual fresh-campaign item presence, because `[0x005eb5a4]` was
not measured at the call.

**Amended.** EXP-0221 withdraws the complete-negative headline and runtime
prediction ([`retracted.md`](retracted.md), `ITEM-DOC-069`). The headline read
"Mission 10 starts with document content but no document access item", and the
prediction stated that no item is visible or held and the panel does not open
at mission 10 start. The scoped censuses remain valid, but they did not cover a
string-based conditional producer. Identity, localization, UI route,
collection separation and each named census stand.

### ITEM-DOC-069

- `FUN_004d3f6b` creates the item only when `server+0x0c == 0` and
  `[0x005eb5a4] != 2`. It allocates an Item, constructs it from the exact row
  name, and appends it through `FUN_0050e8e2` to the argument actor's container
  at `+0x7c`.
- The executable has exactly two instruction-validated direct `E8 rel32`
  callers and no absolute dword reference to the routine.
- The shipped-campaign caller `FUN_004d3755` runs after starting-skill and
  derived-stat application and before the same hero pointer receives its
  runtime id, owner and actor-list membership.
- The other caller is `FUN_004d403c`'s optional `Humans` builder, which shipped
  maps do not reach.
- The campaign construction path establishes the first gate, not the second.

**Confidence.** High for the conditional constructor, the destination, the two
encoded caller forms and the relative ordering.

**Unknown.** Actual fresh-campaign item presence; the semantic name and
fresh-campaign value of `[0x005eb5a4] == 2`; and computed indirect entries
carrying no encoded producer address.

## Wear rule

| ID | Claim | Confidence | Status | Evidence |
|---|---|---|---|---|
| ITEM-WEAR-055 | Whether a character may use an item is one 13-instruction predicate over one bit of the character and two bits of the item's shipped `sutableFor` column. | High | ● active (amended, partially retracted) | [EXP-0168](../experiments/EXP-0168-wear-rule/), [EXP-0192](../experiments/EXP-0192-mission20-party-boundary/) |
| ITEM-WEAR-056 | The shipped `sutableFor` distribution over `Armors`, `Shields`, `Weapons` and `MagicItems`, identical on both roots, and the descriptor byte each item therefore carries. | High / Medium | ● active (amended, partially retracted) | [EXP-0168](../experiments/EXP-0168-wear-rule/) |
| ITEM-WEAR-057 | The wear rule is enforced at exactly two sites, both in the client, a shop background and an equipment-doll drop refusal; the simulation applies none of it. | High / Medium | ● active | [EXP-0168](../experiments/EXP-0168-wear-rule/) |
| ITEM-WEAR-058 | A `MagicItems` item's class restriction comes from its row name, not a column: `Book` rows are mage-only, and the doll refuses a spellbook by a second gate. | High / Unknown | ● active | [EXP-0168](../experiments/EXP-0168-wear-rule/) |

### ITEM-WEAR-055

- `FUN_00460440(member, element)` accepts fighter-bit items when
  `member+0x18c & 2` is clear and mage-bit items when it is set.
- `FUN_0045f850` derives that member bit from `(typeID-0x21)&2` inside
  `0x20..0x3f`. Its separate low-type arm also sets it for `typeID` `0x17` and
  `0x18`.
- A Human constructor reaches the `gender+0x21/+0x23` typeID overwrite only in
  non-zero player-character constructor mode; map Humans may retain a table
  typeID, as mission 20 does.
- The polarity still follows for player characters: sex is the addend and class
  the base, so descriptor bit 1 is fighter and bit 2 mage.
- On the item side, the Armor, Shield and Weapon descriptor builders copy
  Data.bin parameter 15 bits independently to `element+0x08`, and their
  four-entry virtual producer set is complete. The embedded descriptor offset
  is fixed by the effects, consume, effect-kind and packed-code
  correspondences.
- Value 0 is a real fourth state meaning usable by neither.

**Confidence.** High for the predicate, the producers, the column and the
corrected player-character polarity.

**Amended.** The universal Human-constructor sentence is refuted
([`retracted.md`](retracted.md), EXP-0192, `PARTY-M20-031`). It read
"`Human`'s constructor builds `typeID = sex + 0x23` … or `sex + 0x21`"; the
overwrite occurs only in non-zero player-character constructor mode, and the
client has a separate low-type arm for `0x17/0x18`. Mission 20 falsifies only
that former universal constructor wording, not the wear rule: the predicate,
the item producers and the player-character bit polarity stand.

### ITEM-WEAR-056

- `tools/wearrule` walks `world.res::data/data.bin` on both preserved roots. The
  group-C title array is 18 strings with the same sha256 on both roots, and
  slot 15 is `sutableFor`.
- Armors, 30 rows: 19 fighter-only, 9 mage-only (`Hat`, `Cap`, `Low Hat`,
  `Cloak`, `Cape`, `Robe`, `Dress`, `Gloves`, `Shoes`), 2 both (`Ring`,
  `Amulet`).
- Shields, 9 rows: all 9 fighter-only.
- Weapons, 27 rows: 19 fighter-only, 2 mage-only (`Staff`, `Shaman Staff`),
  2 both (`BareHands`, `Plasma Sword`), 3 neither (`Sonic Beam`,
  `Flame Thrower`, `Boulder Thrower`). One row (`rem`) has a 17-cell parameter
  array whose every cell is `-1` on both roots, so slot 15 reads `-1` and both
  bits are set, which makes three rows with both bits set.
- MagicItems, 49 rows: 44 usable by both, 5 mage-only (`ITEM-WEAR-058`).
- The two censuses are identical across roots row for row once the root label
  is removed.
- An ordinal reading of the column is excluded twice over. The producers test
  bit 0 and bit 1 independently in one routine with no ordering between them.
  Value 3 occurs on `BareHands`, `Amulet`, `Ring` and `Plasma Sword`, which no
  ordering of fighter below mage places above value 2.

**Confidence.** High for the column identity and the per-row values: the title
array is measured on both roots, and the values come from a walk that consumes
`data.bin` to its last byte. Medium for reading the partition as fighter/mage
from the data alone, which is corpus agreement; the discriminating evidence for
that reading is `ITEM-WEAR-055`'s instruction side, not this census.

**Amended.** The `rem` clause is corrected by a re-measurement dated 2026-08-15
with this experiment's own `data.bin` walker over `world.res::data/data.bin` on
both preserved roots, with no new experiment ([`retracted.md`](retracted.md)).
The clause had read: one row (`rem`) carries no parameter array at all, so the
slot is absent and both bits stay clear. `Weapons[23]` has 17 parameter cells
and every cell is `-1`, and rows 22 and 24 parse with correct names and values
on either side of it, so the array is present and consumed exactly. The
published `group-c-census.csv` value `-1` for the row is correct; its
`no parameter array` verdict is not, because `param()` returns `-1` both for an
absent slot and for a stored `-1`. The two descriptor arms `AND` the slot value
with 1 (`0050cb40`) and with 2 (`0050cb59`) without a sign test, so a stored
`-1` sets both bits. No item is built from the row: its five-word material mask
is `0000` for all five shapes, so no (material, shape) pair is legal
(`ITEM-PICT-050`).

### ITEM-WEAR-057

- `EnumRefs callto:460440` returns 2 hits, 2 distinct owners, 0 in orphan or
  undisassembled code.
- Site 1, presentation only: the shop cell painter `FUN_004a2fe0` calls it at
  `004a3198` and uses the result only to pick a background. Usable draws
  `[0x005ef954]` (`backinv.bmp`); not usable draws `view+0xc8`
  (`backinvg.bmp`, the same background an empty cell gets, `SHOP-SCREEN-036`
  arm b). The icon is blitted either way at `004a3229`, and no other routine on
  the shop screen calls the predicate.
- Site 2, refusal: `FUN_00492410`, the equipment doll's drop handler, slot
  `+0x7c` of the widget vtable at `0x0059a768`, the routine `ITEM-DOC-053`
  identifies by its store into `[member+0x15c + slot*4]`. It calls the
  predicate at `004924e3` and sets its refusal accumulator `EBP` at `004924ec`
  when the result is 0.
- Four further conditions set the same accumulator: no held element
  (`00492455`); `[[this+0x5c] + 0x140] != 1` (`00492461`); `member+0x74` in
  {3, 7, 8} for slot index 0 or 1 only (`00492489`..`004924a2`); and
  `ITEM-WEAR-058`'s spellbook gate. A sixth condition returns 0 earlier and
  separately: `member+0x14 != campaign+0x9b4` (`0049247e`).
- On refusal the routine calls `[screen+0xe8]->vt+0xa4([screen+0x3d0])` and
  returns 0 at `00492539`. It does not reach `this->vt+0x80` at `00492774`, the
  only call in the routine that turns a drop into a move; does not write
  `member+0x15c + slot*4` (`0049275d`); and does not set `member+0x18c` bit 3
  (`0049273a`).
- The simulation: command `0x22`'s handler `FUN_004d5dd8` and every routine its
  equip destination `cmd+0x0d == 1` reaches (`FUN_004f4d98`, `FUN_004f705b`,
  `FUN_004f7099`, `FUN_00508779`, `FUN_0050c8d0` `Armor::Equip`,
  `FUN_0050d2de` `Shield::Equip`, `FUN_0050def2` `Weapon::Equip`) contain no
  read of `actor+0x4c`, no call to `FUN_0051cb70` or `FUN_00523ab0`, no read of
  `member+0x18c` and no read of parameter 15. The `Data.bin` parameters they
  read are 5, 0xc, 0xd and 0xe.

**Confidence.** High for the two sites and the refusal behaviour: each is a
named instruction in a routine read whole, and `callto:` on a routine reached
by direct `CALL` from both owners has none of `docs/INSTRUMENT.md` rule 7's
blind spots. High for the simulation's silence about class: four disjoint
searches over one handler read whole plus its seven callees. Medium for the
client being the only enforcer anywhere: the search covered the equip path and
the whole-image call site set, not every UI route that could reach the same
command.

### ITEM-WEAR-058

- The base `Item` descriptor arm `FUN_005092c9` sets bits 0, 1 and 2 together
  (`0050930d OR CL,0x7`), then clears bit 1 when `item+0x44 == 5`
  (`0050931e CMP ECX,0x5`, `00509329 AND AL,0xfd`). When `item+0x1c` is `-1` it
  writes the byte as `1` instead (`005092d3`, `005092dc`), leaving an item
  usable by neither class.
- `item+0x44` is assigned by three tests on the resolved `MagicItems` record's
  name CString at `record+0x4`: `Potion` gives 3 (`0050824c`, `00508266`),
  `Book` gives 5 (`0050826c`, `00508286`), `Scroll` gives 4 (`0050828c`,
  `005082a6`); no match leaves the constructor's 0.
- `FUN_0056f10f` is `strstr` returning `-1` for no match, and each arm is taken
  only on an exact 0, so all three are prefix tests.
- Corpus, both roots, 49 rows: 13 `Potion`, 5 `Scroll`, 5 `Book`, 26 neither.
  The five spellbooks `Book Air|Water|Fire|Earth|Astral` are therefore
  mage-only, and the five same-school scrolls are not.
- The second gate: `FUN_00492410` also asks the descriptor for characteristic
  `0x2a` (`004924f7 PUSH 0x2a`, `FUN_00484060`) and refuses unless
  `member+0x18c & 2` and `member+0x1c & (1 << value)` (`00492502`,
  `00492515 SHL EDX,CL`, `00492517 TEST EAX,EDX`).
- The characteristic namespace is the effect kind: `FUN_00509363` pushes
  `effect+0x3c` as the id (`005093d3`). The shipped `Magic` collection's row
  `0x2a` is titled `teachSpell`; row `0x29` is `castSpell`, and that one sets
  descriptor bit 4 at `005093f2`.
- Kind 3 additionally clears descriptor bit 5, the has-effects overlay
  (`00509354 AND CL,0xdf`).

**Confidence.** High for the three prefix tests, the bit clear and the second
gate: each is a named instruction, and the string helper is read whole. High
for the effect-kind namespace, which is the shipped `Magic` collection read on
both roots against the ids the emitter pushes.

**Unknown.** What `member+0x1c` holds and who writes it: it is read here and no
writer was located.

## Effects and weapon-borne casts

| ID | Claim | Confidence | Status | Evidence |
|---|---|---|---|---|
| ITEM-CASTSTATE-056 | A weapon-borne cast borrows two actor fields and preserves the item objects, but it rewrites the weapon-owned `Spell`. | High / Medium / Unknown | ● active (amended, partially retracted) | [EXP-0191](../experiments/EXP-0191-itemcast-training/) |
| ITEM-EFFGRAM-070 | An equipment cell has one first-brace frame and an ordered, failure-isolated effect tail. | High / Unknown | ● active | [EXP-0225](../experiments/EXP-0225-item-effect-grammar/) |
| ITEM-EFFPOP-071 | The complete shipped equipment-cell population is 2,386 cells and 260 braced cells per root, with one malformed nested brace. | Medium | ● active | [EXP-0225](../experiments/EXP-0225-item-effect-grammar/) |
| ITEM-EFFOBJ-072 | An Effect is a 0x48-byte Token-derived value object whose equality is narrower than copy. | High | ● active | [EXP-0225](../experiments/EXP-0225-item-effect-grammar/) |
| ITEM-EFFMODE-073 | The effect operand dword is four grammars, and malformed numeric text is not uniformly rejected. | High / Unknown | ● active | [EXP-0225](../experiments/EXP-0225-item-effect-grammar/) |
| ITEM-EFFSPLIT-074 | Item copy and every split deep-copy the ordered Effect list, but omit one serialized item byte, `+0x47`. | High / Unknown | ● active | [EXP-0225](../experiments/EXP-0225-item-effect-grammar/) |
| ITEM-EFFDISP-075 | All 49 accepted effect keys, their operand widths and their state-0 equip/remove consumers are enumerated. | High | ● active | [EXP-0225](../experiments/EXP-0225-item-effect-grammar/) |
| ITEM-CASTLINK-076 | `castSpell` is both an ordered general-list Effect and the source of a separate Weapon-owned Spell. | High / Medium | ● active | [EXP-0225](../experiments/EXP-0225-item-effect-grammar/) |
| ITEM-EFFSAVE-077 | Effect persistence preserves ordered values but omits transient `+0x44`. | High / Medium | ● active | [EXP-0225](../experiments/EXP-0225-item-effect-grammar/) |

### ITEM-CASTSTATE-056

- Both routes set `actor+0x64 = weapon+0x80` and `actor+0x68 = weapon`, and
  clear both after the attempt. The non-null item pointer suppresses mana use
  and the immediate cast-training call.
- `Unit::Serialize` writes `+0x68` as an archive object reference and `+0x64` as
  a raw `u32`; a live-interval save carries both bit patterns. The world LOAD
  lifecycle explicitly remaps `+0x64` through `00527cb0` and `+0x68` through
  `00527cf0`, each replacing a map hit and clearing a miss (`SAV-908`). Later
  pointer validity and native consumption are untested.
- The attached `castSpell` record supplies byte spell id at `effect+0x40` and
  raw signed-`i16` power at `+0x42`. Every apply passes that same `weapon+0x80`
  object to the common fill, which reloads and power-adjusts serialized
  `Spell+0x09` Max Range and rewrites `+0x0e` damage base, `+0x0f` spread and
  `+0x10` duration scratch.
- The effect record, pointer and object identity survive, and no charge is
  decremented; the four filled fields do not remain unchanged.
- A separate post-cast arm destroys an item of kind `0x0e` and its `Spell`; the
  shipped staff is not that kind.
- Prismatic Spray id 14 separates the routes. The caster admission path applies
  the fan while these fields are live, and its release wrapper prevents
  duplication; the fighter rider reaches only that refusing wrapper.
- G2: spell id width, signed power width, serialized range, raw-pointer
  persistence and the destructive-kind test are engine code. The authored
  Data.bin corpus has 48 weapon-borne cells, 46 Human staffs plus Catapult and
  Ballista, and 18 powers `{1,5,10,15,25,30,34,35,40,50,60,63,65,70,82,90,98,99}`
  per root. It does not impose a 10/20 bound or bound runtime stock: the shop
  additionally generates ids `{1,11,13,14,20}` with price-derived random power
  capped at 100.

**Confidence.** High for field lifetime, suppression, serialized representation
and mutation scope: the two setters, two clearers, common-fill call, all fill
stores, serializer arms, award gate and destroy arm are among 201 raw-byte
assertions, and both direct wrapper callers were enumerated. Medium for the
authored equipment population and kind, which are identical corpus
observations on both roots.

**Unknown.** Raw-pointer validity after load; no runtime load witness was
obtained.

**Amended.** The clause that only `+0x68` is remapped is withdrawn
([`retracted.md`](retracted.md), EXP-0336, `SAV-908`). The world LOAD lifecycle
repairs raw `+0x64` through `00527cb0` and `+0x68` through `00527cf0`, with the
same map-hit replacement and miss clear. Cast lifetime, wire widths, the
destructive-kind branch and the later runtime Unknowns stand.

### ITEM-EFFGRAM-070

- `FUN_004db944` takes the right-trimmed text before the first `{` as the item
  head, then takes bytes after it through the first `}`. No close means
  end-of-string, and text after the first close is ignored.
- `FUN_00502d82` appends a comma, splits left to right and appends every
  non-null `FUN_00502eee` result to a dynamic `CObList`.
- Keys are case-folded and resolved against 50 entries, with index 0 a rejected
  sentinel.
- Missing `=`, unknown key and unknown spell reject only that comma element;
  later elements remain. Duplicate valid keys remain as independent ordered
  nodes.
- No sort, deduplication, count comparison or fixed list budget occurs; a
  synthetic 64-effect tail is accepted.

**Confidence.** High. The splitter and both parser routines are read end to
end; delimiter, lookup-bound and append anchors are asserted against both
executables; opaque-string, unordered-map and fixed-budget alternatives are
directly contradicted.

**Unknown.** Close-before-open and allocator (resource) exhaustion, which are
outside the bounded result.

### ITEM-EFFPOP-071

- Per root: `Humans` 215 rows × 10 = 2,150 cells, 909 non-empty and 258
  braced; `Units` 118 × 2 = 236 cells, 26 non-empty and 2 braced. EN and RU
  corresponding raw cells mismatch 0/2,386.
- Each root parses 266 accepted effects and one rejected token, with zero
  authored indeterminate operands, duplicate-key cells or missing closes. The
  maximum list length is three, in two cells.
- Sixteen of 49 keys occur: mind 1, reaction 1, healthMax 2, toHit 9,
  defence 66, absorbtion 7, speed 5, scanRange 7, protectionFire 35,
  protectionWater 22, protectionAir 37, protectionEarth 20, protectionAstral 2,
  skillPike 1, castSpell 48, damageBonus 3.
- The combined count is 520 braced cells but 522 opening-brace characters: row
  79 slot 4 of `Humans` in each root swallows
  `Rare Mithrill Chain Boots{defence=19` into an unknown second token and builds
  only the Gauntlets with defence 7.
- `castSpell` divides into 46 Human and two Unit cells per root.

**Confidence.** Medium. A complete committed walk of both roots, whose
corresponding cell texts are identical; parity is not independent variation.

### ITEM-EFFOBJ-072

- Kind is u8 `+0x3c`, mode u8 `+0x3d`, and `+0x40..+0x43` is one operand union.
- Copy `FUN_0050124a` creates a new identity and preserves the Token base,
  kind, mode, full operand dword, Token `+0x0c` and transient dword `+0x44`.
- Equality `FUN_005012ab` compares only kind, full operand dword and mode.
  Identity, all other Token fields, Token `+0x0c` and Effect `+0x44` are
  ignored.
- Item equality compares two non-stackable lists in lockstep and requires both
  to end together, so element order, count and duplicates matter while copied
  identity does not.

**Confidence.** High. Constructor allocation, every explicit copy store, the
complete three-comparison predicate and the ordered caller are read directly;
identity and whole-Token alternatives are discriminated.

### ITEM-EFFMODE-073

- Kinds 1..40 and 49 parse a signed-decimal prefix: mode 0/8 uses i32, and
  modes 1/2/4 observe i16 magnitude plus u16 `count<<4` ticks.
- Kind 41 is u16 spell id plus observed i16 power. Its power prefix admits plus
  and digits but not minus, so `:-20` becomes zero.
- Kind 42 is u16 spell id plus optional ticks.
- Kinds 43..48 split at the first hyphen into u8 low and u8
  `(high-zeroExtend(low))`, wrapping both bytes; no hyphen leaves zero.
- Exact `permanent`/`singleuse` produce modes 0/8. Otherwise the searched
  substrings `charges`/`duration`/`continuous` produce 4/1/2, and unknown
  suffix text produces 0.
- Empty decimal is zero and `12junk` is 12. A non-empty prefix for which `%d`
  makes no assignment leaves an unwritten operand without rejecting the Effect.
- Out-of-range input is also accepted.

**Confidence.** High for delimiter priority, field widths, narrowing, mode
selection and a true no-assignment destination: the complete routines and all
load-bearing bytes are read.

**Unknown.** The concrete operand produced by out-of-range `%d` input: this
experiment does not establish whether the executable's CRT assigns a saturated,
wrapped or other value.

### ITEM-EFFSPLIT-074

- Base copy `FUN_00508340` walks source order, allocates a fresh 0x48-byte
  Effect per node and invokes `ITEM-EFFOBJ-072`'s copy.
- Armor, Shield and Weapon copies call it. Weapon additionally creates a fresh
  owned Spell from the source Spell id instead of sharing identity or copying
  filled scratch.
- The four split virtuals decrement source count, allocate their concrete
  class, invoke its copy and set clone count to one.
- `Item::Serialize` stores u8 `+0x47`, but base copy neither copies that byte
  nor runs the default Item constructor before its explicit stores. Clone/split
  `+0x47` is therefore indeterminate, and a later save persists it.

**Confidence.** High for the four split paths, deep-copy identity, Weapon
reconstruction and the copy omission: all named routines and vtables are
complete.

**Unknown.** `+0x47`'s semantic meaning, and its actual value in any particular
allocation.

### ITEM-EFFDISP-075

- Armor, Shield and Weapon equip walk list order through `FUN_00508860`;
  unequip does the same through `FUN_005088ab`.
- A fresh grammar Effect has Token `+0x0c=0`, so `FUN_0050177e` dispatches
  multiplier +1 and `FUN_00501938` multiplier -1.
- The 50-slot table at `0x00502846` has 46 distinct arms; every arm reaches
  target `vt+0x50` recompute.
- `effect-keys.tsv` gives each executable/Data.bin key, operand representation,
  jump arm, class gate/clamp, target field and both-root occurrence count.
- Kinds 38..40 diagnose not implemented; kind 41 is a general no-op; 42 teaches
  a missing Spell; 43/49 share damage-base; 44..48 install elemental
  base/spread.
- Token-state values 8, 12 and 17 divert apply/remove before this general rule
  and are not produced by the equipment grammar.

**Confidence.** High. All table slots are byte-asserted on both images, every
arm and the common tail are read, and both ordered walkers have complete
3-owner call populations.

### ITEM-CASTLINK-076

- General kind-41 dispatch is a no-op before recompute.
- `FUN_005089a7` walks the list and returns the first kind-41 Effect. Its three
  direct owners are Weapon construction and the two cast consumers.
- Construction uses the low spell-id byte at Effect `+0x40` to allocate
  `weapon+0x80`; casting reads raw signed-i16 power at `+0x42`.
- Weapon price special-cases kind 41 through Effect `vt+0x4c`.
- Weapon copy/split deep-copy the source Effect and reconstruct a fresh Spell
  from id.
- The source record therefore stays ordered, comparable and serialized in the
  item list, while the derived object owns cast runtime state.

**Confidence.** High for the dual role: the general no-op arm, the finder, the
three-owner population, Weapon build/copy and the price split rule out
general-only and weapon-tail-only alternatives. Medium for the 46 Human plus
two Unit shipped cells per root, which are identical corpus observations.

### ITEM-EFFSAVE-077

- `Effect::Serialize` writes/reads the 37-byte Token head, u8 kind, u8 mode,
  u32 operand and u8 Token `+0x0c`: 44 bytes of body in that order.
- `FUN_00527b70` writes a u32 list count, then archive object references in list
  order; load clears and appends that count in archive order.
- Eight self-framed new-class first-instance witnesses, one in each of four
  saves per root, all decode as kind 8, mode 1, operand `0x03c00064`, Token
  state 0. They establish a real archived Effect body, not a census of later
  object references.
- Arbitrary high-bit words were rejected as unframed tags.

**Confidence.** High for static writer/reader order and the omitted `+0x44`.
Medium for the bounded eight-witness save population, whose EN/RU files are
byte-identical.

## Authored enchantments, death drops and item producers

| ID | Claim | Confidence | Status | Evidence |
|---|---|---|---|---|
| ITEM-AUTHCAST-086 | The shipped Human `castSpell` authoring population is 46 slot-0 Weapon cells per root, not 46 loose Effects or actor instances. | Medium | ● active | [EXP-0227](../experiments/EXP-0227-authored-enchantments-death-drops/) |
| ITEM-AUTHDROP-087 | In a shipped single-player campaign map, a placed non-`NPC` Human's authored staff reaches its death sack as the same object; an `NPC`-named template's staff is destroyed. | High / Medium | ● active | [EXP-0227](../experiments/EXP-0227-authored-enchantments-death-drops/) |
| ITEM-DRAGON-088 | All four Dragon definitions author only an innate Flame Thrower, and death explicitly withholds it. | High / Medium | ● active | [EXP-0227](../experiments/EXP-0227-authored-enchantments-death-drops/) |
| ITEM-DRAGDROP-089 | The enumerated authored-equipment, type-8, script, direct/reference and shared-death populations contain no Dragon staff or enchanted-item death producer. | Medium | ● active | [EXP-0227](../experiments/EXP-0227-authored-enchantments-death-drops/) |
| ITEM-MAGVAL-090 | `FUN_00508486` initially copies an ordinary `MagicItems` row's parameter 0 to `item+0x1c` as signed `i32`, without a sign branch at assignment. | High / Medium | ● active (amended, partially retracted) | [EXP-0227](../experiments/EXP-0227-authored-enchantments-death-drops/), [EXP-0275](../experiments/EXP-0275-consumable-use/) |
| ITEM-PRODUCER-091 | The repaired `EnumRefs` reference-manager and `.rdata` populations contain three packed-item-factory owners. | Medium | ● active | [EXP-0227](../experiments/EXP-0227-authored-enchantments-death-drops/) |

### ITEM-AUTHCAST-086

- Suffix resolution against `Weapons` identifies every head as `Staff` or
  `Shaman Staff`; both rows carry `sutableFor = 2`.
- The cells belong to 36 Data.bin template names without uppercase `NPC` and
  ten with it.
- `human-staff-templates.tsv` retains all 46, including unplaced definitions,
  rather than selecting only a mission witness.

**Confidence.** Medium. A complete Data.bin walk on both roots, whose
corresponding cell strings are identical; parity is not independent variation.

### ITEM-AUTHDROP-087

- Slot 0 constructs a Weapon and equips it at actor `+0x74`. Death accepts its
  nonzero suitability, unequips it into actor `+0x7c`, then hands that
  container to the sack wrapper.
- The case-sensitive template-name search runs after worn equipment is moved
  and destroys the container on a match.
- The other suppression input is `Player+0x5c != 0` in multiplayer
  (`ITEM-DEATH-012`). All 28 campaign maps author the scalar that clears the
  multiplayer flag (`ALM-MODE-070`).
- Campaign corpus, per root: 76 staff-template placements, 49 eligible and 27
  `NPC`-suppressed.
- Loose-map rows retain their owner and placement joins; their non-`NPC`
  eligibility is unevaluated because player state is not in the map.

**Confidence.** High for the object route and the two suppression inputs,
reusing `ITEM-DEATH-012` and `ALM-MODE-070`. Medium for 76/49/27: a complete
two-root campaign corpus with identical relations.

### ITEM-DRAGON-088

- Data.bin `Units` rows 112..115 (`Dragon`, `.2`, `.3`, `.4`) each carry
  `Flame Thrower` in slot 0, an empty slot 1 and no `castSpell`.
- The matching `Weapons` row has `sutableFor = 0`. Death reads column 15 and
  skips the weapon-to-container move exactly on zero.
- The campaign placements resolve by the loader's `(param29,param30)` pair to
  the four rows, with counts 2/6/4/9 per root.

**Confidence.** High for the death gate and the row resolution. Medium for the
four-row and 21-placement populations, which are corpus measurements.

### ITEM-DRAGDROP-089

- The bounded negative covers all 21 campaign placements per root and every
  loose placement:
  - zero type-8 stock elements;
  - zero type-8 ground elements at the actor cell;
  - zero targets among all item/container instants (7 create, 43 remove,
    3 container-to-sack, 0 transfer);
  - no constructor in the shared death body or in the enumerated
    direct/reference producer population.
- The full roots contain 135 EN / 37 RU Dragon placements because their
  loose-map sets differ; all remain in `dragon-placements.tsv`.
- Flame Thrower is excluded by `ITEM-DRAGON-088`, not counted as item evidence.

**Confidence.** Medium, scoped to the named corpus and producer families:
complete joins plus routine reads and caller enumerations. A computed target
outside Ghidra's reference and `.rdata` populations remains a static blind
spot.

### ITEM-MAGVAL-090

- Scroll and Book are name-derived exceptions inside that routine:
  - Scroll resets, sums spell-derived kind-41 values and falls back to
    parameter 0 only on a zero sum;
  - Book reads only the first Effect's signed `+0x40` as a spell id, assigning
    Spells parameter 21, or zero when empty; there is no Book loop
    (`ITEM-VALUE-115`).
- `Quest Item31` therefore leaves construction at 10,000, and ordinary quest
  rows authored as -1 leave it at -1.
- This is not a global denial of sentinel semantics: the later descriptor path
  gives `item+0x1c == -1` a special arm (`ITEM-WEAR-058`).

**Confidence.** High for the narrow initial signed assignment, asserted on both
executables. Medium for the full Scroll/Book routine reading, whose semantic
listing is not committed.

**Amended.** The Book list-traversal clause is refuted
([`retracted.md`](retracted.md), EXP-0275, `ITEM-VALUE-115`). It had read:
Book overwrites with each effect's spell value in order, ending at the last or
zero when empty. Book gets only the first Effect, prices its signed id through
Spells parameter 21 and exits with no back edge; empty is zero. The Scroll loop
and the initial signed assignment stand.

### ITEM-PRODUCER-091

- `callto:4dcf92` returns 3 hits / 3 owners / 0 orphan: `FUN_004e4f3e`,
  `FUN_00539be0` and `FUN_00484160`.
- The first two are type-8 map load and script instant 12.
- In the read portion of the third body the object formats a UI value and is
  immediately passed to its deleting destructor at `004844f2..004844fc`; no
  actor, container or sack sink was found.
- The same population method gives the sack maker 6 calls / 5 owners /
  0 orphan, and the death inventory routine one direct caller.

**Confidence.** Medium for the owner interpretation and bounded exclusivity:
the hit lists are committed, but the transient body listing is not, and
computed targets remain outside the instrument.

## Enchanted-item trail

| ID | Claim | Confidence | Status | Evidence |
|---|---|---|---|---|
| ITEM-STARFLAG-096 | The normal compact item-display route adds bit `0x20` for a nonzero-picture item with a nonempty Effect list, except that a Potion explicitly loses it. | High / Unknown | ● active | [EXP-0256](../experiments/EXP-0256-enchanted-item-star-animation/) |
| ITEM-STARSURF-097 | The extra item-cell pixels have two direct renderer calls: mission inventory, and one shared shop painter reached by shelf, shown-member backpack and trade-tray wrappers. | Medium | ● active | [EXP-0256](../experiments/EXP-0256-enchanted-item-star-animation/) |
| ITEM-STARPIX-098 | The extra layer is procedural RGB `(255,0,255)`, not archive frames. | High / Unknown | ● active | [EXP-0256](../experiments/EXP-0256-enchanted-item-star-animation/) |

### ITEM-STARFLAG-096

- `FUN_00509229` requires nonzero `item+0x40`, then calls virtual `+0x54`. The
  Item, Armor, Shield and Weapon slots all reach `FUN_00509363`, whose nonempty
  `item+0x20` arm ORs `0x20` into compact-record byte `+4`.
- Base Item then tests `item+0x44 == 3` and clears the bit. `ITEM-STACK-003`
  identifies 3 as the `Potion` prefix result.
- The incoming equipment, inventory and shop builders copy compact `+4` to
  display `+8`, and both cell painters gate the extra pixels on
  `display+8 & 0x20`.
- The writer has 10 direct calls across four owners, and the extension has
  four, one per concrete item class.

**Confidence.** High for the selector, the Potion counterexample and the three
copy paths: 44 exact raw anchors and complete named direct-call populations on
byte-identical EN/RU executables.

**Unknown.** An artificial client record, or a producer outside the normal
compact-record population.

### ITEM-STARSURF-097

- `callto:00482a70` is exactly `004833c6` and `004a323b`; `.rdata` contains
  zero dwords equal to the entry.
- The shop painter has three direct wrappers, at `004a430a`, `004a53db` and
  `004a5ece`.
- Of the four non-check `graphics\inventory\` composers, `00481f27` is the
  shared grid-icon loader, and `004838c0`, `00491815`, `004a3571` build the
  mission, character-pane and shop held-item cursors. Those cursors compose
  only the base `.16a`.
- Worn character/world equipment uses the separate `.256` figure path and has
  no call in the same bounded census.

**Confidence.** Medium for surface exclusivity. Exact direct-call,
initialized-pointer, vtable and repaired function populations account for the
named paths on both executable roots, but a runtime-computed target stored in
neither population remains possible.

### ITEM-STARPIX-098

- Every grid constructor builds 1,024 x/y pairs as `rand()/0x1ff+8` in instance
  tables, hence coordinates 8..72.
- The shared CRT source is the LCG `state=state*0x343fd+0x269ec3`, returning
  `(state>>16)&0x7fff`.
- Slot phase `q` selects `(q>>1-age)&0x3ff` and introduces one coordinate every
  two increments. After startup seven centres use alpha
  `63,127,191,255,191,127,63` and repeat coordinate selection after 2,048
  increments.
- Each centre is a 13-point plus/diamond kernel: centre `a`, axial distance one
  `floor(3a/4)`, axial distance two `floor(a/2)`, diagonals `floor(a/4)`.
- Its item-relative footprint is 6..74 inside the 80x80 icon.

**Confidence.** High. Exact constructor, RNG, renderer and kernel bytes, 3/7/13
complete direct-call populations and executable re-derivation agree on both
roots; stored-frame and static-frame alternatives are discriminated.

**Unknown.** The exact seed, and the intervening CRT draws before a grid is
constructed.

## Enchanted-item trail: phase and placement

| ID | Claim | Confidence | Status | Evidence |
|---|---|---|---|---|
| ITEM-STARPHASE-099 | The trail phase belongs to a visible grid slot, not to an item, and the mission and the shop advance it on different clocks. | High / Unknown | ● active | [EXP-0256](../experiments/EXP-0256-enchanted-item-star-animation/) |
| ITEM-STARCOMP-100 | Trail geometry is common and icon-relative, but its placement and occlusion are surface-specific: the mission inventory and the shop differ. | High / Unknown | ● active | [EXP-0256](../experiments/EXP-0256-enchanted-item-star-animation/) |

### ITEM-STARPHASE-099

- `FUN_004829a0` increments the phase of each eligible visible slot. It has
  two callers: the `0x401` message arm and shop paint.
- The mission advances the phase on `0x401`. Ordinary pause posts no `0x401`
  (`SESS-PAUSE-020`), so pause holds it.
- The shop advances it once per paint, with no pause-bit test.
- A selected singleton is omitted from the grid, and its held cursor has no
  trail. A remaining stack continues normally.
- Scroll-down shifts phases toward slot zero and clears the new final slot.
  Scroll-up shifts them upward but leaves phase zero unchanged. This
  asymmetry directly refutes a phase keyed by item identity.

**Confidence.** High for the two static clocks, the mission pause join, the
selected-singleton rule and the asymmetric slot shifts: exact bytes plus the
complete two-call population.

**Unknown.** Actual shop message and repaint scheduling, and presentation while
paused: the shop's repaint cadence and paused-shop coexistence remain runtime
Unknowns.

### ITEM-STARCOMP-100

- Mission origins are `x=left+32+((screenW-240)%80)/2+80*i`, `y=top+6`. The
  mission inventory rectangle is `(0,H-90,W-160,H)`, which gives
  `(32+80i,396)` at 640x480.
- Shop origins are `view.left+cell.left,view.top+cell.top+1`.
- The point blender clips to the renderer-global half-open rectangle
  `[005e4408,005e4410) x [005e440c,005e4414)`. The trail installs no per-cell
  clip; its 6..74 footprint stays inside the 80x80 icon.
- Layer order differs by surface. Mission: base icon, quantity, trail. Shop:
  background, base icon, trail, quantity, then later price layers. One
  universal layer order is therefore false.

**Confidence.** High for the arithmetic, the clip accesses and the instruction
order on the two positive paths, with exact raw anchors on both roots.

**Unknown.** The active global clip values at every caller, and runtime pixel
occlusion.

## Whole-item movement

| ID | Claim | Confidence | Status | Evidence |
|---|---|---|---|---|
| ITEM-WHOLE-128 | A whole extraction and a non-merging insertion move the existing Item pointer; they are not a copy. | High / Medium / Unknown | ● active | [EXP-0285](../experiments/EXP-0285-whole-item-transfer/) |
| ITEM-MERGE-129 | A container merge retains the destination Item, not the incoming Item. | High / Medium | ● active | [EXP-0285](../experiments/EXP-0285-whole-item-transfer/) |
| ITEM-GROUNDMOVE-130 | New-Sack adoption, existing-Sack draining and pickup have different identity and flag outcomes. | High / Medium / Unknown | ● active | [EXP-0285](../experiments/EXP-0285-whole-item-transfer/) |
| ITEM-EQUIPMOVE-131 | Armor, Shield and Weapon equip move the existing Item pointer, but equipment pointer movement does not imply that all nested state is preserved. | High / Medium / Unknown | ● active | [EXP-0285](../experiments/EXP-0285-whole-item-transfer/) |
| ITEM-SPELLMOVE-132 | Weapon equip and unequip change the identity of the Weapon-owned Spell without an Item split. | High / Unknown | ● active | [EXP-0285](../experiments/EXP-0285-whole-item-transfer/) |

### ITEM-WHOLE-128

- `0050ebf9` takes the whole pointer when the requested quantity is at least
  the Item's u16 count (`0050ec52 -> 0050ec86`). Count one with request one
  never reaches split `vt+0x40`.
- Nonmerge `0050e902` stores that pointer in a list node and updates the
  container load.
- The inspected take/add/list bodies contain no Item-receiver store to the
  serialized Token fields or to the Position reached through Item `+0x10`.
- `0050eaaf` drains quantity one repeatedly, so a count above one can still
  split before insertion.
- This bounded helper result does not establish a complete transitive write set
  for command, equipment, sender or actor callbacks.

**Confidence.** High for the branch, the pointer flow and the node/count
stores, verified against both original PE images. Medium for the bounded
no-store result.

**Unknown.** Transitive callbacks, and aliased or malformed graphs.

### ITEM-MERGE-129

- `0050e902` requires incoming stackability, equal item codes and candidate
  stackability.
- It then writes destination count = old destination count + incoming count
  (`0050e9b2`), writes destination Token `+0x08` = old destination flags OR
  incoming flags (`0050e9c5`), and deleting-dispatches the incoming object at
  `0050e9e4`.
- Its direct merge body does not copy the incoming Position or `+0x14` into the
  retained Item.
- Same-pointer extraction is therefore insufficient to establish final Item
  identity, even with count one.
- The rule concerns the named valid-object branch, not pointer aliases or
  deletion side effects through other references.

**Confidence.** High for the predicate, the destination stores and the incoming
deletion. Medium for the bounded absence of other retained-head stores beyond
the inspected body.

### ITEM-GROUNDMOVE-130

- New Sack: `0050f360` stores the supplied container at Sack `+0x40`
  (`0050f39e`).
- Existing Sack: `0050f5aa` instead calls `0050eaaf` when a Sack already exists
  (`0050f6a1`), draining and deleting the supplied container. Death and DropAll
  reach that same branch.
- Pickup writes every incoming Item's Token `+0x08 = 1` at `004f4e9b` before
  draining. A pickup merge therefore leaves the retained destination flags =
  old destination flags OR 1, not OR the incoming pre-pickup word.
- The new Sack allocates and copies its own Position through
  `004f24d4 -> 00544980`; this does not set the contained Item Position.
- The GiveAll/XferItem direct bodies use drain or take/add without the pickup
  stamp.

**Confidence.** High for the branch, the adoption pointer, the Item flag source
and the Sack Position receiver. Medium for direct-body Item Position/reference
nonmutation.

**Unknown.** Registration/sender transitive effects, and actual runtime reach.

### ITEM-EQUIPMOVE-131

- Armor/Shield/Weapon equip and unequip use the existing Item pointer while
  writing actor slots and actor stats. Displaced pieces may merge on
  reinsertion.
- These direct item-class bodies contain no serialized Item Token/Position
  store, but call Effect `+0x40/+0x44` and the actor recompute.
- Base Item equip is different. With positive actor HP it changes zero Effect
  mode `+0x3d` to 8, dispatches Effect attach, then deletes the Item. Its
  unequip is an error path.
- Effect attach can copy into the actor list or change an existing actor
  Effect's operand.
- State8 apply has a referenced-actor callback. State0 dispatch and the actor
  callbacks prevent a universal nested/head no-write conclusion.

**Confidence.** High for the resolved virtual targets, the slot identity flow
and the explicit consume-mode/deletion stores. Medium for the direct-body
no-Item-head-store boundary.

**Unknown.** Complete transitive effects, and aliases.

### ITEM-SPELLMOVE-132

- Equip `0050def2 -> 0050e7a6` finds the first kind41 Effect. If one is found,
  `0050e6f3` deletes any old Spell and stores a new 0x14-byte Spell at Weapon
  `+0x80` (`0050e790`), or null on allocation failure. With no matching
  Effect, that helper writes no Spell field.
- The id is the low byte of Effect `+0x40`. For a valid nonzero id,
  `004fdd96 -> 004fdf57` writes:
  - Spell `+0x08` from that id;
  - `+0x09` from the low byte of Spells row parameter6, after the constructor
    default1;
  - `+0x0a` from parameter18 == 1, after default0;
  - `+0x0c` from the low word of parameter1.
- `005006e3` serializes these four values plus `this`.
- Unequip deletes a nonnull Spell at `0050e38e` and zeroes Weapon `+0x80` at
  `0050e3a0`, without removing the source Effect.
- The Item allocation survives, but its nested archive graph changes. Freed
  addresses may be reused, so numerical key inequality is not required.

**Confidence.** High for the constructor, delete and store paths and their
sources.

**Unknown.** Invalid database ids, allocator failure recovery, runtime
occurrence and transitive callback effects.

## Sack visual frame

| ID | Claim | Confidence | Status | Evidence |
|---|---|---|---|---|
| ITEM-136 | The client message dispatcher's opcode-`0x7a` case writes `CBackPack+0x20`, the ground-sack painter's frame field; that case is not shown to be the field's only writer. | High / Medium | ● active | [EXP-0379](../experiments/EXP-0379-sack-frame-ladder/) |
| ITEM-137 | The sole opcode-`0x7a` producer found in `.text` is `FUN_004e8eb3`'s Sack branch, which writes the frame field from `Sack+0x1c` through `log10`, not from weight or item count. | High / Medium | ● active | [EXP-0379](../experiments/EXP-0379-sack-frame-ladder/) |
| ITEM-138 | For the producer `ITEM-137` found, the frame index ascends: `log10` is monotonic increasing over `Sack+0x1c`'s positive domain, so a larger recomputed value never yields a lower frame index. | High / Medium / Unknown | ● active | [EXP-0379](../experiments/EXP-0379-sack-frame-ladder/) |

### ITEM-136

- `CBackPack+0x20` is the field `SPR256-077` shows the ground-sack painter
  reads at both its draw calls. This write closes the
  server-Sack/client-CBackPack join `ANIM-CATEGORY-084` left unestablished.
- `ANIM-072`'s dispatcher `FUN_004104e8` pops from `[0x5f22d0]` and dispatches
  on the popped message's `+0x9` byte. A replay of its own two-level table lands
  raw opcode `0x7a` on case index 16, block `0x416a64`..`0x417006`. The table
  is:
  - a range check of opcode−3 against limit 187;
  - a byte remap table at `0x41879b`;
  - a dword jump table at `0x4186e3`.
- The block's fall-through target (`0x417006 -> 0x418692`) is `ANIM-072`'s own
  shared match tail.
- Both arms of the block write `CBackPack+0x20 = (byte)msg+0xd`, then clamp it
  to a signed maximum of 5 (`CMP …,5 / JLE`). No lower clamp is needed: the
  source is a zero-extended byte.
  - Existing-record arm: `00416b7f`..`00416b9f`.
  - Newly-constructed-record arm: `00416faf`..`00416fd5`.
- The newly-constructed-record arm's CBackPack comes from `FUN_00459cc0`, whose
  sole direct caller anywhere in `.text` is this arm (`00416bd4`).
- The same arm calls a nine-argument helper, `FUN_00459240` (`00416f95`), that
  stores its own second argument, a literal `1`, to `CBackPack+0x20` at
  `00459294`, eleven instructions before the clamped write overwrites it. This
  second, earlier write inside the same case is harmless to the drawn frame.
- A full-CFG scan of the dispatcher's whole 7325-instruction body finds nine
  stores of the form `mov [reg+0x20], src`: the four named above inside case
  `0x7a`, and five elsewhere, at `00411746`, `00412c90`, `004174a4`,
  `00417cdc` and `00418108`, whose object class EXP-0379 did not verify.
  - `00412c90` falls inside the address range `UNIT-STRUCTCONT-078` (promoted)
    names as the client's own opcode-`0x82` Building arm
    (`004125d2`..`00412cf2`), a different message tag entirely.
  - The other four are unchecked either way.
- The producer side of this join, what writes `msg+0xd`, is `ITEM-137`. This
  claim extends `ANIM-072` rather than restating it: that claim's own Unknown
  names `FUN_004e9e7c`/`FUN_004e8eb3` as unread for server-broadcast transport.

**Confidence.** High for the case resolution, a programmatic replay of the
dispatcher's own tables rather than a hardcoded case address; for the
write/clamp instructions in both arms; for the sole-caller fact for the
no-argument CBackPack constructor; and for the helper's own harmless write.
Medium for "the field's only writer": the full-CFG scan of the dispatcher's own
body finds five further `[reg+0x20]` stores outside case `0x7a`, one inside a
different, already-published message arm and four of unverified class. The
search is bounded to one function, not a proof over every path that can reach
a CBackPack instance.

### ITEM-137

- `FUN_004e8eb3` is the generic notify routine, the one `ANIM-072` names as its
  own unread Unknown for server-broadcast transport.
- `go run ./tools/claim -k "4e8eb3"` (`evidence/corpus-searches.txt`) returns
  16 rows. Excluding `ANIM-072` itself and EXP-0379's own three, three of the
  remaining twelve read inside this routine's body for other purposes:
  - `SAV-POSTLOAD-221`: its Sack- and Building-class arms' handling of the
    Position field;
  - `SAV-678`: its null-recipient broadcast publish step;
  - `UNIT-STRUCTCONT-078`: its Building-class arm's HP-field message write.
- The other nine (`ANIM-071`, `ANIM-BLOW-019`, `SAV-POSTLOAD-223`,
  `TRIG-DROPALL-024`, `TRIG-PROPERTY-036`, `TRIG-CHECK-053`,
  `UNIT-STRUCTUSE-090`, `UNIT-STRUCTUSE-091`, `UNIT-STRUCTZERO-080`) cite it
  only as a call target, without reading its internals. None of the twelve
  reaches its Sack branch or its opcode/`log10` computation.
- The Sack branch is reached only after, in order:
  1. a null/non-null recipient-argument split (`ebp+0xc`, `004e8ebc`);
  2. a null-object early exit (`004e8f1a`..`004e8f3d`);
  3. a virtual call through the object's own vtable slot `+0x2c` (`004e8f4a`),
     where only a zero return proceeds (`004e8f4f`). Sack's own vtable
     (`0x59ca88`, `ITEM-SACK-010`) has an eight-instruction stub at that slot
     (`0x522f60`) that always returns 0;
  4. the class test `0x57272f` (`mov eax,[ecx]; call [eax]`, a virtual
     `GetRuntimeClass`-shaped call, then `call 0x5727e6`, an
     `IsDerivedFrom`-shaped tail) against `0x5c33e8`, `ITEM-SACK-010`'s Sack
     `CRuntimeClass` (`004e8fdc`..`004e8fe9`).
- On a match it stamps `msg+9 = 0x7a` (`004e9001`), then computes `msg+0xd`:
  1. `FILD dword ptr [Sack+0x1c]` (`004e905e`);
  2. a call to `0x557374` (`004e9067`), the `fldlg2`/`fxch`/`fyl2x` log10
     idiom, confirmed by reading that callee's whole CFG: both error branches
     and both exit tails, not only the happy path;
  3. a call to `0x55458c` (`004e906f`), the MSVC `_ftol` truncate-to-integer
     helper, confirmed by its `fnstcw`/`or ah,0xc`/`fldcw`/`fistp`/restore
     shape;
  4. `MOV [msg+0xd],AL` (`004e9077`).
- No weight field and no item-count field is read anywhere in this routine. The
  only Sack-owned input is `+0x1c`, `ITEM-SACK-010`'s recomputed total value.
  This refutes a weight-keyed or count-keyed selector.
- It also refutes a literal `CMP`/`JL` threshold cascade for the six frames:
  `ITEM-136`'s case handler clamps `_ftol(log10(Sack+0x1c))`, a smooth
  `log10` truncated, not a table of six literal breakpoints. The result is a
  six-band ladder at powers of ten, not a smooth curve
  ([`formats/item/sacks.md`](../formats/item/sacks.md)).
- A per-byte-offset decode of the whole `.text` section (1 660 928 offsets) for
  the direct-immediate `mov byte ptr [mem], 0x7a` form, with any
  ModRM/SIB/displacement encoding rather than only the narrower `C6 /0`
  pattern, finds one hit, `004e9001`. The same pass found 6 019
  `mov byte ptr [mem], reg8` sites. A further scan of each for a
  same-register-family `mov reg,0x7a` within the preceding 64 bytes flagged
  none.
- `log10`'s zero/subnormal-zero branch does not compute an ordinary logarithm:
  at `Sack+0x1c = 0` it discards the operand and returns the extended-precision
  negative-infinity bit pattern with an internal error code 2
  (`formats/item/sacks.md`'s domain note).

**Confidence.** High for the reached instructions: the full entry-to-gate chain
(recipient split, object-null exit, vtable `+0x2c` zero-return gate, class
test), the opcode stamp, the complete `log10` function (its whole CFG, both
error branches and both exit tails) and `ftol` are all read at instruction
level. Medium for "the only producer": the per-byte-offset scan supersedes the
narrower `C6 /0` pattern and still finds one site, and the register-mediated
lookback flags none of 6 019 candidates. Neither scan can catch a register
value arriving by arithmetic, by a table, or by a load further back than 64
bytes, so a register-mediated producer is narrowed, not excluded.

**Unknown.** Whether a zero-value Sack occurs in play (`log10`'s zero branch,
read in full) was not traced.

### ITEM-138

- No sign flip, table reversal or subtraction from a fixed frame count sits
  anywhere between the `fyl2x` result and the clamped store `ITEM-137` reads.
  The stored frame index only rises or holds as `Sack+0x1c` rises, in six
  bands at the powers of ten
  ([`formats/item/sacks.md`](../formats/item/sacks.md)).
- This claim establishes the index direction only. Whether a higher index draws
  a visibly bigger sack sprite depends on the backpack sheet's own frame
  ordering, which EXP-0379 did not read: the sheet's frame-to-appearance layout
  is asset-domain and out of the scope declared for EXP-0379.
- One concrete trigger path from a value change to this notifier call is read
  end to end. `FUN_0050f6d2` (`ITEM-DEATH-012`'s wrapper, used by death and
  both drop arms) calls the shared create/merge `FUN_0050f5aa`, whose own call
  census includes `ITEM-SACK-010`'s cell lookup (`0x547c60`) and value
  recompute (`0x50f4b3`). On success the wrapper repeats the cell lookup
  against the sack registry `[0x5f22c8]` and calls `FUN_004e8eb3` directly
  (`0050f70a`).
- `TRIG-DROPALL-024` (active) independently publishes the same
  wrapper-through-create/merge-through-notify chain for its own drop-all
  purpose, which corroborates that this is a real, reached path rather than
  dead code.
- A full `.text` scan for `E8 rel32` sites targeting `0x4e8eb3` finds 13 direct
  callers. `SAV-POSTLOAD-221` (promoted) separately reads one of the other
  twelve, `004ea094` inside `FUN_004ea059`, for a save-load Position-sync
  purpose unrelated to the frame field. So 2 of the 13 direct callers are
  addressed by some published claim (this one and `SAV-POSTLOAD-221`), and 11
  remain wholly untraced by any claim this search found.

**Confidence.** High for the monotonic direction on the traced path and for the
one complete trigger-path read. Medium overall. The finding inherits
`ITEM-137`'s Medium: it holds only for the producer `ITEM-137` found, and is
not established for a register-mediated producer of opcode `0x7a`, which that
scan narrowed but did not exclude. It is also conditioned on whether the
owner-reported visual size change follows the same direction as the decoded
index, which EXP-0379 did not measure.

**Unknown.** Whether the 11 untraced callers reach the notifier by the same
value-change path or a different one, and what `log10`'s error-handler
branches (read in full, see `ITEM-137`) produce as a stored byte when
`Sack+0x1c` is non-positive.

## Open questions

- Whether the interface ever emits the session-space source-3 form of the
  pick-up, the one of its two entry points that requires the sack already
  underfoot (`ITEM-PICK-016`). It is not established because who enqueues a
  command is not (`SESS-CMD-015`).
- Any producer of `item+0x44` outside the bounded item constructor family
  (`ITEM-STACK-003`), and the semantic meaning of serialized `item+0x47`.
  `ITEM-EFFSPLIT-074` establishes that copy omits `+0x47`, leaving its concrete
  clone/split value indeterminate; a running observation could measure one
  value but would not name the field.
- Whether any route other than `Armor::Equip` ever writes `actor+0x19c` or
  `actor+0x1a0`, armour parts 1 and 2 (`ITEM-ARMSLOT-031`).
- The names of the `Armors`/`Shields` columns other than column 4, `Slot`
  (`ITEM-HUMEQ-030`).
- A runtime session. Everything in this ledger is static. Three
  stopwatch-free observations would be independent witnesses: drop a stack two
  cells away and see where it lands, drop three cells away and see it land
  underfoot (`ITEM-DROP-008`), and overload a hero to watch the speed floor at
  6 (`ITEM-LOAD-005`).
- What `member+0x1c` holds on the client unit, and who writes it.
  `FUN_00492410` reads it as a 32-bit mask indexed by the value of a
  `teachSpell` effect (`00492512`, `00492515`), so a spellbook is refused
  unless the corresponding bit is set (`ITEM-WEAR-058`). No writer was located:
  the client unit range `0045a000..00461000` contains no store to that
  displacement. The discriminator is the network opcode that carries it;
  `HERO-APPEAR-048`'s pair carry equipment, not this.
- Which actions the three codes that bar a weapon or shield swap stand for.
  `FUN_00492410` refuses slot indices 0 and 1 when `member+0x74` is 3, 7 or 8
  (`00492490`..`004924a0`, `ITEM-WEAR-057`). `member+0x74` is written from the
  pending-action byte `member+0x84` (`0045d5e4`) and zeroed at `0045cf75`; the
  code set itself is unenumerated. Armour slots are not gated this way.
- What `FUN_005095ee`, `FUN_0050a854` and `FUN_0050ba63` do with `sutableFor`
  (`ITEM-SUIT-035`). Nine instructions read parameter 15; three are the
  descriptor arms and one is the corpse-drop gate. These three test `AND 1`
  alone. `FUN_0050a854` is passed as a function pointer at `0050a97b`, so it is
  a comparator; the other two are unread. A fifth reader, `FUN_0050b6e3`, is
  unreachable: 0 dword occurrences in the whole file and 0 direct relative
  calls in `.text`, against a control that scores 1.
- What `[screen+0xe8]->vt+0xa4` does. It is the only call the refusal path
  makes before returning 0 (`ITEM-WEAR-057`). Identifying it would settle
  whether a refused drag is cancelled or continues.
