# HW

Pure Go Library to control on a low level board access.
This project is purely for research to play with my [Asus Tinker Board](https://www.asus.com/us/Single-Board-Computer/Tinker-Board/) GPIO and my [64x64 RGB LED Matrix](https://www.adafruit.com/product/3649) in Go.

It is inspired by:

* [`go-rpi-rgb-led-matrix`](https://github.com/mcuadros/go-rpi-rgb-led-matrix): The Go binding for [`rpi-rgb-led-matrix`](https://github.com/hzeller/rpi-rgb-led-matrix) an excellent C++ library to control [RGB LED displays](https://learn.adafruit.com/32x16-32x32-rgb-led-matrix/overview) with Raspberry Pi GPIO.
* [`gpio_lib_c`](https://github.com/TinkerBoard/gpio_lib_c): GPIO_LIB is a extension of WiringPi, it can control low speed peripherial of Tinker Board.

This library includes the basic bindings to control the LED Matrix directly and also a convenient Matrix Toolkit with more high level functions. Also some [examples](https://github.com/post-l/hw/tree/master/examples) are included to test the library and the configuration.

To learn about the configuration and the wiring go to the [rpi-rgb-led-matrix](https://github.com/hzeller/rpi-rgb-led-matrix) project. It is highly detailed and well explained.

![Life Gif](life.gif)

Gopher gif from [`egonelbre/gophers`](https://github.com/egonelbre/gophers).

## Matrix Emulation

As part of the library, a small Matrix emulator is provided. The emulator renders a virtual RGB matrix on a window in your desktop, without the need of a real RGB matrix connected to your computer.

To start the examples with the emulator, set the `-emulator` flag.

## License

MIT, see [LICENSE](LICENSE)
