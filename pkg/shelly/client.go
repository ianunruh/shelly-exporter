package shelly

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func NewClient(target, username, password string, httpClient *http.Client) Client {
	return &client{
		target:   target,
		username: username,
		password: password,

		httpClient: httpClient,
	}
}

type Client interface {
	Status() (*Status, error)
}

type client struct {
	target   string
	username string
	password string

	httpClient *http.Client
}

func (c *client) Status() (*Status, error) {
	var result Status
	if err := c.get("/status", &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *client) get(path string, out interface{}) error {
	url := c.target + path

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	// optional basic auth
	if len(c.username) > 0 {
		req.SetBasicAuth(c.username, c.password)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	buffer := &bytes.Buffer{}
	if _, err := io.Copy(buffer, resp.Body); err != nil {
		return err
	}
	body := buffer.Bytes()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return json.Unmarshal(body, out)
}
