package tokenrefresher

import (
	"context"
	"errors"
	"testing"
	"time"

	mockhttpclient "github.com/m1thrandir225/imperium/apps/host/internal/httpclient/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestNewAuthTokenRefresher(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTokenGetter := mockhttpclient.NewMockTokenGetter(ctrl)
	mockTokenRefresher := mockhttpclient.NewMockTokenRefresher(ctrl)

	testCases := []struct {
		name        string
		errExpected bool
		build       func() (*authTokenRefresher, error)
	}{
		{
			name:        "valid refresher",
			errExpected: false,
			build: func() (*authTokenRefresher, error) {
				return newAuthTokenRefresher(mockTokenGetter, mockTokenRefresher)
			},
		},
		{
			name:        "invalid-getter",
			errExpected: true,
			build: func() (*authTokenRefresher, error) {
				return newAuthTokenRefresher(nil, mockTokenRefresher)
			},
		},
		{
			name:        "invalid-refresher",
			errExpected: true,
			build: func() (*authTokenRefresher, error) {
				return newAuthTokenRefresher(mockTokenGetter, nil)
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tR, err := tc.build()
			if tc.errExpected {
				require.Error(t, err)
				require.Nil(t, tR)
				require.Empty(t, tR)
			} else {
				require.NoError(t, err)
				require.NotEmpty(t, tR)
				require.NotNil(t, tR)
			}
		})
	}
}

func TestAuthTokenRefresher_Start(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTokenGetter := mockhttpclient.NewMockTokenGetter(ctrl)
	mockTokenRefresher := mockhttpclient.NewMockTokenRefresher(ctrl)

	testCases := []struct {
		name          string
		setupMocks    func()
		testDuration  time.Duration
		expectedCalls int
		cancelContext bool
		callStop      bool
	}{
		{
			name: "token expires soon - should refresh",
			setupMocks: func() {
				expiresAt := time.Now().Add(2 * time.Minute)
				mockTokenGetter.EXPECT().GetAccessTokenExpiresAt().Return(expiresAt).AnyTimes()
				mockTokenRefresher.EXPECT().RefreshToken(gomock.Any()).Return(nil).Times(1)
			},
			testDuration:  1100 * time.Millisecond,
			expectedCalls: 1,
		},
		{
			name: "token expires later - should not refresh",
			setupMocks: func() {
				// Token expires in 10 minutes (more than 5 minute threshold)
				expiresAt := time.Now().Add(10 * time.Minute)
				mockTokenGetter.EXPECT().GetAccessTokenExpiresAt().Return(expiresAt).AnyTimes()
			},
			testDuration:  1100 * time.Millisecond,
			expectedCalls: 0,
		},
		{
			name: "context cancelled - should stop",
			setupMocks: func() {
				// Setup can be minimal since context will be cancelled
				expiresAt := time.Now().Add(10 * time.Minute)
				mockTokenGetter.EXPECT().GetAccessTokenExpiresAt().Return(expiresAt).AnyTimes()
			},
			testDuration:  500 * time.Millisecond,
			cancelContext: true,
		},
		{
			name: "stop called - should stop",
			setupMocks: func() {
				expiresAt := time.Now().Add(10 * time.Minute)
				mockTokenGetter.EXPECT().GetAccessTokenExpiresAt().Return(expiresAt).AnyTimes()
			},
			testDuration: 500 * time.Millisecond,
			callStop:     true,
		},
		{
			name: "refresh token error - should log but continue",
			setupMocks: func() {
				expiresAt := time.Now().Add(2 * time.Minute)
				mockTokenGetter.EXPECT().GetAccessTokenExpiresAt().Return(expiresAt).AnyTimes()
				mockTokenRefresher.EXPECT().RefreshToken(gomock.Any()).Return(errors.New("refresh failed")).Times(1)
			},
			testDuration: 1100 * time.Millisecond,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			refresher, err := newAuthTokenRefresher(mockTokenGetter, mockTokenRefresher)
			require.NoError(t, err)

			tc.setupMocks()

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			refresher.Start(ctx)

			time.Sleep(tc.testDuration)
			if tc.cancelContext {
				cancel()
			}

			if tc.callStop {
				refresher.Stop()
			}

			time.Sleep(100 * time.Millisecond)
		})
	}
}

func TestAuthTokenRefresher_Stop(t *testing.T) {
}
