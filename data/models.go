package data

import (
	"time"

	"gorm.io/gorm"
)

type Course struct {
	gorm.Model
	Name   string `json:"name"`
	Weight uint   `so tin chi`
}

type Class struct {
	gorm.Model
	InstructorID         uint
	Instructor           User
	CourseId             uint
	Course               Course
	Attendances          []Attendance
	SlotTotal            uint
	SlotEquipped         uint
	OpenForRegisterUntil time.Time
}

type Attendance struct {
	gorm.Model
	StudentID uint
	Student   User
	ClassID   uint
	Grade     float32 `json:"grade"`
}
