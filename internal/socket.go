package internal

import (
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	WebsocketUrl string = "wss://www.avanza.se/_push/cometd"
)

type AvanzaSocket struct {
	socket             *websocket.Conn
	clientID           string
	messageCount       int
	pushSubscriptionID string
	connected          bool
	subscriptions      map[string]Subscription
	cookies            string
	logger             *log.Logger
}

type Subscription struct {
	Callback func(string, map[string]interface{})
	ClientID string
}

func NewAvanzaSocket(pushSubscriptionID, cookies string, logger *log.Logger) (*AvanzaSocket, error) {
	socket := &AvanzaSocket{
		socket:             nil,
		clientID:           "",
		messageCount:       1,
		pushSubscriptionID: pushSubscriptionID,
		connected:          false,
		subscriptions:      make(map[string]Subscription),
		cookies:            cookies,
		logger:             logger,
	}

	err := socket.init()
	if err != nil {
		return nil, err
	}

	return socket, nil
}

func (a *AvanzaSocket) init() error {
	var err error

	dialer := websocket.Dialer{
		Proxy:            http.ProxyFromEnvironment,
		HandshakeTimeout: 10 * time.Second,
	}

	a.socket, _, err = dialer.Dial(WebsocketUrl, http.Header{"Cookie": []string{a.cookies}})
	if err != nil {
		return fmt.Errorf("failed to connect to Avanza WebSocket: %v", err)
	}

	err = a.sendHandshakeMessage()
	if err != nil {
		return err
	}

	go func() {
		err := a.socketMessageHandler()
		if err != nil {
			a.logger.Printf("failed to handle websocket message: %v", err)
		}
	}()

	err = a.waitForWebsocketToBeConnected()
	if err != nil {
		return err
	}

	return nil
}

func (a *AvanzaSocket) waitForWebsocketToBeConnected() error {
	timeoutCount := 40
	timeoutValue := 250 * time.Millisecond

	for i := 0; i < timeoutCount; i++ {
		if a.connected {
			return nil
		}
		time.Sleep(timeoutValue)
	}

	return errors.New("failed to connect to the websocket within the expected timeframe")
}

func (a *AvanzaSocket) sendHandshakeMessage() error {
	return a.send(map[string]interface{}{
		"advice": map[string]interface{}{
			"timeout":  60000,
			"interval": 0,
		},
		"channel":                  "/meta/handshake",
		"ext":                      map[string]string{"subscriptionId": a.pushSubscriptionID},
		"minimumVersion":           "1.0",
		"supportedConnectionTypes": []string{"websocket", "long-polling", "callback-polling"},
		"version":                  "1.0",
	})
}

func (a *AvanzaSocket) handshake(message map[string]interface{}) error {
	successful, ok := message["successful"].(bool)
	if ok && successful {
		a.clientID, _ = message["clientId"].(string)
		err := a.send(map[string]interface{}{
			"advice":         map[string]interface{}{"timeout": 0},
			"channel":        "/meta/connect",
			"clientId":       a.clientID,
			"connectionType": "websocket",
		})
		if err != nil {
			return err
		}
		return nil
	}

	advice, ok := message["advice"].(map[string]interface{})
	if ok && advice["reconnect"] == "handshake" {
		return a.sendHandshakeMessage()
	}

	return nil
}

func (a *AvanzaSocket) sendConnectMessage() error {
	return a.send(map[string]interface{}{
		"channel":        "/meta/connect",
		"clientId":       a.clientID,
		"connectionType": "websocket",
		"id":             a.messageCount,
	})
}

func (a *AvanzaSocket) socketSubscribe(subscriptionString string, callback func(string, map[string]interface{})) error {
	a.subscriptions[subscriptionString] = Subscription{Callback: callback}

	return a.send(map[string]interface{}{
		"channel":      "/meta/subscribe",
		"clientId":     a.clientID,
		"subscription": subscriptionString,
	})
}

func (a *AvanzaSocket) send(message map[string]interface{}) error {
	wrappedMessage := []map[string]interface{}{{
		"id": a.messageCount,
	}}

	for k, v := range message {
		wrappedMessage[0][k] = v
	}

	err := a.socket.WriteJSON(wrappedMessage)
	if err != nil {
		return fmt.Errorf("failed to send message to Avanza WebSocket: %v", err)
	}

	a.messageCount++

	return nil
}

func (a *AvanzaSocket) connect(message map[string]interface{}) error {
	successful, ok := message["successful"].(bool)
	if !ok {
		successful = false
	}
	advice, ok := message["advice"].(map[string]interface{})
	if !ok {
		advice = map[string]interface{}{}
	}
	reconnect := advice["reconnect"] == "retry"
	interval, ok := advice["interval"].(float64)

	connectSuccessful := successful && (!reconnect || interval >= 0)

	if connectSuccessful {
		err := a.send(map[string]interface{}{
			"channel":        "/meta/connect",
			"clientId":       a.clientID,
			"connectionType": "websocket",
		})
		if err != nil {
			return fmt.Errorf("failed to send connect message: %v", err)
		}

		if !a.connected {
			a.connected = true
			err := a.resubscribeExistingSubscriptions()
			if err != nil {
				return fmt.Errorf("failed to resubscribe existing subscriptions: %v", err)
			}
		}
	} else if a.clientID != "" {
		err := a.sendConnectMessage()
		if err != nil {
			return fmt.Errorf("failed to send connect message: %v", err)
		}
	}

	return nil
}

func (a *AvanzaSocket) resubscribeExistingSubscriptions() error {
	for key, value := range a.subscriptions {
		if value["client_id"].(string) != a.clientID {
			callback := value["callback"].(func(string))
			err := a.socketSubscribe(key, callback)
			if err != nil {
				return fmt.Errorf("failed to resubscribe to channel %s: %v", key, err)
			}
		}
	}
	return nil
}

func (a *AvanzaSocket) disconnect(message map[string]interface{}) error {
	err := a.sendHandshakeMessage()
	if err != nil {
		return fmt.Errorf("failed to send handshake message: %v", err)
	}
	return nil
}

func (a *AvanzaSocket) registerSubscription(message map[string]interface{}) error {
	subscription, ok := message["subscription"].(string)
	if !ok {
		return fmt.Errorf("no subscription channel found on subscription message")
	}
	a.subscriptions[subscription]["client_id"] = a.clientID
	return nil
}

func (a *AvanzaSocket) socketMessageHandler() error {
	messageAction := map[string]func(message map[string]interface{}) error{
		"/meta/disconnect": a.disconnect,
		"/meta/handshake":  a.handshake,
		"/meta/connect":    a.connect,
		"/meta/subscribe":  a.registerSubscription,
	}

	for {
		var message []map[string]interface{}
		err := a.socket.ReadJSON(&message)
		if err != nil {
			return fmt.Errorf("failed to read message from Avanza WebSocket: %v", err)
		}

		messageChannel := message[0]["channel"].(string)
		errMsg := message[0]["error"]
		if errMsg != nil {
			log.Printf("Incoming error message on channel %v: %v", messageChannel, errMsg)
			continue
		}

		log.Printf("Incoming message on channel %v: %v", messageChannel, message)

		if action, ok := messageAction[messageChannel]; ok {
			err = action(message[0])
			if err != nil {
				log.Printf("Failed to process incoming message on channel %v: %v", messageChannel, err)
			}
		} else if subscription, ok := a.subscriptions[messageChannel]; ok {
			callback := subscription["callback"].(func(string, map[string]interface{}))
			callback(subscription["id"].(string), message[0])
		} else {
			logger.Warnf("No action or subscription found for incoming message on channel %v", messageChannel)
		}
	}
}

func (a *AvanzaSocket) subscribeToID(channel ChannelType, id string, callback func(string, map[string]interface{})) error {
	return a.subscribeToIDs(channel, []string{id}, callback)
}

func (a *AvanzaSocket) subscribeToIDs(channel ChannelType, ids []string, callback func(string, map[string]interface{})) error {
	validChannelsForMultipleIDs := []ChannelType{
		Orders,
		Deals,
		Positions,
	}

	if len(ids) > 1 && !contains(validChannelsForMultipleIDs, channel) {
		return fmt.Errorf("multiple ids are not supported for channels other than %v", validChannelsForMultipleIDs)
	}

	subscriptionString := fmt.Sprintf("/%v/%v", channel, strings.Join(ids, ","))
	err := a.socketSubscribe(subscriptionString, callback)
	if err != nil {
		return fmt.Errorf("failed to subscribe to ids: %v", err)
	}
	return nil

}

func contains(channels []ChannelType, channel ChannelType) bool {
	for _, c := range channels {
		if c == channel {
			return true
		}
	}
	return false
}
