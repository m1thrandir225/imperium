package httpclient

//const (
//	TokenGetterKey    = "token_getter"
//	TokenRefresherKey = "token_refresher"
//)
//
//var testCtx context.Context
//
//func TestMain(m *testing.M) {
//	token := "test_token"
//	expiresAt := time.Now().Add(time.Hour)
//
//	refreshToken := "refresh_token"
//
//	ctrl := gomock.NewController(m)
//
//	tokenGetter := mockhttpclient.NewMockTokenGetter(ctrl)
//
//	tokenRefresher := NewMockTokenRefresher()
//
//	ctx := context.Background()
//
//	ctx = context.WithValue(ctx, TokenGetterKey, tokenGetter)
//	ctx = context.WithValue(ctx, TokenRefresherKey, tokenRefresher)
//
//	testCtx = ctx
//
//	os.Exit(m.Run())
//}
//
//func GetTokenGetter() TokenGetter {
//	return testCtx.Value(TokenGetterKey).(TokenGetter)
//}
//
//func GetTokenRefresher() TokenRefresher {
//	return testCtx.Value(TokenRefresherKey).(TokenRefresher)
//}
