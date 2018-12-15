package restful2protobuf

import (
	"errors"
	"fmt"
	"github.com/Jeffail/gabs"
	"github.com/iancoleman/strcase"
	"github.com/bradfitz/slice"
	"io/ioutil"
	"strconv"
)


var G_JsonTypes = map[string]bool{"string": true, "int32": true, "int64": true,
                 "double": true, "object": true, "bool": true, "list": true}

type ProtoMessage struct {
	Name string
	Desc string
	Vars []*ProtoVar
}

func (protoMsg *ProtoMessage) SortProtoVars() *ProtoMessage {
	slice.SortInterface(protoMsg.Vars[:], func(i, j int) bool {
		return protoMsg.Vars[i].Order < protoMsg.Vars[j].Order
	})
	return protoMsg
}

func (protoMsg *ProtoMessage) ToStr() string {
	protoSource := ""
	if protoMsg.Desc != "" {
		protoSource += "//" + protoMsg.Desc + "\n"
	}
	protoSource += "message " + protoMsg.Name + "\n"
	protoSource += "{" + "\n"
	for _, protoVar := range protoMsg.Vars {
		if protoVar.Type == "list" {
			protoSource += "    " + "repeated"
			protoSource += " " + protoVar.SubType
		} else if protoVar.IsRequired == true {
			protoSource += "    " + "required"
			protoSource += " " + protoVar.Type
		} else {
			protoSource += "    " + "required"
			protoSource += " " + protoVar.Type
		}

		protoSource += " " + protoVar.Name + " ="
		protoSource += " " + strconv.Itoa(int(protoVar.Order)) + ";"

		if protoVar.Desc != "" {
			protoSource += " //" + protoVar.Desc
		}
		protoSource += "\n"
	}
	protoSource += "}\n"
	return protoSource
}

type ProtoVar struct {
	Order int32
	Name string
	Type string
	SubType string
	IsRequired bool
	Desc string
}

type Json2Protobuf struct {
	JsonObj *gabs.Container
	JsonName string
	JsonDesc string
	Messages []*ProtoMessage
}

func NewJson2Protobuf(JsonName string, JsonDesc string) *Json2Protobuf {
	var json2Proto Json2Protobuf
	json2Proto.JsonObj = gabs.New()
	json2Proto.JsonName = JsonName
	json2Proto.JsonDesc = JsonDesc
	return &json2Proto
}

func (json2Proto *Json2Protobuf) Loads(jsonStr []byte) (error) {
	var err error = nil
	if json2Proto.JsonObj, err = gabs.ParseJSON(jsonStr); err != nil {
		return err
	}

	return nil
}

func (json2Proto *Json2Protobuf) Load(filename string) (error) {
	var (
		content []byte
		err error
	)

	content, err = ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	if json2Proto.JsonObj, err = gabs.ParseJSON(content); err != nil {
		return err
	}

	return nil
}

func (json2Proto *Json2Protobuf) ToProtoMessages() (error) {
	_, err := json2Proto.fields_2_obj(json2Proto.JsonName, json2Proto.JsonDesc, json2Proto.JsonObj)
	return err
}

func (json2Proto *Json2Protobuf) ToStr() (string) {
	protoSource := ""
	for _, protoMsg := range json2Proto.Messages {
		protoSource += protoMsg.ToStr() + "\n"
	}
	return protoSource
}

func (json2Proto *Json2Protobuf) add_message(msg *ProtoMessage) {
	json2Proto.Messages = append(json2Proto.Messages, msg.SortProtoVars())
}

func (json2Proto *Json2Protobuf) fields_2_obj(messageName string, desc string, obj *gabs.Container) (string, error) {
	if obj == nil {
		return "", nil
	}

	if false == obj.ExistsP("fields") {
		return "", errors.New("fields not exist")
	}

	var message ProtoMessage
	message.Name = strcase.ToCamel(messageName)
	message.Desc = desc
	fields, _ := obj.S("fields").Children()
	for _, field := range fields {
		if protoVar, err := json2Proto.field_2_var(field); err != nil {
			return "", err
		} else {
			message.Vars = append(message.Vars, protoVar)
		}
	}
	json2Proto.add_message(&message)
	return message.Name, nil
}

func (json2Proto *Json2Protobuf) field_2_var(field *gabs.Container) (*ProtoVar, error) {
	var protoVar ProtoVar
	if false == field.Exists("name") {
		return nil, errors.New("字段属性中不存在key[name]")
	}

	if false == field.Exists("type") {
		return nil, errors.New("字段属性中不存在key[type]")
	}

	if false == field.Exists("order") {
		return nil, errors.New("字段属性中不存在key[order]")
	}

	protoVar.Type, _ = field.Path("type").Data().(string)
	if _, ok := G_JsonTypes[protoVar.Type]; ok == false {
		return nil, errors.New(fmt.Sprintf("非法的类型[%s]", protoVar.Type))
	}

	order, _ := field.Path("order").Data().(float64)
	protoVar.Order = int32(order)
	protoVar.Name, _ = field.Path("name").Data().(string)
	protoVar.Desc = ""

	protoVar.IsRequired = false
	if true == field.Exists("is_required") {
		protoVar.IsRequired = true
	}

	if true == field.Exists("desc") {
		protoVar.Desc = field.Path("desc").Data().(string)
	}

	if protoVar.Type == "object" {
		if objName, err := json2Proto.fields_2_obj(strcase.ToCamel(protoVar.Name), protoVar.Desc, field); err != nil {
			return nil, err
		} else {
			protoVar.Type = objName
		}
	} else if protoVar.Type == "list" {
		if false == field.Exists("item") {
			return nil, errors.New("字段为list类型，不能没有item属性")
		}

		item := field.Path("item")
		if subType, err := json2Proto.field_2_list(item); err != nil {
			return nil, err
		} else {
			protoVar.SubType = subType
		}
	}

	return &protoVar, nil
}

func (json2Proto *Json2Protobuf) field_2_list(item *gabs.Container) (string, error) {
	if false == item.Exists("type") {
		return "", errors.New("list属性中不存在key[type]")
	}

	itemType := item.Path("type").Data().(string)
	if _, ok := G_JsonTypes[itemType]; ok == false {
		return "", errors.New(fmt.Sprintf("非法的类型[%s]", itemType))
	}

	itemDesc := ""
	if false == item.Exists("desc") {
		itemDesc = item.Path("desc").Data().(string)
	}

	if itemType == "object" {
		if false == item.Exists("name") {
			return "", errors.New("object类型不存在[name]")
		}

		itemName := item.Path("name").Data().(string)
		objType := strcase.ToCamel(itemName)
		json2Proto.fields_2_obj(objType, itemDesc, item)
		return objType, nil
	}

	return itemType, nil
}





