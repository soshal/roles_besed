package middleware

import (
    "role/database"
    "role/models"
    "fmt"
    "net/http"
    "os"
    "strings"

    "github.com/dgrijalva/jwt-go"
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

// AuthMiddleware verifies JWT tokens and retrieves the user from the database.
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
            c.Abort()
            return
        }

        tokenString := strings.TrimPrefix(authHeader, "Bearer ")
        if tokenString == authHeader {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
            c.Abort()
            return
        }

        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
            }
            return []byte(os.Getenv("JWT_PRIVATE_KEY")), nil
        })

        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }

        if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
            userID := uint(claims["user_id"].(float64))
            var user models.User

            // Fetch user from the database
            if err := database.Db.First(&user, userID).Error; err != nil {
                if err == gorm.ErrRecordNotFound {
                    c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
                    c.Abort()
                    return
                }
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
                c.Abort()
                return
            }

            // Use the user variable for further processing if necessary
            fmt.Printf("Authenticated user: %v\n", user)

            // Attach user to the context
            c.Set("user", user)
            c.Next()
        } else {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
            c.Abort()
        }
    }
}

// Authorize middleware checks if the user has the required role and permission.
func Authorize(role, resource, action string) gin.HandlerFunc {
    return func(c *gin.Context) {
        user, exists := c.Get("user")
        if !exists {
            c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized"})
            c.Abort()
            return
        }

        // Assuming the user has Roles field
        roles := user.(models.User).Roles
        // Check if user has the required role and permission
        if hasPermission(roles, role, resource, action) {
            c.Next()
        } else {
            c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
            c.Abort()
        }
    }
}

func hasPermission(userRoles []models.Role, requiredRole, resource, action string) bool {
    for _, role := range userRoles {
        if role.Name == requiredRole {
            for _, permission := range role.Permissions {
                if permission.Resource == resource && permission.Action == action {
                    return true
                }
            }
        }
    }
    return false
}