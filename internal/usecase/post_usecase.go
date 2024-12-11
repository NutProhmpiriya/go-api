package usecase

import (
	"errors"
	"socialnetwork/internal/domain"

	"github.com/google/uuid"
)

type postUseCase struct {
	postRepo domain.PostRepository
	userRepo domain.UserRepository
}

func NewPostUseCase(postRepo domain.PostRepository, userRepo domain.UserRepository) domain.PostUseCase {
	return &postUseCase{
		postRepo: postRepo,
		userRepo: userRepo,
	}
}

func (u *postUseCase) GetPost(id string) (*domain.Post, error) {
	if id == "" {
		return nil, errors.New("invalid post id")
	}
	return u.postRepo.GetByID(id)
}

func (u *postUseCase) GetUserPosts(userID string) ([]*domain.Post, error) {
	if userID == "" {
		return nil, errors.New("invalid user id")
	}
	return u.postRepo.GetByUserID(userID)
}

func (u *postUseCase) CreatePost(post *domain.Post) error {
	if post.UserID == uuid.Nil || post.Content == "" {
		return errors.New("user id and content are required")
	}

	// Verify user exists
	_, err := u.userRepo.GetByID(post.UserID.String())
	if err != nil {
		return errors.New("user not found")
	}

	return u.postRepo.Create(post)
}

func (u *postUseCase) UpdatePost(post *domain.Post) error {
	if post.ID == uuid.Nil {
		return errors.New("invalid post id")
	}

	existingPost, err := u.postRepo.GetByID(post.ID.String())
	if err != nil {
		return err
	}

	// Verify ownership
	if existingPost.UserID != post.UserID {
		return errors.New("unauthorized to update this post")
	}

	// Update allowed fields
	existingPost.Content = post.Content
	existingPost.Media = post.Media

	return u.postRepo.Update(existingPost)
}

func (u *postUseCase) DeletePost(id string) error {
	if id == "" {
		return errors.New("invalid post id")
	}
	return u.postRepo.Delete(id)
}

func (u *postUseCase) GetFeed(page int, limit int) ([]*domain.Post, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	return u.postRepo.GetFeed(page, limit)
}
