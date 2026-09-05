# Source and public edition

## Snapshot

Snapshot k1, exported from the private research repository at commit
`80670a957c38a6b281e8144e6ad76dc81cd1abb9`.

| Field | Count |
|---|---:|
| Claim ids | 1978 |
| Retracted ids | 319 |
| Format pages | 35 |

Each count is recomputed from this snapshot's own exported ledgers, not carried forward
from an earlier snapshot. A commit identity is not a licence, a proof of purchase, or a
proof that no outside material was encountered.

One permanent id is withheld. `RES-HDR-001` has no statement: the research index carried
a placeholder row for it and no ledger row defines it. A placeholder is not a claim, and
a statement is not invented to fill a permanent id, so the id is not published and the
count above is 1978 rather than the 1979 index rows. Whether to promote or retire it is
a research question, not a publication edit.

## Public edition

The public text is an edition of the research text, not the research text verbatim: the
same claims expressed as functional conclusions, without disassembly listings, quoted
game text beyond a short identifier, or reconstructable shipped-content tables. The
edition is tagged `k1` with this snapshot.

[PUBLICATION.md](PUBLICATION.md) states the policy and the remaining review backlog,
[PUBLICATION-CHANGES.md](PUBLICATION-CHANGES.md) records what the edition changed and
what it left open, and [THIRD_PARTY.md](THIRD_PARTY.md) separates an analysed component
from a borrowed factual source. The export refuses a later snapshot that would replace
an edited page with text the edition removed.

The edition changes wording only. Claim ids, evidence identity, confidence grades,
status values, retraction and amendment scope, explicit Unknowns, residual-scope
statements and exact arithmetic are preserved; a factual correction belongs to the
research process, not to this repository. This is a partial pass: `PUBLICATION.md` lists
the ledgers and content tables still awaiting review, and the strict review job is
expected to stay red while that backlog exists.

An `experiments/...` path cited inside claim text names the private research repository
at the commit above. It is not a path inside this repository, and the unabridged evidence
behind an omitted excerpt stays there.

`claims/rom2-asset.md` is not exported, so public pages citing `R2-ASSET-*` have
incomplete public provenance. The gap is recorded, not filled from private material.

## History boundary

This is the first published snapshot of this repository. Every later snapshot is a new
commit and a new tag; publishing one does not rewrite, move or delete an earlier tag, and
does not withdraw earlier published material. The private research record is preserved
whatever the public repository carries.
