package minion // import "github.com/oliverbestmann/go-minion"

import (
	"errors"
	"github.com/esailors/go-datadog"
	"github.com/rcrowley/go-metrics"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

type (
	MetricsConfig struct {
		// Specify how often we want to write metrics back
		SampleInterval time.Duration

		// You might want to specify the hostname of the system.
		// This will override hostname auto detection.
		Hostname string

		// Print metrics to the console?
		Console bool

		Datadog DatadogConfig
	}
	DatadogConfig struct {
		ApiKey string

		// Right now, tags are not supported. This is a problem with
		// the Datadog reporter library. Should be added shortly.
		// FIXME support tags!
		Tags []string
	}
)

// Tries to get the hostname of the system by exploiting different
// sources.
func hostname() (string, error) {
	name, err := os.Hostname()
	if err == nil {
		return name, nil
	}

	if name = os.Getenv("HOSTNAME"); name != "" {
		return name, nil
	}

	output, err := exec.Command("hostname").Output()
	if err == nil {
		return strings.TrimSpace(string(output)), nil
	}

	output, err = exec.Command("uname", "-n").Output()
	if err == nil {
		return strings.TrimSpace(string(output)), nil
	}

	return "", errors.New("Could not determine hostname")
}

func SetupMetrics(r metrics.Registry, config MetricsConfig) {
	if r == nil {
		r = metrics.DefaultRegistry
	}

	metrics.RegisterRuntimeMemStats(r)
	go metrics.CaptureRuntimeMemStats(r, config.SampleInterval)

	if config.Console {
		go metrics.WriteJSON(r, config.SampleInterval, os.Stdout)
	}

	if config.Datadog.ApiKey != "" {
		host := config.Hostname
		if host == "" {
			var err error
			if host, err = hostname(); err != nil {
				log.Print("Could not get hostname")
			}
		}

		client := datadog.Client{host, config.Datadog.ApiKey}
		go client.Reporter(r, config.Datadog.Tags).Start(config.SampleInterval)
	}
}
