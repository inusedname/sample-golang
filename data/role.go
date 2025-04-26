package data

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID                  uint           `gorm:"primarykey"`
	CreatedAt           time.Time      `json:"-"`
	UpdatedAt           time.Time      `json:"-"`
	DeletedAt           gorm.DeletedAt `gorm:"index" json:"-"`
	FullName            string
	Email               string            `gorm:"unique"`
	UserPermissionDefID uint              `json:"-"`
	UserPermissionDef   UserPermissionDef `json:"-"`
	UserCredentialID    uint              `json:"-"`
	UserCredential      UserCredential    `json:"-"`
	StudentInfoID       *uint             `json:"-"`
	StudentInfo         *StudentInfo      `json:"omitempty"`
	TeacherInfoID       *uint             `json:"-"`
	TeacherInfo         *TeacherInfo      `json:"omitempty"`
}

type StudentInfo struct {
	gorm.Model
	Semester string
}

type TeacherInfo struct {
	gorm.Model
	Majors []Course `gorm:"many2many:teacher_majors;"`
}

type UserCredential struct {
	gorm.Model
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserPermissionDef struct {
	gorm.Model
	ReadCourse       bool
	CreateCourse     bool
	CreateAttendance bool
	ViewAttendance   bool
	UpdateCourse     bool
}

func NewUserPermissionDef(role string) UserPermissionDef {
	switch role {
	case "student":
		return UserPermissionDef{
			ViewAttendance:   true,
			CreateAttendance: false,
		}
	case "teacher":
		return UserPermissionDef{
			ReadCourse: true,
		}
	case "manager":
		return UserPermissionDef{
			CreateCourse: true,
		}
	default:
		return UserPermissionDef{}
	}
}
