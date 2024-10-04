package sim

import (
	"image/color"
	"math"
)

const (
	mass    float64 = 1e24
	kmToM   float64 = 1000
	fromSun float64 = 1e6 * kmToM
)

func BodiesSol() []*Body {
	bodies := []*Body{
		{"SOL", 1988400 * mass, 695700 * kmToM, Location{}, Velocity{}},
		{"MERCURY", 0.330 * mass, 4879 * kmToM / 2, Location{57.9 * fromSun, 0}, Velocity{0, 47.4 * kmToM}},
		{"VENUS", 4.87 * mass, 12104 * kmToM / 2, Location{108.2 * fromSun, 0}, Velocity{0, 35.0 * kmToM}},
		{"EARTH", 5.97 * mass, 12756 * kmToM / 2, Location{149.6 * fromSun, 0}, Velocity{0, 29.8 * kmToM}},
		{"MARS", 0.642 * mass, 6792 * kmToM / 2, Location{228.0 * fromSun, 0}, Velocity{0, 24.1 * kmToM}},
		{"JUPITER", 1898 * mass, 142984 * kmToM / 2, Location{778.5 * fromSun, 0}, Velocity{0, 13.1 * kmToM}},
		{"SATURN", 568 * mass, 120536 * kmToM / 2, Location{1432.0 * fromSun, 0}, Velocity{0, 9.7 * kmToM}},
		{"URANUS", 86.8 * mass, 51118 * kmToM / 2, Location{2867.0 * fromSun, 0}, Velocity{0, 6.8 * kmToM}},
		{"NEPTUNE", 102 * mass, 49528 * kmToM / 2, Location{4515.0 * fromSun, 0}, Velocity{0, 5.4 * kmToM}},
		{"PLUTO", 0.0130 * mass, 2376 * kmToM / 2, Location{5906.4 * fromSun, 0}, Velocity{0, 4.7 * kmToM}},
	}
	return bodies
}

func BodiesSunTest() []*Body {
	bodies := []*Body{
		{"SOL", 1988400 * mass, 695700 * kmToM, Location{549.6 * fromSun, -300 * fromSun}, Velocity{}},
		{"SOL", 988400 * mass, 695700 * kmToM, Location{-349.6 * fromSun, 500 * fromSun}, Velocity{}},
		{"SOL", 1088400 * mass, 695700 * kmToM, Location{550.6 * fromSun, 430 * fromSun}, Velocity{}},
	}
	return bodies
}

func ScaledRadius(body *Body, min, max float64) float64 {
	minRadius := (2376.0 / 2) * kmToM // Pluto's radius
	maxRadius := 695700.0 * kmToM     // Sol's radius
	return (body.Radius-minRadius)/(maxRadius-minRadius)*(max-min) + min
}

func LogScaledRadius(body *Body, min, max float64) float64 {
	maxRadius := 695700.0 * kmToM // Sol's radius
	vector := math.Sqrt(body.Radius / maxRadius)
	scaledRadius := vector*(max-min) + min
	return scaledRadius
}

func ScaledLocation(body *Body, scaledMax float64) (float64, float64) {
	maxLocation := 6500 * fromSun
	scaledX := (body.Location.X / maxLocation) * scaledMax
	scaledY := (body.Location.Y / maxLocation) * scaledMax
	return scaledX, scaledY
}

func LogScaledLocation(body *Body, scaledMax float64) (float64, float64) {
	maxLocation := 6000.0 * fromSun
	vector := math.Sqrt(body.Location.X*body.Location.X + body.Location.Y*body.Location.Y)
	scaledVector := math.Sqrt(vector/maxLocation) * scaledMax
	scaledX := scaledVector * body.Location.X / vector
	scaledY := scaledVector * body.Location.Y / vector
	return scaledX, scaledY
}

func BodyColor(name string) color.Color {
	switch name {
	case "SOL":
		return color.RGBA{R: 255, G: 255, B: 0, A: 255} // Yellow
	case "MERCURY":
		return color.RGBA{R: 128, G: 128, B: 128, A: 255} // Grey
	case "VENUS":
		return color.RGBA{R: 255, G: 255, B: 255, A: 255} // White
	case "EARTH":
		return color.RGBA{R: 0, G: 0, B: 255, A: 255} // Blue
	case "MARS":
		return color.RGBA{R: 255, G: 0, B: 0, A: 255} // Red
	case "JUPITER":
		return color.RGBA{R: 255, G: 200, B: 100, A: 255} // Banded pattern, approximated as a warm yellow-brown
	case "SATURN":
		return color.RGBA{R: 200, G: 150, B: 100, A: 255} // Yellow-brown
	case "URANUS":
		return color.RGBA{R: 0, G: 128, B: 128, A: 255} // Blue-green
	case "NEPTUNE":
		return color.RGBA{R: 0, G: 0, B: 128, A: 255} // Deep blue
	case "PLUTO":
		return color.RGBA{R: 128, G: 128, B: 128, A: 255} // Grey-brown
	default:
		return color.White
	}
}
