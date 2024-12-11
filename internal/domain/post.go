package domain

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Content   string    `json:"content"`
	ImageURL  string    `json:"image_url,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PostRepository interface {
	GetByID(id string) (*Post, error)
	GetByUserID(userID string) ([]*Post, error)
	Create(post *Post) error
	Update(post *Post) error
	Delete(id string) error
}

type PostUseCase interface {
	GetPost(id string) (*Post, error)
	GetUserPosts(userID string) ([]*Post, error)
	CreatePost(post *Post) error
	UpdatePost(post *Post) error
	DeletePost(id string) error
}
