package auth

import (
	"context"
	"errors"
	repo "gary/ecom/internal/adapters/postgres/sqlc"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("Invalid Credentials")
	ErrUserAlreadyExists  = errors.New("User with this email already exists")
)

type Service interface {
	RegisterUser(ctx context.Context, tempUser UserRequest) (repo.User, error)
	Login(ctx context.Context, req UserRequest) (string, error)
}

type svc struct {
	repo *repo.Queries
	db   *pgxpool.Pool
	jwt  *JwtManager
}

func NewService(repo *repo.Queries, db *pgxpool.Pool, jwt *JwtManager) Service {
	return &svc{
		repo: repo,
		db:   db,
		jwt:  jwt,
	}
}

func (s *svc) RegisterUser(ctx context.Context, tempUser UserRequest) (repo.User, error) {
	//check if email already exists

	//hash the password using bcrypt
	hash, err := bcrypt.GenerateFromPassword([]byte(tempUser.Password), bcrypt.DefaultCost)
	if err != nil {
		return repo.User{}, err
	}
	hasheduser := repo.CreateUserParams{
		Email:        tempUser.Email,
		PasswordHash: string(hash),
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return repo.User{}, err
	}
	defer tx.Rollback(ctx)

	qtx := s.repo.WithTx(tx)

	createUser, err := qtx.CreateUser(ctx, hasheduser)

	if err != nil {
		return repo.User{}, err
	}

	tx.Commit(ctx)

	return createUser, nil

	//store the email and password using the sqlc function
}

func (s *svc) Login(ctx context.Context, req UserRequest) (string, error) {
	foundUser, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", ErrInvalidCredentials
		}
		return "", err
	}
	if err := s.comparePassword(foundUser.PasswordHash, req.Password); err != nil {
		return "", ErrInvalidCredentials
	}

	token, err := s.jwt.Generate(foundUser.ID)

	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *svc) comparePassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
