{{- define "relationship_to_one_struct_helper" -}}
{{- end -}}

{{- $dot := . -}}
{{- $tableNameSingular := .Table.Name | singular -}}
{{- $modelName := $tableNameSingular | titleCase -}}
{{- $modelNameCamel := $tableNameSingular | camelCase -}}
// {{$modelName}} is an object representing the database table.
type {{$modelName}} struct {
	{{range $column := .Table.Columns -}}
	{{titleCase $column.Name}} {{$column.Type}} `{{generateTags $dot.Tags $column.Name}}boil:"{{$column.Name}}" json:"{{$column.Name}}{{if $column.Nullable}},omitempty{{end}}" toml:"{{$column.Name}}" yaml:"{{$column.Name}}{{if $column.Nullable}},omitempty{{end}}"`
	{{end -}}
	{{- if .Table.IsJoinTable -}}
	{{- else}}
	R *{{$modelNameCamel}}R `{{generateIgnoreTags $dot.Tags}}boil:"-" json:"-" toml:"-" yaml:"-"`
	L {{$modelNameCamel}}L `{{generateIgnoreTags $dot.Tags}}boil:"-" json:"-" toml:"-" yaml:"-"`
	{{end -}}

  readonly *{{$modelName}}
  whitelist []string
  operation string
}

var {{$modelName}}FieldMapping = map[string]string{
	{{range $column := .Table.Columns -}}
  "{{$column.Name}}": "{{titleCase $column.Name}}",
	{{end -}}
}

{{- if .Table.IsJoinTable -}}
{{- else}}
// {{$modelNameCamel}}R is where relationships are stored.
type {{$modelNameCamel}}R struct {
	{{range .Table.FKeys -}}
	{{- $txt := txtsFromFKey $dot.Tables $dot.Table . -}}
	{{$txt.Function.Name}} *{{$txt.ForeignTable.NameGo}}
	{{end -}}

	{{range .Table.ToOneRelationships -}}
	{{- $txt := txtsFromOneToOne $dot.Tables $dot.Table . -}}
	{{$txt.Function.Name}} *{{$txt.ForeignTable.NameGo}}
	{{end -}}

	{{range .Table.ToManyRelationships -}}
	{{- $txt := txtsFromToMany $dot.Tables $dot.Table . -}}
	{{$txt.Function.Name}} {{$txt.ForeignTable.Slice}}
	{{end -}}{{/* range tomany */}}
}

// {{$modelNameCamel}}L is where Load methods for each relationship are stored.
type {{$modelNameCamel}}L struct{}
{{end -}}
