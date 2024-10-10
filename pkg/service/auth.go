package service

import (
	"crypto/sha1"
	"fmt"
	"github.com/vladislavtinishov/todo-app"
	"github.com/vladislavtinishov/todo-app/pkg/repository"
	"github.com/vladislavtinishov/todo-app/utils"
)

const (
	salt = "asfdkj420cd0sdsdlk"
)

type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateUser(user todo.User) (int, error) {
	user.Password = generatePasswordHash(user.Password)
	return s.repo.CreateUser(user)
}

func (s *AuthService) GenerateToken(username, password string) (string, error) {
	user, err := s.repo.GetUser(username, generatePasswordHash(password))

	if err != nil {
		return "", err
	}

	return utils.GenerateJWT(user.Id)
}

func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
