package api

import (
	"fmt"
	"net/http"
)

func (s *Server) handleMainPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Main page")
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Handle Login")
}