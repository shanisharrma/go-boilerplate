package service

import (
	"github.com/shanisharrma/go-boilerplate/internal/app/core/repository"
	"github.com/shanisharrma/go-boilerplate/internal/domain/auth"
	"github.com/shanisharrma/go-boilerplate/internal/server"
)

type Services struct {
	Auth *auth.AuthService
	Job  server.JobRunner
}

func NewServices(s *server.Server, repos *repository.Repositories) (*Services, error) {
	authService := auth.NewAuthService(s)

	return &Services{
		Job:  s.Jobs,
		Auth: authService,
	}, nil
}
