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

	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})

	RequireSignedIn := handle.MiddlewareFunc()

	api := r.Group("/api")
	container.Invoke(func(db gorm.DB) {
		api.POST("/login", handle.LoginHandler)
		api.POST("/signup", viewcontrollers.HandlePostSignup(&db))
		student_group := api.Group("/student", RequireSignedIn)
		student_group.GET("/classes/attended", viewcontrollers.GetAttendedClasses(&db))
		student_group.GET("/classes/courses", viewcontrollers.GetCourses(&db))
		student_group.GET("/classes/attend/:course_id", viewcontrollers.GetAttendableClasses(&db))
	})

	auth := r.Group("/auth", RequireSignedIn)
	auth.GET("/refresh_token", handle.RefreshHandler)
}

func handleNoRoute() func(c *gin.Context) {
	return func(c *gin.Context) {
		claims := jwt.ExtractClaims(c)
		log.Printf("NoRoute claims: %#v\n", claims)
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	}
}
