package controllers

import (
	"api-center/api/middlewares"
	"api-center/configs"
	"net/http"
)

func (s *Server) initializeRoutes(c *configs.Config) {
	// Apply CORS middleware globally
	s.Router.Use(middlewares.SetCORSMiddleware)

	s.Router.HandleFunc("/", s.Home).Methods("GET")
	s.Router.HandleFunc(c.ApiVersion+"/login", middlewares.SetMiddlewareJSON(s.Login)).Methods("POST")
	s.Router.HandleFunc(c.ApiVersion+"/users", middlewares.SetMiddlewareJSON(s.CreateUser)).Methods("POST")
	s.Router.HandleFunc(c.ApiVersion+"/users", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.GetUsers))).Methods("GET")
	s.Router.HandleFunc(c.ApiVersion+"/users/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.GetUser))).Methods("GET")
	s.Router.HandleFunc(c.ApiVersion+"/users/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateUser))).Methods("PUT")
	s.Router.HandleFunc(c.ApiVersion+"/users/{id}", middlewares.SetMiddlewareAuthentication(s.DeleteUser)).Methods("DELETE")

	// Catch-all OPTIONS route for CORS
	s.Router.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}).Methods("OPTIONS")
}
