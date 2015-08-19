package pcp

import (
	"encoding/json"
	"fmt"
	"sort"
)

/*
The PM_* types are taken directly from the PCP source:

https://github.com/performancecopilot/pcp/blob/master/src/include/pcp/pmapi.h#L113

To date, I have only encountered a subset of the below types when interacting
with the pmwebapi. Specifically: FLOAT, 32, U32, U64, DOUBLE, and STRING.

I cannot guarantee the accuracy of the other types 'in the wild', but their
values reflect same naming convention as in those observed.
*/
const (
	PM_TYPE_NOSUPPORT        = "NOSUPPORT"        /* not implemented in this version */
	PM_TYPE_32               = "32"               /* 32-bit signed integer */
	PM_TYPE_U32              = "U32"              /* 32-bit unsigned integer */
	PM_TYPE_64               = "64"               /* 64-bit signed integer */
	PM_TYPE_U64              = "U64"              /* 64-bit unsigned integer */
	PM_TYPE_FLOAT            = "FLOAT"            /* 32-bit floating point */
	PM_TYPE_DOUBLE           = "DOUBLE"           /* 64-bit floating point */
	PM_TYPE_STRING           = "STRING"           /* array of char */
	PM_TYPE_AGGREGATE        = "AGGREGATE"        /* arbitrary binary data (aggregate) */
	PM_TYPE_AGGREGATE_STATIC = "AGGREGATE_STATIC" /* static pointer to aggregate */
	PM_TYPE_EVENT            = "EVENT"            /* packed pmEventArray */
	PM_TYPE_HIGHRES_EVENT    = "HIGHRES_EVENT"    /* packed pmHighResEventArray */
	PM_TYPE_UNKNOWN          = "UNKOWN"           /* used in pmValueBlock, not pmDesc */
	PM_NO_DOMAIN             = -1                 /* this is the instance domain id used to indicate no domain */
	PM_NO_INSTANCE           = -1                 /* instance does not exist */
)

type Metric struct {
	Name        string
	ID          uint32
	Indom       int
	Type        string
	Sem         string
	Units       string
	TextOneline string
	TextHelp    string
}

func (m *Metric) UnmarshalJSON(data []byte) error {
	result := struct {
		Name        *string `json:"name"`
		ID          *uint32 `json:"pmID"`
		Indom       *int    `json:"indom"`
		Type        *string `json:"type"`
		Sem         *string `json:"instant"`
		Units       *string `json:"units"`
		TextOneline *string `json:"text-oneline"`
		TextHelp    *string `json:"text-help"`
	}{}

	json.Unmarshal(data, &result)

	if result.Name != nil {
		m.Name = *result.Name
	} else {
		m.Name = ""
	}

	if result.ID != nil {
		m.ID = *result.ID
	} else {
		m.ID = 0
	}

	if result.Indom != nil {
		m.Indom = *result.Indom
	} else {
		m.Indom = -1
	}

	if result.Type != nil {
		m.Type = *result.Type
	} else {
		m.Type = ""
	}

	if result.Sem != nil {
		m.Sem = *result.Sem
	} else {
		m.Sem = ""
	}

	if result.Units != nil {
		m.Units = *result.Units
	} else {
		m.Units = ""
	}
	if result.TextOneline != nil {
		m.TextOneline = *result.TextOneline
	} else {
		m.TextOneline = ""
	}
	if result.TextHelp != nil {
		m.TextHelp = *result.TextHelp
	} else {
		m.TextHelp = ""
	}
	return nil
}

type MetricList []Metric

func (metrics MetricList) MetricValueType(value *MetricValue) string {
	t := PM_TYPE_UNKNOWN
	i := sort.Search(len(metrics), func(i int) bool {
		return metrics[i].Name >= value.MetricName
	})
	if i != len(metrics) {
		t = metrics[i].Type
	}
	return t
}

func (metrics MetricList) FindMetricByName(name string) *Metric {
	var metric *Metric
	metric = nil

	i := sort.Search(len(metrics), func(i int) bool {
		return metrics[i].Name >= name
	})
	if i != len(metrics) {
		metric = &metrics[i]
	}
	return metric
}

type MetricInstance struct {
	ID    int32 `json:"instance"`
	Name  string
	Value interface{} `json:"value"`
}

type MetricValue struct {
	MetricName string `json:"name"`
	Pmid       uint32 `json:"pmid"`
	Instances  []MetricInstance
}

func (m *MetricValue) UpdateInstanceNames(indom *InstanceDomain) {
	for idx, inst := range m.Instances {
		if inst.ID != PM_NO_INSTANCE {
			i := sort.Search(len(indom.Instances), func(i int) bool {
				return indom.Instances[i].ID >= inst.ID
			})
			if i != len(indom.Instances) {
				m.Instances[idx].Name = indom.Instances[i].Name
			} else {
				fmt.Printf("Failed to find metric name for %v\n", inst)
			}
		} else {
			m.Instances[idx].Name = "UNDEFINED"
		}
	}
}

type InstanceDomainInstance struct {
	ID   int32  `json:"instance"`
	Name string `json:"name"`
}

type InstanceDomain struct {
	ID        int                      `json:"indom"`
	Instances []InstanceDomainInstance `json:"instances"`
}

/* Sorter Methods for Metric lists

Example by metric names:

name := func(m1, m2 *Metric) bool {
	return m1.Name < m2.Name
}
MetricBy(name).Sort(metrics)

*/
type MetricBy func(m1, m2 *Metric) bool

func (by MetricBy) Sort(metrics MetricList) {
	ms := &metricSorter{
		metrics: metrics,
		by:      by,
	}
	sort.Sort(ms)
}

type metricSorter struct {
	metrics MetricList
	by      func(m1, m2 *Metric) bool
}

func (s *metricSorter) Len() int {
	return len(s.metrics)
}

func (s *metricSorter) Swap(i, j int) {
	s.metrics[i], s.metrics[j] = s.metrics[j], s.metrics[i]
}

func (s *metricSorter) Less(i, j int) bool {
	return s.by(&s.metrics[i], &s.metrics[j])
}

/* Sorter methods for Instance Domains
 */
type IDInstanceBy func(d1, d2 *InstanceDomainInstance) bool

func (by IDInstanceBy) Sort(instances []InstanceDomainInstance) {
	ids := &iDInstanceSorter{
		instances: instances,
		by:        by,
	}
	sort.Sort(ids)
}

type iDInstanceSorter struct {
	instances []InstanceDomainInstance
	by        func(d1, d2 *InstanceDomainInstance) bool
}

func (s *iDInstanceSorter) Len() int {
	return len(s.instances)
}

func (s *iDInstanceSorter) Swap(i, j int) {
	s.instances[i], s.instances[j] = s.instances[j], s.instances[i]
}

func (s *iDInstanceSorter) Less(i, j int) bool {
	return s.by(&s.instances[i], &s.instances[j])
}
