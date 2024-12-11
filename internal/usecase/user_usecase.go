package usecase

import (
	"errors"
	"socialnetwork/internal/domain"

	"github.com/google/uuid"
)

type userUseCase struct {
	userRepo domain.UserRepository
}

func NewUserUseCase(userRepo domain.UserRepository) domain.UserUseCase {
	return &userUseCase{
		userRepo: userRepo,
	}
}

func (u *userUseCase) GetUser(id string) (*domain.User, error) {
	if id == "" {
		return nil, errors.New("invalid user id")
	}
	return u.userRepo.GetByID(id)
}

func (u *userUseCase) CreateUser(user *domain.User) error {
	if user.Email == "" || user.Username == "" {
		return errors.New("email and username are required")
	}

	// Check if user with email already exists
	existingUser, _ := u.userRepo.GetByEmail(user.Email)
	if existingUser != nil {
		return errors.New("email already registered")
	}

	return u.userRepo.Create(user)
}

func (u *userUseCase) UpdateUser(user *domain.User) error {
	if user.ID == uuid.Nil {
		return errors.New("invalid user id")
	}

	existingUser, err := u.userRepo.GetByID(user.ID.String())
	if err != nil {
		return err
	}

	// Update only allowed fields
	existingUser.FullName = user.FullName
	existingUser.Bio = user.Bio
	existingUser.Avatar = user.Avatar

	return u.userRepo.Update(existingUser)
}

func (u *userUseCase) PatchUser(id string, patch *domain.UserPatchRequest) error {
	_, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid user id")
	}

	existingUser, err := u.userRepo.GetByID(id)
	if err != nil {
		return err
	}

	// Apply patches only if they are present in the request
	if patch.FullName != nil {
		existingUser.FullName = *patch.FullName
	}
	if patch.Bio != nil {
		existingUser.Bio = *patch.Bio
	}
	if patch.Avatar != nil {
		existingUser.Avatar = *patch.Avatar
	}

	return u.userRepo.Update(existingUser)
}

func (u *userUseCase) DeleteUser(id string) error {
	if id == "" {
		return errors.New("invalid user id")
	}
	return u.userRepo.Delete(id)
}
