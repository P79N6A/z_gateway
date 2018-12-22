package idl

import "github.com/Jeffail/gabs"

// api's object
type ApiObj struct {
	FuncName   string
	FuncParam  ApiParam
	FuncReturn string

	Uri         string
	RequestType string
	Desc        string

	Version     int64

	StructsMap *map[string]*StructObj
}

func (apiObj *ApiObj) ToJson() (*gabs.Container) {
	jsonObj := gabs.New()

	jsonObj.Set(apiObj.FuncName, "method")
	jsonObj.Set(apiObj.RequestType, "request_type")
	jsonObj.Set(apiObj.Uri, "uri")

	//jsonObj.Set(apiObj.FuncParam., "pack_param_type")

	return jsonObj
}
