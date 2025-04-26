package use_cases

import (
	"app/data"
	models "app/data"
	"app/dto"
	"app/utils"

	"gorm.io/gorm"
)

func GetStudentAttendances(db *gorm.DB, studentID uint) []data.Attendance {
	var attendances []models.Attendance
	db.Preload("Student").Preload("Class").Where("student_id = ?", studentID).Find(&attendances)
	return attendances
}

func GetCourses(db *gorm.DB) []dto.CourseDto {
	var courses []models.Course
	db.Find(&courses)
	return utils.Map(courses, dto.FromCourse)
}

func GetClasses(db *gorm.DB, courseID uint) []dto.ClassDto {
	var classes []models.Class
	db.Preload("Course").Preload("Instructor").Where("course_id = ?", courseID).Find(&classes)
	return utils.Map(classes, dto.FromClass)
}

func JoinClass(db *gorm.DB, classID uint, studentID uint) (*models.Attendance, error) {
	attendance := models.Attendance{
		StudentID: studentID,
		ClassID:   classID,
	}
	tx := db.Create(&attendance)
	if tx.Error != nil {
		return nil, tx.Error
	}
	db.Preload("Student").Preload("Class").Find(&attendance)
	return &attendance, nil
}

func GetClassStudents(db *gorm.DB, classID uint) []data.User {
	var attendances []models.Attendance
	db.Preload("Student").Where("class_id = ?", classID).Find(&attendances)
	return utils.Map(attendances, func(attendance models.Attendance) data.User {
		return *attendance.Student
	})
}
