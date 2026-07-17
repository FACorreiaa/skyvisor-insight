package repository

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/FACorreiaa/Aviation-tracker/app/auth"
	"github.com/FACorreiaa/Aviation-tracker/app/models"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"golang.org/x/oauth2"
)

const (
	RedisPrefix     = "user_session:"
	OIDCTokenPrefix = "oidc_token:"
	RandSize        = 32
	MaxAge          = time.Hour * 24 * 60
)

const userColumns = `
	user_id,
	username,
	email,
	password_hash,
	bio,
	image,
	created_at,
	updated_at
`

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

// Logout deletes the session token and any stored OIDC tokens from Redis.
func (a *AccountRepository) Logout(ctx context.Context, token Token) error {
	if err := a.redisClient.Del(ctx, token, OIDCTokenPrefix+token).Err(); err != nil {
		return errors.New("error deleting token")
	}
	return nil
}

// UpsertOIDCUser links a verified OIDC principal to a local user row. Match
// order: existing identity, then verified email (adopting pre-OIDC accounts),
// then a fresh row.
func (a *AccountRepository) UpsertOIDCUser(ctx context.Context, principal auth.Principal) (*models.UserSession, error) {
	if principal.Issuer == "" || principal.Subject == "" {
		return nil, errors.New("identity issuer and subject are required")
	}

	rows, err := a.pgpool.Query(ctx, `
		update "user"
		set email = coalesce(nullif($3, ''), email),
			image = coalesce(nullif($4, ''), image)
		where oidc_issuer = $1 and oidc_subject = $2
		returning `+userColumns,
		principal.Issuer, principal.Subject, strings.ToLower(principal.Email), principal.Picture,
	)
	if err != nil {
		slog.Error("look up OIDC identity", "error", err)
		return nil, errors.New("internal server error")
	}
	user, err := pgx.CollectOneRow(rows, pgx.RowToStructByPos[models.UserSession])
	if err == nil {
		return &user, nil
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		slog.Error("collect OIDC identity", "error", err)
		return nil, errors.New("internal server error")
	}

	if principal.EmailVerified && principal.Email != "" {
		rows, err := a.pgpool.Query(ctx, `
			update "user"
			set oidc_issuer = $1,
				oidc_subject = $2,
				image = coalesce(nullif($4, ''), image)
			where email = $3 and oidc_issuer is null
			returning `+userColumns,
			principal.Issuer, principal.Subject, strings.ToLower(principal.Email), principal.Picture,
		)
		if err != nil {
			slog.Error("link OIDC identity by email", "error", err)
			return nil, errors.New("internal server error")
		}
		user, err := pgx.CollectOneRow(rows, pgx.RowToStructByPos[models.UserSession])
		if err == nil {
			return &user, nil
		}
		if !errors.Is(err, pgx.ErrNoRows) {
			slog.Error("collect linked identity", "error", err)
			return nil, errors.New("internal server error")
		}
	}

	return a.insertOIDCUser(ctx, principal)
}

func (a *AccountRepository) insertOIDCUser(ctx context.Context, principal auth.Principal) (*models.UserSession, error) {
	if principal.Email == "" {
		return nil, errors.New("an email address is required to sign in; grant the email permission and try again")
	}
	base := usernameCandidate(principal)
	for attempt := 0; attempt < 5; attempt++ {
		username := base
		if attempt > 0 {
			username = fmt.Sprintf("%s-%s", base, identitySuffix(principal, attempt))
		}
		rows, err := a.pgpool.Query(ctx, `
			insert into "user" (username, email, oidc_issuer, oidc_subject, image)
			values ($1, $2, $3, $4, nullif($5, ''))
			returning `+userColumns,
			username, strings.ToLower(principal.Email), principal.Issuer, principal.Subject, principal.Picture,
		)
		if err == nil {
			user, collectErr := pgx.CollectOneRow(rows, pgx.RowToStructByPos[models.UserSession])
			if collectErr == nil {
				return &user, nil
			}
			err = collectErr
		}
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			// Another identity already resolved concurrently, or the generated
			// username collided; re-check the identity before retrying.
			existing, lookupErr := a.lookupByIdentity(ctx, principal)
			if lookupErr == nil {
				return existing, nil
			}
			continue
		}
		slog.Error("insert OIDC user", "error", err)
		return nil, errors.New("internal server error")
	}
	return nil, errors.New("unable to allocate a username")
}

func (a *AccountRepository) lookupByIdentity(ctx context.Context, principal auth.Principal) (*models.UserSession, error) {
	rows, err := a.pgpool.Query(ctx, `
		select `+userColumns+` from "user"
		where oidc_issuer = $1 and oidc_subject = $2
		limit 1
	`, principal.Issuer, principal.Subject)
	if err != nil {
		return nil, err
	}
	user, err := pgx.CollectOneRow(rows, pgx.RowToStructByPos[models.UserSession])
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// CreateSession issues a random session token in Redis for the user.
func (a *AccountRepository) CreateSession(ctx context.Context, userID uuid.UUID) (Token, error) {
	tokenBytes := make([]byte, RandSize)
	if _, err := rand.Read(tokenBytes); err != nil {
		slog.Error("generate session token", "error", err)
		return "", errors.New("internal server error")
	}
	token := fmt.Sprintf("%x", tokenBytes)
	if err := a.redisClient.Set(ctx, token, userID.String(), MaxAge).Err(); err != nil {
		slog.Error("store login session", "error", err)
		return "", errors.New("internal server error")
	}
	return token, nil
}

// StoreOIDCToken keeps the OAuth2 token set alongside the session so API
// calls can attach and refresh the user's access token.
func (a *AccountRepository) StoreOIDCToken(ctx context.Context, sessionToken Token, token *oauth2.Token) error {
	payload, err := json.Marshal(token)
	if err != nil {
		return errors.New("encode OIDC token")
	}
	if err := a.redisClient.Set(ctx, OIDCTokenPrefix+sessionToken, payload, MaxAge).Err(); err != nil {
		slog.Error("store OIDC token", "error", err)
		return errors.New("internal server error")
	}
	return nil
}

func (a *AccountRepository) LoadOIDCToken(ctx context.Context, sessionToken Token) (*oauth2.Token, error) {
	payload, err := a.redisClient.Get(ctx, OIDCTokenPrefix+sessionToken).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, errors.New("auth session expired")
		}
		slog.Error("read OIDC token", "error", err)
		return nil, errors.New("internal server error")
	}
	var token oauth2.Token
	if err := json.Unmarshal(payload, &token); err != nil {
		return nil, errors.New("decode OIDC token")
	}
	return &token, nil
}

func (m *MiddlewareRepository) UserFromSessionToken(ctx context.Context, token Token) (*models.UserSession, error) {
	userID, err := m.RedisClient.Get(ctx, token).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, errors.New("auth session expired")
		}

		slog.Error("read login session", "error", err)
		return nil, errors.New("internal server error")
	}

	rows, err := m.Pgpool.Query(ctx, `
		select `+userColumns+` from "user" where user_id = $1 limit 1
	`, userID)
	if err != nil {
		slog.Error("load session user", "error", err)
		return nil, errors.New("internal server error")
	}

	userWithToken, err := pgx.CollectOneRow(rows, pgx.RowToStructByPos[models.UserSession])
	if err != nil {
		return nil, errors.New("internal server error")
	}

	return &userWithToken, nil
}

func usernameCandidate(principal auth.Principal) string {
	for _, candidate := range []string{principal.Name, emailLocalPart(principal.Email)} {
		if slug := slugify(candidate); slug != "" {
			return slug
		}
	}
	return "traveler"
}

func emailLocalPart(email string) string {
	local, _, _ := strings.Cut(email, "@")
	return local
}

func slugify(value string) string {
	var builder strings.Builder
	for _, r := range strings.ToLower(strings.TrimSpace(value)) {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			builder.WriteRune(r)
		case r == ' ', r == '.', r == '-', r == '_':
			builder.WriteRune('-')
		}
	}
	slug := strings.Trim(builder.String(), "-")
	if len(slug) > 32 {
		slug = slug[:32]
	}
	return slug
}

func identitySuffix(principal auth.Principal, attempt int) string {
	sum := sha256.Sum256([]byte(fmt.Sprintf("%s|%s|%d", principal.Issuer, principal.Subject, attempt)))
	return fmt.Sprintf("%x", sum[:3])
}
