package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/choipopik/gRPC-SSO/internal/domain/model"
	"github.com/choipopik/gRPC-SSO/internal/lib/jwt"
	"github.com/choipopik/gRPC-SSO/internal/storage"
	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	log          *slog.Logger
	userCreater  UserCreater
	userProvider UserProvider
	appProvider  AppProvider
	tokenTTL     time.Duration
}

type UserCreater interface {
	CreateUser(ctx context.Context, email string,
		pHash []byte) (userID int64, err error)
}

type UserProvider interface {
	User(ctx context.Context, email string) (model.User, error)
	isAdmin(ctx context.Context, userID int64) (bool, error)
}

type AppProvider interface {
	App(ctx context.Context, appID int) (model.App, error)
}

var (
	ErrInvalidCreds = errors.New("invalid credentials")
	ErrInvalidAppID = errors.New("invalid app id")
	ErrUserExists   = errors.New("user already exists")
)

// New returns a new instance of Auth service.
func New(log *slog.Logger, userCreater UserCreater, userProvider UserProvider,
	appProvider AppProvider, tokenTTL time.Duration) *Auth {
	return &Auth{
		log:          log,
		userCreater:  userCreater,
		userProvider: userProvider,
		appProvider:  appProvider,
		tokenTTL:     tokenTTL,
	}
}

// Login checks if user is existing.
// If so and password is correct returns token, else returns error.
func (a *Auth) Login(ctx context.Context, email string, password string,
	appID int) (string, error) {
	const op = "auth.Login"

	log := a.log.With(
		slog.String("op", op),
	)

	log.Info("Logining user")

	user, err := a.userProvider.User(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			a.log.Warn("user not found", storage.ErrUserNotFound)

			return "", fmt.Errorf("%s: %w", op, ErrInvalidCreds)
		}
		log.Error("failed to get user")

		return "", fmt.Errorf("%s: %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(password)); err != nil {
		a.log.Info("invalid credentials")

		return "", fmt.Errorf("%s: %w", op, ErrInvalidCreds)
	}

	app, err := a.appProvider.App(ctx, appID)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user logged in successefuly")

	token, err := jwt.NewToken(user, app, a.tokenTTL)
	if err != nil {
		a.log.Error("failed to generate token")

		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}

// Register is creating new user in the system and if user
// with given username do not exist - returns user`s id.
func (a *Auth) RegisterUser(ctx context.Context, email string, password string) (int64, error) {
	const op = "auth.RegisterUser"

	log := a.log.With(
		slog.String("op", op),
	)

	log.Info("registering user")

	pHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("error generating password hash")

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := a.userCreater.CreateUser(ctx, email, pHash)
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			a.log.Warn("app not found", storage.ErrUserExists)

			return 0, fmt.Errorf("%s: %w", op, ErrUserExists)
		}

		a.log.Error("error creating new user")

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user successefuly registered")

	return id, nil
}

// IsAdmin checks if user is admin.
func (a *Auth) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	const op = "auth.IsAdmin"

	log := a.log.With(
		slog.String("op", op),
	)

	log.Info("checking if user is admin")

	isAdmin, err := a.userProvider.isAdmin(ctx, userID)
	if err != nil {
		if errors.Is(err, storage.ErrAppNotFounf) {
			a.log.Warn("app not found", storage.ErrAppNotFounf)
			return false, fmt.Errorf("%s: %w", op, ErrInvalidAppID)
		}
		return false, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("checked if user is admin", slog.Bool("is admin", isAdmin))

	return isAdmin, nil
}
