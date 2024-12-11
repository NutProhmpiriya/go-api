package handler

import (
	"net/http"
	"socialnetwork/internal/domain"
	"socialnetwork/internal/middleware"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PostHandler struct {
	postUseCase domain.PostUseCase
	secretKey   string
}

func NewPostHandler(postUseCase domain.PostUseCase, secretKey string) *PostHandler {
	return &PostHandler{
		postUseCase: postUseCase,
		secretKey:   secretKey,
	}
}

func (h *PostHandler) Register(router *gin.RouterGroup) {
	posts := router.Group("/posts")
	posts.Use(middleware.JWTMiddleware(h.secretKey))
	{
		posts.POST("/", h.CreatePost)
		posts.GET("/:id", h.GetPost)
		posts.PUT("/:id", h.UpdatePost)
		posts.DELETE("/:id", h.DeletePost)
		posts.GET("/user/:id", h.GetUserPosts)
	}
}

// @Summary Create new post
// @Description Create a new post with the provided content
// @Tags posts
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param post body domain.Post true "Post content"
// @Success 201 {object} domain.Post
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /posts [post]
// @Security Bearer
func (h *PostHandler) CreatePost(c *gin.Context) {
	var post domain.Post
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Set post metadata
	post.ID = uuid.New()
	post.CreatedAt = time.Now()
	post.UpdatedAt = time.Now()

	// Get user ID from context (assuming it was set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Parse user ID
	parsedUserID, err := uuid.Parse(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	post.UserID = parsedUserID

	if err := h.postUseCase.CreatePost(&post); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, post)
}

// @Summary Get post by ID
// @Description Get post details by post ID
// @Tags posts
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param id path string true "Post ID"
// @Success 200 {object} domain.Post
// @Failure 404 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /posts/{id} [get]
// @Security Bearer
func (h *PostHandler) GetPost(c *gin.Context) {
	id := c.Param("id")

	post, err := h.postUseCase.GetPost(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, post)
}

// @Summary Update post
// @Description Update an existing post
// @Tags posts
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param id path string true "Post ID"
// @Param post body domain.Post true "Updated post content"
// @Success 200 {object} domain.Post
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /posts/{id} [put]
// @Security Bearer
func (h *PostHandler) UpdatePost(c *gin.Context) {
	id := c.Param("id")

	var post domain.Post
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Parse post ID
	parsedID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}
	post.ID = parsedID

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Parse user ID
	parsedUserID, err := uuid.Parse(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	post.UserID = parsedUserID
	post.UpdatedAt = time.Now()

	if err := h.postUseCase.UpdatePost(&post); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, post)
}

// @Summary Delete post
// @Description Delete an existing post
// @Tags posts
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param id path string true "Post ID"
// @Success 204 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /posts/{id} [delete]
// @Security Bearer
func (h *PostHandler) DeletePost(c *gin.Context) {
	id := c.Param("id")

	if err := h.postUseCase.DeletePost(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// @Summary Get user posts
// @Description Get all posts by user ID
// @Tags posts
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param id path string true "User ID"
// @Success 200 {array} domain.Post
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /posts/user/{id} [get]
// @Security Bearer
func (h *PostHandler) GetUserPosts(c *gin.Context) {
	userID := c.Param("id")

	posts, err := h.postUseCase.GetUserPosts(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, posts)
}
