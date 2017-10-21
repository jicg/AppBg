package offorder

import (
	"github.com/jicg/AppBg/middleware/wxpay"
	"github.com/gin-gonic/gin"
	"github.com/jicg/AppBg/db"
	"github.com/jicg/AppBg/bean"
	"encoding/json"
	"github.com/pkg/errors"
)

var (
	payHandler *wxpay.APIHandler
)

func getAPIHandlerByCache() (*bean.APIHandler, error) {
	cache, err := db.GetCache("wxpay_key")
	if err != nil {
		return nil, err
	}
	api := new(bean.APIHandler)
	if err := json.Unmarshal([]byte(cache), api); err != nil {
		return nil, err
	}
	return api, nil
}

func LoadPayHandle(c *gin.Context, call wxpay.BusFunc) (*wxpay.APIHandler, error) {
	if payHandler == nil {
		api, err := getAPIHandlerByCache();
		if err != nil {
			return nil, err
		}
		if len(api.PayKey) == 0 || len(api.MchId) == 0 || len(api.PayKey) == 0 {
			return nil, errors.New("支付参数非法！，请维护好支付参数")
		}
		payHandler = &wxpay.APIHandler{
			AppId:       api.AppId,
			MchId:       api.MchId,
			PayKey:      api.PayKey,
			NotifyUrl:   "http://" + c.Request.Host + "/wxpay",
			InvokeApiIp: c.ClientIP(),
			Busfunc:     call,
		}
	} else {
		payHandler.InvokeApiIp = c.ClientIP()
		payHandler.NotifyUrl = "http://" + c.Request.Host + "/wxpay"
	}
	return payHandler, nil
}

func ReloadPayHandle() {
	payHandler = nil
}
