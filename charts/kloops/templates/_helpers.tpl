{{/* vim: set filetype=mustache: */}}

{{/*
Expand the name of the chatbot.
*/}}
{{- define "chatbot.name" -}}
{{- $name := default "chatbot" .Values.chatbot.nameOverride -}}
{{- printf "%s-%s" .Chart.Name $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Expand the name of the dashboard.
*/}}
{{- define "dashboard.name" -}}
{{- $name := default "dashboard" .Values.dashboard.nameOverride -}}
{{- printf "%s-%s" .Chart.Name $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create version.
*/}}
{{- define "version" -}}
{{ default .Chart.AppVersion .Values.version }}
{{- end -}}

{{/*
Create base labels
*/}}
{{- define "labels.base" -}}
app.kubernetes.io/instance: {{ $.Release.Name }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
app.kubernetes.io/part-of: kloops
{{- end -}}

{{/*
Create component labels
*/}}
{{- define "labels.component" -}}
app.kubernetes.io/component: {{ . }}
{{- end -}}

{{/*
Create version labels
*/}}
{{- define "labels.version" -}}
app.kubernetes.io/version: {{ template "version" . }}
{{- end -}}

{{/*
Create name labels
*/}}
{{- define "labels.name" -}}
app.kubernetes.io/name: {{ . }}
{{- end -}}