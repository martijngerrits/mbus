package mbus

import "fmt"

type TelegramShort WMBusFrame

func NewTelegramShort() *TelegramShort {
	frame := &TelegramShort{
		Start: FRAME_SHORT_START,
		Stop: FRAME_STOP,
	}

	return frame
}

func (f *TelegramShort) CalculateChecksum() byte {
	checksum := f.Control
	for i := 0; i < len(f.Header.Id); i++ {
        checksum += f.Header.Id[i]
    }

	return checksum
	//f.Checksum = f.Control
	//f.Checksum += f.Address
}

func (f *TelegramShort) CalculateLength() int {
	return 0
}

func (f *TelegramShort) Verify() error {
	if f.Start != FRAME_SHORT_START {
		return fmt.Errorf("no frame start")
	}

	if f.Control != CONTROL_MASK_SND_NKE &&
		f.Control != CONTROL_MASK_REQ_UD1 &&
		f.Control != (CONTROL_MASK_REQ_UD1 | CONTROL_MASK_FCB) &&
		f.Control != CONTROL_MASK_REQ_UD2 &&
		f.Control != (CONTROL_MASK_REQ_UD2 | CONTROL_MASK_FCB) {
		return fmt.Errorf("unknown Controle Code 0x%.2x", f.Control)
	}

	return nil
}

func (f *TelegramShort) Encode() ([]byte, int) {
	pack := []byte{
		f.Start,
		f.Control,
		f.Header.Id[0],
		f.Header.Id[1],
		f.Header.Id[2],
		f.Header.Id[3],
		f.Checksum,
		f.Stop,
	}

	return pack, len(pack)
}
