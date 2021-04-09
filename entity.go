package corm

import "reflect"

const LEN int = 5
type Entity struct{
	aType [LEN]string
	aValue [LEN]reflect.Value
}
