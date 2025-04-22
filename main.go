package main

import (
	models "app/data"
	"app/middleware"
	viewcontrollers "app/view_controllers"
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
	db.AutoMigrate(&models.Student{}, &models.Manager{}, &models.Instructor{}, &models.Course{}, &models.Attendance{})

	// Seed the database with initial data
	seedDatabase(db)
	return *db
}

func seedDatabase(db *gorm.DB) {
	// Check if database already has records
	var count int64
	db.Model(&models.Student{}).Count(&count)
	if count > 0 {
		return // Database already seeded
	}

	// Create a new teacher
	teacher := models.Instructor{
		Name:     "Professor Smith",
		Username: "teacher1",
		Password: "password123",
	}
	db.Create(&teacher)

	// Create two new courses
	courses := []models.Course{
		{
			Name:         "Introduction to Programming",
			InstructorID: teacher.ID,
			Semester:     "Spring 2025",
		},
		{
			Name:         "Advanced Algorithms",
			InstructorID: teacher.ID,
			Semester:     "Spring 2025",
		},
	}
	for i := range courses {
		db.Create(&courses[i])
	}

	// Create student with username "admin"
	// Create student with username "admin"
	student := models.Student{
		Name:     "Admin User",
		Username: "admin",
		Password: "admin", // Match the JWT middleware credentials
	}
	db.Create(&student)

	// Create attendance records
	attendances := []models.Attendance{
		{
			StudentID: student.ID,
			CourseID:  courses[0].ID,
			Grade:     95.5,
		},
		{
			StudentID: student.ID,
			CourseID:  courses[1].ID,
			Grade:     87.0,
		},
	}
	for i := range attendances {
		db.Create(&attendances[i])
	}

	log.Println("Database seeded successfully!")
}

func main() {
	container := dig.New()
	if err := container.Provide(NewDB); err != nil {
		log.Fatal(err)
	}

	engine := gin.Default()
	// the jwt middleware
	authMiddleware, err := jwt.New(middleware.InitParams())
	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	// register middleware
	engine.Use(handlerMiddleware(authMiddleware))

	// register route
	registerRoute(container, engine, authMiddleware)

	// start http server
	if err = http.ListenAndServe(":8080", engine); err != nil {
		log.Fatal(err)
	}
}

func registerRoute(container *dig.Container, r *gin.Engine, handle *jwt.GinJWTMiddleware) {
	r.POST("/login", handle.LoginHandler)
	r.NoRoute(handle.MiddlewareFunc(), handleNoRoute())
	r.LoadHTMLGlob("html/*")

	container.Invoke(func(db gorm.DB) {
		r.GET(
			"/classes",
			handle.MiddlewareFunc(),
			func(c *gin.Context) {
				viewcontrollers.GetStudentClasses(c.Writer, c.Request, handle.IdentityHandler(c).(*middleware.User), &db)
			},
		)
		r.GET("/login", func(c *gin.Context) {
			// Render the HTML file
			c.HTML(200, "login.html", nil)
		})
		r.GET("/", func(c *gin.Context) {
			// Render the HTML file
			c.HTML(200, "index.html", nil)
		})
	})

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
