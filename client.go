package pcp

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Client struct {
	Context  *Context
	Endpoint string
	logger   *Logger
	path     string
}

func NewClient(endpoint string, context *Context) *Client {
	return &Client{
		Context:  context,
		Endpoint: endpoint,
		path:     "/pmapi",
		logger:   NewLogger(LOG_INFO),
	}
}

func (c *Client) SetLogLevel(level int) error {
	return c.logger.SetLogLevel(level)
}

func (c *Client) RefreshContext() error {
	c.logger.Debugf("Refreshing context %v", c.Context)
	query, err := c.Context.Query()

	if err != nil {
		return err
	}

	url := fmt.Sprintf(
		"%s%s/%s",
		c.Endpoint,
		c.path,
		query,
	)
	c.logger.Debugf("Generated refresh url: %s", url)

	body, err := c.get(url)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, c.Context); err != nil {
		return err
	}

	return nil
}

func (c *Client) Metrics(query *MetricMetadataQuery) ([]Metric, error) {
	c.logger.Debugln("Fetching metrics for context...")
	result := make(map[string][]Metric)

	body, err := c.getQuery(query)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return result["metrics"], nil
}

func (c *Client) MetricValues(query *MetricValueQuery) (*MetricValueResponse, error) {
	c.logger.Debugln("Fetching metric values....")
	result := MetricValueResponse{}

	body, err := c.getQuery(query)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) getQuery(query Query) ([]byte, error) {
	q, err := query.Query()
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf(
		"%s%s/%d/%s",
		c.Endpoint,
		c.path,
		c.Context.ContextID,
		q,
	)

	return c.get(url)
}

func (c *Client) get(url string) ([]byte, error) {
	resp, err := http.Get(url)

	c.logger.Debugf("Generated url: %s", url)
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
