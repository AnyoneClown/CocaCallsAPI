package api

import (
	"fmt"
	"net/http"
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
	http.HandleFunc("/", s.handleMainPage)
	return http.ListenAndServe(s.listenAddr, nil)
}

func (s *Server) handleMainPage(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "Main page")
}