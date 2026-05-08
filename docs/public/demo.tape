# obacht README demo recording.
#
# Regenerate:
#   make build && vhs demo.tape
#
# Requires: mise install (vhs 0.11.0), ttyd, ffmpeg.
# The demo runs against a curated fixture rule set in demo/ so the
# output is deterministic and never reflects the recorder's real machine.

Output docs/public/demo.gif

Set FontSize 14
Set Width 1400
Set Height 800
Set Theme "Dracula"
Set TypingSpeed 60ms

Hide
Type "export PATH=$PWD/bin:$PATH"
Enter
Type "clear"
Enter
Sleep 300ms
Show

Type "obacht scan --rules-dir demo"
Sleep 600ms
Enter
Sleep 6s
