package mbus

import (
    "context"
    "fmt"
    "github.com/tarm/serial"
)

type SerialConfig serial.Config

type MbusSerialHandle struct {
    MbusHandle
    Fd *serial.Port
}

func NewSerialClient(device string, config SerialConfig) (Handle, error) {
    client := &MbusSerialHandle{
        Fd: nil,
        MbusHandle: MbusHandle{
            MaxDataRetry: 3,
            MaxSearchRetry: 3,
            IsSerial: true,
        },
    }

    if err := client.Open(device, config); err != nil {
        return nil, err
    }

    return client, nil
}

func (handle *MbusSerialHandle) Open(device string, config interface{}) error {
    serialConfig := config.(SerialConfig)

    port, err := serial.OpenPort(&serial.Config{
       Name: device,
       Baud: serialConfig.Baud,
       Size: serialConfig.Size,
       StopBits: serialConfig.StopBits,
       Parity: serialConfig.Parity,
       ReadTimeout: serialConfig.ReadTimeout,
    })
    if err != nil {
        return err
    }

    handle.Fd = port
    return nil
}

func (handle *MbusSerialHandle) Stream(ctx context.Context) chan Frame {
    stream := make(chan Frame, 1024)

    go func() {
        for {
            select {
            case <-ctx.Done():
                close(stream)
                return
            default:
                frame, err := handle.ReceiveFrame()
                if err != nil {
                    fmt.Printf("Got error while receiving frame: %s\n", err)
                } else {
                    stream<-frame
                }
            }
        }
    }()

    return stream
}

func (handle *MbusSerialHandle) Close() error {
    if err := handle.Fd.Close(); err != nil {
        return err
    }

    return nil
}

func (handle *MbusSerialHandle) Send(frame Frame) error {
    //data, _ := frame.Encode()
    //handle.Fd.Write(data)
    return nil
}

func (handle *MbusSerialHandle) ReceiveFrame() (Frame, error) {
    buffer := make([]byte, PACKET_BUFF_SIZE)
    frame := NewWirelessMBusFrame()

    // Clear the buffer for the GC
    defer func() {
        buffer = nil
    }()

    length := 0
    remaining := 1
    timeouts := 0

    for {
        if length + remaining > PACKET_BUFF_SIZE {
            return nil, fmt.Errorf("out of bounds")
        }

        // Since Fd.Read() will read all bytes until the given buffer is full..
        // we use the tmpBuffer to force Fd.Read() to only retrieve len(tmpBuffer) amount bytes
        // tmpBuffer is a dynamic sized buffer on each loop
        // This way we can check each byte on its own, and only start collecting data when there is a frame start byte.

        // We could resize the buffer (buffer = buffer[:remaining] to fix this issue,
        // but then every loop the buffer would be overwritten from index 0,
        // since Fd.Read() will write into the buffer from the start (index 0).
        tmpBuffer := make([]byte, remaining)

        if DEBUG {
            fmt.Println("Waiting for data from serial device...")
        }
        nread, err := handle.Fd.Read(tmpBuffer)
        if err != nil {
            return nil, err
        }

        if nread == 0 {
            timeouts++

            if timeouts >= 3 {
                return nil, fmt.Errorf("timeout")
            }

            // Nothing else to do here
            continue
        }

        // Copy over the data from the temporary remaining buffer to the actual buffer for the M-Bus frame
        for i := 0; i < len(tmpBuffer); i++ {
            buffer[length + i] = tmpBuffer[i]
        }

        // Clear the temporary buffer
        tmpBuffer = nil

        length += nread

        parseResult, err := ParseWirelessMBusData(frame, &buffer, length)
        if err != nil {
            return nil, err
        }

        remaining = parseResult.Remaining

        // We got some useless bytes
        if !parseResult.GotFrame {
            length -= nread
        }

        if remaining == 0 {
            return frame, nil
        }
    }
}
