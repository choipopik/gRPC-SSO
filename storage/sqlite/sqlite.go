package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/choipopik/gRPC-SSO/internal/domain/model"
	"github.com/choipopik/gRPC-SSO/internal/storage"
	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) CreateUser(ctx context.Context, email string, pHash []byte) (int64, error) {
	const op = "storage.sqlite.CreateUser"

	stmt, err := s.db.Prepare("INSERT INTO users(email, pass_hash) VALUES(?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.ExecContext(ctx, email, pHash)
	if err != nil {

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) User(ctx context.Context, email string) (model.User, error) {
	const op = "storage.sqlite.User"
	var user model.User

	stmt, err := s.db.Prepare("SELECT (id, email, pass_hash) FROM users WHERE email = ?")
	if err != nil {
		return model.User{}, fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRowContext(ctx, email)

	if err := row.Scan(&user.ID, &user.Email, &user.PasswordHash); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.User{}, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}

		return model.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (s *Storage) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	const op = "storage.sqlite.IsAdmin"
	var flag bool

	stmt, err := s.db.Prepare("SELECT is_admin FROM users WHERE id = ?")
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRowContext(ctx, userID)

	if err := row.Scan(&flag); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf("%s: %w", op, storage.ErrAppNotFounf)
		}

		return false, fmt.Errorf("%s: %w", op, err)
	}

	return flag, nil
}

func (s *Storage) App(ctx context.Context, appID int) (model.App, error) {
	const op = "storage.sqlite.User"
	var app model.App

	stmt, err := s.db.Prepare("SELECT (id, name, secret) FROM apps WHERE id = ?")
	if err != nil {
		return model.App{}, fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRowContext(ctx, appID)

	if err := row.Scan(&app.ID, &app.Name, &app.Secret); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.App{}, fmt.Errorf("%s: %w", op, storage.ErrAppNotFounf)
		}

		return model.App{}, fmt.Errorf("%s: %w", op, err)
	}

	return app, nil
}
