package k2a

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/IBM/sarama"
	"github.com/go-resty/resty/v2"
	"github.com/iancoleman/strcase"
	"github.com/swaggest/go-asyncapi/spec-2.4.0"
)

type Topic struct {
	NumPartitions     int32
	ReplicationFactor int16
	TopicName         string
	ConfigEntries     map[string]*string
}

func (t *Topic) GetTopicName() string {
	if t == nil {
		var ret string
		return ret
	}
	return t.TopicName
}

type Schema struct {
	Subject    *string `json:"subject,omitempty"`
	Version    *int32  `json:"version,omitempty"`
	Id         *int32  `json:"id,omitempty"`
	SchemaType *string `json:"schemaType,omitempty"`
	Schema     *string `json:"schema,omitempty"`
}

func (s *Schema) GetSchemaType() string {
	if s == nil || s.SchemaType == nil {
		var ret string
		return ret
	}
	return *s.SchemaType
}

func (s *Schema) GetSchema() string {
	if s == nil || s.Schema == nil {
		var ret string
		return ret
	}
	return *s.Schema
}

type ChannelDetails struct {
	topicName          string
	topicDescription   string
	subject            string
	contentType        string
	unmarshalledSchema map[string]any
	schema             *Schema
}

type topicConfigurationExport struct {
	CleanupPolicy       []string `json:"cleanup.policy,omitempty"`
	RetentionTime       int64    `json:"retention.ms,omitempty"`
	RetentionSize       int64    `json:"retention.bytes,omitempty"`
	DeleteRetentionTime int64    `json:"delete.retention.ms,omitempty"`
	MaxMessageSize      int32    `json:"max.message.bytes,omitempty"`
}

type bindings struct {
	channelBindings  spec.ChannelBindingsObject
	messageBinding   spec.MessageBindingsObject
	operationBinding spec.OperationBindingsObject
}

type AccountDetails struct {
	kafkaClusterId          string
	schemaRegistryClusterId string
	topics                  []Topic
	kafkaUrl                string
	schemaRegistryUrl       string
	subjects                []string
	channelDetails          channelDetails
}

func (a *AccountDetails) QuerySubjectSchema() (Schema, error) {
	var nilS Schema
	client := resty.New()
	resp, err := client.R().Get(fmt.Sprintf("%s/subjects/%s/versions/latest", a.schemaRegistryUrl, a.channelDetails.currentSubject))
	if err != nil {
		return nilS, err
	}

	if resp.IsError() {
		return nilS, errors.New(resp.String())
	}

	schema := Schema{}
	if err = json.Unmarshal([]byte(resp.Body()), &schema); err != nil {
		return nilS, err
	}
	return schema, nil
}

func (a *AccountDetails) queryTopicInfo(topics []string) ([]Topic, error) {
	config := sarama.NewConfig()
	brokers := []string{a.kafkaUrl}

	admin, err := sarama.NewClusterAdmin(brokers, config)
	if err != nil {
		return nil, err
	}

	tps, err1 := admin.ListTopics()
	if err1 != nil {
		return nil, err1
	}

	ret := make([]Topic, 0)
	for _, topic := range topics {
		value, ok := tps[topic]
		if ok {
			ret = append(ret, Topic{
				value.NumPartitions,
				value.ReplicationFactor,
				topic,
				value.ConfigEntries,
			})
		}
	}
	return ret, nil
}

func (a *AccountDetails) querySchemaRegistrySubjects() ([]string, error) {
	client := resty.New()
	resp, err := client.R().Get(fmt.Sprintf("%s/subjects", a.schemaRegistryUrl))
	if err != nil {
		return nil, err
	}

	var ret []string
	if err = json.Unmarshal(resp.Body(), &ret); err != nil {
		return nil, err
	}
	return ret, nil
}

func (a *AccountDetails) buildMessageEntity() *spec.MessageEntity {
	entityProducer := new(spec.MessageEntity)
	entityProducer.WithContentType(a.channelDetails.contentType)
	if a.channelDetails.contentType == "application/avro" {
		entityProducer.WithSchemaFormat("application/vnd.apache.avro;version=1.9.0")
	} else if a.channelDetails.contentType == "application/json" {
		(*spec.MessageEntity).WithSchemaFormat(entityProducer, "application/schema+json;version=draft-07")
	}
	entityProducer.WithTags(a.channelDetails.schemaLevelTags...)
	// Name
	entityProducer.WithName(msgName(a.channelDetails.currentTopic.GetTopicName()))
	if a.channelDetails.bindings != nil {
		entityProducer.WithBindings(a.channelDetails.bindings.messageBinding)
	}
	if a.channelDetails.unmarshalledSchema != nil {
		entityProducer.WithPayload(a.channelDetails.unmarshalledSchema)
	}
	return entityProducer
}

func msgName(s string) string {
	return strcase.ToCamel(s) + "Message"
}

func (a *AccountDetails) getSchemaDetails() error {
	schema, err := a.QuerySubjectSchema()
	if err != nil {
		return err
	}

	a.channelDetails.schema = &schema
	if schema.GetSchemaType() == "" {
		schema.SchemaType = ptrString("AVRO")
	}

	switch schema.GetSchemaType() {
	case "PROTOBUF":
		return fmt.Errorf(protobufErrorMessage)
	case "AVRO", "JSON":
		a.channelDetails.contentType = fmt.Sprintf("application/%s", strings.ToLower(schema.GetSchemaType()))
	}

	if err := json.Unmarshal([]byte(*schema.Schema), &a.channelDetails.unmarshalledSchema); err != nil {
		a.channelDetails.unmarshalledSchema, err = handlePrimitiveSchemas(schema.GetSchema(), err)
	}
	return nil
}

func handlePrimitiveSchemas(schema string, err error) (map[string]any, error) {
	unmarshalledSchema := make(map[string]any)
	primitiveTypes := []string{"string", "null", "boolean", "int", "long", "float", "double", "bytes"}
	for _, primitiveType := range primitiveTypes {
		if schema == fmt.Sprintf(`"%s"`, primitiveType) {
			unmarshalledSchema["type"] = primitiveType
			return unmarshalledSchema, nil
		}
	}
	return nil, fmt.Errorf("failed to unmarshal schema: %w", err)
}

func ptrString(s string) *string {
	return &s
}

type channelDetails struct {
	currentTopic            Topic
	currentTopicDescription string
	currentSubject          string
	contentType             string
	schema                  *Schema
	unmarshalledSchema      map[string]any
	mapOfMessageCompat      map[string]any
	topicLevelTags          []spec.Tag
	schemaLevelTags         []spec.Tag
	bindings                *bindings
}

type ChannelBinding struct {
	BindingVersion     string                   `json:"bindingVersion,omitempty"`
	Partitions         int32                    `json:"partitions,omitempty"`
	TopicConfiguration topicConfigurationExport `json:"topicConfiguration,omitempty"`
	XConfigs           map[string]string        `json:"x-configs,omitempty"`
}
