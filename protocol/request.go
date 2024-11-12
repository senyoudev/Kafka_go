package protocol

import "encoding/binary"

const (
	MaxPacketSize     = 1024
	APIVersionsKey    = 18
	MaxSupportedVersion = 4
	UnknownTopicOrPartition = 3
	DescribeTopicPartitionsKey  = 75
)

type Request struct {
	APIKey        uint16
	APIVersion    int16
	CorrelationID uint32
}

type DescribeTopicPartitionsRequest struct {
    Request            // Embed the base Request struct
    ClientID          string
    Topics            []string
    PartitionLimit    uint32
    Cursor            byte   // 0xff for null
}

func ParseDescribeTopicPartitionsRequest(buff []byte) DescribeTopicPartitionsRequest {
    // First parse the basic request info using existing ParseRequest
    baseRequest := ParseRequest(buff)
    
    var req DescribeTopicPartitionsRequest
    req.APIKey = baseRequest.APIKey
    req.APIVersion = baseRequest.APIVersion
    req.CorrelationID = baseRequest.CorrelationID
    
    // Start parsing from after the basic header (after correlation ID)
    offset := 12 // 4 (message size) + 2 (API key) + 2 (version) + 4 (correlation ID)
    
    // Parse Client ID length and content
    clientIDLength := binary.BigEndian.Uint16(buff[offset:])
    offset += 2
    req.ClientID = string(buff[offset:offset+int(clientIDLength)])
    offset += int(clientIDLength)
    
    // Skip tag buffer
    offset++
    
    // Parse topics array
    topicsLength := int(buff[offset]) - 1  // Subtract 1 as per compact array encoding
    offset++
    
    req.Topics = make([]string, topicsLength)
    for i := 0; i < topicsLength; i++ {
        // Parse topic name length and content
        topicLength := int(buff[offset]) - 1  // Subtract 1 for compact string
        offset++
        req.Topics[i] = string(buff[offset:offset+topicLength])
        offset += topicLength
        // Skip topic tag buffer
        offset++
    }
    
    // Parse partition limit
    req.PartitionLimit = binary.BigEndian.Uint32(buff[offset:])
    offset += 4
    
    // Parse cursor
    req.Cursor = buff[offset]
    
    return req
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