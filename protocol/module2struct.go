package protocol

import (
	"errors"
	"fmt"
	"github.com/Jeffail/gabs"
	"github.com/iancoleman/strcase"
	"io/ioutil"
)

type Module2Struct struct {
	JsonObj *gabs.Container

	ApiStructs []*Api2Struct

	StructsMap map[string]*StructObj

	ModuleName string
	Desc string
	Version string
}

func NewModule2Struct(ModuleName string, Desc string, JsonObj *gabs.Container) (*Module2Struct) {
	var module2Struct Module2Struct
	module2Struct.ModuleName = ModuleName
	module2Struct.Desc = Desc
	module2Struct.JsonObj = JsonObj

	module2Struct.StructsMap = make(map[string]*StructObj)

	return &module2Struct
}

func (module2Struct *Module2Struct) Loads(jsonStr []byte) (error) {
	var err error = nil
	if module2Struct.JsonObj, err = gabs.ParseJSON(jsonStr); err != nil {
		return err
	}

	return nil
}

func (module2Struct *Module2Struct) Load(filename string) (error) {
	var (
		content []byte
		err error
	)

	content, err = ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	if module2Struct.JsonObj, err = gabs.ParseJSON(content); err != nil {
		return err
	}

	return nil
}

func (module2Struct *Module2Struct) ToStructs() (error) {
	if false == module2Struct.JsonObj.ExistsP("module_name") {
		return errors.New("module_name不存在")
	}

	if false == module2Struct.JsonObj.ExistsP("apis") {
		return errors.New("interfaces不存在")
	}

	module2Struct.ModuleName = module2Struct.JsonObj.Path("module_name").Data().(string)

	if true == module2Struct.JsonObj.ExistsP("desc") {
		module2Struct.Desc = module2Struct.JsonObj.Path("desc").Data().(string)
	}

	if true == module2Struct.JsonObj.ExistsP("version") {
		module2Struct.Version = module2Struct.JsonObj.Path("version").Data().(string)
	}


	apis, _ := module2Struct.JsonObj.S("apis").Children()

	for _, api := range apis {
		if api.ExistsP("uri") == false {
			return errors.New("接口中uri属性不存在")
		}

		if api.ExistsP("request_type") == false {
			return errors.New("接口中request_type属性不存在")
		}

		if api.ExistsP("params") == false {
			return errors.New("接口中params属性不存在")
		}

		uri := api.Path("uri").Data().(string)
		requestType := api.Path("request_type").Data().(string)
		desc := ""
		if api.ExistsP("desc") == true {
			desc = api.Path("desc").Data().(string)
		}

		params := api.Path("params")

		api2Struct := NewApi2Struct(uri, requestType, desc, &module2Struct.StructsMap, params)
		api2Struct.ToStructs()
		module2Struct.ApiStructs = append(module2Struct.ApiStructs, api2Struct)
	}
	return nil
}

func (module2Struct *Module2Struct) ToProtobufStr() (string) {
	protoSource := ""

	protoSource += "// " + module2Struct.Desc + "\n"

	for _, api := range module2Struct.ApiStructs {
		protoSource += api.ToProtobufStr()
		protoSource += "\n"
	}

	protoSource = fmt.Sprintf("%s\nservice %s {\n", protoSource, strcase.ToCamel(module2Struct.ModuleName))
	for _, api := range module2Struct.ApiStructs {
		protoSource += "\t" + api.ToProtobufFuncStr() + "\n"
	}
	protoSource += "}\n"

	return protoSource
}




