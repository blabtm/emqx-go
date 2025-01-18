package emqx

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Gateway interface {
	Type() string
}

type ExProtoServer struct {
	Bind string `json:"bind"`
}

type ExProtoHandler struct {
	Addr string `json:"address"`
}

type GatewayListener struct {
	ID           string   `json:"id"`
	Type         string   `json:"type"`
	Name         string   `json:"name"`
	Running      bool     `json:"running"`
	Acceptors    int      `json:"acceptors"`
	Proxy        bool     `json:"proxy_protocol"`
	ProxyTimeout string   `json:"proxy_protocol_timeout"`
	Enable       bool     `json:"enable"`
	Bind         string   `json:"bind"`
	MaxConns     int      `json:"max_connections"`
	MaxConnRate  int      `json:"max_conn_rate"`
	EnableAuthN  bool     `json:"enable_authn"`
	Mountpoint   string   `json:"mountpoint"`
	AccessRules  []string `json:"access_rules"`
}

type ExProtoGateway struct {
	Name       string            `json:"name"`
	Timeout    string            `json:"idle_timeout"`
	Mountpoint string            `json:"mountpoint"`
	Enable     bool              `json:"enable"`
	Statistics bool              `json:"enable_stats"`
	Server     ExProtoServer     `json:"server"`
	Handler    ExProtoHandler    `json:"handler"`
	Listeners  []GatewayListener `json:"listeners"`
}

func (*ExProtoGateway) Type() string {
	return "exproto"
}

func (cli *Client) GatewayUpdate(ctx context.Context, gw Gateway) error {
	pay, err := json.Marshal(gw)

	if err != nil {
		return fmt.Errorf("req: %w", err)
	}

	url := cli.Base + "/gateways/" + gw.Type()
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, bytes.NewReader(pay))

	if err != nil {
		return fmt.Errorf("req: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	res, err := cli.Do(ctx, req)

	if err != nil {
		return fmt.Errorf("req: %w", err)
	}

	defer res.Body.Close()

	if res.StatusCode != 204 {
		var buf bytes.Buffer

		if _, err := buf.ReadFrom(res.Body); err != nil {
			return fmt.Errorf("res: %w", err)
		}

		return fmt.Errorf("api: %v", buf.String())
	}

	return nil
}
