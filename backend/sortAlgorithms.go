package colorboxd

import (
	"fmt"
	"math"
)

// isVividEnough uses a predefined relationship between saturation and luminosity, to determine if the current
// saturation and luminosity result in a color which is vivid enough. The relationship was determined by fitting
// a curve to a scatter plot of sat/lum points which were considered the limit of vividness.
func isVividEnough(S, L float64) bool {
	// sLimit := 0.635 - (2.71 * L) + (3.34 * math.Pow(L, 2))
	sLimit := 1.37 - (17.3 * L) + (101 * math.Pow(L, 2)) - (303 * math.Pow(L, 3)) + (473 * math.Pow(L, 4)) - (366 * math.Pow(L, 5)) + (112 * math.Pow(L, 6))
	return S >= sLimit
}

func AlgoBrightDominantHue(colors []Color) float64 {
	prevColorCount := 1.0
	for _, col := range colors {
		if isVividEnough(col.s, col.l) && float64(col.count)/prevColorCount > 0.5 {
			// fmt.Printf(" colour %d satisfies BrightDominantHue", i+1)
			return col.h
		}
		prevColorCount = float64(col.count)
	}
	return colors[0].h
}

// AlgoInverseStep implements the inverse-step sorting method, using just the single most dominant colour.
// It'd be interesting to try to create an inverseStep method that takes into account isVividEnough
func AlgoInverseStep(color Color, reps int) int {
	// fmt.Println(color.hex)
	// fmt.Printf("h: %f; l: %f; v: %f.\n", color.h, color.l, color.v)
	h2 := int((color.h / 360) * float64(reps))
	l2 := int(color.l * float64(reps))
	v2 := int(color.v * float64(reps))

	if h2%2 == 1 {
		v2 = reps - v2
		l2 = reps - l2
	}
	// fmt.Printf("h2: %v; l2: %v; v2: %v.\n", h2, l2, v2)
	// fmt.Printf("h2: %v; l2: %v; v2: %v.\n", 1000*h2, 10*l2, v2)

	return 100*h2 + 10*l2 + v2 // Factors of 10 works since reps < 10 (so h2, l2, v2 will always be single digit)
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
			if isVividEnough(s, l) {
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
			if isVividEnough(s, l) {
				score += int((h / 360) * math.Pow(100, exp))
				fmt.Print("---------------dark: ")
				break
			} else {
				// fmt.Print("still not vivid enough")
			}
		} else if score >= 1000000 {
			if isVividEnough(s, l) {
				score += int(math.Pow(100, exp) - ((h / 360) * math.Pow(100, exp)))
				fmt.Print("---------------light: ")
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
		fmt.Println("No colours, mostly dark vvvvvvvvvvvvvvvvvv")
	}
	if score == 1000000 {
		lum1 := colors[0].l
		score += int(100 * lum1)
		fmt.Println("No colours, mostly light vvvvvvvvvvvvvvvvvv")
	}

	fmt.Printf("%v,", score)

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
		exp := float64(iterations - i)
		h := colors[i].h
		s := colors[i].s
		l := colors[i].l

		if i == 0 {
			if isVividEnough(s, l) {
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
			if isVividEnough(s, l) {
				score += int((h / 360) * math.Pow(100, exp))
				fmt.Printf("color %v; code %s; count %v. ", i+1, colors[i].hex, colors[i].count)
				break
			} else {
				// fmt.Print("still not vivid enough")
			}
		} else if score >= 1000000 {
			if isVividEnough(s, l) {
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
