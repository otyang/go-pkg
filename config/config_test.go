package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type Example_ConfigStruct struct {
	AppName string `env:"APP_NAME" env-default:"Auth"`
	AppPort int    `env:"APP_PORT" env-default:"Auth"`
}

func TestLoad(t *testing.T) {
	actual1 := Example_ConfigStruct{}
	err1 := Load("./file-does-not-exists.lol", &actual1)
	assert.Error(t, err1, "case 1: it should be an error. but it isn't")

	// case 2

	expected2 := Example_ConfigStruct{
		AppName: "testing-service",
		AppPort: 9000,
	}
	actual2 := Example_ConfigStruct{}
	err2 := Load("./_test.env", &actual2)
	assert.Equal(t, expected2, actual2, "case 2: it should be same. but it isn't")
	assert.ErrorIs(t, err2, nil, "case 2: it shouldn't be an error. but it is")
}

func TestMustLoad(t *testing.T) {
	actual1 := Example_ConfigStruct{}
	actualPanic1 := func() {
		MustLoad("./file-does-not-exists.lol", &actual1)
	}
	assert.Panics(t, actualPanic1, "case 1: it should panic. it didnt")

	// case 2
	expected1 := Example_ConfigStruct{
		AppName: "testing-service",
		AppPort: 9000,
	}
	actual2 := Example_ConfigStruct{}
	MustLoad("./_test.env", &actual2)
	assert.Equal(t, expected1, actual2, "case 1: it should panic. it didnt")
}
