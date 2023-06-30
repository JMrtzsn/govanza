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
	sync.Mutex
	conn          *websocket.Conn
	clientID      string
	messageCount  int
	connected     bool
	subscriptions map[string]struct {
		Callback func(string, map[string]interface{})
		ClientID string
	}
	pushSubscriptionID string // TODO: where does this come from?
}

// NewAvanzaSocket creates a new AvanzaSocket instance.
func NewAvanzaSocket(pushSubscriptionID, cookies string) (*AvanzaSocket, error) {
	headers := make(http.Header)
	headers.Add("Cookie", cookies)

	conn, _, err := websocket.DefaultDialer.Dial(webSocketURL, headers)
	if err != nil {
		return nil, err
	}

	s := &AvanzaSocket{
		conn: conn,
		subscriptions: make(map[string]struct {
			Callback func(string, map[string]interface{})
			ClientID string
		}),
		pushSubscriptionID: pushSubscriptionID,
	}

	err = s.sendHandshakeMessage(s.pushSubscriptionID)
	if err != nil {
		return nil, err
	}

	return s, nil
}

// Connect establishes the WebSocket connection and starts listening for messages.
func (s *AvanzaSocket) Connect() error {
	for {
		_, message, err := s.conn.ReadMessage()
		if err != nil {
			log.Println("Failed to read message from websocket:", err)
			return err
		}

		var messages []map[string]interface{}
		err = json.Unmarshal(message, &messages)
		if err != nil {
			log.Println("Failed to unmarshal message:", err)
			continue
		}

		for _, msg := range messages {
			channel, ok := msg["channel"].(string)
			if !ok {
				log.Println("Invalid message format: missing 'channel' field")
				continue
			}
			// TODO: log message
			// TODO: Convert message to struct
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
				s.handleCustomMessage(channel, msg)
			}

			if err != nil {
				log.Println("Failed to handle message:", err)
			}
		}
	}
}

// Close closes the WebSocket connection.
func (s *AvanzaSocket) Close() error {
	if s.conn != nil {
		err := s.conn.Close()
		s.conn = nil
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

	err := s.conn.WriteJSON([]interface{}{message})
	if err != nil {
		return err
	}

	s.messageCount++
	return nil
}

func (s *AvanzaSocket) sendConnectMessage() error {
	message := map[string]interface{}{
		"channel":        "/meta/connect",
		"clientId":       s.clientID,
		"connectionType": "websocket",
		"id":             s.messageCount,
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

	if _, ok := s.subscriptions[subscriptionString]; ok {
		return errors.New("subscription already exists")
	}

	s.subscriptions[subscriptionString] = struct {
		Callback func(string, map[string]interface{})
		ClientID string
	}{
		Callback: callback,
		ClientID: "",
	}

	message := map[string]interface{}{
		"channel":      "/meta/subscribe",
		"clientId":     s.clientID,
		"subscription": subscriptionString,
	}

	return s.send(message)
}

func (s *AvanzaSocket) handleDisconnectMessage() error {
	// TODO: log disconnect message
	return s.sendHandshakeMessage(s.pushSubscriptionID)
}

func (s *AvanzaSocket) handleHandshakeMessage(msg map[string]interface{}) error {
	successful, _ := msg["successful"].(bool)
	if successful {
		s.clientID, _ = msg["clientId"].(string)
		err := s.sendConnectMessage()
		if err != nil {
			return err
		}
	} else {
		advice, _ := msg["advice"].(map[string]interface{})
		reconnect, _ := advice["reconnect"].(string)
		if reconnect == "handshake" {
			err := s.sendHandshakeMessage(s.pushSubscriptionID)
			if err != nil {
				return err
			}
		}
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

		if !s.connected {
			s.connected = true
			err := s.resubscribeExistingSubscriptions()
			if err != nil {
				return err
			}
		}
	} else if s.clientID != "" {
		err := s.sendConnectMessage()
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *AvanzaSocket) resubscribeExistingSubscriptions() error {
	s.Lock()
	defer s.Unlock()

	for key, value := range s.subscriptions {
		if value.ClientID != s.clientID {
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

	if _, ok := s.subscriptions[subscription]; ok {
		s.subscriptions[subscription] = struct {
			Callback func(string, map[string]interface{})
			ClientID string
		}{
			Callback: s.subscriptions[subscription].Callback,
			ClientID: s.clientID,
		}
	}

	return nil
}

// handleCustomMessage handles custom messages for a specific channel.
func (s *AvanzaSocket) handleCustomMessage(channel string, msg map[string]interface{}) {
	// You can implement custom handling logic based on the 'channel' and 'msg' here.
	// For now, let's just print the received message.
	log.Printf("Received custom message on channel %s: %v\n", channel, msg)
}

func contains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}
