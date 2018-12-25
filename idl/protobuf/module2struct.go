package protobuf

import (
	"bytes"
	"github.com/CharellKing/z_gateway/idl"
	"github.com/emicklei/proto"
	log "github.com/sirupsen/logrus"
)


type Module2Struct struct {
	idl.ModuleObj

	Definition *proto.Proto
}

func NewModule2Struct(content []byte) (*Module2Struct) {
	var module2Struct Module2Struct

	module2Struct.Extra = make(map[string]interface{})

	module2Struct.StructsMap = make(map[string]*idl.StructObj)

	reader := bytes.NewReader(content)

	parser := proto.NewParser(reader)

	var err error
	if module2Struct.Definition, err = parser.Parse(); err != nil {
		log.Error("analyze proto file failed: %v", err)
		return nil
	}

	return &module2Struct
}

func (module2Struct *Module2Struct) ToStructs() (*idl.ModuleObj) {

	module2Struct.analyzeModule("", module2Struct.Definition)
	module2Struct.fillSubParams()

	return &module2Struct.ModuleObj
}

func (module2Struct *Module2Struct) analyzeModule(path string, protoElem *proto.Proto) {
	for _, elem := range module2Struct.Definition.Elements {
		if syntax, ok := elem.(*proto.Syntax); ok == true {
			module2Struct.analyzeSyntax(syntax)
		} else if message, ok := elem.(*proto.Message); ok == true {
			module2Struct.analyzeMessage(message)
		} else if service, ok := elem.(*proto.Service); ok == true {
			module2Struct.analyzeService(service)
		}
	}
}

func (module2Struct *Module2Struct) analyzeSyntax(syntax *proto.Syntax) {
	module2Struct.Extra["syntax"] = syntax.Value
}

func (module2Struct *Module2Struct) analyzeService(service *proto.Service) {
	module2Struct.ModuleName = service.Name

	if service.Comment != nil {
		module2Struct.Desc = service.Comment.Message()
	}

	for _, elem := range service.Elements {
		rpc := elem.(*proto.RPC)
		module2Struct.analyzeApi(rpc)

	}
}

func (module2Struct *Module2Struct) analyzeApi(rpc *proto.RPC) {
	var apiObj idl.ApiObj
	apiObj.StructsMap = &module2Struct.StructsMap

	apiObj.FuncName = rpc.Name
	apiObj.FuncParam.ParamName = rpc.RequestType
	apiObj.FuncReturn = rpc.ReturnsType

	for _, option := range rpc.Options {
		if option.Name == "(z_gateway)" {
			for _, apiAttr := range option.AggregatedConstants {
				if apiAttr.Name == "uri" {
					apiObj.Uri = apiAttr.Literal.Source
				} else if apiAttr.Name == "type" {
					apiObj.RequestType = apiAttr.Literal.Source
				}
			}
		}
	}

	module2Struct.ApiObjs = append(module2Struct.ApiObjs, &apiObj)
}

func (module2Struct *Module2Struct) analyzeMessage(message *proto.Message) {
	var structObj idl.StructObj

	structObj.Name = message.Name

	if message.Comment != nil {
		structObj.Desc = message.Comment.Message()
	}

	for _, elem := range message.Elements {
		if syntax, ok := elem.(*proto.Syntax); ok == true {
			module2Struct.analyzeSyntax(syntax)
		} else if message, ok := elem.(*proto.Message); ok == true {
			module2Struct.analyzeMessage(message)
		} else if service, ok := elem.(*proto.Service); ok == true {
			module2Struct.analyzeService(service)
		} else if field, ok := elem.(*proto.NormalField); ok == true {
			module2Struct.analyzeField(&structObj, field)
		}
	}

	module2Struct.StructsMap[structObj.Name] = &structObj

}

func (module2Struct *Module2Struct) analyzeField(structObj *idl.StructObj, field *proto.NormalField) {
	var structVar idl.StructVar

	structVar.StructsMap = &module2Struct.StructsMap

	structVar.Extra = make(map[string]interface{})

	structVar.IsRequired = true
	if field.Optional == false {
		structVar.IsRequired = true
	}


	if field.Repeated == true {
		structVar.Type = "list"
		structVar.SubType = field.Field.Type
	} else {
		structVar.Type = field.Field.Type
	}

	structVar.Name = field.Field.Name
	structVar.Order = int32(field.Field.Sequence)

	if field.Field.Comment != nil {
		structVar.Desc = field.Field.Comment.Message()
	}

	for _, option := range field.Field.Options {
		structVar.Extra[option.Name] = option.Constant.Source
	}

	structObj.Vars = append(structObj.Vars, &structVar)
}

func (module2Struct *Module2Struct) fillSubParams() {
	for _, apiObj := range module2Struct.ApiObjs {
		paramTypes := map[int]bool{}
		paramName := apiObj.FuncParam.ParamName
		if structObj, ok := module2Struct.StructsMap[paramName]; ok == true {
			for _, structVar := range structObj.Vars {
				if paramTypeStr, ok := structVar.Extra["(z_gateway.param_type)"]; ok == true {
					paramType := -1
					if paramTypeStr == "params" {
						paramType = idl.ParamsType
					} else if paramTypeStr == "headers" {
						paramType = idl.HeadersType
					} else if paramTypeStr == "body" {
						paramType = idl.BodyType
					} else {
						log.Error("rpc param(%s)'s sub param(%s) z_gateway.param_type(%s) not in [params, headers, body]",
									paramName, structVar.Name, paramTypeStr)
						continue
					}

					if _, ok := paramTypes[paramType]; ok == false {
						paramTypes[paramType] = true
						apiObj.FuncParam.SubParams = append(apiObj.FuncParam.SubParams,
							idl.SubApiParam{int32(paramType), structVar.Type})
					} else {
						log.Error("rpc param(%s)'s sub param(%s) z_gateway.param_type(%s) existed",
							paramName, structVar.Name, paramTypeStr)
					}

				} else {
					log.Error("rpc param's sub param(%s) not set z_gateway.param_type", structVar.Name)
					continue
				}
			}
		}
	}
}
