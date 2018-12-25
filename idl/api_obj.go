package idl

import (
	"github.com/Jeffail/gabs"
	"github.com/sirupsen/logrus"
)

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

	jsonObj.Set(apiObj.FuncParam.ParamName, "packed_param_name")

	for _, subApiParam := range apiObj.FuncParam.SubParams {
		apiParamTypeStr := ""
		if subApiParam.HttpParamType == ParamsType {
			apiParamTypeStr = "params"
		} else if subApiParam.HttpParamType == HeadersType {
			apiParamTypeStr = "headers"
		} else if subApiParam.HttpParamType == BodyType {
			apiParamTypeStr = "body"
		}

		if apiParamTypeStr != "" {
			if structObj, ok := (*apiObj.StructsMap)[subApiParam.SubParamTypeName]; ok == true {
				structJson := structObj.ToJson()

				// TODO:: 没有设置
				structJson.Set(apiParamTypeStr, "sub_param_type")
				jsonObj.ArrayAppend(structJson.Data(), "sub_params")
			} else {
				logrus.Error("struct(%s) not exist", subApiParam.SubParamTypeName)
			}
		} else {
			logrus.Error("struct(%s)'s sub param type(%s) not exist", subApiParam.SubParamTypeName, apiParamTypeStr)
		}
	}

	return jsonObj
}
