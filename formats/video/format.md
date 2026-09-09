<a id="audio-and-video--public-functional-specification"></a><a id="cutscenes"></a>

# Audio and video resource contract

ROM1 selects music, sound effects and cutscenes through the resource resolver.
Playback state belongs to the game. Music and sound payloads use standard
audio formats; numbered movie payloads use Smacker. The bundled decoder's
internal ABI and implementation are outside this reference.

<a id="scope"></a>

## Resource families

| Surface | Resource | Game state |
|---|---|---|
| Music | `MUSIC.RES`, archive identity `music` | Candidate list, shuffled order, fixed-source mode, enabled flag and volume |
| Sound effects | Sparse numeric entries in `SFX.RES`, plus direct paths | Sample selector, volume, pan, loop/play state, priority and optional frequency |
| Cutscenes | `VIDEO4.RES` or `VIDEO8.RES`, numbered `.smk` members and same-basename `.reg` | Selected namespace, mission identity, scan index, fade/pan and stop state |

— VIDEO-MUSIC-001, VIDEO-MUSIC-002, VIDEO-MUSIC-003, VIDEO-MUSIC-004,
VIDEO-MUSIC-005, VIDEO-MUSIC-006, VIDEO-MUSIC-007, VIDEO-MUSIC-008,
VIDEO-MUSIC-009, VIDEO-MUSIC-010, VIDEO-MUSIC-011, VIDEO-MUSIC-012,
VIDEO-SFX-013, VIDEO-SFX-014, VIDEO-029, VIDEO-030

<a id="shipped-payload-shape"></a><a id="namespace-and-selection"></a>

## Music

Installed music members are stereo PCM RIFF/WAVE, 22050 Hz and 16 bits per
sample. Their payloads have no identified loop metadata. Looping and track
succession are playback state. — VIDEO-MUSIC-001, VIDEO-MUSIC-002

The program has fixed candidate lists for UI/game contexts and a separate
mission/background family. Ordinary playback advances through a randomized
candidate order: a one-member list repeats, and a multi-member list advances
through its shuffled order. Fixed-source mode can retain one candidate.
The context-list assignment to named surface transitions retains Medium
confidence. — VIDEO-MUSIC-003, VIDEO-MUSIC-004, VIDEO-MUSIC-005,
VIDEO-MUSIC-006, TOWN-372

Archive mounting follows the resource resolver and working-directory rules.
A missing music archive is a recoverable resource condition. Enable/disable
and volume are independent settings. The player uses a streaming buffer.
— VIDEO-MUSIC-007, VIDEO-MUSIC-008, VIDEO-MUSIC-009, VIDEO-MUSIC-010,
VIDEO-MUSIC-011, VIDEO-MUSIC-012

<a id="public-receiver-contract"></a><a id="sample-selection"></a>

## Sound effects

The common play boundary receives attenuation/volume, pan, play/loop state,
priority/category and optional playback frequency. Sample selection precedes
that call. — VIDEO-SFX-013, VIDEO-SFX-014, VIDEO-SFX-015

Selectors include registry-backed numbers, unit class/action values,
spell/projectile values, ambient state and literal UI/voice paths. The
registry is sparse; a valid selector can have no installed sample. Ambient
effects use visible terrain/object state and scheduled replay.
— VIDEO-SFX-016, VIDEO-SFX-017, VIDEO-SFX-018, VIDEO-SFX-019,
VIDEO-SFX-020, VIDEO-SFX-021

<a id="resource-selection"></a><a id="sidecar"></a><a id="user-stop-behaviour"></a><a id="decoder-boundary"></a>

## Cutscene sequence

1. Select the low/high video namespace from startup/environment state.
   No general fallback to the other namespace follows selection.
2. Form a numbered movie path from mission identity and an increasing
   two-digit index. Missing or unopenable members can advance the scan.
3. Replace the movie extension with `.reg` to obtain its sidecar. Apply its
   initial position, frame-bounded fades and frame-bounded pans.
4. Supply the resource-backed movie to the decoder. Draw decoded frames into
   the game-owned destination and track palette/frame progress.
5. Stop the complete numbered scan on ordinary key-down, system-key-down,
   mouse-button-down, close or quit events. Natural completion and some
   open/construction failures instead continue the scan.

— VIDEO-029, VIDEO-030, VIDEO-031, VIDEO-032, VIDEO-033, VIDEO-034,
VIDEO-035, VIDEO-036

The decoder boundary needs resource input, frame dimensions/count, palette
changes, frame progression/timing and destination ownership. Sidecar fade/pan
state and user interruption operate around this boundary. The sidecar is a
[REG store](../reg/format.md); its records are not part of the SMK bitstream.
— VIDEO-045, VIDEO-046, VIDEO-047, VIDEO-048, VIDEO-049, VIDEO-050, VIDEO-051

<a id="a-namespace-that-is-not-a-cutscene-surface"></a>

## School training pictures

The four school-training resource families belong to the school-room loader
and painter. They use fixed room-relative anchors; an active sequence
replaces the current still on the same side. These named paths do not pass
the frames to a movie/cutscene surface. — TOWN-427, TOWN-428

The classification covers the named literal/direct/stored-pointer and
field-to-draw paths. Runtime-computed filenames or targets elsewhere are not
excluded. Native visibility and audio remain unverified; the animated arms
make no direct sound request. — TOWN-434

<a id="unknown--bounded-areas"></a>

## Unknowns

All OS/decoder error paths, arbitrary malformed media, hardware/driver timing,
runtime-computed sample/movie names and sidecar branches beyond the named
consumers remain unspecified. No container writer or new audio/movie codec
is defined by this behavioral contract. Standard payload production and
bundled decoder internals are separate subjects.
