package handler

import (
	"context"
	"fmt"
	"gin-boilerplate/comm/cache"
	"gin-boilerplate/comm/db"
	"gin-boilerplate/comm/errors"
	"gin-boilerplate/comm/http"
	"gin-boilerplate/comm/logger"
	"gin-boilerplate/models"
	"gin-boilerplate/types"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jordan-wright/email"
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
			http.Fail(ctx, http.MsgOption("The account or password is incorrect"))
			return
		}
		http.Fail(ctx, http.MsgOption(err.Error()))
		return
	}
	if ok, err := user.CompareHashAndPassword([]byte(loginForm.Password)); !ok || err != nil {
		if !ok {
			http.Fail(ctx, http.MsgOption("The account or password is incorrect"))
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
	if err := user.GenerateFromPassword(bcrypt.DefaultCost); err != nil {
		http.Fail(ctx, http.MsgOption(err.Error()))
		return
	}
	if err := s.InsertUserDB(ctx, session, &user); err != nil {
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
	if err := user.GenerateFromPassword(bcrypt.DefaultCost); err != nil {
		http.Fail(ctx, http.MsgOption(err.Error()))
		return
	}

	if err := s.UpdateUserDB(ctx, session, &user); err != nil {
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
	err := s.sendVerificationEmail(ctx.Request.Context(), sendVerificationEmailRequestForm.FromName, sendVerificationEmailRequestForm.Email, user.Name, sendVerificationEmailRequestForm.Subject, sendVerificationEmailRequestForm.TextContent, token, sendVerificationEmailRequestForm.RedirectUrl, sendVerificationEmailRequestForm.FailureRedirectUrl)
	if err != nil {
		logger.Error(ctx.Request.Context(), err)
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
	if err := s.UpdateUserDB(ctx, session, &user); err != nil {
		http.Fail(ctx, http.MsgOption(err.Error()))
		return
	}
	http.Success(ctx, http.MsgOption("VerifyEmail succeeded"))
}

//SendPasswordResetEmail...
func (s *Handler) SendPasswordResetEmail(ctx *gin.Context) {
	var session = db.GetDB()
	var sendPasswordResetEmailForm types.SendPasswordResetEmailForm
	if err := ctx.ShouldBindJSON(&sendPasswordResetEmailForm); err != nil {
		http.Fail(ctx, http.MsgOption(sendPasswordResetEmailForm.SendPasswordResetEmail(err)))
		return
	}
	where := models.User{}
	where.Email = sendPasswordResetEmailForm.Email
	user := models.User{}
	if err := s.QueryUserDetailDB(ctx, session, &where, &user); err != nil {
		if errors.Is(err, errors.ERecordNotFound) {
			http.Fail(ctx, http.MsgOption("The account was not found"))
			return
		}
		http.Fail(ctx, http.MsgOption(err.Error()))
		return
	}

	var expiry int64 = 1800 // 1800 secs = 30 min
	if sendPasswordResetEmailForm.Expiration > 0 {
		expiry = sendPasswordResetEmailForm.Expiration
	}

	code := shortid.MustGenerate()

	// save the password reset code
	_, err := s.savePasswordResetCode(ctx, user.ID, code, time.Duration(expiry)*time.Second)
	if err != nil {
		http.Fail(ctx, http.MsgOption(err.Error()))
		return
	}
	// save the code in the database and then send via email
	err = s.sendPasswordResetEmail(ctx, user.ID, code, sendPasswordResetEmailForm.FromName, sendPasswordResetEmailForm.Email, user.Name, sendPasswordResetEmailForm.Subject, sendPasswordResetEmailForm.TextContent)
	if err != nil {
		http.Fail(ctx, http.MsgOption(err.Error()))
		return
	}
	http.Success(ctx, http.MsgOption("SendPasswordResetEmail succeeded"))
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

	//check code
	if _, err := s.readPasswordResetCode(ctx, user.ID, resetPasswordRequestForm.Code); err != nil {
		http.Fail(ctx, http.MsgOption(err.Error()))
		return
	}

	//update password
	user.Password = resetPasswordRequestForm.NewPassword
	if err := user.GenerateFromPassword(bcrypt.DefaultCost); err != nil {
		http.Fail(ctx, http.MsgOption(err.Error()))
		return
	}

	if err := s.UpdateUserDB(ctx, session, &user); err != nil {
		http.Fail(ctx, http.MsgOption(err.Error()))
		return
	}

	http.Success(ctx, http.MsgOption("ResetPassword succeeded"))
}

//sendVerificationEmail...
func (s *Handler) sendVerificationEmail(ctx context.Context, fromName, toAddress, toUsername, subject, textContent, token, redirctUrl, failureRedirectUrl string) error {
	uri := "https://baidu.com"
	query := "?token=" + token + "&redirectUrl=" + url.QueryEscape(redirctUrl) + "&failureRedirectUrl=" + url.QueryEscape(failureRedirectUrl)
	textContent = strings.Replace(textContent, "$verification_link", uri+query, -1)
	email := email.Email{
		Subject: subject,
		From:    fromName,
		To:      []string{toAddress},
		Text:    []byte(textContent),
	}
	err := s.sendEmail(ctx, &email)
	fmt.Println(email)
	return err
}

//SendPasswordResetEmail...
func (s *Handler) sendPasswordResetEmail(ctx context.Context, userId uint, codeStr, fromName, toAddress, toUsername, subject, textContent string) error {
	textContent = strings.Replace(textContent, "$code", codeStr, -1)
	email := email.Email{
		Subject: subject,
		From:    fromName,
		To:      []string{toAddress},
		Text:    []byte(textContent),
	}
	err := s.sendEmail(ctx, &email)
	return err
}

func (s *Handler) generatePasswordResetCodeStoreKey(userId uint, code string) string {
	return fmt.Sprintf("user/password-reset-codes/%v-%v", userId, code)
}

//savePasswordResetCode...
func (s *Handler) savePasswordResetCode(ctx context.Context, userId uint, code string, expiry time.Duration) (*types.PasswordResetCode, error) {
	pwcode := types.PasswordResetCode{
		Expires: time.Now().Add(expiry),
		UserID:  userId,
		Code:    code,
	}

	if err := s.Cache.Set(s.generatePasswordResetCodeStoreKey(userId, code), pwcode, expiry); err != nil {
		return nil, err
	}

	return &pwcode, nil
}

// readPasswordResetCode returns the user reset code
func (s *Handler) readPasswordResetCode(ctx context.Context, userId uint, code string) (*types.PasswordResetCode, error) {
	pwcode := types.PasswordResetCode{}
	err := s.Cache.Get(s.generatePasswordResetCodeStoreKey(userId, code), &pwcode)
	if err != nil && err != cache.ErrCacheMiss {
		return nil, err
	}

	// check the expiry
	if pwcode.Expires.Before(time.Now()) {
		return nil, errors.New(errors.ENone, "password reset code expired")
	}

	return &pwcode, nil
}
