package handler

import (
	"github.com/shanisharrma/go-boilerplate/internal/app/core/service"
	"github.com/shanisharrma/go-boilerplate/internal/server"
)

type Handlers struct {
	Health  *HealthHandler
	OpenAPI *OpenAPIHandler
}

func NewHandlers(s *server.Server, services *service.Services) *Handlers {
	return &Handlers{
		Health:  NewHealthHandler(s),
		OpenAPI: NewOpenAPIHandler(s),
	}
}
