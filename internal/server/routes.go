package server

import (
	"net/http"

	_ "github.com/dsha256/plfa/docs/swagger"
	"github.com/julienschmidt/httprouter"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func (s *Server) routes() http.Handler {
	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", s.healthcheckHandler)

	router.HandlerFunc(http.MethodGet, "/v1/tables", s.listTablesHandler)

	router.HandlerFunc("GET", "/v1/swagger/*any", httpSwagger.WrapHandler)

	return s.metrics(s.recoverPanic(s.enableCORS(s.rateLimit(router))))
}
