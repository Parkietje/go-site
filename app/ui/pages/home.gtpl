{{template "base" .}}

{{define "title"}}Home{{end}}

{{define "main"}}
    <main>
        <div class="center-data">
            {{if .User.Account}}
                <img src="data:image/png;base64,{{.PageContent.PNG}}" class="img-fluid image-dashboard resize" />
            {{else}}
                <p>There's nothing to see here yet!</p>
            {{end}}
        </div>
    </main>
{{end}}

{{define "script"}}

{{end}}