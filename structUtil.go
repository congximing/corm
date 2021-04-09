package corm

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

//todo 输入：结构体
//todo 返回：成员名，标签值
func (db*DB)GetStructTypeValue(arg interface{})(interface{},interface{}){
	sType := reflect.TypeOf(arg) //string
	sValue := reflect.ValueOf(arg)	//reflect.Value
	if sValue.Kind() != reflect.Struct {
		panic("need struct kind")
	}
	typeLen := sType.NumField()
	//todo 这个最好不要声明为 interface
	aType :=&db.entity.aType
	aValue :=&db.entity.aValue

	for i:=0;i<typeLen;i++{
		s := sValue.Field(i)
		if s.String() != "<invalid Value>" {
			aType[i] = strings.ToLower(sType.Field(i).Name)
			aValue[i] = s
		}
	}
	return *aType,*aValue
}
//todo 输入 结构体标签值
//todo 返回 实际的值(string)
func RefValue(v reflect.Value) string{
	s := v.String()
	if s != "<invalid Value>" {
		if s == "<int32 Value>"{
			i32 := v.Interface().(int32)
			if i32==0{
				return ""
			}
			return strconv.FormatInt(int64(i32),10)
		}else if s == "<int16 Value>"{
			i16 := v.Interface().(int16)
			if i16==0{
				return ""
			}
			return strconv.FormatInt(int64(i16),10)
		}else if s == "<int8 Value>"{
			i8 := v.Interface().(int8)
			if i8 == 0{
				return ""
			}
			return strconv.FormatInt(int64(i8),10)
		}
		return s
	}
	return ""
}
//todo 输入 结构体
//todo 输出 结构体的所有成员类型
func GetStructType(arg interface{})(interface{}){
	sType := reflect.TypeOf(arg)
	len := sType.NumField()
	var s [LEN]string
	for i:=0 ;i<len ;i++{
		s[i] = sType.Field(i).Type.Name()
		fmt.Println(sType.Field(i).Name)
	}
	return s
}
//todo 输入 结构体
//todo 输出 结构体的所有成员名
func GetStructName(arg interface{})(interface{}){
	sType := reflect.TypeOf(arg)
	var s [LEN]string
	len := sType.NumField()
	for i:=0 ;i<len ;i++{
		s[i] = sType.Field(i).Name
	}
	return s
}
//todo 输入 结构体
//todo 输出 结构体的所有 类型：成员名
func PutStructTypeName(arg interface{})(map[string]string){
	typeName := make(map[string]string)
	sType := reflect.TypeOf(arg)
	len := sType.NumField()
	for i:=0 ;i<len ;i++{
		typeName[sType.Field(i).Name]=sType.Field(i).Type.Name()
	}
	return typeName
}

//todo data2struct 数据注入结构体
//todo 如果没有查到的字段，就不注入
func DataToStruct(db *DB,args []string)interface{}{
	//struct_entity
	t := reflect.ValueOf(db.structBody).Type()
	entity := reflect.New(t).Elem()
	//todo 根据类型和名称，插入具体的值
	structType := reflect.TypeOf(db.structBody)
	//all field
	length := structType.NumField()// 整个结构体字段的长度
	s := len(args)// 存储数据的变量
	str := db.args.([]string)//字段名 长度
	//fmt.Println(len(str))
	for i:=0;i<s;i++{
		fieldName := ToUpper(str[i])
		structName,_ := structType.FieldByName(fieldName)
		fieldType := structName.Type.Name()
			//fmt.Println(fieldName)
			//int32 string int16
			if fieldType == "string"{
				if i >= s && s<length{
					entity.FieldByName(fieldName).SetString("nil")
					continue
				}
				entity.FieldByName(fieldName).SetString(args[i])
			}else if fieldType == "int32" {
				if i >= s && s<length{
					entity.FieldByName(fieldName).SetInt(0)
					continue
				}
				i, _ := strconv.ParseInt(args[i], 10, 32)
				entity.FieldByName(fieldName).SetInt(i)
			}else if fieldType == "int16"{
				if i >= s && s<length{
					entity.FieldByName(fieldName).SetInt(0)
					continue
				}
				i, _ := strconv.ParseInt(args[i], 10, 16)
				entity.FieldByName(fieldName).SetInt(i)
			}

	}
	//todo 打印entity
	return entity
}
//todo 获取结构体的字符串名称
//todo 以后还要加入其他的类型
func ToGetName(module interface{})string{
	typeStr := reflect.TypeOf(module)
	if str,ok := module.(string);ok{
		return str
	} else {
		return typeStr.Name()
	}
}
