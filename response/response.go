package response

import (
	"encoding/binary"
	"fmt"
	"net"
	"github.com/codecrafters-io/kafka-starter-go/errors"
)

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
	response := make([]byte, 19)
	
	binary.BigEndian.PutUint32(response[0:4], correlationID)
	binary.BigEndian.PutUint16(response[4:6], 0)
	
	response[6] = 2
	
	binary.BigEndian.PutUint16(response[7:9], 18)
	binary.BigEndian.PutUint16(response[9:11], 0)
	binary.BigEndian.PutUint16(response[11:13], 4)
	
	response[13] = 0
	binary.BigEndian.PutUint32(response[14:18], 0)
	response[18] = 0
	
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