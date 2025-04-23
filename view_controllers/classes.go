package viewcontrollers

import (
	models "app/data"
	"app/use_cases"
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
