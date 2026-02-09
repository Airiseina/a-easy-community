package api

import (
	"commmunity/app/internal/model"
	"commmunity/app/internal/response"
	"commmunity/app/internal/service/login"
	"commmunity/app/utils"
	"commmunity/app/zlog"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Register(c *gin.Context) {
	var user model.UserRequest
	if err := c.ShouldBindJSON(&user); err != nil {
		zlog.Warn("请求出错了", zap.Error(err))
		response.FailWithCode(c, response.INVALID_PARAMS, response.GetMsg(response.INVALID_PARAMS))
		return
	}
	err, flag1, flag2 := login.Register(user.Account, user.Password, user.Name)
	if err != nil {
		response.FailWithCode(c, response.INTERNAL_ERROR, response.GetMsg(response.INTERNAL_ERROR))
		return
	}
	if flag1 == false {
		response.FailWithMessage(c, "请完善信息")
		return
	}
	if flag2 == false {
		response.FailWithCode(c, response.ERROR_USERNAME_EXIST, response.GetMsg(response.ERROR_USERNAME_EXIST))
		return
	}
	response.Ok(c)
}

func Login(c *gin.Context) {
	var user model.UserRequest
	if err := c.ShouldBindJSON(&user); err != nil {
		zlog.Warn("请求出错了")
		response.FailWithCode(c, response.INVALID_PARAMS, response.GetMsg(response.INVALID_PARAMS))
		return
	}
	err, flag1, flag2 := login.Login(user.Account, user.Password)
	if err != nil {
		response.FailWithCode(c, response.INTERNAL_ERROR, response.GetMsg(response.INTERNAL_ERROR))
		return
	}
	if flag1 == false {
		response.FailWithMessage(c, "请完善信息")
		return
	}
	if flag2 == false {
		response.FailWithCode(c, response.ERROR_USER_NOT_EXIST_OR_PASSWORD_WRONG, response.GetMsg(response.ERROR_USER_NOT_EXIST_OR_PASSWORD_WRONG))
		return
	}
	role, err := login.GetUserRole(user.Account)
	if err != nil {
		response.FailWithCode(c, response.INTERNAL_ERROR, response.GetMsg(response.INTERNAL_ERROR))
		return
	}
	token, refreshToken, err := utils.MakeToken(user.Account, role)
	if err != nil {
		response.FailWithCode(c, response.INTERNAL_ERROR, response.GetMsg(response.INTERNAL_ERROR))
		return
	}
	c.SetCookie(
		"refresh_token", // Cookie 名字
		refreshToken,    // 值
		3600*24*7,       // 过期时间 (秒)，这里设为 7 天
		"/",             // Path
		"localhost",     // Domain (上线换成你的域名)
		false,           // Secure: 本地开发 false (HTTP)，上线必须 true (HTTPS)
		true,            // HttpOnly: 【关键】开启！禁止 JS 读取
	)
	response.OkWithData(c, gin.H{"token": token})
}

func GetProfile(c *gin.Context) {
	account := c.GetString("account")
	info, err := login.GetProfile(account)
	if err != nil {
		response.Fail(c)
		return
	}
	response.OkWithData(c, info)
}

func Logout(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		response.FailWithMessage(c, "token格式不对")
		c.Abort()
		return
	}
	token := parts[1]
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		response.FailWithMessage(c, "请重新登录")
		return
	}
	fmt.Println(token)
	err = login.Logout(token, refreshToken)
	if err != nil {
		response.Fail(c)
		return
	}
	c.SetCookie(
		"refresh_token", // 名字必须一样
		"",              // 值设为空
		-1,              // MaxAge < 0 表示删除
		"/",             // Path 必须一样
		"localhost",     // Domain 必须一样
		false,           // Secure
		true,            // HttpOnly
	)
	response.Ok(c)
}

func DeleteUser(c *gin.Context) {
	account := c.GetString("account")
	err := login.DeleteUser(account)
	if err != nil {
		response.FailWithCode(c, response.INTERNAL_ERROR, response.GetMsg(response.INTERNAL_ERROR))
		return
	}
	response.Ok(c)
}

func ChangePassword(c *gin.Context) {
	var user model.UserPassword
	if err := c.ShouldBindJSON(&user); err != nil {
		zlog.Warn("请求出错了")
		response.FailWithCode(c, response.INVALID_PARAMS, response.GetMsg(response.INVALID_PARAMS))
		return
	}
	user.Account = c.GetString("account")
	err, flag1, flag2 := login.ChangePassword(user.Account, user.FirstPassWord, user.SecondPassWord)
	if err != nil {
		response.FailWithCode(c, response.INTERNAL_ERROR, response.GetMsg(response.INTERNAL_ERROR))
		return
	}
	if flag1 == false {
		response.FailWithMessage(c, "请完善输入的两次密码")
		return
	}
	if flag2 == false {
		response.FailWithMessage(c, "请确保前后密码一致")
		return
	}
	response.Ok(c)
}

func ChangeUserName(c *gin.Context) {
	var user model.UserRequest
	if err := c.ShouldBindJSON(&user); err != nil {
		zlog.Warn("请求出错了")
		response.FailWithCode(c, response.INVALID_PARAMS, response.GetMsg(response.INVALID_PARAMS))
		return
	}
	user.Account = c.GetString("account")
	err, flag1 := login.ChangeName(user.Account, user.Name)
	if err != nil {
		response.FailWithCode(c, response.INTERNAL_ERROR, response.GetMsg(response.INTERNAL_ERROR))
		return
	}
	if flag1 == false {
		response.FailWithMessage(c, "修改昵称不可为空")
		return
	}
	response.Ok(c)
}

// 将储存如硬盘的步骤移到serve层
func ChangeAvatar(c *gin.Context) {
	file, err := c.FormFile("avatar")
	if err != nil {
		zlog.Warn("请求出错了")
		response.FailWithCode(c, response.INVALID_PARAMS, response.GetMsg(response.INVALID_PARAMS))
		return
	}
	ext := filepath.Ext(file.Filename)
	account := c.GetString("account")
	newFileName := fmt.Sprintf("%s_%d%s", account, time.Now().Unix(), ext)
	err = os.MkdirAll("./uploads/avatars", 0755)
	if err != nil {
		zlog.Error("文件夹创建失败", zap.Error(err))
		response.FailWithCode(c, response.INTERNAL_ERROR, response.GetMsg(response.INTERNAL_ERROR))
		return
	}
	savePath := filepath.Join("uploads", "avatars", newFileName)
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		zlog.Error("服务器硬盘出错", zap.Error(err))
		response.FailWithCode(c, response.INTERNAL_ERROR, response.GetMsg(response.INTERNAL_ERROR))
		return
	}
	accessURL := "/static/avatars/" + newFileName
	err = login.ChangeAvatar(account, accessURL)
	if err != nil {
		response.FailWithMessage(c, "头像上传失败")
		return
	}
	response.OkWithData(c, gin.H{"avatar": accessURL})
}

func ChangeIntroduction(c *gin.Context) {
	var user model.UserInfoRequest
	if err := c.ShouldBindJSON(&user); err != nil {
		zlog.Error("请求出错了")
		response.FailWithCode(c, response.INVALID_PARAMS, response.GetMsg(response.INVALID_PARAMS))
		return
	}
	user.Account = c.GetString("account")
	err := login.ChangeIntroduction(user.Account, user.Introduction)
	if err != nil {
		response.FailWithCode(c, response.INTERNAL_ERROR, response.GetMsg(response.INTERNAL_ERROR))
		return
	}
	response.Ok(c)
}

func Muted(c *gin.Context) {
	var user model.UserInfoRequest
	if err := c.ShouldBindJSON(&user); err != nil {
		zlog.Warn("请求出错了")
		response.FailWithCode(c, response.INVALID_PARAMS, response.GetMsg(response.INVALID_PARAMS))
		return
	}
	role := c.MustGet("role").(int)
	i, err := strconv.ParseUint(c.Param("Id"), 10, 64)
	if err != nil {
		zlog.Error("转换失败")
		response.Fail(c)
	}
	id := uint(i)
	err, flag := login.Muted(id, role, user.IsMuted)
	if err != nil {
		response.FailWithCode(c, response.INTERNAL_ERROR, response.GetMsg(response.INTERNAL_ERROR))
		return
	}
	if flag == false {
		response.FailWithMessage(c, "权限不够")
		return
	}
	response.Ok(c)
}
func RefreshToken(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		response.FailWithMessage(c, "请重新登录")
		return
	}
	claims, err := utils.ParseRefreshToken(refreshToken)
	if err != nil {
		response.FailWithMessage(c, "登录已彻底过期，请重新登录")
		return
	}

	// 3. (可选) 检查 Redis 黑名单
	// 如果 Refresh Token 也在黑名单（用户注销了），那也不能换

	// 4. 签发新的 Access Token
	newAccessToken, _, err := utils.MakeToken(claims.Account, claims.Role)
	// 5. (可选) 甚至可以把 Refresh Token 也顺便换个新的 (Token Rotation 策略)
	response.OkWithData(c, gin.H{"access_token": newAccessToken})
}
