package pcp

import (
	"errors"
	"net/url"
	"strings"
)

type Query interface {
	Query() (string, error)
}

type MetricMetadataQuery struct {
	Prefix string
}

type MetricValueQuery struct {
	Names []string
	Pmids []string
}

type TimeStamp struct {
	Seconds      uint64 `json:"s"`
	MicroSeconds uint64 `json:"us"`
}

type MetricValueResponse struct {
	Timestamp *TimeStamp
	Values    []MetricValue
}

func NewMetricMetadataQuery(prefix string) *MetricMetadataQuery {
	return &MetricMetadataQuery{Prefix: prefix}
}

func (m *MetricMetadataQuery) Query() (string, error) {
	query := "_metric"
	u := url.Values{}

	if m.Prefix != "" {
		u.Set("prefix", m.Prefix)
		query += "?" + u.Encode()
	}
	return query, nil
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
