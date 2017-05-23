{{define "domain.show"}}
{{ $titleName := titleCase .Name }}
{{- $respSpecName := printf "%sShow" .Name -}}
{{- $domainName := $titleName -}}
{{- $baseName := printf "%sBase" .Name -}}
{{- $specPkg := pkgName .SpecPkg -}}
{{- $modelPkg := pkgName .ModelPkg -}}
{{- $domainPkg := pkgName .DomainPkg -}}
{{- $primaryModel := titleCase .PrimaryModel -}}

func (d *{{$domainName}}) Show(ctx ctx.TaskContext) (errs []*errors.Error) {
  ctx.WithExecutor(func(exec boil.Executor, networkID int64) (errCtx error) {
    params := ctx.Params()
    id, idErr := params.ID()
    userID := ctx.UserID()

    if idErr != nil {
      errs = append(errs, idErr)
    }

    if lErr := d.load(exec, id, userID, networkID); lErr != nil {
      errs = append(errs, lErr)
    }
    return
  })

  return
}
{{end}}
