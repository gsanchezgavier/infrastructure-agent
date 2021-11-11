package filter_test

import (
	"testing"

	"github.com/newrelic/infrastructure-agent/pkg/integrations/v4/filter"
	"github.com/newrelic/infrastructure-agent/pkg/integrations/v4/protocol"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

var integrationJSON = []byte(`
{
	"name": "com.newrelic.redis",
	"protocol_version": "3",
	"integration_version": "0.0.0",
	"data": [
		{
			"metrics": [
				{
					"cluster.connectedSlaves": 0,
					"cluster.role": "master",
					"db.aofLastBgrewriteStatus": "ok",
					"db.aofLastRewriteTimeMiliseconds": -1,
					"db.aofLastWriteStatus": "ok",
					"db.evictedKeysPerSecond": 0,
					"db.expiredKeysPerSecond": 0,
					"db.keyspaceHitsPerSecond": 0,
					"db.keyspaceMissesPerSecond": 0,
					"db.latestForkMilliseconds": 0,
					"db.rdbBgsaveInProgress": 0,
					"db.rdbChangesSinceLastSave": 0,
					"db.rdbLastBgsaveStatus": "ok",
					"db.rdbLastBgsaveTimeMilliseconds": -1,
					"db.rdbLastSaveTime": 1636636361,
					"db.syncFull": 0,
					"db.syncPartialErr": 0,
					"db.syncPartialOk": 0,
					"event_type": "RedisSample",
					"net.blockedClients": 0,
					"net.clientBiggestInputBufBytes": 0,
					"net.clientLongestOutputList": 0,
					"net.commandsProcessedPerSecond": 0,
					"net.connectedClients": 1,
					"net.connectionsReceivedPerSecond": 0,
					"net.inputBytesPerSecond": 0,
					"net.outputBytesPerSecond": 0,
					"net.pubsubChannels": 0,
					"net.pubsubPatterns": 0,
					"net.rejectedConnectionsPerSecond": 0,
					"port": "6379",
					"software.uptimeMilliseconds": 73000,
					"software.version": "4.0.5",
					"system.maxmemoryBytes": 0,
					"system.memFragmentationRatio": 2.88,
					"system.totalSystemMemoryBytes": 8346603520,
					"system.usedCpuSysChildrenMilliseconds": 0,
					"system.usedCpuSysMilliseconds": 340,
					"system.usedCpuUserChildrenMilliseconds": 0,
					"system.usedCpuUserMilliseconds": 60,
					"system.usedMemoryBytes": 828256,
					"system.usedMemoryLuaBytes": 37888,
					"system.usedMemoryPeakBytes": 828256,
					"system.usedMemoryRssBytes": 2387968
				},
				{
					"event_type": "RedisKeyspaceSample",
					"port": "6379"
				}
			],
			"inventory": {},
			"events": []
		}
	]
}
`)

func Test_Filter(t *testing.T) {
	t.Parallel()

	t.Run("no_filter_returns_same_dataset", func(t *testing.T) {
		t.Parallel()
		data, _ := protocol.ParsePayload(integrationJSON, 3)
		f := filter.SampleFilterer{}
		assert.Equal(t, data, f.Filter(data))
	})

	t.Run("drops_metric_from_sample", func(t *testing.T) {
		t.Parallel()
		data, _ := protocol.ParsePayload(integrationJSON, 3)
		f := filter.SampleFilterer{
			SamplesKeys: filter.SampleKeysFilter{
				Exclude: filter.SampleKeysMatcher{
					SampleNames: filter.Matcher{
						"RedisKeyspaceSample",
					},
					MetricNames: filter.Matcher{
						"port",
					},
				},
			},
		}

		filteredData := f.Filter(data)
		assert.Contains(t, filteredData.DataSets[0].Metrics[0], "port")
		assert.NotContains(t, filteredData.DataSets[0].Metrics[1], "port")
	})

	t.Run("drops_sample", func(t *testing.T) {
		t.Parallel()
		data, _ := protocol.ParsePayload(integrationJSON, 3)
		f := filter.SampleFilterer{
			Samples: filter.SamplesFilter{
				Exclude: filter.SampleMatcher{
					SampleNames: filter.Matcher{
						"RedisKeyspaceSample",
					},
				},
			},
		}
		assert.Len(t, data.DataSets[0].Metrics, 2)
		filteredData := f.Filter(data)
		assert.Len(t, filteredData.DataSets[0].Metrics, 1)
	})
}

func Test_Config(t *testing.T) {
	t.Parallel()

	config := []byte(`
sample_keys:
  exclude:
    sample_names:
    - RedisSample
    metric_names:
    - "system.*"
samples:
  exclude:
    sample_names:
      - RedisKeyspaceSample

`)
	c := filter.SampleFilterer{}

	assert.NoError(t, yaml.Unmarshal(config, &c))
}
