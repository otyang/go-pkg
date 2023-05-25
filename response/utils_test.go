package response

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsJsonErrorGetDetails(t *testing.T) {
	v1 := &json.SyntaxError{
		Offset: 1,
	}
	actualok1, actualerr1 := IsJsonErrorGetDetails(v1)
	assert.Error(t, actualerr1, "should be an error. it isnt")
	assert.Equal(t, true, actualok1, "should be true. it isnt")

	// case 2
	v2 := &json.UnmarshalTypeError{
		Offset: 90, Field: "field",
	}
	expected2 := fmt.Errorf(
		"body contains incorrect JSON type [for field %q] (at character %d)",
		v2.Field,
		v2.Offset,
	)
	actualok2, actualerr2 := IsJsonErrorGetDetails(v2)
	assert.Equal(t, true, actualok2, "should be true. it isnt")
	assert.Equal(t, expected2, actualerr2, "should be same. it isnt")

	// case 3
	expected3 := errors.New("body must not be empty")
	actualok3, actualerr3 := IsJsonErrorGetDetails(io.EOF)
	assert.Equal(t, true, actualok3, "should be true. it isnt")
	assert.Equal(t, expected3, actualerr3, "should be same. it isnt")

	// case 4
	expected4 := errors.New("just a normal error")
	actualok4, actualerr4 := IsJsonErrorGetDetails(errors.New("just a normal error"))
	assert.Equal(t, false, actualok4, "should be true. it isnt")
	assert.Equal(t, expected4, actualerr4, "should be same. it isnt")

	// case 5
	actualok5, actualerr5 := IsJsonErrorGetDetails(nil)
	assert.Equal(t, false, actualok5, "should be true. it isnt")
	assert.Equal(t, nil, actualerr5, "should be same. it isnt")
}

func TestJSON(t *testing.T) {
	data := map[string]string{"x-example": "123456789"}
	// req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()

	response := &Response{
		StatusCode: 200,
		Success:    false,
	}
	JSON(res, response, data)

	//
	expected := "123456789"
	got := res.Header().Get("x-example")

	assert.Equal(t, response.StatusCode, res.Code)
	assert.Equal(t, expected, got)
}
