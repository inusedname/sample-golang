package data

import (
	"gorm.io/gorm"
)

type Course struct {
	gorm.Model
	Name         string `json:"name"`
	InstructorID uint
	Instructor   User
	Attendances  []Attendance
}

type Attendance struct {
	gorm.Model
	StudentID uint
	CourseID  uint
	Student   User
	Course    Course
	Grade     float32 `json:"grade"`
}
