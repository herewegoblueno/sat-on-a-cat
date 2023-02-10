package pkg

type VarIndex int
type ClauseIndex int
type VarState int
type InstanceState int

//Important that int(POS) == int (IPOS), same for INEG also
const (
	POS VarState = iota
	NEG
	DEFAULT //Only ever used in inital parsing
)

type SATClause struct {
	index     ClauseIndex
	instances map[VarIndex]VarState
}

type SATVar struct {
	index             VarIndex
	isPure            bool
	clauseAppearances map[ClauseIndex]VarState

	//These are only used for initial parsing
	lastSeenState VarState
}

type BooleanFormula struct {
	vars        map[VarIndex]SATVar
	clauses     map[ClauseIndex]SATClause
	sat         bool          `default:true` //TODO: does this work??
	unitClauses []ClauseIndex //Unit clauses
}

// vars: [2: SATVar{2, ?, [0: POS]}]
// clauses: [0: SATClause{0, [2: IPOS]}]
