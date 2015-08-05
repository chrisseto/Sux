# Escape codes from popular programs.

Note: Quotes indicate raw text

## Bash
* ^[[?1034h
* "bash-3.2$ "

## Nano
* ^[[?1049h
  - Save cursor as in DECSC and use Alternate Screen Buffer, clearing it first.
  - (This may be disabled by the titeInhibit resource).
  - This combines the effects of the 1 0 4 7  and 1 0 4 8  modes.
  - Use this with terminfo-based applications rather than the 4 7  mode.
* ^[[1;54r
  - DECSTBM
  - Set top and bottom lines of a window
* ^[(B
  - setusg0
  - Set United States G0 character set
* ^[[m
  - SGR0
  - Turn off character attributes
* ^[[4l
  - Reset Mode: Replace Mode
* ^[[?7h
  -  Wraparound Mode (DECAWM)
* ^[[?12l
  - Reset Mode: Send/receive (SRM)
* ^[[?25h
  - Show Cursor (DECTCEM)
* ^[[?1h
  - Application Cursor Keys (DECCKM)
* ^[=
  - Application Keypad (DECKPAM)
* ^[[?1h
  - Application Cursor Keys (DECCKM)
* ^[=
* ^[[?1h
  - Application Cursor Keys (DECCKM)
* ^[=
  - Application Keypad (DECKPAM)
* ^[[39;49m
  - Set foreground color to default (original).
  - Set background color to default (original).
* ^[[39;49m
  - Set foreground color to default (original).
  - Set background color to default (original).
* ^[(B
  - setusg0
  - Set United States G0 character set
* ^[[m
  - SGR0
  - Turn off character attributes
* ^[[H
  - Move cursor to upper left corner  cursorhome
* ^[[2J
  - Erase display
* ^[(B
  - setusg0
  - Set United States G0 character set
* ^[[0;7m
  - Turn off character attributes
  - Inverse
* "  GNU nano 2.0.6                                                         New Buffer                                                                                                                         "
* ^[[53;1H
  - Cursor Position [row;column]
* "^G"
* ^[(B
  - setusg0
  - Set United States G0 character set
* ^[[m
  - SGR0
  - Turn off character attributes
* " Get Help"
* ^[[53;35H
  - Cursor Position [row;column]
* ^[(B
* ^[[0;7m
* "^O"
* ^[(B
* ^[[m
* " WriteOut"
* ^[[53;69H
* ^[(B
* ^[[0;7m^R
* ^[(B
* ^[[m Read File
* ^[[53;103H
* ^[(B
* ^[[0;7m^Y
* ^[(B
* ^[[m
* " Prev Page"
* ^[[53;137H
* ^[(B
* ^[[0;7m^K
* ^[(B
* ^[[m Cut Text
* ^[[53;171H
* ^[(B
* ^[[0;7m^C
* ^[(B
* ^[[m
* " Cur Pos"^M
* ^[[54d
* ^[(B
* ^[[0;7m^X
* ^[(B
* ^[[m
* " Exit"
* ^[[54;35H
* ^[(B
* ^[[0;7m^J
* ^[(B
* ^[[m
* " Justify"
* ^[[54;69H
* ^[(B
* ^[[0;7m^W
* ^[(B
* ^[[m
* " Where Is"
* ^[[54;103H
* ^[(B
* ^[[0;7m^V
* ^[(B
* ^[[m Next Page
* ^[[54;137H
* ^[(B
* ^[[0;7m^U
* ^[(B
* ^[[m UnCut Text
* ^[[54;171H
* ^[(B
* ^[[0;7m^T
* ^[(B
* ^[[m To Spell
* ^[[3d
* ^[[53d
* ^[[J
* ^[[54;204H
* ^[[54;1H
* ^[[?1049l
* ^[[?1l
* ^[>
