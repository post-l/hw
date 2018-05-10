package tinkerboard

import (
	"fmt"
	"os"
	"reflect"
	"syscall"
	"unsafe"

	"github.com/post-l/hw/board"
)

type TinkerBoard struct {
	gpio [][]uint32
	grf  []uint32
	pwm  []uint32
	pmu  []uint32
	cru  []uint32

	gpioMap [][]byte
	grfMap  []byte
	pwmMap  []byte
	pmuMap  []byte
	cruMap  []byte
}

func New() (*TinkerBoard, error) {
	mem, err := os.OpenFile("/dev/mem", os.O_RDWR|os.O_SYNC, 0755)
	if err != nil {
		return nil, err
	}
	defer mem.Close()
	tb := &TinkerBoard{
		gpioMap: make([][]byte, gpioBankLen),
		gpio:    make([][]uint32, gpioBankLen),
	}
	memFd := int(mem.Fd())
	for i := range tb.gpio {
		offset := gpioBaseAddr + int64(i)*gpioLen
		if i > 0 {
			offset += gpioCh
		}
		tb.gpioMap[i], err = syscall.Mmap(memFd, offset, blockSize, syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
		if err != nil {
			tb.Close()
			return nil, fmt.Errorf("unable to map gpio bank %d: %v", i, err)
		}
		tb.gpio[i] = mapToUInt32Slice(tb.gpioMap[i])
	}

	tb.grfMap, err = syscall.Mmap(memFd, RK3288_GRF_PHYS, blockSize, syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	if err != nil {
		tb.Close()
		return nil, fmt.Errorf("unable to map grf: %v", err)
	}
	tb.grf = mapToUInt32Slice(tb.grfMap)

	tb.pwmMap, err = syscall.Mmap(memFd, RK3288_PWM, blockSize, syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	if err != nil {
		tb.Close()
		return nil, fmt.Errorf("unable to map pwm: %v", err)
	}
	tb.pwm = mapToUInt32Slice(tb.pwmMap)

	tb.pmuMap, err = syscall.Mmap(memFd, RK3288_PMU, blockSize, syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	if err != nil {
		tb.Close()
		return nil, fmt.Errorf("unable to map pmu: %v", err)
	}
	tb.pmu = mapToUInt32Slice(tb.pmuMap)

	tb.cruMap, err = syscall.Mmap(memFd, RK3288_CRU, blockSize, syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	if err != nil {
		tb.Close()
		return nil, fmt.Errorf("unable to map cru: %v", err)
	}
	tb.cru = mapToUInt32Slice(tb.cruMap)

	return tb, nil
}

func (tb *TinkerBoard) SetPinMode(pin int, mode board.PinMode) {
	bank, bankPin := gpioToBank(pin)
	cc := tb.gpioClkDisable(bank)

	tb.setGPIOPinMode(pin)
	switch mode {
	case board.Input:
		tb.gpio[bank][GPIO_SWPORTA_DDR_OFFSET/4] &= ^(1 << bankPin)
	case board.Output:
		tb.gpio[bank][GPIO_SWPORTA_DDR_OFFSET/4] |= (1 << bankPin)
	}

	tb.gpioClkEnable(cc)
}

func (tb *TinkerBoard) DigitalRead(pin int) bool {
	bank, bankPin := gpioToBank(pin)
	cc := tb.gpioClkDisable(bank)

	r := tb.gpio[bank][GPIO_EXT_PORTA_OFFSET/4]
	v := ((r & (1 << bankPin)) >> bankPin) != 0

	tb.gpioClkEnable(cc)
	return v
}

func (tb *TinkerBoard) DigitalWrite(pin int, v bool) {
	bank, bankPin := gpioToBank(pin)
	cc := tb.gpioClkDisable(bank)

	if v {
		tb.gpio[bank][GPIO_SWPORTA_DR_OFFSET/4] |= (1 << bankPin)
	} else {
		tb.gpio[bank][GPIO_SWPORTA_DR_OFFSET/4] &= ^(1 << bankPin)
	}

	tb.gpioClkEnable(cc)
}

func (tb *TinkerBoard) DigitalWrites(pvs []board.PinValue) {
	banks := make([]uint32, gpioBankLen)
	maskBanks := make([]uint32, gpioBankLen)
	cruv := uint32(0)
	for _, pv := range pvs {
		pin := pv.Pin
		v := pv.Value
		bank, bankPin := gpioToBank(pin)
		bitPin := uint32(1 << bankPin)
		maskBanks[bank] |= bitPin
		if v {
			banks[bank] |= bitPin
		}
		if bank != 0 {
			cruv |= 1 << (16 + uint32(bank))
		}
	}
	if mask := maskBanks[0]; mask != 0 {
		writeBit := uint32(4)
		regOffset := CRU_CLKGATE17_CON / 4
		v := uint32(1 << (16 + writeBit))
		tb.cru[regOffset] = v
		tb.writeMaskedBits(0, banks[0], mask)
		tb.cru[regOffset] = v
	}
	if cruv != 0 {
		tb.cru[CRU_CLKGATE14_CON/4] = cruv
		for bank := 1; bank < gpioBankLen; bank++ {
			mask := maskBanks[bank]
			if mask == 0 {
				continue
			}
			tb.writeMaskedBits(uint32(bank), banks[bank], mask)
		}
		tb.cru[CRU_CLKGATE14_CON/4] = cruv
	}
}

func (tb *TinkerBoard) writeMaskedBits(bank, value, mask uint32) {
	tb.gpio[bank][GPIO_SWPORTA_DR_OFFSET/4] = tb.gpio[bank][GPIO_SWPORTA_DR_OFFSET/4] & ^(^value&mask) | value
}

type clkCtrl struct {
	regOffset int
	v         uint32
}

func (tb *TinkerBoard) gpioClkDisable(bank uint32) clkCtrl {
	var cc clkCtrl
	writeBit := bank
	regOffset := CRU_CLKGATE14_CON / 4
	if bank == 0 {
		writeBit = 4
		cc.regOffset = CRU_CLKGATE17_CON / 4
	}
	v := uint32(1 << (16 + writeBit))
	tb.cru[regOffset] = v
	return clkCtrl{
		regOffset: regOffset,
		v:         v,
	}
}

func (tb *TinkerBoard) gpioClkEnable(ctrl clkCtrl) {
	tb.cru[ctrl.regOffset] = ctrl.v
}

func (tb *TinkerBoard) Close() error {
	for _, data := range tb.gpioMap {
		if data == nil {
			return nil
		}
		syscall.Munmap(data)
	}
	d := [][]byte{tb.grfMap, tb.pwmMap, tb.pmuMap, tb.cruMap}
	for _, data := range d {
		if data == nil {
			return nil
		}
		syscall.Munmap(data)
	}
	return nil
}

func gpioToBank(gpio int) (uint32, uint32) {
	if gpio < 24 {
		return 0, uint32(gpio)
	}
	return uint32(((gpio - 24) / 32) + 1), uint32((gpio - 24) % 32)
}

func mapToUInt32Slice(m []byte) []uint32 {
	h := (*reflect.SliceHeader)(unsafe.Pointer(&m))
	h.Len /= 4
	h.Cap /= 4
	return *(*[]uint32)(unsafe.Pointer(h))
}

func (tb *TinkerBoard) setGPIOPinMode(pin int) {
	p := uint32(pin)
	switch p {
	//GPIO0
	case GPIO0_C1:
		tb.pmu[PMU_GPIO0C_IOMUX/4] = (tb.pmu[PMU_GPIO0C_IOMUX/4] | (0x03 << ((p%8)*2 + 16))) & (^(0x03 << ((p % 8) * 2)))
	case GPIO1_D0:
	//GPIO5B
	case GPIO5_B0, GPIO5_B1, GPIO5_B2, GPIO5_B3, GPIO5_B4, GPIO5_B5, GPIO5_B6, GPIO5_B7:
		tb.grf[GRF_GPIO5B_IOMUX/4] = (tb.grf[GRF_GPIO5B_IOMUX/4] | (0x03 << ((p%8)*2 + 16))) & (^(0x03 << ((p % 8) * 2)))
	//GPIO5C
	case GPIO5_C0, GPIO5_C1, GPIO5_C2, GPIO5_C3:
		tb.grf[GRF_GPIO5C_IOMUX/4] = (tb.grf[GRF_GPIO5C_IOMUX/4] | (0x03 << ((p%8)*2 + 16))) & (^(0x03 << ((p % 8) * 2)))
	//GPIO6A
	case GPIO6_A1:
		tb.grf[GRF_GPIO6A_IOMUX/4] = (tb.grf[GRF_GPIO6A_IOMUX/4] | (0x0f << ((p%8)*2 + 16))) & (^(0x0f << ((p % 8) * 2)))
		// tb.grf[GRF_GPIO6A_P/4] = ((tb.grf[GRF_GPIO6A_P/4] | (0x03 << (((GPIO6_A2)%8)*2 + 16))) & (^(0x03 << (((GPIO6_A2) % 8) * 2)))) | (0 << (((GPIO6_A2)%8)*2 + 1)) | (0 << (((GPIO6_A2) % 8) * 2))
	case GPIO6_A0, GPIO6_A3, GPIO6_A4:
		tb.grf[GRF_GPIO6A_IOMUX/4] = (tb.grf[GRF_GPIO6A_IOMUX/4] | (0x03 << ((p%8)*2 + 16))) & (^(0x03 << ((p % 8) * 2)))
	//GPIO7A7
	case GPIO7_A7:
		tb.grf[GRF_GPIO7A_IOMUX/4] = (tb.grf[GRF_GPIO7A_IOMUX/4] | (0x03 << ((p%8)*2 + 16))) & (^(0x03 << ((p % 8) * 2)))
	//GPIO7B
	case GPIO7_B0, GPIO7_B1, GPIO7_B2:
		tb.grf[GRF_GPIO7B_IOMUX/4] = (tb.grf[GRF_GPIO7B_IOMUX/4] | (0x03 << ((p%8)*2 + 16))) & (^(0x03 << ((p % 8) * 2)))
	//GPIO7C
	case GPIO7_C1, GPIO7_C2:
		tb.grf[GRF_GPIO7CL_IOMUX/4] = (tb.grf[GRF_GPIO7CL_IOMUX/4] | (0x0f << (16 + (p%8)*4))) & (^(0x0f << ((p % 8) * 4)))
	case GPIO7_C6, GPIO7_C7:
		tb.grf[GRF_GPIO7CH_IOMUX/4] = (tb.grf[GRF_GPIO7CH_IOMUX/4] | (0x0f << (16 + (p%8-4)*4))) & (^(0x0f << ((p%8 - 4) * 4)))
	//GPIO8A
	case GPIO8_A3, GPIO8_A6, GPIO8_A7, GPIO8_A4, GPIO8_A5:
		tb.grf[GRF_GPIO8A_IOMUX/4] = (tb.grf[GRF_GPIO8A_IOMUX/4] | (0x03 << ((p%8)*2 + 16))) & (^(0x03 << ((p % 8) * 2)))
	//GPIO8B
	case GPIO8_B0, GPIO8_B1:
		tb.grf[GRF_GPIO8B_IOMUX/4] = (tb.grf[GRF_GPIO8B_IOMUX/4] | (0x03 << ((p%8)*2 + 16))) & (^(0x03 << ((p % 8) * 2)))
	}
}
