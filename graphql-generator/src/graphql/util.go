package graphql

import (
	"fmt"
	"graphql-generator/db"
	"graphql-generator/model"
	"html/template"
	"io/ioutil"
	"os"
	"strings"
	"unicode"
)

var (
	TypeMapping = map[string]string{
		"int4":      "graphql.Int",
		"numeric":   "graphql.Float",
		"timestamp": "graphql.DateTime",
		"text":      "graphql.String",
	}
)

// keyword must use uppercase
var sql = `
	SELECT id, name AS Name, created_time AS CreatedTime 
	FROM note AS a INNER JOIN user AS b on a.ownerId=b.onwerId
`

var Output = GetEnvWithDefault("output", "d:/tmp/")

func GetDBType(name string, fields []model.Field) string {
	var dbType string
	var hName = name
	if strings.Contains(name, ".") {
		hName = strings.Trim(strings.Split(name, ".")[1], " ")
	} else {
		hName = strings.Trim(name, " ")
	}
	for _, field := range fields {

		if field.Name == hName {
			dbType = field.DBType
			break
		}
	}
	return dbType
}

func ParseSQL(query *model.Query) *model.Query {

	filedsStr := strings.Split(query.SQL, "FROM ")
	mainTable := strings.Split(filedsStr[1], " ")[0]
	// handle main table.
	dbFileds := []model.Field{}
	dbFileds = append(dbFileds, db.TableInfo(mainTable)...)
	// handle join table.
	joinTables := strings.Split(filedsStr[1], " JOIN ")
	if len(joinTables) > 1 {
		for _, t := range joinTables {
			jt := strings.Split(t, " ")[0]
			dbFileds = append(dbFileds, db.TableInfo(jt)...)
		}

	}

	sfields := strings.Split(strings.Trim(filedsStr[0], "SELECT"), ",")
	fields := []model.Field{}
	for _, a := range sfields {
		b := strings.Split(a, "AS")
		dbType := GetDBType(b[0], dbFileds)
		if len(b) > 1 {
			fields = append(fields, model.Field{
				Name:        strings.Trim(b[1], " "),
				Alias:       strings.Trim(b[1], " "),
				DBType:      dbType,
				GraphqlType: TypeMapping[dbType],
			})
		} else {
			fields = append(fields, model.Field{
				Name:        strings.Trim(b[0], " "),
				DBType:      dbType,
				GraphqlType: TypeMapping[dbType],
			})
		}
	}

	return &model.Query{
		ProjectName:      query.ProjectName,
		Prefix:           query.Prefix,
		QueryName:        Ucfirst(query.QueryName),
		SQL:              query.SQL,
		QueryDescription: query.QueryDescription,
		Fields:           fields,
	}

}

func CreateModFile(project string) {
	tpl, _ := template.ParseFiles("templates/go.mod.tpl")
	file := WriteFile(fmt.Sprintf("%s%s/src/graphql/go.mod", Output, project))
	defer file.Close()

	tpl.Execute(file, map[string]any{
		"ProjectName": project,
	})
}

func CreateBasicFile(project string) {
	tpl, _ := template.ParseFiles("templates/basic.go.tpl")
	file := WriteFile(fmt.Sprintf("%s%s/src/graphql/basic.go", Output, project))
	defer file.Close()

	tpl.Execute(file, map[string]any{
		"ProjectName": project,
	})
}

func CreateGraphqlDefiniton(project string, query *model.Query) {
	gFolder := fmt.Sprintf("%s%s/src/graphql", Output, project)
	_, err := os.Stat(gFolder)
	if os.IsNotExist(err) {
		os.MkdirAll(gFolder, 0755)
		CreateModFile(project)
		CreateBasicFile(project)
	}
	tpl, _ := template.ParseFiles("templates/schema.tpl")
	file := WriteFile(fmt.Sprintf("%s%s/src/graphql/%s_schema.go", Output, project, strings.ToLower(query.QueryName)))
	defer file.Close()

	tpl.Execute(file, query)
}

func GenerateRouterFunc(project string, query *model.Query) {
	tpl, _ := template.ParseFiles("templates/router.tpl")

	file := WriteFile(fmt.Sprintf("%s%s/src/route/ghql_%s.go", Output, project, strings.ToLower(query.QueryName)))
	defer file.Close()
	tpl.Execute(file, map[string]any{
		"ProjectName": project,
		"RouterName":  strings.ToLower(query.QueryName),
		"QueryName":   query.QueryName,
		"Prefix":      query.Prefix,
	})
}

func ReplaceMain(project string, query *model.Query) {
	mainName := fmt.Sprintf("%s%s/main.go", Output, project)
	b, err := ioutil.ReadFile(mainName)
	if err != nil {
		panic(err)
	}
	lineBytes := string(b)
	lines := strings.Split(lineBytes, "\n")
	var newLines []string
	for _, l := range lines {
		newLines = append(newLines, l)
		if strings.Index(l, "//auto generate groups") >= 0 {
			newLines = append(newLines, fmt.Sprintf("	route.%sGroupGraphql(e)\n", Ucfirst(query.QueryName)))
		}
	}
	file, err := os.OpenFile(mainName, os.O_WRONLY, 0666)
	defer file.Close()
	file.WriteString(strings.Join(newLines, "\n"))

}

func DeleteFiles(project string, api string) {

	if err := os.Remove(fmt.Sprintf("%s%s/src/route/ghql_%s.go", Output, project, strings.ToLower(api))); err != nil {
		panic(err)
	}

	if err := os.Remove(fmt.Sprintf("%s%s/src/graphql/%s_schema.go", Output, project, strings.ToLower(api))); err != nil {
		panic(err)
	}

}

func UpdateMain(project string, api string) {
	mainName := fmt.Sprintf("%s%s/main.go", Output, project)
	b, err := ioutil.ReadFile(mainName)
	if err != nil {
		panic(err)
	}
	content := string(b)
	newContent := strings.Replace(content, fmt.Sprintf("	route.%sGroupGraphql(e)\n", Ucfirst(api)), "", 1)
	fmt.Println(newContent)

	if err := os.Remove(mainName); err != nil {
		panic(err)
	}

	file := WriteFile(mainName)
	defer file.Close()
	file.WriteString(newContent)

}

// 首字母大写
func Ucfirst(str string) string {
	for i, v := range str {
		return string(unicode.ToUpper(v)) + str[i+1:]
	}
	return ""
}

func GetEnvWithDefault(key string, defValue string) string {
	val, err := os.LookupEnv(key)
	if !err {
		return defValue
	}
	return val

}

func WriteFile(fileName string) *os.File {
	os.OpenFile(fileName, os.O_CREATE, 0o666)
	file, err := os.OpenFile(fileName, os.O_RDWR, 0o666)
	if err != nil {
		panic(err)
	}
	return file
}
