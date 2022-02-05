{{define "sidebar"}}
    {{if .PageContent.Sidebar}}
        <div class="sidebar">
            {{if .PageContent.Sidebar}}
                {{range .PageContent.Sidebar}}
                    <a href="{{.Route}}">{{.Title}}</a>
                {{end}}
            {{end}}
        </div>
    {{end}}
{{end}}