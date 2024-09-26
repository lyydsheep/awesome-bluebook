package integration

import (
	"awesome-bluebook/integration/startup"
	"awesome-bluebook/repository/dao"
	"awesome-bluebook/web"
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"testing"
)

type UserHandlerTestSuite struct {
	suite.Suite
	server *gin.Engine
	db     *gorm.DB
}

func (u *UserHandlerTestSuite) TestUserHandler_Signup() {
	t := u.T()
	testCases := []struct {
		name       string
		before     func(t *testing.T)
		input      string
		after      func(t *testing.T)
		expectCode int
		expectRes  web.Result[int]
	}{
		{
			name: "注册成功",
			before: func(t *testing.T) {

			},
			input: `{
"email":"123@qq.com",
"password":"1234567",
"confirmPassword":"1234567"}`,
			after: func(t *testing.T) {
				var user dao.User
				err := u.db.Where("id=?", 1).First(&user).Error
				if err != nil {
					require.NoError(t, err)
				}
				err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte("1234567"))
				require.NoError(t, err)
				user.Password = ""
				assert.True(t, user.Ctime > 0)
				user.Ctime = 0
				assert.True(t, user.Utime > 0)
				user.Utime = 0
				assert.Equal(t, dao.User{
					Id:    1,
					Email: "123@qq.com",
				}, user)
			},
			expectCode: http.StatusOK,
			expectRes: web.Result[int]{
				Msg: "注册成功",
			},
		},
		{
			name: "Bind失败",
			before: func(t *testing.T) {

			},
			input: `
"email":"123@qq.com",
"password":"1234567",
"confirmPassword":"1234567"}`,
			after: func(t *testing.T) {
			},
			expectCode: http.StatusBadRequest,
			expectRes: web.Result[int]{
				Msg:  "系统错误",
				Code: 5,
			},
		},
		{
			name: "两次密码不一致",
			before: func(t *testing.T) {

			},
			input: `{
"email":"123@qq.com",
"password":"1234567",
"confirmPassword":"12345671"}`,
			after: func(t *testing.T) {
			},
			expectCode: http.StatusOK,
			expectRes: web.Result[int]{
				Msg:  "密码不一致",
				Code: 4,
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			req, err := http.NewRequest(http.MethodPost, "/users/signup", bytes.NewBuffer([]byte(tc.input)))
			req.Header.Set("Content-Type", "application/json")
			require.NoError(t, err)
			resp := httptest.NewRecorder()
			u.server.ServeHTTP(resp, req)
			var res web.Result[int]
			code := resp.Code
			err = json.NewDecoder(resp.Body).Decode(&res)
			require.NoError(t, err)
			assert.Equal(t, tc.expectCode, code)
			assert.Equal(t, tc.expectRes, res)
			tc.after(t)
		})
	}
}

func (u *UserHandlerTestSuite) SetupSuite() {
	startup.InitViper()
	u.server = startup.InitGin()
	hdl := startup.NewUserHandler()
	hdl.RegisterRoutes(u.server)
	u.db = startup.InitDB()
	dao.InitTable(u.db)
}

func (u *UserHandlerTestSuite) TearDownTest() {
	u.db.Exec("truncate table users")
}

func TestUserHandler(t *testing.T) {
	suite.Run(t, new(UserHandlerTestSuite))
}
