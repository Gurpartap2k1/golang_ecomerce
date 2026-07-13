package main

import (
	repo "gary/ecom/internal/adapters/postgres/sqlc"
	"gary/ecom/internal/auth"
	"gary/ecom/internal/middleware"
	"gary/ecom/internal/orders"
	"gary/ecom/internal/products"
	"log/slog"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func (app *application) mount() http.Handler {
	//r := chi.NewRouter()
	mux := http.NewServeMux()

	// A good base middleware stack

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("healthy"))
	})

	//products
	productService := products.NewService(repo.New(app.db))
	productHandler := products.NewHandler(productService, app.logger)
	products.RegisterRoutes(mux, productHandler)

	//orders
	orderService := orders.NewService(repo.New(app.db), app.db)
	orderHandler := orders.NewHandler(orderService, app.logger)

	//protected using auth middleware
	app.jwt.Authenticate(mux)
	orders.RegisterRoutes(mux, orderHandler, app.jwt.Authenticate)

	//auth
	authService := auth.NewService(repo.New(app.db), app.db, app.jwt)
	authHandler := auth.NewHandler(authService, app.logger)
	auth.RegisterRoutes(mux, authHandler)

	return mux
}

func (app *application) run(h http.Handler) error {
	handler := middleware.RequestLogger(app.logger)(h)
	srv := &http.Server{
		Addr:         ":" + app.config.addr,
		Handler:      handler,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	app.logger.Info("Server has started at address", "addr", app.config.addr)

	return srv.ListenAndServe()
}

type application struct {
	config config
	logger *slog.Logger
	db     *pgxpool.Pool
	jwt    *auth.JwtManager
}

type config struct {
	//server configuration
	addr string
	//database configuration
	db dbconfig

	jwt jwtConfig
}

type jwtConfig struct {
	secret string
}

type dbconfig struct {
	//some fields like user, password, host, port number, databse name
	dsn string
}
