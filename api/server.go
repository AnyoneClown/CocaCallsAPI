package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	listenAddr string
}

func NewServer(listenAddr string) *Server {
	return &Server{
		listenAddr: listenAddr,
	}
}

func (s *Server) Start() error {
	r := mux.NewRouter()

	r.HandleFunc("/login", s.handleLogin).Methods("GET")

	authRouter := r.PathPrefix("/").Subrouter()
	authRouter.Use(AuthenticationMiddleware)
	authRouter.HandleFunc("/", s.handleMainPage).Methods("GET")

	return http.ListenAndServe(s.listenAddr, r)
}
