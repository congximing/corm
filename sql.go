package corm

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
)

type Sql struct{
	Select string
	From string
	Where string
	Limit string
	Offset string
	Count string
}
//-------------------Query-----------------------------------------
//todo 对sql的Query()进行封装
func Query(stmt *sql.Stmt,db *DB)( *sql.Rows, error){
	var rows *sql.Rows
	var err error
	aValue := db.entity.aValue//value

	//todo db.entity.aValue没取到值
	var args [LEN]string
	for i,s :=range aValue{
		s := RefValue(s)
		args[i]=s
		//fmt.Println(s+"====")
	}
	size := Size(args)
	//fmt.Println(size)
	//参数匹配
	rows, err = selectQuery(size, stmt, args)

	return rows,err
}
func selectQuery(size int,stmt *sql.Stmt,args [LEN]string)( *sql.Rows, error){
	var rows *sql.Rows
	var err error
	if size==1{
		for i,s := range args{
			if s == ""{
				continue
			}
			rows, err = query1(args[i], stmt)
		}
	}else if size==2{
		rows, err = query2(stmt,args[0],args[1])
	}else if size==3{
		rows, err = query3(stmt, args[0], args[1], args[2])
	}else if size==4{
		rows, err = query4(stmt,args[0],args[1],args[2],args[3])
	}else if size ==0{
		rows, err = stmt.Query()
	}
	return rows,err
}

func query1(args string,stmt *sql.Stmt)( *sql.Rows, error){
	rows, err := stmt.Query(args)
	if err != nil {
		fmt.Println("query1 error")
	}
	return rows,err
}
func query2(stmt *sql.Stmt,args ...string)( *sql.Rows, error){
	rows, err := stmt.Query(args[0],args[1])
	return rows,err
}
func query3(stmt *sql.Stmt,args ...string)( *sql.Rows, error){
	rows, err := stmt.Query(args[0],args[1],args[2])
	return rows,err
}
func query4(stmt *sql.Stmt,args ...string)( *sql.Rows, error){
	rows, err := stmt.Query(args[0],args[1],args[2],args[3])
	return rows,err
}
//---------------------------------------------------------------
//-------------------Exec-----------------------------------------
//todo 对sql的Exec()进行封装
func Exec(stmt *sql.Stmt,db *DB)( sql.Result, error){
	var res sql.Result
	var err error
	aValue := db.entity.aValue//value
	var args [LEN]string
	for i,s :=range aValue{
		s := RefValue(s)
		args[i]=s
	}
	size := Size(args)
	//参数匹配
	res, err = selectExec(size, stmt, args)
	return res,err
}
func selectExec(size int,stmt *sql.Stmt,args [LEN]string)( sql.Result, error){
	var res sql.Result
	var err error
	if size==1{
		args := Remove(args)
		res, err = Exec1(args[0], stmt)
	}else if size==2{
		args := Remove(args)
		//fmt.Println(args[0]+" "+args[1]+" size:2")
		res, err = Exec2(stmt,args[0],args[1])
	}else if size==3{
		args := Remove(args)
		res, err = Exec3(stmt, args[0], args[1], args[2])
	}else if size==4{
		args := Remove(args)
		res, err = Exec4(stmt,args[0],args[1],args[2],args[3])
	}
	return res,err
}

func Exec1(args string,stmt *sql.Stmt)( sql.Result, error){
	res, err := stmt.Exec(args)
	if err != nil {
		fmt.Println("query1 error")
	}
	return res,err
}
func Exec2(stmt *sql.Stmt,args ...string)( sql.Result, error){
	res, err := stmt.Exec(args[0],args[1])
	return res,err
}
func Exec3(stmt *sql.Stmt,args ...string)( sql.Result, error){
	res, err := stmt.Exec(args[0],args[1],args[2])
	return res,err
}
func Exec4(stmt *sql.Stmt,args ...string)( sql.Result, error){
	res, err := stmt.Exec(args[0],args[1],args[2],args[3])
	return res,err
}
//---------------------------------------------------------------
//----------------------------select------------------------------
//todo 根据y返回特定长度
func appendString(x []string,y int)[]string{
	x=make([]string,y,y)
	return x
}
//todo read data from db ===> Json
func Scan(rows *sql.Rows,args []string,ll int,db *DB)(string){
	var jstr string
	l :=len(args)
	//select *
	if len(args)==0 || args[0]=="" || args[0]=="*"{
		l=ll
		//todo 得到结构体
		db.GetStructTypeValue(db.structBody)
		aType := db.entity.aType
		//todo [LEN]string => []string 类型转换 使用slice
		//todo 扩容
		args = appendString(args, len(aType))
		copy(args[:],aType[0:len(aType)])
		db.args=args
		//fmt.Println(args)
	}
	if rows.Next(){
		//fmt.Println("=========")
		if l==1 {
			db.Scan1(rows,args[0])
		}else if l == 2 {
			db.Scan2(rows,args[0],args[1])
		}else if l==3 {
			//todo index out of range [1] with length 1
			db.Scan3(rows,args[0],args[1],args[2])
		}else if l==4 {
			db.Scan4(rows,args[0],args[1],args[2],args[3])
		}
	}
	return jstr
}

func (db*DB)Scan1(rows *sql.Rows,s ...string){
	for rows.Next() {
		rows.Scan(&s[0])
		//fmt.Println("根据查询字段： "+s[0])
		entity := DataToStruct(db,s)
		fmt.Printf("entity: %v\n",entity)
	}
}
func (db*DB)Scan2(rows *sql.Rows,args ... string){
	for rows.Next() {
		rows.Scan(&args[0],&args[1])
		//fmt.Println("根据查询字段： "+args[0]+args[1])
		entity := DataToStruct(db,args)
		fmt.Printf("entity: %v\n",entity)
	}
}
func (db*DB)Scan3(rows *sql.Rows,args ... string){
	rows.Scan(&args[0],&args[1],&args[2])
	//fmt.Println("根据查询字段 得到的数据为： "+args[0]+args[1]+args[2])
	entity := DataToStruct(db,args)
	fmt.Printf("entity: %v\n",entity)
	for rows.Next(){
		rows.Scan(&args[0],&args[1],&args[2])
		//fmt.Println("根据查询字段 得到的数据为： "+args[0]+args[1]+args[2])
		entity := DataToStruct(db,args)
		fmt.Printf("entity: %v\n",entity)
	}
}
func (db*DB)Scan4(rows *sql.Rows,args ...string){
	for rows.Next() {
		rows.Scan(&args[0],&args[1],&args[2],&args[3])
		//fmt.Println("根据查询字段： "+args[0]+args[1]+args[2]+args[3])
	}
}

//todo where k=? and k=? and k=?
func PutCondition(aType interface{},aValue interface{})string{
	var condition string
	var flage =false
	//k=? and k=? and k=?
	if aType == nil{
		return ""
	}
	strV := aValue.([LEN]reflect.Value)
	strT := aType.([LEN]string)

	for i,v := range strV{
		s := RefValue(v)
		if s != ""{
			if !flage {
				condition += "and "+strings.ToLower(strT[i])+"=? "//变小写
				flage=true
			}else{
				condition += "and "+strings.ToLower(strT[i])+"=? "
			}
		}
	}
	return condition
}
//todo (k,k,k,k) values(?,?,?,?)
func InsertCondition(aType interface{},aValue interface{})string{
	var key string
	var value string
	flage := false
	key = "("
	value = "VALUES("

	strV := aValue.([LEN]reflect.Value)
	strT := aType.([LEN]string)
	for i,v := range strV{
		s := RefValue(v)
		if s != ""{
			if flage!=true{
				key += strings.ToLower(strT[i])
				value += "?"
				flage=true
			}else{
				value += ",?"
				key += ","+strings.ToLower(strT[i])
			}
		}
	}
	key += ")"
	value += ")"
	condition := key +" "+ value
	return condition
}
//todo k=v,k=v,k=v
func UpdateCondition(args interface{})string{
	//todo set k=v
	sType := reflect.TypeOf(args) //string
	sValue := reflect.ValueOf(args)	//reflect.Value
	if sValue.Kind() != reflect.Struct {
		panic("need struct kind")
	}
	typeLen := sType.NumField()
	var k [LEN]string
	var v [LEN]reflect.Value

	for i:=0;i<typeLen;i++{
		s := sValue.Field(i)
		if s.String() != "<invalid Value>" {
			k[i] = strings.ToLower(sType.Field(i).Name)
			v[i] = s
		}
	}

	var set string
	flage := false
	for i:=0;i<LEN;i++{
		if RefValue(v[i])!=""{
			if !flage{
				set +=k[i]+"=\""+RefValue(v[i])+"\" "
				flage=true
				continue
			}
			set +=","+k[i]+"=\""+RefValue(v[i])+"\" "
		}
	}
	return set
}