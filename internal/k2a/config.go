package k2a

import "strings"

const DEFAULT_KAFKA_URL = "localhost:9092"
const DEFAULT_SCHEMA_REGISTRY_URL = "http://localhost:8081"
const DEFAULT_FILE_FORMAT_YAML = "yaml"

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
	FileFormat        string
}

func (c *K2AConfig) GetTopics() []string {
	return strings.Split(c.Topics, ",")
}
