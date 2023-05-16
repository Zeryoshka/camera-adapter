package controller

type KeyboardKey byte

const (
	F1Keyboard           KeyboardKey = 58
	F2Keyboard           KeyboardKey = 59
	F3Keyboard           KeyboardKey = 60
	F4Keyboard           KeyboardKey = 61
	F5Keyboard           KeyboardKey = 62
	F6Keyboard           KeyboardKey = 63
	F7Keyboard           KeyboardKey = 64
	F8Keyboard           KeyboardKey = 65
	F9Keyboard           KeyboardKey = 66
	F10Keyboard          KeyboardKey = 67
	F11Keyboard          KeyboardKey = 68
	F12Keyboard          KeyboardKey = 69
	RightArrowKeyboard   KeyboardKey = 79
	LeftArrowKeyboard    KeyboardKey = 80
	DownArrowKeyboard    KeyboardKey = 81
	UpArrowKeyboard      KeyboardKey = 82
	MinusNumLockKeyboard KeyboardKey = 86
	PlusNumLockKeyboard  KeyboardKey = 87
	OneNumLockKeyboard   KeyboardKey = 89
	TwoNumLockKeyboard   KeyboardKey = 90
	ThreeNumLockKeyboard KeyboardKey = 91
	FourNumLockKeyboard  KeyboardKey = 92
	FiveNumLockKeyboard  KeyboardKey = 93
	SixNumLockKeyboard   KeyboardKey = 94
	SevenNumLockKeyboard KeyboardKey = 95
	EightNumLockKeyboard KeyboardKey = 96
	NineNumLockKeyboard  KeyboardKey = 97
	ZeroNumLockKeyboard  KeyboardKey = 98
)

type StatusKeyKeyboardMask byte

const (
	RightControlKeyboardMask StatusKeyKeyboardMask = 0b00010000
)
