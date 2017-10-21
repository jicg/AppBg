package login

import (
	"fmt"
	"github.com/jicg/AppBg/bean/ret"
	"github.com/jicg/AppBg/middleware/jwt"
	"time"
	"github.com/gin-gonic/gin"
	"encoding/json"
	"github.com/jicg/AppBg/db"
	"net/http"
	"github.com/dgrijalva/jwt-go"
)

func Login(c *gin.Context) {
	bs, _ := c.GetRawData()
	puser := new(struct {
		Username string
		Password string
	})
	err := json.Unmarshal(bs, &puser)
	if err != nil {
		fmt.Println(err.Error())
	}
	u := db.QueryUser(&db.User{Name: puser.Username, Pwd: puser.Password})
	if u.ID == 0 {
		c.JSON(http.StatusOK, ret.Login_Error("用户不能存在！"))
		return
	}
	j := jwtauth.NewJWT()
	claims := jwtauth.CustomClaims{
		u.ID,
		u.Name,
		u.Email,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(2 * time.Hour).Unix(),
			Issuer:    "jicg",
		},
	}
	token, err := j.CreateToken(claims)
	if err != nil {
		c.JSON(http.StatusOK, ret.Login_Error(err.Error()))
		c.Abort()
	}
	c.JSON(http.StatusOK, ret.Login_Success(token))
}
