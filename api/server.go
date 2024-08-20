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
    r.StrictSlash(true)

    apiRouter := r.PathPrefix("/api").Subrouter()

    // Router for routes that do not require authentication
    publicRouter := apiRouter.PathPrefix("/").Subrouter()

    authRouter := publicRouter.PathPrefix("/auth").Subrouter()
    authRouter.HandleFunc("/login/", s.handleLogin).Methods("POST")
    authRouter.HandleFunc("/register/", s.handleRegister).Methods("POST")

    // Router for routes that require authentication
    privateRouter := apiRouter.PathPrefix("/").Subrouter()
    privateRouter.Use(AuthenticationMiddleware)
    
    userRouter := privateRouter.PathPrefix("/user").Subrouter()
    userRouter.HandleFunc("/me/", s.handleUserMe).Methods("GET")

    return http.ListenAndServe(s.listenAddr, r)
}
