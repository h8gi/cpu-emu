package main

import (
	"fmt"
)

type OpCode = uint16

const (
	Mov OpCode = iota
	Add
	Sub
	And
	Or
	Sl
	Sr
	Sra
	Ldl
	Ldh
	Cmp
	Je
	Jmp
	Ld
	St
	Hlt
)

type RegIndex = uint16

const (
	Reg0 RegIndex = iota
	Reg1
	Reg2
	Reg3
	Reg4
	Reg5
	Reg6
	Reg7
)

var tokenMap = map[string]uint16{
	"mov":  Mov,
	"add":  Add,
	"sub":  Sub,
	"and":  And,
	"or":   Or,
	"sl":   Sl,
	"sr":   Sr,
	"sra":  Sra,
	"ldl":  Ldl,
	"ldh":  Ldh,
	"cmp":  Cmp,
	"je":   Je,
	"jmp":  Jmp,
	"ld":   Ld,
	"st":   St,
	"hlt":  Hlt,
	"reg0": Reg0,
	"reg1": Reg1,
	"reg2": Reg2,
	"reg3": Reg3,
	"reg4": Reg4,
	"reg5": Reg5,
	"reg6": Reg6,
	"reg7": Reg7,
}

type Assembler struct {
}

// func (asm *Assembler) Assemble(program string) []uint16 {
// 	lines := strings.Split(program, "\n")
// 	for _, l := range lines {
// 		l = strings.TrimSpace(l)
// 		tokens := strings.Split(l, " ")

// 	}
// }

func (asm *Assembler) Mov(ra, rb RegIndex) uint16 {
	return Mov<<11 | ra<<8 | rb<<5
}

func (asm *Assembler) Add(ra, rb RegIndex) uint16 {
	return Add<<11 | ra<<8 | rb<<5
}

func (asm *Assembler) Sub(ra, rb RegIndex) uint16 {
	return Sub<<11 | ra<<8 | rb<<5
}

func (asm *Assembler) And(ra, rb RegIndex) uint16 {
	return And<<11 | ra<<8 | rb<<5
}

func (asm *Assembler) Or(ra, rb RegIndex) uint16 {
	return Or<<11 | ra<<8 | rb<<5
}

func (asm *Assembler) Sl(ra RegIndex) uint16 {
	return Sl<<11 | ra<<8
}

func (asm *Assembler) Sr(ra RegIndex) uint16 {
	return Sr<<11 | ra<<8
}

func (asm *Assembler) Sra(ra RegIndex) uint16 {
	return Sra<<11 | ra<<8
}

func (asm *Assembler) Ldl(ra RegIndex, ival uint16) uint16 {
	return Ldl<<11 | ra<<8 | ival&0x00ff
}

func (asm *Assembler) Ldh(ra RegIndex, ival uint16) uint16 {
	return Ldh<<11 | ra<<8 | ival&0x00ff
}

func (asm *Assembler) Cmp(ra, rb RegIndex) uint16 {
	return Cmp<<11 | ra<<8 | rb<<5
}

func (asm *Assembler) Je(addr uint16) uint16 {
	return Je<<11 | addr&0x00ff
}

func (asm *Assembler) Jmp(addr uint16) uint16 {
	return Jmp<<11 | addr&0x00ff
}

func (asm *Assembler) Ld(ra RegIndex, addr uint16) uint16 {
	return Ld<<11 | ra<<8 | addr&0x00ff
}

func (asm *Assembler) St(ra RegIndex, addr uint16) uint16 {
	return St<<11 | ra<<8 | addr&0x00ff
}

func (asm *Assembler) Hlt() uint16 {
	return Hlt << 11
}

type CpuEmulator struct {
	reg     [8]uint16
	rom     [256]uint16
	ram     [256]uint16
	pc      uint16
	ir      uint16
	flag_eq uint16
}

func (cpu *CpuEmulator) String() string {
	return fmt.Sprintf("%5d %5x %5d %5d %5d %5d",
		cpu.pc, cpu.ir,
		cpu.reg[0], cpu.reg[1], cpu.reg[2], cpu.reg[3])
}

func (cpu *CpuEmulator) GetOpCode() OpCode {
	return cpu.ir >> 11
}

func (cpu *CpuEmulator) GetOpRegA() RegIndex {
	return (cpu.ir >> 8) & 0x0007
}

func (cpu *CpuEmulator) GetOpRegB() RegIndex {
	return (cpu.ir >> 5) & 0x0007
}

func (cpu *CpuEmulator) GetOpData() uint16 {
	return cpu.ir & 0x00ff
}

func (cpu *CpuEmulator) GetOpAddr() uint16 {
	return cpu.ir & 0x00ff
}

func (cpu *CpuEmulator) Run() error {
	fmt.Println("   pc    ir  reg0  reg1  reg2  reg3")

	for cpu.GetOpCode() != Hlt {
		cpu.ir = cpu.rom[cpu.pc]
		fmt.Println(cpu)
		cpu.pc += 1
		switch cpu.GetOpCode() {
		case Mov:
			cpu.reg[cpu.GetOpRegA()] = cpu.reg[cpu.GetOpRegB()]
		case Add:
			cpu.reg[cpu.GetOpRegA()] = cpu.reg[cpu.GetOpRegA()] + cpu.reg[cpu.GetOpRegB()]
		case Sub:
			cpu.reg[cpu.GetOpRegA()] = uint16(int16(cpu.reg[cpu.GetOpRegA()]) -
				int16(cpu.reg[cpu.GetOpRegB()]))
		case And:
			cpu.reg[cpu.GetOpRegA()] = cpu.reg[cpu.GetOpRegA()] & cpu.reg[cpu.GetOpRegB()]
		case Or:
			cpu.reg[cpu.GetOpRegA()] = cpu.reg[cpu.GetOpRegA()] | cpu.reg[cpu.GetOpRegB()]
		case Sl:
			cpu.reg[cpu.GetOpRegA()] = cpu.reg[cpu.GetOpRegA()] << 1
		case Sr:
			cpu.reg[cpu.GetOpRegA()] = cpu.reg[cpu.GetOpRegA()] >> 1
		case Sra:
			cpu.reg[cpu.GetOpRegA()] = (cpu.reg[cpu.GetOpRegA()] & 0x8000) |
				(cpu.reg[cpu.GetOpRegA()] >> 1)
		case Ldl:
			cpu.reg[cpu.GetOpRegA()] = (cpu.reg[cpu.GetOpRegA()] & 0xff00) |
				(cpu.GetOpData() & 0x00ff)
		case Ldh:
			cpu.reg[cpu.GetOpRegA()] = (cpu.GetOpData() << 8) |
				(cpu.reg[cpu.GetOpRegA()] & 0x00ff)
		case Cmp:
			if cpu.reg[cpu.GetOpRegA()] == cpu.reg[cpu.GetOpRegB()] {
				cpu.flag_eq = 1
			} else {
				cpu.flag_eq = 0
			}
		case Je:
			if cpu.flag_eq == 1 {
				cpu.pc = cpu.GetOpAddr()
			}
		case Jmp:
			cpu.pc = cpu.GetOpAddr()
		case Ld:
			cpu.reg[cpu.GetOpRegA()] = cpu.ram[cpu.GetOpAddr()]
		case St:
			cpu.ram[cpu.GetOpAddr()] = cpu.reg[cpu.GetOpRegA()]
		}
	}

	fmt.Printf("ram[64] = %d\n", cpu.ram[64])

	return nil
}

func (cpu *CpuEmulator) Setup() {
	asm := &Assembler{}
	cpu.rom[0] = asm.Ldh(Reg0, 0)
	cpu.rom[1] = asm.Ldl(Reg0, 0)
	cpu.rom[2] = asm.Ldh(Reg1, 0)
	cpu.rom[3] = asm.Ldl(Reg1, 1)
	cpu.rom[4] = asm.Ldh(Reg2, 0)
	cpu.rom[5] = asm.Ldl(Reg2, 0)
	cpu.rom[6] = asm.Ldh(Reg3, 0)
	cpu.rom[7] = asm.Ldl(Reg3, 10)
	cpu.rom[8] = asm.Add(Reg2, Reg1)
	cpu.rom[9] = asm.Add(Reg0, Reg2)
	cpu.rom[10] = asm.St(Reg0, 64)
	cpu.rom[11] = asm.Cmp(Reg2, Reg3)
	cpu.rom[12] = asm.Je(14)
	cpu.rom[13] = asm.Jmp(8)
	cpu.rom[14] = asm.Hlt()
}

func main() {
	cpu := &CpuEmulator{}
	cpu.Setup()
	cpu.Run()
}
