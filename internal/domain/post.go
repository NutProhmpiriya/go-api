package domain

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	ID        uuid.UUID  `json:"id" gorm:"type:uuid;primary_key"`
	UserID    uuid.UUID  `json:"user_id" gorm:"type:uuid"`
	Content   string     `json:"content"`
	Media     []string   `json:"media" gorm:"type:text[]"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

type PostRepository interface {
	GetByID(id string) (*Post, error)
	GetByUserID(userID string) ([]*Post, error)
	Create(post *Post) error
	Update(post *Post) error
	Delete(id string) error
	GetFeed(page int, limit int) ([]*Post, error)
}

type PostUseCase interface {
	GetPost(id string) (*Post, error)
	GetUserPosts(userID string) ([]*Post, error)
	CreatePost(post *Post) error
	UpdatePost(post *Post) error
	DeletePost(id string) error
	GetFeed(page int, limit int) ([]*Post, error)
}
