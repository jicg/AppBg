package pay

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/jicg/AppBg/middleware/wxpay"
	"encoding/json"
	"fmt"
	"github.com/boombuler/barcode/qr"
	"github.com/boombuler/barcode"
	"github.com/jicg/AppBg/bean/ret"
	"image/png"
	log "github.com/sirupsen/logrus"
)

func WxPay(c *gin.Context) {
	orderno := "T0001"
	if no, flag := c.Params.Get("no"); flag {
		orderno = no
	}
	pay := wxpay.APIHandler{
		AppId:       "wx20b2c8a4d6fecd67",
		MchId:       "1363077502",
		PayKey:      "rzale3mm5pnyrj10sv4yzydzeal0xpmz",
		NotifyUrl:   "http://" + c.Request.Host + "/wxpay",
		InvokeApiIp: c.ClientIP(),
		Busfunc:     call,
	}
	str, e := pay.WxUnifyChargeReq(orderno, "夏老板，测试测试测试，测试，测试，测试", 0.01)
	if e != nil {
		log.Error(e.Error())
	}
	if str.Result_code == "FAIL" {
		log.Error("%s", str.Err_code_des)
		c.JSON(http.StatusOK, ret.Error(ret.H{"错误：": str.Err_code_des}))
		return
	}
	qrcode := qrcode(str.Code_url)
	png.Encode(c.Writer, qrcode)
}

func qrcode(base64 string) (barcode.Barcode) {
	code, err := qr.Encode(base64, qr.L, qr.Unicode)
	if err != nil {
		log.Fatal(err.Error())
	}
	if base64 != code.Content() {
		log.Fatal("data differs")
	}
	code, err = barcode.Scale(code, 300, 300)
	if err != nil {
		log.Fatal(err.Error())
	}
	return code
}

func call(handle *wxpay.APIHandler, req *wxpay.WXPayNotifyReq) {
	bytes, _ := json.Marshal(req)
	fmt.Println(string(bytes))
	log.Info(string(bytes))
}
