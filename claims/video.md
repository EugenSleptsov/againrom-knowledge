# VIDEO — music, sound effects and cutscenes

Level 2 public functional ledger. This edition preserves the functional
conclusion, confidence, status and private evidence identity while omitting
instruction listings, executable-address inventories and reconstructable
shipped-content tables. The private research snapshot named by
[`SOURCE.md`](../SOURCE.md) retains the complete evidence. Format of this file:
[registry.md](registry.md). IDs are permanent.

## Music

| ID | Public functional claim | Confidence | Status | Evidence |
|---|---|---|---|---|
| VIDEO-MUSIC-001 | The two preserved ROM1 roots carry the same 110,632,356-byte music archive, exposing exactly 21 file nodes whose matching payloads are byte-identical 21/21. | Medium | ✔ promoted | [EXP-0229](../experiments/EXP-0229-music-runtime/) |
| VIDEO-MUSIC-002 | Every examined shipped music member is a complete ordinary PCM RIFF/WAVE with two channels, 22,050 Hz, 16-bit samples and no observed extra loop/metadata chunks. | Medium | ✔ promoted | [EXP-0229](../experiments/EXP-0229-music-runtime/) |
| VIDEO-MUSIC-003 | Shipped track duration follows PCM frame count and reproduces as `(payloadSize − 44) / 4 / 22050`; the shipped durations span 12.982404 s to 88.693016 s. | Medium | ✔ promoted | [EXP-0229](../experiments/EXP-0229-music-runtime/) |
| VIDEO-MUSIC-004 | Music reachability depends on the ordinary resource resolver and process/environment layout; presence of an archive somewhere in the install tree alone does not prove it is mounted. | High / Medium | ✔ promoted | [EXP-0229](../experiments/EXP-0229-music-runtime/) |
| VIDEO-MUSIC-005 | The recovered music requests use fixed candidate names that join to the shipped music namespace; the recovered request population accounts for the observed candidate families. | High / Medium | ✔ promoted | [EXP-0229](../experiments/EXP-0229-music-runtime/) |
| VIDEO-MUSIC-006 | Recovered music requests replace candidate **lists**, not isolated one-shot tracks. | High / Medium | ✔ promoted | [EXP-0229](../experiments/EXP-0229-music-runtime/) |
| VIDEO-MUSIC-007 | With an initialized player buffer, ordinary mode chooses a randomized candidate order and advances through it at stream end; fixed-source mode retains/reloads the selected candidate. | High | ✔ promoted | [EXP-0229](../experiments/EXP-0229-music-runtime/) |
| VIDEO-MUSIC-008 | The researched school-state bit changes the order of the two school candidates, not which of the two is permitted; ordinary playback still uses the common randomized/list progression rules. | High | ✔ promoted | [EXP-0229](../experiments/EXP-0229-music-runtime/) |
| VIDEO-MUSIC-009 | Text markup includes a numeric music-candidate selector that enters persistent fixed-source mode. | High / Medium | ✔ promoted | [EXP-0229](../experiments/EXP-0229-music-runtime/) |
| VIDEO-MUSIC-010 | Request availability, playback enable and music volume are distinct controls. | High / Unknown | ✔ promoted | [EXP-0229](../experiments/EXP-0229-music-runtime/) |
| VIDEO-MUSIC-011 | Music uses a streaming audio buffer. | High | ✔ promoted | [EXP-0229](../experiments/EXP-0229-music-runtime/) |
| VIDEO-MUSIC-012 | Archive capacity and request reachability are separate: adding a valid music payload does not by itself create a request for it. | High / Medium | ✔ promoted | [EXP-0229](../experiments/EXP-0229-music-runtime/) |

### VIDEO-MUSIC-001

- The archive's SHA-256 is
  `61b7fcd4c4515525725fa6bd45ab4bd7b84453a6a36d36639404ba10fc6dbe86` on both
  roots.
- The 21 matching payloads are byte-identical, not only name- or
  metadata-equal.
- The exact member-name inventory is corpus evidence rather than part of the
  public format grammar.

**Confidence.** Medium — an exhaustive two-root corpus result; agreement alone
does not discriminate against a different lawful release.

### VIDEO-MUSIC-002

**Confidence.** Medium — exhaustive for the examined shipped population, not a
universal WAVE restriction.

### VIDEO-MUSIC-003

No separate authored loop point was found in the examined WAVE payloads: the
absence of `smpl` and of every other non-`fmt `/`data` chunk means no shipped
WAVE supplies one.

**Confidence.** Medium — exhaustive over the shipped corpus, with the
arithmetic fixed by each payload's own header.

### VIDEO-MUSIC-004

**Confidence.** High for the resolver condition / Medium for the
preserved-layout observation.

### VIDEO-MUSIC-005

Exact literal/reference tables are private evidence.

**Confidence.** High for the recovered literal/reference population / Medium
for UI-surface interpretation.

### VIDEO-MUSIC-006

One family builds the larger mission/background list; the others build smaller
context lists. Computed indirect request paths remain possible.

**Confidence.** High within the recovered direct-call population / Medium for
named-surface interpretation.

### VIDEO-MUSIC-007

A one-entry ordinary list therefore repeats.

**Confidence.** High — competing ordinary/fixed-source models are
distinguished by the reached state machine.

### VIDEO-MUSIC-008

**Confidence.** High.

### VIDEO-MUSIC-009

No use was found in the recovered pager corpus; additional unrecovered content
remains possible.

**Confidence.** High for the parser-to-player mechanism / Medium for bounded
shipped absence.

### VIDEO-MUSIC-010

Music volume is independent of SFX and speech volume; one initializer path
remains Unknown.

**Confidence.** High for the distinct controls and ordering / Unknown for the
unresolved initializer.

### VIDEO-MUSIC-011

Stop/reload behavior depends on ordinary versus fixed-source state, and
replacement/destruction pass through that state-aware cleanup path. No stored
absolute pointer to that conditional routine exists in the parsed sections.

**Confidence.** High within the direct-call scope — both state arms, the
import, the call arguments, the refill worker and the destructor are recorded;
a computed indirect caller remains possible.

### VIDEO-MUSIC-012

Recovered live lists are fixed-name lists; runtime-computed names and targets
remain open.

**Confidence.** High for the recovered list-building and resolver mechanisms /
Medium for the negative population boundary, which does not exclude
runtime-computed names, targets or another representation.

## Sound effects

| ID | Public functional claim | Confidence | Status | Evidence |
|---|---|---|---|---|
| VIDEO-SFX-013 | The common SFX playback boundary receives an already selected sample plus volume/attenuation, pan, play/loop state, priority/category and optional frequency. | High | ✔ promoted | [EXP-0230](../experiments/EXP-0230-non-music-sfx-events/) |
| VIDEO-SFX-014 | The static research found a large bounded population of direct callers to the common SFX terminal, with some receiver sources classified and others unresolved. | High / Medium | ✔ promoted | [EXP-0230](../experiments/EXP-0230-non-music-sfx-events/) |
| VIDEO-SFX-015 | Registry capacity is data-sized while event reachability is supplied by compiled selectors, inherited class values, formulas or direct filenames. | High / Medium | ✔ promoted | [EXP-0230](../experiments/EXP-0230-non-music-sfx-events/) |
| VIDEO-SFX-016 | Recovered interface actions select fixed SFX registry slots before entering the common playback boundary. | High / Unknown | ✔ promoted | [EXP-0230](../experiments/EXP-0230-non-music-sfx-events/) |
| VIDEO-SFX-017 | Ambient sound selection is driven by visible terrain/object state plus a scheduled deadline branch. | High / Medium | ✔ promoted | [EXP-0230](../experiments/EXP-0230-non-music-sfx-events/) |
| VIDEO-SFX-018 | Unit swing, cast and hurt sounds do not share one universal five-element source: different actions use inherited class entries, formulas and a voice-bank route with throttling/gates. | High / Medium | ✔ promoted | [EXP-0230](../experiments/EXP-0230-non-music-sfx-events/) |
| VIDEO-SFX-019 | Spell SFX selection combines formula-driven choices with a delayed fixed projectile choice; not every spell is proved to execute every candidate arm. | High / Medium / Unknown | ✔ promoted | [EXP-0230](../experiments/EXP-0230-non-music-sfx-events/) |
| VIDEO-SFX-020 | The two preserved roots agree on the sparse SFX registry and on payloads reached by the recovered registry join, while the complete SFX archives are not byte-identical. | Medium | ✔ promoted | [EXP-0230](../experiments/EXP-0230-non-music-sfx-events/) |
| VIDEO-SFX-021 | A separate direct-filename SFX route exists, so a registry-only sound model is false. | High / Medium | ✔ promoted | [EXP-0230](../experiments/EXP-0230-non-music-sfx-events/) |

### VIDEO-SFX-013

Slot selection occurs upstream.

**Confidence.** High — the alternative model in which a playback argument is
itself the registry slot is excluded.

### VIDEO-SFX-014

This is not claimed as a universal runtime census: fully computed targets,
aliases derived after a bulk copy, and registry-derived receivers passed
through members or parameters all remain possible.

**Confidence.** High for the measured static populations / Medium for their
union as a runtime bound; the unresolved-receiver complement is explicitly
unresolved, not a demonstrated non-registry source.

### VIDEO-SFX-015

Adding a registry row can change an existing selector but does not itself
create a new event.

**Confidence.** High for loader/selector forms / Medium outside recovered
forms.

### VIDEO-SFX-016

Exact shipped slot-to-resource names are corpus/content evidence and are
omitted from this public edition; one selector family's common gate, update
ordering and visible transition labels remain Unknown, because only its
selector tail was recovered.

**Confidence.** High for each recovered selector and for the gates and
orderings stated outside that family / Unknown for that family's common gate,
update ordering and physical labels; none of it is inferred from the registry
path.

### VIDEO-SFX-017

One potential branch is unreachable in the examined shipped object corpus but
is not structurally impossible.

**Confidence.** High for control flow/selectors/timing / Medium for shipped
non-reachability.

### VIDEO-SFX-018

**Confidence.** High for the reached selector/gate logic / Medium for the
exhaustive shipped class-table join.

### VIDEO-SFX-019

Ordering of several retained client tails remains Unknown.

**Confidence.** High for the recovered selectors/gates / Medium for corpus
joins / Unknown for the stated ordering gaps.

### VIDEO-SFX-020

Exact missing-slot and filename inventories are private corpus evidence.

**Confidence.** Medium — exhaustive for two preserved corpora only.

### VIDEO-SFX-021

Filename spelling alone does not identify which producer route owns an event,
and computed names remain open.

**Confidence.** High for the recovered literal/reference population / Medium
for bounded corpus join.

## Cutscenes

| ID | Public functional claim | Confidence | Status | Evidence |
|---|---|---|---|---|
| VIDEO-029 | Startup selects one of two cutscene namespaces from startup/environment media-speed state. | High | ✔ promoted | [EXP-0264](../experiments/EXP-0264-cutscene-smacker-abi/) |
| VIDEO-030 | When the cutscene gate is active, a mission request probes numbered `.smk` members in increasing two-digit order within the selected namespace. | High / Medium | ✔ promoted | [EXP-0264](../experiments/EXP-0264-cutscene-smacker-abi/) |
| VIDEO-031 | The sidecar filename is derived by replacing the movie's **last** extension with `.reg`, preserving earlier dots. | High / Medium / Unknown | ✔ promoted (amended, partially retracted) | [EXP-0264](../experiments/EXP-0264-cutscene-smacker-abi/) |
| VIDEO-032 | Cutscene sidecars provide initial position plus frame-bounded fade and pan records. | High / Medium | ✔ promoted | [EXP-0264](../experiments/EXP-0264-cutscene-smacker-abi/) |
| VIDEO-033 | ROM1 links a bundled Smacker decoder through a fixed imported API surface. | High / Unknown | ✔ promoted | [EXP-0264](../experiments/EXP-0264-cutscene-smacker-abi/) |
| VIDEO-034 | Game-side cutscene setup selects the movie resource, opens decoder state, obtains dimensions, configures a destination/blitter path and enables decoder sound after successful setup. | High / Unknown | ✔ promoted | [EXP-0264](../experiments/EXP-0264-cutscene-smacker-abi/) |
| VIDEO-035 | A game-side frame step waits for readiness/focus, updates sidecar fade/pan state, decodes/presents the frame, checks termination state and advances. | High / Unknown | ✔ promoted | [EXP-0264](../experiments/EXP-0264-cutscene-smacker-abi/) |
| VIDEO-036 | Ordinary key/system-key down, mouse-button down, close and quit messages stop the current numbered cutscene scan. | High / Unknown | ✔ promoted | [EXP-0264](../experiments/EXP-0264-cutscene-smacker-abi/) |

### VIDEO-029

Both archive families may be attempted, but the numbered mission route uses one
selected namespace and has no general fallback to the other.

**Confidence.** High for the recovered branch logic / runtime mount success
outside scope.

### VIDEO-030

Missing/open-failure continuation and user-stop termination are distinct; the
gate producer and visible failure transitions remain Unknown.

**Confidence.** High for gate/format/bounds/branch behavior / Medium for
examined archive population.

### VIDEO-031

The sidecar is loaded before the original movie resource is opened. Exact
shipped sidecar hashes/values are private evidence.

**Confidence.** High for filename derivation and load order / Medium for corpus
equality / Unknown for missing-sidecar presentation.

**Amended.** The former first-dot reading is refuted: EXP-0264's correction in
[`retracted.md`](retracted.md) replaces it with the last-dot rule stated
above. The sidecar-before-open order and the corpus equality stand.

### VIDEO-032

Fade state affects palette presentation; pan state changes presentation
coordinates across the reached frame interval.

**Confidence.** High for reader-to-consumer data flow / Medium for two-root
sidecar census.

### VIDEO-033

The public edition records only this game-side dependency boundary;
export-body/internal decoder details require separate third-party review.

**Confidence.** High within the parsed import/export boundary / Unknown for
computed dynamic invocation.

### VIDEO-034

Internal decoder structure offsets are intentionally omitted.

**Confidence.** High for the reached ROM1-side call order / decoder-internal
representation remains Unknown; decoder internals are also outside the public
scope.

### VIDEO-035

Exact callback meaning, pixels, audio latency and physical cadence remain
Unknown.

**Confidence.** High for recovered ROM1-side ordering / Unknown for physical
presentation and opaque decoder/host semantics.

### VIDEO-036

Natural completion and some construction/open failures follow the continuation
path.

**Confidence.** High for message/return branches / Unknown for physical error
presentation.

## Cutscene decoder boundary

| ID | Public functional claim | Confidence | Status | Evidence |
|---|---|---|---|---|
| VIDEO-045 | The original decoder integration can consume a game-provided resource source rather than requiring a standalone filename. | High | ✔ promoted (amended, partially retracted) | [EXP-0266](../experiments/EXP-0266-smacker-buffer-contract/) |
| VIDEO-046 | For the borrowed-resource route, the source's current origin determines where movie decoding begins. | High / Unknown | ✔ promoted | [EXP-0266](../experiments/EXP-0266-smacker-buffer-contract/) |
| VIDEO-047 | Decoder setup configures a borrowed game destination and frame decode writes into that destination. | High / Unknown | ✔ promoted | [EXP-0266](../experiments/EXP-0266-smacker-buffer-contract/) |
| VIDEO-048 | The researched decoder supports indexed palette output and additional packed-color modes; ROM1's reached caller uses indexed output. | High / Medium / Unknown | ✔ promoted | [EXP-0266](../experiments/EXP-0266-smacker-buffer-contract/) |
| VIDEO-049 | Frame decode and frame advance are separate operations. | High / Unknown | ✔ promoted | [EXP-0266](../experiments/EXP-0266-smacker-buffer-contract/) |
| VIDEO-050 | Decoder readiness/wait behavior has clock- and sound-progress-dependent paths rather than being a single fixed sleep. | High / Unknown | ✔ promoted | [EXP-0266](../experiments/EXP-0266-smacker-buffer-contract/) |
| VIDEO-051 | Decoder close respects ownership of a borrowed game resource source and configured destination; separate game/buffer cleanup owns other resources. | High / Unknown | ✔ promoted | [EXP-0266](../experiments/EXP-0266-smacker-buffer-contract/) |

### VIDEO-045

The public claim is limited to that ownership/input distinction;
bundled-decoder ABI flag values and instruction evidence are private. These
branches do not establish a general malformed-input contract.

**Confidence.** High — the recovered callee body, restored import-name map and
read-helper slice distinguish text, handle and caller-state alternatives, and
native nonzero-offset handle input matches filename input.

**Amended.** The allocation-selector clause is refuted: EXP-0266's correction
in [`retracted.md`](retracted.md) separates the two selectors that clause had
joined. ROM1's borrowed-handle route stands.

### VIDEO-046

Decoder buffering policy is distinct from a compressed-payload length.

**Confidence.** High for the bounded discriminator / I/O-failure and broader
ownership paths Unknown.

### VIDEO-047

ROM1 uses the ordinary indexed-output path; exact decoder descriptor
offsets/mode flags are intentionally omitted.

**Confidence.** High for reached setup/decode behavior / unsupported output
modes and extents Unknown.

### VIDEO-048

Palette changes are exposed to the game-side presentation path. Exact
decoder-internal tables/flags are private third-party evidence.

**Confidence.** High for the reached mode distinction, the retained layout
measurements and the bounded differential sample / Medium for EN/RU sample
parity / Unknown for the isolated third flag value, for all palette-change
sequences and for other video populations.

### VIDEO-049

Decoder-local wrapping behavior does not itself define ROM1's cutscene end
condition; the game-side caller owns the stopping bound.

**Confidence.** High for the retained increment, wrap and latch branches and
for the observed no-advance results / Unknown for native terminal/ring-frame
behaviour and for complete drop-helper semantics.

### VIDEO-050

The public claim does not expose decoder state offsets or callback internals.
The wait helper's own latency is outside the retained proof, so no nonblocking
or wall-clock playback guarantee follows.

**Confidence.** High for the bounded clock arithmetic, branches and return
roles / Unknown for audible timing, wraparound behaviour, sound-callback
semantics and observed ROM1 cadence; no native wait call was made.

### VIDEO-051

Exact allocator flags, structure sizes and internal free lists are private
third-party evidence. In the bounded native sample the borrowed source remains
seekable and the destination bytes remain readable and unchanged after decoder
close.

**Confidence.** High for the bounded ownership and survival checks, which the
native survival results corroborate on the named sample / Unknown for decoding
after destination release, for concurrent calls and for all error-path
lifetimes.

## Option controls and consumers

| ID | Public functional claim | Confidence | Status | Evidence |
|---|---|---|---|---|
| VIDEO-SFX-053 | The recovered Sound Options Acknowledgments checkbox imports the stored Acknowledgement value. | High / Medium | ✔ promoted | [EXP-0364](../experiments/EXP-0364-acknowledgment/) |
| VIDEO-SFX-054 | In the recovered selected-unit response path, a zero Acknowledgement value rejects the response chooser before candidate collection; nonzero values pass that gate. | High / Medium / Unknown | ✔ promoted | [EXP-0364](../experiments/EXP-0364-acknowledgment/) |
| VIDEO-MUSIC-056 | With a playback buffer present, the recovered Random Order setter selects ordinary mode and initializes the candidate permutation0..n-1. | High / Medium / Unknown | ✔ promoted | [EXP-0366](../experiments/EXP-0366-option-consumers/), derived instructions, controlled measurements and two branch mutations |
| VIDEO-OPTIONS-057 | The recovered Sound Options list walks the player current candidate bank and looks up normalized names after their six-character resource prefix in the dictionary loaded from main/text/tunes.txt. | High / Medium / Unknown | ✔ promoted | [EXP-0366](../experiments/EXP-0366-option-consumers/), derived instructions, controlled measurements and two branch mutations |

### VIDEO-SFX-053

A separate dialog message exports the selected control value to that cell
before forwarding the close message. Changing the selected control value alone
does not update the cell in the isolated probe; the close message alone does
not export it. This does not establish every control transaction or native
registry persistence.

**Confidence.** High for the joined copy methods and discriminated message
arms; Medium for their original dialog interpretation.

### VIDEO-SFX-054

The tested caller sends its command before choosing a response, so disabling
the option leaves that command-service event intact. The concrete unit response
method separately refuses elapsed time below 3000 milliseconds since its
preceding response timestamp. Conditional original-instruction execution
discriminates zero/nonzero, the timing boundary and two mutations; it does not
establish all voice-request families or audible playback.

**Confidence.** High for the reached instruction conditions and request
boundary; Medium for the selected original command-path interpretation.
Dialogue, damage/death audio and unexamined callers remain Unknown.

### VIDEO-MUSIC-056

Disabled retains that sequential order. Enabled performs n swaps, each using
two CRT rand()%n indices. It does not select fixed-source mode. An absent buffer
leaves the measured mode, flag and order unchanged. The Sound Options checkbox
writes its configuration value and invokes this setter immediately.

**Confidence.** High for the original setter, control arm and controlled draws;
Medium for dialog interpretation. Native RNG seed, distribution, complete
context activation and listening remain Unknown.

### VIDEO-OPTIONS-057

Selecting a list row only changes the selected index. Play enables music,
reloads only when the selected index differs, then applies gain and starts.
Stop disables music and invokes the state-specific stop or transition method.
These controls operate separately from channel volume.

**Confidence.** High for index walk and measured selection/Play/Stop event
arms; Medium for the title dictionary and dialog interpretation. CString
normalization is a named substitute; Unicode, missing-title behavior, native
list population, fade audibility and registry lifetime remain Unknown.

## Character-generator sounds

| ID | Public functional claim | Confidence | Status | Evidence |
|---|---|---|---|---|
| VIDEO-SFX-058 | On the pre-create page each left press restarts a difficulty button's `level1..3.wav` or a hero button's `char.wav`; Back, Escape, and OK or Enter with a non-empty name request `ok.wav` just before the page's close stops it. | High / Unknown | ✔ promoted | [EXP-0407](../experiments/EXP-0407-chargen-sounds/) |
| VIDEO-SFX-059 | On the detailed page each applied statistic step restarts `+_-.wav`, also on double-click and held-button repeat; a skill press requests its class's member for that skill slot unless it plays, as the school room does. | High / Medium / Unknown | ✔ promoted | [EXP-0407](../experiments/EXP-0407-chargen-sounds/) |
| VIDEO-SFX-060 | Every `chrgen` request found passes SFX volume, pan 0, no loop and priority 128; at most one instance of each member plays, and members share 16 channels, displacing only a lower-priority sound when all are busy. | High / Medium / Unknown | ✔ promoted | [EXP-0407](../experiments/EXP-0407-chargen-sounds/) |

### VIDEO-SFX-058

- A press on difficulty button 0, 1 or 2 (`UNIT-GATE-014`) stops a playing
  `level1.wav`, `level2.wav` or `level3.wav`, rewinds it and requests it again,
  whether or not that button was already selected.
- A press on any of the four hero buttons (`SESS-HERO-013`) does the same with
  `char.wav`. A comparison with the current hero guards only the name-field
  update; the request follows on both branches.
- OK and Enter with a non-empty name (`TEXT-CHARGEN-028`), and the amulet Back
  control and Escape (`HERO-CHARGEN-083`), each request their own `ok.wav`
  sample unless it plays and then send the page its continue or Back message.
  With an empty name, OK and Enter request nothing and send nothing.
  The open page closes inside that call: the close stops the new instance and
  deletes the page's samples before the press or key routine returns. Escape
  reaches the page through the frame's key forwarding after the
  character-generation transition (`MENU-ESC-010`, `TOWN-373`).
- Motion, release, the right button and a double-click's second click request
  no `chrgen` member, and the page's open routine makes no playback call.
  Typing the name requests non-`chrgen` letter sounds.

**Confidence.** High for the controls, events, arguments and replay rules:
every arm is read with its switch table, and the only selection comparison
guards other work. High that the close stops `ok.wav` within the same call on a
page whose open ran; the open sets the page's open flag on every path.

**Unknown.** Whether any of that `ok.wav` instance is audible. A press reaching
the page while its open flag is clear was not searched for.

### VIDEO-SFX-059

- An increase applies only when its cost fits the free points and the value is
  below 45; a decrease only when the value is above 15. Each applied step stops
  and rewinds a playing `+_-.wav` and requests it; a refused step is silent.
- The statistic panel treats a double-click's second click and the held-button
  repeat as presses. While the left button stays down, the first cursor tick
  more than 150 ms after the last mouse message posts the repeat, and each
  later tick more than 66 ms after the tick that posted the previous repeat
  posts it again; any mouse message restores the 150 ms delay. The other
  panels and pages ignore both.
- A skill press selects index 0..4 and requests the member loaded at that index
  unless it plays: `fsword`, `faxe`, `fclub`, `fpike`, `fbow`, or with the
  hero's class bit set `mfire`, `mwater`, `mair`, `mearth`, `mastral`. It stores
  index + 1 as the draft skill (`SAV-934`); no comparison with the previous
  selection guards the request.
- The school room requests the same ten members on a left press in a skill
  column unless playing, also when the press deselects the column. Its member
  follows the stored skill slot (`TOWN-GENERAL-106`) as on the detailed page:
  1 sword or fire, 2 axe or water, 3 club or air, 4 pike or earth, 5 bow or
  astral.
- Motion and release request no `chrgen` member on the page, its panels or the
  school room.

**Confidence.** High for the step conditions, event routing, index-to-member
and slot-to-member mapping and replay rules. Medium for the on-screen position
of each skill index and school column: the hit masks were not rendered.

**Unknown.** The cursor tick rate and the system double-click interval, which
set the repeat cadence and decide when a second click counts as a
double-click. The meaning of the school room field that blocks its press when
non-zero.

### VIDEO-SFX-060

- The request goes through the common SFX boundary (`VIDEO-SFX-013`). It loads
  the sample on first use, takes the first of its duplicate buffers that is not
  playing and needs a channel: the first idle one of 16, or, when all 16 play,
  the one holding the lowest priority below 128, which is stopped. Otherwise
  the request is dropped, so a `chrgen` request never displaces another at 128.
- Each `chrgen` caller either restarts or skips a playing instance of its own
  member, so at most one instance of each member plays; different members
  overlap.
- The SFX volume option only sets the volume value; no `chrgen` request tests
  it or an enable flag. A static initializer sets it to −700 at process start;
  a settings load replaces it with the registry value `SoundSfxPos` when
  `HKLM\SOFTWARE\1C\Allods` opens; the option's control stores −(p − m)²/m,
  truncated, for control position p and control value m.
- `ok.wav` is also requested by a press on any of the main menu's eight buttons
  and on the Hall of Fame OK, each skipping a playing instance.
- The searched population is every reference to the 16 literals and every
  playback call in the code of the five owning surfaces; 80 playback calls
  elsewhere were not attributed.

**Confidence.** High for the arguments, the buffer and channel rules and the
channel count, set once from a constant 16. High for the static initializer
and the registry load of the volume value, each read whole. Medium that no
other code requests these members: a sample pointer copied out of an unswept
field is not excluded.

**Unknown.** The runtime volume value: whether the registry holds
`SoundSfxPos`, and any other write through the settings object's address.
Mixing, device behaviour and audibility.

## Music at mission loss and campaign LOAD

Terms. The player is the object at `frame+0xc8`. The music gate is `0x005eb478`, the same gate
`VIDEO-MUSIC-006` reads. The stop is `004530f3`, the replace is `00452c34` and the start is
`00453089`. The mask is `frame+0x3dc`. Window messages are posted through one import slot,
`[0x632f5c]`, so every hop below is queued, not nested.

| ID | Public functional claim | Confidence | Status | Evidence |
|---|---|---|---|---|
| VIDEO-MUSIC-061 | No arm read on the mission-lost panel's display path and no direct call from it reaches a music routine (126, 162 and 62 computed-call functions unread); its one sound call is a fixed SFX. | High / Medium | ● active | [EXP-0426](../experiments/EXP-0426-loss-load-music/) |
| VIDEO-MUSIC-062 | Exit to Main Menu stops the player through the teardown, then requests the one-entry menu list from the `0x421` arm; both happen after the choice, not at display. | High / Medium | ● active | [EXP-0426](../experiments/EXP-0426-loss-load-music/) |
| VIDEO-MUSIC-063 | In the arms read, Load Game from the failure panel makes no music call when chosen or when the save dialog opens; the stop comes at selection, and cancelling takes the Exit route. | High / Medium | ● active | [EXP-0426](../experiments/EXP-0426-loss-load-music/) |
| VIDEO-MUSIC-064 | LOAD of a mission save stops the player, then, if `00477c00` reaches `00478956` with the gate set, requests the twelve-entry `B00`..`B11` list and starts it. | High / Medium | ● active | [EXP-0426](../experiments/EXP-0426-loss-load-music/) |
| VIDEO-MUSIC-065 | LOAD of a town save makes no music request in the load routine; it posts `0x42e`, whose arm requests the one-entry Town list. | High / Medium | ● active | [EXP-0426](../experiments/EXP-0426-loss-load-music/) |
| VIDEO-MUSIC-066 | The menu and Town owners skip the replace when the player's list already starts with their entry; the `B00`..`B11` owner never compares, so a mission LOAD that reaches it re-requests. | High / Medium | ● active | [EXP-0426](../experiments/EXP-0426-loss-load-music/) |

### VIDEO-MUSIC-061

- Order, EN and RU executables byte-identical: the client `0xb4` arm at `00410623`
  posts `0x433` with wParam `0xff`; the `0x433` arm at `004738bb` stores
  `frame+0x414 = 0xff` and posts `0x431`; the `0x431` arm at `00474140`, taken only while
  latch `frame+0x3c0` is zero, sets the latch, constructs `00446de2`, stores the panel at
  `frame+0x110`, shows it with `00476810` and then calls the SFX terminal `00453b08`.
- `00476810` sets mask bit 8 (bit `0x8000` for a panel stored at `frame+0x3a8`), remaps
  the screen and calls virtual slots of the panel. It calls no music routine.
- The SFX call carries a fixed registry slot, 16, in `EXP-0230` (event `mission-failed`,
  `VIDEO-SFX-013`, `VIDEO-SFX-014`). It does not call the stop, the replace or the start.
- Direct-call closures from the reporter `004d8963` (686 functions), the client
  dispatcher `004104e8` (574), the panel constructor (59), the show routine (94) and the
  panel's vtable slots `+0x48`, `+0x78`, `+0x80`, `+0x84` reach none of the twenty-four
  music routines marked in `evidence/closures.tsv`.
- The stop, the replace, the start and fourteen owner or teardown entries have no stored
  address in any section of either executable and no `E9` transfer; each has only direct
  `E8` references (`evidence/scan-en.tsv`, `evidence/scan-ru.tsv`, identical).
- Bound: no direct-call path from the reporter, the client dispatcher, the panel constructor, the show routine or the panel's vtable slots reaches a music routine, and none of the arms read calls one. The computed-call functions unread are 126 in the reporter closure, 162 in the `0xb4` handler closure and 62 in the SFX terminal closure.
- The player keeps its previous list and state at display. Where a track ends during the
  panel, `VIDEO-MUSIC-007`'s end-of-stream advance applies; it is not a loss request.

**Confidence.** High for the arm-by-arm reads and the three `E8` posts. Medium that no
request exists between the reporter and the panel: the closures hold 126 functions with
a computed call that were not read, and the first `0xb4` hop through the session queue is
inherited from `MISSION-DEFEAT-046`.

**Unknown.** Whether the stream worker keeps running while the panel is up. Which sample
SFX slot 16 plays and whether it is audible over the music.

### VIDEO-MUSIC-062

- The panel result `0x445` makes the `+0x110` close arm post `0x41e`. The `0x41e` arm at
  `004745fe` calls teardown `00479dd0`, then posts `0x421` (`MISSION-DEFEAT-046`).
- `00479dd0` calls `00478a80` only when mask bit 0 is set. `00478a80` calls the stop on
  the player at `00478a9e` when the gate is non-zero, then clears mask bit 0.
- The stop takes no argument. With a buffer at `player+0x9c` and a non-empty list at
  `player+0x10`, ordinary mode (`player+0x18 == 0`)
  stops the buffer through the DirectSound buffer's stop slot and sets `player+0xc = 1`;
  fixed-source mode reloads the selected candidate instead (`VIDEO-MUSIC-007`).
- The `0x421` arm at `00473a76` calls the menu owner `004769c0` only when the mask is zero.
  The owner requests a one-entry list, `music\menu.wav`, with the replace at `00476bc3` and
  starts at `00476bdb`, when the gate is non-zero.
- The menu owner skips the replace when the player's list is non-empty and its first entry
  equals `music\menu.wav` (`VIDEO-MUSIC-066`); the start still runs.

**Confidence.** High for the call order and the gates, and for the read fact that the
`+0x110` close arm clears mask bit 8 before it posts `0x41e` (`evidence/listing.txt`). Medium for the
remaining mask bits: that `frame+0x3dc` is exactly bit 0 on the mission screen before the panel (and so is
zero after teardown) rests on the mission-screen state in `MISSION-STOP-016`, not on a
witnessed value.

**Unknown.** Any other bit set in the mask at the failure, which would skip the menu request.
The audible gap between the stop and the menu start.

### VIDEO-MUSIC-063

- Result `0x446` makes the `+0x110` close arm post `0x418`. The `0x418` arm at `00473468`
  builds the save-selection dialog, stores it at `frame+0x128` and shows it with `00476810`.
  It calls no music routine in the arm read (`evidence/listing.txt`).
- A selected entry makes the `+0x128` close arm post `0x419`. The `0x419` arm at
  `004734cc` calls `00479dd0` only when the mask equals 1, which stops the player as in
  `VIDEO-MUSIC-062`, then calls `00476340`, and for mask zero sends command `0x445` to
  `frame+0x100` and calls the load path `00478af0`. The mask is zero after a taken teardown.
- A cancelled dialog (result other than `0x445`) posts `0x41e` when the session pointer
  `0x005cd758` is non-zero and `frame+0x414 == 0xff`, which the failure arm stored. It then
  follows `VIDEO-MUSIC-062` exactly.

**Confidence.** High for the arm reads and the exclusive close chain. Medium that the mask
equals 1 when the dialog closes (see `VIDEO-MUSIC-062`).

**Unknown.** The result code of a dialog dismissed by other means than the two buttons.

### VIDEO-MUSIC-064

- The load routine `00478af0` reads `CurrentState`/`InBattle` with default 1 (`SAV-914`).
  A non-zero value calls `00477c00` with argument 1 at `00478bba`.
- `00477c00`, when the gate is non-zero, calls the stop on the player at `00477c4d`
  before it rebuilds the session. Between the stop and `00478956` a wait loop
  (`00477e84..00477f1c`) has three exits that return zero without a request: a pumped message
  of id `0x12` (`00477ea2`), a 60,000 ms timeout (`00477edb`, to `00478a64`) and a zero
  result of `004104e8(0x64)` (`00477f08`). Only if the routine reaches `00478956`, and the
  phase is not 3, does it call the list owner `0047cc00`, which itself requires the gate.
- `0047cc00`, when the gate is non-zero, builds one list of twelve strings,
  `music\B00.wav` through `music\B11.wav` in that order, replaces the player's list at
  `0047cd1d` and starts at `0047cd35`. Argument 1 suppresses the map load and the `0x442`
  post (`SESS-START-036`).
- Direct closure from `00478af0` reaches the stop, the replace and the start only through
  `00477c00` and `0047cc00` (`evidence/closures.tsv`). Nothing posted by this branch
  requests music.
- The replace is the player's usual one: a random initial candidate and ordinary progression
  (`VIDEO-MUSIC-007`).

**Confidence.** High for the call order, the three exits and the arguments. A failed mission-start
exit therefore leaves the player stopped with no replacement. Medium for the closure's absence
clause: 400 functions of the load closure hold a computed call that was not read.

**Unknown.** The document loader's own effect on a running player beyond its direct
closure, which reached no music routine.

### VIDEO-MUSIC-065

- A zero `InBattle` value takes the town branch of `00478af0`: `0041db20`, `004d88f1`,
  `004104e8(100)`, `00477650`, `0041da75`, `004d88f1`, then post `0x42e`
  (`SHOP-TOWN-022`). Their direct closures reach no music routine.
- The `0x42e` arm at `00474ce3` calls `00477130`, which is a surface transition (`TOWN-372`) and, when the
  gate is non-zero, requests the one-entry list `music\Town.wav` with the replace at
  `0047720f` and starts at `00477227`.
- The `0x42e` arm runs from the message queue after the load routine returns. The difference
  from `VIDEO-MUSIC-064` is therefore a post, not a synchronous call.
- The Town owner skips the replace when the player's list is non-empty and starts with
  `music\Town.wav` (`VIDEO-MUSIC-066`).
- No stop is issued before the replace on this branch except the one the replace routine
  itself runs first (`VIDEO-MUSIC-006`); with a skipped replace the running track continues.

**Confidence.** High for the order and arguments. Medium for the absence of any other
request: 234 functions of the `004d88f1` closure alone hold a computed call that was not read.

**Unknown.** Whether a later message in a town session replaces the list again before the
player hears the Town list.

### VIDEO-MUSIC-066

- The menu owner `004769c0` and the Town owner `00477130` call string comparison `00553d80`
  on the first entry of the player's list against their own literal and request the replace
  only when the list is empty or the strings differ.
- Owner `0047cc00` does no comparison and replaces unconditionally whenever the gate is non-zero.
- Consequences: a mission LOAD that reaches `00478956` with the gate set re-requests the list
  even when the `B` list is already playing; a town LOAD while the Town list plays keeps the
  running track, since that branch has no preceding stop. On Exit the teardown stop runs first
  (`VIDEO-MUSIC-062`), so the skipped menu replace avoids a list rebuild, not the stop.
- The comparison reads the first array element. Whether the random initial selection
  (`VIDEO-MUSIC-007`) reorders that array is not read, which is immaterial for the two
  one-entry lists these owners build.

**Confidence.** High for the two comparisons and the unconditional `B` request. Medium for
the equality semantics of `00553d80`, read as a zero-on-equal byte comparison with a locale branch
when `[0x00630854]` is non-zero.

**Unknown.** The comparison's case handling under that locale branch, and the other six
list owners' behaviour, which this experiment did not read.
