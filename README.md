# GoEmulate
Game Boy Emulator in GO

This is a work in progress!

GoEmulate is an emulator for the original Game Boy (DMG) written entirely in Go, utilizing the [pixel](https://github.com/faiface/pixel) library for graphics. The goal of the project is to reach basic emulator functionality with minimal code complexity, and very little focus on UI or useability.

Currently no support is planned for GBC/GBA emulation

GoEmulate is a learning project, and the first emulator I have ever even attempted. While I did write all the code it was heavily influenced by https://github.com/Humpheh/goboy


To Do List
- other cartridge types
- code cleanup/resolve todos
- implement sound
- allow loading files from cmd line argument
- reorganize into packages
- more unit tests
- better error handling
- allow loading files from GUI or similar

Completed Tasks
- longer compare with goboy
- access registers by enum
- window scaling
