package middleware

import (
	models "app/data"
	"log"
	"net/http"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

var (
	IdentityKey = "id"
)

func InitParams(db *gorm.DB) *jwt.GinJWTMiddleware {

	return &jwt.GinJWTMiddleware{
		Realm:       "test zone",
		Key:         []byte("secret key"),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: IdentityKey,
		PayloadFunc: payloadFunc(),

		IdentityHandler: identityHandler(db),
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
			c.Redirect(http.StatusFound, "/classes")
		},
	}
}

func payloadFunc() func(data any) jwt.MapClaims {
	log.Printf("payloadFunc")
	return func(data any) jwt.MapClaims {
		if v, ok := data.(*models.User); ok {
			return jwt.MapClaims{
				IdentityKey: v.ID,
			}
		}
		return jwt.MapClaims{}
	}
}

func authenticator(db *gorm.DB) func(c *gin.Context) (any, error) {
	log.Printf("authenticator")
	return func(c *gin.Context) (any, error) {
		var loginVals login
		if err := c.ShouldBind(&loginVals); err != nil {
			return "", jwt.ErrMissingLoginValues
		}
		userID := loginVals.Username
		password := loginVals.Password

		var userCred models.UserCredential
		result := db.Where("username = ? AND password = ?", userID, password).First(&userCred)

		if result.Error != nil {
			log.Printf("Database error: %v", result.Error)
			return nil, jwt.ErrFailedAuthentication
		}
		var user models.User
		db.Where("user_credential_id = ?", userCred.ID).First(&user)
		return &user, nil
	}
}

func identityHandler(db *gorm.DB) func(c *gin.Context) interface{} {
	return func(c *gin.Context) interface{} {
		claims := jwt.ExtractClaims(c)
		var user models.User
		db.Preload("UserPermissionDef").Where("id = ?", claims[IdentityKey]).First(&user)
		return &user
	}
}

func unauthorized() func(c *gin.Context, code int, message string) {
	return func(c *gin.Context, code int, message string) {
		log.Printf("unauthorized, redirect to login.")
		c.Redirect(http.StatusFound, "/login")
	}
}
