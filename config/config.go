package config

import (
	"errors"
	"log"
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

type ServerConfig struct {
	Addr            string
	WriteTimeout    time.Duration
	ReadTimeout     time.Duration
	IdleTimeout     time.Duration
	GracefulTimeout time.Duration
	SessionKey      string
}

func NewConfig() (*Config, error) {
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

	return &Config{
		Log:      NewLogConfig(),
		Database: database,
		Server:   server,
		Redis:    redisClient,
	}, nil
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
	err := godotenv.Load(".env")
	if err != nil {
		log.Println(err)
		log.Fatal("Error loading .env file")
	}

	if os.Getenv("APP_ENV") == "dev" {
		if err != nil {
			log.Println(err)
			log.Fatal("Error loading .env file")
		}
	}

	host := GetEnv("DB_HOST", "localhost")
	port, err := strconv.Atoi(GetEnv("DB_PORT", "5435"))
	if err != nil {
		return nil, errors.New("invalid DB_PORT")
	}
	user := GetEnv("DB_USER", "postgres")
	pass := GetEnv("DB_PASS", "postgres")
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
	host := GetEnv("REDIS_HOST", "redis:6381")
	pass := GetEnv("REDIS_PASS", "qwerty")
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
	addr := GetEnv("ADDR", "0.0.0.0:6969")
	writeTimeout, err := time.ParseDuration(GetEnv("write_timeout", "15s"))
	if err != nil {
		return nil, errors.New("invalid WRITE_TIMEOUT")
	}
	readTimeout, err := time.ParseDuration(GetEnv("read_timeout", "15s"))
	if err != nil {
		return nil, errors.New("invalid READ_TIMEOUT")
	}
	idleTimeout, err := time.ParseDuration(GetEnv("idle_timeout", "60s"))
	if err != nil {
		return nil, errors.New("invalid IDLE_TIMEOUT")
	}
	gracefulTimeout, err := time.ParseDuration(GetEnv("graceful_timeout", "5s"))
	if err != nil {
		return nil, errors.New("invalid GRACEFUL_TIMEOUT")
	}
	sessionKey := GetEnv("session_key", "super-secret")

	return &ServerConfig{
		Addr:            addr,
		GracefulTimeout: gracefulTimeout,
		WriteTimeout:    writeTimeout,
		ReadTimeout:     readTimeout,
		IdleTimeout:     idleTimeout,
		SessionKey:      sessionKey,
	}, nil
}
