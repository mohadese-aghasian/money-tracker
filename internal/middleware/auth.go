package middleware

import (
	"errors"
	"money-tracker/internal/config"
	"money-tracker/internal/constants"
	"money-tracker/internal/entity"

	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type UserContext struct {
	ID          uint   `json:"id"`
	Email       string `json:"email"`
	UserName    string `json:"user_name`
	LevelManage int8   `json:"level_manage"`
	Name        string `json:"name"`
}

// AuthMiddleware validates the JWT token
func AuthMiddleware(requiredLevel []int8) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the token from the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			c.Abort()
			return
		}

		// Ensure the token has the "Bearer " prefix
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer {token}"})
			c.Abort()
			return
		}

		// Validate the token
		_, err := validateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		var userToken entity.UserToken
		if err := config.DB.Where("token = ?", tokenString).First(&userToken).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token is invalid or expired"})
			c.Abort()
			return
		}

		var user entity.User
		if err := config.DB.Where("id = ? AND status_id = ?", userToken.UserID, constants.StatusActive).Select("id", "level_manage", "user_name").First(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find user"})
			c.Abort()
			return
		}

		// ✅ Check if the user has the required level (with single number)
		// if user.LevelManage != requiredLevel {
		// 	c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to access this resource"})
		// 	c.Abort()
		// 	return
		// }

		//with array
		isAllowed := false
		for _, level := range requiredLevel {
			if user.LevelManage == level {
				isAllowed = true
				break
			}
		}

		// If not allowed, return Forbidden
		if !isAllowed {
			c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to access this resource"})
			c.Abort()
			return
		}

		// userContext := UserContext{
		// 	ID:          user.ID,
		// 	Email:       user.Email,
		// 	LevelManage: user.LevelManage,
		// 	Name:        user.Name,
		// 	UserName:    user.UserName,
		// }

		c.Set("user", user)
		// Add the claims to the context
		// c.Set("user_id", claims["user_id"])

		c.Next()
	}
}

// validateToken validates the JWT and returns the claims
func validateToken(tokenString string) (jwt.MapClaims, error) {
	// Get the secret key from the environment
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		return nil, errors.New("JWT_SECRET environment variable not set")
	}

	// Parse and validate the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secretKey), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	// Extract claims from the token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}