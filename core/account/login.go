package account

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"log"
	"log/slog"
	"time"
)

type LoginForm struct {
	Email    string `form:"email" validate:"required,email"`
	Password string `form:"password" validate:"required"`
}

func (a *Accounts) Login(ctx context.Context, form LoginForm) (*Token, error) {
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
	user, err := pgx.CollectOneRow[User](rows, pgx.RowToStructByPos[User])
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

	tokenBytes := make([]byte, RAND_SIZE)
	if _, err := rand.Read(tokenBytes); err != nil {
		slog.Error("Error generating token", "err", err)
		return nil, errors.New("internal server error")
	}

	token := Token(fmt.Sprintf("%x", tokenBytes))
	log.Printf("Generated token: %s", token)

	//if _, err := a.pgpool.Exec(
	//	ctx,
	//	`
	//	insert into user_token (user_id, token, context)
	//	values ($1, $2, $3)
	//	`,
	//	user.ID,
	//	token,
	//	"auth",
	//); err != nil {
	//	slog.Error("Error inserting token", "err", err)
	//	return nil, errors.New("internal server error")
	//}

	// Store the session token in Redis
	//key := REDIS_PREFIX + string(token)
	err = a.redisClient.Set(ctx, string(token), (user.ID).String(), MAX_AGE).Err()
	if err != nil {
		log.Println("Error inserting token into Redis:", err)
		return nil, errors.New("internal server error")
	}

	log.Println("Token successfully inserted into Redis")
	return &token, nil
}

func (a *Accounts) UserFromSessionToken(ctx context.Context, token Token) (*User, error) {
	//key := REDIS_PREFIX + string(token)
	// Retrieve user ID from Redis
	fmt.Println("Retrieving user ID from Redis for token:", token)
	userID, err := a.redisClient.Get(ctx, token).Result()
	fmt.Println("Retrieved user ID from Redis:", userID)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, errors.New("auth session expired")
		}

		log.Println("Error querying user ID with token from Redis:", err)
		return nil, errors.New("internal server error")
	}

	// Retrieve user details from your data store (PostgreSQL in this case)
	rows, err := a.pgpool.Query(
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

	userWithToken, err := pgx.CollectOneRow(rows, pgx.RowToStructByPos[User])
	if err != nil {
		return nil, errors.New("internal server error")
	}

	// Check if the session has expired
	if userWithToken.CreatedAt == nil || time.Since(*userWithToken.CreatedAt) > MAX_AGE {
		return nil, errors.New("auth session expired")
	}

	return &userWithToken, nil
}
