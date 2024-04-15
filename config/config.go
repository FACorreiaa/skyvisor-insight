package config

import (
	"errors"
	"log"
	"log/slog"
	"net/url"
	"os"
	"strconv"
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

func NewLogConfig() *LogConfig {
	var level slog.Level
	levelStr := os.Getenv("LOG_LEVEL")
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

	format := os.Getenv("LOG_FORMAT")
	if format != "json" {
		format = "text"
	}

	return &LogConfig{
		Level:  level,
		Format: format,
	}
}

// NewDatabaseConfig TODO VIBER CONFIG

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

	host := os.Getenv("DB_HOST")
	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		return nil, errors.New("invalid DB_PORT")
	}
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	dbname := os.Getenv("DB_NAME")
	schema := os.Getenv("")

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
	host := os.Getenv("REDIS_HOST")
	pass := os.Getenv("REDIS_PASS")
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
	err := godotenv.Load(".env.compose")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	addr := os.Getenv("ADDR")
	writeTimeout, err := time.ParseDuration(os.Getenv("write_timeout"))
	if err != nil {
		return nil, errors.New("invalid WRITE_TIMEOUT")
	}
	readTimeout, err := time.ParseDuration(os.Getenv("read_timeout"))
	if err != nil {
		return nil, errors.New("invalid READ_TIMEOUT")
	}
	idleTimeout, err := time.ParseDuration(os.Getenv("idle_timeout"))
	if err != nil {
		return nil, errors.New("invalid IDLE_TIMEOUT")
	}
	gracefulTimeout, err := time.ParseDuration(os.Getenv("graceful_timeout"))
	if err != nil {
		return nil, errors.New("invalid GRACEFUL_TIMEOUT")
	}
	sessionKey := os.Getenv("session_key")

	return &ServerConfig{
		Addr:            addr,
		GracefulTimeout: gracefulTimeout,
		WriteTimeout:    writeTimeout,
		ReadTimeout:     readTimeout,
		IdleTimeout:     idleTimeout,
		SessionKey:      sessionKey,
	}, nil
}
