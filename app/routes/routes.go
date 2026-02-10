package routes

import (
	"commmunity/app/internal/api"
	"commmunity/app/middleware"

	"github.com/gin-gonic/gin"
)

func Routes() {
	r := gin.Default()
	r.Use(middleware.CorsMiddleWare())
	r.Static("/static", "./uploads")
	account := r.Group("/account")
	{
		account.POST("/register", api.Register)
		account.POST("/login", api.Login)
		account.POST("/refresh", api.RefreshToken)
	}
	protected := r.Group("/account/protected")
	protected.Use(middleware.JwtAuthMiddleware())
	{
		// 显示信息
		protected.GET("/profile", api.GetProfile)
		// 修改用户名
		protected.PATCH("/username", api.ChangeUserName)
		// 修改密码
		protected.POST("/password-change", api.ChangePassword)
		// 修改头像
		protected.POST("/avatar", api.ChangeAvatar)
		// 修改简介
		protected.PATCH("/introduction", api.ChangeIntroduction)
		// 注销账户
		protected.DELETE("/delete-user", api.DeleteUser)
		//	退出登录
		protected.POST("/logout", api.Logout)
		//发布文章
		protected.POST("/posts", api.CreatePost)
		//论坛主页
		protected.GET("/posts", api.GetPostList)
		//文章主页
		protected.GET("/posts/:postId", api.GetPostDetail)
		//发表评论
		protected.POST("/posts/:postId", api.CreateComment)
		//查找作者主页
		protected.GET("/users/:Id", api.GetUserProfile)
		//删除文章
		protected.DELETE("/posts/:postId", api.DeletePost)
		//删除评论
		protected.DELETE("/posts/:postId/:posterId/:commentId", api.DeleteComment)
		//禁言
		protected.POST("/muted/:Id", api.Muted)
		//为文章添加图片
		protected.POST("/upload", api.UploadImage)
	}
	r.Run(":8080")
}
