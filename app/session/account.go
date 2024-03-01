package session

import (
	"context"
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

const (
	RedisPrefix = "user_session:"
	RandSize    = 32
	MaxAge      = time.Hour * 24 * 60
)

type Token = string

type RepositoryAccount struct {
	pgpool      *pgxpool.Pool
	redisClient *redis.Client
	validator   *validator.Validate
	sessions    *sessions.CookieStore
}

func NewAccounts(
	pgpool *pgxpool.Pool,
	redisClient *redis.Client,
	validator *validator.Validate,
	sessions *sessions.CookieStore,

) *RepositoryAccount {
	return &RepositoryAccount{
		pgpool:      pgpool,
		redisClient: redisClient,
		validator:   validator,
		sessions:    sessions,
	}
}

type User struct {
	ID           uuid.UUID
	Username     string
	Email        string
	PasswordHash []byte
	Bio          string
	Image        *string
	CreatedAt    *time.Time
	UpdatedAt    *time.Time
}

type UserToken struct {
	Token     string
	CreatedAt *time.Time
	User      *User
}

// Logout deletes the user token from the Redis store.
func (a *RepositoryAccount) Logout(ctx context.Context, token Token) error {
	// userKey := REDIS_PREFIX + string(token)

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
