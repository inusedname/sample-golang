package authorizer

import (
	"github.com/gin-gonic/gin"
)

func AuthorizeRole(requiredRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}
