package api

import (
	"fmt"
	"net/http"
)

func (s *Server) handleMainPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Main page")
}