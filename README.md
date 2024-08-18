# Go gRPC Server Template

## What is this?

This is a starter template to kick off a repo holding a gRPC server written in golang.

It includes

1) A top level main that listens for SIGINT and SIGTERM and gracefully closes its database connection and the gRPC server
2) Has a simple embedded data store [badger](https://dgraph.io/docs/badger)
3) Has a starting point for configuration with a [interface wrapper](./config//config.go) around [viper](https://github.com/spf13/viper)
4) Has a [github workflow](https://github.com/curium-rocks/flows/blob/main/.github/workflows/golang.yml)
5) Has a [dev container](./.devcontainer/devcontainer.json)
6) Has a [Dockerfile](./Dockerfile)
7) Has a [Makefile](./Makefile) for common tasks such as building, testing, linting
8) Is automatically updated with renovate [renovate.json](./renovate.json)

## How do I change the .proto and update the associated code?

You can run `make generate-grpc-code` in the dev container and it will re-generate the golang code under api/v1 to match what's specified in the .proto file.


## What configuration properties are available?

| Configuration Property       | Default Value       | Description                           |
|------------------------------|---------------------|---------------------------------------|
| `database.path`              | `data/db`           | Path to the database file             |
| `server.port`                | `50051`             | Port on which the server listens      |
| `server.address`             | `localhost`         | Address on which the server listens   |
| `server.tls.enabled`         | `false`             | Enable TLS for the server             |
| `server.tls.cert`            | `""`                | TLS certificate content               |
| `server.tls.cert_path`       | `""`                | Path to the TLS certificate file      |
| `server.tls.key`             | `""`                | TLS key content                       |
| `server.tls.key_path`        | `""`                | Path to the TLS key file              |
| `server.tls.ca`              | `""`                | CA certificate content                |
| `server.tls.ca_path`         | `""`                | Path to the CA certificate file       |

### How to set configuration values

You can set these configuration properties using either environment variables or a YAML configuration file. Viper is used to extract the configuration properties.

#### Using Environment Variables

Set the environment variables with the corresponding configuration property names in uppercase and replace dots with underscores. For example:

```sh
export DATABASE_PATH="custom/db/path"
export SERVER_PORT="8080"
export SERVER_ADDRESS="0.0.0.0"
export SERVER_TLS_ENABLED="true"
export SERVER_TLS_CERT="your_cert_content"
export SERVER_TLS_CERT_PATH="/path/to/cert"
export SERVER_TLS_KEY="your_key_content"
export SERVER_TLS_KEY_PATH="/path/to/key"
export SERVER_TLS_CA="your_ca_content"
export SERVER_TLS_CA_PATH="/path/to/ca"
```

#### Using a config file

Create a `config` file in the same directory that the process runs in.

Below is an example

``` yaml
database:
  path: "custom/db/path"

server:
  port: 8080
  address: "0.0.0.0"
  tls:
    enabled: true
    cert: "your_cert_content"
    cert_path: "/path/to/cert"
    key: "your_key_content"
    key_path: "/path/to/key"
    ca: "your_ca_content"
    ca_path: "/path/to/ca"
```

#### Certs/Keys

`server.tls.cert` and the matching fields without the `_path` suffix, are expected to be string values in PEM format.
These will be used first, if they are empty, then the associated path variable is used, the expectation is the file provided in the path property is also PEM formatted.
