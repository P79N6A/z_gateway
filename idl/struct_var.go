package idl

import "github.com/Jeffail/gabs"

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
}

func (structVar *StructVar) StructVarIsEqual(otherStructVar *StructVar) bool {
	return structVar.Order == otherStructVar.Order &&
		structVar.Name == otherStructVar.Name &&
		structVar.Type == otherStructVar.Type &&
		structVar.SubType == otherStructVar.SubType &&
		structVar.IsRequired == otherStructVar.IsRequired &&
		structVar.Default == otherStructVar.Default
}

func (structVar *StructVar) ToJson() *gabs.Container{
	jsonObj := gabs.New()
	jsonObj.Set(structVar.Order, "order")
	jsonObj.Set(structVar.Name, "name")
	jsonObj.Set(structVar.Type, "type")
	jsonObj.Set(structVar.IsRequired, "is_required")
	jsonObj.Set(structVar.Desc, "desc")

	if structVar.Default != nil {
		jsonObj.Set(structVar.Default, "default")
	}

	for k, v := range structVar.Extra {
		jsonObj.Set(v, "extra." + k)
	}

	if structVar.Type == "list" {

		// TODO:: structVar.SubType 为普通类型

		// TODO:: structVar.SubType 为object

		// TODO:: structVar.SubType 为item类型


	}

	// TODO:: structVar.Type为object类型

	return jsonObj
}


