package main

import (
	models "app/data"
	"app/middleware"
	"app/routes"
	"log"
	"net/http"
	"slices"

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
	db.AutoMigrate(&models.User{}, &models.UserCredential{}, &models.UserPermissionDef{}, &models.Course{}, &models.Class{}, &models.Attendance{}, &models.StudentInfo{}, &models.TeacherInfo{})

	return *db
}

func seedDatabase(db *gorm.DB) {
	if db.Find(&models.User{}).RowsAffected > 0 {
		return
	}
	users := []models.User{
		{
			FullName: "John Doe",
			Email:    "john.doe@example.com",
			UserCredential: models.UserCredential{
				Username: "john.doe",
				Password: "password",
			},
			UserPermissionDef: models.NewUserPermissionDef("student"),
			StudentInfo: &models.StudentInfo{
				Semester: "1",
			},
		},
		{
			FullName: "Jane Doe",
			Email:    "jane.doe@example.com",
			UserCredential: models.UserCredential{
				Username: "jane.doe",
				Password: "password",
			},
			UserPermissionDef: models.NewUserPermissionDef("teacher"),
			StudentInfo: &models.StudentInfo{
				Semester: "1",
			},
		},
		{
			FullName: "Teacher 1",
			Email:    "teacher1@example.com",
			UserCredential: models.UserCredential{
				Username: "teacher1",
				Password: "password",
			},
			UserPermissionDef: models.NewUserPermissionDef("teacher"),
		},
		{
			FullName: "Teacher 2",
			Email:    "teacher2@example.com",
			UserCredential: models.UserCredential{
				Username: "teacher2",
				Password: "password",
			},
			UserPermissionDef: models.NewUserPermissionDef("teacher"),
		},
		{
			FullName: "Admin",
			Email:    "admin@example.com",
			UserCredential: models.UserCredential{
				Username: "admin",
				Password: "password",
			},
			UserPermissionDef: models.NewUserPermissionDef("manager"),
		},
	}
	courses := []models.Course{
		{
			Name:   "Math",
			Weight: 3,
		},
		{
			Name:   "Physics",
			Weight: 3,
		},
		{
			Name:   "Chemistry",
			Weight: 3,
		},
	}
	classes := []models.Class{
		{
			InstructorID: 3,
			CourseId:     1,
			SlotTotal:    10,
			SlotEquipped: 0,
		},
		{
			InstructorID: 4,
			CourseId:     1,
			SlotTotal:    10,
			SlotEquipped: 0,
		},
	}

	db.Create(&users)
	db.Create(&courses)
	db.Create(&classes)

	idx := slices.IndexFunc(users, func(user models.User) bool {
		return user.FullName == "Teacher 1"
	})
	users[idx].TeacherInfo = &models.TeacherInfo{
		Majors: []models.Course{
			courses[1], courses[2],
		},
	}
	users[idx+1].TeacherInfo = &models.TeacherInfo{
		Majors: []models.Course{
			courses[0], courses[1], courses[2],
		},
	}
	db.Save(&users[idx])
	db.Save(&users[idx+1])
	log.Println("Database seeded")
}

func main() {
	container := dig.New()
	if err := container.Provide(NewDB); err != nil {
		log.Fatal(err)
	}

	engine := gin.Default()
	// the jwt middleware
	container.Invoke(func(db gorm.DB) {
		seedDatabase(&db)
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
