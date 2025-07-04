package resp

type ResultV2 struct {
	Result    bool        `json:"result"`
	Code      int         `json:"code"`
	Msg       string      `json:"message"`
	Data      interface{} `json:"data"`
	RequestId string      `json:"request_id"`
}

const (
	Succ        = 1
	Fail        = 0
	MSG_SUCCESS = "success"
)

// S 快速构造成功的返回对象
func S(data interface{}, reqId string) ResultV2 {
	r := ResultV2{RequestId: reqId, Result: true, Data: data, Code: Succ, Msg: MSG_SUCCESS}
	return r
}

// E 快速构造错误的返回对象
func E(msg string, reqId string) ResultV2 {
	r := ResultV2{Result: false, Data: "", Code: Fail, Msg: msg, RequestId: reqId}
	return r
}

func NewResultV2(result bool, data interface{}, c int, reqId string, m ...string) ResultV2 {
	r := ResultV2{Result: result, Data: data, Code: c, RequestId: reqId}

	if e, ok := data.(error); ok {
		if m == nil {
			r.Msg = e.Error()
		}
	} else {
		r.Msg = "SUCCESS"
	}
	if len(m) > 0 {
		r.Msg = m[0]
	}

	return r
}
