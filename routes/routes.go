package routes

import (
	viewcontrollers "app/view_controllers"
	"log"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"go.uber.org/dig"
	"gorm.io/gorm"
)

func RegisterRoute(container *dig.Container, r *gin.Engine, handle *jwt.GinJWTMiddleware) {
	r.LoadHTMLGlob("html/*")
	r.NoRoute(handle.MiddlewareFunc(), handleNoRoute())

	r.POST("/login", handle.LoginHandler)
	r.GET("/login", func(c *gin.Context) {
		c.HTML(200, "login.html", nil)
	})
	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})

	container.Invoke(func(db gorm.DB) {
		r.GET("/signup", viewcontrollers.HandleGetSignup)
		r.POST("/signup", viewcontrollers.HandlePostSignup(&db))
		r.GET(
			"/classes",
			handle.MiddlewareFunc(),
			viewcontrollers.GetStudentClasses(&db),
		)
	})

	auth := r.Group("/auth", handle.MiddlewareFunc())
	auth.GET("/refresh_token", handle.RefreshHandler)
}

func handleNoRoute() func(c *gin.Context) {
	return func(c *gin.Context) {
		claims := jwt.ExtractClaims(c)
		log.Printf("NoRoute claims: %#v\n", claims)
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	}
}
