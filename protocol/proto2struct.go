package protocol
//
//import (
//	"errors"
//	"fmt"
//	"github.com/CharellKing/z_gateway/util"
//	"io/ioutil"
//	"unicode"
//)
//
////type ProtobufFunc struct {
////	Func
////}
//
//type Protobuf2Struct struct {
//	Contents []rune
//
//	Types map[string]bool
//
//	Row int
//	Col int
//	Index int
//}
//
//func NewProtobuf2Struct(Contents []rune) (*Protobuf2Struct) {
//	var proto2Struct Protobuf2Struct
//	proto2Struct.Contents = Contents
//	proto2Struct.Row = 1
//	proto2Struct.Col = 1
//	proto2Struct.Index = 0
//	return &proto2Struct
//}
//
//func (proto2Struct *Protobuf2Struct) Load(filename string) (error) {
//	contents, err := ioutil.ReadFile(filename)
//	if err != nil {
//		return err
//	}
//
//	proto2Struct.Contents = []rune(string(contents))
//	return nil
//}
//
//func (proto2Struct *Protobuf2Struct) ToStructs() (error) {
//	return nil
//}
//
//func (proto2Struct *Protobuf2Struct) analyze_file() (error) {
//	return nil
//}
//
//func (proto2Struct *Protobuf2Struct) analyze_message() (error) {
//	return nil
//}
//
////func (proto2Struct *Protobuf2Struct) analyze_service() (string, [][]rune, [][]rune, error) {
////	return "", [][]rune,
////}
//
//func (proto2Struct *Protobuf2Struct) analyze_func() (error) {
//	funcKey := proto2Struct.analyze_identifer()
//	if string(funcKey) != "rpc" {
//		return proto2Struct.error("接口缺少关键字rpc")
//	}
//
//	funcName := proto2Struct.analyze_identifer()
//	if util.IsIdentifier(funcName) == false {
//		return proto2Struct.error("函数名称不为标识符")
//	}
//
//	funcParams, err := proto2Struct.analyze_params()
//	if err != nil {
//		return err
//	}
//
//	funcRetKey := proto2Struct.analyze_identifer()
//	if string(funcRetKey) != "return" {
//		return proto2Struct.error("函数缺少return")
//	}
//
//	funcRetParams, err := proto2Struct.analyze_params()
//	if err != nil {
//		return err
//	}
//
//
//	// return funcName, funcParams, funcRetParams, nil
//	return nil
//}
//
//func (proto2Struct *Protobuf2Struct) analyze_params() ([][]rune, error) {
//	var endCh *rune;
//	var funcParams [][]rune
//
//	ch := proto2Struct.next_non_space_ch()
//	if *ch != '(' {
//		return funcParams, proto2Struct.error("函数缺少(")
//	}
//
//	for {
//		funcParam := proto2Struct.analyze_identifer()
//		if len(funcParam) > 0 {
//			funcParams = append(funcParams, funcParam)
//		} else {
//			ch := proto2Struct.next_non_space_ch()
//			if ch == nil || *ch == ')' {
//				endCh = ch
//				break
//			}
//		}
//	}
//
//	if endCh == nil {
//		return funcParams, proto2Struct.error("函数参数列表必须以)结尾")
//	}
//
//	return funcParams, nil
//}
//
//func (proto2Struct *Protobuf2Struct) analyze_type() ([]rune, error) {
//	identifier := proto2Struct.analyze_identifer()
//	if false == util.IsIdentifier(identifier) {
//		return identifier, proto2Struct.error(fmt.Sprintf("非标识符(%v)", identifier))
//	}
//
//	if _, ok := util.G_PROTO_TYPES[string(identifier)]; ok == false {
//		return identifier, proto2Struct.error(fmt.Sprintf("类型不存在(%v)", identifier))
//	}
//	return identifier, nil
//}
//
//func (proto2Struct *Protobuf2Struct) analyze_var() ([]rune, error) {
//	identifier := proto2Struct.analyze_identifer()
//	if false == util.IsIdentifier(identifier) {
//		return identifier, proto2Struct.error("非标识符")
//	}
//	return identifier, nil
//}
//
//func (proto2Struct *Protobuf2Struct) analyze_identifer() ([]rune) {
//	proto2Struct.skip_space()
//
//	var identifier []rune
//	for {
//		ch := proto2Struct.next_ch()
//		if ch == nil {
//			break
//		}
//
//		if unicode.IsDigit(*ch) || unicode.IsLetter(*ch) || *ch == '_' {
//			identifier = append(identifier, *ch)
//			continue
//		}
//
//		break
//	}
//
//	return identifier
//}
//
//func (proto2Struct *Protobuf2Struct) skip_space() {
//	for {
//		ch := proto2Struct.next_ch()
//		if ch == nil {
//			break
//		}
//
//		if unicode.IsSpace(*ch) {
//			continue
//		}
//
//		break
//	}
//}
//
//func (proto2Struct *Protobuf2Struct) next_non_space_ch() (*rune) {
//	proto2Struct.skip_space()
//	return proto2Struct.next_ch()
//}
//
//func (proto2Struct *Protobuf2Struct) next_ch() (*rune) {
//	if proto2Struct.Index >= len(proto2Struct.Contents) {
//		return nil
//	}
//
//	ch := proto2Struct.Contents[proto2Struct.Index]
//	if ch == '\n' || ch == '\r' {
//		proto2Struct.Row ++
//		proto2Struct.Col = 0
//	} else {
//		proto2Struct.Col++
//	}
//
//	proto2Struct.Index++
//	return &ch
//}
//
//func (proto2Struct *Protobuf2Struct) error(msg string) error {
//	return errors.New(fmt.Sprintf("[%d, %d]%s", proto2Struct.Row, proto2Struct.Col, msg))
//}
