package application

import (
	"context"
	"errors"
	"fmt"

	"github.com/your-org/boilerplate-go/internal/user/domain"
)

// UserService handles user business logic
type UserService struct {
	userRepo domain.UserRepository
}

// NewUserService creates a new UserService
func NewUserService(userRepo domain.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// CreateUser creates a new user
func (s *UserService) CreateUser(ctx context.Context, name, email string) (*domain.User, error) {
	if name == "" {
		return nil, errors.New("name is required")
	}
	if email == "" {
		return nil, errors.New("email is required")
	}

	// Check if user already exists
	existingUser, err := s.userRepo.GetByEmail(ctx, email)
	if err == nil && existingUser != nil {
		return nil, fmt.Errorf("user with email %s already exists", email)
	}

	user := &domain.User{
		Name:  name,
		Email: email,
	}

	return s.userRepo.Create(ctx, user)
}

// GetUser retrieves a user by ID
func (s *UserService) GetUser(ctx context.Context, id uint) (*domain.User, error) {
	if id == 0 {
		return nil, errors.New("invalid user ID")
	}

	return s.userRepo.GetByID(ctx, id)
}

// GetUserByEmail retrieves a user by email
func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	if email == "" {
		return nil, errors.New("email is required")
	}

	return s.userRepo.GetByEmail(ctx, email)
}

// UpdateUser updates an existing user
func (s *UserService) UpdateUser(ctx context.Context, id uint, name, email string) (*domain.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if name != "" {
		user.Name = name
	}
	if email != "" {
		user.Email = email
	}

	err = s.userRepo.Update(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// DeleteUser deletes a user
func (s *UserService) DeleteUser(ctx context.Context, id uint) error {
	if id == 0 {
		return errors.New("invalid user ID")
	}

	return s.userRepo.Delete(ctx, id)
}

// ListUsers retrieves users with pagination
func (s *UserService) ListUsers(ctx context.Context, limit, offset int) ([]*domain.User, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	return s.userRepo.List(ctx, limit, offset)
}
