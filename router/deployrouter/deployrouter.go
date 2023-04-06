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
		session := sessions.Default(c)

		if session.Get("user") == nil {
			c.Redirect(http.StatusFound, "/auth/sign-in")
			return
		}

		userID := session.Get("user").(int)

		regions, _, _, err := vultr.Region.List(c, &govultr.ListOptions{})
		if err != nil {
			panic(err)
		}

		settings, err := userService.GetSettings(c, userID)
		if err != nil {
			panic(err)
		}

		if c.Request.Method == http.MethodPost {
			region := c.PostForm("region")

			regionFound := false

			for _, o := range regions {
				if o.ID == region {
					regionFound = true
					break
				}
			}

			if settings.ZeroTierNetworkId == nil {
				session.AddFlash("Update your settings first", "danger")
				session.Save()
			} else if !(settings.DeployStatus == 0 || settings.DeployStatus == 4) {
				session.AddFlash("Already deployed, you can't do that", "danger")
				session.Save()
			} else if !regionFound {
				session.AddFlash("Selected region is currently unavailable", "danger")
				session.Save()
			} else {
				if err := userService.StartDeploying(c, userID, region); err != nil {
					panic(err)
				}

				session.AddFlash("Deploying in progress", "success")
				session.Save()
			}

			c.Redirect(http.StatusFound, "/deploy/submit")
			return
		}

		var lastRegion *govultr.Region

		if settings.DeployRegion != nil {
			for i := range regions {
				if regions[i].ID == *settings.DeployRegion {
					lastRegion = &regions[i]
					break
				}
			}
		}

		dangers := session.Flashes("danger")
		successes := session.Flashes("success")

		session.Save()

		c.HTML(http.StatusOK, "deploy/submit", gin.H{
			"Title":    "Deploy",
			"Regions":  regions,
			"Settings": settings,
			"Alert": gin.H{
				"Successes": successes,
				"Dangers":   dangers,
			},
			"LastRegion": lastRegion,
			"Status": gin.H{
				"CanDeploy": settings.ZeroTierNetworkId != nil && (settings.DeployStatus == 0 || settings.DeployStatus == 4),
				"CanStop":   settings.DeployStatus == 3,
				"Failed":    settings.DeployStatus == 4,
				"Verbose":   deployservice.DeployStatusVerbose[settings.DeployStatus],
			},
		})
	}
}

func stop(deployService *deployservice.Service, vultr *govultr.Client, userService *userservice.Service) func(*gin.Context) {
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

		if settings.DeployStatus != 3 {
			session.AddFlash("You can't stop instance right now", "danger")
			session.Save()

			c.Redirect(http.StatusFound, "/deploy/submit")
			return
		}

		if c.Request.Method == http.MethodPost {
			if err := userService.Stop(c, userID); err != nil {
				panic(err)
			}

			session.AddFlash("Your instance will be stopped in a moment", "submit")
			session.Save()

			c.Redirect(http.StatusFound, "/deploy/submit")
			return
		}

		dangers := session.Flashes("danger")
		successes := session.Flashes("success")

		session.Save()

		c.HTML(http.StatusOK, "deploy/stop", gin.H{
			"Title":    "Deploy",
			"Settings": settings,
			"Alert": gin.H{
				"Successes": successes,
				"Dangers":   dangers,
			},
		})
	}
}

func Init(r *gin.RouterGroup, services *service.Services) {
	r.Use(middleware.AuthRequired)
	{
		r.GET("/submit", submit(services.DeployService, services.Vultr, services.UserService))
		r.POST("/submit", submit(services.DeployService, services.Vultr, services.UserService))
		r.GET("/stop", stop(services.DeployService, services.Vultr, services.UserService))
		r.POST("/stop", stop(services.DeployService, services.Vultr, services.UserService))
	}
}
