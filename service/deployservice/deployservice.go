package deployservice

import (
	"context"

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

var DeployStatusVerbose = map[int]string{
	0: "ready to deploy",
	1: "received deploy request",
	2: "deploying in progress",
	3: "deployed",
	4: "deploy failed",
	5: "stopping",
}

func (service *Service) GetDeployIPByHash(ctx context.Context, hash string) (*user.User, error) {
	row := service.DB.QueryRowContext(ctx, "SELECT deployIP, deployStatus FROM public.user WHERE hash = $1", hash)
	if err := row.Err(); err != nil {
		return nil, err
	}

	var user user.User

	if err := row.Scan(&user.DeployIP, &user.DeployStatus); err != nil {
		return nil, err
	}

	return &user, nil
}
