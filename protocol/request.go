package protocol

import "encoding/binary"

const (
	MaxPacketSize     = 1024
	APIVersionsKey    = 18
	MaxSupportedVersion = 4
)

type Request struct {
	APIKey        uint16
	APIVersion    int16
	CorrelationID uint32
}

func ParseRequest(buff []byte) Request {
	return Request{
		APIKey:        binary.BigEndian.Uint16(buff[4:6]),
		APIVersion:    int16(binary.BigEndian.Uint16(buff[6:8])),
		CorrelationID: binary.BigEndian.Uint32(buff[8:12]),
	}
}

func IsVersionSupported(version int16) bool {
	return version >= 0 && version <= MaxSupportedVersion
}