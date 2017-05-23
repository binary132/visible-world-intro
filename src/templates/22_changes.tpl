{{- $tableNameSingular := .Table.Name | singular | titleCase -}}
{{- $modelName := $tableNameSingular | titleCase -}}
{{- $varNameSingular := .Table.Name | singular | camelCase -}}
{{- $schemaTable := .Table.Name | .SchemaTable}}
// InsertG a single record. See Insert for whitelist behavior description.
func (o *{{$tableNameSingular}}) Changes()(ch *Changeset,err error) {
  ch = &Changeset{Table: "{{.Table.Name}}",
      Changes: []*ChangeItem{}, Operation: o.Operation()}
	v := reflect.Indirect(reflect.ValueOf(o.readonly))
	vnew := reflect.Indirect(reflect.ValueOf(o))

  for _, c := range o.Whitelist() {
    if f, ok := {{$modelName}}FieldMapping[c]; ok {
      var before, after interface{}
      if v.IsValid() {
        before = v.FieldByName(f).Interface()
      }

      if vnew.IsValid() {
        after = vnew.FieldByName(f).Interface()
      }

      chitem := &ChangeItem{Name: c}
      if o.operation == "DELETE" {
        chitem.Before = before
      }else{
        chitem.Before = before
        chitem.After = after
      }

      if !reflect.DeepEqual(chitem.Before, chitem.After) {
        ch.Changes = append(ch.Changes, chitem)
      }
    }
  }

  return
}

// Calculates changed columns on the object
func (o *{{$tableNameSingular}}) Whitelist()(wl []string) {
	if len(o.whitelist) > 0 {
		return o.whitelist
	}

	// Calculates changed columns as whitelist
	v := reflect.Indirect(reflect.ValueOf(o.readonly))
	vnew := reflect.Indirect(reflect.ValueOf(o))

	for _, c := range {{$varNameSingular}}Columns {
		if f, ok := {{$modelName}}FieldMapping[c]; ok {
			var before, after interface{}
			if v.IsValid() {
				before = v.FieldByName(f).Interface()
			}

			if vnew.IsValid() {
				after = vnew.FieldByName(f).Interface()
			}
			if !reflect.DeepEqual(before, after) || o.operation == "DELETE"{
				wl = append(wl, c)
			}
		}
	}

	return
}

func (o *{{$tableNameSingular}}) Operation() string {
  return o.operation
}

// Generated change history hook for models
func init() {
	chFunc := func(exec boil.Executor, s *{{$modelName}}) error {
    if s == nil || exec == nil {
      return nil
    }

		ch, _ := s.Changes()
		if changeable, ok := exec.(Changeable); ok {
			changeable.AddChange(ch)
		}

		return nil
	}

  afterSel := func(exec boil.Executor, s *{{$modelName}}) error {
    if s == nil || exec == nil {
      return nil
    }

    s.readonly = &{{$modelName}}{}
    *s.readonly = *s
    return nil
  }

  beforeInsert := func(exec boil.Executor, s *{{$modelName}}) error {
    if s == nil || exec == nil {
      return nil
    }

    s.operation = "INSERT"
    return nil
  }

  beforeUpdate := func(exec boil.Executor, s *{{$modelName}}) error {
    if s == nil || exec == nil {
      return nil
    }

    s.operation = "UPDATE"
    return nil
  }

  beforeUpsert := func(exec boil.Executor, s *{{$modelName}}) error {
    if s == nil || exec == nil {
      return nil
    }

    s.operation = "UPSERT"
    return nil
  }

  afterDelete := func(exec boil.Executor, s *{{$modelName}}) error {
    if s == nil || exec == nil {
      return nil
    }

    s.operation = "DELETE"
    return nil
  }

	Add{{$modelName}}Hook(boil.AfterSelectHook, afterSel)
	Add{{$modelName}}Hook(boil.BeforeInsertHook, beforeInsert)
	Add{{$modelName}}Hook(boil.AfterInsertHook, chFunc)
	Add{{$modelName}}Hook(boil.BeforeUpdateHook, beforeUpdate)
	Add{{$modelName}}Hook(boil.AfterUpdateHook, chFunc)
	Add{{$modelName}}Hook(boil.BeforeUpsertHook, beforeUpsert)
	Add{{$modelName}}Hook(boil.AfterUpsertHook, chFunc)
	Add{{$modelName}}Hook(boil.AfterDeleteHook, afterDelete)
	Add{{$modelName}}Hook(boil.AfterDeleteHook, chFunc)
}
