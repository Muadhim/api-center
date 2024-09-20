package controllers

import (
	"api-center/api/middlewares"
	"api-center/configs"
	"net/http"
)

func (s *Server) initializeRoutes(c *configs.Config) {
	// Apply CORS middleware globally
	s.Router.Use(middlewares.SetCORSMiddleware)

	// root
	s.Router.HandleFunc("/", s.Home).Methods("GET")

	// login
	s.Router.HandleFunc(c.ApiVersion+"/login", middlewares.SetMiddlewareJSON(s.Login)).Methods("POST")

	// users
	s.Router.HandleFunc(c.ApiVersion+"/users", middlewares.SetMiddlewareJSON(s.CreateUser)).Methods("POST")
	// s.Router.HandleFunc(c.ApiVersion+"/users", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.GetUsers))).Methods("GET")
	s.Router.HandleFunc(c.ApiVersion+"/users/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.GetUser))).Methods("GET")
	s.Router.HandleFunc(c.ApiVersion+"/users/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateUser))).Methods("PUT")
	s.Router.HandleFunc(c.ApiVersion+"/users/{id}", middlewares.SetMiddlewareAuthentication(s.DeleteUser)).Methods("DELETE")

	// projects
	s.Router.HandleFunc(c.ApiVersion+"/project", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.CreateProject))).Methods("POST")
	s.Router.HandleFunc(c.ApiVersion+"/project/{id}/members", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateProjectMembers))).Methods("PUT")
	s.Router.HandleFunc(c.ApiVersion+"/project/{id}/members", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.DeleteProjectMembers))).Methods("DELETE")
	s.Router.HandleFunc(c.ApiVersion+"/project/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.DeleteProject))).Methods("DELETE")
	s.Router.HandleFunc(c.ApiVersion+"/projects", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.GetProjects))).Methods("GET")
	s.Router.HandleFunc(c.ApiVersion+"/project/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.GetProjectByID))).Methods("GET")
	s.Router.HandleFunc(c.ApiVersion+"/project/{id}/invite", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.GenerateInviteToken))).Methods("GET")
	s.Router.HandleFunc(c.ApiVersion+"/project/join", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.ValidateInvitationToken))).Methods("POST")
	// project tree
	s.Router.HandleFunc(c.ApiVersion+"/project/{id}/tree", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.GetProjectTree))).Methods("GET")

	// groups
	s.Router.HandleFunc(c.ApiVersion+"/group", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.CreateGroup))).Methods("POST")
	s.Router.HandleFunc(c.ApiVersion+"/groups", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.GetGroups))).Methods("GET")
	s.Router.HandleFunc(c.ApiVersion+"/group/{id}/invite", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.GenInviteGroupToken))).Methods("GET")
	s.Router.HandleFunc(c.ApiVersion+"/group/join", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.ValidateInvitationGroupToken))).Methods("POST")

	// folders
	s.Router.HandleFunc(c.ApiVersion+"/project-folder", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.CreateProjectFolder))).Methods("POST")
	s.Router.HandleFunc(c.ApiVersion+"/project-folder/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.DeleteProjectFolder))).Methods("DELETE")
	s.Router.HandleFunc(c.ApiVersion+"/project-folder", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateProjectFolder))).Methods("PUT")

	// project api
	s.Router.HandleFunc(c.ApiVersion+"/project-api", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.CreateProjectApi))).Methods("POST")
	s.Router.HandleFunc(c.ApiVersion+"/project-api/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.DeleteProjectApi))).Methods("DELETE")
	s.Router.HandleFunc(c.ApiVersion+"/project-api", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateProjectApi))).Methods("PUT")

	// Catch-all OPTIONS route for CORS
	s.Router.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}).Methods("OPTIONS")
}
