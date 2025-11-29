package repository

import (
	"context"
	"errors"
	"fmt"
	"shiftkerja-backend/internal/core/entity"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresShiftRepo struct {
	DB *pgxpool.Pool
}

func NewPostgresShiftRepo(db *pgxpool.Pool) *PostgresShiftRepo {
	return &PostgresShiftRepo{DB: db}
}

// CreateShift inserts a new shift into the database
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

// GetShiftByID retrieves a shift by its ID
func (r *PostgresShiftRepo) GetShiftByID(ctx context.Context, id int64) (*entity.Shift, error) {
	query := `
		SELECT id, owner_id, title, description, pay_rate, lat, lng, status, created_at
		FROM shifts
		WHERE id = $1
	`
	var shift entity.Shift
	err := r.DB.QueryRow(ctx, query, id).Scan(
		&shift.ID,
		&shift.OwnerID,
		&shift.Title,
		&shift.Description,
		&shift.PayRate,
		&shift.Lat,
		&shift.Lng,
		&shift.Status,
		&shift.CreatedAt,
	)
	
	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("shift not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get shift: %w", err)
	}
	
	return &shift, nil
}

// GetShiftsByOwner retrieves all shifts posted by a business owner
func (r *PostgresShiftRepo) GetShiftsByOwner(ctx context.Context, ownerID int64) ([]entity.Shift, error) {
	query := `
		SELECT id, owner_id, title, description, pay_rate, lat, lng, status, created_at
		FROM shifts
		WHERE owner_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.DB.Query(ctx, query, ownerID)
	if err != nil {
		return nil, fmt.Errorf("failed to query shifts: %w", err)
	}
	defer rows.Close()
	
	var shifts []entity.Shift
	for rows.Next() {
		var shift entity.Shift
		err := rows.Scan(
			&shift.ID,
			&shift.OwnerID,
			&shift.Title,
			&shift.Description,
			&shift.PayRate,
			&shift.Lat,
			&shift.Lng,
			&shift.Status,
			&shift.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan shift: %w", err)
		}
		shifts = append(shifts, shift)
	}
	
	return shifts, nil
}

// UpdateShiftStatus updates the status of a shift
func (r *PostgresShiftRepo) UpdateShiftStatus(ctx context.Context, id int64, status string) error {
	query := `UPDATE shifts SET status = $1 WHERE id = $2`
	_, err := r.DB.Exec(ctx, query, status, id)
	if err != nil {
		return fmt.Errorf("failed to update shift status: %w", err)
	}
	return nil
}

// ApplyForShift creates a new application
func (r *PostgresShiftRepo) ApplyForShift(ctx context.Context, shiftID, workerID int64) error {
	// Check if already applied
	checkQuery := `SELECT COUNT(*) FROM applications WHERE shift_id = $1 AND worker_id = $2`
	var count int
	err := r.DB.QueryRow(ctx, checkQuery, shiftID, workerID).Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check existing application: %w", err)
	}
	if count > 0 {
		return errors.New("you have already applied to this shift")
	}

	query := `
		INSERT INTO applications (shift_id, worker_id, status)
		VALUES ($1, $2, 'PENDING')
	`
	_, err = r.DB.Exec(ctx, query, shiftID, workerID)
	if err != nil {
		return fmt.Errorf("failed to submit application: %w", err)
	}
	return nil
}

// GetApplicationsByWorker retrieves all applications for a worker with shift details
func (r *PostgresShiftRepo) GetApplicationsByWorker(ctx context.Context, workerID int64) ([]entity.Application, error) {
	query := `
		SELECT 
			a.id, a.shift_id, a.worker_id, a.status, a.created_at,
			s.title, s.pay_rate
		FROM applications a
		JOIN shifts s ON a.shift_id = s.id
		WHERE a.worker_id = $1
		ORDER BY a.created_at DESC
	`
	rows, err := r.DB.Query(ctx, query, workerID)
	if err != nil {
		return nil, fmt.Errorf("failed to query applications: %w", err)
	}
	defer rows.Close()
	
	var applications []entity.Application
	for rows.Next() {
		var app entity.Application
		err := rows.Scan(
			&app.ID,
			&app.ShiftID,
			&app.WorkerID,
			&app.Status,
			&app.CreatedAt,
			&app.ShiftTitle,
			&app.ShiftPayRate,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan application: %w", err)
		}
		applications = append(applications, app)
	}
	
	return applications, nil
}

// GetApplicationsByShift retrieves all applications for a shift with worker details
func (r *PostgresShiftRepo) GetApplicationsByShift(ctx context.Context, shiftID int64) ([]entity.Application, error) {
	query := `
		SELECT 
			a.id, a.shift_id, a.worker_id, a.status, a.created_at,
			u.full_name, u.email
		FROM applications a
		JOIN users u ON a.worker_id = u.id
		WHERE a.shift_id = $1
		ORDER BY a.created_at ASC
	`
	rows, err := r.DB.Query(ctx, query, shiftID)
	if err != nil {
		return nil, fmt.Errorf("failed to query applications: %w", err)
	}
	defer rows.Close()
	
	var applications []entity.Application
	for rows.Next() {
		var app entity.Application
		err := rows.Scan(
			&app.ID,
			&app.ShiftID,
			&app.WorkerID,
			&app.Status,
			&app.CreatedAt,
			&app.WorkerName,
			&app.WorkerEmail,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan application: %w", err)
		}
		applications = append(applications, app)
	}
	
	return applications, nil
}

// UpdateApplicationStatus updates the status of an application
func (r *PostgresShiftRepo) UpdateApplicationStatus(ctx context.Context, applicationID int64, status string) error {
	query := `UPDATE applications SET status = $1 WHERE id = $2`
	_, err := r.DB.Exec(ctx, query, status, applicationID)
	if err != nil {
		return fmt.Errorf("failed to update application status: %w", err)
	}
	return nil
}

// GetApplicationByID retrieves an application by its ID
func (r *PostgresShiftRepo) GetApplicationByID(ctx context.Context, id int64) (*entity.Application, error) {
	query := `
		SELECT id, shift_id, worker_id, status, created_at
		FROM applications
		WHERE id = $1
	`
	var app entity.Application
	err := r.DB.QueryRow(ctx, query, id).Scan(
		&app.ID,
		&app.ShiftID,
		&app.WorkerID,
		&app.Status,
		&app.CreatedAt,
	)
	
	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("application not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get application: %w", err)
	}
	
	return &app, nil
}

// UpdateShift updates shift details
func (r *PostgresShiftRepo) UpdateShift(ctx context.Context, shift *entity.Shift) error {
	query := `
		UPDATE shifts
		SET title = $1, description = $2, pay_rate = $3, lat = $4, lng = $5, status = $6
		WHERE id = $7
		RETURNING id
	`
	var id int64
	err := r.DB.QueryRow(ctx, query,
		shift.Title,
		shift.Description,
		shift.PayRate,
		shift.Lat,
		shift.Lng,
		shift.Status,
		shift.ID,
	).Scan(&id)

	if err == pgx.ErrNoRows {
		return fmt.Errorf("shift not found")
	}
	if err != nil {
		return fmt.Errorf("failed to update shift: %w", err)
	}
	return nil
}

// DeleteShift deletes a shift by ID
func (r *PostgresShiftRepo) DeleteShift(ctx context.Context, id int64) error {
	// First delete all applications for this shift
	_, err := r.DB.Exec(ctx, "DELETE FROM applications WHERE shift_id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete applications: %w", err)
	}

	// Then delete the shift
	result, err := r.DB.Exec(ctx, "DELETE FROM shifts WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete shift: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("shift not found")
	}

	return nil
}

// DeleteApplication deletes an application by ID
func (r *PostgresShiftRepo) DeleteApplication(ctx context.Context, id int64) error {
	result, err := r.DB.Exec(ctx, "DELETE FROM applications WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete application: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("application not found")
	}

	return nil
}