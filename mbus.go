package mbus

import (
	"bytes"
	"context"
	"fmt"
	"math"
	"time"
)

const (
	PACKET_BUFF_SIZE = 2048
)

var (
	DEBUG = false
)

type Frame interface {
	//------------------------------------------------------------------------------
	/// Calculate the checksum of the M-Bus frame. The checksum algorithm is the
	/// arithmetic sum of the frame content, without using carry. Which content
	/// that is included in the checksum calculation depends on the frame type.
	//------------------------------------------------------------------------------
	//CalculateChecksum() byte
	//
	//CalculateLength() int

	//Verify() error

	//Encode() ([]byte, int)

	//DecodeTariff() (int, error)
	//DecodeDevice() (int, error)
	//DecodeUnit() (string, error)

	DataParse() error

	DecodeFrame() (*DecodedFrame, error)

	DecodeManufacturer() (string, error)
	DecodeSerialNumber() (string, error)
	DecodeProductName() (string, error)
	DecodeDeviceType() (string, error)
	DecodeDataRecords() ([]DecodedDataRecord, error)

	ProtocolVersion() (int, error)

	HasEncryptionMode() bool
	DecryptData(key []byte) error
	IsDecrypted() bool
}

type MbusHandle struct {
	Fd interface{} // Can be either Serial or TCP

	MaxDataRetry   int
	MaxSearchRetry int
	IsSerial       bool
}

type Handle interface {
	Open(device string, config interface{}) error
	Close() error

	// Send
	Send(frame Frame) error

	// Receive frame through a channel
	Stream(ctx context.Context) chan Frame
	// Receive a single frame from the serial buffer
	ReceiveFrame() (Frame, error)
}

type Device struct {
	SerialNumber string
	AESKey       []byte
}

type DataInformationBlock struct {
	DIF  byte
	DIFe []byte

	NDIFe int
}

type ValueInformationBlock struct {
	VIF  byte
	VIFe []byte

	NVIFe int

	Custom string
}

type DataRecord struct {
	DIB DataInformationBlock
	VIB ValueInformationBlock

	Data     []byte
	DataSize int

	Timestamp time.Time
}

type VariableData struct {
	DataRecords []*DataRecord

	Data     *[]byte
	DataSize int

	MoreRecordsFollow bool
}

type MbusDataFrame struct {
	Type  int
	Error int

	Variable *VariableData
	Fixed    byte
}

type DecodedDataRecord struct {
	Function      string
	StorageNumber int

	Tariff int
	Device int

	Unit     string
	Exponent float64
	Type     string
	Quantity string

	Value    string
	RawValue []byte
}

type DecodedFrame struct {
	SerialNumber string
	Manufacturer string
	ProductName  string
	Version      int
	DeviceType   string
	AccessNumber int16

	Signature int16

	Status         int
	ReadableStatus string

	DataRecords []DecodedDataRecord

	ParsedAt time.Time
}

func DeviceTypeLookup(deviceType byte) (string, error) {
	buffer := bytes.Buffer{}
	var err error

	switch deviceType {
	case VARIABLE_DATA_MEDIUM_OTHER:
		_, err = fmt.Fprint(&buffer, "Other")
		break
	case VARIABLE_DATA_MEDIUM_OIL:
		_, err = fmt.Fprint(&buffer, "Oild")
		break
	case VARIABLE_DATA_MEDIUM_ELECTRICITY:
		_, err = fmt.Fprint(&buffer, "Electricity")
		break
	case VARIABLE_DATA_MEDIUM_GAS:
		_, err = fmt.Fprint(&buffer, "Gas")
		break
	case VARIABLE_DATA_MEDIUM_HEAT_OUT:
		_, err = fmt.Fprint(&buffer, "Heat: Outlet")
		break
	case VARIABLE_DATA_MEDIUM_STEAM:
		_, err = fmt.Fprint(&buffer, "Steam")
		break
	case VARIABLE_DATA_MEDIUM_HOT_WATER:
		_, err = fmt.Fprint(&buffer, "Warm water (30-90°C)")
		break
	case VARIABLE_DATA_MEDIUM_WATER:
		_, err = fmt.Fprint(&buffer, "Warm")
		break
	case VARIABLE_DATA_MEDIUM_HEAT_COST:
		_, err = fmt.Fprint(&buffer, "Heat Cost Allocator")
		break
	case VARIABLE_DATA_MEDIUM_COMPR_AIR:
		_, err = fmt.Fprint(&buffer, "Compressed Air")
		break
	case VARIABLE_DATA_MEDIUM_COOL_OUT:
		_, err = fmt.Fprint(&buffer, "Cooling load meter: Outlet")
		break
	case VARIABLE_DATA_MEDIUM_COOL_IN:
		_, err = fmt.Fprint(&buffer, "Cooling load meter: Inlet")
		break
	case VARIABLE_DATA_MEDIUM_HEAT_IN:
		_, err = fmt.Fprint(&buffer, "Heat: Inlet")
		break
	case VARIABLE_DATA_MEDIUM_HEAT_COOL:
		_, err = fmt.Fprint(&buffer, "Heat / Cooling load meter")
		break
	case VARIABLE_DATA_MEDIUM_BUS:
		_, err = fmt.Fprint(&buffer, "Bus / System")
		break
	case VARIABLE_DATA_MEDIUM_UNKNOWN:
		_, err = fmt.Fprint(&buffer, "Unknown Device type")
		break
	case VARIABLE_DATA_MEDIUM_IRRIGATION:
		_, err = fmt.Fprint(&buffer, "Irrigation Water")
		break
	case VARIABLE_DATA_MEDIUM_WATER_LOGGER:
		_, err = fmt.Fprint(&buffer, "Water Logger")
		break
	case VARIABLE_DATA_MEDIUM_GAS_LOGGER:
		_, err = fmt.Fprint(&buffer, "Gas Logger")
		break
	case VARIABLE_DATA_MEDIUM_GAS_CONV:
		_, err = fmt.Fprint(&buffer, "Gas Converter")
		break
	case VARIABLE_DATA_MEDIUM_COLORIFIC:
		_, err = fmt.Fprint(&buffer, "Calorific value")
		break
	case VARIABLE_DATA_MEDIUM_BOIL_WATER:
		_, err = fmt.Fprint(&buffer, "Hot water (>90°C)")
		break
	case VARIABLE_DATA_MEDIUM_COLD_WATER:
		_, err = fmt.Fprint(&buffer, "Cold water")
		break
	case VARIABLE_DATA_MEDIUM_DUAL_WATER:
		_, err = fmt.Fprint(&buffer, "Dual water")
		break
	case VARIABLE_DATA_MEDIUM_PRESSURE:
		_, err = fmt.Fprint(&buffer, "Pressure")
		break
	case VARIABLE_DATA_MEDIUM_ADC:
		_, err = fmt.Fprint(&buffer, "A/D Converter")
		break
	case VARIABLE_DATA_MEDIUM_SMOKE:
		_, err = fmt.Fprint(&buffer, "Smoke Detector")
		break
	case VARIABLE_DATA_MEDIUM_ROOM_SENSOR:
		_, err = fmt.Fprint(&buffer, "Ambient Sensor")
		break
	case VARIABLE_DATA_MEDIUM_GAS_DETECTOR:
		_, err = fmt.Fprint(&buffer, "Gas Detector")
		break
	case VARIABLE_DATA_MEDIUM_BREAKER_E:
		_, err = fmt.Fprint(&buffer, "Breaker: Electricity")
		break
	case VARIABLE_DATA_MEDIUM_VALVE:
		_, err = fmt.Fprint(&buffer, "Valve: Gas or Water")
		break
	case VARIABLE_DATA_MEDIUM_CUSTOMER_UNIT:
		_, err = fmt.Fprint(&buffer, "Customer Unit: Display Device")
		break
	case VARIABLE_DATA_MEDIUM_WASTE_WATER:
		_, err = fmt.Fprint(&buffer, "Waste Water")
		break
	case VARIABLE_DATA_MEDIUM_GARBAGE:
		_, err = fmt.Fprint(&buffer, "Garbage")
		break
	case VARIABLE_DATA_MEDIUM_VOC:
		_, err = fmt.Fprint(&buffer, "VOC Sensor")
		break
	case VARIABLE_DATA_MEDIUM_SERVICE_UNIT:
		_, err = fmt.Fprint(&buffer, "Service Unit")
		break
	case VARIABLE_DATA_MEDIUM_RC_SYSTEM:
		_, err = fmt.Fprint(&buffer, "Radio Converter: System")
		break
	case VARIABLE_DATA_MEDIUM_RC_METER:
		_, err = fmt.Fprint(&buffer, "Radio Converter: Meter")
		break
	case 0x22, 0x23, 0x24, 0x26, 0x27, 0x2A, 0x2C, 0x2D, 0x2E, 0x2F, 0x31, 0x32, 0x33, 0x34, 0x38, 0x39, 0x3A, 0x3B, 0x3C, 0x3D, 0x3E, 0x3F:
		_, err = fmt.Fprint(&buffer, "Reserved")
		break

	// add more ...
	default:
		_, err = fmt.Fprintf(&buffer, "Unknown medium (0x%.2x)", deviceType)
		break
	}

	if err != nil {
		return "", err
	}

	return buffer.String(), nil
}

func DataLengthLookup(dif byte) int {
	switch dif & DATA_RECORD_DIF_MASK_DATA {
	case 0x0:
		return 0
	case 0x1:
		return 1
	case 0x2:
		return 2
	case 0x3:
		return 3
	case 0x4:
		return 4
	case 0x5:
		return 4
	case 0x6:
		return 6
	case 0x7:
		return 8
	case 0x8:
		return 0
	case 0x9:
		return 1
	case 0xA:
		return 2
	case 0xB:
		return 3
	case 0xC:
		return 4
	case 0xD:
		return 0
	case 0xE:
		return 6
	case 0xF:
		return 8
	default:
		return 0x00
	}
}

func unitPrefix(exp int) string {
	switch exp {
	case 0:
		return "0"
	case -3:
		return "m"
	case -6:
		return "my"
	case 1:
		return "10 "
	case 2:
		return "100 "
	case 3:
		return "k"
	case 4:
		return "10 k"
	case 5:
		return "100 k"
	case 6:
		return "M"
	case 9:
		return "G"
	default:
		return fmt.Sprintf("1e%d ", exp)
	}
}

func unitDurationNN(pp int) string {
	switch pp {
	case 0:
		return "hour(s)"
	case 1:
		return "day(s)"
	case 2:
		return "month(s)"
	case 3:
		return "year(s)"
	default:
		return "error: out-of-range"
	}
}

func (dib *DataInformationBlock) IsEndOfUserData() bool {
	return dib.DIF == 0x0F || dib.DIF == 0x1F
}

func (dib *DataInformationBlock) HasMoreRecordsToFollow() bool {
	return dib.DIF == DIB_DIF_MORE_RECORDS_FOLLOW
}

func (dib *DataInformationBlock) IsManufacturerSpecific() bool {
	return dib.DIF == 0x0F
}

func (vib *ValueInformationBlock) UnitLookup() VIF {
	var code int

	if vib.VIF == 0xFB {
		code = int(vib.VIFe[1])&DIB_VIF_WITHOUT_EXTENSION | 0x200
	} else if vib.VIF == 0xFD {
		code = int(vib.VIFe[1])&DIB_VIF_WITHOUT_EXTENSION | 0x100
	} else if vib.VIF == 0x7C {
		//var unit string
		//DecodeASCII(vib.Custom, &unit)
		return VIF{
			Exp: 1,
			//Unit:        unit,
			Unit:        vib.Custom,
			Type:        VIFUnit["VARIABLE_VIF"],
			VIFUnitDesc: "",
		}
	} else if vib.VIF == 0xFC {
		//  && (vib->vife[0] & 0x78) == 0x70

		// Disable this for now as it is implicit
		// from 0xFC
		// if vif & vtf_ebm {}
		code := vib.VIFe[0] & DIB_VIF_WITHOUT_EXTENSION
		var factor float64

		if 0x70 <= code && code <= 0x77 {
			factor = math.Pow10((int(vib.VIFe[0]) & 0x07) - 6)
		} else if 0x78 <= code && code <= 0x7B {
			factor = math.Pow10((int(vib.VIFe[0]) & 0x03) - 3)
		} else if code == 0x7D {
			// A bit unnecessary
			factor = 1
		}

		return VIF{
			Exp:         factor,
			Unit:        vib.Custom,
			Type:        VIFUnit["VARIABLE_VIF"],
			VIFUnitDesc: "",
		}
	} else {
		code = int(vib.VIF) & DIB_VIF_WITHOUT_EXTENSION
	}

	return VIFTable[code]
}

//func (vib *ValueInformationBlock) UnitLookupFB() string {
//    var prefix string
//
//    switch int(vib.VIFe[1]) & DIB_VIF_WITHOUT_EXTENSION | 0x200 {
//    case 0x0:
//    case 0x0 + 1:
//        // E000 000n
//        n := 0x01 & vib.VIFe[0]
//        if n == 0 {
//            prefix = "0.1 "
//        }
//        return fmt.Sprintf("Energy (%sMWh)", prefix)
//    case 0x2:
//    case 0x2 + 1:
//        // E000 001n
//    case 0x4:
//    case 0x4 + 1:
//    case 0x4 + 2:
//    case 0x4 + 3:
//        // E000 01nn
//        return fmt.Sprintf("Reserved (0x%.2x)", vib.VIFe[0])
//    case 0x8:
//    case 0x8 + 1:
//        // E000 100n
//        n := 0x01 & vib.VIFe[0]
//        if n == 0 {
//            prefix = "0.1 "
//        }
//        return fmt.Sprintf("Energy (%sGJ)", prefix)
//    case 0xA:
//    case 0xA + 1:
//    case 0xC:
//    case 0xC + 1:
//    case 0xC + 2:
//    case 0xC + 3:
//        // E000 101n
//        // E000 11nn
//        return fmt.Sprintf("Reserved (0x%.2x)", vib.VIFe[0])
//    case 0x10:
//    case 0x10 + 1:
//        // E001 000n
//        n := 0x01 & vib.VIFe[0]
//        return fmt.Sprintf("Volume (%sm3)", unitPrefix(int(n) + 2))
//    case 0x12:
//    case 0x12 + 1:
//    case 0x14:
//    case 0x14 + 1:
//    case 0x14 + 2:
//    case 0x14 + 3:
//        // E001 001n
//        // E001 01nn
//        return fmt.Sprintf("Reserved (0x%.2x)", vib.VIFe[0])
//    case 0x18:
//    case 0x18 + 1:
//        // E001 100n
//        n := 0x01 & vib.VIFe[0]
//        return fmt.Sprintf("Mass (%st)", unitPrefix(int(n) + 2))
//    case 0x1A:
//    case 0x1B:
//    case 0x1C:
//    case 0x1D:
//    case 0x1E:
//    case 0x1F:
//    case 0x20:
//        // E001 1010 to E010 0000, Reserved
//        return fmt.Sprintf("Reserved (0x%.2x)", vib.VIFe[0])
//    case 0x21:
//        // E010 0001
//        return "Volume (0.1 feet^3)"
//    case 0x22:
//    case 0x23:
//        // E010 0010
//        // E010 0011
//        n := 0x01 & vib.VIFe[0]
//        if n == 0 {
//            prefix = "0.1 "
//        }
//        return fmt.Sprintf("Volume (%samerican gallon)", prefix)
//    case 0x24:
//        // E010 0100
//        return fmt.Sprintf("Volume flow (0.001 american gallon/min)")
//    case 0x25:
//        // E010 0101
//        return fmt.Sprintf("Volume flow (american gallon/min)")
//    case 0x26:
//        // E010 0110
//        return fmt.Sprintf("Volume flow (american gallon/h)")
//    case 0x27:
//        // E010 0111, Reserved
//        return fmt.Sprintf("Reserved (0x%.2x)", vib.VIFe[0])
//    case 0x28:
//    case 0x28 + 1:
//        // E010 100n
//        n := 0x01 & vib.VIFe[0]
//        if n == 0 {
//            prefix = "0.1 "
//        }
//        return fmt.Sprintf("Power (%sMW)", prefix)
//    case 0x2A:
//    case 0x2A + 1:
//    case 0x2C:
//    case 0x2C + 1:
//    case 0x2C + 2:
//    case 0x2C + 3:
//        // E010 101n, Reserved
//        // E010 11nn, Reserved
//        return fmt.Sprintf("Reserved (0x%.2x)", vib.VIFe[0])
//    case 0x30:
//    case 0x30 + 1:
//        // E011 000n
//        n := 0x01 & vib.VIFe[0]
//        if n == 0 {
//            prefix = "0.1 "
//        }
//
//        return fmt.Sprintf("Power (%sGJ/h)", prefix)
//    case 0x32:
//    case 0x33:
//    case 0x34:
//    case 0x35:
//    case 0x36:
//    case 0x37:
//    case 0x38:
//    case 0x39:
//    case 0x3A:
//    case 0x3B:
//    case 0x3C:
//    case 0x3D:
//    case 0x3E:
//    case 0x3F:
//
//    case 0x40:
//    case 0x41:
//    case 0x42:
//    case 0x43:
//    case 0x44:
//    case 0x45:
//    case 0x46:
//    case 0x47:
//    case 0x48:
//    case 0x49:
//    case 0x4A:
//    case 0x4B:
//    case 0x4C:
//    case 0x4D:
//    case 0x4E:
//    case 0x4F:
//
//    case 0x52:
//    case 0x53:
//    case 0x54:
//    case 0x55:
//    case 0x56:
//    case 0x57:
//        // E011 0010 to E101 0111
//        return fmt.Sprintf("Reserved (0x%.2x)", vib.VIFe[0])
//    case 0x58:
//    case 0x58 + 1:
//    case 0x58 + 2:
//    case 0x58 + 3:
//        // E101 10nn
//        n := 0x03 & vib.VIFe[0]
//        return fmt.Sprintf("Flow Temperature (%s degree F)", unitPrefix(int(n) -3))
//    case 0x5C:
//    case 0x5C + 1:
//    case 0x5C + 2:
//    case 0x5C + 3:
//        // E101 11nn
//        n := 0x03 & vib.VIFe[0]
//        return fmt.Sprintf("Return Temperature (%s degree F)", unitPrefix(int(n) -3))
//    case 0x60:
//    case 0x60 + 1:
//    case 0x60 + 2:
//    case 0x60 + 3:
//        // E110 00nn
//        n := 0x03 & vib.VIFe[0]
//        return fmt.Sprintf("Temperature Difference (%s degree F)", unitPrefix(int(n) -3))
//    case 0x64:
//    case 0x64 + 1:
//    case 0x64 + 2:
//    case 0x64 + 3:
//        // E110 01nn
//        n := 0x03 & vib.VIFe[0]
//        return fmt.Sprintf("External Temperature (%s degree F)", unitPrefix(int(n) -3))
//    case 0x68:
//    case 0x69:
//    case 0x6A:
//    case 0x6B:
//    case 0x6C:
//    case 0x6D:
//    case 0x6E:
//    case 0x6F:
//        // E110 1nnn
//        return fmt.Sprintf("Reserved (0x%.2x)", vib.VIFe[0])
//    case 0x70:
//    case 0x70 + 1:
//    case 0x70 + 2:
//    case 0x70 + 3:
//        // E111 00nn
//        n := 0x03 & vib.VIFe[0]
//        return fmt.Sprintf("Cold / Warm Temperature Limit (%s degree F)", unitPrefix(int(n) -3))
//    case 0x74:
//    case 0x74 + 1:
//    case 0x74 + 2:
//    case 0x74 + 3:
//        // E111 00nn
//        n := 0x03 & vib.VIFe[0]
//        return fmt.Sprintf("Cold / Warm Temperature Limit (%s degree C)", unitPrefix(int(n) -3))
//    case 0x78, 0x78 + 1, 0x78 + 2, 0x78 + 3, 0x78 + 4, 0x78 + 5, 0x78 + 6, 0x78 + 7:
//        // E111 1nnn
//        n := 0x07 & vib.VIFe[0]
//        return fmt.Sprintf("cumul. count max power (%s W)", unitPrefix(int(n) - 3))
//    default:
//        return fmt.Sprintf("Unrecognized VIF 0xFB extension: 0x%.2x", vib.VIFe[0])
//    }
//
//    return "unknown unit"
//}
//
//func (vib *ValueInformationBlock) UnitLookupFD() string {
//    maskedVIFe0 := int(vib.VIFe[1]) & DIB_VIF_WITHOUT_EXTENSION | 0x100
//
//    if maskedVIFe0 & 0x7C == 0x00 {
//        // VIFE = E000 00nn	Credit of 10nn-3 of the nominal local legal currency units
//        n := maskedVIFe0 & 0x03
//        return fmt.Sprintf("Credit of %s of the nominal local legal currency units", unitPrefix(int(n - 3)))
//    } else if maskedVIFe0 & 0x7C == 0x04 {
//        // VIFE = E000 01nn Debit of 10nn-3 of the nominal local legal currency units
//        n := maskedVIFe0 & 0x03
//        return fmt.Sprintf("Debit of %s of the nominal local legal currency units", unitPrefix(int(n - 3)))
//    } else if maskedVIFe0 == 0x08 {
//        // E000 1000
//        return fmt.Sprintf("Access Number (transmission count)")
//    } else if maskedVIFe0 == 0x09 {
//        // E000 1001
//        return fmt.Sprintf("Medium (as in fixed header)")
//    } else if maskedVIFe0 == 0x0A {
//        // E000 1010
//        return fmt.Sprintf("Manufacturer (as in fixed header)")
//    } else if maskedVIFe0 == 0x0B {
//        // E000 1010
//        return fmt.Sprintf("Parameter set identification")
//    } else if maskedVIFe0 == 0x0C {
//        // E000 1100
//        return fmt.Sprintf("Model / Version")
//    } else if maskedVIFe0 == 0x0D {
//        // E000 1100
//        return fmt.Sprintf("Hardware version")
//    } else if maskedVIFe0 == 0x0E {
//        // E000 1101
//        return fmt.Sprintf("Firmware version")
//    } else if maskedVIFe0 == 0x0F {
//        // E000 1101
//        return fmt.Sprintf("Software version")
//    } else if maskedVIFe0 == 0x10 {
//        // VIFE = E001 0000 Customer location
//        return fmt.Sprintf("Customer location")
//    } else if maskedVIFe0 == 0x11 {
//        // VIFE = E001 0001 Customer
//        return fmt.Sprintf("Customer")
//    } else if maskedVIFe0 == 0x12 {
//        // VIFE = E001 0010	Access Code User
//        return fmt.Sprintf("Access Code User")
//    } else if maskedVIFe0 == 0x13 {
//        // VIFE = E001 0011	Access Code Operator
//        return fmt.Sprintf("Access Code Operator")
//    } else if maskedVIFe0 == 0x14 {
//        // VIFE = E001 0100	Access Code System Operator
//        return fmt.Sprintf("Access Code System Operator")
//    } else if maskedVIFe0 == 0x15 {
//        // VIFE = E001 0101	Access Code Developer
//        return fmt.Sprintf("Access Code Developer")
//    } else if maskedVIFe0 == 0x16 {
//        // VIFE = E001 0110 Password
//        return fmt.Sprintf("Password")
//    } else if maskedVIFe0 == 0x17 {
//        // VIFE = E001 0111 Error flags
//        return fmt.Sprintf("Error flags")
//    } else if maskedVIFe0 == 0x18 {
//        // VIFE = E001 1000	Error mask
//        return fmt.Sprintf("Error mask")
//    } else if maskedVIFe0 == 0x19 {
//        // VIFE = E001 1001	Reserved
//        return fmt.Sprintf("Reserved")
//    } else if maskedVIFe0 == 0x1A {
//        // VIFE = E001 1010 Digital output (binary)
//        return fmt.Sprintf("Digital output (binary)")
//    } else if maskedVIFe0 == 0x1B {
//        // VIFE = E001 1011 Digital input (binary)
//        return fmt.Sprintf("Digital input (binary)")
//    } else if maskedVIFe0 == 0x1C {
//        // VIFE = E001 1100	Baudrate [Baud]
//        return fmt.Sprintf("Baudrate")
//    } else if maskedVIFe0 == 0x1D {
//        // VIFE = E001 1101	response delay time [bittimes]
//        return fmt.Sprintf("response delay time")
//    } else if maskedVIFe0 == 0x1E {
//        // VIFE = E001 1110	Retry
//        return fmt.Sprintf("Retry")
//    } else if maskedVIFe0 == 0x1F {
//        // VIFE = E001 1111	Reserved
//        return fmt.Sprintf("Reserved")
//    } else if maskedVIFe0 == 0x20 {
//        // VIFE = E010 0000	First storage # for cyclic storage
//        return fmt.Sprintf("First storage # for cyclic storage")
//    } else if maskedVIFe0 == 0x21 {
//        // VIFE = E010 0001	Last storage # for cyclic storage
//        return fmt.Sprintf("Last storage # for cyclic storage")
//    } else if maskedVIFe0 == 0x22 {
//        // VIFE = E010 0010	Size of storage block
//        return fmt.Sprintf("Size of storage block")
//    } else if maskedVIFe0 == 0x23 {
//        // VIFE = E010 0011	Reserved
//        return fmt.Sprintf("Reserved")
//    } else if maskedVIFe0 & 0x7C == 0x24 {
//        // VIFE = E010 01nn	Storage interval [sec(s)..day(s)]
//        n := maskedVIFe0 & 0x03
//        return fmt.Sprintf("Storage interval %s", unitDurationNN(int(n)))
//    } else if maskedVIFe0 == 0x28 {
//        // VIFE = E010 1000	Storage interval month(s)
//        return fmt.Sprintf("Storage interval month(s)")
//    } else if maskedVIFe0 == 0x29 {
//        // VIFE = E010 1001	Storage interval year(s)
//        return fmt.Sprintf("Storage interval year(s)")
//    } else if maskedVIFe0 == 0x2A {
//        // VIFE = E010 1010	Reserved
//        return fmt.Sprintf("Reserved")
//    } else if maskedVIFe0 == 0x2B {
//        // VIFE = E010 1011	Reserved
//        return fmt.Sprintf("Reserved")
//    } else if maskedVIFe0 & 0x7C == 0x2C {
//        // VIFE = E010 11nn	Duration since last readout [sec(s)..day(s)]
//        n := maskedVIFe0 & 0x03
//        return fmt.Sprintf("Duration since last readout %s", unitDurationNN(int(n)))
//    } else if maskedVIFe0 == 0x30 {
//        // VIFE = E011 0000	Start (date/time) of tariff
//        return fmt.Sprintf("Start (date/time) of tariff")
//    } else if maskedVIFe0 & 0x7C == 0x30 {
//        // VIFE = E011 00nn	Duration of tariff (nn=01 ..11: min to days)
//        n := maskedVIFe0 & 0x03
//        return fmt.Sprintf("Duration of tariff %s", unitDurationNN(int(n)))
//    } else if maskedVIFe0 & 0x7C == 0x34 {
//        // VIFE = E011 01nn	Period of tariff [sec(s) to day(s)]
//        n := maskedVIFe0 & 0x03
//        return fmt.Sprintf("Period of tariff %s", unitDurationNN(int(n)))
//    } else if maskedVIFe0 == 0x38 {
//        // VIFE = E011 1000	Period of tariff months(s)
//        return fmt.Sprintf("Period of tariff months(s)")
//    } else if maskedVIFe0 == 0x39 {
//        // VIFE = E011 1001	Period of tariff year(s)
//        return fmt.Sprintf("Period of tariff year(s)")
//    } else if maskedVIFe0 == 0x3A {
//        // VIFE = E011 1010	dimensionless / no VIF
//        return fmt.Sprintf("dimensionless / no VIF")
//    } else if maskedVIFe0 == 0x3B {
//        // VIFE = E011 1011	Reserved
//        return fmt.Sprintf("Reserved")
//    } else if maskedVIFe0 & 0x7C == 0x3C {
//        // VIFE = E011 11xx	Reserved
//        return fmt.Sprintf("Reserved")
//    } else if maskedVIFe0 & 0x70 == 0x40 {
//        // VIFE = E100 nnnn 10^(nnnn-9) V
//        n := maskedVIFe0 & 0x0F
//        return fmt.Sprintf("%s V", unitPrefix(int(n - 9)))
//    } else if maskedVIFe0 & 0x70 == 0x50 {
//        // VIFE = E101 nnnn 10nnnn-12 A
//        n := maskedVIFe0 & 0x0F
//        return fmt.Sprintf("%s A", unitPrefix(int(n - 12)))
//    } else if maskedVIFe0 == 0x60 {
//        // VIFE = E110 0000	Reset counter
//        return fmt.Sprintf("Reset counter")
//    } else if maskedVIFe0 == 0x61 {
//        // VIFE = E110 0001	Cumulation counter
//        return fmt.Sprintf("Cumulation counter")
//    } else if maskedVIFe0 == 0x62 {
//        // VIFE = E110 0010	Control signal
//        return fmt.Sprintf("Control signal")
//    } else if maskedVIFe0 == 0x63 {
//        // VIFE = E110 0011	Day of week
//        return fmt.Sprintf("Day of week")
//    } else if maskedVIFe0 == 0x64 {
//        // VIFE = E110 0100	Week number
//        return fmt.Sprintf("Week number")
//    } else if maskedVIFe0 == 0x65 {
//        // VIFE = E110 0101	Time point of day change
//        return fmt.Sprintf("Time point of day change")
//    } else if maskedVIFe0 == 0x66 {
//        // VIFE = E110 0110	State of parameter activation
//        return fmt.Sprintf("State of parameter activation")
//    } else if maskedVIFe0 == 0x67 {
//        // VIFE = E110 0111	Special supplier information
//        return fmt.Sprintf("Special supplier information")
//    } else if maskedVIFe0 & 0x7C == 0x68 {
//        // VIFE = E110 10pp	Duration since last cumulation [hour(s)..years(s)]Ž
//        n := maskedVIFe0 & 0x03
//        return fmt.Sprintf("Duration since last cumulation %s", unitDurationNN(int(n)))
//    } else if maskedVIFe0 & 0x7C == 0x6C {
//        // VIFE = E110 11pp	Operating time battery [hour(s)..years(s)]Ž
//        n := maskedVIFe0 & 0x03
//        return fmt.Sprintf("Operating time battery %s", unitDurationNN(int(n)))
//    } else if maskedVIFe0 == 0x70 {
//        // VIFE = E111 0000	Date and time of battery change
//        return fmt.Sprintf("Date and time of battery change")
//    } else if maskedVIFe0 & 0x70 == 0x70 {
//        // VIFE = E111 nnn Reserved
//        return fmt.Sprintf("Reserved VIF extension")
//    } else {
//        return fmt.Sprintf("Unrecognized VIF 0xFD extension: 0x%.2x", maskedVIFe0)
//    }
//}
