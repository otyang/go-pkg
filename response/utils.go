package response

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// IsJsonErrorGetDetails takes an error and check if it is a JSON Error.
// If it is, it returns true and the type of json Error
// Else it returns the error back as it came
func IsJsonErrorGetDetails(err error) (ok bool, e error) {
	var syntaxError *json.SyntaxError
	var unmarshalTypeError *json.UnmarshalTypeError

	if err != nil {
		switch {
		case errors.As(err, &syntaxError):
			return true, fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)
		case errors.Is(err, io.ErrUnexpectedEOF):
			return true, errors.New("body contains badly-formed JSON")
		case errors.As(err, &unmarshalTypeError):
			return true, fmt.Errorf(
				"body contains incorrect JSON type [for field %q] (at character %d)",
				unmarshalTypeError.Field,
				unmarshalTypeError.Offset,
			)

		case errors.Is(err, io.EOF):
			return true, errors.New("body must not be empty")

		default:
			return false, err
		}
	}
	return false, nil
}

// JSON send output when using standard library
// it marshals 'response'struct to JSON, escapes HTML & sets the
// Content-Type as 'application/json' all via the standard library.
func JSON(w http.ResponseWriter, vPtr *Response, headers map[string]string) {
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(true)

	if err := enc.Encode(vPtr); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	for index, value := range headers {
		w.Header().Add(index, value)
	}
	w.WriteHeader(vPtr.StatusCode)
	w.Write(buf.Bytes())
}
