package account

import (
	"context"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"time"
)

const (
	REDIS_PREFIX = "user_session:"
	RAND_SIZE    = 32
	MAX_AGE      = time.Hour * 24 * 60
)

type Token = string

type Accounts struct {
	pgpool      *pgxpool.Pool
	redisClient *redis.Client
	validator   *validator.Validate
}

func NewAccounts(
	pgpool *pgxpool.Pool,
	redisClient *redis.Client,
	validator *validator.Validate,

) *Accounts {
	return &Accounts{
		pgpool:      pgpool,
		redisClient: redisClient,
		validator:   validator,
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

// Logout deletes the user token from the Redis store
func (h *Accounts) Logout(ctx context.Context, token Token) error {
	//userKey := REDIS_PREFIX + string(token)

	// Check if the token exists
	exists, err := h.redisClient.Exists(ctx, token).Result()
	if err != nil {
		return fmt.Errorf("error checking token existence: %w", err)
	}

	if exists == 0 {
		// Token not found, consider it already logged out
		return nil
	}

	// Delete the token
	if err := h.redisClient.Del(ctx, token).Err(); err != nil {
		return fmt.Errorf("error deleting token: %w", err)
	}

	return nil
}
