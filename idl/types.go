package idl


var BASIC_JSON_TYPES = map[string]bool {
	"bool": true, "string": true,
	"int8": true, "uint8": true,
	"int16": true, "uint16": true,
	"int32": true, "uint32": true,
	"int64": true, "uint64": true,
	"float32": true, "float64": true,
}

var JSON_TYPES = map[string]bool {
	"bool": true, "string": true,
	"int8": true, "uint8": true,
	"int16": true, "uint16": true,
	"int32": true, "uint32": true,
	"int64": true, "uint64": true,
	"float32": true, "float64": true,
	"object": true, "list": true,
}

var PROTO_2_STRUCTS_TYPE = map[string] string{
	"double": "float64",
	"float":  "float32",
	"int64":  "int64",
	"int32":  "int32",
	"uint64": "uint64",
	"uint32":   "uint32",
	"sint64":   "int64",
	"sint32":   "int32",
	"fixed64":  "int64",
	"fixed32":  "int32",
	"sfixed64": "int64",
	"sfixed32": "int32",
	"bool":     "bool",
	"string":   "string",
	"bytes":    "string",
	"list":     "list",
}
