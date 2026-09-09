# SAV world, terrain and session

[Format reference](format.md) · [Document](document.md) · [Token](token.md)

These records exist only in the document's nonzero world-present arm.
The external ALM supplies the baseline map; SAV supplies overlays, state and
file-local relations. — SAV-SHAPE-023, SAV-LOAD-057

## Load order

1. Deserialize live Player/Group/actor graphs.
2. Deserialize dead actors.
3. Rebuild the global actor list.
4. Deserialize Buildings, then SpellEffects.
5. Construct the named map and ALM terrain; publish terrain.
6. Restore saved block rows, cell rows and terrain identity.
7. Restore session, then Sacks.
8. Rebuild triggers.
9. Run live-actor, then dead-actor post-load hooks.
10. Rebind object keys in every cell record.
11. Run session, then SpellEffect post-load hooks.

Saved cell rows overlay ALM-created rows before Sacks load. Cell identity
repair waits until Sacks and triggers exist. Actor/Building/SpellEffect keys
are already bound before terrain construction. — SAV-CELLLOAD-108,
SAV-CELLLOAD-109

## Block array

```text
u16 count
count * u32: (packedCell << 16) | (dynamic << 8) | static
```

| Part | Rule |
|---|---|
| packedCell | 256-stride index `(row<<8)\|col` |
| Store order | Strictly increasing, sweep window `[0x0807,0xedee)` |
| Inclusion | Only cells with `dynamic > 0x0f` |
| static byte | Terrain `+0x10000` plane |
| dynamic byte | Terrain `+0x20000` plane |

Bits 0..3 are map block bits, bits 4/5 runtime flags and bits 6/7 occupancy
bits, appearing in the dynamic plane in the observed records. This is a delta:
map-only low block bits and cells outside the sweep are omitted. Construct
ALM planes first, then apply the saved rows. A zero baseline loses terrain
blocking. — SAV-BLOCK-011, SAV-BLOCK-012, TERR-PASS-049, TERR-PASS-053

The original writer has a fixed u16 count here, unlike the following table.
No wide-count arm is specified for this array. The fixed sweep window does
not expand for a wider map. — SAV-BLOCK-011, TERR-PASS-053

## Cell table

```text
Count(n)
n * { u16 packedCell; raw payload[52]; }
u32 terrainIdentityKey
```

The hash at terrain `+0x540b4`, serializer `0054ff70`, writes two bytes from
node `+8` and 52 from node `+c` per record. The table itself has no additional
tail; the following dword belongs to the enclosing terrain serializer.
Cell keys match the block rows whose **static** byte has bit 5 in the promoted
relation. A writer must preserve that relation. — SAV-CELLREC-017,
SAV-TERRKEY-056, SAV-WRITERAUDIT-380

Payload offsets in the following table are hexadecimal.

| Payload offset | Width | Meaning / later rule |
|---:|---:|---|
| `+00` | 1 | Cost-plane baseline captured at creation; recomputation reads it |
| `+01` | 1 | Static-plane baseline captured at creation; recomputation reads it |
| `+02` | 1 | Occupied area-layer count; attach/detach recounts six slots |
| `+03` | 1 | Constructor zero; semantic consumer Unknown |
| `+04` | 4 | Movement-domain-1/2 actor key |
| `+08` | 4 | Movement-domain-3 actor key |
| `+0c` | 4 | Building key |
| `+10` | 4 | Sack key |
| `+14..+28` | 6 × 4 | SpellEffect keys indexed by map layer |
| `+2c` | 1 | Operation: nonzero except 26 arms entry casting; 26 arms relocation |
| `+2d` | 1 | Trigger power |
| `+2e..+2f` | 2 | Temporary-caster source X/Y |
| `+30..+31` | 2 | Relocation X/Y, read for operation 26 |
| `+32..+33` | 2 | Constructor zero; semantic consumer Unknown |

LOAD reads a zero-initialized 13-dword scratch, finds/allocates the node and
copies all 13 dwords. It does not clear the constructed hash: saved nodes
overwrite, saved-only nodes insert, construction-only nodes survive.
— SAV-CELLLOAD-109, SAV-CELLLOAD-110

After Sacks/triggers exist, `0054d4f0 -> 0054d6a0` repairs all ten object dwords
in every node of that union. A hit replaces the key; zero/miss retains the
saved word. Other payload bytes remain unchanged. Entry operation/power/source
feed the temporary unit/cell caster; operation 26 reads relocation coordinates.
Residue meanings and aliased/computed consumers remain Unknown.
— SAV-CELLLOAD-111, SAV-CELLLOAD-112,
SAV-CELLLOAD-113, TRIG-CELLTAIL-035

## Session block

After the cell table, terrain `00544a60` writes its own u32 identity key and
LOAD binds it to the fresh terrain. Session then transfers exactly 4,374 bytes.
— SAV-TERRKEY-056, SAV-CELLREC-032, SAV-SESS-031

| Session-wire offset | Bytes | Runtime source | Meaning |
|---:|---:|---|---|
| 0 | 400 | `+bd34` | 100 signed trigger-result slots |
| 400 | 1000 | `+bec4` | Fire-once trigger latches |
| 1400 | 48 | `+08` | Raw state |
| 1448 | 400 | `+a828` | Raw state |
| 1848 | 2508 | `+a9bc` | Self-pointer at block+4; 50×50 diplomacy matrix starts at block+8 |
| 4356 | 1 | `+a48` | Raw byte |
| 4357 | 1 | `+a49` | Raw byte |
| 4358 | 4 | `+a4c` | Raw dword |
| 4362 | 4 | `+b3ac` | Win counter |
| 4366 | 4 | `+b3b0` | Raw dword |
| 4370 | 4 | `+b3b4` | Lose counter |

The no-world branch has no session/trigger block. Outcome nevertheless survives
in Player `+3c`. — SAV-SESS-031, SAV-FLAG-027

## Session repair and trigger constants

LOAD restores every trigger slot before rebuilding the map's trigger programme.
Located literal slot writes cover 90..93: 90..92 at construction, 93 there and
at two runtime sites followed by check/pattern evaluation. Other established
slot accesses use script-selected indices; other literal slot writes remain
Unknown. — SAV-650, SAV-658

During rebuild, check opcode `0x10002` writes node `+48` into the compiled
slot and advances it; `0x10003` does neither on the selected arm. Successful
normal checks/constants advance the slot; rejected normal checks do not.
Later `0053b870` repairs session `+a9c0` (wire bytes 1852..1855) to the live
session address. It contains no register-array access and is not a call to
the adjacent register writer. Whole-LOAD execution, arbitrary indices and
unexpanded helpers remain separate boundaries. — SAV-710, SAV-711

## Reconstruction inputs

| Source | Examples | Required treatment |
|---|---|---|
| Direct SAV state | Clocks, actors, dead state, session latches, campaign collections | Restore until a named consumer replaces it |
| SAV identity | Archive graph, terrain and cell keys | Bind after target objects exist |
| SAV overlay | Block/cell rows, Fog | Construct external baseline first |
| External definitions | ALM terrain, compiled triggers, Data pointers, MapPoint objects | Select by stored map/index/relation or class-specific branch |

The located positive routes require both SAV and external state. Joint
sufficiency for arbitrary reconstruction remains Unknown: constructor outputs,
opaque direct fields and the first computed ticks are unresolved. A different
route input independent of shape is not established by a fixed-bytes crossing;
the loader can select resume/fresh construction from SAV shape itself.
— SAV-RECON-268, SAV-RECON-269, SAV-RECON-270, SAV-SUFF-302, SAV-SUFF-303

## Actor crossing and local cell refusal

Stage-zero reference repair does not center Position or reset the mover's
crossing bytes. Mover `+7c` uses hit-only identity replacement; the progress-3
arm can continue its saved step. First actual loaded scheduling remains
Unknown. — SAV-ACTORINPUT-547, SAV-HUMRESUME-460, MOVE-STEP-040

| Operation | Located local behavior |
|---|---|
| Existing domain-1/2 entry | May dispatch its trigger before refusing occupied payload `+04` |
| New cell entry | Zero 52 bytes, capture current Cost/Static baselines, then take post-creation store path |
| Existing domain-3 entry | No corresponding entry-trigger dispatch |
| Footprint refusal | Stop iteration; retain earlier cell writes and mover caches already written |
| Detach | Clear any nonzero domain slot, even a different actor; missing node/zero slot refuses |
| Delete after successful clear | Requires four occupant slots, layer count and operation byte all zero; remaining tail does not retain the node |

Reusing a cell preserves its saved baseline/tail. These local operations do
not roll back earlier work. — SAV-CELLENTRY-582, SAV-CELLFAIL-583,
SAV-CELLLEAVE-584

The following cache offsets are relative to the mover at actor `+154`;
Position is addressed by actor `+10`.

| Coordinate | Entry cache / accessor | Successful-detach cache / direct source |
|---|---|---|
| Cell X | `+82`, cached low byte of `00544a10(Position)` | `+86`, Position byte `+00` |
| Cell Y | `+83`, cached low byte of `00544a20(Position)` | `+87`, Position byte `+01` |
| Sub-cell X | `+84`, low byte of `005449e0(Position)` | `+88`, Position byte `+04` |
| Sub-cell Y | `+85`, low byte of `005449f0(Position)` | `+89`, Position byte `+05` |

Entry caches cell X/Y before `0054a620`, writes that call's AX to mover word
`+72`, then writes `+82/+83` and obtains full-coordinate low bytes for `+84/+85`.
Those reads need not be one atomic snapshot. Detach reloads Position after
recompute/optional removal and directly stores its four bytes in table order
at `0054538a/0054539c/005453ae/005453c0`.
— SAV-TOKENPOS-074, SAV-CELLFAIL-583, SAV-CELLLEAVE-584

The step ignores detach/entry refusal before its centered cleanup. It can
empty the dynamic route at actor `+178` and clear progress/claim fields even
when the destination refused entry; the static route at `+15c` survives the
named instructions. SAVE then emits current Position12, both lists, mover180
and cell payload52 without reconciling the discrepancy. The lists precede the
mover on the wire. This is a conditional local result, not a safe authored
state or an original-resave observation. — SAV-CROSSNEXT-585

## Sack cell transitions

Sack uses Position `+02` as its cell key.

| Operation | Located local behavior |
|---|---|
| Register | Refuse Dynamic bit0; refuse an existing nonzero payload `+10`, even the same pointer |
| Reuse empty slot | Write Sack pointer without recompute |
| Create record | Zero52, capture current Cost/Static, set present flag, write Sack, recompute |
| Remove missing record | Return 0 |
| Remove existing record | Clear `+10` without zero/equality check, recompute, evaluate deletion |
| Delete empty record | Require four occupant slots, count `+02` and operation `+2c` zero |

Saved baselines must survive reuse. Other residue and individual layer pointers
are not separate deletion retainers. Successful deletion restores Cost/Static
baselines while preserving current Static bit4; it does not restore Dynamic
or clear Dynamic bit5 after recomputation.
— SAV-SACKENTRY-590, SAV-SACKREMOVE-591, SAV-SACKPLANES-592

The merge/create accessor also gates on Static bit5. Two named caller slices
ignore cell-removal return before later merge/append/transfer/destructor work.
A local cell result therefore does not establish a complete item transfer or
next SAVE. — SAV-SACKCALLER-593
