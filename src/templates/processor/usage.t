{{define "proc.usage.creator" -}}
{{- $titleName := titleCase .Name -}}
{{- $procName := printf "%sCreator" $titleName -}}
// {{ $procName }} creates one {{ $titleName }}
{{- end}}

{{define "proc.usage.updater" -}}
{{- $titleName := titleCase .Name -}}
{{- $procName := printf "%sUpdater" $titleName -}}
// {{ $procName }} updates one {{ $titleName }}
{{- end}}

{{define "proc.usage.upserter" -}}
{{- $titleName := titleCase .Name -}}
{{- $procName := printf "%sUpserter" $titleName -}}
// {{ $procName }} upserts one {{ $titleName }}
{{- end}}

{{define "proc.usage.reader" -}}
{{- $titleName := titleCase .Name -}}
{{- $procName := printf "%sReader" $titleName -}}
// {{ $procName }} reads one {{ $titleName }}
{{- end}}

{{define "proc.usage.search" -}}
{{- $titleName := titleCase .Name -}}
{{- $procName := printf "%sSearcher" $titleName -}}
// {{ $procName }} searches one {{ $titleName }}
{{- end}}

{{define "proc.usage.one_rel_processor" -}}
{{- $titleName := titleCase .Name -}}
{{- $procName := printf "%sOneRelProcessor" $titleName -}}
// {{ $procName }} handles Create/Update/Delete on one relation of {{ $titleName }}
{{- end}}

{{define "proc.usage.rels_processor" -}}
{{- $titleName := titleCase .Name -}}
{{- $procName := printf "%sRelsProcessor" $titleName -}}
// {{ $procName }} handles Create/Update/Delete relations of {{ $titleName }}
{{- end}}

{{define "proc.usage.rels_lister" -}}
{{- $titleName := titleCase .Name -}}
{{- $procName := printf "%sRelsLister" $titleName -}}
// {{ $procName }} lists relation items of {{ $titleName }}
{{- end}}
