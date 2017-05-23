{{- $titleName := titleCase .Name -}}
{{- $procName := printf "%sUpdater" $titleName -}}
{{- $domainName := $titleName -}}
{{- $domainPkg := pkgName .DomainPkg -}}

{{template "base" $procName}}

func (c *{{$procName}}) Process(ctx ctx.TaskContext) (out interface{}, newEntity bool, errs []*errors.Error) {
  // initalize a domain object
  d := &{{$domainPkg}}.{{$domainName}}{}
  if errs = d.Update(ctx); len(errs) > 0{
    return nil, false, errs
  }

	return d, false, nil
}
