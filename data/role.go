package data

import "gorm.io/gorm"

type User struct {
	gorm.Model
	FullName            string
	Email               string `gorm:"primaryKey"`
	UserPermissionDefID uint
	UserPermissionDef   UserPermissionDef
	UserCredentialID    uint
	UserCredential      UserCredential
	StudentInfoID       uint
	StudentInfo         StudentInfo
}

type StudentInfo struct {
	gorm.Model
	semester string
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
