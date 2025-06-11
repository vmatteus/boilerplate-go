package presentation

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/your-org/boilerplate-go/internal/user/application"
)

// UserController handles HTTP requests for users
type UserController struct {
	userService *application.UserService
	logger      zerolog.Logger
}

// NewUserController creates a new UserController
func NewUserController(userService *application.UserService, logger zerolog.Logger) *UserController {
	return &UserController{
		userService: userService,
		logger:      logger,
	}
}

// CreateUser handles POST /users
func (c *UserController) CreateUser(ctx *gin.Context) {
	var req CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := c.userService.CreateUser(ctx.Request.Context(), req.Name, req.Email)
	if err != nil {
		c.logger.Error().Err(err).Msg("Failed to create user")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}

	ctx.JSON(http.StatusCreated, response)
}

// GetUser handles GET /users/:id
func (c *UserController) GetUser(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := c.userService.GetUser(ctx.Request.Context(), uint(id))
	if err != nil {
		if err.Error() == "user not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.logger.Error().Err(err).Msg("Failed to get user")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}

	ctx.JSON(http.StatusOK, response)
}

// UpdateUser handles PUT /users/:id
func (c *UserController) UpdateUser(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req UpdateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := c.userService.UpdateUser(ctx.Request.Context(), uint(id), req.Name, req.Email)
	if err != nil {
		if err.Error() == "user not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.logger.Error().Err(err).Msg("Failed to update user")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}

	ctx.JSON(http.StatusOK, response)
}

// DeleteUser handles DELETE /users/:id
func (c *UserController) DeleteUser(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	err = c.userService.DeleteUser(ctx.Request.Context(), uint(id))
	if err != nil {
		if err.Error() == "user not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.logger.Error().Err(err).Msg("Failed to delete user")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

// ListUsers handles GET /users
func (c *UserController) ListUsers(ctx *gin.Context) {
	limitStr := ctx.DefaultQuery("limit", "10")
	offsetStr := ctx.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 10
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		offset = 0
	}

	users, err := c.userService.ListUsers(ctx.Request.Context(), limit, offset)
	if err != nil {
		c.logger.Error().Err(err).Msg("Failed to list users")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	userResponses := make([]UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		}
	}

	page := (offset / limit) + 1
	response := ListUsersResponse{
		Users: userResponses,
		Page:  page,
		Limit: limit,
	}

	ctx.JSON(http.StatusOK, response)
}

// RegisterRoutes registers user routes
func (c *UserController) RegisterRoutes(router *gin.RouterGroup) {
	users := router.Group("/users")
	{
		users.POST("", c.CreateUser)
		users.GET("/:id", c.GetUser)
		users.PUT("/:id", c.UpdateUser)
		users.DELETE("/:id", c.DeleteUser)
		users.GET("", c.ListUsers)
	}
}
