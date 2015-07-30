package pcp

import "sort"

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
)

type Metric struct {
	Name        string `json:"name"`
	ID          uint32 `json:"pmID"`
	Indom       uint32 `json:"indom"`
	Type        string `json:"type"`
	Sem         string `json:"instant"`
	Units       string `json:"units"`
	TextOneline string `json:"text-oneline"`
	TextHelp    string `json:"text-help"`
}

type MetricInstance struct {
	ID    int32       `json:"instance"`
	Value interface{} `json:"value"`
}

type MetricValue struct {
	Name      string `json:"name"`
	Pmid      uint32 `json:"pmid"`
	Instances []MetricInstance
}

type InstanceDomainInstance struct {
	ID   int32  `json:"instance"`
	Name string `json:"name"`
}

type InstanceDomain struct {
	ID        uint32                   `json:"indom"`
	Instances []InstanceDomainInstance `json:"instances"`
}

func MetricValueType(metrics []Metric, value *MetricValue) string {
	t := PM_TYPE_UNKNOWN
	i := sort.Search(len(metrics), func(i int) bool {
		return metrics[i].Name >= value.Name
	})
	if i != len(metrics) {
		t = metrics[i].Type
	}
	return t
}

/* Sorter Methods for Metric lists

Example by metric names:

name := func(m1, m2 *Metric) bool {
	return m1.Name < m2.Name
}
By(name).Sort(metrics)

*/
type By func(m1, m2 *Metric) bool

func (by By) Sort(metrics []Metric) {
	ms := &metricSorter{
		metrics: metrics,
		by:      by,
	}
	sort.Sort(ms)
}

type metricSorter struct {
	metrics []Metric
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
