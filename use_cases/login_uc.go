package use_cases

import (
	models "app/data"
	"app/errors"

	"gorm.io/gorm"
)

func CreateUser(db *gorm.DB, username string, password string, email string, fullName string) error {
	if db.Find(&models.User{}, "email = ?", email).RowsAffected != 0 {
		return errors.ErrUserAlreadyExists{}
	}
	user := models.User{
		FullName: fullName,
		Email:    email,
		UserCredential: models.UserCredential{
			Username: username,
			Password: password,
		},
		UserPermissionDef: models.NewUserPermissionDef("student"),
	}
	db.Create(&user)
	return nil
}
