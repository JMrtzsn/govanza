package internal

import (
	"crypto/rand"
	"encoding/base32"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TODO: use public library for TOTP?
func TestTOTP(t *testing.T) {
	t.Run("Assert that encoded keys generated on different times have different TOTP codes", func(t *testing.T) {
		secret := make([]byte, 10)
		_, err := rand.Read(secret)
		if err != nil {
			t.Fatalf("Error generating secret key: %v", err)
		}

		encoded := base32.StdEncoding.EncodeToString(secret)

		numDigitsList := []int{6, 14}
		for _, numDigits := range numDigitsList {
			code := TOTP(encoded, 2, numDigits)
			time.Sleep(2 * time.Second)
			code2 := TOTP(encoded, 2, numDigits)

			assert.NotEqual(t, code, code2, "Expected different TOTP codes")
		}
	})
	t.Run("Assert that string keys generated on different times have different TOTP codes", func(t *testing.T) {
		secret := make([]byte, 10)
		_, err := rand.Read(secret)
		if err != nil {
			t.Fatalf("Error generating secret key: %v", err)
		}

		numDigitsList := []int{6, 14}
		for _, numDigits := range numDigitsList {
			code := TOTP(secret, 2, numDigits)
			time.Sleep(2 * time.Second)
			code2 := TOTP(secret, 2, numDigits)

			assert.NotEqual(t, code, code2, "Expected different TOTP codes")
		}
	})
}
