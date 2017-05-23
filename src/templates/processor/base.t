{{ define "base" }}
{{ $procName := . -}}

func New{{$procName}}(db *sql.DB) Processor {
	return &{{$procName}}{}
}

type {{$procName}} struct{}

func (c *{{$procName}}) Name() string {
	return "{{$procName}}"
}
{{ end}}
