package ingest

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/konveyor/migration-harness/internal/goose"
	"github.com/konveyor/migration-harness/internal/logging"
)

type IngestReport struct {
	RoundsCompleted int    `json:"rounds_completed"`
	Termination     string `json:"termination_reason"`
	ContextFile     string `json:"context_file"`
}

func Run(ctx context.Context, repoDir, runDir, request string, recipesDir string, runner goose.Runner) (*IngestReport, error) {
	logging.Header("Step 1b: Ingest (Interview Loop)")

	logging.Info("1b-a. scanning graphify output for pre-fill...")
	scan, err := ScanGraphifyOutput(repoDir, runDir)
	if err != nil {
		return nil, fmt.Errorf("scan graphify: %w", err)
	}

	rubric := DefaultRubric()
	PrefillCategories(rubric, scan)

	logging.Ok("1b-a. pre-filled %d categories from graph", countPrefilled(rubric))
	logging.Info("Coverage after scan:\n%s", CoverageReport(rubric))

	if RequiredCoverage(rubric) {
		logging.Ok("1b-a. all required categories resolved from graph — skipping interview")
		return writeResult(repoDir, runDir, rubric, scan, request, 0, "coverage_from_scan")
	}

	recipePath := filepath.Join(recipesDir, "ingest.yaml")
	if _, err := os.Stat(recipePath); os.IsNotExist(err) {
		logging.Warn("ingest recipe not found at %s — generating context from scan only", recipePath)
		return writeResult(repoDir, runDir, rubric, scan, request, 0, "no_recipe")
	}

	var rounds []QARound
	termination := "coverage_threshold"

	for round := 1; round <= MaxRounds; round++ {
		logging.Info("1b-b. interview round %d/%d...", round, MaxRounds)

		gaps := LowestScoring(rubric, MaxQuestionsPerRnd)
		if len(gaps) == 0 {
			termination = "coverage_threshold"
			break
		}

		prompt := buildInterviewPrompt(round, gaps, scan, request, rubric)

		result, err := runner.RunRecipe(ctx, recipePath, 5, map[string]string{
			"round":    fmt.Sprintf("%d", round),
			"prompt":   prompt,
			"repo_dir": repoDir,
		})
		if err != nil {
			logging.Warn("interview round %d failed: %v", round, err)
			termination = "runner_error"
			break
		}

		qar := parseInterviewResponse(round, result, gaps)
		rounds = append(rounds, qar)

		updateScoresFromResponse(rubric, qar)

		logging.Info("Coverage after round %d:\n%s", round, CoverageReport(rubric))

		if RequiredCoverage(rubric) {
			logging.Info("1b-c. running free-form round...")
			freeFormPrompt := buildFreeFormPrompt(scan, request)
			freeResult, err := runner.RunRecipe(ctx, recipePath, 3, map[string]string{
				"round":    fmt.Sprintf("%d", round+1),
				"prompt":   freeFormPrompt,
				"repo_dir": repoDir,
			})
			if err == nil {
				ffRound := parseFreeFormResponse(round+1, freeResult)
				rounds = append(rounds, ffRound)
			}
			termination = "coverage_threshold"
			break
		}

		if round == MaxRounds {
			termination = "max_rounds"
		}
	}

	logging.Ok("1b. interview complete (%s after %d rounds)", termination, len(rounds))

	report, err := writeResult(repoDir, runDir, rubric, scan, request, len(rounds), termination)
	if err != nil {
		return nil, err
	}
	report.RoundsCompleted = len(rounds)
	return report, nil
}

func countPrefilled(cats []Category) int {
	n := 0
	for _, c := range cats {
		if c.Score > 0 {
			n++
		}
	}
	return n
}

func buildInterviewPrompt(round int, gaps []Category, scan *ScanResult, request string, rubric []Category) string {
	var b strings.Builder
	fmt.Fprintf(&b, "You are conducting interview round %d for a migration project.\n", round)
	fmt.Fprintf(&b, "Migration request: %s\n\n", request)

	b.WriteString("What we know so far from scanning the codebase:\n")
	fmt.Fprintf(&b, "- Source stack: %s\n", scan.SourceStack)
	if len(scan.CustomDeps) > 0 {
		fmt.Fprintf(&b, "- Custom dependencies: %s\n", strings.Join(scan.CustomDeps, ", "))
	}
	if len(scan.EntryPoints) > 0 {
		fmt.Fprintf(&b, "- Entry points: %s\n", strings.Join(scan.EntryPoints, "; "))
	}
	if len(scan.GodNodes) > 0 {
		fmt.Fprintf(&b, "- High-risk nodes (20+ edges): %s\n", strings.Join(scan.GodNodes, ", "))
	}

	if len(scan.TopNodes) > 0 {
		b.WriteString("\nTop nodes from the code graph (use these to identify messaging, persistence, and framework patterns):\n")
		for _, n := range scan.TopNodes {
			fmt.Fprintf(&b, "  - %s\n", n)
		}
	}

	if scan.GraphReport != "" {
		reportSnippet := scan.GraphReport
		if len(reportSnippet) > 3000 {
			reportSnippet = reportSnippet[:3000] + "\n... (truncated)"
		}
		b.WriteString("\nGraph report summary:\n")
		b.WriteString(reportSnippet)
		b.WriteString("\n")
	}

	b.WriteString("\nCategories needing answers (ask 1 question per category, max 4 total):\n")
	for _, g := range gaps {
		fmt.Fprintf(&b, "- [%s] %s (current score: %d%%)\n", g.ID, g.Name, g.Score)
		if g.Answers != "" {
			fmt.Fprintf(&b, "  Known so far: %s\n", g.Answers)
		}
	}

	b.WriteString("\nRules:\n")
	b.WriteString("- If graphify already partially answered a category, frame as confirmation\n")
	b.WriteString("- For unknown custom dependencies, ask what they do and whether to replace or remove\n")
	b.WriteString("- Present questions as a numbered batch\n")
	b.WriteString("- Return your questions as JSON with format: {\"questions\": [{\"category_id\": \"R2\", \"question\": \"...\"}]}\n")

	return b.String()
}

func buildFreeFormPrompt(scan *ScanResult, request string) string {
	var b strings.Builder
	b.WriteString("All required categories are now resolved. This is the final free-form round.\n\n")
	fmt.Fprintf(&b, "Migration request: %s\n\n", request)
	b.WriteString("Ask the user ONE open-ended question:\n")
	b.WriteString("\"Before I start planning, is there anything else I should know?\n")
	b.WriteString(" For example: custom build steps, deployment scripts, known gotchas,\n")
	b.WriteString(" or things that surprised previous developers.\"\n\n")
	b.WriteString("Return as JSON: {\"questions\": [{\"category_id\": \"FREEFORM\", \"question\": \"...\"}]}\n")
	return b.String()
}

func parseInterviewResponse(round int, raw json.RawMessage, gaps []Category) QARound {
	qar := QARound{Round: round}

	var resp struct {
		Questions []struct {
			CategoryID string `json:"category_id"`
			Question   string `json:"question"`
		} `json:"questions"`
		Answers []struct {
			CategoryID string `json:"category_id"`
			Answer     string `json:"answer"`
		} `json:"answers"`
	}

	if err := json.Unmarshal(raw, &resp); err != nil {
		for _, g := range gaps {
			qar.Questions = append(qar.Questions, fmt.Sprintf("[%s] %s", g.ID, g.Name))
		}
		return qar
	}

	for _, q := range resp.Questions {
		qar.Questions = append(qar.Questions, fmt.Sprintf("[%s] %s", q.CategoryID, q.Question))
	}
	for _, a := range resp.Answers {
		qar.Answers = append(qar.Answers, fmt.Sprintf("[%s] %s", a.CategoryID, a.Answer))
	}

	return qar
}

func parseFreeFormResponse(round int, raw json.RawMessage) QARound {
	qar := QARound{Round: round}
	qar.Questions = append(qar.Questions, "[FREEFORM] Is there anything else I should know?")
	return qar
}

func updateScoresFromResponse(rubric []Category, qar QARound) {
	for _, ans := range qar.Answers {
		for i := range rubric {
			catPrefix := "[" + rubric[i].ID + "]"
			if strings.HasPrefix(ans, catPrefix) {
				answer := strings.TrimPrefix(ans, catPrefix+" ")
				rubric[i].Answers += "; " + answer
				rubric[i].Score += 40
				if rubric[i].Score > 100 {
					rubric[i].Score = 100
				}
			}
		}
	}
}

func writeResult(repoDir, runDir string, rubric []Category, scan *ScanResult, request string, rounds int, termination string) (*IngestReport, error) {
	contextPath := filepath.Join(repoDir, "INGEST_CONTEXT.md")

	var b strings.Builder
	b.WriteString("# Migration Ingest Context\n\n")
	fmt.Fprintf(&b, "| Field | Value |\n")
	fmt.Fprintf(&b, "|-------|-------|\n")
	fmt.Fprintf(&b, "| Generated | %s |\n", time.Now().Format("2006-01-02 15:04"))
	fmt.Fprintf(&b, "| Source application | %s |\n", repoDir)
	fmt.Fprintf(&b, "| Migration request | %s |\n", request)
	fmt.Fprintf(&b, "| Rounds completed | %d of %d max |\n", rounds, MaxRounds)
	fmt.Fprintf(&b, "| Termination reason | %s |\n\n", termination)

	b.WriteString("---\n\n")

	for _, c := range rubric {
		if !c.Required {
			continue
		}
		fmt.Fprintf(&b, "## %s (%s)\n", c.Name, c.ID)
		fmt.Fprintf(&b, "**Score**: %d%%\n\n", c.Score)
		if c.Answers != "" {
			fmt.Fprintf(&b, "%s\n\n", c.Answers)
		} else {
			b.WriteString("_Not addressed_\n\n")
		}
	}

	b.WriteString("## Assumptions\n\n")
	b.WriteString("| Assumption | Confidence | Category |\n")
	b.WriteString("|-----------|-----------|----------|\n")
	for _, c := range rubric {
		if c.Required && c.Score < 100 && c.Score >= CoverageThreshold {
			confidence := "MEDIUM"
			if c.Score >= 90 {
				confidence = "HIGH"
			}
			fmt.Fprintf(&b, "| %s details partially assumed | %s | %s |\n", c.Name, confidence, c.ID)
		} else if c.Required && c.Score < CoverageThreshold {
			fmt.Fprintf(&b, "| %s not fully resolved | LOW | %s |\n", c.Name, c.ID)
		}
	}
	b.WriteString("\n")

	if len(scan.MigrationOrder) > 0 {
		b.WriteString("## Migration Order (from graph topology)\n\n")
		for i, f := range scan.MigrationOrder {
			fmt.Fprintf(&b, "%d. %s\n", i+1, f)
			if i >= 29 {
				fmt.Fprintf(&b, "... and %d more files\n", len(scan.MigrationOrder)-30)
				break
			}
		}
		b.WriteString("\n")
	}

	if err := os.WriteFile(contextPath, []byte(b.String()), 0644); err != nil {
		return nil, fmt.Errorf("write INGEST_CONTEXT.md: %w", err)
	}

	copyFile(contextPath, filepath.Join(runDir, "INGEST_CONTEXT.md"))

	resultJSON := &IngestResult{
		SourceStack:     scan.SourceStack,
		RoundsCompleted: rounds,
		Termination:     termination,
	}
	if len(scan.MigrationOrder) > 0 {
		resultJSON.MigrationOrder = scan.MigrationOrder
	}
	if len(scan.EntryPoints) > 0 {
		resultJSON.EntryPoint = strings.Join(scan.EntryPoints, "; ")
	}

	jsonData, err := json.MarshalIndent(resultJSON, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal ingest.json: %w", err)
	}
	jsonPath := filepath.Join(runDir, "ingest.json")
	if err := os.WriteFile(jsonPath, jsonData, 0644); err != nil {
		return nil, fmt.Errorf("write ingest.json: %w", err)
	}

	logging.Ok("INGEST_CONTEXT.md written to %s", contextPath)

	return &IngestReport{
		RoundsCompleted: rounds,
		Termination:     termination,
		ContextFile:     contextPath,
	}, nil
}

func copyFile(src, dst string) {
	data, err := os.ReadFile(src)
	if err != nil {
		return
	}
	os.WriteFile(dst, data, 0644)
}
