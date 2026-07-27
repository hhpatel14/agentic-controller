package ingest

const (
	MaxRounds          = 6
	MaxQuestionsPerRnd = 4
	CoverageThreshold  = 80
)

type Category struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Required    bool   `json:"required"`
	Score       int    `json:"score"`
	Prefillable bool   `json:"prefillable"`
	Answers     string `json:"answers,omitempty"`
}

type Assumption struct {
	Description string `json:"description"`
	Confidence  string `json:"confidence"`
	Category    string `json:"category"`
	Risk        string `json:"risk_if_wrong"`
}

type QARound struct {
	Round     int      `json:"round"`
	Questions []string `json:"questions"`
	Answers   []string `json:"answers"`
}

type IngestResult struct {
	SourceStack     string       `json:"source_stack"`
	TargetStack     string       `json:"target_stack"`
	CustomDeps      []CustomDep  `json:"custom_deps"`
	Messaging       string       `json:"messaging"`
	DataPersistence string       `json:"data_persistence"`
	EntryPoint      string       `json:"entry_point"`
	MigrationOrder  []string     `json:"migration_order"`
	DesignDecisions []string     `json:"design_decisions"`
	Assumptions     []Assumption `json:"assumptions"`
	Rounds          []QARound    `json:"rounds"`
	RoundsCompleted int          `json:"rounds_completed"`
	Termination     string       `json:"termination_reason"`
}

type CustomDep struct {
	Name     string `json:"name"`
	Location string `json:"location"`
	Purpose  string `json:"purpose"`
	Strategy string `json:"migration_strategy"`
}
