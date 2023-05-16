package controller

import (
	"strconv"
	"strings"
)

func NewKeyboardProfileKey(s string) *KeyboardProfileKey {
	keyCode, err := strconv.ParseInt(s, 0, 8)
	if err == nil {
		return &KeyboardProfileKey{
			isModifier: false,
			code:       byte(keyCode),
			Line:       s,
		}
	}

	if strings.HasPrefix(s, "m") {
		mask, err := strconv.ParseInt(s[1:], 0, 8)
		if err == nil {
			return &KeyboardProfileKey{
				isModifier: true,
				mask:       byte(mask),
				Line:       s,
			}
		}
	}

	return nil
}

// Use only one of mask and code
type KeyboardProfileKey struct {
	isModifier bool
	mask       byte
	code       byte
	Line       string
}

func (k *KeyboardProfileKey) KeyIsPressed(modyByte byte, pressed map[byte]struct{}) bool {
	if k.isModifier {
		return modyByte&k.mask != 0
	}
	_, isPressed := pressed[k.code]
	return isPressed
}
