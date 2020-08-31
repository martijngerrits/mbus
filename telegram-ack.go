package mbus

import "fmt"

type TelegramACK WMBusFrame

func NewTelegramACK() *TelegramACK {
	frame := &TelegramACK{
		Start: FRAME_LONG_START,
		Stop: FRAME_STOP,
	}

	return frame
}

func (f *TelegramACK) CalculateChecksum() byte {
	checksum := f.Control
	for i := 0; i < len(f.Header.Id); i++ {
        checksum += f.Header.Id[i]
    }
	checksum += f.ControlInformation

	return checksum
	//f.Checksum = f.Control
	//f.Checksum += f.Address
	//f.Checksum += f.ControlInformation
}

func (f *TelegramACK) CalculateLength() int {
	return 0
}

func (f *TelegramACK) Verify() error {
	if f.Start != FRAME_ACK_START {
		return fmt.Errorf("no frame start")
	}

	return nil
}

func (f *TelegramACK) Encode() ([]byte, int) {
	return []byte{}, 0
}
