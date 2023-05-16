package camera

import "fmt"

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

func (c *PTZMoveCommand) String() string {
	return "PTZMoveCommand"
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

func (c *PTZPresetCommand) String() string {
	return fmt.Sprintf("PTZPresetCommand(num: %d, set: %t)", c.PresetNumber, c.SetPreset)
}
