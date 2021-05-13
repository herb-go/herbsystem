package herbsystem

type Command string

const CommandStart = Command("start")
const CommandStop = Command("stop")

type Action struct {
	Command interface{}
	Process Process
}

func NewAction() *Action {
	return &Action{}
}

func CreateAction(cmd interface{}, p Process) *Action {
	return &Action{
		Command: cmd,
		Process: p,
	}
}

func CreateStartAction(p Process) *Action {
	return CreateAction(CommandStart, p)
}

func CreateStopAction(p Process) *Action {
	return CreateAction(CommandStop, p)
}
