{{define "base"}}
<!doctype html>
<html lang='en'>
    <head>
        {{template "header"}}
    </head>
    <body>
        <header>
            <h1><a href='/'>Krakhol</a></h1>
        </header>
        <nav>
            <a href='/'>Home</a>
            {{if .PageContent.Navigation}}
                {{range .PageContent.Navigation}}
                    <a href='{{.Route}}'>{{.Title}}</a>
                {{end}}
            {{end}}
        </nav>
        <main>
            {{template "main" .}}
        </main>
        {{template "footer" .}}
        <script src="/static/js/main.js" type="text/javascript"></script>
    </body>
</html>
{{end}}