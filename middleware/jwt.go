package middleware

import (
	"app/constants"
	models "app/data"
	mystrconv "app/utils/strconv"
	"fmt"
	"log"
	"net/http"
	"time"

	ginjwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

func InitParams(db *gorm.DB) *ginjwt.GinJWTMiddleware {

	return &ginjwt.GinJWTMiddleware{
		Realm:           "test zone",
		Key:             []byte("secret key"),
		Timeout:         time.Hour,
		MaxRefresh:      time.Hour,
		IdentityKey:     constants.IdentityKey,
		PayloadFunc:     payloadFunc(),
		IdentityHandler: identityHandler(),
		Authenticator:   authenticator(db),
		Unauthorized:    unauthorized(),
		TokenLookup:     "header: Authorization, query: token, cookie: jwt",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,

		SendCookie:     true,
		SecureCookie:   false,
		CookieHTTPOnly: true,
		CookieDomain:   "localhost:8080",
		CookieSameSite: http.SameSiteDefaultMode,

		LoginResponse: func(c *gin.Context, code int, message string, time time.Time) {
		},
	}
}

type MyClaims struct {
	UserID uint
}

func identityHandler() func(*gin.Context) any {
	return func(c *gin.Context) any {
		mapClaims := ginjwt.ExtractClaims(c)
		user_id, _ := mystrconv.Atoui(fmt.Sprintf("%f", mapClaims[constants.IdentityKey]))
		return MyClaims{
			UserID: user_id,
		}
	}
}

func payloadFunc() func(data any) ginjwt.MapClaims {
	return func(data any) ginjwt.MapClaims {
		if v, ok := data.(*models.User); ok {
			return ginjwt.MapClaims{
				constants.IdentityKey: v.ID,
			}
		}
		return ginjwt.MapClaims{}
	}
}

func authenticator(db *gorm.DB) func(c *gin.Context) (any, error) {
	return func(c *gin.Context) (any, error) {
		var loginVals login
		if err := c.ShouldBind(&loginVals); err != nil {
			return "", ginjwt.ErrMissingLoginValues
		}
		userID := loginVals.Username
		password := loginVals.Password

		var userCred models.UserCredential
		result := db.Where("username = ? AND password = ?", userID, password).First(&userCred)

		if result.Error != nil {
			log.Printf("Database error: %v", result.Error)
			return nil, ginjwt.ErrFailedAuthentication
		}
		var user models.User
		db.Where("user_credential_id = ?", userCred.ID).First(&user)
		return &user, nil
	}
}

func unauthorized() func(c *gin.Context, code int, message string) {
	return func(c *gin.Context, code int, message string) {
	}
}
