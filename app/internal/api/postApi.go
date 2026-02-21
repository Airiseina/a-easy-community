package api

import (
	"commmunity/app/internal/model"
	"commmunity/app/internal/response"
	"commmunity/app/internal/service/controller"
	"commmunity/app/internal/service/feed"
	"commmunity/app/zlog"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
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

func UploadImage(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		zlog.Warn("请求出错了")
		response.FailWithCode(c, response.INVALID_PARAMS, response.GetMsg(response.INVALID_PARAMS))
		return
	}
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".jpg" && ext != ".png" && ext != ".gif" && ext != ".jpeg" {
		zlog.Warn("插入图片格式不对")
		response.FailWithMessage(c, "请插入正确的图片")
		return
	}
	fileName := uuid.New().String() + ext
	err = os.MkdirAll("./uploads/post_file", 0755)
	if err != nil {
		zlog.Error("文件夹创建失败", zap.Error(err))
		response.FailWithCode(c, response.INTERNAL_ERROR, response.GetMsg(response.INTERNAL_ERROR))
		return
	}
	savePath := filepath.Join("uploads", "post_file", fileName)
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		zlog.Error("服务器硬盘出错", zap.Error(err))
		response.FailWithCode(c, response.INTERNAL_ERROR, response.GetMsg(response.INTERNAL_ERROR))
		return
	}
	url := "http://localhost:8080/static/post_file/" + fileName
	response.OkWithData(c, gin.H{"url": url})
}

func GetPostList(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		zlog.Warn("请求出错了")
		response.FailWithCode(c, response.INVALID_PARAMS, response.GetMsg(response.INVALID_PARAMS))
		return
	}
	pageSize := 10
	offset := (page - 1) * pageSize
	posts, err := controller.GetPostList(offset, pageSize)
	if err != nil {
		response.FailWithCode(c, response.INTERNAL_ERROR, response.GetMsg(response.INTERNAL_ERROR))
		return
	}
	if len(posts) == 0 {
		response.FailWithMessage(c, "该页没有对应数据")
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
	account := c.GetString("account")
	post, err := controller.GetPostDetail(account, postIdInt)
	if err != nil {
		response.FailWithCode(c, response.INTERNAL_ERROR, response.GetMsg(response.INTERNAL_ERROR))
		return
	}
	if post.PostID == 0 {
		response.FailWithMessage(c, "未找到该文章")
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
	if user.Account == "" {
		response.FailWithMessage(c, "未找到该用户")
		return
	}
	response.OkWithData(c, user)
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

func ToggleLike(c *gin.Context) {
	postId, err := strconv.ParseUint(c.Param("postId"), 10, 64)
	if err != nil {
		zlog.Error("转换失败")
		response.Fail(c)
	}
	postIdInt := uint(postId)
	account := c.GetString("account")
	isLike, count, err := controller.ToggleLike(postIdInt, account)
	if err != nil {
		response.FailWithCode(c, response.INTERNAL_ERROR, response.GetMsg(response.INTERNAL_ERROR))
		return
	}
	response.OkWithData(c, gin.H{"isLike": isLike, "count": count})
}

func GetFollowingPost(c *gin.Context) {
	account := c.GetString("account")
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		zlog.Warn("请求出错了")
		response.FailWithCode(c, response.INVALID_PARAMS, response.GetMsg(response.INVALID_PARAMS))
		return
	}
	pageSize := 10
	offset := (page - 1) * pageSize
	posts, err := feed.GetFollowingPosts(account, offset, pageSize)
	if err != nil {
		response.FailWithCode(c, response.INTERNAL_ERROR, response.GetMsg(response.INTERNAL_ERROR))
		return
	}
	if len(posts) == 0 {
		response.OkWithData(c, "你的暂时没有关注的对象哦")
		return
	}
	response.OkWithData(c, posts)
}

func GetHotRank(c *gin.Context) {
	var hotPosts []interface{}
	posts, score, err := controller.GetHotRank()
	if err != nil {
		response.FailWithCode(c, response.INTERNAL_ERROR, response.GetMsg(response.INTERNAL_ERROR))
		return
	}
	for _, post := range posts {
		if s, ok := score[post.PostID]; ok {
			hotPosts = append(hotPosts, gin.H{"post": post, "score": s})
		}
	}
	response.OkWithData(c, hotPosts)
}

func SearchPosts(c *gin.Context) {
	keyword := c.Query("keyword")
	if keyword == "" {
		response.FailWithMessage(c, "搜索关键词不能为空")
		return
	}
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		zlog.Warn("请求出错了")
		response.FailWithCode(c, response.INVALID_PARAMS, response.GetMsg(response.INVALID_PARAMS))
		return
	}
	pageSize := 10
	offset := (page - 1) * pageSize
	posts, err := controller.SearchPosts(keyword, offset, pageSize)
	if err != nil {
		response.FailWithCode(c, response.INTERNAL_ERROR, response.GetMsg(response.INTERNAL_ERROR))
		return
	}
	if len(posts) == 0 {
		response.OkWithData(c, "暂无相关内容")
		return
	}
	response.OkWithData(c, posts)
}

func SetPostPaid(c *gin.Context) {
	postId, err := strconv.ParseUint(c.Param("postId"), 10, 64)
	if err != nil {
		zlog.Error("转换失败")
		response.Fail(c)
	}
	postIdInt := uint(postId)
	role := c.MustGet("role").(int)
	flag, err := controller.SetPostPaid(role, postIdInt)
	if err != nil {
		response.FailWithCode(c, response.INTERNAL_ERROR, response.GetMsg(response.INTERNAL_ERROR))
		return
	}
	if !flag {
		response.FailWithMessage(c, "你暂无该权限")
		return
	}
	response.Ok(c)
}

func SetVip(c *gin.Context) {
	i, err := strconv.ParseUint(c.Param("Id"), 10, 64)
	if err != nil {
		zlog.Error("转换失败")
		response.Fail(c)
	}
	id := uint(i)
	role := c.MustGet("role").(int)
	flag, err := controller.PayVip(role, id)
	if err != nil {
		response.FailWithCode(c, response.INTERNAL_ERROR, response.GetMsg(response.INTERNAL_ERROR))
		return
	}
	if !flag {
		response.FailWithMessage(c, "暂无权限")
		return
	}
	response.Ok(c)
}

func AiSummary(c *gin.Context) {
	postId, err := strconv.ParseUint(c.Param("postId"), 10, 64)
	if err != nil {
		zlog.Error("转换失败")
		response.Fail(c)
		return
	}
	postIdInt := uint(postId)
	account := c.GetString("account")
	postSummary, err := controller.AiSummary(account, postIdInt)
	if err != nil {
		response.FailWithCode(c, response.INTERNAL_ERROR, response.GetMsg(response.INTERNAL_ERROR))
		return
	}
	if postSummary == "" {
		response.FailWithMessage(c, "请解锁会员使用")
		return
	}
	response.OkWithData(c, gin.H{"postSummary": postSummary})
}
