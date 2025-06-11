package application_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/your-org/boilerplate-go/internal/user/application"
	"github.com/your-org/boilerplate-go/internal/user/domain"
)

// MockUserRepository is a mock implementation of UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	args := m.Called(ctx, user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id uint) (*domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) List(ctx context.Context, limit, offset int) ([]*domain.User, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.User), args.Error(1)
}

func (m *MockUserRepository) Count(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

func TestUserService_CreateUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := application.NewUserService(mockRepo)
	ctx := context.Background()

	t.Run("successful user creation", func(t *testing.T) {
		user := &domain.User{
			ID:    1,
			Name:  "John Doe",
			Email: "john@example.com",
		}

		mockRepo.On("GetByEmail", ctx, "john@example.com").Return(nil, errors.New("user not found"))
		mockRepo.On("Create", ctx, mock.MatchedBy(func(u *domain.User) bool {
			return u.Name == "John Doe" && u.Email == "john@example.com"
		})).Return(user, nil)

		result, err := service.CreateUser(ctx, "John Doe", "john@example.com")

		assert.NoError(t, err)
		assert.Equal(t, user, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("user creation with empty name", func(t *testing.T) {
		result, err := service.CreateUser(ctx, "", "john@example.com")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "name is required", err.Error())
	})

	t.Run("user creation with empty email", func(t *testing.T) {
		result, err := service.CreateUser(ctx, "John Doe", "")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "email is required", err.Error())
	})

	t.Run("user creation with existing email", func(t *testing.T) {
		existingUser := &domain.User{
			ID:    1,
			Name:  "Existing User",
			Email: "john@example.com",
		}

		// Reset mock for this test
		mockRepo := new(MockUserRepository)
		service := application.NewUserService(mockRepo)

		mockRepo.On("GetByEmail", ctx, "john@example.com").Return(existingUser, nil)

		result, err := service.CreateUser(ctx, "John Doe", "john@example.com")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "already exists")
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_GetUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := application.NewUserService(mockRepo)
	ctx := context.Background()

	t.Run("successful user retrieval", func(t *testing.T) {
		user := &domain.User{
			ID:    1,
			Name:  "John Doe",
			Email: "john@example.com",
		}

		mockRepo.On("GetByID", ctx, uint(1)).Return(user, nil)

		result, err := service.GetUser(ctx, 1)

		assert.NoError(t, err)
		assert.Equal(t, user, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("user retrieval with invalid ID", func(t *testing.T) {
		result, err := service.GetUser(ctx, 0)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "invalid user ID", err.Error())
	})
}
