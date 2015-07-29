package main

import (
	"flag"
	"os"
	"strings"

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
	client := pcp.NewClient(endpoint)
	client.SetLogLevel(pcp.LOG_DEBUG)

	err := client.RefreshContext(context)
	logger.Debugln(context)
	if err != nil {
		logger.Errorf("Received error refreshing context: %s", err)
		os.Exit(1)
	}

	metrics, err := client.Metrics(context, "")

	if err != nil {
		logger.Errorf("Received error retrieving metrics: %s", err)
		os.Exit(1)
	}
	logger.Infof("Retrieved %d unique metrics from contenxt", len(metrics))

	wanted := strings.Join(func() []string {
		names := []string{}
		for _, metric := range metrics[:5] {
			names = append(names, metric.Name)
		}
		return names
	}(), ",")

	logger.Infoln(wanted)

}
