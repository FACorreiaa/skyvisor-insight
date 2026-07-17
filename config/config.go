package config

import (
	"errors"
	"log/slog"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Log      *LogConfig
	Database *DatabaseConfig
	Redis    *RedisConfig
	Server   *ServerConfig
	OIDC     *OIDCConfig
	API      *APIConfig
}

// APIConfig points the web app at skyvisor-api. Nil disables API-backed
// features (trips, assistant, billing).
type APIConfig struct {
	BaseURL string
}

type LogConfig struct {
	Level  slog.Level
	Format string
}

type DatabaseConfig struct {
	ConnectionURL string
}

type RedisConfig struct {
	Host     string
	Password string
	DB       int
}

// OIDCConfig configures the Authorization Code + PKCE login flow. When it is
// absent the web app boots, but signing in is unavailable.
type OIDCConfig struct {
	IssuerURL    string
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Audience     string
}

type ServerConfig struct {
	Addr            string
	WriteTimeout    time.Duration
	ReadTimeout     time.Duration
	IdleTimeout     time.Duration
	GracefulTimeout time.Duration
	SessionKey      string
	CookieSecure    bool
}

func NewConfig() (*Config, error) {
	if err := LoadEnvironment(); err != nil {
		return nil, err
	}

	database, err := NewDatabaseConfig()
	if err != nil {
		return nil, err
	}

	server, err := NewServerConfig()
	if err != nil {
		return nil, err
	}

	redisClient, err := NewRedisConfig()
	if err != nil {
		return nil, err
	}

	oidc, err := NewOIDCConfig()
	if err != nil {
		return nil, err
	}

	return &Config{
		Log:      NewLogConfig(),
		Database: database,
		Server:   server,
		Redis:    redisClient,
		OIDC:     oidc,
		API:      NewAPIConfig(),
	}, nil
}

func NewAPIConfig() *APIConfig {
	baseURL := GetEnv("SKYVISOR_API_URL", "")
	if baseURL == "" {
		return nil
	}
	return &APIConfig{BaseURL: baseURL}
}

// NewOIDCConfig returns nil without error when no OIDC variables are set, and
// fails when the configuration is only partially present.
func NewOIDCConfig() (*OIDCConfig, error) {
	cfg := &OIDCConfig{
		IssuerURL:    GetEnv("OIDC_ISSUER_URL", ""),
		ClientID:     GetEnv("OIDC_CLIENT_ID", ""),
		ClientSecret: GetEnv("OIDC_CLIENT_SECRET", ""),
		RedirectURL:  GetEnv("OIDC_REDIRECT_URL", ""),
		Audience:     GetEnv("OIDC_AUDIENCE", ""),
	}
	if cfg.IssuerURL == "" && cfg.ClientID == "" && cfg.RedirectURL == "" {
		return nil, nil
	}
	if cfg.IssuerURL == "" || cfg.ClientID == "" || cfg.RedirectURL == "" {
		return nil, errors.New("OIDC_ISSUER_URL, OIDC_CLIENT_ID, and OIDC_REDIRECT_URL must be set together")
	}
	issuer, err := url.Parse(cfg.IssuerURL)
	if err != nil || issuer.Scheme != "https" || issuer.Host == "" {
		return nil, errors.New("OIDC_ISSUER_URL must be an absolute HTTPS URL")
	}
	redirect, err := url.Parse(cfg.RedirectURL)
	if err != nil || (redirect.Scheme != "https" && redirect.Scheme != "http") || redirect.Host == "" {
		return nil, errors.New("OIDC_REDIRECT_URL must be an absolute URL")
	}
	return cfg, nil
}

// LoadEnvironment loads an optional local .env file. Production configuration
// remains environment-only and missing .env files are expected.
func LoadEnvironment() error {
	err := godotenv.Load(".env")
	if err == nil || errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return errors.New("load .env")
}

func GetEnv(key, defaultVal string) string {
	key = strings.ToUpper(key)
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultVal
}

func NewLogConfig() *LogConfig {
	var level slog.Level
	levelStr := GetEnv("LOG_LEVEL", "info")
	switch levelStr {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	format := GetEnv("LOG_FORMAT", "text")
	if format != "json" {
		format = "text"
	}

	return &LogConfig{
		Level:  level,
		Format: format,
	}
}

func NewDatabaseConfig() (*DatabaseConfig, error) {
	if connectionURL := GetEnv("DATABASE_URL", ""); connectionURL != "" {
		return &DatabaseConfig{ConnectionURL: connectionURL}, nil
	}

	host := GetEnv("DB_HOST", "localhost")
	port, err := strconv.Atoi(GetEnv("DB_PORT", "5435"))
	if err != nil {
		return nil, errors.New("invalid DB_PORT")
	}
	user := GetEnv("DB_USER", "postgres")
	pass := GetEnv("DB_PASS", "")
	if pass == "" {
		return nil, errors.New("DB_PASS is required")
	}
	dbname := GetEnv("DB_NAME", "aviation-tracker-dev")
	schema := GetEnv("DB_SCHEMA", "")

	query := url.Values{
		"sslmode":  []string{"disable"},
		"timezone": []string{"utc"},
	}
	if schema != "" {
		query.Add("search_path", schema)
	}
	connURL := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(user, pass),
		Host:     host + ":" + strconv.Itoa(port),
		Path:     dbname,
		RawQuery: query.Encode(),
	}
	return &DatabaseConfig{
		ConnectionURL: connURL.String(),
	}, nil
}

func NewRedisConfig() (*RedisConfig, error) {
	host := GetEnv("REDIS_HOST", "127.0.0.1:6381")
	pass := GetEnv("REDIS_PASS", "")
	// rdb := redis.NewClient(&redis.Options{
	//	Addr:     host,
	//	Password: pass, // no password set
	//	DB:       0,    // use default DB
	// })

	return &RedisConfig{
		Host:     host,
		Password: pass,
		DB:       0,
	}, nil
}

func NewServerConfig() (*ServerConfig, error) {
	addr := GetEnv("ADDR", "127.0.0.1:6969")
	writeTimeout, err := time.ParseDuration(GetEnv("WRITE_TIMEOUT", "15s"))
	if err != nil {
		return nil, errors.New("invalid WRITE_TIMEOUT")
	}
	readTimeout, err := time.ParseDuration(GetEnv("READ_TIMEOUT", "15s"))
	if err != nil {
		return nil, errors.New("invalid READ_TIMEOUT")
	}
	idleTimeout, err := time.ParseDuration(GetEnv("IDLE_TIMEOUT", "60s"))
	if err != nil {
		return nil, errors.New("invalid IDLE_TIMEOUT")
	}
	gracefulTimeout, err := time.ParseDuration(GetEnv("GRACEFUL_TIMEOUT", "5s"))
	if err != nil {
		return nil, errors.New("invalid GRACEFUL_TIMEOUT")
	}
	sessionKey := GetEnv("SESSION_KEY", "")
	if len(sessionKey) < 32 {
		return nil, errors.New("SESSION_KEY must contain at least 32 characters")
	}
	cookieSecure, err := strconv.ParseBool(GetEnv("COOKIE_SECURE", "false"))
	if err != nil {
		return nil, errors.New("invalid COOKIE_SECURE")
	}

	return &ServerConfig{
		Addr:            addr,
		GracefulTimeout: gracefulTimeout,
		WriteTimeout:    writeTimeout,
		ReadTimeout:     readTimeout,
		IdleTimeout:     idleTimeout,
		SessionKey:      sessionKey,
		CookieSecure:    cookieSecure,
	}, nil
}
