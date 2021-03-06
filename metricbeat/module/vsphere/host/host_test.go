package host

import (
	"testing"

	"github.com/elastic/beats/libbeat/common"
	mbtest "github.com/elastic/beats/metricbeat/mb/testing"

	"github.com/stretchr/testify/assert"
	"github.com/vmware/vic/pkg/vsphere/simulator"
)

func TestFetchEventContents(t *testing.T) {

	model := simulator.ESX()

	err := model.Create()
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	ts := model.Service.NewServer()
	defer ts.Close()

	urlSimulator := ts.URL.Scheme + "://" + ts.URL.Host + ts.URL.Path

	config := map[string]interface{}{
		"module":     "vsphere",
		"metricsets": []string{"host"},
		"hosts":      []string{urlSimulator},
		"username":   "user",
		"password":   "pass",
		"insecure":   true,
	}

	f := mbtest.NewEventsFetcher(t, config)

	events, err := f.Fetch()

	event := events[0]
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	t.Logf("%s/%s event: %+v", f.Module().Name(), f.Name(), event.StringToPrint())

	assert.EqualValues(t, "ha-datacenter", event["datacenter"])
	assert.EqualValues(t, "localhost.localdomain", event["name"])

	cpu := event["cpu"].(common.MapStr)

	cpuUsed := cpu["used"].(common.MapStr)
	assert.EqualValues(t, 67, cpuUsed["mhz"])

	cpuTotal := cpu["total"].(common.MapStr)
	assert.EqualValues(t, 4588, cpuTotal["mhz"])

	cpuFree := cpu["free"].(common.MapStr)
	assert.EqualValues(t, 4521, cpuFree["mhz"])

	memory := event["memory"].(common.MapStr)

	memoryUsed := memory["used"].(common.MapStr)
	assert.EqualValues(t, uint64(1472200704), memoryUsed["bytes"])

	memoryTotal := memory["total"].(common.MapStr)
	assert.EqualValues(t, uint64(4294430720), memoryTotal["bytes"])

	memoryFree := memory["free"].(common.MapStr)
	assert.EqualValues(t, uint64(2822230016), memoryFree["bytes"])

}
