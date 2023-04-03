package service

import (
	"github.com/ClassAxion/parrot-disco-as-a-service/service/authservice"
	"github.com/ClassAxion/parrot-disco-as-a-service/service/dashboardservice"
	"github.com/ClassAxion/parrot-disco-as-a-service/service/deployservice"
	"github.com/ClassAxion/parrot-disco-as-a-service/service/userservice"
	"github.com/vultr/govultr/v3"
)

type Services struct {
	DashboardService *dashboardservice.Service
	AuthService      *authservice.Service
	UserService      *userservice.Service
	DeployService    *deployservice.Service
	Vultr            *govultr.Client
}
