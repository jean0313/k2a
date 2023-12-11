package k2a

import "strings"

const DEFAULT_KAFKA_URL = "localhost:9092"
const DEFAULT_SCHEMA_REGISTRY_URL = "http://localhost:8081"

type K2AConfig struct {
	KafkaUrl          string
	SchemaRegistryUrl string
	Topics            string
	File              string
	SpecVersion       string
	Certificate       string
	KeyFile           string
	CAFile            string
	TLSSkipVerify     bool
	UseTLS            bool
}

func (c *K2AConfig) GetTopics() []string {
	return strings.Split(c.Topics, ",")
}

func (c *K2AConfig) WithDefaults() {
	if c.KafkaUrl == "" {
		c.KafkaUrl = DEFAULT_KAFKA_URL
	}

	if c.SchemaRegistryUrl == "" {
		c.SchemaRegistryUrl = DEFAULT_SCHEMA_REGISTRY_URL
	}

	if c.SpecVersion == "" {
		c.SpecVersion = "1.0.0"
	}
}
