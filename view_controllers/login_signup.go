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

func HandlePostSignup(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		if c.Request.ParseForm() != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		username := c.Request.FormValue("username")
		password := c.Request.FormValue("password")
		email := c.Request.FormValue("email")
		fullName := c.Request.FormValue("fullName")
		role := c.Request.FormValue("role")
		if err := use_cases.CreateUser(db, username, password, email, fullName, role); err != nil {
			c.HTML(http.StatusBadRequest, "signup.html", signUpData{
				Error: err.Error(),
			})
			return
		}
		c.Redirect(http.StatusFound, "/login")
	}
}
