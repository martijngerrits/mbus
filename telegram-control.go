package mbus

//import "fmt"
//
//type TelegramControl MBusFrame
//
//func NewTelegramControl() *TelegramControl {
//	frame := &TelegramControl{
//		Start1: FRAME_LONG_START,
//		Start2: FRAME_LONG_START,
//		Stop: FRAME_STOP,
//	}
//
//	return frame
//}
//
//func (f *TelegramControl) CalculateChecksum() byte {
//	checksum := f.Control
//	checksum += f.Address
//	checksum += f.ControlInformation
//
//	return checksum
//	//f.Checksum = f.Control
//	//f.Checksum += f.Address
//	//f.Checksum += f.ControlInformation
//}
//
//func (f *TelegramControl) CalculateLength() int {
//	return 3
//}
//
//func (f *TelegramControl) Verify() error {
//	if f.Start1 != FRAME_LONG_START || f.Start2 != FRAME_LONG_START {
//		return fmt.Errorf("no frame start")
//	}
//
//	if f.Control != CONTROL_MASK_SND_UD &&
//		f.Control != (CONTROL_MASK_SND_UD | CONTROL_MASK_FCB) &&
//		f.Control != CONTROL_MASK_RSP_UD &&
//		f.Control != (CONTROL_MASK_RSP_UD | CONTROL_MASK_DFC) &&
//		f.Control != (CONTROL_MASK_RSP_UD | CONTROL_MASK_ACD) &&
//		f.Control != (CONTROL_MASK_RSP_UD | CONTROL_MASK_DFC | CONTROL_MASK_ACD) {
//		return fmt.Errorf("unkown Control Code 0x%.2x", f.Control)
//	}
//
//	if f.Length1 != f.Length2 {
//		return fmt.Errorf("frame length 1 != 2")
//	}
//
//	if int(f.Length1) != f.CalculateLength() {
//		return fmt.Errorf("frame length1 != calc length")
//	}
//
//	if f.Stop != FRAME_STOP {
//		return fmt.Errorf("no frame stop")
//	}
//
//	checksum := f.CalculateChecksum()
//	if f.Checksum != checksum {
//		return fmt.Errorf("invalid checksum (0x%.2x != 0x%.2x)", f.Checksum, checksum)
//	}
//
//	return nil
//}
//
//func (f *TelegramControl) Encode() ([]byte, int) {
//	pack := []byte{
//		f.Start1,
//		f.Length1,
//		f.Length2,
//		f.Start2,
//		f.Control,
//		f.Address,
//		f.ControlInformation,
//		f.Checksum,
//		f.Stop,
//	}
//
//	return pack, len(pack)
//}
