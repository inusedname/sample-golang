package main

import (
	models "app/data"
	"app/middleware"
	"log"
	"net/http"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&models.Student{}, &models.Manager{}, &models.Instructor{}, &models.Course{}, &models.Attendance{})
	db.Create(&models.Student{
		Name:     "Nguyen Van A",
		Username: "a",
		Password: "a",
	})
	engine := gin.Default()
	// the jwt middleware
	authMiddleware, err := jwt.New(middleware.InitParams())
	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	// register middleware
	engine.Use(handlerMiddleware(authMiddleware))

	// register route
	registerRoute(engine, authMiddleware)

	// start http server
	if err = http.ListenAndServe(":8080", engine); err != nil {
		log.Fatal(err)
	}
}

func registerRoute(r *gin.Engine, handle *jwt.GinJWTMiddleware) {
	r.POST("/login", handle.LoginHandler)
	r.NoRoute(handle.MiddlewareFunc(), handleNoRoute())

	auth := r.Group("/auth", handle.MiddlewareFunc())
	auth.GET("/refresh_token", handle.RefreshHandler)
	auth.GET("/hello", func(c *gin.Context) {
		helloHandler(c, handle)
	})
}

func handlerMiddleware(authMiddleware *jwt.GinJWTMiddleware) gin.HandlerFunc {
	return func(context *gin.Context) {
		errInit := authMiddleware.MiddlewareInit()
		if errInit != nil {
			log.Fatal("authMiddleware.MiddlewareInit() Error:" + errInit.Error())
		}
	}
}

func handleNoRoute() func(c *gin.Context) {
	return func(c *gin.Context) {
		claims := jwt.ExtractClaims(c)
		log.Printf("NoRoute claims: %#v\n", claims)
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	}
}

func helloHandler(c *gin.Context, authMiddleware *jwt.GinJWTMiddleware) {
	claims := jwt.ExtractClaims(c)
	user, _ := authMiddleware.IdentityHandler(c).(*middleware.User)
	c.JSON(200, gin.H{
		"userID":   claims[middleware.IdentityKey],
		"userName": user.UserName,
		"text":     "Hello World.",
	})
}
