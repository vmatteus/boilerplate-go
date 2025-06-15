package application_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/your-org/boilerplate-go/internal/config"
	"github.com/your-org/boilerplate-go/internal/logger"
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

// createTestLogger creates a logger instance for testing
func createTestLogger() *logger.Logger {
	cfg := config.LoggerConfig{
		Level:    "error", // Use error level to reduce test output
		Format:   "json",
		Provider: "stdout",
	}
	appLogger := logger.InitLogger(cfg)
	return &appLogger
}

func TestUserService_CreateUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	testLogger := createTestLogger()
	service := application.NewUserService(mockRepo, testLogger)
	ctx := context.Background()

	t.Run("successful user creation", func(t *testing.T) {
		user := &domain.User{
			ID:    1,
			Name:  "John Doe",
			Email: "john@example.com",
		}

		mockRepo.On("GetByEmail", mock.Anything, "john@example.com").Return(nil, errors.New("user not found"))
		mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(u *domain.User) bool {
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
		testLogger := createTestLogger()
		service := application.NewUserService(mockRepo, testLogger)

		mockRepo.On("GetByEmail", mock.Anything, "john@example.com").Return(existingUser, nil)

		result, err := service.CreateUser(ctx, "John Doe", "john@example.com")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "already exists")
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_GetUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	testLogger := createTestLogger()
	service := application.NewUserService(mockRepo, testLogger)
	ctx := context.Background()

	t.Run("successful user retrieval", func(t *testing.T) {
		user := &domain.User{
			ID:    1,
			Name:  "John Doe",
			Email: "john@example.com",
		}

		mockRepo.On("GetByID", mock.Anything, uint(1)).Return(user, nil)

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

func TestUserService_GetUserByEmail(t *testing.T) {
	mockRepo := new(MockUserRepository)
	testLogger := createTestLogger()
	service := application.NewUserService(mockRepo, testLogger)
	ctx := context.Background()

	t.Run("successful user retrieval by email", func(t *testing.T) {
		user := &domain.User{
			ID:    1,
			Name:  "John Doe",
			Email: "john@example.com",
		}

		mockRepo.On("GetByEmail", mock.Anything, "john@example.com").Return(user, nil)

		result, err := service.GetUserByEmail(ctx, "john@example.com")

		assert.NoError(t, err)
		assert.Equal(t, user, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("user retrieval with empty email", func(t *testing.T) {
		result, err := service.GetUserByEmail(ctx, "")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "email is required", err.Error())
	})
}

func TestUserService_UpdateUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	testLogger := createTestLogger()
	service := application.NewUserService(mockRepo, testLogger)
	ctx := context.Background()

	t.Run("successful user update", func(t *testing.T) {
		user := &domain.User{
			ID:    1,
			Name:  "John Doe",
			Email: "john@example.com",
		}

		mockRepo.On("GetByID", mock.Anything, uint(1)).Return(user, nil)
		mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(u *domain.User) bool {
			return u.ID == 1 && u.Name == "John Smith" && u.Email == "johnsmith@example.com"
		})).Return(nil)

		result, err := service.UpdateUser(ctx, 1, "John Smith", "johnsmith@example.com")

		assert.NoError(t, err)
		assert.Equal(t, "John Smith", result.Name)
		assert.Equal(t, "johnsmith@example.com", result.Email)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_DeleteUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	testLogger := createTestLogger()
	service := application.NewUserService(mockRepo, testLogger)
	ctx := context.Background()

	t.Run("successful user deletion", func(t *testing.T) {
		mockRepo.On("Delete", mock.Anything, uint(1)).Return(nil)

		err := service.DeleteUser(ctx, 1)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("user deletion with invalid ID", func(t *testing.T) {
		err := service.DeleteUser(ctx, 0)

		assert.Error(t, err)
		assert.Equal(t, "invalid user ID", err.Error())
	})
}

func TestUserService_ListUsers(t *testing.T) {
	t.Run("successful user listing", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		testLogger := createTestLogger()
		service := application.NewUserService(mockRepo, testLogger)
		ctx := context.Background()

		users := []*domain.User{
			{ID: 1, Name: "John Doe", Email: "john@example.com"},
			{ID: 2, Name: "Jane Doe", Email: "jane@example.com"},
		}

		mockRepo.On("List", mock.Anything, 10, 0).Return(users, nil)

		result, err := service.ListUsers(ctx, 10, 0)

		assert.NoError(t, err)
		assert.Equal(t, users, result)
		assert.Len(t, result, 2)
		mockRepo.AssertExpectations(t)
	})

	t.Run("user listing with default pagination", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		testLogger := createTestLogger()
		service := application.NewUserService(mockRepo, testLogger)
		ctx := context.Background()

		users := []*domain.User{
			{ID: 1, Name: "John Doe", Email: "john@example.com"},
		}

		mockRepo.On("List", mock.Anything, 10, 0).Return(users, nil)

		result, err := service.ListUsers(ctx, 0, -1)

		assert.NoError(t, err)
		assert.Equal(t, users, result)
		mockRepo.AssertExpectations(t)
	})
}
