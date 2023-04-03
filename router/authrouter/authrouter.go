package authrouter

import (
	"net/http"

	"github.com/ClassAxion/parrot-disco-as-a-service/service"
	"github.com/ClassAxion/parrot-disco-as-a-service/service/authservice"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func signin(authService *authservice.Service) func(*gin.Context) {
	return func(c *gin.Context) {
		session := sessions.Default(c)

		if c.Request.Method == http.MethodPost {
			email := c.PostForm("email")
			password := c.PostForm("password")

			if email == "" || password == "" {
				session.AddFlash("Please fill all fields correctly", "danger")
				session.Save()
			} else {
				user, err := authService.Login(c, email, password)
				if err != nil {
					session.AddFlash("Please fill all fields correctly", "danger")
					session.Save()
				} else {
					session.Set("user", user.ID)
					session.Save()

					c.Redirect(http.StatusFound, "/dashboard")
					return
				}
			}

			c.Redirect(http.StatusFound, "/auth/sign-in")
			return
		}

		dangers := session.Flashes("danger")
		successes := session.Flashes("success")

		session.Save()

		c.HTML(http.StatusOK, "auth/signin", gin.H{
			"Title": "Sign in",
			"Alert": gin.H{
				"Successes": successes,
				"Dangers":   dangers,
			},
		})
	}
}

func signup(authService *authservice.Service) func(*gin.Context) {
	return func(c *gin.Context) {
		session := sessions.Default(c)

		if c.Request.Method == http.MethodPost {
			name := c.PostForm("name")
			email := c.PostForm("email")
			password := c.PostForm("password")

			if name == "" || email == "" || password == "" {
				session.AddFlash("Please fill all fields correctly", "danger")
				session.Save()
			} else if len(password) < 8 || len(password) > 50 {
				session.AddFlash("Please fill all fields correctly", "danger")
				session.Save()
			} else if err := authService.Register(c, name, email, password); err != nil {
				session.AddFlash("An unknown error occured", "danger")
				session.Save()
			} else {
				session.AddFlash("Account created - you can now log in", "success")
				session.Save()

				c.Redirect(http.StatusFound, "/auth/sign-in")
				return
			}

			c.Redirect(http.StatusFound, "/auth/sign-up")
			return
		}

		dangers := session.Flashes("danger")
		successes := session.Flashes("success")

		session.Save()

		c.HTML(http.StatusOK, "auth/signup", gin.H{
			"Title": "Sign up",
			"Alert": gin.H{
				"Successes": successes,
				"Dangers":   dangers,
			},
		})
	}
}

func signout(authService *authservice.Service) func(*gin.Context) {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		session.Delete("user")

		session.Save()

		c.Redirect(http.StatusFound, "/auth/sign-in")
	}
}

func Init(r *gin.RouterGroup, services *service.Services) {
	r.GET("/sign-in", signin(services.AuthService))
	r.POST("/sign-in", signin(services.AuthService))
	r.GET("/sign-out", signout(services.AuthService))
	r.GET("/sign-up", signup(services.AuthService))
	r.POST("/sign-up", signup(services.AuthService))
}
