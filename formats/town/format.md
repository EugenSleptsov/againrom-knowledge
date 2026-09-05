# Town reactions and tavern interior animation

This page describes executable state, not a stored file format. Exterior contracts are
conditional on the named town pointer handler and paint hub being reached.
Physical-pointer message delivery and audio-device results remain Unknown.

## Selector and labels

Message `200h` can reach town vtable `0059b000+4c`, wrapper `004b4c50`, then
`004b7080`. Each delivered update samples the mask and writes `town+b4`;
there is no enter-only comparison. Mask bytes map to selectors
`80h→2`, `90h→1`, `a0h→8`, `b0h→16`, `c0h→4`; other values map to `-1`.
Selectors 1, 2 and 4 draw Shop, Tavern and School label pictures respectively.
A delivered blank-mask update removes those labels but does not reset active
shop/tavern/school frames or enable bits. Focus loss is not that proved update.
(`TOWN-399`)

## Entrance state

| entrance | arm and retained state | advance and end | conditional sound request |
|---|---|---|---|
| Shop | Selector 1: `rand()%100 > 95` sets bit 1 if clear. Sprite `+1d8`, frame `+1e0`, shipped count 30. | One increment per admitted hub; equality with sprite count resets frame 0 and clears bit 1. Leave does not stop the active cycle. | Sound latch `+ac=0` permits `SFX/Town/Shop/Enter.wav`; clears school latch. (`TOWN-400`) |
| Tavern | Selector 2 ORs bit 2 on every delivered update. Sprite `+160`, frame `+168`, count 10. | One increment per admitted hub; count equality resets frame 0 and clears bit 2. Another update can rearm. | Frame 0 before increment permits `SFX/Town/Point.wav` and cancellation calls for shop/school sounds. (`TOWN-401`) |
| School exterior | Selector 4 ORs bit 4 without resetting fighter/mage frames `+1cc/+1d4` or static directions `005f20c8/005f20cc`. Both sprites have 11 frames. | At endpoint 0, remainder 96..99 starts +1; at 10 it starts -1, otherwise endpoint wait. Return to 0 with -1 clears shared bit 4. The hub still calls the second helper that tick; either clear can pause both thereafter. | Latch `+b0=0` permits `SFX/Town/School/Point.wav`; clears shop latch. Leave does not itself reverse the figures. (`TOWN-402`) |
| Gate | The hub always calls the gate helper. If availability helper returns -1, select `door/T08.bmp` and return. Otherwise poll current pointer mask every update. | Inside selector 8 decrements `+1a0` toward 0, outside increments toward 8; `+19c` selects `door/T00..T08.bmp`. Re-entry reverses existing progress. | Latch `+1a4` changes request `GateUp.wav`/`GateDn.wav`. (`TOWN-403`) |

Guard animation is separate from the door: eight-frame sprite `+e4`, frame
`+e8`, direction `+ec`. Pointer updates over an unavailable gate set -1;
other recovered arms set +1. Its unconditional hub helper clamps at endpoints,
stops on overshoot, and uses `Guard1.wav`/`Guard2.wav` on direction-latch changes.
The same hub can request then release a guard sound. (`TOWN-403`)

All sound requests above are conditional calls to `00453b08`, not guarantees
of audible output. Null pointers, status lookup, audio availability and indirect
buffer operations can prevent playback. Latch behavior is separate from the
animation frame. (`TOWN-399`, `TOWN-400`, `TOWN-401`, `TOWN-402`, `TOWN-403`)

## Clock and independent ambience

Paint admits a single hub when unsigned elapsed time is strictly greater than
67 ms. No catch-up loop exists in that block. Sign and fluger are independent
random triggers, bits `40h`/`20h`, with 10/8 frames and conditional Flag/Flugel
sounds. They are not shop entrance reactions. Birds and the other town wildlife
retain separate arming mechanisms. Stars belong to selector 16, outside the
four entrance selectors. (`TOWN-404`)

The hub has nine non-bird tested flag bits, invoking ten gated helpers because
bit 4 calls two. Gate and guard helpers are unconditional. This corrects only
the eight-count clause of `TOWN-158`. (`TOWN-405`)

### Bird episode

Bird state is split between the view and process statics. The view owns clock
`+b8`, group `+bc`, count `+c0`, a nine-pointer array `+c8`, three progress
words `+d8/+dc/+e0`, and active bit `80h`. Process statics `005f20b8/bc/c0`
own the once latch, admitted-hub clock and next episode delay. Paint arms only
when strict elapsed exceeds the retained 1000..2999-ms delay and bit `80h` is
clear. Entry clears the view flags and resets the view clock; the next arm
zeros all progress, while no direct exact-address static reset is known.
(`TOWN-415`, `TOWN-417`)

The loader binds nine `TownBirds/Birds1..9/sprites.16a` sheets, each 57
frames. An arm selects group0..2 and count1..3; composition uses
`array[group*3+i]` for `i<count`. Every admitted hub increments all three
progress words once while active. Paint draws each selected nonterminal bird
through the class-specific `.16a` frame-draw slot+18, then draws keyed
`Town_add.bmp`, including on
the final active paint. It clears bit `80h` only after every selected sprite is
terminal. Count1 conditionally requests `Birds1.wav`; count2/3 requests
`Birds2.wav`, both repeat0. (`TOWN-415`, `TOWN-416`)

### Horse, baba and dervish

The exterior loader selects one of five horse positions, one of four baba
positions and a different one of four dervish positions. It binds the chosen
horse's A1..A3 sheets at `+128`, baba's A1..A2 at `+fc`, and the selected
dervish sheet at `+150`. All installed horse sheets have15 frames, baba A1
has31 and A2 has32, and dervish has30. (`TOWN-439`)

| family | entry and activation | admitted hub progression | terminal / idle |
|---|---|---|---|
| Baba | selects A1, current -1 and delay2000..3999 ms; a strict elapsed paint test rerolls A1/A2 and delay2000..6999, then sets bit `200h` | one forward frame through the selected sheet | count31/32 clears the bit and writes -1; retained sheet draws frame0 until another delay |
| Horse | selects A1, current -1 and its own delay2000..3999 ms; a separate strict test rerolls A1/A2/A3 and delay2000..6999, then sets bit `100h` | one forward frame through a 15-frame sheet | count15 clears the bit and writes -1; retained sheet draws frame0 until another delay |
| Dervish | current0 and bit `400h` are set on entry | `(current+1)%30` | wraps without clearing; a null sheet writes -1 but retains the bit |

Every eligible own-paint runs the two delay tests; the admitted hub gives each
active family at most one step and has no catch-up loop. Horse selector `+144`
is retained but dormant at entry current -1, then overwritten by every horse
arm before its frame-specific sound gates. Baba and horse therefore combine
entry selection with independent paint-time reselection; dervish is
entry-active and continuously cyclic. (`TOWN-440`, `TOWN-441`, `TOWN-442`,
`TOWN-443`, `TOWN-444`)

The sound loader binds Horse2/Horse3/Horse1 to `+a0/+a4/+a8`. Non-null and
status0 gates request Horse1 at A3 frame1 and Horse2 at A1 frame14, A2
frames8/14 and A3 frame14, repeat0 and priority128. No direct Horse3 request
receiver or baba/dervish-specific request is present in the bounded town
closure; indirect/computed and runtime audio results remain open.
(`TOWN-445`)

Leave and destruction release the three family assets and horse sound
objects; re-entry reloads positions/assets, restores horse and baba idle A1
state and arms dervish at0. The families share only flag word `+208` and the
process-static admitted-hub clock. Exhaustive mask-selector outputs are
-1/1/2/4/8/16, so its default OR does not directly arm bits
`100h/200h/400h`; horse/baba delay arms and dervish entry arm remain separate.
(`TOWN-446`, `TOWN-447`)

### Statue star and crowd

The star loader binds nine 64×44 `Town/stars/S00..S08.bmp` pictures at `+1ac`,
initializes current `+1c0` through `-1→0`, and stores S00 in selected pointer
`+1bc`. That pointer draws at view-relative `(340,288)` whenever non-null.
Selector16 ORs bit `10h`; the admitted hub then advances S01..S08. A step
beginning at current0 conditionally requests `Stars.wav`, repeat0. The step
reaching current9 hides the pointer and clears the bit; it does not wrap.
(`TOWN-418`)

Each terminal star call increments process-static `005f20c4`. On the tenth it
resets current and the static to0, but still ends hidden; a subsequent selector
arm begins the visible sequence again. Entry independently restores S00 and
does not directly clear that exact static. Thus the star is rearm-driven, not
an autonomous cycle. Alias/bulk reset and physical pointer delivery remain
Unknown. (`TOWN-419`)

Crowd has no reached paint/progression owner in the bounded 39-function town
range. Entry loads `Crowd.wav` at `+74` and conditionally requests repeat1;
reload, leave and destruction reach detach/release/zero cleanup. This is a
bounded sound-route result, not a whole-image absence of a visual crowd.
Bird and star share only the admitted hub after separate activation; crowd
bypasses it on the measured route. (`TOWN-420`)

For mouse-family routing, a non-null `control+34` handler bypasses the child
broadcast, including after a zero return. Zero still permits the final vtable
slot decision. With a null handler, the broadcast runs. The keyboard `+38`
path differs: a zero handler result permits a child broadcast. This corrects
the identical-fallback clause of `TOWN-211`. (`TOWN-406`)

## Limits

Complete residence behavior depends on event delivery. No static conclusion
here proves that a stationary physical pointer receives no updates. Frame and
latch persistence through focus changes, screen reuse, untaken indirect calls,
or changing gate availability is not closed. School interiors are outside this
contract. (`TOWN-399`, `TOWN-402`, `TOWN-403`)

## Tavern interior draw and clocks

The central child owns four interior sequences. Its painter first draws
`CenterArea.bmp`, then candle and cauldron, then a mode-selected tender episode.
Coordinates below are relative to the parent origin. Each sequence stores an
index at `+18` and a cached picture pointer at `+14`; drawing uses the pointer,
not a fresh index lookup. Loaded pointers must be non-null. (`TOWN-407`)

| series | child offset | loaded files | position | progression when the central painter runs |
|---|---|---|---|---|
| Candle | `1dc` | `candle/t0000..t0009.bmp` | `(160,48)` | Draw current, then share the strict >100-ms gate with cauldron; `(index+1) % (count-1)` visits 0..8. (`TOWN-408`) |
| Cauldron | `20c` | `cauldron/t0000..t0020.bmp` | `(420,160)` | The same admitted step visits 0..19. Neither series selects its last loaded entry through this cycle. (`TOWN-408`) |
| Tender breath | `23c` | `tender/breath/br0001..br0024.bmp` | `(240,152)` | Mode2 draws before one strict >83-ms forward step. Completion clears mode/index but retains the last cached picture. (`TOWN-410`) |
| Tender drink | `26c` | `tender/drink/dr0001..dr0040.bmp` | `(240,152)` | Mode1 draws before a bounded forward/reverse step. The top step immediately reverses to index38; reverse completion at0 disables the episode. (`TOWN-411`) |

Tender delay is `3000 + rand()/16`, range3000..5047 ms from the measured
15-bit random result. Strict elapsed >delay arms mode1/direction1 for odd
delay, or mode2 for even delay. That test does not require idle mode. The
>83-ms tender gate resets the same tender timestamp after one step; completion
chooses the next delay. Mode0 draws neither tender series. These are conditional
gates, not measured frame rates or uniform random probabilities. (`TOWN-409`)

Two consequences matter. A long paint gap can rearm a descending drink toward
ascent. A completed breath retains cached `br0024.bmp` while its index is0;
the next episode first draws that cached last picture and then selects
`br0002.bmp`, unless a reload or another writer intervenes. No catch-up loop
is present in either measured animation gate. (`TOWN-408`, `TOWN-410`, `TOWN-411`)

## Tavern lifetime and sound boundary

Entry reloads all four arrays and resets each index and cached picture to the
first entry. Timestamps, mode and direction reside at static addresses, not in
those arrays. Exact-address references to eight named globals all belong to
the central painter; this does not exclude indirect or overlapping writes.
Leave passes all four objects to a cleanup helper whose body is not closed
here. Actual release, destruction and cross-visit state remain Unknown.
(`TOWN-412`)

Conditional requests are steam (>10000-ms separate clock), chair and
`Town/Shop/Breath.wav` (breath arm), drink (index30 on every paint, either
direction), glotok (reverse completion), water (parent own-paint, repeat1),
and enter (entry, repeat0). Other named requests use repeat0. A non-null sound,
status-zero result and further audio-buffer gates are required. Request
sites do not prove audible playback. (`TOWN-413`)

The local painter is time-driven. Its complete invocation route and cadence,
indirect interaction effects and all alias writers are not closed. The party
pickers and read central input bodies contain no direct interior-state reset;
their complete callee effects remain Unknown. The selected-character preview
is outside this contract. (`TOWN-414`)

## School training presentation

The four `movies/training/{mage,fighter}/{m,tr}` BMP series belong to the
school-room object. The school loader calls the only direct family loaders:
mage `tr0000..tr0022` fills pointer array `+e4/+e8`, mage `m0001..m0011`
builds sequence `+164`, fighter `tr0000..tr0018` fills `+198/+19c`, and
fighter `m0001..m0009` builds sequence `+218`. Both shipped roots contain
these four gapless groups, 62 BMPs per root. (`TOWN-427`)

The room painter gives each class one lower-area anchor. The mage side is
room-relative `(0,200)` and fighter side `(320,200)`. On each side an active
`m` bit selects the sequence's cached picture; otherwise the current `tr`
picture is drawn at the same anchor. The modes replace one another; they are
not layered and are not a separate movie surface. (`TOWN-428`)

| family | activation | admitted progression | terminal behavior |
|---|---|---|---|
| Mage `tr` | entry/current-class transition, bit1 | one modulo-23 step under strict elapsed `>83` ms | wrap to index0 clears bit1 |
| Fighter `tr` | entry/current-class transition, bit4 | one modulo-19 step under the same gate | wrap to index0 clears bit4 |
| Mage `m` | idle side, strict elapsed `>3000+rand()/10`, bit2 | forward indices0..10; terminal cached10 hold; reverse starts from index8 to7 | reverse terminal0 clears bit2 before the `m` draw |
| Fighter `m` | independent idle side, same delay shape, bit8 | forward0..8; terminal cached8 hold; reverse7..0 | reverse terminal0 clears bit8 before the `m` draw |

`tr` starts after a class change and participates in the pending-transition
input block until its counter returns to zero and the column reaches that
class endpoint. Mage index at least5 can start column step+1; fighter index at
least6 can start step-1. Those starts conditionally request
`SFX\Town\School\Rotate.wav`. (`TOWN-429`)

Each idle armer requires both the `tr` and `m` bit for its own side to be
clear. A busy side refreshes its timestamp. With the accepted 15-bit random
result, `3000+rand()/10` spans 3000..6276 ms; the comparison is strict.
Arming sets direction1 and sequence index -1. The shared paint gate admits at
most one step per side and has no elapsed-time catch-up. (`TOWN-430`)

At mage `m` ascent completion, the direction changes to reverse and samples
hold threshold `((rand()*20)/0x7fff)%20 + 20`, range20..39. This is not
`20+rand()%20`, and no uniform distribution is established. The sequence
index is then set to8 while the cached pointer remains terminal index10. The
terminal stays selected during hold; the first reverse selection is index7,
so return omits indices9 and8. Fighter samples the same threshold, retains
index8 at its terminal and reverses normally to7. An active same-side `tr`
forces an ascending `m` toward reverse rather than instantly removing it.
(`TOWN-431`, `TOWN-432`)

Entry rebuilds both sequences, resets the local flags/counters and primes the
two `tr` pointers at index0. Leave and destruction call both family cleanup
routines. The idle clocks, random extras and hold counters are static state
and are not explicitly cleared by the read lifecycle bodies; object identity,
alias writers and cross-visit timing remain Unknown. (`TOWN-433`)

The literal/direct-call/draw population converges on this one school object:
one owner per format literal, one direct caller per family loader, no parsed
code pointer to either loader, and one room painter consuming all fields. No
`m` arm/step/draw body directly requests a sound. Computed filenames, targets,
other school audio paths, delivered paint cadence, visible output and audible
results remain open. (`TOWN-434`)
