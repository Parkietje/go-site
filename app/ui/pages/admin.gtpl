{{template "base" .}}

{{define "title"}}
  Admin
{{end}}

{{define "main"}}
  {{template "admin" .}}
  <main>
    {{template "admin-data" .}}
  </main>
{{end}}

{{define "script"}}
  {{template "admin-controller" .}}
{{end}}