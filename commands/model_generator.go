package commands

import (
	"flag"
	"fmt"
	"github.com/astaxie/beego/orm"
	"os"
	"path/filepath"
	"strings"
)

//ModelGeneratorCommand 根据连接的表信息生成model
func ModelGeneratorCommand() {
	if len(os.Args) < 2 || os.Args[1] != "model_generator" {
		return
	}
	defer func() {
		if err := recover(); err != nil {
			panic(err)
		}
		os.Exit(0)
	}()

	flag.String("table", "", "指定某张表")
	flag.Bool("force", false, "是否强制执行")
	flag.String("db", "", "指定库名")

	flag.CommandLine.Parse(os.Args[2:])
	giveTable := flag.Lookup("table").Value.String()
	force := flag.Lookup("force").Value.String()
	dbName := flag.Lookup("db").Value.String()

	//显示所有的表名
	o := orm.NewOrm()
	if dbName == "" {
		o.Raw("select database()").QueryRow(&dbName)
	}

	tables := make([]string, 0)
	db, _ := orm.GetDB("default")
	_, err := db.Exec("use "+dbName)
	if err != nil {
		panic(err)
	}

	_, err = o.Raw("show tables").QueryRows(&tables)
	if err != nil {
		panic(err)
	}

	//遍历表明，获取表结构
	for _,table := range tables {
		if giveTable != "" && table != giveTable {
			continue
		}

		if force != "true" {
			//判断表结构对应的model文件是否存在
			_, err := os.Stat(fmt.Sprintf("./models/%s/%s.go", dbName, table))
			if err == nil {
				fmt.Println("exist:", table)
				continue
			}
		}

		//获取表字段类型生成model文件
		fmt.Println("generator:", table)
		generatorFromTable(dbName, table)
	}

	//不存在则生成文件
}

func generatorFromTable(dbName, table string)  {
	o := orm.NewOrm()

	columns := make([]orm.Params, 0)
	_, err := o.Raw("SELECT TABLE_CATALOG,TABLE_SCHEMA,TABLE_NAME,COLUMN_NAME,ORDINAL_POSITION,COLUMN_DEFAULT,IS_NULLABLE,DATA_TYPE,CHARACTER_MAXIMUM_LENGTH,CHARACTER_OCTET_LENGTH,NUMERIC_PRECISION,NUMERIC_SCALE,DATETIME_PRECISION,CHARACTER_SET_NAME,COLLATION_NAME,COLUMN_TYPE,COLUMN_KEY,EXTRA,PRIVILEGES,COLUMN_COMMENT,GENERATION_EXPRESSION,SRS_ID FROM INFORMATION_SCHEMA.COLUMNS WHERE table_name = ? and TABLE_SCHEMA=? ORDER BY ORDINAL_POSITION asc", table, dbName).Values(&columns)
	if err != nil {
		panic(err)
	}
	//显示列信息
	imports := make(map[string]int)
	imports["github.com/astaxie/beego/orm"] = 1
	//相关列信息 name,type,size,auto,pk,comment
	orms := make([]map[string]interface{}, 0)
	for _, column := range columns {
		tmp := make(map[string]interface{})
		//获取字段名
		tmp["name"] = column["COLUMN_NAME"].(string)
		//获取类型和尺寸
		switch column["DATA_TYPE"].(string) {
		case "varchar":
			tmp["type"] = "string"
			tmp["size"] = column["CHARACTER_MAXIMUM_LENGTH"]
		case "bigint":
			tmp["type"] = "int64"
		case "int":
			tmp["type"] = "int64"
		case "mediumint":
			tmp["type"] = "int64"
		case "tinyint":
			tmp["type"] = "int"
		case "smallint":
			tmp["type"] = "int"
		case "float":
			tmp["type"] = "float64"
		case "datetime":
			tmp["type"] = "time.Time"
			imports["time"] = 1
		default:
			if strings.Index(column["DATA_TYPE"].(string), "int") != -1 {
				tmp["type"] = "int64"
			} else {
				tmp["type"] = "string"
			}
		}
		//判断是否非负
		if strings.Index(column["COLUMN_TYPE"].(string), "unsigned") > -1 && tmp["type"] != "string" {
			tmp["type"] = "u"+tmp["type"].(string)
		}

		//是否自增
		if column["EXTRA"].(string) == "auto_increment" {
			tmp["auto"] = 1
		}
		//是否主键
		if column["COLUMN_KEY"].(string) == "PRI" {
			tmp["pk"] = 1
		}
		//注释
		tmp["comment"] = column["COLUMN_COMMENT"]

		orms = append(orms, tmp)
	}

	//组织代码:需要的包
	codes := fmt.Sprintf("package %s\n\n", dbName)
	if len(imports) > 0 {
		codes += "import ( \n"
		for k,_ := range imports  {
			codes +=fmt.Sprintf("\t\"%s\"\n", k)
		}
		codes +=")\n\n"
	}

	//组织代码：定义models
	tableUp := toUpword(table)
	codes += fmt.Sprintf("type %s struct {\n", tableUp)
	if len(orms) > 0 {
		for _,v := range orms {
			//相关列信息 name,type,size,auto,pk,comment
			zhujie := ""
			tmp := ""
			if v["size"] != nil {
				tmp += "size("+v["size"].(string)+");"
			}
			if v["auto"] != nil {
				tmp += "auto;"
			}
			if v["pk"] != nil {
				tmp += "pk;"
			}
			if v["type"] == "time.Time" {
				tmp += "type(datetime);"
			}
			if len(tmp) > 0 {
				tmp = strings.TrimRight(tmp, ";")
				zhujie += fmt.Sprintf(" orm:\"%s\"", tmp)
			}
			if v["comment"] != nil && v["comment"].(string) != "" {
				zhujie += fmt.Sprintf(" description:\"%s\"", v["comment"])
			}
			upName := toUpword(v["name"].(string))
			if upName == "TableName" {
				upName = "TableNames"
			}
			codes +=fmt.Sprintf("\t%s %s `json:\"%s\"%s`\n", upName, v["type"], v["name"], zhujie)
		}
	}
	codes += "}\n\n"

	//注册表结构
	codes += "func init() {\n"
	codes += fmt.Sprintf("\torm.RegisterModel(new(%s))\n}\n\n", tableUp)

	//定义表名
	codes += fmt.Sprintf("func (this *%s) TableName() string {\n", tableUp)
	codes += fmt.Sprintf("\treturn \"%s.%s\"\n}\n\n", dbName, table)

	//写入文件
	fp := fmt.Sprintf("./models/%s/%s.go", dbName, table)
	err = os.MkdirAll(filepath.Dir(fp), 0666)
	if err != nil {
		panic(err)
	}
	f, err := os.Create(fp)
	if err != nil {
		panic(err)
	}
	_, err = f.WriteString(codes)
	if err != nil {
		panic(err)
	}
	f.Close()
}

func toUpword(word string) string {
	tbs := strings.Split(word, "_")
	tableUp := ""
	for _,v := range tbs {
		rv := []rune(v)
		tmp := strings.ToUpper(string(rv[0]))+string(rv[1:])
		tableUp += tmp
	}
	return tableUp
}