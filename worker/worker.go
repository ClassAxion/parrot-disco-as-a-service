package worker

import (
	"context"

	"github.com/ClassAxion/parrot-disco-as-a-service/internal/config"
	"github.com/ClassAxion/parrot-disco-as-a-service/internal/database"
	"github.com/ClassAxion/parrot-disco-as-a-service/service"
	"github.com/mlvzk/gopgs/pgkv"
)

type WorkerContext struct {
	Config   *config.Config
	DB       *database.Database
	Services *service.Services
	KeyValue *pgkv.Store
}

type Worker func(context.Context, WorkerContext) error
