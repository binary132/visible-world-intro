{{- $titleName := titleCase .Name -}}
{{- $procName := printf "%sReader" $titleName -}}
{{- $domainName := $titleName -}}
{{- $domainPkg := pkgName .DomainPkg -}}

{{template "base" $procName}}

func (c *{{$procName}}) Do(ctx ctx.TaskContext) (out interface{}, errs []*errors.Error) {
  // initalize a domain object
  d := &{{$domainPkg}}.{{$domainName}}{}
  if errs = d.Show(ctx); len(errs) > 0{
    return nil, errs
  }

	return d, nil
}
