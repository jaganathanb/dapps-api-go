package controllers

import "github.com/jaganathanb/dapps-api-go/api/middlewares"

func (s *Server) initializeRoutes() {

	// Home Route
	s.Router.HandleFunc("/", middlewares.SetMiddlewareJSON(s.Home)).Methods("GET")

	// Login Route
	s.Router.HandleFunc("/login", middlewares.SetMiddlewareJSON(s.Login)).Methods("POST")

	//Users routes
	s.Router.HandleFunc("/users", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.CreateUser))).Methods("POST")
	s.Router.HandleFunc("/users", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.GetUsers))).Methods("GET")
	s.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.GetUser))).Methods("GET")
	s.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateUser))).Methods("PUT")
	s.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareAuthentication(s.DeleteUser)).Methods("DELETE")

	//GSTs routes
	s.Router.HandleFunc("/gsts", middlewares.SetMiddlewareJSON(s.CreateGst)).Methods("POST")
	s.Router.HandleFunc("/gsts", middlewares.SetMiddlewareJSON(s.GetGsts)).Methods("GET")
	s.Router.HandleFunc("/gsts/{id}", middlewares.SetMiddlewareJSON(s.GetGst)).Methods("GET")
	s.Router.HandleFunc("/gsts/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateGst))).Methods("PUT")
	s.Router.HandleFunc("/gsts/{id}", middlewares.SetMiddlewareAuthentication(s.DeleteGst)).Methods("DELETE")
}
