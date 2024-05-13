package config

import (
	"log"
	"log/slog"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
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

func GetProdEnv() bool {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	mode := os.Getenv("MODE")

	return mode == "production"
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

func NewDatabaseConfig() (*DatabaseConfig, error) {
	env := GetProdEnv()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	user := os.Getenv("DB_USER")
	host := os.Getenv("DB_HOST")
	pass := os.Getenv("DB_PASS")
	port := os.Getenv("DB_PORT")
	name := os.Getenv("DB_NAME")
	schema := os.Getenv("SCHEMA")
	if env {
		connURL := os.Getenv("DB_PG_PROD")
		return &DatabaseConfig{
			ConnectionURL: connURL,
		}, nil
	}

	println(user)
	println(host)
	println(pass)

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
		Host:     host + ":" + port,
		Path:     name,
		RawQuery: query.Encode(),
	}
	println("DB")
	println(connURL.String())
	return &DatabaseConfig{
		ConnectionURL: connURL.String(),
	}, nil
}

func NewRedisConfig() (*RedisConfig, error) {
	var host, pass string

	env := GetProdEnv()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	if env {
		opt, err := redis.ParseURL(os.Getenv("UPSTASH_URL"))
		if err != nil {
			return nil, err
		}
		host = opt.Addr
		pass = opt.Password
	} else {
		host = os.Getenv("REDIS_HOST")
		pass = os.Getenv("REDIS_PASSWORD")
	}
	println(host)
	println(pass)

	return &RedisConfig{
		Host:     host,
		Password: pass,
		DB:       0,
	}, nil
}

func NewServerConfig() (*ServerConfig, error) {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get environment variable values
	addr := os.Getenv("SERVER_ADDR")
	port := os.Getenv("SERVER_PORT")
	gtStr := os.Getenv("SERVER_GRACEFUL_TIMEOUT")
	wtStr := os.Getenv("SERVER_WRITE_TIMEOUT")
	rtStr := os.Getenv("SERVER_READ_TIMEOUT")
	itStr := os.Getenv("SERVER_IDLE_TIMEOUT")
	sessionKey := os.Getenv("session_key")

	// Convert string values to integers
	gt, _ := strconv.Atoi(gtStr)
	wt, _ := strconv.Atoi(wtStr)
	rt, _ := strconv.Atoi(rtStr)
	it, _ := strconv.Atoi(itStr)

	// Convert integers to time.Duration
	gracefulTimeout := time.Duration(gt) * time.Second
	writeTimeout := time.Duration(wt) * time.Second
	readTimeout := time.Duration(rt) * time.Second
	idleTimeout := time.Duration(it) * time.Second

	return &ServerConfig{
		Addr:            addr + ":" + port,
		GracefulTimeout: gracefulTimeout,
		WriteTimeout:    writeTimeout,
		ReadTimeout:     readTimeout,
		IdleTimeout:     idleTimeout,
		SessionKey:      sessionKey,
	}, nil
}
