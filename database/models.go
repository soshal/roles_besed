package database

import (
    "gorm.io/gorm"
)

type User struct {
    gorm.Model
    Username string `gorm:"unique; not null"`
    Password string `gorm:"not null"`
    RoleID   uint
    Role     Role
}

type Role struct {
    gorm.Model
    Name        string `gorm:"unique;not null"`
    Permissions []Permission
}

type Permission struct {
    gorm.Model
    RoleID  uint
    Resource string `gorm:"not null"`
    Action  string `gorm:"not null"`
    Role    Role
}

func Migrate(db *gorm.DB) {
    db.AutoMigrate(&User{}, &Role{}, &Permission{})
}
