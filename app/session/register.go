package session

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/FACorreiaa/Aviation-tracker/app/models"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
)

func (a *AccountRepository) RegisterNewAccount(ctx context.Context, form models.RegisterForm) (*Token, error) {
	if err := a.validator.Struct(form); err != nil {
		slog.Warn("Validation error")
		return nil, err
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("Error hashing password", "err", err)
		return nil, errors.New("internal server error")
	}

	var user models.UserSession
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
		user, err = pgx.CollectOneRow(row, pgx.RowToStructByPos[models.UserSession])
		if err != nil {
			return errors.New("error inserting user")
		}

		tokenBytes := make([]byte, RandSize)
		if _, err = rand.Read(tokenBytes); err != nil {
			return errors.New("error generating token")
		}
		token = Token(fmt.Sprintf("%x", tokenBytes))

		// Store the session token in Redis
		// redisKey := fmt.Sprintf("user_session:%s", token)
		if err := a.redisClient.Set(ctx, token, user.ID, time.Hour*24*7).Err(); err != nil {
			return errors.New("error inserting token into Redis")
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
