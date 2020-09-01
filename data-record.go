package mbus

import (
	"bytes"
	"fmt"
)

func (dr *DataRecord) DecodeRecordFunction() string {
	switch dr.DIB.DIF & DATA_RECORD_DIF_MASK_FUNCTION {
	case 0x00:
		return "Instantaneous value"
	case 0x10:
		return "Maximum value"
	case 0x20:
		return "Minimum value"
	case 0x30:
		return "Value during error state"
	default:
		return "Unknown"
	}
}

// Return the storage number for a variable-length data record
func (dr *DataRecord) DecodeStorageNumber() int {
	bitIndex := 0
	result := 0

	result |= int(dr.DIB.DIF & DATA_RECORD_DIF_MASK_STORAGE_NO >> 6)
	bitIndex++

	for i := 0; i < dr.DIB.NDIFe; i++ {
		result |= int(dr.DIB.DIFe[i]&DATA_RECORD_DIFE_MASK_STORAGE_NO) << bitIndex
		bitIndex += 4
	}

	return result
}

// Return the tariff for a variable-length data record
func (dr *DataRecord) DecodeTariff() (int, error) {
	if dr.DIB.NDIFe < 0 {
		return 0, fmt.Errorf("could not decode tariff")
	}

	bitIndex := 0
	result := 0

	for i := 0; i < dr.DIB.NDIFe; i++ {
		result |= int(dr.DIB.DIFe[i]&DATA_RECORD_DIFE_MASK_TARIFF>>4) << bitIndex
		bitIndex += 2
	}

	return result, nil
}

func (dr *DataRecord) DecodeDevice() (int, error) {
	if dr.DIB.NDIFe == 0 {
		return 0, fmt.Errorf("could not decode device")
	}

	bitIndex := 0
	result := 0

	for i := 0; i < dr.DIB.NDIFe; i++ {
		result |= int(dr.DIB.DIFe[i]&DATA_RECORD_DIFE_MASK_DEVICE>>6) << bitIndex
		bitIndex++
	}

	return result, nil
}

func (dr *DataRecord) DecodeUnit() (VIF, error) {
	unit := dr.VIB.UnitLookup()

	return unit, nil
}

func (dr *DataRecord) DecodeValue() (string, interface{}, error) {
	buffer := bytes.Buffer{}
	var rawValue interface{}
	var intValue int
	//var floatValue float64
	//var timeValue time.Time

	var err error

	vif := dr.VIB.VIF & DIB_DIF_WITHOUT_EXTENSION
	//vife := dr.VIB.VIFe[0] & DIB_DIF_WITHOUT_EXTENSION

	unit, err := dr.DecodeUnit()
	if err != nil {
		return "", rawValue, err
	}

	switch dr.DIB.DIF & DATA_RECORD_DIF_MASK_DATA {
	case 0x00:
		return "", rawValue, nil

	// 1 byte integer (8 bit)
	case 0x01:
		if err := DecodeInt(dr.Data, 1, &intValue); err != nil {
			return "", rawValue, err
		}

		if DEBUG {
			fmt.Printf("DIF 0x%.2x was decoded using 1 byte integer\n", dr.DIB.DIF)
		}
		break

	// 2 byte (16 bit)
	case 0x02:
		// E110 1100  Time Point (date)
		if vif == 0x6C {
			//mbus_data_tm_decode(&time, record->data, 2);
			//snprintf(buff, sizeof(buff), "%04d-%02d-%02d",
			//    (time.tm_year + 1900),
			//    (time.tm_mon + 1),
			//    time.tm_mday);
		} else {
			// 2 byte integer
			if err := DecodeInt(dr.Data, 2, &intValue); err != nil {
				return "", rawValue, err
			}

			if DEBUG {
				fmt.Printf("DIF 0x%.2x was decoded using 2 byte integer\n", dr.DIB.DIF)
			}

			rawValue = float64(intValue) * unit.Exp
			_, err = fmt.Fprintf(&buffer, "%.2f", rawValue)
		}
		break

	// 3 byte integer (24 bit)
	case 0x03:
		if err := DecodeInt(dr.Data, 3, &intValue); err != nil {
			return "", rawValue, err
		}

		if DEBUG {
			fmt.Printf("DIF 0x%.2x was decoded using 3 byte integer\n", dr.DIB.DIF)
		}

		_, err = fmt.Fprintf(&buffer, "%d", intValue)
		break

	// 4 byte (32 bit)
	case 0x04:
		// E110 1101  Time Point (date/time)
		// E011 0000  Start (date/time) of tariff
		// E111 0000  Date and time of battery change
		if vif == 0x6D {
			// || (dr.VIB.VIF == 0xFD && vife == 0x30) || (dr.VIB.VIF == 0xFD && vife == 0x70) {

			//mbus_data_tm_decode(&time, record->data, 4);
			//snprintf(buff, sizeof(buff), "%04d-%02d-%02dT%02d:%02d:%02d",
			//    (time.tm_year + 1900),
			//    (time.tm_mon + 1),
			//    time.tm_mday,
			//    time.tm_hour,
			//    time.tm_min,
			//    time.tm_sec);

			// 4 byte integer
		} else {
			if err := DecodeInt(dr.Data, 4, &intValue); err != nil {
				return "", rawValue, err
			}

			if DEBUG {
				fmt.Printf("DIF 0x%.2x was decoded using 4 byte integer\n", dr.DIB.DIF)
			}

			_, err = fmt.Fprintf(&buffer, "%d", intValue)
			break
		}
		// 4 Byte Real (32 bit)
		//case 0x05:
		//    float_val = mbus_data_float_decode(record->data);
		//
		//    snprintf(buff, sizeof(buff), "%f", float_val);
		//
		//    if (debug)
		//    printf("%s: DIF 0x%.2x was decoded using 4 byte Real\n", __PRETTY_FUNCTION__, record->drh.dib.dif);
		//
		//    break;
		//

	case 0x09, // 2 digit BCD (8 bit)
		0x0A, // 4 digit BCD (16 bit)
		0x0B, // 6 digit BCD (24 bit)
		0x0C, // 8 digit BCD (32 bit)
		0x0E: // 12 digit BCD (48 bit)
		if err := DecodeBCDHEX(dr.Data, dr.DataSize, &intValue); err != nil {
			return "", rawValue, err
		}

		if DEBUG {
			fmt.Printf("DIF 0x%.2x was decoded using %d digit BCD\n", dr.DIB.DIF, dr.DataSize*2)
		}

		_, err = fmt.Fprintf(&buffer, "%X", intValue)
		break
	default:
		err = fmt.Errorf("unkown DIF (0x%.2X)", dr.DIB.DIF)
		break
	}

	// Check error
	if err != nil {
		return "", rawValue, err
	}

	return buffer.String(), rawValue, nil
}
