package test

import (
	"fmt"
	"github.com/CharellKing/z_gateway/restful2protobuf"
	"testing"
)

func TestRestful2Protobuf(t *testing.T) {
	json2Proto := restful2protobuf.NewJson2Protobuf("RequestBody", "请求body")
	if err := json2Proto.Load("/Users/ck/.gvm/pkgsets/go1.9.2/global/src/github.com/CharellKing/z_gateway/test/body.json"); err != nil {
		fmt.Println(err)
	}

	if err := json2Proto.ToProtoMessages(); err != nil {
		fmt.Println(err)
	}

	fmt.Println(json2Proto.ToStr())
}