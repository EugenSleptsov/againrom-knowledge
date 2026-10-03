# Claim registry — MOVE (unit movement and path selection)

Level 2 ledger. Index: [registry.md](registry.md) · spec: [`formats/move/format.md`](../formats/move/format.md). Legend: `✔ promoted` (also in the spec) · `● active` · `✖ retracted`. IDs are permanent.

Not a file format — the **simulation area** that moves a unit from A to B. It sits on top of the
planes `TERR-PASS-049…053` derived, the mover block `TERR-MOVE-054…058` identified, and the
`data/map.reg` reader `RES-CODE-020` names. Opened by
[EXP-0054](../experiments/EXP-0054-unit-movement/), which read the routines at instruction level;
the tick-order half — what fixes the walk order that *is* the contention priority — closed by
[EXP-0055](../experiments/EXP-0055-tick-order/) (`MOVE-TICK-013…017`); the fallback half — which
cell a failed search settles for, and whether anything distributes destinations across a group — by
[EXP-0068](../experiments/EXP-0068-alt-goal/) (`MOVE-ALT-018…022`, `MOVE-ORDER-023`); and the
**movement domain** itself — what the mover's one dispatch byte selects, and what each domain is
permitted to enter — by [EXP-0078](../experiments/EXP-0078-movement-domains/)
(`MOVE-DOM-024…028`); and the **rate** — what makes one unit faster than another, and against which
clock — by [EXP-0093](../experiments/EXP-0093-move-rate/) (`MOVE-RATE-029…034`), which also closes
`TERR-MOVE-058`'s Unknown and corrects a clause of `MOVE-STEP-010`. **When** that group term applies — the two gates
`MOVE-GROUP-030` located and could not identify — by [EXP-0094](../experiments/EXP-0094-group-rate/)
(`MOVE-GATE-035…037`), which also finds that a group move **does** distribute destinations and that nothing clears
the term.

Vocabulary used below is the **engine's own**, taken from its `data/map.reg` key names: a *static*
search reads the block plane without occupancy and cannot see units; a *dynamic* search reads the
plane that carries occupancy. Each produces its own route list on the actor.

| ID | Claim | Confidence | Status | Evidence |
|----|-------|-----------|--------|----------|
| MOVE-SEARCH-001 | The route search is `FUN_00541e80`, and it is a double-buffered label-correcting wave — not Dijkstra, not A\*. | High | ● active | [EXP-0054](../experiments/EXP-0054-unit-movement/) |
| MOVE-COST-002 | The step costs, read out of the image, and there is no distance estimate at all. | High | ● active | [EXP-0054](../experiments/EXP-0054-unit-movement/) |
| MOVE-TERM-003 | Termination: three exits, one of them a generation budget — and no node budget, no frontier bound, no closed-set limit. | High / Medium | ● active (amended, partially retracted) | [EXP-0054](../experiments/EXP-0054-unit-movement/), [EXP-0068](../experiments/EXP-0068-alt-goal/) |
| MOVE-ROUTE-004 | The route is extracted by walking the label field downhill from the goal, and its tie-break is asymmetric and favours the last straight step. | High | ● active | [EXP-0054](../experiments/EXP-0054-unit-movement/) |
| MOVE-PLANE-005 | Two searches over two planes, two routes on the actor — and the static search is structurally unit-blind. | High / Medium | ● active | [EXP-0054](../experiments/EXP-0054-unit-movement/) |
| MOVE-PARAM-006 | The six `[Path Finding]` scalars — the engine's own names, and the shipped file agrees with the code defaults 7/7. | High | ● active | [EXP-0054](../experiments/EXP-0054-unit-movement/) |
| MOVE-CLAIM-007 | There IS a reservation: a unit marks the cell it intends to enter, before it enters it. | High / Unknown | ● active | [EXP-0054](../experiments/EXP-0054-unit-movement/) |
| MOVE-WAIT-008 | When the next cell is occupied the unit waits, facing it — and it neither pushes, swaps, nor steps aside. At step time there is no collision test at all. | High / Medium | ● active | [EXP-0054](../experiments/EXP-0054-unit-movement/) |
| MOVE-TICK-009 | The per-tick update imposes no priority, so insertion order decides which actor gets a contested cell. | High | ● active (partially retracted) | [EXP-0054](../experiments/EXP-0054-unit-movement/), [EXP-0055](../experiments/EXP-0055-tick-order/), [EXP-0192](../experiments/EXP-0192-mission20-party-boundary/) |
| MOVE-STEP-010 | Movement is sub-cell, 1/256 of a cell per axis, and the position is a byte pair per axis. | High | ● active (amended, partially retracted) | [EXP-0054](../experiments/EXP-0054-unit-movement/) |
| MOVE-SPEED-011 | Speed does not feed the search. | High | ● active | [EXP-0054](../experiments/EXP-0054-unit-movement/) |
| MOVE-REFRESH-012 | When a route is recomputed — per cell for the dynamic one, per target change for the static one, and nothing is staggered. | High / Medium | ● active | [EXP-0054](../experiments/EXP-0054-unit-movement/) |
| MOVE-TICK-013 | The tick loop's container is a pooled doubly-linked list embedded in a 0x20-byte manager at `[0x00609558]`, its walk is head→tail, and its order is pure insertion history — the only insert operation that exists for it is AddTail. | High | ● active | [EXP-0055](../experiments/EXP-0055-tick-order/) |
| MOVE-TICK-014 | What fixes the order at map load: the creators' call sequence — campaign party first, then the map's type-6 records in ascending record order — and the list is per-*ticking-thing*, not total. | High / Medium | ● active (contested) | [EXP-0055](../experiments/EXP-0055-tick-order/) |
| MOVE-TICK-015 | How the order changes during play — and death leaves no hole. | High | ● active | [EXP-0055](../experiments/EXP-0055-tick-order/) |
| MOVE-ID-016 | The runtime id (`actor+0x04`, `SAV-ID-015`'s) is a lowest-free-bit bitmap allocation — creation-ordered only until the first death decays, and never the walk order. | High | ● active | [EXP-0055](../experiments/EXP-0055-tick-order/) |
| MOVE-TICK-017 | The walk order is NOT preserved across save/load — the saved stream never carries it, and the loader rebuilds it grouped by player. | High / Medium | ● active | [EXP-0055](../experiments/EXP-0055-tick-order/) |
| MOVE-ALT-018 | The substitution is not the caller's — `FUN_00541e80` does it itself, in three branches, and `altTarget` is a target *actor*, not a cell. | High | ● active | [EXP-0068](../experiments/EXP-0068-alt-goal/) |
| MOVE-ALT-019 | Picker A (`FUN_0054bac0`): expanding square rings around the requested cell, and the whole ring is scanned before the best in it is taken. | High | ● active | [EXP-0068](../experiments/EXP-0068-alt-goal/) |
| MOVE-ALT-020 | Picker B (`FUN_0054b420`): the contact ring around the target actor, entered where the line between the two actors crosses it and walked in both directions. | High / Medium | ● active | [EXP-0068](../experiments/EXP-0068-alt-goal/) |
| MOVE-ALT-021 | What a candidate is tested against, and what the choice is measured from — and the answer to both is the label plane, so the substitute is the cell cheapest to reach *from the mover*, not the cell nearest the click. | High | ● active | [EXP-0068](../experiments/EXP-0068-alt-goal/) |
| MOVE-ALT-022 | The substitute is never written back — it is consumed by one route extraction and forgotten. | High / Medium | ● active | [EXP-0068](../experiments/EXP-0068-alt-goal/) |
| MOVE-ORDER-023 | There is no multi-unit destination distribution anywhere: no formation, no offset table, no spread. A group move writes the same cell into every member's order block. | Medium | ● active (partially retracted) | [EXP-0068](../experiments/EXP-0068-alt-goal/), **[EXP-0094](../experiments/EXP-0094-group-rate/)** |
| MOVE-DOM-024 | The movement domain is one byte with six consumers, and passability is only one of them. | High / Medium | ● active | [EXP-0078](../experiments/EXP-0078-movement-domains/) |
| MOVE-DOM-025 | The mask is installed once, at spawn, is the mover's only passability state, and survives a save because the mover is serialized whole. | High / Medium | ● active | [EXP-0078](../experiments/EXP-0078-movement-domains/) |
| MOVE-DOM-026 | What each domain is permitted to enter, terrain and non-terrain counted apart — and the axis is not terrain, it is which of the three block bits the mask carries. | High / Medium | ● active | [EXP-0078](../experiments/EXP-0078-movement-domains/) |
| MOVE-DOM-027 | Nothing narrows a mover's verdict after the mask: no height term, no corner rule, no order-time check, no step-time check. | High / Medium | ● active | [EXP-0078](../experiments/EXP-0078-movement-domains/) |
| MOVE-DOM-028 | Which shipped classes can be non-ground — and that a player can never command one. | Medium / Unknown | ● active | [EXP-0078](../experiments/EXP-0078-movement-domains/) |
| MOVE-RATE-029 | The whole rate law, composed, with the instruction for every term — and the composition is five terms, not one. | High | ● active | [EXP-0093](../experiments/EXP-0093-move-rate/) |
| MOVE-GROUP-030 | The speed override `TERR-MOVE-058` left Unknown is a GROUP term, and it is the minimum `Speed` over the group's members. | High / Medium | ● active (amended, partially retracted) | [EXP-0093](../experiments/EXP-0093-move-rate/), **[EXP-0094](../experiments/EXP-0094-group-rate/)** |
| MOVE-TURN-031 | Stepping requires matching facing; route destruction and final default0x10 were false interpretations. | High | ● active (partially retracted) | [EXP-0093](../experiments/EXP-0093-move-rate/), [EXP-0324](../experiments/EXP-0324-turn-continuation/), [EXP-0329](../experiments/EXP-0329-mover-rate-binding/); [retraction](retracted.md) |
| MOVE-CLOCK-032 | The displacement is one step per actor per SUB-TICK, it is counted in ticks and never measured in milliseconds, and the presentation tick that drives animation is issued from the same loop iteration. | High / Medium | ● active | [EXP-0093](../experiments/EXP-0093-move-rate/) |
| MOVE-LIMIT-033 | The customisation limits of the rate (goal G2), by complete enumeration of the law's input domain rather than by sample. | High / Medium | ● active | [EXP-0093](../experiments/EXP-0093-move-rate/) |
| MOVE-DIR-034 | The eight directions, the diagonal constant, and the per-axis step — the last three things a consumer needs to compute the next position. | High | ● active | [EXP-0093](../experiments/EXP-0093-move-rate/) |
| MOVE-GATE-035 | The group rate term applies exactly when the group moves in FORMATION — one local flag decides both, and it is governed by two gates neither of which is data. | High | ● active | [EXP-0094](../experiments/EXP-0094-group-rate/) |
| MOVE-FORM-036 | A group move DOES distribute destinations: there is a formation, there is a per-member offset table, and there is a spread test — `MOVE-ORDER-023`'s headline is refuted, through the blind spot that row's own confidence cell named. | High | ● active | [EXP-0094](../experiments/EXP-0094-group-rate/) |
| MOVE-GROUP-037 | The complete writer set of `grpAI+0x44`, and the finding that NOTHING clears it — so a group keeps a rate it can no longer justify. | High / Medium | ● active | [EXP-0094](../experiments/EXP-0094-group-rate/) |
| MOVE-AREA-038 | The complete movement-side reader set for both block planes, and the one channel through which an area effect reaches it. | High / Medium | ● active | [EXP-0174](../experiments/EXP-0174-area-movement/) |
| MOVE-GATE-039 | (rom.exe) No actor is displaced except from inside one call of the order machine, so an actor the machine refuses cannot move at all. | High / Medium | ● active | [EXP-0179](../experiments/EXP-0179-effect-action/) |
| MOVE-STEP-040 | (rom.exe) The refusal begins on a cell centre, so a step already in flight completes. | High | ● active | [EXP-0179](../experiments/EXP-0179-effect-action/) |
| MOVE-EFFLIST-041 | (rom.exe) No routine in the movement path reads the drawable's effect list, so the list is not a second immobilisation mechanism. | High | ● active | [EXP-0182](../experiments/EXP-0182-effect-at-actor/) |
| MOVE-072 | Movement rule recorded under MOVE-072. | High / Medium | ● active | [EXP-0386](../experiments/EXP-0386-engagement-under-occupancy/EXP-0386.md) |
| MOVE-073 | Movement rule recorded under MOVE-073. | High / Medium / Unknown | ● active | [EXP-0386](../experiments/EXP-0386-engagement-under-occupancy/EXP-0386.md) |

### MOVE-SEARCH-001

**The route search is `FUN_00541e80`, and it is a double-buffered label-correcting wave — not Dijkstra, not A\*.** `__thiscall(world, actor, srcX, srcY, dstX, dstY, staticFlag, altTarget)`. State: a **u16 label plane at `world+0x30000`** (65 536 cells, 256-stride, cleared to `0xffff` by `REP STOSD` of `0x8000` dwords at `00541f8d`), `label[src] = 0` at `00541fcc`; and **two frontier lists** used alternately — each generation relaxes every cell of the current list and appends every cell whose label **improved** to the other list, then swaps. Two list shapes for the same two slots: parallel **byte** arrays x`world+0x50008` / y`+0x51008` and x`+0x52008` / y`+0x53008` (4096 entries each, the `0x1000` stride being the x→y offset) when the mover's footprint side is > 1, and **u16 packed-cell** arrays `+0x545b4` / `+0x565b4` when it is 1; both shapes share the u16 fill counts `+0x54008` / `+0x5400a`. There is **no priority queue, no ordering of any kind, and no closed set** — a cell is re-expanded every time its label improves, and the relaxation is the same nine-cell scan (centre included) in five places: `FUN_0054bd20` (static plane, n×n), `FUN_0054bf10` / `FUN_0054c100` (dynamic, n×n), `FUN_0054c2f0` / `FUN_0054c6f0` (dynamic, 1×1, 8 neighbours unrolled), plus two copies inlined in the driver itself — the static 1×1 pair, one half-pass per list (pushes at `00542448…005427b6` and `00542887…00542bf5`, eight neighbours each), and the static n×n half-pass built on the predicate `FUN_00543060` (push at `005421bc`/`005421d1`)

**Confidence.** High (the plane's clear, its stride, the seed, both list shapes, both counters and the improve-then-append form are named instructions, `evidence/rom-move-excerpt.md` §2–3. The live rivals are excluded by *absent* instructions, which is why they are listed: a priority queue needs an extract-min — there is none; A\* needs `h` added to the label — the only distance term computed is compared against the generation counter, never added; a "seeded at the goal" reading is excluded by the callers, which pass the unit's own `pos+0`/`pos+1` as `src` and the ordered cell as `dst`)

**Original status.** ● active

**Evidence.** [EXP-0054](../experiments/EXP-0054-unit-movement/)

### MOVE-COST-002

**The step costs, read out of the image, and there is no distance estimate at all.** For a mover whose `vt+0x20()` returns **1** (movementType 1, the ordinary ground domain): a **straight** step adds `cost[dst]`, a **diagonal** step adds `cost[dst] + (cost[dst]>>1)`. For **every other** mover: flat **2** straight and **3** diagonal, the cost plane not read. `cost[]` is the byte plane at `world+0` (`TERR-COST-052`), indexed at the **destination** cell, read **inline** (`0054bde5 MOV CL,byte ptr [EAX + ECX*0x1]`). The diagonal is the straight cost `× 3/2` **truncated** (`0054bdf0 SHR AL,0x1`), so with the shipped cost bytes (6…16, `TERR-COST-052`) a diagonal is 9…24, and on a cost byte of 1 the two would be equal; in the flat arm the ratio is exactly 3:2. All arithmetic is 16-bit (`ADD AX,DI`) and the improvement test is unsigned strictly-less (`0054be92 CMP AX,word ptr [ESI + ECX*0x2]` / `0054be99 JNC`), which is what makes `0xffff` the unlabelled marker. **No heuristic term exists**: the Chebyshev distance `D = max(\|Δx\|,\|Δy\|)` is computed once (`00541f18…00541f4d`) and used **only** to size the iteration budget (`MOVE-TERM-003`). The same four arms appear in seven further places (`MOVE-SEARCH-001`, `MOVE-ROUTE-004`) with identical constants

**Confidence.** High (each of the four arms, the `SHR`, the operand widths and the compare are transcribed instructions; the rivals are excluded by the listing rather than by ratio-fitting — a table lookup or an FPU √2 would need instructions the arm does not contain, and "no straight/diagonal distinction" is refuted by the two separate arms)

**Original status.** ● active

**Evidence.** [EXP-0054](../experiments/EXP-0054-unit-movement/)

### MOVE-TERM-003

**Termination: three exits, one of them a generation budget — and no node budget, no frontier bound, no closed-set limit.** Tested once per generation, *before* the generation runs (`00541fff`, `0054201c`, `00542040`): (a) `label[goal] != 0xffff`; (b) the current frontier is empty; (c) the generation counter (`+2` per iteration, one iteration = both half-passes) reaches the budget. Budget = `scalar + D` for the n×n arms and `max(scalar, D>>2) + D` for the 1×1 arms, `scalar` being `StaticScanAhead` (5) or `DynamicScanAhead` (3) per `MOVE-PARAM-006`; the static 1×1 arm substitutes a flat **1000** when `actor[5]->[0x28] == 0` *and* the goal's own footprint is free on the static plane (`00542366`). *(Amended 2026-07-31 by EXP-0072, which named the term: `actor[5]` is `actor+0x14`, **the owning `Player`**, and `Player+0x28` is **0 exactly when a human participant owns the unit** — 1 or 2 for every scenario-authored owner, the constructor's own default being 1 (`UNIT-OWNER-009`). **So the flat 1000 is the NORMAL budget for a player-ordered move, not a rare case**, and an engine that implements only the two computed forms gives a human player's units `max(scalar, D>>2)` rings of slack where the game gives them a thousand generations — a long land detour round an obstacle fails and the unit does not move. Read at instruction level: `005422fe MOV EAX,[EBP + 0x14]` / `00542301 MOV ECX,[EAX + 0x28]` / `00542304 TEST ECX,ECX` / `00542306 JNZ 0054236b`, `EBP` being the actor — the same block uses `[EBP + 0x154]` for the mover and calls `[EBP]`'s `vt+0x1c`, the footprint getter. The second half is an `n×n` scan of the **goal** footprint against the mover's own domain mask `mover+0x5` on the static plane `world+0x10000` (`00542320`…`00542349`), `n` being `tokenSize`; any masked cell jumps past the override, and `n <= 0` skips straight to it. The store is `0054235d MOV dword ptr [ESP + 0x20],0x3e8`, into the same slot the two computed forms write at `005422ed`/`005422fa`.)* Since one generation advances the wave by one ring, the slack over the straight-line distance is only `max(scalar, D>>2)` rings — a detour costing more than that fails, and **the search itself** then substitutes a nearby goal and extracts to it (`FUN_0054bac0` / `FUN_0054b420`, `MOVE-ALT-018`…`022`). *(Amended 2026-07-31 by EXP-0068: this clause said "the **caller** then substitutes". It does not — both pickers are called from inside `FUN_00541e80`'s own tail and from nowhere else in the image, and no caller ever sees a substitute; `claims/retracted.md`.)* **The frontier arrays are unguarded**: the push is `MOV AX,[count]; MOV [x+AX]; MOV [y+AX]; INC [count]` with no comparison (`0054bea6…0054bec2`), the count is a u16, and each array is 4096 bytes — the four `0x1000` immediates in the whole driver are all the x→y stride, not a bound. Exit (a) also means the route is taken from the **first** generation in which the goal acquires any label, so with unequal step costs the label is not necessarily minimal when the search stops

**Confidence.** High (the three tests, both budget forms, the 1000 override and the unguarded push are named instructions; "no capacity test" is an immediate scan of the driver and of all five relaxation routines) / Medium (that a >4096-cell generation is therefore reachable on a shipped map: the overflow follows from the code, but no shipped-map census of generation widths was run)

**Original status.** ● active (amended)

**Evidence.** [EXP-0054](../experiments/EXP-0054-unit-movement/), [EXP-0068](../experiments/EXP-0068-alt-goal/)

**Amended.** The caller-substitution clause is partially retracted. Search selects its own substitute; MOVE-ALT-018 through MOVE-ALT-022 and retracted.md record the correction. Budget arithmetic stands.

### MOVE-ROUTE-004

**The route is extracted by walking the label field downhill from the goal, and its tie-break is asymmetric and favours the last straight step.** `FUN_005436b0` (dynamic list) and `FUN_005433a0` (static list) are the same code: start at the cell the search terminated on, and repeatedly choose the 3×3 neighbour minimising `label[nb] + stepcost(current cell)` — the *same* four cost arms as `MOVE-COST-002` — until the seed cell is reached. The **straight** arm accepts on `<=` (`0054385c CMP SI,…` / `00543861 JA skip`) and the **diagonal** arm only on `<` (`00543808` / `0054380d JNC skip`), over the scan order Δx = −1,0,+1 (outer) × Δy = −1,0,+1 (inner), **the centre included** — so among equal candidates the *last straight one in scan order* wins, a diagonal never displaces an equal straight, and the walk is not stable under a relabelling that only permutes equal costs. Guards: a candidate must be labelled (`CMP SI,0xffff`) and inside `8 ≤ x ≤ W+8`, `8 ≤ y ≤ H+8` (`world+0x50000` = W, `+0x50004` = H) — **the only rectangle test in the whole search**; passability is *not* re-tested and there is **no corner rule**, so a diagonal step between two blocked cells is legal. Output: a doubly-linked list of 12-byte pooled nodes (`+0x00` toward the goal, `+0x04` toward the unit, `+0x08` the packed cell `(y<<8)\|x`), built goal-first, with the unit's own cell unlinked before returning; a walk longer than **1000** steps discards the whole list (`0054372d`) — reachable, because the centre is a candidate

**Confidence.** High (both compares, the scan order, the bounds, the node layout and the 1000 cap are named instructions)

**Original status.** ● active

**Evidence.** [EXP-0054](../experiments/EXP-0054-unit-movement/)

### MOVE-PLANE-005

**Two searches over two planes, two routes on the actor — and the static search is structurally unit-blind.** `staticFlag != 0` reads the **static** block plane `world+0x10000` and writes the static list at `actor+0x160/0x164/0x168/0x16c/0x170` (tail/head/count/free/pool); `staticFlag == 0` reads the **dynamic** plane `world+0x20000` and writes the dynamic list at `actor+0x17c/0x180/0x184/0x188/0x18c`. Both use the same mask byte `mover+0x05`, and the reason the split works is the mask's own bits: `0x41` = terrain bit 0 **+ ground-occupancy bit 6**, `0x44` = static-object bit 2 + bit 6, `0x82` = border bit 1 **+ air-occupancy bit 7** (`TERR-PASS-051`) — the occupancy half of every mask is inert on the static plane, because bits 6/7 are never set there. Writers of those two bits, image-wide: `FUN_005456d0`, `FUN_0054abb0`, `FUN_0054ac70`, `FUN_0054ad20`, `FUN_0054af40` — **every one writing `world+0x20000`**

**Confidence.** High (the two plane displacements are in the relaxation routines' `TEST` operands; the mask bit meanings are `TERR-PASS-051`'s, re-derived here from the same three constants) / Medium (the writer list as an *enumeration*: two instruments, each with a blind spot the other covers only partly — `EnumRefs disp:10000`/`disp:20000` (60 hits, 19 owners, 0 orphan) cannot see a write whose address was folded into a register first, which is exactly how `FUN_0054ad20` writes; and `EnumRefs re:` over the four read-modify-write forms (`OR …,0x40/0x80`, `AND …,0xbf/0x7f`) **misses `FUN_0054af40` entirely**, which clears with a load/AND/store pair. Both routines are in the list because they were read, not because a scan found them)

**Original status.** ● active

**Evidence.** [EXP-0054](../experiments/EXP-0054-unit-movement/)

### MOVE-PARAM-006

**The six `[Path Finding]` scalars — the engine's own names, and the shipped file agrees with the code defaults 7/7.** `FUN_00547eb0` reads them through the `.reg` int getter `FUN_004ccc20(section, key, default)` (`RES-CODE-020`) into `world+0x585b4` `StaticScanAhead` (default **5**), `+0x585b8` `DynamicScanAhead` (**3**), `+0x585bc` `StaticRefreshRate` (**16**), `+0x585c0` `DynamicRefreshRate` (**32**), `+0x585c4` `DynamicByStaticLookup` (**3**), `+0x585c8` `StaticIsntNeeded` (**5**); `SpeedMultiplier` → `+0x58db4` (**8**) is the same section and is `TERR-MOVE-056`'s. `world.res:data/map.reg` `[Path Finding]` ships **all seven equal to those defaults, identically in the EN and RU roots** (`evidence/pathfinding-params.md`; `[Scanning] ScanShift = 7` beside them). Names are clamped to 15 characters in the file, matching the 15-char compare `TERR-COST-052` names

**Confidence.** High (the seven immediates and their store addresses are in one function's listing; the file values are the file's own bytes under EXP-0032's framing, and the 7/7 agreement is an independent check on the name→default pairing — a mis-attributed default would disagree with the file. The reading is from `tools/regcorpus`. It was **not** taken from `tools/regdump`, which until 2026-07-30 walked the pre-EXP-0032 framing — an 8-byte-shifted window pairing each key with the *next* record's value — and read this very section as `StaticScanAhead = 3`. That tool now carries the corrected framing and reads 5, agreeing with `regcorpus` and with the code default; the hazard is recorded because the figure was chosen to avoid it, not because it still exists)

**Original status.** ● active

**Evidence.** [EXP-0054](../experiments/EXP-0054-unit-movement/)

### MOVE-CLAIM-007

**There IS a reservation: a unit marks the cell it intends to enter, before it enters it.** `FUN_00549990` reads the next route cell from the dynamic list's tail into `mover+0x06`, and when that differs from `mover+0x80` it calls `FUN_0054af40(mover+0x80)` to drop the previous claim and **`FUN_0054ad20(mover+0x06)` to OR the mover's own occupancy bit over the n×n footprint of a cell the unit has not yet reached**, then stores it in `mover+0x80`. Every other unit's *dynamic* search therefore treats that cell as blocked (`MOVE-PLANE-005`), and the static search does not. The claim is released together with the vacated cell when the unit lands centred on the new one (`FUN_005495f0`: `mover+0xaa/0xac/0xa8 = 0`, `FUN_0054af40(mover+0x70)`, `mover+0x80 = 0`, `mover+0xa6 = 0`). The set/clear pair is **asymmetric**: `FUN_0054ad20` ORs unconditionally, while `FUN_0054af40` skips any cell inside the unit's *current* footprint — so releasing a claim never un-occupies the ground the unit is standing on. So the answer to "is there any structure recording a unit's intended position" is **yes, `mover+0x80`**, and the shared dynamic plane is the channel; `mover+0xa6` records the cell whose footprint currently carries the bits

**Confidence.** High (each call, its argument, the guard forms and the two mover fields are named instructions, `evidence/rom-move-excerpt.md` §6) / **Unknown** (whether a unit's own outstanding claim can block its own next dynamic search: the search lifts only its *current*-cell occupancy — `FUN_0054ac70` at `00542c49`, restored at `00542eae` — and never touches `mover+0x80`, but no reachability argument was built for the case where the claimed cell is not on the new route)

**Original status.** ● active

**Evidence.** [EXP-0054](../experiments/EXP-0054-unit-movement/)

### MOVE-WAIT-008

**When the next cell is occupied the unit waits, facing it — and it neither pushes, swaps, nor steps aside. At step time there is no collision test at all.** In `FUN_00548f70`, on a tick where a dynamic route is needed (list empty, or the per-tick counter `mover+0x78` over `DynamicRefreshRate`) *and* the next static waypoint is **exactly Chebyshev-1 away** *and* its n×n footprint is blocked on the dynamic plane, `FUN_00533540` decides: nonzero ⇒ the unit turns toward the cell (`FUN_0054a7d0` → `FUN_0054a210`), sets `actor[0x56]+8 = 10`, and **returns — no step, no re-search**; zero ⇒ fall through to the dynamic re-search `FUN_00549a90`. When the waypoint is farther than 1 the test degenerates to cell 0 (a border cell, always blocked) and `FUN_00533540` answers for that cell instead, so the wait is reachable only in the adjacent case. Nothing anywhere in the movement path displaces another unit. And `FUN_00549990` — the routine that actually steps — **tests no plane**: it claims the next cell and moves regardless of what is in it, so two units whose dynamic routes were both computed before either claimed can walk into the same cell

**Confidence.** High (the caller's control flow: the Chebyshev-1 gate, the footprint test's plane, the turn-and-return, and the absence of any plane test in `FUN_00549990`) / Medium (**which** blocker yields wait vs re-search: `FUN_00533540`'s arms were decompiled — the cell-record lookup `FUN_005463d0`, the actor list at the cell, the match on the other mover's `+0x70`, the `1`/`4` state byte — but its two helpers `FUN_00545c30` and `FUN_00544a00` were not read, so the verdict table is not pinned)

**Original status.** ● active

**Evidence.** [EXP-0054](../experiments/EXP-0054-unit-movement/)

### MOVE-TICK-009

**The per-tick update imposes no priority, so insertion order decides which actor gets a contested cell.** Simulation actors share `vt+0x18 = FUN_004f37be`, which reaches the order state machine. `FUN_0050fca9` walks one pooled linked list, invokes that virtual on each active element, and contains no sort, priority key or reordering. Its `[0x21,0x3f]` typeID test lies only in the `+0x54==0x10` teardown arm; it is not what proves the update elements' class, and it does not contain every Human because zero-mode map Humans can retain lower table typeIDs (`PARTY-M20-031`). A movement claim is written immediately, so the first actor visited claims the cell and later searches see it. The list is AddTail-only: map load inserts party then type-6 record order, death unlinks, spawns append, and save/load regroups by Player (`MOVE-TICK-013`..`017`)

**Confidence.** High for the virtual, whole loop and insertion-order writers; the corrected typeID scope is independently witnessed by original saves

**Original status.** ● active (corrected)

**Evidence.** [EXP-0054](../experiments/EXP-0054-unit-movement/), [EXP-0055](../experiments/EXP-0055-tick-order/), [EXP-0192](../experiments/EXP-0192-mission20-party-boundary/)

**Amended.** The teardown typeID range as proof of every list member being Human is partially retracted. Update order and list mechanics stand; PARTY-M20-031 and retracted.md record the correction.

### MOVE-STEP-010

**Movement is sub-cell, 1/256 of a cell per axis, and the position is a byte pair per axis.** The unit's position object is `actor[4]` = `*(actor+0x10)`: `+0x00` x cell, `+0x01` y cell, `+0x02` the packed cell `(y<<8)\|x`, `+0x04` x fraction, `+0x05` y fraction, **`0x80` = centred**. `FUN_00548c60` treats `(cell<<8)\|fraction` as one 16-bit value per axis and adds the signed per-axis step `mover+0xb0` / `mover+0xb1` (`TERR-MOVE-056` derives their magnitudes) to it, increments `mover+0xac`, and when the *cell* byte changes calls `FUN_0054abb0` on the ~~new~~ **old** cell (`00548da0` loads `word[pos+2]` before `00548dda` rewrites it; **corrected by [EXP-0093]**, `claims/retracted.md`) plus `FUN_00544d00` after the write; when `mover+0xac >= mover+0xaa` it snaps both fractions back to `0x80`. So a cell transit takes `mover+0xaa = ceil(256 / max\|step\|)` ticks and ends exactly at the cell centre, and a unit is *between* cells for most of its life — the packed cell is the near cell, not a rounded one. `FUN_00548f70` refuses to do anything but continue the transit while either fraction is not `0x80`. The steps are **signed** bytes (`00548cac`/`00548cb6 MOVSX`), and one direction-keyed special case (`00548cc7`/`00548cd5 CMP CX,0x1` / `0x5` on `mover+0xae`) subtracts 1 from the combined x at `00548ce3` when both fractions land on 0

**Confidence.** High (the 16-bit add, the snap condition, the boundary-crossing test and the field offsets are named instructions; the argument of `FUN_0054abb0` was a name, not a read, and fell)

**Original status.** ● active (amended)

**Evidence.** [EXP-0054](../experiments/EXP-0054-unit-movement/)

**Amended.** The boundary hook receives the old cell, not the new cell. That argument clause is partially retracted; MOVE-CLOCK-032 and retracted.md record the correction. The other step relations stand.

### MOVE-SPEED-011

**Speed does not feed the search.** The search reads exactly two things about the mover — the footprint side `vt+0x1c()` and the mask `mover+0x05` — plus the cost plane; no speed term enters any label. Speed enters only *after* a route exists: `FUN_0054d210` computes the per-step duration (`TERR-MOVE-056`: `SpeedMultiplier × unitSpeed`, height-tilted, divided by the two cells' mean cost), and `FUN_005495f0` recomputes `mover+0x72 = (actor+0x8c << 3) / cost[cell]` on each cell entry for a movementType-1 mover (0 for types outside 1..3, the raw class speed for 2 and 3). So two units of different speed choose the **same** route and traverse it at different rates, and the per-class speed source is `Data.bin` (`TERR-MOVE-057`/`058`) — **except under a group order **issued in formation**, where the rate is the group's slowest member's `Speed` and the two units traverse it at the SAME rate** ([EXP-0093], `MOVE-GROUP-030`; the *when* is `MOVE-GATE-035`, and because nothing clears the term the exception outlives the order that set it — `MOVE-GROUP-037`). The full composition is `MOVE-RATE-029`

**Confidence.** High (the search's whole mover interface is two virtual calls and one byte; the speed sites are named instructions)

**Original status.** ● active

**Evidence.** [EXP-0054](../experiments/EXP-0054-unit-movement/)

### MOVE-REFRESH-012

**When a route is recomputed — per cell for the dynamic one, per target change for the static one, and nothing is staggered.** (a) The **dynamic** route is destroyed on **every completed cell transit** — `FUN_005495f0` frees the whole list the moment the unit lands centred — so a dynamic search runs at least once per cell stepped. (b) It is also recomputed when `mover+0x78`, incremented once per `FUN_00548f70` tick and zeroed on re-search, exceeds `DynamicRefreshRate` (32) — the path that matters while a unit is waiting or turning rather than stepping. (c) The **static** route is recomputed only when the ordered target differs from `mover+0x74`, or when `mover+0x09` — the number of dynamic re-searches since the last static one — exceeds `StaticRefreshRate` (16); a static recompute frees the dynamic list. (d) The dynamic search does **not** aim at the final goal: `FUN_00549a90` takes the static list's tail waypoint when it is more than `DynamicByStaticLookup` (3) cells away, else the node `DynamicByStaticLookup+1` further along, and uses the final goal `mover+0x76` only when the static list holds `StaticIsntNeeded` (5) nodes or fewer; a static waypoint is popped once the unit is within `DynamicByStaticLookup` of it. Every counter is per-unit and reset on use, so **there is no stagger and no shared phase** — two units ordered in the same tick re-search in the same ticks. Nothing invalidates a stored route on terrain change or on target *movement*: only a new target cell (c), a completed step (a), the tick counter (b), and a blocked adjacent waypoint via `MOVE-WAIT-008`

**Confidence.** High (each counter, its threshold, its increment and its reset are named instructions, and the target-selection arms are `FUN_00549a90`'s own three branches) / Medium (calling `mover+0x78` a *stuck* counter: it increments on every tick of an unfinished move order, stepping or not, so it is a re-search period that a stalled unit reaches sooner — no separate give-up counter was found, and `mover+0x98`/`mover+0x90`, set on "no static route" and "arrived", are read but not traced to an effect -- *`mover+0x98` traced by [EXP-0099]: it is set on an empty route by `FUN_005492a0` at `005494a4` and consumed only by `FUN_005310e0`'s epilogue, which cancels the pending order and re-acquires within reach; see `AI-ROUTE-045`. `mover+0x90` traced by [EXP-0386]: it is a plain boolean, "actor already stands on the route's own resolved endpoint", not a route-list-emptiness test; see `MOVE-072`*)

**Original status.** ● active

**Evidence.** [EXP-0054](../experiments/EXP-0054-unit-movement/)

### MOVE-TICK-013

**The tick loop's container is a pooled doubly-linked list embedded in a 0x20-byte manager at `[0x00609558]`, its walk is head→tail, and its order is pure insertion history — the only insert operation that exists for it is AddTail.** Layout (ctors `FUN_0050fb5e` → `FUN_00518c30`): vptr `0x59caf0` @+0, then an MFC-shaped CObList @+4 — head @+8, tail @+0xc, count @+0x10, node free-list @+0x14, CPlex block chain @+0x18, blockSize 10 @+0x1c. Nodes are 12 bytes (`next`@+0, `prev`@+4, `element`@+8), pooled in CPlex blocks of 10 (`FUN_00519130` → `FUN_00570ba6(…,10,0xc)`), freed nodes pushed LIFO on the free list. The iterator (`FUN_00518160`/`FUN_005183a0`) advances **before** the element is processed, which is what makes the teardown's self-removal safe. The same class serves three more populations: the **dead list** at `*(world+0xc)`, a **per-player list** at `*(player+0x20)`, and the 0x48-byte **group** objects on `player+0x24`. Mutation surface, complete: insert = `FUN_00518de0` (AddTail; NewNode has exactly two callers, the other an inlined tail-append on a caller-local list in `FUN_0053ddd0`); remove = unlink (`FUN_00519710`/`FUN_005184e0`) via the by-value wrappers; RemoveAll; the serialize pair. **No AddHead, no InsertBefore/After, no sort and no compare-driven reorder exists in the family**

**Confidence.** High (every field offset, both ctors, the node functions and the iterator are named instructions, `evidence/rom-tick-excerpt.md` §2–3; the array rival is killed by the pointer walk and the unlink, the AddHead rival by `callto:519130` = 2 callers, the id-keyed rival by the iterator never reading `actor+4`. Enumeration instrument: `EnumRefs callto:` on every family function — 0 hits in orphan/undisassembled code)

**Original status.** ● active

**Evidence.** [EXP-0055](../experiments/EXP-0055-tick-order/)

### MOVE-TICK-014

**What fixes the order at map load: the creators' call sequence — campaign party first, then the map's type-6 records in ascending record order — and the list is per-*ticking-thing*, not total.** Every creator inserts through one routine, `FUN_0050fc0a([0x00609558], actor)` (AddTail + id, six callers, each loading the global immediately before the call): the party placement `FUN_004d403c`, the map type-6 walk `FUN_004e26bb` (index ascends 1..count over the map's record array; the record's own unit id is stored to `actor+8` at `004e2efe`), a trigger-parameter spawner `FUN_004f164c`, a roster walk `FUN_00504da1`, spawn-by-name `FUN_004d8fcd`, and sack creation `FUN_004feadb` — one of whose callers is `FUN_004f37be`, the actor tick itself, so a mid-tick death-drop appends to the very list being walked. Humans, Units and Sacks are members; **no Building insert into this list was found**, though buildings draw ids from the same allocator (direct `FUN_004d9fed` callers with no list insert: `FUN_0050f3c3`, `FUN_00504af0`, `FUN_0051018e`) — which is how `SAV-ID-015`'s ids interleave buildings 2..19 between the hero and the units while the tick walk carries no buildings

**Confidence.** High (the six-caller enumeration: `callto:50fc0a`, 0 orphan, each site's `ECX=[0x00609558]` shown by `refto:609558` adjacency; the unit walk's ascending index and the `actor+8` store are named instructions) / Medium (that the record array the walk indexes = the `.alm` type-6 **file** order: the accessor `FUN_0051c8b0` was not chased into the map loader — `SAV-ID-015`'s 35/35 file-order ids corroborate on one map) / Medium ("buildings never enter this list" — an absence over the enumerated insert surface, blind to an insert through a wrapper this experiment did not identify; whether buildings tick from another container — `world+0x2c`, `*(world+0)/(+4)/(+8)` — is open)

**Original status.** ● active (contested)

**Evidence.** [EXP-0055](../experiments/EXP-0055-tick-order/)

**Amended.** Building membership remains contested: no Building insert was found over the named insert surface, but untraced wrappers and other tick containers remain open. The record-array to file-order link also remains Medium.

### MOVE-TICK-015

**How the order changes during play — and death leaves no hole.** (a) The teardown arm (`actor+0x54 == 0x10`) removes the actor's node by **unlink** (`0050fd44`) — nothing is compacted, swapped from the tail, or replaced; every survivor keeps its relative position — then AddTails the actor to the **dead list** `*(world+0xc)` (`0050fd59`), where `FUN_0050fed6` runs corpse decay per tick. (b) New spawns/summons append at the **tail**. (c) **Leaving the map** (into a building/transport, `FUN_004f47e6`) unlinks from the tick list and sets `actor+0x4c` bit 3; **returning** (`FUN_004f4865`/`FUN_004f4905`, "Unit can't return to map — no free place") re-AddTails — the returning unit moves to the **back** of the walk, losing its old contention position. (d) An **owner change** (`FUN_004d1e14`) moves the actor between the players' own lists and groups but never touches the tick list — its walk position survives a conversion. So mid-game, the walk order = spawn order, minus the dead and the garrisoned, plus returners and newcomers at the tail

**Confidence.** High (each arm is a raw listing with the `ECX` chain shown, `evidence/rom-tick-excerpt.md` §3, §5, §6; the replace-in-place and periodic-resort rivals are excluded by the complete mutation surface of `MOVE-TICK-013`)

**Original status.** ● active

**Evidence.** [EXP-0055](../experiments/EXP-0055-tick-order/)

### MOVE-ID-016

**The runtime id (`actor+0x04`, `SAV-ID-015`'s) is a lowest-free-bit bitmap allocation — creation-ordered only until the first death decays, and never the walk order.** `FUN_004d9fed` scans the bitmap `0x62c7e0` for the first clear bit (id 0 pre-marked by `FUN_004d9f4a`, which runs at world construction and again at save-load start), marks it, and `FUN_0050fc0a` stores it at insert time. `FUN_004d9fac` clears a bit; its actor-side caller is the **corpse decay**: at stage 5 (`actor+0x94` reaching −600) `FUN_004f52ee` frees the id and zeroes `actor+4` (`004f53f6`/`004f5401`) — which is why long-dead units serialize with id 0, and why the **next** spawn after a decay reuses the lowest freed id while still appending at the tail. The id space is shared with non-ticking placeables (three direct allocator callers with no tick-list insert). On load the id is read back from each actor's head and **re-marked** (`FUN_00510e5c` at `SetAt`… listing §7), so ids round-trip exactly

**Confidence.** High (allocator, set, clear, reset, the insert-time store, the decay-time free and the load-time re-mark are named instructions; `refto:62c7e0` = exactly 4 touching functions)

**Original status.** ● active

**Evidence.** [EXP-0055](../experiments/EXP-0055-tick-order/)

### MOVE-TICK-017

**The walk order is NOT preserved across save/load — the saved stream never carries it, and the loader rebuilds it grouped by player.** The world Serialize (`FUN_004d0cb7`) stores: the players (each Player storing its groups inline, each group WriteObject-ing **its own actor list head→tail** — `FUN_00511089` → `FUN_005102f4` → `FUN_00511938` → `FUN_00527ac0`), then the dead list; the tick list appears nowhere in the store arm's exhaustively-walked call sequence, and EXP-0048's byte-level tiling of four saves corroborates (393/393 objects under players/groups/dead, ≤ 70 bytes unattributed — no room for a 400-word reference list). On load, Player::Serialize rebuilds `*(player+0x20)` from its groups (groups in list order, actors in group order, `005113f9`), and the world's load arm then rebuilds the tick list: **for each player in players-manager order, append `*(player+0x20)` head→tail, skipping off-map (`+0x4c` bit 3) actors** (`004d10ab…004d1120`). Consequence, joint with `MOVE-TICK-009`: a fresh map walks global creation order (players interleaved as the records interleave); a loaded game walks player-by-player, group-by-group — relative order within one group survives, the cross-player interleave does not — so **which of two contending units of different players moves first can change across one save/load cycle**

**Confidence.** High (both rebuild loops and the store sequence are raw listings, §7–§8; the round-trip rival is killed by the absence of the list from the store arm *and* by the rebuild's existence) / Medium (the regrouping's *observable* consequence on a real contention: derived from the loops, not yet demonstrated in a live session — the named follow-up)

**Original status.** ● active

**Evidence.** [EXP-0055](../experiments/EXP-0055-tick-order/)

### MOVE-ALT-018

**The substitution is not the caller's — `FUN_00541e80` does it itself, in three branches, and `altTarget` is a target *actor*, not a cell.** The two pickers have **two and one call sites, all three inside `FUN_00541e80`** (`00542f4c`, `00543028`, `00543011`) and are reached from nowhere else in the image. They run in the tail entered when the goal is still unlabelled (`00542eca CMP word ptr [ESI + ECX*0x2],0xffff` / `00542ed0 JZ`). Which one runs is decided by the `staticFlag` argument at `00542f2a` and the `altTarget` argument at `00543001`: **static** → picker A around the requested cell with bound `(D>>2) + 4` (`00542f3a SHR AX,0x2` / `00542f3e ADD EAX,0x4`), `altTarget` **ignored**; **dynamic** → `altTarget != 0` ? picker B(`mover`, `altTarget`) : picker A with bound **8** (`00543020 PUSH 0x8`). The answer is fed straight to the ordinary route extraction with the substitute as `(dstX, dstY)` — `FUN_005433a0` at `00542ff2`, `FUN_005436b0` at `00543045` (`MOVE-ROUTE-004`) — and a zero answer means no route (`00542fc2…00542fcc` frees the static list through `FUN_0051c510(actor+0x15c)`; the dynamic branch simply returns). **`altTarget` is an object**: picker B reads `arg2+0x10` as the position record (`MOVE-STEP-010`) and calls `arg2->vt+0x1c()` for its footprint side. The only caller that passes a non-zero value is `FUN_005492a0`, the go-to-actor order, via `FUN_00549a90(mover, targetActor)` at `0054959c`; `FUN_00548f70` passes 0 at both `00549056` and `00549265`, and `FUN_005492a0`'s own **static** call at `00549450` passes the actor into an argument the branch above discards

**Confidence.** High (the two `EnumRefs callto:` sweeps — 2 hits / 1 owner and 1 hit / 1 owner, 0 in orphan or undisassembled code — bound the surface; the three branch tests, both bounds and the two extraction calls are named instructions, `evidence/rom-altgoal-excerpt.md` §3. The rival "a caller substitutes" is excluded by the absent call sites, not by plausibility; the rival "`altTarget` is a cell" by the two dereferences)

**Original status.** ● active

**Evidence.** [EXP-0068](../experiments/EXP-0068-alt-goal/)

### MOVE-ALT-019

**Picker A (`FUN_0054bac0`): expanding square rings around the requested cell, and the whole ring is scanned before the best in it is taken.** `__thiscall(world, mover, unused, cell, limit)` — the second argument is `0` at both call sites and is never read. Rings `r = 1, 2, 3, …`; for each `i = -r..r` it probes the four cells `(x+i, y+r)`, `(x+i, y-r)`, `(x+r, y+i)`, `(x-r, y+i)` — the four `LEA`/`ADD` forms at `0054bb3b`, `0054bb59`, `0054bb7a`, `0054bba2` over the packed cell `(y<<8)\|x`, with `i` walked by `0054bbc9 INC EBX` / `0054bbca CMP EBX,ECX` (`ECX = r+1`) and `i`'s start decremented once per ring at `0054bbe4 DEC EBX`. Each probe is **one instruction**, `MOV DX,word ptr [ESI + EDX*0x2 + 0x30000]`, and the keep test is `CMP DX,AX` / `JNC` — the **strict** minimum. The ring is scanned to the end; if it produced any labelled cell at all the loop exits (`0054bbda CMP AX,0xffff` / `0054bbde JC`), otherwise `r` grows while `r+1 < limit` (`0054bbe5 CMP ECX,EDX` / `0054bbef JL`), and `limit <= 1` skips everything (`0054baf6`). The **centre is never probed** (`r` starts at 1), and correctly so — the branch is only entered when the goal is unlabelled. Return: the packed cell, or **0** if no ring had one (`0054bd00 INC AX` / `NEG AX` / `SBB EAX,EAX` / `AND EAX,ECX`). **The mover's footprint does not change the scan**: `0054bae2 CMP ECX,0x2` on `mover->vt+0x1c()` selects between two copies of the same code whose 199-byte inner bodies differ in exactly **4 bytes**, the operands of a swapped register-restore pair, with no jump displacement among them (`evidence/image.txt`)

**Confidence.** High (every instruction above is transcribed from a full listing of a 595-byte routine, not from an enumeration, and the routine's bytes are hashed in `evidence/image.txt`. Check **D1** (`evidence/ring.txt`) re-executes the four side formulas and confirms they cover the Chebyshev ring of radius `r` exactly for `r = 1..8`, corners twice and nothing else duplicated — a sign error would break it. The rivals die on absent instructions: a spiral or a raster needs a different index recurrence than `±r`/`±(r<<8)`; a distance-sorted list needs a sort, and there is none; "stop at the first labelled cell" is excluded because the exit test sits **after** the inner loop, not inside it)

**Original status.** ● active

**Evidence.** [EXP-0068](../experiments/EXP-0068-alt-goal/)

### MOVE-ALT-020

**Picker B (`FUN_0054b420`): the contact ring around the target actor, entered where the line between the two actors crosses it and walked in both directions.** `__thiscall(world, mover, target)`. The box is `x ∈ [tx - nM, tx + nT]`, `y ∈ [ty - nM, ty + nT]`, built from four `vt+0x1c` calls — two on the **target** (`0054b451`, `0054b47d`, giving `nT`) and two on the **mover** (`0054b4ac`, `0054b4d4`, giving `nM`) — and it grows by one in each direction per ring, for **8 rings** (`0054b6b8 CMP ECX,0x8` / `JGE`, expansion at `0054b9c7…0054b9ef`), the loop continuing only while nothing has been found (`0054b9f9 CMP ECX,0xffff` / `JZ`). The entry cell is where the straight line between the two actors' **fine footprint centres** meets that box: `FUN_0054b2c0` returns a 16-way direction code, folded to a quadrant by `(dir + 2) >> 2 & 3` (`0054b50d…0054b516`), which selects one of four edges through the inline jump table at `0054ba48`; the crossing is `slope × fine_edge + intercept` → `__ftol` → `SAR 8`, with `slope` = `dy/dx` or `dx/dy` per quadrant (table `0054ba38`, two distinct arms) and a divide-by-zero guard that turns `dx == 0` into `1` (`0054b596 FCOM float ptr [0x0059cda0]` = `0.0` / `0054b5a3 FSUB float ptr [0x0059cda4]` = `-1.0`). Two walkers then run the perimeter from that cell in opposite directions, stepping by the 12-entry `(dx,dy)` byte table at `world+0x54186` — `(+1,0) (0,+1) (-1,0) (0,-1)` **three times**, written only by the world constructor (`EnumRefs disp:54186/54187/54188` → 2 / 2 / 1 hits, 0 orphan; `EBX = 0` from `00547ee5`, the only write to `EBX`/`BX`/`BL` before its epilogue) — turning at a corner through the two 12-entry jump tables at `0054ba58` / `0054ba88`, each of which is four handlers repeated three times. Each probe is again **one** label-plane read (`0054b7ff`, `0054b841`) kept by `CMP ECX,EAX` / `JGE`, the strict minimum. Return: the packed cell `(y<<8)\|x` (`0054ba13`), or **0** (`0054ba2d XOR AX,AX`)

**Confidence.** High (the box, the 8-ring bound, the two walkers, the direction table, the four jump tables and both label probes are named instructions over a full listing of a 1559-byte routine whose bytes are hashed; the jump tables are read out of the image, not off the decompiler, in `evidence/image.txt`. Check **D2** (`evidence/box.txt`): at ring 0 the box perimeter is exactly the set of mover origins whose `nM × nM` footprint touches the target's `nT × nT` footprint without overlapping it, for every `(nM,nT)` in `1..4` — so "contact ring" is measured, not asserted) / Medium (the entry-cell arithmetic: the four FPU arms were read and their *shape* — line, slope, intercept, truncate — is transcribed, but the quadrant→edge mapping was not checked cell by cell against a worked example, and no consumer of the entry cell's exact value was tested. What would lift it: re-execute the four arms over a grid of `(dx,dy)` and confirm the entry cell always lies on the box)

**Original status.** ● active

**Evidence.** [EXP-0068](../experiments/EXP-0068-alt-goal/)

### MOVE-ALT-021

**What a candidate is tested against, and what the choice is measured from — and the answer to both is the label plane, so the substitute is the cell cheapest to reach *from the mover*, not the cell nearest the click.** Across **both** pickers read in full, the only world displacements are `+0x30000` (the u16 label plane, `MOVE-SEARCH-001`) — eight sites in `FUN_0054bac0`, two in `FUN_0054b420` — and `+0x54186/7` (picker B's direction table). Neither routine contains a passability predicate, a cost read or an occupancy test, and neither calls one: the complete callee set is the footprint virtual `vt+0x1c`, `FUN_004f280d`/`FUN_004f2849` (which read `pos+0/1` and `pos+4/5` and add `(side-1)<<7`, reaching only `FUN_005449e0`/`FUN_005449f0`), `FUN_0054b2c0` (which calls only those two) and `__ftol`. So **the test is "this cell carries a label"** — the verdict of the wave that has just failed, which already folded in the plane it ran on, the mover's footprint and its mask (`MOVE-PLANE-005`, `TERR-PASS-051`) — and **the metric is the label itself**, the accumulated step cost from the mover's own start cell under `MOVE-COST-002`. Two consequences a consumer must implement. (a) The rival "the nearest free cell to the clicked cell, by Chebyshev, scanning outward" is **half right and half wrong**: the ring order *is* Chebyshev outward from the requested cell, but within a ring the whole ring is scanned and the cheapest-from-the-mover cell wins, and "free" is *not* passability — a perfectly passable cell that the wave's generation budget never reached carries no label and is invisible to the picker. (b) Substitution is **per-mover**, on that mover's own failed wave: on the dynamic branch the labels were laid over `world+0x20000`, which carries other movers' claims (`MOVE-CLAIM-007`), so a claimed cell is never labelled and can never be returned; on the static branch `world+0x10000` carries no occupancy, so a static substitute is chosen as if the mover were alone on the map

**Confidence.** High (an **exhaustive** read of two hashed byte ranges, which is a stronger instrument than any enumeration and has no orphan-code blind spot: the claim is not "no other reader was found" but "these are all the instructions". The three rivals are excluded by absent instructions — a distance-to-click metric needs a subtraction of the two cells and there is none; a distance-to-mover metric needs the mover's position, which picker A never loads; a fresh passability test needs the block plane, whose displacement appears nowhere in either routine)

**Original status.** ● active

**Evidence.** [EXP-0068](../experiments/EXP-0068-alt-goal/)

### MOVE-ALT-022

**The substitute is never written back — it is consumed by one route extraction and forgotten.** Over the whole tail `00542eb3…00543054` the only stores reaching the actor are the route list the extraction builds (`actor+0x160…` / `actor+0x17c…`, `MOVE-PLANE-005`) and `FUN_0051c510(actor+0x15c)` freeing the static list when no substitute was found. `mover+0x74` (the last ordered static target) and `mover+0x76` (the final goal) are **not** touched, the order block `actor+0x158` is not touched, and the picker's return value never leaves `FUN_00541e80` — it is pushed to the extraction call and to nothing else. The static branch does *measure* the substitute against the request — `FUN_0054d100(substitute, requested)`, the Chebyshev distance of two packed cells — with a tolerance of **1**, raised to **2** when the requested cell has static-plane **bit 5** set (`00542f61 TEST byte ptr [ESI + ECX*0x1 + 0x10000],0x20`, one of `TERR-PASS-051`'s 13 bit-5 sites) *and* the cell-record lookup `FUN_0054fc10(world+0x540b4, cell)` returns a record whose `+0x0c` is non-zero; exceeding the tolerance queues a UI message (`FUN_004ea281` with `ECX = 0x603c28`) and **does not cancel the move** (`00542fbd TEST BX,BX` gates the extraction on the substitute existing, not on the tolerance)

**Confidence.** High (the tail is listed in full and hashed, so "no store" is an exhaustive read rather than an enumeration; the tolerance test, both constants and the message call are named instructions) / Medium (the run-time consequence — that the order still points at the original cell, so the unit re-runs the same failing search whenever `MOVE-REFRESH-012`'s counters fire, and substitutes again from wherever it now stands: this follows from the two rows but was not observed. `FUN_0054fc10`'s record and its `+0x0c` are named, not read)

**Original status.** ● active

**Evidence.** [EXP-0068](../experiments/EXP-0068-alt-goal/)

### MOVE-ORDER-023

**There is no multi-unit destination distribution anywhere: no formation, no offset table, no spread. A group move writes the same cell into every member's order block.** Every site that sets the move opcode is the instruction form `MOV byte ptr [reg + 0x8],0x1` on an order block loaded from `actor+0x158` — **16 sites over 12 owners, 0 in orphan or undisassembled code**. Nine owners write one actor's order per call. The three that walk a container of actors are the group orders `FUN_00536ef0`, `FUN_005370f0` and `FUN_00537510`, and only `FUN_005370f0` takes a cell: `__thiscall(issuer, group, cell)`, reached from the group order machine `FUN_00533ae0` case 2 with the group order's own `+0x0a`. Its cell is loaded **once**, `00537121 MOV BP,word ptr [ESP + 0x18]`, *before* the loop head at `00537126`; that is the **only** write to `BP`/`EBP` in the routine; the loop body contains no counter, no index and no arithmetic on it; and the store is `00537196 MOV word ptr [EDX + 0xa],BP` for every member the list walk reaches. The other two re-issue each actor's *own* already-stored cell (`order+0x00` at `005370a2`, `order+0x0a` at `005376cd`). So the spread a player sees when several units are sent to one tile is produced entirely by `MOVE-ALT-019`/`021` — each mover substituting independently against a plane that already carries the earlier movers' claims — and not by anything at order time

**Confidence.** Medium (the enumeration is over **one instruction form**: a move order set with the opcode in a register, through a store wide enough to span `+0x8` and its neighbour, or by copying a whole order block, is invisible to it — `EnumRefs re:` sees orphan code, so that half of the blind spot is closed, but the encoding half is not. What is High inside the row is `FUN_005370f0` itself: the single `BP` load, its position before the loop head, and the absence of any other `BP` write are a full listing of a 226-byte routine, hashed in `evidence/image.txt`. What would lift the row: a `disp:158` sweep classified owner by owner, or locating the order block's class and enumerating its writers through the vtable)

**Original status.** ● partially retracted — **the headline is refuted by [EXP-0094]**, through this cell's own named blind spot: the per-member issue routine `FUN_0052f4f0` sets the move opcode from a **register** (`0052f5eb MOV byte ptr [EDX + 0x8],AL`, `EAX` loaded `1` at `0052f5ca`), and on the formation arm of `FUN_005340a0`/`FUN_00534390` it is called once per member with `target + (memberCell − centroidCell)`. Formation, offset table and spread test all exist (`MOVE-FORM-036`, `MOVE-GATE-035`, retracted.md). What survives untouched: the three routines this row read are still right about themselves

**Evidence.** [EXP-0068](../experiments/EXP-0068-alt-goal/), **[EXP-0094](../experiments/EXP-0094-group-rate/)**

**Amended.** The no-formation headline is partially retracted. MOVE-FORM-036 and MOVE-GATE-035 establish the counterexample; retracted.md records its scope. The three selected routines retain their bounded local conclusions.

### MOVE-DOM-024

**The movement domain is one byte with six consumers, and passability is only one of them.** `actor+0x4a`, read through slot **`+0x20`** of all three simulation-actor vtables (`0x59c3c0`/`0x59c448`/`0x59c4d0`, whose `+0x20` all hold `FUN_00523230` = `MOV AL,byte ptr [ECX+0x4a] ; RET`; `TERR-MOVE-055`). What the code branches on it for: (a) the **block mask** — `FUN_0054b120`, `MOVE-DOM-025`; (b) **which occupancy bit** the mover sets and tests on the dynamic plane — `FUN_0054abb0`/`FUN_0054ac70`/`FUN_0054ad20`/`FUN_0054af40`, `< 3 → 0x40`, `== 3 → 0x80`; (c) the **step-cost arm** — `== 1` reads the cost plane, every other value takes flat 2/3 (`MOVE-COST-002`), at nine sites across the driver, the three relaxation routines and the two route extractions; (d) the **speed** — `FUN_0054d210` skips the height tilt and the cost divide for `!= 1`, and `FUN_0054a620` (read whole) returns `(class speed << 3)/cost[cell]` for 1, the **raw class speed** for 2 and 3, and **0** for 0 and for anything above 3; (e) the **AI's target choice** — a candidate whose domain is 3 counts one cell farther away when the decider's is not (`0052e18f…0052e19e`, `AI-ACQUIRE-002`), which is the image's own word for *flier* and the only attestation of that name that does not come from the mask; (f) the **corpse** — the actor tick reads the byte directly at `004f38fc` and, for `> 1`, stores `0xfc18` (−1000) into the decay counter `actor+0x94` (`004f3907`), which the next test at `004f391a` turns into `actor+0x54 = 0x10`, teardown — so a non-ground mover leaves no corpse, which is the other half of `ANIM-DEATH-007`'s “the four `movementType > 1` classes leave no corpse at all”

**Confidence.** High (the getter and its three vtable slots, the mask arms, the two occupancy immediates, `FUN_0054a620`'s four arms and the corpse slam are named instructions, `evidence/rom-domain-excerpt.md` §1, §2, §6; the getter's identity comes from a `.rdata` slot read, the instrument that *can* see a virtual-only target) / Medium (that the consumer list is **complete**: two instruments, `EnumRefs disp:4a` (32 hits / 0 orphan — the direct reads) and `re:CALL dword ptr [E..+0x20]` (229 hits / 128 owners / 0 orphan, of which the actor-family receivers are listed in §6). Neither sees a call through a computed pointer, and six of the listed owners — `FUN_0052af90`, `FUN_0052e4d0`, `FUN_00535d30`, `FUN_00535e90`, `FUN_00537ee0`, `FUN_004fd37a` — were located and **not read**)

**Original status.** ● active

**Evidence.** [EXP-0078](../experiments/EXP-0078-movement-domains/)

### MOVE-DOM-025

**The mask is installed once, at spawn, is the mover's only passability state, and survives a save because the mover is serialized whole.** `FUN_0054b120` is `__thiscall(mover, actor)`: it calls `actor->vt+0x20()`, then `DEC EAX ; JZ` three times — `1 → MOV byte ptr [ESI+0x5],0x41` (`0054b152`), `2 → 0x44` (`0054b14a`), `3 → 0x82` (`0054b142`) — and a value outside `{1,2,3}` falls through to `0054b156` and **stores nothing**, leaving whatever was there. It has **2 call sites** (`EnumRefs callto:54b120`, 2 owners, 0 orphan): `FUN_004f59de` @`004f6ad9`, the Units-table spawn constructor's tail, and `FUN_004f6ded` @`004f6f00`, the Humans-table defaults routine — both inside object construction, neither reachable afterwards. The byte's **complete direct-write population** is four instructions: the mover constructor `FUN_00545c00` @`00545c0e` (`0x41`, after a `REP STOSD` of `0x2d` dwords = the mover's `0xb4` bytes) and `FUN_0054b120`'s three (`EnumRefs re:` on the store form — 33 hits / 25 owners / 0 orphan, every other hit on a different object, most of them `MOVE-STEP-010`'s `pos+0x04`/`+0x05` fraction pair written as a `[+0x4]`/`[+0x5]` couple). **So the mask is a cache with no refresh** — and the rival that follows from that, *a deserialized actor keeps the constructor's `0x41`*, is real all the way to its last step and then false: `Unit::CreateObject` (`FUN_004f2afd`) runs the **default** ctor `FUN_004f2bb4`, which pushes the literal `"null"` (`0x62e83c`, `StrDump`) into the `Data.bin` name search; no shipped Units row is named `null`; and the search's failure arm reports through `FUN_00437040` and `004f5c96 JMP 0x004f6af0` — **past** the mask install. It dies because `Unit::Serialize` `FUN_00510518` calls `FUN_0054d4c0` at `005105d8`, *before* its own `CArchive::IsStoring` test at `005105f2`, and `FUN_0054d4c0` is eleven instructions that push `0xb4` and the mover and call `CArchive::Read` (`0057ba16`) or `Write` (`0057b908`) — the mask is byte 5 of that block and round-trips verbatim

**Confidence.** High (every instruction above is transcribed from full listings of five routines, `evidence/rom-domain-excerpt.md` §2 and §7; the `"null"` literal is `StrDump`'s and the three `CRuntimeClass` records are read out of the PE section table with no disassembler in the loop; the save rival is refuted by a **present** instruction, not by an absent one) / Medium (the write enumeration's blind spot, which is exactly what the rival tripped on: a `re:` sweep over one store form cannot see a **wholesale copy** of the mover, and `FUN_0054d4c0` is such a copy)

**Original status.** ● active

**Evidence.** [EXP-0078](../experiments/EXP-0078-movement-domains/)

### MOVE-DOM-026

**What each domain is permitted to enter, terrain and non-terrain counted apart — and the axis is not terrain, it is which of the three block bits the mask carries.** Over the 38 shipped maps' static block plane, rebuilt as the ingest (`TERR-PASS-049`/`050`) and the structure pass (`TERR-STRUCT-068`/`070`/`071`) build it and classified by the arm that set each byte (`tools/domainstop -mode stops`, `evidence/stops.txt`): domain **1** (`0x41`) is stopped by **199 361** terrain cells (Mountain 117 572, water 77 845, tile-word bit 13 3 944) **and** 242 179 non-terrain (type-3 object 69 676, structure 13 527, border 158 976) = **441 540** of 880 704 (50.13 %); domain **2** (`0x44`) by **0** terrain and the same 242 179 non-terrain (27.50 %); domain **3** (`0x82`) by **0** terrain and **158 976** — the border alone — (18.05 %). So **domain 2 passes water and mountain and is stopped by every object, building and ground occupant, while domain 3 is stopped by nothing on shipped terrain except the 8-cell border and another domain-3 occupant**: they are not two grades of one thing, and they disagree on **83 203** cells. Runtime occupancy (bits 6/7) is by construction invisible to a static census and is stated separately: the same masks make bit 6 block domains 1 and 2 and bit 7 block domain 3 (`MOVE-DOM-024`(b), `TERR-PASS-051`). Consequence a route consumer must reproduce: the 8-connected free set — the search steps 8 neighbours and applies **no** corner rule (`MOVE-ROUTE-004`) — falls into **more than one piece on 23 of 38 maps for domain 1 and on 0 of 38 for domain 3**, whose largest piece is the whole interior on every map (`evidence/cross.txt`)

**Confidence.** High (the rule is `TERR-PASS-051`'s named instructions and the arm attribution is the ingest's own four assignments plus the structure pass's `OR 5`/`AND 0xfa`; the discriminator is that **bit 1 is set on 158 976 cells and 0 of them by any arm but the border stamp** — an exhaustive re-execution, not a sample — so domain 3's stop set is measured rather than inferred, and the `==`-comparison rival scores **0 / 0 / 0** blocked cells) / Medium (the figures are one corpus on one root; and the attribution reports the arm that wrote a byte **last**, which the ingest's four assignments make well defined but which cannot separate two arms that agree on a cell)

**Original status.** ● active

**Evidence.** [EXP-0078](../experiments/EXP-0078-movement-domains/)

### MOVE-DOM-027

**Nothing narrows a mover's verdict after the mask: no height term, no corner rule, no order-time check, no step-time check.** (a) Both passability predicates were read **end to end** — `FUN_00543060` (static plane) and `FUN_00541dd0` (dynamic) are instruction-for-instruction identical apart from the plane displacement: `n = actor->vt+0x1c()`, `n <= 0` returns free, and the `n × n` walk's only memory operand is `TEST byte ptr [cell + plane],CL` with `CL` **re-loaded from `mover[0x154]+5` on every cell**. No cost read, no height read, no edge or adjacency term, no reachability seed. (b) The **height plane has six references image-wide** (`EnumRefs disp:9451c`, 5 owners, 0 orphan): the line-of-sight pair `FUN_00546c70`/`FUN_00546d20` (`AI-SIGHT-006`), the constructor's memset, the ingest's write, and the two reads in `FUN_0054d210` (`TERR-MOVE-056`) — **none in the search**, so height gates sight and speed and never passability. (c) The two block planes have **no reader below `0x00523be0`** over `EnumRefs disp:10000` (78 hits / 33 owners) + `disp:20000` (60 / 19), 0 orphan — every other owner lies in `0x00541dd0…0x0054f680` — so no interface, input, campaign or view routine can consult passability, and there is **no click-time or order-time refusal of a destination anywhere in the image**; a group move writes one loop-invariant cell into every member's order block with no test at all (`MOVE-ORDER-023`). (d) The step routine `FUN_00549990` tests no plane (`MOVE-WAIT-008`), and the mask's only other consumer at step time is the claim, which *writes* (`MOVE-CLAIM-007`)

**Confidence.** High for (a) and (d) — full listings, and “no other operand” is an exhaustive read rather than an enumeration / Medium for (b) and (c) as **enumerations**: a `disp:` sweep cannot see an access through a pointer the routine `LEA`'d first (the height plane's one `LEA` is the memset, whose destination never leaves the constructor; both block planes' `LEA` hits are followed to their stores in `TERR-PASS-073`) nor a wholesale copy of the containing structure

**Original status.** ● active

**Evidence.** [EXP-0078](../experiments/EXP-0078-movement-domains/)

### MOVE-DOM-028

**Which shipped classes can be non-ground — and that a player can never command one.** `Data.bin` **Humans**: 210 parameterised rows, `movementType` (slot `0x16`) is **−1 on all 210**, so the base constructor's `1` stands and every human, hero included, is a ground mover. `Data.bin` **Units**: 57 named rows, `movementType` (slot `0x20`) is 1 on 40 and non-1 on **16** — the four difficulty variants each of **Ghost** and **Bee** (2) and **Bat_Sonic** and **Dragon** (3); no row anywhere carries a value outside `{−1,1,2,3}`. The two domain-3 classes are exactly the two with `units.reg` `Z != 0`, i.e. exactly the two allocated as `CAirUnit` and drawn in layer 3 (`REG-UNITS-061`) — an independent second witness for the name *air*, the third being `AI-ACQUIRE-002`'s distance penalty (`MOVE-DOM-024`(e)). The player's own roster contains none of them: the hero is a Human, and the tavern's fifteen hireable types are `Catapult` and `Ballista` (Units, `movementType` 1) plus thirteen `NPC%02d_%d` **Humans** (`MERC-TYPE-001`). Over the 38 maps, all **1739** non-ground type-6 placements (1024 domain-2, 715 domain-3, resolved by `ALM-CLS-038`'s key rule; 4 of 8094 records unresolved) belong to scenario-authored groups — `Monsters`, `Nocturnal`, `Beasts`, `Enemy`, `Ghosts`, `Wild`, `Evil-Creatures`, … — whose `Player+0x28` is 1 or 2, never a human participant's 0 (`UNIT-OWNER-009`). **The corroboration that would fail if the domain assignment were wrong:** scored against its own mask, **0 of 1739** non-ground placements stands on a cell that blocks it, while **258** of them (213 domain-2, 45 domain-3) stand on cells the **ground** mask blocks — against 4 of 6351 for ground placements on either scoring (`evidence/who.txt`)

**Confidence.** Medium (the value tables and the placement scoring are a corpus census on one root, which caps here by rule even though the null model fails hard — a wrong domain assignment cannot make 258 in-terrain placements come out 0/1739 — and the hireable roster is `MERC-TYPE-001`'s reading rather than this experiment's) / **Unknown** (whether a *trigger* can hand a non-ground actor to a human participant mid-mission: the spawners `FUN_004f164c` and `FUN_004d8fcd` and the owner-change `FUN_004d1e14` exist and were not walked over the shipped trigger set, and `Control Spirit` spawns a `Ghost` owned by the caster — `MAGIC-SING-019` — which is a domain-**2** mover a player *can* end up owning)

**Original status.** ● active

**Evidence.** [EXP-0078](../experiments/EXP-0078-movement-domains/)

### MOVE-RATE-029

**The whole rate law, composed, with the instruction for every term — and the composition is five terms, not one.** `FUN_0054d210(world; actor, u16 srcCell, u8 facing)` is called from **one site image-wide** (`EnumRefs callto:0054d210` → 1 hit / 1 owner / 0 orphan, `FUN_00549990` @`00549a5b`), at the **start of a cell transit**, so the rate is frozen for the whole transit and re-derived at the next cell. `dir = ((facing + 0x10) >> 5) & 7` (`0054d213`…`0054d21b`); `dst = src + (i16)world[0x58ec0 + 4·dir]` (`0054d23c`). Then `domain = actor->vt+0x20()` (`0054d249`) forks. **Ground (`domain == 1`)**: `speed` = `[[actor+0x70]+0x3c]+0x44` **zero-extended** when that byte is nonzero (`0054d2db`…`0054d2e2`, the group term, `MOVE-GROUP-030`), else `(i16)actor+0x8c` (`0054d2e9`, the class `Speed` column); `d = clamp((i8)(height[src] − height[dst]), −32, +32)` over `world+0x9451c` (`0054d2ac`…`0054d2d2`); `v = SpeedMultiplier · speed` with `SpeedMultiplier = world+0x58db4` (`0054d2f0`/`0054d2f6`, `MOVE-PARAM-006`, ships 8 in **both** roots); `v += (v·d) >> 6` **arithmetic** shift, uphill (`d < 0`) reducing (`0054d305`…`0054d31a`); `c = ((u8)(cost[src] + cost[dst])) >> 1`, a **byte** add that wraps at 256, and `if c == 0 then c = 8` (`0054d334`…`0054d342`); `v = v / c` **signed** (`0054d358`); `v = clamp(v, 1, 63)` (`0054d35c`…`0054d36b`). **Non-ground (`domain != 1`)**: `v = (speed·8)/8` — the identity — from the **same two sources in the same order** (`0054d25c`…`0054d296`), with **no** multiplier, **no** slope, **no** cost lookup, then the same clamp. Consequence a consumer must carry: the two arms **agree exactly** when `SpeedMultiplier == 8` and mean cost is 8, which is what ships, so no corpus can separate them; move either and only the ground movers change. Then `mover+0xa8 = v`, `mover+0xae = dir` (`0054d37c`/`0054d389`), the per-axis steps by `MOVE-DIR-034`, and `mover+0xaa = ceil(256 / step)` (`0054d441`…`0054d468`). `cost(cell)` is `FUN_0054e5e0`, reading the plane at **`world+0x00000`** and **not a pure read** — with `block[cell] & 0x20` and a nonzero byte at the cell's `world+0x540b4` record `+0xe` it shifts the stored byte right by 2 and writes it back (`0054e654`/`0054e657`)

**Confidence.** High (every term, its address and its arithmetic form — `SAR` against `SHR`, the byte-width add, the signed `IDIV`, the zero- against sign-extension of the two speed sources — is transcribed from the raw listing in `evidence/listings.md` §1–2, and the single call site is an `EnumRefs callto:` with 0 orphan)

**Original status.** ● active

**Evidence.** [EXP-0093](../experiments/EXP-0093-move-rate/)

### MOVE-GROUP-030

**The speed override `TERR-MOVE-058` left Unknown is a GROUP term, and it is the minimum `Speed` over the group's members.** `actor+0x70` is the AI **group** (`AI-GROUP-009`) and `grp+0x3c` its `0x50`-byte AI record, so `FUN_0054d210`'s `[[actor+0x70]+0x3c]+0x44` is `grpAI+0x44` — the byte `AI-MOVE-023` already reads as a flag at `005376bb`. Both group-order setters compute it the same way: a local initialised to **`0xfa`** (`005340c0`, `005343b0`), then a per-member loop `CMP word ptr [ESI + 0x8c],AX ; JGE ; MOV CL,byte ptr [ESI + 0x8c] ; MOV byte ptr [ESP + 0x13],CL` (`005342fa`…`0053430f`, and `005345ea`…`005345ff`), then `grpAI+0x44 = min` (`00534353`, `00534643`) beside `grpAI+0x20 = 4`, the Move order. It is cleared to **0** together with the order code when the group's per-member pass ends (`FUN_005355c0` @`00535694`/`0053569c`). So **a unit under a group order moves at its slowest companion's rate and reverts to its own class `Speed` when the order clears**, and because the store is `AND EAX,0xff` while the class fallback is `MOVSX`, the group term is **unsigned** where the class term is signed. Instrument for the writer set: `EnumRefs disp:44`, 74 byte-wide hits image-wide, of which exactly four in the AI module `0x53xxxx` touch `[reg+0x44]` as a byte — two writes, one clear, one compare; blind, as every `disp:` sweep is, to a wholesale structure copy

**Confidence.** High (the initialisation, the comparison, the store, the clear and the two call sites are named instructions, and the `min` shape is the loop's own arithmetic) / **Medium** (that a *running* game reaches this store on any particular order: the gate at `[ESP+0x18]` is derived from `[[grp+0x44]+0x30]+0x1f == 2` and a spread test against `world+0xa824`, neither identified — see the write-up's open questions) / ~~**Unknown** (whether any path writes `grpAI+0x44` while `grpAI+0x20` is not a move order)~~ — **both close, and one clause of this row falls, by [EXP-0094]**. The gate is a single formation flag: the two-indirection comparison is the owning `Player`'s formation mode `[[grp+0x44]+0x30]+0x1f`, default **2** (`AI-FORM-037`), and the spread test is a **Chebyshev distance in whole cells to the group centroid** against `AImanager+0xa824` — a compile-time **2** on the **AI manager**, not `world+0xa824`, which is the wrong object (`AI-SPREAD-038`, retracted.md). The store fires **iff** the group moves in formation (`MOVE-GATE-035`), and no path writes the byte outside a group Move/Swarm-2 order except the record's own constructor and its raw `CArchive` (de)serialization (`MOVE-GROUP-037`). **And the clear is retracted**: `FUN_005355c0` has no caller by any of four instruments (`AI-DEAD-036`), so a unit does **not** revert to its own `Speed` when the order ends

**Original status.** ● active (amended)

**Evidence.** [EXP-0093](../experiments/EXP-0093-move-rate/), **[EXP-0094](../experiments/EXP-0094-group-rate/)**

**Amended.** The group-rate clear and world-field ownership clauses are partially retracted. MOVE-GROUP-037 and the correction in retracted.md state the retained field and AI-manager owner.

### MOVE-TURN-031

**Stepping requires matching facing; route destruction and final default0x10 were false interpretations.** `FUN_00549990` derives next-cell facing in 32-unit steps and tests current against desired before its rate/step branch. `0054a296 CALL 00548e20` advances current facing without touching a route; the short-arc snap requires the full mover+a0 dword to be zero. MOVE-TURN-044 gives the local turn update. The table destination binding stands: Units slot9 at004f5704 and Humans slot7 at004f9817 address mover+a. The old clause that67 unparameterized rows keep0x10 as the resulting actor rate is partially retracted. The mover constructor writes 16, then Unit defaults write8 to the same allocated object; Human derive later copies low8(actor+8c). MOVE-RATE-052 and MOVE-RATE-053 distinguish these producers. Earlier raw table statistics were not remeasured by these corrections and do not establish final live rates.

**Confidence.** **High** for the named local instructions and table destinations. Whole-action elapsed time, final corpus-wide actor rates and first post-LOAD scheduling are not established.

**Original status.** ● partially retracted

**Evidence.** [EXP-0093](../experiments/EXP-0093-move-rate/), [EXP-0324](../experiments/EXP-0324-turn-continuation/), [EXP-0329](../experiments/EXP-0329-mover-rate-binding/); [retraction](retracted.md)

**Amended.** Route destruction, unconditional snap, elapsed-time and final-default-rate clauses are partially retracted. MOVE-TURN-044 and MOVE-RATE-052/MOVE-RATE-053 state the narrower relations; retracted.md records their scope.

### MOVE-CLOCK-032

**The displacement is one step per actor per SUB-TICK, it is counted in ticks and never measured in milliseconds, and the presentation tick that drives animation is issued from the same loop iteration.** `FUN_00548c60` has **2 call sites image-wide** (`EnumRefs callto:00548c60`, 0 orphan): `FUN_005495f0` @`00549605` and `FUN_00549990` @`00549a72`. Every path through the actor update reaches exactly one of them — `FUN_00548f70` calls `FUN_005495f0` and **returns immediately** when either fraction is not `0x80` (`00548fcd`…`00548fe4`), and otherwise reaches `FUN_00549990` once at its tail (`0054928d`); `FUN_005310e0`'s arms are a jump table. Above that the chain is single-caller at every step: `FUN_0050fca9` ← `FUN_004d891a` (1 hit) ← `FUN_004d2551` (`004d25df`, unconditional), and `FUN_004d891a` increments `server+0x04`, the sub-tick counter. **So one actor advances by `mover+0xb0`/`+0xb1` exactly once per sub-tick**, and the arrival test is `mover+0xac >= mover+0xaa` — a tick count against a tick count (`00548deb`…`00548df9`), with **no elapsed-time term anywhere in the routine that writes the position**. Three consequences, all of which a consumer gets wrong by modelling seconds: cells **per tick** is invariant; cells **per second** is `1000 / (T · campaign+0x3f0)` and therefore moves with the `[GameOptions] Speed` index; frame rate cannot enter, because `FUN_004753c0` is a **deadline loop that repeats** rather than a scaler (`00475461 JBE`); and `campaign+0x3dc & 1` clear stops **both** the simulation and the `0x401` presentation tick (`00471665`/`00471667`), so a paused unit does not drift and its animation does not either. On the normal play arm the simulation sub-tick (`00475415`) and the `0x401` tick (`00475431`) are two calls of **one** iteration in a fixed order, simulation first — they are not two clocks but two consumers of one pacer (`SESS-PACE-018` bounds where that stops being true)

**Confidence.** High (the two call sites, the single-caller chain, the tick-against-tick arrival test and the two calls' adjacency in one loop body are named instructions and `callto:` enumerations with 0 orphan) / **Medium** (that no *other* arm of the engine advances an actor twice in one pacer iteration: `FUN_004d891a`'s second caller `FUN_004d214b` does run it sixteen times in a loop, and is a profiler reached from `FUN_004d0c2e`, not from the idle handler — read, but its own caller set was not enumerated)

**Original status.** ● active

**Evidence.** [EXP-0093](../experiments/EXP-0093-move-rate/)

### MOVE-LIMIT-033

**The customisation limits of the rate (goal G2), by complete enumeration of the law's input domain rather than by sample.** (a) **The clamp `[1,63]`** (`0054d361`, `0054d36b`) is a hard-coded pair of immediates; at `SpeedMultiplier` 8 and mean cost 8 the ceiling first bites at `Speed = 63`, against a shipped maximum of 35, so it is a guard on shipped data and a wall on authored data. (b) **The transit quantisation is the severe one.** `mover+0xaa` is a whole number of sub-ticks and the fractions are snapped to the cell centre when it is reached, so the surplus travel is discarded and every `v` in a class crosses a cell in the same time: `v ∈ [1,63]` yields only **27** distinct transit times and the largest class is **12** values wide (`v` 52..63 all take 5 ticks). Against the shipped `Speed` column that is **22 distinct values → 15 distinct transit times at cost 8, 13 at cost 6, 11 at cost 16**. (c) **The slope term is a `>>6`**, so a height delta under 4 changes nothing at `Speed` 16. (d) **A diagonal is not `√2`**: the `0.707` truncation and the round-up together leave the diagonal transit between **0.97×** and **1.15×** of `√2 ×` the straight one over the shipped speeds. Which of these move a shipped file's bytes if lifted: `Speed` and `RotationSpeed` **yes** (`world.res:data/data.bin`), `SpeedMultiplier` and the `Cost*` alphabet **yes** (`world.res:data/map.reg`); the clamp, the `>>6`, the `256` grid, the `0.707` and the nine-arm tick ladder **no** — none of them is carried by any shipped file. Shipped alphabets, identical in the EN and RU roots (0 differing rows over 333): `Speed` 22 distinct values in 8..35 (mode 19), `RotationSpeed` 16 distinct in 8..23 (mode 19), 67 rows unparameterised in both, `movementType` `−1:277 1:40 2:8 3:8`

**Confidence.** High (the collapse classes are a **complete enumeration** of the law's domain by a re-implementation whose every step cites its instruction, `tools/movespeed -mode limits`; the byte-moves-or-not column is a statement about which file carries each value, and each is named) / **Medium** (the shipped alphabets: two roots, one version of each, under EXP-0049's `data.bin` grammar and EXP-0032's `.reg` framing)

**Original status.** ● active

**Evidence.** [EXP-0093](../experiments/EXP-0093-move-rate/)

### MOVE-DIR-034

**The eight directions, the diagonal constant, and the per-axis step — the last three things a consumer needs to compute the next position.** `world+0x58eb0` and `world+0x58eb8` are two 8-byte tables written once by the world init `FUN_00547eb0` (`00547f86`…`00548003`; `EnumRefs disp:58eb0` 4 hits / 1 write, `disp:58eb8` 5 hits / 1 write, 0 orphan): `dx = [0,+1,+1,+1,0,−1,−1,−1]`, `dy = [−1,−1,0,+1,+1,+1,0,−1]` — clockwise from north, all in `{−1,0,+1}`. The packed-cell delta table at `world+0x58ec0` is built from those two in an 8-iteration loop as `(dy << 8) + dx` (`0054800f`…`0054802e`), matching `MOVE-STEP-010`'s `(y<<8)\|x` packing. The step stores: **straight** (`dx·dy == 0`) is an 8-bit `IMUL` of `v` by the component (`0054d3f1`/`0054d3f3`), so the step is `±v`; **diagonal** is `trunc(v · component · K)` through `FILD` / `FMUL double ptr [0x0059cd98]` / the CRT `ftol` at `0x0055458c` (`0054d3b9`…`0054d3c3`). **`K` is `0.707` exactly, not `√2/2`** — the eight bytes at `0x0059cd98`, read out of `rom.exe` through the PE section table without a disassembler, are `39 b4 c8 76 be 9f e6 3f` = `0.70699999999999996` (`evidence/diag-constant.txt`)

**Confidence.** High (the sixteen immediate stores, the delta loop, both step forms and the `FMUL` operand are named instructions; the constant's value is the shipped file's own bytes)

**Original status.** ● active

**Evidence.** [EXP-0093](../experiments/EXP-0093-move-rate/)

### MOVE-GATE-035

**The group rate term applies exactly when the group moves in FORMATION — one local flag decides both, and it is governed by two gates neither of which is data.** `FUN_005340a0`/`FUN_00534390` seed `[ESP+0x18] = 1` (`005340b8`, `005343a8`) and then run **gate 2**, the per-player formation mode `[[grp+0x44]+0x30]+0x1f` (`AI-FORM-037`): `== 2` takes the conditional arm, `!= 2` sets the flag **to the mode byte itself** and skips the spread test entirely (`005341df`…`005341e8`, `005344cf`…`005344d8`), so `0` disables formation and any other nonzero value forces it. On the conditional arm the routine computes the group centroid with `FUN_00533210` — the members' 1/256-cell positions summed, divided **unsigned** by `grp+0x0c`, high byte of each taken — and per member runs **gate 1**: `FUN_0054d100(memberCell, centreCell)` is `max(\|dx\|,\|dy\|)`, Chebyshev in whole cells, and `CMP AL,byte ptr [ECX + 0xa824]` / `JBE` clears the flag for the whole group the moment one member exceeds the threshold (`005341a0`/`005341a6`, `00534490`/`00534496`; the threshold is `AI-SPREAD-038`). The **same** flag is then tested twice: at `00534296` it forks the per-member loop (`MOVE-FORM-036`) and at `00534346` it gates the rate store `00534353` — and the running minimum is updated **only inside the formation arm**, because the plain arm jumps over it (`0053437e` → `0053438e` → `00534315`). Consequence for a consumer: `MOVE-RATE-029`'s group speed source is live iff the last group Move/Swarm-2 order this group received was issued in formation

**Confidence.** High (both gates, the fork and the store gate are cited instructions in two routines read end to end; the mode byte's identity rests on an `EnumRefs "re:.*\+ 0x1f\].*"` sweep of **8 hits / 8 owners / 0 orphan** image-wide, four of them stack, and the threshold on a `disp:a824` sweep of **3 hits / 3 owners / 0 orphan**)

**Original status.** ● active

**Evidence.** [EXP-0094](../experiments/EXP-0094-group-rate/)

### MOVE-FORM-036

**A group move DOES distribute destinations: there is a formation, there is a per-member offset table, and there is a spread test — `MOVE-ORDER-023`'s headline is refuted, through the blind spot that row's own confidence cell named.** On the formation arm each member's order block records `ord+0x24 = memberCellX − centroidCellX` (`005342b9`, `005345a9`) and `ord+0x26 = memberCellY − centroidCellY` (`005342d6`, `005345c6`) as `i16`, then the routine reads them back **as bytes** and issues `FUN_0052f4f0(member, X + dx, Y + dy)` (`005342e4`…`005342f5`, `005345d4`…`005345e5`); on the plain arm every member gets the unmodified `(X, Y)` (`0053437e`, `0053466e`). `FUN_0052f4f0` clamps the destination into the playable rectangle `world+0x58ee0..+0x58ee3`, writes `ord+0x0a = (y<<8)\|x`, `actor+0x50 = 1` and `ord+0x08 = 1` — **and that last store is `MOV byte ptr [EDX + 0x8],AL` at `0052f5eb` with `EAX` loaded `1` at `0052f5ca`**, i.e. the move opcode arrives in a register, which is exactly the encoding `MOVE-ORDER-023`'s `MOV byte ptr [reg + 0x8],0x1` enumeration states it cannot see. The offsets are scratch: `EnumRefs disp:26` is 11 hits image-wide and exactly **4** lie in `0x0052c000..0x0055a000`, the two stores and the two byte reads inside these two routines, so nothing re-forms a group later

**Confidence.** High (both arms and the clamp are full listings of routines read end to end; the refutation is an instruction the earlier row's instrument is documented as excluding, not a disagreement about what a listing says)

**Original status.** ● active

**Evidence.** [EXP-0094](../experiments/EXP-0094-group-rate/)

### MOVE-GROUP-037

**The complete writer set of `grpAI+0x44`, and the finding that NOTHING clears it — so a group keeps a rate it can no longer justify.** Instruction-level writers, from `EnumRefs disp:44` over the whole image (**630 hits / 270 owners / 1 orphan**, of which 48 are byte-wide on a register base and exactly **6** have a `[grp+0x3c]` base): `00534353` and `00534643`, the two gated stores (`MOVE-GATE-035`); `00535694`, the only clear — **and its routine `FUN_005355c0` is unreachable** (`AI-DEAD-036`). The remaining two hits are reads: `005376bb` in the group order-4 arm and `0054d262`/`0054d2db` in the rate routine. Two further writers exist that **no displacement sweep can see**, and they are named rather than left to the blind spot: the record's constructor `FUN_0052cc40` zeroes the whole `0x50` bytes (`0052cc45`…`0052cc50`, `REP STOSD` × `0x14`), and `FUN_005391d0`, the record's `Serialize`, is a **raw `CArchive::Write`/`::Read` of `this` for `0x50` bytes** (`005391fa`/`00539250`) — so the group rate is in every save carrying a group, and a load writes it wholesale. Lifecycle, each hop a cited routine: a **player** order always allocates a new group (`FUN_0050fa33` → `new(0x48)` → `FUN_0050f81b` → `FUN_0052cc40`), so it starts at 0; a **scenario** group is the authored one and persists, and guard / aggressive / stand-ground / patrol / roam write only `grpAI+0x20`; a member dying or being re-ordered goes through `FUN_0050f996`, which unlinks it and zeroes `actor+0x70` and touches nothing else, **so the group keeps the departed member's `Speed`**. The order ending does not clear it either — `FUN_00537510`'s only use of the byte is as a boolean guarding `ord+0x6c = 1`, and over the image `disp:6c` (496 hits / 248 owners) puts **three** hits in the simulation modules, all of them stores

**Confidence.** High (the writer and reader sets are printed enumerations with their instrument and both of that instrument's blind spots discharged by name; the unreachability is `AI-DEAD-036`) / **Medium** (that a *running* game ever leaves a stale value in play: the persistence argument is a reading of which routines write the byte, not an observation of a session)

**Original status.** ● active

**Evidence.** [EXP-0094](../experiments/EXP-0094-group-rate/)

### MOVE-AREA-038

**The complete movement-side reader set for both block planes, and the one channel through which an area effect reaches it.** Instrument: `EnumRefs disp:10000` (78 hits / 33 owners) and `disp:20000` (60 / 19), 0 hits in orphan or undisassembled code, re-run here and agreeing with `MOVE-DOM-027`'s counts. Every plane read in the movement path is a `TEST byte ptr [cell + plane],reg` whose register is `mover+0x5` reloaded per cell, in six routines: `FUN_0054bd20` and `FUN_00543060` on the static plane, `FUN_0054bf10`, `FUN_0054c100`, `FUN_0054c2f0`, `FUN_0054c6f0` and `FUN_00541dd0` on the dynamic plane, plus the copies inlined in the driver `FUN_00541e80` (`MOVE-SEARCH-001`). **No movement routine reads an area-effect layer slot, a cell-record occupant slot, or a plane bit by immediate.** The mask population is `FUN_0054b120`'s three stores `0x41`, `0x44`, `0x82` (`MOVE-DOM-025`), so the only plane bits that can ever block are 0, 1, 2, 6 and 7. An area effect therefore reaches the search through exactly two channels, both of them writes performed by `FUN_005456d0` when it recomputes a cell from its record (`MAGIC-WALLBLOCK-045`, `MAGIC-AREACOST-046`): the `wall_of_earth` layer slot `payload+0x20` sets bits 0 and 2 on **both** planes, and any occupied layer slot multiplies the cell's cost byte by 4, which `MOVE-COST-002`'s `movementType == 1` arm reads inline at the destination cell. The area module's own two plane writes reach neither: bit 5 (`0054e8a2`) and bit 4 (`0054e09e`, `0054e0a9`) are in no mask, and bit 4's only consuming reader is a serializer (`TERR-PASS-148`). **`FUN_00541dd0` is a standalone footprint query on the dynamic plane** — instruction-for-instruction `FUN_00543060` apart from the displacement, per `MOVE-DOM-027` — with 5 call sites in 2 owners, `FUN_0052e6d0` and `FUN_0054e220`; it is not part of the search

**Confidence.** High for the reader set as an exhaustive read of six listings, and for the two channels, each a named instruction in `FUN_005456d0` / Medium for the enumeration's completeness, whose blind spot is the one `TERR-PASS-073` names: a `disp:` sweep cannot see an access through a pointer the routine `LEA`'d first, and the two address-of-cell getters `FUN_00523be0` and `FUN_005453e0` had their callers enumerated (1 and 2, neither a mover) rather than being covered by the sweep

**Original status.** ● active

**Evidence.** [EXP-0174](../experiments/EXP-0174-area-movement/)

### MOVE-GATE-039

**(rom.exe) No actor is displaced except from inside one call of the order machine, so an actor the machine refuses cannot move at all.** `MOVE-CLOCK-032` names `FUN_00548c60` as the routine that writes the position. Closing the call graph above it with `EnumRefs callto:`, whole image, 0 orphan at every step: `callto:548c60` → 2 hits, `FUN_005495f0` @`00549605` and `FUN_00549990` @`00549a72`; `callto:5495f0` → 3, `FUN_005310e0` @`005311f6`, `FUN_00548f70` @`00548fda`, `FUN_005492a0` @`005492dc`; `callto:549990` → 2, `FUN_00548f70` @`0054928d`, `FUN_005492a0` @`005495e1`; `callto:548f70` → 2, `FUN_005310e0` @`0053129c`, `FUN_00531970` @`005319a3`; `callto:5492a0` → 2, `FUN_00532100` @`0053211b`, `FUN_005310e0` @`005313c0`; `callto:531970` → 1 and `callto:532100` → 2, all of them inside `FUN_005310e0`. `callto:5310e0` → 2, `FUN_004f37be` @`004f398c` and `FUN_00531070` @`0053109b`. `FUN_00531070` is a loop that calls the machine for every element of a list and its own `callto:` is **0 hits**, which `INSTRUMENT.md` rule 7 says is a tell rather than an answer — so every section of the image was scanned for the stored little-endian dword (`tools/effectgate -mode refs`): `0x00531070` appears **0** times, so nothing reaches it; `0x005310e0` appears 0 times, consistent with one direct `CALL`; and `0x004f37be` appears **3** times, at `0x0059c3d8`, `0x0059c460` and `0x0059c4e8`, offset `0x18` in the three actor vtables, which reproduces `MOVE-TICK-009`'s slot 6 from a different instrument. Inside the machine, only one arm above the order switch reaches a displacement — progress arm 3 at `005311ef`, which calls `FUN_005495f0`. Progress arm 4, the refusal `MAGIC-ACTGATE-079` describes, calls nothing. **So a refused actor is not slowed and does not drift: no instruction that writes its position can execute.** The same closure bounds the refusal's one escape (`MAGIC-ACTGATE-079`): `mover+0x98`, the flag that admits the machine's un-gated tail, has exactly three setters — `00549092` in `FUN_00548f70`, `005494a4` in `FUN_005492a0` and `00549b8d` in `FUN_00549a90`, whose two callers are those same two routines (`EnumRefs disp:98`, 215 hits / 140 owners / 4 orphan; `disp:99`, 0 hits) — and all three sit inside executors the closure puts behind an order arm. So while an actor is refused the flag can never be raised again, and the tail's own `00531673` clears it

**Confidence.** **High** (six `callto:` enumerations with 0 orphan hits at every step, and the two results that a `callto:` cannot settle — a zero and an all-`RDATA-SLOT` set — settled by a raw scan of every section for the stored dword, which is the instrument that can see a pointer table) / **Medium** that no *other* routine writes an actor's position without going through `FUN_00548c60`: that identification is `MOVE-CLOCK-032`'s and was cited rather than re-derived

**Original status.** ● active

**Evidence.** [EXP-0179](../experiments/EXP-0179-effect-action/)

### MOVE-STEP-040

**(rom.exe) The refusal begins on a cell centre, so a step already in flight completes.** `FUN_005310e0` writes the parking value only when the order object's progress byte is already 0: `00531167 MOV CL,byte ptr [EAX + 0x9]` / `0053116a TEST CL,CL` / `0053116c JNZ 0x00531171` (rel8 `0x03`). `AI-ORDER-039` establishes that progress 0 is the tick the actor stands on a cell centre, `FUN_00545c30` testing `pos+0x4 == pos+0x5 == 0x80`. An actor in transit is in progress state 3, whose arm at `005311ef` calls `FUN_005495f0` — the one arm above the order switch that still advances a step — sets `actor+0x54 = 1`, and clears the progress byte when `FUN_00545c30` reports arrival (`005311fb`…`0053121c`). The state is entered by the walk arm itself at `005312b7` and by `005313db`, `005319ba` and `00532132`, each immediately after the same `FUN_00545c30` test fails. So the sequence is: the effect attaches, the actor finishes the cell it is entering, the progress byte returns to 0, and the gate fires on the following tick. **The observable is a stop aligned to the grid, never a stop between cells**

**Confidence.** **High** (the order of the two tests inside one routine, the arm's own call and its clearing condition are named instructions, byte-asserted at their own addresses; the cell-centre meaning of progress 0 is `AI-ORDER-039`'s, cited)

**Original status.** ● active

**Evidence.** [EXP-0179](../experiments/EXP-0179-effect-action/)

### MOVE-EFFLIST-041

**(rom.exe) No routine in the movement path reads the drawable's effect list, so the list is not a second immobilisation mechanism.** `MOVE-GATE-039` and `MAGIC-ACTGATE-079` place the only refusal in the order machine, on `actor+0x144`. The effect list at `drawable+0x124`/`+0x128`/`+0x12c` was checked separately. It has one reader, `FUN_004599d0`, whose complete call-site set is fourteen addresses over six routines, every one of them between `0x0040a13c` and `0x004619ab` in presentation code (`MAGIC-EFFLOOK-083`). A displacement sweep of the list's own three fields over the whole image, `EnumRefs disp:124 disp:128 disp:12c imm:124`, returns 71, 78, 130 and 10 hits. Every hit in the `0x0053xxxx` and `0x0054xxxx` family -- the order machine, the movers, the AI -- is a **byte** access at `+0x12c` on the simulation actor, which is the reach field `FUN_005327d0` compares a distance against at `00532a07`, while `FUN_004599d0` reads `drawable+0x12c` as a **dword**. The `imm:124` hits are three stack frames, two `LEA`s at `+0x1244` and `+0x124c`, and five `ADD reg,0x124` of which two are `FUN_004104e8`'s own message arms. No simulation-family access to the list was found.

**Confidence.** **High**, scoped to what was searched: two instruments with different blind spots for the reader's call sites (a reference index on the repaired function table, and a raw `E8 rel32` scan of every executable section) agreeing address for address, and a displacement sweep over the whole image reduced by access width and by base class rather than by address family alone. Blind spot: a read of the list through a pointer held in a register with no displacement, or through a computed call, is invisible to both

**Original status.** ● active

**Evidence.** [EXP-0182](../experiments/EXP-0182-effect-at-actor/)

### MOVE-072

**`mover+0x90` is not "route list non-empty" — it is a boolean "actor already stands on the route's own resolved endpoint" flag. Image-wide, `1` is written only at `0054910c`, inside `FUN_00548f70`'s own walk-arm tail; `0` is written at eight further sites across seven routines, `FUN_00548f70` itself included — not "only in `FUN_00548f70`," as this row first had it.** A whole-image write-site census (`evidence/mover90-write-census.json`, this round) finds every store to `[reg+0x90]` with no index register across the whole `.text` section by linear-sweep capstone decode — sequential decode with no control-flow following, whose own blind spot is a byte sequence that is really embedded data, or a misaligned instruction tail, decoding as a spurious hit — then keeps the sites whose base register was loaded from `[x+0x154]` (the mover-pointer fetch) within the 25 preceding decoded instructions. It finds 150 raw `disp+0x90` stores and nine through a mover pointer: `0054910c` (value 1) and `00549061` (value 0), both inside `FUN_00548f70`, read whole this experiment; `00537612` and `0053769d` (both value 0), inside `FUN_00537510`, the Move arm's own per-member routine, also read whole this experiment; and five sites outside this experiment's nine-routine scope, each storing 0 — `0052f5dd` (`FUN_0052f4f0`), `0052f6b1` (`FUN_0052f600`), `00530951` (`FUN_005308a0`), `00530a06` (`FUN_00530970`), `005358cb` (`FUN_005357a0`). Attribution for those five follows `Image.body()` recursive descent seeded from those five entries, which independently confirms each site falls inside that routine's own reachable body; this experiment did not read any of the five whole, so nothing beyond "this routine writes 0 here" is claimed for them. `FUN_00548f70`'s own two writers work as before: `00549061 MOV dword ptr [EAX+0x90],EBX` (`EBX=0`) clears it unconditionally right after the static search `FUN_00541e80` returns (`00549056`), regardless of that search's own result; `0054910c MOV dword ptr [EAX+0x90],1` sets it, later in the same tail, only when the actor's own current cell (`[actor+0x10]+2`) equals `mover+0x76` (`00549106 CMP CX,word ptr [EAX+0x76]` / `00549108 JNE`) — and only on that branch does the routine return immediately (`0054910d`..`0054911b`), skipping the rest of the tail, including the call into the stepper `FUN_00549990`. `mover+0x76` itself, set earlier in the same tail (`00549088`/`005490ae`), takes the cell word from the search result node at `[actor+0x164]+8` when the static search's own count field `[actor+0x168]` is nonzero, or collapses to the actor's own current cell when it is zero — the same branch that sets `mover+0x98` (`AI-ROUTE-045`'s route-search-failure flag) to 1 (`00549092`). So `mover+0x90 = 1` fires in exactly two circumstances, both inside `FUN_00548f70`: a genuine "already there," or the degenerate case where a total search failure collapsed `mover+0x76` onto the actor's own cell. The gloss `AI-MOVE-023` gives the same flag, at its own read of `FUN_00537510`'s `0053761d CMP dword ptr [EAX+0x90],EBX` branch, is corrected the same way — see `claims/retracted.md`, whose row is extended this round to also cover that routine's own arrival-branch gloss, "route list `pth+0x90` cleared": the two 0-writes this census finds inside `FUN_00537510` (`00537612`, `0053769d`) clear `mover+0x90`, not a route list, either

**Confidence.** High for the two `FUN_00548f70` writers and their conditions (cited instructions inside one routine read whole this round) and for the census's own site-and-value findings (a deterministic scan, independently reproduced from `regen.sh`) / Medium for the routine attribution of the five sites outside this experiment's own nine-routine scope, which rests on recursive descent from entries the census script did not itself discover, not on a fresh whole-body read of those five routines

**Original status.** ● active

**Evidence.** [EXP-0386](../experiments/EXP-0386-engagement-under-occupancy/EXP-0386.md)

### MOVE-073

**`FUN_005492a0`, the pursuit order's own walk/search routine, has no `MOVE-WAIT-008`-style wait-and-turn branch anywhere in its 247 instructions — but it is not a universal "every path reaches the stepper" either: three named exits return without calling it, and one of the routines it calls can raise `mover+0x98` against an occupied approach through a mechanism this row did not previously read.** The three exits: not centred (`005492d2`/`005492d7` false, `005492d9` calls `FUN_005495f0` and returns at `005492e8`); already within stop distance (`005492ef ja 0x54930d` taken, turns via `FUN_0054a680`/`FUN_0054a210` and returns at `0054930a`); and the routine's own reused-route match, where the current cell already equals both `mover+0x76` and `mover+0x8c` (`0054954a`/`00549550`, both `jne` not taken), which increments the progress counter `mover+0x09`, clears `mover+0x7c`, and returns at `00549571` without calling the stepper. Every other path — a refreshed dynamic route, or the static search running at all — does reach the stepper `FUN_00549990` unconditionally, at `005495e1`. The static search this routine runs is at `00549450`, not `00549056` as this row first had it (`00549056` is `FUN_00548f70`'s own call to the same search); `MOVE-PLANE-005`'s own unit-blind static plane means a destination cell occupied only by another actor does not fail that search — it still returns a route ending at the literal requested cell — and this routine's own direct setter (`005494a4`) fires only on a **total** static-search failure (`[esi+0x168] == 0` after the call at `00549450`). But this routine also calls `FUN_00549a90` at `0054959c`, passing the target actor as `altTarget` — Picker B's own contact-ring search (`MOVE-ALT-018`/`MOVE-ALT-020`), read whole this round for `AI-335`'s own correction — and that routine's own `mover+0x98 = 1` store at `00549b8d` (`MOVE-GATE-039`'s third setter) is reachable from here too, under the same two-part gate `AI-335` names: the inner dynamic-plane search against the goal comes back empty, and the branch that picked that goal aimed straight at the final destination. **So "never on mere occupancy" is unsupported**, and what an every-approach-occupied pursuit actually produces is Unknown — narrower than either "always waits like an ordinary walk" or a flat "fails" — because no reading traces Picker B's own contact-ring search for the specific case where every ring cell is itself occupied. `FUN_00532100`, the out-of-position wrapper `AI-PURSUE-040` already reads, confirms the wait-and-turn absence at its own call site: it calls `FUN_005492a0` once (`0053211b`), tests only whether the actor ended up centered (`FUN_00545c30`, `00532123`), and on failure writes `ord+0x09 = 3` at `00532132` — entry into progress state 3, the in-transit step-along-the-path arm `MOVE-STEP-040`/`AI-ORDER-039` already name, not a retry marker — before unconditionally setting `actor+0x54 = 1` on every path

**Confidence.** High for the structural exit enumeration (both routines read whole this round — 247 and 21 instructions, 0 unresolved — every exit from `FUN_005492a0`'s tail traced by hand) and for the corrected static-call address and progress-state gloss / Medium for what an occupied-but-passable approach ordinarily produces through the unconditional stepper call, which composes `MOVE-PLANE-005`/`TERR-PASS-051`'s already-published occupancy-blindness with this round's caller rather than re-deriving the passability test itself / Unknown, explicitly, for the every-approach-occupied scenario specifically: it depends on Picker B's own contact-ring search, which this experiment did not trace for that case

**Original status.** ● active

**Evidence.** [EXP-0386](../experiments/EXP-0386-engagement-under-occupancy/EXP-0386.md)


## Centered turn continuation

| ID | Claim | Confidence | Status | Evidence |
|---|---|---|---|---|
| MOVE-TURN-044 | The turn call advances orientation; its byte estimate is recomputed, not consumed as a countdown. | High | ✔ promoted | [EXP-0324](../experiments/EXP-0324-turn-continuation/), complete instruction listings and labelled derived transition examples |
| MOVE-RATE-052 | Mover construction, actor defaults and table binding are separate producers. | High | ✔ promoted | [EXP-0329](../experiments/EXP-0329-mover-rate-binding/), instruction, binding and vector measurements |
| MOVE-RATE-053 | Human derive replaces the mover byte with the low byte of the current actor speed word. | High | ✔ promoted | [EXP-0329](../experiments/EXP-0329-mover-rate-binding/), complete derive listing and discriminating transfer vectors |
| MOVE-RATE-054 | Effect selector18 adjusts the target's mover byte before invoking that same target's derive. | High | ✔ promoted | [EXP-0329](../experiments/EXP-0329-mover-rate-binding/), original table values and effect-to-call vectors |
| MOVE-RATE-055 | The reached byte consumer changes facing and an estimate; its actor comes from an argument, not incoming ECX. | High | ✔ promoted | [EXP-0329](../experiments/EXP-0329-mover-rate-binding/), all-file-backed-section references and bounded instruction vectors |

### MOVE-TURN-044

**The turn call advances orientation; its byte estimate is recomputed, not consumed as a countdown.** `0054a210..0054a332` writes desired facing to mover+1 and computes the pre-step shorter arc. When the DWORD at+a0 is zero it clears BYTE+9d; only that inactive arm snaps an arc<=32, writing current=desired and BYTE+a4=1. Otherwise it calls the complete leaf `00548e20..00548e95`, which reads BYTE current+0, desired+1 and RotationSpeed+a. The leaf snaps if the shorter arc is strictly less than the rate; otherwise it adds or subtracts the rate modulo256 in the direction reaching desired, taking addition for the exact128 tie. It contains no call or route/list access. The caller writes BYTE+a4=ceil(pre-step shorter arc/rate), sets DWORD+a0=1, increments BYTE+9d modulo256, then clears DWORD+a0 if the updated facings match. Nonzero rate is required on the division arm; zero-rate gameplay validity is Unknown. `00549990..00549a8a` tests facing equality before the rate/position step; its mismatch arm calls0054a210 and returns, even when that turn finishes. The ordinary move caller reaches00549990 at0054928d; the selected order arm calls that move routine at0053129c. This establishes local repeatable progression and a later-call step boundary, not an unconditional next-server-tick promise for every order.

**Confidence.** **High** for the complete local turn/update and next-cell branch; **Medium** for the bounded ordinary-caller interpretation. Full route replanning, callback effects, all order scheduling, malformed/zero-rate lifecycle and post-LOAD first dispatch remain **Unknown**.

**Original status.** ✔ promoted

**Evidence.** [EXP-0324](../experiments/EXP-0324-turn-continuation/), complete instruction listings and labelled derived transition examples

### MOVE-RATE-052

**Mover construction, actor defaults and table binding are separate producers.** Unit initialize004f30a2 requests 180 bytes at004f3110, passes the nonnull allocation to00545c00, then stores that constructor's returned pointer in actor+154 at004f3154. The mover constructor clears180 bytes and sets BYTE+a=16. The same initializer calls004f3317 at004f32ea; its004f3367 loads actor+154 and004f336d writes BYTE+a=8. Unit table004f5604 passes that allocated byte address to00523410 at004f5704 after nine sequential slot advances; Human table004f974d does so at004f9817 after seven. Helper00523410 reads the selected dword, skips a -1 assignment, otherwise stores low8, and always advances the cursor. A missing table value therefore preserves the incoming byte, not an independently chosen16.

**Confidence.** **High** for these complete helpers and selected caller prefixes: both images and 12 table vectors distinguish sentinel, zero, truncation and receiver aliases. Allocation-vector return is synthetic; allocation failure, all construction paths and final spawned rates remain **Unknown**.

**Original status.** ✔ promoted

**Evidence.** [EXP-0329](../experiments/EXP-0329-mover-rate-binding/), instruction, binding and vector measurements

### MOVE-RATE-053

**Human derive replaces the mover byte with the low byte of the current actor speed word.** Derive entry004f7dfc stores its actor receiver atEBP-14. At004f85d7..004f85ec it loads ECX=[actor+154], AL=BYTE[actor+8c], then writes [ECX+a]=AL. This follows its Reaction-based speed stores, type19/21 addition, carried-load penalty and004f54f8 fold of the speed modifier at actor+d8 into WORD+8c. The fold's signed-negative branch clears the modifier, not the already stored negative speed word. A complete derive listing and bounded original-instruction slices distinguish low-byte copy from an independent retained RotationSpeed, saturation or pointer alias. Unit vt+50=004f5946 has no such store; Humanoid/Human vt+50=004f7dfc.

**Confidence.** **High** for the stated local sources, widths and selected vtable slots. The eight composed speed-slice vectors do not execute intervening derive code; arbitrary inputs, final spellbook callback effects, all event reachability and first post-LOAD recomputation remain **Unknown**.

**Original status.** ✔ promoted

**Evidence.** [EXP-0329](../experiments/EXP-0329-mover-rate-binding/), complete derive listing and discriminating transfer vectors

### MOVE-RATE-054

**Effect selector18 adjusts the target's mover byte before invoking that same target's derive.** Original dispatch table00502846 slot18 selects005020dc. Target is argument1 atEBP+8; the arm loads [target+154], adds low8 of the computed effect value to BYTE+a, and stores modulo256. The value is magnitude times signed16(effect+40) when mode+3d meets the original masks1/2/4, otherwise magnitude times DWORD(effect+40), with wrapped32-bit product. At00502833 ECX is the same target and the call is [target.vtable+50]. For the measured Unit slot this selects004f5946; Humanoid/Human select004f7dfc, whose later store is MOVE-RATE-053. The direct effect adjustment is therefore an intermediate store, not proof of a surviving Human rate bonus.

**Confidence.** **High** for selector, arithmetic, exact target and outgoing virtual call; five original-instruction vectors per image end immediately before that call. Item application scheduling, duration/removal, intervening callbacks and full post-effect values remain **Unknown**.

**Original status.** ✔ promoted

**Evidence.** [EXP-0329](../experiments/EXP-0329-mover-rate-binding/), original table values and effect-to-call vectors

### MOVE-RATE-055

**The reached byte consumer changes facing and an estimate; its actor comes from an argument, not incoming ECX.** Leaf00548e20 reads actor fromESP+4 and loads its+154 intoESI before reading BYTE+a at00548e56. It updates only current facing BYTE+0. All three raw E8 callers of this leaf decode in selected windows:005499a3 when actor+184 is zero,00549c93 after a requested-facing comparison, and0054a296 in the turn caller. The first path returns directly after the leaf. The turn caller also reads the same allocated byte at0054a2bb/0054a2d5 as quotient/remainder divisor for BYTE+a4. Selected translation-rate body0054d210 reads actor+8c or the formation override and does not directly read mover+a. Separate rate values7/19 change facing with Position unchanged in the selected no-route-count path.

**Confidence.** **High** for exact receiver and local consumer actions. Thirty-four selected instruction/data windows do not close the whole mover graph: the turn caller has 11 raw direct callers, eight outside the selection; rebased/indirect accesses, complete scheduling and elapsed gameplay time remain **Unknown**.

**Original status.** ✔ promoted

**Evidence.** [EXP-0329](../experiments/EXP-0329-mover-rate-binding/), all-file-backed-section references and bounded instruction vectors


## Restored actor event ordering

| ID | Claim | Confidence | Status | Evidence |
|---|---|---|---|---|
| MOVE-EVENT-060 | The common actor tick admits the selected turn consumer after effects and explicit actor/order gates. | High | ✔ promoted | [EXP-0330](../experiments/EXP-0330-first-mover-event/), selected complete bodies, original class/order tables and relative-transfer checks |
| MOVE-EVENT-061 | A selected attached-effect route can invoke derive on the very actor whose order consumer follows; the derive target is class-specific. | High | ✔ promoted | [EXP-0330](../experiments/EXP-0330-first-mover-event/), Effect/actor slots, constructor assignments and complete selected tick/apply/derive listings |

### MOVE-EVENT-060

**The common actor tick admits the selected turn consumer after effects and explicit actor/order gates.** Original Unit, Humanoid and Human tables all hold004f37be at+18. Manager0050fca9 passes its current actor as ECX at0050fd1b only when actor+54!=16. That actor's tick first dispatches its attached entries through+38 with the same actor as argument1; only afterward does signed WORD(actor+94)>0 reach the order path. Predicate00523360 returns whether actor+3c is zero, so only a nonzero+3c admits004f398c ->005310e0(actor). With the order-machine prefix admitted, order+9=0 and the status hold not installed, its selected table maps order+8=10 to0053154b. Unequal mover+0/+1 passes this same actor to0054a210 at0053156a; equality clears the pending byte instead. The turn's inactive short-arc arm snaps without reading+a; its other arm calls00548e20 and reads the same actor's allocated mover+a for the estimate.

**Confidence.** **High** for the named original slots, receiver aliases, branch table and local control order. The prefix's imported critical-section call, position helpers, earlier effect/dispatcher calls and mutable receiver state remain explicit boundaries. This is conditional reachability, not an assertion that order10 or either turn arm is first after a real LOAD.

**Original status.** ✔ promoted

**Evidence.** [EXP-0330](../experiments/EXP-0330-first-mover-event/), selected complete bodies, original class/order tables and relative-transfer checks

### MOVE-EVENT-061

**A selected attached-effect route can invoke derive on the very actor whose order consumer follows; the derive target is class-specific.** The two measured Effect tables0059c688/0059c6e0 hold0050134f at+38,0050177e at+40 and00501a22 at+48. If an attached receiver has one of these tables, mode+3d includes mask0x02 and its incoming unsigned WORD+42 is a multiple of 8, the tick calls+40 before decrementing duration. For identities other than 8,12,17,0050177e forwards the unchanged actor argument and multiplier1 through+48; generic00501a22 calls that actor's+50 at00502833. Selector0's measured table arm jumps directly to this tail, demonstrating that a direct magnitude store is not required to reach derive. Original archive constructors install Unit/Humanoid/Human tables0059c3c0/0059c448/0059c4d0; their+50 entries are004f5946/004f7dfc/004f7dfc. Thus, if both this generic call and the later order consumer are reached in one actor tick, the same actor's derive call precedes that consumer. The Human/Humanoid derive's reached004f85e9 stores low8(actor+8c) into mover+a; the selected Unit derive body has no such store.

**Confidence.** **High** for the conditional selected-table chain, exact actor and class distinction. An empty effect list, an ineligible pulse, a special identity or another virtual receiver need not take that chain. Actual restored effect population, intervening derive/notification callbacks, normal completion and the first global mover access remain **Unknown**; no native event was observed.

**Original status.** ✔ promoted

**Evidence.** [EXP-0330](../experiments/EXP-0330-first-mover-event/), Effect/actor slots, constructor assignments and complete selected tick/apply/derive listings


`EXP-0386` was allocated ids `72`..`73` of `claims/move.md` (2 ids) and spent both,
`MOVE-072`..`MOVE-073`. **No id is returned unused.** The next free `move.md` id is
therefore `MOVE-074`.

## Seed-cell route boundary

Member offsets, addresses and packed cells below are hexadecimal. Counts,
pool sizes and flag values are decimal.

| ID | Claim | Confidence | Status | Evidence |
|---|---|---|---|---|
| MOVE-080 | Both selected original extractors leave an initially empty list empty when endpoint equals seed; one node is allocated and removed, with count 0 to 1 to 0. | High | ● active | [EXP-0422](../experiments/EXP-0422-seed-route/) |
| MOVE-081 | FUN_00541e80 clears the static list before its seed shortcut; the dynamic entry and both extractors require an explicit input-list premise. | High / Medium | ● active | [EXP-0422](../experiments/EXP-0422-seed-route/) |
| MOVE-082 | In the selected zero-radius FUN_00548f70 static-search tail, substitute=seed and no substitute both take zero count: resolved=current cell, mover+90=1 and mover+98=1. | High / Medium / Unknown | ● active | [EXP-0422](../experiments/EXP-0422-seed-route/) |
| MOVE-083 | A centred zero-radius FUN_00548f70 request equal to the current cell returns before search and preserves route lists and prior mover flags. | High / Unknown | ● active | [EXP-0422](../experiments/EXP-0422-seed-route/) |

### MOVE-080

+4/+8/+c of the static embedded list at actor+15c are first/last/count;
the dynamic list starts at actor+178. Traversal follows node+0 toward the
resolved endpoint. Older movement prose calls these first/last ends tail/head.
The offset and direction relations stand; this card states its naming.

Static FUN_005433a0 calls real NewNode at 005433c1, stores the endpoint word
at 005433e0, and links it at 005433e4..005433f3. Its seed test at 00543406 and
0054340a reaches 00543643, then 0054365b. That tail removes the first node,
updates both end pointers, decrements count at 00543679 and stores it at
0054367d. Zero count reaches full pool teardown at 0054368f..005436a3.
Dynamic FUN_005436b0 has the corresponding operations shifted by 310 hex.
NewNode 0051c590 increments the original embedded count at 0051c643/0051c649;
the pool helpers 00570ba6 and 00570bc6 execute, with only heap malloc/free cut.

Private replay crosses both extractors, empty/old input lists, zero through
three transitions and two memory placements: 32 rows. Initially empty seed
extraction yields first=last=0, count=0, and no stored cells. The ordinary
controls retain exactly the one, two or three non-seed cells. A pre-existing
one-node list retains that old node after seed extraction; extraction does
not clear its input. Allocations are A5-filled, not silently zeroed nodes.

**Confidence.** High for the conditional instruction/list relation. Complete
byte binding, real count/link operations, independent memory walks and
retained-seed/omitted-decrement/backlink loss controls exclude a sentinel
node, a retained seed and count-only normalization. The input list and heap
services are explicit.

**Unknown.** Native allocation failure, corrupt list inputs, and actor
reachability are unmeasured. No universal route result is claimed without
the input-list premise.

### MOVE-081

FUN_00541e80's static arm clears actor+160/+164/+168/+16c/+170 at
00541ecb..00541efe, before coordinate equality at 00541f5c..00541f6f. Equality
returns at 0054304a without calling either extractor. The dynamic arm skips
that static teardown. Its requested-seed shortcut preserves an old dynamic
list. Picker A's nonzero seed answer instead reaches the relevant extractor
through 00542ff2 static or00543045 dynamic. A picker miss clears the static
list through 00542fcc but returns directly on the dynamic arm.

The actual driver/picker/search replay crosses both paths, empty/old lists,
four topologies and two placements: 32 rows. Requested seed, substitute=seed
and no substitute all yield zero nodes from empty input. An ordinary admitted
endpoint 1013 retains 1011,1012,1013 from seed 1010. With an old dynamic cell
1414, those three boundary cases retain 1414; the ordinary route has that
old suffix. Static search clears the old cell in all four cases.

**Confidence.** High for the named local teardown/shortcut/extraction relations.
Medium for reaching the selected topologies under footprint 1, domain 2,
synthetic owner/geometry and declared allocation/notification services. The
code images are identical, one population. Actual helper execution and the
old-list control discriminate an unconditional dynamic-clear model.

**Unknown.** Native caller input-list invariants, notification side effects,
other footprint/domain branches and picker B are unmeasured.

### MOVE-082

FUN_00548f70 calls actual static search at 00549056. It clears mover+90 at
00549061, writes the requested cell to+74 at 0054906d, and tests count at
00549071..00549079. Zero count loads the current Position cell into+76 at
00549088 and sets+98=1 at 00549092. Nonzero count takes the resolved cell
from actor+164's node+8 at 0054909e..005490ae; it does not clear+98.
The following local teardown clears the dynamic list. The comparison at
00549106 then sets+90=1 at 0054910c and returns when+76 equals current cell.

Eight full-entry caller replay rows cross four topologies and two placements.
Both substitute=seed for blocked goal 1011 and no substitute for blocked
goal 3030 leave empty static/dynamic lists, requested goal in+74, seed 1010
in+76, and+90/+98 both 1. Ordinary goal 1013 leaves three static cells,
+76=1013,+90=0 and unchanged prior+98. Its next dynamic-refresh call is
the stopping boundary, not a witnessed downstream action.

**Confidence.** High for local branch use given the count and Position inputs.
Medium for topology reach under the declared private services. The
three-node control and wrong-end loss exclude selection from the first node.
Unknown for native actor/order continuation.

**Unknown.** A zero-count admitted endpoint and true search failure fit the
same caller outputs. The flags alone do not identify semantic search success.
Native reachability, notification effects and the subsequent order result
remain open. AI-335's nonempty-substitute clause is partially retracted.

### MOVE-083

FUN_00548f70 samples current Position cell and the requested cell, then tests
Position fractions for 128/128 at 00548fcd..00548fd5. With equal cells the
computed distance is 0. Range argument 0 reaches 00548ff1's branch directly
to00549292, before the search call and all selected route/flag writes.

The requested-seed full-entry controls preserve supplied static/dynamic old
lists, mover+74=5555,+76=6666,+90=7 and+98=0. The deliberate late-shortcut
loss instead enters search and produces+90/+98 both 1. This distinguishes
a caller-entry shortcut from extraction of a seed endpoint.

**Confidence.** High for the centred zero-radius local entry branch and its
skipped stores. Unknown for native invocation of that input and later order
behaviour. No nonzero-range or mid-transit arm is generalized.

**Unknown.** Native reachability, actor/order cadence, callback effects and
final actor behaviour remain unobserved. A reachable original breakpoint
trace of the same centred zero-radius entry that reaches search or rewrites
these fields refutes this local prediction.

`EXP-0434` was allocated ids `84`..`86` of `claims/move.md` (3 ids) and spent all three,
`MOVE-084`..`MOVE-086`. Its `MAGIC-229`..`MAGIC-234` ids are returned unused because
`claims/magic.md` is still in the table format. The next free `move.md` id is `MOVE-087`.

## Area layers and the cost byte

Addresses, offsets and cell values below are hexadecimal unless a count or a cost is named. Every
quoted instruction is asserted byte for byte on both editions' `rom.exe`, which are identical.

| ID | Claim | Confidence | Status | Evidence |
|---|---|---|---|---|
| MOVE-084 | A transit start reads the cell costs with FUN_0054d210 before any FUN_005456d0 recompute, and every actor tick runs before every area-effect tick; the recompute comes later, at the cell-crossing step. | High / Medium / Unknown | ● active | [EXP-0434](../experiments/EXP-0434-area-cost-order/) |
| MOVE-085 | FUN_0054e5e0 divides a layered cell's cost byte by four once per read, whatever the layer count, Wall of Earth included; two reads with no recompute between them divide twice, and no call-graph edge forbids that. | High / Medium / Unknown | ● active | [EXP-0434](../experiments/EXP-0434-area-cost-order/) |
| MOVE-086 | The cost byte is restored from payload+0 or the CostCracked constant and no cost plane is saved, so a loaded layered cell keeps its ingest value until a recompute; two inline readers divide by it with no zero guard, and decay can zero it. | High / Medium / Unknown | ● active | [EXP-0434](../experiments/EXP-0434-area-cost-order/) |

### MOVE-084

FUN_00549990 (transit start) runs in this order. When the word at movement-record `+0x80` differs from
the target-cell word at `+6`, it calls FUN_0054af40 at `005499fb` (claim clear, skipped when `+0x80`
is 0, `005499f2`) and FUN_0054ad20 at `00549a0e` (claim set), then FUN_00544300 at `00549a32` (heading
for the target cell, stored at record `+1`). When record byte `+0` equals `+1` (`00549a4b`) it calls
FUN_0054d210 at `00549a5b` and then FUN_00548c60 (the step) at `00549a72`. Otherwise it calls
FUN_0054a210 at `00549a80` and reads no cost.

The claim routines write the dynamic plane, not the cost plane: set `0054ae23 OR byte ptr [EAX],0x80`
and `0054af04 OR byte ptr [EAX],0x40`; clear `0054b021 AND CL,0x7f` and `0054b0e6 AND CL,0xbf`. No
direct-call chain from FUN_0054af40, FUN_0054ad20 or FUN_00544300 reaches FUN_005456d0, FUN_0054e5e0,
FUN_0054d210, FUN_0054e730 or FUN_0054e9e0 (`callgraph.txt`). FUN_0054af40 and FUN_0054ad20 each make
two indirect calls on the actor; vtable `+0x1c` and `+0x20` of the three actor vtables are
FUN_00523210 and FUN_00523230, which return actor bytes `+0x49` and `+0x4a` (`0052321a`, `0052323a`).

FUN_0054d210 divides by the mean of two FUN_0054e5e0 reads (`0054d323` source cell, `0054d32f`
destination cell). The mean is `(a+b)>>1` in 8 bits and a result of 0 becomes 8 (`0054d338`..`0054d342`),
so this divide cannot see a zero. The speed term is clamped to 63 (`0054d366`) and the per-tick step
is that term times a table value, times the double `0.707` at `0059cd98` when both table values are
non-zero (`TERR-MOVE-056`). FUN_00548c60 adds the step
bytes (record `+0xb0`, `+0xb1`) to the position and compares the cell bytes (`00548cf3`). Equal cells
jump past every recompute (`00548d30`); a changed cell runs FUN_00545230 (`00548d7e`, release),
FUN_0054abb0 (`00548da8`) and FUN_00544d00 (`00548de0`, occupy); the release and the occupy reach
FUN_005456d0 (`005452ce`, `00545043`, `00545118`). FUN_00549990 is called only by FUN_00548f70
(`0054928d`) and FUN_005492a0 (`005495e1`), and FUN_005492a0 sends any position whose sub-cell bytes
`+4` and `+5` are not both `0x80` to FUN_005495f0 instead (`005492cf`..`005492d7`); FUN_00548f70's
equal test is `MOVE-083`'s. A transit start therefore begins at sub-cell `0x80/0x80`. A step is
the speed term, at most 63, times a direction factor, so it is at most 63 units on an axis (44 on a
diagonal with the `0.707` factor) and cannot change the cell from `0x80`. The crossing recompute follows the cost
read by at least one tick.

FUN_004d891a calls FUN_004d1d86 once per tick (`004d893e`). That routine calls vtable `+0x18` of
each element of the list at `this+0x2c` (`004d1db7`), tests byte `+0x136` after each call
(`004d1dbf`), and only after the list is exhausted calls FUN_00510247 (`004d1deb`), which calls
vtable `+0x18` of each area effect (`005102ac`). The direct-caller closure of FUN_0054d210,
FUN_0054e5e0 and FUN_00548c60 has one root, FUN_004f37be, in `+0x18` of the vtables at `0059c3c0`,
`0059c448` and `0059c4d0` (`0059c3d8`, `0059c460`, `0059c4e8`). The closure of FUN_0054e730 and
FUN_0054e9e0, the layer add and removal, has one root, FUN_004fc9b2, in `+0x18` of the vtable at
`0059c5f0` (`0059c608`). Within it the layers are laid when effect byte `+0x48` is 0 (`004fc9e9`,
`004fcb21 CALL FUN_004fccd8`; `+0x48` is set to 1 at `004fd093` and `004fd1b0`) and removed by
FUN_004fd28d when the word at `+0x4c` reaches 0 (`004fc9ff`, `004fcb10`).

**Confidence.** High for the call order inside FUN_00549990, the dynamic-plane stores of the claim
routines, the getter bodies and the tick order, each a quoted instruction. Medium for the absence of
a direct-call chain from the claim routines to FUN_005456d0, since each makes two indirect calls
that are not followed, and for the closures
and the cell-crossing bound: the instrument is direct `CALL rel32` plus aligned data dwords, so
indirect calls, other `vtable+0x18` call sites, the identity of the `this+0x2c` list's elements with
the actor vtables' objects and the direction tables at `world+0x58eb0` and `world+0x58eb8`, which `TERR-MOVE-056` records as
immediates written by the world constructor, are not re-read here. The sub-cell bound assumes their
entries have magnitude at most 1.

**Unknown.** Observed ticks. An effect added to the area list during the actor walk of the same
tick. On which tick of a transit the crossing falls, which depends on the speed.

### MOVE-085

FUN_0054e5e0 returns the plane byte unchanged unless static bit 5 (`0054e612 TEST AL,0x20`) is set.
For a cell with a record it zero-fills 13 dwords at `ESP+0xc` (`0054e61d`), copies the 13 dwords of the
record there (`0054e644`), reads the byte at `ESP+0xe` (`0054e64a`) and, when it is non-zero, runs
`SHR AL,2` and stores the result to the cost plane (`0054e654`). Scratch `+2` is payload `+2` (the
copy is the record itself). The divide runs once per read and does not use the count's value.

Payload `+2` is written at six sites (`payload2.txt`): zeroed then incremented once per non-null
slot of the six in FUN_0054e730 (`0054e7ae`, `0054e7ba`, second arm `0054e900`, `0054e90c`) and
FUN_0054e9e0 (`0054ea84`, `0054ea90`); a seventh store, `00543e4d`, rewrites the byte read at
`00543e48`. The one-byte displacement-2 operands over `0x540000..0x552000` also read it at `005452f4`,
`00545c71`, `00547abe`, `0054dab5`, `0054ddbc` and `0054eb25`.

Wall of Earth is spell 19, layer index 3 at payload `+0x20` (`MAGIC-MAPLAYER-040`, `TERR-CELLREC-146`).
FUN_005456d0 walks the six slots from `+0x14` and runs `005457ea SHL byte ptr [EAX],2` once per
non-null one (`005457dd`..`005457f1`), so the Wall of Earth slot takes the multiply like any other,
and FUN_0054e5e0 reads only the count, so it takes the divide the same way. A separate test of that
slot at `005457f3` ORs bits into the static and dynamic planes and does not touch the cost byte. The recompute first resets
the byte from payload `+0` (`0054572c`), or from the `CostCracked` constant at `world+0x5417b` on the
footprint-clearing arm (`005457d3`, `005457d9`; `TERR-STRUCT-071`), so a recompute gives that source
`<< 2` per layer, in 8 bits, whatever was read before. The shipped `CostCracked` value is 6, inside
the range of the table below.

`costbyte.tsv` applies both routines to the shipped range 6..16 (`TERR-COST-052`). One layer: read 1
returns the baseline, read 2 returns `c >> 2` (1..4 for 6..16), read 3 returns 0 except for 16 (1).
Two layers: read 1 returns `c << 2` in 8 bits, so 16 gives 0. Three layers give 0 for 8, 12 and 16.

FUN_0054d210 calls FUN_0054e5e0 for the two cells only on its domain-1 arm (`0054d254`). The
recompute call sites are 16 direct calls in the map module (`xref-recompute.txt`): the crossing
release and occupy (FUN_00545230, FUN_00544ec0), the area module (FUN_0054e730, FUN_0054e9e0), the
building routines (FUN_0054d790, FUN_0054dc70), the sack routines (FUN_005477b0, FUN_005479b0,
FUN_0054f680) and FUN_00545de0. FUN_00545230 skips the recompute when the actor's slot is empty (`00545291`,
`005452a5`); FUN_00544ec0 skips it when the slot is taken (`00545001`, `005450d6`). From a transit
start the call graph reaches a recompute of its two cells only through FUN_00548c60, at the crossing.
No edge orders another actor's transit start between the read and that recompute, and none forbids
it.

**Confidence.** High for the scratch identity, the writers within the instrument, the one divide
per read, the slot-blind divide and multiply, and the table, all quoted instructions and arithmetic.
Medium for the absence of a forbidding edge: its population is the bodies of FUN_00549990,
FUN_0054d210, FUN_00548c60 and the closures of `MOVE-084`; the guards in FUN_00548f70 and
FUN_005492a0 on a destination another actor has claimed were not read. The writer sweep misses word
and dword stores that overlap the byte and `REP MOVSD` record copies.

**Unknown.** Whether two reads of one cell without a recompute happen in play. A runtime trace of
FUN_0054e5e0 and FUN_005456d0 over one cloud confirms or refutes the decay.

### MOVE-086

The cost byte is written from payload `+0` by FUN_005456d0 (`0054572c`), from the `CostCracked`
constant at `world+0x5417b` by the same routine's footprint-clearing arm (`005457d3 MOV DL,byte ptr
[EBX+0x5417b]`, `005457d9 MOV byte ptr [EAX],DL`; `TERR-STRUCT-071`), and by `0054eb77` in
FUN_0054e9e0 (from `MOV DL,byte ptr [EBP]` at `0054eb66`), and the sweep of byte operands of the
form `[reg+reg]` over `0x540000..0x552000` (`byteidx-map-module.txt`) shows further stores at
`00541d28`, `00541d3f`, `00541d46`, `0054555e`, `00547b16`, `00547d74`, `00547dac`, `00547df9`,
`00547e8a`, `0054db06` and `0054de0e`, left unclassified here. The sweep does not match the write of
the ingest FUN_00548720 (`TERR-COST-052`), which uses another operand form. Baselines are read from
the plane when a record is made (`00545468`, `0054d8ad`, `0054e874`, `0054f7ca`). A cell whose static
bit 5 is clear has no record and no divide (`0054e612`).

No cost plane is saved. The sweep finds no displacement-0 byte operand between `00543ea9` and `00545468`, a span that holds
FUN_00544a60, and `SAV-BLOCK-011` and `SAV-CELLREC-017` describe only the two flag planes and the
52-byte records with baseline, count and slots. LOAD in FUN_004d0cb7 builds the terrain first
(`004d130f CALL FUN_005417f0`, ingest), then runs FUN_00544a60 (`004d1345`). After it the direct-call
chains to FUN_005456d0 are `004d13b4` FUN_00539310 > FUN_00539900 > FUN_00539be0 > FUN_0054f680
and `004d143f` FUN_004e3591 > FUN_005476d0 > FUN_00545230 (`callgraph-load.txt`): the sack cells
and an actor release. A tick lays layers only for an effect whose `+0x48` is 0 (`MOVE-084`). A saved
layered cell therefore keeps the ingest byte, with a non-zero count, until a recompute; its first
FUN_0054e5e0 read returns `c >> 2` and stores it.

Two inline readers divide by the byte with no zero test. FUN_005495f0 runs after FUN_00548c60 and
reads the cost of the cell word (`00549674 MOV CL,byte ptr [EDX+EAX]`, `00549682 IDIV ECX`; the arm is
selected by `0054960d` and `00549628`). FUN_0054a620, called at `00544d58` before the occupy
recompute, does the same (`0054a672`, `0054a677`). `IDIV` by 0 raises a divide error. The byte is 0
after a recompute for two layers at cost 16 and three layers at 8, 12 and 16, and by decay with a
non-zero count after read 3 of a one-layer cell at cost 6..15 and after read 4 of a two-layer cell
at cost 6..15 (`costbyte.tsv`), so a layer count of one or two reaches the zero divisor if
two or three reads of the cell occur with no recompute (`MOVE-085`).

**Confidence.** High for the two divide sites and their missing guards and for the load sequence.
Medium for the restore sources, the absence of a cost plane in the serializer and the population of
readers and writers: the instrument is the `[reg+reg]` sweep over the map module, classified by
hand, with no original run, no corpus read and whole-image writes outside `0x540000..0x552000`
unclassified. The `[reg]` displacement-0 store form (about 85 byte stores and read-modify-writes in
`0x540000..0x552000`, among them `0054572f`, `005457d9`, `0054e657` and `0054eb80`) was not swept or
classified; only the sites quoted here are named. Other index forms are not seen.

**Unknown.** Whether LOAD restores the area effects with `+0x48` set, so that no layer is laid again; the effect save routine was not read. What the process does on a zero divisor, since no exception handler was read.
Whether a decayed byte, 0 after read 3 or 4, reaches an inline reader; that needs the double read of
`MOVE-085` and a mover starting or crossing in that cell. Whether a
shipped or generated save holds a layered cell. Whether two or three layers meet on a cost-8 or
cost-16 cell in play; `MAGIC-MAPLAYER-040`'s conflict rules do not forbid it.

## Footprint position and cell transit

Addresses are hexadecimal. Every quoted instruction is asserted byte for byte on both editions' `rom.exe`,
which are identical (162 rows, 0 mismatches). The instrument is capstone disassembly of address ranges plus
direct `CALL rel32` cross-references and a linear byte-store sweep; indirect calls are not followed.

`EXP-0443` was allocated ids `87`..`89` of `claims/move.md` (3 ids) and spent all three, `MOVE-087`..`MOVE-089`.
The next free `move.md` id is `MOVE-090`.

| ID | Claim | Confidence | Status | Evidence |
|---|---|---|---|---|
| MOVE-087 | A mover of footprint side n stores the footprint's top-left cell and a sub-cell offset; the fine point P = cell*256 + sub is that corner, and the centre read for range, edge gap and bearing is P + (n-1)*128 per axis. | High / Medium / Unknown | ● active | [EXP-0443](../experiments/EXP-0443-footprint-slot/EXP-0443.md) |
| MOVE-088 | A cell crossing releases the old footprint, claims, rewrites the position and occupies the new one in one FUN_00548c60 call; an empty release slot or a taken occupy slot skips that cell's recompute, and no deferral exists in those bodies. | High / Medium / Unknown | ● active (amended) | [EXP-0443](../experiments/EXP-0443-footprint-slot/EXP-0443.md) |
| MOVE-089 | Actor removal through FUN_005476d0 releases the footprint at the stored position, recomputing each held cell, and clears the claim bits only after every release succeeds; release is blind to which actor holds the slot. | High / Medium / Unknown | ● active | [EXP-0443](../experiments/EXP-0443-footprint-slot/EXP-0443.md) |

### MOVE-087

The position record at `actor+0x10` holds byte `+0` cell x, `+1` cell y, word `+2` the packed cell, and bytes `+4`
and `+5` the sub-cell x and y. `FUN_005449e0` returns `(cell x << 8) + sub x` and `FUN_005449f0` the same for y
(`005449e4`..`005449ec`). Call that 16-bit value P. The footprint side n is actor byte `+0x49`, read through
vtable `+0x1c` (`MOVE-084`).

The stored cell is the origin of the occupied cells. The occupy routine `FUN_00544d00` reads n once
(`00544d12`), then calls `FUN_00544ec0` for every offset pair 0..n-1 added to the cell bytes from
`FUN_00544a10` and `FUN_00544a20` (`00544e57 ADD EDX,[EBP-0x10]`, `00544e63 ADD EAX,[EBP-0x14]`). The step
release loop in `FUN_00548c60` adds the same offsets to the stored cell bytes (`00548d67`, `00548d78`), and the
claim loops of `FUN_0054abb0` run 0..n-1 from the packed cell (`0054abf6`..`0054ac5c`). The footprint is
therefore the n by n block whose top-left cell is the stored cell. A centre-cell anchor and a stored footprint
centre are both excluded in these bodies: the step, occupy, release and claim loops read here add offsets 0..n-1 to
the stored cell and none subtracts one.

`FUN_004f280d` returns the x centre and `FUN_004f2849` the y centre: P, plus `(n-1) << 7`
(`004f2831 CALL [EDX+0x1c]`, `004f2839 SUB EAX,1`, `004f283c SHL EAX,7`, `004f283f ADD ESI,EAX`), in 16 bits.
Three consumers read it:

- `FUN_004fb702`, the range test, takes the absolute difference of the two x centres and of the two y centres
  (`005545c0` is `NEG` on a negative argument), keeps the larger, subtracts `((n1+n2) << 7) - 0x100`
  (`004fb7a3`, `004fb7a6`), and returns 1 when the result is at most `0x180`, else `(v + 0x40) >> 8`
  (`004fb7b4`, `004fb7c6`, `004fb7c9`).
- `FUN_0054a960`, the edge gap, builds each centre as `(2*cell + n + 0x1ff) << 7` plus the sub-cell, masked to
  16 bits (`0054a994`, `0054a9a9`, `0054a9ac`, `0054aa02`). Since `0x1ff*128 = 0x10000 - 128`, that value is
  P + (n-1)*128 modulo 65536, the same centre. It subtracts the axes, takes the absolute value, subtracts
  `(n1+n2) << 7` (`0054aa1e`), clamps at 0, keeps the larger axis, and returns `(v >> 8) + 1` (`0054aa4f`,
  `0054aa53`).
- `FUN_0054a680`, the bearing, subtracts the centres from the same two routines (`0054a689`..`0054a6bb`).

The direct-call readers of these routines are `FUN_004fb702` at `004fb8af`, `004fba6f`, `004fbc16`;
`FUN_0054a960` in `FUN_0052ab60`, `FUN_0052e110`, `FUN_0052e4d0`, `FUN_0053d9b0`, `FUN_0053ddd0`, `FUN_00548ea0`
and `FUN_005492a0`; and `FUN_0054a680` in 11 owners (`xref-position.txt`).

A step adds the signed step bytes to P itself (`00548cb4 ADD EAX,EBP`, `00548cc5 ADD EBX,EBP`) and stores the
resulting cell bytes, sub-cell bytes and packed word back (`00548d12`..`00548d2c` when the cell is unchanged,
`00548db9`..`00548dda` at a crossing). The same offset moves the corner and the centre whatever n is. At arrival
both sub-cell bytes return to `0x80` (`00548dfe`..`00548e03`), as they do in the constructor and in the placement
and portal writers (`005445ea`, `0054984b`, `0054984f`). A resting mover therefore has centre
cell*256 + 0x80 + (n-1)*128: for n = 2 that is the grid line shared by its two cell columns, and for n = 3 the
middle of its middle cell, which is the geometric centre of the block in both cases.

Fourteen routines store the four position bytes within 24 instructions of one another (`posstores.txt`):
`00543d30`, `00544530`, `00544550`, `005445d0`, `00544600`, `00544760`, `005447d0`, `00544810`, `00544980`,
`005449b0`, `00548720`, `00548c60`, `005495f0` and `0054eec0`.

**Confidence.** High for the centre formulas, the footprint origin and the step arithmetic, each an asserted
instruction sequence read whole, and for excluding a centre-cell anchor and a stored centre in the step, occupy, release and claim loops read. Medium for the
writer list and the reader lists: the instrument is a linear sweep with a 24-instruction window and direct
calls only, and it misses stores split over longer windows, changed base registers, `REP MOVS` copies of the
record, computed indices and indirect callers.

**Unknown.** Which shipped or authored actors have n above 1; no census was run. Observed positions of a size
above 1 in play. The full bodies of the 13 writers other than `FUN_00548c60`, which were located and only their quoted stores asserted.

### MOVE-088

`FUN_00548c60` (the step) runs, when the cell bytes differ after the step is added (`00548cf3 CMP DL,CL`):

1. a release loop over the n by n footprint at the old position, one `FUN_00545230` call per cell
   (`00548d7e`), which stops at the first call that returns 0 (`00548d83 TEST EAX,EAX`, then `JE 0x548d99`);
2. `FUN_0054abb0` with the old packed cell (`00548da8`), which stores the claim cell at `mover+0xa6`
   (`0054abcd`) and ORs the dynamic bits over the n by n footprint (`MOVE-084`);
3. the new cell bytes, sub-cell bytes and packed word (`00548db9`..`00548dda`);
4. `FUN_00544d00`, the occupy (`00548de0`). Its return value is not tested before the arrival recentre at
   `00548de5`..`00548e03`.

All four are in one call of the step, and the loop never yields. The two release outcomes decide the recompute
of one cell. `FUN_00545230` selects the slot from the actor's movement domain (`0054527f CALL [EAX+0x20]`):
domain 1 or 2 uses record `+0x04`, domain 3 uses `+0x08`. A missing cell record returns 0 at `00545267`, before any
slot test. A domain of 0 or above 3 takes neither arm (`00545284`, `0054528c`): it skips the slot test and runs the
recompute at `005452ce` unconditionally. A slot that is empty returns 0 at `00545293` or `005452a7`, before the
clear and before the recompute. A slot that holds any
non-zero value is cleared (`00545299`, `005452ad`) and `FUN_005456d0` runs at `005452ce`; the routine then returns
1. The slot test is a non-zero test and does not compare the held pointer with the actor, so the release clears
a slot another actor holds.

The occupy `FUN_00544d00` calls `FUN_00544ec0` for each cell in row-major order and continues only on a
return of 1 (`00544e73 TEST EAX,EAX`, `00544e75 JNE 0x544e9d`); the first return of 0 ends it with 0 and the
later cells are not entered (`MOVE-087` gives the order). In the domain 1 and 2 arm of `FUN_00544ec0`, a cell record whose byte `+0x2c` is non-zero and not `0x1a`
(`00544f4a`, `00544f5a`) first builds a temporary caster through `FUN_004d20c8` or `FUN_004d2105`
(`00544fb1`, `00544ff9`; `MAGIC-235`); the builder appends a new object to the list at `this+0x2c` (`004d1fb5`).
That happens before the slot test, whether or not the slot is taken. The domain 3 arm (`005450b8`..`0054511d`)
contains no such call. A slot that is already taken (`00545001 CMP [EAX+4],0`, `005450d6 CMP [EDX+8],0`) calls
`FUN_0054a200` and returns 0 with no slot store and no recompute. `FUN_0054a200` and `FUN_0054a1f0` are each one `RET 4` (`0054a1f0`, `0054a200`). An empty slot is
written with the actor (`00545021`, `005450f6`) and `FUN_005456d0` runs at `00545043` or `00545118`.

A deferred recompute needs a stored mark or a queue entry that a later routine reads. The release, the occupy,
the per-cell routine, the stubs and the step contain no store of that kind. Their stores are the slot, the
write-back of the record copy through `FUN_0054fc70`, the position, the claim cell, the cached bytes at
`mover+0x86`..`+0x89` (`0054538a`..`005453c0`) and, in the per-cell occupy, the trigger caster's construction and
list append. None is read here as a recompute mark; the consumers of that list were not read. The recompute therefore runs in the tick of the crossing, in the same call, for every
cell whose release found a held slot and every cell whose occupy found an empty one. A cell whose release found
an empty slot, or whose occupy found a taken one, gets no recompute from that call, and the first such
release or occupy skips the cells after it in that footprint. The mover's `+0x76` word and `FUN_00548720`, which
has its own release, claim and occupy call sites (`00548b0b`, `00548b35`, `00548b77`), were not read.

**Confidence.** High for the order inside the step, the two slot branches of the release and of the per-cell
occupy, the no-op stubs, the loop exits and the untested occupy return, each an asserted instruction. Medium
for the absence of deferral: the population is the bodies named above and the two stubs; stores outside
them, the consumers of the trigger caster list, indirect callers, and any routine that reads the bytes at `mover+0x86`..`+0x89` were not enumerated.

**Unknown.** An observed tick. Which tick of a transit holds the crossing, which depends on speed
(`MOVE-084`). Whether `FUN_00548720` follows the same order.

**Amended.** The unread `FUN_00548720` clause is closed by `MOVE-093`: Ghidra's `FUN_00548720` is an unrelated cost helper, and the second release, claim and occupy call sites (`00548b0b`, `00548b35`, `00548b77`) belong to an unreferenced run at `00548900`..`00548c4f` that follows the same order. `MOVE-094` gives the crossing tick for the Unknown above. The rest of the claim stands.

### MOVE-089

`FUN_005476d0` is the removal of an actor from the cell records. It reads n through vtable `+0x1c`
(`005476e1`) and runs the same row-major release loop as the step, at the stored position, one
`FUN_00545230` call per cell (`00547726`). A return of 0 ends the routine with 0 (`0054772d JE 0x54779a`), so
the claim state below is skipped. After every release has returned 1 it reads the word at `mover+0xa6`
(`00547747`). A zero word, or a word equal to the packed cell at position `+2` (`00547756`), ends the routine
with 1. Any other word is passed to `FUN_0054af40` (`00547760`), which clears the dynamic bit (`0x80` for
domain 3, `0x40` for domain 1 or 2) over the n by n cells at that claim cell except those inside the current
footprint (`0054b021`, `0054b0e6`, `MOVE-084`); `+0xa6`, `+0x80` and `+0x76` are then zeroed
(`0054776d`, `0054777a`, `00547787`).

The four direct callers are `004e3c9b` in `FUN_004e3591`, `004f4808` in `FUN_004f47e6`, `004f4f9a` in the
teardown `FUN_004f4f5d`, and `004fb183` in `FUN_004fb0b5`. The teardown then calls `FUN_00548e10` at `004f4fa5`,
which sets both sub-cell bytes to `0x80`. When every footprint cell's release finds a held slot, each is cleared with its recompute, so a removed mover
leaves no occupancy in those cells and no cost contribution: the alternative that the mover's contribution stays
in the plane after removal is rejected for that all-cells-held case, and a mover removed after a completed
crossing is fully released by this call. A release that finds an empty slot ends the routine early and leaves the
later cells' slots and the claim bits `0x40` or `0x80` as they were. A crossing is not a state that can be cut short between ticks, since release, position rewrite and
occupy are one call (`MOVE-088`).

A footprint that was refused part-way keeps its row-major prefix (`SAV-CELLFAIL-583`). The release at removal
walks the whole footprint at the new corner. Where a refused cell is held by another actor, the non-zero slot
is cleared and recomputed, since the release tests no identity; the walk then reaches a never-entered cell,
finds it empty, and returns 0, leaving the claim state in place. That consequence is read from the bytes and
not run.

**Confidence.** High for the loop, the exit on an empty slot, the claim clear and its conditions and the
teardown recentre, each an asserted instruction. Medium for the consequences on a refused footprint, which
combine these bytes with `SAV-CELLFAIL-583` and `MOVE-088` without a run, and for the four callers, a direct
`CALL rel32` list that misses indirect callers.

**Unknown.** Whether any caller tests the return of 0. Whether each caller removes the actor in the tick in which
it leaves the world. Whether two actors reach one cell in play and one of
them is then removed. Native behaviour of the claim bits after the early return.

## Walk step at order recovery

Evidence is a static read of `rom.exe` (one image on both lawful installs) in Ghidra 12.1.2 headless; no process
was run. `EXP-0451` was allocated ids `90`..`92` of `claims/move.md` and spent `MOVE-090`.

| ID | Claim | Confidence | Status | Evidence |
|---|---|---|---|---|
| MOVE-090 | A fresh walk call of FUN_00548f70 writes a sub-cell step in that call only when the facing byte already equals the direction to the first path node; otherwise it turns and the step comes at the next call. | High / Medium / Unknown | ● active | [EXP-0451](../experiments/EXP-0451-attack-cycle/EXP-0451.md) |

### MOVE-090

- Entry: executor row 1 at `00531270` calls `FUN_00548f70(actor, ord+0x0a, 0)` at `0053129c`. `FUN_00545c30` is the centred test (both sub-offsets 0x80); a mover that is not centred stores progress 3.
- A fresh walk forces a full search; the dynamic list is empty, so the near search `FUN_00549a90` runs and then the stepper `FUN_00549990`.
- The stepper computes the facing to the head node with `FUN_00544300` into `mover+1`. When `mover+0` equals `mover+1` it calls `FUN_0054d210` (step vector `mover+0xb0` and `mover+0xb1`, count `mover+0xaa`) and `FUN_00548c60`, which writes the sub-cell position at `+4` and `+5` in that call. Otherwise `FUN_0054a210` turns only; it snaps the facing at once when the difference is under 0x21, and the step is written by the next call. An empty dynamic list calls `FUN_00548e20`, a turn only.

**Confidence.** High for the control flow, each branch read from `rom-walk.txt`. Medium for the composed outcome at the recovery tick, since the facing at that tick is a run-time value.

**Unknown.** The facing byte's usual value when a recovered actor was last stopped; whether the turn path has a second step in the same call under any condition not shown in the listed routines.

**Evidence.** [EXP-0451](../experiments/EXP-0451-attack-cycle/EXP-0451.md), `evidence/listings/rom-walk.txt`, `evidence/listings/rom-row-callees.txt`

## Second step routine, transit ticks and teardown claim state

Evidence is a static read of `rom.exe` (one image on both lawful installs) with a capstone sweep; no process was run. `EXP-0452` was allocated ids `93`..`96` of `claims/move.md` and spent all four.

| ID | Claim | Confidence | Status | Evidence |
|---|---|---|---|---|
| MOVE-093 | The step routine at 00548900..00548c4f, with no direct reference found, releases the old footprint, sets the claim at the old cell, rewrites the position, then occupies, with no arrival compare and no test of the occupy return. | High / Medium / Unknown | ● active | [EXP-0452](../experiments/EXP-0452-area-cost/EXP-0452.md) |
| MOVE-094 | In the live step a transit from the centre crosses on tick ceil(128/s) on a positive axis and floor(128/s)+1 on a negative one, start tick 1; the step calls release and occupy on that tick only, whatever the slots hold. | High / Medium | ● active | [EXP-0452](../experiments/EXP-0452-area-cost/EXP-0452.md) |
| MOVE-095 | A removal whose release ends early on an empty slot returns before the claim clear: the transit's claim bits and the later cells' slots stay as they were; no store of the claim bits follows in the teardown body or FUN_00548e10. | Medium / Unknown | ● active | [EXP-0452](../experiments/EXP-0452-area-cost/EXP-0452.md) |
| MOVE-096 | A unit with health at or below 0 is not stepped; its slot, claim bits and +0xa6 stay frozen through the death countdown, and the teardown releases at the stored cell and clears the claim when +0xa6 is not the current packed cell. | High / Medium / Unknown | ● active | [EXP-0452](../experiments/EXP-0452-area-cost/EXP-0452.md) |

### MOVE-093

`MOVE-088` left `FUN_00548720` unread. Ghidra's `FUN_00548720` ends at `005488cb` and is an unrelated cost helper; the second step routine is the code after it, `00548900`..`00548c4f`, which has no function entry. Its order, from `evidence/listings/rom-step.txt`:

1. the move computation, reading record `+0x72` (`005489e4`, `00548a0a`); a crossing is a flip of bit 0x80 of the moving axis's sub byte together with a cell change (`00548a6c`..`00548a9c`);
2. on a crossing, a release loop over the n by n footprint, one `FUN_00545230` call per cell (`00548b0b`, row-major), stopping at the first return of 0 (`00548b10`, `00548b12`);
3. `FUN_0054abb0` with the old packed cell (`00548b35`);
4. the position rewrite: cell bytes, sub bytes, packed word (`00548b3d`..`00548b73`);
5. `FUN_00544d00`, the occupy (`00548b77`), then `RET 4`. The return is not tested.

The stride `+0x72` is recomputed at the crossing (`00548bd0`..`00548c49`). A branch with no crossing writes the position only (`00548b86`); a same-cell half-flip from an off-centre start sets both sub bytes to 0x80 (`00548bb9`..`00548bc2`). The routine has no `+0xaa` arrival compare.

The order matches the live step of `MOVE-088` (release `00548d7e`, claim `00548da8`, position, occupy `00548de0`). References: `evidence/listings/orphanrefs.txt` finds no `E8`, `E9` or `0F 8x` rel32 into the range from outside it and one dword in the whole file with a value in the range (file offset `0x174d83`, inside a `CALL` operand). `xref.txt` and the displacement scan `scan.txt` (value 0x72) find no other reader of record `+0x72` in `0x4d0000..0x54ffff`; stores to it are at `00544d66`, `00548c49`, `00549641`, `00549651`, `0054968c`.

**Confidence.** High for the order and the untested occupy return, each asserted instruction on both roots (`asserts-en.txt`, `asserts-ru.txt`). Medium for "unreferenced": the population is direct relative transfers and literal dwords over the file; a computed jump into the range was not excluded.

**Unknown.** Whether the range is ever executed. Whether a pre-release build called it.

**Evidence.** [EXP-0452](../experiments/EXP-0452-area-cost/EXP-0452.md), `evidence/listings/rom-step.txt`, `evidence/listings/orphanrefs.txt`, `evidence/listings/xref.txt`

### MOVE-094

The live step `FUN_00548c60` runs once per tick: `FUN_00549990` on the start tick (`+0xac` reset to 0, then the step) and `FUN_00548f70` then `FUN_005495f0` on later ticks while the sub bytes are not both 0x80 (`00548fda`, `00549605`). It calls release, claim and occupy only on the call whose cell bytes differ after the step is added (`00548cf3`); otherwise it writes the position and jumps to `00548de5`.

Starting from sub 0x80 with signed axis step s (`FUN_0054d210`, clamp 1..63, `TERR-MOVE-056` gives the shipped range 4..32), the crossing tick counted with the start tick as 1 is ceil(128/s) for a positive axis and floor(128/s)+1 for a negative one; the transit length is N = ceil(256/s), where the `+0xaa` compare recentres to 0x80 (`evidence/listings/crossing.txt`, v = 1..63). At v = 16: N = 16, positive tick 8, negative tick 9. Anti-diagonal directions 1 and 5 add 0xffff to the fine X when both sub bytes are 0 (`00548cd3`..`00548ce8`); diagonal steps are `FUN_0054d210`'s truncated v*0.707, so the table is exact for the four axis directions only.

The slot outcomes are those of `MOVE-088`: an empty release slot returns 0 and skips that cell's recompute; a taken occupy slot returns 0 with no slot store and no recompute. Neither defers: the step calls them on no other tick (`MOVE-088`). The scope is the step `FUN_00548c60`, not every caller of the occupy or the claim clear.

**Confidence.** High for the tick rule given the start state and the step bytes: the compare and the add are asserted instructions and the table is arithmetic over them. Medium for the diagonal and tie-break paths, read but not tabulated.

**Unknown.** The arrival path of `FUN_005495f0`, which the walk calls on later ticks: it calls the claim clear `FUN_0054af40` at `005496e0` and `0054980b` and the occupy `FUN_00544d00` at `00549853` (the occupy when the cell byte is `0x1a`, `00549788`, `0054984b`). These sites are in `rom-transit.txt` and `xref.txt`, were read but not traced, and are not reached through a crossing.

**Evidence.** [EXP-0452](../experiments/EXP-0452-area-cost/EXP-0452.md), `evidence/listings/rom-step.txt`, `evidence/listings/rom-transit.txt`, `evidence/listings/crossing.txt`

### MOVE-095

`FUN_005476d0` (the removal) runs the release loop over the footprint and returns 0 at the first release that returns 0 (`0054772d JE 0x54779a`), before the claim clear. The clear (`FUN_0054af40`, `00547760`) is reached only after the whole loop and only when `+0xa6` is non-zero and not the current packed cell. `FUN_00548e10`, which follows the removal in the teardown, sets both sub bytes to 0x80 and writes no plane byte.

The plane byte at `world+0x20000 + packed cell` carries the claim bits 0x40 and 0x80 (`MOVE-084`). It is rewritten whole by `FUN_005456d0` (callers `00545043`, `005450a2`, `00545118`, `00545174`, `005452ce`, `00545e00`, `00545f59`, `00547984`, `00547a98`, `0054d994`) or cleared by `FUN_0054af40` and `FUN_0054ac70`. None of these is called in the teardown body or in `FUN_00548e10` after the early return. The claim bits at the `+0xa6` footprint therefore read as set. Slots of the cells after the empty one keep what they held: zero, or a pointer to another unit, or a stale pointer to the removed unit.

Plane readers found by displacement sweeps (`scan.txt`, values 0x20000 and 0x200): `541dd0`, `54bf10`, `54c100`, `54c2f0`, `54c6f0`, `54e220`, `548f70` (`005491df`), `549880`, `549a90`, `54a060`, `54e070`, `54ed30`, `54eec0`, `546fa0`, `5470a0`.

**Confidence.** Medium: the stores and the early return are asserted; the claim bits reading composes them without a run, and no reader of the stale bits was traced to an effect.

**Unknown.** Four teardown callees after `004f4f9a` were not read: the two virtual calls through `[edx+0x40]` (`004f4fc2`, `004f5001`), `FUN_0050e8e2` (called twice) and `FUN_0051ab50`. What each plane reader does with a stale claim bit. Slot-pointer dereferencers were not enumerated. Observed play.

**Evidence.** [EXP-0452](../experiments/EXP-0452-area-cost/EXP-0452.md), `evidence/listings/rom-removal.txt`, `evidence/listings/rom-claim.txt`, `evidence/listings/rom-slots.txt`, `evidence/listings/scan.txt`

### MOVE-096

The actor tick `FUN_004f37be` does nothing when action `+0x54` is 0x10. When health `+0x94` is at or below 0 it takes the dying branch, which does not call the executor `FUN_005310e0` (its direct rel32 and dword references are `004f398c`, in the health above 0 branch, and `0053109b` in `FUN_00530c90`; `xref.txt`), so no walk and no step runs. The first dying tick sets `+0x13c` to 1, calls `FUN_004f4c97`, halves `+0xbe` and sets `+0x6c`; later ticks count `+0x6c` down. At zero, if `+0x4a` is above 1 health is set to 0xfc18; when health is at or below -10 the actor sets `+0x54` to 0x10 and calls the teardown `FUN_004f4f5d`.

Through the countdown the occupancy slot, the claim bits, `+0xa6`, `+0x80` and `+0xac` are unchanged. The teardown calls `FUN_005476d0` (`004f4f9a`, return not tested): it releases at the stored position (the new cell if the transit had crossed), then clears the claim at `+0xa6` (the target cell before the crossing, the old cell after it) unless it equals the current packed cell (`00547756 CMP AX,[EDX+2]`; for a footprint wider than 1 a `+0xa6` cell inside the footprint but not the packed cell is cleared), zeroes `+0xa6`, `+0x80` and `+0x76`, and `FUN_00548e10` recentres. When the release ends early, `MOVE-095` applies. Other callers of `FUN_005476d0` (`004e3c9b`, `004f4808`, `004fb183`) were not covered.

**Confidence.** High for the dying branch's exclusion of the executor and the teardown order, asserted instructions. Medium for the stored position at death: it composes the step's store order with the freeze and no run was observed.

**Unknown.** Other writers of `+0x54`. Observed play. Whether a unit dies in the tick it crosses.

**Evidence.** [EXP-0452](../experiments/EXP-0452-area-cost/EXP-0452.md), `evidence/listings/rom-death.txt`, `evidence/listings/rom-removal.txt`
