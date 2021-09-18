{{template "base" .}}

{{define "title"}}Home{{end}}

{{define "main"}}
{{if .User.Account}}

    <form action="http://localhost:8080" target="_blank" rel="noopener noreferrer">
        <input type="submit" name="submit" value="Go to Fileserver" />
    </form>


    <img src="data:image/png;base64,{{.PageContent.PNG}}" class="img-fluid image-dashboard" />

{{else}}
    <p>There's nothing to see here yet!</p>
{{end}}
{{end}}