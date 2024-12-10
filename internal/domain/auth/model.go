package auth

type RegisterRequest struct {
	Email    string `json:"email,required" example:"user@example.com" binding:"required,email"`
	Username string `json:"username,required" example:"user123" binding:"required,min=3,max=16"`
	Password string `json:"password,required" example:"password123" binding:"required,min=8,max=32"`
}

type LoginRequest struct {
	Username string `json:"username,required" example:"user123" binding:"required"`
	Password string `json:"password,required" example:"password123" binding:"required"`
}
