package utils

import (
	"strings"
	"time"

	"github.com/Poloni84Learning/ebook-store/config"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims định nghĩa custom claims structure
type JWTClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateJWT tạo JWT token với thông tin user
func GenerateJWT(userID uint, username string, role string) (string, error) {
	cfg := config.LoadConfig()

	claims := JWTClaims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(cfg.JWTExpiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "ebook-store-api",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.JWTSecret))
}

// ParseJWT phân tích và validate JWT token
func ParseJWT(tokenString string) (*JWTClaims, error) {
	cfg := config.LoadConfig()

	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(cfg.JWTSecret), nil
	})

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, err
}

// RefreshJWT làm mới token khi gần hết hạn
func RefreshJWT(tokenString string) (string, error) {
	claims, err := ParseJWT(tokenString)
	if err != nil {
		return "", err
	}

	// Chỉ refresh nếu token còn hạn trong 30 phút tới
	if time.Until(claims.ExpiresAt.Time) > 30*time.Minute {
		return tokenString, nil
	}

	return GenerateJWT(claims.UserID, claims.Username, claims.Role)
}

// extractToken lấy token từ header hoặc cookie
func ExtractToken(c *gin.Context) string {
	// Ưu tiên từ header
	bearerToken := c.GetHeader("Authorization")
	if len(bearerToken) > 7 && strings.HasPrefix(bearerToken, "Bearer ") {
		return bearerToken[7:]
	}

	// Fallback lấy từ cookie
	if token, err := c.Cookie("jwt"); err == nil {
		return token
	}

	return ""
}
