# DIALOGUE — the window a mission's script speaks through — specification (partial)

Level 3. Promoted, evidence-backed claims only. Ledger: `claims/dialogue.md`.

**Status: partial (◐).** Read at instruction level: both announcement senders and their shared
transport, the client dispatcher's whole opcode map, the session window's `0x433` arm with both
of its negatives, the panel class end to end (13 functions, 0 orphan bytes), the tag scan's
vocabulary, and the text control's wrap-and-clamp. `DLG-MSGNUM-025`…`DLG-SOUND-028` add the mission number's own field, the
eight sex-and-class conditionals and the `sound=` consumer. Not specified: what the four
npc-flag conditionals *do* once matched; the `%s` of the speech name; the lose chain past
`0x41e`; whether text past the clamp is reachable by any route; where the campaign's list
of valid mission numbers is loaded from.

This is not a file format. It is the seam between the script runtime (`formats/trigger`), the
mission (`formats/mission`) and the container (`formats/res`). The glyph rule — byte → glyph →
pixel and the font atlases — is specified by `TEXT-CONV-001` (partially retracted)…`TEXT-TILDE-009` and is deliberately absent here.

## The one thing a consumer must not get wrong

**An announcement whose text does not ship is silent, and an announcement that arrives while a
dialog is open is discarded.** Neither is an error path: both are ordinary returns with no
message, no fallback and no retry. A consumer that logs, falls back, or queues shows the player
something the original never showed — and on the English release that is 19 of the campaign's
242 announcements.

## The path

```
script instant 2          FUN_004ea1b3   opcode 0xb6 into the static message at 0x00609c38
mission outcome           FUN_004ea2f6   opcode = argument (0xb4 lose, 0xb5 win), SAME object
      |                   FUN_004e74fe   broadcast to every recipient at this+0x18b8
      v
client dispatcher         FUN_004104e8   pool 0x005f22d0; opcode-3 bounded 0xbb;
                                         index 0x0041879b[188] -> jump 0x004186e3[46]
      0xb6 -> PostMessage(hwnd, 0x433, number, 0)
      0xb4 -> PostMessage(hwnd, 0x433, 0xff,   0)      <- the sentinel
      0xb5 -> PostMessage(hwnd, 0x430, 0,      0)
      v
session window            FUN_00473110   0x416..0x485 via 0x004751b8[112] -> 0x004750d0[58]
      0x433 arm:  campaign+0x414 = wParam            (ALWAYS, sentinel included)
                  wParam == 0xff      -> post 0x431  (the lose panel)
                  campaign+0x3dc & 8  -> return      (a dialog is open: DROP)
                  open main\text\battle\m<mission>\event<NN>.txt
                  absent              -> return      (silence)
                  present             -> FUN_004217be(relative name)
```

### Where `<mission>` and `<NN>` come from

`<NN>` is the script parameter, carried as a **dword** the whole way — written
`004ea1f4 MOV dword ptr [EDX + 0xa],EAX` and read back `004105d9 MOV ECX,dword ptr [EAX + 0xa]`
— and stored in `campaign+0x414`. `<mission>` is **`campaign+0x660`** and never comes from the
packet. That one field also formats `%d.alm`, `main\text\battle\m%d\briefing.txt` and `m%d`,
so the text directory number is the map file number by construction. It is the `+0x118` of the
record embedded at `campaign+0x548`, written only through that record's own setters; nothing in
the image writes displacement `0x660`.

**Reserved value.** `<NN> = 255` is the mission-lost sentinel and reaches the lose panel, not a
text window. `<NN>` is otherwise unconstrained; the shipped corpus uses 0..25.

**Mode.** The `0x433` arm carries no `campaign+0x6bc` test, so it runs in any session. Only the
map load is mode-gated: `campaign+0x6bc == 2` builds `<n>.alm` from `campaign+0x660`, and any
other value loads the map named by the CString at `campaign+0x6b4`.

## The window

The failure panel is **not** the following dialogue class. Its constructor `00446de2` installs
vtable `00598b78`; slot `+0x48` is `00447063`, forwarding the base result/close mechanism.
It offers Exit to Main Menu (`0x445`) and Load Game (`0x446`). Base close posts `0x44c`; the
frontend matches the stored failure-panel pointer `+0x110`, then chooses `0x41e` teardown/menu
or `0x418` save selection. There is no failure-panel `00446d35 -> 0x41d` hop: that method is
in a neighboring class (`DLG-PATH-002` (amended), `MISSION-DEFEAT-046`). Both outcome panels set bit 8;
the campaign idle gate pauses stepping while shown (`SESS-DEFEAT-065`).

`FUN_004217be` is the class's **only** constructor. Rects are `{left, top, right, bottom}`.

```
panel     id  9   ( 30,120)-(610,360)   0x84 bytes, vtable 0x0059b798
portrait  id 12   ( 30, 54)-(118,168)   only when panel+0x7c
text      id 10   (128, 36)-(428,172)   with a portrait
          id 10   ( 48, 36)-(428,172)   without one
button    id 11   (200,172)-(280,198)   label = main.txt line 77, command 0x46f
```

Six call sites, all six shipping their resources: mission events (`battle\m%d\event%02d`), the
inn's NPCs (`inn\NPC\npc%02dm%d`), the mercenary hall (`inn\mercenary\npc%02d`, `…\npc35`), the
shop (`shop\npc31m%d`) and the training hall (`training\npc34m%d`).

## The speaker's figure

Child 12's picture is built by `FUN_00421b46(npcId)` and comes in two forms, chosen by the
speaker's own flags word. `DLG-FIGURE-020`, `DLG-FIGURE-021`, `DLG-SPEAKER-022`.

```
speaker = FUN_00422399(npcId)              a live actor matching the npc<n> section's
                                           Flags predicate, or 0
          ?: FUN_00421f46(npcId)           else a drawable synthesised from npc.reg

surfaces  FUN_00429520(0x58,0x6c)  =  88 x 108   returned to child 12
          FUN_00429520(0xa0,0xf0)  = 160 x 240   composition canvas

speaker+0x18c & 0x11 == 0   flat: graphics\infowindow\<InfoPicture><face>.bmp into the canvas
speaker+0x18c & 0x11 != 0   figure: speaker->vt+0x80(0, canvas, 0)

blit      72 x 96 window of the canvas -> (8,7) of the 88 x 108 surface
          source top 0x90 - PortraitY1, or (0x24,0x8c)-(0x6c,0xe8) when the key is absent
```

`vt+0x80` is `FUN_0045ed10`, the same figure compositor the world view uses. It is `RET 0xc`;
its three parameters are the picture surface, an optional stencil surface and an optional third
surface, and none of them selects layers. The dialogue passes a null stencil, so the click-map
pass does not run and the colour pass draws every occupied slot, the head slot included.

A synthesised speaker's twelve visible-equipment slots are zero, so its figure is the face sheet
alone. A live speaker's figure carries whatever it is wearing at that moment.

## Lifecycle

```
show      vt+0x80 -> advance one part; the RETURN IS IGNORED, so a file with no
                     part 1 opens on the literal "Nothing to say"
input     button 0x46f, or key 0x0d (RETURN), or key 0x1b (ESCAPE) -> the same command
          -> another part remains: set the text, refresh the face, STAY OPEN
          -> none remains: post 0x445 -> vt+0x84 (free the text) -> post 0x44c
                                      -> FUN_004757b0 clears campaign+0x3dc bit 3
timer     none. 13 functions in the class, not one references 0x113.
overlap   impossible: FUN_00476810 sets bit 3 on show and the 0x433 arm drops on it.
```

## Content

The file is scanned, not parsed: find `<` … `>`, lowercase the body, substring-search it.
Fifteen literals, in test order:

```
part=%d  npc=  iamfemale  iammale  iammage  iamfighter  npcalive=  npcdead=
female  male  mage  fighter  sound=  tune=  tips=
```

`part=%d` is a **substring** test, so `part=10` satisfies a search for part 1 — the first
matching tag in file order wins. Shipped data never trips it (0 of 513 EN / 518 RU parts).
`npcalive=`, `npcdead=`, `sound=` and `tune=` are never used on either root.

Two different `npc` tests decide the face: the constructor's `Find("npc")` over the **whole
lowercased file** picks the layout once, and the per-part `Find("npc=")` refreshes the portrait.

### The eight sex-and-class conditionals

Each **abandons the tag and resumes the scan** when the tag and the bit disagree; the scan is
one forward pass, so another tag declaring the same part can still match, and only when no
tag matches does the routine return 0 and the window close. `iam*` tests the player's own hero;
the bare four test the speaker, and are skipped entirely unless the npc section has a `Start`
key and a speaker resolves.

```
actor+0x18c bit 0x4 = female        bit 0x2 = spellcaster
iamfemale   player  reject if clear    female   speaker  reject if clear
iammale     player  reject if set      male     speaker  reject if set  (and Find("female") == -1)
iammage     player  reject if clear    mage     speaker  reject if clear
iamfighter  player  reject if set      fighter  speaker  reject if set
```

**The literals nest and an arm that does not reject falls through.** `Find` is a substring test,
so `iamfemale` also satisfies `female` and `male`, `iammale` satisfies `male`, `iammage`
satisfies `mage` and `iamfighter` satisfies `fighter`. Only the `male` arm guards against it,
and only against `female`. A part tagged for a female player therefore also requires a female
speaker, and an `iamfemale`/`iammale` pair read to a female player by a male speaker matches
neither tag. Nine bodies on each root carry `iamfemale` and nine carry `iammale`.

### Sound

`sound=` is an out-parameter, not an effect: the scan takes the value to the next `"` or `;` and
hands it back. The pager plays it as `speech\` + the resource path's directory prefix + the
value + `.wav`. When the tag is absent the pager composes the name itself, and **which format it
uses is decided by the resource's leaf name**: a leaf beginning `event` gives `npc%02de%sp%d`,
any other leaf gives `%sp%d`. That is the only family test in the whole surface that reads the
file name rather than the directory.

## Measurement

The text control wraps into a line array and computes

```
ctrl+0x8c = min( (bottom - top) / pitch , lineCount )        pitch = fontHeight + 2
```

It has a scroll setter (`vt+0x80`, notify `0x46d`) but this window builds it no scrollbar and
answers no `0x46d`. The 300 px width is not a fixed character count: wrapping calls the text
measurer (`TEXT-API-007`), whose per-glyph advance is `.dat[glyph] + spacing`
(`SPR16A-FONT-018`). The visible-line formula above uses the resulting `lineCount`, not a fixed
pixel-to-character conversion.

## Language

`rom.exe` is byte-identical in both roots. Everything above is compiled. The two
language-carrying inputs are the event files and `main.res::text/main.txt`, which has **274
lines on both roots** — the indices 77/140/141 are compiled constants and a root with a
different line count would break them. The event corpora differ by three files, all RU-only:
`m100/event09`, `m130/event07`, `m150/event10`.
