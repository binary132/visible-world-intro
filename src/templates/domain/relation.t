{{define "domain.relation"}}
{{ $titleName := titleCase .Name }}
{{- $validatorFuncName := "validateSpec" -}}
{{- $reqSpecName := "Relation" -}}
{{- $reqRelationSpecName := $reqSpecName -}}
{{- $domainName := $titleName -}}
{{- $baseName := printf "%sBase" .Name -}}
{{- $specPkg := pkgName .SpecPkg -}}
{{- $modelPkg := pkgName .ModelPkg -}}
{{- $domainPkg := pkgName .DomainPkg -}}
{{- $primaryModel := titleCase .PrimaryModel -}}

func (d *{{$domainName}}) Relations(ctx ctx.TaskContext) (rel interface{}, errs []*errors.Error) {
  var s *{{$specPkg}}.{{$reqSpecName}}
  switch ctx.Method() {
  case "POST", "PUT":
    if vErr :=  d.{{$validatorFuncName}}(ctx, &validator.{{$reqSpecName}}{}); vErr != nil {
      errs = append(errs, vErr)
      return
    }

    s = &{{$specPkg}}.{{$reqSpecName}}{}
    ctx.Read(s)
  }

	ctx.WithExecutor(func(exec boil.Executor, networkID int64) (errCtx error) {
    switch ctx.Method() {
    case "POST":
      {{template "domain.usage.relations.func.1"}}
    case "PUT":
      {{template "domain.usage.relations.func.1"}}
    case "GET":
      {{template "domain.usage.relations.func.2"}}
    }

    return
  })

  return
}

func (d *{{$domainName}}) OneRelation(ctx ctx.TaskContext) (rel interface{}, errs []*errors.Error) { 
	ctx.WithExecutor(func(exec boil.Executor, networkID int64) (errCtx error) {
    switch ctx.Method() {
    case "POST":
      {{template "domain.usage.one_relation.func.1"}}
    case "PUT":
      {{template "domain.usage.one_relation.func.1"}}
    case "DELETE":
      {{template "domain.usage.one_relation.func.2"}}
    }

    return
  })

  return
}

func (d *{{$domainName}}) setRelation(exec boil.Executor, name string, rels *{{$specPkg}}.{{$reqRelationSpecName}}) (err *errors.Error) {
  if rels == nil {
    return
  }

  var ids []int64
  for _, item := range rels.Items {
    ids = append(ids, item.ID)
  }

  return d.setRelationByIDs(exec, name, ids...)
}

{{template "domain.usage.set_relation.func"}}
func (d *{{$domainName}}) setRelationByIDs(exec boil.Executor, name string, ids ...int64) (err *errors.Error) {
  return
}

func (d *{{$domainName}}) getRelation(exec boil.Executor, name string, params parser.Params) (rels interface{}, err *errors.Error) {
  return
}

func (d *{{$domainName}}) addRelation(exec boil.Executor, name string, rels *{{$specPkg}}.{{$reqRelationSpecName}}) (err *errors.Error) {
  if rels == nil {
    return
  }

  var ids []int64
  for _, item := range rels.Items {
    ids = append(ids, item.ID)
  }

  return d.addRelationByIDs(exec, name, ids...)
}

{{template "domain.usage.add_relation.func"}}
func (d *{{$domainName}}) addRelationByIDs(exec boil.Executor, name string, ids ...int64) (err *errors.Error) {
  return
}

func (d *{{$domainName}}) deleteRelation(exec boil.Executor, name string, rels *{{$specPkg}}.{{$reqRelationSpecName}}) (err *errors.Error) {
  if rels == nil {
    return
  }

  var ids []int64
  for _, item := range rels.Items {
    ids = append(ids, item.ID)
  }

  return d.deleteRelationByIDs(exec, name, ids...)
}

{{template "domain.usage.delete_relation.func"}}
func (d *{{$domainName}}) deleteRelationByIDs(exec boil.Executor, name string, ids ...int64) (err *errors.Error) {
  return
}
{{end}}
