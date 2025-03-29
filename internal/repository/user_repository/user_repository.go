package user_repository

import (
	"context"
	"cult/internal/domain"
	"cult/internal/repository"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"log/slog"

	"github.com/jackc/pgx/v5"
)

type UserRepository struct {
	db  *pgx.Conn
	log *slog.Logger
}

func NewUserRepository(db *pgx.Conn, log *slog.Logger) *UserRepository {
	return &UserRepository{
		db:  db,
		log: log,
	}
}

// SaveUser implements UserSaver interface
func (r *UserRepository) SaveUser(ctx context.Context, phone string, passHash []byte) (uuid.UUID, error) {
	const op = "UserRepository.SaveUser"

	query := `
		INSERT INTO users (phone, pass_hash, user_type)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	var id uuid.UUID
	err := r.db.QueryRow(ctx, query, phone, string(passHash), domain.RegularUserType.String()).Scan(&id)
	if err != nil {
		if isUniqueViolation(err) {
			return uuid.Nil, fmt.Errorf("%s: %w", op, repository.ErrUserExists)
		}
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

// User implements UserProvider interface
func (r *UserRepository) User(ctx context.Context, phone string) (domain.User, error) {
	const op = "UserRepository.User"

	query := `
		SELECT id, phone, pass_hash
		FROM users
		WHERE phone = $1
	`

	var user domain.User
	err := r.db.QueryRow(ctx, query, phone).Scan(
		&user.ID,
		&user.Phone,
		&user.PassHash,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, fmt.Errorf("%s: %w", op, repository.ErrUserNotFound)
		}
		return domain.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

// IsAdmin implements UserProvider interface
func (r *UserRepository) IsAdmin(ctx context.Context, userID uuid.UUID) (bool, error) {
	const op = "UserRepository.IsAdmin"

	query := `
		SELECT EXISTS (
			SELECT 1 FROM users
			WHERE id = $1 AND user_type = 'admin'
		)
	`

	var isAdmin bool
	err := r.db.QueryRow(ctx, query, userID).Scan(&isAdmin)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return isAdmin, nil
}

// helper function to check for unique violation
func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505" // unique_violation
	}
	return false
}
