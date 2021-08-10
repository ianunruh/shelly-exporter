package collect

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"

	"github.com/ianunruh/shelly-exporter/pkg/shelly"
)

func TestCollector(t *testing.T) {
	log := zaptest.NewLogger(t)

	status, err := loadStatus()
	require.NoError(t, err)

	client := &mockClient{
		status: status,
	}

	collector := NewCollector(client, log)

	descCh := make(chan *prometheus.Desc, 12)
	collector.Describe(descCh)
	close(descCh)

	var descriptors []*prometheus.Desc
	for desc := range descCh {
		descriptors = append(descriptors, desc)
	}
	assert.Len(t, descriptors, 12)

	metricCh := make(chan prometheus.Metric, 12)
	collector.Collect(metricCh)
	close(metricCh)

	var metrics []prometheus.Metric
	for metric := range metricCh {
		fmt.Println(metric)
		metrics = append(metrics, metric)
	}
	assert.Len(t, metrics, 12)
}

type mockClient struct {
	status *shelly.Status
}

func (c mockClient) Status() (*shelly.Status, error) {
	return c.status, nil
}

func loadStatus() (*shelly.Status, error) {
	encoded, err := ioutil.ReadFile("testdata/status.json")
	if err != nil {
		return nil, err
	}

	var status shelly.Status
	if err := json.Unmarshal(encoded, &status); err != nil {
		return nil, err
	}

	return &status, nil
}
