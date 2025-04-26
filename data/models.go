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
	ID        uint           `gorm:"primaryKey"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	StudentID uint           `gorm:"unique" json:"-"`
	Student   *User
	ClassID   uint `gorm:"unique" json:"-"`
	Class     *Class
	Grade     float32 `json:"grade"`
}
