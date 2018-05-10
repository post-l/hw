package tinkerboard

const (
	gpioLen      int64 = 0x00010000
	gpioCh       int64 = 0x00020000
	gpioBaseAddr int64 = 0xff750000
	gpioBankLen        = 9
)

const (
	blockSize = 4096

	RK3288_GRF_PHYS int64 = 0xff770000
	RK3288_PWM      int64 = 0xff680000
	RK3288_PMU      int64 = 0xff730000
	RK3288_CRU      int64 = 0xff760000
)

const (
	CRU_CLKGATE14_CON = 0x0198
	CRU_CLKGATE17_CON = 0x01a4
)

const (
	GPIO0_C1 = 17 //7----->17

	GPIO1_D0 = (24 + 24)

	GPIO5_B0 = (8 + 152)  //8----->160
	GPIO5_B1 = (9 + 152)  //9----->161
	GPIO5_B2 = (10 + 152) //10----->162
	GPIO5_B3 = (11 + 152) //11----->163
	GPIO5_B4 = (12 + 152) //12----->164
	GPIO5_B5 = (13 + 152) //13----->165
	GPIO5_B6 = (14 + 152) //14----->166
	GPIO5_B7 = (15 + 152) //15----->167
	GPIO5_C0 = (16 + 152) //16----->168
	GPIO5_C1 = (17 + 152) //17----->169
	GPIO5_C2 = (18 + 152) //18----->170
	GPIO5_C3 = (19 + 152) //19----->171

	GPIO6_A0 = (184)     //0----->184
	GPIO6_A1 = (1 + 184) //1----->185
	GPIO6_A2 = (2 + 184) //2----->186
	GPIO6_A3 = (3 + 184) //3----->187
	GPIO6_A4 = (4 + 184) //4----->188

	GPIO7_A0 = (0 + 216)  //0----->216
	GPIO7_A7 = (7 + 216)  //7----->223
	GPIO7_B0 = (8 + 216)  //8----->224
	GPIO7_B1 = (9 + 216)  //9----->225
	GPIO7_B2 = (10 + 216) //10----->226
	GPIO7_C1 = (17 + 216) //17----->233
	GPIO7_C2 = (18 + 216) //18----->234
	GPIO7_C6 = (22 + 216) //22----->238
	GPIO7_C7 = (23 + 216) //23----->239

	GPIO8_A3 = (3 + 248) //3----->251
	GPIO8_A4 = (4 + 248) //4----->252
	GPIO8_A5 = (5 + 248) //5----->253
	GPIO8_A6 = (6 + 248) //6----->254
	GPIO8_A7 = (7 + 248) //7----->255
	GPIO8_B0 = (8 + 248) //8----->256
	GPIO8_B1 = (9 + 248) //9----->257
)

const (
	GPIO_SWPORTA_DR_OFFSET    = 0x0000
	GPIO_SWPORTA_DDR_OFFSET   = 0x0004
	GPIO_INTEN_OFFSET         = 0x0030
	GPIO_INTMASK_OFFSET       = 0x0034
	GPIO_INTTYPE_LEVEL_OFFSET = 0x0038
	GPIO_INT_POLARITY_OFFSET  = 0x003c
	GPIO_INT_STATUS_OFFSET    = 0x0040
	GPIO_INT_RAWSTATUS_OFFSET = 0x0044
	GPIO_DEBOUNCE_OFFSET      = 0x0048
	GPIO_PORTA_EOF_OFFSET     = 0x004c
	GPIO_EXT_PORTA_OFFSET     = 0x0050
	GPIO_LS_SYNC_OFFSET       = 0x0060
)

const (
	PMU_GPIO0C_IOMUX = 0x008c
)

const (
	GRF_GPIO5B_IOMUX  = 0x0050
	GRF_GPIO5C_IOMUX  = 0x0054
	GRF_GPIO6A_IOMUX  = 0x005c
	GRF_GPIO6B_IOMUX  = 0x0060
	GRF_GPIO6C_IOMUX  = 0x0064
	GRF_GPIO7A_IOMUX  = 0x006c
	GRF_GPIO7B_IOMUX  = 0x0070
	GRF_GPIO7CL_IOMUX = 0x0074
	GRF_GPIO7CH_IOMUX = 0x0078
	GRF_GPIO8A_IOMUX  = 0x0080
	GRF_GPIO8B_IOMUX  = 0x0084
)

const (
	GRF_GPIO6A_P = 0x0190
)
