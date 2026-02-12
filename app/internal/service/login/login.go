package login

import (
	"commmunity/app/internal/db/global"
	"commmunity/app/internal/model"
	"commmunity/app/utils"
	"commmunity/app/zlog"
	"time"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

func Register(account string, password string, name string) (error, bool, bool) { //第一个bool判断信息是否完善，第二个bool判断用户是否可使用该账户
	if account == "" || password == "" || name == "" {
		return nil, false, false
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		zlog.Error("密码哈希失败", zap.String("password", password), zap.Error(err))
		return err, false, false
	}
	users, err := global.User.GetUser(account)
	if users != nil || err != nil {
		zlog.Warn("用户已存在", zap.String("account", account))
		return nil, true, false
	}
	err = global.User.CreateUser(account, string(hash), name)
	if err != nil {
		zlog.Error("用户创建失败", zap.String("account", account), zap.String("name", name), zap.Error(err))
		return err, false, false
	}
	return nil, true, true
}

// 密码最多输入5次

func Login(account string, password string) (error, bool, bool) { //第一个bool判断传入信息是否合格，第二个bool判断账户和密码是否对
	if account == "" || password == "" {
		return nil, false, false
	}
	user, err := global.User.GetUser(account)
	if err != nil {
		return err, false, false
	}
	if user == nil {
		return nil, true, false
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Hash), []byte(password))
	if err != nil {
		zlog.Warn("密码错误", zap.String("password", password))
		return nil, true, false
	}
	return nil, true, true
}

type UserProfileDTO struct {
	Account      string        `json:"account"`
	Name         string        `json:"name"`
	Introduction string        `json:"introduction"`
	Avatar       string        `json:"avatar"`
	Role         int           `json:"role"`
	IsMuted      bool          `json:"isMuted"`
	Posts        []UserPostDTO `json:"posts"`
}
type UserPostDTO struct {
	PostID       uint   `json:"post_id"`
	Title        string `json:"title"`
	ViewCount    uint   `json:"view_count"`
	LikeCount    uint   `json:"like_count"`
	CommentCount uint   `json:"comment_count"`
}

func GetProfile(account string) (UserProfileDTO, error) {
	user, err := global.User.GetProfile(account)
	if err != nil {
		return UserProfileDTO{}, err
	}
	userPostDTO := make([]UserPostDTO, 0, len(user.Posts))
	for _, post := range user.Posts {
		userPostDTO = append(userPostDTO, UserPostDTO{
			PostID:       post.ID,
			Title:        post.Title,
			ViewCount:    post.ViewCount,
			LikeCount:    post.LikeCount,
			CommentCount: post.CommentCount,
		})
	}
	userProfile := UserProfileDTO{
		Account:      user.Account,
		Name:         user.UserProfile.Name,
		Introduction: user.UserProfile.Introduction,
		Avatar:       user.UserProfile.Avatar,
		Role:         user.Role,
		IsMuted:      user.UserProfile.IsMuted,
		Posts:        userPostDTO,
	}
	return userProfile, nil
}

func Logout(token string, refreshToken string) error {
	claim, err := utils.ParseToken(token)
	if err != nil {
		return err
	}
	timeLeft := claim.StandardClaims.ExpiresAt - time.Now().Unix()
	if timeLeft <= 0 {
		return nil
	}
	err = global.UserRedis.AddToBlacklist(token, time.Duration(timeLeft)*time.Second)
	if err != nil {
		zlog.Error("token拉入黑名单失败", zap.String("token", token))
		return err
	}
	newClaim, err := utils.ParseRefreshToken(refreshToken)
	if err != nil {
		return err
	}
	timeLeft1 := newClaim.StandardClaims.ExpiresAt - time.Now().Unix()
	if timeLeft1 <= 0 {
		return nil
	}
	err = global.UserRedis.AddToBlacklist(refreshToken, time.Duration(timeLeft1)*time.Second)
	if err != nil {
		zlog.Error("refreshToken拉入黑名单失败", zap.String("refreshToken", refreshToken))
		return err
	}
	return nil
}

func DeleteUser(account string) error {
	err := global.User.DeleteUser(account)
	if err != nil {
		zlog.Error("用户删除失败", zap.String("account", account))
		return err
	}
	return global.User.DeleteUser(account)
}

// 短时间内应只修改一次，旧密码验证,改密踢人

func ChangePassword(account string, firstPassword string, secondPassword string) (error, bool, bool) { //第一个bool检查传入信息是否合格，第二个检查两次密码是否相同
	if firstPassword == "" || secondPassword == "" {
		return nil, false, false
	}
	if firstPassword != secondPassword {
		return nil, true, false
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(firstPassword), bcrypt.DefaultCost)
	if err != nil {
		zlog.Error("修改密码哈希失败", zap.Error(err))
		return err, false, false
	}
	err = global.User.ChangePassword(account, string(hash))
	if err != nil {
		return err, false, false
	}
	return nil, true, true
}

// 一天仅能修改5次

func ChangeName(account string, newName string) (error, bool) {
	if newName == "" {
		zlog.Warn("修改用户名不能为空")
		return nil, false
	}
	err := global.User.ChangeUserName(account, newName)
	if err != nil {
		return err, false
	}
	return nil, true
}

func ChangeAvatar(account string, avatar string) error {
	return global.User.ChangeAvatar(account, avatar)
}

func ChangeIntroduction(account string, introduction string) error {
	return global.User.ChangeIntroduction(account, introduction)
}

func GetUserRole(account string) (int, error) {
	user, err := global.User.GetUser(account)
	if user == nil || err != nil {
		return 0, err
	}
	role := user.Role
	return role, nil
}

func IsTokenValid(tokenStr string) bool {
	if global.UserRedis.IsInBlacklist(tokenStr) {
		return false
	}
	return true
}

func Muted(userId uint, role int, isMuted bool) (error, bool) {
	if role == model.RoleAdmin {
		err := global.User.Muted(userId, isMuted)
		if err != nil {
			return err, false
		}
		return nil, true
	}
	return nil, false
}
