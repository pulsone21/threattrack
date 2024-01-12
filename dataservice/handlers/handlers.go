package handlers

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/pulsone21/threattrack/dataservice/storage"
	utils "github.com/pulsone21/threattrack/lib/utils"
)

func CreateHandlers(r *mux.Router, s storage.Storage) {
	inc_sR := r.PathPrefix("/incidents").Subrouter()
	inc_sR.HandleFunc("", utils.CreateApiHandleFunc(s.HandleGetAll)).Methods(http.MethodGet)
	inc_sR.HandleFunc("", utils.CreateApiHandleFunc(s.HandleCreate)).Methods(http.MethodPost)
	inc_sR.HandleFunc("/{id}", utils.CreateApiHandleFunc(s.HandleGetID)).Methods(http.MethodGet)
	inc_sR.HandleFunc("/{id}", utils.CreateApiHandleFunc(s.HandleUpdate)).Methods(http.MethodPut)
	inc_sR.HandleFunc("/{id}", utils.CreateApiHandleFunc(s.HandleDelete)).Methods(http.MethodDelete)
	inc_sR.HandleFunc("/query", utils.CreateApiHandleFunc(s.HandleGetQuery)).Methods(http.MethodGet)

	it_sR := r.PathPrefix("/incidenttypes").Subrouter()
	it_sR.HandleFunc("", utils.CreateApiHandleFunc(s.HandleGetAll)).Methods(http.MethodGet)
	it_sR.HandleFunc("", utils.CreateApiHandleFunc(s.HandleCreate)).Methods(http.MethodPost)
	it_sR.HandleFunc("/{id}", utils.CreateApiHandleFunc(s.HandleGetID)).Methods(http.MethodGet)
	it_sR.HandleFunc("/{id}", utils.CreateApiHandleFunc(s.HandleUpdate)).Methods(http.MethodPut)
	it_sR.HandleFunc("/{id}", utils.CreateApiHandleFunc(s.HandleDelete)).Methods(http.MethodDelete)

	usr_sR := r.PathPrefix("/users").Subrouter()
	usr_sR.HandleFunc("", utils.CreateApiHandleFunc(s.HandleGetAll)).Methods(http.MethodGet)
	usr_sR.HandleFunc("", utils.CreateApiHandleFunc(s.HandleCreate)).Methods(http.MethodPost)
	usr_sR.HandleFunc("/{id}", utils.CreateApiHandleFunc(s.HandleGetID)).Methods(http.MethodGet)
	usr_sR.HandleFunc("/{id}", utils.CreateApiHandleFunc(s.HandleUpdate)).Methods(http.MethodPut)
	usr_sR.HandleFunc("/{id}", utils.CreateApiHandleFunc(s.HandleDelete)).Methods(http.MethodDelete)

	task_sR := r.PathPrefix("/tasks").Subrouter()
	task_sR.HandleFunc("", utils.CreateApiHandleFunc(s.HandleGetAll)).Methods(http.MethodGet)
	task_sR.HandleFunc("", utils.CreateApiHandleFunc(s.HandleCreate)).Methods(http.MethodPost)
	task_sR.HandleFunc("/{id}", utils.CreateApiHandleFunc(s.HandleGetID)).Methods(http.MethodGet)
	task_sR.HandleFunc("/{id}", utils.CreateApiHandleFunc(s.HandleUpdate)).Methods(http.MethodPut)
	task_sR.HandleFunc("/{id}", utils.CreateApiHandleFunc(s.HandleDelete)).Methods(http.MethodDelete)
	task_sR.HandleFunc("/query", utils.CreateApiHandleFunc(s.HandleGetQuery)).Methods(http.MethodGet)
}
