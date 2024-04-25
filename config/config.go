package config

import (
	"fmt"
	"log"
	"log/slog"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

type StructConfig struct {
	Mode   string `mapstructure:"mode"`
	Dotenv string `mapstructure:"dotenv"`
	Server struct {
		Addr            string        `mapstructure:"addr"`
		Port            int           `mapstructure:"port"`
		WriteTimeout    time.Duration `mapstructure:"write_timeout"`
		ReadTimeout     time.Duration `mapstructure:"read_timeout"`
		IdleTimeout     time.Duration `mapstructure:"idle_timeout"`
		GracefulTimeout time.Duration `mapstructure:"graceful_timeout"`
	} `mapstructure:"server"`
	Postgres struct {
		DBName string `mapstructure:"name"`
		User   string `mapstructure:"user"`
		Port   int    `mapstructure:"port"`
		Host   string `mapstructure:"host"`
	} `mapstructure:"postgres"`
	Redis struct {
		Host string `mapstructure:"host"`
		DB   int    `mapstructure:"db"`
	} `mapstructure:"postgres"`
	Pprof struct {
		Addr string `mapstructure:"addr"`
		Port int    `mapstructure:"port"`
	}
}

func InitConfig() (StructConfig, error) {
	var config StructConfig
	v := viper.New()

	v.AddConfigPath("config")
	v.AddConfigPath("app/config")
	v.AddConfigPath("entrypoint/config")

	v.SetConfigName("config_dev")
	fmt.Println("Configuration paths:", v.ConfigFileUsed())

	if err := v.ReadInConfig(); err != nil {
		return config, err
	}
	if err := v.Unmarshal(&config); err != nil {
		return config, err
	}
	return config, nil
}

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
	err := godotenv.Load(".env.compose")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	env := os.Getenv("APP_ENV")
	if env == "production" {
		return true
	}

	return false
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
	err := godotenv.Load(".env.compose")
	if err != nil {
		log.Println(err)
		log.Fatal("Error loading .env file")
	}

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
	err := godotenv.Load(".env.compose")
	env := GetProdEnv()

	if env {
		connURL := os.Getenv("DB_PG_PROD")
		return &DatabaseConfig{
			ConnectionURL: connURL,
		}, nil
	}
	if err != nil {
		log.Println(err)
		log.Fatal("Error loading .env file")
	}

	cfg, err := InitConfig()
	if err != nil {
		log.Println(err)
		log.Fatal("Error loading yml config")
	}

	pass := os.Getenv("DB_PASS")
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
		User:     url.UserPassword(cfg.Postgres.User, pass),
		Host:     cfg.Postgres.Host + ":" + strconv.Itoa(cfg.Postgres.Port),
		Path:     cfg.Postgres.DBName,
		RawQuery: query.Encode(),
	}
	println("DB")
	println(connURL.String())
	return &DatabaseConfig{
		ConnectionURL: connURL.String(),
	}, nil
}

func NewRedisConfig() (*RedisConfig, error) {
	err := godotenv.Load(".env.compose")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var host, pass string

	env := GetProdEnv()

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

	return &RedisConfig{
		Host:     host,
		Password: pass,
		DB:       0,
	}, nil
}

func NewServerConfig() (*ServerConfig, error) {
	env := GetProdEnv()
	println(env)
	// TODO prepare app for prod

	err := godotenv.Load(".env.compose")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg, err := InitConfig()
	if err != nil {
		log.Println(err)
		log.Fatal("Error loading yml config")
	}
	sessionKey := os.Getenv("session_key")

	return &ServerConfig{
		Addr:            cfg.Server.Addr + ":" + strconv.Itoa(cfg.Server.Port),
		GracefulTimeout: cfg.Server.GracefulTimeout,
		WriteTimeout:    cfg.Server.WriteTimeout,
		ReadTimeout:     cfg.Server.ReadTimeout,
		IdleTimeout:     cfg.Server.IdleTimeout,
		SessionKey:      sessionKey,
	}, nil
}
