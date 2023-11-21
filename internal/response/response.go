package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

const (
	successCode = 0
	errorCode   = 10000

	successMsg = "成功"
	errMsg     = "失败"
)

func Result(code int, data any, msg string, c *gin.Context) {
	// 开始时间
	c.JSON(http.StatusOK, response{
		code,
		msg,
		data,
	})
}

func Ok(c *gin.Context) {
	Result(successCode, map[string]any{}, successMsg, c)
}

func OkWithMessage(message string, c *gin.Context) {
	Result(successCode, map[string]any{}, message, c)
}

func OkWithData(data any, c *gin.Context) {
	Result(successCode, data, successMsg, c)
}

func OkWithDetailed(data any, message string, c *gin.Context) {
	Result(successCode, data, message, c)
}

func Fail(c *gin.Context) {
	Result(errorCode, map[string]any{}, errMsg, c)
}

func FailWithMessage(message string, c *gin.Context) {
	Result(errorCode, map[string]any{}, message, c)
}

func FailWithDetailed(data any, message string, c *gin.Context) {
	Result(errorCode, data, message, c)
}

func FailWithCode(code int, c *gin.Context) {
	Result(code, map[string]any{}, mapMsg[code], c)
}
