<a id="rom2-text--identity-survey"></a>

# ROM2 text resources

ROM2 has CRLF-delimited string resources under `main.res:text/` and
`patch.res:patch.txt`. Their container form matches the positional resource
family; ROM2's byte-to-glyph conversion and named character encoding remain
Unknown. — R2-ASSET-010, R2-ASSET-011

<a id="container-result"></a>

## Resource families

The fifteen ROM1 table names from `main.txt` through `credits.txt` remain
present, together with patch.txt. Additional resources include numbered
`missionNN.txt` files and `globalmap.txt`, `help.txt`, `itemserv.txt`,
`quest.txt`, `town.txt`, `docs/1.txt`. The installed mission numbers are sparse;
do not derive a contiguous list from the filename pattern. — R2-ASSET-010

<a id="encoding-result"></a>

## Stored-byte domain

The preserved ROM2 text-like archive nodes (`.txt`, `.ini`, `.lst` and
extensionless) contain high bytes in `0xe0..0xff` and 31 values of
`0xc0..0xdf` (`0xda` absent). They contain none in `0x80..0xbf`.
This is a property of that input population, not a valid-byte restriction.
— R2-ASSET-011

Only the `0xe0..0xef` portion overlaps the two source ranges of the ROM1
Russian display transform. ROM1 TEXT-CONV-001 is partially retracted and is
not a ROM2 conversion rule. A matching byte range does not establish a
matching encoding or atlas. — R2-ASSET-011, TEXT-CONV-001

<a id="not-yet-surveyed"></a>

## Reading and rewriting

Keep source strings as bytes and preserve resource/table ordering. Do not
convert using a guessed code page or reverse a presumed glyph transform.
Native ROM2 load ordering, exact delimiter-edge behavior, font selection,
conversion and independent text-writer acceptance remain Unknown.
— R2-ASSET-010, R2-ASSET-011
