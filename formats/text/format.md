<a id="text--public-functional-specification"></a>

# Text bytes, tables and input

ROM1 stores one-byte strings. The active language selector controls input and
display conversion; font indexing follows conversion. Positional string
tables retain source byte order. Glyph framing and advances are defined in
[SPR16A](../spr16a/format.md). — TEXT-CONV-001 (partially retracted),
TEXT-LANG-002, TEXT-STRTAB-023, SPR16A-FONT-018

## Encoding model

ROM1 uses a one-byte text pipeline. The active language selector chooses whether the
Russian conversion rules are applied.

For display in the Russian mode, the byte transform is:

```text
0x80..0xAF -> byte + 0x30
0xE0..0xEF -> byte + 0x10
otherwise  -> unchanged
```

In the other selector state the display transform is the identity. The transformed byte
selects a font record by:

```text
record = uint8(transformedByte - 0x20)
```

This is a byte operation, not Unicode decoding. A compatible implementation may expose
Unicode internally, but conversion to/from the original resources must preserve the
original one-byte semantics. — `TEXT-CONV-001` (partially retracted), `TEXT-LANG-002`,
`TEXT-INDEX-003`, `TEXT-SEL0-012`

## Input conversion

In Russian mode the keyboard/input conversion into stored bytes is:

```text
< 0x80      -> unchanged
0xC0..0xEF  -> byte - 0x40
0xF0..0xFF  -> byte - 0x10
otherwise   -> unchanged
```

The game's Russian lowercase helper additionally maps the two Cyrillic uppercase ranges
before falling back to the ordinary single-byte lowercase operation:

```text
0x80..0x8F -> byte + 0x20
0x90..0x9F -> byte + 0x50
otherwise  -> ordinary single-byte lowercase
```

These transforms are byte functions; no multibyte encoding is involved.

The display transform is not injective in Russian mode: distinct stored bytes can map to
the same font record. A reimplementation should therefore keep stored text bytes and
rendered-glyph identity conceptually separate rather than normalizing the source data. —
`TEXT-DOM-010`

## Font indexing and bounds

The original font access path indexes records directly after the `-0x20` transform and
does not provide a general substitute-glyph/clamp rule. Therefore callers are
responsible for splitting control characters and supplying bytes valid for the chosen
font.

A compatible implementation may add defensive bounds checks for safety, but that is an
implementation divergence and must not be mistaken for a property of the original file
format.

Font record counts, cell geometry and sprite framing are documented in the SPR16A
specification rather than duplicated here. A font is two resources — the glyph sprite
set and a separate per-record advance table — so the drawn cell width is not the advance
a text measurer must add. — `SPR16A-FONT-018`

## Markup byte

The tilde byte has markup semantics in the relevant text-drawing path:

- a doubled tilde represents a literal tilde glyph;
- a lone tilde activates the original rule/markup behaviour and is not measured like a
  normal character advance.

A text measurer must mirror the draw path's markup handling or widths can diverge.

## Resource string tables

[Pointer-hover help](hover.md) describes the common delayed display route,
control-specific sources, line layout and the bounds of the control inventory.
— TEXT-HOVER-048, TEXT-HOVERSET-049, TEXT-HOVERCHAR-050,
TEXT-HOVERROOM-051, TEXT-HOVERTEXT-052, TEXT-HOVERPAINT-053

The game loads multiple CRLF-delimited text resources into positional tables. Entries
are addressed by table-local or global numeric indices depending on the consumer.
Loading is byte-preserving; conversion happens when text is displayed or entered, not
when the resource file is parsed. — `TEXT-STRTAB-023`

Table rules:

- table ordering is semantically significant;
- inserting/removing a line in a positional table can renumber later entries;
- consumers may apply their own numeric formatting and visibility rules after selecting
  a string entry;
- path-like tables and UI-prose tables use the same basic storage mechanism but should
  not be assumed interchangeable;
- a consumer must not infer a language-independent semantic key from the displayed
  wording alone.

## Character-name entry

The original character-name control stores a bounded one-byte string and applies the
input conversion before appending accepted bytes. Backspace is handled as editing rather
than a stored character; control bytes below the printable range are rejected in the
reached input path. — `TEXT-NAMEIN-024`

Two different typed byte sequences can therefore become visually identical under the
Russian display transform: composing the input conversion with the display conversion
leaves a bounded set of reachable glyph records that more than one keystroke sequence
can reach. A compatible save/editor implementation should preserve the stored bytes
rather than replacing them solely from rendered appearance. — `TEXT-COLL-025`

## Save-label entry

The SAVE/LOAD chooser's list-item label is not the character-name control's bounded,
filtering class above. Its recognized producer path is a generic list-control item-text
copy that filters no byte value other than the NUL terminator, distinct from the
character-name entry's rejection of control bytes below the printable range.
— `TEXT-SAVELABEL-057`

The character-name control's own vtable is referenced nowhere in the executable image by
a full instruction-reference search, ruling out that specific class as the save-label
producer for any construction reached through the literal-vtable-store idiom this
codebase uses elsewhere; a class constructed by some other means is not excluded by this
search. — `TEXT-SAVELABEL-055`

The byte-indexed display-conversion selector documented above for other text surfaces is
not called, directly or through its only wrapper, by any traced save-label chooser code
path. — `TEXT-SAVELABEL-054`

Neither the byte-indexed selector's own draw functions nor this image's GDI text-out
import surface is reached by any traced save-label chooser code path either, and the MFC
GDI-wrapper vtables that would carry the latter are themselves never constructed. By
elimination among these named mechanisms, the field's pixels are consistent with native
Win32/MFC list-control default painting outside this executable's own code — an
elimination among catalogued candidates, not a positive trace, and it does not exclude an
uncatalogued in-game draw routine reached only by virtual dispatch this search cannot
enumerate. Because the two lawful executables are the same file, EN and RU cannot differ
in a mechanism this image does not exhibit for this field. What the draw path does with a
byte outside 7-bit printable ASCII, and what bounds the field's *drawn* (as opposed to
retrieved) length, both require observing a running original and are not established by
static analysis. — `TEXT-SAVELABEL-058`, `TEXT-SAVELABEL-059`, `TEXT-SAVELABEL-060`,
`TEXT-SAVELABEL-061`

## Character-generation and UI labels

ROM1 mixes three presentation mechanisms:

1. positional string-table entries;
2. captions baked into bitmap resources;
3. pictorial controls with text used only for hover/help or other secondary surfaces.

A replacement engine must not assume every visible caption has a corresponding string
entry, or that every descriptive string is persistently drawn next to its control.
— `TEXT-CHARGEN-027`, `TEXT-CHARGEN-028`, `TEXT-CHARGEN-029`,
`TEXT-UI-032`, `TEXT-UI-033`, `TEXT-UI-034`, `TEXT-UI-035`, `TEXT-UI-036`, `TEXT-UI-037`,
`TEXT-UI-038`, `TEXT-UI-039`, `TEXT-UI-040`, `TEXT-UI-041`, `TEXT-UI-042`, `TEXT-UI-043`,
`TEXT-UI-044`, `TEXT-UI-045`, `TEXT-UI-046`, `TEXT-UI-047`

Where a control does take its caption from the string tables, the reached constructors
read fixed global indices and copy the bytes into storage the control owns, rather than
borrowing the loader's pointers; the control's destructor releases that storage. The
pressed and unpressed paint branches then select presentation only, not a different
caption source, and a caption element can be replaced later by a selection writer while
the surrounding numeric strings are refreshed as separate paint arguments. A consumer
must therefore treat caption identity as element position in the control's own array,
not as the rendered wording. — `TOWN-383`, `TOWN-384`, `TOWN-385`, `TOWN-391`,
`TOWN-392`, `TOWN-393`

<a id="corpus-observations-versus-rules"></a>

## Unknown / bounded areas

The published rules do not establish:

- a Unicode encoding for the original resources;
- safe behaviour for arbitrary malformed byte values or out-of-range font indices;
- semantic equality of every EN/RU string-table entry;
- a universal UI labelling mechanism;
- coverage of strings embedded in every possible non-text resource type;
- whether the save-label chooser's byte-to-glyph mapping can differ between EN and RU for
  this field: the two lawful executables are byte-identical (`TEXT-SAVELABEL-061`), so any
  such difference would have to come from outside `rom.exe` — this image's own draw
  mechanism is eliminated by `TEXT-SAVELABEL-058`/`-059`, so the mapping is a property of
  whatever paints the field, not settled by that elimination;
- which native control paints the save-label chooser's field (a catalogued in-image
  mechanism is excluded by elimination, `TEXT-SAVELABEL-058`/`-061`, but the specific
  native control is not identified), what it does with a byte outside 7-bit printable
  ASCII, and what bounds the field's drawn length — all three need a running original.

## Decode and encode sequence

1. Parse CRLF-delimited resource tables as bytes, retaining positional order.
2. Apply the active selector's display transform to each display byte.
3. Handle controls and markup on the appropriate text path, then index
   `uint8(transformedByte-0x20)` in the selected font.
4. Measure using that font's advance sidecar and the draw path's markup rule.

For keyboard input, apply the input transform before appending accepted bytes;
backspace edits the current string. Preserve stored bytes on save or resource
rewrite. Display conversion is not injective and therefore has no unique
inverse suitable for reconstructing original text. — TEXT-STRTAB-023,
TEXT-INDEX-003, TEXT-DOM-010, TEXT-NAMEIN-024, TEXT-COLL-025
