package template

import (
	"html/template"
	"net/http"
)

var templates = template.Must(template.ParseFiles("index.html"))

func Execute(w http.ResponseWriter, templatename string, data interface{}) {
	templates.ExecuteTemplate(w, templatename, data)
}
