package models

import (
    "gorm.io/gorm"
)

type User struct {
    gorm.Model
    Username string `json:"username"`
    Email    string `json:"email"`
    Password string `json:"password"`
    Roles    []Role `gorm:"many2many:user_roles;" json:"roles"`
}

type Role struct {
    gorm.Model
    Name        string       `json:"name"`
    Permissions []Permission `gorm:"many2many:role_permissions;" json:"permissions"`
}

type Permission struct {
    gorm.Model
    Resource string `json:"resource"`
    Action   string `json:"action"`
}
