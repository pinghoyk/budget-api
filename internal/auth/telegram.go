package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// TelegramAuthData represents the data received from Telegram Login Widget
type TelegramAuthData struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name,omitempty"`
	Username  string `json:"username,omitempty"`
	PhotoURL  string `json:"photo_url,omitempty"`
	AuthDate  int64  `json:"auth_date"`
	Hash      string `json:"hash"`
}

// VerifyTelegramAuth verifies the authentication data from Telegram
func VerifyTelegramAuth(data TelegramAuthData, botToken string) error {
	// Check if auth is not too old (24 hours)
	now := time.Now().Unix()
	if now-data.AuthDate > 86400 {
		return fmt.Errorf("authentication data is too old")
	}

	// Create data check string
	var parts []string
	
	if data.FirstName != "" {
		parts = append(parts, fmt.Sprintf("first_name=%s", data.FirstName))
	}
	if data.LastName != "" {
		parts = append(parts, fmt.Sprintf("last_name=%s", data.LastName))
	}
	if data.Username != "" {
		parts = append(parts, fmt.Sprintf("username=%s", data.Username))
	}
	if data.PhotoURL != "" {
		parts = append(parts, fmt.Sprintf("photo_url=%s", data.PhotoURL))
	}
	parts = append(parts, fmt.Sprintf("auth_date=%d", data.AuthDate))
	parts = append(parts, fmt.Sprintf("id=%d", data.ID))
	
	// Sort alphabetically
	sort.Strings(parts)
	
	// Join with newline
	dataCheckString := strings.Join(parts, "\n")
	
	// Create secret key from bot token
	secretKey := sha256.Sum256([]byte(botToken))
	
	// Create HMAC
	h := hmac.New(sha256.New, secretKey[:])
	h.Write([]byte(dataCheckString))
	hash := hex.EncodeToString(h.Sum(nil))
	
	// Compare hashes
	if hash != data.Hash {
		return fmt.Errorf("invalid hash")
	}
	
	return nil
}

// JWTClaims represents the JWT claims
type JWTClaims struct {
	UserID     int64  `json:"user_id"`
	TelegramID int64  `json:"telegram_id"`
	jwt.RegisteredClaims
}

// GenerateJWT generates a JWT token for a user
func GenerateJWT(userID, telegramID int64, jwtSecret string) (string, error) {
	claims := JWTClaims{
		UserID:     userID,
		TelegramID: telegramID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)), // 7 days
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

// ValidateJWT validates a JWT token and returns the claims
func ValidateJWT(tokenString, jwtSecret string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})
	
	if err != nil {
		return nil, err
	}
	
	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}
	
	return nil, fmt.Errorf("invalid token")
}
