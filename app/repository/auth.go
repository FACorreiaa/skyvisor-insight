package repository

import (
	"errors"
	"fmt"
	"log"
	"log/slog"
	"time"

	"context"

	"crypto/rand"

	"github.com/FACorreiaa/Aviation-tracker/app/models"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/sessions"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

const (
	RedisPrefix = "user_session:"
	RandSize    = 32
	MaxAge      = time.Hour * 24 * 60
)

type Token = string

type AccountRepository struct {
	pgpool      *pgxpool.Pool
	redisClient *redis.Client
	validator   *validator.Validate
	sessions    *sessions.CookieStore
}

func NewAccountRepository(db *pgxpool.Pool,
	redisClient *redis.Client,
	validator *validator.Validate,
	sessions *sessions.CookieStore,
) *AccountRepository {
	return &AccountRepository{
		pgpool:      db,
		redisClient: redisClient,
		validator:   validator,
		sessions:    sessions,
	}
}

// Logout deletes the user token from the Redis store.
func (a *AccountRepository) Logout(ctx context.Context, token Token) error {
	// userKey := RedisPrefix + string(token)

	// Check if the token exists
	exists, err := a.redisClient.Exists(ctx, token).Result()
	if err != nil {
		return errors.New("error checking token existence")
	}

	if exists == 0 {
		// Token not found, consider it already logged out
		return nil
	}

	// Delete the token
	if err = a.redisClient.Del(ctx, token).Err(); err != nil {
		return errors.New("error deleting token")
	}

	return nil
}

func (a *AccountRepository) Login(ctx context.Context, form models.LoginForm) (*Token, error) {
	if err := a.validator.Struct(form); err != nil {
		return nil, err
	}

	rows, _ := a.pgpool.Query(
		ctx,
		`
		select
			user_id,
			username,
			email,
			password_hash,
			bio,
			image,
			created_at,
			updated_at
		from "user" where email = $1 limit 1
		`,
		form.Email,
	)
	user, err := pgx.CollectOneRow[models.UserSession](rows, pgx.RowToStructByPos[models.UserSession])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("invalid email or password")
		}

		slog.Error("Error querying user", "err", err)
		return nil, errors.New("internal server error")
	}

	if err := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(form.Password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	tokenBytes := make([]byte, RandSize)
	if _, err = rand.Read(tokenBytes); err != nil {
		slog.Error("Error generating token", "err", err)
		return nil, errors.New("internal server error")
	}

	token := fmt.Sprintf("%x", tokenBytes)
	log.Printf("Generated token: %s", token)

	// if _, err := a.pgpool.Exec(
	//	ctx,
	//	`
	//	insert into user_token (user_id, token, context)
	//	values ($1, $2, $3)
	//	`,
	//	user.ID,
	//	token,
	//	"auth",
	// ); err != nil {
	//	slog.Error("Error inserting token", "err", err)
	//	return nil, errors.New("internal server error")
	//}

	// Store the session token in Redis
	// key := RedisPrefix + string(token)
	err = a.redisClient.Set(ctx, token, (user.ID).String(), MaxAge).Err()
	if err != nil {
		log.Println("Error inserting token into Redis:", err)
		return nil, errors.New("internal server error")
	}

	log.Println("Token successfully inserted into Redis")
	return &token, nil
}

func (m *MiddlewareRepository) UserFromSessionToken(ctx context.Context, token Token) (*models.UserSession, error) {
	// key := RedisPrefix + string(token)
	// Retrieve user ID from Redis
	log.Println("Retrieving user ID from Redis for token:", token)
	userID, err := m.RedisClient.Get(ctx, token).Result()
	log.Println("Retrieved user ID from Redis:", userID)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, errors.New("auth session expired")
		}

		log.Println("Error querying user ID with token from Redis:", err)
		return nil, errors.New("internal server error")
	}

	// Retrieve user details from your data store (PostgreSQL in this case)
	rows, err := m.Pgpool.Query(
		ctx,
		`
		select
			user_id,
			username,
			email,
			password_hash,
			bio,
			image,
			created_at,
			updated_at
		from "user" where user_id = $1 limit 1
		`,
		userID,
	)
	if err != nil {
		log.Println("Error querying user from PostgreSQL:", err)
		return nil, errors.New("internal server error")
	}

	userWithToken, err := pgx.CollectOneRow(rows, pgx.RowToStructByPos[models.UserSession])
	if err != nil {
		return nil, errors.New("internal server error")
	}

	// Check if the session has expired
	if userWithToken.CreatedAt == nil || time.Since(*userWithToken.CreatedAt) > MaxAge {
		return nil, errors.New("auth session expired")
	}

	return &userWithToken, nil
}

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
		token = fmt.Sprintf("%x", tokenBytes)

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
