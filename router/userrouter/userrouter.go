package userrouter

import (
	"net/http"
	"strconv"

	"github.com/ClassAxion/parrot-disco-as-a-service/internal/database/user"
	"github.com/ClassAxion/parrot-disco-as-a-service/internal/middleware"
	"github.com/ClassAxion/parrot-disco-as-a-service/service"
	"github.com/ClassAxion/parrot-disco-as-a-service/service/userservice"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func settings(userService *userservice.Service) func(*gin.Context) {
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

		if c.Request.Method == http.MethodPost {
			hash := c.PostForm("hash")
			zeroTierNetworkId := c.PostForm("zeroTierNetworkId")
			zeroTierDiscoIP := c.PostForm("zeroTierDiscoIP")

			if !(settings.DeployStatus == 0 || settings.DeployStatus == 4) {
				session.AddFlash("You cannot modify settings now", "danger")
				session.Save()
			} else if hash == "" || zeroTierNetworkId == "" || zeroTierDiscoIP == "" {
				session.AddFlash("Please fill all fields correctly", "danger")
				session.Save()
			} else {
				var homeLocation *user.Location

				homeLocationLatitudeRaw := c.PostForm("homeLocationLatitude")
				homeLocationLongitudeRaw := c.PostForm("homeLocationLongitude")
				homeLocationAltitudeRaw := c.PostForm("homeLocationAltitude")

				useHomeLocation := homeLocationLatitudeRaw != "" || homeLocationLongitudeRaw != "" || homeLocationAltitudeRaw != ""

				if useHomeLocation && (homeLocationLatitudeRaw == "" || homeLocationLongitudeRaw == "" || homeLocationAltitudeRaw == "") {
					session.AddFlash("Please fill home location fields correctly", "danger")
					session.Save()
				} else {
					if useHomeLocation {
						homeLocation = &user.Location{}

						if val, err := strconv.ParseFloat(homeLocationLatitudeRaw, 64); err == nil {
							homeLocation.Latitude = val
						}

						if val, err := strconv.ParseFloat(homeLocationLongitudeRaw, 64); err == nil {
							homeLocation.Longitude = val
						}

						if val, err := strconv.Atoi(homeLocationAltitudeRaw); err == nil {
							homeLocation.Altitude = val
						}
					}

					if useHomeLocation && (homeLocation.Latitude == 0 || homeLocation.Longitude == 0 || homeLocation.Altitude == 0) {
						session.AddFlash("Please fill home location fields correctly", "danger")
						session.Save()
					} else if err := userService.SaveSettings(c, userID, hash, zeroTierNetworkId, zeroTierDiscoIP, homeLocation); err != nil {
						session.AddFlash("Please fill all fields correctly", "danger")
						session.Save()
					} else {
						session.AddFlash("Settings saved", "success")
						session.Save()
					}
				}
			}

			c.Redirect(http.StatusFound, "/user/settings")
			return
		}

		dangers := session.Flashes("danger")
		successes := session.Flashes("success")

		session.Save()

		c.HTML(http.StatusOK, "user/settings", gin.H{
			"Title": "Settings",
			"User":  settings,
			"Alert": gin.H{
				"Successes": successes,
				"Dangers":   dangers,
			},
			"CanChange": settings.DeployStatus == 0 || settings.DeployStatus == 4,
		})
	}
}

func Init(r *gin.RouterGroup, services *service.Services) {
	r.Use(middleware.AuthRequired)
	{
		r.GET("/settings", settings(services.UserService))
		r.POST("/settings", settings(services.UserService))
	}
}
