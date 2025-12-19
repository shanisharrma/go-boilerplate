package service

import "github.com/shanisharrma/go-boilerplate/internal/server"

type AuthService struct {
	server *server.Server
}

func NewAuthService(s *server.Server) *AuthService {
	return &AuthService{
		server: s,
	}
}
