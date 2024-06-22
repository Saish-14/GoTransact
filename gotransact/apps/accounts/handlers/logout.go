package handlers

import (
	"context"
	logger "gotransact/log"
	"net/http"
	"time"

	account_utils "gotransact/apps/accounts/utils"
	base_utils "gotransact/apps/base/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

var (
	ctx = context.Background()
	rdb = redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Redis server address
	})
)

func LogoutHandler(c *gin.Context) {
	logger.InfoLogger.WithFields(logrus.Fields{
		"method": c.Request.Method,
		"url":    c.Request.URL.String(),
	}).Info("attempted Logout")
	authHeader := c.GetHeader("Authorization")

	status, message, data := Logout(authHeader)

	c.JSON(status, base_utils.Response{
		Status:  status,
		Message: message,
		Data:    data,
	})
}

func Logout(authHeader string) (int, string, map[string]interface{}) {
	logger.InfoLogger.WithFields(logrus.Fields{
		"message": "in logout func",
	}).Info("attempted logout method")

	if authHeader == "" {
		return http.StatusUnauthorized, "authorization header missing", map[string]interface{}{}
	}

	//tokenStr := authHeader[len("Bearer "):]

	_, err := account_utils.VerifyPasetoToken(authHeader)
	if err != nil {
		return http.StatusUnauthorized, "invalid token", map[string]interface{}{}
	}

	// Blacklist the token by storing it in Redis with an expiration time
	err = rdb.Set(ctx, authHeader, "Blacklisted", 24*time.Hour).Err() // adjust expiration time as needed
	if err != nil {
		return http.StatusInternalServerError, "failed to blacklist token", map[string]interface{}{}
	}
	logger.InfoLogger.WithFields(logrus.Fields{"logout": "success"}).Info("Logged out successfully")
	return http.StatusOK, "logged out successfully", map[string]interface{}{}
}