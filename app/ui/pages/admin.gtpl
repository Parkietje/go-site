{{template "base" .}}

{{define "title"}}
  Admin
{{end}}

{{define "main"}}
  {{template "admin-menu" .}}
  <main>
    {{template "admin-menu-data" .}}
  </main>
{{end}}

{{define "script"}}
  {{template "hello-controller" .}}
  {{template "admin-menu-controller" .}}
{{end}}