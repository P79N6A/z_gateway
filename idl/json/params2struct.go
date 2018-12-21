package json

import (
	"errors"
	"fmt"
	"github.com/CharellKing/z_gateway/idl"
	"github.com/CharellKing/z_gateway/util"
	"github.com/Jeffail/gabs"
	"github.com/iancoleman/strcase"
	"io/ioutil"
)

type Param2Struct struct {
	idl.ParamStructObj

	JsonObj *gabs.Container
}

func NewParam2Struct(name string, desc string, structsMap *map[string]*idl.StructObj,
	                 jsonObj *gabs.Container) *Param2Struct {
	var param2Struct Param2Struct
	param2Struct.Name = name
	param2Struct.Desc = desc
	param2Struct.StructsMap = structsMap
	param2Struct.JsonObj = jsonObj
	return &param2Struct
}

func (param2Struct *Param2Struct) Loads(jsonStr []byte) (error) {
	var err error = nil
	if param2Struct.JsonObj, err = gabs.ParseJSON(jsonStr); err != nil {
		return err
	}

	return nil
}

func (param2Struct *Param2Struct) Load(filename string) (error) {
	var (
		content []byte
		err error
	)

	content, err = ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	if param2Struct.JsonObj, err = gabs.ParseJSON(content); err != nil {
		return err
	}

	return nil
}

func (param2Struct *Param2Struct) ToStructs() (error) {
	prefixName := ""
	_, err := param2Struct.fields_2_obj(prefixName, param2Struct.Name, param2Struct.Desc, param2Struct.JsonObj)
	return err
}

func (param2Struct *Param2Struct) add_struct(prefixStructName string, shortStructName string, structObj *idl.StructObj) {
	structName := strcase.ToCamel(shortStructName)
	existStructObj, ok := (*param2Struct.StructsMap)[structName]

	if false == ok {
		structObj.Name = structName
		(*param2Struct.StructsMap)[structName] = structObj
		return
	}

	if true == existStructObj.StructIsEqual(structObj) {
		return
	}

	structName = strcase.ToCamel(prefixStructName + "_" + shortStructName)
	structObj.Name = structName
	(*param2Struct.StructsMap)[structName] = structObj
}

func (param2Struct *Param2Struct) fields_2_obj(prefixName string, structName string, desc string,
	                                             obj *gabs.Container) (string, error) {
	if obj == nil {
		return "", nil
	}

	if false == obj.ExistsP("fields") {
		return "", errors.New("fields not exist")
	}

	var structObj idl.StructObj
	structObj.Name = strcase.ToCamel(structName)
	structObj.Desc = desc
	fields, _ := obj.S("fields").Children()
	for _, field := range fields {
		if protoVar, err := param2Struct.field_2_var(prefixName, field); err != nil {
			return "", err
		} else {
			structObj.Vars = append(structObj.Vars, protoVar)
		}
	}

	param2Struct.add_struct(prefixName, structName, &structObj)
	return structObj.Name, nil
}

func(param2Struct *Param2Struct) get_default(protoType string, defaultContainer *gabs.Container) (interface{}, error) {
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

func (param2Struct *Param2Struct) field_2_var(prefixName string, field *gabs.Container) (*idl.StructVar, error) {
	var protoVar idl.StructVar
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
		if protoVar.Default, err = param2Struct.get_default(protoVar.Type, field.Path("default")); err != nil {
			return nil, err
		}
	}

	if protoVar.Type == "object" {
		if objName, err := param2Struct.fields_2_obj(prefixName, strcase.ToCamel(protoVar.Name), protoVar.Desc, field); err != nil {
			return nil, err
		} else {
			protoVar.Type = objName
		}
	} else if protoVar.Type == "list" {
		if false == field.Exists("item") {
			return nil, errors.New("字段为list类型，不能没有item属性")
		}

		item := field.Path("item")
		if subType, err := param2Struct.field_2_list(prefixName + "_" + "Item", item); err != nil {
			return nil, err
		} else {
			protoVar.SubType = subType
		}
	}

	return &protoVar, nil
}

func (param2Struct *Param2Struct) field_2_list(prefixName string, item *gabs.Container) (string, error) {
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
		param2Struct.fields_2_obj(prefixName, objType, itemDesc, item)
		return objType, nil
	}

	return itemType, nil
}





