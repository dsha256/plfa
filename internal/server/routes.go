package server

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (s *Server) routes() http.Handler {
	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", s.healthcheckHandler)

	router.HandlerFunc(http.MethodGet, "/v1/tables", s.listTablesHandler)

	return s.metrics(s.recoverPanic(s.enableCORS(s.rateLimit(router))))
}
