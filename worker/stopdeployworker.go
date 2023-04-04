package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/ClassAxion/parrot-disco-as-a-service/internal/database"
	"github.com/ClassAxion/parrot-disco-as-a-service/internal/database/user"
)

func getUsersForStop(ctx context.Context, db *database.Database) ([]user.User, error) {
	selector, err := db.QueryContext(ctx, "SELECT id, deployID FROM public.user WHERE deployStatus = 5 OR ((deployStatus = 3 OR deployStatus = 2) AND deployedAt < NOW() - INTERVAL '3' HOUR)")
	if err != nil {
		return nil, err
	}

	defer selector.Close()

	users := []user.User{}

	for selector.Next() {
		if err := selector.Err(); err != nil {
			return nil, err
		}

		var user user.User

		if err := selector.Scan(&user.ID, &user.DeployID); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func stopDeploy(ctx context.Context, wc WorkerContext) error {
	users, err := getUsersForStop(ctx, wc.DB)
	if err != nil {
		return err
	}

	for _, user := range users {
		if err := wc.Services.Vultr.Instance.Delete(ctx, *user.DeployID); err != nil {
			if !strings.Contains(err.Error(), "404") {
				return err
			}

		}

		if _, err := wc.DB.ExecContext(ctx, "UPDATE public.user SET deployStatus = 0, deployedAt = NULL WHERE id = $1", user.ID); err != nil {
			return err
		}
	}

	return nil
}

func StopDeploy() Worker {
	return func(ctx context.Context, wc WorkerContext) error {
		for {
			isOurError := false

			_, err := wc.KeyValue.GetOrSet(ctx, "StopDeploy", time.Now().Add(-(time.Minute * 1)), func() ([]byte, error) {
				if err := stopDeploy(ctx, wc); err != nil {
					isOurError = true
					return nil, fmt.Errorf("failed to stop deploy: %w", err)
				}

				log.Println("StopDeploy")

				return json.Marshal(time.Now())
			})

			if err != nil {
				if isOurError {
					log.Println(err)
					continue
				}

				return fmt.Errorf("failed to get or set key: %w", err)
			}

			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(time.Minute):
			}
		}
	}
}
