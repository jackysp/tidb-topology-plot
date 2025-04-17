package main

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"sort"
	"strings"

	"github.com/pingcap/tiup/pkg/cluster/spec"
	"gopkg.in/yaml.v2"
)

func main() {
	// Read the YAML file
	yamlFile, err := os.ReadFile("poc-cluster-config.yaml")
	if err != nil {
		log.Fatalf("Error reading YAML file: %v", err)
	}

	// Parse the YAML file
	var cluster spec.Specification
	err = yaml.Unmarshal(yamlFile, &cluster)
	if err != nil {
		log.Fatalf("Error parsing YAML file: %v", err)
	}

	// Generate DOT topology file
	if err := generateDot(cluster); err != nil {
		log.Fatalf("Error generating DOT topology: %v", err)
	}
	fmt.Println("DOT topology saved as topology.dot")

	// Render SVG from DOT file
	if err := generateSvg(); err != nil {
		log.Fatalf("Error generating SVG topology: %v", err)
	}
	fmt.Println("SVG topology saved as topology.svg")
}

// generateSvg executes Graphviz dot to produce SVG output
func generateSvg() error {
	cmd := exec.Command("dot", "-Tsvg", "topology.dot", "-o", "topology.svg")
	return cmd.Run()
}

// generateDot creates a Graphviz DOT topology grouping components by host
func generateDot(cluster spec.Specification) error {
	// build component groups by host dynamically
	groups := make(map[string][]string)
	// group each instance by host and include port in key (e.g., "tidb:4000")
	for _, comp := range cluster.ComponentsByStartOrder() {
		for _, inst := range comp.Instances() {
			host := inst.GetHost()
			name := comp.Name()
			port := inst.GetMainPort()
			groups[host] = append(groups[host], fmt.Sprintf("%s:%d", name, port))
		}
	}
	// sort hosts
	hosts := make([]string, 0, len(groups))
	for h := range groups {
		hosts = append(hosts, h)
	}
	sortHosts(hosts)
	// sort component nodes within each host alphabetically
	for _, host := range hosts {
		sort.Strings(groups[host])
	}

	var b strings.Builder
	// shapes for each component type
	shapeMap := map[string]string{
		"pd":           "circle",
		"tikv":         "box",
		"tidb":         "ellipse",
		"prometheus":   "diamond",
		"grafana":      "octagon",
		"alertmanager": "hexagon",
	}

	b.WriteString("digraph topology {\n")
	b.WriteString("  compound=true;\n")
	b.WriteString("  rankdir=TB;\n")
	b.WriteString("  nodesep=0.5;\n")
	b.WriteString("  ranksep=0.5;\n\n")
	// create cluster for each host
	for i, host := range hosts {
		b.WriteString(fmt.Sprintf("  subgraph cluster_%d {\n", i))
		b.WriteString(fmt.Sprintf("    label=\"%s\";\n", host))
		// define component nodes with unique IDs and stack vertically
		for j, c := range groups[host] {
			// c format: "component:port"
			id := fmt.Sprintf("%s_%s_%d", host, strings.ReplaceAll(c, ":", "_"), j)
			parts := strings.SplitN(c, ":", 2)
			compName := parts[0]
			portStr := parts[1]
			label := fmt.Sprintf("%s:%s", strings.Title(compName), portStr)
			// lookup shape by component name
			shape, ok := shapeMap[compName]
			if !ok {
				shape = "box"
			}
			b.WriteString(fmt.Sprintf("    \"%s\" [label=\"%s\" shape=\"%s\"];\n", id, label, shape))
			// stack nodes vertically with invisible edges
			if j > 0 {
				prevId := fmt.Sprintf("%s_%s_%d", host, strings.ReplaceAll(groups[host][j-1], ":", "_"), j-1)
				b.WriteString(fmt.Sprintf("    \"%s\" -> \"%s\" [style=invis, constraint=true];\n", prevId, id))
			}
		}
		b.WriteString("  }\n\n")
	}

	b.WriteString("\n")

	b.WriteString("}")
	return os.WriteFile("topology.dot", []byte(b.String()), 0644)
}

// sortHosts sorts host strings by numeric IP order if possible, fallback to lexical
func sortHosts(hosts []string) {
	sort.Slice(hosts, func(i, j int) bool {
		h1, _, e1 := net.SplitHostPort(hosts[i])
		if e1 != nil {
			h1 = hosts[i]
		}
		h2, _, e2 := net.SplitHostPort(hosts[j])
		if e2 != nil {
			h2 = hosts[j]
		}
		ip1, ip2 := net.ParseIP(h1), net.ParseIP(h2)
		if ip1 != nil && ip2 != nil {
			if cmp := bytes.Compare(ip1, ip2); cmp != 0 {
				return cmp < 0
			}
		}
		return hosts[i] < hosts[j]
	})
}
