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
				assert.True(t, user.Ctime > 0)
				user.Ctime = 0
				assert.True(t, user.Utime > 0)
				user.Utime = 0
				assert.Equal(t, dao.User{
					Id:       1,
					Email:    "123@qq.com",
					Password: "1234567",
				}, user)
			},
			expectCode: http.StatusOK,
			expectRes: web.Result[int]{
				Msg: "登录成功",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			hdl := web.NewUserHandler()
			hdl.RegisterRoutes(u.server)
			req, err := http.NewRequest(http.MethodPost, "/users/signup", bytes.NewBuffer([]byte(tc.input)))
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
	u.db = startup.InitDB()
	dao.InitTable(u.db)
}

func (u *UserHandlerTestSuite) TearDownTest() {
	u.db.Exec("truncate table webook")
}

func TestUserHandler(t *testing.T) {
	suite.Run(t, new(UserHandlerTestSuite))
}
