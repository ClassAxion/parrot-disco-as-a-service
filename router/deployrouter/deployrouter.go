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

		if c.Request.Method == http.MethodPost {
			region := c.PostForm("region")
			rememberRegion := c.PostForm("rememberRegion") != ""

			regionFound := false

			for _, o := range regions {
				if o.ID == region {
					regionFound = true
					break
				}
			}

			if settings.DeployStatus == 3 {
				session.AddFlash("Already deployed, you can't do that", "danger")
				session.Save()
			} else if !regionFound {
				session.AddFlash("Selected region is currently unavailable", "danger")
				session.Save()
			} else {
				if rememberRegion {
					if err := userService.UpdateDefaultRegion(c, userID, region); err != nil {
						panic(err)
					}
				}

				if err := userService.StartDeploying(c, userID); err != nil {
					panic(err)
				}

				session.AddFlash("Deploying in progress", "success")
				session.Save()
			}

			c.Redirect(http.StatusFound, "/deploy/submit")
			return
		}

		var defaultRegion *govultr.Region

		if settings.DefaultRegion != nil {
			for i := range regions {
				if regions[i].ID == *settings.DefaultRegion {
					defaultRegion = &regions[i]
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
			"DefaultRegion": defaultRegion,
			"Status": gin.H{
				"CanDeploy": settings.DeployStatus == 0 || settings.DeployStatus == 4,
				"CanStop":   settings.DeployStatus == 3,
				"Verbose":   deployservice.DeployStatusVerbose[settings.DeployStatus],
			},
		})
	}
}

func stop(deployService *deployservice.Service, vultr *govultr.Client, userService *userservice.Service) func(*gin.Context) {
	return func(c *gin.Context) {
		session := sessions.Default(c)

		userID := session.Get("user").(int)

		settings, err := userService.GetSettings(c, userID)
		if err != nil {
			panic(err)
		}

		if settings.DeployStatus != 3 {
			session.AddFlash("You can't to that", "danger")
			session.Save()

			c.Redirect(http.StatusFound, "/deploy/submit")
			return
		}

		if c.Request.Method == http.MethodPost {
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
