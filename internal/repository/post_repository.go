package repository

import (
	"socialnetwork/internal/domain"

	"gorm.io/gorm"
)

type postRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) domain.PostRepository {
	return &postRepository{
		db: db,
	}
}

func (r *postRepository) GetByID(id string) (*domain.Post, error) {
	var post domain.Post
	if err := r.db.Where("id = ?", id).First(&post).Error; err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *postRepository) GetByUserID(userID string) ([]*domain.Post, error) {
	var posts []*domain.Post
	if err := r.db.Where("user_id = ?", userID).Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

func (r *postRepository) Create(post *domain.Post) error {
	return r.db.Create(post).Error
}

func (r *postRepository) Update(post *domain.Post) error {
	return r.db.Save(post).Error
}

func (r *postRepository) Delete(id string) error {
	return r.db.Delete(&domain.Post{}, "id = ?", id).Error
}

func (r *postRepository) GetFeed(page int, limit int) ([]*domain.Post, error) {
	var posts []*domain.Post
	offset := (page - 1) * limit

	if err := r.db.Order("created_at desc").
		Offset(offset).
		Limit(limit).
		Find(&posts).Error; err != nil {
		return nil, err
	}

	return posts, nil
}
