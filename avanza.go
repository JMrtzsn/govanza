package avanza

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/JMrtzsn/govanza/internal"
	"hash"
	"math"
	"net/http"
	"strings"
	"time"
)

// TODO https://pkg.go.dev/github.com/xlzd/gotp?utm_source=godoc

const (
	BaseURL            = "https://www.avanza.se"
	MinInactiveMinutes = 30
	MaxInactiveMinutes = 60 * 24
)

type Avanza struct {
	AuthenticationTimeout int
	Session               *http.Client

	Credentials           map[string]string
	AuthenticationSession string
	PushSubscriptionID    string
	CustomerID            string

	Socket *internal.AvanzaSocket
}

func NewAvanza(credentials map[string]string) (*Avanza, error) {

	session := &http.Client{}
	avanza := &Avanza{
		AuthenticationTimeout: MaxInactiveMinutes,
		Session:               session,
		Credentials:           credentials,
	}

	responseBody, err := avanza.authenticate()
	if err != nil {
		return nil, err
	}

	avanza.AuthenticationSession = responseBody["authenticationSession"].(string)
	avanza.PushSubscriptionID = responseBody["pushSubscriptionId"].(string)
	avanza.CustomerID = responseBody["customerId"].(string)

	socket, err := internal.NewAvanzaSocket(avanza.PushSubscriptionID, session)
	if err != nil {
		return nil, err
	}
	avanza.Socket = socket

	return avanza, nil
}

func (avanza *Avanza) authenticate() (map[string]interface{}, error) {
	data := map[string]interface{}{
		"maxInactiveMinutes": avanza.AuthenticationTimeout,
		"username":           avanza.Credentials["username"],
		"password":           avanza.Credentials["password"],
	}

	routeAuthenticationPath := fmt.Sprintf("%s%s", BaseURL, internal.AuthenticationPath.String())

	response, err := avanza.sendRequest(http.MethodPost, routeAuthenticationPath, data)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// TODO implement struct for response
	var responseBody map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&responseBody)
	if err != nil {
		return nil, err
	}

	// No second factor required, continue with normal login
	if _, ok := responseBody["twoFactorLogin"]; !ok {
		return responseBody, nil
	}

	tfaMethod := responseBody["twoFactorLogin"].(map[string]interface{})["method"].(string)

	if tfaMethod != "TOTP" {
		return nil, fmt.Errorf("unsupported two factor method %s", tfaMethod)
	}

	return avanza.validate2FA()
}

func (avanza *Avanza) validate2FA() (map[string]interface{}, error) {
	var totpCode string
	if totpSecret, ok := avanza.Credentials["totpSecret"]; ok {
		totpCode = GenerateTOTPCode(totpSecret)
	} else if totpCodeValue, ok := avanza.Credentials["totpCode"]; ok {
		totpCode = totpCodeValue
	}

	if totpCode == "" {
		return nil, fmt.Errorf("failed to get TOTP code")
	}

	data := map[string]interface{}{
		"method":   "TOTP",
		"totpCode": totpCode,
	}

	response, err := avanza.sendRequest(http.MethodPost, RouteTOTPPath, data)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var responseBody map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&responseBody)
	if err != nil {
		return nil, err
	}

	return responseBody, nil
}

func (avanza *Avanza) sendRequest(method string, path string, data map[string]interface{}) (*http.Response, error) {
	method = strings.ToUpper(method)
	url := fmt.Sprintf("%s%s", BaseURL, path)

	var body []byte
	if data != nil {
		body, _ = json.Marshal(data)
	}

	request, err := http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-AuthenticationSession", avanza.AuthenticationSession)

	if avanza.Socket != nil && avanza.Socket.Connected {
		request.Header.Add("X-SecurityToken", avanza.Socket.SecurityToken)
	}

	response, err := avanza.Session.Do(request)
	if err != nil {
		return nil, err
	}

	if response.StatusCode >= 400 {
		return nil, fmt.Errorf("request failed with status code %d", response.StatusCode)
	}

	return response, nil
}

func GenerateTOTPCode(secret string) string {
	totp := NewTOTP(secret, sha1.New, 6, 30)
	return totp.GenerateCode()
}

func NewTOTP(secret string, hashFunc func() hash.Hash, digits int, period int) *TOTP {
	return &TOTP{
		Secret:   secret,
		HashFunc: hashFunc,
		Digits:   digits,
		Period:   period,
	}
}

type TOTP struct {
	Secret   string
	HashFunc func() hash.Hash
	Digits   int
	Period   int
}

func (t *TOTP) GenerateCode() string {
	timeInterval := time.Now().Unix() / int64(t.Period)
	secret := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(t.Secret)

	hmac := hmac.New(t.HashFunc, secret)
	binary.Write(hmac, binary.BigEndian, timeInterval)
	hash := hmac.Sum(nil)

	offset := hash[len(hash)-1] & 0x0F
	code := int64(binary.BigEndian.Uint32(hash[offset:offset+4]) & 0x7FFFFFFF)

	code %= int64(math.Pow10(t.Digits))
	format := fmt.Sprintf("%%0%dd", t.Digits)

	return fmt.Sprintf(format, code)
}
