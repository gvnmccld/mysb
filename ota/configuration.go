package ota

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"strings"
)

// Configuration - The MySysBootloader Firmware Config Request
type Configuration struct {
	Type    uint16
	Version uint16
	Blocks  uint16
	Crc     uint16
}

// NewConfiguration - Loads a string; computes type/version/blocks/crc
func NewConfiguration(payload string) *Configuration {
	t := Configuration{}
	b, err := hex.DecodeString(payload)
	if err != nil {
		return &t
	}

	r := bytes.NewReader(b)
	binary.Read(r, binary.LittleEndian, &t)

	return &t
}

func (t *Configuration) String() string {
	w := new(bytes.Buffer)
	binary.Write(w, binary.LittleEndian, t)
	return strings.ToUpper(hex.EncodeToString(w.Bytes()))
}
