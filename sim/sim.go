package sim

type Body interface {
	Location() (float64, float64)
	Mass() float64
}

type Simulator interface {
	Step(float64)
	Bodies() []Body
}
