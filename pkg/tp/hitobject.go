package tp

import osu "github.com/natsukagami/go-osu-parser"

// Vector2 represents a 2D point or vector.
type Vector2 struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// HitObjectBase encapsulates the common fields of all osu! hit objects.
type HitObjectBase struct {
	// The time at which the object is to be hit.
	StartTime int `json:"startTime"`

	// For spannable objects, when the object ends; otherwise equal to StartTime.
	EndTime int `json:"endTime"`

	// The type flags of this object.
	Type HitObjectType `json:"type"`

	// Hitsound flags for this object.
	SoundType HitObjectSoundType `json:"soundType"`

	// Number of segments (e.g. slider repeats +1).
	SegmentCount int `json:"segmentCount"`

	// Length in gamefield pixels (only for spans).
	SpatialLength float64 `json:"spatialLength"`

	// If New Combo, how many extra colours to cycle by.
	ComboColourOffset int `json:"comboColourOffset"`

	// Zero-based index for combo colour (for more than RGB triples).
	ComboColourIndex int `json:"comboColourIndex"`

	// Gamefield position of this object.
	Position Vector2 `json:"position"`

	// For spans, the ending position. For simple objects, may equal Position.
	EndPosition Vector2 `json:"endPosition"`

	// Current height in a stack of notes. Zero means no stack.
	StackCount int `json:"stackCount"`

	// The combo number displayed (one-based).
	ComboNumber int `json:"comboNumber"`

	// True if this object is the last in its combo.
	LastInCombo bool `json:"lastInCombo"`
}

func NewHitObjectBase(hitObject osu.HitObject) *HitObjectBase {
	startTime := hitObject.StartTime
	endTime := hitObject.EndTime
	if endTime == 0 {
		endTime = hitObject.StartTime
	}

	startPosition := Vector2{
		X: hitObject.Position.X,
		Y: hitObject.Position.Y,
	}
	endPosition := Vector2{
		X: hitObject.EndPosition.X,
		Y: hitObject.EndPosition.Y,
	}
	if endPosition.X == 0 && endPosition.Y == 0 {
		endPosition = startPosition
	}

	return &HitObjectBase{
		StartTime:         startTime,
		EndTime:           endTime,
		Type:              NewHitObjectType(hitObject.ObjectName, hitObject.NewCombo),
		SoundType:         NewHitObjectSoundType(hitObject.SoundTypes),
		SegmentCount:      hitObject.RepeatCount + 1,
		SpatialLength:     hitObject.PixelLength,
		Position:          startPosition,
		EndPosition:       endPosition,
		ComboColourOffset: 0,     // TODO
		ComboColourIndex:  0,     // TODO
		StackCount:        0,     // TODO
		ComboNumber:       0,     // TODO
		LastInCombo:       false, // TODO
	}
}
