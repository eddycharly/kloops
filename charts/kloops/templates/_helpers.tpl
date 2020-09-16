{{- define "chatbot.name" -}}
{{- $name := default "chatbot" .Values.chatbot.nameOverride -}}
{{- printf "%s-%s" .Chart.Name $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "dashboard.name" -}}
{{- $name := default "dashboard" .Values.dashboard.nameOverride -}}
{{- printf "%s-%s" .Chart.Name $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}
