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
	InstanceDomain int
	Name           string
	Instances      []uint32
	INames         []string
}

func NewInstanceDomainQuery(id int) *InstanceDomainQuery {
	return &InstanceDomainQuery{InstanceDomain: id}
}

func (id *InstanceDomainQuery) Query() (string, error) {
	if id.InstanceDomain <= PM_NO_DOMAIN && id.Name == "" {
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
	if id.InstanceDomain > 0 {
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

type MetricValueResponseList []*MetricValueResponse

func (m MetricValueResponseList) InstanceNames() []string {
	names := []string{}
	nameMap := make(map[string]bool)

	for _, response := range m {
		for _, value := range response.Values {
			for _, inst := range value.Instances {
				nameMap[inst.Name] = true
			}
		}
	}
	for key, _ := range nameMap {
		names = append(names, key)
	}
	return names
}

func (m MetricValueResponseList) InstanceFilter(
	match func(instance MetricInstance) bool,
) map[string][]MetricInstance {

	result := make(map[string][]MetricInstance)

	for _, response := range m {
		for _, value := range response.Values {
			for _, inst := range value.Instances {
				if match(inst) {
					result[value.MetricName] = append(result[value.MetricName], inst)
				}
			}

		}
	}
	return result
}

func (m MetricValueResponseList) MetricNames() []string {
	result := []string{}
	mapResult := make(map[string]struct{})

	for _, response := range m {
		for _, value := range response.Values {
			mapResult[value.MetricName] = struct{}{}
		}
	}

	for name, _ := range mapResult {
		result = append(result, name)
	}

	return result
}

func (m MetricValueResponseList) MetricValueByInstance() map[string]map[string]interface{} {
	metricsMap := make(map[string]map[string]interface{})

	for _, response := range m {
		for _, value := range response.Values {
			for _, inst := range value.Instances {
				if _, ok := metricsMap[inst.Name]; !ok {
					metricsMap[inst.Name] = make(map[string]interface{})
				}
				metricsMap[inst.Name][value.MetricName] = inst.Value

				// Set the timestamp, only once
				if _, ok := metricsMap[inst.Name]["time"]; !ok {
					metricsMap[inst.Name]["time"] = response.Timestamp.Seconds
				}
			}
		}
	}
	return metricsMap
}
