# Againrom knowledge

Research-derived documentation of **Rage of Mages 1** and selected **Rage of
Mages II** file formats, interfaces and engine behaviour, with a small claim-reading
tool. The project studies compatibility with independently obtained installations;
this repository is not a distribution of either game or a decoder SDK.

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

ROM2 coverage is incomplete: the k1 export omitted `claims/rom2-asset.md` while retaining
pages that depend on it. This edition records that gap rather than inventing replacement
claims or copying an unreviewed private ledger. Such pages are not self-contained
implementation authorities until their cited findings have been published.

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
