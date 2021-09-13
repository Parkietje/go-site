{{template "base" .}}

{{define "title"}}Login{{end}}

{{define "main"}}

{{if not .SessionCookie}}
    <h2>Log in</h2>
    <form action="/login" method="post">
    Username:<input type="text" name="username">
    Password:<input type="password" name="password">
    <input type="submit" value="Login">
    </form>
{{else}}
    <h2>Log out</h2>
    <form action="/logout">
    <input type="submit" value="Logout" />
    </form>
{{end}}

{{end}}