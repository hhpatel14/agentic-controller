package ingest

import (
	"testing"
)

func TestDefaultRubric(t *testing.T) {
	rubric := DefaultRubric()
	if len(rubric) != 12 {
		t.Errorf("expected 12 categories, got %d", len(rubric))
	}

	requiredCount := 0
	for _, c := range rubric {
		if c.Required {
			requiredCount++
		}
	}
	if requiredCount != 6 {
		t.Errorf("expected 6 required categories, got %d", requiredCount)
	}
}

func TestRequiredCoverage_AllBelow(t *testing.T) {
	rubric := DefaultRubric()
	if RequiredCoverage(rubric) {
		t.Error("should not pass with all scores at 0")
	}
}

func TestRequiredCoverage_AllAbove(t *testing.T) {
	rubric := DefaultRubric()
	for i := range rubric {
		if rubric[i].Required {
			rubric[i].Score = 80
		}
	}
	if !RequiredCoverage(rubric) {
		t.Error("should pass with all required scores at 80")
	}
}

func TestRequiredCoverage_OneBelow(t *testing.T) {
	rubric := DefaultRubric()
	for i := range rubric {
		if rubric[i].Required {
			rubric[i].Score = 90
		}
	}
	rubric[1].Score = 50
	if RequiredCoverage(rubric) {
		t.Error("should not pass with one required score below threshold")
	}
}

func TestLowestScoring_ReturnsRequired(t *testing.T) {
	rubric := DefaultRubric()
	rubric[0].Score = 90
	rubric[1].Score = 10
	rubric[2].Score = 30
	rubric[3].Score = 50
	rubric[4].Score = 70
	rubric[5].Score = 85

	gaps := LowestScoring(rubric, 3)
	if len(gaps) != 3 {
		t.Fatalf("expected 3 gaps, got %d", len(gaps))
	}
	if gaps[0].ID != "R2" {
		t.Errorf("expected R2 as lowest, got %s", gaps[0].ID)
	}
}

func TestLowestScoring_FallsToOptional(t *testing.T) {
	rubric := DefaultRubric()
	for i := range rubric {
		if rubric[i].Required {
			rubric[i].Score = 100
		}
	}

	gaps := LowestScoring(rubric, 3)
	for _, g := range gaps {
		if g.Required {
			t.Errorf("should return optional categories when all required are met, got required %s", g.ID)
		}
	}
}

func TestCoverageReport_Format(t *testing.T) {
	rubric := DefaultRubric()
	rubric[0].Score = 100
	rubric[1].Score = 0
	rubric[2].Score = 50
	rubric[3].Score = 80

	report := CoverageReport(rubric)
	if report == "" {
		t.Error("coverage report should not be empty")
	}
	if len(report) < 100 {
		t.Error("coverage report seems too short")
	}
}

func TestPrefillCategories(t *testing.T) {
	rubric := DefaultRubric()
	scan := &ScanResult{
		SourceStack: "Java (28 files), Maven",
		EntryPoints: []string{"RestApplication -> src/main/java/RestApplication.java"},
		CustomDeps:  []string{"lib/audit-logging-library-1.0.0.jar"},
	}

	PrefillCategories(rubric, scan)

	r1 := findCategory(rubric, "R1")
	if r1.Score != 70 {
		t.Errorf("R1 should be 70 after pre-fill, got %d", r1.Score)
	}

	r3 := findCategory(rubric, "R3")
	if r3.Score != 30 {
		t.Errorf("R3 should be 30 with custom deps detected, got %d", r3.Score)
	}

	r4 := findCategory(rubric, "R4")
	if r4.Score != 0 {
		t.Errorf("R4 should be 0 (LLM analyzes), got %d", r4.Score)
	}

	r5 := findCategory(rubric, "R5")
	if r5.Score != 0 {
		t.Errorf("R5 should be 0 (LLM analyzes), got %d", r5.Score)
	}

	r6 := findCategory(rubric, "R6")
	if r6.Score < 90 {
		t.Errorf("R6 should be >=90 with entry points, got %d", r6.Score)
	}
}

func TestPrefillCategories_NoCustomDeps(t *testing.T) {
	rubric := DefaultRubric()
	scan := &ScanResult{
		SourceStack: "Java (28 files)",
		CustomDeps:  []string{},
	}

	PrefillCategories(rubric, scan)

	r3 := findCategory(rubric, "R3")
	if r3.Score != 80 {
		t.Errorf("R3 should be 80 when no custom deps, got %d", r3.Score)
	}
}

func TestUpdateScoresFromResponse(t *testing.T) {
	rubric := DefaultRubric()
	rubric[1].Score = 30

	qar := QARound{
		Round: 1,
		Answers: []string{
			"[R2] Quarkus 3.x with resteasy-jackson",
		},
	}

	updateScoresFromResponse(rubric, qar)

	r2 := findCategory(rubric, "R2")
	if r2.Score != 70 {
		t.Errorf("R2 should be 70 after answer (+40), got %d", r2.Score)
	}
}

func TestUpdateScoresFromResponse_Cap100(t *testing.T) {
	rubric := DefaultRubric()
	rubric[1].Score = 80

	qar := QARound{
		Round: 1,
		Answers: []string{
			"[R2] Quarkus 3.17 with all extensions",
		},
	}

	updateScoresFromResponse(rubric, qar)

	r2 := findCategory(rubric, "R2")
	if r2.Score > 100 {
		t.Errorf("score should cap at 100, got %d", r2.Score)
	}
}

func findCategory(rubric []Category, id string) *Category {
	for i := range rubric {
		if rubric[i].ID == id {
			return &rubric[i]
		}
	}
	return nil
}
