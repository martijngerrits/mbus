package mbus

import "fmt"

type ParseReturn struct {
	Remaining int
	GotFrame  bool
}

// https://oms-group.org/fileadmin/files/download4all/specification/Vol2/4.2.1/OMS-Spec_Vol2_AnnexN_C042.pdf
func ParseWirelessMBusData(frame *WMBusFrame, data *[]byte, dataSize int) (ParseReturn, error) {
	var frameOffset = 12

	if dataSize <= 0 {
		return ParseReturn{
			Remaining: -1,
			GotFrame:  false,
		}, fmt.Errorf("got no data")
	}

	if DEBUG {
		fmt.Printf("Attempting to parse binary data [size = %d]\n", dataSize)

		for i := 0; i < dataSize; i++ {
			fmt.Printf("%.2X ", (*data)[i]&0xFF)
		}
		fmt.Println()
	}

	frame.Start = (*data)[0]

	switch frame.Start {
	case FRAME_ACK_START:
		frame.Type = FRAME_TYPE_ACK

		if dataSize < FRAME_BASE_SIZE_SHORT {
			// OK, got a valid short packet start, but we need more data
			return ParseReturn{
				Remaining: FRAME_BASE_SIZE_SHORT - dataSize,
				GotFrame:  true,
			}, nil
		}

		if dataSize > FRAME_BASE_SIZE_SHORT {
			return ParseReturn{
				Remaining: -2,
				GotFrame:  true,
			}, fmt.Errorf("too much data in frame")
		}
		break
	case FRAME_SHORT_START:
		frame.Type = FRAME_TYPE_SHORT
		break
	case FRAME_LONG_START:
		frame.Type = FRAME_TYPE_LONG

		if dataSize < 3 {
			// OK, got a valid long/control packet start, but we need
			// more data to determine the length
			return ParseReturn{
				Remaining: 3 - dataSize,
				GotFrame:  true,
			}, nil
		}
		break
	default:
		frame.Type = FRAME_TYPE_ANY
		break
	}

	//************************************
	// Link Layer
	//************************************
	frame.Length = (*data)[1]
	frame.Control = (*data)[2]

	frame.Header.Manufacturer = []byte{(*data)[3], (*data)[4]}
	// The next 4 bytes hold the id (serial number) of the device - LSB first
	frame.Header.Id = []byte{(*data)[5], (*data)[6], (*data)[7], (*data)[8]}

	frame.Header.Version = (*data)[9]
	frame.Header.DeviceType = (*data)[10]

	//************************************
	// Network Layer
	//************************************
	frame.ControlInformation = (*data)[11]

	switch frame.ControlInformation {
	// Short header
	case 0x61, 0x65, 0x6A, 0x6E, 0x74, 0x7A, 0x7B, 0x7D, 0x7F, 0x8A:
		// https://github.com/ganehag/pyMeterBus/blob/bc853aa38ac6b10301bdf97f13ac25b36985316f/meterbus/wtelegram_body.py#L323
		frame.Header.AccessNumber = (*data)[12]
		frame.Header.Status = (*data)[13]
		frame.Header.NEncryptedBlocks = int((*data)[14])
		frame.Header.EncryptionMode = (*data)[15]

		frameOffset += 4 // Excluded the 2 bytes AES Encryption verification
		break
	// Long header
	case 0x60, 0x64, 0x6B, 0x6F, 0x72, 0x37, 0x75, 0x7C, 0x7E, 0x80, 0x8B:
		// https://github.com/ganehag/pyMeterBus/blob/bc853aa38ac6b10301bdf97f13ac25b36985316f/meterbus/wtelegram_body.py#L352
		frameOffset += 12 // Excluded the 2 bytes AES Encryption verification
		break
	// Manufacturer specific layer
	case 0xAA:
		// @TODO: implement manufacturer specific layer
		break
	default:
		return ParseReturn{
			Remaining: -1,
			GotFrame:  false,
		}, fmt.Errorf("no valid Control Information byte: %.2X", frame.ControlInformation)
	}

	//************************************
	// Data Blocks
	//************************************
	// Set the data size
	frame.DataSize = int(frame.Length) - frameOffset

	// According to the data size, we can determine the Frame Type more accurately
	if frame.DataSize == 0 {
		frame.Type = FRAME_TYPE_CONTROL
	} else {
		frame.Type = FRAME_TYPE_LONG
	}

	// Reserve space for the data
	frame.Data = make([]byte, frame.DataSize)
	// Copy over the data
	for i := 0; i < frame.DataSize; i++ {
		frame.Data[i] = (*data)[frameOffset+i]
	}

	// Get the checksum
	frame.Checksum = (*data)[dataSize-2]
	// The last byte is the stop byte
	frame.Stop = (*data)[dataSize-1]

	var err error
	switch frame.Start {
	case FRAME_ACK_START:
		validate := TelegramACK(*frame)
		err = validate.Verify()
		break
	case FRAME_SHORT_START:
		validate := TelegramShort(*frame)
		err = validate.Verify()
		break
	case FRAME_LONG_START:
		validate := TelegramLong(*frame)
		err = validate.Verify()
		break
	}

	if err != nil {
		return ParseReturn{
			Remaining: -3,
			GotFrame:  false,
		}, err
	}

	// Successfully parsed data
	return ParseReturn{
		Remaining: 0,
		GotFrame:  true,
	}, nil
}

//func ParseWirelessMBusData(frame *WMBusFrame, data *[]byte, dataSize int) (ParseReturn, error) {
//    var length int
//
//    if dataSize <= 0 {
//        return ParseReturn{
//            Remaining: -1,
//            GotFrame:  false,
//        }, fmt.Errorf("got no data")
//    }
//
//    if DEBUG {
//        fmt.Printf("Attempting to parse binary data [size = %d]\n", dataSize)
//
//        for i := 0; i < dataSize; i++ {
//            fmt.Printf("%.2X ", (*data)[i] & 0xFF)
//        }
//        fmt.Println()
//    }
//
//    switch (*data)[0] {
//    case FRAME_ACK_START:
//        // OK, got a valid ack frame, require no more data
//        frame.Start = (*data)[0]
//        frame.Type = FRAME_TYPE_ACK
//
//        return ParseReturn{
//            Remaining: 0,
//            GotFrame: true,
//        }, nil
//
//    case FRAME_SHORT_START:
//        if dataSize < FRAME_BASE_SIZE_SHORT {
//            // OK, got a valid short packet start, but we need more data
//            return ParseReturn{
//                Remaining: FRAME_BASE_SIZE_SHORT - dataSize,
//                GotFrame:  true,
//            }, nil
//        }
//
//        if dataSize > FRAME_BASE_SIZE_SHORT {
//            return ParseReturn{
//                Remaining: -2,
//                GotFrame:  true,
//            }, fmt.Errorf("too much data in frame")
//        }
//
//        frame.Start = (*data)[0]
//        frame.Length = (*data)[1]
//        frame.Control = (*data)[2]
//
//        frame.Header.Manufacturer = []byte{(*data)[3], (*data)[4]}
//
//        // The next 4 bytes hold the id (serial number) of the device - LSB first
//        frame.Header.Id = []byte{(*data)[8], (*data)[7], (*data)[6], (*data)[5]}
//
//        frame.Header.Version = (*data)[9]
//        frame.Header.DeviceType = (*data)[10]
//
//        //frame.Checksum = (*data)[3 + frameOffset]
//        frame.Checksum = (*data)[dataSize - 2]
//        frame.Stop = (*data)[dataSize - 1]
//
//        frame.Type = FRAME_TYPE_SHORT
//
//        validate := TelegramShort(*frame)
//        if err := validate.Verify(); err != nil {
//            return ParseReturn{
//                Remaining: -3,
//                GotFrame:  false,
//            }, err
//        }
//
//        // Successfully parsed data
//        return ParseReturn{
//            Remaining: 0,
//            GotFrame:  true,
//        }, nil
//
//    // case FRAME_CONTROL_START: A control frame and a Long frame have the same start byte 0x68
//    case FRAME_LONG_START:
//        if dataSize < 3 {
//            // OK, got a valid long/control packet start, but we need
//            // more data to determine the length
//            return ParseReturn{
//                Remaining: 3 - dataSize,
//                GotFrame:  true,
//            }, nil
//        }
//
//        frame.Start = (*data)[0]
//        frame.Length = (*data)[1]
//        frame.Control = (*data)[2]
//
//        // Early verify control to see if we did'nt get a FRAME_LONG_START byte
//        // in the middle of another frame which was still in the read buffer
//        validate := TelegramLong(*frame)
//        if err := validate.VerifyControl(); err != nil {
//            return ParseReturn{
//                Remaining: -2,
//                GotFrame:  false,
//            }, err
//        }
//
//        if frame.Length < 3 {
//            // not a valid M-bus frame
//            return ParseReturn{
//                Remaining: -2,
//                GotFrame:  false,
//            }, fmt.Errorf("invalid M-Bus frame length")
//        }
//
//        // Make up for the Start & Stop bytes and the Length byte itself,
//        // those are not included in the Length calculation.
//        if int(frame.Length) != dataSize - 3 {
//            return ParseReturn{
//                // Normally the entire frame exists of {Start+Length+Data+Stop} which results in a remaining
//                // of Length + 3, but at this point we already got 3 bytes, so the remaining length is:
//                // Length + 3 - 3, or just: Length
//                Remaining: int(frame.Length),
//                GotFrame:  true,
//            }, nil
//        }
//
//        frame.Header.Manufacturer = []byte{(*data)[3], (*data)[4]}
//
//        // The next 4 bytes hold the id (serial number) of the device - LSB first
//        frame.Header.Id = []byte{(*data)[5], (*data)[6], (*data)[7], (*data)[8]}
//
//        frame.Header.Version = (*data)[9]
//        frame.Header.DeviceType = (*data)[10]
//        frame.ControlInformation = (*data)[11]
//
//
//
//
//        frame.Header.AccessNumber = (*data)[12]
//        frame.Header.Status = (*data)[13]
//        frame.Header.NEncryptedBlocks = int((*data)[14])
//        frame.Header.EncryptionMode = (*data)[15]
//
//        d i := 0; i < frame.DataSize; i++ {
//            frame.Data[i] = (*data)[16 + i]
//        }
//
//        frame.Checksum = (*data)[dataSize - 2]
//        // The last byte is the stop byte
//        frame.Stop = (*data)[dataSize - 1]
//
//        if frame.DataSize == 0 {
//            frame.Type = FRAME_TYPE_CONTROL
//        } else {
//            frame.Type = FRAME_TYPE_LONG
//        }
//
//        validate = TelegramLong(*frame)
//        if err := validate.Verify(); err != nil {
//            return ParseReturn{
//                Remaining: -3,
//                GotFrame:  false,
//            }, err
//        }
//
//        // Successfully parsed data
//        return ParseReturn{
//            Remaining: 0,
//            GotFrame:  true,
//        }, nil
//
//    default:
//        return ParseReturn{
//            Remaining: 1,
//            GotFrame:  false,
//        }, nil // fmt.Errorf("invalid M-Bus frame start")
//    }
//}

//func ParseWiredMBusData(frame *WiredMBusFrame, data *[]byte, dataSize int) (int, error) {
//    var length int
//
//    if dataSize <= 0 {
//        return -1, nil
//    }
//
//    fmt.Printf("Attempting to parse binary data [size = %d]\n", dataSize)
//
//    for i := 0; i < dataSize; i++ {
//        fmt.Printf("%.2X ", (*data)[i] & 0xFF)
//    }
//    fmt.Println()
//
//    switch (*data)[0] {
//    case FRAME_ACK_START:
//        // OK, got a valid ack frame, require no more data
//        frame.Start1 = (*data)[0]
//        frame.Type = FRAME_TYPE_ACK
//
//        return 0, nil
//
//    case FRAME_SHORT_START:
//        if dataSize < FRAME_BASE_SIZE_SHORT {
//            // OK, got a valid short packet start, but we need more data
//            return FRAME_BASE_SIZE_SHORT - dataSize, nil
//        }
//
//        if dataSize != FRAME_BASE_SIZE_SHORT {
//            return -2, fmt.Errorf("too much data in frame")
//        }
//
//        frame.Start1 = (*data)[0]
//        frame.Control = (*data)[1]
//        frame.Address = (*data)[2]
//        frame.Checksum = (*data)[3]
//        frame.Stop = (*data)[4]
//
//        frame.Type = FRAME_TYPE_SHORT
//
//        validate := TelegramShort(*frame)
//        if err := validate.Verify(); err != nil {
//            return -3, err
//        }
//
//        // Successfully parsed data
//        return 0, nil
//
//    case FRAME_LONG_START:
//        if dataSize < 3 {
//            // OK, got a valid long/control packet start, but we need
//            // more data to determine the length
//            return 3 - dataSize, nil
//        }
//
//        frame.Start1 = (*data)[0]
//        frame.Length1 = (*data)[1]
//        frame.Length2 = (*data)[2]
//
//        if frame.Length1 < 3 || frame.Length1 != frame.Length2 {
//            // not a valid M-bus frame
//            return -2, fmt.Errorf("invalid M-Bus frame length")
//        }
//
//        // check length of packet:
//        length = int(frame.Length1)
//
//        if dataSize < FRAME_FIXED_SIZE_LONG + length {
//            fmt.Printf("OK, but we need more data. %d %d \n", dataSize, length)
//            // OK, but we need more data
//            return FRAME_FIXED_SIZE_LONG + length - dataSize, nil
//        }
//
//        if dataSize > FRAME_FIXED_SIZE_LONG + length {
//            return -2, fmt.Errorf("too much data in frame")
//        }
//
//        frame.Start2 = (*data)[3]
//        frame.Control = (*data)[4]
//        frame.Address = (*data)[5]
//        frame.ControlInformation = (*data)[6]
//
//        frame.DataSize = length - 3
//        for i := 0; i < frame.DataSize; i++ {
//            frame.Data[i] = (*data)[7 + i]
//        }
//
//        frame.Checksum = (*data)[dataSize - 2]
//        frame.Stop = (*data)[dataSize - 1]
//
//        if frame.DataSize == 0 {
//            frame.Type = FRAME_TYPE_CONTROL
//        } else {
//            frame.Type = FRAME_TYPE_LONG
//        }
//
//        validate := TelegramLong(*frame)
//        if err := validate.Verify(); err != nil {
//           return -3, err
//        }
//
//        // Successfully parsed data
//        return 0, nil
//
//    default:
//        return -4, fmt.Errorf("invalid M-Bus frame start")
//    }
//}

//
//func (frame *MBusFrame) Verify() error {
//    switch frame.Type {
//    case FRAME_TYPE_ACK:
//        if frame.Start1 != FRAME_ACK_START {
//            return fmt.Errorf("no valid ack type")
//        }
//
//    case FRAME_TYPE_SHORT:
//        if frame.Start1 != FRAME_SHORT_START {
//            return fmt.Errorf("no frame start")
//        }
//
//        if frame.Control != CONTROL_MASK_SND_NKE &&
//            frame.Control != CONTROL_MASK_REQ_UD1 &&
//            frame.Control != (CONTROL_MASK_REQ_UD1 | CONTROL_MASK_FCB) &&
//            frame.Control != CONTROL_MASK_REQ_UD2 &&
//            frame.Control != (CONTROL_MASK_REQ_UD2 | CONTROL_MASK_FCB) {
//            return fmt.Errorf("unknown control code 0x%.2x", frame.Control)
//        }
//
//    case FRAME_TYPE_CONTROL:
//    case FRAME_TYPE_LONG:
//        if frame.Start1 != FRAME_CONTROL_START || frame.Start1 != FRAME_LONG_START {
//            return fmt.Errorf("no frame start")
//        }
//
//        if frame.Control != CONTROL_MASK_SND_UD &&
//            frame.Control != (CONTROL_MASK_SND_UD | CONTROL_MASK_FCB) &&
//            frame.Control != CONTROL_MASK_RSP_UD &&
//            frame.Control != (CONTROL_MASK_RSP_UD | CONTROL_MASK_DFC) &&
//            frame.Control != (CONTROL_MASK_RSP_UD | CONTROL_MASK_ACD) &&
//            frame.Control != (CONTROL_MASK_RSP_UD | CONTROL_MASK_DFC | CONTROL_MASK_ACD) {
//            return fmt.Errorf("unknown control code 0x%.2x", frame.Control)
//        }
//
//        if frame.Length1 != frame.CalculateLength()
//    }
//    return nil
//}
