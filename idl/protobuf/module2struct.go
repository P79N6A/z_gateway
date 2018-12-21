package protobuf

import (
	"fmt"
	"github.com/emicklei/proto"
	"strings"
)


type Protobuf2Struct struct {
	Content string
}

func NewProtobuf2Struct(Content string) (*Protobuf2Struct) {
	proto2Struct := Protobuf2Struct{Content}

	return &proto2Struct
}

func HandleService(s *proto.Service) {
	fmt.Println("hello")
}

func HandleMessage(m *proto.Message) {
	fmt.Println("hello")
}

func HandleEnum(m *proto.Enum) {
	fmt.Println("hello")
}

func (proto2Struct *Protobuf2Struct) ToStructs() error {
	reader := strings.NewReader(proto2Struct.Content)

	parser := proto.NewParser(reader)
	definition, err := parser.Parse()
	if err != nil {
		return err
	}

	proto.Walk(definition, proto.WithService(HandleService),
							proto.WithEnum(HandleEnum),
							proto.WithMessage(HandleMessage))
	return nil
}




