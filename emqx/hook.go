package emqx

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Hook struct {
	Name      string `json:"name"`
	Enable    bool   `json:"enable"`
	Addr      string `json:"url"`
	Timeout   string `json:"request_timeout,omitempty"`
	Action    string `json:"failed_action,omitempty"`
	Reconnect string `json:"auto_reconnect,omitempty"`
	PoolSize  int    `json:"pool_size"`
}

func (c *Client) HookUpdate(ctx context.Context, hook *Hook) error {
	pay, err := json.Marshal(hook)

	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/exhooks/%s", c.Base, hook.Name)
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, bytes.NewReader(pay))

	if err != nil {
		return err
	}

	res, err := c.Do(ctx, req)

	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode == 200 {
		return nil
	}

	msg, err := io.ReadAll(res.Body)

	if err != nil {
		return fmt.Errorf("%w: [could not read response] %v", fmt.Errorf(res.Status), err)
	}

	return fmt.Errorf("%w: %v", fmt.Errorf(res.Status), string(msg))
}

func (c *Client) HookGet(ctx context.Context, name string) (*Hook, error) {
	url := c.Base + "/exhooks/" + name
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)

	if err != nil {
		return nil, err
	}

	res, err := c.Do(ctx, req)

	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	pay, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, fmt.Errorf("%w: [could not read response] %v", fmt.Errorf(res.Status), err)
	}

	if res.StatusCode == 200 {
		hook := &Hook{}

		if err := json.Unmarshal(pay, hook); err != nil {
			return nil, fmt.Errorf("%w: [could not parse response] %v", fmt.Errorf(res.Status), err)
		}

		return hook, nil
	}

	return nil, fmt.Errorf("%w: %v", fmt.Errorf(res.Status), string(pay))
}

func (c *Client) HookCreate(ctx context.Context, hook *Hook) error {
	pay, err := json.Marshal(hook)

	if err != nil {
		return err
	}

	url := c.Base + "/exhooks"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(pay))

	if err != nil {
		return err
	}

	res, err := c.Do(ctx, req)

	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode == 200 {
		return nil
	}

	msg, err := io.ReadAll(res.Body)

	if err != nil {
		return fmt.Errorf("%w: [could not read response] %v", fmt.Errorf(res.Status), err)
	}

	return fmt.Errorf("%w: %v", fmt.Errorf(res.Status), string(msg))

}
