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

func WrapStartAction(h func()) *Action {
	return CreateAction(CommandStart, Wrap(h))
}

func WrapStopAction(h func()) *Action {
	return CreateAction(CommandStop, Wrap(h))
}

func WrapStartOrPanicAction(h func() error) *Action {
	return CreateAction(CommandStart, WrapOrPanic(h))
}

func WrapStopOrPanicAction(h func() error) *Action {
	return CreateAction(CommandStop, WrapOrPanic(h))
}
