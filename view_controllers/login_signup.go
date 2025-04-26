package viewcontrollers

import (
	"app/use_cases"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type signUpData struct {
	Error string
}

func HandleGetSignup(c *gin.Context) {
	c.HTML(http.StatusOK, "signup.html", nil)
}

type signUpPayload struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
	Email    string `form:"email" json:"email" binding:"required"`
	FullName string `form:"fullName" json:"fullName" binding:"required"`
	Role     string `form:"role" json:"role" binding:"required"`
}

func HandlePostSignup(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		if c.Request.ParseForm() != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		var payload signUpPayload
		if err := c.ShouldBind(&payload); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		username := payload.Username
		password := payload.Password
		email := payload.Email
		fullName := payload.FullName
		role := payload.Role

		if err := use_cases.CreateUser(db, username, password, email, fullName, role); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
	}
}
