# ROM2 format specifications (level 3)

**☑ Identity survey complete for the seven declared surfaces below.** The survey
uses the preserved Russian ROM2 root and compares container/record layout,
located extensions and bounded nonidentity with ROM1. A difference is a result;
it does not have to become byte identity to close a survey. The session row
specifies the existing ROM2 header only, with no ROM1 protocol-equivalence claim.
This is not a complete ROM2 decoder, game or networking implementation.

Each page uses R2-ASSET/R2-SESSION claims as evidence. ROM1 claims may appear as
explicit comparison references, never as proof of ROM2 behavior. The main index
is [formats/README.md](../README.md). Structural compiled-code comparison remains
the separate `claims/rom2-engine.md` ledger.

| Surface | Established survey result | Explicit boundary | Page |
|---|---|---|---|
| RES | Eleven containers tile under the compared grammar,0 violations; client/server directory and size differences measured (`R2-ASSET-001`, `R2-ASSET-012`) | Same-size payload equality is not measured | [RES](../rom2-res/format.md) |
| ALM |83 complete record chains with the type0 overhang; per-type layouts, dispatcher, metadata growth and type6 widths mapped (`R2-ASSET-002`, `R2-ASSET-003`, `R2-ASSET-017` through `R2-ASSET-024`, `R2-ASSET-027`, `R2-SESSION-010`) | New-type field meanings and the writer's declared-size split remain Unknown | [ALM](../rom2-alm/format.md) |
| Data.bin |Both client/server A–H groups tile exactly; only C's raw block widens10→14; schema and2556-byte delta located (`R2-ASSET-029`, `R2-ASSET-030`, `R2-ASSET-031`) | `templates.bin` is a bounded divergent input at offset35004, not a proved Data.bin instance (`R2-ASSET-032`); wider table semantics remain open | [Data.bin](../rom2-databin/format.md) |
| Inline REG |19/19 nested registries tile exactly (`R2-ASSET-006`) | Value/kind semantics and key-set identity are outside the layout survey | [REG](../rom2-reg/format.md) |
| Sprite / palette |623/623 .16a and159/159 palettes match framing; .256 has1916 exact,8 empty stubs,5 ROM1-identical residues; all3 .16 frame streams validate,1 tiles and2 have ROM1-identical residues (`R2-ASSET-007`, `R2-ASSET-008`, `R2-ASSET-009`) | Pixel/RLE and colour-value re-derivation excluded | [Sprite/palette](../rom2-spr/format.md) |
| Text |67 main text files plus patch; container extension and byte-range nonidentity measured,66.3617% overlap (`R2-ASSET-010`, `R2-ASSET-011`) | Encoding name, converter and font mapping remain open | [Text](../rom2-text/format.md) |
| Session header |Eight-byte frame and four consumed fields specified;142-byte payload/150-byte total bounds reconciled (`R2-SESSION-003`) | Payload/opcodes, protocol continuation and ROM1 equivalence excluded | [Header](../rom2-net/format.md) |

The survey's input population is one preserved ROM2 locale. No English ROM2
behavior is inferred. Future format decoding or networking research needs its
own question and authority; these checkmarks do not schedule it.
