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
	user := model.User{
		Account: account,
		Hash:    string(hash),
		Name:    name,
	}
	users, err := global.User.GetUser(account)
	if users != nil || err != nil {
		zlog.Warn("用户已存在", zap.String("account", account))
		return nil, true, false
	}
	err = global.User.CreateUser(user)
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

func GetProfile(account string) (*model.User, error, bool) {
	user, err := global.User.GetUser(account)
	if err != nil {
		return nil, err, false
	}
	if user == nil {
		return nil, nil, false
	}
	return user, nil, true
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
	err = global.Redis.AddToBlacklist(token, time.Duration(timeLeft)*time.Second)
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
	err = global.Redis.AddToBlacklist(refreshToken, time.Duration(timeLeft1)*time.Second)
	if err != nil {
		zlog.Error("refreshToken拉入黑名单失败", zap.String("refreshToken", refreshToken))
		return err
	}
	return nil
}

// 注销后的用户应立即消除登录jwt

func DeleteUser(account string) error {
	err := global.User.DeleteUser(account)
	if err != nil {
		return err
	}
	return nil
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
	user := model.User{
		Account: account,
		Hash:    string(hash),
	}
	err = global.User.ChangePassword(user)
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
	user := model.User{
		Account: account,
		Name:    newName,
	}
	err := global.User.ChangeUserName(user)
	if err != nil {
		return err, false
	}
	return nil, true
}

func ChangeAvatar(account string, avatar string) error {
	user := model.User{
		Account: account,
		Avatar:  avatar,
	}
	err := global.User.ChangeAvatar(user)
	if err != nil {
		return err
	}
	return nil
}

func ChangeIntroduction(account string, introduction string) error {
	user := model.User{
		Account:      account,
		Introduction: introduction,
	}
	err := global.User.ChangeIntroduction(user)
	if err != nil {
		return err
	}
	return nil
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
	if global.Redis.IsInBlacklist(tokenStr) {
		return false
	}
	return true
}
