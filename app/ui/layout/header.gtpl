{{define "header"}}
    <header>
        <h1><a href='/'>Krakhol</a></h1>
        {{if .User.Account}}
            <h2><a href='/login'>Logged in as {{.User.Account}}</a></h2>
        {{end}}
    </header>
        <nav>
            <a href='/'>Home</a>
            {{if .PageContent.Navigation}}
                {{range .PageContent.Navigation}}
                    <a href='{{.Route}}'>{{.Title}}</a>
                {{end}}
            {{end}}
        </nav>
{{end}}