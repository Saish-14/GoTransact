package utils

import (
	// "fmt"
	// "gotransact/apps/Accounts/models"
	// "gotransact/logger"
	"context"
	"fmt"
	"gotransact/apps/accounts/models"
	"time"

	"aidanwoods.dev/go-paseto"
	"github.com/go-redis/redis/v8"
	"golang.org/x/crypto/bcrypt"
	// "github.com/sirupsen/logrus"
)

var (
	ctx = context.Background()
	rdb = redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Redis server address
	})
)

// implementing creation of token
var secretKey = paseto.NewV4AsymmetricSecretKey() // don't share this!!!
var publicKey = secretKey.Public()

func GeneratePasetoToken(user models.User) (string, error) {
	now := time.Now()
	exp := now.Add(24 * time.Hour)
	token := paseto.NewToken()
	token.SetIssuedAt(now)
	token.SetExpiration(exp)

	token.Set("User", user)
	signed := token.V4Sign(secretKey, nil)
	return signed, nil
}

func VerifyPasetoToken(signed string) (any, error) {
	val, err := rdb.Get(ctx, signed).Result()
	if err == nil && val == "Blacklisted" {
		return nil, fmt.Errorf("token has been revoked")
	}

	parser := paseto.NewParser()
	parser.AddRule(paseto.NotExpired())
	token, err := parser.ParseV4Public(publicKey, signed, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to parse token %w", err)
	}
	var user models.User
	err = token.Get("User", &user)
	if err != nil {
		return nil, fmt.Errorf("subject claim not found in token")
	}
	return user, nil
}



func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}