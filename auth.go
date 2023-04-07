package rest

import (
	b64 "encoding/base64"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

const (
	DefaultExpireLoginSessionTime = 24 // 默认 登陆有效时长24h
)

var (
	ExpireLoginSessionTime = DefaultExpireLoginSessionTime
	CustomLoginExpireTime  = os.Getenv("CUSTOME_LOGIN_EXPIRE_TIME")
)

func init() {
	if CustomLoginExpireTime != "" {
		ExpireLoginSessionTime, _ = strconv.Atoi(CustomLoginExpireTime)
	}
}

// RenderLogin 返回登陆信息
func RenderLogin(c *gin.Context, accountId string, passwordFromDB []byte, passwordFromReq string, secretKey string, host string) {

	logCtx := log.WithField("account_id", accountId)
	// 检查密码hash是否相同
	if err := bcrypt.CompareHashAndPassword(passwordFromDB, []byte(passwordFromReq)); err != nil {
		// c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": "code is no correct"})
		c.JSON(http.StatusBadRequest, gin.H{
			"message":    "the password or account isn't correct",
			"account_id": accountId,
			"status":     "error",
			"title":      "An error occurred.",
		})
		logCtx.Error(err)
		return
	}

	// 声明密码签名
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    accountId,
		ExpiresAt: time.Now().Add(time.Hour * time.Duration(ExpireLoginSessionTime)).Unix(), //1 day
	})

	sDec, err := b64.StdEncoding.DecodeString(secretKey)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message":    "the sign isn't correct",
			"account_id": accountId,
			"status":     "error",
			"title":      "An error occurred.",
		})
		logCtx.Error(err)
		return
	}
	// 签名密钥
	token, err := claims.SignedString(sDec)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message":    "the encode to base64 isn't correct",
			"account_id": accountId,
			"status":     "error",
			"title":      "An error occurred.",
			"error":      err.Error(),
		})
		logCtx.Error(err)
		return
	}

	// 设置缓存过期时间
	c.SetCookie("jwt", token, 3600*ExpireLoginSessionTime, "/", host, false, false)

	c.JSON(http.StatusCreated, gin.H{
		"message":    "sign in success",
		"account_id": accountId,
		"status":     "success",
		"title":      "Sign In.",
	})
	logCtx.Debug("password check done")
	// TODO 发生登陆消息到mq中记录登陆信息
}
