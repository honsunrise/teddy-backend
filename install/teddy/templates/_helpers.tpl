{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}
{{- define "teddy.name" -}}
    {{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "teddy.fullname" -}}
    {{- if .Values.fullnameOverride -}}
        {{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" -}}
    {{- else -}}
        {{- $name := default .Chart.Name .Values.nameOverride -}}
        {{- if contains $name .Release.Name -}}
            {{- .Release.Name | trunc 63 | trimSuffix "-" -}}
        {{- else -}}
            {{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
        {{- end -}}
    {{- end -}}
{{- end -}}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "teddy.chart" -}}
    {{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "teddy.apis.base.name" -}}
  {{- printf "%s" "api-base" -}}
{{- end -}}

{{- define "teddy.apis.content.name" -}}
  {{- printf "%s" "api-content" -}}
{{- end -}}

{{- define "teddy.apis.message.name" -}}
  {{- printf "%s" "api-message" -}}
{{- end -}}

{{- define "teddy.apis.uaa.name" -}}
  {{- printf "%s" "api-uaa" -}}
{{- end -}}

{{- define "teddy.services.captcha.name" -}}
  {{- printf "%s" "api-base" -}}
{{- end -}}

{{- define "teddy.services.content.name" -}}
  {{- printf "%s" "api-content" -}}
{{- end -}}

{{- define "teddy.services.message.name" -}}
  {{- printf "%s" "api-message" -}}
{{- end -}}

{{- define "teddy.services.uaa.name" -}}
  {{- printf "%s" "api-uaa" -}}
{{- end -}}