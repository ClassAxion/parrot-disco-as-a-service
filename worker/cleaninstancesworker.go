package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/ClassAxion/parrot-disco-as-a-service/internal/database"
	"github.com/ClassAxion/parrot-disco-as-a-service/internal/database/user"
	"github.com/vultr/govultr/v3"
)

func getUsedInstances(ctx context.Context, db *database.Database) ([]user.User, error) {
	selector, err := db.QueryContext(ctx, "SELECT deployID FROM public.user WHERE deployID IS NOT NULL")
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

		if err := selector.Scan(&user.DeployID); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func cleanInstances(ctx context.Context, wc WorkerContext) error {
	users, err := getUsedInstances(ctx, wc.DB)
	if err != nil {
		return err
	}

	exists := map[string]bool{}

	for _, u := range users {
		if u.DeployID != nil {
			exists[*u.DeployID] = true
		}
	}

	instances, _, _, err := wc.Services.Vultr.Instance.List(ctx, &govultr.ListOptions{})
	if err != nil {
		return err
	}

	for _, instance := range instances {
		if _, exists := exists[instance.ID]; !exists {
			wc.Services.Vultr.Instance.Delete(ctx, instance.ID)

			log.Printf("Removing abandoned instance %s", instance.ID)
		}
	}

	return nil
}

func CleanInstances() Worker {
	return func(ctx context.Context, wc WorkerContext) error {
		for {
			isOurError := false

			_, err := wc.KeyValue.GetOrSet(ctx, "CleanInstances", time.Now().Add(-(time.Minute * 1)), func() ([]byte, error) {
				if err := cleanInstances(ctx, wc); err != nil {
					isOurError = true
					return nil, fmt.Errorf("failed to clean instances: %w", err)
				}

				// log.Println("CleanInstances")

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
