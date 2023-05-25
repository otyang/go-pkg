package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"github.com/otyang/go-pkg/response"
)

func main() {
	success := true
	httpStatusCode := 400
	msg := "Request was successful"
	var errorCode *string = nil
	data := struct{ userName string }{"exampleUser"}

	resp := response.New(httpStatusCode, success, msg, errorCode, data)
	v, _ := json.Marshal(resp)
	fmt.Println(string(v))
	/*{
		"success": true,
		"message": "Request was successful",
		"data": {
			"username": "exampleUser",
		}
	}*/

	// A shortcut helper function to create an error response.
	rsp := response.NewError(http.StatusBadRequest, "this is a bad request", "1000_validation")
	fmt.Println(rsp)

	// An shortcut helper function to create a success response.
	rsp = response.NewSuccess(http.StatusOK, "request was successful", data)
	fmt.Println(rsp)

	// Other helper functions
	rsp = response.Unauthorized("an error occured on our end", "un_authorized")
	fmt.Println(rsp)

	// Note by default even if the error codes are not defined. response gives a default error code
	// default response code for internal server error is 'internal_Server_Error' so this means
	// out put of 'example1' and 'example2' are same
	example1 := response.InternalServerError("an error occured on our end", "internal_server_error")
	example2 := response.InternalServerError("an error occured on our end", "")
	if reflect.DeepEqual(example1, example2) {
		fmt.Println("yes same")
	}

	// to disable default errorCode do this
	response.UseDefaultCode = false
	example3 := response.InternalServerError("an error occured on our end", "")
	fmt.Println(example3) // now the errorCode parameter will be empty
}

// Returning Response as an error:
// you can return response as an error, because it satisfy the golang error interface
func ListSearchEnginesEndPoint() error {
	return response.InternalServerError("error connecting to database", "")
}
