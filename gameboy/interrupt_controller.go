package gameboy

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
var interruptAddresses = [5]uint16{
	0x0040, // VBlank Address
	0x0048, // LCDStat Address
	0x0050, // Timer Address
	0x0058, // Serial Address
	0x0060, // Joypad Address
}

// Sets a bit in the interrupt flag register, which will trigger an interrupt
// if the Interrupt Enable register and Interrupt Master Enable flag allow it
func (gb *Gameboy) SetInterruptRequestFlag(flag uint8) {
	gb.memory.set(IF, gb.memory.get(IF)|flag)
}

// Run interrupts for the processor and return number of cycles it took (4 MHz clock cycles)
func (gb *Gameboy) RunInterrupts() int {
	request := gb.memory.get(IF)
	enabled := gb.memory.get(IE)

	// Determine if CPU should resume from a HALT. This can happen even if MasterEnable is false
	if gb.halted {
		if (request&enabled)&0b00011111 != 0 {
			gb.halted = false
		}
	}

	if gb.pendingInterruptEnable {
		// If delayed enable request was made, enable but do not process interrupts until next time
		gb.pendingInterruptEnable = false
		if !gb.interruptMasterEnable {
			gb.interruptMasterEnable = true
			return 0
		}
	}

	if !gb.interruptMasterEnable {
		return 0
	}

	interruptPerformed := false

	// Service interrupts if requested
	for i := 0; i < 5; i++ {
		if (request&(1<<i)) != 0 && (enabled&(1<<i)) != 0 {
			// interrupt `i` is enabled
			// reset this bit in the interrupt flag register
			gb.memory.set(IF, request & ^(1<<i))
			// disable interrupts globally
			gb.interruptMasterEnable = false

			// push current PC to the stack and jump to the interrupt address
			gb.pushToStack16(gb.cpu.getRegister16(regPC))
			gb.cpu.setRegister16(regPC, interruptAddresses[i])

			// CPU will service only 1 interrupt at a time, highest priority first
			interruptPerformed = true
			break
		}
	}

	if interruptPerformed {
		return 20
	}
	return 0
}
