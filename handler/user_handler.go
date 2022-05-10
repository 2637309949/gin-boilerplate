package handler

import (
	"gin-boilerplate/comm/db"
	"gin-boilerplate/comm/errors"
	"gin-boilerplate/comm/http"
	"gin-boilerplate/models"
	"gin-boilerplate/types"

	"github.com/gin-gonic/gin"
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
