package microwallet

import "github.com/jacobsa/go-serial/serial"

// DefaultConfig use linux default serial port /dev/tty.SLAB_USBtoUART
func DefaultConfig() (*serial.OpenOptions, error) {
	return &serial.OpenOptions{
		PortName:        "/dev/tty.SLAB_USBtoUART",
		BaudRate:        115200,
		DataBits:        8,
		StopBits:        1,
		MinimumReadSize: 1,
	}, nil
}
