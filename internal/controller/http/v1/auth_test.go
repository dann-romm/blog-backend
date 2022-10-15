package v1

import (
	"blog-backend/internal/service"
	"blog-backend/internal/service/mocks"
	"blog-backend/pkg/validator"
	"bytes"
	"context"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthRoutes_SignUp(t *testing.T) {
	type args struct {
		ctx   context.Context
		input service.AuthCreateUserInput
	}

	type MockBehaviour func(s *mocks.MockAuth, args args)

	testCases := []struct {
		name            string
		args            args
		inputBody       string
		mockBehaviour   MockBehaviour
		wantStatusCode  int
		wantRequestBody string
	}{
		{
			name: "OK",
			args: args{
				ctx: context.Background(),
				input: service.AuthCreateUserInput{
					Name:     "Test Name",
					Username: "test",
					Password: "Qwerty!1",
					Email:    "test@example.com",
				},
			},
			inputBody: `{"name":"Test Name","username":"test","password":"Qwerty!1","email":"test@example.com"}`,
			mockBehaviour: func(s *mocks.MockAuth, args args) {
				s.EXPECT().CreateUser(args.ctx, args.input).Return(1, nil)
			},
			wantStatusCode:  200,
			wantRequestBody: `{"id":1}` + "\n",
		},
		{
			name:            "Invalid password: not provided",
			args:            args{},
			inputBody:       `{"name":"Test Name","username":"test","email":"test@example.com"}`,
			mockBehaviour:   func(s *mocks.MockAuth, args args) {},
			wantStatusCode:  400,
			wantRequestBody: `{"message":"field password is required"}` + "\n",
		},
		{
			name:            "Invalid password: too short",
			args:            args{},
			inputBody:       `{"name":"Test Name","username":"test","password":"Qw!1","email":"test@example.com"}`,
			mockBehaviour:   func(s *mocks.MockAuth, args args) {},
			wantStatusCode:  400,
			wantRequestBody: `{"message":"field password must be between 8 and 32 characters"}` + "\n",
		},
		{
			name:            "Invalid password: too long",
			args:            args{},
			inputBody:       `{"name":"Test Name","username":"test","password":"Qwerty!123456789012345678901234567890","email":"test@example.com"}`,
			mockBehaviour:   func(s *mocks.MockAuth, args args) {},
			wantStatusCode:  400,
			wantRequestBody: `{"message":"field password must be between 8 and 32 characters"}` + "\n",
		},
		{
			name:            "Invalid password: no uppercase",
			args:            args{},
			inputBody:       `{"name":"Test Name","username":"test","password":"qwerty!1","email":"test@example.com"}`,
			mockBehaviour:   func(s *mocks.MockAuth, args args) {},
			wantStatusCode:  400,
			wantRequestBody: `{"message":"field password must contain at least 1 uppercase letter(s)"}` + "\n",
		},
		{
			name:            "Invalid password: no lowercase",
			args:            args{},
			inputBody:       `{"name":"Test Name","username":"test","password":"QWERTY!1","email":"test@example.com"}`,
			mockBehaviour:   func(s *mocks.MockAuth, args args) {},
			wantStatusCode:  400,
			wantRequestBody: `{"message":"field password must contain at least 1 lowercase letter(s)"}` + "\n",
		},
		{
			name:            "Invalid password: no digit",
			args:            args{},
			inputBody:       `{"name":"Test Name","username":"test","password":"Qwerty!!","email":"test@example.com"}`,
			mockBehaviour:   func(s *mocks.MockAuth, args args) {},
			wantStatusCode:  400,
			wantRequestBody: `{"message":"field password must contain at least 1 digit(s)"}` + "\n",
		},
		{
			name:            "Invalid password: no special character",
			args:            args{},
			inputBody:       `{"name":"Test Name","username":"test","password":"Qwerty11","email":"test@example.com"}`,
			mockBehaviour:   func(s *mocks.MockAuth, args args) {},
			wantStatusCode:  400,
			wantRequestBody: `{"message":"field password must contain at least 1 special character(s)"}` + "\n",
		},
		{
			name:            "Invalid email: not provided",
			args:            args{},
			inputBody:       `{"name":"Test Name","username":"test","password":"Qwerty!1"}`,
			mockBehaviour:   func(s *mocks.MockAuth, args args) {},
			wantStatusCode:  400,
			wantRequestBody: `{"message":"field email is required"}` + "\n",
		},
		{
			name:            "Invalid email: no @",
			args:            args{},
			inputBody:       `{"name":"Test Name","username":"test","password":"Qwerty!1","email":"testexample.com"}`,
			mockBehaviour:   func(s *mocks.MockAuth, args args) {},
			wantStatusCode:  400,
			wantRequestBody: `{"message":"field email must be a valid email address"}` + "\n",
		},
		{
			name:            "Invalid name: not provided",
			args:            args{},
			inputBody:       `{"username":"test","password":"Qwerty!1","email":"test@example.com"}`,
			mockBehaviour:   func(s *mocks.MockAuth, args args) {},
			wantStatusCode:  400,
			wantRequestBody: `{"message":"field name is required"}` + "\n",
		},
		{
			name:            "Invalid name: too short",
			args:            args{},
			inputBody:       `{"name":"T","username":"test","password":"Qwerty!1","email":"test@example.com"}`,
			mockBehaviour:   func(s *mocks.MockAuth, args args) {},
			wantStatusCode:  400,
			wantRequestBody: `{"message":"field name must be at least 4 characters"}` + "\n",
		},
		{
			name:            "Invalid name: too long",
			args:            args{},
			inputBody:       `{"name": "Test Name Test Name Test Name Test Name", "username":"test","password":"Qwerty!1","email":"test@example.com"}`,
			mockBehaviour:   func(s *mocks.MockAuth, args args) {},
			wantStatusCode:  400,
			wantRequestBody: `{"message":"field name must be at most 32 characters"}` + "\n",
		},
		{
			name:            "Invalid name: invalid characters",
			args:            args{},
			inputBody:       `{"name": "Test Name!", "username":"test","password":"Qwerty!1","email":"test@example.com"}`,
			mockBehaviour:   func(s *mocks.MockAuth, args args) {},
			wantStatusCode:  400,
			wantRequestBody: `{"message":"field name must contain only letters"}` + "\n",
		},
		{
			name:            "Invalid username: not provided",
			args:            args{},
			inputBody:       `{"name":"Test Name","password":"Qwerty!1","email":"test@example.com"}`,
			mockBehaviour:   func(s *mocks.MockAuth, args args) {},
			wantStatusCode:  400,
			wantRequestBody: `{"message":"field username is required"}` + "\n",
		},
		{
			name:            "Invalid username: too short",
			args:            args{},
			inputBody:       `{"name":"Test Name","username":"t","password":"Qwerty!1","email":"test@example.com"}`,
			mockBehaviour:   func(s *mocks.MockAuth, args args) {},
			wantStatusCode:  400,
			wantRequestBody: `{"message":"field username must be at least 4 characters"}` + "\n",
		},
		{
			name:            "Invalid username: too long",
			args:            args{},
			inputBody:       `{"name":"Test Name","username":"testtesttesttesttesttesttesttesttest","password":"Qwerty!1","email":"test@example.com"}`,
			mockBehaviour:   func(s *mocks.MockAuth, args args) {},
			wantStatusCode:  400,
			wantRequestBody: `{"message":"field username must be at most 32 characters"}` + "\n",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// init deps
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// create service mock
			auth := mocks.NewMockAuth(ctrl)
			tc.mockBehaviour(auth, tc.args)
			services := &service.Services{Auth: auth}

			// create test server
			e := echo.New()
			e.Validator = validator.NewCustomValidator()
			g := e.Group("/auth")
			newAuthRoutes(g, services.Auth)

			// create request
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/auth/sign-up", bytes.NewBufferString(tc.inputBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			// execute request
			e.ServeHTTP(w, req)

			// check response
			assert.Equal(t, tc.wantStatusCode, w.Code)
			assert.Equal(t, tc.wantRequestBody, w.Body.String())
		})
	}
}
