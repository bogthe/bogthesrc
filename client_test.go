package bogthesrc

import (
	"testing"
)

func TestClientBasics(t *testing.T) {
	t.Run("Can create a new request", func(t *testing.T) {
		c := NewClient(nil)
		req, err := c.NewRequest("GET", "posts")
		ua := req.Header.Get("User-Agent")

		if req.Method != "GET" {
			t.Errorf("Unexpected method, got: %s", req.Method)
		}

		if ua != userAgent {
			t.Errorf("Unexpected UA, got: %s", ua)
		}

		if err != nil {
			t.Errorf("Failed with: %s", err)
		}
	})
}
