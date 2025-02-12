package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrNotFound         = errors.New("resource not found")
	ErrUniqueConstraint = errors.New("unique constraint violated")
	ErrForeignKey       = errors.New("foreign key constraint violated")
)

func HandleDBError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}

	var pqErr *pgconn.PgError
	if errors.As(err, &pqErr) {
		switch pqErr.Code {
		case "23505": // Unique violation
			return fmt.Errorf("%w: %s", ErrUniqueConstraint, pqErr.Detail)
		case "23503": // Foreign key violation
			return fmt.Errorf("%w: %s", ErrForeignKey, pqErr.Detail)
		}
	}

	// Return the original error if not handled specifically.
	return err
}
