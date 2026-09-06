# ROM2 text — identity survey

Level 3. Promoted, evidence-backed claims only. Basis: `R2-ASSET-010` (container),
`R2-ASSET-011` (encoding). Cross-reference only, not evidence: ROM1's own
[`formats/text/format.md`](../text/format.md) (`TEXT-STRTAB-023`,
`TEXT-CONV-001`, partially retracted).

Seen as: `main.res:text/*.txt` (CRLF-delimited string tables) and
`patch.res:patch.txt`.

## Container result

Preserved and extended. All 15 of TEXT-STRTAB-023's named files (`main.txt`
through `credits.txt`) are present under `main.res:text/`, and `patch.res:patch.txt`
is present — matching ROM1's own load-order population exactly by name. The
directory is also substantially larger: **67** files live under `main.res:text/`
(**52** beyond the 15 ROM1 names found there, not 54 beyond 14), plus the 1 file
under `patch.res` (68 corpus-wide). Of the 52 beyond ROM1: **46** are
`missionNN.txt` campaign files, `mission10.txt` through `mission110.txt` (not
every number in between present), and **6** are new non-mission names:
`globalmap.txt`, `help.txt`, `itemserv.txt`, `quest.txt`, `town.txt`,
`docs/1.txt`. `R2-ASSET-010`.

## Encoding result

Does not stay inside TEXT-CONV-001's (partially retracted) two source blocks `{0x80..0xAF, 0xE0..0xEF}`
the way ROM1's own text does, and its in-domain share is not spread over both
blocks the way ROM1's is. Measured fresh, with this survey's own tool and its
own population definition (`.txt`/`.ini`/`.lst`/extensionless nodes across all 11
containers — narrower than TEXT-FIT2-013's population, which also folds in `.reg`
and `data.bin` embedded strings), as a same-tool positive control on both ROM1
roots first: ROM1 EN scores 1/1 high bytes in-domain (100.0000%) — with only 1
high byte total, this control carries no discriminating information and is
reported only for completeness. ROM1 RU scores 87446/87446 (100.0000%,
independently corroborating TEXT-FIT2-013's conclusion): 62933 of those in
`0x80..0xAF` and 24513 in `0xE0..0xEF` — RU's in-domain share spans BOTH blocks.
ROM2 scores only 210405/317058 (**66.3617%**) in-domain, and **0 (zero)** of its
317058 high bytes fall in `0x80..0xAF` at all (against RU's 62933) — ROM2's
entire in-domain share of 210405 sits inside the single shared `0xE0..0xEF`
window. The remaining 106653 out-of-domain bytes (33.6%) populate 47 of the 64
values in the two blocks the converter never touches — all 16 of `0xF0..0xFF`,
31 of 32 of `0xC0..0xDF` (only `0xda` absent), and **none** of the 16 values
`0xB0..0xBF`. `R2-ASSET-011`.

## Not yet surveyed

Which converter (if any) ROM2's client applies to these bytes, or which font
atlas draws them — this is a byte-range corpus measurement, not an
`allods2.exe` instruction read. Which specific 8-bit encoding the populated
`0xC0..0xFF` shape corresponds to is not named; only the measured byte-range
shape is reported.
