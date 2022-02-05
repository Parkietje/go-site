{{template "base" .}}

{{define "title"}}Home{{end}}

{{define "main"}}
    <main>
    {{if .User.Account}}
        <img src="data:image/png;base64,{{.PageContent.PNG}}" class="img-fluid image-dashboard" />
    {{else}}
        <p>There's nothing to see here yet!</p>
    {{end}}
    </main>
{{end}}

{{define "script"}}

{{end}}