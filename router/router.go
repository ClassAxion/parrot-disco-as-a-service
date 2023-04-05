package router

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ClassAxion/parrot-disco-as-a-service/router/authrouter"
	"github.com/ClassAxion/parrot-disco-as-a-service/router/dashboardrouter"
	"github.com/ClassAxion/parrot-disco-as-a-service/router/deployrouter"
	"github.com/ClassAxion/parrot-disco-as-a-service/router/userrouter"
	"github.com/ClassAxion/parrot-disco-as-a-service/service"
	"github.com/gin-gonic/gin"
)

func New(
	r *gin.Engine,
	services *service.Services,
) *gin.Engine {
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/dashboard")
	})

	dashboardrouter.Init(r.Group("/dashboard"), services)
	authrouter.Init(r.Group("/auth"), services)
	userrouter.Init(r.Group("/user"), services)
	deployrouter.Init(r.Group("/deploy"), services)

	r.StaticFS("/public", http.Dir("./public/"))

	r.GET("/:hash", func(c *gin.Context) {
		host := c.Request.Host
		hash := c.Param("hash")

		log.Println(host, hash)

		if host == "flight.parrotdisco.pl" {
			user, _ := services.DeployService.GetDeployIPByHash(c, hash)

			if user != nil {
				if user.DeployStatus == 3 && user.DeployIP != nil {
					c.Redirect(http.StatusFound, fmt.Sprintf("http://%s:8000/", *user.DeployIP))
				} else if user.DeployStatus == 1 || user.DeployStatus == 2 {
					c.Data(http.StatusNotFound, "plain/text; charset=utf-8", []byte("deploying in progress"))
				} else if user.DeployStatus == 4 {
					c.Data(http.StatusNotFound, "plain/text; charset=utf-8", []byte("deploying failed"))
				} else {
					c.Data(http.StatusNotFound, "plain/text; charset=utf-8", []byte("not deployed"))
				}

				return
			}
		}

		c.AbortWithStatus(http.StatusNotFound)
	})

	return r
}
