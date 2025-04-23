package use_cases

import (
	"app/data"
	models "app/data"

	"gorm.io/gorm"
)

func GetStudentAttendances(db *gorm.DB, studentID uint) []data.Attendance {
	var attendances []models.Attendance
	db.Where("student_id = ?", studentID).Find(&attendances)
	return attendances
}
