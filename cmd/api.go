package main

import (
	repo "gary/ecom/internal/adapters/postgres/sqlc"
	"gary/ecom/internal/auth"
	"gary/ecom/internal/orders"
	"gary/ecom/internal/products"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"
)

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("healthy"))
	})
	//products
	productService := products.NewService(repo.New(app.db))
	productHandler := products.NewHandler(productService)
	products.RegisterRoutes(r, productHandler)

	//orders
	orderService := orders.NewService(repo.New(app.db), app.db)
	orderHandler := orders.NewHandler(orderService)
	//protected using auth middleware
	r.Group(func(r chi.Router) {
		r.Use(app.jwt.Authenticate)
		orders.RegisterRoutes(r, orderHandler)
	})

	//auth
	authService := auth.NewService(repo.New(app.db), app.db, app.jwt)
	authHandler := auth.NewHandler(authService)
	auth.RegisterRoutes(r, authHandler)

	return r
}

func (app *application) run(h http.Handler) error {
	srv := &http.Server{
		Addr:         ":" + app.config.addr,
		Handler:      h,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	log.Printf("Server has started at address: %s", app.config.addr)

	return srv.ListenAndServe()
}

type application struct {
	config config
	db     *pgx.Conn
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
