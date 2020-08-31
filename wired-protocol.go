package mbus

import (
    "fmt"
    "time"
)

// https://github.com/rscada/libmbus/blob/027f6fb6899b902bdd7d0b3230ecccc24f6bc6c3/mbus/mbus-protocol.h#L75
type MBusFrame struct {
    Start1 byte
    Length1 byte
    Length2 byte
    Start2 byte
    Control byte
    Address byte
    ControlInformation byte
    Checksum byte
    Stop byte

    Data []byte
    DataSize int

    // Holds the Data Records
    FrameData MbusDataFrame

    Type int

    Timestamp time.Time

    Next *MBusFrame
}

func NewWiredMBusFrame() *MBusFrame {
    return &MBusFrame{
        Data: make([]byte, FRAME_DATA_LENGTH),
        Next: nil,
    }
}


func (frame *MBusFrame) DecodeDeviceRecords() error {
    if DEBUG {
        fmt.Println("Decoding DR")
        for i := 0; i < frame.DataSize; i++ {
            fmt.Printf("%.2X ", frame.Data[i])
        }
    }

    direction := frame.Control & CONTROL_MASK_DIR
    fmt.Printf("\n0x%.2X\n", direction)
    // Check the direction of the frame ( Slave 2 Master )
    if direction == CONTROL_MASK_DIR_S2M {
        if frame.ControlInformation == CONTROL_INFO_ERROR_GENERAL {
            frame.FrameData.Type = DATA_TYPE_ERROR

            if frame.DataSize > 0 {
                frame.FrameData.Error = int(frame.Data[0])
            } else {
                frame.FrameData.Error = 0
            }
        } else if frame.ControlInformation == CONTROL_INFO_RESP_FIXED {
            if frame.DataSize == 0 {
                return fmt.Errorf("got not data")
            }

            frame.FrameData.Type = DATA_TYPE_FIXED
            // return mbus_data_fixed_parse(frame, &(data->data_fix));
        } else if frame.ControlInformation == CONTROL_INFO_RESP_VARIABLE {
            if frame.DataSize == 0 {
                return fmt.Errorf("got not data")
            }

            frame.FrameData.Type = DATA_TYPE_VARIABLE
            // return mbus_data_variable_parse(frame, &(data->data_var));
        } else {
            return fmt.Errorf("unknown control information 0x%.2X", frame.ControlInformation)
        }
    }

    return fmt.Errorf("wrong direction in frame (M2S -> master to slave)")
}
