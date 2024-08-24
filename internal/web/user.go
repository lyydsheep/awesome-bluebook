package web

import (
	"errors"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
	"github.com/lyydsheep/awesome-bluebook/internal/domain"
	"github.com/lyydsheep/awesome-bluebook/internal/repository"
	"github.com/lyydsheep/awesome-bluebook/internal/service"
	"net/http"
)

type UserHandler struct {
	UserSvc               service.BasicUserService
	passwordRegexpPattern *regexp.Regexp
	emailRegexpPattern    *regexp.Regexp
}

func NewUserHandler(UserSvc service.BasicUserService) *UserHandler {
	const passwordRegexpPattern = "^(?=.*\\d).{6,18}$"
	const emailRegexpPattern = "^[A-Za-z0-9]+([_\\.][A-Za-z0-9]+)*@([A-Za-z0-9\\-]+\\.)+[A-Za-z]{2,6}$"
	return &UserHandler{
		// 预编译正则表达式
		UserSvc:               UserSvc,
		passwordRegexpPattern: regexp.MustCompile(passwordRegexpPattern, regexp.None),
		emailRegexpPattern:    regexp.MustCompile(emailRegexpPattern, regexp.None),
	}
}

func (h *UserHandler) RegisterRouters(server *gin.Engine) {
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
	err := ctx.Bind(&req)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{Code: 5, Msg: "系统错误"})
		return
	}
	// 校验前端传来的数据
	ok, err := h.emailRegexpPattern.MatchString(req.Email)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{Code: 5, Msg: "系统错误"})
		return
	}
	if !ok {
		ctx.JSON(http.StatusOK, Result{Code: 4, Msg: "邮箱格式不对"})
		return
	}
	if req.Password != req.ConfirmPassword {
		ctx.JSON(http.StatusOK, Result{Code: 4, Msg: "两次密码不一致"})
		return
	}
	ok, err = h.passwordRegexpPattern.MatchString(req.Password)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{Code: 5, Msg: "系统错误"})
		return
	}
	if !ok {
		ctx.JSON(http.StatusOK, Result{Code: 4, Msg: "密码至少包含数字且6位"})
		return
	}
	err = h.UserSvc.SignUp(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	switch {
	case errors.Is(err, repository.ErrDuplicateUser):
		ctx.JSON(http.StatusOK, Result{Code: 4, Msg: "邮箱冲突了哦"})
		return
	case err == nil:
		ctx.JSON(http.StatusOK, Result{Msg: "注册成功"})
		return
	default:
		ctx.JSON(http.StatusOK, Result{Code: 5, Msg: "系统错误"})
		return
	}
}
