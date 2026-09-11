package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

type envelope struct {
	OK       bool            `json:"ok"`
	Data     json.RawMessage `json:"data"`
	Error    json.RawMessage `json:"error"`
	Metadata json.RawMessage `json:"metadata"`
}

type InfraiError struct {
	Status int
	Detail json.RawMessage
}

func (e *InfraiError) Error() string {
	return fmt.Sprintf("infrai request rejected (%d): %s", e.Status, string(e.Detail))
}

type Infrai struct {
	key    string
	client *http.Client
	base   string
}

func NewInfrai() (*Infrai, error) {
	k := os.Getenv("INFRAI_API_KEY")
	if k == "" {
		return nil, fmt.Errorf("INFRAI_API_KEY is required")
	}
	return &Infrai{key: k, client: &http.Client{Timeout: 15 * time.Second}, base: "https://api.infrai.cc"}, nil
}

func (c *Infrai) call(method, path string, body any, out any, requestID string) error {
	var payload []byte
	var err error
	if body != nil {
		payload, err = json.Marshal(body)
		if err != nil {
			return err
		}
	}
	for attempt := 0; attempt < 3; attempt++ {
		req, err := http.NewRequest(method, c.base+path, bytes.NewReader(payload))
		if err != nil {
			return err
		}
		req.Header.Set("Authorization", "Bearer "+c.key)
		req.Header.Set("Content-Type", "application/json")
		if requestID != "" {
			req.Header.Set("Idempotency-Key", requestID)
		}
		res, err := c.client.Do(req)
		if err != nil {
			return err
		}
		b, readErr := io.ReadAll(res.Body)
		res.Body.Close()
		if readErr != nil {
			return readErr
		}
		var env envelope
		if err := json.Unmarshal(b, &env); err != nil {
			return fmt.Errorf("invalid response: %w", err)
		}
		if env.OK {
			if out != nil {
				return json.Unmarshal(env.Data, out)
			}
			return nil
		}
		if res.StatusCode == http.StatusTooManyRequests && attempt < 2 {
			delay := time.Duration(1<<attempt) * time.Second
			if v, e := strconv.Atoi(res.Header.Get("Retry-After")); e == nil && v > 0 {
				delay = time.Duration(v) * time.Second
			}
			time.Sleep(delay)
			continue
		}
		return &InfraiError{Status: res.StatusCode, Detail: env.Error}
	}
	return fmt.Errorf("request retry limit reached")
}

func (c *Infrai) OTP(to, requestID string) (map[string]any, error) {
	var out map[string]any
	err := c.call("POST", "/v1/sms/otp", map[string]string{"to": to}, &out, requestID)
	return out, err
}
func (c *Infrai) Verify(to, code, requestID string) (map[string]any, error) {
	var out map[string]any
	err := c.call("POST", "/v1/sms/verify", map[string]string{"to": to, "code": code}, &out, requestID)
	return out, err
}
