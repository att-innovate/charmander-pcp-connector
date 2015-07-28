package pcp

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

type Client struct {
	Endpoint string
	path     string
	logger   *Logger
}

func NewClient(endpoint string) *Client {
	return &Client{
		Endpoint: endpoint,
		path:     "/pmapi",
		logger:   NewLogger(LOG_INFO),
	}
}

func (c *Client) SetLogLevel(level int) error {
	return c.logger.SetLogLevel(level)
}

func (c *Client) RefreshContext(context *Context) error {
	c.logger.Debugf("Refreshing context %v", *context)
	url := fmt.Sprintf(
		"%s%s/%s?%s",
		c.Endpoint,
		c.path,
		"context",
		context.params(),
	)
	c.logger.Debugf("Generated refresh url: %s", url)
	resp, err := http.Get(url)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	c.logger.Debugln(string(body))
	return nil
}
