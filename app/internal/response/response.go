package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code int         `json:"code"` // 业务状态码
	Msg  string      `json:"msg"`  // 提示信息
	Data interface{} `json:"data"` // 数据
}

const (
	SUCCESS                                = 0 //成功
	ERROR                                  = 1 //错误
	INTERNAL_ERROR                         = 2 //内部错误
	INVALID_PARAMS                         = 3 //参数错误
	ERROR_USERNAME_EXIST                   = 4 //用户已存在
	ERROR_USER_NOT_EXIST_OR_PASSWORD_WRONG = 5 //用户不存在或密码错误
	ERROR_AUTH_CHECK_TOKEN_FAIL            = 6 //jwt鉴权错误
)

var Error = map[int]string{
	INTERNAL_ERROR:                         "系统繁忙",
	INVALID_PARAMS:                         "请求出错",
	ERROR_USERNAME_EXIST:                   "用户已存在",
	ERROR_USER_NOT_EXIST_OR_PASSWORD_WRONG: "账号或密码错误",
	ERROR_AUTH_CHECK_TOKEN_FAIL:            "Token鉴权失败",
}

func Result(c *gin.Context, code int, msg string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code: code,
		Msg:  msg,
		Data: data,
	})
}

func Ok(c *gin.Context) {
	Result(c, SUCCESS, "操作成功", map[string]interface{}{})
}

func OkWithData(c *gin.Context, data interface{}) {
	Result(c, SUCCESS, "操作成功", data)
}

func Fail(c *gin.Context) {
	Result(c, ERROR, "操作失败", map[string]interface{}{})
}

func FailWithMessage(c *gin.Context, message string) {
	Result(c, ERROR, message, map[string]interface{}{})
}

func FailWithCode(c *gin.Context, code int, message string) {
	Result(c, code, message, map[string]interface{}{})
}

func GetMsg(code int) string {
	msg, _ := Error[code]
	return msg
}
