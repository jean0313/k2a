package k2a

import "strings"

type K2AConfig struct {
	KafkaUrl          string
	SchemaRegistryUrl string
	Topics            string
	File              string
	SpecVersion       string
}

func (c *K2AConfig) GetTopics() []string {
	return strings.Split(c.Topics, ",")
}
