package setting

import (
	"github.com/gin-gonic/gin"
	"github.com/jicg/AppBg/bean/ret"
	"net/http"
	"encoding/json"
	"github.com/jicg/AppBg/bean"
	"github.com/jicg/AppBg/db"
)

func GetSetting(context *gin.Context) {
	cache, err := db.GetCache("wxpay_key")
	if err != nil {
		context.JSON(http.StatusOK, ret.Error(err.Error()))
		return
	}
	api := new(bean.APIHandler)
	json.Unmarshal([]byte(cache), api)
	context.JSON(http.StatusOK, ret.Success(api))
}

func Update(context *gin.Context) {
	var (
		bs  []byte
		err error
	)
	if bs, err = context.GetRawData(); err != nil {
		context.JSON(http.StatusOK, ret.Error(err.Error()))
		return
	}
	api := new(bean.APIHandler)
	json.Unmarshal(bs, api)
	if len(api.AppId) == 0 {
		context.JSON(http.StatusOK, "appid 不能为空!")
		return
	}
	if len(api.MchId) == 0 {
		context.JSON(http.StatusOK, "machid 不能为空!")
		return
	}
	if len(api.PayKey) == 0 {
		context.JSON(http.StatusOK, "appkey 不能为空!")
		return
	}
	if err := db.SaveCache("wxpay_key", string(bs)); err != nil {
		context.JSON(http.StatusOK, err.Error());
		return
	}
	context.JSON(http.StatusOK, ret.Success("更新成功"))
}
