package internal

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"hash"
	"strings"
	"time"
)

// TOTP - Time-based one-time password generates a one-time password (OTP) that uses the current time as a source of uniqueness
func TOTP(secret interface{}, timeStep int64, numDigits int) string {
	counter := time.Now().Unix() / timeStep

	var secretBytes []byte
	switch secret.(type) {
	case string:
		secretBytes, _ = base32.StdEncoding.DecodeString(
			strings.ToUpper(secret.(string)),
		)
	case []byte:
		secretBytes = secret.([]byte)
	}

	counterBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(counterBytes, uint64(counter))

	h := hmac.New(func() hash.Hash { return sha1.New() }, secretBytes)
	h.Write(counterBytes)
	newHash := h.Sum(nil)

	offset := newHash[len(newHash)-1] & 0x0f
	binaryValue := (int32(newHash[offset]&0x7f) << 24) |
		(int32(newHash[offset+1]) << 16) |
		(int32(newHash[offset+2]) << 8) |
		int32(newHash[offset+3])

	code := int64(binaryValue) % int64(pow10(numDigits))

	return fmt.Sprintf("%0*d", numDigits, code)
}

func pow10(n int) int {
	if n == 0 {
		return 1
	}
	return 10 * pow10(n-1)
}
