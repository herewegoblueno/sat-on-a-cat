package pkg

func (b *BooleanFormula) SolveFormula(initialState *BooleanFormulaState) (bool, *BooleanFormulaState) {
	//First assign everyone who has a sign larger than 1 some watched variables...
	for clauseIndex, clause := range b.Clauses {

		if len(clause.Instances) < 2 {
			continue
		}

		watchedLiteralHolder := WatchedLiterals{}
		//Ugly way to get two random keys...
		counter := 0
		for varIndex := range clause.Instances {
			counter++
			if counter == 1 {
				watchedLiteralHolder.right = varIndex
				initialState.VariablesKeepingTrackOfWhereTheyreBeingWatched[varIndex] = append(initialState.VariablesKeepingTrackOfWhereTheyreBeingWatched[varIndex], clauseIndex)
			} else if counter == 2 {
				watchedLiteralHolder.left = varIndex
				initialState.VariablesKeepingTrackOfWhereTheyreBeingWatched[varIndex] = append(initialState.VariablesKeepingTrackOfWhereTheyreBeingWatched[varIndex], clauseIndex)
			} else {
				break
			}
		}
		initialState.ClauseWatchedLiterals[clauseIndex] = watchedLiteralHolder
	}

	DebugLine("before solving state")
	PrintBooleanFormulaState(initialState)
	DebugLine("map of ClauseWatchedLiterals", initialState.ClauseWatchedLiterals)
	DebugLine("map of VariablesKeepingTrackOfWhereTheyreBeingWatched", initialState.VariablesKeepingTrackOfWhereTheyreBeingWatched)
	DebugLine("map of UnitClauses", initialState.UnitClauses)
	DebugLine("now solve", initialState.Sat)

	//Now solve the state...
	return initialState.SolveFromState()
}

func (state *BooleanFormulaState) SolveFromState() (bool, *BooleanFormulaState) {
	//Clear out the unit clauses first
	if !state.Sat {
		return false, state
	}

	err := state.UnitClauseElimination()
	if !state.Sat || err != nil {
		return false, nil
	}

	//Clear out the unit clauses first
	state.PureLiteralElimination()
	if !state.Sat {
		return false, nil
	}

	// if there are no more clauses, terminate before branching
	if len(state.DeletedClauses) == len(state.Formula.Clauses) {
		return true, state
	}

	// TODO: could make a better heuristic as opposed to arbitrarily picking a variable to branch on
	// loop over the variables in formula for branching
	for idx := range state.Formula.Vars {
		_, ok := state.Assignments[idx]

		if !ok {
			stateCopy := state.Copy()
			assignment := stateCopy.AssignmentFromDynamicLargestCombinedSum(idx)
			stateCopy.AssignmentPropagation(idx, assignment)
			solved, solvedState := stateCopy.SolveFromState()

			if solved {
				return solved, solvedState
			}

			stateCopy = state.Copy()
			stateCopy.AssignmentPropagation(idx, Negate(assignment))
			solved, solvedState = stateCopy.SolveFromState()

			if solved {
				return solved, solvedState
			}
		}

		//This loop will bottom out if we've gone branched on all the
		//unassigned variables and none of them lead to sat
	}

	return false, state
}

func (state *BooleanFormulaState) AssignmentFromDynamicLargestCombinedSum(variable VarIndex) VarState {
	posCount := 0
	negCount := 0
	for clauseIdx, clause := range state.Formula.Clauses {
		if _, ok := state.DeletedClauses[clauseIdx]; !ok {
			if clause.Instances[variable] == POS {
				posCount += 1
			} else {
				negCount += 1
			}
		}
	}
	if posCount > negCount {
		return POS
	} else {
		return NEG
	}
}
