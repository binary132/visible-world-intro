{{define "domain.create"}}
{{- $validatorFuncName := "validateSpec" -}}
{{- $titleName := titleCase .Name -}}
{{- $reqSpecName := printf "%sCreate" .Name -}}
{{- $respSpecName := printf "%sShow" .Name -}}
{{- $domainName := $titleName -}}
{{- $specPkg := pkgName .SpecPkg -}}
{{- $modelPkg := pkgName .ModelPkg -}}
{{- $domainPkg := pkgName .DomainPkg -}}
{{- $primaryModel := titleCase .PrimaryModel -}}

{{template "domain.usage.create.func" .}}
func (d *{{$domainName}}) Create(ctx ctx.TaskContext) (errs []*errors.Error) {
  {{template "domain.usage.create.func.1" .}}
  if vErr :=  d.{{$validatorFuncName}}(ctx, &validator.{{$reqSpecName}}{}); vErr != nil {
    errs = append(errs, vErr)
    return
  }

  {{template "domain.usage.create.func.2" .}}
  o := &{{$specPkg}}.{{$reqSpecName}}{}
	ctx.Read(o)

  userID := ctx.UserID()
  {{template "domain.usage.create.func.3" .}}
	ctx.WithExecutor(func(exec boil.Executor, networkID int64) (errCtx error) {
    {{template "domain.usage.create.func.4" .}}
    if lErr := d.loadCreateSpec(exec, userID, networkID, o); lErr != nil {
      errs = append(errs, lErr)
      return
    }

    {{template "domain.usage.create.func.5" .}}
    if insErr := d.model().Insert(exec); insErr != nil {
      errs = append(errs, errors.New(errors.INTERNAL_PROCESSOR_ERROR, "site", insErr.Error()))
    }
    return
  })

  {{template "domain.usage.create.func.6" .}}
  ctx.WithRequest(func(exec boil.Executor, networkID int64) (errCtx error) {
    if lErr := d.load(exec, d.ID, userID, networkID); lErr != nil {
      errs = append(errs, lErr)
    }
    return
  })

  return
}

{{template "domain.usage.load_spec.func"}}
func (d *{{$domainName}}) loadCreateSpec(exec boil.Executor, userID, networkID int64,s *{{$specPkg}}.{{$reqSpecName}}) (err *errors.Error) {
  return
}

{{end}}
