package colorboxd

import (
	"fmt"
	"math"

	"github.com/lucasb-eyer/go-colorful"
)

// isVividEnough uses a predefined relationship between saturation and luminosity, to determine if the current
// saturation and luminosity result in a color which is vivid enough. The relationship was determined by fitting
// a curve to a scatter plot of sat/lum points which were considered the limit of vividness.
// mode=1 is used for determining if most dominant colour is vivid enough.
// mode=2 is used for determining if second most dominant colour is vivid enough.
func isVividEnough(S, L float64, mode int) bool {
	var sLimit float64
	switch mode {
	case 1:
		sLimit = 1.37 - (17.3 * L) + (101 * math.Pow(L, 2)) - (303 * math.Pow(L, 3)) + (473 * math.Pow(L, 4)) - (366 * math.Pow(L, 5)) + (112 * math.Pow(L, 6))
	case 2:
		// sLimit = 1.65 - (15.2 * L) + (93.7 * math.Pow(L, 2)) - (307 * math.Pow(L, 3)) + (527 * math.Pow(L, 4)) - (441 * math.Pow(L, 5)) + (143 * math.Pow(L, 6))
		sLimit = 1.6 - (15.2 * L) + (93.7 * math.Pow(L, 2)) - (307 * math.Pow(L, 3)) + (527 * math.Pow(L, 4)) - (441 * math.Pow(L, 5)) + (143 * math.Pow(L, 6))
	}
	return S >= sLimit
}

func CheckVividness(rgb colorful.Color) float64 {
	r := rgb.R
	g := rgb.G
	b := rgb.B
	if r == g && g == b {
		return 0
	}
	return math.Sqrt(math.Pow(r, 2) + math.Pow(g, 2) + math.Pow(b, 2) - r*g - r*b - g*b)
}

func chooseDominantVividColour(colors []Color) Color {
	for _, c := range colors {
		ratio := float64(colors[0].count) / float64(c.count)
		if CheckVividness(c.rgb) >= 0.25 && ratio <= 1.5 {
			return c
		}
	}
	return colors[0]
}

// Sorts by the hue of the most common colour
func AlgoHue(colors []Color) float64 {
	return colors[0].h
}

// Sorts by the most common hue with high vividness, provided there
// are enough occurrences of it relative to the most common hue.
func AlgoBrightDominantHue(colors []Color) float64 {
	color := chooseDominantVividColour(colors)
	return color.h
}

// Sorts by the inverse-step sorting method, using just the single most dominant colour.
// It'd be interesting to incoroporate a check if the second most dominant colour is vivid, and frequent enough.
func AlgoInverseStep(colors []Color, reps int) int {
	color := colors[0]

	h2 := int((color.h / 360) * float64(reps))
	l2 := int(color.l * float64(reps))
	v2 := int(color.v * float64(reps))

	// fmt.Printf("h2: %v; l2: %v; v2: %v.", h2, l2, v2)
	// fmt.Printf("---- %v ----", h2/2)

	if h2%2 == 1 {
		v2 = reps - v2
		l2 = reps - l2
	}

	// fmt.Print("    ", 100*h2+10*l2+v2, "     ")

	return 10000*h2 + 100*l2 + v2 // Factors of 100 to allow reps >= 10
}

func AlgoInverseStepV2(colors []Color, reps int) int {
	// fmt.Println(color.hex)
	// fmt.Printf("%s,%v,%f; ", colors[0].hex, colors[0].count, CheckVividness(colors[0].rgb))
	// fmt.Printf("%s,%v,%f; ", colors[1].hex, colors[1].count, CheckVividness(colors[1].rgb))
	// fmt.Printf("%s,%v,%f; ", colors[2].hex, colors[2].count, CheckVividness(colors[2].rgb))

	color := chooseDominantVividColour(colors)

	h2 := int((color.h / 360) * float64(reps))
	l2 := int(color.l * float64(reps))
	v2 := int(color.v * float64(reps))

	if h2%2 == 1 {
		v2 = reps - v2
		l2 = reps - l2
	}
	// fmt.Printf("h2: %v; l2: %v; v2: %v.\n", h2, l2, v2)
	// fmt.Printf("h2: %v; l2: %v; v2: %v.\n", 1000*h2, 10*l2, v2)

	return 10000*h2 + 100*l2 + v2 // Factors of 100 to allow reps >= 10
}

func AlgoBRBW1(colors []Color) int {
	iterations := 3
	if len(colors) < iterations {
		iterations = len(colors)
	}
	score := 0

	for i := 0; i < iterations; i++ {
		exp := float64(iterations - i)
		h := colors[i].h
		s := colors[i].s
		l := colors[i].l

		if i == 0 {
			if isVividEnough(s, l, 1) {
				score += int((h / 360) * math.Pow(100, exp))
				break
			} else {
				if l > 0.5 {
					score += 1000000
				} else {
					score -= 10000
				}
			}
		} else if score <= -10000 {
			if isVividEnough(s, l, 1) {
				score += int((h / 360) * math.Pow(100, exp))
				// fmt.Print("---------------dark: ")
				break
			} else {
				// fmt.Print("still not vivid enough")
			}
		} else if score >= 1000000 {
			if isVividEnough(s, l, 1) {
				score += int(math.Pow(100, exp) - ((h / 360) * math.Pow(100, exp)))
				// fmt.Print("---------------light: ")
				break
			} else {
				// fmt.Print("still not vivid enough")
			}
		}
	}

	// Handle cases where no vivid colours are present
	if score == -10000 {
		lum1 := colors[0].l
		score -= (100 - int(100*lum1))
		// fmt.Println("No colours, mostly dark vvvvvvvvvvvvvvvvvv")
	}
	if score == 1000000 {
		lum1 := colors[0].l
		score += int(100 * lum1)
		// fmt.Println("No colours, mostly light vvvvvvvvvvvvvvvvvv")
	}

	// fmt.Printf("%v,", score)

	return score
}

// This variant takes into account the occurrence count of each color.
// If the second color is vivid and occurs enough, just use that value instead
func AlgoBRBW2(colors []Color) int {
	iterations := 3
	if len(colors) < iterations {
		iterations = len(colors)
	}
	score := 0

	for i := 0; i < iterations; i++ {
		exp := float64(3 - i) // hardcoded 3 used since iterations may not always be 3
		h := colors[i].h
		s := colors[i].s
		l := colors[i].l

		if i == 0 {
			if isVividEnough(s, l, 1) {
				score += int((h / 360) * math.Pow(100, exp))
				break
			} else {
				if l > 0.5 {
					score += 1000000
				} else {
					score -= 10000
				}
			}
		} else if isVividEnough(s, l, 2) && colors[i].count >= 1500 {
			fmt.Print("--------------")
			score = int((h / 360) * math.Pow(100, 3)) // Hardcoded 3 as we want to make this colour top priority
			break
		} else if score <= -10000 {
			if isVividEnough(s, l, 1) {
				score += int((h / 360) * math.Pow(100, exp))
				fmt.Printf("color %v; code %s; count %v. ", i+1, colors[i].hex, colors[i].count)
				break
			} else {
				// fmt.Print("still not vivid enough")
			}
		} else if score >= 1000000 {
			if isVividEnough(s, l, 1) {
				score += int(math.Pow(100, exp) - ((h / 360) * math.Pow(100, exp)))
				fmt.Printf("color %v; code %s; count %v. ", i+1, colors[i].hex, colors[i].count)
				break
			} else {
				// fmt.Print("still not vivid enough")
			}
		}
	}

	// Handle cases where no vivid colours are present
	if score == -10000 {
		lum1 := colors[0].l
		score -= (100 - int(100*lum1))
		fmt.Print("vvvvvvvv mostly dark ")
	}
	if score == 1000000 {
		lum1 := colors[0].l
		score += int(100 * lum1)
		fmt.Print("^^^^^^^^ mostly light ")
	}

	fmt.Printf("Score: %v.", score)

	return score
}
