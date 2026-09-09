# Againrom knowledge

Research-derived documentation of **Rage of Mages 1** and selected **Rage of
Mages II** file formats, interfaces and engine behaviour, with a small claim-reading
tool. The project studies compatibility with independently obtained installations;
this repository is not a distribution of either game or a decoder SDK.

## SAV readiness

This is delivery status for Againrom engine `5bccb0d`, not evidence about ROM1.
AGS is the engine's native save format. "From scratch" means without an imported
SAV or a copied SAV document; lawful installed game resources are still required.

| Capability | Status | Exact scope |
|---|---|---|
| SAV writer alongside AGS | **Ready for the selected save workflow** | Town: AGS, SAV or both from one captured state. Mission: exact AGS, or SAV as an ordinary return to town with recovered HP/mana, retained XP/items/money and an unfinished, restartable mission. The live mission is unchanged. |
| Full original SAV reader | **Accepted for city/world structural reading and the agreed restoration scope** | The complete known document is read, including the object graph and opaque fields. The 13-structure restoration milestone is accepted; the recorded EN/RU gate resumes all 62 discovered valid world saves with zero measured mismatches or refusals. Remaining field meanings, unobserved serializer classes and non-corpus gameplay behaviour are not declared complete. |
| SAV creator from scratch | **Supported cities: ready. Full mission worlds: not ready** | The city producer constructs the document from current state without an imported SAV. The owner accepted the generated EN town case through tavern, SAVE, restart and LOAD. Complete source-free mission construction remains paused. |

City SAV requires a supported settled campaign state, starting at main chapter30
in the shipped campaign. Pre-town and completed-campaign states, incompatible
graphs and other named writer limits refuse explicit SAV; AGS remains available.
Original RU acceptance and arbitrary generated-state interoperability remain open.

The converter separately supports bounded current-world SAV output from a previously
imported mission. That path retains source authority for unresolved state; it is
neither the mission save dialog's city-return policy nor a creator from scratch.
See the [save workflow](https://github.com/EugenSleptsov/againrom-engine/blob/5bccb0de8535fed48940b9178ee74f879651e0c6/docs/1173/story.md),
[converter](https://github.com/EugenSleptsov/againrom-engine/blob/5bccb0de8535fed48940b9178ee74f879651e0c6/pkg/game/saveconvert.go)
and [SAV format reference](formats/sav/format.md).

## Sources and publication boundary

The research source policy permits the owner's game installations, their runtime
observations, project probes, and documented general knowledge of standard formats
and platform interfaces. It excludes other game reimplementations and their derived
specifications as factual authorities. The exact research snapshot is recorded in
[SOURCE.md](SOURCE.md).

That is a source policy and provenance trail, **not proof that no contributor has
previously encountered outside material**, and not legal clearance for every finding.
Hashes establish the exported bytes; reproducing a result does not by itself establish
its historical independence. A separate repository does not establish personnel
separation or retroactively change how an earlier result was obtained.

The k1-derived corpus still contains legacy instruction excerpts, quoted technical
strings and internal research references. This editorial pass removes selected
unnecessary reproduction; it does **not** certify the remaining corpus or Git history
as free of original expression. See [the change record](PUBLICATION-CHANGES.md),
[publication policy](PUBLICATION.md) and [component/source distinctions](THIRD_PARTY.md).

## Contents

- `claims/`: permanent IDs, findings, confidence, status and evidence references.
- `formats/`: implementation-facing descriptions derived from claims. Scope and
  qualifications remain part of each description; a fact about the observed files is
  not necessarily a format-wide rule.
- `tools/claim/`: a reader that prints individual claims and their retraction entries.
- `SOURCE.md`: upstream snapshot identity and the scope of this public edition.
- `scripts/`: project-authored, non-mutating publication-review tooling.

The ROM2 functional edition includes all 28 `R2-ASSET` claim rows and the seven
bounded survey pages. Their claim dependencies are present; their full evidence
remains in the private source named by [SOURCE.md](SOURCE.md). Layout/header
coverage does not establish every field meaning, runtime consumer or protocol.

## Reading and checking

Go 1.21 or newer is needed for the reader; it has no third-party module dependencies.
Run from the repository root:

```sh
go run ./tools/claim SAV-SACKENTRY-590
go run ./tools/claim -k 'sight range'
go run ./tools/claim -stats
```

The supplied Python 3.9+ check is heuristic. It reads tracked working-tree files and
reports binary/data candidates, instruction excerpts and unresolved inline links or
claim references, without changing files or printing their contents:

```sh
python3 -m unittest discover -s scripts -p 'test_*.py'
python3 scripts/publication_check.py --json
python3 scripts/publication_check.py --strict
```

`--strict` also fails on review candidates. Neither a zero exit status nor absence of
matches means legal clearance, absence of embedded game content, or independent origin.
The existing corpus is expected to produce findings; there is no accepted-all baseline.

## Confidence and corrections

Confidence concerns the strength and scope of research evidence, not permission to
publish or reuse material. A row may have High confidence for one clause and Medium or
Unknown for another. Read its current wording and qualification, not just the first
rating word.

A claim ID is permanent. `claims/retracted.md` records withdrawn or corrected claims,
including clause-limited corrections. Publication editing must preserve IDs, confidence
levels, uncertainty and the meaning and scope of corrections. Original evidence and
unabridged research records remain in the private research snapshot.

An `experiments/...` citation names a private research record at the revision identified
in `SOURCE.md`; it is not a working public evidence link. Public verification of those
records is therefore limited.

Corrections to facts about the original still originate in the research process.
Publication wording and tooling may be proposed in a public branch. Accepted editorial
transformations must be incorporated into the export process before the next snapshot,
so a raw copy does not silently restore removed material. See [CONTRIBUTING.md](CONTRIBUTING.md).

## Licensing

The project's own documentation is offered under [CC BY 4.0](LICENSE). Project-authored
code in `tools/` and `scripts/` is offered under [Apache-2.0](LICENSE-CODE).

These grants cover only rights the contributors are entitled to license. They do not
relicense protected game or third-party expression quoted in a finding, grant rights to
original assets or binaries, or grant trademark rights. This clarification does not
withdraw or add restrictions to the licenses already granted for project-owned work.
See [NOTICE](NOTICE) and [THIRD_PARTY.md](THIRD_PARTY.md).
