{{template "base" .}}

{{define "title"}}Deploy{{end}}

{{define "main"}}
    <main>
        <div class="center-data shiftleft">
            {{template "file-upload"}}
        </div>
    </main>
{{end}}

{{define "script"}}
    {{template "file-upload-controller"}}
{{end}}