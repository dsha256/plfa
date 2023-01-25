package server

import (
	"net/http"
)

type systemInfo struct {
	Version string `json:"version"`
}

type envelopeHealthCheck struct {
	Status     string     `json:"status"`
	SystemInfo systemInfo `json:"system_info"`
}

// healthcheckHandler shows if the service is up or not
//
// @Summary Check if the service is up
// @Schemes
// @Description Check if the service is up
// @Accept json
// @Produce json
// @Success 200 {object} envelopeHealthCheck
// @Router /healthcheck [get]
func (s *Server) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	env := envelopeHealthCheck{
		Status: "available",
		SystemInfo: systemInfo{
			Version: version,
		},
	}
	err := s.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		s.serverErrorResponse(w, r, err)
	}
}
