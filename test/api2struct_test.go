package test

import (
	"fmt"
	"github.com/CharellKing/z_gateway/protocol"
	"testing"
)

func TestApi2Struct(t *testing.T) {
	var structMap map[string]*protocol.StructObj
	structMap = make(map[string]*protocol.StructObj)

	api2Struct := protocol.NewApi2Struct("/api/user/get_task_list", "POST", "获取任务list", &structMap, nil)
	if err := api2Struct.Load("/Users/ck/.gvm/pkgsets/go1.9.2/global/src/github.com/CharellKing/z_gateway/test/api.json"); err != nil {
		fmt.Println(err)
	}

	if err := api2Struct.ToStructs(); err != nil {
		fmt.Println(err)
	}

	fmt.Println(api2Struct.ToProtobufStr())

	fmt.Println(api2Struct.ToProtobufFuncStr())
}