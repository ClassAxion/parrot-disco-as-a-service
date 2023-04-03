package deployservice

import (
	"context"

	"github.com/ClassAxion/parrot-disco-as-a-service/internal/database"
)

type Service struct {
	DB *database.Database
}

func New(db *database.Database) *Service {
	return &Service{
		DB: db,
	}
}

func (service *Service) GetDeployIPByHash(ctx context.Context, hash string) (*string, error) {
	row := service.DB.QueryRowContext(ctx, "SELECT deployIP FROM public.user WHERE hash = $1", hash)
	if err := row.Err(); err != nil {
		return nil, err
	}

	var deployIP string

	if err := row.Scan(&deployIP); err != nil {
		return nil, err
	}

	return &deployIP, nil
}
