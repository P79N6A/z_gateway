package test

import (
	"fmt"
	"github.com/CharellKing/z_gateway/protocol"
	"testing"
)

func TestJson2Struct(t *testing.T) {
	json2Struct := protocol.NewParams2Struct("RequestBody", "请求body", nil)
	if err := json2Struct.Load("/Users/ck/.gvm/pkgsets/go1.9.2/global/src/github.com/CharellKing/z_gateway/test/param.json"); err != nil {
		fmt.Println(err)
	}

	if err := json2Struct.ToStructs(); err != nil {
		fmt.Println(err)
	}

	fmt.Println(json2Struct.ToProtobufStr())

}