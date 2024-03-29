{{- $tableNameSingular := .Table.Name | singular | titleCase -}}
{{- $varNameSingular := .Table.Name | singular | camelCase -}}
{{- $schemaTable := .Table.Name | .SchemaTable}}
// UpsertG attempts an insert, and does an update or ignore on conflict.
func (o *{{$tableNameSingular}}) UpsertG({{if ne .DriverName "mysql"}}updateOnConflict bool, conflictColumns []string, {{end}}updateColumns []string,	whitelist ...string) error {
	return o.Upsert(boil.GetDB(), {{if ne .DriverName "mysql"}}updateOnConflict, conflictColumns, {{end}}updateColumns, whitelist...)
}

// UpsertGP attempts an insert, and does an update or ignore on conflict. Panics on error.
func (o *{{$tableNameSingular}}) UpsertGP({{if ne .DriverName "mysql"}}updateOnConflict bool, conflictColumns []string, {{end}}updateColumns []string,	whitelist ...string) {
	if err := o.Upsert(boil.GetDB(), {{if ne .DriverName "mysql"}}updateOnConflict, conflictColumns, {{end}}updateColumns, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpsertP attempts an insert using an executor, and does an update or ignore on conflict.
// UpsertP panics on error.
func (o *{{$tableNameSingular}}) UpsertP(exec boil.Executor, {{if ne .DriverName "mysql"}}updateOnConflict bool, conflictColumns []string, {{end}}updateColumns []string,	whitelist ...string) {
	if err := o.Upsert(exec, {{if ne .DriverName "mysql"}}updateOnConflict, conflictColumns, {{end}}updateColumns, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
func (o *{{$tableNameSingular}}) Upsert(exec boil.Executor, {{if ne .DriverName "mysql"}}updateOnConflict bool, conflictColumns []string, {{end}}updateColumns []string, whitelist ...string) error {
	if o == nil {
		return errors.New("{{.PkgName}}: no {{.Table.Name}} provided for upsert")
	}
  o.whitelist = whitelist
  o.operation = "UPSERT"

	{{- template "timestamp_upsert_helper" . }}

	{{if not .NoHooks -}}
	if err := o.doBeforeUpsertHooks(exec); err != nil {
		return err
	}
	{{- end}}

	nzDefaults := queries.NonZeroDefaultSet({{$varNameSingular}}ColumnsWithDefault, o)

	// Build cache key in-line uglily - mysql vs postgres problems
	buf := strmangle.GetBuffer()
	{{if ne .DriverName "mysql" -}}
	if updateOnConflict {
		buf.WriteByte('t')
	} else {
		buf.WriteByte('f')
	}
	buf.WriteByte('.')
	for _, c := range conflictColumns {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	{{end -}}
	for _, c := range updateColumns {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range whitelist {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range nzDefaults {
		buf.WriteString(c)
	}
	key := buf.String()
	strmangle.PutBuffer(buf)

	{{$varNameSingular}}UpsertCacheMut.RLock()
	cache, cached := {{$varNameSingular}}UpsertCache[key]
	{{$varNameSingular}}UpsertCacheMut.RUnlock()

	var err error

	if !cached {
		var ret []string
		whitelist, ret = strmangle.InsertColumnSet(
			{{$varNameSingular}}Columns,
			{{$varNameSingular}}ColumnsWithDefault,
			{{$varNameSingular}}ColumnsWithoutDefault,
			nzDefaults,
			whitelist,
		)
		update := strmangle.UpdateColumnSet(
			{{$varNameSingular}}Columns,
			{{$varNameSingular}}PrimaryKeyColumns,
			updateColumns,
		)
		if len(update) == 0 {
			return errors.New("{{.PkgName}}: unable to upsert {{.Table.Name}}, could not build update column list")
		}

		{{if ne .DriverName "mysql" -}}
		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len({{$varNameSingular}}PrimaryKeyColumns))
			copy(conflict, {{$varNameSingular}}PrimaryKeyColumns)
		}
		cache.query = queries.BuildUpsertQueryPostgres(dialect, "{{$schemaTable}}", updateOnConflict, ret, update, conflict, whitelist)
		{{- else -}}
		cache.query = queries.BuildUpsertQueryMySQL(dialect, "{{.Table.Name}}", update, whitelist)
		cache.retQuery = fmt.Sprintf(
			"SELECT %s FROM {{.LQ}}{{.Table.Name}}{{.RQ}} WHERE {{whereClause .LQ .RQ 0 .Table.PKey.Columns}}",
			strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, ret), ","),
		)
		{{- end}}

		cache.valueMapping, err = queries.BindMapping({{$varNameSingular}}Type, {{$varNameSingular}}Mapping, whitelist)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping({{$varNameSingular}}Type, {{$varNameSingular}}Mapping, ret)
			if err != nil {
				return err
			}
		}
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)
	var returns []interface{}
	if len(cache.retMapping) != 0 {
		returns = queries.PtrsFromMapping(value, cache.retMapping)
	}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.query)
		fmt.Fprintln(boil.DebugWriter, vals)
	}

	{{if .UseLastInsertID -}}
	{{- $canLastInsertID := .Table.CanLastInsertID -}}
	{{if $canLastInsertID -}}
	result, err := exec.Exec(cache.query, vals...)
	{{else -}}
	_, err = exec.Exec(cache.query, vals...)
	{{- end}}
	if err != nil {
		return errors.Wrap(err, "{{.PkgName}}: unable to upsert for {{.Table.Name}}")
	}

	{{if $canLastInsertID -}}
	var lastID int64
	{{- end}}
	var identifierCols []interface{}

	if len(cache.retMapping) == 0 {
		goto CacheNoHooks
	}

	{{if $canLastInsertID -}}
	lastID, err = result.LastInsertId()
	if err != nil {
		return ErrSyncFail
	}

	{{$colName := index .Table.PKey.Columns 0 -}}
	{{- $col := .Table.GetColumn $colName -}}
	{{- $colTitled := $colName | titleCase}}
	o.{{$colTitled}} = {{$col.Type}}(lastID)
	if lastID != 0 && len(cache.retMapping) == 1 && cache.retMapping[0] == {{$varNameSingular}}Mapping["{{$colTitled}}"] {
		goto CacheNoHooks
	}
	{{- end}}

	identifierCols = []interface{}{
		{{range .Table.PKey.Columns -}}
		o.{{. | titleCase}},
		{{end -}}
	}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.retQuery)
		fmt.Fprintln(boil.DebugWriter, identifierCols...)
	}

	err = exec.QueryRow(cache.retQuery, identifierCols...).Scan(returns...)
	if err != nil {
		return errors.Wrap(err, "{{.PkgName}}: unable to populate default values for {{.Table.Name}}")
	}
	{{- else}}
	if len(cache.retMapping) != 0 {
		err = exec.QueryRow(cache.query, vals...).Scan(returns...)
	} else {
		_, err = exec.Exec(cache.query, vals...)
	}
	if err != nil {
		return errors.Wrap(err, "{{.PkgName}}: unable to upsert for {{.Table.Name}}")
	}
	{{- end}}

{{if .UseLastInsertID -}}
CacheNoHooks:
{{end -}}
	if !cached {
		{{$varNameSingular}}UpsertCacheMut.Lock()
		{{$varNameSingular}}UpsertCache[key] = cache
		{{$varNameSingular}}UpsertCacheMut.Unlock()
	}

	{{if not .NoHooks -}}
	return o.doAfterUpsertHooks(exec)
	{{- else -}}
	return nil
	{{- end}}
}
