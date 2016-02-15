package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/att-innovate/charmander-pcp"
)

var (
	endpoint string
	hostname string
	logger   *pcp.Logger
)

func init() {
	logger = pcp.NewLogger(pcp.LOG_DEBUG)
	flag.StringVar(&endpoint, "endpoint", "", "endpoint to retrieve metrics")
	flag.StringVar(&hostname, "hostname", "local", "hostname to retrieve metrics")
	flag.Parse()
	validate()
}

func validate() {
	if endpoint == "" {
		logger.Errorln("You must provide an endpoint for the pcp client!")
		flag.Usage()
		os.Exit(1)
	}
}

func main() {
	logger.Infoln("Starting...")

	context := pcp.NewContext(hostname, "")
	client := pcp.NewClient(endpoint, context)
	client.SetLogLevel(pcp.LOG_DEBUG)
	// var query pcp.Query

	err := client.RefreshContext()
	logger.Debugln(client.Context)

	if err != nil {
		logger.Errorf("Received error refreshing context: %s", err)
		os.Exit(1)
	}

	metrics, err := client.Metrics(pcp.NewMetricQuery(""))

	if err != nil {
		logger.Errorf("Received error retrieving metrics: %s", err)
		os.Exit(1)
	}

	logger.Infof("Retrieved %d unique metrics from context", len(metrics))

	// Get values for first 5 metrics by name
	var names []string
	for _, metric := range metrics[:50] {
		names = append(names, metric.Name)
	}

	metric_values_query := pcp.NewMetricValueQuery(names, []string{})
	resp, err := client.MetricValues(metric_values_query)

	logger.Debugln(metrics.MetricValueType(resp.Values[0]))

	if err != nil {
		logger.Errorf("Received error retrieving metric values: %s\n", err)
	}

	// update the metric values with their metric names
	for _, value := range resp.Values {
		metric := metrics.FindMetricByName(value.MetricName)
		indom, err := client.GetIndomForMetric(metric)
		logger.Debugln(indom)
		if err != nil {
			logger.Errorf("Failed to find Instance Domain for metric: %s", err)
		}
		value.UpdateInstanceNames(indom)
		logger.Debugln(value)
	}

	// Get all stats for containers on a host.
	q3 := pcp.NewMetricValueQuery([]string{"containers.name"}, []string{})
	containers, err := client.MetricValues(q3)
	if err != nil {
		logger.Errorf("Container query failed: %s", err)
	}
	// /pmapi/_context?hostspec=local:?container=fooba
	// get a name
	cname := containers.Values[0].Instances[0].Value
	// create a new context
	spec := fmt.Sprintf("local:?container=%s", cname)
	c_context := pcp.NewContext("", spec)
	c_client := pcp.NewClient(endpoint, c_context)
	c_client.SetLogLevel(pcp.LOG_DEBUG)
	err = c_client.RefreshContext()
	if err != nil {
		logger.Errorf("Query failed with error: %s", err)
	}
	names = []string{
		"cgroup.cpuacct.stat.user",
		"cgroup.cpuacct.stat.system",
		"cgroup.memory.usage",
	}
	for _, name := range names {
		c_query := pcp.NewMetricValueQuery([]string{name}, []string{})
		resp, err = c_client.MetricValues(c_query)

		if err != nil {
			logger.Errorf("Query failed with error: %s", err)
		}
		logger.Debugln(resp.Values)
	}

}
