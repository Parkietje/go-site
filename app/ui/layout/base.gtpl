{{define "base"}}
<!doctype html>
<html lang='en'>
    <head>
        <meta charset='utf-8'>
        <title>{{template "title" .}}</title>
        <link rel='stylesheet' href='/static/css/main.css'>
        <link rel='shortcut icon' href='/static/img/favicon.ico' type='image/x-icon'>
        <link rel='stylesheet' href='https://fonts.googleapis.com/css?family=Ubuntu+Mono:400,700'>
        <script type="module">
            import { Application, Controller } from "https://unpkg.com/@hotwired/stimulus/dist/stimulus.js"
            window.Stimulus = Application.start()
            console.log("Hello!")
            {{template "script" .}}
        </script>
    </head>
    <body>
        <div class="page-container">
            <div class="content-wrap">
                {{template "header" .}}
                <div class="container-1">
                    {{template "main" .}}
                </div>
            </div>
            {{template "footer" .}}
        </div>
        <script src="/static/js/main.js" type="text/javascript"></script>
    </body>
</html>
{{end}}