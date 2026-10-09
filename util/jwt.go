package util

import (
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func jwtSecret() []byte {
	return []byte(os.Getenv("JWT_SECRET"))
}

// 密码加密
func HashPassword(pwd string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	return string(b), err
}

// 校验密码
func CheckPassword(hashPwd, pwd string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashPwd), []byte(pwd)) == nil
}

// 生成JWT
func GenToken(userId uint) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userId,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret())
}

// 解析token
func ParseToken(tokenStr string) (uint, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret(), nil
	})
	if err != nil {
		return 0, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return 0, errors.New("token invalid")
	}
	uid, ok := claims["user_id"].(float64)
	if !ok {
		return 0, errors.New("user id invalid")
	}
	return uint(uid), nil
}

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := ""
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			parts := strings.SplitN(authHeader, " ", 2)
			if !(len(parts) == 2 && parts[0] == "Bearer") {
				c.JSON(http.StatusUnauthorized, gin.H{"msg": "token格式错误"})
				c.Abort()
				return
			}
			tokenStr = parts[1]
		} else if q := c.Query("token"); q != "" {
			// 供 <a href> 直接下载等无法自定义 Header 的场景
			tokenStr = q
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "未登录"})
			c.Abort()
			return
		}
		uid, err := ParseToken(tokenStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "token失效"})
			c.Abort()
			return
		}
		// 把 userId 存到上下文，后续 controller 拿
		c.Set("user_id", uid)
		c.Next()
	}
}
