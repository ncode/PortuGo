package main

type artifact struct {
	Path   string `json:"path"`
	SHA256 string `json:"sha256"`
}

type reference struct {
	Product          string `json:"product"`
	Version          string `json:"version"`
	SourceURL        string `json:"sourceURL"`
	AcquiredDate     string `json:"acquiredDate"`
	ArchiveSHA256    string `json:"archiveSHA256"`
	ExecutableSHA256 string `json:"executableSHA256"`
	OS               string `json:"os"`
	Culture          string `json:"culture"`
}

type review struct {
	Reason string `json:"reason"`
	Link   string `json:"link"`
}

type inventoryItem struct {
	ID     string   `json:"id"`
	Kind   string   `json:"kind"`
	Link   string   `json:"link"`
	Probes []string `json:"probes"`
}

type evidence struct {
	State         string          `json:"state"`
	Accepted      *bool           `json:"accepted,omitempty"`
	CapturedAt    string          `json:"capturedAt,omitempty"`
	Raw           artifact        `json:"raw"`
	Normalized    artifact        `json:"normalized"`
	Normalizer    string          `json:"normalizer"`
	GUIOnly       bool            `json:"guiOnly,omitempty"`
	Screenshot    *artifact       `json:"screenshot,omitempty"`
	Transcription *artifact       `json:"transcription,omitempty"`
	Review        *review         `json:"review,omitempty"`
	Generated     []generatedFile `json:"generated,omitempty"`
}

type diagnostic struct {
	Code   string `json:"code"`
	Line   int    `json:"line"`
	Column int    `json:"column,omitempty"`
}

type generatedFile struct {
	Path    string   `json:"path"`
	Content artifact `json:"content"`
}

// observation also defines the state/host adapter boundary used by group 3.
// Subprocess replay uses the deterministic execution adapter for these channels.
type observation struct {
	ExitCode    int             `json:"exitCode"`
	Stdout      artifact        `json:"stdout"`
	Diagnostics []diagnostic    `json:"diagnostics,omitempty"`
	Generated   []generatedFile `json:"generated,omitempty"`
	State       *artifact       `json:"state,omitempty"`
	HostTrace   *artifact       `json:"hostTrace,omitempty"`
}

type implementation struct {
	State    string      `json:"state"`
	Tests    []string    `json:"tests,omitempty"`
	Expected observation `json:"expected"`
	Review   *review     `json:"review,omitempty"`
}

type probe struct {
	ID             string          `json:"id"`
	OwnerGroup     int             `json:"ownerGroup"`
	Tasks          []string        `json:"tasks"`
	Source         artifact        `json:"source"`
	Input          artifact        `json:"input"`
	Files          []generatedFile `json:"files,omitempty"`
	TimeoutMS      int             `json:"timeoutMS"`
	MaxSteps       uint64          `json:"maxSteps,omitempty"`
	Evidence       evidence        `json:"evidence"`
	Implementation implementation  `json:"implementation"`
}

type disposition struct {
	ID          string `json:"id"`
	Replacement string `json:"replacement,omitempty"`
	Review      review `json:"review"`
}

type manifest struct {
	Version           int             `json:"version"`
	RecorderVersion   string          `json:"recorderVersion"`
	NormalizerVersion string          `json:"normalizerVersion"`
	Reference         reference       `json:"reference"`
	TasksPath         string          `json:"tasksPath"`
	InventorySources  []string        `json:"inventorySources"`
	Inventory         []inventoryItem `json:"inventory"`
	Probes            []probe         `json:"probes"`
	Retired           []disposition   `json:"retired,omitempty"`
}
