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
)

// ref: https://swaggo.github.io/swaggo.io/declarative_comments_format/api_operation.html
// @Summary Login
// @Description user login
// @Tags users
// @Accept  json
// @Produce  json
// @Router /api/v1/login [POST]
func (h *Handler) Login(ctx *gin.Context) {
	var session = db.GetDB()
	var token types.Token
	var loginForm types.LoginForm
	if err := ctx.ShouldBindJSON(&loginForm); err != nil {
		logger.Error(ctx.Request.Context(), err)
		http.Fail(ctx, http.MsgOption(loginForm.Login(err)))
		return
	}

	//Check if exit
	where := models.User{
		Email: loginForm.Email,
	}
	user := models.User{}
	if err := h.QueryUserDetailDB(ctx, session, &where, &user); err != nil {
		if errors.Is(err, errors.ERecordNotFound) {
			http.Fail(ctx, http.MsgOption("The account or password is incorrect"))
			return
		}
		logger.Error(ctx.Request.Context(), err)
		http.Fail(ctx, http.MsgOption(err.Error()))
		return
	}
	if ok := user.CompareHashAndPassword(loginForm.Password); !ok {
		http.Fail(ctx, http.MsgOption("The account or password is incorrect"))
		return
	}

	//Generate the JWT auth token
	tokenDetails, err := h.CreateToken(user.ID)
	if err != nil {
		logger.Error(ctx.Request.Context(), err)
		http.Fail(ctx, http.MsgOption(err.Error()))
		return
	}

	err = h.CreateAuth(user.ID, tokenDetails)
	if err == nil {
		token.AccessToken = tokenDetails.AccessToken
		token.RefreshToken = tokenDetails.RefreshToken
	}

	profile := types.UserProfile{
		Name:   user.Name,
		Email:  user.Email,
		Avatar: user.Avatar,
	}

	tk := gin.H{
		"access_token":  token.AccessToken,
		"refresh_token": token.RefreshToken,
		"profile":       profile,
	}
	http.Success(ctx, http.FlatOption(tk))
}

// ref: https://swaggo.github.io/swaggo.io/declarative_comments_format/api_operation.html
// @Summary Register
// @Description user register
// @Tags users
// @Accept  json
// @Produce  json
// @Router /api/v1/register [POST]
func (h *Handler) Register(ctx *gin.Context) {
	var session = db.GetDB()
	var registerForm types.RegisterForm
	if err := ctx.ShouldBindJSON(&registerForm); err != nil {
		logger.Error(ctx.Request.Context(), err)
		http.Fail(ctx, http.MsgOption(registerForm.Register(err)))
		return
	}

	//Check if exit
	where := models.User{
		Email: registerForm.Email,
	}
	user := models.User{
		Name:  registerForm.Name,
		Email: registerForm.Email,
	}
	if err := h.QueryUserDetailDB(ctx, session, &where, &user); err == nil {
		http.Fail(ctx, http.MsgOption("Email already exists"))
		return
	}
	password, err := user.GenerateFromPassword(registerForm.Password)
	if err != nil {
		logger.Error(ctx.Request.Context(), err)
		http.Fail(ctx, http.MsgOption(err.Error()))
		return
	}
	user.Password = password
	if err := h.InsertUserDB(ctx, session, &user); err != nil {
		logger.Error(ctx.Request.Context(), err)
		http.Fail(ctx, http.MsgOption(err.Error()))
		return
	}
	user.Name = registerForm.Name
	user.Email = registerForm.Email
	http.Success(ctx, http.DataOption(user))
}

// ref: https://swaggo.github.io/swaggo.io/declarative_comments_format/api_operation.html
// @Summary Logout
// @Description user logout
// @Tags users
// @Accept  json
// @Produce  json
// @Router /api/v1/logout [POST]
func (h *Handler) Logout(ctx *gin.Context) {
	au, err := h.ExtractTokenMetadata(ctx.Request)
	if err != nil {
		logger.Error(ctx.Request.Context(), err)
		http.Fail(ctx, http.MsgOption("User not logged in"))
		return
	}

	err = h.DeleteAuth(au.AccessUUID)
	if err != nil { //if any goes wrong
		logger.Error(ctx.Request.Context(), err)
		http.Fail(ctx, http.MsgOption("Invalid request"))
		return
	}
	http.Success(ctx, http.MsgOption("User not logged in"))
}

// ref: https://swaggo.github.io/swaggo.io/declarative_comments_format/api_operation.html
// @Summary Logout
// @Description user updatePassword
// @Tags users
// @Accept  json
// @Produce  json
// @Router /api/v1/updatePassword [POST]
func (h *Handler) UpdatePassword(ctx *gin.Context) {
	var session = db.GetDB()
	var updatePasswordForm types.UpdatePasswordForm
	if err := ctx.ShouldBindJSON(&updatePasswordForm); err != nil {
		logger.Error(ctx.Request.Context(), err)
		http.Fail(ctx, http.MsgOption(updatePasswordForm.Login(err)))
		return
	}

	if updatePasswordForm.NewPassword != updatePasswordForm.ConfirmPassword {
		logger.Errorf(ctx.Request.Context(), "Passwords don't match")
		http.Fail(ctx, http.MsgOption("Passwords don't match"))
		return
	}

	where := models.User{}
	where.ID = updatePasswordForm.UserId
	user := models.User{}
	if err := h.QueryUserDetailDB(ctx, session, &where, &user); err != nil {
		if errors.Is(err, errors.ERecordNotFound) {
			logger.Errorf(ctx.Request.Context(), "The account was not found")
			http.Fail(ctx, http.MsgOption("The account was not found"))
			return
		}
		http.Fail(ctx, http.MsgOption(err.Error()))
		return
	}

	if ok := user.CompareHashAndPassword(updatePasswordForm.OldPassword); !ok {
		logger.Errorf(ctx.Request.Context(), "Wrong original password")
		http.Fail(ctx, http.MsgOption("Wrong original password"))
		return
	}

	password, err := user.GenerateFromPassword(updatePasswordForm.NewPassword)
	if err != nil {
		logger.Error(ctx.Request.Context(), err)
		http.Fail(ctx, http.MsgOption(err.Error()))
		return
	}
	user.Password = password
	if err := h.UpdateUserDB(ctx, session, &user); err != nil {
		logger.Error(ctx.Request.Context(), err)
		http.Fail(ctx, http.MsgOption(err.Error()))
		return
	}

	http.Success(ctx, http.MsgOption("Update succeeded"))
}

func (h *Handler) generateVerificationTokenStoreKey(token string) string {
	return fmt.Sprintf("user/verification-token/%s", token)
}

// ref: https://swaggo.github.io/swaggo.io/declarative_comments_format/api_operation.html
// @Summary SendVerificationEmail
// @Description send verification email
// @Tags users
// @Accept  json
// @Produce  json
// @Router /api/v1/sendVerificationEmail [POST]
func (h *Handler) SendVerificationEmail(ctx *gin.Context) {
	var session = db.GetDB()
	var sendVerificationEmailRequestForm types.SendVerificationEmailRequestForm
	if err := ctx.ShouldBindJSON(&sendVerificationEmailRequestForm); err != nil {
		http.Fail(ctx, http.MsgOption(sendVerificationEmailRequestForm.SendVerificationEmail(err)))
		return
	}

	where := models.User{}
	where.Email = sendVerificationEmailRequestForm.Email
	user := models.User{}
	if err := h.QueryUserDetailDB(ctx, session, &where, &user); err != nil {
		if errors.Is(err, errors.ERecordNotFound) {
			http.Fail(ctx, http.MsgOption("The account was not found"))
			return
		}
		logger.Error(ctx.Request.Context(), err)
		http.Fail(ctx, http.MsgOption(err.Error()))
		return
	}

	token := shortid.MustGenerate()
	if err := h.Cache.Set(h.generateVerificationTokenStoreKey(token), user.ID, 10*time.Minute); err != nil {
		logger.Error(ctx.Request.Context(), err)
		http.Fail(ctx, http.MsgOption(err.Error()))
		return
	}

	//send email
	err := h.sendVerificationEmail(ctx.Request.Context(), sendVerificationEmailRequestForm.FromName, sendVerificationEmailRequestForm.Email, user.Name, sendVerificationEmailRequestForm.Subject, sendVerificationEmailRequestForm.TextContent, token, sendVerificationEmailRequestForm.RedirectUrl, sendVerificationEmailRequestForm.FailureRedirectUrl)
	if err != nil {
		logger.Error(ctx.Request.Context(), err)
		http.Fail(ctx, http.MsgOption(err.Error()))
		return
	}

	http.Success(ctx, http.MsgOption("SendVerificationEmail succeeded"))
}

// ref: https://swaggo.github.io/swaggo.io/declarative_comments_format/api_operation.html
// @Summary VerifyEmail
// @Description verify email
// @Tags users
// @Accept  json
// @Produce  json
// @Router /api/v1/verifyEmail [POST]
func (h *Handler) VerifyEmail(ctx *gin.Context) {
	var session = db.GetDB()
	var verifyEmailRequestForm types.VerifyEmailRequestForm
	if err := ctx.ShouldBindJSON(&verifyEmailRequestForm); err != nil {
		logger.Error(ctx.Request.Context(), err)
		http.Fail(ctx, http.MsgOption(verifyEmailRequestForm.VerifyEmail(err)))
		return
	}

	userId := uint(0)
	if err := h.Cache.Get(h.generateVerificationTokenStoreKey(verifyEmailRequestForm.Token), &userId); err != nil {
		http.Fail(ctx, http.MsgOption(err.Error()))
		return
	}

	where := models.User{}
	where.ID = userId
	user := models.User{}
	if err := h.QueryUserDetailDB(ctx, session, &where, &user); err != nil {
		if errors.Is(err, errors.ERecordNotFound) {
			http.Fail(ctx, http.MsgOption("The account was not found"))
			return
		}
		logger.Error(ctx.Request.Context(), err)
		http.Fail(ctx, http.MsgOption(err.Error()))
		return
	}
	user.Verified = 1
	if err := h.UpdateUserDB(ctx, session, &user); err != nil {
		logger.Error(ctx.Request.Context(), err)
		http.Fail(ctx, http.MsgOption(err.Error()))
		return
	}
	http.Success(ctx, http.MsgOption("VerifyEmail succeeded"))
}

// ref: https://swaggo.github.io/swaggo.io/declarative_comments_format/api_operation.html
// @Summary SendPasswordResetEmail
// @Description send password reset email
// @Tags users
// @Accept  json
// @Produce  json
// @Router /api/v1/sendPasswordResetEmail [POST]
func (h *Handler) SendPasswordResetEmail(ctx *gin.Context) {
	var session = db.GetDB()
	var sendPasswordResetEmailForm types.SendPasswordResetEmailForm
	if err := ctx.ShouldBindJSON(&sendPasswordResetEmailForm); err != nil {
		http.Fail(ctx, http.MsgOption(sendPasswordResetEmailForm.SendPasswordResetEmail(err)))
		return
	}
	where := models.User{}
	where.Email = sendPasswordResetEmailForm.Email
	user := models.User{}
	if err := h.QueryUserDetailDB(ctx, session, &where, &user); err != nil {
		if errors.Is(err, errors.ERecordNotFound) {
			http.Fail(ctx, http.MsgOption("The account was not found"))
			return
		}
		logger.Error(ctx.Request.Context(), err)
		http.Fail(ctx, http.MsgOption(err.Error()))
		return
	}

	var expiry int64 = 1800 // 1800 secs = 30 min
	if sendPasswordResetEmailForm.Expiration > 0 {
		expiry = sendPasswordResetEmailForm.Expiration
	}

	code := shortid.MustGenerate()

	// save the password reset code
	_, err := h.savePasswordResetCode(ctx, user.ID, code, time.Duration(expiry)*time.Second)
	if err != nil {
		logger.Error(ctx.Request.Context(), err)
		http.Fail(ctx, http.MsgOption(err.Error()))
		return
	}
	// save the code in the database and then send via email
	err = h.sendPasswordResetEmail(ctx, user.ID, code, sendPasswordResetEmailForm.FromName, sendPasswordResetEmailForm.Email, user.Name, sendPasswordResetEmailForm.Subject, sendPasswordResetEmailForm.TextContent)
	if err != nil {
		logger.Error(ctx.Request.Context(), err)
		http.Fail(ctx, http.MsgOption(err.Error()))
		return
	}
	http.Success(ctx, http.MsgOption("SendPasswordResetEmail succeeded"))
}

// ref: https://swaggo.github.io/swaggo.io/declarative_comments_format/api_operation.html
// @Summary ResetPassword
// @Description reset password
// @Tags users
// @Accept  json
// @Produce  json
// @Router /api/v1/resetPassword [POST]
func (h *Handler) ResetPassword(ctx *gin.Context) {
	var session = db.GetDB()
	var resetPasswordRequestForm types.ResetPasswordRequestForm
	if err := ctx.ShouldBindJSON(&resetPasswordRequestForm); err != nil {
		logger.Error(ctx.Request.Context(), err)
		http.Fail(ctx, http.MsgOption(resetPasswordRequestForm.ResetPassword(err)))
		return
	}
	where := models.User{}
	where.Email = resetPasswordRequestForm.Email
	user := models.User{}
	if err := h.QueryUserDetailDB(ctx, session, &where, &user); err != nil {
		if errors.Is(err, errors.ERecordNotFound) {
			http.Fail(ctx, http.MsgOption("The account was not found"))
			return
		}
		logger.Error(ctx.Request.Context(), err)
		http.Fail(ctx, http.MsgOption(err.Error()))
		return
	}

	//check code
	if _, err := h.readPasswordResetCode(ctx, user.ID, resetPasswordRequestForm.Code); err != nil {
		logger.Error(ctx.Request.Context(), err)
		http.Fail(ctx, http.MsgOption(err.Error()))
		return
	}

	//update password
	password, err := user.GenerateFromPassword(resetPasswordRequestForm.NewPassword)
	if err != nil {
		logger.Error(ctx.Request.Context(), err)
		http.Fail(ctx, http.MsgOption(err.Error()))
		return
	}
	user.Password = password
	if err := h.UpdateUserDB(ctx, session, &user); err != nil {
		logger.Error(ctx.Request.Context(), err)
		http.Fail(ctx, http.MsgOption(err.Error()))
		return
	}

	http.Success(ctx, http.MsgOption("ResetPassword succeeded"))
}

//sendVerificationEmail...
func (h *Handler) sendVerificationEmail(ctx context.Context, fromName, toAddress, toUsername, subject, textContent, token, redirctUrl, failureRedirectUrl string) error {
	uri := "https://baidu.com"
	query := "?token=" + token + "&redirectUrl=" + url.QueryEscape(redirctUrl) + "&failureRedirectUrl=" + url.QueryEscape(failureRedirectUrl)
	textContent = strings.Replace(textContent, "$verification_link", uri+query, -1)
	email := email.Email{
		Subject: subject,
		From:    fromName,
		To:      []string{toAddress},
		Text:    []byte(textContent),
	}
	err := h.sendEmail(ctx, &email)
	return err
}

//SendPasswordResetEmail...
func (h *Handler) sendPasswordResetEmail(ctx context.Context, userId uint, codeStr, fromName, toAddress, toUsername, subject, textContent string) error {
	textContent = strings.Replace(textContent, "$code", codeStr, -1)
	email := email.Email{
		Subject: subject,
		From:    fromName,
		To:      []string{toAddress},
		Text:    []byte(textContent),
	}
	err := h.sendEmail(ctx, &email)
	return err
}

func (h *Handler) generatePasswordResetCodeStoreKey(userId uint, code string) string {
	return fmt.Sprintf("user/password-reset-codes/%v-%v", userId, code)
}

//savePasswordResetCode...
func (h *Handler) savePasswordResetCode(ctx context.Context, userId uint, code string, expiry time.Duration) (*types.PasswordResetCode, error) {
	pwcode := types.PasswordResetCode{
		Expires: time.Now().Add(expiry),
		UserID:  userId,
		Code:    code,
	}

	if err := h.Cache.Set(h.generatePasswordResetCodeStoreKey(userId, code), pwcode, expiry); err != nil {
		return nil, err
	}
	return &pwcode, nil
}

// readPasswordResetCode returns the user reset code
func (h *Handler) readPasswordResetCode(ctx context.Context, userId uint, code string) (*types.PasswordResetCode, error) {
	pwcode := types.PasswordResetCode{}
	err := h.Cache.Get(h.generatePasswordResetCodeStoreKey(userId, code), &pwcode)
	if err != nil && err != cache.ErrCacheMiss {
		return nil, err
	}

	// check the expiry
	if pwcode.Expires.Before(time.Now()) {
		return nil, errors.New(errors.ENone, "password reset code expired")
	}
	return &pwcode, nil
}
