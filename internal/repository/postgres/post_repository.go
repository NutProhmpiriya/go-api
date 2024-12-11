package postgres

import (
	"socialnetwork/internal/domain"

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
	
	if err := r.db.Where("id = ?", uid).First(&post).Error; err != nil {
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
	
	if err := r.db.Where("user_id = ?", uid).Order("created_at DESC").Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

func (r *postRepository) Create(post *domain.Post) error {
	post.ID = uuid.New()
	return r.db.Create(post).Error
}

func (r *postRepository) Update(post *domain.Post) error {
	return r.db.Save(post).Error
}

func (r *postRepository) Delete(id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return r.db.Delete(&domain.Post{}, "id = ?", uid).Error
}

func (r *postRepository) GetFeed(userID string, page, limit int) ([]*domain.Post, error) {
	var posts []*domain.Post
	offset := (page - 1) * limit

	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	err = r.db.Where("user_id = ?", uid).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&posts).Error

	if err != nil {
		return nil, err
	}
	return posts, nil
}
