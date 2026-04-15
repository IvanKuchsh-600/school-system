package postgres

import (
	"auth-service/internal/entities"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(connStr string) (*UserRepository, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("open connection: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("ping database: %w", err)
	}

	return &UserRepository{db: db}, nil
}

// Create - добавить нового пользователя
func (r *UserRepository) Create(user *entities.User) error {
	query := `
       INSERT INTO users (email, password_hash, role)
       VALUES ($1, $2, $3)
       RETURNING id
   `

	var id int64
	err := r.db.QueryRow(query, user.Email, user.PasswordHash, user.Role).Scan(&id)
	if err != nil {
		return fmt.Errorf("create user: %w", entities.ErrDatabaseOperation)
	}

	return nil
}

func (r *UserRepository) GetByEmail(email string) (*entities.User, error) {
	query := `
	       SELECT email, password_hash, role
	       FROM users
	       WHERE email = $1
	   `

	var user entities.User
	err := r.db.QueryRow(query, email).Scan(
		&user.Email,
		&user.PasswordHash,
		&user.Role,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("get user by email: %w", entities.ErrDatabaseOperation)
	}

	return &user, nil
}

//
//// GetByID - получить пользователя по ID
//func (r *UserRepository) GetByID(id int64) (*models.User, error) {
//	query := `
//        SELECT id, email, password_hash, role, created_at
//        FROM users
//        WHERE id = $1
//    `
//
//	var user models.User
//	err := r.db.QueryRow(query, id).Scan(
//		&user.ID,
//		&user.Email,
//		&user.PasswordHash,
//		&user.Role,
//		&user.CreatedAt,
//	)
//
//	if err != nil {
//		if err == sql.ErrNoRows {
//			return nil, nil // пользователь не найден
//		}
//		return nil, fmt.Errorf("get user by id: %w", err)
//	}
//
//	return &user, nil
//}
//
// GetByEmail - получить пользователя по email
//
//// Update - обновить пользователя
//func (r *UserRepository) Update(user *models.User) error {
//	query := `
//        UPDATE users
//        SET email = $1, password_hash = $2, role = $3
//        WHERE id = $4
//    `
//
//	result, err := r.db.Exec(query, user.Email, user.PasswordHash, user.Role, user.ID)
//	if err != nil {
//		return fmt.Errorf("update user: %w", err)
//	}
//
//	rowsAffected, err := result.RowsAffected()
//	if err != nil {
//		return fmt.Errorf("get rows affected: %w", err)
//	}
//
//	if rowsAffected == 0 {
//		return fmt.Errorf("user with id %d not found", user.ID)
//	}
//
//	return nil
//}
//
//// Delete - удалить пользователя (мягкое удаление, если есть поле deleted_at)
//func (r *UserRepository) Delete(id int64) error {
//	query := `DELETE FROM users WHERE id = $1`
//
//	result, err := r.db.Exec(query, id)
//	if err != nil {
//		return fmt.Errorf("delete user: %w", err)
//	}
//
//	rowsAffected, err := result.RowsAffected()
//	if err != nil {
//		return fmt.Errorf("get rows affected: %w", err)
//	}
//
//	if rowsAffected == 0 {
//		return fmt.Errorf("user with id %d not found", id)
//	}
//
//	return nil
//}
//
//// ListAll - получить всех пользователей
//func (r *UserRepository) ListAll() ([]*models.User, error) {
//	query := `
//        SELECT id, email, password_hash, role, created_at
//        FROM users
//        ORDER BY id
//    `
//
//	rows, err := r.db.Query(query)
//	if err != nil {
//		return nil, fmt.Errorf("list users: %w", err)
//	}
//	defer rows.Close()
//
//	var users []*models.User
//	for rows.Next() {
//		var user models.User
//		err := rows.Scan(
//			&user.ID,
//			&user.Email,
//			&user.PasswordHash,
//			&user.Role,
//			&user.CreatedAt,
//		)
//		if err != nil {
//			return nil, fmt.Errorf("scan user: %w", err)
//		}
//		users = append(users, &user)
//	}
//
//	if err = rows.Err(); err != nil {
//		return nil, fmt.Errorf("rows error: %w", err)
//	}
//
//	return users, nil
//}
//
//// GetByRole - получить пользователей по роли
//func (r *UserRepository) GetByRole(role string) ([]*models.User, error) {
//	query := `
//        SELECT id, email, password_hash, role, created_at
//        FROM users
//        WHERE role = $1
//        ORDER BY id
//    `
//
//	rows, err := r.db.Query(query, role)
//	if err != nil {
//		return nil, fmt.Errorf("get users by role: %w", err)
//	}
//	defer rows.Close()
//
//	var users []*models.User
//	for rows.Next() {
//		var user models.User
//		err := rows.Scan(
//			&user.ID,
//			&user.Email,
//			&user.PasswordHash,
//			&user.Role,
//			&user.CreatedAt,
//		)
//		if err != nil {
//			return nil, fmt.Errorf("scan user: %w", err)
//		}
//		users = append(users, &user)
//	}
//
//	return users, nil
//}
