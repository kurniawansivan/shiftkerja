package repository

import (
	"context"
	"fmt"
	"shiftkerja-backend/internal/core/entity"
	"github.com/jackc/pgx/v5"
)

type PostgresUserRepo struct {
	DB *pgx.Conn
}

func NewPostgresUserRepo(db *pgx.Conn) *PostgresUserRepo {
	return &PostgresUserRepo{DB: db}
}

func (r *PostgresUserRepo) CreateUser(ctx context.Context, user *entity.User) error {
	query := `
		INSERT INTO users (email, password_hash, full_name, role)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`
	// We scan the generated ID back into the struct
	err := r.DB.QueryRow(ctx, query, user.Email, user.PasswordHash, user.FullName, user.Role).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}
	return nil
}

func (r *PostgresUserRepo) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	query := `SELECT id, email, password_hash, full_name, role, created_at FROM users WHERE email = $1`
	
	var user entity.User
	err := r.DB.QueryRow(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.FullName, &user.Role, &user.CreatedAt,
	)
	
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // Not found
		}
		return nil, err
	}
	return &user, nil
}