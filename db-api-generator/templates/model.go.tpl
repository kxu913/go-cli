package model
{{if and .ContainsArray .ContainsTime}}import (
	"time"
	"github.com/lib/pq"
){{else if .ContainsTime}}
import (
	"time"
)
{{else if .ContainsArray}}
import (
	"github.com/lib/pq"
)
{{ end }}

type {{.TableName}} struct { {{range $index,$Field := .Fields }} 
	{{$Field.FieldCaml}}	{{$Field.GoType}} {{ end }}
}
