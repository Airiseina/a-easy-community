package utils

import (
	"commmunity/app/zlog"
	"errors"
	"math/rand"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var jwtKey = viper.GetString("jwtKey")
var jwtRefreshKey = viper.GetString("jwtRefreshKey")

type MyClaims struct {
	Account string `json:"account"`
	Role    int    `json:"role"`
	Type    string `json:"type"`
	jwt.StandardClaims
}

func MakeToken(account string, role int) (string, string, error) {
	claim := &MyClaims{
		Account: account,
		Role:    role,
		Type:    "access",
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Add(-5 * time.Second).Unix(),
			ExpiresAt: time.Now().Add(15 * time.Minute).Unix(),
			Issuer:    "Tom",
		},
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	token, err := t.SignedString([]byte(jwtKey))
	if err != nil {
		zlog.Error("token生成失败", zap.Error(err))
		return "", "", err
	}
	newClaim := &MyClaims{
		Account: account,
		Role:    role,
		Type:    "refresh",
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Add(-5 * time.Second).Unix(),
			ExpiresAt: time.Now().Add(7 * 24 * time.Hour).Unix(),
			Issuer:    "Tom",
		},
	}
	rT := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaim)
	refreshToken, err := rT.SignedString([]byte(jwtRefreshKey))
	if err != nil {
		zlog.Error("refreshToken生成失败", zap.Error(err))
		return "", "", err
	}
	return token, refreshToken, nil
}

func ParseToken(tokenString string) (*MyClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtKey), nil
	})
	if err != nil {
		zlog.Error("无效token", zap.Error(err))
		return nil, err
	}
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		return claims, nil
	}
	zlog.Error("token已过期")
	return nil, errors.New("token已过期")
}

func ParseRefreshToken(tokenString string) (*MyClaims, error) {
	refreshToken, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtRefreshKey), nil
	})
	if err != nil {
		zlog.Error("无效refresh", zap.Error(err))
		return nil, err
	}
	if claims, ok := refreshToken.Claims.(*MyClaims); ok && refreshToken.Valid {
		return claims, nil
	}
	zlog.Error("refreshToken已过期")
	return nil, errors.New("refreshToken已过期")
}

func RandomDuration(rt int) time.Duration {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	randTime := time.Duration(r.Intn(rt)) * time.Minute
	return randTime
}

func TruncateContent(content string, maxLines int, maxChars int) string {
	lines := strings.Split(content, "\n")
	var resultBuilder strings.Builder
	charCount := 0
	for i, line := range lines {
		if i >= maxLines {
			break
		}
		if strings.Contains(line, "![") || strings.Contains(line, "<img") {
			continue
		}
		lineLen := utf8.RuneCountInString(line)
		if charCount+lineLen > maxChars {
			break
		}
		resultBuilder.WriteString(line)
		resultBuilder.WriteString("\n")
		charCount += lineLen
	}
	resultBuilder.WriteString("\n\n> 🔒 **剩余内容为付费会员专享，请升级后查看...**")
	return resultBuilder.String()
}
