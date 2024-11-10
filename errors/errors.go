package errors

import "encoding/binary"


type ErrorCode interface {
	Code() uint16
	Message() string
}

type APIError struct {
	errorCode uint16
	message string
}

func (e *APIError) Code() uint16 {
	return e.errorCode
}

func (e *APIError) Message() string {
	return e.message
}

func NewAPIError(code uint16, message string) *APIError {
	return &APIError {
		errorCode: code,
		message: message,
	}
}

func ErrorCodeToBytes(err ErrorCode) []byte {
	response := make([]byte, 2)
	binary.BigEndian.PutUint16(response, err.Code())
	return response
}
