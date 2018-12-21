package test

import (
	"fmt"
	"github.com/CharellKing/z_gateway/idl/json"
	"os"
	"testing"
)

func TestModule2Struct(t *testing.T) {
	goPath := os.Getenv("GOPATH")
	module2Struct := json.NewModule2Struct("", "", nil)
	if err := module2Struct.Load(goPath + "/src/github.com/CharellKing/z_gateway/test/json/module.json"); err != nil {
		fmt.Println(err)
	}

	if err := module2Struct.ToStructs(); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(module2Struct.ToProtobufStr())
}