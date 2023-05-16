package camera

type Command interface {
	Type() CommandType
	String() string
}

type CommandType int

const (
	EmptyCommandType CommandType = iota
	PTZMoveCommandType
	PTZPresetCommandType
	SetDeviceCommandType
)

func (a CommandType) String() string {
	return [...]string{"EmptyCommandType", "PTZMoveCommandType", "PTZPresetCommandType"}[a]
}
