package data

import (
	"app/constants"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetUser(db *gorm.DB, c *gin.Context) User {
	userId := jwt.ExtractClaims(c)[constants.IdentityKey]
	var user User
	db.Preload("UserPermissionDef").Where("id = ?", userId).First(&user)
	return user
}
