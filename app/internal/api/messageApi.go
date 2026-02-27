package api

import (
	"commmunity/app/internal/response"
	"commmunity/app/internal/service/controller"
	"commmunity/app/zlog"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetHistoryMessage(c *gin.Context) {
	formUserId := c.MustGet("userId").(uint)
	i, err := strconv.ParseUint(c.Param("Id"), 10, 64)
	if err != nil {
		zlog.Error("转换失败")
		response.Fail(c)
		return
	}
	toUserId := uint(i)
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		zlog.Warn("请求出错了")
		response.FailWithCode(c, response.INVALID_PARAMS, response.GetMsg(response.INVALID_PARAMS))
		return
	}
	pageSize := 10
	offset := (page - 1) * pageSize
	messages, err := controller.GetHistoryMessage(formUserId, toUserId, offset, pageSize)
	if err != nil {
		response.FailWithCode(c, response.INTERNAL_ERROR, response.GetMsg(response.INTERNAL_ERROR))
		return
	}
	response.OkWithData(c, messages)
}

func GetNotice(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		zlog.Warn("请求出错了")
		response.FailWithCode(c, response.INVALID_PARAMS, response.GetMsg(response.INVALID_PARAMS))
		return
	}
	pageSize := 10
	offset := (page - 1) * pageSize
	userId := c.MustGet("userId").(uint)
	notices, err := controller.GetUnreadNotice(userId, offset, pageSize)
	if err != nil {
		response.FailWithCode(c, response.INTERNAL_ERROR, response.GetMsg(response.INTERNAL_ERROR))
		return
	}
	response.OkWithData(c, notices)
}
