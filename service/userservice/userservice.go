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
	row := service.DB.QueryRowContext(ctx, "SELECT hash, zeroTierNetworkId, zeroTierDiscoIP, homeLocation, deployRegion, deployStatus, share_location FROM public.user WHERE id = $1", ID)
	if err := row.Err(); err != nil {
		return nil, err
	}

	var user user.User

	var location *[]byte

	if err := row.Scan(&user.Hash, &user.ZeroTierNetworkId, &user.ZeroTierDiscoIP, &location, &user.DeployRegion, &user.DeployStatus, &user.ShareLocation); err != nil {
		return nil, err
	}

	if location != nil {
		if err := json.Unmarshal(*location, &user.HomeLocation); err != nil {
			return nil, err
		}
	}

	return &user, nil
}

func (service *Service) SaveSettings(ctx context.Context, userID int, hash string, zeroTierNetworkId string, zeroTierDiscoIP string, homeLocation *user.Location, shareLocation bool) error {
	var data *[]byte

	if homeLocation != nil {
		val, err := json.Marshal(homeLocation)
		if err != nil {
			return err
		}

		data = &val
	}

	return service.DB.UpdateQuery(ctx, "UPDATE public.user SET hash = $2, zeroTierNetworkId = $3, zeroTierDiscoIP = $4, homeLocation = $5, share_location = $6 WHERE id = $1", userID, hash, zeroTierNetworkId, zeroTierDiscoIP, data, shareLocation)
}

func (service *Service) StartDeploying(ctx context.Context, userID int, deployRegion string) error {
	return service.DB.UpdateQuery(ctx, "UPDATE public.user SET deployStatus = 1, deployRegion = $2 WHERE id = $1", userID, deployRegion)
}

func (service *Service) Stop(ctx context.Context, userID int) error {
	return service.DB.UpdateQuery(ctx, "UPDATE public.user SET deployStatus = 5 WHERE id = $1", userID)
}
