package postgres

import (
	"socialnetwork/internal/domain"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type postRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) domain.PostRepository {
	return &postRepository{db: db}
}

func (r *postRepository) GetByID(id string) (*domain.Post, error) {
	var post domain.Post
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	
	if err := r.db.Where("id = ? AND deleted_at IS NULL", uid).First(&post).Error; err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *postRepository) GetByUserID(userID string) ([]*domain.Post, error) {
	var posts []*domain.Post
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}
	
	if err := r.db.Where("user_id = ? AND deleted_at IS NULL", uid).Order("created_at DESC").Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

func (r *postRepository) Create(post *domain.Post) error {
	post.ID = uuid.New()
	post.CreatedAt = time.Now()
	post.UpdatedAt = time.Now()
	return r.db.Create(post).Error
}

func (r *postRepository) Update(post *domain.Post) error {
	post.UpdatedAt = time.Now()
	return r.db.Model(post).Where("id = ? AND deleted_at IS NULL", post.ID).
		Updates(map[string]interface{}{
			"content":    post.Content,
			"media":      post.Media,
			"updated_at": post.UpdatedAt,
		}).Error
}

func (r *postRepository) Delete(id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return r.db.Model(&domain.Post{}).
		Where("id = ? AND deleted_at IS NULL", uid).
		Update("deleted_at", time.Now()).Error
}

func (r *postRepository) GetFeed(page, limit int) ([]*domain.Post, error) {
	var posts []*domain.Post
	offset := (page - 1) * limit

	if err := r.db.Where("deleted_at IS NULL").
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&posts).Error; err != nil {
		return nil, err
	}

	return posts, nil
}
