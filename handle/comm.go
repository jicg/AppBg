package handle

import (
	"github.com/gin-gonic/gin"
	"github.com/jicg/AppBg/middleware/jwt"
	"github.com/jicg/AppBg/bean"
	"github.com/jicg/AppBg/db"
)

func FindUserId(c *gin.Context) (uint) {
	claims := c.MustGet("claims").(*jwtauth.CustomClaims)
	return claims.ID
}

func FindUser(c *gin.Context) (*bean.User) {
	u := &db.User{}
	u.ID = FindUserId(c)
	user := new(bean.User)
	db.FindUser(u, user)
	return user
}
