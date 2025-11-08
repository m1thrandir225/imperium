package programs

import (
	"context"
	"os"
	"testing"

	"github.com/m1thrandir225/imperium/apps/host/internal/util"
)

var testCtx context.Context

const (
	RawgApiKey = "rawg_api_key"
)

func TestMain(m *testing.M) {
	rawgApiKey := os.Getenv("RAWG_API_KEY")
	if util.IsEmptyString(rawgApiKey) {
		panic("RAWG_API_KEY is required for tests")
	}

	testCtx = context.WithValue(context.Background(), RawgApiKey, rawgApiKey)

	code := m.Run()

	os.Exit(code)
}

func GetRAWGApiKey() string {
	return testCtx.Value(RawgApiKey).(string)
}
