package k2a

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/swaggest/go-asyncapi/reflector/asyncapi-2.4.0"
	"github.com/swaggest/go-asyncapi/spec-2.4.0"
)

const protobufErrorMessage = "protobuf is not supported"

func ExportAsyncApi(config *K2AConfig) ([]byte, error) {
	details, err := GetAccountDetails(config)
	if err != nil {
		return nil, err
	}

	reflector := createAsyncReflector(config)
	messages := make(map[string]spec.Message)
	channelCount := 0

	for _, topic := range details.topics {
		for _, subject := range details.subjects {
			if subject != fmt.Sprintf("%s-value", topic.TopicName) || strings.HasPrefix(topic.TopicName, "_") {
				continue
			}

			details.channelDetails = channelDetails{
				currentTopic:   topic,
				currentSubject: subject,
			}

			if err := processChannelDetails(details); err != nil {
				if err.Error() == protobufErrorMessage {
					continue
				}
				return nil, err
			}

			channelCount++
			messages[*details.channelDetails.schema.Name] = spec.Message{
				OneOf1: &spec.MessageOneOf1{MessageEntity: details.buildMessageEntity()},
			}
			reflector, err = addChannel(reflector, details.channelDetails)
			if err != nil {
				return nil, err
			}
		}
	}

	if channelCount == 0 {
		reflector.Schema.Channels = map[string]spec.ChannelItem{}
	}

	addComponents(reflector, messages)

	if config.FileFormat == DEFAULT_FILE_FORMAT_YAML {
		return reflector.Schema.MarshalYAML()
	} else {
		return reflector.Schema.MarshalJSON()
	}
}

func addComponents(reflector asyncapi.Reflector, messages map[string]spec.Message) asyncapi.Reflector {
	reflector.Schema.WithComponents(spec.Components{
		Messages: messages,
		SecuritySchemes: &spec.ComponentsSecuritySchemes{
			MapOfComponentsSecuritySchemesWDValues: map[string]spec.ComponentsSecuritySchemesWD{
				"schemaRegistry": {
					SecurityScheme: &spec.SecurityScheme{
						UserPassword: &spec.UserPassword{
							MapOfAnything: map[string]any{
								"x-configs": any(map[string]string{
									"basic.auth.user.info": "{{SCHEMA_REGISTRY_API_KEY}}:{{SCHEMA_REGISTRY_API_SECRET}}",
								}),
							},
						},
					},
				},
				"broker": {
					SecurityScheme: &spec.SecurityScheme{
						UserPassword: &spec.UserPassword{
							MapOfAnything: map[string]any{
								"x-configs": any(map[string]string{
									"security.protocol": "sasl_ssl",
									"sasl.mechanisms":   "PLAIN",
									"sasl.username":     "{{CLUSTER_API_KEY}}",
									"sasl.password":     "{{CLUSTER_API_SECRET}}",
								}),
							},
						},
					},
				},
			},
		},
	})
	return reflector
}

func addChannel(reflector asyncapi.Reflector, details channelDetails) (asyncapi.Reflector, error) {
	channel := asyncapi.ChannelInfo{
		Name: details.currentTopic.GetTopicName(),
		BaseChannelItem: &spec.ChannelItem{
			Description: details.currentTopicDescription,
			Subscribe: &spec.Operation{
				ID:   strcase.ToCamel(details.currentTopic.GetTopicName()) + "Subscribe",
				Tags: details.topicLevelTags,
			},
		},
	}
	if details.mapOfMessageCompat != nil {
		channel.BaseChannelItem.MapOfAnything = details.mapOfMessageCompat
	}
	if details.unmarshalledSchema != nil {
		channel.BaseChannelItem.Subscribe.Message = &spec.Message{Reference: &spec.Reference{Ref: "#/components/messages/" + *details.schema.Name}}
	}
	if details.bindings != nil {
		if details.bindings.operationBinding.Kafka != nil {
			channel.BaseChannelItem.Subscribe.Bindings = &details.bindings.operationBinding
		}
		if details.bindings.channelBindings.Kafka != nil {
			channel.BaseChannelItem.Bindings = &details.bindings.channelBindings
		}
	}
	err := reflector.AddChannel(channel)
	return reflector, err
}

func processChannelDetails(details *AccountDetails) error {
	if err := details.getSchemaDetails(); err != nil {
		if err.Error() == protobufErrorMessage {
			return err
		}
		return fmt.Errorf("failed to get schema details: %w", err)
	}

	err := processBindings(details)
	if err != nil {
		return err
	}

	//TODO schema registry compatibility
	return nil
}

func processBindings(details *AccountDetails) error {
	customConfigMap := make(map[string]string)
	topicConfigMap := make(map[string]any)

	var err error
	topic := details.channelDetails.currentTopic
	for configKey, configValue := range topic.ConfigEntries {
		switch configKey {
		case "cleanup.policy":
			topicConfigMap[configKey] = strings.Split(*configValue, ",")
		case "max.message.bytes":
			topicConfigMap[configKey], err = strconv.ParseInt(*configValue, 10, 32)
			if err != nil {
				return err
			}
		case "delete.retention.ms", "retention.bytes", "retention.ms":
			topicConfigMap[configKey], err = strconv.ParseInt(*configValue, 10, 64)
			if err != nil {
				return err
			}
		default:
			customConfigMap[configKey] = *configValue
		}
	}

	topicConfigs := topicConfigurationExport{}
	jsonString, err := json.Marshal(topicConfigMap)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(jsonString, &topicConfigs); err != nil {
		return err
	}

	var channelBindings any = ChannelBinding{
		BindingVersion:     "0.4.0",
		Partitions:         details.channelDetails.currentTopic.NumPartitions,
		Replicas:           details.channelDetails.currentTopic.ReplicationFactor,
		TopicConfiguration: topicConfigs,
		XConfigs:           customConfigMap,
	}
	messageBindings := spec.MessageBindingsObject{
		Kafka: &spec.KafkaMessage{
			Key: &spec.KafkaMessageKey{Schema: map[string]any{"type": "string"}},
		}}
	operationBindings := spec.OperationBindingsObject{Kafka: &spec.KafkaOperation{
		GroupID:  &spec.KafkaOperationGroupID{Schema: map[string]any{"type": "string"}},
		ClientID: &spec.KafkaOperationClientID{Schema: map[string]any{"type": "string"}},
	}}
	bindings := &bindings{
		messageBinding:   messageBindings,
		operationBinding: operationBindings,
	}
	bindings.channelBindings = spec.ChannelBindingsObject{Kafka: &channelBindings}
	details.channelDetails.bindings = bindings
	return nil
}

func createAsyncReflector(config *K2AConfig) asyncapi.Reflector {
	return asyncapi.Reflector{
		Schema: &spec.AsyncAPI{
			DefaultContentType: "application/json",
			Servers: map[string]spec.ServersAdditionalProperties{
				"cluster": {Server: &spec.Server{
					URL:         config.KafkaUrl,
					Description: "Kafka instance.",
					Protocol:    "kafka",
					Security:    []map[string][]string{{"broker": []string{}}},
				}},
				"schema-registry": {Server: &spec.Server{
					URL:         config.SchemaRegistryUrl,
					Description: "Kafka Schema Registry Server",
					Protocol:    "kafka",
					Security:    []map[string][]string{{"schemaRegistry": []string{}}},
				}},
			},
			Info: spec.Info{
				Version:     config.SpecVersion,
				Title:       "Async API Specification Document for Applications",
				Description: "Async API Specification Document for Applications",
				Contact: &spec.Contact{
					Name: "API Support",
				},
				License: &spec.License{
					Name: "Apache 2.0",
					URL:  "https://www.apache.org/licenses/LICENSE-2.0.html",
				},
			},
		},
	}
}

func CreateAccountDetails(config *K2AConfig) *AccountDetails {
	details := new(AccountDetails)
	details.kafkaUrl = config.KafkaUrl
	details.schemaRegistryUrl = config.SchemaRegistryUrl
	details.config = config
	return details
}

func GetAccountDetails(config *K2AConfig) (*AccountDetails, error) {
	details := CreateAccountDetails(config)

	topics, err := details.queryTopicInfo(config.GetTopics())
	if err != nil {
		return nil, err
	}
	details.topics = topics

	subjects, err := details.querySchemaRegistrySubjects()
	if err != nil {
		return nil, err
	}

	details.subjects = subjects
	return details, nil
}
