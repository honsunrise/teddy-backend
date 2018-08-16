package controllers

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/zhsyourai/URCF-engine/api/http/controllers/shard"
	"github.com/zhsyourai/teddy-backend/common/gin-jwt"
	"github.com/zhsyourai/teddy-backend/common/types"
	"github.com/zhsyourai/teddy-backend/uaa/services"
	"net/http"
	"time"
)

func NewAccountController(middleware *gin_jwt.JwtMiddleware, generator *gin_jwt.JwtGenerator) *AccountController {
	return &AccountController{
		service:    services.NewAccountService(),
		middleware: middleware,
		generator:  generator,
	}
}

// AccountController is our /uaa controller.
type AccountController struct {
	service    services.AccountService
	middleware *gin_jwt.JwtMiddleware
	generator  *gin_jwt.JwtGenerator
}

func (c *AccountController) Handler(root *gin.RouterGroup) {
	root.GET("/register", c.RegisterHandler)
	root.POST("/login", c.LoginHandler)
	root.POST("/logout", c.middleware.Handler, c.LogoutHandler)
	root.POST("/change_password", c.middleware.Handler, c.ChangePassword)
}

func (c *AccountController) RegisterHandler(ctx *gin.Context) {
	registerRequest := &types.RegisterRequest{}
	ctx.Bind(registerRequest)

	user, err := c.service.Register(registerRequest.Username, registerRequest.Password, registerRequest.Roles)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	registerResponse := types.RegisterResponse{
		Username:   user.Username,
		Role:       user.Roles,
		CreateDate: user.CreateDate,
	}

	ctx.JSON(http.StatusOK, registerResponse)
}

func (c *AccountController) LoginHandler(ctx *gin.Context) {
	loginRequest := &types.LoginRequest{}
	ctx.Bind(loginRequest)

	user, err := c.service.Verify(loginRequest.Username, loginRequest.Password)
	if err != nil {
		ctx.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	token, err := c.generator.GenerateJwt(time.Hour, time.Hour*3, jwt.MapClaims{
		"uid":      user.ID,
		"username": user.Username,
		"roles":    user.Roles,
	})

	if err != nil {
		ctx.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"access_token": token,
		"type":         "bearer",
		"expire_in":    time.Hour.Seconds(),
	})
}

func (c *AccountController) ChangePassword(ctx *gin.Context) {
	changePasswordRequest := &shard.ChangePasswordRequest{}
	ctx.Bind(changePasswordRequest)

	token, err := c.middleware.ExtractToken(ctx)
	if err != nil {
		ctx.AbortWithError(http.StatusUnauthorized, err)
		return
	}
	claims := token.Claims.(jwt.MapClaims)

	if err := c.service.ChangePassword(claims["uid"].(string),
		changePasswordRequest.OldPassword, changePasswordRequest.NewPassword); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
	}
}

func (c *AccountController) LogoutHandler(ctx *gin.Context) {
	//token, ok := ctx.Get("jwt")
	//if !ok {
	//}
	//jwtToken := token.(*jwt.Token)
	//claims := jwtToken.Claims.(jwt.MapClaims)
}
