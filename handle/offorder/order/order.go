package order

import (
	"strconv"
	"github.com/jicg/AppBg/bean/ret"
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/jicg/AppBg/db"
	"github.com/jicg/AppBg/handle/offorder"
	"github.com/jicg/AppBg/middleware/wxpay"
	log "github.com/sirupsen/logrus"
	"image/png"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"encoding/json"
	"gopkg.in/olahol/melody.v1"
	"math/rand"
	"sync"
	"fmt"
	"github.com/xuri/excelize"
	"time"
)

var (
	m        = newMelody()
	lock     = new(sync.Mutex)
	wshelper = map[int][]*melody.Session{}
)

type WsAction struct {
	Action int    `json:"action"`
	Msg    string `json:"msg"`
}

func newMelody() *melody.Melody {
	mm := melody.New()
	mm.HandleMessage(func(session *melody.Session, bytes []byte) {
		lock.Lock()
		defer lock.Unlock()
		fmt.Println("HandleMessage", string(bytes))
		param := new(WsAction)
		if err := json.Unmarshal(bytes, param); err != nil {
			return
		}
		bs, _ := json.Marshal(param)
		fmt.Println(string(bs))
		if param.Action == 1 {
			id, _ := strconv.Atoi(param.Msg)
			if wshelper[id] == nil {
				wshelper[id] = []*melody.Session{}
			}
			wshelper[id] = append(wshelper[id], session)
		}
	})
	mm.HandleClose(func(session *melody.Session, i int, s string) error {
		lock.Lock()
		defer lock.Unlock()
		for k, v := range wshelper {
			temp := []*melody.Session{}
			for _, s := range v {
				if s != session && !s.IsClosed() {
					temp = append(temp, s)
				}
			}
			wshelper[k] = temp
		}
		return nil
	})
	return mm
}

func QueryOrders(context *gin.Context) {
	orderno, _ := context.GetQuery("orderno");
	optname, _ := context.GetQuery("optname");
	remark, _ := context.GetQuery("remark");
	product, _ := context.GetQuery("product");
	datebeg, _ := context.GetQuery("datebeg");
	dateend, _ := context.GetQuery("dateend");
	sttype := uint8(0);
	if sttype_str, flag := context.GetQuery("sttype"); flag {
		tmp, _ := strconv.ParseUint(sttype_str, 10, 8)
		sttype = uint8(tmp)
	}
	status := uint8(0);
	if status_str, flag := context.GetQuery("status"); flag {
		tmp, _ := strconv.ParseUint(status_str, 10, 8)
		status = uint8(tmp)
	}
	pageIndex := int(0)
	if pageIndex_str, flag := context.GetQuery("pageIndex"); flag {
		pageIndex, _ = strconv.Atoi(pageIndex_str)
	}
	pageSize := int(0)
	if pageSize_str, flag := context.GetQuery("pageSize"); flag {
		pageSize, _ = strconv.Atoi(pageSize_str)
	}
	data, cnt := db.QueryOfforder(orderno, datebeg, dateend, optname, remark, product, sttype, status, (pageIndex-1)*pageSize, pageSize)

	//time.Sleep(50 * 1e9)

	context.JSON(http.StatusOK, ret.Success(ret.H{
		"total": cnt,
		"data":  data,
	}))
}

func QueryXlsOrders(context *gin.Context) {
	orderno, _ := context.GetQuery("orderno");
	optname, _ := context.GetQuery("optname");
	remark, _ := context.GetQuery("remark");
	product, _ := context.GetQuery("product");
	datebeg, _ := context.GetQuery("datebeg");
	dateend, _ := context.GetQuery("dateend");
	sttype := uint8(0);
	if sttype_str, flag := context.GetQuery("sttype"); flag {
		tmp, _ := strconv.ParseUint(sttype_str, 10, 8)
		sttype = uint8(tmp)
	}
	status := uint8(0);
	if status_str, flag := context.GetQuery("status"); flag {
		tmp, _ := strconv.ParseUint(status_str, 10, 8)
		status = uint8(tmp)
	}

	datas := db.QueryOfforderWithOutPage(orderno, datebeg, dateend, optname, remark, product, sttype, status)

	xlsx := excelize.NewFile()
	if style, err := xlsx.NewStyle(`{"font":{"bold":true,family":"Berlin Sans FB Demi"}}`); err == nil {
		xlsx.SetCellStyle("Sheet1", "A1", "E1", style)
	}

	xlsx.SetCellValue("Sheet1", "A1", "序号")
	xlsx.SetCellValue("Sheet1", "B1", "订单编号")
	xlsx.SetCellValue("Sheet1", "C1", "员工")
	xlsx.SetCellValue("Sheet1", "D1", "出货点")
	xlsx.SetCellValue("Sheet1", "E1", "商品详情")
	xlsx.SetCellValue("Sheet1", "F1", "金额")
	xlsx.SetCellValue("Sheet1", "G1", "状态")
	xlsx.SetCellValue("Sheet1", "H1", "备注")
	for index, v := range datas {
		sttype_name := ""
		if v.Sttype == 1 {
			sttype_name = "公司"
		}
		if v.Sttype == 2 {
			sttype_name = "门店"
		}
		status_name := ""
		if v.Status == 1 {
			status_name = "未付款"
		}
		if v.Status == 2 {
			status_name = "已付款"
		}

		indexStr := index + 2
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("%s[%d]", "A", indexStr), v.ID)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("%s[%d]", "B", indexStr), v.Orderno)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("%s[%d]", "C", indexStr), v.Optname)

		xlsx.SetCellValue("Sheet1", fmt.Sprintf("%s[%d]", "D", indexStr), sttype_name)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("%s[%d]", "E", indexStr), v.Product)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("%s[%d]", "F", indexStr), v.Amt)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("%s[%d]", "G", indexStr), status_name)
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("%s[%d]", "H", indexStr), v.Remark)
	}

	xlsx.SetActiveSheet(1)
	context.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	context.Header("Content-Disposition", "attachment;filename=offorders.xlsx")
	context.Header("Cache-Control", "max-age=0")
	xlsx.Write(context.Writer)
}

func Add(context *gin.Context) {
	order := new(db.Offorder)
	bs, _ := context.GetRawData();
	if err := json.Unmarshal(bs, order); err != nil {
		context.JSON(http.StatusOK, ret.Error(err.Error()))
		return
	}
	order.Status = 1
	order.Billdate = time.Now().Format("20060102");
	if len(order.Orderno) == 0 {

		if (order.Sttype == 1) {
			order.Orderno = "Q" + order.Billdate + getRandomString(4);
		} else {
			order.Orderno = "X" + order.Billdate + getRandomString(4);
		}
	}
	if len(order.Optname) == 0 {
		context.JSON(http.StatusOK, ret.Error("员工不能为空！"))
		return
	}
	if len(order.Product) == 0 {
		context.JSON(http.StatusOK, ret.Error("商品详情不能为空！"))
		return
	}
	if order.Amt <= 0 {
		context.JSON(http.StatusOK, ret.Error("金额非法！"))
		return
	}
	if err := db.AddOfforder(order); err != nil {
		context.JSON(http.StatusOK, ret.Error(err.Error()))
		return
	}
	context.JSON(http.StatusOK, ret.Success(order))
}

func GetById(context *gin.Context) {
	idstr, _ := context.Params.Get("id");
	var (
		id    int64
		err   error
		order *db.Offorder
	)
	if id, err = strconv.ParseInt(idstr, 10, 64); err != nil {
		context.JSON(http.StatusOK, ret.Error(err.Error()))
		return
	}
	if order, err = db.GetOfforderById(int(id)); err != nil {
		context.JSON(http.StatusOK, ret.Error(err.Error()))
		return
	}
	context.JSON(http.StatusOK, ret.Success(order))
}

func Update(context *gin.Context) {

	idstr, _ := context.Params.Get("id");
	var (
		id       int64
		err      error
		neworder *db.Offorder
		order    *db.Offorder
	)
	if id, err = strconv.ParseInt(idstr, 10, 64); err != nil {
		context.JSON(http.StatusOK, ret.Error(err.Error()))
		return
	}

	if order, _ = db.GetOfforderById(int(id)); order == nil {
		context.JSON(http.StatusOK, ret.Error("单据不存在！"))
		return
	}
	if order.Status > 1 {
		context.JSON(http.StatusOK, ret.Error("单据已经支付完成，不允许删除！"))
		return
	}

	neworder = new(db.Offorder)
	bs, _ := context.GetRawData();
	if err = json.Unmarshal(bs, neworder); err != nil {
		context.JSON(http.StatusOK, ret.Error(err.Error()))
		return
	}

	if len(neworder.Optname) == 0 {
		context.JSON(http.StatusOK, ret.Error("员工不能为空！"))
		return
	}
	if len(neworder.Product) == 0 {
		context.JSON(http.StatusOK, ret.Error("商品详情不能为空！"))
		return
	}
	if neworder.Amt <= 0 {
		context.JSON(http.StatusOK, ret.Error("金额非法！"))
		return
	}
	neworder.Orderno = order.Orderno
	neworder.ID = uint(id);
	if err := db.UpdateOfforder(neworder); err != nil {
		context.JSON(http.StatusOK, ret.Error(err.Error()))
		return
	}
	context.JSON(http.StatusOK, ret.Success(order))
}

func Del(context *gin.Context) {
	idstr, _ := context.Params.Get("id");
	var (
		id    int64
		err   error
		order *db.Offorder
	)
	if id, err = strconv.ParseInt(idstr, 10, 64); err != nil {
		context.JSON(http.StatusOK, ret.Error(err.Error()))
		return
	}

	if order, _ = db.GetOfforderById(int(id)); order == nil {
		context.JSON(http.StatusOK, ret.Error("单据不存在！"))
		return
	}
	if order.Status > 1 {
		context.JSON(http.StatusOK, ret.Error("单据已经支付完成，不允许删除！"))
		return
	}
	if err = db.DeleteOfforder(id); err != nil {
		context.JSON(http.StatusOK, ret.Error(err.Error()))
		return
	}
	context.JSON(http.StatusOK, ret.Success())
}

func WxPay(context *gin.Context) {
	idstr, _ := context.Params.Get("id");
	var (
		id        int64
		err       error
		apiHandle *wxpay.APIHandler
		order     *db.Offorder
		wxcode    *wxpay.UnifyOrderResp
	)
	if id, err = strconv.ParseInt(idstr, 10, 64); err != nil {
		context.JSON(http.StatusOK, ret.Error(err.Error()))
		return
	}
	if apiHandle, err = offorder.LoadPayHandle(context, call); err != nil {
		context.JSON(http.StatusOK, ret.Error(err.Error()))
		return
	}
	if order, err = db.GetOfforderById(int(id)); err != nil {
		context.JSON(http.StatusOK, ret.Error(err.Error()))
		return
	}
	if wxcode, err = apiHandle.WxUnifyChargeReqWithDetail(order.Orderno, order.Product, order.Remark, order.Amt); err != nil {
		context.JSON(http.StatusOK, ret.Error(err.Error()))
		return
	}
	if wxcode.Result_code == "FAIL" {
		log.Error("%s", wxcode.Err_code_des)
		context.JSON(http.StatusOK, ret.Error(wxcode.Err_code_des))
		return
	}
	qrcode := qrcode(wxcode.Code_url)
	png.Encode(context.Writer, qrcode)
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
	orderno := req.Out_trade_no;
	order, err := db.GetOfforderByOrderno(orderno)
	id := int(order.ID)
	//重置 session
	for k, v := range wshelper {
		temp := []*melody.Session{}
		for _, s := range v {
			if !s.IsClosed() {
				temp = append(temp, s)
			}
		}
		wshelper[k] = temp
	}

	if err != nil {
		log.Error("fail :%v", err)
		if wshelper[id] != nil {
			m.BroadcastMultiple(ret.Error(err.Error()).Bytes(), wshelper[id])
		}
		return
	}
	order.Status = 2
	if err := db.UpdateOfforder(order); err != nil {
		if wshelper[id] != nil {
			m.BroadcastMultiple(ret.Error(err.Error()).Bytes(), wshelper[id])
		}
		return
	}
	if wshelper[id] != nil {
		m.BroadcastMultiple(ret.Success(order).Bytes(), wshelper[id])
	}
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

func WsPayHandler(c *gin.Context) {
	m.HandleRequest(c.Writer, c.Request)
}

func Text(c *gin.Context) {
	idstr, _ := c.Params.Get("id");
	var (
		id  int64
		err error
	)
	if id, err = strconv.ParseInt(idstr, 10, 64); err != nil {
		c.JSON(http.StatusOK, ret.Error(err.Error()))
		return
	}
	order, err := db.GetOfforderById(int(id))
	fmt.Println(id)
	m.BroadcastMultiple(ret.Success(order).Bytes(), wshelper[int(id)])
}
