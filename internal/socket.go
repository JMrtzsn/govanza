package internal

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strings"
	"sync"
)

const (
	webSocketURL = "wss://www.avanza.se/_push/cometd"
)

// AvanzaSocket represents the Avanza WebSocket client.
type AvanzaSocket struct {
	sync.Mutex                             // Mutex for thread safety
	ClientID           string              // Initialized in handshake message
	Conn               *websocket.Conn     // The WebSocket connection
	Connected          bool                // True if the socket is Connected
	MessageCount       int                 // Keeps check of the number of messages sent
	Logger             *log.Logger         // Logger for logging
	PushSubscriptionID string              // TODO: where does this come from?
	Subscriptions      map[string]struct { // Map of subscriptions
		Callback func(string, map[string]interface{})
		ClientID string
	}
}

// NewAvanzaSocket creates a new AvanzaSocket instance with the given logger.
// If logger is nil, the default logger is used.
func NewAvanzaSocket(pushSubscriptionID, cookies string, reconnectLimit int, logger *log.Logger) (*AvanzaSocket, error) {
	headers := make(http.Header)
	headers.Add("Cookie", cookies)

	conn, _, err := websocket.DefaultDialer.Dial(webSocketURL, headers)
	if err != nil {
		return nil, err
	}

	if logger == nil {
		logger = log.Default()
	}

	s := &AvanzaSocket{
		Conn: conn,
		Subscriptions: make(map[string]struct {
			Callback func(string, map[string]interface{})
			ClientID string
		}),
		PushSubscriptionID: pushSubscriptionID,
		Logger:             logger,
	}

	err = s.sendHandshakeMessage(s.PushSubscriptionID)
	if err != nil {
		return nil, err
	}

	return s, nil
}

// Listen establishes the WebSocket connection and starts listening for messages.
func (s *AvanzaSocket) Listen() error {
	for {
		_, message, err := s.Conn.ReadMessage()
		if err != nil {
			s.Logger.Println("Failed to read message from websocket:", err)
			return err
		}

		var messages []map[string]interface{}
		err = json.Unmarshal(message, &messages)
		if err != nil {
			s.Logger.Println("Failed to unmarshal message:", err)
			continue
		}

		for _, msg := range messages {
			channel, ok := msg["channel"].(string)
			if !ok {
				s.Logger.Println("Invalid message format: missing 'channel' field")
				continue
			}

			s.Logger.Println("Received message on channel:", channel)
			s.Logger.Println("Message:", msg)
			switch channel {
			case "/meta/disconnect":
				err = s.handleDisconnectMessage()
			case "/meta/handshake":
				err = s.handleHandshakeMessage(msg)
			case "/meta/connect":
				err = s.handleConnectMessage(msg)
			case "/meta/subscribe":
				err = s.handleSubscribeMessage(msg)
			default:
				s.Logger.Println("Unknown channel:", channel)
			}

			if err != nil {
				s.Logger.Println("Failed to handle message:", err)
			}
		}
	}
}

// Close closes the WebSocket connection.
func (s *AvanzaSocket) Close() error {
	if s.Conn != nil {
		err := s.Conn.Close()
		s.Conn = nil
		return err
	}
	return nil
}

// SubscribeToID subscribes to a channel with a single ID.
func (s *AvanzaSocket) SubscribeToID(channel, id string, callback func(string, map[string]interface{})) error {
	return s.SubscribeToIDs(channel, []string{id}, callback)
}

// SubscribeToIDs subscribes to a channel with multiple IDs.
func (s *AvanzaSocket) SubscribeToIDs(channel string, ids []string, callback func(string, map[string]interface{})) error {
	if len(ids) == 0 {
		return errors.New("no IDs provided")
	}

	validChannelsForMultipleIDs := []string{
		"orders",
		"deals",
		"positions",
	}

	if len(ids) > 1 && !contains(validChannelsForMultipleIDs, channel) {
		return errors.New("multiple IDs are not supported for this channel")
	}

	subscriptionString := "/" + channel + "/" + strings.Join(ids, ",")
	return s.socketSubscribe(subscriptionString, callback)
}

func (s *AvanzaSocket) send(message interface{}) error {
	s.Lock()
	defer s.Unlock()

	err := s.Conn.WriteJSON([]interface{}{message})
	if err != nil {
		return err
	}

	s.MessageCount++
	return nil
}

func (s *AvanzaSocket) sendConnectMessage() error {
	message := map[string]interface{}{
		"channel":        "/meta/connect",
		"clientId":       s.ClientID,
		"connectionType": "websocket",
		"id":             s.MessageCount,
	}

	return s.send(message)
}

func (s *AvanzaSocket) sendHandshakeMessage(pushSubscriptionID string) error {
	message := map[string]interface{}{
		"advice": map[string]interface{}{
			"timeout":  60000,
			"interval": 0,
		},
		"channel":                  "/meta/handshake",
		"ext":                      map[string]interface{}{"subscriptionId": pushSubscriptionID},
		"minimumVersion":           "1.0",
		"supportedConnectionTypes": []string{"websocket", "long-polling", "callback-polling"},
		"version":                  "1.0",
	}

	return s.send(message)
}

func (s *AvanzaSocket) socketSubscribe(subscriptionString string, callback func(string, map[string]interface{})) error {
	s.Lock()
	defer s.Unlock()

	if _, ok := s.Subscriptions[subscriptionString]; ok {
		return errors.New("subscription already exists")
	}

	s.Subscriptions[subscriptionString] = struct {
		Callback func(string, map[string]interface{})
		ClientID string
	}{
		Callback: callback,
		ClientID: "",
	}

	message := map[string]interface{}{
		"channel":      "/meta/subscribe",
		"clientId":     s.ClientID,
		"subscription": subscriptionString,
	}

	return s.send(message)
}

func (s *AvanzaSocket) handleDisconnectMessage() error {
	// TODO: log disconnect message
	return s.sendHandshakeMessage(s.PushSubscriptionID)
}

func (s *AvanzaSocket) handleHandshakeMessage(msg map[string]interface{}) error {
	successful, _ := msg["successful"].(bool)
	if successful {
		s.ClientID, _ = msg["clientId"].(string)
		err := s.sendConnectMessage()
		if err != nil {
			return err
		}
	} else {
		advice, _ := msg["advice"].(map[string]interface{})
		reconnect, _ := advice["reconnect"].(string)
		if reconnect == "handshake" {
			err := s.sendHandshakeMessage(s.PushSubscriptionID)
			if err != nil {
				return err
			}
		}
		return errors.New("handshake failed")
	}
	return nil
}

func (s *AvanzaSocket) handleConnectMessage(msg map[string]interface{}) error {
	successful, _ := msg["successful"].(bool)
	advice, _ := msg["advice"].(map[string]interface{})
	reconnect := advice["reconnect"].(string) == "retry"
	interval, _ := advice["interval"].(float64)

	if successful && (advice == nil || reconnect && interval >= 0) {
		err := s.sendConnectMessage()
		if err != nil {
			return err
		}

		if !s.Connected {
			s.Connected = true
			err := s.resubscribeExistingSubscriptions()
			if err != nil {
				return err
			}
		}
	} else if s.ClientID != "" {
		err := s.sendConnectMessage()
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *AvanzaSocket) resubscribeExistingSubscriptions() error {
	for key, value := range s.Subscriptions {
		if value.ClientID != s.ClientID {
			err := s.socketSubscribe(key, value.Callback)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *AvanzaSocket) handleSubscribeMessage(msg map[string]interface{}) error {
	subscription, ok := msg["subscription"].(string)
	if !ok || subscription == "" {
		return errors.New("no subscription channel found on subscription message")
	}

	s.Lock()
	defer s.Unlock()

	if _, ok := s.Subscriptions[subscription]; ok {
		s.Subscriptions[subscription] = struct {
			Callback func(string, map[string]interface{})
			ClientID string
		}{
			Callback: s.Subscriptions[subscription].Callback,
			ClientID: s.ClientID,
		}
	}

	return nil
}

func contains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}
