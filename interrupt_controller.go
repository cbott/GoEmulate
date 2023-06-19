package main

// Code for interrupt controller

// Interrupt memory addresses
const (
	IE                 = 0xFFFF // Interrupt Enable
	IF                 = 0xFF0F // Interrupt Flag
	Interrupt_vblank   = 1 << 0
	Interrupt_lcd_stat = 1 << 1
	Interrupt_timer    = 1 << 2
	Interrupt_serial   = 1 << 3
	Interrupt_joypad   = 1 << 4
)

// These are the addresses the CPU will jump to when executing the interrupt
const (
	InterruptVBlankAddress  = 0x40
	InterruptLCDStatAddress = 0x48
	InterruptTimerAddress   = 0x50
	InterruptSerialAddress  = 0x58
	InterruptJoypadAddress  = 0x60
)

func (gb *Gameboy) SetInterruptRequestFlag(flag uint8) {
	// Sets a bit in the interrupt flag register, which will trigger an interrupt
	// if the Interrupt Enable register and Interrupt Master Enable flag allow it
	gb.memory.set(IF, gb.memory.get(IF)|flag)
}
