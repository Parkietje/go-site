{{template "base" .}}

{{define "title"}}Home{{end}}

{{define "main"}}
<h2>Home</h2>
{{if .Img}}
    <div class="col-4 col-sm-4 col-md-3 col-xl-2 center">
        <img src="data:image/png;base64,{{.Img}}" class="img-fluid image-dashboard" />
    </div>
{{else}}
    <p>There's nothing to see here yet!</p>
{{end}}
{{end}}