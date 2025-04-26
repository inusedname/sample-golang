package main

import (
	models "app/data"
	"app/middleware"
	"app/routes"
	"log"
	"net/http"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"go.uber.org/dig"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewDB() gorm.DB {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&models.User{}, &models.UserCredential{}, &models.UserPermissionDef{}, &models.Course{}, &models.Class{}, &models.Attendance{})

	return *db
}

func main() {
	container := dig.New()
	if err := container.Provide(NewDB); err != nil {
		log.Fatal(err)
	}

	engine := gin.Default()
	// the jwt middleware
	container.Invoke(func(db gorm.DB) {
		authMiddleware, err := jwt.New(middleware.InitParams(&db))
		if err != nil {
			log.Fatal("JWT Error:" + err.Error())
		}

		// register middleware
		engine.Use(handlerMiddleware(authMiddleware))

		// register route
		routes.RegisterRoute(container, engine, authMiddleware)

		// start http server
		if err = http.ListenAndServe(":8080", engine); err != nil {
			log.Fatal(err)
		}
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
