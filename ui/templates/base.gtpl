{{define "base"}}
<!doctype html>
<html lang='en'>
    <head>
        <meta charset='utf-8'>
        <title>{{template "title" .}}</title>
        <link rel='stylesheet' href='/static/css/main.css'>
        <link rel='shortcut icon' href='/static/img/favicon.ico' type='image/x-icon'>
        <link rel='stylesheet' href='https://fonts.googleapis.com/css?family=Ubuntu+Mono:400,700'>
    </head>
    <body>
        <header>
            <h1><a href='/'>Krakhol</a></h1>
        </header>
        <nav>
            <a href='/'>Home</a>
            <a href='/login'>Login</a>
            {{if .Nav}}
                {{range .Nav}}
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