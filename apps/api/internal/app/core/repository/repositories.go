package repository

import "github.com/shanisharrma/go-boilerplate/internal/app/server"

type Repositories struct {
}

func NewRepositories(s *server.Server) *Repositories {
	return &Repositories{}
}
