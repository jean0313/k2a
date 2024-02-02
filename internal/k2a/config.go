package k2a

import (
	"fmt"
	"strings"
)

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
	TopicOperationMap map[string]string
}

func (c *K2AConfig) GetTopics() []string {
	keys := make([]string, 0, len(c.TopicOperationMap))
	for k := range c.TopicOperationMap {
		keys = append(keys, k)
	}
	return keys
}

func (c *K2AConfig) InitTopicOperations() error {
	ops := strings.Split(c.Topics, ",")
	c.TopicOperationMap = make(map[string]string)
	for _, op := range ops {
		topic, operation, err := extractTopicAndOp(op)
		if err != nil {
			return err
		}
		c.TopicOperationMap[topic] = operation
	}
	return nil
}

func extractTopicAndOp(param string) (string, string, error) {
	arr := strings.Split(param, ":")
	if len(arr) == 1 {
		return arr[0], "Publish", nil
	}

	topic := arr[0]
	op := arr[1]
	if op == "p" {
		return topic, "Publish", nil
	} else if op == "s" {
		return topic, "Subscribe", nil
	}
	return "", "", fmt.Errorf("operation [%s] not support, only support 's' and 'p', s is for subscribe, p is for publish ", op)
}
