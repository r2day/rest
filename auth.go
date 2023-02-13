package rest

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

const (
	DefaultExpireLoginSessionTime = 24 // 默认 登陆有效时长24h
)

var (
	ExpireLoginSessionTime = DefaultExpireLoginSessionTime
	CustomLoginExpireTime = os.Getenv("CUSTOME_LOGIN_EXPIRE_TIME")
)

func init() {
	if CustomLoginExpireTime != "" {
		ExpireLoginSessionTime, _ = strconv.Atoi(CustomLoginExpireTime)
	}
}

// RenderLogin 返回登陆信息
func RenderLogin(c *gin.Context, accountId string, passwordFromDB []byte, passwordFromReq string, secretKey string, host string) {
	
	// 检查密码hash是否相同
	if err := bcrypt.CompareHashAndPassword(passwordFromDB, []byte(passwordFromReq)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 声明密码签名
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    accountId,
		ExpiresAt: time.Now().Add(time.Hour * time.Duration(ExpireLoginSessionTime)).Unix(), //1 day
	})

	// 签名密钥
	token, err := claims.SignedString([]byte(secretKey))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 设置缓存过期时间
	c.SetCookie("jwt", token, 3600 * ExpireLoginSessionTime, "/", host, false, false)

	// 
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "login success",
		"expire": ExpireLoginSessionTime,
		"expire_unit": "hours",
		"user": accountId,
	})

	// TODO 发生登陆消息到mq中记录登陆信息
}