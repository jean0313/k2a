# Cli to Generate AsyncApi Specification from Kafka

## Setup
- Kafka Cluster
- SchemaRegistry Cluster

### Init
```bash
# kafka broker: localhost:9092
# init kafka test topic
bin/kafka-topics.sh --bootstrap-server localhost:9092 \
    --create \
    --topic test-topic \
    --partitions 5 \
    --replication-factor 1 \
    --config retention.ms=720000 \
    --config delete.retention.ms=360000 \
    --config max.message.bytes=64000 \
    --config flush.messages=1

# schema registry cluster: http://localhost:8081
# init registry schema
curl localhost:8081/subjects/test-topic-value/versions \
    -H 'Content-Type: application/json' \
    -d '{ "schema": "{ \"type\": \"record\", \"name\": \"test-topic\", \"fields\": [ { \"type\": \"string\", \"name\": \"field1\" }, { \"type\": \"int\", \"name\": \"field2\", \"default\": 0}] }", "schemaType": "AVRO"}'    
```

## Build

### Pre-Requisites
go version go1.21.5


```bash
# for macos x64
make macos

# for linux x64
make linux

# for windows x64
make windows
```

## Run

### cli k2a help
```
Export an AsyncAPI specification for a Kafka cluster and Schema Registry.

Usage:
  cli k2a [flags]

Examples:
cli k2a --kurl prod.kafka.com:9092 --rurl http://prod.schema-registry.com --topics demo,sample

Flags:
      --ca-file string        The optional certificate authority file for TLS client authentication
      --cert string           The optional certificate file for client authentication
      --file string           Output file name (default "k2a.yaml")
  -h, --help                  help for k2a
      --key-file string       The optional key file for client authentication
      --kurl string           Kafka cluster broker url (default "localhost:9092")
      --rurl string           Schema registry url (default "http://localhost:8081")
      --spec-version string   Version number of the output file. (default 1.0.0)
      --tls-skip-verify       Whether to skip TLS server cert verification (default true)
      --topics string         Topics to export
      --use-tls               Use TLS to communicate with the kafka cluster

```

### cli generate example
```bash
cli k2a --topics test-topic

```

### cli ws help
```
Start web server to export topics, browser to localhost:8080

Usage:
  cli ws [flags]

Flags:
  -h, --help          help for ws
      --port string   server port to listen (default "8080")
```

![ws](docs/ws.jpg)


## Verify
[![Studio](docs/studio.jpg)](https://studio.asyncapi.com/)

### Output sample
```yaml
asyncapi: 2.4.0
info:
  title: Async API Specification Document for Applications
  version: 1.0.0
  description: Async API Specification Document for Applications
  contact:
    name: API Support
  license:
    name: Apache 2.0
    url: https://www.apache.org/licenses/LICENSE-2.0.html
servers:
  cluster:
    url: localhost:9092
    description: Kafka instance.
    protocol: kafka
    security:
    - broker: []
  schema-registry:
    url: http://localhost:8081
    description: Kafka Schema Registry Server
    protocol: kafka
    security:
    - schemaRegistry: []
defaultContentType: application/json
channels:
  demo:
    subscribe:
      operationId: DemoSubscribe
      bindings:
        kafka:
          bindingVersion: 0.3.0
          groupId:
            type: string
          clientId:
            type: string
      message:
        $ref: '#/components/messages/DemoMessage'
    bindings:
      kafka:
        bindingVersion: 0.4.0
        partitions: 1
        replicas: 1
        topicConfiguration: {}
  test-rep:
    subscribe:
      operationId: TestRepSubscribe
      bindings:
        kafka:
          bindingVersion: 0.3.0
          groupId:
            type: string
          clientId:
            type: string
      message:
        $ref: '#/components/messages/TestRepMessage'
    bindings:
      kafka:
        bindingVersion: 0.4.0
        partitions: 5
        replicas: 1
        topicConfiguration:
          retention.ms: 1.2345e+07
          delete.retention.ms: 3.6e+06
          max.message.bytes: 64000
        x-configs:
          flush.messages: "1"
  test-topic:
    subscribe:
      operationId: TestTopicSubscribe
      bindings:
        kafka:
          bindingVersion: 0.3.0
          groupId:
            type: string
          clientId:
            type: string
      message:
        $ref: '#/components/messages/TestTopicMessage'
    bindings:
      kafka:
        bindingVersion: 0.4.0
        partitions: 1
        replicas: 1
        topicConfiguration:
          max.message.bytes: 64000
        x-configs:
          flush.messages: "1"
components:
  messages:
    DemoMessage:
      schemaFormat: application/schema+json;version=draft-07
      contentType: application/json
      messageId: DemoMessage
      payload:
        properties:
          name:
            type: string
        type: object
      name: DemoMessage
      bindings:
        kafka:
          bindingVersion: 0.3.0
          key:
            type: string
    TestRepMessage:
      schemaFormat: application/vnd.apache.avro;version=1.9.0
      contentType: application/avro
      messageId: TestRepMessage
      payload:
        fields:
        - name: name
          type: string
        - default: 20000
          name: money
          type: int
        name: account
        type: record
      name: TestRepMessage
      bindings:
        kafka:
          bindingVersion: 0.3.0
          key:
            type: string
    TestTopicMessage:
      schemaFormat: application/vnd.apache.avro;version=1.9.0
      contentType: application/avro
      messageId: TestTopicMessage
      payload:
        fields:
        - name: name
          type: string
        - default: 20
          name: age
          type: int
        name: test
        type: record
      name: TestTopicMessage
      bindings:
        kafka:
          bindingVersion: 0.3.0
          key:
            type: string
  securitySchemes:
    broker:
      type: userPassword
      x-configs:
        sasl.mechanisms: PLAIN
        sasl.password: '{{CLUSTER_API_SECRET}}'
        sasl.username: '{{CLUSTER_API_KEY}}'
        security.protocol: sasl_ssl
    schemaRegistry:
      type: userPassword
      x-configs:
        basic.auth.user.info: '{{SCHEMA_REGISTRY_API_KEY}}:{{SCHEMA_REGISTRY_API_SECRET}}'
```