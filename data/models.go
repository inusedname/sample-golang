package data

import (
	"gorm.io/gorm"
)

type Student struct {
	gorm.Model
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Manager struct {
	gorm.Model
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Instructor struct {
	gorm.Model
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Course struct {
	gorm.Model
	Name         string `json:"name"`
	InstructorID uint
	Instructor   Instructor
	Semester     string `json:"semester"`
}

type Attendance struct {
	gorm.Model
	StudentID uint
	CourseID  uint
	Student   Student
	Course    Course
	Grade     float32 `json:"grade"`
}
