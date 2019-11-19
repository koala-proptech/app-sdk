package apps_sdk

import (
	"context"
	"fmt"
	"net/http"
)

type RequestCreateNotificationMessage struct {
	NotificationMessage *NotificationMessage `json:"notification_message,omitempty"`
	UserId              string               `json:"user_id,omitempty"`
}

type NotificationMessage struct {
	Message    string `json:"message,omitempty"`
	Title      string `json:"title,omitempty"`
	Type       string `json:"type,omitempty"`
	InstanceId string `json:"instance_id,omitempty"`
}

func (c *Client) CreateNotification(ctx context.Context, token string, req RequestCreateNotificationMessage) (*Response, error) {
	url := fmt.Sprintf("%s/notifications/internal/messages", c.url)
	return c.walk(http.MethodPost, url, token, req)
}
