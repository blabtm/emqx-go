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

type ExProtoGateway struct {
	Name       string         `json:"name,omitempty"`
	Server     ExProtoServer  `json:"server"`
	Handler    ExProtoHandler `json:"handler"`
	Mountpoint string         `json:"mountpoint,omitempty"`
	Enable     bool           `json:"enable"`
	Timeout    string         `json:"idle_timeout,omitempty"`
	Statistics bool           `json:"enable_stats"`
}

func (*ExProtoGateway) Type() string {
	return "exproto"
}

func (cli *Client) GatewayUpdate(ctx context.Context, gtw Gateway) error {
	pay, err := json.Marshal(gtw)

	if err != nil {
		return err
	}

	url := cli.Base + "/gateways/" + gtw.Type()
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, bytes.NewReader(pay))

	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	res, err := cli.Do(ctx, req)

	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode != 204 {
		var buf bytes.Buffer

		if _, err := buf.ReadFrom(res.Body); err != nil {
			return err
		}

		return fmt.Errorf("api: %v", buf.String())
	}

	return nil
}
