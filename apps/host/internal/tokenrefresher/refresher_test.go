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

func TestNew(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTokenGetter := mockhttpclient.NewMockTokenGetter(ctrl)
	mockTokenRefresher := mockhttpclient.NewMockTokenRefresher(ctrl)

	testCases := []struct {
		name        string
		errExpected bool
		build       func() (Refresher, error)
	}{
		{
			name:        "New() - valid configuration",
			errExpected: false,
			build: func() (Refresher, error) {
				return New(mockTokenGetter, mockTokenRefresher)
			},
		},
		{
			name:        "New() - invalid token getter",
			errExpected: true,
			build: func() (Refresher, error) {
				return New(nil, mockTokenRefresher)
			},
		},
		{
			name:        "New() - invalid token refresher",
			errExpected: true,
			build: func() (Refresher, error) {
				return New(mockTokenGetter, nil)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			refresher, err := tc.build()

			if tc.errExpected {
				require.Error(t, err)
				require.Empty(t, refresher)
				require.Nil(t, refresher)
			} else {
				require.NoError(t, err)
				require.NotNil(t, refresher)
				require.NotEmpty(t, refresher)
			}
		})
	}
}

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
		testDuration  time.Duration
		cancelContext bool
		callStop      bool
		description   string
	}{
		{
			name:          "context cancelled - should stop goroutine",
			testDuration:  100 * time.Millisecond,
			cancelContext: true,
			description:   "Should exit goroutine when context is cancelled",
		},
		{
			name:         "stop called - should stop goroutine",
			testDuration: 100 * time.Millisecond,
			callStop:     true,
			description:  "Should exit goroutine when Stop() is called",
		},
		{
			name:         "start without cancellation - goroutine should be running",
			testDuration: 100 * time.Millisecond,
			description:  "Goroutine should continue running when not cancelled or stopped",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			refresher, err := newAuthTokenRefresher(mockTokenGetter, mockTokenRefresher)
			require.NoError(t, err)

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			require.NotPanics(t, func() {
				refresher.Start(ctx)
			})

			time.Sleep(tc.testDuration)

			if tc.cancelContext {
				cancel()
			}

			if tc.callStop {
				refresher.Stop()
			}

			time.Sleep(50 * time.Millisecond)

			if tc.callStop {
				select {
				case <-refresher.stop:
				default:
					t.Error("stop channel should be closed after Stop() is called")
				}
			}
		})
	}
}
func TestAuthTokenRefresher_Start_TokenRefreshLogic(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTokenGetter := mockhttpclient.NewMockTokenGetter(ctrl)
	mockTokenRefresher := mockhttpclient.NewMockTokenRefresher(ctrl)

	refresher, err := newAuthTokenRefresher(mockTokenGetter, mockTokenRefresher)
	require.NoError(t, err)

	testCases := []struct {
		name          string
		expiresIn     time.Duration
		shouldRefresh bool
		refreshError  error
		description   string
	}{
		{
			name:          "token expires in 2 minutes - should refresh",
			expiresIn:     2 * time.Minute,
			shouldRefresh: true,
			description:   "Should refresh when token expires in less than 5 minutes",
		},
		{
			name:          "token expires in 10 minutes - should not refresh",
			expiresIn:     10 * time.Minute,
			shouldRefresh: false,
			description:   "Should not refresh when token expires in more than 5 minutes",
		},
		{
			name:          "token expires in 1 minute with refresh error",
			expiresIn:     1 * time.Minute,
			shouldRefresh: true,
			refreshError:  errors.New("refresh failed"),
			description:   "Should handle refresh errors gracefully",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			expiresAt := time.Now().Add(tc.expiresIn)

			if tc.shouldRefresh {
				mockTokenGetter.EXPECT().GetAccessTokenExpiresAt().Return(expiresAt)
				mockTokenRefresher.EXPECT().RefreshToken(gomock.Any()).Return(tc.refreshError)
			} else {
				mockTokenGetter.EXPECT().GetAccessTokenExpiresAt().Return(expiresAt)
			}

			expiresAtResult := refresher.getter.GetAccessTokenExpiresAt()
			if time.Until(expiresAtResult) < 5*time.Minute {
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				err := refresher.refresher.RefreshToken(ctx)
				cancel()

				if tc.refreshError != nil {
					require.Error(t, err)
					require.Equal(t, tc.refreshError, err)
				} else {
					require.NoError(t, err)
				}
			}
		})
	}
}

func TestAuthTokenRefresher_Stop(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTokenGetter := mockhttpclient.NewMockTokenGetter(ctrl)
	mockTokenRefresher := mockhttpclient.NewMockTokenRefresher(ctrl)

	testCases := []struct {
		name        string
		setupTest   func(*authTokenRefresher)
		expectPanic bool
	}{
		{
			name: "stop once- should close channel",
			setupTest: func(atr *authTokenRefresher) {

			},
			expectPanic: false,
		},
		{
			name: "stop twice - should not panic",
			setupTest: func(atr *authTokenRefresher) {
				atr.Stop()
			},
			expectPanic: false,
		},
		{
			name: "stop with running start - should stop gracefully",
			setupTest: func(atr *authTokenRefresher) {
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()

				expiresAt := time.Now().Add(10 * time.Minute)
				mockTokenGetter.EXPECT().GetAccessTokenExpiresAt().Return(expiresAt).AnyTimes()
				atr.Start(ctx)
				time.Sleep(50 * time.Millisecond)
			},
			expectPanic: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			refresher, err := newAuthTokenRefresher(mockTokenGetter, mockTokenRefresher)
			require.NoError(t, err)

			tc.setupTest(refresher)

			if tc.expectPanic {
				require.Panics(t, func() {
					refresher.Stop()
				})
			} else {
				require.NotPanics(t, func() {
					refresher.Stop()
				})

				select {
				case <-refresher.stop:
				default:
					t.Error("stop channel should be closed")
				}
			}

		})
	}
}

func TestAuthTokenRefresher_StartStop_Integration(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTokenGetter := mockhttpclient.NewMockTokenGetter(ctrl)
	mockTokenRefresher := mockhttpclient.NewMockTokenRefresher(ctrl)

	refresher, err := newAuthTokenRefresher(mockTokenGetter, mockTokenRefresher)
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	refresher.Start(ctx)

	time.Sleep(100 * time.Millisecond)

	refresher.Stop()

	time.Sleep(50 * time.Millisecond)

	select {
	case <-refresher.stop:
	default:
		t.Error("stop channel should be closed after Stop() is called")
	}
}
