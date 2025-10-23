package user

import (
	"time"

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
	u := User{
		Role:        RoleInspector,
		Surname:     "Хрунин",
		Name:        "Дмитрий",
		Patronymic:  "Алексеевич",
		PhoneNumber: "+79371234567",
		Email:       "test@gmail.com",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	users := make([]User, 0, len(userIDs))
	for _, id := range userIDs {
		u.ID = id
		users = append(users, u)
	}

	return users, nil
}
