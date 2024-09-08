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
	listenAddr     string
	storage        storage.CockroachDB
	defaultHandler *handlers.DefaultHandler
}

func NewServer(listenAddr string, storage storage.CockroachDB) *Server {
	return &Server{
		listenAddr: listenAddr,
		storage:    storage,
		defaultHandler: &handlers.DefaultHandler{
			Storage: &storage,
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
	authRouter.HandleFunc("/register/", s.defaultHandler.HandleRegister).Methods("POST")
	authRouter.HandleFunc("/google/", s.defaultHandler.OauthGoogleLogin).Methods("GET")
	authRouter.HandleFunc("/{provider}/callback/", s.defaultHandler.OauthGoogleCallback).Methods("GET")

	// JWT router for login and verify token
	jwtRouter := apiRouter.PathPrefix("/jwt").Subrouter()
	jwtRouter.HandleFunc("/create/", s.defaultHandler.HandleJWTCreate).Methods("POST")
	jwtRouter.HandleFunc("/verify/", s.defaultHandler.HandleJWTVerify).Methods("POST")

	// User routes
	userRouter := apiRouter.PathPrefix("/users").Subrouter()
	userRouter.Use(AuthenticationMiddleware)
	userRouter.HandleFunc("/", s.defaultHandler.GetUsers).Methods("GET")
	userRouter.HandleFunc("/{userID}/", s.defaultHandler.GetUser).Methods("GET")
	userRouter.HandleFunc("/picture/{userID}/", s.defaultHandler.UpdateUserProfilePicture).Methods("PUT")

	// Subscription routes
	subRouter := apiRouter.PathPrefix("/subscription").Subrouter()
	subRouter.Use(AuthenticationMiddleware)
	subRouter.HandleFunc("/", s.defaultHandler.CreateSubscription).Methods("POST")

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
