{{define "domain.base"}}
{{ $titleName := titleCase .Name }}
{{- $validatorFuncName := "validateSpec" -}}
{{- $reqRelationSpecName := "Relation" -}}
{{- $linkItemSpecName := "LinkItem" -}}
{{- $reqSpecName := printf "%sCreate" $titleName -}}
{{- $respSpecName := printf "%sShow" $titleName -}}
{{- $respSliceSpecName := printf "%sList" .Name -}}
{{- $domainName := $titleName -}}
{{- $specPkg := pkgName .SpecPkg -}}
{{- $modelPkg := pkgName .ModelPkg -}}
{{- $domainPkg := pkgName .DomainPkg -}}
{{- $primaryModel := titleCase .PrimaryModel -}}
{{- $sliceDomainName := plural $titleName -}}

type {{$domainName}} {{$modelPkg}}.{{$primaryModel}}

{{template "domain.usage.marshaljson.func"}}
func (d *{{$domainName}}) MarshalJSON() (out []byte, err error) {
  o, errSpec := d.toSpec()
  if errSpec != nil {
    return nil, errSpec
  }

	return json.Marshal(o)
}

func (d *{{$domainName}}) toSpec() (s *{{$specPkg}}.{{$respSpecName}}, err *errors.Error) {
  s = &{{$specPkg}}.{{$respSpecName}} {
  // TODO: fullfill the spec
  }

  return
}

func (d *{{$domainName}}) model() *{{$modelPkg}}.{{$primaryModel}} {
	return (*{{$modelPkg}}.{{$primaryModel}})(d)
}

func (d *{{$domainName}}) load(exec boil.Executor, id, userID, networkID int64) (err *errors.Error) {
	o, e1 := {{$modelPkg}}.{{plural $primaryModel}}(exec,
    // TODO: add more filters here, i.e.
    //   qm.Where("group_type = 'SITE'")
		qm.Where("id = ? AND network_id = ?", id, networkID),
    // TODO: Load associatted data, i.e.
		// qm.Load("SiteSectionGroupAttributeData.SiteSectionAttributeLabel"),
    ).One()
	if e1 != nil {
		return errors.New(errors.DATA_ENTITY_NOT_FOUND, "id", "{{titleCase .Name}} not found")
	}

  *d = ({{$domainName}})(*o)

  return
}

func (d *{{$domainName}}) {{$validatorFuncName}}(ctx ctx.TaskContext, v SpecValidator) (errs *errors.Error) {
	readErr := ctx.Read(v)
	if readErr != nil {
		return errors.New(errors.DATA_JSON_PARSE_FAIL, "", readErr.Error())
	}

	formats := strfmt.NewFormats()
	if vErr := v.Validate(formats); vErr != nil {
        return NewSchemaValidationError(vErr)
	}

  return
}

type {{$sliceDomainName}} struct {
  {{plural $titleName}} []*{{$domainName}}
  URL        *parser.URL
	Page       int64
	PerPage    int64
	TotalCount int64
	TotalPage  int64
}

func (d *{{$sliceDomainName}}) MarshalJSON() (out []byte, err error) {
	o := &{{$specPkg}}.{{$respSliceSpecName}}{}

	for _, md := range d.{{plural $titleName}} {
		s, errSpec := md.toSpec()
    if errSpec != nil {
      return nil, errSpec
    }

		o.Items = append(o.Items, s)
	}

  d.paginate(o)
  d.buildLinks(o)

	if o.Items == nil {
		o.Items = []*{{$specPkg}}.{{$respSpecName}}{}
	}

	return json.Marshal(o)
}

func (d *{{$sliceDomainName}}) paginate(s *{{$specPkg}}.{{$respSliceSpecName}}) {
  if d.Page <= 0 {
    d.Page = 1
  }

  if d.PerPage <= 0 {
    d.PerPage = 10
  }

  s.Page = &d.Page
  s.PerPage = &d.PerPage
  s.TotalCount = &d.TotalCount

	if c := d.TotalCount / d.PerPage; (c * d.PerPage) < d.TotalCount  && d.TotalCount > 0 {
		d.TotalPage = d.TotalCount / d.PerPage + 1
	} else {
		d.TotalPage = d.TotalCount / d.PerPage
	}

  s.TotalPage = &d.TotalPage
}

func (d *{{$sliceDomainName}}) buildLinks(s *{{$specPkg}}.{{$respSliceSpecName}}) {
  if d.URL == nil {
    return
  }

  s.Links = append(s.Links, &{{$specPkg}}.{{$linkItemSpecName}}{
			Rel:  "self",
			Href: fmt.Sprintf(PAGINATION_FORMAT, d.URL.StringWithoutQuery(), d.Page, d.PerPage),
  })

  if d.Page > 1 {
    s.Links = append(s.Links, &{{$specPkg}}.{{$linkItemSpecName}}{
        Rel:  "prev",
        Href: fmt.Sprintf(PAGINATION_FORMAT, d.URL.StringWithoutQuery(), d.Page-1, d.PerPage),
    })
  }

	if d.Page < d.TotalPage {
		s.Links = append(s.Links, &{{$specPkg}}.{{$linkItemSpecName}}{
			Rel:  "next",
			Href: fmt.Sprintf(PAGINATION_FORMAT, d.URL.StringWithoutQuery(), d.Page+1, d.PerPage),
		})
		s.Links = append(s.Links, &{{$specPkg}}.{{$linkItemSpecName}}{
			Rel:  "last",
			Href: fmt.Sprintf(PAGINATION_FORMAT, d.URL.StringWithoutQuery(), d.TotalPage, d.PerPage),
		})
	}
}
{{end}}

