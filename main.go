package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/simple"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"github.com/pingcap/tiup/pkg/cluster/spec"
)

func main() {
	// Read the YAML file
	yamlFile, err := ioutil.ReadFile("poc-cluster-config.yaml.txt")
	if err != nil {
		log.Fatalf("Error reading YAML file: %v", err)
	}

	// Parse the YAML file
	var cluster spec.Specification
	err = yaml.Unmarshal(yamlFile, &cluster)
	if err != nil {
		log.Fatalf("Error parsing YAML file: %v", err)
	}

	// Create a new graph
	g := simple.NewDirectedGraph()

	// Add nodes to the graph
	addNodesToGraph(g, cluster)

	// Create a plot
	p, err := plot.New()
	if err != nil {
		log.Fatalf("Error creating plot: %v", err)
	}

	// Add the graph to the plot
	err = addGraphToPlot(p, g)
	if err != nil {
		log.Fatalf("Error adding graph to plot: %v", err)
	}

	// Save the plot as an SVG file
	err = p.Save(10*vg.Inch, 10*vg.Inch, "topology.svg")
	if err != nil {
		log.Fatalf("Error saving plot: %v", err)
	}

	fmt.Println("Topology plot saved as topology.svg")
}

func addNodesToGraph(g *simple.DirectedGraph, cluster spec.Specification) {
	// Add PD nodes
	for _, pd := range cluster.PDServers {
		node := g.NewNode()
		g.AddNode(node)
	}

	// Add TiKV nodes
	for _, tikv := range cluster.TiKVServers {
		node := g.NewNode()
		g.AddNode(node)
	}

	// Add TiDB nodes
	for _, tidb := range cluster.TiDBServers {
		node := g.NewNode()
		g.AddNode(node)
	}

	// Add Prometheus nodes
	for _, prometheus := range cluster.PrometheusServers {
		node := g.NewNode()
		g.AddNode(node)
	}

	// Add Grafana nodes
	for _, grafana := range cluster.GrafanaServers {
		node := g.NewNode()
		g.AddNode(node)
	}

	// Add Alertmanager nodes
	for _, alertmanager := range cluster.AlertmanagerServers {
		node := g.NewNode()
		g.AddNode(node)
	}
}

func addGraphToPlot(p *plot.Plot, g *simple.DirectedGraph) error {
	// Create a new graph plotter
	graphPlotter, err := plotter.NewGraph(g)
	if err != nil {
		return err
	}

	// Add the graph plotter to the plot
	p.Add(graphPlotter)

	return nil
}
