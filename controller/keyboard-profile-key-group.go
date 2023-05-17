package controller

type KeyTuple struct {
	key   *KeyboardProfileKey
	index uint
}

func NewKeyboardProfileKeyGroup(arrS []string) *KeyboardProfileKeyGroup {
	atLeastOne := false
	keyProfileMap := make(map[byte]*KeyTuple)
	statusKeyProfileMap := make(map[byte]*KeyTuple)
	for i, s := range arrS {
		keyProfile := NewKeyboardProfileKey(s)
		if keyProfile == nil {
			continue
		} else if keyProfile.isModifier {
			statusKeyProfileMap[keyProfile.mask] = &KeyTuple{keyProfile, uint(i)}
			atLeastOne = true
		} else {
			keyProfileMap[keyProfile.code] = &KeyTuple{keyProfile, uint(i)}
			atLeastOne = true
		}
	}
	if !atLeastOne {
		return nil
	}
	return &KeyboardProfileKeyGroup{
		keyProfileMap:       keyProfileMap,
		statusKeyProfileMap: statusKeyProfileMap,
	}
}

type KeyboardProfileKeyGroup struct {
	keyProfileMap       map[byte]*KeyTuple
	statusKeyProfileMap map[byte]*KeyTuple
}

func (g *KeyboardProfileKeyGroup) PressedOneKey(
	statusByte byte, pressed map[byte]struct{},
) (*KeyboardProfileKey, uint) {
	mask := byte(0b1)
	if statusByte != 0 {
		for {
			keyTuple, ok := g.statusKeyProfileMap[mask]
			if ok && keyTuple.key != nil && keyTuple.key.mask == mask {
				return keyTuple.key, keyTuple.index
			}
			if mask == 0b10000000 {
				break
			}
			mask <<= 1
		}
	}
	for keyCode := range pressed {
		keyTuple, ok := g.keyProfileMap[keyCode]
		if ok && keyTuple.key != nil {
			return keyTuple.key, keyTuple.index
		}
	}
	return nil, 0
}
