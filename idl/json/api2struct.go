package json

import (
	"fmt"
	"github.com/CharellKing/z_gateway/idl"
	"github.com/CharellKing/z_gateway/util"
	"github.com/Jeffail/gabs"
	"github.com/iancoleman/strcase"
	"io/ioutil"
)

type Api2Struct struct {
	idl.ApiObj
	RequestParams2Struct *Param2Struct
	RequestHeaders2Struct *Param2Struct
	RequestBody2Struct *Param2Struct

	ResponseBody2Struct *Param2Struct

	JsonObj *gabs.Container
}

func NewApi2Struct(uri string, requestType string, desc string, structsMap *map[string]*idl.StructObj,
	               jsonObj *gabs.Container) (*Api2Struct) {
	var api2Struct Api2Struct
	api2Struct.Uri = uri
	api2Struct.RequestType = requestType
	api2Struct.Desc = desc
	api2Struct.StructsMap = structsMap
	api2Struct.JsonObj = jsonObj
	return &api2Struct
}

func (api2Struct *Api2Struct) Loads(jsonStr []byte) (error) {
	var err error = nil
	if api2Struct.JsonObj, err = gabs.ParseJSON(jsonStr); err != nil {
		return err
	}

	return nil
}

func (api2Struct *Api2Struct) Load(filename string) (error) {
	var (
		content []byte
		err error
	)

	content, err = ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	if api2Struct.JsonObj, err = gabs.ParseJSON(content); err != nil {
		return err
	}

	return nil
}

func (api2Struct *Api2Struct) ToStructs() (error) {
	apiName := util.GetUriLastSeg(api2Struct.Uri)

	if api2Struct.JsonObj.ExistsP("request_params") {
		requestParamsName := strcase.ToCamel(fmt.Sprintf("%s_%s", apiName, "RequestParams"))
		api2Struct.FuncParamNames = append(api2Struct.FuncParamNames, requestParamsName)

		requestParamsObj := api2Struct.JsonObj.Path("request_params")
		api2Struct.RequestParams2Struct = NewParam2Struct(requestParamsName, "", api2Struct.StructsMap, requestParamsObj)
		if err := api2Struct.RequestParams2Struct.ToStructs(); err != nil {
			return err
		}
	}

	if api2Struct.JsonObj.ExistsP("request_headers") {
		requestHeadersName := strcase.ToCamel(fmt.Sprintf("%s_%s", apiName, "RequestHeaders"))
		api2Struct.FuncParamNames = append(api2Struct.FuncParamNames, requestHeadersName)

		requestHeadersObj := api2Struct.JsonObj.Path("request_headers")
		api2Struct.RequestHeaders2Struct = NewParam2Struct(requestHeadersName, "", api2Struct.StructsMap, requestHeadersObj)
		if err := api2Struct.RequestHeaders2Struct.ToStructs(); err != nil {
			return err
		}
	}

	if api2Struct.JsonObj.ExistsP("request_body") {
		requestBodyName := strcase.ToCamel(fmt.Sprintf("%s_%s", apiName, "RequestBody"))
		api2Struct.FuncParamNames = append(api2Struct.FuncParamNames, requestBodyName)

		requestBodyObj := api2Struct.JsonObj.Path("request_body")
		api2Struct.RequestBody2Struct = NewParam2Struct(requestBodyName, "", api2Struct.StructsMap, requestBodyObj)
		if err := api2Struct.RequestBody2Struct.ToStructs(); err != nil {
			return err
		}
	}

	if api2Struct.JsonObj.ExistsP("response_body") {
		responseBodyName := strcase.ToCamel(fmt.Sprintf("%s_%s", apiName, "ResponseBody"))
		api2Struct.FuncRetNames = append(api2Struct.FuncRetNames, responseBodyName)

		repsonseBodyObj := api2Struct.JsonObj.Path("response_body")
		api2Struct.ResponseBody2Struct = NewParam2Struct(responseBodyName, "", api2Struct.StructsMap, repsonseBodyObj)
		if err := api2Struct.ResponseBody2Struct.ToStructs(); err != nil {
			return err
		}
	}

	return nil
}
