# Fog and visibility

[Reference](format.md)

## Visibility bits

Bits 15/14 carry explored/current visibility. They gate animated objects
and partial terrain repaint. Installed ALM words leave these bits clear;
the runtime modifies the light-parser tile plane at `[[view+0x80]+0x0c]`.
That pointer belongs to the render object layout, distinct from the main
reader's simulation input. — TERR-TILE-079, TERR-FOG-080, TERR-FOG-081,
ALM-TILEVIEW-122

The light parser clears bit 13 only, so an authored high pair `01` can
survive that operation. The known stamp and clear have these local effects:

| Input pair (15/14) | After light parse | Immediately after clear | After an admitted stamp |
|---|---|---|---|
| 00 | 00 | 00 | 11 |
| 01 | 01 | 00 | 11 |
| 10 | 10 | 10 | 11 |
| 11 | 11 | 10 | 11 |

The clear column ends before restamping in the same event. A full event's
final grid and its native interval are separate observations.
— ALM-TILEVIEW-122, TERR-TILECLEAR-169, TERR-DRAWSTAMP-170

The global `01`-unreachability clause of TERR-FOG-082 is partially retracted.
Its shroud classifier remains exact: `11` selects level0, `10` selects level8,
and both `00` and raw `01` select the default level16.

Drawable `00459f50` ORs the four corner words before testing `0xc000` and
`0x8000`. Combined `0xc000` sets drawable state 0; `0x8000` sets 1;
`0x0000` or `0x4000` leaves 2 in a one-cell footprint without the permission
bypass. A corner `0x8000` plus another `0x4000` passes as `0xc000`, although
neither word contains both bits. This is not an all-corners-current test.
Downstream redraw is separate from that stored state. — TERR-DRAWGATE-171

The stamp's drawable virtual dimensions read class `+0xd0`; they are not
the simulation actor's footprint virtual. Both named drawable vtables reach
the original decay, permission, sight and position guards. The observed
stores require an admitted LOS cell. Bounded guard controls used a fixed
synthetic LOS field; native LOS extent and event cadence remain Unknown.
— TERR-DRAWSTAMP-170

```
state  11  currently in sight     10  explored, not in sight     00  never seen
set    FUN_00462f90 = CUnit/CAirUnit vt+0x48, tail-jumping to the body at 004598c0:
       OR byte ptr [tile + idx*2 + {1,3,W*2+1,W*2-1}], 0xc0    0045997a 0045997f 0045998b 00459990
       -- ONE immediate, both bits, so a single stamped word satisfies the four-corner OR.
       SECOND writer of bit 15: 00478870 OR word ptr [ESI],DI in the save-load path
       FUN_00477c00, restoring the saved run-length record (SAV-FOG-061, TERR-FOG-145).
       The saved-value source is a register.
clear  bit 14 only, map-wide, on a period: 0040eb44 AND word ptr [ESI],0xbfff  (FUN_0040eaee)
       No bit15 clear is identified in the immediate-form writer set.
       Register-sourced clears are outside that negative. The shroud reader
       distinguishes explored0x8000 from visible0xc000.

guards, all three, before anything is stamped
  0  00462f90 CMP byte ptr [ECX+0x15a],0x3 ; JNC        decay stage >= 3 -> nothing
  1  004598e0 TEST byte ptr [[[view+0x9b4]+0x38] + [[this+0x14]+0x4]*2], 0x8
       bit 3 of a per-player u16 on the LOCAL participant's record, indexed by the
       drawable's owning Player+0x04. FUN_004757b0 writes flags[localPlayer] = 0x0a at
       session setup and its refresh loop preserves bit 3 (0047607e AND EDX,0x8); other
       writers are the single-entry message arm at004162ec and the bulk-copy
       message arm at00416067 described below. In single player: the local player only.
  2  004598ea CMP word ptr [this+0x102],0x0             sight; the ctor leaves it 0
vt+0x44 = FUN_00462f80 repeats 0..2 per tick and adds
  3  (this+0x50 & 0x1f) == 0x10 and (this+0x54 & 0x1f) == 0x10     an eighth-of-a-cell grid
  4  this+0x8/+0xc != this+0xc0/+0xc4                              it moved since last stamp

the stamped set — CMapView+0x17cc, 41x41 int32, memset to 0 by FUN_00403c8c on EVERY stamp
  seed   mask[20][20] (= view+0x24ec) = (sight >> (8-k)) + (1 << (k-1)),  k = view+0x3f38 = 7
                                        so the unit of the field is 1/128 cell
  walk   Chebyshev rings r = 1..19, four edges each; stop at the first fully blocked ring
  cell   FUN_00403b77:  mask[dx][dy] = mask[pred[dx][dy]] - (cost[dx][dy] + h(cell) - h(obs))
         visible iff > 0.  h is map+0x10, signed bytes (the .alm type-2 Altitudes grid);
         h(obs) is sampled ONCE, so the term is per-cell altitude vs the observer's, NOT a
         slope along the ray.
  tables built once in the CMapView ctor by FUN_00403718:
         pred  view+0xaa8   41x41 int8 pairs, one step toward the observer; three zones,
                            j < i>>1 -> (-1,0), j > i<<1 -> (0,-1), else (-1,-1), mirrored
         cost  view+0x3210  41x41 int16 = ftol(128 * sqrt(i^2+j^2) / max(i,j))
                            axis 128, diagonal 181 -> the revealed region is a DISC
  margin ring cells outside [7, W-7) x [7, H-7) are never evaluated and stay 0
  clip   19 cells; max shipped scanRange is 12, so it never binds on shipped data

sight   drawable+0x102 <- actor+0xa4 by the state sync (0047c2a3 / 0047c2aa)
        hero:     FUN_004f7dfc writes ftol(((mind+reaction)/25 + 4) * 256)   -> 1/256 cell
        non-hero: the streamer's slot 11 targets actor+0xa5, the HIGH BYTE   -> whole cells
                  Data.bin Units title 11 "scanRange", Humans title 9 "ScanRange"
```

Consumer consequences: `Index` is in range on 82/82 shipped classes but a decoder should still
bounds-check it; `DeadObject` is a subscript into the same class array; `FireObject` is **not** a
class reference (`REG-OBJ-047`). **And the animated arm is not dead**: it fires on exactly the
cells the local player can currently see, so a shipped map's fires and trees animate inside the
field of view and hold frame 0 outside it. A consumer that implements only the plain arm
reproduces every shipped map at load time and never afterwards.

**Shroud / fog** (`TERR-FOG-037`): after the terrain, the same cell quad is darkened in place by
`FUN_00451200` (Gouraud level ramp, `dst = LUT[level][dst]`), or by `FUN_00450cf0` (level `0x10`,
degenerate quads filled with colour 0) or `FUN_00451710` (level `8`, degenerate quads halved as
`(px>>1) & mask`). These reuse the terrain quad and the same step-table edge walk.

**What the renderer branches on, and the levels it turns the pair into**
(amended `TERR-FOG-082`, `TERR-FOG-083`, `TERR-FOG-084`, `TERR-FOG-085`). The shroud pass does not read the
tile word directly: `FUN_00404135` projects the pair onto a **per-vertex** dword grid every frame.

```
grid    CMapView+0xa0, one dword per lattice vertex, (cols+7)*(rows+11) entries
        allocated with +0xa4 by FUN_00402c88 at +0x70 = (cols+7)*(rows+11)*4 bytes
        memset to 0 at 004043d7, filled, then memcpy'd to +0xa4 at 004046f2
        FUN_00404135 is its only content writer; FUN_00402f85 frees both

classify per vertex, off [CMapView+0x80]+0xc  (= the light-parser render plane)
        004044d0 MOV DX,word ptr [ECX+EAX*2] ; 004044d4 AND EDX,0xc000
        == 0xc000  -> level 0    00404644     11  in sight        NO shroud drawn at all
        == 0x8000  -> level 8    00404674     10  explored        half brightness
        otherwise  -> level 0x10 00404695     00  never seen      black
        Raw 01 takes the same default arm as 00. The known stamp writes 11;
        the periodic clear maps 01 to 00 and 11 to 10. This does not rule
        out raw authored 01 or an additional writer.

dispatch FUN_00407b1a reads the four corner vertices (0040c1ce..0040c305) and branches
        (0040c421..0040c679) on all-equal-0 / all-equal-0x10 / all-equal-8 / otherwise.
        A cell whose corners DISAGREE is a gradient -- one fog value per cell cannot draw it.

table   [0x005e8420], built once by FUN_0044ba10: 17 rows (L = 0..16) of `stride` u16,
        stride = 65536 or 8192 by [0x005eb570]   (0044ba7f SHL ECX,4 ; ADD ECX,EAX ; SHL ECX,1)
        LUT[L][px] = each channel index scaled (v * (16-L)) >> 4, saturated at 0xff, repacked
        row 0 = identity   row 8 = (px>>1) & [0x005eb578]   row 16 = 0
        checked against the engine's own two table-free fast paths over all 65536 pixel
        values, RGB565 (mask 0x7bef) and RGB555 (0x3def): 65536/65536 on all three rows.
        [0x005eb578] = sum over channels of (0x7f >> (8-bits)) << shift   (0044bd29..0044bd85)

clock   FUN_0040eaee clears bit 14 over the WHOLE map (W*H words) then re-stamps every
        drawable in +0x9b8. Its one caller FUN_0040dcdd gates it:
          0040dd54 CMP dword ptr [0x005eb588],0x0 -- nonzero SKIPS the clear, freezing the layer
          0040dd8d AND EAX,0x1f ; JNZ             -- CMapView+0xa70, the ANIM-CLOCK-001 counter
        => once every 32 PRESENTATION ticks ~ 2 s at the default speed index, ~4 s at the
        slowest. Neither server+0x00 nor server+0x04 appears on this path: the fog is not
        simulation state, is not hashed, and stops when rendering stops.

gates   1. the shroud pixels above
        2. FUN_004597f0: OR the four corner words, & 0xc000 != 0xc000 -> drawable+0x78 = 1,
           the sprite pass's guard (00459814..00459839). Aggregate 0xc000 stores 0.
           One 0xc000 corner suffices but is not necessary: separate 0x8000 and 0x4000
           corners also pass. +0x7c latches +0x10c on change.
        3. AI visibility uses a separate array (below); passability reads the block
           planes (MOVE-DOM-027). The direct immediate-reference negative does not
           cover byte-wide tests of the high half.

persist Bit15 is saved; bit14 is recomputed. The save tail's &YA1 registry
        stores Fog.FirstState (int32) and Fog.Data (int32[]) as run lengths
        over W*H cells in idx=col+row*W order. Runs sum to the cell count.
        Store: FUN_00478c40; load: FUN_00477c00, OR-ing the decoded word.
        The ALM plane supplies the initial state; the saved runs restore
        exploration. The 32-tick clear and stamp rebuild bit14.
        The no-persistence clause of TERR-FOG-087 is retracted; its earlier
        compressed-body size argument did not describe this tail record.
        — SAV-FOG-061, TERR-FOG-145

edge    the map border is black because it is NEVER SEEN -- level 16, the same mechanism at
        its maximum, and TERR-FOG-080's 7-cell stamp margin is why it is never lit. There is
        no separate edge treatment.

reveal  permission to stamp is bit 3 of [[mapView+0x9b4]+0x38][player], and that array has
        THREE writers: FUN_004757b0's setup store, dispatcher opcode 33 (0x21, arm 004160c0)
        writing ONE entry at 004162ec, and dispatcher opcode 45 (0x2d, arm 00415c85) which
        resizes the array and memcpy's it WHOLESALE out of the message body --
        00416067 memcpy([[view+0x9b4]+0x38], msg+0xe, [msg+0xa]*2). Opcode = 3 + (slot -
        0x004186e3)/4 from the dispatcher's own SUB/JMP. A wholesale copy carries no
        displacement in the bulk copy instruction.
```

TERR-FOG-086 is partially retracted only for its individual-corner equivalence.
The aggregate gate, field stores and latch described above remain supported;
the admitted raw domain is separate from the known stamp/clear-produced states.

**The AI's vision is a second implementation of this same algorithm, not this one**
(`TERR-FOG-088`). It lives in an object embedded at `world+0x58ee8` — so `AI-SIGHT-006`'s byte
map at `+0x2a008` **is** `world+0x82ef0` — with the same recurrence and its own `pred`
(`+0x22000`), `cost` (`+0x28000`) and accumulator (`+0x24000`, whose centre cell `[20·64+20]` is
the address that row calls the field `fog+0x25450`). Its `k` is **`[Scanning] ScanShift` from
`World\Data\map.reg`** (ships 7), while the view's `k` is the compiled constant at `00402b61`;
they agree only because the shipped value equals the code default, so **editing `ScanShift` moves
the AI's sight and leaves the player's fog untouched**. Its byte map is cleared by
`FUN_005474f0`, whose two callers stamp **one actor** (`FUN_0052dea0`) or **a whole group**
(`FUN_005365e0`). Different storage, clock, consumer and lifetime — do not merge them.

Both of the server's tables are built by **`FUN_00546790`**, called once from the object's init
`FUN_00547510` in each of the four world constructors, before the map load; there is no other
writer and no other reader than `FUN_00546c70` (`TERR-SIGHT-115`, `AI-LOS-087`…`AI-LOS-091` —
the layout and the region are specified in [AI](../ai/format.md), which is where the stamp is
consumed). The same init builds a **fourth** grid at `+0x20000` holding `di² + dj²`, which serves
a disc test in `FUN_00546fa0`/`FUN_005470a0` and has nothing to do with line of sight.
The view and server use the same predecessor/cost algorithm, including
four final stores repairing cells(+1,0) and(-1,0). At k7, flat-ground
scanRange1..6 yields 9,21,45,69,105,145 cells; range 19 yields 1253. The
127-cell radius 6 result in `TERR-FOG-080` is withdrawn. The AI seeds from
the high byte at actor+0xa5 and truncates fractional sight, while drawn fog
uses the full u16. Their playable rectangles also differ by one cell on
each edge. — TERR-FOG-117, AI-SIGHT-093, AI-SIGHT-094, TERR-FOG-118

**The playable rectangle** (`TERR-SIGHT-116`) is four bytes at `world+0x58ee0..0x58ee3` written by
the map load `FUN_00548550` as `(8, 8, W-9, H-9)`, with the same corners packed as words at
`+0x58ee4` = `0x808` and `+0x58ee6`. Every sight ring cell is tested against them; the compares
are **byte-wide** while the height and visibility indices are exact 32-bit, so on a 256-wide map a
true column of `-9..-11` wraps to `245..247`, passes, and reads the previous row.
