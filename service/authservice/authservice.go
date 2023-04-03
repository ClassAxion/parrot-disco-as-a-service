package authservice

import (
	"context"
	"fmt"

	"github.com/ClassAxion/parrot-disco-as-a-service/internal/database"
	"github.com/ClassAxion/parrot-disco-as-a-service/internal/database/user"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	DB *database.Database
}

func New(db *database.Database) *Service {
	return &Service{
		DB: db,
	}
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	return string(bytes), err
}

func (service *Service) Register(ctx context.Context, name string, email string, password string) error {
	hashed, err := HashPassword(password)
	if err != nil {
		return err
	}

	return service.DB.UpdateQuery(ctx, "INSERT INTO public.user (name, email, password) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING", name, email, hashed)
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (service *Service) Login(ctx context.Context, email string, password string) (*user.User, error) {
	row := service.DB.QueryRowContext(ctx, "SELECT id, password FROM public.user WHERE email = $1", email)
	if err := row.Err(); err != nil {
		return nil, err
	}

	var user user.User

	if err := row.Scan(&user.ID, &user.Password); err != nil {
		return nil, err
	}

	if !CheckPasswordHash(password, user.Password) {
		return nil, fmt.Errorf("authentication failed")
	}

	return &user, nil
}
