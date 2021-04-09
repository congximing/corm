package corm

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
)
//功能：数据库与struct之间的类型转换
//todo 字段类型的转换
//todo 将struct转换成数据库类型
func (db*DB)StructToDb(str string)string{
	if str=="uint"{
		str="int"
	}else if str=="float32"{
		str="float"
	}else if str == "string"{

	}else if str=="int16"{
		str="int"
	}else if str=="int8"{
		str="int"
	}else if str=="int64"{
		str="bigint"
	}

	return str
}

//todo 首字母大写
func ToUpper(str string)string{
	s := strings.ToUpper(str[:1])
	s2 :=str[1:len(str)]
	return s+s2
}

//---------------------json convert----------------------------
//todo map2json
func (db *DB)Map2Json(mapKV map[string]string)(string){
	var s string
	s="{"
	for k,v := range mapKV{
		s+="\""+k+"\":"+"\""+v+"\","
	}
	s=s[:len(s)-1]
	s+="}"
	fmt.Print(s)
	return s
}

//todo struct2json
func (db *DB)Struct2Json(arg interface{})string{
	marshal, err := json.Marshal(arg)
	if err != nil {
		log.Fatal("struct to json failed! #typeConv.go")
	}
	return string(marshal)
}

