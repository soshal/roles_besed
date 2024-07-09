package database

import (
    "fmt"
    "log"
    "os"

    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

var Db *gorm.DB

func InitDb() *gorm.DB {
    Db = connectDB()
    Migrate(Db)
    return Db
}

func connectDB() *gorm.DB {
    var err error
    host := os.Getenv("DB_HOST")
    username := os.Getenv("DB_USER")
    password := os.Getenv("DB_PASSWORD")
    dbname := os.Getenv("DB_NAME")
    port := os.Getenv("DB_PORT")

    dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC", host, username, password, dbname, port)
    //log.Println("dsn : ", dsn)
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

    if err != nil {
        log.Fatal("Error connecting to database :", err)
        return nil
    }
    log.Println("Successfully connected to the database")

    return db
}
