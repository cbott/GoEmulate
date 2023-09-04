# GoEmulate
Game Boy Emulator in GO

This is a work in progress!

GoEmulate is an emulator for the original Game Boy (DMG) written entirely in Go, utilizing the [pixel](https://github.com/faiface/pixel) library for graphics and [Oto](https://github.com/ebitengine/oto) for sound. The goal of the project is to reach basic emulator functionality with minimal code complexity, and very little focus on UI or useability.

Currently no support is planned for GBC/GBA emulation

GoEmulate is a learning project, and the first emulator I have written. The code was heavily influenced by all the projects listed below as development resources and I enocourage you to use them directly.


To Do List
----------
- reorganize into packages
- fix display bug with window over background
- CPU save states
- allow loading files from cmd line argument
- more unit tests
- better error handling
- allow loading files from GUI or similar
- other cartridge types
- implement serial

Completed Tasks
- RAM saving
- longer compare with goboy
- access registers by enum
- window scaling
- MBC3 cartridge type
- code cleanup/resolve todos
- implement sound

Setup
-----
Per pixel requirements https://github.com/faiface/pixel#requirements
```
sudo apt install libgl1-mesa-dev
sudo apt install xorg-dev
```

Development Resources
---
Other emulators used for comparison
- (go) https://github.com/Humpheh/goboy
- (java) https://github.com/trekawek/coffee-gb
- (rust) https://github.com/Gekkio/mooneye-gb
- (c) https://github.com/AntonioND/giibiiadvance

Game Boy docs:
- https://izik1.github.io/gbops/index.html
- http://bgb.bircd.org/pandocs.htm#cpuinstructionset
- http://marc.rawer.de/Gameboy/Docs/GBCPUman.pdf
- https://gbdev.io/pandocs/Specifications.html
