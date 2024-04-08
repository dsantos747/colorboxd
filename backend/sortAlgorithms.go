package colorboxd

import (
	"math"

	"github.com/lucasb-eyer/go-colorful"
)

// Returns the color vividness according to the vividness equation proposed by Tieling Chen (top dawg)
// in his paper "A measurement of the overall vividness of a color image based on RGB color model"
func checkVividness(rgb colorful.Color) float64 {
	r := rgb.R
	g := rgb.G
	b := rgb.B
	if r == g && g == b {
		return 0
	}
	return math.Sqrt(math.Pow(r, 2) + math.Pow(g, 2) + math.Pow(b, 2) - r*g - r*b - g*b)
}

// Chooses the first colour which satisfies vividness and count ratio
// If no match, return nil.
func dominantVividColour(colors []Color, vividness float64) *Color {
	for i, c := range colors {
		ratio := float64(colors[0].count) / float64(c.count)
		if checkVividness(c.rgb) >= vividness && ratio <= 1.5 {
			return &colors[i]
		}
	}
	return nil
}

// Chooses the first colour which satisfies vividness and count ratio
// If no match, returns the most dominant colour.
func getDominantVividColour(colors []Color, vividness float64) Color {
	color := dominantVividColour(colors, vividness)
	if color != nil {
		return *color
	}
	return colors[0]
}

// Sorts by the hue of the most dominant colour
func AlgoHue(colors []Color) int {
	return int(colors[0].h * 100)
}

// Sorts by the luminosity of the most dominant colour
func AlgoLuminosity(colors []Color) int {
	return int(colors[0].l * 100)
}

// Sorts by the most common hue with high vividness, provided there
// are enough occurrences of it relative to the most dominant hue.
func AlgoBrightDominantHue(colors []Color) int {
	color := getDominantVividColour(colors, 0.25)
	return int(color.h * 100)
}

// Sorts by the inverse-step sorting method, using just the single most dominant colour.
func AlgoInverseStep(colors []Color, reps int) int {
	color := colors[0]

	h2 := int((color.h / 360) * float64(reps))
	l2 := int(color.l * float64(reps))
	v2 := int(color.v * float64(reps))

	if h2%2 == 1 {
		v2 = reps - v2
		l2 = reps - l2
	}

	return 10000*h2 + 100*l2 + v2 // Factors of 100 to allow reps >= 10
}

// Sorts by the inverse-step sorting method, using the most common vivid colour.
func AlgoInverseStepV2(colors []Color, reps int) int {
	color := getDominantVividColour(colors, 0.25)

	h2 := int((color.h / 360) * float64(reps))
	l2 := int(color.l * float64(reps))
	v2 := int(color.v * float64(reps))

	if h2%2 == 1 {
		v2 = reps - v2
		l2 = reps - l2
	}

	return 10000*h2 + 100*l2 + v2 // Factors of 100 to allow reps >= 10
}

// Sorts by putting all the blacks first, then colours (red to blue), then whites.
// Hence, BRBW. Blacks are further sorted so that, if a secondary colour is present,
// that is sorted inversely. White is sorted with a similar logic. The full range is
// Black.blue -> Black.red -> Red -> Blue -> White.blue -> White.red.
func AlgoBRBW1(colors []Color) int {
	score := 0
	order := math.Pow(100, 3)

	domCol := dominantVividColour(colors, 0.15) // Reduced vividness threshold, to reduce amount of posters in black/white zones. Consider adding a secondary lum check here
	if domCol != nil {
		return int(((*domCol).h / 360) * order)
	}

	domCol = dominantVividColour(colors, 0.05)
	if domCol != nil && domCol.l > 0.05 && domCol.l < 0.85 { // Lum check, to reduce amount in black/white zones
		return int(((*domCol).h / 360) * order)
	}

	if colors[0].l > 0.5 {
		score = 1000000
	} else {
		score = -10000
	}

	order = math.Pow(100, 2)
	for i := 1; i < len(colors); i++ {
		if checkVividness(colors[i].rgb) >= 0.01 {
			score += int((colors[i].h / 360) * order)
			return score
		}
	}

	// Tack posters with no vivid colours onto the very ends
	if score == -10000 {
		lum1 := colors[0].l
		score -= int(100 * lum1) // white-black - then the remaining blacks
	}
	if score == 1000000 {
		lum1 := colors[0].l
		score += (100 - int(100*lum1)) // the remaining whites - then white-black
	}

	return score
}

// Similar to BRBW, with a different method of reducing the amount of films in the black/white zone
func AlgoBRBW2(colors []Color) int {
	score := 0
	order := math.Pow(100, 3)

	domCol := dominantVividColour(colors, 0.08) // Greatly reduced vividness threshold, to reduce amount of posters in black/white zones.
	if domCol != nil {
		return int(((*domCol).h / 360) * order)
	}

	if colors[0].l > 0.5 {
		score = 1000000
	} else {
		score = -10000
	}

	order = math.Pow(100, 2)
	for i := 1; i < len(colors); i++ {
		if checkVividness(colors[i].rgb) >= 0.01 {
			score += int((colors[i].h / 360) * order)
			return score
		}
	}

	// Tack posters with no vivid colours onto the very ends
	if score == -10000 {
		lum1 := colors[0].l
		score -= int(100 * lum1) // white-black - then the remaining blacks
	}
	if score == 1000000 {
		lum1 := colors[0].l
		score += (100 - int(100*lum1)) // the remaining whites - then white-black
	}

	return score
}
