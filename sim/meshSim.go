package sim

import (
	"math"
	"math/rand"
)

type location struct{ x, y float64 }
type velocity struct{ x, y float64 }
type body struct {
	mass, radius float64
	location     location
	velocity     velocity

	gridX, gridY             int
	gridOffsetX, gridOffsetY float64
}

type meshSimulator struct {
	bodies []*body

	gridScale, gridX, gridY int
	massGrid                []float64
	vectorGrid              []velocity
}

func NewMeshSimulator(gridScale, gridX, gridY, bodyCount int) meshSimulator {
	gridPoints := gridX * gridY
	massGrid := make([]float64, gridPoints)
	vectorGrid := make([]velocity, gridPoints)
	bodies := GenerateBodies(gridScale, gridX, gridY, bodyCount)

	return meshSimulator{bodies, gridScale, gridX, gridY, massGrid, vectorGrid}
}

func GenerateBodies(gridScale, gridX, gridY, bodyCount int) []*body {
	bodies := make([]*body, bodyCount)

	for i := range bodies {
		b := &body{}

		// Generate random location within the grid
		b.location.x = float64(rand.Intn(gridX*gridScale)) + rand.Float64()
		b.location.y = float64(rand.Intn(gridY*gridScale)) + rand.Float64()

		// Generate random velocity
		b.velocity.x = rand.Float64()*6 - 3
		b.velocity.y = rand.Float64()*6 - 3

		// Set mass and radius to some default values
		b.mass = rand.Float64()*2 + 1
		b.radius = 1.0

		bodies[i] = b
	}

	return bodies
}

func (b *body) Location() (float64, float64) {
	return b.location.x, b.location.y
}
func (b *body) Mass() float64 {
	return b.mass
}

func (s *meshSimulator) gridCoordinates(index int) (int, int) {
	y := index / s.gridX
	x := index % s.gridX
	return x * s.gridScale, y * s.gridScale
}

func (s *meshSimulator) Bodies() []Body {
	// bodies := make([]Body, s.gridX*s.gridY)
	// for i, _ := range s.massGrid {
	// 	x, y := s.gridCoordinates(i)
	// 	v := s.vectorGrid[i]
	// 	bodies[i] = &body{
	// 		// mass:     mass + 1,
	// 		mass:     math.Sqrt(v.x*v.x+v.y*v.y)*20 + .5,
	// 		location: location{float64(x), float64(y)},
	// 	}
	// }
	// return bodies

	bodies := make([]Body, len(s.bodies))
	for i, b := range s.bodies {
		bodies[i] = b
	}
	return bodies
}

func (s *meshSimulator) Step(timeStep float64) {
	s.calculateBodyGridLocations()
	s.calculateGridMass()
	s.calculateVectorGrid()
	s.calculateBodyAccelerations(timeStep)
	s.calculateBodyLocation(timeStep)
}

func (s *meshSimulator) calculateBodyGridLocations() {
	for _, b := range s.bodies {
		x, y := b.location.x, b.location.y
		gridX := int(x) / s.gridScale
		gridY := int(y) / s.gridScale

		dx := x - float64(gridX*s.gridScale)
		dy := y - float64(gridY*s.gridScale)

		if gridX >= s.gridX {
			gridX = s.gridX - 1
			b.location.x = float64(s.gridX*s.gridScale - 1)
			b.velocity = velocity{}
		} else if gridX < 0 {
			gridX = 0
			b.location.x = 0
			b.velocity = velocity{}
		}
		if gridY >= s.gridY {
			gridY = s.gridY - 1
			b.location.y = float64(s.gridY*s.gridScale - 1)
			b.velocity = velocity{}
		} else if gridY < 0 {
			gridY = 0
			b.location.y = 0
			b.velocity = velocity{}
		}

		b.gridX = gridX
		b.gridY = gridY
		b.gridOffsetX = dx
		b.gridOffsetY = dy
	}
}

func (s *meshSimulator) calculateGridMass() {
	for i := range s.massGrid {
		s.massGrid[i] = 0
	}

	gridScale := float64(s.gridScale)
	for _, b := range s.bodies {
		// Calculate the weights for each grid point
		wx1 := 1 - b.gridOffsetX/gridScale
		wx2 := b.gridOffsetX / gridScale
		wy1 := 1 - b.gridOffsetY/gridScale
		wy2 := b.gridOffsetY / gridScale

		// Calculate the indices of the surrounding grid points
		i1 := b.gridY*s.gridX + b.gridX
		i2 := b.gridY*s.gridX + (b.gridX + 1)
		i3 := (b.gridY+1)*s.gridX + b.gridX
		i4 := (b.gridY+1)*s.gridX + (b.gridX + 1)

		if i1 >= len(s.vectorGrid) || i1 < 0 || i2 >= len(s.vectorGrid) || i2 < 0 ||
			i3 >= len(s.vectorGrid) || i3 < 0 || i4 >= len(s.vectorGrid) || i4 < 0 {
			continue
		}

		// Add the mass to the surrounding grid points
		s.massGrid[i1] += b.mass * wx1 * wy1
		s.massGrid[i2] += b.mass * wx2 * wy1
		s.massGrid[i3] += b.mass * wx1 * wy2
		s.massGrid[i4] += b.mass * wx2 * wy2
	}
}

func (s *meshSimulator) calculateVectorGrid() {
	G := 1.0 // gravitational constant

	for i := range s.vectorGrid {
		s.vectorGrid[i].x = 0
		s.vectorGrid[i].y = 0
	}
	for i := range s.massGrid {
		for j, mass := range s.massGrid {
			if i == j || mass == 0 {
				continue
			}

			dx := float64(((j % s.gridX) - (i % s.gridX)) * s.gridScale)
			dy := float64(((j / s.gridX) - (i / s.gridX)) * s.gridScale)
			distance := math.Sqrt(dx*dx + dy*dy)

			acceleration := G * mass / (distance * distance)
			s.vectorGrid[i].x += acceleration * dx / distance
			s.vectorGrid[i].y += acceleration * dy / distance
		}
	}
}

func (s *meshSimulator) calculateBodyAccelerations(timeStep float64) {
	gridScale := float64(s.gridScale)
	for _, b := range s.bodies {
		ax := 0.0
		ay := 0.0

		i1 := b.gridY*s.gridX + b.gridX
		i2 := b.gridY*s.gridX + (b.gridX + 1)
		i3 := (b.gridY+1)*s.gridX + b.gridX
		i4 := (b.gridY+1)*s.gridX + (b.gridX + 1)

		if i1 >= len(s.vectorGrid) || i1 < 0 || i2 >= len(s.vectorGrid) || i2 < 0 ||
			i3 >= len(s.vectorGrid) || i3 < 0 || i4 >= len(s.vectorGrid) || i4 < 0 {
			continue
		}

		wx1 := 1 - b.gridOffsetX/gridScale
		wx2 := b.gridOffsetX / gridScale
		wy1 := 1 - b.gridOffsetY/gridScale
		wy2 := b.gridOffsetY / gridScale

		ax += s.vectorGrid[i1].x * wx1 * wy1
		ax += s.vectorGrid[i2].x * wx2 * wy1
		ax += s.vectorGrid[i3].x * wx1 * wy2
		ax += s.vectorGrid[i4].x * wx2 * wy2

		ay += s.vectorGrid[i1].y * wx1 * wy1
		ay += s.vectorGrid[i2].y * wx2 * wy1
		ay += s.vectorGrid[i3].y * wx1 * wy2
		ay += s.vectorGrid[i4].y * wx2 * wy2

		b.velocity.x += ax * timeStep
		b.velocity.y += ay * timeStep
	}
}

func (s *meshSimulator) calculateBodyLocation(timeStep float64) {
	for _, b := range s.bodies {
		b.location.x += b.velocity.x * timeStep
		b.location.y += b.velocity.y * timeStep

	}
}
