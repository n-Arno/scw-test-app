package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
	"text/template"
)

func multiline(s string) string {
	return strings.Replace(s, "\n", "<br/>", -1)
}

var headerHtml = `<html><head><title>scw-test-app</title>
<style>` + simpleCss + `
</style></head><body><header><h1>Weyland-Yutani</h1></header>
`

var footerHtml = `<footer><p>Source code is <a href="https://github.com/n-Arno/scw-test-app">here</a></p>
<p>CSS leveraging <a href="https://simplecss.org/">Simple CSS</a></p></footer></body></html>`

var newsHtml = headerHtml + `<h1>News</h1>
{{ range . }}
<article>
<h2>{{.Title}}</h2>
<p>{{.Content | multiline}}</p>
</article>
{{ end }}
` + footerHtml

var adminHtml = headerHtml + `<h1>Admin Panel</h1><ul>
<li><a href="/admin/news">Add news</a></li>
<li><a href="/admin/config">Configure App</a></li>
<li><a href="/metrics">Metrics</a></li>
<li><a href="/">Back</a></li>
</ul>` + footerHtml

var adminNewsHtml = headerHtml + `<h1>Add news</h1>
  <form action="/admin/news" method="POST">
  <input type="text" placeholder="News Title" name="title" id="title" size="40"/></br/>
  <textarea rows="10" cols="40"  name="content" id="content"></textarea><br/>
  <button type="submit">Add News</button>
  </form>
<br/><a href="/admin/">Back</a>` + footerHtml

var addNewsHtml = headerHtml + `<h2>Created!</h2>
<br/><a href="/admin/news">Back</a>` + footerHtml

var adminConfigHtml = headerHtml + `<h1>Edit config</h1>
  <form action="/admin/config" method="POST">
  <p>DB config</p>
  <input type="text" placeholder="" name="db_port" id="db_port" value="{{.Db.Port}}"/><br/>
  <input type="text" placeholder="" name="db_host" id="db_host" value="{{.Db.Host}}"/><br/>
  <input type="text" placeholder="" name="db_name" id="db_name" value="{{.Db.Name}}"/><br/>
  <input type="text" placeholder="" name="db_user" id="db_user" value="{{.Db.User}}"/><br/>
  <input type="password" placeholder="" name="db_pass" id="db_pass" value="{{.Db.Pass}}"/><br/>
  <button type="submit">Update config</button>
  </form>
<br/><a href="/admin/">Back</a>` + footerHtml

var updateConfigHtml = headerHtml + `<h2>Updated!</h2><ul>
<li><a href="/admin/">Back to admin</a></li>
<li><a href="/">Back to main page</a></li>
</ul>` + footerHtml

var errorDbHtml = headerHtml + `<h2>Error!</h2>
<p>Access to database may not be correctly configured. Please update in admin panel:</p>
<br/><a href="/admin/config">Go to Admin</a>` + footerHtml

func html(tmpl string, data any) ([]byte, error) {
	var funcMap = template.FuncMap{
		"multiline": multiline,
	}
	var b bytes.Buffer
	t, _ := template.New("").Funcs(funcMap).Parse(tmpl)
	err := t.Execute(&b, data)
	return b.Bytes(), err
}

func NewsHTML(data any) ([]byte, error) {
	return html(newsHtml, data)
}

func AdminHTML() ([]byte, error) {
	return html(adminHtml, nil)
}

func AdminNewsHTML() ([]byte, error) {
	return html(adminNewsHtml, nil)
}

func AddNewsHTML() ([]byte, error) {
	return html(addNewsHtml, nil)
}

func AdminConfigHTML(data any) ([]byte, error) {
	return html(adminConfigHtml, data)
}

func UpdateConfigHTML() ([]byte, error) {
	return html(updateConfigHtml, nil)
}

func ErrorDbHTML() ([]byte, error) {
	return html(errorDbHtml, nil)
}

func simpleAnswer(w http.ResponseWriter, httpStatus int, anyStruct any) {

	result, err := json.MarshalIndent(anyStruct, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(httpStatus)
	w.Write(result)
}
