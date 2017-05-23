{{define "domain.search"}}
{{ $titleName := titleCase .Name }}
{{- $respSpecName := printf "%sList" .Name -}}
{{- $domainName := $titleName -}}
{{- $specPkg := pkgName .SpecPkg -}}
{{- $modelPkg := pkgName .ModelPkg -}}
{{- $domainPkg := pkgName .DomainPkg -}}
{{- $primaryModel := titleCase .PrimaryModel -}}
{{ $titleName := titleCase .Name }}
{{- $sliceDomainName := plural $titleName -}}

{{ template "domain.usage.search.func" $sliceDomainName }}
func (d *{{$domainName}}) Search(ctx ctx.TaskContext) (slice *{{$sliceDomainName}}, errs []*errors.Error) {
  return
}

// TODO: get the wrapped slice object by ids
func (d *{{$domainName}}) listByIDs(exec boil.Executor, ids ...interface{}) (o *{{$sliceDomainName}}, err *errors.Error) {
  return
}
{{end}}
