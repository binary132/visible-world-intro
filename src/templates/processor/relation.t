{{- $titleName := titleCase .Name -}}
{{- $procName := printf "%sOneRelProcessor" $titleName -}}
{{- $domainName := $titleName -}}
{{- $domainPkg := pkgName .DomainPkg -}}

{{template "base" $procName}}

func (c *{{$procName}}) Process(ctx ctx.TaskContext) (out interface{}, newEntity bool, errs []*errors.Error) {
  // initalize a domain object
  d := &{{$domainPkg}}.{{$domainName}}{}
  if out, errs = d.OneRelation(ctx); len(errs) > 0{
    return nil, false, errs
  }

	return
}

{{- $titleName := titleCase .Name -}}
{{- $procName := printf "%sRelsProcessor" $titleName -}}
{{- $domainName := $titleName -}}
{{- $domainPkg := pkgName .DomainPkg -}}

{{template "base" $procName}}

func (c *{{$procName}}) Process(ctx ctx.TaskContext) (out interface{}, newEntity bool, errs []*errors.Error) {
  // initalize a domain object
  d := &{{$domainPkg}}.{{$domainName}}{}
  if out, errs = d.Relations(ctx); len(errs) > 0{
    return nil, false, errs
  }

	return
}

{{- $titleName := titleCase .Name -}}
{{- $procName := printf "%sRelsLister" $titleName -}}
{{- $domainName := $titleName -}}
{{- $domainPkg := pkgName .DomainPkg -}}

{{template "base" $procName}}

func (c *{{$procName}}) Do(ctx ctx.TaskContext) (out interface{}, errs []*errors.Error) {
  // initalize a domain object
  d := &{{$domainPkg}}.{{$domainName}}{}
  if out, errs = d.Relations(ctx); len(errs) > 0{
    return nil, errs
  }

	return
}
