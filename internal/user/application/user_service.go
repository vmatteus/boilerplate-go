package application

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/your-org/boilerplate-go/internal/logger"
	"github.com/your-org/boilerplate-go/internal/user/domain"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

// UserService handles user business logic
type UserService struct {
	userRepo domain.UserRepository
	logger   *logger.Logger
}

// NewUserService creates a new UserService
func NewUserService(userRepo domain.UserRepository, logger *logger.Logger) *UserService {
	return &UserService{
		userRepo: userRepo,
		logger:   logger,
	}
}

// CreateUser creates a new user
func (s *UserService) CreateUser(ctx context.Context, name, email string) (*domain.User, error) {
	ctx, span := otel.Tracer("user-service").Start(ctx, "UserService.CreateUser")
	defer span.End()

	// Add span attributes
	span.SetAttributes(
		attribute.String("user.email", email),
		attribute.String("user.name", name),
	)

	start := time.Now()

	// Log start of operation
	s.logger.LogInfo(ctx, "Starting user creation", map[string]interface{}{
		"email": email,
		"name":  name,
	})

	if name == "" {
		s.logger.LogError(ctx, "User creation failed: name is required", nil, map[string]interface{}{
			"validation_error": "name_required",
			"email":            email,
		})
		return nil, errors.New("name is required")
	}
	if email == "" {
		s.logger.LogError(ctx, "User creation failed: email is required", nil, map[string]interface{}{
			"validation_error": "email_required",
			"name":             name,
		})
		return nil, errors.New("email is required")
	}

	// Check if user already exists
	existingUser, err := s.userRepo.GetByEmail(ctx, email)
	if err == nil && existingUser != nil {
		s.logger.LogWarn(ctx, "User creation failed: user already exists", map[string]interface{}{
			"email":       email,
			"existing_id": existingUser.ID,
		})
		return nil, fmt.Errorf("user with email %s already exists", email)
	}

	user := &domain.User{
		Name:  name,
		Email: email,
	}

	createdUser, err := s.userRepo.Create(ctx, user)
	if err != nil {
		s.logger.LogError(ctx, "Failed to create user in repository", err, map[string]interface{}{
			"email": email,
			"name":  name,
		})
		return nil, err
	}

	duration := time.Since(start)
	s.logger.LogInfo(ctx, "User created successfully", map[string]interface{}{
		"user_id":  createdUser.ID,
		"email":    createdUser.Email,
		"duration": duration.Milliseconds(),
	})

	return createdUser, nil
}

// GetUser retrieves a user by ID
func (s *UserService) GetUser(ctx context.Context, id uint) (*domain.User, error) {
	ctx, span := otel.Tracer("user-service").Start(ctx, "UserService.GetUser")
	defer span.End()

	span.SetAttributes(attribute.Int64("user.id", int64(id)))

	s.logger.LogDebug(ctx, "Retrieving user by ID", map[string]interface{}{
		"user_id": id,
	})

	if id == 0 {
		s.logger.LogError(ctx, "Invalid user ID provided", nil, map[string]interface{}{
			"user_id": id,
		})
		return nil, errors.New("invalid user ID")
	}

	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.LogError(ctx, "Failed to retrieve user from repository", err, map[string]interface{}{
			"user_id": id,
		})
		return nil, err
	}

	s.logger.LogInfo(ctx, "User retrieved successfully", map[string]interface{}{
		"user_id": user.ID,
		"email":   user.Email,
	})

	return user, nil
}

// GetUserByEmail retrieves a user by email
func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	ctx, span := otel.Tracer("user-service").Start(ctx, "UserService.GetUserByEmail")
	defer span.End()

	span.SetAttributes(attribute.String("user.email", email))

	s.logger.LogDebug(ctx, "Retrieving user by email", map[string]interface{}{
		"email": email,
	})

	if email == "" {
		s.logger.LogError(ctx, "Email is required for user lookup", nil)
		return nil, errors.New("email is required")
	}

	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		s.logger.LogError(ctx, "Failed to retrieve user by email", err, map[string]interface{}{
			"email": email,
		})
		return nil, err
	}

	s.logger.LogInfo(ctx, "User retrieved by email successfully", map[string]interface{}{
		"user_id": user.ID,
		"email":   user.Email,
	})

	return user, nil
}

// UpdateUser updates an existing user
func (s *UserService) UpdateUser(ctx context.Context, id uint, name, email string) (*domain.User, error) {
	ctx, span := otel.Tracer("user-service").Start(ctx, "UserService.UpdateUser")
	defer span.End()

	span.SetAttributes(
		attribute.Int64("user.id", int64(id)),
		attribute.String("user.new_email", email),
		attribute.String("user.new_name", name),
	)

	s.logger.LogInfo(ctx, "Starting user update", map[string]interface{}{
		"user_id":   id,
		"new_name":  name,
		"new_email": email,
	})

	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.LogError(ctx, "Failed to retrieve user for update", err, map[string]interface{}{
			"user_id": id,
		})
		return nil, err
	}

	oldName := user.Name
	oldEmail := user.Email

	if name != "" {
		user.Name = name
	}
	if email != "" {
		user.Email = email
	}

	err = s.userRepo.Update(ctx, user)
	if err != nil {
		s.logger.LogError(ctx, "Failed to update user in repository", err, map[string]interface{}{
			"user_id":   id,
			"old_name":  oldName,
			"old_email": oldEmail,
			"new_name":  user.Name,
			"new_email": user.Email,
		})
		return nil, err
	}

	s.logger.LogInfo(ctx, "User updated successfully", map[string]interface{}{
		"user_id":   user.ID,
		"old_name":  oldName,
		"old_email": oldEmail,
		"new_name":  user.Name,
		"new_email": user.Email,
	})

	return user, nil
}

// DeleteUser deletes a user
func (s *UserService) DeleteUser(ctx context.Context, id uint) error {
	ctx, span := otel.Tracer("user-service").Start(ctx, "UserService.DeleteUser")
	defer span.End()

	span.SetAttributes(attribute.Int64("user.id", int64(id)))

	s.logger.LogInfo(ctx, "Starting user deletion", map[string]interface{}{
		"user_id": id,
	})

	if id == 0 {
		s.logger.LogError(ctx, "Invalid user ID for deletion", nil, map[string]interface{}{
			"user_id": id,
		})
		return errors.New("invalid user ID")
	}

	err := s.userRepo.Delete(ctx, id)
	if err != nil {
		s.logger.LogError(ctx, "Failed to delete user from repository", err, map[string]interface{}{
			"user_id": id,
		})
		return err
	}

	s.logger.LogInfo(ctx, "User deleted successfully", map[string]interface{}{
		"user_id": id,
	})

	return nil
}

// ListUsers retrieves users with pagination
func (s *UserService) ListUsers(ctx context.Context, limit, offset int) ([]*domain.User, error) {
	ctx, span := otel.Tracer("user-service").Start(ctx, "UserService.ListUsers")
	defer span.End()

	span.SetAttributes(
		attribute.Int("pagination.limit", limit),
		attribute.Int("pagination.offset", offset),
	)

	s.logger.LogDebug(ctx, "Listing users with pagination", map[string]interface{}{
		"limit":  limit,
		"offset": offset,
	})

	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	users, err := s.userRepo.List(ctx, limit, offset)
	if err != nil {
		s.logger.LogError(ctx, "Failed to retrieve users list", err, map[string]interface{}{
			"limit":  limit,
			"offset": offset,
		})
		return nil, err
	}

	s.logger.LogInfo(ctx, "Users list retrieved successfully", map[string]interface{}{
		"count":  len(users),
		"limit":  limit,
		"offset": offset,
	})

	return users, nil
}
