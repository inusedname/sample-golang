package viewcontrollers

import (
	models "app/data"
	"app/middleware"
	"app/use_cases"
	"app/utils/strconv"
	"html/template"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type studentClassesData struct {
	Attendances []models.Attendance
}

func GetStudentClasses(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		tmpl := template.Must(template.ParseFiles("html/student_attendances.html"))
		var user = models.GetUser(db, c)
		if !user.UserPermissionDef.ViewAttendance {
			c.HTML(403, "unauthorized.html", gin.H{
				"Reason": "You don't have permission to view classes",
			})
			return
		}
		attendances := use_cases.GetStudentAttendances(db, user.ID)
		tmpl.Execute(c.Writer, studentClassesData{
			Attendances: attendances,
		})
	}
}

func GetAttendedClasses(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		user := models.GetUser(db, c)
		if !user.UserPermissionDef.ViewAttendance {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "You don't have permission to view classes"})
			return
		}
		userID := c.MustGet("id").(middleware.MyClaims).UserID
		log.Printf("User ID: %v", userID)
		classes := use_cases.GetStudentAttendances(db, userID)
		c.JSON(http.StatusOK, gin.H{
			"classes": classes,
		})
	}
}

func GetCourses(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		courses := use_cases.GetCourses(db)
		c.JSON(http.StatusOK, gin.H{
			"courses": courses,
		})
	}
}

func GetAttendableClasses(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		courseID, err := strconv.Atoui(c.Param("course_id"))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		classes := use_cases.GetClasses(db, courseID)
		c.JSON(http.StatusOK, gin.H{
			"classes": classes,
		})
	}
}
