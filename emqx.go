package emqx

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"time"
)

const (
	base = "http://%s:%d/api/v5"

	DefaultHost    = "localhost"
	DefaultPort    = 18083
	DefaultTimeout = 5 * time.Second
)

type Option func(*Client) error

func WithHost(h string) Option {
	return func(c *Client) error {
		c.host = h
		return nil
	}
}

func WithPort(p int) Option {
	return func(c *Client) error {
		c.port = p
		return nil
	}
}

func WithUser(u string) Option {
	return func(c *Client) error {
		c.user = u
		return nil
	}
}

func WithPass(p string) Option {
	return func(c *Client) error {
		c.pass = p
		return nil
	}
}

func WithTimeout(t time.Duration) Option {
	return func(c *Client) error {
		c.tout = t
		return nil
	}
}

func WithClient(c *http.Client) Option {
	return func(cli *Client) error {
		if c == nil {
			return fmt.Errorf("client: nil")
		}

		cli.con = c
		return nil
	}
}

func WithLogger(l *slog.Logger) Option {
	return func(c *Client) error {
		if l == nil {
			return fmt.Errorf("logger: nil")
		}

		c.log = l
		return nil
	}
}

type Client struct {
	Base string

	log *slog.Logger
	con *http.Client
	mux sync.Mutex

	host string
	port int
	user string
	pass string
	tout time.Duration
}

func NewClient(base string, opts ...Option) (*Client, error) {
	cli := &Client{
		log: slog.Default(),
		con: &http.Client{},

		host: DefaultHost,
		port: DefaultPort,
		tout: DefaultTimeout,
	}

	for _, opt := range opts {
		if err := opt(cli); err != nil {
			return nil, fmt.Errorf("opt: %w", err)
		}
	}

	cli.Base = fmt.Sprintf(base, cli.host, cli.port)

	return cli, nil
}

func (c *Client) Host() string {
	return c.host
}

func (c *Client) Port() int {
	return c.port
}

func (cli *Client) Do(ctx context.Context, req *http.Request) (res *http.Response, err error) {
	cli.mux.Lock()
	defer cli.mux.Unlock()

	if cli.user != "" {
		req.SetBasicAuth(cli.user, cli.pass)
	}

	for {
		res, err = cli.con.Do(req)

		if err == nil {
			break
		}

		slog.Error("attempt", "req", req, "err", err)

		if ctx.Err() != nil {
			break
		}

		time.Sleep(cli.tout)
	}

	return
}
