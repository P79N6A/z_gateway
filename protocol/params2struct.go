package protocol

import (
	"errors"
	"fmt"
	"github.com/CharellKing/z_gateway/util"
	"github.com/Jeffail/gabs"
	"github.com/iancoleman/strcase"
	"github.com/bradfitz/slice"
	"io/ioutil"
	"strconv"
)


type StructObj struct {
	Name string
	Desc string
	Vars []*StructVar
}

func (structObj *StructObj) SortStructVars() *StructObj {
	slice.SortInterface(structObj.Vars[:], func(i, j int) bool {
		return structObj.Vars[i].Order < structObj.Vars[j].Order
	})
	return structObj
}

func (structObj *StructObj) ToStr() string {
	protoSource := ""
	if structObj.Desc != "" {
		protoSource += "//" + structObj.Desc + "\n"
	}
	protoSource += "message " + structObj.Name + "\n"
	protoSource += "{" + "\n"
	for _, protoVar := range structObj.Vars {
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
		protoSource += " " + strconv.Itoa(int(protoVar.Order))

		if protoVar.Default != nil {
			protoSource = fmt.Sprintf("%s [default=%v]", protoSource, protoVar.Default)
		}

		protoSource += ";"

		if protoVar.Desc != "" {
			protoSource += " //" + protoVar.Desc
		}
		protoSource += "\n"
	}
	protoSource += "}\n"
	return protoSource
}

type StructVar struct {
	Order int32
	Name string
	Type string
	SubType string
	IsRequired bool
	Desc string
	Default interface{}
}

type Params2Struct struct {
	JsonObj *gabs.Container
	JsonName string
	JsonDesc string
	Structs []*StructObj
}

func NewParams2Struct(JsonName string, JsonDesc string, JsonObj *gabs.Container) *Params2Struct {
	var params2Struct Params2Struct
	params2Struct.JsonObj = JsonObj
	params2Struct.JsonName = JsonName
	params2Struct.JsonDesc = JsonDesc
	return &params2Struct
}

func (params2Struct *Params2Struct) Loads(jsonStr []byte) (error) {
	var err error = nil
	if params2Struct.JsonObj, err = gabs.ParseJSON(jsonStr); err != nil {
		return err
	}

	return nil
}

func (params2Struct *Params2Struct) Load(filename string) (error) {
	var (
		content []byte
		err error
	)

	content, err = ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	if params2Struct.JsonObj, err = gabs.ParseJSON(content); err != nil {
		return err
	}

	return nil
}

func (params2Struct *Params2Struct) ToStructs() (error) {
	_, err := params2Struct.fields_2_obj(params2Struct.JsonName, params2Struct.JsonDesc, params2Struct.JsonObj)
	return err
}

func (params2Struct *Params2Struct) ToProtobufStr() (string) {
	protoSource := ""
	for _, structObj := range params2Struct.Structs {
		protoSource += structObj.ToStr() + "\n"
	}
	return protoSource
}

func (params2Struct *Params2Struct) add_message(msg *StructObj) {
	params2Struct.Structs = append(params2Struct.Structs, msg.SortStructVars())
}

func (params2Struct *Params2Struct) fields_2_obj(messageName string, desc string, obj *gabs.Container) (string, error) {
	if obj == nil {
		return "", nil
	}

	if false == obj.ExistsP("fields") {
		return "", errors.New("fields not exist")
	}

	var message StructObj
	message.Name = strcase.ToCamel(messageName)
	message.Desc = desc
	fields, _ := obj.S("fields").Children()
	for _, field := range fields {
		if protoVar, err := params2Struct.field_2_var(field); err != nil {
			return "", err
		} else {
			message.Vars = append(message.Vars, protoVar)
		}
	}
	params2Struct.add_message(&message)
	return message.Name, nil
}

func(params2Struct *Params2Struct) get_default(protoType string, defaultContainer *gabs.Container) (interface{}, error) {
	if protoType == "string" {
		return defaultContainer.Data().(string), nil
	}

	if protoType == "bool" {
		return 	defaultContainer.Data().(bool), nil
	}

	if protoType == "list" || protoType == "object" {
		return nil, nil
	}

	val := defaultContainer.Data().(float64)
	if protoType == "float64" {
		return val, nil
	}

	if util.G_TYPES_RANGE[protoType](val) {
		return util.G_CONVERT_TYPES[protoType](val), nil
	} else {
		return nil, errors.New(fmt.Sprintf("%v不在%s类型范围之内", val, protoType))
	}

	return nil, nil
}

func (params2Struct *Params2Struct) field_2_var(field *gabs.Container) (*StructVar, error) {
	var protoVar StructVar
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
	if _, ok := util.G_JSON_TYPES[protoVar.Type]; ok == false {
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

	protoVar.Default = nil
	if true == field.Exists("default") {
		var err error
		if protoVar.Default, err = params2Struct.get_default(protoVar.Type, field.Path("default")); err != nil {
			return nil, err
		}
	}

	if protoVar.Type == "object" {
		if objName, err := params2Struct.fields_2_obj(strcase.ToCamel(protoVar.Name), protoVar.Desc, field); err != nil {
			return nil, err
		} else {
			protoVar.Type = objName
		}
	} else if protoVar.Type == "list" {
		if false == field.Exists("item") {
			return nil, errors.New("字段为list类型，不能没有item属性")
		}

		item := field.Path("item")
		if subType, err := params2Struct.field_2_list(item); err != nil {
			return nil, err
		} else {
			protoVar.SubType = subType
		}
	}

	return &protoVar, nil
}

func (params2Struct *Params2Struct) field_2_list(item *gabs.Container) (string, error) {
	if false == item.Exists("type") {
		return "", errors.New("list属性中不存在key[type]")
	}

	itemType := item.Path("type").Data().(string)
	if _, ok := util.G_JSON_TYPES[itemType]; ok == false {
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
		params2Struct.fields_2_obj(objType, itemDesc, item)
		return objType, nil
	}

	return itemType, nil
}





