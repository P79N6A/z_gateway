package idl

import (
	"github.com/Jeffail/gabs"
	"github.com/bradfitz/slice"
)

// struct's object
type StructObj struct {
	Name string
	Desc string
	Vars []*StructVar
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

	for i := 0; i < len(structObj.Vars); i++ {
		if false == structObj.Vars[i].StructVarIsEqual(otherStructObj.Vars[i]) {
			return false
		}
	}
	return true
}

func (structObj *StructObj) ToStruct() *gabs.Container {
	return nil
}