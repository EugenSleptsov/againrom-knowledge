# SAV Diary

[Format reference](format.md) · [Player](player.md) · [Actors](actors.md)

## Wire fields

| Order | Wire | Runtime member |
|---:|---|---|
| 1 | `Count(n)`, raw `4*n` | Embedded CDWordArray at `+04` |
| 2 | `Count(m)`, raw `2*m` | Embedded CWordArray at `+18` |
| 3 | u32 saved-address reference | `+2c` |

Counts are independent. LOAD sizes each array from its own count, reads its
payload, then resolves `+2c` through the saved-address map: hit installs the
pointer, miss/zero becomes null. It enforces neither equal lengths nor a
complementary element relation. — SAV-DIARY-042, SAV-PLDIARY-054, SAV-847

A Player owns one Diary at `+40`, dispatched directly in its suffix. A
Humanoid/Human can also reference a separate Diary at actor `+1e4` through
CArchive. Actor-owned Diaries default to a null owner; the two ownership
routes are not interchangeable. — SAV-HUMAN-043, SAV-846

## Initialization and index domain

A fresh Diary sizes both arrays from the Units collection count and initializes
dword 0 and word 1024. That capacity does not establish a single actor index
domain: base Unit construction writes a Units ordinal into actor `+0c`,
while Human construction writes a Humans ordinal. The mutation rejects an
index above 63 only when actor virtual `+30` is nonzero: false for Unit,
true for Humanoid/Human. — SAV-667, SAV-668, SAV-844

The array accessors perform pointer arithmetic without a length check. No
safe malformed count/index relationship follows from the field widths.
— SAV-842

## Mutation and notification

The located mutation first decrements a nonzero word, then increments the
dword only while its unsigned old value is at most 16. From defaults,
17 admitted calls yield `(17,1007)` and 18 yield `(17,1006)`. Thus
`word=1024-dword` is not an invariant. LOAD must preserve both arrays independently. — SAV-842, SAV-847

The located caller is actor teardown under removed actor `T+54==16`.
With `S=T+40`, it requires S type word `+0e` in 33..63 and signed health
`+94>=0`. The receiver is `S->Player(+14)->Diary(+40)` and argument is T.
Standard damage can record its source in T `+40` before death; this does not
establish a universal lethal-blow/every-kill interpretation.
— SAV-843

Non-null Diary `+2c` admits notification when the dword increments to
2,4,6,8,10,12,14,16. The downstream packet builder reads the supplied Player's
own array and writes 17 words under opcode 186:

1. For `i>=64`, obtain `typeID` and `face` from Units row i.
2. For `typeID=64..80`, select output word `typeID - 64`.
3. OR `min(lowByte(dword[i]>>1),7) << (4*(face-1))` into that word.
4. If `Player+68>10`, select all `0xffff` instead.

Null owner suppresses notification without suppressing mutation. Receiving UI
and visible cadence remain Unknown. — SAV-845

## Remaining scope

Additional actor-owned Diary consumers remain Unknown. Other uses of the shared word accessor need not
receive a Diary; one argument-taking path receives a fresh stack collection.
— SAV-846, SAV-854, SAV-855

Indirect/computed/inlined access and unclassified receivers remain open.
The word's gameplay meaning, full event attribution and the first post-LOAD
consumer are Unknown. The wire grammar requires neither an owner self-reference
nor equal/complementary arrays. — SAV-843, SAV-847, SAV-854, SAV-855
