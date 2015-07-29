package pcp

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
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

	body, err := c.get(url)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, context); err != nil {
		return err
	}

	return nil
}

func (c *Client) Metrics(context *Context, prefix string) ([]Metric, error) {
	c.logger.Debugln("Fetching metrics for context...")
	result := make(map[string][]Metric)

	u := fmt.Sprintf(
		"%s%s/%d/%s",
		c.Endpoint,
		c.path,
		context.ContextID,
		"_metric",
	)
	if prefix != "" {
		v := url.Values{}
		v.Set("prefix", prefix)
		u += "?" + v.Encode()
	}
	c.logger.Debugf("Generated metrics url: %s", u)

	body, err := c.get(u)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return result["metrics"], nil
}

func (c *Client) get(url string) ([]byte, error) {
	resp, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
