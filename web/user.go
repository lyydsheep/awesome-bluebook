package web

import (
	"awesome-bluebook/domain"
	"awesome-bluebook/service"
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type UserHandler struct {
	svc service.UserService
	l   *zap.Logger
}

func NewUserHandler(svc service.UserService, l *zap.Logger) *UserHandler {
	return &UserHandler{
		svc: svc,
		l:   l,
	}
}

func (h *UserHandler) RegisterRoutes(server *gin.Engine) {
	sg := server.Group("/users")
	sg.POST("/signup", h.Signup)
}

func (h *UserHandler) Signup(ctx *gin.Context) {
	type Req struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		h.l.Info("Bind失败", zap.Error(err))
		ctx.JSON(http.StatusOK, Result[int]{
			Msg:  "系统错误",
			Code: 5,
		})
		return
	}
	if req.Password != req.ConfirmPassword {
		h.l.Info("两次密码不一致")
		ctx.JSON(http.StatusOK, Result[int]{
			Code: 4,
			Msg:  "密码不一致",
		})
		return
	}
	err := h.svc.Signup(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	switch {
	case err == nil:
		ctx.JSON(http.StatusOK, Result[int]{
			Msg: "注册成功",
		})
	case errors.Is(err, service.UserDuplicateErr):
		ctx.JSON(http.StatusOK, Result[int]{
			Msg:  "邮箱已存在",
			Code: 4,
		})
	}
}
