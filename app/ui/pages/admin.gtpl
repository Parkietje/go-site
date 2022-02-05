{{template "base" .}}

{{define "title"}}
  Admin
{{end}}

{{define "main"}}
  {{template "hello" .}}
  {{template "sidenav" .}}
{{end}}

{{define "script"}}
  {{template "hello-controller" .}}
{{end}}