package middleware

import (
	"context"
	"fmt"
	"github.com/ashkanamani/rate-limiter/internal/limiter"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"strings"
)

func ExtractIdentity(c *gin.Context) string {
	// Extract JWT token from Authorization header
	authHeader := c.GetHeader("Authorization")

	if strings.HasPrefix(authHeader, "Bearer ") {
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Parse JWT token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			// Return your JWT secret here
			return []byte("jwt-config-secret"), nil
		})
		if err != nil || !token.Valid {
			log.Println("Invalid JWT token")
			return ""
		}
		// Extract user ID from the "sub" claim
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if sub, ok := claims["sub"].(string); ok {
				return sub // return the user ID
			}
		}
	}
	// If JWT is invalid or not present, fallback to the client IP
	return c.ClientIP() // Return the client's IP as the fallback identity
}

func NewRateLimiterMiddleware(limiter limiter.Limiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := ExtractIdentity(c)
		if key == "" {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Unable to extract identity",
			})
			c.Abort()
			return
		}
		ctx := context.Background()

		allowed, retryAfter, err := limiter.AllowRequest(ctx, key)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "rate limiter internal error",
			})
			return
		}
		if !allowed {
			c.Header("Retry-After", fmt.Sprintf("%d", retryAfter))
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "rate limiter reached",
			})
			return
		}
		c.Next()
	}
}
