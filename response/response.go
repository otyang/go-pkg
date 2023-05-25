package response

import (
	"bytes"
	"encoding/gob"
	"net/http"
)

var UseDefaultCode = true

// Response this struct represents the response structure.
type Response struct {
	StatusCode int     `json:"-"`
	Success    bool    `json:"success"`
	Message    string  `json:"message"`
	ErrorCode  *string `json:"errorCode,omitempty"`
	Data       any     `json:"data,omitempty"`
}

// Error, lets satisfy golang's error interface.
// so we can return 'response' as an error type.
func (a Response) Error() string {
	return a.Message
}

// OverideDefaultMsg helps overide or set a message useful
// in shorthand variables dfinition
func (a *Response) OverideDefaultMsg(msg string) *Response {
	a.Message = msg
	return a
}

// New creates an response struct
func New(statusCode int, success bool, msg string, errorCode *string, data any) *Response {
	if !UseDefaultCode {
		errorCode = nil
	}
	return &Response{
		StatusCode: statusCode,
		Success:    success,
		Message:    msg,
		ErrorCode:  errorCode,
		Data:       data,
	}
}

// NewError creates an error response
func NewError(statusCode int, msg string, errorCode string) *Response {
	return New(statusCode, false, msg, &errorCode, nil)
}

// NewSuccess creates a success response
func NewSuccess(statusCode int, msg string, data any) *Response {
	if msg == "" {
		msg = "Request was successful"
	}
	return New(statusCode, true, msg, nil, data)
}

// OK creates a success response with (HTTP 200) code
func Ok(msg string, data any) *Response {
	if msg == "" {
		msg = "Request was successful"
	}
	return NewSuccess(http.StatusOK, msg, data)
}

// Created creates a success response with (HTTP 201) code
func Created(msg string, data any) *Response {
	if msg == "" {
		msg = "Request was successful"
	}
	return NewSuccess(http.StatusCreated, msg, data)
}

// BadRequest creates a error response with (HTTP 400) code
func BadRequest(msg string, errorCode string) *Response {
	if msg == "" {
		msg = "Your request is in a bad format"
	}

	return NewError(
		http.StatusCreated,
		msg,
		setErrCode(errorCode, "bad_request"),
	)
}

// Unauthorized creates a error response with (HTTP 401) code
func Unauthorized(msg string, errorCode string) *Response {
	if msg == "" {
		msg = "You are not authenticated to perform the requested action"
	}

	return NewError(
		http.StatusUnauthorized,
		msg,
		setErrCode(errorCode, "un_authorized"),
	)
}

// Forbidden creates a error response with (HTTP 403) code
func Forbidden(msg string, errorCode string) *Response {
	if msg == "" {
		msg = "You are not authorized to perform the requested action"
	}

	return NewError(
		http.StatusForbidden,
		msg,
		setErrCode(errorCode, "forbidden"),
	)
}

// NotFound creates a error response with (HTTP 404) code
func NotFound(msg string, errorCode string) *Response {
	if msg == "" {
		msg = "The requested resource was not found"
	}

	return NewError(
		http.StatusNotFound,
		msg,
		setErrCode(errorCode, "not_found"),
	)
}

// Conflict creates a error response with (HTTP 409) code
func Conflict(msg string, errorCode string) *Response {
	if msg == "" {
		msg = "The requested resource was not found"
	}

	return NewError(
		http.StatusConflict,
		msg,
		setErrCode(errorCode, "conflict"),
	)
}

// InternalServerError creates a error response with (HTTP 500)code
func InternalServerError(msg string, errorCode string) *Response {
	if msg == "" {
		msg = "Something went wrong on our end."
	}

	return NewError(
		http.StatusInternalServerError,
		msg,
		setErrCode(errorCode, "internal_server_error"),
	)
}

func setErrCode(errorCode string, defaultCode string) string {
	if errorCode == "" {
		return defaultCode
	}
	return errorCode
}

// AsByte returns the response struct as byte
func (r *Response) AsByte() ([]byte, error) {
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(r)
	if err != nil {
		return buf.Bytes(), nil
	}
	return buf.Bytes(), nil
}

// Decode byte to response
func DecodeToResponse(s []byte, vPtr any) error {
	dec := gob.NewDecoder(bytes.NewReader(s))
	err := dec.Decode(&vPtr)
	return err
}
