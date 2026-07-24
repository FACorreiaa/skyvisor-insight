package db

import (
	"context"
	"crypto/md5" //nolint:gosec // content fingerprint for migration drift only
	"embed"
	"fmt"
	"io/fs"
	"log/slog"
	"sort"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	uuid "github.com/vgarvardt/pgx-google-uuid/v5"
)

//go:embed migrations/*.sql
var migrationFS embed.FS

const retries = 25

// migrationAdvisoryLock is an arbitrary 64-bit key used with pg_advisory_lock
// so concurrent app starts serialize schema changes.
const migrationAdvisoryLock int64 = 0x534b595649534f52 // "SKYVISOR" as hex-ish constant

// Init Init.
func Init(connectionURL string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(connectionURL)
	if err != nil {
		return nil, err
	}
	cfg.AfterConnect = func(_ context.Context, conn *pgx.Conn) error {
		uuid.Register(conn.TypeMap())
		return nil
	}

	return pgxpool.NewWithConfig(context.Background(), cfg)
}

func InitRedis(host, password string, db int) (*redis.Client, error) {
	return redis.NewClient(&redis.Options{
		Addr:     host,
		Password: password,
		DB:       db,
	}), nil
}

// Migrate applies embedded SQL files under migrations/ that have not been
// recorded in _migrations. Already-applied files must not change content
// (MD5 fingerprint). Concurrent callers are serialized with an advisory lock.
func Migrate(ctx context.Context, conn *pgxpool.Pool) error {
	if ctx == nil {
		ctx = context.Background()
	}

	slog.Info("Running migrations")

	// Hold the lock for the whole run so two processes cannot interleave
	// "check applied" and "apply file".
	if _, err := conn.Exec(ctx, `select pg_advisory_lock($1)`, migrationAdvisoryLock); err != nil {
		return fmt.Errorf("acquire migration lock: %w", err)
	}
	defer func() {
		if _, err := conn.Exec(context.Background(), `select pg_advisory_unlock($1)`, migrationAdvisoryLock); err != nil {
			slog.Error("release migration lock", "error", err)
		}
	}()

	files, err := fs.ReadDir(migrationFS, "migrations")
	if err != nil {
		return fmt.Errorf("list migrations: %w", err)
	}
	names := make([]string, 0, len(files))
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		names = append(names, file.Name())
	}
	sort.Strings(names)

	if _, err = conn.Exec(ctx, `
		create table if not exists _migrations (
			name text primary key,
			hash text not null,
			created_at timestamp default now()
		);
	`); err != nil {
		return fmt.Errorf("create _migrations table: %w", err)
	}

	rows, err := conn.Query(ctx, `select name, hash from _migrations`)
	if err != nil {
		return fmt.Errorf("read applied migrations: %w", err)
	}
	appliedMigrations := make(map[string]string)
	var name, hash string
	_, err = pgx.ForEachRow(rows, []any{&name, &hash}, func() error {
		appliedMigrations[name] = hash
		return nil
	})
	if err != nil {
		return fmt.Errorf("scan applied migrations: %w", err)
	}

	applied, skipped := 0, 0
	for _, fileName := range names {
		contents, err := migrationFS.ReadFile("migrations/" + fileName)
		if err != nil {
			return fmt.Errorf("read migration %s: %w", fileName, err)
		}

		contentHash := fmt.Sprintf("%x", md5.Sum(contents)) //nolint:gosec

		if prevHash, ok := appliedMigrations[fileName]; ok {
			if prevHash != contentHash {
				return fmt.Errorf(
					"migration %q was already applied but its contents changed (hash %s → %s); add a new migration file instead of editing applied SQL",
					fileName, prevHash, contentHash,
				)
			}
			skipped++
			slog.Info("migration already applied", "name", fileName)
			continue
		}

		err = pgx.BeginFunc(ctx, conn, func(tx pgx.Tx) error {
			if _, err = tx.Exec(ctx, string(contents)); err != nil {
				return fmt.Errorf("execute %s: %w", fileName, err)
			}
			if _, err := tx.Exec(ctx, `insert into _migrations (name, hash) values ($1, $2)`,
				fileName, contentHash); err != nil {
				return fmt.Errorf("record %s: %w", fileName, err)
			}
			return nil
		})
		if err != nil {
			return err
		}
		applied++
		slog.Info("migration applied", "name", fileName)
	}

	slog.Info("Migrations finished", "applied", applied, "skipped", skipped, "total", len(names))
	return nil
}

// WaitForDB Small hack to wait for database to start inside docker.
func WaitForDB(ctx context.Context, pgpool *pgxpool.Pool) error {
	var lastErr error
	for attempts := 1; attempts <= retries; attempts++ {
		if err := pgpool.Ping(ctx); err == nil {
			return nil
		} else {
			lastErr = err
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(time.Duration(attempts) * 100 * time.Millisecond):
		}
	}
	return fmt.Errorf("database unavailable after %d attempts: %w", retries, lastErr)
}
