package repository

import (
	"context"
	"fmt"
	"shiftkerja-backend/internal/core/entity"

	"github.com/jackc/pgx/v5"
)

type PostgresShiftRepo struct {
	DB *pgx.Conn
}

func NewPostgresShiftRepo(db *pgx.Conn) *PostgresShiftRepo {
	return &PostgresShiftRepo{DB: db}
}

func (r *PostgresShiftRepo) CreateShift(ctx context.Context, shift *entity.Shift) error {
	query := `
		INSERT INTO shifts (owner_id, title, description, pay_rate, lat, lng, status)
		VALUES ($1, $2, $3, $4, $5, $6, 'OPEN')
		RETURNING id, created_at
	`
	err := r.DB.QueryRow(ctx, query,
		shift.OwnerID,
		shift.Title,
		shift.Description,
		shift.PayRate,
		shift.Lat,
		shift.Lng,
	).Scan(&shift.ID, &shift.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to insert shift: %w", err)
	}
	return nil
}