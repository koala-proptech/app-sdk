package apps_sdk

import (
	"context"
	"fmt"
	"net/http"
)

func (c *Client) TemplateById(ctx context.Context, token string, templateId string) (*Response, error) {
	url := fmt.Sprintf("%s/notifications/template/%s", c.url, templateId)
	return c.talk(http.MethodGet, url, token, nil)
}
