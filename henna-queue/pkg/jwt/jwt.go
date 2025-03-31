package jwt

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
)

// UserClaims 用户JWT声明
type UserClaims struct {
	UserID string `json:"user_id"`
	jwt.StandardClaims
}

// AdminClaims 管理员JWT声明
type AdminClaims struct {
	AdminID uint  `json:"admin_id"`
	ShopID  uint  `json:"shop_id"`
	Role    int8  `json:"role"`
	jwt.StandardClaims
}

// GenerateUserToken 生成用户token
func GenerateUserToken(userID string) (string, error) {
	// 获取配置
	secret := viper.GetString("app.jwt_secret")
	expiration := viper.GetInt("app.jwt_expiration")
	
	// 创建声明
	claims := UserClaims{
		userID,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(expiration) * time.Second).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}
	
	// 生成token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// GenerateAdminToken 生成管理员token
func GenerateAdminToken(adminID, shopID uint, role int8) (string, error) {
	// 获取配置
	secret := viper.GetString("app.jwt_secret")
	expiration := viper.GetInt("app.jwt_expiration")
	
	// 创建声明
	claims := AdminClaims{
		adminID,
		shopID,
		role,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(expiration) * time.Second).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}
	
	// 生成token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ParseUserToken 解析用户token
func ParseUserToken(tokenString string) (*UserClaims, error) {
	// 获取秘钥
	secret := viper.GetString("app.jwt_secret")
	
	// 解析token
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	
	if err != nil {
		return nil, err
	}
	
	// 验证并返回声明
	if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		return claims, nil
	}
	
	return nil, errors.New("无效的token")
}

// ParseAdminToken 解析管理员token
func ParseAdminToken(tokenString string) (*AdminClaims, error) {
	// 获取秘钥
	secret := viper.GetString("app.jwt_secret")
	
	// 解析token
	token, err := jwt.ParseWithClaims(tokenString, &AdminClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	
	if err != nil {
		return nil, err
	}
	
	// 验证并返回声明
	if claims, ok := token.Claims.(*AdminClaims); ok && token.Valid {
		return claims, nil
	}
	
	return nil, errors.New("无效的token")
} 