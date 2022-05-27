package types

import (
	"encoding/json"

	"github.com/go-playground/validator/v10"
)

// Send an email with a verification code to reset password.
// Call "ResetPassword" endpoint once user provides the code.
type SendPasswordResetEmailForm struct {
	UserValidator
	// email address to send reset for
	Email string `form:"email" json:"email" binding:"required"`
	// subject of the email
	Subject string `form:"subject" json:"subject" binding:"required"`
	// Text content of the email. Don't forget to include the string '$code' which will be replaced by the real verification link
	// HTML emails are not available currently.
	TextContent string `form:"text_content" json:"text_content" binding:"required"`
	// Display name of the sender for the email. Note: the email address will still be 'noreply@email.m3ocontent.com'
	FromName string `form:"from_name" json:"from_name" binding:"required"`
	// Number of secs that the password reset email is valid for, defaults to 1800 secs (30 mins)
	Expiration int64 `form:"expiration" json:"expiration" binding:"required"`
}

// Reset password with the code sent by the "SendPasswordResetEmail" endpoint.
type ResetPasswordRequestForm struct {
	UserValidator
	// the email to reset the password for
	Email string `form:"email" json:"email" binding:"required"`
	// The code from the verification email
	Code string `form:"code" json:"code" binding:"required"`
	// the new password
	NewPassword string `form:"new_password" json:"new_password" binding:"required,min=6,max=12"`
	// confirm new password
	ConfirmPassword string `form:"confirm_password" json:"confirm_password" binding:"required,min=6,max=12"`
}

//SendVerificationEmailRequestForm...
type SendVerificationEmailRequestForm struct {
	UserValidator
	// email address to send the verification code
	Email string `form:"email" json:"email,omitempty" binding:"required"`
	// subject of the email
	Subject string `form:"subject" json:"subject,omitempty" binding:"required"`
	// Text content of the email. Don't forget to include the string '$micro_verification_link' which will be replaced by the real verification link
	// HTML emails are not available currently.
	TextContent string `form:"text_content" json:"text_content,omitempty" binding:"required"`
	// The url to redirect to after successful verification
	RedirectUrl string `form:"redirect_url" json:"redirect_url,omitempty" binding:"required"`
	// The url to redirect to incase of failure
	FailureRedirectUrl string `form:"failure_redirect_url" json:"failure_redirect_url,omitempty" binding:"required"`
	// Display name of the sender for the email. Note: the email address will still be 'noreply@email.m3ocontent.com'
	FromName string `form:"from_name" json:"from_name,omitempty" binding:"required"`
}

//VerifyEmailRequestForm ...
type VerifyEmailRequestForm struct {
	UserValidator
	// the token
	Token string `form:"token" json:"token" binding:"required"`
}

//LoginForm ...
type UpdatePasswordForm struct {
	UserValidator
	// the account id
	UserId uint `form:"user_id" json:"user_id" binding:"required"`
	// the old password
	OldPassword string `form:"old_password" json:"old_password" binding:"required"`
	// the new password
	NewPassword string `form:"new_password" json:"new_password" binding:"required,min=6,max=12"`
	// confirm new password
	ConfirmPassword string `form:"confirm_password" json:"confirm_password" binding:"required,min=6,max=12"`
}

//LoginForm ...
type LoginForm struct {
	UserValidator
	Email    string `form:"email" json:"email" binding:"required,email"`
	Password string `form:"password" json:"password" binding:"required,min=3,max=50"`
}

//RegisterForm ...
type RegisterForm struct {
	UserValidator
	Name     string `form:"name" json:"name" binding:"required,min=3,max=20,fullName"` //fullName rule is in validator.go
	Email    string `form:"email" json:"email" binding:"required,email"`
	Password string `form:"password" json:"password" binding:"required,min=3,max=50"`
}

type UserProfile struct {
	Name   string `form:"name" json:"name"`
	Email  string `form:"email" json:"email"`
	Avatar string `form:"avatar" json:"avatar"`
}

//UserValidator ...
type UserValidator struct{}

//Token ...
func (f UserValidator) Token(tag string, errMsg ...string) (message string) {
	switch tag {
	case "required":
		if len(errMsg) == 0 {
			return "Please enter your name"
		}
		return errMsg[0]
	case "min", "max":
		return "Your name should be between 3 to 20 characters"
	case "fullName":
		return "Name should not include any special characters or numbers"
	default:
		return "Something went wrong, please try again later"
	}
}

//Name ...
func (f UserValidator) Name(tag string, errMsg ...string) (message string) {
	switch tag {
	case "required":
		if len(errMsg) == 0 {
			return "Please enter your name"
		}
		return errMsg[0]
	case "min", "max":
		return "Your name should be between 3 to 20 characters"
	case "fullName":
		return "Name should not include any special characters or numbers"
	default:
		return "Something went wrong, please try again later"
	}
}

//Email ...
func (f UserValidator) Email(tag string, errMsg ...string) (message string) {
	switch tag {
	case "required":
		if len(errMsg) == 0 {
			return "Please enter your email"
		}
		return errMsg[0]
	case "min", "max", "email":
		return "Please enter a valid email"
	default:
		return "Something went wrong, please try again later"
	}
}

//Password ...
func (f UserValidator) Password(tag string) (message string) {
	switch tag {
	case "required":
		return "Please enter your password"
	case "min", "max":
		return "Your password should be between 3 and 50 characters"
	case "eqfield":
		return "Your passwords does not match"
	default:
		return "Something went wrong, please try again later"
	}
}

//OldPassword ...
func (f UserValidator) OldPassword(tag string) (message string) {
	switch tag {
	case "required":
		return "Please enter your old password"
	case "min", "max":
		return "Your password should be between 3 and 50 characters"
	case "eqfield":
		return "Your passwords does not match"
	default:
		return "Something went wrong, please try again later"
	}
}

//NewPassword ...
func (f UserValidator) NewPassword(tag string) (message string) {
	switch tag {
	case "required":
		return "Please enter your password"
	case "min", "max":
		return "Your password should be between 3 and 50 characters"
	case "eqfield":
		return "Your passwords does not match"
	default:
		return "Something went wrong, please try again later"
	}
}

//UserId ...
func (f UserValidator) UserId(tag string) (message string) {
	switch tag {
	case "required":
		return "Please enter your userid"
	default:
		return "Something went wrong, please try again later"
	}
}

//Signin ...
func (f UserValidator) Login(err error) string {
	switch err.(type) {
	case validator.ValidationErrors:
		if _, ok := err.(*json.UnmarshalTypeError); ok {
			return "Something went wrong, please try again later"
		}
		for _, err := range err.(validator.ValidationErrors) {
			if err.Field() == "Email" {
				return f.Email(err.Tag())
			}
			if err.Field() == "Password" {
				return f.Password(err.Tag())
			}
		}
	default:
		return "Invalid request"
	}
	return "Something went wrong, please try again later"
}

//Register ...
func (f UserValidator) Register(err error) string {
	switch err.(type) {
	case validator.ValidationErrors:
		if _, ok := err.(*json.UnmarshalTypeError); ok {
			return "Something went wrong, please try again later"
		}
		for _, err := range err.(validator.ValidationErrors) {
			if err.Field() == "Name" {
				return f.Name(err.Tag())
			}
			if err.Field() == "Email" {
				return f.Email(err.Tag())
			}
			if err.Field() == "Password" {
				return f.Password(err.Tag())
			}
		}
	default:
		return "Invalid request"
	}
	return "Something went wrong, please try again later"
}

//Register ...
func (f UserValidator) UpdatePassword(err error) string {
	switch err.(type) {
	case validator.ValidationErrors:
		if _, ok := err.(*json.UnmarshalTypeError); ok {
			return "Something went wrong, please try again later"
		}
		for _, err := range err.(validator.ValidationErrors) {
			if err.Field() == "NewPassword" {
				return f.NewPassword(err.Tag())
			}
			if err.Field() == "OldPassword" {
				return f.OldPassword(err.Tag())
			}
			if err.Field() == "UserId" {
				return f.UserId(err.Tag())
			}
		}
	default:
		return "Invalid request"
	}
	return "Something went wrong, please try again later"
}

//VerifyEmail ...
func (f UserValidator) VerifyEmail(err error) string {
	switch err.(type) {
	case validator.ValidationErrors:
		if _, ok := err.(*json.UnmarshalTypeError); ok {
			return "Something went wrong, please try again later"
		}
		for _, err := range err.(validator.ValidationErrors) {
			if err.Field() == "Token" {
				return f.Token(err.Tag())
			}
		}
	default:
		return "Invalid request"
	}
	return "Something went wrong, please try again later"
}

//SendVerificationEmail ...
func (f UserValidator) SendVerificationEmail(err error) string {
	switch err.(type) {
	case validator.ValidationErrors:
		if _, ok := err.(*json.UnmarshalTypeError); ok {
			return "Something went wrong, please try again later"
		}
		for _, err := range err.(validator.ValidationErrors) {
			if err.Field() == "Token" {
				return f.Token(err.Tag())
			}
		}
	default:
		return "Invalid request"
	}
	return "Something went wrong, please try again later"
}

//ResetPassword ...
func (f UserValidator) ResetPassword(err error) string {
	switch err.(type) {
	case validator.ValidationErrors:
		if _, ok := err.(*json.UnmarshalTypeError); ok {
			return "Something went wrong, please try again later"
		}
		for _, err := range err.(validator.ValidationErrors) {
			if err.Field() == "Token" {
				return f.Token(err.Tag())
			}
		}
	default:
		return "Invalid request"
	}
	return "Something went wrong, please try again later"
}

//ResetPassword ...
func (f UserValidator) SendPasswordResetEmail(err error) string {
	switch err.(type) {
	case validator.ValidationErrors:
		if _, ok := err.(*json.UnmarshalTypeError); ok {
			return "Something went wrong, please try again later"
		}
		for _, err := range err.(validator.ValidationErrors) {
			if err.Field() == "Token" {
				return f.Token(err.Tag())
			}
		}
	default:
		return "Invalid request"
	}
	return "Something went wrong, please try again later"
}
