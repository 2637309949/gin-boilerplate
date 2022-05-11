package handler

import (
	"fmt"
	"gin-boilerplate/comm/db"
	"gin-boilerplate/comm/errors"
	"gin-boilerplate/comm/http"
	"gin-boilerplate/models"
	"gin-boilerplate/types"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/teris-io/shortid"
	"golang.org/x/crypto/bcrypt"
)

//Login ...
func (s *Handler) Login(ctx *gin.Context) {
	var session = db.GetDB()
	var token types.Token
	var loginForm types.LoginForm
	if err := ctx.ShouldBindJSON(&loginForm); err != nil {
		http.Fail(ctx, http.MsgOption(loginForm.Login(err)))
		return
	}

	//Check if exit
	where := models.User{
		Email: loginForm.Email,
	}
	user := models.User{}
	if err := s.QueryUserDetailDB(ctx, session, &where, &user); err != nil {
		if errors.Is(err, errors.ERecordNotFound) {
			http.Fail(ctx, http.MsgOption("账号或密码错误!"))
			return
		}
		http.Fail(ctx, http.MsgOption(err.Error()))
		return
	}
	if ok, err := user.CompareHashAndPassword([]byte(loginForm.Password)); !ok || err != nil {
		if !ok {
			http.Fail(ctx, http.MsgOption("账号或密码错误!"))
			return
		}
		http.Fail(ctx, http.MsgOption(err.Error()))
		return
	}

	//Generate the JWT auth token
	tokenDetails, err := s.CreateToken(user.ID)
	if err != nil {
		http.Fail(ctx, http.MsgOption(err.Error()))
		return
	}

	err = s.CreateAuth(user.ID, tokenDetails)
	if err == nil {
		token.AccessToken = tokenDetails.AccessToken
		token.RefreshToken = tokenDetails.RefreshToken
	}

	tk := gin.H{
		"user":  user,
		"token": token,
	}
	http.Success(ctx, http.FlatOption(tk))
}

//Register ...
func (s *Handler) Register(ctx *gin.Context) {
	var session = db.GetDB()
	var registerForm types.RegisterForm
	if validationErr := ctx.ShouldBindJSON(&registerForm); validationErr != nil {
		http.Fail(ctx, http.MsgOption(registerForm.Register(validationErr)))
		return
	}

	//Check if exit
	where := models.User{
		Email: registerForm.Email,
	}
	user := models.User{
		Name:     registerForm.Name,
		Email:    registerForm.Email,
		Password: registerForm.Password,
	}
	if err := s.QueryUserDetailDB(ctx, session, &where, &user); err == nil {
		http.Fail(ctx, http.MsgOption("Email already exists"))
		return
	}
	if err := user.GenerateFromPassword(bcrypt.DefaultCost); err == nil {
		http.Fail(ctx, http.MsgOption(err.Error()))
		return
	}
	if err := s.InsertUserDB(ctx, session, &user); err == nil {
		http.Fail(ctx, http.MsgOption(err.Error()))
		return
	}
	user.Name = registerForm.Name
	user.Email = registerForm.Email
	http.Success(ctx, http.DataOption(user))
}

//Logout ...
func (s *Handler) Logout(ctx *gin.Context) {
	au, err := s.ExtractTokenMetadata(ctx.Request)
	if err != nil {
		http.Fail(ctx, http.MsgOption("User not logged in"))
		return
	}

	err = s.DeleteAuth(au.AccessUUID)
	if err != nil { //if any goes wrong
		http.Fail(ctx, http.MsgOption("Invalid request"))
		return
	}
	http.Success(ctx, http.MsgOption("User not logged in"))
}

//UpdatePassword...
func (s *Handler) UpdatePassword(ctx *gin.Context) {
	var session = db.GetDB()
	var updatePasswordForm types.UpdatePasswordForm
	if err := ctx.ShouldBindJSON(&updatePasswordForm); err != nil {
		http.Fail(ctx, http.MsgOption(updatePasswordForm.Login(err)))
		return
	}

	if updatePasswordForm.NewPassword != updatePasswordForm.ConfirmPassword {
		http.Fail(ctx, http.MsgOption("Passwords don't match"))
		return
	}

	where := models.User{}
	where.ID = updatePasswordForm.UserId
	user := models.User{}
	if err := s.QueryUserDetailDB(ctx, session, &where, &user); err != nil {
		if errors.Is(err, errors.ERecordNotFound) {
			http.Fail(ctx, http.MsgOption("The account was not found"))
			return
		}
		http.Fail(ctx, http.MsgOption(err.Error()))
		return
	}

	if ok, err := user.CompareHashAndPassword([]byte(updatePasswordForm.OldPassword)); !ok || err != nil {
		if !ok {
			http.Fail(ctx, http.MsgOption("Wrong original password"))
			return
		}
		http.Fail(ctx, http.MsgOption(err.Error()))
		return
	}
	user.Password = updatePasswordForm.NewPassword
	if err := user.GenerateFromPassword(bcrypt.DefaultCost); err == nil {
		http.Fail(ctx, http.MsgOption(err.Error()))
		return
	}

	if err := s.UpdateUserDB(ctx, session, &user); err == nil {
		http.Fail(ctx, http.MsgOption(err.Error()))
		return
	}

	http.Success(ctx, http.MsgOption("Update succeeded"))
}

func (s *Handler) generateVerificationTokenStoreKey(token string) string {
	return fmt.Sprintf("user/verification-token/%s", token)
}

//SendVerificationEmail...
func (s *Handler) SendVerificationEmail(ctx *gin.Context) {
	var session = db.GetDB()
	var sendVerificationEmailRequestForm types.SendVerificationEmailRequestForm
	if err := ctx.ShouldBindJSON(&sendVerificationEmailRequestForm); err != nil {
		http.Fail(ctx, http.MsgOption(sendVerificationEmailRequestForm.SendVerificationEmail(err)))
		return
	}

	where := models.User{}
	where.Email = sendVerificationEmailRequestForm.Email
	user := models.User{}
	if err := s.QueryUserDetailDB(ctx, session, &where, &user); err != nil {
		if errors.Is(err, errors.ERecordNotFound) {
			http.Fail(ctx, http.MsgOption("The account was not found"))
			return
		}
		http.Fail(ctx, http.MsgOption(err.Error()))
		return
	}

	token := shortid.MustGenerate()
	if err := s.Cache.Set(s.generateVerificationTokenStoreKey(token), user.ID, 10*time.Minute); err != nil {
		http.Fail(ctx, http.MsgOption(err.Error()))
		return
	}

	//send email
	err := s.SendEmail(sendVerificationEmailRequestForm.FromName, sendVerificationEmailRequestForm.Email, user.Name, sendVerificationEmailRequestForm.Subject, sendVerificationEmailRequestForm.TextContent, token, sendVerificationEmailRequestForm.RedirectUrl, sendVerificationEmailRequestForm.FailureRedirectUrl)
	if err != nil {
		http.Fail(ctx, http.MsgOption(err.Error()))
		return
	}

	http.Success(ctx, http.MsgOption("SendVerificationEmail succeeded"))
}

//VerifyEmail...
func (s *Handler) VerifyEmail(ctx *gin.Context) {
	var session = db.GetDB()
	var verifyEmailRequestForm types.VerifyEmailRequestForm
	if err := ctx.ShouldBindJSON(&verifyEmailRequestForm); err != nil {
		http.Fail(ctx, http.MsgOption(verifyEmailRequestForm.VerifyEmail(err)))
		return
	}

	userId := uint(0)
	if err := s.Cache.Get(s.generateVerificationTokenStoreKey(verifyEmailRequestForm.Token), &userId); err != nil {
		http.Fail(ctx, http.MsgOption(err.Error()))
		return
	}

	where := models.User{}
	where.ID = userId
	user := models.User{}
	if err := s.QueryUserDetailDB(ctx, session, &where, &user); err != nil {
		if errors.Is(err, errors.ERecordNotFound) {
			http.Fail(ctx, http.MsgOption("The account was not found"))
			return
		}
		http.Fail(ctx, http.MsgOption(err.Error()))
		return
	}
	user.Verified = 1
	if err := s.UpdateUserDB(ctx, session, &user); err == nil {
		http.Fail(ctx, http.MsgOption(err.Error()))
		return
	}
	http.Success(ctx, http.MsgOption("VerifyEmail succeeded"))
}

//SendPasswordResetEmail...
func (s *Handler) SendPasswordResetEmail(ctx *gin.Context) {

}

//ResetPassword...
func (s *Handler) ResetPassword(ctx *gin.Context) {
	var session = db.GetDB()
	var resetPasswordRequestForm types.ResetPasswordRequestForm
	if err := ctx.ShouldBindJSON(&resetPasswordRequestForm); err != nil {
		http.Fail(ctx, http.MsgOption(resetPasswordRequestForm.ResetPassword(err)))
		return
	}
	where := models.User{}
	where.Email = resetPasswordRequestForm.Email
	user := models.User{}
	if err := s.QueryUserDetailDB(ctx, session, &where, &user); err != nil {
		if errors.Is(err, errors.ERecordNotFound) {
			http.Fail(ctx, http.MsgOption("The account was not found"))
			return
		}
		http.Fail(ctx, http.MsgOption(err.Error()))
		return
	}
	//check reset code

	//update password
	user.Password = resetPasswordRequestForm.NewPassword
	if err := user.GenerateFromPassword(bcrypt.DefaultCost); err == nil {
		http.Fail(ctx, http.MsgOption(err.Error()))
		return
	}

	if err := s.UpdateUserDB(ctx, session, &user); err == nil {
		http.Fail(ctx, http.MsgOption(err.Error()))
		return
	}

	http.Success(ctx, http.MsgOption("ResetPassword succeeded"))
}
