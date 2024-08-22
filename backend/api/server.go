package api

import (
	"net/http"

	"github.com/AnyoneClown/CocaCallsAPI/storage"
	"github.com/AnyoneClown/CocaCallsAPI/utils"
	"github.com/gorilla/mux"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/google"
	"github.com/rs/cors"
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

    goth.UseProviders(
        google.New(
            utils.GetEnvVariable("CLIENT_ID"), 
            utils.GetEnvVariable("CLIENT_SECRET"),
            utils.GetEnvVariable("CLIENT_CALLBACK_URL"),
        ),
    )

    apiRouter := r.PathPrefix("/api").Subrouter()

    // Auth router for all operations like OAuth 2.0, register
    authRouter := apiRouter.PathPrefix("/auth").Subrouter()
    authRouter.HandleFunc("/register/", s.handleRegister).Methods("POST")
    authRouter.HandleFunc("/google/", s.oauthGoogleLogin).Methods("GET")
    authRouter.HandleFunc("/{provider}/callback/", s.oauthGoogleCallback).Methods("GET")

    // JWT router for login and verify token
    jwtRouter := apiRouter.PathPrefix("/jwt").Subrouter()
    jwtRouter.HandleFunc("/create/", s.handleJWTCreate).Methods("POST")
    jwtRouter.HandleFunc("/verify/", s.handleJWTVerify).Methods("POST")

    // Configure CORS
    corsOptions := cors.New(cors.Options{
        AllowedOrigins:   []string{"http://localhost:3000"},
        AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowedHeaders:   []string{"Content-Type", "Authorization"},
        ExposedHeaders:   []string{"Link"},
        AllowCredentials: true,
    })

    return http.ListenAndServe(s.listenAddr, corsOptions.Handler(r))
}
