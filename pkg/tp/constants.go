package tp

// HitObjectType is a bit-flag enum for the type of hit object.
type HitObjectType int

const (
	HitObjectTypeNormal         HitObjectType = 1
	HitObjectTypeSlider         HitObjectType = 2
	HitObjectTypeNewCombo       HitObjectType = 4
	HitObjectTypeNormalNewCombo HitObjectType = 5 // Normal | NewCombo
	HitObjectTypeSliderNewCombo HitObjectType = 6 // Slider | NewCombo
	HitObjectTypeSpinner        HitObjectType = 8
	HitObjectTypeColourHax      HitObjectType = 112
	HitObjectTypeHold           HitObjectType = 128
	HitObjectTypeManiaLong      HitObjectType = 128
)

func NewHitObjectType(objectType string, newCombo bool) HitObjectType {
	var result HitObjectType

	switch objectType {
	case "slider":
		result = HitObjectTypeSlider
	case "spinner":
		result = HitObjectTypeSpinner
	case "circle":
		result = HitObjectTypeNormal
	default:
		result = HitObjectTypeNormal
	}

	if newCombo {
		result |= HitObjectTypeNewCombo
	}
	return result
}

// HitObjectSoundType is a bit-flag enum for hitsounds.
type HitObjectSoundType int

const (
	HitObjectSoundTypeNone    HitObjectSoundType = 0
	HitObjectSoundTypeNormal  HitObjectSoundType = 1
	HitObjectSoundTypeWhistle HitObjectSoundType = 2
	HitObjectSoundTypeFinish  HitObjectSoundType = 4
	HitObjectSoundTypeClap    HitObjectSoundType = 8
)

func NewHitObjectSoundType(soundType []string) HitObjectSoundType {
	var result HitObjectSoundType
	for _, sound := range soundType {
		switch sound {
		case "normal":
			result |= HitObjectSoundTypeNormal
		case "whistle":
			result |= HitObjectSoundTypeWhistle
		case "finish":
			result |= HitObjectSoundTypeFinish
		case "clap":
			result |= HitObjectSoundTypeClap
		default:
			// idk how that would happen but ok
		}
	}
	return result
}
