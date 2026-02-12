package routes

import (
	"commmunity/app/internal/api"
	"commmunity/app/internal/cron"
	"commmunity/app/middleware"
	"context"
	"time"

	"github.com/gin-gonic/gin"
)

func Routes() {
	cronLikeManager := cron.NewCronManager(1 * time.Minute)
	cronLikeManager.Start(context.Background(), cron.SyncPostLikes)
	cronViewManager := cron.NewCronManager(5 * time.Second)
	cronViewManager.Start(context.Background(), cron.SyncView)
	cronHotRankManager := cron.NewCronManager(5 * time.Second)
	cronHotRankManager.Start(context.Background(), cron.RefreshHot)
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
		protected.POST("/password-change", middleware.RateLimitingMiddleware("changePassword", 5*time.Second, 1), api.ChangePassword)
		// 修改头像
		protected.POST("/avatar", api.ChangeAvatar)
		// 修改简介
		protected.PATCH("/introduction", api.ChangeIntroduction)
		// 注销账户
		protected.DELETE("/delete-user", api.DeleteUser)
		//	退出登录
		protected.POST("/logout", api.Logout)
		//发布文章
		protected.POST("/posts", middleware.RateLimitingMiddleware("createPost", 5*time.Second, 1), api.CreatePost)
		//论坛主页
		protected.GET("/posts", api.GetPostList)
		//文章主页
		protected.GET("/posts/:postId", api.GetPostDetail)
		//发表评论
		protected.POST("/posts/:postId", middleware.RateLimitingMiddleware("createComment", 3*time.Second, 1), api.CreateComment)
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
		//点赞
		protected.POST("posts/:postId/like", middleware.RateLimitingMiddleware("like", 5*time.Second, 2), api.ToggleLike)
		//关注
		protected.POST("/follow/:Id", api.Follow)
		//关注列表
		protected.GET("/following", api.GetFollowings)
		//粉丝列表
		protected.GET("/follow", api.GetFollowers)
		//关注动态
		protected.GET("/following_post", api.GetFollowingPost)
		//热度榜
		protected.GET("/hot_rank", api.GetHotRank)
	}
	r.Run(":8080")
}
