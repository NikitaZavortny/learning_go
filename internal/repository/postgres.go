package repository

import (
	"auth-server/internal/model"
	"database/sql"
	"fmt"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(email, passwordHash string, activationlink string) (*model.User, error) {
	query := `INSERT INTO users (email, password_hash, activated, activation_link) VALUES ($1, $2, FALSE, $3) RETURNING id, email, created_at`

	user := &model.User{}
	err := r.db.QueryRow(query, email, passwordHash, activationlink).Scan(&user.ID, &user.Email, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) GetUserByEmail(email string) (*model.User, error) {
	query := `SELECT id, email, password_hash, created_at, activated FROM users WHERE email = $1`

	user := &model.User{}
	err := r.db.QueryRow(query, email).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.Activated)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) GetUserByID(id int) (*model.User, error) {
	query := `SELECT id, email, password_hash, created_at, activated FROM users WHERE id = $1`

	user := &model.User{}
	err := r.db.QueryRow(query, id).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.Activated)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) GetUserByLink(link string) (*model.User, error) {
	query := `SELECT id, email, password_hash, created_at, activated FROM users WHERE activation_link = $1`

	user := &model.User{}
	err := r.db.QueryRow(query, link).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.Activated)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) Activate(link string) (*model.User, error) {

	query := `UPDATE users SET activated = true WHERE activation_link = $1`
	fmt.Println("works")
	user := &model.User{}
	err := r.db.QueryRow(query, link).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.Activated)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}

	return user, nil
}
