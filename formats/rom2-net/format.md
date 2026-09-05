# ROM2 session wire frame — identity survey

Level 3. Promoted, evidence-backed claims only. Basis: `R2-SESSION-003`. No ROM1
cross-reference: no ROM1 wire-protocol counterpart was surveyed by the
experiment that produced this page.

Seen as: the 8-byte record header `CBufferManager::ReceiveData` reads at the
start of every inbound socket record, in both `allods2.exe` and `a2server.exe`.

## Result

**One 8-byte header, four fixed-offset fields.** Relative to the header's own
start: `+0` is a `u16` payload length, observed valid range 1..142; `+4` is a
`u8` codec selector (0 skips decoding); `+5` is a `u16` codec input size,
passed to the codec as its own input-length argument; `+7` is a `u8`
passthrough byte the codec never reads, copied verbatim into the decoded
record. Both binaries read these four fields identically, once each
transport's own record-copy displacement is accounted for. `R2-SESSION-003`.

**The 142-byte cap is one bound restated, not a second bound.** A second
transport on the same binary (`a2server.exe`'s DirectPlay message handler)
checks its own total message length against 8..150 bytes; 150 is the same
142-byte payload cap plus the 8-byte header that transport has not yet split
off the total, not an independently chosen constant. `R2-SESSION-003`.

## Not yet surveyed

The record's own payload grammar past the codec/passthrough fields; whether
any wire message carries a type or opcode field at all, downstream of the
queue this header feeds; a ROM1 counterpart, if one exists.
