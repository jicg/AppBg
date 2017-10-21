package user

import (
	"github.com/gin-gonic/gin"
	"github.com/jicg/AppBg/bean/ret"
	"net/http"
	"github.com/jicg/AppBg/handle"
	"github.com/jicg/AppBg/db"
	"encoding/json"
	"github.com/jicg/AppBg/bean"
)

func GetUser(c *gin.Context) {
	c.JSON(http.StatusOK, ret.Success(ret.H{"user": handle.FindUser(c)}));
}

func ChangePwd(c *gin.Context) {
	id := handle.FindUserId(c)
	bs, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusOK, ret.Error(err.Error()))
		return
	}
	pwdchange := &bean.PwdChange{}
	if err = json.Unmarshal(bs, pwdchange); err != nil {
		c.JSON(http.StatusOK, ret.Error(err.Error()))
		return
	}
	if len(pwdchange.Pwd) == 0 {
		c.JSON(http.StatusOK, ret.Error("密码不能为空"))
		return
	}
	if err = db.ChangePwd(id, pwdchange.Pwd); err != nil {
		c.JSON(http.StatusOK, ret.Error(err.Error()))
		return
	}
	c.JSON(http.StatusOK, ret.Success());
}
