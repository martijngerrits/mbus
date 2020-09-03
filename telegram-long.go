package mbus

import (
	"fmt"
)

type TelegramLong WMBusFrame

func NewTelegramLong() *TelegramLong {
	frame := &TelegramLong{
		Start: FRAME_LONG_START,
		Stop:  FRAME_STOP,
	}

	return frame
}

func (f *TelegramLong) CalculateChecksum() byte {
	checksum := f.Control

	checksum += f.ControlInformation

	//if f.HeaderEnabled {
	//	checksum += f.Manufacturer[0]
	//	checksum += f.Manufacturer[1]
	//
	//	for i := 0; i < len(f.Address); i++ {
	//		checksum += f.Address[i]
	//	}
	//
	//	checksum += f.Version
	//	checksum += f.DeviceType
	//}

	for i := 0; i < f.DataSize; i++ {
		checksum += f.Data[i]
	}

	return checksum

	//f.Checksum = f.Control
	//f.Checksum += f.Address
	//f.Checksum += f.ControlInformation
	//
	//for i:= 0; uint64(i) < f.DataSize; i++ {
	//	f.Checksum += f.Data[i]
	//}
}

func (f *TelegramLong) CalculateLength() int {
	// C (Control) byte
	// CI (Control Information) byte
	// 5 bytes Network Layer;
	//   - CI
	//   - Acc
	//   - Status
	//   - NEncryptedBlocks
	//   - Encryption mode
	addSize := 15

	if f.RSSIEnabled {
		addSize += 1
	}

	if f.CRCEnabled {
		addSize += 1
	}

	return f.DataSize + addSize
}

func (f *TelegramLong) VerifyControl() error {
	if f.Control != CONTROL_MASK_SND_UD &&
		f.Control != CONTROL_MASK_SND_NR &&
		f.Control != (CONTROL_MASK_SND_UD|CONTROL_MASK_FCB) &&
		f.Control != CONTROL_MASK_RSP_UD &&
		f.Control != (CONTROL_MASK_RSP_UD|CONTROL_MASK_DFC) &&
		f.Control != (CONTROL_MASK_RSP_UD|CONTROL_MASK_ACD) &&
		f.Control != (CONTROL_MASK_RSP_UD|CONTROL_MASK_DFC|CONTROL_MASK_ACD) {
		return fmt.Errorf("unkown Control Code 0x%.2x", f.Control)
	}

	return nil
}

func (f *TelegramLong) Verify() error {
	if f.Start != FRAME_LONG_START {
		return fmt.Errorf("no frame start")
	}

	if err := f.VerifyControl(); err != nil {
		return err
	}

	calcLength := f.CalculateLength()
	if int(f.Length) != calcLength {
		return fmt.Errorf("frame length (%d) != calc length (%d)", f.Length, calcLength)
	}

	if f.Stop != FRAME_STOP {
		return fmt.Errorf("no frame stop")
	}

	// @TODO: Fix calculate checksum
	//checksum := f.CalculateChecksum()
	//if f.Checksum != checksum {
	//	return fmt.Errorf("invalid checksum (0x%.2x != 0x%.2x)", f.Checksum, checksum)
	//}

	return nil
}

func (f *TelegramLong) Encode() ([]byte, int) {
	pack := []byte{
		f.Start,
		f.Length,
		f.Control,
		f.Header.Id[0],
		f.Header.Id[1],
		f.Header.Id[2],
		f.Header.Id[3],
		f.ControlInformation,
	}

	// Total pack size
	offset := 7

	for i := 0; i < f.DataSize; i++ {
		pack = append(pack, f.Data[i])
		offset++
	}

	pack = append(pack, f.Checksum)
	offset++

	pack = append(pack, f.Stop)
	offset++

	return pack, offset
}
