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
	cronViewManager := cron.NewCronManager(5 * time.Minute)
	cronViewManager.Start(context.Background(), cron.SyncView)
	cronHotRankManager := cron.NewCronManager(5 * time.Hour)
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
		protected.GET("/profile", api.GetProfile)                                                                                     // 获取个人信息
		protected.PATCH("/username", api.ChangeUserName)                                                                              // 修改用户名
		protected.POST("/avatar", api.ChangeAvatar)                                                                                   // 修改头像
		protected.PATCH("/introduction", api.ChangeIntroduction)                                                                      // 修改简介
		protected.POST("/logout", api.Logout)                                                                                         // 退出登录
		protected.DELETE("/delete-user", api.DeleteUser)                                                                              // 注销账户
		protected.POST("/password-change", middleware.RateLimitingMiddleware("changePassword", 5*time.Second, 1), api.ChangePassword) // 修改密码
		protected.GET("/users/:Id", api.GetUserProfile)                                                                               // 查看指定用户主页
		protected.POST("/muted/:Id", api.Muted)                                                                                       // 禁言用户
		protected.POST("vip/:Id", api.SetVip)                                                                                         //设置vip用户
	}
	{
		protected.GET("/posts", api.GetPostList)                                                                    // 论坛主页/帖子列表
		protected.POST("/posts/:postId/summary", api.AiSummary)                                                     //总结文章（VIP专属）
		protected.GET("/posts/:postId", api.GetPostDetail)                                                          // 文章详情
		protected.GET("/search", api.SearchPosts)                                                                   // 搜索帖子
		protected.POST("/upload", api.UploadImage)                                                                  // 上传文章图片
		protected.GET("/hot_rank", api.GetHotRank)                                                                  // 热度榜单
		protected.POST("/posts", middleware.RateLimitingMiddleware("createPost", 5*time.Second, 1), api.CreatePost) // 发布帖子
		protected.DELETE("/posts/:postId", api.DeletePost)                                                          // 删除帖子
		protected.POST("/paid-post/:postId", api.SetPostPaid)                                                       //设置需要花费文章
	}
	{
		protected.POST("/posts/:postId", middleware.RateLimitingMiddleware("createComment", 3*time.Second, 1), api.CreateComment) // 发表评论
		protected.DELETE("/posts/:postId/:posterId/:commentId", api.DeleteComment)                                                // 删除评论
	}
	{
		protected.POST("posts/:postId/like", middleware.RateLimitingMiddleware("like", 5*time.Second, 2), api.ToggleLike) //点赞
		protected.POST("/follow/:Id", api.Follow)                                                                         // 关注/取消关注用户
		protected.GET("/following", api.GetFollowings)                                                                    // 我的关注列表
		protected.GET("/follow", api.GetFollowers)                                                                        // 我的粉丝列表
		protected.GET("/following_post", api.GetFollowingPost)                                                            // 关注人的动态
	}

	r.Run(":8080")
}
