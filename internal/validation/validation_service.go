package validation

import (
	"fmt"
	"unicode"

	"github.com/go-playground/validator/v10"
	service_errors "vitaliiPsl/synthesizer/internal/error"
	"vitaliiPsl/synthesizer/internal/logger"
	"vitaliiPsl/synthesizer/internal/requests"
)

const MinPasswordLength = 8

type ValidationService struct {
	v *validator.Validate
}

func NewValidationService() *ValidationService {
	return &ValidationService{
		v: validator.New(),
	}
}

func (vs *ValidationService) ValidateSignUpRequest(request *requests.SignUpRequest) error {
	if err := vs.v.Struct(request); err != nil {
		logger.Logger.Error("Request didn't pass validation", "message", err.Error())
		return service_errors.NewErrBadRequest("Request didn't pass validation")
	}

	return validatePassword(request.Password)
}

func (vs *ValidationService) ValidateSignInRequest(request *requests.SignInRequest) error {
	if err := vs.v.Struct(request); err != nil {
		logger.Logger.Error("Request didn't pass validation", "message", err.Error())
		return service_errors.NewErrBadRequest("Request didn't pass validation")
	}

	return validatePassword(request.Password)
}

func (vs *ValidationService) ValidateVerificationTokenRequest(request *requests.VerificationTokenRequest) error {
	if err := vs.v.Struct(request); err != nil {
		logger.Logger.Error("Request didn't pass validation", "message", err.Error())
		return service_errors.NewErrBadRequest("Request didn't pass validation")
	}

	return nil
}

func (vs *ValidationService) ValidateEmailVerificationRequest(request *requests.EmailVerificationRequest) error {
	if err := vs.v.Struct(request); err != nil {
		logger.Logger.Error("Request didn't pass validation", "message", err.Error())
		return service_errors.NewErrBadRequest("Request didn't pass validation")
	}

	return nil
}

func (vs *ValidationService) ValidatePasswordResetRequest(request *requests.PasswordResetRequest) error {
	if err := vs.v.Struct(request); err != nil {
		logger.Logger.Error("Request didn't pass validation", "message", err.Error())
		return service_errors.NewErrBadRequest("Request didn't pass validation")
	}

	return validatePassword(request.Password)
}

func (vs *ValidationService) ValidateSynthesisRequest(request *requests.SynthesisRequest) error {
	if err := vs.v.Struct(request); err != nil {
		logger.Logger.Error("Request didn't pass validation", "message", err.Error())
		return service_errors.NewErrBadRequest("Request didn't pass validation")
	}

	return nil
}

func (vs *ValidationService) ValidateModelRequest(request *requests.ModelRequest) error {
	if err := vs.v.Struct(request); err != nil {
		logger.Logger.Error("Request didn't pass validation", "message", err.Error())
		return service_errors.NewErrBadRequest("Request didn't pass validation")
	}

	return nil
}

func validatePassword(password string) error {
	var hasUpper, hasLower, hasNumber, hasSpecial bool

	if len(password) < MinPasswordLength {
		return service_errors.NewErrBadRequest(fmt.Sprintf("password must be at least %d characters long", MinPasswordLength))
	}

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !(hasUpper && hasLower && hasNumber && hasSpecial) {
		return service_errors.NewErrBadRequest("Password must include at least one uppercase letter, one lowercase letter, one digit, and one special character")
	}

	return nil
}
