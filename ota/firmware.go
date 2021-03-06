package ota

import (
	"fmt"
	"github.com/kierdavis/ihex-go"
	"os"
)

const firmwareBlockSize uint16 = 16

// Firmware - The MySysBootloader Firmware Calculations
type Firmware struct {
	Blocks uint16
	Crc    uint16
	data   []byte
}

// NewFirmware - Loads a filename; computes block count and crc
func NewFirmware(filename string) *Firmware {
	file, err := os.Open(filename)
	if err != nil {
		return &Firmware{}
	}
	defer file.Close()

	blocks := uint16(0)
	crc := uint16(0)
	data := []byte{}
	start := uint16(0)
	end := uint16(0)

	scanner := ihex.NewDecoder(file)
	for scanner.Scan() {
		record := scanner.Record()
		if record.Type != ihex.Data {
			continue
		}

		if start == 0 && end == 0 {
			start = record.Address
			end = record.Address
		}

		for record.Address > end {
			data = append(data, 255)
			end++
		}

		data = append(data, record.Data...)

		end += uint16(len(record.Data))
	}

	pad := end % 128
	for i := uint16(0); i < 128-pad; i++ {
		data = append(data, 255)
		end++
	}

	blocks = uint16(end-start) / firmwareBlockSize
	crc = 0xFFFF
	for i := 0; i < len(data); i++ {
		crc = (crc ^ uint16(data[i]&0xFF))
		for j := 0; j < 8; j++ {
			a001 := (crc & 1) > 0
			crc = (crc >> 1)
			if a001 {
				crc = (crc ^ 0xA001)
			}
		}
	}

	return &Firmware{
		Blocks: blocks,
		Crc:    crc,
		data:   data,
	}
}

// Data - Gets a specific block from the firmware data
func (f Firmware) Data(block uint16) ([]byte, error) {
	fromBlock := block * firmwareBlockSize
	toBlock := fromBlock + firmwareBlockSize
	if dataLen := uint16(len(f.data)); dataLen < toBlock {
		return []byte{}, fmt.Errorf("Block %d cannot be found in the firmware data.", block)
	}

	return f.data[fromBlock:toBlock], nil
}
