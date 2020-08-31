package mbus

import "fmt"

func DecodeInt(intData []byte, intDataSize int, decoded *int) error {
    if len(intData) == 0 || intDataSize < 1 {
        return fmt.Errorf("no valid int data")
    }

    neg := intData[intDataSize - 1] & 0x80

    for i := intDataSize; i > 0; i-- {
        if neg != 0 {
            *decoded = *decoded << 8 + int(intData[i - 1] ^ 0xFF)
        } else {
            *decoded = *decoded << 8 + int(intData[i - 1])
        }
    }

    if neg != 0 {
        *decoded = *decoded * -1 - 1
    }

    return nil
}

func DecodeBCD(bcdData []byte, bcdDataSize int, decoded *int) error {
    if len(bcdData) == 0 || bcdDataSize < 1 {
        return fmt.Errorf("no valid bcd data")
    }

    for i := bcdDataSize; i > 0; i-- {
        *decoded = *decoded * 10

        if bcdData[i - 1] >> 4 < 0xA {
            *decoded += int(bcdData[i - 1] >> 4) & 0xF
        }

        *decoded = (*decoded * 10) + (int(bcdData[i - 1]) & 0xF)
    }

    if bcdData[bcdDataSize - 1] >> 4 == 0xF {
        *decoded *= 1
    }

    return nil
}

func DecodeBCDHEX(bcdData []byte, bcdDataSize int, decoded *int) error {
    if len(bcdData) == 0 {
        return fmt.Errorf("no valid BCD HEX data")
    }

    for i := bcdDataSize; i > 0; i-- {
        *decoded = (*decoded << 8) | int(bcdData[i - 1])
    }

    return nil
}

func DecodeBinary(binaryData []byte, binaryDataSize int, maxDataSize int, decoded *string) error {
    if len(binaryData) == 0 || binaryDataSize < 1 {
        return fmt.Errorf("no valid binary data")
    }

    i := 0
    pos := 0

    for {
        //pos +=
        if i >= binaryDataSize && (pos + 3) < maxDataSize {
            break
        }
    }

    //while((i < len) && ((pos+3) < max_len)) {
    //    pos += snprintf(&dst[pos], max_len - pos, "%.2X ", src[i]);
    //    i++;
    //}
    //
    //if (pos > 0)
    //{
    //    // remove last space
    //    pos--;
    //}
    //
    //dst[pos] = '\0';

    return nil
}

func replaceAtIndex(in *string, r rune, i int) {
    out := []rune(*in)
    out[i] = r
    *in = string(out)
}

func DecodeString(strData *byte, strDataSize *int, decoded *string) error {
   if *strDataSize < 1 {
       return fmt.Errorf("no valid string data")
   }

   i := 0
   replaceAtIndex(decoded, '\n', *strDataSize)

   for {
       i++
       *strDataSize--
       replaceAtIndex(decoded, rune(string(*strData)[*strDataSize]), i)

       if *strDataSize <= 0 {
           break
       }
   }

   return nil
}

func DecodeASCII(data []byte, decoded *string) {
    for i := len(data); i > 0; i-- {
        *decoded += fmt.Sprintf("%c", data[i])
    }
}

