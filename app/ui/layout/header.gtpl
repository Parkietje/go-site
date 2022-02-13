{{define "header"}}
    <header>
        <div class="row">
            <div class="column">
                <h1><a href='/'>Pipelines</a></h1>
            </div>
            <div class="column">
            </div>
        </div>
    </header>
        <nav>
            {{if .PageContent.Navigation}}
                {{range .PageContent.Navigation}}
                    <a href='{{.Route}}'>{{.Title}}</a>
                {{end}}
            {{end}}
            {{if .User.Account}}
                <div class="topnav-right">
                    <a href='/login'>Logged in as: {{.User.Account}}</a>
                </div>    
            {{end}}
        </nav>
{{end}}