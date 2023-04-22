package user

import "time"

type Location struct {
	Latitude  float64
	Longitude float64
	Altitude  int
}

type User struct {
	ID                int
	Email             string
	Password          string
	Name              string
	Hash              *string
	ZeroTierNetworkId *string
	ZeroTierDiscoIP   *string
	HomeLocation      *Location
	DeployStatus      int
	// 0 - stopped (can deploy)
	// 1 - deploy request
	// 2 - deploying
	// 3 - deployed (can stop)
	// 4 - deploy failed (can deploy)
	// 5 - stop request
	DeployIP       *string
	DeployPassword *string
	DeployID       *string
	DeployRegion   *string
	DeployedAt     *time.Time
	ShareLocation  bool
}
