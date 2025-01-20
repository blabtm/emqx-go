package emqx

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
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

	url := fmt.Sprintf("%s/gateways/%s", cli.Base, gtw.Type())
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, bytes.NewReader(pay))

	if err != nil {
		return err
	}

	res, err := cli.Do(ctx, req)

	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode == 204 {
		return nil
	}

	msg, err := io.ReadAll(res.Body)

	if err != nil {
		return fmt.Errorf("%w: [could not read response] %v", fmt.Errorf(res.Status), err)
	}

	return fmt.Errorf("%w: %v", fmt.Errorf(res.Status), msg)
}
