package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID  `json:"id" gorm:"type:uuid;primary_key"`
	Username  string     `json:"username" gorm:"uniqueIndex"`
	Email     string     `json:"email" gorm:"uniqueIndex"`
	Password  string     `json:"password"`
	FullName  string     `json:"full_name"`
	Bio       string     `json:"bio"`
	Avatar    string     `json:"avatar"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"index"`
}

type UserPatchRequest struct {
	FullName *string `json:"full_name,omitempty"`
	Bio      *string `json:"bio,omitempty"`
	Avatar   *string `json:"avatar,omitempty"`
}

type UserRepository interface {
	GetByID(id string) (*User, error)
	GetByEmail(email string) (*User, error)
	Create(user *User) error
	Update(user *User) error
	Delete(id string) error
}

type UserUseCase interface {
	GetUser(id string) (*User, error)
	CreateUser(user *User) error
	UpdateUser(user *User) error
	PatchUser(id string, patch *UserPatchRequest) error
	DeleteUser(id string) error
}
