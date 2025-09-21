package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 标准响应结构
type Response struct {
	Code int         `json:"code"` // 业务状态码
	Msg  string      `json:"msg"`  // 提示信息
	Data interface{} `json:"data"` // 返回数据
}

// Success 成功响应
func Success(c *gin.Context, msg string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code: http.StatusOK,
		Msg:  msg,
		Data: data,
	})
}

// Error 失败响应（默认500错误）
func Error(c *gin.Context, msg string) {
	c.JSON(http.StatusInternalServerError, Response{
		Code: http.StatusInternalServerError,
		Msg:  msg,
		Data: nil,
	})
}

// Custom 自定义响应
func Custom(c *gin.Context, code int, msg string, data interface{}) {
	c.JSON(code, Response{
		Code: code,
		Msg:  msg,
		Data: data,
	})
}

// ParamError 参数错误响应（预定义400错误）
func ParamError(c *gin.Context) {
	c.JSON(http.StatusBadRequest, Response{
		Code: http.StatusBadRequest,
		Msg:  "请求参数错误",
		Data: nil,
	})
}

// AuthError 鉴权失败响应（预定义401错误）
func AuthError(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, Response{
		Code: http.StatusUnauthorized,
		Msg:  "身份验证失败",
		Data: nil,
	})
}
