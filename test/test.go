package main

import (
	"fmt"
	"os"

	"github.com/emicklei/proto"
)

func main() {
	goPath := os.Getenv("GOPATH")
	reader, _ := os.Open(goPath + "/src/github.com/CharellKing/z_gateway/test/protobuf/user.proto")
	defer reader.Close()

	parser := proto.NewParser(reader)
	definition, _ := parser.Parse()

	proto.Walk(definition,
		proto.WithService(handleService),
		proto.WithMessage(handleMessage),
		proto.WithEnum(handleEnum))
}

func handleService(s *proto.Service) {
	fmt.Println(s.Name)
}

func handleMessage(m *proto.Message) {
	fmt.Println(m.Name)
}

func handleEnum(m *proto.Enum) {
	fmt.Println(m.Name)
}