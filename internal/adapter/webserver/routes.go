package webserver

import (
	"github.com/anushasankaranarayanan/book-tracker-service/internal/adapter/webserver/swagger"
	"net/http"

	"github.com/anushasankaranarayanan/book-tracker-service/internal/adapter/webserver/probes"

	"github.com/gin-gonic/gin"
)

func (s *Server) Routes() error {
	r := gin.Default()

	r.Group("/api/v1").
		Use(gin.Logger()).
		POST("/book", s.AddBook).
		GET("/book/:id", s.GetBook).
		GET("/book", s.ListBooks).
		PUT("/book", s.UpdateBook).
		GET("/genre", s.GroupBooksByGenre)

	r.Group("/api/v1/probes").
		GET("/liveness", probes.Liveness)

	r.Group("/api/v1/openapi").
		GET("/", swagger.Build).
		GET("/:resource", swagger.Build)

	http.Handle("/", r)

	return nil
}
