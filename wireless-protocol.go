package mbus

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"time"
)

type WMBusHeader struct {
	// LSB first
	Manufacturer []byte // 2 bytes

	// LSB first
	Id []byte // 4 bytes

	Version    byte
	DeviceType byte

	AccessNumber byte
	Status       byte

	NEncryptedBlocks int
	EncryptionMode   byte
}

type WMBusFrame struct {
	Start   byte
	Stop    byte
	Length  byte
	Control byte

	Header WMBusHeader

	ControlInformation byte

	// Holds to unprocessed bytes
	Data     []byte
	DataSize int

	// Holds the Data Records
	FrameData MbusDataFrame

	Checksum byte

	Type int

	Timestamp time.Time

	CRCEnabled  bool
	RSSIEnabled bool
}

func NewWirelessMBusFrame() *WMBusFrame {
	return &WMBusFrame{
		Header: WMBusHeader{
			Manufacturer: make([]byte, 2),
			Id:           make([]byte, 4),
		},
		FrameData:   MbusDataFrame{},
		CRCEnabled:  true,
		RSSIEnabled: false,
	}
}

func (frame *WMBusFrame) DecodeSerialNumber() (string, error) {
	var serialNumber int
	if err := DecodeBCDHEX(frame.Header.Id, 4, &serialNumber); err != nil {
		return "", err
	}

	return fmt.Sprintf("%X", serialNumber), nil
}

func (frame *WMBusFrame) DecodeManufacturer() (string, error) {
	var manufacturerId int

	if err := DecodeInt(frame.Header.Manufacturer, len(frame.Header.Manufacturer), &manufacturerId); err != nil {
		return "", err
	}

	return fmt.Sprintf(
		"%s%s%s%s",
		string(((manufacturerId>>10)&0x001F)+64),
		string(((manufacturerId>>5)&0x001F)+64),
		string((manufacturerId&0x001F)+64),
		"",
	), nil
}

func (frame *WMBusFrame) DecodeProductName() (string, error) {
	manufacturer, err := frame.DecodeManufacturer()
	if err != nil {
		return "", err
	}

	_, ok := products[manufacturer]
	if !ok {
		return "", fmt.Errorf("could not find manufacturer: %s", manufacturer)
	}

	if productName, ok := products[manufacturer][frame.Header.Version]; ok {
		return productName, nil
	}

	return "", fmt.Errorf("could not find device type: 0x%.2X, for manufacturer: %s", frame.Header.DeviceType, manufacturer)
}

func (frame *WMBusFrame) DecodeStatus() (string, error) {
	switch frame.Header.Status {
	case 0x00, 0x02, 0x10, 0x20, 0x80:
		return "OK", nil
	case 0x04:
		return "Low battery", nil
	case 0x08:
		return "Permanent error/Sabotage enclosure", nil
	case 0x40:
		return "Sabotage enclosure", nil
	default:
		return "", fmt.Errorf("invalid status")
	}
}

func (frame *WMBusFrame) DecodeDeviceType() (string, error) {
	deviceType, err := DeviceTypeLookup(frame.Header.DeviceType)
	if err != nil {
		return "", err
	}

	return deviceType, nil
}

func (frame *WMBusFrame) ProtocolVersion() (int, error) {
	return int(frame.Header.Version), nil
}

// Method that will return true if there is an encryption mode present
func (frame *WMBusFrame) HasEncryptionMode() bool {
	return frame.Header.EncryptionMode&0x0F != 0
}

// Method that will check if the 2 first data bytes are 0x2F,
// This will indicate if the data is decrypted or not
func (frame *WMBusFrame) IsDecrypted() bool {
	// Only check if first 2 bytes are AES filler bytes when there is an encryption mode set
	if frame.HasEncryptionMode() {
		return frame.Data[0] == 0x2F && frame.Data[1] == 0x2F
	}

	// Always return true when no encryption mode has been set, data was never encrypted
	return true
}

// Returns the IV in Little Endian
// The IV is derived from the manufacturer bytes, the device address and
// the access number from the data header. Note, that None is being
// returned if the current mode does not specify an IV or the IV for that
// specific mode is not implemented.
// Currently implemented IVs are:
//   - IV for mode 2 encryption
//   - IV for mode 4 encryption
//   - IV for mode 5 encryption
func (frame *WMBusFrame) CryptoIV() ([]byte, error) {
	var iv []byte

	switch frame.Header.EncryptionMode & 0x0F {
	case 2:
		iv = []byte{
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		}
		break
	case 4:
		iv = []byte{
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		}
		break
	case 5:
		// According to prEN 13757-3 the IV for mode 5 is setup as follows
		// LSB 1   2   3   4   5   6   7   8   9   10  11  12  13  14  MSB
		// Man Man ID  ..  ..  ID  Ver Med Acc ..  ..  ..  ..  ..  ..  Acc
		// LSB MSB LSB         MSB sio ium
		iv = []byte{
			frame.Header.Manufacturer[0],
			frame.Header.Manufacturer[1],
			frame.Header.Id[0],
			frame.Header.Id[1],
			frame.Header.Id[2],
			frame.Header.Id[3],
			frame.Header.Version,
			frame.Header.DeviceType,
		}

		// The last 8 bytes hold the Access Number
		for i := 0; i < 8; i++ {
			iv = append(iv, frame.Header.AccessNumber)
		}
		break
	}

	if iv != nil {
		return iv, nil
	}

	return nil, fmt.Errorf("unkown encryption mode: 0x%.2X", frame.Header.EncryptionMode)
}

func (frame *WMBusFrame) DecryptData(key []byte) error {
	// No need to decrypt if no Encryption Mode has been set in the Frame
	if !frame.HasEncryptionMode() {
		return nil
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	if len(frame.Data) < aes.BlockSize {
		return fmt.Errorf("ciphertext block size is too short")
	}

	// Get the Crypto IV
	iv, err := frame.CryptoIV()
	if err != nil {
		return err
	}

	if len(frame.Data)%aes.BlockSize != 0 {
		return fmt.Errorf("data length is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	// Replace encrypted bytes with decrypted bytes
	mode.CryptBlocks(frame.Data, frame.Data)

	if DEBUG {
		fmt.Println("Result of decoded data blocks:")
		for i := 0; i < len(frame.Data); i++ {
			fmt.Printf("%.2X ", frame.Data[i])
		}
		fmt.Println()
	}

	// Check if the decryption was successful
	if !frame.IsDecrypted() {
		return fmt.Errorf("error decoding data blocks, first 2 bytes should have the value: 0x2F. Check that you provided the correct key")
	}

	return nil
}

func (frame *WMBusFrame) DataParse() error {
	if DEBUG {
		fmt.Println("Decoding Frame data...")
	}

	// Frame data is encrypted and not yet unencrypted
	if frame.HasEncryptionMode() && !frame.IsDecrypted() {
		return fmt.Errorf("data is not yet decrypted, call `frame.DecryptData(key []byte)` with the correct key")
	}

	if err := frame.DataVariableParse(); err != nil {
		return err
	}

	// recordStartPosition := 0
	//
	//variableRecord := &VariableData{
	//	MoreRecordsFollow: false,
	//}
	//
	//frame.FrameData.Variable = variableRecord
	//
	//for {
	//	if recordStartPosition > frame.DataSize {
	//		break
	//	}
	//
	//	if err := frame.DataVariableParse(&recordStartPosition); err != nil {
	//		return err
	//	}
	//}

	return nil
}

// https://github.com/rscada/libmbus/blob/6edab86078b33f6c870215df2fb605b8fb2fab60/mbus/mbus-protocol.c#L3038
func (frame *WMBusFrame) DataVariableParse() error {
	variableRecord := &VariableData{
		MoreRecordsFollow: false,
	}

	frame.FrameData.Variable = variableRecord

	i := 0
	dr := 0

	for {
		if i > frame.DataSize-1 {
			break
		}

		// Skip filler bytes (0x2F)
		// First 2 encryption verification bytes for example ( if present )
		// And the filler bytes at the end to create the aes blocksize length
		if frame.Data[i]&0xFF == DIB_DIF_IDLE_FILLER {
			if DEBUG {
				if i < 2 {
					fmt.Printf("Skipping encryption verification byte\n")
				} else {
					fmt.Printf("Skipping filler byte\n")
				}
			}

			i++
			continue
		}

		// read and parse DIB (= DIF + DIFE)
		record := &DataRecord{}

		if DEBUG {
			dr++
			fmt.Printf("DR %d :: ", dr)
		}

		// DIF
		record.DIB.DIF = frame.Data[i]

		if DEBUG {
			fmt.Printf("DIB.DIF: 0x%.2X; ", frame.Data[i])
		}

		if record.DIB.DIF == DIB_DIF_MANUFACTURER_SPECIFIC || record.DIB.DIF == DIB_DIF_MORE_RECORDS_FOLLOW {
			if record.DIB.DIF&0xFF == DIB_DIF_MORE_RECORDS_FOLLOW {
				variableRecord.MoreRecordsFollow = true
			}

			i++

			record.DataSize = frame.DataSize - i
			record.Data = make([]byte, record.DataSize)

			for j := 0; j < record.DataSize; j++ {
				i++
				record.Data[j] = frame.Data[i]
			}

			variableRecord.DataRecords = append(variableRecord.DataRecords, record)
			continue
		}

		record.DataSize = DataLengthLookup(record.DIB.DIF)

		if DEBUG {
			fmt.Printf("Data length lookup: %d; ", record.DataSize)
		}

		record.DIB.NDIFe = 0
		record.DIB.DIFe = make([]byte, 10)

		for {
			if i > frame.DataSize || frame.Data[i]&DIB_DIF_EXTENSION_BIT == 0 {
				break
			}

			if record.DIB.NDIFe >= len(record.DIB.DIFe) {
				return fmt.Errorf("too many DIFE")
			}

			dife := frame.Data[i+1]
			record.DIB.DIFe[record.DIB.NDIFe] = dife

			record.DIB.NDIFe++
			i++
		}
		i++

		if i > frame.DataSize {
			return fmt.Errorf("premature end of record at DIF")
		}

		// read and parse VIB (= VIF + VIFE)

		// VIF
		record.VIB.VIF = frame.Data[i]

		if DEBUG {
			fmt.Printf("VIB.VIF: 0x%.2X; ", record.VIB.VIF)
		}

		if record.VIB.VIF&DIB_VIF_WITHOUT_EXTENSION == 0x7C {
			i++
			variableVIFLength := int(frame.Data[i])
			if variableVIFLength > len(record.VIB.Custom) {
				return fmt.Errorf("too long variable length VIF")
			}

			if i+variableVIFLength > frame.DataSize {
				return fmt.Errorf("premature end of record at variable length VIF")
			}

			if err := DecodeString(&frame.Data[i], &variableVIFLength, &record.VIB.Custom); err != nil {
				return err
			}

			i += variableVIFLength
		}

		// VIFE
		record.VIB.NVIFe = 0

		if record.VIB.VIF&DIB_VIF_EXTENSION_BIT != 0 {
			record.VIB.VIFe = make([]byte, 10)
			record.VIB.VIFe[0] = frame.Data[i]
			record.VIB.NVIFe++

			for {
				if i > frame.DataSize || frame.Data[i]&DIB_DIF_EXTENSION_BIT == 0 {
					break
				}

				if record.VIB.NVIFe >= len(record.VIB.VIFe) {
					return fmt.Errorf("too many VIFE")
				}

				vife := frame.Data[i+1]
				record.VIB.VIFe[record.VIB.NVIFe] = vife

				record.VIB.NVIFe++
				i++
			}
			// This should not be here
			// https://github.com/rscada/libmbus/blob/6edab86078b33f6c870215df2fb605b8fb2fab60/mbus/mbus-protocol.c#L3202
			//i++
		}

		if DEBUG {
			if record.VIB.NVIFe > 0 {
				fmt.Printf("VIFe:")
				for i := 0; i < record.VIB.NVIFe-1; i++ {
					fmt.Printf(" 0x%.2X", record.VIB.VIFe[i])
				}
				fmt.Printf("; ")
			} else {
				fmt.Printf("No VIF extension; ")
			}
		}

		if i > frame.DataSize {
			return fmt.Errorf("premature end of record at VIF.")
		}

		// re-calculate data length, if of variable length type
		// 0x0D => Flag for variable data length
		if record.DIB.DIF&DATA_RECORD_DIF_MASK_DATA == 0x0D {
			if frame.Data[i] <= 0xBF {
				i++
				record.DataSize = int(frame.Data[i])
			} else if frame.Data[i] >= 0xC0 && frame.Data[i] <= 0xCF {
				i++
				record.DataSize = (int(frame.Data[i]) - 0xC0) * 2
			} else if frame.Data[i] >= 0xD0 && frame.Data[i] <= 0xDF {
				i++
				record.DataSize = (int(frame.Data[i]) - 0xD0) * 2
			} else if frame.Data[i] >= 0xE0 && frame.Data[i] <= 0xEF {
				i++
				record.DataSize = int(frame.Data[i]) - 0xE0
			} else if frame.Data[i] >= 0xF0 && frame.Data[i] <= 0xFA {
				i++
				record.DataSize = int(frame.Data[i]) - 0xF0
			}
		}

		if DEBUG {
			fmt.Printf("Record datasize: %d; ", record.DataSize)
		}

		if i+record.DataSize > frame.DataSize {
			return fmt.Errorf("premature end of record at data.")
		}

		// Reserve the Record DataSize for the Data byte slice
		record.Data = make([]byte, record.DataSize)

		if DEBUG {
			fmt.Printf("Record data:")
		}

		// Copy the data over
		for j := 0; j < record.DataSize; j++ {
			i++
			record.Data[j] = frame.Data[i]

			if DEBUG {
				fmt.Printf(" 0x%.2X", record.Data[j])
			}
		}

		if DEBUG {
			fmt.Println()
		}

		frame.FrameData.Variable.DataRecords = append(frame.FrameData.Variable.DataRecords, record)

		i++
	}

	return nil
}

// https://github.com/ganehag/pyMeterBus/blob/b41ba61fde1cd9195cd0ec45abacb4ceb154afcd/meterbus/telegram_body.py#L55
//func (frame *WMBusFrame) DataVariableParse(startPos *int) error {
//	lowerBoundary := 0
//	upperBoundary := 0
//
//	record := &DataRecord{}
//
//	record.DIB.DIF = frame.Data[*startPos]
//	if record.DIB.IsEndOfUserData() {
//		if record.DIB.HasMoreRecordsToFollow() {
//			frame.FrameData.Variable.DataRecords = append(frame.FrameData.Variable.DataRecords, record)
//		}
//
//		if record.DIB.IsManufacturerSpecific() {
//			record.Data = make([]byte, len(frame.Data[*startPos:]))
//			record.Data = frame.Data[*startPos:]
//
//			frame.FrameData.Variable.DataRecords = append(frame.FrameData.Variable.DataRecords, record)
//		}
//
//		*startPos = frame.DataSize
//	} else if record.DIB.DIF == DIB_DIF_EXTENSION_BIT {
//		*startPos++
//		return nil
//	}
//
//	return nil
//}

func (frame *WMBusFrame) DecodeDataRecords() ([]DecodedDataRecord, error) {
	//decodedDataRecords := make([]DecodedDataRecord, len(frame.FrameData.Variable.DataRecords))
	var decodedDataRecords []DecodedDataRecord

	for _, record := range frame.FrameData.Variable.DataRecords {
		decodedDataRecord := DecodedDataRecord{
			Function:      record.DecodeRecordFunction(),
			StorageNumber: record.DecodeStorageNumber(),
			//Quantity:      "",
		}

		// Decode the tariff
		tariff, err := record.DecodeTariff()
		if err != nil {
			return nil, err
		} else {
			// Decode the device
			device, _ := record.DecodeDevice()
			decodedDataRecord.Device = device
		}
		decodedDataRecord.Tariff = tariff

		// Decode Unit
		unit, err := record.DecodeUnit()
		if err != nil {
			return nil, err
		}
		decodedDataRecord.Unit = unit.Unit
		decodedDataRecord.Exponent = unit.Exp
		decodedDataRecord.Type = string(unit.Type)

		value, raw, err := record.DecodeValue()
		if err != nil {
			return nil, err
		}
		decodedDataRecord.Value = value
		decodedDataRecord.RawValue = raw

		// Append to the rest
		decodedDataRecords = append(decodedDataRecords, decodedDataRecord)
	}

	return decodedDataRecords, nil
}

func (frame *WMBusFrame) DecodeFrame() (*DecodedFrame, error) {
	decodedFrame := &DecodedFrame{
		ParsedAt: time.Now(),
		Version:  int(frame.Header.Version),
	}

	//Decode serial number
	serialNumber, err := frame.DecodeSerialNumber()
	if err != nil {
		return nil, err
	}
	decodedFrame.SerialNumber = serialNumber

	// Decode manufacturer
	manufacturer, err := frame.DecodeManufacturer()
	if err != nil {
		return nil, err
	}
	decodedFrame.Manufacturer = manufacturer

	// Decode product name
	productName, err := frame.DecodeProductName()
	if err != nil {
		return nil, err
	}
	decodedFrame.ProductName = productName

	// Decode device type
	deviceType, err := frame.DecodeDeviceType()
	if err != nil {
		return nil, err
	}
	decodedFrame.DeviceType = deviceType

	// Decode status
	status, err := frame.DecodeStatus()
	if err != nil {
		return nil, err
	}
	decodedFrame.Status = int(frame.Header.Status)
	decodedFrame.ReadableStatus = status

	// Decode data records
	decodedDeviceRecords, err := frame.DecodeDataRecords()
	if err != nil {
		return nil, err
	}
	decodedFrame.DataRecords = decodedDeviceRecords

	return decodedFrame, nil
}
