package main

import (
	"context"
	"gary/ecom/internal/auth"
	"gary/ecom/internal/env"
	"gary/ecom/internal/pkg/logger"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	ctx := context.Background()

	cfg := config{
		addr: "8080",
		db: dbconfig{
			dsn: env.GetString("GOOSE_DBSTRING", "host=localhost user=postgres password=postgres dbname=ecom sslmode=disable"),
		},
		jwt: jwtConfig{
			secret: env.GetString("JWT_SECRET", "my-super-secret-jwt-string-for-authentication"),
		},
	}
	jwtManager := auth.NewJwtManager(cfg.jwt.secret)

	//Database

	pool, err := pgxpool.New(ctx, cfg.db.dsn)
	if err != nil {
		slog.Error(
			"failed to connect to database",
			"error", err,
		)

		os.Exit(1)
	}

	defer pool.Close()

	api := application{
		config: cfg,
		logger: logger.New(),
		db:     pool,
		jwt:    jwtManager,
	}

	if err := api.run(api.mount()); err != nil {
		slog.Error("Server failed to start", "error", err)
		os.Exit(1)
	}
}
