//package account
//
//import (
//	"context"
//	"crypto/rand"
//	"errors"
//	"fmt"
//	"github.com/jackc/pgerrcode"
//	"github.com/jackc/pgx/v5"
//	"github.com/jackc/pgx/v5/pgconn"
//	"golang.org/x/crypto/bcrypt"
//	"log/slog"
//)
//
//type RegisterForm struct {
//	Username        string `form:"username" validate:"required"`
//	Email           string `form:"email" validate:"required,email"`
//	Password        string `form:"password" validate:"required,min=8,max=72"`
//	PasswordConfirm string `form:"password_confirm" validate:"required,eqfield=Password"`
//}
//
//func (a *Accounts) RegisterNewAccount(ctx context.Context, form RegisterForm) (*Token, error) {
//	if err := a.validator.Struct(form); err != nil {
//		slog.Warn("Validation error")
//		return nil, err
//	}
//
//	passwordHash, err := bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.DefaultCost)
//	if err != nil {
//		slog.Error("Error hashing password", "err", err)
//		return nil, errors.New("internal server error")
//	}
//
//	var user User
//	var token Token
//
//	err = pgx.BeginFunc(ctx, a.pgpool, func(tx pgx.Tx) error {
//		row, _ := tx.Query(
//			ctx,
//			`
//			insert into "user" (username, email, password_hash)
//				values ($1, $2, $3)
//			returning
//				user_id,
//				username,
//				email,
//				password_hash,
//				bio,
//				image,
//				created_at,
//				updated_at
//			`,
//			form.Username,
//			form.Email,
//			passwordHash,
//		)
//		user, err = pgx.CollectOneRow(row, pgx.RowToStructByPos[User])
//		if err != nil {
//			return fmt.Errorf("error inserting user: %w", err)
//		}
//
//		tokenBytes := make([]byte, RAND_SIZE)
//		if _, err := rand.Read(tokenBytes); err != nil {
//			return fmt.Errorf("error generating token: %w", err)
//		}
//		token = Token(fmt.Sprintf("%x", tokenBytes))
//
//		if _, err := tx.Exec(
//			ctx,
//			`insert into user_token (user_id, token, context) values ($1, $2, $3)`,
//			user.ID,
//			token,
//			"auth",
//		); err != nil {
//			return fmt.Errorf("error inserting token: %w", err)
//		}
//
//		return nil
//	})
//
//	if err != nil {
//		var pgErr *pgconn.PgError
//		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
//			return nil, errors.New("username or email already taken")
//		}
//
//		slog.Error("Error creating account", "err", err)
//		return nil, errors.New("internal server error")
//	}
//
//	slog.Info("Created account", "user_id", user.ID)
//	return &token, nil
//}

package account

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
)

type RegisterForm struct {
	Username        string `form:"username" validate:"required"`
	Email           string `form:"email" validate:"required,email"`
	Password        string `form:"password" validate:"required,min=8,max=72"`
	PasswordConfirm string `form:"password_confirm" validate:"required,eqfield=Password"`
}

func (a *Accounts) RegisterNewAccount(ctx context.Context, form RegisterForm) (*Token, error) {
	if err := a.validator.Struct(form); err != nil {
		slog.Warn("Validation error")
		return nil, err
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("Error hashing password", "err", err)
		return nil, errors.New("internal server error")
	}

	var user User
	var token Token

	err = pgx.BeginFunc(ctx, a.pgpool, func(tx pgx.Tx) error {
		row, _ := tx.Query(
			ctx,
			`
			insert into "user" (username, email, password_hash)
				values ($1, $2, $3)
			returning
				user_id,
				username,
				email,
				password_hash,
				bio,
				image,
				created_at,
				updated_at
			`,
			form.Username,
			form.Email,
			passwordHash,
		)
		user, err = pgx.CollectOneRow(row, pgx.RowToStructByPos[User])
		if err != nil {
			return fmt.Errorf("error inserting user: %w", err)
		}

		tokenBytes := make([]byte, RAND_SIZE)
		if _, err := rand.Read(tokenBytes); err != nil {
			return fmt.Errorf("error generating token: %w", err)
		}
		token = Token(fmt.Sprintf("%x", tokenBytes))

		// Store the session token in Redis
		//redisKey := fmt.Sprintf("user_session:%s", token)
		if err := a.redisClient.Set(ctx, token, user.ID, time.Hour*24*7).Err(); err != nil {
			return fmt.Errorf("error inserting token into Redis: %w", err)
		}

		return nil
	})

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return nil, errors.New("username or email already taken")
		}

		slog.Error("Error creating account", "err", err)
		return nil, errors.New("internal server error")
	}

	slog.Info("Created account", "user_id", user.ID)
	return &token, nil
}
