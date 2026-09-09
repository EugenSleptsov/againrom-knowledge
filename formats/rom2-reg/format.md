<a id="rom2-inline-registry-ya1-nested--identity-survey"></a>

# ROM2 inline REG store (`&YA1`)

Nested ROM2 registries use the inline REG envelope. It is distinct from the
tail-registry RES envelope despite the shared magic. Resource context and
record geometry select the parser. — R2-ASSET-006

<a id="result"></a>

## Layout

| Position | Width | Data |
|---|---:|---|
| 0 | 24 bytes | Six-dword header; magic `&YA1`, root fields, record count at 0x10 |
| 0x18 | 32×R bytes | Record array |
| 0x18+32×R | 4 bytes | Pool byte length |
| 0x1c+32×R | poolLen bytes | Pool |

`registryEnd = 24 + 32*R + 4 + poolLen`. This framing covers the preserved
nested registries. The installed `data/ai.reg` and `data/map.reg` lengths are
156 and 1020 bytes in both world archives; equal lengths do not prove equal
values. — R2-ASSET-006

## Read and write order

Read the six header dwords, R records, u32 pool length and pool bytes.
The corresponding structural emitter writes those same regions in order.
Raw record and pool values must remain opaque where ROM2 type/lookup semantics
have not been established. [ROM1 REG](../reg/format.md) supplies a layout
cross-reference, not an authority for unverified ROM2 value behavior.
REG-REC-032's ROM1 record layout is retained; its lookup clause is partially
retracted and cannot be imported as an unconditional ROM2 rule.
— R2-ASSET-006, REG-FMT-031, REG-REC-032

<a id="not-yet-surveyed"></a>

## Unknowns

ROM2 record kinds, value interpretation, key-set differences, native lookup,
sort and application/writer acceptance are unspecified. Recognition by a
nested `&YA1` payload is broader than a `.reg` extension filter.
