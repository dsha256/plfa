package server

import (
	"net/http"
)

// listTablesHandler lists all the existing pragmatic tables from the repository.
// Some listing's best practises (e.g. Pagination) are ignored cause of the task specific requirements - all existing
// data is required as a one big chunk.
func (s *Server) listTablesHandler(w http.ResponseWriter, r *http.Request) {
	pts, err := s.repo.ListTables()
	if err != nil {
		s.serverErrorResponse(w, r, err)
	}

	err = s.writeJSON(w, http.StatusOK, envelope{"tables": pts}, nil)
	if err != nil {
		s.serverErrorResponse(w, r, err)
	}
}
