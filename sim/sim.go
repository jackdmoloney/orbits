package sim

import "math"

type Location struct{ X, Y float64 }
type Velocity struct{ X, Y float64 }

type Body struct {
	Name     string
	Mass     float64
	Radius   float64
	Location Location
	Velocity Velocity
}

type Simulator struct {
	Bodies []*Body
}

func (s *Simulator) Step(timeStep float64) {
	s.calculateAccelerations(timeStep)
	s.progress(timeStep)
}

func (s *Simulator) progress(timeStep float64) {
	for _, body := range s.Bodies {
		body.Location.X += body.Velocity.X * timeStep
		body.Location.Y += body.Velocity.Y * timeStep
	}
}

func (s *Simulator) calculateAccelerations(timeStep float64) {
	G := 6.67430e-11 // gravitational constant

	for i, sourceBody := range s.Bodies {
		sourceBodyAcceleration := Velocity{X: 0, Y: 0}
		for j, targetBody := range s.Bodies {
			if i == j { // don't calculate acceleration of body on itself
				continue
			}

			distance := sourceBody.distanceTo(targetBody)
			accelerationMagnitude := G * targetBody.Mass / (distance * distance)
			accelerationDirectionX := (targetBody.Location.X - sourceBody.Location.X) / distance
			accelerationDirectionY := (targetBody.Location.Y - sourceBody.Location.Y) / distance

			sourceBodyAcceleration.X += accelerationMagnitude * accelerationDirectionX * timeStep
			sourceBodyAcceleration.Y += accelerationMagnitude * accelerationDirectionY * timeStep
		}

		// store the calculated acceleration on sourceBody
		sourceBody.Velocity.X += sourceBodyAcceleration.X
		sourceBody.Velocity.Y += sourceBodyAcceleration.Y
	}
}

func (b *Body) distanceTo(other *Body) float64 {
	dx := other.Location.X - b.Location.X
	dy := other.Location.Y - b.Location.Y
	distanceSquared := dx*dx + dy*dy
	distance := math.Sqrt(distanceSquared)

	minDistance := b.Radius + other.Radius
	if distance < minDistance {
		distance = minDistance
	}

	return distance
}
