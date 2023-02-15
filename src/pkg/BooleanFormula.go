package pkg

type VarIndex int
type ClauseIndex int
type VarState int
type InstanceState int

const (
	POS VarState = iota
	NEG
	DEFAULT //Only ever used in initial parsing
)

type SATClause struct {
	//--Immutable (after parsing finishes)
	Index     ClauseIndex
	Instances map[VarIndex]VarState
}

type SATVar struct {
	//--Immutable (after parsing finishes)
	Index             VarIndex
	ClauseAppearances map[ClauseIndex]VarState

	//These are only used for initial parsing
	LastSeenState VarState
}

type BooleanFormula struct {
	//--Immutable (after parsing finishes)
	Vars    map[VarIndex]*SATVar
	Clauses map[ClauseIndex]*SATClause
}

type WatchedLiterals struct {
	left  VarIndex
	right VarIndex
}

type BooleanFormulaState struct {
	//--Mutable and copied during branching
	Formula                                        *BooleanFormula
	Assignments                                    map[VarIndex]VarState
	ClauseWatchedLiterals                          map[ClauseIndex]WatchedLiterals //Won't contain unit clauses
	VariablesKeepingTrackOfWhereTheyreBeingWatched map[VarIndex][]ClauseIndex
	DeletedClauses                                 map[ClauseIndex]bool     //Bootleg set
	UnitClauses                                    map[ClauseIndex]VarIndex //Bootleg set
	PureVariables                                  map[VarIndex]VarState
	Sat                                            bool `default:true`
}
