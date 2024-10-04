package main

import (
	"fmt"
	"log"

	"github.com/jackdmoloney/orbits/sim"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Game struct {
	simulator    sim.Simulator
	layoutSize   int
	timeStep     int
	stepsPerTick int
	timeElapsed  int
}

func (g *Game) Update() error {
	timeStep := g.timeStep
	stepsPerTick := g.stepsPerTick
	for i := 0; i < stepsPerTick; i++ {
		g.simulator.Step(float64(timeStep))
	}
	g.timeElapsed += stepsPerTick * timeStep

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	scale := float64(g.layoutSize) / 2
	for _, body := range g.simulator.Bodies {
		locationX, locationY := sim.LogScaledLocation(body, scale)
		vector.DrawFilledCircle(
			screen,
			float32(locationX+scale),
			float32(locationY+scale),
			float32(sim.LogScaledRadius(body, 2, 20)),
			sim.BodyColor(body.Name),
			false,
		)
	}

	ebitenutil.DebugPrint(screen, fmt.Sprintf(
		"TPS: %0.2f (Target: 60)\nTime Simulated: %s\nSimulated per Second: %s\nSimulation Time Step: %s\nSimulation Steps per Second: %d",
		ebiten.ActualTPS(),
		SecondsToTime(g.timeElapsed),
		SecondsToTime(g.timeStep*g.stepsPerTick*60),
		SecondsToTime(g.timeStep),
		g.stepsPerTick*60,
	))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.layoutSize, g.layoutSize
}

func main() {
	bodies := sim.BodiesSol()
	// bodies := sim.BodiesSunTest()
	simulator := sim.Simulator{Bodies: bodies}

	game := Game{
		layoutSize:   600,
		simulator:    simulator,
		timeStep:     60,
		stepsPerTick: 5000,
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
