package viewcontrollers

import (
	models "app/data"
	"app/middleware"
	"errors"
	"html/template"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type studentClassesData struct {
	Attendances []models.Attendance
}

func GetStudentClasses(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		tmpl := template.Must(template.ParseFiles("html/student_attendances.html"))
		var attendances []models.Attendance
		var userId = jwt.ExtractClaims(c)[middleware.IdentityKey]
		var user models.User
		db.Preload("UserPermissionDef").Where("id = ?", userId).First(&user)
		if !user.UserPermissionDef.ViewAttendance {
			c.AbortWithError(403, errors.New("you don't have permission to view attendances"))
			return
		}
		db.Where("student_id = ?", userId).Find(&attendances)
		tmpl.Execute(c.Writer, studentClassesData{
			Attendances: attendances,
		})
	}
}
