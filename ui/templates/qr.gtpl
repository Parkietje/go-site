{{template "base" .}}

{{define "title"}}2FA{{end}}

{{define "main"}}
<h2>QR Code</h2>
<form action="/auth" method="post">
    Token:<input type="number" name="token">
    <input type="hidden" name="account" value={{.Account}}>
    <input type="submit" value="Auth">
</form>
<div class="col-4 col-sm-4 col-md-3 col-xl-2 center">
    <img src="data:image/png;base64,{{.PNG}}" class="img-fluid image-dashboard" />
</div>
{{end}}