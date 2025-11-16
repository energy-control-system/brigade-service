package user

import (
	"fmt"
	"net/http"

	"github.com/sunshineOfficial/golib/goctx"
	"github.com/sunshineOfficial/golib/gohttp"
)

type Client struct {
	client  gohttp.Client
	baseURL string
}

func NewClient(client gohttp.Client, baseURL string) *Client {
	return &Client{
		client:  client,
		baseURL: baseURL,
	}
}

func (c *Client) GetUsersByIDs(ctx goctx.Context, userIDs []int) ([]User, error) {
	url, err := gohttp.AddIntQuery(c.baseURL+"/users", "ids", userIDs...)
	if err != nil {
		return nil, fmt.Errorf("AddIntQuery: %w", err)
	}

	var response []User
	status, err := c.client.DoJson(ctx, http.MethodGet, url, nil, &response)
	if err != nil {
		return nil, fmt.Errorf("c.client.DoJson: %w", err)
	}

	if status != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", status)
	}

	return response, nil
}
