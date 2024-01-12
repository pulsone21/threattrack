package gui

import (
	"fmt"
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
)

type Server struct {
	*mux.Router
	Port string
}

type PageHandler interface {
	createHandles(*Server)
	loadTemplate(string) (*template.Template, error)
}

func CreateServer(port, backendAdress string) *Server {
	server := &Server{
		Router: mux.NewRouter(),
		Port:   port,
	}

	CreateIncidentHandler(server, backendAdress)
	CreateIndicatorHandler(server, backendAdress)
	server.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("frontend/static/"))))
	fmt.Println("FrontEnd Created")
	return server
}

func (s *Server) Run() {
	fmt.Printf("Serving Webserver at https://localhost:%s", s.Port)
	panic(http.ListenAndServe(fmt.Sprintf(":%s", s.Port), s))
}
