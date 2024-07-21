package controllers

import "net/http"

func (s *Server) InitializeStaticAsset() {
	a := http.StripPrefix("/static/", http.FileServer(http.Dir("./assets/")))
	s.Router.PathPrefix("/static/").Handler(a)
}
