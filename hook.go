package emqx

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Hook struct {
	Name      string `json:"name"`
	Enable    bool   `json:"enable,omitempty"`
	Addr      string `json:"url"`
	Timeout   string `json:"request_timeout,omitempty"`
	Action    string `json:"failed_action,omitempty"`
	Reconnect string `json:"auto_reconnect,omitempty"`
	PoolSize  int    `json:"pool_size,omitempty"`
}

func (c *Client) HookUpdate(ctx context.Context, h *Hook) error {
	pay, err := json.Marshal(h)

	if err != nil {
		return fmt.Errorf("req: %w", err)
	}

	url := c.Base + "/exhooks/" + h.Name
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, bytes.NewReader(pay))

	if err != nil {
		return fmt.Errorf("req: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	res, err := c.Do(ctx, req)

	if err != nil {
		return fmt.Errorf("req: %w", err)
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		var buf bytes.Buffer

		if _, err := buf.ReadFrom(res.Body); err != nil {
			return fmt.Errorf("res: %w", err)
		}

		return fmt.Errorf("api: %v", buf.String())
	}

	return nil
}

func (c *Client) HookGet(ctx context.Context, name string) (*Hook, error) {
	url := c.Base + "/exhooks/" + name
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)

	if err != nil {
		return nil, fmt.Errorf("req: %w", err)
	}

	res, err := c.Do(ctx, req)

	if err != nil {
		return nil, fmt.Errorf("req: %w", err)
	}

	defer res.Body.Close()

	buf := bytes.Buffer{}
	pay := &Hook{}

	if _, err := buf.ReadFrom(res.Body); err != nil {
		return nil, fmt.Errorf("res: %w", err)
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("api: %v", buf.String())
	}

	if err := json.Unmarshal(buf.Bytes(), pay); err != nil {
		return nil, fmt.Errorf("res: %w", err)
	}

	return pay, nil
}

func (c *Client) HookCreate(ctx context.Context, h *Hook) error {
	pay, err := json.Marshal(h)

	if err != nil {
		return fmt.Errorf("req: %w", err)
	}

	url := c.Base + "/exhooks"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(pay))

	if err != nil {
		return fmt.Errorf("req: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	res, err := c.Do(ctx, req)

	if err != nil {
		return fmt.Errorf("req: %w", err)
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		var buf bytes.Buffer

		if _, err := buf.ReadFrom(res.Body); err != nil {
			return fmt.Errorf("res: %w", err)
		}

		return fmt.Errorf("api: %v", buf.String())
	}

	return nil

}
