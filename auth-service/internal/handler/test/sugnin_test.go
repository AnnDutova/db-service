package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	mocks "github.com/stretchr/testify/mock"

	"auth-service/api/internal/handler"
	"auth-service/api/pkg/model"
	"auth-service/api/pkg/model/mock"
)

func TestSignin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUserService := new(mock.MockUserService)
	mockTokenService := new(mock.MockTokenService)

	router := gin.Default()

	handler.NewHandler(&handler.Config{
		R:            router,
		UserService:  mockUserService,
		TokenService: mockTokenService,
	})

	t.Run("Bad request data", func(t *testing.T) {
		rr := httptest.NewRecorder()

		reqBody, err := json.Marshal(gin.H{
			"email":    "notanemail",
			"password": "short",
		})
		assert.NoError(t, err)

		request, err := http.NewRequest(http.MethodPost, "/signin", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(rr, request)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		mockUserService.AssertNotCalled(t, "Signin")
		mockTokenService.AssertNotCalled(t, "NewTokensFromUser")
	})

	t.Run("Error Returned from UserService.Signin", func(t *testing.T) {
		email := "bob@bob.com"
		password := "pwdoesnotmatch123"

		mockUSArgs := mocks.Arguments{
			mocks.AnythingOfType("*context.emptyCtx"),
			&model.User{Email: email, Password: password},
		}

		mockError := model.UnauthorizedError("invalid email/password combo")

		mockUserService.On("Signin", mockUSArgs...).Return(mockError)

		rr := httptest.NewRecorder()

		reqBody, err := json.Marshal(gin.H{
			"email":    email,
			"password": password,
		})
		assert.NoError(t, err)

		request, err := http.NewRequest(http.MethodPost, "/signin", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(rr, request)

		mockUserService.AssertCalled(t, "Signin", mockUSArgs...)
		mockTokenService.AssertNotCalled(t, "NewTokensFromUser")
		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Successful Token Creation", func(t *testing.T) {
		email := "bob@bob.com"
		password := "pwworksgreat123"

		mockUSArgs := mocks.Arguments{
			mocks.AnythingOfType("*context.emptyCtx"),
			&model.User{Email: email, Password: password},
		}

		mockUserService.On("Signin", mockUSArgs...).Return(nil)

		mockTSArgs := mocks.Arguments{
			mocks.AnythingOfType("*context.emptyCtx"),
			&model.User{Email: email, Password: password},
			"",
		}

		mockTokenPair := &model.Token{
			IDToken:      model.IDToken{SS: "idToken"},
			RefreshToken: model.RefreshToken{SS: "refreshToken"},
		}

		mockTokenService.On("NewToken", mockTSArgs...).Return(mockTokenPair, nil)

		rr := httptest.NewRecorder()

		reqBody, err := json.Marshal(gin.H{
			"email":    email,
			"password": password,
		})
		assert.NoError(t, err)

		request, err := http.NewRequest(http.MethodPost, "/signin", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(rr, request)

		respBody, err := json.Marshal(gin.H{
			"tokens": mockTokenPair,
		})
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockUserService.AssertCalled(t, "Signin", mockUSArgs...)
		mockTokenService.AssertCalled(t, "NewToken", mockTSArgs...)
	})

	t.Run("Failed Token Creation", func(t *testing.T) {
		email := "cannotproducetoken@bob.com"
		password := "cannotproducetoken"

		mockUSArgs := mocks.Arguments{
			mocks.AnythingOfType("*context.emptyCtx"),
			&model.User{Email: email, Password: password},
		}

		mockUserService.On("Signin", mockUSArgs...).Return(nil)

		mockTSArgs := mocks.Arguments{
			mocks.AnythingOfType("*context.emptyCtx"),
			&model.User{Email: email, Password: password},
			"",
		}

		mockError := model.InternalError()
		mockTokenService.On("NewToken", mockTSArgs...).Return(nil, mockError)
		rr := httptest.NewRecorder()

		reqBody, err := json.Marshal(gin.H{
			"email":    email,
			"password": password,
		})
		assert.NoError(t, err)

		request, err := http.NewRequest(http.MethodPost, "/signin", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(rr, request)

		respBody, err := json.Marshal(gin.H{
			"error": mockError,
		})
		assert.NoError(t, err)

		assert.Equal(t, mockError.Status, rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockUserService.AssertCalled(t, "Signin", mockUSArgs...)
		mockTokenService.AssertCalled(t, "NewToken", mockTSArgs...)
	})
}
