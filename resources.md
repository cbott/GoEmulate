> Links
Go GB emulator but out of date and won't build. Has a number of helpful links
https://github.com/Humpheh/goboy

OpenGL-based graphics library written in Go
https://github.com/faiface/pixel

Java GB emulator
https://github.com/trekawek/coffee-gb

Rust GB emulator
https://github.com/Gekkio/mooneye-gb/tree/master


GB docs:
- http://bgb.bircd.org/pandocs.htm#cpuinstructionset
- http://marc.rawer.de/Gameboy/Docs/GBCPUman.pdf
- https://pastraiser.com/cpu/gameboy/gameboy_opcodes.html



> Setup
Per pixel requirements https://github.com/faiface/pixel#requirements
- sudo apt install libgl1-mesa-dev
- sudo apt install xorg-dev

go get .
go build
./GoEmulate
