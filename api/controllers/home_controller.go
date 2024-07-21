package controllers

import (
	"html/template"
	"net/http"
)

func (s *Server) Home(w http.ResponseWriter, r *http.Request) {
	var data = map[string]interface{}{
		"title": "Welcome to API-Center",
	}

	var tmpl = template.Must(template.ParseFiles("assets/views/index.html", "assets/views/_header.html"))

	var err = tmpl.ExecuteTemplate(w, "index", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
