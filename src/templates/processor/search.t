{{- $titleName := titleCase .Name -}}
{{- $procName := printf "%sSearcher" $titleName -}}
{{- $domainName := $titleName -}}
{{- $domainPkg := pkgName .DomainPkg -}}

{{template "base" $procName}}

func (c *{{$procName}}) Do(ctx ctx.TaskContext) (out interface{}, errs []*errors.Error) {
  // initalize a domain object
  d := &{{$domainPkg}}.{{$domainName}}{}
  out, errs = d.Search(ctx)

	return
}
