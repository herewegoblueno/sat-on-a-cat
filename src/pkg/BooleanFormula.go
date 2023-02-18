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
	Vars    map[VarIndex]*SATVar       //Immutable after parsing
	Clauses map[ClauseIndex]*SATClause //Immutable after parsing

	VarBranchingOrderOriginal        []VarIndex
	VarBranchingOrderShuffleDistance float64 //[0, 1]
	VarBranchingOrderShuffleChance   int     //[0, 100]

	BacktrackCounter              int
	BacktrackingLimit             int
	BacktrackingLimitIncreaseRate int
}

type WatchedLiterals struct {
	left  VarIndex
	right VarIndex
}

//--Mutable and copied during branching
type BooleanFormulaState struct {
	Sat     bool `default:true`
	Parent  *BooleanFormulaState
	Formula *BooleanFormula

	Depth                    int
	VarBranchingOrderLocal   []VarIndex
	VarBranchingOrderPointer *[]VarIndex //Can point to parent's VarBranchingOrderLocal, or own

	Assignments                                    map[VarIndex]VarState
	ClauseWatchedLiterals                          map[ClauseIndex]WatchedLiterals //Won't contain unit clauses
	VariablesKeepingTrackOfWhereTheyreBeingWatched map[VarIndex][]ClauseIndex
	DeletedClauses                                 map[ClauseIndex]bool     //Bootleg set
	UnitClauses                                    map[ClauseIndex]VarIndex //Bootleg set
	PureVariables                                  map[VarIndex]VarState
}
