package pcp

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

/*
Querier provides an interface for generating PCP Queries. It must provide the
Query() method.

Query()
*/
type Querier interface {
	Query() (string, error)
}

type MetricQuery struct {
	Prefix string
}

func NewMetricQuery(prefix string) *MetricQuery {
	return &MetricQuery{Prefix: prefix}
}

func (m *MetricQuery) Query() (string, error) {
	query := "_metric"
	u := url.Values{}

	if m.Prefix != "" {
		u.Set("prefix", m.Prefix)
		query += "?" + u.Encode()
	}
	return query, nil
}

type MetricValueQuery struct {
	Names []string
	Pmids []string
}

func NewMetricValueQuery(names []string, pmids []string) *MetricValueQuery {
	return &MetricValueQuery{Names: names, Pmids: pmids}
}

func (m *MetricValueQuery) Query() (string, error) {
	if len(m.Names) == 0 && len(m.Pmids) == 0 {
		e := errors.New("You must provide at least one PMID or Name for the query!")
		return "", e
	}

	u := url.Values{}
	names := strings.Join(m.Names, ",")
	pmids := strings.Join(m.Pmids, ",")

	if names != "" {
		u.Set("names", names)
	}
	if pmids != "" {
		u.Set("pmids", pmids)
	}

	query := "_fetch" + "?" + u.Encode()
	return query, nil
}

type InstanceDomainQuery struct {
	InstanceDomain uint32
	Name           string
	Instances      []uint32
	INames         []string
}

func NewInstanceDomainQuery(id uint32) *InstanceDomainQuery {
	return &InstanceDomainQuery{InstanceDomain: id}
}

func (id *InstanceDomainQuery) Query() (string, error) {
	if id.InstanceDomain == 0 && id.Name == "" {
		e := errors.New("You must provide at least one InstanceDomain or Name for the query!")
		return "", e
	}

	u := url.Values{}
	instances := strings.Join(func() []string {
		instances := []string{}
		for _, inst := range id.Instances {
			instances = append(instances, fmt.Sprintf("%d", inst))
		}
		return instances
	}(), ",")
	inames := strings.Join(id.INames, ",")

	if inames != "" {
		u.Set("iname", inames)
	}
	if instances != "" {
		u.Set("instance", instances)
	}
	if id.InstanceDomain != 0 {
		u.Set("indom", fmt.Sprintf("%d", id.InstanceDomain))
	}
	if id.Name != "" {
		u.Set("name", id.Name)
	}
	query := "_indom" + "?" + u.Encode()
	return query, nil
}

type TimeStamp struct {
	Seconds      uint64 `json:"s"`
	MicroSeconds uint64 `json:"us"`
}

type MetricValueResponse struct {
	Timestamp *TimeStamp
	Values    []*MetricValue
}
