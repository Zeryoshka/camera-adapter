package camera

type Command struct {
	Action CameraAction
}

type CameraAction int

const (
	EmptyAction CameraAction = iota
	PanLeftAction
	PanRightAction
	PanStopAction
	TiltUpAction
	TiltDownAction
	TiltStopAction
	ZoomInAction
	ZoomOutAction
	ZoomStopAction
)

func (a CameraAction) String() string {
	return [...]string{
		"EmptyAction", "PanLeftAction", "PanRightAction", "PanStopAction", "TiltUpAction", "TiltDownAction", "TiltStopAction",
		"ZoomInAction", "ZoomOutAction", "ZoomStopAction",
	}[a]
}
