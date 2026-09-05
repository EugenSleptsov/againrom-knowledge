# Publication hygiene: first editorial pass

Base: the unedited research-derived export at private research source
`80670a957c38a6b281e8144e6ad76dc81cd1abb9`. That raw tree is not part of this
repository's published history; the first published tree is the edition described here.

## Changes in this branch

- README/NOTICE now distinguish the actual corpus and source policy from an absolute
  no-excerpts/no-exposure or legal-clearance assurance. Scope includes selected ROM2
  research and static, runtime and corpus evidence.
- SOURCE identifies the editorial edition separately from k1, records the omitted
  R2-ASSET ledger, and preserves original manifest counts as reported metadata.
- THIRD_PARTY distinguishes an external source from a bundled component being analyzed,
  including Smacker, and limits license statements to rights the project can grant.
- FAME's full ten-entry shipped name/score table and redundant name list are omitted
  from the public descriptions. The wire layout, count/size measurements, field-width
  evidence and uncertainty remain. The old table plus the published layout and zero
  tails reconstructs 228 bytes with SHA-256
  `1845f7728c6ee8b27c6a7bf086f0e09481412468705df3386bd922a8db1cbe9d`;
  retaining the digest does not redistribute that payload.
- FAME's literal record-stride and name-length instruction excerpts are replaced with
  descriptions of the same operations. No new experiment is asserted; its eight claim
  IDs, status values and confidence qualifications are preserved.
- The reader receives a minimal Go module file so the documented package invocation
  does not depend on a parent research checkout. No `require` dependency is added.
- A non-mutating publication check and synthetic tests report candidates and broken
  references without treating a regex result as evidence of infringement or clearance.

## Not completed or changed

This is not a whole-corpus sanitization. Other claim ledgers, retractions and the long
registry still require a semantic review. The R2-ASSET public evidence gap is recorded,
not repaired by copying an unreviewed private file or inventing stub claims. Component-
specific legal bases, purchase/terms verification and complete historical-object scans
are not certified by this pass.

No existing license text, factual research source, original-game rule or retraction
record is changed. The private research record is unchanged and complete. Publishing an
edited first tree is not a withdrawal of anything the private record holds, and it does
not retroactively establish a clean room.

## Semantic corrections after the first pass

A review of the first pass against the raw export found semantic drift that the
publication rule does not permit, and this edition corrects it. `claims/video.md`
regains one status cell (`VIDEO-045`, which clause the amendment refuted), the
`Medium` grade for EN/RU sample parity (`VIDEO-048`), the decoder-internal Unknown
(`VIDEO-034`) and the residual-scope statements of `VIDEO-SFX-014`, `VIDEO-SFX-016`,
`VIDEO-MUSIC-011`, `VIDEO-MUSIC-012`, `VIDEO-049`, `VIDEO-050` and `VIDEO-051`; an
Unknown the raw export did not carry was removed from `VIDEO-045`. Exact arithmetic
returns as functional knowledge: the Russian lowercase helper's two byte ranges in
`formats/text/format.md`, and the member count, archive length, archive digest and
duration arithmetic of `VIDEO-MUSIC-001` and `VIDEO-MUSIC-003`. Canonical format pages
again cite every claim id literally instead of abbreviating a family as a range, and
`formats/res/format.md` regains the citations it dropped. `RES-HDR-001` stays withheld
and [SOURCE.md](SOURCE.md) says so.

None of this reinstates an instruction listing, a shipped-content table or quoted game
text. Confidence levels, statuses and Unknowns are restored to what the research record
states, never raised.
