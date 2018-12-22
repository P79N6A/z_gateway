package idl

import "github.com/Jeffail/gabs"

// module's object
type ModuleObj struct {
	ApiObjs []*ApiObj

	StructsMap map[string]*StructObj

	ModuleName string
	Desc       string
	Version    int64
	Extra      map[string]interface{} // 存储一些额外的信息
}

func (moduleObj *ModuleObj) ToJson() *gabs.Container {
	jsonObj := gabs.New()

	jsonObj.Set(moduleObj.ModuleName, "name")
	jsonObj.Set(moduleObj.Desc, "desc")
	jsonObj.Set(moduleObj.Version, "version")

	for k, v := range moduleObj.Extra {
		jsonObj.Set(v, k)
	}

	for _, apiObj := range moduleObj.ApiObjs {
		jsonObj.ArrayAppend(apiObj.ToJson(), "apis")
	}
	return jsonObj
}


