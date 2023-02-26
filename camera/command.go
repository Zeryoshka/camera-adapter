package camera

type Command interface {
	Type() CommandType
}

type CommandType int

const (
	EmptyCommandType CommandType = iota
	PTZMoveCommandType
	PTZPresetCommandType
)

func (a CommandType) String() string {
	return [...]string{"EmptyCommandType", "PTZMoveCommandType", "PTZPresetCommandType"}[a]
}

type PTZMoveCommand struct {
	PanMove  int
	TiltMove int
	ZoomMove int
}

func NewPTZMoveCommand(PanMove, TiltMove, ZoomMove int) *PTZMoveCommand {
	return &PTZMoveCommand{
		PanMove:  PanMove,
		TiltMove: TiltMove,
		ZoomMove: ZoomMove,
	}
}

func (c *PTZMoveCommand) Type() CommandType {
	return PTZMoveCommandType
}

type PTZPresetCommand struct {
	SetPreset    bool
	PresetNumber uint
}

func NewPTZPresetCommand(setPreset bool, PresetNumber uint) *PTZPresetCommand {
	return &PTZPresetCommand{
		SetPreset:    setPreset,
		PresetNumber: PresetNumber,
	}
}

func (c *PTZPresetCommand) Type() CommandType {
	return PTZPresetCommandType
}
