package idl

const (
	ParamsType  = iota
	HeadersType = iota
	BodyType    = iota
)

// api's param

const (
	HttpParams = iota
	HttpHeaders = iota
	HttpBody = iota
)

type SubApiParam struct {
	HttpParamType int32
	SubParamName string
}

type ApiParam struct {
	ParamName string // 参数类型名称

	SubParams []SubApiParam // 子参数类型
}

