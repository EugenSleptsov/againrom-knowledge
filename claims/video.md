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

## Map command voices

| ID | Public functional claim | Confidence | Status | Evidence |
|---|---|---|---|---|
| VIDEO-067 | Ten builder tails and the selection tail call the speaker chooser, then one voice reader: move, attack, swarm, patrol, town `vt+0x6c`; guard, stand ground, defend `+0x70`; retreat `+0x74`; pickup `+0x7c`; selection `+0x78`; cast none. | High / Medium / Unknown | ● active | [EXP-0445](../experiments/EXP-0445-command-voices/) |
| VIDEO-068 | The speaker is one random member of the highest non-empty tier: hero-shaped (`+0x18c` bit 0x1), else armed human, else unarmed human; a member needs `+0x7c` set and health `+0xfc` above 0; no owner, distance or visibility field is read. | High / Medium / Unknown | ● active | [EXP-0445](../experiments/EXP-0445-command-voices/) |
| VIDEO-069 | The selection reply is `vt+0x78` (`select1` or `select2`, 2000 ms), played by the chooser speaker when Shift is up and the summary bits allow it; it shares the stamp `+0x190` with every other voice. | High / Medium / Unknown | ● active | [EXP-0445](../experiments/EXP-0445-command-voices/) |
| VIDEO-070 | Guard, Stand Ground and the defend order play the speaker bank's `defend` slot (`+0x1c`), Retreat plays `retreat` (`+0x18`) and a pickup order plays `idle` (`+0x20`), each behind the shared 3000 ms stamp. | High / Medium | ● active | [EXP-0445](../experiments/EXP-0445-command-voices/) |

### VIDEO-067

```
gesture (builder)            order  chooser call  reader   bank slot
move       (FUN_0041b620)    0x16   0041b789      +0x6c    command1..3 / defend
attack     (FUN_0041b7b1)    0x19   0041b8f7      +0x6c    command1..3 / defend
swarm      (FUN_0041b916)    0x1a   0041ba67      +0x6c    command1..3 / defend
patrol     (FUN_0041bf6c)    0x1d   0041c0bd      +0x6c    command1..3 / defend
town       (FUN_0041c796)    0x24   0041c8db      +0x6c    command1..3 / defend
guard      (FUN_0041ba8f)    0x17   0041bbc9      +0x70    defend
stand grd  (FUN_0041bbe6)    0x18   0041bd20      +0x70    defend
defend     (FUN_0041bd3d)    0x1b   0041bf4d      +0x70    defend
retreat    (FUN_0041c63f)    0x14   0041c779      +0x74    retreat
pickup     (FUN_0041c4be)    0x21   0041c620      +0x7c    idle
selection  (FUN_0041a2d5)    -      0041a9da      +0x78    select1 / select2
cast       (FUN_0041c0dc, FUN_0041c2d6)  0x1f 0x26 0x25 0x1e   none
```

- **Chooser callers.** `FUN_00422b5e` has exactly eleven direct callers, the eleven rows above with a call (a byte scan for `E8` and the decoded sweep agree, 0 other sites; `evidence/xref-voice.txt`). Each builder stores its order byte at `[edx+9]` (`evidence/opcode-stores.txt`), sends the order (`FUN_004e74fe`) and then calls the chooser, so a refused reply never holds back the order. The reader of each tail is the `CALL [reg+slot]` after the chooser's null test.
- **Cast.** The two cast builders hold no chooser call; their only indirect calls are `vt+0x7c` of the object `FUN_00573196` returns (`evidence/calls-cast-builders.txt`).
- **Input surfaces.** The map click `FUN_00419ec1` reaches the builders by cursor (`AI-CLICK-050`); the minimap handler `FUN_0048fdb0` reaches move, attack, swarm, defend and patrol (`AI-MINIMAP-062`); the command panel `FUN_0041b439`, which has 14 direct callers, reaches guard, stand ground and retreat through the jump table `0x0041b608` for buttons 3, 7 and 8 (`AI-PANEL-123`). The three panel builders have no other caller.
- **Draw-state gate.** Move and swarm, and only these, skip the reader when the chosen unit's `+0x74` equals 1 (`0041b79a`, `0041ba78`); `ANIM-STATE-002` reads draw state 1 as the move action. The chooser has already picked the unit, so a moving speaker silences the reply and no second unit is tried.
- **Executed.** Each tail, started at its order push and stopped before its epilogue, with one selected unit in each of the hero, hero-shaped mage, mercenary and peasant banks: one command-service call, then the slot of the table. Acknowledgement 0 removes every reply and keeps the command; a clock of 2500 ms (stamp 0) removes every reply; `+0x74 = 1` removes move and swarm only.
- **Reader census.** `evidence/vcalls-census.txt` lists the 340 indirect calls whose displacement is `0x6c`, `0x70`, `0x74`, `0x78` or `0x7c`. 224 follow a call to `FUN_00573196` (222 of them `vt+0x7c` of the object it returns). 11 are the tails above. 94 have a `PUSH` within the six preceding instructions; the five routines take no stack argument (`RET` without a count). The remaining 11 (`evidence/vcalls-classification.txt`) include four user-interface routines that compare or store the returned value (`004abeb8`, `004abecb`, `004ad7f7`, `004ad80a`), the routine `004f38b3` that stores the result, and `00446d2c`, whose receiver is a child-control lookup.

**Confidence.** **High** for the eleven chooser callers, their orders, readers and slots: decoded whole and executed on both roots with a sentinel per slot. **Medium** that no other path reaches the five slots on a unit: the census sorts the other sites by window and by use of the result, not by receiver proof. **Unknown** the receivers of `00444305` and of the four sites in `0x00573..0x00583` (`0057315d`, `0057b4f9`, `005801db`, `0058309b`), which were not read.

### VIDEO-068

- **Gates.** `FUN_00422b5e` returns 0 unless the session's `+0x6bc` is 2 (`00422bc5`, the campaign; `DLG-ENTRY-016`) and the Acknowledgement cell is non-zero (`00422bd5`; `VIDEO-SFX-054`).
- **Population.** It walks the selection map `view+0x9b8` (`AI-SELECT-065`) from the first bucket in bucket-chain order. For each member `unit = value`:
  - `[unit+0x7c]` must be non-zero (`00422d46`);
  - `MOVSX word [unit+0xfc]` must be above 0 (`00422d65`, `00422d6e`), so health 0 and every negative value are excluded;
  - `+0x18c & 1` appends the unit to list A (`00422d77`..`00422d9b`);
  - otherwise, while A is empty, `+0x18c & 0x10` appends it to list B when `+0x15c` is non-zero (`00422dc7`..`00422de7`), else to list C while B is empty (`00422dee`..`00422e17`);
  - a member with neither bit is dropped.
- **Pick.** After the walk: list A if non-empty (`00422e2a`), else list B (`00422fe7`), else list C (`004230ca`), else 0. The index is `rand() * n / 0x7fff` (`00422e45`..`00422e52`, repeated at `00423002` and `004230e1`), near-uniform (counts differ by at most two of the 32768 draws; draw 32767 selects the slot one past the last member). Without a hero-shaped unit the order is armed humans, then unarmed humans. A and B or C members speak from the bank `HERO-APPEAR-055` gives them: the hero or mage bank for A, the mercenary bank for B, the peasant bank for C, the mage bank for a B or C drawable with bit 0x2.
- **Fields read.** The only unit displacements the chooser reads, over all 20 executed cases and in the listing, are `+0x7c`, `+0xfc`, `+0x15c` and `+0x18c`; it reads no owner, player, position, distance or visibility field and writes nothing (`evidence/probe.json`, `access`). The stamp is not read, so a unit on cooldown can be picked and then stay silent; no second candidate is chosen.
- **Executed.** Hero before or after an armed and an unarmed human: the hero for every draw. Armed before or after unarmed: the armed one. A monster alone (neither bit): no speaker and no `rand()` call. A mage bit with bit 0x10 and no hero bit: tier B or C by `+0x15c`. Three heroes: draws 0, 10922 | 10923, 21845 | 21846 and 32766 give members 0, 0 | 1, 1 | 2 and 2. Health 0 and -10 excluded, 1 included. `+0x7c = 0` excluded. Empty map, Acknowledgement 0 and session mode 1: no speaker.
- **Ownership.** The chooser does not test it. Upstream, the command panel is disabled for a selection whose primary object another player owns (summary bit `0x4`, `AI-PANEL-061`) and the hover cursor gives select or default when that bit is set (`AI-CURSOR-226`); a click selection can hold a foreign unit (`AI-SELECT-122`).

**Confidence.** **High** for the gates, per-member tests, tier order and index formula: the routine is read whole and executed on both roots over 20 discriminating cases, including both list orders. **Medium** for the meaning of the bits and the flag: `+0x18c` bit 0x1 is the hero-shaped drawable (`HERO-APPEAR-041`), bit 0x10 follows a wire class below 0x1a (`ANIM-096`), and `+0x7c` as "selected" is inferred (`AI-SELECT-065`). **Unknown** the content of the list slot one past the last element, which a draw of 32767 (1 in 32768) selects; the executed stub returns 0 there. **Unknown** whether any order can be issued from a foreign selection, and whether a structure can be a member that reaches the chooser.

### VIDEO-069

- **Caller.** The only call of `vt+0x78` on a unit is `0041a9f0`, after the chooser call `0041a9da`, at the end of `FUN_0041a2d5`. Its three callers (`00419f88`, `0041a0f6`, `0041a2ca`) are in the map-click handler `FUN_00419ec1`: the click with an empty selection, the select cursor arm and the drag-rectangle branch (`AI-CLICK-050`). The key and digit group selections are not callers.
- **Gates.** After the selection summary is rebuilt (`0041a8f0`), the tail tests: the Shift latch `[0x005eb55c]` clear (`0041a996`, `AI-KEYMOD-059`), the selection count `view+0x140` non-zero (`0041a9a5`), `view+0x144 & 1` (`0041a9ba`) and `view+0x144 & 4` clear (`0041a9cd`). The chooser's own gates of `VIDEO-068` follow.
- **Executed.** All 16 combinations of Shift 0 or 1, count 0 or 1 and flags 0, 1, 4 and 5: only Shift 0, count 1, flags 1 plays (`select1` for a hero at draw 0).
- **Sharing.** The reader is `FUN_0045e9c0` (`ANIM-119`): it uses the chooser's speaker and the stamp `unit+0x190` of every other voice, and its threshold is 2000 ms. Executed on one unit: a move reply at 5000 ms then a selection reply at 6999 ms is silent, and at 7000 ms plays; a selection reply at 5000 ms then a move reply at 7999 ms is silent, and at 8000 ms plays; a move reply then a retreat reply is silent at 2999 ms and plays at 3000 ms.

**Confidence.** **High** for the gate matrix and the shared stamp: the original tail and readers executed on both roots. **Medium** for the bit meanings (`+0x144` bit 0x1 follows the class-name test "CUnit", bit 0x4 is ownership of the primary object; `AI-PANEL-061`, which does not establish the other bits) and for the absence of other callers (the eleven-caller census of `VIDEO-067`). **Unknown** which input branches of the routine reach the tail: it was executed from its first gate.

### VIDEO-070

- **Slots.** Guard (`0x17`) and Stand Ground (`0x18`) are panel builders; the defend order (`0x1b`) is a map and minimap builder. All three call `vt+0x70`, `FUN_0045ea70`, which reads `[bank+0x1c]`, `defend.wav` (`ANIM-094`). Retreat (`0x14`, panel) calls `vt+0x74`, `[bank+0x18]`, `retreat.wav`. A pickup order (`0x21`, map click) calls `vt+0x7c`, `[bank+0x20]`, `idle.wav`. None draws from `rand()`.
- **Banks.** The bank is the speaker's (`HERO-APPEAR-055`): `mf_hero` or `ff_hero`, `m_mage` or `f_mage`, `mf_merc` or `ff_merc`, `m_peasant` or `f_peasant`. Executed with one selected unit in a hero, hero-shaped mage, mercenary and peasant bank: guard, stand ground and defend decode to `defend` of that bank, retreat to `retreat`, pickup to `idle`.
- **Stamp.** The three share the 3000 ms stamp with every other voice, so after a reply of any gesture that unit's next command, defend, retreat or idle reply is refused for 3000 ms and its next selection reply for 2000 ms.
- **Other stance changes.** A stance or retreat set by anything other than these builders reaches none of the five slots on a unit within the census of `VIDEO-067`.

**Confidence.** **High** for the slot of each gesture and the bank rule: decoded and executed on both roots. **Medium** for the last bullet, which depends on the census.

## Cutscene presentation

| ID | Public functional claim | Confidence | Status | Evidence |
|---|---|---|---|---|
| VIDEO-071 | The player copies the sidecar start to the blit source origin once after open, arms each fade at its start frame, scales the decoder palette per frame, and adds each pan step to the source origin after the frame advance. | High / Medium / Unknown | ● active | [EXP-0447](../experiments/EXP-0447-cutscene-presentation/) |
| VIDEO-072 | A movie is drawn doubled when twice its width fits the 640 wide output region or twice its height fits the 360 high region; a 480 high movie then takes its own size as the blit region. Display size does not enter. | High / Medium / Unknown | ● active | [EXP-0447](../experiments/EXP-0447-cutscene-presentation/) |
| VIDEO-073 | Frame pacing is the decoder wait: the sound playback position when a track is open and on, else a timer in 10 microsecond units from the header interval. The player adds no floor or ceiling to the interval. | High / Medium / Unknown | ● active | [EXP-0447](../experiments/EXP-0447-cutscene-presentation/) |
| VIDEO-074 | Both roots hold 33 movies each; every header has flags 0 and one 16-bit 22050 Hz audio track with the compression flag set, mono in 15 or 16 movies and stereo in the rest. Intervals: -4000, -6666, -6673, -8333. | High | ● active | [EXP-0447](../experiments/EXP-0447-cutscene-presentation/) |

### VIDEO-071

Frame step `FUN_004ae9b0`, called by the player `FUN_0047a600` with argument 1, runs in this order:

1. Wait: the loop `FUN_004ae990` repeats `SmackWait` until it returns 0.
2. `SmackBufferFocused` on the buffer; a zero return ends the step with 1, before any state below changes.
3. Fade arm `FUN_004af260`, then pan state `FUN_004af2d0`.
4. Fade apply `FUN_004ae920` on the decoder palette (`handle+0x6c`). With a fade active the scaled palette goes to `SmackBlitSetPalette`; otherwise the raw palette goes there when `handle+0x68` is non-zero.
5. `SmackDoFrame`, display surface lock (`vt+0x64`; a failed lock at `004aea5e` ends the movie, return 0), `SmackBlit`, unlock (`vt+0x80`).
6. If the decoder frame counter `handle+0x374` equals frames minus 1 the step returns 0 and no advance runs. Otherwise `SmackNextFrame`, then pan add `FUN_004af330`, return 1.

- **Reader.** `FUN_004aee90` reads, with default 0 for every key (strings `005c0880` to `005c090c`): the `Common` keys `startx`, `starty`, `nFadings`, `nPanaramings` into player `+0x2c`, `+0x30`, `+0x14`, `+0x20`. It allocates `count * 16` zeroed bytes per family and fills record `i` from section `Fading<i+1>` or `Panaraming<i+1>`: a dword start frame at `+0`, a dword end frame at `+4`, then two floats (`startfade`, `endfade`) or two dwords (`stepx`, `stepy`). The section and key schema is `REG-CUT-053`.
- **Blit arguments.** `004aea71` to `004aeab1` pushes the `SmackBlit` parameters from the player: source x `+0x348`, source y `+0x34c`, source width `+0x350`, source height `+0x354`, destination x `+0x358`, destination y `+0x35c`. The executed blit probe of `VIDEO-072` shows the source rectangle is read at the source coordinates and the destination origin is where it lands, which separates source from destination.
- **Start.** After `FUN_004aeb10` succeeds, `004af1f5` copies `+0x2c` and `+0x30` to the source origin `+0x348` and `+0x34c`. The destination `+0x358` and `+0x35c` is set by the player before the loader, and no sidecar value reaches it. The copy follows the doubling arm of `VIDEO-072`, so a start is never halved. All 33 shipped sidecars hold `startx = starty = 0` (`REG-CUT-053`).
- **Fade.** The arm tests that the record index `+0x18` is below `+0x14` and that the record start equals the frame counter. It sets active `+0x34`, advances the index, stores remaining `+0x344 = end - start`, factor `+0x33c = startfade` and delta `+0x340 = (endfade - startfade) / (end - start)`. Each apply adds delta to the factor, multiplies all 768 palette bytes by it with truncation (`__ftol`, no clamp) into `+0x38`, and decrements remaining; at 0 the fade is inactive. A fade from frame `s` to `e` is applied on frames `s` to `e - 1`: frame `s` already shows `startfade + delta` and frame `e - 1` shows `endfade`, within float32 accumulation and with the bytes truncated (a fade-in ending at 1.0 can give 254 for 255).
- **Pan.** Per step, with the record index `+0x24` below `+0x20`: a record start equal to the counter loads `+0x360` and `+0x364` with the record's step, and a record end equal to the counter zeroes both and advances the index. After each frame advance `FUN_004af330` adds the two steps to `+0x348` and `+0x34c`, with no clamp. The first blit with a moved origin is frame `s + 1`; frame `e` is drawn at the start origin plus `(e - s)` steps.
- **Order.** Both families consume records strictly in index order. A fade whose start frame is already past, or a pan whose end frame is already past, never matches, and no later record of that family runs.
- **Executed.** None of these ROM routines was executed. The blit and wait calls they make were executed for `VIDEO-072` and `VIDEO-073`.

**Confidence.** **High** for the reader layout, the consumer of each value, the frame step order and the arithmetic: whole routines read on the one `rom.exe` both roots hold, with the key strings read from its image. The live alternatives (gamma factor, blit scale, destination offset, rectangle) are excluded because the fade writes only the palette and the pan only the source origin. **Medium** that the palette scaled at step 4 is the previous decode's: `SmackDoFrame` is the routine that writes `handle+0x6c` (`10005d64`, `10005d70`), and the state before the first decode is not read. **Unknown** the reader's result for a missing sidecar file or a negative count, the result for equal start and end frames, and how the blitter clips an out-of-range source origin.

### VIDEO-072

- **Region.** The player calls the loader with width 0x280 and height 0x168 (`0047a687`), stored in `+0x350` and `+0x354`, and the destination origin `(0x005ea210, 0x005ea214 + 0x3c)`. After the loader returns, a movie of height 0x1e0 replaces the region with its own size (`FUN_004aee70`, `0047a6a6` to `0047a6c5`) and the origin with `(0x005ea210, 0x005ea214)`. The doubling decision below runs inside the loader, before that replacement, so it uses 640x360 for every movie.
- **Decision.** `FUN_004aeb10` compares unsigned. If `2 * width > region width` and `2 * height > region height` it calls `SmackBlitOpen(mode | 1)` and sets doubled `+0x368 = 0`. Otherwise it calls `SmackBlitOpen(mode | 2)`, shifts source x, y, width and height right by 1, halves the destination origin toward zero, and sets `+0x368 = 1`; close `FUN_004aecf0` shifts them back. `mode` comes from `FUN_004ae860`: 0 for an 8-bit surface, `0x80000000` for 5-5-5, `0xc0000000` for 5-6-5, else the message `Unsupported pixel format.`.
- **Inputs.** Movie width and height, and the two region constants. The surface pixel format enters only through `mode`, and no window size or movie flag is read. The display size words `-640`, `-800` and `-1024` set only the origin, `(width - 640) / 2` and `(height - 480) / 2` (`FUN_00471790`).
- **Blit executed.** `tools/smackblit` calls the installed decoder's `SmackBlitOpen` and `SmackBlit` on synthetic indexed memory, source rectangle 5x3 at (3,1), destination (7,2). Mode 1 writes 15 cells, a plain copy. Mode 2 writes 60 cells, each source pixel as a 2x2 block, and doubles the destination origin to (14,4), so the halving above is undone by the blitter.
- **Shipped population.** Census of `VIDEO-074`. Doubled by the rule: the `320x180` movies, 4 on EN (`INTRO/04` and `M150/01`, each in both containers) and 2 on RU (`INTRO/04`), drawn as 640x360. Not doubled: 22 `640x360` and 4 `800x360` on EN, 24 and 4 on RU, and 3 `640x480` logos on each root. An `800x360` movie shows a 640 wide window of its source.

**Confidence.** **High** for the decision, the constants, the halving arithmetic and the indexed blit kinds: whole routines read and the blitter executed. **Medium** that no other call of the loader sets a different region: `FUN_004aee90` has the single call site `0047a687`, but a computed call is outside the read. **Unknown** the blit output for the two 16-bit surface modes, whether ROM1 ever runs in one, and the effect of the halving on an odd destination origin, which was not executed.

### VIDEO-073

- **Player.** Step argument 1 selects the blocking wait; the only pacing in the player is the `FUN_004ae990` loop. All pending messages are drained before each step (`0047a6e8`) and none during the wait; `WM_CLOSE`, `WM_QUIT`, key down, system key down and left or right button down (`0x10`, `0x12`, `0x100`, `0x104`, `0x201`, `0x204`) end the movie at `0047a746`. The constructor calls `SmackSoundUseDirectSound` (`004ae7cb`), which installs the readback `100096b0`. The player opens with track flags `0xff000` (bit `0x80`, the frame-rate override, is clear), so the header interval is used, and it turns sound on after open (`SmackSoundOnOff(h, 1)`, `004aecb5`).
- **Interval.** Open reads the signed interval `v` at header `+0x10` (`100045a6`, `100045ae`; the open and `SmackWait` bodies `10007863` to `1000796f` are the `VIDEO-050` listing). For `v >= 0` the state value `+0x424` is `v * 100` modulo 2^32; for `v < 0` it is `-v`. The unit is 10 microseconds, so `-6666` is 66.66 ms and a positive value is milliseconds. The clock is `timeGetTime * 100`.
- **Timer wait.** Used when no track opened (`+0x444 = -1`). The deadline `+0x420` starts at -1: the first poll sets it to now plus the interval and returns not ready. After that a poll with now below the deadline is not ready. Otherwise it is ready and the deadline becomes: now plus interval if now is within one interval after the deadline; deadline plus interval if now is up to one further interval late; now plus interval if later. A non-zero interval sets the latch `+0x41c`, which keeps later polls ready until `SmackNextFrame` clears it.
- **Sound wait.** Used when a track opened and sound is on. Each `SmackDoFrame` adds the bytes per frame to the track targets and arms `+0x448` (`10005dd0`). The wait is not ready while target `(+0x20 >> 10) + +0x24` exceeds the playback position plus 8; becoming ready clears `+0x448`, and while it is clear polls are ready. The position readback `[0x10015f7c]` is `100096b0` for the DirectSound backend, interpolated with `timeGetTime`. The first track that opens is the clock.
- **Audio ahead or behind.** Behind: the wait holds the frame until the position reaches the target, with no timeout. Ahead: the wait returns ready without sleeping, and `SmackDoFrame` skips the decode (`10006180` returns 1 once the position passes the second target `+0x28`) unless open flag `0x400` is set. ROM1 does not set it and ignores the return, so it blits the unchanged buffer and advances. Sound off (`+0x44c` non-zero) makes every poll ready, with no pacing at all.
- **No audio.** A movie whose track word lacks bit 30 or has rate 0 opens no track and uses the timer. Executed: `tools/smackpace` polled `SmackWait` on header-patched scratch copies of one shipped movie, opened with no track flag. The shipped interval -6666 became 6666 units and the first ready poll of frame 1 came at 66 ms.
- **Unfocused.** `SmackBufferFocused` is 0 unless the buffer's window is the foreground window, or, for a direct-draw buffer, its focus word `+0x460` is non-zero. The step then returns before decode or advance, so the movie holds its frame and the loop spins on ready polls.
- **Degenerate intervals.** Frame 1, first ready poll, cap 450 ms. Header 0 gives 0 units, ready at once, and the latch is never set. Headers -1 and 1 give 1 and 100 units, ready within 5 ms. Header 100 gives 100 ms, ready at 99 to 100 ms. Header -100000 gives 1 s, not ready within the cap. Header 2147483647 gives 4294967196 units, a deadline 100 units in the past, ready at once. Header 42949673 gives 4 units, ready at once; a positive value of 42949673 or more wraps. Header -2147483648 negates to 2147483648 units, and whether it waits depends on the boot-time phase of `timeGetTime * 100`, so it is not in the reproduced evidence. With sound open, interval 0 arms nothing: no wait and no skip. The player reads no interval and applies no floor or ceiling.

**Confidence.** **High** for the interval derivation, the timer wait arithmetic and the absence of a floor in the player: read in full, and executed on the installed decoder through `timeGetTime`, timer path only. **Medium** for the sound wait, the skip rule and the unfocused behaviour: read in full, but no sound output was opened, so no sound position was observed. **Unknown** the audible timing and the cursor's real advance rate, the result when no audio device or backend is available, the effect of losing focus while sound plays, and what the final-frame return does to queued sound.

### VIDEO-074

- **Population.** Every `.res` and `.lm` container on each root plus loose `.smk` files (none found): 12 EN and 11 RU archives opened, 33 `.smk` nodes per root, 18 in `VIDEO4.RES` and 15 in `VIDEO8.RES`. No node of another name carries the `SMK2` magic. `rom.exe` also names `video4\rom.smk` (`005be518`, `00474a54`); no container node or loose file has that name, and what the call does without it is Unknown. The header fields are those the decoder's open routine reads: `+0x14` flags, `+0x48 + 4 i` audio word of track `i`.
- **Audio word.** Bit 31 compressed, bit 30 present, bit 29 16-bit, bit 28 stereo, low 24 bits the rate. Track 0 of every movie is `e0005622` (flag set, 16-bit, mono, 22050 Hz) or `f0005622` (the same, stereo); bits 27 to 24 are clear in 66 of 66 and the decoder reads only bits 31 to 28 and the rate. Tracks 1 to 6 hold `00000000` in all 66 headers. The decoder routes bit 31 to the compressed-audio routine `10011f30` (`10006028`, `10006051`) and copies raw otherwise. The stored coding is not decoded here; the routine reads a bit-tree stream, consistent with Smacker audio, which is inference.
- **By container.** `VIDEO8.RES` holds 15 stereo movies on each root. `VIDEO4.RES` holds 15 mono and 3 stereo on EN, 16 mono and 2 stereo on RU. Its stereo movies are the logos `LOGOS/1c`, `buka` and `nival` on EN, and `buka` and `nival` on RU; RU `LOGOS/1c` is mono.
- **Header.** Flags are 0 in 66 of 66 (no ring, no y-scale). Intervals: `-6666` in 30 EN and 29 RU, `-6673` (`M120/01`, both containers) in 2 on each root, `-4000` (`LOGOS/buka`) in 1 on each root, `-8333` (RU `LOGOS/1c`) in 1. Frames run 75 to 1309 and durations 5.0 to 87.3 s. No interval is zero or positive.

**Confidence.** **High** for the census: every node of the named archives, parsed by `tools/smkcensus`. It does not cover movie files outside `Allods/VIDEO4.RES` and `Allods/VIDEO8.RES`, or the audio payload.
