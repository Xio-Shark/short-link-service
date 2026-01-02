package handler

import "github.com/gin-gonic/gin"

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func respondOK(c *gin.Context, data interface{}) {
	c.JSON(200, Response{Code: 0, Message: "ok", Data: data})
}

func respondError(c *gin.Context, httpStatus int, code int, message string) {
	c.JSON(httpStatus, Response{Code: code, Message: message})
}
