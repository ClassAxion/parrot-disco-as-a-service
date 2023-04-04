package userservice

import (
	"context"
	"encoding/json"

	"github.com/ClassAxion/parrot-disco-as-a-service/internal/database"
	"github.com/ClassAxion/parrot-disco-as-a-service/internal/database/user"
)

type Service struct {
	DB *database.Database
}

func New(db *database.Database) *Service {
	return &Service{
		DB: db,
	}
}

func (service *Service) GetUser(ctx context.Context, ID int) (*user.User, error) {
	row := service.DB.QueryRowContext(ctx, "SELECT name FROM public.user WHERE id = $1", ID)
	if err := row.Err(); err != nil {
		return nil, err
	}

	var user user.User

	if err := row.Scan(&user.Name); err != nil {
		return nil, err
	}

	return &user, nil
}

func (service *Service) GetSettings(ctx context.Context, ID int) (*user.User, error) {
	row := service.DB.QueryRowContext(ctx, "SELECT hash, zeroTierNetworkId, zeroTierDiscoIP, homeLocation, defaultRegion, deployStatus FROM public.user WHERE id = $1", ID)
	if err := row.Err(); err != nil {
		return nil, err
	}

	var user user.User

	var location *[]byte

	if err := row.Scan(&user.Hash, &user.ZeroTierNetworkId, &user.ZeroTierDiscoIP, &location, &user.DefaultRegion, &user.DeployStatus); err != nil {
		return nil, err
	}

	if location != nil {
		if err := json.Unmarshal(*location, &user.HomeLocation); err != nil {
			return nil, err
		}
	}

	return &user, nil
}

func (service *Service) SaveSettings(ctx context.Context, userID int, hash string, zeroTierNetworkId string, zeroTierDiscoIP string, homeLocation *user.Location) error {
	var data *[]byte

	if homeLocation != nil {
		val, err := json.Marshal(homeLocation)
		if err != nil {
			return err
		}

		data = &val
	}

	return service.DB.UpdateQuery(ctx, "UPDATE public.user SET hash = $2, zeroTierNetworkId = $3, zeroTierDiscoIP = $4, homeLocation = $5 WHERE id = $1", userID, hash, zeroTierNetworkId, zeroTierDiscoIP, data)
}

func (service *Service) UpdateDefaultRegion(ctx context.Context, userID int, defaultRegion string) error {
	return service.DB.UpdateQuery(ctx, "UPDATE public.user SET defaultRegion = $2 WHERE id = $1", userID, defaultRegion)
}

func (service *Service) StartDeploying(ctx context.Context, userID int) error {
	return service.DB.UpdateQuery(ctx, "UPDATE public.user SET deployStatus = 1 WHERE id = $1", userID)
}
