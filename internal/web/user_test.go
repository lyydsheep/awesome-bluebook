package web

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/lyydsheep/awesome-bluebook/internal/service"
	svcmocks "github.com/lyydsheep/awesome-bluebook/internal/web/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserHandler_Signup(t *testing.T) {
	testCases := []struct {
		name           string
		mock           func(ctrl *gomock.Controller) service.BasicUserService
		input          string
		expectedCode   int
		expectedResult Result
	}{
		{
			name: "注册成功",
			mock: func(ctrl *gomock.Controller) service.BasicUserService {
				svc := svcmocks.NewMockBasicUserService(ctrl)
				svc.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return(nil)
				return svc
			},
			input:        `{"email":"123@qq.com","password":"123456","confirmPassword":"123456"}`,
			expectedCode: 200,
			expectedResult: Result{
				Msg: "注册成功",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			svc := tc.mock(ctrl)
			h := NewUserHandler(svc)
			gin.SetMode(gin.ReleaseMode)
			server := gin.Default()
			h.RegisterRouters(server)
			req, err := http.NewRequest(http.MethodPost, "/users/signup", bytes.NewBuffer([]byte(tc.input)))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()
			server.ServeHTTP(resp, req)
			assert.Equal(t, tc.expectedCode, resp.Code)
			decoder := json.NewDecoder(resp.Body)
			var res Result
			err = decoder.Decode(&res)
			require.NoError(t, err)
			assert.Equal(t, tc.expectedResult, res)
		})
	}
}
