package main

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"time"

	"github.com/jackdmoloney/orbits/sim"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Game struct {
	simulator    sim.Simulator
	layoutSize   int
	timeStep     float64
	stepsPerTick int
	timeElapsed  float64

	simulationTime string
}

func (g *Game) Update() error {
	timeStep := g.timeStep
	stepsPerTick := g.stepsPerTick
	now := time.Now()

	for i := 0; i < stepsPerTick; i++ {
		g.simulator.Step(timeStep)
	}
	g.timeElapsed += float64(stepsPerTick) * timeStep

	g.simulationTime = fmt.Sprint(time.Since(now))
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	bodies := g.simulator.Bodies()
	now := time.Now()

	for _, body := range bodies[:1000] {
		x, y := body.Location()
		radius := math.Sqrt(body.Mass())

		vector.DrawFilledCircle(
			screen,
			float32(x),
			float32(y),
			float32(radius),
			color.RGBA{A: 255, R: 255, G: 255, B: 255},
			false,
		)
	}

	ebitenutil.DebugPrint(screen, fmt.Sprintf(
		"TPS: %0.2f (Target: 60)\nSimulation Time: %s\nRender Time: %s\nBodies Simulated: %d\nBodies Rendered: %d",
		ebiten.ActualTPS(),
		g.simulationTime,
		time.Since(now),
		100000,
		1000,
	))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.layoutSize, g.layoutSize
}

func main() {
	layout := 400
	gridSize := 20
	simulator := sim.NewMeshSimulator(gridSize, layout/gridSize+1, layout/gridSize+1, 100000)

	game := Game{
		layoutSize:   layout,
		simulator:    &simulator,
		timeStep:     0.01,
		stepsPerTick: 1,
		timeElapsed:  0,
	}

	ebiten.SetWindowSize(1200, 1200)
	ebiten.SetWindowTitle("Orbits")
	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}

func SecondsToTime(seconds int) string {
	if seconds >= 31536000 {
		years := float64(seconds) / 31536000
		return fmt.Sprintf("%.2f years", years)
	} else if seconds >= 86400 {
		days := seconds / 86400
		return fmt.Sprintf("%d days", days)
	} else if seconds >= 3600 {
		hours := seconds / 3600
		return fmt.Sprintf("%d hours", hours)
	} else if seconds >= 60 {
		minutes := seconds / 60
		return fmt.Sprintf("%d minutes", minutes)
	} else {
		return fmt.Sprintf("%d seconds", seconds)
	}
}
