package response

import (
	"bytes"
	"encoding/gob"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_setErrCode(t *testing.T) {
	UseDefaultCode = true
	expected := "default_for_bad_request"
	got := setErrCode("", "default_for_bad_request")
	assert.Equal(t, expected, got)

	// case 2
	expected2 := "bad_request"
	got2 := setErrCode("bad_request", "default_for_bad_request")
	assert.Equal(t, expected2, got2)

	// case 3
	UseDefaultCode = false
	expected3 := ""
	got3 := setErrCode("", "")
	assert.Equal(t, expected3, got3)
}

func TestNew(t *testing.T) {
	expected1 := &Response{
		StatusCode: http.StatusOK,
		Success:    true,
		Message:    "new",
		ErrorCode:  nil,
		Data:       nil,
	}

	got1 := New(http.StatusOK, true, "new", nil, nil)
	assert.Equal(t, expected1, got1, "TestNew case 1: should be same. why isnt it?")
}

func TestNewError(t *testing.T) {
	ec := "bad_request"
	expected2 := &Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
		Message:    "new-error",
		ErrorCode:  &ec,
		Data:       nil,
	}

	UseDefaultCode = true
	got2 := NewError(http.StatusBadRequest, "new-error", "bad_request")
	assert.Equal(t, expected2, got2, "TestNew case 2: should be same. why isnt it?")
}

func TestNewSuccess(t *testing.T) {
	data := struct{ ID string }{ID: "hello world"}
	expected3 := &Response{
		StatusCode: http.StatusBadRequest,
		Success:    true,
		Message:    "new-success",
		ErrorCode:  nil,
		Data:       data,
	}

	got3 := NewSuccess(http.StatusBadRequest, "new-success", data)
	assert.Equal(t, expected3, got3, "TestNew case 3: should be same. why isnt it?")
}

func TestResponse_AsByte(t *testing.T) {

	asByte := func(r any) ([]byte, error) {
		buf := bytes.Buffer{}
		enc := gob.NewEncoder(&buf)
		err := enc.Encode(r)
		if err != nil {
			return buf.Bytes(), nil
		}
		return buf.Bytes(), nil
	}

	r := &Response{
		StatusCode: http.StatusOK,
		Success:    true,
		Message:    "its bytes",
		ErrorCode:  nil,
		Data:       nil,
	}

	got, err := r.AsByte()

	expected, _ := asByte(r)

	assert.Equal(t, nil, err, "error should be nil. but it isnt")
	assert.Equal(t, expected, got, "testing for byte")

}
