//go:build devtoken

package token

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/spf13/viper"
)

// TestGenerateLocalDevToken generates a short-lived token for local integration testing.
// It is excluded from normal builds and tests; run with -tags devtoken explicitly.
func TestGenerateLocalDevToken(t *testing.T) {
	configFile := os.Getenv("FORUM_GATEWAY_CONFIG")
	if configFile == "" {
		t.Fatal("FORUM_GATEWAY_CONFIG is required")
	}

	viper.Reset()
	viper.SetConfigFile(configFile)
	if err := viper.ReadInConfig(); err != nil {
		t.Fatalf("read gateway config: %v", err)
	}

	// Make sure the test does not accidentally reuse a key from another test.
	jwtKey = ""
	value, err := GenerateToken(&TokenPayload{
		Id:      145,
		Role:    1, // constvar.Normal
		Expired: time.Hour,
	})
	if err != nil {
		t.Fatalf("generate dev token: %v", err)
	}

	fmt.Println(value)
}
