package router

// BaseController 基类
type BaseController struct {
	Data interface{}
}

type response struct {
	Code    int         `json:"code"`
	Message string      `json:"msg"`
	Data    interface{} `json:"data"`
}

// Success ????
func Success(message string, data interface{}) response {
	res := response{}
	// res.Code = 200
	res.Code = 0
	res.Message = message
	res.Data = data
	return res
}

// Error ????
func Error(message string, code int, data interface{}) response {
	res := response{}
	res.Code = code
	res.Message = message
	res.Data = data
	return res
}
