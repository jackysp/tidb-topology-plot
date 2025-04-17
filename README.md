# TiDB Cluster Topology Plot Generator

This tool generates a topology diagram for a TiDB cluster based on a YAML configuration file. It produces a Graphviz DOT file and renders an SVG using the `dot` command.

## Prerequisites

- Go (1.18+)
- Graphviz (for the `dot` command)

On macOS, you can install Graphviz with Homebrew:

```sh
brew install graphviz
```

## Building

In the project root, fetch dependencies and build:

```sh
go mod tidy
go build -o topology-plot main.go
```

## Running

Provide your YAML file as the first argument. For example:

```sh
./topology-plot poc-cluster-config.yaml
```

This will generate:

- `topology.dot` — the Graphviz DOT representation
- `topology.svg` — the rendered SVG diagram

## Example

```sh
go run main.go poc-cluster-config.yaml
```

After running, open `topology.svg` in your browser or image viewer to see the cluster layout.

## Configuration Format

See `poc-cluster-config.yaml` for an example of the YAML schema, which follows the TiUP specification for TiDB cluster topology.
