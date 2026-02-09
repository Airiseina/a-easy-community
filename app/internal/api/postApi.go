package api

import (
	"commmunity/app/internal/model"
	"commmunity/app/internal/response"
	"commmunity/app/internal/service/controller"
	"commmunity/app/zlog"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreatePost(c *gin.Context) {
	var post model.PostRequest
	if err := c.ShouldBindJSON(&post); err != nil {
		zlog.Warn("请求出错了")
		response.FailWithCode(c, response.INVALID_PARAMS, response.GetMsg(response.INVALID_PARAMS))
		return
	}
	account := c.GetString("account")
	err, flag := controller.CreatePost(account, post.Title, post.Content)
	if err != nil {
		response.FailWithCode(c, response.INTERNAL_ERROR, response.GetMsg(response.INTERNAL_ERROR))
		return
	}
	if !flag {
		response.FailWithMessage(c, "你已被禁言")
		return
	}
	response.Ok(c)
}

func GetPostList(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		zlog.Warn("请求出错了")
		response.FailWithCode(c, response.INVALID_PARAMS, response.GetMsg(response.INVALID_PARAMS))
		return
	}
	pageSize := 10 //每页十条文章
	offset := (page - 1) * pageSize
	posts, err := controller.GetPostList(offset, pageSize)
	if err != nil {
		response.FailWithCode(c, response.INTERNAL_ERROR, response.GetMsg(response.INTERNAL_ERROR))
		return
	}
	response.OkWithData(c, posts)
}

func GetPostDetail(c *gin.Context) {
	postId, err := strconv.ParseUint(c.Param("postId"), 10, 64)
	if err != nil {
		zlog.Error("转换失败")
		response.Fail(c)
		return
	}
	postIdInt := uint(postId)
	post, err := controller.GetPostDetail(postIdInt)
	if err != nil {
		response.FailWithCode(c, response.INTERNAL_ERROR, response.GetMsg(response.INTERNAL_ERROR))
		return
	}
	response.OkWithData(c, post)
}

func CreateComment(c *gin.Context) {
	var comment model.CommentRequest
	if err := c.ShouldBindJSON(&comment); err != nil {
		zlog.Warn("请求失败")
		response.FailWithCode(c, response.INVALID_PARAMS, response.GetMsg(response.INVALID_PARAMS))
		return
	}
	account := c.GetString("account")
	postId, err := strconv.ParseUint(c.Param("postId"), 10, 64)
	if err != nil {
		zlog.Error("转换失败")
		response.Fail(c)
	}
	postIdInt := uint(postId)
	err, flag := controller.CreateComment(account, postIdInt, comment.Content)
	if err != nil {
		response.FailWithCode(c, response.INTERNAL_ERROR, response.GetMsg(response.INTERNAL_ERROR))
		return
	}
	if !flag {
		response.FailWithMessage(c, "你已被禁言")
		return
	}
	response.Ok(c)
}

func GetUserProfile(c *gin.Context) {
	i, err := strconv.ParseUint(c.Param("Id"), 10, 64)
	if err != nil {
		zlog.Error("转换失败")
		response.Fail(c)
	}
	id := uint(i)
	user, err := controller.GetUserProfile(id)
	if err != nil {
		response.FailWithCode(c, response.INTERNAL_ERROR, response.GetMsg(response.INTERNAL_ERROR))
		return
	}
	response.OkWithData(c, *user)
}

func DeletePost(c *gin.Context) {
	account := c.GetString("account")
	role := c.MustGet("role").(int)
	postId, err := strconv.ParseUint(c.Param("postId"), 10, 64)
	if err != nil {
		zlog.Error("转换失败")
		response.Fail(c)
	}
	postIdInt := uint(postId)
	err, flag := controller.DeletePost(account, postIdInt, role)
	if err != nil {
		response.FailWithCode(c, response.INTERNAL_ERROR, response.GetMsg(response.INTERNAL_ERROR))
		return
	}
	if !flag {
		response.FailWithMessage(c, "你没有该权限")
		return
	}
	response.Ok(c)
}

func DeleteComment(c *gin.Context) {
	commentId, err := strconv.ParseUint(c.Param("commentId"), 10, 64)
	if err != nil {
		zlog.Error("转换失败")
		response.Fail(c)
	}
	commentIdInt := uint(commentId)
	account := c.GetString("account")
	role := c.MustGet("role").(int)
	err, flag := controller.DeleteComment(account, commentIdInt, role)
	if err != nil {
		response.FailWithCode(c, response.INTERNAL_ERROR, response.GetMsg(response.INTERNAL_ERROR))
		return
	}
	if !flag {
		response.FailWithMessage(c, "你没有该权限")
		return
	}
	response.Ok(c)
}
