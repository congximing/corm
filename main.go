package corm

import (
	"database/sql"
	"fmt"
	"log"
	"reflect"
	"strings"
)

type DB struct{
	Parent *sql.DB
	//Module interface{}
	DbName string
	//<k,map>
	DbKV map[string]map[string]string //不安全
	flage bool
	stmt	*sql.Stmt
	sql		Sql
	entity Entity //存放结构体的名称与值==>where

	structBody interface{}//存放结构体
	args interface{}
	structList	[]interface{}
}


//"root:324130@tcp(127.0.0.1)/book?charset=utf8&multiStatements=true")"
//todo 创建连接
func Connect(dialect string,user string,pwd string,host string,dbName string)( *DB, error){
	args := user+":"+pwd+"@tcp("+host+")/"+dbName+"?charset=utf8&multiStatements=true"
	dbSQL,err :=sql.Open(dialect,args)
	db :=&DB{
		Parent: dbSQL,
		DbName: dbName,
	}
	return db,err
}

//---------------------------dml------------------------------------
//todo 插入
//INSERT INTO user (id, name, age) VALUES (?, ?, ?)
//写一个含有主键的struct，检查struct是否有主键，尽量不要从DB查
func (db *DB)Insert(arg interface{})(interface{}){
	db.Where(arg)
	sql := "INSERT INTO "+db.sql.From+" "
	//todo (k,k,k,k) values(?,?,?,?) 在这里判断主键，没必要，插入必须要加主键，抛个异常
	//todo 暂时不支持，不含有主键的insert操作,并且数据没有约束的问题。
	insertCondition := InsertCondition(db.entity.aType,db.entity.aValue)
	sql += insertCondition
	fmt.Println("insert sql : "+sql)// INSERT INTO user (id,name,age) VALUES(?,?,?)
	stmt, err := db.Parent.Prepare(sql)
	defer stmt.Close()
	if err!=nil {
		log.Fatalf("insert data error: %v\n", err)
		return db
	}
	//todo 插入参数 stmt.Query()
	rows, err := Query(stmt, db)
	fmt.Println()
	defer rows.Close()
	if err!=nil {
		log.Fatalf("insert data error: %v\n", err)
		return db
	}

	return db
}
//todo 删除
//db.Delete("DELETE FROM user WHERE id = ?")
func (db*DB)Delete()(interface{}){
	sql := "Delete from "+db.sql.From+" where "
	condtition := PutCondition(db.entity.aType,db.entity.aValue)
	sql += condtition
	fmt.Println(sql)
	stmt, err := db.Parent.Prepare(sql)
	defer stmt.Close()
	rows, err := Query(stmt, db)
	defer rows.Close()
	if err != nil {
		log.Fatalf("delete data error : %v\n",err)
		return db
	}
	if err!=nil{
		log.Fatalf("query delete data failed: %v\n",err)
		return db
	}

	return db
}

//todo 修改
//todo update user set k=v where k=?,设置字段没有预编译
func (db *DB)Update(args interface{})(interface{}){
	//todo 1.prepare
	sql := "update "+db.sql.From+" set "
	condition := PutCondition(db.entity.aType, db.entity.aValue)
	//todo set k=v
	set := UpdateCondition(args)
	sql = sql+set+"where 1=1 "+condition
	stmt, err := db.Parent.Prepare(sql)
	//log.Fatal(sql)
	defer stmt.Close()
	if err!=nil{
		log.Fatal("Update prepare failed")
		return db
	}
	//todo 2.exec
	_, err = Exec(stmt, db)
	if err!=nil{
		log.Fatal("update failed")
		return db
	}
	log.Fatal("update success")
	return db
}

//"select *from user where id =?"
//查询字段
func (db *DB)Select(args ...string)(interface{}){
	//todo 1.预编译SQL
	if len(args) == 0 || args[0]==""{
		db.sql.Select="select *"
		//args[0]="*" //todo error
	}else{
		flage := false
		for _,s:=range args{
			if !flage{
				db.sql.Select="select "+s+" "
				//str[i]=s
				flage=true
				continue
			}
			db.sql.Select+=","+s
		}
	}
	sql := db.sql.Select+" from "+db.sql.From+" where 1=1 "
	//todo where
	condition := PutCondition(db.entity.aType, db.entity.aValue)
	sql +=condition
	//todo limit
	if db.sql.Limit!=""{
		sql += " limit "+db.sql.Limit
	}
	stmt, err := db.Parent.Prepare(sql)
	defer stmt.Close()
	fmt.Println(sql)
	if err != nil {
		log.Fatal("select data error : %v\n",stmt)
		return db
	}
	//todo 2.参数注入并执行
	rows, err := Query(stmt,db)
	//fmt.Println(rows)
	defer rows.Close()
	if err != nil {
		log.Fatal("select query data error :%v\n",rows)
		return db
	}
	//todo 3. 查询字段
	db.args=args
	PrintRows(rows,db,args)
	return db
}

//todo 根据主键，自增长
func (db *DB)autoIncreatment(id string)int{
	stmt, err := db.Parent.Prepare("select ? from user order by id desc limit 1")
	if err!=nil{
		log.Fatal("none data")
		return 1
	}
	rows, err := stmt.Query(id)
	var i int
	rows.Scan(&i)
	return i
}
//todo 根据字段查询每一个数据
func PrintRows(rows *sql.Rows,db *DB,str []string){
	//todo 查看表结构
	PutStructTypeName(db.structBody)
	//根据参数个数来返回对应的方法。
	l :=reflect.TypeOf(db.structBody).NumField()
	Scan(rows,str,l,db)
}

//---------------------------ddl------------------------------------
//todo 可以创建多个表
func (db*DB)CreateTable(modules  ...interface{})(*DB){
	for _,module := range modules{
		//1. 解析结构创建表
		//	表名{
		//		列名		类型
		//	}
		clzz := reflect.TypeOf(module)//class info
		var len int= clzz.NumField()
		var columnName  [20]string
		var columnType  [20]string
		for i :=0;i<len;i++{
			field := clzz.Field(i)
			//get columnName
			columnName[i]=field.Name
			//get columnType
			columnType[i]=field.Type.Name()
		}
		db.From(module)
		//2. insert sql
		//var tableName  string =
		sql := "DROP TABLE IF EXISTS "+db.sql.From+";" +
			"CREATE TABLE "+db.sql.From+" ("
		for i:=0;i<len;i++{
			//"  `user_id` int(11) ," +
			//" `name` varchar(30) ," +
			if i==len-1{
				t := db.StructToDb(columnType[i])
				sql +=columnName[i]+" "+t
			}else {
				t := db.StructToDb(columnType[i])
				sql +=columnName[i]+" "+t+","
			}
		}
		sql +=") ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='产品表';"
		fmt.Println(sql+" 157")
		//3. 创建表
		db.Parent.Exec(sql)

	}
	return db
}

//todo DropTable
func (db*DB)DropTable(modules ...interface{})(*DB){
	//从module中提取表名
	for _,module := range modules{
		str := ToGetName(module)
		db.Parent.Exec("DROP TABLE "+str)
	}
	return db
}
//todo UpdateTable
func (db*DB)UpdateTable()(*DB){

	return db
}
//todo 数据库的列名和类型
func (db*DB)GetTableInfo(modules ...interface{})(*DB){
	for _,module := range modules{
		table_name := ToGetName(module)
		str :="SELECT COLUMN_NAME fName,DATA_TYPE dataType FROM information_schema.columns  WHERE table_schema = ? AND table_name = ?; "
		stmt, err := db.Parent.Prepare(str)
		defer stmt.Close()
		if err!=nil{
			log.Fatal("prepare: desc table error: ",err)
			return db
		}
		rows, err := stmt.Query(db.DbName,table_name)
		defer rows.Close()
		if err != nil {
			log.Fatal("row: desc table error:  ",err)
			return db
		}
		if err!=nil{

		}

		db.putKV(rows,table_name)
	}

	return db
}
//store  <table,<fName,dataType>>
func (db *DB)putKV(rows *sql.Rows,table_name string){
	var a,b string
	var dbMap map[string]string
	dbMap=make(map[string]string)
	for rows.Next(){
		rows.Scan(&a,&b)
		fmt.Println(a+": "+b)
		if a != " " {
			dbMap[a]=b
		}

	}
	db.initDbKV()
	db.DbKV[table_name]=dbMap
}
func (db*DB)initDbKV(){
	if db.flage==false{
		db.DbKV=make(map[string]map[string]string)
	}
}

//store <table,list>
//todo mysql ==> json,
//todo 给结构体设置数据没什么用，主要是给前端传数据
func (db*DB)MysqlToJson(module interface{})(string){
	db.GetTableInfo(module)
	struct_name := ToGetName(module)
	//根据db中的字段进行生成json
	json := db.Map2Json(db.DbKV[struct_name])
	//返回json
	return json
}

//todo 输入 User{}，"User"
//todo 输出 table_name
func (db *DB)From(arg interface{})(*DB){
	table_name := strings.ToLower(ToGetName(arg))
	db.sql.From=table_name
	db.structBody=arg
	return db
}

//todo where Field
func (db *DB)Where(stru interface{})(*DB){
	if s,ok := stru.(string);ok{
		if s=="all"{
			return db
		}
		log.Fatalf("input error ,please input \"all\" to select all rows!")
		return db
	}
	db.GetStructTypeValue(stru)
	return db
}

func (db *DB)Limit(arg string)(*DB){
	db.sql.Limit=arg
	return db
}
func (db *DB)Offset()(*DB){
	return db
}

func (db *DB)Count()(*DB){
	return db
}

func (db *DB)OrderBy()(*DB){
	return db
}

func (db *DB)GroupBy()(*DB){
	return db
}

func (db *DB)Having()(*DB){
	return db
}

