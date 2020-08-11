package herbsystem

type Stage int

const (
	StageNew Stage = iota
	StageReady
	StageConfiguring
	StageStarting
	StageRunning
	StageStoping
)

var StageNames = map[Stage]string{
	StageNew:         "New",
	StageReady:       "Ready",
	StageConfiguring: "Configuring",
	StageStarting:    "Starting",
	StageRunning:     "Running",
	StageStoping:     "Stoping",
}
