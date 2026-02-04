package api

import (
	"commmunity/app/internal/model"
	"commmunity/app/internal/response"
	"commmunity/app/internal/service/login"
	"commmunity/app/utils"
	"commmunity/app/zlog"
	"fmt"
	"net/http"
	"path/filepath"
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

// 增加refresh token
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
	info, err, flag1 := login.GetProfile(account)
	if err != nil {
		response.Fail(c)
		return
	}
	if flag1 == false {
		response.FailWithCode(c, response.ERROR_USER_NOT_EXIST_OR_PASSWORD_WRONG, response.GetMsg(response.ERROR_USER_NOT_EXIST_OR_PASSWORD_WRONG))
		return
	}
	userInfo := map[string]interface{}{
		"Account":      (*info).Account,
		"Name":         (*info).Name,
		"Introduction": (*info).Introduction,
		"Avatar":       (*info).Avatar,
		"Role":         (*info).Role,
		"IsMuted":      (*info).IsMuted,
	}
	response.OkWithData(c, userInfo)
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

func UploadAvatar(c *gin.Context) {
	// 1. 找文件
	// "avatar" 是我们和前端约好的暗号。前端传文件时，字段名必须叫 avatar。
	file, err := c.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "大哥，你没传图片啊！"})
		return
	}

	// 2. 也是为了安全：重命名文件
	// 为什么？如果两个用户都上传了 "me.jpg"，后传的会把先传的覆盖掉！
	// 所以我们要把文件名改成唯一的。
	// 策略：用户ID + 当前时间戳 + 原来的后缀名

	// 获取文件后缀，比如 ".jpg" 或 ".png"
	ext := filepath.Ext(file.Filename)

	// 假设从中间件拿到了当前用户ID (c.Get("userID"))，这里先假设是 10086
	userID := 10086

	// 生成新名字：10086_1762345678.jpg
	newFileName := fmt.Sprintf("%d_%d%s", userID, time.Now().Unix(), ext)

	// 3. 拼凑保存路径
	// 告诉 Go：存到当前目录下的 uploads/avatars 文件夹里
	savePath := filepath.Join("uploads", "avatars", newFileName)

	// 4. 执行保存动作 (存硬盘)
	// 这一步才是真正把二进制流写进硬盘
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器硬盘坏了，存不进去"})
		return
	}

	// 5. 生成访问链接 (存数据库用)
	// 这一步很重要！我们不能把硬盘路径 "uploads/avatars/..." 告���前端。
	// 我们要给一个 URL，比如 "/static/avatars/..."
	// 为什么是 /static？这是我们下一步要配置的“传送门”。
	accessURL := "/static/avatars/" + newFileName

	// 6. 更新数据库 (这里写伪代码)
	// database.UpdateUserAvatar(userID, accessURL)

	// 7. 告诉前端成功了，顺便把新头像地址给他
	c.JSON(http.StatusOK, gin.H{
		"msg": "上传成功",
		"url": accessURL,
	})
}
