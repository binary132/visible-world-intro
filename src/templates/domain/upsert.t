{{define "domain.upsert"}}
{{ $titleName := titleCase .Name }}
{{- $reqSpecName := printf "%sUpsert" .Name -}}
{{- $respSpecName := printf "%sShow" .Name -}}
{{- $domainName := $titleName -}}
{{- $baseName := printf "%sBase" .Name -}}
{{- $validatorFuncName := "validateSpec" -}}
{{- $specPkg := pkgName .SpecPkg -}}
{{- $modelPkg := pkgName .ModelPkg -}}
{{- $domainPkg := pkgName .DomainPkg -}}
{{- $primaryModel := titleCase .PrimaryModel -}}

func (d *{{$domainName}}) Upsert(ctx ctx.TaskContext) (newEntity bool, errs []*errors.Error) {
  if vErr :=  d.{{$validatorFuncName}}(ctx, &validator.{{$reqSpecName}}{}); vErr != nil {
    errs = append(errs, vErr)
    return
  }

  s := &{{$specPkg}}.{{$reqSpecName}}{}
	ctx.Read(s)


  userID := ctx.UserID()
	ctx.WithExecutor(func(exec boil.Executor, networkID int64) (errCtx error) {
    params := ctx.Params()
    id, idErr := params.ID()
    if idErr != nil {
      errs = append(errs, idErr)
      return
    }

    loadErr := d.load(exec, id, userID, networkID)
    if loadErr != nil {
      errs = append(errs, loadErr)
      return
    }

    if lErr := d.loadUpsertSpec(exec, userID, networkID, s); lErr != nil {
      errs = append(errs, lErr)
      return
    }

    return
  })

  return
}

{{template "domain.usage.load_spec.func"}}
func (d *{{$domainName}}) loadUpsertSpec(exec boil.Executor, userID, networkID int64, s *{{$specPkg}}.{{$reqSpecName}}) (err *errors.Error) {
  return
}
{{end}}
