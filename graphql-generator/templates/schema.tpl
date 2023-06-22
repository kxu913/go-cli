package graphql

import (
	"{{.ProjectName}}/db"

	"github.com/graphql-go/graphql"
)

var {{.QueryName}} = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Task",
		Fields: graphql.Fields{
			{{range $index,$Field := .Fields }} 
			"{{$Field.Name}}": &graphql.Field{
				Type: {{$Field.GraphqlType}},
			},{{- end }}
		},
	},
)

var {{.QueryName}}Cfg = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "{{.QueryName}}Cfg",
		Fields: graphql.Fields{

			"{{.QueryName}}": &graphql.Field{
				Type: graphql.NewList({{.QueryName}}),
				Description: "{{.QueryDescription}}",
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					return db.Query(`{{.SQL}}`), nil
				},
			},
		},
	},
)

var {{.QueryName}}Schema, _ = graphql.NewSchema(
	graphql.SchemaConfig{
		Query: {{.QueryName}}Cfg,
	},
)