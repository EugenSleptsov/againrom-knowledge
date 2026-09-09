# World and character-panel appearance

[Reference](format.md)

## World appearance (`HERO-APPEAR-040`…`HERO-APPEAR-046`)

The class discriminator above decides *stats*. What decides the **drawn unit class** is a separate
chain, it lives on the client, and it does not read the stat discriminator at all.

What the server sends is not a drawn class. `FUN_004e7de3` pushes, under field-mask bit `0x4000`,
`FUN_00521890(actor)` — nine instructions returning `word[actor+0xe]`, the `typeID` — and then the
face byte `actor+0x4b`. The client parks them in `[0x005cd5c8]` / `[0x005cd5cc]` and assigns the
first to `drawable+0x20`, which is the subscript the frame selector `FUN_0045bf00` uses on the
`units.reg` class array `0x005eb674`.

`FUN_0045f850` then splits on that id:

| id | meaning | what happens to `+0x20` |
|---|---|---|
| `< 0x1a`, `>= 0x40` | a class a map places | left alone — it **is** the `units.reg` `ID` |
| `[0x20, 0x40)` | a player's character | discarded; see below |

The shipped roster occupies `1..27` and `64..80`, so `[0x20,0x40)` holds no class and is free for the
four archetypes `0x21..0x24`. For them the routine computes `typeID - 0x21`, files bit 0 into
`drawable+0x18c` bit 2 and bit 1 into bit 1 — the **same** assignment the chargen record uses, so
bit 2 is the sex axis and bit 1 the mage axis — sets bits 0 and 3, and stores the literal `1`.

`FUN_0045fb00`, the next instruction in the caller, does the real work. It returns at once unless
`+0x18c & 1`. Otherwise it builds a **body name** from a twelve-slot array of visible equipment at
`drawable+0x15c`, which the equipment messages fill from twelve `u16` at `msg+0xc`:

```
name = registryString[(byte[slot0 + 6] & 0x1f) - 1]        slot 0 absent -> index 0
if slot1 != 0            name += "_"
if mage and name == "unarmed"   name = (action == 6 ? "mage_st" : "mage")
```

and then maps the name to the drawn class with a seventeen-arm chain:

| name | id | name | id | name | id |
|---|---|---|---|---|---|
| `unarmed` | 1 | `axeman_` | 8 | `pikeman_` | 13 |
| `unarmed_` | 2 | `axeman2h` | 9 | `archer` | 14 |
| `swordsman` | 3 | `clubman` | 10 | `bowman` | 14 |
| `swordsman_` | 4 | `clubman_` | 11 | `xbowman` | 15 |
| `swordsman2h` | 5 | `pikeman` | 12 | `mage_st` | 24 |
| `axeman` | 7 | | | `mage` | 23 |

Seventeen names, sixteen distinct ids -- `archer` and `bowman` are the only pair that share one. `bowman` ships no sheet on either root and is dead.

**The pixels do not come from `units.reg`'s `File`.** The same routine composes

```
graphics\units\ <dir> \ <name> \sprites.256      -> drawable+0x194
graphics\units\ <dir> \ <name> \spritesb.256     -> drawable+0x198
```

where `<dir>` is `heroes` for a mage, `heroes_l` for a fighter with nothing in slot 7, and otherwise
`material.reg`'s `Material[word[slot7 + 6] >> 12].Path` — sixteen blocks carrying `heroes` for
0–7 and 14 and `heroes_l` for 8–13 and 15. So the class record supplies only the **geometry**:
phases, `Width`/`Height`, `CenterX`/`CenterY`, `Flip` and the three timelines.

`+0x19c` caches the last name and `+0x1ac` the last material index; nothing is reloaded unless one of
them changed.

**It is recomputed while the mission runs**, on the ordinary field-masked state message (opcodes
`0x6c`/`0x6e`/`0x6f`/`0x70`) as well as on the two equipment messages (`0x9c`, `0x76`) and from the
character screen. The server picks between the two equipment opcodes on the same `[0x21,0x40)` test.

**A consumer must therefore not store a drawn class on a player's character.** It must keep the
twelve visible-equipment slots, recompute the name whenever they change, load the sheet by path, and
take the geometry from the resulting `units.reg` record. For a unit a map places it must do none of
this: that actor's drawn class is its `typeID` and nothing derives it.

Claims: `HERO-APPEAR-040`…`HERO-APPEAR-046`, `UNIT-APPEAR-030`. The name list's order and the
second sheet's consumer follow the figure-equipment rules below.

## Figure equipment slots (`HERO-APPEAR-047`…`HERO-APPEAR-055`, `HERO-FIGURE-062`)

**Who sends them.** `FUN_004e873b` sends nothing at all unless `actor->vt+0x30()` — the humanoid
predicate — is true, so a non-humanoid actor's twelve slots stay null for its whole life. With no
named recipient it walks the player list and calls itself once per player.

**Which field becomes which slot.** Both opcodes visit `actor+0x74`, `actor+0x78`, then
`actor+0x198+4i` for `i = 3..12` — `ITEM-EQUIP-006`'s own 1..12 equipment numbering. Wire slot `k`
is equipment slot `k+1`, and an empty slot is sent as 0 rather than omitted.

**The two opcodes do not carry the same bytes.**

```
0x9c  twelve bare u16 at msg+0xc+2k, each = word[item+0x40]
0x76  twelve serialised item records from msg+0x13, each 7 + len bytes:
        rec+0  u16   word[item+0x40]        -> slot+0x06
        rec+2  u16   item+0x42, the count   -> slot+0x10
        rec+4  u8    bit7 = item+0x14 != 0  -> slot+0x08
                     bit6 = item+0x8  != 0
        rec+5  u8                           -> slot+0x09
        rec+6  u8    a length               -> slot+0x0a
        rec+7  ...   len bytes              -> malloc'd at slot+0x0c
      rec+0 == 0 means the slot is empty.  An empty slot is sent as a
      throwaway Item constructed, serialised and destroyed.
```

**The `u16`.** Every bit is consumed, by `FUN_00483c80`, which turns it into a seven-digit name:

```
  bits 15..12  A = v >> 12          the material.reg index  (the world-appearance slot-7 axis)
  bits 11..8   B = (v >> 8) & 0xf   the equipment slot, 1..12
  bits  7..5   C = (v >> 5) & 7     unnamed
  bits  4..0   D = v & 0x1f         the item's Data.bin definition row index

  name = sprintf("%02d%02d%1d%02d", A, B, C, D)
  or     sprintf("%02d%02d%03d",    A, B, v & 0xff)   when B == 14
```

`FUN_00525d60` assembles it from `item+0x46` (A), the item's kind (B), `item+0x45` (C) and
`item+0x0c` (D), each masked to a byte and `OR`ed — so a definition row index above 31 would corrupt
C.

**What the name addresses.** Two trees, and they are different pictures:

```
graphics\equipment\ <figure> \ primary   \ <name>.256    the info window's figure, per slot
graphics\equipment\ <figure> \ secondary \ <name>.256    slots 3, 8, 9, and 7 for a mage
graphics\equipment\ <figure> \ <face>.256                 the head, from drawable+0x24
graphics\inventory\ <name>.16a                            the icon; missing -> "Invalid item weared "

<figure> is one of mfighter, mmage, ffighter, fmage, chosen by (drawable+0x18c & 6) —
the mage bit and the SEX bit together.
```

`FUN_0045ed10`, virtual slot `+0x80` of the drawable class, loops all twelve slots and builds these.
**It is not the world renderer.** Every one of the 928 shipped equipment sheets holds exactly one
frame of 160×240, against a hero body sheet's 129–216 frames of 24×40 to 40×48, and the routine's
identified callers are the info-window routines. In the world nothing is composited over the body:
a held weapon is visible because the world body name changes.

**The name list is a shipped file.** `main\text\heropicture.txt`, loaded into the object at
`0x5ea220` at startup, 25 non-blank lines, identical on both roots, indexed by `D - 1`:

```
 0 unarmed      5 swordsman2h  10 axeman      15 pikeman   20 archer
 1 swordsman    6 clubman      11 axeman2h    16 pikeman   21 xbowman
 2 swordsman    7 clubman      12 mage_st     17 axeman    22 Sonic Beam
 3 swordsman    8 clubman      13 mage_st     18 axeman2h  23 Flame Thrower
 4 swordsman    9 clubman      14 pikeman     19 archer    24 swordsman
```

`bowman` never occurs, which is why the seventeenth world-appearance arm is dead. Two lines name no shipped body
directory. The list is 25 long while D is five bits wide.

**The second sheet is drawn.** `FUN_0045bf00` blits `drawable+0x194` and then `drawable+0x198` at
the same frame index behind one `+0x18c & 1` gate, the second centred from its own width and height
and the class record's `+0x30`/`+0x38`. Its six arguments are `dstX`, `dstY`, frame,
`[0x005eb4a0]`, sun shear, mirror — `TERR-SPR-038`'s push list, so the fourth is
`TERR-SPR-066`'s **shroud level** and the second sheet belongs to the silhouette family.
**`0x26` is not one of them** (figure composition below): `HERO-FIGURE-062` corrects the partially retracted
`HERO-APPEAR-056` by identifying it as the preceding gate's key, not a blit argument.

**A consumer must therefore** keep twelve slots per humanoid actor, fill them from the equipment
slots offset by one, decode each `u16` into the four fields above, use D of slot 0 against
`heropicture.txt` for the body name, A of slot 7 against `material.reg` for the directory, and treat
`graphics\equipment` as a portrait tree that never touches the world sprite.

Claims: `HERO-APPEAR-047`…`HERO-APPEAR-055`, `HERO-FIGURE-062`, `ITEM-APPEAR-023`, `ITEM-APPEAR-024`,
`UNIT-APPEAR-031`, `SPR256-EQUIP-042`.
Not established here: field C beyond the figure-composition clause below; which garment each of slots 2..6 and 8..11 is.

## Figure composition and hit map (`HERO-FIGURE-057`…`HERO-FIGURE-064`)

`FUN_0045ed10` builds up to twelve layers per figure and then draws each of them **twice**, into two
caller-supplied surfaces bound in turn with `vt+0x28`:

```
  pass 1   vt+0x18 = FUN_00428e60   ->  FUN_0044ea40 / FUN_0044db00
           (x, y, frame, palRow, mode)      palette = layer+0x1c + (palRow << 9)
           the picture

  pass 2   vt+0x40 = FUN_00429190   ->  FUN_0044dd40
           (x, y, frame, tag)               tag is the SIXTH argument of the blitter
           a byte-per-pixel stencil: for an opaque run FUN_0044dd40 copies no source
           pixel at all, it stores the tag  (0044dda2 / 0044dda8)
```

**The tag is the equipment slot number**, `drawable index + 1`, over all 28 draw sites; the head
layer takes 0. So the second surface is a map from screen position to equipment slot, which is what
an info window with clickable equipment needs, and it is a second, independent derivation of the
`k+1` slot map in the wire-equipment rules above.

**Draw order is program order and depends on the mage bit** (`0045f1fd`):

```
  mage       p7 head p11 p9 s9 p3 s3 p6 p4 p8 p0 p5 s7
  non-mage   p11 p10 p6 p3 p4 p8 p9 p0 p7 s3 p5 s8 s9   then  p1 or p0
```

Three asymmetries a consumer must reproduce rather than smooth over: index **2** (equipment slot 3)
is drawn by neither half and ships no sheet; index **10** (slot 11) is drawn only by the non-mage
half and ships no sheet; index **8** (slot 9) is **tagged but not colour-blitted for a mage**.

**Which held layer is on top.** `FUN_00483e70` takes slot 0's field D, indexes `heropicture.txt`,
and returns 1 for `bowman`, `archer`, `xbowman`, `axeman2h`, `swordsman2h`, `mage_st` — every
two-handed and ranged body. Predicate 1 → index **0** last; predicate 0 → index **1** last. Field D
is the only field of the `u16` this routine reads.

**The list accessor has no bound.** `FUN_004687f0` is `[[0x005eb3d4] + ([list+0xc] + i)*4]`, six
instructions, no test — and this caller runs only when slot 0 is **occupied**, so `D = 0` indexes at
**−1**. A consumer must clamp at both ends; `D ≥ 26` runs off the other.

**And the same drawable field feeds sound.** `FUN_0045e890` picks one of eight 48-byte voice banks
on the mage bit, the composed-figure bit, **whether slot 0 is occupied**, and the sex bit:

```
  mage bit set                          -> 0x5f1b20   m_mage    / f_mage
  composed-figure bit set, mage clear    -> 0x5f1bf8   mf_hero   / ff_hero
  neither, slot 0 occupied               -> 0x5f1c58   mf_merc   / ff_merc
  neither, slot 0 empty                  -> 0x5f1b98   m_peasant / f_peasant
```

`FUN_004af540` fills them at start-up from `sfx\<bank>\` + `select1/2`, `command1/2`, `retreat`,
`defend`, `idle`, `easy`, `hard`, `die`. The slot-0 condition can therefore
select different voice banks for armed and unarmed humanoids.

Claims: `HERO-FIGURE-057`…`HERO-FIGURE-064`, `UNIT-FIGURE-032`, `ITEM-APPEAR-025`.
Not established: what the effect-list keys `0x26`/`0x30` name; what `[0x005ef990]`/`[0x005ef994]`
are; which garment each of indices 2..6 and 8..11 is.
