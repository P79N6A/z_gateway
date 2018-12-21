package idl

import (
	"fmt"
	"github.com/CharellKing/z_gateway/util"
	"github.com/bradfitz/slice"
	"github.com/iancoleman/strcase"
	"strconv"
)

// struct' variable
type StructVar struct {
	Order int32
	Name string
	Type string
	SubType string
	IsRequired bool
	Desc string
	Default interface{}
}

type EnumVar struct {
	Order int32
	Name  string
	Value string
	Type  string
	Desc  string
}

func (structVar *StructVar) StructVarIsEqual(otherStructVar *StructVar) bool {
	return structVar.Order == otherStructVar.Order &&
		structVar.Name == otherStructVar.Name &&
		structVar.Type == otherStructVar.Type &&
		structVar.SubType == otherStructVar.SubType &&
		structVar.IsRequired == otherStructVar.IsRequired &&
		structVar.Default == otherStructVar.Default
}

// struct's object
type StructObj struct {
	Name string
	Desc string
	Vars []*StructVar
}

type EnumObj struct {
	Name string
	Desc string
	Vars []*EnumVar
}
func (structObj *StructObj) SortStructVars() *StructObj {
	slice.SortInterface(structObj.Vars[:], func(i, j int) bool {
		return structObj.Vars[i].Order < structObj.Vars[j].Order
	})
	return structObj
}

func (structObj *StructObj) StructIsEqual(otherStructObj *StructObj) bool {
	if len(structObj.Vars) != len(otherStructObj.Vars) {
		return false
	}

	for i := 0; i < len(structObj.Vars); i ++ {
		if false == structObj.Vars[i].StructVarIsEqual(otherStructObj.Vars[i]) {
			return false
		}
	}
	return true
}

func (structObj *StructObj) ToProtobufStr() string {
	protoSource := ""
	if structObj.Desc != "" {
		protoSource += "//" + structObj.Desc + "\n"
	}
	protoSource += "message " + structObj.Name + "\n"
	protoSource += "{" + "\n"
	for _, protoVar := range structObj.Vars {
		if protoVar.Type == "list" {
			protoSource += "    " + "repeated"
			protoSource += " " + protoVar.SubType
		} else if protoVar.IsRequired == true {
			protoSource += "    " + "required"
			protoSource += " " + protoVar.Type
		} else {
			protoSource += "    " + "required"
			protoSource += " " + protoVar.Type
		}

		protoSource += " " + protoVar.Name + " ="
		protoSource += " " + strconv.Itoa(int(protoVar.Order))

		if protoVar.Default != nil {
			protoSource = fmt.Sprintf("%s [default=%v]", protoSource, protoVar.Default)
		}

		protoSource += ";"

		if protoVar.Desc != "" {
			protoSource += " //" + protoVar.Desc
		}
		protoSource += "\n"
	}
	protoSource += "}\n"
	return protoSource
}

type ParamStructObj struct {
	StructObj

	StructsMap *map[string]*StructObj
}

func (paramStructObj *ParamStructObj) ToProtobufStr() (string) {
	protoSource := ""
	for _, structObj := range *(paramStructObj.StructsMap) {
		protoSource += structObj.ToProtobufStr() + "\n"
	}
	return protoSource
}

// api's object
type ApiObj struct {

	FuncParamNames []string
	FuncRetNames []string

	Uri string
	RequestType string
	Desc string
	StructsMap *map[string]*StructObj
}

func (apiObj *ApiObj) StructsToProtobufStr() (string) {
	protoSource := ""

	for _, funcParamName := range apiObj.FuncParamNames {
		if structObj, ok := (*apiObj.StructsMap)[funcParamName]; ok == true {
			protoSource += structObj.ToProtobufStr()
			protoSource += "\n"
		}
	}

	for _, funcRetName := range apiObj.FuncRetNames {
		if structObj, ok := (*apiObj.StructsMap)[funcRetName]; ok == true {
			protoSource += structObj.ToProtobufStr()
			protoSource += "\n"
		}
	}

	return protoSource
}

func (apiObj *ApiObj) FuncToProtobufStr() (string) {
	apiName := util.GetUriLastSeg(apiObj.Uri)

	protoSource := ""
	protoSource += fmt.Sprintf("rpc %s(", strcase.ToCamel(apiName))

	isFirst := true
	for _, funcParamName := range apiObj.FuncParamNames {
		if isFirst {
			protoSource += funcParamName
			isFirst = false
		} else {
			protoSource += ", " + funcParamName
		}
	}

	protoSource += ") "


	if len(apiObj.FuncRetNames) >= 0 {
		isFirst := true
		protoSource += "return ("
		for _, funcRetName := range apiObj.FuncRetNames {
			if isFirst {
				protoSource += funcRetName
				isFirst = false
			} else {
				protoSource += ", " + funcRetName
			}
		}
		protoSource += ")"
	}

	protoSource += ";\n"
	if apiObj.Desc != "" {
		protoSource += "// " + apiObj.Desc
	}

	return protoSource
}


// module's object
type ModuleObj struct {
	ApiObjs []*ApiObj

	StructsMap map[string]*StructObj
	EnumsMap map[string]*EnumObj

	ModuleName string
	Desc string
	Version string
	Extra map[string]string // 存储一些额外的信息
}


func (module2Struct *ModuleObj) ToProtobufStr() (string) {
	protoSource := ""

	protoSource += "// " + module2Struct.Desc + "\n"

	for _, structObj := range module2Struct.StructsMap {
		protoSource += structObj.ToProtobufStr()
		protoSource += "\n"
	}

	protoSource = fmt.Sprintf("%s\nservice %s {\n", protoSource, strcase.ToCamel(module2Struct.ModuleName))
	for _, api := range module2Struct.ApiObjs {
		protoSource += "\t" + api.FuncToProtobufStr() + "\n"
	}
	protoSource += "}\n"

	return protoSource
}

