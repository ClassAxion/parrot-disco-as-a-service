package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/ClassAxion/parrot-disco-as-a-service/internal/database"
	"github.com/ClassAxion/parrot-disco-as-a-service/internal/database/user"
	"github.com/helloyi/go-sshclient"
	"github.com/vultr/govultr/v3"
)

func getUsersForDeploy(ctx context.Context, db *database.Database) ([]user.User, error) {
	selector, err := db.QueryContext(ctx, "SELECT id, email, deployRegion, hash, zeroTierNetworkId, zeroTierDiscoIP, homeLocation FROM public.user WHERE deployStatus = 1")
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

		var location *[]byte

		if err := selector.Scan(&user.ID, &user.Email, &user.DeployRegion, &user.Hash, &user.ZeroTierNetworkId, &user.ZeroTierDiscoIP, &location); err != nil {
			return nil, err
		}

		if location != nil {
			if err := json.Unmarshal(*location, &user.HomeLocation); err != nil {
				return nil, err
			}
		}

		users = append(users, user)
	}

	return users, nil
}

func startDeploy(ctx context.Context, wc WorkerContext) error {
	users, err := getUsersForDeploy(ctx, wc.DB)
	if err != nil {
		return err
	}

	for _, user := range users {
		log.Printf("Creating instance for user %s", user.Email)

		if _, err := wc.DB.ExecContext(ctx, "UPDATE public.user SET deployStatus = 2 WHERE id = $1", user.ID); err != nil {
			return err
		}

		instance, _, err := wc.Services.Vultr.Instance.Create(ctx, &govultr.InstanceCreateReq{
			Label:   fmt.Sprintf("flight-%s-%d-%s", *user.Hash, user.ID, user.Email),
			Backups: "disabled",
			OsID:    1743,
			// Plan:    "vc2-1c-1gb",
			Plan:   "vc2-2c-4gb",
			Region: *user.DeployRegion,
		})
		if err != nil {
			return err
		}

		if _, err := wc.DB.ExecContext(ctx, "UPDATE public.user SET deployID = $2, deployPassword = $3, deployedAt = NOW() WHERE id = $1", user.ID, instance.ID, instance.DefaultPassword); err != nil {
			return err
		}

		var IP string

		for {
			res, _, err := wc.Services.Vultr.Instance.Get(ctx, instance.ID)
			if err != nil {
				continue
			}

			if res.Status == "active" && (res.ServerStatus == "ok" || res.ServerStatus == "installingbooting") {
				IP = res.MainIP

				break
			}

			time.Sleep(time.Second * 30)
		}

		log.Println("Instance created, connecting..")

		if _, err := wc.DB.ExecContext(ctx, "UPDATE public.user SET deployIP = $2 WHERE id = $1", user.ID, IP); err != nil {
			return err
		}

		var client *sshclient.Client

		for {
			client, err = sshclient.DialWithPasswd(fmt.Sprintf("%s:22", IP), "root", instance.DefaultPassword)
			if err == nil {
				break
			}

			time.Sleep(time.Second * 15)
		}

		log.Println("Connected! Deploying app..")

		if _, err := client.Cmd("curl -fsSL https://get.docker.com -o get-docker.sh && sh get-docker.sh").Output(); err != nil {
			return fmt.Errorf("failed to install docker")
		}

		if _, err := client.Cmd("git clone https://github.com/ClassAxion/parrot-disco-devops").Output(); err != nil {
			return fmt.Errorf("failed to clone git repo")
		}

		if _, err := client.Cmd("touch ~/parrot-disco-devops/.env").Output(); err != nil {
			return fmt.Errorf("failed to create env")
		}

		if _, err := client.Cmd(fmt.Sprintf("echo DISCO_IP=%s >> ~/parrot-disco-devops/.env", *user.ZeroTierDiscoIP)).Output(); err != nil {
			return fmt.Errorf("failed to set disco IP")
		}

		if _, err := client.Cmd(fmt.Sprintf("echo ZEROTIER_ID=%s >> ~/parrot-disco-devops/.env", *user.ZeroTierNetworkId)).Output(); err != nil {
			return fmt.Errorf("failed to set zerotier ID")
		}

		if user.HomeLocation != nil {
			if _, err := client.Cmd(fmt.Sprintf("echo HOME_LOCATION=%.5f,%.5f,%d >> ~/parrot-disco-devops/.env", user.HomeLocation.Latitude, user.HomeLocation.Longitude, user.HomeLocation.Altitude)).Output(); err != nil {
				return fmt.Errorf("failed to set home location")
			}
		}

		if _, err := client.Cmd("cd ~/parrot-disco-devops && docker compose up -d").Output(); err != nil {
			return fmt.Errorf("failed to start containers")
		}

		if _, err := client.Cmd("ufw disable").Output(); err != nil {
			return fmt.Errorf("failed to disable ufw")
		}

		defer client.Close()

		if _, err := wc.DB.ExecContext(ctx, "UPDATE public.user SET deployStatus = 3 WHERE id = $1", user.ID); err != nil {
			return err
		}

		log.Println("Deployed successfully")
	}

	return nil
}

func StartDeploy() Worker {
	return func(ctx context.Context, wc WorkerContext) error {
		for {
			isOurError := false

			_, err := wc.KeyValue.GetOrSet(ctx, "StartDeploy", time.Now().Add(-(time.Minute * 1)), func() ([]byte, error) {
				if err := startDeploy(ctx, wc); err != nil {
					isOurError = true
					return nil, fmt.Errorf("failed to start deploy: %w", err)
				}

				// log.Println("StartDeploy")

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
