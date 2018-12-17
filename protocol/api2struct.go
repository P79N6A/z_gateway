package protocol

import (
	"fmt"
	"github.com/CharellKing/z_gateway/util"
	"github.com/Jeffail/gabs"
	"github.com/iancoleman/strcase"
	"io/ioutil"
)

type Api2Struct struct {
	JsonObj *gabs.Container

	RequestParams *Params2Struct
	RequestHeaders *Params2Struct
	RequestBody *Params2Struct
	ResponseBoby *Params2Struct

	Uri string
	RequestType string
	Desc string
	StructsMap *map[string]*StructObj
}

func NewApi2Struct(Uri string, RequestType string, Desc string, structsMap *map[string]*StructObj, JsonObj *gabs.Container) (*Api2Struct) {
	var api2Struct Api2Struct
	api2Struct.Uri = Uri
	api2Struct.RequestType = RequestType
	api2Struct.Desc = Desc
	api2Struct.StructsMap = structsMap
	api2Struct.JsonObj = JsonObj
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
		requestParamsObj := api2Struct.JsonObj.Path("request_params")
		api2Struct.RequestParams = NewParams2Struct(requestParamsName, "",  api2Struct.StructsMap, requestParamsObj)
		if err := api2Struct.RequestParams.ToStructs(); err != nil {
			return err
		}
	}

	if api2Struct.JsonObj.ExistsP("request_headers") {
		requestHeadersName := strcase.ToCamel(fmt.Sprintf("%s_%s", apiName, "RequestHeaders"))
		requestHeadersObj := api2Struct.JsonObj.Path("request_headers")
		api2Struct.RequestHeaders = NewParams2Struct(requestHeadersName, "", api2Struct.StructsMap, requestHeadersObj)
		if err := api2Struct.RequestHeaders.ToStructs(); err != nil {
			return err
		}
	}

	if api2Struct.JsonObj.ExistsP("request_body") {
		requestBodyName := strcase.ToCamel(fmt.Sprintf("%s_%s", apiName, "RequestBody"))

		requestBodyObj := api2Struct.JsonObj.Path("request_body")
		api2Struct.RequestBody = NewParams2Struct(requestBodyName, "", api2Struct.StructsMap, requestBodyObj)
		if err := api2Struct.RequestBody.ToStructs(); err != nil {
			return err
		}
	}

	if api2Struct.JsonObj.ExistsP("response_body") {
		responseBodyName := strcase.ToCamel(fmt.Sprintf("%s_%s", apiName, "ResponseBody"))
		repsonseBodyObj := api2Struct.JsonObj.Path("response_body")
		api2Struct.ResponseBoby = NewParams2Struct(responseBodyName, "", api2Struct.StructsMap, repsonseBodyObj)
		if err := api2Struct.ResponseBoby.ToStructs(); err != nil {
			return err
		}
	}

	return nil
}

func (api2Struct *Api2Struct) ToProtobufStr() (string) {
	protoSource := ""

	if api2Struct.RequestParams != nil {
		protoSource += api2Struct.RequestParams.ToProtobufStr()
		protoSource += "\n"
	}


	if api2Struct.RequestHeaders != nil {
		protoSource += api2Struct.RequestHeaders.ToProtobufStr()
		protoSource += "\n"
	}

	if api2Struct.RequestBody != nil {
		protoSource += api2Struct.RequestBody.ToProtobufStr()
		protoSource += "\n"
	}

	if api2Struct.ResponseBoby != nil {
		protoSource += api2Struct.ResponseBoby.ToProtobufStr()
		protoSource += "\n"
	}

	return protoSource
}

func (api2Struct *Api2Struct) ToProtobufFuncStr() (string) {
	apiName := util.GetUriLastSeg(api2Struct.Uri)

	protoSource := ""
	protoSource += fmt.Sprintf("rpc %s(", strcase.ToCamel(apiName))

	isFirst := true
	if api2Struct.RequestParams != nil {
		protoSource += api2Struct.RequestParams.JsonName
		isFirst = false
	}

	if isFirst == false {
		protoSource += ", "
	}

	if api2Struct.RequestHeaders != nil {
		protoSource += api2Struct.RequestHeaders.JsonName
		isFirst = false
	}

	if isFirst == false {
		protoSource += ", "
	}

	if api2Struct.RequestBody != nil {
		protoSource += api2Struct.RequestBody.JsonName
		isFirst = false
	}

	protoSource += ") "


	if api2Struct.ResponseBoby != nil {
		protoSource += fmt.Sprintf(" return (%s)", api2Struct.RequestBody.JsonName)
	}

	protoSource += ";\n"
	if api2Struct.Desc != "" {
		protoSource += api2Struct.Desc
	}

	return protoSource
}
