# Audio and video — public functional specification

Level 3 functional edition. This page describes the behaviour that an independent
consumer needs to reproduce at the **game boundary**. Detailed instruction listings,
private experiment output and internal implementation details of bundled third-party
code are intentionally not reproduced here.

Primary ledgers: [`claims/video.md`](../../claims/video.md) and the claim IDs named
below. Source snapshot and publication boundary: [`SOURCE.md`](../../SOURCE.md) and
[`PUBLICATION.md`](../../PUBLICATION.md).

## Scope

ROM1 uses three relevant media surfaces:

1. music streams stored in `MUSIC.RES`;
2. sound effects stored primarily in `SFX.RES` plus some directly named resources;
3. numbered cutscenes stored in `VIDEO4.RES` / `VIDEO8.RES`, with `.reg` sidecars and
   Smacker (`.smk`) movie payloads.

The public specification documents the game's resource-selection, state and timing
contracts. It does **not** publish a reverse-engineered implementation specification of
the bundled Smacker decoder. A reimplementation should use independently lawful
support for the standard media format rather than copy bundled decoder expression.
See [`THIRD_PARTY.md`](../../THIRD_PARTY.md).

---

# Music

Promoted claim family: `VIDEO-MUSIC-001`, `VIDEO-MUSIC-002`, `VIDEO-MUSIC-003`,
`VIDEO-MUSIC-004`, `VIDEO-MUSIC-005`, `VIDEO-MUSIC-006`, `VIDEO-MUSIC-007`,
`VIDEO-MUSIC-008`, `VIDEO-MUSIC-009`, `VIDEO-MUSIC-010`, `VIDEO-MUSIC-011`,
`VIDEO-MUSIC-012`.

## Shipped payload shape

The shipped music members examined by the research are ordinary PCM RIFF/WAVE streams:

- stereo;
- 22,050 Hz;
- 16 bits per sample;
- no loop metadata was observed in the payloads.

Looping and succession are therefore player state, not a property a compatible reader
must infer from WAVE metadata.

The exact shipped filenames, file lengths and track durations are corpus evidence, not
part of the wire grammar, and are omitted from this public edition. Consumers discover
members through the resource namespace rather than a copied content table.

## Namespace and selection

Music resources resolve under the `music` archive identity. The original program uses
fixed candidate lists for several UI/game contexts and a separate mission/background
candidate family. Reading those context lists as named surface transitions inherits
`TOWN-372` at Medium. Ordinary playback advances through a randomized candidate order;
a one-member list therefore repeats, while a multi-member list moves through that
shuffled order. A fixed-source mode can retain one selected candidate.

Availability of the archive is environment-dependent: mounting follows the game's
ordinary resource resolver and working-directory rules rather than an abstract
"always mounted" guarantee. A consumer should treat missing music as a recoverable
resource condition.

Music enable/disable and volume are independent state. The original uses a streaming
buffer and keeps looping/succession state outside the WAVE payload.

---

# Sound effects

Promoted claim family: `VIDEO-SFX-013`, `VIDEO-SFX-014`, `VIDEO-SFX-015`,
`VIDEO-SFX-016`, `VIDEO-SFX-017`, `VIDEO-SFX-018`, `VIDEO-SFX-019`, `VIDEO-SFX-020`,
`VIDEO-SFX-021`.

## Public receiver contract

The game reaches a common sound-play boundary after choosing a sample. The observable
parameters are:

- attenuation/volume;
- pan;
- play/loop state;
- priority/category;
- optional playback frequency.

A compatible engine may represent these differently internally; interoperability
requires reproducing the resulting selection and audible behaviour rather than the
original object's memory layout.

## Sample selection

The game selects effects through several routes:

- registry-backed numeric selectors;
- unit-class / action-driven selectors;
- spell/projectile selectors;
- ambient selectors;
- directly named resource paths for some interface and voice effects.

The shipped registry contains sparse entries, and not every possible selector resolves
to a shipped sample. Missing samples therefore must not be treated as proof of a parser
failure.

Ambient effects depend on visible terrain/object state and use scheduled replay rather
than being embedded into terrain files as continuously playing audio.

The complete shipped filename/slot table is research evidence and is not reproduced in
this public functional edition.

---

# Cutscenes

Promoted claims: `VIDEO-029`, `VIDEO-030`, `VIDEO-031`, `VIDEO-032`, `VIDEO-033`,
`VIDEO-034`, `VIDEO-035`, `VIDEO-036`, and the bounded game-side portions of
`VIDEO-045`, `VIDEO-046`, `VIDEO-047`, `VIDEO-048`, `VIDEO-049`, `VIDEO-050`,
`VIDEO-051`.

## Resource selection

The game chooses one of two video namespaces corresponding to the low/high media sets.
The selected namespace is controlled by startup/environment state; there is no general
fallback from one namespace to the other after selection.

A mission cutscene request searches numbered movie members using the mission identity
and an increasing two-digit index. Missing/unopenable members can advance the scan;
a normal user-stop condition terminates the scan.

## A namespace that is not a cutscene surface

The school training picture namespace is not consumed by the cutscene/movie surface in
the measured population. Its four resource families are loaded by two family loaders
whose sole direct caller is the school-room loader, and they bind the loaded series to
fields the school-room painter later reads. Both preserved roots ship the same gapless
populations for those families. — `TOWN-427`

That painter places the two sides at fixed room-relative anchors. An active animated
sequence replaces the same side's current still picture; no frame is handed to a
separate cutscene object, movie player or transition surface. — `TOWN-428`

The classification is bounded to literal references, direct calls, parsed stored code
pointers and the reached field-to-draw chains. Inside that population it excludes the
separate-movie, data-only and mixed-family ownership models, but not a computed filename
or a computed target elsewhere. Original visibility and audio are not witnessed, and the
animated arms issue no direct sound request of their own. — `TOWN-434`

## Sidecar

For a requested movie, the game derives a same-basename `.reg` sidecar by replacing the
last filename extension. The sidecar supplies presentation metadata including:

- initial screen position;
- a sequence of frame-bounded fade records;
- a sequence of frame-bounded pan records.

These records affect presentation around movie decoding; they are not part of the SMK
compressed bitstream itself.

## User stop behaviour

The original cutscene wrapper treats ordinary key-down/system-key-down, mouse-button
down, close and quit events as stop conditions for the current numbered cutscene scan.
Natural movie completion and some open/construction failures follow the continuation
path instead. — `VIDEO-036`

## Decoder boundary

The original game supplies a resource-backed movie source to a bundled Smacker decoder,
renders decoded frames into a game-owned destination, observes palette/frame progress,
and advances until the game-side stopping condition is met.

For interoperability, a replacement needs the following functional properties:

- decode the selected `.smk` resource from its resource origin;
- expose frame dimensions/count and palette changes as required by the movie;
- render frames into the game's presentation surface;
- support frame progression/timing sufficiently for sidecar fade/pan state and user
  interruption to remain synchronized;
- preserve ownership boundaries between game-managed resource/destination objects and
  decoder-managed state.

The private research identified additional behaviour of the bundled decoder under
`VIDEO-045`, `VIDEO-046`, `VIDEO-047`, `VIDEO-048`, `VIDEO-049`, `VIDEO-050` and
`VIDEO-051`. **Those decoder-internal ABI, memory-layout and instruction
level details are intentionally not reproduced here.** They are not needed to define
ROM1's public game-side media contract and require a separate third-party rights review
before any broader publication.

A replacement implementation should rely on an independently available lawful Smacker
implementation/specification or another decoder capable of the required `.smk` input,
not on copied expression from the bundled library.

## Unknown / bounded areas

The following remain explicitly outside the proved public contract:

- every OS- and decoder-error presentation path;
- exact audible timing under all hardware/driver combinations;
- all malformed-media behaviour;
- every possible runtime-computed sound/movie name;
- media SDK behaviour not necessary to explain ROM1's observed calls;
- sidecar behaviours not reached by the preserved corpus.

These boundaries are intentional. Absence of a rule here must not be filled by copying
third-party decoder code or by promoting a private experiment detail into a public
format rule without review.
