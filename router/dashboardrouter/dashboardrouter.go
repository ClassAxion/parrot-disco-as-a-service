package dashboardrouter

import (
	"net/http"

	"github.com/ClassAxion/parrot-disco-as-a-service/internal/middleware"
	"github.com/ClassAxion/parrot-disco-as-a-service/service"
	"github.com/ClassAxion/parrot-disco-as-a-service/service/dashboardservice"
	"github.com/ClassAxion/parrot-disco-as-a-service/service/deployservice"
	"github.com/ClassAxion/parrot-disco-as-a-service/service/userservice"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func dashboard(dashboardService *dashboardservice.Service, userService *userservice.Service) func(*gin.Context) {
	return func(c *gin.Context) {
		session := sessions.Default(c)

		if session.Get("user") == nil {
			c.Redirect(http.StatusFound, "/auth/sign-in")
			return
		}

		userID := session.Get("user").(int)

		settings, err := userService.GetSettings(c, userID)
		if err != nil {
			panic(err)
		}

		c.HTML(http.StatusOK, "dashboard/index", gin.H{
			"Title":    "Homepage",
			"Settings": settings,
			"Status": gin.H{
				"CanDeploy": settings.DeployStatus == 0 || settings.DeployStatus == 4,
				"CanStop":   settings.DeployStatus == 3,
				"Verbose":   deployservice.DeployStatusVerbose[settings.DeployStatus],
			},
		})
	}
}

func Init(r *gin.RouterGroup, services *service.Services) {
	r.Use(middleware.AuthRequired)
	{
		r.GET("/", dashboard(services.DashboardService, services.UserService))
	}
}
