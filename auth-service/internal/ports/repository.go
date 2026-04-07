package ports

import "auth-service/internal/entities"

// UserRepository - порт (интерфейс), который говорит:
// "мне нужен способ сохранять и находить пользователей"
// Реализацию (адаптер) предоставит внешний слой
type UserRepository interface {
	Create(user *entities.User) error
	FindByEmail(email string) (*entities.User, error)
	FindByID(id int64) (*entities.User, error)
	Delete(id int64) error
}
