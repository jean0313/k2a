# Cli for Generate AsyncApi Specification from Kafka

## Setup
- Kafka Cluster
- SchemaRegistry Cluster

## Build
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
      --spec-version string   Version number of the output file. (default 1.0.0) (default "1.0.0")
      --tls-skip-verify       Whether to skip TLS server cert verification (default true)
      --topics string         Topics to export
      --use-tls               Use TLS to communicate with the kafka cluster

```