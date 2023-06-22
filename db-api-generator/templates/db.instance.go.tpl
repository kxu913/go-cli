package db

import (

	"{{.ProjectName}}/model"
)

{{- $m :=.TableInfo.LastFieldIndex}}{{- $db :=.TableInfo.DBTable}}

const Selected{{.TableInfo.TableName}}Fields = `{{range $index,$Field := .TableInfo.Fields }} 
	{{ if lt $index $m}}{{ $db}}.{{$Field.Field}} as "{{$Field.FieldCaml}}",{{ else }}{{ $db}}.{{$Field.Field}} as "{{$Field.FieldCaml}}"{{- end }}
	{{- end }}
`

func Create{{.TableInfo.TableName}}(m *model.{{.TableInfo.TableName}}) *model.{{.TableInfo.TableName}} {
	g := OpenDBByGorm()
	defer Close(g)
	g.Table("{{.TableInfo.DBTable}}").Create(m)
	return m
}

func Get{{.TableInfo.TableName}}(id int) *model.{{.TableInfo.TableName}} {
	g := OpenDBByGorm()
	defer Close(g)
	rtn := &model.{{.TableInfo.TableName}}{}
	err := g.Table("{{.TableInfo.DBTable}}").Where("id=?", id).Take(rtn).Error
	if err != nil {
		return nil
	}
	return rtn
}


func Get{{.TableInfo.TableName}}s() []model.{{.TableInfo.TableName}} {
	g := OpenDBByGorm()
	defer Close(g)
	rtn := []model.{{.TableInfo.TableName}}{}
	g.Table("{{.TableInfo.DBTable}}").Find(&rtn)
	return rtn
}

