package pkg

type VarState int
type VarIndex int

const (
	POS VarState = iota
	NEG
	MIX
)

type SATVar struct {
	index             VarIndex
	isPure            bool
	clauseAppearances map[int]VarState
}
