package ingest

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/konveyor/migration-harness/internal/detect"
)

type ScanResult struct {
	SourceStack    string
	Manifests      detect.Manifests
	EntryPoints    []string
	MigrationOrder []string
	CustomDeps     []string
	GodNodes       []string
	TopNodes       []string
	GraphReport    string
}

func ScanGraphifyOutput(repoDir, runDir string) (*ScanResult, error) {
	result := &ScanResult{}

	detectPath := filepath.Join(runDir, "detect.json")
	data, err := os.ReadFile(detectPath)
	if err != nil {
		return nil, fmt.Errorf("read detect.json: %w", err)
	}

	var dr detect.DetectResult
	if err := json.Unmarshal(data, &dr); err != nil {
		return nil, fmt.Errorf("parse detect.json: %w", err)
	}

	result.Manifests = dr.Manifests
	result.SourceStack = inferSourceStack(dr)

	reportPath := filepath.Join(runDir, "GRAPH_REPORT.md")
	if reportData, err := os.ReadFile(reportPath); err == nil {
		result.GraphReport = string(reportData)
	}

	graphPath := filepath.Join(runDir, "graph.json")
	graphData, err := os.ReadFile(graphPath)
	if err != nil {
		return result, nil
	}

	var g detect.GraphJSON
	if err := json.Unmarshal(graphData, &g); err != nil {
		return result, nil
	}

	result.EntryPoints = findEntryPoints(g)
	result.GodNodes = findGodNodes(g)
	result.TopNodes = collectTopNodes(g, 30)
	result.MigrationOrder = computeMigrationOrder(g)
	result.CustomDeps = findCustomDeps(repoDir)

	return result, nil
}

func findEntryPoints(g detect.GraphJSON) []string {
	var entries []string
	for _, node := range g.Nodes {
		if node.Degree > 0 && node.SourceFile != "" {
			label := strings.ToLower(node.Label)
			if strings.Contains(label, "main") ||
				strings.Contains(label, "startup") ||
				strings.Contains(label, "program") ||
				strings.Contains(label, "application") ||
				strings.Contains(label, "app") {
				entries = append(entries, node.Label+" -> "+node.SourceFile)
			}
		}
	}
	return entries
}

func findGodNodes(g detect.GraphJSON) []string {
	var gods []string
	for _, node := range g.Nodes {
		if node.Degree > 20 {
			gods = append(gods, node.Label)
		}
	}
	return gods
}

func collectTopNodes(g detect.GraphJSON, limit int) []string {
	type nodeInfo struct {
		label  string
		file   string
		degree int
	}

	var nodes []nodeInfo
	for _, n := range g.Nodes {
		if n.SourceFile != "" {
			nodes = append(nodes, nodeInfo{n.Label, n.SourceFile, n.Degree})
		}
	}

	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].degree > nodes[j].degree
	})

	var top []string
	for i, n := range nodes {
		if i >= limit {
			break
		}
		top = append(top, fmt.Sprintf("%s (%d edges) -> %s", n.label, n.degree, n.file))
	}
	return top
}

func inferSourceStack(dr detect.DetectResult) string {
	var parts []string
	if dr.Files.Java > 0 {
		parts = append(parts, fmt.Sprintf("Java (%d files)", dr.Files.Java))
	}
	if dr.Files.Python > 0 {
		parts = append(parts, fmt.Sprintf("Python (%d files)", dr.Files.Python))
	}
	if dr.Files.TypeScript > 0 {
		parts = append(parts, fmt.Sprintf("TypeScript (%d files)", dr.Files.TypeScript))
	}
	if dr.Files.JavaScript > 0 {
		parts = append(parts, fmt.Sprintf("JavaScript (%d files)", dr.Files.JavaScript))
	}
	if dr.Files.Go > 0 {
		parts = append(parts, fmt.Sprintf("Go (%d files)", dr.Files.Go))
	}
	if dr.Files.CSharp > 0 {
		parts = append(parts, fmt.Sprintf("C# (%d files)", dr.Files.CSharp))
	}
	if dr.Files.Ruby > 0 {
		parts = append(parts, fmt.Sprintf("Ruby (%d files)", dr.Files.Ruby))
	}
	if dr.Files.Rust > 0 {
		parts = append(parts, fmt.Sprintf("Rust (%d files)", dr.Files.Rust))
	}

	if dr.Manifests.PomXML {
		parts = append(parts, "Maven")
	}
	if dr.Manifests.PackageJSON {
		parts = append(parts, "npm")
	}
	if dr.Manifests.GoMod {
		parts = append(parts, "Go modules")
	}
	if dr.Manifests.PyprojectTOML || dr.Manifests.RequirementsTXT || dr.Manifests.SetupPy {
		parts = append(parts, "Python packaging")
	}
	if dr.Manifests.CargoTOML {
		parts = append(parts, "Cargo")
	}
	if dr.Manifests.Gemfile {
		parts = append(parts, "Bundler")
	}

	if len(parts) == 0 {
		return "unknown"
	}
	return strings.Join(parts, ", ")
}

func computeMigrationOrder(g detect.GraphJSON) []string {
	type nodeScore struct {
		label  string
		file   string
		degree int
	}

	var nodes []nodeScore
	for _, n := range g.Nodes {
		if n.SourceFile != "" {
			nodes = append(nodes, nodeScore{n.Label, n.SourceFile, n.Degree})
		}
	}

	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].degree < nodes[j].degree
	})

	seen := make(map[string]bool)
	var order []string
	for _, n := range nodes {
		if !seen[n.file] {
			seen[n.file] = true
			order = append(order, n.file)
		}
	}
	return order
}

func findCustomDeps(repoDir string) []string {
	var deps []string

	libDirs := []string{"lib", "vendor", "libs", "third_party", "packages"}
	depExts := []string{".jar", ".dll", ".so", ".dylib", ".whl", ".egg"}

	for _, dir := range libDirs {
		dirPath := filepath.Join(repoDir, dir)
		entries, err := os.ReadDir(dirPath)
		if err != nil {
			continue
		}
		for _, e := range entries {
			name := e.Name()
			for _, ext := range depExts {
				if strings.HasSuffix(strings.ToLower(name), ext) {
					deps = append(deps, dir+"/"+name)
					break
				}
			}
		}
	}

	manifestChecks := map[string][]string{
		"pom.xml":    {"<scope>system</scope>", "<systemPath>"},
		"*.csproj":   {"<Reference Include=", "<HintPath>"},
		"setup.py":   {"sys.path", "ext_modules"},
		"Cargo.toml": {"path = \""},
	}

	for manifest, patterns := range manifestChecks {
		var paths []string
		if strings.Contains(manifest, "*") {
			matches, _ := filepath.Glob(filepath.Join(repoDir, manifest))
			paths = matches
		} else {
			paths = []string{filepath.Join(repoDir, manifest)}
		}
		for _, p := range paths {
			data, err := os.ReadFile(p)
			if err != nil {
				continue
			}
			content := string(data)
			for _, pattern := range patterns {
				if strings.Contains(content, pattern) {
					deps = append(deps, fmt.Sprintf("local/vendored dependency in %s", filepath.Base(p)))
					break
				}
			}
		}
	}

	return deps
}

func PrefillCategories(rubric []Category, scan *ScanResult) {
	for i := range rubric {
		switch rubric[i].ID {
		case "R1":
			if scan.SourceStack != "unknown" {
				rubric[i].Score = 70
				rubric[i].Answers = "Pre-filled from graphify: " + scan.SourceStack
			}
		case "R3":
			if len(scan.CustomDeps) > 0 {
				rubric[i].Score = 30
				rubric[i].Answers = "Detected: " + strings.Join(scan.CustomDeps, ", ")
			} else {
				rubric[i].Score = 80
				rubric[i].Answers = "No custom dependencies detected"
			}
		case "R4":
			rubric[i].Score = 0
			rubric[i].Answers = "LLM will analyze graph for messaging/async patterns"
		case "R5":
			rubric[i].Score = 0
			rubric[i].Answers = "LLM will analyze graph for persistence patterns"
		case "R6":
			if len(scan.EntryPoints) > 0 {
				rubric[i].Score = 90
				rubric[i].Answers = "Entry points: " + strings.Join(scan.EntryPoints, "; ")
			}
			if len(scan.MigrationOrder) > 0 {
				rubric[i].Score = 100
			}
		}
	}
}
