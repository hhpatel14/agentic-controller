package ingest

func DefaultRubric() []Category {
	return []Category{
		{ID: "R1", Name: "Source stack confirmation", Required: true, Prefillable: true},
		{ID: "R2", Name: "Target stack details", Required: true, Prefillable: false},
		{ID: "R3", Name: "Custom/proprietary dependencies", Required: true, Prefillable: true},
		{ID: "R4", Name: "Messaging & async patterns", Required: true, Prefillable: true},
		{ID: "R5", Name: "Data persistence strategy", Required: true, Prefillable: true},
		{ID: "R6", Name: "Entry point & dependency order", Required: true, Prefillable: true},
		{ID: "O1", Name: "Authentication & authorization", Required: false, Prefillable: false},
		{ID: "O2", Name: "External integrations", Required: false, Prefillable: false},
		{ID: "O3", Name: "Deployment target", Required: false, Prefillable: false},
		{ID: "O4", Name: "Design preferences", Required: false, Prefillable: false},
		{ID: "O5", Name: "Team context", Required: false, Prefillable: false},
		{ID: "O6", Name: "Timeline & constraints", Required: false, Prefillable: false},
	}
}

func RequiredCoverage(categories []Category) bool {
	for _, c := range categories {
		if c.Required && c.Score < CoverageThreshold {
			return false
		}
	}
	return true
}

func LowestScoring(categories []Category, n int) []Category {
	var required []Category
	for _, c := range categories {
		if c.Required && c.Score < CoverageThreshold {
			required = append(required, c)
		}
	}

	if len(required) == 0 {
		var optional []Category
		for _, c := range categories {
			if !c.Required && c.Score < CoverageThreshold {
				optional = append(optional, c)
			}
		}
		required = optional
	}

	sortByScore(required)

	if len(required) > n {
		return required[:n]
	}
	return required
}

func sortByScore(cats []Category) {
	for i := 0; i < len(cats); i++ {
		for j := i + 1; j < len(cats); j++ {
			if cats[j].Score < cats[i].Score {
				cats[i], cats[j] = cats[j], cats[i]
			}
		}
	}
}

func CoverageReport(categories []Category) string {
	var report string
	report += "Category                              Score   Status\n"
	report += "───────────────────────────────────────────────────────\n"
	for _, c := range categories {
		if !c.Required {
			continue
		}
		status := "✗ Needs input"
		if c.Score >= 100 {
			status = "✓ Resolved"
		} else if c.Score >= CoverageThreshold {
			status = "✓ Sufficient"
		} else if c.Score >= 50 {
			status = "~ Partial"
		}
		report += padRight(c.Name, 38) + padLeft(itoa(c.Score)+"%", 5) + "   " + status + "\n"
	}
	return report
}

func padRight(s string, n int) string {
	for len(s) < n {
		s += " "
	}
	return s
}

func padLeft(s string, n int) string {
	for len(s) < n {
		s = " " + s
	}
	return s
}

func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	s := ""
	neg := i < 0
	if neg {
		i = -i
	}
	for i > 0 {
		s = string(rune('0'+i%10)) + s
		i /= 10
	}
	if neg {
		s = "-" + s
	}
	return s
}
