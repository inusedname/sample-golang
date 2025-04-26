package use_cases

import (
	"app/data"
	models "app/data"

	"gorm.io/gorm"
)

func GetStudentAttendances(db *gorm.DB, studentID uint) []data.Attendance {
	var attendances []models.Attendance
	db.Preload("Student").Where("student_id = ?", studentID).Find(&attendances)
	return attendances
}

func GetCourses(db *gorm.DB) []models.Course {
	var courses []models.Course
	db.Find(&courses)
	return courses
}

func GetClasses(db *gorm.DB, courseID uint) []models.Class {
	var classes []models.Class
	db.Where("course_id = ?", courseID).Find(&classes)
	return classes
}

func JoinClass(db *gorm.DB, classID uint, studentID uint) models.Attendance {
	attendance := models.Attendance{
		StudentID: studentID,
		ClassID:   classID,
	}
	db.Create(&attendance)
	return attendance
}
