package server

import (
	"net/http"
)

func (s *Server) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	env := envelope{
		"status": "available",
		"system_info": map[string]string{
			"version": version,
		},
	}
	err := s.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		s.serverErrorResponse(w, r, err)
	}
}
