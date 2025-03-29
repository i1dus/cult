package auth

import (
	"context"
	"cult/internal/domain"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"log/slog"
	"time"

	"cult/internal/lib/jwt"
	"cult/internal/lib/logger/sl"
	"cult/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	log            *slog.Logger
	userRepository UserRepository
	tokenTTL       time.Duration
	Secret         string
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type UserRepository interface {
	SaveUser(ctx context.Context, id uuid.UUID, phone string, passHash []byte) (uuid.UUID, error)
	User(ctx context.Context, phone string) (domain.User, error)
}

func New(log *slog.Logger, userRepo UserRepository, tokenTTL time.Duration, secret string) *Auth {
	return &Auth{
		userRepository: userRepo,
		log:            log,
		tokenTTL:       tokenTTL,
		Secret:         secret,
	}
}

func (a *Auth) Login(ctx context.Context, phoneNumber string, password string) (uuid.UUID, string, error) {
	const op = "Auth.Login"

	log := a.log.With(
		slog.String("op", op),
		slog.String("phone_number", phoneNumber),
	)

	log.Info("attempting to login user")

	user, err := a.userRepository.User(ctx, phoneNumber)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			a.log.Warn("user not found", sl.Err(err))

			return uuid.Nil, "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		a.log.Error("failed to get user", sl.Err(err))

		return uuid.Nil, "", fmt.Errorf("%s: %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		a.log.Info("invalid credentials", sl.Err(err))

		return uuid.Nil, "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	log.Info("user logged in successfully")

	token, err := jwt.NewToken(user, a.Secret, a.tokenTTL)
	if err != nil {
		a.log.Error("failed to generate token", sl.Err(err))

		return uuid.Nil, "", fmt.Errorf("%s: %w", op, err)
	}

	return user.ID, token, nil
}

// RegisterNewUser registers new user in the system and returns user ID.
// If user with given username already exists, returns error.
func (a *Auth) RegisterNewUser(ctx context.Context, phoneNumber string, pass string) (uuid.UUID, error) {
	const op = "Auth.RegisterNewUser"

	log := a.log.With(
		slog.String("op", op),
		slog.String("phoneNumber", phoneNumber),
	)

	log.Info("registering user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate password hash", sl.Err(err))

		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	id, err := a.userRepository.SaveUser(ctx, uuid.New(), phoneNumber, passHash)
	if err != nil {
		log.Error("failed to save user", sl.Err(err))

		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}
