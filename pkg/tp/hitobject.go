package tp

import (
	"math"

	osu "github.com/natsukagami/go-osu-parser"
)

// Vector2 represents a 2D point or vector.
type Vector2 struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

func (v Vector2) Add(other Vector2) Vector2 {
	return Vector2{X: v.X + other.X, Y: v.Y + other.Y}
}

func (v Vector2) Sub(other Vector2) Vector2 {
	return Vector2{X: v.X - other.X, Y: v.Y - other.Y}
}

func (v Vector2) Mul(scalar float64) Vector2 {
	return Vector2{X: v.X * scalar, Y: v.Y * scalar}
}

func (v Vector2) Length() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

func (v Vector2) Normalize() Vector2 {
	length := v.Length()
	if length == 0 {
		return Vector2{X: 0, Y: 0}
	}
	return Vector2{X: v.X / length, Y: v.Y / length}
}

// TpHitObject represents a hit object with strain values for difficulty calculation
type TpHitObject struct {
	HitObject                  osu.HitObject
	LazySliderLengthFirst      float64
	LazySliderLengthSubsequent float64
	NormalizedStartPosition    Vector2
	NormalizedEndPosition      Vector2
	Strains                    [2]float64
}

/* Helper functions to work with osu.HitObject types */

func getStartPosition(hitObject osu.HitObject) Vector2 {
	return Vector2{X: hitObject.Position.X, Y: hitObject.Position.Y}
}

func getEndPosition(hitObject osu.HitObject) Vector2 {
	endPos := Vector2{X: hitObject.EndPosition.X, Y: hitObject.EndPosition.Y}
	if endPos.X == 0 && endPos.Y == 0 {
		return getStartPosition(hitObject)
	}
	return endPos
}

func getEndTime(hitObject osu.HitObject) int {
	if hitObject.EndTime == 0 {
		return hitObject.StartTime
	}
	return hitObject.EndTime
}

const (
	// Almost the normed diameter of a circle (104 osu pixel). That is -after- position transforming.
	almostDiameter = 90.0

	// Pseudo threshold values to distinguish between "singles" and "streams". Of course the border can not be defined clearly, therefore the algorithm
	// has a smooth transition between those values. They also are based on tweaking and general feedback.
	streamSpacingThreshold = 110.0
	singleSpacingThreshold = 125.0

	// In milliseconds. The smaller the value, the more accurate sliders are approximated. 0 leads to an infinite loop, so use something bigger.
	lazySliderStepLength = 1
)

var (
	// Factor by how much speed / aim strain decays per second. Those values are results of tweaking a lot and taking into account general feedback.
	// Opinionated observation: Speed is easier to maintain than accurate jumps.
	decayBase = [2]float64{0.3, 0.15}

	// Scaling values for weightings to keep aim and speed difficulty in balance. Found from testing a very large map pool (containing all ranked maps) and keeping the
	// average values the same.
	spacingWeightScaling = [2]float64{1400, 26.25}
)

// NewTpHitObject creates a new TpHitObject from an osu.HitObject
func NewTpHitObject(hitObject osu.HitObject, circleRadius float64, beatmap *osu.Beatmap) *TpHitObject {
	obj := &TpHitObject{
		HitObject: hitObject,
		Strains:   [2]float64{1, 1},
	}

	// We will scale everything by this factor, so we can assume a uniform CircleSize among beatmaps.
	scalingFactor := 52.0 / circleRadius
	startPos := getStartPosition(hitObject)
	obj.NormalizedStartPosition = startPos.Mul(scalingFactor)

	hitObjectType := NewHitObjectType(hitObject)

	// Calculate approximation of lazy movement on the slider
	if hitObjectType&HitObjectTypeSlider != 0 {
		// Not sure if this is correct, but here we do not need 100% exact values. This comes pretty darn close in my tests.
		sliderFollowCircleRadius := circleRadius * 3

		segmentLength := hitObject.PixelLength / float64(hitObject.RepeatCount)
		segmentEndTime := float64(hitObject.StartTime) + segmentLength

		// For simplifying this step we use actual osu! coordinates and simply scale the length, that we obtain by the ScalingFactor later
		cursorPos := startPos

		// Actual computation of the first lazy curve
		for time := float64(hitObject.StartTime) + lazySliderStepLength; time < segmentEndTime; time += lazySliderStepLength {
			targetPos := positionAtTime(hitObject, int(time))
			difference := targetPos.Sub(cursorPos)
			distance := difference.Length()

			// Did we move away too far?
			if distance > sliderFollowCircleRadius {
				difference = difference.Normalize() // Obtain the direction of difference. We do no longer need the actual difference
				distance -= sliderFollowCircleRadius
				cursorPos = cursorPos.Add(difference.Mul(distance)) // We move the cursor just as far as needed to stay in the follow circle
				obj.LazySliderLengthFirst += distance
			}
		}

		obj.LazySliderLengthFirst *= scalingFactor

		// If we have an odd amount of repetitions the current position will be the end of the slider. Note that this will -always- be triggered if
		// hitObject.RepeatCount <= 1, because hitObject.RepeatCount can not be smaller than 1. Therefore NormalizedEndPosition will always be initialized
		if hitObject.RepeatCount%2 == 1 {
			obj.NormalizedEndPosition = cursorPos.Mul(scalingFactor)
		}

		// If we have more than one segment, then we also need to compute the length ob subsequent lazy curves. They are different from the first one, since the first
		// one starts right at the beginning of the slider.
		if hitObject.RepeatCount > 1 {
			// Use the next segment
			segmentEndTime += segmentLength

			for time := segmentEndTime - segmentLength + lazySliderStepLength; time < segmentEndTime; time += lazySliderStepLength {
				targetPos := positionAtTime(hitObject, int(time))
				difference := targetPos.Sub(cursorPos)
				distance := difference.Length()

				// Did we move away too far?
				if distance > sliderFollowCircleRadius {
					// Yep, we need to move the cursor
					difference = difference.Normalize() // Obtain the direction of difference. We do no longer need the actual difference
					distance -= sliderFollowCircleRadius
					cursorPos = cursorPos.Add(difference.Mul(distance)) // We move the cursor just as far as needed to stay in the follow circle
					obj.LazySliderLengthSubsequent += distance
				}
			}

			obj.LazySliderLengthSubsequent *= scalingFactor

			// If we have an even amount of repetitions the current position will be the end of the slider
			if hitObject.RepeatCount%2 == 0 {
				obj.NormalizedEndPosition = cursorPos.Mul(scalingFactor)
			}
		}
	} else {
		// We have a normal HitCircle or a spinner
		endPos := getEndPosition(hitObject)
		obj.NormalizedEndPosition = endPos.Mul(scalingFactor)
	}

	return obj
}

// CalculateStrains calculates the strain values for this object
func (obj *TpHitObject) CalculateStrains(previousHitObject *TpHitObject, timeRate float64) {
	obj.calculateSpecificStrain(previousHitObject, DifficultyTypeSpeed, timeRate)
	obj.calculateSpecificStrain(previousHitObject, DifficultyTypeAim, timeRate)
}

// positionAtTime returns the position of a hit object at a given time
// For sliders, this linearly interpolates between start and end position
func positionAtTime(hitObject osu.HitObject, time int) Vector2 {
	hitObjectType := NewHitObjectType(hitObject)
	if hitObjectType&HitObjectTypeSlider == 0 {
		return getStartPosition(hitObject)
	}

	// Simple linear interpolation for sliders,
	// which should match the C# service
	if time <= hitObject.StartTime {
		return getStartPosition(hitObject)
	}

	endTime := getEndTime(hitObject)
	if time >= endTime {
		return getEndPosition(hitObject)
	}

	// Linear interpolation
	startPos := getStartPosition(hitObject)
	endPos := getEndPosition(hitObject)
	progress := float64(time-hitObject.StartTime) / float64(endTime-hitObject.StartTime)
	return Vector2{
		X: startPos.X + (endPos.X-startPos.X)*progress,
		Y: startPos.Y + (endPos.Y-startPos.Y)*progress,
	}
}

func (obj *TpHitObject) calculateSpecificStrain(previousHitObject *TpHitObject, diffType DifficultyType, timeRate float64) {
	timeElapsed := float64(obj.HitObject.StartTime-previousHitObject.HitObject.StartTime) / timeRate
	decay := math.Pow(decayBase[int(diffType)], timeElapsed/1000)
	addition := 1.0

	objType := NewHitObjectType(obj.HitObject)
	if objType&HitObjectTypeSlider != 0 {
		switch diffType {
		case DifficultyTypeSpeed:
			// For speed strain we treat the whole slider as a single spacing entity, since "Speed" is about how hard it is to click buttons fast.
			// The spacing weight exists to differentiate between being able to easily alternate or having to single.
			totalDistance := previousHitObject.LazySliderLengthFirst +
				previousHitObject.LazySliderLengthSubsequent*float64(previousHitObject.HitObject.RepeatCount-1) +
				obj.distanceTo(previousHitObject)
			addition = spacingWeight(totalDistance, diffType) * spacingWeightScaling[int(diffType)]
		case DifficultyTypeAim:
			// For Aim strain we treat each slider segment and the jump after the end of the slider as separate jumps, since movement-wise there is no difference
			// to multiple jumps.
			addition = (spacingWeight(previousHitObject.LazySliderLengthFirst, diffType) +
				spacingWeight(previousHitObject.LazySliderLengthSubsequent, diffType)*float64(previousHitObject.HitObject.RepeatCount-1) +
				spacingWeight(obj.distanceTo(previousHitObject), diffType)) * spacingWeightScaling[int(diffType)]
		}
	} else if objType&HitObjectTypeNormal != 0 {
		addition = spacingWeight(obj.distanceTo(previousHitObject), diffType) * spacingWeightScaling[int(diffType)]
	}

	// Scale addition by the time, that elapsed. Filter out HitObjects that are too close to be played anyway to avoid crazy values by division through close to zero.
	// You will never find maps that require this amongst ranked maps.
	addition /= math.Max(timeElapsed, 50)

	obj.Strains[int(diffType)] = previousHitObject.Strains[int(diffType)]*decay + addition
}

func (obj *TpHitObject) distanceTo(other *TpHitObject) float64 {
	// Scale the distance by circle size
	return obj.NormalizedStartPosition.Sub(other.NormalizedEndPosition).Length()
}

func spacingWeight(distance float64, diffType DifficultyType) float64 {
	// Caution: The subjective values are strong with this one
	switch diffType {
	case DifficultyTypeSpeed:
		var weight float64
		if distance > singleSpacingThreshold {
			weight = 2.5
		} else if distance > streamSpacingThreshold {
			weight = 1.6 + 0.9*(distance-streamSpacingThreshold)/(singleSpacingThreshold-streamSpacingThreshold)
		} else if distance > almostDiameter {
			weight = 1.2 + 0.4*(distance-almostDiameter)/(streamSpacingThreshold-almostDiameter)
		} else if distance > almostDiameter/2 {
			weight = 0.95 + 0.25*(distance-almostDiameter/2)/(almostDiameter/2)
		} else {
			weight = 0.95
		}
		return weight
	case DifficultyTypeAim:
		return math.Pow(distance, 0.99)
	default:
		return 0
	}
}
