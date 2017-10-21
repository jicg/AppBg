package bean

type User struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type QueryOrder struct {
	Orderno   string  `json:"orderno"`
	Optname   string  `json:"optname"`
	Remark    string  `json:"remark"`
	Amt       float32 `json:"amt"`
	Status    uint8   `json:"status"`
	Product   string  `json:"product"`
	Sttype    uint8   `json:"sttype"`
	PageIndex int     `json:"pageIndex"`
	PageSize  int     `json:"pageSize"`
}

type APIHandler struct {
	AppId  string `json:"appid"`
	MchId  string `json:"mchid"`
	PayKey string `json:"paykey"`
}

type PwdChange struct {
	Pwd string `json:"pwd"`
}
