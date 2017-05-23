{{- $titleName := titleCase .Name -}}
{{- $procName := printf "%sUpserter" $titleName -}}
{{- $domainName := $titleName -}}
{{- $domainPkg := pkgName .DomainPkg -}}

{{template "base" $procName}}

func (c *{{$procName}}) Process(ctx ctx.TaskContext) (out interface{}, newEntity bool, errs []*errors.Error) {
  // initalize a domain object
  d := &{{$domainPkg}}.{{$domainName}}{}
  if newEntity, errs = d.Upsert(ctx); len(errs) > 0{
    return nil, false, errs
  } else {
    return d, newEntity, nil
  }
}
