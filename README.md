# TiDB Cluster Topology Plot Generator

This tool generates a topology plot based on a TiDB cluster topology file in YAML format. The generated plot is saved as an SVG file.

## How to Run the Tool

1. Ensure you have Go installed on your system. You can download and install Go from [https://golang.org/dl/](https://golang.org/dl/).

2. Clone this repository and navigate to the project directory:

   ```sh
   git clone https://github.com/githubnext/workspace-blank.git
   cd workspace-blank
   ```

3. Install the required Go packages:

   ```sh
   go get gopkg.in/yaml.v2
   go get gonum.org/v1/gonum/graph
   go get gonum.org/v1/plot
   ```

4. Run the tool with the provided YAML file:

   ```sh
   go run main.go
   ```

## Example

An example YAML file (`poc-cluster-config.yaml.txt`) is provided in the repository. The tool will generate a topology plot based on this file and save it as `topology.svg`.

## Description

The TiDB Cluster Topology Plot Generator is a tool designed to visualize the topology of a TiDB cluster. It parses a YAML file containing the cluster configuration and generates a graph representation of the cluster using the `gonum/graph` library. The graph is then visualized and saved as an SVG file using the `gonum/plot` library.

The generated plot provides a hierarchical layout of the cluster nodes, with different shapes and colors representing different node types (e.g., TiDB, TiKV, PD, Prometheus, Grafana, Alertmanager). The plot helps in understanding the structure and dependencies of the TiDB cluster.
