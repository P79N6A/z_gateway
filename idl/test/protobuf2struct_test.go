package test

import (
	"github.com/CharellKing/z_gateway/idl/protobuf"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestProto2Struct(t *testing.T) {
	goPath := os.Getenv("GOPATH")

	content, err := ioutil.ReadFile(goPath + "/src/github.com/CharellKing/z_gateway/test/protobuf/test.proto")
	if err != nil {
		log.Fatal(err)
	}

	module2Struct := protobuf.NewModule2Struct(content)
	if module2Struct == nil {
		return
	}

	moduleObj := module2Struct.ToStructs()
	moduleObj.ToJson()
}
