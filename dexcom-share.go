// Package dexcomshare provides an interface to the Dexcom Share API.
// The API is not public. You have been warned.
package dexcomshare

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

// Client is a client for the Dexcom Share API.
type Client struct {
	Username  string
	Password  string
	accountID string // accountID is the account ID returned by the authentication endpoint.
	sessionID string // sessionID is the session ID returned by the login endpoint.
	client    *http.Client
}

// Option is a function that configures a Client.
type Option func(*Client)

// WithClient configures a Client with a custom http.Client.
func WithClient(client *http.Client) Option {
	return func(c *Client) {
		c.client = client
	}
}

// NewClient creates a new Dexcom Share client.
func NewClient(username, password string, options ...Option) (*Client, error) {
	client := &Client{
		Username: username,
		Password: password,
		client:   &http.Client{},
	}

	for _, option := range options {
		option(client)
	}

	err := client.authenticate()
	if err != nil {
		return nil, err
	}

	err = client.login()
	if err != nil {
		return nil, err
	}

	return client, nil
}

type authenticateRequest struct {
	AccountName string `json:"accountName"`
	Password    string `json:"password"`
	Application string `json:"applicationId"`
}

// authenticate returns an account ID
func (c *Client) authenticate() error {
	b, err := json.Marshal(authenticateRequest{
		AccountName: c.Username,
		Password:    c.Password,
		Application: applicationID,
	})
	if err != nil {
		return err
	}

	request, err := http.NewRequest(http.MethodPost, baseURL+"/"+authenticateEndpoint, bytes.NewReader(b))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(request)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("authentication failed")
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &c.accountID)
	if err != nil {
		return err
	}

	return nil
}

type loginRequest struct {
	AccountID   string `json:"accountId"`
	Password    string `json:"password"`
	Application string `json:"applicationId"`
}

// login returns a session ID. It must be called after authenticate.
func (c *Client) login() error {
	b, err := json.Marshal(loginRequest{
		AccountID:   c.accountID,
		Password:    c.Password,
		Application: applicationID,
	})
	if err != nil {
		return err
	}

	request, err := http.NewRequest(http.MethodPost, baseURL+"/"+loginIDEndpoint, bytes.NewReader(b))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(request)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("login failed")
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &c.sessionID)
	if err != nil {
		return err
	}

	return nil
}

// GlucoseEntry represents a glucose entry.
type GlucoseEntry struct {
	Value int    `json:"Value"`
	Trend string `json:"Trend"`
	DT    string `json:"DT"`
}

type readGlucoseRequest struct {
	SessionID string `json:"sessionId"`
	Minutes   int    `json:"minutes"`
	MaxCount  int    `json:"maxCount"`
}

// ReadGlucose returns a list of glucose entries.
func (c *Client) ReadGlucose(minutes, maxCount int) ([]GlucoseEntry, error) {
	b, err := json.Marshal(readGlucoseRequest{
		SessionID: c.sessionID,
		Minutes:   minutes,
		MaxCount:  maxCount,
	})
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(http.MethodPost, baseURL+"/"+readGlucoseEndpoint, bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("read glucose failed")
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var entries []GlucoseEntry
	err = json.Unmarshal(data, &entries)
	if err != nil {
		return nil, err
	}

	return entries, nil
}
