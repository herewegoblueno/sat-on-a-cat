package pkg

type ClauseIndex int

type SATClause struct {
	index     ClauseIndex
	instances map[VarIndex]VarState
}
