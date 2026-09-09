# Cell records and Building attachment

[Reference](format.md)

## Structures on the block plane (`TERR-STRUCT-068`…`072`, `TERR-PASS-073`)

The three tile-word arms, the type-3 cell and the border above are the whole of what the **ingest**
writes. A placed structure never goes through it. The sim class is the one the image names
`Building` (`CRuntimeClass 0x5c32d0`, `0x6c` bytes, vptr `0x59c738`; `Shop` derives from it), and it
attaches **in its own constructor**, after the ingest, to a per-cell record:

```
FUN_004e1924 (map load)  004e1c3f world ctor -> FUN_00547eb0 -> FUN_00548550 + FUN_00547d40
                         004e1c66 MOV [0x005f22c8], world        <- published, no null check later
                         004e1c76 FUN_004e2462   the .alm type-4 walker
   Building ctor FUN_005042b6 -> FUN_0050445d -> FUN_0054dbc0 -> FUN_0054d790 -> FUN_005456d0
   ~Building     FUN_005045e7 -> FUN_0054dc70 -> (payload+0x0c = 0, FUN_005456d0, free the record)
```

**The cell record**, `world+0x540b4`, keyed by the `u16` cell index, `0x34`-byte payload:

```
+0x00  u8    terrain-baseline COST byte        snapshotted at record creation (0054d8b7)
+0x01  u8    terrain-baseline STATIC block     snapshotted at record creation (0054d8c2)
+0x04  ptr   ground occupant  -> dynamic bit 6
+0x08  ptr   air occupant     -> dynamic bit 7
+0x0c  ptr   the Building
+0x10  ptr   set/cleared by FUN_005477b0 / FUN_005479b0; no plane arm reads it
+0x14..+0x28 six area-effect layer slots; each non-null one does cost <<= 2. +0x20 also blocks.
             SOURCED by TERR-STRUCT-074 and TERR-STRUCT-078: 005457e0 MOV ESI,0x6 ; 005457e5 CMP [ECX],0x0 ;
             005457ea SHL byte ptr [EAX],0x2 ; 005457ed ADD ECX,4 ; 005457f0 DEC ESI ; JNZ.
             WHICH SLOT IS WHICH (TERR-CELLREC-146): slot = +0x14 + 4*FUN_004fd8c3(spellId), and
             the two tables that function dispatches on are read out of the shipped image
             by tools/areamove:
               +0x14 spell  3 Wall of Fire     +0x24 spell 12 Light
               +0x18 spell  7 Freezing Cloud   +0x28 spell 17 Darkness
               +0x1c spell  8 Poison Cloud     +0x20 spell 19 WALL OF EARTH  <- the blocking one
             So a Wall of Earth is the only area effect that writes passability, and it does
             so through this arm rather than through anything in the area module. The
             registration is FUN_0054e730 (which then calls the recompute at 0054e829 or
             0054e97b) and the removal FUN_0054e9e0 (0054ea7c clears the slot, 0054eaff
             recomputes). MAGIC-WALLBLOCK-045, MAGIC-AREACOST-046.
+0x2c  u8
```

**Who takes `+0x04` and who takes `+0x08`** (`TERR-CELLREC-146`). Not *ground* and *air*
as such: the selector is the actor's movement-domain byte, read through `vt+0x20` in both directions
— `FUN_00544ec0` at `00544efc` (`JBE` rejects 0 and below, `<= 2` takes `+0x04`, `== 3` takes
`+0x08`) and `FUN_00545230` at `0054527f` (clears `+0x08` for 3, `+0x04` for 1 and 2). Domain 2 —
`Ghost` and `Bee`, `MOVE-DOM-028` — therefore shares the slot with ordinary ground movers. Both
slots hold at most one actor and a taken slot fails the entry (`00545001`, `005450d6`). `+0x10` is
the **sack** slot: `FUN_005477b0`'s one caller `FUN_0050f715` is the sack registration of
`ITEM-SACK-010`, and no plane arm reads the slot. Accessors by returned displacement: `+0x04`
`FUN_005463d0` / `FUN_00546520` / `FUN_0054a180`; `+0x08` `FUN_00546590`; `+0x0c` `FUN_0054df10`;
`+0x10` `FUN_00547be0` / `FUN_00547c60` / `FUN_00547cd0`.

**An actor occupies every cell of its footprint** (`TERR-FOOTPRINT-147`), the same shape
the building attach uses. `FUN_00544d00(map, actor)` reads the footprint side once
(`00544d12 CALL dword ptr [EDX + 0x1c]`), runs two nested loops both bounded by it (`00544e2f`,
`00544e49`) and calls `FUN_00544ec0` once per covered cell at `00544e6e`; a refusal from any covered
cell stops further iteration (`00544e75` to `00544e99 XOR EAX,EAX`), without local
rollback of earlier cells or mover+72/+82..85. Its five callers include the
sub-cell step `FUN_00548c60` and the arrival `FUN_005495f0`, so the `n x n` record entries are
rewritten on every cell transit, and `FUN_005456d0` derives dynamic bits 6 and 7 per cell from the
slot rather than from a separate footprint walk.

This is prefix-preserving failure, not atomic entry. A conflict at each ordinal
of a2x2 footprint leaves the earlier successful actor slots in place under the
selected normal-return paths. — SAV-CELLFAIL-583

Existing domain 1/2 actor entry checks the trigger before testing occupied+04;
its caster call can precede an eventual refusal. Domain 3 checks+08 without that
trigger arm. A missing record takes creation instead:52 zeroed bytes, current
cost/static baselines in+00/+01, then refetch and actor store. The creation path
does not revisit the existing-record trigger branch. Existing record reuse keeps
its baseline, tail and residue. — SAV-CELLENTRY-582

Sack registration is a separate programme. It reads Position+02 through
Sack+10, refuses Dynamic bit 0, and rejects every nonzero existing+10 slot.
Writing an empty existing+10 returns without recompute. Missing-record creation
zeroes 52 bytes, captures current Cost/Static before setting Static bit 5, stores
the Sack, then recomputes. — SAV-SACKENTRY-590

Sack removal clears a present record's+10 without testing zero or identity,
then recomputes. Missing records return 0. Its deletion predicate tests the four
occupant slots, byte+02 and byte+2c, not the six layer pointers or other residue.
— SAV-SACKREMOVE-591

The deletion arm restores Cost and Static from payload+00/+01 and preserves
current Static bit 4. Dynamic is not copied from the restored baseline and its
record-present bit 5 is not cleared: it retains the preceding recompute result,
plus the conditional bit 4 OR. The recompute itself never reads Sack+10; its
other payload inputs remain active. Allocation-release effects are outside this
direct write-set. — SAV-SACKPLANES-592

Sack lookup separately requires Static bit 5 before hash lookup. The creation
caller uses that lookup; registration's caller dispatches append or deletion,
and two selected removal callers continue without testing removal's result.
Those dispatches do not establish complete caller side effects.
— SAV-SACKCALLER-593

Detach tests the selected actor slot only for nonzero, not equality to the
actor argument. It clears and recomputes before testing whether four occupant
slots, layer-count+02 and operation+2c are all zero. Only that predicate enters
record deletion; other tail/residue bytes do not retain it. Successful detach
copies current Position cell/fractions into mover+86..89; missing-node/zero-slot
refusal returns before those stores. — SAV-CELLLEAVE-584

Entry+82/+83 cache the low bytes from cell-X/Y accessors before
`0054a620`; +84/+85 cache full-X/Y low bytes after it (Position+04/+05
under `SAV-TOKENPOS-074`). Word+72 receives the callback's AX. This order
does not establish an atomic snapshot or a fixed runtime+72 value.
After recompute and optional removal, detach copies Position+00/+01/+04/+05
to+86/+87/+88/+89, reloading actor+10 for each byte.
— SAV-CELLFAIL-583, SAV-CELLLEAVE-584

**Bit 4** (`TERR-PASS-148`). Written by `FUN_0054e070`, four instructions that OR `0x10`
into both planes for one cell, called from the area module's per-cell add and blast on a fourth
argument meaning *the inner effect does damage*. `wall_of_earth` never reaches it. No mover mask
contains bit 4, and its one consuming reader is `FUN_0054e140`, which run-length encodes it over the
map interior into a `CArchive`. It is transmitted state, not passability.

The recompute reads the record through a **52-byte snapshot**: `005456ed` calls the map lookup and
`00545704 REP MOVSD` copies 13 dwords from `record+0x0c` into the scratch at `world+0x5402c`, and
every `payload+N` below is really `scratch+N`. A lookup miss returns at `005456f4` having written
no plane byte at all.

**The recompute** `FUN_005456d0(world, cellIndex)` — the only routine that turns a record into plane
bytes, and a no-op on a cell that has no record:

```
static  := payload+0x01 ; cost := payload+0x00        the terrain baseline, restored first
static  |= 0x20 ; dynamic := static                   bit 5 = this cell has a record
payload+0x04 ? dynamic |= 0x40                        ground occupant
payload+0x08 ? dynamic |= 0x80                        air occupant
if payload+0x0c:                                      a Building stands here
    pos = *(u8**)(obj + 0x10)                             ; a POINTER, see "the position object"
    bit = (cellRow - pos[1]) * obj[0x60] + (cellCol - pos[0])   ; SHL masks the count to 5 bits,
                                                          ; which is the only reason the engine's
                                                          ; own 32-bit intermediate is harmless
    if obj[0x64] & (1 << bit):  static |= 5 ; dynamic |= 5            ; BLOCKS
    else:                       static &= 0xfa ; dynamic &= 0xfa      ; OPENS
                                cost := costTable[5] (CostCracked, shipped 6)

  POLARITY, settled by branch displacement (`TERR-STRUCT-078`). It is decided by
  00545793 = 74 20: a JZ whose rel8 fixes the target at 00545795 + 0x20 = 005457b5, the
  block that begins b1 fa = MOV CL,0xfa. ZF is set when TEST finds the bit CLEAR, so
  CLEAR takes the AND-0xfa arm and SET falls through to OR 5. The fall-through measures
  exactly 0x20 bytes and the taken block exactly 0x28 (005457b3 = eb 28), so neither can
  be misaligned by a decode. The shipped Data.bin column title "Passability" is therefore
  the INVERSE of what the code does with it: a set bit is impassable.

  both arms load the operand once and then read-modify-write each plane separately:
    BLOCKS  0054579b MOV DL,0x5   ; 0054579d OR  BL,DL ; 0054579f MOV [EAX+0x10000],BL
                                  ; 005457ab OR  CL,DL ; 005457ad MOV [EAX+0x20000],CL
    OPENS   005457bb MOV CL,0xfa  ; 005457bd AND DL,CL ; 005457bf MOV [EAX+0x10000],DL
                                  ; 005457cb AND DL,CL ; 005457cd MOV [EAX+0x20000],DL
  There is no `AND r/m8,0xfa` in the routine: 0xfa reaches the byte only through CL.
for p in payload+0x14 .. +0x28: if p: cost <<= 2
if payload+0x20:  static |= 5 ; dynamic |= 5
if oldStatic & 0x10: static |= 0x10 ; dynamic |= 0x10  bit 4 is carried across every recompute
```

So a block byte is **derived state**, recomputed per cell from (terrain baseline, occupants,
building). A consumer must keep the baseline, not only the current byte, or it cannot demolish.

**The footprint** is a rectangle plus two 32-bit masks, all four from the class's `Data.bin`
Buildings entry (`FUN_0050445d`, params 0/1/4/5, the file's own column titles):

```
obj+0x60  u8    sizeX  the extent along the plane key's LOW byte  (param 0)
obj+0x61  u8    sizeY  the extent along its HIGH byte             (param 1)
obj+0x64  u32   "Passability"     the BLOCKING set    tested by FUN_005456d0 @00545790
obj+0x68  u32   "BuildingPresent" the ATTACH set      tested by FUN_0054dbc0 @0054dc06
```

`FUN_0054dbc0` walks `row = 0..sizeY-1`, `col = 0..sizeX-1`, bit index running continuously, and
attaches `cell = low16(((objRow+row) << 8) + objCol+col)` for each set bit of `BuildingPresent`.
This is addition, not bitwise OR: synthetic out-of-byte coordinates can carry into the other
coordinate. Neither registration body clips against the authored map dimensions
(`UNIT-STRUCTCELL-070`).
`sizeX` bounds the **inner** loop and is added to the position object's byte 0 — the same axis the
tile grid strides by 1 and the `.alm` type-4 record's `+0x00` carries (`ALM-OBJ-061`); the shipped
masks corroborate it, since `Horisontal Bridge` (6×4) and `Vertical Bridge` (4×6) are deck patterns
only under this assignment.
`FUN_0054d790` refuses a cell that already carries a building and the walk then aborts. A footprint over 32 cells aliases — the 11×4 `Castle` folds bits 32..43
onto 0..11. The `.alm` extension arm overrides `(w,h)` from the record and sets
`Passability = 0`, `BuildingPresent = 0xffffffff`; `kind == 0x21` is class id 33,
`Vertical Wooden Bridge`, so the arm exists to let a map author size a bridge, and it opens every
cell of the rectangle. **The arm is selected by `(w & 0xff) + (h & 0xff) > 0`** at
`00504509`…`0050450f`, not by the kind (`TERR-STRUCT-090`) — the caller-supplied bytes are file
`+0x14` → `obj+0x60` and `+0x18` → `obj+0x61` (`ALM-OBJ-062`), and the other two callers of
`FUN_0050445d` push literal zeros, so this arm has exactly one reachable caller. The
`Passability = 0` store happens on **both** sides of the arm's own `w·h > 32` test
(`00504559`, `00504568`, same immediate): the branch is dead (`TERR-STRUCT-077`).

Registration is not transactional. A collision aborts immediately and leaves earlier accepted
references in place; the constructor ignores the return value (`UNIT-STRUCTCELL-070`). Thus the
rectangle, mask-selected cells and successfully attached cells are different sets. Installed collisions do not establish a partial-prefix runtime case; the
conditional non-transactional rule is retained (`UNIT-AREAPOP-075`).

`~Building` independently calls `FUN_0054dc70`. This walks the current dimensions and mask again,
returns at a missing cell record, and clears an existing record's `+0xc` without checking pointer
identity. It does not replay a saved successful-attachment list. Therefore cleanup of a partially
registered object has a conditional ownership hazard; actual destructor reach and HP-to-destruction
ordering remain Unknown (`UNIT-STRUCTDETACH-074`). Ring and blast consumers read each current
cell reference anew, so a surviving alias can be visited repeatedly (`UNIT-AREAVISIT-071`).

The cell accessor does not read Building HP. Recompute reads Position, width and
the blocking mask, but no HP or class gate. Conditional on an unchanged reference,
mask and other cell inputs, HP 1, 0 and -1 therefore produce the same lookup and
blocking/opening contribution (`UNIT-STRUCTNEXT-079`). This conditional consumer
contract is not a lethal-hit lifetime rule: the bounded caller search leaves the
actual reference/mask survival boundary Unknown (`UNIT-STRUCTBOUND-081`).

**Two call sites, three callers** (`ALM-CLS-063`). `FUN_0050445d` is called only from
`FUN_005042b6` (`Building`'s ctor) and `FUN_0050433f`, but the `Shop` ctor `FUN_00505cbd` opens by
calling `FUN_005042b6` at `00505cea` — so a `kind ∈ {0x22,0x23}` placement attaches a footprint too,
always from the table. On the shipped maps that is 450 cells over 50 placements: a 3×3 square with
the bottom-left cell open as the doorway, 400 blocking cells and 373 cells that a ground mover could
otherwise walk through.

**The attach does not always recompute** (`TERR-STRUCT-076`). `FUN_0054d790` has two exits and only
the one that had to **create** the cell record ends `0054d994 CALL FUN_005456d0`. If a record was
already there with a null `payload+0x0c`, the building pointer is stored (`0054d809`) and the
routine returns 1 at `0054d85f` with no plane byte written — the cell keeps its old bytes until
something else recomputes it. At map load no record exists before the type-4 walk, so 0 of the
17 057 shipped attachments take that exit; a consumer that places a building at runtime must decide
what to do about it. The recompute's callers are eleven routines, not two:
`FUN_00544ec0` (×4), `FUN_00545230`, `FUN_005477b0`, `FUN_005479b0`, `FUN_0054d790`, `FUN_0054dc70`,
`FUN_0054e730` (×2), `FUN_0054e9e0`, `FUN_0054f680`, `FUN_00545de0`, `FUN_00545f50`.

**The position object** (`TERR-STRUCT-075`). `obj+0x10` is a pointer, allocated and stored by the
base actor constructor (`004f2523 PUSH 0xc`, `004f2525 CALL 0x00572824`, `004f2562 MOV [ECX+0x10],
EDX`) and written whole by `FUN_00544550`:

```
+0x00  u8   col           +0x01  u8   row
+0x02  u16  (row<<8)|col  the packed cell index the record map is keyed by
+0x04  u8   0x80          +0x05  u8   0x80    the half-cell sub-position
+0x08  ptr  the world
```

The chain from the file is complete **and unpermuted** (`ALM-OBJ-061`): the `.alm` type-4 record's
`+0x00` is read into the local at `[EBP-0x68]` (`00512903`) and `+0x04` into `[EBP-0x74]`
(`00512914`); those are `SAR ...,0x8`-ed at `005129ae`/`005129a7` and pushed **last** and
second-to-last, so they become the in-memory record's `+0x00` and `+0x02`; `FUN_004e2462` pushes
`byte[rec+0x00]` last into `FUN_00544550` (`004e25b8`, `004e25b1`), which makes it `[ESP+0x4]` and
therefore **byte 0**. So file `+0x00` is `col`. **The anchor is the object's own cell and the
footprint runs right and down** — nothing on that chain subtracts `(w-1)/2`, which is what excludes
the centred rival.

A Building can make terrain passable: its AND0xfa mask clears two ingest
bits in each plane. Structure registration must run after terrain ingest
before movement consumes the final planes. The masks read as pictures (`#` blocks, `o` opens):

```
38 Horisontal Bridge 6x4      52 Vertical Bridge 4x6      11 Church 4x3      27 Magic Symbol 3x3
   ######                        #oo#                        o##o               ooo
   oooooo                        #oo#                        ####               ooo
   oooooo                        #oo#                        o###               ooo
   ######                        #oo# (x6)
```

The ground search expands the full 3×3 neighborhood with no diagonal
corner rule. Bridge openings can connect regions disconnected after terrain
ingest; passability is the final state after structure attachment.
— TERR-STRUCT-074

**Save/load.** `Building::Serialize` (`FUN_005047c6`) round-trips `+0x40/+0x42/+0x44/+0x46/+0x48/
+0x60/+0x61/+0x64/+0x68` and does **not** re-attach; MFC's `CreateObject` ctor builds by the name
`"null"`, which misses the table and skips the attach. It does not need to: the world's Serialize
(`FUN_00544a60`) stores both plane sweeps (`TERR-PASS-053`) **and** the cell-record map itself
(`00544baa` storing, `00544c4c` loading) — the `u16 key + 0x34` records described by `SAV-CELLREC-017`.
The saved record keys correspond to occupied bit 5 cells in the named save
contract; broader arbitrary-state equivalence is not established.
