package rest

import (
	"bytes"
	"context"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/magiconair/properties/assert"
	"github.com/wilfridterry/contact-list/internal/domain"
	mock_rest "github.com/wilfridterry/contact-list/internal/transport/rest/mocks"
	"go.uber.org/mock/gomock"
)

func TestHandler_signUp(t *testing.T) {
	type mockBehavior func(s *mock_rest.MockAuth, user *domain.SignUpInput)

	testTable := []struct {
		name                string
		inputBody           string
		inputUser           domain.SignUpInput
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:      "OK",
			inputBody: `{"name": "Test", "email": "test@test.com", "password": "qwerty"}`,
			inputUser: domain.SignUpInput{
				Name:     "Test",
				Email:    "test@test.com",
				Password: "qwerty",
			},
			mockBehavior: func(s *mock_rest.MockAuth, userInput *domain.SignUpInput) {
				s.EXPECT().SignUp(context.Background(), userInput).Return(domain.User{
					ID:    1,
					Name:  "Test",
					Email: "test@test.com",
				}, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"message": "Created", "user": {"id": 1}}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			// Init Deps
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_rest.NewMockAuth(c)
			testCase.mockBehavior(auth, &testCase.inputUser)

			handler := NewHandler(&mock_rest.MockContacts{}, auth)

			// Test Server

			r := gin.New()
			r.POST("/sign-up", handler.signUp)

			// Perform
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/sign-up", bytes.NewBufferString(testCase.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}
