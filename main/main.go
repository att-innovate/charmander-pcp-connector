package main

import (
	"flag"
	"os"

	"github.com/jameskyle/pcp"
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

	context := pcp.NewContext(hostname)
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

	for _, v := range resp.Values {
		logger.Debugln(pcp.MetricValueType(metrics, v))
	}

	if err != nil {
		logger.Errorf("Received error retrieving metric values: %s\n", err)
	}

	// Fetch indoms
	// First, get the indom for a metric
	indom := metrics[0].Indom

	// Then create a InstanceDomainQuery
	q2 := pcp.NewInstanceDomainQuery(indom)

	// Retrieve results
	indoms_result, err := client.InstanceDomain(q2)
	if err != nil {
		logger.Errorf("Received error retrieving indoms: %s\n", err)
	}

	logger.Debugln(indoms_result)
}
