package pkg

func (state *BooleanFormulaState) Debug_CheckAssignmentIsSat() bool {
	for clauseIdx, clause := range state.Formula.Clauses {
		satisfiedClause := false
		for literal, asgn := range clause.Instances {
			if state.Assignments[literal] == asgn {
				// found one satisfying
				satisfiedClause = true
				break
			}
		}
		if !satisfiedClause {
			// curr assignments don't satisfy clause
			DebugLine("error: assignment doesn't satisfy", clauseIdx, state.Formula.Clauses[clauseIdx])
			return false
		}
	}
	// if it reaches here, every clause is satisfied, so it is SAT
	return true
}
