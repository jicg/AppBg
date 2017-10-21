package ret

import (
	"encoding/json"
)

type Ret struct {
	Code int32       `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func (ret *Ret) Bytes() ([]byte) {
	bs, _ := json.Marshal(ret)
	return bs
}

type H map[string]interface{}

func Success(args ...interface{}) *Ret {
	ret := &Ret{
		Code: 0,
		Msg:  "操作成功",
		Data: nil,
	}
	for _, value := range args {

		switch value.(type) {
		case int, int32, int64, uint:
			ret.Code = value.(int32)
		case string:
			ret.Msg = value.(string)
		default:
			ret.Data = value.(interface{})
		}
	}
	return ret
}

func Error(args ...interface{}) *Ret {
	ret := &Ret{
		Code: -1,
		Msg:  "操作失败",
		Data: nil,
	}
	for index, value := range args {
		switch args[index].(type) {
		case int, int32, int64, uint:
			code := value.(int32);
			if (code >= 0) {
				code = code * -1;
			}
			ret.Code = code
		case string:
			ret.Msg = value.(string)
		default:
			ret.Data = value.(interface{})
		}
	}
	return ret
}
