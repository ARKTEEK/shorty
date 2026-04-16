package models

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Message  string `json:"message"`
}

type LoginResponse struct {
	UserID  int64  `json:"user_id"`
	Email   string `json:"email"`
	Token   string `json:"token"`
	Message string `json:"message"`
}

type DeactivateRequest struct {
	UserID   int64  `json:"user_id"`
	Password string `json:"password"`
}

type DeactivateResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
