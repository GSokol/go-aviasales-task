{{/* vim: set filetype=mustache: */}}

{{- define "go-aviasales-task.fullname" -}}
{{- printf "%s-%s" .Release.Name .Chart.Name -}}
{{- end -}}
