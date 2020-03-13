package apps_sdk

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/pkg/errors"
)

const (
	mimeApplicationJson = "application/json"
)

type (
	httpClient interface {
		Do(r *http.Request) (*http.Response, error)
	}
	Option func(*Client)
	Client struct {
		uid, token, clientId, url string
		httpClient                httpClient
		uriService                *Service
	}
	Service struct {
		Booking string `json:"booking"`
	}
	ErrorResponse struct {
		Code    string            `json:"code"`
		Message string            `json:"message"`
		Reasons map[string]string `json:"reasons"`
		Details []interface{}     `json:"details, omitempty"`
	}
	Response struct {
		RequestID string                 `json:"request_id"`
		Status    int                    `json:"status"`
		Content   map[string]interface{} `json:"content, omitempty"`
		Error     *ErrorResponse         `json:"error,omitempty"`
	}
)

func OptionHttpClient(client httpClient) func(*Client) {
	return func(c *Client) {
		c.httpClient = client
	}
}

func New(uid, url, clientId string, options ...Option) (*Client, error) {
	s := &Client{
		uid:        uid,
		url:        url,
		clientId:   clientId,
		httpClient: &http.Client{},
	}
	for _, opt := range options {
		opt(s)
	}
	return s, nil
}

func (c *Client) build(method, url, token string, payload io.Reader) (*http.Request, error) {
	r, err := http.NewRequest(method, url, payload)
	if err != nil {
		return nil, errors.WithMessage(err, "failed creating request")
	}
	r.Header.Set("X-Auth-Token", token)
	r.Header.Set("Content-Type", mimeApplicationJson)
	r.Header.Set("Accept", mimeApplicationJson)
	r.Header.Set("Cache-Control", "no-cache")
	return r, nil
}

func (c *Client) request(r *http.Request) (*Response, error) {
	resp, err := c.httpClient.Do(r)
	if err != nil {
		return nil, errors.WithMessage(err, "failed communication with upstream")
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusForbidden || resp.StatusCode == http.StatusUnauthorized {
		return &Response{Status: resp.StatusCode}, errors.New(http.StatusText(resp.StatusCode))
	}
	var rs Response
	if err := json.NewDecoder(resp.Body).Decode(&rs); err != nil {
		return nil, errors.WithMessage(err, "failed decoding response")
	}
	return &rs, nil
}

func (c *Client) talk(method, url, token string, payload interface{}) (*Response, error) {
	var ir io.Reader
	if nil != payload {
		b, err := json.Marshal(payload)
		if err != nil {
			return nil, errors.WithMessage(err, "failed encoding request payload")
		}
		ir = bytes.NewReader(b)
	}
	r, err := c.build(method, url, token, ir)
	if err != nil {
		return nil, err
	}
	return c.request(r)
}
