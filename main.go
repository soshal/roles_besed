package main

import (
    "role/database"
    "role/middleware"
    "role/models"
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"
    "log"
    "fmt"
    "github.com/dgrijalva/jwt-go"
    "os"
    "time"
	"gorm.io/gorm"
)

func generateToken(userID uint) (string, error) {
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id": userID,
        "exp":     time.Now().Add(time.Hour * 24).Unix(), // Token expiration time
    })

    tokenString, err := token.SignedString([]byte(os.Getenv("JWT_PRIVATE_KEY")))
    if err != nil {
        return "", err
    }

    return tokenString, nil
}

func main() {
    loadEnv()
    db := loadDatabase()
    
    // Automatically migrate the schema
    models.AutoMigrate(db)
    
    token, err := generateToken(1) // Replace 1 with the actual user ID you want to generate a token for
    if err != nil {
        log.Fatalf("Error generating token: %v", err)
    }

    log.Println("Generated Token:", token)
    
    serveApplication()
}

func loadEnv() {
    err := godotenv.Load(".env")
    if err != nil {
        log.Fatal("Error loading .env file")
    }
    log.Println(".env file loaded successfully")
}

func loadDatabase() *gorm.DB {
    db := database.InitDb()
    return db
}

func serveApplication() {
    router := gin.Default()

    router.LoadHTMLGlob("templates/*")
    router.Use(middleware.AuthMiddleware())

    router.GET("/customer", func(c *gin.Context) {
        c.HTML(http.StatusOK, "customer.html", nil)
    })

    router.GET("/billing", func(c *gin.Context) {
        c.HTML(http.StatusOK, "billing.html", nil)
    })

    router.GET("/payroll", func(c *gin.Context) {
        c.HTML(http.StatusOK, "payroll.html", nil)
    })

    router.GET("/user", func(c *gin.Context) {
        c.HTML(http.StatusOK, "user.html", nil)
    })

    customer := router.Group("/customer")
    {
        customer.GET("/", middleware.Authorize("sales", "customer", "read"), getCustomer)
        customer.POST("/", middleware.Authorize("sales", "customer", "write"), createCustomer)
    }

    billing := router.Group("/billing")
    {
        billing.GET("/", middleware.Authorize("sales", "billing", "read"), getBilling)
        billing.POST("/", middleware.Authorize("sales", "billing", "write"), createBilling)
    }

    payroll := router.Group("/payroll")
    {
        payroll.GET("/", middleware.Authorize("hr", "payroll", "read"), getPayroll)
        payroll.POST("/", middleware.Authorize("hr", "payroll", "write"), createPayroll)
    }

    user := router.Group("/user")
    {
        user.GET("/", middleware.Authorize("admin", "user", "read"), getUser)
        user.POST("/", middleware.Authorize("admin", "user", "write"), createUser)
    }

    router.Run(":8000")
    fmt.Println("Server running on port 8000")
}

func getCustomer(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{"message": "Get customer"})
}

func createCustomer(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{"message": "Create customer"})
}

func getBilling(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{"message": "Get billing"})
}

func createBilling(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{"message": "Create billing"})
}

func getPayroll(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{"message": "Get payroll"})
}

func createPayroll(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{"message": "Create payroll"})
}

func getUser(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{"message": "Get user"})
}

func createUser(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{"message": "Create user"})
}
