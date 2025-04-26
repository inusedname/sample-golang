package viewcontrollers

import (
	models "app/data"
	"app/middleware"
	"app/use_cases"
	"app/utils/strconv"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type studentClassesData struct {
	Attendances []models.Attendance
}

func GetTimetable(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		user := models.GetUser(db, c)
		if !user.UserPermissionDef.ViewAttendance {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "You don't have permission to view classes"})
			return
		}
		userID := c.MustGet("id").(middleware.MyClaims).UserID
		log.Printf("User ID: %v", userID)
		attendances := use_cases.GetStudentAttendances(db, userID)
		c.JSON(http.StatusOK, gin.H{
			"attendances": attendances,
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
		classes := use_cases.GetOpenClasses(db, courseID)
		c.JSON(http.StatusOK, gin.H{
			"classes": classes,
		})
	}
}

func JoinClass(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		studentId := c.MustGet("id").(middleware.MyClaims).UserID
		classID, err := strconv.Atoui(c.Param("course_id"))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		log.Printf("Student ID: %v", studentId)
		attendance, error := use_cases.JoinClass(db, classID, studentId)
		if error != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": error.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"attendance": attendance,
		})
	}
}

type ClassInfo struct {
	Students []models.User
}

func GetClassInfo(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		classID, err := strconv.Atoui(c.Param("class_id"))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		students := use_cases.GetClassStudents(db, classID)
		c.JSON(http.StatusOK, gin.H{
			"students": students,
		})
	}
}
