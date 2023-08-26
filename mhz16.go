package mhz16

import (
	"fmt"

	"github.com/tarm/serial"
)

type Mhz16 struct {
	port serial.Port
}

var (
	getCommand         = []byte{0xff, 0x01, 0x86, 0x00, 0x00, 0x00, 0x00, 0x00, 0x79}
	calibrationCommand = []byte{0xff, 0x01, 0x87, 0x00, 0x00, 0x00, 0x00, 0x00, 0x78}
)

func Open(portName string) (*Mhz16, error) {
	options := &serial.Config{
		Name:        portName,
		Baud:        9600,
		ReadTimeout: 0,
	}
	port, err := serial.OpenPort(options)
	if err != nil {
		return nil, fmt.Errorf("serial port cannot open: %w", err)
	}
	m := Mhz16{port: *port}
	return &m, nil
}

func (m *Mhz16) Close() error {
	err := m.port.Close()
	return err
}

func (m *Mhz16) GetCo2() (int, error) {
	n, err := m.port.Write(getCommand)
	if err != nil {
		return 0, fmt.Errorf("serial write error: %w", err)
	}
	res := make([]byte, 9)
	readed := 0
	for {
		n, err = m.port.Read(res[readed:9])
		if err != nil {
			return 0, fmt.Errorf("serial read error: %w", err)
		}
		if n == 0 {
			break
		}
		readed += n
		if readed >= 9 {
			break
		}
	}
	checksum := 0xff & (^(res[1] + res[2] + res[3] + res[4] + res[5] + res[6] + res[7]) + 1)
	if res[8] != checksum {
		return 0, fmt.Errorf("checksum not match error")
	}
	result := int(res[2])<<8 + int(res[3])
	return result, nil
}

func (m *Mhz16) ZeroCalibration() error {
	n, err := m.port.Write(calibrationCommand)
	if err != nil || n != 9 {
		return fmt.Errorf("serial write error: %w", err)
	}
	return nil
}
