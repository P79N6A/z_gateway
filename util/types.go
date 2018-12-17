package util

import "unicode"

var G_JSON_TYPES = map[string]bool {
	"bool": true,
	"int8": true, "uint8": true,
	"int16": true, "uint16": true,
	"int32": true, "uint32": true,
	"int64": true, "uint64": true,
	"float32": true, "float64": true,
	"object": true, "list": true,
	"string": true,
}

var G_TYPES_RANGE = map[string]func(float64)bool {
	"int8": func(val float64) bool {return val >= -128 && val <= 127},
	"uint8": func(val float64) bool {return val >= 0 && val <= 255},
	"int16": func(val float64) bool {return val >= -32768 && val <= 32767},
	"uint16": func(val float64) bool {return val >= 0 && val <= 65535},
	"int32": func(val float64) bool {return val >= -2147483648 && val <= 2147483647},
	"uint32": func(val float64) bool {return val >= 0 && val <= 4294967295},
	"float32": func(val float64) bool {return val >= -3.40E+38 && val <= +3.40E+38},
}

var G_CONVERT_TYPES = map[string] func(float64)interface{} {
	"int8": func(val float64) interface{} {return int8(val)},
	"uint8": func(val float64) interface{} {return uint8(val)},
	"int16": func(val float64) interface{} {return int16(val)},
	"uint16": func(val float64) interface{} {return uint16(val)},
	"int32": func(val float64) interface{} {return int32(val)},
	"uint32": func(val float64) interface{} {return uint32(val)},
	"float32": func(val float64) interface{} {return float32(val)},
	"float64": func(val float64) interface{} {return val},
}

var G_PROTO_TYPES = map[string] bool {
	"double": true, "float": true,
	"int32": true, "int64": true,
	"uint32": true, "uint64": true,
	"sint32": true, "sint64": true,
	"fixed32": true, "fixed64": true,
	"sfixed32": true, "sfixed64": true,
	"bool": true, "string": true,
	"bytes": true}

func IsIdentifier(str []rune) bool {
	count := 0
	for i, ch := range str {
		if i == 0 && (ch != '_' && !unicode.IsLetter(ch)) {
			return false
		}

		if unicode.IsDigit(ch) || unicode.IsLetter(ch) {
			count ++
		}
	}

	return count > 0
}
