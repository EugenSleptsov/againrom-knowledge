# SAV Unit, Humanoid and Human

[Format reference](format.md) · [Token](token.md) · [Human state](human-state.md)

## Unit programme

`Unit::Serialize` is `00510518`. Each row follows the previous one without
padding. Sources are offsets in the actor unless an indirection is shown.
`list32` and `Count` are different primitives. — SAV-UNITPROG-156

| Order | Wire | Source |
|---:|---|---|
| 1 | Token37 | Base Token |
| 2 | `list32<Effect>` | `+20` |
| 3 | `Count(n)`, `n*u16` | Static route at `+15c` |
| 4 | `Count(n)`, `n*u16` | Dynamic route at `+178` |
| 5 | raw 24 | `+a6` |
| 6 | raw 22 | `+be` |
| 7 | raw 24 | `+114` |
| 8 | raw 64 | `+d4` |
| 9 | raw 180 | Mover at `*(+154)` |
| 10 | raw 148 | Order at `*(+158)` |
| 11 | `Count(n)`, `n*u16` | List at `*(*(this+158)+90)`, inside order serializer `005390c0` |
| 12 | u8 each | `+49,+4a,+4b,+4c` |
| 13 | raw 4 each | `+50,+54,+58` |
| 14 | u8 each | `+60,+61,+6c` |
| 15 | Two `objref` operations | `+74,+78` |
| 16 | CString | `+80`, name |
| 17 | u16 each | `+84,+86,+88,+8a,+8c,+8e,+90,+92,+94,+96,+98,+9a,+9c,+9e` |
| 18 | u8 each | `+a2,+a3` |
| 19 | u16 each | `+a0,+a4` |
| 20 | u8, u32, three u8 | `+12c,+130,+134,+135,+136` |
| 21 | u32, u8, u32, u32 | `+138,+13c,+148,+144` |
| 22 | `objref` | `+68` |
| 23 | u8 presence; if nonzero, direct container | `*(+7c)` |
| 24 | u8 presence; if nonzero, direct Spellbook | `*(+140)` |
| 25 | Four u32, then u8 | `+5c,+64,+44,+40,+48` |

The six raw block writes in rows 5..10 total 462 bytes. SAVE first copies
`u8 +14c` into `u32 +148`; this is a local SAVE mutation. Presence producers
emit 0/1 while the load arms accept every nonzero byte. The inventory presence
flag stores no container identity: a preconstructed container receives its
count, references, insertion index and load. — SAV-UNITPROG-156,
SAV-CITYSTORE-516, SAV-DEADLOAD-130

Humanoid appends raw XP24 at `+1cc`, twelve equipment `objref`s at
`+198+4*i` for `i=1..12`, then a Diary `objref` at `+1e4`. Human inherits
that whole programme and adds no bytes in its own store arm. Inventory,
equipment, Spellbook and Diary are separate constructs.
— SAV-HUMAN-043, SAV-CARRY-050

The actor Diary reference may resolve to an existing archive object. The
typed LOAD does not establish exclusive ownership; a new object is registered
before its virtual serializer. The bounded actor-method search reaches SAVE
and destruction but no ordinary independent Diary array consumer. — SAV-982,
SAV-984

Actor `+1e4` also equals `+1cc+4*6`. The two selected progress branches lack
a local index<=5 check, so index6 would address the Diary pointer word as a
scalar. Native reach of that index and a resulting pointer change remain
Unknown. It defines neither a seventh skill nor Diary array meaning. — SAV-983

## Scalar meanings

| Runtime field | Wire width | Established meaning |
|---|---:|---|
| `+84,+86,+88,+8a` | 2 each | Body, Reaction, Mind, Spirit |
| `+8c` | 2 | Speed |
| `+8e` | 2 | Actor's own carried weight |
| `+90` | 2 | Derived load: own weight plus half the container's running weight |
| `+92` | 2 | Capacity |
| `+94,+96,+98` | 2 each | Signed health, maximum and regeneration period |
| `+9a,+9c,+9e` | 2 each | Signed mana, maximum and regeneration period |
| `+a2,+a3` | 1 each | Separately saved regeneration remainders, unsigned on reload |
| `+a4` | 2 | Sight radius in 1/256-cell units |
| `+12c` | 1 | Range, including equipment delta |
| `+130` | 4 | Aggregate experience, simulation i32 |
| `+13c` | 1 | Actor stage |
| `+6c` | 1 | Signed timer; countdown while dying |

The archive restores weight, load and sight without recomputation. The named
Unit constructor capacity is based on default Body 30 and stays 300 on that
arm; that is not a general relation to restored Body. The direct weight-delta
helper adds a 16-bit delta and recalculates load. Its 64000 threshold is a
signed compare; a wrapped negative load removes the speed penalty in the
located arithmetic; native acceptance of such loaded values remains Unknown.
— SAV-UNITFLD-049, UNIT-CTOR-004, SAV-635, ITEM-LOAD-005, HERO-SIGHT-007,
AI-SIGHT-092, SAV-792, SAV-793, SAV-REGENWIDTH-528, SAV-REGENWIRE-532

A stored load of 181 with own weight 178 and an empty container can survive
original load/resave; recomputing would destroy that state. Sight is not
restricted to whole cells. Its named constructor gives `0x0500`; LOAD can retain
zero without that becoming a general constructor/default value. — SAV-794, SAV-795, SAV-796

## Skills and experience

| Source | Wire position | Meaning |
|---|---|---|
| `+a8+2*i`, `i=0..5` | Bytes 2..13 of raw-a6 | Six u16 levels |
| `+130` | Scalar row 20 | Aggregate XP, simulation i32 |
| `+1cc+4*i`, `i=0..5` | First 24 bytes of Humanoid suffix | Six XP values, simulation i32 |

The level word and XP dword share index i. Initializers and award paths update
them and the aggregate, but LOAD copies all three independently. The word at
`+a6` is not level slot zero. `aggregate=sum(XP)` is not enforced by LOAD. The original programme fixes six slots; extending
the XP tail shifts the following references and requires an external versioned
extension. — SAV-HEROXP-063, SAV-HEROSKILL-064

Weapon-borne spell awards use the same fields; no separate item-cast skill or
award-source tag is stored. Progress belongs to the containing Human object,
not its Group position. Player `+34` selects a primary character without
changing other members' shape. — HERO-ITEMSKILL-096, SAV-HEROID-065

## Definition and presentation binding

The restored row byte `+0c` and type word `+0e` are separate selectors. Let
`U=[00609ba8]` and `H=[00609bbc]` be Units/Humans array bases.

| Exact loaded class | Definition pointer `+3c` after local serializer return |
|---|---|
| Unit | `U + 48*u8(+0c)` |
| Humanoid | 0 |
| Human | `H + 48*(u16(+0e)<33 ? u8(+0c) : 5)` |

The arithmetic wraps at 32 bits and has no collection-size check. The fixed
Human row-5 branch covers every zero-extended word 33..65535. These branches
do not rewrite the selectors; Item has a different size-checked lookup.
— SAV-ACTORBIND-544

Unit copies lowByte(`+148`) to `+14c`; exact Unit then clears the byte,
Humanoid preserves it, and Human retains a nonzero byte only with a nonnull
definition whose name contains `NPC`. The substring need not start the name;
the suffix preserves values such as 2/255 rather than Booleanizing them.
The dword remains unchanged until later SAVE widens the effective byte into
it. A sender uses the byte to select a separate collection at `005f0758`.
Other character modes and full runtime remain outside this local result.
— SAV-ACTORDISPLAY-545

Archive creation invokes constructors first. Human's creator uses
`004f9065("Man_Unarmed",0,0)`; constructor definition/equipment/type setup is
not rerun by the load suffix. Actor `+154/+158` allocation identities are
fresh; their 180/148 bytes come from SAV, and the order helper replaces its
`+90` list. These constructor outputs are not post-load defaults.
— SAV-ACTORCTOR-546

Footprint/domain `+49/+4a`, face/class `+4b/+4c` and mover mask `+5` restore
independently. That mask selects the passability plane. — TERR-PASS-051

The selected post-read hooks do not normalize their consistency.
Stage zero admits key repair in the fresh order object for
`+0c,+10,+18,+20,+28,+30,+68`; hits replace, misses retain raw words, without
class/lifetime checks. Other stages skip that repair. — SAV-ACTORINPUT-547,
SAV-HUMRESUME-460

Exact Humanoid's null definition is not established as safe: its shared Unit
`vt+58` uses `+3c+8` without a local null guard. Exact-class acceptance, loaded
reach and the first frame/move/save chronology remain Unknown.
— SAV-ACTORLIMIT-548

## Mover and first event

Mover `0054d4c0` transfers 180 bytes without field normalization, including
current/desired bytes `+0/+1`, rotation speed `+a`, counter byte `+9d`, active
dword `+a0` and estimate byte `+a4`. Unit passes its `+154` receiver and the
same archive. The selected Unit/Humanoid/Human serializers contain no later
direct replacement of rotation-speed byte `+a`; their virtual/embedded calls
and later Human derive are separate mutation boundaries. — SAV-TURNLOAD-822,
SAV-MOVRATE-866, MOVE-TURN-044

World LOAD rebuilds Player and global actor membership using the same actor
pointers, excluding actor `+4c` mask `0x08`. The insertion helper `0050fbee`
does not include the creator helper's `+4` assignment. Post-load `actor+24`
uses the base hook; stage zero reaches mover `+7c` reference repair. These
selected bodies provide no unconditional derive-before-read guarantee.
— SAV-LOADREG-878, SAV-LOADHOOK-879

Common tick processes attached Effect `+38` before HP/order admission.
The selected pending-order-10 turn requires positive signed HP, nonzero
`+3c` and its admitted order prefix. Its inactive snap skips the allocated
byte read; an admitted continuous Effect can call actor `+50` first.
Frontend, world, phase and command callbacks precede the selected actor suffix,
so the first actual restored mover event, its value and native timing remain
Unknown. — MOVE-EVENT-060, MOVE-EVENT-061, SAV-FIRSTMOVE-880

## Dead actors and source references

Dead actors leave the owner graph and remain in the exact top-level dead list.
LOAD creates/resolves them through CArchive, then the later manager callback
runs `actor->vt+24`, repairing Position and `+5c,+64,+44,+68,+40` references.
It does not join an authored ALM record or derive actor state. Timer `+6c`,
signed health `+94` and stage `+13c` remain separately restored fields.
— SAV-DEATH-051, SAV-DEADLOAD-124, SAV-DEADLOAD-125, SAV-DEADLOAD-126

Stage-5 health is not restricted to -10001: the decay consumer retains every
signed value below -10000. Producers of the other observed residues remain
Unknown. Death moves/adopts or drains the old inventory into a Sack and gives
the corpse a fresh container; the saved corpse container is not the Sack
container identity. — SAV-DEADLOAD-127, SAV-DEADLOAD-129, SAV-DEADLOAD-130,
ITEM-DEATH-012

Unit `+40` is a raw old actor-address key; `+64` is temporary attribution;
`+68` is initially an archive object reference; signed `+48` is an attribution
byte. World lifecycle separately map-repairs `+40/+64/+68`, nulling misses;
no-world LOAD clears `+40`, and `+48` is retained at that boundary. The `+68`
second lookup uses its then-current word, so archive resolution alone does
not prove its survival. Ordinary item-cast completion clears `+64/+68`.
— SAV-908, MAGIC-ITEMKILL-117, ITEM-CASTSTATE-056

### Consequence receivers

The selected Unit callbacks `+48/+60/+64` are empty. Humanoid/Human use
`004f7a0c/004f783f/004f78cc`; the latter two keep the source actor as the
receiver of progress slot `+5c ->004f72d7`. Constructor and vtable provenance
identify these classes. A restored address-map hit or actor type word does
not validate an arbitrary receiver's class or allocation lifetime. — SAV-958

The source callback pair belongs to actor manager `0050fca9`, after its
admitted actor tick/removal work. Humanoid/Human `+48` can add victim `+1c`
to source-owner Player `+38` and enter notification dispatch. It bypasses that
helper for server `+0c==0`, a nonzero victim class predicate, or negative
source HP. The manager then rereads signed victim `+48` and victim `+40`
for `+60`, without another null/type/health check. Notification callbacks
remain an unresolved mutation boundary. — SAV-959

A `+60` refusal does not necessarily suppress the manager's later accounting.
After it returns, the manager rereads the current source and owner, updates
that Player's selected counter, and calls the Player-owned Diary with the
victim. Missing victim owner or the paired multiplayer/owner gate can refuse
progress while that later Diary operation still runs on the admitted prefix.
Arithmetic and notification cuts limit the joined evidence to Medium;
native awards and first post-LOAD recipients remain Unknown. — SAV-960

Before the current actor's tick, the manager clears nonzero actor `+40` when
the referenced source's `+14` is null. It dereferences that source to make
the check; this is not an allocation-lifetime test. — SAV-962

Before the callback pair, the removed victim's attached-list loop sets each reached
Effect's counter `+42=1` and source `+44=0` before its tick callback. This is
not a global walk of Effects referring to the removed actor. Nested removal,
derived Effect and notification callbacks, actual destruction, full LOAD
ordering and the `+68` second lookup remain open. — SAV-962
