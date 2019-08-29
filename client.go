package bogthesrc

import (
	"net/http"
	"net/url"
)

// define Client struct
type Client struct {
	BaseUrl   *url.URL
	UserAgent string

	httpClient *http.Client
}

const (
	version   = "0.0.1"
	userAgent = "bogthesrc-client" + version
)

// NewClient
func NewClient(client *http.Client) *Client {
	if client == nil {
		client = http.DefaultClient
	}

	return &Client{
		BaseUrl:    &url.URL{Scheme: "http", Host: "bogthesrc.co.uk", Path: "/api"},
		UserAgent:  userAgent,
		httpClient: client,
	}
}

// client.NewRequest - creates a new request http.NewRequest for method, relativeUrl, jsonBody?
func (c *Client) NewRequest(method, relativeUrl string) (*http.Request, error) {
	rel, err := url.Parse(relativeUrl)
	if err != nil {
		return nil, err
	}

	u := c.BaseUrl.ResolveReference(rel)

	req, err := http.NewRequest(method, u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("User-Agent", c.UserAgent)
	return req, nil
}

// client.Do - sends an API request and receives an API response, returns the response or error
func (c *Client) Do(request *http.Request) (*http.Response, error) {
	resp, err := c.httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return resp, nil
}
