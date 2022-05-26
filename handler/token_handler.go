package handler

import (
	"fmt"
	"gin-boilerplate/comm/http"
	"gin-boilerplate/types"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/teris-io/shortid"
)

var (
	userIDKey = "userID"
)

//TokenValid ...
func (h *Handler) TokenValid(ctx *gin.Context) {
	tokenAuth, err := h.ExtractTokenMetadata(ctx.Request)
	//Token either expired or not valid
	if err != nil {
		http.Fail(ctx, http.MsgOption("Please login first"), http.StatusOption(http.StatusUnauthorized))
		ctx.Abort()
		return
	}

	userID, err := h.FetchAuth(tokenAuth)
	//Token does not exists in Redis (User logged out or expired)
	if err != nil {
		http.Fail(ctx, http.MsgOption("Please login first"), http.StatusOption(http.StatusUnauthorized))
		ctx.Abort()
		return
	}

	//To be called from GetUserID()
	ctx.Set(userIDKey, userID)

	//Next middle
	ctx.Next()
}

//Refresh ...
func (h *Handler) Refresh(ctx *gin.Context) {
	var tokenForm types.Token
	if ctx.ShouldBindJSON(&tokenForm) != nil {
		http.Fail(ctx, http.MsgOption("Invalid form"))
		return
	}
	//verify the token
	token, err := jwt.Parse(tokenForm.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("REFRESH_SECRET")), nil
	})

	//if there is an error, the token must have expired
	if err != nil {
		http.Fail(ctx, http.MsgOption("Invalid authorization, please login again"), http.StatusOption(http.StatusUnauthorized))
		return
	}
	//is token valid?
	if token != nil && !token.Valid {
		http.Fail(ctx, http.MsgOption("Invalid authorization, please login again"), http.StatusOption(http.StatusUnauthorized))
		return
	}

	//Since token is valid, get the uuid:
	claims, ok := token.Claims.(jwt.MapClaims) //the token claims should conform to MapClaims
	if !ok || !token.Valid {
		http.Fail(ctx, http.MsgOption("Invalid authorization, please login again"), http.StatusOption(http.StatusUnauthorized))
		return
	}
	refreshUUID, ok := claims["refresh_uuid"].(string) //convert the interface to string
	if !ok {
		http.Fail(ctx, http.MsgOption("Invalid authorization, please login again"), http.StatusOption(http.StatusUnauthorized))
		return
	}
	userID, err := strconv.ParseInt(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
	if err != nil {
		http.Fail(ctx, http.MsgOption("Invalid authorization, please login again"), http.StatusOption(http.StatusUnauthorized))
		return
	}

	//Delete the previous Refresh Token
	err = h.DeleteAuth(refreshUUID)
	if err != nil { //if any goes wrong
		http.Fail(ctx, http.MsgOption("Invalid authorization, please login again"), http.StatusOption(http.StatusUnauthorized))
		return
	}

	//Create new pairs of refresh and access tokens
	ts, err := h.CreateToken(uint(userID))
	if err != nil {
		http.Fail(ctx, http.MsgOption("Invalid authorization, please login again"), http.StatusOption(http.StatusUnauthorized))
		return
	}

	//save the tokens metadata to redis
	err = h.CreateAuth(uint(userID), ts)
	if err != nil {
		http.Fail(ctx, http.MsgOption("Invalid authorization, please login again"), http.StatusOption(http.StatusUnauthorized))
		return
	}
	tk := gin.H{
		"access_token":  ts.AccessToken,
		"refresh_token": ts.RefreshToken,
	}
	http.Success(ctx, http.FlatOption(tk))
}

//CreateToken ...
func (h *Handler) CreateToken(userID uint) (*types.TokenDetails, error) {
	td := &types.TokenDetails{}
	td.AtExpires = time.Now().Add(time.Minute * 15).Unix()
	td.AccessUUID = shortid.MustGenerate()

	td.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
	td.RefreshUUID = shortid.MustGenerate()

	var err error
	//Creating Access Token
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = td.AccessUUID
	atClaims["user_id"] = userID
	atClaims["exp"] = td.AtExpires

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return nil, err
	}
	//Creating Refresh Token
	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUUID
	rtClaims["user_id"] = userID
	rtClaims["exp"] = td.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(os.Getenv("REFRESH_SECRET")))
	if err != nil {
		return nil, err
	}
	return td, nil
}

//CreateAuth ...
func (h *Handler) CreateAuth(userid uint, td *types.TokenDetails) error {
	at := time.Unix(td.AtExpires, 0) //converting Unix to UTC(to Time object)
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()

	err := h.Cache.Set(td.AccessUUID, strconv.Itoa(int(userid)), at.Sub(now))
	if err != nil {
		return err
	}
	err = h.Cache.Set(td.RefreshUUID, strconv.Itoa(int(userid)), rt.Sub(now))
	if err != nil {
		return err
	}
	return nil
}

//ExtractTokenMetadata ...
func (h *Handler) ExtractTokenMetadata(r *http.Request) (*types.AccessDetails, error) {
	token, err := h.VerifyToken(r)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUUID, ok := claims["access_uuid"].(string)
		if !ok {
			return nil, err
		}
		userID, err := strconv.ParseInt(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
		if err != nil {
			return nil, err
		}
		return &types.AccessDetails{
			AccessUUID: accessUUID,
			UserID:     userID,
		}, nil
	}
	return nil, err
}

//VerifyToken ...
func (h *Handler) VerifyToken(r *http.Request) (*jwt.Token, error) {
	//Make sure that the token method conform to "SigningMethodHMAC"
	tokenString := h.ExtractToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

//ExtractToken ...
func (h *Handler) ExtractToken(r *http.Request) string {
	//normally Authorization the_token_xxx
	bearToken := r.Header.Get("Authorization")
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

//FetchAuth ...
func (h *Handler) FetchAuth(authD *types.AccessDetails) (int64, error) {
	var userid int64
	err := h.Cache.Get(authD.AccessUUID, &userid)
	if err != nil {
		return 0, err
	}
	return userid, nil
}

//DeleteAuth ...
func (h *Handler) DeleteAuth(givenUUID string) error {
	err := h.Cache.Delete(givenUUID)
	if err != nil {
		return err
	}
	return nil
}
