package viewcontrollers

import (
	models "app/data"
	"app/middleware"
	"html/template"
	"net/http"

	"gorm.io/gorm"
)

type studentClassesData struct {
	Attendances []models.Attendance
}

func GetStudentClasses(w http.ResponseWriter, r *http.Request, user *middleware.User, db *gorm.DB) {
	tmpl := template.Must(template.ParseFiles("html/student_attendances.html"))
	var attendances []models.Attendance
	var student models.Student
	db.Where("username = ?", user.UserName).First(&student)
	db.Where("student_id = ?", student.ID).Find(&attendances)
	tmpl.Execute(w, studentClassesData{
		Attendances: attendances,
	})
}
