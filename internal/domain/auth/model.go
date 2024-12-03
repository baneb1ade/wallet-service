package auth

type RegisterRequest struct {
	Email    string `json:"email,required" binding:"required,email"`
	Username string `json:"username,required" binding:"required,min=3,max=16"`
	Password string `json:"password,required" binding:"required,min=8,max=32"`
}

type LoginRequest struct {
	Username string `json:"username,required" binding:"required"`
	Password string `json:"password,required" binding:"required"`
}
