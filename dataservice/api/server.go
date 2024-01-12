package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pulsone21/threattrack/dataService/handlers"
	"github.com/pulsone21/threattrack/dataService/storage"
)

type Server struct {
	*mux.Router
	listenAddr string
	store      storage.Storage
}

func NewServer(listenAddress string, storage storage.Storage) *Server {
	return &Server{
		Router:     mux.NewRouter(),
		listenAddr: listenAddress,
		store:      storage,
	}
}

func (s *Server) Run() error {
	handlers.CreateHandlers(s.Router, s.store)
	return http.ListenAndServe(fmt.Sprintf(":%s", s.listenAddr), s)
}
