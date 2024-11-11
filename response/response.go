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

// https://forum.codecrafters.io/t/question-about-handle-apiversions-requests-stage/1743/4
//
// ApiVersions Response (Version: CodeCrafters) =>
//	error_code num_of_api_keys [api_keys] throttle_time_ms TAG_BUFFER
// error_code => INT16
// num_of_api_keys => INT8
// api_keys => api_key min_version max_version
//   api_key => INT16
// 	 min_version => INT16
//   max_version => INT16
// _tagged_fields
// throttle_time_ms => INT32
// _tagged_fields
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