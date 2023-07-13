package internal_test

import (
	"bytes"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/JMrtzsn/govanza/internal"
)

// TODO add tests for different messages (should connect to live server)
func TestNewAvanzaSocket(t *testing.T) {
	pushSubscriptionID := "12345"
	cookies := "test cookies"
	reconnectLimit := 5

	// Create a buffer to capture log output
	var logBuffer bytes.Buffer

	// Create a logger that writes to the buffer
	logger := log.New(&logBuffer, "", log.LstdFlags)

	socket, err := internal.NewAvanzaSocket(pushSubscriptionID, cookies, reconnectLimit, logger)
	assert.NotNil(t, socket, "Expected non-nil socket")
	assert.NoError(t, err, "Unexpected error")

	go func() {
		err := socket.Listen()
		assert.NoError(t, err, "Unexpected error")
	}()

	time.Sleep(2 * time.Second)

	t.Run("Test Connect", func(t *testing.T) {
		socket.Logger.Println("Test Connect")
		// Get the log messages from the buffer
		logMessages := logBuffer.String()

		// Assert on the log messages
		expectedLogMessage := "handshake failed"
		assert.Contains(t, logMessages, expectedLogMessage, "Log message not found")
	})

	socket.Close()
}

// TODO add tests for different messages (should connect to live server)
