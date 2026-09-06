# TEXT — public functional specification

Level 3 functional edition. Promoted claim families live in
[`claims/text.md`](../../claims/text.md). Sprite/font framing is documented separately
under [`formats/spr16a`](../spr16a/format.md).

This page publishes the byte-level interoperability rules needed to display and accept
text. It deliberately omits copied UI prose, complete shipped string tables and other
reconstructable content inventories.

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

The game loads multiple CRLF-delimited text resources into positional tables. Entries
are addressed by table-local or global numeric indices depending on the consumer.
Loading is byte-preserving; conversion happens when text is displayed or entered, not
when the resource file is parsed. — `TEXT-STRTAB-023`

Important compatibility consequences:

- table ordering is semantically significant;
- inserting/removing a line in a positional table can renumber later entries;
- consumers may apply their own numeric formatting and visibility rules after selecting
  a string entry;
- path-like tables and UI-prose tables use the same basic storage mechanism but should
  not be assumed interchangeable;
- a consumer must not infer a language-independent semantic key from the displayed
  wording alone.

The public edition intentionally does **not** reproduce the complete shipped filename
list, line counts, global offsets, localized phrases or other tables that would recreate
substantial parts of the game's textual resources. Those are corpus evidence, not the
text-format grammar.

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

## Character-generation and UI labels

ROM1 mixes three presentation mechanisms:

1. positional string-table entries;
2. captions baked into bitmap resources;
3. pictorial controls with text used only for hover/help or other secondary surfaces.

A replacement engine must not assume every visible caption has a corresponding string
entry, or that every descriptive string is persistently drawn next to its control.
Specific shipped wording and bitmap caption contents are intentionally omitted from the
public functional edition. — `TEXT-CHARGEN-027`, `TEXT-CHARGEN-028`, `TEXT-CHARGEN-029`,
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

## Corpus observations versus rules

The private evidence includes complete EN/RU byte censuses, atlas-coverage measurements
and comparisons of localized tables. Those measurements support the transforms above
but are **not themselves wire-format requirements**. This public page therefore keeps
only the constraints needed for an independent consumer:

- single-byte resource storage;
- selector-dependent input/display conversion;
- direct font-record indexing;
- positional string tables;
- caller-owned formatting/markup semantics.

Exact shipped prose, full content inventories, hashes and reconstructable translation
matrices stay in private evidence unless separately reviewed for publication.

## Unknown / bounded areas

The published rules do not establish:

- a Unicode encoding for the original resources;
- safe behaviour for arbitrary malformed byte values or out-of-range font indices;
- semantic equality of every EN/RU string-table entry;
- a universal UI labelling mechanism;
- coverage of strings embedded in every possible non-text resource type.

Do not fill these gaps from a third-party reimplementation or by copying game-resource
text. Preserve them as Unknown until independently researched and reviewed.
