package viewcontrollers

import (
	models "app/data"
	"app/middleware"
	"html/template"

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
		var user = c.MustGet(middleware.IdentityKey).(*models.User)
		if !user.UserPermissionDef.ViewAttendance {
			c.HTML(403, "unauthorized.html", gin.H{
				"Reason": "You don't have permission to view classes",
			})
			return
		}
		db.Where("student_id = ?", user.ID).Find(&attendances)
		tmpl.Execute(c.Writer, studentClassesData{
			Attendances: attendances,
		})
	}
}
