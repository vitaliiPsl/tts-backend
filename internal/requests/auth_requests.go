package requests

type SignUpRequest struct {
	Username string `json:"username"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type SignInRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type SignInWithSSORequest struct {
	Code string `json:"code" validate:"required"`
}

type EmailVerificationRequest struct {
	Token    string `json:"token" validate:"required"`
}
