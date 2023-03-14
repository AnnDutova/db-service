package mock

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"auth-service/api/pkg/model"
)

type MockTokenService struct {
	mock.Mock
}

func (m *MockTokenService) NewToken(ctx context.Context, u *model.User, prevTokenID string) (*model.Token, error) {
	ret := m.Called(ctx, u, prevTokenID)

	var r0 *model.Token
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*model.Token)
	}

	var r1 error

	if ret.Get(1) != nil {
		r1 = ret.Get(1).(error)
	}

	return r0, r1
}

func (m *MockTokenService) Signout(ctx context.Context, uid uuid.UUID) error {
	ret := m.Called(ctx, uid)

	var r0 error
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(error)
	}

	return r0
}

func (m *MockTokenService) ValidateIDToken(tokenString string) (*model.User, error) {
	ret := m.Called(tokenString)

	var r0 *model.User
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*model.User)
	}

	var r1 error

	if ret.Get(1) != nil {
		r1 = ret.Get(1).(error)
	}

	return r0, r1
}

func (m *MockTokenService) ValidateRefreshToken(refreshTokenString string) (*model.RefreshToken, error) {
	ret := m.Called(refreshTokenString)

	var r0 *model.RefreshToken
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*model.RefreshToken)
	}

	var r1 error

	if ret.Get(1) != nil {
		r1 = ret.Get(1).(error)
	}

	return r0, r1
}
