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