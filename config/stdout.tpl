{{ if eq .Req "runner" -}}
ID: {{ .ID }}, Description: {{ .Description }}, Active: {{ .Active }}, IP: {{ .IPAddress }}, Shared: {{ .IsShared}}, Status: {{ .Status }}
{{ else if eq .Req "groupRunner" -}}
Runner ID: {{ .ID}}
	Type: {{ .Name }}
	Runner: {{ .Description }}
	Active: {{ .Active }}, IP: {{ .IPAddress }}, Shared: {{ .IsShared }}, Status: {{ .Status }}
{{ else if eq .Req "user" -}}
ID: {{ .ID }}, Name: {{ .Name }}
{{ else if eq .Req "project" -}}
ID: {{ .ID }}, Name: {{ .Name }}, Last Activity: {{ .LastActivity }} h
{{ end }}