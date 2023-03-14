package service

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	mocks "github.com/stretchr/testify/mock"

	"auth-service/api/pkg/model"
	"auth-service/api/pkg/model/mock"
)

func TestNewToken(t *testing.T) {
	var idExp int64 = 15 * 60
	var refreshExp int64 = 3 * 24 * 3600

	priv, _ := os.ReadFile("../../config/rsa_private_test.pem")
	privKey, _ := jwt.ParseRSAPrivateKeyFromPEM(priv)
	pub, _ := os.ReadFile("../../config/rsa_public_test.pem")
	pubKey, _ := jwt.ParseRSAPublicKeyFromPEM(pub)
	secret := "anotsorandomtestsecret"

	mockTokenRepository := new(mock.MockTokenRepository)

	tokenService := NewTokenService(&TSConfig{
		TokenRepository:       mockTokenRepository,
		PrivateKey:            privKey,
		PublicKey:             pubKey,
		RefreshSecret:         secret,
		RefreshExpirationSecs: refreshExp,
		IDExpirationSecs:      idExp,
	})

	uid, _ := uuid.NewRandom()
	u := &model.User{
		UID:      uid,
		Email:    "bob@bob.com",
		Password: "blarghedymcblarghface",
	}

	uidErrorCase, _ := uuid.NewRandom()
	uErrorCase := &model.User{
		UID:      uidErrorCase,
		Email:    "failure@failure.com",
		Password: "blarghedymcblarghface",
	}
	prevID := "a_previous_tokenID"

	setSuccessArguments := mocks.Arguments{
		mocks.AnythingOfType("*context.emptyCtx"),
		u.UID.String(),
		mocks.AnythingOfType("string"),
		mocks.AnythingOfType("time.Duration"),
	}

	setErrorArguments := mocks.Arguments{
		mocks.AnythingOfType("*context.emptyCtx"),
		uidErrorCase.String(),
		mocks.AnythingOfType("string"),
		mocks.AnythingOfType("time.Duration"),
	}

	deleteWithPrevIDArguments := mocks.Arguments{
		mocks.AnythingOfType("*context.emptyCtx"),
		u.UID.String(),
		prevID,
	}

	// mock call argument/responses
	mockTokenRepository.On("SetRefreshToken", setSuccessArguments...).Return(nil)
	mockTokenRepository.On("SetRefreshToken", setErrorArguments...).Return(fmt.Errorf("Error setting refresh token"))
	mockTokenRepository.On("DeleteRefreshToken", deleteWithPrevIDArguments...).Return(nil)

	t.Run("Returns a token pair with proper values", func(t *testing.T) {
		ctx := context.Background()
		tokenPair, err := tokenService.NewToken(ctx, u, prevID)
		assert.NoError(t, err)

		mockTokenRepository.AssertCalled(t, "SetRefreshToken", setSuccessArguments...)
		mockTokenRepository.AssertCalled(t, "DeleteRefreshToken", deleteWithPrevIDArguments...)

		var s string
		assert.IsType(t, s, tokenPair.IDToken.SS)

		idTokenClaims := &idTokenCustomClaims{}

		_, err = jwt.ParseWithClaims(tokenPair.IDToken.SS, idTokenClaims, func(token *jwt.Token) (interface{}, error) {
			return pubKey, nil
		})

		assert.NoError(t, err)

		// assert claims on idToken
		expectedClaims := []interface{}{
			u.UID,
			u.Email,
			u.Name,
			u.ImageURL,
		}
		actualIDClaims := []interface{}{
			idTokenClaims.User.UID,
			idTokenClaims.User.Email,
			idTokenClaims.User.Name,
			idTokenClaims.User.ImageURL,
		}

		assert.ElementsMatch(t, expectedClaims, actualIDClaims)
		assert.Empty(t, idTokenClaims.User.Password) // password should never be encoded to json

		expiresAt := time.Unix(idTokenClaims.RegisteredClaims.ExpiresAt.Unix(), 0)
		expectedExpiresAt := time.Now().Add(time.Duration(idExp) * time.Second)
		assert.WithinDuration(t, expectedExpiresAt, expiresAt, 5*time.Second)

		refreshTokenClaims := &refreshTokenCustomClaims{}
		_, err = jwt.ParseWithClaims(tokenPair.RefreshToken.SS, refreshTokenClaims, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})

		assert.IsType(t, s, tokenPair.RefreshToken.SS)

		// assert claims on refresh token
		assert.NoError(t, err)
		assert.Equal(t, u.UID, refreshTokenClaims.UID)

		expiresAt = time.Unix(refreshTokenClaims.RegisteredClaims.ExpiresAt.Unix(), 0)
		expectedExpiresAt = time.Now().Add(time.Duration(refreshExp) * time.Second)
		assert.WithinDuration(t, expectedExpiresAt, expiresAt, 5*time.Second)
	})

	t.Run("Error setting refresh token", func(t *testing.T) {
		ctx := context.Background()
		_, err := tokenService.NewToken(ctx, uErrorCase, "")
		assert.Error(t, err) // should return an error

		mockTokenRepository.AssertCalled(t, "SetRefreshToken", setErrorArguments...)
		mockTokenRepository.AssertNotCalled(t, "DeleteRefreshToken")
	})

	t.Run("Empty string provided for prevID", func(t *testing.T) {
		ctx := context.Background()
		_, err := tokenService.NewToken(ctx, u, "")
		assert.NoError(t, err)

		mockTokenRepository.AssertCalled(t, "SetRefreshToken", setSuccessArguments...)
		mockTokenRepository.AssertNotCalled(t, "DeleteRefreshToken")
	})
}

func TestValidateIDToken(t *testing.T) {
	var idExp int64 = 15 * 60

	priv, _ := os.ReadFile("../../config/rsa_private_test.pem")
	privKey, _ := jwt.ParseRSAPrivateKeyFromPEM(priv)
	pub, _ := os.ReadFile("../../config/rsa_public_test.pem")
	pubKey, _ := jwt.ParseRSAPublicKeyFromPEM(pub)

	// instantiate a common token service to be used by all tests
	tokenService := NewTokenService(&TSConfig{
		PrivateKey:       privKey,
		PublicKey:        pubKey,
		IDExpirationSecs: idExp,
	})

	// include password to make sure it is not serialized
	// since json tag is "-"
	uid, _ := uuid.NewRandom()
	u := &model.User{
		UID:      uid,
		Email:    "bob@bob.com",
		Password: "blarghedymcblarghface",
	}

	t.Run("Valid token", func(t *testing.T) {
		ss, _ := GenerateIDToken(u, privKey, idExp)

		uFromToken, err := tokenService.ValidateIDToken(ss)
		assert.NoError(t, err)

		assert.ElementsMatch(
			t,
			[]interface{}{u.Email, u.Name, u.UID, u.ImageURL},
			[]interface{}{uFromToken.Email, uFromToken.Name, uFromToken.UID, uFromToken.ImageURL},
		)
	})

	t.Run("Expired token", func(t *testing.T) {
		// maybe not the best approach to depend on utility method
		// token will be valid for 15 minutes
		ss, _ := GenerateIDToken(u, privKey, -1) // expires one second ago

		expectedErr := model.UnauthorizedError("Unable to verify user from idToken")

		_, err := tokenService.ValidateIDToken(ss)
		assert.EqualError(t, err, expectedErr.Message)
	})

	t.Run("Invalid signature", func(t *testing.T) {
		// maybe not the best approach to depend on utility method
		// token will be valid for 15 minutes
		ss, _ := GenerateIDToken(u, privKey, -1) // expires one second ago

		expectedErr := model.UnauthorizedError("Unable to verify user from idToken")

		_, err := tokenService.ValidateIDToken(ss)
		assert.EqualError(t, err, expectedErr.Message)
	})

	// TODO - Add other invalid token types
}
