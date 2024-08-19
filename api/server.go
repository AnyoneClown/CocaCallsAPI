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

	apiRouter := r.PathPrefix("/api").Subrouter()

	// Router for routes that do not require authentication
	publicRouter := apiRouter.PathPrefix("/auth").Subrouter()
	publicRouter.HandleFunc("/login", s.handleLogin).Methods("POST")
	publicRouter.HandleFunc("/register", s.handleRegister).Methods("POST")

	// Router for routes that require authentication
	privateRouter := apiRouter.PathPrefix("/").Subrouter()
	privateRouter.Use(AuthenticationMiddleware)
	privateRouter.HandleFunc("/", s.handleMainPage).Methods("GET")

	return http.ListenAndServe(s.listenAddr, r)
}
