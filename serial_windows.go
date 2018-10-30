package microwallet

import (
	"fmt"

	"github.com/jacobsa/go-serial/serial"
	xserial "go.bug.st/serial.v1"
)

// DefaultConfig search serial port COM1 - COM256
func DefaultConfig() (*serial.OpenOptions, error) {
	ports, err := xserial.GetPortsList()
	if err != nil {
		return nil, err
	}
	if len(ports) == 0 {
		return nil, fmt.Errorf("No COM port detected")
	}
	return &serial.OpenOptions{
		PortName:        ports[0],
		BaudRate:        115200,
		DataBits:        8,
		StopBits:        1,
		MinimumReadSize: 1,
	}, nil
}
