package api

import (
	"net/http"

	"github.com/AnyoneClown/CocaCallsAPI/api/handlers"
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
    authHandler *handlers.AuthHandler
}

func NewServer(listenAddr string, storage storage.CockroachDB) *Server {
	return &Server{
		listenAddr: listenAddr,
		storage:    storage,
        authHandler: &handlers.AuthHandler{
            Storage: storage,
        },
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
    authRouter.HandleFunc("/register/", s.authHandler.HandleRegister).Methods("POST")
    authRouter.HandleFunc("/google/", s.authHandler.OauthGoogleLogin).Methods("GET")
    authRouter.HandleFunc("/{provider}/callback/", s.authHandler.OauthGoogleCallback).Methods("GET")

    // JWT router for login and verify token
    jwtRouter := apiRouter.PathPrefix("/jwt").Subrouter()
    jwtRouter.HandleFunc("/create/", s.authHandler.HandleJWTCreate).Methods("POST")
    jwtRouter.HandleFunc("/verify/", s.authHandler.HandleJWTVerify).Methods("POST")

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
