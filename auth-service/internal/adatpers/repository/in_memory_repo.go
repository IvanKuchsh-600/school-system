package repository

import (
	"auth-service/internal/entities"
	"sync"
)

type InMemoryUserRepo struct {
	users  map[int64]*entities.User
	nextID int64
	mu     sync.RWMutex
}

// NewInMemoryUserRepo - конструктор адаптера
func NewInMemoryUserRepo() *InMemoryUserRepo {
	return &InMemoryUserRepo{
		users:  make(map[int64]*entities.User),
		nextID: 1,
	}
}

// Create реализует ports.UserRepository
func (r *InMemoryUserRepo) Create(user *entities.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	user.ID = r.nextID
	r.nextID++
	r.users[user.ID] = user
	return nil
}

// FindByEmail реализует ports.UserRepository
func (r *InMemoryUserRepo) FindByEmail(email string) (*entities.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, u := range r.users {
		if u.Email == email {
			// Возвращаем копию, чтобы избежать изменений извне
			return r.copyUser(u), nil
		}
	}
	return nil, nil
}

// FindByID реализует ports.UserRepository
func (r *InMemoryUserRepo) FindByID(id int64) (*entities.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.users[id]
	if !exists {
		return nil, nil
	}
	return r.copyUser(user), nil
}

// Delete реализует ports.UserRepository
func (r *InMemoryUserRepo) Delete(id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.users, id)
	return nil
}

// copyUser создает копию пользователя
func (r *InMemoryUserRepo) copyUser(u *entities.User) *entities.User {
	return &entities.User{
		ID:           u.ID,
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		Role:         u.Role,
	}
}
