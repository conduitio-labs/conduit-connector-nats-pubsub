# Conduit Connector NATS PubSub

### General

The [NATS](https://nats.io/) PubSub connector is one of [Conduit](https://github.com/ConduitIO/conduit) plugins. It provides both, a source and a destination NATS PubSub connector.

### Prerequisites

- [Go](https://go.dev/) 1.18
- (optional) [golangci-lint](https://github.com/golangci/golangci-lint) 1.45.2
- [NATS](https://nats.io/download/) 2.8.4 and [nats.go](https://github.com/nats-io/nats.go) library v1.16.0

### How to build it

Run `make`.

### Testing

Run `make test` to run all the unit and integration tests, which require Docker and Docker Compose to be installed and running. The command will handle starting and stopping docker containers for you.

## Source

### Connection and authentication

The NATS PubSub connector connects to a NATS server or a cluster with the required parameters `urls`, `subject` and `mode`. If your NATS server has a configured authentication you can pass an authentication details in the connection URL. For example, for a token authentication the url will look like: `nats://mytoken@127.0.0.1:4222`, and for a username/password authentication: `nats://username:password@127.0.0.1:4222`. But if your server is using [NKey](https://docs.nats.io/using-nats/developer/connecting/nkey) or [Credentials file](https://docs.nats.io/using-nats/developer/connecting/creds) for authentication you must configure them via seperate [configuration](#configuration) parameters.

### Receiving messages

The connector listening on a subject receives messages published on that subject. If the connector is not actively listening on the subject, the message is not received. The connector can use the [wildcard](https://docs.nats.io/nats-concepts/subjects#wildcards) tokens such as `*` and `>` to match a single token or to match the tail of a subject.

### Position handling

The position is a random binary marshaled UUIDv4.

### Configuration

The config passed to Configure can contain the following fields.

| name                      | description                                                                                                                                                                                                                                       | required | default                            |
| ------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | -------- | ---------------------------------- |
| `urls`                    | A list of connection URLs joined by comma. Must be a valid URLs.<br />Examples:<br />`nats://127.0.0.1:1222`<br />`nats://127.0.0.1:1222,nats://127.0.0.1:1223`<br />`nats://myname:password@127.0.0.1:4222`<br />`nats://mytoken@127.0.0.1:4222` | **true** |                                    |
| `subject`                 | A name of a subject from which or to which the connector should read write.                                                                                                                                                                       | **true** |                                    |
| `connectionName`          | Optional connection name which will come in handy when it comes to monitoring                                                                                                                                                                     | false    | `conduit-connection-<random_uuid>` |
| `nkeyPath`                | A path pointed to a [NKey](https://docs.nats.io/using-nats/developer/connecting/nkey) pair. Must be a valid file path. Requird if your NATS server is using NKey authentication.                                                                  | false    |                                    |
| `credentialsFilePath`     | A path pointed to a [credentials file](https://docs.nats.io/using-nats/developer/connecting/creds). Must be a valid file path                                                                                                                     | false    |                                    |
| `tlsClientCertPath`       | A path pointed to a TLS client certificate, must be present if tlsClientPrivateKeyPath field is also present. Must be a valid file path. Required if your NATS server is using TLS.                                                               | false    |                                    |
| `tlsClientPrivateKeyPath` | A path pointed to a TLS client private key, must be present if tlsClientCertPath field is also present. Must be a valid file path. Required if your NATS server is using TLS.                                                                     | false    |                                    |
| `tlsRootCACertPath`       | A path pointed to a TLS root certificate, provide if you want to verify server’s identity. Must be a valid file path                                                                                                                              | false    |                                    |
| `bufferSize`              | A buffer size for consumed messages. It must be set to avoid the [slow consumers](https://docs.nats.io/running-a-nats-service/nats_admin/slow_consumers) problem. Minimum allowed value is `64`                                                   | false    | `1024`                             |

## Destination

### Sending messages

The connector sends message synchronously, one by one.

### Configuration

The config passed to Configure can contain the following fields.

| name                      | description                                                                                                                                                                                                                                       | required | default                            |
| ------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | -------- | ---------------------------------- |
| `urls`                    | A list of connection URLs joined by comma. Must be a valid URLs.<br />Examples:<br />`nats://127.0.0.1:1222`<br />`nats://127.0.0.1:1222,nats://127.0.0.1:1223`<br />`nats://myname:password@127.0.0.1:4222`<br />`nats://mytoken@127.0.0.1:4222` | **true** |                                    |
| `subject`                 | A name of a subject from which or to which the connector should read write.                                                                                                                                                                       | **true** |                                    |
| `connectionName`          | Optional connection name which will come in handy when it comes to monitoring                                                                                                                                                                     | false    | `conduit-connection-<random_uuid>` |
| `nkeyPath`                | A path pointed to a [NKey](https://docs.nats.io/using-nats/developer/connecting/nkey) pair. Must be a valid file path. Requird if your NATS server is using NKey authentication.                                                                  | false    |                                    |
| `credentialsFilePath`     | A path pointed to a [credentials file](https://docs.nats.io/using-nats/developer/connecting/creds). Must be a valid file path                                                                                                                     | false    |                                    |
| `tlsClientCertPath`       | A path pointed to a TLS client certificate, must be present if tlsClientPrivateKeyPath field is also present. Must be a valid file path. Required if your NATS server is using TLS.                                                               | false    |                                    |
| `tlsClientPrivateKeyPath` | A path pointed to a TLS client private key, must be present if tlsClientCertPath field is also present. Must be a valid file path. Required if your NATS server is using TLS.                                                                     | false    |                                    |
| `tlsRootCACertPath`       | A path pointed to a TLS root certificate, provide if you want to verify server’s identity. Must be a valid file path                                                                                                                              | false    |                                    |
