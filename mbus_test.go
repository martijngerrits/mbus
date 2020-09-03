package mbus

import (
	"context"
	"fmt"
	"testing"
)

var testFrame = []byte{
	// After RadioCrafts module reset + Memory Configuration -> CONFIG_INTERFACE => 0x0C ( Add start/stop byte and CRC )
	//0x68, // Start
	//
	////************************************
	//// Link Layer
	////************************************
	//0x50, // Length 0x50 = 80 bytes -> Header{CI + M (2) + ID (4) + Version + Device Type + CI} = 10 + Data length (68) + CRC (2)
	//0x44, // Control
	//0x33, 0x30, // Manufacturer ( MSB first )
	//0x53, 0x56, 0x02, 0x00, // Device Serial Number ( MSB first ) => LAN-WMBUS-E-VOC
	//0x01, // Protocol Version
	//0x2B, // Device Type
	//
	////************************************
	//// Network Layer
	////************************************
	//0x7A, // Control Information
	//0x51, // Access no. ( transmission counter )
	//0x00, // Status byte ( 0x00 = No errors, 0x04 = Low battery, 0x08 = Permanent error/Sabotage enclosure, 0x40 = Sabotage alarm )
	//0x40, // 0x40 = 64 => No. Encrypted Blocks
	//0x25, // Encryption mode 5 + Synchronized
	//
	////************************************
	//// Data Blocks (64 bytes) ( Encrypted )
	////************************************
	//0x33, 0x57, 0x08, 0x48, 0xEE, 0x40, 0x44, 0x68,
	//0x60, 0x96, 0xF8, 0x5D, 0xF1, 0x6C, 0x2F, 0x48,
	//0x09, 0xCA, 0xB7, 0x8E, 0x54, 0x2C, 0x9B, 0x7C,
	//0x3C, 0x36, 0x7C, 0xD0, 0x1A, 0x71, 0xC7, 0x06,
	//0x98, 0xB3, 0x6D, 0xC4, 0x04, 0x82, 0x36, 0x6C,
	//0x44, 0x41, 0xA6, 0x0C, 0xB1, 0x88, 0x5E, 0xD5,
	//0xE8, 0x5D, 0x89, 0x74, 0x78, 0x05, 0xA9, 0x72,
	//0x84, 0xA1, 0x2C, 0xC5, 0x6C, 0x6D, 0xB5, 0x56,
	//
	//// Decrypted Data Records:
	////  Index | Bytes | Data bytes        | Type
	//// _______|_______|___________________|_______________________________________
	////    0   |   2   | 2F 2F             | AES Verification bytes
	////    2   |   4   | 02 65 9D 0B       | DR1, 16 bit integer
	////    6   |   4   | 42 65 9A 0B       | DR2, 16 bit integer + storage 1
	////   10   |   5   | 82 01 65 52 0B    | DR3, 16 bit integer
	////   15   |   5   | 02 FB 1A 2F 02    | DR4, 16 bit integer
	////   20   |   5   | 42 FB 1A 2F 02    | DR5, 16 bit integer + storage 1
	////   25   |   6   | 82 01 FB 1A 4B 02 | DR6, 16 bit integer
	////   31   |   5   | 02 FD 3A 4E 01    | DR7, 16 bit integer
	////   36   |   5   | 42 FD 3A 4F 01    | DR8, 16 bit integer + storage 1
	////   41   |   6   | 82 01 FD 3A 43 01 | DR9, 16 bit integer
	////   47   |   5   | 02 FD 0F 04 00    | DR10, 16 bit integer
	////   51   |  12   | 2F 2F 2F 2F 2F 2F | AES Encryption Filler bytes
	////        |       | 2F 2F 2F 2F 2F 2F |
	//
	//0xC2, 0xF9, // CRC
	//
	//0x16, // Stop

	// LAN-WMBUS-G2-OOP
	0x68,
	0x20,
	0x44,
	0x33, 0x30,
	0x54, 0x56, 0x02, 0x00,
	0x0B,
	0x02,
	0x7A, 0x87, 0x00, 0x10, 0x25, 0xD6, 0xF4,
	0x2D, 0xD2, 0x66, 0x0C, 0x65, 0x6E, 0xEB, 0x46,
	0x3D, 0xD8, 0xC2, 0x64, 0xC3, 0x0E, 0xD7, 0xCD,
	0x16,

	//0x68,
	//0x47, 0x44,
	//0x93, 0x15,
	//0x68, 0x77, 0x02, 0x35, 0x42, 0x03, 0x7A, 0x7B, 0x00, 0x30, 0xAF, 0x7C, 0x35, 0xB5, 0x94, 0x60, 0xC1, 0x1C, 0xD2, 0xBD, 0xC0, 0x5F, 0xBD, 0xE6, 0xC4, 0xD3, 0x7B, 0x93, 0x18, 0x9B, 0xFD, 0x65, 0x20, 0x7B, 0x5B, 0x78, 0xEA, 0xE4, 0x24, 0xDC, 0x8B, 0xA7, 0x59, 0x75, 0x80, 0xB6, 0x8C, 0x32, 0xD1, 0x84, 0x4F, 0x67, 0x57, 0x7E, 0x9E, 0x4A, 0x35, 0xE6, 0xEB, 0x04, 0xFD, 0x08, 0x59, 0x82, 0x00, 0x00, 0xED, 0x2D, 0x16,

	//0x68,
	//0x0B,
	//0x08, // Control = CONTROL_MASK_RSP_UD
	//0x06, 0x0A, // Manufacturer
	//0x00, 0xB0, 0x00, 0xB0, // Id
	//0x01, // Version
	//0x00, // Device
	//0xC0, // Control information
	//0x00,
	//0x16,
}

var devices = []Device{
	{
		SerialNumber: "25653",
		AESKey: []byte{
			0x3D, 0x19, 0x76, 0x69, 0xB5, 0x3C, 0x3E, 0xA9,
			0xA2, 0x61, 0x5B, 0x28, 0x5C, 0xA1, 0x72, 0x1A,
		},
	},
	{
		SerialNumber: "25654",
		AESKey: []byte{
			0x27, 0xF9, 0x27, 0x62, 0xF6, 0x6A, 0x41, 0xCB,
			0x26, 0x71, 0x31, 0xDB, 0x09, 0x12, 0x22, 0x46,
		},
	},
}

// Expected results based on the test frame above
var expectedManufacturer = "LAS"
var expectedProductName = "LAN-WMBUS-E-VOC"

func FindAESKeyForSerialNumber(serialNumber string) ([]byte, error) {
	var key = make([]byte, 0, 16)
	for _, device := range devices {
		if device.SerialNumber == serialNumber {
			key = device.AESKey
		}
	}

	if len(key) == 0 {
		return nil, fmt.Errorf("could not find AES key for serial number: %s\n", serialNumber)
	}

	return key, nil
}

func TestNewSerialClient(t *testing.T) {
	client, err := NewSerialClient("/tmp/ttyConBridge", SerialConfig{
		//BaudRate: 19200,
		//DataBits: 8,
		Baud:     19200,
		Size:     8,
		StopBits: 1,
		Parity:   'N',
	})

	if err != nil {
		t.Fatal(err)
	}

	DEBUG = true

	defer client.Close()

	fmt.Println("Serial connection made")

	// Create context with cancel
	ctx, cancel := context.WithCancel(context.Background())
	stream := client.Stream(ctx)

	i := 0
	for frame := range stream {
		serialNumber, err := frame.DecodeSerialNumber()
		if err != nil {
			fmt.Println("something went wrong while decoding serial number, skipping frame")
			continue
		}

		key, err := FindAESKeyForSerialNumber(serialNumber)
		if err != nil {
			t.Fatal(err)
		}

		// Decrypt the data blocks
		if err := frame.DecryptData(key); err != nil {
			t.Fatal(err)
		}

		if err := frame.DataParse(); err != nil {
			t.Fatal(err)
		}

		decodedFrame, err := frame.DecodeFrame()
		if err != nil {
			t.Fatal(err)
		}

		t.Logf("Serial number: %s\n", decodedFrame.SerialNumber)
		t.Logf("Manufacturer: %s\n", decodedFrame.Manufacturer)
		t.Logf("Product name: %s\n", decodedFrame.ProductName)
		t.Logf("Device type: %s\n", decodedFrame.DeviceType)
		t.Logf("Status: %.2X (%s)\n", decodedFrame.Status, decodedFrame.ReadableStatus)
		t.Logf("Version: %d\n", decodedFrame.Version)
		t.Logf("ParsedAt: %s", decodedFrame.ParsedAt)
		for i := 0; i < len(decodedFrame.DataRecords); i++ {
			r := decodedFrame.DataRecords[i]

			t.Logf("DR %d :: %s: %s %s (%d)\n", i+1, r.Function, r.Value, r.Unit, r.StorageNumber)
		}

		i++

		if i == 3333333 {
			cancel()
			break
		}
	}
}

func TestDecodeManufacturer(t *testing.T) {
	frame := NewWirelessMBusFrame()

	if _, err := ParseWirelessMBusData(frame, &testFrame, len(testFrame)); err != nil {
		t.Fatal(err)
	}

	manufacturer, err := frame.DecodeManufacturer()
	if err != nil {
		t.Fatal(err)
	}

	if manufacturer != expectedManufacturer {
		t.Fatal(
			fmt.Errorf(
				"decoded manufacturer does not match expected: '%s', got: %s",
				expectedManufacturer,
				manufacturer,
			),
		)
	}
}

func TestProductName(t *testing.T) {
	frame := NewWirelessMBusFrame()

	if _, err := ParseWirelessMBusData(frame, &testFrame, len(testFrame)); err != nil {
		t.Fatal(err)
	}

	productName, err := frame.DecodeProductName()
	if err != nil {
		t.Fatal(err)
	}

	if productName != expectedProductName {
		t.Fatal(
			fmt.Errorf(
				"product name does not match expected: '%s', got: %s",
				expectedProductName,
				productName,
			),
		)
	}
}

func TestDecodeFrame(t *testing.T) {
	frame := NewWirelessMBusFrame()

	// Parse the bytes
	if _, err := ParseWirelessMBusData(frame, &testFrame, len(testFrame)); err != nil {
		t.Fatal(err)
	}

	serialNumber, err := frame.DecodeSerialNumber()
	if err != nil {
		t.Fatal(err)
	}

	if frame.HasEncryptionMode() {
		key, err := FindAESKeyForSerialNumber(serialNumber)
		if err != nil {
			t.Fatal(err)
		}

		// Decrypt the data blocks
		if err := frame.DecryptData(key); err != nil {
			t.Fatal(err)
		}
	}

	if err := frame.DataParse(); err != nil {
		t.Fatal(err)
	}

	decodedFrame, err := frame.DecodeFrame()
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Serial number: %s\n", decodedFrame.SerialNumber)
	t.Logf("Manufacturer: %s\n", decodedFrame.Manufacturer)
	t.Logf("Product name: %s\n", decodedFrame.ProductName)
	t.Logf("Device type: %s\n", decodedFrame.DeviceType)
	t.Logf("Status: %.2X (%s)\n", decodedFrame.Status, decodedFrame.ReadableStatus)
	t.Logf("Version: %d\n", decodedFrame.Version)
	t.Logf("ParsedAt: %s", decodedFrame.ParsedAt)
	for i := 0; i < len(decodedFrame.DataRecords); i++ {
		r := decodedFrame.DataRecords[i]

		t.Logf("DR %d :: %s: %s %s (%d)\n", i+1, r.Function, r.Value, r.Unit, r.StorageNumber)
	}
}
