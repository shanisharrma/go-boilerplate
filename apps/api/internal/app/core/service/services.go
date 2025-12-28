package service

import (
	"github.com/shanisharrma/go-boilerplate/internal/app/core/repository"
	"github.com/shanisharrma/go-boilerplate/internal/app/server"
	"github.com/shanisharrma/go-boilerplate/internal/app/worker/job"
	"github.com/shanisharrma/go-boilerplate/internal/domain/auth"
)

type Services struct {
	Auth *auth.AuthService
	Job  *job.JobService
}

func NewServices(s *server.Server, repos *repository.Repositories) (*Services, error) {
	authService := auth.NewAuthService(s)

	return &Services{
		Job:  s.Job,
		Auth: authService,
	}, nil
}
