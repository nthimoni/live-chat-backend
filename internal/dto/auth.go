package dto

type UserDTO struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

type RegisterUserRequest struct {
	Username string `json:"username" validate:"required,min=3"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"` // password policy is verified in the service layer
}

type RegisterUserResponse struct {
	User  UserDTO `json:"user"`
	Token string  `json:"token"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"` // password policy is verified in the service layer
}

type LoginResponse struct {
	User  UserDTO `json:"user"`
	Token string  `json:"token"`
}
