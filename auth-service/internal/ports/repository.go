package ports

import "auth-service/internal/entities"

type UserRepository interface {
	Create(user *entities.User) error
	GetByEmail(email string) (*entities.User, error)
}
