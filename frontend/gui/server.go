package gui

import (
	"context"
	"fmt"
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
	"github.com/pulsone21/threattrack/frontend/templates"
	"github.com/pulsone21/threattrack/lib/entities"
	"github.com/pulsone21/threattrack/lib/utils"
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

	server.Router.HandleFunc("/", utils.CreateHtmlHandleFunc(serveIndexPage)).Methods("GET")
	CreateIncidentHandler(server, backendAdress)
	CreateIndicatorHandler(server, backendAdress)
	server.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("frontend/static/"))))

	fmt.Println("FrontEnd Created")
	return server
}

func (s *Server) Run() {
	fmt.Printf("Serving Webserver at https://localhost:%s\n", s.Port)
	panic(http.ListenAndServe(fmt.Sprintf(":%s", s.Port), s))
}

func serveIndexPage(ctx context.Context, w http.ResponseWriter, r *http.Request) *entities.ApiError {
	fmt.Println("Serving Index Page")
	page := templates.Index()
	if err := page.Render(ctx, w); err != nil {
		return entities.InternalServerError(err, "/")
	}
	return nil
}
