package api

import (
	"net/http"

	"github.com/AnyoneClown/CocaCallsAPI/storage"
	"github.com/gorilla/mux"
)

type Server struct {
	listenAddr string
	storage    storage.CockroachDB
}

func NewServer(listenAddr string, storage storage.CockroachDB) *Server {
	return &Server{
		listenAddr: listenAddr,
		storage:    storage,
	}
}

func (s *Server) Start() error {
	r := mux.NewRouter()

	// Router for routes that do not require authentication
	publicRouter := r.PathPrefix("/").Subrouter()
	publicRouter.HandleFunc("/login", s.handleLogin).Methods("GET")
	publicRouter.HandleFunc("/register", s.handleRegister).Methods("POST")

	// Router for routes that require authentication
	privateRouter := r.PathPrefix("/").Subrouter()
	privateRouter.Use(AuthenticationMiddleware)
	privateRouter.HandleFunc("/", s.handleMainPage).Methods("GET")

	return http.ListenAndServe(s.listenAddr, r)
}
