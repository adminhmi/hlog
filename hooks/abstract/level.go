package abstract

type LevelCheck struct {
	level Level
}

func NewLevelCheck(level Level) LevelCheck {
	return LevelCheck{
		level: level,
	}
}

// Level are all possible logging levels in increasing order, starting with DebugLevel
type Level int

const (
	DebugLevel Level = iota - 1
	InfoLevel
	WarnLevel
	ErrorLevel
	PanicLevel
	FatalLevel
)

// Check returns true if the supplied logging level should be logged.
// Because the logging levels are increasing this is a simple greater equals check.
func (l LevelCheck) Check(level Level) bool {
	return level >= l.level
}
