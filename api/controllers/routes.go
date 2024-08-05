package controllers

import (
	"api-center/api/middlewares"
	"api-center/configs"
)

func (s *Server) initializeRoutes(c *configs.Config) {
	// Home Route
	s.Router.HandleFunc("/", s.Home).Methods("GET")

	// Login Route
	s.Router.HandleFunc(c.ApiVersion+"/login", middlewares.SetMiddlewareJSON(s.Login)).Methods("POST")

	// // Users Route
	s.Router.HandleFunc(c.ApiVersion+"/users", middlewares.SetMiddlewareJSON(s.CreateUser)).Methods("POST")
	s.Router.HandleFunc(c.ApiVersion+"/users", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.GetUsers))).Methods("GET")
	s.Router.HandleFunc(c.ApiVersion+"/users/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.GetUser))).Methods("GET")
	s.Router.HandleFunc(c.ApiVersion+"/users/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateUser))).Methods("PUT")
	s.Router.HandleFunc(c.ApiVersion+"/users/{id}", middlewares.SetMiddlewareAuthentication(s.DeleteUser)).Methods("DELETE")
}
