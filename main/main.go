package main

import (
	"flag"
	"fmt"
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

	metrics, err := client.Metrics(pcp.NewMetricMetadataQuery(""))

	if err != nil {
		logger.Errorf("Received error retrieving metrics: %s", err)
		os.Exit(1)
	}
	types := make(map[string]bool)
	for _, metric := range metrics {
		types[metric.Type] = true
		logger.Infof("%s\n", metric.Name)
	}
	for key := range types {
		fmt.Println(key)
	}
	logger.Infof("Retrieved %d unique metrics from context", len(metrics))

	// Get values for first 5 metrics by name
	var names []string
	for _, metric := range metrics[:50] {
		names = append(names, metric.Name)
	}

	metric_values_query := pcp.NewMetricValueQuery(names, []string{})
	_, err = client.MetricValues(metric_values_query)
	if err != nil {
		logger.Errorf("Received error retrieving metric values: %s", err)
	}

}
