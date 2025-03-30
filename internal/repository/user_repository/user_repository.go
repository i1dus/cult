package user_repository

import (
	"context"
	"cult/internal/domain"
	"cult/internal/repository"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"strings"

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
func (r *UserRepository) SaveUser(ctx context.Context, _ uuid.UUID, phone string, passHash []byte) (uuid.UUID, error) {
	const op = "UserRepository.SaveUser"

	query := `
		INSERT INTO users (phone, pass_hash, user_type)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	var userID uuid.UUID
	err := r.db.QueryRow(ctx, query, phone, string(passHash), domain.RegularUserType.String()).Scan(&userID)
	if err != nil {
		if isUniqueViolation(err) {
			return uuid.Nil, fmt.Errorf("%s: %w", op, repository.ErrUserExists)
		}
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	return userID, nil
}

// UserByID implements UserProvider interface
func (r *UserRepository) UserByID(ctx context.Context, userID uuid.UUID) (domain.User, error) {
	const op = "UserRepository.UserByID"

	query := `
        SELECT 
            id, 
            phone, 
            name,
            surname,
            patronymic,
            address,
            user_type,
            pass_hash
        FROM users
        WHERE id = $1
    `

	var user domain.User

	err := r.db.QueryRow(ctx, query, userID).Scan(
		&user.ID,
		&user.Phone,
		&user.Name,
		&user.Surname,
		&user.Patronymic,
		&user.Address,
		&user.UserType,
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

func (r *UserRepository) UserByPhone(ctx context.Context, phoneNumber string) (domain.User, error) {
	const op = "UserRepository.UserByPhone"

	query := `
        SELECT 
            id, 
            phone, 
            name,
            surname,
            patronymic,
            address,
            user_type,
            pass_hash
        FROM users
        WHERE phone = $1 
    `

	var user domain.User

	err := r.db.QueryRow(ctx, query, phoneNumber).Scan(
		&user.ID,
		&user.Phone,
		&user.Name,
		&user.Surname,
		&user.Patronymic,
		&user.Address,
		&user.UserType,
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

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505" // unique_violation
	}
	return false
}

// UpdateUser updates user information with optional fields
func (r *UserRepository) UpdateUser(ctx context.Context, userID uuid.UUID, update domain.UserUpdate) error {
	const op = "UserRepository.UpdateUser"

	query := "UPDATE users SET "
	params := []interface{}{}
	setClauses := []string{}
	paramCount := 1

	// Build SET clauses for provided fields
	if update.Name != nil {
		setClauses = append(setClauses, fmt.Sprintf("name = $%d", paramCount+1))
		params = append(params, *update.Name)
		paramCount++
	}
	if update.Surname != nil {
		setClauses = append(setClauses, fmt.Sprintf("surname = $%d", paramCount+1))
		params = append(params, *update.Surname)
		paramCount++
	}
	if update.Patronymic != nil {
		setClauses = append(setClauses, fmt.Sprintf("patronymic = $%d", paramCount+1))
		params = append(params, *update.Patronymic)
		paramCount++
	}
	if update.Phone != nil {
		setClauses = append(setClauses, fmt.Sprintf("phone = $%d", paramCount+1))
		params = append(params, *update.Phone)
		paramCount++
	}
	if update.Address != nil {
		setClauses = append(setClauses, fmt.Sprintf("address = $%d", paramCount+1))
		params = append(params, *update.Address)
		paramCount++
	}

	if len(setClauses) == 0 {
		return fmt.Errorf("%s: %w", op, repository.ErrNoFieldsToUpdate)
	}

	query += strings.Join(setClauses, ", ") + " WHERE id = $1"
	params = append([]interface{}{userID}, params...)

	tag, err := r.db.Exec(ctx, query, params...)
	if err != nil {
		if isUniqueViolation(err) {
			return fmt.Errorf("%s: %w", op, repository.ErrUserExists)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, repository.ErrUserNotFound)
	}

	return nil
}
