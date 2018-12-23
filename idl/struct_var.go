package idl

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


