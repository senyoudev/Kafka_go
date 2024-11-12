package response

import (
	"encoding/binary"
	"fmt"
	"net"
	"github.com/codecrafters-io/kafka-starter-go/errors"
)

type APIKey struct {
    Key         uint16
    MinVersion  uint16
    MaxVersion  uint16
    TaggedField byte
}

var apiKeys = []APIKey{
    {
        Key:         18, // APIVersions
        MinVersion:  0,
        MaxVersion:  4,
		TaggedField: 0,
    },
    {
        Key:         75, // DescribeTopicPartitions
        MinVersion:  0,
        MaxVersion:  0,
		TaggedField: 0,
    },
}

func CreateDescribeTopicPartitionsResponse(correlationId uint32, topicName string) []byte {
    totalSize := 4 + // correlation ID
                  1 + // first tag buffer
                  4 + // throttle time
                  1 + // array length (varint)
                  2 + // error code
                  2 + len(topicName) + // topic name length + content
                  16 + // topic ID (UUID)
                  1 + // is internal
                  1 + // partitions array length
                  4 + // topic authorized operations
                  1 + // topic tag buffer
                  1 + // next cursor (nullable)
                  1   // final tag buffer
	
	response := make([]byte, totalSize)
	offset := 0

	// Write correlationId
	binary.BigEndian.PutUint32(response[offset:], correlationId)
    offset += 4

	// First tag buffer (empty)
    response[offset] = 0x00
    offset++

    // Throttle time (0)
    binary.BigEndian.PutUint32(response[offset:], 0)
    offset += 4

	// Array length (length + 1 as varint)
    response[offset] = 0x02 // Array length of 1 encoded as varint
    offset++

    // Error code (UNKNOWN_TOPIC = 3)
    binary.BigEndian.PutUint16(response[offset:], errors.UNKNOWN_TOPIC_OR_PARTITION.Code())
    offset += 2

	topicLength := len(topicName) + 1
    response[offset] = byte(topicLength)
    offset++

    // Topic name content
    copy(response[offset:], topicName)
    offset += len(topicName)

    // Topic ID (16 zero bytes)
    zeroUUID := make([]byte, 16)
    copy(response[offset:], zeroUUID)
    offset += 16

    // Is Internal (false)
    response[offset] = 0x00
    offset++

	// Partitions array (empty)
    response[offset] = 0x01 // Empty array encoded as varint
    offset++

    // Topic authorized operations
    binary.BigEndian.PutUint32(response[offset:], 0x00000df8)
    offset += 4

	    // Topic tag buffer
    response[offset] = 0x00
    offset++

    // Next cursor (null)
    response[offset] = 0xff
    offset++

    // Final tag buffer
    response[offset] = 0x00

    return response
}

// https://forum.codecrafters.io/t/question-about-handle-apiversions-requests-stage/1743/4
func CreateAPIVersionsResponse(correlationID uint32) []byte {
	totalSize := 4 + // Correlation ID
			2 + // Error Code
			1 + // Number of API Keys
			(len(apiKeys) * 7) + // Each API Key entry (2+2+2+1 bytes)
			4 + // Throttle time
			1   // Final tagged fields
	response := make([]byte, totalSize)
	offset := 0
	
	// CorrelationId
	binary.BigEndian.PutUint32(response[offset:], correlationID)
	offset += 4
	// error code(no code = 0)
	binary.BigEndian.PutUint16(response[offset:], 0)
	offset += 2
	
	response[offset] = byte(len(apiKeys) + 1)
	offset++

	for _, api := range apiKeys {
        binary.BigEndian.PutUint16(response[offset:], api.Key)
        offset += 2
        binary.BigEndian.PutUint16(response[offset:], api.MinVersion)
        offset += 2
        binary.BigEndian.PutUint16(response[offset:], api.MaxVersion)
        offset += 2
		response[offset] = api.TaggedField
		offset++
    }
	
	binary.BigEndian.PutUint32(response[offset:], 0) // throttle time
	offset += 4
	response[offset] = 0 // final tagged fields
	
	return response
}

func CreateErrorResponse(correlationID uint32, err errors.ErrorCode) []byte {
	response := make([]byte, 6)
	binary.BigEndian.PutUint32(response[0:4], correlationID)
	
	if err != nil {
		errorBytes := errors.ErrorCodeToBytes(err)
		copy(response[4:], errorBytes)
		fmt.Printf("Error response: %s (code: %d)\n", err.Message(), err.Code())
	} else {
		binary.BigEndian.PutUint16(response[4:], 0)
	}
	
	return response
}

func SendResponse(conn net.Conn, response []byte) {
	length := make([]byte, 4)
	binary.BigEndian.PutUint32(length, uint32(len(response)))
	
	conn.Write(length)
	conn.Write(response)
}

func SendErrorResponse(conn net.Conn, correlationID uint32, err errors.ErrorCode) {
	response := CreateErrorResponse(correlationID, err)
	SendResponse(conn, response)
}