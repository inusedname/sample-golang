package middleware

import (
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// User demo
type User struct {
	UserName  string
	FirstName string
	LastName  string
}

type login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

var (
	IdentityKey = "id"
)

func InitParams() *jwt.GinJWTMiddleware {

	return &jwt.GinJWTMiddleware{
		Realm:       "test zone",
		Key:         []byte("secret key"),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: IdentityKey,
		PayloadFunc: payloadFunc(),

		IdentityHandler: identityHandler(),
		Authenticator:   authenticator(),
		Authorizator:    authorizator(),
		Unauthorized:    unauthorized(),
		TokenLookup:     "header: Authorization, query: token, cookie: jwt",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	}
}

func identityHandler() func(c *gin.Context) any {
	return func(c *gin.Context) any {
		claims := jwt.ExtractClaims(c)
		return &User{
			UserName: claims[IdentityKey].(string),
		}
	}
}

func payloadFunc() func(data any) jwt.MapClaims {
	return func(data any) jwt.MapClaims {
		if v, ok := data.(*User); ok {
			return jwt.MapClaims{
				IdentityKey: v.UserName,
			}
		}
		return jwt.MapClaims{}
	}
}

func authenticator() func(c *gin.Context) (any, error) {
	return func(c *gin.Context) (any, error) {
		var loginVals login
		if err := c.ShouldBind(&loginVals); err != nil {
			return "", jwt.ErrMissingLoginValues
		}
		userID := loginVals.Username
		password := loginVals.Password

		if (userID == "admin" && password == "admin") || (userID == "test" && password == "test") {
			return &User{
				UserName:  userID,
				LastName:  "Bo-Yi",
				FirstName: "Wu",
			}, nil
		}
		return nil, jwt.ErrFailedAuthentication
	}
}

func unauthorized() func(c *gin.Context, code int, message string) {
	return func(c *gin.Context, code int, message string) {
		c.JSON(code, gin.H{
			"code":    code,
			"message": message,
		})
	}
}

func authorizator() func(data any, c *gin.Context) bool {
	return func(data any, c *gin.Context) bool {
		if v, ok := data.(*User); ok && v.UserName == "admin" {
			return true
		}
		return false
	}
}
