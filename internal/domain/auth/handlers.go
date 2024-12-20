package auth

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
)

// Register godoc
// @Summary      Register a new user
// @Description  Register a new user with email, username, and password
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      RegisterRequest  true  "Registration request"
// @Success      200      {object}  map[string]string       "User registered successfully"
// @Router       /api/v1/auth/register/ [post]
func Register(s ServiceAuth, wc WalletCreator, v *validator.Validate) func(c *gin.Context) {
	return func(c *gin.Context) {
		var req RegisterRequest

		if err := c.ShouldBind(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}

		if err := v.Struct(req); err != nil {
			var validationErrors validator.ValidationErrors
			errors.As(err, &validationErrors)
			invalidFields := make([]string, len(validationErrors))

			for i, fieldError := range validationErrors {
				invalidFields[i] = fieldError.Field()
			}

			c.JSON(http.StatusBadRequest, gin.H{
				"error":  "Validation failed",
				"fields": invalidFields,
			})
			return
		}
		_, err := s.Register(context.Background(), wc, req.Email, req.Username, req.Password)
		if err != nil {
			st, ok := status.FromError(err)
			if !ok {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Unknown error"})
				return
			}

			switch st.Code() {
			case codes.AlreadyExists:
				c.JSON(http.StatusBadRequest, gin.H{"error": "Username or email already exists"})
			case codes.InvalidArgument:
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid username or password"})
			case codes.Internal:
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			}
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
	}
}

// Login godoc
// @Summary      User login
// @Description  Authenticate a user and return a JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      LoginRequest  true  "Login request"
// @Success      200      {object}  map[string]string       "JWT token"
// @Failure      400      {object}  map[string]interface{}  "Invalid request or validation failed"
// @Failure      401      {object}  map[string]string       "Invalid username or password"
// @Failure      500      {object}  map[string]string       "Internal server error"
// @Router       /api/v1/auth/login/ [post]
func Login(s ServiceAuth, v *validator.Validate) func(c *gin.Context) {
	return func(c *gin.Context) {
		var req LoginRequest
		if err := c.ShouldBind(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}

		if err := v.Struct(req); err != nil {
			var validationErrors validator.ValidationErrors
			errors.As(err, &validationErrors)
			invalidFields := make([]string, len(validationErrors))

			for i, fieldError := range validationErrors {
				invalidFields[i] = fieldError.Field()
			}

			c.JSON(http.StatusBadRequest, gin.H{
				"error":  "Validation failed",
				"fields": invalidFields,
			})
			return
		}
		token, err := s.Login(context.Background(), req.Username, req.Password)
		if err != nil {
			st, ok := status.FromError(err)
			if !ok {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Unknown error"})
				return
			}

			switch st.Code() {
			case codes.Unauthenticated:
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
			case codes.InvalidArgument:
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
			case codes.Internal:
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			}
			return
		}
		c.JSON(http.StatusOK, gin.H{"token": token})
	}
}
