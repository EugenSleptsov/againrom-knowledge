# Againrom knowledge

Public, derived research about **Rage of Mages 1** file formats and engine behaviour:
clean-room claims and specifications, exported one-way from a private research
repository. This repository publishes conclusions and their evidence citations. It
carries no research narrative, no raw evidence, and no game data.

## Clean room and interoperability

Every claim here was derived by observing a lawfully owned install and reading its
compiled code, never by consulting another implementation, a third-party
specification, or a reverse-engineering write-up of Rage of Mages. The
interoperability purpose is exact: understanding the file formats and behaviour a
legally obtained copy of the game already exhibits, so that independent, clean-room
software can read and reproduce them. No game asset, executable, or install byte is
committed anywhere in this repository.

## What is here

- `claims/` — one ledger per format or subsystem. Each row is a stable, permanent ID, a
  claim, a confidence level, a status, and an evidence citation.
- `formats/` — one folder per format or subsystem, holding the specification distilled
  from claims that survived promotion. Only promoted, evidence-backed facts land here.
- `tools/claim/` — the reader. `go run ./tools/claim <ID>` prints one claim row, and its
  retraction state if any, instead of a whole ledger.
- `SOURCE.md` — which private-repository commit this snapshot was exported from, and
  the counts that describe it.

## Reading a claim

Read one claim at a time, for example `go run ./tools/claim SAV-SACKENTRY-590`. This
prints the row from its ledger, marked if a later row in `claims/retracted.md`
overturned it. Ledgers run to hundreds of rows; the reader exists so that a one-claim
question never costs opening a whole file.

## Confidence

Confidence measures whether the evidence discriminates between live alternatives, not
how persuasive a claim reads.

| Level | Bar |
|---|---|
| High | Live alternatives are ruled out: an instruction-level reading of the compiled consumer, corpus agreement, and a failed falsification attempt. |
| Medium | Corpus-consistent and internally coherent, but another fitting model remains live. |
| Low | Partial or single-instance support. |
| Unknown | The evidence does not decide. |

## Retraction

A claim ID is permanent once allocated; it is never reused or deleted. When later
evidence overturns a claim, `claims/retracted.md` preserves the former wording, the
confidence it was believed at, the overturning evidence, and the corrected truth, and
the per-format ledger row is amended to the corrected statement. Nothing in that
history is edited to look better in hindsight. `go run ./tools/claim <ID>` reads a
claim's current row together with its retraction state in one call.

## Evidence

Claim text cites the experiment that earned it, commonly as a path such as
`../experiments/EXP-0293-sack-cell-lifecycle/`. That path does not resolve inside this
repository: raw evidence, disassembly, and experiment write-ups live in the private
research repository this snapshot was exported from. `SOURCE.md` names its exact
commit.

This repository is a one-way export of that repository's `claims/`, `formats/`, and
`tools/claim/`. It is never edited directly here. A correction to a claim starts as a
falsifiable question against the private research repository, not an edit to a file in
this one.

## License

Documentation (`claims/`, `formats/`, `SOURCE.md`, this file) is licensed under
[CC BY 4.0](LICENSE). Code (`tools/claim/`) is licensed under
[Apache License 2.0](LICENSE-CODE). See `NOTICE` for the trademark and asset statement.
