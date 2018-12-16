package test

import (
	"fmt"
	"github.com/CharellKing/z_gateway/protocol"
	"testing"
)

func TestModule2Struct(t *testing.T) {
	module2Struct := protocol.NewModule2Struct("", "", nil)
	if err := module2Struct.Load("/Users/ck/.gvm/pkgsets/go1.9.2/global/src/github.com/CharellKing/z_gateway/test/module.json"); err != nil {
		fmt.Println(err)
	}

	if err := module2Struct.ToStructs(); err != nil {
		fmt.Println(err)
	}

	fmt.Println(module2Struct.ToProtobufStr())
}