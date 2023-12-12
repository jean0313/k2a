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
	UserName          string
	Password          string
}

func (c *K2AConfig) GetTopics() []string {
	return strings.Split(c.Topics, ",")
}
