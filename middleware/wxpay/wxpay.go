package wxpay

import (
	"sort"
	"crypto/md5"
	"strings"
	"encoding/hex"
	"encoding/xml"
	"net/http"
	"bytes"
	"io/ioutil"
	"time"
	"math/rand"
	"fmt"
	log4go  "github.com/sirupsen/logrus"
)

type UnifyOrderReq struct {
	Appid            string `xml:"appid"`
	Body             string `xml:"body"`
	Mch_id           string `xml:"mch_id"`
	Nonce_str        string `xml:"nonce_str"`
	Notify_url       string `xml:"notify_url"`
	Trade_type       string `xml:"trade_type"`
	Spbill_create_ip string `xml:"spbill_create_ip"`
	Total_fee        int    `xml:"total_fee"`
	Out_trade_no     string `xml:"out_trade_no"`
	Sign             string `xml:"sign"`
	//Device_info      string `xml:"device_info"`
	Detail string `xml:"detail"`
	//Attach           string `xml:"attach"`
}

type UnifyOrderResp struct {
	Return_code  string `xml:"return_code"`
	Return_msg   string `xml:"return_msg"`
	Appid        string `xml:"appid"`
	Mch_id       string `xml:"mch_id"`
	Nonce_str    string `xml:"nonce_str"`
	Sign         string `xml:"sign"`
	Result_code  string `xml:"result_code"`
	Prepay_id    string `xml:"prepay_id"`
	Trade_type   string `xml:"trade_type"`
	Code_url     string `xml:"code_url"`
	Err_code     string `xml:"err_code"`
	Err_code_des string `xml:"err_code_des"`
}

var (
	//	//地址
	WX_UNIFIEDORDER_API = "https://api.mch.weixin.qq.com/pay/unifiedorder"
	//	WX_APP_ID           = ""
	//	WX_MCH_ID           = ""
	//	WX_PAY_KEY          = ""
	//	WX_NOTIFY_URL       = ""
	//	WX_INVOKE_API_IP    = ""
)

type APIHandler struct {
	AppId       string
	MchId       string
	PayKey      string
	NotifyUrl   string
	InvokeApiIp string
	Tradetype   string
	Busfunc     BusFunc
}

func (o *APIHandler) WxUnifyChargeReq(orderNo string, body string, price float32) (*UnifyOrderResp, error) {
	return o.WxUnifyChargeReqWithDetail(orderNo, body, body, price)
}

func (o *APIHandler) WxUnifyChargeReqWithDetail(orderNo string, body string, detail string, price float32) (*UnifyOrderResp, error) {
	if len(o.Tradetype) == 0 {
		o.Tradetype = "NATIVE"
	}
	nonceStr := getRandomString(32)
	//请求UnifiedOrder的代码
	var yourReq UnifyOrderReq
	yourReq.Appid = o.AppId //微信开放平台我们创建出来的app的app id
	yourReq.Detail = detail
	yourReq.Body = body
	yourReq.Mch_id = o.MchId
	yourReq.Nonce_str = nonceStr
	yourReq.Notify_url = o.NotifyUrl
	yourReq.Trade_type = "NATIVE"
	yourReq.Spbill_create_ip = o.InvokeApiIp
	yourReq.Total_fee = int(price * 100) //单位是分
	yourReq.Out_trade_no = orderNo
	var m map[string]interface{}
	m = make(map[string]interface{}, 0)
	m["appid"] = yourReq.Appid
	m["body"] = yourReq.Body
	m["mch_id"] = yourReq.Mch_id
	m["notify_url"] = yourReq.Notify_url
	m["trade_type"] = yourReq.Trade_type
	m["spbill_create_ip"] = yourReq.Spbill_create_ip
	m["total_fee"] = yourReq.Total_fee
	m["out_trade_no"] = yourReq.Out_trade_no
	m["nonce_str"] = yourReq.Nonce_str
	m["detail"] = yourReq.Detail
	yourReq.Sign = wxpayCalcSign(m, o.PayKey) // 这个是计算wxpay签名的函数上面已贴出 WX_PAY_KEY
	bytes_req, err := xml.Marshal(yourReq)
	if err != nil {
		log4go.Error("wxUnifyChargeReq(): xml.Marshal error:%s", err)
		return nil, err
	}
	str_req := string(bytes_req)
	//wxpay的unifiedorder接口需要http body中xmldoc的根节点是<xml></xml>这种，所以这里需要replace一下
	str_req = strings.Replace(str_req, "XUnifyOrderReq", "xml", -1)
	bytes_req = []byte(str_req)
	//发送unified order请求.
	req, err := http.NewRequest("POST", WX_UNIFIEDORDER_API, bytes.NewReader(bytes_req))
	if err != nil {
		log4go.Error("wxUnifyChargeReq(): http.NewRequest error:%s", err)
		return nil, err
	}
	req.Header.Set("Accept", "application/xml")
	//这里的http header的设置是必须设置的.
	req.Header.Set("Content-Type", "application/xml;charset=utf-8")
	c := http.Client{}
	resp, _err := c.Do(req)
	if _err != nil {
		log4go.Error("wxUnifyChargeReq(): http.Do error:%s", _err)
		return nil, err
	}
	defer resp.Body.Close()
	bs, err := ioutil.ReadAll(resp.Body)
	xmlResp := &UnifyOrderResp{}
	_err = xml.Unmarshal(bs, xmlResp)
	if _err != nil {
		return nil, _err
	}
	if handles == nil {
		handles = make(HandlersMap, 0);
	}
	handles[orderNo] = o
	return xmlResp, nil
}

func getRandomString(size int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < size; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func wxpayVerifySign(needVerifyM map[string]interface{}, key string, sign string) bool {
	signCalc := wxpayCalcSign(needVerifyM, key)

	log4go.Infoln("计算出来的sign: %v", signCalc)
	log4go.Infoln("微信异步通知sign: %v", sign)
	if sign == signCalc {
		return true
	}
	return false
}

//wxpay计算签名的函数
func wxpayCalcSign(mReq map[string]interface{}, key string) (sign string) {
	log4go.Info("wxpayCalcSign()...API KEY:%s", key)
	//STEP 1, 对key进行升序排序.
	sorted_keys := make([]string, 0)
	for k, _ := range mReq {
		sorted_keys = append(sorted_keys, k)
	}
	sort.Strings(sorted_keys)
	//STEP2, 对key=value的键值对用&连接起来，略过空值
	var signStrings string
	for _, k := range sorted_keys {
		value := fmt.Sprintf("%v", mReq[k])
		if value != "" {
			signStrings = signStrings + k + "=" + value + "&"
		}
	}
	//STEP3, 在键值对的最后加上key=API_KEY
	if key != "" {
		signStrings = signStrings + "key=" + key
	}
	log4go.Info("wxpayCalcSign()...signStrings:%s", signStrings)
	//STEP4, 进行MD5签名并且将所有字符转为大写.
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(signStrings))
	cipherStr := md5Ctx.Sum(nil)
	upperSign := strings.ToUpper(hex.EncodeToString(cipherStr))
	return upperSign
}
