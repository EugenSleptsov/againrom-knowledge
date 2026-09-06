# TRIGGER — the mission-script runtime — specification

Level 3. Promoted, evidence-backed claims only. Ledger: `claims/trigger.md`.

**Status: complete for the shipped campaign's reachable operation vocabulary (`TRIG-CLOSURE-037`,
`TRIG-INSTCENSUS-046`).** The bounded census walks
all 28 archived type-7 maps on both roots with exact tiling and records 52 reachable operation
rows per root. Its completeness clause includes instant 13 (`TRIG-TAKEITEM-038`…`TRIG-ITEMTEST-040`)
and instants 16, 17, 18, 32 and 33 (`TRIG-OFFMAP-041`…`TRIG-CELLEFFECT-045`), each of which had only the vocabulary-table
claim `TRIG-ACT-004` indexed against it and therefore no statement of what the arm does. This does not claim that every dormant arm
is specified: operations absent from the campaign, dead dispatch arms, the catalogue-only group
subcommand 18, and authored actions no trigger reaches stay outside that closure. The full
separation is maintained by `TRIG-CLOSURE-037` and `TRIG-INSTCENSUS-046`.

This is not a file format. The bytes a map stores are `formats/alm` and `claims/alm.md`; the
clock and the session object are `formats/session`; the save container is `formats/sav`.

## The one thing a consumer must not get wrong

**There is no evaluator for victory.** Winning and losing are two integers that an ordinary
*action* increments. A mission with no action that increments them can never be won and can only
be lost by the hero dying. Implementing "the engine checks whether the objective is met" produces
a game that never ends.

## The machine

The map's authored script is **compiled once at load** into three flat stores on the session
object, and never read from the map again:

```
session+0xc2ac   check list      each element an 80-byte record
session+0xc2b0   instant array   each element a 72-byte record
session+0xc2b4   pattern list    each element a 24-byte record
session+0xbd34   int  slot[100]  the result/variable register file
session+0xbec4   byte latch[1000]  one per pattern, indexed by the trigger's array position
```

Both array sizes are hard: nothing bounds-checks them at build time. A map with more than 100
conditions or more than 1000 triggers corrupts memory. The shipped maxima are 64 and 36.

Slot 93 (`session+0xbea8`) is written by a literal, non-indexed displacement at three engine
sites, and it is not the only slot written this way: slots 90, 91 and 92 are too, but only once
each, at session construction, never again at runtime (`SAV-658`). Slot 93 alone is also written
this way at two runtime sites, both immediately and unconditionally followed by a call to the
check evaluator and then the pattern evaluator (below). Every slot other than 90–93 is written
only generically, at a script-supplied index, by the check and pattern evaluators executing a
map's own authored data. A check's own compiled-order subscript stays within the shipped maxima
above, so no check in the analyzed corpus reaches slot 93; whether a pattern's own index
assignment is bounded the same way was not established. — SAV-650, SAV-651 (partially
retracted), SAV-658

## The pass

Once per **full tick** — before any actor is walked, ≈ 992 ms at the shipped speed default
(`TRIG-EVAL-001`) — the engine reads, increments and stores slot 93 (`session+0xbea8`). It then
calls the check evaluator, then the pattern evaluator, with no other instruction between any of
the three steps. Two engine sites run exactly this triple. `FUN_005336a0` runs it, reached
unconditionally from the dedicated-server loop `FUN_004d214b` and again from a
`server+0x04 % 16 == 6` gate inside `FUN_004d2551` — the same function that advances the
full-tick counter, on its own separate `% 16 == 15` gate (`SESS-TICK-004`). A second, minimal
routine, `FUN_005394e0`, runs the identical triple as its entire body; no caller was located for
it in `.text` or as a whole-image dword reference, so whether or how it runs is Unknown. — SAV-650

1. **Every check** is evaluated and writes its own slot. One 22-arm dispatch on the check's
   opcode (1..22); an opcode outside that range, and the two dead arms 11 and 13, write nothing,
   so the slot keeps its previous value.
2. **Every pattern** is evaluated. Up to three triples `(slotA, slotB, code)`; the six codes are
   `0 ==`, `1 !=`, `2 >`, `3 <`, `4 >=`, `5 <=` and anything above 5 is false. The triples are
   **ANDed with short-circuit**; there is no OR.
3. A pattern that passes sets its **latch** and runs its up-to-four instants, in slot order. One
   34-arm dispatch on the instant's opcode (1..34); arm 9 is dead and out-of-range logs
   `"Script: Bad instant %d"`.

## Firing discipline

`trigger+0xb4` in the map file is the once flag.

- **1** — fire at most once for the whole session. The latch gates re-entry.
- **0** — the latch is cleared at the head of every pass, so the pattern fires **again on every
  full tick** its conditions hold.

Shipped corpus: 387 one-shot, 34 repeating.

## What must be in a save

`slot[100]`, `latch[1000]` and the three outcome integers are in the **world half** of a save
(`SAV-SESS-031`): present in mid-mission saves and absent between missions. **A mid-mission load
restores this session block first (`004d13b4`) and rebuilds the map trigger programme later
(`004d143f`).** The builder can preset at least an authored-variable result slot, so builder-first
followed by a blanket saved-slot overwrite is not equivalent. A consumer that saves mid-mission
without the latch array re-fires every one-shot trigger on reload. — TRIG-SAVE-008 (load
order and corpus count amended), SAV-RECON-268

The selected positive overwrite is the authored opcode `0x10002` check arm: it copies node+0x48
into the current compiled register slot and advances that slot. Successfully compiled earlier
normal checks and constants determine its index; a resolution-rejected normal check does not
advance it. Later `0053b870` repairs `session+0xa9c0` only. Its adjacency to `0053b880` does not prove
a LOAD call to that separate register writer. — SAV-652 (partially retracted), SAV-710,
SAV-711

## Mission end

```
instant 4  ->  session+0xb3ac++     win
instant 5  ->  session+0xb3b4++     lose
check   18 ->  session+0xb3b4++     lose, when its unit is dead; writes no slot
```

At full-tick reporting, script loss `== 1` precedes win `== 1`, but only after the active-human,
latch and primary-fall gates. Latch `player+0x3c >= 2` cannot reach either counter arm; campaign
mode returns without automatic latch reset. Loss-to-win is therefore not enabled merely by
incrementing the loss counter past 1. With a living primary and latch below 2, counter outcomes
send `0xb4`/`0xb5` and write latch 2/1 without a mode predicate (`TRIG-END-009`, amended).

An earlier branch tests the primary pointer `player+0x34` and its stage, not every companion.
It writes latch 2 and sets owned actors' HP to -50. Its failure packet is restricted to joined
participants with actual `server+0x0c == 0`; nonzero-mode recovery has separate entry-active and
HP-below--53 gates. Campaign panels also pause frontend ticking (`MISSION-DEFEAT-045`,
`SESS-DEFEAT-065`).

## How a node's ten authored slots reach an arm

Every arm below is written in terms of `p0`, `p1`, … and of six reference fields. **Those are
not the file's `Par0..Par9`.** The builder loops all ten `(value, type)` slots and dispatches
`type − 2`, bounded by 7, through an 8-entry table (`TRIG-PARAM-030`):

```
type 2  -> group      first -> rec+0x34, a second one -> rec+0x3c
type 3  -> player     first -> rec+0x38, a second one -> rec+0x3c
type 4  -> unit       first -> rec+0x30, a second one -> rec+0x3c
type 8  -> item       rec+0x40, a WORD, and the stored value is value + 0xe18
type 9  -> building   rec+0x44, resolved by NAME
everything else       the next plain-int slot, rec+0x08 + 4*k
```

`k` is a **separate** counter that advances only on the last line, so a reference-typed slot
takes no `p` index while an *unused* slot — the editor's `<None>`, type 0 — does:

```
p_index(i) = i - (number of reference-typed slots before i)
```

The file's own layout is the opposite (`ALM-TRIG-045`: parameters are stored by slot, not
packed), and reading an arm against the file rather than against this rule silently changes
what it does. The standing example is instant 23, whose every shipped node authors `Player`
in `Par0` and `Money` in `Par1`: read by file slot, the arm adds a player index to a purse
and never reads the money at all.

The item slot is an **offset, not a code**. `0xe18` decodes under `ITEM-CODE-029` as class
14, index 24, so an authored `Item` value `V` names class-14 index `24 + V`. The shipped
labels agree: `V = 11` is `"Add item35 …"`, `V = 13` is `"Add item Tooth(37) …"`. `V` above
231 carries into the class nibble and the factory returns null.

**No arm null-tests a reference field**, and none needs to: a node whose reference fails to
resolve is dropped whole at build time and logged (`TRIG-BIND-010`).

## The vocabularies

The promoted core state and outcome operations below are `TRIG-ACT-004` (instants
3, 8, 4 and 5) and `TRIG-COND-003` (check 19 and build-time opcode `0x10002`):

```
instant 3   slot[p0] = p1        set a mission variable
instant 8   slot[p0]++           increment it
instant 4   win
instant 5   lose
check   19  slot = slot[p0]      read a variable back into a comparable slot
check   0x10002  (build-time)    declare a variable and preset it to value[0]
```

A condition's constant operand is itself a check node with opcode `0x10002`, given a slot and
preset once at build. There is no immediate operand anywhere: **a comparison is always slot
against slot.**

Four more arms, read at instruction level and each with a trap in it
(`TRIG-SACK-022`…`TRIG-GIVEALL-025`):

```
check   14   slot = (a sack exists at cell (p0,p1)) ? 1 : 0
             p0 and p1 are read as BYTES; the sack itself is never yielded
instant 2    broadcast a "show message p0" packet to every client
             writes no slot, no counter, no latch and no actor field
instant 20   move the unit's whole item container to a sack at the unit's OWN cell,
             merging into a sack already there, and give the unit a fresh empty one
instant 28   move the whole container from the first Target_Unit to the second,
             destroy the source, re-seat the giver, notify each owner separately
```

The container is `actor+0x7c` (`ITEM-CONT-004`) and the routine that moves it **deletes**
it, which is why both instants must re-seat the giver. A consumer maintaining a cell→sack
index for pick-ups already has everything check 14 needs; one that does not will implement
the check against the wrong structure.

Three more (`TRIG-ADDITEM-027` (amended)…`TRIG-NEAREST-029`):

```
instant 12   obj = itemFactory(rec+0x40); if obj, add it to the Target_Unit's container
             a CREATION, not a transfer -- nothing is taken from anywhere.
             The owner notification runs even when the code decodes to null.
instant 23   Target_Player money += p0, then a packet carrying the NEW balance
             addressed to that player alone. An accumulate, no clamp, no floor.
check   15   slot = min over every unit the Target_Player owns of
             chebyshev(unit cell, (p0,p1)); 0xff when the player owns none
             p0 and p1 are read as BYTES. There is no unit parameter.
```

Check 15 is the vocabulary's only *"how close is anybody"* test and is the second
most-authored check in the corpus; the trap is its empty answer, because `0xff` under the
usual `<= constant` pattern is the opposite verdict from the `0` a naive implementation
returns.

### The three item arms (`TRIG-TAKEITEM-038`, `TRIG-XFERITEM-039`, `TRIG-ITEMTEST-040`)

Instants 11, 12 and 13 all end by resyncing a unit's inventory to its owner, and they are
easy to confuse. They differ in what they do to the container and in how many packets they
emit.

```
instant 11   take ONE unit of item (u16)rec+0x40 out of the FIRST Target_Unit's
             container; notify that unit's owner UNCONDITIONALLY; then, only if
             something was taken, add it to the SECOND Target_Unit's container and
             notify that owner too. One packet on a miss, two on a hit.
             No shipped map authors this opcode.
instant 12   obj = itemFactory((u16)rec+0x40); if obj, add it to the Target_Unit's
             container. A CREATION -- nothing is taken from anywhere. One packet,
             emitted even when the code decodes to null.
instant 13   take ONE unit of item (u16)rec+0x40 out of the Target_Unit's container
             and FREE it through the object's own scalar deleting destructor
             (vtable slot +0x04, argument 1). One packet, emitted even on a miss.
check   12   1 iff the Target_Unit's container holds item (u16)rec+0x40, else 0
check   17   the same 64 bytes as check 12, emitted twice. Only 17 is authored.
```

The take is `FUN_0050eb82` and arms 11 and 13 are its only callers in the image. It
**detaches**: a stack of `n > 1` is split, one unit off through the item's `vt+0x40`, and a
count of `0` or `1` is unlinked from the container list outright. Either way the container's
running load `+0x20` drops by `(i16)+0x4a * (u16)+0x42`. So instant 13 removes exactly one
unit of a stack, not the stack, and a consumer must unlink before freeing. Checks 12 and 17
run the same lookup without detaching, so either is the non-destructive predicate for
"would instant 13 destroy something here".

The item field is a **word** at `rec+0x40`, written by the builder as `0xe18 + V` from the
authored `Item` value `V`, and it is matched against `(u16)item+0x40` on the object. Instant
11's destination comes from `rec+0x3c`, which the builder fills with the *second* reference
of whichever kind reached its own second slot -- unit, group and player each have their own
first/second flag and all three store to that one field. So `rec+0x3c` holds a unit only when
the node declares two `Target_Unit` parameters, which is what the editor's own declaration
for the opcode does.

### The campaign-reachable closure arms (`TRIG-CLOSURE-037`)

The last campaign-reachable helpers are exact enough to implement directly. Parameters below
are the builder-packed `p` sequence, not the authored `Par` positions.

```
check 9:
    if subject.order.pending != 5: slot = 0
    else:                          slot = subject.order.target.mapUnitID

instant 21: castFromCell(p0, p1, p2, p3, spell=p4, power=(p5 != 0 ? p5 : 99))
instant 24: castFromCellAtUnit(p0, p1, Target_Unit, spell=p2, power=p3)

instant 29:
    key = (u16)(((u16)p1 << 8) + (u16)p0)     // a 16-bit ADD, not an OR
    for effect in sixEffectSlotsOfCellRecord(key):
        if effect != null && (u32)effect.spellID == p2: effect.life = u16(p3)

instant 30:
    for effect in Target_Unit.attachedEffects:
        if effect.spellID == u8(p0): effect.duration = u16(p1)

instant 34:
    if p0 == 6:  Target_Unit.health     = u16(p1)
    if p0 == 15: Target_Unit.defence    = u16(p1)
    if p0 == 16: Target_Unit.absorption = u16(p1)
    notify(Target_Unit)                 // also for every unsupported selector
```

The temporary-caster builders place `power` in `skill[Sphere]`. The accessor
`FUN_004fe11e` reads byte 0 of Spell-row parameter 2, whose title is `Sphere`, and the caller
masks that result to a byte before indexing the six-word skill vector. Shipped spells use only
Sphere 1..5; a custom Sphere-0 row writes General.

Instants 21 and 24 allocate a temporary actor at the source cell. A failed cast attempt retries
on a later actor tick; a successful one counts down, applies the spell, and removes the actor from
`session+0x2c`. The original save format serializes the resulting effects but not that in-flight
actor. A deterministic consumer therefore includes it in live simulation state without inventing
save persistence (`TRIG-CAST-033`).

These helpers are not script-exclusive (`TRIG-CASTACTOR-044`, exclusivity clause
retracted). A map-authored type-9 cell record also reaches the unit-target helper
through actor footprint attachment (`UNIT-M10CELL-054`, `UNIT-M10ENTRY-055`).
The two direct constructor callers do not constrain their own upstream callers.

Instant 25 writes a persistent actor-entry trigger:

```
cell = (u8(p3) << 8) | u8(p2)
tail = {u8(p0), u8(p1), 0, u8(p3), 0, u8(p3)}
writeOrCreateCellRecord(cell, tail)
```

Actor attachment accepts `tail[0] != 0 && tail[0] != 26`. It uses `tail[0]` as the spell id,
`tail[1]` as power and `tail[2:4]` as the temporary caster's source x/y; the entering actor or its
current cell is the target, selected through the spell table. A separate arrival reader accepts
`tail[0] == 26` and treats `tail[4:6]` as relocation x/y. The two shipped nodes write 3 and 9, so
they qualify for actor-entry casting rather than relocation. No controlled original load witnessed
entry into either restored cell or completion of its effect (`TRIG-CELLTAIL-035`,
`SAV-CELLLOAD-111`).

Check 4's shipped selector 6 reads signed current health; check 21 reads signed current building
health. Their field identities are supplied by `HERO-HEALTH-032`, `ALM-CLS-053` and
`SAV-BLDG-037`. This is claim composition, not a second decode of either arm.

### Map presence — instants 16, 17, 18, 32, 33 (`TRIG-OFFMAP-041`…`TRIG-MAPGROUP-043`)

Being on the map is one bit on the actor and one membership in one list. Two routines own
both, and the script reaches them five ways.

```
FUN_004f47e6(actor):                        // instant 16, and instant 18's first half
    if actor.flags4c & 0x08: return         // idempotent
    map.clearFootprint(actor)               // n x n, n from actor vtable+0x1c
    actor.flags4c |= 0x08
    onMapList.remove(actor)                 // [0x00609558]+4 = world+0xa4554
    packet(0x74, actor.id, 1, addressee = 0)          // broadcast
    packet(0x74, actor.id, 0, addressee = actor.owner)

FUN_004f4905(actor):                        // instant 17
    x, y = actor.pos.x, actor.pos.y         // RETAINED from before the removal
    if !place(actor, x, y, r = 0) && !place(actor, x, y, r = 3):
        log "Unit can't return to map - no free place"
        return 0                            // still off the map, nothing else changed
    onMapList.append(actor)
    actor.flags4c &= ~0x08
    notifyEachClient(actor); sendActorDescription(actor, mask 0x20)
    packet(0x74, actor.id, 0, addressee = 0)
    return 1

place(actor, x, y, r):                      // FUN_004f4604
    actor.pos.set(x, y)
    repeat (r*r)/2 + 2 times:               // 2 at r = 0, 6 at r = 3; the counter is
        cx = x - r/2 + rand(r)              // raised after a FAILED attempt and then
        cy = y - r/2 + rand(r)              // compared JLE against (r*r)/2 + 1
        if map.tryOccupy(actor, cx, cy): commit; return 1
                                            // at r = 0 both attempts are exactly (x, y)
    for cy in [y - r/2 .. y + r/2]:         // skipped entirely when r == 0
        for cx in [x - r/2 .. x + r/2]:
            if map.tryOccupy(actor, cx, cy): commit; return 1
    return 0
```

`rand(n)` is `MISSION-DROP-002`'s inclusive `(rand() * (n+1)) / 32768`.

```
instant 32:  for member in Target_Group: FUN_004f47e6(member)
instant 33:  for member in Target_Group: FUN_004f4905(member)
instant 18:  FUN_004f47e6(first Target_Unit)
             FUN_004f4865(second Target_Unit, first.pos.x, first.pos.y, r = 3)
```

`FUN_004f4865` is `FUN_004f4905`'s tail with the cell supplied by the caller and one
attempt at `r = 3`; its failure string is `"Unit can't enter map - no free place"`.

Four consequences a consumer has to build for:

- **The return cell is never authored.** No node of opcode 17 or 18 carries a coordinate.
  Instant 17 returns the unit to the cell it stood on, and instant 18 sends the second unit
  to the cell the first stood on. Storing the cell at removal time is therefore mandatory:
  the engine keeps it because it never clears the position object.
- **The return can fail, and the failure is silent in the game.** The unit stays off the
  map, the trigger has already latched, and nothing retries.
- **Nothing else about the unit changes.** Group membership, the item container `+0x7c`, the
  owner `+0x14`, the position object `+0x10` and health are untouched, in both directions.
  A group whose members are all off the map still answers its full count to check 1, which
  is an unfiltered read of `group+0x0c` that only death lowers.
- **Instant 18 needs two `Target_Unit` parameters**, because `rec+0x3c` is the builder's
  second reference of any kind. All six shipped nodes declare two.

The shipped campaign uses the pair as reinforcement and ambush: on `131.alm` two triggers
whose condition is `constant 2 == constant 2` fire on the first pass and take five units off
the map, and a later trigger returns four of them when a group is wiped out and a mission
variable does not hold 1.

### The cell record instant 29 reaches

Instant 29's lookup is the same 52-byte dynamic cell record `TRIG-CELLTAIL-035`'s writer
creates and `SAV-CELLREC-017` serializes.

```
FUN_0054ec40(map, key):                     // key is 16 bits: (y << 8) + x
    if !(map.cellFlags[key] & 0x20): return null      // flag array at map+0x10000
    bucket = (key >> 4) % map.bucketCount             // map+0x540bc
    for node in map.buckets[bucket]:                  // map+0x540b8
        if node.key == key:                           // word at node+0x08
            copy 52 bytes from node+0x0c to map+0x5402c
            return map+0x5402c + 0x14
    return null
```

So dwords 5..10 of the record are six area-effect pointers, and the copy is of pointers:
instant 29's write reaches the live effect. The flag array is 65536 entries and the key is
16 bits with `x` in the low byte, so **the addressable cell space is 256 × 256**. The
largest shipped campaign map is exactly 256 × 256.

**Instant 29 sends no packet.** Neither does instant 21: both mutate an object outside the
session's register file and neither calls a notification helper, so the change is invisible
to a client until the effect's own tick, or the cast actor's own ticks, change what is
drawn. The five map-presence arms all notify.

### Action 6: the three member helpers, and the gate on subcommand 10

All three arms first stop every member and set `grpAI+0x20 = 0`, then walk the group again.

**Subcommand 10 scores each member against the named unit before it issues anything.** The
second walk calls the target-cost routine `FUN_00535d30(member, namedUnit)` and compares the
result with `0xffffff`, the value returned only for a **0** cell of the 4x4 preference matrix
at `AImanager+0xb94` (`AI-PREF-070`, `AI-COST-071`). The matrix's only zeros are
`M[1][3]` and `M[2][3]`, so the veto needs a candidate of movement domain 3.

```
for member in Target_Group:
    if cost(member, Target_Unit) == 0xffffff:      # preference matrix cell is 0
        member.state    = 0xc                      # acquire with no leash
        member.ord[0x00] = member.cell
        member.ord[0x08] = 0
    elif member == Target_Unit:
        member.state    = 0xc                      # same disposition
        member.ord[0x00] = member.cell
        member.ord[0x08] = 0
    else:
        member.state    = 3                        # engage the named unit
        member.ord[0x0c] = Target_Unit
        member.ord[0x14] = (u8)member.reach
        member.ord[0x08] = 0
```

A vetoed member therefore **does not attack the unit the script named**. The same gate is in
the player's attack order `0x19` (`AI-CMD-054`). No shipped node reaches it: every authored
`Victim` is a Humans actor and therefore domain 1 (`TRIG-GRPARM-047`).

**Subcommands 11 and 15 are one shape with two constants** (`AI-DEFEND-111`,
`AI-FOLLOWSET-116`, `AI-FOLLOWRANGE-115`):

```
for member in Target_Group:
    if member == Target_Unit:
        member.state    = 0xc ; member.ord[0x00] = member.cell ; member.ord[0x08] = 0
    else:
        member.state    = 8 for subcommand 11, 0x11 for subcommand 15
        member.ord[0x10] = Target_Unit
        member.ord[0x08] = 0
        member.ord[0x70] = (u8)p1 if (u8)p1 != 0 else 3
```

`p1` is read as a **byte** out of the node's 32-bit int, in the arm and again in the helper, so
an authored 256 behaves as 0 and coerces to 3. The player's defend order passes a literal 0 and
so is always range 3, which makes the script surface the wider of the two.

Authored-but-unreferenced nodes do not make any of the three helpers reachable.

### Instants 19 and 22: the ownership move

Both dispatch to one routine, `FUN_004d1e14`; instant 22 calls it once per member of the
`Target_Group` (`PARTY-JOIN-025`).

```
instant 19: giveActor(Target_Unit,  Target_Player)
instant 22: for member in Target_Group: giveActor(member, Target_Player)

giveActor(actor, newOwner):
    if actor.group:        actor.group.remove(actor)
    actor.owner.flat.remove(actor)          // Player+0x20
    actor.owner = newOwner                  // actor+0x14
    newOwner.flat.append(actor)             // Player+0x20
    g = newGroup()                          // 0x48 bytes
    newOwner.groups.append(g)               // Player+0x24
    g.insert(actor)
    actor.visMask &= ~oldOwner.mask         // actor+0x18, both owners cleared
    actor.visMask &= ~newOwner.mask
    broadcast(actor)
```

It writes no recruit flag, does not touch `Player+0x34`, and does not allocate a new runtime id.
The actor therefore lands in a group of its own, which is what a save of a joined actor shows.

Whether the actor survives the mission is **not** decided here. At mission end the server keeps an
actor only when `0x21 <= actor.typeID < 0x40`; the client first applies an independent keep-bit
filter. The range is **not** the output of every Humans placement. The Humans constructor streams
Data.bin slot 16 and overwrites it with a player-character value only in non-zero constructor mode:
definition-id and explicit typeID placements pass zero, while an npc placement passes its exact
`Hero` flag result. Mission 20 therefore transfers four Humans that retain `0x17,0x0a,0x0a,0x0a`
and are all removed (`PARTY-M20-030`, `PARTY-M20-031`). Mission 40's `npc25` survives because that
npc record is flagged `Hero`. The transfer node itself writes no persistence flag.

### Closure boundary

Checks are runtime-reachable merely by being built because every check is evaluated. Actions and
action-6 subcommands are runtime-reachable only when a surviving trigger names their node id.
The matrix keeps these separate from authored-but-unreferenced actions, dead table arms and
operations absent from the campaign. EN and RU expose the same operation set and reachability
classes, but their counts are not interchangeable: check 10 is authored 5/4 and check 18 10/9.
The historical three-arm remainder is already `TRIG-ADDITEM-027` (amended), `TRIG-MONEY-028`
and `TRIG-NEAREST-029`; the closure retains those claim identities rather than allocating
replacements (`TRIG-CLOSURE-037`).

The **editor's** names for these arms live in `Description Checks.ini` / `Description
Instants.ini` at the EN install root. Those are Map Editor files — `rom.exe` contains none
of their literals and the RU root ships neither them nor the editor — so a name in them is a
hypothesis about an arm, never a reading of it (`TRIG-CAT-026`). `Get sack` is the standing
example: it names the right index and misdescribes the result. The same names and kinds do
reach the runtime, because the editor stamps them into every authored node
(`node+0x9c + 64*i` name, `node+0x74 + 4*i` kind).

### The two population checks count the LIVING, and neither tests health

`1 How many units contains this group` and `8 How many units this player have` are the only two
checks that count a population, and a consumer that implements either as an unfiltered walk over
the map's placements gets a number that never falls. Both are unfiltered — and both answer
*living* anyway, because the container they measure is one death removes from.

```
check 1   slot = *(u32*)(group + 0x0c)      five instructions, 0x0053957b, no loop, no test
check 8   slot = count of every member of every group of that player, body = INC EDI
```

`group + 0x0c` is the group's own `CObList` element count: zeroed by the list constructor,
raised by one in the node allocator, lowered by one in the node freer, and returned by the
class's own `GetCount`. The `AddTail` behind `AddMember` is taken on the group itself, so the
counted list is the group's base and not the patrol list embedded at `group + 0x20`.

**The removal is per SUB-tick, not per full tick, and it comes after the script pass.** The
sub-tick body ticks every actor on the on-map list and, for one whose tick left it dead
(`actor+0x54 == 0x10`), removes it from its group **before** unlinking it from the world and
appending it to the dead list. The tick driver takes the phase from `server+0x04` before the
sub-tick increments it, so on the phase-6 iteration the whole script pass runs first and the reap
follows in the same iteration. A member that dies during the script pass is gone from the group
by the end of that sub-tick, and the next pass is sixteen sub-ticks away: **no pass ever sees a
corpse in a group.**

A consumer therefore owes two things that are one behaviour — a plain read of a live member
count, and an unlink at death in the same tick. Implementing either alone gives the wrong answer.

Two lifetime rules go with it. An emptied group is **destroyed** when its owning player's
`+0x28` is zero (a human participant) and **kept** otherwise, so a scenario group survives its
last member and reads 0 — which is what makes a `== 0` comparison reachable at all. And check 1
has **no null path**: it dereferences its group unconditionally. What keeps that safe is the
binder, which never builds a check whose `Target_Group` failed to resolve; the trigger that names
the unbuilt check then resolves to **slot 0**, which is the hazard that replaces it.

Corpus, identical in both roots: 96 check-1 nodes on 24 maps, 92 bound condition pairs, of which
**81 are true only when the count reaches zero** (`== 0` 47, `< 1` 31, `<= 0` 3).

## Identifier bands

`Target_Unit` is three id spaces, and a consumer that treats it as one will fail to resolve four
fifths of the shipped references:

```
value <  10001    a unit id from the map's own type-6 records
10001..11000      a hero ordinal, resolved against the live player list at run time
                  (and unconditionally unresolvable when the multiplayer flag is set)
value >  11000    an index into a static name table in the executable
```

## The drop table — where the player lands

The authored action `0x10002` is never dispatched. At build time its `value[0]` and `value[1]` are
truncated to bytes, packed `(y << 8) | x`, and appended to a `CWordArray` at scriptObj`+0x28` =
`mapObj+0x6c`. **Every shipped map carries exactly one, 38/38** — campaign and skirmish alike, and
on 9 of the 10 loose maps it is the map's only script node.

`FUN_004d403c` consumes it, and three details decide whether a reimplementation lands the player
where the engine does (`MISSION-DROP-002`):

```
if (mapObj+0xc != 0 && (u16)player+0x60 != 0)   x,y = low,high byte of player+0x60
else if (array.GetSize() > 0)                   i   = FUN_00504003(GetSize()-1)   <- RANDOM index
                                                x,y = low,high byte of array[i]
if (x * y == 0)                                 x   = FUN_00504003(0x46) + 0x1e   <- 30..100
                                                y   = FUN_00504003(0x46) + 0x1e   <- independently
                                                log "no drop location in .alm - random used"
```

`FUN_00504003(n)` is `(rand() * (n+1)) / 32768` — **inclusive of `n`**, and not a modulo. The
array is **not** read at `[0]`.

## Starting the mission

The map places nobody for the player. Roster slot 1 owns no type-6 record on any campaign map,
and `FUN_004d403c` iterates the **player's own unit list** at `player+0x20`, placing the hero
(`player+0x34`) at the drop cell exactly and everything else within
`ftol(max(5.0, sqrt(nUnits) + K))` of it. A consumer that builds the party from the map file
starts every mission empty (`MISSION-START-001`).

The routine's second half — a by-name lookup of a node called **`"Humans"`** whose children are
`label#id` strings with their own `(x, y, radius)`, or `param[1] == -1` meaning "at the drop
cell, radius 8" — is the `.ini` overlay's entry point below, and **no shipped map has such a
node**, so on shipped data it never runs.

## Ending it

`instant 4` is the only thing that wins, and **every campaign map carries exactly one**
(`MISSION-WIN-003`). None of the ten loose maps carries any: a skirmish map cannot be won,
only lost. A win action may be referenced by more than one trigger — four on `81.alm` — so
bind by node id.

`check 18` is the "protect this unit" objective and it is authored as a trigger with **no
action at all**: the check writes no slot, so the pattern never passes, and its only effect is
the lose it raises when its unit dies. A check is armed by being **authored**, not by being
referenced — the builder gives every check node a slot and the pass evaluates every check
(`MISSION-VIP-004`).

## The register file is shared

The build assigns check node *i* the slot *i*, constants included, and an authored **variable**
addresses the same `session[0xbd34 + p0*4]` array with its own literal number. Any variable index
below a map's check-node count is overwritten every full tick. `60.alm` ships exactly that
collision — variables 32 and 33 against 64 check nodes (`MISSION-SLOT-008`).

## Mission text

`instant 2` carries a **number**, not a string. `FUN_00473110` opens
`main.res::text/battle/m<mission>/event<NN>.txt`; the same block holds `briefing.txt`,
`briefmap.txt`, `title.txt` and `tips<NN>.txt`. The files are markup with `<NPC=n,Part=k,…>`
tags whose `n` is an `npc.reg` section number (`MISSION-TEXT-005`). **A number whose file does
not ship is a silent no-op** — the window is never built and nothing is logged; the window
itself, its lifecycle and the tag vocabulary are `formats/dialogue`
(`DLG-WIN-001`…`DLG-MARKUP-007`).

## `World\Mission\<n>.ini`

Built and opened at every map load, absent from the install, and its absence is a **silent
no-op** — the open fails, the reader returns immediately, and the one call site discards the
result. It is a sectioned text overlay whose sections the engine looks for by name:
`Humans.Hero` (`Name`, `Bag.`, `Armor`, `Weapon`, `Shield`, `Item`), `Outposts`, `Patrol`,
`StandGround`, `Items`, `Players`, `Mission`, `Monsters`. **Nothing in the trigger machinery reads
it**; its consumers are hero import, patrol paths, outposts and item setup.
