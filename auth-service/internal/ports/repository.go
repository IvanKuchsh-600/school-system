package ports

import "auth-service/internal/entities"

type UserRepository interface {
	Create(user *entities.User) error
	FindByEmail(email string) (*entities.User, error)
	FindByID(id int64) (*entities.User, error)
	Delete(id int64) error
}
