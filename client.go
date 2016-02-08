package pcp

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Client struct {
	Context *Context
	Host    string
	Port    uint16
	logger  *Logger
	path    string
}

func NewClient(host string, port uint16, context *Context) *Client {
	return &Client{
		Context: context,
		Host:    host,
		Port:    port,
		path:    "/pmapi",
		logger:  NewLogger(LOG_INFO),
	}
}

func (c *Client) Endpoint() string {
	return fmt.Sprintf("http://%s:%d", c.Host, c.Port)
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
		c.Endpoint(),
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

func (c *Client) Metrics(query *MetricQuery) (MetricList, error) {
	c.logger.Debugln("Fetching metrics for context...")
	result := make(map[string]MetricList)
	var metrics MetricList

	body, err := c.getQuery(query)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	metrics = result["metrics"]
	name := func(m1, m2 *Metric) bool {
		return m1.Name < m2.Name
	}
	MetricBy(name).Sort(metrics)
	return metrics, nil
}

func (c *Client) MetricValues(query *MetricValueQuery) (*MetricValueResponse, error) {
	c.logger.Debugln("Fetching metric values....")
	result := MetricValueResponse{}

	body, err := c.getQuery(query)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(body, &result); err != nil {
	    switch e := err.(type) {
	    case *json.UnmarshalTypeError:
                c.logger.Debugln("UnmarshalTypeError: Value[%s] Type[%v]\n", e.Value, e.Type)
	    case *json.InvalidUnmarshalError:
		c.logger.Debugln("InvalidUnmarshalError: Type[%v]\n", e.Type)
	    default:
	        c.logger.Debugln(err)
	    }
	    return nil, err
	}
	return &result, nil
}

func (c *Client) InstanceDomain(query *InstanceDomainQuery) (*InstanceDomain, error) {
	c.logger.Debugln("Fetching instance domains...")
	var indom InstanceDomain

	body, err := c.getQuery(query)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(body, &indom); err != nil {
		return nil, err
	}
	ids := func(d1, d2 *InstanceDomainInstance) bool {
		return d1.ID < d2.ID
	}
	IDInstanceBy(ids).Sort(indom.Instances)
	return &indom, nil
}

func (c *Client) getQuery(query Querier) ([]byte, error) {
	q, err := query.Query()
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf(
		"%s%s/%d/%s",
		c.Endpoint(),
		c.path,
		c.Context.ContextID,
		q,
	)
	c.logger.Debugf("Executing Query: %s\n", url)
	return c.get(url)
}

func (c *Client) get(url string) ([]byte, error) {
	c.logger.Debugf("Generated url: %s", url)
	resp, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	c.logger.Debugf("Query Raw Result: %s", string(body))
	return body, nil
}

func (c *Client) GetIndomForMetric(metric *Metric) (*InstanceDomain, error) {
	var indom *InstanceDomain
	var err error
	indom = &InstanceDomain{ID: PM_NO_DOMAIN}

	if metric.Indom != PM_NO_DOMAIN {
		query := NewInstanceDomainQuery(metric.Indom)
		indom, err = c.InstanceDomain(query)

		if err != nil {
			return nil, err
		}
	}

	return indom, nil
}
