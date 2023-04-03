package deployrouter

import (
	"net/http"

	"github.com/ClassAxion/parrot-disco-as-a-service/internal/middleware"
	"github.com/ClassAxion/parrot-disco-as-a-service/service"
	"github.com/ClassAxion/parrot-disco-as-a-service/service/deployservice"
	"github.com/ClassAxion/parrot-disco-as-a-service/service/userservice"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/vultr/govultr/v3"
)

func submit(deployService *deployservice.Service, vultr *govultr.Client, userService *userservice.Service) func(*gin.Context) {
	return func(c *gin.Context) {
		regions, _, _, err := vultr.Region.List(c, &govultr.ListOptions{})
		if err != nil {
			panic(err)
		}

		session := sessions.Default(c)

		userID := session.Get("user").(int)

		settings, err := userService.GetSettings(c, userID)
		if err != nil {
			panic(err)
		}

		c.HTML(http.StatusOK, "deploy/submit", gin.H{
			"Title":    "Deploy",
			"Regions":  regions,
			"Settings": settings,
		})
	}
}

func Init(r *gin.RouterGroup, services *service.Services) {
	r.Use(middleware.AuthRequired)
	{
		r.GET("/submit", submit(services.DeployService, services.Vultr, services.UserService))
	}
}
