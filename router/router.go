package router

import (
	"fmt"
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
		host := c.GetHeader("Host")
		hash := c.Param("hash")

		if host == "flight.parrotdisco.pl" {
			deployIP, _ := services.DeployService.GetDeployIPByHash(c, hash)
			if deployIP != nil {
				c.Redirect(http.StatusFound, fmt.Sprintf("http://%s:8000/", *deployIP))
			}
		}

		c.AbortWithStatus(http.StatusNotFound)
	})

	return r
}
