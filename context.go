package pcp

import (
	"fmt"
	"net/url"
)

type Context struct {
	Hostname    string
	Hostspec    string
	Local       string
	Username    string
	Password    string
	Archivefile string
	PollTimeout int32
	ContextID   uint32 `json:"context"`
}

func NewContext(hostname string) *Context {
	c := &Context{
		Hostname: hostname,
	}

	return c
}

func (c *Context) setIfNotEmpty(v *url.Values, key string, value string) {
	if value != "" {
		v.Set(key, value)
	}
}

func (c *Context) params() string {
	v := url.Values{}

	v.Set("hostname", c.Hostname)
	c.setIfNotEmpty(&v, "hostspec", c.Hostspec)
	c.setIfNotEmpty(&v, "local", c.Local)
	c.setIfNotEmpty(&v, "archivefile", c.Archivefile)
	if c.PollTimeout != 0 {
		v.Set("polltimeout", fmt.Sprintf("%d", c.PollTimeout))
	}
	return v.Encode()
}
