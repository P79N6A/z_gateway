package test

import (
	"fmt"
	"github.com/CharellKing/z_gateway/idl"
	"github.com/CharellKing/z_gateway/idl/json"
	"testing"
)

func TestJson2Struct(t *testing.T) {
	var structMap map[string]*idl.StructObj
	structMap = make(map[string]*idl.StructObj)

	param2Struct := json.NewParam2Struct("RequestBody", "请求body", &structMap,nil)
	if err := param2Struct.Load("/Users/ck/.gvm/pkgsets/go1.9.2/global/src/github.com/CharellKing/z_gateway/test/param.json"); err != nil {
		fmt.Println(err)
	}

	if err := param2Struct.ToStructs(); err != nil {
		fmt.Println(err)
	}

	fmt.Println(param2Struct.ToProtobufStr())

}