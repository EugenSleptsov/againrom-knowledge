# Publication policy

This is a maintainer review policy, not a legal opinion or an additional condition on
reuse of work licensed under CC BY 4.0 or Apache-2.0.

## Preserve useful knowledge

Prefer independently written descriptions of file layouts, API/wire contracts,
algorithms, state transitions, exact arithmetic and observed behaviour. Keep the claim
ID, evidence identity, confidence, scope, edge cases and explicit Unknowns. A finite
corpus observation must not become a universal rule during shortening.

Addresses and offsets are not categorically prohibited: a wire offset may be essential
to interoperability, and a research locator may be necessary to trace evidence. Review
literal instruction sequences separately from independently expressed functional rules.
Do not use a global regular expression to delete instructions: some legacy rows express
an essential condition only in the cited operation.

## Material needing a separate decision

Do not export original game/decoder binaries, assets, complete resource payloads or
reconstructable content tables merely because they are text rather than binary files.
Review long code excerpts, disassembly/decompilation, quoted game text and encoded data
before export. Keep unnecessary original expression in private evidence. Any retained
excerpt needs a recorded purpose, proportionate scope, source and applicable permission
or other legal basis; a technical confidence grade is not that basis.

A paraphrase alone does not decide whether obtaining or publicly disclosing the
underlying information is permitted. Likewise "old game", "noncommercial", "clean
room" and "interoperability" are not blanket clearances. Apply the appropriate rules
to the actual method, purpose, component and jurisdiction rather than making an
all-repository legal guarantee.

## Evidence and authoring

Maintain two linked records: the unabridged private research snapshot and the public
functional edition. Record material omissions and editorial transformations without
fabricating experiments or changing what an old confidence level meant. Preserve
retractions and their clause scope; do not erase a former error to improve appearances.

Factual corrections belong in the research process. Publication corrections can be
reviewed here, but must feed the upstream export transformation before the next release.
Do not maintain two independently edited factual ledgers.

## Before approving an export

1. Record input identity, acquisition/terms review, method and purpose for the affected
   game or component. Separately assess the need and basis for public disclosure. Do
   not infer these from a SHA or a general README statement.
2. Inspect content as well as filenames, including reconstructable tables, quotations,
   code excerpts and embedded encodings. Run the supplied check as an aid, not a verdict.
3. Confirm that all public normative claims have available authorities or are clearly
   marked incomplete; internal experiment citations are allowed but not publicly
   reproducible evidence. Check license notices and public claims against real contents.
4. Verify claim IDs, confidence levels and retraction scope against the source; inspect
   the semantic diff. Review branches, tags and history separately when withdrawal of
   earlier material is required.

## Tool limits and current backlog

`python3 scripts/publication_check.py --strict` examines tracked working-tree content,
not previous commits, Git LFS storage, releases, attachments or caches. It reports only
locations/codes, not the potentially sensitive source excerpts. Inline-link parsing and
instruction detection are heuristic; reference-style links, unusual encodings and
paraphrased copied expression can be missed. Small legitimate technical fragments may
be flagged. No automatic deletion or approved-all baseline is provided.

### Completed in `publication/legal-hygiene-k1`

- `claims/registry.md` was replaced by a public index; private experiment/review narrative
  was removed from that surface.
- `claims/fame.md` and `formats/fame/format.md` no longer reproduce the shipped default
  hall-of-fame table or unnecessary instruction excerpts.
- `formats/res/format.md` now states the RES wire grammar and resolver semantics without
  executable-address/disassembly evidence or the complete shipped archive inventory.
- `formats/text/format.md` keeps byte conversion, indexing and table semantics while
  omitting localized prose and reconstructable shipped string/content tables.
- `claims/video.md` and `formats/video/format.md` now expose the ROM1 game-side media
  contract without bundled-Smacker structure offsets, instruction evidence or complete
  media-name/content inventories. Decoder-internal evidence remains private pending a
  separate third-party review.
- README/NOTICE/license scope, source mapping, third-party boundary, contribution rules,
  review tooling, synthetic tests and read-only CI were added or corrected in the first
  pass.

### Remaining review

This edition is still a **partial pass**. Remaining work includes:

- instruction-heavy ledgers other than `claims/video.md` and the already-minimized FAME
  ledger;
- `claims/retracted.md`, while preserving exact retraction/narrowing scope;
- copied/reconstructable content tables in other claim/format areas;
- the missing public `R2-ASSET` ledger and dependent ROM2 format pages;
- component-specific publication bases and acquisition/terms records;
- a final semantic check that public claim IDs, confidence and Unknown boundaries still
  match the pinned private source after editorial shortening.

Neither this checklist nor a green synthetic-test job is legal clearance. The strict
review job is expected to remain red while unresolved review candidates exist.

A content cleanup commit does not remove its parent tree. Changes to history, `k1`,
visibility or published copies require a distinct, explicit decision; preserve private
provenance rather than destroying evidence.
