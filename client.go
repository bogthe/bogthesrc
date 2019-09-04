package bogthesrc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strings"

	"github.com/bogthe/bogthesrc/router"
	"github.com/google/go-querystring/query"
)

// define Client struct
type Client struct {
	BaseUrl   *url.URL
	UserAgent string

	httpClient *http.Client
}

type ApiClient struct {
	Client *Client
	Posts  PostService
}

type ListOptions struct {
	PerPage int `url:",omitempty" json:",omitempty"`
	Page    int `url:",omitempty" json:",omitempty"`
}

func (lo ListOptions) PageOrDefault() int {
	if lo.Page <= 0 {
		return 1
	}

	return lo.Page
}

func (lo ListOptions) PerPageOrDefault() int {
	if lo.PerPage <= 0 {
		return DefaultPerPage
	}

	return lo.PerPage
}

func (lo ListOptions) Offset() int {
	return (lo.PageOrDefault() - 1) * lo.PerPageOrDefault()
}

const (
	version        = "0.0.1"
	userAgent      = "bogthesrc-client" + version
	DefaultPerPage = 60
)

// NewClient
func NewClient(client *http.Client) *Client {
	if client == nil {
		client = http.DefaultClient
	}

	envPort := os.Getenv("PORT")
	if envPort != "" {
		envPort = fmt.Sprintf("localhost:%s", envPort)
	} else {
		envPort = "localhost:5000"
	}

	return &Client{
		BaseUrl:    &url.URL{Scheme: "http", Host: envPort, Path: "/api/"},
		UserAgent:  userAgent,
		httpClient: client,
	}
}

func NewApiClient(client *Client) *ApiClient {
	if client == nil {
		client = NewClient(nil)
	}

	posts := &postService{client}
	return &ApiClient{
		Client: client,
		Posts:  posts,
	}
}

var apiRouter = router.API()

func (c *Client) url(apiRouteName string, routeVars map[string]string, opt interface{}) (*url.URL, error) {
	route := apiRouter.Get(apiRouteName)
	if route == nil {
		return nil, fmt.Errorf("Route not found %s", apiRouteName)
	}

	i := 0
	routeList := make([]string, len(routeVars)*2)
	for key, value := range routeVars {
		routeList[i*2] = key
		routeList[i*2+1] = value
		i++
	}

	newUrl, err := route.URL(routeList...)
	if err != nil {
		return nil, err
	}

	if opt != nil {
		err = addOptions(newUrl, opt)
		if err != nil {
			return nil, err
		}
	}

	newUrl.Path = strings.TrimPrefix(newUrl.Path, "/")
	return newUrl, nil
}

func (c *Client) NewRequest(method, relativeUrl string, body interface{}) (*http.Request, error) {
	rel, err := url.Parse(relativeUrl)
	if err != nil {
		return nil, err
	}

	u := c.BaseUrl.ResolveReference(rel)

	buf := new(bytes.Buffer)
	if body != nil {
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	req.Header.Add("User-Agent", c.UserAgent)
	return req, nil
}

// client.Do - sends an API request and receives an API response, returns the response or error
func (c *Client) Do(request *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if v != nil {
		if bp, ok := v.(*[]byte); ok {
			*bp, err = ioutil.ReadAll(resp.Body)
		} else {
			err = json.NewDecoder(resp.Body).Decode(v)
		}
	}

	if err != nil {
		return nil, fmt.Errorf("Error reading response: %s %s %s", request.Method, request.URL.RequestURI(), err)
	}

	return resp, nil
}

func addOptions(u *url.URL, opt interface{}) error {
	v := reflect.ValueOf(opt)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return nil
	}

	queryString, err := query.Values(opt)
	if err != nil {
		return err
	}

	u.RawQuery = queryString.Encode()
	return nil
}
