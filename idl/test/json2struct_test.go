package test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/CharellKing/z_gateway/idl"
	"github.com/sirupsen/logrus"
)

func TestJson2Struct(t *testing.T) {
	goPath := os.Getenv("GOPATH")

	content, err := ioutil.ReadFile(goPath + "/src/github.com/CharellKing/z_gateway/idl/test/samples/module.json")
	if err != nil {
		logrus.Error(err)
	}

	module2Struct := json.NewModule2Struct(nil)
	if err := module2Struct.Loads(content); err != nil {
		logrus.Error(err)
	}

	module2Struct.ToStructs()

	return
}
