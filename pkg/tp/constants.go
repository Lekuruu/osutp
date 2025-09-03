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

const (
	NoMod       uint32 = 0
	NoFail      uint32 = 1 << 0
	Easy        uint32 = 1 << 1
	NoVideo     uint32 = 1 << 2
	Hidden      uint32 = 1 << 3
	HardRock    uint32 = 1 << 4
	SuddenDeath uint32 = 1 << 5
	DoubleTime  uint32 = 1 << 6
	Relax       uint32 = 1 << 7
	HalfTime    uint32 = 1 << 8
	Nightcore   uint32 = 1 << 9
	Flashlight  uint32 = 1 << 10
	Autoplay    uint32 = 1 << 11
	SpunOut     uint32 = 1 << 12
	Autopilot   uint32 = 1 << 13
	Perfect     uint32 = 1 << 14
	Key4        uint32 = 1 << 15
	Key5        uint32 = 1 << 16
	Key6        uint32 = 1 << 17
	Key7        uint32 = 1 << 18
	Key8        uint32 = 1 << 19
	FadeIn      uint32 = 1 << 20
	Random      uint32 = 1 << 21
	Cinema      uint32 = 1 << 22
	Target      uint32 = 1 << 23
	Key9        uint32 = 1 << 24
	KeyCoop     uint32 = 1 << 25
	Key1        uint32 = 1 << 26
	Key3        uint32 = 1 << 27
	Key2        uint32 = 1 << 28
	ScoreV2     uint32 = 1 << 29
	Mirror      uint32 = 1 << 30
)
