package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/imnzr/sim-service-project/models"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserById(ctx context.Context, id uint) (*models.User, error)
	// UpdateUserEmail()
	// UpdateUserUsername()
}

type UserRepositoryImplementation struct {
	db *sql.DB
}

// GetUserById implements UserRepository.
func (u *UserRepositoryImplementation) GetUserById(ctx context.Context, id uint) (*models.User, error) {
	query := "SELECT id, username, email, password FROM `users` WHERE id = ?"

	rows, err := u.db.QueryContext(ctx, query, id)
	if err != nil {
		log.Printf("failed to execute query get user by id: %v", err)
		return nil, fmt.Errorf("failed to get user by id")
	}
	defer rows.Close()

	if rows.Next() {
		user := &models.User{}
		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.Password,
		)
		if err != nil {
			log.Printf("failed to scan user row: %v", err)
			return nil, fmt.Errorf("failed scan user")
		}
		return user, nil
	}
	return nil, nil
}

// GetUserByEmail implements UserRepository.
func (u *UserRepositoryImplementation) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := "SELECT id, username, email, password FROM `users` WHERE email = ?"

	rows, err := u.db.QueryContext(ctx, query, email)
	if err != nil {
		log.Printf("failed to execute query get user by email: %v", err)
		return nil, fmt.Errorf("failed to get user by email")
	}
	defer rows.Close()

	if rows.Next() {
		user := &models.User{}
		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.Password,
		)
		if err != nil {
			log.Printf("failed to scan user row: %v", err)
			return nil, fmt.Errorf("failed scan user")
		}
		return user, nil
	}
	return nil, nil
}

// CreateUser implements UserRepository.
func (u *UserRepositoryImplementation) CreateUser(ctx context.Context, user *models.User) error {
	query := "INSERT INTO users(username, email, password) VALUES(?,?,?)"

	result, err := u.db.ExecContext(ctx, query, user.Username, user.Email, user.Password)
	if err != nil {
		log.Println("failed to execute query create user :%w", err)
		return fmt.Errorf("failed to create user")
	}
	LastInsertID, err := result.LastInsertId()
	if err != nil {
		log.Println("failed to get last insert id: %w", err)
		return fmt.Errorf("failed to create user")
	}

	user.ID = uint(LastInsertID)
	return nil
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &UserRepositoryImplementation{
		db: db,
	}
}
