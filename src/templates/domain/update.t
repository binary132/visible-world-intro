{{define "domain.update"}}
{{ $titleName := titleCase .Name }}
{{- $reqSpecName := printf "%sUpdate" .Name -}}
{{- $respSpecName := printf "%sShow" .Name -}}
{{- $validatorFuncName := "validateSpec" -}}
{{- $domainName := $titleName -}}
{{- $baseName := printf "%sBase" .Name -}}
{{- $specPkg := pkgName .SpecPkg -}}
{{- $modelPkg := pkgName .ModelPkg -}}
{{- $domainPkg := pkgName .DomainPkg -}}
{{- $primaryModel := titleCase .PrimaryModel -}}

{{template "domain.usage.update.func" .}}
func (d *{{$domainName}}) Update(ctx ctx.TaskContext) (errs []*errors.Error) {
  {{template "domain.usage.update.func.1" .}}
  if vErr :=  d.{{$validatorFuncName}}(ctx, &validator.{{$reqSpecName}}{}); vErr != nil {
    errs = append(errs, vErr)
    return
  }

  s := &{{$specPkg}}.{{$reqSpecName}}{}
	ctx.Read(s)


  userID := ctx.UserID()
	ctx.WithExecutor(func(exec boil.Executor, networkID int64) (errCtx error) {
    {{template "domain.usage.update.func.2" .}}
    id, idErr := ctx.Params().ID()
    if idErr != nil {
      errs = append(errs, idErr)
      return
    }

    {{template "domain.usage.update.func.3" .}}
    loadErr := d.load(exec, id, userID, networkID)
    if loadErr != nil {
      errs = append(errs, loadErr)
      return
    }

    {{template "domain.usage.update.func.4" .}}
    if lErr := d.loadUpdateSpec(exec, userID, networkID, s); lErr != nil {
      errs = append(errs, lErr)
      return
    }

    {{template "domain.usage.update.func.5" .}}
    d.model().Update(exec)

    return
  })

  {{template "domain.usage.update.func.6" .}}
  ctx.WithRequest(func(exec boil.Executor, networkID int64) (errCtx error) {
    if lErr := d.load(exec, d.ID, userID, networkID); lErr != nil {
      errs = append(errs, lErr)
    }
    return
  })

  return
}

{{template "domain.usage.load_spec.func"}}
func (d *{{$domainName}}) loadUpdateSpec(exec boil.Executor, userID, networkID int64, s *{{$specPkg}}.{{$reqSpecName}}) (err *errors.Error) {
  return
}
{{end}}
