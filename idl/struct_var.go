package idl

import (
	"github.com/Jeffail/gabs"
)

// struct' variable
type StructVar struct {
	Order      int32
	Name       string
	Type       string
	SubType    string
	IsRequired bool
	Desc       string
	Default    interface{}
	Extra      map[string]interface{}

	// TODO:: StructMaps
	StructsMap map[string]*StructObj
}

func (structVar *StructVar) StructVarIsEqual(otherStructVar *StructVar) bool {
	return structVar.Order == otherStructVar.Order &&
		structVar.Name == otherStructVar.Name &&
		structVar.Type == otherStructVar.Type &&
		structVar.SubType == otherStructVar.SubType &&
		structVar.IsRequired == otherStructVar.IsRequired &&
		structVar.Default == otherStructVar.Default
}

func (structVar *StructVar) ItemToJson() *gabs.Container {
	if _, ok := BASIC_JSON_TYPES[structVar.SubType]; ok == true {
		jsonObj := gabs.New()
		jsonObj.Set(structVar.Type, "type")
		return jsonObj
	}

	if structObj, ok := structVar.StructsMap[structVar.SubType]; ok == true {
		return structObj.ToJson()
	}

	return nil
}


func (structVar *StructVar) ToJson() *gabs.Container{
	jsonObj := gabs.New()
	jsonObj.Set(structVar.Order, "order")
	jsonObj.Set(structVar.Name, "name")
	jsonObj.Set(structVar.IsRequired, "is_required")
	jsonObj.Set(structVar.Desc, "desc")

	if _, ok := BASIC_JSON_TYPES[structVar.Type]; ok == true {
		jsonObj.Set(structVar.Type, "type")
		return jsonObj

	}

	if structVar.Type == "list" {
		jsonObj.Set(structVar.Type, "type")
		jsonObj.Set(structVar.ItemToJson(), "item")
		return jsonObj
	}

	if structObj, ok := structVar.StructsMap[structVar.Type]; ok == true {
		jsonObj.Set("object", "type")
		jsonObj.Set(structVar.Type, "name")

		jsonObj.Merge(structObj.ToJson())
	}

	return jsonObj
}


