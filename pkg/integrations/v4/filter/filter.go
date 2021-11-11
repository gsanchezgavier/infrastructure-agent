package filter

import (
	"fmt"

	"github.com/newrelic/infrastructure-agent/pkg/integrations/v4/protocol"
)

type Config struct {
	SampleFilterer `yaml:",inline" json:",inline"`
}

type SampleFilterer struct {
	Samples     SamplesFilter    `yaml:"samples" json:"samples"`
	SamplesKeys SampleKeysFilter `yaml:"sample_keys" json:"sample_keys"`
}

func (s SampleFilterer) Filter(data protocol.PluginDataV3) protocol.PluginDataV3 {
	// for i, pluginDataSet := range data.DataSets {
	// 	newMetrics := []protocol.MetricData{}

	// 	for _, metrics := range pluginDataSet.Metrics {
	// 		// include logic not implemented
	// 		sampleName := fmt.Sprintf("%v", metrics["event_type"])

	// 		if s.Samples.Exclude.match(fmt.Sprintf("%v", sampleName)) {
	// 			break
	// 		}

	// 		for key := range metrics {
	// 			if s.SamplesKeys.SampleName.match(sampleName) &&
	// 				s.SamplesKeys.KeyName.Exclude.match(key) {
	// 				delete(metrics, key)
	// 			}
	// 		}

	// 		newMetrics = append(newMetrics, metrics)
	// 	}
	// 	data.DataSets[i].Metrics = newMetrics
	// }
	for i, pluginDataSet := range data.DataSets {
		newMetrics := []protocol.MetricData{}

		for _, metrics := range pluginDataSet.Metrics {
			// include logic not implemented
			sampleName := fmt.Sprintf("%v", metrics["event_type"])

			if s.Samples.Exclude.SampleNames.match(fmt.Sprintf("%v", sampleName)) {
				break
			}

			for key := range metrics {
				if s.SamplesKeys.Exclude.SampleNames.match(sampleName) &&
					s.SamplesKeys.Exclude.MetricNames.match(key) {
					delete(metrics, key)
				}
			}

			newMetrics = append(newMetrics, metrics)
		}
		data.DataSets[i].Metrics = newMetrics
	}

	return data
}

type SampleKeysFilter struct {
	Exclude SampleKeysMatcher `yaml:"exclude" json:"exclude"`
}

type SampleKeysMatcher struct {
	// SampleNamesList []string `yaml:"sample_names" json:"sample_names"`
	// MetricNamesList []string `yaml:"metric_names" json:"metric_names"`
	SampleNames Matcher `yaml:"sample_names" json:"sample_names"`
	MetricNames Matcher `yaml:"metric_names" json:"metric_names"`
	matchType   string  `yaml:"match_type" json:"match_type"`
}

type SamplesFilter struct {
	Exclude SampleMatcher `yaml:"exclude" json:"exclude"`
}

type SampleMatcher struct {
	// SampleNamesList []string `yaml:"sample_names" json:"sample_names"`
	SampleNames Matcher `yaml:"sample_names" json:"sample_names"`
	matchType   string  `yaml:"match_type" json:"match_type"`
}

type Matcher []string

// type Matcher struct {
// 	// Patterns []*regexp.Regexp
// 	Patterns []string
// }

func (m Matcher) match(s string) bool {
	for _, p := range m {
		if p == s {
			return true
		}
	}
	return false
}

// func (m Matcher) match(s string) bool {
// 	for _, p := range m.Patterns {
// 		if p.MatchString(s) {
// 			return true
// 		}
// 	}
// 	return false
// }
