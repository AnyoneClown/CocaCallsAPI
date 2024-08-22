package api

import (
	"net/http"

	"github.com/AnyoneClown/CocaCallsAPI/storage"
	"github.com/gorilla/mux"
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

    apiRouter := r.PathPrefix("/api").Subrouter()

    // Router for routes that do not require authentication
    publicRouter := apiRouter.PathPrefix("/").Subrouter()

    authRouter := publicRouter.PathPrefix("/auth").Subrouter()
    authRouter.HandleFunc("/register/", s.handleRegister).Methods("POST")

    jwtRouter := publicRouter.PathPrefix("/jwt").Subrouter()
    jwtRouter.HandleFunc("/create/", s.handleJWTCreate).Methods("POST")
    jwtRouter.HandleFunc("/verify/", s.handleJWTVerify).Methods("POST")

    // Router for routes that require authentication
    privateRouter := apiRouter.PathPrefix("/").Subrouter()
    privateRouter.Use(AuthenticationMiddleware)

    userRouter := privateRouter.PathPrefix("/user/").Subrouter()
    userRouter.HandleFunc("/me/", s.handleUserMe).Methods("GET")

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
